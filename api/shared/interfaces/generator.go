package interfaces

type Generator interface {
	GeneratePromptResponse(prompt string) (string, error)
}
