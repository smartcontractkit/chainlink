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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/config"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func executeEVMReadTest(t *testing.T, testEnv *TestEnvironment) {
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
	const workflowFileLocation = "./evmread/main.go"
	lggr := framework.L
	var workflowsWg sync.WaitGroup
	var successfulWorkflowRuns atomic.Int32
	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		if _, ok := enabledChains[bcOutput.BlockchainOutput.ChainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM read workflow test", bcOutput.BlockchainOutput.ChainID)
			continue
		}
		workflowName := "evm-read-workflow-" + bcOutput.BlockchainOutput.ChainID

		workflowConfig := configureEVMReadWorkflow(t, lggr, bcOutput)

		lggr.Info().Msg("Proceeding to register workflow...")
		compileAndDeployWorkflow(t, testEnv, lggr, fmt.Sprintf("evmreadtest-%d", bcOutput.ChainID), &workflowConfig, workflowFileLocation)

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

	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		workflowsWg.Wait()
		cancel()
	}()
	logWorkflowLogs(ctx, t, testEnv)
	require.Equal(t, len(enabledChains), int(successfulWorkflowRuns.Load()), "Not all workflows executed successfully")
}

func logWorkflowLogs(ctx context.Context, t *testing.T, testEnv *TestEnvironment) {
	beholder, err := NewBeholder(framework.L, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
	require.NoError(t, err, "failed to create beholder instance")

	// We are interested in UserLogs (successful execution)
	// or BaseMessage with specific error message (engine initialization failure)
	beholderMessageTypes := map[string]func() proto.Message{
		"workflows.v1.UserLogs": func() proto.Message {
			return &workflowevents.UserLogs{}
		},
		"BaseMessage": func() proto.Message {
			return &commonevents.BaseMessage{}
		},
	}

	lggr := framework.L
	beholderMsgChan, beholderErrChan := beholder.SubscribeToBeholderMessages(ctx, beholderMessageTypes)
	// Check the beholder logs for the expected messages
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-beholderErrChan:
			require.FailNowf(t, "Kafka error received from Kafka %s", err.Error())
		case msg := <-beholderMsgChan:
			switch typedMsg := msg.(type) {
			case *commonevents.BaseMessage:
				lggr.Info().Msgf("Received BaseMessage from Beholder: %s", typedMsg.Msg)
			case *workflowevents.UserLogs:
				for _, logLine := range typedMsg.LogLines {
					lggr.Info().Msgf("Received workflow msg: %s", logLine.Message)
				}
			default:
				lggr.Info().Msgf("Received unknown message of type '%T'", msg)
			}
		}
	}
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
	lggr.Info().Msgf("Deploying message emitter for chain %s", chain.BlockchainOutput.ChainID)
	msgEmitterContractAddr, tx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(chain.SethClient.NewTXOpts(), chain.SethClient.Client)
	require.NoError(t, err)
	lggr.Info().Msgf("Deployed message emitter for chain %s at %s", chain.BlockchainOutput.ChainID, msgEmitterContractAddr.String())
	_, err = chain.SethClient.WaitMined(t.Context(), lggr, chain.SethClient.Client, tx)
	require.NoError(t, err)
	lggr.Printf("Emitting event to be picked up by workflow for chain %s", chain.BlockchainOutput.ChainID)
	emittingTx, err := msgEmitter.EmitMessage(chain.SethClient.NewTXOpts(), "Initial message to be read by workflow")
	require.NoError(t, err)
	emittingReceipt, err := chain.SethClient.WaitMined(t.Context(), lggr, chain.SethClient.Client, emittingTx)
	require.NoError(t, err)
	lggr.Info().Msgf("Updating nonces for chain %s", chain.BlockchainOutput.ChainID)
	// force update nonces to ensure the transfer works
	require.NoError(t, chain.SethClient.NonceManager.UpdateNonces())
	const expectedBalance = 10
	pk, err := crypto.GenerateKey()
	require.NoError(t, err)
	accountAddr := crypto.PubkeyToAddress(pk.PublicKey)
	lggr.Info().Msgf("Funding account %s for BalanceAt read test for chain %s", accountAddr.Hex(), chain.BlockchainOutput.ChainID)
	err = chain.SethClient.TransferETHFromKey(t.Context(), 0, accountAddr.Hex(), big.NewInt(expectedBalance), nil)
	require.NoError(t, err, "failed to transfer ETH to contract %s", msgEmitterContractAddr.String())
	marshalledTx, err := emittingTx.MarshalBinary()
	require.NoError(t, err)
	return config.Config{
		ContractAddress:  msgEmitterContractAddr.Bytes(),
		ChainSelector:    chain.ChainSelector,
		AccountAddress:   accountAddr.Bytes(),
		ExpectedBalance:  big.NewInt(expectedBalance),
		ExpectedReceipt:  emittingReceipt,
		TxHash:           emittingReceipt.TxHash.Bytes(),
		ExpectedBinaryTx: marshalledTx,
	}
}
