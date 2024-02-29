// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_convenience

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

type ConditionalTrigger struct {
	BlockNum  uint32
	BlockHash [32]byte
}

type Log struct {
	Index       *big.Int
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

type LogTrigger struct {
	LogBlockHash [32]byte
	TxHash       [32]byte
	LogIndex     uint32
	BlockNum     uint32
	BlockHash    [32]byte
}

type LogTriggerConfig struct {
	ContractAddress common.Address
	FilterSelector  uint8
	Topic0          [32]byte
	Topic1          [32]byte
	Topic2          [32]byte
	Topic3          [32]byte
}

type OnchainConfig struct {
	PaymentPremiumPPB      uint32
	FlatFeeMicroLink       uint32
	CheckGasLimit          uint32
	StalenessSeconds       *big.Int
	GasCeilingMultiplier   uint16
	MinUpkeepSpend         *big.Int
	MaxPerformGas          uint32
	MaxCheckDataSize       uint32
	MaxPerformDataSize     uint32
	MaxRevertDataSize      uint32
	FallbackGasPrice       *big.Int
	FallbackLinkPrice      *big.Int
	Transcoder             common.Address
	Registrars             []common.Address
	UpkeepPrivilegeManager common.Address
	ChainModule            common.Address
	ReorgProtectionEnabled bool
}

type OnchainConfigLegacy struct {
	PaymentPremiumPPB      uint32
	FlatFeeMicroLink       uint32
	CheckGasLimit          uint32
	StalenessSeconds       *big.Int
	GasCeilingMultiplier   uint16
	MinUpkeepSpend         *big.Int
	MaxPerformGas          uint32
	MaxCheckDataSize       uint32
	MaxPerformDataSize     uint32
	MaxRevertDataSize      uint32
	FallbackGasPrice       *big.Int
	FallbackLinkPrice      *big.Int
	Transcoder             common.Address
	Registrars             []common.Address
	UpkeepPrivilegeManager common.Address
}

type Report struct {
	FastGasWei   *big.Int
	LinkNative   *big.Int
	UpkeepIds    []*big.Int
	GasLimits    []*big.Int
	Triggers     [][]byte
	PerformDatas [][]byte
}

var AutomationConvenienceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structConditionalTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_conditionalTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_log\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"logBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"logIndex\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structLogTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structLogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structOnchainConfigLegacy\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_onChainConfig21\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structOnchainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_onChainConfig22Plus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"gasLimits\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"triggers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"performDatas\",\"type\":\"bytes[]\"}],\"internalType\":\"structReport\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var AutomationConvenienceABI = AutomationConvenienceMetaData.ABI

type AutomationConvenience struct {
	address common.Address
	abi     abi.ABI
	AutomationConvenienceCaller
	AutomationConvenienceTransactor
	AutomationConvenienceFilterer
}

type AutomationConvenienceCaller struct {
	contract *bind.BoundContract
}

type AutomationConvenienceTransactor struct {
	contract *bind.BoundContract
}

type AutomationConvenienceFilterer struct {
	contract *bind.BoundContract
}

type AutomationConvenienceSession struct {
	Contract     *AutomationConvenience
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationConvenienceCallerSession struct {
	Contract *AutomationConvenienceCaller
	CallOpts bind.CallOpts
}

type AutomationConvenienceTransactorSession struct {
	Contract     *AutomationConvenienceTransactor
	TransactOpts bind.TransactOpts
}

type AutomationConvenienceRaw struct {
	Contract *AutomationConvenience
}

type AutomationConvenienceCallerRaw struct {
	Contract *AutomationConvenienceCaller
}

type AutomationConvenienceTransactorRaw struct {
	Contract *AutomationConvenienceTransactor
}

func NewAutomationConvenience(address common.Address, backend bind.ContractBackend) (*AutomationConvenience, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationConvenienceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationConvenience(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationConvenience{address: address, abi: abi, AutomationConvenienceCaller: AutomationConvenienceCaller{contract: contract}, AutomationConvenienceTransactor: AutomationConvenienceTransactor{contract: contract}, AutomationConvenienceFilterer: AutomationConvenienceFilterer{contract: contract}}, nil
}

func NewAutomationConvenienceCaller(address common.Address, caller bind.ContractCaller) (*AutomationConvenienceCaller, error) {
	contract, err := bindAutomationConvenience(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationConvenienceCaller{contract: contract}, nil
}

func NewAutomationConvenienceTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationConvenienceTransactor, error) {
	contract, err := bindAutomationConvenience(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationConvenienceTransactor{contract: contract}, nil
}

func NewAutomationConvenienceFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationConvenienceFilterer, error) {
	contract, err := bindAutomationConvenience(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationConvenienceFilterer{contract: contract}, nil
}

func bindAutomationConvenience(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationConvenienceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationConvenience *AutomationConvenienceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationConvenience.Contract.AutomationConvenienceCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationConvenience *AutomationConvenienceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.AutomationConvenienceTransactor.contract.Transfer(opts)
}

func (_AutomationConvenience *AutomationConvenienceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.AutomationConvenienceTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationConvenience *AutomationConvenienceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationConvenience.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationConvenience *AutomationConvenienceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.contract.Transfer(opts)
}

func (_AutomationConvenience *AutomationConvenienceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) ConditionalTrigger(opts *bind.TransactOpts, arg0 ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_conditionalTrigger", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) ConditionalTrigger(arg0 ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.ConditionalTrigger(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) ConditionalTrigger(arg0 ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.ConditionalTrigger(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_log", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) Log(arg0 Log) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.Log(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) Log(arg0 Log) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.Log(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) LogTrigger(opts *bind.TransactOpts, arg0 LogTrigger) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_logTrigger", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) LogTrigger(arg0 LogTrigger) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.LogTrigger(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) LogTrigger(arg0 LogTrigger) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.LogTrigger(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) LogTriggerConfig(opts *bind.TransactOpts, arg0 LogTriggerConfig) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_logTriggerConfig", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) LogTriggerConfig(arg0 LogTriggerConfig) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.LogTriggerConfig(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) LogTriggerConfig(arg0 LogTriggerConfig) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.LogTriggerConfig(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) OnChainConfig21(opts *bind.TransactOpts, arg0 OnchainConfigLegacy) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_onChainConfig21", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) OnChainConfig21(arg0 OnchainConfigLegacy) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.OnChainConfig21(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) OnChainConfig21(arg0 OnchainConfigLegacy) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.OnChainConfig21(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) OnChainConfig22Plus(opts *bind.TransactOpts, arg0 OnchainConfig) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_onChainConfig22Plus", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) OnChainConfig22Plus(arg0 OnchainConfig) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.OnChainConfig22Plus(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) OnChainConfig22Plus(arg0 OnchainConfig) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.OnChainConfig22Plus(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactor) Report(opts *bind.TransactOpts, arg0 Report) (*types.Transaction, error) {
	return _AutomationConvenience.contract.Transact(opts, "_report", arg0)
}

func (_AutomationConvenience *AutomationConvenienceSession) Report(arg0 Report) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.Report(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenienceTransactorSession) Report(arg0 Report) (*types.Transaction, error) {
	return _AutomationConvenience.Contract.Report(&_AutomationConvenience.TransactOpts, arg0)
}

func (_AutomationConvenience *AutomationConvenience) Address() common.Address {
	return _AutomationConvenience.address
}

type AutomationConvenienceInterface interface {
	ConditionalTrigger(opts *bind.TransactOpts, arg0 ConditionalTrigger) (*types.Transaction, error)

	Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error)

	LogTrigger(opts *bind.TransactOpts, arg0 LogTrigger) (*types.Transaction, error)

	LogTriggerConfig(opts *bind.TransactOpts, arg0 LogTriggerConfig) (*types.Transaction, error)

	OnChainConfig21(opts *bind.TransactOpts, arg0 OnchainConfigLegacy) (*types.Transaction, error)

	OnChainConfig22Plus(opts *bind.TransactOpts, arg0 OnchainConfig) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, arg0 Report) (*types.Transaction, error)

	Address() common.Address
}
