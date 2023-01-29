package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/upper-institute/flipbook/internal/drivers/postgres/database"
	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type postgresStore struct {
	db               *sql.DB
	queries          *database.Queries
	defaultBatchSize int32
}

func (p *postgresStore) Write(ctx context.Context, req *flipbookv1.Event_AppendRequest, sem helpers.SortedEventMap) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	for partKeyStr, evs := range sem {

		partKey := uuid.MustParse(partKeyStr)

		firstEv := evs[0]

		for _, ev := range evs {

			raw, err := proto.Marshal(ev)
			if err != nil {
				panic(err)
			}

			switch firstEv.SortingKeyType {

			case flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE:

				err = qtx.InsertEventSequence(ctx, database.InsertEventSequenceParams{
					PartitionKey: partKey,
					SortingKey:   ev.SortingKey,
					EventPayload: raw,
				})

			case flipbookv1.SortingKeyType_SORTING_KEY_ARBITRARY_NUMBER:

				err = qtx.UpsertEvent(ctx, database.UpsertEventParams{
					PartitionKey: partKey,
					SortingKey:   ev.SortingKey,
					EventPayload: raw,
				})

			}

			if err != nil {
				return err
			}

		}

	}

	return tx.Commit()

}

func (p *postgresStore) Read(ctx context.Context, req *flipbookv1.Event_IterateRequest, evCh chan *flipbookv1.Event) error {

	partKey := uuid.MustParse(req.PartitionKey)

	batchSize := int32(req.BatchSize)

	if batchSize <= 0 {
		batchSize = p.defaultBatchSize
	}

	startSortKey := req.Query.StartSortingKey

	for {

		evsCount := 0

		switch req.Query.Stop {

		case flipbookv1.QueryStop_QUERY_STOP_EXACT:

			rows, err := p.queries.RangeEventsWithStop(ctx, database.RangeEventsWithStopParams{
				PartitionKey:    partKey,
				StartSortingKey: startSortKey,
				StopSortingKey:  req.Query.StopSortingKey,
				Limit:           batchSize,
			})
			if err != nil {
				return err
			}

			evsCount = len(rows)

			for _, row := range rows {

				pld := &anypb.Any{}

				if err := proto.Unmarshal(row.EventPayload, pld); err != nil {
					return err
				}

				startSortKey = row.SortingKey + 1

				evCh <- &flipbookv1.Event{
					PartitionKey: req.PartitionKey,
					SortingKey:   row.SortingKey,
					EventPayload: pld,
				}

			}

			if startSortKey >= req.Query.StopSortingKey {
				return nil
			}

		case flipbookv1.QueryStop_QUERY_STOP_LATEST:

			rows, err := p.queries.RangeEvents(ctx, database.RangeEventsParams{
				PartitionKey:    partKey,
				StartSortingKey: startSortKey,
				Limit:           batchSize,
			})
			if err != nil {
				return err
			}

			evsCount = len(rows)

			for _, row := range rows {

				pld := &anypb.Any{}

				if err := proto.Unmarshal(row.EventPayload, pld); err != nil {
					return err
				}

				startSortKey = row.SortingKey + 1

				evCh <- &flipbookv1.Event{
					PartitionKey: req.PartitionKey,
					SortingKey:   row.SortingKey,
					EventPayload: pld,
				}

			}

		}

		if evsCount == 0 {
			break
		}

	}

	return nil
}

func (p *postgresStore) Tail(ctx context.Context, req *flipbookv1.Event_GetLatestRequest) (*flipbookv1.Event, error) {

	partKey := uuid.MustParse(req.PartitionKey)

	dbEv, err := p.queries.GetLastEvent(ctx, partKey)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	ev := &flipbookv1.Event{
		PartitionKey: dbEv.PartitionKey.String(),
		SortingKey:   dbEv.SortingKey,
		EventPayload: &anypb.Any{},
	}

	if err := proto.Unmarshal(dbEv.EventPayload, ev.EventPayload); err != nil {
		panic(err)
	}

	return ev, nil
}
