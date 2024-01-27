package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func AddPartToItem(
	itemType constants.ItemType,
	itemID string,
	partTitle string,
	partIndex uint16,
) error {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	var tableName string

	if itemType == constants.Report {
		tableName := os.Getenv(constants.ReportTable)
		report, err := GetReport(itemID)

		if err != nil {
			return fmt.Errorf("error getting report from DynamoDB: %v", err)
		}

		if report == nil {
			return fmt.Errorf("report not found: %v", err)
		}

		// Check if its nil, as in its empty in dynamodb.
		if report.Parts == nil {
			report.Parts = []models.ReportPart{}
		}

		newPart := models.ReportPart{
			Title:    partTitle,
			Sections: []models.ReportSection{},
		}

		err = insertReportPart(report, newPart, int(partIndex))

		if err != nil {
			return fmt.Errorf("error inserting report part: %v", err)
		}

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

	} else if itemType == constants.Template {
		tableName = os.Getenv(constants.TemplateTable)
		template, err := GetTemplate(itemID)

		if err != nil {
			return fmt.Errorf("error getting template from DynamoDB: %v", err)
		}

		if template == nil {
			return fmt.Errorf("template not found: %v", err)
		}

		// Check if its nil, as in its empty in dynamodb.
		if template.Parts == nil {
			template.Parts = []models.TemplatePart{}
		}

		newPart := models.TemplatePart{
			Title:    partTitle,
			Sections: []models.TemplateSection{},
		}

		err = insertTemplatePart(template, newPart, int(partIndex))

		if err != nil {
			return fmt.Errorf("error inserting report part: %v", err)
		}

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
	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	return nil
}

func UpdatePartInItem(
	itemType constants.ItemType,
	itemID string,
	oldIndex uint16,
	newIndex uint16,
	newTitle string,
) error {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	if itemType == constants.Report {
		tableName := os.Getenv(constants.ReportTable)
		report, err := GetReport(itemID)

		if err != nil {
			return fmt.Errorf("error getting report from DynamoDB: %v", err)
		}

		if report == nil {
			return fmt.Errorf("report not found: %v", err)
		}

		err = moveReportPart(report, int(oldIndex), int(newIndex))

		if err != nil {
			return fmt.Errorf("error moving report part: %v", err)
		}

		report.Parts[newIndex].Title = newTitle

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

	} else if itemType == constants.Template {
		tableName := os.Getenv(constants.TemplateTable)
		template, err := GetTemplate(itemID)

		if err != nil {
			return fmt.Errorf("error getting template from DynamoDB: %v", err)
		}

		if template == nil {
			return fmt.Errorf("template not found: %v", err)
		}

		err = moveTemplatePart(template, int(oldIndex), int(newIndex))

		if err != nil {
			return fmt.Errorf("error moving report part: %v", err)
		}

		template.Parts[newIndex].Title = newTitle

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
	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}
}

func insertReportPart(report *models.Report, part models.ReportPart, index int) error {
	if index < 0 || index > len(report.Parts) {
		// Handle the error or ignore if the index is out of bounds
		return fmt.Errorf("unable to insert part into report. index out of bounds")
	}
	report.Parts = append(report.Parts[:index], append([]models.ReportPart{part}, report.Parts[index:]...)...)
	return nil
}

func moveReportPart(report *models.Report, fromIndex, toIndex int) error {
	if fromIndex < 0 || fromIndex >= len(report.Parts) || toIndex < 0 || toIndex >= len(report.Parts) {
		// Handle the error or ignore if indices are out of bounds
		return fmt.Errorf("unable to move part in report. index out of bounds")
	}

	// Remove the part from the current position
	part := report.Parts[fromIndex]
	report.Parts = append(report.Parts[:fromIndex], report.Parts[fromIndex+1:]...)

	// Reinsert the part at the new position
	report.Parts = append(report.Parts[:toIndex], append([]models.ReportPart{part}, report.Parts[toIndex:]...)...)
	return nil
}

func insertTemplatePart(template *models.Template, part models.TemplatePart, index int) error {
	if index < 0 || index > len(template.Parts) {
		return fmt.Errorf("unable to insert part into template. index out of bounds")
	}
	template.Parts = append(template.Parts[:index], append([]models.TemplatePart{part}, template.Parts[index:]...)...)
	return nil
}

func moveTemplatePart(template *models.Template, fromIndex, toIndex int) error {
	if fromIndex < 0 || fromIndex >= len(template.Parts) || toIndex < 0 || toIndex >= len(template.Parts) {
		// Handle the error or ignore if indices are out of bounds
		return fmt.Errorf("unable to move part in template. index out of bounds")
	}

	// Remove the part from the current position
	part := template.Parts[fromIndex]
	template.Parts = append(template.Parts[:fromIndex], template.Parts[fromIndex+1:]...)

	// Reinsert the part at the new position
	template.Parts = append(template.Parts[:toIndex], append([]models.TemplatePart{part}, template.Parts[toIndex:]...)...)
	return nil
}
