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

var KeeperRegistryCheckUpkeepGasUsageWrapper12MockMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"emitOwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitOwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"measureCheckGas\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"s_mockGas\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_mockPayload\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_mockResult\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setMeasureCheckGasResult\",\"inputs\":[{\"name\":\"result\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"gas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b506106b9806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063846811931161005b57806384681193146100d5578063b019b4e8146100f2578063b023145014610105578063f7420bc21461011a57600080fd5b80632dae06f51461008257806356343496146100975780636bf49030146100b3575b600080fd5b610095610090366004610466565b61012d565b005b6100a060025481565b6040519081526020015b60405180910390f35b6100c66100c1366004610556565b610174565b6040516100aa939291906105e4565b6000546100e29060ff1681565b60405190151581526020016100aa565b610095610100366004610433565b610227565b61010d610285565b6040516100aa919061060f565b610095610128366004610433565b610313565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016841515179055815161016c906001906020850190610371565b506002555050565b6000606060008060009054906101000a900460ff16600160025481805461019a90610629565b80601f01602080910402602001604051908101604052809291908181526020018280546101c690610629565b80156102135780601f106101e857610100808354040283529160200191610213565b820191906000526020600020905b8154815290600101906020018083116101f657829003601f168201915b505050505091509250925092509250925092565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6001805461029290610629565b80601f01602080910402602001604051908101604052809291908181526020018280546102be90610629565b801561030b5780601f106102e05761010080835404028352916020019161030b565b820191906000526020600020905b8154815290600101906020018083116102ee57829003601f168201915b505050505081565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b82805461037d90610629565b90600052602060002090601f01602090048101928261039f57600085556103e5565b82601f106103b857805160ff19168380011785556103e5565b828001600101855582156103e5579182015b828111156103e55782518255916020019190600101906103ca565b506103f19291506103f5565b5090565b5b808211156103f157600081556001016103f6565b803573ffffffffffffffffffffffffffffffffffffffff8116811461042e57600080fd5b919050565b6000806040838503121561044657600080fd5b61044f8361040a565b915061045d6020840161040a565b90509250929050565b60008060006060848603121561047b57600080fd5b8335801515811461048b57600080fd5b9250602084013567ffffffffffffffff808211156104a857600080fd5b818601915086601f8301126104bc57600080fd5b8135818111156104ce576104ce61067d565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156105145761051461067d565b8160405282815289602084870101111561052d57600080fd5b826020860160208301376000602084830101528096505050505050604084013590509250925092565b6000806040838503121561056957600080fd5b8235915061045d6020840161040a565b6000815180845260005b8181101561059f57602081850181015186830182015201610583565b818111156105b1576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b83151581526060602082015260006105ff6060830185610579565b9050826040830152949350505050565b6020815260006106226020830184610579565b9392505050565b600181811c9082168061063d57607f821691505b60208210811415610677577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeeperRegistryCheckUpkeepGasUsageWrapper12MockABI = KeeperRegistryCheckUpkeepGasUsageWrapper12MockMetaData.ABI

var KeeperRegistryCheckUpkeepGasUsageWrapper12MockBin = KeeperRegistryCheckUpkeepGasUsageWrapper12MockMetaData.Bin

func DeployKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistryCheckUpkeepGasUsageWrapper12Mock, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapper12MockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryCheckUpkeepGasUsageWrapper12MockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryCheckUpkeepGasUsageWrapper12Mock{address: address, abi: *parsed, KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller: KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor: KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer: KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer{contract: contract}}, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12Mock struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller
	KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor
	KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapper12Mock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockCallerSession struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapper12Mock
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockCallerRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(address common.Address, backend bind.ContractBackend) (*KeeperRegistryCheckUpkeepGasUsageWrapper12Mock, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryCheckUpkeepGasUsageWrapper12MockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper12Mock{address: address, abi: abi, KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller: KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor: KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer: KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer{contract: contract}}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer{contract: contract}, nil
}

func bindKeeperRegistryCheckUpkeepGasUsageWrapper12Mock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapper12MockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller) SMockGas(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Call(opts, &out, "s_mockGas")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) SMockGas() (*big.Int, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SMockGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCallerSession) SMockGas() (*big.Int, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SMockGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller) SMockPayload(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Call(opts, &out, "s_mockPayload")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) SMockPayload() ([]byte, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SMockPayload(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCallerSession) SMockPayload() ([]byte, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SMockPayload(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCaller) SMockResult(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Call(opts, &out, "s_mockResult")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) SMockResult() (bool, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SMockResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockCallerSession) SMockResult() (bool, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SMockResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.EmitOwnershipTransferred(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.EmitOwnershipTransferred(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor) MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Transact(opts, "measureCheckGas", id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactor) SetMeasureCheckGasResult(opts *bind.TransactOpts, result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.Transact(opts, "setMeasureCheckGasResult", result, payload, gas)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockSession) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SetMeasureCheckGasResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, result, payload, gas)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockTransactorSession) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.Contract.SetMeasureCheckGasResult(&_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.TransactOpts, result, payload, gas)
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12MockFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12Mock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper12Mock *KeeperRegistryCheckUpkeepGasUsageWrapper12Mock) Address() common.Address {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper12Mock.address
}

type KeeperRegistryCheckUpkeepGasUsageWrapper12MockInterface interface {
	SMockGas(opts *bind.CallOpts) (*big.Int, error)

	SMockPayload(opts *bind.CallOpts) ([]byte, error)

	SMockResult(opts *bind.CallOpts) (bool, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	SetMeasureCheckGasResult(opts *bind.TransactOpts, result bool, payload []byte, gas *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapper12MockOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
