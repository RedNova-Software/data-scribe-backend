package util

import (
	"api/shared/constants"
	"api/shared/interfaces"
	"api/shared/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// AddSectionToReport adds a Section to a Part with a specific index in a specified report.
func AddSectionToReport(reportID string, partIndex int, sectionIndex int, newSection models.ReportSection, userID string) error {
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

	err = insertSectionInReport(report, partIndex, sectionIndex, newSection)

	if err != nil {
		return fmt.Errorf("error inserting report section: %v", err)
	}

	// Update last modified
	report.LastModified = GetCurrentTime()

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
func AddSectionToTemplate(templateID string, partIndex int, sectionIndex int, newSection models.TemplateSection, userID string) error {
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

	err = insertSectionInTemplate(template, partIndex, sectionIndex, newSection)

	if err != nil {
		return fmt.Errorf("error inserting report section: %v", err)
	}

	// Update last modified
	template.LastModified = GetCurrentTime()

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

// AddSectionToReport adds a Section to a Part with a specific index in a specified report.
func DeleteSectionFromItem(itemType constants.ItemType,
	itemID string,
	partIndex int,
	sectionIndex int,
	userID string) error {

	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	var tableName string

	if itemType == constants.Report {
		tableName := os.Getenv(constants.ReportTable)
		report, err := GetReport(itemID, userID)

		if err != nil {
			return fmt.Errorf("error getting report from DynamoDB: %v", err)
		}

		if report == nil {
			return fmt.Errorf("report not found: %v", err)
		}

		err = deleteReportSection(report, partIndex, sectionIndex)

		if err != nil {
			return fmt.Errorf("error deleteing report part: %v", err)
		}

		// Update last modified
		report.LastModified = GetCurrentTime()

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
		template, err := GetTemplate(itemID, userID)

		if err != nil {
			return fmt.Errorf("error getting template from DynamoDB: %v", err)
		}

		if template == nil {
			return fmt.Errorf("template not found: %v", err)
		}

		err = deleteTemplateSection(template, partIndex, sectionIndex)

		if err != nil {
			return fmt.Errorf("error deleteing report part: %v", err)
		}

		// Update last modified
		template.LastModified = GetCurrentTime()

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

func UpdateSectionInReport(
	reportID string,
	oldPartIndex int,
	newPartIndex int,
	oldSectionIndex int,
	newSectionIndex int,
	newSectionTitle string,
	newQuestions []models.ReportQuestion,
	newTextOutputs []models.ReportTextOutput,
	deleteGeneratedOutput bool,
	userID string,
) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	report, err := GetReport(reportID, userID)
	if err != nil {
		return fmt.Errorf("error getting report: %v", err)
	}

	if report == nil {
		return fmt.Errorf("report not found")
	}

	updatedSection := &report.Parts[oldPartIndex].Sections[oldSectionIndex]

	// Update the title and questions of the section
	updatedSection.Title = newSectionTitle
	updatedSection.Questions = newQuestions

	// If deleteGeneratedOutput is true or the type is not Generator, update the TextOutputs as is
	if deleteGeneratedOutput {
		updatedSection.TextOutputs = newTextOutputs
		// Rest output generated since we're wiping all outputs
		updatedSection.OutputGenerated = false
	} else {
		// Otherwise, update selectively
		for i, newTextOutput := range newTextOutputs {
			if newTextOutput.Type == models.Generator {
				// Check if the corresponding text output already exists
				if i < len(updatedSection.TextOutputs) && updatedSection.TextOutputs[i].Type == models.Generator {
					// Keep the existing Result
					newTextOutput.Result = updatedSection.TextOutputs[i].Result
				}
			}
			updatedSection.TextOutputs[i] = newTextOutput
		}
	}

	if oldPartIndex != newPartIndex || oldSectionIndex != newSectionIndex {
		err = moveSectionInReport(report, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex)
		if err != nil {
			return err
		}
	}

	// Update last modified
	report.LastModified = GetCurrentTime()

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
	oldPartIndex int,
	newPartIndex int,
	oldSectionIndex int,
	newSectionIndex int,
	newSectionTitle string,
	newQuestions []models.TemplateQuestion,
	newTextOutputs []models.TemplateTextOutput,
	userID string,
) error {
	tableName := os.Getenv(constants.TemplateTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	template, err := GetTemplate(templateID, userID)
	if err != nil {
		return fmt.Errorf("error getting template: %v", err)
	}

	if template == nil {
		return fmt.Errorf("template not found")
	}

	updatedSection := &template.Parts[oldPartIndex].Sections[oldSectionIndex]

	// Update the qualities of the section
	updatedSection.Title = newSectionTitle
	updatedSection.Questions = newQuestions
	updatedSection.TextOutputs = newTextOutputs

	if oldPartIndex != newPartIndex || oldSectionIndex != newSectionIndex {
		err = moveSectionInTemplate(template, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex)
		if err != nil {
			return err
		}
	}

	// Update last modified
	template.LastModified = GetCurrentTime()

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

func GenerateSection(reportID string, partIndex int, sectionIndex int, answers []models.Answer, generateAIOutput bool, userID string) error {
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

	section, err := GetReportSection(report, partIndex, sectionIndex)

	if err != nil {
		return fmt.Errorf("error getting section: %v", err)
	}

	// Reset the text output results so that they can be created from input again
	ResetTextOutputResults(section, generateAIOutput)

	GenerateSectionStaticText(section, answers)

	if generateAIOutput {
		generator := OpenAiGenerator{}

		err = GenerateSectionGeneratorText(generator, section, answers)
		if err != nil {
			log.Panicf("error creating generator outputs: %v", err)
			return fmt.Errorf("error creating generator outputs: %v", err)
		}
	}

	// Set output generated after all sections generated successfully
	section.OutputGenerated = true

	// Update last modified
	report.LastModified = GetCurrentTime()

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

	// Used to make the concurrent channel
	numberOfGeneratorSections := 0
	for _, textOutput := range section.TextOutputs {
		if textOutput.Type == models.Generator {
			numberOfGeneratorSections++
		}
	}

	// Splice answers into prompts
	for _, answer := range answers {
		// Find the matching question
		question, err := GetReportQuestion(section.Questions, answer.QuestionIndex)
		if err != nil {
			continue
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

	type generateResult struct {
		Index  int
		Result string
		Err    error
	}

	// Create a channel for communication
	resultsChan := make(chan generateResult, numberOfGeneratorSections)

	log.Print("Starting GPT Generation\n")

	for i, textOutput := range section.TextOutputs {
		if textOutput.Type == models.Generator {
			go func(index int, input string) {
				log.Print("Generating TextOutput: " + strconv.Itoa(index) + "\n")
				result, err := generator.GeneratePromptResponse(input)
				// Send a Result struct to the channel
				resultsChan <- generateResult{Index: index, Result: result, Err: err}
			}(i, textOutput.Input)
		}
	}

	// Process the results
	for i := 0; i < numberOfGeneratorSections; i++ {
		result := <-resultsChan
		log.Print("Processing Result: " + strconv.Itoa(result.Index) + "\n")
		log.Print("Result: " + result.Result + "\n")

		if result.Err != nil {
			log.Print("Err: " + string(result.Err.Error()) + "\n")
			section.TextOutputs[result.Index].Result = "err generating: " + result.Err.Error()
		} else {
			section.TextOutputs[result.Index].Result = result.Result
		}
	}

	// Remember to close the channel
	close(resultsChan)

	log.Print("Generation Finished")

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
func GetReportQuestion(questions []models.ReportQuestion, index int) (*models.ReportQuestion, error) {
	if index < len(questions) {
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
func GetReportSection(report *models.Report, partIndex int, sectionIndex int) (*models.ReportSection, error) {
	if partIndex < len(report.Parts) {
		part := &report.Parts[partIndex]
		if sectionIndex < len(part.Sections) {
			return &part.Sections[sectionIndex], nil
		}
		return nil, errors.New("section not found")
	}
	return nil, errors.New("part not found")
}

// ResetTextOutputResults sets all TextOutput.Result fields to an empty string in the provided section.
func ResetTextOutputResults(section *models.ReportSection, generateAIOutput bool) {
	if section == nil {
		return // or handle the error as you see fit
	}

	for i := range section.TextOutputs {
		if generateAIOutput {
			section.TextOutputs[i].Result = ""
		} else {
			if section.TextOutputs[i].Type == models.Static {
				section.TextOutputs[i].Result = ""
			}
		}

	}
}

func insertSectionInReport(report *models.Report, partIndex int, sectionIndex int, section models.ReportSection) error {
	if partIndex < 0 || partIndex >= len(report.Parts) {
		// Handle out of range partIndex
		return fmt.Errorf("unable to insert section into template. part index out of bounds")
	}

	part := &report.Parts[partIndex]
	if sectionIndex < -1 || sectionIndex > len(part.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to insert section into template. section index out of bounds")
	}

	// The first insert will be inserting into nil
	if part.Sections == nil {
		part.Sections = []models.ReportSection{}
	}

	// Special case handling for inserting at the beginning
	if sectionIndex == -1 {
		part.Sections = append([]models.ReportSection{section}, part.Sections...)
		return nil
	}

	// Adjust the index to insert after the specified sectionIndex
	if sectionIndex != -1 {
		sectionIndex++
	}

	// Insert the section
	part.Sections = append(part.Sections[:sectionIndex], append([]models.ReportSection{section}, part.Sections[sectionIndex:]...)...)
	return nil
}

func deleteReportSection(report *models.Report, partIndex int, sectionIndex int) error {
	// Check if partIndex is within the range of the Parts slice
	if partIndex < 0 || partIndex >= len(report.Parts) {
		return errors.New("partIndex is out of range")
	}

	// Get the part from the report
	part := report.Parts[partIndex]

	// Check if sectionIndex is within the range of the Sections slice in the part
	if sectionIndex < 0 || sectionIndex >= len(part.Sections) {
		return errors.New("sectionIndex is out of range")
	}

	// Remove the section at the specified index
	part.Sections = append(part.Sections[:sectionIndex], part.Sections[sectionIndex+1:]...)

	// Update the part in the report
	report.Parts[partIndex] = part

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

	// Handle the special case where newSectionIndex is -1
	if newSectionIndex == -1 {
		newSectionIndex = 0
	} else {
		// Increment to insert after the specified index
		newSectionIndex++
	}

	// Adjust newSectionIndex if moving within the same part and the target position comes after the removed section
	if oldPartIndex == newPartIndex && newSectionIndex > oldSectionIndex {
		newSectionIndex--
	}

	// Handle removal after adjusting to avoid index out of range issues
	oldPart.Sections = append(oldPart.Sections[:oldSectionIndex], oldPart.Sections[oldSectionIndex+1:]...)

	// Insert the section into the new part
	newPart := &report.Parts[newPartIndex]
	if newSectionIndex < 0 || newSectionIndex > len(newPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in new part")
	}

	// Insert the section
	newPart.Sections = append(newPart.Sections[:newSectionIndex], append([]models.ReportSection{section}, newPart.Sections[newSectionIndex:]...)...)
	return nil
}

func deleteTemplateSection(template *models.Template, partIndex int, sectionIndex int) error {
	// Check if partIndex is within the range of the Parts slice
	if partIndex < 0 || partIndex >= len(template.Parts) {
		return errors.New("partIndex is out of range")
	}

	// Get the part from the report
	part := template.Parts[partIndex]

	// Check if sectionIndex is within the range of the Sections slice in the part
	if sectionIndex < 0 || sectionIndex >= len(part.Sections) {
		return errors.New("sectionIndex is out of range")
	}

	// Remove the section at the specified index
	part.Sections = append(part.Sections[:sectionIndex], part.Sections[sectionIndex+1:]...)

	// Update the part in the report
	template.Parts[partIndex] = part

	return nil
}

func insertSectionInTemplate(template *models.Template, partIndex int, sectionIndex int, section models.TemplateSection) error {
	if partIndex < 0 || partIndex >= len(template.Parts) {
		// Handle out of range partIndex
		return fmt.Errorf("unable to insert section into template. part index out of bounds")
	}

	part := &template.Parts[partIndex]
	if sectionIndex < -1 || sectionIndex > len(part.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to insert section into template. section index out of bounds")
	}

	// The first insert will be inserting into nil
	if part.Sections == nil {
		part.Sections = []models.TemplateSection{}
	}

	// Special case handling for inserting at the beginning
	if sectionIndex == -1 {
		part.Sections = append([]models.TemplateSection{section}, part.Sections...)
		return nil
	}

	// Adjust the index to insert after the specified sectionIndex
	if sectionIndex != -1 {
		sectionIndex++
	}

	// Insert the section
	part.Sections = append(part.Sections[:sectionIndex], append([]models.TemplateSection{section}, part.Sections[sectionIndex:]...)...)
	return nil
}

func moveSectionInTemplate(template *models.Template, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex int) error {
	if oldPartIndex < 0 || oldPartIndex >= len(template.Parts) || newPartIndex < 0 || newPartIndex >= len(template.Parts) {
		// Handle out of range indices
		return fmt.Errorf("unable to move section in report. part index out of bounds")
	}

	// Remove the section from the old part
	oldPart := &template.Parts[oldPartIndex]
	if oldSectionIndex < 0 || oldSectionIndex >= len(oldPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in old part")
	}
	section := oldPart.Sections[oldSectionIndex]

	// Handle the special case where newSectionIndex is -1
	if newSectionIndex == -1 {
		newSectionIndex = 0
	} else {
		// Increment to insert after the specified index
		newSectionIndex++
	}

	// Adjust newSectionIndex if moving within the same part and the target position comes after the removed section
	if oldPartIndex == newPartIndex && newSectionIndex > oldSectionIndex {
		newSectionIndex--
	}

	// Handle removal after adjusting to avoid index out of range issues
	oldPart.Sections = append(oldPart.Sections[:oldSectionIndex], oldPart.Sections[oldSectionIndex+1:]...)

	// Insert the section into the new part
	newPart := &template.Parts[newPartIndex]
	if newSectionIndex < 0 || newSectionIndex > len(newPart.Sections) {
		// Handle out of range sectionIndex
		return fmt.Errorf("unable to move section in report. section index out of bounds in new part")
	}

	// Insert the section
	newPart.Sections = append(newPart.Sections[:newSectionIndex], append([]models.TemplateSection{section}, newPart.Sections[newSectionIndex:]...)...)
	return nil
}
