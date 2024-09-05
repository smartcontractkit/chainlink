package loadvrfv2

import (
	"math/rand"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"

	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	vrfv2_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
)

type BHSTestGun struct {
	contracts  *vrfcommon.VRFContracts
	subIDs     []uint64
	keyHash    [32]byte
	testConfig *vrfv2_config.Config
	logger     zerolog.Logger
	sethClient *seth.Client
}

func NewBHSTestGun(
	contracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	subIDs []uint64,
	testConfig *vrfv2_config.Config,
	logger zerolog.Logger,
	sethClient *seth.Client,
) *BHSTestGun {
	return &BHSTestGun{
		contracts:  contracts,
		subIDs:     subIDs,
		keyHash:    keyHash,
		testConfig: testConfig,
		logger:     logger,
		sethClient: sethClient,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *BHSTestGun) Call(_ *wasp.Generator) *wasp.Response {
	_, err := vrfv2.RequestRandomness(
		m.logger,
		m.contracts.VRFV2Consumers[0],
		m.contracts.CoordinatorV2,
		m.subIDs[0],
		&vrfcommon.VRFKeyData{KeyHash: m.keyHash},
		*m.testConfig.General.MinimumConfirmations,
		*m.testConfig.General.CallbackGasLimit,
		*m.testConfig.General.NumberOfWords,
		*m.testConfig.General.RandomnessRequestCountPerRequest,
		*m.testConfig.General.RandomnessRequestCountPerRequestDeviation,
		seth_utils.AvailableSethKeyNum(m.sethClient),
	)
	//todo - might need to store randRequestBlockNumber and blockhash to verify that it was stored in BHS contract at the end of the test
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	return &wasp.Response{}
}

type SingleHashGun struct {
	contracts  *vrfcommon.VRFContracts
	keyHash    [32]byte
	subIDs     []uint64
	testConfig *vrfv2_config.Config
	logger     zerolog.Logger
	sethClient *seth.Client
}

func NewSingleHashGun(
	contracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	subIDs []uint64,
	testConfig *vrfv2_config.Config,
	logger zerolog.Logger,
	sethClient *seth.Client,
) *SingleHashGun {
	return &SingleHashGun{
		contracts:  contracts,
		keyHash:    keyHash,
		subIDs:     subIDs,
		testConfig: testConfig,
		logger:     logger,
		sethClient: sethClient,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleHashGun) Call(_ *wasp.Generator) *wasp.Response {
	//todo - should work with multiple consumers and consumers having different keyhashes and wallets

	vrfv2Config := m.testConfig.General
	//randomly increase/decrease randomness request count per TX
	randomnessRequestCountPerRequest := deviateValue(*vrfv2Config.RandomnessRequestCountPerRequest, *vrfv2Config.RandomnessRequestCountPerRequestDeviation)
	_, _, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
		m.logger,
		//the same consumer is used for all requests and in all subs
		m.contracts.VRFV2Consumers[0],
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
		seth_utils.AvailableSethKeyNum(m.sethClient),
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
