package workflow_registry_syncer_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"

	"github.com/stretchr/testify/require"
)

func Test_SecretsWorker(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = syncer.NewWorkflowRegistryDS(db, lggr)

		giveTicker     = time.NewTicker(500 * time.Millisecond)
		giveSecretsURL = "https://original-url.com"
		donID          = uint32(1)
		giveWorkflow   = RegisterWorkflowCMD{
			Name:       "test-wf",
			DonID:      donID,
			Status:     uint8(1),
			SecretsURL: giveSecretsURL,
		}
		giveContents = "contents"
		wantContents = "updated contents"
		fetcherFn    = func(_ context.Context, _ string) ([]byte, error) {
			return []byte(wantContents), nil
		}
		contractName            = syncer.ContractName
		forceUpdateSecretsEvent = string(syncer.ForceUpdateSecretsEvent)
	)

	defer giveTicker.Stop()

	// fill ID with randomd data
	var giveID [32]byte
	_, err := rand.Read((giveID)[:])
	require.NoError(t, err)
	giveWorkflow.ID = giveID

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	lggr.Infof("deployed workflow registry at %s\n", wfRegistryAddr.Hex())

	// Build the ContractReader config
	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			contractName: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{forceUpdateSecretsEvent},
				},
				ContractABI: workflow_registry_wrapper.WorkflowRegistryABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					forceUpdateSecretsEvent: {
						ChainSpecificName: forceUpdateSecretsEvent,
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

	err = contractReader.Bind(ctx, []types.BoundContract{{Name: contractName, Address: wfRegistryAddr.Hex()}})
	require.NoError(t, err)

	// Seed the DB
	hash, err := crypto.Keccak256(append(backendTH.ContractsOwner.From[:], []byte(giveSecretsURL)...))
	require.NoError(t, err)
	giveHash := hex.EncodeToString(hash)

	gotID, err := orm.Create(ctx, giveSecretsURL, giveHash, giveContents)
	require.NoError(t, err)

	gotSecretsURL, err := orm.GetSecretsURLByID(ctx, gotID)
	require.NoError(t, err)
	require.Equal(t, giveSecretsURL, gotSecretsURL)

	// verify the DB
	contents, err := orm.GetContents(ctx, giveSecretsURL)
	require.NoError(t, err)
	require.Equal(t, contents, giveContents)

	// Create the worker
	worker := syncer.NewWorkflowRegistry(
		lggr,
		orm,
		contractReader,
		fetcherFn,
		wfRegistryAddr.Hex(),
		nil,
		nil,
		syncer.WithTicker(giveTicker.C),
	)

	servicetest.Run(t, worker)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	// generate a log event
	requestForceUpdateSecrets(t, backendTH, wfRegistryC, giveSecretsURL)

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		secrets, err := orm.GetContents(ctx, giveSecretsURL)
		lggr.Debugf("got secrets %v", secrets)
		require.NoError(t, err)
		return secrets == wantContents
	}, 5*time.Second, time.Second)
}

func updateAuthorizedAddress(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	addresses []common.Address,
	_ bool,
) {
	t.Helper()
	_, err := wfRegC.UpdateAuthorizedAddresses(th.ContractsOwner, addresses, true)
	require.NoError(t, err, "failed to update authorised addresses")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
	gotAddresses, err := wfRegC.GetAllAuthorizedAddresses(&bind.CallOpts{
		From: th.ContractsOwner.From,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, addresses, gotAddresses)
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
	gotDons, err := wfRegC.GetAllAllowedDONs(&bind.CallOpts{
		From: th.ContractsOwner.From,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, donIDs, gotDons)
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
