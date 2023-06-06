package ethereum

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
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
)

var PerformDataCheckerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"expectedData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_expectedData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"expectedData\",\"type\":\"bytes\"}],\"name\":\"setExpectedData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516107d33803806107d383398101604081905261002f916100e2565b8051610042906001906020840190610049565b5050610202565b828054610055906101b1565b90600052602060002090601f01602090048101928261007757600085556100bd565b82601f1061009057805160ff19168380011785556100bd565b828001600101855582156100bd579182015b828111156100bd5782518255916020019190600101906100a2565b506100c99291506100cd565b5090565b5b808211156100c957600081556001016100ce565b600060208083850312156100f557600080fd5b82516001600160401b038082111561010c57600080fd5b818501915085601f83011261012057600080fd5b815181811115610132576101326101ec565b604051601f8201601f19908116603f0116810190838211818310171561015a5761015a6101ec565b81604052828152888684870101111561017257600080fd5b600093505b828410156101945784840186015181850187015292850192610177565b828411156101a55760008684830101525b98975050505050505050565b600181811c908216806101c557607f821691505b602082108114156101e657634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b6105c2806102116000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806361bc221a1161005057806361bc221a146100945780636e04ff0d146100b05780638d1a93c2146100d157600080fd5b80632aa0f7951461006c5780634585e33b14610081575b600080fd5b61007f61007a366004610304565b6100e6565b005b61007f61008f366004610304565b6100f7565b61009d60005481565b6040519081526020015b60405180910390f35b6100c36100be366004610304565b610145565b6040516100a79291906104c4565b6100d96101bf565b6040516100a791906104e7565b6100f26001838361024d565b505050565b600160405161010691906103f1565b6040518091039020828260405161011e9291906103e1565b604051809103902014156101415760008054908061013b83610555565b91905055505b5050565b60006060600160405161015891906103f1565b604051809103902084846040516101709291906103e1565b604051809103902014848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b600180546101cc90610501565b80601f01602080910402602001604051908101604052809291908181526020018280546101f890610501565b80156102455780601f1061021a57610100808354040283529160200191610245565b820191906000526020600020905b81548152906001019060200180831161022857829003601f168201915b505050505081565b82805461025990610501565b90600052602060002090601f01602090048101928261027b57600085556102df565b82601f106102b2578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008235161785556102df565b828001600101855582156102df579182015b828111156102df5782358255916020019190600101906102c4565b506102eb9291506102ef565b5090565b5b808211156102eb57600081556001016102f0565b6000806020838503121561031757600080fd5b823567ffffffffffffffff8082111561032f57600080fd5b818501915085601f83011261034357600080fd5b81358181111561035257600080fd5b86602082850101111561036457600080fd5b60209290920196919550909350505050565b6000815180845260005b8181101561039c57602081850181015186830182015201610380565b818111156103ae576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8183823760009101908152919050565b600080835481600182811c91508083168061040d57607f831692505b6020808410821415610446577f4e487b710000000000000000000000000000000000000000000000000000000086526022600452602486fd5b81801561045a5760018114610489576104b6565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616895284890196506104b6565b60008a81526020902060005b868110156104ae5781548b820152908501908301610495565b505084890196505b509498975050505050505050565b82151581526040602082015260006104df6040830184610376565b949350505050565b6020815260006104fa6020830184610376565b9392505050565b600181811c9082168061051557607f821691505b6020821081141561054f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156105ae577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fea164736f6c6343000806000a",
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
	return address, tx, &PerformDataChecker{PerformDataCheckerCaller: PerformDataCheckerCaller{contract: contract}, PerformDataCheckerTransactor: PerformDataCheckerTransactor{contract: contract}, PerformDataCheckerFilterer: PerformDataCheckerFilterer{contract: contract}}, nil
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
	parsed, err := abi.JSON(strings.NewReader(PerformDataCheckerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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
