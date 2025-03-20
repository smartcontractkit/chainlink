package syncer

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
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

func newMockFetcher(m map[string]mockFetchResp) FetcherFunc {
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

func (m *mockDecrypter) registerMock(data []byte, owner string, output map[string]string, err error) {
	input := string(data) + owner
	m.mocks[input] = decryptSecretsOutput{output: output, err: err}
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
		h := NewEventHandler(lggr, mockORM, fetcher, nil, nil, emitter, clockwork.NewFakeClock(), workflowkey.Key{}, rl, workflowLimits)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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

		h := NewEventHandler(lggr, mockORM, fetcher, nil, nil, emitter, clockwork.NewFakeClock(), workflowkey.Key{}, rl, workflowLimits)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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

		h := NewEventHandler(lggr, mockORM, nil, nil, nil, emitter, clockwork.NewFakeClock(), workflowkey.Key{}, rl, workflowLimits)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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
		h := NewEventHandler(lggr, mockORM, fetcher, nil, nil, emitter, clockwork.NewFakeClock(), workflowkey.Key{}, rl, workflowLimits)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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
		h := NewEventHandler(lggr, mockORM, fetcher, nil, nil, emitter, clockwork.NewFakeClock(), workflowkey.Key{}, rl, workflowLimits)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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

	defaultValidationFn := func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, wfOwner []byte, wfName string, wfID string) {
		err := h.workflowRegisteredEvent(ctx, event)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := h.orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
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
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, wfOwner []byte, wfName string, wfID string) {
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
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, wfOwner []byte, wfName string, wfID string) {
				me := &mockEngine{}
				h.engineRegistry.Add(wfID, me)
				err := h.workflowRegisteredEvent(ctx, event)
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
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, wfOwner []byte, wfName string, wfID string) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)

				// Verify the record is updated in the database
				dbSpec, err := h.orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
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
	fetcher         FetcherFunc
	Event           func([]byte) WorkflowRegistryWorkflowRegisteredV1
	validationFn    func(t *testing.T, ctx context.Context, event WorkflowRegistryWorkflowRegisteredV1, h *eventHandler, wfOwner []byte, wfName string, wfID string)
	engineFactoryFn func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error)
}

func testRunningWorkflow(t *testing.T, tc testCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = NewWorkflowRegistryDS(db, lggr)
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

		store := wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)
		h := NewEventHandler(lggr, orm, fetcher, store, registry, emitter, clockwork.NewFakeClock(),
			workflowkey.Key{}, rl, workflowLimits, opts...)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets

		tc.validationFn(t, ctx, event, h, wfOwner, "workflow-name", wfID)
	})
}

func Test_workflowDeletedHandler(t *testing.T) {
	t.Run("success deleting existing engine and spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = NewWorkflowRegistryDS(db, lggr)
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
		store := wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)
		h := NewEventHandler(
			lggr,
			orm,
			fetcher,
			store,
			registry,
			emitter,
			clockwork.NewFakeClock(),
			workflowkey.Key{},
			rl,
			workflowLimits,
			WithEngineRegistry(er),
		)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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
			orm     = NewWorkflowRegistryDS(db, lggr)
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
		store := wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)
		h := NewEventHandler(
			lggr,
			orm,
			fetcher,
			store,
			registry,
			emitter,
			clockwork.NewFakeClock(),
			workflowkey.Key{},
			rl,
			workflowLimits,
			WithEngineRegistry(er),
		)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets

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
			orm     = NewWorkflowRegistryDS(db, lggr)
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
		store := wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)
		h := NewEventHandler(
			lggr,
			orm,
			fetcher,
			store,
			registry,
			emitter,
			clockwork.NewFakeClock(),
			workflowkey.Key{},
			rl,
			workflowLimits,
			WithEngineRegistry(er),
		)
		decrypter := newMockDecrypter()
		h.decryptSecrets = decrypter.decryptSecrets
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

func Test_Handler_SecretsFor(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	decodedWorkflowName := "decodedName"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	url := "http://example.com"
	hash := hex.EncodeToString([]byte(url))
	secretsPayload, err := generateSecrets(workflowOwner, map[string][]string{"Foo": []string{"Bar"}}, encryptionKey)
	require.NoError(t, err)
	secretsID, err := orm.Create(testutils.Context(t), url, hash, string(secretsPayload))
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			url: mockFetchResp{Err: errors.New("could not fetch")},
		},
	}
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
	require.NoError(t, err)
	h := NewEventHandler(
		lggr,
		orm,
		fetcher.Fetch,
		wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock()),
		capabilities.NewRegistry(lggr),
		custmsg.NewLabeler(),
		clockwork.NewFakeClock(),
		encryptionKey,
		rl,
		workflowLimits,
	)
	expectedSecrets := map[string]string{
		"Foo": "Bar",
	}
	decrypter := newMockDecrypter()
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)
	h.decryptSecrets = decrypter.decryptSecrets
	gotSecrets, err := h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)
}

func Test_Handler_SecretsFor_RefreshesSecrets(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	decodedWorkflowName := "decodedName"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	secretsPayload, err := generateSecrets(workflowOwner, map[string][]string{"Foo": []string{"Bar"}}, encryptionKey)
	require.NoError(t, err)

	url := "http://example.com"
	hash := hex.EncodeToString([]byte(url))

	secretsID, err := orm.Create(testutils.Context(t), url, hash, string(secretsPayload))
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	secretsPayload, err = generateSecrets(workflowOwner, map[string][]string{"Baz": []string{"Bar"}}, encryptionKey)
	require.NoError(t, err)
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			url: mockFetchResp{Body: secretsPayload},
		},
	}
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
	require.NoError(t, err)
	h := NewEventHandler(
		lggr,
		orm,
		fetcher.Fetch,
		wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock()),
		capabilities.NewRegistry(lggr),
		custmsg.NewLabeler(),
		clockwork.NewFakeClock(),
		encryptionKey,
		rl,
		workflowLimits,
	)
	expectedSecrets := map[string]string{
		"Baz": "Bar",
	}
	decrypter := newMockDecrypter()
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)
	h.decryptSecrets = decrypter.decryptSecrets

	gotSecrets, err := h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)
}

func Test_Handler_SecretsFor_RefreshLogic(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	decodedWorkflowName := "decodedName"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	secretsPayload, err := generateSecrets(workflowOwner, map[string][]string{"Foo": []string{"Bar"}}, encryptionKey)
	require.NoError(t, err)

	url := "http://example.com"
	hash := hex.EncodeToString([]byte(url))

	secretsID, err := orm.Create(testutils.Context(t), url, hash, string(secretsPayload))
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			url: {
				Body: secretsPayload,
			},
		},
	}
	clock := clockwork.NewFakeClock()
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(syncerlimiter.Config{Global: 200, PerOwner: 200})
	require.NoError(t, err)
	h := NewEventHandler(
		lggr,
		orm,
		fetcher.Fetch,
		wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock()),
		capabilities.NewRegistry(lggr),
		custmsg.NewLabeler(),
		clock,
		encryptionKey,
		rl,
		workflowLimits,
	)
	expectedSecrets := map[string]string{
		"Foo": "Bar",
	}
	decrypter := newMockDecrypter()
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)
	h.decryptSecrets = decrypter.decryptSecrets

	gotSecrets, err := h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)

	// Now stub out an unparseable response, since we already fetched it recently above, we shouldn't need to refetch
	// SecretsFor should still succeed.
	fetcher.responseMap[url] = mockFetchResp{}

	gotSecrets, err = h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)

	secretsPayload, err = generateSecrets(workflowOwner, map[string][]string{"Baz": []string{"Bar"}}, encryptionKey)
	require.NoError(t, err)
	fetcher.responseMap[url] = mockFetchResp{
		Body: secretsPayload,
	}

	expectedSecrets = map[string]string{
		"Baz": "Bar",
	}
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)

	// Now advance so that we hit the freshness limit
	clock.Advance(48 * time.Hour)

	gotSecrets, err = h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)
	assert.Equal(t, expectedSecrets, gotSecrets)
}

func generateSecrets(workflowOwner string, secretsMap map[string][]string, encryptionKey workflowkey.Key) ([]byte, error) {
	sm, secretsEnvVars, err := secrets.EncryptSecretsForNodes(
		workflowOwner,
		secretsMap,
		map[string][32]byte{
			"p2pId": encryptionKey.PublicKey(),
		},
		secrets.SecretsConfig{},
	)
	if err != nil {
		return nil, err
	}
	return json.Marshal(secrets.EncryptedSecretsResult{
		EncryptedSecrets: sm,
		Metadata: secrets.Metadata{
			WorkflowOwner:          workflowOwner,
			EnvVarsAssignedToNodes: secretsEnvVars,
			NodePublicEncryptionKeys: map[string]string{
				"p2pId": encryptionKey.PublicKeyString(),
			},
		},
	})
}
