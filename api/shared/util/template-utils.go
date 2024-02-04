package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
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

func GetTemplate(templateID string, userID string) (*models.Template, error) {
	isAuthorized, err := isUserAuthorizedForItem(constants.Template, templateID, userID)

	if err != nil {
		return nil, fmt.Errorf("error getting authentication status for item: %v", err)
	}

	if !isAuthorized {
		return nil, fmt.Errorf("user is not authorized for item")
	}

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

	var template *models.Template

	err = dynamodbattribute.UnmarshalMap(result.Item, &template)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling dynamo item into template: %v", err)
	}

	if template.IsDeleted {
		return nil, fmt.Errorf("report is deleted. cannot fetch")
	}

	// Ensure all nil parts and nil sections are returned as an empty list
	// This is an annoyance due to the way dynamodb marshalls empty lists
	// When we create an empty report, the parts will be null in dynamodb.
	// Same for an empty part, the sections will be null. So, to return a list
	// to the frontend, we need to set it explicitly here.
	// https://github.com/aws/aws-sdk-go/issues/682
	ensureNonNullTemplateFields(template)

	return template, nil
}

func GetAllTemplates(userID string, deletedTemplatesOnly bool) ([]*models.TemplateMetadata, error) {
	tableName := os.Getenv(constants.TemplateTable)

	// Create a new DynamoDB session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.USEast2),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new session: %v", err)
	}

	dynamoDBClient := dynamodb.New(sess)

	// Fields to retrieve
	fields := []string{
		constants.TemplateIDField,
		constants.TitleField,
		constants.OwnedByUserIDField,
		constants.SharedWithIDsField,
		constants.CreatedAtField,
		constants.LastModifiedAtField,
	}

	projectionExpression := strings.Join(fields, ", ")

	var filterExpression string

	if deletedTemplatesOnly {
		filterExpression =
			constants.OwnedByUserIDField +
				" = :userID AND " +
				constants.IsDeletedField + " = :isDeleted"
	} else {
		filterExpression = "(" +
			constants.OwnedByUserIDField +
			" = :userID OR contains(" + constants.SharedWithIDsField + ", :userID)) AND " +
			constants.IsDeletedField + " = :isDeleted"
	}

	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String(filterExpression),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				S: aws.String(userID),
			},
			":isDeleted": {
				BOOL: aws.Bool(deletedTemplatesOnly),
			},
		},
		ProjectionExpression: aws.String(projectionExpression),
	}

	result, err := dynamoDBClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error querying DynamoDB table: %v", err)
	}

	templates := []*models.TemplateMetadata{}

	for _, item := range result.Items {
		var templateMetadata *models.TemplateMetadata
		err = dynamodbattribute.UnmarshalMap(item, &templateMetadata)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		// Needed to extract the SharedWithIDs field
		var template *models.Template
		err = dynamodbattribute.UnmarshalMap(item, &template)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		createTemplateMetadataSharedWith(templateMetadata, template.SharedWithIDs)

		setTemplateMetadataOwnerUserName(templateMetadata)

		ensureNonNullTemplateMetadataFields(templateMetadata)

		templates = append(templates, templateMetadata)
	}

	return templates, nil
}

func ConvertTemplateToReport(templateID, reportTitle, reportCity, reportType, userID string) error {
	isAuthorized, err := isUserAuthorizedForItem(constants.Template, templateID, userID)

	if err != nil {
		return fmt.Errorf("error getting authentication status for item: %v", err)
	}

	if !isAuthorized {
		return fmt.Errorf("user is not authorized for item")
	}

	reportTableName := os.Getenv(constants.ReportTable)

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	template, err := GetTemplate(templateID, userID)

	if err != nil {
		return fmt.Errorf("error getting template from DynamoDB: %v", err)
	}

	if template == nil {
		return fmt.Errorf("template not found: %v", err)
	}

	ownerNickName, err := GetUserNickname(userID)

	if err != nil {
		ownerNickName = "*Error Fetching Nickname*"
	}

	// Create report metadata
	newReport := &models.Report{
		ReportID:   uuid.New().String(),
		IsDeleted:  false,
		Title:      reportTitle,
		City:       reportCity,
		ReportType: reportType,
		OwnedBy: models.User{
			UserID:       userID,
			UserNickName: ownerNickName,
		},
		SharedWithIDs:  make([]string, 0),
		CreatedAt:      GetCurrentTime(),
		LastModifiedAt: GetCurrentTime(),
		// Create empty parts for filling
		Parts:           make([]models.ReportPart, 0),
		CSVID:           "no-csv-id", // Needed because CSVID is a GSI, and cannot be null
		CSVColumnsS3Key: "no-csv-s3-key",
	}

	var reportParts []models.ReportPart

	// Extract template parts into report
	for i := range template.Parts {
		templatePart := &template.Parts[i]

		reportPart := &models.ReportPart{
			Title:    templatePart.Title,
			Sections: make([]models.ReportSection, 0),
		}

		// Iterate through each section in templatePart
		for _, templateSection := range templatePart.Sections {
			reportSection := models.ReportSection{
				Title:           templateSection.Title,
				OutputGenerated: false,
				Questions:       make([]models.ReportQuestion, len(templateSection.Questions)),
				CSVData:         make([]models.ReportCSVData, len(templateSection.CSVData)),
				TextOutputs:     make([]models.ReportTextOutput, len(templateSection.TextOutputs)),
				ChartOutputs:    make([]models.ReportChartOutput, len(templateSection.ChartOutputs)),
			}

			// Convert TemplateQuestions to ReportQuestions
			for j, question := range templateSection.Questions {
				reportSection.Questions[j] = models.ReportQuestion{
					Label:    question.Label,
					Question: question.Question,
					Answer:   "", // Initialize with empty answer
				}
			}

			// Convert TemplateTextOutputs to ReportTextOutputs
			for k, textOutput := range templateSection.TextOutputs {
				reportSection.TextOutputs[k] = models.ReportTextOutput{
					Title:  textOutput.Title,
					Type:   textOutput.Type,
					Input:  textOutput.Input,
					Result: "", // Initialize with empty result
				}
			}

			// Convert TemplateCsvData to ReportCsvData
			for l, data := range templateSection.CSVData {
				reportSection.CSVData[l] = models.ReportCSVData{
					Label: data.Label,
					ConfigOneDim: models.ReportOneDimConfig{
						AggregateValueLabel: data.ConfigOneDim.AggregateValueLabel,
						Description:         data.ConfigOneDim.Description,
						OperationType:       data.ConfigOneDim.OperationType,
						Column:              "",
						AcceptedValues:      make([]string, 0),
					},
				}
			}

			// Convert TemplateChartOutputs to ReportChartOutputs
			for h, chart := range templateSection.ChartOutputs {
				newDependentColumns := make([]models.ReportOneDimConfig, len(chart.Config.DependentColumns))

				for u, templateDependentColumn := range chart.Config.DependentColumns {
					newDependentColumns[u] = models.ReportOneDimConfig{
						AggregateValueLabel: templateDependentColumn.AggregateValueLabel,
						Description:         templateDependentColumn.Description,
						OperationType:       templateDependentColumn.OperationType,
						Column:              "",
						AcceptedValues:      make([]string, 0),
					}
				}

				reportSection.ChartOutputs[h] = models.ReportChartOutput{
					Title:         chart.Title,
					Type:          chart.Type,
					XAxisTitle:    chart.XAxisTitle,
					YAxisTitle:    chart.YAxisTitle,
					CartesianGrid: chart.CartesianGrid,
					Config: models.ReportTwoDimConfig{
						IndependentColumnLabel: chart.Config.IndependentColumnLabel,
						DependentColumns:       newDependentColumns,
					},
					Results: make([]map[string]interface{}, 0),
				}
			}

			// Append the converted section to reportPart
			reportPart.Sections = append(reportPart.Sections, reportSection)
		}

		// Append the converted part to reportParts
		reportParts = append(reportParts, *reportPart)
	}

	newReport.Parts = reportParts

	reportAV, err := dynamodbattribute.MarshalMap(newReport)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      reportAV,
		TableName: aws.String(reportTableName),
	}

	_, err = dynamoDBClient.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}

func ensureNonNullTemplateFields(template *models.Template) {
	// Check if Parts is nil, if so, initialize it as an empty slice
	if template.Parts == nil {
		template.Parts = []models.TemplatePart{}
	}

	if template.SharedWithIDs == nil {
		template.SharedWithIDs = []string{}
	}

	// Iterate over each part
	for i := range template.Parts {
		part := &template.Parts[i]

		// Check if Sections is nil, if so, initialize it as an empty slice
		if part.Sections == nil {
			part.Sections = []models.TemplateSection{}
		}

		// Iterate over each section
		for j := range part.Sections {
			section := &part.Sections[j]

			// Check if Questions is nil, if so, initialize it as an empty slice
			if section.Questions == nil {
				section.Questions = []models.TemplateQuestion{}
			}

			// Check if TextOutputs is nil, if so, initialize it as an empty slice
			if section.TextOutputs == nil {
				section.TextOutputs = []models.TemplateTextOutput{}
			}
		}
	}
}

func ensureNonNullTemplateMetadataFields(template *models.TemplateMetadata) {

	if template.SharedWith == nil {
		template.SharedWith = []models.User{}
	}
}

func createTemplateMetadataSharedWith(templateMetadata *models.TemplateMetadata, sharedWithIDs []string) {
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
	templateMetadata.SharedWith = sharedWithUsers
}

func setTemplateMetadataOwnerUserName(templateMetadata *models.TemplateMetadata) {
	userNickName, err := GetUserNickname(templateMetadata.OwnedBy.UserID)

	if err != nil {
		userNickName = "*Error Fetching Nickname*"
	}

	templateMetadata.OwnedBy.UserNickName = userNickName
}
