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

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/contracts"
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
	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		chainID := bcOutput.BlockchainOutput.ChainID
		if _, ok := enabledChains[chainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM Read workflow test", chainID)
			continue
		}

		lggr.Info().Msg("Creating EVM Read workflow configuration...")
		workflowConfig := configureEVMReadWorkflow(t, lggr, bcOutput)
		workflowName := fmt.Sprintf("evm-read-workflow-%s-%04d", chainID, rand.Intn(10000))
		t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, workflowName, &workflowConfig, workflowFileLocation)

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

func validateWorkflowExecution(t *testing.T, lggr zerolog.Logger, testEnv *ttypes.TestEnvironment, bcOutput *cre.WrappedBlockchainOutput, workflowName string, workflowConfig evm_config.Config) {
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
	addresses, addrErr := t_helpers.CreateAndFundAddresses(t, lggr, numberOfAddressesToCreate, amountToFund, chainSethClient, chain, nil)
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
