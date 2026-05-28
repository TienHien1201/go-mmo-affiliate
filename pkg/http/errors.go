package xhttp

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Field   string `json:"field,omitempty"`
    Status  int    `json:"-"` 
    Err     error  `json:"-"` 
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

func NewAppError(status int, code, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Status:  status,
    }
}

func (e *AppError) Wrap(err error) *AppError {
    e.Err = err
    return e
}

func (e *AppError) WithField(field string) *AppError {
    e.Field = field
    return e
}

func NotFoundErrorf(format string, a ...any) *AppError {
    return NewAppError(http.StatusNotFound, "ERR_NOT_FOUND", fmt.Sprintf(format, a...))
}

func BadRequestErrorf(format string, a ...any) *AppError {
    return NewAppError(http.StatusBadRequest, "ERR_BAD_REQUEST", fmt.Sprintf(format, a...))
}

func UnauthorizedErrorf(format string, a ...any) *AppError {
    return NewAppError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", fmt.Sprintf(format, a...))
}

func ForbiddenErrorf(format string, a ...any) *AppError {
    return NewAppError(http.StatusForbidden, "ERR_FORBIDDEN", fmt.Sprintf(format, a...))
}

func InternalErrorf(format string, a ...any) *AppError {
    return NewAppError(http.StatusInternalServerError, "ERR_INTERNAL", fmt.Sprintf(format, a...))
}

func IsAppError(err error) (*AppError, bool) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        return appErr, true
    }
    return nil, false
}