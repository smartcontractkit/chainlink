// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ccip_encoding_utils

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

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type RMNRemoteReport struct {
	DestChainId                 *big.Int
	DestChainSelector           uint64
	RmnRemoteContractAddress    common.Address
	OfframpAddress              common.Address
	RmnHomeContractConfigDigest [32]byte
	MerkleRoots                 []InternalMerkleRoot
}

var EncodingUtilsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"DoNotDeploy\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rmnReportVersion\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnRemoteContractAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offrampAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"internalType\":\"structRMNRemote.Report\",\"name\":\"rmnReport\",\"type\":\"tuple\"}],\"name\":\"_rmnReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051636f1e4f5f60e01b815260040160405180910390fdfe",
}

var EncodingUtilsABI = EncodingUtilsMetaData.ABI

var EncodingUtilsBin = EncodingUtilsMetaData.Bin

func DeployEncodingUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EncodingUtils, error) {
	parsed, err := EncodingUtilsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EncodingUtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EncodingUtils{address: address, abi: *parsed, EncodingUtilsCaller: EncodingUtilsCaller{contract: contract}, EncodingUtilsTransactor: EncodingUtilsTransactor{contract: contract}, EncodingUtilsFilterer: EncodingUtilsFilterer{contract: contract}}, nil
}

type EncodingUtils struct {
	address common.Address
	abi     abi.ABI
	EncodingUtilsCaller
	EncodingUtilsTransactor
	EncodingUtilsFilterer
}

type EncodingUtilsCaller struct {
	contract *bind.BoundContract
}

type EncodingUtilsTransactor struct {
	contract *bind.BoundContract
}

type EncodingUtilsFilterer struct {
	contract *bind.BoundContract
}

type EncodingUtilsSession struct {
	Contract     *EncodingUtils
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type EncodingUtilsCallerSession struct {
	Contract *EncodingUtilsCaller
	CallOpts bind.CallOpts
}

type EncodingUtilsTransactorSession struct {
	Contract     *EncodingUtilsTransactor
	TransactOpts bind.TransactOpts
}

type EncodingUtilsRaw struct {
	Contract *EncodingUtils
}

type EncodingUtilsCallerRaw struct {
	Contract *EncodingUtilsCaller
}

type EncodingUtilsTransactorRaw struct {
	Contract *EncodingUtilsTransactor
}

func NewEncodingUtils(address common.Address, backend bind.ContractBackend) (*EncodingUtils, error) {
	abi, err := abi.JSON(strings.NewReader(EncodingUtilsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEncodingUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EncodingUtils{address: address, abi: abi, EncodingUtilsCaller: EncodingUtilsCaller{contract: contract}, EncodingUtilsTransactor: EncodingUtilsTransactor{contract: contract}, EncodingUtilsFilterer: EncodingUtilsFilterer{contract: contract}}, nil
}

func NewEncodingUtilsCaller(address common.Address, caller bind.ContractCaller) (*EncodingUtilsCaller, error) {
	contract, err := bindEncodingUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EncodingUtilsCaller{contract: contract}, nil
}

func NewEncodingUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*EncodingUtilsTransactor, error) {
	contract, err := bindEncodingUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EncodingUtilsTransactor{contract: contract}, nil
}

func NewEncodingUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*EncodingUtilsFilterer, error) {
	contract, err := bindEncodingUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EncodingUtilsFilterer{contract: contract}, nil
}

func bindEncodingUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EncodingUtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_EncodingUtils *EncodingUtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EncodingUtils.Contract.EncodingUtilsCaller.contract.Call(opts, result, method, params...)
}

func (_EncodingUtils *EncodingUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EncodingUtils.Contract.EncodingUtilsTransactor.contract.Transfer(opts)
}

func (_EncodingUtils *EncodingUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EncodingUtils.Contract.EncodingUtilsTransactor.contract.Transact(opts, method, params...)
}

func (_EncodingUtils *EncodingUtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EncodingUtils.Contract.contract.Call(opts, result, method, params...)
}

func (_EncodingUtils *EncodingUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EncodingUtils.Contract.contract.Transfer(opts)
}

func (_EncodingUtils *EncodingUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EncodingUtils.Contract.contract.Transact(opts, method, params...)
}

func (_EncodingUtils *EncodingUtilsTransactor) RmnReport(opts *bind.TransactOpts, rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _EncodingUtils.contract.Transact(opts, "_rmnReport", rmnReportVersion, rmnReport)
}

func (_EncodingUtils *EncodingUtilsSession) RmnReport(rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _EncodingUtils.Contract.RmnReport(&_EncodingUtils.TransactOpts, rmnReportVersion, rmnReport)
}

func (_EncodingUtils *EncodingUtilsTransactorSession) RmnReport(rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _EncodingUtils.Contract.RmnReport(&_EncodingUtils.TransactOpts, rmnReportVersion, rmnReport)
}

func (_EncodingUtils *EncodingUtils) Address() common.Address {
	return _EncodingUtils.address
}

type EncodingUtilsInterface interface {
	RmnReport(opts *bind.TransactOpts, rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error)

	Address() common.Address
}
