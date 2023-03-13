package db

import "time"

type TaskStatus int64

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
	ID           string     `json:"id,omitempty"`
	UID          string     `json:"uid,omitempty"`
	TimeCreated  time.Time  `json:"time_created,omitempty"`
	LastModified time.Time  `json:"last_modified,omitempty"`
	Content      string     `json:"content"`
	Status       TaskStatus `json:"status"`
}
