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

type EditPartRequest struct {
	ItemType  constants.ItemType `json:"itemType"`
	ItemID    string             `json:"itemID"`
	OldIndex  uint16             `json:"oldPartIndex"`
	NewIndex  uint16             `json:"newPartIndex"`
	PartTitle string             `json:"partTitle"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req EditPartRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.PartTitle == "" || req.OldIndex < 0 || req.NewIndex < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, partTitle, and partIndex are required.",
		}, nil
	}

	if req.ItemType == constants.Report {
		err = util.UpdatePartInItem(constants.Report, req.ItemID, req.OldIndex, req.PartTitle, req.NewIndex)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    constants.CorsHeaders,
				Body:       "Internal Server Error: " + err.Error(),
			}, nil
		}
	} else if req.ItemType == constants.Template {
		err = util.UpdatePartInItem(constants.Template, req.ItemID, req.OldIndex, req.PartTitle, req.NewIndex)

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
		Body:       "Part edited successfully",
	}, nil
}

func main() {
	lambda.Start(Handler)
}