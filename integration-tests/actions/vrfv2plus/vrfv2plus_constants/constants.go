package vrfv2plus_constants

import (
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"math/big"
)

var (
	LinkEthFeedResponse                     = big.NewInt(1e18)
	MinimumConfirmations                    = uint16(3)
	RandomnessRequestCountPerRequest        = uint16(1)
	VRFSubscriptionFundingAmountLink        = big.NewInt(10)
	VRFSubscriptionFundingAmountNativeToken = big.NewInt(1)
	ChainlinkNodeFundingAmountEth           = big.NewFloat(0.1)
	NumberOfWords                           = uint32(3)
	CallbackGasLimit                        = uint32(1000000)
	MaxGasLimitVRFCoordinatorConfig         = uint32(2.5e6)
	StalenessSeconds                        = uint32(86400)
	GasAfterPaymentCalculation              = uint32(33825)

	VRFCoordinatorV2PlusFeeConfig = vrf_coordinator_v2plus.VRFCoordinatorV2PlusFeeConfig{
		FulfillmentFlatFeeLinkPPM: 500,
		FulfillmentFlatFeeEthPPM:  500,
	}

	VRFCoordinatorV2PlusUpgradedVersionFeeConfig = vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionFeeConfig{
		FulfillmentFlatFeeLinkPPM: 500,
		FulfillmentFlatFeeEthPPM:  500,
	}

	WrapperGasOverhead                      = uint32(50_000)
	CoordinatorGasOverhead                  = uint32(52_000)
	WrapperPremiumPercentage                = uint8(25)
	WrapperMaxNumberOfWords                 = uint8(10)
	WrapperConsumerFundingAmountNativeToken = big.NewFloat(1)

	WrapperConsumerFundingAmountLink = big.NewInt(10)
)
