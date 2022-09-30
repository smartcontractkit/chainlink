// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_beacon_proxy

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

var VRFBeaconProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"implementation_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000f9238038062000f9283398101604081905262000034916200051b565b82828282816200006660017f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbd620005fb565b60008051602062000f4b8339815191521462000086576200008662000621565b6200009482826000620000fb565b50620000c4905060017fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6104620005fb565b60008051602062000f2b83398151915214620000e457620000e462000621565b620000ef8262000138565b5050505050506200068a565b620001068362000193565b600082511180620001145750805b156200013357620001318383620001d560201b6200023e1760201c565b505b505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6200016362000204565b604080516001600160a01b03928316815291841660208301520160405180910390a162000190816200023d565b50565b6200019e81620002f2565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060620001fd838360405180606001604052806027815260200162000f6b6027913962000395565b9392505050565b60006200022e60008051602062000f2b83398151915260001b6200047260201b620001fa1760201c565b546001600160a01b0316919050565b6001600160a01b038116620002a85760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b80620002d160008051602062000f2b83398151915260001b6200047260201b620001fa1760201c565b80546001600160a01b0319166001600160a01b039290921691909117905550565b62000308816200047560201b6200026a1760201c565b6200036c5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016200029f565b80620002d160008051602062000f4b83398151915260001b6200047260201b620001fa1760201c565b6060833b620003f65760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b60648201526084016200029f565b600080856001600160a01b03168560405162000413919062000637565b600060405180830381855af49150503d806000811462000450576040519150601f19603f3d011682016040523d82523d6000602084013e62000455565b606091505b509092509050620004688282866200047b565b9695505050505050565b90565b3b151590565b606083156200048c575081620001fd565b8251156200049d5782518084602001fd5b8160405162461bcd60e51b81526004016200029f919062000655565b80516001600160a01b0381168114620004d157600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b8381101562000509578181015183820152602001620004ef565b83811115620001315750506000910152565b6000806000606084860312156200053157600080fd5b6200053c84620004b9565b92506200054c60208501620004b9565b60408501519092506001600160401b03808211156200056a57600080fd5b818601915086601f8301126200057f57600080fd5b815181811115620005945762000594620004d6565b604051601f8201601f19908116603f01168101908382118183101715620005bf57620005bf620004d6565b81604052828152896020848701011115620005d957600080fd5b620005ec836020830160208801620004ec565b80955050505050509250925092565b6000828210156200061c57634e487b7160e01b600052601160045260246000fd5b500390565b634e487b7160e01b600052600160045260246000fd5b600082516200064b818460208701620004ec565b9190910192915050565b602081526000825180602084015262000676816040850160208701620004ec565b601f01601f19169190910160400192915050565b610891806200069a6000396000f3fe60806040526004361061005e5760003560e01c80635c60da1b116100435780635c60da1b146100a85780638f283970146100d9578063f851a440146100f95761006d565b80633659cfe6146100755780634f1ef286146100955761006d565b3661006d5761006b61010e565b005b61006b61010e565b34801561008157600080fd5b5061006b61009036600461071b565b610128565b61006b6100a3366004610736565b610165565b3480156100b457600080fd5b506100bd6101cc565b6040516001600160a01b03909116815260200160405180910390f35b3480156100e557600080fd5b5061006b6100f436600461071b565b6101fd565b34801561010557600080fd5b506100bd61021d565b610116610270565b610126610121610320565b61032a565b565b61013061034e565b6001600160a01b0316330361015d5761015a81604051806020016040528060008152506000610381565b50565b61015a61010e565b61016d61034e565b6001600160a01b031633036101c4576101bf8383838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250610381915050565b505050565b6101bf61010e565b60006101d661034e565b6001600160a01b031633036101f2576101ed610320565b905090565b6101fa61010e565b90565b61020561034e565b6001600160a01b0316330361015d5761015a816103ac565b600061022761034e565b6001600160a01b031633036101f2576101ed61034e565b6060610263838360405180606001604052806027815260200161083560279139610400565b9392505050565b3b151590565b61027861034e565b6001600160a01b031633036101265760405162461bcd60e51b815260206004820152604260248201527f5472616e73706172656e745570677261646561626c6550726f78793a2061646d60448201527f696e2063616e6e6f742066616c6c6261636b20746f2070726f7879207461726760648201527f6574000000000000000000000000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b60006101ed6104eb565b3660008037600080366000845af43d6000803e808015610349573d6000f35b3d6000fd5b60007fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b546001600160a01b0316919050565b61038a83610513565b6000825111806103975750805b156101bf576103a6838361023e565b50505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6103d561034e565b604080516001600160a01b03928316815291841660208301520160405180910390a161015a81610553565b6060833b6104765760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f60448201527f6e747261637400000000000000000000000000000000000000000000000000006064820152608401610317565b600080856001600160a01b03168560405161049191906107e5565b600060405180830381855af49150503d80600081146104cc576040519150601f19603f3d011682016040523d82523d6000602084013e6104d1565b606091505b50915091506104e182828661062b565b9695505050505050565b60007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610372565b61051c81610664565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6001600160a01b0381166105cf5760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610317565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b039290921691909117905550565b6060831561063a575081610263565b82511561064a5782518084602001fd5b8160405162461bcd60e51b81526004016103179190610801565b803b6106d85760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e7472616374000000000000000000000000000000000000006064820152608401610317565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6105f2565b80356001600160a01b038116811461071657600080fd5b919050565b60006020828403121561072d57600080fd5b610263826106ff565b60008060006040848603121561074b57600080fd5b610754846106ff565b9250602084013567ffffffffffffffff8082111561077157600080fd5b818601915086601f83011261078557600080fd5b81358181111561079457600080fd5b8760208285010111156107a657600080fd5b6020830194508093505050509250925092565b60005b838110156107d45781810151838201526020016107bc565b838111156103a65750506000910152565b600082516107f78184602087016107b9565b9190910192915050565b60208152600082518060208401526108208160408501602087016107b9565b601f01601f1916919091016040019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220a1ed96031f9cbfad60f98c50fd3df87b7d87d4c13cac62a08addc168afb9d30e64736f6c634300080f0033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
}

var VRFBeaconProxyABI = VRFBeaconProxyMetaData.ABI

var VRFBeaconProxyBin = VRFBeaconProxyMetaData.Bin

func DeployVRFBeaconProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, _admin common.Address, _data []byte) (common.Address, *types.Transaction, *VRFBeaconProxy, error) {
	parsed, err := VRFBeaconProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFBeaconProxyBin), backend, _logic, _admin, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFBeaconProxy{VRFBeaconProxyCaller: VRFBeaconProxyCaller{contract: contract}, VRFBeaconProxyTransactor: VRFBeaconProxyTransactor{contract: contract}, VRFBeaconProxyFilterer: VRFBeaconProxyFilterer{contract: contract}}, nil
}

type VRFBeaconProxy struct {
	address common.Address
	abi     abi.ABI
	VRFBeaconProxyCaller
	VRFBeaconProxyTransactor
	VRFBeaconProxyFilterer
}

type VRFBeaconProxyCaller struct {
	contract *bind.BoundContract
}

type VRFBeaconProxyTransactor struct {
	contract *bind.BoundContract
}

type VRFBeaconProxyFilterer struct {
	contract *bind.BoundContract
}

type VRFBeaconProxySession struct {
	Contract     *VRFBeaconProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFBeaconProxyCallerSession struct {
	Contract *VRFBeaconProxyCaller
	CallOpts bind.CallOpts
}

type VRFBeaconProxyTransactorSession struct {
	Contract     *VRFBeaconProxyTransactor
	TransactOpts bind.TransactOpts
}

type VRFBeaconProxyRaw struct {
	Contract *VRFBeaconProxy
}

type VRFBeaconProxyCallerRaw struct {
	Contract *VRFBeaconProxyCaller
}

type VRFBeaconProxyTransactorRaw struct {
	Contract *VRFBeaconProxyTransactor
}

func NewVRFBeaconProxy(address common.Address, backend bind.ContractBackend) (*VRFBeaconProxy, error) {
	abi, err := abi.JSON(strings.NewReader(VRFBeaconProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFBeaconProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxy{address: address, abi: abi, VRFBeaconProxyCaller: VRFBeaconProxyCaller{contract: contract}, VRFBeaconProxyTransactor: VRFBeaconProxyTransactor{contract: contract}, VRFBeaconProxyFilterer: VRFBeaconProxyFilterer{contract: contract}}, nil
}

func NewVRFBeaconProxyCaller(address common.Address, caller bind.ContractCaller) (*VRFBeaconProxyCaller, error) {
	contract, err := bindVRFBeaconProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxyCaller{contract: contract}, nil
}

func NewVRFBeaconProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFBeaconProxyTransactor, error) {
	contract, err := bindVRFBeaconProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxyTransactor{contract: contract}, nil
}

func NewVRFBeaconProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFBeaconProxyFilterer, error) {
	contract, err := bindVRFBeaconProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxyFilterer{contract: contract}, nil
}

func bindVRFBeaconProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFBeaconProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFBeaconProxy *VRFBeaconProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFBeaconProxy.Contract.VRFBeaconProxyCaller.contract.Call(opts, result, method, params...)
}

func (_VRFBeaconProxy *VRFBeaconProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.VRFBeaconProxyTransactor.contract.Transfer(opts)
}

func (_VRFBeaconProxy *VRFBeaconProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.VRFBeaconProxyTransactor.contract.Transact(opts, method, params...)
}

func (_VRFBeaconProxy *VRFBeaconProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFBeaconProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.contract.Transfer(opts)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.contract.Transact(opts, method, params...)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.Transact(opts, "admin")
}

func (_VRFBeaconProxy *VRFBeaconProxySession) Admin() (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Admin(&_VRFBeaconProxy.TransactOpts)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Admin(&_VRFBeaconProxy.TransactOpts)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

func (_VRFBeaconProxy *VRFBeaconProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.ChangeAdmin(&_VRFBeaconProxy.TransactOpts, newAdmin)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.ChangeAdmin(&_VRFBeaconProxy.TransactOpts, newAdmin)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.Transact(opts, "implementation")
}

func (_VRFBeaconProxy *VRFBeaconProxySession) Implementation() (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Implementation(&_VRFBeaconProxy.TransactOpts)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Implementation(&_VRFBeaconProxy.TransactOpts)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

func (_VRFBeaconProxy *VRFBeaconProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.UpgradeTo(&_VRFBeaconProxy.TransactOpts, newImplementation)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.UpgradeTo(&_VRFBeaconProxy.TransactOpts, newImplementation)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

func (_VRFBeaconProxy *VRFBeaconProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.UpgradeToAndCall(&_VRFBeaconProxy.TransactOpts, newImplementation, data)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.UpgradeToAndCall(&_VRFBeaconProxy.TransactOpts, newImplementation, data)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.RawTransact(opts, calldata)
}

func (_VRFBeaconProxy *VRFBeaconProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Fallback(&_VRFBeaconProxy.TransactOpts, calldata)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Fallback(&_VRFBeaconProxy.TransactOpts, calldata)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconProxy.contract.RawTransact(opts, nil)
}

func (_VRFBeaconProxy *VRFBeaconProxySession) Receive() (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Receive(&_VRFBeaconProxy.TransactOpts)
}

func (_VRFBeaconProxy *VRFBeaconProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _VRFBeaconProxy.Contract.Receive(&_VRFBeaconProxy.TransactOpts)
}

type VRFBeaconProxyAdminChangedIterator struct {
	Event *VRFBeaconProxyAdminChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconProxyAdminChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconProxyAdminChanged)
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
		it.Event = new(VRFBeaconProxyAdminChanged)
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

func (it *VRFBeaconProxyAdminChangedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log
}

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*VRFBeaconProxyAdminChangedIterator, error) {

	logs, sub, err := _VRFBeaconProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxyAdminChangedIterator{contract: _VRFBeaconProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFBeaconProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _VRFBeaconProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconProxyAdminChanged)
				if err := _VRFBeaconProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) ParseAdminChanged(log types.Log) (*VRFBeaconProxyAdminChanged, error) {
	event := new(VRFBeaconProxyAdminChanged)
	if err := _VRFBeaconProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconProxyBeaconUpgradedIterator struct {
	Event *VRFBeaconProxyBeaconUpgraded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconProxyBeaconUpgradedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconProxyBeaconUpgraded)
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
		it.Event = new(VRFBeaconProxyBeaconUpgraded)
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

func (it *VRFBeaconProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log
}

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*VRFBeaconProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _VRFBeaconProxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxyBeaconUpgradedIterator{contract: _VRFBeaconProxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *VRFBeaconProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _VRFBeaconProxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconProxyBeaconUpgraded)
				if err := _VRFBeaconProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) ParseBeaconUpgraded(log types.Log) (*VRFBeaconProxyBeaconUpgraded, error) {
	event := new(VRFBeaconProxyBeaconUpgraded)
	if err := _VRFBeaconProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconProxyUpgradedIterator struct {
	Event *VRFBeaconProxyUpgraded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconProxyUpgradedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconProxyUpgraded)
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
		it.Event = new(VRFBeaconProxyUpgraded)
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

func (it *VRFBeaconProxyUpgradedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log
}

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VRFBeaconProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _VRFBeaconProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconProxyUpgradedIterator{contract: _VRFBeaconProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VRFBeaconProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _VRFBeaconProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconProxyUpgraded)
				if err := _VRFBeaconProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

func (_VRFBeaconProxy *VRFBeaconProxyFilterer) ParseUpgraded(log types.Log) (*VRFBeaconProxyUpgraded, error) {
	event := new(VRFBeaconProxyUpgraded)
	if err := _VRFBeaconProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFBeaconProxy *VRFBeaconProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFBeaconProxy.abi.Events["AdminChanged"].ID:
		return _VRFBeaconProxy.ParseAdminChanged(log)
	case _VRFBeaconProxy.abi.Events["BeaconUpgraded"].ID:
		return _VRFBeaconProxy.ParseBeaconUpgraded(log)
	case _VRFBeaconProxy.abi.Events["Upgraded"].ID:
		return _VRFBeaconProxy.ParseUpgraded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFBeaconProxyAdminChanged) Topic() common.Hash {
	return common.HexToHash("0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f")
}

func (VRFBeaconProxyBeaconUpgraded) Topic() common.Hash {
	return common.HexToHash("0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e")
}

func (VRFBeaconProxyUpgraded) Topic() common.Hash {
	return common.HexToHash("0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b")
}

func (_VRFBeaconProxy *VRFBeaconProxy) Address() common.Address {
	return _VRFBeaconProxy.address
}

type VRFBeaconProxyInterface interface {
	Admin(opts *bind.TransactOpts) (*types.Transaction, error)

	ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Implementation(opts *bind.TransactOpts) (*types.Transaction, error)

	UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error)

	UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterAdminChanged(opts *bind.FilterOpts) (*VRFBeaconProxyAdminChangedIterator, error)

	WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFBeaconProxyAdminChanged) (event.Subscription, error)

	ParseAdminChanged(log types.Log) (*VRFBeaconProxyAdminChanged, error)

	FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*VRFBeaconProxyBeaconUpgradedIterator, error)

	WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *VRFBeaconProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error)

	ParseBeaconUpgraded(log types.Log) (*VRFBeaconProxyBeaconUpgraded, error)

	FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VRFBeaconProxyUpgradedIterator, error)

	WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VRFBeaconProxyUpgraded, implementation []common.Address) (event.Subscription, error)

	ParseUpgraded(log types.Log) (*VRFBeaconProxyUpgraded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
