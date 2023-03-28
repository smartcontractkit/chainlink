package ocrcommon_test

import (
	"testing"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/stretchr/testify/require"

	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_Peerstore_Start(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	peerID, err := p2ppeer.Decode("12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X")
	require.NoError(t, err)

	nonExistentP2PPeerID, err := p2ppeer.Decode("12D3KooWAdCzaesXyezatDzgGvCngqsBqoUqnV9PnVc46jsVt2i9")
	require.NoError(t, err)

	err = utils.JustError(db.Exec(`INSERT INTO p2p_peers (id, addr, created_at, updated_at, peer_id) VALUES
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.1/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		$1
	),
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		$2
	),
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		$3
	)
	`, p2pkey.PeerID(peerID), p2pkey.PeerID(peerID), p2pkey.PeerID(nonExistentP2PPeerID)))
	require.NoError(t, err)

	cfg := configtest.NewTestGeneralConfig(t)
	wrapper, err := ocrcommon.NewPeerstoreWrapper(db, 1*time.Second, p2pkey.PeerID(peerID), logger.TestLogger(t), cfg)
	require.NoError(t, err)

	err = wrapper.Start()
	require.NoError(t, err)

	require.Equal(t, 1, wrapper.Peerstore.PeersWithAddrs().Len())

	peerID, err = p2ppeer.Decode("12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)

	maddrs := wrapper.Peerstore.Addrs(peerID)

	require.Len(t, maddrs, 2)
}

func Test_Peerstore_WriteToDB(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	peerID, err := p2ppeer.Decode("12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X")
	require.NoError(t, err)

	cfg := configtest.NewTestGeneralConfig(t)
	wrapper, err := ocrcommon.NewPeerstoreWrapper(db, 1*time.Second, p2pkey.PeerID(peerID), logger.TestLogger(t), cfg)
	require.NoError(t, err)

	maddr, err := ma.NewMultiaddr("/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)
	newPeerID, err := p2ppeer.Decode("12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)

	wrapper.Peerstore.AddAddr(newPeerID, maddr, p2ppeerstore.PermanentAddrTTL)

	err = wrapper.WriteToDB()
	require.NoError(t, err)

	peers := make([]ocrcommon.P2PPeer, 0)
	err = db.Select(&peers, `SELECT * FROM p2p_peers`)
	require.NoError(t, err)
	require.Equal(t, 1, len(peers))

	peer := peers[0]
	require.Equal(t, "12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph", peer.ID)
	require.Equal(t, "/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph", peer.Addr)
	require.Equal(t, p2pkey.PeerID(peerID).Raw(), peer.PeerID)
}
