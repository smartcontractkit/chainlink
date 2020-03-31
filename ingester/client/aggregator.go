package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	AggregatorV1Filename = "aggregation.v1.json"
	AggregatorV2Filename = "aggregation.v2.json"

	LatestRoundFnName         = "latestRound"
	OraclesInstanceVar        = "oracles"
	ResponseReceivedEventName = "ResponseReceived"
	SubmissionReceivedName    = "SubmissionReceived"
	MaxOracleCount            = 45

	UnmarshalEmptyStringError = "abi: attempting to unmarshall an empty string while arguments are expected"
)

type OracleMapping map[common.Address]string

type AggregatorOracle struct {
	Name    string
	Address common.Address
}

func (om OracleMapping) AggregatorOracle(address common.Address) *AggregatorOracle {
	name, ok := om[address]
	if !ok {
		name = "Unknown"
	}
	return &AggregatorOracle{
		Name:    name,
		Address: address,
	}
}

type Aggregator interface {
	Name() string
	Address() common.Address
	LatestRound() (*big.Int, error)
	SubscribeToSubmissionReceived(chan<- types.Log) (Subscription, error)
	UnmarshalSubmissionReceivedEvent(types.Log) (*SubmissionReceivedEvent, error)
}

type aggregator struct {
	name    string
	client  ETH
	feedsUI FeedsUI
	abi     *abi.ABI
	address common.Address
}

func NewAggregator(client ETH, feedsUI FeedsUI, name string, address common.Address) (Aggregator, error) {
	aabi, err := client.ABI(AggregatorV2Filename)
	if err != nil {
		return nil, err
	}
	return &aggregator{
		name:    name,
		client:  client,
		feedsUI: feedsUI,
		abi:     &aabi,
		address: address,
	}, nil
}

func (a *aggregator) Name() string {
	return a.name
}

func (a *aggregator) Address() common.Address {
	return a.address
}

func (a *aggregator) LatestRound() (*big.Int, error) {
	var round *big.Int
	return round, a.client.Call(a.address, a.abi, LatestRoundFnName, &round)
}

func (a *aggregator) SubscribeToSubmissionReceived(logChan chan<- types.Log) (Subscription, error) {
	e := a.abi.Events[SubmissionReceivedName]
	q := ethereum.FilterQuery{
		Addresses: []common.Address{a.address},
		Topics:    [][]common.Hash{{e.ID()}},
	}
	sub, err := a.client.SubscribeToLogs(logChan, q)
	if err != nil {
		return nil, err
	}
	return sub, err
}

// UnmarshalSubmissionReceivedEvent hydrates a struct from the log topics emitted from the solidity event
func (a *aggregator) UnmarshalSubmissionReceivedEvent(log types.Log) (*SubmissionReceivedEvent, error) {
	sr := &SubmissionReceivedEvent{}
	if len(log.Topics) == 3 {
		sr.Answer = log.Topics[0].String()
		sr.RoundID = log.Topics[1].Big()
		sr.Oracle = common.BytesToAddress(log.Topics[2].Bytes())
	} else {
		return sr, errors.New("invalid log type while un-marshaling submission received, expected 3 topics")
	}
	return sr, nil
}
