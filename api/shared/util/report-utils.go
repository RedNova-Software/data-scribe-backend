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
	"github.com/google/uuid"
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

func GetAllReports(userID string) ([]*models.ReportMetadata, error) {
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
		constants.SharedWithIDsField,
		constants.CreatedField,
		constants.LastModifiedField,
	}

	projectionExpression := strings.Join(fields, ", ")

	// Use FilterExpression for nested attributes
	filterExpression := constants.OwnerUserIDField + " = :userID OR contains(" + constants.SharedWithIDsField + ", :userID)"

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

	reports := []*models.ReportMetadata{}

	for _, item := range result.Items {
		var reportMetadata *models.ReportMetadata

		err = dynamodbattribute.UnmarshalMap(item, &reportMetadata)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		// Needed to extract the SharedWithIDs field
		var report *models.Template
		err = dynamodbattribute.UnmarshalMap(item, &report)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		createReportMetadataSharedWith(reportMetadata, report.SharedWithIDs)

		setReportMetadataOwnerUserName(reportMetadata)

		ensureNonNullReportMetadataFields(reportMetadata)

		reports = append(reports, reportMetadata)
	}

	return reports, nil
}

func SetReportShared(reportID string, userIDs []string, userID string) error {
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

	// Eventually you could refactor this into checking dynamodb first so we don't incur costs for loading data, but this should never be called
	// so it's not worth it right now.
	if !isUserOwnerOfReport(report, userID) {
		return fmt.Errorf("user is not the owner of this report. cannot share with others")
	}

	report.SharedWithIDs = userIDs

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

func ConvertReportToTemplate(reportID, templateTitle, userID string) error {
	templateTableName := os.Getenv(constants.TemplateTable)

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

	if !isUserAuthorizedForReport(report, userID) {
		return fmt.Errorf("user does not have access to this report. cannot convert to template")
	}

	ownerNickName, err := GetUserNickname(userID)

	if err != nil {
		ownerNickName = "*Error Fetching Nickname*"
	}

	// Create template metadata
	newTemplate := &models.Template{
		TemplateID: uuid.New().String(),
		Title:      templateTitle,
		OwnedBy: models.User{
			UserID:       userID,
			UserNickName: ownerNickName,
		},
		SharedWithIDs: make([]string, 0),
		Created:       GetCurrentTime(),
		LastModified:  GetCurrentTime(),
		// Create empty parts for filling
		Parts: make([]models.TemplatePart, 0),
	}

	var templateParts []models.TemplatePart

	// Extract report parts into template
	for i := range report.Parts {
		reportPart := &report.Parts[i]

		templatePart := &models.TemplatePart{
			Title:    reportPart.Title,
			Sections: make([]models.TemplateSection, 0),
		}

		// Iterate through each section in reportPart
		for _, reportSection := range reportPart.Sections {
			templateSection := models.TemplateSection{
				Title:       reportSection.Title,
				Questions:   make([]models.TemplateQuestion, len(reportSection.Questions)),
				TextOutputs: make([]models.TemplateTextOutput, len(reportSection.TextOutputs)),
			}

			// Convert ReportQuestions to TemplateQuestions
			for j, question := range reportSection.Questions {
				templateSection.Questions[j] = models.TemplateQuestion{
					Label:    question.Label,
					Question: question.Question,
				}
			}

			// Convert ReportTextOutputs to TemplateTextOutputs
			for k, textOutput := range reportSection.TextOutputs {
				templateSection.TextOutputs[k] = models.TemplateTextOutput{
					Title: textOutput.Title,
					Type:  textOutput.Type,
					Input: textOutput.Input,
				}
			}

			// Append the converted section to templatePart
			templatePart.Sections = append(templatePart.Sections, templateSection)
		}

		// Append the converted part to templateParts
		templateParts = append(templateParts, *templatePart)
	}

	newTemplate.Parts = templateParts

	templateAV, err := dynamodbattribute.MarshalMap(newTemplate)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      templateAV,
		TableName: aws.String(templateTableName),
	}

	_, err = dynamoDBClient.PutItem(input)

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

	if report.SharedWithIDs == nil {
		report.SharedWithIDs = []string{}
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

func ensureNonNullReportMetadataFields(report *models.ReportMetadata) {
	if report.SharedWith == nil {
		report.SharedWith = []models.User{}
	}
}

func createReportMetadataSharedWith(reportMetadata *models.ReportMetadata, sharedWithIDs []string) {
	// Create a slice to hold User structs
	var sharedWithUsers []models.User

	// Iterate through each UserID in the SharedWith field of the report
	for _, userID := range sharedWithIDs {
		// Call the getUserName function to get the UserNickName
		userNickName, err := GetUserNickname(userID)

		if err != nil {
			userNickName = "*Error Fetching Nickname*"
		}

		// Create a User struct for each UserID
		user := models.User{
			UserID:       userID,
			UserNickName: userNickName,
		}

		// Append the User struct to the slice
		sharedWithUsers = append(sharedWithUsers, user)
	}

	// Update the SharedWith field of the report with the slice of User structs
	reportMetadata.SharedWith = sharedWithUsers
}

func setReportMetadataOwnerUserName(reportMetadata *models.ReportMetadata) {
	userNickName, err := GetUserNickname(reportMetadata.OwnedBy.UserID)

	if err != nil {
		userNickName = "*Error Fetching Nickname*"
	}

	reportMetadata.OwnedBy.UserNickName = userNickName
}

// isUserAuthorizedForReport checks if a given userID is the owner of the report or is in the shared users list
func isUserAuthorizedForReport(report *models.Report, userID string) bool {
	// Check if the user is the owner
	if report.OwnedBy.UserID == userID {
		return true
	}

	// Check if the user is in the shared list
	for _, id := range report.SharedWithIDs {
		if id == userID {
			return true
		}
	}

	return false
}

func isUserOwnerOfReport(report *models.Report, userID string) bool {
	// Check if the user is the owner
	return report.OwnedBy.UserID == userID
}
