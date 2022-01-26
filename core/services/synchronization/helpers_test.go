package synchronization

import (
	"net/url"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
)

// NewTestTelemetryIngressClient calls NewTelemetryIngressClient and injects telemClient.
func NewTestTelemetryIngressClient(t *testing.T, url *url.URL, serverPubKeyHex string, ks keystore.CSA, logging bool, telemClient telemPb.TelemClient) TelemetryIngressClient {
	tc := NewTelemetryIngressClient(url, serverPubKeyHex, ks, logging, logger.TestLogger(t))
	tc.(*telemetryIngressClient).telemClient = telemClient
	return tc
}
