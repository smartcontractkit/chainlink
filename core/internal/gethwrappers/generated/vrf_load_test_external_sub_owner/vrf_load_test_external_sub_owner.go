// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_load_test_external_sub_owner

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405161076d38038061076d83398101604081905261002f91610080565b6001600160601b0319606092831b8116608081905260a052911b1660c052600180546001600160a01b031916331790556100b3565b80516001600160a01b038116811461007b57600080fd5b919050565b6000806040838503121561009357600080fd5b61009c83610064565b91506100aa60208401610064565b90509250929050565b60805160601c60a05160601c60c05160601c6106766100f76000396000609101526000818160f501526101d501526000818161029a015261030201526106766000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c80633b2bcbf1116100505780633b2bcbf1146100f0578063dc1670db14610117578063f2fde38b1461012e57600080fd5b80630ef284da146100775780631b6b6d231461008c5780631fe543e3146100dd575b600080fd5b61008a610085366004610534565b610141565b005b6100b37f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61008a6100eb366004610445565b610282565b6100b37f000000000000000000000000000000000000000000000000000000000000000081565b61012060005481565b6040519081526020016100d4565b61008a61013c3660046103ef565b610341565b60015473ffffffffffffffffffffffffffffffffffffffff16331461016557600080fd5b60005b8161ffff168161ffff161015610279576040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810184905267ffffffffffffffff8816602482015261ffff8616604482015263ffffffff8088166064830152851660848201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561022e57600080fd5b505af1158015610242573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610266919061042c565b5080610271816105b0565b915050610168565b50505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610333576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b61033d82826103ac565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461036557600080fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805490806103bb836105d2565b91905055505050565b803561ffff811681146103d657600080fd5b919050565b803563ffffffff811681146103d657600080fd5b60006020828403121561040157600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461042557600080fd5b9392505050565b60006020828403121561043e57600080fd5b5051919050565b6000806040838503121561045857600080fd5b8235915060208084013567ffffffffffffffff8082111561047857600080fd5b818601915086601f83011261048c57600080fd5b81358181111561049e5761049e61063a565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156104e1576104e161063a565b604052828152858101935084860182860187018b101561050057600080fd5b600095505b83861015610523578035855260019590950194938601938601610505565b508096505050505050509250929050565b60008060008060008060c0878903121561054d57600080fd5b863567ffffffffffffffff8116811461056557600080fd5b9550610573602088016103db565b9450610581604088016103c4565b935061058f606088016103db565b9250608087013591506105a460a088016103c4565b90509295509295509295565b600061ffff808316818114156105c8576105c861060b565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156106045761060461060b565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RawFulfillRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, requestId, randomWords)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RawFulfillRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, requestId, randomWords)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "requestRandomWords", _subId, _callbackGasLimit, _requestConfirmations, _numWords, _keyHash, _requestCount)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) RequestRandomWords(_subId uint64, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RequestRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, _subId, _callbackGasLimit, _requestConfirmations, _numWords, _keyHash, _requestCount)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) RequestRandomWords(_subId uint64, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.RequestRandomWords(&_VRFLoadTestExternalSubOwner.TransactOpts, _subId, _callbackGasLimit, _requestConfirmations, _numWords, _keyHash, _requestCount)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.TransferOwnership(&_VRFLoadTestExternalSubOwner.TransactOpts, newOwner)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwnerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFLoadTestExternalSubOwner.Contract.TransferOwnership(&_VRFLoadTestExternalSubOwner.TransactOpts, newOwner)
}

func (_VRFLoadTestExternalSubOwner *VRFLoadTestExternalSubOwner) Address() common.Address {
	return _VRFLoadTestExternalSubOwner.address
}

type VRFLoadTestExternalSubOwnerInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	Address() common.Address
}
