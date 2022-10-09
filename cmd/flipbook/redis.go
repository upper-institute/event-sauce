package flipbook

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	redisdriver "github.com/upper-institute/event-sauce/internal/snapshotstore/redis"
)

var (
	redisShards []string
)

func redisSnapshotStoreBackend() error {

	addrs := map[string]string{}

	for i, address := range redisShards {
		addrs[fmt.Sprintf("shard%d", i)] = address
	}

	redisClient := redis.NewRing(&redis.RingOptions{
		Addrs: addrs,
	})

	snapshotStoreService.Backend = &redisdriver.RedisSnapshotStore{
		Redis: redisClient,
	}

	return nil

}

func init() {

	startCmd.PersistentFlags().StringArrayVar(&redisShards, "redisShards", []string{}, "Redis servers addresses list to build the Ring Cluster")

}
