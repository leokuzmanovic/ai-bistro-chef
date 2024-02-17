package api

type ConversationResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type CreateConversationMessageRequest struct {
	Message  string `json:"message" validate:"required"`
	ImageUrl string `json:"imageUrl"`
}

type ConversationMessageResponse struct {
	Message   string `json:"message"`
	MessageId string `json:"messageId"`
}

type ConversationMessageShortResponse struct {
	Message string `json:"message"`
}

type PrepareAssistantRequest struct {
	ForceUpdate bool `json:"forceUpdate"`
}
