package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutNewReport(tableName string, report models.Report) error {
	dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))
	if err != nil {
		return err
	}

	reportAV, err := dynamodbattribute.MarshalMap(report)
	if err != nil {
		return err
	}

	// Needed to set "Parts" to empty list
	// For more info, see https://github.com/aws/aws-sdk-go/issues/682
	reportAV["Parts"] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

	input := &dynamodb.PutItemInput{
		Item:      reportAV,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDBClient.PutItem(input)

	if err != nil {
		return err
	}
	return nil
}

func AddPartToReport(
	tableName string,
	reportID string,
	newPart models.Part,
) error {
	dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))
	if err != nil {
		return err
	}

	newPartAV, err := dynamodbattribute.MarshalMap(newPart)
	if err != nil {
		return err
	}

	newPartAV["Sections"] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ReportID": {
				S: aws.String(reportID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#p": aws.String("Parts"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newPart": {
				L: []*dynamodb.AttributeValue{
					{M: newPartAV},
				},
			},
		},
		UpdateExpression: aws.String("SET #p = list_append(#p, :newPart)"),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	_, err = dynamoDBClient.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetReport(tableName, keyName, keyValue string) (item map[string]*dynamodb.AttributeValue, err error) {
	dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))

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

	return result.Item, nil
}

func GetAllReports(tableName, projectionExpression string) ([]models.Report, error) {
	dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))

	if err != nil {
		return nil, err
	}

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

func newDynamoDBClient(region string) (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}
