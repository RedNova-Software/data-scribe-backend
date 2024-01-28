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

	if !isUserAuthorizedForTemplate(template, userID) {
		return nil, fmt.Errorf("user does not have access to this report")
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

func GetAllTemplates(userID string) ([]*models.TemplateMetadata, error) {
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
		constants.OwnerUserIDField,
		constants.SharedWithIDsField,
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

		createTemplateSharedWith(templateMetadata, template.SharedWithIDs)

		ensureNonNullTemplateMetadataFields(templateMetadata)

		templates = append(templates, templateMetadata)
	}

	return templates, nil
}

func SetTemplateShared(templateID string, userIDs []string, userID string) error {
	tableName := os.Getenv(constants.TemplateTable)

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	template, err := GetTemplate(templateID, userID)

	if err != nil {
		return fmt.Errorf("error getting report from DynamoDB: %v", err)
	}

	if template == nil {
		return fmt.Errorf("report not found: %v", err)
	}

	template.SharedWithIDs = userIDs

	av, err := dynamodbattribute.MarshalMap(template)
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

func createTemplateSharedWith(template *models.TemplateMetadata, sharedWithIDs []string) {
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
	template.SharedWith = sharedWithUsers
}

// isUserAuthorizedForReport checks if a given userID is the owner of the template or is in the shared users list
func isUserAuthorizedForTemplate(template *models.Template, userID string) bool {
	// Check if the user is the owner
	if template.OwnedBy.UserID == userID {
		return true
	}

	// Check if the user is in the shared list
	for _, id := range template.SharedWithIDs {
		if id == userID {
			return true
		}
	}

	return false
}
