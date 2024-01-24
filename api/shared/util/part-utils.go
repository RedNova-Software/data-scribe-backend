package util

import (
	"api/shared/constants"
	"api/shared/models"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func AddPartToItem(
	itemType constants.ItemType,
	itemID string,
	newPart models.ReportPart,
) error {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	newPartAV, err := dynamodbattribute.MarshalMap(newPart)
	if err != nil {
		return err
	}

	newPartAV[constants.SectionsField] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

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

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			itemKey: {
				S: aws.String(itemID),
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
func ModifyItemPartIndices(itemType constants.ItemType, itemID string, newIndex uint16, increment bool) (bool, error) {
	var tableName string
	updated := false

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return false, fmt.Errorf("error getting dynamodb client: %v", err)
	}

	if itemType == constants.Report {
		tableName = os.Getenv(constants.ReportTable)
		report, err := GetReport(itemID)

		if err != nil {
			return false, fmt.Errorf("error getting report from DynamoDB: %v", err)
		}

		if report == nil {
			return false, fmt.Errorf("report not found: %v", err)
		}

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

	} else if itemType == constants.Template {
		tableName = os.Getenv(constants.TemplateTable)
		template, err := GetTemplate(itemID)

		if err != nil {
			return false, fmt.Errorf("error getting template from DynamoDB: %v", err)
		}

		if template == nil {
			return false, fmt.Errorf("template not found: %v", err)
		}

		for i, part := range template.Parts {
			if increment && part.Index >= newIndex {
				template.Parts[i].Index++
				updated = true
			} else if !increment && part.Index > newIndex {
				template.Parts[i].Index--
				updated = true
			}
		}

		if updated {
			av, err := dynamodbattribute.MarshalMap(template)
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
	} else {
		return false, fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	return false, nil
}

// If increment is true, all Part indices equal to or above the newIndex will be incremented.
// If false, everything larger will be decremented.
func ModifyPartSectionIndices(itemType constants.ItemType, itemID string, partIndex uint16, newSectionIndex uint16, increment bool) (bool, error) {
	var tableName string
	updated := false

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return false, fmt.Errorf("error getting dynamodb client: %v", err)
	}

	if itemType == constants.Report {
		tableName = os.Getenv(constants.ReportTable)
		report, err := GetReport(itemID)

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

	} else if itemType == constants.Template {
		tableName = os.Getenv(constants.TemplateTable)
		template, err := GetTemplate(itemID)

		if err != nil {
			return false, fmt.Errorf("error getting report from DynamoDB: %v", err)
		}

		if template == nil {
			return false, fmt.Errorf("report not found: %v", err)
		}

		var part *models.TemplatePart

		for _, searchPart := range template.Parts {
			if searchPart.Index == partIndex {
				part = &searchPart
				break
			}
		}

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
			av, err := dynamodbattribute.MarshalMap(template)
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

	} else {
		return false, fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	return false, nil
}

// AddSectionToReportPart adds a Section to a Part with a specific index in a specified report.
func AddSectionToReportPart(reportID string, partIndex uint16, newSection models.ReportSection) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
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

// AddSectionToReportPart adds a Section to a Part with a specific index in a specified template.
func AddSectionToTemplatePart(templateID string, partIndex uint16, newSection models.TemplateSection) error {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	// Retrieve the Report item
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.TemplateIDField: {
				S: aws.String(templateID),
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
	var template models.Template
	err = dynamodbattribute.UnmarshalMap(result.Item, &template)
	if err != nil {
		return fmt.Errorf("unable to unmarshal dynamodb template: %v", err)
	}

	// Find the Part and add the Section
	partFound := false
	for i, part := range template.Parts {
		if part.Index == partIndex {
			template.Parts[i].Sections = append(template.Parts[i].Sections, newSection)
			partFound = true
			break
		}
	}

	if !partFound {
		return errors.New("part not found")
	}

	// Marshal the updated Report back to a map
	updatedTemplate, err := dynamodbattribute.MarshalMap(template)
	if err != nil {
		return fmt.Errorf("unable to marshall template into dynamodb attribute: %v", err)
	}

	// Update the Report in the DynamoDB table
	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      updatedTemplate,
	})
	if err != nil {
		return fmt.Errorf("error updating template item in dynamodb: %v", err)
	}

	return nil
}
