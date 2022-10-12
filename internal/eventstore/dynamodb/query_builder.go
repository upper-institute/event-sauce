package dynamodb

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	apiv1 "github.com/upper-institute/flipbook/pkg/api/v1"
)

type QueryBuilder struct {
	params                *dynamodb.ScanInput
	naturalTimestampIndex string
	id                    string
}

var (
	filterVersion          = aws.String("#v <= :end_range")
	filterNaturalTimestamp = aws.String("#nt <= :end_range")
)

func (q *QueryBuilder) SetVersionRange(start int64, end int64, operator apiv1.QueryOperator) {

	if operator == apiv1.QueryOperator_QUERY_OPERATOR_LESS_THAN_OR_EQUAL {

		q.params.ExpressionAttributeNames = map[string]string{
			"#v": "version",
		}

		q.params.ExpressionAttributeValues = map[string]types.AttributeValue{
			":end_range": &types.AttributeValueMemberN{Value: strconv.FormatInt(end, 10)},
		}

		q.params.FilterExpression = filterVersion

	}

	q.params.ExclusiveStartKey["version"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(start, 10)}

}

func (q *QueryBuilder) SetNaturalTimestampRange(start time.Time, end time.Time, operator apiv1.QueryOperator) {

	q.params.IndexName = aws.String(q.naturalTimestampIndex)

	if operator == apiv1.QueryOperator_QUERY_OPERATOR_LESS_THAN_OR_EQUAL {

		q.params.ExpressionAttributeNames = map[string]string{
			"#nt": "naturalTimestamp",
		}

		q.params.ExpressionAttributeValues = map[string]types.AttributeValue{
			":end_range": &types.AttributeValueMemberS{Value: end.Format(time.RFC3339Nano)},
		}

		q.params.FilterExpression = filterNaturalTimestamp
	}

	q.params.ExclusiveStartKey["version"] = zeroVersionAttr
	q.params.ExclusiveStartKey["naturalTimestamp"] = &types.AttributeValueMemberS{Value: start.Format(time.RFC3339Nano)}

}
