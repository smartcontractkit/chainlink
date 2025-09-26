package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"

	evm_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmread-negative/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// regression
const (
	// find returned errors in the logs of the workflow
	balanceAtFunction                          = "BalanceAt"
	expectedBalanceAtError                     = "balanceAt errored"
	callContractInvalidAddressToReadFunction   = "CallContract - invalid address to read"
	expectedCallContractInvalidAddressToRead   = "balances=&[+0]" // expecting empty array of balances
	callContractInvalidBRContractAddress       = "CallContract - invalid balance reader contract address"
	expectedCallContractInvalidContractAddress = "got expected error for invalid balance reader contract address"
	estimateGasInvalidToAddress                = "EstimateGas - invalid 'to' address"
)

type evmNegativeTest struct {
	name           string
	invalidInput   string
	functionToTest string
	expectedError  string
}

var evmNegativeTests = []evmNegativeTest{
	// CallContract - invalid address to read
	// Some invalid inputs are skipped (empty, symbols, "0x", "0x0") as they may map to the zero address and return a balance instead of empty.
	{"a letter", "a", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"a number", "1", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"short address", "0x123456789012345678901234567890123456789", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"long address", "0x12345678901234567890123456789012345678901", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", callContractInvalidAddressToReadFunction, expectedCallContractInvalidAddressToRead},

	// CallContract - invalid balance reader contract address
	// TODO: Uncomment tests after https://smartcontract-it.atlassian.net/browse/CRE-943
	// "empty" will default to the 0-address which is valid but has no contract deployed, so we expect an error.
	// {"empty", "", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	// {"a letter", "a", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	// {"a symbol", "/", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	{"a number", "1", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	// {"empty hex", "0x", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress}, // we do not care if anything but contract may be at this address
	// {"cut hex", "0x0", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},  // we do not care if anything but contract may be at this address
	{"short address", "0x123456789012345678901234567890123456789", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	{"long address", "0x12345678901234567890123456789012345678901", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", callContractInvalidBRContractAddress, expectedCallContractInvalidContractAddress},

	// EstimateGas - invalid 'to' address
	// do not use 1, short, long addresses because common.Address will convert them to a valid address
	// also it does not make sense to use invalid CallMsg.Data because any bytes will be correctly processed
	{"empty", "", estimateGasInvalidToAddress, "EVM error StackUnderflow"},
	{"a letter", "a", estimateGasInvalidToAddress, "EVM error PrecompileError"},
	{"a symbol", "/", estimateGasInvalidToAddress, "EVM error StackUnderflow"},
	{"not authored contract", "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512", estimateGasInvalidToAddress, "execution reverted"},
	{"cut hex", "0x", estimateGasInvalidToAddress, "EVM error StackUnderflow"}, // equivalent to "0x0"

	// BalanceAt
	// TODO: Move BalanceAt to the top after fixing https://smartcontract-it.atlassian.net/browse/CRE-934
	{"empty", "", balanceAtFunction, expectedBalanceAtError},
	{"a letter", "a", balanceAtFunction, expectedBalanceAtError},
	{"a symbol", "/", balanceAtFunction, expectedBalanceAtError},
	{"a number", "1", balanceAtFunction, expectedBalanceAtError},
	{"empty hex", "0x", balanceAtFunction, expectedBalanceAtError},
	{"cut hex", "0x0", balanceAtFunction, expectedBalanceAtError},
	{"short address", "0x123456789012345678901234567890123456789", balanceAtFunction, expectedBalanceAtError},
	{"long address", "0x12345678901234567890123456789012345678901", balanceAtFunction, expectedBalanceAtError},
	{"invalid address", "0x1234567890abcdefg1234567890abcdef123456", balanceAtFunction, expectedBalanceAtError},
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
		workflowName := "evm-read-fail-workflow-" + chainID
		t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

		expectedError := evmNegativeTest.expectedError
		timeout := 75 * time.Second
		err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedError, testLogger, messageChan, kafkaErrChan, timeout)
		require.NoError(t, err, "EVM Read Fail test failed")
		testLogger.Info().Msg("EVM Read Fail test successfully completed")
	}
}
