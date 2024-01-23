package util

import (
	"api/shared/constants"
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAiGenerator struct{}

func (g OpenAiGenerator) GeneratePromptResponse(prompt string) (string, error) {
	key := os.Getenv(constants.OpenAIKey) // Assuming OpenAIKey is the environment variable name

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo, // Ensure this constant is defined in the `openai` package
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser, // Ensure this constant is defined
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
