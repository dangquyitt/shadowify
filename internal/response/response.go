package response

import (
	"shadowify/internal/apperr"
	"shadowify/internal/pagination"
)

type Response struct {
	Code       string                 `json:"code"`
	Data       any                    `json:"data,omitempty"`
	Metadata   any                    `json:"metadata,omitempty"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
	Errors     []*apperr.AppErr       `json:"errors,omitempty"`
}

func NewSuccessResponse(data any) *Response {
	return &Response{
		Code: "success",
		Data: data,
	}
}

func NewErrorResponse(errs ...*apperr.AppErr) *Response {
	return &Response{
		Code:   "error",
		Errors: errs,
	}
}

func (r *Response) WithMetadata(metadata any) *Response {
	r.Metadata = metadata
	return r
}

func (r *Response) WithPagination(pagination *pagination.Pagination) *Response {
	r.Pagination = pagination
	return r
}

func (r *Response) WithErrors(errs ...*apperr.AppErr) *Response {
	r.Errors = append(r.Errors, errs...)
	return r
}

func (r *Response) WithError(err *apperr.AppErr) *Response {
	r.Errors = append(r.Errors, err)
	return r
}
