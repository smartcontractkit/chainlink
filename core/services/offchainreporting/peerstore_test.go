package offchainreporting_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func Test_Peerstore_Start(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Deferring the constraint avoids having to insert an entire set of jobs/specs
	require.NoError(t, store.DB.Exec(`SET CONSTRAINTS p2p_peers_job_id_fkey DEFERRED`).Error)
	err := store.DB.Exec(`INSERT INTO p2p_peers (id, addr, created_at, updated_at, job_id) VALUES
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.1/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		1
	),
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		1
	),
	(
		'12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		'/ip4/127.0.0.2/tcp/12000/p2p/12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph',
		NOW(),
		NOW(),
		2
	)
	`).Error
	require.NoError(t, err)

	wrapper, err := offchainreporting.NewPeerstoreWrapper(context.Background(), store.DB, 1*time.Second, 1)
	require.NoError(t, err)

	err = wrapper.Start()
	require.NoError(t, err)

	require.Equal(t, 1, wrapper.Peerstore.PeersWithAddrs().Len())

	peerID, err := p2ppeer.Decode("12D3KooWL1yndUw9T2oWXjhfjdwSscWA78YCpUdduA3Cnn4dCtph")
	require.NoError(t, err)

	maddrs := wrapper.Peerstore.Addrs(peerID)

	require.Len(t, maddrs, 2)
}
