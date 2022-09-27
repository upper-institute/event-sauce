package redis

import (
	"context"

	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-redis/redis/v8"
)

type PaprikaServiceServer struct {
	apiv1.UnimplementedPaprikaServiceServer

	Redis redis.UniversalClient
}

func (p *PaprikaServiceServer) Set(ctx context.Context, snapshot *apiv1.Snapshot) (*emptypb.Empty, error) {

	encodedSnapshot, err := proto.Marshal(snapshot)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while marshaling snapshot to protobuf: %v", err)
	}

	err = p.Redis.Set(ctx, snapshot.Id, encodedSnapshot, snapshot.Ttl.AsDuration()).Err()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "redis set error: %v", err)
	}

	return &emptypb.Empty{}, nil

}

func (p *PaprikaServiceServer) Get(ctx context.Context, getRequest *apiv1.Snapshot_GetRequest) (*apiv1.Snapshot, error) {

	encodedSnapshot, err := p.Redis.Get(ctx, getRequest.Id).Bytes()

	switch {

	case err == redis.Nil:
		return nil, status.Error(codes.NotFound, "snapshot not found")

	case err != nil:
		return nil, status.Errorf(codes.Internal, "redis get: %v", err)

	}

	snapshot := &apiv1.Snapshot{}

	err = proto.Unmarshal(encodedSnapshot, snapshot)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while unmarshaling protobuf to snapshot: %v", err)
	}

	return snapshot, nil
}
