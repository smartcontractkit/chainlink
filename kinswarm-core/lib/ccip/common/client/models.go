package client

import (
	"bytes"
	"fmt"
)

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
	TerminallyStuck                          // The error returned when a transaction is or could get terminally stuck in the mempool without any chance of inclusion.
	sendTxReturnCodeLen                      // tracks the number of errors. Must always be last
)

// sendTxSevereErrors - error codes which signal that transaction would never be accepted in its current form by the node
var sendTxSevereErrors = []SendTxReturnCode{Fatal, Underpriced, Unsupported, ExceedsMaxFee, FeeOutOfValidRange, Unknown}

// sendTxSuccessfulCodes - error codes which signal that transaction was accepted by the node
var sendTxSuccessfulCodes = []SendTxReturnCode{Successful, TransactionAlreadyKnown}

func (c SendTxReturnCode) String() string {
	switch c {
	case Successful:
		return "Successful"
	case Fatal:
		return "Fatal"
	case Retryable:
		return "Retryable"
	case Underpriced:
		return "Underpriced"
	case Unknown:
		return "Unknown"
	case Unsupported:
		return "Unsupported"
	case TransactionAlreadyKnown:
		return "TransactionAlreadyKnown"
	case InsufficientFunds:
		return "InsufficientFunds"
	case ExceedsMaxFee:
		return "ExceedsMaxFee"
	case FeeOutOfValidRange:
		return "FeeOutOfValidRange"
	case TerminallyStuck:
		return "TerminallyStuck"
	default:
		return fmt.Sprintf("SendTxReturnCode(%d)", c)
	}
}

type NodeTier int

const (
	Primary = NodeTier(iota)
	Secondary
)

func (n NodeTier) String() string {
	switch n {
	case Primary:
		return "primary"
	case Secondary:
		return "secondary"
	default:
		return fmt.Sprintf("NodeTier(%d)", n)
	}
}

// syncStatus - defines problems related to RPC's state synchronization. Can be used as a bitmask to define multiple issues
type syncStatus int

const (
	// syncStatusSynced - RPC is fully synced
	syncStatusSynced = 0
	// syncStatusNotInSyncWithPool - RPC is lagging behind the highest block observed within the pool of RPCs
	syncStatusNotInSyncWithPool syncStatus = 1 << iota
	// syncStatusNoNewHead - RPC failed to produce a new head for too long
	syncStatusNoNewHead
	// syncStatusNoNewFinalizedHead - RPC failed to produce a new finalized head for too long
	syncStatusNoNewFinalizedHead
	syncStatusLen
)

func (s syncStatus) String() string {
	if s == syncStatusSynced {
		return "Synced"
	}
	var result bytes.Buffer
	for i := syncStatusNotInSyncWithPool; i < syncStatusLen; i = i << 1 {
		if i&s == 0 {
			continue
		}
		result.WriteString(i.string())
		result.WriteString(",")
	}
	result.Truncate(result.Len() - 1)
	return result.String()
}

func (s syncStatus) string() string {
	switch s {
	case syncStatusNotInSyncWithPool:
		return "NotInSyncWithRPCPool"
	case syncStatusNoNewHead:
		return "NoNewHead"
	case syncStatusNoNewFinalizedHead:
		return "NoNewFinalizedHead"
	default:
		return fmt.Sprintf("syncStatus(%d)", s)
	}
}
