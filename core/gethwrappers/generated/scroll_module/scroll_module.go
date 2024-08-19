// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package scroll_module

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var ScrollModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"InvalidL1FeeCoefficient\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"L1FeeCoefficientSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_l1FeeCoefficient\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"setL1FeeCalculation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805460ff60a01b1916601960a21b17905534801561002357600080fd5b50338060008161007a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100aa576100aa816100b2565b50505061015b565b336001600160a01b0382160361010a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610071565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61093c8061016a6000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638da5cb5b11610076578063de9ee35e1161005b578063de9ee35e1461017e578063f22559a014610194578063f2fde38b146101a757600080fd5b80638da5cb5b1461011f57806390bd5c741461014757600080fd5b80637810d12a116100a75780637810d12a146100ef57806379ba50971461010257806385df51fd1461010c57600080fd5b806312544140146100c357806357e871e7146100e9575b600080fd5b6100d66100d13660046106b1565b6101ba565b6040519081526020015b60405180910390f35b436100d6565b6100d66100fd3660046106b1565b6101cb565b61010a61020a565b005b6100d661011a3660046106b1565b61030c565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e0565b60015461016c9074010000000000000000000000000000000000000000900460ff1681565b60405160ff90911681526020016100e0565b6040805161afc8815260aa6020820152016100e0565b61010a6101a23660046106ca565b610339565b61010a6101b53660046106ed565b610404565b60006101c582610418565b92915050565b600060646101d883610418565b600154610200919074010000000000000000000000000000000000000000900460ff16610752565b6101c59190610769565b60015473ffffffffffffffffffffffffffffffffffffffff163314610290576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60004382101580610327575061010061032583436107a4565b115b1561033457506000919050565b504090565b610341610539565b60648160ff161115610384576040517f1a8a06a000000000000000000000000000000000000000000000000000000000815260ff82166004820152602401610287565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000060ff8416908102919091179091556040519081527f29ec9e31de0d3fe0208a7ccb792bbc26a854f123146110daa3a77219cb74a5549060200160405180910390a150565b61040c610539565b610415816105bc565b50565b600080610426836004610752565b67ffffffffffffffff81111561043e5761043e6107b7565b6040519080825280601f01601f191660200182016040528015610468576020820181803683370190505b50905073530000000000000000000000000000000000000273ffffffffffffffffffffffffffffffffffffffff166349948e0e826040518060c00160405280608c81526020016108a4608c91396040516020016104c692919061080a565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016104f19190610839565b602060405180830381865afa15801561050e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610532919061088a565b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610287565b565b3373ffffffffffffffffffffffffffffffffffffffff82160361063b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610287565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156106c357600080fd5b5035919050565b6000602082840312156106dc57600080fd5b813560ff8116811461053257600080fd5b6000602082840312156106ff57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461053257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176101c5576101c5610723565b60008261079f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b818103818111156101c5576101c5610723565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60005b838110156108015781810151838201526020016107e9565b50506000910152565b6000835161081c8184602088016107e6565b8351908301906108308183602088016107e6565b01949350505050565b60208152600082518060208401526108588160408501602087016107e6565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60006020828403121561089c57600080fd5b505191905056feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
}

var ScrollModuleABI = ScrollModuleMetaData.ABI

var ScrollModuleBin = ScrollModuleMetaData.Bin

func DeployScrollModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ScrollModule, error) {
	parsed, err := ScrollModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ScrollModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ScrollModule{address: address, abi: *parsed, ScrollModuleCaller: ScrollModuleCaller{contract: contract}, ScrollModuleTransactor: ScrollModuleTransactor{contract: contract}, ScrollModuleFilterer: ScrollModuleFilterer{contract: contract}}, nil
}

type ScrollModule struct {
	address common.Address
	abi     abi.ABI
	ScrollModuleCaller
	ScrollModuleTransactor
	ScrollModuleFilterer
}

type ScrollModuleCaller struct {
	contract *bind.BoundContract
}

type ScrollModuleTransactor struct {
	contract *bind.BoundContract
}

type ScrollModuleFilterer struct {
	contract *bind.BoundContract
}

type ScrollModuleSession struct {
	Contract     *ScrollModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ScrollModuleCallerSession struct {
	Contract *ScrollModuleCaller
	CallOpts bind.CallOpts
}

type ScrollModuleTransactorSession struct {
	Contract     *ScrollModuleTransactor
	TransactOpts bind.TransactOpts
}

type ScrollModuleRaw struct {
	Contract *ScrollModule
}

type ScrollModuleCallerRaw struct {
	Contract *ScrollModuleCaller
}

type ScrollModuleTransactorRaw struct {
	Contract *ScrollModuleTransactor
}

func NewScrollModule(address common.Address, backend bind.ContractBackend) (*ScrollModule, error) {
	abi, err := abi.JSON(strings.NewReader(ScrollModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindScrollModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ScrollModule{address: address, abi: abi, ScrollModuleCaller: ScrollModuleCaller{contract: contract}, ScrollModuleTransactor: ScrollModuleTransactor{contract: contract}, ScrollModuleFilterer: ScrollModuleFilterer{contract: contract}}, nil
}

func NewScrollModuleCaller(address common.Address, caller bind.ContractCaller) (*ScrollModuleCaller, error) {
	contract, err := bindScrollModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleCaller{contract: contract}, nil
}

func NewScrollModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*ScrollModuleTransactor, error) {
	contract, err := bindScrollModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleTransactor{contract: contract}, nil
}

func NewScrollModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*ScrollModuleFilterer, error) {
	contract, err := bindScrollModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleFilterer{contract: contract}, nil
}

func bindScrollModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ScrollModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ScrollModule *ScrollModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScrollModule.Contract.ScrollModuleCaller.contract.Call(opts, result, method, params...)
}

func (_ScrollModule *ScrollModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.Contract.ScrollModuleTransactor.contract.Transfer(opts)
}

func (_ScrollModule *ScrollModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScrollModule.Contract.ScrollModuleTransactor.contract.Transact(opts, method, params...)
}

func (_ScrollModule *ScrollModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScrollModule.Contract.contract.Call(opts, result, method, params...)
}

func (_ScrollModule *ScrollModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.Contract.contract.Transfer(opts)
}

func (_ScrollModule *ScrollModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScrollModule.Contract.contract.Transact(opts, method, params...)
}

func (_ScrollModule *ScrollModuleCaller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ScrollModule.Contract.BlockHash(&_ScrollModule.CallOpts, n)
}

func (_ScrollModule *ScrollModuleCallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ScrollModule.Contract.BlockHash(&_ScrollModule.CallOpts, n)
}

func (_ScrollModule *ScrollModuleCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) BlockNumber() (*big.Int, error) {
	return _ScrollModule.Contract.BlockNumber(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) BlockNumber() (*big.Int, error) {
	return _ScrollModule.Contract.BlockNumber(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getCurrentL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) GetCurrentL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetCurrentL1Fee(&_ScrollModule.CallOpts, dataSize)
}

func (_ScrollModule *ScrollModuleCallerSession) GetCurrentL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetCurrentL1Fee(&_ScrollModule.CallOpts, dataSize)
}

func (_ScrollModule *ScrollModuleCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_ScrollModule *ScrollModuleSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ScrollModule.Contract.GetGasOverhead(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ScrollModule.Contract.GetGasOverhead(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetMaxL1Fee(&_ScrollModule.CallOpts, dataSize)
}

func (_ScrollModule *ScrollModuleCallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetMaxL1Fee(&_ScrollModule.CallOpts, dataSize)
}

func (_ScrollModule *ScrollModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) Owner() (common.Address, error) {
	return _ScrollModule.Contract.Owner(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) Owner() (common.Address, error) {
	return _ScrollModule.Contract.Owner(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) SL1FeeCoefficient(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "s_l1FeeCoefficient")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) SL1FeeCoefficient() (uint8, error) {
	return _ScrollModule.Contract.SL1FeeCoefficient(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) SL1FeeCoefficient() (uint8, error) {
	return _ScrollModule.Contract.SL1FeeCoefficient(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.contract.Transact(opts, "acceptOwnership")
}

func (_ScrollModule *ScrollModuleSession) AcceptOwnership() (*types.Transaction, error) {
	return _ScrollModule.Contract.AcceptOwnership(&_ScrollModule.TransactOpts)
}

func (_ScrollModule *ScrollModuleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ScrollModule.Contract.AcceptOwnership(&_ScrollModule.TransactOpts)
}

func (_ScrollModule *ScrollModuleTransactor) SetL1FeeCalculation(opts *bind.TransactOpts, coefficient uint8) (*types.Transaction, error) {
	return _ScrollModule.contract.Transact(opts, "setL1FeeCalculation", coefficient)
}

func (_ScrollModule *ScrollModuleSession) SetL1FeeCalculation(coefficient uint8) (*types.Transaction, error) {
	return _ScrollModule.Contract.SetL1FeeCalculation(&_ScrollModule.TransactOpts, coefficient)
}

func (_ScrollModule *ScrollModuleTransactorSession) SetL1FeeCalculation(coefficient uint8) (*types.Transaction, error) {
	return _ScrollModule.Contract.SetL1FeeCalculation(&_ScrollModule.TransactOpts, coefficient)
}

func (_ScrollModule *ScrollModuleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ScrollModule.contract.Transact(opts, "transferOwnership", to)
}

func (_ScrollModule *ScrollModuleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ScrollModule.Contract.TransferOwnership(&_ScrollModule.TransactOpts, to)
}

func (_ScrollModule *ScrollModuleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ScrollModule.Contract.TransferOwnership(&_ScrollModule.TransactOpts, to)
}

type ScrollModuleL1FeeCoefficientSetIterator struct {
	Event *ScrollModuleL1FeeCoefficientSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ScrollModuleL1FeeCoefficientSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScrollModuleL1FeeCoefficientSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(ScrollModuleL1FeeCoefficientSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *ScrollModuleL1FeeCoefficientSetIterator) Error() error {
	return it.fail
}

func (it *ScrollModuleL1FeeCoefficientSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ScrollModuleL1FeeCoefficientSet struct {
	Coefficient uint8
	Raw         types.Log
}

func (_ScrollModule *ScrollModuleFilterer) FilterL1FeeCoefficientSet(opts *bind.FilterOpts) (*ScrollModuleL1FeeCoefficientSetIterator, error) {

	logs, sub, err := _ScrollModule.contract.FilterLogs(opts, "L1FeeCoefficientSet")
	if err != nil {
		return nil, err
	}
	return &ScrollModuleL1FeeCoefficientSetIterator{contract: _ScrollModule.contract, event: "L1FeeCoefficientSet", logs: logs, sub: sub}, nil
}

func (_ScrollModule *ScrollModuleFilterer) WatchL1FeeCoefficientSet(opts *bind.WatchOpts, sink chan<- *ScrollModuleL1FeeCoefficientSet) (event.Subscription, error) {

	logs, sub, err := _ScrollModule.contract.WatchLogs(opts, "L1FeeCoefficientSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ScrollModuleL1FeeCoefficientSet)
				if err := _ScrollModule.contract.UnpackLog(event, "L1FeeCoefficientSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_ScrollModule *ScrollModuleFilterer) ParseL1FeeCoefficientSet(log types.Log) (*ScrollModuleL1FeeCoefficientSet, error) {
	event := new(ScrollModuleL1FeeCoefficientSet)
	if err := _ScrollModule.contract.UnpackLog(event, "L1FeeCoefficientSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ScrollModuleOwnershipTransferRequestedIterator struct {
	Event *ScrollModuleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ScrollModuleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScrollModuleOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(ScrollModuleOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *ScrollModuleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ScrollModuleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ScrollModuleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ScrollModule *ScrollModuleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ScrollModuleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ScrollModule.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleOwnershipTransferRequestedIterator{contract: _ScrollModule.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ScrollModule *ScrollModuleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ScrollModuleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ScrollModule.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ScrollModuleOwnershipTransferRequested)
				if err := _ScrollModule.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_ScrollModule *ScrollModuleFilterer) ParseOwnershipTransferRequested(log types.Log) (*ScrollModuleOwnershipTransferRequested, error) {
	event := new(ScrollModuleOwnershipTransferRequested)
	if err := _ScrollModule.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ScrollModuleOwnershipTransferredIterator struct {
	Event *ScrollModuleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ScrollModuleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ScrollModuleOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(ScrollModuleOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *ScrollModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ScrollModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ScrollModuleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ScrollModule *ScrollModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ScrollModuleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ScrollModule.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleOwnershipTransferredIterator{contract: _ScrollModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ScrollModule *ScrollModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ScrollModuleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ScrollModule.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ScrollModuleOwnershipTransferred)
				if err := _ScrollModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_ScrollModule *ScrollModuleFilterer) ParseOwnershipTransferred(log types.Log) (*ScrollModuleOwnershipTransferred, error) {
	event := new(ScrollModuleOwnershipTransferred)
	if err := _ScrollModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_ScrollModule *ScrollModule) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ScrollModule.abi.Events["L1FeeCoefficientSet"].ID:
		return _ScrollModule.ParseL1FeeCoefficientSet(log)
	case _ScrollModule.abi.Events["OwnershipTransferRequested"].ID:
		return _ScrollModule.ParseOwnershipTransferRequested(log)
	case _ScrollModule.abi.Events["OwnershipTransferred"].ID:
		return _ScrollModule.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ScrollModuleL1FeeCoefficientSet) Topic() common.Hash {
	return common.HexToHash("0x29ec9e31de0d3fe0208a7ccb792bbc26a854f123146110daa3a77219cb74a554")
}

func (ScrollModuleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ScrollModuleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_ScrollModule *ScrollModule) Address() common.Address {
	return _ScrollModule.address
}

type ScrollModuleInterface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SL1FeeCoefficient(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetL1FeeCalculation(opts *bind.TransactOpts, coefficient uint8) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterL1FeeCoefficientSet(opts *bind.FilterOpts) (*ScrollModuleL1FeeCoefficientSetIterator, error)

	WatchL1FeeCoefficientSet(opts *bind.WatchOpts, sink chan<- *ScrollModuleL1FeeCoefficientSet) (event.Subscription, error)

	ParseL1FeeCoefficientSet(log types.Log) (*ScrollModuleL1FeeCoefficientSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ScrollModuleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ScrollModuleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ScrollModuleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ScrollModuleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ScrollModuleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ScrollModuleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
