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

type SetSectionResponseRequest struct {
	ReportID             string                       `json:"reportID"`
	PartIndex            int                          `json:"partIndex"`
	SectionIndex         int                          `json:"sectionIndex"`
	Answers              []models.Answer              `json:"answers"`
	CsvDataResponses     []models.CsvDataResponse     `json:"csvDataResponses"`
	ChartOutputResponses []models.ChartOutputResponse `json:"chartOutputResponses"`
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

	var req SetSectionResponseRequest
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ReportID == "" || req.PartIndex < 0 || req.SectionIndex < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: reportID, sectionTitle, and sectionIndex are required.",
		}, nil
	}

	err = util.SetReportSectionResponses(req.ReportID, req.PartIndex, req.SectionIndex, req.Answers, req.CsvDataResponses, req.ChartOutputResponses, userID)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Error setting section responses: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    constants.CorsHeaders,
		Body:       "Section responses set successfully",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
