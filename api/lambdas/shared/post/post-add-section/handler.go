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

type AddSectionToPartRequest struct {
	ItemType     constants.ItemType `json:"itemType"`
	ItemID       string             `json:"itemID"`
	PartIndex    uint16             `json:"partIndex"`
	SectionIndex uint16             `json:"sectionIndex"`
	SectionTitle string             `json:"sectionTitle"`
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
	var req AddSectionToPartRequest

	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.PartIndex < 0 || req.SectionTitle == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, sectionTitle, and partIndex are required.",
		}, nil
	}

	if req.ItemType == constants.Report {
		var contents ReportSectionContents

		err := json.Unmarshal([]byte(request.Body), &contents)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    constants.CorsHeaders,
				Body:       "Bad Request: " + err.Error(),
			}, nil
		}

		newSection := models.ReportSection{
			Title:       req.SectionTitle,
			Index:       req.SectionIndex,
			Questions:   contents.Questions,
			TextOutputs: contents.TextOutputs,
		}
		err = util.AddSectionToReportPart(req.ItemID, req.PartIndex, newSection)
	} else if req.ItemType == constants.Template {
		var contents TemplateSectionContents

		err := json.Unmarshal([]byte(request.Body), &contents)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    constants.CorsHeaders,
				Body:       "Bad Request: " + err.Error(),
			}, nil
		}

		newSection := models.TemplateSection{
			Title:       req.SectionTitle,
			Index:       req.SectionIndex,
			Questions:   contents.Questions,
			TextOutputs: contents.TextOutputs,
		}

		err = util.AddSectionToTemplatePart(req.ItemID, req.PartIndex, newSection)
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType must be 'report' or 'template' ",
		}, nil
	}

	updatedIndices, err := util.ModifyPartSectionIndices(req.ItemType, req.ItemID, req.PartIndex, req.SectionIndex, true) // Increment all index values equal and above this section

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	if err != nil {
		if updatedIndices {
			util.ModifyPartSectionIndices(req.ItemType, req.ItemID, req.PartIndex, req.SectionIndex, false) // Return indices of parts back to normal. Maybe handle this response soon
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
		Body:       "Section added successfully to report with ID: " + req.ItemID + "and part with index: " + fmt.Sprint(req.PartIndex),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
