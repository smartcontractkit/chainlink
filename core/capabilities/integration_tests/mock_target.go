package integration_tests

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func mockEthereumTestnetSepoliaTarget(_ *testing.T, reportsSink chan capabilities.CapabilityResponse) capabilities.TargetCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_ethereum-testnet-sepolia@1.0.0",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting ethereum sepolia testnet",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			m := req.Inputs.Underlying["report"].(*values.Map)
			resp := capabilities.CapabilityResponse{
				Value: m,
			}

			reportsSink <- resp

			return resp, nil
		},
	)
}

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

func (m *mockCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	cr, err := m.transform(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan capabilities.CapabilityResponse, 10)

	m.response <- cr
	ch <- cr
	close(ch)
	return ch, nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
