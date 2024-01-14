package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract the ReportID from the query string parameters
	// Access the reportID query string parameter
	reportID := request.QueryStringParameters["reportID"]

	// Check if ReportID is provided
	if reportID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Missing reportID from path.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	tableName := os.Getenv(string(constants.ReportTable))
	report, err := util.GetReport(tableName, "ReportID", reportID)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting report by ReportID: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	if report == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Report not found",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Marshal the report into JSON
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling report into JSON: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Return the report in the response body
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(reportJSON),
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
