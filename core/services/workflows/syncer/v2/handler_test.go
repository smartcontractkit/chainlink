package v2

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basicaction"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	confworkflowtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow/server"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
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

func (m *mockFetcher) RetrieveURL(ctx context.Context, req *storage_service.DownloadArtifactRequest) (string, error) {
	m.calledMap[req.Id]++
	return string(m.responseMap[req.Id+"-"+req.Type.String()].Body), m.responseMap[req.Id+"-"+req.Type.String()].Err
}

func (m *mockFetcher) Calls(identifier string) int {
	return m.calledMap[identifier]
}

func (m *mockFetcher) FetcherFunc() types.FetcherFunc {
	return m.Fetch
}

func (m *mockFetcher) RetrieverFunc() types.LocationRetrieverFunc {
	return m.RetrieveURL
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

type mockDrainableEngine struct {
	mockEngine
	draining         atomic.Bool
	activeExecutions atomic.Int32
	drainCalls       atomic.Int32
	closeCalls       atomic.Int32
	drainStartedAtNs atomic.Int64
}

func (m *mockDrainableEngine) Drain() bool {
	started := m.draining.CompareAndSwap(false, true)
	m.draining.Store(true)
	m.drainCalls.Add(1)
	m.drainStartedAtNs.CompareAndSwap(0, time.Now().UnixNano())
	return started
}

func (m *mockDrainableEngine) ActiveExecutions() int32 {
	return m.activeExecutions.Load()
}

func (m *mockDrainableEngine) DrainStartedAt() (time.Time, bool) {
	ns := m.drainStartedAtNs.Load()
	if ns == 0 {
		return time.Time{}, false
	}

	return time.Unix(0, ns), true
}

func (m *mockDrainableEngine) Close() error {
	m.closeCalls.Add(1)
	return m.CloseErr
}

// mockEngineFactory returns a standard mock engine factory for tests.
// It sends nil to initDone to signal successful initialization.
func mockEngineFactory(ctx context.Context, wfid string, owner string, name types.WorkflowName, tag string, config []byte, binary []byte, binaryURL string, initDone chan<- error) (services.Service, error) {
	if initDone != nil {
		initDone <- nil
	}
	return &mockEngine{}, nil
}

func Test_Handler(t *testing.T) {
	t.Run("fails with unsupported event type", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		lf := limits.Factory{Logger: lggr}
		emitter := custmsg.NewLabeler()
		wfStore := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))

		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)

		giveEvent := Event{
			Head: Head{
				Hash:      "0x123",
				Height:    "123",
				Timestamp: 1234567890,
			},
		}
		retriever := func(_ context.Context, _ *storage_service.DownloadArtifactRequest) (string, error) {
			return "", nil
		}
		fetcher := func(_ context.Context, _ string, _ ghcapabilities.Request) ([]byte, error) {
			return []byte("contents"), nil
		}

		store, err := artifacts.NewStore(lggr, mockORM, fetcher, retriever, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		featureFlags, err := v2.NewFeatureFlags(lf, nil)
		require.NoError(t, err)

		h, err := NewEventHandler(lggr, wfStore, nil, true, registry, NewEngineRegistry(), emitter, limiters, featureFlags, rl, workflowLimits, store, workflowEncryptionKey, &testDonNotifier{})
		require.NoError(t, err)

		err = h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "event type unsupported")
	})
}

const (
	binaryLocation = "test/simple/cmd/testmodule.wasm"
	binaryCmd      = "core/capabilities/compute/test/simple/cmd"
	noTeeV2Cmd     = "core/services/workflows/test/wasm/v2/cmd/without_tee"
	withTeeV2Cmd   = "core/services/workflows/test/wasm/v2/cmd/with_tee"
)

func Test_workflowRegisteredHandler(t *testing.T) {
	t.Parallel()
	binaryURLFactory := func(wfID string) string {
		return "http://example.com/" + wfID + "/binary"
	}
	configURLFactory := func(wfID string) string {
		return "http://example.com/" + wfID + "/config"
	}
	config := []byte("")
	wfOwner := testutils.NewAddress().Bytes()

	binary := wasmtest.CreateTestBinary(binaryCmd, true, t)
	encodedBinary := []byte(base64.StdEncoding.EncodeToString(binary))
	workflowTag := "workflow-tag"
	signedURLParameter := "?auth=abc123"

	defaultValidationFn := func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, _ *mockFetcher) {
		err := h.workflowRegisteredEvent(ctx, event)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := s.GetWorkflowSpec(ctx, wfID.Hex())
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, wfName, dbSpec.WorkflowName)
		require.Equal(t, workflowTag, dbSpec.WorkflowTag)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(wfID)
		require.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)
	}

	defaultValidationFnWithFetch := func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string) {
		defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

		// Verify that the URLs have been called
		require.Equal(t, 1, fetcher.Calls(binaryURL+signedURLParameter))
		require.Equal(t, 1, fetcher.Calls(configURL+signedURLParameter))
	}

	tt := []testCase{
		{
			Name: "success with active workflow registered",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			engineFactoryFn:  mockEngineFactory,
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				wfIDString := hex.EncodeToString(wfID)
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(wfIDString),
					ConfigURL:     configURLFactory(wfIDString),
				}
			},
			validationFn: defaultValidationFnWithFetch,
		},
		{
			Name: "correctly generates the workflow name",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name types.WorkflowName, tag string, config []byte, binary []byte, binaryURL string, initDone chan<- error) (services.Service, error) {
				if _, err := hex.DecodeString(name.Hex()); err != nil {
					return nil, fmt.Errorf("invalid workflow name: %w", err)
				}
				want := hex.EncodeToString([]byte(pkgworkflows.HashTruncateName(name.String())))
				if want != name.Hex() {
					return nil, fmt.Errorf("invalid workflow name: doesn't match, got %s, want %s", name.Hex(), want)
				}
				if initDone != nil {
					initDone <- nil
				}
				return &mockEngine{}, nil
			},
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
			validationFn: defaultValidationFnWithFetch,
		},
		{
			Name: "fails to start engine",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			engineFactoryFn: func(ctx context.Context, wfid string, owner string, name types.WorkflowName, tag string, config []byte, binary []byte, binaryURL string, initDone chan<- error) (services.Service, error) {
				if initDone != nil {
					initDone <- nil
				}
				return &mockEngine{StartErr: assert.AnError}, nil
			},
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string,
			) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.Error(t, err)
				require.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			Name: "succeeds if correct engine already exists",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string) {
				me := &mockEngine{}
				err := h.engineRegistry.Add(wfID, event.Source, me)
				require.NoError(t, err)
				err = h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)
			},
		},
		{
			Name: "handles incorrect engine already exists",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
			engineFactoryFn: mockEngineFactory,
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string) {
				me := &mockEngine{}
				oldWfIDBytes := [32]byte{0, 1, 2, 3, 5}
				err := h.engineRegistry.Add(oldWfIDBytes, event.Source, me)
				require.NoError(t, err)
				err = h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)
				engineInRegistry, ok := h.engineRegistry.Get(wfID)
				assert.True(t, ok)
				require.Equal(t, engineInRegistry.WorkflowID, wfID)
			},
		},
		{
			Name: "success with paused workflow registered",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusPaused,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string,
			) {
				err := h.workflowRegisteredEvent(ctx, event)
				require.NoError(t, err)

				// Verify the record is updated in the database
				dbSpec, err := s.GetWorkflowSpec(ctx, wfID.Hex())
				require.NoError(t, err)
				require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
				require.Equal(t, wfName, dbSpec.WorkflowName)
				require.Equal(t, job.WorkflowSpecStatusPaused, dbSpec.Status)

				// Verify there is no running engine
				_, ok := h.engineRegistry.Get(wfID)
				assert.False(t, ok)
			},
		},
		{
			Name: "same wf ID, different status",
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
			engineFactoryFn: mockEngineFactory,
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler,
				s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string,
			) {
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
				dbSpec, err := s.GetWorkflowSpec(ctx, wfID.Hex())
				require.NoError(t, err)
				require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
				require.Equal(t, wfName, dbSpec.WorkflowName)

				// This reflects the event status, not what was previously stored in the DB
				require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

				_, ok := h.engineRegistry.Get(wfID)
				assert.True(t, ok)
			},
		},
		{
			Name:             "skips fetch if config url is missing",
			GiveConfig:       make([]byte, 0),
			ConfigURLFactory: func(string) string { return "" },
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
				})
			},
			engineFactoryFn: mockEngineFactory,
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string) {
				defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

				// Verify that the URLs have been called
				require.Equal(t, 1, fetcher.Calls(binaryURL+signedURLParameter))
				require.Equal(t, 0, fetcher.Calls(configURL+signedURLParameter))
			},
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					WorkflowTag:   workflowTag,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
				}
			},
		},
		{
			Name:             "skips fetching if same DB entry exists",
			GiveConfig:       config,
			ConfigURLFactory: configURLFactory,
			BinaryURLFactory: binaryURLFactory,
			GiveBinary:       binary,
			WFOwner:          wfOwner,
			fetcherFactory: func(wfID []byte) *mockFetcher {
				wfIDString := hex.EncodeToString(wfID)
				signedBinaryURL := binaryURLFactory(wfIDString) + signedURLParameter
				signedConfigURL := configURLFactory(wfIDString) + signedURLParameter
				return newMockFetcher(map[string]mockFetchResp{
					wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
					wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
					signedBinaryURL:                      {Body: encodedBinary, Err: nil},
					signedConfigURL:                      {Body: config, Err: nil},
				})
			},
			engineFactoryFn: mockEngineFactory,
			validationFn: func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string) {
				// Create the record in the database
				entry := &job.WorkflowSpec{
					Workflow:      hex.EncodeToString(binary),
					Config:        string(config),
					WorkflowID:    hex.EncodeToString(event.WorkflowID[:]),
					Status:        job.WorkflowSpecStatusActive,
					WorkflowOwner: hex.EncodeToString(event.WorkflowOwner),
					WorkflowName:  event.WorkflowName,
					WorkflowTag:   workflowTag,
					SpecType:      job.WASMFile,
					BinaryURL:     binaryURL,
					ConfigURL:     configURL,
				}
				_, err := s.UpsertWorkflowSpec(ctx, entry)
				require.NoError(t, err)

				defaultValidationFn(t, ctx, event, h, s, wfOwner, wfName, wfID, fetcher)

				// Verify that the URLs have not been called
				require.Equal(t, 0, fetcher.Calls(binaryURL+signedURLParameter))
				require.Equal(t, 0, fetcher.Calls(configURL+signedURLParameter))
			},
			Event: func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent {
				return WorkflowRegisteredEvent{
					Status:        WorkflowStatusActive,
					WorkflowID:    [32]byte(wfID),
					WorkflowOwner: wfOwner,
					WorkflowName:  wfName,
					BinaryURL:     binaryURLFactory(hex.EncodeToString(wfID)),
					ConfigURL:     configURLFactory(hex.EncodeToString(wfID)),
				}
			},
		},
	}

	for _, tc := range tt {
		testRunningWorkflow(t, tc)
	}
}

func Test_workflowRegisteredHandler_confidentialRouting(t *testing.T) {
	payload, err := anypb.New(&basictrigger.Outputs{CoolOutput: "foo"})
	require.NoError(t, err)

	triggerResponse := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			TriggerType: "basic-test-capture",
			ID:          "id",
			Payload:     payload,
		},
	}

	t.Run("confidential workflow module is hooked correctly", func(t *testing.T) {
		var (
			ctx                   = t.Context()
			lggr                  = logger.TestLogger(t)
			lf                    = limits.Factory{Logger: lggr}
			db                    = pgtest.NewSqlxDB(t)
			orm                   = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter               = custmsg.NewLabeler()
			binary                = wasmtest.CreateTestBinary(withTeeV2Cmd, true, t)
			encodedBinary         = []byte(base64.StdEncoding.EncodeToString(binary))
			config                = []byte("")
			workflowName          = testutils.RandomizeName(t.Name())
			workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
		)
		wfOwner := testutils.NewAddress().Bytes()

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")
		require.NoError(t, err)
		wfIDString := hex.EncodeToString(giveWFID[:])

		binaryURL := "http://example.com/" + wfIDString + "/binary"
		configURL := "http://example.com/" + wfIDString + "/config"
		signedURLParameter := "?auth=abc123"
		signedBinaryURL := binaryURL + signedURLParameter
		signedConfigURL := configURL + signedURLParameter

		fetcher := newMockFetcher(map[string]mockFetchResp{
			wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
			wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
			signedBinaryURL:                      {Body: encodedBinary, Err: nil},
			signedConfigURL:                      {Body: config, Err: nil},
		})
		artifactStore, err := artifacts.NewStore(lggr, orm, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		er := NewEngineRegistry()

		wfStore := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		trigger := &fireOnceTrigger{testActionBase{CapabilityInfo: commoncap.MustNewCapabilityInfo("basic-test-trigger@1.0.0", commoncap.CapabilityTypeCombined, "test capture")}, triggerResponse}
		require.NoError(t, registry.Add(ctx, trigger))
		action := &captureAction{
			t:              t,
			testActionBase: testActionBase{CapabilityInfo: commoncap.MustNewCapabilityInfo("basic-test-action@1.0.0", commoncap.CapabilityTypeCombined, "test action")},
		}
		require.NoError(t, registry.Add(ctx, action))

		executeRequest := &sdk.ExecuteRequest{
			Config:          config,
			Request:         &sdk.ExecuteRequest_Trigger{Trigger: &sdk.Trigger{Payload: payload}},
			MaxResponseSize: 100000,
		}

		confidential := &confidentialCap{
			CapabilityInfo: commoncap.MustNewCapabilityInfo("confidential-workflows@1.0.0-alpha", commoncap.CapabilityTypeCombined, "test confidential cap"),
			t:              t,
			expected: &confworkflowtypes.ConfidentialWorkflowRequest{
				Execution: &confworkflowtypes.WorkflowExecution{
					WorkflowId:        wfIDString,
					BinaryHash:        v2.ComputeBinaryHash(binary),
					SdkExecuteRequest: executeRequest,
					Owner:             hex.EncodeToString(wfOwner),
					BinaryUrl:         binaryURL,
					Requirements:      &sdk.Requirements{Tee: &sdk.Tee{Item: &sdk.Tee_AnyRegions{}}},
				},
				BinaryUrl: binaryURL,
			},
		}

		require.NoError(t, registry.Add(ctx, server.NewClientServer(confidential)))
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)
		featureFlags, err := v2.NewFeatureFlags(lf, nil)
		require.NoError(t, err)

		h, err := NewEventHandler(lggr, wfStore, nil, true, registry, er, emitter, limiters, featureFlags, rl, workflowLimits, artifactStore, workflowEncryptionKey, &testDonNotifier{},
			WithEngineRegistry(er),
		)
		require.NoError(t, err)
		servicetest.Run(t, h)

		event := WorkflowRegisteredEvent{
			Status:        WorkflowStatusActive,
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   "workflow-tag",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
		}

		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: hex.EncodeToString(wfOwner), Workflow: wfIDString})
		err = h.workflowRegisteredEvent(ctx, event)
		require.NoError(t, err)

		assert.Eventually(t, confidential.ran.Load, 10*time.Second, time.Millisecond)

		// The workflow run is delegated, and we simulate that delegation.
		// Therefore, no callback is made during this test unless something went wrong.
		assert.False(t, action.ran.Load())
	})

	t.Run("non-confidential workflow module is hooked correctly", func(t *testing.T) {
		var (
			ctx     = t.Context()
			lggr    = logger.TestLogger(t)
			lf      = limits.Factory{Logger: lggr}
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary                = wasmtest.CreateTestBinary(noTeeV2Cmd, true, t)
			encodedBinary         = []byte(base64.StdEncoding.EncodeToString(binary))
			config                = []byte("")
			wfOwner               = testutils.NewAddress().Bytes()
			workflowName          = testutils.RandomizeName(t.Name())
			workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
		)

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")
		require.NoError(t, err)
		wfIDString := hex.EncodeToString(giveWFID[:])

		binaryURL := "http://example.com/" + wfIDString + "/binary"
		configURL := "http://example.com/" + wfIDString + "/config"
		signedURLParameter := "?auth=abc123"
		signedBinaryURL := binaryURL + signedURLParameter
		signedConfigURL := configURL + signedURLParameter

		fetcher := newMockFetcher(map[string]mockFetchResp{
			wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
			wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
			signedBinaryURL:                      {Body: encodedBinary, Err: nil},
			signedConfigURL:                      {Body: config, Err: nil},
		})
		artifactStore, err := artifacts.NewStore(lggr, orm, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		er := NewEngineRegistry()

		wfStore := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		trigger := &fireOnceTrigger{testActionBase{CapabilityInfo: commoncap.MustNewCapabilityInfo("basic-test-trigger@1.0.0", commoncap.CapabilityTypeCombined, "test capture")}, triggerResponse}
		require.NoError(t, registry.Add(ctx, trigger))
		action := &captureAction{
			t:              t,
			testActionBase: testActionBase{CapabilityInfo: commoncap.MustNewCapabilityInfo("basic-test-action@1.0.0", commoncap.CapabilityTypeCombined, "test action")},
		}
		action.shouldRun.Store(true)
		require.NoError(t, registry.Add(ctx, action))

		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)
		featureFlags, err := v2.NewFeatureFlags(lf, nil)
		require.NoError(t, err)

		h, err := NewEventHandler(lggr, wfStore, nil, true, registry, er, emitter, limiters, featureFlags, rl, workflowLimits, artifactStore, workflowEncryptionKey, &testDonNotifier{},
			WithEngineRegistry(er),
		)
		require.NoError(t, err)
		servicetest.Run(t, h)

		event := WorkflowRegisteredEvent{
			Status:        WorkflowStatusActive,
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   "workflow-tag",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
		}

		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: hex.EncodeToString(wfOwner), Workflow: wfIDString})
		err = h.workflowRegisteredEvent(ctx, event)
		require.NoError(t, err)

		assert.Eventually(t, action.ran.Load, 10*time.Second, time.Millisecond)
	})
}

type testCase struct {
	Name             string
	BinaryURLFactory func(string) string
	GiveBinary       []byte
	GiveConfig       []byte
	ConfigURLFactory func(string) string
	WFOwner          []byte
	fetcherFactory   func(wfID []byte) *mockFetcher
	Event            func(wfID []byte, wfName string, wfOwner []byte) WorkflowRegisteredEvent
	validationFn     func(t *testing.T, ctx context.Context, event WorkflowRegisteredEvent, h *eventHandler, s *artifacts.Store, wfOwner []byte, wfName string, wfID types.WorkflowID, fetcher *mockFetcher, binaryURL string, configURL string)
	engineFactoryFn  func(ctx context.Context, wfid string, owner string, name types.WorkflowName, tag string, config []byte, binary []byte, binaryURL string, initDone chan<- error) (services.Service, error)
}

func testRunningWorkflow(t *testing.T, tc testCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			lf      = limits.Factory{Logger: lggr}
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary                = tc.GiveBinary
			config                = tc.GiveConfig
			workflowName          = testutils.RandomizeName(t.Name())
			workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))

			fetcherFactory = tc.fetcherFactory
		)
		wfOwner := testutils.NewAddress().Bytes()

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")
		require.NoError(t, err)

		event := tc.Event(giveWFID[:], workflowName, wfOwner)

		er := NewEngineRegistry()
		opts := []func(*eventHandler){
			WithEngineRegistry(er),
		}
		if tc.engineFactoryFn != nil {
			opts = append(opts, WithEngineFactoryFn(tc.engineFactoryFn))
		}

		store := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)

		fetcher := fetcherFactory(giveWFID[:])
		artifactStore, err := artifacts.NewStore(lggr, orm, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		h, err := NewEventHandler(lggr, store, nil, true, registry, NewEngineRegistry(), emitter, limiters, nil, rl, workflowLimits, artifactStore, workflowEncryptionKey, &testDonNotifier{}, opts...)
		require.NoError(t, err)
		servicetest.Run(t, h)

		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: hex.EncodeToString(wfOwner), Workflow: hex.EncodeToString(giveWFID[:])})
		tc.validationFn(t, ctx, event, h, artifactStore, wfOwner, workflowName, giveWFID, fetcher, tc.BinaryURLFactory(hex.EncodeToString(giveWFID[:])), tc.ConfigURLFactory(hex.EncodeToString(giveWFID[:])))
	})
}

func Test_customerFacingError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		assert.NoError(t, customerFacingError(nil))
	})

	t.Run("ArtifactFetchError returns deterministic customer message", func(t *testing.T) {
		fetchErr := &types.ArtifactFetchError{
			ArtifactType: "binary",
			URL:          "https://storage.example.com/binary.wasm?Expires=123&Signature=nodeSpecificSig",
			Err:          errors.New("connection refused"),
		}
		got := customerFacingError(fetchErr)
		require.Error(t, got)
		assert.Equal(t, "Internal error: failed to fetch workflow binary from storage. Contact support if this persists.", got.Error())
	})

	t.Run("wrapped ArtifactFetchError is still detected", func(t *testing.T) {
		fetchErr := &types.ArtifactFetchError{
			ArtifactType: "config",
			URL:          "https://storage.example.com/config.yaml?Expires=456&Signature=abc",
			Err:          errors.New("timeout"),
		}
		wrapped := fmt.Errorf("createWorkflowSpec: %w", fetchErr)
		got := customerFacingError(wrapped)
		assert.Contains(t, got.Error(), "workflow config")
		assert.NotContains(t, got.Error(), "Expires")
	})

	t.Run("non-ArtifactFetchError passes through unchanged", func(t *testing.T) {
		original := errors.New("some other error")
		assert.Equal(t, original, customerFacingError(original))
	})
}

type mockArtifactStore struct {
	artifactStore              *artifacts.Store
	deleteWorkflowArtifactsErr error
}

func (m *mockArtifactStore) FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryURL, configURL string) ([]byte, []byte, error) {
	return m.artifactStore.FetchWorkflowArtifacts(ctx, workflowID, binaryURL, configURL)
}

func (m *mockArtifactStore) GetWorkflowSpec(ctx context.Context, workflowID string) (*job.WorkflowSpec, error) {
	return m.artifactStore.GetWorkflowSpec(ctx, workflowID)
}

func (m *mockArtifactStore) UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	return m.artifactStore.UpsertWorkflowSpec(ctx, spec)
}

func (m *mockArtifactStore) DeleteWorkflowArtifacts(ctx context.Context, workflowID string) error {
	if m.deleteWorkflowArtifactsErr != nil {
		return m.deleteWorkflowArtifactsErr
	}
	return m.artifactStore.DeleteWorkflowArtifacts(ctx, workflowID)
}

func (m *mockArtifactStore) DeleteWorkflowArtifactsBatch(ctx context.Context, workflowIDs []string) error {
	return m.artifactStore.DeleteWorkflowArtifactsBatch(ctx, workflowIDs)
}

func newMockArtifactStore(as *artifacts.Store, deleteWorkflowArtifactsErr error) WorkflowArtifactsStore {
	return &mockArtifactStore{
		artifactStore:              as,
		deleteWorkflowArtifactsErr: deleteWorkflowArtifactsErr,
	}
}

func Test_workflowDeletedHandler(t *testing.T) {
	t.Parallel()
	t.Run("success deleting existing engine and spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			lf      = limits.Factory{Logger: lggr}
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary        = wasmtest.CreateTestBinary(binaryCmd, true, t)
			encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
			config        = []byte("")
			workflowName  = testutils.RandomizeName(t.Name())

			workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
		)
		wfOwner := testutils.NewAddress().Bytes()

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")
		require.NoError(t, err)
		wfIDString := hex.EncodeToString(giveWFID[:])

		var (
			binaryURL          = "http://example.com/" + wfIDString + "/binary"
			configURL          = "http://example.com/" + wfIDString + "/config"
			signedURLParameter = "?auth=abc123"
			signedBinaryURL    = binaryURL + signedURLParameter
			signedConfigURL    = configURL + signedURLParameter
			fetcher            = newMockFetcher(map[string]mockFetchResp{
				wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
				wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
				signedBinaryURL:                      {Body: encodedBinary, Err: nil},
				signedConfigURL:                      {Body: config, Err: nil},
			})
		)

		require.NoError(t, err)

		active := WorkflowRegisteredEvent{
			Status:        WorkflowStatusActive,
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   "workflow-tag",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
		}

		er := NewEngineRegistry()
		store := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)

		artifactStore, err := artifacts.NewStore(lggr, orm, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		h, err := NewEventHandler(lggr, store, nil, true, registry, NewEngineRegistry(), emitter, limiters, nil, rl, workflowLimits, artifactStore, workflowEncryptionKey, &testDonNotifier{},
			WithEngineRegistry(er),
			WithEngineFactoryFn(mockEngineFactory),
		)
		require.NoError(t, err)
		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: hex.EncodeToString(wfOwner), Workflow: wfIDString})
		err = h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, types.WorkflowID(giveWFID).Hex())
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, workflowName, dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(types.WorkflowID(giveWFID))
		assert.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)

		deleteEvent := WorkflowDeletedEvent{
			WorkflowID: giveWFID,
		}
		err = h.workflowDeletedEvent(ctx, deleteEvent)
		require.NoError(t, err)

		// Verify the record is deleted in the database
		_, err = orm.GetWorkflowSpec(ctx, types.WorkflowID(giveWFID).Hex())
		require.Error(t, err)

		// Verify the engine is deleted
		_, ok = h.engineRegistry.Get(types.WorkflowID(giveWFID))
		assert.False(t, ok)
	})

	t.Run("success deleting non-existing workflow spec", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			lf      = limits.Factory{Logger: lggr}
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary                = wasmtest.CreateTestBinary(binaryCmd, true, t)
			config                = []byte("")
			workflowName          = testutils.RandomizeName(t.Name())
			workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
		)
		wfOwner := testutils.NewAddress().Bytes()

		fetcher := newMockFetcher(map[string]mockFetchResp{})

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")
		require.NoError(t, err)

		er := NewEngineRegistry()
		store := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)
		artifactStore, err := artifacts.NewStore(lggr, orm, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		h, err := NewEventHandler(lggr, store, nil, true, registry, NewEngineRegistry(), emitter, limiters, nil, rl, workflowLimits, artifactStore, workflowEncryptionKey, &testDonNotifier{}, WithEngineRegistry(er))
		require.NoError(t, err)

		deleteEvent := WorkflowDeletedEvent{
			WorkflowID: giveWFID,
		}
		err = h.workflowDeletedEvent(ctx, deleteEvent)
		require.NoError(t, err)

		// Verify the record is deleted in the database
		_, err = orm.GetWorkflowSpec(ctx, types.WorkflowID(giveWFID).Hex())
		require.Error(t, err)
	})

	t.Run("removes from DB before engine registry", func(t *testing.T) {
		var (
			ctx     = testutils.Context(t)
			lggr    = logger.TestLogger(t)
			lf      = limits.Factory{Logger: lggr}
			db      = pgtest.NewSqlxDB(t)
			orm     = artifacts.NewWorkflowRegistryDS(db, lggr)
			emitter = custmsg.NewLabeler()

			binary                = wasmtest.CreateTestBinary(binaryCmd, true, t)
			encodedBinary         = []byte(base64.StdEncoding.EncodeToString(binary))
			config                = []byte("")
			workflowName          = testutils.RandomizeName(t.Name())
			workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))

			failWith = "mocked fail DB delete"
		)
		wfOwner := testutils.NewAddress().Bytes()

		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")

		require.NoError(t, err)
		wfIDString := hex.EncodeToString(giveWFID[:])

		var (
			binaryURL          = "http://example.com/" + wfIDString + "/binary"
			configURL          = "http://example.com/" + wfIDString + "/config"
			signedURLParameter = "?auth=abc123"
			signedBinaryURL    = binaryURL + signedURLParameter
			signedConfigURL    = configURL + signedURLParameter
			fetcher            = newMockFetcher(map[string]mockFetchResp{
				wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
				wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
				signedBinaryURL:                      {Body: encodedBinary, Err: nil},
				signedConfigURL:                      {Body: config, Err: nil},
			})
		)

		active := WorkflowRegisteredEvent{
			Status:        WorkflowStatusActive,
			WorkflowID:    giveWFID,
			WorkflowOwner: wfOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   "workflow-tag",
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
		}

		er := NewEngineRegistry()
		store := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
		registry := capabilities.NewRegistry(lggr)
		registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		rl, err := ratelimiter.NewRateLimiter(rlConfig)
		require.NoError(t, err)
		workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
		require.NoError(t, err)

		artifactStore, err := artifacts.NewStore(lggr, orm, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		mockAS := newMockArtifactStore(artifactStore, errors.New(failWith))

		h, err := NewEventHandler(lggr, store, nil, true, registry, NewEngineRegistry(), emitter, limiters, nil, rl, workflowLimits, mockAS, workflowEncryptionKey, &testDonNotifier{},
			WithEngineRegistry(er),
			WithEngineFactoryFn(mockEngineFactory),
		)
		require.NoError(t, err)
		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: hex.EncodeToString(wfOwner), Workflow: wfIDString})
		err = h.workflowRegisteredEvent(ctx, active)
		require.NoError(t, err)

		// Verify the record is updated in the database
		dbSpec, err := orm.GetWorkflowSpec(ctx, types.WorkflowID(giveWFID).Hex())
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(wfOwner), dbSpec.WorkflowOwner)
		require.Equal(t, workflowName, dbSpec.WorkflowName)
		require.Equal(t, job.WorkflowSpecStatusActive, dbSpec.Status)

		// Verify the engine is started
		engine, ok := h.engineRegistry.Get(types.WorkflowID(giveWFID))
		assert.True(t, ok)
		err = engine.Ready()
		require.NoError(t, err)

		deleteEvent := WorkflowDeletedEvent{
			WorkflowID: giveWFID,
		}
		err = h.workflowDeletedEvent(ctx, deleteEvent)
		require.Error(t, err, failWith)

		// Verify the record is still in the DB
		_, err = orm.GetWorkflowSpec(ctx, types.WorkflowID(giveWFID).Hex())
		require.NoError(t, err)

		// Verify the engine is still running
		_, ok = h.engineRegistry.Get(giveWFID)
		assert.True(t, ok)
	})
}

type stubWorkflowArtifactsStore struct {
	spec        *job.WorkflowSpec
	deleteErr   error
	deleteCalls atomic.Int32
}

func (s *stubWorkflowArtifactsStore) FetchWorkflowArtifacts(context.Context, string, string, string) ([]byte, []byte, error) {
	return nil, nil, nil
}

func (s *stubWorkflowArtifactsStore) GetWorkflowSpec(context.Context, string) (*job.WorkflowSpec, error) {
	if s.spec == nil {
		return nil, errors.New("not found")
	}
	return s.spec, nil
}

func (s *stubWorkflowArtifactsStore) UpsertWorkflowSpec(context.Context, *job.WorkflowSpec) (int64, error) {
	return 1, nil
}

func (s *stubWorkflowArtifactsStore) DeleteWorkflowArtifacts(context.Context, string) error {
	s.deleteCalls.Add(1)
	return s.deleteErr
}

func (s *stubWorkflowArtifactsStore) DeleteWorkflowArtifactsBatch(context.Context, []string) error {
	return nil
}

func Test_workflowDeletedEvent_DrainInProgress(t *testing.T) {
	t.Parallel()

	workflowID := types.WorkflowID{1}
	drainable := &mockDrainableEngine{}
	drainable.activeExecutions.Store(2)
	artifactStore := &stubWorkflowArtifactsStore{}
	registry := NewEngineRegistry()
	require.NoError(t, registry.Add(workflowID, "test-source", drainable))

	h := &eventHandler{
		lggr:                   logger.TestLogger(t),
		engineRegistry:         registry,
		workflowArtifactsStore: artifactStore,
	}

	err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: workflowID})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDrainInProgress)
	assert.Equal(t, int32(1), drainable.drainCalls.Load())
	assert.Equal(t, int32(0), drainable.closeCalls.Load())
	assert.Equal(t, int32(0), artifactStore.deleteCalls.Load())
	_, ok := registry.Get(workflowID)
	assert.True(t, ok)
}

func Test_workflowDeletedEvent_IgnoresErrAlreadyStopped(t *testing.T) {
	t.Parallel()

	workflowID := types.WorkflowID{2}
	drainable := &mockDrainableEngine{}
	drainable.CloseErr = services.ErrAlreadyStopped
	artifactStore := &stubWorkflowArtifactsStore{}
	registry := NewEngineRegistry()
	require.NoError(t, registry.Add(workflowID, "test-source", drainable))

	h := &eventHandler{
		lggr:                   logger.TestLogger(t),
		engineRegistry:         registry,
		workflowArtifactsStore: artifactStore,
	}

	err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: workflowID})
	require.NoError(t, err)
	assert.Equal(t, int32(1), drainable.closeCalls.Load())
	assert.Equal(t, int32(1), artifactStore.deleteCalls.Load())
	_, ok := registry.Get(workflowID)
	assert.False(t, ok)
}

func Test_workflowRegisteredEvent_DrainingEngineNotTreatedAsHealthy(t *testing.T) {
	t.Parallel()

	workflowID := types.WorkflowID{3}
	drainable := &mockDrainableEngine{
		mockEngine: mockEngine{
			CloseErr: assert.AnError,
		},
	}
	require.True(t, drainable.Drain())

	registry := NewEngineRegistry()
	require.NoError(t, registry.Add(workflowID, "test-source", drainable))

	artifactStore := &stubWorkflowArtifactsStore{
		spec: &job.WorkflowSpec{
			WorkflowID: workflowID.Hex(),
			Status:     job.WorkflowSpecStatusActive,
		},
	}
	h := &eventHandler{
		lggr:                   logger.TestLogger(t),
		engineRegistry:         registry,
		workflowArtifactsStore: artifactStore,
		tracer:                 noop.NewTracerProvider().Tracer(""),
	}

	err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:     WorkflowStatusActive,
		WorkflowID: workflowID,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "could not clean up old engine")
	assert.Equal(t, int32(1), drainable.closeCalls.Load())
}

// mockLinkingService implements the LinkingServiceServer interface for testing
type mockLinkingService struct {
	linkingclient.UnimplementedLinkingServiceServer
	orgID string
}

func (m *mockLinkingService) GetOrganizationFromWorkflowOwner(ctx context.Context, req *linkingclient.GetOrganizationFromWorkflowOwnerRequest) (*linkingclient.GetOrganizationFromWorkflowOwnerResponse, error) {
	return &linkingclient.GetOrganizationFromWorkflowOwnerResponse{
		OrganizationId: m.orgID,
	}, nil
}

func Test_Handler_OrganizationID(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	emitter := custmsg.NewLabeler()
	ctx := testutils.Context(t)

	// Set up mock gRPC server for linking service
	mockLinking := &mockLinkingService{orgID: "test-org"}
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "localhost:0")
	require.NoError(t, err)
	s := grpc.NewServer()
	linkingclient.RegisterLinkingServiceServer(s, mockLinking)
	go func() {
		assert.NoError(t, s.Serve(lis))
	}()
	defer s.Stop()
	linkingURL := lis.Addr().String()

	var (
		lggr          = logger.TestLogger(t)
		lf            = limits.Factory{Logger: lggr}
		mockORM       = mocks.NewORM(t)
		binary        = wasmtest.CreateTestBinary(binaryCmd, true, t)
		encodedBinary = []byte(base64.StdEncoding.EncodeToString(binary))
		config        = []byte("")
		workflowName  = testutils.RandomizeName(t.Name())

		workflowEncryptionKey = workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	)
	wfOwner := testutils.NewAddress().Bytes()

	giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, workflowName, binary, config, "")

	require.NoError(t, err)
	wfIDString := hex.EncodeToString(giveWFID[:])

	// Set up artifact fetcher using existing mockFetcher pattern
	signedBinaryURL := "http://example.com/" + wfIDString + "/binary?auth=abc123"
	signedConfigURL := "http://example.com/" + wfIDString + "/config?auth=abc123"

	fetcher := newMockFetcher(map[string]mockFetchResp{
		wfIDString + "-ARTIFACT_TYPE_BINARY": {Body: []byte(signedBinaryURL), Err: nil},
		wfIDString + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(signedConfigURL), Err: nil},
		signedBinaryURL:                      {Body: encodedBinary, Err: nil},
		signedConfigURL:                      {Body: config, Err: nil},
	})

	// Mock ORM responses
	mockORM.EXPECT().GetWorkflowSpec(mock.Anything, types.WorkflowID(giveWFID).Hex()).Return(nil, errors.New("not found"))
	mockORM.EXPECT().UpsertWorkflowSpec(mock.Anything, mock.AnythingOfType("*job.WorkflowSpec")).Return(int64(1), nil)

	// Set up handler
	er := NewEngineRegistry()
	store := store.NewInMemoryStore(lggr, clockwork.NewFakeClock())
	registry := capabilities.NewRegistry(lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
	limiters, err := v2.NewLimiters(lf, nil)
	require.NoError(t, err)
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 200, PerOwner: 200}, lf)
	require.NoError(t, err)

	artifactStore, err := artifacts.NewStore(lggr, mockORM, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
		ArtifactStorageHost: "example.com",
	}))
	require.NoError(t, err)

	// Create gRPC client and orgResolver
	conn, err := grpc.NewClient(linkingURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	linkingClient := linkingclient.NewLinkingServiceClient(conn)
	orgResolverConfig := orgresolver.Config{
		URL:                           linkingURL,
		TLSEnabled:                    false,
		WorkflowRegistryAddress:       "0x1234567890abcdef",
		WorkflowRegistryChainSelector: 1,
	}
	orgResolver, err := orgresolver.NewOrgResolverWithClient(orgResolverConfig, linkingClient, lggr)
	require.NoError(t, err)
	defer orgResolver.Close()

	h, err := NewEventHandler(lggr, store, nil, true, registry, er, emitter, limiters, nil, rl, workflowLimits, artifactStore, workflowEncryptionKey, &testDonNotifier{},
		WithEngineRegistry(er),
		WithEngineFactoryFn(mockEngineFactory),
		WithOrgResolver(orgResolver),
	)
	require.NoError(t, err)

	// Handle workflow registered event
	event := WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    giveWFID,
		WorkflowOwner: wfOwner,
		WorkflowName:  workflowName,
		WorkflowTag:   "workflow-tag",
		BinaryURL:     "http://example.com/" + wfIDString + "/binary",
		ConfigURL:     "http://example.com/" + wfIDString + "/config",
	}
	// Convert to WorkflowActivatedEvent and call through Handle method to test the full flow
	activatedEvent := WorkflowActivatedEvent(event)
	err = h.Handle(ctx, Event{
		Name: WorkflowActivated,
		Data: activatedEvent,
		Head: Head{
			Hash:      "0x123",
			Height:    "123",
			Timestamp: 1234567890,
		},
	})
	require.NoError(t, err)

	// Verify that WorkflowActivated message was emitted with orgID
	allMessages := observer.Messages(t)

	var orgIDFound bool
	for _, msg := range allMessages {
		if msg.Attrs["beholder_entity"] == "workflows.v2.WorkflowActivated" {
			var payload eventsv2.WorkflowActivated
			require.NoError(t, proto.Unmarshal(msg.Body, &payload))

			if payload.Workflow != nil && payload.Workflow.WorkflowKey != nil && payload.Workflow.WorkflowKey.OrganizationID == "test-org" {
				orgIDFound = true
				break
			}
		}
	}
	require.True(t, orgIDFound, "Expected WorkflowActivated message with orgID to be emitted")

	// Test deletion event
	t.Run("WorkflowDeleted event includes org ID in labels", func(t *testing.T) {
		deleteObserver := beholdertest.NewObserver(t)
		deleteEmitter := custmsg.NewLabeler()

		mockDeleteORM := mocks.NewORM(t)
		spec := &job.WorkflowSpec{
			WorkflowID:    hex.EncodeToString(giveWFID[:]),
			WorkflowOwner: hex.EncodeToString(wfOwner),
			WorkflowName:  workflowName,
		}

		mockDeleteORM.EXPECT().GetWorkflowSpec(mock.Anything, types.WorkflowID(giveWFID).Hex()).Return(spec, nil)
		mockDeleteORM.EXPECT().DeleteWorkflowSpec(mock.Anything, types.WorkflowID(giveWFID).Hex()).Return(nil)

		deleteArtifactStore, err := artifacts.NewStore(lggr, mockDeleteORM, fetcher.FetcherFunc(), fetcher.RetrieverFunc(), clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), lf, artifacts.WithConfig(artifacts.StoreConfig{
			ArtifactStorageHost: "example.com",
		}))
		require.NoError(t, err)

		hDelete, err := NewEventHandler(lggr, store, nil, true, registry, er, deleteEmitter, limiters, nil, rl, workflowLimits, deleteArtifactStore, workflowEncryptionKey, &testDonNotifier{},
			WithEngineRegistry(er),
			WithEngineFactoryFn(mockEngineFactory),
			WithOrgResolver(orgResolver),
		)
		require.NoError(t, err)

		err = hDelete.Handle(ctx, Event{
			Name: WorkflowDeleted,
			Data: WorkflowDeletedEvent{WorkflowID: giveWFID},
			Head: Head{
				Hash:      "0x456",
				Height:    "456",
				Timestamp: 1234567890,
			},
		})
		require.NoError(t, err)

		// Verify that WorkflowDeleted message was emitted with orgID
		deleteMessages := deleteObserver.Messages(t)
		var deleteOrgIDFound bool
		for _, msg := range deleteMessages {
			if msg.Attrs["beholder_entity"] == "workflows.v2.WorkflowDeleted" {
				var payload eventsv2.WorkflowDeleted
				require.NoError(t, proto.Unmarshal(msg.Body, &payload))

				if payload.Workflow != nil && payload.Workflow.WorkflowKey != nil && payload.Workflow.WorkflowKey.OrganizationID == "test-org" {
					deleteOrgIDFound = true
					break
				}
			}
		}
		require.True(t, deleteOrgIDFound, "Expected WorkflowDeleted message with orgID to be emitted")
	})
}

type testActionBase struct {
	commoncap.CapabilityInfo
}

var _ commoncap.ExecutableAndTriggerCapability = (*testActionBase)(nil)

func (t *testActionBase) AckEvent(_ context.Context, _ string, _ string, _ string) error { return nil }
func (t *testActionBase) RegisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	panic("not implemented for this test")
}
func (t *testActionBase) UnregisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) error {
	return nil
}

func (t *testActionBase) RegisterToWorkflow(_ context.Context, _ commoncap.RegisterToWorkflowRequest) error {
	return nil
}

func (t *testActionBase) UnregisterFromWorkflow(_ context.Context, _ commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}

func (t *testActionBase) Execute(_ context.Context, _ commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	panic("not implemented for this test")
}

type fireOnceTrigger struct {
	testActionBase
	triggerResponse commoncap.TriggerResponse
}

func (t *fireOnceTrigger) RegisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	ch := make(chan commoncap.TriggerResponse, 1)
	ch <- t.triggerResponse
	return ch, nil
}

type captureAction struct {
	ran       atomic.Bool
	shouldRun atomic.Bool
	t         *testing.T
	testActionBase
}

var _ commoncap.ExecutableAndTriggerCapability = (*captureAction)(nil)

func (t *captureAction) Execute(_ context.Context, _ commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	assert.True(t.t, t.shouldRun.Load(), "Execute was called when it should not have been")
	t.ran.Store(true)
	result, err := anypb.New(&basicaction.Outputs{AdaptedThing: "result"})
	require.NoError(t.t, err)
	return commoncap.CapabilityResponse{Payload: result}, nil
}

type confidentialCap struct {
	commoncap.CapabilityInfo
	t        *testing.T
	ran      atomic.Bool
	expected *confworkflowtypes.ConfidentialWorkflowRequest
}

func (c *confidentialCap) Execute(_ context.Context, _ commoncap.RequestMetadata, input *confworkflowtypes.ConfidentialWorkflowRequest) (*commoncap.ResponseAndMetadata[*confworkflowtypes.ConfidentialWorkflowResponse], caperrors.Error) {
	// execution ID differs on every run.
	assert.NotEmpty(c.t, input.Execution.ExecutionId)
	input.Execution.ExecutionId = ""

	assert.True(c.t, proto.Equal(c.expected, input), "WorkflowExecution mismatch")
	c.ran.Store(true)

	triggerPayload, err := anypb.New(&basictrigger.Config{Name: "test", Number: 0})
	require.NoError(c.t, err)
	execResult := &sdk.ExecutionResult{
		Result: &sdk.ExecutionResult_TriggerSubscriptions{
			TriggerSubscriptions: &sdk.TriggerSubscriptionRequest{
				Subscriptions: []*sdk.TriggerSubscription{
					{
						Id:      "basic-test-capture@1.0.0",
						Payload: triggerPayload,
						Method:  "Trigger",
					},
				},
			},
		},
	}

	return &commoncap.ResponseAndMetadata[*confworkflowtypes.ConfidentialWorkflowResponse]{
		Response: &confworkflowtypes.ConfidentialWorkflowResponse{SdkExecutionResult: execResult},
	}, nil
}

func (c *confidentialCap) ProvidedTees(_ context.Context, _ commoncap.RequestMetadata, _ *emptypb.Empty) (*commoncap.ResponseAndMetadata[*confworkflowtypes.ProvidedTeesResponse], caperrors.Error) {
	return &commoncap.ResponseAndMetadata[*confworkflowtypes.ProvidedTeesResponse]{
		Response: &confworkflowtypes.ProvidedTeesResponse{Tee: []*sdk.TeeTypeAndRegions{{Type: sdk.TeeType_TEE_TYPE_AWS_NITRO, Regions: []string{"us-west-2"}}}},
	}, nil
}

func (c *confidentialCap) Start(_ context.Context) error {
	return nil
}

func (c *confidentialCap) Close() error {
	return nil
}

func (c *confidentialCap) HealthReport() map[string]error {
	return map[string]error{}
}

func (c *confidentialCap) Name() string {
	return "confidential workflow"
}

func (c *confidentialCap) Description() string {
	return "confidential workflow"
}

func (c *confidentialCap) Ready() error {
	return nil
}

func (c *confidentialCap) Initialise(_ context.Context, _ core.StandardCapabilitiesDependencies) error {
	return nil
}

var _ server.ClientCapability = &confidentialCap{}
