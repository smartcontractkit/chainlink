// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple_log_upkeep_counter_wrapper

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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gethwrappers/generated"
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

type Log struct {
	Index       *big.Int
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var SimpleLogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeToPerform\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeToPerform\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060006001819055438155600281905560035561088a806100326000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806361bc221a1161005b57806361bc221a146100d4578063806b984f146100dd578063917d895f146100e6578063c6066f0d146100ef57600080fd5b80632cb158641461008257806340691db41461009e5780634585e33b146100bf575b600080fd5b61008b60025481565b6040519081526020015b60405180910390f35b6100b16100ac366004610384565b6100f8565b60405161009592919061055e565b6100d26100cd366004610312565b61012a565b005b61008b60035481565b61008b60005481565b61008b60015481565b61008b60045481565b6000606060018460405160200161010f91906105db565b604051602081830303815290604052915091505b9250929050565b60025461013657436002555b436000556003546101489060016107f0565b6003556000805460015561015e828401846103f1565b90508060200151426101709190610808565b6004819055600254600054600154600354604080519485526020850193909352918301526060820152608081019190915232907f4874b8dd61a40fe23599b4360a9a824d7081742fca9f555bcee3d389c4f4bd659060a00160405180910390a2505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146101f957600080fd5b919050565b600082601f83011261020f57600080fd5b8135602067ffffffffffffffff82111561022b5761022b61084e565b8160051b61023a8282016106d6565b83815282810190868401838801850189101561025557600080fd5b600093505b8584101561027857803583526001939093019291840191840161025a565b50979650505050505050565b600082601f83011261029557600080fd5b813567ffffffffffffffff8111156102af576102af61084e565b6102e060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016106d6565b8181528460208386010111156102f557600080fd5b816020850160208301376000918101602001919091529392505050565b6000806020838503121561032557600080fd5b823567ffffffffffffffff8082111561033d57600080fd5b818501915085601f83011261035157600080fd5b81358181111561036057600080fd5b86602082850101111561037257600080fd5b60209290920196919550909350505050565b6000806040838503121561039757600080fd5b823567ffffffffffffffff808211156103af57600080fd5b9084019061010082870312156103c457600080fd5b909250602084013590808211156103da57600080fd5b506103e785828601610284565b9150509250929050565b60006020828403121561040357600080fd5b813567ffffffffffffffff8082111561041b57600080fd5b90830190610100828603121561043057600080fd5b6104386106ac565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261047060a084016101d5565b60a082015260c08301358281111561048757600080fd5b610493878286016101fe565b60c08301525060e0830135828111156104ab57600080fd5b6104b787828601610284565b60e08301525095945050505050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8311156104f857600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b821515815260006020604081840152835180604085015260005b8181101561059457858101830151858201606001528201610578565b818111156105a6576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b6020815281356020820152602082013560408201526040820135606082015260608201356080820152608082013560a082015273ffffffffffffffffffffffffffffffffffffffff61062f60a084016101d5565b1660c0820152600061064460c0840184610725565b6101008060e086015261065c610120860183856104c6565b925061066b60e087018761078c565b92507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086850301828701526106a1848483610515565b979650505050505050565b604051610100810167ffffffffffffffff811182821017156106d0576106d061084e565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561071d5761071d61084e565b604052919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261075a57600080fd5b830160208101925035905067ffffffffffffffff81111561077a57600080fd5b8060051b360383131561012357600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126107c157600080fd5b830160208101925035905067ffffffffffffffff8111156107e157600080fd5b80360383131561012357600080fd5b600082198211156108035761080361081f565b500190565b60008282101561081a5761081a61081f565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var SimpleLogUpkeepCounterABI = SimpleLogUpkeepCounterMetaData.ABI

var SimpleLogUpkeepCounterBin = SimpleLogUpkeepCounterMetaData.Bin

func DeploySimpleLogUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleLogUpkeepCounter, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleLogUpkeepCounterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleLogUpkeepCounter{address: address, abi: *parsed, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

type SimpleLogUpkeepCounter struct {
	address common.Address
	abi     abi.ABI
	SimpleLogUpkeepCounterCaller
	SimpleLogUpkeepCounterTransactor
	SimpleLogUpkeepCounterFilterer
}

type SimpleLogUpkeepCounterCaller struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterTransactor struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterFilterer struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterSession struct {
	Contract     *SimpleLogUpkeepCounter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterCallerSession struct {
	Contract *SimpleLogUpkeepCounterCaller
	CallOpts bind.CallOpts
}

type SimpleLogUpkeepCounterTransactorSession struct {
	Contract     *SimpleLogUpkeepCounterTransactor
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterRaw struct {
	Contract *SimpleLogUpkeepCounter
}

type SimpleLogUpkeepCounterCallerRaw struct {
	Contract *SimpleLogUpkeepCounterCaller
}

type SimpleLogUpkeepCounterTransactorRaw struct {
	Contract *SimpleLogUpkeepCounterTransactor
}

func NewSimpleLogUpkeepCounter(address common.Address, backend bind.ContractBackend) (*SimpleLogUpkeepCounter, error) {
	abi, err := abi.JSON(strings.NewReader(SimpleLogUpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSimpleLogUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounter{address: address, abi: abi, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

func NewSimpleLogUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*SimpleLogUpkeepCounterCaller, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterCaller{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleLogUpkeepCounterTransactor, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterTransactor{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleLogUpkeepCounterFilterer, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterFilterer{contract: contract}, nil
}

func bindSimpleLogUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) CheckLog(opts *bind.CallOpts, log Log, arg1 []byte) (bool, []byte, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "checkLog", log, arg1)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckLog(log Log, arg1 []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, arg1)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) CheckLog(log Log, arg1 []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, arg1)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) TimeToPerform(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "timeToPerform")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

type SimpleLogUpkeepCounterPerformingUpkeepIterator struct {
	Event *SimpleLogUpkeepCounterPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
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
		it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
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

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleLogUpkeepCounterPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	TimeToPerform *big.Int
	Raw           types.Log
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterPerformingUpkeepIterator{contract: _SimpleLogUpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleLogUpkeepCounterPerformingUpkeep)
				if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error) {
	event := new(SimpleLogUpkeepCounterPerformingUpkeep)
	if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SimpleLogUpkeepCounter.abi.Events["PerformingUpkeep"].ID:
		return _SimpleLogUpkeepCounter.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SimpleLogUpkeepCounterPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x4874b8dd61a40fe23599b4360a9a824d7081742fca9f555bcee3d389c4f4bd65")
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) Address() common.Address {
	return _SimpleLogUpkeepCounter.address
}

type SimpleLogUpkeepCounterInterface interface {
	CheckLog(opts *bind.CallOpts, log Log, arg1 []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TimeToPerform(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
