// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package batch_vrf_coordinator_v2plus

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

type VRFTypesRequestCommitmentV2Plus struct {
	BlockNum         uint64
	SubId            *big.Int
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	ExtraArgs        []byte
}

var BatchVRFCoordinatorV2PlusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"ErrorReturned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"RawErrorReturned\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2PlusFulfill\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRFTypes.Proof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus[]\",\"name\":\"rcs\",\"type\":\"tuple[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610ba2380380610ba283398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b608051610b11610091600039600081816040015261011a0152610b116000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633b2bcbf11461003b5780636abb17211461008b575b600080fd5b6100627f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61009e6100993660046103dc565b6100a0565b005b82811461010d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b83811015610336577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663301f42e9868684818110610166576101666104a5565b90506101a0020185858581811061017f5761017f6104a5565b905060200281019061019191906104d4565b60006040518463ffffffff1660e01b81526004016101b19392919061069b565b6020604051808303816000875af192505050801561020a575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682019092526102079181019061074c565b60015b61032457610216610781565b806308c379a0036102a9575061022a610817565b8061023557506102ab565b600061026887878581811061024c5761024c6104a5565b90506101a002018036038101906102639190610959565b61033d565b9050807f4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e8360405161029a9190610a61565b60405180910390a25050610326565b505b3d8080156102d5576040519150601f19603f3d011682016040523d82523d6000602084013e6102da565b606091505b5060006102f287878581811061024c5761024c6104a5565b9050807fbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b58360405161029a9190610a61565b505b61032f81610a74565b9050610110565b5050505050565b60008061034d83600001516103ac565b905080836080015160405160200161036f929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b6000816040516020016103bf9190610ad3565b604051602081830303815290604052805190602001209050919050565b600080600080604085870312156103f257600080fd5b843567ffffffffffffffff8082111561040a57600080fd5b818701915087601f83011261041e57600080fd5b81358181111561042d57600080fd5b8860206101a08302850101111561044357600080fd5b60209283019650945090860135908082111561045e57600080fd5b818701915087601f83011261047257600080fd5b81358181111561048157600080fd5b8860208260051b850101111561049657600080fd5b95989497505060200194505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff4183360301811261050857600080fd5b9190910192915050565b60408183375050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461053f57600080fd5b919050565b803563ffffffff8116811461053f57600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000813567ffffffffffffffff8082168083146105bd57600080fd5b8552602084810135908601526105d560408501610544565b915063ffffffff8083166040870152806105f160608701610544565b1660608701525073ffffffffffffffffffffffffffffffffffffffff6106196080860161051b565b16608086015260a084013591507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301821261065657600080fd5b602091840191820191358181111561066d57600080fd5b80360383131561067c57600080fd5b60c060a087015261069160c087018285610558565b9695505050505050565b60006101e06040868437604080870160408501376080860135608084015260a086013560a084015260c086013560c084015273ffffffffffffffffffffffffffffffffffffffff6106ee60e0880161051b565b1660e084015261010060408188018286013750610140610712818501828901610512565b50610180808701358185015250806101a0840152610732818401866105a1565b9150506107446101c083018415159052565b949350505050565b60006020828403121561075e57600080fd5b81516bffffffffffffffffffffffff8116811461077a57600080fd5b9392505050565b600060033d111561079a5760046000803e5060005160e01c5b90565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff821117156108105761081061079d565b6040525050565b600060443d10156108255790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff816024840111818411171561087357505050505090565b828501915081518181111561088b5750505050505090565b843d87010160208285010111156108a55750505050505090565b6108b4602082860101876107cc565b509095945050505050565b604051610120810167ffffffffffffffff811182821017156108e3576108e361079d565b60405290565b600082601f8301126108fa57600080fd5b6040516040810181811067ffffffffffffffff8211171561091d5761091d61079d565b806040525080604084018581111561093457600080fd5b845b8181101561094e578035835260209283019201610936565b509195945050505050565b60006101a0828403121561096c57600080fd5b6109746108bf565b61097e84846108e9565b815261098d84604085016108e9565b60208201526080830135604082015260a0830135606082015260c083013560808201526109bc60e0840161051b565b60a08201526101006109d0858286016108e9565b60c08301526109e38561014086016108e9565b60e083015261018084013581830152508091505092915050565b6000815180845260005b81811015610a2357602081850181015186830182015201610a07565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061077a60208301846109fd565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610acc577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b60408101818360005b6002811015610afb578151835260209283019290910190600101610adc565b5050509291505056fea164736f6c6343000813000a",
}

var BatchVRFCoordinatorV2PlusABI = BatchVRFCoordinatorV2PlusMetaData.ABI

var BatchVRFCoordinatorV2PlusBin = BatchVRFCoordinatorV2PlusMetaData.Bin

func DeployBatchVRFCoordinatorV2Plus(auth *bind.TransactOpts, backend bind.ContractBackend, coordinatorAddr common.Address) (common.Address, *types.Transaction, *BatchVRFCoordinatorV2Plus, error) {
	parsed, err := BatchVRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchVRFCoordinatorV2PlusBin), backend, coordinatorAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchVRFCoordinatorV2Plus{address: address, abi: *parsed, BatchVRFCoordinatorV2PlusCaller: BatchVRFCoordinatorV2PlusCaller{contract: contract}, BatchVRFCoordinatorV2PlusTransactor: BatchVRFCoordinatorV2PlusTransactor{contract: contract}, BatchVRFCoordinatorV2PlusFilterer: BatchVRFCoordinatorV2PlusFilterer{contract: contract}}, nil
}

type BatchVRFCoordinatorV2Plus struct {
	address common.Address
	abi     abi.ABI
	BatchVRFCoordinatorV2PlusCaller
	BatchVRFCoordinatorV2PlusTransactor
	BatchVRFCoordinatorV2PlusFilterer
}

type BatchVRFCoordinatorV2PlusCaller struct {
	contract *bind.BoundContract
}

type BatchVRFCoordinatorV2PlusTransactor struct {
	contract *bind.BoundContract
}

type BatchVRFCoordinatorV2PlusFilterer struct {
	contract *bind.BoundContract
}

type BatchVRFCoordinatorV2PlusSession struct {
	Contract     *BatchVRFCoordinatorV2Plus
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BatchVRFCoordinatorV2PlusCallerSession struct {
	Contract *BatchVRFCoordinatorV2PlusCaller
	CallOpts bind.CallOpts
}

type BatchVRFCoordinatorV2PlusTransactorSession struct {
	Contract     *BatchVRFCoordinatorV2PlusTransactor
	TransactOpts bind.TransactOpts
}

type BatchVRFCoordinatorV2PlusRaw struct {
	Contract *BatchVRFCoordinatorV2Plus
}

type BatchVRFCoordinatorV2PlusCallerRaw struct {
	Contract *BatchVRFCoordinatorV2PlusCaller
}

type BatchVRFCoordinatorV2PlusTransactorRaw struct {
	Contract *BatchVRFCoordinatorV2PlusTransactor
}

func NewBatchVRFCoordinatorV2Plus(address common.Address, backend bind.ContractBackend) (*BatchVRFCoordinatorV2Plus, error) {
	abi, err := abi.JSON(strings.NewReader(BatchVRFCoordinatorV2PlusABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBatchVRFCoordinatorV2Plus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2Plus{address: address, abi: abi, BatchVRFCoordinatorV2PlusCaller: BatchVRFCoordinatorV2PlusCaller{contract: contract}, BatchVRFCoordinatorV2PlusTransactor: BatchVRFCoordinatorV2PlusTransactor{contract: contract}, BatchVRFCoordinatorV2PlusFilterer: BatchVRFCoordinatorV2PlusFilterer{contract: contract}}, nil
}

func NewBatchVRFCoordinatorV2PlusCaller(address common.Address, caller bind.ContractCaller) (*BatchVRFCoordinatorV2PlusCaller, error) {
	contract, err := bindBatchVRFCoordinatorV2Plus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2PlusCaller{contract: contract}, nil
}

func NewBatchVRFCoordinatorV2PlusTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchVRFCoordinatorV2PlusTransactor, error) {
	contract, err := bindBatchVRFCoordinatorV2Plus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2PlusTransactor{contract: contract}, nil
}

func NewBatchVRFCoordinatorV2PlusFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchVRFCoordinatorV2PlusFilterer, error) {
	contract, err := bindBatchVRFCoordinatorV2Plus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2PlusFilterer{contract: contract}, nil
}

func bindBatchVRFCoordinatorV2Plus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchVRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchVRFCoordinatorV2Plus.Contract.BatchVRFCoordinatorV2PlusCaller.contract.Call(opts, result, method, params...)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.BatchVRFCoordinatorV2PlusTransactor.contract.Transfer(opts)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.BatchVRFCoordinatorV2PlusTransactor.contract.Transact(opts, method, params...)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchVRFCoordinatorV2Plus.Contract.contract.Call(opts, result, method, params...)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.contract.Transfer(opts)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.contract.Transact(opts, method, params...)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BatchVRFCoordinatorV2Plus.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusSession) COORDINATOR() (common.Address, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.COORDINATOR(&_BatchVRFCoordinatorV2Plus.CallOpts)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusCallerSession) COORDINATOR() (common.Address, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.COORDINATOR(&_BatchVRFCoordinatorV2Plus.CallOpts)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusTransactor) FulfillRandomWords(opts *bind.TransactOpts, proofs []VRFTypesProof, rcs []VRFTypesRequestCommitmentV2Plus) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.contract.Transact(opts, "fulfillRandomWords", proofs, rcs)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusSession) FulfillRandomWords(proofs []VRFTypesProof, rcs []VRFTypesRequestCommitmentV2Plus) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.FulfillRandomWords(&_BatchVRFCoordinatorV2Plus.TransactOpts, proofs, rcs)
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusTransactorSession) FulfillRandomWords(proofs []VRFTypesProof, rcs []VRFTypesRequestCommitmentV2Plus) (*types.Transaction, error) {
	return _BatchVRFCoordinatorV2Plus.Contract.FulfillRandomWords(&_BatchVRFCoordinatorV2Plus.TransactOpts, proofs, rcs)
}

type BatchVRFCoordinatorV2PlusErrorReturnedIterator struct {
	Event *BatchVRFCoordinatorV2PlusErrorReturned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BatchVRFCoordinatorV2PlusErrorReturnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchVRFCoordinatorV2PlusErrorReturned)
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
		it.Event = new(BatchVRFCoordinatorV2PlusErrorReturned)
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

func (it *BatchVRFCoordinatorV2PlusErrorReturnedIterator) Error() error {
	return it.fail
}

func (it *BatchVRFCoordinatorV2PlusErrorReturnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BatchVRFCoordinatorV2PlusErrorReturned struct {
	RequestId *big.Int
	Reason    string
	Raw       types.Log
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusFilterer) FilterErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2PlusErrorReturnedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2Plus.contract.FilterLogs(opts, "ErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2PlusErrorReturnedIterator{contract: _BatchVRFCoordinatorV2Plus.contract, event: "ErrorReturned", logs: logs, sub: sub}, nil
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusFilterer) WatchErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2PlusErrorReturned, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2Plus.contract.WatchLogs(opts, "ErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BatchVRFCoordinatorV2PlusErrorReturned)
				if err := _BatchVRFCoordinatorV2Plus.contract.UnpackLog(event, "ErrorReturned", log); err != nil {
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

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusFilterer) ParseErrorReturned(log types.Log) (*BatchVRFCoordinatorV2PlusErrorReturned, error) {
	event := new(BatchVRFCoordinatorV2PlusErrorReturned)
	if err := _BatchVRFCoordinatorV2Plus.contract.UnpackLog(event, "ErrorReturned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BatchVRFCoordinatorV2PlusRawErrorReturnedIterator struct {
	Event *BatchVRFCoordinatorV2PlusRawErrorReturned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BatchVRFCoordinatorV2PlusRawErrorReturnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchVRFCoordinatorV2PlusRawErrorReturned)
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
		it.Event = new(BatchVRFCoordinatorV2PlusRawErrorReturned)
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

func (it *BatchVRFCoordinatorV2PlusRawErrorReturnedIterator) Error() error {
	return it.fail
}

func (it *BatchVRFCoordinatorV2PlusRawErrorReturnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BatchVRFCoordinatorV2PlusRawErrorReturned struct {
	RequestId    *big.Int
	LowLevelData []byte
	Raw          types.Log
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusFilterer) FilterRawErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2PlusRawErrorReturnedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2Plus.contract.FilterLogs(opts, "RawErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &BatchVRFCoordinatorV2PlusRawErrorReturnedIterator{contract: _BatchVRFCoordinatorV2Plus.contract, event: "RawErrorReturned", logs: logs, sub: sub}, nil
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusFilterer) WatchRawErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2PlusRawErrorReturned, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _BatchVRFCoordinatorV2Plus.contract.WatchLogs(opts, "RawErrorReturned", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BatchVRFCoordinatorV2PlusRawErrorReturned)
				if err := _BatchVRFCoordinatorV2Plus.contract.UnpackLog(event, "RawErrorReturned", log); err != nil {
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

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2PlusFilterer) ParseRawErrorReturned(log types.Log) (*BatchVRFCoordinatorV2PlusRawErrorReturned, error) {
	event := new(BatchVRFCoordinatorV2PlusRawErrorReturned)
	if err := _BatchVRFCoordinatorV2Plus.contract.UnpackLog(event, "RawErrorReturned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2Plus) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BatchVRFCoordinatorV2Plus.abi.Events["ErrorReturned"].ID:
		return _BatchVRFCoordinatorV2Plus.ParseErrorReturned(log)
	case _BatchVRFCoordinatorV2Plus.abi.Events["RawErrorReturned"].ID:
		return _BatchVRFCoordinatorV2Plus.ParseRawErrorReturned(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BatchVRFCoordinatorV2PlusErrorReturned) Topic() common.Hash {
	return common.HexToHash("0x4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e")
}

func (BatchVRFCoordinatorV2PlusRawErrorReturned) Topic() common.Hash {
	return common.HexToHash("0xbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b5")
}

func (_BatchVRFCoordinatorV2Plus *BatchVRFCoordinatorV2Plus) Address() common.Address {
	return _BatchVRFCoordinatorV2Plus.address
}

type BatchVRFCoordinatorV2PlusInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	FulfillRandomWords(opts *bind.TransactOpts, proofs []VRFTypesProof, rcs []VRFTypesRequestCommitmentV2Plus) (*types.Transaction, error)

	FilterErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2PlusErrorReturnedIterator, error)

	WatchErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2PlusErrorReturned, requestId []*big.Int) (event.Subscription, error)

	ParseErrorReturned(log types.Log) (*BatchVRFCoordinatorV2PlusErrorReturned, error)

	FilterRawErrorReturned(opts *bind.FilterOpts, requestId []*big.Int) (*BatchVRFCoordinatorV2PlusRawErrorReturnedIterator, error)

	WatchRawErrorReturned(opts *bind.WatchOpts, sink chan<- *BatchVRFCoordinatorV2PlusRawErrorReturned, requestId []*big.Int) (event.Subscription, error)

	ParseRawErrorReturned(log types.Log) (*BatchVRFCoordinatorV2PlusRawErrorReturned, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
