// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_gas_module

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

var OptimismGasModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"l1_gas_calculation_mode\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidMode\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"getTxL1GasFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_l1_gas_calculation_mode\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"l1_gas_calculation_mode\",\"type\":\"uint8\"}],\"name\":\"setL1GasCalculationMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610f6c380380610f6c83398101604081905261002f91610185565b33806000816100855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b5576100b5816100dc565b50506001805460ff909316600160a01b0260ff60a01b1990931692909217909155506101af565b336001600160a01b038216036101345760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561019757600080fd5b815160ff811681146101a857600080fd5b9392505050565b610dae806101be6000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c80638dd6fd2e116100505780638dd6fd2e146100cf578063bd72c23b14610106578063f2fde38b1461011957600080fd5b806302c693371461007757806379ba50971461009d5780638da5cb5b146100a7575b600080fd5b61008a610085366004610960565b61012c565b6040519081526020015b60405180910390f35b6100a5610302565b005b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610094565b6001546100f49074010000000000000000000000000000000000000000900460ff1681565b60405160ff9091168152602001610094565b6100a5610114366004610a2f565b610404565b6100a5610127366004610a59565b610496565b60015460009074010000000000000000000000000000000000000000900460ff1681036102205773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e83604051806080016040528060478152602001610d5b604791396040516020016101ae929190610ab3565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016101d99190610ae2565b602060405180830381865afa1580156101f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021a9190610b33565b92915050565b6001805474010000000000000000000000000000000000000000900460ff16900361024f5761021a82516104aa565b60015474010000000000000000000000000000000000000000900460ff166002036102d05773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663f1c7a58b604784516102b29190610b7b565b6040518263ffffffff1660e01b81526004016101d991815260200190565b6040517fa0042b1700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff163314610388576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61040c6107b9565b60028160ff16111561044a576040517fa0042b1700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001805460ff90921674010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff909216919091179055565b61049e6107b9565b6104a78161083c565b50565b60008073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663519b4bd36040518163ffffffff1660e01b8152600401602060405180830381865afa15801561050c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105309190610b33565b73420000000000000000000000000000000000001573ffffffffffffffffffffffffffffffffffffffff1663c59859186040518163ffffffff1660e01b8152600401602060405180830381865afa15801561058f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105b39190610b8e565b6105be906010610bb4565b63ffffffff166105ce9190610bdc565b9050600073420000000000000000000000000000000000001573ffffffffffffffffffffffffffffffffffffffff1663f82061406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610631573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106559190610b33565b73420000000000000000000000000000000000001573ffffffffffffffffffffffffffffffffffffffff166368d5dca66040518163ffffffff1660e01b8152600401602060405180830381865afa1580156106b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106d89190610b8e565b63ffffffff166106e89190610bdc565b905060006106f68284610b7b565b61070260476044610b7b565b61070c9087610b7b565b6107169190610bdc565b905073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015610777573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061079b9190610b33565b6107a690600a610d13565b6107b09082610d1f565b95945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461083a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161037f565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036108bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161037f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60006020828403121561097257600080fd5b813567ffffffffffffffff8082111561098a57600080fd5b818401915084601f83011261099e57600080fd5b8135818111156109b0576109b0610931565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156109f6576109f6610931565b81604052828152876020848701011115610a0f57600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208284031215610a4157600080fd5b813560ff81168114610a5257600080fd5b9392505050565b600060208284031215610a6b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610a5257600080fd5b60005b83811015610aaa578181015183820152602001610a92565b50506000910152565b60008351610ac5818460208801610a8f565b835190830190610ad9818360208801610a8f565b01949350505050565b6020815260008251806020840152610b01816040850160208701610a8f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b600060208284031215610b4557600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561021a5761021a610b4c565b600060208284031215610ba057600080fd5b815163ffffffff81168114610a5257600080fd5b63ffffffff818116838216028082169190828114610bd457610bd4610b4c565b505092915050565b808202811582820484141761021a5761021a610b4c565b600181815b80851115610c4c57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115610c3257610c32610b4c565b80851615610c3f57918102915b93841c9390800290610bf8565b509250929050565b600082610c635750600161021a565b81610c705750600061021a565b8160018114610c865760028114610c9057610cac565b600191505061021a565b60ff841115610ca157610ca1610b4c565b50506001821b61021a565b5060208310610133831016604e8410600b8410161715610ccf575081810a61021a565b610cd98383610bf3565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115610d0b57610d0b610b4c565b029392505050565b6000610a528383610c54565b600082610d55577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
}

var OptimismGasModuleABI = OptimismGasModuleMetaData.ABI

var OptimismGasModuleBin = OptimismGasModuleMetaData.Bin

func DeployOptimismGasModule(auth *bind.TransactOpts, backend bind.ContractBackend, l1_gas_calculation_mode uint8) (common.Address, *types.Transaction, *OptimismGasModule, error) {
	parsed, err := OptimismGasModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OptimismGasModuleBin), backend, l1_gas_calculation_mode)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OptimismGasModule{address: address, abi: *parsed, OptimismGasModuleCaller: OptimismGasModuleCaller{contract: contract}, OptimismGasModuleTransactor: OptimismGasModuleTransactor{contract: contract}, OptimismGasModuleFilterer: OptimismGasModuleFilterer{contract: contract}}, nil
}

type OptimismGasModule struct {
	address common.Address
	abi     abi.ABI
	OptimismGasModuleCaller
	OptimismGasModuleTransactor
	OptimismGasModuleFilterer
}

type OptimismGasModuleCaller struct {
	contract *bind.BoundContract
}

type OptimismGasModuleTransactor struct {
	contract *bind.BoundContract
}

type OptimismGasModuleFilterer struct {
	contract *bind.BoundContract
}

type OptimismGasModuleSession struct {
	Contract     *OptimismGasModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismGasModuleCallerSession struct {
	Contract *OptimismGasModuleCaller
	CallOpts bind.CallOpts
}

type OptimismGasModuleTransactorSession struct {
	Contract     *OptimismGasModuleTransactor
	TransactOpts bind.TransactOpts
}

type OptimismGasModuleRaw struct {
	Contract *OptimismGasModule
}

type OptimismGasModuleCallerRaw struct {
	Contract *OptimismGasModuleCaller
}

type OptimismGasModuleTransactorRaw struct {
	Contract *OptimismGasModuleTransactor
}

func NewOptimismGasModule(address common.Address, backend bind.ContractBackend) (*OptimismGasModule, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismGasModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismGasModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismGasModule{address: address, abi: abi, OptimismGasModuleCaller: OptimismGasModuleCaller{contract: contract}, OptimismGasModuleTransactor: OptimismGasModuleTransactor{contract: contract}, OptimismGasModuleFilterer: OptimismGasModuleFilterer{contract: contract}}, nil
}

func NewOptimismGasModuleCaller(address common.Address, caller bind.ContractCaller) (*OptimismGasModuleCaller, error) {
	contract, err := bindOptimismGasModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismGasModuleCaller{contract: contract}, nil
}

func NewOptimismGasModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismGasModuleTransactor, error) {
	contract, err := bindOptimismGasModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismGasModuleTransactor{contract: contract}, nil
}

func NewOptimismGasModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismGasModuleFilterer, error) {
	contract, err := bindOptimismGasModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismGasModuleFilterer{contract: contract}, nil
}

func bindOptimismGasModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismGasModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismGasModule *OptimismGasModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismGasModule.Contract.OptimismGasModuleCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismGasModule *OptimismGasModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.OptimismGasModuleTransactor.contract.Transfer(opts)
}

func (_OptimismGasModule *OptimismGasModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.OptimismGasModuleTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismGasModule *OptimismGasModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismGasModule.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismGasModule *OptimismGasModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.contract.Transfer(opts)
}

func (_OptimismGasModule *OptimismGasModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismGasModule *OptimismGasModuleCaller) GetTxL1GasFees(opts *bind.CallOpts, data []byte) (*big.Int, error) {
	var out []interface{}
	err := _OptimismGasModule.contract.Call(opts, &out, "getTxL1GasFees", data)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismGasModule *OptimismGasModuleSession) GetTxL1GasFees(data []byte) (*big.Int, error) {
	return _OptimismGasModule.Contract.GetTxL1GasFees(&_OptimismGasModule.CallOpts, data)
}

func (_OptimismGasModule *OptimismGasModuleCallerSession) GetTxL1GasFees(data []byte) (*big.Int, error) {
	return _OptimismGasModule.Contract.GetTxL1GasFees(&_OptimismGasModule.CallOpts, data)
}

func (_OptimismGasModule *OptimismGasModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismGasModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismGasModule *OptimismGasModuleSession) Owner() (common.Address, error) {
	return _OptimismGasModule.Contract.Owner(&_OptimismGasModule.CallOpts)
}

func (_OptimismGasModule *OptimismGasModuleCallerSession) Owner() (common.Address, error) {
	return _OptimismGasModule.Contract.Owner(&_OptimismGasModule.CallOpts)
}

func (_OptimismGasModule *OptimismGasModuleCaller) SL1GasCalculationMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _OptimismGasModule.contract.Call(opts, &out, "s_l1_gas_calculation_mode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_OptimismGasModule *OptimismGasModuleSession) SL1GasCalculationMode() (uint8, error) {
	return _OptimismGasModule.Contract.SL1GasCalculationMode(&_OptimismGasModule.CallOpts)
}

func (_OptimismGasModule *OptimismGasModuleCallerSession) SL1GasCalculationMode() (uint8, error) {
	return _OptimismGasModule.Contract.SL1GasCalculationMode(&_OptimismGasModule.CallOpts)
}

func (_OptimismGasModule *OptimismGasModuleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismGasModule.contract.Transact(opts, "acceptOwnership")
}

func (_OptimismGasModule *OptimismGasModuleSession) AcceptOwnership() (*types.Transaction, error) {
	return _OptimismGasModule.Contract.AcceptOwnership(&_OptimismGasModule.TransactOpts)
}

func (_OptimismGasModule *OptimismGasModuleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OptimismGasModule.Contract.AcceptOwnership(&_OptimismGasModule.TransactOpts)
}

func (_OptimismGasModule *OptimismGasModuleTransactor) SetL1GasCalculationMode(opts *bind.TransactOpts, l1_gas_calculation_mode uint8) (*types.Transaction, error) {
	return _OptimismGasModule.contract.Transact(opts, "setL1GasCalculationMode", l1_gas_calculation_mode)
}

func (_OptimismGasModule *OptimismGasModuleSession) SetL1GasCalculationMode(l1_gas_calculation_mode uint8) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.SetL1GasCalculationMode(&_OptimismGasModule.TransactOpts, l1_gas_calculation_mode)
}

func (_OptimismGasModule *OptimismGasModuleTransactorSession) SetL1GasCalculationMode(l1_gas_calculation_mode uint8) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.SetL1GasCalculationMode(&_OptimismGasModule.TransactOpts, l1_gas_calculation_mode)
}

func (_OptimismGasModule *OptimismGasModuleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OptimismGasModule.contract.Transact(opts, "transferOwnership", to)
}

func (_OptimismGasModule *OptimismGasModuleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.TransferOwnership(&_OptimismGasModule.TransactOpts, to)
}

func (_OptimismGasModule *OptimismGasModuleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OptimismGasModule.Contract.TransferOwnership(&_OptimismGasModule.TransactOpts, to)
}

type OptimismGasModuleOwnershipTransferRequestedIterator struct {
	Event *OptimismGasModuleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismGasModuleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismGasModuleOwnershipTransferRequested)
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
		it.Event = new(OptimismGasModuleOwnershipTransferRequested)
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

func (it *OptimismGasModuleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OptimismGasModuleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismGasModuleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OptimismGasModule *OptimismGasModuleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismGasModuleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismGasModule.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OptimismGasModuleOwnershipTransferRequestedIterator{contract: _OptimismGasModule.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OptimismGasModule *OptimismGasModuleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OptimismGasModuleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismGasModule.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismGasModuleOwnershipTransferRequested)
				if err := _OptimismGasModule.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OptimismGasModule *OptimismGasModuleFilterer) ParseOwnershipTransferRequested(log types.Log) (*OptimismGasModuleOwnershipTransferRequested, error) {
	event := new(OptimismGasModuleOwnershipTransferRequested)
	if err := _OptimismGasModule.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OptimismGasModuleOwnershipTransferredIterator struct {
	Event *OptimismGasModuleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismGasModuleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismGasModuleOwnershipTransferred)
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
		it.Event = new(OptimismGasModuleOwnershipTransferred)
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

func (it *OptimismGasModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OptimismGasModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismGasModuleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OptimismGasModule *OptimismGasModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismGasModuleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismGasModule.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OptimismGasModuleOwnershipTransferredIterator{contract: _OptimismGasModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OptimismGasModule *OptimismGasModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OptimismGasModuleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OptimismGasModule.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismGasModuleOwnershipTransferred)
				if err := _OptimismGasModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OptimismGasModule *OptimismGasModuleFilterer) ParseOwnershipTransferred(log types.Log) (*OptimismGasModuleOwnershipTransferred, error) {
	event := new(OptimismGasModuleOwnershipTransferred)
	if err := _OptimismGasModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OptimismGasModule *OptimismGasModule) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OptimismGasModule.abi.Events["OwnershipTransferRequested"].ID:
		return _OptimismGasModule.ParseOwnershipTransferRequested(log)
	case _OptimismGasModule.abi.Events["OwnershipTransferred"].ID:
		return _OptimismGasModule.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OptimismGasModuleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OptimismGasModuleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OptimismGasModule *OptimismGasModule) Address() common.Address {
	return _OptimismGasModule.address
}

type OptimismGasModuleInterface interface {
	GetTxL1GasFees(opts *bind.CallOpts, data []byte) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SL1GasCalculationMode(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetL1GasCalculationMode(opts *bind.TransactOpts, l1_gas_calculation_mode uint8) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismGasModuleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OptimismGasModuleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OptimismGasModuleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OptimismGasModuleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OptimismGasModuleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OptimismGasModuleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
