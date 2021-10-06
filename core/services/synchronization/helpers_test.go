package synchronization

import (
	"net/url"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
)

// NewTestTelemetryIngressClient calls NewTelemetryIngressClient and injects telemClient.
func NewTestTelemetryIngressClient(url *url.URL, serverPubKeyHex string, ks keystore.CSA, logging bool, telemClient telemPb.TelemClient) TelemetryIngressClient {
	tc := NewTelemetryIngressClient(url, serverPubKeyHex, ks, logging)
	tc.(*telemetryIngressClient).telemClient = telemClient
	return tc
}
