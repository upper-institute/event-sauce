package eventstore

import (
	"context"

	"github.com/upper-institute/flipbook/internal/validation"
	apiv1 "github.com/upper-institute/flipbook/pkg/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Backend[Q validation.QueryBuilder] interface {
	Append(ctx context.Context, id string, version int64, events []*apiv1.Event) error
	GetQueryBuilder(id string) Q
	Scan(builder Q, stream apiv1.EventStore_ScanServer) error
	Latest(ctx context.Context, id string) (*apiv1.Event, error)
}

type EventStoreServer[Q validation.QueryBuilder] struct {
	apiv1.UnimplementedEventStoreServer

	Backend Backend[Q]
}

func (e *EventStoreServer[Q]) Append(ctx context.Context, req *apiv1.Event_AppendRequest) (*emptypb.Empty, error) {

	err := validation.IsValidEventSlice(req.Events)

	if err != nil {
		return nil, err
	}

	version := req.Events[0].Version
	id := req.Events[0].Id

	err = e.Backend.Append(ctx, id, version, req.Events)

	if err != nil {
		return nil, validation.FallbackGRPCError(err, validation.BackendAppendErr)
	}

	return &emptypb.Empty{}, nil
}

func (e *EventStoreServer[Q]) Scan(req *apiv1.Event_ScanRequest, stream apiv1.EventStore_ScanServer) error {

	builder := e.Backend.GetQueryBuilder(req.Id)

	err := validation.IsValidScanRequest(req, builder)

	if err != nil {
		return err
	}

	err = e.Backend.Scan(builder, stream)

	if err != nil {
		return validation.FallbackGRPCError(err, validation.BackendScanErr)
	}

	return nil

}

func (e *EventStoreServer[Q]) Latest(ctx context.Context, req *apiv1.Event_LatestRequest) (*apiv1.Event, error) {

	err := validation.IsValidID(req.Id)

	if err != nil {
		return nil, err
	}

	event, err := e.Backend.Latest(ctx, req.Id)

	if err != nil {
		return nil, validation.FallbackGRPCError(err, validation.BackendLatestErr)
	}

	if event == nil {
		return nil, validation.LatestNotFoundErr
	}

	return event, nil

}
