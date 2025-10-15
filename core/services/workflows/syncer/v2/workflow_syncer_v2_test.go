package v2

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	coretestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	corecaps "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	wfstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	wfTypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

func Test_InitialStateSyncV2(t *testing.T) {
	lggr := logger.TestLogger(t)
	backendTH := testutils.NewEVMBackendTH(t)
	donID := uint32(1)
	donFamily := "A"

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)
	// setup the initial contract state
	updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily})
	updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily)

	// Add requests to ensure we go above the MaxResultsPerQuery
	activeAllowlistedRequestsCount := int(MaxResultsPerQuery + 1)
	expiryTimestamp := time.Now().Add(24 * time.Hour)
	for i := 0; i < activeAllowlistedRequestsCount; i++ {
		createSecretsRequestParams, marshalErr := json.Marshal(vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
						Key:       strconv.Itoa(i),
						Namespace: "active",
					},
					EncryptedValue: "encrypted-value",
				},
			},
		})
		require.NoError(t, marshalErr)

		allowlistRequest(t, backendTH, wfRegistryC, allowlistRequestParams{
			Request: jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				Params: (*json.RawMessage)(&createSecretsRequestParams),
			},
			Owner:           backendTH.ContractsOwner.From,
			ExpiryTimestamp: expiryTimestamp,
		})
	}

	// The number of workflows should be greater than the workflow registry contracts pagination limit to ensure
	// that the syncer will query the contract multiple times to get the full list of workflows
	numberWorkflows := 250
	for i := range numberWorkflows {
		var workflowID [32]byte
		_, err = rand.Read((workflowID)[:])
		require.NoError(t, err)
		workflow := RegisterWorkflowCMDV2{
			Name:      fmt.Sprintf("test-wf-%d", i),
			Tag:       "sometag",
			ID:        workflowID,
			Status:    WorkflowStatusActive,
			DonFamily: donFamily,
			BinaryURL: "someurl",
			KeepAlive: false,
		}
		workflow.ID = workflowID
		upsertWorkflowV2(t, backendTH, wfRegistryC, workflow)
	}
	testEventHandler := newTestEvtHandler(nil)

	// Create the worker
	worker, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID:       donID,
				Families: []string{donFamily},
			},
			err: nil,
		},
		NewEngineRegistry(),
	)
	require.NoError(t, err)

	servicetest.Run(t, worker)

	require.Eventually(t, func() bool {
		return len(testEventHandler.GetEvents()) == numberWorkflows
	}, tests.WaitTimeout(t), time.Second)

	for _, event := range testEventHandler.GetEvents() {
		assert.Equal(t, WorkflowActivated, event.Name)
	}

	assert.Len(t,
		worker.GetAllowlistedRequests(context.Background()),
		activeAllowlistedRequestsCount,
		"synced allowlisted requests do not match expectations",
	)
}

func Test_RegistrySyncer_SkipsEventsNotBelongingToDONV2(t *testing.T) {
	var (
		lggr      = logger.TestLogger(t)
		backendTH = testutils.NewEVMBackendTH(t)

		giveBinaryURL   = "https://original-url.com"
		donID           = uint32(1)
		donFamily1      = "A"
		donFamily2      = "B"
		skippedWorkflow = RegisterWorkflowCMDV2{
			Name:      "test-wf2",
			Status:    WorkflowStatusActive,
			BinaryURL: giveBinaryURL,
			Tag:       "sometag",
			DonFamily: donFamily2,
			KeepAlive: false,
		}
		giveWorkflow = RegisterWorkflowCMDV2{
			Name:      "test-wf",
			Status:    WorkflowStatusActive,
			BinaryURL: "someurl",
			Tag:       "sometag",
			DonFamily: donFamily1,
			KeepAlive: false,
		}
		wantContents = "updated contents"
	)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	from := [20]byte(backendTH.ContractsOwner.From)
	id, err := pkgworkflows.GenerateWorkflowID(from[:], "test-wf", []byte(wantContents), []byte(""), "")
	require.NoError(t, err)
	giveWorkflow.ID = id

	from = [20]byte(backendTH.ContractsOwner.From)
	id, err = pkgworkflows.GenerateWorkflowID(from[:], "test-wf", []byte(wantContents), []byte("dummy config"), "")
	require.NoError(t, err)
	skippedWorkflow.ID = id

	handler := newTestEvtHandler(nil)

	worker, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID:       donID,
				Families: []string{donFamily1},
			},
			err: nil,
		},
		NewEngineRegistry(),
	)
	require.NoError(t, err)

	// setup the initial contract state
	updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily1)
	updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily1, donFamily2})

	servicetest.Run(t, worker)

	// generate a log event
	upsertWorkflowV2(t, backendTH, wfRegistryC, skippedWorkflow)
	upsertWorkflowV2(t, backendTH, wfRegistryC, giveWorkflow)

	require.Eventually(t, func() bool {
		// we process events in order, and should only receive 1 event
		// the first is skipped as it belongs to another don.
		return len(handler.GetEvents()) == 1
	}, tests.WaitTimeout(t), time.Second)
}

func Test_RegistrySyncer_WorkflowRegistered_InitiallyPausedV2(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		emitter   = custmsg.NewLabeler()
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = artifacts.NewWorkflowRegistryDS(db, lggr)

		giveBinaryURL = "https://original-url.com"
		donID         = uint32(1)
		donFamily     = "A"
		giveWorkflow  = RegisterWorkflowCMDV2{
			Name:      "test-wf",
			Status:    WorkflowStatusPaused,
			BinaryURL: giveBinaryURL,
			Tag:       "sometag",
			DonFamily: donFamily,
			KeepAlive: false,
		}
		wantContents = "updated contents"
		fetcherFn    = func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte(base64.StdEncoding.EncodeToString([]byte(wantContents))), nil
		}
		retrieverFn = func(ctx context.Context, req *storage_service.DownloadArtifactRequest) (string, error) {
			return "", nil
		}
		workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	from := [20]byte(backendTH.ContractsOwner.From)
	id, err := pkgworkflows.GenerateWorkflowID(from[:], "test-wf", []byte(wantContents), []byte(""), "")
	require.NoError(t, err)
	giveWorkflow.ID = id

	er := NewEngineRegistry()
	limiters, err := v2.NewLimiters(limits.Factory{}, nil)
	require.NoError(t, err)
	rl, err := ratelimiter.NewRateLimiter(rlConfig, limits.Factory{})
	require.NoError(t, err)

	wl, err := syncerlimiter.NewWorkflowLimits(lggr, wlConfig, limits.Factory{})
	require.NoError(t, err)
	wfStore := wfstore.NewInMemoryStore(lggr, clockwork.NewFakeClock())
	capRegistry := corecaps.NewRegistry(lggr)
	capRegistry.SetLocalRegistry(&corecaps.TestMetadataRegistry{})
	store, err := artifacts.NewStore(lggr, orm, fetcherFn, retrieverFn, clockwork.NewFakeClock(), workflowkey.Key{}, emitter, artifacts.WithConfig(artifacts.StoreConfig{
		ArtifactStorageHost: "storage.chain.link",
	}))
	require.NoError(t, err)

	handler, err := NewEventHandler(lggr, wfStore, nil, true, capRegistry, er, emitter, limiters, rl, wl, store, workflowEncryptionKey)
	require.NoError(t, err)

	worker, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID:       donID,
				Families: []string{donFamily},
			},
			err: nil,
		},
		er,
	)
	require.NoError(t, err)

	// setup the initial contract state
	updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily)
	updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily})

	servicetest.Run(t, worker)

	// generate a log event
	upsertWorkflowV2(t, backendTH, wfRegistryC, giveWorkflow)

	// Paused workflows should generate no events
	time.Sleep(5 * time.Second)
	_, ok := er.Get(wfTypes.WorkflowID(id))
	require.False(t, ok)
	_, err = orm.GetWorkflowSpec(ctx, wfTypes.WorkflowID(id).Hex())
	require.ErrorContains(t, err, "sql: no rows in result set")
}

func Test_RegistrySyncer_WorkflowRegistered_InitiallyActivatedV2(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		emitter   = custmsg.NewLabeler()
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = artifacts.NewWorkflowRegistryDS(db, lggr)

		giveBinaryURL = "https://original-url.com"
		donID         = uint32(1)
		donFamily     = "A"
		giveWorkflow  = RegisterWorkflowCMDV2{
			Name:      "test-wf",
			Status:    WorkflowStatusActive,
			BinaryURL: giveBinaryURL,
			Tag:       "sometag",
			DonFamily: donFamily,
			KeepAlive: false,
		}
		wantContents = "updated contents"
		fetcherFn    = func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte(base64.StdEncoding.EncodeToString([]byte(wantContents))), nil
		}
		retrieverFn = func(ctx context.Context, req *storage_service.DownloadArtifactRequest) (string, error) {
			return "", nil
		}
		workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	from := [20]byte(backendTH.ContractsOwner.From)
	id, err := pkgworkflows.GenerateWorkflowID(from[:], "test-wf", []byte(wantContents), []byte(""), "")
	require.NoError(t, err)
	giveWorkflow.ID = id

	er := NewEngineRegistry()
	limiters, err := v2.NewLimiters(limits.Factory{}, nil)
	require.NoError(t, err)
	rl, err := ratelimiter.NewRateLimiter(rlConfig, limits.Factory{})
	require.NoError(t, err)
	wl, err := syncerlimiter.NewWorkflowLimits(lggr, wlConfig, limits.Factory{})
	require.NoError(t, err)
	wfStore := wfstore.NewInMemoryStore(lggr, clockwork.NewFakeClock())
	capRegistry := corecaps.NewRegistry(lggr)
	capRegistry.SetLocalRegistry(&corecaps.TestMetadataRegistry{})
	store, err := artifacts.NewStore(lggr, orm, fetcherFn, retrieverFn, clockwork.NewFakeClock(), workflowkey.Key{}, emitter, artifacts.WithConfig(artifacts.StoreConfig{
		ArtifactStorageHost: "storage.chain.link",
	}))
	require.NoError(t, err)

	handler, err := NewEventHandler(lggr, wfStore, nil, true, capRegistry, er,
		emitter, limiters, rl, wl, store, workflowEncryptionKey, WithStaticEngine(&mockService{}))
	require.NoError(t, err)

	worker, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID:       donID,
				Families: []string{donFamily},
			},
			err: nil,
		},
		er,
	)
	require.NoError(t, err)

	// setup the initial contract state
	updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily)
	updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily})

	servicetest.Run(t, worker)

	// generate a log event
	upsertWorkflowV2(t, backendTH, wfRegistryC, giveWorkflow)

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		_, ok := er.Get(wfTypes.WorkflowID(id))
		if !ok {
			return false
		}

		_, err = orm.GetWorkflowSpec(ctx, wfTypes.WorkflowID(id).Hex())
		return err == nil
	}, tests.WaitTimeout(t), time.Second)
}

func Test_StratReconciliation_InitialStateSyncV2(t *testing.T) {
	t.Run("with heavy load", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		backendTH := testutils.NewEVMBackendTH(t)
		donID := uint32(1)
		donFamily := "A"

		// Deploy a test workflow_registry
		wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
		backendTH.Backend.Commit()
		require.NoError(t, err)

		// setup the initial contract state
		updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily)
		updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily})

		// Use a high number of workflows
		// Tested up to 7_000
		numberWorkflows := 1_000
		for i := range numberWorkflows {
			var workflowID [32]byte
			_, err = rand.Read((workflowID)[:])
			require.NoError(t, err)
			workflow := RegisterWorkflowCMDV2{
				Name:      fmt.Sprintf("test-wf-%d", i),
				Status:    WorkflowStatusActive,
				BinaryURL: "someurl",
				Tag:       "sometag",
				DonFamily: donFamily,
				KeepAlive: false,
			}
			workflow.ID = workflowID
			upsertWorkflowV2(t, backendTH, wfRegistryC, workflow)
		}

		testEventHandler := newTestEvtHandler(nil)

		// Create the worker
		worker, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return backendTH.NewContractReader(ctx, t, bytes)
			},
			wfRegistryAddr.Hex(),
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			testEventHandler,
			&testDonNotifier{
				don: capabilities.DON{
					ID:       donID,
					Families: []string{donFamily},
				},
				err: nil,
			},
			NewEngineRegistry(),
			WithRetryInterval(1*time.Second),
		)
		require.NoError(t, err)

		servicetest.Run(t, worker)

		require.Eventually(t, func() bool {
			return len(testEventHandler.GetEvents()) == numberWorkflows
		}, 30*time.Second, 1*time.Second)

		for _, event := range testEventHandler.GetEvents() {
			assert.Equal(t, WorkflowActivated, event.Name)
		}
	})
}

func Test_StratReconciliation_RetriesWithBackoffV2(t *testing.T) {
	lggr := logger.TestLogger(t)
	backendTH := testutils.NewEVMBackendTH(t)
	donID := uint32(1)
	donFamily := "A"

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	// setup the initial contract state
	updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily)
	updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily})

	var workflowID [32]byte
	_, err = rand.Read((workflowID)[:])
	require.NoError(t, err)
	workflow := RegisterWorkflowCMDV2{
		Name:      "test-wf",
		Status:    WorkflowStatusActive,
		BinaryURL: "someurl",
		Tag:       "sometag",
		DonFamily: donFamily,
		KeepAlive: false,
	}
	workflow.ID = workflowID
	upsertWorkflowV2(t, backendTH, wfRegistryC, workflow)

	var retryCount int
	testEventHandler := newTestEvtHandler(func() error {
		if retryCount <= 1 {
			retryCount++
			return errors.New("error handling event")
		}
		return nil
	})

	// Create the worker
	worker, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID:       donID,
				Families: []string{donFamily},
			},
			err: nil,
		},
		NewEngineRegistry(),
		WithRetryInterval(1*time.Second),
	)
	require.NoError(t, err)

	servicetest.Run(t, worker)

	require.Eventually(t, func() bool {
		return len(testEventHandler.GetEvents()) == 1
	}, 30*time.Second, 1*time.Second)

	event := testEventHandler.GetEvents()[0]
	assert.Equal(t, WorkflowActivated, event.Name)

	assert.Equal(t, 1, retryCount)
}

// Links owner account to Workflow Registry contract
func updateAuthorizedAddressV2(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper_v2.WorkflowRegistry,
	ownerAddress common.Address,
	donFamily string,
) {
	t.Helper()

	// First, allow signer
	_, err := wfRegC.UpdateAllowedSigners(th.ContractsOwner, []common.Address{ownerAddress}, true)
	require.NoError(t, err)

	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()

	// Double check that signer has been allowed
	isAllowed, err := wfRegC.IsAllowedSigner(&bind.CallOpts{
		From: th.ContractsOwner.From,
	}, ownerAddress)
	require.NoError(t, err)
	require.True(t, isAllowed)

	requestTypeLink := uint8(0)
	typeAndVersion := "WorkflowRegistry 2.0.0"
	chainID, err := th.Backend.Client().ChainID(t.Context())
	require.NoError(t, err)
	linkContract := wfRegC.Address()
	validityTimestamp := big.NewInt(time.Now().UTC().Add(1 * time.Hour).Unix()) // block timestamp + 1 hour
	proof := generateOwnershipProofHash(ownerAddress.String(), donFamily, "1")

	// Prepare a list of ABI arguments in the exact order as expected by the Solidity contract
	arguments, err := prepareABIArguments()
	require.NoError(t, err)
	packed, err := arguments.Pack(
		requestTypeLink,
		ownerAddress,
		chainID,
		linkContract,
		typeAndVersion,
		validityTimestamp,
		proof,
	)
	require.NoError(t, err)

	// Hash the concatenated result using SHA256, Solidity contract will use keccak256()
	hash, err := crypto.Keccak256(packed)
	require.NoError(t, err)

	// Prepare a message that can be verified in a Solidity contract.
	// For a signature to be recoverable, it must follow the EIP-191 standard.
	// The message must be prefixed with "\x19Ethereum Signed Message:\n" followed by the length of the message.
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	messageHash, err := crypto.Keccak256([]byte(prefixedMessage))
	require.NoError(t, err)

	// Sign the message
	signature, err := th.ContractsOwnerSign(messageHash)
	require.NoError(t, err)
	// For Ethereum - add 27 to the recovery ID
	signature[64] += 27

	_, err = wfRegC.LinkOwner(th.ContractsOwner, validityTimestamp, proof, signature)
	require.NoError(t, err)

	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()

	_, err = HandleRevertData(err)
	require.ErrorIs(t, err, ErrCouldNotDecode)

	// Double check that owner has been linked
	isLinked, err := wfRegC.IsOwnerLinked(&bind.CallOpts{
		From: th.ContractsOwner.From,
	}, ownerAddress)
	require.NoError(t, err)
	require.True(t, isLinked)
}

func updateAllowedDONsV2(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper_v2.WorkflowRegistry,
	donFamilies []string,
) {
	t.Helper()
	for _, donFamily := range donFamilies {
		workflowLimit := uint32(10000)
		userLimit := uint32(1000)
		_, err := wfRegC.SetDONLimit(th.ContractsOwner, donFamily, workflowLimit, userLimit)
		require.NoError(t, err, "failed to update DONs")

		th.Backend.Commit()
		th.Backend.Commit()
		th.Backend.Commit()

		// Double check that DON has limit set
		maxWorkflows, err := wfRegC.GetMaxWorkflowsPerDON(&bind.CallOpts{
			From: th.ContractsOwner.From,
		}, donFamily)
		require.NoError(t, err)
		require.Equal(t, workflowLimit, maxWorkflows.MaxWorkflows, "max workflows mismatch")
	}
}

type RegisterWorkflowCMDV2 struct {
	Name       string   // Required
	Tag        string   // Required
	ID         [32]byte // Required
	Status     uint8    // Required
	DonFamily  string   // Required, must match a DON family that has had a limit set
	BinaryURL  string   // Required
	ConfigURL  string   // Optional
	Attributes []byte   // Optional
	KeepAlive  bool     // Optional, will default to false
}

func upsertWorkflowV2(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper_v2.WorkflowRegistry,
	input RegisterWorkflowCMDV2,
) {
	t.Helper()
	_, err := wfRegC.UpsertWorkflow(
		th.ContractsOwner,
		input.Name,
		input.Tag,
		input.ID,
		input.Status,
		input.DonFamily,
		input.BinaryURL,
		input.ConfigURL,
		input.Attributes,
		input.KeepAlive,
	)
	require.NoError(t, err, "failed to register workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()

	workflow, err := wfRegC.GetWorkflow(&bind.CallOpts{
		From: th.ContractsOwner.From,
	},
		th.ContractsOwner.From, input.Name, input.Tag,
	)
	require.NoError(t, err, "failed to register workflow")
	require.Equal(t, input.Name, workflow.WorkflowName, "workflow name mismatch")
	require.Equal(t, input.Tag, workflow.Tag, "workflow tag mismatch")
	require.Equal(t, input.ID, workflow.WorkflowId, "workflow ID mismatch")
	require.Equal(t, input.Status, workflow.Status, "workflow status mismatch")
	require.Equal(t, input.BinaryURL, workflow.BinaryUrl, "workflow binary URL mismatch")
	require.Equal(t, input.ConfigURL, workflow.ConfigUrl, "workflow config URL mismatch")
	// From contract comment:
	// For ACTIVE workflows this will resolve to the correct DON label.
	// For PAUSED/never‑assigned workflows the label is the empty string.
	if input.Status == WorkflowStatusActive {
		require.Equal(t, input.DonFamily, workflow.DonFamily, "workflow DON family mismatch")
	}
}

// Generates a hash for the ownership proof based on the workflow owner address, organization ID, and nonce.
func generateOwnershipProofHash(
	workflowOwnerAddress, organizationID, nonce string,
) [32]byte {
	data := workflowOwnerAddress + organizationID + nonce
	hash := sha256.Sum256([]byte(data))
	return hash
}

type allowlistRequestParams struct {
	Request         jsonrpc.Request[json.RawMessage]
	Owner           common.Address
	ExpiryTimestamp time.Time
}

func allowlistRequest(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper_v2.WorkflowRegistry,
	input allowlistRequestParams,
) {
	t.Helper()
	totalAllowlistedRequestsBefore, err := wfRegC.TotalAllowlistedRequests(&bind.CallOpts{
		From: th.ContractsOwner.From,
	})
	require.NoError(t, err, "failed to get total allowlisted requests")

	requestDigest, err := input.Request.Digest()
	require.NoError(t, err)
	requestDigestBytes, err := hex.DecodeString(requestDigest)
	require.NoError(t, err)

	_, err = wfRegC.AllowlistRequest(
		th.ContractsOwner,
		[32]byte(requestDigestBytes),
		uint32(input.ExpiryTimestamp.Unix()), //nolint:gosec // safe conversion
	)
	require.NoError(t, err, "failed to register allowlisted request")
	th.Backend.Commit()

	totalAllowlistedRequestsAfter, err := wfRegC.TotalAllowlistedRequests(&bind.CallOpts{
		From: th.ContractsOwner.From,
	})
	require.NoError(t, err, "failed to get total allowlisted requests")
	require.Equal(t, totalAllowlistedRequestsBefore.Uint64()+1, totalAllowlistedRequestsAfter.Uint64(), "total allowlisted requests mismatch")
}

// Prepare the ABI arguments, in the exact order as expected by the Solidity contract.
func prepareABIArguments() (*abi.Arguments, error) {
	arguments := abi.Arguments{}

	uint8Type, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create uint8 type: %w", err)
	}

	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create address type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create bytes32 type: %w", err)
	}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create uint256 type: %w", err)
	}

	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create string type: %w", err)
	}

	arguments = append(arguments, abi.Argument{Type: uint8Type})   // request type
	arguments = append(arguments, abi.Argument{Type: addressType}) // owner address
	arguments = append(arguments, abi.Argument{Type: uint256Type}) // chain ID
	arguments = append(arguments, abi.Argument{Type: addressType}) // address of the contract
	arguments = append(arguments, abi.Argument{Type: stringType})  // version string
	arguments = append(arguments, abi.Argument{Type: uint256Type}) // validity timestamp
	arguments = append(arguments, abi.Argument{Type: bytes32Type}) // ownership proof hash

	return &arguments, nil
}
