package util_test

import "api/shared/models"

// Mock data for testing
func mockStaticData() (*models.ReportSection, []models.Answer) {
	section := &models.ReportSection{
		Title: "Test Section",
		Questions: []models.ReportQuestion{
			{
				Label:    "question1",
				Question: "What is 2+2?",
				Answer:   "",
			},
			{
				Label:    "question2",
				Question: "What is the capital of France?",
				Answer:   "",
			},
		},
		TextOutputs: []models.ReportTextOutput{
			{
				Title:  "Output1",
				Type:   models.Static,
				Input:  "The answer to 2 + 2 = **question1",
				Result: "",
			},
			{
				Title:  "Output2",
				Type:   models.Static,
				Input:  "**question2 is the capital of France.",
				Result: "",
			},
		},
	}

	answers := []models.Answer{
		{
			Answer: "4",
		},
		{
			Answer: "Paris",
		},
	}

	return section, answers
}

func mockGeneratorData() (*models.ReportSection, []models.Answer) {
	section := &models.ReportSection{
		Title: "Section Two - Generator",
		Questions: []models.ReportQuestion{
			{
				Label:    "questionOne",
				Question: "What's your favourite color?",
				Answer:   "",
			},
			{
				Label:    "questionTwo",
				Question: "What's your favourite city?",
				Answer:   "",
			},
		},
		TextOutputs: []models.ReportTextOutput{
			{
				Title:  "Generator One",
				Input:  "Tell me about this color: **questionOne",
				Type:   models.Generator,
				Result: "",
			},
			{
				Title:  "Generator Two",
				Input:  "Tell me about this city: **questionTwo",
				Type:   models.Generator,
				Result: "",
			},
		},
	}

	answers := []models.Answer{
		{
			Answer: "Blue",
		},
		{
			Answer: "Toronto",
		},
	}

	return section, answers
}
