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

type CCIPHomeOCR3Config struct {
	PluginType            uint8
	ChainSelector         uint64
	FRoleDON              uint8
	OffchainConfigVersion uint64
	OfframpAddress        []byte
	RmnHomeAddress        []byte
	Nodes                 []CCIPHomeOCR3Node
	OffchainConfig        []byte
}

type CCIPHomeOCR3Node struct {
	P2pId          [32]byte
	SignerKey      []byte
	TransmitterKey []byte
}

type IRMNRemoteSignature struct {
	R [32]byte
	S [32]byte
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type OffRampCommitReport struct {
	PriceUpdates  InternalPriceUpdates
	MerkleRoots   []InternalMerkleRoot
	RmnSignatures []IRMNRemoteSignature
}

type RMNRemoteReport struct {
	DestChainId                 *big.Int
	DestChainSelector           uint64
	RmnRemoteContractAddress    common.Address
	OfframpAddress              common.Address
	RmnHomeContractConfigDigest [32]byte
	MerkleRoots                 []InternalMerkleRoot
}

var ICCIPEncodingUtilsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"exposeCommitReport\",\"inputs\":[{\"name\":\"commitReport\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]},{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structIRMNRemote.Signature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeOCR3Config\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Config[]\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeRmnReport\",\"inputs\":[{\"name\":\"rmnReportVersion\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rmnReport\",\"type\":\"tuple\",\"internalType\":\"structRMNRemote.Report\",\"components\":[{\"name\":\"destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemoteContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"offrampAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

var ICCIPEncodingUtilsABI = ICCIPEncodingUtilsMetaData.ABI

type ICCIPEncodingUtils struct {
	address common.Address
	abi     abi.ABI
	ICCIPEncodingUtilsCaller
	ICCIPEncodingUtilsTransactor
	ICCIPEncodingUtilsFilterer
}

type ICCIPEncodingUtilsCaller struct {
	contract *bind.BoundContract
}

type ICCIPEncodingUtilsTransactor struct {
	contract *bind.BoundContract
}

type ICCIPEncodingUtilsFilterer struct {
	contract *bind.BoundContract
}

type ICCIPEncodingUtilsSession struct {
	Contract     *ICCIPEncodingUtils
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ICCIPEncodingUtilsCallerSession struct {
	Contract *ICCIPEncodingUtilsCaller
	CallOpts bind.CallOpts
}

type ICCIPEncodingUtilsTransactorSession struct {
	Contract     *ICCIPEncodingUtilsTransactor
	TransactOpts bind.TransactOpts
}

type ICCIPEncodingUtilsRaw struct {
	Contract *ICCIPEncodingUtils
}

type ICCIPEncodingUtilsCallerRaw struct {
	Contract *ICCIPEncodingUtilsCaller
}

type ICCIPEncodingUtilsTransactorRaw struct {
	Contract *ICCIPEncodingUtilsTransactor
}

func NewICCIPEncodingUtils(address common.Address, backend bind.ContractBackend) (*ICCIPEncodingUtils, error) {
	abi, err := abi.JSON(strings.NewReader(ICCIPEncodingUtilsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindICCIPEncodingUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICCIPEncodingUtils{address: address, abi: abi, ICCIPEncodingUtilsCaller: ICCIPEncodingUtilsCaller{contract: contract}, ICCIPEncodingUtilsTransactor: ICCIPEncodingUtilsTransactor{contract: contract}, ICCIPEncodingUtilsFilterer: ICCIPEncodingUtilsFilterer{contract: contract}}, nil
}

func NewICCIPEncodingUtilsCaller(address common.Address, caller bind.ContractCaller) (*ICCIPEncodingUtilsCaller, error) {
	contract, err := bindICCIPEncodingUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICCIPEncodingUtilsCaller{contract: contract}, nil
}

func NewICCIPEncodingUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*ICCIPEncodingUtilsTransactor, error) {
	contract, err := bindICCIPEncodingUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICCIPEncodingUtilsTransactor{contract: contract}, nil
}

func NewICCIPEncodingUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*ICCIPEncodingUtilsFilterer, error) {
	contract, err := bindICCIPEncodingUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICCIPEncodingUtilsFilterer{contract: contract}, nil
}

func bindICCIPEncodingUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICCIPEncodingUtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICCIPEncodingUtils.Contract.ICCIPEncodingUtilsCaller.contract.Call(opts, result, method, params...)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.Contract.ICCIPEncodingUtilsTransactor.contract.Transfer(opts)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.Contract.ICCIPEncodingUtilsTransactor.contract.Transact(opts, method, params...)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICCIPEncodingUtils.Contract.contract.Call(opts, result, method, params...)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.Contract.contract.Transfer(opts)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.Contract.contract.Transact(opts, method, params...)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsCaller) ExposeCommitReport(opts *bind.CallOpts, commitReport OffRampCommitReport) ([]byte, error) {
	var out []interface{}
	err := _ICCIPEncodingUtils.contract.Call(opts, &out, "exposeCommitReport", commitReport)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsSession) ExposeCommitReport(commitReport OffRampCommitReport) ([]byte, error) {
	return _ICCIPEncodingUtils.Contract.ExposeCommitReport(&_ICCIPEncodingUtils.CallOpts, commitReport)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsCallerSession) ExposeCommitReport(commitReport OffRampCommitReport) ([]byte, error) {
	return _ICCIPEncodingUtils.Contract.ExposeCommitReport(&_ICCIPEncodingUtils.CallOpts, commitReport)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsCaller) ExposeOCR3Config(opts *bind.CallOpts, config []CCIPHomeOCR3Config) ([]byte, error) {
	var out []interface{}
	err := _ICCIPEncodingUtils.contract.Call(opts, &out, "exposeOCR3Config", config)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsSession) ExposeOCR3Config(config []CCIPHomeOCR3Config) ([]byte, error) {
	return _ICCIPEncodingUtils.Contract.ExposeOCR3Config(&_ICCIPEncodingUtils.CallOpts, config)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsCallerSession) ExposeOCR3Config(config []CCIPHomeOCR3Config) ([]byte, error) {
	return _ICCIPEncodingUtils.Contract.ExposeOCR3Config(&_ICCIPEncodingUtils.CallOpts, config)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsTransactor) ExposeRmnReport(opts *bind.TransactOpts, rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.contract.Transact(opts, "exposeRmnReport", rmnReportVersion, rmnReport)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsSession) ExposeRmnReport(rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.Contract.ExposeRmnReport(&_ICCIPEncodingUtils.TransactOpts, rmnReportVersion, rmnReport)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtilsTransactorSession) ExposeRmnReport(rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _ICCIPEncodingUtils.Contract.ExposeRmnReport(&_ICCIPEncodingUtils.TransactOpts, rmnReportVersion, rmnReport)
}

func (_ICCIPEncodingUtils *ICCIPEncodingUtils) Address() common.Address {
	return _ICCIPEncodingUtils.address
}

type ICCIPEncodingUtilsInterface interface {
	ExposeCommitReport(opts *bind.CallOpts, commitReport OffRampCommitReport) ([]byte, error)

	ExposeOCR3Config(opts *bind.CallOpts, config []CCIPHomeOCR3Config) ([]byte, error)

	ExposeRmnReport(opts *bind.TransactOpts, rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error)

	Address() common.Address
}
