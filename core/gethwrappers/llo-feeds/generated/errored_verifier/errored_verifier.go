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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"activateFeed\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"deactivateFeed\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"setConfigFromSource\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610c2e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063b70d929d11610076578063e7db9c2a1161005b578063e7db9c2a146101d1578063e84f128e146101e4578063f01072211461021a57600080fd5b8063b70d929d14610188578063ded6307c146101be57600080fd5b80633dd86430116100a75780633dd864301461014d578063564a0a7a1461016257806394d959801461017557600080fd5b806301ffc9a7146100c35780633d3ac1b51461012d575b600080fd5b6101186100d136600461059a565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61014061013b366004610741565b610228565b604051610124919061078f565b61016061015b3660046107fb565b610292565b005b6101606101703660046107fb565b6102f4565b610160610183366004610814565b610356565b61019b6101963660046107fb565b6103b8565b604080519315158452602084019290925263ffffffff1690820152606001610124565b6101606101cc366004610814565b610447565b6101606101df3660046109f1565b6104a9565b6101f76101f23660046107fb565b61050b565b6040805163ffffffff948516815293909216602084015290820152606001610124565b6101606101df366004610b24565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f4661696c656420746f207665726966790000000000000000000000000000000060448201526060906064015b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4661696c656420746f20616374697661746520666565640000000000000000006044820152606401610289565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f4661696c656420746f20646561637469766174652066656564000000000000006044820152606401610289565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4661696c656420746f206465616374697661746520636f6e66696700000000006044820152606401610289565b60008060006040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610289906020808252602c908201527f4661696c656420746f20676574206c617465737420636f6e666967206469676560408201527f737420616e642065706f63680000000000000000000000000000000000000000606082015260800190565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f4661696c656420746f20616374697661746520636f6e666967000000000000006044820152606401610289565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4661696c656420746f2073657420636f6e6669670000000000000000000000006044820152606401610289565b60008060006040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102899060208082526023908201527f4661696c656420746f20676574206c617465737420636f6e666967206465746160408201527f696c730000000000000000000000000000000000000000000000000000000000606082015260800190565b6000602082840312156105ac57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146105dc57600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610635576106356105e3565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610682576106826105e3565b604052919050565b600082601f83011261069b57600080fd5b813567ffffffffffffffff8111156106b5576106b56105e3565b6106e660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161063b565b8181528460208386010111156106fb57600080fd5b816020850160208301376000918101602001919091529392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461073c57600080fd5b919050565b6000806040838503121561075457600080fd5b823567ffffffffffffffff81111561076b57600080fd5b6107778582860161068a565b92505061078660208401610718565b90509250929050565b600060208083528351808285015260005b818110156107bc578581018301518582016040015282016107a0565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561080d57600080fd5b5035919050565b6000806040838503121561082757600080fd5b50508035926020909101359150565b803563ffffffff8116811461073c57600080fd5b600067ffffffffffffffff821115610864576108646105e3565b5060051b60200190565b600082601f83011261087f57600080fd5b8135602061089461088f8361084a565b61063b565b82815260059290921b840181019181810190868411156108b357600080fd5b8286015b848110156108d5576108c881610718565b83529183019183016108b7565b509695505050505050565b600082601f8301126108f157600080fd5b8135602061090161088f8361084a565b82815260059290921b8401810191818101908684111561092057600080fd5b8286015b848110156108d55780358352918301918301610924565b803560ff8116811461073c57600080fd5b803567ffffffffffffffff8116811461073c57600080fd5b600082601f83011261097557600080fd5b8135602061098561088f8361084a565b82815260069290921b840181019181810190868411156109a457600080fd5b8286015b848110156108d557604081890312156109c15760008081fd5b6109c9610612565b6109d282610718565b81526109df85830161094c565b818601528352918301916040016109a8565b60008060008060008060008060008060006101608c8e031215610a1357600080fd5b8b359a5060208c01359950610a2a60408d01610718565b9850610a3860608d01610836565b975067ffffffffffffffff8060808e01351115610a5457600080fd5b610a648e60808f01358f0161086e565b97508060a08e01351115610a7757600080fd5b610a878e60a08f01358f016108e0565b9650610a9560c08e0161093b565b95508060e08e01351115610aa857600080fd5b610ab88e60e08f01358f0161068a565b9450610ac76101008e0161094c565b9350806101208e01351115610adb57600080fd5b610aec8e6101208f01358f0161068a565b9250806101408e01351115610b0057600080fd5b50610b128d6101408e01358e01610964565b90509295989b509295989b9093969950565b600080600080600080600080610100898b031215610b4157600080fd5b88359750602089013567ffffffffffffffff80821115610b6057600080fd5b610b6c8c838d0161086e565b985060408b0135915080821115610b8257600080fd5b610b8e8c838d016108e0565b9750610b9c60608c0161093b565b965060808b0135915080821115610bb257600080fd5b610bbe8c838d0161068a565b9550610bcc60a08c0161094c565b945060c08b0135915080821115610be257600080fd5b610bee8c838d0161068a565b935060e08b0135915080821115610c0457600080fd5b50610c118b828c01610964565b915050929598509295989093965056fea164736f6c6343000813000a",
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
