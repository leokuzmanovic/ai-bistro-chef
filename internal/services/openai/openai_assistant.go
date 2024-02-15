package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/configuration"
	errs "github.com/leokuzmanovic/ai-bistro-chef/internal/errors"
	openai "github.com/sashabaranov/go-openai"
)

type openaiAssistantImpl struct {
	client           *openai.Client
	assistantConfig  *configuration.OpenAiAssistantConfig
	assistantBuilder openaiAssistantBuilder
	threadRunManager openaiThreadRunManager
	assistant        *openai.Assistant
}

func NewOpenaiAssistantImpl(localRecipesPath string, assistantConfig *configuration.OpenAiAssistantConfig) *openaiAssistantImpl {
	c := openai.NewClient(*assistantConfig.GetOpenAiToken())
	return &openaiAssistantImpl{
		client:           c,
		assistantConfig:  assistantConfig,
		assistantBuilder: newOpenaiAssistantBuilderImpl(c, localRecipesPath, assistantConfig),
		threadRunManager: newOpenaiThreadRunManagerImpl(c, assistantConfig.GetThreadRunTimeout()),
	}
}

func (s *openaiAssistantImpl) SetupAssistant() error {
	if s.assistantConfig.GetAssistantId() != nil {
		s.assistant = &openai.Assistant{ID: *s.assistantConfig.GetAssistantId()}
		return nil
	}

	assistant, err := s.assistantBuilder.initAssistant()
	if err != nil {
		return err
	}
	s.assistant = assistant
	return nil
}

func (s *openaiAssistantImpl) CreateNewThread(ctx context.Context) (string, error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	var thread openai.Thread
	var err error

	err = backoff.Retry(func() error {
		thread, err = s.client.CreateThread(ctx, openai.ThreadRequest{})
		if err != nil {
			fmt.Println("Error while creating thread: ", err.Error())
			return err
		}
		return nil
	}, b)

	if err != nil {
		return "", err
	}
	fmt.Println("Thread created")
	return thread.ID, err
}

func (s *openaiAssistantImpl) AddMessageToThread(ctx context.Context, threadId, message string) (string, error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	var msg openai.Message
	var err error
	var opErr error

	_ = backoff.Retry(func() error {
		msg, err = s.client.CreateMessage(ctx, threadId, openai.MessageRequest{
			Role:    string(openai.ThreadMessageRoleUser),
			Content: message,
		})
		if err == nil {
			return nil
		}
		opErr = err
		if errs.IsOpenAINotFoundError(opErr) {
			opErr = errs.ThreadNotFoundError{}
			return nil
		} else if errs.IsOpenAICannotAddMessageToRunningThreadError(opErr) {
			opErr = s.threadRunManager.cancelAllThreadRuns(ctx, threadId)
			if opErr != nil {
				return opErr
			}
			b.Reset()
			return errors.New("retry")
		}
		return opErr
	}, b)

	if opErr != nil {
		fmt.Println("Error while adding message to thread: ", opErr.Error())
		return "", opErr
	}

	fmt.Println("Message added to thread")
	return msg.ID, nil
}

func (s *openaiAssistantImpl) StartThreadProcessing(ctx context.Context, threadId, messageId string) error {
	if s.assistant == nil {
		return errs.AssistantNotReadyError{}
	}
	return s.threadRunManager.startThreadRun(ctx, threadId, messageId, s.assistant.ID)
}

func (s *openaiAssistantImpl) GetMessageResponse(ctx context.Context, threadId, messageId string, timeout time.Duration) string {
	return s.threadRunManager.getThreadRunResult(ctx, threadId, messageId, timeout)
}
