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

type DeletePartRequest struct {
	ItemType constants.ItemType `json:"itemType"`
	ItemID   string             `json:"itemID"`
	Index    int                `json:"partIndex"`
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

	var req DeletePartRequest
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.Index < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, and partIndex are required.",
		}, nil
	}

	if req.ItemType == constants.Report {
		err = util.DeletePartFromItem(constants.Report, req.ItemID, req.Index, userID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    constants.CorsHeaders,
				Body:       "Internal Server Error: " + err.Error(),
			}, nil
		}
	} else if req.ItemType == constants.Template {
		err = util.DeletePartFromItem(constants.Template, req.ItemID, req.Index, userID)
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
		Body:       "Part added successfully to item with ID: " + req.ItemID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
