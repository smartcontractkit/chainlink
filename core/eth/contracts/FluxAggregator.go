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

type FluxAggregator struct {
	*eth.ConnectedContract
	ethClient eth.Client
	address   common.Address
}

type LogNewRound struct {
	RoundID   *big.Int       `abi:"roundId"`
	StartedBy common.Address `abi:"startedBy"`
	StartedAt *big.Int       `abi:"startedAt"`
	Address   common.Address `abi:"-"`
}

type LogRoundDetailsUpdated struct {
	PaymentAmount  *big.Int       `abi:"paymentAmount"`
	MinAnswerCount *big.Int       `abi:"minAnswerCount"`
	MaxAnswerCount *big.Int       `abi:"maxAnswerCount"`
	RestartDelay   *big.Int       `abi:"restartDelay"`
	Timeout        *big.Int       `abi:"timeout"`
	Address        common.Address `abi:"-"`
}

type LogAnswerUpdated struct {
	Current   *big.Int       `abi:"current"`
	RoundID   *big.Int       `abi:"roundId"`
	Timestamp *big.Int       `abi:"timestamp"`
	Address   common.Address `abi:"-"`
}

func NewFluxAggregator(ethClient eth.Client, address common.Address) (*FluxAggregator, error) {
	contract, err := eth.GetV6Contract(eth.FluxAggregatorName)
	if err != nil {
		return nil, err
	}
	connectedContract := contract.Connect(ethClient, address)
	return &FluxAggregator{connectedContract, ethClient, address}, nil
}

func (fa *FluxAggregator) SubscribeToLogs(fromBlock *big.Int) (*LogSubscription, error) {
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
					err = fa.Contract.UnpackLog(&decodedLog, "NewRound", rawLog)
					decodedLog.Address = fa.address
					chLogs <- MaybeDecodedLog{decodedLog, err}

				case models.AggregatorRoundDetailsUpdatedLogTopic20191220:
					var decodedLog LogRoundDetailsUpdated
					err = fa.Contract.UnpackLog(&decodedLog, "RoundDetailsUpdated", rawLog)
					decodedLog.Address = fa.address
					chLogs <- MaybeDecodedLog{decodedLog, err}

				case models.AggregatorAnswerUpdatedLogTopic20191220:
					var decodedLog LogAnswerUpdated
					err = fa.Contract.UnpackLog(&decodedLog, "AnswerUpdated", rawLog)
					decodedLog.Address = fa.address
					chLogs <- MaybeDecodedLog{decodedLog, err}

				default:
					// @@TODO: warn?
				}
			}
		}
	}()

	return &LogSubscription{subscription, chLogs}, nil
}

// Price returns the current price at the given aggregator address.
func (fa *FluxAggregator) Price(precision int32) (decimal.Decimal, error) {
	var result string
	err := fa.Call(&result, "latestAnswer")
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, "unable to encode message call")
	}

	raw, err := newDecimalFromString(result)
	if err != nil {
		return decimal.Decimal{}, errors.Wrapf(err, "unable to fetch aggregator price from %s", fa.address.Hex())
	}
	precisionDivisor := dec10.Pow(decimal.NewFromInt32(precision))
	return raw.Div(precisionDivisor), nil

}

// LatestRound returns the latest round at the given aggregator address.
func (fa *FluxAggregator) LatestRound() (*big.Int, error) {
	var result *big.Int
	err := fa.Call(&result, "latestRound")
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// ReportingRound returns the reporting round at the given aggregator address.
func (fa *FluxAggregator) ReportingRound() (*big.Int, error) {
	var result *big.Int
	err := fa.Call(&result, "reportingRound")
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// TimedOutStatus returns the a boolean indicating whether the provided round has timed out or not
func (fa *FluxAggregator) TimedOutStatus(round *big.Int) (bool, error) {
	var result bool
	err := fa.Call(&result, "getTimedOutStatus", round)
	if err != nil {
		return false, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}

// LatestSubmission returns the latest submission as a tuple, (answer, round)
// for a given oracle address.
func (fa *FluxAggregator) LatestSubmission(oracleAddress common.Address) (*big.Int, *big.Int, error) {
	var result struct {
		LatestAnswer        *big.Int
		LatestReportedRound *big.Int
	}

	err := fa.Call(&result, "latestSubmission", oracleAddress)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to encode message call")
	}
	return result.LatestAnswer, result.LatestReportedRound, nil
}
