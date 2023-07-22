package main

import (
	"app/wenda/handler"
	"context"
	"flag"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/joho/godotenv"
)

func load_env() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}
}

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	dev := flag.Bool("dev", false, "")
	flag.Parse()
	if *dev {
		load_env()
		router := handler.Router()
		router.Run()
		return
	}

	// Load ENV
	router := handler.Router()
	// Run
	//router.Run("localhost:8080")
	ginLambda = ginadapter.New(router)
	lambda.Start(Handler)
}
