// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_inbox

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

var ArbitrumInboxMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"messageNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"InboxMessageDelivered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"messageNum\",\"type\":\"uint256\"}],\"name\":\"InboxMessageDeliveredFromOrigin\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"allowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"contractIBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"}],\"name\":\"calculateRetryableSubmissionFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProxyAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIBridge\",\"name\":\"_bridge\",\"type\":\"address\"},{\"internalType\":\"contractISequencerInbox\",\"name\":\"_sequencerInbox\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"isAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxDataSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"sendContractTransaction\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"messageData\",\"type\":\"bytes\"}],\"name\":\"sendL2Message\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"messageData\",\"type\":\"bytes\"}],\"name\":\"sendL2MessageFromOrigin\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"sendUnsignedTransaction\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sequencerInbox\",\"outputs\":[{\"internalType\":\"contractISequencerInbox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"user\",\"type\":\"address[]\"},{\"internalType\":\"bool[]\",\"name\":\"val\",\"type\":\"bool[]\"}],\"name\":\"setAllowList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_allowListEnabled\",\"type\":\"bool\"}],\"name\":\"setAllowListEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var ArbitrumInboxABI = ArbitrumInboxMetaData.ABI

type ArbitrumInbox struct {
	address common.Address
	abi     abi.ABI
	ArbitrumInboxCaller
	ArbitrumInboxTransactor
	ArbitrumInboxFilterer
}

type ArbitrumInboxCaller struct {
	contract *bind.BoundContract
}

type ArbitrumInboxTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumInboxFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumInboxSession struct {
	Contract     *ArbitrumInbox
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumInboxCallerSession struct {
	Contract *ArbitrumInboxCaller
	CallOpts bind.CallOpts
}

type ArbitrumInboxTransactorSession struct {
	Contract     *ArbitrumInboxTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumInboxRaw struct {
	Contract *ArbitrumInbox
}

type ArbitrumInboxCallerRaw struct {
	Contract *ArbitrumInboxCaller
}

type ArbitrumInboxTransactorRaw struct {
	Contract *ArbitrumInboxTransactor
}

func NewArbitrumInbox(address common.Address, backend bind.ContractBackend) (*ArbitrumInbox, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumInboxABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumInbox(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumInbox{address: address, abi: abi, ArbitrumInboxCaller: ArbitrumInboxCaller{contract: contract}, ArbitrumInboxTransactor: ArbitrumInboxTransactor{contract: contract}, ArbitrumInboxFilterer: ArbitrumInboxFilterer{contract: contract}}, nil
}

func NewArbitrumInboxCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumInboxCaller, error) {
	contract, err := bindArbitrumInbox(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumInboxCaller{contract: contract}, nil
}

func NewArbitrumInboxTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumInboxTransactor, error) {
	contract, err := bindArbitrumInbox(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumInboxTransactor{contract: contract}, nil
}

func NewArbitrumInboxFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumInboxFilterer, error) {
	contract, err := bindArbitrumInbox(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumInboxFilterer{contract: contract}, nil
}

func bindArbitrumInbox(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrumInboxMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbitrumInbox *ArbitrumInboxRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumInbox.Contract.ArbitrumInboxCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumInbox *ArbitrumInboxRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.ArbitrumInboxTransactor.contract.Transfer(opts)
}

func (_ArbitrumInbox *ArbitrumInboxRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.ArbitrumInboxTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumInbox *ArbitrumInboxCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumInbox.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.contract.Transfer(opts)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) AllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "allowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) AllowListEnabled() (bool, error) {
	return _ArbitrumInbox.Contract.AllowListEnabled(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) AllowListEnabled() (bool, error) {
	return _ArbitrumInbox.Contract.AllowListEnabled(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) Bridge() (common.Address, error) {
	return _ArbitrumInbox.Contract.Bridge(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) Bridge() (common.Address, error) {
	return _ArbitrumInbox.Contract.Bridge(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) CalculateRetryableSubmissionFee(opts *bind.CallOpts, dataLength *big.Int, baseFee *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "calculateRetryableSubmissionFee", dataLength, baseFee)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) CalculateRetryableSubmissionFee(dataLength *big.Int, baseFee *big.Int) (*big.Int, error) {
	return _ArbitrumInbox.Contract.CalculateRetryableSubmissionFee(&_ArbitrumInbox.CallOpts, dataLength, baseFee)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) CalculateRetryableSubmissionFee(dataLength *big.Int, baseFee *big.Int) (*big.Int, error) {
	return _ArbitrumInbox.Contract.CalculateRetryableSubmissionFee(&_ArbitrumInbox.CallOpts, dataLength, baseFee)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) GetProxyAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "getProxyAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) GetProxyAdmin() (common.Address, error) {
	return _ArbitrumInbox.Contract.GetProxyAdmin(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) GetProxyAdmin() (common.Address, error) {
	return _ArbitrumInbox.Contract.GetProxyAdmin(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) IsAllowed(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "isAllowed", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) IsAllowed(user common.Address) (bool, error) {
	return _ArbitrumInbox.Contract.IsAllowed(&_ArbitrumInbox.CallOpts, user)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) IsAllowed(user common.Address) (bool, error) {
	return _ArbitrumInbox.Contract.IsAllowed(&_ArbitrumInbox.CallOpts, user)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) MaxDataSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "maxDataSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) MaxDataSize() (*big.Int, error) {
	return _ArbitrumInbox.Contract.MaxDataSize(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) MaxDataSize() (*big.Int, error) {
	return _ArbitrumInbox.Contract.MaxDataSize(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCaller) SequencerInbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumInbox.contract.Call(opts, &out, "sequencerInbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumInbox *ArbitrumInboxSession) SequencerInbox() (common.Address, error) {
	return _ArbitrumInbox.Contract.SequencerInbox(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxCallerSession) SequencerInbox() (common.Address, error) {
	return _ArbitrumInbox.Contract.SequencerInbox(&_ArbitrumInbox.CallOpts)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) Initialize(opts *bind.TransactOpts, _bridge common.Address, _sequencerInbox common.Address) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "initialize", _bridge, _sequencerInbox)
}

func (_ArbitrumInbox *ArbitrumInboxSession) Initialize(_bridge common.Address, _sequencerInbox common.Address) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.Initialize(&_ArbitrumInbox.TransactOpts, _bridge, _sequencerInbox)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) Initialize(_bridge common.Address, _sequencerInbox common.Address) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.Initialize(&_ArbitrumInbox.TransactOpts, _bridge, _sequencerInbox)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "pause")
}

func (_ArbitrumInbox *ArbitrumInboxSession) Pause() (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.Pause(&_ArbitrumInbox.TransactOpts)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) Pause() (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.Pause(&_ArbitrumInbox.TransactOpts)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) SendContractTransaction(opts *bind.TransactOpts, gasLimit *big.Int, maxFeePerGas *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "sendContractTransaction", gasLimit, maxFeePerGas, to, value, data)
}

func (_ArbitrumInbox *ArbitrumInboxSession) SendContractTransaction(gasLimit *big.Int, maxFeePerGas *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendContractTransaction(&_ArbitrumInbox.TransactOpts, gasLimit, maxFeePerGas, to, value, data)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) SendContractTransaction(gasLimit *big.Int, maxFeePerGas *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendContractTransaction(&_ArbitrumInbox.TransactOpts, gasLimit, maxFeePerGas, to, value, data)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) SendL2Message(opts *bind.TransactOpts, messageData []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "sendL2Message", messageData)
}

func (_ArbitrumInbox *ArbitrumInboxSession) SendL2Message(messageData []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendL2Message(&_ArbitrumInbox.TransactOpts, messageData)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) SendL2Message(messageData []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendL2Message(&_ArbitrumInbox.TransactOpts, messageData)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) SendL2MessageFromOrigin(opts *bind.TransactOpts, messageData []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "sendL2MessageFromOrigin", messageData)
}

func (_ArbitrumInbox *ArbitrumInboxSession) SendL2MessageFromOrigin(messageData []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendL2MessageFromOrigin(&_ArbitrumInbox.TransactOpts, messageData)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) SendL2MessageFromOrigin(messageData []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendL2MessageFromOrigin(&_ArbitrumInbox.TransactOpts, messageData)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) SendUnsignedTransaction(opts *bind.TransactOpts, gasLimit *big.Int, maxFeePerGas *big.Int, nonce *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "sendUnsignedTransaction", gasLimit, maxFeePerGas, nonce, to, value, data)
}

func (_ArbitrumInbox *ArbitrumInboxSession) SendUnsignedTransaction(gasLimit *big.Int, maxFeePerGas *big.Int, nonce *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendUnsignedTransaction(&_ArbitrumInbox.TransactOpts, gasLimit, maxFeePerGas, nonce, to, value, data)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) SendUnsignedTransaction(gasLimit *big.Int, maxFeePerGas *big.Int, nonce *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SendUnsignedTransaction(&_ArbitrumInbox.TransactOpts, gasLimit, maxFeePerGas, nonce, to, value, data)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) SetAllowList(opts *bind.TransactOpts, user []common.Address, val []bool) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "setAllowList", user, val)
}

func (_ArbitrumInbox *ArbitrumInboxSession) SetAllowList(user []common.Address, val []bool) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SetAllowList(&_ArbitrumInbox.TransactOpts, user, val)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) SetAllowList(user []common.Address, val []bool) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SetAllowList(&_ArbitrumInbox.TransactOpts, user, val)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) SetAllowListEnabled(opts *bind.TransactOpts, _allowListEnabled bool) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "setAllowListEnabled", _allowListEnabled)
}

func (_ArbitrumInbox *ArbitrumInboxSession) SetAllowListEnabled(_allowListEnabled bool) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SetAllowListEnabled(&_ArbitrumInbox.TransactOpts, _allowListEnabled)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) SetAllowListEnabled(_allowListEnabled bool) (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.SetAllowListEnabled(&_ArbitrumInbox.TransactOpts, _allowListEnabled)
}

func (_ArbitrumInbox *ArbitrumInboxTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumInbox.contract.Transact(opts, "unpause")
}

func (_ArbitrumInbox *ArbitrumInboxSession) Unpause() (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.Unpause(&_ArbitrumInbox.TransactOpts)
}

func (_ArbitrumInbox *ArbitrumInboxTransactorSession) Unpause() (*types.Transaction, error) {
	return _ArbitrumInbox.Contract.Unpause(&_ArbitrumInbox.TransactOpts)
}

type ArbitrumInboxInboxMessageDeliveredIterator struct {
	Event *ArbitrumInboxInboxMessageDelivered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbitrumInboxInboxMessageDeliveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbitrumInboxInboxMessageDelivered)
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
		it.Event = new(ArbitrumInboxInboxMessageDelivered)
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

func (it *ArbitrumInboxInboxMessageDeliveredIterator) Error() error {
	return it.fail
}

func (it *ArbitrumInboxInboxMessageDeliveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbitrumInboxInboxMessageDelivered struct {
	MessageNum *big.Int
	Data       []byte
	Raw        types.Log
}

func (_ArbitrumInbox *ArbitrumInboxFilterer) FilterInboxMessageDelivered(opts *bind.FilterOpts, messageNum []*big.Int) (*ArbitrumInboxInboxMessageDeliveredIterator, error) {

	var messageNumRule []interface{}
	for _, messageNumItem := range messageNum {
		messageNumRule = append(messageNumRule, messageNumItem)
	}

	logs, sub, err := _ArbitrumInbox.contract.FilterLogs(opts, "InboxMessageDelivered", messageNumRule)
	if err != nil {
		return nil, err
	}
	return &ArbitrumInboxInboxMessageDeliveredIterator{contract: _ArbitrumInbox.contract, event: "InboxMessageDelivered", logs: logs, sub: sub}, nil
}

func (_ArbitrumInbox *ArbitrumInboxFilterer) WatchInboxMessageDelivered(opts *bind.WatchOpts, sink chan<- *ArbitrumInboxInboxMessageDelivered, messageNum []*big.Int) (event.Subscription, error) {

	var messageNumRule []interface{}
	for _, messageNumItem := range messageNum {
		messageNumRule = append(messageNumRule, messageNumItem)
	}

	logs, sub, err := _ArbitrumInbox.contract.WatchLogs(opts, "InboxMessageDelivered", messageNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbitrumInboxInboxMessageDelivered)
				if err := _ArbitrumInbox.contract.UnpackLog(event, "InboxMessageDelivered", log); err != nil {
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

func (_ArbitrumInbox *ArbitrumInboxFilterer) ParseInboxMessageDelivered(log types.Log) (*ArbitrumInboxInboxMessageDelivered, error) {
	event := new(ArbitrumInboxInboxMessageDelivered)
	if err := _ArbitrumInbox.contract.UnpackLog(event, "InboxMessageDelivered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbitrumInboxInboxMessageDeliveredFromOriginIterator struct {
	Event *ArbitrumInboxInboxMessageDeliveredFromOrigin

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbitrumInboxInboxMessageDeliveredFromOriginIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbitrumInboxInboxMessageDeliveredFromOrigin)
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
		it.Event = new(ArbitrumInboxInboxMessageDeliveredFromOrigin)
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

func (it *ArbitrumInboxInboxMessageDeliveredFromOriginIterator) Error() error {
	return it.fail
}

func (it *ArbitrumInboxInboxMessageDeliveredFromOriginIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbitrumInboxInboxMessageDeliveredFromOrigin struct {
	MessageNum *big.Int
	Raw        types.Log
}

func (_ArbitrumInbox *ArbitrumInboxFilterer) FilterInboxMessageDeliveredFromOrigin(opts *bind.FilterOpts, messageNum []*big.Int) (*ArbitrumInboxInboxMessageDeliveredFromOriginIterator, error) {

	var messageNumRule []interface{}
	for _, messageNumItem := range messageNum {
		messageNumRule = append(messageNumRule, messageNumItem)
	}

	logs, sub, err := _ArbitrumInbox.contract.FilterLogs(opts, "InboxMessageDeliveredFromOrigin", messageNumRule)
	if err != nil {
		return nil, err
	}
	return &ArbitrumInboxInboxMessageDeliveredFromOriginIterator{contract: _ArbitrumInbox.contract, event: "InboxMessageDeliveredFromOrigin", logs: logs, sub: sub}, nil
}

func (_ArbitrumInbox *ArbitrumInboxFilterer) WatchInboxMessageDeliveredFromOrigin(opts *bind.WatchOpts, sink chan<- *ArbitrumInboxInboxMessageDeliveredFromOrigin, messageNum []*big.Int) (event.Subscription, error) {

	var messageNumRule []interface{}
	for _, messageNumItem := range messageNum {
		messageNumRule = append(messageNumRule, messageNumItem)
	}

	logs, sub, err := _ArbitrumInbox.contract.WatchLogs(opts, "InboxMessageDeliveredFromOrigin", messageNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbitrumInboxInboxMessageDeliveredFromOrigin)
				if err := _ArbitrumInbox.contract.UnpackLog(event, "InboxMessageDeliveredFromOrigin", log); err != nil {
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

func (_ArbitrumInbox *ArbitrumInboxFilterer) ParseInboxMessageDeliveredFromOrigin(log types.Log) (*ArbitrumInboxInboxMessageDeliveredFromOrigin, error) {
	event := new(ArbitrumInboxInboxMessageDeliveredFromOrigin)
	if err := _ArbitrumInbox.contract.UnpackLog(event, "InboxMessageDeliveredFromOrigin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ArbitrumInbox *ArbitrumInbox) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ArbitrumInbox.abi.Events["InboxMessageDelivered"].ID:
		return _ArbitrumInbox.ParseInboxMessageDelivered(log)
	case _ArbitrumInbox.abi.Events["InboxMessageDeliveredFromOrigin"].ID:
		return _ArbitrumInbox.ParseInboxMessageDeliveredFromOrigin(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ArbitrumInboxInboxMessageDelivered) Topic() common.Hash {
	return common.HexToHash("0xff64905f73a67fb594e0f940a8075a860db489ad991e032f48c81123eb52d60b")
}

func (ArbitrumInboxInboxMessageDeliveredFromOrigin) Topic() common.Hash {
	return common.HexToHash("0xab532385be8f1005a4b6ba8fa20a2245facb346134ac739fe9a5198dc1580b9c")
}

func (_ArbitrumInbox *ArbitrumInbox) Address() common.Address {
	return _ArbitrumInbox.address
}

type ArbitrumInboxInterface interface {
	AllowListEnabled(opts *bind.CallOpts) (bool, error)

	Bridge(opts *bind.CallOpts) (common.Address, error)

	CalculateRetryableSubmissionFee(opts *bind.CallOpts, dataLength *big.Int, baseFee *big.Int) (*big.Int, error)

	GetProxyAdmin(opts *bind.CallOpts) (common.Address, error)

	IsAllowed(opts *bind.CallOpts, user common.Address) (bool, error)

	MaxDataSize(opts *bind.CallOpts) (*big.Int, error)

	SequencerInbox(opts *bind.CallOpts) (common.Address, error)

	Initialize(opts *bind.TransactOpts, _bridge common.Address, _sequencerInbox common.Address) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	SendContractTransaction(opts *bind.TransactOpts, gasLimit *big.Int, maxFeePerGas *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error)

	SendL2Message(opts *bind.TransactOpts, messageData []byte) (*types.Transaction, error)

	SendL2MessageFromOrigin(opts *bind.TransactOpts, messageData []byte) (*types.Transaction, error)

	SendUnsignedTransaction(opts *bind.TransactOpts, gasLimit *big.Int, maxFeePerGas *big.Int, nonce *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error)

	SetAllowList(opts *bind.TransactOpts, user []common.Address, val []bool) (*types.Transaction, error)

	SetAllowListEnabled(opts *bind.TransactOpts, _allowListEnabled bool) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterInboxMessageDelivered(opts *bind.FilterOpts, messageNum []*big.Int) (*ArbitrumInboxInboxMessageDeliveredIterator, error)

	WatchInboxMessageDelivered(opts *bind.WatchOpts, sink chan<- *ArbitrumInboxInboxMessageDelivered, messageNum []*big.Int) (event.Subscription, error)

	ParseInboxMessageDelivered(log types.Log) (*ArbitrumInboxInboxMessageDelivered, error)

	FilterInboxMessageDeliveredFromOrigin(opts *bind.FilterOpts, messageNum []*big.Int) (*ArbitrumInboxInboxMessageDeliveredFromOriginIterator, error)

	WatchInboxMessageDeliveredFromOrigin(opts *bind.WatchOpts, sink chan<- *ArbitrumInboxInboxMessageDeliveredFromOrigin, messageNum []*big.Int) (event.Subscription, error)

	ParseInboxMessageDeliveredFromOrigin(log types.Log) (*ArbitrumInboxInboxMessageDeliveredFromOrigin, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
