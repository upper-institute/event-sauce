package internal

import (
	"context"
	"io"

	"github.com/upper-institute/flipbook/internal/drivers"
	"github.com/upper-institute/flipbook/internal/exceptions"
	"github.com/upper-institute/flipbook/internal/helpers"
	"github.com/upper-institute/flipbook/internal/logging"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type eventStore struct {
	flipbookv1.UnimplementedEventStoreServer
	driver drivers.StoreDriver
	log    *zap.SugaredLogger
}

func NewEventStore(driver drivers.StoreDriver) flipbookv1.EventStoreServer {
	return &eventStore{
		driver: driver,
		log:    logging.Logger.Sugar().Named("EventStore"),
	}
}

func (s *eventStore) Append(ctx context.Context, req *flipbookv1.Event_AppendRequest) (*emptypb.Empty, error) {

	sem := helpers.NewSortedEventMap(req.Events)

	if err := exceptions.ThrowInvalidAppendRequestErr(req, sem); err != nil {
		return nil, err
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
			s.log.Debug("IO EOF for event store iterate request")
			return nil
		case err != nil:
			return err
		}

		if err := exceptions.ThrowInvalidIterateRequestErr(req); err != nil {
			return err
		}

		evCh := make(chan *flipbookv1.Event)
		endCh := make(chan error)

		go func() {
			endCh <- s.driver.Read(context.Background(), req, evCh)
			close(endCh)
			close(evCh)
		}()

	read:
		for {

			select {

			case ev := <-evCh:
				if ev == nil {
					s.log.Debug("Event is nil")
					break read
				}
				if err = srv.Send(ev); err != nil {
					s.log.Debug("Send event")
					break read
				}

			case readErr := <-endCh:
				s.log.Debug("End channel signaled")
				if readErr != nil {
					err = exceptions.ThrowDriverReadErr(readErr)
				}
				break read

			}

		}
		if err != nil {
			return err
		}

		err = srv.Send(&flipbookv1.Event{SortingKey: -1})
		if err != nil {
			return err
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
