package ocrcommon

import (
	"context"
	"fmt"
	"io"
	"net"

	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocrnetworkingtypes "github.com/smartcontractkit/libocr/networking/types"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PeerWrapperConfig interface {
	config.P2PNetworking
	config.P2PV1Networking
	config.P2PV2Networking
	pg.QConfig
	OCRTraceLogging() bool
	FeatureOffchainReporting() bool
}

type (
	peerAdapterOCR1 struct {
		ocr1types.BinaryNetworkEndpointFactory
		ocr1types.BootstrapperFactory
	}

	peerAdapterOCR2 struct {
		ocr2types.BinaryNetworkEndpointFactory
		ocr2types.BootstrapperFactory
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		utils.StartStopOnce
		keyStore      keystore.Master
		config        PeerWrapperConfig
		db            *sqlx.DB
		lggr          logger.Logger
		PeerID        p2pkey.PeerID
		pstoreWrapper *Pstorewrapper

		// Used at shutdown to stop all of this peer's goroutines
		peerCloser io.Closer

		// OCR1 peer adapter
		Peer1 *peerAdapterOCR1

		// OCR2 peer adapter
		Peer2 *peerAdapterOCR2
	}
)

func ValidatePeerWrapperConfig(config PeerWrapperConfig) error {
	switch config.P2PNetworkingStack() {
	case ocrnetworking.NetworkingStackV1:
		// Note: If P2PListenPort isn't set, the peer wrapper will generate a random one itself.
		return nil
	case ocrnetworking.NetworkingStackV2:
		if len(config.P2PV2ListenAddresses()) == 0 {
			return errors.New("networking stack v2 selected but no P2P.V2.ListenAddresses specified")
		}
	case ocrnetworking.NetworkingStackV1V2:
		if len(config.P2PV2ListenAddresses()) == 0 {
			return errors.New("networking stack v1v2 selected but no P2P.V2.ListenAddresses specified")
		}
	default:
		return errors.New("unknown networking stack")
	}
	return nil
}

// NewSingletonPeerWrapper creates a new peer based on the p2p keys in the keystore
// It currently only supports one peerID/key
// It should be fairly easy to modify it to support multiple peerIDs/keys using e.g. a map
func NewSingletonPeerWrapper(keyStore keystore.Master, config PeerWrapperConfig, db *sqlx.DB, lggr logger.Logger) *SingletonPeerWrapper {
	return &SingletonPeerWrapper{
		keyStore: keyStore,
		config:   config,
		db:       db,
		lggr:     lggr.Named("SingletonPeerWrapper"),
	}
}

func (p *SingletonPeerWrapper) IsStarted() bool {
	return p.State() == utils.StartStopOnce_Started
}

// Start starts SingletonPeerWrapper.
func (p *SingletonPeerWrapper) Start(context.Context) error {
	return p.StartOnce("SingletonPeerWrapper", func() error {
		peerConfig, err := p.peerConfig()
		if err != nil {
			return err
		}

		p.lggr.Debugw("Creating OCR/OCR2 Peer", "config", peerConfig)
		// Note: creates and starts the peer
		peer, err := ocrnetworking.NewPeer(peerConfig)
		if err != nil {
			return errors.Wrap(err, "error calling NewPeer")
		}
		p.Peer1 = &peerAdapterOCR1{
			peer.OCR1BinaryNetworkEndpointFactory(),
			peer.OCR1BootstrapperFactory(),
		}
		p.Peer2 = &peerAdapterOCR2{
			peer.OCR2BinaryNetworkEndpointFactory(),
			peer.OCR2BootstrapperFactory(),
		}
		p.peerCloser = peer
		return nil
	})
}

func (p *SingletonPeerWrapper) peerConfig() (ocrnetworking.PeerConfig, error) {
	// Peer wrapper panics if no p2p keys are present.
	if ks, err := p.keyStore.P2P().GetAll(); err == nil && len(ks) == 0 {
		return ocrnetworking.PeerConfig{}, errors.Errorf("No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	}
	key, err := p.keyStore.P2P().GetOrFirst(p.config.P2PPeerID())
	if err != nil {
		return ocrnetworking.PeerConfig{}, err
	}
	p.PeerID = key.PeerID()

	p2pPort := p.config.P2PListenPort()
	// We need to start the peer store wrapper if v1 is required.
	// Also fallback to listen params if announce params not specified.
	ns := p.config.P2PNetworkingStack()
	// NewPeer requires that these are both set or unset, otherwise it will error out.
	v1AnnounceIP, v1AnnouncePort := p.config.P2PAnnounceIP(), p.config.P2PAnnouncePort()
	var peerStore p2ppeerstore.Peerstore
	if ns == ocrnetworking.NetworkingStackV1 || ns == ocrnetworking.NetworkingStackV1V2 {
		p.pstoreWrapper, err = NewPeerstoreWrapper(p.db, p.config.P2PPeerstoreWriteInterval(), p.PeerID, p.lggr, p.config)
		if err != nil {
			return ocrnetworking.PeerConfig{}, errors.Wrap(err, "could not make new pstorewrapper")
		}
		if err = p.pstoreWrapper.Start(); err != nil {
			return ocrnetworking.PeerConfig{}, errors.Wrap(err, "failed to start peer store wrapper")
		}

		peerStore = p.pstoreWrapper.Peerstore

		// Use a random port if the port hasn't been set explicitly.
		if p2pPort == 0 {
			port, perr := p.randomPort()
			if perr != nil {
				return ocrnetworking.PeerConfig{}, perr
			}
			p2pPort = port

			p.lggr.Warnw(
				fmt.Sprintf("P2PListenPort was not set, listening on random port %d. A new random port will be generated on every boot, for stability it is recommended to set P2PListenPort to a fixed value in your environment", p2pPort),
				"p2pPort",
				p2pPort,
			)
		}

		// Support someone specifying only the announce IP but leaving out
		// the port.
		// We _should not_ handle the case of someone specifying only the
		// port but leaving out the IP, because the listen IP is typically
		// an unspecified IP (https://pkg.go.dev/net#IP.IsUnspecified) and
		// using that for the announce IP will cause other peers to not be
		// able to connect.
		if v1AnnounceIP != nil && v1AnnouncePort == 0 {
			v1AnnouncePort = p2pPort
		}
	}

	// Discover DB is only required for v2
	var discovererDB ocrnetworkingtypes.DiscovererDatabase
	if ns == ocrnetworking.NetworkingStackV2 || ns == ocrnetworking.NetworkingStackV1V2 {
		discovererDB = NewDiscovererDatabase(p.db.DB, p2ppeer.ID(p.PeerID))
	}

	peerConfig := ocrnetworking.PeerConfig{
		NetworkingStack: p.config.P2PNetworkingStack(),
		PrivKey:         key.PrivKey,
		Logger:          logger.NewOCRWrapper(p.lggr, p.config.OCRTraceLogging(), func(string) {}),

		// V1 config
		V1ListenIP:                         p.config.P2PListenIP(),
		V1ListenPort:                       p2pPort,
		V1AnnounceIP:                       v1AnnounceIP,
		V1AnnouncePort:                     v1AnnouncePort,
		V1Peerstore:                        peerStore,
		V1DHTAnnouncementCounterUserPrefix: p.config.P2PDHTAnnouncementCounterUserPrefix(),

		// V2 config
		V2ListenAddresses:    p.config.P2PV2ListenAddresses(),
		V2AnnounceAddresses:  p.config.P2PV2AnnounceAddresses(), // NewPeer will handle the fallback to listen addresses for us.
		V2DeltaReconcile:     p.config.P2PV2DeltaReconcile().Duration(),
		V2DeltaDial:          p.config.P2PV2DeltaDial().Duration(),
		V2DiscovererDatabase: discovererDB,

		V1EndpointConfig: ocrnetworking.EndpointConfigV1{
			IncomingMessageBufferSize: p.config.P2PIncomingMessageBufferSize(),
			OutgoingMessageBufferSize: p.config.P2POutgoingMessageBufferSize(),
			NewStreamTimeout:          p.config.P2PNewStreamTimeout(),
			DHTLookupInterval:         p.config.P2PDHTLookupInterval(),
			BootstrapCheckInterval:    p.config.P2PBootstrapCheckInterval(),
		},

		V2EndpointConfig: ocrnetworking.EndpointConfigV2{
			IncomingMessageBufferSize: p.config.P2PIncomingMessageBufferSize(),
			OutgoingMessageBufferSize: p.config.P2POutgoingMessageBufferSize(),
		},
	}

	return peerConfig, nil
}

func (p *SingletonPeerWrapper) randomPort() (uint16, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("unexpected ResolveTCPAddr error generating random P2PListenPort: %w", err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("unexpected ListenTCP error generating random P2PListenPort: %w", err)
	}
	defer l.Close()

	return uint16(l.Addr().(*net.TCPAddr).Port), nil
}

// Close closes the peer and peerstore
func (p *SingletonPeerWrapper) Close() error {
	return p.StopOnce("SingletonPeerWrapper", func() (err error) {
		if p.peerCloser != nil {
			err = p.peerCloser.Close()
		}
		if p.pstoreWrapper != nil {
			err = multierr.Combine(err, p.pstoreWrapper.Close())
		}
		return err
	})
}

func (p *SingletonPeerWrapper) Name() string {
	return p.lggr.Name()
}

func (p *SingletonPeerWrapper) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
}

func (p *SingletonPeerWrapper) Config() PeerWrapperConfig {
	return p.config
}
