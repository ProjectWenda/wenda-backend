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
