package p2p_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestDon2DonSharedPeer_WithRealSingletonPeerWrapper(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	k, err := keyStore.P2P().Create(t.Context())
	require.NoError(t, err)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.P2P.V2.Enabled = new(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = new(k.PeerID())
	})
	pw := ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), db, logger.TestLogger(t))

	require.NoError(t, pw.Start(t.Context()))
	defer pw.Close()

	require.Equal(t, k.PeerID(), pw.PeerID)
	sp := p2p.NewDon2DonSharedPeer(pw, nil, logger.TestLogger(t))
	require.NoError(t, sp.Start(t.Context()))

	_, peerID2 := newKeyPair(t)
	_, peerID3 := newKeyPair(t)
	donPairs := []p2ptypes.DonPair{{
		{ID: 1, Members: []ragetypes.PeerID{ragetypes.PeerID(k.PeerID()), peerID2}},
		{ID: 2, Members: []ragetypes.PeerID{peerID2, peerID3}},
	}}
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.NoError(t, sp.Close())
}

func TestDon2DonSharedPeer_ErrorOnNilSingletonPeerWrapper(t *testing.T) {
	sp := p2p.NewDon2DonSharedPeer(nil, nil, logger.TestLogger(t))
	require.Error(t, sp.Start(t.Context()))
}

func TestDon2DonSharedPeer_UpdateConnectionsByDONs(t *testing.T) {
	pw := ocrcommon.NewSingletonPeerWrapper(nil, nil, nil, nil, logger.TestLogger(t)) // nils are ok, we won't Start() it
	_, myPeerID := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	_, peerID3 := newKeyPair(t)
	_, peerID4 := newKeyPair(t)
	mockPGFactory := mockPeerGroupFactory{}
	pw.PeerGroupFactory = &mockPGFactory
	pw.PeerID = p2pkey.PeerID(myPeerID)

	sp := p2p.NewDon2DonSharedPeer(pw, nil, logger.TestLogger(t))
	require.NoError(t, sp.Start(t.Context()))

	donPairs := []p2ptypes.DonPair{{
		{ID: 1, Members: []ragetypes.PeerID{myPeerID, peerID2}},
		{ID: 2, Members: []ragetypes.PeerID{peerID2, peerID3}},
	}}
	// Adding a new DON pair
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.Equal(t, 1, mockPGFactory.newDonGroupCounter)
	require.Equal(t, 2, mockPGFactory.newNodeGroupCounter) // myPeer is connected to peers 2 and 3
	require.Equal(t, 0, mockPGFactory.closedGroupCounter)
	require.Equal(t, 2, mockPGFactory.newStreamCounter)

	// No changes expected when updating the same group
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.Equal(t, 1, mockPGFactory.newDonGroupCounter)
	require.Equal(t, 2, mockPGFactory.newNodeGroupCounter)
	require.Equal(t, 0, mockPGFactory.closedGroupCounter)
	require.Equal(t, 2, mockPGFactory.newStreamCounter)

	// Expect a change when DON membership changes
	donPairs[0][1].Members[1] = peerID4
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.Equal(t, 2, mockPGFactory.newDonGroupCounter)  // update of existing group
	require.Equal(t, 3, mockPGFactory.newNodeGroupCounter) // one new connection to peer 4
	require.Equal(t, 2, mockPGFactory.closedGroupCounter)  // close old DON group + peer 2 group
	require.Equal(t, 3, mockPGFactory.newStreamCounter)    // one new connection to peer 4

	// Expect a change when a new DON pair is added
	donPairs = append(donPairs, [2]capabilities.DON{
		{ID: 3, Members: []ragetypes.PeerID{myPeerID, peerID3}},
		{ID: 2, Members: []ragetypes.PeerID{peerID2, peerID3}},
	})
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.Equal(t, 3, mockPGFactory.newDonGroupCounter)
	require.Equal(t, 4, mockPGFactory.newNodeGroupCounter) // re-create connection to peer 2
	require.Equal(t, 2, mockPGFactory.closedGroupCounter)
	require.Equal(t, 4, mockPGFactory.newStreamCounter) // re-create connection to peer 2

	require.NoError(t, sp.Close())
	require.Equal(t, 2+2+3, mockPGFactory.closedGroupCounter) // closed 2 DON groups and 3 node groups
}

// TestDon2DonSharedPeer_DONMembershipChange reproduces the failure
// seen in production when a DON's on-chain membership changes:
//
//	failed to create discovery group: asked to add group with digest we already have
//	(digest: 000f0000000300000005...)
//
// discoveryGroups is keyed by pairID(), which hashes the DON IDs *and* both member lists,
// but the peer group is registered in libocr under donPairDigest(), which is derived from
// the DON IDs alone. That mapping is many-to-one, so a membership change produces a new
// pairID while the digest stays the same: the "already exists" check misses, and libocr
// then rejects the create because it still holds the digest for the previous group.
func TestDon2DonSharedPeer_DONMembershipChange(t *testing.T) {
	t.Parallel()
	pw := ocrcommon.NewSingletonPeerWrapper(nil, nil, nil, nil, logger.TestLogger(t)) // nils are ok, we won't Start() it
	_, myPeerID := newKeyPair(t)
	_, peerID2 := newKeyPair(t)
	_, peerID3 := newKeyPair(t)
	_, peerID4 := newKeyPair(t)
	mockPGFactory := mockPeerGroupFactory{}
	pw.PeerGroupFactory = &mockPGFactory
	pw.PeerID = p2pkey.PeerID(myPeerID)

	sp := p2p.NewDon2DonSharedPeer(pw, nil, logger.TestLogger(t))
	require.NoError(t, sp.Start(t.Context()))
	defer func() { require.NoError(t, sp.Close()) }()

	donPairs := []p2ptypes.DonPair{{
		{ID: 3, Members: []ragetypes.PeerID{myPeerID, peerID2, peerID3}},
		{ID: 5, Members: []ragetypes.PeerID{peerID2, peerID3}},
	}}
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.Equal(t, 1, mockPGFactory.newDonGroupCounter)

	// DON 3's membership changes on chain: peerID3 is replaced by peerID4. The DON IDs are
	// unchanged, so donPairDigest() still yields the digest for the pair (3, 5).
	donPairs[0][0].Members = []ragetypes.PeerID{myPeerID, peerID2, peerID4}
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}),
		"membership change must not fail: the group for the previous membership has to be closed, "+
			"releasing its digest, before a group for the new membership is created")

	// The stale group is gone and exactly one discovery group for the pair remains.
	require.Equal(t, 2, mockPGFactory.newDonGroupCounter, "a discovery group for the new membership should have been created")

	// A further update with unchanged membership must be a no-op.
	require.NoError(t, sp.UpdateConnectionsByDONs(t.Context(), donPairs, p2ptypes.StreamConfig{}))
	require.Equal(t, 2, mockPGFactory.newDonGroupCounter)
}

type mockPeerGroupFactory struct {
	newDonGroupCounter  int // large - more than 2 members
	newNodeGroupCounter int // small - 2 members
	closedGroupCounter  int
	newStreamCounter    int
}

func (m *mockPeerGroupFactory) NewPeerGroup(
	configDigest ocr2types.ConfigDigest,
	peerIDs []string,
	bootstrappers []commontypes.BootstrapperLocator,
) (networking.PeerGroup, error) {
	if len(peerIDs) > 2 {
		m.newDonGroupCounter++
	} else {
		m.newNodeGroupCounter++
	}
	return &mockPeerGroup{groupFactory: m}, nil
}

type mockPeerGroup struct {
	groupFactory *mockPeerGroupFactory
}

func (m *mockPeerGroup) NewStream(remotePeerID string, newStreamArgs networking.NewStreamArgs) (networking.Stream, error) {
	m.groupFactory.newStreamCounter++
	return &mockStream{msgCh: make(chan []byte)}, nil
}

func (m *mockPeerGroup) Close() error {
	m.groupFactory.closedGroupCounter++
	return nil
}

type mockStream struct {
	msgCh chan []byte
}

func (m *mockStream) SendMessage(data []byte) {
	m.msgCh <- data
}

func (m *mockStream) ReceiveMessages() <-chan []byte {
	return m.msgCh
}

func (m *mockStream) Close() error {
	close(m.msgCh)
	return nil
}

func newKeyPair(t *testing.T) (ed25519.PrivateKey, ragetypes.PeerID) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	peerID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	require.NoError(t, err)
	return privKey, peerID
}
