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
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//
// Return types:
//     nil
//
type ETHTxTask struct {
	BaseTask `mapstructure:",squash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Data     string `json:"data"`
	GasLimit string `json:"gasLimit"`
	TxMeta   string `json:"txMeta"`

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
	CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, meta interface{}, strategy bulletprooftxmanager.TxStrategy) (etx bulletprooftxmanager.EthTx, err error)
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
		fromAddrs AddressSliceParam
		toAddr    AddressParam
		data      BytesParam
		gasLimit  Uint64Param
		txMetaMap MapParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&fromAddrs, From(VarExpr(t.From, vars), JSONWithVarExprs(t.From, vars, false), NonemptyString(t.From), nil)), "from"),
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		errors.Wrap(ResolveParam(&gasLimit, From(VarExpr(t.GasLimit, vars), NonemptyString(t.GasLimit), t.config.EthGasLimitDefault())), "gasLimit"),
		errors.Wrap(ResolveParam(&txMetaMap, From(VarExpr(t.TxMeta, vars), JSONWithVarExprs(t.TxMeta, vars, false), MapParam{})), "txMeta"),
	)
	if err != nil {
		return Result{Error: err}
	}

	var txMeta models.EthTxMetaV2

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

	_, err = t.txManager.CreateEthTransaction(t.db, fromAddr, common.Address(toAddr), []byte(data), uint64(gasLimit), &txMeta, strategy)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrTaskRunFailed, "while creating transaction: %v", err)}
	}
	// TODO(spook): once @archseer's "async jobs" work is merged, return the tx hash of
	// the successful EthTxAttempt
	return Result{Value: nil}
}
