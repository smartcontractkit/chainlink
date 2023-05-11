package gateway

type ErrorCode int

const (
	NoError ErrorCode = iota
	UserMessageParseError
	FatalError
)
