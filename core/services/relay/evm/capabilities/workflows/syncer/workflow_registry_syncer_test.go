package workflow_registry_syncer_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"

	"github.com/stretchr/testify/require"
)

func Test_SecretsWorker(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = syncer.NewWorkflowRegistryDS(db, lggr)

		giveTicker     = signalers.MakeTicker(ctx.Done(), 500*time.Millisecond)
		giveSecretsURL = "https://original-url.com"
		giveHash       = syncer.Keccak256Hash(giveSecretsURL)
		giveWorkflow   = RegisterWorkflowCMD{
			Name:       "test-wf",
			DonID:      uint32(1),
			Status:     uint8(0),
			SecretsURL: giveSecretsURL,
		}
		giveContents = "contents"
		wantContents = "updated contents"
		fetcherFn    = func(_ context.Context, _ string) ([]byte, error) {
			return []byte(wantContents), nil
		}
		contractName = "WorkflowRegistry"
		eventName    = "WorkflowForceUpdateSecretsRequestedV1"
		giveCfg      = syncer.ContractEventPollerConfig{
			ContractName:      contractName,
			ContractEventName: eventName,
			QueryCount:        20,
		}
	)

	// fill ID with randomd data
	var giveID [32]byte
	rand.Read((giveID)[:])
	giveWorkflow.ID = giveID

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend)
	require.NoError(t, err)

	lggr.Infof("deployed workflow registry at %s\n", wfRegistryAddr.Hex())
	giveCfg.ContractAddress = wfRegistryAddr.Hex()

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

	contractReaderCfgBytes, err := json.Marshal(contractReaderCfg)
	require.NoError(t, err)

	contractReader, err := backendTH.NewContractReader(ctx, t, contractReaderCfgBytes)
	require.NoError(t, err)

	giveCfg.StartBlockNum = uint64(0)

	// Seed the DB
	_, err = orm.Update(ctx, giveSecretsURL, giveContents)
	require.NoError(t, err)

	worker := syncer.NewWorkflowRegistry(
		lggr,
		orm,
		contractReader,
		fetcherFn,
		giveCfg,
		syncer.WithTicker(giveTicker),
	)

	// generate a log event
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{1}, true)
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	servicetest.Run(t, worker)

	go func() {
		for {
			<-time.After(time.Second)
			requestForceUpdateSecrets(t, backendTH, wfRegistryC, giveSecretsURL)
			backendTH.Backend.Commit()
		}
	}()

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		secrets, err := orm.GetArtifactByHash(ctx, giveHash)
		lggr.Debugf("got secrets %v", secrets)
		require.NoError(t, err)
		return secrets.Contents == wantContents
	}, 5*time.Second, time.Second)
}

func updateAuthorizedAddress(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	addresses []common.Address,
	allowed bool,
) {
	t.Helper()
	_, err := wfRegC.UpdateAuthorizedAddresses(th.ContractsOwner, addresses, allowed)
	require.NoError(t, err, "failed to update authorised addresses")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func updateAllowedDONs(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	donIDs []uint32,
	allowed bool,
) {
	t.Helper()
	_, err := wfRegC.UpdateAllowedDONs(th.ContractsOwner, donIDs, allowed)
	require.NoError(t, err, "failed to update DONs")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

type RegisterWorkflowCMD struct {
	Name       string
	ID         [32]byte
	DonID      uint32
	Status     uint8
	BinaryURL  string
	ConfigURL  string
	SecretsURL string
}

func registerWorkflow(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	input RegisterWorkflowCMD,
) {
	t.Helper()
	_, err := wfRegC.RegisterWorkflow(th.ContractsOwner, input.Name, input.ID, input.DonID,
		input.Status, input.BinaryURL, input.ConfigURL, input.SecretsURL)
	require.NoError(t, err, "failed to register workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func requestForceUpdateSecrets(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	secretsURL string,
) {
	_, err := wfRegC.RequestForceUpdateSecrets(th.ContractsOwner, secretsURL)
	require.NoError(t, err)
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func miner(t *testing.T, tb *testutils.EVMBackendTH, d time.Duration) {
	t.Helper()

	ctx := coretestutils.Context(t)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(d).C:
			tb.Backend.Commit()
		}
	}
}
