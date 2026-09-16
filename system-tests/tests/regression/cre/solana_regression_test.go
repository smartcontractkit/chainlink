package cre

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	ks_sol "github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana"
	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	sollogtrigger_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/solana/sollogtrigger-negative/config"
	solread_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/solana/solread-negative/config"
	solwrite_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/solana/solwrite-negative/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// Solana regression tests need a topology that runs a Solana chain and the Solana capability.
const solanaRegressionConfigPath = "/configs/workflow-don-solana.toml"

const (
	solanaWriteNegativeWorkflowFile      = "./solana/solwrite-negative/main.go"
	solanaReadNegativeWorkflowFile       = "./solana/solread-negative/main.go"
	solanaLogTriggerNegativeWorkflowFile = "./solana/sollogtrigger-negative/main.go"
)

// Read and write cases run in batches: one workflow deployment covers a whole batch, instead of
// paying the compile/deploy/trigger cost per case. A batch runs in a single workflow execution,
// so it is bounded by the engine's per-execution call limits - and the engine counts the calls
// the capability rejects too, which is every call these tests make.
const (
	// solanaMaxReadCallsPerBatch mirrors PerWorkflow.ChainRead.CallLimit.
	solanaMaxReadCallsPerBatch = 15
	// solanaMaxWriteCallsPerBatch mirrors PerWorkflow.ChainWrite.TargetsLimit.
	solanaMaxWriteCallsPerBatch = 10
)

// solanaTestNameTemplate is a template for Solana negative test names to avoid duplication,
// e.g. "Solana.<Function> fails with <invalid input>".
const solanaTestNameTemplate = "Solana.%s fails with %s"

// solanaNegativeTest is one negative case. For read cases an empty expectedError means the
// request is well-formed, so the capability may answer with either an error or an empty
// response - returning real data is the only failure.
//
// The tables below hold one case per validation branch, not one per malformed input: the
// capability checks a key's length, not its shape, so an empty, one-byte, short and long key all
// take the same branch and prove the same thing. Inputs are only repeated across methods where
// the code path genuinely differs - a singular field (GetBalance.Addr) and a repeated one
// (GetMultipleAccounts.Accounts) are validated by different converters, so both are covered.
type solanaNegativeTest struct {
	name           string
	invalidInput   string
	functionToTest string
	expectedError  string
}

// nonExistentSolanaAccount is a well-formed 32-byte key that no test ever creates or funds,
// so the account behind it never exists on the local validator.
const nonExistentSolanaAccount = "beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef"

// Malformed keys and signatures: they are passed to the capability as raw bytes, so unlike EVM
// addresses they are never padded or truncated into something valid.
var (
	solanaShortKey         = strings.Repeat("ab", 31)
	solanaShortSignature   = strings.Repeat("ab", 63)
	solanaUnknownSignature = strings.Repeat("ab", 64)
)

// solanaFundingAmount is 1 SOL, enough to make the funded address exist on-chain as a
// system-owned (and therefore non-executable) account.
var solanaFundingAmount = big.NewInt(1_000_000_000)

//////////////////////////////////////////////////////
// WRITE NEGATIVE TESTS
//////////////////////////////////////////////////////

// ...Function variables must literally match the names of the switch-case statements in the
// workflow (solana/solwrite-negative/main.go). In each case the Solana capability's WriteReport
// is called with one invalid input, and the capability is expected to reject the request during
// validation, before anything is submitted on-chain.
const (
	solanaWriteReportInvalidReceiverLength   = "WriteReport - invalid receiver length"
	solanaWriteReportZeroReceiver            = "WriteReport - zero receiver"
	solanaWriteReportNonExistentReceiver     = "WriteReport - non-existent receiver"
	solanaWriteReportNonExecutableReceiver   = "WriteReport - non-executable receiver"
	solanaWriteReportTooFewRemainingAccounts = "WriteReport - too few remaining accounts"
	solanaWriteReportInvalidRemainingAccount = "WriteReport - invalid remaining account key"
	solanaWriteReportForwarderStateMismatch  = "WriteReport - mismatched forwarder state"
	solanaWriteReportAccountHashMismatch     = "WriteReport - remaining accounts hash mismatch"
	solanaWriteReportComputeLimitExceeded    = "WriteReport - compute limit exceeded"
)

const (
	expectedSolanaReceiverNotPublicKey    = "received public key is not 32 bytes long"
	expectedSolanaReceiverEmpty           = "receiver public key is empty"
	expectedSolanaReceiverMissing         = "receiver account does not exist"
	expectedSolanaReceiverNonExecutable   = "receiver account is non-executable"
	expectedSolanaTooFewRemainingAccounts = "expected accounts meta length > 2"
	expectedSolanaInvalidRemainingAccount = "public key must be exactly 32 bytes"
	expectedSolanaForwarderStateMismatch  = "doesn't match configured forwarder state"
	expectedSolanaAccountHashMismatch     = "remaining account hash mismatch"
	expectedSolanaComputeLimitExceeded    = "provided compute config exceeds limit"
)

var solanaNegativeTestsWriteReportInvalidReceiver = []solanaNegativeTest{
	// WriteReport - receiver the forwarder can never deliver to. The capability checks the
	// receiver's length, so one malformed length stands in for empty/short/long keys.
	{"31 bytes (short) key", solanaShortKey, solanaWriteReportInvalidReceiverLength, expectedSolanaReceiverNotPublicKey},
	{"zero key", "", solanaWriteReportZeroReceiver, expectedSolanaReceiverEmpty},
	{"non-existent account", nonExistentSolanaAccount, solanaWriteReportNonExistentReceiver, expectedSolanaReceiverMissing},
	// A funded system account exists on-chain but is not a program, so it can never receive a report.
	{"funded system account", "", solanaWriteReportNonExecutableReceiver, expectedSolanaReceiverNonExecutable},
}

var solanaNegativeTestsWriteReportInvalidPayload = []solanaNegativeTest{
	// WriteReport - remaining accounts that do not match the forwarder's expected layout.
	// Index 0 must be the configured forwarder state and index 1 the forwarder authority PDA.
	{"only forwarder state", "", solanaWriteReportTooFewRemainingAccounts, expectedSolanaTooFewRemainingAccounts},
	{"31 bytes (short) key", solanaShortKey, solanaWriteReportInvalidRemainingAccount, expectedSolanaInvalidRemainingAccount},
	{"foreign forwarder state", "", solanaWriteReportForwarderStateMismatch, expectedSolanaForwarderStateMismatch},
	// The report commits to a hash of the accounts, so submitting a different account list must be rejected.
	{"accounts not matching the report hash", "", solanaWriteReportAccountHashMismatch, expectedSolanaAccountHashMismatch},
	// WriteReport - compute limit above PerWorkflow.ChainWrite.Solana.GasLimit (300_000 CUs by
	// default). Any value over the limit takes the same branch, so one is enough.
	{"above the configured limit", "300001", solanaWriteReportComputeLimitExceeded, expectedSolanaComputeLimitExceeded},
}

// SolanaWriteFailsTest deploys the Solana write negative workflow once for a whole batch of cases.
// The workflow runs every case in a single execution and checks each of them against its expected
// error, so the test waits for one "batch passed" log. A failing case fails the batch with a log
// that names the case and what it actually got back.
func SolanaWriteFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, batchName string, tests []solanaNegativeTest) {
	testLogger := framework.L

	cases := toSolanaWriteCases(tests)

	solChain := mustSolanaChainInEnv(t, testEnv)
	chainSelector := solChain.ChainSelector()
	dataStore := testEnv.CreEnvironment.CldfEnvironment.DataStore

	forwarderProgramID := mustGetSolanaContract(t, dataStore, chainSelector, ks_sol.ForwarderContract)
	forwarderState := mustGetSolanaContract(t, dataStore, chainSelector, ks_sol.ForwarderState)

	userLogsCh, baseMessageCh := startSolanaChipSink(t, testLogger)

	testLogger.Info().Msgf("Creating Solana Write Fail workflow configuration for batch '%s' (%d cases)...", batchName, len(cases))
	workflowConfig := solwrite_negative_config.Config{
		ChainSelector: chainSelector,
		BatchName:     batchName,
		Cases:         cases,
		// The forwarder program is a deployed, executable program, so it is a valid receiver for
		// the cases that are not about the receiver itself.
		Receiver:           forwarderProgramID,
		ForwarderProgramID: forwarderProgramID,
		ForwarderState:     forwarderState,
	}

	if batchNeedsNonExecutableAccount(cases) {
		workflowConfig.NonExecutableAccount = mustCreateFundedSolanaAccount(t, testLogger, solChain)
	}

	workflowName := fmt.Sprintf("sol-write-fail-workflow-%d-%04d", chainSelector, rand.Intn(10000))
	workflowID := t_helpers.CompileAndDeployWorkflow(
		t,
		testEnv,
		testLogger,
		workflowName,
		&workflowConfig,
		solanaWriteNegativeWorkflowFile,
	)

	t_helpers.WatchWorkflowLogs(
		t,
		testLogger,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		solanaWriteBatchPassedLog(batchName, len(cases)),
		2*time.Minute,
		t_helpers.WithUserLogWorkflowID(workflowID),
	)
	testLogger.Info().Msgf("Solana Write Fail test successfully completed for batch '%s'", batchName)
}

// runSolanaWriteNegativeTestSuite runs one batch of Solana write negative cases.
func runSolanaWriteNegativeTestSuite(t *testing.T, batchName string, tests []solanaNegativeTest) {
	requireSolanaBatchFitsCallLimit(t, batchName, len(tests), solanaMaxWriteCallsPerBatch, "writes")

	if parallelEnabled {
		t.Parallel()
	}
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetTestConfig(t, solanaRegressionConfigPath))

	SolanaWriteFailsTest(t, testEnv, batchName, tests)
}

func Test_CRE_V2_Solana_WriteReport_Invalid_Receiver_Regression(t *testing.T) {
	runSolanaWriteNegativeTestSuite(t, "invalid-receiver", solanaNegativeTestsWriteReportInvalidReceiver)
}

func Test_CRE_V2_Solana_WriteReport_Invalid_Payload_Regression(t *testing.T) {
	runSolanaWriteNegativeTestSuite(t, "invalid-payload", solanaNegativeTestsWriteReportInvalidPayload)
}

//////////////////////////////////////////////////////
// READ NEGATIVE TESTS
//////////////////////////////////////////////////////

// ...Function variables must literally match the names of the switch-case statements in the
// workflow (solana/solread-negative/main.go).
const (
	solanaGetAccountInfoInvalidAccount      = "GetAccountInfoWithOpts - invalid account"
	solanaGetAccountInfoNonExistentAccount  = "GetAccountInfoWithOpts - non-existent account"
	solanaGetAccountInfoInvalidCommitment   = "GetAccountInfoWithOpts - invalid commitment"
	solanaGetAccountInfoInvalidEncoding     = "GetAccountInfoWithOpts - invalid encoding"
	solanaGetBalanceInvalidAddress          = "GetBalance - invalid address"
	solanaGetBalanceNonExistentAccount      = "GetBalance - non-existent account"
	solanaGetMultipleAccountsInvalidAccount = "GetMultipleAccountsWithOpts - invalid account"
	solanaGetMultipleAccountsTooManyItems   = "GetMultipleAccountsWithOpts - too many accounts"
	solanaGetProgramAccountsInvalidProgram  = "GetProgramAccounts - invalid program"
	solanaGetProgramAccountsNonProgram      = "GetProgramAccounts - non-program account"
	solanaGetBlockInvalidSlot               = "GetBlock - invalid slot"
	solanaGetSlotHeightInvalidCommitment    = "GetSlotHeight - invalid commitment"
	solanaGetTransactionInvalidSignature    = "GetTransaction - invalid signature"
	solanaGetTransactionUnknownSignature    = "GetTransaction - unknown signature"
	solanaGetSignatureStatusesInvalidSig    = "GetSignatureStatuses - invalid signature"
	solanaGetSignatureStatusesTooManyItems  = "GetSignatureStatuses - too many signatures"
	solanaGetFeeForMessageInvalidMessage    = "GetFeeForMessage - invalid message"
)

const (
	expectedSolanaNilAddress        = "address can't be nil"
	expectedSolanaInvalidPublicKey  = "invalid public key: got %d bytes, expected 32"
	expectedSolanaNilSignature      = "signature can't be nil"
	expectedSolanaInvalidSignature  = "invalid signature: got %d bytes, expected 64"
	expectedSolanaUnknownCommitment = "unknown commitment type"
	expectedSolanaUnknownEncoding   = "unknown encoding type"
	// the batch limiter names the setting it enforces, the same way the EVM log query limiter does
	expectedSolanaBatchItemLimit = "BatchItemLimit"
	// acceptErrorOrEmpty marks cases whose request is well-formed, so the capability is free to
	// answer with either an error or an empty response.
	acceptErrorOrEmpty = ""
)

// solanaBatchItemLimitExceeded is one more than PerWorkflow.ChainRead.Solana.BatchItemLimit (100).
const solanaBatchItemLimitExceeded = "101"

// invalidSolanaEnumValue is outside the range of every commitment/encoding enum.
const invalidSolanaEnumValue = "99"

var solanaNegativeTestsGetAccountInfo = []solanaNegativeTest{
	// GetAccountInfoWithOpts - an unset account reaches the capability as a nil field, a
	// malformed one fails the length check: two branches, two cases.
	{"unset account", "", solanaGetAccountInfoInvalidAccount, expectedSolanaNilAddress},
	{"31 bytes (short) key", solanaShortKey, solanaGetAccountInfoInvalidAccount, invalidPublicKeyError(31)},
	// A well-formed key of an account that was never created is a valid request, so the answer is
	// expected to be an error or an empty account rather than account data.
	{"non-existent account", nonExistentSolanaAccount, solanaGetAccountInfoNonExistentAccount, acceptErrorOrEmpty},
	// This is the only read that carries both opts enums, so it covers them for all of them.
	{"unknown commitment", invalidSolanaEnumValue, solanaGetAccountInfoInvalidCommitment, expectedSolanaUnknownCommitment},
	{"unknown encoding", invalidSolanaEnumValue, solanaGetAccountInfoInvalidEncoding, expectedSolanaUnknownEncoding},
}

var solanaNegativeTestsGetBalance = []solanaNegativeTest{
	// GetBalance - its own converter, so the address branches are worth covering again here
	{"unset address", "", solanaGetBalanceInvalidAddress, expectedSolanaNilAddress},
	{"31 bytes (short) key", solanaShortKey, solanaGetBalanceInvalidAddress, invalidPublicKeyError(31)},
	// Solana reports zero lamports for an account that was never funded.
	{"non-existent account", nonExistentSolanaAccount, solanaGetBalanceNonExistentAccount, acceptErrorOrEmpty},
}

var solanaNegativeTestsGetMultipleAccounts = []solanaNegativeTest{
	// GetMultipleAccountsWithOpts - a repeated field, validated per element, so one malformed key
	// is enough to reject the whole batch
	{"31 bytes (short) key in the batch", solanaShortKey, solanaGetMultipleAccountsInvalidAccount, invalidPublicKeyError(31)},
	// Every key of this batch is valid, so it can only be rejected by the batch item limiter.
	{"over the batch item limit", solanaBatchItemLimitExceeded, solanaGetMultipleAccountsTooManyItems, expectedSolanaBatchItemLimit},
}

var solanaNegativeTestsGetProgramAccounts = []solanaNegativeTest{
	// GetProgramAccounts - program that is not a usable public key
	{"unset program", "", solanaGetProgramAccountsInvalidProgram, expectedSolanaNilAddress},
	{"31 bytes (short) key", solanaShortKey, solanaGetProgramAccountsInvalidProgram, invalidPublicKeyError(31)},
	// Nothing owns accounts under a key that is not a program.
	{"account that owns nothing", nonExistentSolanaAccount, solanaGetProgramAccountsNonProgram, acceptErrorOrEmpty},
}

var solanaNegativeTestsGetBlockAndSlotHeight = []solanaNegativeTest{
	// GetBlock - a slot the local validator has never produced
	{"slot never produced", "18446744073709551615", solanaGetBlockInvalidSlot, acceptErrorOrEmpty},
	// GetSlotHeight takes nothing but a commitment, so that is all there is to get wrong.
	{"unknown commitment", invalidSolanaEnumValue, solanaGetSlotHeightInvalidCommitment, expectedSolanaUnknownCommitment},
}

var solanaNegativeTestsGetTransaction = []solanaNegativeTest{
	// GetTransaction - signature that is not a usable 64-byte signature
	{"unset signature", "", solanaGetTransactionInvalidSignature, expectedSolanaNilSignature},
	{"63 bytes (short) signature", solanaShortSignature, solanaGetTransactionInvalidSignature, invalidSignatureError(63)},
	// A well-formed signature of a transaction that was never submitted.
	{"unknown signature", solanaUnknownSignature, solanaGetTransactionUnknownSignature, acceptErrorOrEmpty},
}

var solanaNegativeTestsGetSignatureStatuses = []solanaNegativeTest{
	// GetSignatureStatuses - the repeated counterpart of GetTransaction's signature check
	{"63 bytes (short) signature in the batch", solanaShortSignature, solanaGetSignatureStatusesInvalidSig, invalidSignatureError(63)},
	{"over the batch item limit", solanaBatchItemLimitExceeded, solanaGetSignatureStatusesTooManyItems, expectedSolanaBatchItemLimit},
}

var solanaNegativeTestsGetFeeForMessage = []solanaNegativeTest{
	// GetFeeForMessage - anything that is not an encoded Solana message. The capability passes
	// the string through, so every rejection comes from the RPC and takes the same path.
	{"unset message", "", solanaGetFeeForMessageInvalidMessage, acceptErrorOrEmpty},
	{"not an encoded message", "not-a-base64-message", solanaGetFeeForMessageInvalidMessage, acceptErrorOrEmpty},
}

// Read cases are grouped into batches that each run in a single workflow execution, and every
// batch must stay within solanaMaxReadCallsPerBatch - when adding cases, either keep a batch
// under it or split the batch.
var (
	// 10 calls
	solanaReadBatchAccountCalls = concatSolanaTests(
		solanaNegativeTestsGetAccountInfo,
		solanaNegativeTestsGetBalance,
		solanaNegativeTestsGetMultipleAccounts,
	)
	// 12 calls
	solanaReadBatchProgramBlockAndTxCalls = concatSolanaTests(
		solanaNegativeTestsGetProgramAccounts,
		solanaNegativeTestsGetBlockAndSlotHeight,
		solanaNegativeTestsGetTransaction,
		solanaNegativeTestsGetSignatureStatuses,
		solanaNegativeTestsGetFeeForMessage,
	)
)

// SolanaReadFailsTest deploys the Solana read negative workflow once for a whole batch of cases.
// The workflow runs every case in a single execution and checks each of them against its expected
// outcome, so the test waits for one "batch passed" log. A failing case fails the batch with a log
// that names the case and what it actually got back.
func SolanaReadFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, batchName string, tests []solanaNegativeTest) {
	testLogger := framework.L

	cases := toSolanaReadCases(tests)

	solChain := mustSolanaChainInEnv(t, testEnv)
	chainSelector := solChain.ChainSelector()

	userLogsCh, baseMessageCh := startSolanaChipSink(t, testLogger)

	testLogger.Info().Msgf("Creating Solana Read Fail workflow configuration for batch '%s' (%d cases)...", batchName, len(cases))
	workflowConfig := solread_negative_config.Config{
		ChainSelector: chainSelector,
		BatchName:     batchName,
		Cases:         cases,
	}

	workflowName := fmt.Sprintf("sol-read-fail-workflow-%d-%04d", chainSelector, rand.Intn(10000))
	workflowID := t_helpers.CompileAndDeployWorkflow(
		t,
		testEnv,
		testLogger,
		workflowName,
		&workflowConfig,
		solanaReadNegativeWorkflowFile,
	)

	t_helpers.WatchWorkflowLogs(
		t,
		testLogger,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		solanaReadBatchPassedLog(batchName, len(cases)),
		2*time.Minute,
		t_helpers.WithUserLogWorkflowID(workflowID),
	)
	testLogger.Info().Msgf("Solana Read Fail test successfully completed for batch '%s'", batchName)
}

// runSolanaReadNegativeTestSuite runs one batch of Solana read negative cases.
func runSolanaReadNegativeTestSuite(t *testing.T, batchName string, tests []solanaNegativeTest) {
	requireSolanaBatchFitsCallLimit(t, batchName, len(tests), solanaMaxReadCallsPerBatch, "chain reads")

	if parallelEnabled {
		t.Parallel()
	}
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetTestConfig(t, solanaRegressionConfigPath))

	SolanaReadFailsTest(t, testEnv, batchName, tests)
}

func Test_CRE_V2_Solana_Read_Account_Calls_Regression(t *testing.T) {
	runSolanaReadNegativeTestSuite(t, "account-calls", solanaReadBatchAccountCalls)
}

func Test_CRE_V2_Solana_Read_Program_Block_And_Tx_Calls_Regression(t *testing.T) {
	runSolanaReadNegativeTestSuite(t, "program-block-and-tx-calls", solanaReadBatchProgramBlockAndTxCalls)
}

//////////////////////////////////////////////////////
// LOG TRIGGER NEGATIVE TESTS
//////////////////////////////////////////////////////

// ...Function variables must literally match the names of the switch-case statements in the
// workflow (solana/sollogtrigger-negative/main.go). Every case here is rejected while the
// workflow engine registers the trigger, so the failure shows up as an engine init error.
//
// Unlike reads and writes these cases cannot be batched: a rejected registration fails the whole
// engine, so a workflow can only ever demonstrate one rejected filter.
const (
	solanaLogTriggerInvalidAddress  = "LogTrigger - invalid address"
	solanaLogTriggerEmptyEventName  = "LogTrigger - empty event name"
	solanaLogTriggerEmptyFilterName = "LogTrigger - empty filter name"
	solanaLogTriggerEmptyIDL        = "LogTrigger - empty IDL"
	solanaLogTriggerConflictSubkeys = "LogTrigger - conflicting subkey filters"
	solanaLogTriggerInvalidCPIDest  = "LogTrigger - invalid CPI destination address"
)

const (
	expectedSolanaInvalidTriggerAddress = "invalid address length: expected 32 bytes, got %d"
	expectedSolanaEmptyEventName        = "event name cannot be empty"
	expectedSolanaEmptyFilterName       = "filter name cannot be empty"
	expectedSolanaEmptyIDL              = "event idl json cannot be empty"
	expectedSolanaConflictingSubkeys    = "subkey 0 has conflicting equality filters"
	expectedSolanaInvalidCPIDest        = "invalid cpi filter destination address length"
)

// solanaLogReadTestProgramID is the log_read_test program deployed by the Solana topology, also
// used by the Solana log trigger smoke test. The negative cases only need a well-formed address,
// but using a real program keeps every field other than the one under test realistic.
const solanaLogReadTestProgramID = "J1zQwrBNBngz26jRPNWsUSZMHJwBwpkoDitXRV95LdK4"

// Each of these cases costs a whole workflow deployment, so the table is kept to exactly one
// case per branch of the capability's filter validation.
var solanaNegativeTestsLogTrigger = []solanaNegativeTest{
	// LogTrigger - address that is not a 32-byte public key
	{"31 bytes (short) address", solanaShortKey, solanaLogTriggerInvalidAddress, invalidTriggerAddressError(31)},
	// LogTrigger - filter fields that must be set
	{"no event name", "", solanaLogTriggerEmptyEventName, expectedSolanaEmptyEventName},
	{"no filter name", "", solanaLogTriggerEmptyFilterName, expectedSolanaEmptyFilterName},
	{"no contract IDL", "", solanaLogTriggerEmptyIDL, expectedSolanaEmptyIDL},
	// LogTrigger - a subkey cannot be required to equal two different values at once
	{"two equality filters on one subkey", "", solanaLogTriggerConflictSubkeys, expectedSolanaConflictingSubkeys},
	// LogTrigger - the CPI destination address is a public key too
	{"31 bytes (short) CPI address", solanaShortKey, solanaLogTriggerInvalidCPIDest, expectedSolanaInvalidCPIDest},
}

// SolanaLogTriggerFailsTest deploys the Solana log trigger negative workflow for a single case.
// The trigger is rejected during registration, so the workflow engine never initialises and the
// failure is asserted on the engine init error rather than on a user log.
func SolanaLogTriggerFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, solanaNegativeTest solanaNegativeTest) {
	testLogger := framework.L

	solChain := mustSolanaChainInEnv(t, testEnv)
	chainSelector := solChain.ChainSelector()

	userLogsCh, baseMessageCh := startSolanaChipSink(t, testLogger)
	// drain user logs channel in the background, we are not asserting anything on it
	t_helpers.IgnoreUserLogs(t.Context(), userLogsCh)

	testLogger.Info().Msg("Creating Solana LogTrigger Fail workflow configuration...")
	workflowConfig := sollogtrigger_negative_config.Config{
		ChainSelector:  chainSelector,
		FunctionToTest: solanaNegativeTest.functionToTest,
		InvalidInput:   solanaNegativeTest.invalidInput,
		ProgramID:      solanago.MustPublicKeyFromBase58(solanaLogReadTestProgramID),
	}

	workflowName := fmt.Sprintf("sol-logtrigger-fail-workflow-%d-%04d", chainSelector, rand.Intn(10000))
	workflowID := t_helpers.CompileAndDeployWorkflow(
		t,
		testEnv,
		testLogger,
		workflowName,
		&workflowConfig,
		solanaLogTriggerNegativeWorkflowFile,
		// the log trigger lives on the capabilities DON, so it needs workflow artifacts copied there
		t_helpers.WithArtifactCopyDONTypes(crelib.WorkflowDON, crelib.CapabilitiesDON),
	)

	_ = t_helpers.WatchBaseMessages(
		t,
		testLogger,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		2*time.Minute,
		t_helpers.WithBaseMessageWorkflowID(workflowID),
		t_helpers.WithBaseMessageLabelContains("err", solanaNegativeTest.expectedError),
	)
	testLogger.Info().Msgf("Solana LogTrigger Fail test successfully completed for test case %s", solanaNegativeTest.name)
}

func Test_CRE_V2_Solana_LogTrigger_Invalid_Filter_Regression(t *testing.T) {
	for _, tCase := range solanaNegativeTestsLogTrigger {
		testName := fmt.Sprintf(solanaTestNameTemplate, tCase.functionToTest, tCase.name)
		t.Run(testName, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetTestConfig(t, solanaRegressionConfigPath))

			SolanaLogTriggerFailsTest(t, testEnv, tCase)
		})
	}
}

// ─── shared helpers ───────────────────────────────────────────────────────────

// toSolanaReadCases converts the test tables into the workflow's case list. A case with no
// expected error accepts an error or an empty response, see solanaNegativeTest.
func toSolanaReadCases(tests []solanaNegativeTest) []solread_negative_config.Case {
	cases := make([]solread_negative_config.Case, 0, len(tests))
	for _, test := range tests {
		outcome := solread_negative_config.OutcomeError
		if test.expectedError == acceptErrorOrEmpty {
			outcome = solread_negative_config.OutcomeErrorOrEmpty
		}

		cases = append(cases, solread_negative_config.Case{
			Name:           test.name,
			FunctionToTest: test.functionToTest,
			InvalidInput:   test.invalidInput,
			Outcome:        outcome,
			ExpectedError:  test.expectedError,
		})
	}
	return cases
}

func toSolanaWriteCases(tests []solanaNegativeTest) []solwrite_negative_config.Case {
	cases := make([]solwrite_negative_config.Case, 0, len(tests))
	for _, test := range tests {
		cases = append(cases, solwrite_negative_config.Case{
			Name:           test.name,
			FunctionToTest: test.functionToTest,
			InvalidInput:   test.invalidInput,
			ExpectedError:  test.expectedError,
		})
	}
	return cases
}

// requireSolanaBatchFitsCallLimit fails the test before the environment is booted, so that a
// batch grown past what one workflow execution may call is reported immediately instead of as a
// call-limit error two minutes into the run.
func requireSolanaBatchFitsCallLimit(t *testing.T, batchName string, cases, limit int, calls string) {
	t.Helper()

	require.LessOrEqualf(t, cases, limit,
		"batch %q has %d cases, which is over the %d %s a single workflow execution is allowed - split the batch",
		batchName, cases, limit, calls)
}

func concatSolanaTests(tables ...[]solanaNegativeTest) []solanaNegativeTest {
	var all []solanaNegativeTest
	for _, table := range tables {
		all = append(all, table...)
	}
	return all
}

// batchNeedsNonExecutableAccount reports whether the batch contains the case that needs an
// existing, non-executable account, so that funding one is skipped for the batches that do not.
func batchNeedsNonExecutableAccount(cases []solwrite_negative_config.Case) bool {
	for _, testCase := range cases {
		if testCase.FunctionToTest == solanaWriteReportNonExecutableReceiver {
			return true
		}
	}
	return false
}

// solanaReadBatchPassedLog and solanaWriteBatchPassedLog must match the log the workflows emit
// once every case of a batch passed.
func solanaReadBatchPassedLog(batchName string, cases int) string {
	return fmt.Sprintf("Solana read negative batch '%s' passed (%d cases)", batchName, cases)
}

func solanaWriteBatchPassedLog(batchName string, cases int) string {
	return fmt.Sprintf("Solana write negative batch '%s' passed (%d cases)", batchName, cases)
}

// startSolanaChipSink boots a per-test CHiP sink and returns the channels the workflow logs and
// engine messages are published to.
func startSolanaChipSink(t *testing.T, testLogger zerolog.Logger) (chan *workflowevents.UserLogs, chan *commonevents.BaseMessage) {
	t.Helper()

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))

	t.Cleanup(func() {
		// can't use t.Context() here because it will have been cancelled before the cleanup function is called
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(
			ctx,
			server,
			userLogsCh,
			baseMessageCh,
		)
	})

	return userLogsCh, baseMessageCh
}

func invalidPublicKeyError(gotBytes int) string {
	return fmt.Sprintf(expectedSolanaInvalidPublicKey, gotBytes)
}

func invalidSignatureError(gotBytes int) string {
	return fmt.Sprintf(expectedSolanaInvalidSignature, gotBytes)
}

func invalidTriggerAddressError(gotBytes int) string {
	return fmt.Sprintf(expectedSolanaInvalidTriggerAddress, gotBytes)
}

// mustSolanaChainInEnv returns the single Solana chain of the test environment.
func mustSolanaChainInEnv(t *testing.T, testEnv *ttypes.TestEnvironment) *solana.Blockchain {
	t.Helper()

	var solChain *solana.Blockchain
	for _, blockchain := range testEnv.CreEnvironment.Blockchains {
		if !blockchain.IsFamily(chainselectors.FamilySolana) {
			continue
		}
		require.IsType(t, &solana.Blockchain{}, blockchain, "expected Solana blockchain type")
		// we assume we always have just 1 Solana chain
		solChain = blockchain.(*solana.Blockchain)
		break
	}
	require.NotNil(t, solChain, "Solana blockchain not found in test environment")

	return solChain
}

// mustGetSolanaContract reads a deployed Solana program/account address from the datastore.
func mustGetSolanaContract(t *testing.T, ds datastore.DataStore, chainSelector uint64, contractType datastore.ContractType) solanago.PublicKey {
	t.Helper()

	key := datastore.NewAddressRefKey(
		chainSelector,
		contractType,
		semver.MustParse("1.0.0"),
		ks_sol.DefaultForwarderQualifier,
	)
	contract, err := ds.Addresses().Get(key)
	require.NoError(t, err, "failed to get '%s' address from the datastore", contractType)

	return solanago.MustPublicKeyFromBase58(contract.Address)
}

// mustCreateFundedSolanaAccount funds a fresh address so that it exists on-chain as a
// system-owned, non-executable account.
func mustCreateFundedSolanaAccount(t *testing.T, testLogger zerolog.Logger, solChain blockchains.Blockchain) solanago.PublicKey {
	t.Helper()

	addresses, err := t_helpers.CreateAndFundAddressesSolana(t, testLogger, 1, solanaFundingAmount, solChain)
	require.NoError(t, err, "failed to create and fund a Solana address")
	require.Len(t, addresses, 1, "failed to create the correct number of Solana addresses")

	return addresses[0]
}
