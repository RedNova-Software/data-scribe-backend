package main

import (
	"api/shared/constants"
	"api/shared/models"
	"api/shared/util"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
		}, nil
	}

	tableName := os.Getenv(string(constants.ReportTable))
	item, err := util.GetItem(tableName, "ReportID", reportID)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting report by ReportID: " + err.Error(),
		}, nil
	}

	if item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Report not found",
		}, nil
	}

	var report models.Report

	err = dynamodbattribute.UnmarshalMap(item, &report)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error unmarshalling dynamo item into report: " + err.Error(),
		}, nil
	}

	// Marshal the report into JSON
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling report into JSON: " + err.Error(),
		}, nil
	}

	// Return the report in the response body
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(reportJSON),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
