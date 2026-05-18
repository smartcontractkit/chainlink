// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package configurator

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

var ConfiguratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610d21806101576000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379ba50971161005057806379ba5097146101355780638da5cb5b1461013d578063f2fde38b1461016557600080fd5b806301ffc9a714610077578063181f5a77146100e157806344a0b2ad14610120575b600080fd5b6100cc6100853660046106c9565b7fffffffff00000000000000000000000000000000000000000000000000000000167f44a0b2ad000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601281527f436f6e666967757261746f7220302e342e300000000000000000000000000000602082015290516100d89190610776565b61013361012e3660046109d8565b610178565b005b610133610289565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100d8565b610133610173366004610ab0565b610386565b85518460ff16806000036101b8576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610202576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b61020d816003610afa565b8211610265578161021f826003610afa565b61022a906001610b17565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260248201526044016101f9565b61026d61039a565b61027e8946308b8b8b8b8b8b61041d565b505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461030a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016101f9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61038e61039a565b61039781610526565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461041b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016101f9565b565b60008981526002602052604081208054909190829082906104479067ffffffffffffffff16610b2a565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055905060006104828c8c8c858d8d8d8d8d8d61061b565b90508b7fa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da8460000160089054906101000a900463ffffffff1683858d8d8d8d8d8d6040516104d899989796959493929190610bd2565b60405180910390a2505080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000004363ffffffff1602179055505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036105a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016101f9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808b8b8b8b8b8b8b8b8b8b6040516020016106419a99989796959493929190610c67565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b6000602082840312156106db57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461070b57600080fd5b9392505050565b6000815180845260005b818110156107385760208185018101518683018201520161071c565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061070b6020830184610712565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156107ff576107ff610789565b604052919050565b600067ffffffffffffffff82111561082157610821610789565b5060051b60200190565b803573ffffffffffffffffffffffffffffffffffffffff8116811461084f57600080fd5b919050565b600082601f83011261086557600080fd5b8135602061087a61087583610807565b6107b8565b82815260059290921b8401810191818101908684111561089957600080fd5b8286015b848110156108bb576108ae8161082b565b835291830191830161089d565b509695505050505050565b600082601f8301126108d757600080fd5b813560206108e761087583610807565b82815260059290921b8401810191818101908684111561090657600080fd5b8286015b848110156108bb578035835291830191830161090a565b803560ff8116811461084f57600080fd5b600082601f83011261094357600080fd5b813567ffffffffffffffff81111561095d5761095d610789565b61098e60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016107b8565b8181528460208386010111156109a357600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461084f57600080fd5b600080600080600080600060e0888a0312156109f357600080fd5b87359650602088013567ffffffffffffffff80821115610a1257600080fd5b610a1e8b838c01610854565b975060408a0135915080821115610a3457600080fd5b610a408b838c016108c6565b9650610a4e60608b01610921565b955060808a0135915080821115610a6457600080fd5b610a708b838c01610932565b9450610a7e60a08b016109c0565b935060c08a0135915080821115610a9457600080fd5b50610aa18a828b01610932565b91505092959891949750929550565b600060208284031215610ac257600080fd5b61070b8261082b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610b1157610b11610acb565b92915050565b80820180821115610b1157610b11610acb565b600067ffffffffffffffff808316818103610b4757610b47610acb565b6001019392505050565b600081518084526020808501945080840160005b83811015610b9757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610b65565b509495945050505050565b600081518084526020808501945080840160005b83811015610b9757815187529582019590820190600101610bb6565b600061012063ffffffff8c1683528a602084015267ffffffffffffffff808b166040850152816060850152610c098285018b610b51565b91508382036080850152610c1d828a610ba2565b915060ff881660a085015283820360c0850152610c3a8288610712565b90861660e08501528381036101008501529050610c578185610712565b9c9b505050505050505050505050565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b166060850152816080850152610cb48285018b610b51565b915083820360a0850152610cc8828a610ba2565b915060ff881660c085015283820360e0850152610ce58288610712565b9086166101008501528381036101208501529050610d038185610712565b9d9c5050505050505050505050505056fea164736f6c6343000813000a",
}

var ConfiguratorABI = ConfiguratorMetaData.ABI

var ConfiguratorBin = ConfiguratorMetaData.Bin

func DeployConfigurator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Configurator, error) {
	parsed, err := ConfiguratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfiguratorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Configurator{address: address, abi: *parsed, ConfiguratorCaller: ConfiguratorCaller{contract: contract}, ConfiguratorTransactor: ConfiguratorTransactor{contract: contract}, ConfiguratorFilterer: ConfiguratorFilterer{contract: contract}}, nil
}

type Configurator struct {
	address common.Address
	abi     abi.ABI
	ConfiguratorCaller
	ConfiguratorTransactor
	ConfiguratorFilterer
}

type ConfiguratorCaller struct {
	contract *bind.BoundContract
}

type ConfiguratorTransactor struct {
	contract *bind.BoundContract
}

type ConfiguratorFilterer struct {
	contract *bind.BoundContract
}

type ConfiguratorSession struct {
	Contract     *Configurator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ConfiguratorCallerSession struct {
	Contract *ConfiguratorCaller
	CallOpts bind.CallOpts
}

type ConfiguratorTransactorSession struct {
	Contract     *ConfiguratorTransactor
	TransactOpts bind.TransactOpts
}

type ConfiguratorRaw struct {
	Contract *Configurator
}

type ConfiguratorCallerRaw struct {
	Contract *ConfiguratorCaller
}

type ConfiguratorTransactorRaw struct {
	Contract *ConfiguratorTransactor
}

func NewConfigurator(address common.Address, backend bind.ContractBackend) (*Configurator, error) {
	abi, err := abi.JSON(strings.NewReader(ConfiguratorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindConfigurator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Configurator{address: address, abi: abi, ConfiguratorCaller: ConfiguratorCaller{contract: contract}, ConfiguratorTransactor: ConfiguratorTransactor{contract: contract}, ConfiguratorFilterer: ConfiguratorFilterer{contract: contract}}, nil
}

func NewConfiguratorCaller(address common.Address, caller bind.ContractCaller) (*ConfiguratorCaller, error) {
	contract, err := bindConfigurator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorCaller{contract: contract}, nil
}

func NewConfiguratorTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfiguratorTransactor, error) {
	contract, err := bindConfigurator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorTransactor{contract: contract}, nil
}

func NewConfiguratorFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfiguratorFilterer, error) {
	contract, err := bindConfigurator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorFilterer{contract: contract}, nil
}

func bindConfigurator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfiguratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Configurator *ConfiguratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Configurator.Contract.ConfiguratorCaller.contract.Call(opts, result, method, params...)
}

func (_Configurator *ConfiguratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Configurator.Contract.ConfiguratorTransactor.contract.Transfer(opts)
}

func (_Configurator *ConfiguratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Configurator.Contract.ConfiguratorTransactor.contract.Transact(opts, method, params...)
}

func (_Configurator *ConfiguratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Configurator.Contract.contract.Call(opts, result, method, params...)
}

func (_Configurator *ConfiguratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Configurator.Contract.contract.Transfer(opts)
}

func (_Configurator *ConfiguratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Configurator.Contract.contract.Transact(opts, method, params...)
}

func (_Configurator *ConfiguratorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Configurator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Configurator *ConfiguratorSession) Owner() (common.Address, error) {
	return _Configurator.Contract.Owner(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorCallerSession) Owner() (common.Address, error) {
	return _Configurator.Contract.Owner(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Configurator.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Configurator *ConfiguratorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Configurator.Contract.SupportsInterface(&_Configurator.CallOpts, interfaceId)
}

func (_Configurator *ConfiguratorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Configurator.Contract.SupportsInterface(&_Configurator.CallOpts, interfaceId)
}

func (_Configurator *ConfiguratorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Configurator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Configurator *ConfiguratorSession) TypeAndVersion() (string, error) {
	return _Configurator.Contract.TypeAndVersion(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorCallerSession) TypeAndVersion() (string, error) {
	return _Configurator.Contract.TypeAndVersion(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "acceptOwnership")
}

func (_Configurator *ConfiguratorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Configurator.Contract.AcceptOwnership(&_Configurator.TransactOpts)
}

func (_Configurator *ConfiguratorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Configurator.Contract.AcceptOwnership(&_Configurator.TransactOpts)
}

func (_Configurator *ConfiguratorTransactor) SetConfig(opts *bind.TransactOpts, donId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "setConfig", donId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorSession) SetConfig(donId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.Contract.SetConfig(&_Configurator.TransactOpts, donId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorTransactorSession) SetConfig(donId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.Contract.SetConfig(&_Configurator.TransactOpts, donId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "transferOwnership", to)
}

func (_Configurator *ConfiguratorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Configurator.Contract.TransferOwnership(&_Configurator.TransactOpts, to)
}

func (_Configurator *ConfiguratorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Configurator.Contract.TransferOwnership(&_Configurator.TransactOpts, to)
}

type ConfiguratorConfigSetIterator struct {
	Event *ConfiguratorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorConfigSet)
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
		it.Event = new(ConfiguratorConfigSet)
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

func (it *ConfiguratorConfigSetIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorConfigSet struct {
	ConfigId                  [32]byte
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	OffchainTransmitters      [][32]byte
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ConfiguratorConfigSetIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "ConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorConfigSetIterator{contract: _Configurator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *ConfiguratorConfigSet, configId [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "ConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorConfigSet)
				if err := _Configurator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseConfigSet(log types.Log) (*ConfiguratorConfigSet, error) {
	event := new(ConfiguratorConfigSet)
	if err := _Configurator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConfiguratorOwnershipTransferRequestedIterator struct {
	Event *ConfiguratorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorOwnershipTransferRequested)
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
		it.Event = new(ConfiguratorOwnershipTransferRequested)
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

func (it *ConfiguratorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorOwnershipTransferRequestedIterator{contract: _Configurator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorOwnershipTransferRequested)
				if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfiguratorOwnershipTransferRequested, error) {
	event := new(ConfiguratorOwnershipTransferRequested)
	if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConfiguratorOwnershipTransferredIterator struct {
	Event *ConfiguratorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorOwnershipTransferred)
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
		it.Event = new(ConfiguratorOwnershipTransferred)
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

func (it *ConfiguratorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorOwnershipTransferredIterator{contract: _Configurator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorOwnershipTransferred)
				if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseOwnershipTransferred(log types.Log) (*ConfiguratorOwnershipTransferred, error) {
	event := new(ConfiguratorOwnershipTransferred)
	if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Configurator *Configurator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Configurator.abi.Events["ConfigSet"].ID:
		return _Configurator.ParseConfigSet(log)
	case _Configurator.abi.Events["OwnershipTransferRequested"].ID:
		return _Configurator.ParseOwnershipTransferRequested(log)
	case _Configurator.abi.Events["OwnershipTransferred"].ID:
		return _Configurator.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ConfiguratorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da")
}

func (ConfiguratorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ConfiguratorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_Configurator *Configurator) Address() common.Address {
	return _Configurator.address
}

type ConfiguratorInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, donId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ConfiguratorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *ConfiguratorConfigSet, configId [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*ConfiguratorConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ConfiguratorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ConfiguratorOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
