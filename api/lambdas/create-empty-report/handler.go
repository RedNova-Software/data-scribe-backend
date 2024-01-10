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

type ReportRequest struct {
	ReportType string `json:"reportType"` // JSON tags should match the JSON keys
	Title      string `json:"title"`
	City       string `json:"city"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
		}, nil
	}

	var req ReportRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ReportType == "" || req.Title == "" || req.City == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: reportType, title and city are required.",
		}, nil
	}

	reportID := uuid.New().String()
	report := models.Report{
		ReportID:   reportID,
		ReportType: req.ReportType,
		Title:      req.Title,
		City:       req.City,
		Parts:      []models.Part{}, // assuming models.Part is a defined struct
	}

	tableName := os.Getenv(string(constants.ReportTable))
	err = util.PutItem(report, tableName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Empty report created successfully with ID: " + reportID,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
