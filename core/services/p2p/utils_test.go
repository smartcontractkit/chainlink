package p2p_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
)

const (
	peerID1 = "peer1"
	peerID2 = "peer2"
	ann1    = "announcement1"
	ann2    = "announcement2"
)

func TestInMemoryDiscovererDatabase(t *testing.T) {
	db := p2p.NewInMemoryDiscovererDatabase()
	require.NoError(t, db.StoreAnnouncement(testutils.Context(t), peerID1, []byte(ann1)))
	require.NoError(t, db.StoreAnnouncement(testutils.Context(t), peerID2, []byte(ann2)))
	state, err := db.ReadAnnouncements(testutils.Context(t), []string{peerID1, peerID2})
	require.NoError(t, err)
	require.Equal(t, map[string][]byte{peerID1: []byte(ann1), peerID2: []byte(ann2)}, state)
	state, err = db.ReadAnnouncements(testutils.Context(t), []string{peerID2})
	require.NoError(t, err)
	require.Equal(t, map[string][]byte{peerID2: []byte(ann2)}, state)
}
