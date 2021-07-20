package offchainreporting

import (
	"net"
	"strings"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type NetworkingConfig interface {
	OCRBootstrapCheckInterval() time.Duration
	OCRDHTLookupInterval() int
	OCRIncomingMessageBufferSize() int
	OCRNewStreamTimeout() time.Duration
	OCROutgoingMessageBufferSize() int
	OCRTraceLogging() bool
	P2PAnnounceIP() net.IP
	P2PAnnouncePort() uint16
	P2PBootstrapPeers([]string) ([]string, error)
	P2PDHTAnnouncementCounterUserPrefix() uint32
	P2PListenIP() net.IP
	P2PListenPort() uint16
	P2PNetworkingStack() ocrnetworking.NetworkingStack
	P2PPeerID(override *p2pkey.PeerID) (p2pkey.PeerID, error)
	P2PPeerstoreWriteInterval() time.Duration
	P2PV2AnnounceAddresses() []string
	P2PV2Bootstrappers() []ocrtypes.BootstrapperLocator
	P2PV2DeltaDial() models.Duration
	P2PV2DeltaReconcile() models.Duration
	P2PV2ListenAddresses() []string
}

type (
	peer interface {
		ocrtypes.BootstrapperFactory
		ocrtypes.BinaryNetworkEndpointFactory
		Close() error
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		keyStore *keystore.OCR
		config   NetworkingConfig
		db       *gorm.DB

		pstoreWrapper *Pstorewrapper
		PeerID        p2pkey.PeerID
		Peer          peer

		utils.StartStopOnce
	}
)

// NewSingletonPeerWrapper creates a new peer based on the p2p keys in the keystore
// It currently only supports one peerID/key
// It should be fairly easy to modify it to support multiple peerIDs/keys using e.g. a map
func NewSingletonPeerWrapper(keyStore *keystore.OCR, config *config.Config, db *gorm.DB) *SingletonPeerWrapper {
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
		p2pkeys := p.keyStore.DecryptedP2PKeys()
		listenPort := p.config.P2PListenPort()
		if listenPort == 0 {
			return errors.New("failed to instantiate oracle or bootstrapper service. If FEATURE_OFFCHAIN_REPORTING is on, then P2P_LISTEN_PORT is required and must be set to a non-zero value")
		}

		if len(p2pkeys) == 0 {
			return nil
		}

		var key p2pkey.Key
		var matched bool
		checkedKeys := []string{}
		configuredPeerID, err := p.config.P2PPeerID(nil)
		if err != nil {
			return errors.Wrap(err, "failed to start peer wrapper")
		}
		for _, k := range p2pkeys {
			var peerID p2pkey.PeerID
			peerID, err = k.GetPeerID()
			if err != nil {
				return errors.Wrap(err, "unexpectedly failed to get peer ID from key")
			}
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
				return errors.Errorf("multiple p2p keys found but peer ID was not set. You must specify P2P_PEER_ID if you have more than one key. Keys available: %s", keys)
			}
			return errors.Errorf("multiple p2p keys found but none matched the given P2P_PEER_ID of '%s'. Keys available: %s", configuredPeerID, keys)
		}

		p.PeerID, err = key.GetPeerID()
		if err != nil {
			return errors.Wrap(err, "could not get peer ID")
		}
		p.pstoreWrapper, err = NewPeerstoreWrapper(p.db, p.config.P2PPeerstoreWriteInterval(), p.PeerID)
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
		if p.config.P2PAnnounceIP() != nil && p.config.P2PAnnouncePort() != 0 {
			announcePort = p.config.P2PAnnouncePort()
		} else if p.config.P2PAnnounceIP() != nil {
			announcePort = listenPort
		}

		peerLogger := NewLogger(logger.Default, p.config.OCRTraceLogging(), func(string) {})

		p.Peer, err = ocrnetworking.NewPeer(ocrnetworking.PeerConfig{
			NetworkingStack:      p.config.P2PNetworkingStack(),
			PrivKey:              key.PrivKey,
			V1ListenIP:           p.config.P2PListenIP(),
			V1ListenPort:         listenPort,
			V1AnnounceIP:         p.config.P2PAnnounceIP(),
			V1AnnouncePort:       announcePort,
			Logger:               peerLogger,
			V1Peerstore:          p.pstoreWrapper.Peerstore,
			V2ListenAddresses:    p.config.P2PV2ListenAddresses(),
			V2AnnounceAddresses:  p.config.P2PV2AnnounceAddresses(),
			V2DeltaReconcile:     p.config.P2PV2DeltaReconcile().Duration(),
			V2DeltaDial:          p.config.P2PV2DeltaDial().Duration(),
			V2DiscovererDatabase: discovererDB,
			EndpointConfig: ocrnetworking.EndpointConfig{
				IncomingMessageBufferSize: p.config.OCRIncomingMessageBufferSize(),
				OutgoingMessageBufferSize: p.config.OCROutgoingMessageBufferSize(),
				NewStreamTimeout:          p.config.OCRNewStreamTimeout(),
				DHTLookupInterval:         p.config.OCRDHTLookupInterval(),
				BootstrapCheckInterval:    p.config.OCRBootstrapCheckInterval(),
			},
			V1DHTAnnouncementCounterUserPrefix: p.config.P2PDHTAnnouncementCounterUserPrefix(),
		})
		if err != nil {
			return errors.Wrap(err, "error calling NewPeer")
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
