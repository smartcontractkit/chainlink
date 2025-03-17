// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_compatible_utils

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

type IAutomationV21PlusCommonConditionalTrigger struct {
	BlockNum  uint32
	BlockHash [32]byte
}

type IAutomationV21PlusCommonLogTrigger struct {
	LogBlockHash [32]byte
	TxHash       [32]byte
	LogIndex     uint32
	BlockNum     uint32
	BlockHash    [32]byte
}

type IAutomationV21PlusCommonLogTriggerConfig struct {
	ContractAddress common.Address
	FilterSelector  uint8
	Topic0          [32]byte
	Topic1          [32]byte
	Topic2          [32]byte
	Topic3          [32]byte
}

type IAutomationV21PlusCommonReport struct {
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

var AutomationCompatibleUtilsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"_conditionalTrigger\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.ConditionalTrigger\",\"components\":[{\"name\":\"blockNum\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_log\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structLog\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"source\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"topics\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_logTrigger\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.LogTrigger\",\"components\":[{\"name\":\"logBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"logIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockNum\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_logTriggerConfig\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.LogTriggerConfig\",\"components\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"filterSelector\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"topic0\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"topic1\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"topic2\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"topic3\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_report\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.Report\",\"components\":[{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"upkeepIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"gasLimits\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"triggers\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"performDatas\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610672806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c8063776f306111610050578063776f30611461008d578063e65d65461461009b578063e9720a49146100a957600080fd5b806321f373d71461006c5780634b6df2941461007f575b600080fd5b61007d61007a3660046101ab565b50565b005b61007d61007a366004610231565b61007d61007a366004610288565b61007d61007a3660046104a3565b61007d61007a366004610590565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff81118282101715610109576101096100b7565b60405290565b604051610100810167ffffffffffffffff81118282101715610109576101096100b7565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561017a5761017a6100b7565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146101a657600080fd5b919050565b600060c082840312156101bd57600080fd5b6101c56100e6565b6101ce83610182565b8152602083013560ff811681146101e457600080fd5b8060208301525060408301356040820152606083013560608201526080830135608082015260a083013560a08201528091505092915050565b803563ffffffff811681146101a657600080fd5b60006040828403121561024357600080fd5b6040516040810181811067ffffffffffffffff82111715610266576102666100b7565b6040526102728361021d565b8152602083013560208201528091505092915050565b600060a0828403121561029a57600080fd5b60405160a0810181811067ffffffffffffffff821117156102bd576102bd6100b7565b806040525082358152602083013560208201526102dc6040840161021d565b60408201526102ed6060840161021d565b6060820152608083013560808201528091505092915050565b600067ffffffffffffffff821115610320576103206100b7565b5060051b60200190565b600082601f83011261033b57600080fd5b8135602061035061034b83610306565b610133565b82815260059290921b8401810191818101908684111561036f57600080fd5b8286015b8481101561038a5780358352918301918301610373565b509695505050505050565b600082601f8301126103a657600080fd5b813567ffffffffffffffff8111156103c0576103c06100b7565b6103f160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610133565b81815284602083860101111561040657600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261043457600080fd5b8135602061044461034b83610306565b82815260059290921b8401810191818101908684111561046357600080fd5b8286015b8481101561038a57803567ffffffffffffffff8111156104875760008081fd5b6104958986838b0101610395565b845250918301918301610467565b6000602082840312156104b557600080fd5b813567ffffffffffffffff808211156104cd57600080fd5b9083019060c082860312156104e157600080fd5b6104e96100e6565b823581526020830135602082015260408301358281111561050957600080fd5b6105158782860161032a565b60408301525060608301358281111561052d57600080fd5b6105398782860161032a565b60608301525060808301358281111561055157600080fd5b61055d87828601610423565b60808301525060a08301358281111561057557600080fd5b61058187828601610423565b60a08301525095945050505050565b6000602082840312156105a257600080fd5b813567ffffffffffffffff808211156105ba57600080fd5b9083019061010082860312156105cf57600080fd5b6105d761010f565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261060f60a08401610182565b60a082015260c08301358281111561062657600080fd5b6106328782860161032a565b60c08301525060e08301358281111561064a57600080fd5b61065687828601610395565b60e0830152509594505050505056fea164736f6c6343000813000a",
}

var AutomationCompatibleUtilsABI = AutomationCompatibleUtilsMetaData.ABI

var AutomationCompatibleUtilsBin = AutomationCompatibleUtilsMetaData.Bin

func DeployAutomationCompatibleUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AutomationCompatibleUtils, error) {
	parsed, err := AutomationCompatibleUtilsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationCompatibleUtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationCompatibleUtils{address: address, abi: *parsed, AutomationCompatibleUtilsCaller: AutomationCompatibleUtilsCaller{contract: contract}, AutomationCompatibleUtilsTransactor: AutomationCompatibleUtilsTransactor{contract: contract}, AutomationCompatibleUtilsFilterer: AutomationCompatibleUtilsFilterer{contract: contract}}, nil
}

type AutomationCompatibleUtils struct {
	address common.Address
	abi     abi.ABI
	AutomationCompatibleUtilsCaller
	AutomationCompatibleUtilsTransactor
	AutomationCompatibleUtilsFilterer
}

type AutomationCompatibleUtilsCaller struct {
	contract *bind.BoundContract
}

type AutomationCompatibleUtilsTransactor struct {
	contract *bind.BoundContract
}

type AutomationCompatibleUtilsFilterer struct {
	contract *bind.BoundContract
}

type AutomationCompatibleUtilsSession struct {
	Contract     *AutomationCompatibleUtils
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationCompatibleUtilsCallerSession struct {
	Contract *AutomationCompatibleUtilsCaller
	CallOpts bind.CallOpts
}

type AutomationCompatibleUtilsTransactorSession struct {
	Contract     *AutomationCompatibleUtilsTransactor
	TransactOpts bind.TransactOpts
}

type AutomationCompatibleUtilsRaw struct {
	Contract *AutomationCompatibleUtils
}

type AutomationCompatibleUtilsCallerRaw struct {
	Contract *AutomationCompatibleUtilsCaller
}

type AutomationCompatibleUtilsTransactorRaw struct {
	Contract *AutomationCompatibleUtilsTransactor
}

func NewAutomationCompatibleUtils(address common.Address, backend bind.ContractBackend) (*AutomationCompatibleUtils, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationCompatibleUtilsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationCompatibleUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationCompatibleUtils{address: address, abi: abi, AutomationCompatibleUtilsCaller: AutomationCompatibleUtilsCaller{contract: contract}, AutomationCompatibleUtilsTransactor: AutomationCompatibleUtilsTransactor{contract: contract}, AutomationCompatibleUtilsFilterer: AutomationCompatibleUtilsFilterer{contract: contract}}, nil
}

func NewAutomationCompatibleUtilsCaller(address common.Address, caller bind.ContractCaller) (*AutomationCompatibleUtilsCaller, error) {
	contract, err := bindAutomationCompatibleUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationCompatibleUtilsCaller{contract: contract}, nil
}

func NewAutomationCompatibleUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationCompatibleUtilsTransactor, error) {
	contract, err := bindAutomationCompatibleUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationCompatibleUtilsTransactor{contract: contract}, nil
}

func NewAutomationCompatibleUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationCompatibleUtilsFilterer, error) {
	contract, err := bindAutomationCompatibleUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationCompatibleUtilsFilterer{contract: contract}, nil
}

func bindAutomationCompatibleUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationCompatibleUtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationCompatibleUtils.Contract.AutomationCompatibleUtilsCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.AutomationCompatibleUtilsTransactor.contract.Transfer(opts)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.AutomationCompatibleUtilsTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationCompatibleUtils.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.contract.Transfer(opts)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactor) ConditionalTrigger(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonConditionalTrigger) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.contract.Transact(opts, "_conditionalTrigger", arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsSession) ConditionalTrigger(arg0 IAutomationV21PlusCommonConditionalTrigger) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.ConditionalTrigger(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorSession) ConditionalTrigger(arg0 IAutomationV21PlusCommonConditionalTrigger) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.ConditionalTrigger(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactor) Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.contract.Transact(opts, "_log", arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsSession) Log(arg0 Log) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.Log(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorSession) Log(arg0 Log) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.Log(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactor) LogTrigger(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonLogTrigger) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.contract.Transact(opts, "_logTrigger", arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsSession) LogTrigger(arg0 IAutomationV21PlusCommonLogTrigger) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.LogTrigger(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorSession) LogTrigger(arg0 IAutomationV21PlusCommonLogTrigger) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.LogTrigger(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactor) LogTriggerConfig(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonLogTriggerConfig) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.contract.Transact(opts, "_logTriggerConfig", arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsSession) LogTriggerConfig(arg0 IAutomationV21PlusCommonLogTriggerConfig) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.LogTriggerConfig(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorSession) LogTriggerConfig(arg0 IAutomationV21PlusCommonLogTriggerConfig) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.LogTriggerConfig(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactor) Report(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonReport) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.contract.Transact(opts, "_report", arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsSession) Report(arg0 IAutomationV21PlusCommonReport) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.Report(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtilsTransactorSession) Report(arg0 IAutomationV21PlusCommonReport) (*types.Transaction, error) {
	return _AutomationCompatibleUtils.Contract.Report(&_AutomationCompatibleUtils.TransactOpts, arg0)
}

func (_AutomationCompatibleUtils *AutomationCompatibleUtils) Address() common.Address {
	return _AutomationCompatibleUtils.address
}

type AutomationCompatibleUtilsInterface interface {
	ConditionalTrigger(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonConditionalTrigger) (*types.Transaction, error)

	Log(opts *bind.TransactOpts, arg0 Log) (*types.Transaction, error)

	LogTrigger(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonLogTrigger) (*types.Transaction, error)

	LogTriggerConfig(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonLogTriggerConfig) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, arg0 IAutomationV21PlusCommonReport) (*types.Transaction, error)

	Address() common.Address
}
