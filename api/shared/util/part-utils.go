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
	partIndex int,
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

		err = insertReportPart(report, newPart, partIndex)

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

		err = insertTemplatePart(template, newPart, partIndex)

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
	oldIndex int,
	newIndex int,
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

		report.Parts[oldIndex].Title = newTitle

		if oldIndex != newIndex {
			err = moveReportPart(report, oldIndex, newIndex)
		}

		if err != nil {
			return fmt.Errorf("error moving report part: %v", err)
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

		template.Parts[oldIndex].Title = newTitle

		if oldIndex != newIndex {
			err = moveTemplatePart(template, oldIndex, newIndex)
		}

		if err != nil {
			return fmt.Errorf("error moving report part: %v", err)
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

		return nil
	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}
}

func insertReportPart(report *models.Report, part models.ReportPart, index int) error {
	if index < -1 || index > len(report.Parts) {
		// Handle the error or ignore if the index is out of bounds
		return fmt.Errorf("unable to insert part into report. index out of bounds")
	}

	index++
	report.Parts = append(report.Parts[:index], append([]models.ReportPart{part}, report.Parts[index:]...)...)
	return nil
}

func moveReportPart(report *models.Report, fromIndex, toIndex int) error {
	if fromIndex < 0 || fromIndex >= len(report.Parts) || toIndex < -1 || toIndex > len(report.Parts) {
		// Handle the error or ignore if indices are out of bounds
		return fmt.Errorf("unable to move part in report. index out of bounds")
	}

	// Check if fromIndex and toIndex are the same, in which case, do not move
	if fromIndex == toIndex {
		return nil // No action needed as the part is already in the desired position
	}

	// Remove the part from the current position
	part := report.Parts[fromIndex]
	report.Parts = append(report.Parts[:fromIndex], report.Parts[fromIndex+1:]...)

	// If toIndex is the last index, simply append the part to the end
	if toIndex == len(report.Parts) {
		report.Parts = append(report.Parts, part)
		return nil
	}

	// Adjust toIndex if it is greater than fromIndex
	if toIndex > fromIndex {
		toIndex--
	} else {
		// Increment toIndex to insert after the specified index
		toIndex++
	}

	// Reinsert the part at the new position
	report.Parts = append(report.Parts[:toIndex], append([]models.ReportPart{part}, report.Parts[toIndex:]...)...)
	return nil
}

func insertTemplatePart(template *models.Template, part models.TemplatePart, index int) error {
	if index < -1 || index > len(template.Parts) {
		return fmt.Errorf("unable to insert part into template. index out of bounds")
	}
	index++
	template.Parts = append(template.Parts[:index], append([]models.TemplatePart{part}, template.Parts[index:]...)...)
	return nil
}

func moveTemplatePart(template *models.Template, fromIndex, toIndex int) error {
	if fromIndex < 0 || fromIndex >= len(template.Parts) || toIndex < -1 || toIndex > len(template.Parts) {
		// Handle the error or ignore if indices are out of bounds
		return fmt.Errorf("unable to move part in report. index out of bounds")
	}

	// Check if fromIndex and toIndex are the same, in which case, do not move
	if fromIndex == toIndex {
		return nil // No action needed as the part is already in the desired position
	}

	// Remove the part from the current position
	part := template.Parts[fromIndex]
	template.Parts = append(template.Parts[:fromIndex], template.Parts[fromIndex+1:]...)

	// If toIndex is the last index, simply append the part to the end
	if toIndex == len(template.Parts) {
		template.Parts = append(template.Parts, part)
		return nil
	}

	// Adjust toIndex if it is greater than fromIndex
	if toIndex > fromIndex {
		toIndex--
	} else {
		// Increment toIndex to insert after the specified index
		toIndex++
	}

	// Reinsert the part at the new position
	template.Parts = append(template.Parts[:toIndex], append([]models.TemplatePart{part}, template.Parts[toIndex:]...)...)
	return nil
}
