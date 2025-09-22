package cre

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	evm_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread-negative/config"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/contracts"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// smoke
func ExecuteEVMReadTest(t *testing.T, testEnv *TestEnvironment) {
	lggr := framework.L
	const workflowFileLocation = "./evm/evmread/main.go"
	enabledChains := getEVMEnabledChains(t, testEnv)

	var workflowsWg sync.WaitGroup
	var successfulWorkflowRuns atomic.Int32
	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		chainID := bcOutput.BlockchainOutput.ChainID
		if _, ok := enabledChains[chainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM Read workflow test", chainID)
			continue
		}

		lggr.Info().Msg("Creating EVM Read workflow configuration...")
		workflowConfig := configureEVMReadWorkflow(t, lggr, bcOutput)
		workflowName := "evm-read-workflow-" + chainID
		compileAndDeployWorkflow(t, testEnv, lggr, workflowName, &workflowConfig, workflowFileLocation)

		workflowsWg.Add(1)
		go func(bcOutput *cre.WrappedBlockchainOutput) {
			defer workflowsWg.Done()
			validateWorkflowExecution(t, lggr, testEnv, bcOutput, workflowName, workflowConfig) //nolint:testifylint // TODO: consider refactoring
			successfulWorkflowRuns.Add(1)
		}(bcOutput)
	}

	// wait for all workflows to complete
	workflowsWg.Wait()
	require.Equal(t, len(enabledChains), int(successfulWorkflowRuns.Load()), "Not all workflows executed successfully")
}

func validateWorkflowExecution(t *testing.T, lggr zerolog.Logger, testEnv *TestEnvironment, bcOutput *cre.WrappedBlockchainOutput, workflowName string, workflowConfig evm_config.Config) {
	forwarderAddress, _, err := crecontracts.FindAddressesForChain(testEnv.CreEnvironment.CldfEnvironment.ExistingAddresses, bcOutput.ChainSelector, keystonechangeset.KeystoneForwarder.String()) //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
	require.NoError(t, err, "failed to find forwarder address for chain %s", bcOutput.ChainSelector)

	forwarderContract, err := forwarder.NewKeystoneForwarder(forwarderAddress, bcOutput.SethClient.Client)
	require.NoError(t, err, "failed to instantiate forwarder contract")

	msgEmitterAddr := common.BytesToAddress(workflowConfig.ContractAddress)

	timeout := 5 * time.Minute
	tick := 3 * time.Second
	require.Eventually(t, func() bool {
		lggr.Info().Msgf("Waiting for workflow '%s' to finish", workflowName)
		ctx, cancel := context.WithTimeout(t.Context(), timeout)
		defer cancel()
		isSubmitted := isReportSubmittedByWorkflow(ctx, t, forwarderContract, msgEmitterAddr, workflowConfig)
		if !isSubmitted {
			lggr.Warn().Msgf("Forwarder has not received any reports from a workflow '%s' yet (delay is permissible due to latency in event propagation, waiting).", workflowName)
			return false
		}

		if isSubmitted {
			lggr.Info().Msgf("🎉 Workflow %s executed successfully on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
			return true
		}

		// if there are no more filtered reports, stop
		return !isReportSubmittedByWorkflow(ctx, t, forwarderContract, msgEmitterAddr, workflowConfig)
	}, timeout, tick, "workflow %s did not execute within the timeout %s", workflowName, timeout.String())
}

func configureEVMReadWorkflow(t *testing.T, lggr zerolog.Logger, chain *cre.WrappedBlockchainOutput) evm_config.Config {
	t.Helper()

	chainID := chain.BlockchainOutput.ChainID
	chainSethClient := chain.SethClient

	lggr.Info().Msgf("Deploying message emitter for chain %s", chainID)
	msgEmitterContractAddr, tx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(chainSethClient.NewTXOpts(), chainSethClient.Client)
	require.NoError(t, err, "failed to deploy message emitter contract")

	lggr.Info().Msgf("Deployed message emitter for chain '%s' at '%s'", chainID, msgEmitterContractAddr.String())
	_, err = chainSethClient.WaitMined(t.Context(), lggr, chainSethClient.Client, tx)
	require.NoError(t, err, "failed to get message emitter deployment tx")

	lggr.Printf("Emitting event to be picked up by workflow for chain '%s'", chainID)
	emittingTx, err := msgEmitter.EmitMessage(chainSethClient.NewTXOpts(), "Initial message to be read by workflow")
	require.NoError(t, err, "failed to emit message from contract '%s'", msgEmitterContractAddr.String())

	emittingReceipt, err := chainSethClient.WaitMined(t.Context(), lggr, chainSethClient.Client, emittingTx)
	require.NoError(t, err, "failed to get message emitter event tx")

	lggr.Info().Msgf("Updating nonces for chain %s", chainID)
	// force update nonces to ensure the transfer works
	require.NoError(t, chainSethClient.NonceManager.UpdateNonces(), "failed to update nonces for chain %s", chainID)

	// create and fund an address to be used by the workflow
	amountToFund := big.NewInt(0).SetUint64(10) // 10 wei
	numberOfAddressesToCreate := 1
	addresses, addrErr := createAndFundAddresses(t, lggr, numberOfAddressesToCreate, amountToFund, chainSethClient)
	require.NoError(t, addrErr, "failed to create and fund new addresses")
	require.Len(t, addresses, numberOfAddressesToCreate, "failed to create the correct number of addresses")

	marshalledTx, err := emittingTx.MarshalBinary()
	require.NoError(t, err)

	accountAddress := addresses[0].Bytes()
	return evm_config.Config{
		ContractAddress:  msgEmitterContractAddr.Bytes(),
		ChainSelector:    chain.ChainSelector,
		AccountAddress:   accountAddress,
		ExpectedBalance:  amountToFund,
		ExpectedReceipt:  emittingReceipt,
		TxHash:           emittingReceipt.TxHash.Bytes(),
		ExpectedBinaryTx: marshalledTx,
	}
}

// isReportSubmittedByWorkflow checks if a report has been submitted by the workflow by filtering the ReportProcessed events
func isReportSubmittedByWorkflow(ctx context.Context, t *testing.T, forwarderContract *forwarder.KeystoneForwarder, msgEmitterAddr common.Address, cfg evm_config.Config) bool {
	iter, err := forwarderContract.FilterReportProcessed(
		&bind.FilterOpts{
			Start:   cfg.ExpectedReceipt.BlockNumber.Uint64(),
			End:     nil,
			Context: ctx,
		},
		[]common.Address{msgEmitterAddr}, nil, nil)

	require.NoError(t, err, "failed to filter forwarder events")
	require.NoError(t, iter.Error(), "error during iteration of forwarder events")

	return iter.Next()
}

// regression
const (
	// find returned errors in the logs of the workflow
	balanceAtFunction                        = "BalanceAt"
	expectedBalanceAtError                   = "balanceAt errored"
	callContractInvalidAddressToReadFunction = "CallContract - invalid address to read"
	expectedCallContractInvalidAddressToRead = "balances=&[+0]" // expecting empty array of balances
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

	// BalanceAt
	// TODO: Move BalanceAt tests after fixing consensus crash because of invalid address
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

func EVMReadFailsTest(t *testing.T, testEnv *TestEnvironment, evmNegativeTest evmNegativeTest) {
	testLogger := framework.L
	const workflowFileLocation = "./evm/evmread-negative/main.go"
	enabledChains := getEVMEnabledChains(t, testEnv)

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

		listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)
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
		compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

		expectedError := evmNegativeTest.expectedError
		timeout := 2 * time.Minute
		err := assertBeholderMessage(listenerCtx, t, expectedError, testLogger, messageChan, kafkaErrChan, timeout)
		require.NoError(t, err, "EVM Read Fail test failed")
		testLogger.Info().Msg("EVM Read Fail test successfully completed")
	}
}

func getEVMEnabledChains(t *testing.T, testEnv *TestEnvironment) map[string]struct{} {
	t.Helper()

	enabledChains := map[string]struct{}{}
	for _, nodeSet := range testEnv.Config.NodeSets {
		require.NoError(t, nodeSet.ParseChainCapabilities())
		if nodeSet.ChainCapabilities == nil || nodeSet.ChainCapabilities[cre.EVMCapability] == nil {
			continue
		}

		for _, chainID := range nodeSet.ChainCapabilities[cre.EVMCapability].EnabledChains {
			strChainID := strconv.FormatUint(chainID, 10)
			enabledChains[strChainID] = struct{}{}
		}
	}
	require.NotEmpty(t, enabledChains, "No chains have EVM capability enabled in any node set")
	return enabledChains
}
