package internal

import (
	"github.com/upper-institute/flipbook/internal/drivers"
	"github.com/upper-institute/flipbook/internal/drivers/aws_dynamodb"
	"github.com/upper-institute/flipbook/internal/drivers/postgres"
	"github.com/upper-institute/flipbook/internal/drivers/redis"
	"github.com/upper-institute/flipbook/internal/helpers"
)

const (
	storeDriver_flag = "store.driver"
)

var wellKnownStoreDrivers = drivers.StoreDriverMap{
	"redis":        &redis.RedisStoreAdapter{},
	"aws_dynamodb": &aws_dynamodb.AwsDynamodbStoreAdapter{},
	"postgres":     &postgres.PostgresStoreAdapter{},
}

func BindWellKnownDrivers(binder helpers.FlagBinder) {

	for _, adapter := range wellKnownStoreDrivers {
		adapter.Bind(binder)
	}

	binder.BindString(storeDriver_flag, "", "Define which driver used for the persistence layer of event store")

}

func GetStoreAdapter(getter helpers.FlagGetter) drivers.StoreDriverAdapter {

	storeDriver := getter.GetString(storeDriver_flag)

	adapter, ok := wellKnownStoreDrivers[storeDriver]

	if !ok {
		return nil
	}

	return adapter

}
