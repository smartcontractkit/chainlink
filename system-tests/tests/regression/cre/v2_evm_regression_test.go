package cre

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"

	evm_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmread-negative/config"
	evm_write_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmwrite-negative/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// regression
const (
	// ...Function variables should literally match the name of the switch-case statements in the workflow (evm/evmread-negative/main.go)
	// in each case the corresponding evm capability function is called with the invalid input
	balanceAtFunction                          = "BalanceAt"
	expectedBalanceAtError                     = "balanceAt errored"
	callContractInvalidAddressToReadFunction   = "CallContract - invalid address to read"
	expectedCallContractInvalidAddressToRead   = "balances=&[+0]" // expecting empty array of balances
	callContractInvalidBRContractAddress       = "CallContract - invalid balance reader contract address"
	expectedCallContractInvalidContractAddress = "got expected empty response for invalid balance reader contract address"
	estimateGasInvalidToAddress                = "EstimateGas - invalid 'to' address"
	filterLogsInvalidAddresses                 = "FilterLogs - invalid addresses"
	expectedFilterLogsInvalidAddresses         = "got expected error or empty logs"
	filterLogsInvalidFromBlock                 = "FilterLogs - invalid FromBlock"
	expectedFilterLogsInvalidFromBlock         = "got expected error for FilterLogs with invalid fromBlock"
	filterLogsInvalidToBlock                   = "FilterLogs - invalid ToBlock"
	expectedFilterLogsInvalidToBlock           = "got expected error for FilterLogs with invalid toBlock"
	getTransactionByHashInvalidHash            = "GetTransactionByHash - invalid hash"
	getTransactionReceiptInvalidHash           = "GetTransactionReceipt - invalid hash"
	expectedGetTransactionByHashInvalidHash    = "not found"
	headerByNumberInvalidBlock                 = "HeaderByNumber - invalid block number"
	writeReportInvalidReceiver                 = "WriteReport - invalid receiver"
	expectedWriteReportInvalidReceiver         = "RECEIVER_CONTRACT_EXECUTION_STATUS_REVERTED"
	writeReportCorruptReceiverAddress          = "WriteReport - corrupt receiver address"
	expectedWriteReportCorruptReceiverAddress  = "failed to execute capability: received address is not 20 bytes long"
	writeReportInvalidGas                      = "WriteReport - invalid gas"
	expectedWriteReportInvalidGas              = "failed to execute capability"
	writeReportRandomTimestamps                = "WriteReport - random timestamps"
)

type evmNegativeTest struct {
	name           string
	invalidInput   string
	functionToTest string
	expectedError  string
}

var evmNegativeTestsBalanceAtInvalidAddress = []evmNegativeTest{
	// BalanceAt
	// TODO: Move BalanceAt to the top after fixing https://smartcontract-it.atlassian.net/browse/CRE-934
	{"a letter", "a", balanceAtFunction, expectedBalanceAtError},
	{"a symbol", "/", balanceAtFunction, expectedBalanceAtError},
	{"a number", "1", balanceAtFunction, expectedBalanceAtError},
	{"empty hex", "0x", balanceAtFunction, expectedBalanceAtError},
	{"cut hex", "0x0", balanceAtFunction, expectedBalanceAtError},
	{"short address", "0x123456789012345678901234567890123456789", balanceAtFunction, expectedBalanceAtError},
	{"long address", "0x12345678901234567890123456789012345678901", balanceAtFunction, expectedBalanceAtError},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", balanceAtFunction, expectedBalanceAtError},
}

var evmNegativeTestsCallContractInvalidAddressToRead = []evmNegativeTest{
	// CallContract - invalid address to read
	// Some invalid inputs are skipped (empty, symbols, "0x", "0x0") as they may map to the zero address and return a balance instead of empty.
	{"a letter", "a", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"a number", "1", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"short address", "0x123456789012345678901234567890123456789", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"long address", "0x12345678901234567890123456789012345678901", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
}

var evmNegativeTestsCallContractInvalidBalanceReaderContract = []evmNegativeTest{
	// CallContract - invalid balance reader contract address
	{"empty", "", callContractInvalidBRContractAddress, "EVM error OpcodeNotFound"}, // equivalent to "0x", "0x0", we do not care if anything but contract may be at this address
	{"a letter", "a", callContractInvalidBRContractAddress, "EVM error PrecompileError"},
	{"a symbol", "/", callContractInvalidBRContractAddress, "EVM error OpcodeNotFound"},
	{"a number", "1", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	{"short address", "0x123456789012345678901234567890123456789", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	{"long address", "0x12345678901234567890123456789012345678901", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
}

var evmNegativeTestsEstimateGasInvalidToAddress = []evmNegativeTest{
	// EstimateGas - invalid 'to' address
	// do not use 1, short, long addresses because common.Address will convert them to a valid address
	// also it does not make sense to use invalid CallMsg.Data because any bytes will be correctly processed
	{"empty", "", estimateGasInvalidToAddress, "EVM error StackUnderflow"},
	{"a letter", "a", estimateGasInvalidToAddress, "EVM error PrecompileError"},
	{"a symbol", "/", estimateGasInvalidToAddress, "EVM error StackUnderflow"},
	{"not authored contract", "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512", estimateGasInvalidToAddress, "execution reverted"},
	{"cut hex", "0x", estimateGasInvalidToAddress, "EVM error StackUnderflow"}, // equivalent to "0x0"
}

var evmNegativeTestsFilterLogsWithInvalidAddress = []evmNegativeTest{
	// FilterLogs - invalid addresses.
	// Note: evm.FilterLogs does not validate addresses if they are correctly formatted
	// Since filtering is applied over blocks/logs — if no logs match, the result is just empty, which is a valid response
	// do not use empty, 1, short, long addresses because common.Address will convert them to a valid address
	{"a letter", "a", filterLogsInvalidAddresses, expectedFilterLogsInvalidAddresses},
	{"a number", "1", filterLogsInvalidAddresses, expectedFilterLogsInvalidAddresses},
	{"a symbol", "/", filterLogsInvalidAddresses, expectedFilterLogsInvalidAddresses},
	{"short address", "0x123456789012345678901234567890123456789", filterLogsInvalidAddresses, expectedFilterLogsInvalidAddresses},
	{"long address", "0x12345678901234567890123456789012345678901", filterLogsInvalidAddresses, expectedFilterLogsInvalidAddresses},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", filterLogsInvalidAddresses, expectedFilterLogsInvalidAddresses},
}

var evmNegativeTestsFilterLogsWithInvalidFromBlock = []evmNegativeTest{
	// FilterLogs - invalid FromBlock/ToBlock values
	// Distance between blocks should not be more than 100
	{"negative number", "-1", filterLogsInvalidFromBlock, "block number -1 is not supported"},
	{"zero", "0", filterLogsInvalidFromBlock, "block number 0 is not supported"},
	{"very large number", "9223372036854775808", filterLogsInvalidFromBlock, "is not an int64"}, // int64 max + 1
	{"non-numeric string", "abc", filterLogsInvalidFromBlock, "toBlock 150 is less than fromBlock"},
	{"empty string", "", filterLogsInvalidFromBlock, "toBlock 150 is less than fromBlock"},
	{"decimal", "100.5", filterLogsInvalidFromBlock, "toBlock 150 is less than fromBlock"},
	{"fromBlock greater than toBlock by more than 100", "49", filterLogsInvalidFromBlock, "PerWorkflow.ChainRead.LogQueryBlockLimit limited for workflow"}, // toBlock is 150, so distance is 100+
}

var evmNegativeTestsFilterLogsWithInvalidToBlock = []evmNegativeTest{
	// FilterLogs - invalid toBlock values
	// Distance between blocks should not be more than 100
	{"negative number", "-1", filterLogsInvalidToBlock, "block number -1 is not supported"},
	{"zero", "0", filterLogsInvalidToBlock, "block number 0 is not supported"},
	{"less then FromBlock", "1", filterLogsInvalidToBlock, "toBlock 1 is less than fromBlock"},
	{"very large number", "9223372036854775808", filterLogsInvalidToBlock, "is not an int64"}, // int64 max + 1
	{"non-numeric string", "abc", filterLogsInvalidToBlock, "PerWorkflow.ChainRead.LogQueryBlockLimit limited for workflow"},
	{"empty string", "", filterLogsInvalidToBlock, "PerWorkflow.ChainRead.LogQueryBlockLimit limited for workflow"}, // equivalent to "current block"
	{"decimal", "100.5", filterLogsInvalidToBlock, "PerWorkflow.ChainRead.LogQueryBlockLimit limited for workflow"},
	{"toBlock greater than fromBlock by more than 100", "103", filterLogsInvalidToBlock, "PerWorkflow.ChainRead.LogQueryBlockLimit limited for workflow"}, // fromBlock is 2
}

var evmNegativeTestsGetTransactionByHashInvalidHash = []evmNegativeTest{
	// GetTransactionByHash - invalid hash (requires 32 bytes)
	{"empty", "", getTransactionByHashInvalidHash, "hash can't be nil"}, // equivalent to whitespace " "
	{"a symbol", ";", getTransactionByHashInvalidHash, "hash can't be nil"},
	{"a char", "0xz", getTransactionByHashInvalidHash, "hash can't be nil"},         // equivalent to any alfa-numeric string/character
	{"null-0-like hex", "0x", getTransactionByHashInvalidHash, "hash can't be nil"}, // equivalent to "0x0", empty
	{"31 bytes (short) non-0x-prefixed", "12345678901234567890123456789012345678901234567890123456789012", getTransactionByHashInvalidHash, "got 31 bytes, expected 32"},
	{"33 bytes (long) non-0x-prefixed", "12345678901234567890123456789012345678901234567890123456789012345", getTransactionByHashInvalidHash, "got 33 bytes, expected 32"},
	{"malformed (non-hex) correct length", "0x123gggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", getTransactionByHashInvalidHash, "got 2 bytes, expected 32"}, // produces x01#
	{"short hash", "0x647b7f17f9edba01d1f75ce071d0bc10173bc66b5d072f28b644275bf13bb99", getTransactionByHashInvalidHash, "RPC call failed: not found"},
	{"non-existent hash", "0x1234567890123456789012345678901234567890123456789012345678901234", getTransactionByHashInvalidHash, "RPC call failed: not found"},
}

var evmNegativeTestsGetTransactionReceiptInvalidHash = []evmNegativeTest{
	// GetTransactionReceipt - invalid hash (requires 32 bytes)
	{"empty", "", getTransactionReceiptInvalidHash, "hash can't be nil"}, // equivalent to whitespace " "
	{"a symbol", ";", getTransactionReceiptInvalidHash, "hash can't be nil"},
	{"a char", "0xz", getTransactionReceiptInvalidHash, "hash can't be nil"},         // equivalent to any alfa-numeric string/character
	{"null-0-like hex", "0x", getTransactionReceiptInvalidHash, "hash can't be nil"}, // equivalent to "0x0", empty
	{"31 bytes (short) non-0x-prefixed", "12345678901234567890123456789012345678901234567890123456789012", getTransactionReceiptInvalidHash, "got 31 bytes, expected 32"},
	{"33 bytes (long) non-0x-prefixed", "12345678901234567890123456789012345678901234567890123456789012345", getTransactionReceiptInvalidHash, "got 33 bytes, expected 32"},
	{"malformed (non-hex) correct length", "0x123gggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", getTransactionReceiptInvalidHash, "got 2 bytes, expected 32"}, // produces x01#
	{"short hash", "0x647b7f17f9edba01d1f75ce071d0bc10173bc66b5d072f28b644275bf13bb99", getTransactionReceiptInvalidHash, "RPC call failed: not found"},
	{"non-existent hash", "0x1234567890123456789012345678901234567890123456789012345678901234", getTransactionReceiptInvalidHash, "RPC call failed: not found"},
}

var evmNegativeTestsHeaderByNumberInvalidBlock = []evmNegativeTest{
	// HeaderByNumber - invalid block number
	// empty, non-numeric string, decimal will return nil, when parsed to big.Int,
	// nil is a valid param for searching the latest block, and won't error.
	{"negative number", "-1", headerByNumberInvalidBlock, "block number -1 is not supported"},
	{"zero", "0", headerByNumberInvalidBlock, "block number 0 is not supported"},
	{"int overflownumber", "9223372036854775808", headerByNumberInvalidBlock, "is not an int64"},             // int64 max + 1
	{"not existing block)", "9223372036854775807", headerByNumberInvalidBlock, "RPC call failed: not found"}, // int64 max
}

func EVMReadFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, evmNegativeTest evmNegativeTest) {
	testLogger := framework.L
	const workflowFileLocation = "./evm/evmread-negative/main.go"
	enabledChains := t_helpers.GetEVMEnabledChains(t, testEnv)

	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		chainID := bcOutput.BlockchainOutput.ChainID
		chainSelector := bcOutput.ChainSelector
		creEnvironment := testEnv.CreEnvironment
		if _, ok := enabledChains[chainID]; !ok {
			testLogger.Info().Msgf("Skipping chain %s as it is not enabled for EVM Read workflow test", chainID)
			continue
		}

		testLogger.Info().Msgf("Deploying additional contracts to chain %s (%d)", chainID, chainSelector)
		readBalancesAddress, rbOutput, rbErr := crecontracts.DeployReadBalancesContract(testLogger, chainSelector, creEnvironment)
		require.NoError(t, rbErr, "failed to deploy Read Balances contract on chain %d", chainSelector)
		crecontracts.MergeAllDataStores(creEnvironment, rbOutput, rbOutput)

		listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)
		testLogger.Info().Msg("Creating EVM Read Fail workflow configuration...")
		workflowConfig := evm_negative_config.Config{
			ChainSelector:  bcOutput.ChainSelector,
			FunctionToTest: evmNegativeTest.functionToTest,
			InvalidInput:   evmNegativeTest.invalidInput,
			BalanceReader: evm_negative_config.BalanceReader{
				BalanceReaderAddress: readBalancesAddress,
			},
		}
		workflowName := fmt.Sprintf("evm-read-fail-workflow-%s-%04d", chainID, rand.Intn(10000))
		t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

		expectedError := evmNegativeTest.expectedError
		timeout := 2 * time.Minute
		err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedError, testLogger, messageChan, kafkaErrChan, timeout)
		require.NoError(t, err, "EVM Read Fail test failed")
		testLogger.Info().Msg("EVM Read Fail test successfully completed")
	}
}

//////////////////////////////////////////////////////
// WRITE REPORT NEGATIVE TESTS
//////////////////////////////////////////////////////

var evmNegativeTestsWriteReportInvalidReceiver = []evmNegativeTest{
	// WriteReport - invalid receiver
	// symbols, numbers, 0-addresses are skipped as they are equivalent values for common.Address
	{"empty", "", writeReportInvalidReceiver, expectedWriteReportInvalidReceiver},
	{"short address", "0x123456789012345678901234567890123456789", writeReportInvalidReceiver, expectedWriteReportInvalidReceiver},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef1234567", writeReportInvalidReceiver, expectedWriteReportInvalidReceiver},
	{"not a contract", "0x9b516F6741Dd1889A3Db4DC276aD349F0DC403C8", writeReportInvalidReceiver, expectedWriteReportInvalidReceiver},
}

var evmNegativeTestsWriteReportCorruptReceiverAddress = []evmNegativeTest{
	// WriteReport - corrupt receiver address
	// malformed values
	{"empty", "", writeReportCorruptReceiverAddress, expectedWriteReportCorruptReceiverAddress},
	{"short address", "0x1234", writeReportCorruptReceiverAddress, expectedWriteReportCorruptReceiverAddress},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef1234567", writeReportCorruptReceiverAddress, expectedWriteReportCorruptReceiverAddress},
}

var evmNegativeTestsWriteReportInvalidGas = []evmNegativeTest{
	// WriteReport - corrupt receiver address
	// malformed values
	{"zero", "0", writeReportInvalidGas, expectedWriteReportInvalidGas},
	{"low", "100000", writeReportInvalidGas, expectedWriteReportInvalidGas},
	{"too high", "100000000000", writeReportInvalidGas, expectedWriteReportInvalidGas},
}

func EVMWriteFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, evmNegativeTest evmNegativeTest) {
	testLogger := framework.L
	const workflowFileLocation = "./evm/evmwrite-negative/main.go"
	enabledChains := t_helpers.GetEVMEnabledChains(t, testEnv)

	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		chainID := bcOutput.ChainID
		chainSelector := bcOutput.ChainSelector
		creEnvironment := testEnv.CreEnvironment
		if _, ok := enabledChains[strconv.FormatUint(chainID, 10)]; !ok {
			testLogger.Info().Msgf("Skipping chain %d as it is not enabled for EVM Read workflow test", chainID)
			continue
		}

		forwarderAddress, _, forwarderErr := crecontracts.FindAddressesForChain(creEnvironment.CldfEnvironment.ExistingAddresses, chainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
		require.NoError(t, forwarderErr, "failed to find Forwarder address for chain %d", chainSelector)

		workflowOwner := bcOutput.SethClient.MustGetRootKeyAddress()
		workflowName := fmt.Sprintf("evm-write-fail-workflow-%d-%04d", chainID, rand.Intn(10000))
		feedID := "018e16c38e000320000000000000000000000000000000000000000000000000" // 32 hex characters (16 bytes)
		dataFeedsCacheAddress := deployAndConfigureEVMContracts(t, testLogger, chainSelector, chainID, creEnvironment, workflowOwner, workflowName, feedID, forwarderAddress)

		listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)
		testLogger.Info().Msg("Creating EVM Write Regression workflow configuration...")
		workflowConfig := evm_write_negative_config.Config{
			FeedID:         feedID,
			ChainSelector:  bcOutput.ChainSelector,
			FunctionToTest: evmNegativeTest.functionToTest,
			InvalidInput:   evmNegativeTest.invalidInput,
			DataFeedsCache: evm_write_negative_config.DataFeedsCache{
				DataFeedsCacheAddress: dataFeedsCacheAddress,
			},
		}
		t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

		expectedError := evmNegativeTest.expectedError
		timeout := 2 * time.Minute
		err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedError, testLogger, messageChan, kafkaErrChan, timeout)
		require.NoError(t, err, "EVM Write Regression test failed")
		testLogger.Info().Msg("EVM Write Regression test successfully completed")
	}
}

func deployAndConfigureEVMContracts(t *testing.T, testLogger zerolog.Logger, chainSelector uint64, chainID uint64, creEnvironment *cre.Environment, workflowOwner common.Address, uniqueWorkflowName string, feedID string, forwarderAddress common.Address) common.Address {
	testLogger.Info().Msgf("Deploying additional contracts to chain %d (%d)", chainID, chainSelector)
	dfAddress, dfOutput, dfErr := crecontracts.DeployDataFeedsCacheContract(testLogger, chainSelector, creEnvironment)
	require.NoError(t, dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)
	crecontracts.MergeAllDataStores(creEnvironment, dfOutput, dfOutput)

	testLogger.Info().Msgf("Configuring Data Feeds Cache contract for EVM Write Regression test and feed ID %s", feedID)
	configInput := &cre.ConfigureDataFeedsCacheInput{
		CldEnv:                creEnvironment.CldfEnvironment,
		ChainSelector:         chainSelector,
		FeedIDs:               []string{feedID},
		Descriptions:          []string{"EVM Write Regression test"},
		DataFeedsCacheAddress: dfAddress,
		AdminAddress:          workflowOwner,
		AllowedSenders:        []common.Address{forwarderAddress},
		AllowedWorkflowNames:  []string{uniqueWorkflowName},
		AllowedWorkflowOwners: []common.Address{workflowOwner},
	}
	_, dfConfigErr := crecontracts.ConfigureDataFeedsCache(testLogger, configInput)
	require.NoError(t, dfConfigErr, "failed to configure Data Feeds Cache contract")
	testLogger.Info().Msg("Data Feeds Cache contract configured successfully.")

	return dfAddress
}
