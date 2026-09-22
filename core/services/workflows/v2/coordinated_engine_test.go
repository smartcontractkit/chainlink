package v2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// noopAcknowledger satisfies v2.Acknowledger without doing anything.
type noopAcknowledger struct{}

func (noopAcknowledger) Ack(_ context.Context, _, _, _ string) error { return nil }

func TestNewCoordinatedEngine_RequiresAcknowledger(t *testing.T) {
	t.Parallel()

	cfg := defaultTestConfig(t, nil)
	cfg.TriggerAcknowledger = nil

	_, err := v2.NewCoordinatedEngine(cfg)
	require.EqualError(t, err, "trigger acknowledger not set")
}

func TestNewCoordinatedEngine_Succeeds(t *testing.T) {
	t.Parallel()

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	cfg := defaultTestConfig(t, nil)
	cfg.CapRegistry = capreg
	cfg.TriggerAcknowledger = noopAcknowledger{}

	engine, err := v2.NewCoordinatedEngine(cfg)
	require.NoError(t, err)
	require.NotNil(t, engine)
	require.Equal(t, "WorkflowEngine.WorkflowCoordinatedEngine", engine.Name())

	// Compile-time interface assertions.
	var _ v2.WorkflowEngine = engine
	var _ v2.EventSink = engine
}
