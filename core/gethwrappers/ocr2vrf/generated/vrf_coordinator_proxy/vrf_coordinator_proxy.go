// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_proxy

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

var VRFCoordinatorProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"implementation_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000f9238038062000f9283398101604081905262000034916200051b565b82828282816200006660017f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbd620005fb565b60008051602062000f4b8339815191521462000086576200008662000621565b6200009482826000620000fb565b50620000c4905060017fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6104620005fb565b60008051602062000f2b83398151915214620000e457620000e462000621565b620000ef8262000138565b5050505050506200068a565b620001068362000193565b600082511180620001145750805b156200013357620001318383620001d560201b6200023e1760201c565b505b505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6200016362000204565b604080516001600160a01b03928316815291841660208301520160405180910390a162000190816200023d565b50565b6200019e81620002f2565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060620001fd838360405180606001604052806027815260200162000f6b6027913962000395565b9392505050565b60006200022e60008051602062000f2b83398151915260001b6200047260201b620001fa1760201c565b546001600160a01b0316919050565b6001600160a01b038116620002a85760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b80620002d160008051602062000f2b83398151915260001b6200047260201b620001fa1760201c565b80546001600160a01b0319166001600160a01b039290921691909117905550565b62000308816200047560201b6200026a1760201c565b6200036c5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016200029f565b80620002d160008051602062000f4b83398151915260001b6200047260201b620001fa1760201c565b6060833b620003f65760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b60648201526084016200029f565b600080856001600160a01b03168560405162000413919062000637565b600060405180830381855af49150503d806000811462000450576040519150601f19603f3d011682016040523d82523d6000602084013e62000455565b606091505b509092509050620004688282866200047b565b9695505050505050565b90565b3b151590565b606083156200048c575081620001fd565b8251156200049d5782518084602001fd5b8160405162461bcd60e51b81526004016200029f919062000655565b80516001600160a01b0381168114620004d157600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b8381101562000509578181015183820152602001620004ef565b83811115620001315750506000910152565b6000806000606084860312156200053157600080fd5b6200053c84620004b9565b92506200054c60208501620004b9565b60408501519092506001600160401b03808211156200056a57600080fd5b818601915086601f8301126200057f57600080fd5b815181811115620005945762000594620004d6565b604051601f8201601f19908116603f01168101908382118183101715620005bf57620005bf620004d6565b81604052828152896020848701011115620005d957600080fd5b620005ec836020830160208801620004ec565b80955050505050509250925092565b6000828210156200061c57634e487b7160e01b600052601160045260246000fd5b500390565b634e487b7160e01b600052600160045260246000fd5b600082516200064b818460208701620004ec565b9190910192915050565b602081526000825180602084015262000676816040850160208701620004ec565b601f01601f19169190910160400192915050565b610891806200069a6000396000f3fe60806040526004361061005e5760003560e01c80635c60da1b116100435780635c60da1b146100a85780638f283970146100d9578063f851a440146100f95761006d565b80633659cfe6146100755780634f1ef286146100955761006d565b3661006d5761006b61010e565b005b61006b61010e565b34801561008157600080fd5b5061006b61009036600461071b565b610128565b61006b6100a3366004610736565b610165565b3480156100b457600080fd5b506100bd6101cc565b6040516001600160a01b03909116815260200160405180910390f35b3480156100e557600080fd5b5061006b6100f436600461071b565b6101fd565b34801561010557600080fd5b506100bd61021d565b610116610270565b610126610121610320565b61032a565b565b61013061034e565b6001600160a01b0316330361015d5761015a81604051806020016040528060008152506000610381565b50565b61015a61010e565b61016d61034e565b6001600160a01b031633036101c4576101bf8383838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250610381915050565b505050565b6101bf61010e565b60006101d661034e565b6001600160a01b031633036101f2576101ed610320565b905090565b6101fa61010e565b90565b61020561034e565b6001600160a01b0316330361015d5761015a816103ac565b600061022761034e565b6001600160a01b031633036101f2576101ed61034e565b6060610263838360405180606001604052806027815260200161083560279139610400565b9392505050565b3b151590565b61027861034e565b6001600160a01b031633036101265760405162461bcd60e51b815260206004820152604260248201527f5472616e73706172656e745570677261646561626c6550726f78793a2061646d60448201527f696e2063616e6e6f742066616c6c6261636b20746f2070726f7879207461726760648201527f6574000000000000000000000000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b60006101ed6104eb565b3660008037600080366000845af43d6000803e808015610349573d6000f35b3d6000fd5b60007fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b546001600160a01b0316919050565b61038a83610513565b6000825111806103975750805b156101bf576103a6838361023e565b50505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6103d561034e565b604080516001600160a01b03928316815291841660208301520160405180910390a161015a81610553565b6060833b6104765760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f60448201527f6e747261637400000000000000000000000000000000000000000000000000006064820152608401610317565b600080856001600160a01b03168560405161049191906107e5565b600060405180830381855af49150503d80600081146104cc576040519150601f19603f3d011682016040523d82523d6000602084013e6104d1565b606091505b50915091506104e182828661062b565b9695505050505050565b60007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610372565b61051c81610664565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6001600160a01b0381166105cf5760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610317565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b039290921691909117905550565b6060831561063a575081610263565b82511561064a5782518084602001fd5b8160405162461bcd60e51b81526004016103179190610801565b803b6106d85760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e7472616374000000000000000000000000000000000000006064820152608401610317565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6105f2565b80356001600160a01b038116811461071657600080fd5b919050565b60006020828403121561072d57600080fd5b610263826106ff565b60008060006040848603121561074b57600080fd5b610754846106ff565b9250602084013567ffffffffffffffff8082111561077157600080fd5b818601915086601f83011261078557600080fd5b81358181111561079457600080fd5b8760208285010111156107a657600080fd5b6020830194508093505050509250925092565b60005b838110156107d45781810151838201526020016107bc565b838111156103a65750506000910152565b600082516107f78184602087016107b9565b9190910192915050565b60208152600082518060208401526108208160408501602087016107b9565b601f01601f1916919091016040019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220fdab69a43392777d3611f0f3ede32ba76ab72fb8da12c6f315dcc7607006b40264736f6c634300080f0033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
}

var VRFCoordinatorProxyABI = VRFCoordinatorProxyMetaData.ABI

var VRFCoordinatorProxyBin = VRFCoordinatorProxyMetaData.Bin

func DeployVRFCoordinatorProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, _admin common.Address, _data []byte) (common.Address, *types.Transaction, *VRFCoordinatorProxy, error) {
	parsed, err := VRFCoordinatorProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorProxyBin), backend, _logic, _admin, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorProxy{VRFCoordinatorProxyCaller: VRFCoordinatorProxyCaller{contract: contract}, VRFCoordinatorProxyTransactor: VRFCoordinatorProxyTransactor{contract: contract}, VRFCoordinatorProxyFilterer: VRFCoordinatorProxyFilterer{contract: contract}}, nil
}

type VRFCoordinatorProxy struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorProxyCaller
	VRFCoordinatorProxyTransactor
	VRFCoordinatorProxyFilterer
}

type VRFCoordinatorProxyCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorProxyTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorProxyFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorProxySession struct {
	Contract     *VRFCoordinatorProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorProxyCallerSession struct {
	Contract *VRFCoordinatorProxyCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorProxyTransactorSession struct {
	Contract     *VRFCoordinatorProxyTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorProxyRaw struct {
	Contract *VRFCoordinatorProxy
}

type VRFCoordinatorProxyCallerRaw struct {
	Contract *VRFCoordinatorProxyCaller
}

type VRFCoordinatorProxyTransactorRaw struct {
	Contract *VRFCoordinatorProxyTransactor
}

func NewVRFCoordinatorProxy(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorProxy, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxy{address: address, abi: abi, VRFCoordinatorProxyCaller: VRFCoordinatorProxyCaller{contract: contract}, VRFCoordinatorProxyTransactor: VRFCoordinatorProxyTransactor{contract: contract}, VRFCoordinatorProxyFilterer: VRFCoordinatorProxyFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorProxyCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorProxyCaller, error) {
	contract, err := bindVRFCoordinatorProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxyCaller{contract: contract}, nil
}

func NewVRFCoordinatorProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorProxyTransactor, error) {
	contract, err := bindVRFCoordinatorProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxyTransactor{contract: contract}, nil
}

func NewVRFCoordinatorProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorProxyFilterer, error) {
	contract, err := bindVRFCoordinatorProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxyFilterer{contract: contract}, nil
}

func bindVRFCoordinatorProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorProxy.Contract.VRFCoordinatorProxyCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.VRFCoordinatorProxyTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.VRFCoordinatorProxyTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.Transact(opts, "admin")
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) Admin() (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Admin(&_VRFCoordinatorProxy.TransactOpts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Admin(&_VRFCoordinatorProxy.TransactOpts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.ChangeAdmin(&_VRFCoordinatorProxy.TransactOpts, newAdmin)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.ChangeAdmin(&_VRFCoordinatorProxy.TransactOpts, newAdmin)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.Transact(opts, "implementation")
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) Implementation() (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Implementation(&_VRFCoordinatorProxy.TransactOpts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Implementation(&_VRFCoordinatorProxy.TransactOpts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.UpgradeTo(&_VRFCoordinatorProxy.TransactOpts, newImplementation)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.UpgradeTo(&_VRFCoordinatorProxy.TransactOpts, newImplementation)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.UpgradeToAndCall(&_VRFCoordinatorProxy.TransactOpts, newImplementation, data)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.UpgradeToAndCall(&_VRFCoordinatorProxy.TransactOpts, newImplementation, data)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.RawTransact(opts, calldata)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Fallback(&_VRFCoordinatorProxy.TransactOpts, calldata)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Fallback(&_VRFCoordinatorProxy.TransactOpts, calldata)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorProxy.contract.RawTransact(opts, nil)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxySession) Receive() (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Receive(&_VRFCoordinatorProxy.TransactOpts)
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _VRFCoordinatorProxy.Contract.Receive(&_VRFCoordinatorProxy.TransactOpts)
}

type VRFCoordinatorProxyAdminChangedIterator struct {
	Event *VRFCoordinatorProxyAdminChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorProxyAdminChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorProxyAdminChanged)
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
		it.Event = new(VRFCoordinatorProxyAdminChanged)
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

func (it *VRFCoordinatorProxyAdminChangedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*VRFCoordinatorProxyAdminChangedIterator, error) {

	logs, sub, err := _VRFCoordinatorProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxyAdminChangedIterator{contract: _VRFCoordinatorProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorProxyAdminChanged)
				if err := _VRFCoordinatorProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) ParseAdminChanged(log types.Log) (*VRFCoordinatorProxyAdminChanged, error) {
	event := new(VRFCoordinatorProxyAdminChanged)
	if err := _VRFCoordinatorProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorProxyBeaconUpgradedIterator struct {
	Event *VRFCoordinatorProxyBeaconUpgraded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorProxyBeaconUpgradedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorProxyBeaconUpgraded)
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
		it.Event = new(VRFCoordinatorProxyBeaconUpgraded)
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

func (it *VRFCoordinatorProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*VRFCoordinatorProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _VRFCoordinatorProxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxyBeaconUpgradedIterator{contract: _VRFCoordinatorProxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _VRFCoordinatorProxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorProxyBeaconUpgraded)
				if err := _VRFCoordinatorProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) ParseBeaconUpgraded(log types.Log) (*VRFCoordinatorProxyBeaconUpgraded, error) {
	event := new(VRFCoordinatorProxyBeaconUpgraded)
	if err := _VRFCoordinatorProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorProxyUpgradedIterator struct {
	Event *VRFCoordinatorProxyUpgraded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorProxyUpgradedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorProxyUpgraded)
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
		it.Event = new(VRFCoordinatorProxyUpgraded)
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

func (it *VRFCoordinatorProxyUpgradedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VRFCoordinatorProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _VRFCoordinatorProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorProxyUpgradedIterator{contract: _VRFCoordinatorProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _VRFCoordinatorProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorProxyUpgraded)
				if err := _VRFCoordinatorProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

func (_VRFCoordinatorProxy *VRFCoordinatorProxyFilterer) ParseUpgraded(log types.Log) (*VRFCoordinatorProxyUpgraded, error) {
	event := new(VRFCoordinatorProxyUpgraded)
	if err := _VRFCoordinatorProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorProxy.abi.Events["AdminChanged"].ID:
		return _VRFCoordinatorProxy.ParseAdminChanged(log)
	case _VRFCoordinatorProxy.abi.Events["BeaconUpgraded"].ID:
		return _VRFCoordinatorProxy.ParseBeaconUpgraded(log)
	case _VRFCoordinatorProxy.abi.Events["Upgraded"].ID:
		return _VRFCoordinatorProxy.ParseUpgraded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorProxyAdminChanged) Topic() common.Hash {
	return common.HexToHash("0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f")
}

func (VRFCoordinatorProxyBeaconUpgraded) Topic() common.Hash {
	return common.HexToHash("0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e")
}

func (VRFCoordinatorProxyUpgraded) Topic() common.Hash {
	return common.HexToHash("0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b")
}

func (_VRFCoordinatorProxy *VRFCoordinatorProxy) Address() common.Address {
	return _VRFCoordinatorProxy.address
}

type VRFCoordinatorProxyInterface interface {
	Admin(opts *bind.TransactOpts) (*types.Transaction, error)

	ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Implementation(opts *bind.TransactOpts) (*types.Transaction, error)

	UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error)

	UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterAdminChanged(opts *bind.FilterOpts) (*VRFCoordinatorProxyAdminChangedIterator, error)

	WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorProxyAdminChanged) (event.Subscription, error)

	ParseAdminChanged(log types.Log) (*VRFCoordinatorProxyAdminChanged, error)

	FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*VRFCoordinatorProxyBeaconUpgradedIterator, error)

	WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error)

	ParseBeaconUpgraded(log types.Log) (*VRFCoordinatorProxyBeaconUpgraded, error)

	FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VRFCoordinatorProxyUpgradedIterator, error)

	WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorProxyUpgraded, implementation []common.Address) (event.Subscription, error)

	ParseUpgraded(log types.Log) (*VRFCoordinatorProxyUpgraded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
