package db

import (
	"fmt"
	"os"

	supa "github.com/nedpals/supabase-go"
)

func create_client() *supa.Client {
	return supa.CreateClient(os.Getenv("SUPABASE_API_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
}

func SelectAll() map[string]interface{} {
	supabase := create_client()

	var results map[string]interface{}
	err := supabase.DB.From("tasks").Select("*").Single().Execute(&results)
	if err != nil {
		panic(err)
	}

	fmt.Println(results) // Selected rows
	return results
}

func InsertUser(user User) []User {
	supabase := create_client()

	fmt.Println(user)
	var results []User
	supabase.DB.From("users").Insert(user).Execute(&results)

	// TODO: seems like this returns an error even if it isn't actually failing
	// see: https://github.com/nedpals/supabase-go/issues/3
	// will need to write a small library ourselves to handle interaction with
	// supabase rather than using this

	fmt.Println(results)
	return results
}
