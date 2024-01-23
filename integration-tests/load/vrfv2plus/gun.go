package loadvrfv2plus

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

/* SingleHashGun is a gun that constantly requests randomness for one feed  */

type SingleHashGun struct {
	contracts  *vrfv2plus.VRFV2_5Contracts
	keyHash    [32]byte
	subIDs     []*big.Int
	testConfig types.VRFv2PlusTestConfig
	logger     zerolog.Logger
}

func NewSingleHashGun(
	contracts *vrfv2plus.VRFV2_5Contracts,
	keyHash [32]byte,
	subIDs []*big.Int,
	testConfig types.VRFv2PlusTestConfig,
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

	billingType, err := selectBillingType(*m.testConfig.GetVRFv2PlusConfig().General.SubscriptionBillingType)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}

	//randomly increase/decrease randomness request count per TX
	randomnessRequestCountPerRequest := deviateValue(*m.testConfig.GetVRFv2PlusConfig().General.RandomnessRequestCountPerRequest, *m.testConfig.GetVRFv2PlusConfig().General.RandomnessRequestCountPerRequestDeviation)
	_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
		//the same consumer is used for all requests and in all subs
		m.contracts.LoadTestConsumers[0],
		m.contracts.Coordinator,
		&vrfv2plus.VRFV2PlusData{VRFV2PlusKeyData: vrfv2plus.VRFV2PlusKeyData{KeyHash: m.keyHash}},
		//randomly pick a subID from pool of subIDs
		m.subIDs[randInRange(0, len(m.subIDs)-1)],
		//randomly pick payment type
		billingType,
		*m.testConfig.GetVRFv2PlusConfig().General.MinimumConfirmations,
		*m.testConfig.GetVRFv2PlusConfig().General.CallbackGasLimit,
		*m.testConfig.GetVRFv2PlusConfig().General.NumberOfWords,
		randomnessRequestCountPerRequest,
		*m.testConfig.GetVRFv2PlusConfig().General.RandomnessRequestCountPerRequestDeviation,
		m.testConfig.GetVRFv2PlusConfig().General.RandomWordsFulfilledEventTimeout.Duration,
		m.logger,
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

func selectBillingType(billingType string) (bool, error) {
	switch vrfv2plus_config.BillingType(billingType) {
	case vrfv2plus_config.BillingType_Link:
		return false, nil
	case vrfv2plus_config.BillingType_Native:
		return true, nil
	case vrfv2plus_config.BillingType_Link_and_Native:
		return randBool(), nil
	default:
		return false, fmt.Errorf("invalid billing type: %s", billingType)
	}
}
