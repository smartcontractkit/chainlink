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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var VRFV2WrapperConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vrfV2Wrapper\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"WrappedRequestFulfilled\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"makeRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b50604051610fd8380380610fd883398101604081905261002f91610197565b6001600160a01b03808316608052811660a05233806000816100985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100c8576100c8816100d2565b50505050506101ca565b336001600160a01b0382160361012a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b038116811461019257600080fd5b919050565b600080604083850312156101aa57600080fd5b6101b38361017b565b91506101c16020840161017b565b90509250929050565b60805160a051610dd4610204600039600081816101b8015281816102d40152818161068a0152610793015260006106600152610dd46000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80638da5cb5b1161005b5780638da5cb5b146100c5578063a168fa89146100ed578063d8a4676f1461012c578063f2fde38b1461014e57600080fd5b80630c09b832146100825780631fe543e3146100a857806379ba5097146100bd575b600080fd5b610095610090366004610aad565b610161565b6040519081526020015b60405180910390f35b6100bb6100b6366004610b28565b6102bc565b005b6100bb61036e565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161009f565b6101176100fb366004610c10565b6002602052600090815260409020805460019091015460ff1682565b6040805192835290151560208301520161009f565b61013f61013a366004610c10565b61046b565b60405161009f93929190610c64565b6100bb61015c366004610c8e565b61057d565b600061016b610591565b610176848484610614565b6040805160608101918290527f4306d3540000000000000000000000000000000000000000000000000000000090915263ffffffff86166064820152909150807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634306d35460848301602060405180830381865afa158015610213573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102379190610ccb565b8152600060208083018290526040805183815280830182529381019390935284825260028082529183902084518155848201516001820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790559284015180516102b2938501929190910190610a34565b5050509392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610360576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792056524620563220777261707065722063616e2066756c66696c6c0060448201526064015b60405180910390fd5b61036a8282610828565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103ef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610357565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60008181526002602052604081205481906060906104e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610357565b6000848152600260208181526040808420815160608101835281548152600182015460ff161515818501529381018054835181860281018601855281815292949386019383018282801561055857602002820191906000526020600020905b815481526020019060010190808311610544575b5050509190925250508151602083015160409093015190989297509550909350505050565b610585610591565b61058e8161093f565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610612576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610357565b565b6040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff8416600482015260009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691634000aea0917f00000000000000000000000000000000000000000000000000000000000000009190821690634306d35490602401602060405180830381865afa1580156106d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106f89190610ccb565b6040805163ffffffff808b16602083015261ffff8a169282019290925290871660608201526080016040516020818303038152906040526040518463ffffffff1660e01b815260040161074d93929190610ce4565b6020604051808303816000875af115801561076c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107909190610d7c565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fc2a88c36040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107fc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108209190610ccb565b949350505050565b60008281526002602052604090205461089d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610357565b6000828152600260208181526040909220600181810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905583516108f093919092019190840190610a34565b50600082815260026020526040908190205490517f6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b916109339185918591610d9e565b60405180910390a15050565b3373ffffffffffffffffffffffffffffffffffffffff8216036109be576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610357565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610a6f579160200282015b82811115610a6f578251825591602001919060010190610a54565b50610a7b929150610a7f565b5090565b5b80821115610a7b5760008155600101610a80565b803563ffffffff81168114610aa857600080fd5b919050565b600080600060608486031215610ac257600080fd5b610acb84610a94565b9250602084013561ffff81168114610ae257600080fd5b9150610af060408501610a94565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060408385031215610b3b57600080fd5b8235915060208084013567ffffffffffffffff80821115610b5b57600080fd5b818601915086601f830112610b6f57600080fd5b813581811115610b8157610b81610af9565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610bc457610bc4610af9565b604052918252848201925083810185019189831115610be257600080fd5b938501935b82851015610c0057843584529385019392850192610be7565b8096505050505050509250929050565b600060208284031215610c2257600080fd5b5035919050565b600081518084526020808501945080840160005b83811015610c5957815187529582019590820190600101610c3d565b509495945050505050565b8381528215156020820152606060408201526000610c856060830184610c29565b95945050505050565b600060208284031215610ca057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610cc457600080fd5b9392505050565b600060208284031215610cdd57600080fd5b5051919050565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b81811015610d3457858101830151858201608001528201610d18565b81811115610d46576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b600060208284031215610d8e57600080fd5b81518015158114610cc457600080fd5b838152606060208201526000610db76060830185610c29565b905082604083015294935050505056fea164736f6c634300080d000a",
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
	return address, tx, &VRFV2WrapperConsumerExample{VRFV2WrapperConsumerExampleCaller: VRFV2WrapperConsumerExampleCaller{contract: contract}, VRFV2WrapperConsumerExampleTransactor: VRFV2WrapperConsumerExampleTransactor{contract: contract}, VRFV2WrapperConsumerExampleFilterer: VRFV2WrapperConsumerExampleFilterer{contract: contract}}, nil
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
	parsed, err := abi.JSON(strings.NewReader(VRFV2WrapperConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2WrapperConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2WrapperConsumerExampleOwnershipTransferred, error)

	FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2WrapperConsumerExampleWrappedRequestFulfilledIterator, error)

	WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperConsumerExampleWrappedRequestFulfilled) (event.Subscription, error)

	ParseWrappedRequestFulfilled(log types.Log) (*VRFV2WrapperConsumerExampleWrappedRequestFulfilled, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
