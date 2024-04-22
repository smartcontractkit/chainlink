package types

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	MethodRegisterTrigger   = "RegisterTrigger"
	MethodUnRegisterTrigger = "UnregisterTrigger"
	MethodTriggerEvent      = "TriggerEvent"
)

//go:generate mockery --quiet --name Dispatcher --output ./mocks/ --case=underscore
type Dispatcher interface {
	SetReceiver(capabilityId string, donId string, receiver Receiver) error
	RemoveReceiver(capabilityId string, donId string)
	Send(peerID p2ptypes.PeerID, msgBody *MessageBody) error
}

//go:generate mockery --quiet --name Receiver --output ./mocks/ --case=underscore
type Receiver interface {
	Receive(msg *MessageBody)
}

type Aggregator interface {
	Aggregate(eventID string, responses [][]byte) (commoncap.CapabilityResponse, error)
}

// NOTE: this type will become part of the Registry (KS-108)
type DON struct {
	ID      string
	Members []p2ptypes.PeerID
	F       uint8
}
