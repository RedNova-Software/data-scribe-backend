package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	users, err := util.GetAllUsers()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting all users: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	responseBody, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Return the users in the response body
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
