// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_utils_2_3

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

type AutomationRegistryBase23BillingConfig struct {
	Active           bool
	GasFeePPB        uint32
	FlatFeeMicroLink *big.Int
	PriceFeed        common.Address
}

type AutomationRegistryBase23ConditionalTrigger struct {
	BlockNum  uint32
	BlockHash [32]byte
}

type AutomationRegistryBase23LogTrigger struct {
	LogBlockHash [32]byte
	TxHash       [32]byte
	LogIndex     uint32
	BlockNum     uint32
	BlockHash    [32]byte
}

type AutomationRegistryBase23OnchainConfig struct {
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

type AutomationRegistryBase23Report struct {
	FastGasWei   *big.Int
	LinkNative   *big.Int
	UpkeepIds    []*big.Int
	GasLimits    []*big.Int
	Triggers     [][]byte
	PerformDatas [][]byte
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structAutomationRegistryBase2_3.ConditionalTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_conditionalTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_log\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"logBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"logIndex\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structAutomationRegistryBase2_3.LogTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structLogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"priceFeed\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"_onChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"gasLimits\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"triggers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"performDatas\",\"type\":\"bytes[]\"}],\"internalType\":\"structAutomationRegistryBase2_3.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610a0b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063776f306111610050578063776f3061146100ab578063e65d6546146100b9578063e9720a49146100c757600080fd5b806321f373d7146100775780634b6df2941461008a578063568aee5014610098575b600080fd5b610088610085366004610219565b50565b005b6100886100853660046102a1565b6100886100a6366004610492565b505050565b610088610085366004610655565b61008861008536600461083c565b610088610085366004610929565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff81118282101715610127576101276100d5565b60405290565b6040516080810167ffffffffffffffff81118282101715610127576101276100d5565b604051610220810167ffffffffffffffff81118282101715610127576101276100d5565b604051610100810167ffffffffffffffff81118282101715610127576101276100d5565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101df576101df6100d5565b604052919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461008557600080fd5b8035610214816101e7565b919050565b600060c0828403121561022b57600080fd5b610233610104565b823561023e816101e7565b8152602083013560ff8116811461025457600080fd5b8060208301525060408301356040820152606083013560608201526080830135608082015260a083013560a08201528091505092915050565b803563ffffffff8116811461021457600080fd5b6000604082840312156102b357600080fd5b6040516040810181811067ffffffffffffffff821117156102d6576102d66100d5565b6040526102e28361028d565b8152602083013560208201528091505092915050565b803562ffffff8116811461021457600080fd5b803561ffff8116811461021457600080fd5b80356bffffffffffffffffffffffff8116811461021457600080fd5b600067ffffffffffffffff821115610353576103536100d5565b5060051b60200190565b600082601f83011261036e57600080fd5b8135602061038361037e83610339565b610198565b82815260059290921b840181019181810190868411156103a257600080fd5b8286015b848110156103c65780356103b9816101e7565b83529183019183016103a6565b509695505050505050565b8035801515811461021457600080fd5b600082601f8301126103f257600080fd5b8135602061040261037e83610339565b82815260079290921b8401810191818101908684111561042157600080fd5b8286015b848110156103c6576080818903121561043e5760008081fd5b61044661012d565b61044f826103d1565b815261045c85830161028d565b85820152604061046d8184016102f8565b90820152606082810135610480816101e7565b90820152835291830191608001610425565b6000806000606084860312156104a757600080fd5b833567ffffffffffffffff808211156104bf57600080fd5b9085019061022082880312156104d457600080fd5b6104dc610150565b6104e58361028d565b81526104f36020840161028d565b60208201526105046040840161028d565b6040820152610515606084016102f8565b60608201526105266080840161030b565b608082015261053760a0840161031d565b60a082015261054860c0840161028d565b60c082015261055960e0840161028d565b60e082015261010061056c81850161028d565b9082015261012061057e84820161028d565b90820152610140838101359082015261016080840135908201526101806105a6818501610209565b908201526101a083810135838111156105be57600080fd5b6105ca8a82870161035d565b8284015250506101c06105de818501610209565b908201526101e06105f0848201610209565b908201526102006106028482016103d1565b908201529450602086013591508082111561061c57600080fd5b6106288783880161035d565b9350604086013591508082111561063e57600080fd5b5061064b868287016103e1565b9150509250925092565b600060a0828403121561066757600080fd5b60405160a0810181811067ffffffffffffffff8211171561068a5761068a6100d5565b806040525082358152602083013560208201526106a96040840161028d565b60408201526106ba6060840161028d565b6060820152608083013560808201528091505092915050565b600082601f8301126106e457600080fd5b813560206106f461037e83610339565b82815260059290921b8401810191818101908684111561071357600080fd5b8286015b848110156103c65780358352918301918301610717565b600082601f83011261073f57600080fd5b813567ffffffffffffffff811115610759576107596100d5565b61078a60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610198565b81815284602083860101111561079f57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126107cd57600080fd5b813560206107dd61037e83610339565b82815260059290921b840181019181810190868411156107fc57600080fd5b8286015b848110156103c657803567ffffffffffffffff8111156108205760008081fd5b61082e8986838b010161072e565b845250918301918301610800565b60006020828403121561084e57600080fd5b813567ffffffffffffffff8082111561086657600080fd5b9083019060c0828603121561087a57600080fd5b610882610104565b82358152602083013560208201526040830135828111156108a257600080fd5b6108ae878286016106d3565b6040830152506060830135828111156108c657600080fd5b6108d2878286016106d3565b6060830152506080830135828111156108ea57600080fd5b6108f6878286016107bc565b60808301525060a08301358281111561090e57600080fd5b61091a878286016107bc565b60a08301525095945050505050565b60006020828403121561093b57600080fd5b813567ffffffffffffffff8082111561095357600080fd5b90830190610100828603121561096857600080fd5b610970610174565b82358152602083013560208201526040830135604082015260608301356060820152608083013560808201526109a860a08401610209565b60a082015260c0830135828111156109bf57600080fd5b6109cb878286016106d3565b60c08301525060e0830135828111156109e357600080fd5b6109ef8782860161072e565b60e0830152509594505050505056fea164736f6c6343000813000a",
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

func (_AutomationUtils *AutomationUtilsTransactor) ConditionalTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase23ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_conditionalTrigger", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) ConditionalTrigger(arg0 AutomationRegistryBase23ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.ConditionalTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) ConditionalTrigger(arg0 AutomationRegistryBase23ConditionalTrigger) (*types.Transaction, error) {
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

func (_AutomationUtils *AutomationUtilsTransactor) LogTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase23LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_logTrigger", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) LogTrigger(arg0 AutomationRegistryBase23LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.LogTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) LogTrigger(arg0 AutomationRegistryBase23LogTrigger) (*types.Transaction, error) {
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

func (_AutomationUtils *AutomationUtilsTransactor) OnChainConfig(opts *bind.TransactOpts, arg0 AutomationRegistryBase23OnchainConfig, arg1 []common.Address, arg2 []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_onChainConfig", arg0, arg1, arg2)
}

func (_AutomationUtils *AutomationUtilsSession) OnChainConfig(arg0 AutomationRegistryBase23OnchainConfig, arg1 []common.Address, arg2 []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationUtils.Contract.OnChainConfig(&_AutomationUtils.TransactOpts, arg0, arg1, arg2)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) OnChainConfig(arg0 AutomationRegistryBase23OnchainConfig, arg1 []common.Address, arg2 []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationUtils.Contract.OnChainConfig(&_AutomationUtils.TransactOpts, arg0, arg1, arg2)
}

func (_AutomationUtils *AutomationUtilsTransactor) Report(opts *bind.TransactOpts, arg0 AutomationRegistryBase23Report) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_report", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) Report(arg0 AutomationRegistryBase23Report) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Report(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) Report(arg0 AutomationRegistryBase23Report) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Report(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtils) Address() common.Address {
	return _AutomationUtils.address
}

type AutomationUtilsInterface interface {
	ConditionalTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase23ConditionalTrigger) (*types.Transaction, error)

	Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error)

	LogTrigger(opts *bind.TransactOpts, arg0 AutomationRegistryBase23LogTrigger) (*types.Transaction, error)

	LogTriggerConfig(opts *bind.TransactOpts, arg0 LogTriggerConfig) (*types.Transaction, error)

	OnChainConfig(opts *bind.TransactOpts, arg0 AutomationRegistryBase23OnchainConfig, arg1 []common.Address, arg2 []AutomationRegistryBase23BillingConfig) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, arg0 AutomationRegistryBase23Report) (*types.Transaction, error)

	Address() common.Address
}
