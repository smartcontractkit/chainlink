package offchainreporting

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type (
	peer interface {
		ocrtypes.BootstrapperFactory
		ocrtypes.BinaryNetworkEndpointFactory
		Close() error
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		keyStore *keystore.OCRKeyStore
		config   *orm.Config
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
func NewSingletonPeerWrapper(keyStore *keystore.OCRKeyStore, config *orm.Config, db *gorm.DB) *SingletonPeerWrapper {
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
			PrivKey:      key.PrivKey,
			ListenIP:     p.config.P2PListenIP(),
			ListenPort:   listenPort,
			AnnounceIP:   p.config.P2PAnnounceIP(),
			AnnouncePort: announcePort,
			Logger:       peerLogger,
			Peerstore:    p.pstoreWrapper.Peerstore,
			EndpointConfig: ocrnetworking.EndpointConfig{
				IncomingMessageBufferSize: p.config.OCRIncomingMessageBufferSize(),
				OutgoingMessageBufferSize: p.config.OCROutgoingMessageBufferSize(),
				NewStreamTimeout:          p.config.OCRNewStreamTimeout(),
				DHTLookupInterval:         p.config.OCRDHTLookupInterval(),
				BootstrapCheckInterval:    p.config.OCRBootstrapCheckInterval(),
			},
			DHTAnnouncementCounterUserPrefix: p.config.P2PDHTAnnouncementCounterUserPrefix(),
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
