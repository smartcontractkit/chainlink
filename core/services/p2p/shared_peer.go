package p2p

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// Don2DonSharedPeer is a wrapper around PeerGroupFactory created by SingletonPeerWrapper.
// It manages peer groups and their streams for Don2Don traffic.
//
// Don2DonSharedPeer creates two types of PeerGroups:
//   - Discovery Groups - each consisting of all peers from exactly two DONs.
//     Those groups are used for peer discovery and won’t have any Streams created within them.
//   - Messaging Groups (“pairs”) - each consisting of exactly two peers - the current peer
//     and one remote peer. Each of those groups will always have exactly one Stream
//     created within them, which will be used to exchange Don2Don messages.
//
// Thread-safety: after calling Start(), Send() can be called concurrently from multiple goroutines,
// UpdateConnectionsByDONs() can be called concurrently with Send() but only from a single goroutine.
type don2DonSharedPeer struct {
	services.Service
	srvcEng              *services.Engine
	singletonPeerWrapper *ocrcommon.SingletonPeerWrapper
	bootstrappers        []ocrcommontypes.BootstrapperLocator
	lggr                 logger.Logger

	// fields derived from config and dependencies
	pgFactory   networking.PeerGroupFactory
	myID        ragetypes.PeerID
	isBootstrap bool

	recvCh          chan p2ptypes.Message
	discoveryGroups map[string]networking.PeerGroup // keyed by donPairHash()
	remotePeers     map[ragetypes.PeerID]*remotePeer
	mu              sync.RWMutex // protects discoveryGroups and remotePeers

	metrics *SharedPeerMetrics
}

var _ p2ptypes.SharedPeer = &don2DonSharedPeer{}

const (
	// TODO(CRE-755): replace these with don2don values when available in libocr
	streamNamePrefix = "ccip-rmn/"
	digestPrefix     = ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo
)

type remotePeer struct {
	// A PeerGroup with exactly two members, connecting our peer with a single remote peer.
	peerPairGroup networking.PeerGroup
	// Stream managed by the PeerGroup.
	stream networking.Stream
}

func NewDon2DonSharedPeer(singletonPeerWrapper *ocrcommon.SingletonPeerWrapper, bootstrappers []ocrcommontypes.BootstrapperLocator, lggr logger.Logger) *don2DonSharedPeer {
	sp := &don2DonSharedPeer{
		singletonPeerWrapper: singletonPeerWrapper,
		bootstrappers:        bootstrappers,
		recvCh:               make(chan p2ptypes.Message, defaultRecvChSize),
		discoveryGroups:      make(map[string]networking.PeerGroup),
		remotePeers:          make(map[ragetypes.PeerID]*remotePeer),
		lggr:                 lggr.Named("Don2DonSharedPeer"),
	}
	sp.Service, sp.srvcEng = services.Config{
		Name:  "Don2DonSharedPeer",
		Start: sp.start,
		Close: sp.close,
	}.NewServiceEngine(sp.lggr)
	return sp
}

func (sp *don2DonSharedPeer) start(ctx context.Context) error {
	sp.lggr.Info("Starting Don2DonSharedPeer ...")
	sp.pgFactory = sp.singletonPeerWrapper.PeerGroupFactory
	if sp.pgFactory == nil {
		return errors.New("PeerGroupFactory is not set in SingletonPeerWrapper. It's possible that SingletonPeerWrapper was not started before Don2DonSharedPeer or somehow failed to initialize")
	}
	sp.myID = ragetypes.PeerID(sp.singletonPeerWrapper.PeerID)
	if (sp.myID == ragetypes.PeerID{}) {
		return errors.New("PeerID is not set in SingletonPeerWrapper. Was it started before Don2DonSharedPeer?")
	}
	myIDStr := sp.myID.String()
	for _, bootstrapper := range sp.bootstrappers {
		if bootstrapper.PeerID == myIDStr {
			sp.isBootstrap = true
			break
		}
	}
	metrics, err := initSharedPeerMetrics()
	if err != nil {
		return fmt.Errorf("failed to init Don2DonSharedPeer metrics: %w", err)
	}
	sp.metrics = metrics
	sp.lggr.Infow("Started Don2DonSharedPeer", "peerId", myIDStr, "isBootstrap", sp.isBootstrap)
	return nil
}

func (sp *don2DonSharedPeer) close() error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.lggr.Info("Closing Don2DonSharedPeer ...")
	for name, group := range sp.discoveryGroups {
		if group.Close() != nil {
			sp.lggr.Errorw("failed to close discovery group", "name", name)
		}
	}
	sp.discoveryGroups = make(map[string]networking.PeerGroup)
	for pid, peer := range sp.remotePeers {
		if peer.peerPairGroup.Close() != nil {
			sp.lggr.Errorw("failed to close messaging group peer", "remotePeerID", pid.String())
		}
	}
	sp.remotePeers = make(map[ragetypes.PeerID]*remotePeer)
	close(sp.recvCh)
	sp.lggr.Info("Closed Don2DonSharedPeer")
	return nil
}

func (sp *don2DonSharedPeer) ID() ragetypes.PeerID {
	return sp.myID
}

func (sp *don2DonSharedPeer) IsBootstrap() bool {
	return sp.isBootstrap
}

func (sp *don2DonSharedPeer) Send(peerID ragetypes.PeerID, msg []byte) error {
	sp.mu.RLock()
	rp, ok := sp.remotePeers[peerID]
	sp.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no stream to remote peer id: %s, from local peer id: %s", peerID.String(), sp.myID.String())
	}
	rp.stream.SendMessage(msg)
	return nil
}

func (sp *don2DonSharedPeer) Receive() <-chan p2ptypes.Message {
	return sp.recvCh
}

func (sp *don2DonSharedPeer) UpdateConnections(peers map[ragetypes.PeerID]p2ptypes.StreamConfig) error {
	return errors.New("UpdateConnections is not supported, use UpdateConnectionsByDON instead")
}

func (sp *don2DonSharedPeer) UpdateConnectionsByDONs(ctx context.Context, donPairs []p2ptypes.DonPair, streamConfig p2ptypes.StreamConfig) error {
	sp.lggr.Infow("UpdateConnectionsByDONs", "numDonPairs", len(donPairs))
	startTs := time.Now().UnixMilli()

	desiredDONPairsIDs := make(map[string]struct{})
	for _, dp := range donPairs {
		pairID := pairID(dp[0], dp[1])
		desiredDONPairsIDs[pairID] = struct{}{}
	}
	desiredRemotePeers := make(map[ragetypes.PeerID]struct{})
	for _, dp := range donPairs {
		if slices.Contains(dp[0].Members, sp.myID) {
			for _, pid := range dp[1].Members {
				desiredRemotePeers[pid] = struct{}{}
			}
		}
		if slices.Contains(dp[1].Members, sp.myID) {
			for _, pid := range dp[0].Members {
				desiredRemotePeers[pid] = struct{}{}
			}
		}
	}
	delete(desiredRemotePeers, sp.myID) // don't accidentally create a stream to ourselves

	// Most of the time, there are no changes to peer groups and we can short-circuit.
	// If changes are needed, let's perform updates under a write lock.
	// A large change is expected only once, on startup, when all groups are created.
	if !sp.stateEqual(desiredDONPairsIDs, desiredRemotePeers) {
		err := sp.updateConnections(donPairs, desiredDONPairsIDs, desiredRemotePeers, streamConfig)
		if err != nil {
			sp.metrics.groupUpdateFailureCounter.Add(ctx, 1)
			return fmt.Errorf("failed to update connections by DONs: %w", err)
		}
	}

	sp.mu.RLock()
	defer sp.mu.RUnlock()
	sp.metrics.discoveryGroups.Record(ctx, int64(len(sp.discoveryGroups)))
	sp.metrics.messagingGroups.Record(ctx, int64(len(sp.remotePeers)))
	sp.metrics.groupUpdateDurationMs.Record(ctx, time.Now().UnixMilli()-startTs)
	sp.lggr.Info("UpdateConnectionsByDONs done")
	return nil
}

// Return true if current state is equal to desider state.
func (sp *don2DonSharedPeer) stateEqual(desiredDONPairsIDs map[string]struct{}, desiredRemotePeers map[ragetypes.PeerID]struct{}) bool {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	if len(sp.discoveryGroups) != len(desiredDONPairsIDs) || len(sp.remotePeers) != len(desiredRemotePeers) {
		return false
	}
	for donPairID := range desiredDONPairsIDs {
		if _, ok := sp.discoveryGroups[donPairID]; !ok {
			return false
		}
	}
	for pid := range desiredRemotePeers {
		if _, ok := sp.remotePeers[pid]; !ok {
			return false
		}
	}
	return true
}

func (sp *don2DonSharedPeer) updateConnections(donPairs []p2ptypes.DonPair, desiredDONPairsIDs map[string]struct{}, desiredRemotePeers map[ragetypes.PeerID]struct{}, streamConfig p2ptypes.StreamConfig) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	// Phase 1 - create discovery groups. Those groups consist of peers from exactly two DONs.
	for _, dp := range donPairs {
		pairID := pairID(dp[0], dp[1])
		if _, ok := sp.discoveryGroups[pairID]; ok {
			continue // group already exists, no need to recreate it
		}
		peerSet := make(map[string]struct{})
		for _, pid := range dp[0].Members {
			peerSet[pid.String()] = struct{}{}
		}
		for _, pid := range dp[1].Members {
			peerSet[pid.String()] = struct{}{}
		}
		peers := make([]string, 0, len(peerSet))
		for pidStr := range peerSet {
			peers = append(peers, pidStr) // no duplicate peers
		}

		if _, ok := sp.discoveryGroups[pairID]; !ok {
			digest := donPairDigest(dp[0].ID, dp[1].ID)
			peerGroup, err := sp.pgFactory.NewPeerGroup(digest, peers, sp.bootstrappers)
			if err != nil {
				sp.lggr.Errorw("failed to create discovery group", "digest", digest, "err", err)
				return fmt.Errorf("failed to create discovery group: %w", err)
			}
			sp.lggr.Infow("Created discovery group", "donPairId", pairID, "don1", dp[0].ID, "don2", dp[1].ID, "numPeers", len(peers))
			sp.discoveryGroups[pairID] = peerGroup
		}
	}

	// Remove obsolete groups
	for donPairID := range sp.discoveryGroups {
		if _, ok := desiredDONPairsIDs[donPairID]; !ok {
			peerGroup := sp.discoveryGroups[donPairID]
			err := peerGroup.Close()
			if err != nil {
				sp.lggr.Errorw("failed to close discovery group", "donPairId", donPairID, "err", err)
			}
			delete(sp.discoveryGroups, donPairID)
			sp.lggr.Infow("Closed discovery group", "donPairId", donPairID)
		}
	}

	// Phase 2 - create messaging groups for Don2Don streams. Each group consists of exactly two Peers.
	if !sp.isBootstrap {
		for remotePID := range desiredRemotePeers {
			if _, exists := sp.remotePeers[remotePID]; exists {
				continue // stream to this peer already exists
			}
			// Create a group with only our peer and the remote peer
			peers := []string{sp.myID.String(), remotePID.String()}
			digest := nodePairDigest(sp.myID, remotePID)
			peerGroup, err := sp.pgFactory.NewPeerGroup(digest, peers, sp.bootstrappers)
			if err != nil {
				sp.lggr.Errorw("failed to create remote peer group", "digest", digest, "err", err)
				return fmt.Errorf("failed to create remote peer group: %w", err)
			}
			idStr1, idStr2 := peerIDStrings(sp.myID, remotePID)
			cfg := networking.NewStreamArgs1{
				StreamName:         streamNamePrefix + idStr1 + "-" + idStr2,
				OutgoingBufferSize: streamConfig.OutgoingMessageBufferSize,
				IncomingBufferSize: streamConfig.IncomingMessageBufferSize,
				MaxMessageLength:   streamConfig.MaxMessageLenBytes,
				MessagesLimit:      streamConfig.MessageRateLimiter,
				BytesLimit:         streamConfig.BytesRateLimiter,
			}
			stream, err := peerGroup.NewStream(remotePID.String(), cfg)
			if err != nil {
				sp.lggr.Errorw("failed to create stream for remote peer", "peerID", remotePID, "err", err)
				peerGroup.Close()
				return fmt.Errorf("failed to create stream for remote peer: %w", err)
			}
			sp.remotePeers[remotePID] = &remotePeer{
				peerPairGroup: peerGroup,
				stream:        stream,
			}
			sp.srvcEng.Go(func(srvcCtx context.Context) {
				sp.recvLoopSingle(srvcCtx, remotePID, stream.ReceiveMessages())
			})
			sp.lggr.Infow("Created stream to remote peer", "remotePeerID", remotePID)
		}

		// Remove obsolete remote peers
		for pid := range sp.remotePeers {
			if _, ok := desiredRemotePeers[pid]; !ok {
				rp := sp.remotePeers[pid]
				if rp.peerPairGroup != nil {
					rp.peerPairGroup.Close() // closes the stream
				}
				delete(sp.remotePeers, pid)
				sp.lggr.Infow("Closed stream to remote peer", "remotePeerID", pid)
			}
		}
	}
	return nil
}

func (sp *don2DonSharedPeer) recvLoopSingle(ctx context.Context, pid ragetypes.PeerID, ch <-chan []byte) {
	sp.lggr.Infow("starting recvLoopSingle", "peerID", pid)
	for {
		select {
		case <-ctx.Done():
			sp.lggr.Infow("stopped - exiting recvLoopSingle", "peerID", pid)
			return
		case msg, ok := <-ch:
			if !ok {
				sp.lggr.Infow("channel closed - exiting recvLoopSingle", "peerID", pid)
				return
			}
			sp.recvCh <- p2ptypes.Message{Sender: pid, Payload: msg}
		}
	}
}

// Create a unique ID based of both DON IDs and all sorted DON members to uniquely represent a DON pair
func pairID(donA, donB capabilities.DON) string {
	if donA.ID > donB.ID {
		donA, donB = donB, donA
	}
	memberIDsA := make([]string, 0, len(donA.Members))
	for _, pid := range donA.Members {
		memberIDsA = append(memberIDsA, pid.String())
	}
	slices.Sort(memberIDsA)
	hashPeersA := sha256.Sum256([]byte(strings.Join(memberIDsA, ",")))
	memberIDsB := make([]string, 0, len(donB.Members))
	for _, pid := range donB.Members {
		memberIDsB = append(memberIDsB, pid.String())
	}
	slices.Sort(memberIDsB)
	hashPeersB := sha256.Sum256([]byte(strings.Join(memberIDsB, ",")))
	return fmt.Sprintf("%d-%d-%s-%s", donA.ID, donB.ID, hex.EncodeToString(hashPeersA[:]), hex.EncodeToString(hashPeersB[:]))
}

func donPairDigest(donID1, donID2 uint32) ocr2types.ConfigDigest {
	// Create a digest based on sorted DON IDs
	if donID1 > donID2 {
		donID1, donID2 = donID2, donID1
	}
	var digest ocr2types.ConfigDigest
	binary.BigEndian.PutUint16(digest[:], uint16(digestPrefix))
	binary.BigEndian.PutUint32(digest[2:], donID1)
	binary.BigEndian.PutUint32(digest[6:], donID2)
	return digest
}

func nodePairDigest(peerID1, peerID2 ragetypes.PeerID) ocr2types.ConfigDigest {
	id1Str, id2Str := peerIDStrings(peerID1, peerID2)
	var digest ocr2types.ConfigDigest
	binary.BigEndian.PutUint16(digest[:], uint16(digestPrefix))
	combinedHash := sha256.Sum256(([]byte(id1Str + id2Str)))
	copy(digest[2:], combinedHash[:30])
	return digest
}

func peerIDStrings(peerID1, peerID2 ragetypes.PeerID) (string, string) {
	id1Str := peerID1.String()
	id2Str := peerID2.String()
	if id1Str > id2Str {
		id1Str, id2Str = id2Str, id1Str
	}
	return id1Str, id2Str
}
