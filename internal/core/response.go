package core

type status string

const (
	_StatusSuccess status = "success"
	_StatusError   status = "error"
)

type Response struct {
	Status status  `json:"status"`
	Data   any     `json:"data,omitempty"`
	Errors []error `json:"errors,omitempty"`
}

func NewSuccessResponse(data any) *Response {
	return &Response{
		Status: _StatusSuccess,
		Data:   data,
	}
}

func NewErrorResponse(erros ...error) *Response {
	return &Response{
		Status: _StatusError,
		Errors: erros,
	}
}

func (r *Response) AddError(err ...error) *Response {
	r.Errors = append(r.Errors, err...)
	return r
}
