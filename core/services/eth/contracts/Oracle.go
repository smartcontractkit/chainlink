package contracts

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Oracle --output ../../../internal/mocks/ --case=underscore

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

	oracleLogTypes = map[common.Hash]interface{}{
		OracleRequestLogTopic20190207:       &LogOracleRequest{},
		CancelOracleRequestLogTopic20190207: &LogCancelOracleRequest{},
	}
)

type (
	Oracle interface {
		eth.ConnectedContract
	}

	oracle struct {
		eth.ConnectedContract
		ethClient eth.Client
		address   common.Address
	}

	OracleRequest struct {
		SpecID             common.Hash
		Requester          common.Address
		RequestID          common.Hash
		Payment            *assets.Link
		CallbackAddr       common.Address
		CallbackFunctionID models.FunctionSelector
		CancelExpiration   time.Time
		DataVersion        *big.Int
		Data               []byte
	}

	LogOracleRequest struct {
		types.Log
		SpecId             [32]byte
		Requester          common.Address
		RequestId          [32]byte
		Payment            *big.Int
		CallbackAddr       common.Address
		CallbackFunctionId [4]byte
		CancelExpiration   *big.Int
		DataVersion        *big.Int
		Data               []byte
	}

	LogCancelOracleRequest struct {
		types.Log
		RequestId [32]byte
	}
)

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
		eth.NewDecodingLogListener(o, oracleLogTypes, listener),
	)
}

func (re LogOracleRequest) ToOracleRequest() OracleRequest {
	req := OracleRequest{}

	req.SpecID = common.BytesToHash(re.SpecId[:])
	req.Requester = re.Requester
	req.RequestID = common.BytesToHash(re.RequestId[:])
	req.Payment = (*assets.Link)(re.Payment)
	req.CallbackAddr = re.CallbackAddr
	req.CallbackFunctionID = models.BytesToFunctionSelector(re.CallbackFunctionId[:])
	req.CancelExpiration = time.Unix(re.CancelExpiration.Int64(), 0)
	req.DataVersion = re.DataVersion
	req.Data = re.Data

	return req
}

func (o OracleRequest) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["specId"] = o.SpecID
	m["requester"] = o.Requester
	m["requestId"] = o.RequestID
	m["payment"] = o.Payment
	m["callbackAddr"] = o.CallbackAddr
	m["callbackFunctionId"] = o.CallbackFunctionID
	m["cancelExpiration"] = o.CancelExpiration
	m["dataVersion"] = o.DataVersion
	m["data"] = o.Data
	return m
}
