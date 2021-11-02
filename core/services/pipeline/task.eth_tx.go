package pipeline

import (
	"context"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
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
	Simulate         string `json:"simulate" mapstructure:"simulate"`

	keyStore ETHKeyStore
	chainSet evm.ChainSet
}

//go:generate mockery --name ETHKeyStore --output ./mocks/ --case=underscore
//go:generate mockery --name TxManager --output ./mocks/ --case=underscore

type ETHKeyStore interface {
	GetRoundRobinAddress(addrs ...common.Address) (common.Address, error)
}

type TxManager interface {
	CreateEthTransaction(db *gorm.DB, newTx bulletprooftxmanager.NewTx) (etx bulletprooftxmanager.EthTx, err error)
}

var _ Task = (*ETHTxTask)(nil)

func (t *ETHTxTask) Type() TaskType {
	return TaskTypeETHTx
}

func (t *ETHTxTask) Run(_ context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	chain, err := getChainByString(t.chainSet, t.EVMChainID)
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
		simulate              BoolParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&fromAddrs, From(VarExpr(t.From, vars), JSONWithVarExprs(t.From, vars, false), NonemptyString(t.From), nil)), "from"),
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		errors.Wrap(ResolveParam(&gasLimit, From(VarExpr(t.GasLimit, vars), NonemptyString(t.GasLimit), cfg.EvmGasLimitDefault())), "gasLimit"),
		errors.Wrap(ResolveParam(&txMetaMap, From(VarExpr(t.TxMeta, vars), JSONWithVarExprs(t.TxMeta, vars, false), MapParam{})), "txMeta"),
		errors.Wrap(ResolveParam(&maybeMinConfirmations, From(t.MinConfirmations)), "minConfirmations"),
		errors.Wrap(ResolveParam(&simulate, From(VarExpr(t.Simulate, vars), NonemptyString(t.Simulate), false)), "simulate"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	var minConfirmations uint64
	if min, isSet := maybeMinConfirmations.Uint64(); isSet {
		minConfirmations = min
	} else {
		minConfirmations = cfg.MinRequiredOutgoingConfirmations()
	}

	var txMeta bulletprooftxmanager.EthTxMeta

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
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
		return Result{Error: errors.Wrapf(ErrBadInput, "txMeta: %v", err)}, runInfo
	}

	err = decoder.Decode(txMetaMap)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "txMeta: %v", err)}, runInfo
	}

	fromAddr, err := t.keyStore.GetRoundRobinAddress(fromAddrs...)
	if err != nil {
		err = errors.Wrap(err, "ETHTxTask failed to get fromAddress")
		lggr.Error(err)
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while querying keystore: %v", err)}, retryableRunInfo()
	}

	// NOTE: This can be easily adjusted later to allow job specs to specify the details of which strategy they would like
	strategy := bulletprooftxmanager.NewSendEveryStrategy(bool(simulate))

	newTx := bulletprooftxmanager.NewTx{
		FromAddress:    fromAddr,
		ToAddress:      common.Address(toAddr),
		EncodedPayload: []byte(data),
		GasLimit:       uint64(gasLimit),
		Meta:           &txMeta,
		Strategy:       strategy,
	}

	if minConfirmations > 0 {
		// Store the task run ID so we can resume the pipeline when tx is confirmed
		newTx.PipelineTaskRunID = &t.uuid
		newTx.MinConfirmations = null.Uint32From(uint32(minConfirmations))
	}

	_, err = txManager.CreateEthTransaction(newTx)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while creating transaction: %v", err)}, retryableRunInfo()
	}

	if minConfirmations > 0 {
		return Result{}, pendingRunInfo()
	}

	return Result{Value: nil}, runInfo
}
