package api

const (
	RESOURCE_CONVERSATION   = "conversation"
	RESOURCE_TEXT_MESSAGE   = "text-message"
	RESOURCE_VISUAL_MESSAGE = "visual-message"

	ENDPOINT_CONVERSATION         = "/" + RESOURCE_CONVERSATION
	ENDPOINT_CONVERSATION_MESSAGE = "/" + RESOURCE_CONVERSATION + "/:conversationId/" + RESOURCE_TEXT_MESSAGE
	ENDPOINT_CONVERSATION_VISUAL  = "/" + RESOURCE_CONVERSATION + "/:conversationId/" + RESOURCE_VISUAL_MESSAGE
)

/*
POST /conversation
POST /conversation/:conversationId/text-message
POST /conversation/:conversationId/visual-message


- creates an assistant (chef) if not existing already
- trains it with predefined recipes which represent the favourite foods chef can make (knowledge base)
- prompt chef with text or image to ask about the suggestion for a meal, chef should respond with recipe
- if no meals can be recommended from the knowledge base, chef can browse internet to find something

*/
