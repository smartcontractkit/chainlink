package telemetry

import (
	"context"
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
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	me := tm.GenMonitoringEndpoint("", "", "network-1", "network-1-chainID-1")
	require.Equal(t, "*telemetry.IngressAgentBatch", reflect.TypeOf(me).String())

	tic = setupMockConfig(t, false)
	tic.On("Endpoints").Return([]config.TelemetryIngressEndpoint{te})
	tm = NewManager(tic, ks, lggr)
	require.Equal(t, "*synchronization.telemetryIngressClient", reflect.TypeOf(tm.endpoints[0].client).String())
	me = tm.GenMonitoringEndpoint("", "", "network-1", "network-1-chainID-1")
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
			pubKey:      "network-1-chainID-1-pub-key",
			shouldError: false,
		},
		{
			network:     "NETWORK-1",
			chainID:     "NETWORK-1-CHAINID-2",
			url:         "http://network-1-chainID-2.test",
			pubKey:      "network-1-chainID-2-pub-key",
			shouldError: false,
		},
		{
			network:     "NETWORK-2",
			chainID:     "NETWORK-2-CHAINID-1",
			url:         "http://network-2-chainID-1.test",
			pubKey:      "network-2-chainID-1-pub-key",
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
			pubKey:        "network-1-chainID-1-pub-key",
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

	require.Nil(t, m.Start(context.Background()))
	testutils.WaitForLogMessageCount(t, logObs, "error connecting error while dialing dial tcp: lookup http: no such host, waiting then retrying", 3)

	hr := m.HealthReport()
	require.Equal(t, 4, len(hr))
	require.Nil(t, m.Close())
	time.Sleep(time.Second * 1)
}

func TestCorrectEndpointRouting(t *testing.T) {
	tic := setupMockConfig(t, true)
	tic.On("Endpoints").Return(nil)

	lggr, obsLogs := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	ks := mocks3.NewCSA(t)

	tm := NewManager(tic, ks, lggr)
	tm.endpoints = make([]*telemetryEndpoint, 4)

	clientN1c1 := mocks2.NewTelemetryIngressBatchClient(t)
	clientN1c1Sent := make([]synchronization.TelemPayload, 0)
	clientN1c1.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		clientN1c1Sent = append(clientN1c1Sent, args[0].(synchronization.TelemPayload))
	})
	tm.endpoints[0] = &telemetryEndpoint{
		StartStopOnce: utils.StartStopOnce{},
		ChainID:       "NETWORK-1-CHAINID-1",
		Network:       "NETWORK-1",
		client:        clientN1c1,
	}

	clientN1c2 := mocks2.NewTelemetryIngressBatchClient(t)
	clientN1c2Sent := make([]synchronization.TelemPayload, 0)
	clientN1c2.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		clientN1c2Sent = append(clientN1c2Sent, args[0].(synchronization.TelemPayload))
	})
	tm.endpoints[1] = &telemetryEndpoint{
		StartStopOnce: utils.StartStopOnce{},
		ChainID:       "NETWORK-1-CHAINID-2",
		Network:       "NETWORK-1",
		client:        clientN1c2,
	}

	clientN2c1 := mocks2.NewTelemetryIngressBatchClient(t)
	clientN2c1Sent := make([]synchronization.TelemPayload, 0)
	clientN2c1.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		clientN2c1Sent = append(clientN2c1Sent, args[0].(synchronization.TelemPayload))
	})
	tm.endpoints[2] = &telemetryEndpoint{
		StartStopOnce: utils.StartStopOnce{},
		ChainID:       "NETWORK-2-CHAINID-1",
		Network:       "NETWORK-2",
		client:        clientN2c1,
	}

	clientN2c2 := mocks2.NewTelemetryIngressBatchClient(t)
	clientN2c2Sent := make([]synchronization.TelemPayload, 0)
	clientN2c2.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		clientN2c2Sent = append(clientN2c2Sent, args[0].(synchronization.TelemPayload))
	})
	tm.endpoints[3] = &telemetryEndpoint{
		StartStopOnce: utils.StartStopOnce{},
		ChainID:       "NETWORK-2-CHAINID-2",
		Network:       "NETWORK-2",
		client:        clientN2c2,
	}

	//Unknown networks or chainID
	noopEndpoint := tm.GenMonitoringEndpoint("some-contractID", "some-type", "unknown-network", "unknown-chainID")
	require.Equal(t, "*telemetry.NoopAgent", reflect.TypeOf(noopEndpoint).String())
	require.Equal(t, 1, obsLogs.Len())
	require.Contains(t, obsLogs.TakeAll()[0].Message, "no telemetry endpoint found")

	noopEndpoint = tm.GenMonitoringEndpoint("some-contractID", "some-type", "network-1", "unknown-chainID")
	require.Equal(t, "*telemetry.NoopAgent", reflect.TypeOf(noopEndpoint).String())
	require.Equal(t, 1, obsLogs.Len())
	require.Contains(t, obsLogs.TakeAll()[0].Message, "no telemetry endpoint found")

	noopEndpoint = tm.GenMonitoringEndpoint("some-contractID", "some-type", "network-2", "network-1-chainID-1")
	require.Equal(t, "*telemetry.NoopAgent", reflect.TypeOf(noopEndpoint).String())
	require.Equal(t, 1, obsLogs.Len())
	require.Contains(t, obsLogs.TakeAll()[0].Message, "no telemetry endpoint found")

	//Known networks and chainID
	n1c1 := tm.GenMonitoringEndpoint("contractID-network-1-chainID-1", "telemType1", "network-1", "network-1-chainID-1")
	n1c2 := tm.GenMonitoringEndpoint("contractID-network-1-chainID-2", "telemType2", "network-1", "network-1-chainID-2")
	n2c1 := tm.GenMonitoringEndpoint("contractID-network-2-chainID-1", "telemType3", "network-2", "network-2-chainID-1")
	n2c2 := tm.GenMonitoringEndpoint("contractID-network-2-chainID-2", "telemType4", "network-2", "network-2-chainID-2")
	require.Equal(t, 0, obsLogs.Len())

	n1c1.SendLog([]byte("network-1-chainID-1"))
	require.Equal(t, 1, len(clientN1c1Sent))
	require.Equal(t, "network-1", clientN1c1Sent[0].Network)
	require.Equal(t, "network-1-chainID-1", clientN1c1Sent[0].ChainID)
	require.Equal(t, "contractID-network-1-chainID-1", clientN1c1Sent[0].ContractID)
	require.Equal(t, "telemType1", string(clientN1c1Sent[0].TelemType))
	require.Equal(t, []byte("network-1-chainID-1"), clientN1c1Sent[0].Telemetry)

	n1c2.SendLog([]byte("network-1-chainID-2"))
	require.Equal(t, 1, len(clientN1c2Sent))
	require.Equal(t, "network-1", clientN1c2Sent[0].Network)
	require.Equal(t, "network-1-chainID-2", clientN1c2Sent[0].ChainID)
	require.Equal(t, "contractID-network-1-chainID-2", clientN1c2Sent[0].ContractID)
	require.Equal(t, "telemType2", string(clientN1c2Sent[0].TelemType))
	require.Equal(t, []byte("network-1-chainID-2"), clientN1c2Sent[0].Telemetry)

	n2c1.SendLog([]byte("network-2-chainID-1"))
	require.Equal(t, 1, len(clientN2c1Sent))
	require.Equal(t, "network-2", clientN2c1Sent[0].Network)
	require.Equal(t, "network-2-chainID-1", clientN2c1Sent[0].ChainID)
	require.Equal(t, "contractID-network-2-chainID-1", clientN2c1Sent[0].ContractID)
	require.Equal(t, "telemType3", string(clientN2c1Sent[0].TelemType))
	require.Equal(t, []byte("network-2-chainID-1"), clientN2c1Sent[0].Telemetry)

	n2c2.SendLog([]byte("network-2-chainID-2"))
	require.Equal(t, 1, len(clientN2c2Sent))
	require.Equal(t, "network-2", clientN2c2Sent[0].Network)
	require.Equal(t, "network-2-chainID-2", clientN2c2Sent[0].ChainID)
	require.Equal(t, "contractID-network-2-chainID-2", clientN2c2Sent[0].ContractID)
	require.Equal(t, "telemType4", string(clientN2c2Sent[0].TelemType))
	require.Equal(t, []byte("network-2-chainID-2"), clientN2c2Sent[0].Telemetry)

}
