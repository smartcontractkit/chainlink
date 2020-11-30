package adapters

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tidwall/gjson"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	// DataFormatBytes instructs the EthTx Adapter to treat the input value as a
	// bytes string, rather than a hexadecimal encoded bytes32
	DataFormatBytes = "bytes"
)

// EthTx holds the Address to send the result to and the FunctionSelector
// to execute.
type EthTx struct {
	ToAddress common.Address `json:"address"`
	// NOTE: FromAddress is deprecated and kept for backwards compatibility, new job specs should use fromAddresses
	FromAddress      common.Address          `json:"fromAddress,omitempty"`
	FromAddresses    []common.Address        `json:"fromAddresses,omitempty"`
	FunctionSelector models.FunctionSelector `json:"functionSelector"`
	// DataPrefix is typically a standard first argument
	// to chainlink callback calls - usually the requestID
	DataPrefix hexutil.Bytes `json:"dataPrefix"`
	DataFormat string        `json:"format"`
	GasLimit   uint64        `json:"gasLimit,omitempty"`

	// Optional list of desired encodings for ResultCollectKey arguments.
	// i.e. ["uint256", "bytes32"]
	ABIEncoding []string `json:"abiEncoding"`

	// MinRequiredOutgoingConfirmations only works with bulletprooftxmanager
	MinRequiredOutgoingConfirmations uint64 `json:"minRequiredOutgoingConfirmations,omitempty"`
}

// TaskType returns the type of Adapter.
func (e *EthTx) TaskType() models.TaskType {
	return TaskTypeEthTx
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (e *EthTx) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	trtx, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
	if err != nil {
		err = errors.Wrap(err, "FindEthTaskRunTxByTaskRunID failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}
	if trtx != nil {
		return e.checkForConfirmation(*trtx, input, store)
	}
	return e.insertEthTx(input, store)
}

func (e *EthTx) checkForConfirmation(trtx models.EthTaskRunTx,
	input models.RunInput, store *strpkg.Store) models.RunOutput {
	switch trtx.EthTx.State {
	case models.EthTxConfirmed:
		return e.checkEthTxForReceipt(trtx.EthTx.ID, input, store)
	case models.EthTxFatalError:
		return models.NewRunOutputError(trtx.EthTx.GetError())
	default:
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}
}

func (e *EthTx) pickFromAddress(input models.RunInput, store *strpkg.Store) (common.Address, error) {
	if len(e.FromAddresses) > 0 {
		if e.FromAddress != utils.ZeroAddress {
			logger.Warnf("task spec for task run %s specified both fromAddress and fromAddresses."+
				" fromAddress is deprecated, it will be ignored and fromAddresses used instead. "+
				"Specifying both of these keys in a job spec may result in an error in future versions of Chainlink", input.TaskRunID())
		}
		return store.GetRoundRobinAddress(e.FromAddresses...)
	}
	if e.FromAddress == utils.ZeroAddress {
		return store.GetRoundRobinAddress()
	}
	logger.Warnf(`DEPRECATION WARNING: task spec for task run %s specified a fromAddress of %s. fromAddress has been deprecated and will be removed in a future version of Chainlink. Please use fromAddresses instead. You can pin a job to one address simply by using only one element, like so:
{
	"type": "EthTx",
	"fromAddresses": ["%s"],
}
`, input.TaskRunID(), e.FromAddress.Hex(), e.FromAddress.Hex())
	return e.FromAddress, nil
}

func (e *EthTx) insertEthTx(input models.RunInput, store *strpkg.Store) models.RunOutput {
	var (
		txData, encodedPayload []byte
		err                    error
	)
	if e.ABIEncoding != nil {
		// The requestID is present as the first element of DataPrefix (from the oracle event log).
		// We prepend it as a magic first argument of for the consumer contract.
		data, errPrepend := input.Data().PrependAtArrayKey(models.ResultCollectionKey, e.DataPrefix[:32])
		if errPrepend != nil {
			return models.NewRunOutputError(err)
		}
		// Encode the calldata for the consumer contract.
		txData, err = getTxDataUsingABIEncoding(e.ABIEncoding, data.Get(models.ResultCollectionKey).Array())
	} else {
		txData, err = getTxData(e, input)
	}
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed while constructing EthTx data")
		return models.NewRunOutputError(err)
	}

	taskRunID := input.TaskRunID()
	toAddress := e.ToAddress
	fromAddress, err := e.pickFromAddress(input, store)
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed to pickFromAddress")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	if e.ABIEncoding != nil {
		// Encode the calldata for the operator/oracle contract. Note that the last argument is nested calldata, calldata
		// for the consumer contract.
		// [hash(fulfillOracleRequest2...)[:4]]	[..............................dataPrefix...............................] [call data]
		// [hash(fulfillOracleRequest2...)[:4]] [requestID] [payment] [callbackAddress] [callbackFunctionId] [expiration] [call data]
		// 6 = requestID + payment + callbackAddress + callbackFunctionId + expiration + offset itself
		payloadOffset := utils.EVMWordUint64(utils.EVMWordByteLen * 6)
		encodedPayload = append(append(append(e.FunctionSelector.Bytes(), e.DataPrefix...), payloadOffset...), utils.EVMEncodeBytes(txData)...)
	} else {
		encodedPayload = append(append(e.FunctionSelector.Bytes(), e.DataPrefix...), txData...)
	}

	var gasLimit uint64
	if e.GasLimit == 0 {
		gasLimit = store.Config.EthGasLimitDefault()
	} else {
		gasLimit = e.GasLimit
	}

	if err := store.IdempotentInsertEthTaskRunTx(taskRunID, fromAddress, toAddress, encodedPayload, gasLimit); err != nil {
		err = errors.Wrap(err, "insertEthTx failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
}

func (e *EthTx) checkEthTxForReceipt(ethTxID int64, input models.RunInput, s *strpkg.Store) models.RunOutput {
	var minRequiredOutgoingConfirmations uint64
	if e.MinRequiredOutgoingConfirmations == 0 {
		minRequiredOutgoingConfirmations = s.Config.MinRequiredOutgoingConfirmations()
	} else {
		minRequiredOutgoingConfirmations = e.MinRequiredOutgoingConfirmations
	}

	hash, err := getConfirmedTxHash(ethTxID, s.DB, minRequiredOutgoingConfirmations)

	if err != nil {
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	if hash == nil {
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}

	hexHash := (*hash).Hex()

	output := input.Data()
	output, err = output.MultiAdd(models.KV{
		"result": hexHash,
		// HACK: latestOutgoingTxHash is used for backwards compatibility with the stats pusher
		"latestOutgoingTxHash": hexHash,
	})
	if err != nil {
		err = errors.Wrap(err, "checkEthTxForReceipt failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}
	return models.NewRunOutputComplete(output)
}

func getConfirmedTxHash(ethTxID int64, db *gorm.DB, minRequiredOutgoingConfirmations uint64) (*common.Hash, error) {
	receipt := models.EthReceipt{}
	err := db.
		Joins("INNER JOIN eth_tx_attempts ON eth_tx_attempts.hash = eth_receipts.tx_hash AND eth_tx_attempts.eth_tx_id = ?", ethTxID).
		Joins("INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state = 'confirmed'").
		Where("eth_receipts.block_number <= (SELECT max(number) - ? FROM heads)", minRequiredOutgoingConfirmations).
		First(&receipt).
		Error

	if err == nil {
		return &receipt.TxHash, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return nil, errors.Wrap(err, "getConfirmedTxHash failed")

}

var (
	ErrInvalidABIEncoding = errors.New("invalid abi encoding")
	// A base set of supported types, expand as needed.
	// The corresponding go type is the type we need to pass into abi.Arguments.PackValues.
	solidityTypeToGoType = map[string]reflect.Type{
		"int256":  reflect.TypeOf(big.Int{}),
		"uint256": reflect.TypeOf(big.Int{}),
		"bool":    reflect.TypeOf(false),
		"bytes32": reflect.TypeOf([32]byte{}),
		"bytes4":  reflect.TypeOf([4]byte{}),
		"bytes":   reflect.TypeOf([]byte{}),
		"address": reflect.TypeOf(common.Address{}),
	}
	jsonTypes = map[gjson.Type]map[string]struct{}{
		gjson.String: {
			"bytes32": {},
			"bytes4":  {},
			"bytes":   {},
			"address": {},
		},
		gjson.True: {
			"bool": {},
		},
		gjson.False: {
			"bool": {},
		},
		gjson.Number: {
			"uint256": {},
			"int256":  {},
		},
	}
	supportedSolidityTypes []string
)

func init() {
	for k := range solidityTypeToGoType {
		supportedSolidityTypes = append(supportedSolidityTypes, k)
	}
}

// Note we need to include the data prefix handling here because
// if dynamic types (such as bytes) are used, the offset will be affected.
func getTxDataUsingABIEncoding(encodingSpec []string, jsonValues []gjson.Result) ([]byte, error) {
	var arguments abi.Arguments
	if len(jsonValues) != len(encodingSpec) {
		return nil, errors.Errorf("number of collectors %d != number of types in ABI encoding %d", len(jsonValues), len(encodingSpec))
	}
	var values = make([]interface{}, len(jsonValues))
	for i, argType := range encodingSpec {
		if _, supported := solidityTypeToGoType[argType]; !supported {
			return nil, errors.Wrapf(ErrInvalidABIEncoding, "%v is unsupported, supported types are %v", argType, supportedSolidityTypes)
		}
		if _, ok := jsonTypes[jsonValues[i].Type][argType]; !ok {
			return nil, errors.Wrapf(ErrInvalidABIEncoding, "can't convert %v to %v", jsonValues[i].Value(), argType)
		}
		t, err := abi.NewType(argType, "", nil)
		if err != nil {
			return nil, errors.Errorf("err %v on arg type %s index %d", err, argType, i)
		}
		arguments = append(arguments, abi.Argument{
			Type: t,
		})

		switch jsonValues[i].Type {
		case gjson.String:
			// Only supports hex strings.
			b, err := hexutil.Decode(jsonValues[i].String())
			if err != nil {
				return nil, errors.Wrapf(ErrInvalidABIEncoding, "can't convert %v to %v, bytes should be 0x-prefixed hex strings", jsonValues[i].Value(), argType)
			}
			if argType == "bytes32" {
				if len(b) != 32 {
					return nil, errors.Wrapf(ErrInvalidABIEncoding, "can't convert %v to %v", jsonValues[i].Value(), argType)
				}
				var arg [32]byte
				copy(arg[:], b)
				values[i] = arg
			} else if argType == "address" {
				if !common.IsHexAddress(jsonValues[i].String()) || len(b) != 20 {
					return nil, errors.Wrapf(ErrInvalidABIEncoding, "invalid address %v", jsonValues[i].String())
				}
				values[i] = common.HexToAddress(jsonValues[i].String())
			} else if argType == "bytes4" {
				if len(b) != 4 {
					return nil, errors.Wrapf(ErrInvalidABIEncoding, "can't convert %v to %v", jsonValues[i].Value(), argType)
				}
				var arg [4]byte
				copy(arg[:], b)
				values[i] = arg
			} else if argType == "bytes" {
				values[i] = b
			}
		case gjson.Number:
			values[i] = big.NewInt(jsonValues[i].Int()) // JSON specs can't actually handle 256bit numbers only 64bit?
		case gjson.False, gjson.True:
			// Note we can potentially use this cast strategy to support more types
			if reflect.TypeOf(jsonValues[i].Value()).ConvertibleTo(solidityTypeToGoType[argType]) {
				values[i] = reflect.ValueOf(jsonValues[i].Value()).Convert(solidityTypeToGoType[argType]).Interface()
			} else {
				return nil, errors.Wrapf(ErrInvalidABIEncoding, "can't convert %v to %v", jsonValues[i].Value(), argType)
			}
		default:
			// Complex types, array or object. Support as needed
			return nil, errors.Wrapf(ErrInvalidABIEncoding, "can't convert %v to %v", jsonValues[i].Value(), argType)
		}
	}
	packedArgs, err := arguments.PackValues(values)
	if err != nil {
		return nil, err
	}
	return utils.ConcatBytes(packedArgs), nil
}

// getTxData returns the data to save against the callback encoded according to
// the dataFormat parameter in the job spec
func getTxData(e *EthTx, input models.RunInput) ([]byte, error) {
	result := input.Result()
	if e.DataFormat == "" {
		return common.HexToHash(result.Str).Bytes(), nil
	}

	output, err := utils.EVMTranscodeJSONWithFormat(result, e.DataFormat)
	if err != nil {
		return []byte{}, err
	}
	// If data format is "bytes" then we have dynamic types,
	// which involve specifying the location of the data portion of the arg.
	// i.e. callback(reqID bytes32, bytes arg)
	if e.DataFormat == DataFormatBytes || len(e.DataPrefix) > 0 {
		// If we do not have a data prefix (reqID), encoding is:
		// [4byte fs][0x00..20][arg 1].
		payloadOffset := utils.EVMWordUint64(utils.EVMWordByteLen)
		if len(e.DataPrefix) > 0 {
			// If we have a data prefix (reqID), encoding is:
			// [4byte fs][0x00..40][reqID][arg1]
			payloadOffset = utils.EVMWordUint64(utils.EVMWordByteLen * 2)
			return append(payloadOffset, output...), nil
		}
		return append(payloadOffset, output...), nil
	}
	return output, nil
}
