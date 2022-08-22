package coordinator

import (
	"github.com/ethereum/go-ethereum/common"

	vrf_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_coordinator"
)

type topics struct {
	randomnessRequestedTopic            common.Hash
	randomnessFulfillmentRequestedTopic common.Hash
	randomWordsFulfilledTopic           common.Hash
	configSetTopic                      common.Hash
	newTransmissionTopic                common.Hash
}

func newTopics() topics {
	return topics{
		randomnessRequestedTopic:            vrf_wrapper.VRFBeaconCoordinatorRandomnessRequested{}.Topic(),
		randomnessFulfillmentRequestedTopic: vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested{}.Topic(),
		randomWordsFulfilledTopic:           vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled{}.Topic(),
		configSetTopic:                      vrf_wrapper.VRFBeaconCoordinatorConfigSet{}.Topic(),
		newTransmissionTopic:                vrf_wrapper.VRFBeaconCoordinatorNewTransmission{}.Topic(),
	}
}
