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

// var ginLambda *ginadapter.GinLambda

// func init() {
// 	// Load ENV
// 	db.InitDB()
// 	router := handler.Router()
// 	// Run
// 	//router.Run("localhost:8080")
// 	ginLambda = ginadapter.New(router)
// }

// func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	// If no name is provided in the HTTP request body, throw an error
// 	return ginLambda.ProxyWithContext(ctx, req)
// }

func main() {
	//lambda.Start(Handler)
	load_env()
	db.InitDB()
	router := handler.Router()
	router.Run()
}
