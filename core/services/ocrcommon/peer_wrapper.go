package ocrcommon

import (
	"io"
	"net"

	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/sqlx"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocrnetworkingtypes "github.com/smartcontractkit/libocr/networking/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"go.uber.org/multierr"
)

type PeerWrapperConfig interface {
	config.P2PNetworking
	OCRTraceLogging() bool
	LogSQL() bool
}

type (
	peerAdapter struct {
		io.Closer
		ocrtypes.BinaryNetworkEndpointFactory
		ocrtypes.BootstrapperFactory
	}

	peerAdapter2 struct {
		io.Closer
		ocr2types.BinaryNetworkEndpointFactory
		ocr2types.BootstrapperFactory
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		keyStore keystore.Master
		config   PeerWrapperConfig
		db       *sqlx.DB
		lggr     logger.Logger
		PeerID   p2pkey.PeerID

		// V1 peer
		pstoreWrapper *Pstorewrapper
		Peer          *peerAdapter

		// V2 peer
		Peer2 *peerAdapter2

		utils.StartStopOnce
	}
)

// TODO
func ValidatePeerWrapperConfig(config PeerWrapperConfig) error {
	switch config.P2PNetworkingStack() {
	case ocrnetworking.NetworkingStackV1:
	case ocrnetworking.NetworkingStackV2:
	case ocrnetworking.NetworkingStackV1V2:
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

func (p *SingletonPeerWrapper) Start() error {
	return p.StartOnce("SingletonPeerWrapper", func() (err error) {
		// If there are no keys, permit the peer to start without a key
		// TODO: This appears only a requirement for the tests, normally the node
		// always ensures a key is available on boot. We should update the tests
		// but there is a lot of them...
		if ks, err := p.keyStore.P2P().GetAll(); err == nil && len(ks) == 0 {
			p.lggr.Warn("No P2P keys found in keystore. Peer wrapper will not be fully initialized")
			return nil
		}
		key, err := p.keyStore.P2P().Get(p.config.P2PPeerID())
		if err != nil {
			return err
		}
		p.PeerID = key.PeerID()

		// We need to start the peer store wrapper if v1 is required.
		// Also fallback to listen params if announce params not specified.
		var announcePort uint16
		var announceIP net.IP
		var peerStore p2ppeerstore.Peerstore
		if p.config.P2PNetworkingStack() == ocrnetworking.NetworkingStackV1 || p.config.P2PNetworkingStack() == ocrnetworking.NetworkingStackV1V2 {
			p.pstoreWrapper, err = NewPeerstoreWrapper(p.db, p.config.P2PPeerstoreWriteInterval(), p.PeerID, p.lggr, p.config)
			if err != nil {
				return errors.Wrap(err, "could not make new pstorewrapper")
			}
			if err = p.pstoreWrapper.Start(); err != nil {
				return errors.Wrap(err, "failed to start peer store wrapper")
			}

			peerStore = p.pstoreWrapper.Peerstore
			announcePort = p.config.P2PListenPort()
			if p.config.P2PAnnouncePort() != 0 {
				announcePort = p.config.P2PAnnouncePort()
			}
			announceIP = p.config.P2PListenIP()
			if p.config.P2PAnnounceIP() != nil {
				announceIP = p.config.P2PAnnounceIP()
			}
		}

		// Discover DB is only required for v2
		// Also fallback to listen addresses if announce not specified
		var discovererDB ocrnetworkingtypes.DiscovererDatabase
		var announceAddresses []string
		if p.config.P2PNetworkingStack() == ocrnetworking.NetworkingStackV2 {
			discovererDB = NewDiscovererDatabase(p.db.DB, p2ppeer.ID(p.PeerID))
			announceAddresses = p.config.P2PV2ListenAddresses()
			if len(p.config.P2PV2AnnounceAddresses()) != 0 {
				announceAddresses = p.config.P2PV2AnnounceAddresses()
			}
		}

		peerConfig := ocrnetworking.PeerConfig{
			NetworkingStack: p.config.P2PNetworkingStack(),
			PrivKey:         key.PrivKey,
			Logger:          logger.NewOCRWrapper(p.lggr, p.config.OCRTraceLogging(), func(string) {}),

			// V1 config
			V1ListenIP:                         p.config.P2PListenIP(),
			V1ListenPort:                       p.config.P2PListenPort(),
			V1AnnounceIP:                       announceIP,
			V1AnnouncePort:                     announcePort,
			V1Peerstore:                        peerStore,
			V1DHTAnnouncementCounterUserPrefix: p.config.P2PDHTAnnouncementCounterUserPrefix(),

			// V2 config
			V2ListenAddresses:    p.config.P2PV2ListenAddresses(),
			V2AnnounceAddresses:  announceAddresses,
			V2DeltaReconcile:     p.config.P2PV2DeltaReconcile().Duration(),
			V2DeltaDial:          p.config.P2PV2DeltaDial().Duration(),
			V2DiscovererDatabase: discovererDB,

			EndpointConfig: ocrnetworking.EndpointConfig{
				// V1 and V2 config
				IncomingMessageBufferSize: p.config.P2PIncomingMessageBufferSize(),
				OutgoingMessageBufferSize: p.config.P2POutgoingMessageBufferSize(),

				// V1 Config
				NewStreamTimeout:       p.config.P2PNewStreamTimeout(),
				DHTLookupInterval:      p.config.P2PDHTLookupInterval(),
				BootstrapCheckInterval: p.config.P2PBootstrapCheckInterval(),
			},
		}

		p.lggr.Debugw("Creating OCR/OCR2 Peer", "config", peerConfig)
		peer, err := ocrnetworking.NewPeer(peerConfig)
		if err != nil {
			return errors.Wrap(err, "error calling NewPeer")
		}
		p.Peer = &peerAdapter{
			peer,
			peer.OCR1BinaryNetworkEndpointFactory(),
			peer.OCR1BootstrapperFactory(),
		}
		p.Peer2 = &peerAdapter2{
			peer,
			peer.OCR2BinaryNetworkEndpointFactory(),
			peer.OCR2BootstrapperFactory(),
		}
		return nil
	})
}

// Close closes the peer and peerstore
func (p *SingletonPeerWrapper) Close() error {
	return p.StopOnce("SingletonPeerWrapper", func() (err error) {
		if p.Peer != nil {
			err = p.Peer.Close()
		}
		if p.Peer2 != nil {
			err = p.Peer2.Close()
		}

		if p.pstoreWrapper != nil {
			err = multierr.Combine(err, p.pstoreWrapper.Close())
		}

		return err
	})
}
