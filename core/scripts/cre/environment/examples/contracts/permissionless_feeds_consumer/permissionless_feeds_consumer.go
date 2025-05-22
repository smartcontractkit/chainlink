// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package permissionless_feeds_consumer

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

// PermissionlessFeedsConsumerMetaData contains all meta data concerning the PermissionlessFeedsConsumer contract.
var PermissionlessFeedsConsumerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint224\",\"name\":\"price\",\"type\":\"uint224\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"FeedReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5033805f8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000085576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200007c906200029e565b60405180910390fd5b815f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146200010a5762000109816200011360201b60201c565b5b5050506200032c565b3373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000184576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200017b906200030c565b60405180910390fd5b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff165f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b5f82825260208201905092915050565b7f43616e6e6f7420736574206f776e657220746f207a65726f00000000000000005f82015250565b5f6200028660188362000240565b9150620002938262000250565b602082019050919050565b5f6020820190508181035f830152620002b78162000278565b9050919050565b7f43616e6e6f74207472616e7366657220746f2073656c660000000000000000005f82015250565b5f620002f460178362000240565b91506200030182620002be565b602082019050919050565b5f6020820190508181035f8301526200032581620002e6565b9050919050565b610f3d806200033a5f395ff3fe608060405234801561000f575f80fd5b5060043610610060575f3560e01c806301ffc9a71461006457806331d98b3f1461009457806379ba5097146100c5578063805f2132146100cf5780638da5cb5b146100eb578063f2fde38b14610109575b5f80fd5b61007e60048036038101906100799190610887565b610125565b60405161008b91906108cc565b60405180910390f35b6100ae60048036038101906100a99190610918565b6101f6565b6040516100bc929190610997565b60405180910390f35b6100cd6102bb565b005b6100e960048036038101906100e49190610a1f565b61044a565b005b6100f361062d565b6040516101009190610adc565b60405180910390f35b610123600480360381019061011e9190610b1f565b610654565b005b5f7f805f2132000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806101ef57507f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b9050919050565b5f805f60025f8581526020019081526020015f206040518060400160405290815f82015f9054906101000a90047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020015f8201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff16815250509050805f015181602001519250925050915091565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461034a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161034190610ba4565b60405180910390fd5b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050335f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505f60015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a350565b5f828281019061045a9190610dc3565b90505f5b815181101561062557604051806040016040528083838151811061048557610484610e0a565b5b6020026020010151604001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020018383815181106104c7576104c6610e0a565b5b60200260200101516020015163ffffffff1681525060025f8484815181106104f2576104f1610e0a565b5b60200260200101515f015181526020019081526020015f205f820151815f015f6101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1602179055506020820151815f01601c6101000a81548163ffffffff021916908363ffffffff16021790555090505081818151811061059a57610599610e0a565b5b60200260200101515f01517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d8383815181106105d9576105d8610e0a565b5b6020026020010151604001518484815181106105f8576105f7610e0a565b5b602002602001015160200151604051610612929190610997565b60405180910390a280600101905061045e565b505050505050565b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61065c610668565b610665816106f7565b50565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106f5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106ec90610e81565b60405180910390fd5b565b3373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610765576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161075c90610ee9565b60405180910390fd5b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff165f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b61086681610832565b8114610870575f80fd5b50565b5f813590506108818161085d565b92915050565b5f6020828403121561089c5761089b61082a565b5b5f6108a984828501610873565b91505092915050565b5f8115159050919050565b6108c6816108b2565b82525050565b5f6020820190506108df5f8301846108bd565b92915050565b5f819050919050565b6108f7816108e5565b8114610901575f80fd5b50565b5f81359050610912816108ee565b92915050565b5f6020828403121561092d5761092c61082a565b5b5f61093a84828501610904565b91505092915050565b5f7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff82169050919050565b61097381610943565b82525050565b5f63ffffffff82169050919050565b61099181610979565b82525050565b5f6040820190506109aa5f83018561096a565b6109b76020830184610988565b9392505050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f8401126109df576109de6109be565b5b8235905067ffffffffffffffff8111156109fc576109fb6109c2565b5b602083019150836001820283011115610a1857610a176109c6565b5b9250929050565b5f805f8060408587031215610a3757610a3661082a565b5b5f85013567ffffffffffffffff811115610a5457610a5361082e565b5b610a60878288016109ca565b9450945050602085013567ffffffffffffffff811115610a8357610a8261082e565b5b610a8f878288016109ca565b925092505092959194509250565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610ac682610a9d565b9050919050565b610ad681610abc565b82525050565b5f602082019050610aef5f830184610acd565b92915050565b610afe81610abc565b8114610b08575f80fd5b50565b5f81359050610b1981610af5565b92915050565b5f60208284031215610b3457610b3361082a565b5b5f610b4184828501610b0b565b91505092915050565b5f82825260208201905092915050565b7f4d7573742062652070726f706f736564206f776e6572000000000000000000005f82015250565b5f610b8e601683610b4a565b9150610b9982610b5a565b602082019050919050565b5f6020820190508181035f830152610bbb81610b82565b9050919050565b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610c0882610bc2565b810181811067ffffffffffffffff82111715610c2757610c26610bd2565b5b80604052505050565b5f610c39610821565b9050610c458282610bff565b919050565b5f67ffffffffffffffff821115610c6457610c63610bd2565b5b602082029050602081019050919050565b5f80fd5b610c8281610979565b8114610c8c575f80fd5b50565b5f81359050610c9d81610c79565b92915050565b610cac81610943565b8114610cb6575f80fd5b50565b5f81359050610cc781610ca3565b92915050565b5f60608284031215610ce257610ce1610c75565b5b610cec6060610c30565b90505f610cfb84828501610904565b5f830152506020610d0e84828501610c8f565b6020830152506040610d2284828501610cb9565b60408301525092915050565b5f610d40610d3b84610c4a565b610c30565b90508083825260208201905060608402830185811115610d6357610d626109c6565b5b835b81811015610d8c5780610d788882610ccd565b845260208401935050606081019050610d65565b5050509392505050565b5f82601f830112610daa57610da96109be565b5b8135610dba848260208601610d2e565b91505092915050565b5f60208284031215610dd857610dd761082a565b5b5f82013567ffffffffffffffff811115610df557610df461082e565b5b610e0184828501610d96565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000005f82015250565b5f610e6b601683610b4a565b9150610e7682610e37565b602082019050919050565b5f6020820190508181035f830152610e9881610e5f565b9050919050565b7f43616e6e6f74207472616e7366657220746f2073656c660000000000000000005f82015250565b5f610ed3601783610b4a565b9150610ede82610e9f565b602082019050919050565b5f6020820190508181035f830152610f0081610ec7565b905091905056fea2646970667358221220f9ba75e24e4d273d983e4d5101485cdedbaa9f6f87a2d973e79b08df6f56ca4e64736f6c63430008180033",
}

// PermissionlessFeedsConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use PermissionlessFeedsConsumerMetaData.ABI instead.
var PermissionlessFeedsConsumerABI = PermissionlessFeedsConsumerMetaData.ABI

// PermissionlessFeedsConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PermissionlessFeedsConsumerMetaData.Bin instead.
var PermissionlessFeedsConsumerBin = PermissionlessFeedsConsumerMetaData.Bin

// DeployPermissionlessFeedsConsumer deploys a new Ethereum contract, binding an instance of PermissionlessFeedsConsumer to it.
func DeployPermissionlessFeedsConsumer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PermissionlessFeedsConsumer, error) {
	parsed, err := PermissionlessFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PermissionlessFeedsConsumerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PermissionlessFeedsConsumer{PermissionlessFeedsConsumerCaller: PermissionlessFeedsConsumerCaller{contract: contract}, PermissionlessFeedsConsumerTransactor: PermissionlessFeedsConsumerTransactor{contract: contract}, PermissionlessFeedsConsumerFilterer: PermissionlessFeedsConsumerFilterer{contract: contract}}, nil
}

// PermissionlessFeedsConsumer is an auto generated Go binding around an Ethereum contract.
type PermissionlessFeedsConsumer struct {
	PermissionlessFeedsConsumerCaller     // Read-only binding to the contract
	PermissionlessFeedsConsumerTransactor // Write-only binding to the contract
	PermissionlessFeedsConsumerFilterer   // Log filterer for contract events
}

// PermissionlessFeedsConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PermissionlessFeedsConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PermissionlessFeedsConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PermissionlessFeedsConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PermissionlessFeedsConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PermissionlessFeedsConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PermissionlessFeedsConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PermissionlessFeedsConsumerSession struct {
	Contract     *PermissionlessFeedsConsumer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// PermissionlessFeedsConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PermissionlessFeedsConsumerCallerSession struct {
	Contract *PermissionlessFeedsConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// PermissionlessFeedsConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PermissionlessFeedsConsumerTransactorSession struct {
	Contract     *PermissionlessFeedsConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// PermissionlessFeedsConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PermissionlessFeedsConsumerRaw struct {
	Contract *PermissionlessFeedsConsumer // Generic contract binding to access the raw methods on
}

// PermissionlessFeedsConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PermissionlessFeedsConsumerCallerRaw struct {
	Contract *PermissionlessFeedsConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// PermissionlessFeedsConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PermissionlessFeedsConsumerTransactorRaw struct {
	Contract *PermissionlessFeedsConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPermissionlessFeedsConsumer creates a new instance of PermissionlessFeedsConsumer, bound to a specific deployed contract.
func NewPermissionlessFeedsConsumer(address common.Address, backend bind.ContractBackend) (*PermissionlessFeedsConsumer, error) {
	contract, err := bindPermissionlessFeedsConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumer{PermissionlessFeedsConsumerCaller: PermissionlessFeedsConsumerCaller{contract: contract}, PermissionlessFeedsConsumerTransactor: PermissionlessFeedsConsumerTransactor{contract: contract}, PermissionlessFeedsConsumerFilterer: PermissionlessFeedsConsumerFilterer{contract: contract}}, nil
}

// NewPermissionlessFeedsConsumerCaller creates a new read-only instance of PermissionlessFeedsConsumer, bound to a specific deployed contract.
func NewPermissionlessFeedsConsumerCaller(address common.Address, caller bind.ContractCaller) (*PermissionlessFeedsConsumerCaller, error) {
	contract, err := bindPermissionlessFeedsConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumerCaller{contract: contract}, nil
}

// NewPermissionlessFeedsConsumerTransactor creates a new write-only instance of PermissionlessFeedsConsumer, bound to a specific deployed contract.
func NewPermissionlessFeedsConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*PermissionlessFeedsConsumerTransactor, error) {
	contract, err := bindPermissionlessFeedsConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumerTransactor{contract: contract}, nil
}

// NewPermissionlessFeedsConsumerFilterer creates a new log filterer instance of PermissionlessFeedsConsumer, bound to a specific deployed contract.
func NewPermissionlessFeedsConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*PermissionlessFeedsConsumerFilterer, error) {
	contract, err := bindPermissionlessFeedsConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumerFilterer{contract: contract}, nil
}

// bindPermissionlessFeedsConsumer binds a generic wrapper to an already deployed contract.
func bindPermissionlessFeedsConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PermissionlessFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PermissionlessFeedsConsumer.Contract.PermissionlessFeedsConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.PermissionlessFeedsConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.PermissionlessFeedsConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PermissionlessFeedsConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.contract.Transact(opts, method, params...)
}

// GetPrice is a free data retrieval call binding the contract method 0x31d98b3f.
//
// Solidity: function getPrice(bytes32 feedId) view returns(uint224, uint32)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCaller) GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error) {
	var out []interface{}
	err := _PermissionlessFeedsConsumer.contract.Call(opts, &out, "getPrice", feedId)

	if err != nil {
		return *new(*big.Int), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

// GetPrice is a free data retrieval call binding the contract method 0x31d98b3f.
//
// Solidity: function getPrice(bytes32 feedId) view returns(uint224, uint32)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _PermissionlessFeedsConsumer.Contract.GetPrice(&_PermissionlessFeedsConsumer.CallOpts, feedId)
}

// GetPrice is a free data retrieval call binding the contract method 0x31d98b3f.
//
// Solidity: function getPrice(bytes32 feedId) view returns(uint224, uint32)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCallerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _PermissionlessFeedsConsumer.Contract.GetPrice(&_PermissionlessFeedsConsumer.CallOpts, feedId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PermissionlessFeedsConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerSession) Owner() (common.Address, error) {
	return _PermissionlessFeedsConsumer.Contract.Owner(&_PermissionlessFeedsConsumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCallerSession) Owner() (common.Address, error) {
	return _PermissionlessFeedsConsumer.Contract.Owner(&_PermissionlessFeedsConsumer.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PermissionlessFeedsConsumer.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PermissionlessFeedsConsumer.Contract.SupportsInterface(&_PermissionlessFeedsConsumer.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PermissionlessFeedsConsumer.Contract.SupportsInterface(&_PermissionlessFeedsConsumer.CallOpts, interfaceId)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.AcceptOwnership(&_PermissionlessFeedsConsumer.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.AcceptOwnership(&_PermissionlessFeedsConsumer.TransactOpts)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes metadata, bytes rawReport) returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.contract.Transact(opts, "onReport", metadata, rawReport)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes metadata, bytes rawReport) returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.OnReport(&_PermissionlessFeedsConsumer.TransactOpts, metadata, rawReport)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes metadata, bytes rawReport) returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactorSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.OnReport(&_PermissionlessFeedsConsumer.TransactOpts, metadata, rawReport)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.TransferOwnership(&_PermissionlessFeedsConsumer.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PermissionlessFeedsConsumer.Contract.TransferOwnership(&_PermissionlessFeedsConsumer.TransactOpts, to)
}

// PermissionlessFeedsConsumerFeedReceivedIterator is returned from FilterFeedReceived and is used to iterate over the raw logs and unpacked data for FeedReceived events raised by the PermissionlessFeedsConsumer contract.
type PermissionlessFeedsConsumerFeedReceivedIterator struct {
	Event *PermissionlessFeedsConsumerFeedReceived // Event containing the contract specifics and raw log

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
func (it *PermissionlessFeedsConsumerFeedReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionlessFeedsConsumerFeedReceived)
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
		it.Event = new(PermissionlessFeedsConsumerFeedReceived)
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
func (it *PermissionlessFeedsConsumerFeedReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PermissionlessFeedsConsumerFeedReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PermissionlessFeedsConsumerFeedReceived represents a FeedReceived event raised by the PermissionlessFeedsConsumer contract.
type PermissionlessFeedsConsumerFeedReceived struct {
	FeedId    [32]byte
	Price     *big.Int
	Timestamp uint32
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeedReceived is a free log retrieval operation binding the contract event 0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d.
//
// Solidity: event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*PermissionlessFeedsConsumerFeedReceivedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _PermissionlessFeedsConsumer.contract.FilterLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumerFeedReceivedIterator{contract: _PermissionlessFeedsConsumer.contract, event: "FeedReceived", logs: logs, sub: sub}, nil
}

// WatchFeedReceived is a free log subscription operation binding the contract event 0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d.
//
// Solidity: event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *PermissionlessFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _PermissionlessFeedsConsumer.contract.WatchLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PermissionlessFeedsConsumerFeedReceived)
				if err := _PermissionlessFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
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

// ParseFeedReceived is a log parse operation binding the contract event 0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d.
//
// Solidity: event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) ParseFeedReceived(log types.Log) (*PermissionlessFeedsConsumerFeedReceived, error) {
	event := new(PermissionlessFeedsConsumerFeedReceived)
	if err := _PermissionlessFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PermissionlessFeedsConsumerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the PermissionlessFeedsConsumer contract.
type PermissionlessFeedsConsumerOwnershipTransferRequestedIterator struct {
	Event *PermissionlessFeedsConsumerOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *PermissionlessFeedsConsumerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionlessFeedsConsumerOwnershipTransferRequested)
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
		it.Event = new(PermissionlessFeedsConsumerOwnershipTransferRequested)
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
func (it *PermissionlessFeedsConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PermissionlessFeedsConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PermissionlessFeedsConsumerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the PermissionlessFeedsConsumer contract.
type PermissionlessFeedsConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionlessFeedsConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumerOwnershipTransferRequestedIterator{contract: _PermissionlessFeedsConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PermissionlessFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PermissionlessFeedsConsumerOwnershipTransferRequested)
				if err := _PermissionlessFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*PermissionlessFeedsConsumerOwnershipTransferRequested, error) {
	event := new(PermissionlessFeedsConsumerOwnershipTransferRequested)
	if err := _PermissionlessFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PermissionlessFeedsConsumerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PermissionlessFeedsConsumer contract.
type PermissionlessFeedsConsumerOwnershipTransferredIterator struct {
	Event *PermissionlessFeedsConsumerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PermissionlessFeedsConsumerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionlessFeedsConsumerOwnershipTransferred)
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
		it.Event = new(PermissionlessFeedsConsumerOwnershipTransferred)
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
func (it *PermissionlessFeedsConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PermissionlessFeedsConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PermissionlessFeedsConsumerOwnershipTransferred represents a OwnershipTransferred event raised by the PermissionlessFeedsConsumer contract.
type PermissionlessFeedsConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionlessFeedsConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PermissionlessFeedsConsumerOwnershipTransferredIterator{contract: _PermissionlessFeedsConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PermissionlessFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PermissionlessFeedsConsumerOwnershipTransferred)
				if err := _PermissionlessFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_PermissionlessFeedsConsumer *PermissionlessFeedsConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*PermissionlessFeedsConsumerOwnershipTransferred, error) {
	event := new(PermissionlessFeedsConsumerOwnershipTransferred)
	if err := _PermissionlessFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
