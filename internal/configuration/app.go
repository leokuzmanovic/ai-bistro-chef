package configuration

import "time"

type AppConfig struct {
	AssistantConfig *OpenAiAssistantConfig
}

type OpenAiAssistantConfig struct {
}

func (s *AppConfig) GetOpenAiAssistantConfig() *OpenAiAssistantConfig {
	return &OpenAiAssistantConfig{}
}

func (s *OpenAiAssistantConfig) GetAssistantId() *string {
	return nil
}

func (s *OpenAiAssistantConfig) GetAssistantName() *string {
	name := "ai-bistro-chef-assistant-v1"
	return &name
}

func (s *OpenAiAssistantConfig) GetAssistantDescription() *string {
	description := "AI Bistro Chef Assistant"
	return &description
}

func (s *OpenAiAssistantConfig) GetAssistantInstructions() *string {
	instructions := `Your are an experienced chef with more than 20 year of experience cooking different cousines, and you are asked for consultation on what to cook and how.
		Write a recipe for the user once they ask for suggestions with or without providing the ingredients they have at the moment.
		Use the recipes from the files you are provided as recommended meals but if you are not able to use them for the given ingredients list, then you can recommend other meals as well.
		This recipe should contain the ingredients list, cooking time, and a list of instructions in the bullet point format.
		Use clear and concise language and write in a confident yet slightly humorous tone.`
	return &instructions
}

func (s *OpenAiAssistantConfig) GetAssistantGptModel() *string {
	instructions := "gpt-4-1106-preview"
	return &instructions
}

func (s *OpenAiAssistantConfig) GetOpenAiToken() *string {
	return nil
}

func (s *OpenAiAssistantConfig) GetThreadRunTimeout() time.Duration {
	return 10 * time.Minute
}

func (s *AppConfig) GetLocalRecipesPath() string {
	return "../recipes"
}
