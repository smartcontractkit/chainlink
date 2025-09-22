package syncer

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
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

func Test_generateReconciliationEvents(t *testing.T) {
	t.Run("WorkflowRegisteredEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// The only event is WorkflowRegisteredEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegisteredEvent, events[0].EventType)
		expectedRegisteredEvent := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedRegisteredEvent, events[0].Data)
	})

	t.Run("WorkflowUpdatedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyEvent,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata gets updated
		wfID2 := [32]byte{2}
		status := uint8(0)
		binaryURL2 := "b2"
		configURL := "c1"
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL2,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// The only event is WorkflowUpdatedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowUpdatedEvent, events[0].EventType)
		expectedUpdatedEvent := WorkflowUpdatedV1{
			OldWorkflowID: wfID,
			NewWorkflowID: wfID2,
			WorkflowOwner: owner,
			DonID:         donID,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL2,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedUpdatedEvent, events[0].Data)
	})

	t.Run("WorkflowDeletedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyEvent,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata is empty
		metadata := []GetWorkflowMetadata{}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// The only event is WorkflowDeletedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeletedEvent, events[0].EventType)
		expectedUpdatedEvent := WorkflowDeletedV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			WorkflowName:  wfName,
		}
		require.Equal(t, expectedUpdatedEvent, events[0].Data)
	})

	t.Run("No change", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// The only event is WorkflowRegisteredEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegisteredEvent, events[0].EventType)
		expectedRegisteredEvent := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedRegisteredEvent, events[0].Data)

		// Add the workflow to the engine registry as the handler would
		err = er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
		require.NoError(t, err)

		// Repeated ticks do not make any new events
		events, err = wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("A paused workflow doesn't start a new workflow", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)
		// No events
		require.Empty(t, events)
	})

	t.Run("A paused workflow deletes a running workflow", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// The only event is WorkflowDeletedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeletedEvent, events[0].EventType)
		expectedUpdatedEvent := WorkflowDeletedV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			WorkflowName:  wfName,
		}
		require.Equal(t, expectedUpdatedEvent, events[0].Data)
	})

	t.Run("pending delete events are handled when workflow metadata no longer exists", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
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
		event := WorkflowDeletedV1{
			WorkflowID: wfID,
		}
		pendingEvents := map[string]*reconciliationEvent{
			idFor(owner, wfName): {
				Event: Event{
					Data:      event,
					EventType: WorkflowDeletedEvent,
				},
				id:          idFor(owner, wfName),
				signature:   fmt.Sprintf("%s-%s", WorkflowDeletedEvent, hex.EncodeToString(wfID[:])),
				nextRetryAt: time.Now(),
				retryCount:  5,
			},
		}

		// No workflows in metadata
		metadata := []GetWorkflowMetadata{}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeletedEvent, events[0].EventType)
		require.Empty(t, pendingEvents)
	})

	t.Run("pending create events are handled when workflow metadata no longer exists", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
		event := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    "",
		}
		pendingEvents := map[string]*reconciliationEvent{
			idFor(owner, wfName): {
				Event: Event{
					Data:      event,
					EventType: WorkflowRegisteredEvent,
				},
				id:          idFor(owner, wfName),
				signature:   fmt.Sprintf("%s-%s", WorkflowRegisteredEvent, hex.EncodeToString(wfID[:])),
				nextRetryAt: time.Now(),
				retryCount:  5,
			},
		}

		// The workflow then gets removed
		metadata := []GetWorkflowMetadata{}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)
		require.Empty(t, events)
		require.Empty(t, pendingEvents)
	})

	t.Run("delete events are handled before any other events", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{1}
		wfName := "wf name 1"
		err := er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
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

		// The first workflow is delete and a second workflow is added
		wfID2 := [32]byte{2}
		wfName2 := "wf name 2"
		binaryURL := "b1"
		configURL := "c1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				DonID:        donID,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName2,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   "",
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// Delete event happens before create event
		require.Equal(t, WorkflowDeletedEvent, events[0].EventType)
		require.Equal(t, WorkflowRegisteredEvent, events[1].EventType)
	})

	t.Run("reconciles with a pending event if it has the same signature", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		id := idFor(owner, wfName)
		event := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowRegisteredEvent, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			id: {
				Event: Event{
					Data:      event,
					EventType: WorkflowRegisteredEvent,
				},
				signature:   signature,
				id:          idFor(owner, wfName),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		// The only event is WorkflowRegisteredEvent
		// Since there's a failing event in the pendingEvents queue, we should expect to see
		// that event returned to us.
		require.Empty(t, pendingEvents)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegisteredEvent, events[0].EventType)
		require.Equal(t, event, events[0].Data)
		require.Equal(t, retryCount, events[0].retryCount)
		require.Equal(t, nextRetryAt, events[0].nextRetryAt)
	})

	t.Run("removes pending event if different signature", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       WorkflowStatusPaused,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		// Now let's emit an event that changes the signature; this should remove the event
		// from the pending queue. Since the status is paused, an event won't be emitted.
		id := idFor(owner, wfName)
		wfID = [32]byte{2}
		event := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowRegisteredEvent, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			id: {
				Event: Event{
					Data:      event,
					EventType: WorkflowRegisteredEvent,
				},
				signature:   signature,
				id:          idFor(owner, wfName),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		require.Empty(t, pendingEvents)
		require.Empty(t, events)
	})

	t.Run("removes pending event if the workflow ID changed", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
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
				SyncStrategy: SyncStrategyEvent,
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
		secretsURL := "s1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}
		// Now let's emit an event that changes the signature; this should remove the event
		// from the pending queue.
		id := idFor(owner, wfName)
		wfID = [32]byte{2}
		event := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowRegisteredEvent, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			id: {
				Event: Event{
					Data:      event,
					EventType: WorkflowRegisteredEvent,
				},
				signature:   signature,
				id:          idFor(owner, wfName),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, donID)
		require.NoError(t, err)

		require.Empty(t, pendingEvents)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowRegisteredEvent, events[0].EventType)
		wfID = [32]byte{1}
		expectedEvent := WorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedEvent, events[0].Data)
		require.Equal(t, 0, events[0].retryCount)
	})
}
