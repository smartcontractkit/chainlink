// Package shim contains implementations of internal types in terms of the external types
package shim

import (
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type SerializingEndpoint struct {
	chTelemetry           chan<- *serialization.TelemetryWrapper
	configDigest          types.ConfigDigest
	endpoint              commontypes.BinaryNetworkEndpoint
	logger                commontypes.Logger
	reportingPluginLimits types.ReportingPluginLimits

	mutex        sync.Mutex
	subprocesses subprocesses.Subprocesses
	started      bool
	closed       bool
	closedChOut  bool
	chCancel     chan struct{}
	chOut        chan protocol.MessageWithSender
	taper        loghelper.LogarithmicTaper
}

var _ protocol.NetworkEndpoint = (*SerializingEndpoint)(nil)

func NewSerializingEndpoint(
	chTelemetry chan<- *serialization.TelemetryWrapper,
	configDigest types.ConfigDigest,
	endpoint commontypes.BinaryNetworkEndpoint,
	logger commontypes.Logger,
	reportingPluginLimits types.ReportingPluginLimits,
) *SerializingEndpoint {
	return &SerializingEndpoint{
		chTelemetry,
		configDigest,
		endpoint,
		logger,
		reportingPluginLimits,

		sync.Mutex{},
		subprocesses.Subprocesses{},
		false,
		false,
		false,
		make(chan struct{}),
		make(chan protocol.MessageWithSender),
		loghelper.LogarithmicTaper{},
	}
}

func (n *SerializingEndpoint) sendTelemetry(t *serialization.TelemetryWrapper) {
	select {
	case n.chTelemetry <- t:
		n.taper.Reset(func(oldCount uint64) {
			n.logger.Info("SerializingEndpoint: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		n.taper.Trigger(func(newCount uint64) {
			n.logger.Warn("SerializingEndpoint: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (n *SerializingEndpoint) serialize(msg protocol.Message) ([]byte, *serialization.MessageWrapper) {
	if !msg.CheckSize(n.reportingPluginLimits) {
		n.logger.Error("SerializingEndpoint: Dropping outgoing message because it fails size check", commontypes.LogFields{
			"message": msg,
			"limits":  n.reportingPluginLimits,
		})
		return nil, nil
	}
	sMsg, pbm, err := serialization.Serialize(msg)
	if err != nil {
		n.logger.Error("SerializingEndpoint: Failed to serialize", commontypes.LogFields{
			"message": msg,
		})
		return nil, nil
	}
	return sMsg, pbm
}

func (n *SerializingEndpoint) deserialize(raw []byte) (protocol.Message, *serialization.MessageWrapper, error) {
	m, pbm, err := serialization.Deserialize(raw)
	if err != nil {
		return nil, nil, err
	}

	if !m.CheckSize(n.reportingPluginLimits) {
		return nil, nil, fmt.Errorf("message failed size check")
	}

	return m, pbm, nil
}

// Start starts the SerializingEndpoint. It will also start the underlying endpoint.
func (n *SerializingEndpoint) Start() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.started {
		return fmt.Errorf("cannot start already started SerializingEndpoint")
	}
	n.started = true

	if err := n.endpoint.Start(); err != nil {
		return fmt.Errorf("error while starting SerializingEndpoint: %w", err)
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
					n.logger.Error("SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"message": raw,
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
				case n.chOut <- protocol.MessageWithSender{m, raw.Sender}:
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
func (n *SerializingEndpoint) Close() error {
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

func (n *SerializingEndpoint) SendTo(msg protocol.Message, to commontypes.OracleID) {
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

func (n *SerializingEndpoint) Broadcast(msg protocol.Message) {
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

func (n *SerializingEndpoint) Receive() <-chan protocol.MessageWithSender {
	return n.chOut
}
