package redis

import (
	"context"

	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-redis/redis/v8"
)

type PaprikaServiceServer struct {
	apiv1.UnimplementedPaprikaServiceServer

	Redis redis.UniversalClient
}

func (p *PaprikaServiceServer) Set(context.Context, *apiv1.Snapshot) (*apiv1.Snapshot_SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (p *PaprikaServiceServer) Get(context.Context, *apiv1.Snapshot_GetRequest) (*apiv1.Snapshot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
