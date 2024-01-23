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
	"github.com/google/uuid"
)

type CreateTemplateRequest struct {
	Title string `json:"title"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	var req CreateTemplateRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	if req.Title == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: title is required.",
			Headers:    constants.CorsHeaders,
		}, nil
	}

	templateID := uuid.New().String()

	template := models.Template{
		TemplateID: templateID,
		Title:      req.Title,
		Parts:      make([]models.TemplatePart, 0),
	}

	err = util.PutNewTemplate(template)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error: " + err.Error(),
			Headers:    constants.CorsHeaders,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Empty template created successfully with ID: " + templateID,
		Headers:    constants.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
