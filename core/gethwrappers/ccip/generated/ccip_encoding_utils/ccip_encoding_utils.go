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
	PriceUpdates         InternalPriceUpdates
	BlessedMerkleRoots   []InternalMerkleRoot
	UnblessedMerkleRoots []InternalMerkleRoot
	RmnSignatures        []IRMNRemoteSignature
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
	ABI: "[{\"type\":\"function\",\"name\":\"exposeCommitReport\",\"inputs\":[{\"name\":\"commitReport\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]},{\"name\":\"blessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"unblessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structIRMNRemote.Signature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeOCR3Config\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Config[]\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeRmnReport\",\"inputs\":[{\"name\":\"rmnReportVersion\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rmnReport\",\"type\":\"tuple\",\"internalType\":\"structRMNRemote.Report\",\"components\":[{\"name\":\"destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemoteContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"offrampAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

var EncodingUtilsABI = EncodingUtilsMetaData.ABI

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

func (_EncodingUtils *EncodingUtilsCaller) ExposeCommitReport(opts *bind.CallOpts, commitReport OffRampCommitReport) ([]byte, error) {
	var out []interface{}
	err := _EncodingUtils.contract.Call(opts, &out, "exposeCommitReport", commitReport)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_EncodingUtils *EncodingUtilsSession) ExposeCommitReport(commitReport OffRampCommitReport) ([]byte, error) {
	return _EncodingUtils.Contract.ExposeCommitReport(&_EncodingUtils.CallOpts, commitReport)
}

func (_EncodingUtils *EncodingUtilsCallerSession) ExposeCommitReport(commitReport OffRampCommitReport) ([]byte, error) {
	return _EncodingUtils.Contract.ExposeCommitReport(&_EncodingUtils.CallOpts, commitReport)
}

func (_EncodingUtils *EncodingUtilsCaller) ExposeOCR3Config(opts *bind.CallOpts, config []CCIPHomeOCR3Config) ([]byte, error) {
	var out []interface{}
	err := _EncodingUtils.contract.Call(opts, &out, "exposeOCR3Config", config)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_EncodingUtils *EncodingUtilsSession) ExposeOCR3Config(config []CCIPHomeOCR3Config) ([]byte, error) {
	return _EncodingUtils.Contract.ExposeOCR3Config(&_EncodingUtils.CallOpts, config)
}

func (_EncodingUtils *EncodingUtilsCallerSession) ExposeOCR3Config(config []CCIPHomeOCR3Config) ([]byte, error) {
	return _EncodingUtils.Contract.ExposeOCR3Config(&_EncodingUtils.CallOpts, config)
}

func (_EncodingUtils *EncodingUtilsTransactor) ExposeRmnReport(opts *bind.TransactOpts, rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _EncodingUtils.contract.Transact(opts, "exposeRmnReport", rmnReportVersion, rmnReport)
}

func (_EncodingUtils *EncodingUtilsSession) ExposeRmnReport(rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _EncodingUtils.Contract.ExposeRmnReport(&_EncodingUtils.TransactOpts, rmnReportVersion, rmnReport)
}

func (_EncodingUtils *EncodingUtilsTransactorSession) ExposeRmnReport(rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error) {
	return _EncodingUtils.Contract.ExposeRmnReport(&_EncodingUtils.TransactOpts, rmnReportVersion, rmnReport)
}

func (_EncodingUtils *EncodingUtils) Address() common.Address {
	return _EncodingUtils.address
}

type EncodingUtilsInterface interface {
	ExposeCommitReport(opts *bind.CallOpts, commitReport OffRampCommitReport) ([]byte, error)

	ExposeOCR3Config(opts *bind.CallOpts, config []CCIPHomeOCR3Config) ([]byte, error)

	ExposeRmnReport(opts *bind.TransactOpts, rmnReportVersion [32]byte, rmnReport RMNRemoteReport) (*types.Transaction, error)

	Address() common.Address
}
