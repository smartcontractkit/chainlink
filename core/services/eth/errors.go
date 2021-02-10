package eth

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// fatal means this transaction can never be accepted even with a different nonce or higher gas price
type SendError struct {
	fatal bool
	err   error
}

func (s *SendError) Error() string {
	return s.err.Error()
}

func (s *SendError) StrPtr() *string {
	e := s.err.Error()
	return &e
}

// Fatal indicates whether the error should be considered fatal or not
// Fatal errors mean that no matter how many times the send is retried, no node
// will ever accept it
func (s *SendError) Fatal() bool {
	return s != nil && s.fatal
}

// Parity errors
var (
	// Non-fatal
	parTooCheapToReplace    = regexp.MustCompile("^Transaction gas price .+is too low. There is another transaction with same nonce in the queue")
	parLimitReached         = "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
	parAlreadyImported      = "Transaction with the same hash was already imported."
	parOld                  = "Transaction nonce is too low. Try incrementing the nonce."
	parInsufficientGasPrice = regexp.MustCompile("^Transaction gas price is too low. It does not satisfy your node's minimal gas price")
	parInsufficientEth      = regexp.MustCompile("^(Insufficient funds. The account you tried to send transaction from does not have enough funds.|Insufficient balance for transaction.)")

	// Fatal
	parInsufficientGas  = regexp.MustCompile("^Transaction gas is too low. There is not enough gas to cover minimal cost of the transaction")
	parGasLimitExceeded = regexp.MustCompile("^Transaction cost exceeds current gas limit. Limit:")
	parInvalidSignature = regexp.MustCompile("^Invalid signature")
	parInvalidGasLimit  = "Supplied gas is beyond limit."
	parSenderBanned     = "Sender is banned in local queue."
	parRecipientBanned  = "Recipient is banned in local queue."
	parCodeBanned       = "Code is banned in local queue."
	parNotAllowed       = "Transaction is not permitted."
	parTooBig           = "Transaction is too big, see chain specification for the limit."
	parInvalidRlp       = regexp.MustCompile("^Invalid RLP data:")
)

// IsReplacementUnderpriced indicates that a transaction already exists in the mempool with this nonce but a different gas price or payload
func (s *SendError) IsReplacementUnderpriced() bool {
	return s != nil && s.err != nil && (s.Error() == "replacement transaction underpriced" || parTooCheapToReplace.MatchString(s.Error()))
}

func (s *SendError) IsNonceTooLowError() bool {
	return s != nil && s.err != nil && ((s.Error() == "nonce too low") || s.Error() == parOld)
}

// Geth/parity returns this error if the transaction is already in the node's mempool
func (s *SendError) IsTransactionAlreadyInMempool() bool {
	return s != nil && s.err != nil && (strings.HasPrefix(strings.ToLower(s.Error()), "known transaction") || s.Error() == "already known" || s.Error() == parAlreadyImported)
}

// IsTerminallyUnderpriced indicates that this transaction is so far
// underpriced the node won't even accept it in the first place
func (s *SendError) IsTerminallyUnderpriced() bool {
	return s != nil && s.err != nil && (s.Error() == "transaction underpriced" || parInsufficientGasPrice.MatchString(s.Error()))
}

func (s *SendError) IsTemporarilyUnderpriced() bool {
	return s != nil && s.err != nil && s.Error() == parLimitReached
}

func (s *SendError) IsInsufficientEth() bool {
	if s == nil || s.err == nil {
		return false
	}
	switch s.Error() {
	// Geth
	case "insufficient funds for transfer", "insufficient funds for gas * price + value", "insufficient balance for transfer":
		return true
	default:
		// Parity
		return parInsufficientEth.MatchString(s.Error())
	}
}

func NewFatalSendErrorS(s string) *SendError {
	return &SendError{err: errors.New(s), fatal: true}
}

func NewFatalSendError(e error) *SendError {
	if e == nil {
		return nil
	}
	return &SendError{err: errors.WithStack(e), fatal: true}
}

func NewSendErrorS(s string) *SendError {
	return NewSendError(errors.New(s))
}

func NewSendError(e error) *SendError {
	if e == nil {
		return nil
	}
	fatal := isFatalSendError(e)
	return &SendError{err: errors.WithStack(e), fatal: fatal}
}

// Geth/parity returns these errors if the transaction failed in such a way that:
// 1. It can NEVER be included into a block
// 2. Resending the transaction even with higher gas price will never change that outcome
func isFatalSendError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	// Geth errors
	// See: https://github.com/ethereum/go-ethereum/blob/b9df7ecdc3d3685180ceb29665bab59e9f614da5/core/tx_pool.go#L516
	case "exceeds block gas limit", "invalid sender", "negative value", "oversized data", "gas uint64 overflow", "intrinsic gas too low", "nonce too high":
		return true
	default:
		return isParityFatal(err.Error())
	}
}

// See: https://github.com/openethereum/openethereum/blob/master/rpc/src/v1/helpers/errors.rs#L420
func isParityFatal(s string) bool {
	return s == parInvalidGasLimit ||
		s == parSenderBanned ||
		s == parRecipientBanned ||
		s == parCodeBanned ||
		s == parNotAllowed ||
		s == parTooBig ||
		(parInsufficientGas.MatchString(s) ||
			parGasLimitExceeded.MatchString(s) ||
			parInvalidSignature.MatchString(s) ||
			parInvalidRlp.MatchString(s))
}
