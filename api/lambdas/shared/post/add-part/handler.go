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

type AddPartRequest struct {
	ItemType  constants.ItemType `json:"itemType"`
	ItemID    string             `json:"itemID"`
	Index     uint16             `json:"partIndex"`
	PartTitle string             `json:"partTitle"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req AddPartRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: " + err.Error(),
		}, nil
	}

	if req.ItemType == "" || req.ItemID == "" || req.PartTitle == "" || req.Index < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, partTitle, and partIndex are required.",
		}, nil
	}

	updatedIndices, err := util.ModifyItemPartIndices(req.ItemType, req.ItemID, req.Index, true) // Increment all index values equal and above this part

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	newPart := models.ReportPart{
		Title:    req.PartTitle,
		Index:    req.Index,
		Sections: []models.ReportSection{},
	}

	err = util.AddPartToItem(req.ItemType, req.ItemID, newPart)

	if err != nil {
		if updatedIndices {
			util.ModifyItemPartIndices(req.ItemType, req.ItemID, req.Index, false) // Return indices of parts back to normal. Maybe handle this response soon
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    constants.CorsHeaders,
			Body:       "Internal Server Error: " + err.Error(),
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
