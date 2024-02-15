package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff"
	errs "github.com/leokuzmanovic/ai-bistro-chef/internal/errors"
	openai "github.com/sashabaranov/go-openai"
)

type openaiThreadRunManager interface {
	startThreadRun(context.Context, string, string, string) error
	cancelAllThreadRuns(context.Context, string) error
	getThreadRunResult(context.Context, string, string, time.Duration) string
}
type openaiThreadRunManagerImpl struct {
	client            *openai.Client
	resultsMap        *sync.Map
	runMonitorTimeout time.Duration
}

func newOpenaiThreadRunManagerImpl(client *openai.Client, runMonitorTimeout time.Duration) *openaiThreadRunManagerImpl {
	return &openaiThreadRunManagerImpl{
		client:            client,
		runMonitorTimeout: runMonitorTimeout,
		resultsMap:        new(sync.Map),
	}
}

func (s *openaiThreadRunManagerImpl) startThreadRun(ctx context.Context, threadId, messageId string, assistantId string) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	var opErr error
	var run openai.Run

	_ = backoff.Retry(func() error {
		run, opErr = s.client.CreateRun(ctx, threadId, openai.RunRequest{
			AssistantID: assistantId,
		})
		if opErr == nil {
			s.startRunResultPolling(context.Background(), threadId, messageId, &run)
			return nil
		}
		if errs.IsOpenAINotFoundError(opErr) {
			opErr = errs.ThreadNotFoundError{}
			return nil
		} else if errs.IsOpenAIThreadHasActiveRunError(opErr) {
			fmt.Println("Thread has active run, cancelling all runs")
			opErr = s.cancelAllThreadRuns(ctx, threadId)
			if opErr != nil {
				fmt.Println("Error while cancelling all runs: ", opErr.Error())
				return opErr
			}
			b.Reset()
			return errors.New("retry")
		} else if opErr != nil {
			return opErr
		}
		return opErr
	}, b)

	if opErr != nil {
		fmt.Println("Error while creating run: ", opErr.Error())
		return opErr
	}
	fmt.Println("Run created")
	return nil
}

func (s *openaiThreadRunManagerImpl) startRunResultPolling(ctx context.Context, threadId, messageId string, run *openai.Run) {
	go func(ctx context.Context, threadId, messageId string, run *openai.Run) {
		rc := make(chan bool, 1)
		var expired atomic.Bool = atomic.Bool{}
		expired.Store(false)

		// start the polling, but monitor for cancleation signals
		go func() {
			for {
				if expired.Load() {
					break
				}
				s.tryGetRunResult(ctx, threadId, messageId, run, &expired)
			}
			rc <- true
		}()

		select {
		case <-ctx.Done():
			s.updateThreadRunResult(ctx, threadId, messageId, "", errors.New("context cancelled"))
			expired.Store(true)
			fmt.Println("Context cancelled")
		case <-time.After(s.runMonitorTimeout):
			s.updateThreadRunResult(ctx, threadId, messageId, "", errors.New("run monitor timeout"))
			expired.Store(true)
			fmt.Println("Run monitor timeout")
		case <-rc:
			fmt.Println("Run result stored")
		}
	}(ctx, threadId, messageId, run)
}

func (s *openaiThreadRunManagerImpl) tryGetRunResult(ctx context.Context, threadId, messageId string, run *openai.Run, stopSignal *atomic.Bool) {
	var err error
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 1 * time.Minute

	_ = backoff.Retry(func() error {
		if stopSignal.Load() {
			fmt.Println("Stop signal received")
			return nil
		}
		run, err = s.tryRetrieveFinishedRun(ctx, run.ID, threadId)
		if err != nil && errors.Is(err, errs.ThreadNotFoundError{}) {
			s.updateThreadRunResult(ctx, threadId, messageId, "", err)
			stopSignal.Store(true)
			fmt.Println("Thread not found")
			return nil
		}
		if err != nil || run == nil {
			fmt.Println("run is not ready")
			return errors.New("run is not ready")
		}

		// from this point on, run is completed
		if run.Status == openai.RunStatusCompleted {
			limit := 10
			order := "desc"
			messages, err := s.client.ListMessage(ctx, threadId, &limit, &order, &messageId, nil)
			if errs.IsOpenAINotFoundError(err) {
				s.updateThreadRunResult(ctx, threadId, messageId, "", err)
				stopSignal.Store(true)
				fmt.Println("Thread not found")
				return nil
			} else if err != nil {
				fmt.Println("Error while retrieving messages: ", err.Error())
				return err
			}

			resultMessage := ""
			if len(messages.Messages) <= 0 {
				fmt.Println("No response from the assistant")
				resultMessage = "No response from the assistant"
			} else {
				for _, message := range messages.Messages {
					if message.ID == messageId {
						break
					}
					if message.Role == openai.ChatMessageRoleAssistant {
						// because we are iterating in reverse order, we need to prepend the new message
						resultMessage = message.Content[0].Text.Value + "\n\n" + resultMessage
					} else {
						break
					}
				}
				if len(resultMessage) > 0 {
					resultMessage = resultMessage[:len(resultMessage)-2]
				}
				fmt.Println("Assistant response parsed")
			}

			s.updateThreadRunResult(ctx, threadId, messageId, resultMessage, nil)
			stopSignal.Store(true)
			return nil
		}
		fmt.Println("Run is not completed")
		return nil
	}, b)
}

func (s *openaiThreadRunManagerImpl) updateThreadRunResult(ctx context.Context, threadId, messageId, result string, err error) {
	if err != nil {
		s.resultsMap.Store(threadId+messageId, getErrorResult(err))
	} else {
		s.resultsMap.Store(threadId+messageId, result)
	}
}

func getErrorResult(err error) string {
	// todo: make sure we do not expose internal errors to the client
	return fmt.Sprintf("Error: %s", err)
}

func (s *openaiThreadRunManagerImpl) tryRetrieveFinishedRun(ctx context.Context, runId string, threadId string) (*openai.Run, error) {
	var err error
	var run openai.Run

	run, err = s.client.RetrieveRun(ctx, threadId, runId)
	if errs.IsOpenAINotFoundError(err) {
		fmt.Println("Thread not found")
		return &run, errs.ThreadNotFoundError{}
	} else if err != nil {
		fmt.Println("Error while retrieving run: ", err.Error())
		return &run, err
	}
	if run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress || errs.IsOpenAIRateLimitExcededError(run.LastError) {
		fmt.Println("Continuing to poll for run results")
		return &run, fmt.Errorf("continue polling")
	}
	return &run, err
}

func (s *openaiThreadRunManagerImpl) cancelAllThreadRuns(ctx context.Context, threadId string) error {
	limit := 20
	order := "desc"

	runListResponse, err := s.client.ListRuns(context.Background(), threadId, openai.Pagination{
		Limit: &limit,
		Order: &order,
	})
	if err != nil {
		fmt.Println("Error while retrieving runs: ", err.Error())
		return err
	}

	if runListResponse.Runs == nil || len(runListResponse.Runs) == 0 {
		fmt.Println("No runs found")
		return nil
	}

	for _, run := range runListResponse.Runs {
		if run.Status == openai.RunStatusInProgress || run.Status == openai.RunStatusQueued {
			_, err = s.client.CancelRun(ctx, threadId, run.ID)
			if errs.IsOpenAICannotCancelFinishedRunError(err) {
				fmt.Println("Run is already finished")
				continue
			} else if err != nil {
				fmt.Println("Error while cancelling run: ", err.Error())
				return err
			}
		}
	}
	return nil
}

func (s *openaiThreadRunManagerImpl) getThreadRunResult(ctx context.Context, threadId, messageId string, timeout time.Duration) string {
	rc := make(chan string, 1)
	result := ""
	expired := false

	go func() {
		for {
			if expired {
				break
			}
			if val, ok := s.resultsMap.Load(threadId + messageId); ok {
				rc <- val.(string)
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled on getting run result")
		expired = true
	case <-time.After(timeout):
		fmt.Println("Timeout on getting run result")
		expired = true
	case result = <-rc:
		fmt.Println("Run result retrieved")
	}

	return result
}
