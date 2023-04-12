package db

import (
	e "app/wenda/errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

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

func add_object(in interface{}, table_name string) error {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		log.Printf("Failed to marshal task %s", err)
		return e.ErrInvalidStructure
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table_name),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %s\n", err)
		return e.ErrDBAddFail
	}

	return nil
}

func GetUserTasks(uid string) ([]Task, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []Task{}, err
	}

	params := filter_tasks_by_discordID(discord_id)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return []Task{}, e.ErrDBQueryFail
	}

	var tasks []Task

	for _, i := range result.Items {
		task := Task{}
		if err := dynamodbattribute.UnmarshalMap(i, &task); err != nil {
			log.Printf("Failed to unmarshal user task %s", i)
			return []Task{}, e.ErrInvalidStructure
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetUserTaskByID(uid string, task_id string) (Task, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return Task{}, err
	}

	params := filter_task_by_id(discord_id, task_id)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return Task{}, e.ErrDBQueryFail
	}

	task := Task{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &task); err != nil {
		log.Println("Failed to unmarshal user data")
		return Task{}, e.ErrInvalidStructure
	}

	return task, nil
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
	log.Println("Successfully added " + task.Content + " to table " + table_name)
	return nil
}

func UpdateTaskDate(uid string, task_id string, task_date time.Time) error {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return err
	}
	table_name := "tasks"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
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
		UpdateExpression:    aws.String("set taskDate = :date, lastModified = :modified"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		return e.ErrTaskUpdateFail
	}

	log.Println("Successfully updated task " + task_id)
	return nil
}

func UpdateTask(uid string, task_id string, content string, status int, task_date time.Time) error {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return err
	}
	table_name := "tasks"

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
		return e.ErrTaskUpdateFail
	}

	log.Println("Successfully updated task " + task_id)
	return nil
}

func DeleteTask(uid string, task_id string) error {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return err
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
		return e.ErrTaskDeleteFail
	}

	log.Println("Deleted " + task_id)
	return nil
}
