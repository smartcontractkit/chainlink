package syncer

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	wfstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/mocks"
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

func (m *mockFetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	return m.responseMap[url].Body, m.responseMap[url].Err
}

func newMockFetcher(m map[string]mockFetchResp) FetcherFunc {
	return (&mockFetcher{responseMap: m}).Fetch
}

func Test_Handler(t *testing.T) {
	lggr := logger.TestLogger(t)
	emitter := custmsg.NewLabeler()
	t.Run("success", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
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

		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)
		mockORM.EXPECT().Update(matches.AnyContext, giveHash, "contents").Return(int64(1), nil)
		h := newEventHandler(lggr, mockORM, fetcher, nil, nil, nil, emitter, nil)
		err = h.Handle(ctx, giveEvent)
		require.NoError(t, err)
	})

	t.Run("fails with unsupported event type", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)

		giveEvent := WorkflowRegistryEvent{}
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}

		h := newEventHandler(lggr, mockORM, fetcher, nil, nil, nil, emitter, nil)
		err := h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "event type unsupported")
	})

	t.Run("fails to get secrets url", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		h := newEventHandler(lggr, mockORM, nil, nil, nil, nil, emitter, nil)
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

		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return nil, assert.AnError
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)
		h := newEventHandler(lggr, mockORM, fetcher, nil, nil, nil, emitter, nil)
		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to update secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
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

		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().GetSecretsURLByHash(matches.AnyContext, giveHash).Return(giveURL, nil)
		mockORM.EXPECT().Update(matches.AnyContext, giveHash, "contents").Return(0, assert.AnError)
		h := newEventHandler(lggr, mockORM, fetcher, nil, nil, nil, emitter, nil)
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
	t.Run("success with paused workflow registered", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary     = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)
			config     = []byte("")
			secretsURL = "http://example.com"
			binaryURL  = "http://example.com/binary"
			configURL  = "http://example.com/config"
			wfOwner    = []byte("0xOwner")

			fetcher = newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: binary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			})
		)

		giveWFID := workflowID(binary, config, []byte(secretsURL))

		b, err := hex.DecodeString(giveWFID)
		require.NoError(t, err)
		wfID := make([]byte, 32)
		copy(wfID, b)

		paused := WorkflowRegistryWorkflowRegisteredV1{
			Status:        uint8(1),
			WorkflowID:    [32]byte(wfID),
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}

		h := &eventHandler{
			lggr:    lggr,
			orm:     orm,
			fetcher: fetcher,
			emitter: emitter,
		}
		err = h.workflowRegisteredEvent(ctx, paused)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusPaused, dbSpec.Status)
	})

	t.Run("success with active workflow registered", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			db      = pgtest.NewSqlxDB(t)
			orm     = NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary     = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)
			config     = []byte("")
			secretsURL = "http://example.com"
			binaryURL  = "http://example.com/binary"
			configURL  = "http://example.com/config"
			wfOwner    = []byte("0xOwner")

			fetcher = newMockFetcher(map[string]mockFetchResp{
				binaryURL:  {Body: binary, Err: nil},
				configURL:  {Body: config, Err: nil},
				secretsURL: {Body: []byte("secrets"), Err: nil},
			})
		)

		giveWFID := workflowID(binary, config, []byte(secretsURL))

		b, err := hex.DecodeString(giveWFID)
		require.NoError(t, err)
		wfID := make([]byte, 32)
		copy(wfID, b)

		active := WorkflowRegistryWorkflowRegisteredV1{
			Status:        uint8(0),
			WorkflowID:    [32]byte(wfID),
			WorkflowOwner: wfOwner,
			WorkflowName:  "workflow-name",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}

		er := newEngineRegistry()
		store := wfstore.NewDBStore(db, lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		h := &eventHandler{
			lggr:           lggr,
			orm:            orm,
			fetcher:        fetcher,
			emitter:        emitter,
			engineRegistry: er,
			capRegistry:    registry,
			workflowStore:  store,
		}
		err = h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, hex.EncodeToString(wfOwner), "workflow-name")
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, "workflow-name", dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, err := h.engineRegistry.Get(giveWFID)
		require.NoError(t, err)
		err = engine.Ready()
		require.NoError(t, err)
	})
}
