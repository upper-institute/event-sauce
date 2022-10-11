package validation

import (
	"time"

	"github.com/google/uuid"
	apiv1 "github.com/upper-institute/flipbook/pkg/api/v1"
)

func IsValidID(id string) error {

	_, err := uuid.Parse(id)

	if err != nil {
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

	if len(events) == 0 {
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

type QueryBuilder interface {
	SetVersionRange(start int64, end int64, operator apiv1.QueryOperator)
	SetNaturalTimestampRange(start time.Time, end time.Time, operator apiv1.QueryOperator)
}

func IsValidScanRequest(req *apiv1.Event_ScanRequest, builder QueryBuilder) error {

	err := IsValidID(req.Id)

	if err != nil {
		return err
	}

	if req.Start == nil {
		return MissingStartQueryErr
	}

	if req.End == nil {
		return MissingEndQueryErr
	}

	switch start := req.Start.Parameter.(type) {
	case *apiv1.QueryParameter_Version:

		end, ok := req.End.Parameter.(*apiv1.QueryParameter_Version)

		if !ok {
			return InconsistentRangeErr
		}

		err = IsValidVersion(start.Version)

		if err != nil {
			return err
		}

		err = IsValidVersion(end.Version)

		if err != nil {
			return err
		}

		if start.Version > end.Version {
			return InvalidVersionRangeErr
		}

		builder.SetVersionRange(start.Version, end.Version, req.EndOperator)

	case *apiv1.QueryParameter_NaturalTimestamp:

		end, ok := req.End.Parameter.(*apiv1.QueryParameter_NaturalTimestamp)

		if !ok {
			return InconsistentRangeErr
		}

		if !start.NaturalTimestamp.IsValid() || end.NaturalTimestamp.IsValid() {
			return InvalidQueryNaturalTimestampErr
		}

		startTimestamp := start.NaturalTimestamp.AsTime()
		endTimestamp := end.NaturalTimestamp.AsTime()

		if startTimestamp.After(endTimestamp) {
			return InvalidNaturalTimestampRangeErr
		}

		builder.SetNaturalTimestampRange(startTimestamp, endTimestamp, req.EndOperator)

	}

	return nil

}
