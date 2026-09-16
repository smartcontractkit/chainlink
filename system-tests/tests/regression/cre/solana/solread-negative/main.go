//go:build wasip1

package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/solana/solread-negative/config"
)

// Function names, they must literally match the `functionToTest` values used by the
// regression test (system-tests/tests/regression/cre/solana_regression_test.go).
const (
	getAccountInfoInvalidAccount      = "GetAccountInfoWithOpts - invalid account"
	getAccountInfoNonExistentAccount  = "GetAccountInfoWithOpts - non-existent account"
	getAccountInfoInvalidCommitment   = "GetAccountInfoWithOpts - invalid commitment"
	getAccountInfoInvalidEncoding     = "GetAccountInfoWithOpts - invalid encoding"
	getBalanceInvalidAddress          = "GetBalance - invalid address"
	getBalanceNonExistentAccount      = "GetBalance - non-existent account"
	getMultipleAccountsInvalidAccount = "GetMultipleAccountsWithOpts - invalid account"
	getMultipleAccountsTooManyItems   = "GetMultipleAccountsWithOpts - too many accounts"
	getProgramAccountsInvalidProgram  = "GetProgramAccounts - invalid program"
	getProgramAccountsNonProgram      = "GetProgramAccounts - non-program account"
	getBlockInvalidSlot               = "GetBlock - invalid slot"
	getSlotHeightInvalidCommitment    = "GetSlotHeight - invalid commitment"
	getTransactionInvalidSignature    = "GetTransaction - invalid signature"
	getTransactionUnknownSignature    = "GetTransaction - unknown signature"
	getSignatureStatusesInvalidSig    = "GetSignatureStatuses - invalid signature"
	getSignatureStatusesTooManyItems  = "GetSignatureStatuses - too many signatures"
	getFeeForMessageInvalidMessage    = "GetFeeForMessage - invalid message"
)

// maxReadCallsPerExecution mirrors PerWorkflow.ChainRead.CallLimit. The engine counts every
// attempted chain read, including the ones the capability rejects, so a batch that goes over the
// limit would start failing with a call-limit error instead of the error each case is about.
const maxReadCallsPerExecution = 15

func main() {
	wasm.NewRunner(parseConfig).Run(RunSolanaReadNegativeWorkflow)
}

func parseConfig(b []byte) (config.Config, error) {
	var wfCfg config.Config
	if err := yaml.Unmarshal(b, &wfCfg); err != nil {
		return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return wfCfg, nil
}

func RunSolanaReadNegativeWorkflow(wfCfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onSolanaReadTrigger,
		),
	}, nil
}

// outcome is what a single case got back from the capability.
type outcome struct {
	// empty reports whether the capability answered successfully but with nothing in it.
	empty bool
	// err is the capability error, nil when the call succeeded.
	err error
}

// onSolanaReadTrigger runs every case of the batch in one execution and verifies each of them
// in place, so that the test only has to wait for a single "batch passed" log.
func onSolanaReadTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (any, error) {
	runtime.Logger().Info("onSolanaReadTrigger called", "batch", wfCfg.BatchName, "cases", len(wfCfg.Cases), "payload", payload)

	if len(wfCfg.Cases) == 0 {
		return nil, fmt.Errorf("batch '%s' has no cases to run", wfCfg.BatchName)
	}
	if len(wfCfg.Cases) > maxReadCallsPerExecution {
		return nil, fmt.Errorf("batch '%s' has %d cases, which is over the %d chain reads a single execution is allowed",
			wfCfg.BatchName, len(wfCfg.Cases), maxReadCallsPerExecution)
	}

	client := solana.Client{ChainSelector: wfCfg.ChainSelector}

	for _, testCase := range wfCfg.Cases {
		result, err := runCase(client, runtime, testCase)
		if err != nil {
			runtime.Logger().Error(fmt.Sprintf("Solana read negative batch '%s' could not run case '%s' (%s): %v",
				wfCfg.BatchName, testCase.Name, testCase.FunctionToTest, err))
			return nil, fmt.Errorf("could not run case '%s' (%s): %w", testCase.Name, testCase.FunctionToTest, err)
		}

		if err := verify(testCase, result); err != nil {
			runtime.Logger().Error(fmt.Sprintf("Solana read negative batch '%s' failed on case '%s' (%s): %v",
				wfCfg.BatchName, testCase.Name, testCase.FunctionToTest, err))
			return nil, fmt.Errorf("case '%s' (%s) failed: %w", testCase.Name, testCase.FunctionToTest, err)
		}

		runtime.Logger().Info(fmt.Sprintf("case '%s' (%s) passed: %s", testCase.Name, testCase.FunctionToTest, describe(result)))
	}

	// The test waits for this line, and for this line only.
	runtime.Logger().Info(fmt.Sprintf("Solana read negative batch '%s' passed (%d cases)", wfCfg.BatchName, len(wfCfg.Cases)))
	return fmt.Sprintf("batch '%s' passed", wfCfg.BatchName), nil
}

// verify checks what the capability answered against what the case expects. The error it returns
// describes the mismatch, so a failing batch names the case and what it actually got.
func verify(testCase config.Case, result outcome) error {
	switch testCase.Outcome {
	case config.OutcomeError:
		if result.err == nil {
			return fmt.Errorf("expected an error containing %q, got a successful response", testCase.ExpectedError)
		}
		if !strings.Contains(result.err.Error(), testCase.ExpectedError) {
			return fmt.Errorf("expected an error containing %q, got: %v", testCase.ExpectedError, result.err)
		}
	case config.OutcomeErrorOrEmpty:
		if result.err == nil && !result.empty {
			return fmt.Errorf("expected an error or an empty response, got data back")
		}
	default:
		return fmt.Errorf("unknown outcome: %d", testCase.Outcome)
	}

	return nil
}

func describe(result outcome) string {
	if result.err != nil {
		return fmt.Sprintf("got expected error: %v", result.err)
	}
	return "got expected empty response"
}

// runCase performs the capability call of a single case. The returned error means the case could
// not be run at all (for example its input could not be decoded); the capability's own error is
// reported through the outcome instead, because for these cases it is the expected result.
func runCase(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	runtime.Logger().Info("running case", "name", testCase.Name, "functionToTest", testCase.FunctionToTest, "invalidInput", testCase.InvalidInput)

	switch testCase.FunctionToTest {
	case getAccountInfoInvalidAccount, getAccountInfoNonExistentAccount:
		return getAccountInfo(client, runtime, testCase)
	case getAccountInfoInvalidCommitment:
		return getAccountInfoWithInvalidCommitment(client, runtime, testCase)
	case getAccountInfoInvalidEncoding:
		return getAccountInfoWithInvalidEncoding(client, runtime, testCase)
	case getBalanceInvalidAddress, getBalanceNonExistentAccount:
		return getBalance(client, runtime, testCase)
	case getMultipleAccountsInvalidAccount:
		return getMultipleAccountsWithInvalidAccount(client, runtime, testCase)
	case getMultipleAccountsTooManyItems:
		return getMultipleAccountsOverBatchLimit(client, runtime, testCase)
	case getProgramAccountsInvalidProgram, getProgramAccountsNonProgram:
		return getProgramAccounts(client, runtime, testCase)
	case getBlockInvalidSlot:
		return getBlock(client, runtime, testCase)
	case getSlotHeightInvalidCommitment:
		return getSlotHeightWithInvalidCommitment(client, runtime, testCase)
	case getTransactionInvalidSignature, getTransactionUnknownSignature:
		return getTransaction(client, runtime, testCase)
	case getSignatureStatusesInvalidSig:
		return getSignatureStatusesWithInvalidSignature(client, runtime, testCase)
	case getSignatureStatusesTooManyItems:
		return getSignatureStatusesOverBatchLimit(client, runtime, testCase)
	case getFeeForMessageInvalidMessage:
		return getFeeForMessage(client, runtime, testCase)
	default:
		return outcome{}, fmt.Errorf("the provided name for function to test in regression Solana Read Workflow did not match any known functions: %s", testCase.FunctionToTest)
	}
}

// getAccountInfo reads an account that is either not a 32-byte public key or was never created.
func getAccountInfo(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	account, err := decodeHex(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode account '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetAccountInfoWithOpts(runtime, &solana.GetAccountInfoWithOptsRequest{
		Account: account,
		Opts:    defaultAccountInfoOpts(),
	}).Await()

	return outcome{empty: reply == nil || reply.Value == nil, err: err}, nil
}

// getAccountInfoWithInvalidCommitment sends a commitment level the capability does not know.
func getAccountInfoWithInvalidCommitment(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	commitment, err := parseInt32(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode commitment '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetAccountInfoWithOpts(runtime, &solana.GetAccountInfoWithOptsRequest{
		Account: validAccount(),
		Opts: &solana.GetAccountInfoOpts{
			Encoding:   solana.EncodingType_ENCODING_TYPE_JSON_PARSED,
			Commitment: solana.CommitmentType(commitment),
		},
	}).Await()

	return outcome{empty: reply == nil || reply.Value == nil, err: err}, nil
}

// getAccountInfoWithInvalidEncoding sends an encoding the capability does not know.
func getAccountInfoWithInvalidEncoding(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	encoding, err := parseInt32(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode encoding '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetAccountInfoWithOpts(runtime, &solana.GetAccountInfoWithOptsRequest{
		Account: validAccount(),
		Opts: &solana.GetAccountInfoOpts{
			Encoding:   solana.EncodingType(encoding),
			Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
		},
	}).Await()

	return outcome{empty: reply == nil || reply.Value == nil, err: err}, nil
}

// getBalance reads the balance of an address that is either not a 32-byte public key or was
// never funded - Solana reports zero lamports for the latter rather than failing.
func getBalance(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	address, err := decodeHex(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode address '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetBalance(runtime, &solana.GetBalanceRequest{
		Addr:       address,
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}).Await()

	return outcome{empty: reply == nil || reply.Value == 0, err: err}, nil
}

// getMultipleAccountsWithInvalidAccount puts a malformed key into an otherwise valid batch.
func getMultipleAccountsWithInvalidAccount(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	account, err := decodeHex(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode account '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetMultipleAccountsWithOpts(runtime, &solana.GetMultipleAccountsWithOptsRequest{
		Accounts: [][]byte{account, validAccount()},
		Opts:     defaultMultipleAccountsOpts(),
	}).Await()

	return outcome{empty: reply == nil || len(reply.Value) == 0, err: err}, nil
}

// getMultipleAccountsOverBatchLimit asks for more accounts than
// PerWorkflow.ChainRead.Solana.BatchItemLimit allows. Every key in the batch is valid, so the
// request can only be rejected by the batch limiter.
func getMultipleAccountsOverBatchLimit(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	count, err := strconv.Atoi(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode account count '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetMultipleAccountsWithOpts(runtime, &solana.GetMultipleAccountsWithOptsRequest{
		Accounts: generateAccounts(count),
		Opts:     defaultMultipleAccountsOpts(),
	}).Await()

	return outcome{empty: reply == nil || len(reply.Value) == 0, err: err}, nil
}

// getProgramAccounts reads accounts owned by a key that is either not a public key at all or is
// not a program, in which case nothing is owned under it.
func getProgramAccounts(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	program, err := decodeHex(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode program '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetProgramAccounts(runtime, &solana.GetProgramAccountsRequest{
		Program: program,
		Opts: &solana.GetProgramAccountsOpts{
			Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
		},
	}).Await()

	return outcome{empty: reply == nil || len(reply.Value) == 0, err: err}, nil
}

// getBlock reads a slot the local validator has never produced.
func getBlock(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	slot, err := strconv.ParseUint(testCase.InvalidInput, 10, 64)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode slot '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetBlock(runtime, &solana.GetBlockRequest{
		Slot: slot,
		Opts: &solana.GetBlockOpts{
			Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
		},
	}).Await()

	return outcome{empty: reply == nil || len(reply.Blockhash) == 0, err: err}, nil
}

// getSlotHeightWithInvalidCommitment sends a commitment level the capability does not know.
func getSlotHeightWithInvalidCommitment(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	commitment, err := parseInt32(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode commitment '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetSlotHeight(runtime, &solana.GetSlotHeightRequest{
		Commitment: solana.CommitmentType(commitment),
	}).Await()

	return outcome{empty: reply == nil || reply.Height == 0, err: err}, nil
}

// getTransaction reads a transaction by a signature that is either not 64 bytes or belongs to a
// transaction that was never submitted.
func getTransaction(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	signature, err := decodeHex(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode signature '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetTransaction(runtime, &solana.GetTransactionRequest{
		Signature: signature,
	}).Await()

	return outcome{empty: reply == nil || reply.Slot == 0, err: err}, nil
}

// getSignatureStatusesWithInvalidSignature puts a malformed signature into the batch.
func getSignatureStatusesWithInvalidSignature(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	signature, err := decodeHex(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode signature '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetSignatureStatuses(runtime, &solana.GetSignatureStatusesRequest{
		Sigs: [][]byte{signature},
	}).Await()

	return outcome{empty: reply == nil || len(reply.Results) == 0, err: err}, nil
}

// getSignatureStatusesOverBatchLimit asks for more statuses than
// PerWorkflow.ChainRead.Solana.BatchItemLimit allows.
func getSignatureStatusesOverBatchLimit(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	count, err := strconv.Atoi(testCase.InvalidInput)
	if err != nil {
		return outcome{}, fmt.Errorf("failed to decode signature count '%s': %w", testCase.InvalidInput, err)
	}

	reply, err := client.GetSignatureStatuses(runtime, &solana.GetSignatureStatusesRequest{
		Sigs: generateSignatures(count),
	}).Await()

	return outcome{empty: reply == nil || len(reply.Results) == 0, err: err}, nil
}

// getFeeForMessage estimates the fee of something that is not an encoded Solana message.
func getFeeForMessage(client solana.Client, runtime cre.Runtime, testCase config.Case) (outcome, error) {
	reply, err := client.GetFeeForMessage(runtime, &solana.GetFeeForMessageRequest{
		Message:    testCase.InvalidInput,
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}).Await()

	return outcome{empty: reply == nil || reply.Fee == 0, err: err}, nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func defaultAccountInfoOpts() *solana.GetAccountInfoOpts {
	return &solana.GetAccountInfoOpts{
		Encoding:   solana.EncodingType_ENCODING_TYPE_JSON_PARSED,
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}
}

func defaultMultipleAccountsOpts() *solana.GetMultipleAccountsOpts {
	return &solana.GetMultipleAccountsOpts{
		Encoding:   solana.EncodingType_ENCODING_TYPE_JSON_PARSED,
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}
}

// validAccount is the system program, a well-formed key that always exists on-chain. It is used
// by the cases whose invalid input is something other than the account itself.
func validAccount() []byte {
	return make([]byte, solanaPublicKeyLength)
}

const (
	solanaPublicKeyLength = 32
	solanaSignatureLength = 64
)

// generateAccounts builds distinct, well-formed public keys to fill a batch request with.
func generateAccounts(count int) [][]byte {
	accounts := make([][]byte, count)
	for i := range accounts {
		account := make([]byte, solanaPublicKeyLength)
		binary.BigEndian.PutUint64(account, uint64(i)+1)
		accounts[i] = account
	}
	return accounts
}

// generateSignatures builds distinct, well-formed signatures to fill a batch request with.
func generateSignatures(count int) [][]byte {
	signatures := make([][]byte, count)
	for i := range signatures {
		signature := make([]byte, solanaSignatureLength)
		binary.BigEndian.PutUint64(signature, uint64(i)+1)
		signatures[i] = signature
	}
	return signatures
}

func parseInt32(s string) (int32, error) {
	parsed, err := strconv.ParseInt(s, 10, 32)
	return int32(parsed), err
}

// decodeHex decodes a hex-encoded byte string of any length. An empty string decodes to no bytes,
// which reaches the capability as a nil field.
func decodeHex(s string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}
