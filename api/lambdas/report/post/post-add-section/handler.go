package main

import (
	"api/shared/constants"
	"api/shared/models"
	"api/shared/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type AddSectionRequest struct {
	ReportID     string                    `json:"reportID"`
	PartIndex    uint16                    `json:"partIndex"`
	SectionIndex uint16                    `json:"sectionIndex"`
	SectionTitle string                    `json:"sectionTitle"`
	Questions    []models.ReportQuestion   `json:"questions"`
	TextOutputs  []models.ReportTextOutput `json:"textOutputs"`
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

	if req.ReportID == "" || req.PartIndex < 0 || req.SectionTitle == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: reportID, sectionTitle, and partIndex are required.",
		}, nil
	}

	tableName := os.Getenv(constants.ReportTable)

	updatedIndices, err := util.ModifyPartSectionIndices(tableName, req.ReportID, req.PartIndex, req.SectionIndex, true) // Increment all index values equal and above this section

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	newSection := models.ReportSection{
		Title:       req.SectionTitle,
		Index:       req.SectionIndex,
		Questions:   req.Questions,
		TextOutputs: req.TextOutputs,
	}

	err = util.AddSectionToPart(tableName, req.ReportID, req.PartIndex, newSection)

	if err != nil {
		if updatedIndices {
			util.ModifyPartSectionIndices(tableName, req.ReportID, req.PartIndex, req.SectionIndex, false) // Return indices of parts back to normal. Maybe handle this response soon
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Error adding section to part: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    constants.CorsHeaders,
		Body:       "Section added successfully to report with ID: " + req.ReportID + "and part with index: " + fmt.Sprint(req.PartIndex),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
