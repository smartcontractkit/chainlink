// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_proxy_admin

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

var VRFProxyAdminMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractTransparentUpgradeableProxy\",\"name\":\"proxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeProxyAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractTransparentUpgradeableProxy\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"getProxyAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractTransparentUpgradeableProxy\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"getProxyImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractTransparentUpgradeableProxy\",\"name\":\"proxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractTransparentUpgradeableProxy\",\"name\":\"proxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061001a3361001f565b61006f565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6108a08061007e6000396000f3fe60806040526004361061007b5760003560e01c80639623609d1161004e5780639623609d1461011157806399a88ec414610124578063f2fde38b14610144578063f3b7dead1461016457600080fd5b8063204e1c7a14610080578063715018a6146100bc5780637eff275e146100d35780638da5cb5b146100f3575b600080fd5b34801561008c57600080fd5b506100a061009b366004610685565b610184565b6040516001600160a01b03909116815260200160405180910390f35b3480156100c857600080fd5b506100d161022e565b005b3480156100df57600080fd5b506100d16100ee3660046106a9565b610299565b3480156100ff57600080fd5b506000546001600160a01b03166100a0565b6100d161011f366004610711565b61036c565b34801561013057600080fd5b506100d161013f3660046106a9565b610446565b34801561015057600080fd5b506100d161015f366004610685565b6104e7565b34801561017057600080fd5b506100a061017f366004610685565b6105c9565b6000806000836001600160a01b03166040516101c3907f5c60da1b00000000000000000000000000000000000000000000000000000000815260040190565b600060405180830381855afa9150503d80600081146101fe576040519150601f19603f3d011682016040523d82523d6000602084013e610203565b606091505b50915091508161021257600080fd5b8080602001905181019061022691906107e7565b949350505050565b6000546001600160a01b0316331461028d5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b6102976000610608565b565b6000546001600160a01b031633146102f35760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610284565b6040517f8f2839700000000000000000000000000000000000000000000000000000000081526001600160a01b038281166004830152831690638f283970906024015b600060405180830381600087803b15801561035057600080fd5b505af1158015610364573d6000803e3d6000fd5b505050505050565b6000546001600160a01b031633146103c65760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610284565b6040517f4f1ef2860000000000000000000000000000000000000000000000000000000081526001600160a01b03841690634f1ef28690349061040f9086908690600401610804565b6000604051808303818588803b15801561042857600080fd5b505af115801561043c573d6000803e3d6000fd5b5050505050505050565b6000546001600160a01b031633146104a05760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610284565b6040517f3659cfe60000000000000000000000000000000000000000000000000000000081526001600160a01b038281166004830152831690633659cfe690602401610336565b6000546001600160a01b031633146105415760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610284565b6001600160a01b0381166105bd5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610284565b6105c681610608565b50565b6000806000836001600160a01b03166040516101c3907ff851a44000000000000000000000000000000000000000000000000000000000815260040190565b600080546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6001600160a01b03811681146105c657600080fd5b60006020828403121561069757600080fd5b81356106a281610670565b9392505050565b600080604083850312156106bc57600080fd5b82356106c781610670565b915060208301356106d781610670565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060006060848603121561072657600080fd5b833561073181610670565b9250602084013561074181610670565b9150604084013567ffffffffffffffff8082111561075e57600080fd5b818601915086601f83011261077257600080fd5b813581811115610784576107846106e2565b604051601f8201601f19908116603f011681019083821181831017156107ac576107ac6106e2565b816040528281528960208487010111156107c557600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b6000602082840312156107f957600080fd5b81516106a281610670565b6001600160a01b038316815260006020604081840152835180604085015260005b8181101561084157858101830151858201606001528201610825565b81811115610853576000606083870101525b50601f01601f19169290920160600194935050505056fea2646970667358221220531200594c54aee0c61b34c0b1b4cb8a3f4c6438e0071312d99f618f467910bf64736f6c634300080f0033",
}

var VRFProxyAdminABI = VRFProxyAdminMetaData.ABI

var VRFProxyAdminBin = VRFProxyAdminMetaData.Bin

func DeployVRFProxyAdmin(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFProxyAdmin, error) {
	parsed, err := VRFProxyAdminMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFProxyAdminBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFProxyAdmin{VRFProxyAdminCaller: VRFProxyAdminCaller{contract: contract}, VRFProxyAdminTransactor: VRFProxyAdminTransactor{contract: contract}, VRFProxyAdminFilterer: VRFProxyAdminFilterer{contract: contract}}, nil
}

type VRFProxyAdmin struct {
	address common.Address
	abi     abi.ABI
	VRFProxyAdminCaller
	VRFProxyAdminTransactor
	VRFProxyAdminFilterer
}

type VRFProxyAdminCaller struct {
	contract *bind.BoundContract
}

type VRFProxyAdminTransactor struct {
	contract *bind.BoundContract
}

type VRFProxyAdminFilterer struct {
	contract *bind.BoundContract
}

type VRFProxyAdminSession struct {
	Contract     *VRFProxyAdmin
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFProxyAdminCallerSession struct {
	Contract *VRFProxyAdminCaller
	CallOpts bind.CallOpts
}

type VRFProxyAdminTransactorSession struct {
	Contract     *VRFProxyAdminTransactor
	TransactOpts bind.TransactOpts
}

type VRFProxyAdminRaw struct {
	Contract *VRFProxyAdmin
}

type VRFProxyAdminCallerRaw struct {
	Contract *VRFProxyAdminCaller
}

type VRFProxyAdminTransactorRaw struct {
	Contract *VRFProxyAdminTransactor
}

func NewVRFProxyAdmin(address common.Address, backend bind.ContractBackend) (*VRFProxyAdmin, error) {
	abi, err := abi.JSON(strings.NewReader(VRFProxyAdminABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFProxyAdmin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFProxyAdmin{address: address, abi: abi, VRFProxyAdminCaller: VRFProxyAdminCaller{contract: contract}, VRFProxyAdminTransactor: VRFProxyAdminTransactor{contract: contract}, VRFProxyAdminFilterer: VRFProxyAdminFilterer{contract: contract}}, nil
}

func NewVRFProxyAdminCaller(address common.Address, caller bind.ContractCaller) (*VRFProxyAdminCaller, error) {
	contract, err := bindVRFProxyAdmin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFProxyAdminCaller{contract: contract}, nil
}

func NewVRFProxyAdminTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFProxyAdminTransactor, error) {
	contract, err := bindVRFProxyAdmin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFProxyAdminTransactor{contract: contract}, nil
}

func NewVRFProxyAdminFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFProxyAdminFilterer, error) {
	contract, err := bindVRFProxyAdmin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFProxyAdminFilterer{contract: contract}, nil
}

func bindVRFProxyAdmin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFProxyAdminABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFProxyAdmin *VRFProxyAdminRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFProxyAdmin.Contract.VRFProxyAdminCaller.contract.Call(opts, result, method, params...)
}

func (_VRFProxyAdmin *VRFProxyAdminRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.VRFProxyAdminTransactor.contract.Transfer(opts)
}

func (_VRFProxyAdmin *VRFProxyAdminRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.VRFProxyAdminTransactor.contract.Transact(opts, method, params...)
}

func (_VRFProxyAdmin *VRFProxyAdminCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFProxyAdmin.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.contract.Transfer(opts)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.contract.Transact(opts, method, params...)
}

func (_VRFProxyAdmin *VRFProxyAdminCaller) GetProxyAdmin(opts *bind.CallOpts, proxy common.Address) (common.Address, error) {
	var out []interface{}
	err := _VRFProxyAdmin.contract.Call(opts, &out, "getProxyAdmin", proxy)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFProxyAdmin *VRFProxyAdminSession) GetProxyAdmin(proxy common.Address) (common.Address, error) {
	return _VRFProxyAdmin.Contract.GetProxyAdmin(&_VRFProxyAdmin.CallOpts, proxy)
}

func (_VRFProxyAdmin *VRFProxyAdminCallerSession) GetProxyAdmin(proxy common.Address) (common.Address, error) {
	return _VRFProxyAdmin.Contract.GetProxyAdmin(&_VRFProxyAdmin.CallOpts, proxy)
}

func (_VRFProxyAdmin *VRFProxyAdminCaller) GetProxyImplementation(opts *bind.CallOpts, proxy common.Address) (common.Address, error) {
	var out []interface{}
	err := _VRFProxyAdmin.contract.Call(opts, &out, "getProxyImplementation", proxy)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFProxyAdmin *VRFProxyAdminSession) GetProxyImplementation(proxy common.Address) (common.Address, error) {
	return _VRFProxyAdmin.Contract.GetProxyImplementation(&_VRFProxyAdmin.CallOpts, proxy)
}

func (_VRFProxyAdmin *VRFProxyAdminCallerSession) GetProxyImplementation(proxy common.Address) (common.Address, error) {
	return _VRFProxyAdmin.Contract.GetProxyImplementation(&_VRFProxyAdmin.CallOpts, proxy)
}

func (_VRFProxyAdmin *VRFProxyAdminCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFProxyAdmin.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFProxyAdmin *VRFProxyAdminSession) Owner() (common.Address, error) {
	return _VRFProxyAdmin.Contract.Owner(&_VRFProxyAdmin.CallOpts)
}

func (_VRFProxyAdmin *VRFProxyAdminCallerSession) Owner() (common.Address, error) {
	return _VRFProxyAdmin.Contract.Owner(&_VRFProxyAdmin.CallOpts)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactor) ChangeProxyAdmin(opts *bind.TransactOpts, proxy common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.contract.Transact(opts, "changeProxyAdmin", proxy, newAdmin)
}

func (_VRFProxyAdmin *VRFProxyAdminSession) ChangeProxyAdmin(proxy common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.ChangeProxyAdmin(&_VRFProxyAdmin.TransactOpts, proxy, newAdmin)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorSession) ChangeProxyAdmin(proxy common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.ChangeProxyAdmin(&_VRFProxyAdmin.TransactOpts, proxy, newAdmin)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFProxyAdmin.contract.Transact(opts, "renounceOwnership")
}

func (_VRFProxyAdmin *VRFProxyAdminSession) RenounceOwnership() (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.RenounceOwnership(&_VRFProxyAdmin.TransactOpts)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.RenounceOwnership(&_VRFProxyAdmin.TransactOpts)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_VRFProxyAdmin *VRFProxyAdminSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.TransferOwnership(&_VRFProxyAdmin.TransactOpts, newOwner)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.TransferOwnership(&_VRFProxyAdmin.TransactOpts, newOwner)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactor) Upgrade(opts *bind.TransactOpts, proxy common.Address, implementation common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.contract.Transact(opts, "upgrade", proxy, implementation)
}

func (_VRFProxyAdmin *VRFProxyAdminSession) Upgrade(proxy common.Address, implementation common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.Upgrade(&_VRFProxyAdmin.TransactOpts, proxy, implementation)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorSession) Upgrade(proxy common.Address, implementation common.Address) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.Upgrade(&_VRFProxyAdmin.TransactOpts, proxy, implementation)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactor) UpgradeAndCall(opts *bind.TransactOpts, proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFProxyAdmin.contract.Transact(opts, "upgradeAndCall", proxy, implementation, data)
}

func (_VRFProxyAdmin *VRFProxyAdminSession) UpgradeAndCall(proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.UpgradeAndCall(&_VRFProxyAdmin.TransactOpts, proxy, implementation, data)
}

func (_VRFProxyAdmin *VRFProxyAdminTransactorSession) UpgradeAndCall(proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFProxyAdmin.Contract.UpgradeAndCall(&_VRFProxyAdmin.TransactOpts, proxy, implementation, data)
}

type VRFProxyAdminOwnershipTransferredIterator struct {
	Event *VRFProxyAdminOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFProxyAdminOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFProxyAdminOwnershipTransferred)
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
		it.Event = new(VRFProxyAdminOwnershipTransferred)
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

func (it *VRFProxyAdminOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFProxyAdminOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFProxyAdminOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log
}

func (_VRFProxyAdmin *VRFProxyAdminFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VRFProxyAdminOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VRFProxyAdmin.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VRFProxyAdminOwnershipTransferredIterator{contract: _VRFProxyAdmin.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFProxyAdmin *VRFProxyAdminFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFProxyAdminOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VRFProxyAdmin.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFProxyAdminOwnershipTransferred)
				if err := _VRFProxyAdmin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFProxyAdmin *VRFProxyAdminFilterer) ParseOwnershipTransferred(log types.Log) (*VRFProxyAdminOwnershipTransferred, error) {
	event := new(VRFProxyAdminOwnershipTransferred)
	if err := _VRFProxyAdmin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFProxyAdmin *VRFProxyAdmin) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFProxyAdmin.abi.Events["OwnershipTransferred"].ID:
		return _VRFProxyAdmin.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFProxyAdminOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFProxyAdmin *VRFProxyAdmin) Address() common.Address {
	return _VRFProxyAdmin.address
}

type VRFProxyAdminInterface interface {
	GetProxyAdmin(opts *bind.CallOpts, proxy common.Address) (common.Address, error)

	GetProxyImplementation(opts *bind.CallOpts, proxy common.Address) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	ChangeProxyAdmin(opts *bind.TransactOpts, proxy common.Address, newAdmin common.Address) (*types.Transaction, error)

	RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	Upgrade(opts *bind.TransactOpts, proxy common.Address, implementation common.Address) (*types.Transaction, error)

	UpgradeAndCall(opts *bind.TransactOpts, proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VRFProxyAdminOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFProxyAdminOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFProxyAdminOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
