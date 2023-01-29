package helpers

import (
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
)

type SortedEventMap map[string][]*flipbookv1.Event

func NewSortedEventMap(evs []*flipbookv1.Event) SortedEventMap {

	s := make(SortedEventMap)

	for _, ev := range evs {
		s.Set(ev)
	}

	return s

}

func (s SortedEventMap) GetKeys() []string {

	keys := []string{}

	for key := range s {
		keys = append(keys, key)
	}

	return keys

}

func (s SortedEventMap) Set(ev *flipbookv1.Event) {

	if _, ok := s[ev.PartitionKey]; !ok {
		s[ev.PartitionKey] = make([]*flipbookv1.Event, 0)
	}

	s[ev.PartitionKey] = append(s[ev.PartitionKey], ev)

}

func (s SortedEventMap) GetByPartitionKey(partKey string) []*flipbookv1.Event {

	evs, ok := s[partKey]

	if !ok {
		return nil
	}

	return evs

}
