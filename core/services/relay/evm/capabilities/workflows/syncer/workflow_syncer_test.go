package workflow_registry_syncer_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	rand2 "math/rand/v2"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"

	"github.com/stretchr/testify/require"

	crypto2 "github.com/ethereum/go-ethereum/crypto"
)

var rlConfig = ratelimiter.Config{
	GlobalRPS:      1000.0,
	GlobalBurst:    1000,
	PerSenderRPS:   30.0,
	PerSenderBurst: 30,
}

var wlConfig = syncerlimiter.Config{
	Global:   200,
	PerOwner: 200,
}

type testEvtHandler struct {
	events []syncer.Event
	mux    sync.Mutex
	errFn  func() error
}

func (m *testEvtHandler) Close() error { return nil }

func (m *testEvtHandler) Handle(ctx context.Context, event syncer.Event) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.events = append(m.events, event)
	if m.errFn != nil {
		return m.errFn()
	}
	return nil
}

func (m *testEvtHandler) ClearEvents() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.events = make([]syncer.Event, 0)
}

func (m *testEvtHandler) GetEvents() []syncer.Event {
	m.mux.Lock()
	defer m.mux.Unlock()

	eventsCopy := make([]syncer.Event, len(m.events))
	copy(eventsCopy, m.events)

	return eventsCopy
}

func newTestEvtHandler(errFn func() error) *testEvtHandler {
	return &testEvtHandler{
		errFn:  errFn,
		events: make([]syncer.Event, 0),
	}
}

type testWorkflowRegistryContractLoader struct {
}

type testDonNotifier struct {
	don capabilities.DON
	err error
}

func (t *testDonNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	return t.don, t.err
}

func (m *testWorkflowRegistryContractLoader) LoadWorkflows(ctx context.Context, don capabilities.DON) (*types.Head, error) {
	return &types.Head{
		Height:    "0",
		Hash:      nil,
		Timestamp: 0,
	}, nil
}

func Test_EventHandlerStateSync(t *testing.T) {
	lggr := logger.TestLogger(t)
	backendTH := testutils.NewEVMBackendTH(t)
	donID := uint32(1)

	eventPollTicker := time.NewTicker(50 * time.Millisecond)
	defer eventPollTicker.Stop()

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	// Create some initial static state
	numberWorkflows := 20
	for i := 0; i < numberWorkflows; i++ {
		var workflowID [32]byte
		_, err = rand.Read((workflowID)[:])
		require.NoError(t, err)
		workflow := RegisterWorkflowCMD{
			Name:       fmt.Sprintf("test-wf-%d", i),
			DonID:      donID,
			Status:     uint8(1),
			SecretsURL: "someurl",
			BinaryURL:  "someurl",
		}
		workflow.ID = workflowID
		registerWorkflow(t, backendTH, wfRegistryC, workflow)
	}

	testEventHandler := newTestEvtHandler(nil)

	// Create the registry
	registry, err := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.Config{
			QueryCount:   20,
			SyncStrategy: syncer.SyncStrategyEvent,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.NewEngineRegistry(),
		syncer.WithTicker(eventPollTicker.C),
	)
	require.NoError(t, err)

	servicetest.Run(t, registry)

	require.Eventually(t, func() bool {
		numEvents := len(testEventHandler.GetEvents())
		return numEvents == numberWorkflows
	}, tests.WaitTimeout(t), time.Second)

	for _, event := range testEventHandler.GetEvents() {
		assert.Equal(t, syncer.WorkflowRegisteredEvent, event.EventType)
	}

	testEventHandler.ClearEvents()

	// Create different event types for a number of workflows and confirm that the event handler processes them in order
	numberOfEventCycles := 50
	for i := 0; i < numberOfEventCycles; i++ {
		var workflowID [32]byte
		_, err = rand.Read((workflowID)[:])
		require.NoError(t, err)
		workflow := RegisterWorkflowCMD{
			Name:       "test-wf-register-event",
			DonID:      donID,
			Status:     uint8(1),
			SecretsURL: "",
			BinaryURL:  "someurl",
		}
		workflow.ID = workflowID

		// Generate events of different types with some jitter
		registerWorkflow(t, backendTH, wfRegistryC, workflow)
		time.Sleep(time.Millisecond * time.Duration(rand2.IntN(10)))
		data := append(backendTH.ContractsOwner.From.Bytes(), []byte(workflow.Name)...)
		workflowKey := crypto2.Keccak256Hash(data)
		activateWorkflow(t, backendTH, wfRegistryC, workflowKey)
		time.Sleep(time.Millisecond * time.Duration(rand2.IntN(10)))
		pauseWorkflow(t, backendTH, wfRegistryC, workflowKey)
		time.Sleep(time.Millisecond * time.Duration(rand2.IntN(10)))
		var newWorkflowID [32]byte
		_, err = rand.Read((newWorkflowID)[:])
		require.NoError(t, err)
		updateWorkflow(t, backendTH, wfRegistryC, workflowKey, newWorkflowID, workflow.BinaryURL+"2", workflow.ConfigURL, workflow.SecretsURL)
		time.Sleep(time.Millisecond * time.Duration(rand2.IntN(10)))
		deleteWorkflow(t, backendTH, wfRegistryC, workflowKey)
	}

	// Confirm the expected number of events are received in the correct order
	require.Eventually(t, func() bool {
		events := testEventHandler.GetEvents()
		numEvents := len(events)
		expectedNumEvents := 5 * numberOfEventCycles

		if numEvents == expectedNumEvents {
			// verify the events are the expected types in the expected order
			for idx, event := range events {
				switch idx % 5 {
				case 0:
					assert.Equal(t, syncer.WorkflowRegisteredEvent, event.EventType)
				case 1:
					assert.Equal(t, syncer.WorkflowActivatedEvent, event.EventType)
				case 2:
					assert.Equal(t, syncer.WorkflowPausedEvent, event.EventType)
				case 3:
					assert.Equal(t, syncer.WorkflowUpdatedEvent, event.EventType)
				case 4:
					assert.Equal(t, syncer.WorkflowDeletedEvent, event.EventType)
				}
			}
			return true
		}

		return false
	}, tests.WaitTimeout(t), time.Second)
}
func Test_InitialStateSync(t *testing.T) {
	lggr := logger.TestLogger(t)
	backendTH := testutils.NewEVMBackendTH(t)
	donID := uint32(1)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)
	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)
	// The number of workflows should be greater than the workflow registry contracts pagination limit to ensure
	// that the syncer will query the contract multiple times to get the full list of workflows
	numberWorkflows := 250
	for i := 0; i < numberWorkflows; i++ {
		var workflowID [32]byte
		_, err = rand.Read((workflowID)[:])
		require.NoError(t, err)
		workflow := RegisterWorkflowCMD{
			Name:       fmt.Sprintf("test-wf-%d", i),
			DonID:      donID,
			Status:     uint8(1),
			SecretsURL: "someurl",
			BinaryURL:  "someurl",
		}
		workflow.ID = workflowID
		registerWorkflow(t, backendTH, wfRegistryC, workflow)
	}
	testEventHandler := newTestEvtHandler(nil)

	// Create the worker
	worker, err := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.Config{
			QueryCount:   20,
			SyncStrategy: syncer.SyncStrategyEvent,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.NewEngineRegistry(),
		syncer.WithTicker(make(chan time.Time)),
	)
	require.NoError(t, err)

	servicetest.Run(t, worker)

	require.Eventually(t, func() bool {
		return len(testEventHandler.GetEvents()) == numberWorkflows
	}, tests.WaitTimeout(t), time.Second)

	for _, event := range testEventHandler.GetEvents() {
		assert.Equal(t, syncer.WorkflowRegisteredEvent, event.EventType)
	}
}

func Test_SecretsWorker(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-732")
	tc := []struct {
		ss syncer.SyncStrategy
	}{
		{ss: syncer.SyncStrategyEvent},
		{ss: syncer.SyncStrategyReconciliation},
	}

	for _, tt := range tc {
		t.Run(string(tt.ss), func(t *testing.T) {
			var (
				ctx       = coretestutils.Context(t)
				lggr      = logger.TestLogger(t)
				emitter   = custmsg.NewLabeler()
				backendTH = testutils.NewEVMBackendTH(t)
				db        = pgtest.NewSqlxDB(t)
				orm       = artifacts.NewWorkflowRegistryDS(db, lggr)

				encryptionKey  = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
				workflowOwner  = backendTH.ContractsOwner.From.Hex()
				beforeContents = "contents"
				afterContents  = "updated contents"
				giveTicker     = time.NewTicker(500 * time.Millisecond)
				giveSecretsURL = "https://original-url.com"
				donID          = uint32(1)
				giveWorkflow   = RegisterWorkflowCMD{
					Name:       "test-wf",
					DonID:      donID,
					Status:     uint8(0),
					SecretsURL: giveSecretsURL,
					BinaryURL:  "someurl",
				}
			)

			beforeSecretsPayload := encryptSecrets(t, workflowOwner, map[string][]string{
				"SECRET_A": {beforeContents},
			}, encryptionKey)
			afterSecretsPayload := encryptSecrets(t, workflowOwner, map[string][]string{
				"SECRET_A": {afterContents},
			}, encryptionKey)
			fetcherFn := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
				return afterSecretsPayload, nil
			}

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

			// Seed the DB
			hash, err := crypto.Keccak256(append(backendTH.ContractsOwner.From[:], []byte(giveSecretsURL)...))
			require.NoError(t, err)
			giveHash := hex.EncodeToString(hash)

			gotID, err := orm.Create(ctx, giveSecretsURL, giveHash, string(beforeSecretsPayload))
			require.NoError(t, err)

			gotSecretsURL, err := orm.GetSecretsURLByID(ctx, gotID)
			require.NoError(t, err)
			require.Equal(t, giveSecretsURL, gotSecretsURL)

			// verify the DB
			contents, err := orm.GetContents(ctx, giveSecretsURL)
			require.NoError(t, err)
			require.Equal(t, string(beforeSecretsPayload), contents)
			rl, err := ratelimiter.NewRateLimiter(rlConfig)
			require.NoError(t, err)

			wl, err := syncerlimiter.NewWorkflowLimits(lggr, wlConfig)
			require.NoError(t, err)

			store := artifacts.NewStore(lggr, orm, fetcherFn, clockwork.NewFakeClock(), encryptionKey, emitter)

			engineRegistry := syncer.NewEngineRegistry()

			evtHandler, err := syncer.NewEventHandler(lggr, nil, nil, engineRegistry,
				emitter, rl, wl, store)
			require.NoError(t, err)
			handler := &testSecretsWorkEventHandler{
				wrappedHandler: evtHandler,
				registeredCh:   make(chan syncer.Event, 1),
			}

			worker, err := syncer.NewWorkflowRegistry(
				lggr,
				func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
					return backendTH.NewContractReader(ctx, t, bytes)
				},
				wfRegistryAddr.Hex(),
				syncer.Config{
					QueryCount:   20,
					SyncStrategy: tt.ss,
				},
				handler,
				&testDonNotifier{
					don: capabilities.DON{
						ID: donID,
					},
					err: nil,
				},
				engineRegistry,
				syncer.WithTicker(giveTicker.C),
			)
			require.NoError(t, err)

			// setup contract state to allow the secrets to be updated
			updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
			updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)
			registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

			servicetest.Run(t, worker)

			// wait for the workflow to be registered
			<-handler.registeredCh

			// generate a log event
			requestForceUpdateSecrets(t, backendTH, wfRegistryC, giveSecretsURL)

			// Require the secrets contents to eventually be updated
			require.Eventually(t, func() bool {
				secrets, err := orm.GetContents(ctx, giveSecretsURL)
				lggr.Debugf("got secrets %v", secrets)
				require.NoError(t, err)
				return secrets == string(afterSecretsPayload)
			}, tests.WaitTimeout(t), time.Second)
		})
	}
}

func Test_RegistrySyncer_SkipsEventsNotBelongingToDON(t *testing.T) {
	var (
		lggr      = logger.TestLogger(t)
		backendTH = testutils.NewEVMBackendTH(t)

		giveTicker      = time.NewTicker(500 * time.Millisecond)
		giveBinaryURL   = "https://original-url.com"
		donID           = uint32(1)
		otherDonID      = uint32(2)
		skippedWorkflow = RegisterWorkflowCMD{
			Name:      "test-wf2",
			DonID:     otherDonID,
			Status:    uint8(1),
			BinaryURL: giveBinaryURL,
		}
		giveWorkflow = RegisterWorkflowCMD{
			Name:      "test-wf",
			DonID:     donID,
			Status:    uint8(1),
			BinaryURL: "someurl",
		}
		wantContents = "updated contents"
	)

	defer giveTicker.Stop()

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
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

	worker, err := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.Config{
			QueryCount:   20,
			SyncStrategy: syncer.SyncStrategyEvent,
		},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.NewEngineRegistry(),
		syncer.WithTicker(giveTicker.C),
	)
	require.NoError(t, err)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID, otherDonID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	servicetest.Run(t, worker)

	// generate a log event
	registerWorkflow(t, backendTH, wfRegistryC, skippedWorkflow)
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	require.Eventually(t, func() bool {
		// we process events in order, and should only receive 1 event
		// the first is skipped as it belongs to another don.
		return len(handler.GetEvents()) == 1
	}, tests.WaitTimeout(t), time.Second)
}

func Test_RegistrySyncer_WorkflowRegistered_InitiallyPaused(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		emitter   = custmsg.NewLabeler()
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = artifacts.NewWorkflowRegistryDS(db, lggr)

		giveTicker    = time.NewTicker(500 * time.Millisecond)
		giveBinaryURL = "https://original-url.com"
		donID         = uint32(1)
		giveWorkflow  = RegisterWorkflowCMD{
			Name:      "test-wf",
			DonID:     donID,
			Status:    uint8(1),
			BinaryURL: giveBinaryURL,
		}
		wantContents = "updated contents"
		fetcherFn    = func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte(base64.StdEncoding.EncodeToString([]byte(wantContents))), nil
		}
	)

	defer giveTicker.Stop()

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	from := [20]byte(backendTH.ContractsOwner.From)
	id, err := pkgworkflows.GenerateWorkflowID(from[:], "test-wf", []byte(wantContents), []byte(""), "")
	require.NoError(t, err)
	giveWorkflow.ID = id

	er := syncer.NewEngineRegistry()
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)

	wl, err := syncerlimiter.NewWorkflowLimits(lggr, wlConfig)
	require.NoError(t, err)

	store := artifacts.NewStore(lggr, orm, fetcherFn, clockwork.NewFakeClock(), workflowkey.Key{}, emitter)

	handler, err := syncer.NewEventHandler(lggr, nil, nil, er, emitter, rl, wl, store)
	require.NoError(t, err)

	worker, err := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.Config{
			QueryCount:   20,
			SyncStrategy: syncer.SyncStrategyEvent,
		},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		er,
		syncer.WithTicker(giveTicker.C),
	)
	require.NoError(t, err)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	servicetest.Run(t, worker)

	// generate a log event
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		_, ok := er.Get(syncer.EngineRegistryKey{Owner: backendTH.ContractsOwner.From.Bytes(), Name: "test-wf"})
		if ok {
			return false
		}

		owner := strings.ToLower(backendTH.ContractsOwner.From.Hex()[2:])
		_, err := orm.GetWorkflowSpec(ctx, owner, "test-wf")
		return err == nil
	}, tests.WaitTimeout(t), time.Second)
}

type mockService struct{}

func (m *mockService) Start(context.Context) error { return nil }

func (m *mockService) Close() error { return nil }

func (m *mockService) HealthReport() map[string]error { return map[string]error{"svc": nil} }

func (m *mockService) Ready() error { return nil }

func (m *mockService) Name() string { return "svc" }

func Test_RegistrySyncer_WorkflowRegistered_InitiallyActivated(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		emitter   = custmsg.NewLabeler()
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = artifacts.NewWorkflowRegistryDS(db, lggr)

		giveTicker    = time.NewTicker(500 * time.Millisecond)
		giveBinaryURL = "https://original-url.com"
		donID         = uint32(1)
		giveWorkflow  = RegisterWorkflowCMD{
			Name:      "test-wf",
			DonID:     donID,
			Status:    uint8(0),
			BinaryURL: giveBinaryURL,
		}
		wantContents = "updated contents"
		fetcherFn    = func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte(base64.StdEncoding.EncodeToString([]byte(wantContents))), nil
		}
	)

	defer giveTicker.Stop()

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	from := [20]byte(backendTH.ContractsOwner.From)
	id, err := pkgworkflows.GenerateWorkflowID(from[:], "test-wf", []byte(wantContents), []byte(""), "")
	require.NoError(t, err)
	giveWorkflow.ID = id

	er := syncer.NewEngineRegistry()
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	wl, err := syncerlimiter.NewWorkflowLimits(lggr, wlConfig)
	require.NoError(t, err)

	store := artifacts.NewStore(lggr, orm, fetcherFn, clockwork.NewFakeClock(), workflowkey.Key{}, emitter)

	handler, err := syncer.NewEventHandler(lggr, nil, nil, er,
		emitter, rl, wl, store, syncer.WithStaticEngine(&mockService{}))
	require.NoError(t, err)

	worker, err := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.Config{
			QueryCount:   20,
			SyncStrategy: syncer.SyncStrategyEvent,
		},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		er,
		syncer.WithTicker(giveTicker.C),
	)
	require.NoError(t, err)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	servicetest.Run(t, worker)

	// generate a log event
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		_, ok := er.Get(syncer.EngineRegistryKey{Owner: backendTH.ContractsOwner.From.Bytes(), Name: "test-wf"})
		if !ok {
			return false
		}

		owner := strings.ToLower(backendTH.ContractsOwner.From.Hex()[2:])
		_, err = orm.GetWorkflowSpec(ctx, owner, "test-wf")
		return err == nil
	}, tests.WaitTimeout(t), time.Second)
}

func Test_StratReconciliation_InitialStateSync(t *testing.T) {
	t.Run("with heavy load", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		backendTH := testutils.NewEVMBackendTH(t)
		donID := uint32(1)

		// Deploy a test workflow_registry
		wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
		backendTH.Backend.Commit()
		require.NoError(t, err)

		// setup contract state to allow the secrets to be updated
		updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
		updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

		// Use a high number of workflows
		// Tested up to 7_000
		numberWorkflows := 1_000
		for i := 0; i < numberWorkflows; i++ {
			var workflowID [32]byte
			_, err = rand.Read((workflowID)[:])
			require.NoError(t, err)
			workflow := RegisterWorkflowCMD{
				Name:       fmt.Sprintf("test-wf-%d", i),
				DonID:      donID,
				Status:     uint8(0),
				SecretsURL: "someurl",
				BinaryURL:  "someurl",
			}
			workflow.ID = workflowID
			registerWorkflow(t, backendTH, wfRegistryC, workflow)
		}

		testEventHandler := newTestEvtHandler(nil)

		// Create the worker
		worker, err := syncer.NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
				return backendTH.NewContractReader(ctx, t, bytes)
			},
			wfRegistryAddr.Hex(),
			syncer.Config{
				QueryCount:   20,
				SyncStrategy: syncer.SyncStrategyReconciliation,
			},
			testEventHandler,
			&testDonNotifier{
				don: capabilities.DON{
					ID: donID,
				},
				err: nil,
			},
			syncer.NewEngineRegistry(),
			syncer.WithRetryInterval(1*time.Second),
		)
		require.NoError(t, err)

		servicetest.Run(t, worker)

		require.Eventually(t, func() bool {
			return len(testEventHandler.GetEvents()) == numberWorkflows
		}, 30*time.Second, 1*time.Second)

		for _, event := range testEventHandler.GetEvents() {
			assert.Equal(t, syncer.WorkflowRegisteredEvent, event.EventType)
		}
	})
}

func Test_StratReconciliation_RetriesWithBackoff(t *testing.T) {
	lggr := logger.TestLogger(t)
	backendTH := testutils.NewEVMBackendTH(t)
	donID := uint32(1)

	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	var workflowID [32]byte
	_, err = rand.Read((workflowID)[:])
	require.NoError(t, err)
	workflow := RegisterWorkflowCMD{
		Name:       "test-wf",
		DonID:      donID,
		Status:     uint8(0),
		SecretsURL: "someurl",
		BinaryURL:  "someurl",
	}
	workflow.ID = workflowID
	registerWorkflow(t, backendTH, wfRegistryC, workflow)

	var retryCount int
	testEventHandler := newTestEvtHandler(func() error {
		if retryCount <= 1 {
			retryCount++
			return errors.New("error handling event")
		}
		return nil

	})

	// Create the worker
	worker, err := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.Config{
			QueryCount:   20,
			SyncStrategy: syncer.SyncStrategyReconciliation,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.NewEngineRegistry(),
		syncer.WithRetryInterval(1*time.Second),
	)
	require.NoError(t, err)

	servicetest.Run(t, worker)

	require.Eventually(t, func() bool {
		return len(testEventHandler.GetEvents()) == 1
	}, 30*time.Second, 1*time.Second)

	event := testEventHandler.GetEvents()[0]
	assert.Equal(t, syncer.WorkflowRegisteredEvent, event.EventType)

	assert.Equal(t, 1, retryCount)
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

func activateWorkflow(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	workflowKey [32]byte,
) {
	t.Helper()
	_, err := wfRegC.ActivateWorkflow(th.ContractsOwner, workflowKey)
	require.NoError(t, err, "failed to activate workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func pauseWorkflow(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	workflowKey [32]byte,
) {
	t.Helper()
	_, err := wfRegC.PauseWorkflow(th.ContractsOwner, workflowKey)
	require.NoError(t, err, "failed to pause workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func deleteWorkflow(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	workflowKey [32]byte,
) {
	t.Helper()
	_, err := wfRegC.DeleteWorkflow(th.ContractsOwner, workflowKey)
	require.NoError(t, err, "failed to delete workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func updateWorkflow(
	t *testing.T,
	th *testutils.EVMBackendTH,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string,
) {
	t.Helper()
	_, err := wfRegC.UpdateWorkflow(th.ContractsOwner, workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
	require.NoError(t, err, "failed to update workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

type evtHandler interface {
	io.Closer
	Handle(ctx context.Context, event syncer.Event) error
}

type testSecretsWorkEventHandler struct {
	wrappedHandler evtHandler
	registeredCh   chan syncer.Event
}

func (m *testSecretsWorkEventHandler) Close() error { return m.wrappedHandler.Close() }

func (m *testSecretsWorkEventHandler) Handle(ctx context.Context, event syncer.Event) error {
	switch {
	case event.EventType == syncer.ForceUpdateSecretsEvent:
		return m.wrappedHandler.Handle(ctx, event)
	case event.EventType == syncer.WorkflowRegisteredEvent:
		m.registeredCh <- event
		return nil
	default:
		panic(fmt.Sprintf("unexpected event type: %v", event.EventType))
	}
}

func encryptSecrets(t *testing.T, workflowOwner string, secretsMap map[string][]string, encryptionKey workflowkey.Key) []byte {
	sm, secretsEnvVars, err := secrets.EncryptSecretsForNodes(
		workflowOwner,
		secretsMap,
		map[string][32]byte{
			"p2pId": encryptionKey.PublicKey(),
		},
		secrets.SecretsConfig{},
	)
	require.NoError(t, err)

	secretsPayload, err := json.Marshal(secrets.EncryptedSecretsResult{
		EncryptedSecrets: sm,
		Metadata: secrets.Metadata{
			WorkflowOwner:          workflowOwner,
			EnvVarsAssignedToNodes: secretsEnvVars,
			NodePublicEncryptionKeys: map[string]string{
				"p2pId": encryptionKey.PublicKeyString(),
			},
		},
	})
	require.NoError(t, err)
	return secretsPayload
}
