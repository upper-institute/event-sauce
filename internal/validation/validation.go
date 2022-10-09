package validation

import (
	"github.com/google/uuid"
	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
)

func IsValidID(id string) error {

	_, err := uuid.Parse(id)

	if err == nil {
		return InvalidIDErr
	}

	return nil

}

func IsValidVersion(version int64) error {

	if version < 1 {
		return InvalidVersionErr
	}

	return nil

}

func IsValidEventSlice(events []*apiv1.Event) error {

	if events == nil || len(events) == 0 {
		return EmptyEventArrayErr
	}

	id := events[0].Id
	version := events[0].Version

	for _, event := range events {

		if err := IsValidVersion(event.Version); err != nil {
			return err
		}

		if version != event.Version {
			return InconsistentVersionErr
		}

		if err := IsValidID(event.Id); err != nil {
			return err
		}

		if id != event.Id {
			return InconsistentIDErr
		}

		version++

	}

	return nil

}

func IsValidSnapshot(snapshot *apiv1.Snapshot) error {

	err := IsValidID(snapshot.Id)

	if err != nil {
		return err
	}

	err = IsValidVersion(snapshot.Version)

	if err != nil {
		return err
	}

	return nil

}
