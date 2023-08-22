// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr

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

var OCR2DRMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyUrl\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"}]",
	Bin: "0x602d6037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea164736f6c6343000806000a",
}

var OCR2DRABI = OCR2DRMetaData.ABI

var OCR2DRBin = OCR2DRMetaData.Bin

func DeployOCR2DR(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR2DR, error) {
	parsed, err := OCR2DRMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DRBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2DR{OCR2DRCaller: OCR2DRCaller{contract: contract}, OCR2DRTransactor: OCR2DRTransactor{contract: contract}, OCR2DRFilterer: OCR2DRFilterer{contract: contract}}, nil
}

type OCR2DR struct {
	address common.Address
	abi     abi.ABI
	OCR2DRCaller
	OCR2DRTransactor
	OCR2DRFilterer
}

type OCR2DRCaller struct {
	contract *bind.BoundContract
}

type OCR2DRTransactor struct {
	contract *bind.BoundContract
}

type OCR2DRFilterer struct {
	contract *bind.BoundContract
}

type OCR2DRSession struct {
	Contract     *OCR2DR
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DRCallerSession struct {
	Contract *OCR2DRCaller
	CallOpts bind.CallOpts
}

type OCR2DRTransactorSession struct {
	Contract     *OCR2DRTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DRRaw struct {
	Contract *OCR2DR
}

type OCR2DRCallerRaw struct {
	Contract *OCR2DRCaller
}

type OCR2DRTransactorRaw struct {
	Contract *OCR2DRTransactor
}

func NewOCR2DR(address common.Address, backend bind.ContractBackend) (*OCR2DR, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DRABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DR(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DR{address: address, abi: abi, OCR2DRCaller: OCR2DRCaller{contract: contract}, OCR2DRTransactor: OCR2DRTransactor{contract: contract}, OCR2DRFilterer: OCR2DRFilterer{contract: contract}}, nil
}

func NewOCR2DRCaller(address common.Address, caller bind.ContractCaller) (*OCR2DRCaller, error) {
	contract, err := bindOCR2DR(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRCaller{contract: contract}, nil
}

func NewOCR2DRTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DRTransactor, error) {
	contract, err := bindOCR2DR(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRTransactor{contract: contract}, nil
}

func NewOCR2DRFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DRFilterer, error) {
	contract, err := bindOCR2DR(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DRFilterer{contract: contract}, nil
}

func bindOCR2DR(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR2DRMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OCR2DR *OCR2DRRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DR.Contract.OCR2DRCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DR *OCR2DRRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DR.Contract.OCR2DRTransactor.contract.Transfer(opts)
}

func (_OCR2DR *OCR2DRRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DR.Contract.OCR2DRTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DR *OCR2DRCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DR.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DR *OCR2DRTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DR.Contract.contract.Transfer(opts)
}

func (_OCR2DR *OCR2DRTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DR.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DR *OCR2DR) Address() common.Address {
	return _OCR2DR.address
}

type OCR2DRInterface interface {
	Address() common.Address
}
