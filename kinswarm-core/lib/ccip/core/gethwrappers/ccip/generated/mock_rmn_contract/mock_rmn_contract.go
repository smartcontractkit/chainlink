// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_rmn_contract

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

var MockRMNContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"CustomError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"blessVoteAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"curseVoteAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"curseUnvoteAddr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"blessWeight\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"curseWeight\",\"type\":\"uint8\"}],\"internalType\":\"structRMN.Voter[]\",\"name\":\"voters\",\"type\":\"tuple[]\"},{\"internalType\":\"uint16\",\"name\":\"blessWeightThreshold\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"curseWeightThreshold\",\"type\":\"uint16\"}],\"internalType\":\"structRMN.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMN.TaggedRoot\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"isBlessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"curseVoteAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"cursesHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"forceUnvote\",\"type\":\"bool\"}],\"internalType\":\"structRMN.UnvoteToCurseRecord[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"ownerUnvoteToCurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"curseVoteAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"cursesHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"forceUnvote\",\"type\":\"bool\"}],\"internalType\":\"structRMN.UnvoteToCurseRecord[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"ownerUnvoteToCurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"setRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"voteToCurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"voteToCurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610ed7806101576000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063618af128116100815780637a7c27491161005b5780637a7c2749146102b55780638da5cb5b146102c8578063f2fde38b146102f057600080fd5b8063618af1281461020a578063794860871461024357806379ba5097146102ad57600080fd5b8063397796f7116100b2578063397796f7146101ba5780633f42ab73146101c25780634d616771146101d957600080fd5b8063119a3527146100d9578063257174dc1461012b5780632cbc26bb14610192575b600080fd5b6101296100e73660046107fe565b50600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055565b005b6101296101393660046109db565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000016600090815260066020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905550565b6101a56101a0366004610a29565b610303565b60405190151581526020015b60405180910390f35b6101a56103b7565b6101ca610424565b6040516101b193929190610a4b565b6101a56101e7366004610b1e565b5060015474010000000000000000000000000000000000000000900460ff161590565b610129610218366004610b36565b50600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff169055565b610129610251366004610b73565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000016600090815260066020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905550565b610129610565565b6101296102c3366004610b96565b610662565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101b1565b6101296102fe366004610c49565b610672565b60006002805461031290610c64565b1590506103575760026040517f5a4ff67100000000000000000000000000000000000000000000000000000000815260040161034e9190610cb1565b60405180910390fd5b60015474010000000000000000000000000000000000000000900460ff16806103b157507fffffffffffffffffffffffffffffffff00000000000000000000000000000000821660009081526006602052604090205460ff165b92915050565b6000600280546103c690610c64565b1590506104025760026040517f5a4ff67100000000000000000000000000000000000000000000000000000000815260040161034e9190610cb1565b5060015474010000000000000000000000000000000000000000900460ff1690565b6040805160608082018352815260006020820181905291810182905281906005546040805160038054608060208202840181019094526060830181815263ffffffff8087169664010000000090041694929392849284929184919060009085015b828210156105315760008481526020908190206040805160a08101825260038602909201805473ffffffffffffffffffffffffffffffffffffffff90811684526001808301548216858701526002909201549081169284019290925260ff74010000000000000000000000000000000000000000830481166060850152750100000000000000000000000000000000000000000090920490911660808301529083529092019101610485565b505050908252506001919091015461ffff808216602084015262010000909104166040909101529296919550919350915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161034e565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600261066e8282610db0565b5050565b61067a610686565b61068381610709565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610707576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161034e565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610788576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161034e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561081057600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff8111828210171561086957610869610817565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156108b6576108b6610817565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146108e257600080fd5b919050565b600082601f8301126108f857600080fd5b8135602067ffffffffffffffff82111561091457610914610817565b610922818360051b0161086f565b8281526060928302850182019282820191908785111561094157600080fd5b8387015b8581101561099e5781818a03121561095d5760008081fd5b610965610846565b61096e826108be565b81528582013586820152604080830135801515811461098d5760008081fd5b908201528452928401928101610945565b5090979650505050505050565b80357fffffffffffffffffffffffffffffffff00000000000000000000000000000000811681146108e257600080fd5b600080604083850312156109ee57600080fd5b823567ffffffffffffffff811115610a0557600080fd5b610a11858286016108e7565b925050610a20602084016109ab565b90509250929050565b600060208284031215610a3b57600080fd5b610a44826109ab565b9392505050565b63ffffffff84811682528316602080830191909152606060408084018290528451848301839052805160c0860181905260009491820190859060e08801905b80831015610af1578351805173ffffffffffffffffffffffffffffffffffffffff9081168452868201518116878501528782015116878401528781015160ff908116898501526080918201511690830152928401926001929092019160a090910190610a8a565b509288015161ffff908116608089015260409098015190971660a090960195909552979650505050505050565b600060408284031215610b3057600080fd5b50919050565b600060208284031215610b4857600080fd5b813567ffffffffffffffff811115610b5f57600080fd5b610b6b848285016108e7565b949350505050565b60008060408385031215610b8657600080fd5b82359150610a20602084016109ab565b60006020808385031215610ba957600080fd5b823567ffffffffffffffff80821115610bc157600080fd5b818501915085601f830112610bd557600080fd5b813581811115610be757610be7610817565b610c17847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161086f565b91508082528684828501011115610c2d57600080fd5b8084840185840137600090820190930192909252509392505050565b600060208284031215610c5b57600080fd5b610a44826108be565b600181811c90821680610c7857607f821691505b602082108103610b30577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000602080835260008454610cc581610c64565b8060208701526040600180841660008114610ce75760018114610d2157610d51565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00851660408a0152604084151560051b8a01019550610d51565b89600052602060002060005b85811015610d485781548b8201860152908301908801610d2d565b8a016040019650505b509398975050505050505050565b601f821115610dab576000816000526020600020601f850160051c81016020861015610d885750805b601f850160051c820191505b81811015610da757828155600101610d94565b5050505b505050565b815167ffffffffffffffff811115610dca57610dca610817565b610dde81610dd88454610c64565b84610d5f565b602080601f831160018114610e315760008415610dfb5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610da7565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610e7e57888601518255948401946001909101908401610e5f565b5085821015610eba57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000818000a",
}

var MockRMNContractABI = MockRMNContractMetaData.ABI

var MockRMNContractBin = MockRMNContractMetaData.Bin

func DeployMockRMNContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockRMNContract, error) {
	parsed, err := MockRMNContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockRMNContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockRMNContract{address: address, abi: *parsed, MockRMNContractCaller: MockRMNContractCaller{contract: contract}, MockRMNContractTransactor: MockRMNContractTransactor{contract: contract}, MockRMNContractFilterer: MockRMNContractFilterer{contract: contract}}, nil
}

type MockRMNContract struct {
	address common.Address
	abi     abi.ABI
	MockRMNContractCaller
	MockRMNContractTransactor
	MockRMNContractFilterer
}

type MockRMNContractCaller struct {
	contract *bind.BoundContract
}

type MockRMNContractTransactor struct {
	contract *bind.BoundContract
}

type MockRMNContractFilterer struct {
	contract *bind.BoundContract
}

type MockRMNContractSession struct {
	Contract     *MockRMNContract
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockRMNContractCallerSession struct {
	Contract *MockRMNContractCaller
	CallOpts bind.CallOpts
}

type MockRMNContractTransactorSession struct {
	Contract     *MockRMNContractTransactor
	TransactOpts bind.TransactOpts
}

type MockRMNContractRaw struct {
	Contract *MockRMNContract
}

type MockRMNContractCallerRaw struct {
	Contract *MockRMNContractCaller
}

type MockRMNContractTransactorRaw struct {
	Contract *MockRMNContractTransactor
}

func NewMockRMNContract(address common.Address, backend bind.ContractBackend) (*MockRMNContract, error) {
	abi, err := abi.JSON(strings.NewReader(MockRMNContractABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockRMNContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockRMNContract{address: address, abi: abi, MockRMNContractCaller: MockRMNContractCaller{contract: contract}, MockRMNContractTransactor: MockRMNContractTransactor{contract: contract}, MockRMNContractFilterer: MockRMNContractFilterer{contract: contract}}, nil
}

func NewMockRMNContractCaller(address common.Address, caller bind.ContractCaller) (*MockRMNContractCaller, error) {
	contract, err := bindMockRMNContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockRMNContractCaller{contract: contract}, nil
}

func NewMockRMNContractTransactor(address common.Address, transactor bind.ContractTransactor) (*MockRMNContractTransactor, error) {
	contract, err := bindMockRMNContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockRMNContractTransactor{contract: contract}, nil
}

func NewMockRMNContractFilterer(address common.Address, filterer bind.ContractFilterer) (*MockRMNContractFilterer, error) {
	contract, err := bindMockRMNContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockRMNContractFilterer{contract: contract}, nil
}

func bindMockRMNContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockRMNContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockRMNContract *MockRMNContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockRMNContract.Contract.MockRMNContractCaller.contract.Call(opts, result, method, params...)
}

func (_MockRMNContract *MockRMNContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRMNContract.Contract.MockRMNContractTransactor.contract.Transfer(opts)
}

func (_MockRMNContract *MockRMNContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockRMNContract.Contract.MockRMNContractTransactor.contract.Transact(opts, method, params...)
}

func (_MockRMNContract *MockRMNContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockRMNContract.Contract.contract.Call(opts, result, method, params...)
}

func (_MockRMNContract *MockRMNContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRMNContract.Contract.contract.Transfer(opts)
}

func (_MockRMNContract *MockRMNContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockRMNContract.Contract.contract.Transact(opts, method, params...)
}

func (_MockRMNContract *MockRMNContractCaller) GetConfigDetails(opts *bind.CallOpts) (GetConfigDetails,

	error) {
	var out []interface{}
	err := _MockRMNContract.contract.Call(opts, &out, "getConfigDetails")

	outstruct := new(GetConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Config = *abi.ConvertType(out[2], new(RMNConfig)).(*RMNConfig)

	return *outstruct, err

}

func (_MockRMNContract *MockRMNContractSession) GetConfigDetails() (GetConfigDetails,

	error) {
	return _MockRMNContract.Contract.GetConfigDetails(&_MockRMNContract.CallOpts)
}

func (_MockRMNContract *MockRMNContractCallerSession) GetConfigDetails() (GetConfigDetails,

	error) {
	return _MockRMNContract.Contract.GetConfigDetails(&_MockRMNContract.CallOpts)
}

func (_MockRMNContract *MockRMNContractCaller) IsBlessed(opts *bind.CallOpts, arg0 IRMNTaggedRoot) (bool, error) {
	var out []interface{}
	err := _MockRMNContract.contract.Call(opts, &out, "isBlessed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockRMNContract *MockRMNContractSession) IsBlessed(arg0 IRMNTaggedRoot) (bool, error) {
	return _MockRMNContract.Contract.IsBlessed(&_MockRMNContract.CallOpts, arg0)
}

func (_MockRMNContract *MockRMNContractCallerSession) IsBlessed(arg0 IRMNTaggedRoot) (bool, error) {
	return _MockRMNContract.Contract.IsBlessed(&_MockRMNContract.CallOpts, arg0)
}

func (_MockRMNContract *MockRMNContractCaller) IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error) {
	var out []interface{}
	err := _MockRMNContract.contract.Call(opts, &out, "isCursed", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockRMNContract *MockRMNContractSession) IsCursed(subject [16]byte) (bool, error) {
	return _MockRMNContract.Contract.IsCursed(&_MockRMNContract.CallOpts, subject)
}

func (_MockRMNContract *MockRMNContractCallerSession) IsCursed(subject [16]byte) (bool, error) {
	return _MockRMNContract.Contract.IsCursed(&_MockRMNContract.CallOpts, subject)
}

func (_MockRMNContract *MockRMNContractCaller) IsCursed0(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MockRMNContract.contract.Call(opts, &out, "isCursed0")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockRMNContract *MockRMNContractSession) IsCursed0() (bool, error) {
	return _MockRMNContract.Contract.IsCursed0(&_MockRMNContract.CallOpts)
}

func (_MockRMNContract *MockRMNContractCallerSession) IsCursed0() (bool, error) {
	return _MockRMNContract.Contract.IsCursed0(&_MockRMNContract.CallOpts)
}

func (_MockRMNContract *MockRMNContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockRMNContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockRMNContract *MockRMNContractSession) Owner() (common.Address, error) {
	return _MockRMNContract.Contract.Owner(&_MockRMNContract.CallOpts)
}

func (_MockRMNContract *MockRMNContractCallerSession) Owner() (common.Address, error) {
	return _MockRMNContract.Contract.Owner(&_MockRMNContract.CallOpts)
}

func (_MockRMNContract *MockRMNContractTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "acceptOwnership")
}

func (_MockRMNContract *MockRMNContractSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockRMNContract.Contract.AcceptOwnership(&_MockRMNContract.TransactOpts)
}

func (_MockRMNContract *MockRMNContractTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockRMNContract.Contract.AcceptOwnership(&_MockRMNContract.TransactOpts)
}

func (_MockRMNContract *MockRMNContractTransactor) OwnerUnvoteToCurse(opts *bind.TransactOpts, arg0 []RMNUnvoteToCurseRecord, subject [16]byte) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "ownerUnvoteToCurse", arg0, subject)
}

func (_MockRMNContract *MockRMNContractSession) OwnerUnvoteToCurse(arg0 []RMNUnvoteToCurseRecord, subject [16]byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.OwnerUnvoteToCurse(&_MockRMNContract.TransactOpts, arg0, subject)
}

func (_MockRMNContract *MockRMNContractTransactorSession) OwnerUnvoteToCurse(arg0 []RMNUnvoteToCurseRecord, subject [16]byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.OwnerUnvoteToCurse(&_MockRMNContract.TransactOpts, arg0, subject)
}

func (_MockRMNContract *MockRMNContractTransactor) OwnerUnvoteToCurse0(opts *bind.TransactOpts, arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "ownerUnvoteToCurse0", arg0)
}

func (_MockRMNContract *MockRMNContractSession) OwnerUnvoteToCurse0(arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error) {
	return _MockRMNContract.Contract.OwnerUnvoteToCurse0(&_MockRMNContract.TransactOpts, arg0)
}

func (_MockRMNContract *MockRMNContractTransactorSession) OwnerUnvoteToCurse0(arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error) {
	return _MockRMNContract.Contract.OwnerUnvoteToCurse0(&_MockRMNContract.TransactOpts, arg0)
}

func (_MockRMNContract *MockRMNContractTransactor) SetRevert(opts *bind.TransactOpts, err []byte) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "setRevert", err)
}

func (_MockRMNContract *MockRMNContractSession) SetRevert(err []byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.SetRevert(&_MockRMNContract.TransactOpts, err)
}

func (_MockRMNContract *MockRMNContractTransactorSession) SetRevert(err []byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.SetRevert(&_MockRMNContract.TransactOpts, err)
}

func (_MockRMNContract *MockRMNContractTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "transferOwnership", to)
}

func (_MockRMNContract *MockRMNContractSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockRMNContract.Contract.TransferOwnership(&_MockRMNContract.TransactOpts, to)
}

func (_MockRMNContract *MockRMNContractTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockRMNContract.Contract.TransferOwnership(&_MockRMNContract.TransactOpts, to)
}

func (_MockRMNContract *MockRMNContractTransactor) VoteToCurse(opts *bind.TransactOpts, arg0 [32]byte) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "voteToCurse", arg0)
}

func (_MockRMNContract *MockRMNContractSession) VoteToCurse(arg0 [32]byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.VoteToCurse(&_MockRMNContract.TransactOpts, arg0)
}

func (_MockRMNContract *MockRMNContractTransactorSession) VoteToCurse(arg0 [32]byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.VoteToCurse(&_MockRMNContract.TransactOpts, arg0)
}

func (_MockRMNContract *MockRMNContractTransactor) VoteToCurse0(opts *bind.TransactOpts, arg0 [32]byte, subject [16]byte) (*types.Transaction, error) {
	return _MockRMNContract.contract.Transact(opts, "voteToCurse0", arg0, subject)
}

func (_MockRMNContract *MockRMNContractSession) VoteToCurse0(arg0 [32]byte, subject [16]byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.VoteToCurse0(&_MockRMNContract.TransactOpts, arg0, subject)
}

func (_MockRMNContract *MockRMNContractTransactorSession) VoteToCurse0(arg0 [32]byte, subject [16]byte) (*types.Transaction, error) {
	return _MockRMNContract.Contract.VoteToCurse0(&_MockRMNContract.TransactOpts, arg0, subject)
}

type MockRMNContractOwnershipTransferRequestedIterator struct {
	Event *MockRMNContractOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockRMNContractOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockRMNContractOwnershipTransferRequested)
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
		it.Event = new(MockRMNContractOwnershipTransferRequested)
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

func (it *MockRMNContractOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MockRMNContractOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockRMNContractOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockRMNContract *MockRMNContractFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockRMNContractOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockRMNContract.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockRMNContractOwnershipTransferRequestedIterator{contract: _MockRMNContract.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MockRMNContract *MockRMNContractFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockRMNContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockRMNContract.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockRMNContractOwnershipTransferRequested)
				if err := _MockRMNContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MockRMNContract *MockRMNContractFilterer) ParseOwnershipTransferRequested(log types.Log) (*MockRMNContractOwnershipTransferRequested, error) {
	event := new(MockRMNContractOwnershipTransferRequested)
	if err := _MockRMNContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockRMNContractOwnershipTransferredIterator struct {
	Event *MockRMNContractOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockRMNContractOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockRMNContractOwnershipTransferred)
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
		it.Event = new(MockRMNContractOwnershipTransferred)
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

func (it *MockRMNContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MockRMNContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockRMNContractOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockRMNContract *MockRMNContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockRMNContractOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockRMNContract.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockRMNContractOwnershipTransferredIterator{contract: _MockRMNContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MockRMNContract *MockRMNContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockRMNContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockRMNContract.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockRMNContractOwnershipTransferred)
				if err := _MockRMNContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MockRMNContract *MockRMNContractFilterer) ParseOwnershipTransferred(log types.Log) (*MockRMNContractOwnershipTransferred, error) {
	event := new(MockRMNContractOwnershipTransferred)
	if err := _MockRMNContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MockRMNContract *MockRMNContract) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockRMNContract.abi.Events["OwnershipTransferRequested"].ID:
		return _MockRMNContract.ParseOwnershipTransferRequested(log)
	case _MockRMNContract.abi.Events["OwnershipTransferred"].ID:
		return _MockRMNContract.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockRMNContractOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MockRMNContractOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_MockRMNContract *MockRMNContract) Address() common.Address {
	return _MockRMNContract.address
}

type MockRMNContractInterface interface {
	GetConfigDetails(opts *bind.CallOpts) (GetConfigDetails,

		error)

	IsBlessed(opts *bind.CallOpts, arg0 IRMNTaggedRoot) (bool, error)

	IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error)

	IsCursed0(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OwnerUnvoteToCurse(opts *bind.TransactOpts, arg0 []RMNUnvoteToCurseRecord, subject [16]byte) (*types.Transaction, error)

	OwnerUnvoteToCurse0(opts *bind.TransactOpts, arg0 []RMNUnvoteToCurseRecord) (*types.Transaction, error)

	SetRevert(opts *bind.TransactOpts, err []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	VoteToCurse(opts *bind.TransactOpts, arg0 [32]byte) (*types.Transaction, error)

	VoteToCurse0(opts *bind.TransactOpts, arg0 [32]byte, subject [16]byte) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockRMNContractOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockRMNContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MockRMNContractOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockRMNContractOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockRMNContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MockRMNContractOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
