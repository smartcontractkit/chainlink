package coordinator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type topics struct {
	randomnessRequestedTopic            common.Hash
	randomnessFulfillmentRequestedTopic common.Hash
	randomWordsFulfilledTopic           common.Hash
	configSetTopic                      common.Hash
	newTransmissionTopic                common.Hash
}

func newTopics() (topics, error) {
	requestedEvent, ok := vrfABI.Events[randomnessRequestedEvent]
	if !ok {
		return topics{}, fmt.Errorf("could not find event %s in vrfABI %+v", randomnessRequestedEvent, vrfABI.Events)
	}

	fulfillmentRequestedEvent, ok := vrfABI.Events[randomnessFulfillmentRequestedEvent]
	if !ok {
		return topics{}, fmt.Errorf("could not find event %s in vrfABI %+v", randomnessFulfillmentRequestedEvent, vrfABI.Events)
	}

	fulfilledEvent, ok := vrfABI.Events[randomWordsFulfilledEvent]
	if !ok {
		return topics{}, fmt.Errorf("could not find event %s in vrfABI %+v", randomWordsFulfilledEvent, vrfABI.Events)
	}

	transmissionEvent, ok := vrfABI.Events[newTransmissionEvent]
	if !ok {
		return topics{}, fmt.Errorf("could not find event %s in vrfABI %+v", newTransmissionEvent, vrfABI.Events)
	}

	configSet, ok := vrfABI.Events[configSetEvent]
	if !ok {
		return topics{}, fmt.Errorf("could not find event %s in vrfABI %+v", configSetEvent, vrfABI.Events)
	}

	dkgConfigSet, ok := dkgABI.Events[configSetEvent]
	if !ok {
		return topics{}, fmt.Errorf("could not find event %s in dkgABI %+v", configSetEvent, dkgABI.Events)
	}

	// DKG set config and VRF set config should be equal
	if dkgConfigSet.ID != configSet.ID {
		return topics{}, fmt.Errorf(
			"invariant violation: dkg ConfigSet topic (%s) != vrf ConfigSet topic (%s)",
			dkgConfigSet.ID.Hex(), configSet.ID.Hex())
	}

	return topics{
		randomnessRequestedTopic:            requestedEvent.ID,
		randomnessFulfillmentRequestedTopic: fulfillmentRequestedEvent.ID,
		randomWordsFulfilledTopic:           fulfilledEvent.ID,
		configSetTopic:                      configSet.ID,
		newTransmissionTopic:                transmissionEvent.ID,
	}, nil
}
