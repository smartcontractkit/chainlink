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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"ErrorReturned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"RawErrorReturned\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRFTypes.Proof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus[]\",\"name\":\"rcs\",\"type\":\"tuple[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610cc7380380610cc783398101604081905261002f91610044565b60601b6001600160601b031916608052610074565b60006020828403121561005657600080fd5b81516001600160a01b038116811461006d57600080fd5b9392505050565b60805160601c610c2f610098600039600081816040015261011d0152610c2f6000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633b2bcbf11461003b5780636abb17211461008b575b600080fd5b6100627f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61009e610099366004610665565b6100a0565b005b805182511461010f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b8251811015610330577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663330987b384838151811061016957610169610b00565b602002602001015184848151811061018357610183610b00565b60200260200101516040518363ffffffff1660e01b81526004016101a8929190610930565b602060405180830381600087803b1580156101c257600080fd5b505af1925050508015610210575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261020d918101906107c5565b60015b61031c5761021c610b5e565b806308c379a014156102a15750610231610b7a565b8061023c57506102a3565b600061026085848151811061025357610253610b00565b6020026020010151610335565b9050807f4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e83604051610292919061091d565b60405180910390a2505061031e565b505b3d8080156102cd576040519150601f19603f3d011682016040523d82523d6000602084013e6102d2565b606091505b5060006102ea85848151811061025357610253610b00565b9050807fbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b583604051610292919061091d565b505b8061032881610aa0565b915050610112565b505050565b60008061034583600001516103a4565b9050808360800151604051602001610367929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b6000816040516020016103b79190610909565b604051602081830303815290604052805190602001209050919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146103f857600080fd5b919050565b600082601f83011261040e57600080fd5b8135602061041b82610a0b565b604080516104298382610a55565b8481528381019250868401600586901b8801850189101561044957600080fd5b60005b8681101561053957813567ffffffffffffffff8082111561046c57600080fd5b818b01915060c0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848f030112156104a457600080fd5b86516104af81610a2f565b8984013583811681146104c157600080fd5b8152838801358a82015260606104d8818601610651565b8983015260806104e9818701610651565b8284015260a091506104fc8287016103d4565b9083015291840135918383111561051257600080fd5b6105208f8c858801016105bf565b908201528852505050938501939085019060010161044c565b509098975050505050505050565b600082601f83011261055857600080fd5b6040516040810181811067ffffffffffffffff8211171561057b5761057b610b2f565b806040525080838560408601111561059257600080fd5b60005b60028110156105b4578135835260209283019290910190600101610595565b509195945050505050565b600082601f8301126105d057600080fd5b813567ffffffffffffffff8111156105ea576105ea610b2f565b60405161061f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160182610a55565b81815284602083860101111561063457600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff811681146103f857600080fd5b600080604080848603121561067957600080fd5b833567ffffffffffffffff8082111561069157600080fd5b818601915086601f8301126106a557600080fd5b813560206106b282610a0b565b85516106be8282610a55565b83815282810191508583016101a0808602880185018d10156106df57600080fd5b600097505b858810156107945780828e0312156106fb57600080fd5b6107036109e1565b61070d8e84610547565b815261071b8e8b8501610547565b8682015260808301358a82015260a0830135606082015260c0830135608082015261074860e084016103d4565b60a082015261010061075c8f828601610547565b60c083015261076f8f6101408601610547565b60e08301526101808401359082015284526001979097019692840192908101906106e4565b509098505050870135935050808311156107ad57600080fd5b50506107bb858286016103fd565b9150509250929050565b6000602082840312156107d757600080fd5b81516bffffffffffffffffffffffff811681146107f357600080fd5b9392505050565b8060005b600281101561081d5781518452602093840193909101906001016107fe565b50505050565b6000815180845260005b818110156108495760208185018101518683018201520161082d565b8181111561085b576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b67ffffffffffffffff8151168252602081015160208301526000604082015163ffffffff8082166040860152806060850151166060860152505073ffffffffffffffffffffffffffffffffffffffff608083015116608084015260a082015160c060a085015261090160c0850182610823565b949350505050565b6040810161091782846107fa565b92915050565b6020815260006107f36020830184610823565b60006101c06109408386516107fa565b602085015161095260408501826107fa565b5060408501516080840152606085015160a0840152608085015160c084015273ffffffffffffffffffffffffffffffffffffffff60a08601511660e084015260c08501516101006109a5818601836107fa565b60e087015191506109ba6101408601836107fa565b860151610180850152506101a083018190526109d88184018561088e565b95945050505050565b604051610120810167ffffffffffffffff81118282101715610a0557610a05610b2f565b60405290565b600067ffffffffffffffff821115610a2557610a25610b2f565b5060051b60200190565b60c0810181811067ffffffffffffffff82111715610a4f57610a4f610b2f565b60405250565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff82111715610a9957610a99610b2f565b6040525050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610af9577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060033d1115610b775760046000803e5060005160e01c5b90565b600060443d1015610b885790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff8160248401118184111715610bd657505050505090565b8285019150815181811115610bee5750505050505090565b843d8701016020828501011115610c085750505050505090565b610c1760208286010187610a55565b50909594505050505056fea164736f6c6343000806000a",
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
	return address, tx, &BatchVRFCoordinatorV2Plus{BatchVRFCoordinatorV2PlusCaller: BatchVRFCoordinatorV2PlusCaller{contract: contract}, BatchVRFCoordinatorV2PlusTransactor: BatchVRFCoordinatorV2PlusTransactor{contract: contract}, BatchVRFCoordinatorV2PlusFilterer: BatchVRFCoordinatorV2PlusFilterer{contract: contract}}, nil
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
