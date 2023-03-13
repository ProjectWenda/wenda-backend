package api

import (
	"fmt"
	"net/http"
	"strconv"

	"app/wenda/db"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	uid := c.Query("uid")
	// SELECT * FROM tasks WITH tasks.uid == uid
	users_tasks := db.SelectUserTasks(uid)
	c.IndentedJSON(http.StatusOK, users_tasks)
}

func GetTaskByID(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("task_id")
	// SELECT * FROM tasks WITH tasks.uid == uid AND task.id == id
	user_task := db.SelectUserTaskByID(uid, taskid)
	if (user_task == db.Task{}) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
	}
	c.IndentedJSON(http.StatusOK, user_task)
}

func PostTask(c *gin.Context) {
	uid := c.Query("uid")
	var newTask db.Task
	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "JSON formatted incorrectly in postTasks"})
	}
	newTask.DiscordID = db.SelectDiscordID(uid)
	newTask.ID = db.InsertTask(newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func UpdateTask(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("task_id")
	content := c.Query("content")
	status, err := strconv.Atoi(c.Query("status"))
	fmt.Println("STATUS", status, err)
	// verify status is valid
	if err != nil || (status != 0 && status != 1 && status != 2) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "status should be 0, 1, or 2"})
		return
	}
	if !db.UpdateTask(uid, taskid, content, status) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, db.SelectUserTaskByID(uid, taskid))
}

func DeleteTask(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("task_id")
	task := db.SelectUserTaskByID(uid, taskid)
	if !db.DeleteTask(uid, taskid) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + taskid + " not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}
