package controllers

import (
	"github.com/labstack/echo/v4"
	er "github.com/leokuzmanovic/ai-bistro-chef/internal/errors"
)

func WireControllers(e *echo.Echo) {
	e.HTTPErrorHandler = er.GlobalErrorHandler
	wireConversations(e)
}
