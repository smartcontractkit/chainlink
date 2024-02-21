// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rebalancer_report_encoder

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

type IRebalancerCrossChainRebalancerArgs struct {
	RemoteRebalancer    common.Address
	LocalBridge         common.Address
	RemoteToken         common.Address
	RemoteChainSelector uint64
	Enabled             bool
}

type IRebalancerLiquidityInstructions struct {
	SendLiquidityParams    []IRebalancerSendLiquidityParams
	ReceiveLiquidityParams []IRebalancerReceiveLiquidityParams
}

type IRebalancerReceiveLiquidityParams struct {
	Amount              *big.Int
	RemoteChainSelector uint64
	BridgeData          []byte
}

type IRebalancerSendLiquidityParams struct {
	Amount              *big.Int
	NativeBridgeFee     *big.Int
	RemoteChainSelector uint64
	BridgeData          []byte
}

var RebalancerReportEncoderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nativeBridgeFee\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"internalType\":\"structIRebalancer.SendLiquidityParams[]\",\"name\":\"sendLiquidityParams\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"internalType\":\"structIRebalancer.ReceiveLiquidityParams[]\",\"name\":\"receiveLiquidityParams\",\"type\":\"tuple[]\"}],\"internalType\":\"structIRebalancer.LiquidityInstructions\",\"name\":\"instructions\",\"type\":\"tuple\"}],\"name\":\"exposeForEncoding\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainRebalancers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structIRebalancer.CrossChainRebalancerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var RebalancerReportEncoderABI = RebalancerReportEncoderMetaData.ABI

type RebalancerReportEncoder struct {
	address common.Address
	abi     abi.ABI
	RebalancerReportEncoderCaller
	RebalancerReportEncoderTransactor
	RebalancerReportEncoderFilterer
}

type RebalancerReportEncoderCaller struct {
	contract *bind.BoundContract
}

type RebalancerReportEncoderTransactor struct {
	contract *bind.BoundContract
}

type RebalancerReportEncoderFilterer struct {
	contract *bind.BoundContract
}

type RebalancerReportEncoderSession struct {
	Contract     *RebalancerReportEncoder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RebalancerReportEncoderCallerSession struct {
	Contract *RebalancerReportEncoderCaller
	CallOpts bind.CallOpts
}

type RebalancerReportEncoderTransactorSession struct {
	Contract     *RebalancerReportEncoderTransactor
	TransactOpts bind.TransactOpts
}

type RebalancerReportEncoderRaw struct {
	Contract *RebalancerReportEncoder
}

type RebalancerReportEncoderCallerRaw struct {
	Contract *RebalancerReportEncoderCaller
}

type RebalancerReportEncoderTransactorRaw struct {
	Contract *RebalancerReportEncoderTransactor
}

func NewRebalancerReportEncoder(address common.Address, backend bind.ContractBackend) (*RebalancerReportEncoder, error) {
	abi, err := abi.JSON(strings.NewReader(RebalancerReportEncoderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRebalancerReportEncoder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RebalancerReportEncoder{address: address, abi: abi, RebalancerReportEncoderCaller: RebalancerReportEncoderCaller{contract: contract}, RebalancerReportEncoderTransactor: RebalancerReportEncoderTransactor{contract: contract}, RebalancerReportEncoderFilterer: RebalancerReportEncoderFilterer{contract: contract}}, nil
}

func NewRebalancerReportEncoderCaller(address common.Address, caller bind.ContractCaller) (*RebalancerReportEncoderCaller, error) {
	contract, err := bindRebalancerReportEncoder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RebalancerReportEncoderCaller{contract: contract}, nil
}

func NewRebalancerReportEncoderTransactor(address common.Address, transactor bind.ContractTransactor) (*RebalancerReportEncoderTransactor, error) {
	contract, err := bindRebalancerReportEncoder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RebalancerReportEncoderTransactor{contract: contract}, nil
}

func NewRebalancerReportEncoderFilterer(address common.Address, filterer bind.ContractFilterer) (*RebalancerReportEncoderFilterer, error) {
	contract, err := bindRebalancerReportEncoder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RebalancerReportEncoderFilterer{contract: contract}, nil
}

func bindRebalancerReportEncoder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RebalancerReportEncoderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RebalancerReportEncoder *RebalancerReportEncoderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RebalancerReportEncoder.Contract.RebalancerReportEncoderCaller.contract.Call(opts, result, method, params...)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RebalancerReportEncoder.Contract.RebalancerReportEncoderTransactor.contract.Transfer(opts)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RebalancerReportEncoder.Contract.RebalancerReportEncoderTransactor.contract.Transact(opts, method, params...)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RebalancerReportEncoder.Contract.contract.Call(opts, result, method, params...)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RebalancerReportEncoder.Contract.contract.Transfer(opts)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RebalancerReportEncoder.Contract.contract.Transact(opts, method, params...)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCaller) ExposeForEncoding(opts *bind.CallOpts, instructions IRebalancerLiquidityInstructions) error {
	var out []interface{}
	err := _RebalancerReportEncoder.contract.Call(opts, &out, "exposeForEncoding", instructions)

	if err != nil {
		return err
	}

	return err

}

func (_RebalancerReportEncoder *RebalancerReportEncoderSession) ExposeForEncoding(instructions IRebalancerLiquidityInstructions) error {
	return _RebalancerReportEncoder.Contract.ExposeForEncoding(&_RebalancerReportEncoder.CallOpts, instructions)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCallerSession) ExposeForEncoding(instructions IRebalancerLiquidityInstructions) error {
	return _RebalancerReportEncoder.Contract.ExposeForEncoding(&_RebalancerReportEncoder.CallOpts, instructions)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCaller) GetAllCrossChainRebalancers(opts *bind.CallOpts) ([]IRebalancerCrossChainRebalancerArgs, error) {
	var out []interface{}
	err := _RebalancerReportEncoder.contract.Call(opts, &out, "getAllCrossChainRebalancers")

	if err != nil {
		return *new([]IRebalancerCrossChainRebalancerArgs), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRebalancerCrossChainRebalancerArgs)).(*[]IRebalancerCrossChainRebalancerArgs)

	return out0, err

}

func (_RebalancerReportEncoder *RebalancerReportEncoderSession) GetAllCrossChainRebalancers() ([]IRebalancerCrossChainRebalancerArgs, error) {
	return _RebalancerReportEncoder.Contract.GetAllCrossChainRebalancers(&_RebalancerReportEncoder.CallOpts)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCallerSession) GetAllCrossChainRebalancers() ([]IRebalancerCrossChainRebalancerArgs, error) {
	return _RebalancerReportEncoder.Contract.GetAllCrossChainRebalancers(&_RebalancerReportEncoder.CallOpts)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCaller) GetLiquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RebalancerReportEncoder.contract.Call(opts, &out, "getLiquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_RebalancerReportEncoder *RebalancerReportEncoderSession) GetLiquidity() (*big.Int, error) {
	return _RebalancerReportEncoder.Contract.GetLiquidity(&_RebalancerReportEncoder.CallOpts)
}

func (_RebalancerReportEncoder *RebalancerReportEncoderCallerSession) GetLiquidity() (*big.Int, error) {
	return _RebalancerReportEncoder.Contract.GetLiquidity(&_RebalancerReportEncoder.CallOpts)
}

func (_RebalancerReportEncoder *RebalancerReportEncoder) Address() common.Address {
	return _RebalancerReportEncoder.address
}

type RebalancerReportEncoderInterface interface {
	ExposeForEncoding(opts *bind.CallOpts, instructions IRebalancerLiquidityInstructions) error

	GetAllCrossChainRebalancers(opts *bind.CallOpts) ([]IRebalancerCrossChainRebalancerArgs, error)

	GetLiquidity(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
