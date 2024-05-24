// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_malicious_migrator

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

var VRFV2PlusMaliciousMigratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516102cb3803806102cb83398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b610238806100936000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80638ea9811714610030575b600080fd5b61004361003e36600461011b565b610045565b005b600080546040805160c081018252838152602080820185905281830185905260608201859052608082018590528251908101835293845260a0810193909352517f9b1c385e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911691639b1c385e916100d49190600401610158565b6020604051808303816000875af11580156100f3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101179190610212565b5050565b60006020828403121561012d57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461015157600080fd5b9392505050565b6000602080835283518184015280840151604084015261ffff6040850151166060840152606084015163ffffffff80821660808601528060808701511660a0860152505060a084015160c08085015280518060e086015260005b818110156101cf57828101840151868201610100015283016101b2565b5061010092506000838287010152827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116860101935050505092915050565b60006020828403121561022457600080fd5b505191905056fea164736f6c6343000813000a",
}

var VRFV2PlusMaliciousMigratorABI = VRFV2PlusMaliciousMigratorMetaData.ABI

var VRFV2PlusMaliciousMigratorBin = VRFV2PlusMaliciousMigratorMetaData.Bin

func DeployVRFV2PlusMaliciousMigrator(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFV2PlusMaliciousMigrator, error) {
	parsed, err := VRFV2PlusMaliciousMigratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusMaliciousMigratorBin), backend, _vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusMaliciousMigrator{address: address, abi: *parsed, VRFV2PlusMaliciousMigratorCaller: VRFV2PlusMaliciousMigratorCaller{contract: contract}, VRFV2PlusMaliciousMigratorTransactor: VRFV2PlusMaliciousMigratorTransactor{contract: contract}, VRFV2PlusMaliciousMigratorFilterer: VRFV2PlusMaliciousMigratorFilterer{contract: contract}}, nil
}

type VRFV2PlusMaliciousMigrator struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusMaliciousMigratorCaller
	VRFV2PlusMaliciousMigratorTransactor
	VRFV2PlusMaliciousMigratorFilterer
}

type VRFV2PlusMaliciousMigratorCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusMaliciousMigratorTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusMaliciousMigratorFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusMaliciousMigratorSession struct {
	Contract     *VRFV2PlusMaliciousMigrator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusMaliciousMigratorCallerSession struct {
	Contract *VRFV2PlusMaliciousMigratorCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusMaliciousMigratorTransactorSession struct {
	Contract     *VRFV2PlusMaliciousMigratorTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusMaliciousMigratorRaw struct {
	Contract *VRFV2PlusMaliciousMigrator
}

type VRFV2PlusMaliciousMigratorCallerRaw struct {
	Contract *VRFV2PlusMaliciousMigratorCaller
}

type VRFV2PlusMaliciousMigratorTransactorRaw struct {
	Contract *VRFV2PlusMaliciousMigratorTransactor
}

func NewVRFV2PlusMaliciousMigrator(address common.Address, backend bind.ContractBackend) (*VRFV2PlusMaliciousMigrator, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusMaliciousMigratorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusMaliciousMigrator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusMaliciousMigrator{address: address, abi: abi, VRFV2PlusMaliciousMigratorCaller: VRFV2PlusMaliciousMigratorCaller{contract: contract}, VRFV2PlusMaliciousMigratorTransactor: VRFV2PlusMaliciousMigratorTransactor{contract: contract}, VRFV2PlusMaliciousMigratorFilterer: VRFV2PlusMaliciousMigratorFilterer{contract: contract}}, nil
}

func NewVRFV2PlusMaliciousMigratorCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusMaliciousMigratorCaller, error) {
	contract, err := bindVRFV2PlusMaliciousMigrator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusMaliciousMigratorCaller{contract: contract}, nil
}

func NewVRFV2PlusMaliciousMigratorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusMaliciousMigratorTransactor, error) {
	contract, err := bindVRFV2PlusMaliciousMigrator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusMaliciousMigratorTransactor{contract: contract}, nil
}

func NewVRFV2PlusMaliciousMigratorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusMaliciousMigratorFilterer, error) {
	contract, err := bindVRFV2PlusMaliciousMigrator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusMaliciousMigratorFilterer{contract: contract}, nil
}

func bindVRFV2PlusMaliciousMigrator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusMaliciousMigratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusMaliciousMigrator.Contract.VRFV2PlusMaliciousMigratorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.Contract.VRFV2PlusMaliciousMigratorTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.Contract.VRFV2PlusMaliciousMigratorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusMaliciousMigrator.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorTransactor) SetCoordinator(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.contract.Transact(opts, "setCoordinator", arg0)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorSession) SetCoordinator(arg0 common.Address) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.Contract.SetCoordinator(&_VRFV2PlusMaliciousMigrator.TransactOpts, arg0)
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorTransactorSession) SetCoordinator(arg0 common.Address) (*types.Transaction, error) {
	return _VRFV2PlusMaliciousMigrator.Contract.SetCoordinator(&_VRFV2PlusMaliciousMigrator.TransactOpts, arg0)
}

type VRFV2PlusMaliciousMigratorCoordinatorSetIterator struct {
	Event *VRFV2PlusMaliciousMigratorCoordinatorSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusMaliciousMigratorCoordinatorSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusMaliciousMigratorCoordinatorSet)
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
		it.Event = new(VRFV2PlusMaliciousMigratorCoordinatorSet)
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

func (it *VRFV2PlusMaliciousMigratorCoordinatorSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusMaliciousMigratorCoordinatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusMaliciousMigratorCoordinatorSet struct {
	VrfCoordinator common.Address
	Raw            types.Log
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorFilterer) FilterCoordinatorSet(opts *bind.FilterOpts) (*VRFV2PlusMaliciousMigratorCoordinatorSetIterator, error) {

	logs, sub, err := _VRFV2PlusMaliciousMigrator.contract.FilterLogs(opts, "CoordinatorSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusMaliciousMigratorCoordinatorSetIterator{contract: _VRFV2PlusMaliciousMigrator.contract, event: "CoordinatorSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorFilterer) WatchCoordinatorSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusMaliciousMigratorCoordinatorSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusMaliciousMigrator.contract.WatchLogs(opts, "CoordinatorSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusMaliciousMigratorCoordinatorSet)
				if err := _VRFV2PlusMaliciousMigrator.contract.UnpackLog(event, "CoordinatorSet", log); err != nil {
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

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigratorFilterer) ParseCoordinatorSet(log types.Log) (*VRFV2PlusMaliciousMigratorCoordinatorSet, error) {
	event := new(VRFV2PlusMaliciousMigratorCoordinatorSet)
	if err := _VRFV2PlusMaliciousMigrator.contract.UnpackLog(event, "CoordinatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigrator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusMaliciousMigrator.abi.Events["CoordinatorSet"].ID:
		return _VRFV2PlusMaliciousMigrator.ParseCoordinatorSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusMaliciousMigratorCoordinatorSet) Topic() common.Hash {
	return common.HexToHash("0xd1a6a14209a385a964d036e404cb5cfb71f4000cdb03c9366292430787261be6")
}

func (_VRFV2PlusMaliciousMigrator *VRFV2PlusMaliciousMigrator) Address() common.Address {
	return _VRFV2PlusMaliciousMigrator.address
}

type VRFV2PlusMaliciousMigratorInterface interface {
	SetCoordinator(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error)

	FilterCoordinatorSet(opts *bind.FilterOpts) (*VRFV2PlusMaliciousMigratorCoordinatorSetIterator, error)

	WatchCoordinatorSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusMaliciousMigratorCoordinatorSet) (event.Subscription, error)

	ParseCoordinatorSet(log types.Log) (*VRFV2PlusMaliciousMigratorCoordinatorSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
