// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_usdc_token_transmitter

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

var MockE2EUSDCTransmitterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_localDomain\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"localDomain\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextAvailableNonce\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"s_shouldSucceed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sendMessage\",\"inputs\":[{\"name\":\"destinationDomain\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"recipient\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"messageBody\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendMessageWithCaller\",\"inputs\":[{\"name\":\"destinationDomain\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"recipient\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"destinationCaller\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"messageBody\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setShouldSucceed\",\"inputs\":[{\"name\":\"shouldSucceed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x60e0346100c857601f610ab038819003918201601f19168301916001600160401b038311848410176100cd578084926060946040528339810103126100c857610047816100e3565b906040610056602083016100e3565b9101516001600160a01b03811692908390036100c85760805260a052600160ff19600054161760005560c0526040516109bb90816100f5823960805181818161011f0152818161063c0152610704015260a05181818161014901528181610418015261072e015260c051816105550152f35b600080fd5b634e487b7160e01b600052604160045260246000fd5b519063ffffffff821682036100c85756fe6080604052600436101561001257600080fd5b6000803560e01c80630ba469bc1461066057806354fd4d501461060157806357ecfd28146104c45780637a642935146104845780638371744e1461043c5780638d3638f4146103dd5780639e31ddb61461036e5763f7259a751461007557600080fd5b3461036b5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b576100ac61086b565b6044359160243560643567ffffffffffffffff8111610367576100d3903690600401610883565b9480156102e3576100e2610921565b94831561028557866020976101ed946094947fffffffff0000000000000000000000000000000000000000000000000000000097604051988996817f000000000000000000000000000000000000000000000000000000000000000060e01b168e890152817f000000000000000000000000000000000000000000000000000000000000000060e01b16602489015260e01b1660288701527fffffffffffffffff0000000000000000000000000000000000000000000000008b60c01b16602c87015233603487015260548601526074850152848401378101858382015203017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826108b1565b604051918483528151918286850152815b838110610271575050827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f846040948585977f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b0369901015201168101030190a167ffffffffffffffff60405191168152f35b8181018701518582016040015286016101fe565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f526563697069656e74206d757374206265206e6f6e7a65726f000000000000006044820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f44657374696e6174696f6e2063616c6c6572206d757374206265206e6f6e7a6560448201527f726f0000000000000000000000000000000000000000000000000000000000006064820152fd5b8280fd5b80fd5b503461036b5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b576004358015158091036103d95760ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00835416911617815580f35b5080fd5b503461036b57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b57602060405163ffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b503461036b57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b5767ffffffffffffffff6020915460081c16604051908152f35b503461036b57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b5760ff60209154166040519015158152f35b503461036b5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b5760043567ffffffffffffffff81116103d957610514903690600401610883565b60243567ffffffffffffffff81116105fd57610534903690600401610883565b505060b8116103d9578173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691823b156103d957604460a4918360405195869485937f40c10f19000000000000000000000000000000000000000000000000000000008552013560601c6004840152600160248401525af180156105f2579160ff91816020946105e2575b505054166040519015158152f35b6105eb916108b1565b38816105d4565b6040513d84823e3d90fd5b8380fd5b503461036b57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b57602060405163ffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b503461036b5760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261036b5761069861086b565b906024359060443567ffffffffffffffff81116103d9576106bd903690600401610883565b90926106c7610921565b93811561028557602095837fffffffff00000000000000000000000000000000000000000000000000000000946094936107d395604051978895817f000000000000000000000000000000000000000000000000000000000000000060e01b168d880152817f000000000000000000000000000000000000000000000000000000000000000060e01b16602488015260e01b1660288601527fffffffffffffffff0000000000000000000000000000000000000000000000008a60c01b16602c8601523360348601526054850152876074850152848401378101858382015203017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826108b1565b604051918483528151918286850152815b838110610857575050827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f846040948585977f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b0369901015201168101030190a167ffffffffffffffff60405191168152f35b8181018701518582016040015286016107e4565b6004359063ffffffff8216820361087e57565b600080fd5b9181601f8401121561087e5782359167ffffffffffffffff831161087e576020838186019501011161087e57565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176108f257604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60005467ffffffffffffffff8160081c16906001820167ffffffffffffffff811161097f5768ffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff9160081b1691161760005590565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c634300081a000a",
}

var MockE2EUSDCTransmitterABI = MockE2EUSDCTransmitterMetaData.ABI

var MockE2EUSDCTransmitterBin = MockE2EUSDCTransmitterMetaData.Bin

func DeployMockE2EUSDCTransmitter(auth *bind.TransactOpts, backend bind.ContractBackend, _version uint32, _localDomain uint32, token common.Address) (common.Address, *types.Transaction, *MockE2EUSDCTransmitter, error) {
	parsed, err := MockE2EUSDCTransmitterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockE2EUSDCTransmitterBin), backend, _version, _localDomain, token)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockE2EUSDCTransmitter{address: address, abi: *parsed, MockE2EUSDCTransmitterCaller: MockE2EUSDCTransmitterCaller{contract: contract}, MockE2EUSDCTransmitterTransactor: MockE2EUSDCTransmitterTransactor{contract: contract}, MockE2EUSDCTransmitterFilterer: MockE2EUSDCTransmitterFilterer{contract: contract}}, nil
}

type MockE2EUSDCTransmitter struct {
	address common.Address
	abi     abi.ABI
	MockE2EUSDCTransmitterCaller
	MockE2EUSDCTransmitterTransactor
	MockE2EUSDCTransmitterFilterer
}

type MockE2EUSDCTransmitterCaller struct {
	contract *bind.BoundContract
}

type MockE2EUSDCTransmitterTransactor struct {
	contract *bind.BoundContract
}

type MockE2EUSDCTransmitterFilterer struct {
	contract *bind.BoundContract
}

type MockE2EUSDCTransmitterSession struct {
	Contract     *MockE2EUSDCTransmitter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockE2EUSDCTransmitterCallerSession struct {
	Contract *MockE2EUSDCTransmitterCaller
	CallOpts bind.CallOpts
}

type MockE2EUSDCTransmitterTransactorSession struct {
	Contract     *MockE2EUSDCTransmitterTransactor
	TransactOpts bind.TransactOpts
}

type MockE2EUSDCTransmitterRaw struct {
	Contract *MockE2EUSDCTransmitter
}

type MockE2EUSDCTransmitterCallerRaw struct {
	Contract *MockE2EUSDCTransmitterCaller
}

type MockE2EUSDCTransmitterTransactorRaw struct {
	Contract *MockE2EUSDCTransmitterTransactor
}

func NewMockE2EUSDCTransmitter(address common.Address, backend bind.ContractBackend) (*MockE2EUSDCTransmitter, error) {
	abi, err := abi.JSON(strings.NewReader(MockE2EUSDCTransmitterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockE2EUSDCTransmitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTransmitter{address: address, abi: abi, MockE2EUSDCTransmitterCaller: MockE2EUSDCTransmitterCaller{contract: contract}, MockE2EUSDCTransmitterTransactor: MockE2EUSDCTransmitterTransactor{contract: contract}, MockE2EUSDCTransmitterFilterer: MockE2EUSDCTransmitterFilterer{contract: contract}}, nil
}

func NewMockE2EUSDCTransmitterCaller(address common.Address, caller bind.ContractCaller) (*MockE2EUSDCTransmitterCaller, error) {
	contract, err := bindMockE2EUSDCTransmitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTransmitterCaller{contract: contract}, nil
}

func NewMockE2EUSDCTransmitterTransactor(address common.Address, transactor bind.ContractTransactor) (*MockE2EUSDCTransmitterTransactor, error) {
	contract, err := bindMockE2EUSDCTransmitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTransmitterTransactor{contract: contract}, nil
}

func NewMockE2EUSDCTransmitterFilterer(address common.Address, filterer bind.ContractFilterer) (*MockE2EUSDCTransmitterFilterer, error) {
	contract, err := bindMockE2EUSDCTransmitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTransmitterFilterer{contract: contract}, nil
}

func bindMockE2EUSDCTransmitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockE2EUSDCTransmitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockE2EUSDCTransmitter.Contract.MockE2EUSDCTransmitterCaller.contract.Call(opts, result, method, params...)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.MockE2EUSDCTransmitterTransactor.contract.Transfer(opts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.MockE2EUSDCTransmitterTransactor.contract.Transact(opts, method, params...)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockE2EUSDCTransmitter.Contract.contract.Call(opts, result, method, params...)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.contract.Transfer(opts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.contract.Transact(opts, method, params...)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCaller) LocalDomain(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _MockE2EUSDCTransmitter.contract.Call(opts, &out, "localDomain")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) LocalDomain() (uint32, error) {
	return _MockE2EUSDCTransmitter.Contract.LocalDomain(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCallerSession) LocalDomain() (uint32, error) {
	return _MockE2EUSDCTransmitter.Contract.LocalDomain(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCaller) NextAvailableNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _MockE2EUSDCTransmitter.contract.Call(opts, &out, "nextAvailableNonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) NextAvailableNonce() (uint64, error) {
	return _MockE2EUSDCTransmitter.Contract.NextAvailableNonce(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCallerSession) NextAvailableNonce() (uint64, error) {
	return _MockE2EUSDCTransmitter.Contract.NextAvailableNonce(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCaller) SShouldSucceed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MockE2EUSDCTransmitter.contract.Call(opts, &out, "s_shouldSucceed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) SShouldSucceed() (bool, error) {
	return _MockE2EUSDCTransmitter.Contract.SShouldSucceed(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCallerSession) SShouldSucceed() (bool, error) {
	return _MockE2EUSDCTransmitter.Contract.SShouldSucceed(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCaller) Version(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _MockE2EUSDCTransmitter.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) Version() (uint32, error) {
	return _MockE2EUSDCTransmitter.Contract.Version(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterCallerSession) Version() (uint32, error) {
	return _MockE2EUSDCTransmitter.Contract.Version(&_MockE2EUSDCTransmitter.CallOpts)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactor) ReceiveMessage(opts *bind.TransactOpts, message []byte, arg1 []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.contract.Transact(opts, "receiveMessage", message, arg1)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) ReceiveMessage(message []byte, arg1 []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.ReceiveMessage(&_MockE2EUSDCTransmitter.TransactOpts, message, arg1)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactorSession) ReceiveMessage(message []byte, arg1 []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.ReceiveMessage(&_MockE2EUSDCTransmitter.TransactOpts, message, arg1)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactor) SendMessage(opts *bind.TransactOpts, destinationDomain uint32, recipient [32]byte, messageBody []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.contract.Transact(opts, "sendMessage", destinationDomain, recipient, messageBody)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) SendMessage(destinationDomain uint32, recipient [32]byte, messageBody []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.SendMessage(&_MockE2EUSDCTransmitter.TransactOpts, destinationDomain, recipient, messageBody)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactorSession) SendMessage(destinationDomain uint32, recipient [32]byte, messageBody []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.SendMessage(&_MockE2EUSDCTransmitter.TransactOpts, destinationDomain, recipient, messageBody)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactor) SendMessageWithCaller(opts *bind.TransactOpts, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, messageBody []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.contract.Transact(opts, "sendMessageWithCaller", destinationDomain, recipient, destinationCaller, messageBody)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) SendMessageWithCaller(destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, messageBody []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.SendMessageWithCaller(&_MockE2EUSDCTransmitter.TransactOpts, destinationDomain, recipient, destinationCaller, messageBody)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactorSession) SendMessageWithCaller(destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, messageBody []byte) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.SendMessageWithCaller(&_MockE2EUSDCTransmitter.TransactOpts, destinationDomain, recipient, destinationCaller, messageBody)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactor) SetShouldSucceed(opts *bind.TransactOpts, shouldSucceed bool) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.contract.Transact(opts, "setShouldSucceed", shouldSucceed)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterSession) SetShouldSucceed(shouldSucceed bool) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.SetShouldSucceed(&_MockE2EUSDCTransmitter.TransactOpts, shouldSucceed)
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterTransactorSession) SetShouldSucceed(shouldSucceed bool) (*types.Transaction, error) {
	return _MockE2EUSDCTransmitter.Contract.SetShouldSucceed(&_MockE2EUSDCTransmitter.TransactOpts, shouldSucceed)
}

type MockE2EUSDCTransmitterMessageSentIterator struct {
	Event *MockE2EUSDCTransmitterMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockE2EUSDCTransmitterMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockE2EUSDCTransmitterMessageSent)
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
		it.Event = new(MockE2EUSDCTransmitterMessageSent)
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

func (it *MockE2EUSDCTransmitterMessageSentIterator) Error() error {
	return it.fail
}

func (it *MockE2EUSDCTransmitterMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockE2EUSDCTransmitterMessageSent struct {
	Message []byte
	Raw     types.Log
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterFilterer) FilterMessageSent(opts *bind.FilterOpts) (*MockE2EUSDCTransmitterMessageSentIterator, error) {

	logs, sub, err := _MockE2EUSDCTransmitter.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTransmitterMessageSentIterator{contract: _MockE2EUSDCTransmitter.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *MockE2EUSDCTransmitterMessageSent) (event.Subscription, error) {

	logs, sub, err := _MockE2EUSDCTransmitter.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockE2EUSDCTransmitterMessageSent)
				if err := _MockE2EUSDCTransmitter.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitterFilterer) ParseMessageSent(log types.Log) (*MockE2EUSDCTransmitterMessageSent, error) {
	event := new(MockE2EUSDCTransmitterMessageSent)
	if err := _MockE2EUSDCTransmitter.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockE2EUSDCTransmitter.abi.Events["MessageSent"].ID:
		return _MockE2EUSDCTransmitter.ParseMessageSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockE2EUSDCTransmitterMessageSent) Topic() common.Hash {
	return common.HexToHash("0x8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036")
}

func (_MockE2EUSDCTransmitter *MockE2EUSDCTransmitter) Address() common.Address {
	return _MockE2EUSDCTransmitter.address
}

type MockE2EUSDCTransmitterInterface interface {
	LocalDomain(opts *bind.CallOpts) (uint32, error)

	NextAvailableNonce(opts *bind.CallOpts) (uint64, error)

	SShouldSucceed(opts *bind.CallOpts) (bool, error)

	Version(opts *bind.CallOpts) (uint32, error)

	ReceiveMessage(opts *bind.TransactOpts, message []byte, arg1 []byte) (*types.Transaction, error)

	SendMessage(opts *bind.TransactOpts, destinationDomain uint32, recipient [32]byte, messageBody []byte) (*types.Transaction, error)

	SendMessageWithCaller(opts *bind.TransactOpts, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, messageBody []byte) (*types.Transaction, error)

	SetShouldSucceed(opts *bind.TransactOpts, shouldSucceed bool) (*types.Transaction, error)

	FilterMessageSent(opts *bind.FilterOpts) (*MockE2EUSDCTransmitterMessageSentIterator, error)

	WatchMessageSent(opts *bind.WatchOpts, sink chan<- *MockE2EUSDCTransmitterMessageSent) (event.Subscription, error)

	ParseMessageSent(log types.Log) (*MockE2EUSDCTransmitterMessageSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
