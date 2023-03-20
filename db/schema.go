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
	DiscordID   string `json:"discordID"`
	DiscordName string `json:"discordName"`
}

type Task struct {
	ID           string    `json:"taskID"`
	DiscordID    string    `json:"discordID"`
	TimeCreated  time.Time `json:"timeCreated"`
	LastModified time.Time `json:"lastModified"`
	Content      string    `json:"content"`
	Status       int8      `json:"taskStatus"`
	TaskDate     time.Time `json:"taskDate"`
}

type DBTask struct {
	ID           string `json:"taskID"`
	DiscordID    string `json:"discordID"`
	TimeCreated  string `json:"timeCreated"`
	LastModified string `json:"lastModified"`
	Content      string `json:"content"`
	Status       int8   `json:"taskStatus"`
	TaskDate     string `json:"taskDate"`
}
