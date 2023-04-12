package db

import (
	"app/wenda/utils"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func query_order_by_date(discord_id string, task_date string) *dynamodb.QueryInput {
	log.Println(discord_id, task_date)
	table_name := "task_order"
	expr, _ := expression.NewBuilder().
		WithKeyCondition(
			expression.KeyAnd(
				expression.Key("taskDate").Equal(expression.Value(task_date)), expression.Key("discordID").Equal(expression.Value(discord_id)),
			)).
		WithProjection(order_proj).
		Build()

	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table_name),
	}
	return queryInput
}

func get_order(discord_id string, date string) (TaskOrder, error) {
	params := query_order_by_date(discord_id, date)

	result, err := svc.Query(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return TaskOrder{}, nil
	}

	if len(result.Items) == 0 {
		return TaskOrder{
			DiscordID: discord_id,
			TaskDate:  date,
			Order:     make([]string, 0),
		}, nil
	}

	var ord TaskOrder
	if dynamodbattribute.UnmarshalMap(result.Items[0], &ord); err != nil {
		log.Println("Failed to unmarshal task order")
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

	if len(ord.Order) == 0 {
		input := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"discordID": {
					S: aws.String(ord.DiscordID),
				},
				"taskDate": {
					S: aws.String(ord.TaskDate),
				},
			},
			TableName: aws.String(table_name),
		}
		if _, err = svc.DeleteItem(input); err != nil {
			log.Println("failed to delete >_<")
			return err
		}
		return nil
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

func UpdateTaskOrder(uid string, task_id string, init_date time.Time, new_date time.Time, next_task_id string, prev_task_id string) ([]string, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []string{}, errors.New("failed to get discord ID")
	}

	init_ord, err := get_order(discord_id, init_date.Format(no_time_layout))
	if err != nil {
		return []string{}, err
	}

	init_ord.Order = utils.Remove(init_ord.Order, task_id)
	var new_ord TaskOrder
	if new_date == init_date {
		new_ord = init_ord
	} else {
		// Moving across days
		err := UpdateTaskDate(uid, task_id, new_date)
		if err != nil {
			return []string{}, err
		}

		new_ord, err = get_order(discord_id, new_date.Format(no_time_layout))
		if err != nil {
			return []string{}, err
		}
	}

	if len(new_ord.Order) == 0 {
		new_ord.Order = append(new_ord.Order, task_id)
	} else {
		new_ord.Order = utils.InsertBetween(new_ord.Order, task_id, prev_task_id, next_task_id)
	}

	if new_date != init_date {
		if err := update_order(init_ord); err != nil {
			return []string{}, err
		}
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

func RemoveTaskOrder(uid string, task_id string, date string) ([]string, error) {
	discord_id, err := GetDiscordID(uid)
	if err != nil {
		log.Printf("Failed to get discord id for %s", uid)
		return []string{}, errors.New("failed to get discord ID")
	}

	ord, err := get_order(discord_id, date)
	if err != nil {
		return []string{}, err
	}

	ord.Order = utils.Remove(ord.Order, task_id)

	if err := update_order(ord); err != nil {
		return []string{}, err
	}

	return ord.Order, nil
}
