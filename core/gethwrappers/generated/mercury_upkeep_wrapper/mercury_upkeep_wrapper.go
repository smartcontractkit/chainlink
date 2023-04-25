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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedLabel\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feedList\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"queryLabel\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"query\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"MercuryLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"MercuryEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"mercuryCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"input\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200122938038062001229833981016040819052620000349162000215565b6000828155600182905560038190554360025560048190556005556040805180820190915260098152683332b2b224a229ba3960b91b60208201526007906200007e9082620002df565b506040805160808101825260188183018181527f4554482d5553442d415242495452554d2d544553544e4554000000000000000060608401528252825180840190935282527f4254432d5553442d415242495452554d2d544553544e45540000000000000000602083810191909152810191909152620001039060069060026200013e565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152600890620001359082620002df565b505050620003ab565b82805482825590600052602060002090810192821562000189579160200282015b82811115620001895782518290620001789082620002df565b50916020019190600101906200015f565b50620001979291506200019b565b5090565b8082111562000197576000620001b28282620001bc565b506001016200019b565b508054620001ca9062000250565b6000825580601f10620001db575050565b601f016020900490600052602060002090810190620001fb9190620001fe565b50565b5b80821115620001975760008155600101620001ff565b600080604083850312156200022957600080fd5b505080516020909101519092909150565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200026557607f821691505b6020821081036200028657634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002da57600081815260208120601f850160051c81016020861015620002b55750805b601f850160051c820191505b81811015620002d657828155600101620002c1565b5050505b505050565b81516001600160401b03811115620002fb57620002fb6200023a565b62000313816200030c845462000250565b846200028c565b602080601f8311600181146200034b5760008415620003325750858301515b600019600386901b1c1916600185901b178555620002d6565b600085815260208120601f198616915b828110156200037c578886015182559484019460019091019084016200035b565b50858210156200039b5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610e6e80620003bb6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80636250a13a1161009757806386e330af1161006657806386e330af146101ca578063917d895f146101dd578063947a36fb146101e6578063d832d92f146101ef57600080fd5b80636250a13a1461019d5780636e04ff0d146101a6578063806b984f146101b957806380f4df1b146101c257600080fd5b80634a5479f3116100d35780634a5479f31461014b5780634ad8c9a61461016b5780634d6954451461018c57806361bc221a1461019457600080fd5b80631e34c585146100fa5780632cb158641461011c5780634585e33b14610138575b600080fd5b61011a6101083660046105b9565b60009182556001556004819055600555565b005b61012560045481565b6040519081526020015b60405180910390f35b61011a6101463660046105db565b610207565b61015e61015936600461064d565b610273565b60405161012f91906106e0565b61017e610179366004610832565b61031f565b60405161012f92919061090b565b61015e6103df565b61012560055481565b61012560005481565b61017e6101b43660046105db565b6103ec565b61012560025481565b61015e61048e565b61011a6101d836600461092e565b61049b565b61012560035481565b61012560015481565b6101f76104b2565b604051901515815260200161012f565b60045460000361021657436004555b43600255600554610228906001610a23565b600555604051339032907f297347099d64607af3294cc51e20c531fee83039a87428f2eea9cb85f397119a906102619086908690610a84565b60405180910390a35050600254600355565b6006818154811061028357600080fd5b90600052602060002001600091509050805461029e90610a98565b80601f01602080910402602001604051908101604052809291908181526020018280546102ca90610a98565b80156103175780601f106102ec57610100808354040283529160200191610317565b820191906000526020600020905b8154815290600101906020018083116102fa57829003601f168201915b505050505081565b6040805160008082526020820190925260609060005b855181101561038e578186828151811061035157610351610aeb565b602002602001015160405160200161036a929190610b1a565b6040516020818303038152906040529150808061038690610b49565b915050610335565b5080846040516020016103a2929190610b1a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b6007805461029e90610a98565b600060606103f86104b2565b610444576000848481818080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509597509195506103d8945050505050565b6007600660084387876040517f62e8a50d00000000000000000000000000000000000000000000000000000000815260040161048596959493929190610c1c565b60405180910390fd5b6008805461029e90610a98565b80516104ae9060069060208401906104f4565b5050565b60006004546000036104c45750600190565b6000546004546104d49043610ce1565b1080156104ef57506001546002546104ec9043610ce1565b10155b905090565b82805482825590600052602060002090810192821561053a579160200282015b8281111561053a578251829061052a9082610d47565b5091602001919060010190610514565b5061054692915061054a565b5090565b8082111561054657600061055e8282610567565b5060010161054a565b50805461057390610a98565b6000825580601f10610583575050565b601f0160209004906000526020600020908101906105a191906105a4565b50565b5b8082111561054657600081556001016105a5565b600080604083850312156105cc57600080fd5b50508035926020909101359150565b600080602083850312156105ee57600080fd5b823567ffffffffffffffff8082111561060657600080fd5b818501915085601f83011261061a57600080fd5b81358181111561062957600080fd5b86602082850101111561063b57600080fd5b60209290920196919550909350505050565b60006020828403121561065f57600080fd5b5035919050565b60005b83811015610681578181015183820152602001610669565b83811115610690576000848401525b50505050565b600081518084526106ae816020860160208601610666565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006106f36020830184610696565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610770576107706106fa565b604052919050565b600067ffffffffffffffff821115610792576107926106fa565b5060051b60200190565b600067ffffffffffffffff8311156107b6576107b66106fa565b6107e760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601610729565b90508281528383830111156107fb57600080fd5b828260208301376000602084830101529392505050565b600082601f83011261082357600080fd5b6106f38383356020850161079c565b6000806040838503121561084557600080fd5b823567ffffffffffffffff8082111561085d57600080fd5b818501915085601f83011261087157600080fd5b8135602061088661088183610778565b610729565b82815260059290921b840181019181810190898411156108a557600080fd5b8286015b848110156108dd578035868111156108c15760008081fd5b6108cf8c86838b0101610812565b8452509183019183016108a9565b50965050860135925050808211156108f457600080fd5b5061090185828601610812565b9150509250929050565b82151581526040602082015260006109266040830184610696565b949350505050565b6000602080838503121561094157600080fd5b823567ffffffffffffffff8082111561095957600080fd5b818501915085601f83011261096d57600080fd5b813561097b61088182610778565b81815260059190911b8301840190848101908883111561099a57600080fd5b8585015b838110156109e7578035858111156109b65760008081fd5b8601603f81018b136109c85760008081fd5b6109d98b898301356040840161079c565b84525091860191860161099e565b5098975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115610a3657610a366109f4565b500190565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000610926602083018486610a3b565b600181811c90821680610aac57607f821691505b602082108103610ae5577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008351610b2c818460208801610666565b835190830190610b40818360208801610666565b01949350505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610b7a57610b7a6109f4565b5060010190565b60008154610b8e81610a98565b808552602060018381168015610bab5760018114610be357610c11565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550610c11565b866000528260002060005b85811015610c095781548a8201860152908301908401610bee565b890184019650505b505050505092915050565b60a081526000610c2f60a0830189610b81565b6020838203818501528189548084528284019150828160051b8501018b6000528360002060005b83811015610ca1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552610c8f8383610b81565b94860194925060019182019101610c56565b50508681036040880152610cb5818c610b81565b9450505050508560608401528281036080840152610cd4818587610a3b565b9998505050505050505050565b600082821015610cf357610cf36109f4565b500390565b601f821115610d4257600081815260208120601f850160051c81016020861015610d1f5750805b601f850160051c820191505b81811015610d3e57828155600101610d2b565b5050505b505050565b815167ffffffffffffffff811115610d6157610d616106fa565b610d7581610d6f8454610a98565b84610cf8565b602080601f831160018114610dc85760008415610d925750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610d3e565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610e1557888601518255948401946001909101908401610df6565b5085821015610e5157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c634300080f000a",
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
	Origin common.Address
	Sender common.Address
	Data   []byte
	Raw    types.Log
}

func (_MercuryUpkeep *MercuryUpkeepFilterer) FilterMercuryEvent(opts *bind.FilterOpts, origin []common.Address, sender []common.Address) (*MercuryUpkeepMercuryEventIterator, error) {

	var originRule []interface{}
	for _, originItem := range origin {
		originRule = append(originRule, originItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MercuryUpkeep.contract.FilterLogs(opts, "MercuryEvent", originRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepMercuryEventIterator{contract: _MercuryUpkeep.contract, event: "MercuryEvent", logs: logs, sub: sub}, nil
}

func (_MercuryUpkeep *MercuryUpkeepFilterer) WatchMercuryEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryEvent, origin []common.Address, sender []common.Address) (event.Subscription, error) {

	var originRule []interface{}
	for _, originItem := range origin {
		originRule = append(originRule, originItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MercuryUpkeep.contract.WatchLogs(opts, "MercuryEvent", originRule, senderRule)
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
	return common.HexToHash("0x297347099d64607af3294cc51e20c531fee83039a87428f2eea9cb85f397119a")
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

	FilterMercuryEvent(opts *bind.FilterOpts, origin []common.Address, sender []common.Address) (*MercuryUpkeepMercuryEventIterator, error)

	WatchMercuryEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryEvent, origin []common.Address, sender []common.Address) (event.Subscription, error)

	ParseMercuryEvent(log types.Log) (*MercuryUpkeepMercuryEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
