package core

import (
	"fmt"
)

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Params  map[string]any `json:"params"`
	Field   string         `json:"field"`
	cause   error          `json:"-"`
}

func NewError(code string, message ...string) *Error {
	e := &Error{
		Code: code,
	}
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

func (e *Error) AddParam(key string, value any) *Error {
	if e.Params == nil {
		e.Params = make(map[string]any)
	}
	e.Params[key] = value
	return e
}

func (e *Error) WithCause(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Params:  e.Params,
		cause:   err,
	}
}

func (e *Error) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("code=%s, message=%v", e.Code, e.Message)
	}
	return fmt.Sprintf("code=%s, message=%v, internal=%v", e.Code, e.Message, e.cause)
}

func (e *Error) Unwrap() error {
	return e.cause
}
