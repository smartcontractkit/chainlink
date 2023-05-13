package ocr2vrf_constants

import (
	"math/big"
	"time"

	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
)

var (
	LinkEthFeedResponse                          = big.NewInt(1e18)
	LinkFundingAmount                            = big.NewInt(100)
	BeaconPeriodBlocksCount                      = big.NewInt(3)
	EthFundingAmount                             = big.NewFloat(1)
	NumberOfRandomWordsToRequest                 = uint16(2)
	ConfirmationDelay                            = big.NewInt(1)
	RandomnessFulfilmentTransmissionEventTimeout = time.Minute * 6
	RandomnessRedeemTransmissionEventTimeout     = time.Minute * 5
	//keyId can be any random value
	KeyID = "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc"

	CoordinatorConfig = &ocr2vrftypes.CoordinatorConfig{
		CacheEvictionWindowSeconds: 60,
		BatchGasLimit:              5_000_000,
		CoordinatorOverhead:        50_000,
		CallbackOverhead:           50_000,
		BlockGasOverhead:           50_000,
		LookbackBlocks:             1_000,
	}
	VRFBeaconAllowedConfirmationDelays = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
)
