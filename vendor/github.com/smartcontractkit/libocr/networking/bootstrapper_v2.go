package networking

import (
	"fmt"
	"io"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/libocr/internal/loghelper"
)

var (
	_ commontypes.Bootstrapper = &bootstrapperV2{}
)

type bootstrapperState int

const (
	_ bootstrapperState = iota
	bootstrapperUnstarted
	bootstrapperStarted
	bootstrapperClosed
)

type bootstrapperV2 struct {
	peer            *concretePeerV2
	v2peerIDs       []ragetypes.PeerID
	v2bootstrappers []ragetypes.PeerInfo
	logger          loghelper.LoggerWithContext
	configDigest    ocr2types.ConfigDigest
	registration    io.Closer
	state           bootstrapperState

	stateMu *sync.Mutex
	f       int
}

func newBootstrapperV2(
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeerV2,
	v2peerIDs []ragetypes.PeerID,
	v2bootstrappers []ragetypes.PeerInfo,
	f int,
	registration io.Closer,
) (*bootstrapperV2, error) {
	logger = logger.MakeChild(commontypes.LogFields{
		"id":           "bootstrapperV2",
		"configDigest": configDigest.Hex(),
	})

	logger.Info("BootstrapperV2: Initialized", commontypes.LogFields{
		"bootstrappers": v2bootstrappers,
		"oracles":       v2peerIDs,
	})

	return &bootstrapperV2{
		peer,
		v2peerIDs,
		v2bootstrappers,
		logger,
		configDigest,
		registration,
		bootstrapperUnstarted,
		new(sync.Mutex),
		f,
	}, nil
}

func (b *bootstrapperV2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			b.Close()
		}
	}()

	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		return fmt.Errorf("cannot start bootstrapperV2 that is not unstarted, state was: %d", b.state)
	}

	b.state = bootstrapperStarted

	b.logger.Info("BootstrapperV2: Started listening", nil)
	succeeded = true
	return nil
}

func (b *bootstrapperV2) Close() error {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	if b.state != bootstrapperStarted {
		return fmt.Errorf("cannot close bootstrapperV2 that is not started, state was: %d", b.state)
	}
	b.state = bootstrapperClosed

	if err := b.registration.Close(); err != nil {
		return fmt.Errorf("could not unregister bootstrapperV2: %w", err)
	}
	return nil
}
