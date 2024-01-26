package util

import (
	"api/shared/constants"
	"api/shared/interfaces"
	"api/shared/models"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func UpdateSectionInReport(
	reportID string,
	oldPartIndex uint16,
	newPartIndex uint16,
	oldSectionIndex uint16,
	newSectionIndex uint16,
	newSectionTitle string,
	newQuestions []models.ReportQuestion,
	newTextOutputs []models.ReportTextOutput,
) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	report, err := GetReport(reportID)
	if err != nil {
		return fmt.Errorf("error getting report: %v", err)
	}

	if report == nil {
		return fmt.Errorf("report not found")
	}

	var oldPart, newPart *models.ReportPart
	var sectionToMove *models.ReportSection

	// Find the old part of the section
	for i, part := range report.Parts {
		if part.Index == oldPartIndex {
			oldPart = &report.Parts[i]
			// Get the section to update
			for j, section := range part.Sections {
				if section.Index == oldSectionIndex {
					sectionToMove = &oldPart.Sections[j]
					break
				}
			}
		}
		// Get new part of the section
		if part.Index == newPartIndex {
			newPart = &report.Parts[i]
		}
	}

	// We don't check for if the old part is equal to the new part
	// because it doesn't matter. The logic works if we remove the section from a part
	// and add it back to the same part.
	if newPart == nil {
		return fmt.Errorf("new part not found")
	}
	// Remove the section from the old part
	oldPart.Sections = removeSectionFromReport(oldPart.Sections, oldSectionIndex)
	// Add the section to the new part
	newPart.Sections = addSectionToReport(newPart.Sections, *sectionToMove, newSectionIndex)

	// Update the qualities of the section
	for i, section := range newPart.Sections {
		if section.Index == newSectionIndex {
			newPart.Sections[i].Title = newSectionTitle
			newPart.Sections[i].Questions = newQuestions
			newPart.Sections[i].TextOutputs = newTextOutputs
			newPart.Sections[i].OutputGenerated = false
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

}

func removeSectionFromReport(sections []models.ReportSection, indexToRemove uint16) []models.ReportSection {
	// Create a new slice for the updated sections
	newSections := make([]models.ReportSection, 0, len(sections))

	// First, remove the section with the specified index
	for _, section := range sections {
		if section.Index != indexToRemove {
			newSections = append(newSections, section)
		}
	}

	// Update the index for sections that come after the removed section
	for i := range newSections {
		if newSections[i].Index > indexToRemove {
			newSections[i].Index--
		}
	}

	return newSections
}

func addSectionToReport(sections []models.ReportSection, newSection models.ReportSection, indexToAdd uint16) []models.ReportSection {
	// Create a new slice to hold the updated list of sections
	newSections := make([]models.ReportSection, 0, len(sections)+1)

	// Add the new section first
	newSection.Index = indexToAdd
	newSections = append(newSections, newSection)

	// Iterate over the existing sections and add them to the new slice
	for _, section := range sections {
		// If the current section's index is greater or equal to the indexToAdd, increment its index
		if section.Index >= indexToAdd {
			section.Index++
		}

		// Add the current section to the new sections list
		newSections = append(newSections, section)
	}

	return newSections
}

func UpdateSectionInTemplate(
	templateID string,
	oldPartIndex uint16,
	newPartIndex uint16,
	oldSectionIndex uint16,
	newSectionIndex uint16,
	newSectionTitle string,
	newQuestions []models.TemplateQuestion,
	newTextOutputs []models.TemplateTextOutput,
) error {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	template, err := GetTemplate(templateID)
	if err != nil {
		return fmt.Errorf("error getting template: %v", err)
	}

	if template == nil {
		return fmt.Errorf("template not found")
	}

	var oldPart, newPart *models.TemplatePart
	var sectionToMove *models.TemplateSection

	// Find the old part of the section
	for i, part := range template.Parts {
		if part.Index == oldPartIndex {
			oldPart = &template.Parts[i]
			// Get the section to update
			for j, section := range part.Sections {
				if section.Index == oldSectionIndex {
					sectionToMove = &oldPart.Sections[j]
					break
				}
			}
		}
		// Get new part of the section
		if part.Index == newPartIndex {
			newPart = &template.Parts[i]
		}
	}

	// We don't check for if the old part is equal to the new part
	// because it doesn't matter. The logic works if we remove the section from a part
	// and add it back to the same part.
	if newPart == nil {
		return fmt.Errorf("new part not found")
	}
	// Remove the section from the old part
	oldPart.Sections = removeSectionFromTemplate(oldPart.Sections, oldSectionIndex)
	// Add the section to the new part
	newPart.Sections = addSectionToTemplate(newPart.Sections, *sectionToMove, newSectionIndex)

	// Update the qualities of the section
	for i, section := range newPart.Sections {
		if section.Index == newSectionIndex {
			newPart.Sections[i].Title = newSectionTitle
			newPart.Sections[i].Questions = newQuestions
			newPart.Sections[i].TextOutputs = newTextOutputs
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

}

func removeSectionFromTemplate(sections []models.TemplateSection, indexToRemove uint16) []models.TemplateSection {
	// Create a new slice for the updated sections
	newSections := make([]models.TemplateSection, 0, len(sections))

	// First, remove the section with the specified index
	for _, section := range sections {
		if section.Index != indexToRemove {
			newSections = append(newSections, section)
		}
	}

	// Update the index for sections that come after the removed section
	for i := range newSections {
		if newSections[i].Index > indexToRemove {
			newSections[i].Index--
		}
	}

	return newSections
}

func addSectionToTemplate(sections []models.TemplateSection, newSection models.TemplateSection, indexToAdd uint16) []models.TemplateSection {
	// Create a new slice to hold the updated list of sections
	newSections := make([]models.TemplateSection, 0, len(sections)+1)

	// Add the new section first
	newSection.Index = indexToAdd
	newSections = append(newSections, newSection)

	// Iterate over the existing sections and add them to the new slice
	for _, section := range sections {
		// If the current section's index is greater or equal to the indexToAdd, increment its index
		if section.Index >= indexToAdd {
			section.Index++
		}

		// Add the current section to the new sections list
		newSections = append(newSections, section)
	}

	return newSections
}

func GenerateSection(reportID string, partIndex uint16, sectionIndex uint16, answers []models.Answer) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	report, err := GetReport(reportID)

	if err != nil {
		return fmt.Errorf("error getting report from DynamoDB: %v", err)
	}

	if report == nil {
		return fmt.Errorf("report not found: %v", err)
	}

	section, err := GetSection(report, partIndex, sectionIndex)

	if err != nil {
		return fmt.Errorf("error getting section: %v", err)
	}

	// Reset the text output results so that they can be created from input again
	ResetTextOutputResults(section)

	GenerateSectionStaticText(section, answers)

	generator := OpenAiGenerator{}

	err = GenerateSectionGeneratorText(generator, section, answers)

	if err != nil {
		return fmt.Errorf("error creating generator outputs: %v", err)
	}

	// Set output generated after all sections generated successfully
	section.OutputGenerated = true

	// Update the report in DynamoDB
	updatedReport, err := dynamodbattribute.MarshalMap(report)
	if err != nil {
		return err
	}

	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      updatedReport,
	})
	return err
}

func GenerateSectionStaticText(section *models.ReportSection, answers []models.Answer) {
	// Iterate over each answer
	for _, answer := range answers {
		// Find the matching question
		question, err := FindQuestion(section.Questions, answer.QuestionIndex)
		if err != nil {
			continue // or handle the error as you see fit
		}

		// Update question answer
		question.Answer = answer.Answer

		// Generate static text
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Static {
				// Assuming GenerateStaticText modifies textOutput in place
				GenerateStaticText(&section.TextOutputs[i], question.Label, answer.Answer)
			}
		}
	}
}

func GenerateSectionGeneratorText(generator interfaces.Generator, section *models.ReportSection, answers []models.Answer) error {
	// Preserve original inputs as to be able to re-create the section later with different answers to the questions.
	originalInputs := []string{}
	for _, textOutput := range section.TextOutputs {
		if textOutput.Type == models.Generator {
			originalInputs = append(originalInputs, textOutput.Input)
		}
	}

	// Splice answers into prompts
	for _, answer := range answers {
		// Find the matching question
		question, err := FindQuestion(section.Questions, answer.QuestionIndex)
		if err != nil {
			continue // or handle the error as you see fit
		}

		// Update question answer
		question.Answer = answer.Answer

		// Generate the inputs with question answers spliced in
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Generator {
				GenerateGeneratorInput(&section.TextOutputs[i], question.Label, answer.Answer)
			}
		}
	}

	// Generate the outputs
	for i, textOutput := range section.TextOutputs {
		if textOutput.Type == models.Generator {
			// Assuming GenerateStaticText modifies textOutput in place
			result, err := generator.GeneratePromptResponse(section.TextOutputs[i].Input)

			if err != nil {
				return err
			}

			section.TextOutputs[i].Result = result
		}
	}

	// Restore original inputs
	originalInputIndex := 0
	for i, textOutput := range section.TextOutputs {
		if textOutput.Type == models.Generator {
			section.TextOutputs[i].Input = originalInputs[originalInputIndex]
			originalInputIndex += 1
		}
	}

	return nil
}

// FindQuestion finds a question by its index in a slice of questions
func FindQuestion(questions []models.ReportQuestion, index uint16) (*models.ReportQuestion, error) {
	for i := range questions {
		if questions[i].Index == index {
			return &questions[i], nil
		}
	}
	return nil, errors.New("question not found")
}

// GenerateStaticText processes a TextOutput, splicing in answers into static text outputs.
func GenerateStaticText(textOutput *models.ReportTextOutput, questionLabel, answer string) {
	// Define the pattern to be replaced
	pattern := "**" + questionLabel

	// Replace the pattern with the answer in textOutput.Input
	// If first pass, set it to the input, else set it to the generated output replaced.
	// This way, you can splice question answers in multiple outputs
	if textOutput.Result == "" {
		textOutput.Result = strings.ReplaceAll(textOutput.Input, pattern, answer)
	} else {
		textOutput.Result = strings.ReplaceAll(textOutput.Result, pattern, answer)
	}
}

// This function allows users to define answers in their openai prompts as well
func GenerateGeneratorInput(textOutput *models.ReportTextOutput, questionLabel, answer string) {
	// Define the pattern to be replaced
	pattern := "**" + questionLabel

	// Replace the pattern with the answer in textOutput.Input
	// If first pass, set it to the input, else set it to the generated output replaced.
	// This way, you can splice question answers in multiple inputs
	textOutput.Input = strings.ReplaceAll(textOutput.Input, pattern, answer)
}

// GetSection returns the section from a report based on partIndex and sectionIndex.
func GetSection(report *models.Report, partIndex uint16, sectionIndex uint16) (*models.ReportSection, error) {
	// Search for the part with the matching index
	var foundPart *models.ReportPart
	for i := range report.Parts {
		if report.Parts[i].Index == partIndex {
			foundPart = &report.Parts[i]
			break
		}
	}

	if foundPart == nil {
		return nil, errors.New("part not found")
	}

	// Search for the section with the matching index within the found part
	for i := range foundPart.Sections {
		if foundPart.Sections[i].Index == sectionIndex {
			return &foundPart.Sections[i], nil
		}
	}

	return nil, errors.New("section not found")
}

// ResetTextOutputResults sets all TextOutput.Result fields to an empty string in the provided section.
func ResetTextOutputResults(section *models.ReportSection) {
	if section == nil {
		return // or handle the error as you see fit
	}

	for i := range section.TextOutputs {
		section.TextOutputs[i].Result = ""
	}
}
