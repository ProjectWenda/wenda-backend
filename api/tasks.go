package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app/wenda/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var tasks = []db.Task{
	{ID: "1", UID: "$2a$10$xesny5lQml6vHUltQ7Diw.iOARAQrr3nw5GBqehg6BzWAKbm8r2AC", TimeCreated: time.Now(), LastModified: time.Now(), Content: "Test Message 1", Status: db.ToDo},
	{ID: "1", UID: "gaming", TimeCreated: time.Now(), LastModified: time.Now(), Content: "Hello Hello", Status: db.ToDo},
}

func GetTasks(c *gin.Context) {
	uid := c.Query("uid")

	users_tasks := []db.Task{}
	// SELECT * FROM tasks WITH tasks.uid == uid
	for _, task := range tasks {
		if task.UID == uid {
			users_tasks = append(users_tasks, task)
		}
	}

	c.IndentedJSON(http.StatusOK, users_tasks)
}

func GetTaskByID(c *gin.Context) {
	taskid := c.Query("task_id")
	uid := c.Query("uid")

	// SELECT * FROM tasks WITH tasks.uid == uid AND task.id == id
	for _, task := range tasks {
		if task.UID == uid && task.ID == taskid {
			c.IndentedJSON(http.StatusOK, task)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
}

func PostTask(c *gin.Context) {
	uid := c.Query("uid")
	var newTask db.Task

	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println("JSON formatted incorrectly in postTasks")
	}

	// Place UID and Task ID
	newTask.UID = uid
	newTask.ID = uuid.New().String()
	newTask.TimeCreated = time.Now()
	newTask.LastModified = time.Now()

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func UpdateTask(c *gin.Context) {
	uid := c.Query("uid")
	taskid := c.Query("task_id")
	content := c.Query("content")
	status, err := strconv.Atoi(c.Query("status"))

	// verify status is valid
	if err != nil || (status != 0 && status != 1 && status != 2) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "status should be 0, 1, or 2"})
		return
	}

	for i, task := range tasks {
		if task.UID == uid && task.ID == taskid {
			tasks[i].Content = content
			tasks[i].Status = db.TaskStatus(status)
			tasks[i].LastModified = time.Now()
			c.IndentedJSON(http.StatusOK, task)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
}

func DeleteTask(c *gin.Context) {
	uid := c.Query("uid")
	taskid := c.Query("task_id")

	for i, task := range tasks {
		if task.UID == uid && task.ID == taskid {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.IndentedJSON(http.StatusOK, task)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
}
