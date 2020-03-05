package contracts

import (
	"math/big"

	"chainlink/core/eth"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

//go:generate mockery -name FluxAggregator -output ../../internal/mocks/ -case=underscore

type FluxAggregator interface {
	eth.ConnectedContract
	RoundState() (FluxAggregatorRoundState, error)
}

type fluxAggregator struct {
	eth.ConnectedContract
	ethClient eth.Client
	address   common.Address
}

type LogNewRound struct {
	eth.Log
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
}

type LogRoundDetailsUpdated struct {
	eth.Log
	PaymentAmount  *big.Int
	MinAnswerCount uint32
	MaxAnswerCount uint32
	RestartDelay   uint32
	Timeout        uint32
}

type LogAnswerUpdated struct {
	eth.Log
	Current   *big.Int
	RoundId   *big.Int
	Timestamp *big.Int
}

var fluxAggregatorLogTypes = map[common.Hash]interface{}{
	models.AggregatorNewRoundLogTopic20191220:            LogNewRound{},
	models.AggregatorRoundDetailsUpdatedLogTopic20191220: LogRoundDetailsUpdated{},
	models.AggregatorAnswerUpdatedLogTopic20191220:       LogAnswerUpdated{},
}

func NewFluxAggregator(address common.Address, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) (FluxAggregator, error) {
	codec, err := eth.GetV6ContractCodec(eth.FluxAggregatorName)
	if err != nil {
		return nil, err
	}
	connectedContract := eth.NewConnectedContract(codec, address, ethClient, logBroadcaster)
	return &fluxAggregator{connectedContract, ethClient, address}, nil
}

func (fa *fluxAggregator) SubscribeToLogs(listener eth.LogListener) eth.UnsubscribeFunc {
	return fa.ConnectedContract.SubscribeToLogs(
		eth.NewDecodingLogListener(fa, fluxAggregatorLogTypes, listener),
	)
}

type FluxAggregatorRoundState struct {
	ReportableRoundID *big.Int `abi:"_reportableRoundId"`
	EligibleToSubmit  bool     `abi:"_eligibleToSubmit"`
	LatestAnswer      *big.Int `abi:"_latestAnswer"`
}

func (fa *fluxAggregator) RoundState() (FluxAggregatorRoundState, error) {
	var result FluxAggregatorRoundState
	err := fa.Call(&result, "roundState")
	if err != nil {
		return FluxAggregatorRoundState{}, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// Price returns the current price at the given aggregator address.
func (fa *fluxAggregator) LatestAnswer(precision int32) (decimal.Decimal, error) {
	var result *big.Int
	err := fa.Call(&result, "latestAnswer")
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, "unable to encode message call")
	}
	return utils.DecimalFromBigInt(result, precision), nil
}

// LatestRound returns the latest round at the given aggregator address.
func (fa *fluxAggregator) LatestRound() (*big.Int, error) {
	var result *big.Int
	err := fa.Call(&result, "latestRound")
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// ReportingRound returns the reporting round at the given aggregator address.
func (fa *fluxAggregator) ReportingRound() (*big.Int, error) {
	var result *big.Int
	err := fa.Call(&result, "reportingRound")
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// TimedOutStatus returns the a boolean indicating whether the provided round has timed out or not
func (fa *fluxAggregator) TimedOutStatus(round *big.Int) (bool, error) {
	var result bool
	err := fa.Call(&result, "getTimedOutStatus", round)
	if err != nil {
		return false, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// LatestSubmission returns the latest submission as a tuple, (answer, round)
// for a given oracle address.
func (fa *fluxAggregator) LatestSubmission(oracleAddress common.Address) (*big.Int, *big.Int, error) {
	result := make([]interface{}, 2)
	err := fa.Call(&result, "latestSubmission", oracleAddress)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to encode message call")
	}
	latestAnswer := result[0].(*big.Int)
	latestReportedRound := result[1].(*big.Int)
	return latestAnswer, latestReportedRound, nil
}
