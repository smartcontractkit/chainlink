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

const (
	NonceTooLow = iota
	ReplacementTransactionUnderpriced
	LimitReached
	TransactionAlreadyInMempool
	TerminallyUnderpriced
	InsufficientEth
	TooExpensive
	FeeTooLow
	FeeTooHigh
	Fatal
)

type ClientErrors = map[int]*regexp.Regexp

// Parity
// See: https://github.com/openethereum/openethereum/blob/master/rpc/src/v1/helpers/errors.rs#L420
var parFatal = regexp.MustCompile(`^Transaction gas is too low. There is not enough gas to cover minimal cost of the transaction|^Transaction cost exceeds current gas limit. Limit:|^Invalid signature|Recipient is banned in local queue.|Supplied gas is beyond limit|Sender is banned in local queue|Code is banned in local queue|Transaction is not permitted|Transaction is too big, see chain specification for the limit|^Invalid RLP data`)
var parity = ClientErrors{
	NonceTooLow:                       regexp.MustCompile("^Transaction nonce is too low. Try incrementing the nonce."),
	ReplacementTransactionUnderpriced: regexp.MustCompile("^Transaction gas price .+is too low. There is another transaction with same nonce in the queue"),
	LimitReached:                      regexp.MustCompile("There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."),
	TransactionAlreadyInMempool:       regexp.MustCompile("Transaction with the same hash was already imported."),
	TerminallyUnderpriced:             regexp.MustCompile("^Transaction gas price is too low. It does not satisfy your node's minimal gas price"),
	InsufficientEth:                   regexp.MustCompile("^(Insufficient funds. The account you tried to send transaction from does not have enough funds.|Insufficient balance for transaction.)"),
	Fatal:                             parFatal,
}

// Geth
// See: https://github.com/ethereum/go-ethereum/blob/b9df7ecdc3d3685180ceb29665bab59e9f614da5/core/tx_pool.go#L516
var gethFatal = regexp.MustCompile(`(: |^)(exceeds block gas limit|invalid sender|negative value|oversized data|gas uint64 overflow|intrinsic gas too low|nonce too high)$`)
var geth = ClientErrors{
	NonceTooLow:                       regexp.MustCompile(`(: |^)nonce too low$`),
	ReplacementTransactionUnderpriced: regexp.MustCompile(`(: |^)replacement transaction underpriced$`),
	TransactionAlreadyInMempool:       regexp.MustCompile(`(: |^)(?i)(known transaction|already known)`),
	TerminallyUnderpriced:             regexp.MustCompile(`(: |^)transaction underpriced$`),
	InsufficientEth:                   regexp.MustCompile(`(: |^)(insufficient funds for transfer|insufficient funds for gas \* price \+ value|insufficient balance for transfer)$`),
	TooExpensive:                      regexp.MustCompile(`(: |^)tx fee \([0-9\.]+ ether\) exceeds the configured cap \([0-9\.]+ ether\)$`),
	Fatal:                             gethFatal,
}

// Arbitrum
// https://github.com/OffchainLabs/arbitrum/blob/cac30586bc10ecc1ae73e93de517c90984677fdb/packages/arb-evm/evm/result.go#L158
var arbitrumFatal = regexp.MustCompile(`(: |^)(invalid message format|forbidden sender address|execution reverted: error code)$`)
var arbitrum = ClientErrors{
	// TODO: Arbitrum returns this in case of low or high nonce. Update this when Arbitrum fix it
	NonceTooLow: regexp.MustCompile(`(: |^)invalid transaction nonce$`),
	// TODO: Is it terminally or replacement?
	TerminallyUnderpriced: regexp.MustCompile(`(: |^)gas price too low$`),
	InsufficientEth:       regexp.MustCompile(`(: |^)not enough funds for gas`),
	Fatal:                 arbitrumFatal,
}

var optimism = ClientErrors{
	FeeTooLow:  regexp.MustCompile(`(: |^)fee too low: \d+, use at least tx.gasLimit = \d+ and tx.gasPrice = \d+$`),
	FeeTooHigh: regexp.MustCompile(`(: |^)fee too large: \d+$`),
}

var clients = []ClientErrors{parity, geth, arbitrum, optimism}

func (s *SendError) is(errorType int) bool {
	if s == nil || s.err == nil {
		return false
	}
	str := s.CauseStr()
	for _, client := range clients {
		if _, ok := client[errorType]; !ok {
			continue
		}
		if client[errorType].MatchString(str) {
			return true
		}
	}
	return false
}

var hexDataRegex = regexp.MustCompile(`0x\w+$`)

// IsReplacementUnderpriced indicates that a transaction already exists in the mempool with this nonce but a different gas price or payload
func (s *SendError) IsReplacementUnderpriced() bool {
	return s.is(ReplacementTransactionUnderpriced)
}

func (s *SendError) IsNonceTooLowError() bool {
	return s.is(NonceTooLow)
}

// Geth/parity returns this error if the transaction is already in the node's mempool
func (s *SendError) IsTransactionAlreadyInMempool() bool {
	return s.is(TransactionAlreadyInMempool)
}

// IsTerminallyUnderpriced indicates that this transaction is so far underpriced the node won't even accept it in the first place
func (s *SendError) IsTerminallyUnderpriced() bool {
	return s.is(TerminallyUnderpriced)
}

func (s *SendError) IsTemporarilyUnderpriced() bool {
	return s.is(LimitReached)
}

func (s *SendError) IsInsufficientEth() bool {
	return s.is(InsufficientEth)
}

// IsTooExpensive returns true if the transaction and gas price are combined in
// some way that makes the total transaction too expensive for the eth node to
// accept at all. No amount of retrying at this or higher gas prices can ever
// succeed.
func (s *SendError) IsTooExpensive() bool {
	return s.is(TooExpensive)
}

// IsFeeTooLow is an optimism-specific error returned when total fee is too low
func (s *SendError) IsFeeTooLow() bool {
	return s.is(FeeTooLow)
}

// IsFeeTooHigh is an optimism-specific error returned when total fee is too high
func (s *SendError) IsFeeTooHigh() bool {
	return s.is(FeeTooHigh)
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
	for _, client := range clients {
		if _, ok := client[Fatal]; !ok {
			continue
		}
		if client[Fatal].MatchString(str) {
			return true
		}
	}
	return false
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
