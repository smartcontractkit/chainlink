// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gas_wrapper_mock

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

var KeeperRegistryCheckUpkeepGasUsageWrapperMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"emitOwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitOwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"measureCheckGas\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"s_mockGas\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_mockPayload\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_mockResult\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setMeasureCheckGasResult\",\"inputs\":[{\"name\":\"result\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"gas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b506106b9806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063846811931161005b57806384681193146100d5578063b019b4e8146100f2578063b023145014610105578063f7420bc21461011a57600080fd5b80632dae06f51461008257806356343496146100975780636bf49030146100b3575b600080fd5b610095610090366004610466565b61012d565b005b6100a060025481565b6040519081526020015b60405180910390f35b6100c66100c1366004610556565b610174565b6040516100aa939291906105e4565b6000546100e29060ff1681565b60405190151581526020016100aa565b610095610100366004610433565b610227565b61010d610285565b6040516100aa919061060f565b610095610128366004610433565b610313565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016841515179055815161016c906001906020850190610371565b506002555050565b6000606060008060009054906101000a900460ff16600160025481805461019a90610629565b80601f01602080910402602001604051908101604052809291908181526020018280546101c690610629565b80156102135780601f106101e857610100808354040283529160200191610213565b820191906000526020600020905b8154815290600101906020018083116101f657829003601f168201915b505050505091509250925092509250925092565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6001805461029290610629565b80601f01602080910402602001604051908101604052809291908181526020018280546102be90610629565b801561030b5780601f106102e05761010080835404028352916020019161030b565b820191906000526020600020905b8154815290600101906020018083116102ee57829003601f168201915b505050505081565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b82805461037d90610629565b90600052602060002090601f01602090048101928261039f57600085556103e5565b82601f106103b857805160ff19168380011785556103e5565b828001600101855582156103e5579182015b828111156103e55782518255916020019190600101906103ca565b506103f19291506103f5565b5090565b5b808211156103f157600081556001016103f6565b803573ffffffffffffffffffffffffffffffffffffffff8116811461042e57600080fd5b919050565b6000806040838503121561044657600080fd5b61044f8361040a565b915061045d6020840161040a565b90509250929050565b60008060006060848603121561047b57600080fd5b8335801515811461048b57600080fd5b9250602084013567ffffffffffffffff808211156104a857600080fd5b818601915086601f8301126104bc57600080fd5b8135818111156104ce576104ce61067d565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156105145761051461067d565b8160405282815289602084870101111561052d57600080fd5b826020860160208301376000602084830101528096505050505050604084013590509250925092565b6000806040838503121561056957600080fd5b8235915061045d6020840161040a565b6000815180845260005b8181101561059f57602081850181015186830182015201610583565b818111156105b1576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b83151581526060602082015260006105ff6060830185610579565b9050826040830152949350505050565b6020815260006106226020830184610579565b9392505050565b600181811c9082168061063d57607f821691505b60208210811415610677577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeeperRegistryCheckUpkeepGasUsageWrapperABI = KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.ABI

var KeeperRegistryCheckUpkeepGasUsageWrapperBin = KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.Bin

func DeployKeeperRegistryCheckUpkeepGasUsageWrapper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistryCheckUpkeepGasUsageWrapper, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryCheckUpkeepGasUsageWrapperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryCheckUpkeepGasUsageWrapper{address: address, abi: *parsed, KeeperRegistryCheckUpkeepGasUsageWrapperCaller: KeeperRegistryCheckUpkeepGasUsageWrapperCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperTransactor: KeeperRegistryCheckUpkeepGasUsageWrapperTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperFilterer: KeeperRegistryCheckUpkeepGasUsageWrapperFilterer{contract: contract}}, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapper struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryCheckUpkeepGasUsageWrapperCaller
	KeeperRegistryCheckUpkeepGasUsageWrapperTransactor
	KeeperRegistryCheckUpkeepGasUsageWrapperFilterer
}

type KeeperRegistryCheckUpkeepGasUsageWrapperCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapper
}

type KeeperRegistryCheckUpkeepGasUsageWrapperCallerRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperCaller
}

type KeeperRegistryCheckUpkeepGasUsageWrapperTransactorRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapper(address common.Address, backend bind.ContractBackend) (*KeeperRegistryCheckUpkeepGasUsageWrapper, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryCheckUpkeepGasUsageWrapperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper{address: address, abi: abi, KeeperRegistryCheckUpkeepGasUsageWrapperCaller: KeeperRegistryCheckUpkeepGasUsageWrapperCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperTransactor: KeeperRegistryCheckUpkeepGasUsageWrapperTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperFilterer: KeeperRegistryCheckUpkeepGasUsageWrapperFilterer{contract: contract}}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCheckUpkeepGasUsageWrapperCaller, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperCaller{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryCheckUpkeepGasUsageWrapperTransactor, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperTransactor{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryCheckUpkeepGasUsageWrapperFilterer, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperFilterer{contract: contract}, nil
}

func bindKeeperRegistryCheckUpkeepGasUsageWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCaller) SMockGas(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Call(opts, &out, "s_mockGas")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) SMockGas() (*big.Int, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SMockGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession) SMockGas() (*big.Int, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SMockGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCaller) SMockPayload(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Call(opts, &out, "s_mockPayload")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) SMockPayload() ([]byte, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SMockPayload(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession) SMockPayload() ([]byte, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SMockPayload(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCaller) SMockResult(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Call(opts, &out, "s_mockResult")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) SMockResult() (bool, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SMockResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession) SMockResult() (bool, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SMockResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.EmitOwnershipTransferred(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.EmitOwnershipTransferred(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "measureCheckGas", id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) SetMeasureCheckGasResult(opts *bind.TransactOpts, result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "setMeasureCheckGasResult", result, payload, gas)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SetMeasureCheckGasResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, result, payload, gas)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.SetMeasureCheckGasResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, result, payload, gas)
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapper.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapper) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryCheckUpkeepGasUsageWrapper.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapper.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryCheckUpkeepGasUsageWrapper.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapper.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapper) Address() common.Address {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.address
}

type KeeperRegistryCheckUpkeepGasUsageWrapperInterface interface {
	SMockGas(opts *bind.CallOpts) (*big.Int, error)

	SMockPayload(opts *bind.CallOpts) ([]byte, error)

	SMockResult(opts *bind.CallOpts) (bool, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	SetMeasureCheckGasResult(opts *bind.TransactOpts, result bool, payload []byte, gas *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
