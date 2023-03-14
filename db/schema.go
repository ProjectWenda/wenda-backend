package db

import "time"

type TaskStatus int8

const (
	ToDo TaskStatus = iota
	Completed
	Archived
)

type User struct {
	UID         string `json:"uid"`
	Token       string `json:"token"`
	DiscordID   string `json:"discord_id"`
	DiscordName string `json:"discord_name"`
}

type Task struct {
	ID           int8      `json:"id"`
	DiscordID    string    `json:"uid"`
	TimeCreated  time.Time `json:"time_created"`
	LastModified time.Time `json:"last_modified"`
	Content      string    `json:"content"`
	Status       int8      `json:"status"`
	TaskDate     time.Time `json:"task_date"`
}
