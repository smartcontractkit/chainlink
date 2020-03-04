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
	SubscribeToLogs(fromBlock *big.Int) (LogSubscription, error)
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
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Address   common.Address
}

type LogRoundDetailsUpdated struct {
	PaymentAmount  *big.Int
	MinAnswerCount *big.Int
	MaxAnswerCount *big.Int
	RestartDelay   *big.Int
	Timeout        *big.Int
	Address        common.Address
}

type LogAnswerUpdated struct {
	Current   *big.Int
	RoundID   *big.Int
	Timestamp *big.Int
	Address   common.Address
}

func NewFluxAggregator(ethClient eth.Client, address common.Address) (FluxAggregator, error) {
	contract, err := eth.GetV6Contract(eth.FluxAggregatorName)
	if err != nil {
		return nil, err
	}
	connectedContract := eth.NewConnectedContract(contract, ethClient, address)
	return &fluxAggregator{connectedContract, ethClient, address}, nil
}

func (fa *fluxAggregator) SubscribeToLogs(fromBlock *big.Int) (LogSubscription, error) {
	var (
		filterQuery = ethereum.FilterQuery{
			FromBlock: fromBlock,
			Addresses: utils.WithoutZeroAddresses([]common.Address{fa.address}),
			Topics:    [][]common.Hash{{models.AggregatorNewRoundLogTopic20191220, models.AggregatorRoundDetailsUpdatedLogTopic20191220, models.AggregatorAnswerUpdatedLogTopic20191220}},
		}
		chLogs    = make(chan MaybeDecodedLog)
		chRawLogs = make(chan eth.Log)
	)

	subscription, err := fa.ethClient.SubscribeToLogs(chRawLogs, filterQuery)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(chLogs)
		for {
			select {
			case err := <-subscription.Err():
				chLogs <- MaybeDecodedLog{nil, err}

			case rawLog, stillOpen := <-chRawLogs:
				// @@TODO: make sure that calling .Unsubscribe() closes chRawLogs
				if !stillOpen {
					select {
					case err := <-subscription.Err():
						chLogs <- MaybeDecodedLog{nil, err}
					default:
					}
					return
				}

				switch rawLog.Topics[0] {
				case models.AggregatorNewRoundLogTopic20191220:
					var decodedLog LogNewRound
					err = fa.UnpackLog(&decodedLog, "NewRound", rawLog)
					decodedLog.Address = fa.address
					chLogs <- MaybeDecodedLog{decodedLog, err}

				case models.AggregatorRoundDetailsUpdatedLogTopic20191220:
					var decodedLog LogRoundDetailsUpdated
					err = fa.UnpackLog(&decodedLog, "RoundDetailsUpdated", rawLog)
					decodedLog.Address = fa.address
					chLogs <- MaybeDecodedLog{decodedLog, err}

				case models.AggregatorAnswerUpdatedLogTopic20191220:
					var decodedLog LogAnswerUpdated
					err = fa.UnpackLog(&decodedLog, "AnswerUpdated", rawLog)
					decodedLog.Address = fa.address
					chLogs <- MaybeDecodedLog{decodedLog, err}

				default:
					// @@TODO: warn?
				}
			}
		}
	}()

	return &logSubscription{subscription, chLogs}, nil
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
