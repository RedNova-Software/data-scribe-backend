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

type UpdatedSectionRequest struct {
	ItemType              constants.ItemType `json:"itemType"`
	ItemID                string             `json:"itemID"`
	OldPartIndex          int                `json:"oldPartIndex"`
	NewPartIndex          int                `json:"newPartIndex"`
	OldSectionIndex       int                `json:"oldSectionIndex"`
	NewSectionIndex       int                `json:"newSectionIndex"`
	NewSectionTitle       string             `json:"newSectionTitle"`
	DeleteGeneratedOutput bool               `json:"deleteGeneratedOutput"`
}

type ReportSectionContents struct {
	Questions   []models.ReportQuestion    `json:"questions"`
	TextOutputs []models.ReportTextOutput  `json:"textOutputs"`
	CSVData     []models.ReportCSVData     `json:"csvData"`
	ChartOuput  []models.ReportChartOutput `json:"chartOutputs"`
}

type TemplateSectionContents struct {
	Questions   []models.TemplateQuestion    `json:"questions"`
	TextOutputs []models.TemplateTextOutput  `json:"textOutputs"`
	CSVData     []models.TemplateCSVData     `json:"csvData"`
	ChartOuput  []models.TemplateChartOutput `json:"chartOutputs"`
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

	var req UpdatedSectionRequest

	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.OldPartIndex < 0 || req.NewPartIndex < 0 || req.OldSectionIndex < 0 || req.NewSectionIndex < -1 || req.NewSectionTitle == "" {
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

		err = util.UpdateSectionInReport(
			req.ItemID,
			req.OldPartIndex,
			req.NewPartIndex,
			req.OldSectionIndex,
			req.NewSectionIndex,
			req.NewSectionTitle,
			sectionContents.Questions,
			sectionContents.TextOutputs,
			sectionContents.CSVData,
			sectionContents.ChartOuput,
			req.DeleteGeneratedOutput,
			userID)

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

		err = util.UpdateSectionInTemplate(
			req.ItemID,
			req.OldPartIndex,
			req.NewPartIndex,
			req.OldSectionIndex,
			req.NewSectionIndex,
			req.NewSectionTitle,
			sectionContents.Questions,
			sectionContents.TextOutputs,
			sectionContents.CSVData,
			sectionContents.ChartOuput,
			userID)

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
		Body:       "Section updated successfully in report with ID: " + req.ItemID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
