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
	report.LastModifiedAt = GetCurrentTime()

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
	template.LastModifiedAt = GetCurrentTime()

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
		report.LastModifiedAt = GetCurrentTime()

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
		template.LastModifiedAt = GetCurrentTime()

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
	newCSVData []models.ReportCSVData,
	newChartOutputs []models.ReportChartOutput,
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

	updateReportTextOutputs(updatedSection, newTextOutputs, deleteGeneratedOutput)

	// Update csv data and chart outputs
	updatedSection.CSVData = newCSVData
	updatedSection.ChartOutputs = newChartOutputs

	if oldPartIndex != newPartIndex || oldSectionIndex != newSectionIndex {
		err = moveSectionInReport(report, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex)
		if err != nil {
			return err
		}
	}

	// Update last modified
	report.LastModifiedAt = GetCurrentTime()

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
	newCSVData []models.TemplateCSVData,
	newChartOutputs []models.TemplateChartOutput,
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
	updatedSection.CSVData = newCSVData
	updatedSection.ChartOutputs = newChartOutputs

	if oldPartIndex != newPartIndex || oldSectionIndex != newSectionIndex {
		err = moveSectionInTemplate(template, oldPartIndex, oldSectionIndex, newPartIndex, newSectionIndex)
		if err != nil {
			return err
		}
	}

	// Update last modified
	template.LastModifiedAt = GetCurrentTime()

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

func SetReportSectionResponses(reportID string,
	partIndex int,
	sectionIndex int,
	questionAnswers []models.Answer,
	csvDataResponses []models.CsvDataResponse,
	chartOutputResponses []models.ChartOutputResponse,
	userID string) error {

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

	// First, update question answers
	if section.Questions != nil {
		for i := range section.Questions {
			section.Questions[i].Answer = questionAnswers[i].Answer
		}
	}

	// Next, update csv data responses
	if section.CSVData != nil {
		for i := range section.CSVData {
			section.CSVData[i].OperationColumn = csvDataResponses[i].OperationColumn
			section.CSVData[i].AcceptedValues = csvDataResponses[i].AcceptedValues
			section.CSVData[i].FilterColumns = csvDataResponses[i].FilterColumns
		}
	}

	// Next, update chart output responses
	if section.ChartOutputs != nil {
		for i := range section.ChartOutputs {
			section.ChartOutputs[i].IndependentColumn = chartOutputResponses[i].IndependentColumn
			section.ChartOutputs[i].AcceptedValues = chartOutputResponses[i].AcceptedValues
			section.ChartOutputs[i].FilterColumns = chartOutputResponses[i].FilterColumns

			// Update the dependent columns
			for j := range section.ChartOutputs[i].DependentColumns {
				section.ChartOutputs[i].DependentColumns[j].Column = chartOutputResponses[i].DependentColumns[j].Column
				section.ChartOutputs[i].DependentColumns[j].AcceptedValues = chartOutputResponses[i].DependentColumns[j].AcceptedValues
				section.ChartOutputs[i].DependentColumns[j].FilterColumns = chartOutputResponses[i].DependentColumns[j].FilterColumns
			}
		}
	}

	// Update last modified
	report.LastModifiedAt = GetCurrentTime()

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

func GenerateSection(reportID string, partIndex int, sectionIndex int, generateAIOutput bool, userID string) error {
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

	// Load CSV file from S3
	csvFile, err := GetCSVFileHandle(report.CSVID)
	if err != nil {
		return fmt.Errorf("error loading CSV from S3: %v", err)
	}

	// Generate csv data results from csv
	err = GenerateSectionCsvDataResults(csvFile, section)

	if err != nil {
		return fmt.Errorf("error generating section csv data results: %v", err)
	}

	// Reset the text output results so that they can be created from input again
	ResetTextOutputResults(section, generateAIOutput)

	GenerateSectionStaticText(section, &report.GlobalQuestions)

	if generateAIOutput {
		generator := OpenAiGenerator{}

		err = GenerateSectionGeneratorText(generator, section, &report.GlobalQuestions)
		if err != nil {
			log.Panicf("error creating generator outputs: %v", err)
			return fmt.Errorf("error creating generator outputs: %v", err)
		}
	}

	// Generate Chart Output from CSV
	err = GenerateChartOutputResults(csvFile, section)

	if err != nil {
		return fmt.Errorf("error generating section csv data results: %v", err)
	}

	// Set output generated after all sections generated successfully
	section.OutputGenerated = true

	// Update last modified
	report.LastModifiedAt = GetCurrentTime()

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

func GenerateSectionCsvDataResults(csvFile *os.File, section *models.ReportSection) error {
	// Iterate over each csv data
	for index := range section.CSVData {

		err := AnalyzeOneDimensionalData(csvFile, &section.CSVData[index])
		if err != nil {
			return fmt.Errorf("error generating section csv data results: %v", err)
		}
	}
	return nil
}

func GenerateChartOutputResults(csvFile *os.File, section *models.ReportSection) error {
	// Iterate over each chart output
	for index := range section.ChartOutputs {

		err := AnalyzeTwoDimensionalData(csvFile, &section.ChartOutputs[index])
		if err != nil {
			return fmt.Errorf("error generating section chart output results: %v", err)
		}
	}
	return nil
}

func GenerateSectionStaticText(section *models.ReportSection, globalQuestions *[]models.ReportQuestion) {
	// Iterate over each question, splicing answer into text
	for _, question := range section.Questions {

		// Generate static text
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Static {
				GenerateStaticText(&section.TextOutputs[i], question.Label, question.Answer)
			}
		}
	}

	// Iterate over each csv data, splicing result into text
	for _, csvData := range section.CSVData {

		// Generate static text
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Static {
				GenerateStaticText(&section.TextOutputs[i], csvData.Label, csvData.Result)
			}
		}
	}

	// Iterate over each global question, splicing answer into text
	for _, question := range *globalQuestions {

		// Generate static text
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Static {
				GenerateStaticText(&section.TextOutputs[i], question.Label, question.Answer)
			}
		}
	}
}

func GenerateSectionGeneratorText(generator interfaces.Generator, section *models.ReportSection, globalQuestions *[]models.ReportQuestion) error {
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

	// Splice question answers into prompts
	for _, question := range section.Questions {

		// Generate the inputs with question answers spliced in
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Generator {
				GenerateGeneratorInput(&section.TextOutputs[i], question.Label, question.Answer)
			}
		}
	}

	// Splice the csv data results into prompts
	for _, csvData := range section.CSVData {

		// Generate the inputs with question answers spliced in
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Generator {
				GenerateGeneratorInput(&section.TextOutputs[i], csvData.Label, csvData.Result)
			}
		}
	}

	// Splice global question answers into prompts
	for _, question := range *globalQuestions {

		// Generate the inputs with question answers spliced in
		for i, textOutput := range section.TextOutputs {
			if textOutput.Type == models.Generator {
				GenerateGeneratorInput(&section.TextOutputs[i], question.Label, question.Answer)
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
				log.Printf("input after splicing: %v", input)

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
func GenerateStaticText(textOutput *models.ReportTextOutput, label, value string) {
	// Define the pattern to be replaced
	pattern := "**" + label

	// Replace the pattern with the answer in textOutput.Input
	// If first pass, set it to the input, else set it to the generated output replaced.
	// This way, you can splice question answers in multiple outputs
	if textOutput.Result == "" {
		textOutput.Result = strings.ReplaceAll(textOutput.Input, pattern, value)
	} else {
		textOutput.Result = strings.ReplaceAll(textOutput.Result, pattern, value)
	}
}

// This function allows users to define answers in their openai prompts as well
func GenerateGeneratorInput(textOutput *models.ReportTextOutput, label, value string) {
	// Define the pattern to be replaced
	pattern := "**" + label

	// Replace the pattern with the answer in textOutput.Input
	// If first pass, set it to the input, else set it to the generated output replaced.
	// This way, you can splice question answers in multiple inputs
	textOutput.Input = strings.ReplaceAll(textOutput.Input, pattern, value)
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

func updateReportTextOutputs(section *models.ReportSection, newTextOutputs []models.ReportTextOutput, clearGeneratorResult bool) {
	for _, newTextOutput := range newTextOutputs {
		found := false
		for i, existingOutput := range section.TextOutputs {
			// Check if the ReportTextOutput already exists (by Title and Type)
			if existingOutput.Title == newTextOutput.Title && existingOutput.Type == newTextOutput.Type {
				found = true
				// Update existing ReportTextOutput
				if clearGeneratorResult && newTextOutput.Type == models.Generator {
					section.TextOutputs[i].Result = "" // Clear Result if specified and type is Generator
				} else {
					section.TextOutputs[i].Result = newTextOutput.Result
				}
				section.TextOutputs[i].Input = newTextOutput.Input
				break
			}
		}
		// If the ReportTextOutput is not found, add it to the section
		if !found {
			if clearGeneratorResult && newTextOutput.Type == models.Generator {
				newTextOutput.Result = "" // Clear Result if specified and type is Generator
			}
			section.TextOutputs = append(section.TextOutputs, newTextOutput)
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
