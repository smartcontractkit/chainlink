package synchronization

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	telemPb "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

// NewTestTelemetryIngressClient calls NewTelemetryIngressClient and injects telemClient.
func NewTestTelemetryIngressClient(t *testing.T, url *url.URL, serverPubKeyHex string, csaKeyStore keystore.CSA, telemClient telemPb.TelemClient) TelemetryService {
	tc := NewTelemetryIngressClient(url, serverPubKeyHex, csaKeyStore, logger.TestLogger(t), 100)
	tc.(*telemetryIngressClient).telemClient = telemClient
	return tc
}

// NewTestTelemetryIngressBatchClient calls NewTelemetryIngressBatchClient and injects telemClient.
func NewTestTelemetryIngressBatchClient(t *testing.T, url *url.URL, serverPubKeyHex string, csaKeyStore keystore.CSA, logging bool, telemClient telemPb.TelemClient, sendInterval time.Duration, uniconn bool) TelemetryService {
	tc := NewTelemetryIngressBatchClient(url, serverPubKeyHex, csaKeyStore, logging, logger.TestLogger(t), 100, 50, sendInterval, time.Second, uniconn)
	tc.(*telemetryIngressBatchClient).closeFn = func() error { return nil }
	tc.(*telemetryIngressBatchClient).telemClient = telemClient
	return tc
}

// NewTestChipIngressBatchClient calls NewChipIngressBatchClient with a fixed test CSA key and default sizing.
func NewTestChipIngressBatchClient(t *testing.T, chipClient chipingress.Client, logging bool, sendInterval time.Duration) ChipIngressService {
	c, err := NewChipIngressBatchClient(chipClient, "deadbeef", ChipIngressBatchConfig{
		BufferSize:         10_000,
		MaxBatchSize:       1_000,
		MaxConcurrentSends: 10,
		SendInterval:       sendInterval,
		SendTimeout:        time.Second,
		DrainTimeout:       5 * time.Second,
		Logging:            logging,
	}, logger.TestLogger(t))
	require.NoError(t, err)
	return c
}
