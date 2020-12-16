package offchainreporting

import (
	"sync"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type (
	peer interface {
		ocrtypes.BootstrapperFactory
		ocrtypes.BinaryNetworkEndpointFactory
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		keyStore *KeyStore
		config   *orm.Config
		db       *gorm.DB

		pstoreWrapper *Pstorewrapper
		PeerID        models.PeerID
		Peer          peer

		startMu *sync.Mutex
		started bool
	}
)

// NewSingletonPeerWrapper creates a new peer based on the p2p keys in the keystore
// It currently only supports one peerID/key
// It should be fairly easy to modify it to support multiple peerIDs/keys using e.g. a map
func NewSingletonPeerWrapper(keyStore *KeyStore, config *orm.Config, db *gorm.DB) *SingletonPeerWrapper {
	return &SingletonPeerWrapper{keyStore, config, db, nil, "", nil, new(sync.Mutex), false}
}

func (p *SingletonPeerWrapper) IsStarted() bool {
	p.startMu.Lock()
	defer p.startMu.Unlock()
	return p.started
}

func (p *SingletonPeerWrapper) Start() (err error) {
	p.startMu.Lock()
	defer p.startMu.Unlock()

	if p.started {
		return errors.New("already started")
	}

	p.started = true

	p2pkeys := p.keyStore.DecryptedP2PKeys()
	listenPort := p.config.P2PListenPort()
	if listenPort == 0 {
		return errors.New("failed to instantiate oracle or bootstrapper service, P2P_LISTEN_PORT is required and must be set to a non-zero value")
	}

	if len(p2pkeys) == 0 {
		return nil
	}
	if len(p2pkeys) > 1 {
		return errors.New("more than one p2p key is not currently supported")
	}

	key := p2pkeys[0]
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
}
func (p SingletonPeerWrapper) Close() error {
	p.startMu.Lock()
	defer p.startMu.Unlock()
	if !p.started {
		return errors.New("already stopped")
	}

	p.started = false
	if p.pstoreWrapper != nil {
		return p.pstoreWrapper.Close()
	}
	return nil
}
