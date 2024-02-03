package main

import (
	"api/shared/constants"
	"api/shared/models"
	"api/shared/util"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type GetUniqueCsvColumnsResponse struct {
	ColumnsMap models.CsvDataColumnUniqueValuesMap
}

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
	reportID := request.QueryStringParameters["reportID"]

	// Check if ReportID is provided
	if reportID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Missing reportID from query string.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	csvColumnsS3Key, err := util.GetReportCsvColumnsS3Key(reportID, userID)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting csvColumnsS3Key by ReportID: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	if csvColumnsS3Key == "no-csv-s3-key" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "csv id not set for report",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Marshal the report into JSON
	columnValues, err := util.GetColumnValuesMapFromS3(csvColumnsS3Key)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting column values map from s3: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Marshal the report into JSON
	columnValuesJSON, err := json.Marshal(columnValues)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling column values into JSON: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Return the report in the response body
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(columnValuesJSON),
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
