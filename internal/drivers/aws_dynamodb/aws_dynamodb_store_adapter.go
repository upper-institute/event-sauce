package aws_dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/upper-institute/flipbook/internal/drivers"
	"github.com/upper-institute/flipbook/internal/helpers"
)

const (
	tableName_flag   = "aws.dynamodb.table.name"
	ensureTable_flag = "aws.dynamodb.table.ensure"
	batchSize_flag   = "aws.dynamodb.batch.size"
)

type AwsDynamodbStoreAdapter struct {
}

func (d *AwsDynamodbStoreAdapter) Bind(binder helpers.FlagBinder) {

	binder.BindString(tableName_flag, "events", "DynamoDB table name to use as event store")
	binder.BindBool(ensureTable_flag, true, "Ensure the event store DynamoDB table exists")
	binder.BindInt64(batchSize_flag, 50, "Default batch size to limit at DynamoDB queries")

}

func (d *AwsDynamodbStoreAdapter) New(getter helpers.FlagGetter) (drivers.StoreDriver, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	tableName := aws.String(getter.GetString(tableName_flag))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if getter.GetBool(ensureTable_flag) {

		_, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: tableName,
		})

		dynamoErr := &types.ResourceNotFoundException{}

		if errors.As(err, dynamoErr) {

			_, err = dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
				TableName:  tableName,
				TableClass: types.TableClassStandard,
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String(partitionKey_field), AttributeType: types.ScalarAttributeTypeS},
					{AttributeName: aws.String(sortingKey_field), AttributeType: types.ScalarAttributeTypeN},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String(partitionKey_field), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String(sortingKey_field), KeyType: types.KeyTypeRange},
				},
				Tags: []types.Tag{{Key: aws.String("created-by"), Value: aws.String("flipbook")}},
			})

		}

		if err != nil {
			return nil, err
		}

	}

	return &awsDynamodbStore{
		dynamo:           dynamoClient,
		tableName:        tableName,
		defaultBatchSize: getter.GetInt64(batchSize_flag),
	}, nil

}

func (d *AwsDynamodbStoreAdapter) Destroy(store drivers.StoreDriver) error {
	return nil
}
