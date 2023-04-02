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
const no_time_layout = "2006-01-02"

type TaskDayOutput struct {
	Tasks []db.Task `json:"tasks"`
}

func GetTasks(c *gin.Context) {
	uid := c.Query("uid")
	users_tasks, err := db.GetUserTasks(uid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to retrieve users tasks"})
		return
	}
	if len(users_tasks) == 0 {
		c.IndentedJSON(http.StatusOK, []string{})
		return
	}

	// Map tasks to date
	task_output := make(map[string]TaskDayOutput)
	task_id_to_data := make(map[string][]db.Task)

	for _, task := range users_tasks {
		task_date := task.TaskDate.Format(no_time_layout)
		if entry, exists := task_output[task_date]; exists {
			entry.Tasks = append(task_output[task_date].Tasks, task)
			task_output[task_date] = entry
		} else {
			task_output[task_date] = TaskDayOutput{
				Tasks: append([]db.Task{}, task),
			}
		}

		if entry, exists := task_id_to_data[task.ID]; exists {
			entry = append(entry, task)
			task_id_to_data[task.ID] = entry
		} else {
			task_id_to_data[task.ID] = append([]db.Task{}, task)
		}
	}

	// Output tasks based on sort
	for task_date := range task_output {
		sort_order, err := db.GetTaskOrder(uid, task_date)
		if err != nil {
			fmt.Println("failed to get task order for " + task_date)
			//c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Failed to get task order"})
			continue
		}
		fmt.Println(sort_order)

		var sorted_tasks []db.Task

		for _, id := range sort_order {
			sorted_tasks = append(sorted_tasks, task_id_to_data[id]...)
		}

		if len(sort_order) == 0 {
			task_output[task_date] = TaskDayOutput{
				Tasks: make([]db.Task, 0),
			}
			continue
		}

		task_output[task_date] = TaskDayOutput{
			Tasks: sorted_tasks,
		}
	}

	fmt.Println(task_output)

	c.IndentedJSON(http.StatusOK, task_output)
}

func GetTaskByID(c *gin.Context) {
	uid, task_id := c.Query("uid"), c.Query("taskID")
	user_task, err := db.GetUserTaskByID(uid, task_id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to retrieve task"})
		return
	}
	if (user_task == db.Task{}) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task with id " + task_id + " not found"})
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to get discord ID for user"})
		return
	}

	new_task.ID = uuid.New().String()
	new_task.DiscordID = discord_id
	new_task.TimeCreated = time.Now()
	new_task.LastModified = time.Now()

	db.AddTask(new_task)
	db.AppendTaskOrder(uid, new_task.ID, new_task.TaskDate.Format(no_time_layout))
	c.IndentedJSON(http.StatusCreated, new_task)
}

func UpdateTask(c *gin.Context) {
	uid, task_id := c.Query("uid"), c.Query("taskID")
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

	if err := db.UpdateTask(uid, task_id, content, status, task_date); err != nil {
		fmt.Println("[PUT] failed to update")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to update " + task_id + " in db (maybe wrong id?)"})
		return
	}
	task, err := db.GetUserTaskByID(uid, task_id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error getting updated task " + task_id})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	uid, task_id := c.Query("uid"), c.Query("taskID")
	task, err := db.GetUserTaskByID(uid, task_id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error getting task to delete " + task_id})
		return
	}
	if err := db.DeleteTask(uid, task_id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to delete " + task_id + " in db (maybe wrong id?)"})
		return
	}
	db.RemoveTaskOrder(uid, task_id, task.TaskDate.Format(no_time_layout))

	c.IndentedJSON(http.StatusOK, task)
}

func ChangeOrder(c *gin.Context) {
	uid, task_id, init_datestr, new_datestr, next_task_id, prev_task_id := c.Query("uid"), c.Query("taskID"), c.Query("initialDate"), c.Query("newDate"), c.Query("nextTaskID"), c.Query("prevTaskID")
	init_date, err := time.Parse(time_layout, init_datestr)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": init_datestr + " is provided in invalid format"})
		return
	}
	new_date, err := time.Parse(time_layout, new_datestr)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": new_datestr + " is provided in invalid format"})
		return
	}

	ord, err := db.UpdateTaskOrder(uid, task_id, init_date.Format(no_time_layout), new_date.Format(no_time_layout), next_task_id, prev_task_id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to update order " + task_id})
		return
	}

	c.IndentedJSON(http.StatusOK, ord)
}
