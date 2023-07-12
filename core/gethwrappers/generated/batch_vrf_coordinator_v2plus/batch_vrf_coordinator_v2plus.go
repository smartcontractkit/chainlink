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
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	NativePayment    bool
}

var BatchVRFCoordinatorV2PlusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"ErrorReturned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"RawErrorReturned\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRFTypes.Proof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus[]\",\"name\":\"rcs\",\"type\":\"tuple[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610be1380380610be183398101604081905261002f91610044565b60601b6001600160601b031916608052610074565b60006020828403121561005657600080fd5b81516001600160a01b038116811461006d57600080fd5b9392505050565b60805160601c610b49610098600039600081816040015261011d0152610b496000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633b2bcbf11461003b5780634c34f6ef1461008b575b600080fd5b6100627f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61009e610099366004610596565b6100a0565b005b805182511461010f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b8251811015610330577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663827f1c5d84838151811061016957610169610a1b565b602002602001015184848151811061018357610183610a1b565b60200260200101516040518363ffffffff1660e01b81526004016101a89291906107ec565b602060405180830381600087803b1580156101c257600080fd5b505af1925050508015610210575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261020d918101906106fa565b60015b61031c5761021c610a79565b806308c379a014156102a15750610231610a94565b8061023c57506102a3565b600061026085848151811061025357610253610a1b565b6020026020010151610335565b9050807f4dcab4ce0e741a040f7e0f9b880557f8de685a9520d4bfac272a81c3c3802b2e8360405161029291906107d9565b60405180910390a2505061031e565b505b3d8080156102cd576040519150601f19603f3d011682016040523d82523d6000602084013e6102d2565b606091505b5060006102ea85848151811061025357610253610a1b565b9050807fbfd42bb5a1bf8153ea750f66ea4944f23f7b9ae51d0462177b9769aa652b61b58360405161029291906107d9565b505b80610328816109bb565b915050610112565b505050565b60008061034583600001516103a4565b9050808360800151604051602001610367929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b6000816040516020016103b791906107c5565b604051602081830303815290604052805190602001209050919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146103f857600080fd5b919050565b600082601f83011261040e57600080fd5b8135602061041b82610905565b604080516104298382610970565b848152838101925086840160c0808702890186018a101561044957600080fd5b60005b878110156104e35781838c03121561046357600080fd5b845161046e81610929565b6104778461057e565b815261048488850161057e565b8882015261049386850161056a565b8682015260606104a481860161056a565b9082015260806104b58582016103d4565b9082015260a08481013580151581146104cd57600080fd5b908201528652948601949181019160010161044c565b50919998505050505050505050565b600082601f83011261050357600080fd5b6040516040810181811067ffffffffffffffff8211171561052657610526610a4a565b806040525080838560408601111561053d57600080fd5b60005b600281101561055f578135835260209283019290910190600101610540565b509195945050505050565b803563ffffffff811681146103f857600080fd5b803567ffffffffffffffff811681146103f857600080fd5b600080604083850312156105a957600080fd5b823567ffffffffffffffff808211156105c157600080fd5b818501915085601f8301126105d557600080fd5b813560206105e282610905565b6040516105ef8282610970565b83815282810191508583016101a0808602880185018c101561061057600080fd5b600097505b858810156106ca5780828d03121561062c57600080fd5b6106346108f6565b61063e8d846104f2565b815261064d8d604085016104f2565b868201526080830135604082015260a0830135606082015260c0830135608082015260e061067c8185016103d4565b60a08301526101006106908f8287016104f2565b60c08401526106a38f61014087016104f2565b91830191909152610180840135908201528452600197909701969284019290810190610615565b509097505050860135925050808211156106e357600080fd5b506106f0858286016103fd565b9150509250929050565b60006020828403121561070c57600080fd5b81516bffffffffffffffffffffffff8116811461072857600080fd5b9392505050565b8060005b6002811015610752578151845260209384019390910190600101610733565b50505050565b60008151808452602060005b8281101561077f578481018201518682018301528101610764565b828111156107905760008284880101525b50807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401168601019250505092915050565b604081016107d3828461072f565b92915050565b6020815260006107286020830184610758565b60006102608201905061080082855161072f565b6020840151610812604084018261072f565b5060408401516080830152606084015160a0830152608084015160c083015273ffffffffffffffffffffffffffffffffffffffff60a08501511660e083015260c08401516101006108658185018361072f565b60e0860151915061087a61014085018361072f565b85015161018084015250825167ffffffffffffffff9081166101a08401526020840151166101c0830152604083015163ffffffff9081166101e0840152606084015116610200830152608083015173ffffffffffffffffffffffffffffffffffffffff1661022083015260a08301511515610240830152610728565b6040516109028161094f565b90565b600067ffffffffffffffff82111561091f5761091f610a4a565b5060051b60200190565b60c0810181811067ffffffffffffffff8211171561094957610949610a4a565b60405250565b610120810167ffffffffffffffff8111828210171561094957610949610a4a565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff821117156109b4576109b4610a4a565b6040525050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610a14577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060033d11156109025760046000803e5060005160e01c90565b600060443d1015610aa25790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff8160248401118184111715610af057505050505090565b8285019150815181811115610b085750505050505090565b843d8701016020828501011115610b225750505050505090565b610b3160208286010187610970565b50909594505050505056fea164736f6c6343000806000a",
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
