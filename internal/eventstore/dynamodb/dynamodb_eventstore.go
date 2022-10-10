package dynamodb

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/upper-institute/event-sauce/internal/validation"
	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DynamoDBEventStore struct {
	DynamoDB  *dynamodb.Client
	TableName string
}

var (
	zeroVersionAttr           = &types.AttributeValueMemberN{Value: "0"}
	latestProjection          = aws.String("latest")
	latestWithEventProjection = aws.String("latest, eventData")
	scanProjection            = aws.String("version, naturalTimestamp, storeTimestamp, payloadTypeUrl, payloadData")
)

func (e *DynamoDBEventStore) Append(ctx context.Context, id string, version int64, events []*apiv1.Event) error {

	tableName := aws.String(e.TableName)

	idAttr := &types.AttributeValueMemberS{Value: id}

	latest := int64(0)
	current := version - 1

	var err error

	if current == 0 {

		_, err = e.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: tableName,
			Item: map[string]types.AttributeValue{
				"id":        idAttr,
				"version":   zeroVersionAttr,
				"latest":    zeroVersionAttr,
				"eventData": &types.AttributeValueMemberNULL{Value: true},
			},
			ConditionExpression: aws.String("attribute_not_exists(id)"),
		})

	} else {
		latest, err = e.getLatestVersion(ctx, idAttr, nil)
	}

	if err != nil {
		return err
	}

	if current != latest {
		return validation.LatestVersionMismatchErr
	}

	latestEvent := events[len(events)-1]

	encodedLatest, err := proto.Marshal(latestEvent)

	if err != nil {
		return err
	}

	transactItems := []types.TransactWriteItem{
		{
			Update: &types.Update{
				Key: map[string]types.AttributeValue{
					"id":      idAttr,
					"version": zeroVersionAttr,
				},
				UpdateExpression:    aws.String("SET latest = :latest, eventData = :eventData"),
				TableName:           tableName,
				ConditionExpression: aws.String("latest = :current"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":current":   &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(latest), 10)},
					":eventData": &types.AttributeValueMemberB{Value: encodedLatest},
				},
			},
		},
	}

	now := &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339Nano)}

	var (
		latestAttr       *types.AttributeValueMemberN
		naturalTimestamp types.AttributeValue
	)

	for _, event := range events {

		latestAttr = &types.AttributeValueMemberN{Value: strconv.FormatInt(version, 10)}

		naturalTimestamp = &types.AttributeValueMemberNULL{Value: true}

		if event.NaturalTimestamp.IsValid() {
			naturalTimestamp = &types.AttributeValueMemberS{Value: event.NaturalTimestamp.AsTime().Format(time.RFC3339Nano)}
		}

		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName: tableName,
				Item: map[string]types.AttributeValue{
					"id":               idAttr,
					"version":          latestAttr,
					"naturalTimestamp": naturalTimestamp,
					"storeTimestamp":   now,
					"payloadTypeUrl":   &types.AttributeValueMemberS{Value: event.Payload.TypeUrl},
					"payloadData":      &types.AttributeValueMemberB{Value: event.Payload.Value},
				},
			},
		})

		version++

	}

	transactItems[0].Update.ExpressionAttributeValues[":latest"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(version, 10)}

	_, err = e.DynamoDB.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {

		switch err := err.(type) {
		case *smithy.OperationError:
			if strings.Contains(err.Error(), "ConditionalCheckFailed") {
				return validation.LatestVersionMismatchErr
			}
		}

	}

	return err
}

func unmarshalEvent(item map[string]types.AttributeValue) *apiv1.Event {

	event := &apiv1.Event{}

	var (
		versionAttr              = item["version"].(*types.AttributeValueMemberN)
		naturalTimestampAttr, ok = item["naturalTimestamp"].(*types.AttributeValueMemberS)
		storeTimestampAttr       = item["storeTimestamp"].(*types.AttributeValueMemberS)
		payloadTypeUrlAttr       = item["payloadTypeUrl"].(*types.AttributeValueMemberS)
		payloadDataAttr          = item["payloadData"].(*types.AttributeValueMemberB)
	)

	if ok {

		naturalTimestamp, err := time.Parse(time.RFC3339Nano, naturalTimestampAttr.Value)

		if err != nil {
			panic(err)
		}

		event.NaturalTimestamp = timestamppb.New(naturalTimestamp)

	}

	version, err := strconv.Atoi(versionAttr.Value)

	if err != nil {
		panic(err)
	}

	storeTimestamp, err := time.Parse(time.RFC3339Nano, storeTimestampAttr.Value)

	if err != nil {
		panic(err)
	}

	event.Version = int64(version)
	event.StoreTimestamp = timestamppb.New(storeTimestamp)
	event.Payload = &anypb.Any{
		TypeUrl: payloadTypeUrlAttr.Value,
		Value:   payloadDataAttr.Value,
	}

	return event

}

// func (e *DynamoDBEventStore) Scan(*apiv1.Event_ScanRequest, apiv1.EventStore_ScanServer) error {

// 	e.DynamoDB.Scan()

// }

func (e *DynamoDBEventStore) getLatestVersion(ctx context.Context, idAttr types.AttributeValue, event *apiv1.Event) (int64, error) {

	projection := latestProjection

	if event != nil {
		projection = latestWithEventProjection
	}

	res, err := e.DynamoDB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(e.TableName),
		Key: map[string]types.AttributeValue{
			"id":      idAttr,
			"version": zeroVersionAttr,
		},
		ProjectionExpression: projection,
	})

	if err != nil || res.Item == nil || len(res.Item) == 0 {
		return 0, err
	}

	latestAttr := res.Item["latest"].(*types.AttributeValueMemberN)

	latest, err := strconv.Atoi(latestAttr.Value)

	if err != nil {
		return 0, err
	}

	if latest > 0 && event != nil {

		encodedAttr := res.Item["latest"].(*types.AttributeValueMemberB)

		err = proto.Unmarshal(encodedAttr.Value, event)

		if err != nil {
			return 0, err
		}

	}

	return int64(latest), nil

}

func (e *DynamoDBEventStore) Latest(ctx context.Context, id string) (*apiv1.Event, error) {

	idAttr := &types.AttributeValueMemberS{Value: id}

	event := &apiv1.Event{}

	latest, err := e.getLatestVersion(ctx, idAttr, event)

	if err != nil || latest == 0 {
		return nil, err
	}

	return event, nil

}
