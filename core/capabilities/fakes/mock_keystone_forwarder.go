// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fakes

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

// Reference imports to suppress errors if they are not otherwise used.
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

// IRouterTransmissionInfo is an auto generated low-level Go binding around an user-defined struct.
type IRouterTransmissionInfo struct {
	TransmissionId  [32]byte
	State           uint8
	Transmitter     common.Address
	InvalidReceiver bool
	Success         bool
	GasLimit        *big.Int
}

// MockKeystoneForwarderMetaData contains all meta data concerning the MockKeystoneForwarder contract.
var MockKeystoneForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"AlreadyAttempted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientGasForRouting\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedForwarder\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"ReportProcessed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"addForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"invalidReceiver\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint80\",\"name\":\"gasLimit\",\"type\":\"uint80\"}],\"internalType\":\"structIRouter.TransmissionInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"removeForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"reportContext\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"validatedReport\",\"type\":\"bytes\"}],\"name\":\"route\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161561009757610097816100b9565b5050306000908152600260205260409020805460ff1916600117905550610162565b336001600160a01b038216036101115760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6111e3806101716000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80635c41d2fe116100815780638da5cb5b1161005b5780638da5cb5b14610329578063abcef55414610347578063f2fde38b1461038057600080fd5b80635c41d2fe1461023557806379ba5097146102485780638864b8641461025057600080fd5b8063272cbd93116100b2578063272cbd9314610163578063354bdd66146101835780634d93172d1461022257600080fd5b806311289565146100d9578063181f5a77146100ee578063233fd52d14610140575b600080fd5b6100ec6100e7366004610d6d565b610393565b005b61012a6040518060400160405280601b81526020017f4d6f636b4b657973746f6e65466f7277617264657220312e302e30000000000081525081565b6040516101379190610e48565b60405180910390f35b61015361014e366004610eb5565b6105af565b6040519015158152602001610137565b610176610171366004610f50565b610751565b6040516101379190610fe4565b610214610191366004610f50565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201526000906056016040516020818303038152906040528051906020012090509392505050565b604051908152602001610137565b6100ec61023036600461108c565b610957565b6100ec61024336600461108c565b6109d3565b6100ec610a52565b61030461025e366004610f50565b6040805160609490941b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660208086019190915260348501939093527fffff000000000000000000000000000000000000000000000000000000000000919091166054840152805160368185030181526056909301815282519282019290922060009081526003909152205473ffffffffffffffffffffffffffffffffffffffff1690565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610137565b60005473ffffffffffffffffffffffffffffffffffffffff16610304565b61015361035536600461108c565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b6100ec61038e36600461108c565b610b54565b606d8510156103ce576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600061041289898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610b6892505050565b6040805160608f901b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602080830191909152603482018690527fffff000000000000000000000000000000000000000000000000000000000000841660548301528251603681840301815260569092019092528051910120929550935060009250309163233fd52d9150338d8d8d602d90606d926104b4939291906110ae565b8f8f606d9080926104c7939291906110ae565b6040518863ffffffff1660e01b81526004016104e99796959493929190611121565b6020604051808303816000875af1158015610508573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061052c9190611182565b9050817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916838b73ffffffffffffffffffffffffffffffffffffffff167f3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b58460405161059b911515815260200190565b60405180910390a450505050505050505050565b600087815260036020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff89161790555a600089815260036020526040808220805469ffffffffffffffffffff949094167601000000000000000000000000000000000000000000000275ffffffffffffffffffffffffffffffffffffffffffff909416939093179092559051819061066e9088908890889088906024016111a4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f805f21320000000000000000000000000000000000000000000000000000000017815281519192506000918291828c5af160009a8b5260036020526040909a2080547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff1675010000000000000000000000000000000000000000008c151502179055509798975050505050505050565b6040805160c0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840183905284519088901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001681830152603481018790527fffff000000000000000000000000000000000000000000000000000000000000861660548201528451603681830301815260568201808752815191840191909120808552600390935285842060d68301909652945473ffffffffffffffffffffffffffffffffffffffff811680875274010000000000000000000000000000000000000000820460ff9081161515607685015275010000000000000000000000000000000000000000008304161515609684015276010000000000000000000000000000000000000000000090910469ffffffffffffffffffff1660b690920191909152929390929091906108af575060006108d7565b8160200151156108c1575060026108d7565b81604001516108d15760036108d4565b60015b90505b6040518060c001604052808481526020018260038111156108fa576108fa610fb5565b8152602001836000015173ffffffffffffffffffffffffffffffffffffffff168152602001836020015115158152602001836040015115158152602001836060015169ffffffffffffffffffff1681525093505050509392505050565b61095f610b83565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b6109db610b83565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff163314610ad8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b5c610b83565b610b6581610c06565b50565b60218101516045820151608b90920151909260c09290921c91565b60005473ffffffffffffffffffffffffffffffffffffffff163314610c04576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610acf565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610c85576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610acf565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803573ffffffffffffffffffffffffffffffffffffffff81168114610d1f57600080fd5b919050565b60008083601f840112610d3657600080fd5b50813567ffffffffffffffff811115610d4e57600080fd5b602083019150836020828501011115610d6657600080fd5b9250929050565b60008060008060008060006080888a031215610d8857600080fd5b610d9188610cfb565b9650602088013567ffffffffffffffff80821115610dae57600080fd5b610dba8b838c01610d24565b909850965060408a0135915080821115610dd357600080fd5b610ddf8b838c01610d24565b909650945060608a0135915080821115610df857600080fd5b818a0191508a601f830112610e0c57600080fd5b813581811115610e1b57600080fd5b8b60208260051b8501011115610e3057600080fd5b60208301945080935050505092959891949750929550565b60006020808352835180602085015260005b81811015610e7657858101830151858201604001528201610e5a565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b600080600080600080600060a0888a031215610ed057600080fd5b87359650610ee060208901610cfb565b9550610eee60408901610cfb565b9450606088013567ffffffffffffffff80821115610f0b57600080fd5b610f178b838c01610d24565b909650945060808a0135915080821115610f3057600080fd5b50610f3d8a828b01610d24565b989b979a50959850939692959293505050565b600080600060608486031215610f6557600080fd5b610f6e84610cfb565b92506020840135915060408401357fffff00000000000000000000000000000000000000000000000000000000000081168114610faa57600080fd5b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b81518152602082015160c082019060048110611029577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8060208401525073ffffffffffffffffffffffffffffffffffffffff604084015116604083015260608301511515606083015260808301511515608083015260a083015161108560a084018269ffffffffffffffffffff169052565b5092915050565b60006020828403121561109e57600080fd5b6110a782610cfb565b9392505050565b600080858511156110be57600080fd5b838611156110cb57600080fd5b5050820193919092039150565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b878152600073ffffffffffffffffffffffffffffffffffffffff808916602084015280881660408401525060a0606083015261116160a0830186886110d8565b82810360808401526111748185876110d8565b9a9950505050505050505050565b60006020828403121561119457600080fd5b815180151581146110a757600080fd5b6040815260006111b86040830186886110d8565b82810360208401526111cb8185876110d8565b97965050505050505056fea164736f6c6343000816000a",
}

// MockKeystoneForwarderABI is the input ABI used to generate the binding from.
// Deprecated: Use MockKeystoneForwarderMetaData.ABI instead.
var MockKeystoneForwarderABI = MockKeystoneForwarderMetaData.ABI

// MockKeystoneForwarderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockKeystoneForwarderMetaData.Bin instead.
var MockKeystoneForwarderBin = MockKeystoneForwarderMetaData.Bin

// DeployMockKeystoneForwarder deploys a new Ethereum contract, binding an instance of MockKeystoneForwarder to it.
func DeployMockKeystoneForwarder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockKeystoneForwarder, error) {
	parsed, err := MockKeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockKeystoneForwarderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockKeystoneForwarder{MockKeystoneForwarderCaller: MockKeystoneForwarderCaller{contract: contract}, MockKeystoneForwarderTransactor: MockKeystoneForwarderTransactor{contract: contract}, MockKeystoneForwarderFilterer: MockKeystoneForwarderFilterer{contract: contract}}, nil
}

// MockKeystoneForwarder is an auto generated Go binding around an Ethereum contract.
type MockKeystoneForwarder struct {
	MockKeystoneForwarderCaller     // Read-only binding to the contract
	MockKeystoneForwarderTransactor // Write-only binding to the contract
	MockKeystoneForwarderFilterer   // Log filterer for contract events
}

// MockKeystoneForwarderCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockKeystoneForwarderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockKeystoneForwarderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockKeystoneForwarderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockKeystoneForwarderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockKeystoneForwarderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockKeystoneForwarderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockKeystoneForwarderSession struct {
	Contract     *MockKeystoneForwarder // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MockKeystoneForwarderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockKeystoneForwarderCallerSession struct {
	Contract *MockKeystoneForwarderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// MockKeystoneForwarderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockKeystoneForwarderTransactorSession struct {
	Contract     *MockKeystoneForwarderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// MockKeystoneForwarderRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockKeystoneForwarderRaw struct {
	Contract *MockKeystoneForwarder // Generic contract binding to access the raw methods on
}

// MockKeystoneForwarderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockKeystoneForwarderCallerRaw struct {
	Contract *MockKeystoneForwarderCaller // Generic read-only contract binding to access the raw methods on
}

// MockKeystoneForwarderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockKeystoneForwarderTransactorRaw struct {
	Contract *MockKeystoneForwarderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockKeystoneForwarder creates a new instance of MockKeystoneForwarder, bound to a specific deployed contract.
func NewMockKeystoneForwarder(address common.Address, backend bind.ContractBackend) (*MockKeystoneForwarder, error) {
	contract, err := bindMockKeystoneForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarder{MockKeystoneForwarderCaller: MockKeystoneForwarderCaller{contract: contract}, MockKeystoneForwarderTransactor: MockKeystoneForwarderTransactor{contract: contract}, MockKeystoneForwarderFilterer: MockKeystoneForwarderFilterer{contract: contract}}, nil
}

// NewMockKeystoneForwarderCaller creates a new read-only instance of MockKeystoneForwarder, bound to a specific deployed contract.
func NewMockKeystoneForwarderCaller(address common.Address, caller bind.ContractCaller) (*MockKeystoneForwarderCaller, error) {
	contract, err := bindMockKeystoneForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderCaller{contract: contract}, nil
}

// NewMockKeystoneForwarderTransactor creates a new write-only instance of MockKeystoneForwarder, bound to a specific deployed contract.
func NewMockKeystoneForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*MockKeystoneForwarderTransactor, error) {
	contract, err := bindMockKeystoneForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderTransactor{contract: contract}, nil
}

// NewMockKeystoneForwarderFilterer creates a new log filterer instance of MockKeystoneForwarder, bound to a specific deployed contract.
func NewMockKeystoneForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*MockKeystoneForwarderFilterer, error) {
	contract, err := bindMockKeystoneForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderFilterer{contract: contract}, nil
}

// bindMockKeystoneForwarder binds a generic wrapper to an already deployed contract.
func bindMockKeystoneForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockKeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockKeystoneForwarder *MockKeystoneForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockKeystoneForwarder.Contract.MockKeystoneForwarderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockKeystoneForwarder *MockKeystoneForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.MockKeystoneForwarderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockKeystoneForwarder *MockKeystoneForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.MockKeystoneForwarderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockKeystoneForwarder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.contract.Transact(opts, method, params...)
}

// GetTransmissionId is a free data retrieval call binding the contract method 0x354bdd66.
//
// Solidity: function getTransmissionId(address receiver, bytes32 workflowExecutionId, bytes2 reportId) pure returns(bytes32)
func (_MockKeystoneForwarder *MockKeystoneForwarderCaller) GetTransmissionId(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	var out []interface{}
	err := _MockKeystoneForwarder.contract.Call(opts, &out, "getTransmissionId", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetTransmissionId is a free data retrieval call binding the contract method 0x354bdd66.
//
// Solidity: function getTransmissionId(address receiver, bytes32 workflowExecutionId, bytes2 reportId) pure returns(bytes32)
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) GetTransmissionId(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _MockKeystoneForwarder.Contract.GetTransmissionId(&_MockKeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

// GetTransmissionId is a free data retrieval call binding the contract method 0x354bdd66.
//
// Solidity: function getTransmissionId(address receiver, bytes32 workflowExecutionId, bytes2 reportId) pure returns(bytes32)
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerSession) GetTransmissionId(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _MockKeystoneForwarder.Contract.GetTransmissionId(&_MockKeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

// GetTransmissionInfo is a free data retrieval call binding the contract method 0x272cbd93.
//
// Solidity: function getTransmissionInfo(address receiver, bytes32 workflowExecutionId, bytes2 reportId) view returns((bytes32,uint8,address,bool,bool,uint80))
func (_MockKeystoneForwarder *MockKeystoneForwarderCaller) GetTransmissionInfo(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error) {
	var out []interface{}
	err := _MockKeystoneForwarder.contract.Call(opts, &out, "getTransmissionInfo", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new(IRouterTransmissionInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IRouterTransmissionInfo)).(*IRouterTransmissionInfo)

	return out0, err

}

// GetTransmissionInfo is a free data retrieval call binding the contract method 0x272cbd93.
//
// Solidity: function getTransmissionInfo(address receiver, bytes32 workflowExecutionId, bytes2 reportId) view returns((bytes32,uint8,address,bool,bool,uint80))
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) GetTransmissionInfo(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error) {
	return _MockKeystoneForwarder.Contract.GetTransmissionInfo(&_MockKeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

// GetTransmissionInfo is a free data retrieval call binding the contract method 0x272cbd93.
//
// Solidity: function getTransmissionInfo(address receiver, bytes32 workflowExecutionId, bytes2 reportId) view returns((bytes32,uint8,address,bool,bool,uint80))
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerSession) GetTransmissionInfo(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error) {
	return _MockKeystoneForwarder.Contract.GetTransmissionInfo(&_MockKeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

// GetTransmitter is a free data retrieval call binding the contract method 0x8864b864.
//
// Solidity: function getTransmitter(address receiver, bytes32 workflowExecutionId, bytes2 reportId) view returns(address)
func (_MockKeystoneForwarder *MockKeystoneForwarderCaller) GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	var out []interface{}
	err := _MockKeystoneForwarder.contract.Call(opts, &out, "getTransmitter", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTransmitter is a free data retrieval call binding the contract method 0x8864b864.
//
// Solidity: function getTransmitter(address receiver, bytes32 workflowExecutionId, bytes2 reportId) view returns(address)
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _MockKeystoneForwarder.Contract.GetTransmitter(&_MockKeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

// GetTransmitter is a free data retrieval call binding the contract method 0x8864b864.
//
// Solidity: function getTransmitter(address receiver, bytes32 workflowExecutionId, bytes2 reportId) view returns(address)
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _MockKeystoneForwarder.Contract.GetTransmitter(&_MockKeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

// IsForwarder is a free data retrieval call binding the contract method 0xabcef554.
//
// Solidity: function isForwarder(address forwarder) view returns(bool)
func (_MockKeystoneForwarder *MockKeystoneForwarderCaller) IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _MockKeystoneForwarder.contract.Call(opts, &out, "isForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsForwarder is a free data retrieval call binding the contract method 0xabcef554.
//
// Solidity: function isForwarder(address forwarder) view returns(bool)
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _MockKeystoneForwarder.Contract.IsForwarder(&_MockKeystoneForwarder.CallOpts, forwarder)
}

// IsForwarder is a free data retrieval call binding the contract method 0xabcef554.
//
// Solidity: function isForwarder(address forwarder) view returns(bool)
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _MockKeystoneForwarder.Contract.IsForwarder(&_MockKeystoneForwarder.CallOpts, forwarder)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MockKeystoneForwarder *MockKeystoneForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockKeystoneForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) Owner() (common.Address, error) {
	return _MockKeystoneForwarder.Contract.Owner(&_MockKeystoneForwarder.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerSession) Owner() (common.Address, error) {
	return _MockKeystoneForwarder.Contract.Owner(&_MockKeystoneForwarder.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_MockKeystoneForwarder *MockKeystoneForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockKeystoneForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) TypeAndVersion() (string, error) {
	return _MockKeystoneForwarder.Contract.TypeAndVersion(&_MockKeystoneForwarder.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_MockKeystoneForwarder *MockKeystoneForwarderCallerSession) TypeAndVersion() (string, error) {
	return _MockKeystoneForwarder.Contract.TypeAndVersion(&_MockKeystoneForwarder.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockKeystoneForwarder.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.AcceptOwnership(&_MockKeystoneForwarder.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.AcceptOwnership(&_MockKeystoneForwarder.TransactOpts)
}

// AddForwarder is a paid mutator transaction binding the contract method 0x5c41d2fe.
//
// Solidity: function addForwarder(address forwarder) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactor) AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.contract.Transact(opts, "addForwarder", forwarder)
}

// AddForwarder is a paid mutator transaction binding the contract method 0x5c41d2fe.
//
// Solidity: function addForwarder(address forwarder) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.AddForwarder(&_MockKeystoneForwarder.TransactOpts, forwarder)
}

// AddForwarder is a paid mutator transaction binding the contract method 0x5c41d2fe.
//
// Solidity: function addForwarder(address forwarder) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.AddForwarder(&_MockKeystoneForwarder.TransactOpts, forwarder)
}

// RemoveForwarder is a paid mutator transaction binding the contract method 0x4d93172d.
//
// Solidity: function removeForwarder(address forwarder) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactor) RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.contract.Transact(opts, "removeForwarder", forwarder)
}

// RemoveForwarder is a paid mutator transaction binding the contract method 0x4d93172d.
//
// Solidity: function removeForwarder(address forwarder) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.RemoveForwarder(&_MockKeystoneForwarder.TransactOpts, forwarder)
}

// RemoveForwarder is a paid mutator transaction binding the contract method 0x4d93172d.
//
// Solidity: function removeForwarder(address forwarder) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.RemoveForwarder(&_MockKeystoneForwarder.TransactOpts, forwarder)
}

// Report is a paid mutator transaction binding the contract method 0x11289565.
//
// Solidity: function report(address receiver, bytes rawReport, bytes reportContext, bytes[] signatures) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactor) Report(opts *bind.TransactOpts, receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _MockKeystoneForwarder.contract.Transact(opts, "report", receiver, rawReport, reportContext, signatures)
}

// Report is a paid mutator transaction binding the contract method 0x11289565.
//
// Solidity: function report(address receiver, bytes rawReport, bytes reportContext, bytes[] signatures) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) Report(receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.Report(&_MockKeystoneForwarder.TransactOpts, receiver, rawReport, reportContext, signatures)
}

// Report is a paid mutator transaction binding the contract method 0x11289565.
//
// Solidity: function report(address receiver, bytes rawReport, bytes reportContext, bytes[] signatures) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorSession) Report(receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.Report(&_MockKeystoneForwarder.TransactOpts, receiver, rawReport, reportContext, signatures)
}

// Route is a paid mutator transaction binding the contract method 0x233fd52d.
//
// Solidity: function route(bytes32 transmissionId, address transmitter, address receiver, bytes metadata, bytes validatedReport) returns(bool)
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactor) Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _MockKeystoneForwarder.contract.Transact(opts, "route", transmissionId, transmitter, receiver, metadata, validatedReport)
}

// Route is a paid mutator transaction binding the contract method 0x233fd52d.
//
// Solidity: function route(bytes32 transmissionId, address transmitter, address receiver, bytes metadata, bytes validatedReport) returns(bool)
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.Route(&_MockKeystoneForwarder.TransactOpts, transmissionId, transmitter, receiver, metadata, validatedReport)
}

// Route is a paid mutator transaction binding the contract method 0x233fd52d.
//
// Solidity: function route(bytes32 transmissionId, address transmitter, address receiver, bytes metadata, bytes validatedReport) returns(bool)
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.Route(&_MockKeystoneForwarder.TransactOpts, transmissionId, transmitter, receiver, metadata, validatedReport)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.TransferOwnership(&_MockKeystoneForwarder.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_MockKeystoneForwarder *MockKeystoneForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockKeystoneForwarder.Contract.TransferOwnership(&_MockKeystoneForwarder.TransactOpts, to)
}

// MockKeystoneForwarderForwarderAddedIterator is returned from FilterForwarderAdded and is used to iterate over the raw logs and unpacked data for ForwarderAdded events raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderForwarderAddedIterator struct {
	Event *MockKeystoneForwarderForwarderAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockKeystoneForwarderForwarderAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockKeystoneForwarderForwarderAdded)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockKeystoneForwarderForwarderAdded)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockKeystoneForwarderForwarderAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockKeystoneForwarderForwarderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockKeystoneForwarderForwarderAdded represents a ForwarderAdded event raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderForwarderAdded struct {
	Forwarder common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterForwarderAdded is a free log retrieval operation binding the contract event 0x0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e7.
//
// Solidity: event ForwarderAdded(address indexed forwarder)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*MockKeystoneForwarderForwarderAddedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.FilterLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderForwarderAddedIterator{contract: _MockKeystoneForwarder.contract, event: "ForwarderAdded", logs: logs, sub: sub}, nil
}

// WatchForwarderAdded is a free log subscription operation binding the contract event 0x0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e7.
//
// Solidity: event ForwarderAdded(address indexed forwarder)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *MockKeystoneForwarderForwarderAdded, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.WatchLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockKeystoneForwarderForwarderAdded)
				if err := _MockKeystoneForwarder.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
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

// ParseForwarderAdded is a log parse operation binding the contract event 0x0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e7.
//
// Solidity: event ForwarderAdded(address indexed forwarder)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) ParseForwarderAdded(log types.Log) (*MockKeystoneForwarderForwarderAdded, error) {
	event := new(MockKeystoneForwarderForwarderAdded)
	if err := _MockKeystoneForwarder.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockKeystoneForwarderForwarderRemovedIterator is returned from FilterForwarderRemoved and is used to iterate over the raw logs and unpacked data for ForwarderRemoved events raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderForwarderRemovedIterator struct {
	Event *MockKeystoneForwarderForwarderRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockKeystoneForwarderForwarderRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockKeystoneForwarderForwarderRemoved)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockKeystoneForwarderForwarderRemoved)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockKeystoneForwarderForwarderRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockKeystoneForwarderForwarderRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockKeystoneForwarderForwarderRemoved represents a ForwarderRemoved event raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderForwarderRemoved struct {
	Forwarder common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterForwarderRemoved is a free log retrieval operation binding the contract event 0xb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d38.
//
// Solidity: event ForwarderRemoved(address indexed forwarder)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*MockKeystoneForwarderForwarderRemovedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.FilterLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderForwarderRemovedIterator{contract: _MockKeystoneForwarder.contract, event: "ForwarderRemoved", logs: logs, sub: sub}, nil
}

// WatchForwarderRemoved is a free log subscription operation binding the contract event 0xb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d38.
//
// Solidity: event ForwarderRemoved(address indexed forwarder)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *MockKeystoneForwarderForwarderRemoved, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.WatchLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockKeystoneForwarderForwarderRemoved)
				if err := _MockKeystoneForwarder.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
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

// ParseForwarderRemoved is a log parse operation binding the contract event 0xb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d38.
//
// Solidity: event ForwarderRemoved(address indexed forwarder)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) ParseForwarderRemoved(log types.Log) (*MockKeystoneForwarderForwarderRemoved, error) {
	event := new(MockKeystoneForwarderForwarderRemoved)
	if err := _MockKeystoneForwarder.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockKeystoneForwarderOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderOwnershipTransferRequestedIterator struct {
	Event *MockKeystoneForwarderOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockKeystoneForwarderOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockKeystoneForwarderOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockKeystoneForwarderOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockKeystoneForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockKeystoneForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockKeystoneForwarderOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockKeystoneForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderOwnershipTransferRequestedIterator{contract: _MockKeystoneForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockKeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockKeystoneForwarderOwnershipTransferRequested)
				if err := _MockKeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*MockKeystoneForwarderOwnershipTransferRequested, error) {
	event := new(MockKeystoneForwarderOwnershipTransferRequested)
	if err := _MockKeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockKeystoneForwarderOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderOwnershipTransferredIterator struct {
	Event *MockKeystoneForwarderOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockKeystoneForwarderOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockKeystoneForwarderOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockKeystoneForwarderOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockKeystoneForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockKeystoneForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockKeystoneForwarderOwnershipTransferred represents a OwnershipTransferred event raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockKeystoneForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderOwnershipTransferredIterator{contract: _MockKeystoneForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockKeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockKeystoneForwarderOwnershipTransferred)
				if err := _MockKeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*MockKeystoneForwarderOwnershipTransferred, error) {
	event := new(MockKeystoneForwarderOwnershipTransferred)
	if err := _MockKeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockKeystoneForwarderReportProcessedIterator is returned from FilterReportProcessed and is used to iterate over the raw logs and unpacked data for ReportProcessed events raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderReportProcessedIterator struct {
	Event *MockKeystoneForwarderReportProcessed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockKeystoneForwarderReportProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockKeystoneForwarderReportProcessed)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockKeystoneForwarderReportProcessed)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockKeystoneForwarderReportProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockKeystoneForwarderReportProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockKeystoneForwarderReportProcessed represents a ReportProcessed event raised by the MockKeystoneForwarder contract.
type MockKeystoneForwarderReportProcessed struct {
	Receiver            common.Address
	WorkflowExecutionId [32]byte
	ReportId            [2]byte
	Result              bool
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterReportProcessed is a free log retrieval operation binding the contract event 0x3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b5.
//
// Solidity: event ReportProcessed(address indexed receiver, bytes32 indexed workflowExecutionId, bytes2 indexed reportId, bool result)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (*MockKeystoneForwarderReportProcessedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.FilterLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportIdRule)
	if err != nil {
		return nil, err
	}
	return &MockKeystoneForwarderReportProcessedIterator{contract: _MockKeystoneForwarder.contract, event: "ReportProcessed", logs: logs, sub: sub}, nil
}

// WatchReportProcessed is a free log subscription operation binding the contract event 0x3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b5.
//
// Solidity: event ReportProcessed(address indexed receiver, bytes32 indexed workflowExecutionId, bytes2 indexed reportId, bool result)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *MockKeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _MockKeystoneForwarder.contract.WatchLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockKeystoneForwarderReportProcessed)
				if err := _MockKeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
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

// ParseReportProcessed is a log parse operation binding the contract event 0x3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b5.
//
// Solidity: event ReportProcessed(address indexed receiver, bytes32 indexed workflowExecutionId, bytes2 indexed reportId, bool result)
func (_MockKeystoneForwarder *MockKeystoneForwarderFilterer) ParseReportProcessed(log types.Log) (*MockKeystoneForwarderReportProcessed, error) {
	event := new(MockKeystoneForwarderReportProcessed)
	if err := _MockKeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
