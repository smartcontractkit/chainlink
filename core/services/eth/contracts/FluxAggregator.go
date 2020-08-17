package contracts

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

//go:generate mockery --name FluxAggregator --output ../../../internal/mocks/ --case=underscore

type FluxAggregator interface {
	eth.ConnectedContract
	RoundState(oracle common.Address, roundID uint32) (FluxAggregatorRoundState, error)
}

const (
	// FluxAggregatorName is the name of Chainlink's Ethereum contract for
	// aggregating numerical data such as prices.
	FluxAggregatorName = "FluxAggregator"
)

var (
	// AggregatorNewRoundLogTopic20191220 is the NewRound filter topic for
	// the FluxAggregator as of Dec. 20th 2019. Eagerly fails if not found.
	AggregatorNewRoundLogTopic20191220 = eth.MustGetV6ContractEventID("FluxAggregator", "NewRound")
	// AggregatorAnswerUpdatedLogTopic20191220 is the AnswerUpdated filter topic for
	// the FluxAggregator as of Dec. 20th 2019. Eagerly fails if not found.
	AggregatorAnswerUpdatedLogTopic20191220 = eth.MustGetV6ContractEventID("FluxAggregator", "AnswerUpdated")
)

type fluxAggregator struct {
	eth.ConnectedContract
	ethClient eth.Client
	address   common.Address
}

type LogNewRound struct {
	eth.GethRawLog
	RoundId   *big.Int
	StartedBy common.Address
	// seconds since unix epoch
	StartedAt *big.Int
}

type LogAnswerUpdated struct {
	eth.GethRawLog
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
}

var fluxAggregatorLogTypes = map[common.Hash]eth.Log{
	AggregatorNewRoundLogTopic20191220:      &LogNewRound{},
	AggregatorAnswerUpdatedLogTopic20191220: &LogAnswerUpdated{},
}

func NewFluxAggregator(address common.Address, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) (FluxAggregator, error) {
	codec, err := eth.GetV6ContractCodec(FluxAggregatorName)
	if err != nil {
		return nil, err
	}
	connectedContract := eth.NewConnectedContract(codec, address, ethClient, logBroadcaster)
	return &fluxAggregator{connectedContract, ethClient, address}, nil
}

func (fa *fluxAggregator) SubscribeToLogs(listener eth.LogListener) (connected bool, _ eth.UnsubscribeFunc) {
	return fa.ConnectedContract.SubscribeToLogs(
		eth.NewDecodingLogListener(fa, fluxAggregatorLogTypes, listener),
	)
}

type FluxAggregatorRoundState struct {
	ReportableRoundID uint32   `abi:"_roundId"`
	EligibleToSubmit  bool     `abi:"_eligibleToSubmit"`
	LatestAnswer      *big.Int `abi:"_latestSubmission"`
	Timeout           uint64   `abi:"_timeout"`
	StartedAt         uint64   `abi:"_startedAt"`
	AvailableFunds    *big.Int `abi:"_availableFunds"`
	PaymentAmount     *big.Int `abi:"_paymentAmount"`
	OracleCount       uint8    `abi:"_oracleCount"`
}

func (rs FluxAggregatorRoundState) TimesOutAt() uint64 {
	return rs.StartedAt + rs.Timeout
}

func (fa *fluxAggregator) RoundState(oracle common.Address, roundID uint32) (FluxAggregatorRoundState, error) {
	var result FluxAggregatorRoundState
	err := fa.Call(&result, "oracleRoundState", oracle, roundID)
	if err != nil {
		return FluxAggregatorRoundState{}, errors.Wrap(err, "unable to encode message call")
	}
	return result, nil
}
