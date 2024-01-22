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

var VRFV2TransparentUpgradeableProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052604051620011863803806200118683398101604081905262000026916200045e565b8282828281620000398282600062000053565b506200004790508262000090565b505050505050620005d6565b6200005e83620000eb565b6000825111806200006c5750805b156200008b576200008983836200012d60201b620002bd1760201c565b505b505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620000bb6200015c565b604080516001600160a01b03928316815291841660208301520160405180910390a1620000e88162000195565b50565b620000f6816200024a565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606200015583836040518060600160405280602781526020016200115f60279139620002fe565b9392505050565b6000620001866000805160206200113f83398151915260001b6200037d60201b620002e91760201c565b546001600160a01b0316919050565b6001600160a01b038116620002005760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b80620002296000805160206200113f83398151915260001b6200037d60201b620002e91760201c565b80546001600160a01b0319166001600160a01b039290921691909117905550565b62000260816200038060201b620002ec1760201c565b620002c45760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401620001f7565b80620002297f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc60001b6200037d60201b620002e91760201c565b6060600080856001600160a01b0316856040516200031d91906200053e565b600060405180830381855af49150503d80600081146200035a576040519150601f19603f3d011682016040523d82523d6000602084013e6200035f565b606091505b50909250905062000373868383876200038f565b9695505050505050565b90565b6001600160a01b03163b151590565b6060831562000400578251620003f8576001600160a01b0385163b620003f85760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401620001f7565b50816200040c565b6200040c838362000414565b949350505050565b815115620004255781518083602001fd5b8060405162461bcd60e51b8152600401620001f791906200055c565b80516001600160a01b03811681146200045957600080fd5b919050565b6000806000606084860312156200047457600080fd5b6200047f8462000441565b92506200048f6020850162000441565b60408501519092506001600160401b0380821115620004ad57600080fd5b818601915086601f830112620004c257600080fd5b815181811115620004d757620004d7620005c0565b604051601f8201601f19908116603f01168101908382118183101715620005025762000502620005c0565b816040528281528960208487010111156200051c57600080fd5b6200052f83602083016020880162000591565b80955050505050509250925092565b600082516200055281846020870162000591565b9190910192915050565b60208152600082518060208401526200057d81604085016020870162000591565b601f01601f19169190910160400192915050565b60005b83811015620005ae57818101518382015260200162000594565b83811115620000895750506000910152565b634e487b7160e01b600052604160045260246000fd5b610b5980620005e66000396000f3fe60806040523661001357610011610017565b005b6100115b61001f610308565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156102b35760607fffffffff00000000000000000000000000000000000000000000000000000000600035167f3659cfe6000000000000000000000000000000000000000000000000000000008114156100b0576100a9610348565b91506102ab565b7fffffffff0000000000000000000000000000000000000000000000000000000081167f4f1ef286000000000000000000000000000000000000000000000000000000001415610102576100a961039f565b7fffffffff0000000000000000000000000000000000000000000000000000000081167f8f283970000000000000000000000000000000000000000000000000000000001415610154576100a96103e5565b7fffffffff0000000000000000000000000000000000000000000000000000000081167ff851a4400000000000000000000000000000000000000000000000000000000014156101a6576100a9610416565b7fffffffff0000000000000000000000000000000000000000000000000000000081167f5c60da1b0000000000000000000000000000000000000000000000000000000014156101f8576100a9610463565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604260248201527f5472616e73706172656e745570677261646561626c6550726f78793a2061646d60448201527f696e2063616e6e6f742066616c6c6261636b20746f2070726f7879207461726760648201527f6574000000000000000000000000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b815160208301f35b6102bb610477565b565b60606102e28383604051806060016040528060278152602001610b2660279139610487565b9392505050565b90565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b60007fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b5473ffffffffffffffffffffffffffffffffffffffff16919050565b606061035261050c565b60006103613660048184610aa0565b81019061036e9190610938565b905061038b81604051806020016040528060008152506000610517565b505060408051602081019091526000815290565b60606000806103b13660048184610aa0565b8101906103be9190610953565b915091506103ce82826001610517565b604051806020016040528060008152509250505090565b60606103ef61050c565b60006103fe3660048184610aa0565b81019061040b9190610938565b905061038b81610543565b606061042061050c565b600061042a610308565b6040805173ffffffffffffffffffffffffffffffffffffffff831660208201529192500160405160208183030381529060405291505090565b606061046d61050c565b600061042a6105a7565b6102bb6104826105a7565b6105b6565b60606000808573ffffffffffffffffffffffffffffffffffffffff16856040516104b19190610a33565b600060405180830381855af49150503d80600081146104ec576040519150601f19603f3d011682016040523d82523d6000602084013e6104f1565b606091505b5091509150610502868383876105da565b9695505050505050565b34156102bb57600080fd5b6105208361067f565b60008251118061052d5750805b1561053e5761053c83836102bd565b505b505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f61056c610308565b6040805173ffffffffffffffffffffffffffffffffffffffff928316815291841660208301520160405180910390a16105a4816106cc565b50565b60006105b16107d8565b905090565b3660008037600080366000845af43d6000803e8080156105d5573d6000f35b3d6000fd5b6060831561066d5782516106665773ffffffffffffffffffffffffffffffffffffffff85163b610666576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016102a2565b5081610677565b6106778383610800565b949350505050565b61068881610844565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b73ffffffffffffffffffffffffffffffffffffffff811661076f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016102a2565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905550565b60007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc61032c565b8151156108105781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a29190610a4f565b73ffffffffffffffffffffffffffffffffffffffff81163b6108e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e74726163740000000000000000000000000000000000000060648201526084016102a2565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610792565b803573ffffffffffffffffffffffffffffffffffffffff8116811461093357600080fd5b919050565b60006020828403121561094a57600080fd5b6102e28261090f565b6000806040838503121561096657600080fd5b61096f8361090f565b9150602083013567ffffffffffffffff8082111561098c57600080fd5b818501915085601f8301126109a057600080fd5b8135818111156109b2576109b2610af6565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156109f8576109f8610af6565b81604052828152886020848701011115610a1157600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60008251610a45818460208701610aca565b9190910192915050565b6020815260008251806020840152610a6e816040850160208701610aca565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60008085851115610ab057600080fd5b83861115610abd57600080fd5b5050820193919092039150565b60005b83811015610ae5578181015183820152602001610acd565b8381111561053c5750506000910152565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a164736f6c6343000806000ab53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
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
	return address, tx, &VRFV2TransparentUpgradeableProxy{address: address, abi: *parsed, VRFV2TransparentUpgradeableProxyCaller: VRFV2TransparentUpgradeableProxyCaller{contract: contract}, VRFV2TransparentUpgradeableProxyTransactor: VRFV2TransparentUpgradeableProxyTransactor{contract: contract}, VRFV2TransparentUpgradeableProxyFilterer: VRFV2TransparentUpgradeableProxyFilterer{contract: contract}}, nil
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
	parsed, err := VRFV2TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
