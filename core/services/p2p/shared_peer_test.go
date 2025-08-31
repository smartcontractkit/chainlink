package p2p_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestDon2DonSharedPeer_WithRealSingletonPeerWrapper(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	k, err := keyStore.P2P().Create(testutils.Context(t))
	require.NoError(t, err)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = ptr(k.PeerID())
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

func ptr[T any](t T) *T { return &t }
