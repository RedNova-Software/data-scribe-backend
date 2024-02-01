package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID, err := util.ExtractUserID(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	// Extract the query string parameterss
	itemType := constants.ItemType(request.QueryStringParameters["itemType"])
	itemID := request.QueryStringParameters["itemID"]

	if itemID == "" || itemType == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType and itemID are required.",
		}, nil
	}

	if itemType == constants.Report || itemType == constants.Template {
		err = util.SetItemDeleted(itemType, itemID, true, userID)
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
		Body:       "Item marked for deletion successfully",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
