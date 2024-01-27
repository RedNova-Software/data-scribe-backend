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

// AddSectionToReport adds a Section to a Part with a specific index in a specified report.
func AddSectionToReport(reportID string, partIndex uint16, sectionIndex uint16, newSection models.ReportSection) error {
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

	err = insertSectionInReport(report, int(partIndex), int(sectionIndex), newSection)

	if err != nil {
		return fmt.Errorf("error inserting report section: %v", err)
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
func AddSectionToTemplate(templateID string, partIndex uint16, sectionIndex uint16, newSection models.TemplateSection) error {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	template, err := GetTemplate(templateID)

	if err != nil {
		return fmt.Errorf("error getting report from DynamoDB: %v", err)
	}

	if template == nil {
		return fmt.Errorf("report not found: %v", err)
	}

	err = insertSectionInTemplate(template, int(partIndex), int(sectionIndex), newSection)

	if err != nil {
		return fmt.Errorf("error inserting report section: %v", err)
	}

	// Marshal the updated Template back to a map
	updatedTemplate, err := dynamodbattribute.MarshalMap(template)
	if err != nil {
		return fmt.Errorf("unable to marshall template into dynamodb attribute: %v", err)
	}

	// Update the Template in the DynamoDB table
	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      updatedTemplate,
	})
	if err != nil {
		return fmt.Errorf("error updating template item in dynamodb: %v", err)
	}

	return nil
}

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

	err = moveSectionInReport(report, int(oldPartIndex), int(oldSectionIndex), int(newPartIndex), int(newSectionIndex))

	updatedSection := &report.Parts[newPartIndex].Sections[newSectionIndex]

	// Update the qualities of the section
	updatedSection.Title = newSectionTitle
	updatedSection.Questions = newQuestions
	updatedSection.TextOutputs = newTextOutputs
	updatedSection.OutputGenerated = false

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

	err = moveSectionInTemplate(template, int(oldPartIndex), int(oldSectionIndex), int(newPartIndex), int(newSectionIndex))

	updatedSection := &template.Parts[newPartIndex].Sections[newSectionIndex]

	// Update the qualities of the section
	updatedSection.Title = newSectionTitle
	updatedSection.Questions = newQuestions
	updatedSection.TextOutputs = newTextOutputs

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

	section, err := GetReportSection(report, partIndex, sectionIndex)

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
		question, err := GetReportQuestion(section.Questions, answer.QuestionIndex)
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
		question, err := GetReportQuestion(section.Questions, answer.QuestionIndex)
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

// GetReportQuestion finds a question by its index in a slice of questions
func GetReportQuestion(questions []models.ReportQuestion, index uint16) (*models.ReportQuestion, error) {
	if int(index) < len(questions) {
		return &questions[index], nil
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

// GetReportSection returns the section from a report based on partIndex and sectionIndex.
func GetReportSection(report *models.Report, partIndex uint16, sectionIndex uint16) (*models.ReportSection, error) {
	if int(partIndex) < len(report.Parts) {
		part := &report.Parts[partIndex]
		if int(sectionIndex) < len(part.Sections) {
			return &part.Sections[sectionIndex], nil
		}
		return nil, errors.New("section not found")
	}
	return nil, errors.New("part not found")
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

func insertSectionInReport(report *models.Report, partIndex int, sectionIndex int, section models.ReportSection) error {
	if partIndex < 0 || partIndex >= len(report.Parts) {
		// Handle out of range partIndex
		return fmt.Errorf("unable to insert section into report. part index out of bounds")
	}

	part := &report.Parts[partIndex]
	if sectionIndex < 0 || sectionIndex > len(part.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to insert section into report. section index out of bounds")
	}

	// The first insert will be inserting into nil
	if part.Sections == nil {
		part.Sections = []models.ReportSection{}
	}

	part.Sections = append(part.Sections[:sectionIndex], append([]models.ReportSection{section}, part.Sections[sectionIndex:]...)...)
	return nil
}

func moveSectionInReport(report *models.Report, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex int) error {
	if oldPartIndex < 0 || oldPartIndex >= len(report.Parts) || newPartIndex < 0 || newPartIndex >= len(report.Parts) {
		// Handle out of range indices
		return fmt.Errorf("unable to move section in report. part index out of bounds")
	}

	// Remove the section from the old part
	oldPart := &report.Parts[oldPartIndex]
	if oldSectionIndex < 0 || oldSectionIndex >= len(oldPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in old part")
	}
	section := oldPart.Sections[oldSectionIndex]
	oldPart.Sections = append(oldPart.Sections[:oldSectionIndex], oldPart.Sections[oldSectionIndex+1:]...)

	// Insert the section into the new part
	newPart := &report.Parts[newPartIndex]
	if newSectionIndex < 0 || newSectionIndex > len(newPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in new part")
	}
	newPart.Sections = append(newPart.Sections[:newSectionIndex], append([]models.ReportSection{section}, newPart.Sections[newSectionIndex:]...)...)
	return nil
}

func insertSectionInTemplate(report *models.Template, partIndex int, sectionIndex int, section models.TemplateSection) error {
	if partIndex < 0 || partIndex >= len(report.Parts) {
		// Handle out of range partIndex
		return fmt.Errorf("unable to insert section into report. part index out of bounds")
	}

	part := &report.Parts[partIndex]
	if sectionIndex < 0 || sectionIndex > len(part.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to insert section into report. section index out of bounds")
	}

	// The first insert will be inserting into nil
	if part.Sections == nil {
		part.Sections = []models.TemplateSection{}
	}

	part.Sections = append(part.Sections[:sectionIndex], append([]models.TemplateSection{section}, part.Sections[sectionIndex:]...)...)
	return nil
}

func moveSectionInTemplate(report *models.Template, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex int) error {
	if oldPartIndex < 0 || oldPartIndex >= len(report.Parts) || newPartIndex < 0 || newPartIndex >= len(report.Parts) {
		// Handle out of range indices
		return fmt.Errorf("unable to move section in report. part index out of bounds")
	}

	// Remove the section from the old part
	oldPart := &report.Parts[oldPartIndex]
	if oldSectionIndex < 0 || oldSectionIndex >= len(oldPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in old part")
	}
	section := oldPart.Sections[oldSectionIndex]
	oldPart.Sections = append(oldPart.Sections[:oldSectionIndex], oldPart.Sections[oldSectionIndex+1:]...)

	// Insert the section into the new part
	newPart := &report.Parts[newPartIndex]
	if newSectionIndex < 0 || newSectionIndex > len(newPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in new part")
	}
	newPart.Sections = append(newPart.Sections[:newSectionIndex], append([]models.TemplateSection{section}, newPart.Sections[newSectionIndex:]...)...)
	return nil
}
