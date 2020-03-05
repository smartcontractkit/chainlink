package contracts

import (
	"math/big"

	"chainlink/core/eth"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

//go:generate mockery -name FluxAggregator -output ../../internal/mocks/ -case=underscore

type FluxAggregator interface {
	eth.ConnectedContract
	SubscribeToLogs(chMaybeLogs chan<- eth.MaybeDecodedLog) eth.UnsubscribeFunc
	LatestAnswer(precision int32) (decimal.Decimal, error)
	LatestRound() (*big.Int, error)
	ReportingRound() (*big.Int, error)
	TimedOutStatus(round *big.Int) (bool, error)
	LatestSubmission(oracleAddress common.Address) (*big.Int, *big.Int, error)
}

type fluxAggregator struct {
	eth.ConnectedContract
	ethClient eth.Client
	address   common.Address
}

type LogNewRound struct {
	Raw *eth.Log

	RoundID   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
}

type LogRoundDetailsUpdated struct {
	Raw *eth.Log

	PaymentAmount  *big.Int
	MinAnswerCount *big.Int
	MaxAnswerCount *big.Int
	RestartDelay   *big.Int
	Timeout        *big.Int
	Address        common.Address
}

type LogAnswerUpdated struct {
	Raw *eth.Log

	Current   *big.Int
	RoundID   *big.Int
	Timestamp *big.Int
	Address   common.Address
}

var fluxAggregatorLogTypes = map[common.Hash]interface{}{
	models.AggregatorNewRoundLogTopic20191220:            LogNewRound{},
	models.AggregatorRoundDetailsUpdatedLogTopic20191220: LogRoundDetailsUpdated{},
	models.AggregatorAnswerUpdatedLogTopic20191220:       LogAnswerUpdated{},
}

func NewFluxAggregator(address common.Address, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) (FluxAggregator, error) {
	contract, err := eth.GetV6Contract(eth.FluxAggregatorName)
	if err != nil {
		return nil, err
	}
	connectedContract := eth.NewConnectedContract(contract, address, ethClient, logBroadcaster)
	return &fluxAggregator{connectedContract, ethClient, address}, nil
}

func (fa *fluxAggregator) SubscribeToLogs(chMaybeLogs chan<- eth.MaybeDecodedLog) eth.UnsubscribeFunc {
	listener := NewDecodingLogListener(fa, fluxAggregatorLogTypes, func(decodedLog interface{}, err error) {
		chMaybeLogs <- eth.MaybeDecodedLog{decodedLog, err}
	})
	return fa.ConnectedContract.SubscribeToLogs(listener)
}

// Price returns the current price at the given aggregator address.
func (fa *fluxAggregator) LatestAnswer(precision int32) (decimal.Decimal, error) {
	var result *big.Int
	err := fa.Call(&result, "latestAnswer")
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, "unable to encode message call")
	}
	raw := decimal.NewFromBigInt(result, 0)
	precisionDivisor := dec10.Pow(decimal.NewFromInt32(precision))
	return raw.Div(precisionDivisor), nil

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
