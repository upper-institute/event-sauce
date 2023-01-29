package redis

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/upper-institute/flipbook/internal/drivers"
	"github.com/upper-institute/flipbook/internal/helpers"
)

const (
	redisUrls_flag    = "redis.urls"
	readTimeout_flag  = "redis.read.timeout"
	writeTimeout_flag = "redis.write.timeout"
	batchSize_flag    = "redis.batch.size"
)

type RedisStoreAdapter struct {
}

func (r *RedisStoreAdapter) Bind(binder helpers.FlagBinder) {

	binder.BindStringSlice(redisUrls_flag, []string{"redis://localhost:6379?db=0"}, "Redis servers to connect (if more than one, it'll be a ring)")
	binder.BindDuration(readTimeout_flag, 10*time.Second, "Timeout for read operations")
	binder.BindDuration(writeTimeout_flag, 10*time.Second, "Timeout for write operations")
	binder.BindInt64(batchSize_flag, 50, "Default batch size to scan from Redis")

}

func (r *RedisStoreAdapter) New(getter helpers.FlagGetter) (drivers.StoreDriver, error) {

	redisUrls := getter.GetStringSlice(redisUrls_flag)
	urls := []*url.URL{}

	for _, redisUrlStr := range redisUrls {

		redisUrl, err := url.Parse(redisUrlStr)
		if err != nil {
			return nil, err
		}

		urls = append(urls, redisUrl)

	}

	if len(urls) == 0 {
		panic("Redis store initialized with zero urls, this should never happen")
	}

	firstUrl := urls[0]

	store := &redisStore{
		defaultBatchSize: getter.GetInt64(batchSize_flag),
	}

	password, db, err := getPasswordAndDbFromUrl(firstUrl)
	if err != nil {
		return nil, err
	}

	if len(urls) == 1 {
		store.redis = redis.NewClient(&redis.Options{
			Addr:         firstUrl.Host,
			Username:     firstUrl.User.Username(),
			Password:     password,
			MaxRetries:   -1,
			DB:           int(db),
			ReadTimeout:  getter.GetDuration(readTimeout_flag),
			WriteTimeout: getter.GetDuration(writeTimeout_flag),
		})
	}

	if len(urls) > 1 {

		addrs := make(map[string]string)

		for _, redisUrl := range urls {

			ensurePortOnUrl(redisUrl)

			addrs[redisUrl.Query().Get("nodeId")] = redisUrl.Host

		}

		store.redis = redis.NewRing(&redis.RingOptions{
			Addrs:        addrs,
			Username:     firstUrl.User.Username(),
			Password:     password,
			MaxRetries:   -1,
			DB:           int(db),
			ReadTimeout:  getter.GetDuration(readTimeout_flag),
			WriteTimeout: getter.GetDuration(writeTimeout_flag),
		})
	}

	return store, nil
}

func (r *RedisStoreAdapter) Destroy(store drivers.StoreDriver) error {
	return nil
}

func ensurePortOnUrl(redisUrl *url.URL) {
	if len(redisUrl.Port()) == 0 {
		redisUrl.Host = fmt.Sprintf("%s:%d", redisUrl.Host, 6379)
	}
}

func getPasswordAndDbFromUrl(redisUrl *url.URL) (string, int, error) {

	ensurePortOnUrl(redisUrl)

	password, _ := redisUrl.User.Password()
	dbStr := redisUrl.Query().Get("db")

	if len(dbStr) == 0 {
		dbStr = "0"
	}

	db, err := strconv.ParseInt(dbStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("Invalid Redis db parameter in URL: %s", dbStr)
	}

	return password, int(db), nil

}
