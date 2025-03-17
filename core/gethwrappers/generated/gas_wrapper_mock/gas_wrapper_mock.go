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

var KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"measureCheckGas\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_mockGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_mockPayload\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_mockResult\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"}],\"name\":\"setMeasureCheckGasResult\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506106b9806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063846811931161005b57806384681193146100d5578063b019b4e8146100f2578063b023145014610105578063f7420bc21461011a57600080fd5b80632dae06f51461008257806356343496146100975780636bf49030146100b3575b600080fd5b610095610090366004610466565b61012d565b005b6100a060025481565b6040519081526020015b60405180910390f35b6100c66100c1366004610556565b610174565b6040516100aa939291906105e4565b6000546100e29060ff1681565b60405190151581526020016100aa565b610095610100366004610433565b610227565b61010d610285565b6040516100aa919061060f565b610095610128366004610433565b610313565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016841515179055815161016c906001906020850190610371565b506002555050565b6000606060008060009054906101000a900460ff16600160025481805461019a90610629565b80601f01602080910402602001604051908101604052809291908181526020018280546101c690610629565b80156102135780601f106101e857610100808354040283529160200191610213565b820191906000526020600020905b8154815290600101906020018083116101f657829003601f168201915b505050505091509250925092509250925092565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6001805461029290610629565b80601f01602080910402602001604051908101604052809291908181526020018280546102be90610629565b801561030b5780601f106102e05761010080835404028352916020019161030b565b820191906000526020600020905b8154815290600101906020018083116102ee57829003601f168201915b505050505081565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b82805461037d90610629565b90600052602060002090601f01602090048101928261039f57600085556103e5565b82601f106103b857805160ff19168380011785556103e5565b828001600101855582156103e5579182015b828111156103e55782518255916020019190600101906103ca565b506103f19291506103f5565b5090565b5b808211156103f157600081556001016103f6565b803573ffffffffffffffffffffffffffffffffffffffff8116811461042e57600080fd5b919050565b6000806040838503121561044657600080fd5b61044f8361040a565b915061045d6020840161040a565b90509250929050565b60008060006060848603121561047b57600080fd5b8335801515811461048b57600080fd5b9250602084013567ffffffffffffffff808211156104a857600080fd5b818601915086601f8301126104bc57600080fd5b8135818111156104ce576104ce61067d565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156105145761051461067d565b8160405282815289602084870101111561052d57600080fd5b826020860160208301376000602084830101528096505050505050604084013590509250925092565b6000806040838503121561056957600080fd5b8235915061045d6020840161040a565b6000815180845260005b8181101561059f57602081850181015186830182015201610583565b818111156105b1576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b83151581526060602082015260006105ff6060830185610579565b9050826040830152949350505050565b6020815260006106226020830184610579565b9392505050565b600181811c9082168061063d57607f821691505b60208210811415610677577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeeperRegistryCheckUpkeepGasUsageWrapperMockABI = KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.ABI

var KeeperRegistryCheckUpkeepGasUsageWrapperMockBin = KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.Bin

func DeployKeeperRegistryCheckUpkeepGasUsageWrapperMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistryCheckUpkeepGasUsageWrapperMock, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryCheckUpkeepGasUsageWrapperMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryCheckUpkeepGasUsageWrapperMock{address: address, abi: *parsed, KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller: KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor: KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer: KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer{contract: contract}}, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMock struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller
	KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor
	KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapperMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockCallerSession struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperMock
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockCallerRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperMock(address common.Address, backend bind.ContractBackend) (*KeeperRegistryCheckUpkeepGasUsageWrapperMock, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryCheckUpkeepGasUsageWrapperMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapperMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperMock{address: address, abi: abi, KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller: KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor: KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer: KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer{contract: contract}}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperMockCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapperMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapperMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapperMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer{contract: contract}, nil
}

func bindKeeperRegistryCheckUpkeepGasUsageWrapperMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller) SMockGas(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Call(opts, &out, "s_mockGas")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) SMockGas() (*big.Int, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SMockGas(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCallerSession) SMockGas() (*big.Int, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SMockGas(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller) SMockPayload(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Call(opts, &out, "s_mockPayload")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) SMockPayload() ([]byte, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SMockPayload(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCallerSession) SMockPayload() ([]byte, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SMockPayload(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCaller) SMockResult(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Call(opts, &out, "s_mockResult")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) SMockResult() (bool, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SMockResult(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockCallerSession) SMockResult() (bool, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SMockResult(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.EmitOwnershipTransferred(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.EmitOwnershipTransferred(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, from, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor) MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Transact(opts, "measureCheckGas", id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactor) SetMeasureCheckGasResult(opts *bind.TransactOpts, result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.Transact(opts, "setMeasureCheckGasResult", result, payload, gas)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockSession) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SetMeasureCheckGasResult(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, result, payload, gas)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockTransactorSession) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.Contract.SetMeasureCheckGasResult(&_KeeperRegistryCheckUpkeepGasUsageWrapperMock.TransactOpts, result, payload, gas)
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMockFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapperMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryCheckUpkeepGasUsageWrapperMock.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryCheckUpkeepGasUsageWrapperMock.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapperMock *KeeperRegistryCheckUpkeepGasUsageWrapperMock) Address() common.Address {
	return _KeeperRegistryCheckUpkeepGasUsageWrapperMock.address
}

type KeeperRegistryCheckUpkeepGasUsageWrapperMockInterface interface {
	SMockGas(opts *bind.CallOpts) (*big.Int, error)

	SMockPayload(opts *bind.CallOpts) ([]byte, error)

	SMockResult(opts *bind.CallOpts) (bool, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	SetMeasureCheckGasResult(opts *bind.TransactOpts, result bool, payload []byte, gas *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperMockOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
