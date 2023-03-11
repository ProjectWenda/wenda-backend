package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskStatus int64

const (
	ToDo TaskStatus = iota
	Completed
	Archived
)

type task struct {
	ID           int        `json:"id"`
	TimeCreated  time.Time  `json:"time_created"`
	LastModified time.Time  `json:"last_modified"`
	Content      string     `json:"content"`
	Status       TaskStatus `json:"status"`
}

var tasks = []task{
	{ID: 1, TimeCreated: time.Now(), LastModified: time.Now(), Content: "Test Message 1", Status: ToDo},
	{ID: 2, TimeCreated: time.Now(), LastModified: time.Now(), Content: "Hello Hello", Status: ToDo},
}

func GetTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}
