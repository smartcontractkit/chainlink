//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	solanago "github.com/gagliardetto/solana-go"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	solanabindings "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana/bindings"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/solana/solwrite-negative/config"
)

// Function names, they must literally match the `functionToTest` values used by the
// regression test (system-tests/tests/regression/cre/solana_regression_test.go).
const (
	writeReportInvalidReceiverLength   = "WriteReport - invalid receiver length"
	writeReportZeroReceiver            = "WriteReport - zero receiver"
	writeReportNonExistentReceiver     = "WriteReport - non-existent receiver"
	writeReportNonExecutableReceiver   = "WriteReport - non-executable receiver"
	writeReportTooFewRemainingAccounts = "WriteReport - too few remaining accounts"
	writeReportInvalidRemainingAccount = "WriteReport - invalid remaining account key"
	writeReportForwarderStateMismatch  = "WriteReport - mismatched forwarder state"
	writeReportAccountHashMismatch     = "WriteReport - remaining accounts hash mismatch"
	writeReportComputeLimitExceeded    = "WriteReport - compute limit exceeded"
)

// defaultComputeLimit stays below PerWorkflow.ChainWrite.Solana.GasLimit (300_000), so that
// cases which are not about the compute limit are never rejected by the compute limiter first.
const defaultComputeLimit = uint32(290_000)

// maxWriteCallsPerExecution mirrors PerWorkflow.ChainWrite.TargetsLimit. The engine counts every
// attempted WriteReport, including the ones the capability rejects, so a batch that goes over the
// limit would start failing with a call-limit error instead of the error each case is about.
const maxWriteCallsPerExecution = 10

// reportPayload is an arbitrary, small payload: every case here is expected to fail during
// request validation, so the payload is never CPI-ed into the receiver program.
var reportPayload = []byte{1, 2, 3}

func main() {
	wasm.NewRunner(parseConfig).Run(RunSolanaWriteNegativeWorkflow)
}

func parseConfig(b []byte) (config.Config, error) {
	var wfCfg config.Config
	if err := yaml.Unmarshal(b, &wfCfg); err != nil {
		return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return wfCfg, nil
}

func RunSolanaWriteNegativeWorkflow(wfCfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onSolanaWriteTrigger,
		),
	}, nil
}

// writeAttempt describes a single WriteReport call: what is sent to the capability and what the
// report commits to. hashedAccounts is what the report's account hash is computed over; it is
// identical to remainingAccounts for every case except the account-hash-mismatch one.
type writeAttempt struct {
	receiver          []byte
	remainingAccounts []*solana.AccountMeta
	hashedAccounts    []*solana.AccountMeta
	computeLimit      uint32
}

// onSolanaWriteTrigger runs every case of the batch in one execution and verifies each of them
// in place, so that the test only has to wait for a single "batch passed" log.
func onSolanaWriteTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (any, error) {
	runtime.Logger().Info("onSolanaWriteTrigger called", "batch", wfCfg.BatchName, "cases", len(wfCfg.Cases), "payload", payload)

	if len(wfCfg.Cases) == 0 {
		return nil, fmt.Errorf("batch '%s' has no cases to run", wfCfg.BatchName)
	}
	if len(wfCfg.Cases) > maxWriteCallsPerExecution {
		return nil, fmt.Errorf("batch '%s' has %d cases, which is over the %d writes a single execution is allowed",
			wfCfg.BatchName, len(wfCfg.Cases), maxWriteCallsPerExecution)
	}

	client := solana.Client{ChainSelector: wfCfg.ChainSelector}
	runtime.Logger().Info("Solana client created", "chainSelector", wfCfg.ChainSelector)

	for _, testCase := range wfCfg.Cases {
		callErr, err := runCase(client, runtime, wfCfg, testCase)
		if err != nil {
			runtime.Logger().Error(fmt.Sprintf("Solana write negative batch '%s' could not run case '%s' (%s): %v",
				wfCfg.BatchName, testCase.Name, testCase.FunctionToTest, err))
			return nil, fmt.Errorf("could not run case '%s' (%s): %w", testCase.Name, testCase.FunctionToTest, err)
		}

		if err := verify(testCase, callErr); err != nil {
			runtime.Logger().Error(fmt.Sprintf("Solana write negative batch '%s' failed on case '%s' (%s): %v",
				wfCfg.BatchName, testCase.Name, testCase.FunctionToTest, err))
			return nil, fmt.Errorf("case '%s' (%s) failed: %w", testCase.Name, testCase.FunctionToTest, err)
		}

		runtime.Logger().Info(fmt.Sprintf("case '%s' (%s) passed: got expected error: %v", testCase.Name, testCase.FunctionToTest, callErr))
	}

	// The test waits for this line, and for this line only.
	runtime.Logger().Info(fmt.Sprintf("Solana write negative batch '%s' passed (%d cases)", wfCfg.BatchName, len(wfCfg.Cases)))
	return fmt.Sprintf("batch '%s' passed", wfCfg.BatchName), nil
}

// runCase performs the WriteReport call of a single case. The returned error means the case could
// not be run at all (for example its input could not be decoded); the capability's own error is
// returned as callErr, because for these cases it is the expected result.
func runCase(client solana.Client, runtime cre.Runtime, wfCfg config.Config, testCase config.Case) (callErr error, err error) {
	runtime.Logger().Info("running case", "name", testCase.Name, "functionToTest", testCase.FunctionToTest, "invalidInput", testCase.InvalidInput)

	attempt, err := buildAttempt(wfCfg, testCase)
	if err != nil {
		return nil, err
	}

	return writeReport(runtime, client, attempt)
}

// verify checks what the capability answered against what the case expects. The error it returns
// describes the mismatch, so a failing batch names the case and what it actually got.
func verify(testCase config.Case, callErr error) error {
	if callErr == nil {
		return fmt.Errorf("expected an error containing %q, got a successful write", testCase.ExpectedError)
	}
	if !strings.Contains(callErr.Error(), testCase.ExpectedError) {
		return fmt.Errorf("expected an error containing %q, got: %v", testCase.ExpectedError, callErr)
	}
	return nil
}

// buildAttempt returns the (intentionally invalid) WriteReport inputs for the given case.
func buildAttempt(wfCfg config.Config, testCase config.Case) (writeAttempt, error) {
	remainingAccounts, err := validRemainingAccounts(wfCfg)
	if err != nil {
		return writeAttempt{}, err
	}

	attempt := writeAttempt{
		receiver:          wfCfg.Receiver.Bytes(),
		remainingAccounts: remainingAccounts,
		computeLimit:      defaultComputeLimit,
	}

	switch testCase.FunctionToTest {
	case writeReportInvalidReceiverLength, writeReportNonExistentReceiver:
		// The test supplies the raw receiver bytes: a key that is not 32 bytes long, or a
		// well-formed key of an account that was never created on-chain.
		receiver, decodeErr := decodeHex(testCase.InvalidInput)
		if decodeErr != nil {
			return writeAttempt{}, fmt.Errorf("failed to decode receiver '%s': %w", testCase.InvalidInput, decodeErr)
		}
		attempt.receiver = receiver

	case writeReportZeroReceiver:
		attempt.receiver = make([]byte, solanago.PublicKeyLength)

	case writeReportNonExecutableReceiver:
		attempt.receiver = wfCfg.NonExecutableAccount.Bytes()

	case writeReportTooFewRemainingAccounts:
		// The forwarder layout requires at least the state and the authority accounts.
		attempt.remainingAccounts = remainingAccounts[:1]

	case writeReportInvalidRemainingAccount:
		invalidKey, decodeErr := decodeHex(testCase.InvalidInput)
		if decodeErr != nil {
			return writeAttempt{}, fmt.Errorf("failed to decode remaining account key '%s': %w", testCase.InvalidInput, decodeErr)
		}
		attempt.remainingAccounts = append(attempt.remainingAccounts, &solana.AccountMeta{PublicKey: invalidKey})

	case writeReportForwarderStateMismatch:
		// Index 0 must be the forwarder state the capability is configured with.
		tampered := cloneAccounts(remainingAccounts)
		tampered[0] = &solana.AccountMeta{PublicKey: wfCfg.Receiver.Bytes()}
		attempt.remainingAccounts = tampered

	case writeReportAccountHashMismatch:
		// The report commits to the valid account list, but a different list is submitted.
		attempt.hashedAccounts = remainingAccounts
		attempt.remainingAccounts = append(cloneAccounts(remainingAccounts), &solana.AccountMeta{PublicKey: wfCfg.ForwarderProgramID.Bytes()})

	case writeReportComputeLimitExceeded:
		// Everything else stays valid so that validation passes and the compute limiter is reached.
		computeLimit, parseErr := strconv.ParseUint(testCase.InvalidInput, 10, 32)
		if parseErr != nil {
			return writeAttempt{}, fmt.Errorf("failed to parse compute limit '%s': %w", testCase.InvalidInput, parseErr)
		}
		attempt.computeLimit = uint32(computeLimit) //nolint:gosec // ParseUint above bounds the value to 32 bits

	default:
		return writeAttempt{}, fmt.Errorf("the provided name for function to test in regression Solana Write Workflow did not match any known functions: %s", testCase.FunctionToTest)
	}

	if attempt.hashedAccounts == nil {
		attempt.hashedAccounts = attempt.remainingAccounts
	}

	return attempt, nil
}

// validRemainingAccounts builds the account list in the layout the keystone forwarder expects:
// index 0 is the forwarder state, index 1 the forwarder authority PDA, index 2+ receiver accounts.
func validRemainingAccounts(wfCfg config.Config) ([]*solana.AccountMeta, error) {
	authority, err := deriveForwarderAuthority(wfCfg.ForwarderState, wfCfg.Receiver, wfCfg.ForwarderProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive forwarder authority: %w", err)
	}

	return []*solana.AccountMeta{
		{PublicKey: wfCfg.ForwarderState.Bytes()},
		{PublicKey: authority.Bytes()},
		{PublicKey: wfCfg.Receiver.Bytes(), IsWritable: true},
	}, nil
}

func deriveForwarderAuthority(forwarderState, receiverProgram, forwarderProgram solanago.PublicKey) (solanago.PublicKey, error) {
	seeds := [][]byte{
		[]byte("forwarder"),
		forwarderState[:],
		receiverProgram[:],
	}
	authority, _, err := solanago.FindProgramAddress(seeds, forwarderProgram)
	return authority, err
}

// writeReport submits the report and returns the capability error, which is the expected result
// for every case here. The second error means the report could not be generated at all.
func writeReport(runtime cre.Runtime, client solana.Client, attempt writeAttempt) (callErr error, err error) {
	report, err := generateForwarderReport(runtime, attempt.hashedAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}

	runtime.Logger().Info("Writing report...",
		"receiver", hex.EncodeToString(attempt.receiver),
		"remainingAccounts", formatAccounts(attempt.remainingAccounts),
		"computeLimit", attempt.computeLimit,
	)

	reply, callErr := client.WriteReport(runtime, &solana.WriteCreReportRequest{
		Receiver:          attempt.receiver,
		RemainingAccounts: attempt.remainingAccounts,
		ComputeConfig:     &solana.ComputeConfig{ComputeLimit: attempt.computeLimit},
		Report:            report,
	}).Await()
	if callErr == nil {
		runtime.Logger().Info("WriteReport unexpectedly succeeded", "write_output", reply)
	}

	return callErr, nil
}

// generateForwarderReport produces a report whose payload is the Borsh-encoded forwarder report,
// committing to the hash of the given accounts - the same encoding the generated Solana bindings use.
func generateForwarderReport(runtime cre.Runtime, hashedAccounts []*solana.AccountMeta) (*cre.Report, error) {
	forwarderReport := solanabindings.ForwarderReport{
		AccountHash: solanabindings.CalculateAccountsHash(hashedAccounts),
		Payload:     reportPayload,
	}

	encoded, err := forwarderReport.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to encode forwarder report: %w", err)
	}

	return runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "solana",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
}

func cloneAccounts(accounts []*solana.AccountMeta) []*solana.AccountMeta {
	cloned := make([]*solana.AccountMeta, len(accounts))
	copy(cloned, accounts)
	return cloned
}

func formatAccounts(accounts []*solana.AccountMeta) string {
	formatted := make([]string, len(accounts))
	for i, account := range accounts {
		formatted[i] = hex.EncodeToString(account.GetPublicKey())
	}
	return strings.Join(formatted, ",")
}

// decodeHex decodes a hex-encoded byte string of any length. An empty string decodes to no bytes,
// which is how the "empty receiver" style cases are expressed.
func decodeHex(s string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}
