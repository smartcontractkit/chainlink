// Package types re-exports the DON-to-DON dispatcher types now living in
// capabilities/libs/x/don2don/types, so callers that only need the type/interface definitions (not
// the dispatcher implementation itself, which moved) can keep importing this path unchanged.
package types

import (
	don2dontypes "github.com/smartcontractkit/capabilities/libs/x/don2don/types"
)

const (
	MethodRegisterTrigger          = don2dontypes.MethodRegisterTrigger
	MethodUnregisterTrigger        = don2dontypes.MethodUnregisterTrigger
	MethodTriggerRegistrationCheck = don2dontypes.MethodTriggerRegistrationCheck
	MethodTriggerEvent             = don2dontypes.MethodTriggerEvent
	MethodExecute                  = don2dontypes.MethodExecute
	MethodTriggerEventAck          = don2dontypes.MethodTriggerEventAck
)

type (
	Dispatcher      = don2dontypes.Dispatcher
	Receiver        = don2dontypes.Receiver
	ReceiverService = don2dontypes.ReceiverService
	Aggregator      = don2dontypes.Aggregator
	DON             = don2dontypes.DON
	MessageHasher   = don2dontypes.MessageHasher

	MessageBody                      = don2dontypes.MessageBody
	MessageBody_TriggerEventMetadata = don2dontypes.MessageBody_TriggerEventMetadata
	TriggerEventMetadata             = don2dontypes.TriggerEventMetadata
)
