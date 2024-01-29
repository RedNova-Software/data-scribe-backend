package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ConvertItemRequest struct {
	ItemType constants.ItemType `json:"itemType"`
	ItemID   string             `json:"itemID"`
	Title    string             `json:"title"`
}

type ConvertTemplateToReportItems struct {
	ReportType string `json:"reportType"`
	City       string `json:"city"`
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

	var req ConvertItemRequest
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.Title == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, and title are required.",
		}, nil
	}

	if req.ItemType == constants.Report {
		err = util.ConvertReportToTemplate(req.ItemID, req.Title, userID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    constants.CorsHeaders,
				Body:       "Internal Server Error: " + err.Error(),
			}, nil
		}
	} else if req.ItemType == constants.Template {
		var items ConvertTemplateToReportItems
		err = json.Unmarshal([]byte(request.Body), &items)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    constants.CorsHeaders,
				Body:       "Bad Request: " + err.Error(),
			}, nil
		}

		err = util.ConvertTemplateToReport(req.ItemID, req.Title, items.City, items.ReportType, userID)
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
		Body:       "Part added successfully to report with ID: " + req.ItemID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
