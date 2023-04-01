package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Cache discord ID
var uid_to_discordID map[string]string

func init() {
	uid_to_discordID = make(map[string]string)
}

func filter_users_by_uid(uid string) *dynamodb.ScanInput {
	table_name := "users"
	filt := expression.Name("uid").Equal(expression.Value(uid))
	return form_params(filt, user_proj, table_name)
}

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
	if id, exists := uid_to_discordID[uid]; exists {
		return id, nil
	}
	user, err := GetUser(uid)
	if err != nil {
		log.Printf("Failed to get user %s", err)
		return "", err
	}
	uid_to_discordID[uid] = user.DiscordID
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
