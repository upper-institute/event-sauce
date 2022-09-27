package cmd

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	redisdriver "github.com/upper-institute/event-sauce/internal/redis"
	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
)

var (
	redisShards []string

	redisCmd = &cobra.Command{
		Use:   "redis",
		Short: "Use Redis as a Snapshot Store",
		RunE: func(cmd *cobra.Command, args []string) error {

			addrs := map[string]string{}

			for i, address := range redisShards {
				addrs[fmt.Sprintf("shard%d", i)] = address
			}

			rdb := redis.NewRing(&redis.RingOptions{
				Addrs: addrs,
			})

			paprikaService := &redisdriver.PaprikaServiceServer{
				Redis: rdb,
			}

			apiv1.RegisterPaprikaServiceServer(grpcServer, paprikaService)

			serveGrpcServer()

			return nil

		},
	}
)

func init() {

	rootCmd.AddCommand(redisCmd)

	redisCmd.PersistentFlags().StringArrayVar(&redisShards, "redisShards", []string{}, "Redis servers addresses list to build the Ring Cluster")

}
