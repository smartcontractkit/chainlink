// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_sub_owner

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

var VRFV2PlusExternalSubOwnerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610db9380380610db983398101604081905261002f916101d0565b8133806000816100865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b6576100b68161010a565b5050600280546001600160a01b039384166001600160a01b03199182161790915560038054958416958216959095179094555060048054929091169183169190911790556007805490911633179055610203565b6001600160a01b0381163314156101635760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146101cb57600080fd5b919050565b600080604083850312156101e357600080fd5b6101ec836101b4565b91506101fa602084016101b4565b90509250929050565b610ba7806102126000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063e89e106a1161005b578063e89e106a146100ea578063f2fde38b14610101578063f6eaffc814610114578063f793a70e1461012757600080fd5b80631fe543e31461008d5780632d6d99f3146100a257806379ba5097146100b55780638da5cb5b146100bd575b600080fd5b6100a061009b366004610933565b61013a565b005b6100a06100b03660046108ce565b6101c0565b6100a06102f0565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100f360065481565b6040519081526020016100e1565b6100a061010f3660046108ac565b6103ed565b6100f3610122366004610901565b610401565b6100a0610135366004610a22565b610422565b60025473ffffffffffffffffffffffffffffffffffffffff1633146101b2576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6101bc828261053f565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610200575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610284573361022560005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff938416600482015291831660248301529190911660448201526064016101a9565b6002805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090921673ffffffffffffffffffffffffffffffffffffffff90931692909217179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610371576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016101a9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6103f56105c2565b6103fe81610645565b50565b6005818154811061041157600080fd5b600091825260209091200154905081565b61042a6105c2565b60006040518060c001604052808481526020018867ffffffffffffffff1681526020018661ffff1681526020018763ffffffff1681526020018563ffffffff168152602001610488604051806020016040528086151581525061073b565b90526003546040517f596b8b8800000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff169063596b8b88906104e1908490600401610aa1565b602060405180830381600087803b1580156104fb57600080fd5b505af115801561050f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610533919061091a565b60065550505050505050565b60065482146105aa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016101a9565b80516105bd9060059060208401906107f7565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610643576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016101a9565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156106c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016101a9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161077491511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b828054828255906000526020600020908101928215610832579160200282015b82811115610832578251825591602001919060010190610817565b5061083e929150610842565b5090565b5b8082111561083e5760008155600101610843565b803573ffffffffffffffffffffffffffffffffffffffff8116811461087b57600080fd5b919050565b803563ffffffff8116811461087b57600080fd5b803567ffffffffffffffff8116811461087b57600080fd5b6000602082840312156108be57600080fd5b6108c782610857565b9392505050565b600080604083850312156108e157600080fd5b6108ea83610857565b91506108f860208401610894565b90509250929050565b60006020828403121561091357600080fd5b5035919050565b60006020828403121561092c57600080fd5b5051919050565b6000806040838503121561094657600080fd5b8235915060208084013567ffffffffffffffff8082111561096657600080fd5b818601915086601f83011261097a57600080fd5b81358181111561098c5761098c610b6b565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156109cf576109cf610b6b565b604052828152858101935084860182860187018b10156109ee57600080fd5b600095505b83861015610a115780358552600195909501949386019386016109f3565b508096505050505050509250929050565b60008060008060008060c08789031215610a3b57600080fd5b610a4487610894565b9550610a5260208801610880565b9450604087013561ffff81168114610a6957600080fd5b9350610a7760608801610880565b92506080870135915060a08701358015158114610a9357600080fd5b809150509295509295509295565b6000602080835283518184015267ffffffffffffffff8185015116604084015261ffff6040850151166060840152606084015163ffffffff80821660808601528060808701511660a0860152505060a084015160c08085015280518060e086015260005b81811015610b225782810184015186820161010001528301610b05565b81811115610b3557600061010083880101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169390930161010001949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusExternalSubOwnerExampleABI = VRFV2PlusExternalSubOwnerExampleMetaData.ABI

var VRFV2PlusExternalSubOwnerExampleBin = VRFV2PlusExternalSubOwnerExampleMetaData.Bin

func DeployVRFV2PlusExternalSubOwnerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusExternalSubOwnerExample, error) {
	parsed, err := VRFV2PlusExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusExternalSubOwnerExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusExternalSubOwnerExample{VRFV2PlusExternalSubOwnerExampleCaller: VRFV2PlusExternalSubOwnerExampleCaller{contract: contract}, VRFV2PlusExternalSubOwnerExampleTransactor: VRFV2PlusExternalSubOwnerExampleTransactor{contract: contract}, VRFV2PlusExternalSubOwnerExampleFilterer: VRFV2PlusExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusExternalSubOwnerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusExternalSubOwnerExampleCaller
	VRFV2PlusExternalSubOwnerExampleTransactor
	VRFV2PlusExternalSubOwnerExampleFilterer
}

type VRFV2PlusExternalSubOwnerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusExternalSubOwnerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusExternalSubOwnerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusExternalSubOwnerExampleSession struct {
	Contract     *VRFV2PlusExternalSubOwnerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusExternalSubOwnerExampleCallerSession struct {
	Contract *VRFV2PlusExternalSubOwnerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusExternalSubOwnerExampleTransactorSession struct {
	Contract     *VRFV2PlusExternalSubOwnerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusExternalSubOwnerExampleRaw struct {
	Contract *VRFV2PlusExternalSubOwnerExample
}

type VRFV2PlusExternalSubOwnerExampleCallerRaw struct {
	Contract *VRFV2PlusExternalSubOwnerExampleCaller
}

type VRFV2PlusExternalSubOwnerExampleTransactorRaw struct {
	Contract *VRFV2PlusExternalSubOwnerExampleTransactor
}

func NewVRFV2PlusExternalSubOwnerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusExternalSubOwnerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExample{address: address, abi: abi, VRFV2PlusExternalSubOwnerExampleCaller: VRFV2PlusExternalSubOwnerExampleCaller{contract: contract}, VRFV2PlusExternalSubOwnerExampleTransactor: VRFV2PlusExternalSubOwnerExampleTransactor{contract: contract}, VRFV2PlusExternalSubOwnerExampleFilterer: VRFV2PlusExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusExternalSubOwnerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusExternalSubOwnerExampleCaller, error) {
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusExternalSubOwnerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusExternalSubOwnerExampleTransactor, error) {
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusExternalSubOwnerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusExternalSubOwnerExampleFilterer, error) {
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusExternalSubOwnerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusExternalSubOwnerExample.Contract.VRFV2PlusExternalSubOwnerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.VRFV2PlusExternalSubOwnerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.VRFV2PlusExternalSubOwnerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusExternalSubOwnerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusExternalSubOwnerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.Owner(&_VRFV2PlusExternalSubOwnerExample.CallOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.Owner(&_VRFV2PlusExternalSubOwnerExample.CallOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusExternalSubOwnerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRandomWords(&_VRFV2PlusExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRandomWords(&_VRFV2PlusExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusExternalSubOwnerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRequestId(&_VRFV2PlusExternalSubOwnerExample.CallOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRequestId(&_VRFV2PlusExternalSubOwnerExample.CallOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.AcceptOwnership(&_VRFV2PlusExternalSubOwnerExample.TransactOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.AcceptOwnership(&_VRFV2PlusExternalSubOwnerExample.TransactOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "requestRandomWords", subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "setConfig", _vrfCoordinator, _subId)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SetConfig(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, _vrfCoordinator, _subId)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SetConfig(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, _vrfCoordinator, _subId)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.TransferOwnership(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, to)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.TransferOwnership(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, to)
}

type VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusExternalSubOwnerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusExternalSubOwnerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusExternalSubOwnerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested)
				if err := _VRFV2PlusExternalSubOwnerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested)
	if err := _VRFV2PlusExternalSubOwnerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusExternalSubOwnerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferred)
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

func (it *VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusExternalSubOwnerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusExternalSubOwnerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator{contract: _VRFV2PlusExternalSubOwnerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusExternalSubOwnerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusExternalSubOwnerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferred)
				if err := _VRFV2PlusExternalSubOwnerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusExternalSubOwnerExampleOwnershipTransferred)
	if err := _VRFV2PlusExternalSubOwnerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusExternalSubOwnerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusExternalSubOwnerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusExternalSubOwnerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusExternalSubOwnerExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusExternalSubOwnerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExample) Address() common.Address {
	return _VRFV2PlusExternalSubOwnerExample.address
}

type VRFV2PlusExternalSubOwnerExampleInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusExternalSubOwnerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusExternalSubOwnerExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
