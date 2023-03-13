package main

import (
	"fmt"

	"app/wenda/api"
	"app/wenda/db"

	"github.com/gin-contrib/cors"
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

func bind_discord(router *gin.Engine) {
	router.GET("/user", api.GetDiscordUser)
	router.GET("/friends", api.GetDiscordFriends)
}

func set_cors(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:8080"}
	router.Use(cors.New(config))
}

func main() {
	// Load ENV
	load_env()
	db.DB()
	// fmt.Println(db.SelectAll())
	// db.SelectDiscordID("$2a$10$7XRZE0DdXqlhhPXiBNJBae2unc36BOgZex3nG5aEYSVo8OmP7yo4i")
	// fmt.Println(db.SelectUserTasks("150708634370703360"))
	// Create router
	router := gin.Default()
	// CORS config
	set_cors(router)
	// Endpoints
	bind_task_crud(router)
	bind_discord(router)
	router.GET("/auth", api.GetAuth)
	// Run
	router.Run("localhost:8080")
}
