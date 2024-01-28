package util

import (
	"api/shared/constants"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func UpdateItemTitle(itemType constants.ItemType, itemID, newTitle string, userID string) error {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	var tableName string
	var itemKey string

	if itemType == constants.Report {
		tableName = os.Getenv(constants.ReportTable)
		itemKey = constants.ReportIDField

	} else if itemType == constants.Template {
		tableName = os.Getenv(constants.TemplateTable)
		itemKey = constants.TemplateIDField
	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	currentUnixTime := GetCurrentTime()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			itemKey: {
				S: aws.String(itemID),
			},
		},
		UpdateExpression: aws.String("set Title = :t, LastModified = :lm"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(newTitle),
			},
			":lm": {
				N: aws.String(strconv.FormatInt(currentUnixTime, 10)),
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
