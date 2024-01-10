package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"api/shared/constants"
	"api/shared/models"
	"api/shared/util"
	"github.com/google/uuid"
	"net/http"
	"os"
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
		}, nil
	}

	var req AddReportRequest
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
		Parts: []models.Part{
			{
				Title: "Part Title",
				Headers: []models.Header{
					{
						Title: "Header Title",
						Sections: []models.Section{
							{
								Title: "Section Title",
								Questions: []models.Question{}, // An empty slice of Questions, you can add questions here as needed
							},
						},
					},
				},
			},
		},
	}
	

	tableName := os.Getenv(string(constants.ReportTable))
	err = util.PutItem(tableName, report)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Empty report created successfully with ID: " + reportID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}