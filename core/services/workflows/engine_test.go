package workflows

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockCapability struct {
	capabilities.CapabilityInfo
	capabilities.CallbackExecutable
	response  chan capabilities.CapabilityResponse
	transform func(capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error)
}

func newMockCapability(info capabilities.CapabilityInfo, transform func(capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error)) *mockCapability {
	return &mockCapability{
		transform:      transform,
		CapabilityInfo: info,
		response:       make(chan capabilities.CapabilityResponse, 10),
	}
}

func (m *mockCapability) Execute(ctx context.Context, ch chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	cr, err := m.transform(req)
	if err != nil {
		return err
	}

	ch <- cr
	close(ch)
	m.response <- cr
	return nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

type mockTriggerCapability struct {
	capabilities.CapabilityInfo
	triggerEvent capabilities.CapabilityResponse
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

func (m *mockTriggerCapability) RegisterTrigger(ctx context.Context, ch chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	ch <- m.triggerEvent
	return nil
}

func (m *mockTriggerCapability) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	return nil
}

func TestEngineWithHardcodedWorkflow(t *testing.T) {
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			"on_mercury_report",
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
			"v1.0.0",
		),
	}
	require.NoError(t, reg.Add(ctx, trigger))

	consensus := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
			"v3.0.0",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{
				Value: req.Inputs.Underlying["observations"],
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, consensus))

	target1 := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_polygon-testnet-mumbai",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon mumbai testnet",
			"v1.0.0",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			list := req.Inputs.Underlying["report"].(*values.List)
			return capabilities.CapabilityResponse{
				Value: list.Underlying[0],
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, target1))

	target2 := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_ethereum-testnet-sepolia",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting ethereum sepolia testnet",
			"v1.0.0",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			list := req.Inputs.Underlying["report"].(*values.List)
			return capabilities.CapabilityResponse{
				Value: list.Underlying[0],
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, target2))

	lggr := logger.TestLogger(t)
	eng, err := NewEngine(lggr, reg)
	require.NoError(t, err)

	resp, err := values.NewMap(map[string]any{
		"123": decimal.NewFromFloat(1.00),
		"456": decimal.NewFromFloat(1.25),
		"789": decimal.NewFromFloat(1.50),
	})
	require.NoError(t, err)
	cr := capabilities.CapabilityResponse{
		Value: resp,
	}
	trigger.triggerEvent = cr

	err = eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()
	assert.Equal(t, cr, <-target1.response)
	assert.Equal(t, cr, <-target2.response)
}
