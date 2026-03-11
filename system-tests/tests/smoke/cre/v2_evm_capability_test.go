package cre

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

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
	testCases := make([]evm_config.TestCase, 0, evm_config.TestCaseLen)
	for tc := range evm_config.TestCaseLen {
		testCases = append(testCases, tc)
	}

	ExecuteEVMReadTestForCases(t, testEnv, testCases)
}

func ExecuteEVMReadTestForCases(t *testing.T, testEnv *ttypes.TestEnvironment, testCases []evm_config.TestCase) {
	require.NoError(t, evm_config.ValidateReadBucketRegistry(), "invalid EVM read bucket registry; assign each testcase exactly once")
	require.NotEmpty(t, testCases, "no EVM read testcases selected")

	seen := make(map[evm_config.TestCase]struct{}, len(testCases))
	for _, tc := range testCases {
		require.GreaterOrEqualf(t, tc, evm_config.TestCase(0), "invalid testcase %d", tc)
		require.Lessf(t, tc, evm_config.TestCaseLen, "invalid testcase %d", tc)
		if _, alreadySeen := seen[tc]; alreadySeen {
			require.Failf(t, "duplicate testcase", "testcase %q selected more than once", tc.String())
		}

		seen[tc] = struct{}{}
	}

	lggr := framework.L
	const workflowFileLocation = "./evm/evmread/main.go"
	enabledChains := t_helpers.GetEVMEnabledChains(t, testEnv)

	userLogsCh := makeSinkCh[*workflowevents.UserLogs]()
	baseMessageCh := makeSinkCh[*commonevents.BaseMessage]()

	// `./logs` folder inside `smoke/cre` is uploaded as artifact in GH
	server := t_helpers.StartChipTestSink(t, t_helpers.GetLoggingPublishFn(lggr, userLogsCh, baseMessageCh, "./logs/evm_read_workflow_test.log"))
	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		chainID := bcOutput.CtfOutput().ChainID
		if _, ok := enabledChains[chainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM Read workflow test", chainID)
			continue
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Read %s on chain %s", tc.String(), chainID), func(t *testing.T) {
				workflowName := fmt.Sprintf("evm-read-workflow-%s-%04d", chainID, rand.Intn(10000))
				lggr.Info().
					Str("workflow_name", workflowName).
					Str("chain_id", chainID).
					Str("test_case", tc.String()).
					Msg("Creating EVM Read workflow configuration...")
				require.IsType(t, &evm.Blockchain{}, bcOutput, "expected EVM blockchain type")
				evmChain := bcOutput.(*evm.Blockchain)
				workflowConfig := configureEVMReadWorkflow(t, lggr, evmChain, tc, workflowName)
				t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, workflowName, &workflowConfig, workflowFileLocation)

				validateWorkflowExecution(t, lggr, testEnv, evmChain, workflowName, common.BytesToAddress(workflowConfig.ContractAddress), workflowConfig.ExpectedReceipt.BlockNumber.Uint64())
			})
		}
	}
}

func makeSinkCh[T any]() chan T {
	c := make(chan T, 1)
	go func() {
		//nolint:revive //drain the channel to prevent blocking. Content is processed elsewhere.
		for range c {
		}
	}()

	return c
}

func configureEVMReadWorkflow(t *testing.T, lggr zerolog.Logger, chain *evm.Blockchain, testCase evm_config.TestCase, workflowName string) evm_config.Config {
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
		TestCase:         testCase,
		WorkflowName:     workflowName,
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
	forwarderAddress := crecontracts.MustGetAddressFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore, blockchain.ChainSelector(), keystonechangeset.KeystoneForwarder.String(), testEnv.CreEnvironment.ContractVersions[keystonechangeset.KeystoneForwarder.String()], "")
	forwarderContract, err := forwarder.NewKeystoneForwarder(common.HexToAddress(forwarderAddress), blockchain.SethClient.Client)
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
	}, timeout, tick, "workflow %s did not execute within the timeout %s. Check logs of parent test.", workflowName, timeout.String())
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
	if err != nil {
		lggr.Info().Msgf("Failed to emit transaction for chain '%s': %v", chainID, err)
		return 0
	}

	emittingReceipt, err := sethClient.WaitMined(t.Context(), lggr, sethClient.Client, emittingTx)
	if err != nil {
		lggr.Info().Msgf("Failed to emit receipt for chain '%s': %v", chainID, err)
		return 0
	}
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

func connectCapabilitiesDB(t *testing.T) *sql.DB {
	// TODO: Only if the TOPOLOGY is a remote trigger
	// TODO: For local trigger, we should watch the db in wf node
	t.Helper()
	db, err := sql.Open(
		"postgres",
		"host=localhost port=13100 user=chainlink password=thispasswordislongenough dbname=db_0 sslmode=disable")
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	t.Log("connected to Capabilities db")
	return db
}

// trackedEvent holds the identity of a specific event row we are tracking through
// its insert → ACK → delete lifecycle.
type trackedEvent struct {
	TriggerID string
	EventID   string
}

// watchTriggerLifecycle polls the trigger_pending_events table until it spots a
// newly inserted row, records that row's (trigger_id, event_id), and then waits
// for that specific row to be deleted (i.e. ACKed by BaseTrigger).
func watchTriggerLifecycle(
	ctx context.Context,
	db *sql.DB,
	pollInterval time.Duration,
) (<-chan trackedEvent, <-chan trackedEvent, <-chan error) {
	insertedCh := make(chan trackedEvent, 1)
	deletedCh := make(chan trackedEvent, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(insertedCh)
		defer close(deletedCh)
		defer close(errCh)

		var tracked *trackedEvent

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if tracked == nil {
					var triggerID, eventID string
					err := db.QueryRowContext(ctx,
						`SELECT trigger_id, event_id FROM cre.trigger_pending_events ORDER BY first_at ASC LIMIT 1`,
					).Scan(&triggerID, &eventID)
					if err == sql.ErrNoRows {
						continue
					}
					if err != nil {
						errCh <- fmt.Errorf("polling for insert: %w", err)
						return
					}

					tracked = &trackedEvent{TriggerID: triggerID, EventID: eventID}
					fmt.Printf("SEEN INSERTED trigger_id=%s event_id=%s\n", triggerID, eventID)
					insertedCh <- *tracked
					continue
				}

				var exists bool
				err := db.QueryRowContext(ctx,
					`SELECT EXISTS(SELECT 1 FROM cre.trigger_pending_events WHERE trigger_id = $1 AND event_id = $2)`,
					tracked.TriggerID, tracked.EventID,
				).Scan(&exists)
				if err != nil {
					errCh <- fmt.Errorf("polling for delete: %w", err)
					return
				}

				if !exists {
					fmt.Printf("SEEN DELETED trigger_id=%s event_id=%s\n", tracked.TriggerID, tracked.EventID)
					deletedCh <- *tracked
					return
				}
			}
		}
	}()

	return insertedCh, deletedCh, errCh
}

func ExecuteEVMLogTriggerTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	const workflowFileLocation = "./evm/logtrigger/main.go"
	lggr := framework.L

	enabledChains := t_helpers.GetEVMEnabledChains(t, testEnv)
	chainsToTest := make(map[string]blockchains.Blockchain)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(lggr, userLogsCh, baseMessageCh))

	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		chainID := bcOutput.CtfOutput().ChainID
		if _, ok := enabledChains[chainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM LogTrigger workflow test", chainID)
			continue
		}
		chainsToTest[chainID] = bcOutput
	}

	capDB := connectCapabilitiesDB(t)
	_, err := capDB.ExecContext(t.Context(), `DELETE FROM cre.trigger_pending_events`)
	require.NoError(t, err, "failed to clean trigger_pending_events before test")

	insertedCh, deletedCh, errCh := watchTriggerLifecycle(t.Context(), capDB, 50*time.Millisecond)

	successfulLogTriggerChains := make([]string, 0, len(chainsToTest))
	for chainID, bcOutput := range chainsToTest {
		lggr.Info().Msgf("Creating EVM LogTrigger workflow configuration for chain %s", chainID)
		workflowConfig, msgEmitter := configureEVMLogTriggerWorkflow(t, lggr, bcOutput)

		workflowName := fmt.Sprintf("evm-logTrigger-workflow-%s-%04d", chainID, rand.Intn(10000))
		lggr.Info().Msgf("About to deploy Workflow %s on chain %s", workflowName, chainID)
		t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, workflowName, &workflowConfig, workflowFileLocation)

		message := "Data for log trigger"
		// start background event emission every 10s while WatchWorkflowLogs is running, so that the workflow has events to pick up eventually
		var emittedEventCount int64
		ticker := time.NewTicker(10 * time.Second)

		// create a context that will be cancelled as soon as we either find the log we are looking for or timeout
		emitCtx, emitCancelFn := context.WithCancel(t.Context())
		go func() {
			defer func() {
				emitCancelFn()
				ticker.Stop()
			}()
			for {
				select {
				case <-emitCtx.Done():
					return
				case <-ticker.C:
					lggr.Info().Msgf("About to emit event #%d for chain %s", emittedEventCount, chainID)
					blockNumber := emitEvent(t, lggr, chainID, bcOutput, msgEmitter, message, workflowConfig)
					lggr.Info().Msgf("Event emitted for chain %s at blockNumber %d", chainID, blockNumber)
					emittedEventCount++
				}
			}
		}()
		expectedUserLog := "OnTrigger decoded message: message:" + message

		t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedUserLog, 4*time.Minute)
		emitCancelFn()

		// Verify base trigger persists trigger event
		require.Eventually(t, func() bool {
			select {
			case err := <-errCh:
				require.NoError(t, err)
				return false
			case ev, ok := <-insertedCh:
				if !ok {
					return false
				}
				t.Logf("event inserted into trigger_pending_events (trigger_id=%s, event_id=%s)", ev.TriggerID, ev.EventID)
				return true
			default:
				return false
			}
		}, 30*time.Second, time.Second, "event was never inserted")

		// Verify ACK occurs on base trigger via deletion of the tracked event
		require.Eventually(t, func() bool {
			select {
			case err := <-errCh:
				require.NoError(t, err)
				return false
			case ev, ok := <-deletedCh:
				if !ok {
					return false
				}
				t.Logf("event deleted from trigger_pending_events (trigger_id=%s, event_id=%s)", ev.TriggerID, ev.EventID)
				return true
			default:
				return false
			}
		}, 2*time.Minute, time.Second, "tracked event was never deleted (ACK did not occur)")

		lggr.Info().Msgf("🎉 LogTrigger Workflow %s executed successfully on chain %s", workflowName, chainID)
		successfulLogTriggerChains = append(successfulLogTriggerChains, chainID)
	}

	require.Lenf(t, successfulLogTriggerChains, len(chainsToTest),
		"Not all workflows executed successfully. Successful chains: %v, All chains to test: %v",
		successfulLogTriggerChains, keysFromMap(chainsToTest))

	lggr.Info().Msgf("✅ LogTrigger test ran for chains: %v", successfulLogTriggerChains)
}
