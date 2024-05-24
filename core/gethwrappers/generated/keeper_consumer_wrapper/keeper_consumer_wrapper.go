// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_consumer_wrapper

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

var KeeperConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"updateInterval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimeStamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161035a38038061035a83398101604081905261002f9161003f565b6080524260015560008055610058565b60006020828403121561005157600080fd5b5051919050565b6080516102e8610072600039600060cc01526102e86000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806361bc221a1161005057806361bc221a1461009d5780636e04ff0d146100a6578063947a36fb146100c757600080fd5b80633f3b3b271461006c5780634585e33b14610088575b600080fd5b61007560015481565b6040519081526020015b60405180910390f35b61009b6100963660046101b3565b6100ee565b005b61007560005481565b6100b96100b43660046101b3565b610103565b60405161007f929190610225565b6100757f000000000000000000000000000000000000000000000000000000000000000081565b6000546100fc90600161029b565b6000555050565b6000606061010f610157565b6001848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b321580159061017a57503273111111111111111111111111111111111111111114155b156101b1576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b600080602083850312156101c657600080fd5b823567ffffffffffffffff808211156101de57600080fd5b818501915085601f8301126101f257600080fd5b81358181111561020157600080fd5b86602082850101111561021357600080fd5b60209290920196919550909350505050565b821515815260006020604081840152835180604085015260005b8181101561025b5785810183015185820160600152820161023f565b5060006060828601015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509392505050565b808201808211156102d5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000810000a",
}

var KeeperConsumerABI = KeeperConsumerMetaData.ABI

var KeeperConsumerBin = KeeperConsumerMetaData.Bin

func DeployKeeperConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, updateInterval *big.Int) (common.Address, *types.Transaction, *KeeperConsumer, error) {
	parsed, err := KeeperConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperConsumerBin), backend, updateInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperConsumer{address: address, abi: *parsed, KeeperConsumerCaller: KeeperConsumerCaller{contract: contract}, KeeperConsumerTransactor: KeeperConsumerTransactor{contract: contract}, KeeperConsumerFilterer: KeeperConsumerFilterer{contract: contract}}, nil
}

type KeeperConsumer struct {
	address common.Address
	abi     abi.ABI
	KeeperConsumerCaller
	KeeperConsumerTransactor
	KeeperConsumerFilterer
}

type KeeperConsumerCaller struct {
	contract *bind.BoundContract
}

type KeeperConsumerTransactor struct {
	contract *bind.BoundContract
}

type KeeperConsumerFilterer struct {
	contract *bind.BoundContract
}

type KeeperConsumerSession struct {
	Contract     *KeeperConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperConsumerCallerSession struct {
	Contract *KeeperConsumerCaller
	CallOpts bind.CallOpts
}

type KeeperConsumerTransactorSession struct {
	Contract     *KeeperConsumerTransactor
	TransactOpts bind.TransactOpts
}

type KeeperConsumerRaw struct {
	Contract *KeeperConsumer
}

type KeeperConsumerCallerRaw struct {
	Contract *KeeperConsumerCaller
}

type KeeperConsumerTransactorRaw struct {
	Contract *KeeperConsumerTransactor
}

func NewKeeperConsumer(address common.Address, backend bind.ContractBackend) (*KeeperConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumer{address: address, abi: abi, KeeperConsumerCaller: KeeperConsumerCaller{contract: contract}, KeeperConsumerTransactor: KeeperConsumerTransactor{contract: contract}, KeeperConsumerFilterer: KeeperConsumerFilterer{contract: contract}}, nil
}

func NewKeeperConsumerCaller(address common.Address, caller bind.ContractCaller) (*KeeperConsumerCaller, error) {
	contract, err := bindKeeperConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerCaller{contract: contract}, nil
}

func NewKeeperConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperConsumerTransactor, error) {
	contract, err := bindKeeperConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerTransactor{contract: contract}, nil
}

func NewKeeperConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperConsumerFilterer, error) {
	contract, err := bindKeeperConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerFilterer{contract: contract}, nil
}

func bindKeeperConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperConsumer *KeeperConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumer.Contract.KeeperConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperConsumer *KeeperConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.KeeperConsumerTransactor.contract.Transfer(opts)
}

func (_KeeperConsumer *KeeperConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.KeeperConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperConsumer *KeeperConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperConsumer *KeeperConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.contract.Transfer(opts)
}

func (_KeeperConsumer *KeeperConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperConsumer *KeeperConsumerCaller) CheckUpkeep(opts *bind.CallOpts, checkData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "checkUpkeep", checkData)

	outstruct := new(CheckUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_KeeperConsumer *KeeperConsumerSession) CheckUpkeep(checkData []byte) (CheckUpkeep,

	error) {
	return _KeeperConsumer.Contract.CheckUpkeep(&_KeeperConsumer.CallOpts, checkData)
}

func (_KeeperConsumer *KeeperConsumerCallerSession) CheckUpkeep(checkData []byte) (CheckUpkeep,

	error) {
	return _KeeperConsumer.Contract.CheckUpkeep(&_KeeperConsumer.CallOpts, checkData)
}

func (_KeeperConsumer *KeeperConsumerCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumer *KeeperConsumerSession) Counter() (*big.Int, error) {
	return _KeeperConsumer.Contract.Counter(&_KeeperConsumer.CallOpts)
}

func (_KeeperConsumer *KeeperConsumerCallerSession) Counter() (*big.Int, error) {
	return _KeeperConsumer.Contract.Counter(&_KeeperConsumer.CallOpts)
}

func (_KeeperConsumer *KeeperConsumerCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumer *KeeperConsumerSession) Interval() (*big.Int, error) {
	return _KeeperConsumer.Contract.Interval(&_KeeperConsumer.CallOpts)
}

func (_KeeperConsumer *KeeperConsumerCallerSession) Interval() (*big.Int, error) {
	return _KeeperConsumer.Contract.Interval(&_KeeperConsumer.CallOpts)
}

func (_KeeperConsumer *KeeperConsumerCaller) LastTimeStamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "lastTimeStamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumer *KeeperConsumerSession) LastTimeStamp() (*big.Int, error) {
	return _KeeperConsumer.Contract.LastTimeStamp(&_KeeperConsumer.CallOpts)
}

func (_KeeperConsumer *KeeperConsumerCallerSession) LastTimeStamp() (*big.Int, error) {
	return _KeeperConsumer.Contract.LastTimeStamp(&_KeeperConsumer.CallOpts)
}

func (_KeeperConsumer *KeeperConsumerTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.contract.Transact(opts, "performUpkeep", performData)
}

func (_KeeperConsumer *KeeperConsumerSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.PerformUpkeep(&_KeeperConsumer.TransactOpts, performData)
}

func (_KeeperConsumer *KeeperConsumerTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.PerformUpkeep(&_KeeperConsumer.TransactOpts, performData)
}

type CheckUpkeep struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_KeeperConsumer *KeeperConsumer) Address() common.Address {
	return _KeeperConsumer.address
}

type KeeperConsumerInterface interface {
	CheckUpkeep(opts *bind.CallOpts, checkData []byte) (CheckUpkeep,

		error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastTimeStamp(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	Address() common.Address
}
