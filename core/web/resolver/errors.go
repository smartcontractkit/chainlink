package resolver

import (
	"database/sql"

	"github.com/pkg/errors"
)

type ErrorCode string

const (
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrorCodeUnprocessable  ErrorCode = "UNPROCESSABLE"
	ErrorCodeStatusConflict ErrorCode = "STATUS_CONFLICT"
)

type ResolvesNotFoundError struct {
	err     error
	message string
}

// ToNotFoundError resolves to the not found error resolver
func (e *ResolvesNotFoundError) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if e.err != nil && errors.Is(e.err, sql.ErrNoRows) {
		return NewNotFoundError(e.message), true
	}

	return nil, false
}

type NotFoundErrorResolver struct {
	message string
	code    ErrorCode
}

func NewNotFoundError(message string) *NotFoundErrorResolver {
	return &NotFoundErrorResolver{
		message: message,
		code:    ErrorCodeNotFound,
	}
}

func (r *NotFoundErrorResolver) Message() string {
	return r.message
}

func (r *NotFoundErrorResolver) Code() ErrorCode {
	return r.code
}

type InputErrorResolver struct {
	path    string
	message string
}

func NewInputError(path, message string) *InputErrorResolver {
	return &InputErrorResolver{
		path:    path,
		message: message,
	}
}

func (r *InputErrorResolver) Path() string {
	return r.path
}

func (r *InputErrorResolver) Message() string {
	return r.message
}

func (r *InputErrorResolver) Code() ErrorCode {
	return ErrorCodeInvalidInput
}

// InputErrorsResolver groups a slice of input errors
type InputErrorsResolver struct {
	iers []*InputErrorResolver
}

func NewInputErrors(iers []*InputErrorResolver) *InputErrorsResolver {
	return &InputErrorsResolver{
		iers: iers,
	}
}

func (r *InputErrorsResolver) Errors() []*InputErrorResolver {
	return r.iers
}
