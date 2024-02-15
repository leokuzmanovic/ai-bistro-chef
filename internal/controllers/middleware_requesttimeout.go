package controllers

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getMiddlewareFunction(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:                    TimeoutSkipper,
		OnTimeoutRouteErrorHandler: TimeoutErrorHandler,
		Timeout:                    timeout,
	})
}

func TimeoutErrorHandler(err error, c echo.Context) {
	fmt.Println("TimeoutErrorHandler - request timed out")
}

func TimeoutSkipper(c echo.Context) bool {
	//TODO: make this dependent on the environment (dev, test, prod)
	return true
}
