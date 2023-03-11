package main

import (
	auth "app/wenda/api"
	"fmt"
	"os"

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
	//fmt.Println("hi")
	router := gin.Default()
	router.GET("/auth", auth.GetAuth)
	router.Run("localhost:8080")
}
