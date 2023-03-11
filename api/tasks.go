package api

import (
	"fmt"
	"net/http"
	"strconv"
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
	ID           string     `json:"id"`
	TimeCreated  time.Time  `json:"time_created"`
	LastModified time.Time  `json:"last_modified"`
	Content      string     `json:"content"`
	Status       TaskStatus `json:"status"`
}

var tasks = []task{
	{ID: "1", TimeCreated: time.Now(), LastModified: time.Now(), Content: "Test Message 1", Status: ToDo},
	{ID: "2", TimeCreated: time.Now(), LastModified: time.Now(), Content: "Hello Hello", Status: ToDo},
}

// TODO: tasks specific to each user

func GetTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	taskid := c.Query("task_id")
	//uid := c.Query("uid")

	for _, task := range tasks {
		if task.ID == taskid {
			c.IndentedJSON(http.StatusOK, task)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
}

func PostTask(c *gin.Context) {
	//uid := c.Query("uid")
	var newTask task

	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println("JSON formatted incorrectly in postTasks")
	}

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func UpdateTask(c *gin.Context) {
	//uid := c.Query("uid")
	taskid := c.Query("task_id")
	content := c.Query("new_content")
	status, err := strconv.Atoi(c.Query("status"))

	// verify status is valid
	if err != nil || (status != 0 && status != 1 && status != 2) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "status should be 0, 1, or 2"})
		return
	}

	for i, task := range tasks {
		if task.ID == taskid {
			tasks[i].Content = content
			tasks[i].Status = TaskStatus(status)
			c.IndentedJSON(http.StatusOK, task)
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
}

func DeleteTask(c *gin.Context) {
	//uid := c.Query("uid")
	taskid := c.Query("task_id")

	for i, task := range tasks {
		if task.ID == taskid {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted task " + taskid})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
}
