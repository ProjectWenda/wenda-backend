package main

import (
	"fmt"

	"app/wenda/db"
	"app/wenda/handler"

	"github.com/joho/godotenv"
)

func load_env() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	// Load ENV
	load_env()
	db.InitDB()
	router := handler.Router()
	// Run
	router.Run("localhost:8080")
}
