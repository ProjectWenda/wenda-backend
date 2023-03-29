package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var svc *dynamodb.DynamoDB
var (
	user_proj expression.ProjectionBuilder
	task_proj expression.ProjectionBuilder
)

const time_layout = "2006-01-02T15:04:05Z"

func init() {
	// Initialize DB connection
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)

	// User and task cols
	user_proj = expression.NamesList(expression.Name("uid"), expression.Name("discordID"), expression.Name("discordName"), expression.Name("token"))
	task_proj = expression.NamesList(
		expression.Name("taskID"), expression.Name("content"), expression.Name("discordID"), expression.Name("lastModified"), expression.Name("taskStatus"), expression.Name("taskDate"), expression.Name("timeCreated"),
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

func filter_users_by_uid(uid string) *dynamodb.ScanInput {
	table_name := "users"
	filt := expression.Name("uid").Equal(expression.Value(uid))
	return form_params(filt, user_proj, table_name)
}

func filter_tasks_by_discordID(discord_id string) *dynamodb.ScanInput {
	table_name := "tasks"
	filt := expression.Name("discordID").Equal(expression.Value(discord_id))
	return form_params(filt, task_proj, table_name)
}

func filter_task_by_id(discord_id string, task_id string) *dynamodb.ScanInput {
	table_name := "tasks"
	filt := expression.And(
		expression.Name("discordID").Equal(expression.Value(discord_id)),
		expression.Name("taskID").Equal(expression.Value(task_id)),
	)
	return form_params(filt, task_proj, table_name)
}

func filter_task_by_three_id(discord_id string, t1 string, t2 string, t3 string) *dynamodb.ScanInput {
	table_name := "tasks"
	filt := expression.And(
		expression.Name("discordID").Equal(expression.Value(discord_id)),
		expression.Or(
			expression.Name("taskID").Equal(expression.Value(t1)),
			expression.Name("taskID").Equal(expression.Value(t2)),
			expression.Name("taskID").Equal(expression.Value(t3)),
		),
	)
	return form_params(filt, task_proj, table_name)
}

// QUERIES
func add_object(in interface{}, table_name string) error {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		log.Printf("Failed to marshal task %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table_name),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %s\n", err)
		return err
	}

	return nil
}

// USERS
func GetUser(uid string) (User, error) {
	params := filter_users_by_uid(uid)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return User{}, errors.New("query failed")
	}

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		log.Printf("Failed to unmarshal user data")
		return User{}, errors.New("failed to unmarshal")
	}

	return user, nil
}

func GetDiscordID(uid string) (string, error) {
	user, err := GetUser(uid)
	if err != nil {
		log.Printf("Failed to get user %s", err)
		return "", err
	}
	return user.DiscordID, nil
}

func GetUserToken(uid string) (string, error) {
	user, err := GetUser(uid)
	if err != nil {
		log.Printf("Failed to get user %s", err)
		return "", err
	}
	return user.Token, nil
}

func AddUser(user User) error {
	table_name := "users"
	err := add_object(user, table_name)
	if err != nil {
		log.Printf("failed to add user %s", err)
		return err
	}
	fmt.Println("Successfully added " + user.DiscordName + " to table " + table_name)
	return nil
}

// TASKS
func GetUserTasks(uid string) ([]Task, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []Task{}, errors.New("failed to get discord ID")
	}

	params := filter_tasks_by_discordID(discord_id)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return []Task{}, errors.New("query failed")
	}

	var tasks []Task

	for _, i := range result.Items {
		task := Task{}
		if err := dynamodbattribute.UnmarshalMap(i, &task); err != nil {
			log.Printf("Failed to unmarshal user task %s", i)
			return []Task{}, errors.New("failed to unmarshal data")
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetUserTaskByID(uid string, task_id string) (Task, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return Task{}, errors.New("failed to get discord ID")
	}

	params := filter_task_by_id(discord_id, task_id)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return Task{}, errors.New("query failed")
	}

	task := Task{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &task); err != nil {
		log.Printf("Failed to unmarshal user data")
		return Task{}, errors.New("failed to unmarshal data")
	}

	fmt.Println(task.Content)
	return task, nil
}

func GetThreeTask(uid string, t1 string, t2 string, t3 string) ([]Task, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []Task{}, errors.New("failed to get discord ID")
	}

	params := filter_task_by_three_id(discord_id, t1, t2, t3)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return []Task{}, errors.New("query failed")
	}

	tasks := make([]Task, 3)

	for _, i := range result.Items {
		task := Task{}
		if err := dynamodbattribute.UnmarshalMap(i, &task); err != nil {
			log.Printf("Failed to unmarshal user task %s", i)
			return []Task{}, errors.New("failed to unmarshal data")
		}
		switch task.ID {
		case t1:
			tasks[0] = task
		case t2:
			tasks[1] = task
		case t3:
			tasks[2] = task
		}
	}

	return tasks, nil
}

func AddTask(task Task) error {
	table_name := "tasks"
	formatted_task := DBTask{
		ID:           task.ID,
		DiscordID:    task.DiscordID,
		TimeCreated:  task.TimeCreated.Format(time_layout),
		LastModified: task.LastModified.Format(time_layout),
		Content:      task.Content,
		Status:       task.Status,
		TaskDate:     task.TaskDate.Format(time_layout),
	}
	err := add_object(formatted_task, table_name)
	if err != nil {
		log.Printf("Failed to add task %s", err)
		return err
	}
	fmt.Println("Successfully added " + task.Content + " to table " + table_name)
	return nil
}

func UpdateTask(uid string, task_id string, content string, status int, task_date time.Time) error {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return errors.New("failed to get discord ID")
	}
	table_name := "tasks"

	fmt.Println(status)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":content": {
				S: aws.String(content),
			},
			":status": {
				N: aws.String(fmt.Sprint(status)),
			},
			":date": {
				S: aws.String(task_date.Format(time_layout)),
			},
			":modified": {
				S: aws.String(time.Now().Format(time_layout)),
			},
			":id": {
				S: aws.String(discord_id),
			},
		},
		TableName: aws.String(table_name),
		Key: map[string]*dynamodb.AttributeValue{
			"taskID": {
				S: aws.String(task_id),
			},
		},
		ConditionExpression: aws.String("discordID = :id"),
		ReturnValues:        aws.String("UPDATED_NEW"),
		UpdateExpression:    aws.String("set content = :content, taskStatus = :status, taskDate = :date, lastModified = :modified"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		return err
	}

	fmt.Println("Successfully updated task " + task_id)
	return nil
}

func UpdateSortOrder(uid string, task_id string, sort_order string) error {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return errors.New("failed to get discord ID")
	}
	table_name := "tasks"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":order": {
				S: aws.String(sort_order),
			},
			":id": {
				S: aws.String(discord_id),
			},
		},
		TableName: aws.String(table_name),
		Key: map[string]*dynamodb.AttributeValue{
			"taskID": {
				S: aws.String(task_id),
			},
		},
		ConditionExpression: aws.String("discordID = :id"),
		ReturnValues:        aws.String("UPDATED_NEW"),
		UpdateExpression:    aws.String("set sortOrder = :order"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		return err
	}

	fmt.Println("Successfully updated task " + task_id)
	return nil
}

func DeleteTask(uid string, task_id string) error {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return errors.New("failed to get discord ID")
	}
	table_name := "tasks"

	input := &dynamodb.DeleteItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(discord_id),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"taskID": {
				S: aws.String(task_id),
			},
		},
		ConditionExpression: aws.String("discordID = :id"),
		TableName:           aws.String(table_name),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		log.Printf("Got error calling DeleteItem: %s", err)
		return err
	}

	fmt.Println("Deleted " + task_id)
	return nil
}
