package services

import (
	"github.com/leokuzmanovic/ai-bistro-chef/internal/configuration"
	di "github.com/leokuzmanovic/ai-bistro-chef/internal/dependencyinjection"
	openaiAssistant "github.com/leokuzmanovic/ai-bistro-chef/internal/services/openai"
)

func Wire() {
	var appConfig *configuration.AppConfig = di.Get[*configuration.AppConfig]()

	var openaiAssistant = openaiAssistant.NewOpenaiAssistantImpl(appConfig.GetLocalRecipesPath(), appConfig.GetOpenAiAssistantConfig())

	var assistantService AssistantService = NewAssistantServiceImpl(openaiAssistant)
	di.Register(assistantService)
}
