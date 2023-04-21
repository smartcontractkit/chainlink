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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedLabel\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feedList\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"queryLabel\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"query\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"MercuryLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"MercuryEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"mercuryCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"input\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620010cd380380620010cd833981016040819052620000349162000215565b6000828155600182905560038190554360025560048190556005556040805180820190915260098152683332b2b224a229ba3960b91b60208201526007906200007e9082620002df565b506040805160808101825260188183018181527f4554482d5553442d415242495452554d2d544553544e4554000000000000000060608401528252825180840190935282527f4254432d5553442d415242495452554d2d544553544e45540000000000000000602083810191909152810191909152620001039060069060026200013e565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152600890620001359082620002df565b505050620003ab565b82805482825590600052602060002090810192821562000189579160200282015b82811115620001895782518290620001789082620002df565b50916020019190600101906200015f565b50620001979291506200019b565b5090565b8082111562000197576000620001b28282620001bc565b506001016200019b565b508054620001ca9062000250565b6000825580601f10620001db575050565b601f016020900490600052602060002090810190620001fb9190620001fe565b50565b5b80821115620001975760008155600101620001ff565b600080604083850312156200022957600080fd5b505080516020909101519092909150565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200026557607f821691505b6020821081036200028657634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002da57600081815260208120601f850160051c81016020861015620002b55750805b601f850160051c820191505b81811015620002d657828155600101620002c1565b5050505b505050565b81516001600160401b03811115620002fb57620002fb6200023a565b62000313816200030c845462000250565b846200028c565b602080601f8311600181146200034b5760008415620003325750858301515b600019600386901b1c1916600185901b178555620002d6565b600085815260208120601f198616915b828110156200037c578886015182559484019460019091019084016200035b565b50858210156200039b5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610d1280620003bb6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80636250a13a1161009757806386e330af1161006657806386e330af146101ca578063917d895f146101dd578063947a36fb146101e6578063d832d92f146101ef57600080fd5b80636250a13a1461019d5780636e04ff0d146101a6578063806b984f146101b957806380f4df1b146101c257600080fd5b80634a5479f3116100d35780634a5479f31461014b5780634ad8c9a61461016b5780634d6954451461018c57806361bc221a1461019457600080fd5b80631e34c585146100fa5780632cb158641461011c5780634585e33b14610138575b600080fd5b61011a610108366004610502565b60009182556001556004819055600555565b005b61012560045481565b6040519081526020015b60405180910390f35b61011a610146366004610524565b610207565b61015e610159366004610596565b610271565b60405161012f919061061a565b61017e61017936600461076c565b61031d565b60405161012f929190610845565b61015e610328565b61012560055481565b61012560005481565b61017e6101b4366004610524565b610335565b61012560025481565b61015e6103d7565b61011a6101d8366004610868565b6103e4565b61012560035481565b61012560015481565b6101f76103fb565b604051901515815260200161012f565b60045460000361021657436004555b4360025560055461022890600161095d565b60055560405132907f6a4ace417839897225ff8d5ff745e4dbe37f7ffe26749f7632cb1c9d719cab899061025f90859085906109be565b60405180910390a25050600254600355565b6006818154811061028157600080fd5b90600052602060002001600091509050805461029c906109d2565b80601f01602080910402602001604051908101604052809291908181526020018280546102c8906109d2565b80156103155780601f106102ea57610100808354040283529160200191610315565b820191906000526020600020905b8154815290600101906020018083116102f857829003601f168201915b505050505081565b6001815b9250929050565b6007805461029c906109d2565b600060606103416103fb565b61038d576000848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959750919550610321945050505050565b6007600660084387876040517f62e8a50d0000000000000000000000000000000000000000000000000000000081526004016103ce96959493929190610ac0565b60405180910390fd5b6008805461029c906109d2565b80516103f790600690602084019061043d565b5050565b600060045460000361040d5750600190565b60005460045461041d9043610b85565b10801561043857506001546002546104359043610b85565b10155b905090565b828054828255906000526020600020908101928215610483579160200282015b8281111561048357825182906104739082610beb565b509160200191906001019061045d565b5061048f929150610493565b5090565b8082111561048f5760006104a782826104b0565b50600101610493565b5080546104bc906109d2565b6000825580601f106104cc575050565b601f0160209004906000526020600020908101906104ea91906104ed565b50565b5b8082111561048f57600081556001016104ee565b6000806040838503121561051557600080fd5b50508035926020909101359150565b6000806020838503121561053757600080fd5b823567ffffffffffffffff8082111561054f57600080fd5b818501915085601f83011261056357600080fd5b81358181111561057257600080fd5b86602082850101111561058457600080fd5b60209290920196919550909350505050565b6000602082840312156105a857600080fd5b5035919050565b6000815180845260005b818110156105d5576020818501810151868301820152016105b9565b818111156105e7576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061062d60208301846105af565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156106aa576106aa610634565b604052919050565b600067ffffffffffffffff8211156106cc576106cc610634565b5060051b60200190565b600067ffffffffffffffff8311156106f0576106f0610634565b61072160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601610663565b905082815283838301111561073557600080fd5b828260208301376000602084830101529392505050565b600082601f83011261075d57600080fd5b61062d838335602085016106d6565b6000806040838503121561077f57600080fd5b823567ffffffffffffffff8082111561079757600080fd5b818501915085601f8301126107ab57600080fd5b813560206107c06107bb836106b2565b610663565b82815260059290921b840181019181810190898411156107df57600080fd5b8286015b84811015610817578035868111156107fb5760008081fd5b6108098c86838b010161074c565b8452509183019183016107e3565b509650508601359250508082111561082e57600080fd5b5061083b8582860161074c565b9150509250929050565b821515815260406020820152600061086060408301846105af565b949350505050565b6000602080838503121561087b57600080fd5b823567ffffffffffffffff8082111561089357600080fd5b818501915085601f8301126108a757600080fd5b81356108b56107bb826106b2565b81815260059190911b830184019084810190888311156108d457600080fd5b8585015b83811015610921578035858111156108f05760008081fd5b8601603f81018b136109025760008081fd5b6109138b89830135604084016106d6565b8452509186019186016108d8565b5098975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600082198211156109705761097061092e565b500190565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000610860602083018486610975565b600181811c908216806109e657607f821691505b602082108103610a1f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60008154610a32816109d2565b808552602060018381168015610a4f5760018114610a8757610ab5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550610ab5565b866000528260002060005b85811015610aad5781548a8201860152908301908401610a92565b890184019650505b505050505092915050565b60a081526000610ad360a0830189610a25565b6020838203818501528189548084528284019150828160051b8501018b6000528360002060005b83811015610b45577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552610b338383610a25565b94860194925060019182019101610afa565b50508681036040880152610b59818c610a25565b9450505050508560608401528281036080840152610b78818587610975565b9998505050505050505050565b600082821015610b9757610b9761092e565b500390565b601f821115610be657600081815260208120601f850160051c81016020861015610bc35750805b601f850160051c820191505b81811015610be257828155600101610bcf565b5050505b505050565b815167ffffffffffffffff811115610c0557610c05610634565b610c1981610c1384546109d2565b84610b9c565b602080601f831160018114610c6c5760008415610c365750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610be2565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610cb957888601518255948401946001909101908401610c9a565b5085821015610cf557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c634300080f000a",
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

func (_MercuryUpkeep *MercuryUpkeepCaller) FeedLabel(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "feedLabel")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) FeedLabel() (string, error) {
	return _MercuryUpkeep.Contract.FeedLabel(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) FeedLabel() (string, error) {
	return _MercuryUpkeep.Contract.FeedLabel(&_MercuryUpkeep.CallOpts)
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

func (_MercuryUpkeep *MercuryUpkeepCaller) QueryLabel(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "queryLabel")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) QueryLabel() (string, error) {
	return _MercuryUpkeep.Contract.QueryLabel(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) QueryLabel() (string, error) {
	return _MercuryUpkeep.Contract.QueryLabel(&_MercuryUpkeep.CallOpts)
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

	FeedLabel(opts *bind.CallOpts) (string, error)

	Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	QueryLabel(opts *bind.CallOpts) (string, error)

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
