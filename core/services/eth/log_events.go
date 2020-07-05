package eth

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
)

// From evm-contracts/v0.6/Oracle.sol
// event OracleRequest(
//   bytes32 indexed specId,
//   address requester,
//   bytes32 requestId,
//   uint256 payment,
//   address callbackAddr,
//   bytes4 callbackFunctionId,
//   uint256 cancelExpiration,
//   uint256 dataVersion,
//   bytes data
// );

// OracleRequestEvent represents an emitted Oracle request log
type OracleRequestEvent struct {
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

// DecodeOracleRequestLogEvent attempts to decode the log into an OracleRequestEvent
// Returns error if the log data is not the correct shape
func DecodeOracleRequestLogEvent(rawLog models.Log) (models.OracleRequest, error) {
	reqEvent := OracleRequestEvent{}

	oracleCodec, err := GetV6ContractCodec("Oracle")
	if err != nil {
		return models.OracleRequest{}, errors.Wrap(err, "DecodeOracleRequestLogEvent failed in GetV6ContractCodec")
	}
	err = oracleCodec.UnpackLog(&reqEvent, "OracleRequest", rawLog)
	if err != nil {
		return models.OracleRequest{}, errors.Wrap(err, "DecodeOracleRequestLogEvent failed in UnpackLog")
	}
	return reqEvent.toOracleRequest(), nil
}

func (re OracleRequestEvent) toOracleRequest() models.OracleRequest {
	req := models.OracleRequest{}

	req.SpecID = common.BytesToHash(re.SpecId[:])
	req.Requester = re.Requester
	req.RequestID = common.BytesToHash(re.RequestId[:])
	req.Payment = *assets.NewLink(re.Payment.Int64())
	req.CallbackAddr = re.CallbackAddr
	req.CallbackFunctionID = models.BytesToFunctionSelector(re.CallbackFunctionId[:])
	req.CancelExpiration = time.Unix(re.CancelExpiration.Int64(), 0)
	req.DataVersion = re.DataVersion
	req.Data = re.Data

	return req
}
