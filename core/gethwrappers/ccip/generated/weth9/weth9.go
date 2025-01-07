// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package weth9

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

var WETH9MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"guy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dst\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"guy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60806040523461011457610014600054610119565b601f81116100cb575b507f577261707065642045746865720000000000000000000000000000000000001a60005560015461004e90610119565b601f8111610081575b6008630ae8aa8960e31b016001556002805460ff19166012179055604051610a5a90816101548239f35b6001600052601f0160051c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101905b8181106100bf5750610057565b600081556001016100b2565b60008052601f0160051c7f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563908101905b818110610108575061001d565b600081556001016100fb565b600080fd5b90600182811c92168015610149575b602083101461013357565b634e487b7160e01b600052602260045260246000fd5b91607f169161012856fe60806040526004361015610062575b361561001957600080fd5b33600052600360205260406000206100323482546108a0565b90556040513481527fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c60203392a2005b60003560e01c806306fdde0314610696578063095ea7b3146105ef57806318160ddd146105b557806323b872dd146105685780632e1a7d4d146104aa578063313ce5671461046b57806370a082311461040657806395d89b4114610208578063a9059cbb146101b8578063d0e30db0146101755763dd62ed3e0361000e57346101705760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705761011761081e565b73ffffffffffffffffffffffffffffffffffffffff610134610841565b9116600052600460205273ffffffffffffffffffffffffffffffffffffffff604060002091166000526020526020604060002054604051908152f35b600080fd5b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705733600052600360205260406000206100323482546108a0565b346101705760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705760206101fe6101f461081e565b60243590336108ad565b6040519015158152f35b346101705760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610170576000604051908091600154928360011c600185169485156103fc575b6020821086146103cf57839495828552908160001461036f57506001146102f6575b5003601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210176102c9576102c59250604052604051918291826107b6565b0390f35b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b600185528491507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061035357505081016020017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0610275565b602091935080600191548385880101520191019091839261031f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660208581019190915291151560051b840190910191507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09050610275565b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526022600452fd5b90607f1690610253565b346101705760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705773ffffffffffffffffffffffffffffffffffffffff61045261081e565b1660005260036020526020604060002054604051908152f35b346101705760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017057602060ff60025416604051908152f35b346101705760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017057600435336000526003602052806040600020541061017057336000526003602052604060002061050a828254610864565b9055806000811561055f575b600080809381933390f115610553576040519081527f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b6560203392a2005b6040513d6000823e3d90fd5b506108fc610516565b346101705760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705760206101fe6105a461081e565b6105ac610841565b604435916108ad565b346101705760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017057602047604051908152f35b346101705760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705761062661081e565b73ffffffffffffffffffffffffffffffffffffffff6024359133600052600460205260406000208282166000526020528260406000205560405192835216907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560203392a3602060405160018152f35b346101705760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705760006040519080918154928360011c600185169485156107ac575b6020821086146103cf57839495828552908160001461036f5750600114610751575003601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210176102c9576102c59250604052604051918291826107b6565b848052602085208592505b81831061079057505081016020017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0610275565b602091935080600191548385880101520191019091839261075c565b90607f16906106e0565b9190916020815282519283602083015260005b8481106108085750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006040809697860101520116010190565b80602080928401015160408286010152016107c9565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361017057565b6024359073ffffffffffffffffffffffffffffffffffffffff8216820361017057565b9190820391821161087157565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9190820180921161087157565b73ffffffffffffffffffffffffffffffffffffffff16908160005260036020528260406000205410610170573382141580610a03575b610963575b602073ffffffffffffffffffffffffffffffffffffffff7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9284600052600383526040600020610939878254610864565b90551693846000526003825260406000206109558282546108a0565b9055604051908152a3600190565b816000526004602052604060002073ffffffffffffffffffffffffffffffffffffffff3316600052602052826040600020541061017057602073ffffffffffffffffffffffffffffffffffffffff7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9284600052600483526040600020823316600052835260406000206109f8878254610864565b9055925050506108e8565b50816000526004602052604060002073ffffffffffffffffffffffffffffffffffffffff33166000526020526fffffffffffffffffffffffffffffffff60406000205414156108e356fea164736f6c634300081a000a",
}

var WETH9ABI = WETH9MetaData.ABI

var WETH9Bin = WETH9MetaData.Bin

func DeployWETH9(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WETH9, error) {
	parsed, err := WETH9MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WETH9Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WETH9{address: address, abi: *parsed, WETH9Caller: WETH9Caller{contract: contract}, WETH9Transactor: WETH9Transactor{contract: contract}, WETH9Filterer: WETH9Filterer{contract: contract}}, nil
}

type WETH9 struct {
	address common.Address
	abi     abi.ABI
	WETH9Caller
	WETH9Transactor
	WETH9Filterer
}

type WETH9Caller struct {
	contract *bind.BoundContract
}

type WETH9Transactor struct {
	contract *bind.BoundContract
}

type WETH9Filterer struct {
	contract *bind.BoundContract
}

type WETH9Session struct {
	Contract     *WETH9
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type WETH9CallerSession struct {
	Contract *WETH9Caller
	CallOpts bind.CallOpts
}

type WETH9TransactorSession struct {
	Contract     *WETH9Transactor
	TransactOpts bind.TransactOpts
}

type WETH9Raw struct {
	Contract *WETH9
}

type WETH9CallerRaw struct {
	Contract *WETH9Caller
}

type WETH9TransactorRaw struct {
	Contract *WETH9Transactor
}

func NewWETH9(address common.Address, backend bind.ContractBackend) (*WETH9, error) {
	abi, err := abi.JSON(strings.NewReader(WETH9ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindWETH9(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WETH9{address: address, abi: abi, WETH9Caller: WETH9Caller{contract: contract}, WETH9Transactor: WETH9Transactor{contract: contract}, WETH9Filterer: WETH9Filterer{contract: contract}}, nil
}

func NewWETH9Caller(address common.Address, caller bind.ContractCaller) (*WETH9Caller, error) {
	contract, err := bindWETH9(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WETH9Caller{contract: contract}, nil
}

func NewWETH9Transactor(address common.Address, transactor bind.ContractTransactor) (*WETH9Transactor, error) {
	contract, err := bindWETH9(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WETH9Transactor{contract: contract}, nil
}

func NewWETH9Filterer(address common.Address, filterer bind.ContractFilterer) (*WETH9Filterer, error) {
	contract, err := bindWETH9(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WETH9Filterer{contract: contract}, nil
}

func bindWETH9(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WETH9MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_WETH9 *WETH9Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WETH9.Contract.WETH9Caller.contract.Call(opts, result, method, params...)
}

func (_WETH9 *WETH9Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.Contract.WETH9Transactor.contract.Transfer(opts)
}

func (_WETH9 *WETH9Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WETH9.Contract.WETH9Transactor.contract.Transact(opts, method, params...)
}

func (_WETH9 *WETH9CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WETH9.Contract.contract.Call(opts, result, method, params...)
}

func (_WETH9 *WETH9TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.Contract.contract.Transfer(opts)
}

func (_WETH9 *WETH9TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WETH9.Contract.contract.Transact(opts, method, params...)
}

func (_WETH9 *WETH9Caller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WETH9 *WETH9Session) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WETH9.Contract.Allowance(&_WETH9.CallOpts, arg0, arg1)
}

func (_WETH9 *WETH9CallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WETH9.Contract.Allowance(&_WETH9.CallOpts, arg0, arg1)
}

func (_WETH9 *WETH9Caller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WETH9 *WETH9Session) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WETH9.Contract.BalanceOf(&_WETH9.CallOpts, arg0)
}

func (_WETH9 *WETH9CallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WETH9.Contract.BalanceOf(&_WETH9.CallOpts, arg0)
}

func (_WETH9 *WETH9Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_WETH9 *WETH9Session) Decimals() (uint8, error) {
	return _WETH9.Contract.Decimals(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) Decimals() (uint8, error) {
	return _WETH9.Contract.Decimals(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WETH9 *WETH9Session) Name() (string, error) {
	return _WETH9.Contract.Name(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) Name() (string, error) {
	return _WETH9.Contract.Name(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WETH9 *WETH9Session) Symbol() (string, error) {
	return _WETH9.Contract.Symbol(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) Symbol() (string, error) {
	return _WETH9.Contract.Symbol(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WETH9 *WETH9Session) TotalSupply() (*big.Int, error) {
	return _WETH9.Contract.TotalSupply(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) TotalSupply() (*big.Int, error) {
	return _WETH9.Contract.TotalSupply(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Transactor) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "approve", guy, wad)
}

func (_WETH9 *WETH9Session) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Approve(&_WETH9.TransactOpts, guy, wad)
}

func (_WETH9 *WETH9TransactorSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Approve(&_WETH9.TransactOpts, guy, wad)
}

func (_WETH9 *WETH9Transactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "deposit")
}

func (_WETH9 *WETH9Session) Deposit() (*types.Transaction, error) {
	return _WETH9.Contract.Deposit(&_WETH9.TransactOpts)
}

func (_WETH9 *WETH9TransactorSession) Deposit() (*types.Transaction, error) {
	return _WETH9.Contract.Deposit(&_WETH9.TransactOpts)
}

func (_WETH9 *WETH9Transactor) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "transfer", dst, wad)
}

func (_WETH9 *WETH9Session) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Transfer(&_WETH9.TransactOpts, dst, wad)
}

func (_WETH9 *WETH9TransactorSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Transfer(&_WETH9.TransactOpts, dst, wad)
}

func (_WETH9 *WETH9Transactor) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "transferFrom", src, dst, wad)
}

func (_WETH9 *WETH9Session) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.TransferFrom(&_WETH9.TransactOpts, src, dst, wad)
}

func (_WETH9 *WETH9TransactorSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.TransferFrom(&_WETH9.TransactOpts, src, dst, wad)
}

func (_WETH9 *WETH9Transactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "withdraw", wad)
}

func (_WETH9 *WETH9Session) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Withdraw(&_WETH9.TransactOpts, wad)
}

func (_WETH9 *WETH9TransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Withdraw(&_WETH9.TransactOpts, wad)
}

func (_WETH9 *WETH9Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.contract.RawTransact(opts, nil)
}

func (_WETH9 *WETH9Session) Receive() (*types.Transaction, error) {
	return _WETH9.Contract.Receive(&_WETH9.TransactOpts)
}

func (_WETH9 *WETH9TransactorSession) Receive() (*types.Transaction, error) {
	return _WETH9.Contract.Receive(&_WETH9.TransactOpts)
}

type WETH9ApprovalIterator struct {
	Event *WETH9Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Approval)
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
		it.Event = new(WETH9Approval)
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

func (it *WETH9ApprovalIterator) Error() error {
	return it.fail
}

func (it *WETH9ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Approval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WETH9ApprovalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return &WETH9ApprovalIterator{contract: _WETH9.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WETH9Approval, src []common.Address, guy []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Approval)
				if err := _WETH9.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseApproval(log types.Log) (*WETH9Approval, error) {
	event := new(WETH9Approval)
	if err := _WETH9.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WETH9DepositIterator struct {
	Event *WETH9Deposit

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9DepositIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Deposit)
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
		it.Event = new(WETH9Deposit)
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

func (it *WETH9DepositIterator) Error() error {
	return it.fail
}

func (it *WETH9DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Deposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WETH9DepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &WETH9DepositIterator{contract: _WETH9.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WETH9Deposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Deposit)
				if err := _WETH9.contract.UnpackLog(event, "Deposit", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseDeposit(log types.Log) (*WETH9Deposit, error) {
	event := new(WETH9Deposit)
	if err := _WETH9.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WETH9TransferIterator struct {
	Event *WETH9Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Transfer)
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
		it.Event = new(WETH9Transfer)
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

func (it *WETH9TransferIterator) Error() error {
	return it.fail
}

func (it *WETH9TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Transfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WETH9TransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &WETH9TransferIterator{contract: _WETH9.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WETH9Transfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Transfer)
				if err := _WETH9.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseTransfer(log types.Log) (*WETH9Transfer, error) {
	event := new(WETH9Transfer)
	if err := _WETH9.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WETH9WithdrawalIterator struct {
	Event *WETH9Withdrawal

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9WithdrawalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Withdrawal)
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
		it.Event = new(WETH9Withdrawal)
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

func (it *WETH9WithdrawalIterator) Error() error {
	return it.fail
}

func (it *WETH9WithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Withdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WETH9WithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &WETH9WithdrawalIterator{contract: _WETH9.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WETH9Withdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Withdrawal)
				if err := _WETH9.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseWithdrawal(log types.Log) (*WETH9Withdrawal, error) {
	event := new(WETH9Withdrawal)
	if err := _WETH9.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_WETH9 *WETH9) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _WETH9.abi.Events["Approval"].ID:
		return _WETH9.ParseApproval(log)
	case _WETH9.abi.Events["Deposit"].ID:
		return _WETH9.ParseDeposit(log)
	case _WETH9.abi.Events["Transfer"].ID:
		return _WETH9.ParseTransfer(log)
	case _WETH9.abi.Events["Withdrawal"].ID:
		return _WETH9.ParseWithdrawal(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (WETH9Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (WETH9Deposit) Topic() common.Hash {
	return common.HexToHash("0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c")
}

func (WETH9Transfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (WETH9Withdrawal) Topic() common.Hash {
	return common.HexToHash("0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65")
}

func (_WETH9 *WETH9) Address() common.Address {
	return _WETH9.address
}

type WETH9Interface interface {
	Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Name(opts *bind.CallOpts) (string, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error)

	Deposit(opts *bind.TransactOpts) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WETH9ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *WETH9Approval, src []common.Address, guy []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*WETH9Approval, error)

	FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WETH9DepositIterator, error)

	WatchDeposit(opts *bind.WatchOpts, sink chan<- *WETH9Deposit, dst []common.Address) (event.Subscription, error)

	ParseDeposit(log types.Log) (*WETH9Deposit, error)

	FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WETH9TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *WETH9Transfer, src []common.Address, dst []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*WETH9Transfer, error)

	FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WETH9WithdrawalIterator, error)

	WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WETH9Withdrawal, src []common.Address) (event.Subscription, error)

	ParseWithdrawal(log types.Log) (*WETH9Withdrawal, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
