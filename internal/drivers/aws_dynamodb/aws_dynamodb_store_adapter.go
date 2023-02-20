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
	"github.com/upper-institute/flipbook/internal/logging"
)

const (
	tableName_flag   = "aws.dynamodb.table.name"
	ensureTable_flag = "aws.dynamodb.table.ensure"
	batchSize_flag   = "aws.dynamodb.batch.size"
	endpointUrl_flag = "aws.dynamodb.endpoint.url"
)

type AwsDynamodbStoreAdapter struct {
}

func (d *AwsDynamodbStoreAdapter) Bind(binder helpers.FlagBinder) {

	binder.BindString(tableName_flag, "events", "DynamoDB table name to use as event store")
	binder.BindString(endpointUrl_flag, "", "DynamoDB endpoint URL (useful to use with dynamodb-local)")
	binder.BindBool(ensureTable_flag, true, "Ensure the event store DynamoDB table exists")
	binder.BindInt64(batchSize_flag, 50, "Default batch size to limit at DynamoDB queries")

}

func (d *AwsDynamodbStoreAdapter) New(getter helpers.FlagGetter) (drivers.StoreDriver, error) {

	log := logging.Logger.Sugar().Named("AwsDynamodbDriver")

	endpointUrl := getter.GetString(endpointUrl_flag)

	// https://dave.dev/blog/2021/07/14-07-2021-awsddb/
	// https://gist.github.com/davedotdev/0a4a2361b0e04f63dcda91645db1de83
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpointUrl,
				SigningRegion: "us-east-1",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfgOpts := []func(*config.LoadOptions) error{}

	if len(endpointUrl) > 0 {

		log.Infow("DynamoDB endpoint url", "endpoint_url", endpointUrl)

		cfgOpts = append(cfgOpts, config.WithEndpointResolver(customResolver))

	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
	if err != nil {
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	tableNameStr := getter.GetString(tableName_flag)

	tableName := aws.String(tableNameStr)

	if getter.GetBool(ensureTable_flag) {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log.Infow("Ensure DynamoDB table", "table_name", tableNameStr)

		_, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: tableName,
		})

		dynamoErr := new(types.ResourceNotFoundException)

		if err != nil && errors.As(err, &dynamoErr) {

			_, err = dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
				TableName:   tableName,
				TableClass:  types.TableClassStandard,
				BillingMode: types.BillingModePayPerRequest,
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
		log:              log,
	}, nil

}

func (d *AwsDynamodbStoreAdapter) Destroy(store drivers.StoreDriver) error {
	return nil
}
