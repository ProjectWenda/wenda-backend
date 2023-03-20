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

const time_layout = "2006-01-02T15:04:05Z"

func GetTasks(c *gin.Context) {
	uid := c.Query("uid")
	// SELECT * FROM tasks WITH tasks.uid == uid
	users_tasks, err := db.GetUserTasks(uid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to retrieve users tasks"})
		return
	}
	if len(users_tasks) == 0 {
		c.IndentedJSON(http.StatusOK, []string{})
		return
	}
	c.IndentedJSON(http.StatusOK, users_tasks)
}

func GetTaskByID(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("taskID")
	// SELECT * FROM tasks WITH tasks.uid == uid AND task.id == id
	user_task, err := db.GetUserTaskByID(uid, taskid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to retrieve task"})
		return
	}
	if (user_task == db.Task{}) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, user_task)
}

func PostTask(c *gin.Context) {
	uid := c.Query("uid")
	var new_task db.Task
	if err := c.BindJSON(&new_task); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "JSON formatted incorrectly"})
		return
	}
	discord_id, err := db.GetDiscordID(uid)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to get discord ID for user"})
		return
	}

	new_task.ID = uuid.New().String()
	new_task.DiscordID = discord_id
	new_task.TimeCreated = time.Now()
	new_task.LastModified = time.Now()

	db.AddTask(new_task)
	c.IndentedJSON(http.StatusCreated, new_task)
}

func UpdateTask(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("taskID")
	content := c.Query("content")
	time_str := c.Query("taskDate")
	task_date, err := time.Parse(time_layout, time_str)
	if err != nil {
		fmt.Println("[PUT] incorrectly formatted time")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "incorrectly formatted time"})
		return
	}
	status, err := strconv.Atoi(c.Query("taskStatus"))
	// verify status is valid
	if err != nil || (status != 0 && status != 1 && status != 2) {
		fmt.Println("[PUT] incorrectly formatted status")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "status should be 0, 1, or 2"})
		return
	}

	if err := db.UpdateTask(uid, taskid, content, status, task_date); err != nil {
		fmt.Println("[PUT] failed to update")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to update " + taskid + " in db (maybe wrong id?)"})
		return
	}
	task, err := db.GetUserTaskByID(uid, taskid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error getting updated task " + taskid})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("taskID")
	task, err := db.GetUserTaskByID(uid, taskid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error getting task to delete " + taskid})
		return
	}
	if err := db.DeleteTask(uid, taskid); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to delete " + taskid + " in db (maybe wrong id?)"})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}
