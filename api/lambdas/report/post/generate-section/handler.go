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

type AddSectionRequest struct {
	ReportID             string          `json:"reportID"`
	PartIndex            int             `json:"partIndex"`
	SectionIndex         int             `json:"sectionIndex"`
	Answers              []models.Answer `json:"answers"`
	RegenGeneratedOutput bool            `json:"regenGeneratedOutput"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req AddSectionRequest
	err := json.Unmarshal([]byte(request.Body), &req)
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

	err = util.GenerateSection(req.ReportID, req.PartIndex, req.SectionIndex, req.Answers, req.RegenGeneratedOutput)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Error generating section: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    constants.CorsHeaders,
		Body:       "Section generated successfully",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
