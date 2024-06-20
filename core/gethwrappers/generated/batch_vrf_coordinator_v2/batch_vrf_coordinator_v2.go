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
	Bin: "0x60a060405234801561001057600080fd5b50604051610b6e380380610b6e83398101604081905261002f91610044565b60601b6001600160601b031916608052610074565b60006020828403121561005657600080fd5b81516001600160a01b038116811461006d57600080fd5b9392505050565b60805160601c610ad6610098600039600081816055015261011b0152610ad66000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806308b2da0a1461003b5780633b2bcbf114610050575b600080fd5b61004e6100493660046104ae565b6100a0565b005b6100777f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b8051821461010e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b8281101561033c577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663af198b97858584818110610167576101676109a7565b90506101a00201848481518110610180576101806109a7565b60200260200101516040518363ffffffff1660e01b81526004016101a59291906107b9565b602060405180830381600087803b1580156101bf57600080fd5b505af192505050801561020d575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261020a918101906106c5565b60015b61032857610219610a05565b806308c379a014156102ad575061022e610a21565b8061023957506102af565b600061026c868685818110610250576102506109a7565b90506101a002018036038101906102679190610621565b610342565b9050807f4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e8360405161029e91906107a6565b60405180910390a2505061032a565b505b3d8080156102d9576040519150601f19603f3d011682016040523d82523d6000602084013e6102de565b606091505b5060006102f6868685818110610250576102506109a7565b9050807fbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b58360405161029e91906107a6565b505b8061033481610947565b915050610111565b50505050565b60008061035283600001516103b1565b9050808360800151604051602001610374929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b6000816040516020016103c49190610775565b604051602081830303815290604052805190602001209050919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461040557600080fd5b919050565b600082601f83011261041b57600080fd5b6040516040810181811067ffffffffffffffff8211171561043e5761043e6109d6565b806040525080838560408601111561045557600080fd5b60005b6002811015610477578135835260209283019290910190600101610458565b509195945050505050565b803563ffffffff8116811461040557600080fd5b803567ffffffffffffffff8116811461040557600080fd5b600080600060408085870312156104c457600080fd5b843567ffffffffffffffff808211156104dc57600080fd5b818701915087601f8301126104f057600080fd5b8135818111156104ff57600080fd5b602089816101a08402860101111561051657600080fd5b80840197508196508089013593508284111561053157600080fd5b838901935089601f85011261054557600080fd5b8335915082821115610559576105596109d6565b8451925061056c818360051b01846108fc565b81835280830184820160a0808502870184018d101561058a57600080fd5b60009650865b8581101561060e5781838f0312156105a6578788fd5b88516105b1816108d6565b6105ba84610496565b81526105c7868501610496565b868201526105d68a8501610482565b8a82015260606105e7818601610482565b9082015260806105f88582016103e1565b9082015284529284019291810191600101610590565b5050505050508093505050509250925092565b60006101a0828403121561063457600080fd5b61063c6108ac565b610646848461040a565b8152610655846040850161040a565b60208201526080830135604082015260a0830135606082015260c0830135608082015261068460e084016103e1565b60a08201526101006106988582860161040a565b60c08301526106ab85610140860161040a565b60e083015261018084013581830152508091505092915050565b6000602082840312156106d757600080fd5b81516bffffffffffffffffffffffff811681146106f357600080fd5b9392505050565b6040818337600060408301525050565b6000815180845260005b8181101561073057602081850181015186830182015201610714565b81811115610742576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60408101818360005b600281101561079d57815183526020928301929091019060010161077e565b50505092915050565b6020815260006106f3602083018461070a565b61024081016040848337604082016000815260408086018237506080848101359083015260a0808501359083015260c0808501359083015273ffffffffffffffffffffffffffffffffffffffff61081260e086016103e1565b1660e08301526101006108298184018287016106fa565b5061014061083b8184018287016106fa565b506101808481013590830152825167ffffffffffffffff9081166101a08401526020840151166101c0830152604083015163ffffffff9081166101e0840152606084015116610200830152608083015173ffffffffffffffffffffffffffffffffffffffff166102208301526106f3565b604051610120810167ffffffffffffffff811182821017156108d0576108d06109d6565b60405290565b60a0810181811067ffffffffffffffff821117156108f6576108f66109d6565b60405250565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff82111715610940576109406109d6565b6040525050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156109a0577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060033d1115610a1e5760046000803e5060005160e01c5b90565b600060443d1015610a2f5790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff8160248401118184111715610a7d57505050505090565b8285019150815181811115610a955750505050505090565b843d8701016020828501011115610aaf5750505050505090565b610abe602082860101876108fc565b50909594505050505056fea164736f6c6343000806000a",
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
	return address, tx, &BatchVRFCoordinatorV2{address: address, abi: *parsed, BatchVRFCoordinatorV2Caller: BatchVRFCoordinatorV2Caller{contract: contract}, BatchVRFCoordinatorV2Transactor: BatchVRFCoordinatorV2Transactor{contract: contract}, BatchVRFCoordinatorV2Filterer: BatchVRFCoordinatorV2Filterer{contract: contract}}, nil
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
	parsed, err := BatchVRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
