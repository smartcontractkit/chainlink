// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier_proxy

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
)

// VerifierProxyMetaData contains all meta data concerning the VerifierProxy contract.
var VerifierProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"ConfigDigestAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierInvalid\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"VerifierNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAccessController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAccessController\",\"type\":\"address\"}],\"name\":\"AccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierUnset\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"initializeVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"name\":\"setAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"currentConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"unsetVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"verifierResponse\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610da0380380610da083398101604081905261002f91610188565b33806000816100855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b5576100b5816100de565b5050600380546001600160a01b0319166001600160a01b039390931692909217909155506101b8565b6001600160a01b0381163314156101375760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561019a57600080fd5b81516001600160a01b03811681146101b157600080fd5b9392505050565b610bd9806101c76000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806379ba50971161007157806379ba50971461014b5780638da5cb5b146101535780638e760afe14610164578063eeb7b24814610177578063f08391d8146101a0578063f2fde38b146101b357600080fd5b806316d6b5f6146100ae578063181f5a77146100d85780632cc9947714610110578063674145f6146101255780636e91409414610138575b600080fd5b6003546001600160a01b03165b6040516001600160a01b0390911681526020015b60405180910390f35b604080518082019091526013815272566572696669657250726f787920302e302e3160681b60208201525b6040516100cf9190610b16565b61012361011e366004610926565b6101c6565b005b6101236101333660046108f6565b6102b4565b6101236101463660046108dd565b610439565b6101236104e6565b6000546001600160a01b03166100bb565b610103610172366004610948565b610590565b6100bb6101853660046108dd565b6000908152600260205260409020546001600160a01b031690565b6101236101ae366004610897565b610722565b6101236101c1366004610897565b610784565b60008181526002602052604090205481906001600160a01b0316801561021657604051631bae8ff360e11b8152600481018390526001600160a01b03821660248201526044015b60405180910390fd5b6000848152600260205260409020546001600160a01b0316331461024d57604051631decfebb60e31b815260040160405180910390fd5b6000838152600260209081526040918290208054336001600160a01b0319909116811790915582518781529182018690528183015290517fbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf9181900360600190a150505050565b6102bc610798565b806001600160a01b0381166102e45760405163d92e233d60e01b815260040160405180910390fd5b6040516301ffc9a760e01b8152633d3ac1b560e01b60048201526001600160a01b038216906301ffc9a79060240160206040518083038186803b15801561032a57600080fd5b505afa15801561033e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061036291906108bb565b61037f57604051633ad8293d60e11b815260040160405180910390fd5b60008381526002602052604090205483906001600160a01b031680156103ca57604051631bae8ff360e11b8152600481018390526001600160a01b038216602482015260440161020d565b600085815260026020908152604080832080546001600160a01b0319166001600160a01b03891690811790915581519384529183018890528201527fbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf9060600160405180910390a15050505050565b610441610798565b6000818152600260205260409020546001600160a01b03168061047a5760405163b151802b60e01b81526004810183905260240161020d565b6000828152600260205260409081902080546001600160a01b0319169055517f11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c906104da90849084909182526001600160a01b0316602082015260400190565b60405180910390a15050565b6001546001600160a01b031633146105395760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015260640161020d565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6003546060906001600160a01b0316801580159061062c5750604051630d629b5f60e31b81526001600160a01b03821690636b14daf8906105da9033906000903690600401610abc565b60206040518083038186803b1580156105f257600080fd5b505afa158015610606573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061062a91906108bb565b155b1561064a57604051631decfebb60e31b815260040160405180910390fd5b60006106568486610b29565b6000818152600260205260409020549091506001600160a01b0316806106925760405163b151802b60e01b81526004810183905260240161020d565b604051633d3ac1b560e01b81526001600160a01b03821690633d3ac1b5906106c290899089903390600401610aea565b600060405180830381600087803b1580156106dc57600080fd5b505af11580156106f0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261071891908101906109ba565b9695505050505050565b61072a610798565b600380546001600160a01b038381166001600160a01b031983168117909355604080519190921680825260208201939093527f953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b691016104da565b61078c610798565b610795816107ed565b50565b6000546001600160a01b031633146107eb5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015260640161020d565b565b6001600160a01b0381163314156108465760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161020d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156108a957600080fd5b81356108b481610b8e565b9392505050565b6000602082840312156108cd57600080fd5b815180151581146108b457600080fd5b6000602082840312156108ef57600080fd5b5035919050565b6000806040838503121561090957600080fd5b82359150602083013561091b81610b8e565b809150509250929050565b6000806040838503121561093957600080fd5b50508035926020909101359150565b6000806020838503121561095b57600080fd5b823567ffffffffffffffff8082111561097357600080fd5b818501915085601f83011261098757600080fd5b81358181111561099657600080fd5b8660208285010111156109a857600080fd5b60209290920196919550909350505050565b6000602082840312156109cc57600080fd5b815167ffffffffffffffff808211156109e457600080fd5b818401915084601f8301126109f857600080fd5b815181811115610a0a57610a0a610b78565b604051601f8201601f19908116603f01168101908382118183101715610a3257610a32610b78565b81604052828152876020848701011115610a4b57600080fd5b610a5c836020830160208801610b48565b979650505050505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60008151808452610aa8816020860160208601610b48565b601f01601f19169290920160200192915050565b6001600160a01b0384168152604060208201819052600090610ae19083018486610a67565b95945050505050565b604081526000610afe604083018587610a67565b905060018060a01b0383166020830152949350505050565b6020815260006108b46020830184610a90565b80356020831015610b4257600019602084900360031b1b165b92915050565b60005b83811015610b63578181015183820152602001610b4b565b83811115610b72576000848401525b50505050565b634e487b7160e01b600052604160045260246000fd5b6001600160a01b038116811461079557600080fdfea26469706673582212206399aad0de5ba41dcbbbeb20aab5d31ce982d3f893aa208aa0b7d76be001153c64736f6c63430008060033",
}

// VerifierProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierProxyMetaData.ABI instead.
var VerifierProxyABI = VerifierProxyMetaData.ABI

// VerifierProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierProxyMetaData.Bin instead.
var VerifierProxyBin = VerifierProxyMetaData.Bin

// DeployVerifierProxy deploys a new Ethereum contract, binding an instance of VerifierProxy to it.
func DeployVerifierProxy(auth *bind.TransactOpts, backend bind.ContractBackend, accessController common.Address) (common.Address, *types.Transaction, *VerifierProxy, error) {
	parsed, err := VerifierProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierProxyBin), backend, accessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifierProxy{VerifierProxyCaller: VerifierProxyCaller{contract: contract}, VerifierProxyTransactor: VerifierProxyTransactor{contract: contract}, VerifierProxyFilterer: VerifierProxyFilterer{contract: contract}}, nil
}

// VerifierProxy is an auto generated Go binding around an Ethereum contract.
type VerifierProxy struct {
	VerifierProxyCaller     // Read-only binding to the contract
	VerifierProxyTransactor // Write-only binding to the contract
	VerifierProxyFilterer   // Log filterer for contract events
}

// VerifierProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierProxySession struct {
	Contract     *VerifierProxy    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierProxyCallerSession struct {
	Contract *VerifierProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VerifierProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierProxyTransactorSession struct {
	Contract     *VerifierProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VerifierProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierProxyRaw struct {
	Contract *VerifierProxy // Generic contract binding to access the raw methods on
}

// VerifierProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierProxyCallerRaw struct {
	Contract *VerifierProxyCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierProxyTransactorRaw struct {
	Contract *VerifierProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifierProxy creates a new instance of VerifierProxy, bound to a specific deployed contract.
func NewVerifierProxy(address common.Address, backend bind.ContractBackend) (*VerifierProxy, error) {
	contract, err := bindVerifierProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifierProxy{VerifierProxyCaller: VerifierProxyCaller{contract: contract}, VerifierProxyTransactor: VerifierProxyTransactor{contract: contract}, VerifierProxyFilterer: VerifierProxyFilterer{contract: contract}}, nil
}

// NewVerifierProxyCaller creates a new read-only instance of VerifierProxy, bound to a specific deployed contract.
func NewVerifierProxyCaller(address common.Address, caller bind.ContractCaller) (*VerifierProxyCaller, error) {
	contract, err := bindVerifierProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyCaller{contract: contract}, nil
}

// NewVerifierProxyTransactor creates a new write-only instance of VerifierProxy, bound to a specific deployed contract.
func NewVerifierProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierProxyTransactor, error) {
	contract, err := bindVerifierProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyTransactor{contract: contract}, nil
}

// NewVerifierProxyFilterer creates a new log filterer instance of VerifierProxy, bound to a specific deployed contract.
func NewVerifierProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierProxyFilterer, error) {
	contract, err := bindVerifierProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyFilterer{contract: contract}, nil
}

// bindVerifierProxy binds a generic wrapper to an already deployed contract.
func bindVerifierProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VerifierProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VerifierProxy *VerifierProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifierProxy.Contract.VerifierProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VerifierProxy *VerifierProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierProxy.Contract.VerifierProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VerifierProxy *VerifierProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifierProxy.Contract.VerifierProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VerifierProxy *VerifierProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifierProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VerifierProxy *VerifierProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VerifierProxy *VerifierProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifierProxy.Contract.contract.Transact(opts, method, params...)
}

// GetAccessController is a free data retrieval call binding the contract method 0x16d6b5f6.
//
// Solidity: function getAccessController() view returns(address accessController)
func (_VerifierProxy *VerifierProxyCaller) GetAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "getAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAccessController is a free data retrieval call binding the contract method 0x16d6b5f6.
//
// Solidity: function getAccessController() view returns(address accessController)
func (_VerifierProxy *VerifierProxySession) GetAccessController() (common.Address, error) {
	return _VerifierProxy.Contract.GetAccessController(&_VerifierProxy.CallOpts)
}

// GetAccessController is a free data retrieval call binding the contract method 0x16d6b5f6.
//
// Solidity: function getAccessController() view returns(address accessController)
func (_VerifierProxy *VerifierProxyCallerSession) GetAccessController() (common.Address, error) {
	return _VerifierProxy.Contract.GetAccessController(&_VerifierProxy.CallOpts)
}

// GetVerifier is a free data retrieval call binding the contract method 0xeeb7b248.
//
// Solidity: function getVerifier(bytes32 configDigest) view returns(address)
func (_VerifierProxy *VerifierProxyCaller) GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "getVerifier", configDigest)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVerifier is a free data retrieval call binding the contract method 0xeeb7b248.
//
// Solidity: function getVerifier(bytes32 configDigest) view returns(address)
func (_VerifierProxy *VerifierProxySession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _VerifierProxy.Contract.GetVerifier(&_VerifierProxy.CallOpts, configDigest)
}

// GetVerifier is a free data retrieval call binding the contract method 0xeeb7b248.
//
// Solidity: function getVerifier(bytes32 configDigest) view returns(address)
func (_VerifierProxy *VerifierProxyCallerSession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _VerifierProxy.Contract.GetVerifier(&_VerifierProxy.CallOpts, configDigest)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VerifierProxy *VerifierProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VerifierProxy *VerifierProxySession) Owner() (common.Address, error) {
	return _VerifierProxy.Contract.Owner(&_VerifierProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VerifierProxy *VerifierProxyCallerSession) Owner() (common.Address, error) {
	return _VerifierProxy.Contract.Owner(&_VerifierProxy.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_VerifierProxy *VerifierProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_VerifierProxy *VerifierProxySession) TypeAndVersion() (string, error) {
	return _VerifierProxy.Contract.TypeAndVersion(&_VerifierProxy.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_VerifierProxy *VerifierProxyCallerSession) TypeAndVersion() (string, error) {
	return _VerifierProxy.Contract.TypeAndVersion(&_VerifierProxy.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VerifierProxy *VerifierProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VerifierProxy *VerifierProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifierProxy.Contract.AcceptOwnership(&_VerifierProxy.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VerifierProxy *VerifierProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifierProxy.Contract.AcceptOwnership(&_VerifierProxy.TransactOpts)
}

// InitializeVerifier is a paid mutator transaction binding the contract method 0x674145f6.
//
// Solidity: function initializeVerifier(bytes32 configDigest, address verifierAddress) returns()
func (_VerifierProxy *VerifierProxyTransactor) InitializeVerifier(opts *bind.TransactOpts, configDigest [32]byte, verifierAddress common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "initializeVerifier", configDigest, verifierAddress)
}

// InitializeVerifier is a paid mutator transaction binding the contract method 0x674145f6.
//
// Solidity: function initializeVerifier(bytes32 configDigest, address verifierAddress) returns()
func (_VerifierProxy *VerifierProxySession) InitializeVerifier(configDigest [32]byte, verifierAddress common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.InitializeVerifier(&_VerifierProxy.TransactOpts, configDigest, verifierAddress)
}

// InitializeVerifier is a paid mutator transaction binding the contract method 0x674145f6.
//
// Solidity: function initializeVerifier(bytes32 configDigest, address verifierAddress) returns()
func (_VerifierProxy *VerifierProxyTransactorSession) InitializeVerifier(configDigest [32]byte, verifierAddress common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.InitializeVerifier(&_VerifierProxy.TransactOpts, configDigest, verifierAddress)
}

// SetAccessController is a paid mutator transaction binding the contract method 0xf08391d8.
//
// Solidity: function setAccessController(address accessController) returns()
func (_VerifierProxy *VerifierProxyTransactor) SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "setAccessController", accessController)
}

// SetAccessController is a paid mutator transaction binding the contract method 0xf08391d8.
//
// Solidity: function setAccessController(address accessController) returns()
func (_VerifierProxy *VerifierProxySession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetAccessController(&_VerifierProxy.TransactOpts, accessController)
}

// SetAccessController is a paid mutator transaction binding the contract method 0xf08391d8.
//
// Solidity: function setAccessController(address accessController) returns()
func (_VerifierProxy *VerifierProxyTransactorSession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetAccessController(&_VerifierProxy.TransactOpts, accessController)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x2cc99477.
//
// Solidity: function setVerifier(bytes32 currentConfigDigest, bytes32 newConfigDigest) returns()
func (_VerifierProxy *VerifierProxyTransactor) SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "setVerifier", currentConfigDigest, newConfigDigest)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x2cc99477.
//
// Solidity: function setVerifier(bytes32 currentConfigDigest, bytes32 newConfigDigest) returns()
func (_VerifierProxy *VerifierProxySession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetVerifier(&_VerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x2cc99477.
//
// Solidity: function setVerifier(bytes32 currentConfigDigest, bytes32 newConfigDigest) returns()
func (_VerifierProxy *VerifierProxyTransactorSession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetVerifier(&_VerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VerifierProxy *VerifierProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VerifierProxy *VerifierProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.TransferOwnership(&_VerifierProxy.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VerifierProxy *VerifierProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.TransferOwnership(&_VerifierProxy.TransactOpts, to)
}

// UnsetVerifier is a paid mutator transaction binding the contract method 0x6e914094.
//
// Solidity: function unsetVerifier(bytes32 configDigest) returns()
func (_VerifierProxy *VerifierProxyTransactor) UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "unsetVerifier", configDigest)
}

// UnsetVerifier is a paid mutator transaction binding the contract method 0x6e914094.
//
// Solidity: function unsetVerifier(bytes32 configDigest) returns()
func (_VerifierProxy *VerifierProxySession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.UnsetVerifier(&_VerifierProxy.TransactOpts, configDigest)
}

// UnsetVerifier is a paid mutator transaction binding the contract method 0x6e914094.
//
// Solidity: function unsetVerifier(bytes32 configDigest) returns()
func (_VerifierProxy *VerifierProxyTransactorSession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.UnsetVerifier(&_VerifierProxy.TransactOpts, configDigest)
}

// Verify is a paid mutator transaction binding the contract method 0x8e760afe.
//
// Solidity: function verify(bytes signedReport) returns(bytes verifierResponse)
func (_VerifierProxy *VerifierProxyTransactor) Verify(opts *bind.TransactOpts, signedReport []byte) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "verify", signedReport)
}

// Verify is a paid mutator transaction binding the contract method 0x8e760afe.
//
// Solidity: function verify(bytes signedReport) returns(bytes verifierResponse)
func (_VerifierProxy *VerifierProxySession) Verify(signedReport []byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.Verify(&_VerifierProxy.TransactOpts, signedReport)
}

// Verify is a paid mutator transaction binding the contract method 0x8e760afe.
//
// Solidity: function verify(bytes signedReport) returns(bytes verifierResponse)
func (_VerifierProxy *VerifierProxyTransactorSession) Verify(signedReport []byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.Verify(&_VerifierProxy.TransactOpts, signedReport)
}

// VerifierProxyAccessControllerSetIterator is returned from FilterAccessControllerSet and is used to iterate over the raw logs and unpacked data for AccessControllerSet events raised by the VerifierProxy contract.
type VerifierProxyAccessControllerSetIterator struct {
	Event *VerifierProxyAccessControllerSet // Event containing the contract specifics and raw log

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
func (it *VerifierProxyAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyAccessControllerSet)
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
		it.Event = new(VerifierProxyAccessControllerSet)
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
func (it *VerifierProxyAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierProxyAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierProxyAccessControllerSet represents a AccessControllerSet event raised by the VerifierProxy contract.
type VerifierProxyAccessControllerSet struct {
	OldAccessController common.Address
	NewAccessController common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterAccessControllerSet is a free log retrieval operation binding the contract event 0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6.
//
// Solidity: event AccessControllerSet(address oldAccessController, address newAccessController)
func (_VerifierProxy *VerifierProxyFilterer) FilterAccessControllerSet(opts *bind.FilterOpts) (*VerifierProxyAccessControllerSetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyAccessControllerSetIterator{contract: _VerifierProxy.contract, event: "AccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchAccessControllerSet is a free log subscription operation binding the contract event 0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6.
//
// Solidity: event AccessControllerSet(address oldAccessController, address newAccessController)
func (_VerifierProxy *VerifierProxyFilterer) WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierProxyAccessControllerSet)
				if err := _VerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
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

// ParseAccessControllerSet is a log parse operation binding the contract event 0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6.
//
// Solidity: event AccessControllerSet(address oldAccessController, address newAccessController)
func (_VerifierProxy *VerifierProxyFilterer) ParseAccessControllerSet(log types.Log) (*VerifierProxyAccessControllerSet, error) {
	event := new(VerifierProxyAccessControllerSet)
	if err := _VerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierProxyOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the VerifierProxy contract.
type VerifierProxyOwnershipTransferRequestedIterator struct {
	Event *VerifierProxyOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *VerifierProxyOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyOwnershipTransferRequested)
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
		it.Event = new(VerifierProxyOwnershipTransferRequested)
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
func (it *VerifierProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierProxyOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the VerifierProxy contract.
type VerifierProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VerifierProxy *VerifierProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyOwnershipTransferRequestedIterator{contract: _VerifierProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VerifierProxy *VerifierProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierProxyOwnershipTransferRequested)
				if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_VerifierProxy *VerifierProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifierProxyOwnershipTransferRequested, error) {
	event := new(VerifierProxyOwnershipTransferRequested)
	if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierProxyOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the VerifierProxy contract.
type VerifierProxyOwnershipTransferredIterator struct {
	Event *VerifierProxyOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VerifierProxyOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyOwnershipTransferred)
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
		it.Event = new(VerifierProxyOwnershipTransferred)
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
func (it *VerifierProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierProxyOwnershipTransferred represents a OwnershipTransferred event raised by the VerifierProxy contract.
type VerifierProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VerifierProxy *VerifierProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyOwnershipTransferredIterator{contract: _VerifierProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VerifierProxy *VerifierProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierProxyOwnershipTransferred)
				if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_VerifierProxy *VerifierProxyFilterer) ParseOwnershipTransferred(log types.Log) (*VerifierProxyOwnershipTransferred, error) {
	event := new(VerifierProxyOwnershipTransferred)
	if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierProxyVerifierSetIterator is returned from FilterVerifierSet and is used to iterate over the raw logs and unpacked data for VerifierSet events raised by the VerifierProxy contract.
type VerifierProxyVerifierSetIterator struct {
	Event *VerifierProxyVerifierSet // Event containing the contract specifics and raw log

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
func (it *VerifierProxyVerifierSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyVerifierSet)
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
		it.Event = new(VerifierProxyVerifierSet)
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
func (it *VerifierProxyVerifierSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierProxyVerifierSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierProxyVerifierSet represents a VerifierSet event raised by the VerifierProxy contract.
type VerifierProxyVerifierSet struct {
	OldConfigDigest [32]byte
	NewConfigDigest [32]byte
	VerifierAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterVerifierSet is a free log retrieval operation binding the contract event 0xbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf.
//
// Solidity: event VerifierSet(bytes32 oldConfigDigest, bytes32 newConfigDigest, address verifierAddress)
func (_VerifierProxy *VerifierProxyFilterer) FilterVerifierSet(opts *bind.FilterOpts) (*VerifierProxyVerifierSetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyVerifierSetIterator{contract: _VerifierProxy.contract, event: "VerifierSet", logs: logs, sub: sub}, nil
}

// WatchVerifierSet is a free log subscription operation binding the contract event 0xbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf.
//
// Solidity: event VerifierSet(bytes32 oldConfigDigest, bytes32 newConfigDigest, address verifierAddress)
func (_VerifierProxy *VerifierProxyFilterer) WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierSet) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierProxyVerifierSet)
				if err := _VerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
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

// ParseVerifierSet is a log parse operation binding the contract event 0xbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf.
//
// Solidity: event VerifierSet(bytes32 oldConfigDigest, bytes32 newConfigDigest, address verifierAddress)
func (_VerifierProxy *VerifierProxyFilterer) ParseVerifierSet(log types.Log) (*VerifierProxyVerifierSet, error) {
	event := new(VerifierProxyVerifierSet)
	if err := _VerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierProxyVerifierUnsetIterator is returned from FilterVerifierUnset and is used to iterate over the raw logs and unpacked data for VerifierUnset events raised by the VerifierProxy contract.
type VerifierProxyVerifierUnsetIterator struct {
	Event *VerifierProxyVerifierUnset // Event containing the contract specifics and raw log

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
func (it *VerifierProxyVerifierUnsetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyVerifierUnset)
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
		it.Event = new(VerifierProxyVerifierUnset)
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
func (it *VerifierProxyVerifierUnsetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierProxyVerifierUnsetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierProxyVerifierUnset represents a VerifierUnset event raised by the VerifierProxy contract.
type VerifierProxyVerifierUnset struct {
	ConfigDigest    [32]byte
	VerifierAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterVerifierUnset is a free log retrieval operation binding the contract event 0x11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c.
//
// Solidity: event VerifierUnset(bytes32 configDigest, address verifierAddress)
func (_VerifierProxy *VerifierProxyFilterer) FilterVerifierUnset(opts *bind.FilterOpts) (*VerifierProxyVerifierUnsetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyVerifierUnsetIterator{contract: _VerifierProxy.contract, event: "VerifierUnset", logs: logs, sub: sub}, nil
}

// WatchVerifierUnset is a free log subscription operation binding the contract event 0x11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c.
//
// Solidity: event VerifierUnset(bytes32 configDigest, address verifierAddress)
func (_VerifierProxy *VerifierProxyFilterer) WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierUnset) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierProxyVerifierUnset)
				if err := _VerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
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

// ParseVerifierUnset is a log parse operation binding the contract event 0x11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c.
//
// Solidity: event VerifierUnset(bytes32 configDigest, address verifierAddress)
func (_VerifierProxy *VerifierProxyFilterer) ParseVerifierUnset(log types.Log) (*VerifierProxyVerifierUnset, error) {
	event := new(VerifierProxyVerifierUnset)
	if err := _VerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
