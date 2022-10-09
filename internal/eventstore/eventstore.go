package eventstore

import (
	"context"

	"github.com/upper-institute/event-sauce/internal/validation"
	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Backend interface {
	Append(ctx context.Context, id string, version int64, events []*apiv1.Event) error
	Latest(ctx context.Context, id string) (*apiv1.Event, error)
}

type EventStoreServer struct {
	apiv1.UnimplementedEventStoreServer

	Backend Backend
}

func (e *EventStoreServer) Append(ctx context.Context, req *apiv1.Event_AppendRequest) (*emptypb.Empty, error) {

	err := validation.IsValidEventSlice(req.Events)

	if err != nil {
		return nil, err
	}

	version := req.Events[0].Version
	id := req.Events[0].Id

	err = e.Backend.Append(ctx, id, version, req.Events)

	if err != nil {

		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		return nil, validation.BackendAppendErr
	}

	return &emptypb.Empty{}, nil
}

// func (e *EventStoreServer) Scan(req *apiv1.Event_ScanRequest, stream apiv1.EventStore_ScanServer) error {

// 	// e.DynamoDB.Scan()

// }

func (e *EventStoreServer) Latest(ctx context.Context, req *apiv1.Event_LatestRequest) (*apiv1.Event, error) {

	err := validation.IsValidID(req.Id)

	if err != nil {
		return nil, err
	}

	event, err := e.Backend.Latest(ctx, req.Id)

	if err != nil {

		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		return nil, validation.BackendLatestErr

	}

	if event == nil {
		return nil, validation.LatestNotFoundErr
	}

	return event, nil

}
