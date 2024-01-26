package main

import (
	"api/shared/constants"
	"api/shared/models"
	"api/shared/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type EditSectionRequest struct {
	ItemType        constants.ItemType `json:"itemType"`
	ItemID          string             `json:"itemID"`
	OldPartIndex    uint16             `json:"oldPartIndex"`
	NewPartIndex    uint16             `json:"newPartIndex"`
	OldSectionIndex uint16             `json:"oldSectionIndex"`
	NewSectionIndex uint16             `json:"newSectionIndex"`
	NewSectionTitle string             `json:"newSectionTitle"`
}

type ReportSectionContents struct {
	Questions   []models.ReportQuestion   `json:"questions"`
	TextOutputs []models.ReportTextOutput `json:"textOutputs"`
}

type TemplateSectionContents struct {
	Questions   []models.TemplateQuestion   `json:"questions"`
	TextOutputs []models.TemplateTextOutput `json:"textOutputs"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req EditSectionRequest

	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.OldPartIndex < 0 || req.NewPartIndex < 0 || req.OldSectionIndex < 0 || req.NewSectionIndex < 0 || req.NewSectionTitle == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, oldPartIndex, newPartIndex, oldSectionIndex, newSectionIndex, and newSectionTitle are required.",
		}, nil
	}

	if req.ItemType == constants.Report {
		var sectionContents ReportSectionContents

		err := json.Unmarshal([]byte(request.Body), &sectionContents)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    constants.CorsHeaders,
				Body:       "Bad Request: " + err.Error(),
			}, nil
		}

		err = util.UpdateSectionInReport(req.ItemID, req.OldPartIndex, req.NewPartIndex, req.OldSectionIndex, req.NewSectionIndex, req.NewSectionTitle, sectionContents.Questions, sectionContents.TextOutputs)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    constants.CorsHeaders,
				Body:       "Internal Server Error: " + err.Error(),
			}, nil
		}

	} else if req.ItemType == constants.Template {
		var sectionContents TemplateSectionContents

		err := json.Unmarshal([]byte(request.Body), &sectionContents)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    constants.CorsHeaders,
				Body:       "Bad Request: " + err.Error(),
			}, nil
		}

		err = util.UpdateSectionInTemplate(req.ItemID, req.OldPartIndex, req.NewPartIndex, req.OldSectionIndex, req.NewSectionIndex, req.NewSectionTitle, sectionContents.Questions, sectionContents.TextOutputs)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    constants.CorsHeaders,
				Body:       "Internal Server Error: " + err.Error(),
			}, nil
		}

	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType must be 'report' or 'template' ",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    constants.CorsHeaders,
		Body:       "Section added successfully to report with ID: " + req.ItemID + "and part with index: " + fmt.Sprint(req.OldPartIndex),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
