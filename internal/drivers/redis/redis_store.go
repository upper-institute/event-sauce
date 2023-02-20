package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/upper-institute/flipbook/internal/exceptions"
	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type redisStore struct {
	redis            redis.UniversalClient
	defaultBatchSize int64
	log              *zap.SugaredLogger
}

func (r *redisStore) Write(ctx context.Context, req *flipbookv1.Event_AppendRequest, sem helpers.SortedEventMap) error {

	txf := func(tx *redis.Tx) error {

		for partKey, evs := range sem {

			log := r.log.With("partition_key", partKey)

			firstEv := evs[0]

			log.Debugw("Inserting events", "sorting_key_type", firstEv.SortingKeyType.String(), "first_event_sorting_key", firstEv.SortingKey)

			if firstEv.SortingKeyType == flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE {

				hGetRes, err := tx.HGet(ctx, partKey, latestSortingKey_field).Result()

				switch {
				case err == redis.Nil:
					hGetRes = "-1"
				case err != nil:
					return err
				}

				sortKey := unmarshalInt64(hGetRes)

				if err := exceptions.ThrowSortingKeyMismatchErr(firstEv.SortingKey, sortKey); err != nil {
					return err
				}

			}

			_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

				hSetVals := []interface{}{}

				zAddArgs := redis.ZAddArgs{
					NX:      false,
					Members: []redis.Z{},
				}

				if firstEv.SortingKeyType == flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE {
					zAddArgs.NX = true
				}

				for i, ev := range evs {

					sortKey := marshalInt64(ev.SortingKey)

					zAddArgs.Members = append(zAddArgs.Members, redis.Z{
						Score:  float64(ev.SortingKey),
						Member: sortKey,
					})

					evPayload, err := proto.Marshal(ev.EventPayload)
					if err != nil {
						panic(err)
					}

					hSetVals = append(hSetVals, sortKey, evPayload)

					if i == len(evs)-1 {

						raw, err := proto.Marshal(ev)
						if err != nil {
							panic(err)
						}

						hSetVals = append(
							hSetVals,
							latestEvent_field, raw,
							latestSortingKey_field, marshalInt64(ev.SortingKey),
						)
					}

				}

				pipe.HSet(ctx, partKey, hSetVals...)
				pipe.ZAddArgs(ctx, marshalSortedSetKey(partKey), zAddArgs)

				return nil

			})

			if err != nil {
				return err
			}

		}

		return nil

	}

	return r.redis.Watch(ctx, txf, GetKeysFromSortedEventMap(sem)...)

}

func (r *redisStore) Read(ctx context.Context, req *flipbookv1.Event_IterateRequest, evCh chan *flipbookv1.Event) error {

	zRangeArgs := redis.ZRangeArgs{
		Key:     marshalSortedSetKey(req.PartitionKey),
		ByScore: true,
		Count:   req.BatchSize,
		Start:   req.Query.StartSortingKey,
		Stop:    "+inf",
	}

	if req.Query.Stop == flipbookv1.QueryStop_QUERY_STOP_EXACT {
		zRangeArgs.Stop = req.Query.StopSortingKey
	}

	if zRangeArgs.Count <= 0 {
		zRangeArgs.Count = r.defaultBatchSize
	}

	for {

		sortKeys, err := r.redis.ZRangeArgs(ctx, zRangeArgs).Result()

		switch {
		case err == redis.Nil:
			return nil
		case err != nil:
			return err

		}

		if len(sortKeys) == 0 {
			break
		}

		r.log.Debugw("HMGet partition keys", "partition_key", req.PartitionKey, "sort_keys", sortKeys)

		evPayloads, err := r.redis.HMGet(ctx, req.PartitionKey, sortKeys...).Result()
		if err != nil {
			return err
		}

		if len(evPayloads) == 0 {
			break
		}

		for i, pld := range evPayloads {

			evPayload := []byte(pld.(string))

			sortKey := sortKeys[i]

			ev := &flipbookv1.Event{
				PartitionKey: req.PartitionKey,
				SortingKey:   unmarshalInt64(sortKey),
				EventPayload: &anypb.Any{},
			}

			zRangeArgs.Start = ev.SortingKey + 1

			err = proto.Unmarshal(evPayload, ev.EventPayload)
			if err != nil {
				panic(err)
			}

			select {
			case evCh <- ev:
				continue
			case <-ctx.Done():
				return ctx.Err()
			}

		}

	}

	return nil
}

func (r *redisStore) Tail(ctx context.Context, req *flipbookv1.Event_GetLatestRequest) (*flipbookv1.Event, error) {

	hGetRes, err := r.redis.HGet(ctx, req.PartitionKey, latestEvent_field).Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(hGetRes) == 0 {
		return nil, nil
	}

	ev := &flipbookv1.Event{}

	err = proto.Unmarshal([]byte(hGetRes), ev)
	if err != nil {
		return nil, err
	}

	return ev, nil

}
