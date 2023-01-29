package redis

import (
	"strconv"
	"strings"

	"github.com/upper-institute/flipbook/internal/helpers"
)

const (
	latestEvent_field      = "lev"
	latestSortingKey_field = "lsk"
	sorteSetKey_suffix     = "-sst"
)

func unmarshalInt64(src string) int64 {

	t, err := strconv.ParseInt(src, 10, 32)
	if err != nil {
		panic(err)
	}

	return t
}

func marshalInt64(src int64) string {
	return strconv.FormatInt(src, 10)
}

func marshalSortedSetKey(partKey string) string {
	key := strings.Builder{}
	key.WriteString(partKey)
	key.WriteString(sorteSetKey_suffix)
	return key.String()
}

func GetKeysFromSortedEventMap(sem helpers.SortedEventMap) []string {

	keys := sem.GetKeys()
	sstKeys := []string{}

	for _, key := range keys {
		sstKeys = append(sstKeys, marshalSortedSetKey(key))
	}

	return append(keys, sstKeys...)

}
