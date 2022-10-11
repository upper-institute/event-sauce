package snapshotstore

import (
	"context"
	"time"

	"github.com/upper-institute/flipbook/internal/validation"
	apiv1 "github.com/upper-institute/flipbook/pkg/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Backend interface {
	Set(ctx context.Context, id string, encoded []byte, ttl time.Duration) error
	Get(ctx context.Context, id string) ([]byte, error)
}

type SnapshotStoreServer struct {
	apiv1.UnimplementedSnapshotStoreServer

	Backend Backend
}

func (s *SnapshotStoreServer) Set(ctx context.Context, snapshot *apiv1.Snapshot) (*emptypb.Empty, error) {

	err := validation.IsValidSnapshot(snapshot)

	if err != nil {
		return nil, err
	}

	encodedSnapshot, err := proto.Marshal(snapshot)

	if err != nil {
		return nil, validation.ProtoMarshalErr
	}

	if !snapshot.Ttl.IsValid() {
		snapshot.Ttl = durationpb.New(0)
	}

	err = s.Backend.Set(ctx, snapshot.Id, encodedSnapshot, snapshot.Ttl.AsDuration())

	if err != nil {
		return nil, validation.FallbackGRPCError(err, validation.BackendSetErr)
	}

	return &emptypb.Empty{}, nil

}

func (s *SnapshotStoreServer) Get(ctx context.Context, req *apiv1.Snapshot_GetRequest) (*apiv1.Snapshot, error) {

	err := validation.IsValidID(req.Id)

	if err != nil {
		return nil, err
	}

	encoded, err := s.Backend.Get(ctx, req.Id)

	if err != nil {
		return nil, validation.FallbackGRPCError(err, validation.BackendGetErr)
	}

	if encoded == nil {
		return nil, validation.SnapshotNotFoundErr
	}

	snapshot := &apiv1.Snapshot{}

	err = proto.Unmarshal(encoded, snapshot)

	if err != nil {
		return nil, err
	}

	return snapshot, nil

}
