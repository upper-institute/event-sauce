package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/upper-institute/flipbook/internal/exceptions"
	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type redisStore struct {
	redis            redis.UniversalClient
	defaultBatchSize int64
}

func (r *redisStore) Write(ctx context.Context, req *flipbookv1.Event_AppendRequest, sem helpers.SortedEventMap) error {

	txf := func(tx *redis.Tx) error {

		for partKey, evs := range sem {

			firstEv := evs[0]

			if firstEv.SortingKeyType == flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE {

				hGetRes, err := tx.HGet(ctx, partKey, latestSortingKey_field).Result()

				switch {
				case err == redis.Nil:
					hGetRes = "0"
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

					raw, err := proto.Marshal(ev)
					if err != nil {
						panic(err)
					}

					sortKey := marshalInt64(ev.SortingKey)

					zAddArgs.Members = append(zAddArgs.Members, redis.Z{
						Score:  float64(ev.SortingKey),
						Member: sortKey,
					})

					hSetVals = append(hSetVals, sortKey, ev.EventPayload)

					if i == len(evs)-1 {

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
		Key:     req.PartitionKey,
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

		evPayloads, err := r.redis.HMGet(ctx, req.PartitionKey, sortKeys...).Result()
		if err != nil {
			return err
		}

		for i, sortKey := range sortKeys {

			evPayload := evPayloads[i].([]byte)

			ev := &flipbookv1.Event{
				PartitionKey: req.PartitionKey,
				SortingKey:   unmarshalInt64(sortKey),
				EventPayload: &anypb.Any{},
			}

			err = proto.Unmarshal(evPayload, ev.EventPayload)
			if err != nil {
				panic(err)
			}

			evCh <- ev

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
