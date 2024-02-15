package api

const (
	RESOURCE_CONVERSATIONS = "conversations"
	RESOURCE_MESSAGE       = "message"

	ENDPOINT_CONVERSATIONS                     = "/" + RESOURCE_CONVERSATIONS
	ENDPOINT_CONVERSATIONS_MESSAGE             = "/" + RESOURCE_CONVERSATIONS + "/:conversationId/" + RESOURCE_MESSAGE
	ENDPOINT_CONVERSATIONS_RESPONSE_BY_MESSAGE = "/" + RESOURCE_CONVERSATIONS + "/:conversationId/message/:messageId"
)

/*
POST /conversations
POST /conversations/:conversationId/message-text
POST /conversations/:conversationId/message-text-visual

- creates an assistant (chef) if not existing already
- trains it with predefined recipes which represent the favourite foods chef can make (knowledge base)
- prompt chef with text or image to ask about the suggestion for a meal, chef should respond with recipe
- if no meals can be recommended from the knowledge base, chef can browse internet to find something

*/
