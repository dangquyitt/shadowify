package apperr

import "fmt"

type AppErr struct {
	Code    string         `json:"code"`
	Message string         `json:"message,omitempty"`
	Field   string         `json:"field,omitempty"`
	Params  map[string]any `json:"params,omitempty"`
	cause   error          `json:"-"`
}

func NewAppErr(code string, message ...string) *AppErr {
	e := &AppErr{
		Code: code,
	}
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

func (e *AppErr) WithParam(key string, value any) *AppErr {
	if e.Params == nil {
		e.Params = make(map[string]any)
	}
	e.Params[key] = value
	return e
}

func (e *AppErr) WithMessage(message string) *AppErr {
	e.Message = message
	return e
}

func (e *AppErr) WithCause(err error) *AppErr {
	e.cause = err
	return e
}

func (e *AppErr) WithField(field string) *AppErr {
	e.Field = field
	return e
}

func (e *AppErr) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("code: %s, message: %s, field: %s, params: %v, cause: %v",
			e.Code, e.Message, e.Field, e.Params, e.cause)
	}
	return fmt.Sprintf("code: %s, message: %s, field: %s, params: %v",
		e.Code, e.Message, e.Field, e.Params)
}

func (e *AppErr) Unwrap() error {
	return e.cause
}
