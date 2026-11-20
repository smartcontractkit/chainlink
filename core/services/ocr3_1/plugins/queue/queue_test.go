package queue

import (
	"context"
	"iter"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// libocr calls NewReportingPlugin with no recover, so a panic here takes down
// the whole node. This guards the constructor and the plugin lifecycle.
func TestConsensusQueuePluginFactory_NewReportingPlugin(t *testing.T) {
	t.Parallel()

	cfg := &PluginConfig{Limits: ocr3_1types.ReportingPluginLimits{MaxKeyValueModifiedKeys: 1}}
	factory, err := NewConsensusQueuePluginFactory(logger.Test(t), cfg, limits.Factory{}, noopCoordinator{}, noopQueue{})
	require.NoError(t, err)

	plugin, info, err := factory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{N: 4, F: 1}, nil)
	require.NoError(t, err)
	require.NotNil(t, plugin)

	info1, ok := info.(ocr3_1types.ReportingPluginInfo1)
	require.True(t, ok, "unexpected ReportingPluginInfo type %T", info)
	require.Equal(t, "ConsensusQueueReportingPlugin", info1.Name)
	require.Equal(t, cfg.Limits, info1.Limits)

	require.NoError(t, plugin.Close())
}

type noopCoordinator struct{}

func (n noopCoordinator) Start(ctx context.Context) error {
	return nil
}

func (n noopCoordinator) Close() error {
	return nil
}

func (n noopCoordinator) Ready() error {
	return nil
}

func (n noopCoordinator) HealthReport() map[string]error {
	return map[string]error{}
}

func (n noopCoordinator) Name() string {
	return "NoopTriggerCoordinator"
}

func (n noopCoordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	return nil
}

func (n noopCoordinator) RegisterTriggers(ctx context.Context, subscriber triggers.Subscriber, params triggers.RegistrationParams) ([]string, error) {
	return nil, nil
}

func (n noopCoordinator) UnregisterTriggers(ctx context.Context, workflowID string) error {
	return nil
}

func (n noopCoordinator) AvailableCapacity(ctx context.Context) (triggers.Capacity, error) {
	return triggers.Capacity{}, nil
}

func (n noopCoordinator) StartExecution(ctx context.Context, event triggers.CoordinatedEvent) error {
	return nil
}

type noopQueue struct{}

func (n noopQueue) Put(ctx context.Context, event triggers.CoordinatedEvent) error {
	return nil
}

func (n noopQueue) Drain(ctx context.Context) ([]triggers.CoordinatedEvent, error) {
	return nil, nil
}

func (n noopQueue) Stats(ctx context.Context) (triggers.Stats, error) {
	return triggers.Stats{}, nil
}

func (n noopQueue) Sorted() iter.Seq[triggers.CoordinatedEvent] {
	return func(yield func(triggers.CoordinatedEvent) bool) {}
}

func (n noopQueue) Remove(ctx context.Context, triggerEventIDs iter.Seq[string]) {}
