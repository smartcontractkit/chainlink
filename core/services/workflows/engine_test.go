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

type mockTriggerCapability struct {
	capabilities.CapabilityInfo
	ch chan<- capabilities.CapabilityResponse
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

func (m *mockTriggerCapability) RegisterTrigger(ctx context.Context, ch chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	m.ch = ch
	return nil
}

func (m *mockTriggerCapability) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	return nil
}

func TestEngineWithHardcodedWorkflow(t *testing.T) {
	ctx := context.Background()
	reg := coreCap.NewRegistry()

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
			"off-chain-reporting",
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

	target := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_polygon_mainnet",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon mainnet",
			"v1.0.0",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {

			list := req.Inputs.Underlying["report"].(*values.List)
			return capabilities.CapabilityResponse{
				Value: list.Underlying[0],
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, target))

	lggr := logger.TestLogger(t)
	eng, err := NewEngine(lggr, reg)
	require.NoError(t, err)

	err = eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	resp, err := values.NewMap(map[string]any{
		"123": decimal.NewFromFloat(1.00),
		"456": decimal.NewFromFloat(1.25),
		"789": decimal.NewFromFloat(1.50),
	})
	require.NoError(t, err)
	cr := capabilities.CapabilityResponse{
		Value: resp,
	}
	trigger.ch <- cr
	assert.Equal(t, cr, <-target.response)
}
