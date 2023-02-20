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
	"go.uber.org/zap"
)

type awsDynamodbStore struct {
	dynamo           *dynamodb.Client
	tableName        *string
	defaultBatchSize int64
	log              *zap.SugaredLogger
}

func (d *awsDynamodbStore) Write(ctx context.Context, req *flipbookv1.Event_AppendRequest, sem helpers.SortedEventMap) error {

	putItens := []types.TransactWriteItem{}

	for partKey, evs := range sem {

		log := d.log.With("partition_key", partKey)

		for _, ev := range evs {

			putItens = append(putItens, types.TransactWriteItem{
				Put: &types.Put{
					TableName: d.tableName,
					Item:      marshalEvent(ev),
				},
			})

		}

		firstEv := evs[0]
		lastEv := evs[len(evs)-1]

		log.Debugw("Inserting events", "sorting_key_type", firstEv.SortingKeyType.String(), "first_event_sorting_key", firstEv.SortingKey)

		exprBuilder := expression.NewBuilder()

		if firstEv.SortingKeyType == flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE {

			if firstEv.SortingKey == 0 {
				exprBuilder = exprBuilder.WithCondition(expression.AttributeNotExists(
					expression.Name(partitionKey_field),
				))
			} else {
				exprBuilder = exprBuilder.WithCondition(expression.Equal(
					expression.Name(lastSortingKey_field),
					expression.Value(marshalSortingKey(firstEv.SortingKey-1)),
				))
			}

		}

		expr, err := exprBuilder.Build()
		if err != nil {
			return err
		}

		putItens = append(putItens, types.TransactWriteItem{
			Put: &types.Put{
				TableName: d.tableName,
				Item: map[string]types.AttributeValue{
					partitionKey_field:   marshalPartitionKey(partKey),
					sortingKey_field:     marshalSortingKey(metaSortingKey_value),
					lastSortingKey_field: marshalSortingKey(lastEv.SortingKey),
				},
				ConditionExpression:       expr.Condition(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
			},
		})

	}

	_, err := d.dynamo.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: putItens,
	})

	if err != nil {

		txErr := &types.TransactionCanceledException{}

		if errors.As(err, &txErr) {

			reason := txErr.CancellationReasons[len(txErr.CancellationReasons)-1]

			if aws.ToString(reason.Code) == "ConditionalCheckFailed" {
				return exceptions.ThrowOptimisticLockFailedErr()
			}

		}

		return err

	}

	return nil

}

func (d *awsDynamodbStore) Read(ctx context.Context, req *flipbookv1.Event_IterateRequest, evCh chan *flipbookv1.Event) error {

	d.log.Debug("Starting read loop")

	batchSize := req.BatchSize

	if batchSize <= 0 {
		batchSize = d.defaultBatchSize
	}

	partKey := req.PartitionKey
	startSortKey := req.Query.StartSortingKey - 1

	projEx := expression.NamesList(
		expression.Name(sortingKey_field),
		expression.Name(eventPayload_field),
	)

	var exclusiveStartKey map[string]types.AttributeValue = nil

	for {

		keyEx := expression.Key(partitionKey_field).Equal(expression.Value(partKey))

		d.log.Debug("Starting pagination loop")

		switch req.Query.Stop {

		case flipbookv1.QueryStop_QUERY_STOP_EXACT:
			keyEx = keyEx.And(expression.KeyLessThanEqual(
				expression.Key(sortingKey_field),
				expression.Value(req.Query.StopSortingKey),
			))

		case flipbookv1.QueryStop_QUERY_STOP_LATEST:
			if exclusiveStartKey == nil {
				exclusiveStartKey = map[string]types.AttributeValue{
					partitionKey_field: marshalPartitionKey(partKey),
					sortingKey_field:   marshalSortingKey(startSortKey),
				}
			}

		}

		expr, err := expression.NewBuilder().
			WithKeyCondition(keyEx).
			WithProjection(projEx).
			Build()
		if err != nil {
			return err
		}

		d.log.Debug("Send DynamoDB query")

		q := &dynamodb.QueryInput{
			TableName:                 d.tableName,
			ProjectionExpression:      expr.Projection(),
			ConsistentRead:            aws.Bool(false),
			KeyConditionExpression:    expr.KeyCondition(),
			ExclusiveStartKey:         exclusiveStartKey,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			Limit:                     aws.Int32(int32(batchSize)),
			ScanIndexForward:          aws.Bool(true),
		}

		queryRes, err := d.dynamo.Query(ctx, q)
		if err != nil {
			return err
		}

		d.log.Debug("DynamoDB query executed")

		if len(queryRes.Items) == 0 {
			d.log.Debugw(
				"DynamoDB query returned empty items list",
				"key_condition_expression", *q.KeyConditionExpression,
				"exclusive_start_key", exclusiveStartKey,
				"limit", *q.Limit,
			)
			break
		}

		for _, attrs := range queryRes.Items {

			d.log.Debug("New item decoded")

			ev := unmarshalEvent(attrs)

			ev.PartitionKey = partKey

			select {
			case <-ctx.Done():
				return ctx.Err()
			case evCh <- ev:
				continue
			}

		}

		d.log.Debugw("Last evaluated key", "last_evaluated_key", queryRes.LastEvaluatedKey)

		if queryRes.LastEvaluatedKey == nil {
			break
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
