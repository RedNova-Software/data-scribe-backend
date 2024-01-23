package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetTemplate(tableName, keyName, keyValue string) (*models.Template, error) {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)

	if err != nil {
		return nil, err
	}

	// Create a DynamoDB input structure for the GetItem operation.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(keyValue),
			},
		},
	}

	result, err := dynamoDBClient.GetItem(input)
	// Execute the GetItem operation.
	if err != nil {
		return nil, fmt.Errorf("error getting item from DynamoDB: %v", err)
	}

	// Check if the item was found.
	if result.Item == nil {
		return nil, nil // Item not found
	}

	var template models.Template

	err = dynamodbattribute.UnmarshalMap(result.Item, &template)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling dynamo item into template: %v", err)
	}

	return &template, nil
}

func GetAllTemplates(tableName string) ([]models.Report, error) {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)

	if err != nil {
		return nil, err
	}

	// Fields to retrieve
	fields := []string{
		constants.TemplateIDField,
		constants.TitleField,
	}

	projectionExpression := strings.Join(fields, ", ")

	// Create a DynamoDB ScanInput with the ProjectionExpression
	input := &dynamodb.ScanInput{
		TableName:            aws.String(tableName),
		ProjectionExpression: aws.String(projectionExpression),
	}

	result, err := dynamoDBClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error scanning DynamoDB table: %v", err)
	}

	reports := []models.Report{}

	for _, item := range result.Items {
		var report models.Report
		err = dynamodbattribute.UnmarshalMap(item, &report)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		reports = append(reports, report)
	}

	return reports, nil
}
