package config

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"

	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	"github.com/libp2p/go-libp2p/p2p/host/relay"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"

	autonat "github.com/libp2p/go-libp2p-autonat"
	blankhost "github.com/libp2p/go-libp2p-blankhost"
	circuit "github.com/libp2p/go-libp2p-circuit"
	discovery "github.com/libp2p/go-libp2p-discovery"
	swarm "github.com/libp2p/go-libp2p-swarm"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"

	logging "github.com/ipfs/go-log"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("p2p-config")

// AddrsFactory is a function that takes a set of multiaddrs we're listening on and
// returns the set of multiaddrs we should advertise to the network.
type AddrsFactory = bhost.AddrsFactory

// NATManagerC is a NATManager constructor.
type NATManagerC func(network.Network) bhost.NATManager

type RoutingC func(host.Host) (routing.PeerRouting, error)

// AutoNATConfig defines the AutoNAT behavior for the libp2p host.
type AutoNATConfig struct {
	ForceReachability   *network.Reachability
	EnableService       bool
	ThrottleGlobalLimit int
	ThrottlePeerLimit   int
	ThrottleInterval    time.Duration
}

// Config describes a set of settings for a libp2p node
//
// This is *not* a stable interface. Use the options defined in the root
// package.
type Config struct {
	// UserAgent is the identifier this node will send to other peers when
	// identifying itself, e.g. via the identify protocol.
	//
	// Set it via the UserAgent option function.
	UserAgent string

	PeerKey crypto.PrivKey

	Transports         []TptC
	Muxers             []MsMuxC
	SecurityTransports []MsSecC
	Insecure           bool
	PSK                pnet.PSK

	RelayCustom bool
	Relay       bool
	RelayOpts   []circuit.RelayOpt

	ListenAddrs     []ma.Multiaddr
	AddrsFactory    bhost.AddrsFactory
	ConnectionGater connmgr.ConnectionGater

	ConnManager connmgr.ConnManager
	NATManager  NATManagerC
	Peerstore   peerstore.Peerstore
	Reporter    metrics.Reporter

	DisablePing bool

	Routing RoutingC

	EnableAutoRelay bool
	AutoNATConfig
	StaticRelays []peer.AddrInfo
}

func (cfg *Config) makeSwarm(ctx context.Context) (*swarm.Swarm, error) {
	if cfg.Peerstore == nil {
		return nil, fmt.Errorf("no peerstore specified")
	}

	// Check this early. Prevents us from even *starting* without verifying this.
	if pnet.ForcePrivateNetwork && len(cfg.PSK) == 0 {
		log.Error("tried to create a libp2p node with no Private" +
			" Network Protector but usage of Private Networks" +
			" is forced by the enviroment")
		// Note: This is *also* checked the upgrader itself so it'll be
		// enforced even *if* you don't use the libp2p constructor.
		return nil, pnet.ErrNotInPrivateNetwork
	}

	if cfg.PeerKey == nil {
		return nil, fmt.Errorf("no peer key specified")
	}

	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(cfg.PeerKey.GetPublic())
	if err != nil {
		return nil, err
	}

	if err := cfg.Peerstore.AddPrivKey(pid, cfg.PeerKey); err != nil {
		return nil, err
	}
	if err := cfg.Peerstore.AddPubKey(pid, cfg.PeerKey.GetPublic()); err != nil {
		return nil, err
	}

	// TODO: Make the swarm implementation configurable.
	swrm := swarm.NewSwarm(ctx, pid, cfg.Peerstore, cfg.Reporter, cfg.ConnectionGater)
	return swrm, nil
}

func (cfg *Config) addTransports(ctx context.Context, h host.Host) (err error) {
	swrm, ok := h.Network().(transport.TransportNetwork)
	if !ok {
		// Should probably skip this if no transports.
		return fmt.Errorf("swarm does not support transports")
	}
	upgrader := new(tptu.Upgrader)
	upgrader.PSK = cfg.PSK
	upgrader.ConnGater = cfg.ConnectionGater
	if cfg.Insecure {
		upgrader.Secure = makeInsecureTransport(h.ID(), cfg.PeerKey)
	} else {
		upgrader.Secure, err = makeSecurityTransport(h, cfg.SecurityTransports)
		if err != nil {
			return err
		}
	}

	upgrader.Muxer, err = makeMuxer(h, cfg.Muxers)
	if err != nil {
		return err
	}

	tpts, err := makeTransports(h, upgrader, cfg.ConnectionGater, cfg.Transports)
	if err != nil {
		return err
	}
	for _, t := range tpts {
		err = swrm.AddTransport(t)
		if err != nil {
			return err
		}
	}

	if cfg.Relay {
		err := circuit.AddRelayTransport(ctx, h, upgrader, cfg.RelayOpts...)
		if err != nil {
			h.Close()
			return err
		}
	}

	return nil
}

// NewNode constructs a new libp2p Host from the Config.
//
// This function consumes the config. Do not reuse it (really!).
func (cfg *Config) NewNode(ctx context.Context) (host.Host, error) {
	swrm, err := cfg.makeSwarm(ctx)
	if err != nil {
		return nil, err
	}

	h, err := bhost.NewHost(ctx, swrm, &bhost.HostOpts{
		ConnManager:  cfg.ConnManager,
		AddrsFactory: cfg.AddrsFactory,
		NATManager:   cfg.NATManager,
		EnablePing:   !cfg.DisablePing,
		UserAgent:    cfg.UserAgent,
	})

	if err != nil {
		swrm.Close()
		return nil, err
	}

	// XXX: This is the only sane way to get a context out that's guaranteed
	// to be canceled when we shut down.
	//
	// TODO: Stop using contexts to stop services. This is just lazy.
	// Contexts are for canceling requests, services should be managed
	// explicitly.
	ctx = swrm.Context()

	if cfg.Relay {
		// If we've enabled the relay, we should filter out relay
		// addresses by default.
		//
		// TODO: We shouldn't be doing this here.
		oldFactory := h.AddrsFactory
		h.AddrsFactory = func(addrs []ma.Multiaddr) []ma.Multiaddr {
			return oldFactory(relay.Filter(addrs))
		}
	}

	err = cfg.addTransports(ctx, h)
	if err != nil {
		h.Close()
		return nil, err
	}

	// TODO: This method succeeds if listening on one address succeeds. We
	// should probably fail if listening on *any* addr fails.
	if err := h.Network().Listen(cfg.ListenAddrs...); err != nil {
		h.Close()
		return nil, err
	}

	// Configure routing and autorelay
	var router routing.PeerRouting
	if cfg.Routing != nil {
		router, err = cfg.Routing(h)
		if err != nil {
			h.Close()
			return nil, err
		}
	}

	// Note: h.AddrsFactory may be changed by AutoRelay, but non-relay version is
	// used by AutoNAT below.
	addrF := h.AddrsFactory
	if cfg.EnableAutoRelay {
		if !cfg.Relay {
			h.Close()
			return nil, fmt.Errorf("cannot enable autorelay; relay is not enabled")
		}

		hop := false
		for _, opt := range cfg.RelayOpts {
			if opt == circuit.OptHop {
				hop = true
				break
			}
		}

		if !hop && len(cfg.StaticRelays) > 0 {
			_ = relay.NewAutoRelay(ctx, h, nil, router, cfg.StaticRelays)
		} else {
			if router == nil {
				h.Close()
				return nil, fmt.Errorf("cannot enable autorelay; no routing for discovery")
			}

			crouter, ok := router.(routing.ContentRouting)
			if !ok {
				h.Close()
				return nil, fmt.Errorf("cannot enable autorelay; no suitable routing for discovery")
			}

			discovery := discovery.NewRoutingDiscovery(crouter)

			if hop {
				// advertise ourselves
				relay.Advertise(ctx, discovery)
			} else {
				_ = relay.NewAutoRelay(ctx, h, discovery, router, cfg.StaticRelays)
			}
		}
	}

	autonatOpts := []autonat.Option{
		autonat.UsingAddresses(func() []ma.Multiaddr {
			return addrF(h.AllAddrs())
		}),
	}
	if cfg.AutoNATConfig.ThrottleInterval != 0 {
		autonatOpts = append(autonatOpts,
			autonat.WithThrottling(cfg.AutoNATConfig.ThrottleGlobalLimit, cfg.AutoNATConfig.ThrottleInterval),
			autonat.WithPeerThrottling(cfg.AutoNATConfig.ThrottlePeerLimit))
	}
	if cfg.AutoNATConfig.EnableService {
		autonatPrivKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}

		// Pull out the pieces of the config that we _actually_ care about.
		// Specifically, don't setup things like autorelay, listeners,
		// identify, etc.
		autoNatCfg := Config{
			Transports:         cfg.Transports,
			Muxers:             cfg.Muxers,
			SecurityTransports: cfg.SecurityTransports,
			Insecure:           cfg.Insecure,
			PSK:                cfg.PSK,
			ConnectionGater:    cfg.ConnectionGater,
			Reporter:           cfg.Reporter,
			PeerKey:            autonatPrivKey,

			Peerstore: pstoremem.NewPeerstore(),
		}

		dialer, err := autoNatCfg.makeSwarm(ctx)
		if err != nil {
			h.Close()
			return nil, err
		}
		dialerHost := blankhost.NewBlankHost(dialer)
		err = autoNatCfg.addTransports(ctx, dialerHost)
		if err != nil {
			dialerHost.Close()
			h.Close()
			return nil, err
		}
		// NOTE: We're dropping the blank host here but that's fine. It
		// doesn't really _do_ anything and doesn't even need to be
		// closed (as long as we close the underlying network).
		autonatOpts = append(autonatOpts, autonat.EnableService(dialerHost.Network()))
	}
	if cfg.AutoNATConfig.ForceReachability != nil {
		autonatOpts = append(autonatOpts, autonat.WithReachability(*cfg.AutoNATConfig.ForceReachability))
	}

	if h.AutoNat, err = autonat.New(ctx, h, autonatOpts...); err != nil {
		h.Close()
		return nil, fmt.Errorf("cannot enable autorelay; autonat failed to start: %v", err)
	}

	// start the host background tasks
	h.Start()

	if router != nil {
		return routed.Wrap(h, router), nil
	}
	return h, nil
}

// Option is a libp2p config option that can be given to the libp2p constructor
// (`libp2p.New`).
type Option func(cfg *Config) error

// Apply applies the given options to the config, returning the first error
// encountered (if any).
func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}
