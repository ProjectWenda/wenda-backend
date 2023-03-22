package handler

import (
	"app/wenda/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:8080", "http://127.0.0.1:5173", "https://068a6tnn1g.execute-api.us-east-2.amazonaws.com/staging"}
	config.ExposeHeaders = []string{"set-cookie"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
}

func Router() *gin.Engine {
	// Create router
	router := gin.Default()
	// CORS config
	set_cors(router)
	// Endpoints
	bind_task_crud(router)
	bind_discord(router)
	router.GET("/auth", api.GetAuth)
	return router
}
