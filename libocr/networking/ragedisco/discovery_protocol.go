package ragedisco

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/commontypes"
	nettypes "github.com/smartcontractkit/libocr/networking/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type incomingMessage struct {
	payload WrappableMessage
	from    ragetypes.PeerID
}

type outgoingMessage struct {
	payload WrappableMessage
	to      ragetypes.PeerID
}

type discoveryProtocolState int

const (
	_ discoveryProtocolState = iota
	discoveryProtocolUnstarted
	discoveryProtocolStarted
	discoveryProtocolClosed
)

// discoveryProtocolLocked contains a subset of a discoveryProtocol's state
// that requires the discoveryProtocol lock to be held in order to access or
// modify
type discoveryProtocolLocked struct {
	bestAnnouncement        map[ragetypes.PeerID]Announcement
	groups                  map[types.ConfigDigest]*group
	bootstrappers           map[ragetypes.PeerID]map[ragetypes.Address]int
	numGroupsByOracle       map[ragetypes.PeerID]int
	numGroupsByBootstrapper map[ragetypes.PeerID]int
}

type discoveryProtocol struct {
	stateMu sync.Mutex
	state   discoveryProtocolState

	deltaReconcile     time.Duration
	chIncomingMessages <-chan incomingMessage
	chOutgoingMessages chan<- outgoingMessage
	chConnectivity     chan<- connectivityMsg
	chInternalBump     chan Announcement
	privKey            ed25519.PrivateKey
	ownID              ragetypes.PeerID
	ownAddrs           []ragetypes.Address

	lock   sync.RWMutex
	locked discoveryProtocolLocked

	db nettypes.DiscovererDatabase

	processes subprocesses.Subprocesses
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    loghelper.LoggerWithContext
}

const (
	announcementVersionWarnThreshold = 100e6

	saveInterval       = 2 * time.Minute
	reportInitialDelay = 10 * time.Second
	reportInterval     = 5 * time.Minute
)

func newDiscoveryProtocol(
	deltaReconcile time.Duration,
	chIncomingMessages <-chan incomingMessage,
	chOutgoingMessages chan<- outgoingMessage,
	chConnectivity chan<- connectivityMsg,
	privKey ed25519.PrivateKey,
	ownAddrs []ragetypes.Address,
	db nettypes.DiscovererDatabase,
	logger loghelper.LoggerWithContext,
) (*discoveryProtocol, error) {
	ownID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain peer id from private key: %w", err)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	return &discoveryProtocol{
		sync.Mutex{},
		discoveryProtocolUnstarted,
		deltaReconcile,
		chIncomingMessages,
		chOutgoingMessages,
		chConnectivity,
		make(chan Announcement),
		privKey,
		ownID,
		ownAddrs,
		sync.RWMutex{},
		discoveryProtocolLocked{
			make(map[ragetypes.PeerID]Announcement),
			make(map[types.ConfigDigest]*group),
			make(map[ragetypes.PeerID]map[ragetypes.Address]int),
			make(map[ragetypes.PeerID]int),
			make(map[ragetypes.PeerID]int),
		},
		db,
		subprocesses.Subprocesses{},
		ctx,
		ctxCancel,
		logger.MakeChild(commontypes.LogFields{"id": "discoveryProtocol"}),
	}, nil
}

func (p *discoveryProtocol) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			p.Close()
		}
	}()

	p.stateMu.Lock()
	defer p.stateMu.Unlock()
	if p.state != discoveryProtocolUnstarted {
		return fmt.Errorf("cannot start discoveryProtocol that is not unstarted, state was: %v", p.state)
	}
	p.state = discoveryProtocolStarted

	p.lock.Lock()
	defer p.lock.Unlock()
	_, _, err := p.lockedBumpOwnAnnouncement()
	if err != nil {
		return fmt.Errorf("failed to bump own announcement: %w", err)
	}
	p.processes.Go(p.recvLoop)
	p.processes.Go(p.sendLoop)
	p.processes.Go(p.saveLoop)
	p.processes.Go(p.statusReportLoop)
	succeeded = true
	return nil
}

func formatAnnouncementsForReport(allIDs map[ragetypes.PeerID]struct{}, baSigned map[ragetypes.PeerID]Announcement) (string, int) {
	maybeAnnouncementById := make(map[ragetypes.PeerID][]unsignedAnnouncement)
	// For every peer, its array will length 0 when they are undetected or
	// length 1 when they are detected. Note that we can't use pointers instead
	// because printing does not dereference the values.
	// Example for peers A, B where we have an detected B but not A:
	// map[A:[] B:[{Addrs:[1.2.3.4:1234] Counter:0}]]
	undetected := 0
	for id := range allIDs {
		ann, exists := baSigned[id]
		if exists {
			maybeAnnouncementById[id] = append(maybeAnnouncementById[id], ann.unsignedAnnouncement)
		} else {
			maybeAnnouncementById[id] = nil
			undetected++
		}
	}
	return fmt.Sprintf("%+v", maybeAnnouncementById), undetected
}

func (p *discoveryProtocol) statusReportLoop() {
	chDone := p.ctx.Done()
	timer := time.After(reportInitialDelay)
	for {
		select {
		case <-timer:
			func() {
				p.lock.RLock()
				defer p.lock.RUnlock()
				uniquePeersToDetect := make(map[ragetypes.PeerID]struct{})
				for id, cnt := range p.locked.numGroupsByOracle {
					if cnt == 0 {
						continue
					}
					uniquePeersToDetect[id] = struct{}{}
				}

				reportStr, undetected := formatAnnouncementsForReport(uniquePeersToDetect, p.locked.bestAnnouncement)
				p.logger.Info("DiscoveryProtocol: Status report", commontypes.LogFields{
					"statusByPeer":    reportStr,
					"peersToDetect":   len(uniquePeersToDetect),
					"peersUndetected": undetected,
					"peersDetected":   len(uniquePeersToDetect) - undetected,
				})
				timer = time.After(reportInterval)
			}()
		case <-chDone:
			return
		}
	}
}

// Peer A is allowed to learn about an Announcement by peer B if B is an oracle node in
// one of the groups A participates in.
func (p *discoveryProtocol) lockedAllowedPeers(ann Announcement) (ps []ragetypes.PeerID) {
	annPeerID, err := ann.PeerID()
	if err != nil {
		p.logger.Warn("Failed to obtain peer id from announcement", reason(err))
		return
	}
	peers := make(map[ragetypes.PeerID]struct{})
	for _, g := range p.locked.groups {
		if !g.hasOracle(annPeerID) {
			continue
		}
		for _, pid := range g.peerIDs() {
			peers[pid] = struct{}{}
		}
	}
	for pid := range peers {
		if pid == p.ownID {
			continue
		}
		ps = append(ps, pid)
	}
	return
}

func (p *discoveryProtocol) addGroup(digest types.ConfigDigest, onodes []ragetypes.PeerID, bnodes []ragetypes.PeerInfo) error {
	var newPeerIDs []ragetypes.PeerID
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, exists := p.locked.groups[digest]; exists {
		return fmt.Errorf("asked to add group with digest we already have (digest: %s)", digest.Hex())
	}
	newGroup := group{oracleNodes: onodes, bootstrapperNodes: bnodes}
	p.locked.groups[digest] = &newGroup
	for _, oid := range onodes {
		if p.locked.numGroupsByOracle[oid] == 0 {
			newPeerIDs = append(newPeerIDs, oid)
		}
		p.locked.numGroupsByOracle[oid]++
	}
	for _, bs := range bnodes {
		p.locked.numGroupsByBootstrapper[bs.ID]++
		for _, addr := range bs.Addrs {
			if _, exists := p.locked.bootstrappers[bs.ID]; !exists {
				p.locked.bootstrappers[bs.ID] = make(map[ragetypes.Address]int)
			}
			p.locked.bootstrappers[bs.ID][addr]++
		}
	}
	for _, pid := range newGroup.peerIDs() {
		// it's ok to send connectivityAdd messages multiple times
		select {
		case p.chConnectivity <- connectivityMsg{connectivityAdd, pid}:
		case <-p.ctx.Done():
			return nil
		}
	}

	// we hold lock here
	if err := p.lockedLoadFromDB(newPeerIDs); err != nil {
		// db-level errors are not prohibitive
		p.logger.Warn("DiscoveryProtocol: Failed to load announcements from db", commontypes.LogFields{"configDigest": digest, "error": err})
	}
	return nil
}

func (p *discoveryProtocol) lockedLoadFromDB(ragePeerIDs []ragetypes.PeerID) error {
	// The database may have been set to nil, and we don't necessarily need it to function.
	if len(ragePeerIDs) == 0 || p.db == nil {
		return nil
	}
	p.logger.Info("Loading announcements from db", commontypes.LogFields{"peerIDs": ragePeerIDs})
	strPeerIDs := make([]string, len(ragePeerIDs))
	for i, pid := range ragePeerIDs {
		strPeerIDs[i] = pid.String()
	}
	annByID, err := p.db.ReadAnnouncements(p.ctx, strPeerIDs)
	if err != nil {
		return err
	}
	var loaded, found []string
	for peerID, dbannBytes := range annByID {
		found = append(found, peerID)
		dbann, err := deserializeSignedAnnouncement(dbannBytes)
		if err != nil {
			p.logger.Error("Failed to deserialize signed announcement from db", commontypes.LogFields{
				"announcementPeerID": peerID,
				"announcementBytes":  hex.EncodeToString(dbannBytes),
				"error":              err,
			})
			continue
		}
		if err := p.lockedProcessAnnouncement(dbann); err != nil {
			p.logger.Error("Failed to process announcement from db", commontypes.LogFields{
				"announcement": dbann,
				"error":        err,
			})
			continue
		}
		loaded = append(loaded, peerID)
	}
	p.logger.Info("Loaded announcements from db", commontypes.LogFields{
		"queried":    ragePeerIDs,
		"numQueried": len(ragePeerIDs),
		"found":      found,
		"numFound":   len(found),
		"loaded":     loaded,
		"numLoaded":  len(loaded),
	})
	return nil
}

func (p *discoveryProtocol) saveAnnouncementToDB(ann Announcement) error {
	if p.db == nil {
		return nil
	}
	ser, err := ann.serialize()
	if err != nil {
		return err
	}
	pid, err := ann.PeerID()
	if err != nil {
		return err
	}
	return p.db.StoreAnnouncement(p.ctx, pid.String(), ser)
}

func (p *discoveryProtocol) saveToDB() error {
	if p.db == nil {
		return nil
	}
	p.lock.RLock()
	defer p.lock.RUnlock()

	var allErrors error
	for _, ann := range p.locked.bestAnnouncement {
		allErrors = multierr.Append(allErrors, p.saveAnnouncementToDB(ann))
	}
	return allErrors
}

func (p *discoveryProtocol) saveLoop() {
	if p.db == nil {
		return
	}
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "saveLoop"})
	logger.Debug("Entering", nil)
	defer logger.Debug("Exiting", nil)
	for {
		select {
		case <-time.After(saveInterval):
		case <-p.ctx.Done():
			return
		}

		if err := p.saveToDB(); err != nil {
			logger.Warn("Failed to save announcements to db", reason(err))
		}
	}
}

func (p *discoveryProtocol) removeGroup(digest types.ConfigDigest) error {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "removeGroup"})
	logger.Trace("Called", nil)
	p.lock.Lock()
	defer p.lock.Unlock()

	goneGroup, exists := p.locked.groups[digest]
	if !exists {
		return fmt.Errorf("can't remove group that is not registered (digest: %s)", digest.Hex())
	}

	delete(p.locked.groups, digest)

	for _, oid := range goneGroup.oracleIDs() {
		p.locked.numGroupsByOracle[oid]--
		if p.locked.numGroupsByOracle[oid] == 0 {
			if ann, exists := p.locked.bestAnnouncement[oid]; exists {
				if err := p.saveAnnouncementToDB(ann); err != nil {
					p.logger.Warn("Failed to save announcement from removed group to DB", reason(err))
				}
			}
			if oid != p.ownID {
				delete(p.locked.bestAnnouncement, oid)
			}
			delete(p.locked.numGroupsByOracle, oid)
		}
	}

	for _, binfo := range goneGroup.bootstrapperNodes {
		bid := binfo.ID

		p.locked.numGroupsByBootstrapper[bid]--
		if p.locked.numGroupsByBootstrapper[bid] == 0 {
			delete(p.locked.numGroupsByBootstrapper, bid)
			delete(p.locked.bootstrappers, bid)
			continue
		}
		for _, addr := range binfo.Addrs {
			p.locked.bootstrappers[bid][addr]--
			if p.locked.bootstrappers[bid][addr] == 0 {
				delete(p.locked.bootstrappers[bid], addr)
			}
		}
	}

	// Cleanup connections for peers we don't have in any group anymore.
	for _, pid := range goneGroup.peerIDs() {
		if p.locked.numGroupsByOracle[pid]+p.locked.numGroupsByBootstrapper[pid] == 0 {
			select {
			case p.chConnectivity <- connectivityMsg{connectivityRemove, pid}:
			case <-p.ctx.Done():
				return nil
			}
		}
	}

	return nil
}

func (p *discoveryProtocol) FindPeer(peer ragetypes.PeerID) ([]ragetypes.Address, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	var addrs []ragetypes.Address
	// The addresses we know from local configuration take priority — useful for overriding addresses in disaster
	// scenarios
	if baddrs, ok := p.locked.bootstrappers[peer]; ok {
		for baddr := range baddrs {
			addrs = append(addrs, baddr)
		}
	}
	// Followed by the addresses obtained by the best announcement
	if ann, ok := p.locked.bestAnnouncement[peer]; ok {
		addrs = append(addrs, ann.Addrs...)
	}
	return dedup(addrs), nil
}

func (p *discoveryProtocol) recvLoop() {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "recvLoop"})
	logger.Debug("Entering", nil)
	defer logger.Debug("Exiting", nil)
	for {
		select {
		case <-p.ctx.Done():
			return
		case msg := <-p.chIncomingMessages:
			logger := logger.MakeChild(commontypes.LogFields{"remotePeerID": msg.from})
			switch v := msg.payload.(type) {
			case *Announcement:
				announcement := *v
				logger.Trace("Received announcement", commontypes.LogFields{"announcement": announcement})
				if err := p.processAnnouncement(announcement); err != nil {
					logger.Warn("Failed to process announcement", commontypes.LogFields{
						"announcement": announcement,
						"error":        err,
					})
				}
			case *reconcile:
				reconcile := *v
				// logger.Trace("Received reconcile", commontypes.LogFields{"reconcile": reconcile})
				for _, ann := range reconcile.Anns {
					if err := p.processAnnouncement(ann); err != nil {

						logger.Warn("Failed to process announcement from reconcile", commontypes.LogFields{
							"reconcile":    reconcile,
							"announcement": ann,
							"error":        err,
						})
					}
				}
			default:
				logger.Warn("Received unknown message type", commontypes.LogFields{"msg": v})
			}
		}
	}
}

// processAnnouncement locks lock for its whole lifetime.
func (p *discoveryProtocol) processAnnouncement(ann Announcement) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.lockedProcessAnnouncement(ann)
}

// lockedProcessAnnouncement requires lock to be held.
func (p *discoveryProtocol) lockedProcessAnnouncement(ann Announcement) error {
	logger := p.logger.MakeChild(commontypes.LogFields{
		"in":           "processAnnouncement",
		"announcement": ann,
	})
	pid, err := ann.PeerID()
	if err != nil {
		return fmt.Errorf("failed to obtain peer id: %w", err)
	}

	if p.locked.numGroupsByOracle[pid] == 0 {
		return fmt.Errorf("peer %s is not an oracle in any of our jobs; perhaps whoever sent this is running a job that includes us and this peer, but we are not running that job", pid)
	}

	err = ann.verify()
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	if localann, exists := p.locked.bestAnnouncement[pid]; !exists || localann.Counter <= ann.Counter {
		if exists && pid != p.ownID && localann.Counter == ann.Counter {
			return nil
		}
		p.locked.bestAnnouncement[pid] = ann
		if pid == p.ownID {
			bumpedann, better, err := p.lockedBumpOwnAnnouncement()
			if err != nil {
				return fmt.Errorf("failed to bump own announcement: %w", err)
			}

			if !better {
				return nil
			}

			logger.Info("Received better announcement for us - bumped", nil)
			select {
			case p.chInternalBump <- *bumpedann:
			case <-p.ctx.Done():
				return nil
			}
		} else {
			logger.Info("Received better announcement for peer", nil)
			select {
			case p.chConnectivity <- connectivityMsg{connectivityAdd, pid}:
			case <-p.ctx.Done():
				return nil
			}
		}
	}

	return nil
}

func (p *discoveryProtocol) sendToAllowedPeers(ann Announcement) {
	p.lock.RLock()
	allowedPeers := p.lockedAllowedPeers(ann)
	p.lock.RUnlock()
	for _, pid := range allowedPeers {
		select {
		case p.chOutgoingMessages <- outgoingMessage{ann, pid}:
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *discoveryProtocol) sendLoop() {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "sendLoop"})
	logger.Debug("Entering", nil)
	defer logger.Debug("Exiting", nil)
	tick := time.After(0)
	for {
		select {
		case <-p.ctx.Done():
			return
		case ourann := <-p.chInternalBump:
			logger.Info("Our announcement was bumped - broadcasting", commontypes.LogFields{"announcement": ourann})
			p.sendToAllowedPeers(ourann)
		case <-tick:
			logger.Debug("Starting reconciliation", nil)
			reconcileByPeer := make(map[ragetypes.PeerID]*reconcile)
			func() {
				p.lock.RLock()
				defer p.lock.RUnlock()
				for _, ann := range p.locked.bestAnnouncement {
					for _, pid := range p.lockedAllowedPeers(ann) {
						if _, exists := reconcileByPeer[pid]; !exists {
							reconcileByPeer[pid] = &reconcile{Anns: []Announcement{}}
						}
						r := reconcileByPeer[pid]
						r.Anns = append(r.Anns, ann)
					}
				}
			}()

			for pid, rec := range reconcileByPeer {
				select {
				case p.chOutgoingMessages <- outgoingMessage{rec, pid}:
					logger.Trace("Sending reconcile", commontypes.LogFields{"remotePeerID": pid, "reconcile": rec})
				case <-p.ctx.Done():
					return
				}
			}
			tick = time.After(p.deltaReconcile)
		}
	}
}

// lockedBumpOwnAnnouncement requires lock to be held by the caller.
func (p *discoveryProtocol) lockedBumpOwnAnnouncement() (*Announcement, bool, error) {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "lockedBumpOwnAnnouncement"})
	oldann, exists := p.locked.bestAnnouncement[p.ownID]
	newctr := uint64(0)

	if exists {
		if equalAddrs(oldann.Addrs, p.ownAddrs) {
			return nil, false, nil
		}
		// Counter is uint64, and it only changes when a peer's
		// addresses change. We assume a peer will not change addresses
		// more than 2**64 times.
		newctr = oldann.Counter + 1
	}
	newann := unsignedAnnouncement{Addrs: p.ownAddrs, Counter: newctr}
	if newctr > announcementVersionWarnThreshold {
		logger.Warn("New announcement version too big!", commontypes.LogFields{"announcement": newann})
	}
	sann, err := newann.sign(p.privKey)
	if err != nil {
		return nil, false, fmt.Errorf("failed to sign own announcement: %w", err)
	}
	logger.Info("DiscoveryProtocol: Replacing our own announcement", commontypes.LogFields{"announcement": sann})
	p.locked.bestAnnouncement[p.ownID] = sann
	return &sann, true, nil
}

func (p *discoveryProtocol) Close() error {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "Close"})
	p.stateMu.Lock()
	defer p.stateMu.Unlock()
	if p.state != discoveryProtocolStarted {
		return fmt.Errorf("cannot close discoveryProtocol that is not started, state was: %v", p.state)
	}
	p.state = discoveryProtocolClosed

	logger.Debug("Exiting", nil)
	defer logger.Debug("Exited", nil)
	p.ctxCancel()
	p.processes.Wait()
	return nil
}
