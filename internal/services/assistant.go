package services

import (
	"context"
	"fmt"
	"time"
)

type AssistantService interface {
	PrepareAssistant() error
	StartNewConversation(context.Context) (string, error)
	AskAssistant(context.Context, string, string) (string, string, error)
	GetMessageResponseFromAssistant(context.Context, string, string, time.Duration) string
}

type aiAssistant interface {
	SetupAssistant() error
	CreateNewThread(context.Context) (string, error)
	AddMessageToThread(context.Context, string, string) (string, error)
	StartThreadProcessing(context.Context, string, string) error
	GetMessageResponse(context.Context, string, string, time.Duration) string
}

type AssistantServiceImpl struct {
	aiAssistant aiAssistant
}

func NewAssistantServiceImpl(aiAssistant aiAssistant) *AssistantServiceImpl {
	return &AssistantServiceImpl{
		aiAssistant: aiAssistant,
	}
}

func (s *AssistantServiceImpl) PrepareAssistant() error {
	err := s.aiAssistant.SetupAssistant()
	if err != nil {
		return err
	}
	return nil
}

func (s *AssistantServiceImpl) StartNewConversation(ctx context.Context) (string, error) {
	return s.aiAssistant.CreateNewThread(ctx)
}

func (s *AssistantServiceImpl) AskAssistant(ctx context.Context, conversationId, message string) (string, string, error) {
	messageId, err := s.aiAssistant.AddMessageToThread(ctx, conversationId, message)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Message added to thread")

	err = s.aiAssistant.StartThreadProcessing(ctx, conversationId, messageId)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Thread processing started")

	response := s.aiAssistant.GetMessageResponse(ctx, conversationId, messageId, 2*time.Second)
	return response, messageId, nil
}

func (s *AssistantServiceImpl) GetMessageResponseFromAssistant(ctx context.Context, conversationId, messageId string, timeout time.Duration) string {
	response := s.aiAssistant.GetMessageResponse(ctx, conversationId, messageId, timeout)
	return response
}
