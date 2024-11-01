package syncer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	secretMocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/secrets/mocks"
	"golang.org/x/crypto/sha3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Keccak256HashToCommonHash(data []byte) common.Hash {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	var hashBytes [32]byte
	copy(hashBytes[:], hash.Sum(nil))
	return common.BytesToHash(hashBytes[:])
}

func Test_fetchSecrets(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		backend   = testutils.NewEVMBackendTH(t)
		updater   = secretMocks.NewUpdater[UpdateSecretsCommand](t)
		fetcherFn = func(_ context.Context, _ string) ([]byte, error) {
			return nil, assert.AnError
		}
		giveTimer    = make(chan struct{})
		contractName = "WorkflowRegistry"
		eventName    = "WorkflowForceUpdateSecretsRequestedV1"
	)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backend.ContractsOwner, backend.Backend)
	require.NoError(t, err)

	_ = wfRegistryAddr

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

	contractReaderConfigEncoded, err := json.Marshal(contractReaderCfg)
	require.NoError(t, err, "failed to marshal contract reader config")

	contractReader, err := backend.NewContractReader(ctx, t, contractReaderConfigEncoded)
	require.NoError(t, err)

	probes := fetchSecrets(
		ctx,
		giveTimer,
		contractReader,
		updater,
		fetcherFn,
	)

	// generate a log event
	requestForceUpdateSecrets(t, backend, wfRegistryC, "https://some-url.com")
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
	_, err := wfRegC.RequestForceUpdateSecrets(backend.ContractsOwner, Keccak256HashToCommonHash([]byte(secretsURL)).String())
	require.NoError(t, err, "failed to request force update secrets")
	backend.Backend.Commit()
}
