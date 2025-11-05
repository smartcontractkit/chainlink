package cre

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/contracts"
	evm_logTrigger_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/logtrigger/config"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/contracts"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// smoke
func ExecuteEVMReadTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	lggr := framework.L
	const workflowFileLocation = "./evm/evmread/main.go"
	enabledChains := t_helpers.GetEVMEnabledChains(t, testEnv)

	var workflowsWg sync.WaitGroup
	var successfulWorkflowRuns atomic.Int32
	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		chainID := bcOutput.CtfOutput().ChainID
		if _, ok := enabledChains[chainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM Read workflow test", chainID)
			continue
		}

		lggr.Info().Msg("Creating EVM Read workflow configuration...")
		require.IsType(t, &evm.Blockchain{}, bcOutput, "expected EVM blockchain type")
		evmChain := bcOutput.(*evm.Blockchain)
		workflowConfig := configureEVMReadWorkflow(t, lggr, evmChain)
		workflowName := fmt.Sprintf("evm-read-workflow-%s-%04d", chainID, rand.Intn(10000))
		t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, workflowName, &workflowConfig, workflowFileLocation)

		workflowsWg.Add(1)
		go func(evmChain *evm.Blockchain) {
			defer workflowsWg.Done()
			validateWorkflowExecution(t, lggr, testEnv, evmChain, workflowName, common.BytesToAddress(workflowConfig.ContractAddress), workflowConfig.ExpectedReceipt.BlockNumber.Uint64()) //nolint:testifylint // TODO: consider refactoring
			successfulWorkflowRuns.Add(1)
		}(evmChain)
	}

	// wait for all workflows to complete
	workflowsWg.Wait()
	require.Equal(t, len(enabledChains), int(successfulWorkflowRuns.Load()), "Not all workflows executed successfully")
}

func configureEVMReadWorkflow(t *testing.T, lggr zerolog.Logger, chain *evm.Blockchain) evm_config.Config {
	t.Helper()

	chainID := chain.CtfOutput().ChainID
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
	addresses, addrErr := t_helpers.CreateAndFundAddresses(t, lggr, numberOfAddressesToCreate, amountToFund, chain, nil)
	require.NoError(t, addrErr, "failed to create and fund new addresses")
	require.Len(t, addresses, numberOfAddressesToCreate, "failed to create the correct number of addresses")

	marshalledTx, err := emittingTx.MarshalBinary()
	require.NoError(t, err)

	accountAddress := addresses[0].Bytes()
	return evm_config.Config{
		ContractAddress:  msgEmitterContractAddr.Bytes(),
		ChainSelector:    chain.ChainSelector(),
		AccountAddress:   accountAddress,
		ExpectedBalance:  amountToFund,
		ExpectedReceipt:  emittingReceipt,
		TxHash:           emittingReceipt.TxHash.Bytes(),
		ExpectedBinaryTx: marshalledTx,
	}
}

func validateWorkflowExecution(t *testing.T, lggr zerolog.Logger, testEnv *ttypes.TestEnvironment, blockchain *evm.Blockchain, workflowName string, msgEmitterAddr common.Address, startBlock uint64) {
	forwarderAddress, _, err := crecontracts.FindAddressesForChain(testEnv.CreEnvironment.CldfEnvironment.ExistingAddresses, blockchain.ChainSelector(), keystonechangeset.KeystoneForwarder.String()) //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
	require.NoError(t, err, "failed to find forwarder address for chain %s", blockchain.ChainSelector)

	forwarderContract, err := forwarder.NewKeystoneForwarder(forwarderAddress, blockchain.SethClient.Client)
	require.NoError(t, err, "failed to instantiate forwarder contract")

	timeout := 5 * time.Minute
	tick := 3 * time.Second
	require.Eventually(t, func() bool {
		lggr.Info().Msgf("Waiting for workflow '%s' to finish", workflowName)
		ctx, cancel := context.WithTimeout(t.Context(), timeout)
		defer cancel()
		isSubmitted := isReportSubmittedByWorkflow(ctx, t, forwarderContract, msgEmitterAddr, startBlock)
		if !isSubmitted {
			lggr.Warn().Msgf("Forwarder has not received any reports from a workflow '%s' yet (delay is permissible due to latency in event propagation, waiting).", workflowName)
			return false
		}

		if isSubmitted {
			lggr.Info().Msgf("🎉 Workflow %s executed successfully on chain %s", workflowName, blockchain.CtfOutput().ChainID)
			return true
		}

		// if there are no more filtered reports, stop
		return !isReportSubmittedByWorkflow(ctx, t, forwarderContract, msgEmitterAddr, startBlock)
	}, timeout, tick, "workflow %s did not execute within the timeout %s", workflowName, timeout.String())
}

// isReportSubmittedByWorkflow checks if a report has been submitted by the workflow by filtering the ReportProcessed events
func isReportSubmittedByWorkflow(ctx context.Context, t *testing.T, forwarderContract *forwarder.KeystoneForwarder, msgEmitterAddr common.Address, startBlock uint64) bool {
	iter, err := forwarderContract.FilterReportProcessed(
		&bind.FilterOpts{
			Start:   startBlock,
			End:     nil,
			Context: ctx,
		},
		[]common.Address{msgEmitterAddr}, nil, nil)

	require.NoError(t, err, "failed to filter forwarder events")
	require.NoError(t, iter.Error(), "error during iteration of forwarder events")

	return iter.Next()
}

func keysFromMap(m map[string]blockchains.Blockchain) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func emitEvent(t *testing.T, lggr zerolog.Logger, chainID string, bcOutput blockchains.Blockchain, msgEmitter *evmreadcontracts.MessageEmitter, expectedUserLog string, workflowConfig evm_logTrigger_config.Config) uint64 {
	lggr.Info().Msgf("Emitting event to be picked up by workflow for chain '%s'", chainID)
	sethClient := bcOutput.(*evm.Blockchain).SethClient
	emittingTx, err := msgEmitter.EmitMessage(sethClient.NewTXOpts(), expectedUserLog)
	require.NoError(t, err, "failed to emit message from contract '%s'", workflowConfig.Addresses[0])

	emittingReceipt, err := sethClient.WaitMined(t.Context(), lggr, sethClient.Client, emittingTx)
	require.NoError(t, err, "failed to get message emitter event tx")
	lggr.Info().Msgf("Transaction for chain '%s' mined at '%d' with emitted message %q", chainID, emittingReceipt.BlockNumber.Uint64(), expectedUserLog)
	return emittingReceipt.BlockNumber.Uint64()
}

func configureEVMLogTriggerWorkflow(t *testing.T, lggr zerolog.Logger, chain blockchains.Blockchain) (evm_logTrigger_config.Config, *evmreadcontracts.MessageEmitter) {
	t.Helper()

	evmChain := chain.(*evm.Blockchain)
	chainID := evmChain.CtfOutput().ChainID
	chainSethClient := evmChain.SethClient

	lggr.Info().Msgf("Deploying message emitter for chain %s", chainID)
	msgEmitterContractAddr, tx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(chainSethClient.NewTXOpts(), chainSethClient.Client)
	require.NoError(t, err, "failed to deploy message emitter contract")

	lggr.Info().Msgf("Deployed message emitter for chain '%s' at '%s'", chainID, msgEmitterContractAddr.String())
	_, err = chainSethClient.WaitMined(t.Context(), lggr, chainSethClient.Client, tx)
	require.NoError(t, err, "failed to get message emitter deployment tx")

	abiDef, err := contracts.MessageEmitterMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}

	eventName := "MessageEmitted"
	topicFromABI := abiDef.Events[eventName].ID
	eventSigMessageEmitted := topicFromABI.Hex()
	lggr.Info().Msgf("Topic0 (ABI): %s", eventSigMessageEmitted)

	return evm_logTrigger_config.Config{
		ChainSelector: evmChain.ChainSelector(),
		Addresses:     []string{msgEmitterContractAddr.Hex()},
		Topics: []struct {
			Values []string `yaml:"values"`
		}{
			{Values: []string{eventSigMessageEmitted}},
		},
		Event: eventName,
		Abi:   evmreadcontracts.MessageEmitterMetaData.ABI,
	}, msgEmitter
}

func ExecuteEVMLogTriggerTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	const workflowFileLocation = "./evm/logtrigger/main.go"
	lggr := framework.L

	enabledChains := t_helpers.GetEVMEnabledChains(t, testEnv)
	chainsToTest := make(map[string]blockchains.Blockchain)
	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		chainID := bcOutput.CtfOutput().ChainID
		if _, ok := enabledChains[chainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM LogTrigger workflow test", chainID)
			continue
		}
		chainsToTest[chainID] = bcOutput
	}

	successfulLogTriggerChains := make([]string, 0, len(chainsToTest))
	for chainID, bcOutput := range chainsToTest {
		lggr.Info().Msgf("Creating EVM LogTrigger workflow configuration for chain %s", chainID)
		workflowConfig, msgEmitter := configureEVMLogTriggerWorkflow(t, lggr, bcOutput)

		workflowName := fmt.Sprintf("evm-logTrigger-workflow-%s-%04d", chainID, rand.Intn(10000))
		lggr.Info().Msgf("About to deploy Workflow %s on chain %s", workflowName, chainID)
		t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, workflowName, &workflowConfig, workflowFileLocation)

		time.Sleep(30 * time.Second) // wait for trigger to be registered

		message := "Data for log trigger"
		startBlock := emitEvent(t, lggr, chainID, bcOutput, msgEmitter, message, workflowConfig)
		evmChain := bcOutput.(*evm.Blockchain)
		validateWorkflowExecution(t, lggr, testEnv, evmChain, workflowName, common.HexToAddress(workflowConfig.Addresses[0]), startBlock)

		lggr.Info().Msgf("🎉 LogTrigger Workflow %s executed successfully on chain %s", workflowName, chainID)
		successfulLogTriggerChains = append(successfulLogTriggerChains, chainID)
	}

	require.Lenf(t, successfulLogTriggerChains, len(chainsToTest),
		"Not all workflows executed successfully. Successful chains: %v, All chains to test: %v",
		successfulLogTriggerChains, keysFromMap(chainsToTest))

	lggr.Info().Msgf("✅ LogTrigger test ran for chains: %v", successfulLogTriggerChains)
}
