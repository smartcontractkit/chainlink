package cre

import (
	"context"
	"fmt"
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

	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/config"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func executeEVMReadTest(t *testing.T, testEnv *TestEnvironment) {
	lggr := framework.L
	const workflowFileLocation = "./evmread/main.go"
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
		forwarderAddress, _, err := crecontracts.FindAddressesForChain(testEnv.FullCldEnvOutput.Environment.ExistingAddresses, bcOutput.ChainSelector, keystonechangeset.KeystoneForwarder.String()) //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
		require.NoError(t, err)

		// validate workflow execution
		go func(bcOutput *cre.WrappedBlockchainOutput) {
			defer workflowsWg.Done()
			err := validateWorkflowExecution(t, lggr, bcOutput, workflowName, forwarderAddress, workflowConfig)
			if err != nil {
				lggr.Error().Msgf("Workflow %s execution failed on chain %s: %v", workflowName, bcOutput.BlockchainOutput.ChainID, err)
				return
			}

			lggr.Info().Msgf("Workflow %s executed successfully on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
			successfulWorkflowRuns.Add(1)
		}(bcOutput)
	}

	_, messageChan, kafkaErrChan := startBeholder(t, lggr, testEnv)
	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		workflowsWg.Wait()
		cancel()
	}()
	logBeholderMessages(ctx, t, lggr, testEnv, messageChan, kafkaErrChan)
	require.Equal(t, len(enabledChains), int(successfulWorkflowRuns.Load()), "Not all workflows executed successfully")
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

func validateWorkflowExecution(t *testing.T, lggr zerolog.Logger, bcOutput *cre.WrappedBlockchainOutput, workflowName string, forwarderAddr common.Address, cfg config.Config) error {
	forwarderContract, err := forwarder.NewKeystoneForwarder(forwarderAddr, bcOutput.SethClient.Client)
	if err != nil {
		return fmt.Errorf("failed to create forwarder contract instance: %w", err)
	}
	msgEmitterAddr := common.BytesToAddress(cfg.ContractAddress)
	isWorkflowFinished := func(ctx context.Context) (bool, error) {
		iter, err := forwarderContract.FilterReportProcessed(&bind.FilterOpts{
			Start:   cfg.ExpectedReceipt.BlockNumber.Uint64(),
			End:     nil,
			Context: ctx,
		}, []common.Address{msgEmitterAddr}, nil, nil)
		if err != nil {
			return false, fmt.Errorf("failed to filter forwarder: %w", err)
		}

		if iter.Error() != nil {
			return false, fmt.Errorf("error while filtering forwarder: %w", iter.Error())
		}

		return iter.Next(), nil
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Minute)
	defer cancel()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			lggr.Info().Msgf("Checking if workflow %s executed on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
			ok, err := isWorkflowFinished(ctx)
			if err != nil {
				lggr.Error().Msgf("Error checking workflow execution: %v", err)
				continue
			}

			if ok {
				lggr.Info().Msgf("Workflow %s executed successfully on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("workflow %s did not execute on chain %s within the timeout", workflowName, bcOutput.BlockchainOutput.ChainID)
		}
	}
}

func configureEVMReadWorkflow(t *testing.T, lggr zerolog.Logger, chain *cre.WrappedBlockchainOutput) config.Config {
	t.Helper()
	chainID := chain.BlockchainOutput.ChainID
	chainSethClient := chain.SethClient

	lggr.Info().Msgf("Deploying message emitter for chain %s", chainID)
	msgEmitterContractAddr, tx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(chain.SethClient.NewTXOpts(), chain.SethClient.Client)
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

	amountToFund := big.NewInt(0).SetUint64(10) // 10 wei
	numberOfAddressesToCreate := 1
	addresses, addrErr := createAndFundAddresses(t, lggr, numberOfAddressesToCreate, amountToFund, chainSethClient)
	require.NoError(t, addrErr, "failed to create and fund new addresses")
	require.Len(t, addresses, numberOfAddressesToCreate, "failed to create the correct number of addresses")

	marshalledTx, err := emittingTx.MarshalBinary()
	require.NoError(t, err)

	accountAddress := addresses[0].Bytes()
	return config.Config{
		ContractAddress:  msgEmitterContractAddr.Bytes(),
		ChainSelector:    chain.ChainSelector,
		AccountAddress:   accountAddress,
		ExpectedBalance:  amountToFund,
		ExpectedReceipt:  emittingReceipt,
		TxHash:           emittingReceipt.TxHash.Bytes(),
		ExpectedBinaryTx: marshalledTx,
	}
}
