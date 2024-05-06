package loadvrfv2plus

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	vrfcommon "github.com/smartcontractkit/chainlink/e2e-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/e2e-tests/actions/vrf/vrfv2plus"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/e2e-tests/testconfig/vrfv2plus"
)

type BHSTestGun struct {
	contracts  *vrfcommon.VRFContracts
	keyHash    [32]byte
	subIDs     []*big.Int
	testConfig *vrfv2plus_config.Config
	logger     zerolog.Logger
}

func NewBHSTestGun(
	contracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	subIDs []*big.Int,
	testConfig *vrfv2plus_config.Config,
	logger zerolog.Logger,
) *BHSTestGun {
	return &BHSTestGun{
		contracts:  contracts,
		subIDs:     subIDs,
		keyHash:    keyHash,
		testConfig: testConfig,
		logger:     logger,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *BHSTestGun) Call(_ *wasp.Generator) *wasp.Response {
	vrfv2PlusConfig := m.testConfig.General
	billingType, err := selectBillingType(*vrfv2PlusConfig.SubscriptionBillingType)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	_, err = vrfv2plus.RequestRandomness(
		//the same consumer is used for all requests and in all subs
		m.contracts.VRFV2PlusConsumer[0],
		m.contracts.CoordinatorV2Plus,
		&vrfcommon.VRFKeyData{KeyHash: m.keyHash},
		//randomly pick a subID from pool of subIDs
		m.subIDs[randInRange(0, len(m.subIDs)-1)],
		billingType,
		vrfv2PlusConfig,
		m.logger,
	)
	//todo - might need to store randRequestBlockNumber and blockhash to verify that it was stored in BHS contract at the end of the test
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	return &wasp.Response{}
}

/* SingleHashGun is a gun that constantly requests randomness for one feed  */
type SingleHashGun struct {
	contracts  *vrfcommon.VRFContracts
	keyHash    [32]byte
	subIDs     []*big.Int
	testConfig *vrfv2plus_config.Config
	logger     zerolog.Logger
}

func NewSingleHashGun(
	contracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	subIDs []*big.Int,
	testConfig *vrfv2plus_config.Config,
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
	vrfv2PlusConfig := m.testConfig.General
	billingType, err := selectBillingType(*vrfv2PlusConfig.SubscriptionBillingType)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}

	//randomly increase/decrease randomness request count per TX
	reqCount := deviateValue(*m.testConfig.General.RandomnessRequestCountPerRequest, *m.testConfig.General.RandomnessRequestCountPerRequestDeviation)
	m.testConfig.General.RandomnessRequestCountPerRequest = &reqCount
	_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
		//the same consumer is used for all requests and in all subs
		m.contracts.VRFV2PlusConsumer[0],
		m.contracts.CoordinatorV2Plus,
		&vrfcommon.VRFKeyData{KeyHash: m.keyHash},
		//randomly pick a subID from pool of subIDs
		m.subIDs[randInRange(0, len(m.subIDs)-1)],
		billingType,
		vrfv2PlusConfig,
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
