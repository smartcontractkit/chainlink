// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package perform_data_checker_wrapper

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

var PerformDataCheckerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"expectedData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_expectedData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"expectedData\",\"type\":\"bytes\"}],\"name\":\"setExpectedData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516108d33803806108d383398101604081905261002f91610058565b600161003b82826101aa565b5050610269565b634e487b7160e01b600052604160045260246000fd5b6000602080838503121561006b57600080fd5b82516001600160401b038082111561008257600080fd5b818501915085601f83011261009657600080fd5b8151818111156100a8576100a8610042565b604051601f8201601f19908116603f011681019083821181831017156100d0576100d0610042565b8160405282815288868487010111156100e857600080fd5b600093505b8284101561010a57848401860151818501870152928501926100ed565b600086848301015280965050505050505092915050565b600181811c9082168061013557607f821691505b60208210810361015557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101a557600081815260208120601f850160051c810160208610156101825750805b601f850160051c820191505b818110156101a15782815560010161018e565b5050505b505050565b81516001600160401b038111156101c3576101c3610042565b6101d7816101d18454610121565b8461015b565b602080601f83116001811461020c57600084156101f45750858301515b600019600386901b1c1916600185901b1785556101a1565b600085815260208120601f198616915b8281101561023b5788860151825594840194600190910190840161021c565b50858210156102595787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61065b806102786000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806361bc221a1161005057806361bc221a146100945780636e04ff0d146100b05780638d1a93c2146100d157600080fd5b80632aa0f7951461006c5780634585e33b14610081575b600080fd5b61007f61007a36600461024d565b6100e6565b005b61007f61008f36600461024d565b6100f8565b61009d60005481565b6040519081526020015b60405180910390f35b6100c36100be36600461024d565b610145565b6040516100a7929190610323565b6100d96101bf565b6040516100a79190610346565b60016100f3828483610430565b505050565b6001604051610107919061054b565b6040518091039020828260405161011f9291906105df565b6040518091039020036101415760008054908061013b836105ef565b91905055505b5050565b600060606001604051610158919061054b565b604051809103902084846040516101709291906105df565b604051809103902014848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b600180546101cc9061038f565b80601f01602080910402602001604051908101604052809291908181526020018280546101f89061038f565b80156102455780601f1061021a57610100808354040283529160200191610245565b820191906000526020600020905b81548152906001019060200180831161022857829003601f168201915b505050505081565b6000806020838503121561026057600080fd5b823567ffffffffffffffff8082111561027857600080fd5b818501915085601f83011261028c57600080fd5b81358181111561029b57600080fd5b8660208285010111156102ad57600080fd5b60209290920196919550909350505050565b6000815180845260005b818110156102e5576020818501810151868301820152016102c9565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b821515815260406020820152600061033e60408301846102bf565b949350505050565b60208152600061035960208301846102bf565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600181811c908216806103a357607f821691505b6020821081036103dc577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156100f357600081815260208120601f850160051c810160208610156104095750805b601f850160051c820191505b8181101561042857828155600101610415565b505050505050565b67ffffffffffffffff83111561044857610448610360565b61045c83610456835461038f565b836103e2565b6000601f8411600181146104ae57600085156104785750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355610544565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156104fd57868501358255602094850194600190920191016104dd565b5086821015610538577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60008083546105598161038f565b6001828116801561057157600181146105a4576105d3565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841687528215158302870194506105d3565b8760005260208060002060005b858110156105ca5781548a8201529084019082016105b1565b50505082870194505b50929695505050505050565b8183823760009101908152919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610647577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fea164736f6c6343000810000a",
}

var PerformDataCheckerABI = PerformDataCheckerMetaData.ABI

var PerformDataCheckerBin = PerformDataCheckerMetaData.Bin

func DeployPerformDataChecker(auth *bind.TransactOpts, backend bind.ContractBackend, expectedData []byte) (common.Address, *types.Transaction, *PerformDataChecker, error) {
	parsed, err := PerformDataCheckerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PerformDataCheckerBin), backend, expectedData)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PerformDataChecker{address: address, abi: *parsed, PerformDataCheckerCaller: PerformDataCheckerCaller{contract: contract}, PerformDataCheckerTransactor: PerformDataCheckerTransactor{contract: contract}, PerformDataCheckerFilterer: PerformDataCheckerFilterer{contract: contract}}, nil
}

type PerformDataChecker struct {
	address common.Address
	abi     abi.ABI
	PerformDataCheckerCaller
	PerformDataCheckerTransactor
	PerformDataCheckerFilterer
}

type PerformDataCheckerCaller struct {
	contract *bind.BoundContract
}

type PerformDataCheckerTransactor struct {
	contract *bind.BoundContract
}

type PerformDataCheckerFilterer struct {
	contract *bind.BoundContract
}

type PerformDataCheckerSession struct {
	Contract     *PerformDataChecker
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PerformDataCheckerCallerSession struct {
	Contract *PerformDataCheckerCaller
	CallOpts bind.CallOpts
}

type PerformDataCheckerTransactorSession struct {
	Contract     *PerformDataCheckerTransactor
	TransactOpts bind.TransactOpts
}

type PerformDataCheckerRaw struct {
	Contract *PerformDataChecker
}

type PerformDataCheckerCallerRaw struct {
	Contract *PerformDataCheckerCaller
}

type PerformDataCheckerTransactorRaw struct {
	Contract *PerformDataCheckerTransactor
}

func NewPerformDataChecker(address common.Address, backend bind.ContractBackend) (*PerformDataChecker, error) {
	abi, err := abi.JSON(strings.NewReader(PerformDataCheckerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPerformDataChecker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PerformDataChecker{address: address, abi: abi, PerformDataCheckerCaller: PerformDataCheckerCaller{contract: contract}, PerformDataCheckerTransactor: PerformDataCheckerTransactor{contract: contract}, PerformDataCheckerFilterer: PerformDataCheckerFilterer{contract: contract}}, nil
}

func NewPerformDataCheckerCaller(address common.Address, caller bind.ContractCaller) (*PerformDataCheckerCaller, error) {
	contract, err := bindPerformDataChecker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PerformDataCheckerCaller{contract: contract}, nil
}

func NewPerformDataCheckerTransactor(address common.Address, transactor bind.ContractTransactor) (*PerformDataCheckerTransactor, error) {
	contract, err := bindPerformDataChecker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PerformDataCheckerTransactor{contract: contract}, nil
}

func NewPerformDataCheckerFilterer(address common.Address, filterer bind.ContractFilterer) (*PerformDataCheckerFilterer, error) {
	contract, err := bindPerformDataChecker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PerformDataCheckerFilterer{contract: contract}, nil
}

func bindPerformDataChecker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PerformDataCheckerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_PerformDataChecker *PerformDataCheckerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerformDataChecker.Contract.PerformDataCheckerCaller.contract.Call(opts, result, method, params...)
}

func (_PerformDataChecker *PerformDataCheckerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.PerformDataCheckerTransactor.contract.Transfer(opts)
}

func (_PerformDataChecker *PerformDataCheckerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.PerformDataCheckerTransactor.contract.Transact(opts, method, params...)
}

func (_PerformDataChecker *PerformDataCheckerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerformDataChecker.Contract.contract.Call(opts, result, method, params...)
}

func (_PerformDataChecker *PerformDataCheckerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.contract.Transfer(opts)
}

func (_PerformDataChecker *PerformDataCheckerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.contract.Transact(opts, method, params...)
}

func (_PerformDataChecker *PerformDataCheckerCaller) CheckUpkeep(opts *bind.CallOpts, checkData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _PerformDataChecker.contract.Call(opts, &out, "checkUpkeep", checkData)

	outstruct := new(CheckUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_PerformDataChecker *PerformDataCheckerSession) CheckUpkeep(checkData []byte) (CheckUpkeep,

	error) {
	return _PerformDataChecker.Contract.CheckUpkeep(&_PerformDataChecker.CallOpts, checkData)
}

func (_PerformDataChecker *PerformDataCheckerCallerSession) CheckUpkeep(checkData []byte) (CheckUpkeep,

	error) {
	return _PerformDataChecker.Contract.CheckUpkeep(&_PerformDataChecker.CallOpts, checkData)
}

func (_PerformDataChecker *PerformDataCheckerCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PerformDataChecker.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PerformDataChecker *PerformDataCheckerSession) Counter() (*big.Int, error) {
	return _PerformDataChecker.Contract.Counter(&_PerformDataChecker.CallOpts)
}

func (_PerformDataChecker *PerformDataCheckerCallerSession) Counter() (*big.Int, error) {
	return _PerformDataChecker.Contract.Counter(&_PerformDataChecker.CallOpts)
}

func (_PerformDataChecker *PerformDataCheckerCaller) SExpectedData(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PerformDataChecker.contract.Call(opts, &out, "s_expectedData")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_PerformDataChecker *PerformDataCheckerSession) SExpectedData() ([]byte, error) {
	return _PerformDataChecker.Contract.SExpectedData(&_PerformDataChecker.CallOpts)
}

func (_PerformDataChecker *PerformDataCheckerCallerSession) SExpectedData() ([]byte, error) {
	return _PerformDataChecker.Contract.SExpectedData(&_PerformDataChecker.CallOpts)
}

func (_PerformDataChecker *PerformDataCheckerTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _PerformDataChecker.contract.Transact(opts, "performUpkeep", performData)
}

func (_PerformDataChecker *PerformDataCheckerSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.PerformUpkeep(&_PerformDataChecker.TransactOpts, performData)
}

func (_PerformDataChecker *PerformDataCheckerTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.PerformUpkeep(&_PerformDataChecker.TransactOpts, performData)
}

func (_PerformDataChecker *PerformDataCheckerTransactor) SetExpectedData(opts *bind.TransactOpts, expectedData []byte) (*types.Transaction, error) {
	return _PerformDataChecker.contract.Transact(opts, "setExpectedData", expectedData)
}

func (_PerformDataChecker *PerformDataCheckerSession) SetExpectedData(expectedData []byte) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.SetExpectedData(&_PerformDataChecker.TransactOpts, expectedData)
}

func (_PerformDataChecker *PerformDataCheckerTransactorSession) SetExpectedData(expectedData []byte) (*types.Transaction, error) {
	return _PerformDataChecker.Contract.SetExpectedData(&_PerformDataChecker.TransactOpts, expectedData)
}

type CheckUpkeep struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_PerformDataChecker *PerformDataChecker) Address() common.Address {
	return _PerformDataChecker.address
}

type PerformDataCheckerInterface interface {
	CheckUpkeep(opts *bind.CallOpts, checkData []byte) (CheckUpkeep,

		error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	SExpectedData(opts *bind.CallOpts) ([]byte, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetExpectedData(opts *bind.TransactOpts, expectedData []byte) (*types.Transaction, error)

	Address() common.Address
}
