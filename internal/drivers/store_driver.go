package drivers

import (
	"context"

	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
)

type StoreDriver interface {
	Write(context.Context, *flipbookv1.Event_AppendRequest, helpers.SortedEventMap) error
	Read(context.Context, *flipbookv1.Event_IterateRequest, chan *flipbookv1.Event) error
	Tail(context.Context, *flipbookv1.Event_GetLatestRequest) (*flipbookv1.Event, error)
}

type StoreDriverAdapter interface {
	Bind(helpers.FlagBinder)
	New(helpers.FlagGetter) (StoreDriver, error)
	Destroy(StoreDriver) error
}

type StoreDriverMap map[string]StoreDriverAdapter
