package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var svc *dynamodb.DynamoDB
var (
	user_proj  expression.ProjectionBuilder
	task_proj  expression.ProjectionBuilder
	order_proj expression.ProjectionBuilder
)

const time_layout = "2006-01-02T15:04:05Z"
const no_time_layout = "2006-01-02"

func init() {
	// Initialize DB connection
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)

	// User and task cols
	user_proj = expression.NamesList(
		expression.Name("uid"), expression.Name("discordID"), expression.Name("discordName"), expression.Name("token"),
	)
	task_proj = expression.NamesList(
		expression.Name("taskID"), expression.Name("content"), expression.Name("discordID"), expression.Name("lastModified"), expression.Name("taskStatus"), expression.Name("taskDate"), expression.Name("timeCreated"),
	)
	order_proj = expression.NamesList(
		expression.Name("discordID"), expression.Name("taskDate"), expression.Name("taskOrder"),
	)
}

// FILTERS
func form_params(filt expression.ConditionBuilder, proj expression.ProjectionBuilder, table_name string) *dynamodb.ScanInput {
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	// TODO: error catch here somehow
	if err != nil {
		log.Printf("Failed to build query %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table_name),
	}
	return params
}
