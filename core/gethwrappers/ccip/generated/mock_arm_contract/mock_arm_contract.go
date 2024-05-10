// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_arm_contract

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

type IRMNTaggedRoot struct {
	CommitStore common.Address
	Root        [32]byte
}

type RMNConfig struct {
	Voters               []RMNVoter
	BlessWeightThreshold uint16
	CurseWeightThreshold uint16
}

type RMNUnvoteToCurseRecord struct {
	CurseVoteAddr common.Address
	CursesHash    [32]byte
	ForceUnvote   bool
}

type RMNVoter struct {
	BlessVoteAddr   common.Address
	CurseVoteAddr   common.Address
	CurseUnvoteAddr common.Address
	BlessWeight     uint8
	CurseWeight     uint8
}

var MockARMContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"CustomError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"blessVoteAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"curseVoteAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"curseUnvoteAddr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"blessWeight\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"curseWeight\",\"type\":\"uint8\"}],\"internalType\":\"structRMN.Voter[]\",\"name\":\"voters\",\"type\":\"tuple[]\"},{\"internalType\":\"uint16\",\"name\":\"blessWeightThreshold\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"curseWeightThreshold\",\"type\":\"uint16\"}],\"internalType\":\"structRMN.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMN.TaggedRoot\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"isBlessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"curseVoteAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"cursesHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"forceUnvote\",\"type\":\"bool\"}],\"internalType\":\"structRMN.UnvoteToCurseRecord[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"ownerUnvoteToCurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"setRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"voteToCurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610ccf806101576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063618af128116100765780637a7c27491161005b5780637a7c2749146101ce5780638da5cb5b146101e1578063f2fde38b1461020957600080fd5b8063618af1281461018d57806379ba5097146101c657600080fd5b80633f42ab73116100a75780633f42ab73146101325780634d6167711461014957806358babe331461017a57600080fd5b8063119a3527146100c3578063397796f714610115575b600080fd5b6101136100d13660046106d1565b50600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055565b005b61011d61021c565b60405190151581526020015b60405180910390f35b61013a610292565b604051610129939291906106ea565b61011d6101573660046107bc565b5060015474010000000000000000000000000000000000000000900460ff161590565b61011d6101883660046106d1565b6103ca565b61011361019b3660046108a4565b50600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff169055565b610113610438565b6101136101dc36600461098e565b610535565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610129565b610113610217366004610a41565b610545565b60006002805461022b90610a63565b1590506102705760026040517f5a4ff6710000000000000000000000000000000000000000000000000000000081526004016102679190610ab0565b60405180910390fd5b5060015474010000000000000000000000000000000000000000900460ff1690565b6040805160608082018352808252600060208084018290528385018290526005548551600380549384028201608090810190985294810183815263ffffffff808416986401000000009094041696959194919385939192859285015b8282101561039a5760008481526020908190206040805160a08101825260038602909201805473ffffffffffffffffffffffffffffffffffffffff90811684526001808301548216858701526002909201549081169284019290925260ff740100000000000000000000000000000000000000008304811660608501527501000000000000000000000000000000000000000000909204909116608083015290835290920191016102ee565b505050908252506001919091015461ffff8082166020840152620100009091041660409091015292939192919050565b6000600280546103d990610a63565b1590506104155760026040517f5a4ff6710000000000000000000000000000000000000000000000000000000081526004016102679190610ab0565b505060015474010000000000000000000000000000000000000000900460ff1690565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610267565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60026105418282610ba8565b5050565b61054d610559565b610556816105dc565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610267565b565b3373ffffffffffffffffffffffffffffffffffffffff82160361065b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610267565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156106e357600080fd5b5035919050565b63ffffffff84811682528316602080830191909152606060408084018290528451848301839052805160c0860181905260009491820190859060e08801905b80831015610790578351805173ffffffffffffffffffffffffffffffffffffffff9081168452868201518116878501528782015116878401528781015160ff908116898501526080918201511690830152928401926001929092019160a090910190610729565b509288015161ffff9081166080890152939097015190921660a090950194909452509195945050505050565b6000604082840312156107ce57600080fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610826576108266107d4565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610873576108736107d4565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461089f57600080fd5b919050565b600060208083850312156108b757600080fd5b823567ffffffffffffffff808211156108cf57600080fd5b818501915085601f8301126108e357600080fd5b8135818111156108f5576108f56107d4565b610903848260051b0161082c565b8181528481019250606091820284018501918883111561092257600080fd5b938501935b828510156109825780858a03121561093f5760008081fd5b610947610803565b6109508661087b565b81528686013587820152604080870135801515811461096f5760008081fd5b9082015284529384019392850192610927565b50979650505050505050565b600060208083850312156109a157600080fd5b823567ffffffffffffffff808211156109b957600080fd5b818501915085601f8301126109cd57600080fd5b8135818111156109df576109df6107d4565b610a0f847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161082c565b91508082528684828501011115610a2557600080fd5b8084840185840137600090820190930192909252509392505050565b600060208284031215610a5357600080fd5b610a5c8261087b565b9392505050565b600181811c90821680610a7757607f821691505b6020821081036107ce577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000602080835260008454610ac481610a63565b80848701526040600180841660008114610ae55760018114610b1d57610b4b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a01019550610b4b565b896000528660002060005b85811015610b435781548b8201860152908301908801610b28565b8a0184019650505b509398975050505050505050565b601f821115610ba357600081815260208120601f850160051c81016020861015610b805750805b601f850160051c820191505b81811015610b9f57828155600101610b8c565b5050505b505050565b815167ffffffffffffffff811115610bc257610bc26107d4565b610bd681610bd08454610a63565b84610b59565b602080601f831160018114610c295760008415610bf35750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610b9f565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610c7657888601518255948401946001909101908401610c57565b5085821015610cb257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000813000a",
}

var MockARMContractABI = MockARMContractMetaData.ABI

var MockARMContractBin = MockARMContractMetaData.Bin

func DeployMockARMContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockARMContract, error) {
	parsed, err := MockARMContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockARMContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockARMContract{address: address, abi: *parsed, MockARMContractCaller: MockARMContractCaller{contract: contract}, MockARMContractTransactor: MockARMContractTransactor{contract: contract}, MockARMContractFilterer: MockARMContractFilterer{contract: contract}}, nil
}

type MockARMContract struct {
	address common.Address
	abi     abi.ABI
	MockARMContractCaller
	MockARMContractTransactor
	MockARMContractFilterer
}

type MockARMContractCaller struct {
	contract *bind.BoundContract
}

type MockARMContractTransactor struct {
	contract *bind.BoundContract
}

type MockARMContractFilterer struct {
	contract *bind.BoundContract
}

type MockARMContractSession struct {
	Contract     *MockARMContract
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockARMContractCallerSession struct {
	Contract *MockARMContractCaller
	CallOpts bind.CallOpts
}

type MockARMContractTransactorSession struct {
	Contract     *MockARMContractTransactor
	TransactOpts bind.TransactOpts
}

type MockARMContractRaw struct {
	Contract *MockARMContract
}

type MockARMContractCallerRaw struct {
	Contract *MockARMContractCaller
}

type MockARMContractTransactorRaw struct {
	Contract *MockARMContractTransactor
}

func NewMockARMContract(address common.Address, backend bind.ContractBackend) (*MockARMContract, error) {
	abi, err := abi.JSON(strings.NewReader(MockARMContractABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockARMContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockARMContract{address: address, abi: abi, MockARMContractCaller: MockARMContractCaller{contract: contract}, MockARMContractTransactor: MockARMContractTransactor{contract: contract}, MockARMContractFilterer: MockARMContractFilterer{contract: contract}}, nil
}

func NewMockARMContractCaller(address common.Address, caller bind.ContractCaller) (*MockARMContractCaller, error) {
	contract, err := bindMockARMContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockARMContractCaller{contract: contract}, nil
}

func NewMockARMContractTransactor(address common.Address, transactor bind.ContractTransactor) (*MockARMContractTransactor, error) {
	contract, err := bindMockARMContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockARMContractTransactor{contract: contract}, nil
}

func NewMockARMContractFilterer(address common.Address, filterer bind.ContractFilterer) (*MockARMContractFilterer, error) {
	contract, err := bindMockARMContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockARMContractFilterer{contract: contract}, nil
}

func bindMockARMContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockARMContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockARMContract *MockARMContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockARMContract.Contract.MockARMContractCaller.contract.Call(opts, result, method, params...)
}

func (_MockARMContract *MockARMContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockARMContract.Contract.MockARMContractTransactor.contract.Transfer(opts)
}

func (_MockARMContract *MockARMContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockARMContract.Contract.MockARMContractTransactor.contract.Transact(opts, method, params...)
}

func (_MockARMContract *MockARMContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockARMContract.Contract.contract.Call(opts, result, method, params...)
}

func (_MockARMContract *MockARMContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockARMContract.Contract.contract.Transfer(opts)
}

func (_MockARMContract *MockARMContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockARMContract.Contract.contract.Transact(opts, method, params...)
}

func (_MockARMContract *MockARMContractCaller) GetConfigDetails(opts *bind.CallOpts) (GetConfigDetails,

	error) {
	var out []interface{}
	err := _MockARMContract.contract.Call(opts, &out, "getConfigDetails")

	outstruct := new(GetConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Config = *abi.ConvertType(out[2], new(RMNConfig)).(*RMNConfig)

	return *outstruct, err

}

func (_MockARMContract *MockARMContractSession) GetConfigDetails() (GetConfigDetails,

	error) {
	return _MockARMContract.Contract.GetConfigDetails(&_MockARMContract.CallOpts)
}

func (_MockARMContract *MockARMContractCallerSession) GetConfigDetails() (GetConfigDetails,

	error) {
	return _MockARMContract.Contract.GetConfigDetails(&_MockARMContract.CallOpts)
}

func (_MockARMContract *MockARMContractCaller) IsBlessed(opts *bind.CallOpts, arg0 IRMNTaggedRoot) (bool, error) {
	var out []interface{}
	err := _MockARMContract.contract.Call(opts, &out, "isBlessed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockARMContract *MockARMContractSession) IsBlessed(arg0 IRMNTaggedRoot) (bool, error) {
	return _MockARMContract.Contract.IsBlessed(&_MockARMContract.CallOpts, arg0)
}

func (_MockARMContract *MockARMContractCallerSession) IsBlessed(arg0 IRMNTaggedRoot) (bool, error) {
	return _MockARMContract.Contract.IsBlessed(&_MockARMContract.CallOpts, arg0)
}

func (_MockARMContract *MockARMContractCaller) IsCursed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MockARMContract.contract.Call(opts, &out, "isCursed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockARMContract *MockARMContractSession) IsCursed() (bool, error) {
	return _MockARMContract.Contract.IsCursed(&_MockARMContract.CallOpts)
}

func (_MockARMContract *MockARMContractCallerSession) IsCursed() (bool, error) {
	return _MockARMContract.Contract.IsCursed(&_MockARMContract.CallOpts)
}

func (_MockARMContract *MockARMContractCaller) IsCursed0(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _MockARMContract.contract.Call(opts, &out, "isCursed0", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockARMContract *MockARMContractSession) IsCursed0(arg0 [32]byte) (bool, error) {
	return _MockARMContract.Contract.IsCursed0(&_MockARMContract.CallOpts, arg0)
}

func (_MockARMContract *MockARMContractCallerSession) IsCursed0(arg0 [32]byte) (bool, error) {
	return _MockARMContract.Contract.IsCursed0(&_MockARMContract.CallOpts, arg0)
}

func (_MockARMContract *MockARMContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockARMContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockARMContract *MockARMContractSession) Owner() (common.Address, error) {
	return _MockARMContract.Contract.Owner(&_MockARMContract.CallOpts)
}

func (_MockARMContract *MockARMContractCallerSession) Owner() (common.Address, error) {
	return _MockARMContract.Contract.Owner(&_MockARMContract.CallOpts)
}

func (_MockARMContract *MockARMContractTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockARMContract.contract.Transact(opts, "acceptOwnership")
}

func (_MockARMContract *MockARMContractSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockARMContract.Contract.AcceptOwnership(&_MockARMContract.TransactOpts)
}

func (_MockARMContract *MockARMContractTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockARMContract.Contract.AcceptOwnership(&_MockARMContract.TransactOpts)
}

func (_MockARMContract *MockARMContractTransactor) OwnerUnvoteToCurse(opts *bind.TransactOpts, arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error) {
	return _MockARMContract.contract.Transact(opts, "ownerUnvoteToCurse", arg0)
}

func (_MockARMContract *MockARMContractSession) OwnerUnvoteToCurse(arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error) {
	return _MockARMContract.Contract.OwnerUnvoteToCurse(&_MockARMContract.TransactOpts, arg0)
}

func (_MockARMContract *MockARMContractTransactorSession) OwnerUnvoteToCurse(arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error) {
	return _MockARMContract.Contract.OwnerUnvoteToCurse(&_MockARMContract.TransactOpts, arg0)
}

func (_MockARMContract *MockARMContractTransactor) SetRevert(opts *bind.TransactOpts, err []byte) (*types.Transaction, error) {
	return _MockARMContract.contract.Transact(opts, "setRevert", err)
}

func (_MockARMContract *MockARMContractSession) SetRevert(err []byte) (*types.Transaction, error) {
	return _MockARMContract.Contract.SetRevert(&_MockARMContract.TransactOpts, err)
}

func (_MockARMContract *MockARMContractTransactorSession) SetRevert(err []byte) (*types.Transaction, error) {
	return _MockARMContract.Contract.SetRevert(&_MockARMContract.TransactOpts, err)
}

func (_MockARMContract *MockARMContractTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MockARMContract.contract.Transact(opts, "transferOwnership", to)
}

func (_MockARMContract *MockARMContractSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockARMContract.Contract.TransferOwnership(&_MockARMContract.TransactOpts, to)
}

func (_MockARMContract *MockARMContractTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockARMContract.Contract.TransferOwnership(&_MockARMContract.TransactOpts, to)
}

func (_MockARMContract *MockARMContractTransactor) VoteToCurse(opts *bind.TransactOpts, arg0 [32]byte) (*types.Transaction, error) {
	return _MockARMContract.contract.Transact(opts, "voteToCurse", arg0)
}

func (_MockARMContract *MockARMContractSession) VoteToCurse(arg0 [32]byte) (*types.Transaction, error) {
	return _MockARMContract.Contract.VoteToCurse(&_MockARMContract.TransactOpts, arg0)
}

func (_MockARMContract *MockARMContractTransactorSession) VoteToCurse(arg0 [32]byte) (*types.Transaction, error) {
	return _MockARMContract.Contract.VoteToCurse(&_MockARMContract.TransactOpts, arg0)
}

type MockARMContractOwnershipTransferRequestedIterator struct {
	Event *MockARMContractOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockARMContractOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockARMContractOwnershipTransferRequested)
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
		it.Event = new(MockARMContractOwnershipTransferRequested)
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

func (it *MockARMContractOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MockARMContractOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockARMContractOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockARMContract *MockARMContractFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockARMContractOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockARMContract.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockARMContractOwnershipTransferRequestedIterator{contract: _MockARMContract.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MockARMContract *MockARMContractFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockARMContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockARMContract.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockARMContractOwnershipTransferRequested)
				if err := _MockARMContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MockARMContract *MockARMContractFilterer) ParseOwnershipTransferRequested(log types.Log) (*MockARMContractOwnershipTransferRequested, error) {
	event := new(MockARMContractOwnershipTransferRequested)
	if err := _MockARMContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockARMContractOwnershipTransferredIterator struct {
	Event *MockARMContractOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockARMContractOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockARMContractOwnershipTransferred)
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
		it.Event = new(MockARMContractOwnershipTransferred)
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

func (it *MockARMContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MockARMContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockARMContractOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockARMContract *MockARMContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockARMContractOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockARMContract.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockARMContractOwnershipTransferredIterator{contract: _MockARMContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MockARMContract *MockARMContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockARMContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockARMContract.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockARMContractOwnershipTransferred)
				if err := _MockARMContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MockARMContract *MockARMContractFilterer) ParseOwnershipTransferred(log types.Log) (*MockARMContractOwnershipTransferred, error) {
	event := new(MockARMContractOwnershipTransferred)
	if err := _MockARMContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConfigDetails struct {
	Version     uint32
	BlockNumber uint32
	Config      RMNConfig
}

func (_MockARMContract *MockARMContract) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockARMContract.abi.Events["OwnershipTransferRequested"].ID:
		return _MockARMContract.ParseOwnershipTransferRequested(log)
	case _MockARMContract.abi.Events["OwnershipTransferred"].ID:
		return _MockARMContract.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockARMContractOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MockARMContractOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_MockARMContract *MockARMContract) Address() common.Address {
	return _MockARMContract.address
}

type MockARMContractInterface interface {
	GetConfigDetails(opts *bind.CallOpts) (GetConfigDetails,

		error)

	IsBlessed(opts *bind.CallOpts, arg0 IRMNTaggedRoot) (bool, error)

	IsCursed(opts *bind.CallOpts) (bool, error)

	IsCursed0(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OwnerUnvoteToCurse(opts *bind.TransactOpts, arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error)

	SetRevert(opts *bind.TransactOpts, err []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	VoteToCurse(opts *bind.TransactOpts, arg0 [32]byte) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockARMContractOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockARMContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MockARMContractOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockARMContractOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockARMContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MockARMContractOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
