package networking

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/smartcontractkit/chainlink/libocr/networking/knockingtls"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"

	"github.com/libp2p/go-libp2p"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

var (
	_ types.BinaryNetworkEndpointFactory = &concretePeer{}
	_ types.BootstrapperFactory          = &concretePeer{}
)

const (
	dhtPrefix = "/cl_peer_discovery_dht"
)

type PeerConfig struct {
	PrivKey        p2pcrypto.PrivKey
	ListenPort     uint16
	ListenIP       net.IP
	Logger         types.Logger
	Peerstore      p2ppeerstore.Peerstore
	EndpointConfig EndpointConfig
}

type concretePeer struct {
	p2phost.Host
	tls            *knockingtls.KnockingTLSTransport
	gater          *connectionGater
	logger         types.Logger
	endpointConfig EndpointConfig
	registrants    map[types.ConfigDigest]struct{}
	registrantsMu  *sync.Mutex
}

type registrant interface {
	allower
	getConfigDigest() types.ConfigDigest
}

func NewPeer(c PeerConfig) (*concretePeer, error) {
	if c.ListenPort == 0 {
		return nil, errors.New("NewPeer requires a non-zero listen port")
	}

	peerID, err := p2ppeer.IDFromPrivateKey(c.PrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting peer ID from private key")
	}

	ip4 := c.ListenIP.To4()
	if ip4 == nil {
		return nil, errors.Errorf("listen address must be a valid ipv4 address, got: %s", c.ListenIP.String())
	}
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", ip4.String(), c.ListenPort)

	logger := loghelper.MakeLoggerWithContext(c.Logger, types.LogFields{
		"id":         "Peer",
		"peerID":     peerID.Pretty(),
		"listenPort": c.ListenPort,
		"listenAddr": listenAddr,
	})

	gater, err := newConnectionGater(logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create gater")
	}

	tlsId := knockingtls.ID
	tls, err := knockingtls.NewKnockingTLS(logger, c.PrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not create knocking tls")
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.Identity(c.PrivKey),
		libp2p.DisableRelay(),
		libp2p.Security(tlsId, tls),
		libp2p.ConnectionGater(gater),
		libp2p.Peerstore(c.Peerstore),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	logger.Info("Peer: libp2p host booted", nil)

	return &concretePeer{
		Host:           basicHost,
		gater:          gater,
		tls:            tls,
		logger:         logger,
		endpointConfig: c.EndpointConfig,
		registrants:    make(map[types.ConfigDigest]struct{}),
		registrantsMu:  &sync.Mutex{},
	}, nil
}

func (p *concretePeer) MakeEndpoint(configDigest types.ConfigDigest, pids []string, bootstrappers []string, failureThreshold int) (types.BinaryNetworkEndpoint, error) {
	if failureThreshold <= 0 {
		return nil, errors.New("can't set F to 0 or smaller")
	}

	if len(bootstrappers) < 1 {
		return nil, errors.New("requires at least one bootstrapper")
	}
	peerIDs, err := decodePeerIDs(pids)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode peer IDs")
	}

	bnAddrs, err := decodeBootstrappers(bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode bootstrappers")
	}
	return newOCREndpoint(p.logger, configDigest, p, peerIDs, bnAddrs, p.endpointConfig, failureThreshold)
}

func decodeBootstrappers(bootstrappers []string) (bnAddrs []p2ppeer.AddrInfo, err error) {
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

func decodePeerIDs(pids []string) ([]p2ppeer.ID, error) {
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

func (p *concretePeer) MakeBootstrapper(configDigest types.ConfigDigest, pids []string, bootstrappers []string, F int) (types.Bootstrapper, error) {
	if F <= 0 {
		return nil, errors.New("can't set F to zero or smaller")
	}
	peerIDs, err := decodePeerIDs(pids)
	if err != nil {
		return nil, err
	}

	bnAddrs, err := decodeBootstrappers(bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode bootstrappers")
	}

	return newBootstrapper(p.logger, configDigest, p, peerIDs, bnAddrs, F)
}

func (p *concretePeer) register(r registrant) error {
	configDigest := r.getConfigDigest()
	p.logger.Debug("Peer: registering protocol handler", types.LogFields{
		"configDigest": configDigest.Hex(),
	})

	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()

	if _, ok := p.registrants[configDigest]; ok {
		return errors.Errorf("endpoint with getProtocol ID %s has already been registered", configDigest)
	}
	p.registrants[configDigest] = struct{}{}

	p.gater.add(r)

	p.tls.UpdateAllowlist(p.gater.allowlist())

	return nil
}

func (p *concretePeer) deregister(r registrant) error {
	configDigest := r.getConfigDigest()
	p.logger.Debug("Peer: deregistering protocol handler", types.LogFields{
		"ProtocolID": configDigest.Hex(),
	})

	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()

	if _, ok := p.registrants[configDigest]; !ok {
		return errors.Errorf("endpoint with getProtocol ID %s is not currently registered", configDigest)
	}
	p.registrants[configDigest] = struct{}{}

	p.gater.remove(r)

	p.tls.UpdateAllowlist(p.gater.allowlist())

	return nil
}

func (p *concretePeer) PeerID() string {
	return p.ID().Pretty()
}
