package drivers

import flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"

type BusDriver interface {
	Acquire(*flipbookv1.Subscription) error
}
