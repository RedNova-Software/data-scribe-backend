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

type UpdateReportTitleRequest struct {
	ItemType constants.ItemType `json:"itemType"`
	ItemID   string             `json:"itemID"`
	NewTitle string             `json:"newTitle"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req UpdateReportTitleRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.NewTitle == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, and newTitle are required.",
		}, nil
	}

	if req.ItemType == constants.Report {
		err = util.UpdateItemTitle(constants.Report, req.ItemID, req.NewTitle)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    constants.CorsHeaders,
				Body:       "Internal Server Error: " + err.Error(),
			}, nil
		}
	} else if req.ItemType == constants.Template {
		err = util.UpdateItemTitle(constants.Template, req.ItemID, req.NewTitle)

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
		Body:       "Title updated successfully",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
