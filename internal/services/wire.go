package services

import (
	"github.com/leokuzmanovic/ai-bistro-chef/internal/configuration"
	di "github.com/leokuzmanovic/ai-bistro-chef/internal/dependencyinjection"
)

func Wire() {
	var appConfig *configuration.AppConfig = di.Get[*configuration.AppConfig]()

	var assistantService AssistantService = NewAssistantServiceImpl(appConfig.GetOpenAiToken(), appConfig.GetLocalRecipesPath(), appConfig.GetAssistantConfig())
	di.Register(assistantService)
}
