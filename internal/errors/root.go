package errors

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

func LogAndExit(e any) {
	if e == nil {
		return
	}
	err, ok := e.(error)
	if ok {
		fmt.Println("panic is received - exiting for the reason: ", err)
	} else {
		fmt.Println("panic is received - exiting for the reason: ", errors.New("unknown error"))
	}
	os.Exit(1)
}

func IsOpenAINotFoundError(err error) bool {
	var apiErr *openai.APIError
	if ok := errors.As(err, &apiErr); ok {
		return apiErr.HTTPStatusCode == http.StatusNotFound
	}
	return false
}

func GlobalErrorHandler(err error, c echo.Context) {
	code, message, errorType := getErrorDetails(err)

	_ = c.JSON(code, map[string]interface{}{
		"message": message,
		"type":    errorType,
	})
}

func getErrorDetails(err error) (int, string, string) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	errorType := "InternalServerError"

	var appErr AppError
	if ok := errors.As(err, &appErr); ok {
		errorData := appErr.GetErrorData()
		code = errorData.Code
		message = errorData.Message
		errorType = errorData.Type
	} else if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		message = httpError.Message.(string)
		errorType = http.StatusText(code)
	} else {
		// log error message taken from "err.Error()" but do not expose internal errors to the client
	}
	return code, message, errorType
}

type ErrorData struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
}

type AppError interface {
	Error() string
	GetErrorData() ErrorData
}

type InternalServerError struct{}

func (s InternalServerError) Error() string {
	return "InternalServerError"
}

func (s InternalServerError) GetErrorData() ErrorData {
	return ErrorData{Type: "InternalServerError", Code: http.StatusInternalServerError, Message: "Internal Server Error"}
}

type InvalidJsonError struct{}

func (s InvalidJsonError) Error() string {
	return "InvalidJsonError"
}

func (s InvalidJsonError) GetErrorData() ErrorData {
	return ErrorData{Type: "InvalidJsonError", Code: http.StatusBadRequest, Message: "Invalid JSON"}
}

type BadRequestError struct {
}

func (s BadRequestError) Error() string {
	return "BadRequestError"
}

func (s BadRequestError) GetErrorData() ErrorData {
	return ErrorData{Type: "BadRequestError", Code: http.StatusBadRequest, Message: "Invalid request"}
}

type NotFoundError struct {
}

func (s NotFoundError) Error() string {
	return "NotFoundError"
}

func (s NotFoundError) GetErrorData() ErrorData {
	return ErrorData{Type: "NotFoundError", Code: http.StatusNotFound, Message: "Not found"}
}
