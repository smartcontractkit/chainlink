// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_load_test_external_sub_owner

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

var VRFLoadTestExternalSubOwnerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b50604051610aa2380380610aa283398101604081905261002f916101ae565b6001600160601b0319606083901b1660805233806000816100975760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100c7576100c7816100e8565b5050506001600160601b0319606092831b811660a052911b1660c0526101e1565b6001600160a01b0381163314156101415760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161008e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146101a957600080fd5b919050565b600080604083850312156101c157600080fd5b6101ca83610192565b91506101d860208401610192565b90509250929050565b60805160601c60a05160601c60c05160601c61087c610226600039600060a701526000818161010b01526101f00152600081816102b3015261031b015261087c6000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806379ba50971161005b57806379ba50971461012d5780638da5cb5b14610135578063dc1670db14610153578063f2fde38b1461016a57600080fd5b8063096cb17b1461008d5780631b6b6d23146100a25780631fe543e3146100f35780633b2bcbf114610106575b600080fd5b6100a061009b36600461075a565b61017d565b005b6100c97f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100a061010136600461066b565b61029b565b6100c97f000000000000000000000000000000000000000000000000000000000000000081565b6100a061035b565b60005473ffffffffffffffffffffffffffffffffffffffff166100c9565b61015c60025481565b6040519081526020016100ea565b6100a0610178366004610615565b610458565b61018561046c565b60005b8161ffff168161ffff161015610294576040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810184905267ffffffffffffffff8616602482015261ffff8516604482015261c3506064820152600160848201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561024957600080fd5b505af115801561025d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102819190610652565b508061028c816107b6565b915050610188565b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461034d576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61035782826104ef565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610344565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61046061046c565b61046981610508565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146104ed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610344565b565b600280549060006104ff836107d8565b91905055505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610588576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610344565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803561ffff8116811461061057600080fd5b919050565b60006020828403121561062757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461064b57600080fd5b9392505050565b60006020828403121561066457600080fd5b5051919050565b6000806040838503121561067e57600080fd5b8235915060208084013567ffffffffffffffff8082111561069e57600080fd5b818601915086601f8301126106b257600080fd5b8135818111156106c4576106c4610840565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561070757610707610840565b604052828152858101935084860182860187018b101561072657600080fd5b600095505b8386101561074957803585526001959095019493860193860161072b565b508096505050505050509250929050565b6000806000806080858703121561077057600080fd5b843567ffffffffffffffff8116811461078857600080fd5b9350610796602086016105fe565b9250604085013591506107ab606086016105fe565b905092959194509250565b600061ffff808316818114156107ce576107ce610811565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561080a5761080a610811565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFLoadTestExternalSubOwnerABI = VRFLoadTestExternalSubOwnerMetaData.ABI

var VRFLoadTestExternalSubOwnerBin = VRFLoadTestExternalSubOwnerMetaData.Bin

func DeployVRFLoadTestExternalSubOwner(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address) (common.Address, *types.Transaction, *VRFLoadTestExternalSubOwner, error) {
	parsed, err := VRFLoadTestExternalSubOwnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFLoadTestExternalSubOwnerBin), backend, _vrfCoordinator, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFLoadTestExternalSubOwner{VRFLoadTestExternalSubOwnerCaller: VRFLoadTestExternalSubOwnerCaller{contract: contract}, VRFLoadTestExternalSubOwnerTransactor: VRFLoadTestExternalSubOwnerTransactor{contract: contract}, VRFLoadTestExternalSubOwnerFilterer: VRFLoadTestExternalSubOwnerFilterer{contract: contract}}, nil
}

type VRFLoadTestExternalSubOwner struct {
	address common.Address
	abi     abi.ABI
	VRFLoadTestExternalSubOwnerCaller
	VRFLoadTestExternalSubOwnerTransactor
	VRFLoadTestExternalSubOwnerFilterer
}

type VRFLoadTestExternalSubOwnerCaller struct {
	contract *bind.BoundContract
}

type VRFLoadTestExternalSubOwnerTransactor struct {
	contract *bind.BoundContract
}

type VRFLoadTestExternalSubOwnerFilterer struct {
	contract *bind.BoundContract
}

type VRFLoadTestExternalSubOwnerSession struct {
	Contract     *VRFLoadTestExternalSubOwner
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFLoadTestExternalSubOwnerCallerSession struct {
	Contract *VRFLoadTestExternalSubOwnerCaller
	CallOpts bind.CallOpts
}

type VRFLoadTestExternalSubOwnerTransactorSession struct {
	Contract     *VRFLoadTestExternalSubOwnerTransactor
	TransactOpts bind.TransactOpts
}

type VRFLoadTestExternalSubOwnerRaw struct {
	Contract *VRFLoadTestExternalSubOwner
}

type VRFLoadTestExternalSubOwnerCallerRaw struct {
	Contract *VRFLoadTestExternalSubOwnerCaller
}

type VRFLoadTestExternalSubOwnerTransactorRaw struct {
	Contract *VRFLoadTestExternalSubOwnerTransactor
}

func NewVRFLoadTestExternalSubOwner(address common.Address, backend bind.ContractBackend) (*VRFLoadTestExternalSubOwner, error) {
	abi, err := abi.JSON(strings.NewReader(VRFLoadTestExternalSubOwnerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFLoadTestExternalSubOwner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestExternalSubOwner{address: address, abi: abi, VRFLoadTestExternalSubOwnerCaller: VRFLoadTestExternalSubOwnerCaller{contract: contract}, VRFLoadTestExternalSubOwnerTransactor: VRFLoadTestExternalSubOwnerTransactor{contract: contract}, VRFLoadTestExternalSubOwnerFilterer: VRFLoadTestExternalSubOwnerFilterer{contract: contract}}, nil
}

func NewVRFLoadTestExternalSubOwnerCaller(address common.Address, caller bind.ContractCaller) (*VRFLoadTestExternalSubOwnerCaller, error) {
	contract, err := bindVRFLoadTestExternalSubOwner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestExternalSubOwnerCaller{contract: contract}, nil
}

func NewVRFLoadTestExternalSubOwnerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFLoadTestExternalSubOwnerTransactor, error) {
	contract, err := bindVRFLoadTestExternalSubOwner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestExternalSubOwnerTransactor{contract: contract}, nil
}

func NewVRFLoadTestExternalSubOwnerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFLoadTestExternalSubOwnerFilterer, error) {
	contract, err := bindVRFLoadTestExternalSubOwner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestExternalSubOwnerFilterer{contract: contract}, nil
}

func bindVRFLoadTestExternalSubOwner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFLoadTestExternalSubOwnerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFLoadTestExternalSubOwner.Contract.VRFLoadTestExternalSubOwnerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.VRFLoadTestExternalSubOwnerTransactor.contract.Transfer(opts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.VRFLoadTestExternalSubOwnerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFLoadTestExternalSubOwner.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.contract.Transfer(opts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.contract.Transact(opts, method, params...)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFLoadTestExternalSubOwner.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) COORDINATOR() (common.Address, error) {
	return _VRFLoadTestExternalSubOwner.Contract.COORDINATOR(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFLoadTestExternalSubOwner.Contract.COORDINATOR(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFLoadTestExternalSubOwner.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) LINK() (common.Address, error) {
	return _VRFLoadTestExternalSubOwner.Contract.LINK(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCallerSession) LINK() (common.Address, error) {
	return _VRFLoadTestExternalSubOwner.Contract.LINK(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFLoadTestExternalSubOwner.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) Owner() (common.Address, error) {
	return _VRFLoadTestExternalSubOwner.Contract.Owner(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCallerSession) Owner() (common.Address, error) {
	return _VRFLoadTestExternalSubOwner.Contract.Owner(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFLoadTestExternalSubOwner.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) SResponseCount() (*big.Int, error) {
	return _VRFLoadTestExternalSubOwner.Contract.SResponseCount(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFLoadTestExternalSubOwner.Contract.SResponseCount(&_VRFLoadTestExternalSubOwner.CallOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "acceptOwnership")
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.AcceptOwnership(&_VRFLoadTestExternalSubOwner.TransactOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.AcceptOwnership(&_VRFLoadTestExternalSubOwner.TransactOpts)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RawFulfillRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, requestId, randomWords)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RawFulfillRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, requestId, randomWords)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "requestRandomWords", _subId, _requestConfirmations, _keyHash, _requestCount)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) RequestRandomWords(_subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RequestRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, _subId, _requestConfirmations, _keyHash, _requestCount)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) RequestRandomWords(_subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RequestRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, _subId, _requestConfirmations, _keyHash, _requestCount)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.TransferOwnership(&_VRFLoadTestExternalSubOwner.TransactOpts, to)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.TransferOwnership(&_VRFLoadTestExternalSubOwner.TransactOpts, to)
}

type VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator struct {
	Event *VRFLoadTestExternalSubOwnerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFLoadTestExternalSubOwnerOwnershipTransferRequested)
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
		it.Event = new(VRFLoadTestExternalSubOwnerOwnershipTransferRequested)
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

func (it *VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFLoadTestExternalSubOwnerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFLoadTestExternalSubOwner.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator{contract: _VRFLoadTestExternalSubOwner.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFLoadTestExternalSubOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFLoadTestExternalSubOwner.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFLoadTestExternalSubOwnerOwnershipTransferRequested)
				if err := _VRFLoadTestExternalSubOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFLoadTestExternalSubOwnerOwnershipTransferRequested, error) {
	event := new(VRFLoadTestExternalSubOwnerOwnershipTransferRequested)
	if err := _VRFLoadTestExternalSubOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFLoadTestExternalSubOwnerOwnershipTransferredIterator struct {
	Event *VRFLoadTestExternalSubOwnerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFLoadTestExternalSubOwnerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFLoadTestExternalSubOwnerOwnershipTransferred)
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
		it.Event = new(VRFLoadTestExternalSubOwnerOwnershipTransferred)
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

func (it *VRFLoadTestExternalSubOwnerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFLoadTestExternalSubOwnerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFLoadTestExternalSubOwnerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFLoadTestExternalSubOwnerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFLoadTestExternalSubOwner.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestExternalSubOwnerOwnershipTransferredIterator{contract: _VRFLoadTestExternalSubOwner.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFLoadTestExternalSubOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFLoadTestExternalSubOwner.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFLoadTestExternalSubOwnerOwnershipTransferred)
				if err := _VRFLoadTestExternalSubOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFLoadTestExternalSubOwnerOwnershipTransferred, error) {
	event := new(VRFLoadTestExternalSubOwnerOwnershipTransferred)
	if err := _VRFLoadTestExternalSubOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwner) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFLoadTestExternalSubOwner.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFLoadTestExternalSubOwner.ParseOwnershipTransferRequested(log)
	case _VRFLoadTestExternalSubOwner.abi.Events["OwnershipTransferred"].ID:
		return _VRFLoadTestExternalSubOwner.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFLoadTestExternalSubOwnerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFLoadTestExternalSubOwnerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwner) Address() common.Address {
	return _VRFLoadTestExternalSubOwner.address
}

type VRFLoadTestExternalSubOwnerInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFLoadTestExternalSubOwnerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFLoadTestExternalSubOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFLoadTestExternalSubOwnerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFLoadTestExternalSubOwnerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFLoadTestExternalSubOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFLoadTestExternalSubOwnerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
