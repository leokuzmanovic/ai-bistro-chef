package configuration

type AppConfig struct {
	AssistantConfig *AssistantConfig
}

type AssistantConfig struct {
}

func (s *AppConfig) GetAssistantConfig() *AssistantConfig {
	return &AssistantConfig{}
}

func (s *AssistantConfig) GetAssistantName() *string {
	name := "ai-bistro-chef-assistant-v1"
	return &name
}

func (s *AssistantConfig) GetAssistantDescription() *string {
	description := "AI Bistro Chef Assistant"
	return &description
}

func (s *AssistantConfig) GetAssistantInstructions() *string {
	instructions := `Your are an experienced chef with more than 20 year of experience cooking different cousines, and you are asked for consultation on what to cook and how.
		Write a recipe for the user once they ask for suggestions with or without providing the ingredients they have at the moment.
		Use the recipes from the files you are provided as recommended meals but if you are not able to use them for the given ingredients list, then you can recommend other meals as well.
		This recipe should contain the ingredients list, cooking time, and a list of instructions in the bullet point format.
		Use clear and concise language and write in a confident yet slightly humorous tone.`
	return &instructions
}

func (s *AppConfig) GetOpenAiToken() string {
	return "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"
	// return "sk-9Z6Z6QZ6QZ6QZ6QZ6QZ6QZ6Q"
}

func (s *AppConfig) GetLocalRecipesPath() string {
	return "../recipes"
}
