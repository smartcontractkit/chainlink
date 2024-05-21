// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_utils_2_2

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

type AutomationRegistryBase22ConditionalTrigger struct {
	BlockNum  uint32
	BlockHash [32]byte
}

type AutomationRegistryBase22LogTrigger struct {
	LogBlockHash [32]byte
	TxHash       [32]byte
	LogIndex     uint32
	BlockNum     uint32
	BlockHash    [32]byte
}

type AutomationRegistryBase22Report struct {
	FastGasWei   *big.Int
	LinkNative   *big.Int
	UpkeepIds    []*big.Int
	GasLimits    []*big.Int
	Triggers     [][]byte
	PerformDatas [][]byte
}

type IAutomationV21PlusCommonOnchainConfigLegacy struct {
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

type LogTriggerConfig struct {
	ContractAddress common.Address
	FilterSelector  uint8
	Topic0          [32]byte
	Topic1          [32]byte
	Topic2          [32]byte
	Topic3          [32]byte
}

var AutomationUtilsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structAutomationRegistryBase2_2.ConditionalTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_conditionalTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_log\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"logBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"logIndex\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structAutomationRegistryBase2_2.LogTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structLogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_onChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"gasLimits\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"triggers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"performDatas\",\"type\":\"bytes[]\"}],\"internalType\":\"structAutomationRegistryBase2_2.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506108ac806100206000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063776f306111610050578063776f3061146100a6578063e65d6546146100b4578063e9720a49146100c257600080fd5b806321f373d7146100775780632ff92a811461008a5780634b6df29414610098575b600080fd5b6100886100853660046101e8565b50565b005b610088610085366004610345565b61008861008536600461049f565b6100886100853660046104f6565b6100886100853660046106dd565b6100886100853660046107ca565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff81118282101715610122576101226100d0565b60405290565b6040516101e0810167ffffffffffffffff81118282101715610122576101226100d0565b604051610100810167ffffffffffffffff81118282101715610122576101226100d0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101b7576101b76100d0565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146101e357600080fd5b919050565b600060c082840312156101fa57600080fd5b6102026100ff565b61020b836101bf565b8152602083013560ff8116811461022157600080fd5b8060208301525060408301356040820152606083013560608201526080830135608082015260a083013560a08201528091505092915050565b803563ffffffff811681146101e357600080fd5b803562ffffff811681146101e357600080fd5b803561ffff811681146101e357600080fd5b80356bffffffffffffffffffffffff811681146101e357600080fd5b600067ffffffffffffffff8211156102c9576102c96100d0565b5060051b60200190565b600082601f8301126102e457600080fd5b813560206102f96102f4836102af565b610170565b82815260059290921b8401810191818101908684111561031857600080fd5b8286015b8481101561033a5761032d816101bf565b835291830191830161031c565b509695505050505050565b60006020828403121561035757600080fd5b813567ffffffffffffffff8082111561036f57600080fd5b908301906101e0828603121561038457600080fd5b61038c610128565b6103958361025a565b81526103a36020840161025a565b60208201526103b46040840161025a565b60408201526103c56060840161026e565b60608201526103d660808401610281565b60808201526103e760a08401610293565b60a08201526103f860c0840161025a565b60c082015261040960e0840161025a565b60e082015261010061041c81850161025a565b9082015261012061042e84820161025a565b90820152610140838101359082015261016080840135908201526101806104568185016101bf565b908201526101a0838101358381111561046e57600080fd5b61047a888287016102d3565b8284015250506101c091506104908284016101bf565b91810191909152949350505050565b6000604082840312156104b157600080fd5b6040516040810181811067ffffffffffffffff821117156104d4576104d46100d0565b6040526104e08361025a565b8152602083013560208201528091505092915050565b600060a0828403121561050857600080fd5b60405160a0810181811067ffffffffffffffff8211171561052b5761052b6100d0565b8060405250823581526020830135602082015261054a6040840161025a565b604082015261055b6060840161025a565b6060820152608083013560808201528091505092915050565b600082601f83011261058557600080fd5b813560206105956102f4836102af565b82815260059290921b840181019181810190868411156105b457600080fd5b8286015b8481101561033a57803583529183019183016105b8565b600082601f8301126105e057600080fd5b813567ffffffffffffffff8111156105fa576105fa6100d0565b61062b60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610170565b81815284602083860101111561064057600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261066e57600080fd5b8135602061067e6102f4836102af565b82815260059290921b8401810191818101908684111561069d57600080fd5b8286015b8481101561033a57803567ffffffffffffffff8111156106c15760008081fd5b6106cf8986838b01016105cf565b8452509183019183016106a1565b6000602082840312156106ef57600080fd5b813567ffffffffffffffff8082111561070757600080fd5b9083019060c0828603121561071b57600080fd5b6107236100ff565b823581526020830135602082015260408301358281111561074357600080fd5b61074f87828601610574565b60408301525060608301358281111561076757600080fd5b61077387828601610574565b60608301525060808301358281111561078b57600080fd5b6107978782860161065d565b60808301525060a0830135828111156107af57600080fd5b6107bb8782860161065d565b60a08301525095945050505050565b6000602082840312156107dc57600080fd5b813567ffffffffffffffff808211156107f457600080fd5b90830190610100828603121561080957600080fd5b61081161014c565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261084960a084016101bf565b60a082015260c08301358281111561086057600080fd5b61086c87828601610574565b60c08301525060e08301358281111561088457600080fd5b610890878286016105cf565b60e0830152509594505050505056fea164736f6c6343000813000a",
}

var AutomationUtilsABI = AutomationUtilsMetaData.ABI

var AutomationUtilsBin = AutomationUtilsMetaData.Bin

func DeployAutomationUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AutomationUtils, error) {
	parsed, err := AutomationUtilsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationUtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationUtils{address: address, abi: *parsed, AutomationUtilsCaller: AutomationUtilsCaller{contract: contract}, AutomationUtilsTransactor: AutomationUtilsTransactor{contract: contract}, AutomationUtilsFilterer: AutomationUtilsFilterer{contract: contract}}, nil
}

type AutomationUtils struct {
	address common.Address
	abi     abi.ABI
	AutomationUtilsCaller
	AutomationUtilsTransactor
	AutomationUtilsFilterer
}

type AutomationUtilsCaller struct {
	contract *bind.BoundContract
}

type AutomationUtilsTransactor struct {
	contract *bind.BoundContract
}

type AutomationUtilsFilterer struct {
	contract *bind.BoundContract
}

type AutomationUtilsSession struct {
	Contract     *AutomationUtils
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationUtilsCallerSession struct {
	Contract *AutomationUtilsCaller
	CallOpts bind.CallOpts
}

type AutomationUtilsTransactorSession struct {
	Contract     *AutomationUtilsTransactor
	TransactOpts bind.TransactOpts
}

type AutomationUtilsRaw struct {
	Contract *AutomationUtils
}

type AutomationUtilsCallerRaw struct {
	Contract *AutomationUtilsCaller
}

type AutomationUtilsTransactorRaw struct {
	Contract *AutomationUtilsTransactor
}

func NewAutomationUtils(address common.Address, backend bind.ContractBackend) (*AutomationUtils, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationUtilsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationUtils{address: address, abi: abi, AutomationUtilsCaller: AutomationUtilsCaller{contract: contract}, AutomationUtilsTransactor: AutomationUtilsTransactor{contract: contract}, AutomationUtilsFilterer: AutomationUtilsFilterer{contract: contract}}, nil
}

func NewAutomationUtilsCaller(address common.Address, caller bind.ContractCaller) (*AutomationUtilsCaller, error) {
	contract, err := bindAutomationUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationUtilsCaller{contract: contract}, nil
}

func NewAutomationUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationUtilsTransactor, error) {
	contract, err := bindAutomationUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationUtilsTransactor{contract: contract}, nil
}

func NewAutomationUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationUtilsFilterer, error) {
	contract, err := bindAutomationUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationUtilsFilterer{contract: contract}, nil
}

func bindAutomationUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationUtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationUtils *AutomationUtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationUtils.Contract.AutomationUtilsCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationUtils *AutomationUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationUtils.Contract.AutomationUtilsTransactor.contract.Transfer(opts)
}

func (_AutomationUtils *AutomationUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationUtils.Contract.AutomationUtilsTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationUtils *AutomationUtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationUtils.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationUtils *AutomationUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationUtils.Contract.contract.Transfer(opts)
}

func (_AutomationUtils *AutomationUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationUtils.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationUtils *AutomationUtilsTransactor) ConditionalTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase22ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_conditionalTrigger", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) ConditionalTrigger(arg0 AutomationRegistryBase22ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.ConditionalTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) ConditionalTrigger(arg0 AutomationRegistryBase22ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.ConditionalTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactor) Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_log", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) Log(arg0 Log) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Log(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) Log(arg0 Log) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Log(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactor) LogTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase22LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_logTrigger", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) LogTrigger(arg0 AutomationRegistryBase22LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.LogTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) LogTrigger(arg0 AutomationRegistryBase22LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.LogTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactor) LogTriggerConfig(opts *bind.TransactOpts, arg0 LogTriggerConfig) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_logTriggerConfig", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) LogTriggerConfig(arg0 LogTriggerConfig) (*types.Transaction, error) {
	return _AutomationUtils.Contract.LogTriggerConfig(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) LogTriggerConfig(arg0 LogTriggerConfig) (*types.Transaction, error) {
	return _AutomationUtils.Contract.LogTriggerConfig(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactor) OnChainConfig(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonOnchainConfigLegacy) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_onChainConfig", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) OnChainConfig(arg0 IAutomationV21PlusCommonOnchainConfigLegacy) (*types.Transaction, error) {
	return _AutomationUtils.Contract.OnChainConfig(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) OnChainConfig(arg0 IAutomationV21PlusCommonOnchainConfigLegacy) (*types.Transaction, error) {
	return _AutomationUtils.Contract.OnChainConfig(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactor) Report(opts *bind.TransactOpts, arg0 AutomationRegistryBase22Report) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_report", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) Report(arg0 AutomationRegistryBase22Report) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Report(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) Report(arg0 AutomationRegistryBase22Report) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Report(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtils) Address() common.Address {
	return _AutomationUtils.address
}

type AutomationUtilsInterface interface {
	ConditionalTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase22ConditionalTrigger) (*types.Transaction, error)

	Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error)

	LogTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase22LogTrigger) (*types.Transaction, error)

	LogTriggerConfig(opts *bind.TransactOpts, arg0 LogTriggerConfig) (*types.Transaction, error)

	OnChainConfig(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonOnchainConfigLegacy) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, arg0 AutomationRegistryBase22Report) (*types.Transaction, error)

	Address() common.Address
}
