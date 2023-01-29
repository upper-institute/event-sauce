package aws_dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/upper-institute/flipbook/internal/exceptions"
	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
)

type awsDynamodbStore struct {
	dynamo           *dynamodb.Client
	tableName        *string
	defaultBatchSize int64
}

func (d *awsDynamodbStore) Write(ctx context.Context, req *flipbookv1.Event_AppendRequest, sem helpers.SortedEventMap) error {

	putItens := []types.TransactWriteItem{}
	updateItens := []types.TransactWriteItem{}

	for partKey, evs := range sem {

		firstEv := evs[0]
		lastEv := evs[len(evs)-1]

		updateExpr := expression.Set(
			expression.Name(lastSortingKey_field),
			expression.Value(marshalSortingKey(lastEv.SortingKey)),
		)

		exprBuilder := expression.NewBuilder()

		if firstEv.SortingKeyType == flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE {

			condExpr := expression.Equal(
				expression.Name(lastSortingKey_field),
				expression.Value(marshalSortingKey(firstEv.SortingKey-1)),
			)

			exprBuilder = exprBuilder.WithCondition(condExpr)

		}

		expr, err := exprBuilder.
			WithUpdate(updateExpr).
			Build()
		if err != nil {
			return err
		}

		updateItens = append(updateItens, types.TransactWriteItem{
			Update: &types.Update{
				TableName: d.tableName,
				Key: map[string]types.AttributeValue{
					partitionKey_field: marshalPartitionKey(partKey),
					sortingKey_field:   marshalSortingKey(metaSortingKey_value),
				},
				UpdateExpression:          expr.Update(),
				ConditionExpression:       expr.Condition(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
			},
		})

		for _, ev := range evs {

			putItens = append(putItens, types.TransactWriteItem{
				Put: &types.Put{
					TableName: d.tableName,
					Item:      marshalEvent(ev),
				},
			})

		}

	}

	_, err := d.dynamo.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: append(putItens, updateItens...),
	})

	if err != nil {

		dynamoErr := &types.ConditionalCheckFailedException{}

		if errors.As(err, dynamoErr) {
			return exceptions.ThrowOptimisticLockFailedErr()
		}

		return err

	}

	return nil

}

func (d *awsDynamodbStore) Read(ctx context.Context, req *flipbookv1.Event_IterateRequest, evCh chan *flipbookv1.Event) error {

	batchSize := req.BatchSize

	if batchSize <= 0 {
		batchSize = d.defaultBatchSize
	}

	partKey := req.PartitionKey
	startSortKey := req.Query.StartSortingKey

	projEx := expression.NamesList(
		expression.Name(sortingKey_field),
		expression.Name(eventPayload_field),
	)

	var exclusiveStartKey map[string]types.AttributeValue = nil

	for {

		keyEx := expression.Key(partitionKey_field).Equal(expression.Value(partKey))

		switch req.Query.Stop {

		case flipbookv1.QueryStop_QUERY_STOP_EXACT:
			keyEx = keyEx.And(expression.KeyBetween(
				expression.Key(sortingKey_field),
				expression.Value(startSortKey),
				expression.Value(req.Query.StopSortingKey),
			))

		case flipbookv1.QueryStop_QUERY_STOP_LATEST:
			keyEx = keyEx.And(expression.KeyGreaterThanEqual(
				expression.Key(sortingKey_field),
				expression.Value(startSortKey),
			))
		}

		expr, err := expression.NewBuilder().
			WithKeyCondition(keyEx).
			WithProjection(projEx).
			Build()
		if err != nil {
			return err
		}

		queryRes, err := d.dynamo.Query(ctx, &dynamodb.QueryInput{
			ProjectionExpression:      expr.Projection(),
			ConsistentRead:            aws.Bool(false),
			KeyConditionExpression:    expr.KeyCondition(),
			ExclusiveStartKey:         exclusiveStartKey,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			Limit:                     aws.Int32(int32(batchSize)),
			ScanIndexForward:          aws.Bool(true),
		})
		if err != nil {
			return err
		}

		if len(queryRes.Items) == 0 {
			break
		}

		for _, attrs := range queryRes.Items {

			ev := unmarshalEvent(attrs)

			ev.PartitionKey = partKey

			startSortKey = ev.SortingKey + 1

			evCh <- ev

		}

		exclusiveStartKey = queryRes.LastEvaluatedKey

	}

	return nil
}

func (d *awsDynamodbStore) Tail(ctx context.Context, req *flipbookv1.Event_GetLatestRequest) (*flipbookv1.Event, error) {

	queryReq, err := d.newQueryInput(req.PartitionKey)
	if err != nil {
		return nil, err
	}

	queryReq.Limit = aws.Int32(1)
	queryReq.ScanIndexForward = aws.Bool(false)

	queryRes, err := d.dynamo.Query(ctx, queryReq)
	if err != nil {
		return nil, err
	}

	if len(queryRes.Items) == 0 {
		return nil, nil
	}

	attrs := queryRes.Items[0]

	ev := unmarshalEvent(attrs)

	ev.PartitionKey = req.PartitionKey

	return ev, nil

}

func (d *awsDynamodbStore) newQueryInput(partKey string) (*dynamodb.QueryInput, error) {

	keyEx := expression.Key(partitionKey_field).Equal(expression.Value(partKey))
	projEx := expression.NamesList(
		expression.Name(sortingKey_field),
		expression.Name(eventPayload_field),
	)

	expr, err := expression.NewBuilder().
		WithProjection(projEx).
		WithKeyCondition(keyEx).
		Build()
	if err != nil {
		return nil, err
	}

	return &dynamodb.QueryInput{
		TableName:                 d.tableName,
		AttributesToGet:           []string{sortingKey_field, eventPayload_field},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}, nil

}
