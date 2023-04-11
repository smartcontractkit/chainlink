// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_upkeep_wrapper

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

var MercuryUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"feedIDStrList\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"instance\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"MercuryLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"MercuryEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"mercuryCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"input\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000fc338038062000fc383398101604081905262000034916200018b565b6000828155600182905560038190554360029081556004829055600591909155604080516080810182526007818301818152661155120b5554d160ca1b606084015282528251808401909352825266109510cb5554d160ca1b602083810191909152810191909152620000ab9160069190620000b4565b50505062000321565b828054828255906000526020600020908101928215620000ff579160200282015b82811115620000ff5782518290620000ee908262000255565b5091602001919060010190620000d5565b506200010d92915062000111565b5090565b808211156200010d57600062000128828262000132565b5060010162000111565b5080546200014090620001c6565b6000825580601f1062000151575050565b601f01602090049060005260206000209081019062000171919062000174565b50565b5b808211156200010d576000815560010162000175565b600080604083850312156200019f57600080fd5b505080516020909101519092909150565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001db57607f821691505b602082108103620001fc57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200025057600081815260208120601f850160051c810160208610156200022b5750805b601f850160051c820191505b818110156200024c5782815560010162000237565b5050505b505050565b81516001600160401b03811115620002715762000271620001b0565b6200028981620002828454620001c6565b8462000202565b602080601f831160018114620002c15760008415620002a85750858301515b600019600386901b1c1916600185901b1785556200024c565b600085815260208120601f198616915b82811015620002f257888601518255948401946001909101908401620002d1565b5085821015620003115787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610c9280620003316000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636250a13a1161008c57806386e330af1161006657806386e330af146101a4578063917d895f146101b7578063947a36fb146101c0578063d832d92f146101c957600080fd5b80636250a13a1461017f5780636e04ff0d14610188578063806b984f1461019b57600080fd5b80634a5479f3116100bd5780634a5479f3146101355780634ad8c9a61461015557806361bc221a1461017657600080fd5b80631e34c585146100e45780632cb15864146101065780634585e33b14610122575b600080fd5b6101046100f23660046104bc565b60009182556001556004819055600555565b005b61010f60045481565b6040519081526020015b60405180910390f35b6101046101303660046104de565b6101e1565b610148610143366004610550565b61024b565b60405161011991906105d4565b610168610163366004610726565b6102f7565b6040516101199291906107ff565b61010f60055481565b61010f60005481565b6101686101963660046104de565b610302565b61010f60025481565b6101046101b2366004610822565b61039e565b61010f60035481565b61010f60015481565b6101d16103b5565b6040519015158152602001610119565b6004546000036101f057436004555b43600255600554610202906001610917565b60055560405132907f6a4ace417839897225ff8d5ff745e4dbe37f7ffe26749f7632cb1c9d719cab89906102399085908590610978565b60405180910390a25050600254600355565b6006818154811061025b57600080fd5b9060005260206000200160009150905080546102769061098c565b80601f01602080910402602001604051908101604052809291908181526020018280546102a29061098c565b80156102ef5780601f106102c4576101008083540402835291602001916102ef565b820191906000526020600020905b8154815290600101906020018083116102d257829003601f168201915b505050505081565b6001815b9250929050565b6000606061030e6103b5565b61035a576000848481818080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509597509195506102fb945050505050565b60064385856040517fa33ccff600000000000000000000000000000000000000000000000000000000815260040161039594939291906109df565b60405180910390fd5b80516103b19060069060208401906103f7565b5050565b60006004546000036103c75750600190565b6000546004546103d79043610b05565b1080156103f257506001546002546103ef9043610b05565b10155b905090565b82805482825590600052602060002090810192821561043d579160200282015b8281111561043d578251829061042d9082610b6b565b5091602001919060010190610417565b5061044992915061044d565b5090565b80821115610449576000610461828261046a565b5060010161044d565b5080546104769061098c565b6000825580601f10610486575050565b601f0160209004906000526020600020908101906104a491906104a7565b50565b5b8082111561044957600081556001016104a8565b600080604083850312156104cf57600080fd5b50508035926020909101359150565b600080602083850312156104f157600080fd5b823567ffffffffffffffff8082111561050957600080fd5b818501915085601f83011261051d57600080fd5b81358181111561052c57600080fd5b86602082850101111561053e57600080fd5b60209290920196919550909350505050565b60006020828403121561056257600080fd5b5035919050565b6000815180845260005b8181101561058f57602081850181015186830182015201610573565b818111156105a1576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006105e76020830184610569565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610664576106646105ee565b604052919050565b600067ffffffffffffffff821115610686576106866105ee565b5060051b60200190565b600067ffffffffffffffff8311156106aa576106aa6105ee565b6106db60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8601160161061d565b90508281528383830111156106ef57600080fd5b828260208301376000602084830101529392505050565b600082601f83011261071757600080fd5b6105e783833560208501610690565b6000806040838503121561073957600080fd5b823567ffffffffffffffff8082111561075157600080fd5b818501915085601f83011261076557600080fd5b8135602061077a6107758361066c565b61061d565b82815260059290921b8401810191818101908984111561079957600080fd5b8286015b848110156107d1578035868111156107b55760008081fd5b6107c38c86838b0101610706565b84525091830191830161079d565b50965050860135925050808211156107e857600080fd5b506107f585828601610706565b9150509250929050565b821515815260406020820152600061081a6040830184610569565b949350505050565b6000602080838503121561083557600080fd5b823567ffffffffffffffff8082111561084d57600080fd5b818501915085601f83011261086157600080fd5b813561086f6107758261066c565b81815260059190911b8301840190848101908883111561088e57600080fd5b8585015b838110156108db578035858111156108aa5760008081fd5b8601603f81018b136108bc5760008081fd5b6108cd8b8983013560408401610690565b845250918601918601610892565b5098975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561092a5761092a6108e8565b500190565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600061081a60208301848661092f565b600181811c908216806109a057607f821691505b6020821081036109d9577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000606082016060835280875480835260808501915060059250608081841b86010160008a81526020808220825b85811015610adc577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808a8603018752838254610a488161098c565b80885260018281168015610a635760018114610a9a57610ac5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008416888b0152878315158e1b8b01019450610ac5565b868952878920895b84811015610abd5781548c82018b0152908301908901610aa2565b8b0189019550505b509986019992975050509190910190600101610a0d565b505087018a9052508581036040870152610af781888a61092f565b9a9950505050505050505050565b600082821015610b1757610b176108e8565b500390565b601f821115610b6657600081815260208120601f850160051c81016020861015610b435750805b601f850160051c820191505b81811015610b6257828155600101610b4f565b5050505b505050565b815167ffffffffffffffff811115610b8557610b856105ee565b610b9981610b93845461098c565b84610b1c565b602080601f831160018114610bec5760008415610bb65750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610b62565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610c3957888601518255948401946001909101908401610c1a565b5085821015610c7557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c634300080f000a",
}

var MercuryUpkeepABI = MercuryUpkeepMetaData.ABI

var MercuryUpkeepBin = MercuryUpkeepMetaData.Bin

func DeployMercuryUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *MercuryUpkeep, error) {
	parsed, err := MercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryUpkeepBin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryUpkeep{MercuryUpkeepCaller: MercuryUpkeepCaller{contract: contract}, MercuryUpkeepTransactor: MercuryUpkeepTransactor{contract: contract}, MercuryUpkeepFilterer: MercuryUpkeepFilterer{contract: contract}}, nil
}

type MercuryUpkeep struct {
	address common.Address
	abi     abi.ABI
	MercuryUpkeepCaller
	MercuryUpkeepTransactor
	MercuryUpkeepFilterer
}

type MercuryUpkeepCaller struct {
	contract *bind.BoundContract
}

type MercuryUpkeepTransactor struct {
	contract *bind.BoundContract
}

type MercuryUpkeepFilterer struct {
	contract *bind.BoundContract
}

type MercuryUpkeepSession struct {
	Contract     *MercuryUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryUpkeepCallerSession struct {
	Contract *MercuryUpkeepCaller
	CallOpts bind.CallOpts
}

type MercuryUpkeepTransactorSession struct {
	Contract     *MercuryUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type MercuryUpkeepRaw struct {
	Contract *MercuryUpkeep
}

type MercuryUpkeepCallerRaw struct {
	Contract *MercuryUpkeepCaller
}

type MercuryUpkeepTransactorRaw struct {
	Contract *MercuryUpkeepTransactor
}

func NewMercuryUpkeep(address common.Address, backend bind.ContractBackend) (*MercuryUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeep{address: address, abi: abi, MercuryUpkeepCaller: MercuryUpkeepCaller{contract: contract}, MercuryUpkeepTransactor: MercuryUpkeepTransactor{contract: contract}, MercuryUpkeepFilterer: MercuryUpkeepFilterer{contract: contract}}, nil
}

func NewMercuryUpkeepCaller(address common.Address, caller bind.ContractCaller) (*MercuryUpkeepCaller, error) {
	contract, err := bindMercuryUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepCaller{contract: contract}, nil
}

func NewMercuryUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryUpkeepTransactor, error) {
	contract, err := bindMercuryUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepTransactor{contract: contract}, nil
}

func NewMercuryUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryUpkeepFilterer, error) {
	contract, err := bindMercuryUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepFilterer{contract: contract}, nil
}

func bindMercuryUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryUpkeep *MercuryUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryUpkeep.Contract.MercuryUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.MercuryUpkeepTransactor.contract.Transfer(opts)
}

func (_MercuryUpkeep *MercuryUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.MercuryUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.contract.Transfer(opts)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.CheckUpkeep(&_MercuryUpkeep.CallOpts, data)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.CheckUpkeep(&_MercuryUpkeep.CallOpts, data)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Counter() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Counter(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Counter() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Counter(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Eligible() (bool, error) {
	return _MercuryUpkeep.Contract.Eligible(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Eligible() (bool, error) {
	return _MercuryUpkeep.Contract.Eligible(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "feeds", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Feeds(arg0 *big.Int) (string, error) {
	return _MercuryUpkeep.Contract.Feeds(&_MercuryUpkeep.CallOpts, arg0)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Feeds(arg0 *big.Int) (string, error) {
	return _MercuryUpkeep.Contract.Feeds(&_MercuryUpkeep.CallOpts, arg0)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) InitialBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.InitialBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) InitialBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.InitialBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Interval() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Interval(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Interval() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Interval(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) LastBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.LastBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) LastBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.LastBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "mercuryCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) MercuryCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.MercuryCallback(&_MercuryUpkeep.CallOpts, values, extraData)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) MercuryCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.MercuryCallback(&_MercuryUpkeep.CallOpts, values, extraData)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) PreviousPerformBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.PreviousPerformBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.PreviousPerformBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) TestRange() (*big.Int, error) {
	return _MercuryUpkeep.Contract.TestRange(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) TestRange() (*big.Int, error) {
	return _MercuryUpkeep.Contract.TestRange(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_MercuryUpkeep *MercuryUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.PerformUpkeep(&_MercuryUpkeep.TransactOpts, performData)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.PerformUpkeep(&_MercuryUpkeep.TransactOpts, performData)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setConfig", _testRange, _interval)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetConfig(&_MercuryUpkeep.TransactOpts, _testRange, _interval)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetConfig(&_MercuryUpkeep.TransactOpts, _testRange, _interval)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, input []string) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setFeeds", input)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetFeeds(input []string) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetFeeds(&_MercuryUpkeep.TransactOpts, input)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetFeeds(input []string) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetFeeds(&_MercuryUpkeep.TransactOpts, input)
}

type MercuryUpkeepMercuryEventIterator struct {
	Event *MercuryUpkeepMercuryEvent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryUpkeepMercuryEventIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryUpkeepMercuryEvent)
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
		it.Event = new(MercuryUpkeepMercuryEvent)
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

func (it *MercuryUpkeepMercuryEventIterator) Error() error {
	return it.fail
}

func (it *MercuryUpkeepMercuryEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryUpkeepMercuryEvent struct {
	From common.Address
	Data []byte
	Raw  types.Log
}

func (_MercuryUpkeep *MercuryUpkeepFilterer) FilterMercuryEvent(opts *bind.FilterOpts, from []common.Address) (*MercuryUpkeepMercuryEventIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _MercuryUpkeep.contract.FilterLogs(opts, "MercuryEvent", fromRule)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepMercuryEventIterator{contract: _MercuryUpkeep.contract, event: "MercuryEvent", logs: logs, sub: sub}, nil
}

func (_MercuryUpkeep *MercuryUpkeepFilterer) WatchMercuryEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryEvent, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _MercuryUpkeep.contract.WatchLogs(opts, "MercuryEvent", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryUpkeepMercuryEvent)
				if err := _MercuryUpkeep.contract.UnpackLog(event, "MercuryEvent", log); err != nil {
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

func (_MercuryUpkeep *MercuryUpkeepFilterer) ParseMercuryEvent(log types.Log) (*MercuryUpkeepMercuryEvent, error) {
	event := new(MercuryUpkeepMercuryEvent)
	if err := _MercuryUpkeep.contract.UnpackLog(event, "MercuryEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MercuryUpkeep *MercuryUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MercuryUpkeep.abi.Events["MercuryEvent"].ID:
		return _MercuryUpkeep.ParseMercuryEvent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MercuryUpkeepMercuryEvent) Topic() common.Hash {
	return common.HexToHash("0x6a4ace417839897225ff8d5ff745e4dbe37f7ffe26749f7632cb1c9d719cab89")
}

func (_MercuryUpkeep *MercuryUpkeep) Address() common.Address {
	return _MercuryUpkeep.address
}

type MercuryUpkeepInterface interface {
	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetFeeds(opts *bind.TransactOpts, input []string) (*types.Transaction, error)

	FilterMercuryEvent(opts *bind.FilterOpts, from []common.Address) (*MercuryUpkeepMercuryEventIterator, error)

	WatchMercuryEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryEvent, from []common.Address) (event.Subscription, error)

	ParseMercuryEvent(log types.Log) (*MercuryUpkeepMercuryEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
