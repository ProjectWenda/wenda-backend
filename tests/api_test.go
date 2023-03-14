package tests

import (
	"app/wenda/db"
	"app/wenda/handler"
	"fmt"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	//"github.com/steinfletcher/apitest-jsonpath"
)

const UID = "$2a$10$2z9RNlHH.3bN6LK9ITuHBu7pLXQKwVMJ5KTanXy6i1UsJLO2nGNA2"
const task_table = "temp_tasks"

func load_env() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func create_table() bool {
	query := fmt.Sprintf("CREATE TEMP TABLE %s AS SELECT * FROM %s LIMIT 0", task_table, "tasks")
	_, err := db.DB.Exec(query)
	return err == nil
}

func init_test() {
	load_env()
	// Initialize DB
	db.InitDB()
	// Create Test Table
	create_table()
	// Use the test table
	db.TaskTable = task_table
}

func TestPostTask(t *testing.T) {
	init_test()
	body := `{"content": "test task", "status": 1, "taskDate": "2023-06-20"}`
	apitest.New().
		Handler(handler.Router()).
		Post("/task").
		Query("uid", UID).
		JSON(body).
		Expect(t).
		Body(`{'content': 'test task', 'status': 1, 'task_date': '2023-06-20'}`).
		Status(http.StatusOK).
		End()
}

func TestGetTasks(t *testing.T) {
	init_test()
	apitest.New().
		Handler(handler.Router()).
		Get("/tasks").
		Query("uid", UID).
		Expect(t).
		Body(`{}`).
		Status(http.StatusOK).
		End()
}
