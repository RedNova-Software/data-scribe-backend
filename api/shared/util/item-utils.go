package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func SetItemShared(itemType constants.ItemType, itemID string, userIDs []string, userID string) error {

	isOwner, err := isUserOwnerOfItem(itemType, itemID, userID)

	if err != nil {
		return fmt.Errorf("error checking if user is owner of item: %v", err)
	}

	if !isOwner {
		return fmt.Errorf("user is not the owner of this item. cannot share with others")
	}

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	var tableName string
	var av map[string]*dynamodb.AttributeValue

	if itemType == constants.Report {
		tableName = os.Getenv(constants.ReportTable)

		report, err := GetReport(itemID, userID)

		if err != nil {
			return fmt.Errorf("error getting template from DynamoDB: %v", err)
		}

		if report == nil {
			return fmt.Errorf("template not found: %v", err)
		}

		report.SharedWithIDs = userIDs

		av, err = dynamodbattribute.MarshalMap(report)
		if err != nil {
			return err
		}

	} else if itemType == constants.Template {
		tableName = os.Getenv(constants.TemplateTable)

		template, err := GetTemplate(itemID, userID)

		if err != nil {
			return fmt.Errorf("error getting template from DynamoDB: %v", err)
		}

		if template == nil {
			return fmt.Errorf("template not found: %v", err)
		}

		template.SharedWithIDs = userIDs

		av, err = dynamodbattribute.MarshalMap(template)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	updateInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDBClient.PutItem(updateInput)
	if err != nil {
		return err
	}

	return nil
}

func UpdateItemTitle(itemType constants.ItemType, itemID, newTitle string, userID string) error {

	isAuthorized, err := isUserAuthorizedForItem(itemType, itemID, userID)

	if err != nil {
		return fmt.Errorf("error getting authentication status for item: %v", err)
	}

	if !isAuthorized {
		return fmt.Errorf("user is not authorized for item")
	}

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
		UpdateExpression: aws.String("set " + constants.TitleField + " = :t, " + constants.LastModifiedField + " = :lm"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(newTitle),
			},
			":lm": {
				N: aws.String(strconv.FormatInt(currentUnixTime, 10)),
			},
		},
	}

	// Update the item in DynamoDB
	_, err = dynamoDBClient.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}

	return nil
}

func SetItemDeleted(itemType constants.ItemType, itemID, userID string) error {

	isAuthorized, err := isUserOwnerOfItem(itemType, itemID, userID)

	if err != nil {
		return fmt.Errorf("error getting authentication status for item: %v", err)
	}

	if !isAuthorized {
		return fmt.Errorf("user is not authorized for item")
	}

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

	// Set deletion time to 30 days from now
	deletionTime := time.Now().Add(30 * 24 * time.Hour).Unix()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			itemKey: {
				S: aws.String(itemID),
			},
		},
		UpdateExpression: aws.String("set " + constants.IsDeletedField + " :isDel, " + constants.DeleteAtField + " = :delAt"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isDel": {
				BOOL: aws.Bool(true),
			},
			":delAt": {
				N: aws.String(strconv.FormatInt(deletionTime, 10)),
			},
		},
	}

	// Update the item in DynamoDB
	_, err = dynamoDBClient.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}

	return nil
}

func isUserOwnerOfItem(itemType constants.ItemType, itemID, userID string) (bool, error) {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return false, fmt.Errorf("error getting dynamodb client: %v", err)
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
		return false, fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			itemKey: {
				S: aws.String(itemID),
			},
		},
		ProjectionExpression: aws.String(constants.OwnedByUserIDField),
	})
	if err != nil {
		return false, fmt.Errorf("error getting item from DynamoDB: %v", err)
	}

	// Check if the item was found.
	if result.Item == nil {
		return false, nil // Item not found
	}

	if itemType == constants.Report {
		var report *models.Report

		err = dynamodbattribute.UnmarshalMap(result.Item, &report)

		if err != nil {
			return false, fmt.Errorf("error unmarshalling dynamo item into report: %v", err)
		}

		if report.OwnedBy.UserID == userID {
			return true, nil
		}

	} else if itemType == constants.Template {
		var template *models.Template

		err = dynamodbattribute.UnmarshalMap(result.Item, &template)

		if err != nil {
			return false, fmt.Errorf("error unmarshalling dynamo item into template: %v", err)
		}

		if template.OwnedBy.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

func isUserAuthorizedForItem(itemType constants.ItemType, itemID, userID string) (bool, error) {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return false, fmt.Errorf("error getting dynamodb client: %v", err)
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
		return false, fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	// Fields to retrieve
	fields := []string{
		constants.OwnedByUserIDField,
		constants.SharedWithIDsField,
	}

	projectionExpression := strings.Join(fields, ", ")

	// Retrieve the item from DynamoDB
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			itemKey: {
				S: aws.String(itemID),
			},
		},
		ProjectionExpression: aws.String(projectionExpression),
	})
	if err != nil {
		return false, err
	}

	if result.Item == nil {
		return false, fmt.Errorf("item not found")
	}

	if itemType == constants.Report {
		var report *models.Report

		err = dynamodbattribute.UnmarshalMap(result.Item, &report)

		if err != nil {
			return false, fmt.Errorf("error unmarshalling dynamo item into report: %v", err)
		}

		if report.OwnedBy.UserID != userID {
			return false, nil
		}

		// Check if the user is the owner
		if report.OwnedBy.UserID == userID {
			return true, nil
		}

		// Check if the user is in the shared list
		for _, sharedUserID := range report.SharedWithIDs {
			if sharedUserID == userID {
				return true, nil
			}
		}

	} else if itemType == constants.Template {
		var template *models.Template

		err = dynamodbattribute.UnmarshalMap(result.Item, &template)

		if err != nil {
			return false, fmt.Errorf("error unmarshalling dynamo item into template: %v", err)
		}

		if template.OwnedBy.UserID != userID {
			return false, nil
		}

		// Check if the user is the owner
		if template.OwnedBy.UserID == userID {
			return true, nil
		}

		// Check if the user is in the shared list
		for _, sharedUserID := range template.SharedWithIDs {
			if sharedUserID == userID {
				return true, nil
			}
		}
	}

	return false, nil
}
