// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package errored_verifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

type CommonAddressAndWeight struct {
	Addr   common.Address
	Weight uint64
}

var ErroredVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"activateConfig\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"activateFeed\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"deactivateConfig\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"deactivateFeed\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"latestConfigDigestAndEpoch\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structCommon.AddressAndWeight[]\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"weight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setConfigFromSource\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structCommon.AddressAndWeight[]\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"weight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verify\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"FailedToActivateConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToActivateFeed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToDeactivateConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToDeactivateFeed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToGetLatestConfigDetails\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToGetLatestConfigDigestAndEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToSetConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerify\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b50610a58806100206000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063b70d929d11610076578063e7db9c2a1161005b578063e7db9c2a146101d1578063e84f128e146101e4578063f01072211461021a57600080fd5b8063b70d929d14610188578063ded6307c146101be57600080fd5b80633dd86430116100a75780633dd864301461014d578063564a0a7a1461016257806394d959801461017557600080fd5b806301ffc9a7146100c35780633d3ac1b51461012d575b600080fd5b6101186100d13660046103c4565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61014061013b36600461056b565b610228565b60405161012491906105b9565b61016061015b366004610625565b61025c565b005b610160610170366004610625565b61028e565b61016061018336600461063e565b6102c0565b61019b610196366004610625565b6102f2565b604080519315158452602084019290925263ffffffff1690820152606001610124565b6101606101cc36600461063e565b610329565b6101606101df36600461081b565b61035b565b6101f76101f2366004610625565b61038d565b6040805163ffffffff948516815293909216602084015290820152606001610124565b6101606101df36600461094e565b60606040517fcf2e344600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f9601b68300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fa03564b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f8a406e4600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060006040517fbbc0083000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f7adb7c9600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f35e91bf100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060006040517fa06d64a000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000602082840312156103d657600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461040657600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561045f5761045f61040d565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156104ac576104ac61040d565b604052919050565b600082601f8301126104c557600080fd5b813567ffffffffffffffff8111156104df576104df61040d565b61051060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610465565b81815284602083860101111561052557600080fd5b816020850160208301376000918101602001919091529392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461056657600080fd5b919050565b6000806040838503121561057e57600080fd5b823567ffffffffffffffff81111561059557600080fd5b6105a1858286016104b4565b9250506105b060208401610542565b90509250929050565b600060208083528351808285015260005b818110156105e6578581018301518582016040015282016105ca565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561063757600080fd5b5035919050565b6000806040838503121561065157600080fd5b50508035926020909101359150565b803563ffffffff8116811461056657600080fd5b600067ffffffffffffffff82111561068e5761068e61040d565b5060051b60200190565b600082601f8301126106a957600080fd5b813560206106be6106b983610674565b610465565b82815260059290921b840181019181810190868411156106dd57600080fd5b8286015b848110156106ff576106f281610542565b83529183019183016106e1565b509695505050505050565b600082601f83011261071b57600080fd5b8135602061072b6106b983610674565b82815260059290921b8401810191818101908684111561074a57600080fd5b8286015b848110156106ff578035835291830191830161074e565b803560ff8116811461056657600080fd5b803567ffffffffffffffff8116811461056657600080fd5b600082601f83011261079f57600080fd5b813560206107af6106b983610674565b82815260069290921b840181019181810190868411156107ce57600080fd5b8286015b848110156106ff57604081890312156107eb5760008081fd5b6107f361043c565b6107fc82610542565b8152610809858301610776565b818601528352918301916040016107d2565b60008060008060008060008060008060006101608c8e03121561083d57600080fd5b8b359a5060208c0135995061085460408d01610542565b985061086260608d01610660565b975067ffffffffffffffff8060808e0135111561087e57600080fd5b61088e8e60808f01358f01610698565b97508060a08e013511156108a157600080fd5b6108b18e60a08f01358f0161070a565b96506108bf60c08e01610765565b95508060e08e013511156108d257600080fd5b6108e28e60e08f01358f016104b4565b94506108f16101008e01610776565b9350806101208e0135111561090557600080fd5b6109168e6101208f01358f016104b4565b9250806101408e0135111561092a57600080fd5b5061093c8d6101408e01358e0161078e565b90509295989b509295989b9093969950565b600080600080600080600080610100898b03121561096b57600080fd5b88359750602089013567ffffffffffffffff8082111561098a57600080fd5b6109968c838d01610698565b985060408b01359150808211156109ac57600080fd5b6109b88c838d0161070a565b97506109c660608c01610765565b965060808b01359150808211156109dc57600080fd5b6109e88c838d016104b4565b95506109f660a08c01610776565b945060c08b0135915080821115610a0c57600080fd5b610a188c838d016104b4565b935060e08b0135915080821115610a2e57600080fd5b50610a3b8b828c0161078e565b915050929598509295989093965056fea164736f6c6343000813000a",
}

var ErroredVerifierABI = ErroredVerifierMetaData.ABI

var ErroredVerifierBin = ErroredVerifierMetaData.Bin

func DeployErroredVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ErroredVerifier, error) {
	parsed, err := ErroredVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ErroredVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ErroredVerifier{address: address, abi: *parsed, ErroredVerifierCaller: ErroredVerifierCaller{contract: contract}, ErroredVerifierTransactor: ErroredVerifierTransactor{contract: contract}, ErroredVerifierFilterer: ErroredVerifierFilterer{contract: contract}}, nil
}

type ErroredVerifier struct {
	address common.Address
	abi     abi.ABI
	ErroredVerifierCaller
	ErroredVerifierTransactor
	ErroredVerifierFilterer
}

type ErroredVerifierCaller struct {
	contract *bind.BoundContract
}

type ErroredVerifierTransactor struct {
	contract *bind.BoundContract
}

type ErroredVerifierFilterer struct {
	contract *bind.BoundContract
}

type ErroredVerifierSession struct {
	Contract     *ErroredVerifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ErroredVerifierCallerSession struct {
	Contract *ErroredVerifierCaller
	CallOpts bind.CallOpts
}

type ErroredVerifierTransactorSession struct {
	Contract     *ErroredVerifierTransactor
	TransactOpts bind.TransactOpts
}

type ErroredVerifierRaw struct {
	Contract *ErroredVerifier
}

type ErroredVerifierCallerRaw struct {
	Contract *ErroredVerifierCaller
}

type ErroredVerifierTransactorRaw struct {
	Contract *ErroredVerifierTransactor
}

func NewErroredVerifier(address common.Address, backend bind.ContractBackend) (*ErroredVerifier, error) {
	abi, err := abi.JSON(strings.NewReader(ErroredVerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindErroredVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ErroredVerifier{address: address, abi: abi, ErroredVerifierCaller: ErroredVerifierCaller{contract: contract}, ErroredVerifierTransactor: ErroredVerifierTransactor{contract: contract}, ErroredVerifierFilterer: ErroredVerifierFilterer{contract: contract}}, nil
}

func NewErroredVerifierCaller(address common.Address, caller bind.ContractCaller) (*ErroredVerifierCaller, error) {
	contract, err := bindErroredVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ErroredVerifierCaller{contract: contract}, nil
}

func NewErroredVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ErroredVerifierTransactor, error) {
	contract, err := bindErroredVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ErroredVerifierTransactor{contract: contract}, nil
}

func NewErroredVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ErroredVerifierFilterer, error) {
	contract, err := bindErroredVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ErroredVerifierFilterer{contract: contract}, nil
}

func bindErroredVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ErroredVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ErroredVerifier *ErroredVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ErroredVerifier.Contract.ErroredVerifierCaller.contract.Call(opts, result, method, params...)
}

func (_ErroredVerifier *ErroredVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ErroredVerifier.Contract.ErroredVerifierTransactor.contract.Transfer(opts)
}

func (_ErroredVerifier *ErroredVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ErroredVerifier.Contract.ErroredVerifierTransactor.contract.Transact(opts, method, params...)
}

func (_ErroredVerifier *ErroredVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ErroredVerifier.Contract.contract.Call(opts, result, method, params...)
}

func (_ErroredVerifier *ErroredVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ErroredVerifier.Contract.contract.Transfer(opts)
}

func (_ErroredVerifier *ErroredVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ErroredVerifier.Contract.contract.Transact(opts, method, params...)
}

func (_ErroredVerifier *ErroredVerifierCaller) ActivateConfig(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) error {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "activateConfig", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

func (_ErroredVerifier *ErroredVerifierSession) ActivateConfig(arg0 [32]byte, arg1 [32]byte) error {
	return _ErroredVerifier.Contract.ActivateConfig(&_ErroredVerifier.CallOpts, arg0, arg1)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) ActivateConfig(arg0 [32]byte, arg1 [32]byte) error {
	return _ErroredVerifier.Contract.ActivateConfig(&_ErroredVerifier.CallOpts, arg0, arg1)
}

func (_ErroredVerifier *ErroredVerifierCaller) ActivateFeed(opts *bind.CallOpts, arg0 [32]byte) error {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "activateFeed", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_ErroredVerifier *ErroredVerifierSession) ActivateFeed(arg0 [32]byte) error {
	return _ErroredVerifier.Contract.ActivateFeed(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) ActivateFeed(arg0 [32]byte) error {
	return _ErroredVerifier.Contract.ActivateFeed(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCaller) DeactivateConfig(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) error {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "deactivateConfig", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

func (_ErroredVerifier *ErroredVerifierSession) DeactivateConfig(arg0 [32]byte, arg1 [32]byte) error {
	return _ErroredVerifier.Contract.DeactivateConfig(&_ErroredVerifier.CallOpts, arg0, arg1)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) DeactivateConfig(arg0 [32]byte, arg1 [32]byte) error {
	return _ErroredVerifier.Contract.DeactivateConfig(&_ErroredVerifier.CallOpts, arg0, arg1)
}

func (_ErroredVerifier *ErroredVerifierCaller) DeactivateFeed(opts *bind.CallOpts, arg0 [32]byte) error {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "deactivateFeed", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_ErroredVerifier *ErroredVerifierSession) DeactivateFeed(arg0 [32]byte) error {
	return _ErroredVerifier.Contract.DeactivateFeed(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) DeactivateFeed(arg0 [32]byte) error {
	return _ErroredVerifier.Contract.DeactivateFeed(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCaller) LatestConfigDetails(opts *bind.CallOpts, arg0 [32]byte) (uint32, uint32, [32]byte, error) {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "latestConfigDetails", arg0)

	if err != nil {
		return *new(uint32), *new(uint32), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return out0, out1, out2, err

}

func (_ErroredVerifier *ErroredVerifierSession) LatestConfigDetails(arg0 [32]byte) (uint32, uint32, [32]byte, error) {
	return _ErroredVerifier.Contract.LatestConfigDetails(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) LatestConfigDetails(arg0 [32]byte) (uint32, uint32, [32]byte, error) {
	return _ErroredVerifier.Contract.LatestConfigDetails(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts, arg0 [32]byte) (bool, [32]byte, uint32, error) {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "latestConfigDigestAndEpoch", arg0)

	if err != nil {
		return *new(bool), *new([32]byte), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	out2 := *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return out0, out1, out2, err

}

func (_ErroredVerifier *ErroredVerifierSession) LatestConfigDigestAndEpoch(arg0 [32]byte) (bool, [32]byte, uint32, error) {
	return _ErroredVerifier.Contract.LatestConfigDigestAndEpoch(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) LatestConfigDigestAndEpoch(arg0 [32]byte) (bool, [32]byte, uint32, error) {
	return _ErroredVerifier.Contract.LatestConfigDigestAndEpoch(&_ErroredVerifier.CallOpts, arg0)
}

func (_ErroredVerifier *ErroredVerifierCaller) SetConfig(opts *bind.CallOpts, arg0 [32]byte, arg1 []common.Address, arg2 [][32]byte, arg3 uint8, arg4 []byte, arg5 uint64, arg6 []byte, arg7 []CommonAddressAndWeight) error {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "setConfig", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)

	if err != nil {
		return err
	}

	return err

}

func (_ErroredVerifier *ErroredVerifierSession) SetConfig(arg0 [32]byte, arg1 []common.Address, arg2 [][32]byte, arg3 uint8, arg4 []byte, arg5 uint64, arg6 []byte, arg7 []CommonAddressAndWeight) error {
	return _ErroredVerifier.Contract.SetConfig(&_ErroredVerifier.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) SetConfig(arg0 [32]byte, arg1 []common.Address, arg2 [][32]byte, arg3 uint8, arg4 []byte, arg5 uint64, arg6 []byte, arg7 []CommonAddressAndWeight) error {
	return _ErroredVerifier.Contract.SetConfig(&_ErroredVerifier.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

func (_ErroredVerifier *ErroredVerifierCaller) SetConfigFromSource(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int, arg2 common.Address, arg3 uint32, arg4 []common.Address, arg5 [][32]byte, arg6 uint8, arg7 []byte, arg8 uint64, arg9 []byte, arg10 []CommonAddressAndWeight) error {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "setConfigFromSource", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)

	if err != nil {
		return err
	}

	return err

}

func (_ErroredVerifier *ErroredVerifierSession) SetConfigFromSource(arg0 [32]byte, arg1 *big.Int, arg2 common.Address, arg3 uint32, arg4 []common.Address, arg5 [][32]byte, arg6 uint8, arg7 []byte, arg8 uint64, arg9 []byte, arg10 []CommonAddressAndWeight) error {
	return _ErroredVerifier.Contract.SetConfigFromSource(&_ErroredVerifier.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) SetConfigFromSource(arg0 [32]byte, arg1 *big.Int, arg2 common.Address, arg3 uint32, arg4 []common.Address, arg5 [][32]byte, arg6 uint8, arg7 []byte, arg8 uint64, arg9 []byte, arg10 []CommonAddressAndWeight) error {
	return _ErroredVerifier.Contract.SetConfigFromSource(&_ErroredVerifier.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
}

func (_ErroredVerifier *ErroredVerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ErroredVerifier *ErroredVerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ErroredVerifier.Contract.SupportsInterface(&_ErroredVerifier.CallOpts, interfaceId)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ErroredVerifier.Contract.SupportsInterface(&_ErroredVerifier.CallOpts, interfaceId)
}

func (_ErroredVerifier *ErroredVerifierCaller) Verify(opts *bind.CallOpts, arg0 []byte, arg1 common.Address) ([]byte, error) {
	var out []interface{}
	err := _ErroredVerifier.contract.Call(opts, &out, "verify", arg0, arg1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_ErroredVerifier *ErroredVerifierSession) Verify(arg0 []byte, arg1 common.Address) ([]byte, error) {
	return _ErroredVerifier.Contract.Verify(&_ErroredVerifier.CallOpts, arg0, arg1)
}

func (_ErroredVerifier *ErroredVerifierCallerSession) Verify(arg0 []byte, arg1 common.Address) ([]byte, error) {
	return _ErroredVerifier.Contract.Verify(&_ErroredVerifier.CallOpts, arg0, arg1)
}

func (_ErroredVerifier *ErroredVerifier) Address() common.Address {
	return _ErroredVerifier.address
}

type ErroredVerifierInterface interface {
	ActivateConfig(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) error

	ActivateFeed(opts *bind.CallOpts, arg0 [32]byte) error

	DeactivateConfig(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) error

	DeactivateFeed(opts *bind.CallOpts, arg0 [32]byte) error

	LatestConfigDetails(opts *bind.CallOpts, arg0 [32]byte) (uint32, uint32, [32]byte, error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts, arg0 [32]byte) (bool, [32]byte, uint32, error)

	SetConfig(opts *bind.CallOpts, arg0 [32]byte, arg1 []common.Address, arg2 [][32]byte, arg3 uint8, arg4 []byte, arg5 uint64, arg6 []byte, arg7 []CommonAddressAndWeight) error

	SetConfigFromSource(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int, arg2 common.Address, arg3 uint32, arg4 []common.Address, arg5 [][32]byte, arg6 uint8, arg7 []byte, arg8 uint64, arg9 []byte, arg10 []CommonAddressAndWeight) error

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	Verify(opts *bind.CallOpts, arg0 []byte, arg1 common.Address) ([]byte, error)

	Address() common.Address
}
