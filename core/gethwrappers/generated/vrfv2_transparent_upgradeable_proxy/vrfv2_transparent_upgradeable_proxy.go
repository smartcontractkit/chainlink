// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2_transparent_upgradeable_proxy

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

var VRFV2TransparentUpgradeableProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"implementation_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526040516200113b3803806200113b8339810160408190526200002691620004c8565b82828282816200005860017f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbd620005fb565b600080516020620010f48339815191521462000078576200007862000650565b6200008682826000620000ed565b50620000b6905060017fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6104620005fb565b600080516020620010d483398151915214620000d657620000d662000650565b620000e1826200012a565b5050505050506200067c565b620000f88362000185565b600082511180620001065750805b156200012557620001238383620001c760201b620002ff1760201c565b505b505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f62000155620001f6565b604080516001600160a01b03928316815291841660208301520160405180910390a162000182816200022f565b50565b6200019081620002e4565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060620001ef8383604051806060016040528060278152602001620011146027913962000387565b9392505050565b600062000220600080516020620010d483398151915260001b6200046460201b620002731760201c565b546001600160a01b0316919050565b6001600160a01b0381166200029a5760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b80620002c3600080516020620010d483398151915260001b6200046460201b620002731760201c565b80546001600160a01b0319166001600160a01b039290921691909117905550565b620002fa816200046760201b6200032b1760201c565b6200035e5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840162000291565b80620002c3600080516020620010f483398151915260001b6200046460201b620002731760201c565b6060833b620003e85760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b606482015260840162000291565b600080856001600160a01b031685604051620004059190620005a8565b600060405180830381855af49150503d806000811462000442576040519150601f19603f3d011682016040523d82523d6000602084013e62000447565b606091505b5090925090506200045a8282866200046d565b9695505050505050565b90565b3b151590565b606083156200047e575081620001ef565b8251156200048f5782518084602001fd5b8160405162461bcd60e51b8152600401620002919190620005c6565b80516001600160a01b0381168114620004c357600080fd5b919050565b600080600060608486031215620004de57600080fd5b620004e984620004ab565b9250620004f960208501620004ab565b60408501519092506001600160401b03808211156200051757600080fd5b818601915086601f8301126200052c57600080fd5b81518181111562000541576200054162000666565b604051601f8201601f19908116603f011681019083821181831017156200056c576200056c62000666565b816040528281528960208487010111156200058657600080fd5b6200059983602083016020880162000621565b80955050505050509250925092565b60008251620005bc81846020870162000621565b9190910192915050565b6020815260008251806020840152620005e781604085016020870162000621565b601f01601f19169190910160400192915050565b6000828210156200061c57634e487b7160e01b600052601160045260246000fd5b500390565b60005b838110156200063e57818101518382015260200162000624565b83811115620001235750506000910152565b634e487b7160e01b600052600160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b610a48806200068c6000396000f3fe60806040526004361061005e5760003560e01c80635c60da1b116100435780635c60da1b146100a85780638f283970146100e6578063f851a440146101065761006d565b80633659cfe6146100755780634f1ef286146100955761006d565b3661006d5761006b61011b565b005b61006b61011b565b34801561008157600080fd5b5061006b6100903660046108dd565b610135565b61006b6100a33660046108f8565b610196565b3480156100b457600080fd5b506100bd610221565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b3480156100f257600080fd5b5061006b6101013660046108dd565b610276565b34801561011257600080fd5b506100bd6102ba565b610123610331565b61013361012e61041f565b610429565b565b61013d61044d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561018e5761018b8160405180602001604052806000815250600061048d565b50565b61018b61011b565b61019e61044d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610219576102148383838080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506001925061048d915050565b505050565b61021461011b565b600061022b61044d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561026b5761026661041f565b905090565b61027361011b565b90565b61027e61044d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561018e5761018b816104b8565b60006102c461044d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561026b5761026661044d565b60606103248383604051806060016040528060278152602001610a1560279139610519565b9392505050565b3b151590565b61033961044d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610133576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604260248201527f5472616e73706172656e745570677261646561626c6550726f78793a2061646d60448201527f696e2063616e6e6f742066616c6c6261636b20746f2070726f7879207461726760648201527f6574000000000000000000000000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b600061026661062b565b3660008037600080366000845af43d6000803e808015610448573d6000f35b3d6000fd5b60007fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b5473ffffffffffffffffffffffffffffffffffffffff16919050565b61049683610653565b6000825111806104a35750805b15610214576104b283836102ff565b50505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6104e161044d565b6040805173ffffffffffffffffffffffffffffffffffffffff928316815291841660208301520160405180910390a161018b816106a0565b6060833b6105a9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f60448201527f6e747261637400000000000000000000000000000000000000000000000000006064820152608401610416565b6000808573ffffffffffffffffffffffffffffffffffffffff16856040516105d1919061097b565b600060405180830381855af49150503d806000811461060c576040519150601f19603f3d011682016040523d82523d6000602084013e610611565b606091505b50915091506106218282866107ac565b9695505050505050565b60007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610471565b61065c816107ff565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b73ffffffffffffffffffffffffffffffffffffffff8116610743576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610416565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905550565b606083156107bb575081610324565b8251156107cb5782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104169190610997565b803b61088d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e7472616374000000000000000000000000000000000000006064820152608401610416565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610766565b803573ffffffffffffffffffffffffffffffffffffffff811681146108d857600080fd5b919050565b6000602082840312156108ef57600080fd5b610324826108b4565b60008060006040848603121561090d57600080fd5b610916846108b4565b9250602084013567ffffffffffffffff8082111561093357600080fd5b818601915086601f83011261094757600080fd5b81358181111561095657600080fd5b87602082850101111561096857600080fd5b6020830194508093505050509250925092565b6000825161098d8184602087016109e8565b9190910192915050565b60208152600082518060208401526109b68160408501602087016109e8565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60005b83811015610a035781810151838201526020016109eb565b838111156104b2575050600091015256fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a164736f6c6343000806000ab53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
}

var VRFV2TransparentUpgradeableProxyABI = VRFV2TransparentUpgradeableProxyMetaData.ABI

var VRFV2TransparentUpgradeableProxyBin = VRFV2TransparentUpgradeableProxyMetaData.Bin

func DeployVRFV2TransparentUpgradeableProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, admin_ common.Address, _data []byte) (common.Address, *types.Transaction, *VRFV2TransparentUpgradeableProxy, error) {
	parsed, err := VRFV2TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2TransparentUpgradeableProxyBin), backend, _logic, admin_, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2TransparentUpgradeableProxy{VRFV2TransparentUpgradeableProxyCaller: VRFV2TransparentUpgradeableProxyCaller{contract: contract}, VRFV2TransparentUpgradeableProxyTransactor: VRFV2TransparentUpgradeableProxyTransactor{contract: contract}, VRFV2TransparentUpgradeableProxyFilterer: VRFV2TransparentUpgradeableProxyFilterer{contract: contract}}, nil
}

type VRFV2TransparentUpgradeableProxy struct {
	address common.Address
	abi     abi.ABI
	VRFV2TransparentUpgradeableProxyCaller
	VRFV2TransparentUpgradeableProxyTransactor
	VRFV2TransparentUpgradeableProxyFilterer
}

type VRFV2TransparentUpgradeableProxyCaller struct {
	contract *bind.BoundContract
}

type VRFV2TransparentUpgradeableProxyTransactor struct {
	contract *bind.BoundContract
}

type VRFV2TransparentUpgradeableProxyFilterer struct {
	contract *bind.BoundContract
}

type VRFV2TransparentUpgradeableProxySession struct {
	Contract     *VRFV2TransparentUpgradeableProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2TransparentUpgradeableProxyCallerSession struct {
	Contract *VRFV2TransparentUpgradeableProxyCaller
	CallOpts bind.CallOpts
}

type VRFV2TransparentUpgradeableProxyTransactorSession struct {
	Contract     *VRFV2TransparentUpgradeableProxyTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2TransparentUpgradeableProxyRaw struct {
	Contract *VRFV2TransparentUpgradeableProxy
}

type VRFV2TransparentUpgradeableProxyCallerRaw struct {
	Contract *VRFV2TransparentUpgradeableProxyCaller
}

type VRFV2TransparentUpgradeableProxyTransactorRaw struct {
	Contract *VRFV2TransparentUpgradeableProxyTransactor
}

func NewVRFV2TransparentUpgradeableProxy(address common.Address, backend bind.ContractBackend) (*VRFV2TransparentUpgradeableProxy, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2TransparentUpgradeableProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2TransparentUpgradeableProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxy{address: address, abi: abi, VRFV2TransparentUpgradeableProxyCaller: VRFV2TransparentUpgradeableProxyCaller{contract: contract}, VRFV2TransparentUpgradeableProxyTransactor: VRFV2TransparentUpgradeableProxyTransactor{contract: contract}, VRFV2TransparentUpgradeableProxyFilterer: VRFV2TransparentUpgradeableProxyFilterer{contract: contract}}, nil
}

func NewVRFV2TransparentUpgradeableProxyCaller(address common.Address, caller bind.ContractCaller) (*VRFV2TransparentUpgradeableProxyCaller, error) {
	contract, err := bindVRFV2TransparentUpgradeableProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxyCaller{contract: contract}, nil
}

func NewVRFV2TransparentUpgradeableProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2TransparentUpgradeableProxyTransactor, error) {
	contract, err := bindVRFV2TransparentUpgradeableProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxyTransactor{contract: contract}, nil
}

func NewVRFV2TransparentUpgradeableProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2TransparentUpgradeableProxyFilterer, error) {
	contract, err := bindVRFV2TransparentUpgradeableProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxyFilterer{contract: contract}, nil
}

func bindVRFV2TransparentUpgradeableProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFV2TransparentUpgradeableProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2TransparentUpgradeableProxy.Contract.VRFV2TransparentUpgradeableProxyCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.VRFV2TransparentUpgradeableProxyTransactor.contract.Transfer(opts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.VRFV2TransparentUpgradeableProxyTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2TransparentUpgradeableProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.contract.Transfer(opts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.Transact(opts, "admin")
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) Admin() (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Admin(&_VRFV2TransparentUpgradeableProxy.TransactOpts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Admin(&_VRFV2TransparentUpgradeableProxy.TransactOpts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.ChangeAdmin(&_VRFV2TransparentUpgradeableProxy.TransactOpts, newAdmin)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.ChangeAdmin(&_VRFV2TransparentUpgradeableProxy.TransactOpts, newAdmin)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.Transact(opts, "implementation")
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) Implementation() (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Implementation(&_VRFV2TransparentUpgradeableProxy.TransactOpts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Implementation(&_VRFV2TransparentUpgradeableProxy.TransactOpts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.UpgradeTo(&_VRFV2TransparentUpgradeableProxy.TransactOpts, newImplementation)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.UpgradeTo(&_VRFV2TransparentUpgradeableProxy.TransactOpts, newImplementation)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.UpgradeToAndCall(&_VRFV2TransparentUpgradeableProxy.TransactOpts, newImplementation, data)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.UpgradeToAndCall(&_VRFV2TransparentUpgradeableProxy.TransactOpts, newImplementation, data)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.RawTransact(opts, calldata)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Fallback(&_VRFV2TransparentUpgradeableProxy.TransactOpts, calldata)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Fallback(&_VRFV2TransparentUpgradeableProxy.TransactOpts, calldata)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.contract.RawTransact(opts, nil)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxySession) Receive() (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Receive(&_VRFV2TransparentUpgradeableProxy.TransactOpts)
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _VRFV2TransparentUpgradeableProxy.Contract.Receive(&_VRFV2TransparentUpgradeableProxy.TransactOpts)
}

type VRFV2TransparentUpgradeableProxyAdminChangedIterator struct {
	Event *VRFV2TransparentUpgradeableProxyAdminChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2TransparentUpgradeableProxyAdminChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2TransparentUpgradeableProxyAdminChanged)
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
		it.Event = new(VRFV2TransparentUpgradeableProxyAdminChanged)
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

func (it *VRFV2TransparentUpgradeableProxyAdminChangedIterator) Error() error {
	return it.fail
}

func (it *VRFV2TransparentUpgradeableProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2TransparentUpgradeableProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*VRFV2TransparentUpgradeableProxyAdminChangedIterator, error) {

	logs, sub, err := _VRFV2TransparentUpgradeableProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxyAdminChangedIterator{contract: _VRFV2TransparentUpgradeableProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFV2TransparentUpgradeableProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _VRFV2TransparentUpgradeableProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2TransparentUpgradeableProxyAdminChanged)
				if err := _VRFV2TransparentUpgradeableProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) ParseAdminChanged(log types.Log) (*VRFV2TransparentUpgradeableProxyAdminChanged, error) {
	event := new(VRFV2TransparentUpgradeableProxyAdminChanged)
	if err := _VRFV2TransparentUpgradeableProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator struct {
	Event *VRFV2TransparentUpgradeableProxyBeaconUpgraded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2TransparentUpgradeableProxyBeaconUpgraded)
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
		it.Event = new(VRFV2TransparentUpgradeableProxyBeaconUpgraded)
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

func (it *VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

func (it *VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2TransparentUpgradeableProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _VRFV2TransparentUpgradeableProxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator{contract: _VRFV2TransparentUpgradeableProxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *VRFV2TransparentUpgradeableProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _VRFV2TransparentUpgradeableProxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2TransparentUpgradeableProxyBeaconUpgraded)
				if err := _VRFV2TransparentUpgradeableProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) ParseBeaconUpgraded(log types.Log) (*VRFV2TransparentUpgradeableProxyBeaconUpgraded, error) {
	event := new(VRFV2TransparentUpgradeableProxyBeaconUpgraded)
	if err := _VRFV2TransparentUpgradeableProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2TransparentUpgradeableProxyUpgradedIterator struct {
	Event *VRFV2TransparentUpgradeableProxyUpgraded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2TransparentUpgradeableProxyUpgradedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2TransparentUpgradeableProxyUpgraded)
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
		it.Event = new(VRFV2TransparentUpgradeableProxyUpgraded)
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

func (it *VRFV2TransparentUpgradeableProxyUpgradedIterator) Error() error {
	return it.fail
}

func (it *VRFV2TransparentUpgradeableProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2TransparentUpgradeableProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VRFV2TransparentUpgradeableProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _VRFV2TransparentUpgradeableProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2TransparentUpgradeableProxyUpgradedIterator{contract: _VRFV2TransparentUpgradeableProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VRFV2TransparentUpgradeableProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _VRFV2TransparentUpgradeableProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2TransparentUpgradeableProxyUpgraded)
				if err := _VRFV2TransparentUpgradeableProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxyFilterer) ParseUpgraded(log types.Log) (*VRFV2TransparentUpgradeableProxyUpgraded, error) {
	event := new(VRFV2TransparentUpgradeableProxyUpgraded)
	if err := _VRFV2TransparentUpgradeableProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2TransparentUpgradeableProxy.abi.Events["AdminChanged"].ID:
		return _VRFV2TransparentUpgradeableProxy.ParseAdminChanged(log)
	case _VRFV2TransparentUpgradeableProxy.abi.Events["BeaconUpgraded"].ID:
		return _VRFV2TransparentUpgradeableProxy.ParseBeaconUpgraded(log)
	case _VRFV2TransparentUpgradeableProxy.abi.Events["Upgraded"].ID:
		return _VRFV2TransparentUpgradeableProxy.ParseUpgraded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2TransparentUpgradeableProxyAdminChanged) Topic() common.Hash {
	return common.HexToHash("0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f")
}

func (VRFV2TransparentUpgradeableProxyBeaconUpgraded) Topic() common.Hash {
	return common.HexToHash("0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e")
}

func (VRFV2TransparentUpgradeableProxyUpgraded) Topic() common.Hash {
	return common.HexToHash("0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b")
}

func (_VRFV2TransparentUpgradeableProxy *VRFV2TransparentUpgradeableProxy) Address() common.Address {
	return _VRFV2TransparentUpgradeableProxy.address
}

type VRFV2TransparentUpgradeableProxyInterface interface {
	Admin(opts *bind.TransactOpts) (*types.Transaction, error)

	ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Implementation(opts *bind.TransactOpts) (*types.Transaction, error)

	UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error)

	UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterAdminChanged(opts *bind.FilterOpts) (*VRFV2TransparentUpgradeableProxyAdminChangedIterator, error)

	WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFV2TransparentUpgradeableProxyAdminChanged) (event.Subscription, error)

	ParseAdminChanged(log types.Log) (*VRFV2TransparentUpgradeableProxyAdminChanged, error)

	FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*VRFV2TransparentUpgradeableProxyBeaconUpgradedIterator, error)

	WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *VRFV2TransparentUpgradeableProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error)

	ParseBeaconUpgraded(log types.Log) (*VRFV2TransparentUpgradeableProxyBeaconUpgraded, error)

	FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VRFV2TransparentUpgradeableProxyUpgradedIterator, error)

	WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VRFV2TransparentUpgradeableProxyUpgraded, implementation []common.Address) (event.Subscription, error)

	ParseUpgraded(log types.Log) (*VRFV2TransparentUpgradeableProxyUpgraded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
