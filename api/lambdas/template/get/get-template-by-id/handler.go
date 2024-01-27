package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"encoding/json"
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

	// Extract the ReportID from the query string parameters
	// Access the reportID query string parameter
	templateID := request.QueryStringParameters["templateID"]

	// Check if ReportID is provided
	if templateID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Missing templateID from query string.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	template, err := util.GetTemplate(templateID, userID)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting template by TemplateID: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	if template == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Template not found",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Marshal the report into JSON
	templateJSON, err := json.Marshal(template)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling template into JSON: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Return the report in the response body
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(templateJSON),
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
