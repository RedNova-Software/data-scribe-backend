package util_test

import "api/shared/models"

// Mock data for testing
func mockStaticData() (*models.Section, []models.Answer) {
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

func mockGeneratorData() (*models.Section, []models.Answer) {
	section := &models.Section{
		Title: "Section Two - Generator",
		Questions: []models.Question{
			{
				Label:    "questionOne",
				Index:    0,
				Question: "What's your favourite color?",
				Answer:   "",
			},
			{
				Label:    "questionTwo",
				Index:    1,
				Question: "What's your favourite city?",
				Answer:   "",
			},
		},
		TextOutputs: []models.TextOutput{
			{
				Title:  "Generator One",
				Input:  "Tell me about this color: **questionOne",
				Type:   models.Generator,
				Index:  0,
				Result: "",
			},
			{
				Title:  "Generator Two",
				Input:  "Tell me about this city: **questionTwo",
				Type:   models.Generator,
				Index:  1,
				Result: "",
			},
		},
	}

	answers := []models.Answer{
		{
			QuestionIndex: 0,
			Answer:        "Blue",
		},
		{
			QuestionIndex: 1,
			Answer:        "Toronto",
		},
	}

	return section, answers
}
