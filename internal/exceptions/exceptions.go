package exceptions

import (
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ThrowInvalidAppendRequestErr(req *flipbookv1.Event_AppendRequest) error {
	return status.Errorf(codes.InvalidArgument, "")
}

func ThrowInvalidGetLatestRequestErr(req *flipbookv1.Event_GetLatestRequest) error {
	return status.Errorf(codes.InvalidArgument, "")
}

func ThrowInvalidIterateRequestErr(req *flipbookv1.Event_IterateRequest) error {
	return status.Errorf(codes.InvalidArgument, "")
}

func ThrowDriverWriteErr(err error) error {
	return err
}

func ThrowPartitionNotFoundErr(partKey string) error {
	return status.Errorf(codes.NotFound, "")
}

func ThrowSortingKeyMismatchErr(providedSortKey int64, currenSortKey int64) error {

	if providedSortKey == currenSortKey+1 {
		return nil
	}

	return status.Errorf(codes.FailedPrecondition, "")
}

func ThrowOptimisticLockFailedErr() error {
	return status.Errorf(codes.Aborted, "")
}

func ThrowDriverReadErr(err error) error {
	return err
}

func ThrowDriverTailErr(err error) error {
	return err
}
