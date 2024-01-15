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
	"github.com/google/uuid"
)

type AddReportRequest struct {
	ReportType string `json:"reportType"`
	Title      string `json:"title"`
	City       string `json:"city"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	var req AddReportRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	if req.ReportType == "" || req.Title == "" || req.City == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: reportType, title and city are required.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	reportID := uuid.New().String()

	report := models.Report{
		ReportID:   reportID,
		ReportType: req.ReportType,
		Title:      req.Title,
		City:       req.City,
		Parts:      make([]models.Part, 0),
	}

	tableName := os.Getenv(string(constants.ReportTable))

	err = util.PutNewReport(tableName, report)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Empty report created successfully with ID: " + reportID,
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}