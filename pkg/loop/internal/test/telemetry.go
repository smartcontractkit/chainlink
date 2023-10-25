package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
)

var _ grpc.ClientConnInterface = (*mockClientConn)(nil)

type staticTelemetry struct {
	endpoints map[string]staticEndpoint
}

type staticEndpoint struct {
}

func (s staticEndpoint) SendLog(log []byte) {}

func (s staticTelemetry) GenMonitoringEndpoint(network string, chainID string, contractID string, telemType string) commontypes.MonitoringEndpoint {
	s.endpoints[fmt.Sprintf("%s_%s_%s_%s", contractID, telemType, network, chainID)] = staticEndpoint{}
	return s.endpoints[fmt.Sprintf("%s_%s_%s_%s", contractID, telemType, network, chainID)]
}

type mockClientConn struct{}

func (m mockClientConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	return nil
}

func (m mockClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func Telemetry(t *testing.T) {
	mcc := mockClientConn{}
	lggr, ol := logger.TestObserved(t, zapcore.ErrorLevel)
	c := internal.NewTelemetryClient(&mcc, lggr)

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
			error:         "cannot generate monitoring endpoint, contractID is empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "",
			network:       "",
			chainID:       "",
			shouldError:   true,
			error:         "cannot generate monitoring endpoint, telemetryType is empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "",
			chainID:       "",
			shouldError:   true,
			error:         "cannot generate monitoring endpoint, network is empty",
		},
		{
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			network:       "some-network",
			chainID:       "",
			shouldError:   true,
			error:         "cannot generate monitoring endpoint, chainID is empty",
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
		e := c.GenMonitoringEndpoint(test.network, test.chainID, test.contractID, test.telemetryType)
		if test.shouldError {
			require.Nil(t, e)
			require.Equal(t, 1, ol.Len())
			require.Contains(t, ol.TakeAll()[0].Message, test.error)
		} else {
			require.NotNil(t, e)
			require.Equal(t, 0, ol.Len())
			e.SendLog([]byte("some-data"))
			require.Equal(t, 0, ol.Len())
		}
	}

	st := staticTelemetry{
		endpoints: make(map[string]staticEndpoint),
	}
	s := internal.NewTelemetryServer(st)

	type endpointTest struct {
		relayID       *pb.RelayID
		contractID    string
		telemetryType string

		shouldError     bool
		error           string
		endpointsLength int
	}

	endpointTests := []endpointTest{
		{
			relayID:       nil,
			contractID:    "",
			telemetryType: "",
			shouldError:   true,
			error:         "contractID cannot be empty",
		},
		{
			relayID:       nil,
			contractID:    "some-contractID",
			telemetryType: "",
			shouldError:   true,
			error:         "telemetryType cannot be empty",
		},
		{
			relayID:       nil,
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			shouldError:   true,
			error:         "RelayID cannot be nil",
		},
		{
			relayID: &pb.RelayID{
				Network: "",
				ChainId: "",
			},
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			shouldError:   true,
			error:         "RelayID.Network cannot be empty",
		},
		{
			relayID: &pb.RelayID{
				Network: "some-network",
				ChainId: "",
			},
			contractID:    "some-contractID",
			telemetryType: "some-telemetryType",
			shouldError:   true,
			error:         "RelayID.ChainId cannot be empty",
		},
		{
			relayID: &pb.RelayID{
				Network: "some-network",
				ChainId: "some-chainID",
			},
			contractID:      "some-contractID",
			telemetryType:   "some-telemetryType",
			shouldError:     false,
			endpointsLength: 1,
		},
		{
			relayID: &pb.RelayID{
				Network: "some-network",
				ChainId: "some-chainID",
			},
			contractID:      "some-contractID",
			telemetryType:   "some-telemetryType",
			shouldError:     false,
			endpointsLength: 1,
		},
		{
			relayID: &pb.RelayID{
				Network: "some-network",
				ChainId: "some-other-chainID",
			},
			contractID:      "some-contractID",
			telemetryType:   "some-telemetryType",
			shouldError:     false,
			endpointsLength: 2,
		},
	}

	for _, test := range endpointTests {
		_, err := s.Send(context.Background(), &pb.TelemetryMessage{
			RelayID:       test.relayID,
			ContractID:    test.contractID,
			TelemetryType: test.telemetryType,
			Payload:       nil,
		})
		if test.shouldError {
			require.Error(t, err)
			require.Contains(t, err.Error(), test.error)
		} else {
			require.NoError(t, err)
			require.Len(t, st.endpoints, test.endpointsLength)
		}

	}

}
