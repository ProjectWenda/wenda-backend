package db

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var svc *dynamodb.DynamoDB

func init() {
	// Initialize DB connection
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)
}

func filter_users_by_uid(uid string) *dynamodb.ScanInput {
	table_name := "users"
	filt := expression.Name("uid").Equal(expression.Value(uid))

	proj := expression.NamesList(expression.Name("uid"), expression.Name("discordID"), expression.Name("discordName"), expression.Name("token"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Failed to build query %s", err)
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

func GetDiscordID(uid string) string {
	params := filter_users_by_uid(uid)

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return ""
	}

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		log.Fatalf("Failed to unmarshal user data")
		return ""
	}

	fmt.Println(user)
	return user.DiscordID
}
