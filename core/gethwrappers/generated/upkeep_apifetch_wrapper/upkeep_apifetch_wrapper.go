// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_apifetch_wrapper

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

var UpkeepAPIFetchMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"jsonFields\",\"type\":\"string[]\"},{\"internalType\":\"bytes4\",\"name\":\"callbackSelector\",\"type\":\"bytes4\"}],\"name\":\"ChainlinkAPIFetch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"values\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"statusCode\",\"type\":\"uint256\"}],\"name\":\"callback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fields\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"id\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pokemon\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"input\",\"type\":\"string\"}],\"name\":\"setURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"url\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200158b3803806200158b8339810160408190526200003491620001a8565b600082815560018290556003819055436002908155600482815560059290925560408051608081018252808201838152611a5960f21b606083015281528151808301909252928152636e616d6560e01b6020828101919091528301526200009e91600991620000d1565b506040518060600160405280602281526020016200156960229139600690620000c8908262000272565b5050506200033e565b8280548282559060005260206000209081019282156200011c579160200282015b828111156200011c57825182906200010b908262000272565b5091602001919060010190620000f2565b506200012a9291506200012e565b5090565b808211156200012a5760006200014582826200014f565b506001016200012e565b5080546200015d90620001e3565b6000825580601f106200016e575050565b601f0160209004906000526020600020908101906200018e919062000191565b50565b5b808211156200012a576000815560010162000192565b60008060408385031215620001bc57600080fd5b505080516020909101519092909150565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001f857607f821691505b6020821081036200021957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200026d57600081815260208120601f850160051c81016020861015620002485750805b601f850160051c820191505b81811015620002695782815560010162000254565b5050505b505050565b81516001600160401b038111156200028e576200028e620001cd565b620002a6816200029f8454620001e3565b846200021f565b602080601f831160018114620002de5760008415620002c55750858301515b600019600386901b1c1916600185901b17855562000269565b600085815260208120601f198616915b828110156200030f57888601518255948401946001909101908401620002ee565b50858210156200032e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61121b806200034e6000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80636e04ff0d11610097578063947a36fb11610066578063947a36fb146101e2578063af640d0f146101eb578063b772d70a146101f3578063d832d92f1461020657600080fd5b80636e04ff0d1461019c57806371106628146101bd578063806b984f146101d0578063917d895f146101d957600080fd5b80634585e33b116100d35780634585e33b1461016f5780635600f04f1461018257806361bc221a1461018a5780636250a13a1461019357600080fd5b806319137856146101055780631e34c5851461012e5780632cb1586414610150578063362ff8ac14610167575b600080fd5b610118610113366004610804565b61021e565b6040516101259190610897565b60405180910390f35b61014e61013c3660046108b1565b60009182556001556004819055600555565b005b61015960045481565b604051908152602001610125565b6101186102ca565b61014e61017d366004610915565b6102d7565b6101186103f0565b61015960055481565b61015960005481565b6101af6101aa366004610915565b6103fd565b604051610125929190610957565b61014e6101cb366004610a4c565b61053f565b61015960025481565b61015960035481565b61015960015481565b61011861054f565b6101af610201366004610a81565b61055c565b61020e610685565b6040519015158152602001610125565b6009818154811061022e57600080fd5b90600052602060002001600091509050805461024990610b23565b80601f016020809104026020016040519081016040528092919081815260200182805461027590610b23565b80156102c25780601f10610297576101008083540402835291602001916102c2565b820191906000526020600020905b8154815290600101906020018083116102a557829003601f168201915b505050505081565b6008805461024990610b23565b6004546000036102e657436004555b436002556005546102f8906001610ba5565b60055560008061030a83850185610bbd565b6040517f6572726f720000000000000000000000000000000000000000000000000000006020820152919350915060250160405160208183030381529060405280519060200120826040516020016103629190610c21565b60405160208183030381529060405280519060200120036103875760006005556103a2565b60076103938382610c8c565b5060086103a08282610c8c565b505b60025460055460405132927f0b36ab4eaaf3df8d059b63746075640cc0f87f07d011137505160b1de3c983da926103dc9287908790610da6565b60405180910390a250506002546003555050565b6006805461024990610b23565b60006060610409610685565b610455576000848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959750919550610538945050505050565b600061046e60055460016104699190610ba5565b6106c7565b90506000600682604051602001610486929190610de2565b60405160208183030381529060405290508086866040516020016104ab929190610e87565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527fda8b321400000000000000000000000000000000000000000000000000000000825261052f92916009907fb772d70a0000000000000000000000000000000000000000000000000000000090600401610ed4565b60405180910390fd5b9250929050565b600661054b8282610c8c565b5050565b6007805461024990610b23565b6000606061012b83111561059357600160405160200161057b90611036565b6040516020818303038152906040529150915061067b565b6000858560008181106105a8576105a86110af565b90506020028101906105ba91906110de565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093945089925088915060019050818110610606576106066110af565b905060200281019061061891906110de565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060405192935060019261066592508591508490602001611143565b6040516020818303038152906040529350935050505b9550959350505050565b60006004546000036106975750600190565b6000546004546106a79043611168565b1080156106c257506001546002546106bf9043611168565b10155b905090565b60608160000361070a57505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b8115610734578061071e8161117f565b915061072d9050600a836111e6565b915061070e565b60008167ffffffffffffffff81111561074f5761074f610972565b6040519080825280601f01601f191660200182016040528015610779576020820181803683370190505b5090505b84156107fc5761078e600183611168565b915061079b600a866111fa565b6107a6906030610ba5565b60f81b8183815181106107bb576107bb6110af565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506107f5600a866111e6565b945061077d565b949350505050565b60006020828403121561081657600080fd5b5035919050565b60005b83811015610838578181015183820152602001610820565b83811115610847576000848401525b50505050565b6000815180845261086581602086016020860161081d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006108aa602083018461084d565b9392505050565b600080604083850312156108c457600080fd5b50508035926020909101359150565b60008083601f8401126108e557600080fd5b50813567ffffffffffffffff8111156108fd57600080fd5b60208301915083602082850101111561053857600080fd5b6000806020838503121561092857600080fd5b823567ffffffffffffffff81111561093f57600080fd5b61094b858286016108d3565b90969095509350505050565b82151581526040602082015260006107fc604083018461084d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f8301126109b257600080fd5b813567ffffffffffffffff808211156109cd576109cd610972565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715610a1357610a13610972565b81604052838152866020858801011115610a2c57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060208284031215610a5e57600080fd5b813567ffffffffffffffff811115610a7557600080fd5b6107fc848285016109a1565b600080600080600060608688031215610a9957600080fd5b853567ffffffffffffffff80821115610ab157600080fd5b610abd89838a016108d3565b90975095506020880135915080821115610ad657600080fd5b818801915088601f830112610aea57600080fd5b813581811115610af957600080fd5b8960208260051b8501011115610b0e57600080fd5b96999598505060200195604001359392505050565b600181811c90821680610b3757607f821691505b602082108103610b70577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115610bb857610bb8610b76565b500190565b60008060408385031215610bd057600080fd5b823567ffffffffffffffff80821115610be857600080fd5b610bf4868387016109a1565b93506020850135915080821115610c0a57600080fd5b50610c17858286016109a1565b9150509250929050565b60008251610c3381846020870161081d565b9190910192915050565b601f821115610c8757600081815260208120601f850160051c81016020861015610c645750805b601f850160051c820191505b81811015610c8357828155600101610c70565b5050505b505050565b815167ffffffffffffffff811115610ca657610ca6610972565b610cba81610cb48454610b23565b84610c3d565b602080601f831160018114610d0d5760008415610cd75750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610c83565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610d5a57888601518255948401946001909101908401610d3b565b5085821015610d9657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b848152836020820152608060408201526000610dc5608083018561084d565b8281036060840152610dd7818561084d565b979650505050505050565b6000808454610df081610b23565b60018281168015610e085760018114610e3b57610e6a565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168752821515830287019450610e6a565b8860005260208060002060005b85811015610e615781548a820152908401908201610e48565b50505082870194505b505050508351610e7e81836020880161081d565b01949350505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b608081526000610ee7608083018761084d565b602083820381850152610efa828861084d565b91508382036040850152818654808452828401915060058382821b86010160008a8152858120815b85811015610ff5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0898503018752828254610f5d81610b23565b80875260018281168015610f785760018114610faf57610fde565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168d8a01528c8315158b1b8a01019450610fde565b8688528c8820885b84811015610fd65781548f828d01015283820191508e81019050610fb7565b8a018e019550505b50998b019992965050509190910190600101610f22565b5050507fffffffff0000000000000000000000000000000000000000000000000000000089166060890152955061102d945050505050565b95945050505050565b60408152600061107360408301600581527f6572726f72000000000000000000000000000000000000000000000000000000602082015260400190565b82810360208401526108aa81600581527f6572726f72000000000000000000000000000000000000000000000000000000602082015260400190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261111357600080fd5b83018035915067ffffffffffffffff82111561112e57600080fd5b60200191503681900382131561053857600080fd5b604081526000611156604083018561084d565b828103602084015261102d818561084d565b60008282101561117a5761117a610b76565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036111b0576111b0610b76565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826111f5576111f56111b7565b500490565b600082611209576112096111b7565b50069056fea164736f6c634300080f000a68747470733a2f2f706f6b656170692e636f2f6170692f76322f706f6b656d6f6e2f",
}

var UpkeepAPIFetchABI = UpkeepAPIFetchMetaData.ABI

var UpkeepAPIFetchBin = UpkeepAPIFetchMetaData.Bin

func DeployUpkeepAPIFetch(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *UpkeepAPIFetch, error) {
	parsed, err := UpkeepAPIFetchMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepAPIFetchBin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepAPIFetch{UpkeepAPIFetchCaller: UpkeepAPIFetchCaller{contract: contract}, UpkeepAPIFetchTransactor: UpkeepAPIFetchTransactor{contract: contract}, UpkeepAPIFetchFilterer: UpkeepAPIFetchFilterer{contract: contract}}, nil
}

type UpkeepAPIFetch struct {
	address common.Address
	abi     abi.ABI
	UpkeepAPIFetchCaller
	UpkeepAPIFetchTransactor
	UpkeepAPIFetchFilterer
}

type UpkeepAPIFetchCaller struct {
	contract *bind.BoundContract
}

type UpkeepAPIFetchTransactor struct {
	contract *bind.BoundContract
}

type UpkeepAPIFetchFilterer struct {
	contract *bind.BoundContract
}

type UpkeepAPIFetchSession struct {
	Contract     *UpkeepAPIFetch
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepAPIFetchCallerSession struct {
	Contract *UpkeepAPIFetchCaller
	CallOpts bind.CallOpts
}

type UpkeepAPIFetchTransactorSession struct {
	Contract     *UpkeepAPIFetchTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepAPIFetchRaw struct {
	Contract *UpkeepAPIFetch
}

type UpkeepAPIFetchCallerRaw struct {
	Contract *UpkeepAPIFetchCaller
}

type UpkeepAPIFetchTransactorRaw struct {
	Contract *UpkeepAPIFetchTransactor
}

func NewUpkeepAPIFetch(address common.Address, backend bind.ContractBackend) (*UpkeepAPIFetch, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepAPIFetchABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepAPIFetch(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetch{address: address, abi: abi, UpkeepAPIFetchCaller: UpkeepAPIFetchCaller{contract: contract}, UpkeepAPIFetchTransactor: UpkeepAPIFetchTransactor{contract: contract}, UpkeepAPIFetchFilterer: UpkeepAPIFetchFilterer{contract: contract}}, nil
}

func NewUpkeepAPIFetchCaller(address common.Address, caller bind.ContractCaller) (*UpkeepAPIFetchCaller, error) {
	contract, err := bindUpkeepAPIFetch(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchCaller{contract: contract}, nil
}

func NewUpkeepAPIFetchTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepAPIFetchTransactor, error) {
	contract, err := bindUpkeepAPIFetch(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchTransactor{contract: contract}, nil
}

func NewUpkeepAPIFetchFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepAPIFetchFilterer, error) {
	contract, err := bindUpkeepAPIFetch(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchFilterer{contract: contract}, nil
}

func bindUpkeepAPIFetch(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepAPIFetchABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UpkeepAPIFetch *UpkeepAPIFetchRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepAPIFetch.Contract.UpkeepAPIFetchCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.UpkeepAPIFetchTransactor.contract.Transfer(opts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.UpkeepAPIFetchTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepAPIFetch.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.contract.Transfer(opts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Callback(opts *bind.CallOpts, extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "callback", extraData, values, statusCode)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Callback(extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.Callback(&_UpkeepAPIFetch.CallOpts, extraData, values, statusCode)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Callback(extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.Callback(&_UpkeepAPIFetch.CallOpts, extraData, values, statusCode)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.CheckUpkeep(&_UpkeepAPIFetch.CallOpts, data)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.CheckUpkeep(&_UpkeepAPIFetch.CallOpts, data)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Counter() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Counter(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Counter() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Counter(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Eligible() (bool, error) {
	return _UpkeepAPIFetch.Contract.Eligible(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Eligible() (bool, error) {
	return _UpkeepAPIFetch.Contract.Eligible(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Fields(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "fields", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Fields(arg0 *big.Int) (string, error) {
	return _UpkeepAPIFetch.Contract.Fields(&_UpkeepAPIFetch.CallOpts, arg0)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Fields(arg0 *big.Int) (string, error) {
	return _UpkeepAPIFetch.Contract.Fields(&_UpkeepAPIFetch.CallOpts, arg0)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Id(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "id")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Id() (string, error) {
	return _UpkeepAPIFetch.Contract.Id(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Id() (string, error) {
	return _UpkeepAPIFetch.Contract.Id(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) InitialBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.InitialBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) InitialBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.InitialBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Interval() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Interval(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Interval() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Interval(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) LastBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.LastBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) LastBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.LastBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Pokemon(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "pokemon")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Pokemon() (string, error) {
	return _UpkeepAPIFetch.Contract.Pokemon(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Pokemon() (string, error) {
	return _UpkeepAPIFetch.Contract.Pokemon(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.PreviousPerformBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.PreviousPerformBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) TestRange() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.TestRange(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.TestRange(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Url(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "url")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Url() (string, error) {
	return _UpkeepAPIFetch.Contract.Url(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Url() (string, error) {
	return _UpkeepAPIFetch.Contract.Url(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _UpkeepAPIFetch.contract.Transact(opts, "performUpkeep", performData)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.PerformUpkeep(&_UpkeepAPIFetch.TransactOpts, performData)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.PerformUpkeep(&_UpkeepAPIFetch.TransactOpts, performData)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactor) SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepAPIFetch.contract.Transact(opts, "setConfig", _testRange, _interval)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetConfig(&_UpkeepAPIFetch.TransactOpts, _testRange, _interval)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetConfig(&_UpkeepAPIFetch.TransactOpts, _testRange, _interval)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactor) SetURLs(opts *bind.TransactOpts, input string) (*types.Transaction, error) {
	return _UpkeepAPIFetch.contract.Transact(opts, "setURLs", input)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) SetURLs(input string) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetURLs(&_UpkeepAPIFetch.TransactOpts, input)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorSession) SetURLs(input string) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetURLs(&_UpkeepAPIFetch.TransactOpts, input)
}

type UpkeepAPIFetchPerformingUpkeepIterator struct {
	Event *UpkeepAPIFetchPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepAPIFetchPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepAPIFetchPerformingUpkeep)
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
		it.Event = new(UpkeepAPIFetchPerformingUpkeep)
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

func (it *UpkeepAPIFetchPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *UpkeepAPIFetchPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepAPIFetchPerformingUpkeep struct {
	From      common.Address
	LastBlock *big.Int
	Counter   *big.Int
	Id        string
	Name      string
	Raw       types.Log
}

func (_UpkeepAPIFetch *UpkeepAPIFetchFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepAPIFetchPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepAPIFetch.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchPerformingUpkeepIterator{contract: _UpkeepAPIFetch.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_UpkeepAPIFetch *UpkeepAPIFetchFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepAPIFetchPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepAPIFetch.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepAPIFetchPerformingUpkeep)
				if err := _UpkeepAPIFetch.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_UpkeepAPIFetch *UpkeepAPIFetchFilterer) ParsePerformingUpkeep(log types.Log) (*UpkeepAPIFetchPerformingUpkeep, error) {
	event := new(UpkeepAPIFetchPerformingUpkeep)
	if err := _UpkeepAPIFetch.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_UpkeepAPIFetch *UpkeepAPIFetch) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepAPIFetch.abi.Events["PerformingUpkeep"].ID:
		return _UpkeepAPIFetch.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepAPIFetchPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x0b36ab4eaaf3df8d059b63746075640cc0f87f07d011137505160b1de3c983da")
}

func (_UpkeepAPIFetch *UpkeepAPIFetch) Address() common.Address {
	return _UpkeepAPIFetch.address
}

type UpkeepAPIFetchInterface interface {
	Callback(opts *bind.CallOpts, extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	Fields(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	Id(opts *bind.CallOpts) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	Pokemon(opts *bind.CallOpts) (string, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	Url(opts *bind.CallOpts) (string, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetURLs(opts *bind.TransactOpts, input string) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepAPIFetchPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepAPIFetchPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepAPIFetchPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
