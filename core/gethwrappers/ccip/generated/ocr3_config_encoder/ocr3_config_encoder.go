// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr3_config_encoder

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

type CCIPConfigTypesOCR3Config struct {
	PluginType            uint8
	ChainSelector         uint64
	FRoleDON              uint8
	OffchainConfigVersion uint64
	OfframpAddress        []byte
	Nodes                 []CCIPConfigTypesOCR3Node
	OffchainConfig        []byte
}

type CCIPConfigTypesOCR3Node struct {
	P2pId          [32]byte
	SignerKey      []byte
	TransmitterKey []byte
}

var IOCR3ConfigEncoderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"FRoleDON\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offrampAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPConfigTypes.OCR3Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPConfigTypes.OCR3Config[]\",\"name\":\"config\",\"type\":\"tuple[]\"}],\"name\":\"exposeOCR3Config\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var IOCR3ConfigEncoderABI = IOCR3ConfigEncoderMetaData.ABI

type IOCR3ConfigEncoder struct {
	address common.Address
	abi     abi.ABI
	IOCR3ConfigEncoderCaller
	IOCR3ConfigEncoderTransactor
	IOCR3ConfigEncoderFilterer
}

type IOCR3ConfigEncoderCaller struct {
	contract *bind.BoundContract
}

type IOCR3ConfigEncoderTransactor struct {
	contract *bind.BoundContract
}

type IOCR3ConfigEncoderFilterer struct {
	contract *bind.BoundContract
}

type IOCR3ConfigEncoderSession struct {
	Contract     *IOCR3ConfigEncoder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IOCR3ConfigEncoderCallerSession struct {
	Contract *IOCR3ConfigEncoderCaller
	CallOpts bind.CallOpts
}

type IOCR3ConfigEncoderTransactorSession struct {
	Contract     *IOCR3ConfigEncoderTransactor
	TransactOpts bind.TransactOpts
}

type IOCR3ConfigEncoderRaw struct {
	Contract *IOCR3ConfigEncoder
}

type IOCR3ConfigEncoderCallerRaw struct {
	Contract *IOCR3ConfigEncoderCaller
}

type IOCR3ConfigEncoderTransactorRaw struct {
	Contract *IOCR3ConfigEncoderTransactor
}

func NewIOCR3ConfigEncoder(address common.Address, backend bind.ContractBackend) (*IOCR3ConfigEncoder, error) {
	abi, err := abi.JSON(strings.NewReader(IOCR3ConfigEncoderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIOCR3ConfigEncoder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IOCR3ConfigEncoder{address: address, abi: abi, IOCR3ConfigEncoderCaller: IOCR3ConfigEncoderCaller{contract: contract}, IOCR3ConfigEncoderTransactor: IOCR3ConfigEncoderTransactor{contract: contract}, IOCR3ConfigEncoderFilterer: IOCR3ConfigEncoderFilterer{contract: contract}}, nil
}

func NewIOCR3ConfigEncoderCaller(address common.Address, caller bind.ContractCaller) (*IOCR3ConfigEncoderCaller, error) {
	contract, err := bindIOCR3ConfigEncoder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IOCR3ConfigEncoderCaller{contract: contract}, nil
}

func NewIOCR3ConfigEncoderTransactor(address common.Address, transactor bind.ContractTransactor) (*IOCR3ConfigEncoderTransactor, error) {
	contract, err := bindIOCR3ConfigEncoder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IOCR3ConfigEncoderTransactor{contract: contract}, nil
}

func NewIOCR3ConfigEncoderFilterer(address common.Address, filterer bind.ContractFilterer) (*IOCR3ConfigEncoderFilterer, error) {
	contract, err := bindIOCR3ConfigEncoder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IOCR3ConfigEncoderFilterer{contract: contract}, nil
}

func bindIOCR3ConfigEncoder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IOCR3ConfigEncoderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOCR3ConfigEncoder.Contract.IOCR3ConfigEncoderCaller.contract.Call(opts, result, method, params...)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOCR3ConfigEncoder.Contract.IOCR3ConfigEncoderTransactor.contract.Transfer(opts)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOCR3ConfigEncoder.Contract.IOCR3ConfigEncoderTransactor.contract.Transact(opts, method, params...)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOCR3ConfigEncoder.Contract.contract.Call(opts, result, method, params...)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOCR3ConfigEncoder.Contract.contract.Transfer(opts)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOCR3ConfigEncoder.Contract.contract.Transact(opts, method, params...)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderCaller) ExposeOCR3Config(opts *bind.CallOpts, config []CCIPConfigTypesOCR3Config) ([]byte, error) {
	var out []interface{}
	err := _IOCR3ConfigEncoder.contract.Call(opts, &out, "exposeOCR3Config", config)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderSession) ExposeOCR3Config(config []CCIPConfigTypesOCR3Config) ([]byte, error) {
	return _IOCR3ConfigEncoder.Contract.ExposeOCR3Config(&_IOCR3ConfigEncoder.CallOpts, config)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoderCallerSession) ExposeOCR3Config(config []CCIPConfigTypesOCR3Config) ([]byte, error) {
	return _IOCR3ConfigEncoder.Contract.ExposeOCR3Config(&_IOCR3ConfigEncoder.CallOpts, config)
}

func (_IOCR3ConfigEncoder *IOCR3ConfigEncoder) Address() common.Address {
	return _IOCR3ConfigEncoder.address
}

type IOCR3ConfigEncoderInterface interface {
	ExposeOCR3Config(opts *bind.CallOpts, config []CCIPConfigTypesOCR3Config) ([]byte, error)

	Address() common.Address
}
