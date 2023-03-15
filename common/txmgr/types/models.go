package types

type TxmErrorType int

// Generalized transaction manager error types that dictates what should be the next action, depending on the RPC error response.
const (
	Successful TxmErrorType = iota
	Fatal
	Retryable
	Underpriced
	Unknown
	Unsupported
	SuccessfulMissingReceipt
	InsufficientFunds
)
