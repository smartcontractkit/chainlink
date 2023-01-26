package coordinator

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
)

type topics struct {
	// VRF logs
	randomnessRequestedTopic            common.Hash
	randomnessFulfillmentRequestedTopic common.Hash
	randomWordsFulfilledTopic           common.Hash
	outputsServedTopic                  common.Hash

	// OCR logs
	newTransmissionTopic common.Hash
	configSetTopic       common.Hash
}

func newTopics() topics {
	return topics{
		randomnessRequestedTopic:            vrf_coordinator.VRFCoordinatorRandomnessRequested{}.Topic(),
		randomnessFulfillmentRequestedTopic: vrf_coordinator.VRFCoordinatorRandomnessFulfillmentRequested{}.Topic(),
		randomWordsFulfilledTopic:           vrf_coordinator.VRFCoordinatorRandomWordsFulfilled{}.Topic(),
		configSetTopic:                      vrf_beacon.VRFBeaconConfigSet{}.Topic(),
		newTransmissionTopic:                vrf_beacon.VRFBeaconNewTransmission{}.Topic(),
		outputsServedTopic:                  vrf_coordinator.VRFCoordinatorOutputsServed{}.Topic(),
	}
}

func (t topics) allSigs() (rv []common.Hash) {
	e := reflect.ValueOf(t)
	for i := 0; i < e.NumField(); i++ {
		rv = append(rv, e.Field(i).Interface().(common.Hash))
	}
	return rv
}

// vrfSigs returns the topics of the logs directly concerned with the operation
// of the VRF service.
func (t topics) vrfSigs() []common.Hash {
	return []common.Hash{
		t.randomnessRequestedTopic,
		t.randomnessFulfillmentRequestedTopic,
		t.randomWordsFulfilledTopic,
		t.outputsServedTopic,
	}
}
