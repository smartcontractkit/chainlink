package logevent_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	commoncaps "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	commonvalues "github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/test-go/testify/mock"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// Test for Log Event Trigger Capability happy path for EVM
func TestLogEventTriggerEVMHappyPath(t *testing.T) {
	th := testutils.NewEVMLOOPCapabilityTH(t)
	logEventConfig := logevent.LogEventConfig{
		ChainId:        th.ChainID.Uint64(),
		Network:        "evm",
		LookbackBlocks: 1000,
		PollPeriod:     500,
	}

	// Create new contract reader
	reqConfig := logevent.RequestConfig{
		ContractName:      "LogEmitter",
		ContractAddress:   th.LogEmitterAddress.Hex(),
		ContractEventName: "Log1",
	}
	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			reqConfig.ContractName: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{reqConfig.ContractEventName},
				},
				ContractABI: log_emitter.LogEmitterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					reqConfig.ContractEventName: {
						ChainSpecificName: reqConfig.ContractEventName,
						ReadType:          evmtypes.Event,
					},
				},
			},
		},
	}

	// Encode contractReaderConfig as JSON and decode it into a map[string]any for
	// the capability request config. Log Event Trigger capability takes in a
	// []byte as ContractReaderConfig to not depend on evm ChainReaderConfig type
	// and be chain agnostic
	contractReaderCfgBytes, err := json.Marshal(contractReaderCfg)
	require.NoError(t, err)
	contractReaderCfgMap := make(map[string]any)
	err = json.Unmarshal(contractReaderCfgBytes, &contractReaderCfgMap)
	require.NoError(t, err)
	// Encode the config map as JSON to specify in the expected call in mocked object
	// The LogEventTrigger Capability receives a config map, encodes it and
	// calls NewContractReader with it
	contractReaderCfgBytes, _ = json.Marshal(contractReaderCfgMap)

	reqConfig.ContractReaderConfig = contractReaderCfgMap

	config, err := commonvalues.WrapMap(reqConfig)
	require.NoError(t, err)
	req := commoncaps.TriggerRegistrationRequest{
		TriggerID: "logeventtrigger_log1",
		Config:    config,
		Metadata: commoncaps.RequestMetadata{
			ReferenceID: "logeventtrigger",
		},
	}

	// Create a new contract reader to return from mock relayer
	ctx := coretestutils.Context(t)
	contractReader, err := th.NewContractReader(t, ctx, contractReaderCfgBytes)
	require.NoError(t, err)

	// Fetch latest head from simulated backend to return from mock relayer
	height, err := th.EVMClient.LatestBlockHeight(ctx)
	require.NoError(t, err)
	block, err := th.EVMClient.BlockByNumber(ctx, height)
	require.NoError(t, err)

	// Mock relayer to return a New ContractReader instead of gRPC client of a ContractReader
	relayer := commonmocks.NewRelayer(t)
	relayer.On("NewContractReader", mock.Anything, contractReaderCfgBytes).Return(contractReader, nil).Once()
	relayer.On("LatestHead", mock.Anything).Return(commontypes.Head{
		Height:    height.String(),
		Hash:      block.Hash().Bytes(),
		Timestamp: block.Time(),
	}, nil).Once()

	// Create Log Event Trigger Service and register trigger
	logEventTriggerService := logevent.NewLogEventTriggerService(logevent.Params{
		Logger:         th.Lggr,
		Relayer:        relayer,
		LogEventConfig: logEventConfig,
	})
	log1Ch, err := logEventTriggerService.RegisterTrigger(ctx, req)
	require.NoError(t, err)

	// Send a blockchain transaction that emits logs
	go func() {
		_, err :=
			th.LogEmitterContract.EmitLog1(th.ContractsOwner, []*big.Int{big.NewInt(10)})
		require.NoError(t, err)
		th.Backend.Commit()
		th.Backend.Commit()
		th.Backend.Commit()
	}()

	// Wait for logs with a timeout
	timeout := 5 * time.Second
	for {
		select {
		case <-time.After(timeout):
			require.NoError(t, fmt.Errorf("Timeout waiting for Log1 event from ContractReader"))
		case log1 := <-log1Ch:
			require.NoError(t, log1.Err, "error listening for Log1 event from ContractReader")
			require.NotNil(t, log1.Event.Outputs)
		}
	}
}
