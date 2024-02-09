package services

import (
	"github.com/leokuzmanovic/ai-bistro-chef/internal/configuration"
	openai "github.com/sashabaranov/go-openai"
)

const CHATGPT_MODEL = "gpt-4-1106-preview"

type AssistantService interface {
	PrepareAssistant() error
}

type AssistantServiceImpl struct {
	client           *openai.Client
	assistantConfig  *configuration.AssistantConfig
	localRecipesPath string
}

func NewAssistantServiceImpl(openAiToken, localRecipesPath string, assistantConfig *configuration.AssistantConfig) *AssistantServiceImpl {
	c := openai.NewClient(openAiToken)
	return &AssistantServiceImpl{
		client:           c,
		assistantConfig:  assistantConfig,
		localRecipesPath: localRecipesPath,
	}
}

func (s *AssistantServiceImpl) PrepareAssistant() error {

	return PrepareAssistant()

}
