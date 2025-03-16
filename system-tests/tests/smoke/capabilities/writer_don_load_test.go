package capabilities

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	types2 "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/mock_capability"
	pb3 "github.com/smartcontractkit/chainlink/system-tests/lib/mock_capability/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
)

func TestKeystoneWithOCR3Workflow_TwoDons_WriteDON_Mock(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.MockCapability, keystonetypes.OCR3Capability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON, keystonetypes.WorkflowDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{keystonetypes.WriteEVMCapability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: 0,
			},
		}
	}

	setupOutput := setupTestEnvironment(t, testLogger, in, nil, nil, mustSetCapabilitiesFn, nil)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

			logDir := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name())

			removeErr := os.RemoveAll(logDir)
			if removeErr != nil {
				testLogger.Error().Err(removeErr).Msg("failed to remove log directory")
				return
			}

			_, saveErr := framework.SaveContainerLogs(logDir)
			if saveErr != nil {
				testLogger.Error().Err(saveErr).Msg("failed to save container logs")
				return
			}

			debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.DonsWithMetadata))
			for i, donWithMetadata := range setupOutput.donTopology.DonsWithMetadata {
				containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
				for _, output := range setupOutput.nodeOutput[i].Output.CLNodes {
					containerNames = append(containerNames, output.Node.ContainerName)
				}
				debugDons = append(debugDons, &keystonetypes.DebugDon{
					NodesMetadata:  donWithMetadata.NodesMetadata,
					Flags:          donWithMetadata.Flags,
					ContainerNames: containerNames,
				})
			}

			debugInput := keystonetypes.DebugInput{
				DebugDons:        debugDons,
				BlockchainOutput: setupOutput.blockchainOutput,
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	_, err = feeds_consumer.NewKeystoneFeedsConsumer(setupOutput.feedsConsumerAddress, setupOutput.sethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	//Mock capability start
	testLogger.Info().Msg("Connecting to mock capabilities")

	// 0. Connect to mock capabilities
	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	require.NoError(t, mocksClient.ConnectAll([]string{"127.0.0.1:13401", "127.0.0.1:13402", "127.0.0.1:13403", "127.0.0.1:13404"}, true, true))

	// 1. The write_geth-testnet capability is exposed from the write DON and the mock capability is hosted on the "workflow" DON,
	// we must wait until the workflow don registers the remote capability
	targetCapabilityID := "write_geth-testnet@1.0.0"
	targetCapabilityType := pb3.CapabilityType_Target
	ctx := tests.Context(t)
	testLogger.Info().Msg("Waiting for write_geth-testnet@1.0.0 to be exposed")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				found := false
				for _, c := range mocksClient.Nodes {
					list, err2 := c.API.List(ctx, &pb3.ListRequest{})
					require.NoError(t, err2)
					for _, l := range list.CapInfos {
						if l.ID == targetCapabilityID {
							found = true
						}
					}
					if !found {
						break
					}
				}

				if found {
					wg.Done()
					return
				}
			}
		}
	}()

	wg.Wait()

	testLogger.Info().Msg("RegisterToWorkflow on write_geth-testnet@1.0.0")
	// 2. Register workflow to write_geth-testnet@1.0.0 on all nodes of the worflow don
	// For the write target RegisterToWorkflow does not do anything!!
	//require.NoError(t, mocksClient.RegisterToWorkflow(ctx, pb3.RegisterToWorkflowRequest{
	//	ID:                   targetCapabilityID,
	//	CapabilityType:       targetCapabilityType,
	//	RegistrationMetadata: nil, //TODO
	//	Config:               nil, //TODO
	//}))

	// 3. Call Execute

	config, err := values.WrapMap(targets.Config{
		Address:  testutils.NewAddress().String(),
		GasLimit: nil, //TODO
	})
	require.NoError(t, err)
	configBytes, err := mock_capability.MapToBytes(config)
	require.NoError(t, err)

	inputs, err := values.WrapMap(targets.Inputs{
		SignedReport: types2.SignedReport{
			Report:     nil,
			Context:    nil,
			Signatures: nil,
			ID:         nil,
		},
	})
	require.NoError(t, err)
	inputsBytes, err := mock_capability.MapToBytes(inputs)
	require.NoError(t, err)
	for {
		select {
		case <-time.After(time.Second * 15):
			testLogger.Info().Msg("Execute on write_geth-testnet@1.0.0")
			require.NoError(t, mocksClient.Execute(ctx, pb3.ExecutableRequest{
				ID:             targetCapabilityID,
				CapabilityType: targetCapabilityType,
				RequestMetadata: &pb3.Metadata{
					WorkflowID:               "some-workflow",
					WorkflowOwner:            "",
					WorkflowExecutionID:      "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed",
					WorkflowName:             "nonexistent workflow",
					WorkflowDonID:            0,
					WorkflowDonConfigVersion: 0,
					ReferenceID:              "",
					DecodedWorkflowName:      "",
				}, //TODO
				Config: configBytes,
				Inputs: inputsBytes,
			}))
		}
	}
	time.Sleep(time.Minute * 10)
}
