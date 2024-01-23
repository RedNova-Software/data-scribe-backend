package util

import (
	"api/shared/constants"
	"api/shared/models"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutNewReport(tableName string, report models.Report) error {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
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

func AddPartToReport(
	tableName string,
	reportID string,
	newPart models.ReportPart,
) error {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
	}

	newPartAV, err := dynamodbattribute.MarshalMap(newPart)
	if err != nil {
		return err
	}

	newPartAV[constants.SectionsField] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.ReportIDField: {
				S: aws.String(reportID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#p": aws.String(constants.PartsField),
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

// If increment is true, all Part indices equal to or above the newIndex will be incremented.
// If false, everything larger will be decremented.
func ModifyReportPartIndices(tableName string, reportID string, newIndex uint16, increment bool) (bool, error) {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)
	report, err := GetReport(tableName, constants.ReportIDField, reportID)

	if err != nil {
		return false, fmt.Errorf("error getting report from DynamoDB: %v", err)
	}

	if report == nil {
		return false, fmt.Errorf("report not found: %v", err)
	}

	updated := false
	for i, part := range report.Parts {
		if increment && part.Index >= newIndex {
			report.Parts[i].Index++
			updated = true
		} else if !increment && part.Index > newIndex {
			report.Parts[i].Index--
			updated = true
		}
	}

	if updated {
		av, err := dynamodbattribute.MarshalMap(report)
		if err != nil {
			return false, err
		}

		updateInput := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = dynamoDBClient.PutItem(updateInput)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// If increment is true, all Part indices equal to or above the newIndex will be incremented.
// If false, everything larger will be decremented.
func ModifyPartSectionIndices(tableName string, reportID string, partIndex uint16, newSectionIndex uint16, increment bool) (bool, error) {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)
	report, err := GetReport(tableName, constants.ReportIDField, reportID)

	if err != nil {
		return false, fmt.Errorf("error getting report from DynamoDB: %v", err)
	}

	if report == nil {
		return false, fmt.Errorf("report not found: %v", err)
	}

	var part *models.ReportPart

	for _, searchPart := range report.Parts {
		if searchPart.Index == partIndex {
			part = &searchPart
			break
		}
	}

	updated := false
	if part != nil {

		for i, section := range part.Sections {
			if increment && section.Index >= newSectionIndex {
				part.Sections[i].Index++
				updated = true
			} else if !increment && section.Index > newSectionIndex {
				part.Sections[i].Index--
				updated = true
			}
		}
	}

	if updated {
		av, err := dynamodbattribute.MarshalMap(report)
		if err != nil {
			return false, err
		}

		updateInput := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = dynamoDBClient.PutItem(updateInput)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// AddSectionToReportPart adds a Section to a Part with a specific index in a DynamoDB table.
func AddSectionToReportPart(tableName string, reportID string, partIndex uint16, newSection models.ReportSection) error {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
	}

	// Retrieve the Report item
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.ReportIDField: {
				S: aws.String(reportID),
			},
		},
	})
	if err != nil {
		return err
	}

	if result.Item == nil {
		return fmt.Errorf("report not found: %v", err)
	}

	// Unmarshal the Report
	var report models.Report
	err = dynamodbattribute.UnmarshalMap(result.Item, &report)
	if err != nil {
		return fmt.Errorf("unable to unmarshal dynamodb report: %v", err)
	}

	// Find the Part and add the Section
	partFound := false
	for i, part := range report.Parts {
		if part.Index == partIndex {
			report.Parts[i].Sections = append(report.Parts[i].Sections, newSection)
			partFound = true
			break
		}
	}

	if !partFound {
		return errors.New("part not found")
	}

	// Marshal the updated Report back to a map
	updatedReport, err := dynamodbattribute.MarshalMap(report)
	if err != nil {
		return fmt.Errorf("unable to marshall report into dynamodb attribute: %v", err)
	}

	// Update the Report in the DynamoDB table
	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      updatedReport,
	})
	if err != nil {
		return fmt.Errorf("error updating report item in dynamodb: %v", err)
	}

	return nil
}

func GetReport(tableName, keyName, keyValue string) (*models.Report, error) {
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

	var report models.Report

	err = dynamodbattribute.UnmarshalMap(result.Item, &report)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling dynamo item into report: %v", err)
	}

	return &report, nil
}

func GetAllReports(tableName string) ([]models.Report, error) {
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)

	if err != nil {
		return nil, err
	}

	// Fields to retrieve
	fields := []string{
		constants.ReportIDField,
		constants.ReportTypeField,
		constants.TitleField,
		constants.CityField,
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
