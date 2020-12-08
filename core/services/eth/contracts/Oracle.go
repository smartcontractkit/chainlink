package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

//go:generate mockery --name Oracle --output ../../../internal/mocks/ --case=underscore

type Oracle interface {
	eth.ConnectedContract
}

const (
	// OracleName is the name of Chainlink's Ethereum contract for
	// aggregating numerical data such as prices.
	OracleName = "Oracle"
)

var (
	// OracleRequestLogTopic20190207 is the Oracle Request filter topic for
	// the Oracle as of 2019-02-07. Eagerly fails if not found.
	OracleRequestLogTopic20190207 = eth.MustGetV6ContractEventID("Oracle", "OracleRequest")
	// CancelOracleRequestLogTopic20190207 is the Cancel Oracle Request filter topic for
	// the Oracle as of 2019-02-07. Eagerly fails if not found.
	CancelOracleRequestLogTopic20190207 = eth.MustGetV6ContractEventID("Oracle", "CancelOracleRequest")
)

type oracle struct {
	eth.ConnectedContract
	ethClient eth.Client
	address   common.Address
}

var fluxAggregatorLogTypes = map[common.Hash]interface{}{
	OracleRequestLogTopic20190207:       &LogOracleRequest{},
	CancelOracleRequestLogTopic20190207: &LogCancelOracleRequest{},
}

func NewOracle(address common.Address, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) (Oracle, error) {
	codec, err := eth.GetV6ContractCodec(OracleName)
	if err != nil {
		return nil, err
	}
	connectedContract := eth.NewConnectedContract(codec, address, ethClient, logBroadcaster)
	return &oracle{connectedContract, ethClient, address}, nil
}

func (o *oracle) SubscribeToLogs(listener eth.LogListener) (connected bool, _ eth.UnsubscribeFunc) {
	return o.ConnectedContract.SubscribeToLogs(
		eth.NewDecodingLogListener(o, fluxAggregatorLogTypes, listener),
	)
}
