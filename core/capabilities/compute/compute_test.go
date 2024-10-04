package compute

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	cappkg "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func Test_Compute_Start_AddsToRegistry(t *testing.T) {
	log := logger.TestLogger(t)
	registry := capabilities.NewRegistry(log)

	compute := NewAction(log, registry)
	compute.modules.clock = clockwork.NewFakeClock()

	require.NoError(t, compute.Start(tests.Context(t)))

	cp, err := registry.Get(tests.Context(t), CapabilityIDCompute)
	require.NoError(t, err)
	assert.Equal(t, compute, cp)
}

func Test_Compute_Execute_MissingConfig(t *testing.T) {
	log := logger.TestLogger(t)
	registry := capabilities.NewRegistry(log)

	compute := NewAction(log, registry)
	compute.modules.clock = clockwork.NewFakeClock()

	require.NoError(t, compute.Start(tests.Context(t)))

	binary := wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)

	config, err := values.WrapMap(map[string]any{
		"binary": binary,
	})
	require.NoError(t, err)
	req := cappkg.CapabilityRequest{
		Inputs: values.EmptyMap(),
		Config: config,
		Metadata: cappkg.RequestMetadata{
			ReferenceID: "compute",
		},
	}
	_, err = compute.Execute(tests.Context(t), req)
	assert.ErrorContains(t, err, "invalid request: could not find \"config\" in map")
}

func Test_Compute_Execute_MissingBinary(t *testing.T) {
	log := logger.TestLogger(t)
	registry := capabilities.NewRegistry(log)

	compute := NewAction(log, registry)
	compute.modules.clock = clockwork.NewFakeClock()

	require.NoError(t, compute.Start(tests.Context(t)))

	config, err := values.WrapMap(map[string]any{
		"config": []byte(""),
	})
	require.NoError(t, err)
	req := cappkg.CapabilityRequest{
		Inputs: values.EmptyMap(),
		Config: config,
		Metadata: cappkg.RequestMetadata{
			ReferenceID: "compute",
		},
	}
	_, err = compute.Execute(tests.Context(t), req)
	assert.ErrorContains(t, err, "invalid request: could not find \"binary\" in map")
}

func Test_Compute_Execute(t *testing.T) {
	log := logger.TestLogger(t)
	registry := capabilities.NewRegistry(log)

	compute := NewAction(log, registry)
	compute.modules.clock = clockwork.NewFakeClock()

	require.NoError(t, compute.Start(tests.Context(t)))

	binary := wasmtest.CreateTestBinary(binaryCmd, binaryLocation, true, t)

	config, err := values.WrapMap(map[string]any{
		"config": []byte(""),
		"binary": binary,
	})
	require.NoError(t, err)
	inputs, err := values.WrapMap(map[string]any{
		"arg0": map[string]any{
			"cool_output": "foo",
		},
	})
	require.NoError(t, err)
	req := cappkg.CapabilityRequest{
		Inputs: inputs,
		Config: config,
		Metadata: cappkg.RequestMetadata{
			WorkflowID:  "workflowID",
			ReferenceID: "compute",
		},
	}
	resp, err := compute.Execute(tests.Context(t), req)
	assert.NoError(t, err)
	assert.True(t, resp.Value.Underlying["Value"].(*values.Bool).Underlying)

	inputs, err = values.WrapMap(map[string]any{
		"arg0": map[string]any{
			"cool_output": "baz",
		},
	})
	require.NoError(t, err)
	config, err = values.WrapMap(map[string]any{
		"config": []byte(""),
		"binary": binary,
	})
	require.NoError(t, err)
	req = cappkg.CapabilityRequest{
		Inputs: inputs,
		Config: config,
		Metadata: cappkg.RequestMetadata{
			ReferenceID: "compute",
		},
	}
	resp, err = compute.Execute(tests.Context(t), req)
	assert.NoError(t, err)
	assert.False(t, resp.Value.Underlying["Value"].(*values.Bool).Underlying)
}
