package loadvrfv2

import (
	"math/rand"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

/* SingleHashGun is a gun that constantly requests randomness for one feed  */

type SingleHashGun struct {
	contracts  *vrfcommon.VRFContracts
	keyHash    [32]byte
	subIDs     []uint64
	testConfig types.VRFv2TestConfig
	logger     zerolog.Logger
}

func NewSingleHashGun(
	contracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	subIDs []uint64,
	testConfig types.VRFv2TestConfig,
	logger zerolog.Logger,
) *SingleHashGun {
	return &SingleHashGun{
		contracts:  contracts,
		keyHash:    keyHash,
		subIDs:     subIDs,
		testConfig: testConfig,
		logger:     logger,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleHashGun) Call(_ *wasp.Generator) *wasp.Response {
	//todo - should work with multiple consumers and consumers having different keyhashes and wallets

	vrfv2Config := m.testConfig.GetVRFv2Config().General
	//randomly increase/decrease randomness request count per TX
	randomnessRequestCountPerRequest := deviateValue(*vrfv2Config.RandomnessRequestCountPerRequest, *vrfv2Config.RandomnessRequestCountPerRequestDeviation)
	_, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
		m.logger,
		//the same consumer is used for all requests and in all subs
		m.contracts.VRFV2Consumer[0],
		m.contracts.CoordinatorV2,
		//randomly pick a subID from pool of subIDs
		m.subIDs[randInRange(0, len(m.subIDs)-1)],
		&vrfcommon.VRFKeyData{KeyHash: m.keyHash},
		*vrfv2Config.MinimumConfirmations,
		*vrfv2Config.CallbackGasLimit,
		*vrfv2Config.NumberOfWords,
		randomnessRequestCountPerRequest,
		*vrfv2Config.RandomnessRequestCountPerRequestDeviation,
		vrfv2Config.RandomWordsFulfilledEventTimeout.Duration,
	)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	return &wasp.Response{}
}

func deviateValue(requestCountPerTX uint16, deviation uint16) uint16 {
	if randBool() && requestCountPerTX > deviation {
		requestCountPerTX -= uint16(randInRange(0, int(deviation)))
	} else {
		requestCountPerTX += uint16(randInRange(0, int(deviation)))
	}
	return requestCountPerTX
}

func randBool() bool {
	return rand.Intn(2) == 1
}
func randInRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
