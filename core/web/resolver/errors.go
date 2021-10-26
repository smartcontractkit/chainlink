package resolver

type ErrorCode string

const (
	ErrorCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrorCodeInvalidInput ErrorCode = "INVALID_INPUT"
)

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
