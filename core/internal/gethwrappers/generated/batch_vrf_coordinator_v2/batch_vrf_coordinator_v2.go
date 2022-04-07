// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package batch_vrf_coordinator_v2

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

type VRFTypesProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

type VRFTypesRequestCommitment struct {
	BlockNum         uint64
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
}

var BatchVRFCoordinatorV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"ErrorReturned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"RawErrorReturned\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRFTypes.Proof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structVRFTypes.RequestCommitment[]\",\"name\":\"rcs\",\"type\":\"tuple[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610bab380380610bab83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b608051610b1a610091600039600081816055015261011d0152610b1a6000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806308b2da0a1461003b5780633b2bcbf114610050575b600080fd5b61004e610049366004610655565b6100a0565b005b6100777f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b805182511461010f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b8251811015610320577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663af198b97848381518110610169576101696107b4565b6020026020010151848481518110610183576101836107b4565b60200260200101516040518363ffffffff1660e01b81526004016101a892919061080c565b6020604051808303816000875af1925050508015610201575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682019092526101fe9181019061090c565b60015b61030c5761020d61093a565b806308c379a0036102915750610221610955565b8061022c5750610293565b6000610250858481518110610243576102436107b4565b6020026020010151610325565b9050807f4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e836040516102829190610a6a565b60405180910390a2505061030e565b505b3d8080156102bd576040519150601f19603f3d011682016040523d82523d6000602084013e6102c2565b606091505b5060006102da858481518110610243576102436107b4565b9050807fbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b5836040516102829190610a6a565b505b8061031881610a7d565b915050610112565b505050565b6000806103358360000151610394565b9050808360800151604051602001610357929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b6000816040516020016103a79190610adc565b604051602081830303815290604052805190602001209050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810181811067ffffffffffffffff82111715610413576104136103c4565b60405250565b610120810167ffffffffffffffff81118282101715610413576104136103c4565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff8211171561047e5761047e6103c4565b6040525050565b60405161049181610419565b90565b600067ffffffffffffffff8211156104ae576104ae6103c4565b5060051b60200190565b600082601f8301126104c957600080fd5b6040516040810181811067ffffffffffffffff821117156104ec576104ec6103c4565b806040525080604084018581111561050357600080fd5b845b8181101561051d578035835260209283019201610505565b509195945050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461054c57600080fd5b919050565b803567ffffffffffffffff8116811461054c57600080fd5b803563ffffffff8116811461054c57600080fd5b600082601f83011261058e57600080fd5b8135602061059b82610494565b604080516105a9838261043a565b84815260a094850287018401948482019350888611156105c857600080fd5b8488015b868110156106475781818b0312156105e45760008081fd5b83516105ef816103f3565b6105f882610551565b8152610605878301610551565b87820152610614858301610569565b858201526060610625818401610569565b908201526080610636838201610528565b9082015285529385019381016105cc565b509098975050505050505050565b6000806040838503121561066857600080fd5b823567ffffffffffffffff8082111561068057600080fd5b818501915085601f83011261069457600080fd5b813560206106a182610494565b6040516106ae828261043a565b8381526101a0938402860183019383820192508a8511156106ce57600080fd5b958301955b848710156107865780878c0312156106eb5760008081fd5b6106f3610485565b6106fd8c896104b8565b815261070c8c60408a016104b8565b85820152608080890135604083015260a0808a0135606084015260c0808b01358385015260e0925061073f838c01610528565b8285015261010091506107548f838d016104b8565b908401526107668e6101408c016104b8565b9183019190915261018089013590820152835295860195918301916106d3565b509650508601359250508082111561079d57600080fd5b506107aa8582860161057d565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8060005b60028110156108065781518452602093840193909101906001016107e7565b50505050565b6000610240820190506108208285516107e3565b602084015161083260408401826107e3565b5060408401516080830152606084015160a0830152608084015160c083015273ffffffffffffffffffffffffffffffffffffffff60a08501511660e083015260c0840151610100610885818501836107e3565b60e0860151915061089a6101408501836107e3565b85015161018084015250825167ffffffffffffffff9081166101a08401526020840151166101c0830152604083015163ffffffff9081166101e0840152606084015116610200830152608083015173ffffffffffffffffffffffffffffffffffffffff166102208301525b9392505050565b60006020828403121561091e57600080fd5b81516bffffffffffffffffffffffff8116811461090557600080fd5b600060033d11156104915760046000803e5060005160e01c90565b600060443d10156109635790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff81602484011181841117156109b157505050505090565b82850191508151818111156109c95750505050505090565b843d87010160208285010111156109e35750505050505090565b6109f26020828601018761043a565b509095945050505050565b60008151808452602060005b82811015610a24578481018201518682018301528101610a09565b82811115610a355760008284880101525b50807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401168601019250505092915050565b60208152600061090560208301846109fd565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610ad5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b60408101818360005b6002811015610b04578151835260209283019290910190600101610ae5565b5050509291505056fea164736f6c634300080d000a",
}

var BatchVRFCoordinatorV2ABI = BatchVRFCoordinatorV2MetaData.ABI

var BatchVRFCoordinatorV2Bin = BatchVRFCoordinatorV2MetaData.Bin

func DeployBatchVRFCoordinatorV2(auth *bind.TransactOpts, backend bind.ContractBackend, coordinatorAddr common.Address) (common.Address, *types.Transaction, *BatchVRFCoordinatorV2, error) {
	parsed, err := BatchVRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchVRFCoordinatorV2Bin), backend, coordinatorAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchVRFCoordinatorV2{BatchVRFCoordinatorV2Caller: BatchVRFCoordinatorV2Caller{contract: contract}, BatchVRFCoordinatorV2Transactor: BatchVRFCoordinatorV2Transactor{contract: contract}, BatchVRFCoordinatorV2Filterer: BatchVRFCoordinatorV2Filterer{contract: contract}}, nil
}

type BatchVRFCoordinatorV2 struct {
	address common.Address
	abi     abi.ABI
	BatchVRFCoordinatorV2Caller
	BatchVRFCoordinatorV2Transactor
	BatchVRFCoordinatorV2Filterer
}

type BatchVRFCoordinatorV2Caller struct {
	contract *bind.BoundContract
}

type BatchVRFCoordinatorV2Transactor struct {
	contract *bind.BoundContract
}

type BatchVRFCoordinatorV2Filterer struct {
	contract *bind.BoundContract
}

type BatchVRFCoordinatorV2Session struct {
	Contract     *BatchVRFCoordinatorV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BatchVRFCoordinatorV2CallerSession struct {
	Contract *BatchVRFCoordinatorV2Caller
	CallOpts bind.CallOpts
}

type BatchVRFCoordinatorV2TransactorSession struct {
	Contract     *BatchVRFCoordinatorV2Transactor
	TransactOpts bind.TransactOpts
}

type BatchVRFCoordinatorV2Raw struct {
	Contract *BatchVRFCoordinatorV2
}

type BatchVRFCoordinatorV2CallerRaw struct {
	Contract *BatchVRFCoordinatorV2Caller
}

type BatchVRFCoordinatorV2TransactorRaw struct {
	Contract *BatchVRFCoordinatorV2Transactor
}

func NewBatchVRFCoordinatorV2(address common.Address, backend bind.ContractBackend) (*BatchVRFCoordinatorV2, error) {
	abi, err := abi.JSON(strings.NewReader(BatchVRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBatchVRFCoordinatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2{address: address, abi: abi, BatchVRFCoordinatorV2Caller: BatchVRFCoordinatorV2Caller{contract: contract}, BatchVRFCoordinatorV2Transactor: BatchVRFCoordinatorV2Transactor{contract: contract}, BatchVRFCoordinatorV2Filterer: BatchVRFCoordinatorV2Filterer{contract: contract}}, nil
}

func NewBatchVRFCoordinatorV2Caller(address common.Address, caller bind.ContractCaller) (*BatchVRFCoordinatorV2Caller, error) {
	contract, err := bindBatchVRFCoordinatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2Caller{contract: contract}, nil
}

func NewBatchVRFCoordinatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*BatchVRFCoordinatorV2Transactor, error) {
	contract, err := bindBatchVRFCoordinatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2Transactor{contract: contract}, nil
}

func NewBatchVRFCoordinatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*BatchVRFCoordinatorV2Filterer, error) {
	contract, err := bindBatchVRFCoordinatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2Filterer{contract: contract}, nil
}

func bindBatchVRFCoordinatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BatchVRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchVRFCoordinatorV2.Contract.BatchVRFCoordinatorV2Caller.contract.Call(opts, result, method, params...)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.Contract.BatchVRFCoordinatorV2Transactor.contract.Transfer(opts)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.Contract.BatchVRFCoordinatorV2Transactor.contract.Transact(opts, method, params...)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchVRFCoordinatorV2.Contract.contract.Call(opts, result, method, params...)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.Contract.contract.Transfer(opts)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.Contract.contract.Transact(opts, method, params...)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Caller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BatchVRFCoordinatorV2.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Session) COORDINATOR() (common.Address, error) {
	return _BatchVRFCoordinatorV2.Contract.COORDINATOR(&_BatchVRFCoordinatorV2.CallOpts)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2CallerSession) COORDINATOR() (common.Address, error) {
	return _BatchVRFCoordinatorV2.Contract.COORDINATOR(&_BatchVRFCoordinatorV2.CallOpts)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, proofs []VRFTypesProof, rcs []VRFTypesRequestCommitment) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.contract.Transact(opts, "fulfillRandomWords", proofs, rcs)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Session) FulfillRandomWords(proofs []VRFTypesProof, rcs []VRFTypesRequestCommitment) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.Contract.FulfillRandomWords(&_BatchVRFCoordinatorV2.TransactOpts, proofs, rcs)
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2TransactorSession) FulfillRandomWords(proofs []VRFTypesProof, rcs []VRFTypesRequestCommitment) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2.Contract.FulfillRandomWords(&_BatchVRFCoordinatorV2.TransactOpts, proofs, rcs)
}

type BatchVRFCoordinatorV2ErrorReturnedIterator struct {
	Event *BatchVRFCoordinatorV2ErrorReturned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BatchVRFCoordinatorV2ErrorReturnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchVRFCoordinatorV2ErrorReturned)
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
		it.Event = new(BatchVRFCoordinatorV2ErrorReturned)
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

func (it *BatchVRFCoordinatorV2ErrorReturnedIterator) Error() error {
	return it.fail
}

func (it *BatchVRFCoordinatorV2ErrorReturnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BatchVRFCoordinatorV2ErrorReturned struct {
	RequestId *big.Int
	Reason    string
	Raw       types.Log
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Filterer) FilterErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2ErrorReturnedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2.contract.FilterLogs(opts, "ErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2ErrorReturnedIterator{contract: _BatchVRFCoordinatorV2.contract, event: "ErrorReturned", logs: logs, sub: sub}, nil
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Filterer) WatchErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2ErrorReturned, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2.contract.WatchLogs(opts, "ErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BatchVRFCoordinatorV2ErrorReturned)
				if err := _BatchVRFCoordinatorV2.contract.UnpackLog(event, "ErrorReturned", log); err != nil {
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

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Filterer) ParseErrorReturned(log types.Log) (*BatchVRFCoordinatorV2ErrorReturned, error) {
	event := new(BatchVRFCoordinatorV2ErrorReturned)
	if err := _BatchVRFCoordinatorV2.contract.UnpackLog(event, "ErrorReturned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BatchVRFCoordinatorV2RawErrorReturnedIterator struct {
	Event *BatchVRFCoordinatorV2RawErrorReturned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BatchVRFCoordinatorV2RawErrorReturnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchVRFCoordinatorV2RawErrorReturned)
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
		it.Event = new(BatchVRFCoordinatorV2RawErrorReturned)
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

func (it *BatchVRFCoordinatorV2RawErrorReturnedIterator) Error() error {
	return it.fail
}

func (it *BatchVRFCoordinatorV2RawErrorReturnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BatchVRFCoordinatorV2RawErrorReturned struct {
	RequestId    *big.Int
	LowLevelData []byte
	Raw          types.Log
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Filterer) FilterRawErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2RawErrorReturnedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2.contract.FilterLogs(opts, "RawErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2RawErrorReturnedIterator{contract: _BatchVRFCoordinatorV2.contract, event: "RawErrorReturned", logs: logs, sub: sub}, nil
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Filterer) WatchRawErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2RawErrorReturned, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2.contract.WatchLogs(opts, "RawErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BatchVRFCoordinatorV2RawErrorReturned)
				if err := _BatchVRFCoordinatorV2.contract.UnpackLog(event, "RawErrorReturned", log); err != nil {
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

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2Filterer) ParseRawErrorReturned(log types.Log) (*BatchVRFCoordinatorV2RawErrorReturned, error) {
	event := new(BatchVRFCoordinatorV2RawErrorReturned)
	if err := _BatchVRFCoordinatorV2.contract.UnpackLog(event, "RawErrorReturned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BatchVRFCoordinatorV2.abi.Events["ErrorReturned"].ID:
		return _BatchVRFCoordinatorV2.ParseErrorReturned(log)
	case _BatchVRFCoordinatorV2.abi.Events["RawErrorReturned"].ID:
		return _BatchVRFCoordinatorV2.ParseRawErrorReturned(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BatchVRFCoordinatorV2ErrorReturned) Topic() common.Hash {
	return common.HexToHash("0x4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e")
}

func (BatchVRFCoordinatorV2RawErrorReturned) Topic() common.Hash {
	return common.HexToHash("0xbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b5")
}

func (_BatchVRFCoordinatorV2 *BatchVRFCoordinatorV2) Address() common.Address {
	return _BatchVRFCoordinatorV2.address
}

type BatchVRFCoordinatorV2Interface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	FulfillRandomWords(opts *bind.TransactOpts, proofs []VRFTypesProof, rcs []VRFTypesRequestCommitment) (*types.Transaction, error)

	FilterErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2ErrorReturnedIterator, error)

	WatchErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2ErrorReturned, requestId []*big.Int) (event.Subscription, error)

	ParseErrorReturned(log types.Log) (*BatchVRFCoordinatorV2ErrorReturned, error)

	FilterRawErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2RawErrorReturnedIterator, error)

	WatchRawErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2RawErrorReturned, requestId []*big.Int) (event.Subscription, error)

	ParseRawErrorReturned(log types.Log) (*BatchVRFCoordinatorV2RawErrorReturned, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
