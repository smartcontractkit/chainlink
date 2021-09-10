package offchainreporting2

import (
	"io"
	"net"
	"strings"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type NetworkingConfig interface {
	OCR2DHTLookupInterval() int
	OCR2IncomingMessageBufferSize() int
	OCR2NewStreamTimeout() time.Duration
	OCR2OutgoingMessageBufferSize() int
	OCR2TraceLogging() bool
	OCR2P2PAnnounceIP() net.IP
	OCR2P2PAnnouncePort() uint16
	OCR2P2PBootstrapPeers() ([]string, error)
	OCR2P2PDHTAnnouncementCounterUserPrefix() uint32
	OCR2P2PListenIP() net.IP
	OCR2P2PListenPort() uint16
	OCR2P2PNetworkingStack() ocrnetworking.NetworkingStack
	OCR2P2PPeerID() (p2pkey.PeerID, error)
	OCR2P2PPeerstoreWriteInterval() time.Duration
	OCR2P2PV2AnnounceAddresses() []string
	OCR2P2PV2Bootstrappers() []ocrcommontypes.BootstrapperLocator
	OCR2P2PV2DeltaDial() time.Duration
	OCR2P2PV2DeltaReconcile() time.Duration
	OCR2P2PV2ListenAddresses() []string
}

type (
	peer interface {
		ocrtypes.BootstrapperFactory
		ocrtypes.BinaryNetworkEndpointFactory
		Close() error
	}

	peerAdapter struct {
		io.Closer
		ocrtypes.BinaryNetworkEndpointFactory
		ocrtypes.BootstrapperFactory
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		keyStore keystore.Master
		config   NetworkingConfig
		db       *gorm.DB

		pstoreWrapper *ocrcommon.Pstorewrapper
		PeerID        p2pkey.PeerID
		Peer          peer

		utils.StartStopOnce
	}
)

// NewSingletonPeerWrapper creates a new peer based on the p2p keys in the keystore
// It currently only supports one peerID/key
// It should be fairly easy to modify it to support multiple peerIDs/keys using e.g. a map
func NewSingletonPeerWrapper(keyStore keystore.Master, config NetworkingConfig, db *gorm.DB) *SingletonPeerWrapper {
	return &SingletonPeerWrapper{
		keyStore: keyStore,
		config:   config,
		db:       db,
	}
}

func (p *SingletonPeerWrapper) IsStarted() bool {
	return p.State() == utils.StartStopOnce_Started
}

func (p *SingletonPeerWrapper) Start() error {
	return p.StartOnce("SingletonPeerWrapper", func() (err error) {
		p2pkeys, err := p.keyStore.P2P().GetAll()
		if err != nil {
			return nil
		}
		listenPort := p.config.OCR2P2PListenPort()
		if listenPort == 0 {
			return errors.New("failed to instantiate oracle or bootstrapper service. If FEATURE_OFFCHAIN_REPORTING2 is on, then OCR2_P2P_LISTEN_PORT is required and must be set to a non-zero value")
		}

		if len(p2pkeys) == 0 {
			return nil
		}

		var key p2pkey.KeyV2
		var matched bool
		checkedKeys := []string{}
		configuredPeerID, err := p.config.OCR2P2PPeerID()
		if err != nil {
			return errors.Wrap(err, "failed to start peer wrapper")
		}
		for _, k := range p2pkeys {
			var peerID p2pkey.PeerID
			peerID = k.PeerID()
			if peerID == configuredPeerID {
				key = k
				matched = true
				break
			}
			checkedKeys = append(checkedKeys, peerID.String())
		}
		keys := strings.Join(checkedKeys, ", ")
		if !matched {
			if configuredPeerID == "" {
				return errors.Errorf("multiple p2p keys found but peer ID was not set. You must specify OCR2_P2P_PEER_ID if you have more than one key. Keys available: %s", keys)
			}
			return errors.Errorf("multiple p2p keys found but none matched the given OCR2_P2P_PEER_ID of '%s'. Keys available: %s", configuredPeerID, keys)
		}

		p.PeerID = key.PeerID()
		p.pstoreWrapper, err = ocrcommon.NewPeerstoreWrapper(p.db, p.config.OCR2P2PPeerstoreWriteInterval(), p.PeerID)
		if err != nil {
			return errors.Wrap(err, "could not make new pstorewrapper")
		}
		sqlDB, err := p.db.DB()
		if err != nil {
			return err
		}
		discovererDB := NewDiscovererDatabase(sqlDB, p2ppeer.ID(p.PeerID))

		// If the P2PAnnounceIP is set we must also set the P2PAnnouncePort
		// Fallback to P2PListenPort if it wasn't made explicit
		var announcePort uint16
		if p.config.OCR2P2PAnnounceIP() != nil && p.config.OCR2P2PAnnouncePort() != 0 {
			announcePort = p.config.OCR2P2PAnnouncePort()
		} else if p.config.OCR2P2PAnnounceIP() != nil {
			announcePort = listenPort
		}

		peerLogger := ocrcommon.NewLogger(logger.Default, p.config.OCR2TraceLogging(), func(string) {})

		peerConfig := ocrnetworking.PeerConfig{
			NetworkingStack: p.config.OCR2P2PNetworkingStack(),
			PrivKey:         key.PrivKey,
			// XXX: These may be obsolete in the future
			V1ListenIP:           p.config.OCR2P2PListenIP(),
			V1ListenPort:         listenPort,
			V1AnnounceIP:         p.config.OCR2P2PAnnounceIP(),
			V1AnnouncePort:       announcePort,
			Logger:               peerLogger,
			V1Peerstore:          p.pstoreWrapper.Peerstore,
			V2ListenAddresses:    p.config.OCR2P2PV2ListenAddresses(),
			V2AnnounceAddresses:  p.config.OCR2P2PV2AnnounceAddresses(),
			V2DeltaReconcile:     p.config.OCR2P2PV2DeltaReconcile(),
			V2DeltaDial:          p.config.OCR2P2PV2DeltaDial(),
			V2DiscovererDatabase: discovererDB,
			EndpointConfig: ocrnetworking.EndpointConfig{
				IncomingMessageBufferSize: p.config.OCR2IncomingMessageBufferSize(),
				OutgoingMessageBufferSize: p.config.OCR2OutgoingMessageBufferSize(),
				NewStreamTimeout:          p.config.OCR2NewStreamTimeout(),
				DHTLookupInterval:         p.config.OCR2DHTLookupInterval(),
			},
			V1DHTAnnouncementCounterUserPrefix: p.config.OCR2P2PDHTAnnouncementCounterUserPrefix(),
		}
		logger.Debugw("Creating OCR2 Peer", "config", peerConfig)
		peer, err := ocrnetworking.NewPeer(peerConfig)
		if err != nil {
			return errors.Wrap(err, "error calling NewPeer")
		}
		p.Peer = peerAdapter{
			peer,
			peer.GenOCRBinaryNetworkEndpointFactory(),
			peer.GenOCRBootstrapperFactory(),
		}
		return p.pstoreWrapper.Start()
	})
}

// Close closes the peer and peerstore
func (p *SingletonPeerWrapper) Close() error {
	return p.StopOnce("SingletonPeerWrapper", func() (err error) {
		if p.Peer != nil {
			err = p.Peer.Close()
		}

		if p.pstoreWrapper != nil {
			err = multierr.Combine(err, p.pstoreWrapper.Close())
		}

		return err
	})
}
