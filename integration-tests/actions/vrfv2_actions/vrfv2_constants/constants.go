package vrfv2_constants

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
)

var (
	LinkEthFeedResponse              = big.NewInt(1e18)
	MinimumConfirmations             = uint16(3)
	RandomnessRequestCountPerRequest = uint16(1)
	//todo - get Sub id when creating subscription - need to listen for SubscriptionCreated Log
	SubID                            = uint64(1)
	VRFSubscriptionFundingAmountLink = big.NewInt(100)
	ChainlinkNodeFundingAmountEth    = big.NewFloat(0.1)
	NumberOfWords                    = uint32(3)
	MaxGasPriceGWei                  = 1000
	CallbackGasLimit                 = uint32(1000000)
	MaxGasLimitVRFCoordinatorConfig  = uint32(2.5e6)
	StalenessSeconds                 = uint32(86400)
	GasAfterPaymentCalculation       = uint32(33825)

	VRFCoordinatorV2FeeConfig = vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: 500,
		FulfillmentFlatFeeLinkPPMTier2: 500,
		FulfillmentFlatFeeLinkPPMTier3: 500,
		FulfillmentFlatFeeLinkPPMTier4: 500,
		FulfillmentFlatFeeLinkPPMTier5: 500,
		ReqsForTier2:                   big.NewInt(0),
		ReqsForTier3:                   big.NewInt(0),
		ReqsForTier4:                   big.NewInt(0),
		ReqsForTier5:                   big.NewInt(0)}
)
