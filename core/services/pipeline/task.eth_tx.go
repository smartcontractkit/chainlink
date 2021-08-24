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

	db        *gorm.DB
	config    Config
	keyStore  ETHKeyStore
	txManager TxManager
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

func (t *ETHTxTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		fromAddrs             AddressSliceParam
		toAddr                AddressParam
		data                  BytesParam
		gasLimit              Uint64Param
		txMetaMap             MapParam
		maybeMinConfirmations MaybeUint64Param
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&fromAddrs, From(VarExpr(t.From, vars), JSONWithVarExprs(t.From, vars, false), NonemptyString(t.From), nil)), "from"),
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		errors.Wrap(ResolveParam(&gasLimit, From(VarExpr(t.GasLimit, vars), NonemptyString(t.GasLimit), t.config.EvmGasLimitDefault())), "gasLimit"),
		errors.Wrap(ResolveParam(&txMetaMap, From(VarExpr(t.TxMeta, vars), JSONWithVarExprs(t.TxMeta, vars, false), MapParam{})), "txMeta"),
		errors.Wrap(ResolveParam(&maybeMinConfirmations, From(t.MinConfirmations)), "minConfirmations"),
	)
	if err != nil {
		return Result{Error: err}
	}

	var minConfirmations uint64
	if min, isSet := maybeMinConfirmations.Uint64(); isSet {
		minConfirmations = min
	} else {
		minConfirmations = t.config.MinRequiredOutgoingConfirmations()
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
		return Result{Error: errors.Wrapf(ErrBadInput, "txMeta: %v", err)}
	}

	err = decoder.Decode(txMetaMap)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "txMeta: %v", err)}
	}

	fromAddr, err := t.keyStore.GetRoundRobinAddress(fromAddrs...)
	if err != nil {
		err = errors.Wrap(err, "ETHTxTask failed to get fromAddress")
		logger.Error(err)
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while querying keystore: %v", err)}
	}

	// NOTE: This can be easily adjusted later to allow job specs to specify the details of which strategy they would like
	strategy := bulletprooftxmanager.SendEveryStrategy{}

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

	_, err = t.txManager.CreateEthTransaction(t.db, newTx)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while creating transaction: %v", err)}
	}

	if minConfirmations > 0 {
		return Result{Error: ErrPending}
	}

	return Result{Value: nil}
}
