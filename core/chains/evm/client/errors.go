package client

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
)

// fatal means this transaction can never be accepted even with a different nonce or higher gas price
type SendError struct {
	fatal bool
	err   error
}

func (s *SendError) Error() string {
	return s.err.Error()
}

// Fatal indicates whether the error should be considered fatal or not
// Fatal errors mean that no matter how many times the send is retried, no node
// will ever accept it
func (s *SendError) Fatal(configErrors *ClientErrors) bool {
	if configErrors != nil && configErrors.ErrIs(s.err, Fatal) {
		return true
	}
	return s != nil && s.fatal
}

const (
	NonceTooLow = iota
	// Nethermind specific error. Nethermind throws a NonceGap error when the tx nonce is greater than current_nonce + tx_count_in_mempool, instead of keeping the tx in mempool.
	// See: https://github.com/NethermindEth/nethermind/blob/master/src/Nethermind/Nethermind.TxPool/Filters/GapNonceFilter.cs
	NonceTooHigh
	ReplacementTransactionUnderpriced
	LimitReached
	TransactionAlreadyInMempool
	TerminallyUnderpriced
	InsufficientEth
	TxFeeExceedsCap
	// Note: L2FeeTooLow/L2FeeTooHigh/L2Full have a very specific meaning specific
	// to L2s (Arbitrum and clones). Do not implement this for non-L2
	// chains. This is potentially confusing because some RPC nodes e.g.
	// Nethermind implement an error called `FeeTooLow` which has distinct
	// meaning from this one.
	L2FeeTooLow
	L2FeeTooHigh
	L2Full
	TransactionAlreadyMined
	Fatal
	ServiceUnavailable
	OutOfCounters
)

type ClientErrors map[int]*regexp.Regexp

// ErrIs returns true if err matches any provided error types
func (e *ClientErrors) ErrIs(err error, errorTypes ...int) bool {
	if err == nil {
		return false
	}
	for _, errorType := range errorTypes {
		if _, ok := (*e)[errorType]; !ok {
			return false
		}
		if (*e)[errorType].String() == "" {
			return false
		}
		if (*e)[errorType].MatchString(pkgerrors.Cause(err).Error()) {
			return true
		}
	}
	return false
}

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
var gethFatal = regexp.MustCompile(`(: |^)(exceeds block gas limit|invalid sender|negative value|oversized data|gas uint64 overflow|intrinsic gas too low)$`)
var geth = ClientErrors{
	NonceTooLow:                       regexp.MustCompile(`(: |^)nonce too low$`),
	NonceTooHigh:                      regexp.MustCompile(`(: |^)nonce too high$`),
	ReplacementTransactionUnderpriced: regexp.MustCompile(`(: |^)replacement transaction underpriced$`),
	TransactionAlreadyInMempool:       regexp.MustCompile(`(: |^)(?i)(known transaction|already known)`),
	TerminallyUnderpriced:             regexp.MustCompile(`(: |^)transaction underpriced$`),
	InsufficientEth:                   regexp.MustCompile(`(: |^)(insufficient funds for transfer|insufficient funds for gas \* price \+ value|insufficient balance for transfer)$`),
	TxFeeExceedsCap:                   regexp.MustCompile(`(: |^)tx fee \([0-9\.]+ [a-zA-Z]+\) exceeds the configured cap \([0-9\.]+ [a-zA-Z]+\)$`),
	Fatal:                             gethFatal,
}

// Besu
// See: https://github.com/hyperledger/besu/blob/81f25e15f9891787829b532f2fb38c8c43fd6b2e/ethereum/api/src/main/java/org/hyperledger/besu/ethereum/api/jsonrpc/internal/response/JsonRpcError.java
var besuFatal = regexp.MustCompile(`^(Intrinsic gas exceeds gas limit|Transaction gas limit exceeds block gas limit|Invalid signature)$`)
var besu = ClientErrors{
	NonceTooLow:                       regexp.MustCompile(`^Nonce too low$`),
	ReplacementTransactionUnderpriced: regexp.MustCompile(`^Replacement transaction underpriced$`),
	TransactionAlreadyInMempool:       regexp.MustCompile(`^Known transaction$`),
	TerminallyUnderpriced:             regexp.MustCompile(`^Gas price below configured minimum gas price$`),
	InsufficientEth:                   regexp.MustCompile(`^Upfront cost exceeds account balance$`),
	TxFeeExceedsCap:                   regexp.MustCompile(`^Transaction fee cap exceeded$`),
	Fatal:                             besuFatal,
}

// Erigon
// See:
//   - https://github.com/ledgerwatch/erigon/blob/devel/core/tx_pool.go
//   - https://github.com/ledgerwatch/erigon/blob/devel/core/error.go
//   - https://github.com/ledgerwatch/erigon/blob/devel/core/vm/errors.go
//
// Note: some error definitions are unused, many errors are created inline.
var erigonFatal = regexp.MustCompile(`(: |^)(exceeds block gas limit|invalid sender|negative value|oversized data|gas uint64 overflow|intrinsic gas too low)$`)
var erigon = ClientErrors{
	NonceTooLow:                       regexp.MustCompile(`(: |^)nonce too low$`),
	NonceTooHigh:                      regexp.MustCompile(`(: |^)nonce too high$`),
	ReplacementTransactionUnderpriced: regexp.MustCompile(`(: |^)replacement transaction underpriced$`),
	TransactionAlreadyInMempool:       regexp.MustCompile(`(: |^)(block already known|already known)`),
	TerminallyUnderpriced:             regexp.MustCompile(`(: |^)transaction underpriced$`),
	InsufficientEth:                   regexp.MustCompile(`(: |^)(insufficient funds for transfer|insufficient funds for gas \* price \+ value|insufficient balance for transfer)$`),
	TxFeeExceedsCap:                   regexp.MustCompile(`(: |^)tx fee \([0-9\.]+ [a-zA-Z]+\) exceeds the configured cap \([0-9\.]+ [a-zA-Z]+\)$`),
	Fatal:                             erigonFatal,
}

// Arbitrum
// https://github.com/OffchainLabs/arbitrum/blob/cac30586bc10ecc1ae73e93de517c90984677fdb/packages/arb-evm/evm/result.go#L158
// nitro: https://github.com/OffchainLabs/go-ethereum/blob/master/core/state_transition.go
var arbitrumFatal = regexp.MustCompile(`(: |^)(invalid message format|forbidden sender address)$|(: |^)(execution reverted)(:|$)`)
var arbitrum = ClientErrors{
	// TODO: Arbitrum returns this in case of low or high nonce. Update this when Arbitrum fix it
	// https://app.shortcut.com/chainlinklabs/story/16801/add-full-support-for-incorrect-nonce-on-arbitrum
	NonceTooLow:           regexp.MustCompile(`(: |^)invalid transaction nonce$|(: |^)nonce too low(:|$)`),
	NonceTooHigh:          regexp.MustCompile(`(: |^)nonce too high(:|$)`),
	TerminallyUnderpriced: regexp.MustCompile(`(: |^)gas price too low$`),
	InsufficientEth:       regexp.MustCompile(`(: |^)(not enough funds for gas|insufficient funds for gas \* price \+ value)`),
	Fatal:                 arbitrumFatal,
	L2FeeTooLow:           regexp.MustCompile(`(: |^)max fee per gas less than block base fee(:|$)`),
	L2Full:                regexp.MustCompile(`(: |^)(queue full|sequencer pending tx pool full, please try again)(:|$)`),
	ServiceUnavailable:    regexp.MustCompile(`(: |^)502 Bad Gateway: [\s\S]*$`),
}

var celo = ClientErrors{
	TxFeeExceedsCap:       regexp.MustCompile(`(: |^)tx fee \([0-9\.]+ of currency celo\) exceeds the configured cap \([0-9\.]+ [a-zA-Z]+\)$`),
	TerminallyUnderpriced: regexp.MustCompile(`(: |^)gasprice is less than gas price minimum floor`),
	InsufficientEth:       regexp.MustCompile(`(: |^)insufficient funds for gas \* price \+ value \+ gatewayFee$`),
	LimitReached:          regexp.MustCompile(`(: |^)txpool is full`),
}

var metis = ClientErrors{
	L2FeeTooLow: regexp.MustCompile(`(: |^)gas price too low: \d+ wei, use at least tx.gasPrice = \d+ wei$`),
}

// Substrate (Moonriver)
var substrate = ClientErrors{
	NonceTooLow:                 regexp.MustCompile(`(: |^)Pool\(Stale\)$`),
	TransactionAlreadyInMempool: regexp.MustCompile(`(: |^)Pool\(AlreadyImported\)$`),
}

var avalanche = ClientErrors{
	NonceTooLow: regexp.MustCompile(`(: |^)nonce too low: address 0x[0-9a-fA-F]{40} current nonce \([\d]+\) > tx nonce \([\d]+\)$`),
}

// Klaytn
// https://github.com/klaytn/klaytn/blob/dev/blockchain/error.go
// https://github.com/klaytn/klaytn/blob/dev/blockchain/tx_pool.go
var klaytn = ClientErrors{
	NonceTooLow:                       regexp.MustCompile(`(: |^)nonce too low$`),                                                                                    // retry with an increased nonce
	TransactionAlreadyInMempool:       regexp.MustCompile(`(: |^)(known transaction)`),                                                                               // don't send the tx again. The exactly same tx is already in the mempool
	ReplacementTransactionUnderpriced: regexp.MustCompile(`(: |^)replacement transaction underpriced$|there is another tx which has the same nonce in the tx pool$`), // retry with an increased gasPrice or maxFeePerGas. This error happened when there is another tx having higher gasPrice or maxFeePerGas exist in the mempool
	TerminallyUnderpriced:             regexp.MustCompile(`(: |^)(transaction underpriced|^intrinsic gas too low)`),                                                  // retry with an increased gasPrice or maxFeePerGas
	LimitReached:                      regexp.MustCompile(`(: |^)txpool is full`),                                                                                    // retry with few seconds wait
	InsufficientEth:                   regexp.MustCompile(`(: |^)insufficient funds`),                                                                                // stop to send a tx. The sender address doesn't have enough KLAY
	TxFeeExceedsCap:                   regexp.MustCompile(`(: |^)(invalid gas fee cap|max fee per gas higher than max priority fee per gas)`),                        // retry with a valid gasPrice, maxFeePerGas, or maxPriorityFeePerGas. The new value can get from the return of `eth_gasPrice`
	Fatal:                             gethFatal,
}

// Nethermind
// All errors: https://github.com/NethermindEth/nethermind/blob/master/src/Nethermind/Nethermind.TxPool/AcceptTxResult.cs
// All filters: https://github.com/NethermindEth/nethermind/tree/9b68ec048c65f4b44fb863164c0dec3f7780d820/src/Nethermind/Nethermind.TxPool/Filters
var nethermindFatal = regexp.MustCompile(`(: |^)(SenderIsContract|Invalid(, transaction Hash is null)?|Int256Overflow|FailedToResolveSender|GasLimitExceeded(, Gas limit: \d+, gas limit of rejected tx: \d+)?)$`)
var nethermind = ClientErrors{
	// OldNonce: The EOA (externally owned account) that signed this transaction (sender) has already signed and executed a transaction with the same nonce.
	NonceTooLow:  regexp.MustCompile(`(: |^)OldNonce(, Current nonce: \d+, nonce of rejected tx: \d+)?$`),
	NonceTooHigh: regexp.MustCompile(`(: |^)NonceGap(, Future nonce. Expected nonce: \d+)?$`),

	// FeeTooLow/FeeTooLowToCompete: Fee paid by this transaction is not enough to be accepted in the mempool.
	TerminallyUnderpriced: regexp.MustCompile(`(: |^)(FeeTooLow(, MaxFeePerGas too low. MaxFeePerGas: \d+, BaseFee: \d+, MaxPriorityFeePerGas:\d+, Block number: \d+|` +
		`, EffectivePriorityFeePerGas too low \d+ < \d+, BaseFee: \d+|` +
		`, FeePerGas needs to be higher than \d+ to be added to the TxPool. Affordable FeePerGas of rejected tx: \d+.)?|` +
		`FeeTooLowToCompete)$`),

	// AlreadyKnown: A transaction with the same hash has already been added to the pool in the past.
	// OwnNonceAlreadyUsed: A transaction with same nonce has been signed locally already and is awaiting in the pool.
	TransactionAlreadyInMempool: regexp.MustCompile(`(: |^)(AlreadyKnown|OwnNonceAlreadyUsed)$`),

	// InsufficientFunds: Sender account has not enough balance to execute this transaction.
	InsufficientEth:    regexp.MustCompile(`(: |^)InsufficientFunds(, Account balance: \d+, cumulative cost: \d+|, Balance is \d+ less than sending value \+ gas \d+)?$`),
	ServiceUnavailable: regexp.MustCompile(`(: |^)503 Service Unavailable: [\s\S]*$`),
	Fatal:              nethermindFatal,
}

// Harmony
// https://github.com/harmony-one/harmony/blob/main/core/tx_pool.go#L49
var harmonyFatal = regexp.MustCompile("(: |^)(invalid shard|staking message does not match directive message|`from` address of transaction in blacklist|`to` address of transaction in blacklist)$")
var harmony = ClientErrors{
	TransactionAlreadyMined: regexp.MustCompile(`(: |^)transaction already finalized$`),
	Fatal:                   harmonyFatal,
}

var zkSync = ClientErrors{
	NonceTooLow:           regexp.MustCompile(`(?:: |^)nonce too low\..+actual: \d*$`),
	NonceTooHigh:          regexp.MustCompile(`(?:: |^)nonce too high\..+actual: \d*$`),
	TerminallyUnderpriced: regexp.MustCompile(`(?:: |^)(max fee per gas less than block base fee|virtual machine entered unexpected state. please contact developers and provide transaction details that caused this error. Error description: The operator included transaction with an unacceptable gas price)$`),
	InsufficientEth:       regexp.MustCompile(`(?:: |^)(?:insufficient balance for transfer$|insufficient funds for gas + value)`),
	TxFeeExceedsCap:       regexp.MustCompile(`(?:: |^)max priority fee per gas higher than max fee per gas$`),
	// intrinsic gas too low 						- gas limit less than 14700
	// Not enough gas for transaction validation 	- gas limit less than L2 fee
	// Failed to pay the fee to the operator 		- gas limit less than L2+L1 fee
	// Error function_selector = 0x, data = 0x 		- contract call with gas limit of 0
	// can't start a transaction from a non-account - trying to send from an invalid address, e.g. estimating a contract -> contract tx
	// max fee per gas higher than 2^64-1 			- uint64 overflow
	// oversized data 								- data too large
	Fatal:                       regexp.MustCompile(`(?:: |^)(?:exceeds block gas limit|intrinsic gas too low|Not enough gas for transaction validation|Failed to pay the fee to the operator|Error function_selector = 0x, data = 0x|invalid sender. can't start a transaction from a non-account|max(?: priority)? fee per (?:gas|pubdata byte) higher than 2\^64-1|oversized data. max: \d+; actual: \d+)$`),
	TransactionAlreadyInMempool: regexp.MustCompile(`known transaction. transaction with hash .* is already in the system`),
}

var zkEvm = ClientErrors{
	OutOfCounters: regexp.MustCompile(`(?:: |^)not enough .* counters to continue the execution$`),
}

var clients = []ClientErrors{parity, geth, arbitrum, metis, substrate, avalanche, nethermind, harmony, besu, erigon, klaytn, celo, zkSync, zkEvm}

// ClientErrorRegexes returns a map of compiled regexes for each error type
func ClientErrorRegexes(errsRegex config.ClientErrors) *ClientErrors {
	if errsRegex == nil {
		return &ClientErrors{}
	}
	return &ClientErrors{
		NonceTooLow:                       regexp.MustCompile(errsRegex.NonceTooLow()),
		NonceTooHigh:                      regexp.MustCompile(errsRegex.NonceTooHigh()),
		ReplacementTransactionUnderpriced: regexp.MustCompile(errsRegex.ReplacementTransactionUnderpriced()),
		LimitReached:                      regexp.MustCompile(errsRegex.LimitReached()),
		TransactionAlreadyInMempool:       regexp.MustCompile(errsRegex.TransactionAlreadyInMempool()),
		TerminallyUnderpriced:             regexp.MustCompile(errsRegex.TerminallyUnderpriced()),
		InsufficientEth:                   regexp.MustCompile(errsRegex.InsufficientEth()),
		TxFeeExceedsCap:                   regexp.MustCompile(errsRegex.TxFeeExceedsCap()),
		L2FeeTooLow:                       regexp.MustCompile(errsRegex.L2FeeTooLow()),
		L2FeeTooHigh:                      regexp.MustCompile(errsRegex.L2FeeTooHigh()),
		L2Full:                            regexp.MustCompile(errsRegex.L2Full()),
		TransactionAlreadyMined:           regexp.MustCompile(errsRegex.TransactionAlreadyMined()),
		Fatal:                             regexp.MustCompile(errsRegex.Fatal()),
		ServiceUnavailable:                regexp.MustCompile(errsRegex.ServiceUnavailable()),
	}
}

func (s *SendError) is(errorType int, configErrors *ClientErrors) bool {
	if s == nil || s.err == nil {
		return false
	}
	if configErrors != nil && configErrors.ErrIs(s.err, errorType) {
		return true
	}
	for _, client := range clients {
		if client.ErrIs(s.err, errorType) {
			return true
		}
	}
	return false
}

// IsReplacementUnderpriced indicates that a transaction already exists in the mempool with this nonce but a different gas price or payload
func (s *SendError) IsReplacementUnderpriced(configErrors *ClientErrors) bool {
	return s.is(ReplacementTransactionUnderpriced, configErrors)
}

func (s *SendError) IsNonceTooLowError(configErrors *ClientErrors) bool {
	return s.is(NonceTooLow, configErrors)
}

func (s *SendError) IsNonceTooHighError(configErrors *ClientErrors) bool {
	return s.is(NonceTooHigh, configErrors)
}

// IsTransactionAlreadyMined - Harmony returns this error if the transaction has already been mined
func (s *SendError) IsTransactionAlreadyMined(configErrors *ClientErrors) bool {
	return s.is(TransactionAlreadyMined, configErrors)
}

// Geth/parity returns this error if the transaction is already in the node's mempool
func (s *SendError) IsTransactionAlreadyInMempool(configErrors *ClientErrors) bool {
	return s.is(TransactionAlreadyInMempool, configErrors)
}

// IsTerminallyUnderpriced indicates that this transaction is so far underpriced the node won't even accept it in the first place
func (s *SendError) IsTerminallyUnderpriced(configErrors *ClientErrors) bool {
	return s.is(TerminallyUnderpriced, configErrors)
}

func (s *SendError) IsTemporarilyUnderpriced(configErrors *ClientErrors) bool {
	return s.is(LimitReached, configErrors)
}

func (s *SendError) IsInsufficientEth(configErrors *ClientErrors) bool {
	return s.is(InsufficientEth, configErrors)
}

// IsTxFeeExceedsCap returns true if the transaction and gas price are combined in
// some way that makes the total transaction too expensive for the eth node to
// accept at all. No amount of retrying at this or higher gas prices can ever
// succeed.
func (s *SendError) IsTxFeeExceedsCap(configErrors *ClientErrors) bool {
	return s.is(TxFeeExceedsCap, configErrors)
}

// L2FeeTooLow is an l2-specific error returned when total fee is too low
func (s *SendError) L2FeeTooLow(configErrors *ClientErrors) bool {
	return s.is(L2FeeTooLow, configErrors)
}

// IsL2FeeTooHigh is an l2-specific error returned when total fee is too high
func (s *SendError) IsL2FeeTooHigh(configErrors *ClientErrors) bool {
	return s.is(L2FeeTooHigh, configErrors)
}

// IsL2Full is an l2-specific error returned when the queue or mempool is full.
func (s *SendError) IsL2Full(configErrors *ClientErrors) bool {
	return s.is(L2Full, configErrors)
}

// IsServiceUnavailable indicates if the error was caused by a service being unavailable
func (s *SendError) IsServiceUnavailable(configErrors *ClientErrors) bool {
	return s.is(ServiceUnavailable, configErrors)
}

// IsOutOfCounters is a zk chain specific error returned if the transaction is too complex to prove on zk circuits
func (s *SendError) IsOutOfCounters(configErrors *ClientErrors) bool {
	return s.is(OutOfCounters, configErrors)
}

// IsTimeout indicates if the error was caused by an exceeded context deadline
func (s *SendError) IsTimeout() bool {
	if s == nil {
		return false
	}
	if s.err == nil {
		return false
	}
	return pkgerrors.Is(s.err, context.DeadlineExceeded)
}

// IsCanceled indicates if the error was caused by an context cancellation
func (s *SendError) IsCanceled() bool {
	if s == nil {
		return false
	}
	if s.err == nil {
		return false
	}
	return pkgerrors.Is(s.err, context.Canceled)
}

func NewFatalSendError(e error) *SendError {
	if e == nil {
		return nil
	}
	return &SendError{err: pkgerrors.WithStack(e), fatal: true}
}

func NewSendErrorS(s string) *SendError {
	return NewSendError(pkgerrors.New(s))
}

func NewSendError(e error) *SendError {
	if e == nil {
		return nil
	}
	fatal := isFatalSendError(e)
	return &SendError{err: pkgerrors.WithStack(e), fatal: fatal}
}

// Geth/parity returns these errors if the transaction failed in such a way that:
// 1. It will never be included into a block as a result of this send
// 2. Resending the transaction at a different gas price will never change the outcome
func isFatalSendError(err error) bool {
	if err == nil {
		return false
	}
	str := pkgerrors.Cause(err).Error()

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

func (err JsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error { Code = %d, Data = '%v' }", err.Code, err.Data)
	}
	return err.Message
}

func (err *JsonError) String() string {
	return fmt.Sprintf("json-rpc error { Code = %d, Message = '%s', Data = '%v' }", err.Code, err.Message, err.Data)
}

func ExtractRPCErrorOrNil(err error) *JsonError {
	jErr, eErr := ExtractRPCError(err)
	if eErr != nil {
		return nil
	}
	return jErr
}

// ExtractRPCError attempts to extract a full JsonError (including revert reason details)
// from an error returned by a CallContract to an external RPC. As per https://github.com/ethereum/go-ethereum/blob/c49e065fea78a5d3759f7853a608494913e5824e/internal/ethapi/api.go#L974
// CallContract server side for a revert will return an error which contains either:
//   - The error directly from the EVM if there's no data (no revert reason, like an index out of bounds access) which
//     when marshalled will only have a Message.
//   - An error which implements rpc.DataError which when marshalled will have a Data field containing the execution result.
//     If the revert not a custom Error (solidity >= 0.8.0), like require(1 == 2, "revert"), then geth and forks will automatically
//     parse the string and put it in the message. If its a custom error, it's up to the client to decode the Data field which will be
//     the abi encoded data of the custom error, i.e. revert MyCustomError(10) -> keccak(MyCustomError(uint256))[:4] || abi.encode(10).
//
// However, it appears that RPCs marshal this in different ways into a JsonError object received client side,
// some adding "Reverted" prefixes, removing the method signature etc. To avoid RPC specific parsing and support custom errors
// we return the full object returned from the RPC with a String() method that stringifies all fields for logging so no information is lost.
// Some examples:
// kovan (parity)
// { "error": { "code" : -32015, "data": "Reverted 0xABC123...", "message": "VM execution error." } } // revert reason always omitted from message.
// rinkeby / ropsten (geth)
// { "error":  { "code": 3, "data": "0xABC123...", "message": "execution reverted: hello world" } } // revert reason automatically parsed if a simple require and included in message.
func ExtractRPCError(baseErr error) (*JsonError, error) {
	if baseErr == nil {
		return nil, pkgerrors.New("no error present")
	}
	cause := pkgerrors.Cause(baseErr)
	jsonBytes, err := json.Marshal(cause)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to marshal err to json")
	}
	jErr := JsonError{}
	err = json.Unmarshal(jsonBytes, &jErr)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to unmarshal json into jsonError struct (got: %v)", baseErr)
	}
	if jErr.Code == 0 {
		return nil, pkgerrors.Errorf("not a RPCError because it does not have a code (got: %v)", baseErr)
	}
	return &jErr, nil
}

func ClassifySendError(err error, clientErrors config.ClientErrors, lggr logger.SugaredLogger, tx *types.Transaction, fromAddress common.Address, isL2 bool) commonclient.SendTxReturnCode {
	sendError := NewSendError(err)
	if sendError == nil {
		return commonclient.Successful
	}

	configErrors := ClientErrorRegexes(clientErrors)

	if sendError.Fatal(configErrors) {
		lggr.Criticalw("Fatal error sending transaction", "err", sendError, "etx", tx)
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return commonclient.Fatal
	}
	if sendError.IsNonceTooLowError(configErrors) || sendError.IsTransactionAlreadyMined(configErrors) {
		lggr.Debugw(fmt.Sprintf("Transaction already confirmed for this nonce: %d", tx.Nonce()), "err", sendError, "etx", tx)
		// Nonce too low indicated that a transaction at this nonce was confirmed already.
		// Mark it as TransactionAlreadyKnown.
		return commonclient.TransactionAlreadyKnown
	}
	if sendError.IsReplacementUnderpriced(configErrors) {
		lggr.Errorw(fmt.Sprintf("Replacement transaction underpriced for eth_tx %x. "+
			"Please note that using your node's private keys outside of the chainlink node is NOT SUPPORTED and can lead to missed transactions.",
			tx.Hash()), "gasPrice", tx.GasPrice, "gasTipCap", tx.GasTipCap, "gasFeeCap", tx.GasFeeCap, "err", sendError, "etx", tx)

		// Assume success and hand off to the next cycle.
		return commonclient.Successful
	}
	if sendError.IsTransactionAlreadyInMempool(configErrors) {
		lggr.Debugw("Transaction already in mempool", "etx", tx, "err", sendError)
		return commonclient.Successful
	}
	if sendError.IsTemporarilyUnderpriced(configErrors) {
		lggr.Infow("Transaction temporarily underpriced", "err", sendError)
		return commonclient.Successful
	}
	if sendError.IsTerminallyUnderpriced(configErrors) {
		lggr.Errorw("Transaction terminally underpriced", "etx", tx, "err", sendError)
		return commonclient.Underpriced
	}
	if sendError.L2FeeTooLow(configErrors) || sendError.IsL2FeeTooHigh(configErrors) || sendError.IsL2Full(configErrors) {
		if isL2 {
			lggr.Errorw("Transaction fee out of range", "err", sendError, "etx", tx)
			return commonclient.FeeOutOfValidRange
		}
		lggr.Errorw("this error type only handled for L2s", "err", sendError, "etx", tx)
		return commonclient.Unsupported
	}
	if sendError.IsNonceTooHighError(configErrors) {
		// This error occurs when the tx nonce is greater than current_nonce + tx_count_in_mempool,
		// instead of keeping the tx in mempool. This can happen if previous transactions haven't
		// reached the client yet. The correct thing to do is to mark it as retryable.
		lggr.Warnw("Transaction has a nonce gap.", "err", sendError, "etx", tx)
		return commonclient.Retryable
	}
	if sendError.IsInsufficientEth(configErrors) {
		lggr.Criticalw(fmt.Sprintf("Tx %x with type 0x%d was rejected due to insufficient eth: %s\n"+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			tx.Hash(), tx.Type(), sendError.Error(), fromAddress,
		), "err", sendError, "etx", tx)
		return commonclient.InsufficientFunds
	}
	if sendError.IsServiceUnavailable(configErrors) {
		lggr.Errorw(fmt.Sprintf("service unavailable while sending transaction %x", tx.Hash()), "err", sendError, "etx", tx)
		return commonclient.Retryable
	}
	if sendError.IsTimeout() {
		lggr.Errorw(fmt.Sprintf("timeout while sending transaction %x", tx.Hash()), "err", sendError, "etx", tx)
		return commonclient.Retryable
	}
	if sendError.IsCanceled() {
		lggr.Errorw(fmt.Sprintf("context was canceled while sending transaction %x", tx.Hash()), "err", sendError, "etx", tx)
		return commonclient.Retryable
	}
	if sendError.IsTxFeeExceedsCap(configErrors) {
		lggr.Criticalw(fmt.Sprintf("Sending transaction failed: %s", label.RPCTxFeeCapConfiguredIncorrectlyWarning),
			"etx", tx,
			"err", sendError,
			"id", "RPCTxFeeCapExceeded",
		)
		return commonclient.ExceedsMaxFee
	}
	lggr.Criticalw("Unknown error encountered when sending transaction", "err", err, "etx", tx)
	return commonclient.Unknown
}
