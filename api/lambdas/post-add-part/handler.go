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
)

type AddPartRequest struct {
	ReportID  string `json:"reportID"`
	Index     uint16 `json:"partIndex"`
	PartTitle string `json:"partTitle"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req AddPartRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ReportID == "" || req.PartTitle == "" || req.Index < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: reportID, partTitle, and partIndex are required.",
		}, nil
	}

	tableName := os.Getenv(string(constants.ReportTable))

	updatedIndices, err := util.ModifyReportPartIndices(tableName, req.ReportID, req.Index, true) // Increment all index values equal and above this part

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	newPart := models.Part{
		Title:    req.PartTitle,
		Index:    req.Index,
		Sections: []models.Section{},
	}

	err = util.AddPartToReport(tableName, req.ReportID, newPart)
	if err != nil {
		if updatedIndices {
			util.ModifyReportPartIndices(tableName, req.ReportID, req.Index, false) // Return indices of parts back to normal. Maybe handle this response soon
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    constants.CorsHeaders,
		Body:       "Part added successfully to report with ID: " + req.ReportID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
