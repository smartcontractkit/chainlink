package compute

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	cappkg "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	corecapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	fetchBinaryLocation = "test/fetch/cmd/testmodule.wasm"
	fetchBinaryCmd      = "core/capabilities/compute/test/fetch/cmd"
	validRequestUUID    = "d2fe6db9-beb4-47c9-b2d6-d3065ace111e"
)

var defaultConfig = Config{
	ServiceConfig: webapi.ServiceConfig{
		RateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
	},
}

type testHarness struct {
	registry         *corecapabilities.Registry
	connector        *gcmocks.GatewayConnector
	log              logger.Logger
	config           Config
	connectorHandler *webapi.OutgoingConnectorHandler
	compute          *Compute
}

func setup(t *testing.T, config Config) testHarness {
	log := logger.TestLogger(t)
	registry := capabilities.NewRegistry(log)
	connector := gcmocks.NewGatewayConnector(t)
	idGeneratorFn := func() string { return validRequestUUID }
	connector.EXPECT().GatewayIDs().Return([]string{"gateway1"})
	connectorHandler, err := webapi.NewOutgoingConnectorHandler(connector, config.ServiceConfig, ghcapabilities.MethodComputeAction, log)
	require.NoError(t, err)

	fetchFactory, err := NewOutgoingConnectorFetcherFactory(connectorHandler, idGeneratorFn)
	require.NoError(t, err)
	compute, err := NewAction(config, log, registry, fetchFactory)
	require.NoError(t, err)
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

func TestComputeStartAddsToRegistry(t *testing.T) {
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(tests.Context(t)))

	cp, err := th.registry.Get(tests.Context(t), CapabilityIDCompute)
	require.NoError(t, err)
	assert.Equal(t, th.compute, cp)
}

func TestComputeExecuteMissingConfig(t *testing.T) {
	t.Parallel()
	th := setup(t, defaultConfig)
	require.NoError(t, th.compute.Start(tests.Context(t)))

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, simpleBinaryLocation, true, t)

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

func TestComputeExecuteMissingBinary(t *testing.T) {
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

func TestComputeExecute(t *testing.T) {
	t.Parallel()
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(tests.Context(t)))

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, simpleBinaryLocation, true, t)

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
	require.NoError(t, err)
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
	require.NoError(t, err)
	assert.False(t, resp.Value.Underlying["Value"].(*values.Bool).Underlying)
}

func TestComputeFetch(t *testing.T) {
	t.Parallel()
	workflowID := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID := "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	th := setup(t, defaultConfig)

	th.connector.EXPECT().DonID().Return("don-id")
	th.connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
	th.connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

	msgID := strings.Join([]string{
		workflowExecutionID,
		ghcapabilities.MethodComputeAction,
		validRequestUUID,
	}, "/")

	gatewayResp := gatewayResponse(t, msgID)
	th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		th.connectorHandler.HandleGatewayMessage(context.Background(), "gateway1", gatewayResp)
	}).Once()

	require.NoError(t, th.compute.Start(tests.Context(t)))

	binary := wasmtest.CreateTestBinary(fetchBinaryCmd, fetchBinaryLocation, true, t)

	config, err := values.WrapMap(map[string]any{
		"config": []byte(""),
		"binary": binary,
	})
	require.NoError(t, err)

	req := cappkg.CapabilityRequest{
		Config: config,
		Metadata: cappkg.RequestMetadata{
			WorkflowID:          workflowID,
			WorkflowExecutionID: workflowExecutionID,
			ReferenceID:         "compute",
		},
	}

	headers, err := values.NewMap(map[string]any{})
	require.NoError(t, err)
	expected := cappkg.CapabilityResponse{
		Value: &values.Map{
			Underlying: map[string]values.Value{
				"Value": &values.Map{
					Underlying: map[string]values.Value{
						"Body":           values.NewBytes([]byte("response body")),
						"Headers":        headers,
						"StatusCode":     values.NewInt64(200),
						"ErrorMessage":   values.NewString(""),
						"ExecutionError": values.NewBool(false),
					},
				},
			},
		},
	}

	actual, err := th.compute.Execute(tests.Context(t), req)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func gatewayResponse(t *testing.T, msgID string) *api.Message {
	headers := map[string]string{"Content-Type": "application/json"}
	body := []byte("response body")
	responsePayload, err := json.Marshal(ghcapabilities.Response{
		StatusCode:     200,
		Headers:        headers,
		Body:           body,
		ExecutionError: false,
	})
	require.NoError(t, err)
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msgID,
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   responsePayload,
		},
	}
}
