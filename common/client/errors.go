package client

// Provides error classification to external components in a chain agnostic way
// Only exposes the error types that could be set in the transaction error field
type TxError interface {
	error
	IsFatal() bool
	IsTerminallyStuck() bool
}
