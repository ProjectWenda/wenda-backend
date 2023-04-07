package db

import "time"

type TaskStatus int8

const (
	ToDo TaskStatus = iota
	Completed
	Archived
)

type User struct {
	UID          string    `json:"uid"`
	Token        string    `json:"token"`
	DiscordID    string    `json:"discordID"`
	DiscordName  string    `json:"discordName"`
	RefreshToken string    `json:"refreshToken"`
	TimeCreated  time.Time `json:"timeCreated"`
}

type Task struct {
	ID           string    `json:"taskID"`
	DiscordID    string    `json:"discordID"`
	TimeCreated  time.Time `json:"timeCreated"`
	LastModified time.Time `json:"lastModified"`
	Content      string    `json:"content"`
	Status       int       `json:"taskStatus"`
	TaskDate     time.Time `json:"taskDate"`
}

type TaskOrder struct {
	DiscordID string   `json:"discordID"`
	TaskDate  string   `json:"taskDate"`
	Order     []string `json:"taskOrder"`
}

type DBTask struct {
	ID           string `json:"taskID"`
	DiscordID    string `json:"discordID"`
	TimeCreated  string `json:"timeCreated"`
	LastModified string `json:"lastModified"`
	Content      string `json:"content"`
	Status       int    `json:"taskStatus"`
	TaskDate     string `json:"taskDate"`
}

type RelationshipResponse struct {
	ID       string `json:"id"`
	Type     int    `json:"type"`
	Nickname string `json:"nickname"`
	User     struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	} `json:"user"`
}
