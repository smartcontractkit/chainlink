package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"

	"github.com/smartcontractkit/chainlink/v2/core/config"
)

// stubP2P is a minimal config.P2P that returns a configurable PeerID and zero values for the
// fields confidentialRelayPeerID never reads.
type stubP2P struct {
	peerID p2pkey.PeerID
}

func (s stubP2P) Enabled() bool                   { return false }
func (s stubP2P) PeerID() p2pkey.PeerID           { return s.peerID }
func (s stubP2P) V2() config.V2                   { return nil }
func (s stubP2P) IncomingMessageBufferSize() int  { return 0 }
func (s stubP2P) OutgoingMessageBufferSize() int  { return 0 }
func (s stubP2P) TraceLogging() bool              { return false }
func (s stubP2P) EnableExperimentalRageP2P() bool { return false }

// stubCapabilities is a minimal config.Capabilities whose .Peering() returns a configurable
// stubP2P. Other methods return nil because confidentialRelayPeerID does not read them.
type stubCapabilities struct {
	peering config.P2P
}

func (s stubCapabilities) RateLimit() config.EngineExecutionRateLimit            { return nil }
func (s stubCapabilities) Peering() config.P2P                                   { return s.peering }
func (s stubCapabilities) SharedPeering() config.SharedPeering                   { return nil }
func (s stubCapabilities) Dispatcher() config.Dispatcher                         { return nil }
func (s stubCapabilities) ExternalRegistry() config.CapabilitiesExternalRegistry { return nil }
func (s stubCapabilities) WorkflowRegistry() config.CapabilitiesWorkflowRegistry { return nil }
func (s stubCapabilities) GatewayConnector() config.GatewayConnector             { return nil }
func (s stubCapabilities) Local() config.LocalCapabilities                       { return nil }

// stubConfig is a minimal cre.Config whose .P2P() returns a configurable stubP2P. Other
// methods return nil because confidentialRelayPeerID does not read them.
type stubConfig struct {
	p2p config.P2P
}

func (s stubConfig) Billing() config.Billing           { return nil }
func (s stubConfig) Capabilities() config.Capabilities { return nil }
func (s stubConfig) Workflows() config.Workflows       { return nil }
func (s stubConfig) CRE() config.CRE                   { return nil }
func (s stubConfig) P2P() config.P2P                   { return s.p2p }
func (s stubConfig) Sharding() config.Sharding         { return nil }

func peerIDFromByte(b byte) p2pkey.PeerID {
	var id p2pkey.PeerID
	id[0] = b
	return id
}

func TestConfidentialRelayPeerID_PrefersCapabilitiesPeering(t *testing.T) {
	t.Parallel()

	capPeerID := peerIDFromByte(0x11)
	nodePeerID := peerIDFromByte(0x22)

	cfg := stubConfig{p2p: stubP2P{peerID: nodePeerID}}
	capCfg := stubCapabilities{peering: stubP2P{peerID: capPeerID}}

	require.Equal(t, capPeerID, confidentialRelayPeerID(cfg, capCfg))
}

func TestConfidentialRelayPeerID_FallsBackToNodeP2PWhenCapabilitiesPeerIDIsZero(t *testing.T) {
	t.Parallel()

	nodePeerID := peerIDFromByte(0x33)

	cfg := stubConfig{p2p: stubP2P{peerID: nodePeerID}}
	capCfg := stubCapabilities{peering: stubP2P{peerID: p2pkey.PeerID{}}}

	require.Equal(t, nodePeerID, confidentialRelayPeerID(cfg, capCfg))
}

func TestConfidentialRelayPeerID_ReturnsZeroWhenBothAreUnset(t *testing.T) {
	t.Parallel()

	cfg := stubConfig{p2p: stubP2P{peerID: p2pkey.PeerID{}}}
	capCfg := stubCapabilities{peering: stubP2P{peerID: p2pkey.PeerID{}}}

	// Falling through to GetOrFirst(zero) is the documented contract; the resolved peerID
	// here is the zero value and the keystore decides what that means at lookup time.
	require.Equal(t, p2pkey.PeerID{}, confidentialRelayPeerID(cfg, capCfg))
}
