package ragedisco

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/libocr/internal/loghelper"
	nettypes "github.com/smartcontractkit/libocr/networking/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type ragep2pDiscovererState int

const (
	_ ragep2pDiscovererState = iota
	ragep2pDiscovererUnstarted
	ragep2pDiscovererStarted
	ragep2pDiscovererClosed
)

type Ragep2pDiscoverer struct {
	logger            loghelper.LoggerWithContext
	proc              subprocesses.Subprocesses
	ctx               context.Context
	ctxCancel         context.CancelFunc
	deltaReconcile    time.Duration
	announceAddresses []string
	db                nettypes.DiscovererDatabase
	host              *ragep2p.Host
	proto             *discoveryProtocol

	stateMu sync.Mutex
	state   ragep2pDiscovererState

	streamsMu sync.Mutex
	streams   map[ragetypes.PeerID]*ragep2p.Stream

	chIncomingMessages chan incomingMessage
	chOutgoingMessages chan outgoingMessage
	chConnectivity     chan connectivityMsg
}

func NewRagep2pDiscoverer(
	deltaReconcile time.Duration,
	announceAddresses []string,
	db nettypes.DiscovererDatabase,
) *Ragep2pDiscoverer {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &Ragep2pDiscoverer{
		nil, // logger, filled on Start()
		subprocesses.Subprocesses{},
		ctx,
		ctxCancel,
		deltaReconcile,
		announceAddresses,
		db,
		nil, // ragep2p host, filled on Start()
		nil, // discovery protocol, filled on Start()
		sync.Mutex{},
		ragep2pDiscovererUnstarted,
		sync.Mutex{},
		make(map[ragetypes.PeerID]*ragep2p.Stream),
		make(chan incomingMessage),
		make(chan outgoingMessage),
		make(chan connectivityMsg),
	}
}

func (r *Ragep2pDiscoverer) Start(h *ragep2p.Host, privKey ed25519.PrivateKey, logger loghelper.LoggerWithContext) error {
	succeeded := false
	defer func() {
		if !succeeded {
			r.Close()
		}
	}()

	r.logger = logger
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	if r.state != ragep2pDiscovererUnstarted {
		return fmt.Errorf("cannot start Ragep2pDiscoverer that is not unstarted, state was: %v", r.state)
	}
	r.state = ragep2pDiscovererStarted
	r.host = h
	announceAddresses, ok := combinedAnnounceAddrsForDiscoverer(r.logger, r.announceAddresses)
	if !ok {
		return fmt.Errorf("failed to obtain announce addresses")
	}
	proto, err := newDiscoveryProtocol(
		r.deltaReconcile,
		r.chIncomingMessages,
		r.chOutgoingMessages,
		r.chConnectivity,
		privKey,
		announceAddresses,
		r.db,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to construct underlying discovery protocol: %w", err)
	}
	r.proto = proto
	err = r.proto.Start()
	if err != nil {
		return fmt.Errorf("failed to start underlying discovery protocol: %w", err)
	}
	r.proc.Go(r.connectivityLoop)
	r.proc.Go(r.writeLoop)

	succeeded = true
	return nil
}

func (r *Ragep2pDiscoverer) connectivityLoop() {
	var subs subprocesses.Subprocesses
	defer subs.Wait()
	for {
		select {
		case c := <-r.chConnectivity:
			logger := r.logger.MakeChild(commontypes.LogFields{
				"remotePeerID": c.peerID,
			})
			if c.peerID == r.host.ID() {
				break
			}
			r.streamsMu.Lock()
			if c.msgType == connectivityAdd {
				if _, exists := r.streams[c.peerID]; exists {
					r.streamsMu.Unlock()
					break
				}
				// no point in keeping very large buffers, since only
				// the latest messages matter anyways.
				bufferSize := 2
				messagesLimit := ragep2p.TokenBucketParams{
					// we expect one message every deltaReconcile seconds, let's double it
					// for good measure
					2 / r.deltaReconcile.Seconds(),
					// twice the buffer size should be plenty
					2 * uint32(bufferSize),
				}
				// bytesLimit is messagesLimit * maxMessageLength
				bytesLimit := ragep2p.TokenBucketParams{
					messagesLimit.Rate * maxMessageLength,
					messagesLimit.Capacity * maxMessageLength,
				}
				s, err := r.host.NewStream(
					c.peerID,
					"ragedisco/v1",
					bufferSize,
					bufferSize,
					maxMessageLength,
					messagesLimit,
					bytesLimit,
				)
				if err != nil {
					logger.Warn("NewStream failed!", reason(err))
					r.streamsMu.Unlock()
					break
				}
				r.streams[c.peerID] = s
				r.streamsMu.Unlock()
				pid := c.peerID
				subs.Go(func() {
					chDone := r.ctx.Done()
					for {
						select {
						case m, ok := <-s.ReceiveMessages():
							if !ok { // stream Close() will signal us when it's time to go
								return
							}
							w, err := fromProtoWrappedBytes(m)
							if err != nil {
								logger.Warn("Failed to unwrap incoming message", reason(err))
								break
							}
							select {
							case r.chIncomingMessages <- incomingMessage{w, pid}:
							case <-chDone:
								return
							}
						case <-chDone:
							return
						}
					}
				})
			} else {
				if _, exists := r.streams[c.peerID]; !exists {
					logger.Warn("Asked to remove connectivity with peer we don't have a stream for", nil)
					r.streamsMu.Unlock()
					break
				}
				if err := r.streams[c.peerID].Close(); err != nil {
					logger.Warn("Failed to close stream", reason(err))
				}
				delete(r.streams, c.peerID)
				r.streamsMu.Unlock()
			}
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Ragep2pDiscoverer) writeLoop() {
	for {
		select {
		case m := <-r.chOutgoingMessages:
			r.streamsMu.Lock()
			s, exists := r.streams[m.to]
			if !exists {
				r.logger.Warn("Write message to peer we don't have a stream open for", commontypes.LogFields{
					"remotePeerID": m.to,
				})
				r.streamsMu.Unlock()
				break
			}
			r.streamsMu.Unlock()
			bs, err := toBytesWrapped(m.payload)
			if err != nil {
				r.logger.Warn("Failed to convert message to bytes", commontypes.LogFields{"message": m.payload})
				break
			}
			s.SendMessage(bs)
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Ragep2pDiscoverer) Close() error {
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	if r.state != ragep2pDiscovererStarted {
		return fmt.Errorf("cannot close Ragep2pDiscoverer that is not started, state was: %v", r.state)
	}
	r.state = ragep2pDiscovererClosed

	r.ctxCancel()
	r.proc.Wait()
	if r.proto != nil {
		return r.proto.Close()
	}
	return nil
}

func (r *Ragep2pDiscoverer) AddGroup(digest types.ConfigDigest, onodes []ragetypes.PeerID, bnodes []ragetypes.PeerInfo) error {
	r.logger.Info("Ragep2pDiscoverer: Adding group", commontypes.LogFields{
		"configDigest": digest,
		"oracles":      onodes,
		"bootstraps":   bnodes,
	})
	return r.proto.addGroup(digest, onodes, bnodes)
}

// RemoveGroup should not block or panic even if the discoverer is closed.
func (r *Ragep2pDiscoverer) RemoveGroup(digest types.ConfigDigest) error {
	r.logger.Info("Ragep2pDiscoverer: Removing group", commontypes.LogFields{"configDigest": digest})
	return r.proto.removeGroup(digest)
}

func (r *Ragep2pDiscoverer) FindPeer(peer ragetypes.PeerID) ([]ragetypes.Address, error) {
	return r.proto.FindPeer(peer)
}

var _ ragep2p.Discoverer = &Ragep2pDiscoverer{}
