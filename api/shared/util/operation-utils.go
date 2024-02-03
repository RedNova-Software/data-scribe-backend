package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func CreateOperation(operationID string) error {
	tableName := os.Getenv(constants.OperationTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
	}

	operation := models.Operation{
		OperationID: operationID,
		Completed:   false,
		DeleteAt:    time.Now().Add(24 * time.Hour).Unix(), // Set to delete 24 hours from now
	}

	item, err := dynamodbattribute.MarshalMap(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = dynamoDBClient.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

func SetOperationCompleted(operationID string) error {
	tableName := os.Getenv(constants.OperationTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.OperationIDField: {
				S: aws.String(operationID),
			},
		},
		UpdateExpression: aws.String("set " + constants.OperationCompletedField + " = :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				BOOL: aws.Bool(true),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err = dynamoDBClient.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

func GetOperationStatus(operationID string) (bool, error) {
	tableName := os.Getenv(constants.OperationTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return false, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.OperationIDField: {
				S: aws.String(operationID),
			},
		},
		ProjectionExpression: aws.String(constants.OperationCompletedField),
	}

	result, err := dynamoDBClient.GetItem(input)
	if err != nil {
		return false, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return false, nil // Assuming false for non-existent operations
	}

	// Assuming 'Complete' attribute exists and is a boolean
	completeAttr := result.Item[constants.OperationCompletedField]
	if completeAttr == nil || completeAttr.BOOL == nil {
		return false, fmt.Errorf("complete attribute is missing or not a boolean")
	}

	return *completeAttr.BOOL, nil
}
