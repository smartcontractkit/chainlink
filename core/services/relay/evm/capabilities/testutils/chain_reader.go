package testutils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commoncaps "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonvalues "github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// Test harness with EVM backend and chainlink core services like
// Log Poller and Head Tracker
type ContractReaderTH struct {
	BackendTH *EVMBackendTH

	LogEmitterAddress           *common.Address
	LogEmitterContract          *log_emitter.LogEmitter
	LogEmitterContractReader    commontypes.ContractReader
	LogEmitterRegRequest        commoncaps.TriggerRegistrationRequest
	LogEmitterContractReaderCfg []byte
}

// Creates a new test harness for Contract Reader tests
func NewContractReaderTH(t *testing.T) *ContractReaderTH {
	backendTH := NewEVMBackendTH(t)

	// Deploy a test contract LogEmitter for testing ContractReader
	logEmitterAddress, _, _, err :=
		log_emitter.DeployLogEmitter(backendTH.ContractsOwner, backendTH.Backend.Client())
	require.NoError(t, err)
	backendTH.Backend.Commit()
	logEmitter, err := log_emitter.NewLogEmitter(logEmitterAddress, backendTH.Backend.Client())
	require.NoError(t, err)

	// Create new contract reader
	reqConfig := logeventcap.Config{
		ContractName:      "LogEmitter",
		ContractAddress:   logEmitterAddress.Hex(),
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
	var contractReaderCfgMap logeventcap.ConfigContractReaderConfig
	err = json.Unmarshal(contractReaderCfgBytes, &contractReaderCfgMap)
	require.NoError(t, err)
	// Encode the config map as JSON to specify in the expected call in mocked object
	// The LogEventTrigger Capability receives a config map, encodes it and
	// calls NewContractReader with it
	contractReaderCfgBytes, err = json.Marshal(contractReaderCfgMap)
	require.NoError(t, err)

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
	contractReader, err := backendTH.NewContractReader(ctx, t, contractReaderCfgBytes)
	require.NoError(t, err)

	return &ContractReaderTH{
		BackendTH: backendTH,

		LogEmitterAddress:           &logEmitterAddress,
		LogEmitterContract:          logEmitter,
		LogEmitterContractReader:    contractReader,
		LogEmitterRegRequest:        req,
		LogEmitterContractReaderCfg: contractReaderCfgBytes,
	}
}

// Wait for a specific log to be emitted to a response channel by ChainReader
func WaitForLog(lggr logger.Logger, logCh <-chan commoncaps.TriggerResponse, timeout time.Duration) (
	*commoncaps.TriggerResponse, map[string]any, error) {
	select {
	case <-time.After(timeout):
		return nil, nil, fmt.Errorf("timeout waiting for Log1 event from ContractReader")
	case log := <-logCh:
		lggr.Infow("Received log from ContractReader", "event", log.Event.ID)
		if log.Err != nil {
			return nil, nil, fmt.Errorf("error listening for Log1 event from ContractReader: %v", log.Err)
		}
		v := make(map[string]any)
		err := log.Event.Outputs.UnwrapTo(&v)
		if err != nil {
			return nil, nil, fmt.Errorf("error unwrapping log to map: (log %v) %v", log.Event.Outputs, log.Err)
		}
		return &log, v, nil
	}
}

// Get the string value of a key from a generic map[string]any
func GetStrVal(m map[string]any, k string) (string, error) {
	v, ok := m[k]
	if !ok {
		return "", fmt.Errorf("key %s not found", k)
	}
	vstr, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key %s not a string (%T)", k, v)
	}
	return vstr, nil
}

// Get int value of a key from a generic map[string]any
func GetBigIntVal(m map[string]any, k string) (*big.Int, error) {
	v, ok := m[k]
	if !ok {
		return nil, fmt.Errorf("key %s not found", k)
	}
	val, ok := v.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("key %s not a *big.Int (%T)", k, v)
	}
	return val, nil
}

// Get the int value from a map[string]map[string]any
func GetBigIntValL2(m map[string]any, level1Key string, level2Key string) (*big.Int, error) {
	v, ok := m[level1Key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", level1Key)
	}
	level2Map, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("key %s not a map[string]any (%T)", level1Key, v)
	}
	return GetBigIntVal(level2Map, level2Key)
}
