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
	userID, err := util.ExtractUserID(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	templates, err := util.GetAllTemplates(userID)
	if err != nil {
		fmt.Println("Error:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	responseBody, err := json.Marshal(templates)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
