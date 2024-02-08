package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/labstack/echo/v4"
	"github.com/leokuzmanovic/ai-bistro-chef/internal"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/configuration"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/controllers"
	di "github.com/leokuzmanovic/ai-bistro-chef/internal/dependencyinjection"
	err "github.com/leokuzmanovic/ai-bistro-chef/internal/errors"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/services"
)

func doSomething(slice *[]int) {
	*slice = append(*slice, 4)
}

func main() {
	mySlice := []int{1, 2, 3}
	doSomething(&mySlice)

	defer func() {
		err.HandlePanic(recover())
	}()
	config := &configuration.AppConfig{}
	di.Register(config)

	e := echo.New()

	internal.WireDependencies()
	controllers.WireControllers(e)

	prepareAssistant()

	e.Logger.Fatal(e.Start(":5000"))
}

func prepareAssistant() {
	assistantService := di.Get[services.AssistantService]()

	go func(assistantService services.AssistantService) {
		defer func() {
			err.HandlePanic(recover())
		}()
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 10 * time.Second

		err := backoff.Retry(func() error {
			err := assistantService.PrepareAssistant()
			if err != nil {
				return errors.New("error while preparing assistant")
			}
			return nil
		}, b)
		if err != nil {
			fmt.Println("All retries exhausted while preparing the assistant!")
			panic(err)
		}
	}(assistantService)
}
