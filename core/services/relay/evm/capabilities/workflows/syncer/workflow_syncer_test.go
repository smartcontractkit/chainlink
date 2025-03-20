package workflow_registry_syncer_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
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
}

func (m *testEvtHandler) Handle(ctx context.Context, event syncer.Event) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.events = append(m.events, event)
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

func newTestEvtHandler() *testEvtHandler {
	return &testEvtHandler{
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

	testEventHandler := newTestEvtHandler()

	// Create the registry
	registry := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.WorkflowEventPollerConfig{
			QueryCount: 20,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.WithTicker(eventPollTicker.C),
	)

	servicetest.Run(t, registry)

	require.Eventually(t, func() bool {
		numEvents := len(testEventHandler.GetEvents())
		return numEvents == numberWorkflows
	}, tests.WaitTimeout(t), time.Second)

	for _, event := range testEventHandler.GetEvents() {
		assert.Equal(t, syncer.WorkflowRegisteredEvent, event.GetEventType())
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
					assert.Equal(t, syncer.WorkflowRegisteredEvent, event.GetEventType())
				case 1:
					assert.Equal(t, syncer.WorkflowActivatedEvent, event.GetEventType())
				case 2:
					assert.Equal(t, syncer.WorkflowPausedEvent, event.GetEventType())
				case 3:
					assert.Equal(t, syncer.WorkflowUpdatedEvent, event.GetEventType())
				case 4:
					assert.Equal(t, syncer.WorkflowDeletedEvent, event.GetEventType())
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

	testEventHandler := newTestEvtHandler()

	// Create the worker
	worker := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.WorkflowEventPollerConfig{
			QueryCount: 20,
		},
		testEventHandler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.WithTicker(make(chan time.Time)),
	)

	servicetest.Run(t, worker)

	require.Eventually(t, func() bool {
		return len(testEventHandler.GetEvents()) == numberWorkflows
	}, tests.WaitTimeout(t), time.Second)

	for _, event := range testEventHandler.GetEvents() {
		assert.Equal(t, syncer.WorkflowRegisteredEvent, event.GetEventType())
	}
}

func Test_SecretsWorker(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		emitter   = custmsg.NewLabeler()
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = syncer.NewWorkflowRegistryDS(db, lggr)

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
			Status:     uint8(1),
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

	wl, err := syncerlimiter.NewWorkflowLimits(wlConfig)
	require.NoError(t, err)

	handler := &testSecretsWorkEventHandler{
		wrappedHandler: syncer.NewEventHandler(lggr, orm, fetcherFn, nil, nil,
			emitter, clockwork.NewFakeClock(), encryptionKey, rl, wl),
		registeredCh: make(chan syncer.Event, 1),
	}

	worker := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.WorkflowEventPollerConfig{QueryCount: 20},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.WithTicker(giveTicker.C),
	)

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

	handler := newTestEvtHandler()

	worker := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.WorkflowEventPollerConfig{QueryCount: 20},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.WithTicker(giveTicker.C),
	)

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
		orm       = syncer.NewWorkflowRegistryDS(db, lggr)

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

	wl, err := syncerlimiter.NewWorkflowLimits(wlConfig)
	require.NoError(t, err)

	handler := syncer.NewEventHandler(lggr, orm, fetcherFn, nil, nil,
		emitter, clockwork.NewFakeClock(), workflowkey.Key{}, rl, wl, syncer.WithEngineRegistry(er))

	worker := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.WorkflowEventPollerConfig{QueryCount: 20},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.WithTicker(giveTicker.C),
	)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	servicetest.Run(t, worker)

	// generate a log event
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		_, err = er.Get("test-wf")
		if err == nil {
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

type mockEngineFactory struct{}

func (m *mockEngineFactory) new(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error) {
	return &mockService{}, nil
}

func Test_RegistrySyncer_WorkflowRegistered_InitiallyActivated(t *testing.T) {
	var (
		ctx       = coretestutils.Context(t)
		lggr      = logger.TestLogger(t)
		emitter   = custmsg.NewLabeler()
		backendTH = testutils.NewEVMBackendTH(t)
		db        = pgtest.NewSqlxDB(t)
		orm       = syncer.NewWorkflowRegistryDS(db, lggr)

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

	mf := &mockEngineFactory{}
	er := syncer.NewEngineRegistry()
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	wl, err := syncerlimiter.NewWorkflowLimits(wlConfig)
	require.NoError(t, err)
	handler := syncer.NewEventHandler(
		lggr,
		orm,
		fetcherFn,
		nil,
		nil,
		emitter,
		clockwork.NewFakeClock(),
		workflowkey.Key{},
		rl,
		wl,
		syncer.WithEngineRegistry(er),
		syncer.WithEngineFactoryFn(mf.new),
	)

	worker := syncer.NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (syncer.ContractReader, error) {
			return backendTH.NewContractReader(ctx, t, bytes)
		},
		wfRegistryAddr.Hex(),
		syncer.WorkflowEventPollerConfig{QueryCount: 20},
		handler,
		&testDonNotifier{
			don: capabilities.DON{
				ID: donID,
			},
			err: nil,
		},
		syncer.WithTicker(giveTicker.C),
	)

	// setup contract state to allow the secrets to be updated
	updateAllowedDONs(t, backendTH, wfRegistryC, []uint32{donID}, true)
	updateAuthorizedAddress(t, backendTH, wfRegistryC, []common.Address{backendTH.ContractsOwner.From}, true)

	servicetest.Run(t, worker)

	// generate a log event
	registerWorkflow(t, backendTH, wfRegistryC, giveWorkflow)

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		_, err := er.Get("test-wf")
		if err != nil {
			return err != nil
		}

		owner := strings.ToLower(backendTH.ContractsOwner.From.Hex()[2:])
		_, err = orm.GetWorkflowSpec(ctx, owner, "test-wf")
		return err == nil
	}, tests.WaitTimeout(t), time.Second)
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
	Handle(ctx context.Context, event syncer.Event) error
}

type testSecretsWorkEventHandler struct {
	wrappedHandler evtHandler
	registeredCh   chan syncer.Event
}

func (m *testSecretsWorkEventHandler) Handle(ctx context.Context, event syncer.Event) error {
	switch {
	case event.GetEventType() == syncer.ForceUpdateSecretsEvent:
		return m.wrappedHandler.Handle(ctx, event)
	case event.GetEventType() == syncer.WorkflowRegisteredEvent:
		m.registeredCh <- event
		return nil
	default:
		panic(fmt.Sprintf("unexpected event type: %v", event.GetEventType()))
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
