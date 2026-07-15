package cre

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	thelpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	stellarReadWorkflowFile = "./stellar/stellarread/main.go"
)

// Expected user-log lines; values must match the stellar read workflow (stellarread/main.go).
// The ReadContract workflow runs every case in one trigger and emits a single aggregate line.
const (
	expectLogReadContractBatchOK = "Stellar ReadContract batch passed"
	expectLogLatestLedgerOK      = "Stellar read consensus succeeded"
)

// Fixed 32-byte / 2-byte payloads for echo_triple / echo_address (host builds ScVals; fixture echoes).
const (
	tripleAddrHex = "1111111111111111111111111111111111111111111111111111111111111111"
	tripleIDHex   = "2222222222222222222222222222222222222222222222222222222222222222"
	tripleRIDHex  = "abcd"
)

// executeStellarReadLatestLedgerTest exercises GetLatestLedger consensus (its own workflow).
func executeStellarReadLatestLedgerTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	stellarChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	t.Helper()
	lggr := framework.L
	workflowName := thelpers.UniqueStellarWorkflowName("stellar-read-workflow")
	workflowConfig := thelpers.StellarReadWorkflowConfig{
		ChainSelector:     stellarChain.ChainSelector(),
		WorkflowName:      workflowName,
		MinLedgerSequence: 1,
		CronSchedule:      thelpers.StellarTestCronSchedule(),
	}
	workflowID := thelpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, stellarReadWorkflowFile)
	thelpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, thelpers.WorkflowEngineInitErrorLog, expectLogLatestLedgerOK, thelpers.StellarWorkflowTimeout, thelpers.WithUserLogWorkflowID(workflowID))
	lggr.Info().Str("expected_log", expectLogLatestLedgerOK).Msg("Stellar GetLatestLedger capability test passed")
}

// executeStellarReadContractSmokeTest deploys the fixture and runs every happy-path ReadContract case in a single workflow trigger.
func executeStellarReadContractSmokeTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	stellarChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	t.Helper()
	lggr := framework.L
	fixtureID := thelpers.MustDeployStellarReadFixture(t, stellarChain)

	bytesN32 := "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	b, err := hex.DecodeString(bytesN32)
	require.NoError(t, err)

	vals := []uint32{1, 2, 3}
	vec := make(xdr.ScVec, len(vals))
	for i, u := range vals {
		vec[i], err = xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(u))
		require.NoError(t, err)
	}
	vals2 := []uint32{5, 6, 7}
	vec2 := make(xdr.ScVec, len(vals2))
	for i, u := range vals2 {
		vec2[i], err = xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(u))
		require.NoError(t, err)
	}
	steps := []thelpers.StellarReadContractStep{
		{Name: "get_bool", ContractID: fixtureID, Function: "get_bool", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBool, true))},
		{Name: "get_u32", ContractID: fixtureID, Function: "get_u32", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(42)))},
		{Name: "get_i32", ContractID: fixtureID, Function: "get_i32", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvI32, xdr.Int32(-7)))},
		{Name: "get_string", ContractID: fixtureID, Function: "get_string", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString("read-fixture")))},
		{Name: "get_bytes", ContractID: fixtureID, Function: "get_bytes", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes([]byte{0xde, 0xad, 0xbe, 0xef})))},
		{Name: "get_void", ContractID: fixtureID, Function: "get_void"},
		{Name: "get_self", ContractID: fixtureID, Function: "get_self"},

		// Wider scalar return types (distinct XDR encodings).
		{Name: "get_u64", ContractID: fixtureID, Function: "get_u64", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(1_000_000_000_000)))},
		{Name: "get_i64", ContractID: fixtureID, Function: "get_i64", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvI64, xdr.Int64(-1_000_000_000_000)))},
		{Name: "get_u128", ContractID: fixtureID, Function: "get_u128", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvU128, xdr.UInt128Parts{Hi: xdr.Uint64(0), Lo: xdr.Uint64(42)}))},
		// -42 as two's-complement i128: Hi=-1, Lo = 2^64 - 42.
		{Name: "get_i128", ContractID: fixtureID, Function: "get_i128", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvI128, xdr.Int128Parts{Hi: xdr.Int64(-1), Lo: xdr.Uint64(18446744073709551574)}))},
		{Name: "get_symbol", ContractID: fixtureID, Function: "get_symbol", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvSymbol, xdr.ScSymbol("fixture")))},
		{Name: "get_bytes_n32", ContractID: fixtureID, Function: "get_bytes_n32", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(make([]byte, 32))))},
		{Name: "get_bytes_n2", ContractID: fixtureID, Function: "get_bytes_n2", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes([]byte{0xab, 0xcd})))},

		// Collections / composite. map/struct assert consensus+decode (all nodes must agree on
		// identical XDR); exact expected XDR is omitted — entry ordering isn't worth encoding.
		{Name: "get_vec_u32", ContractID: fixtureID, Function: "get_vec_u32", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvVec, &vec))},
		{Name: "get_map", ContractID: fixtureID, Function: "get_map"},
		{Name: "get_sample", ContractID: fixtureID, Function: "get_sample"},

		// Echo / round-trip across arg types.
		{Name: "echo_u32", ContractID: fixtureID, Function: "echo_u32", ArgMode: "echo_u32", ArgU32A: 99, ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(99)))},
		{Name: "sum_u32", ContractID: fixtureID, Function: "sum_u32", ArgMode: "sum_u32", ArgU32A: 2, ArgU32B: 3, ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(5)))},
		{Name: "echo_bool", ContractID: fixtureID, Function: "echo_bool", ArgMode: "echo_bool", ArgBool: true, ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBool, true))},
		{Name: "echo_string", ContractID: fixtureID, Function: "echo_string", ArgMode: "echo_string", ArgString: "hello", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString("hello")))},
		{Name: "echo_bytes", ContractID: fixtureID, Function: "echo_bytes", ArgMode: "echo_bytes", ArgBytesHex: "cafebabe", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes([]byte{0xca, 0xfe, 0xba, 0xbe})))},
		{Name: "echo_i64", ContractID: fixtureID, Function: "echo_i64", ArgMode: "echo_i64", ArgI64: -1234567890123, ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvI64, xdr.Int64(-1234567890123)))},
		{Name: "echo_symbol", ContractID: fixtureID, Function: "echo_symbol", ArgMode: "echo_symbol", ArgSymbol: "echoed", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvSymbol, xdr.ScSymbol("echoed")))},
		{Name: "echo_vec_u32", ContractID: fixtureID, Function: "echo_vec_u32", ArgMode: "echo_vec_u32", ArgVecU32: []uint32{5, 6, 7}, ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvVec, &vec2))},
		{Name: "echo_bytes_n32", ContractID: fixtureID, Function: "echo_bytes_n32", ArgMode: "echo_bytes", ArgBytesHex: bytesN32, ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(b)))},
		{Name: "echo_bytes_n2", ContractID: fixtureID, Function: "echo_bytes_n2", ArgMode: "echo_bytes", ArgBytesHex: "beef", ExpectedResult: marshalScVal(xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes([]byte{0xbe, 0xef})))},
		// echo_triple / echo_address are shape-only (contract echoes the input; XDR-exact
		// asserts would require reconstructing composite/address ScVals for little gain).
		{Name: "echo_triple", ContractID: fixtureID, Function: "echo_triple", ArgMode: "echo_triple", ArgReceiverContractIDHex: tripleAddrHex, ArgWorkflowExecutionIDHex: tripleIDHex, ArgReportIDHex: tripleRIDHex},
		{Name: "echo_address", ContractID: fixtureID, Function: "echo_address", ArgMode: "echo_address", ArgReceiverContractIDHex: tripleAddrHex},
	}

	workflowName := thelpers.UniqueStellarWorkflowName("stellar-rc-batch")
	workflowConfig := thelpers.StellarReadWorkflowConfig{
		ChainSelector: stellarChain.ChainSelector(),
		WorkflowName:  workflowName,
		CronSchedule:  thelpers.StellarTestCronSchedule(),
		ReadKind:      thelpers.StellarReadKindReadContract,
		Cases:         steps,
	}
	workflowID := thelpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, stellarReadWorkflowFile)
	thelpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, thelpers.WorkflowEngineInitErrorLog, expectLogReadContractBatchOK, thelpers.StellarWorkflowTimeout, thelpers.WithUserLogWorkflowID(workflowID))
	lggr.Info().Int("cases", len(steps)).Str("expected_log", expectLogReadContractBatchOK).Msg("Stellar ReadContract capability test passed")
}

// setupStellarScenario provisions an isolated per-test context for a Stellar read scenario:
// its own funded keys (SetupTestEnvironmentWithPerTestKeys), the resolved Stellar chain, and a
// chip sink capturing that scenario's workflow logs. Call it at the top of each scenario subtest.
func setupStellarScenario(t *testing.T, tenv *configuration.TestEnvironment) (
	*configuration.TestEnvironment,
	blockchains.Blockchain,
	<-chan *workflowevents.UserLogs,
	<-chan *commonevents.BaseMessage,
) {
	t.Helper()
	scenarioEnv := thelpers.SetupTestEnvironmentWithPerTestKeys(t, tenv.TestConfig)
	scenarioChain := thelpers.MustStellarChainInEnv(t, scenarioEnv)

	logPath := thelpers.LogFilePath("stellar", t.Name())
	userLogsCh, baseMessageCh := thelpers.StartChipTestSinkWithLogging(t, logPath)
	framework.L.Info().Str("scenario", t.Name()).Str("log_file", logPath).Msg("Starting Stellar scenario")
	return scenarioEnv, scenarioChain, userLogsCh, baseMessageCh
}

// marshalScVal base64-encodes an already-constructed ScVal; panics on error (test-table data).
func marshalScVal(v xdr.ScVal, err error) string {
	if err != nil {
		panic(err)
	}
	b, err := xdr.MarshalBase64(v)
	if err != nil {
		panic(err)
	}
	return b
}
