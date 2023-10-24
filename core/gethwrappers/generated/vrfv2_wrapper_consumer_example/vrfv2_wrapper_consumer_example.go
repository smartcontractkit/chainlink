// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2_wrapper_consumer_example

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

var VRFV2WrapperConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vrfV2Wrapper\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"WrappedRequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"}],\"name\":\"WrapperRequestMade\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"makeRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620011a3380380620011a38339810160408190526200003491620001ac565b6001600160601b0319606083811b821660805282901b1660a0523380600081620000a55760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000d857620000d881620000e3565b5050505050620001e4565b6001600160a01b0381163314156200013e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200009c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001a757600080fd5b919050565b60008060408385031215620001c057600080fd5b620001cb836200018f565b9150620001db602084016200018f565b90509250929050565b60805160601c60a05160601c610f776200022c600039600081816101e30152818161033a015281816107f8015261091f0152600081816104db01526107ce0152610f776000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80638da5cb5b1161005b5780638da5cb5b146100e3578063a168fa891461010b578063d8a4676f1461014a578063f2fde38b1461016c57600080fd5b80630c09b8321461008d5780631fe543e3146100b357806379ba5097146100c85780637a8042bd146100d0575b600080fd5b6100a061009b366004610dc9565b61017f565b6040519081526020015b60405180910390f35b6100c66100c1366004610cda565b610322565b005b6100c66103d4565b6100c66100de366004610ca8565b6104d1565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100aa565b610135610119366004610ca8565b6002602052600090815260409020805460019091015460ff1682565b604080519283529015156020830152016100aa565b61015d610158366004610ca8565b6105d9565b6040516100aa93929190610f11565b6100c661017a366004610c49565b6106eb565b60006101896106ff565b610194848484610782565b6040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff8616600482015290915060009073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690634306d3549060240160206040518083038186803b15801561022557600080fd5b505afa158015610239573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025d9190610cc1565b6040805160608101825282815260006020808301828152845183815280830186528486019081528884526002808452959093208451815590516001820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905591518051959650929491936102e29390850192910190610bd0565b50506040518281528391507f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec49060200160405180910390a2509392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146103c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792056524620563220777261707065722063616e2066756c66696c6c0060448201526064015b60405180910390fd5b6103d082826109c3565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610455576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103bd565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104d96106ff565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb61053460005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401602060405180830381600087803b1580156105a157600080fd5b505af11580156105b5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d09190610c86565b6000818152600260205260408120548190606090610653576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e6400000000000000000000000000000060448201526064016103bd565b6000848152600260208181526040808420815160608101835281548152600182015460ff16151581850152938101805483518186028101860185528181529294938601938301828280156106c657602002820191906000526020600020905b8154815260200190600101908083116106b2575b5050509190925250508151602083015160409093015190989297509550909350505050565b6106f36106ff565b6106fc81610ada565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610780576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103bd565b565b6040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff8416600482015260009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691634000aea0917f00000000000000000000000000000000000000000000000000000000000000009190821690634306d3549060240160206040518083038186803b15801561083d57600080fd5b505afa158015610851573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108759190610cc1565b6040805163ffffffff808b16602083015261ffff8a169282019290925290871660608201526080016040516020818303038152906040526040518463ffffffff1660e01b81526004016108ca93929190610e50565b602060405180830381600087803b1580156108e457600080fd5b505af11580156108f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061091c9190610c86565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fc2a88c36040518163ffffffff1660e01b815260040160206040518083038186803b15801561098357600080fd5b505afa158015610997573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109bb9190610cc1565b949350505050565b600082815260026020526040902054610a38576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e6400000000000000000000000000000060448201526064016103bd565b6000828152600260208181526040909220600181810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690911790558351610a8b93919092019190840190610bd0565b50600082815260026020526040908190205490517f6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b91610ace9185918591610ee8565b60405180910390a15050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610b5a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103bd565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610c0b579160200282015b82811115610c0b578251825591602001919060010190610bf0565b50610c17929150610c1b565b5090565b5b80821115610c175760008155600101610c1c565b803563ffffffff81168114610c4457600080fd5b919050565b600060208284031215610c5b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610c7f57600080fd5b9392505050565b600060208284031215610c9857600080fd5b81518015158114610c7f57600080fd5b600060208284031215610cba57600080fd5b5035919050565b600060208284031215610cd357600080fd5b5051919050565b60008060408385031215610ced57600080fd5b8235915060208084013567ffffffffffffffff80821115610d0d57600080fd5b818601915086601f830112610d2157600080fd5b813581811115610d3357610d33610f3b565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610d7657610d76610f3b565b604052828152858101935084860182860187018b1015610d9557600080fd5b600095505b83861015610db8578035855260019590950194938601938601610d9a565b508096505050505050509250929050565b600080600060608486031215610dde57600080fd5b610de784610c30565b9250602084013561ffff81168114610dfe57600080fd5b9150610e0c60408501610c30565b90509250925092565b600081518084526020808501945080840160005b83811015610e4557815187529582019590820190600101610e29565b509495945050505050565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b81811015610ea057858101830151858201608001528201610e84565b81811115610eb2576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b838152606060208201526000610f016060830185610e15565b9050826040830152949350505050565b8381528215156020820152606060408201526000610f326060830184610e15565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2WrapperConsumerExampleABI = VRFV2WrapperConsumerExampleMetaData.ABI

var VRFV2WrapperConsumerExampleBin = VRFV2WrapperConsumerExampleMetaData.Bin

func DeployVRFV2WrapperConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _vrfV2Wrapper common.Address) (common.Address, *types.Transaction, *VRFV2WrapperConsumerExample, error) {
	parsed, err := VRFV2WrapperConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2WrapperConsumerExampleBin), backend, _link, _vrfV2Wrapper)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2WrapperConsumerExample{address: address, abi: *parsed, VRFV2WrapperConsumerExampleCaller: VRFV2WrapperConsumerExampleCaller{contract: contract}, VRFV2WrapperConsumerExampleTransactor: VRFV2WrapperConsumerExampleTransactor{contract: contract}, VRFV2WrapperConsumerExampleFilterer: VRFV2WrapperConsumerExampleFilterer{contract: contract}}, nil
}

type VRFV2WrapperConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2WrapperConsumerExampleCaller
	VRFV2WrapperConsumerExampleTransactor
	VRFV2WrapperConsumerExampleFilterer
}

type VRFV2WrapperConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2WrapperConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2WrapperConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2WrapperConsumerExampleSession struct {
	Contract     *VRFV2WrapperConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2WrapperConsumerExampleCallerSession struct {
	Contract *VRFV2WrapperConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2WrapperConsumerExampleTransactorSession struct {
	Contract     *VRFV2WrapperConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2WrapperConsumerExampleRaw struct {
	Contract *VRFV2WrapperConsumerExample
}

type VRFV2WrapperConsumerExampleCallerRaw struct {
	Contract *VRFV2WrapperConsumerExampleCaller
}

type VRFV2WrapperConsumerExampleTransactorRaw struct {
	Contract *VRFV2WrapperConsumerExampleTransactor
}

func NewVRFV2WrapperConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFV2WrapperConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2WrapperConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2WrapperConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExample{address: address, abi: abi, VRFV2WrapperConsumerExampleCaller: VRFV2WrapperConsumerExampleCaller{contract: contract}, VRFV2WrapperConsumerExampleTransactor: VRFV2WrapperConsumerExampleTransactor{contract: contract}, VRFV2WrapperConsumerExampleFilterer: VRFV2WrapperConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2WrapperConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2WrapperConsumerExampleCaller, error) {
	contract, err := bindVRFV2WrapperConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleCaller{contract: contract}, nil
}

func NewVRFV2WrapperConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2WrapperConsumerExampleTransactor, error) {
	contract, err := bindVRFV2WrapperConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFV2WrapperConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2WrapperConsumerExampleFilterer, error) {
	contract, err := bindVRFV2WrapperConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFV2WrapperConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2WrapperConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2WrapperConsumerExample.Contract.VRFV2WrapperConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.VRFV2WrapperConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.VRFV2WrapperConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2WrapperConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2WrapperConsumerExample.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Paid = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2WrapperConsumerExample.Contract.GetRequestStatus(&_VRFV2WrapperConsumerExample.CallOpts, _requestId)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2WrapperConsumerExample.Contract.GetRequestStatus(&_VRFV2WrapperConsumerExample.CallOpts, _requestId)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2WrapperConsumerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) Owner() (common.Address, error) {
	return _VRFV2WrapperConsumerExample.Contract.Owner(&_VRFV2WrapperConsumerExample.CallOpts)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2WrapperConsumerExample.Contract.Owner(&_VRFV2WrapperConsumerExample.CallOpts)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2WrapperConsumerExample.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Paid = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2WrapperConsumerExample.Contract.SRequests(&_VRFV2WrapperConsumerExample.CallOpts, arg0)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2WrapperConsumerExample.Contract.SRequests(&_VRFV2WrapperConsumerExample.CallOpts, arg0)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.AcceptOwnership(&_VRFV2WrapperConsumerExample.TransactOpts)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.AcceptOwnership(&_VRFV2WrapperConsumerExample.TransactOpts)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactor) MakeRequest(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.contract.Transact(opts, "makeRequest", _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) MakeRequest(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.MakeRequest(&_VRFV2WrapperConsumerExample.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorSession) MakeRequest(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.MakeRequest(&_VRFV2WrapperConsumerExample.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", _requestId, _randomWords)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) RawFulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2WrapperConsumerExample.TransactOpts, _requestId, _randomWords)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorSession) RawFulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2WrapperConsumerExample.TransactOpts, _requestId, _randomWords)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.TransferOwnership(&_VRFV2WrapperConsumerExample.TransactOpts, to)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.TransferOwnership(&_VRFV2WrapperConsumerExample.TransactOpts, to)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactor) WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.contract.Transact(opts, "withdrawLink", amount)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleSession) WithdrawLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.WithdrawLink(&_VRFV2WrapperConsumerExample.TransactOpts, amount)
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleTransactorSession) WithdrawLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2WrapperConsumerExample.Contract.WithdrawLink(&_VRFV2WrapperConsumerExample.TransactOpts, amount)
}

type VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2WrapperConsumerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperConsumerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2WrapperConsumerExampleOwnershipTransferRequested)
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

func (it *VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperConsumerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator{contract: _VRFV2WrapperConsumerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperConsumerExampleOwnershipTransferRequested)
				if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2WrapperConsumerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2WrapperConsumerExampleOwnershipTransferRequested)
	if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2WrapperConsumerExampleOwnershipTransferredIterator struct {
	Event *VRFV2WrapperConsumerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperConsumerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperConsumerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2WrapperConsumerExampleOwnershipTransferred)
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

func (it *VRFV2WrapperConsumerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperConsumerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperConsumerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperConsumerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleOwnershipTransferredIterator{contract: _VRFV2WrapperConsumerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperConsumerExampleOwnershipTransferred)
				if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2WrapperConsumerExampleOwnershipTransferred, error) {
	event := new(VRFV2WrapperConsumerExampleOwnershipTransferred)
	if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator struct {
	Event *VRFV2WrapperConsumerExampleWrappedRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperConsumerExampleWrappedRequestFulfilled)
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
		it.Event = new(VRFV2WrapperConsumerExampleWrappedRequestFulfilled)
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

func (it *VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperConsumerExampleWrappedRequestFulfilled struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Payment     *big.Int
	Raw         types.Log
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator, error) {

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.FilterLogs(opts, "WrappedRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator{contract: _VRFV2WrapperConsumerExample.contract, event: "WrappedRequestFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleWrappedRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.WatchLogs(opts, "WrappedRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperConsumerExampleWrappedRequestFulfilled)
				if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "WrappedRequestFulfilled", log); err != nil {
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

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) ParseWrappedRequestFulfilled(log types.Log) (*VRFV2WrapperConsumerExampleWrappedRequestFulfilled, error) {
	event := new(VRFV2WrapperConsumerExampleWrappedRequestFulfilled)
	if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "WrappedRequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2WrapperConsumerExampleWrapperRequestMadeIterator struct {
	Event *VRFV2WrapperConsumerExampleWrapperRequestMade

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperConsumerExampleWrapperRequestMadeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperConsumerExampleWrapperRequestMade)
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
		it.Event = new(VRFV2WrapperConsumerExampleWrapperRequestMade)
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

func (it *VRFV2WrapperConsumerExampleWrapperRequestMadeIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperConsumerExampleWrapperRequestMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperConsumerExampleWrapperRequestMade struct {
	RequestId *big.Int
	Paid      *big.Int
	Raw       types.Log
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) FilterWrapperRequestMade(opts *bind.FilterOpts, requestId []*big.Int) (*VRFV2WrapperConsumerExampleWrapperRequestMadeIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.FilterLogs(opts, "WrapperRequestMade", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperConsumerExampleWrapperRequestMadeIterator{contract: _VRFV2WrapperConsumerExample.contract, event: "WrapperRequestMade", logs: logs, sub: sub}, nil
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) WatchWrapperRequestMade(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleWrapperRequestMade, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFV2WrapperConsumerExample.contract.WatchLogs(opts, "WrapperRequestMade", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperConsumerExampleWrapperRequestMade)
				if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "WrapperRequestMade", log); err != nil {
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

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExampleFilterer) ParseWrapperRequestMade(log types.Log) (*VRFV2WrapperConsumerExampleWrapperRequestMade, error) {
	event := new(VRFV2WrapperConsumerExampleWrapperRequestMade)
	if err := _VRFV2WrapperConsumerExample.contract.UnpackLog(event, "WrapperRequestMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Paid        *big.Int
	Fulfilled   bool
	RandomWords []*big.Int
}
type SRequests struct {
	Paid      *big.Int
	Fulfilled bool
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2WrapperConsumerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2WrapperConsumerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2WrapperConsumerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2WrapperConsumerExample.ParseOwnershipTransferred(log)
	case _VRFV2WrapperConsumerExample.abi.Events["WrappedRequestFulfilled"].ID:
		return _VRFV2WrapperConsumerExample.ParseWrappedRequestFulfilled(log)
	case _VRFV2WrapperConsumerExample.abi.Events["WrapperRequestMade"].ID:
		return _VRFV2WrapperConsumerExample.ParseWrapperRequestMade(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2WrapperConsumerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2WrapperConsumerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2WrapperConsumerExampleWrappedRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b")
}

func (VRFV2WrapperConsumerExampleWrapperRequestMade) Topic() common.Hash {
	return common.HexToHash("0x5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec4")
}

func (_VRFV2WrapperConsumerExample *VRFV2WrapperConsumerExample) Address() common.Address {
	return _VRFV2WrapperConsumerExample.address
}

type VRFV2WrapperConsumerExampleInterface interface {
	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	MakeRequest(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2WrapperConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2WrapperConsumerExampleOwnershipTransferred, error)

	FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator, error)

	WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleWrappedRequestFulfilled) (event.Subscription, error)

	ParseWrappedRequestFulfilled(log types.Log) (*VRFV2WrapperConsumerExampleWrappedRequestFulfilled, error)

	FilterWrapperRequestMade(opts *bind.FilterOpts, requestId []*big.Int) (*VRFV2WrapperConsumerExampleWrapperRequestMadeIterator, error)

	WatchWrapperRequestMade(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleWrapperRequestMade, requestId []*big.Int) (event.Subscription, error)

	ParseWrapperRequestMade(log types.Log) (*VRFV2WrapperConsumerExampleWrapperRequestMade, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
