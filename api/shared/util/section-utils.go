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

func GenerateSection(reportID string, partIndex uint16, sectionIndex uint16, answers []models.Answer) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := newDynamoDBClient(constants.USEast2)

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
