// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// VRFABI is the input ABI used to generate the binding from.
const VRFABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"randomValueFromVRFProof\",\"outputs\":[{\"name\":\"output\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// VRFBin is the compiled bytecode used for deploying new contracts.
const VRFBin = `0x608060405234801561001057600080fd5b50610f73806100206000396000f3fe6080604052600436106100405763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663fa8fc6f18114610045575b600080fd5b34801561005157600080fd5b506100f86004803603602081101561006857600080fd5b81019060208101813564010000000081111561008357600080fd5b82018360208201111561009557600080fd5b803590602001918460018302840111640100000000831117156100b757600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061010a945050505050565b60408051918252519081900360200190f35b80516000906101a014610167576040805160e560020a62461bcd02815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b61016f610ecf565b610177610ecf565b61017f610eea565b6000610189610ecf565b610191610ecf565b6000888060200190516101a08110156101a957600080fd5b5060e08101516101808201519198506040890197506080890196509450610100880193506101408801925090506101fc87878760006020020151886001602002015189600260200201518989898961025a565b856040516020018082600260200280838360005b83811015610228578181015183820152602001610210565b505050509050019150506040516020818303038152906040528051906020012060019004975050505050505050919050565b610263896104d5565b15156102b9576040805160e560020a62461bcd02815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015290519081900360640190fd5b6102c2886104d5565b1515610318576040805160e560020a62461bcd02815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015290519081900360640190fd5b610321836104d5565b1515610377576040805160e560020a62461bcd02815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015290519081900360640190fd5b610380826104d5565b15156103d6576040805160e560020a62461bcd02815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015290519081900360640190fd5b6103e2878a8887610501565b1515610438576040805160e560020a62461bcd02815260206004820152601a60248201527f6164647228632a706b2b732a6729e289a05f755769746e657373000000000000604482015290519081900360640190fd5b610440610ecf565b61044a8a8761067a565b9050610454610ecf565b610463898b878b868989610784565b9050610472828c8c8985610916565b89146104c8576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015290519081900360640190fd5b5050505050505050505050565b60208101516000906401000003d0199080096104f88360005b6020020151610a3d565b1490505b919050565b600073ffffffffffffffffffffffffffffffffffffffff82161515610570576040805160e560020a62461bcd02815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015290519081900360640190fd5b835160009070014551231950b75fc4402da1732fc9bebe199085900970014551231950b75fc4402da1732fc9bebe190390506000600286600160200201518115156105b757fe5b06156105c457601c6105c7565b601b5b865190915060009070014551231950b75fc4402da1732fc9bebe199089098751604080516000808252602082810180855289905260ff8816838501526060830194909452608082018590529151939450909260019260a0808401939192601f1981019281900390910190855afa158015610645573d6000803e3d6000fd5b5050604051601f19015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b610682610ecf565b6106dc83836040516020018083600260200280838360005b838110156106b257818101518382015260200161069a565b50505050919091019283525050604080518083038152602092830190915280519101209050610a61565b81526106f16106ec8260006104ee565b610a96565b60208201525b6107028160006104ee565b60208201516401000003d0199080091461075c57805160408051602081810193909352815180820384018152908201909152805191012061074290610a61565b81526107526106ec8260006104ee565b60208201526106f7565b6020810151600290066001141561077e576020810180516401000003d0190390525b92915050565b61078c610ecf565b825186516401000003d0199190030615156107f1576040805160e560020a62461bcd02815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015290519081900360640190fd5b6107fc878988610ac2565b1515610878576040805160e560020a62461bcd02815260206004820152602160248201527f4669727374206d756c7469706c69636174696f6e20636865636b206661696c6560448201527f6400000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b610883848685610ac2565b15156108ff576040805160e560020a62461bcd02815260206004820152602260248201527f5365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c60448201527f6564000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b61090a868484610c14565b98975050505050505050565b600085858584866040516020018086600260200280838360005b83811015610948578181015183820152602001610930565b5050505090500185600260200280838360005b8381101561097357818101518382015260200161095b565b5050505090500184600260200280838360005b8381101561099e578181015183820152602001610986565b5050505090500183600260200280838360005b838110156109c95781810151838201526020016109b1565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401955050505050506040516020818303038152906040528051906020012060019004905095945050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b805b6401000003d01981106104fc57604080516020808201939093528151808203840181529082019091528051910120610a63565b600061077e827f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c610cd6565b6000821515610ad057600080fd5b835160009070014551231950b75fc4402da1732fc9bebe1990850960208601519091506000906001161515610b0657601b610b09565b601c5b9050836040516020018082600260200280838360005b83811015610b37578181015183820152602001610b1f565b50505050905001915050604051602081830303815290604052805190602001206001900473ffffffffffffffffffffffffffffffffffffffff166001600060010283896000600281101515610b8857fe5b60200201516001028660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610be8573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614925050509392505050565b610c1c610ecf565b835160208086015185519186015160009384938493610c3d93909190610d82565b919450925090506401000003d019858209600114610ca5576040805160e560020a62461bcd02815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a00000000000000604482015290519081900360640190fd5b60408051808201909152806401000003d01987860981526020016401000003d0198785099052979650505050505050565b600080610ce1610f09565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152610d13610f28565b60208160c0846005600019fa9250821515610d78576040805160e560020a62461bcd02815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015290519081900360640190fd5b5195945050505050565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000610dc283838585610e62565b9098509050610dd388828e88610e86565b9098509050610de488828c87610e86565b90985090506000610df78d878b85610e86565b9098509050610e0888828686610e62565b9098509050610e1988828e89610e86565b9098509050818114610e4e576401000003d019818a0998506401000003d01982890997506401000003d0198183099650610e52565b8196505b5050505050509450945094915050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b60408051808201825290600290829080388339509192915050565b6060604051908101604052806003906020820280388339509192915050565b60c0604051908101604052806006906020820280388339509192915050565b602060405190810160405280600190602082028038833950919291505056fea165627a7a7230582044d93a58b88b7274286429e4adef12c475e0cda10386d507af47ec3d61d1a8fd0029`

// DeployVRF deploys a new Ethereum contract, binding an instance of VRF to it.
func DeployVRF(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRF, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRF{VRFCaller: VRFCaller{contract: contract}, VRFTransactor: VRFTransactor{contract: contract}, VRFFilterer: VRFFilterer{contract: contract}}, nil
}

// VRF is an auto generated Go binding around an Ethereum contract.
type VRF struct {
	VRFCaller     // Read-only binding to the contract
	VRFTransactor // Write-only binding to the contract
	VRFFilterer   // Log filterer for contract events
}

// VRFCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFSession struct {
	Contract     *VRF              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFCallerSession struct {
	Contract *VRFCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VRFTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFTransactorSession struct {
	Contract     *VRFTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFRaw struct {
	Contract *VRF // Generic contract binding to access the raw methods on
}

// VRFCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFCallerRaw struct {
	Contract *VRFCaller // Generic read-only contract binding to access the raw methods on
}

// VRFTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFTransactorRaw struct {
	Contract *VRFTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRF creates a new instance of VRF, bound to a specific deployed contract.
func NewVRF(address common.Address, backend bind.ContractBackend) (*VRF, error) {
	contract, err := bindVRF(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRF{VRFCaller: VRFCaller{contract: contract}, VRFTransactor: VRFTransactor{contract: contract}, VRFFilterer: VRFFilterer{contract: contract}}, nil
}

// NewVRFCaller creates a new read-only instance of VRF, bound to a specific deployed contract.
func NewVRFCaller(address common.Address, caller bind.ContractCaller) (*VRFCaller, error) {
	contract, err := bindVRF(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCaller{contract: contract}, nil
}

// NewVRFTransactor creates a new write-only instance of VRF, bound to a specific deployed contract.
func NewVRFTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFTransactor, error) {
	contract, err := bindVRF(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTransactor{contract: contract}, nil
}

// NewVRFFilterer creates a new log filterer instance of VRF, bound to a specific deployed contract.
func NewVRFFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFFilterer, error) {
	contract, err := bindVRF(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFFilterer{contract: contract}, nil
}

// bindVRF binds a generic wrapper to an already deployed contract.
func bindVRF(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRF *VRFRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRF.Contract.VRFCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRF *VRFRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRF.Contract.VRFTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRF *VRFRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRF.Contract.VRFTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRF *VRFCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRF.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRF *VRFTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRF.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRF *VRFTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRF.Contract.contract.Transact(opts, method, params...)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRF *VRFCaller) RandomValueFromVRFProof(opts *bind.CallOpts, proof []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRF.contract.Call(opts, out, "randomValueFromVRFProof", proof)
	return *ret0, err
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRF *VRFSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRF.Contract.RandomValueFromVRFProof(&_VRF.CallOpts, proof)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRF *VRFCallerSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRF.Contract.RandomValueFromVRFProof(&_VRF.CallOpts, proof)
}

// VRFAllABI is the input ABI used to generate the binding from.
const VRFAllABI = "[]"

// VRFAllBin is the compiled bytecode used for deploying new contracts.
const VRFAllBin = `0x6080604052348015600f57600080fd5b50603580601d6000396000f3fe6080604052600080fdfea165627a7a723058204219d16ad0f57cfab2731777b441310f5fb75736ccfd413a0bae7d26783f77470029`

// DeployVRFAll deploys a new Ethereum contract, binding an instance of VRFAll to it.
func DeployVRFAll(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFAll, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFAllABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFAllBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFAll{VRFAllCaller: VRFAllCaller{contract: contract}, VRFAllTransactor: VRFAllTransactor{contract: contract}, VRFAllFilterer: VRFAllFilterer{contract: contract}}, nil
}

// VRFAll is an auto generated Go binding around an Ethereum contract.
type VRFAll struct {
	VRFAllCaller     // Read-only binding to the contract
	VRFAllTransactor // Write-only binding to the contract
	VRFAllFilterer   // Log filterer for contract events
}

// VRFAllCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFAllCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFAllTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFAllTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFAllFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFAllFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFAllSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFAllSession struct {
	Contract     *VRFAll           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFAllCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFAllCallerSession struct {
	Contract *VRFAllCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VRFAllTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFAllTransactorSession struct {
	Contract     *VRFAllTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFAllRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFAllRaw struct {
	Contract *VRFAll // Generic contract binding to access the raw methods on
}

// VRFAllCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFAllCallerRaw struct {
	Contract *VRFAllCaller // Generic read-only contract binding to access the raw methods on
}

// VRFAllTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFAllTransactorRaw struct {
	Contract *VRFAllTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFAll creates a new instance of VRFAll, bound to a specific deployed contract.
func NewVRFAll(address common.Address, backend bind.ContractBackend) (*VRFAll, error) {
	contract, err := bindVRFAll(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFAll{VRFAllCaller: VRFAllCaller{contract: contract}, VRFAllTransactor: VRFAllTransactor{contract: contract}, VRFAllFilterer: VRFAllFilterer{contract: contract}}, nil
}

// NewVRFAllCaller creates a new read-only instance of VRFAll, bound to a specific deployed contract.
func NewVRFAllCaller(address common.Address, caller bind.ContractCaller) (*VRFAllCaller, error) {
	contract, err := bindVRFAll(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFAllCaller{contract: contract}, nil
}

// NewVRFAllTransactor creates a new write-only instance of VRFAll, bound to a specific deployed contract.
func NewVRFAllTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFAllTransactor, error) {
	contract, err := bindVRFAll(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFAllTransactor{contract: contract}, nil
}

// NewVRFAllFilterer creates a new log filterer instance of VRFAll, bound to a specific deployed contract.
func NewVRFAllFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFAllFilterer, error) {
	contract, err := bindVRFAll(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFAllFilterer{contract: contract}, nil
}

// bindVRFAll binds a generic wrapper to an already deployed contract.
func bindVRFAll(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFAllABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFAll *VRFAllRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFAll.Contract.VRFAllCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFAll *VRFAllRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFAll.Contract.VRFAllTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFAll *VRFAllRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFAll.Contract.VRFAllTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFAll *VRFAllCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFAll.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFAll *VRFAllTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFAll.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFAll *VRFAllTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFAll.Contract.contract.Transact(opts, method, params...)
}

// VRFTestHelperABI is the input ABI used to generate the binding from.
const VRFTestHelperABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"name\":\"invZ\",\"type\":\"uint256\"}],\"name\":\"affineECAdd_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"zqHash_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"hashToCurve_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"uint256[2]\"},{\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"name\":\"uWitness\",\"type\":\"address\"},{\"name\":\"v\",\"type\":\"uint256[2]\"}],\"name\":\"scalarFromCurve_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"base\",\"type\":\"uint256\"},{\"name\":\"exponent\",\"type\":\"uint256\"}],\"name\":\"bigModExp_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"squareRoot_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"c\",\"type\":\"uint256\"},{\"name\":\"p\",\"type\":\"uint256[2]\"},{\"name\":\"s\",\"type\":\"uint256\"},{\"name\":\"lcWitness\",\"type\":\"address\"}],\"name\":\"verifyLinearCombinationWithGenerator_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"px\",\"type\":\"uint256\"},{\"name\":\"py\",\"type\":\"uint256\"},{\"name\":\"qx\",\"type\":\"uint256\"},{\"name\":\"qy\",\"type\":\"uint256\"}],\"name\":\"projectiveECAdd_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"ySquared_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256[2]\"},{\"name\":\"scalar\",\"type\":\"uint256\"},{\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"ecmulVerify_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"name\":\"c\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"},{\"name\":\"seed\",\"type\":\"uint256\"},{\"name\":\"uWitness\",\"type\":\"address\"},{\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"verifyVRFProof_\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"randomValueFromVRFProof\",\"outputs\":[{\"name\":\"output\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"c\",\"type\":\"uint256\"},{\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"name\":\"cp1Witness\",\"type\":\"uint256[2]\"},{\"name\":\"s\",\"type\":\"uint256\"},{\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"name\":\"sp2Witness\",\"type\":\"uint256[2]\"},{\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"linearCombination_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// VRFTestHelperBin is the compiled bytecode used for deploying new contracts.
const VRFTestHelperBin = `0x608060405234801561001057600080fd5b5061169d806100206000396000f3fe6080604052600436106100c45763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663244f896d81146100c957806324d72ea91461018757806335452450146101c3578063525413cf1461021d5780635de60042146103055780638af046ea1461033557806391d5f6911461035f57806395e6ee92146103e15780639d6f03371461043b578063aa7b2fbb14610465578063ef3b10ec146104f1578063fa8fc6f1146105ed578063fe54f2a2146106a0575b600080fd5b3480156100d557600080fd5b5061014c600480360360a08110156100ec57600080fd5b60408051808201825291830192918183019183906002908390839080828437600092019190915250506040805180820182529295949381810193925090600290839083908082843760009201919091525091945050903591506107899050565b6040518082600260200280838360005b8381101561017457818101518382015260200161015c565b5050505090500191505060405180910390f35b34801561019357600080fd5b506101b1600480360360208110156101aa57600080fd5b50356107a4565b60408051918252519081900360200190f35b3480156101cf57600080fd5b5061014c600480360360608110156101e657600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525091945050903591506107b79050565b34801561022957600080fd5b506101b1600480360361012081101561024157600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525050604080518082018252929594938181019392509060029083908390808284376000920191909152505060408051808201825292959493818101939250906002908390839080828437600092019190915250506040805180820182529295600160a060020a0385351695909490936060820193509160209091019060029083908390808284376000920191909152509194506107d29350505050565b34801561031157600080fd5b506101b16004803603604081101561032857600080fd5b50803590602001356107eb565b34801561034157600080fd5b506101b16004803603602081101561035857600080fd5b50356107f7565b34801561036b57600080fd5b506103cd600480360360a081101561038257600080fd5b60408051808201825283359392830192916060830191906020840190600290839083908082843760009201919091525091945050813592505060200135600160a060020a0316610802565b604080519115158252519081900360200190f35b3480156103ed57600080fd5b5061041d6004803603608081101561040457600080fd5b5080359060208101359060408101359060600135610819565b60408051938452602084019290925282820152519081900360600190f35b34801561044757600080fd5b506101b16004803603602081101561045e57600080fd5b503561083a565b34801561047157600080fd5b506103cd600480360360a081101561048857600080fd5b60408051808201825291830192918183019183906002908390839080828437600092019190915250506040805180820182529295843595909490936060820193509160209091019060029083908390808284376000920191909152509194506108459350505050565b3480156104fd57600080fd5b506105eb60048036036101a081101561051557600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525050604080518082018252929594938181019392509060029083908390808284376000920191909152505060408051808201825292958435956020860135958381013595600160a060020a0360608301351695509293919260c08201929091608001906002908390839080828437600092019190915250506040805180820182529295949381810193925090600290839083908082843760009201919091525091945050903591506108529050565b005b3480156105f957600080fd5b506101b16004803603602081101561061057600080fd5b81019060208101813564010000000081111561062b57600080fd5b82018360208201111561063d57600080fd5b8035906020019184600183028401116401000000008311171561065f57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061086e945050505050565b3480156106ac57600080fd5b5061014c60048036036101608110156106c457600080fd5b604080518082018252833593928301929160608301919060208401906002908390839080828437600092019190915250506040805180820182529295949381810193925090600290839083908082843760009201919091525050604080518082018252929584359590949093606082019350916020909101906002908390839080828437600092019190915250506040805180820182529295949381810193925090600290839083908082843760009201919091525091945050903591506109be9050565b6107916115f9565b61079c8484846109e1565b949350505050565b60006107af82610aa3565b90505b919050565b6107bf6115f9565b6107c98383610ad8565b90505b92915050565b60006107e18686868686610be7565b9695505050505050565b60006107c98383610cf4565b60006107af82610da0565b600061081085858585610dcc565b95945050505050565b600080600061082a87878787610f2b565b9250925092509450945094915050565b60006107af8261100b565b600061079c84848461102f565b610863898989898989898989611167565b505050505050505050565b80516000906101a0146108cb576040805160e560020a62461bcd02815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b6108d36115f9565b6108db6115f9565b6108e3611614565b60006108ed6115f9565b6108f56115f9565b6000888060200190516101a081101561090d57600080fd5b5060e0810151610180820151919850604089019750608089019650945061010088019350610140880192509050610960878787600060200201518860016020020151896002602002015189898989611167565b856040516020018082600260200280838360005b8381101561098c578181015183820152602001610974565b505050509050019150506040516020818303038152906040528051906020012060019004975050505050505050919050565b6109c66115f9565b6109d5888888888888886113e2565b98975050505050505050565b6109e96115f9565b835160208086015185519186015160009384938493610a0a93909190610f2b565b919450925090506401000003d019858209600114610a72576040805160e560020a62461bcd02815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a00000000000000604482015290519081900360640190fd5b60408051808201909152806401000003d01987860981526020016401000003d0198785099052979650505050505050565b805b6401000003d01981106107b257604080516020808201939093528151808203840181529082019091528051910120610aa5565b610ae06115f9565b610b3a83836040516020018083600260200280838360005b83811015610b10578181015183820152602001610af8565b50505050919091019283525050604080518083038152602092830190915280519101209050610aa3565b8152610b55610b508260005b602002015161100b565b610da0565b60208201525b610b66816000610b46565b60208201516401000003d01990800914610bc0578051604080516020818101939093528151808203840181529082019091528051910120610ba690610aa3565b8152610bb6610b50826000610b46565b6020820152610b5b565b602081015160029006600114156107cc576020810180516401000003d01903905292915050565b600085858584866040516020018086600260200280838360005b83811015610c19578181015183820152602001610c01565b5050505090500185600260200280838360005b83811015610c44578181015183820152602001610c2c565b5050505090500184600260200280838360005b83811015610c6f578181015183820152602001610c57565b5050505090500183600260200280838360005b83811015610c9a578181015183820152602001610c82565b5050505090500182600160a060020a0316600160a060020a03166c01000000000000000000000000028152601401955050505050506040516020818303038152906040528051906020012060019004905095945050505050565b600080610cff611633565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152610d31611652565b60208160c0846005600019fa9250821515610d96576040805160e560020a62461bcd02815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015290519081900360640190fd5b5195945050505050565b60006107af827f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c610cf4565b6000600160a060020a0382161515610e2e576040805160e560020a62461bcd02815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015290519081900360640190fd5b835160009070014551231950b75fc4402da1732fc9bebe199085900970014551231950b75fc4402da1732fc9bebe19039050600060028660016020020151811515610e7557fe5b0615610e8257601c610e85565b601b5b865190915060009070014551231950b75fc4402da1732fc9bebe199089098751604080516000808252602082810180855289905260ff8816838501526060830194909452608082018590529151939450909260019260a0808401939192601f1981019281900390910190855afa158015610f03573d6000803e3d6000fd5b5050604051601f190151600160a060020a039081169088161495505050505050949350505050565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000610f6b83838585611568565b9098509050610f7c88828e8861158c565b9098509050610f8d88828c8761158c565b90985090506000610fa08d878b8561158c565b9098509050610fb188828686611568565b9098509050610fc288828e8961158c565b9098509050818114610ff7576401000003d019818a0998506401000003d01982890997506401000003d0198183099650610ffb565b8196505b5050505050509450945094915050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b600082151561103d57600080fd5b835160009070014551231950b75fc4402da1732fc9bebe199085096020860151909150600090600116151561107357601b611076565b601c5b9050836040516020018082600260200280838360005b838110156110a457818101518382015260200161108c565b505050509050019150506040516020818303038152906040528051906020012060019004600160a060020a031660016000600102838960006002811015156110e857fe5b60200201516001028660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611148573d6000803e3d6000fd5b50505060206040510351600160a060020a031614925050509392505050565b611170896115d5565b15156111c6576040805160e560020a62461bcd02815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015290519081900360640190fd5b6111cf886115d5565b1515611225576040805160e560020a62461bcd02815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015290519081900360640190fd5b61122e836115d5565b1515611284576040805160e560020a62461bcd02815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015290519081900360640190fd5b61128d826115d5565b15156112e3576040805160e560020a62461bcd02815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015290519081900360640190fd5b6112ef878a8887610dcc565b1515611345576040805160e560020a62461bcd02815260206004820152601a60248201527f6164647228632a706b2b732a6729e289a05f755769746e657373000000000000604482015290519081900360640190fd5b61134d6115f9565b6113578a87610ad8565b90506113616115f9565b611370898b878b8689896113e2565b905061137f828c8c8985610be7565b89146113d5576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015290519081900360640190fd5b5050505050505050505050565b6113ea6115f9565b825186516401000003d01991900306151561144f576040805160e560020a62461bcd02815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015290519081900360640190fd5b61145a87898861102f565b15156114d6576040805160e560020a62461bcd02815260206004820152602160248201527f4669727374206d756c7469706c69636174696f6e20636865636b206661696c6560448201527f6400000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b6114e184868561102f565b151561155d576040805160e560020a62461bcd02815260206004820152602260248201527f5365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c60448201527f6564000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b6109d58684846109e1565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b60208101516000906401000003d0199080096115f2836000610b46565b1492915050565b60408051808201825290600290829080388339509192915050565b6060604051908101604052806003906020820280388339509192915050565b60c0604051908101604052806006906020820280388339509192915050565b602060405190810160405280600190602082028038833950919291505056fea165627a7a723058206de5d4399dcb23c1aca39c9812cb4891a047a5e7e84f6234978629382a6ea5b30029`

// DeployVRFTestHelper deploys a new Ethereum contract, binding an instance of VRFTestHelper to it.
func DeployVRFTestHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFTestHelper, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFTestHelperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFTestHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFTestHelper{VRFTestHelperCaller: VRFTestHelperCaller{contract: contract}, VRFTestHelperTransactor: VRFTestHelperTransactor{contract: contract}, VRFTestHelperFilterer: VRFTestHelperFilterer{contract: contract}}, nil
}

// VRFTestHelper is an auto generated Go binding around an Ethereum contract.
type VRFTestHelper struct {
	VRFTestHelperCaller     // Read-only binding to the contract
	VRFTestHelperTransactor // Write-only binding to the contract
	VRFTestHelperFilterer   // Log filterer for contract events
}

// VRFTestHelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFTestHelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestHelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFTestHelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestHelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFTestHelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestHelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFTestHelperSession struct {
	Contract     *VRFTestHelper    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFTestHelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFTestHelperCallerSession struct {
	Contract *VRFTestHelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VRFTestHelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFTestHelperTransactorSession struct {
	Contract     *VRFTestHelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VRFTestHelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFTestHelperRaw struct {
	Contract *VRFTestHelper // Generic contract binding to access the raw methods on
}

// VRFTestHelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFTestHelperCallerRaw struct {
	Contract *VRFTestHelperCaller // Generic read-only contract binding to access the raw methods on
}

// VRFTestHelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFTestHelperTransactorRaw struct {
	Contract *VRFTestHelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFTestHelper creates a new instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelper(address common.Address, backend bind.ContractBackend) (*VRFTestHelper, error) {
	contract, err := bindVRFTestHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelper{VRFTestHelperCaller: VRFTestHelperCaller{contract: contract}, VRFTestHelperTransactor: VRFTestHelperTransactor{contract: contract}, VRFTestHelperFilterer: VRFTestHelperFilterer{contract: contract}}, nil
}

// NewVRFTestHelperCaller creates a new read-only instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelperCaller(address common.Address, caller bind.ContractCaller) (*VRFTestHelperCaller, error) {
	contract, err := bindVRFTestHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelperCaller{contract: contract}, nil
}

// NewVRFTestHelperTransactor creates a new write-only instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFTestHelperTransactor, error) {
	contract, err := bindVRFTestHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelperTransactor{contract: contract}, nil
}

// NewVRFTestHelperFilterer creates a new log filterer instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFTestHelperFilterer, error) {
	contract, err := bindVRFTestHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelperFilterer{contract: contract}, nil
}

// bindVRFTestHelper binds a generic wrapper to an already deployed contract.
func bindVRFTestHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFTestHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFTestHelper *VRFTestHelperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFTestHelper.Contract.VRFTestHelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFTestHelper *VRFTestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.VRFTestHelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFTestHelper *VRFTestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.VRFTestHelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFTestHelper *VRFTestHelperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFTestHelper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFTestHelper *VRFTestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFTestHelper *VRFTestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.contract.Transact(opts, method, params...)
}

// AffineECAdd is a free data retrieval call binding the contract method 0x244f896d.
//
// Solidity: function affineECAdd_(uint256[2] p1, uint256[2] p2, uint256 invZ) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCaller) AffineECAdd(opts *bind.CallOpts, p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "affineECAdd_", p1, p2, invZ)
	return *ret0, err
}

// AffineECAdd is a free data retrieval call binding the contract method 0x244f896d.
//
// Solidity: function affineECAdd_(uint256[2] p1, uint256[2] p2, uint256 invZ) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperSession) AffineECAdd(p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.AffineECAdd(&_VRFTestHelper.CallOpts, p1, p2, invZ)
}

// AffineECAdd is a free data retrieval call binding the contract method 0x244f896d.
//
// Solidity: function affineECAdd_(uint256[2] p1, uint256[2] p2, uint256 invZ) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCallerSession) AffineECAdd(p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.AffineECAdd(&_VRFTestHelper.CallOpts, p1, p2, invZ)
}

// BigModExp is a free data retrieval call binding the contract method 0x5de60042.
//
// Solidity: function bigModExp_(uint256 base, uint256 exponent) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) BigModExp(opts *bind.CallOpts, base *big.Int, exponent *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "bigModExp_", base, exponent)
	return *ret0, err
}

// BigModExp is a free data retrieval call binding the contract method 0x5de60042.
//
// Solidity: function bigModExp_(uint256 base, uint256 exponent) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) BigModExp(base *big.Int, exponent *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.BigModExp(&_VRFTestHelper.CallOpts, base, exponent)
}

// BigModExp is a free data retrieval call binding the contract method 0x5de60042.
//
// Solidity: function bigModExp_(uint256 base, uint256 exponent) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) BigModExp(base *big.Int, exponent *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.BigModExp(&_VRFTestHelper.CallOpts, base, exponent)
}

// EcmulVerify is a free data retrieval call binding the contract method 0xaa7b2fbb.
//
// Solidity: function ecmulVerify_(uint256[2] x, uint256 scalar, uint256[2] q) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCaller) EcmulVerify(opts *bind.CallOpts, x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "ecmulVerify_", x, scalar, q)
	return *ret0, err
}

// EcmulVerify is a free data retrieval call binding the contract method 0xaa7b2fbb.
//
// Solidity: function ecmulVerify_(uint256[2] x, uint256 scalar, uint256[2] q) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperSession) EcmulVerify(x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	return _VRFTestHelper.Contract.EcmulVerify(&_VRFTestHelper.CallOpts, x, scalar, q)
}

// EcmulVerify is a free data retrieval call binding the contract method 0xaa7b2fbb.
//
// Solidity: function ecmulVerify_(uint256[2] x, uint256 scalar, uint256[2] q) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCallerSession) EcmulVerify(x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	return _VRFTestHelper.Contract.EcmulVerify(&_VRFTestHelper.CallOpts, x, scalar, q)
}

// HashToCurve is a free data retrieval call binding the contract method 0x35452450.
//
// Solidity: function hashToCurve_(uint256[2] pk, uint256 x) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCaller) HashToCurve(opts *bind.CallOpts, pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "hashToCurve_", pk, x)
	return *ret0, err
}

// HashToCurve is a free data retrieval call binding the contract method 0x35452450.
//
// Solidity: function hashToCurve_(uint256[2] pk, uint256 x) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperSession) HashToCurve(pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.HashToCurve(&_VRFTestHelper.CallOpts, pk, x)
}

// HashToCurve is a free data retrieval call binding the contract method 0x35452450.
//
// Solidity: function hashToCurve_(uint256[2] pk, uint256 x) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCallerSession) HashToCurve(pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.HashToCurve(&_VRFTestHelper.CallOpts, pk, x)
}

// LinearCombination is a free data retrieval call binding the contract method 0xfe54f2a2.
//
// Solidity: function linearCombination_(uint256 c, uint256[2] p1, uint256[2] cp1Witness, uint256 s, uint256[2] p2, uint256[2] sp2Witness, uint256 zInv) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCaller) LinearCombination(opts *bind.CallOpts, c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "linearCombination_", c, p1, cp1Witness, s, p2, sp2Witness, zInv)
	return *ret0, err
}

// LinearCombination is a free data retrieval call binding the contract method 0xfe54f2a2.
//
// Solidity: function linearCombination_(uint256 c, uint256[2] p1, uint256[2] cp1Witness, uint256 s, uint256[2] p2, uint256[2] sp2Witness, uint256 zInv) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperSession) LinearCombination(c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.LinearCombination(&_VRFTestHelper.CallOpts, c, p1, cp1Witness, s, p2, sp2Witness, zInv)
}

// LinearCombination is a free data retrieval call binding the contract method 0xfe54f2a2.
//
// Solidity: function linearCombination_(uint256 c, uint256[2] p1, uint256[2] cp1Witness, uint256 s, uint256[2] p2, uint256[2] sp2Witness, uint256 zInv) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCallerSession) LinearCombination(c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.LinearCombination(&_VRFTestHelper.CallOpts, c, p1, cp1Witness, s, p2, sp2Witness, zInv)
}

// ProjectiveECAdd is a free data retrieval call binding the contract method 0x95e6ee92.
//
// Solidity: function projectiveECAdd_(uint256 px, uint256 py, uint256 qx, uint256 qy) constant returns(uint256, uint256, uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ProjectiveECAdd(opts *bind.CallOpts, px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _VRFTestHelper.contract.Call(opts, out, "projectiveECAdd_", px, py, qx, qy)
	return *ret0, *ret1, *ret2, err
}

// ProjectiveECAdd is a free data retrieval call binding the contract method 0x95e6ee92.
//
// Solidity: function projectiveECAdd_(uint256 px, uint256 py, uint256 qx, uint256 qy) constant returns(uint256, uint256, uint256)
func (_VRFTestHelper *VRFTestHelperSession) ProjectiveECAdd(px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFTestHelper.Contract.ProjectiveECAdd(&_VRFTestHelper.CallOpts, px, py, qx, qy)
}

// ProjectiveECAdd is a free data retrieval call binding the contract method 0x95e6ee92.
//
// Solidity: function projectiveECAdd_(uint256 px, uint256 py, uint256 qx, uint256 qy) constant returns(uint256, uint256, uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ProjectiveECAdd(px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFTestHelper.Contract.ProjectiveECAdd(&_VRFTestHelper.CallOpts, px, py, qx, qy)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperCaller) RandomValueFromVRFProof(opts *bind.CallOpts, proof []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "randomValueFromVRFProof", proof)
	return *ret0, err
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.RandomValueFromVRFProof(&_VRFTestHelper.CallOpts, proof)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperCallerSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.RandomValueFromVRFProof(&_VRFTestHelper.CallOpts, proof)
}

// ScalarFromCurve is a free data retrieval call binding the contract method 0x525413cf.
//
// Solidity: function scalarFromCurve_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ScalarFromCurve(opts *bind.CallOpts, hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "scalarFromCurve_", hash, pk, gamma, uWitness, v)
	return *ret0, err
}

// ScalarFromCurve is a free data retrieval call binding the contract method 0x525413cf.
//
// Solidity: function scalarFromCurve_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) ScalarFromCurve(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ScalarFromCurve(&_VRFTestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

// ScalarFromCurve is a free data retrieval call binding the contract method 0x525413cf.
//
// Solidity: function scalarFromCurve_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ScalarFromCurve(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ScalarFromCurve(&_VRFTestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

// SquareRoot is a free data retrieval call binding the contract method 0x8af046ea.
//
// Solidity: function squareRoot_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) SquareRoot(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "squareRoot_", x)
	return *ret0, err
}

// SquareRoot is a free data retrieval call binding the contract method 0x8af046ea.
//
// Solidity: function squareRoot_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) SquareRoot(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.SquareRoot(&_VRFTestHelper.CallOpts, x)
}

// SquareRoot is a free data retrieval call binding the contract method 0x8af046ea.
//
// Solidity: function squareRoot_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) SquareRoot(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.SquareRoot(&_VRFTestHelper.CallOpts, x)
}

// VerifyLinearCombinationWithGenerator is a free data retrieval call binding the contract method 0x91d5f691.
//
// Solidity: function verifyLinearCombinationWithGenerator_(uint256 c, uint256[2] p, uint256 s, address lcWitness) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCaller) VerifyLinearCombinationWithGenerator(opts *bind.CallOpts, c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "verifyLinearCombinationWithGenerator_", c, p, s, lcWitness)
	return *ret0, err
}

// VerifyLinearCombinationWithGenerator is a free data retrieval call binding the contract method 0x91d5f691.
//
// Solidity: function verifyLinearCombinationWithGenerator_(uint256 c, uint256[2] p, uint256 s, address lcWitness) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperSession) VerifyLinearCombinationWithGenerator(c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	return _VRFTestHelper.Contract.VerifyLinearCombinationWithGenerator(&_VRFTestHelper.CallOpts, c, p, s, lcWitness)
}

// VerifyLinearCombinationWithGenerator is a free data retrieval call binding the contract method 0x91d5f691.
//
// Solidity: function verifyLinearCombinationWithGenerator_(uint256 c, uint256[2] p, uint256 s, address lcWitness) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCallerSession) VerifyLinearCombinationWithGenerator(c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	return _VRFTestHelper.Contract.VerifyLinearCombinationWithGenerator(&_VRFTestHelper.CallOpts, c, p, s, lcWitness)
}

// VerifyVRFProof is a free data retrieval call binding the contract method 0xef3b10ec.
//
// Solidity: function verifyVRFProof_(uint256[2] pk, uint256[2] gamma, uint256 c, uint256 s, uint256 seed, address uWitness, uint256[2] cGammaWitness, uint256[2] sHashWitness, uint256 zInv) constant returns()
func (_VRFTestHelper *VRFTestHelperCaller) VerifyVRFProof(opts *bind.CallOpts, pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	var ()
	out := &[]interface{}{}
	err := _VRFTestHelper.contract.Call(opts, out, "verifyVRFProof_", pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
	return err
}

// VerifyVRFProof is a free data retrieval call binding the contract method 0xef3b10ec.
//
// Solidity: function verifyVRFProof_(uint256[2] pk, uint256[2] gamma, uint256 c, uint256 s, uint256 seed, address uWitness, uint256[2] cGammaWitness, uint256[2] sHashWitness, uint256 zInv) constant returns()
func (_VRFTestHelper *VRFTestHelperSession) VerifyVRFProof(pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	return _VRFTestHelper.Contract.VerifyVRFProof(&_VRFTestHelper.CallOpts, pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
}

// VerifyVRFProof is a free data retrieval call binding the contract method 0xef3b10ec.
//
// Solidity: function verifyVRFProof_(uint256[2] pk, uint256[2] gamma, uint256 c, uint256 s, uint256 seed, address uWitness, uint256[2] cGammaWitness, uint256[2] sHashWitness, uint256 zInv) constant returns()
func (_VRFTestHelper *VRFTestHelperCallerSession) VerifyVRFProof(pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	return _VRFTestHelper.Contract.VerifyVRFProof(&_VRFTestHelper.CallOpts, pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
}

// YSquared is a free data retrieval call binding the contract method 0x9d6f0337.
//
// Solidity: function ySquared_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) YSquared(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "ySquared_", x)
	return *ret0, err
}

// YSquared is a free data retrieval call binding the contract method 0x9d6f0337.
//
// Solidity: function ySquared_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) YSquared(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.YSquared(&_VRFTestHelper.CallOpts, x)
}

// YSquared is a free data retrieval call binding the contract method 0x9d6f0337.
//
// Solidity: function ySquared_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) YSquared(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.YSquared(&_VRFTestHelper.CallOpts, x)
}

// ZqHash is a free data retrieval call binding the contract method 0x24d72ea9.
//
// Solidity: function zqHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ZqHash(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "zqHash_", x)
	return *ret0, err
}

// ZqHash is a free data retrieval call binding the contract method 0x24d72ea9.
//
// Solidity: function zqHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) ZqHash(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ZqHash(&_VRFTestHelper.CallOpts, x)
}

// ZqHash is a free data retrieval call binding the contract method 0x24d72ea9.
//
// Solidity: function zqHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ZqHash(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ZqHash(&_VRFTestHelper.CallOpts, x)
}
