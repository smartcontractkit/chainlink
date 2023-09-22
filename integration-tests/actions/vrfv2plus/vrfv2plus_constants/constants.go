package vrfv2plus_constants

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
)

var (
	LinkEthFeedResponse                     = big.NewInt(1e18)
	MinimumConfirmations                    = uint16(3)
	RandomnessRequestCountPerRequest        = uint16(1)
	VRFSubscriptionFundingAmountLink        = big.NewInt(100)
	VRFSubscriptionFundingAmountNativeToken = big.NewInt(1)
	ChainlinkNodeFundingAmountEth           = big.NewFloat(0.1)
	NumberOfWords                           = uint32(3)
	CallbackGasLimit                        = uint32(1000000)
	MaxGasLimitVRFCoordinatorConfig         = uint32(2.5e6)
	StalenessSeconds                        = uint32(86400)
	GasAfterPaymentCalculation              = uint32(33825)

	VRFCoordinatorV2PlusFeeConfig = vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig{
		FulfillmentFlatFeeLinkPPM:   500,
		FulfillmentFlatFeeNativePPM: 500,
	}

	VRFCoordinatorV2PlusUpgradedVersionFeeConfig = vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionFeeConfig{
		FulfillmentFlatFeeLinkPPM:   500,
		FulfillmentFlatFeeNativePPM: 500,
	}
)
