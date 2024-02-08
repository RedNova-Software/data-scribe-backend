package util_test

import (
	"api/shared/models"
	"api/shared/util"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGenerateSectionStaticText(t *testing.T) {
	section := mockStaticData()
	util.GenerateSectionStaticText(section)

	// Expected results after function execution
	expectedTextOutput := []models.ReportTextOutput{
		{
			Title:  "Text 1",
			Type:   models.Static,
			Input:  "This is test text 1, we're going to splice in every question and csv data label\nq1: **q1\nq2: **q2\nq3: **q3\n\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3",
			Result: "The answer to 2 + 2 = 4",
		},
		{
			Title:  "Text 2",
			Type:   models.Static,
			Input:  "This is test text 2, we're going to splice in every question and csv data label\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3\n\nq1: **q1\nq2: **q2\nq3: **q3\n\n\n",
			Result: "Paris is the capital of France.",
		},
	}

	// Check if the TextOutputs are as expected
	if !reflect.DeepEqual(section.TextOutputs, expectedTextOutput) {
		t.Errorf("TextOutputs were not updated correctly. Got: \n %v \n, want: \n %v", section.TextOutputs, expectedTextOutput)
	}
}

func TestGenerateSectionGeneratorText(t *testing.T) {
	section := mockGeneratorData()

	// Mock function for GeneratePromptResponse
	mockGenerator := MockOpenAiGenerator{}

	err := util.GenerateSectionGeneratorText(mockGenerator, section)
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
func TestAnalyzeOneDimensionalData(t *testing.T) {
	section := mockStaticData()

	csvFile, err := os.Open("../test_files/MasterData.csv")

	if err != nil {
		fmt.Printf("Errror loading file: %v", err)
		return
	}

	util.GenerateSectionCsvDataResults(csvFile, section)

	// Expected results after function execution
	expectedTextOutput := []models.ReportTextOutput{
		{
			Title:  "Text 1",
			Type:   models.Static,
			Input:  "This is test text 1, we're going to splice in every question and csv data label\nq1: **q1\nq2: **q2\nq3: **q3\n\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3",
			Result: "The answer to 2 + 2 = 4",
		},
		{
			Title:  "Text 2",
			Type:   models.Static,
			Input:  "This is test text 2, we're going to splice in every question and csv data label\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3\n\nq1: **q1\nq2: **q2\nq3: **q3\n\n\n",
			Result: "Paris is the capital of France.",
		},
	}

	// Check if the TextOutputs are as expected
	if !reflect.DeepEqual(section.TextOutputs, expectedTextOutput) {
		t.Errorf("TextOutputs were not updated correctly. Got: \n %v \n, want: \n %v", section.TextOutputs, expectedTextOutput)
	}
}

func TestAnalyzeTwoDimensionalData(t *testing.T) {
	section := mockStaticData()

	csvFile, err := os.Open("../test_files/MasterData.csv")

	if err != nil {
		fmt.Printf("Errror loading file: %v", err)
		return
	}

	util.GenerateChartOutputResults(csvFile, section)

	// Expected results after function execution
	expectedTextOutput := []models.ReportTextOutput{
		{
			Title:  "Text 1",
			Type:   models.Static,
			Input:  "This is test text 1, we're going to splice in every question and csv data label\nq1: **q1\nq2: **q2\nq3: **q3\n\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3",
			Result: "The answer to 2 + 2 = 4",
		},
		{
			Title:  "Text 2",
			Type:   models.Static,
			Input:  "This is test text 2, we're going to splice in every question and csv data label\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3\n\nq1: **q1\nq2: **q2\nq3: **q3\n\n\n",
			Result: "Paris is the capital of France.",
		},
	}

	// Check if the TextOutputs are as expected
	if !reflect.DeepEqual(section.TextOutputs, expectedTextOutput) {
		t.Errorf("TextOutputs were not updated correctly. Got: \n %v \n, want: \n %v", section.TextOutputs, expectedTextOutput)
	}
}
