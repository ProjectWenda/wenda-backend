package tests

import (
	"app/wenda/db"
	"app/wenda/handler"
	"fmt"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

const UID = "$2a$10$2z9RNlHH.3bN6LK9ITuHBu7pLXQKwVMJ5KTanXy6i1UsJLO2nGNA2"
const task_table = "temp_tasks"

func load_env() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func create_table() bool {
	query := fmt.Sprintf(`CREATE TEMP TABLE
		%s (
			id bigserial NOT NULL,
			discord_id text NOT NULL,
			time_created timestamp without time zone NOT NULL DEFAULT now(),
			last_modified timestamp without time zone NOT NULL DEFAULT now(),
			content text NOT NULL,
			status bigint NOT NULL,
			task_date timestamp without time zone NOT NULL
		);
		ALTER TABLE
			%s
		ADD
			CONSTRAINT tasks_pkey PRIMARY KEY (id)
		`, task_table, task_table,
	)

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
	body := `{"content": "test task", "status": 1, "taskDate": "2023-06-20T03:00:00Z"}`
	apitest.New().
		Handler(handler.Router()).
		Post("/task").
		Query("uid", UID).
		JSON(body).
		Expect(t).
		Assert(jsonpath.Matches(`$.id`, "1")).
		Assert(jsonpath.Equal(`$.discordID`, "490574905985728523")).
		Assert(jsonpath.Matches(`$.status`, "1")).
		Assert(jsonpath.Equal(`$.content`, "test task")).
		Assert(jsonpath.Equal(`$.taskDate`, "2023-06-20T03:00:00Z")).
		Status(http.StatusCreated).
		End()
}

func create_task(t *testing.T, body string) {
	apitest.New().
		Handler(handler.Router()).
		Post("/task").
		Query("uid", UID).
		JSON(body).
		Expect(t).
		Status(http.StatusCreated).
		End()
}

func TestGetTasks(t *testing.T) {
	init_test()
	// Post some data
	create_task(t, `{"content": "test task", "status": 1, "taskDate": "2023-06-20T00:00:00Z"}`)
	create_task(t, `{"content": "random", "status": 0, "taskDate": "2023-03-20T00:00:00Z"}`)
	create_task(t, `{"content": "aaa", "status": 2, "taskDate": "2025-06-30T01:00:00Z"}`)

	apitest.New().
		Handler(handler.Router()).
		Get("/tasks").
		Query("uid", UID).
		Expect(t).
		Assert(jsonpath.Matches(`$[0].content`, "test task")).
		Assert(jsonpath.Matches(`$[0].status`, "1")).
		Assert(jsonpath.Matches(`$[0].taskDate`, "2023-06-20T00:00:00Z")).
		Assert(jsonpath.Matches(`$[1].content`, "random")).
		Assert(jsonpath.Matches(`$[1].status`, "0")).
		Assert(jsonpath.Matches(`$[1].taskDate`, "2023-03-20T00:00:00Z")).
		Assert(jsonpath.Matches(`$[2].content`, "aaa")).
		Assert(jsonpath.Matches(`$[2].status`, "2")).
		Assert(jsonpath.Matches(`$[2].taskDate`, "2025-06-30T01:00:00Z")).
		Status(http.StatusOK).
		End()
}
