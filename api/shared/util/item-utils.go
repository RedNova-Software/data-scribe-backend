package util

import (
	"api/shared/constants"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func EditItemTitle(itemType constants.ItemType, itemID, newTitle string) error {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	var tableName string
	var itemKey string

	if itemType == constants.Report {
		tableName = constants.ReportTable
		itemKey = constants.ReportIDField

	} else if itemType == constants.Template {
		tableName = constants.TemplateTable
		itemKey = constants.TemplateIDField
	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	// Prepare the update input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			itemKey: {
				S: aws.String(itemID),
			},
		},
		UpdateExpression: aws.String("set Title = :t"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(newTitle),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	// Update the item in DynamoDB
	_, err = dynamoDBClient.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}

	return nil
}
