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

func main() {

	load_env()

	fmt.Println(os.Getenv("CLIENT_ID"))
	router := gin.Default()
	router.GET("/auth", api.GetAuth)
	router.GET("/tasks", api.GetTasks)
	router.Run("localhost:8080")
}
