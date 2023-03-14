package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var TaskTable string = "tasks"

const time_layout = "2006-01-02T00:00:00Z"

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		os.Getenv("AWS_HOST"), 5432,
		os.Getenv("AWS_DBUSER"), os.Getenv("AWS_DBPW"),
		os.Getenv("AWS_DBNAME"),
	)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Failed to open postgres server")
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("Failed to connect to postgres server")
		panic(err)
	}

}

func SelectAllTasks() []Task {
	query := fmt.Sprintf("SELECT * FROM %s", TaskTable)
	rows, err := DB.Query(query)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.DiscordID, &task.TimeCreated, &task.LastModified, &task.Content, &task.Status, &task.TaskDate); err != nil {
			return tasks
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func SelectUserTasks(uid string) []Task {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf("SELECT * FROM %s WHERE discord_id='%s'", TaskTable, discord_id)
	fmt.Println(query)
	rows, err := DB.Query(query)
	if err != nil {
		fmt.Println("Empty")
		return []Task{}
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.DiscordID, &task.TimeCreated, &task.LastModified, &task.Content, &task.Status, &task.TaskDate); err != nil {
			return tasks
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func SelectUserTaskByID(uid string, task_id string) Task {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf("SELECT * FROM %s WHERE discord_id='%s' AND id='%s'", TaskTable, discord_id, task_id)
	var task Task
	err := DB.QueryRow(query).Scan(&task.ID, &task.DiscordID, &task.TimeCreated, &task.LastModified, &task.Content, &task.Status, &task.TaskDate)
	if err != nil {
		return Task{}
	}
	return task
}

func SelectDiscordID(uid string) string {
	query := fmt.Sprintf("SELECT discord_id FROM users WHERE uid='%s'", uid)
	var discord_id string
	err := DB.QueryRow(query).Scan(&discord_id)
	if err != nil {
		panic(err)
	}

	return discord_id
}

func InsertUser(user User) {
	query := fmt.Sprintf("INSERT INTO users (uid, token, discord_id, discord_name) VALUES ('%s', '%s', '%s', '%s')", user.UID, user.Token, user.DiscordID, user.DiscordName)
	_, err := DB.Exec(query)
	if err != nil {
		panic(err)
	}
}

func InsertTask(task Task) int8 {
	query := fmt.Sprintf(
		"INSERT INTO %s (discord_id, content, status, task_date) VALUES ('%s', '%s', '%d', '%s') RETURNING id",
		TaskTable, task.DiscordID, task.Content, task.Status, task.TaskDate.Format(time_layout),
	)
	var id int8
	err := DB.QueryRow(query).Scan(&id)
	if err != nil {
		panic(err)
	}
	return int8(id)
}

func UpdateTask(uid string, task_id string, content string, status int, task_date time.Time) bool {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf(
		"UPDATE %s SET content='%s',status='%d',task_date='%s' WHERE discord_id='%s' AND id='%s'",
		TaskTable, content, status, task_date, discord_id, task_id,
	)
	_, err := DB.Exec(query)
	return err == nil
}

func DeleteTask(uid string, task_id string) bool {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf("DELETE FROM %s WHERE discord_id='%s' AND id='%s'", TaskTable, discord_id, task_id)
	_, err := DB.Exec(query)
	return err == nil
}
