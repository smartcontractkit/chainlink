package api

import "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

type ErrorCode int

const (
	NoError ErrorCode = iota
	UserMessageParseError
	UnsupportedDONIdError
	HandlerError
	RequestTimeoutError
	NodeReponseEncodingError
	FatalError
	UnsupportedMethodError
	InvalidParamsError
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
	case UnsupportedMethodError:
		return "UnsupportedMthodError"
	case InvalidParamsError:
		return "InvalidParamsError"
	default:
		return "UnknownError"
	}
}

// See https://www.jsonrpc.org/specification#error_object
func ToJSONRPCErrorCode(errorCode ErrorCode) int64 {
	gatewayErrorToJSONRPCError := map[ErrorCode]int64{
		NoError:                  0,
		UserMessageParseError:    jsonrpc2.ErrParse,            // Parse Error
		UnsupportedDONIdError:    jsonrpc2.ErrInvalidParams,    // Invalid Params
		InvalidParamsError:       jsonrpc2.ErrInvalidParams,    // Invalid Params
		HandlerError:             jsonrpc2.ErrInvalidRequest,   // Invalid Request
		RequestTimeoutError:      jsonrpc2.ErrServerOverloaded, // Server Error
		NodeReponseEncodingError: jsonrpc2.ErrInternal,         // Internal Error
		FatalError:               jsonrpc2.ErrInternal,         // Internal Error
		UnsupportedMethodError:   jsonrpc2.ErrMethodNotFound,   // Method Not Found
	}

	code, ok := gatewayErrorToJSONRPCError[errorCode]
	if !ok {
		return jsonrpc2.ErrInternal
	}
	return code
}

// See https://go.dev/src/net/http/status.go
func ToHttpErrorCode(errorCode ErrorCode) int {
	gatewayErrorToHttpError := map[ErrorCode]int{
		NoError:                  200, // OK
		UserMessageParseError:    400, // Bad Request
		UnsupportedDONIdError:    400, // Bad Request
		UnsupportedMethodError:   400, // Bad Request
		HandlerError:             400, // Bad Request
		InvalidParamsError:       400, // Bad Request
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
