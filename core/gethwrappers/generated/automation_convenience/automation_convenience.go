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
	Bin: "0x608060405234801561001057600080fd5b50610a79806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063785826fb1161005b578063785826fb146100b1578063e0c4915e146100bf578063e65d6546146100cd578063e9720a49146100db57600080fd5b806321f373d7146100825780634b6df29414610095578063776f3061146100a3575b600080fd5b61009361009036600461022e565b50565b005b6100936100903660046102b6565b61009361009036600461030d565b610093610090366004610474565b6100936100903660046105f6565b6100936100903660046108aa565b610093610090366004610997565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff8111828210171561013b5761013b6100e9565b60405290565b604051610220810167ffffffffffffffff8111828210171561013b5761013b6100e9565b6040516101e0810167ffffffffffffffff8111828210171561013b5761013b6100e9565b604051610100810167ffffffffffffffff8111828210171561013b5761013b6100e9565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101f4576101f46100e9565b604052919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461009057600080fd5b8035610229816101fc565b919050565b600060c0828403121561024057600080fd5b610248610118565b8235610253816101fc565b8152602083013560ff8116811461026957600080fd5b8060208301525060408301356040820152606083013560608201526080830135608082015260a083013560a08201528091505092915050565b803563ffffffff8116811461022957600080fd5b6000604082840312156102c857600080fd5b6040516040810181811067ffffffffffffffff821117156102eb576102eb6100e9565b6040526102f7836102a2565b8152602083013560208201528091505092915050565b600060a0828403121561031f57600080fd5b60405160a0810181811067ffffffffffffffff82111715610342576103426100e9565b80604052508235815260208301356020820152610361604084016102a2565b6040820152610372606084016102a2565b6060820152608083013560808201528091505092915050565b803562ffffff8116811461022957600080fd5b803561ffff8116811461022957600080fd5b80356bffffffffffffffffffffffff8116811461022957600080fd5b600067ffffffffffffffff8211156103e6576103e66100e9565b5060051b60200190565b600082601f83011261040157600080fd5b81356020610416610411836103cc565b6101ad565b82815260059290921b8401810191818101908684111561043557600080fd5b8286015b8481101561045957803561044c816101fc565b8352918301918301610439565b509695505050505050565b8035801515811461022957600080fd5b60006020828403121561048657600080fd5b813567ffffffffffffffff8082111561049e57600080fd5b9083019061022082860312156104b357600080fd5b6104bb610141565b6104c4836102a2565b81526104d2602084016102a2565b60208201526104e3604084016102a2565b60408201526104f46060840161038b565b60608201526105056080840161039e565b608082015261051660a084016103b0565b60a082015261052760c084016102a2565b60c082015261053860e084016102a2565b60e082015261010061054b8185016102a2565b9082015261012061055d8482016102a2565b908201526101408381013590820152610160808401359082015261018061058581850161021e565b908201526101a0838101358381111561059d57600080fd5b6105a9888287016103f0565b8284015250506101c091506105bf82840161021e565b828201526101e091506105d382840161021e565b8282015261020091506105e7828401610464565b91810191909152949350505050565b60006020828403121561060857600080fd5b813567ffffffffffffffff8082111561062057600080fd5b908301906101e0828603121561063557600080fd5b61063d610165565b610646836102a2565b8152610654602084016102a2565b6020820152610665604084016102a2565b60408201526106766060840161038b565b60608201526106876080840161039e565b608082015261069860a084016103b0565b60a08201526106a960c084016102a2565b60c08201526106ba60e084016102a2565b60e08201526101006106cd8185016102a2565b908201526101206106df8482016102a2565b908201526101408381013590820152610160808401359082015261018061070781850161021e565b908201526101a0838101358381111561071f57600080fd5b61072b888287016103f0565b8284015250506101c091506105e782840161021e565b600082601f83011261075257600080fd5b81356020610762610411836103cc565b82815260059290921b8401810191818101908684111561078157600080fd5b8286015b848110156104595780358352918301918301610785565b600082601f8301126107ad57600080fd5b813567ffffffffffffffff8111156107c7576107c76100e9565b6107f860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016101ad565b81815284602083860101111561080d57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261083b57600080fd5b8135602061084b610411836103cc565b82815260059290921b8401810191818101908684111561086a57600080fd5b8286015b8481101561045957803567ffffffffffffffff81111561088e5760008081fd5b61089c8986838b010161079c565b84525091830191830161086e565b6000602082840312156108bc57600080fd5b813567ffffffffffffffff808211156108d457600080fd5b9083019060c082860312156108e857600080fd5b6108f0610118565b823581526020830135602082015260408301358281111561091057600080fd5b61091c87828601610741565b60408301525060608301358281111561093457600080fd5b61094087828601610741565b60608301525060808301358281111561095857600080fd5b6109648782860161082a565b60808301525060a08301358281111561097c57600080fd5b6109888782860161082a565b60a08301525095945050505050565b6000602082840312156109a957600080fd5b813567ffffffffffffffff808211156109c157600080fd5b9083019061010082860312156109d657600080fd5b6109de610189565b8235815260208301356020820152604083013560408201526060830135606082015260808301356080820152610a1660a0840161021e565b60a082015260c083013582811115610a2d57600080fd5b610a3987828601610741565b60c08301525060e083013582811115610a5157600080fd5b610a5d8782860161079c565b60e0830152509594505050505056fea164736f6c6343000813000a",
}

var AutomationConvenienceABI = AutomationConvenienceMetaData.ABI

var AutomationConvenienceBin = AutomationConvenienceMetaData.Bin

func DeployAutomationConvenience(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AutomationConvenience, error) {
	parsed, err := AutomationConvenienceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationConvenienceBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationConvenience{address: address, abi: *parsed, AutomationConvenienceCaller: AutomationConvenienceCaller{contract: contract}, AutomationConvenienceTransactor: AutomationConvenienceTransactor{contract: contract}, AutomationConvenienceFilterer: AutomationConvenienceFilterer{contract: contract}}, nil
}

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
