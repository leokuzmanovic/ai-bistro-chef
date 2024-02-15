package api

type ConversationResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type CreateConversationMessageRequest struct {
	Message string `json:"message" validate:"required"`
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

/*


type AuthLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}
*/
