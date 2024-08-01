package networking

import (
	"fmt"
	"io"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"

	"github.com/smartcontractkit/libocr/internal/loghelper"
)

var (
	_ commontypes.BinaryNetworkEndpoint = &ocrEndpointV2{}
)

type ocrEndpointState int

const (
	_ ocrEndpointState = iota
	ocrEndpointUnstarted
	ocrEndpointStarted
	ocrEndpointClosed
)

type EndpointConfigV2 struct {
	// IncomingMessageBufferSize is the per-remote number of incoming
	// messages to buffer. Any additional messages received on top of those
	// already in the queue will be dropped.
	IncomingMessageBufferSize int

	// OutgoingMessageBufferSize is the per-remote number of outgoing
	// messages to buffer. Any additional messages send on top of those
	// already in the queue will displace the oldest.
	// NOTE: OutgoingMessageBufferSize should be comfortably smaller than remote's
	// IncomingMessageBufferSize to give the remote enough space to process
	// them all in case we regained connection and now send a bunch at once
	OutgoingMessageBufferSize int
}

// ocrEndpointV2 represents a member of a particular feed oracle group
type ocrEndpointV2 struct {
	// configuration and settings
	config              EndpointConfigV2
	peerIDs             []ragetypes.PeerID
	peerMapping         map[commontypes.OracleID]ragetypes.PeerID
	reversedPeerMapping map[ragetypes.PeerID]commontypes.OracleID
	peer                *concretePeerV2
	host                *ragep2p.Host
	configDigest        ocr2types.ConfigDigest
	bootstrappers       []ragetypes.PeerInfo
	f                   int
	ownOracleID         commontypes.OracleID

	// internal and state management
	chSendToSelf chan commontypes.BinaryMessageWithSender
	chClose      chan struct{}
	streams      map[commontypes.OracleID]*ragep2p.Stream
	registration io.Closer
	state        ocrEndpointState

	stateMu sync.RWMutex
	subs    subprocesses.Subprocesses

	// recv is exposed to clients of this network endpoint
	recv chan commontypes.BinaryMessageWithSender

	logger loghelper.LoggerWithContext

	limits BinaryNetworkEndpointLimits
}

func reverseMappingV2(m map[commontypes.OracleID]ragetypes.PeerID) map[ragetypes.PeerID]commontypes.OracleID {
	n := make(map[ragetypes.PeerID]commontypes.OracleID)
	for k, v := range m {
		n[v] = k
	}
	return n
}

// sendToSelfBufferSize is how many messages we will keep in memory that
// are sent to ourself before we start dropping
const sendToSelfBufferSize = 20

func newOCREndpointV2(
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeerV2,
	peerIDs []ragetypes.PeerID,
	v2bootstrappers []ragetypes.PeerInfo,
	config EndpointConfigV2,
	f int,
	limits BinaryNetworkEndpointLimits,
	registration io.Closer,
) (*ocrEndpointV2, error) {
	peerMapping := make(map[commontypes.OracleID]ragetypes.PeerID)
	for i, peerID := range peerIDs {
		peerMapping[commontypes.OracleID(i)] = peerID
	}
	reversedPeerMapping := reverseMappingV2(peerMapping)
	ownOracleID, ok := reversedPeerMapping[peer.peerID]
	if !ok {
		return nil, fmt.Errorf("host peer ID %s is not present in given peerMapping", peer.PeerID())
	}

	chSendToSelf := make(chan commontypes.BinaryMessageWithSender, sendToSelfBufferSize)

	logger = logger.MakeChild(commontypes.LogFields{
		"configDigest": configDigest.Hex(),
		"oracleID":     ownOracleID,
		"id":           "OCREndpointV2",
	})

	logger.Info("OCREndpointV2: Initialized", commontypes.LogFields{
		"bootstrappers": v2bootstrappers,
		"oracles":       peerIDs,
	})

	return &ocrEndpointV2{
		config,
		peerIDs,
		peerMapping,
		reversedPeerMapping,
		peer,
		peer.host,
		configDigest,
		v2bootstrappers,
		f,
		ownOracleID,
		chSendToSelf,
		make(chan struct{}),
		make(map[commontypes.OracleID]*ragep2p.Stream),
		registration,
		ocrEndpointUnstarted,
		sync.RWMutex{},
		subprocesses.Subprocesses{},
		make(chan commontypes.BinaryMessageWithSender),
		logger,
		limits,
	}, nil
}

func streamNameFromConfigDigest(cd ocr2types.ConfigDigest) string {
	return fmt.Sprintf("ocr/%s", cd)
}

// Start the ocrEndpointV2. Should only be called once.
func (o *ocrEndpointV2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			o.Close()
		}
	}()

	o.stateMu.Lock()
	defer o.stateMu.Unlock()

	if o.state != ocrEndpointUnstarted {
		return fmt.Errorf("cannot start ocrEndpointV2 that is not unstarted, state was: %d", o.state)
	}
	o.state = ocrEndpointStarted

	for oid, pid := range o.peerMapping {
		if oid == o.ownOracleID {
			continue
		}
		streamName := streamNameFromConfigDigest(o.configDigest)
		stream, err := o.host.NewStream(
			pid,
			streamName,
			o.config.OutgoingMessageBufferSize,
			o.config.IncomingMessageBufferSize,
			o.limits.MaxMessageLength,
			ragep2p.TokenBucketParams{
				o.limits.MessagesRatePerOracle,
				uint32(o.limits.MessagesCapacityPerOracle),
			},
			ragep2p.TokenBucketParams{
				o.limits.BytesRatePerOracle,
				uint32(o.limits.BytesCapacityPerOracle),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create stream for oracle %v (peer id: %q): %w", oid, pid, err)
		}
		o.streams[oid] = stream
	}

	for oid := range o.streams {
		oid := oid
		o.subs.Go(func() {
			o.runRecv(oid)
		})
	}
	o.subs.Go(func() {
		o.runSendToSelf()
	})

	o.logger.Info("OCREndpointV2: Started listening", nil)
	succeeded = true
	return nil
}

// Receive runloop is per-remote
// This means that each remote gets its own buffered channel, so even if one
// remote goes mad and sends us thousands of messages, we don't drop any
// messages from good remotes
func (o *ocrEndpointV2) runRecv(oid commontypes.OracleID) {
	chRecv := o.streams[oid].ReceiveMessages()
	for {
		select {
		case payload := <-chRecv:
			msg := commontypes.BinaryMessageWithSender{
				Msg:    payload,
				Sender: oid,
			}
			select {
			case o.recv <- msg:
				continue
			case <-o.chClose:
				return
			}
		case <-o.chClose:
			return
		}
	}
}

func (o *ocrEndpointV2) runSendToSelf() {
	for {
		select {
		case <-o.chClose:
			return
		case m := <-o.chSendToSelf:
			select {
			case o.recv <- m:
			case <-o.chClose:
				return
			}
		}
	}
}

// Close should be called to clean up even if Start is never called.
func (o *ocrEndpointV2) Close() error {
	o.stateMu.Lock()
	defer o.stateMu.Unlock()
	if o.state != ocrEndpointStarted {
		return fmt.Errorf("cannot close ocrEndpointV2 that is not started, state was: %d", o.state)
	}
	o.state = ocrEndpointClosed

	o.logger.Debug("OCREndpointV2: Closing", nil)

	o.logger.Debug("OCREndpointV2: Closing streams", nil)
	close(o.chClose)
	o.subs.Wait()

	var allErrors error
	for oid, stream := range o.streams {
		if err := stream.Close(); err != nil {
			allErrors = multierr.Append(allErrors, fmt.Errorf("error while closing stream with oracle %v: %w", oid, err))
		}
	}

	o.logger.Debug("OCREndpointV2: Deregister", nil)
	if err := o.registration.Close(); err != nil {
		allErrors = multierr.Append(allErrors, fmt.Errorf("error closing OCREndpointV2: could not deregister: %w", err))
	}

	o.logger.Debug("OCREndpointV2: Closing o.recv", nil)
	close(o.recv)

	o.logger.Info("OCREndpointV2: Closed", nil)
	return allErrors
}

// SendTo sends a message to the given oracle
// It makes a best effort delivery. If stream is unavailable for any
// reason, it will fill up to outgoingMessageBufferSize then drop messages
// until the stream becomes available again
//
// NOTE: If a stream connection is lost, the buffer will keep only the newest
// messages and drop older ones until the stream opens again.
func (o *ocrEndpointV2) SendTo(payload []byte, to commontypes.OracleID) {
	o.stateMu.RLock()
	state := o.state
	o.stateMu.RUnlock()
	if state != ocrEndpointStarted {
		o.logger.Error("Send on non-started ocrEndpointV2", commontypes.LogFields{"state": state})
		return
	}

	if to == o.ownOracleID {
		o.sendToSelf(payload)
		return
	}

	o.streams[to].SendMessage(payload)
}

func (o *ocrEndpointV2) sendToSelf(payload []byte) {
	m := commontypes.BinaryMessageWithSender{
		payload,
		o.ownOracleID,
	}

	select {
	case o.chSendToSelf <- m:
	default:
		o.logger.Error("Send-to-self buffer is full, dropping message", commontypes.LogFields{
			"remoteOracleID": o.ownOracleID,
		})
	}
}

// Broadcast sends a msg to all oracles in the peer mapping
func (o *ocrEndpointV2) Broadcast(payload []byte) {
	var subs subprocesses.Subprocesses
	defer subs.Wait()
	for oracleID := range o.peerMapping {
		oracleID := oracleID
		subs.Go(func() {
			o.SendTo(payload, oracleID)
		})
	}
}

// Receive gives the channel to receive messages
func (o *ocrEndpointV2) Receive() <-chan commontypes.BinaryMessageWithSender {
	return o.recv
}
