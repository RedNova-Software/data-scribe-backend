package main

import (
	"api/shared/models"
	"api/shared/util"
	"reflect"
	"testing"
)

// Mock data for testing
func mockData() (*models.Section, []models.Answer) {
	section := &models.Section{
		Title: "Test Section",
		Questions: []models.Question{
			{
				Label:    "question1",
				Index:    0,
				Question: "What is 2+2?",
				Answer:   "",
			},
			{
				Label:    "question2",
				Index:    1,
				Question: "What is the capital of France?",
				Answer:   "",
			},
		},
		TextOutputs: []models.TextOutput{
			{
				Title:  "Output1",
				Index:  0,
				Type:   models.Static,
				Input:  "The answer to 2 + 2 = **question1",
				Result: "",
			},
			{
				Title:  "Output2",
				Index:  1,
				Type:   models.Static,
				Input:  "**question2 is the capital of France.",
				Result: "",
			},
		},
	}

	answers := []models.Answer{
		{
			QuestionIndex: 0,
			Answer:        "4",
		},
		{
			QuestionIndex: 1,
			Answer:        "Paris",
		},
	}

	return section, answers
}

func TestGenerateSectionStaticText(t *testing.T) {
	section, answers := mockData()
	util.GenerateSectionStaticText(section, answers)

	// Expected results after function execution
	expectedTextOutput := []models.TextOutput{
		{
			Title:  "Output1",
			Index:  0,
			Type:   models.Static,
			Input:  "The answer to 2 + 2 = **question1",
			Result: "The answer to 2 + 2 = 4",
		},
		{
			Title:  "Output2",
			Index:  1,
			Type:   models.Static,
			Input:  "**question2 is the capital of France.",
			Result: "Paris is the capital of France.",
		},
	}

	// Check if the TextOutputs are as expected
	if !reflect.DeepEqual(section.TextOutputs, expectedTextOutput) {
		t.Errorf("TextOutputs were not updated correctly. Got: \n %v \n, want: \n %v", section.TextOutputs, expectedTextOutput)
	}
}
