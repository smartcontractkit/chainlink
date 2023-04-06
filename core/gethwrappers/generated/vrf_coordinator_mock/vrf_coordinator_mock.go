// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_mock

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

var VRFCoordinatorMockMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequest\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumerContract\",\"type\":\"address\"}],\"name\":\"callBackWithRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161057938038061057983398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b6104e6806100936000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631b6b6d2314610046578063a4c0ed361461008f578063cf55fe97146100a4575b600080fd5b6000546100669073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100a261009d36600461032d565b6100b7565b005b6100a26100b236600461043a565b6101b1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461013d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e0000000000000000000000000060448201526064015b60405180910390fd5b600080828060200190518101906101549190610416565b9150915080828673ffffffffffffffffffffffffffffffffffffffff167fb6a11357fce9fae0b59dcc6e5e4bf50803daf2b17d3b80739767e0c4fdacb444876040516101a291815260200190565b60405180910390a45050505050565b604080516024810185905260448082018590528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f94985ddd00000000000000000000000000000000000000000000000000000000179052600090620324b0805a101561028f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6e6f7420656e6f7567682067617320666f7220636f6e73756d657200000000006044820152606401610134565b60008473ffffffffffffffffffffffffffffffffffffffff16836040516102b6919061046f565b6000604051808303816000865af19150503d80600081146102f3576040519150601f19603f3d011682016040523d82523d6000602084013e6102f8565b606091505b50505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461032857600080fd5b919050565b60008060006060848603121561034257600080fd5b61034b84610304565b925060208401359150604084013567ffffffffffffffff8082111561036f57600080fd5b818601915086601f83011261038357600080fd5b813581811115610395576103956104aa565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156103db576103db6104aa565b816040528281528960208487010111156103f457600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b6000806040838503121561042957600080fd5b505080516020909101519092909150565b60008060006060848603121561044f57600080fd5b833592506020840135915061046660408501610304565b90509250925092565b6000825160005b818110156104905760208186018101518583015201610476565b8181111561049f576000828501525b509190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFCoordinatorMockABI = VRFCoordinatorMockMetaData.ABI

var VRFCoordinatorMockBin = VRFCoordinatorMockMetaData.Bin

func DeployVRFCoordinatorMock(auth *bind.TransactOpts, backend bind.ContractBackend, linkAddress common.Address) (common.Address, *types.Transaction, *VRFCoordinatorMock, error) {
	parsed, err := VRFCoordinatorMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorMockBin), backend, linkAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorMock{VRFCoordinatorMockCaller: VRFCoordinatorMockCaller{contract: contract}, VRFCoordinatorMockTransactor: VRFCoordinatorMockTransactor{contract: contract}, VRFCoordinatorMockFilterer: VRFCoordinatorMockFilterer{contract: contract}}, nil
}

type VRFCoordinatorMock struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorMockCaller
	VRFCoordinatorMockTransactor
	VRFCoordinatorMockFilterer
}

type VRFCoordinatorMockCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorMockTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorMockFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorMockSession struct {
	Contract     *VRFCoordinatorMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorMockCallerSession struct {
	Contract *VRFCoordinatorMockCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorMockTransactorSession struct {
	Contract     *VRFCoordinatorMockTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorMockRaw struct {
	Contract *VRFCoordinatorMock
}

type VRFCoordinatorMockCallerRaw struct {
	Contract *VRFCoordinatorMockCaller
}

type VRFCoordinatorMockTransactorRaw struct {
	Contract *VRFCoordinatorMockTransactor
}

func NewVRFCoordinatorMock(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorMock, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorMock{address: address, abi: abi, VRFCoordinatorMockCaller: VRFCoordinatorMockCaller{contract: contract}, VRFCoordinatorMockTransactor: VRFCoordinatorMockTransactor{contract: contract}, VRFCoordinatorMockFilterer: VRFCoordinatorMockFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorMockCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorMockCaller, error) {
	contract, err := bindVRFCoordinatorMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorMockCaller{contract: contract}, nil
}

func NewVRFCoordinatorMockTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorMockTransactor, error) {
	contract, err := bindVRFCoordinatorMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorMockTransactor{contract: contract}, nil
}

func NewVRFCoordinatorMockFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorMockFilterer, error) {
	contract, err := bindVRFCoordinatorMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorMockFilterer{contract: contract}, nil
}

func bindVRFCoordinatorMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorMock *VRFCoordinatorMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorMock.Contract.VRFCoordinatorMockCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.VRFCoordinatorMockTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.VRFCoordinatorMockTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorMock.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorMock.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorMock *VRFCoordinatorMockSession) LINK() (common.Address, error) {
	return _VRFCoordinatorMock.Contract.LINK(&_VRFCoordinatorMock.CallOpts)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockCallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorMock.Contract.LINK(&_VRFCoordinatorMock.CallOpts)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockTransactor) CallBackWithRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int, consumerContract common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorMock.contract.Transact(opts, "callBackWithRandomness", requestId, randomness, consumerContract)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockSession) CallBackWithRandomness(requestId [32]byte, randomness *big.Int, consumerContract common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.CallBackWithRandomness(&_VRFCoordinatorMock.TransactOpts, requestId, randomness, consumerContract)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockTransactorSession) CallBackWithRandomness(requestId [32]byte, randomness *big.Int, consumerContract common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.CallBackWithRandomness(&_VRFCoordinatorMock.TransactOpts, requestId, randomness, consumerContract)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorMock.contract.Transact(opts, "onTokenTransfer", sender, fee, _data)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockSession) OnTokenTransfer(sender common.Address, fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.OnTokenTransfer(&_VRFCoordinatorMock.TransactOpts, sender, fee, _data)
}

func (_VRFCoordinatorMock *VRFCoordinatorMockTransactorSession) OnTokenTransfer(sender common.Address, fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorMock.Contract.OnTokenTransfer(&_VRFCoordinatorMock.TransactOpts, sender, fee, _data)
}

type VRFCoordinatorMockRandomnessRequestIterator struct {
	Event *VRFCoordinatorMockRandomnessRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorMockRandomnessRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorMockRandomnessRequest)
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
		it.Event = new(VRFCoordinatorMockRandomnessRequest)
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

func (it *VRFCoordinatorMockRandomnessRequestIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorMockRandomnessRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorMockRandomnessRequest struct {
	Sender  common.Address
	KeyHash [32]byte
	Seed    *big.Int
	Fee     *big.Int
	Raw     types.Log
}

func (_VRFCoordinatorMock *VRFCoordinatorMockFilterer) FilterRandomnessRequest(opts *bind.FilterOpts, sender []common.Address, keyHash [][32]byte, seed []*big.Int) (*VRFCoordinatorMockRandomnessRequestIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}
	var seedRule []interface{}
	for _, seedItem := range seed {
		seedRule = append(seedRule, seedItem)
	}

	logs, sub, err := _VRFCoordinatorMock.contract.FilterLogs(opts, "RandomnessRequest", senderRule, keyHashRule, seedRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorMockRandomnessRequestIterator{contract: _VRFCoordinatorMock.contract, event: "RandomnessRequest", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorMock *VRFCoordinatorMockFilterer) WatchRandomnessRequest(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorMockRandomnessRequest, sender []common.Address, keyHash [][32]byte, seed []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}
	var seedRule []interface{}
	for _, seedItem := range seed {
		seedRule = append(seedRule, seedItem)
	}

	logs, sub, err := _VRFCoordinatorMock.contract.WatchLogs(opts, "RandomnessRequest", senderRule, keyHashRule, seedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorMockRandomnessRequest)
				if err := _VRFCoordinatorMock.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
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

func (_VRFCoordinatorMock *VRFCoordinatorMockFilterer) ParseRandomnessRequest(log types.Log) (*VRFCoordinatorMockRandomnessRequest, error) {
	event := new(VRFCoordinatorMockRandomnessRequest)
	if err := _VRFCoordinatorMock.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFCoordinatorMock *VRFCoordinatorMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorMock.abi.Events["RandomnessRequest"].ID:
		return _VRFCoordinatorMock.ParseRandomnessRequest(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorMockRandomnessRequest) Topic() common.Hash {
	return common.HexToHash("0xb6a11357fce9fae0b59dcc6e5e4bf50803daf2b17d3b80739767e0c4fdacb444")
}

func (_VRFCoordinatorMock *VRFCoordinatorMock) Address() common.Address {
	return _VRFCoordinatorMock.address
}

type VRFCoordinatorMockInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	CallBackWithRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int, consumerContract common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, fee *big.Int, _data []byte) (*types.Transaction, error)

	FilterRandomnessRequest(opts *bind.FilterOpts, sender []common.Address, keyHash [][32]byte, seed []*big.Int) (*VRFCoordinatorMockRandomnessRequestIterator, error)

	WatchRandomnessRequest(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorMockRandomnessRequest, sender []common.Address, keyHash [][32]byte, seed []*big.Int) (event.Subscription, error)

	ParseRandomnessRequest(log types.Log) (*VRFCoordinatorMockRandomnessRequest, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
