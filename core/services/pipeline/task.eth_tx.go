package pipeline

import (
	"context"
	"math/big"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
)

//
// Return types:
//     nil
//
type ETHTxTask struct {
	BaseTask         `mapstructure:",squash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Data             string `json:"data"`
	GasLimit         string `json:"gasLimit"`
	TxMeta           string `json:"txMeta"`
	MinConfirmations string `json:"minConfirmations"`
	EVMChainID       string `json:"evmChainID" mapstructure:"evmChainID"`
	TransmitChecker  string `json:"transmitChecker"`

	keyStore ETHKeyStore
	chainSet evm.ChainSet
}

//go:generate mockery --name ETHKeyStore --output ./mocks/ --case=underscore

type ETHKeyStore interface {
	GetRoundRobinAddress(chainID *big.Int, addrs ...common.Address) (common.Address, error)
}

var _ Task = (*ETHTxTask)(nil)

func (t *ETHTxTask) Type() TaskType {
	return TaskTypeETHTx
}

func (t *ETHTxTask) Run(_ context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var chainID StringParam
	err := errors.Wrap(ResolveParam(&chainID, From(VarExpr(t.EVMChainID, vars), NonemptyString(t.EVMChainID), "")), "evmChainID")
	if err != nil {
		return Result{Error: err}, runInfo
	}

	chain, err := getChainByString(t.chainSet, string(chainID))
	if err != nil {
		return Result{Error: errors.Wrapf(err, "failed to get chain by id: %v", t.EVMChainID)}, retryableRunInfo()
	}
	cfg := chain.Config()
	txManager := chain.TxManager()
	_, err = CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		fromAddrs             AddressSliceParam
		toAddr                AddressParam
		data                  BytesParam
		gasLimit              Uint64Param
		txMetaMap             MapParam
		maybeMinConfirmations MaybeUint64Param
		transmitCheckerMap    MapParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&fromAddrs, From(VarExpr(t.From, vars), JSONWithVarExprs(t.From, vars, false), NonemptyString(t.From), nil)), "from"),
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		errors.Wrap(ResolveParam(&gasLimit, From(VarExpr(t.GasLimit, vars), NonemptyString(t.GasLimit), cfg.EvmGasLimitDefault())), "gasLimit"),
		errors.Wrap(ResolveParam(&txMetaMap, From(VarExpr(t.TxMeta, vars), JSONWithVarExprs(t.TxMeta, vars, false), MapParam{})), "txMeta"),
		errors.Wrap(ResolveParam(&maybeMinConfirmations, From(t.MinConfirmations)), "minConfirmations"),
		errors.Wrap(ResolveParam(&transmitCheckerMap, From(VarExpr(t.TransmitChecker, vars), JSONWithVarExprs(t.TransmitChecker, vars, false), MapParam{})), "transmitChecker"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	var minOutgoingConfirmations uint64
	if min, isSet := maybeMinConfirmations.Uint64(); isSet {
		minOutgoingConfirmations = min
	} else {
		minOutgoingConfirmations = cfg.MinRequiredOutgoingConfirmations()
	}

	txMeta, err := decodeMeta(txMetaMap)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	transmitChecker, err := decodeTransmitChecker(transmitCheckerMap)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	fromAddr, err := t.keyStore.GetRoundRobinAddress(chain.ID(), fromAddrs...)
	if err != nil {
		err = errors.Wrap(err, "ETHTxTask failed to get fromAddress")
		lggr.Error(err)
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while querying keystore: %v", err)}, retryableRunInfo()
	}

	// NOTE: This can be easily adjusted later to allow job specs to specify the details of which strategy they would like
	strategy := txmgr.NewSendEveryStrategy()

	newTx := txmgr.NewTx{
		FromAddress:    fromAddr,
		ToAddress:      common.Address(toAddr),
		EncodedPayload: []byte(data),
		GasLimit:       uint64(gasLimit),
		Meta:           txMeta,
		Strategy:       strategy,
		Checker:        transmitChecker,
	}

	if minOutgoingConfirmations > 0 {
		// Store the task run ID, so we can resume the pipeline when tx is confirmed
		newTx.PipelineTaskRunID = &t.uuid
		newTx.MinConfirmations = null.Uint32From(uint32(minOutgoingConfirmations))
	}

	_, err = txManager.CreateEthTransaction(newTx)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while creating transaction: %v", err)}, retryableRunInfo()
	}

	if minOutgoingConfirmations > 0 {
		return Result{}, pendingRunInfo()
	}

	return Result{Value: nil}, runInfo
}

func decodeMeta(metaMap MapParam) (*txmgr.EthTxMeta, error) {
	var txMeta txmgr.EthTxMeta
	metaDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      &txMeta,
		ErrorUnused: true,
		DecodeHook: func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
			switch from {
			case stringType:
				switch to {
				case int32Type:
					i, err2 := strconv.ParseInt(data.(string), 10, 32)
					return int32(i), err2
				case reflect.TypeOf(common.Hash{}):
					return common.HexToHash(data.(string)), nil
				}
			}
			return data, nil
		},
	})
	if err != nil {
		return &txMeta, errors.Wrapf(ErrBadInput, "txMeta: %v", err)
	}

	err = metaDecoder.Decode(metaMap)
	if err != nil {
		return &txMeta, errors.Wrapf(ErrBadInput, "txMeta: %v", err)
	}
	return &txMeta, nil
}

func decodeTransmitChecker(checkerMap MapParam) (txmgr.TransmitCheckerSpec, error) {
	var transmitChecker txmgr.TransmitCheckerSpec
	checkerDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      &transmitChecker,
		ErrorUnused: true,
		DecodeHook: func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
			switch from {
			case stringType:
				switch to {
				case reflect.TypeOf(common.Address{}):
					return common.HexToAddress(data.(string)), nil
				}
			}
			return data, nil
		},
	})
	if err != nil {
		return transmitChecker, errors.Wrapf(ErrBadInput, "transmitChecker: %v", err)
	}

	err = checkerDecoder.Decode(checkerMap)
	if err != nil {
		return transmitChecker, errors.Wrapf(ErrBadInput, "transmitChecker: %v", err)
	}
	return transmitChecker, nil
}
