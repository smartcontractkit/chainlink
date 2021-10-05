// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_external_sub_example

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

var VRFConsumerExternalSubOwnerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"setSubscriptionID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161066638038061066683398101604081905261002f91610137565b606086811b6001600160601b0319166080908152600080546001600160a01b03998a166001600160a01b03199182161782556001805490911698909916979097179097556040805160a08101825296875263ffffffff9586166020880181905261ffff959095169087018190529290941693850184905293909401839052600280546001600160701b0319166801000000000000000090920261ffff60601b1916919091176c010000000000000000000000009094029390931763ffffffff60701b1916600160701b909102179091556003556101ad565b80516001600160a01b038116811461011e57600080fd5b919050565b805163ffffffff8116811461011e57600080fd5b60008060008060008060c0878903121561015057600080fd5b61015987610107565b955061016760208801610107565b945061017560408801610123565b9350606087015161ffff8116811461018c57600080fd5b925061019a60808801610123565b915060a087015190509295509295509295565b60805160601c6104956101d16000396000818160c4015261012c01526104956000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631edebf0b146100465780631fe543e314610091578063e0c86289146100a4575b600080fd5b61008f610054366004610428565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b005b61008f61009f366004610339565b6100ac565b61008f61016b565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461015d576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b61016782826102a8565b5050565b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff68010000000000000000830481166020850181905261ffff6c010000000000000000000000008504168587018190526e010000000000000000000000000000909404909116606085018190526003546080860181905260005496517f5d3b1d3000000000000000000000000000000000000000000000000000000000815260048101919091526024810193909352604483019390935260648201526084810191909152909173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561026a57600080fd5b505af115801561027e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102a29190610320565b60055550565b80516102bb9060049060208401906102c0565b505050565b8280548282559060005260206000209081019282156102fb579160200282015b828111156102fb5782518255916020019190600101906102e0565b5061030792915061030b565b5090565b5b80821115610307576000815560010161030c565b60006020828403121561033257600080fd5b5051919050565b6000806040838503121561034c57600080fd5b8235915060208084013567ffffffffffffffff8082111561036c57600080fd5b818601915086601f83011261038057600080fd5b81358181111561039257610392610459565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156103d5576103d5610459565b604052828152858101935084860182860187018b10156103f457600080fd5b600095505b838610156104175780358552600195909501949386019386016103f9565b508096505050505050509250929050565b60006020828403121561043a57600080fd5b813567ffffffffffffffff8116811461045257600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFConsumerExternalSubOwnerExampleABI = VRFConsumerExternalSubOwnerExampleMetaData.ABI

var VRFConsumerExternalSubOwnerExampleBin = VRFConsumerExternalSubOwnerExampleMetaData.Bin

func DeployVRFConsumerExternalSubOwnerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (common.Address, *types.Transaction, *VRFConsumerExternalSubOwnerExample, error) {
	parsed, err := VRFConsumerExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerExternalSubOwnerExampleBin), backend, vrfCoordinator, link, callbackGasLimit, requestConfirmations, numWords, keyHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerExternalSubOwnerExample{VRFConsumerExternalSubOwnerExampleCaller: VRFConsumerExternalSubOwnerExampleCaller{contract: contract}, VRFConsumerExternalSubOwnerExampleTransactor: VRFConsumerExternalSubOwnerExampleTransactor{contract: contract}, VRFConsumerExternalSubOwnerExampleFilterer: VRFConsumerExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

type VRFConsumerExternalSubOwnerExample struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerExternalSubOwnerExampleCaller
	VRFConsumerExternalSubOwnerExampleTransactor
	VRFConsumerExternalSubOwnerExampleFilterer
}

type VRFConsumerExternalSubOwnerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFConsumerExternalSubOwnerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFConsumerExternalSubOwnerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFConsumerExternalSubOwnerExampleSession struct {
	Contract     *VRFConsumerExternalSubOwnerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerExternalSubOwnerExampleCallerSession struct {
	Contract *VRFConsumerExternalSubOwnerExampleCaller
	CallOpts bind.CallOpts
}

type VRFConsumerExternalSubOwnerExampleTransactorSession struct {
	Contract     *VRFConsumerExternalSubOwnerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerExternalSubOwnerExampleRaw struct {
	Contract *VRFConsumerExternalSubOwnerExample
}

type VRFConsumerExternalSubOwnerExampleCallerRaw struct {
	Contract *VRFConsumerExternalSubOwnerExampleCaller
}

type VRFConsumerExternalSubOwnerExampleTransactorRaw struct {
	Contract *VRFConsumerExternalSubOwnerExampleTransactor
}

func NewVRFConsumerExternalSubOwnerExample(address common.Address, backend bind.ContractBackend) (*VRFConsumerExternalSubOwnerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumerExternalSubOwnerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerExternalSubOwnerExample{address: address, abi: abi, VRFConsumerExternalSubOwnerExampleCaller: VRFConsumerExternalSubOwnerExampleCaller{contract: contract}, VRFConsumerExternalSubOwnerExampleTransactor: VRFConsumerExternalSubOwnerExampleTransactor{contract: contract}, VRFConsumerExternalSubOwnerExampleFilterer: VRFConsumerExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

func NewVRFConsumerExternalSubOwnerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerExternalSubOwnerExampleCaller, error) {
	contract, err := bindVRFConsumerExternalSubOwnerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerExternalSubOwnerExampleCaller{contract: contract}, nil
}

func NewVRFConsumerExternalSubOwnerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerExternalSubOwnerExampleTransactor, error) {
	contract, err := bindVRFConsumerExternalSubOwnerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerExternalSubOwnerExampleTransactor{contract: contract}, nil
}

func NewVRFConsumerExternalSubOwnerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerExternalSubOwnerExampleFilterer, error) {
	contract, err := bindVRFConsumerExternalSubOwnerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerExternalSubOwnerExampleFilterer{contract: contract}, nil
}

func bindVRFConsumerExternalSubOwnerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerExternalSubOwnerExample.Contract.VRFConsumerExternalSubOwnerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.VRFConsumerExternalSubOwnerExampleTransactor.contract.Transfer(opts)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.VRFConsumerExternalSubOwnerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerExternalSubOwnerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.contract.Transfer(opts)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFConsumerExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFConsumerExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.contract.Transact(opts, "requestRandomWords")
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFConsumerExternalSubOwnerExample.TransactOpts)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactorSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFConsumerExternalSubOwnerExample.TransactOpts)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactor) SetSubscriptionID(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.contract.Transact(opts, "setSubscriptionID", subId)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleSession) SetSubscriptionID(subId uint64) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.SetSubscriptionID(&_VRFConsumerExternalSubOwnerExample.TransactOpts, subId)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExampleTransactorSession) SetSubscriptionID(subId uint64) (*types.Transaction, error) {
	return _VRFConsumerExternalSubOwnerExample.Contract.SetSubscriptionID(&_VRFConsumerExternalSubOwnerExample.TransactOpts, subId)
}

func (_VRFConsumerExternalSubOwnerExample *VRFConsumerExternalSubOwnerExample) Address() common.Address {
	return _VRFConsumerExternalSubOwnerExample.address
}

type VRFConsumerExternalSubOwnerExampleInterface interface {
	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	SetSubscriptionID(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	Address() common.Address
}
