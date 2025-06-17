package contracts

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
)

type Coordinator interface {
	ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error)
	ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error)
	Address() string
	WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error)
	WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error)
	FilterRandomWordsFulfilledEvent(opts *bind.FilterOpts, requestId *big.Int) (*CoordinatorRandomWordsFulfilled, error)
}

type Subscription struct {
	Balance       *big.Int
	NativeBalance *big.Int
	ReqCount      uint64
	SubOwner      common.Address
	Consumers     []common.Address
}

type CoordinatorConfigSet struct {
	MinimumRequestConfirmations       uint16
	MaxGasLimit                       uint32
	StalenessSeconds                  uint32
	GasAfterPaymentCalculation        uint32
	FallbackWeiPerUnitLink            *big.Int
	FulfillmentFlatFeeNativePPM       uint32
	FulfillmentFlatFeeLinkDiscountPPM uint32
	NativePremiumPercentage           uint8
	LinkPremiumPercentage             uint8
	FeeConfig                         VRFCoordinatorV2FeeConfig
	Raw                               types.Log
}

type VRFCoordinatorV2FeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}

type RandomWordsFulfilledEventFilter struct {
	RequestIds []*big.Int
	SubIDs     []*big.Int
	Timeout    time.Duration
}

type CoordinatorRandomWordsFulfilled struct {
	RequestId     *big.Int
	OutputSeed    *big.Int
	SubId         string
	Payment       *big.Int
	NativePayment bool
	Success       bool
	OnlyPremium   bool
	Raw           types.Log
}

type CoordinatorRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       string
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	ExtraArgs                   []byte
	Sender                      common.Address
	Raw                         types.Log
}

func parseRequestRandomnessLogs(coordinator Coordinator, logs []*types.Log) (*CoordinatorRandomWordsRequested, error) {
	var randomWordsRequestedEvent *CoordinatorRandomWordsRequested
	var err error
	for _, eventLog := range logs {
		for _, topic := range eventLog.Topics {
			if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested{}.Topic()) == 0 ||
				topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic()) == 0 {
				randomWordsRequestedEvent, err = coordinator.ParseRandomWordsRequested(*eventLog)
				if err != nil {
					return nil, fmt.Errorf("parse RandomWordsRequested log failed, err: %w", err)
				}
			}
		}
	}
	return randomWordsRequestedEvent, nil
}

func ParseRandomWordsFulfilledLogs(coordinator Coordinator, logs []*types.Log) ([]*CoordinatorRandomWordsFulfilled, error) {
	var randomWordsFulfilledEventArr []*CoordinatorRandomWordsFulfilled
	for _, eventLog := range logs {
		for _, topic := range eventLog.Topics {
			if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled{}.Topic()) == 0 ||
				topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled{}.Topic()) == 0 {
				randomWordsFulfilledEvent, err := coordinator.ParseRandomWordsFulfilled(*eventLog)
				if err != nil {
					return nil, fmt.Errorf("parse RandomWordsFulfilled log failed, err: %w", err)
				}
				randomWordsFulfilledEventArr = append(randomWordsFulfilledEventArr, randomWordsFulfilledEvent)
			}
		}
	}
	return randomWordsFulfilledEventArr, nil
}
