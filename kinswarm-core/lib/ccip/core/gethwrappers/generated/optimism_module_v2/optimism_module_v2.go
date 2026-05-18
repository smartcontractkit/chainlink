// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_module_v2

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

var OptimismModuleV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"InvalidL1FeeCoefficient\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"L1FeeCoefficientSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getL1FeeCoefficient\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"coefficient\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"setL1FeeCalculation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805460ff60a01b1916601960a21b17905534801561002357600080fd5b50338060008161007a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100aa576100aa816100b2565b50505061015b565b336001600160a01b0382160361010a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610071565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61073f8061016a6000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638da5cb5b11610076578063de9ee35e1161005b578063de9ee35e1461016a578063f22559a014610180578063f2fde38b1461019357600080fd5b80638da5cb5b1461011f578063d10a944e1461014757600080fd5b80637810d12a116100a75780637810d12a146100ef57806379ba50971461010257806385df51fd1461010c57600080fd5b806312544140146100c357806357e871e7146100e9575b600080fd5b6100d66100d136600461060c565b6101a6565b6040519081526020015b60405180910390f35b436100d6565b6100d66100fd36600461060c565b6101b7565b61010a6101f6565b005b6100d661011a36600461060c565b6102f8565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e0565b60015474010000000000000000000000000000000000000000900460ff166100d6565b60408051616d60815260006020820152016100e0565b61010a61018e366004610625565b610325565b61010a6101a136600461064f565b6103f0565b60006101b182610404565b92915050565b600060646101c483610404565b6001546101ec919074010000000000000000000000000000000000000000900460ff166106b4565b6101b191906106cb565b60015473ffffffffffffffffffffffffffffffffffffffff16331461027c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000438210158061031357506101006103118343610706565b115b1561032057506000919050565b504090565b61032d610494565b60648160ff161115610370576040517f1a8a06a000000000000000000000000000000000000000000000000000000000815260ff82166004820152602401610273565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000060ff8416908102919091179091556040519081527f29ec9e31de0d3fe0208a7ccb792bbc26a854f123146110daa3a77219cb74a5549060200160405180910390a150565b6103f8610494565b61040181610517565b50565b6040517ff1c7a58b0000000000000000000000000000000000000000000000000000000081526004810182905260009073420000000000000000000000000000000000000f9063f1c7a58b90602401602060405180830381865afa158015610470573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101b19190610719565b60005473ffffffffffffffffffffffffffffffffffffffff163314610515576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610273565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610596576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610273565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561061e57600080fd5b5035919050565b60006020828403121561063757600080fd5b813560ff8116811461064857600080fd5b9392505050565b60006020828403121561066157600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461064857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176101b1576101b1610685565b600082610701577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b818103818111156101b1576101b1610685565b60006020828403121561072b57600080fd5b505191905056fea164736f6c6343000813000a",
}

var OptimismModuleV2ABI = OptimismModuleV2MetaData.ABI

var OptimismModuleV2Bin = OptimismModuleV2MetaData.Bin

func DeployOptimismModuleV2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OptimismModuleV2, error) {
	parsed, err := OptimismModuleV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OptimismModuleV2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OptimismModuleV2{address: address, abi: *parsed, OptimismModuleV2Caller: OptimismModuleV2Caller{contract: contract}, OptimismModuleV2Transactor: OptimismModuleV2Transactor{contract: contract}, OptimismModuleV2Filterer: OptimismModuleV2Filterer{contract: contract}}, nil
}

type OptimismModuleV2 struct {
	address common.Address
	abi     abi.ABI
	OptimismModuleV2Caller
	OptimismModuleV2Transactor
	OptimismModuleV2Filterer
}

type OptimismModuleV2Caller struct {
	contract *bind.BoundContract
}

type OptimismModuleV2Transactor struct {
	contract *bind.BoundContract
}

type OptimismModuleV2Filterer struct {
	contract *bind.BoundContract
}

type OptimismModuleV2Session struct {
	Contract     *OptimismModuleV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismModuleV2CallerSession struct {
	Contract *OptimismModuleV2Caller
	CallOpts bind.CallOpts
}

type OptimismModuleV2TransactorSession struct {
	Contract     *OptimismModuleV2Transactor
	TransactOpts bind.TransactOpts
}

type OptimismModuleV2Raw struct {
	Contract *OptimismModuleV2
}

type OptimismModuleV2CallerRaw struct {
	Contract *OptimismModuleV2Caller
}

type OptimismModuleV2TransactorRaw struct {
	Contract *OptimismModuleV2Transactor
}

func NewOptimismModuleV2(address common.Address, backend bind.ContractBackend) (*OptimismModuleV2, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismModuleV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismModuleV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2{address: address, abi: abi, OptimismModuleV2Caller: OptimismModuleV2Caller{contract: contract}, OptimismModuleV2Transactor: OptimismModuleV2Transactor{contract: contract}, OptimismModuleV2Filterer: OptimismModuleV2Filterer{contract: contract}}, nil
}

func NewOptimismModuleV2Caller(address common.Address, caller bind.ContractCaller) (*OptimismModuleV2Caller, error) {
	contract, err := bindOptimismModuleV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2Caller{contract: contract}, nil
}

func NewOptimismModuleV2Transactor(address common.Address, transactor bind.ContractTransactor) (*OptimismModuleV2Transactor, error) {
	contract, err := bindOptimismModuleV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2Transactor{contract: contract}, nil
}

func NewOptimismModuleV2Filterer(address common.Address, filterer bind.ContractFilterer) (*OptimismModuleV2Filterer, error) {
	contract, err := bindOptimismModuleV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2Filterer{contract: contract}, nil
}

func bindOptimismModuleV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismModuleV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismModuleV2 *OptimismModuleV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismModuleV2.Contract.OptimismModuleV2Caller.contract.Call(opts, result, method, params...)
}

func (_OptimismModuleV2 *OptimismModuleV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.OptimismModuleV2Transactor.contract.Transfer(opts)
}

func (_OptimismModuleV2 *OptimismModuleV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.OptimismModuleV2Transactor.contract.Transact(opts, method, params...)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismModuleV2.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismModuleV2 *OptimismModuleV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.contract.Transfer(opts)
}

func (_OptimismModuleV2 *OptimismModuleV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) BlockHash(n *big.Int) ([32]byte, error) {
	return _OptimismModuleV2.Contract.BlockHash(&_OptimismModuleV2.CallOpts, n)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _OptimismModuleV2.Contract.BlockHash(&_OptimismModuleV2.CallOpts, n)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) BlockNumber() (*big.Int, error) {
	return _OptimismModuleV2.Contract.BlockNumber(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) BlockNumber() (*big.Int, error) {
	return _OptimismModuleV2.Contract.BlockNumber(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) GetCurrentL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "getCurrentL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) GetCurrentL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModuleV2.Contract.GetCurrentL1Fee(&_OptimismModuleV2.CallOpts, dataSize)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) GetCurrentL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModuleV2.Contract.GetCurrentL1Fee(&_OptimismModuleV2.CallOpts, dataSize)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) GetGasOverhead() (GetGasOverhead,

	error) {
	return _OptimismModuleV2.Contract.GetGasOverhead(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _OptimismModuleV2.Contract.GetGasOverhead(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) GetL1FeeCoefficient(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "getL1FeeCoefficient")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) GetL1FeeCoefficient() (*big.Int, error) {
	return _OptimismModuleV2.Contract.GetL1FeeCoefficient(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) GetL1FeeCoefficient() (*big.Int, error) {
	return _OptimismModuleV2.Contract.GetL1FeeCoefficient(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModuleV2.Contract.GetMaxL1Fee(&_OptimismModuleV2.CallOpts, dataSize)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModuleV2.Contract.GetMaxL1Fee(&_OptimismModuleV2.CallOpts, dataSize)
}

func (_OptimismModuleV2 *OptimismModuleV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismModuleV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismModuleV2 *OptimismModuleV2Session) Owner() (common.Address, error) {
	return _OptimismModuleV2.Contract.Owner(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2CallerSession) Owner() (common.Address, error) {
	return _OptimismModuleV2.Contract.Owner(&_OptimismModuleV2.CallOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismModuleV2.contract.Transact(opts, "acceptOwnership")
}

func (_OptimismModuleV2 *OptimismModuleV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.AcceptOwnership(&_OptimismModuleV2.TransactOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.AcceptOwnership(&_OptimismModuleV2.TransactOpts)
}

func (_OptimismModuleV2 *OptimismModuleV2Transactor) SetL1FeeCalculation(opts *bind.TransactOpts, coefficient uint8) (*types.Transaction, error) {
	return _OptimismModuleV2.contract.Transact(opts, "setL1FeeCalculation", coefficient)
}

func (_OptimismModuleV2 *OptimismModuleV2Session) SetL1FeeCalculation(coefficient uint8) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.SetL1FeeCalculation(&_OptimismModuleV2.TransactOpts, coefficient)
}

func (_OptimismModuleV2 *OptimismModuleV2TransactorSession) SetL1FeeCalculation(coefficient uint8) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.SetL1FeeCalculation(&_OptimismModuleV2.TransactOpts, coefficient)
}

func (_OptimismModuleV2 *OptimismModuleV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OptimismModuleV2.contract.Transact(opts, "transferOwnership", to)
}

func (_OptimismModuleV2 *OptimismModuleV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.TransferOwnership(&_OptimismModuleV2.TransactOpts, to)
}

func (_OptimismModuleV2 *OptimismModuleV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OptimismModuleV2.Contract.TransferOwnership(&_OptimismModuleV2.TransactOpts, to)
}

type OptimismModuleV2L1FeeCoefficientSetIterator struct {
	Event *OptimismModuleV2L1FeeCoefficientSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismModuleV2L1FeeCoefficientSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismModuleV2L1FeeCoefficientSet)
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
		it.Event = new(OptimismModuleV2L1FeeCoefficientSet)
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

func (it *OptimismModuleV2L1FeeCoefficientSetIterator) Error() error {
	return it.fail
}

func (it *OptimismModuleV2L1FeeCoefficientSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismModuleV2L1FeeCoefficientSet struct {
	Coefficient uint8
	Raw         types.Log
}

func (_OptimismModuleV2 *OptimismModuleV2Filterer) FilterL1FeeCoefficientSet(opts *bind.FilterOpts) (*OptimismModuleV2L1FeeCoefficientSetIterator, error) {

	logs, sub, err := _OptimismModuleV2.contract.FilterLogs(opts, "L1FeeCoefficientSet")
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2L1FeeCoefficientSetIterator{contract: _OptimismModuleV2.contract, event: "L1FeeCoefficientSet", logs: logs, sub: sub}, nil
}

func (_OptimismModuleV2 *OptimismModuleV2Filterer) WatchL1FeeCoefficientSet(opts *bind.WatchOpts, sink chan<- *OptimismModuleV2L1FeeCoefficientSet) (event.Subscription, error) {

	logs, sub, err := _OptimismModuleV2.contract.WatchLogs(opts, "L1FeeCoefficientSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismModuleV2L1FeeCoefficientSet)
				if err := _OptimismModuleV2.contract.UnpackLog(event, "L1FeeCoefficientSet", log); err != nil {
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

func (_OptimismModuleV2 *OptimismModuleV2Filterer) ParseL1FeeCoefficientSet(log types.Log) (*OptimismModuleV2L1FeeCoefficientSet, error) {
	event := new(OptimismModuleV2L1FeeCoefficientSet)
	if err := _OptimismModuleV2.contract.UnpackLog(event, "L1FeeCoefficientSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OptimismModuleV2OwnershipTransferRequestedIterator struct {
	Event *OptimismModuleV2OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismModuleV2OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismModuleV2OwnershipTransferRequested)
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
		it.Event = new(OptimismModuleV2OwnershipTransferRequested)
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

func (it *OptimismModuleV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OptimismModuleV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismModuleV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OptimismModuleV2 *OptimismModuleV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismModuleV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismModuleV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2OwnershipTransferRequestedIterator{contract: _OptimismModuleV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OptimismModuleV2 *OptimismModuleV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OptimismModuleV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismModuleV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismModuleV2OwnershipTransferRequested)
				if err := _OptimismModuleV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OptimismModuleV2 *OptimismModuleV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*OptimismModuleV2OwnershipTransferRequested, error) {
	event := new(OptimismModuleV2OwnershipTransferRequested)
	if err := _OptimismModuleV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OptimismModuleV2OwnershipTransferredIterator struct {
	Event *OptimismModuleV2OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismModuleV2OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismModuleV2OwnershipTransferred)
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
		it.Event = new(OptimismModuleV2OwnershipTransferred)
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

func (it *OptimismModuleV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OptimismModuleV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismModuleV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OptimismModuleV2 *OptimismModuleV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismModuleV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismModuleV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleV2OwnershipTransferredIterator{contract: _OptimismModuleV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OptimismModuleV2 *OptimismModuleV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OptimismModuleV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismModuleV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismModuleV2OwnershipTransferred)
				if err := _OptimismModuleV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OptimismModuleV2 *OptimismModuleV2Filterer) ParseOwnershipTransferred(log types.Log) (*OptimismModuleV2OwnershipTransferred, error) {
	event := new(OptimismModuleV2OwnershipTransferred)
	if err := _OptimismModuleV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_OptimismModuleV2 *OptimismModuleV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OptimismModuleV2.abi.Events["L1FeeCoefficientSet"].ID:
		return _OptimismModuleV2.ParseL1FeeCoefficientSet(log)
	case _OptimismModuleV2.abi.Events["OwnershipTransferRequested"].ID:
		return _OptimismModuleV2.ParseOwnershipTransferRequested(log)
	case _OptimismModuleV2.abi.Events["OwnershipTransferred"].ID:
		return _OptimismModuleV2.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OptimismModuleV2L1FeeCoefficientSet) Topic() common.Hash {
	return common.HexToHash("0x29ec9e31de0d3fe0208a7ccb792bbc26a854f123146110daa3a77219cb74a554")
}

func (OptimismModuleV2OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OptimismModuleV2OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OptimismModuleV2 *OptimismModuleV2) Address() common.Address {
	return _OptimismModuleV2.address
}

type OptimismModuleV2Interface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetL1FeeCoefficient(opts *bind.CallOpts) (*big.Int, error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetL1FeeCalculation(opts *bind.TransactOpts, coefficient uint8) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterL1FeeCoefficientSet(opts *bind.FilterOpts) (*OptimismModuleV2L1FeeCoefficientSetIterator, error)

	WatchL1FeeCoefficientSet(opts *bind.WatchOpts, sink chan<- *OptimismModuleV2L1FeeCoefficientSet) (event.Subscription, error)

	ParseL1FeeCoefficientSet(log types.Log) (*OptimismModuleV2L1FeeCoefficientSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismModuleV2OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OptimismModuleV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OptimismModuleV2OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismModuleV2OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OptimismModuleV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OptimismModuleV2OwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
