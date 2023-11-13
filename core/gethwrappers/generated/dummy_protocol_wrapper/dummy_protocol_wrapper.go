// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dummy_protocol_wrapper

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

var DummyProtocolMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"LimitOrderExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"LimitOrderSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"LimitOrderWithdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"executeLimitOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetContract\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"t0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"t1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"t2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"t3\",\"type\":\"bytes32\"}],\"name\":\"getAdvancedLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"t0\",\"type\":\"bytes32\"}],\"name\":\"getBasicLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"sendLimitedOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"withdrawLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506103e2806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80636bb17e4c116100505780636bb17e4c1461013d5780639ab74b0e1461013d578063af38c9c21461015057600080fd5b80632065ff7e1461006c5780635f35f80b14610081575b600080fd5b61007f61007a3660046102ad565b6101f0565b005b61012761008f3660046102e2565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b604051610134919061033f565b60405180910390f35b61007f61014b3660046102ad565b61023a565b61012761015e3660046103ab565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff94909416808252600060208084018281528486019687526060808601848152608080880186815260a0988901968752895195860197909752925160ff168489015297519083015295519581019590955290519184019190915251828401528051808303909301835260e0909101905290565b8073ffffffffffffffffffffffffffffffffffffffff1682847fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd60405160405180910390a4505050565b8073ffffffffffffffffffffffffffffffffffffffff1682847f3e9c37b3143f2eb7e9a2a0f8091b6de097b62efcfe48e1f68847a832e521750a60405160405180910390a4505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146102a857600080fd5b919050565b6000806000606084860312156102c257600080fd5b83359250602084013591506102d960408501610284565b90509250925092565b60008060008060008060c087890312156102fb57600080fd5b61030487610284565b9550602087013560ff8116811461031a57600080fd5b95989597505050506040840135936060810135936080820135935060a0909101359150565b600060208083528351808285015260005b8181101561036c57858101830151858201604001528201610350565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b600080604083850312156103be57600080fd5b6103c783610284565b94602093909301359350505056fea164736f6c6343000810000a",
}

var DummyProtocolABI = DummyProtocolMetaData.ABI

var DummyProtocolBin = DummyProtocolMetaData.Bin

func DeployDummyProtocol(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DummyProtocol, error) {
	parsed, err := DummyProtocolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DummyProtocolBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DummyProtocol{address: address, abi: *parsed, DummyProtocolCaller: DummyProtocolCaller{contract: contract}, DummyProtocolTransactor: DummyProtocolTransactor{contract: contract}, DummyProtocolFilterer: DummyProtocolFilterer{contract: contract}}, nil
}

type DummyProtocol struct {
	address common.Address
	abi     abi.ABI
	DummyProtocolCaller
	DummyProtocolTransactor
	DummyProtocolFilterer
}

type DummyProtocolCaller struct {
	contract *bind.BoundContract
}

type DummyProtocolTransactor struct {
	contract *bind.BoundContract
}

type DummyProtocolFilterer struct {
	contract *bind.BoundContract
}

type DummyProtocolSession struct {
	Contract     *DummyProtocol
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DummyProtocolCallerSession struct {
	Contract *DummyProtocolCaller
	CallOpts bind.CallOpts
}

type DummyProtocolTransactorSession struct {
	Contract     *DummyProtocolTransactor
	TransactOpts bind.TransactOpts
}

type DummyProtocolRaw struct {
	Contract *DummyProtocol
}

type DummyProtocolCallerRaw struct {
	Contract *DummyProtocolCaller
}

type DummyProtocolTransactorRaw struct {
	Contract *DummyProtocolTransactor
}

func NewDummyProtocol(address common.Address, backend bind.ContractBackend) (*DummyProtocol, error) {
	abi, err := abi.JSON(strings.NewReader(DummyProtocolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDummyProtocol(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DummyProtocol{address: address, abi: abi, DummyProtocolCaller: DummyProtocolCaller{contract: contract}, DummyProtocolTransactor: DummyProtocolTransactor{contract: contract}, DummyProtocolFilterer: DummyProtocolFilterer{contract: contract}}, nil
}

func NewDummyProtocolCaller(address common.Address, caller bind.ContractCaller) (*DummyProtocolCaller, error) {
	contract, err := bindDummyProtocol(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DummyProtocolCaller{contract: contract}, nil
}

func NewDummyProtocolTransactor(address common.Address, transactor bind.ContractTransactor) (*DummyProtocolTransactor, error) {
	contract, err := bindDummyProtocol(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DummyProtocolTransactor{contract: contract}, nil
}

func NewDummyProtocolFilterer(address common.Address, filterer bind.ContractFilterer) (*DummyProtocolFilterer, error) {
	contract, err := bindDummyProtocol(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DummyProtocolFilterer{contract: contract}, nil
}

func bindDummyProtocol(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DummyProtocolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DummyProtocol *DummyProtocolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyProtocol.Contract.DummyProtocolCaller.contract.Call(opts, result, method, params...)
}

func (_DummyProtocol *DummyProtocolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyProtocol.Contract.DummyProtocolTransactor.contract.Transfer(opts)
}

func (_DummyProtocol *DummyProtocolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyProtocol.Contract.DummyProtocolTransactor.contract.Transact(opts, method, params...)
}

func (_DummyProtocol *DummyProtocolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyProtocol.Contract.contract.Call(opts, result, method, params...)
}

func (_DummyProtocol *DummyProtocolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyProtocol.Contract.contract.Transfer(opts)
}

func (_DummyProtocol *DummyProtocolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyProtocol.Contract.contract.Transact(opts, method, params...)
}

func (_DummyProtocol *DummyProtocolCaller) GetAdvancedLogTriggerConfig(opts *bind.CallOpts, targetContract common.Address, selector uint8, t0 [32]byte, t1 [32]byte, t2 [32]byte, t3 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _DummyProtocol.contract.Call(opts, &out, "getAdvancedLogTriggerConfig", targetContract, selector, t0, t1, t2, t3)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_DummyProtocol *DummyProtocolSession) GetAdvancedLogTriggerConfig(targetContract common.Address, selector uint8, t0 [32]byte, t1 [32]byte, t2 [32]byte, t3 [32]byte) ([]byte, error) {
	return _DummyProtocol.Contract.GetAdvancedLogTriggerConfig(&_DummyProtocol.CallOpts, targetContract, selector, t0, t1, t2, t3)
}

func (_DummyProtocol *DummyProtocolCallerSession) GetAdvancedLogTriggerConfig(targetContract common.Address, selector uint8, t0 [32]byte, t1 [32]byte, t2 [32]byte, t3 [32]byte) ([]byte, error) {
	return _DummyProtocol.Contract.GetAdvancedLogTriggerConfig(&_DummyProtocol.CallOpts, targetContract, selector, t0, t1, t2, t3)
}

func (_DummyProtocol *DummyProtocolCaller) GetBasicLogTriggerConfig(opts *bind.CallOpts, targetContract common.Address, t0 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _DummyProtocol.contract.Call(opts, &out, "getBasicLogTriggerConfig", targetContract, t0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_DummyProtocol *DummyProtocolSession) GetBasicLogTriggerConfig(targetContract common.Address, t0 [32]byte) ([]byte, error) {
	return _DummyProtocol.Contract.GetBasicLogTriggerConfig(&_DummyProtocol.CallOpts, targetContract, t0)
}

func (_DummyProtocol *DummyProtocolCallerSession) GetBasicLogTriggerConfig(targetContract common.Address, t0 [32]byte) ([]byte, error) {
	return _DummyProtocol.Contract.GetBasicLogTriggerConfig(&_DummyProtocol.CallOpts, targetContract, t0)
}

func (_DummyProtocol *DummyProtocolTransactor) ExecuteLimitOrder(opts *bind.TransactOpts, orderId *big.Int, amount *big.Int, exchange common.Address) (*types.Transaction, error) {
	return _DummyProtocol.contract.Transact(opts, "executeLimitOrder", orderId, amount, exchange)
}

func (_DummyProtocol *DummyProtocolSession) ExecuteLimitOrder(orderId *big.Int, amount *big.Int, exchange common.Address) (*types.Transaction, error) {
	return _DummyProtocol.Contract.ExecuteLimitOrder(&_DummyProtocol.TransactOpts, orderId, amount, exchange)
}

func (_DummyProtocol *DummyProtocolTransactorSession) ExecuteLimitOrder(orderId *big.Int, amount *big.Int, exchange common.Address) (*types.Transaction, error) {
	return _DummyProtocol.Contract.ExecuteLimitOrder(&_DummyProtocol.TransactOpts, orderId, amount, exchange)
}

func (_DummyProtocol *DummyProtocolTransactor) SendLimitedOrder(opts *bind.TransactOpts, amount *big.Int, price *big.Int, to common.Address) (*types.Transaction, error) {
	return _DummyProtocol.contract.Transact(opts, "sendLimitedOrder", amount, price, to)
}

func (_DummyProtocol *DummyProtocolSession) SendLimitedOrder(amount *big.Int, price *big.Int, to common.Address) (*types.Transaction, error) {
	return _DummyProtocol.Contract.SendLimitedOrder(&_DummyProtocol.TransactOpts, amount, price, to)
}

func (_DummyProtocol *DummyProtocolTransactorSession) SendLimitedOrder(amount *big.Int, price *big.Int, to common.Address) (*types.Transaction, error) {
	return _DummyProtocol.Contract.SendLimitedOrder(&_DummyProtocol.TransactOpts, amount, price, to)
}

func (_DummyProtocol *DummyProtocolTransactor) WithdrawLimit(opts *bind.TransactOpts, amount *big.Int, price *big.Int, from common.Address) (*types.Transaction, error) {
	return _DummyProtocol.contract.Transact(opts, "withdrawLimit", amount, price, from)
}

func (_DummyProtocol *DummyProtocolSession) WithdrawLimit(amount *big.Int, price *big.Int, from common.Address) (*types.Transaction, error) {
	return _DummyProtocol.Contract.WithdrawLimit(&_DummyProtocol.TransactOpts, amount, price, from)
}

func (_DummyProtocol *DummyProtocolTransactorSession) WithdrawLimit(amount *big.Int, price *big.Int, from common.Address) (*types.Transaction, error) {
	return _DummyProtocol.Contract.WithdrawLimit(&_DummyProtocol.TransactOpts, amount, price, from)
}

type DummyProtocolLimitOrderExecutedIterator struct {
	Event *DummyProtocolLimitOrderExecuted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyProtocolLimitOrderExecutedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyProtocolLimitOrderExecuted)
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
		it.Event = new(DummyProtocolLimitOrderExecuted)
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

func (it *DummyProtocolLimitOrderExecutedIterator) Error() error {
	return it.fail
}

func (it *DummyProtocolLimitOrderExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyProtocolLimitOrderExecuted struct {
	OrderId  *big.Int
	Amount   *big.Int
	Exchange common.Address
	Raw      types.Log
}

func (_DummyProtocol *DummyProtocolFilterer) FilterLimitOrderExecuted(opts *bind.FilterOpts, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (*DummyProtocolLimitOrderExecutedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var exchangeRule []interface{}
	for _, exchangeItem := range exchange {
		exchangeRule = append(exchangeRule, exchangeItem)
	}

	logs, sub, err := _DummyProtocol.contract.FilterLogs(opts, "LimitOrderExecuted", orderIdRule, amountRule, exchangeRule)
	if err != nil {
		return nil, err
	}
	return &DummyProtocolLimitOrderExecutedIterator{contract: _DummyProtocol.contract, event: "LimitOrderExecuted", logs: logs, sub: sub}, nil
}

func (_DummyProtocol *DummyProtocolFilterer) WatchLimitOrderExecuted(opts *bind.WatchOpts, sink chan<- *DummyProtocolLimitOrderExecuted, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var exchangeRule []interface{}
	for _, exchangeItem := range exchange {
		exchangeRule = append(exchangeRule, exchangeItem)
	}

	logs, sub, err := _DummyProtocol.contract.WatchLogs(opts, "LimitOrderExecuted", orderIdRule, amountRule, exchangeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyProtocolLimitOrderExecuted)
				if err := _DummyProtocol.contract.UnpackLog(event, "LimitOrderExecuted", log); err != nil {
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

func (_DummyProtocol *DummyProtocolFilterer) ParseLimitOrderExecuted(log types.Log) (*DummyProtocolLimitOrderExecuted, error) {
	event := new(DummyProtocolLimitOrderExecuted)
	if err := _DummyProtocol.contract.UnpackLog(event, "LimitOrderExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyProtocolLimitOrderSentIterator struct {
	Event *DummyProtocolLimitOrderSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyProtocolLimitOrderSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyProtocolLimitOrderSent)
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
		it.Event = new(DummyProtocolLimitOrderSent)
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

func (it *DummyProtocolLimitOrderSentIterator) Error() error {
	return it.fail
}

func (it *DummyProtocolLimitOrderSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyProtocolLimitOrderSent struct {
	Amount *big.Int
	Price  *big.Int
	To     common.Address
	Raw    types.Log
}

func (_DummyProtocol *DummyProtocolFilterer) FilterLimitOrderSent(opts *bind.FilterOpts, amount []*big.Int, price []*big.Int, to []common.Address) (*DummyProtocolLimitOrderSentIterator, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DummyProtocol.contract.FilterLogs(opts, "LimitOrderSent", amountRule, priceRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DummyProtocolLimitOrderSentIterator{contract: _DummyProtocol.contract, event: "LimitOrderSent", logs: logs, sub: sub}, nil
}

func (_DummyProtocol *DummyProtocolFilterer) WatchLimitOrderSent(opts *bind.WatchOpts, sink chan<- *DummyProtocolLimitOrderSent, amount []*big.Int, price []*big.Int, to []common.Address) (event.Subscription, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DummyProtocol.contract.WatchLogs(opts, "LimitOrderSent", amountRule, priceRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyProtocolLimitOrderSent)
				if err := _DummyProtocol.contract.UnpackLog(event, "LimitOrderSent", log); err != nil {
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

func (_DummyProtocol *DummyProtocolFilterer) ParseLimitOrderSent(log types.Log) (*DummyProtocolLimitOrderSent, error) {
	event := new(DummyProtocolLimitOrderSent)
	if err := _DummyProtocol.contract.UnpackLog(event, "LimitOrderSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyProtocolLimitOrderWithdrawnIterator struct {
	Event *DummyProtocolLimitOrderWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyProtocolLimitOrderWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyProtocolLimitOrderWithdrawn)
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
		it.Event = new(DummyProtocolLimitOrderWithdrawn)
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

func (it *DummyProtocolLimitOrderWithdrawnIterator) Error() error {
	return it.fail
}

func (it *DummyProtocolLimitOrderWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyProtocolLimitOrderWithdrawn struct {
	Amount *big.Int
	Price  *big.Int
	From   common.Address
	Raw    types.Log
}

func (_DummyProtocol *DummyProtocolFilterer) FilterLimitOrderWithdrawn(opts *bind.FilterOpts, amount []*big.Int, price []*big.Int, from []common.Address) (*DummyProtocolLimitOrderWithdrawnIterator, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _DummyProtocol.contract.FilterLogs(opts, "LimitOrderWithdrawn", amountRule, priceRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &DummyProtocolLimitOrderWithdrawnIterator{contract: _DummyProtocol.contract, event: "LimitOrderWithdrawn", logs: logs, sub: sub}, nil
}

func (_DummyProtocol *DummyProtocolFilterer) WatchLimitOrderWithdrawn(opts *bind.WatchOpts, sink chan<- *DummyProtocolLimitOrderWithdrawn, amount []*big.Int, price []*big.Int, from []common.Address) (event.Subscription, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _DummyProtocol.contract.WatchLogs(opts, "LimitOrderWithdrawn", amountRule, priceRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyProtocolLimitOrderWithdrawn)
				if err := _DummyProtocol.contract.UnpackLog(event, "LimitOrderWithdrawn", log); err != nil {
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

func (_DummyProtocol *DummyProtocolFilterer) ParseLimitOrderWithdrawn(log types.Log) (*DummyProtocolLimitOrderWithdrawn, error) {
	event := new(DummyProtocolLimitOrderWithdrawn)
	if err := _DummyProtocol.contract.UnpackLog(event, "LimitOrderWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_DummyProtocol *DummyProtocol) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DummyProtocol.abi.Events["LimitOrderExecuted"].ID:
		return _DummyProtocol.ParseLimitOrderExecuted(log)
	case _DummyProtocol.abi.Events["LimitOrderSent"].ID:
		return _DummyProtocol.ParseLimitOrderSent(log)
	case _DummyProtocol.abi.Events["LimitOrderWithdrawn"].ID:
		return _DummyProtocol.ParseLimitOrderWithdrawn(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DummyProtocolLimitOrderExecuted) Topic() common.Hash {
	return common.HexToHash("0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd")
}

func (DummyProtocolLimitOrderSent) Topic() common.Hash {
	return common.HexToHash("0x3e9c37b3143f2eb7e9a2a0f8091b6de097b62efcfe48e1f68847a832e521750a")
}

func (DummyProtocolLimitOrderWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x0a71b8ed921ff64d49e4d39449f8a21094f38a0aeae489c3051aedd63f2c229f")
}

func (_DummyProtocol *DummyProtocol) Address() common.Address {
	return _DummyProtocol.address
}

type DummyProtocolInterface interface {
	GetAdvancedLogTriggerConfig(opts *bind.CallOpts, targetContract common.Address, selector uint8, t0 [32]byte, t1 [32]byte, t2 [32]byte, t3 [32]byte) ([]byte, error)

	GetBasicLogTriggerConfig(opts *bind.CallOpts, targetContract common.Address, t0 [32]byte) ([]byte, error)

	ExecuteLimitOrder(opts *bind.TransactOpts, orderId *big.Int, amount *big.Int, exchange common.Address) (*types.Transaction, error)

	SendLimitedOrder(opts *bind.TransactOpts, amount *big.Int, price *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawLimit(opts *bind.TransactOpts, amount *big.Int, price *big.Int, from common.Address) (*types.Transaction, error)

	FilterLimitOrderExecuted(opts *bind.FilterOpts, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (*DummyProtocolLimitOrderExecutedIterator, error)

	WatchLimitOrderExecuted(opts *bind.WatchOpts, sink chan<- *DummyProtocolLimitOrderExecuted, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (event.Subscription, error)

	ParseLimitOrderExecuted(log types.Log) (*DummyProtocolLimitOrderExecuted, error)

	FilterLimitOrderSent(opts *bind.FilterOpts, amount []*big.Int, price []*big.Int, to []common.Address) (*DummyProtocolLimitOrderSentIterator, error)

	WatchLimitOrderSent(opts *bind.WatchOpts, sink chan<- *DummyProtocolLimitOrderSent, amount []*big.Int, price []*big.Int, to []common.Address) (event.Subscription, error)

	ParseLimitOrderSent(log types.Log) (*DummyProtocolLimitOrderSent, error)

	FilterLimitOrderWithdrawn(opts *bind.FilterOpts, amount []*big.Int, price []*big.Int, from []common.Address) (*DummyProtocolLimitOrderWithdrawnIterator, error)

	WatchLimitOrderWithdrawn(opts *bind.WatchOpts, sink chan<- *DummyProtocolLimitOrderWithdrawn, amount []*big.Int, price []*big.Int, from []common.Address) (event.Subscription, error)

	ParseLimitOrderWithdrawn(log types.Log) (*DummyProtocolLimitOrderWithdrawn, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
