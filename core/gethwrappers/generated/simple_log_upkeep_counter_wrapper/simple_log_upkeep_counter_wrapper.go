// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple_log_upkeep_counter_wrapper

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

type CheckData struct {
	CheckBurnAmount   *big.Int
	PerformBurnAmount *big.Int
	EventSig          [32]byte
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

var SimpleLogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeToPerform\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isRecovered\",\"type\":\"bool\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"checkBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"eventSig\",\"type\":\"bytes32\"}],\"internalType\":\"structCheckData\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_checkDataConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isRecovered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeToPerform\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060006002819055436001556003819055600455610d12806100336000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80637145f11b11610076578063917d895f1161005b578063917d895f1461016b578063c6066f0d14610174578063eb950ce71461017d57600080fd5b80637145f11b1461012f578063806b984f1461016257600080fd5b80634585e33b116100a75780634585e33b1461010057806361bc221a14610115578063697794731461011e57600080fd5b80632cb15864146100c357806340691db4146100df575b600080fd5b6100cc60035481565b6040519081526020015b60405180910390f35b6100f26100ed3660046106c6565b61018a565b6040516100d692919061092d565b61011361010e366004610628565b6102c2565b005b6100cc60045481565b61011361012c36600461066a565b50565b61015261013d36600461060f565b60006020819052908152604090205460ff1681565b60405190151581526020016100d6565b6100cc60015481565b6100cc60025481565b6100cc60055481565b6006546101529060ff1681565b6000606081808061019d8688018861083b565b92509250925060005a905060006101b5600143610c61565b40905060008515610224575b855a6101cd9085610c61565b1015610224578080156101ee575060008281526020819052604090205460ff165b604080516020810185905230918101919091529091506060016040516020818303038152906040528051906020012091506101c1565b8361023260c08d018d610a9d565b600281811061024357610243610ca7565b9050602002013514156102875760018b438c8c60405160200161026994939291906109aa565b604051602081830303815290604052975097505050505050506102ba565b60008b438c8c6040516020016102a094939291906109aa565b604051602081830303815290604052975097505050505050505b935093915050565b6003546102ce57436003555b4360019081556004546102e091610c49565b600455600154600255600080806102f984860186610738565b92509250925082602001514261030f9190610c61565b600555600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556060830151821461037157600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555b60008060008380602001905181019061038a9190610867565b92509250925060005a905060006103a2600143610c61565b40905060008415610411575b845a6103ba9085610c61565b1015610411578080156103db575060008281526020819052604090205460ff165b604080516020810185905230918101919091529091506060016040516020818303038152906040528051906020012091506103ae565b600354600154600254600454600554600654604080519687526020870195909552938501929092526060840152608083015260ff16151560a082015232907f29eff4cb37911c3ea85db4630638cc5474fdd0631ec42215aef1d7ec96c8e63d9060c00160405180910390a25050505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146104ad57600080fd5b919050565b600082601f8301126104c357600080fd5b8135602067ffffffffffffffff8211156104df576104df610cd6565b8160051b6104ee828201610b2f565b83815282810190868401838801850189101561050957600080fd5b600093505b8584101561052c57803583526001939093019291840191840161050e565b50979650505050505050565b60008083601f84011261054a57600080fd5b50813567ffffffffffffffff81111561056257600080fd5b60208301915083602082850101111561057a57600080fd5b9250929050565b600082601f83011261059257600080fd5b813567ffffffffffffffff8111156105ac576105ac610cd6565b6105dd60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610b2f565b8181528460208386010111156105f257600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561062157600080fd5b5035919050565b6000806020838503121561063b57600080fd5b823567ffffffffffffffff81111561065257600080fd5b61065e85828601610538565b90969095509350505050565b60006060828403121561067c57600080fd5b6040516060810181811067ffffffffffffffff8211171561069f5761069f610cd6565b80604052508235815260208301356020820152604083013560408201528091505092915050565b6000806000604084860312156106db57600080fd5b833567ffffffffffffffff808211156106f357600080fd5b90850190610100828803121561070857600080fd5b9093506020850135908082111561071e57600080fd5b5061072b86828701610538565b9497909650939450505050565b60008060006060848603121561074d57600080fd5b833567ffffffffffffffff8082111561076557600080fd5b90850190610100828803121561077a57600080fd5b610782610b05565b82358152602083013560208201526040830135604082015260608301356060820152608083013560808201526107ba60a08401610489565b60a082015260c0830135828111156107d157600080fd5b6107dd898286016104b2565b60c08301525060e0830135828111156107f557600080fd5b61080189828601610581565b60e083015250945060208601359350604086013591508082111561082457600080fd5b5061083186828701610581565b9150509250925092565b60008060006060848603121561085057600080fd5b505081359360208301359350604090920135919050565b60008060006060848603121561087c57600080fd5b8351925060208401519150604084015190509250925092565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8311156108c757600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b821515815260006020604081840152835180604085015260005b8181101561096357858101830151858201606001528201610947565b81811115610975576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b606081528435606082015260208501356080820152604085013560a0820152606085013560c0820152608085013560e082015260006109eb60a08701610489565b61010073ffffffffffffffffffffffffffffffffffffffff821681850152610a1660c0890189610b7e565b925081610120860152610a2e61016086018483610895565b92505050610a3f60e0880188610be5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa085840301610140860152610a758382846108e4565b925050508560208401528281036040840152610a928185876108e4565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610ad257600080fd5b83018035915067ffffffffffffffff821115610aed57600080fd5b6020019150600581901b360382131561057a57600080fd5b604051610100810167ffffffffffffffff81118282101715610b2957610b29610cd6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610b7657610b76610cd6565b604052919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610bb357600080fd5b830160208101925035905067ffffffffffffffff811115610bd357600080fd5b8060051b360383131561057a57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610c1a57600080fd5b830160208101925035905067ffffffffffffffff811115610c3a57600080fd5b80360383131561057a57600080fd5b60008219821115610c5c57610c5c610c78565b500190565b600082821015610c7357610c73610c78565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var SimpleLogUpkeepCounterABI = SimpleLogUpkeepCounterMetaData.ABI

var SimpleLogUpkeepCounterBin = SimpleLogUpkeepCounterMetaData.Bin

func DeploySimpleLogUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleLogUpkeepCounter, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleLogUpkeepCounterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleLogUpkeepCounter{address: address, abi: *parsed, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

type SimpleLogUpkeepCounter struct {
	address common.Address
	abi     abi.ABI
	SimpleLogUpkeepCounterCaller
	SimpleLogUpkeepCounterTransactor
	SimpleLogUpkeepCounterFilterer
}

type SimpleLogUpkeepCounterCaller struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterTransactor struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterFilterer struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterSession struct {
	Contract     *SimpleLogUpkeepCounter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterCallerSession struct {
	Contract *SimpleLogUpkeepCounterCaller
	CallOpts bind.CallOpts
}

type SimpleLogUpkeepCounterTransactorSession struct {
	Contract     *SimpleLogUpkeepCounterTransactor
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterRaw struct {
	Contract *SimpleLogUpkeepCounter
}

type SimpleLogUpkeepCounterCallerRaw struct {
	Contract *SimpleLogUpkeepCounterCaller
}

type SimpleLogUpkeepCounterTransactorRaw struct {
	Contract *SimpleLogUpkeepCounterTransactor
}

func NewSimpleLogUpkeepCounter(address common.Address, backend bind.ContractBackend) (*SimpleLogUpkeepCounter, error) {
	abi, err := abi.JSON(strings.NewReader(SimpleLogUpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSimpleLogUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounter{address: address, abi: abi, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

func NewSimpleLogUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*SimpleLogUpkeepCounterCaller, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterCaller{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleLogUpkeepCounterTransactor, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterTransactor{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleLogUpkeepCounterFilterer, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterFilterer{contract: contract}, nil
}

func bindSimpleLogUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) CheckLog(opts *bind.CallOpts, log Log, checkData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "checkLog", log, checkData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckLog(log Log, checkData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, checkData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) CheckLog(log Log, checkData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, checkData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.DummyMap(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.DummyMap(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) IsRecovered(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "isRecovered")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) IsRecovered() (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.IsRecovered(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) IsRecovered() (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.IsRecovered(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) TimeToPerform(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "timeToPerform")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) CheckDataConfig(opts *bind.TransactOpts, arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "_checkDataConfig", arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckDataConfig(arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckDataConfig(&_SimpleLogUpkeepCounter.TransactOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) CheckDataConfig(arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckDataConfig(&_SimpleLogUpkeepCounter.TransactOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

type SimpleLogUpkeepCounterPerformingUpkeepIterator struct {
	Event *SimpleLogUpkeepCounterPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleLogUpkeepCounterPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	TimeToPerform *big.Int
	IsRecovered   bool
	Raw           types.Log
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterPerformingUpkeepIterator{contract: _SimpleLogUpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleLogUpkeepCounterPerformingUpkeep)
				if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error) {
	event := new(SimpleLogUpkeepCounterPerformingUpkeep)
	if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SimpleLogUpkeepCounter.abi.Events["PerformingUpkeep"].ID:
		return _SimpleLogUpkeepCounter.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SimpleLogUpkeepCounterPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x29eff4cb37911c3ea85db4630638cc5474fdd0631ec42215aef1d7ec96c8e63d")
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) Address() common.Address {
	return _SimpleLogUpkeepCounter.address
}

type SimpleLogUpkeepCounterInterface interface {
	CheckLog(opts *bind.CallOpts, log Log, checkData []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	IsRecovered(opts *bind.CallOpts) (bool, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TimeToPerform(opts *bind.CallOpts) (*big.Int, error)

	CheckDataConfig(opts *bind.TransactOpts, arg0 CheckData) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
