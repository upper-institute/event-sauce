package aws_dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	partitionKey_field   = "partition_key"
	sortingKey_field     = "sorting_key"
	eventPayload_field   = "event_payload"
	lastSortingKey_field = "last_sorting_key"
	metaSortingKey_value = int64(-1)
)

func unmarshalSortingKey(attrs map[string]types.AttributeValue) int64 {

	sortKey := int64(0)

	if err := attributevalue.Unmarshal(attrs[sortingKey_field], &sortKey); err != nil {
		panic(err)
	}

	return sortKey

}

func unmarshalEventPayload(attrs map[string]types.AttributeValue) *anypb.Any {

	evPayload := []byte{}

	if err := attributevalue.Unmarshal(attrs[eventPayload_field], evPayload); err != nil {
		panic(err)
	}

	any := &anypb.Any{}

	if err := proto.Unmarshal(evPayload, any); err != nil {
		panic(err)
	}

	return any

}

func unmarshalEvent(attrs map[string]types.AttributeValue) *flipbookv1.Event {

	return &flipbookv1.Event{
		SortingKey:   unmarshalSortingKey(attrs),
		EventPayload: unmarshalEventPayload(attrs),
	}

}

func marshalPartitionKey(src string) types.AttributeValue {

	partKey, err := attributevalue.Marshal(src)
	if err != nil {
		panic(err)
	}

	return partKey

}

func marshalSortingKey(src int64) types.AttributeValue {

	sortKey, err := attributevalue.Marshal(src)
	if err != nil {
		panic(err)
	}

	return sortKey

}

func marshalEvent(ev *flipbookv1.Event) map[string]types.AttributeValue {

	attrs := make(map[string]types.AttributeValue)

	attrs[partitionKey_field] = marshalPartitionKey(ev.PartitionKey)
	attrs[sortingKey_field] = marshalSortingKey(ev.SortingKey)

	evPayloadRaw, err := proto.Marshal(ev.EventPayload)
	if err != nil {
		panic(err)
	}

	evPayload, err := attributevalue.Marshal(evPayloadRaw)
	if err != nil {
		panic(err)
	}

	attrs[eventPayload_field] = evPayload

	return attrs

}
