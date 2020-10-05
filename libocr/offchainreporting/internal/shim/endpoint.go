package shim

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/serialization"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
)

type SerializingEndpoint struct {
	endpoint     types.BinaryNetworkEndpoint
	logger       types.Logger
	mutex        sync.Mutex
	subprocesses subprocesses.Subprocesses
	started      bool
	closed       bool
	closedChOut  bool
	chCancel     chan struct{}
	chOut        chan protocol.MessageWithSender
}

var _ protocol.NetworkEndpoint = (*SerializingEndpoint)(nil)

func NewSerializingEndpoint(
	endpoint types.BinaryNetworkEndpoint,
	logger types.Logger,
) *SerializingEndpoint {
	return &SerializingEndpoint{
		endpoint,
		logger,
		sync.Mutex{},
		subprocesses.Subprocesses{},
		false,
		false,
		false,
		make(chan struct{}),
		make(chan protocol.MessageWithSender),
	}
}

func (n *SerializingEndpoint) serialize(msg protocol.Message) []byte {
	sMsg, err := serialization.Serialize(msg)
	if err != nil {
		n.logger.Error("SerializingEndpoint: Failed to serialize", types.LogFields{
			"message": msg,
		})
		return nil
	}
	return sMsg
}

func (n *SerializingEndpoint) Start() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.started {
		panic("Cannot start already started SerializingEndpoint")
	}
	n.started = true

	if err := n.endpoint.Start(); err != nil {
		return errors.Wrap(err, "while starting SerializingEndpoint")
	}

	n.subprocesses.Go(func() {
		chRaw := n.endpoint.Receive()
		for {
			select {
			case msg, ok := <-chRaw:
				if !ok {
					n.mutex.Lock()
					defer n.mutex.Unlock()
					n.closedChOut = true
					close(n.chOut)
					return
				}
				deserialized, err := serialization.Deserialize(msg.Msg)
				if err != nil {
					n.logger.Error("SerializingEndpoint: Failed to deserialize", types.LogFields{
						"message": msg,
					})
				}
				select {
				case n.chOut <- protocol.MessageWithSender{deserialized, msg.Sender}:
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

func (n *SerializingEndpoint) SendTo(msg protocol.Message, to types.OracleID) {
	sMsg := n.serialize(msg)
	if sMsg != nil {
		n.endpoint.SendTo(sMsg, to)
	}
}

func (n *SerializingEndpoint) Broadcast(msg protocol.Message) {
	sMsg := n.serialize(msg)
	if sMsg != nil {
		n.endpoint.Broadcast(sMsg)
	}
}

func (n *SerializingEndpoint) Receive() <-chan protocol.MessageWithSender {
	return n.chOut
}
