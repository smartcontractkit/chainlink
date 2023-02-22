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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"jsonFields\",\"type\":\"string[]\"},{\"internalType\":\"bytes4\",\"name\":\"callbackSelector\",\"type\":\"bytes4\"}],\"name\":\"ChainlinkAPIFetch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"gender\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"values\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"statusCode\",\"type\":\"uint256\"}],\"name\":\"callback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fields\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"input\",\"type\":\"string\"}],\"name\":\"setURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"s\",\"type\":\"string\"}],\"name\":\"stringToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"url\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200111538038062001115833981016040819052620000349162000247565b60008281556001829055600381905543600290815560048281556005929092556040805160808101825260068183019081526533b2b73232b960d11b606083015281528151808301909252928152636e616d6560e01b602082810191909152830152620000a491600791620000de565b50604051806060016040528060248152602001620010f1602491398051620000d59160069160209091019062000142565b505050620002a9565b82805482825590600052602060002090810192821562000130579160200282015b828111156200013057825180516200011f91849160209091019062000142565b5091602001919060010190620000ff565b506200013e929150620001cd565b5090565b82805462000150906200026c565b90600052602060002090601f016020900481019282620001745760008555620001bf565b82601f106200018f57805160ff1916838001178555620001bf565b82800160010185558215620001bf579182015b82811115620001bf578251825591602001919060010190620001a2565b506200013e929150620001ee565b808211156200013e576000620001e4828262000205565b50600101620001cd565b5b808211156200013e5760008155600101620001ef565b50805462000213906200026c565b6000825580601f1062000224575050565b601f016020900490600052602060002090810190620002449190620001ee565b50565b600080604083850312156200025b57600080fd5b505080516020909101519092909150565b600181811c908216806200028157607f821691505b60208210811415620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b610e3880620002b96000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80636250a13a11610097578063917d895f11610066578063917d895f146101d9578063947a36fb146101e2578063b772d70a146101eb578063d832d92f146101fe57600080fd5b80636250a13a146101935780636e04ff0d1461019c57806371106628146101bd578063806b984f146101d057600080fd5b80632cb15864116100d35780632cb15864146101665780634585e33b1461016f5780635600f04f1461018257806361bc221a1461018a57600080fd5b806319137856146100fa5780631bd95155146101235780631e34c58514610144575b600080fd5b61010d6101083660046108f4565b610216565b60405161011a9190610ada565b60405180910390f35b610136610131366004610853565b6102c2565b60405190815260200161011a565b61016461015236600461090d565b60009182556001556004819055600555565b005b61013660045481565b61016461017d36600461076f565b610344565b61010d6103e0565b61013660055481565b61013660005481565b6101af6101aa36600461076f565b6103ed565b60405161011a929190610a72565b6101646101cb366004610853565b610493565b61013660025481565b61013660035481565b61013660015481565b6101af6101f93660046107b1565b6104aa565b61020661059f565b604051901515815260200161011a565b6007818154811061022657600080fd5b90600052602060002001600091509050805461024190610d11565b80601f016020809104026020016040519081016040528092919081815260200182805461026d90610d11565b80156102ba5780601f1061028f576101008083540402835291602001916102ba565b820191906000526020600020905b81548152906001019060200180831161029d57829003601f168201915b505050505081565b60008181805b825181101561033c5760008382815181106102e5576102e5610dcd565b016020015160f81c905060308110801590610301575060398111155b1561032957610311603082610cfa565b61031c84600a610cbd565b6103269190610ca5565b92505b508061033481610d65565b9150506102c8565b509392505050565b60045461035057436004555b60055461035e906001610ca5565b60055560008061037083850185610890565b915091503273ffffffffffffffffffffffffffffffffffffffff167f4892f73cdece87001a2b32299bd0ac40d3f60c4faec566074057479dbfebc4ed60045460025460035460055487876040516103cc96959493929190610bf6565b60405180910390a250506002546003555050565b6006805461024190610d11565b6000606060068484604051602001610406929190610a8d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527fda8b321400000000000000000000000000000000000000000000000000000000825261048a92916007907fb772d70a0000000000000000000000000000000000000000000000000000000090600401610b22565b60405180910390fd5b80516104a69060069060208401906105e2565b5050565b600060606000858560008181106104c3576104c3610dcd565b90506020028101906104d59190610c40565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509394508992508891506001905081811061052157610521610dcd565b90506020028101906105339190610c40565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060405192935060019261058092508591508490602001610af4565b6040516020818303038152906040529350935050509550959350505050565b6000600454600014156105b25750600190565b6000546004546105c29043610cfa565b1080156105dd57506001546002546105da9043610cfa565b10155b905090565b8280546105ee90610d11565b90600052602060002090601f0160209004810192826106105760008555610656565b82601f1061062957805160ff1916838001178555610656565b82800160010185558215610656579182015b8281111561065657825182559160200191906001019061063b565b50610662929150610666565b5090565b5b808211156106625760008155600101610667565b60008083601f84011261068d57600080fd5b50813567ffffffffffffffff8111156106a557600080fd5b6020830191508360208285010111156106bd57600080fd5b9250929050565b600082601f8301126106d557600080fd5b813567ffffffffffffffff808211156106f0576106f0610dfc565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561073657610736610dfc565b8160405283815286602085880101111561074f57600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806020838503121561078257600080fd5b823567ffffffffffffffff81111561079957600080fd5b6107a58582860161067b565b90969095509350505050565b6000806000806000606086880312156107c957600080fd5b853567ffffffffffffffff808211156107e157600080fd5b6107ed89838a0161067b565b9097509550602088013591508082111561080657600080fd5b818801915088601f83011261081a57600080fd5b81358181111561082957600080fd5b8960208260051b850101111561083e57600080fd5b96999598505060200195604001359392505050565b60006020828403121561086557600080fd5b813567ffffffffffffffff81111561087c57600080fd5b610888848285016106c4565b949350505050565b600080604083850312156108a357600080fd5b823567ffffffffffffffff808211156108bb57600080fd5b6108c7868387016106c4565b935060208501359150808211156108dd57600080fd5b506108ea858286016106c4565b9150509250929050565b60006020828403121561090657600080fd5b5035919050565b6000806040838503121561092057600080fd5b50508035926020909101359150565b6000815180845260005b8181101561095557602081850181015186830182015201610939565b81811115610967576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8054600090600181811c90808316806109b457607f831692505b60208084108214156109ef577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b838852818015610a065760018114610a3857610a66565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616828a0152604089019650610a66565b876000528160002060005b86811015610a5e5781548b8201850152908501908301610a43565b8a0183019750505b50505050505092915050565b8215158152604060208201526000610888604083018461092f565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b602081526000610aed602083018461092f565b9392505050565b604081526000610b07604083018561092f565b8281036020840152610b19818561092f565b95945050505050565b608081526000610b35608083018761099a565b602083820381850152610b48828861092f565b915083820360408501528186548084528284019150828160051b850101886000528360002060005b83811015610bbb577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552610ba9838361099a565b94860194925060019182019101610b70565b505080955050505050507fffffffff000000000000000000000000000000000000000000000000000000008316606083015295945050505050565b86815285602082015284604082015283606082015260c060808201526000610c2160c083018561092f565b82810360a0840152610c33818561092f565b9998505050505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610c7557600080fd5b83018035915067ffffffffffffffff821115610c9057600080fd5b6020019150368190038213156106bd57600080fd5b60008219821115610cb857610cb8610d9e565b500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610cf557610cf5610d9e565b500290565b600082821015610d0c57610d0c610d9e565b500390565b600181811c90821680610d2557607f821691505b60208210811415610d5f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610d9757610d97610d9e565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a68747470733a2f2f6170692e67656e646572697a652e696f2f3f6e616d653d6368726973",
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

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) StringToUint(opts *bind.CallOpts, s string) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "stringToUint", s)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) StringToUint(s string) (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.StringToUint(&_UpkeepAPIFetch.CallOpts, s)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) StringToUint(s string) (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.StringToUint(&_UpkeepAPIFetch.CallOpts, s)
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
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	Gender        string
	Name          string
	Raw           types.Log
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
	return common.HexToHash("0x4892f73cdece87001a2b32299bd0ac40d3f60c4faec566074057479dbfebc4ed")
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

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	StringToUint(opts *bind.CallOpts, s string) (*big.Int, error)

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
