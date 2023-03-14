package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func DB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		os.Getenv("AWS_HOST"), 5432,
		os.Getenv("AWS_DBUSER"), os.Getenv("AWS_DBPW"),
		os.Getenv("AWS_DBNAME"),
	)
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

}

func SelectAllTasks() []Task {
	rows, err := db.Query("SELECT * FROM tasks")
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
	query := fmt.Sprintf("SELECT * FROM tasks WHERE discord_id='%s'", discord_id)
	rows, err := db.Query(query)
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

func SelectUserTaskByID(uid string, task_id string) Task {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf("SELECT * FROM tasks WHERE discord_id='%s' AND id='%s'", discord_id, task_id)
	var task Task
	err := db.QueryRow(query).Scan(&task.ID, &task.DiscordID, &task.TimeCreated, &task.LastModified, &task.Content, &task.Status, &task.TaskDate)
	if err != nil {
		return Task{}
	}
	return task
}

func SelectDiscordID(uid string) string {
	query := fmt.Sprintf("SELECT discord_id FROM users WHERE uid='%s'", uid)
	var discord_id string
	err := db.QueryRow(query).Scan(&discord_id)
	if err != nil {
		panic(err)
	}

	return discord_id
}

func InsertUser(user User) {
	query := fmt.Sprintf("INSERT INTO users (uid, token, discord_id, discord_name) VALUES ('%s', '%s', '%s', '%s')", user.UID, user.Token, user.DiscordID, user.DiscordName)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func InsertTask(task Task) int8 {
	query := fmt.Sprintf(
		"INSERT INTO tasks (discord_id, content, status, task_date) VALUES ('%s', '%s', '%d', '%s') RETURNING id",
		task.DiscordID, task.Content, task.Status, task.TaskDate,
	)
	var id int8
	err := db.QueryRow(query).Scan(&id)
	if err != nil {
		panic(err)
	}
	return int8(id)
}

func UpdateTask(uid string, task_id string, content string, status int, task_date time.Time) bool {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf(
		"UPDATE tasks SET content='%s',status='%d',task_date='%s' WHERE discord_id='%s' AND id='%s'",
		content, status, task_date, discord_id, task_id,
	)
	_, err := db.Exec(query)
	return err == nil
}

func DeleteTask(uid string, task_id string) bool {
	discord_id := SelectDiscordID(uid)
	query := fmt.Sprintf("DELETE FROM tasks WHERE discord_id='%s' AND id='%s'", discord_id, task_id)
	_, err := db.Exec(query)
	return err == nil
}
