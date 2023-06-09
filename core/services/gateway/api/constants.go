package api

type ErrorCode int

const (
	NoError ErrorCode = iota
	UserMessageParseError
	UnsupportedDONIdError
	InternalHandlerError
	RequestTimeoutError
	NodeReponseEncodingError
	FatalError
)

// See https://www.jsonrpc.org/specification#error_object
func ToJsonRPCErrorCode(errorCode ErrorCode) int {
	gatewayErrorToJsonRPCError := map[ErrorCode]int{
		NoError:                  0,
		UserMessageParseError:    -32700, // Parse Error
		UnsupportedDONIdError:    -32602, // Invalid Params
		InternalHandlerError:     -32000, // Server Error
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
		InternalHandlerError:     500, // Internal Server Error
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
