package syncer

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	wfstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFetchResp struct {
	Body []byte
	Err  error
}

type mockFetcher struct {
	responseMap map[string]mockFetchResp
}

func (m *mockFetcher) Fetch(_ context.Context, mid string, req ghcapabilities.Request) ([]byte, error) {
	return m.responseMap[req.URL].Body, m.responseMap[req.URL].Err
}

func newMockFetcher(m map[string]mockFetchResp) artifacts.FetcherFunc {
	return (&mockFetcher{responseMap: m}).Fetch
}

type mockEngine struct {
	CloseErr error
	ReadyErr error
	StartErr error
}

func (m *mockEngine) Ready() error {
	return m.ReadyErr
}

func (m *mockEngine) Close() error {
	return m.CloseErr
}

func (m *mockEngine) Start(_ context.Context) error {
	return m.StartErr
}

func (m *mockEngine) HealthReport() map[string]error { return nil }

func (m *mockEngine) Name() string { return "mockEngine" }

var rlConfig = ratelimiter.Config{
	GlobalRPS:      1000.0,
	GlobalBurst:    1000,
	PerSenderRPS:   30.0,
	PerSenderBurst: 30,
}

type decryptSecretsOutput struct {
	output map[string]string
	err    error
}
type mockDecrypter struct {
	mocks map[string]decryptSecretsOutput
}

func (m *mockDecrypter) decryptSecrets(data []byte, owner string) (map[string]string, error) {
	input := string(data) + owner
	mock, exists := m.mocks[input]
	if exists {
		return mock.output, mock.err
	}
	return map[string]string{}, nil
}

func newMockDecrypter() *mockDecrypter {
	return &mockDecrypter{
		mocks: map[string]decryptSecretsOutput{},
	}
}

func Test_Handler(t *testing.T) {
	lggr := logger.TestLogger(t)
	emitter := custmsg.NewLabeler()
	t.Run("success", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveURL := "https://original-url.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: WorkflowRegistryForceUpdateSecretsRequestedV1{
				SecretsURLHash: giveBytes,
			},
		}

		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)
		mockORM.EXPECT().Update(matches.AnyContext, giveHash, "contents").Return(int64(1), nil)

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, nil, nil, emitter, rl, workflowLimits, store)

		err = h.Handle(ctx, giveEvent)
		require.NoError(t, err)
	})

	t.Run("fails with unsupported event type", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveEvent := WorkflowRegistryEvent{}
		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte("contents"), nil
		}

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, nil, nil, emitter, rl, workflowLimits, store)

		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "event type unsupported")
	})

	t.Run("fails to get secrets url", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, nil, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, nil, nil, emitter, rl, workflowLimits, store)

		giveURL := "https://original-url.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: WorkflowRegistryForceUpdateSecretsRequestedV1{
				SecretsURLHash: giveBytes,
			},
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return("", assert.AnError)
		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.ErrorContains(t, err, assert.AnError.Error())
	})

	t.Run("fails to fetch contents", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveURL := "http://example.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: WorkflowRegistryForceUpdateSecretsRequestedV1{
				SecretsURLHash: giveBytes,
			},
		}

		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return nil, assert.AnError
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, nil, nil, emitter, rl, workflowLimits, store)

		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to update secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveURL := "http://example.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: WorkflowRegistryForceUpdateSecretsRequestedV1{
				SecretsURLHash: giveBytes,
			},
		}

		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)
		mockORM.EXPECT().Update(matches.AnyContext, giveHash, "contents").Return(0, assert.AnError)

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, nil, nil, emitter, rl, workflowLimits, store)

		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})
}

const (
	binaryLocation = "test/simple/cmd/testmodule.wasm"
	binaryCmd      = "core/capabilities/compute/test/simple/cmd"
)

func Test_workflowRegisteredHandler(t *testing.T) {
	var binaryURL = "http://example.com/binary"
	var secretsURL = "http://example.com/secrets"
	var configURL = "http://example.com/config"
	var config = []byte("")
	var wfOwner = []byte("0xOwner")
	var binary = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)
	var encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
	var workflowName = "workflow-name"

	defaultValidationFn := func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string) {
		err := h.workflowRegisteredEvent(ctx, event)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := s.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, workflowName, dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, err := h.engineRegistry.Get(wfID)
		require.NoError(t, err)
		err = engine.Ready()
		require.NoError(t, err)
	}

	var tt = []testCase{
		{
			Name: "success with active workflow registered",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error) {
				return &mockEngine{}, nil
			},
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: defaultValidationFn,
		},
		{
			Name: "correctly generates the workflow name",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error) {
				if _, err := hex.DecodeString(name.Hex()); err != nil {
					return nil, fmt.Errorf("invalid workflow name: %w", err)
				}
				want := hex.EncodeToString([]byte(pkgworkflows.HashTruncateName(name.String())))
				if want != name.Hex() {
					return nil, fmt.Errorf("invalid workflow name: doesn't match, got %s, want %s", name.Hex(), want)
				}
				return &mockEngine{}, nil
			},
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: defaultValidationFn,
		},
		{
			Name: "fails to start engine",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error) {
				return &mockEngine{StartErr: assert.AnError}, nil
			},
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID string) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.Error(t, err)
				require.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			Name: "fails if running engine exists",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string) {
				me := &mockEngine{}
				err := h.engineRegistry.Add(wfID, me)
				require.NoError(t, err)
				err = h.workflowRegisteredEvent(ctx, event)
				require.Error(t, err)
				require.ErrorContains(t, err, "workflow is already running")
			},
		},
		{
			Name: "success with paused workflow registered",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(1),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID string) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)

				// Verify the record is updated in the database
				dbSpec, err := s.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
				require.NoError(t, err)
				require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
				require.Equal(t, workflowName, dbSpec.WorkflowName)
				require.Equal(t, job.WorkflowSpecStatusPaused, dbSpec.Status)

				// Verify there is no running engine
				_, err = h.engineRegistry.Get(wfID)
				require.Error(t, err)
			},
		},
		{
			Name:       "skips fetch if config url is missing",
			GiveConfig: make([]byte, 0),
			ConfigURL:  "",
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			validationFn: defaultValidationFn,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					SecretsURL:    secretsURL,
				}
			},
		},
		{
			Name:       "skips fetch if secrets url is missing",
			GiveConfig: config,
			ConfigURL:  configURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL: {Body: encodedBinary, Err: nil},
				configURL: {Body: config, Err: nil},
			}),
			validationFn: defaultValidationFn,
			Event: func(wfID []byte) WorkflowRegistryWorkflowRegisteredV1 {
				return WorkflowRegistryWorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
				}
			},
		},
	}

	for _, tc := range tt {
		testRunningWorkflow(t, tc)
	}
}

type testCase struct {
	Name            string
	SecretsURL      string
	BinaryURL       string
	GiveBinary      []byte
	GiveConfig      []byte
	ConfigURL       string
	WFOwner         []byte
	fetcher         artifacts.FetcherFunc
	Event           func([]byte) WorkflowRegistryWorkflowRegisteredV1
	validationFn    func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string)
	engineFactoryFn func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error)
}

func testRunningWorkflow(t *testing.T, tc testCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary     = tc.GiveBinary
			config     = tc.GiveConfig
			secretsURL = tc.SecretsURL
			wfOwner    = tc.WFOwner

			fetcher = tc.fetcher
		)

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "workflow-name", binary, config, secretsURL)
		require.NoError(t, err)

		wfID := hex.EncodeToString(giveWFID[:])

		event := tc.Event(giveWFID[:])

		er := NewEngineRegistry()
		opts := []func(*eventHandler){
			WithEngineRegistry(er),
		}
		if tc.engineFactoryFn != nil {
			opts = append(opts, WithEngineFactoryFn(tc.engineFactoryFn))
		}

		store := wfstore.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, store, registry, emitter, rl, workflowLimits, artifactStore, opts...)
		t.Cleanup(func() { assert.NoError(t, h.Close()) })

		tc.validationFn(t, ctx, event, h, artifactStore, wfOwner, "workflow-name", wfID)
	})
}

func Test_workflowDeletedHandler(t *testing.T) {
	t.Run("success deleting existing engine and spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)
			encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
			config        = []byte("")
			secretsURL    = "http://example.com"
			binaryURL     = "http://example.com/binary"
			configURL     = "http://example.com/config"
			wfOwner       = []byte("0xOwner")

			fetcher = newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			})
		)

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "workflow-name", binary, config, secretsURL)

		require.NoError(t, err)
		wfIDs := hex.EncodeToString(giveWFID[:])

		active := WorkflowRegistryWorkflowRegisteredV1{
			Status:        uint8(0),
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}

		er := NewEngineRegistry()
		store := wfstore.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, store, registry, emitter, rl, workflowLimits, artifactStore, WithEngineRegistry(er))
		err = h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, err := h.engineRegistry.Get(wfIDs)
		require.NoError(t, err)
		err = engine.Ready()
		require.NoError(t, err)

		deleteEvent := WorkflowRegistryWorkflowDeletedV1{
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			DonID:         1,
		}
		err = h.workflowDeletedEvent(ctx, deleteEvent)
		require.NoError(t, err)

		// Verify the record is deleted in the database
		_, err = orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.Error(t, err)

		// Verify the engine is deleted
		_, err = h.engineRegistry.Get(wfIDs)
		require.Error(t, err)
	})
	t.Run("success deleting non-existing workflow spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)
			encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
			config        = []byte("")
			secretsURL    = "http://example.com"
			binaryURL     = "http://example.com/binary"
			configURL     = "http://example.com/config"
			wfOwner       = []byte("0xOwner")

			fetcher = newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			})
		)

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "workflow-name", binary, config, secretsURL)
		require.NoError(t, err)

		er := NewEngineRegistry()
		store := wfstore.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)
		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, store, registry, emitter, rl, workflowLimits, artifactStore, WithEngineRegistry(er))

		deleteEvent := WorkflowRegistryWorkflowDeletedV1{
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			DonID:         1,
		}
		err = h.workflowDeletedEvent(ctx, deleteEvent)
		require.NoError(t, err)

		// Verify the record is deleted in the database
		_, err = orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.Error(t, err)
	})
}

func Test_workflowPausedActivatedUpdatedHandler(t *testing.T) {
	t.Run("success pausing activating and updating existing engine and spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)
			encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
			config        = []byte("")
			updateConfig  = []byte("updated")
			secretsURL    = "http://example.com"
			binaryURL     = "http://example.com/binary"
			configURL     = "http://example.com/config"
			newConfigURL  = "http://example.com/new-config"
			wfOwner       = []byte("0xOwner")

			fetcher = newMockFetcher(map[string]mockFetchResp{
				binaryURL:    {Body: encodedBinary, Err: nil},
				configURL:    {Body: config, Err: nil},
				newConfigURL: {Body: updateConfig, Err: nil},
				secretsURL:   {Body: []byte("secrets"), Err: nil},
			})
		)

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "workflow-name", binary, config, secretsURL)
		require.NoError(t, err)
		updatedWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "workflow-name", binary, updateConfig, secretsURL)
		require.NoError(t, err)

		require.NoError(t, err)
		wfIDs := hex.EncodeToString(giveWFID[:])

		require.NoError(t, err)
		newWFIDs := hex.EncodeToString(updatedWFID[:])

		active := WorkflowRegistryWorkflowRegisteredV1{
			Status:        uint8(0),
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}

		er := NewEngineRegistry()
		store := wfstore.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h := NewEventHandler(lggr, store, registry, emitter, rl, workflowLimits, artifactStore, WithEngineRegistry(er))

		err = h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, err := h.engineRegistry.Get(wfIDs)
		require.NoError(t, err)
		err = engine.Ready()
		require.NoError(t, err)

		// create a paused event
		pauseEvent := WorkflowRegistryWorkflowPausedV1{
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			DonID:         1,
		}
		err = h.workflowPausedEvent(ctx, pauseEvent)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err = orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusPaused, dbSpec.Status)

		// Verify the engine is removed
		_, err = h.engineRegistry.Get(wfIDs)
		require.Error(t, err)

		// create an activated workflow event
		activatedEvent := WorkflowRegistryWorkflowActivatedV1{
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			DonID:         1,
		}

		err = h.workflowActivatedEvent(ctx, activatedEvent)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err = orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, err = h.engineRegistry.Get(wfIDs)
		require.NoError(t, err)
		err = engine.Ready()
		require.NoError(t, err)

		// create an updated event
		updatedEvent := WorkflowRegistryWorkflowUpdatedV1{
			OldWorkflowID: giveWFID,
			NewWorkflowID: updatedWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			BinaryURL:     binaryURL,
			ConfigURL:     newConfigURL,
			SecretsURL:    secretsURL,
			DonID:         1,
		}
		err = h.workflowUpdatedEvent(ctx, updatedEvent)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err = orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)
		require.Equal(t, newWFIDs, dbSpec.WorkflowID)
		require.Equal(t, newConfigURL, dbSpec.ConfigURL)
		require.Equal(t, string(updateConfig), dbSpec.Config)

		// old engine is no longer running
		_, err = h.engineRegistry.Get(wfIDs)
		require.Error(t, err)

		// new engine is started
		engine, err = h.engineRegistry.Get(newWFIDs)
		require.NoError(t, err)
		err = engine.Ready()
		require.NoError(t, err)
	})
}
