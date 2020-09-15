package networking_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/networking"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func Test_NewPeerstore(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	db := store.DB.DB()

	peerstore, err := networking.NewPeerstore(context.Background(), db)
	require.NoError(t, err)

	peerID, err := p2ppeer.Decode("12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)
	multiaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)

	peerstore.AddAddr(peerID, multiaddr, 1*time.Hour)

	peers := peerstore.Peers()
	require.Len(t, peers, 1)
	require.Equal(t, peers[0], peerID)

	addrs := peerstore.Addrs(peerID)
	require.Len(t, addrs, 1)

	require.Equal(t, multiaddr.String(), addrs[0].String())

	// Instantiate a new one to ensure we read from the DB not just from memory
	peerstore, err = networking.NewPeerstore(context.Background(), db)
	require.NoError(t, err)
	peers = peerstore.Peers()
	require.Len(t, peers, 1)
	require.Equal(t, peers[0], peerID)

	addrs = peerstore.Addrs(peerID)
	require.Len(t, addrs, 1)

	require.Equal(t, multiaddr.String(), addrs[0].String())

}
