// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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

// DeviationFlaggingValidatorMetaData contains all meta data concerning the DeviationFlaggingValidator contract.
var DeviationFlaggingValidatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_flags\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"_flaggingThreshold\",\"type\":\"uint24\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"previous\",\"type\":\"uint24\"},{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"current\",\"type\":\"uint24\"}],\"name\":\"FlaggingThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"FlagsAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"THRESHOLD_MULTIPLIER\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"flaggingThreshold\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"flags\",\"outputs\":[{\"internalType\":\"contractFlagsInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"isValid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"_flaggingThreshold\",\"type\":\"uint24\"}],\"name\":\"setFlaggingThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_flags\",\"type\":\"address\"}],\"name\":\"setFlagsAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516109ae3803806109ae8339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b03191633179055610060826001600160e01b0361007916565b610072816001600160e01b0361014116565b5050610215565b6000546001600160a01b031633146100d8576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6002546001600160a01b03908116908216811461013d57600280546001600160a01b0319166001600160a01b0384811691821790925560405190918316907f900aa01828592ab069e4d44e7a36c70ebd476e35f567c7db6a691e503b8029d890600090a35b5050565b6000546001600160a01b031633146101a0576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600154600160a01b900463ffffffff81169062ffffff83811691161461013d576001805463ffffffff60a01b191662ffffff848116600160a01b81029290921790925560405190918316907f985b87e809fd5992ec257eac36f25777ce308055dd9249a0182123d7b5d6633a90600090a35050565b61078a806102246000396000f3fe608060405234801561001057600080fd5b506004361061008e5760003560e01c80630910ce4a1461009357806364cc4aa5146100b457806379ba5097146100d85780638da5cb5b146100e2578063beed9b51146100ea578063eed8a1de1461012d578063f19876951461014f578063f2c0ea9214610175578063f2fde38b1461017d578063ffd93670146101a3575b600080fd5b61009b6101d2565b6040805163ffffffff9092168252519081900360200190f35b6100bc6101e5565b604080516001600160a01b039092168252519081900360200190f35b6100e06101f4565b005b6100bc6102a3565b6101196004803603608081101561010057600080fd5b50803590602081013590604081013590606001356102b2565b604080519115158252519081900360200190f35b6100e06004803603602081101561014357600080fd5b503562ffffff1661033d565b6100e06004803603602081101561016557600080fd5b50356001600160a01b031661040b565b61009b6104cb565b6100e06004803603602081101561019357600080fd5b50356001600160a01b03166104d2565b610119600480360360808110156101b957600080fd5b508035906020810135906040810135906060013561057b565b600154600160a01b900463ffffffff1681565b6002546001600160a01b031681565b6001546001600160a01b0316331461024c576040805162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015290519081900360640190fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b031681565b60006102c08585858561057b565b610331576002546040805163d74af26360e01b815233600482015290516001600160a01b039092169163d74af2639160248082019260009290919082900301818387803b15801561031057600080fd5b505af1158015610324573d6000803e3d6000fd5b5050505060009050610335565b5060015b949350505050565b6000546001600160a01b03163314610395576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b600154600160a01b900463ffffffff81169062ffffff838116911614610407576001805463ffffffff60a01b191662ffffff848116600160a01b81029290921790925560405190918316907f985b87e809fd5992ec257eac36f25777ce308055dd9249a0182123d7b5d6633a90600090a35b5050565b6000546001600160a01b03163314610463576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6002546001600160a01b03908116908216811461040757600280546001600160a01b0319166001600160a01b0384811691821790925560405190918316907f900aa01828592ab069e4d44e7a36c70ebd476e35f567c7db6a691e503b8029d890600090a35050565b620186a081565b6000546001600160a01b0316331461052a576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008361058a57506001610335565b60008061059d868563ffffffff61062616565b90925090506000806105b284620186a0610671565b90925090506000806105ca848b63ffffffff6106c916565b915091506000806105da84610718565b915091508680156105e85750845b80156105f15750825b80156105fa5750805b80156106155750600154600160a01b900463ffffffff168211155b9d9c50505050505050505050505050565b600080828403818412801561063b5750848113155b8061065157506000841215801561065157508481135b1561066357506000915081905061066a565b9150600190505b9250929050565b60008083610685575060009050600161066a565b836000191480156106995750600160ff1b83145b156106a95750600090508061066a565b838302838582816106b657fe5b051461066357506000915081905061066a565b600080826106dc5750600090508061066a565b826000191480156106f05750600160ff1b84145b156107005750600090508061066a565b600083858161070b57fe5b0595600195509350505050565b6000806000831261072e5750819050600161074f565b600160ff1b8314156107455750600090508061074f565b5050600081900360015b91509156fea264697066735822122007573b314c61efd52358d5489743c155f1ee66817e16f9230127c74aad837b6a64736f6c63430006060033",
}

// DeviationFlaggingValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use DeviationFlaggingValidatorMetaData.ABI instead.
var DeviationFlaggingValidatorABI = DeviationFlaggingValidatorMetaData.ABI

// DeviationFlaggingValidatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DeviationFlaggingValidatorMetaData.Bin instead.
var DeviationFlaggingValidatorBin = DeviationFlaggingValidatorMetaData.Bin

// DeployDeviationFlaggingValidator deploys a new Ethereum contract, binding an instance of DeviationFlaggingValidator to it.
func DeployDeviationFlaggingValidator(auth *bind.TransactOpts, backend bind.ContractBackend, _flags common.Address, _flaggingThreshold *big.Int) (common.Address, *types.Transaction, *DeviationFlaggingValidator, error) {
	parsed, err := DeviationFlaggingValidatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DeviationFlaggingValidatorBin), backend, _flags, _flaggingThreshold)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DeviationFlaggingValidator{DeviationFlaggingValidatorCaller: DeviationFlaggingValidatorCaller{contract: contract}, DeviationFlaggingValidatorTransactor: DeviationFlaggingValidatorTransactor{contract: contract}, DeviationFlaggingValidatorFilterer: DeviationFlaggingValidatorFilterer{contract: contract}}, nil
}

// DeviationFlaggingValidator is an auto generated Go binding around an Ethereum contract.
type DeviationFlaggingValidator struct {
	DeviationFlaggingValidatorCaller     // Read-only binding to the contract
	DeviationFlaggingValidatorTransactor // Write-only binding to the contract
	DeviationFlaggingValidatorFilterer   // Log filterer for contract events
}

// DeviationFlaggingValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DeviationFlaggingValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeviationFlaggingValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DeviationFlaggingValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeviationFlaggingValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DeviationFlaggingValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeviationFlaggingValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DeviationFlaggingValidatorSession struct {
	Contract     *DeviationFlaggingValidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DeviationFlaggingValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DeviationFlaggingValidatorCallerSession struct {
	Contract *DeviationFlaggingValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// DeviationFlaggingValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DeviationFlaggingValidatorTransactorSession struct {
	Contract     *DeviationFlaggingValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// DeviationFlaggingValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DeviationFlaggingValidatorRaw struct {
	Contract *DeviationFlaggingValidator // Generic contract binding to access the raw methods on
}

// DeviationFlaggingValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DeviationFlaggingValidatorCallerRaw struct {
	Contract *DeviationFlaggingValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// DeviationFlaggingValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DeviationFlaggingValidatorTransactorRaw struct {
	Contract *DeviationFlaggingValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeviationFlaggingValidator creates a new instance of DeviationFlaggingValidator, bound to a specific deployed contract.
func NewDeviationFlaggingValidator(address common.Address, backend bind.ContractBackend) (*DeviationFlaggingValidator, error) {
	contract, err := bindDeviationFlaggingValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidator{DeviationFlaggingValidatorCaller: DeviationFlaggingValidatorCaller{contract: contract}, DeviationFlaggingValidatorTransactor: DeviationFlaggingValidatorTransactor{contract: contract}, DeviationFlaggingValidatorFilterer: DeviationFlaggingValidatorFilterer{contract: contract}}, nil
}

// NewDeviationFlaggingValidatorCaller creates a new read-only instance of DeviationFlaggingValidator, bound to a specific deployed contract.
func NewDeviationFlaggingValidatorCaller(address common.Address, caller bind.ContractCaller) (*DeviationFlaggingValidatorCaller, error) {
	contract, err := bindDeviationFlaggingValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorCaller{contract: contract}, nil
}

// NewDeviationFlaggingValidatorTransactor creates a new write-only instance of DeviationFlaggingValidator, bound to a specific deployed contract.
func NewDeviationFlaggingValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*DeviationFlaggingValidatorTransactor, error) {
	contract, err := bindDeviationFlaggingValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorTransactor{contract: contract}, nil
}

// NewDeviationFlaggingValidatorFilterer creates a new log filterer instance of DeviationFlaggingValidator, bound to a specific deployed contract.
func NewDeviationFlaggingValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*DeviationFlaggingValidatorFilterer, error) {
	contract, err := bindDeviationFlaggingValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorFilterer{contract: contract}, nil
}

// bindDeviationFlaggingValidator binds a generic wrapper to an already deployed contract.
func bindDeviationFlaggingValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DeviationFlaggingValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeviationFlaggingValidator.Contract.DeviationFlaggingValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.DeviationFlaggingValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.DeviationFlaggingValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeviationFlaggingValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.contract.Transact(opts, method, params...)
}

// THRESHOLDMULTIPLIER is a free data retrieval call binding the contract method 0xf2c0ea92.
//
// Solidity: function THRESHOLD_MULTIPLIER() view returns(uint32)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCaller) THRESHOLDMULTIPLIER(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _DeviationFlaggingValidator.contract.Call(opts, &out, "THRESHOLD_MULTIPLIER")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// THRESHOLDMULTIPLIER is a free data retrieval call binding the contract method 0xf2c0ea92.
//
// Solidity: function THRESHOLD_MULTIPLIER() view returns(uint32)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) THRESHOLDMULTIPLIER() (uint32, error) {
	return _DeviationFlaggingValidator.Contract.THRESHOLDMULTIPLIER(&_DeviationFlaggingValidator.CallOpts)
}

// THRESHOLDMULTIPLIER is a free data retrieval call binding the contract method 0xf2c0ea92.
//
// Solidity: function THRESHOLD_MULTIPLIER() view returns(uint32)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCallerSession) THRESHOLDMULTIPLIER() (uint32, error) {
	return _DeviationFlaggingValidator.Contract.THRESHOLDMULTIPLIER(&_DeviationFlaggingValidator.CallOpts)
}

// FlaggingThreshold is a free data retrieval call binding the contract method 0x0910ce4a.
//
// Solidity: function flaggingThreshold() view returns(uint32)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCaller) FlaggingThreshold(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _DeviationFlaggingValidator.contract.Call(opts, &out, "flaggingThreshold")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// FlaggingThreshold is a free data retrieval call binding the contract method 0x0910ce4a.
//
// Solidity: function flaggingThreshold() view returns(uint32)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) FlaggingThreshold() (uint32, error) {
	return _DeviationFlaggingValidator.Contract.FlaggingThreshold(&_DeviationFlaggingValidator.CallOpts)
}

// FlaggingThreshold is a free data retrieval call binding the contract method 0x0910ce4a.
//
// Solidity: function flaggingThreshold() view returns(uint32)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCallerSession) FlaggingThreshold() (uint32, error) {
	return _DeviationFlaggingValidator.Contract.FlaggingThreshold(&_DeviationFlaggingValidator.CallOpts)
}

// Flags is a free data retrieval call binding the contract method 0x64cc4aa5.
//
// Solidity: function flags() view returns(address)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCaller) Flags(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DeviationFlaggingValidator.contract.Call(opts, &out, "flags")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Flags is a free data retrieval call binding the contract method 0x64cc4aa5.
//
// Solidity: function flags() view returns(address)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) Flags() (common.Address, error) {
	return _DeviationFlaggingValidator.Contract.Flags(&_DeviationFlaggingValidator.CallOpts)
}

// Flags is a free data retrieval call binding the contract method 0x64cc4aa5.
//
// Solidity: function flags() view returns(address)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCallerSession) Flags() (common.Address, error) {
	return _DeviationFlaggingValidator.Contract.Flags(&_DeviationFlaggingValidator.CallOpts)
}

// IsValid is a free data retrieval call binding the contract method 0xffd93670.
//
// Solidity: function isValid(uint256 , int256 _previousAnswer, uint256 , int256 _answer) view returns(bool)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCaller) IsValid(opts *bind.CallOpts, arg0 *big.Int, _previousAnswer *big.Int, arg2 *big.Int, _answer *big.Int) (bool, error) {
	var out []interface{}
	err := _DeviationFlaggingValidator.contract.Call(opts, &out, "isValid", arg0, _previousAnswer, arg2, _answer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValid is a free data retrieval call binding the contract method 0xffd93670.
//
// Solidity: function isValid(uint256 , int256 _previousAnswer, uint256 , int256 _answer) view returns(bool)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) IsValid(arg0 *big.Int, _previousAnswer *big.Int, arg2 *big.Int, _answer *big.Int) (bool, error) {
	return _DeviationFlaggingValidator.Contract.IsValid(&_DeviationFlaggingValidator.CallOpts, arg0, _previousAnswer, arg2, _answer)
}

// IsValid is a free data retrieval call binding the contract method 0xffd93670.
//
// Solidity: function isValid(uint256 , int256 _previousAnswer, uint256 , int256 _answer) view returns(bool)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCallerSession) IsValid(arg0 *big.Int, _previousAnswer *big.Int, arg2 *big.Int, _answer *big.Int) (bool, error) {
	return _DeviationFlaggingValidator.Contract.IsValid(&_DeviationFlaggingValidator.CallOpts, arg0, _previousAnswer, arg2, _answer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DeviationFlaggingValidator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) Owner() (common.Address, error) {
	return _DeviationFlaggingValidator.Contract.Owner(&_DeviationFlaggingValidator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorCallerSession) Owner() (common.Address, error) {
	return _DeviationFlaggingValidator.Contract.Owner(&_DeviationFlaggingValidator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.AcceptOwnership(&_DeviationFlaggingValidator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.AcceptOwnership(&_DeviationFlaggingValidator.TransactOpts)
}

// SetFlaggingThreshold is a paid mutator transaction binding the contract method 0xeed8a1de.
//
// Solidity: function setFlaggingThreshold(uint24 _flaggingThreshold) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactor) SetFlaggingThreshold(opts *bind.TransactOpts, _flaggingThreshold *big.Int) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.contract.Transact(opts, "setFlaggingThreshold", _flaggingThreshold)
}

// SetFlaggingThreshold is a paid mutator transaction binding the contract method 0xeed8a1de.
//
// Solidity: function setFlaggingThreshold(uint24 _flaggingThreshold) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) SetFlaggingThreshold(_flaggingThreshold *big.Int) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.SetFlaggingThreshold(&_DeviationFlaggingValidator.TransactOpts, _flaggingThreshold)
}

// SetFlaggingThreshold is a paid mutator transaction binding the contract method 0xeed8a1de.
//
// Solidity: function setFlaggingThreshold(uint24 _flaggingThreshold) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorSession) SetFlaggingThreshold(_flaggingThreshold *big.Int) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.SetFlaggingThreshold(&_DeviationFlaggingValidator.TransactOpts, _flaggingThreshold)
}

// SetFlagsAddress is a paid mutator transaction binding the contract method 0xf1987695.
//
// Solidity: function setFlagsAddress(address _flags) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactor) SetFlagsAddress(opts *bind.TransactOpts, _flags common.Address) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.contract.Transact(opts, "setFlagsAddress", _flags)
}

// SetFlagsAddress is a paid mutator transaction binding the contract method 0xf1987695.
//
// Solidity: function setFlagsAddress(address _flags) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) SetFlagsAddress(_flags common.Address) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.SetFlagsAddress(&_DeviationFlaggingValidator.TransactOpts, _flags)
}

// SetFlagsAddress is a paid mutator transaction binding the contract method 0xf1987695.
//
// Solidity: function setFlagsAddress(address _flags) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorSession) SetFlagsAddress(_flags common.Address) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.SetFlagsAddress(&_DeviationFlaggingValidator.TransactOpts, _flags)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.TransferOwnership(&_DeviationFlaggingValidator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.TransferOwnership(&_DeviationFlaggingValidator.TransactOpts, _to)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 _previousRoundId, int256 _previousAnswer, uint256 _roundId, int256 _answer) returns(bool)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactor) Validate(opts *bind.TransactOpts, _previousRoundId *big.Int, _previousAnswer *big.Int, _roundId *big.Int, _answer *big.Int) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.contract.Transact(opts, "validate", _previousRoundId, _previousAnswer, _roundId, _answer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 _previousRoundId, int256 _previousAnswer, uint256 _roundId, int256 _answer) returns(bool)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorSession) Validate(_previousRoundId *big.Int, _previousAnswer *big.Int, _roundId *big.Int, _answer *big.Int) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.Validate(&_DeviationFlaggingValidator.TransactOpts, _previousRoundId, _previousAnswer, _roundId, _answer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 _previousRoundId, int256 _previousAnswer, uint256 _roundId, int256 _answer) returns(bool)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorTransactorSession) Validate(_previousRoundId *big.Int, _previousAnswer *big.Int, _roundId *big.Int, _answer *big.Int) (*types.Transaction, error) {
	return _DeviationFlaggingValidator.Contract.Validate(&_DeviationFlaggingValidator.TransactOpts, _previousRoundId, _previousAnswer, _roundId, _answer)
}

// DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator is returned from FilterFlaggingThresholdUpdated and is used to iterate over the raw logs and unpacked data for FlaggingThresholdUpdated events raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator struct {
	Event *DeviationFlaggingValidatorFlaggingThresholdUpdated // Event containing the contract specifics and raw log

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
func (it *DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeviationFlaggingValidatorFlaggingThresholdUpdated)
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
		it.Event = new(DeviationFlaggingValidatorFlaggingThresholdUpdated)
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
func (it *DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeviationFlaggingValidatorFlaggingThresholdUpdated represents a FlaggingThresholdUpdated event raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorFlaggingThresholdUpdated struct {
	Previous *big.Int
	Current  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFlaggingThresholdUpdated is a free log retrieval operation binding the contract event 0x985b87e809fd5992ec257eac36f25777ce308055dd9249a0182123d7b5d6633a.
//
// Solidity: event FlaggingThresholdUpdated(uint24 indexed previous, uint24 indexed current)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) FilterFlaggingThresholdUpdated(opts *bind.FilterOpts, previous []*big.Int, current []*big.Int) (*DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.FilterLogs(opts, "FlaggingThresholdUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorFlaggingThresholdUpdatedIterator{contract: _DeviationFlaggingValidator.contract, event: "FlaggingThresholdUpdated", logs: logs, sub: sub}, nil
}

// WatchFlaggingThresholdUpdated is a free log subscription operation binding the contract event 0x985b87e809fd5992ec257eac36f25777ce308055dd9249a0182123d7b5d6633a.
//
// Solidity: event FlaggingThresholdUpdated(uint24 indexed previous, uint24 indexed current)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) WatchFlaggingThresholdUpdated(opts *bind.WatchOpts, sink chan<- *DeviationFlaggingValidatorFlaggingThresholdUpdated, previous []*big.Int, current []*big.Int) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.WatchLogs(opts, "FlaggingThresholdUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeviationFlaggingValidatorFlaggingThresholdUpdated)
				if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "FlaggingThresholdUpdated", log); err != nil {
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

// ParseFlaggingThresholdUpdated is a log parse operation binding the contract event 0x985b87e809fd5992ec257eac36f25777ce308055dd9249a0182123d7b5d6633a.
//
// Solidity: event FlaggingThresholdUpdated(uint24 indexed previous, uint24 indexed current)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) ParseFlaggingThresholdUpdated(log types.Log) (*DeviationFlaggingValidatorFlaggingThresholdUpdated, error) {
	event := new(DeviationFlaggingValidatorFlaggingThresholdUpdated)
	if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "FlaggingThresholdUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeviationFlaggingValidatorFlagsAddressUpdatedIterator is returned from FilterFlagsAddressUpdated and is used to iterate over the raw logs and unpacked data for FlagsAddressUpdated events raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorFlagsAddressUpdatedIterator struct {
	Event *DeviationFlaggingValidatorFlagsAddressUpdated // Event containing the contract specifics and raw log

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
func (it *DeviationFlaggingValidatorFlagsAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeviationFlaggingValidatorFlagsAddressUpdated)
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
		it.Event = new(DeviationFlaggingValidatorFlagsAddressUpdated)
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
func (it *DeviationFlaggingValidatorFlagsAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeviationFlaggingValidatorFlagsAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeviationFlaggingValidatorFlagsAddressUpdated represents a FlagsAddressUpdated event raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorFlagsAddressUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFlagsAddressUpdated is a free log retrieval operation binding the contract event 0x900aa01828592ab069e4d44e7a36c70ebd476e35f567c7db6a691e503b8029d8.
//
// Solidity: event FlagsAddressUpdated(address indexed previous, address indexed current)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) FilterFlagsAddressUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*DeviationFlaggingValidatorFlagsAddressUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.FilterLogs(opts, "FlagsAddressUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorFlagsAddressUpdatedIterator{contract: _DeviationFlaggingValidator.contract, event: "FlagsAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchFlagsAddressUpdated is a free log subscription operation binding the contract event 0x900aa01828592ab069e4d44e7a36c70ebd476e35f567c7db6a691e503b8029d8.
//
// Solidity: event FlagsAddressUpdated(address indexed previous, address indexed current)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) WatchFlagsAddressUpdated(opts *bind.WatchOpts, sink chan<- *DeviationFlaggingValidatorFlagsAddressUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.WatchLogs(opts, "FlagsAddressUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeviationFlaggingValidatorFlagsAddressUpdated)
				if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "FlagsAddressUpdated", log); err != nil {
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

// ParseFlagsAddressUpdated is a log parse operation binding the contract event 0x900aa01828592ab069e4d44e7a36c70ebd476e35f567c7db6a691e503b8029d8.
//
// Solidity: event FlagsAddressUpdated(address indexed previous, address indexed current)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) ParseFlagsAddressUpdated(log types.Log) (*DeviationFlaggingValidatorFlagsAddressUpdated, error) {
	event := new(DeviationFlaggingValidatorFlagsAddressUpdated)
	if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "FlagsAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeviationFlaggingValidatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorOwnershipTransferRequestedIterator struct {
	Event *DeviationFlaggingValidatorOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *DeviationFlaggingValidatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeviationFlaggingValidatorOwnershipTransferRequested)
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
		it.Event = new(DeviationFlaggingValidatorOwnershipTransferRequested)
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
func (it *DeviationFlaggingValidatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeviationFlaggingValidatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeviationFlaggingValidatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DeviationFlaggingValidatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorOwnershipTransferRequestedIterator{contract: _DeviationFlaggingValidator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DeviationFlaggingValidatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeviationFlaggingValidatorOwnershipTransferRequested)
				if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*DeviationFlaggingValidatorOwnershipTransferRequested, error) {
	event := new(DeviationFlaggingValidatorOwnershipTransferRequested)
	if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeviationFlaggingValidatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorOwnershipTransferredIterator struct {
	Event *DeviationFlaggingValidatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DeviationFlaggingValidatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeviationFlaggingValidatorOwnershipTransferred)
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
		it.Event = new(DeviationFlaggingValidatorOwnershipTransferred)
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
func (it *DeviationFlaggingValidatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeviationFlaggingValidatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeviationFlaggingValidatorOwnershipTransferred represents a OwnershipTransferred event raised by the DeviationFlaggingValidator contract.
type DeviationFlaggingValidatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DeviationFlaggingValidatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DeviationFlaggingValidatorOwnershipTransferredIterator{contract: _DeviationFlaggingValidator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DeviationFlaggingValidatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DeviationFlaggingValidator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeviationFlaggingValidatorOwnershipTransferred)
				if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DeviationFlaggingValidator *DeviationFlaggingValidatorFilterer) ParseOwnershipTransferred(log types.Log) (*DeviationFlaggingValidatorOwnershipTransferred, error) {
	event := new(DeviationFlaggingValidatorOwnershipTransferred)
	if err := _DeviationFlaggingValidator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
