package capabilities

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	types2 "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
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

	// Save forwarder address during setup
	if err := saveForwarderAddressToCache(setupOutput.forwarderAddress); err != nil {
		testLogger.Error().Err(err).Msg("failed to cache forwarder address")
	}

	//Mock capability start
	testLogger.Info().Msg("Connecting to mock capabilities")
	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.InfraType == "docker" {
		mockClientsAddress = []string{"127.0.0.1:13401", "127.0.0.1:13402", "127.0.0.1:13403", "127.0.0.1:13404"}
	} else {
		for i, _ := range setupOutput.nodeOutput[1].CLNodes {
			if i == 0 {
				continue
			}
			mockClientsAddress = append(mockClientsAddress, fmt.Sprintf("%s-%s-%d-mock.main.stage.cldev.sh:443", in.Infra.CRIB.Namespace, setupOutput.nodeOutput[1].NodeSetName, i-1))
		}
	}

	useInsecure := false
	if in.Infra.InfraType == "docker" {
		useInsecure = true
	}

	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, useInsecure, true)) //Capability don ports

	// 1. The write_geth-testnet capability is exposed from the write DON and the mock capability is hosted on the "workflow" DON,
	// we must wait until the workflow don registers the remote capability
	targetCapabilityID := "write_geth-testnet@1.0.0"
	ctx := tests.Context(t)
	testLogger.Info().Msg("Waiting for write_geth-testnet@1.0.0 to be exposed")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				testLogger.Info().Msg("Checking for write_geth-testnet@1.0.0 to be exposed")
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

	labels := map[string]string{
		"go_test_name": "Writer DON Load Test",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 10*time.Minute),
			),
			Gun:                   NewWriteGun(mocksClient, setupOutput.forwarderAddress, "write_geth-testnet@1.0.0"),
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err)
}

func TestReconnectWrite(t *testing.T) {
	testLogger := framework.L

	// Read forwarder address in generateTargetRequest
	forwarderAddress, err := getForwarderAddressFromCache()
	require.NoError(t, err)

	mocksClient, err := mock_capability.NewMockCapabilityControllerFromCache(testLogger, false)
	require.NoError(t, err)

	labels := map[string]string{
		"go_test_name": "Writer DON Load Test",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 10*time.Minute),
			),
			Gun:                   NewWriteGun(mocksClient, forwarderAddress, "write_geth-testnet@1.0.0"),
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err)

}

var _ wasp.Gun = (*WriteGun)(nil)

type WriteGun struct {
	mockClientController *mock_capability.MockCapabilityController
	targetID             string
	forwarderAddress     common.Address
}

func NewWriteGun(mockClientController *mock_capability.MockCapabilityController, forwarderAddress common.Address, targetID string) *WriteGun {
	return &WriteGun{
		mockClientController: mockClientController,
		targetID:             targetID,
		forwarderAddress:     forwarderAddress,
	}
}
func (w WriteGun) Call(l *wasp.Generator) *wasp.Response {

	metadata, config, input, err := generateTargetRequest(w.forwarderAddress)
	if err != nil {
		panic(err)
	}

	err = w.mockClientController.Execute(context.TODO(), pb3.ExecutableRequest{
		ID:              w.targetID,
		CapabilityType:  4,
		RequestMetadata: metadata,
		Config:          config,
		Inputs:          input,
	})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	return &wasp.Response{}
}

func generateTargetRequest(forwarderAddress common.Address) (*pb3.Metadata, []byte, []byte, error) {
	reportID := [2]byte{0x00, 0x01}
	var workflowName [10]byte
	copy(workflowName[:], []byte("name"))
	workflowOwnerString := "219BFD3D78fbb740c614432975CBE829E26C490e"
	workflowOwner := common.HexToAddress(workflowOwnerString)
	reportMetadata := targets.ReportV1Metadata{
		Version:             1,
		WorkflowExecutionID: [32]byte{},
		Timestamp:           0,
		DonID:               0,
		DonConfigVersion:    0,
		WorkflowCID:         [32]byte{},
		WorkflowName:        workflowName,
		WorkflowOwner:       workflowOwner,
		ReportID:            reportID,
	}
	reportMetadataBytes, err := reportMetadata.Encode()
	if err != nil {
		return nil, nil, nil, err
	}

	inputs, err := values.WrapMap(targets.Inputs{
		SignedReport: types2.SignedReport{
			Report:     reportMetadataBytes,
			Context:    []byte{4, 5},
			Signatures: [][]byte{},
			ID:         reportID[:],
		},
	})

	configGasLimit, err2 := values.NewMap(map[string]any{
		"Address":  forwarderAddress.String(),
		"GasLimit": 500000,
	})
	if err2 != nil {
		return nil, nil, nil, err
	}
	configBytes, err2 := mock_capability.MapToBytes(configGasLimit)
	if err2 != nil {
		return nil, nil, nil, err
	}

	inputsBytes, err := mock_capability.MapToBytes(inputs)
	if err != nil {
		return nil, nil, nil, err
	}
	RequestMetadata := pb3.Metadata{
		WorkflowID:               "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0",
		WorkflowOwner:            workflowOwner.String(),
		WorkflowExecutionID:      randomID(),
		WorkflowName:             "nonexistent workflow",
		WorkflowDonID:            0,
		WorkflowDonConfigVersion: 1,
		ReferenceID:              "mock", // What is this?
		DecodedWorkflowName:      "abcdefgasd",
	}

	return &RequestMetadata, configBytes, inputsBytes, nil

}

func randomID() string {
	return random(32)
}

func random(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type CacheData struct {
	ForwarderAddress string
}

func saveForwarderAddressToCache(forwarderAddress common.Address) error {
	cacheDir := ".cache"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data := CacheData{
		ForwarderAddress: forwarderAddress.String(),
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := os.WriteFile(filepath.Join(cacheDir, "forwarder.json"), b, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

func getForwarderAddressFromCache() (common.Address, error) {
	b, err := os.ReadFile(filepath.Join(".cache", "forwarder.json"))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to read cache file: %w", err)
	}

	var data CacheData
	if err := json.Unmarshal(b, &data); err != nil {
		return common.Address{}, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return common.HexToAddress(data.ForwarderAddress), nil
}
