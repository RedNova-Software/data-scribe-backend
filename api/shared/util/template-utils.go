package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutNewTemplate(template models.Template) error {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
	}

	templateAV, err := dynamodbattribute.MarshalMap(template)
	if err != nil {
		return err
	}

	// Needed to set "Parts" to empty list
	// For more info, see https://github.com/aws/aws-sdk-go/issues/682
	templateAV[constants.PartsField] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

	input := &dynamodb.PutItemInput{
		Item:      templateAV,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDBClient.PutItem(input)

	if err != nil {
		return err
	}
	return nil
}

func GetTemplate(templateID string) (*models.Template, error) {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	keyName := constants.TemplateIDField

	if err != nil {
		return nil, fmt.Errorf("error getting dynamodb client: %v", err)
	}

	// Create a DynamoDB input structure for the GetItem operation.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(templateID),
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

func GetAllTemplates() ([]models.Template, error) {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return nil, fmt.Errorf("error getting dynamodb client: %v", err)
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

	templates := []models.Template{}

	for _, item := range result.Items {
		var template models.Template
		err = dynamodbattribute.UnmarshalMap(item, &template)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		templates = append(templates, template)
	}

	return templates, nil
}
