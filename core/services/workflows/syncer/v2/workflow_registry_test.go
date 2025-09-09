package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockService struct{}

func (m *mockService) Start(context.Context) error { return nil }

func (m *mockService) Close() error { return nil }

func (m *mockService) HealthReport() map[string]error { return map[string]error{"svc": nil} }

func (m *mockService) Ready() error { return nil }

func (m *mockService) Name() string { return "svc" }

func Test_generateReconciliationEventsV2(t *testing.T) {
	// Validate that if no engines are on the node in the registry,
	// and we see that the contract has workflow state,
	// that we generate a WorkflowRegisteredEvent
	t.Run("WorkflowRegisteredEvent_whenNoEnginesInRegistry", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := [32]byte{1}
		owner := []byte{}
		createdAt := uint64(1000000)
		status := uint8(0)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		donFamily := "A"
		tag := "tag1"
		attributes := []byte{}
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		// The only event is WorkflowRegisteredEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegistered, events[0].Name)
		expectedRegisteredEvent := WorkflowRegisteredEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			Tag:           tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedRegisteredEvent, events[0].Data)
	})

	t.Run("WorkflowUpdatedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{1}
		wfName := "wf name 1"
		err := er.Add(wfID, &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata gets updated
		wfID2 := [32]byte{2}
		createdAt := uint64(1000000)
		status := uint8(0)
		binaryURL2 := "b2"
		configURL := "c1"
		donFamily := "A"
		tag := "tag1"
		attributes := []byte{}
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL2,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		require.Len(t, events, 2)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		expectedDeletedEvent := WorkflowDeletedEvent{
			WorkflowID: wfID,
		}
		require.Equal(t, expectedDeletedEvent, events[0].Data)
		require.Equal(t, WorkflowRegistered, events[1].Name)
		expectedRegisteredEvent := WorkflowRegisteredEvent{
			WorkflowID:    wfID2,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL2,
			ConfigURL:     configURL,
			Tag:           tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedRegisteredEvent, events[1].Data)
	})

	t.Run("WorkflowDeletedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		err := er.Add(wfID, &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata is empty
		metadata := []WorkflowMetadataView{}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		// The only event is WorkflowDeletedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		expectedDeletedEvent := WorkflowDeletedEvent{
			WorkflowID: wfID,
		}
		require.Equal(t, expectedDeletedEvent, events[0].Data)
	})

	t.Run("No change", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := [32]byte{1}
		owner := []byte{}
		status := uint8(0)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		// The only event is WorkflowRegisteredEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegistered, events[0].Name)
		expectedRegisteredEvent := WorkflowRegisteredEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			Tag:           tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedRegisteredEvent, events[0].Data)

		// Add the workflow to the engine registry as the handler would
		err = er.Add(wfID, &mockService{})
		require.NoError(t, err)

		// Repeated ticks do not make any new events
		events, err = wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("A paused workflow doesn't start a new workflow", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := [32]byte{1}
		owner := []byte{}
		status := uint8(1)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)
		// No events
		require.Empty(t, events)
	})

	t.Run("A paused workflow deletes a running workflow", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(wfID, &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata gets updated
		status := uint8(1)
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		// The only event is WorkflowDeletedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		expectedDeletedEvent := WorkflowDeletedEvent{
			WorkflowID: wfID,
		}
		require.Equal(t, expectedDeletedEvent, events[0].Data)
	})

	t.Run("reconciles with a pending event if it has the same signature", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// The workflow metadata gets updated
		binaryURL := "b1"
		configURL := "c1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		event := WorkflowRegisteredEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			Tag:           tag,
			Attributes:    attributes,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowRegistered, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			event.WorkflowID.Hex(): {
				Event: Event{
					Data: event,
					Name: WorkflowRegistered,
				},
				signature:   signature,
				id:          event.WorkflowID.Hex(),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		// The only event is WorkflowRegisteredEvent
		// Since there's a failing event in the pendingEvents queue, we should expect to see
		// that event returned to us.
		require.Empty(t, pendingEvents)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegistered, events[0].Name)
		require.Equal(t, event, events[0].Data)
		require.Equal(t, retryCount, events[0].retryCount)
		require.Equal(t, nextRetryAt, events[0].nextRetryAt)
	})

	t.Run("a paused workflow clears a pending created event", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// The workflow metadata gets updated
		binaryURL := "b1"
		configURL := "c1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       WorkflowStatusPaused,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}
		// Now let's emit an event with the same signature; this should remove the event
		// from the pending queue.
		event := WorkflowRegisteredEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			Tag:           tag,
			Attributes:    attributes,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowRegistered, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			event.WorkflowID.Hex(): {
				Event: Event{
					Data: event,
					Name: WorkflowRegistered,
				},
				signature:   signature,
				id:          event.WorkflowID.Hex(),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		require.Empty(t, pendingEvents)
		require.Empty(t, events)
	})

	t.Run("delete events are handled before any other events", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{1}
		wfName := "wf name 1"
		err := er.Add(wfID, &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// The workflow gets a new version with updated metadata, which changes the workflow ID
		wfID2 := [32]byte{2}
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)

		// Delete event happens before create event
		require.Equal(t, events[0].Name, WorkflowDeleted)
		require.Equal(t, events[1].Name, WorkflowRegistered)
	})

	t.Run("pending delete events are handled when workflow metadata no longer exists", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		err := er.Add(wfID, &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// A workflow is to be removed, but hits a failure, causing it to stay pending
		event := WorkflowDeletedEvent{
			WorkflowID: wfID,
		}
		pendingEvents := map[string]*reconciliationEvent{
			hex.EncodeToString(wfID[:]): {
				Event: Event{
					Data: event,
					Name: WorkflowDeleted,
				},
				id:          hex.EncodeToString(wfID[:]),
				signature:   fmt.Sprintf("%s-%s-%s", WorkflowDeleted, hex.EncodeToString(wfID[:]), toSpecStatus(WorkflowStatusActive)),
				nextRetryAt: time.Now(),
				retryCount:  5,
			},
		}

		// No workflows in metadata
		metadata := []WorkflowMetadataView{}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		require.Empty(t, pendingEvents)
	})

	t.Run("pending create events are handled when workflow metadata no longer exists", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// A workflow is added, but hits a failure during creation, causing it to stay pending
		binaryURL := "b1"
		configURL := "c1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		event := WorkflowRegisteredEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			Tag:           tag,
			Attributes:    attributes,
		}
		pendingEvents := map[string]*reconciliationEvent{
			hex.EncodeToString(wfID[:]): {
				Event: Event{
					Data: event,
					Name: WorkflowRegistered,
				},
				id:          hex.EncodeToString(wfID[:]),
				signature:   fmt.Sprintf("%s-%s-%s", WorkflowRegistered, hex.EncodeToString(wfID[:]), toSpecStatus(WorkflowStatusActive)),
				nextRetryAt: time.Now(),
				retryCount:  5,
			},
		}

		// The workflow then gets removed
		metadata := []WorkflowMetadataView{}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata)
		require.NoError(t, err)
		require.Empty(t, events)
		require.Empty(t, pendingEvents)
	})
}

func Test_GetAllowlistedRequests(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	workflowDonNotifier := capabilities.NewDonNotifier()
	er := NewEngineRegistry()

	// Mock allowlisted requests
	expectedRequests := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
		{
			RequestDigest:   [32]byte{1, 2, 3},
			Owner:           common.Address{4, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ExpiryTimestamp: 123456789,
		},
		{
			RequestDigest:   [32]byte{7, 8, 9},
			Owner:           common.Address{10, 11, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ExpiryTimestamp: 987654321,
		},
	}

	// Mock contract reader to return expectedRequests
	mockContractReader := &mockContractReader{
		allowlistedRequests: expectedRequests,
	}

	wr, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return mockContractReader, nil
		},
		"",
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		&eventHandler{},
		workflowDonNotifier,
		er,
	)
	require.NoError(t, err)

	// Simulate syncAllowlistedRequests updating the field
	wr.allowListedMu.Lock()
	wr.allowListedRequests = expectedRequests
	wr.allowListedMu.Unlock()

	// Test GetAllowlistedRequests returns the correct data
	got := wr.GetAllowlistedRequests(ctx)
	require.Equal(t, expectedRequests, got)
}

// Mock contract reader implementation
type mockContractReader struct {
	types.ContractReader
	allowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
}

func (m *mockContractReader) GetLatestValueWithHeadData(
	_ context.Context,
	_ string,
	_ primitives.ConfidenceLevel,
	_ interface{},
	result interface{},
) (*types.Head, error) {
	// Simulate returning allowlisted requests
	if res, ok := result.(*struct {
		Requests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
		err      error
	}); ok {
		res.Requests = m.allowlistedRequests
		return &types.Head{Height: "123"}, nil
	}
	return &types.Head{Height: "0"}, nil
}
