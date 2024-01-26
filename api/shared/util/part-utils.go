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
	partTitle string,
	partIndex uint16,
) error {
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	var tableName string
	var itemKey string
	var newPartAV map[string]*dynamodb.AttributeValue

	if itemType == constants.Report {
		tableName = os.Getenv(constants.ReportTable)
		itemKey = constants.ReportIDField

		newPart := models.ReportPart{
			Title:    partTitle,
			Index:    partIndex,
			Sections: []models.ReportSection{},
		}
		newPartAV, err = dynamodbattribute.MarshalMap(newPart)

		if err != nil {
			return err
		}

	} else if itemType == constants.Template {
		tableName = os.Getenv(constants.TemplateTable)
		itemKey = constants.TemplateIDField

		newPart := models.TemplatePart{
			Title:    partTitle,
			Index:    partIndex,
			Sections: []models.TemplateSection{},
		}
		newPartAV, err = dynamodbattribute.MarshalMap(newPart)

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("incorrect item type specified. must be either 'report' or 'template'")
	}

	newPartAV[constants.SectionsField] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}

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

func UpdatePartInItem(
	itemType constants.ItemType,
	itemID string,
	oldIndex uint16,
	newTitle string,
	newIndex uint16,
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

		err = MoveReportPartIndex(report, oldIndex, newIndex)

		if err != nil {
			return fmt.Errorf("error moving report part indices: %v", err)
		}

		for i, part := range report.Parts {
			if part.Index == newIndex {
				report.Parts[i].Title = newTitle
				break
			}
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

		err = MoveTemplatePartIndex(template, oldIndex, newIndex)

		if err != nil {
			return fmt.Errorf("error moving template part indices: %v", err)
		}

		for i, part := range template.Parts {
			if part.Index == newIndex {
				template.Parts[i].Title = newTitle
				break
			}
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

func MoveReportPartIndex(report *models.Report, oldIndex, newIndex uint16) error {
	if oldIndex == newIndex {
		return nil // No change if the indices are the same
	}

	partCount := uint16(len(report.Parts))
	if oldIndex >= partCount || newIndex >= partCount {
		return fmt.Errorf("indices out of range")
	}

	// Find and update the index of the moving part
	movingPartFound := false
	for i := range report.Parts {
		if report.Parts[i].Index == oldIndex {
			report.Parts[i].Index = newIndex
			movingPartFound = true
			break
		}
	}

	if !movingPartFound {
		return fmt.Errorf("part with old index %d not found", oldIndex)
	}

	// Update the indices of other parts affected by the move
	if oldIndex < newIndex {
		// Moving towards a higher index, decrement indices of parts in between
		for i := range report.Parts {
			if report.Parts[i].Index > oldIndex && report.Parts[i].Index <= newIndex {
				report.Parts[i].Index--
			}
		}
	} else {
		// Moving towards a lower index, increment indices of parts in between
		for i := range report.Parts {
			if report.Parts[i].Index < oldIndex && report.Parts[i].Index >= newIndex {
				report.Parts[i].Index++
			}
		}
	}

	return nil
}

func MoveTemplatePartIndex(template *models.Template, oldIndex, newIndex uint16) error {
	if oldIndex == newIndex {
		return nil // No change if the indices are the same
	}

	partCount := uint16(len(template.Parts))
	if oldIndex >= partCount || newIndex >= partCount {
		return fmt.Errorf("indices out of range")
	}

	// Find and update the index of the moving part
	movingPartFound := false
	for i := range template.Parts {
		if template.Parts[i].Index == oldIndex {
			template.Parts[i].Index = newIndex
			movingPartFound = true
			break
		}
	}

	if !movingPartFound {
		return fmt.Errorf("part with old index %d not found", oldIndex)
	}

	// Update the indices of other parts affected by the move
	if oldIndex < newIndex {
		// Moving towards a higher index, decrement indices of parts in between
		for i := range template.Parts {
			if template.Parts[i].Index > oldIndex && template.Parts[i].Index <= newIndex {
				template.Parts[i].Index--
			}
		}
	} else {
		// Moving towards a lower index, increment indices of parts in between
		for i := range template.Parts {
			if template.Parts[i].Index < oldIndex && template.Parts[i].Index >= newIndex {
				template.Parts[i].Index++
			}
		}
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
