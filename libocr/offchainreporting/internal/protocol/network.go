package protocol

import (
	"log"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type NetworkSender interface {
		SendTo(msg Message, to types.OracleID)
		Broadcast(msg Message)
}

type NetworkEndpoint interface {
	NetworkSender
		Receive() <-chan MessageWithSender

			Start() error

			Close() error
}

type SimpleNetwork struct {
	chs []chan MessageWithSender }

func NewSimpleNetwork(n int) *SimpleNetwork {
	s := SimpleNetwork{}
	for i := 0; i < n; i++ {
		s.chs = append(s.chs, make(chan MessageWithSender, 100))
	}
	return &s
}

func (net *SimpleNetwork) Endpoint(id types.OracleID) (NetworkEndpoint, error) {
	return SimpleNetworkEndpoint{
		net,
		id,
	}, nil
}

type SimpleNetworkEndpoint struct {
	net *SimpleNetwork 	id  types.OracleID }

var _ NetworkEndpoint = (*SimpleNetworkEndpoint)(nil)

func (end SimpleNetworkEndpoint) SendTo(msg Message, to types.OracleID) {
	log.Printf("[%v] sending to %v: %T\n", end.id, to, msg)
	end.net.chs[to] <- MessageWithSender{msg, end.id}
}

func (end SimpleNetworkEndpoint) Broadcast(msg Message) {
	log.Printf("[%v] broadcasting: %T\n", end.id, msg)
	for _, ch := range end.net.chs {
		ch <- MessageWithSender{msg, end.id}
	}
}

func (end SimpleNetworkEndpoint) Receive() <-chan MessageWithSender {
	return end.net.chs[end.id]
}

func (SimpleNetworkEndpoint) Start() error { return nil }

func (SimpleNetworkEndpoint) Close() error { return nil }
