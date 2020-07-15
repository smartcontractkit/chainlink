package contracts

import (
	"encoding/hex"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

//go:generate mockery --name FluxAggregator --output ../../../internal/mocks/ --case=underscore

type Oracle interface {
	eth.ConnectedContract
}

const (
	// OracleName is the name of Chainlink's Ethereum contract that can respond to Oracle requests.
	OracleName = "Oracle"
)

var (
	// OracleRequestTopic20190207 was the new OracleRequest filter topic as of 2019-012-07.
	// It is the earlier OracleRequest topic to be used on mainnet.
	OracleRequestTopic20190207 = eth.MustGetV6ContractEventID("Oracle", "OracleRequest")
)

type oracle struct {
	eth.ConnectedContract
	ethClient eth.Client
	address   common.Address
}

type LogOracleRequest struct {
	models.Log
	SpecID             [32]byte       `abi:"specId"`
	Requester          common.Address `abi:"requester"`
	RequestID          [32]byte       `abi:"requestId"`
	Payment            *big.Int       `abi:"payment"`
	CallbackAddr       common.Address `abi:"callbackAddr"`
	CallbackFunctionID [4]byte        `abi:"callbackFunctionId"`
	CancelExpiration   *big.Int       `abi:"cancelExpiration"`
	DataVersion        *big.Int       `abi:"dataVersion"`
	Data               []byte         `abi:"data"`
}

func (l LogOracleRequest) OracleRequest() (models.OracleRequest, error) {
	var payment big.Int
	if l.Payment == nil {
		payment = *big.NewInt(0)
	} else {
		payment = *l.Payment
	}

	jobSpecID, err := l.normalizeJobSpecID()
	if err != nil {
		return models.OracleRequest{}, errors.Wrap(err, "could not decode job spec ID")
	}

	return models.OracleRequest{
		SpecID:             jobSpecID,
		Requester:          l.Requester,
		RequestID:          common.BytesToHash(l.RequestID[:]),
		Payment:            assets.Link(payment),
		CallbackAddr:       l.CallbackAddr,
		CallbackFunctionID: models.BytesToFunctionSelector(l.CallbackFunctionID[:]),
		CancelExpiration:   time.Unix(l.CancelExpiration.Int64(), 0),
		DataVersion:        *utils.NewBig(l.DataVersion),
		Data:               l.Data,
	}, nil
}

func (l LogOracleRequest) normalizeJobSpecID() (uuid.UUID, error) {
	hexEncoded := false
	for i := 16; i < 32; i++ {
		if l.SpecID[i] != 0 {
			hexEncoded = true
			break
		}
	}

	if hexEncoded {
		return decodeHexJobSpecID(l.SpecID)
	}
	return decodePaddedJobSpecID(l.SpecID)
}

func decodePaddedJobSpecID(raw [32]byte) (uuid.UUID, error) {
	return uuid.FromBytes(raw[16:])
}

func decodeHexJobSpecID(raw [32]byte) (uuid.UUID, error) {
	var b []byte
	_, err := hex.Decode(b, raw[:])
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.FromBytes(b)
}

var oracleRequestLogTypes = map[common.Hash]interface{}{
	OracleRequestTopic20190207: LogOracleRequest{},
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
