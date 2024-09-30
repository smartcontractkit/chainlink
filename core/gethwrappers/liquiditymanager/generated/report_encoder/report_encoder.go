// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package report_encoder

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type ILiquidityManagerCrossChainRebalancerArgs struct {
	RemoteRebalancer    common.Address
	LocalBridge         common.Address
	RemoteToken         common.Address
	RemoteChainSelector uint64
	Enabled             bool
}

type ILiquidityManagerLiquidityInstructions struct {
	SendLiquidityParams    []ILiquidityManagerSendLiquidityParams
	ReceiveLiquidityParams []ILiquidityManagerReceiveLiquidityParams
}

type ILiquidityManagerReceiveLiquidityParams struct {
	Amount              *big.Int
	RemoteChainSelector uint64
	ShouldWrapNative    bool
	BridgeData          []byte
}

type ILiquidityManagerSendLiquidityParams struct {
	Amount              *big.Int
	NativeBridgeFee     *big.Int
	RemoteChainSelector uint64
	BridgeData          []byte
}

var ReportEncoderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nativeBridgeFee\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"internalType\":\"structILiquidityManager.SendLiquidityParams[]\",\"name\":\"sendLiquidityParams\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"shouldWrapNative\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"internalType\":\"structILiquidityManager.ReceiveLiquidityParams[]\",\"name\":\"receiveLiquidityParams\",\"type\":\"tuple[]\"}],\"internalType\":\"structILiquidityManager.LiquidityInstructions\",\"name\":\"instructions\",\"type\":\"tuple\"}],\"name\":\"exposeForEncoding\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainRebalancers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var ReportEncoderABI = ReportEncoderMetaData.ABI

type ReportEncoder struct {
	address common.Address
	abi     abi.ABI
	ReportEncoderCaller
	ReportEncoderTransactor
	ReportEncoderFilterer
}

type ReportEncoderCaller struct {
	contract *bind.BoundContract
}

type ReportEncoderTransactor struct {
	contract *bind.BoundContract
}

type ReportEncoderFilterer struct {
	contract *bind.BoundContract
}

type ReportEncoderSession struct {
	Contract     *ReportEncoder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ReportEncoderCallerSession struct {
	Contract *ReportEncoderCaller
	CallOpts bind.CallOpts
}

type ReportEncoderTransactorSession struct {
	Contract     *ReportEncoderTransactor
	TransactOpts bind.TransactOpts
}

type ReportEncoderRaw struct {
	Contract *ReportEncoder
}

type ReportEncoderCallerRaw struct {
	Contract *ReportEncoderCaller
}

type ReportEncoderTransactorRaw struct {
	Contract *ReportEncoderTransactor
}

func NewReportEncoder(address common.Address, backend bind.ContractBackend) (*ReportEncoder, error) {
	abi, err := abi.JSON(strings.NewReader(ReportEncoderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindReportEncoder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReportEncoder{address: address, abi: abi, ReportEncoderCaller: ReportEncoderCaller{contract: contract}, ReportEncoderTransactor: ReportEncoderTransactor{contract: contract}, ReportEncoderFilterer: ReportEncoderFilterer{contract: contract}}, nil
}

func NewReportEncoderCaller(address common.Address, caller bind.ContractCaller) (*ReportEncoderCaller, error) {
	contract, err := bindReportEncoder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReportEncoderCaller{contract: contract}, nil
}

func NewReportEncoderTransactor(address common.Address, transactor bind.ContractTransactor) (*ReportEncoderTransactor, error) {
	contract, err := bindReportEncoder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReportEncoderTransactor{contract: contract}, nil
}

func NewReportEncoderFilterer(address common.Address, filterer bind.ContractFilterer) (*ReportEncoderFilterer, error) {
	contract, err := bindReportEncoder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReportEncoderFilterer{contract: contract}, nil
}

func bindReportEncoder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReportEncoderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ReportEncoder *ReportEncoderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReportEncoder.Contract.ReportEncoderCaller.contract.Call(opts, result, method, params...)
}

func (_ReportEncoder *ReportEncoderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReportEncoder.Contract.ReportEncoderTransactor.contract.Transfer(opts)
}

func (_ReportEncoder *ReportEncoderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReportEncoder.Contract.ReportEncoderTransactor.contract.Transact(opts, method, params...)
}

func (_ReportEncoder *ReportEncoderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReportEncoder.Contract.contract.Call(opts, result, method, params...)
}

func (_ReportEncoder *ReportEncoderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReportEncoder.Contract.contract.Transfer(opts)
}

func (_ReportEncoder *ReportEncoderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReportEncoder.Contract.contract.Transact(opts, method, params...)
}

func (_ReportEncoder *ReportEncoderCaller) ExposeForEncoding(opts *bind.CallOpts, instructions ILiquidityManagerLiquidityInstructions) error {
	var out []interface{}
	err := _ReportEncoder.contract.Call(opts, &out, "exposeForEncoding", instructions)

	if err != nil {
		return err
	}

	return err

}

func (_ReportEncoder *ReportEncoderSession) ExposeForEncoding(instructions ILiquidityManagerLiquidityInstructions) error {
	return _ReportEncoder.Contract.ExposeForEncoding(&_ReportEncoder.CallOpts, instructions)
}

func (_ReportEncoder *ReportEncoderCallerSession) ExposeForEncoding(instructions ILiquidityManagerLiquidityInstructions) error {
	return _ReportEncoder.Contract.ExposeForEncoding(&_ReportEncoder.CallOpts, instructions)
}

func (_ReportEncoder *ReportEncoderCaller) GetAllCrossChainRebalancers(opts *bind.CallOpts) ([]ILiquidityManagerCrossChainRebalancerArgs, error) {
	var out []interface{}
	err := _ReportEncoder.contract.Call(opts, &out, "getAllCrossChainRebalancers")

	if err != nil {
		return *new([]ILiquidityManagerCrossChainRebalancerArgs), err
	}

	out0 := *abi.ConvertType(out[0], new([]ILiquidityManagerCrossChainRebalancerArgs)).(*[]ILiquidityManagerCrossChainRebalancerArgs)

	return out0, err

}

func (_ReportEncoder *ReportEncoderSession) GetAllCrossChainRebalancers() ([]ILiquidityManagerCrossChainRebalancerArgs, error) {
	return _ReportEncoder.Contract.GetAllCrossChainRebalancers(&_ReportEncoder.CallOpts)
}

func (_ReportEncoder *ReportEncoderCallerSession) GetAllCrossChainRebalancers() ([]ILiquidityManagerCrossChainRebalancerArgs, error) {
	return _ReportEncoder.Contract.GetAllCrossChainRebalancers(&_ReportEncoder.CallOpts)
}

func (_ReportEncoder *ReportEncoderCaller) GetLiquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReportEncoder.contract.Call(opts, &out, "getLiquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ReportEncoder *ReportEncoderSession) GetLiquidity() (*big.Int, error) {
	return _ReportEncoder.Contract.GetLiquidity(&_ReportEncoder.CallOpts)
}

func (_ReportEncoder *ReportEncoderCallerSession) GetLiquidity() (*big.Int, error) {
	return _ReportEncoder.Contract.GetLiquidity(&_ReportEncoder.CallOpts)
}

func (_ReportEncoder *ReportEncoder) Address() common.Address {
	return _ReportEncoder.address
}

type ReportEncoderInterface interface {
	ExposeForEncoding(opts *bind.CallOpts, instructions ILiquidityManagerLiquidityInstructions) error

	GetAllCrossChainRebalancers(opts *bind.CallOpts) ([]ILiquidityManagerCrossChainRebalancerArgs, error)

	GetLiquidity(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
