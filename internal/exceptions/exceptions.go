package exceptions

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/upper-institute/flipbook/internal/helpers"
	"github.com/upper-institute/flipbook/internal/logging"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mustMarshalJson(src map[string]interface{}) string {

	raw, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}

	return string(raw)

}

func ThrowInternalErr(err error) error {

	if _, ok := status.FromError(err); ok {
		return err
	}

	errHash := md5.New()
	errHash.Write([]byte(err.Error()))
	hash := hex.EncodeToString(errHash.Sum(nil))

	log := logging.Logger.Named("exceptions").Sugar()

	log.Errorw("Internal error", "hash", hash, "error", err)

	return status.Error(codes.Internal, fmt.Sprintf(`{"internalError":"%s"}`, hash))

}

func ThrowInvalidAppendRequestErr(req *flipbookv1.Event_AppendRequest, sem helpers.SortedEventMap) error {

	errMap := map[string]interface{}{}

	sortingKeyType := flipbookv1.SortingKeyType(-1)

validation:
	for partKey, evs := range sem {

		helpers.SortEventSlice(evs)

		_, err := uuid.Parse(partKey)
		if err != nil {
			errMap["invalidPartitionKey"] = partKey
			break
		}

		if sortingKeyType == -1 {
			sortingKeyType = evs[0].SortingKeyType
		}

		for i, ev := range evs {

			if ev.SortingKey < 0 {
				errMap["invalidSortingKey"] = partKey
				break validation
			}

			if sortingKeyType == flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE && i > 0 {

				if ev.SortingKey-1 != evs[i-1].SortingKey {
					errMap["invalidIncreasingSequence"] = partKey
					break validation
				}

			}

			if ev.SortingKeyType != sortingKeyType {
				errMap["inconsistentSortingKeyType"] = partKey
				break validation
			}

		}

	}

	if len(errMap) == 0 {
		return nil
	}

	return status.Errorf(codes.InvalidArgument, mustMarshalJson(errMap))

}

func ThrowInvalidGetLatestRequestErr(req *flipbookv1.Event_GetLatestRequest) error {

	_, err := uuid.Parse(req.PartitionKey)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, `{"invalidPartitionKey":"%s"}`, req.PartitionKey)
	}

	return nil
}

func ThrowInvalidIterateRequestErr(req *flipbookv1.Event_IterateRequest) error {

	_, err := uuid.Parse(req.PartitionKey)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, `{"invalidPartitionKey":"%s"}`, req.PartitionKey)
	}

	if req.Query == nil {
		return status.Errorf(codes.InvalidArgument, `{"missingQuery":true}`)
	}

	if req.Query.StartSortingKey < 0 {
		return status.Errorf(codes.InvalidArgument, `{"invalidStartSortingKey":%d}`, req.Query.StartSortingKey)
	}

	if req.Query.Stop == flipbookv1.QueryStop_QUERY_STOP_EXACT && req.Query.StartSortingKey > req.Query.StopSortingKey {
		return status.Errorf(codes.InvalidArgument, `{"startGreaterThanStop":true}`)
	}

	return nil
}

func ThrowPartitionNotFoundErr(partKey string) error {
	return status.Errorf(codes.NotFound, `{"notFoundPartitionKey":"%s"}`, partKey)
}

func ThrowSortingKeyMismatchErr(providedSortKey int64, currentSortKey int64) error {

	if providedSortKey == currentSortKey+1 {
		return nil
	}

	return status.Errorf(codes.FailedPrecondition, `{"providedSortingKey":%d,"expectedSortingKey":%d}`, providedSortKey, currentSortKey)
}

func ThrowOptimisticLockFailedErr() error {
	return status.Errorf(codes.Aborted, `{"optimisticLockFailed":true}`)
}

func ThrowDriverWriteErr(err error) error {
	return ThrowInternalErr(err)
}

func ThrowDriverReadErr(err error) error {
	return ThrowInternalErr(err)
}

func ThrowDriverTailErr(err error) error {
	return ThrowInternalErr(err)
}
