package p2p

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/networking/ragedisco"
	nettypes "github.com/smartcontractkit/libocr/networking/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var (
	defaultStreamName = "stream"
	defaultRecvChSize = 10000
)

type PeerConfig struct {
	PrivateKey ed25519.PrivateKey
	// List of <ip>:<port> addresses.
	ListenAddresses []string
	// List of <host>:<port> addresses. If empty, defaults to ListenAddresses.
	AnnounceAddresses []string
	Bootstrappers     []ragetypes.PeerInfo
	// Every DeltaReconcile a Reconcile message is sent to every peer.
	DeltaReconcile time.Duration
	// Dial attempts will be at least DeltaDial apart.
	DeltaDial          time.Duration
	DiscovererDatabase nettypes.DiscovererDatabase
	MetricsRegisterer  prometheus.Registerer
}

type peer struct {
	streams     map[ragetypes.PeerID]*ragep2p.Stream
	cfg         PeerConfig
	isBootstrap bool
	host        *ragep2p.Host
	discoverer  *ragedisco.Ragep2pDiscoverer
	myID        ragetypes.PeerID
	recvCh      chan p2ptypes.Message

	stopCh  services.StopChan
	wg      sync.WaitGroup
	lggr    logger.Logger
	groupID *counter
}

var _ p2ptypes.Peer = &peer{}

func NewPeer(cfg PeerConfig, lggr logger.Logger) (*peer, error) {
	peerID, err := ragetypes.PeerIDFromPrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error extracting v2 peer ID from private key: %w", err)
	}
	isBootstrap := false
	for _, info := range cfg.Bootstrappers {
		if info.ID == peerID {
			isBootstrap = true
			break
		}
	}

	announceAddresses := cfg.AnnounceAddresses
	if len(cfg.AnnounceAddresses) == 0 {
		announceAddresses = cfg.ListenAddresses
	}

	discoverer := ragedisco.NewRagep2pDiscoverer(cfg.DeltaReconcile, announceAddresses, cfg.DiscovererDatabase, cfg.MetricsRegisterer)
	commonLggr := commonlogger.NewOCRWrapper(lggr, true, func(string) {})

	host, err := ragep2p.NewHost(
		ragep2p.HostConfig{DurationBetweenDials: cfg.DeltaDial},
		cfg.PrivateKey,
		cfg.ListenAddresses,
		discoverer,
		commonLggr,
		cfg.MetricsRegisterer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct ragep2p host: %w", err)
	}

	return &peer{
		streams:     make(map[ragetypes.PeerID]*ragep2p.Stream),
		cfg:         cfg,
		isBootstrap: isBootstrap,
		host:        host,
		discoverer:  discoverer,
		myID:        peerID,
		recvCh:      make(chan p2ptypes.Message, defaultRecvChSize),
		stopCh:      make(services.StopChan),
		lggr:        lggr.Named("P2PPeer"),
		groupID:     &counter{},
	}, nil
}

func (p *peer) ID() ragetypes.PeerID {
	return p.myID
}

func (p *peer) UpdateConnections(peers map[ragetypes.PeerID]p2ptypes.StreamConfig) error {
	p.lggr.Infow("updating peer addresses", "peers", peers)
	if !p.isBootstrap {
		if err := p.recreateStreams(peers); err != nil {
			return err
		}
	}
	// updating the group is a small optimization that avoids reconnecting to existing peers
	currentGroupID := p.groupID.Bytes()
	newGroupID := p.groupID.Inc().Bytes()
	peerIDs := []ragetypes.PeerID{}
	for pid := range peers {
		peerIDs = append(peerIDs, pid)
	}
	if err := p.discoverer.AddGroup(newGroupID, peerIDs, p.cfg.Bootstrappers); err != nil {
		p.lggr.Warnw("failed to add group", "groupID", newGroupID)
		return err
	}
	if err := p.discoverer.RemoveGroup(currentGroupID); err != nil {
		p.lggr.Warnw("failed to remove old group", "groupID", currentGroupID)
	}

	return nil
}

func (p *peer) recreateStreams(peers map[ragetypes.PeerID]p2ptypes.StreamConfig) error {
	for pid, cfg := range peers {
		pid := pid
		if pid == p.myID { // don't create a self-stream
			continue
		}
		_, ok := p.streams[pid]
		if ok { // already have a stream with this peer
			continue
		}

		stream, err := p.host.NewStream(
			pid,
			defaultStreamName,
			cfg.OutgoingMessageBufferSize,
			cfg.IncomingMessageBufferSize,
			cfg.MaxMessageLenBytes,
			cfg.MessageRateLimiter,
			cfg.BytesRateLimiter,
		)
		if err != nil {
			return fmt.Errorf("failed to create stream to peer id: %q: %w", pid, err)
		}
		p.lggr.Infow("adding peer", "peerID", pid)
		p.streams[pid] = stream
		p.wg.Add(1)
		go p.recvLoopSingle(pid, stream.ReceiveMessages())
	}
	// remove obsolete streams
	for pid, stream := range p.streams {
		_, ok := peers[pid]
		if !ok {
			p.lggr.Infow("removing peer", "peerID", pid)
			delete(p.streams, pid)
			err := stream.Close()
			if err != nil {
				p.lggr.Errorw("failed to close stream", "peerID", pid, "error", err)
			}
		}
	}
	return nil
}

func (p *peer) Start(ctx context.Context) error {
	err := p.host.Start()
	if err != nil {
		return fmt.Errorf("failed to start ragep2p host: %w", err)
	}
	p.lggr.Info("peer started")
	return nil
}

func (p *peer) recvLoopSingle(pid ragetypes.PeerID, ch <-chan []byte) {
	p.lggr.Infow("starting recvLoopSingle", "peerID", pid)
	defer p.wg.Done()
	for {
		select {
		case <-p.stopCh:
			p.lggr.Infow("stopped - exiting recvLoopSingle", "peerID", pid)
			return
		case msg, ok := <-ch:
			if !ok {
				p.lggr.Infow("channel closed - exiting recvLoopSingle", "peerID", pid)
				return
			}
			p.recvCh <- p2ptypes.Message{Sender: pid, Payload: msg}
		}
	}
}

func (p *peer) Send(peerID ragetypes.PeerID, msg []byte) error {
	stream, ok := p.streams[peerID]
	if !ok {
		return fmt.Errorf("no stream to peer id: %q", peerID)
	}
	stream.SendMessage(msg)
	return nil
}

func (p *peer) Receive() <-chan p2ptypes.Message {
	return p.recvCh
}

func (p *peer) Close() error {
	err := p.host.Close()
	close(p.stopCh)
	p.wg.Wait()
	p.lggr.Info("peer closed")
	return err
}

func (p *peer) Ready() error {
	return nil
}

func (p *peer) HealthReport() map[string]error {
	return nil
}

func (p *peer) Name() string {
	return "P2PPeer"
}
