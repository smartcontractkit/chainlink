package p2p_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/prometheus/client_golang/prometheus"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
)

func TestPeer_CleanStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	port := freeport.GetOne(t)
	privKey, _ := newKeyPair(t)

	reg := prometheus.NewRegistry()
	peerConfig := p2p.PeerConfig{
		PrivateKey:      privKey,
		ListenAddresses: []string{fmt.Sprintf("127.0.0.1:%d", port)},

		DeltaReconcile:     time.Second * 5,
		DeltaDial:          time.Second * 5,
		DiscovererDatabase: p2p.NewInMemoryDiscovererDatabase(),
		MetricsRegisterer:  reg,
	}

	peer, err := p2p.NewPeer(peerConfig, lggr)
	require.NoError(t, err)
	err = peer.Start(testutils.Context(t))
	require.NoError(t, err)
	err = peer.Close()
	require.NoError(t, err)
}

func newKeyPair(t *testing.T) (ed25519.PrivateKey, ragetypes.PeerID) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	peerID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	require.NoError(t, err)
	return privKey, peerID
}
