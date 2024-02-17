package controllers

import (
	"net/http"
	"time"

	"github.com/leokuzmanovic/ai-bistro-chef/api"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/constants"
	di "github.com/leokuzmanovic/ai-bistro-chef/internal/dependencyinjection"
	errs "github.com/leokuzmanovic/ai-bistro-chef/internal/errors"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/services"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
)

type ConversationsController struct {
	assistantService services.AssistantService
}

func wireConversations(e *echo.Echo) {
	controller := &ConversationsController{
		assistantService: di.Get[services.AssistantService](),
	}

	e.POST(api.ENDPOINT_CONVERSATIONS, controller.newConversation, getMiddlewareFunction(constants.TimeoutRegular))
	e.POST(api.ENDPOINT_CONVERSATIONS_MESSAGE, withValidJsonBody(controller.newConversationTextMessage), getMiddlewareFunction(constants.TimeoutRegular))
	e.GET(api.ENDPOINT_CONVERSATIONS_RESPONSE_BY_MESSAGE, controller.getConversationResponseByMessage, getMiddlewareFunction(constants.TimeoutRegular))
}

/*
// @Summary Create new conversation
// @Tags conversations
// @Router /conversations [POST]
// @Success 201
// @Failure 400 {object} errors.AppError "BadRequestError"
*/
func (s *ConversationsController) newConversation(ctx echo.Context) error {
	conversationId, err := s.assistantService.StartNewConversation(ctx.Request().Context())
	if err != nil {
		return errors.Wrap(err, "conversations")
	}

	response := api.ConversationResponse{
		Id:        conversationId,
		CreatedAt: time.Now().String(),
	}
	return errors.Wrap(ctx.JSON(http.StatusCreated, response), "json")
}

/*
// @Summary Create new text message in conversation
// @Tags conversations/:conversationId/message-text
// @Router conversations/{conversationId}/message-text [POST]
// @Param conversationId path string true "conversation id"
// @Produce json
// @Success 200 {object} api.ConversationMessageResponse
// @Success 400 {object} errors.AppError "BadRequestError"
// @Success 404 {object} errors.AppError "NotFoundError"
*/
func (s *ConversationsController) newConversationTextMessage(ctx echo.Context, body *api.CreateConversationMessageRequest) error {
	conversationId := ctx.Param("conversationId")
	if conversationId == "" {
		return &errs.BadRequestError{}
	}

	messageResponse, messageResponseId, err := s.assistantService.AskAssistant(ctx.Request().Context(), conversationId, body)
	if err != nil {
		return errors.Wrap(err, "conversation")
	}

	response := api.ConversationMessageResponse{
		Message:   messageResponse,
		MessageId: messageResponseId,
	}
	return errors.Wrap(ctx.JSON(http.StatusCreated, response), "json")
}

/*
// @Summary Get conversation response by message id
// @Tags conversations/:conversationId/message/:messageId
// @Router conversations/{conversationId}/message/:messageId [GET]
// @Param conversationId path string true "conversation id"
// @Param messageId path string true "message id"
// @Produce json
// @Success 200 {object} api.ConversationMessageResponse
// @Success 400 {object} errors.AppError "BadRequestError"
// @Success 404 {object} errors.AppError "NotFoundError"
*/
func (s *ConversationsController) getConversationResponseByMessage(ctx echo.Context) error {
	conversationId := ctx.Param("conversationId")
	if conversationId == "" {
		return &errs.BadRequestError{}
	}
	messageId := ctx.Param("messageId")
	if messageId == "" {
		return &errs.BadRequestError{}
	}
	timeout := time.Second * 2 // default timeout
	timeoutParam := ctx.QueryParam("timeout")
	if timeoutParam != "" {
		var err error
		timeout, err = time.ParseDuration(timeoutParam)
		if err != nil {
			return &errs.BadRequestError{}
		}
	}

	messageResponse := s.assistantService.GetMessageResponseFromAssistant(ctx.Request().Context(), conversationId, messageId, timeout)
	response := api.ConversationMessageShortResponse{
		Message: messageResponse,
	}
	return errors.Wrap(ctx.JSON(http.StatusOK, response), "json")
}
