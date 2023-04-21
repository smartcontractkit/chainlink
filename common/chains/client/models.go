package client

type SendTxReturnCode int

// SendTxReturnCode is a generalized client error that dictates what should be the next action, depending on the RPC error response.
const (
	Successful              SendTxReturnCode = iota + 1
	Fatal                                    // Unrecoverable error. Most likely the attempt should be thrown away.
	Retryable                                // The error returned by the RPC indicates that if we retry with the same attempt, the tx will eventually go through.
	Underpriced                              // Attempt was underpriced. New estimation is needed with bumped gas price.
	Unknown                                  // Tx failed with an error response that is not recognized by the client.
	Unsupported                              // Attempt failed with an error response that is not supported by the client for the given chain.
	TransactionAlreadyKnown                  // The transaction that was sent has already been received by the RPC.
	InsufficientFunds                        // Tx was rejected due to insufficient funds.
	ExceedsMaxFee                            // Attempt's fee was higher than the node's limit and got rejected.
	FeeOutOfValidRange                       // This error is returned when we use a fee price suggested from an RPC, but the network rejects the attempt due to an invalid range(mostly used by L2 chains). Retry by requesting a new suggested fee price.
)
