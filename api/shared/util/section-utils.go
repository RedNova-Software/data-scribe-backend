package util

import (
	"api/shared/models"
	"errors"
	"strings"
)

func GenerateSectionStaticText(section *models.Section, answers []models.Answer) {
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

// FindQuestion finds a question by its index in a slice of questions
func FindQuestion(questions []models.Question, index uint16) (*models.Question, error) {
	for i := range questions {
		if questions[i].Index == index {
			return &questions[i], nil
		}
	}
	return nil, errors.New("question not found")
}

// GenerateStaticText processes a TextOutput, splicing in answers into static text outputs.
func GenerateStaticText(textOutput *models.TextOutput, questionLabel, answer string) {
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

// GetSection returns the section from a report based on partIndex and sectionIndex.
func GetSection(report *models.Report, partIndex uint16, sectionIndex uint16) (*models.Section, error) {
	// Search for the part with the matching index
	var foundPart *models.Part
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
