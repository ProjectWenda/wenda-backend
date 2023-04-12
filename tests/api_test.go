package tests

import (
	"app/wenda/db"
	"app/wenda/handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

const task_table = "temp_tasks"

var UID string

type PostBody struct {
	Content  string `json:"content"`
	Status   int    `json:"status"`
	TaskDate string `json:"taskDate"`
}

func (body PostBody) str() string {
	bodystr, err := json.Marshal(body)
	if err != nil {
		log.Println("Issue marshaling")
	}
	return string(bodystr)
}

func load_env() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}
	log.Println(os.Getenv("TEST_UID"))
	UID = os.Getenv("TEST_UID")
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
		Assert(jsonpath.Equal(`$.discordID`, "150708634370703360")).
		Assert(jsonpath.Matches(`$.status`, "1")).
		Assert(jsonpath.Equal(`$.content`, "test task")).
		Assert(jsonpath.Equal(`$.taskDate`, "2023-06-20T03:00:00Z")).
		Status(http.StatusCreated).
		End()
}

func create_task(t *testing.T, content string, status int, task_date string) PostBody {
	body := PostBody{
		Content:  content,
		Status:   status,
		TaskDate: task_date,
	}

	apitest.New().
		Handler(handler.Router()).
		Post("/task").
		Query("uid", UID).
		JSON(body.str()).
		Expect(t).
		Status(http.StatusCreated).
		End()
	return body
}

func TestGetTasks(t *testing.T) {
	init_test()
	// Post some data
	var bodies []PostBody
	bodies = append(bodies, create_task(t, "test task", 1, "2023-06-20T00:00:01Z"))
	bodies = append(bodies, create_task(t, "random", 0, "2023-03-20T00:00:00Z"))
	bodies = append(bodies, create_task(t, "aaa", 2, "2025-06-30T01:00:00Z"))

	for i, body := range bodies {
		apitest.New().
			Handler(handler.Router()).
			Get("/tasks").
			Query("uid", UID).
			Expect(t).
			Assert(jsonpath.Matches(fmt.Sprintf(`$[%d].content`, i), body.Content)).
			Assert(jsonpath.Matches(fmt.Sprintf(`$[%d].status`, i), fmt.Sprint(body.Status))).
			Assert(jsonpath.Matches(fmt.Sprintf(`$[%d].taskDate`, i), body.TaskDate)).
			Status(http.StatusOK).
			End()
	}
}

func TestGetTask(t *testing.T) {
	init_test()
	// Post some data
	var bodies []PostBody
	bodies = append(bodies, create_task(t, "test task", 1, "2023-06-20T00:00:01Z"))
	bodies = append(bodies, create_task(t, "random", 0, "2023-03-20T00:00:00Z"))
	bodies = append(bodies, create_task(t, "aaa", 2, "2025-06-30T01:00:00Z"))

	for i, body := range bodies {
		apitest.New().
			Handler(handler.Router()).
			Get("/task").
			Query("uid", UID).
			Query("taskID", fmt.Sprint(i+1)).
			Expect(t).
			Assert(jsonpath.Matches(`$.content`, body.Content)).
			Assert(jsonpath.Matches(`$.status`, fmt.Sprint(body.Status))).
			Assert(jsonpath.Matches(`$.taskDate`, body.TaskDate)).
			Status(http.StatusOK).
			End()
	}
}

func TestPutTask(t *testing.T) {
	init_test()
	create_task(t, "test task", 1, "2023-06-20T00:00:01Z")
	edited_body := PostBody{
		Content:  "new task",
		Status:   0,
		TaskDate: "2023-06-21T00:00:01Z",
	}

	apitest.New().
		Handler(handler.Router()).
		Put("/task").
		Query("uid", UID).
		Query("taskID", fmt.Sprint(1)).
		Query("content", edited_body.Content).
		Query("status", fmt.Sprint(edited_body.Status)).
		Query("taskDate", edited_body.TaskDate).
		Expect(t).
		Assert(jsonpath.Matches(`$.content`, edited_body.Content)).
		Assert(jsonpath.Matches(`$.status`, fmt.Sprint(edited_body.Status))).
		Assert(jsonpath.Matches(`$.taskDate`, edited_body.TaskDate)).
		Status(http.StatusOK).
		End()
}

func TestDeleteTask(t *testing.T) {
	init_test()
	body := create_task(t, "test task", 1, "2023-06-20T00:00:01Z")

	apitest.New().
		Handler(handler.Router()).
		Delete("/task").
		Query("uid", UID).
		Query("taskID", fmt.Sprint(1)).
		Expect(t).
		Assert(jsonpath.Matches(`$.content`, body.Content)).
		Assert(jsonpath.Matches(`$.status`, fmt.Sprint(body.Status))).
		Assert(jsonpath.Matches(`$.taskDate`, body.TaskDate)).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(handler.Router()).
		Get("/task").
		Query("uid", UID).
		Query("taskID", fmt.Sprint(1)).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestUser(t *testing.T) {
	init_test()
	apitest.New().
		Handler(handler.Router()).
		Get("/user").
		Query("uid", UID).
		Expect(t).
		Assert(jsonpath.Matches(`$.username`, "Impact")).
		Assert(jsonpath.Matches(`$.id`, "150708634370703360")).
		Status(http.StatusOK).
		End()
}
