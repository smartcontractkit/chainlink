package networking

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	mplex "github.com/libp2p/go-libp2p-mplex"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	"github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/networking/knockingtls"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
)

const (
	dhtPrefix = "/cl_peer_discovery_dht"
)

// concretePeerV1 represents a libp2p peer with one peer ID listening on one port
type concretePeerV1 struct {
	host   p2phost.Host
	peerID p2ppeer.ID

	tls   *knockingtls.KnockingTLSTransport
	gater *connectionGater

	registrantsMu *sync.Mutex
	registrants   map[ocr1types.ConfigDigest]struct{}

	dhtAnnouncementCounterUserPrefix uint32

	// list of bandwidth limiters, one for each connection to a remote peer.
	bandwidthLimiters *knockingtls.Limiters

	logger         loghelper.LoggerWithContext
	endpointConfig EndpointConfigV1
}

type registrantV1 interface {
	getConfigDigest() ocr1types.ConfigDigest
	allower
}

func newPeerV1(c PeerConfig) (*concretePeerV1, error) {
	peerID, err := p2ppeer.IDFromPrivateKey(c.PrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting v1 peer ID from private key")
	}

	logger := loghelper.MakeRootLoggerWithContext(c.Logger).MakeChild(commontypes.LogFields{
		"id":       "PeerV1",
		"v1peerID": peerID.Pretty(),
	})

	gater, err := newConnectionGater(logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create gater")
	}

	if c.V1ListenPort == 0 {
		return nil, fmt.Errorf("V1ListenPort should not be zero")
	}

	listenAddr, err := makeMultiaddr(c.V1ListenIP, c.V1ListenPort)
	if err != nil {
		return nil, errors.Wrap(err, "could not make listen multiaddr")
	}
	logger = logger.MakeChild(commontypes.LogFields{
		"v1listenPort": c.V1ListenPort,
		"v1listenIP":   c.V1ListenIP.String(),
		"v1listenAddr": listenAddr.String(),
	})

	bandwidthLimiters := knockingtls.NewLimiters(logger)

	tlsID := knockingtls.ID
	tls, err := knockingtls.NewKnockingTLS(logger, c.PrivKey, bandwidthLimiters)
	if err != nil {
		return nil, errors.Wrap(err, "could not create knocking tls")
	}

	addrsFactory, err := makeAddrsFactory(c.V1AnnounceIP, c.V1AnnouncePort)
	if err != nil {
		return nil, errors.Wrap(err, "could not make addrs factory")
	}

	// build a custom upgrader that overrides the default secure transport with knocking TLS
	transportCon := func(upgrader *tptu.Upgrader) transport.Transport {
		betterUpgrader := tptu.Upgrader{
			upgrader.PSK,
			tls,
			upgrader.Muxer,
			upgrader.ConnGater,
		}

		return tcp.NewTCPTransport(&betterUpgrader)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrs(listenAddr),
		libp2p.Identity(c.PrivKey),
		libp2p.DisableRelay(),
		libp2p.Security(tlsID, tls),
		libp2p.ConnectionGater(gater),
		libp2p.Peerstore(c.V1Peerstore),
		libp2p.AddrsFactory(addrsFactory),
		libp2p.Transport(transportCon),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	}

	host, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	logger.Info("PeerV1: libp2p host booted", nil)

	return &concretePeerV1{
		host,
		peerID,
		tls,
		gater,
		&sync.Mutex{},
		make(map[ocr1types.ConfigDigest]struct{}),
		c.V1DHTAnnouncementCounterUserPrefix,
		bandwidthLimiters,
		logger,
		c.V1EndpointConfig,
	}, nil
}

func decodev1PeerIDs(pids []string) ([]p2ppeer.ID, error) {
	peerIDs := make([]p2ppeer.ID, len(pids))
	for i, pid := range pids {
		peerID, err := p2ppeer.Decode(pid)
		if err != nil {
			return nil, errors.Wrapf(err, "error decoding peer ID: %s", pid)
		}
		peerIDs[i] = peerID
	}
	return peerIDs, nil
}

func decodev1Bootstrappers(bootstrappers []string) (bnAddrs []p2ppeer.AddrInfo, err error) {
	bnMAddrs := make([]ma.Multiaddr, len(bootstrappers))
	for i, bNode := range bootstrappers {
		bnMAddr, err := ma.NewMultiaddr(bNode)
		if err != nil {
			return bnAddrs, errors.Wrapf(err, "could not decode peer address %s", bNode)
		}
		bnMAddrs[i] = bnMAddr
	}
	bnAddrs, err = p2ppeer.AddrInfosFromP2pAddrs(bnMAddrs...)
	if err != nil {
		return bnAddrs, errors.Wrap(err, "could not get addrinfos")
	}
	return
}

func (p1 *concretePeerV1) register(r registrantV1) error {
	configDigest := r.getConfigDigest()
	p1.registrantsMu.Lock()
	defer p1.registrantsMu.Unlock()

	p1.logger.Debug("PeerV1: registering v1 protocol handler", commontypes.LogFields{
		"configDigest": configDigest.Hex(),
	})

	if _, ok := p1.registrants[configDigest]; ok {
		p1.logger.Warn("PeerV1: Failed to register endpoint", commontypes.LogFields{"configDigest": configDigest.Hex()})
		return errors.Errorf("v1 endpoint with config digest %s has already been registered", configDigest.Hex())
	}
	p1.registrants[configDigest] = struct{}{}
	p1.gater.add(r)
	p1.tls.UpdateAllowlist(p1.gater.allowlist())
	return nil
}

func (p1 *concretePeerV1) deregister(r registrantV1) error {
	configDigest := r.getConfigDigest()
	p1.registrantsMu.Lock()
	defer p1.registrantsMu.Unlock()
	p1.logger.Debug("PeerV1: deregistering v1 protocol handler", commontypes.LogFields{
		"ProtocolID": configDigest.Hex(),
	})

	if _, ok := p1.registrants[configDigest]; !ok {
		return errors.Errorf("v1 endpoint with config digest %s is not currently registered", configDigest.Hex())
	}
	delete(p1.registrants, configDigest)

	p1.gater.remove(r)
	p1.tls.UpdateAllowlist(p1.gater.allowlist())
	return nil
}

func (p1 *concretePeerV1) PeerID() string {
	return p1.peerID.String()
}

func (p1 *concretePeerV1) ID() p2ppeer.ID {
	return p1.peerID
}

func (p1 *concretePeerV1) Close() error {
	return p1.host.Close()
}

func (p1 *concretePeerV1) newEndpoint(
	configDigest ocr1types.ConfigDigest,
	v1peerIDs []string,
	v1bootstrappers []string,
	f int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if f <= 0 {
		return nil, errors.New("can't set F to 0 or smaller")
	}

	if len(v1bootstrappers) < 1 {
		return nil, errors.New("requires at least one v1 bootstrapper")
	}

	decodedv1PeerIDs, err := decodev1PeerIDs(v1peerIDs)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v1 peer IDs")
	}

	decodedv1Bootstrappers, err := decodev1Bootstrappers(v1bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v1 bootstrappers")
	}

	return newOCREndpointV1(
		p1.logger,
		configDigest,
		p1,
		decodedv1PeerIDs,
		decodedv1Bootstrappers,
		p1.endpointConfig,
		f,
		limits,
	)
}

func (p1 *concretePeerV1) newBootstrapper(
	configDigest ocr1types.ConfigDigest,
	v1peerIDs []string,
	v1bootstrappers []string,
	f int,
) (commontypes.Bootstrapper, error) {
	if f <= 0 {
		return nil, errors.New("can't set f to zero or smaller")
	}
	decodedv1PeerIDs, err := decodev1PeerIDs(v1peerIDs)
	if err != nil {
		return nil, err
	}

	decodedv1Bootstrappers, err := decodev1Bootstrappers(v1bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v1 bootstrappers")
	}

	return newBootstrapperV1(p1.logger, configDigest, p1, decodedv1PeerIDs, decodedv1Bootstrappers, f)
}
