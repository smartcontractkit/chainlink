package loadvrfv2

import (
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	vrfConst "github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/wasp"
)

/* SingleHashGun is a gun that constantly requests randomness for one feed  */

type SingleHashGun struct {
	contracts *vrfv2_actions.VRFV2Contracts
	keyHash   [32]byte
}

func SingleFeedGun(contracts *vrfv2_actions.VRFV2Contracts, keyHash [32]byte) *SingleHashGun {
	return &SingleHashGun{
		contracts: contracts,
		keyHash:   keyHash,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleHashGun) Call(l *wasp.Generator) *wasp.CallResult {
	err := m.contracts.LoadTestConsumer.RequestRandomness(
		m.keyHash,
		vrfConst.SubID,
		vrfConst.MinimumConfirmations,
		vrfConst.CallbackGasLimit,
		vrfConst.NumberOfWords,
		vrfConst.RandomnessRequestCountPerRequest,
	)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	return &wasp.CallResult{}
}
