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
	corecapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

var defaultConfig = corecapabilities.Config{
	RateLimiter: common.RateLimiterConfig{
		GlobalRPS:      100.0,
		GlobalBurst:    100,
		PerSenderRPS:   100.0,
		PerSenderBurst: 100,
	},
}

type testHarness struct {
	registry         *corecapabilities.Registry
	connector        *gcmocks.GatewayConnector
	log              logger.Logger
	config           corecapabilities.Config
	connectorHandler *corecapabilities.OutgoingConnectorHandler
	compute          *Compute
}

func setup(t *testing.T, config corecapabilities.Config) testHarness {
	log := logger.TestLogger(t)
	registry := capabilities.NewRegistry(log)
	connector := gcmocks.NewGatewayConnector(t)
	connectorHandler, err := corecapabilities.NewOutgoingConnectorHandler(connector, config, ghcapabilities.MethodComputeAction, log)
	require.NoError(t, err)

	compute := NewAction(log, registry, connectorHandler)
	compute.modules.clock = clockwork.NewFakeClock()

	return testHarness{
		registry:         registry,
		connector:        connector,
		log:              log,
		config:           config,
		connectorHandler: connectorHandler,
		compute:          compute,
	}
}

func Test_Compute_Start_AddsToRegistry(t *testing.T) {
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(tests.Context(t)))

	cp, err := th.registry.Get(tests.Context(t), CapabilityIDCompute)
	require.NoError(t, err)
	assert.Equal(t, th.compute, cp)
}

func Test_Compute_Execute_MissingConfig(t *testing.T) {
	th := setup(t, defaultConfig)
	require.NoError(t, th.compute.Start(tests.Context(t)))

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
	_, err = th.compute.Execute(tests.Context(t), req)
	assert.ErrorContains(t, err, "invalid request: could not find \"config\" in map")
}

func Test_Compute_Execute_MissingBinary(t *testing.T) {
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(tests.Context(t)))

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
	_, err = th.compute.Execute(tests.Context(t), req)
	assert.ErrorContains(t, err, "invalid request: could not find \"binary\" in map")
}

func Test_Compute_Execute(t *testing.T) {
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(tests.Context(t)))

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
	resp, err := th.compute.Execute(tests.Context(t), req)
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
	resp, err = th.compute.Execute(tests.Context(t), req)
	assert.NoError(t, err)
	assert.False(t, resp.Value.Underlying["Value"].(*values.Bool).Underlying)
}
