package internal

import (
	"context"
	"io"

	"github.com/upper-institute/flipbook/internal/drivers"
	"github.com/upper-institute/flipbook/internal/exceptions"
	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type eventStore struct {
	flipbookv1.UnimplementedEventStoreServer
	driver drivers.StoreDriver
}

func NewEventStore(driver drivers.StoreDriver) flipbookv1.EventStoreServer {
	return &eventStore{
		driver: driver,
	}
}

func (s *eventStore) Append(ctx context.Context, req *flipbookv1.Event_AppendRequest) (*emptypb.Empty, error) {

	if err := exceptions.ThrowInvalidAppendRequestErr(req); err != nil {
		return nil, err
	}

	sem := helpers.NewSortedEventMap(req.Events)

	for _, evs := range sem {
		helpers.SortEventSlice(evs)
	}

	if err := exceptions.ThrowDriverWriteErr(s.driver.Write(ctx, req, sem)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}

func (s *eventStore) Iterate(srv flipbookv1.EventStore_IterateServer) error {

	ctx := srv.Context()

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		if err := exceptions.ThrowInvalidIterateRequestErr(req); err != nil {
			return err
		}

		evCh := make(chan *flipbookv1.Event)
		endCh := make(chan error)

		go func() { endCh <- s.driver.Read(ctx, req, evCh) }()

	read:
		for {

			select {

			case ev := <-evCh:
				if err = srv.Send(ev); err != nil {
					break read
				}

			case readErr := <-endCh:
				if readErr != nil {
					err = exceptions.ThrowDriverReadErr(readErr)
				}
				break read

			}

		}

		close(endCh)

		if evCh != nil {
			close(evCh)
		}

		if err != nil {
			return err
		}

		if req.BatchSize > 0 {
			return nil
		}

	}

}

func (s *eventStore) GetLatest(ctx context.Context, req *flipbookv1.Event_GetLatestRequest) (*flipbookv1.Event, error) {

	if err := exceptions.ThrowInvalidGetLatestRequestErr(req); err != nil {
		return nil, err
	}

	ev, err := s.driver.Tail(ctx, req)

	if err != nil {
		return nil, exceptions.ThrowDriverTailErr(err)
	}

	if ev == nil {
		return nil, exceptions.ThrowPartitionNotFoundErr(req.PartitionKey)
	}

	return ev, nil

}
