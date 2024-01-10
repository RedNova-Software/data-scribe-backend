package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"api/shared/constants"
	"api/shared/models"
	"api/shared/util"
	"net/http"
	"os"
)

type AddPartRequest struct {
	ReportID  string `json:"reportID"`
	PartTitle string `json:"partTitle"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
		}, nil
	}

	var req AddPartRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ReportID == "" || req.PartTitle == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: reportID and partTitle are required.",
		}, nil
	}

	tableName := os.Getenv(string(constants.ReportTable))
	part := map[string]interface{}{
		constants.TitleField.String(): req.PartTitle,
		constants.HeadersField.String(): []models.Header{},
	}

	err = util.AddElementToReportList(tableName, req.ReportID, constants.PartsField.String(), part)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Part added successfully to report with ID: " + req.ReportID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
