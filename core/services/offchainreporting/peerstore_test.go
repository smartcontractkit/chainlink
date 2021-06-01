package offchainreporting_test

import (
	"testing"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/require"
)

func Test_Peerstore_Start(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Deferring the constraint avoids having to insert an entire set of jobs/specs
	require.NoError(t, store.DB.Exec(`SET CONSTRAINTS p2p_peers_peer_id_fkey DEFERRED`).Error)
	err := store.DB.Exec(`INSERT INTO p2p_peers (id, addr, created_at, updated_at, peer_id) VALUES
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.1/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		?
	),
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
	 	?
	),
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		?
	)
	`, cltest.DefaultP2PPeerID, cltest.DefaultP2PPeerID, cltest.NonExistentP2PPeerID).Error
	require.NoError(t, err)

	wrapper, err := offchainreporting.NewPeerstoreWrapper(store.DB, 1*time.Second, p2pkey.PeerID(cltest.DefaultP2PPeerID))
	require.NoError(t, err)

	err = wrapper.Start()
	require.NoError(t, err)

	require.Equal(t, 1, wrapper.Peerstore.PeersWithAddrs().Len())

	peerID, err := p2ppeer.Decode("12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)

	maddrs := wrapper.Peerstore.Addrs(peerID)

	require.Len(t, maddrs, 2)
}

func Test_Peerstore_WriteToDB(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Deferring the constraint avoids having to insert an entire set of jobs/specs
	require.NoError(t, store.DB.Exec(`SET CONSTRAINTS p2p_peers_peer_id_fkey DEFERRED`).Error)

	wrapper, err := offchainreporting.NewPeerstoreWrapper(store.DB, 1*time.Second, p2pkey.PeerID(cltest.DefaultP2PPeerID))
	require.NoError(t, err)

	maddr, err := ma.NewMultiaddr("/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)
	peerID, err := p2ppeer.Decode("12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)

	wrapper.Peerstore.AddAddr(peerID, maddr, p2ppeerstore.PermanentAddrTTL)

	err = wrapper.WriteToDB()
	require.NoError(t, err)

	peers := make([]offchainreporting.P2PPeer, 0)
	result := store.DB.Find(&peers)
	require.Equal(t, int64(1), result.RowsAffected)

	peer := peers[0]
	require.Equal(t, "12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph", peer.ID)
	require.Equal(t, "/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph", peer.Addr)
	require.Equal(t, cltest.DefaultP2PPeerID.Raw(), peer.PeerID)
}
