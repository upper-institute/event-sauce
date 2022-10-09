package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisSnapshotStore struct {
	Redis redis.UniversalClient
}

func (p *RedisSnapshotStore) Set(ctx context.Context, id string, encoded []byte, ttl time.Duration) error {

	return p.Redis.Set(ctx, id, encoded, ttl).Err()

}

func (p *RedisSnapshotStore) Get(ctx context.Context, id string) ([]byte, error) {

	encoded, err := p.Redis.Get(ctx, id).Bytes()

	switch {

	case err == redis.Nil:
		return nil, nil

	case err != nil:
		return nil, err

	}

	return encoded, nil
}
