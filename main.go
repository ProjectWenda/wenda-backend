package main

import (
	"fmt"
	"os"

	"app/wenda/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func load_env() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func bind_task_crud(router *gin.Engine) {
	router.GET("/tasks", api.GetTasks)
	router.GET("/task", api.GetTaskByID)
	router.POST("/task", api.PostTask)
	router.PUT("/task", api.UpdateTask)
	router.DELETE("/task", api.DeleteTask)
}

func main() {

	load_env()

	fmt.Println(os.Getenv("CLIENT_ID"))
	router := gin.Default()
	bind_task_crud(router)
	router.GET("/auth", api.GetAuth)
	router.Run("localhost:8080")
}
