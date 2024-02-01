package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	deletedOnlyString := request.QueryStringParameters["deletedOnly"]

	if deletedOnlyString == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Missing deletedOnly from query string.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	deletedOnlyBool, err := strconv.ParseBool(deletedOnlyString)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: deletedOnly query param must be 'true' or 'false'.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	reports, err := util.GetAllReports(userID, deletedOnlyBool)
	if err != nil {
		fmt.Println("Error:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	responseBody, err := json.Marshal(reports)
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
