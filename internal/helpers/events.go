package helpers

import (
	"sort"

	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
)

func SortEventSlice(evs []*flipbookv1.Event) {
	sort.SliceStable(evs[:], func(i, j int) bool {
		return evs[i].SortingKey < evs[j].SortingKey
	})
}
