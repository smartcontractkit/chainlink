// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package operator_consumer

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ConsumerABI is the input ABI used to generate the binding from.
const ConsumerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"price\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"}],\"name\":\"addExternalRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentPrice\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_price\",\"type\":\"bytes32\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestEthereumPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_callback\",\"type\":\"address\"}],\"name\":\"requestEthereumPriceByCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ConsumerBin is the compiled bytecode used for deploying new contracts.
var ConsumerBin = "0x6080604052600160045534801561001557600080fd5b506040516113133803806113138339818101604052606081101561003857600080fd5b508051602082015160409092015190919061005283610066565b61005b82610088565b600655506100aa9050565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b61125a806100b96000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806383db5cbc1161005b57806383db5cbc146101d55780638dc654a21461027d5780639d1b464a14610285578063e8d5359d1461029f5761007d565b8063042f2b65146100825780635591a608146100a757806374961d4d14610114575b600080fd5b6100a56004803603604081101561009857600080fd5b50803590602001356102d8565b005b6100a5600480360360a08110156100bd57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020810135906040810135907fffffffff0000000000000000000000000000000000000000000000000000000060608201351690608001356103e5565b6100a56004803603606081101561012a57600080fd5b81019060208101813564010000000081111561014557600080fd5b82018360208201111561015757600080fd5b8035906020019184600183028401116401000000008311171561017957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550508235935050506020013573ffffffffffffffffffffffffffffffffffffffff166104ac565b6100a5600480360360408110156101eb57600080fd5b81019060208101813564010000000081111561020657600080fd5b82018360208201111561021857600080fd5b8035906020019184600183028401116401000000008311171561023a57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050913592506105e1915050565b6100a56105f0565b61028d6107ba565b60408051918252519081900360200190f35b6100a5600480360360408110156102b557600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356107c0565b600082815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff163314610356576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806111b66028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a2604051829084907f0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c458503631190600090a35060075550565b604080517f6ee4d55300000000000000000000000000000000000000000000000000000000815260048101869052602481018590527fffffffff0000000000000000000000000000000000000000000000000000000084166044820152606481018390529051869173ffffffffffffffffffffffffffffffffffffffff831691636ee4d5539160848082019260009290919082900301818387803b15801561048c57600080fd5b505af11580156104a0573d6000803e3d6000fd5b50505050505050505050565b6104b4611143565b6006546104e290837f042f2b65000000000000000000000000000000000000000000000000000000006107ca565b905061053e6040518060400160405280600381526020017f67657400000000000000000000000000000000000000000000000000000000008152506040518060800160405280604781526020016111de604791398391906107ef565b604080516001808252818301909252606091816020015b6060815260200190600190039081610555579050509050848160008151811061057a57fe5b60200260200101819052506105cf6040518060400160405280600481526020017f706174680000000000000000000000000000000000000000000000000000000081525082846108129092919063ffffffff16565b6105d9828561087a565b505050505050565b6105ec8282306104ac565b5050565b60006105fa6108aa565b90508073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb338373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561068057600080fd5b505afa158015610694573d6000803e3d6000fd5b505050506040513d60208110156106aa57600080fd5b5051604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff909316600484015260248301919091525160448083019260209291908290030181600087803b15801561072057600080fd5b505af1158015610734573d6000803e3d6000fd5b505050506040513d602081101561074a57600080fd5b50516107b757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015290519081900360640190fd5b50565b60075481565b6105ec82826108c6565b6107d2611143565b6107da611143565b6107e6818686866109ad565b95945050505050565b60808301516107fe9083610a0f565b608083015161080d9082610a0f565b505050565b60808301516108219083610a0f565b61082e8360800151610a26565b60005b815181101561086c5761086482828151811061084957fe5b60200260200101518560800151610a0f90919063ffffffff16565b600101610831565b5061080d8360800151610a31565b6003546000906108a19073ffffffffffffffffffffffffffffffffffffffff168484610a3c565b90505b92915050565b60025473ffffffffffffffffffffffffffffffffffffffff1690565b600081815260056020526040902054819073ffffffffffffffffffffffffffffffffffffffff161561095957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f5265717565737420697320616c72656164792070656e64696e67000000000000604482015290519081900360640190fd5b50600090815260056020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6109b5611143565b6109c58560800151610100610c6a565b505091835273ffffffffffffffffffffffffffffffffffffffff1660208301527fffffffff0000000000000000000000000000000000000000000000000000000016604082015290565b610a1c8260038351610ca4565b61080d8282610d7e565b6107b7816004610d98565b6107b7816007610d98565b6004546040805130606090811b60208084019190915260348084018690528451808503909101815260549093018452825192810192909220908601939093526000838152600590915281812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8816179055905182917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af991a26002805473ffffffffffffffffffffffffffffffffffffffff1690634000aea09086908590610b1d908890610dad565b6040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610b8b578181015183820152602001610b73565b50505050905090810190601f168015610bb85780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610bd957600080fd5b505af1158015610bed573d6000803e3d6000fd5b505050506040513d6020811015610c0357600080fd5b5051610c5a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806111936023913960400191505060405180910390fd5b6004805460010190559392505050565b610c72611178565b6020820615610c875760208206602003820191505b506020828101829052604080518085526000815290920101905290565b60178111610cc557610cbf8360e0600585901b168317610f4b565b5061080d565b60ff8111610cef57610ce2836018611fe0600586901b1617610f4b565b50610cbf83826001610f63565b61ffff8111610d1a57610d0d836019611fe0600586901b1617610f4b565b50610cbf83826002610f63565b63ffffffff8111610d4757610d3a83601a611fe0600586901b1617610f4b565b50610cbf83826004610f63565b67ffffffffffffffff811161080d57610d6b83601b611fe0600586901b1617610f4b565b50610d7883826008610f63565b50505050565b610d86611178565b6108a183846000015151848551610f84565b61080d82601f611fe0600585901b1617610f4b565b6060634042994660e01b6000808560000151866020015187604001518860600151888a6080015160000151604051602401808973ffffffffffffffffffffffffffffffffffffffff1681526020018881526020018781526020018673ffffffffffffffffffffffffffffffffffffffff168152602001857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610e8d578181015183820152602001610e75565b50505050905090810190601f168015610eba5780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909d169c909c17909b525098995050505050505050505092915050565b610f53611178565b6108a1838460000151518461106c565b610f6b611178565b610f7c8485600001515185856110b7565b949350505050565b610f8c611178565b8251821115610f9a57600080fd5b84602001518285011115610fc457610fc485610fbc8760200151878601611115565b60020261112c565b600080865180518760208301019350808887011115610fe35787860182525b505050602084015b6020841061102857805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09093019260209182019101610feb565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b611074611178565b836020015183106110905761109084856020015160020261112c565b8351805160208583010184815350808514156110ad576001810182525b5093949350505050565b6110bf611178565b846020015184830111156110dc576110dc8585840160020261112c565b60006001836101000a03905085518386820101858319825116178152508051848701111561110a5783860181525b509495945050505050565b6000818311156111265750816108a4565b50919050565b81516111388383610c6a565b50610d788382610d7e565b6040805160a081018252600080825260208201819052918101829052606081019190915260808101611173611178565b905290565b60405180604001604052806060815260200160008152509056fe756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c65536f75726365206d75737420626520746865206f7261636c65206f6620746865207265717565737468747470733a2f2f6d696e2d6170692e63727970746f636f6d706172652e636f6d2f646174612f70726963653f6673796d3d455448267473796d733d5553442c4555522c4a5059a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

// DeployConsumer deploys a new Ethereum contract, binding an instance of Consumer to it.
func DeployConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _oracle common.Address, _specId [32]byte) (common.Address, *types.Transaction, *Consumer, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsumerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ConsumerBin), backend, _link, _oracle, _specId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Consumer{ConsumerCaller: ConsumerCaller{contract: contract}, ConsumerTransactor: ConsumerTransactor{contract: contract}, ConsumerFilterer: ConsumerFilterer{contract: contract}}, nil
}

// Consumer is an auto generated Go binding around an Ethereum contract.
type Consumer struct {
	ConsumerCaller     // Read-only binding to the contract
	ConsumerTransactor // Write-only binding to the contract
	ConsumerFilterer   // Log filterer for contract events
}

// ConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsumerSession struct {
	Contract     *Consumer         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsumerCallerSession struct {
	Contract *ConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsumerTransactorSession struct {
	Contract     *ConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsumerRaw struct {
	Contract *Consumer // Generic contract binding to access the raw methods on
}

// ConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsumerCallerRaw struct {
	Contract *ConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// ConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsumerTransactorRaw struct {
	Contract *ConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsumer creates a new instance of Consumer, bound to a specific deployed contract.
func NewConsumer(address common.Address, backend bind.ContractBackend) (*Consumer, error) {
	contract, err := bindConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Consumer{ConsumerCaller: ConsumerCaller{contract: contract}, ConsumerTransactor: ConsumerTransactor{contract: contract}, ConsumerFilterer: ConsumerFilterer{contract: contract}}, nil
}

// NewConsumerCaller creates a new read-only instance of Consumer, bound to a specific deployed contract.
func NewConsumerCaller(address common.Address, caller bind.ContractCaller) (*ConsumerCaller, error) {
	contract, err := bindConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsumerCaller{contract: contract}, nil
}

// NewConsumerTransactor creates a new write-only instance of Consumer, bound to a specific deployed contract.
func NewConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsumerTransactor, error) {
	contract, err := bindConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsumerTransactor{contract: contract}, nil
}

// NewConsumerFilterer creates a new log filterer instance of Consumer, bound to a specific deployed contract.
func NewConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsumerFilterer, error) {
	contract, err := bindConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsumerFilterer{contract: contract}, nil
}

// bindConsumer binds a generic wrapper to an already deployed contract.
func bindConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consumer *ConsumerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Consumer.Contract.ConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consumer *ConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consumer.Contract.ConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consumer *ConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consumer.Contract.ConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consumer *ConsumerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Consumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consumer *ConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consumer *ConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consumer.Contract.contract.Transact(opts, method, params...)
}

// CurrentPrice is a free data retrieval call binding the contract method 0x9d1b464a.
//
// Solidity: function currentPrice() view returns(bytes32)
func (_Consumer *ConsumerCaller) CurrentPrice(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Consumer.contract.Call(opts, out, "currentPrice")
	return *ret0, err
}

// CurrentPrice is a free data retrieval call binding the contract method 0x9d1b464a.
//
// Solidity: function currentPrice() view returns(bytes32)
func (_Consumer *ConsumerSession) CurrentPrice() ([32]byte, error) {
	return _Consumer.Contract.CurrentPrice(&_Consumer.CallOpts)
}

// CurrentPrice is a free data retrieval call binding the contract method 0x9d1b464a.
//
// Solidity: function currentPrice() view returns(bytes32)
func (_Consumer *ConsumerCallerSession) CurrentPrice() ([32]byte, error) {
	return _Consumer.Contract.CurrentPrice(&_Consumer.CallOpts)
}

// AddExternalRequest is a paid mutator transaction binding the contract method 0xe8d5359d.
//
// Solidity: function addExternalRequest(address _oracle, bytes32 _requestId) returns()
func (_Consumer *ConsumerTransactor) AddExternalRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "addExternalRequest", _oracle, _requestId)
}

// AddExternalRequest is a paid mutator transaction binding the contract method 0xe8d5359d.
//
// Solidity: function addExternalRequest(address _oracle, bytes32 _requestId) returns()
func (_Consumer *ConsumerSession) AddExternalRequest(_oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.AddExternalRequest(&_Consumer.TransactOpts, _oracle, _requestId)
}

// AddExternalRequest is a paid mutator transaction binding the contract method 0xe8d5359d.
//
// Solidity: function addExternalRequest(address _oracle, bytes32 _requestId) returns()
func (_Consumer *ConsumerTransactorSession) AddExternalRequest(_oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.AddExternalRequest(&_Consumer.TransactOpts, _oracle, _requestId)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x5591a608.
//
// Solidity: function cancelRequest(address _oracle, bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_Consumer *ConsumerTransactor) CancelRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "cancelRequest", _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x5591a608.
//
// Solidity: function cancelRequest(address _oracle, bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_Consumer *ConsumerSession) CancelRequest(_oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.CancelRequest(&_Consumer.TransactOpts, _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x5591a608.
//
// Solidity: function cancelRequest(address _oracle, bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_Consumer *ConsumerTransactorSession) CancelRequest(_oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.CancelRequest(&_Consumer.TransactOpts, _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

// Fulfill is a paid mutator transaction binding the contract method 0x042f2b65.
//
// Solidity: function fulfill(bytes32 _requestId, bytes32 _price) returns()
func (_Consumer *ConsumerTransactor) Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _price [32]byte) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "fulfill", _requestId, _price)
}

// Fulfill is a paid mutator transaction binding the contract method 0x042f2b65.
//
// Solidity: function fulfill(bytes32 _requestId, bytes32 _price) returns()
func (_Consumer *ConsumerSession) Fulfill(_requestId [32]byte, _price [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.Fulfill(&_Consumer.TransactOpts, _requestId, _price)
}

// Fulfill is a paid mutator transaction binding the contract method 0x042f2b65.
//
// Solidity: function fulfill(bytes32 _requestId, bytes32 _price) returns()
func (_Consumer *ConsumerTransactorSession) Fulfill(_requestId [32]byte, _price [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.Fulfill(&_Consumer.TransactOpts, _requestId, _price)
}

// RequestEthereumPrice is a paid mutator transaction binding the contract method 0x83db5cbc.
//
// Solidity: function requestEthereumPrice(string _currency, uint256 _payment) returns()
func (_Consumer *ConsumerTransactor) RequestEthereumPrice(opts *bind.TransactOpts, _currency string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "requestEthereumPrice", _currency, _payment)
}

// RequestEthereumPrice is a paid mutator transaction binding the contract method 0x83db5cbc.
//
// Solidity: function requestEthereumPrice(string _currency, uint256 _payment) returns()
func (_Consumer *ConsumerSession) RequestEthereumPrice(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPrice(&_Consumer.TransactOpts, _currency, _payment)
}

// RequestEthereumPrice is a paid mutator transaction binding the contract method 0x83db5cbc.
//
// Solidity: function requestEthereumPrice(string _currency, uint256 _payment) returns()
func (_Consumer *ConsumerTransactorSession) RequestEthereumPrice(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPrice(&_Consumer.TransactOpts, _currency, _payment)
}

// RequestEthereumPriceByCallback is a paid mutator transaction binding the contract method 0x74961d4d.
//
// Solidity: function requestEthereumPriceByCallback(string _currency, uint256 _payment, address _callback) returns()
func (_Consumer *ConsumerTransactor) RequestEthereumPriceByCallback(opts *bind.TransactOpts, _currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "requestEthereumPriceByCallback", _currency, _payment, _callback)
}

// RequestEthereumPriceByCallback is a paid mutator transaction binding the contract method 0x74961d4d.
//
// Solidity: function requestEthereumPriceByCallback(string _currency, uint256 _payment, address _callback) returns()
func (_Consumer *ConsumerSession) RequestEthereumPriceByCallback(_currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPriceByCallback(&_Consumer.TransactOpts, _currency, _payment, _callback)
}

// RequestEthereumPriceByCallback is a paid mutator transaction binding the contract method 0x74961d4d.
//
// Solidity: function requestEthereumPriceByCallback(string _currency, uint256 _payment, address _callback) returns()
func (_Consumer *ConsumerTransactorSession) RequestEthereumPriceByCallback(_currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPriceByCallback(&_Consumer.TransactOpts, _currency, _payment, _callback)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_Consumer *ConsumerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "withdrawLink")
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_Consumer *ConsumerSession) WithdrawLink() (*types.Transaction, error) {
	return _Consumer.Contract.WithdrawLink(&_Consumer.TransactOpts)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_Consumer *ConsumerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _Consumer.Contract.WithdrawLink(&_Consumer.TransactOpts)
}

// ConsumerChainlinkCancelledIterator is returned from FilterChainlinkCancelled and is used to iterate over the raw logs and unpacked data for ChainlinkCancelled events raised by the Consumer contract.
type ConsumerChainlinkCancelledIterator struct {
	Event *ConsumerChainlinkCancelled // Event containing the contract specifics and raw log

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
func (it *ConsumerChainlinkCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerChainlinkCancelled)
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
		it.Event = new(ConsumerChainlinkCancelled)
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
func (it *ConsumerChainlinkCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsumerChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsumerChainlinkCancelled represents a ChainlinkCancelled event raised by the Consumer contract.
type ConsumerChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkCancelled is a free log retrieval operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerChainlinkCancelledIterator{contract: _Consumer.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

// WatchChainlinkCancelled is a free log subscription operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsumerChainlinkCancelled)
				if err := _Consumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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

// ParseChainlinkCancelled is a log parse operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) ParseChainlinkCancelled(log types.Log) (*ConsumerChainlinkCancelled, error) {
	event := new(ConsumerChainlinkCancelled)
	if err := _Consumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ConsumerChainlinkFulfilledIterator is returned from FilterChainlinkFulfilled and is used to iterate over the raw logs and unpacked data for ChainlinkFulfilled events raised by the Consumer contract.
type ConsumerChainlinkFulfilledIterator struct {
	Event *ConsumerChainlinkFulfilled // Event containing the contract specifics and raw log

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
func (it *ConsumerChainlinkFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerChainlinkFulfilled)
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
		it.Event = new(ConsumerChainlinkFulfilled)
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
func (it *ConsumerChainlinkFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsumerChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsumerChainlinkFulfilled represents a ChainlinkFulfilled event raised by the Consumer contract.
type ConsumerChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkFulfilled is a free log retrieval operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerChainlinkFulfilledIterator{contract: _Consumer.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

// WatchChainlinkFulfilled is a free log subscription operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsumerChainlinkFulfilled)
				if err := _Consumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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

// ParseChainlinkFulfilled is a log parse operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) ParseChainlinkFulfilled(log types.Log) (*ConsumerChainlinkFulfilled, error) {
	event := new(ConsumerChainlinkFulfilled)
	if err := _Consumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ConsumerChainlinkRequestedIterator is returned from FilterChainlinkRequested and is used to iterate over the raw logs and unpacked data for ChainlinkRequested events raised by the Consumer contract.
type ConsumerChainlinkRequestedIterator struct {
	Event *ConsumerChainlinkRequested // Event containing the contract specifics and raw log

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
func (it *ConsumerChainlinkRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerChainlinkRequested)
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
		it.Event = new(ConsumerChainlinkRequested)
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
func (it *ConsumerChainlinkRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsumerChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsumerChainlinkRequested represents a ChainlinkRequested event raised by the Consumer contract.
type ConsumerChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkRequested is a free log retrieval operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerChainlinkRequestedIterator{contract: _Consumer.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

// WatchChainlinkRequested is a free log subscription operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsumerChainlinkRequested)
				if err := _Consumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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

// ParseChainlinkRequested is a log parse operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_Consumer *ConsumerFilterer) ParseChainlinkRequested(log types.Log) (*ConsumerChainlinkRequested, error) {
	event := new(ConsumerChainlinkRequested)
	if err := _Consumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ConsumerRequestFulfilledIterator is returned from FilterRequestFulfilled and is used to iterate over the raw logs and unpacked data for RequestFulfilled events raised by the Consumer contract.
type ConsumerRequestFulfilledIterator struct {
	Event *ConsumerRequestFulfilled // Event containing the contract specifics and raw log

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
func (it *ConsumerRequestFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerRequestFulfilled)
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
		it.Event = new(ConsumerRequestFulfilled)
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
func (it *ConsumerRequestFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsumerRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsumerRequestFulfilled represents a RequestFulfilled event raised by the Consumer contract.
type ConsumerRequestFulfilled struct {
	RequestId [32]byte
	Price     [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestFulfilled is a free log retrieval operation binding the contract event 0x0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c4585036311.
//
// Solidity: event RequestFulfilled(bytes32 indexed requestId, bytes32 indexed price)
func (_Consumer *ConsumerFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, requestId [][32]byte, price [][32]byte) (*ConsumerRequestFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "RequestFulfilled", requestIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerRequestFulfilledIterator{contract: _Consumer.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

// WatchRequestFulfilled is a free log subscription operation binding the contract event 0x0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c4585036311.
//
// Solidity: event RequestFulfilled(bytes32 indexed requestId, bytes32 indexed price)
func (_Consumer *ConsumerFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *ConsumerRequestFulfilled, requestId [][32]byte, price [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "RequestFulfilled", requestIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsumerRequestFulfilled)
				if err := _Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

// ParseRequestFulfilled is a log parse operation binding the contract event 0x0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c4585036311.
//
// Solidity: event RequestFulfilled(bytes32 indexed requestId, bytes32 indexed price)
func (_Consumer *ConsumerFilterer) ParseRequestFulfilled(log types.Log) (*ConsumerRequestFulfilled, error) {
	event := new(ConsumerRequestFulfilled)
	if err := _Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	return event, nil
}
