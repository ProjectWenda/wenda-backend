package db

import (
	"app/wenda/utils"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func filter_order_by_date(discord_id string, task_date string) *dynamodb.ScanInput {
	table_name := "task_order"
	filt := expression.And(
		expression.Name("discordID").Equal(expression.Value(discord_id)),
		expression.Name("taskDate").Equal(expression.Value(task_date)),
	)
	return form_params(filt, order_proj, table_name)
}

func get_order(discord_id string, date string) (TaskOrder, error) {
	params := filter_order_by_date(discord_id, date)

	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return TaskOrder{}, nil
	}

	if len(result.Items) == 0 {
		return TaskOrder{}, nil
	}

	var ord TaskOrder
	if dynamodbattribute.UnmarshalMap(result.Items[0], &ord); err != nil {
		fmt.Println("Failed to unmarshal task order")
		return TaskOrder{}, errors.New("failed to unmarshal data")
	}

	return ord, nil
}

func update_order(ord TaskOrder) error {
	table_name := "task_order"
	av, err := dynamodbattribute.MarshalMap(ord)
	if err != nil {
		log.Printf("Failed to masrshal order")
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table_name),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		return err
	}
	return nil
}

func GetTaskOrder(uid string, task_date string) ([]string, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []string{}, errors.New("failed to get discord ID")
	}

	ord, err := get_order(discord_id, task_date)
	if err != nil {
		return []string{}, err
	}

	return ord.Order, nil
}

func UpdateTaskOrder(uid string, task_id string, init_date string, new_date string, next_task_id string, prev_task_id string) ([]string, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []string{}, errors.New("failed to get discord ID")
	}

	init_ord, err := get_order(discord_id, init_date)
	if err != nil {
		return []string{}, err
	}

	init_ord.Order = utils.Remove(init_ord.Order, task_id)
	var new_ord TaskOrder
	if new_date == init_date {
		new_ord = init_ord
	} else {
		new_ord, err = get_order(discord_id, new_date)
		if err != nil {
			return []string{}, err
		}
	}
	new_ord.Order = utils.InsertBetween(new_ord.Order, task_id, prev_task_id, next_task_id)

	if err := update_order(init_ord); err != nil {
		return []string{}, err
	}

	if err := update_order(new_ord); err != nil {
		return []string{}, err
	}

	return new_ord.Order, nil
}

func AppendTaskOrder(uid string, task_id string, date string) ([]string, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []string{}, errors.New("failed to get discord ID")
	}

	ord, err := get_order(discord_id, date)
	if err != nil {
		return []string{}, err
	}

	ord.Order = append(ord.Order, task_id)

	if err := update_order(ord); err != nil {
		return []string{}, err
	}

	return ord.Order, nil
}
