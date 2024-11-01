package syncer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	secretMocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/secrets/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_fetchSecrets(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		backendTH = testutils.NewEVMBackendTH(t)
		updater   = secretMocks.NewUpdater[UpdateSecretsCommand](t)
		fetcherFn = func(_ context.Context, _ string) ([]byte, error) {
			return nil, assert.AnError
		}
		giveTimer    = make(chan struct{})
		contractName = "WorkflowRegistry"
		eventName    = "WorkflowForceUpdateSecretsRequestedV1"
	)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend)
	require.NoError(t, err)
	lggr.Infof("deployed workflow registry at %s\n", wfRegistryAddr.Hex())
	backendTH.Backend.Commit()
	backendTH.Backend.Commit()

	// Bind the contract
	/* wfRegistryC, err := workflow_registry_wrapper.NewWorkflowRegistry(wfRegistryAddr, backend.Backend)
	require.NoError(t, err) */

	// Build the ContractReader config
	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			contractName: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{eventName},
				},
				ContractABI: workflow_registry_wrapper.WorkflowRegistryABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					eventName: {
						ChainSpecificName: eventName,
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

	contractReader, err := backendTH.NewContractReader(ctx, t, contractReaderCfgBytes)
	require.NoError(t, err)

	probes := fetchSecrets(
		ctx,
		giveTimer,
		contractReader,
		updater,
		fetcherFn,
	)

	// generate a log event
	requestForceUpdateSecrets(t, backendTH, wfRegistryC, "https://some-url.com")
	_ = wfRegistryC

	select {
	case <-time.After(3 * time.Second):
		t.Fatalf("failed to receive log")
	case l := <-probes.Logs:
		lggr.Infof("got log %+v\n", l)
	}
}

func requestForceUpdateSecrets(
	t *testing.T,
	backend *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	secretsURL string,
) {
	t.Helper()
	_, err := wfRegC.RequestForceUpdateSecrets(backend.ContractsOwner, secretsURL)
	require.NoError(t, err, "failed to request force update secrets")
	backend.Backend.Commit()
}
