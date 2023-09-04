// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_utils_2_1

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

type KeeperRegistryBase21ConditionalTrigger struct {
	BlockNum  uint32
	BlockHash [32]byte
}

type KeeperRegistryBase21LogTrigger struct {
	LogBlockHash [32]byte
	TxHash       [32]byte
	LogIndex     uint32
	BlockNum     uint32
	BlockHash    [32]byte
}

type KeeperRegistryBase21OnchainConfig struct {
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

type KeeperRegistryBase21Report struct {
	FastGasWei   *big.Int
	LinkNative   *big.Int
	UpkeepIds    []*big.Int
	GasLimits    []*big.Int
	Triggers     [][]byte
	PerformDatas [][]byte
}

type Log struct {
	Index       *big.Int
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structKeeperRegistryBase2_1.ConditionalTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_conditionalTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_log\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"logBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"logIndex\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNum\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structKeeperRegistryBase2_1.LogTrigger\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structLogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_logTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structKeeperRegistryBase2_1.OnchainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_onChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"gasLimits\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"triggers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"performDatas\",\"type\":\"bytes[]\"}],\"internalType\":\"structKeeperRegistryBase2_1.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506108be806100206000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063776f306111610050578063776f3061146100a6578063d15a060e146100b4578063e65d6546146100c257600080fd5b806321f373d7146100775780632ff92a811461008a5780634b6df29414610098575b600080fd5b6100886100853660046101e7565b50565b005b610088610085366004610362565b6100886100853660046104bc565b610088610085366004610513565b61008861008536600461067a565b6100886100853660046107c4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516101e0810167ffffffffffffffff81118282101715610123576101236100d0565b60405290565b60405160e0810167ffffffffffffffff81118282101715610123576101236100d0565b60405160c0810167ffffffffffffffff81118282101715610123576101236100d0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101b6576101b66100d0565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146101e257600080fd5b919050565b600060c082840312156101f957600080fd5b60405160c0810181811067ffffffffffffffff8211171561021c5761021c6100d0565b604052610228836101be565b8152602083013560ff8116811461023e57600080fd5b8060208301525060408301356040820152606083013560608201526080830135608082015260a083013560a08201528091505092915050565b803563ffffffff811681146101e257600080fd5b803562ffffff811681146101e257600080fd5b803561ffff811681146101e257600080fd5b80356bffffffffffffffffffffffff811681146101e257600080fd5b600067ffffffffffffffff8211156102e6576102e66100d0565b5060051b60200190565b600082601f83011261030157600080fd5b81356020610316610311836102cc565b61016f565b82815260059290921b8401810191818101908684111561033557600080fd5b8286015b848110156103575761034a816101be565b8352918301918301610339565b509695505050505050565b60006020828403121561037457600080fd5b813567ffffffffffffffff8082111561038c57600080fd5b908301906101e082860312156103a157600080fd5b6103a96100ff565b6103b283610277565b81526103c060208401610277565b60208201526103d160408401610277565b60408201526103e26060840161028b565b60608201526103f36080840161029e565b608082015261040460a084016102b0565b60a082015261041560c08401610277565b60c082015261042660e08401610277565b60e0820152610100610439818501610277565b9082015261012061044b848201610277565b90820152610140838101359082015261016080840135908201526101806104738185016101be565b908201526101a0838101358381111561048b57600080fd5b610497888287016102f0565b8284015250506101c091506104ad8284016101be565b91810191909152949350505050565b6000604082840312156104ce57600080fd5b6040516040810181811067ffffffffffffffff821117156104f1576104f16100d0565b6040526104fd83610277565b8152602083013560208201528091505092915050565b600060a0828403121561052557600080fd5b60405160a0810181811067ffffffffffffffff82111715610548576105486100d0565b8060405250823581526020830135602082015261056760408401610277565b604082015261057860608401610277565b6060820152608083013560808201528091505092915050565b600082601f8301126105a257600080fd5b813560206105b2610311836102cc565b82815260059290921b840181019181810190868411156105d157600080fd5b8286015b8481101561035757803583529183019183016105d5565b600082601f8301126105fd57600080fd5b813567ffffffffffffffff811115610617576106176100d0565b61064860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161016f565b81815284602083860101111561065d57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561068c57600080fd5b813567ffffffffffffffff808211156106a457600080fd5b9083019060e082860312156106b857600080fd5b6106c0610129565b823581526020830135602082015260408301356040820152606083013560608201526106ee608084016101be565b608082015260a08301358281111561070557600080fd5b61071187828601610591565b60a08301525060c08301358281111561072957600080fd5b610735878286016105ec565b60c08301525095945050505050565b600082601f83011261075557600080fd5b81356020610765610311836102cc565b82815260059290921b8401810191818101908684111561078457600080fd5b8286015b8481101561035757803567ffffffffffffffff8111156107a85760008081fd5b6107b68986838b01016105ec565b845250918301918301610788565b6000602082840312156107d657600080fd5b813567ffffffffffffffff808211156107ee57600080fd5b9083019060c0828603121561080257600080fd5b61080a61014c565b823581526020830135602082015260408301358281111561082a57600080fd5b61083687828601610591565b60408301525060608301358281111561084e57600080fd5b61085a87828601610591565b60608301525060808301358281111561087257600080fd5b61087e87828601610744565b60808301525060a08301358281111561089657600080fd5b6108a287828601610744565b60a0830152509594505050505056fea164736f6c6343000810000a",
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
	return address, tx, &AutomationUtils{AutomationUtilsCaller: AutomationUtilsCaller{contract: contract}, AutomationUtilsTransactor: AutomationUtilsTransactor{contract: contract}, AutomationUtilsFilterer: AutomationUtilsFilterer{contract: contract}}, nil
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

func (_AutomationUtils *AutomationUtilsTransactor) ConditionalTrigger(opts *bind.TransactOpts, arg0 KeeperRegistryBase21ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_conditionalTrigger", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) ConditionalTrigger(arg0 KeeperRegistryBase21ConditionalTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.ConditionalTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) ConditionalTrigger(arg0 KeeperRegistryBase21ConditionalTrigger) (*types.Transaction, error) {
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

func (_AutomationUtils *AutomationUtilsTransactor) LogTrigger(opts *bind.TransactOpts, arg0 KeeperRegistryBase21LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_logTrigger", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) LogTrigger(arg0 KeeperRegistryBase21LogTrigger) (*types.Transaction, error) {
	return _AutomationUtils.Contract.LogTrigger(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) LogTrigger(arg0 KeeperRegistryBase21LogTrigger) (*types.Transaction, error) {
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

func (_AutomationUtils *AutomationUtilsTransactor) OnChainConfig(opts *bind.TransactOpts, arg0 KeeperRegistryBase21OnchainConfig) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_onChainConfig", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) OnChainConfig(arg0 KeeperRegistryBase21OnchainConfig) (*types.Transaction, error) {
	return _AutomationUtils.Contract.OnChainConfig(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) OnChainConfig(arg0 KeeperRegistryBase21OnchainConfig) (*types.Transaction, error) {
	return _AutomationUtils.Contract.OnChainConfig(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactor) Report(opts *bind.TransactOpts, arg0 KeeperRegistryBase21Report) (*types.Transaction, error) {
	return _AutomationUtils.contract.Transact(opts, "_report", arg0)
}

func (_AutomationUtils *AutomationUtilsSession) Report(arg0 KeeperRegistryBase21Report) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Report(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtilsTransactorSession) Report(arg0 KeeperRegistryBase21Report) (*types.Transaction, error) {
	return _AutomationUtils.Contract.Report(&_AutomationUtils.TransactOpts, arg0)
}

func (_AutomationUtils *AutomationUtils) Address() common.Address {
	return _AutomationUtils.address
}

type AutomationUtilsInterface interface {
	ConditionalTrigger(opts *bind.TransactOpts, arg0 KeeperRegistryBase21ConditionalTrigger) (*types.Transaction, error)

	Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error)

	LogTrigger(opts *bind.TransactOpts, arg0 KeeperRegistryBase21LogTrigger) (*types.Transaction, error)

	LogTriggerConfig(opts *bind.TransactOpts, arg0 LogTriggerConfig) (*types.Transaction, error)

	OnChainConfig(opts *bind.TransactOpts, arg0 KeeperRegistryBase21OnchainConfig) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, arg0 KeeperRegistryBase21Report) (*types.Transaction, error)

	Address() common.Address
}
