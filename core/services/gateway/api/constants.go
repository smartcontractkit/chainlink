package api

type ErrorCode int

const (
	NoError ErrorCode = iota
	UserMessageParseError
	UnsupportedDONIdError
	HandlerError
	RequestTimeoutError
	NodeReponseEncodingError
	FatalError
)

func (e ErrorCode) String() string {
	switch e {
	case NoError:
		return "NoError"
	case UserMessageParseError:
		return "UserMessageParseError"
	case UnsupportedDONIdError:
		return "UnsupportedDONIdError"
	case HandlerError:
		return "HandlerError"
	case RequestTimeoutError:
		return "RequestTimeoutError"
	case NodeReponseEncodingError:
		return "NodeReponseEncodingError"
	case FatalError:
		return "FatalError"
	default:
		return "UnknownError"
	}
}

// See https://www.jsonrpc.org/specification#error_object
func ToJsonRPCErrorCode(errorCode ErrorCode) int {
	gatewayErrorToJsonRPCError := map[ErrorCode]int{
		NoError:                  0,
		UserMessageParseError:    -32700, // Parse Error
		UnsupportedDONIdError:    -32602, // Invalid Params
		HandlerError:             -32600, // Invalid Request
		RequestTimeoutError:      -32000, // Server Error
		NodeReponseEncodingError: -32603, // Internal Error
		FatalError:               -32000, // Server Error
	}

	code, ok := gatewayErrorToJsonRPCError[errorCode]
	if !ok {
		return -32000
	}
	return code
}

// See https://go.dev/src/net/http/status.go
func ToHttpErrorCode(errorCode ErrorCode) int {
	gatewayErrorToHttpError := map[ErrorCode]int{
		NoError:                  200, // OK
		UserMessageParseError:    400, // Bad Request
		UnsupportedDONIdError:    400, // Bad Request
		HandlerError:             400, // Bad Request
		RequestTimeoutError:      504, // Gateway Timeout
		NodeReponseEncodingError: 500, // Internal Server Error
		FatalError:               500, // Internal Server Error
	}

	code, ok := gatewayErrorToHttpError[errorCode]
	if !ok {
		return 500
	}
	return code
}
