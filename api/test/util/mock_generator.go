package util_test

type MockOpenAiGenerator struct{}

func (m MockOpenAiGenerator) GeneratePromptResponse(prompt string) (string, error) {
	// Mock responses for GeneratePromptResponse (you will need to implement this)
	mockResponses := map[string]string{
		"Tell me about this color: Blue":   "Blue is a calming color",
		"Tell me about this city: Toronto": "Toronto is a vibrant city",
	}
	return mockResponses[prompt], nil
}
