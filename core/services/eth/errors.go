package eth

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
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

// CauseStr returns the string of the original error
func (s *SendError) CauseStr() string {
	if s.err != nil {
		return errors.Cause(s.err).Error()
	}
	return ""
}

// Parity errors
var (
	// Non-fatal
	parTooCheapToReplace    = regexp.MustCompile("^Transaction gas price .+is too low. There is another transaction with same nonce in the queue")
	parLimitReached         = "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
	parAlreadyImported      = "Transaction with the same hash was already imported."
	parNonceTooLow          = "Transaction nonce is too low. Try incrementing the nonce."
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

// Geth and geth-compatible errors
var (
	gethNonceTooLow                       = regexp.MustCompile(`(: |^)nonce too low$`)
	gethReplacementTransactionUnderpriced = regexp.MustCompile(`(: |^)replacement transaction underpriced$`)
	gethKnownTransaction                  = regexp.MustCompile(`(: |^)(?i)(known transaction|already known)`)
	gethTransactionUnderpriced            = regexp.MustCompile(`(: |^)transaction underpriced$`)
	gethInsufficientEth                   = regexp.MustCompile(`(: |^)(insufficient funds for transfer|insufficient funds for gas \* price \+ value|insufficient balance for transfer)$`)
	gethTxFeeExceedsCap                   = regexp.MustCompile(`(: |^)tx fee \([0-9\.]+ ether\) exceeds the configured cap \([0-9\.]+ ether\)$`)

	// Fatal Errors
	// See: https://github.com/ethereum/go-ethereum/blob/b9df7ecdc3d3685180ceb29665bab59e9f614da5/core/tx_pool.go#L516
	gethFatal = regexp.MustCompile(`(: |^)(exceeds block gas limit|invalid sender|negative value|oversized data|gas uint64 overflow|intrinsic gas too low|nonce too high)$`)
)

var hexDataRegex = regexp.MustCompile(`0x\w+$`)

// IsReplacementUnderpriced indicates that a transaction already exists in the mempool with this nonce but a different gas price or payload
func (s *SendError) IsReplacementUnderpriced() bool {
	if s == nil || s.err == nil {
		return false
	}

	str := s.CauseStr()

	switch {
	case gethReplacementTransactionUnderpriced.MatchString(str):
		return true
	case parTooCheapToReplace.MatchString(str):
		return true
	default:
		return false
	}
}

func (s *SendError) IsNonceTooLowError() bool {
	if s == nil || s.err == nil {
		return false
	}

	str := s.CauseStr()
	switch {
	case gethNonceTooLow.MatchString(str):
		return true
	case str == parNonceTooLow:
		return true
	default:
		return false
	}
}

// Geth/parity returns this error if the transaction is already in the node's mempool
func (s *SendError) IsTransactionAlreadyInMempool() bool {
	if s == nil || s.err == nil {
		return false
	}

	str := s.CauseStr()
	switch {
	case gethKnownTransaction.MatchString(str):
		return true
	case str == parAlreadyImported:
		return true
	default:
		return false
	}
}

// IsTerminallyUnderpriced indicates that this transaction is so far
// underpriced the node won't even accept it in the first place
func (s *SendError) IsTerminallyUnderpriced() bool {
	if s == nil || s.err == nil {
		return false
	}

	str := s.CauseStr()
	switch {
	case gethTransactionUnderpriced.MatchString(str):
		return true
	case parInsufficientGasPrice.MatchString(str):
		return true
	default:
		return false
	}
}

func (s *SendError) IsTemporarilyUnderpriced() bool {
	return s != nil && s.err != nil && s.CauseStr() == parLimitReached
}

func (s *SendError) IsInsufficientEth() bool {
	if s == nil || s.err == nil {
		return false
	}

	str := s.CauseStr()
	switch {
	case gethInsufficientEth.MatchString(str):
		return true
	case parInsufficientEth.MatchString(str):
		return true
	default:
		return false
	}
}

// IsTooExpensive returns true if the transaction and gas price are combined in
// some way that makes the total transaction too expensive for the eth node to
// accept at all. No amount of retrying at this or higher gas prices can ever
// succeed.
func (s *SendError) IsTooExpensive() bool {
	if s == nil || s.err == nil {
		return false
	}

	str := s.CauseStr()

	return gethTxFeeExceedsCap.MatchString(str)
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
// 1. It will never be included into a block as a result of this send
// 2. Resending the transaction at a different gas price will never change the outcome
func isFatalSendError(err error) bool {
	if err == nil {
		return false
	}
	str := errors.Cause(err).Error()
	return isGethFatal(str) || isParityFatal(str)
}

func isGethFatal(s string) bool {
	return gethFatal.MatchString(s)
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

// go-ethereum@v1.10.0/rpc/json.go
type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *JsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

// ExtractRevertReasonFromRPCError attempts to extract the revert reason from the response of
// an RPC eth_call that reverted by parsing the message from the "data" field
// ex:
// kovan (parity)
// { "error": { "code" : -32015, "data": "Reverted 0xABC123...", "message": "VM execution error." } } // revert reason always omitted
// rinkeby / ropsten (geth)
// { "error":  { "code": 3, "data": "0x0xABC123...", "message": "execution reverted: hello world" } } // revert reason included in message
func ExtractRevertReasonFromRPCError(err error) (string, error) {
	if err == nil {
		return "", errors.New("no error present")
	}
	cause := errors.Cause(err)
	jsonBytes, err := json.Marshal(cause)
	if err != nil {
		return "", errors.Wrap(err, "unable to marshal err to json")
	}
	jErr := JsonError{}
	err = json.Unmarshal(jsonBytes, &jErr)
	if err != nil {
		return "", errors.Wrap(err, "unable to unmarshal json into jsonError struct")
	}
	dataStr, ok := jErr.Data.(string)
	if !ok {
		return "", errors.New("invalid error type")
	}
	matches := hexDataRegex.FindStringSubmatch(dataStr)
	if len(matches) != 1 {
		return "", errors.New("unknown data payload format")
	}
	hexData := utils.RemoveHexPrefix(matches[0])
	if len(hexData) < 8 {
		return "", errors.New("unknown data payload format")
	}
	revertReasonBytes, err := hex.DecodeString(hexData[8:])
	if err != nil {
		return "", errors.Wrap(err, "unable to decode hex to bytes")
	}
	revertReasonBytes = bytes.Trim(revertReasonBytes, "\x00")
	revertReason := strings.TrimSpace(string(revertReasonBytes))
	return revertReason, nil
}
