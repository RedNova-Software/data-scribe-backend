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

func PutNewReport(report models.Report) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	reportAV, err := dynamodbattribute.MarshalMap(report)
	if err != nil {
		return err
	}

	// Needed to set "Parts" to empty list
	// For more info, see https://github.com/aws/aws-sdk-go/issues/682
	reportAV[constants.PartsField] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

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

func GetReport(reportID string, userID string) (*models.Report, error) {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	const keyName = constants.ReportIDField

	if err != nil {
		return nil, fmt.Errorf("error getting dynamodb client: %v", err)
	}

	// Create a DynamoDB input structure for the GetItem operation.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(reportID),
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

	var report *models.Report

	err = dynamodbattribute.UnmarshalMap(result.Item, &report)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling dynamo item into report: %v", err)
	}

	if !isUserAuthorizedForReport(report, userID) {
		return nil, fmt.Errorf("user does not have access to this report")
	}

	// Ensure all nil parts and nil sections are returned as an empty list
	// This is an annoyance due to the way dynamodb marshalls empty lists
	// When we create an empty report, the parts will be null in dynamodb.
	// Same for an empty part, the sections will be null. So, to return a list
	// to the frontend, we need to set it explicitly here.
	// https://github.com/aws/aws-sdk-go/issues/682
	ensureNonNullReportFields(report)

	return report, nil
}

func GetAllReports(userID string) ([]*models.Report, error) {
	tableName := os.Getenv(constants.ReportTable)

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return nil, fmt.Errorf("error getting dynamodb client: %v", err)
	}

	// Fields to retrieve
	fields := []string{
		constants.ReportIDField,
		constants.ReportTypeField,
		constants.TitleField,
		constants.CityField,
		constants.OwnerUserIDField,
		constants.SharedWithField,
	}

	projectionExpression := strings.Join(fields, ", ")

	// Use FilterExpression for nested attributes
	filterExpression := constants.OwnerUserIDField + " = :userID"

	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String(filterExpression),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				S: aws.String(userID),
			},
		},
		ProjectionExpression: aws.String(projectionExpression),
	}

	result, err := dynamoDBClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error querying DynamoDB table: %v", err)
	}

	reports := []*models.Report{}

	for _, item := range result.Items {
		var report *models.Report
		err = dynamodbattribute.UnmarshalMap(item, &report)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		ensureNonNullReportFields(report)

		reports = append(reports, report)
	}

	return reports, nil
}

func SetReportShared(reportID string, users []models.User, userID string) error {
	tableName := os.Getenv(constants.ReportTable)

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	report, err := GetReport(reportID, userID)

	if err != nil {
		return fmt.Errorf("error getting report from DynamoDB: %v", err)
	}

	if report == nil {
		return fmt.Errorf("report not found: %v", err)
	}

	report.SharedWith = users

	av, err := dynamodbattribute.MarshalMap(report)
	if err != nil {
		return err
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

func ensureNonNullReportFields(report *models.Report) {
	// Check if Parts is nil, if so, initialize it as an empty slice
	if report.Parts == nil {
		report.Parts = []models.ReportPart{}
	}

	if report.SharedWith == nil {
		report.SharedWith = []models.User{}
	}

	// Iterate over each part
	for i := range report.Parts {
		part := &report.Parts[i]

		// Check if Sections is nil, if so, initialize it as an empty slice
		if part.Sections == nil {
			part.Sections = []models.ReportSection{}
		}

		// Iterate over each section
		for j := range part.Sections {
			section := &part.Sections[j]

			// Check if Questions is nil, if so, initialize it as an empty slice
			if section.Questions == nil {
				section.Questions = []models.ReportQuestion{}
			}

			// Check if TextOutputs is nil, if so, initialize it as an empty slice
			if section.TextOutputs == nil {
				section.TextOutputs = []models.ReportTextOutput{}
			}
		}
	}
}

// isUserAuthorizedForReport checks if a given userID is the owner of the report or is in the shared users list
func isUserAuthorizedForReport(report *models.Report, userID string) bool {
	// Check if the user is the owner
	if report.OwnedBy.UserID == userID {
		return true
	}

	// Check if the user is in the shared list
	for _, user := range report.SharedWith {
		if user.UserID == userID {
			return true
		}
	}

	return false
}
