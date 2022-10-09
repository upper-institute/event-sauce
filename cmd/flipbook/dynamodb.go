package flipbook

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbdriver "github.com/upper-institute/event-sauce/internal/eventstore/dynamodb"
)

var (
	dynamodbEventsTable string = "events"
)

func dynamodbEventStoreBackend() error {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsAccessKeyID,
				SecretAccessKey: awsSecretAccessKey,
				SessionToken:    awsSessionToken,
			},
		}),
		config.WithRegion(awsRegion),
	)

	if err != nil {
		return err
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	eventStoreService.Backend = &dynamodbdriver.DynamoDBEventStore{
		DynamoDB:  dynamoClient,
		TableName: dynamodbEventsTable,
	}

	return nil

}

func init() {

	startCmd.PersistentFlags().StringVar(&dynamodbEventsTable, "dynamodbEventsTable", dynamodbEventsTable, "DynamoDB events table name")

}
