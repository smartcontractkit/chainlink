package syncer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

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

func Test_workflowMetadataToEvents(t *testing.T) {
	t.Run("WorkflowRegisteredEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
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

		events := wr.workflowMetadataToEvents(ctx, metadata, donID)

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
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
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

		events := wr.workflowMetadataToEvents(ctx, metadata, donID)

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
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
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

		events := wr.workflowMetadataToEvents(ctx, metadata, donID)

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
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
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

		events := wr.workflowMetadataToEvents(ctx, metadata, donID)

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
		events = wr.workflowMetadataToEvents(ctx, metadata, donID)
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
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
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

		events := wr.workflowMetadataToEvents(ctx, metadata, donID)
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
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
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

		events := wr.workflowMetadataToEvents(ctx, metadata, donID)

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
}
