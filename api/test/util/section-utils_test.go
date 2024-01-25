package util_test

import (
	"api/shared/models"
	"api/shared/util"
	"reflect"
	"testing"
)

func TestGenerateSectionStaticText(t *testing.T) {
	section, answers := mockStaticData()
	util.GenerateSectionStaticText(section, answers)

	// Expected results after function execution
	expectedTextOutput := []models.ReportTextOutput{
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

func TestGenerateSectionGeneratorText(t *testing.T) {
	section, answers := mockGeneratorData()

	// Mock function for GeneratePromptResponse
	mockGenerator := MockOpenAiGenerator{}

	err := util.GenerateSectionGeneratorText(mockGenerator, section, answers)
	if err != nil {
		t.Errorf("GenerateSectionGeneratorText returned an error: %v", err)
	}

	// Check the results
	expectedResults := []string{"Blue is a calming color", "Toronto is a vibrant city"}
	for i, textOutput := range section.TextOutputs {
		if textOutput.Result != expectedResults[i] {
			t.Errorf("Expected result %q, got %q", expectedResults[i], textOutput.Result)
		}
	}
}
