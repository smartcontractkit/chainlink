package resolver

type NotFoundErrorResolver struct {
	message string
	code    string
}

func NewNotFoundError(message string) *NotFoundErrorResolver {
	return &NotFoundErrorResolver{
		message: message,
		code:    "NOT_FOUND",
	}
}

func (r *NotFoundErrorResolver) Message() string {
	return r.message
}

func (r *NotFoundErrorResolver) Code() string {
	return r.code
}
