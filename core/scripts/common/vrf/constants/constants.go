package constants

import (
	"math/big"
)

var (
	SubscriptionBalanceJuels     = "1e19"
	SubscriptionBalanceNativeWei = "1e18"

	// optional flags
	FallbackWeiPerUnitLink      = big.NewInt(6e16)
	BatchFulfillmentEnabled     = true
	MinConfs                    = 3
	NodeSendingKeyFundingAmount = "1e17"
	MaxGasLimit                 = int64(2.5e6)
	StalenessSeconds            = int64(86400)
	GasAfterPayment             = int64(33285)

	//vrfv2
	FlatFeeTier1 = int64(500)
	FlatFeeTier2 = int64(500)
	FlatFeeTier3 = int64(500)
	FlatFeeTier4 = int64(500)
	FlatFeeTier5 = int64(500)
	ReqsForTier2 = int64(0)
	ReqsForTier3 = int64(0)
	ReqsForTier4 = int64(0)
	ReqsForTier5 = int64(0)

	//vrfv2plus
	FlatFeeNativePPM        = uint32(500)
	FlatFeeLinkDiscountPPM  = uint32(100)
	NativePremiumPercentage = uint8(1)
	LinkPremiumPercentage   = uint8(1)
)
