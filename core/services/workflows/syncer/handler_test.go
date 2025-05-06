package syncer

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	wfstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
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
	calledMap   map[string]int
}

func (m *mockFetcher) Fetch(_ context.Context, mid string, req ghcapabilities.Request) ([]byte, error) {
	m.calledMap[req.URL]++
	return m.responseMap[req.URL].Body, m.responseMap[req.URL].Err
}

func (m *mockFetcher) Calls(url string) int {
	return m.calledMap[url]
}

func (m *mockFetcher) FetcherFunc() artifacts.FetcherFunc {
	return m.Fetch
}

func newMockFetcher(m map[string]mockFetchResp) *mockFetcher {
	return &mockFetcher{responseMap: m, calledMap: map[string]int{}}
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveURL := "https://original-url.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := Event{
			EventType: ForceUpdateSecretsEvent,
			Data: ForceUpdateSecretsRequestedV1{
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

		h, err := NewEventHandler(lggr, nil, nil, NewEngineRegistry(), emitter, rl, workflowLimits, store)
		require.NoError(t, err)

		err = h.Handle(ctx, giveEvent)
		require.NoError(t, err)
	})

	t.Run("fails with unsupported event type", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveEvent := Event{}
		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte("contents"), nil
		}

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, nil, nil, NewEngineRegistry(), emitter, rl, workflowLimits, store)
		require.NoError(t, err)

		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "event type unsupported")
	})

	t.Run("fails to get secrets url", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, nil, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, nil, nil, NewEngineRegistry(), emitter, rl, workflowLimits, store)
		require.NoError(t, err)

		giveURL := "https://original-url.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := Event{
			EventType: ForceUpdateSecretsEvent,
			Data: ForceUpdateSecretsRequestedV1{
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveURL := "http://example.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := Event{
			EventType: ForceUpdateSecretsEvent,
			Data: ForceUpdateSecretsRequestedV1{
				SecretsURLHash: giveBytes,
			},
		}

		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return nil, assert.AnError
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)

		decrypter := newMockDecrypter()
		store := artifacts.NewStoreWithDecryptSecretsFn(lggr, mockORM, fetcher, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, nil, nil, NewEngineRegistry(), emitter, rl, workflowLimits, store)
		require.NoError(t, err)
		err = h.Handle(ctx, giveEvent)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to update secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		giveURL := "http://example.com"
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)

		giveHash := hex.EncodeToString(giveBytes)

		giveEvent := Event{
			EventType: ForceUpdateSecretsEvent,
			Data: ForceUpdateSecretsRequestedV1{
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

		h, err := NewEventHandler(lggr, nil, nil, NewEngineRegistry(), emitter, rl, workflowLimits, store)
		require.NoError(t, err)

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
	var binary = wasmtest.CreateTestBinary(binaryCmd, true, t)
	var encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
	var workflowName = "workflow-name"

	defaultValidationFn := func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
		err := h.workflowRegisteredEvent(ctx, event)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := s.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, workflowName, dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: workflowName})
		require.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)
	}

	defaultValidationFnWithFetch := func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
		defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

		// Verify that the URLs have been called
		require.Equal(t, 1, fetcher.Calls(event.BinaryURL))
		require.Equal(t, 1, fetcher.Calls(event.ConfigURL))
		require.Equal(t, 1, fetcher.Calls(event.SecretsURL))
	}

	var tt = []testCase{
		{
			Name: "success with active workflow registered",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name types.WorkflowName, config []byte, binary []byte) (services.Service, error) {
				return &mockEngine{}, nil
			},
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: defaultValidationFnWithFetch,
		},
		{
			Name: "correctly generates the workflow name",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name types.WorkflowName, config []byte, binary []byte) (services.Service, error) {
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
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: defaultValidationFnWithFetch,
		},
		{
			Name: "fails to start engine",
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: encodedBinary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			}),
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name types.WorkflowName, config []byte, binary []byte) (services.Service, error) {
				return &mockEngine{StartErr: assert.AnError}, nil
			},
			GiveConfig: config,
			ConfigURL:  configURL,
			SecretsURL: secretsURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.Error(t, err)
				require.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			Name: "succeeds if correct engine already exists",
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
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				me := &mockEngine{}
				var wfIDBytes [32]byte
				copy(wfIDBytes[:], wfID)
				err := h.engineRegistry.Add(EngineRegistryKey{Owner: wfOwner, Name: workflowName}, me, wfIDBytes)
				require.NoError(t, err)
				err = h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)
			},
		},
		{
			Name: "handles incorrect engine already exists",
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
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				me := &mockEngine{}
				oldWfIDBytes := [32]byte{0, 1, 2, 3, 5}
				err := h.engineRegistry.Add(EngineRegistryKey{Owner: wfOwner, Name: workflowName}, me, oldWfIDBytes)
				require.NoError(t, err)
				err = h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)
				engineInRegistry, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: workflowName})
				assert.True(t, ok)
				require.Equal(t, engineInRegistry.WorkflowID.Hex(), wfID)
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
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(1),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)

				// Verify the record is updated in the database
				dbSpec, err := s.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
				require.NoError(t, err)
				require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
				require.Equal(t, workflowName, dbSpec.WorkflowName)
				require.Equal(t, job.WorkflowSpecStatusPaused, dbSpec.Status)

				// Verify there is no running engine
				_, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
				assert.False(t, ok)
			},
		},
		{
			Name: "same wf ID, different status",
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
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
					SecretsURL:    secretsURL,
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				// Create the record in the database
				entry := &job.WorkflowSpec{
					Workflow:      hex.EncodeToString(binary),
					Config:        string(config),
					WorkflowID:    event.WorkflowID.Hex(),
					Status:        job.WorkflowSpecStatusPaused,
					WorkflowOwner: hex.EncodeToString(event.WorkflowOwner),
					WorkflowName:  event.WorkflowName,
					SpecType:      job.WASMFile,
					BinaryURL:     event.BinaryURL,
					ConfigURL:     event.ConfigURL,
				}
				_, err := s.UpsertWorkflowSpec(ctx, entry)
				require.NoError(t, err)

				err = h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)

				// Verify the record is updated in the database
				dbSpec, err := s.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), workflowName)
				require.NoError(t, err)
				require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
				require.Equal(t, workflowName, dbSpec.WorkflowName)

				// This reflects the event status, not what was previously stored in the DB
				require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

				_, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
				assert.True(t, ok)
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
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

				// Verify that the URLs have been called
				require.Equal(t, 1, fetcher.Calls(event.BinaryURL))
				require.Equal(t, 0, fetcher.Calls(event.ConfigURL))
				require.Equal(t, 1, fetcher.Calls(event.SecretsURL))
			},
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
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
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

				// Verify that the URLs have been called
				require.Equal(t, 1, fetcher.Calls(event.BinaryURL))
				require.Equal(t, 1, fetcher.Calls(event.ConfigURL))
				require.Equal(t, 0, fetcher.Calls(event.SecretsURL))
			},
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
					Status:        uint8(0),
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  workflowName,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
				}
			},
		},
		{
			Name:       "skips fetching if same DB entry exists",
			GiveConfig: config,
			ConfigURL:  configURL,
			BinaryURL:  binaryURL,
			GiveBinary: binary,
			WFOwner:    wfOwner,
			fetcher: newMockFetcher(map[string]mockFetchResp{
				binaryURL: {Body: encodedBinary, Err: nil},
				configURL: {Body: config, Err: nil},
			}),
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher) {
				// Create the record in the database
				entry := &job.WorkflowSpec{
					Workflow:      hex.EncodeToString(binary),
					Config:        string(config),
					WorkflowID:    hex.EncodeToString(event.WorkflowID[:]),
					Status:        job.WorkflowSpecStatusActive,
					WorkflowOwner: hex.EncodeToString(event.WorkflowOwner),
					WorkflowName:  event.WorkflowName,
					SpecType:      job.WASMFile,
					BinaryURL:     event.BinaryURL,
					ConfigURL:     event.ConfigURL,
				}
				_, err := s.UpsertWorkflowSpec(ctx, entry)
				require.NoError(t, err)

				defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

				// Verify that the URLs have not been called
				require.Equal(t, 0, fetcher.Calls(event.BinaryURL))
				require.Equal(t, 0, fetcher.Calls(event.ConfigURL))
				require.Equal(t, 0, fetcher.Calls(event.SecretsURL))
			},
			Event: func(wfID []byte) WorkflowRegisteredV1 {
				return WorkflowRegisteredV1{
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
	fetcher         *mockFetcher
	Event           func([]byte) WorkflowRegisteredV1
	validationFn    func(t *testing.T, ctx context.Context, event WorkflowRegisteredV1, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID string, fetcher *mockFetcher)
	engineFactoryFn func(ctx context.Context, wfid string, owner string, name types.WorkflowName, config []byte, binary []byte) (services.Service, error)
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher.FetcherFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, store, registry, NewEngineRegistry(), emitter, rl, workflowLimits, artifactStore, opts...)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, h.Close()) })

		tc.validationFn(t, ctx, event, h, artifactStore, wfOwner, "workflow-name", wfID, fetcher)
	})
}

type mockArtifactStore struct {
	artifactStore              *artifacts.Store
	deleteWorkflowArtifactsErr error
}

func (m *mockArtifactStore) FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryURL, configURL string) ([]byte, []byte, error) {
	return m.artifactStore.FetchWorkflowArtifacts(ctx, workflowID, binaryURL, configURL)
}
func (m *mockArtifactStore) GetWorkflowSpec(ctx context.Context, workflowOwner string, workflowName string) (*job.WorkflowSpec, error) {
	return m.artifactStore.GetWorkflowSpec(ctx, workflowOwner, workflowName)
}
func (m *mockArtifactStore) UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	return m.artifactStore.UpsertWorkflowSpec(ctx, spec)
}
func (m *mockArtifactStore) UpsertWorkflowSpecWithSecrets(ctx context.Context, entry *job.WorkflowSpec, secretsURL, urlHash, secrets string) (int64, error) {
	return m.artifactStore.UpsertWorkflowSpecWithSecrets(ctx, entry, secretsURL, urlHash, secrets)
}
func (m *mockArtifactStore) DeleteWorkflowArtifacts(ctx context.Context, workflowOwner string, workflowName string, workflowID string) error {
	if m.deleteWorkflowArtifactsErr != nil {
		return m.deleteWorkflowArtifactsErr
	}
	return m.artifactStore.DeleteWorkflowArtifacts(ctx, workflowOwner, workflowName, workflowID)
}
func (m *mockArtifactStore) GetSecrets(ctx context.Context, secretsURL string, workflowID [32]byte, workflowOwner []byte) ([]byte, error) {
	return m.artifactStore.GetSecrets(ctx, secretsURL, workflowID, workflowOwner)
}
func (m *mockArtifactStore) SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return m.artifactStore.SecretsFor(ctx, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID)
}
func (m *mockArtifactStore) GetSecretsURLHash(workflowOwner []byte, secretsURL []byte) ([]byte, error) {
	return m.artifactStore.GetSecretsURLHash(workflowOwner, secretsURL)
}
func (m *mockArtifactStore) GetSecretsURLByID(ctx context.Context, id int64) (string, error) {
	return m.artifactStore.GetSecretsURLByID(ctx, id)
}
func (m *mockArtifactStore) ForceUpdateSecrets(ctx context.Context, secretsURLHash []byte, owner []byte) (string, error) {
	return m.artifactStore.ForceUpdateSecrets(ctx, secretsURLHash, owner)
}

func newMockArtifactStore(as *artifacts.Store, deleteWorkflowArtifactsErr error) WorkflowArtifactsStore {
	return &mockArtifactStore{
		artifactStore:              as,
		deleteWorkflowArtifactsErr: deleteWorkflowArtifactsErr,
	}
}

func Test_workflowDeletedHandler(t *testing.T) {
	t.Run("success deleting existing engine and spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, true, t)
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

		active := WorkflowRegisteredV1{
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher.FetcherFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, store, registry, NewEngineRegistry(), emitter, rl, workflowLimits, artifactStore, WithEngineRegistry(er))
		require.NoError(t, err)
		err =
			h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)

		deleteEvent := WorkflowDeletedV1{
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
		_, ok = h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.False(t, ok)
	})
	t.Run("success deleting non-existing workflow spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, true, t)
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)
		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher.FetcherFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, store, registry, NewEngineRegistry(), emitter, rl, workflowLimits, artifactStore, WithEngineRegistry(er))
		require.NoError(t, err)

		deleteEvent := WorkflowDeletedV1{
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
	t.Run("removes from DB before engine registry", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, true, t)
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

			failWith = "mocked fail DB delete"
		)

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "workflow-name", binary, config, secretsURL)

		require.NoError(t, err)

		active := WorkflowRegisteredV1{
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher.FetcherFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)
		mockAS := newMockArtifactStore(artifactStore, errors.New(failWith))

		h, err := NewEventHandler(lggr, store, registry, NewEngineRegistry(), emitter, rl, workflowLimits, mockAS, WithEngineRegistry(er))
		require.NoError(t, err)
		err =
			h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)

		deleteEvent := WorkflowDeletedV1{
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			DonID:         1,
		}
		err = h.workflowDeletedEvent(ctx, deleteEvent)
		require.Error(t, err, failWith)

		// Verify the record is still in the DB
		_, err = orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)

		// Verify the engine is still running
		_, ok = h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.True(t, ok)
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

			binary        = wasmtest.CreateTestBinary(binaryCmd, true, t)
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
		newWFIDs := hex.EncodeToString(updatedWFID[:])

		active := WorkflowRegisteredV1{
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
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200})
		require.NoError(t, err)

		decrypter := newMockDecrypter()
		artifactStore := artifacts.NewStoreWithDecryptSecretsFn(lggr, orm, fetcher.FetcherFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), decrypter.decryptSecrets)

		h, err := NewEventHandler(lggr, store, registry, NewEngineRegistry(), emitter, rl, workflowLimits, artifactStore, WithEngineRegistry(er))
		require.NoError(t, err)

		err = h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)

		// create a paused event
		pauseEvent := WorkflowPausedV1{
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
		_, ok = h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.False(t, ok)

		// create an updated event
		updatedEvent := WorkflowUpdatedV1{
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

		// new engine is started
		engine, ok = h.engineRegistry.Get(EngineRegistryKey{Owner: wfOwner, Name: dbSpec.WorkflowName})
		assert.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)
		// old engine is no longer running
		require.Equal(t, types.WorkflowID(updatedWFID), engine.WorkflowID)
	})
}
