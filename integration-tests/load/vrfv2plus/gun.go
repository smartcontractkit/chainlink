package loadvrfv2plus

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/wasp"
	"math/big"
)

/* SingleHashGun is a gun that constantly requests randomness for one feed  */

type SingleHashGun struct {
	contracts       *vrfv2plus.VRFV2_5Contracts
	keyHash         [32]byte
	subID           *big.Int
	vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig
	logger          zerolog.Logger
}

func NewSingleHashGun(
	contracts *vrfv2plus.VRFV2_5Contracts,
	keyHash [32]byte,
	subID *big.Int,
	vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig,
	logger zerolog.Logger,
) *SingleHashGun {
	return &SingleHashGun{
		contracts:       contracts,
		keyHash:         keyHash,
		subID:           subID,
		vrfv2PlusConfig: vrfv2PlusConfig,
		logger:          logger,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleHashGun) Call(l *wasp.Generator) *wasp.CallResult {
	//todo - should work with multiple consumers and consumers having different keyhashes and wallets
	_, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
		m.contracts.LoadTestConsumers[0],
		m.contracts.Coordinator,
		&vrfv2plus.VRFV2PlusData{VRFV2PlusKeyData: vrfv2plus.VRFV2PlusKeyData{KeyHash: m.keyHash}},
		m.subID,
		//todo - make this configurable
		m.vrfv2PlusConfig.IsNativePayment,
		m.vrfv2PlusConfig,
		m.logger,
	)

	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}

	return &wasp.CallResult{}
}
