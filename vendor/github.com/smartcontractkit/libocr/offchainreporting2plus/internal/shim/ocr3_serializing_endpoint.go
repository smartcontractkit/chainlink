// Package shim contains implementations of internal types in terms of the external types
package shim

import (
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type OCR3SerializingEndpoint[RI any] struct {
	chTelemetry  chan<- *serialization.TelemetryWrapper
	configDigest types.ConfigDigest
	endpoint     commontypes.BinaryNetworkEndpoint
	maxSigLen    int
	logger       commontypes.Logger
	pluginLimits ocr3types.ReportingPluginLimits
	n, f         int

	mutex        sync.Mutex
	subprocesses subprocesses.Subprocesses
	started      bool
	closed       bool
	closedChOut  bool
	chCancel     chan struct{}
	chOut        chan protocol.MessageWithSender[RI]
	taper        loghelper.LogarithmicTaper
}

var _ protocol.NetworkEndpoint[struct{}] = (*OCR3SerializingEndpoint[struct{}])(nil)

func NewOCR3SerializingEndpoint[RI any](
	chTelemetry chan<- *serialization.TelemetryWrapper,
	configDigest types.ConfigDigest,
	endpoint commontypes.BinaryNetworkEndpoint,
	maxSigLen int,
	logger commontypes.Logger,
	pluginLimits ocr3types.ReportingPluginLimits,
	n, f int,
) *OCR3SerializingEndpoint[RI] {
	return &OCR3SerializingEndpoint[RI]{
		chTelemetry,
		configDigest,
		endpoint,
		maxSigLen,
		logger,
		pluginLimits,
		n, f,

		sync.Mutex{},
		subprocesses.Subprocesses{},
		false,
		false,
		false,
		make(chan struct{}),
		make(chan protocol.MessageWithSender[RI]),
		loghelper.LogarithmicTaper{},
	}
}

func (n *OCR3SerializingEndpoint[RI]) sendTelemetry(t *serialization.TelemetryWrapper) {
	select {
	case n.chTelemetry <- t:
		n.taper.Reset(func(oldCount uint64) {
			n.logger.Info("OCR3SerializingEndpoint: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		n.taper.Trigger(func(newCount uint64) {
			n.logger.Warn("OCR3SerializingEndpoint: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (n *OCR3SerializingEndpoint[RI]) serialize(msg protocol.Message[RI]) ([]byte, *serialization.MessageWrapper) {
	if !msg.CheckSize(n.n, n.f, n.pluginLimits, n.maxSigLen) {
		n.logger.Error("OCR3SerializingEndpoint: Dropping outgoing message because it fails size check", commontypes.LogFields{
			"limits": n.pluginLimits,
		})
		return nil, nil
	}
	sMsg, pbm, err := serialization.Serialize(msg)
	if err != nil {
		n.logger.Error("OCR3SerializingEndpoint: Failed to serialize", commontypes.LogFields{
			"message": msg,
		})
		return nil, nil
	}
	return sMsg, pbm
}

func (n *OCR3SerializingEndpoint[RI]) deserialize(raw []byte) (protocol.Message[RI], *serialization.MessageWrapper, error) {
	m, pbm, err := serialization.Deserialize[RI](raw)
	if err != nil {
		return nil, nil, err
	}

	if !m.CheckSize(n.n, n.f, n.pluginLimits, n.maxSigLen) {
		return nil, nil, fmt.Errorf("message failed size check")
	}

	return m, pbm, nil
}

// Start starts the SerializingEndpoint. It will also start the underlying endpoint.
func (n *OCR3SerializingEndpoint[RI]) Start() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.started {
		return fmt.Errorf("cannot start already started SerializingEndpoint")
	}
	n.started = true

	if err := n.endpoint.Start(); err != nil {
		return fmt.Errorf("error while starting OCR3SerializingEndpoint: %w", err)
	}

	n.subprocesses.Go(func() {
		chRaw := n.endpoint.Receive()
		for {
			select {
			case raw, ok := <-chRaw:
				if !ok {
					n.mutex.Lock()
					defer n.mutex.Unlock()
					n.closedChOut = true
					close(n.chOut)
					return
				}

				m, pbm, err := n.deserialize(raw.Msg)
				if err != nil {
					n.logger.Error("OCR3SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"message": raw,
						"error":   err,
					})
					n.sendTelemetry(&serialization.TelemetryWrapper{
						Wrapped: &serialization.TelemetryWrapper_AssertionViolation{&serialization.TelemetryAssertionViolation{
							Violation: &serialization.TelemetryAssertionViolation_InvalidSerialization{&serialization.TelemetryAssertionViolationInvalidSerialization{
								ConfigDigest:  n.configDigest[:],
								SerializedMsg: raw.Msg,
								Sender:        uint32(raw.Sender),
							}},
						}},
						UnixTimeNanoseconds: time.Now().UnixNano(),
					})
					break
				}

				n.sendTelemetry(&serialization.TelemetryWrapper{
					Wrapped: &serialization.TelemetryWrapper_MessageReceived{&serialization.TelemetryMessageReceived{
						ConfigDigest: n.configDigest[:],
						Msg:          pbm,
						Sender:       uint32(raw.Sender),
					}},
					UnixTimeNanoseconds: time.Now().UnixNano(),
				})

				select {
				case n.chOut <- protocol.MessageWithSender[RI]{m, raw.Sender}:
				case <-n.chCancel:
					return
				}
			case <-n.chCancel:
				return
			}
		}
	})

	return nil
}

// Close closes the SerializingEndpoint. It will also close the underlying endpoint.
func (n *OCR3SerializingEndpoint[RI]) Close() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.started && !n.closed {
		n.closed = true
		close(n.chCancel)
		n.subprocesses.Wait()

		if !n.closedChOut {
			n.closedChOut = true
			close(n.chOut)
		}

		return n.endpoint.Close()
	}

	return nil
}

func (n *OCR3SerializingEndpoint[RI]) SendTo(msg protocol.Message[RI], to commontypes.OracleID) {
	sMsg, pbm := n.serialize(msg)
	if sMsg != nil {
		n.endpoint.SendTo(sMsg, to)
		n.sendTelemetry(&serialization.TelemetryWrapper{
			Wrapped: &serialization.TelemetryWrapper_MessageSent{&serialization.TelemetryMessageSent{
				ConfigDigest:  n.configDigest[:],
				Msg:           pbm,
				SerializedMsg: sMsg,
				Receiver:      uint32(to),
			}},
			UnixTimeNanoseconds: time.Now().UnixNano(),
		})
	}
}

func (n *OCR3SerializingEndpoint[RI]) Broadcast(msg protocol.Message[RI]) {
	sMsg, pbm := n.serialize(msg)
	if sMsg != nil {
		n.endpoint.Broadcast(sMsg)
		n.sendTelemetry(&serialization.TelemetryWrapper{
			Wrapped: &serialization.TelemetryWrapper_MessageBroadcast{&serialization.TelemetryMessageBroadcast{
				ConfigDigest:  n.configDigest[:],
				Msg:           pbm,
				SerializedMsg: sMsg,
			}},
			UnixTimeNanoseconds: time.Now().UnixNano(),
		})
	}
}

func (n *OCR3SerializingEndpoint[RI]) Receive() <-chan protocol.MessageWithSender[RI] {
	return n.chOut
}
