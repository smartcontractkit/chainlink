package telemetry

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	mocks3 "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
)

func setupMockConfig(t *testing.T, useBatchSend bool) *mocks.TelemetryIngress {
	tic := mocks.NewTelemetryIngress(t)
	tic.On("BufferSize").Return(uint(123))
	tic.On("Logging").Return(true)
	tic.On("MaxBatchSize").Return(uint(51))
	tic.On("SendInterval").Return(time.Millisecond * 512)
	tic.On("SendTimeout").Return(time.Second * 7)
	tic.On("UniConn").Return(true)
	tic.On("UseBatchSend").Return(useBatchSend)

	return tic
}

func TestManagerAgents(t *testing.T) {
	tic := setupMockConfig(t, true)
	te := mocks.NewTelemetryIngressEndpoint(t)
	te.On("Network").Return("network-1")
	te.On("ChainID").Return("network-1-chainID-1")
	te.On("ServerPubKey").Return("some-pubkey")
	u, _ := url.Parse("http://some-url.test")
	te.On("URL").Return(u)
	tic.On("Endpoints").Return([]config.TelemetryIngressEndpoint{te})

	lggr, _ := logger.TestLoggerObserved(t, zapcore.InfoLevel)

	ks := mocks3.NewCSA(t)

	tm := NewManager(tic, ks, lggr)
	require.Equal(t, "*synchronization.telemetryIngressBatchClient", reflect.TypeOf(tm.endpoints[0].client).String())
	me := tm.GenMonitoringEndpoint("network-1", "network-1-chainID-1", "", "")
	require.Equal(t, "*telemetry.IngressAgentBatch", reflect.TypeOf(me).String())

	tic = setupMockConfig(t, false)
	tic.On("Endpoints").Return([]config.TelemetryIngressEndpoint{te})
	tm = NewManager(tic, ks, lggr)
	require.Equal(t, "*synchronization.telemetryIngressClient", reflect.TypeOf(tm.endpoints[0].client).String())
	me = tm.GenMonitoringEndpoint("network-1", "network-1-chainID-1", "", "")
	require.Equal(t, "*telemetry.IngressAgent", reflect.TypeOf(me).String())
}

func TestNewManager(t *testing.T) {
	type endpointTest struct {
		network       string
		chainID       string
		url           string
		pubKey        string
		shouldError   bool
		expectedError string
	}

	endpoints := []endpointTest{
		{
			network:     "NETWORK-1",
			chainID:     "NETWORK-1-CHAINID-1",
			url:         "http://network-1-chainID-1.test",
			pubKey:      GetDummyKeyStringWithPrefix("network-1-chainID-1-pub-key"),
			shouldError: false,
		},
		{
			network:     "NETWORK-1",
			chainID:     "NETWORK-1-CHAINID-2",
			url:         "http://network-1-chainID-2.test",
			pubKey:      GetDummyKeyStringWithPrefix("network-1-chainID-2-pub-key"),
			shouldError: false,
		},
		{
			network:     "NETWORK-2",
			chainID:     "NETWORK-2-CHAINID-1",
			url:         "http://network-2-chainID-1.test",
			pubKey:      GetDummyKeyStringWithPrefix("network-2-chainID-1-pub-key"),
			shouldError: false,
		},
		{
			shouldError:   true,
			expectedError: "network cannot be empty",
		},
		{
			network:       "ERROR",
			shouldError:   true,
			expectedError: "chainID cannot be empty",
		},
		{
			network:       "ERROR",
			chainID:       "ERROR",
			shouldError:   true,
			expectedError: "URL cannot be empty",
		},
		{
			network:       "ERROR",
			chainID:       "ERROR",
			url:           "http://error.test",
			shouldError:   true,
			expectedError: "cannot add telemetry endpoint, ServerPubKey cannot be empty",
		},
		{
			network:       "NETWORK-1",
			chainID:       "NETWORK-1-CHAINID-1",
			url:           "http://network-1-chainID-1.test",
			pubKey:        GetDummyKeyStringWithPrefix("network-1-chainID-1-pub-key"),
			shouldError:   true,
			expectedError: "endpoint already exists",
		},
	}

	var mockEndpoints []config.TelemetryIngressEndpoint

	for _, e := range endpoints {
		te := mocks.NewTelemetryIngressEndpoint(t)
		te.On("Network").Maybe().Return(e.network)
		te.On("ChainID").Maybe().Return(e.chainID)
		te.On("ServerPubKey").Maybe().Return(e.pubKey)

		u, _ := url.Parse(e.url)
		if e.url == "" {
			u = nil
		}
		te.On("URL").Maybe().Return(u)
		mockEndpoints = append(mockEndpoints, te)
	}

	tic := setupMockConfig(t, true)
	tic.On("Endpoints").Return(mockEndpoints)

	lggr, logObs := logger.TestLoggerObserved(t, zapcore.InfoLevel)

	ks := mocks3.NewCSA(t)

	ks.On("GetAll").Return([]csakey.KeyV2{csakey.MustNewV2XXXTestingOnly(big.NewInt(0))}, nil)

	m := NewManager(tic, ks, lggr)

	require.Equal(t, uint(123), m.bufferSize)
	require.Equal(t, ks, m.ks)
	require.Equal(t, "TelemetryManager", m.lggr.Name())
	require.Equal(t, true, m.logging)
	require.Equal(t, uint(51), m.maxBatchSize)
	require.Equal(t, time.Millisecond*512, m.sendInterval)
	require.Equal(t, time.Second*7, m.sendTimeout)
	require.Equal(t, true, m.uniConn)
	require.Equal(t, true, m.useBatchSend)

	logs := logObs.TakeAll()
	for i, e := range endpoints {
		if !e.shouldError {
			require.Equal(t, e.network, m.endpoints[i].Network)
			require.Equal(t, e.chainID, m.endpoints[i].ChainID)
			require.Equal(t, e.pubKey, m.endpoints[i].PubKey)
			require.Equal(t, e.url, m.endpoints[i].URL.String())
		} else {
			found := false
			for _, l := range logs {
				if strings.Contains(l.Message, e.expectedError) {
					found = true
				}
			}
			require.Equal(t, true, found, "cannot find log: %s", e.expectedError)
		}
	}

	require.Equal(t, "TelemetryManager", m.Name())

	require.Nil(t, m.Start(testutils.Context(t)))
	t.Cleanup(func() {
		require.NoError(t, m.Close())
	})
	testutils.WaitForLogMessageCount(t, logObs, "error connecting error while dialing dial tcp", 3)

	hr := m.HealthReport()
	require.Equal(t, 4, len(hr))
}

func TestCorrectEndpointRouting(t *testing.T) {
	tic := setupMockConfig(t, true)
	tic.On("Endpoints").Return(nil)

	lggr, obsLogs := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	ks := mocks3.NewCSA(t)

	tm := NewManager(tic, ks, lggr)

	type testEndpoint struct {
		network string
		chainID string
	}

	testEndpoints := []testEndpoint{
		{
			network: "NETWORK-1",
			chainID: "NETWORK-1-CHAINID-1",
		},
		{
			network: "NETWORK-1",
			chainID: "NETWORK-1-CHAINID-2",
		},
		{
			network: "NETWORK-2",
			chainID: "NETWORK-2-CHAINID-1",
		},
		{
			network: "NETWORK-2",
			chainID: "NETWORK-2-CHAINID-2",
		},
	}

	tm.endpoints = make([]*telemetryEndpoint, len(testEndpoints))
	clientSent := make([]synchronization.TelemPayload, 0)
	for i, e := range testEndpoints {
		clientMock := mocks2.NewTelemetryService(t)
		clientMock.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
			clientSent = append(clientSent, synchronization.TelemPayload{
				Telemetry:  args[1].([]byte),
				ContractID: args[2].(string),
				TelemType:  args[3].(synchronization.TelemetryType),
			})
		})

		tm.endpoints[i] = &telemetryEndpoint{
			ChainID: e.chainID,
			Network: e.network,
			client:  clientMock,
		}
	}
	//Unknown networks or chainID
	noopEndpoint := tm.GenMonitoringEndpoint("unknown-network", "unknown-chainID", "some-contractID", "some-type")
	require.Equal(t, "*telemetry.NoopAgent", reflect.TypeOf(noopEndpoint).String())
	require.Equal(t, 1, obsLogs.Len())
	require.Contains(t, obsLogs.TakeAll()[0].Message, "no telemetry endpoint found")

	noopEndpoint = tm.GenMonitoringEndpoint("network-1", "unknown-chainID", "some-contractID", "some-type")
	require.Equal(t, "*telemetry.NoopAgent", reflect.TypeOf(noopEndpoint).String())
	require.Equal(t, 1, obsLogs.Len())
	require.Contains(t, obsLogs.TakeAll()[0].Message, "no telemetry endpoint found")

	noopEndpoint = tm.GenMonitoringEndpoint("network-2", "network-1-chainID-1", "some-contractID", "some-type")
	require.Equal(t, "*telemetry.NoopAgent", reflect.TypeOf(noopEndpoint).String())
	require.Equal(t, 1, obsLogs.Len())
	require.Contains(t, obsLogs.TakeAll()[0].Message, "no telemetry endpoint found")

	//Known networks and chainID
	for i, e := range testEndpoints {
		telemType := fmt.Sprintf("TelemType_%s", e.chainID)
		contractID := fmt.Sprintf("contractID_%s", e.chainID)
		me := tm.GenMonitoringEndpoint(
			e.network,
			e.chainID,
			contractID,
			synchronization.TelemetryType(telemType),
		)
		me.SendLog([]byte(e.chainID))
		require.Equal(t, 0, obsLogs.Len())

		require.Equal(t, i+1, len(clientSent))
		require.Equal(t, contractID, clientSent[i].ContractID)
		require.Equal(t, telemType, string(clientSent[i].TelemType))
		require.Equal(t, []byte(e.chainID), clientSent[i].Telemetry)
	}
}
