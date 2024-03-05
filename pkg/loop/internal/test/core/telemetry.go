package core

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var Telemetry = staticTelemetry{
	staticTelemetryConfig: staticTelemetryConfig{
		chainID:    "some-chainID",
		contractID: "some-contractID",
		network:    "some-network",
		payload:    []byte("some-data"),
		telemType:  "some-telemetryType",
	},
}

var _ testtypes.TelemetryEvaluator = staticTelemetry{}

var _ grpc.ClientConnInterface = (*mockClientConn)(nil)

type staticTelemetryConfig struct {
	chainID    string
	contractID string
	network    string
	payload    []byte
	telemType  string
}

type staticEndpoint struct {
	staticTelemetry
}

func (s staticEndpoint) SendLog(ctx context.Context, log []byte) error {
	return s.staticTelemetry.Send(ctx, s.network, s.chainID, s.contractID, s.telemType, log)
}

type staticTelemetry struct {
	staticTelemetryConfig
}

func (s staticTelemetry) NewEndpoint(ctx context.Context, network string, chainID string, contractID string, telemType string) (types.TelemetryClientEndpoint, error) {
	if network != s.network {
		return nil, fmt.Errorf("expected network %s but got %s", s.network, network)
	}
	if chainID != s.chainID {
		return nil, fmt.Errorf("expected chainID %s but got %s", s.chainID, chainID)
	}
	if contractID != s.contractID {
		return nil, fmt.Errorf("expected contractID %s but got %s", s.contractID, contractID)
	}
	if telemType != s.telemType {
		return nil, fmt.Errorf("expected telemType %s but got %s", s.telemType, telemType)
	}

	return staticEndpoint{
		staticTelemetry: s,
	}, nil
}

func (s staticTelemetry) Send(ctx context.Context, n string, chid string, conid string, t string, p []byte) error {
	if n != s.network {
		return fmt.Errorf("expected %s but got %s", s.network, n)
	}
	if chid != s.chainID {
		return fmt.Errorf("expected %s but got %s", s.chainID, chid)
	}
	if conid != s.contractID {
		return fmt.Errorf("expected %s but got %s", s.contractID, conid)
	}
	if t != s.telemType {
		return fmt.Errorf("expected %s but got %s", s.telemType, t)
	}
	if !bytes.Equal(p, s.payload) {
		return fmt.Errorf("expected %s but got %s", s.payload, p)
	}
	return nil
}

func (s staticTelemetry) Evaluate(ctx context.Context, other types.TelemetryClient) error {
	endpoint, err := other.NewEndpoint(ctx, s.network, s.chainID, s.contractID, s.telemType)
	if err != nil {
		return fmt.Errorf("failed to instantiate endpoint: %w", err)
	}
	err = endpoint.SendLog(ctx, s.payload)
	if err != nil {
		return fmt.Errorf("failed to send log: %w", err)
	}
	return nil
}

func (s staticTelemetry) Expected() types.TelemetryClient {
	return s
}

type mockClientConn struct{}

func (m mockClientConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	return nil
}

func (m mockClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func TestTelemetry(t *testing.T) {
	tsc := internal.NewTelemetryServiceClient(mockClientConn{})
	c := internal.NewTelemetryClient(tsc)

	type sendTest struct {
		contractID    string
		telemetryType string
		network       string
		chainID       string
		payload       []byte

		shouldError bool
		error       string
	}

	sendTests := []sendTest{
		{
			contractID:    "",
			telemetryType: "",
			network:       "",
			chainID:       "",
			payload:       nil,
			shouldError:   true,
			error:         "contractID cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "",
			network:       "",
			chainID:       "",
			payload:       nil,
			shouldError:   true,
			error:         "telemetryType cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "",
			chainID:       "",
			payload:       nil,
			shouldError:   true,
			error:         "network cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "some-network",
			chainID:       "",
			payload:       nil,
			shouldError:   true,
			error:         "chainId cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "some-network",
			chainID:       "some-chainID",
			payload:       nil,
			shouldError:   true,
			error:         "payload cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "some-network",
			chainID:       "some-chainID",
			payload:       []byte("some-data"),
			shouldError:   false,
		},
	}

	for _, test := range sendTests {
		err := c.Send(context.Background(), test.network, test.chainID, test.contractID, test.telemetryType, test.payload)
		if test.shouldError {
			require.ErrorContains(t, err, test.error)
		} else {
			require.NoError(t, err)
		}
	}

	type genMonitoringEndpointTest struct {
		contractID    string
		telemetryType string
		network       string
		chainID       string

		shouldError bool
		error       string
	}

	genMonitoringEndpointTests := []genMonitoringEndpointTest{
		{
			contractID:    "",
			telemetryType: "",
			network:       "",
			chainID:       "",
			shouldError:   true,
			error:         "contractID cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "",
			network:       "",
			chainID:       "",
			shouldError:   true,
			error:         "telemetryType cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "",
			chainID:       "",
			shouldError:   true,
			error:         "network cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "some-network",
			chainID:       "",
			shouldError:   true,
			error:         "chainId cannot be empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "some-network",
			chainID:       "some-chainID",
			shouldError:   false,
		},
	}

	for _, test := range genMonitoringEndpointTests {
		e, err := c.NewEndpoint(context.Background(), test.network, test.chainID, test.contractID, test.telemetryType)
		if test.shouldError {
			require.Nil(t, e)
			require.ErrorContains(t, err, test.error)
		} else {
			require.NotNil(t, e)
			require.Nil(t, err)
			require.Nil(t, e.SendLog(context.Background(), []byte("some-data")))
		}
	}
}
