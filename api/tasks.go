package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app/wenda/db"

	"github.com/gin-gonic/gin"
)

const time_layout = "2006-01-02"

func GetTasks(c *gin.Context) {
	uid := c.Query("uid")
	// SELECT * FROM tasks WITH tasks.uid == uid
	users_tasks := db.SelectUserTasks(uid)
	if len(users_tasks) == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{})
		return
	}
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
	// res, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println(string(res))
	type PostBody struct {
		Content  string `json:"content"`
		Status   int8   `json:"status"`
		TaskDate string `json:"taskDate"`
	}

	var post_body PostBody
	if err := c.BindJSON(&post_body); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "JSON formatted incorrectly"})
		return
	}

	task_date, err := time.Parse(time_layout, post_body.TaskDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "incorrectly formatted time"})
		return
	}

	new_task := db.Task{
		DiscordID: db.SelectDiscordID(uid),
		Content:   post_body.Content,
		Status:    post_body.Status,
		TaskDate:  task_date,
	}
	new_task.ID = db.InsertTask(new_task)
	c.IndentedJSON(http.StatusCreated, new_task)
}

func UpdateTask(c *gin.Context) {
	uid, taskid := c.Query("uid"), c.Query("task_id")
	content := c.Query("content")
	time_str := c.Query("task_date")
	task_date, err := time.Parse(time_layout, time_str)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "incorrectly formatted time"})
		return
	}
	status, err := strconv.Atoi(c.Query("status"))
	// verify status is valid
	if err != nil || (status != 0 && status != 1 && status != 2) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "status should be 0, 1, or 2"})
		return
	}

	if !db.UpdateTask(uid, taskid, content, status, task_date) {
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
