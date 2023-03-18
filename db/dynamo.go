package db

import (
	"errors"
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

// FILTERS
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

func filter_tasks_by_discordID(discord_id string) *dynamodb.ScanInput {
	table_name := "tasks"
	filt := expression.Name("discordID").Equal(expression.Value(discord_id))

	proj := expression.NamesList(
		expression.Name("taskID"), expression.Name("content"), expression.Name("discordID"), expression.Name("lastModified"), expression.Name("status"), expression.Name("taskDate"), expression.Name("timeCreated"),
	)

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

// QUERIES
func GetDiscordID(uid string) (string, error) {
	params := filter_users_by_uid(uid)

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return "", errors.New("query failed")
	}

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		log.Fatalf("Failed to unmarshal user data")
		return "", errors.New("failed to unmarshal")
	}

	fmt.Println("user", user)
	return user.DiscordID, nil
}

func GetUserTasks(uid string) ([]Task, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Fatalf("Failed to get discord id for %s", uid)
		return []Task{}, errors.New("failed to get discord ID")
	}

	params := filter_tasks_by_discordID(discord_id)

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return []Task{}, errors.New("query failed")
	}

	var tasks []Task

	for _, i := range result.Items {
		task := Task{}
		if err := dynamodbattribute.UnmarshalMap(i, &task); err != nil {
			log.Fatalf("Failed to unmarshal user data")
			return []Task{}, errors.New("failed to unmarshal data")
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
