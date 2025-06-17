package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cappkg "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metering"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

const (
	fetchBinaryCmd   = "core/capabilities/compute/test/fetch/cmd"
	validRequestUUID = "d2fe6db9-beb4-47c9-b2d6-d3065ace111e"
)

var defaultConfig = Config{
	ServiceConfig: webapi.ServiceConfig{
		OutgoingRateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
		RateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
	},
}

type testHarness struct {
	registry         *capabilities.Registry
	connector        *gcmocks.GatewayConnector
	log              logger.Logger
	config           Config
	connectorHandler *webapi.OutgoingConnectorHandler
	compute          *Compute
}

func setup(t *testing.T, config Config) testHarness {
	log := logger.Test(t)
	registry := capabilities.NewRegistry(log)
	connector := gcmocks.NewGatewayConnector(t)
	idGeneratorFn := func() string { return validRequestUUID }
	connectorHandler, err := webapi.NewOutgoingConnectorHandler(connector, config.ServiceConfig, ghcapabilities.MethodComputeAction, log, webapi.WithFixedStart())
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

	require.NoError(t, th.compute.Start(t.Context()))

	cp, err := th.registry.Get(t.Context(), CapabilityIDCompute)
	require.NoError(t, err)
	assert.Equal(t, th.compute, cp)
}

func TestComputeExecuteMissingConfig(t *testing.T) {
	t.Parallel()
	th := setup(t, defaultConfig)
	require.NoError(t, th.compute.Start(t.Context()))

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, true, t)

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
	_, err = th.compute.Execute(t.Context(), req)
	assert.ErrorContains(t, err, "invalid request: could not find \"config\" in map")
}

func TestComputeExecuteMissingBinary(t *testing.T) {
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(t.Context()))

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
	_, err = th.compute.Execute(t.Context(), req)
	assert.ErrorContains(t, err, "invalid request: could not find \"binary\" in map")
}

func TestComputeExecute(t *testing.T) {
	t.Parallel()
	th := setup(t, defaultConfig)

	require.NoError(t, th.compute.Start(t.Context()))

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, true, t)

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
	resp, err := th.compute.Execute(t.Context(), req)
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
	resp, err = th.compute.Execute(t.Context(), req)
	require.NoError(t, err)
	assert.False(t, resp.Value.Underlying["Value"].(*values.Bool).Underlying)
	assert.Equal(t, metering.ComputeUnit.Name, resp.Metadata.Metering[0].SpendUnit)
	assert.Equal(t, "0", resp.Metadata.Metering[0].SpendValue)
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

	gatewayResp := gatewayResponse(t, msgID, []byte("response body"))
	th.connector.EXPECT().
		SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).
		Return(nil).
		Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			th.connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).
		Once()

	require.NoError(t, th.compute.Start(t.Context()))

	binary := wasmtest.CreateTestBinary(fetchBinaryCmd, true, t)

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

	headers, err := values.NewMap(map[string]any{
		"Content-Type": "application/json",
	})
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
		Metadata: cappkg.ResponseMetadata{
			Metering: []cappkg.MeteringNodeDetail{
				{
					Peer2PeerID: "",
					SpendUnit:   metering.ComputeUnit.Name,
					SpendValue:  "0",
				},
			},
		},
	}

	actual, err := th.compute.Execute(t.Context(), req)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestCompute_SpendValueRelativeToComputeTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		time               time.Duration
		expectedSpendValue string
	}{
		{time.Duration(0), "0"},
		{time.Second, "1"},
		{2 * time.Second, "2"},
		{2500 * time.Millisecond, "3"},
		{3 * time.Second, "3"},
	}

	workflowID := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID := "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	msgID := strings.Join([]string{
		workflowExecutionID,
		ghcapabilities.MethodComputeAction,
		validRequestUUID,
	}, "/")
	gatewayResp := gatewayResponse(t, msgID, []byte("response body"))
	binary := wasmtest.CreateTestBinary(fetchBinaryCmd, true, t)

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

	for idx := range tests {
		test := tests[idx]

		t.Run(test.time.String(), func(t *testing.T) {
			t.Parallel()

			th := setup(t, defaultConfig)

			th.connector.EXPECT().DonID().Return("don-id")
			th.connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
			th.connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

			th.connector.EXPECT().
				SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).
				Return(nil).
				Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
					th.connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
				}).
				Once().
				After(test.time)

			require.NoError(t, th.compute.Start(t.Context()))

			response, err := th.compute.Execute(t.Context(), req)
			require.NoError(t, err)

			require.Len(t, response.Metadata.Metering, 1)
			assert.Equal(t, test.expectedSpendValue, response.Metadata.Metering[0].SpendValue)
		})
	}
}

func TestComputeFetchMaxResponseSizeBytes(t *testing.T) {
	t.Parallel()
	workflowID := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID := "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"

	th := setup(t, Config{
		ServiceConfig: webapi.ServiceConfig{
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		},
		MaxResponseSizeBytes: 1 * 1024,
	})

	th.connector.EXPECT().DonID().Return("don-id")
	th.connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
	th.connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

	msgID := strings.Join([]string{
		workflowExecutionID,
		ghcapabilities.MethodComputeAction,
		validRequestUUID,
	}, "/")

	gatewayResp := gatewayResponse(t, msgID, make([]byte, 2*1024))
	th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		th.connectorHandler.HandleGatewayMessage(context.Background(), "gateway1", gatewayResp)
	}).Once()

	require.NoError(t, th.compute.Start(t.Context()))

	binary := wasmtest.CreateTestBinary(fetchBinaryCmd, true, t)

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

	_, err = th.compute.Execute(t.Context(), req)
	require.ErrorContains(t, err, fmt.Sprintf("response size %d exceeds maximum allowed size %d", 2092, 1*1024))
}

func gatewayResponse(t *testing.T, msgID string, body []byte) *api.Message {
	headers := map[string]string{"Content-Type": "application/json"}
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
