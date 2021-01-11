package contracts

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockery --name FluxAggregator --output ../../../internal/mocks/ --case=underscore

var FluxAggregatorABI = mustGetABI(flux_aggregator_wrapper.FluxAggregatorABI)

type FluxAggregator interface {
	Address() common.Address
	GetOracles(opts *bind.CallOpts) ([]common.Address, error)
	OracleRoundState(opts *bind.CallOpts, oracle common.Address, roundID uint32) (FluxAggregatorRoundState, error)
	LatestRoundData(opts *bind.CallOpts) (FluxAggregatorRoundData, error)
}

type NewRound = flux_aggregator_wrapper.FluxAggregatorNewRound
type AnswerUpdated = flux_aggregator_wrapper.FluxAggregatorAnswerUpdated

// const (
// 	// FluxAggregatorName is the name of Chainlink's Ethereum contract for
// 	// aggregating numerical data such as prices.
// 	FluxAggregatorName = "FluxAggregator"
// )

// var (
// 	// AggregatorNewRoundLogTopic20191220 is the NewRound filter topic for
// 	// the FluxAggregator as of Dec. 20th 2019. Eagerly fails if not found.
// 	AggregatorNewRoundLogTopic20191220 = eth.MustGetV6ContractEventID("FluxAggregator", "NewRound")
// 	// AggregatorAnswerUpdatedLogTopic20191220 is the AnswerUpdated filter topic for
// 	// the FluxAggregator as of Dec. 20th 2019. Eagerly fails if not found.
// 	AggregatorAnswerUpdatedLogTopic20191220 = eth.MustGetV6ContractEventID("FluxAggregator", "AnswerUpdated")
// )

// type fluxAggregator struct {
// 	ConnectedContract
// 	ethClient eth.Client
// 	address   common.Address
// }

// 	types.Log
// type LogNewRound struct {
// 	RoundId   *big.Int
// 	StartedBy common.Address
// 	// seconds since unix epoch
// 	StartedAt *big.Int
// }

// type LogAnswerUpdated struct {
// 	types.Log
// 	Current   *big.Int
// 	RoundId   *big.Int
// 	UpdatedAt *big.Int
// }

// var fluxAggregatorLogTypes = map[common.Hash]interface{}{
// 	AggregatorNewRoundLogTopic20191220:      &LogNewRound{},
// 	AggregatorAnswerUpdatedLogTopic20191220: &LogAnswerUpdated{},
// }

// func NewFluxAggregator(address common.Address, ethClient eth.Client, logBroadcaster log.Broadcaster) (FluxAggregator, error) {
// 	codec, err := eth.GetV6ContractCodec(FluxAggregatorName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	connectedContract := NewConnectedContract(codec, address, ethClient, logBroadcaster)
// 	return &fluxAggregator{connectedContract, ethClient, address}, nil
// }

// func (fa *fluxAggregator) SubscribeToLogs(listener log.Listener) (connected bool, _ UnsubscribeFunc) {
// 	return fa.ConnectedContract.SubscribeToLogs(
// 		log.NewDecodingLogListener(fa, fluxAggregatorLogTypes, listener),
// 	)
// }

type FluxAggregatorRoundState struct {
	EligibleToSubmit bool     `abi:"_eligibleToSubmit" json:"eligibleToSubmit"`
	RoundId          uint32   `abi:"_roundId" json:"reportableRoundID"`
	LatestSubmission *big.Int `abi:"_latestSubmission" json:"latestAnswer,omitempty"`
	StartedAt        uint64   `abi:"_startedAt" json:"startedAt"`
	Timeout          uint64   `abi:"_timeout" json:"timeout"`
	AvailableFunds   *big.Int `abi:"_availableFunds" json:"availableFunds,omitempty"`
	OracleCount      uint8    `abi:"_oracleCount" json:"oracleCount"`
	PaymentAmount    *big.Int `abi:"_paymentAmount" json:"paymentAmount,omitempty"`
}

type FluxAggregatorRoundData struct {
	RoundId         *big.Int `abi:"roundId" json:"reportableRoundID"`
	Answer          *big.Int `abi:"answer" json:"latestAnswer,omitempty"`
	StartedAt       *big.Int `abi:"startedAt" json:"startedAt"`
	UpdatedAt       *big.Int `abi:"updatedAt" json:"updatedAt"`
	AnsweredInRound *big.Int `abi:"answeredInRound" json:"availableFunds,omitempty"`
}

func (rs FluxAggregatorRoundState) TimesOutAt() uint64 {
	return rs.StartedAt + rs.Timeout
}

// func (fa *fluxAggregator) RoundState(oracle common.Address, roundID uint32) (FluxAggregatorRoundState, error) {
// 	var result FluxAggregatorRoundState
// 	err := fa.Call(&result, "oracleRoundState", oracle, roundID)
// 	if err != nil {
// 		return FluxAggregatorRoundState{}, errors.Wrap(err, "unable to encode message call")
// 	}
// 	return result, nil
// }

// func (fa *fluxAggregator) GetOracles() (oracles []common.Address, err error) {
// 	oracles = make([]common.Address, 0)
// 	err = fa.Call(&oracles, "getOracles")
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error calling flux aggregator getOracles")
// 	}
// 	return oracles, nil
// }

// func (fa *fluxAggregator) LatestRoundData() (FluxAggregatorRoundData, error) {
// 	var result FluxAggregatorRoundData
// 	err := fa.Call(&result, "latestRoundData")
// 	if err != nil {
// 		return FluxAggregatorRoundData{},
// 			errors.Wrap(err, "error calling fluxaggregator#latestRoundData - contract may have 0 rounds")
// 	}
// 	return result, nil
// }

// ////////////////////////////////////////////////////////////////////////////////////////////////////

func NewFluxAggregatorContract(address common.Address, backend bind.ContractBackend) (FluxAggregator, error) {
	wrapper, err := flux_aggregator_wrapper.NewFluxAggregator(address, backend)
	if err != nil {
		return nil, err
	}
	return fluxAggregator{
		address: address,
		// wrapper: *wrapper,
		FluxAggregator: wrapper,
	}, nil
}

type fluxAggregator struct {
	address common.Address
	*flux_aggregator_wrapper.FluxAggregator
	// wrapper flux_aggregator_wrapper.FluxAggregator
}

func (fa fluxAggregator) Address() common.Address {
	return fa.address
}

// func (fa fluxAggregator) Wrapper() flux_aggregator_wrapper.FluxAggregator {
// 	return fa.wrapper
// }

// func (fa fluxAggregator) GetOracles(opts *bind.CallOpts) ([]common.Address, error) {
// 	return fa.wrapper.GetOracles(opts)
// }

func (fa fluxAggregator) OracleRoundState(opts *bind.CallOpts, oracle common.Address, roundID uint32) (FluxAggregatorRoundState, error) {
	res, err := fa.FluxAggregator.OracleRoundState(opts, oracle, roundID)
	return FluxAggregatorRoundState(res), err
}

func (fa fluxAggregator) LatestRoundData(opts *bind.CallOpts) (FluxAggregatorRoundData, error) {
	res, err := fa.FluxAggregator.LatestRoundData(opts)
	return FluxAggregatorRoundData(res), err
}

type FluxAggregatorDecodingLogListener struct {
	wrapper flux_aggregator_wrapper.FluxAggregator
	log.Listener
}

var _ log.Listener = (*FluxAggregatorDecodingLogListener)(nil)

func NewFluxAggregatorDecodingLogListener(
	address common.Address,
	backend bind.ContractBackend,
	innerListener log.Listener,
) (log.Listener, error) {
	wrapper, err := flux_aggregator_wrapper.NewFluxAggregator(address, backend)
	if err != nil {
		return nil, err
	}
	return FluxAggregatorDecodingLogListener{
		wrapper:  *wrapper,
		Listener: innerListener,
	}, nil
}

func (ll FluxAggregatorDecodingLogListener) HandleLog(lb log.Broadcast, err error) {
	if err != nil {
		ll.Listener.HandleLog(lb, err)
		return
	}

	rawLog := lb.RawLog()
	if len(rawLog.Topics) == 0 {
		return
	}
	eventID := rawLog.Topics[0]
	var decodedLog interface{}

	switch eventID {
	case FluxAggregatorABI.Events["NewRound"].ID:
		decodedLog, err = ll.wrapper.ParseNewRound(rawLog)
	case FluxAggregatorABI.Events["AnswerUpdated"].ID:
		decodedLog, err = ll.wrapper.ParseAnswerUpdated(rawLog)
	default:
		logger.Warnf("Unknown topic for FluxAggregator contract: %s", eventID.Hex())
		return // don't pass on unknown/unexpected events
	}

	lb.SetDecodedLog(decodedLog)
	ll.Listener.HandleLog(lb, err)
}
