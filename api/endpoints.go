package api

const (
	RESOURCE_CONVERSATIONS = "conversations"
	RESOURCE_MESSAGE       = "message"

	ENDPOINT_CONVERSATIONS                     = "/" + RESOURCE_CONVERSATIONS
	ENDPOINT_CONVERSATIONS_MESSAGE             = "/" + RESOURCE_CONVERSATIONS + "/:conversationId/" + RESOURCE_MESSAGE
	ENDPOINT_CONVERSATIONS_RESPONSE_BY_MESSAGE = "/" + RESOURCE_CONVERSATIONS + "/:conversationId/message/:messageId"
)
