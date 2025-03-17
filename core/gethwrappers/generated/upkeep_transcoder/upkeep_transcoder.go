// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_transcoder

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

var UpkeepTranscoderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"transcodeUpkeeps\",\"inputs\":[{\"name\":\"fromVersion\",\"type\":\"uint8\",\"internalType\":\"enumUpkeepFormat\"},{\"name\":\"toVersion\",\"type\":\"uint8\",\"internalType\":\"enumUpkeepFormat\"},{\"name\":\"encodedUpkeeps\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidTranscoding\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b5061029b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063181f5a771461003b578063c71249ab1461008d575b600080fd5b6100776040518060400160405280601681526020017f55706b6565705472616e73636f64657220312e302e300000000000000000000081525081565b6040516100849190610245565b60405180910390f35b61007761009b36600461014c565b60608360028111156100af576100af61025f565b8560028111156100c1576100c161025f565b146100f8576040517f90aaccc300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82828080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509298975050505050505050565b80356003811061014757600080fd5b919050565b6000806000806060858703121561016257600080fd5b61016b85610138565b935061017960208601610138565b9250604085013567ffffffffffffffff8082111561019657600080fd5b818701915087601f8301126101aa57600080fd5b8135818111156101b957600080fd5b8860208285010111156101cb57600080fd5b95989497505060200194505050565b6000815180845260005b81811015610200576020818501810151868301820152016101e4565b81811115610212576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061025860208301846101da565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fdfea164736f6c6343000806000a",
}

var UpkeepTranscoderABI = UpkeepTranscoderMetaData.ABI

var UpkeepTranscoderBin = UpkeepTranscoderMetaData.Bin

func DeployUpkeepTranscoder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UpkeepTranscoder, error) {
	parsed, err := UpkeepTranscoderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepTranscoderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepTranscoder{address: address, abi: *parsed, UpkeepTranscoderCaller: UpkeepTranscoderCaller{contract: contract}, UpkeepTranscoderTransactor: UpkeepTranscoderTransactor{contract: contract}, UpkeepTranscoderFilterer: UpkeepTranscoderFilterer{contract: contract}}, nil
}

type UpkeepTranscoder struct {
	address common.Address
	abi     abi.ABI
	UpkeepTranscoderCaller
	UpkeepTranscoderTransactor
	UpkeepTranscoderFilterer
}

type UpkeepTranscoderCaller struct {
	contract *bind.BoundContract
}

type UpkeepTranscoderTransactor struct {
	contract *bind.BoundContract
}

type UpkeepTranscoderFilterer struct {
	contract *bind.BoundContract
}

type UpkeepTranscoderSession struct {
	Contract     *UpkeepTranscoder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepTranscoderCallerSession struct {
	Contract *UpkeepTranscoderCaller
	CallOpts bind.CallOpts
}

type UpkeepTranscoderTransactorSession struct {
	Contract     *UpkeepTranscoderTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepTranscoderRaw struct {
	Contract *UpkeepTranscoder
}

type UpkeepTranscoderCallerRaw struct {
	Contract *UpkeepTranscoderCaller
}

type UpkeepTranscoderTransactorRaw struct {
	Contract *UpkeepTranscoderTransactor
}

func NewUpkeepTranscoder(address common.Address, backend bind.ContractBackend) (*UpkeepTranscoder, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepTranscoderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepTranscoder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoder{address: address, abi: abi, UpkeepTranscoderCaller: UpkeepTranscoderCaller{contract: contract}, UpkeepTranscoderTransactor: UpkeepTranscoderTransactor{contract: contract}, UpkeepTranscoderFilterer: UpkeepTranscoderFilterer{contract: contract}}, nil
}

func NewUpkeepTranscoderCaller(address common.Address, caller bind.ContractCaller) (*UpkeepTranscoderCaller, error) {
	contract, err := bindUpkeepTranscoder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoderCaller{contract: contract}, nil
}

func NewUpkeepTranscoderTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepTranscoderTransactor, error) {
	contract, err := bindUpkeepTranscoder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoderTransactor{contract: contract}, nil
}

func NewUpkeepTranscoderFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepTranscoderFilterer, error) {
	contract, err := bindUpkeepTranscoder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoderFilterer{contract: contract}, nil
}

func bindUpkeepTranscoder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UpkeepTranscoderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_UpkeepTranscoder *UpkeepTranscoderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepTranscoder.Contract.UpkeepTranscoderCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepTranscoder *UpkeepTranscoderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.UpkeepTranscoderTransactor.contract.Transfer(opts)
}

func (_UpkeepTranscoder *UpkeepTranscoderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.UpkeepTranscoderTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepTranscoder *UpkeepTranscoderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepTranscoder.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepTranscoder *UpkeepTranscoderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.contract.Transfer(opts)
}

func (_UpkeepTranscoder *UpkeepTranscoderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepTranscoder *UpkeepTranscoderCaller) TranscodeUpkeeps(opts *bind.CallOpts, fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error) {
	var out []interface{}
	err := _UpkeepTranscoder.contract.Call(opts, &out, "transcodeUpkeeps", fromVersion, toVersion, encodedUpkeeps)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_UpkeepTranscoder *UpkeepTranscoderSession) TranscodeUpkeeps(fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error) {
	return _UpkeepTranscoder.Contract.TranscodeUpkeeps(&_UpkeepTranscoder.CallOpts, fromVersion, toVersion, encodedUpkeeps)
}

func (_UpkeepTranscoder *UpkeepTranscoderCallerSession) TranscodeUpkeeps(fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error) {
	return _UpkeepTranscoder.Contract.TranscodeUpkeeps(&_UpkeepTranscoder.CallOpts, fromVersion, toVersion, encodedUpkeeps)
}

func (_UpkeepTranscoder *UpkeepTranscoderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepTranscoder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepTranscoder *UpkeepTranscoderSession) TypeAndVersion() (string, error) {
	return _UpkeepTranscoder.Contract.TypeAndVersion(&_UpkeepTranscoder.CallOpts)
}

func (_UpkeepTranscoder *UpkeepTranscoderCallerSession) TypeAndVersion() (string, error) {
	return _UpkeepTranscoder.Contract.TypeAndVersion(&_UpkeepTranscoder.CallOpts)
}

func (_UpkeepTranscoder *UpkeepTranscoder) Address() common.Address {
	return _UpkeepTranscoder.address
}

type UpkeepTranscoderInterface interface {
	TranscodeUpkeeps(opts *bind.CallOpts, fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Address() common.Address
}
