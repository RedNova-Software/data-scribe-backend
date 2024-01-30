package main

import (
	"api/shared/constants"
	"api/shared/util"
	"context"
	"net/http"
	"strconv"

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
	partIndexString := request.QueryStringParameters["partIndex"]
	sectionIndexString := request.QueryStringParameters["sectionIndex"]

	partIndex, err := strconv.Atoi(partIndexString)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: unable to parse partIndex. ensure it is an int",
		}, nil
	}

	sectionIndex, err := strconv.Atoi(sectionIndexString)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: unable to parse sectionIndex. ensure it is an int",
		}, nil
	}

	if itemID == "" || itemType == "" || partIndex < 0 || sectionIndex < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    constants.CorsHeaders,
			Body:       "Bad Request: itemType, itemID, and partIndex are required.",
		}, nil
	}

	if itemType == constants.Report || itemType == constants.Template {
		err = util.DeleteSectionFromItem(itemType, itemID, partIndex, sectionIndex, userID)
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
		Body:       "Part delete successfully from item with id: " + itemID,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
