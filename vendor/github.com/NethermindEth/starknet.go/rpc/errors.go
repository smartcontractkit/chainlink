package rpc

import (
	"encoding/json"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

const (
	InvalidJSON    = -32700 // Invalid JSON was received by the server.
	InvalidRequest = -32600 // The JSON sent is not a valid Request object.
	MethodNotFound = -32601 // The method does not exist / is not available.
	InvalidParams  = -32602 // Invalid method parameter(s).
	InternalError  = -32603 // Internal JSON-RPC error.
)

// Err returns an RPCError based on the given code and data.
//
// Parameters:
// - code: an integer representing the error code.
// - data: any data associated with the error.
// Returns
// - *RPCError: a pointer to an RPCError object.
func Err(code int, data any) *RPCError {
	switch code {
	case InvalidJSON:
		return &RPCError{Code: InvalidJSON, Message: "Parse error", Data: data}
	case InvalidRequest:
		return &RPCError{Code: InvalidRequest, Message: "Invalid Request", Data: data}
	case MethodNotFound:
		return &RPCError{Code: MethodNotFound, Message: "Method Not Found", Data: data}
	case InvalidParams:
		return &RPCError{Code: InvalidParams, Message: "Invalid Params", Data: data}
	default:
		return &RPCError{Code: InternalError, Message: "Internal Error", Data: data}
	}
}

// tryUnwrapToRPCErr unwraps the error and checks if it matches any of the given RPC errors.
// If a match is found, the corresponding RPC error is returned.
// If no match is found, the function returns an InternalError with the original error.
//
// Parameters:
// - err: The error to be unwrapped
// - rpcErrors: variadic list of *RPCError objects to be checked
// Returns:
// - error: the original error
func tryUnwrapToRPCErr(err error, rpcErrors ...*RPCError) *RPCError {
	errBytes, errIn := json.Marshal(err)
	if errIn != nil {
		return Err(InternalError, errIn.Error())
	}

	var nodeErr RPCError
	errIn = json.Unmarshal(errBytes, &nodeErr)
	if errIn != nil {
		return Err(InternalError, errIn.Error())
	}

	for _, rpcErr := range rpcErrors {
		if nodeErr.Code == rpcErr.Code && nodeErr.Message == rpcErr.Message {
			return &nodeErr
		}
	}
	return Err(InternalError, err.Error())
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e RPCError) Error() string {
	return e.Message
}

var (
	ErrFailedToReceiveTxn = &RPCError{
		Code:    1,
		Message: "Failed to write transaction",
	}
	ErrNoTraceAvailable = &RPCError{
		Code:    10,
		Message: "No trace available for transaction",
	}
	ErrContractNotFound = &RPCError{
		Code:    20,
		Message: "Contract not found",
	}
	ErrBlockNotFound = &RPCError{
		Code:    24,
		Message: "Block not found",
	}
	ErrInvalidTxnHash = &RPCError{
		Code:    25,
		Message: "Invalid transaction hash",
	}
	ErrInvalidBlockHash = &RPCError{
		Code:    26,
		Message: "Invalid block hash",
	}
	ErrInvalidTxnIndex = &RPCError{
		Code:    27,
		Message: "Invalid transaction index in a block",
	}
	ErrClassHashNotFound = &RPCError{
		Code:    28,
		Message: "Class hash not found",
	}
	ErrHashNotFound = &RPCError{
		Code:    29,
		Message: "Transaction hash not found",
	}
	ErrPageSizeTooBig = &RPCError{
		Code:    31,
		Message: "Requested page size is too big",
	}
	ErrNoBlocks = &RPCError{
		Code:    32,
		Message: "There are no blocks",
	}
	ErrInvalidContinuationToken = &RPCError{
		Code:    33,
		Message: "The supplied continuation token is invalid or unknown",
	}
	ErrTooManyKeysInFilter = &RPCError{
		Code:    34,
		Message: "Too many keys provided in a filter",
	}
	ErrContractError = &RPCError{
		Code:    40,
		Message: "Contract error",
	}
	ErrTxnExec = &RPCError{
		Code:    41,
		Message: "Transaction execution error",
	}
	ErrInvalidContractClass = &RPCError{
		Code:    50,
		Message: "Invalid contract class",
	}
	ErrClassAlreadyDeclared = &RPCError{
		Code:    51,
		Message: "Class already declared",
	}
	ErrInvalidTransactionNonce = &RPCError{
		Code:    52,
		Message: "Invalid transaction nonce",
	}
	ErrInsufficientMaxFee = &RPCError{
		Code:    53,
		Message: "Max fee is smaller than the minimal transaction cost (validation plus fee transfer)",
	}
	ErrInsufficientAccountBalance = &RPCError{
		Code:    54,
		Message: "Account balance is smaller than the transaction's max_fee",
	}
	ErrValidationFailure = &RPCError{
		Code:    55,
		Message: "Account validation failed",
	}
	ErrCompilationFailed = &RPCError{
		Code:    56,
		Message: "Compilation failed",
	}
	ErrContractClassSizeTooLarge = &RPCError{
		Code:    57,
		Message: "Contract class size is too large",
	}
	ErrNonAccount = &RPCError{
		Code:    58,
		Message: "Sender address is not an account contract",
	}
	ErrDuplicateTx = &RPCError{
		Code:    59,
		Message: "A transaction with the same hash already exists in the mempool",
	}
	ErrCompiledClassHashMismatch = &RPCError{
		Code:    60,
		Message: "The compiled class hash did not match the one supplied in the transaction",
	}
	ErrUnsupportedTxVersion = &RPCError{
		Code:    61,
		Message: "The transaction version is not supported",
	}
	ErrUnsupportedContractClassVersion = &RPCError{
		Code:    62,
		Message: "The contract class version is not supported",
	}
	ErrUnexpectedError = &RPCError{
		Code:    63,
		Message: "An unexpected error occurred",
	}
)
