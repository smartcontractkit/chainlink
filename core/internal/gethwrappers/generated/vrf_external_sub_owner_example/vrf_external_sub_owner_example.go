// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_external_sub_owner_example

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

var VRFExternalSubOwnerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516106b43803806106b483398101604081905261002f9161009e565b6001600160601b0319606083901b16608052600080546001600160a01b03199081166001600160a01b039485161790915560018054929093169181169190911790915560048054909116331790556100d1565b80516001600160a01b038116811461009957600080fd5b919050565b600080604083850312156100b157600080fd5b6100ba83610082565b91506100c860208401610082565b90509250929050565b60805160601c6105bf6100f56000396000818160ed015261015501526105bf6000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c8063e89e106a11610050578063e89e106a14610094578063f2fde38b146100af578063f6eaffc8146100c257600080fd5b80631fe543e31461006c5780639561f02314610081575b600080fd5b61007f61007a366004610420565b6100d5565b005b61007f61008f36600461050f565b610194565b61009d60035481565b60405190815260200160405180910390f35b61007f6100bd3660046103b1565b610294565b61009d6100d03660046103ee565b6102ff565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610186576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b6101908282610320565b5050565b60045473ffffffffffffffffffffffffffffffffffffffff1633146101b857600080fd5b6000546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810183905267ffffffffffffffff8716602482015261ffff8516604482015263ffffffff80871660648301528416608482015273ffffffffffffffffffffffffffffffffffffffff90911690635d3b1d309060a401602060405180830381600087803b15801561025257600080fd5b505af1158015610266573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061028a9190610407565b6003555050505050565b60045473ffffffffffffffffffffffffffffffffffffffff1633146102b857600080fd5b600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6002818154811061030f57600080fd5b600091825260209091200154905081565b8051610333906002906020840190610338565b505050565b828054828255906000526020600020908101928215610373579160200282015b82811115610373578251825591602001919060010190610358565b5061037f929150610383565b5090565b5b8082111561037f5760008155600101610384565b803563ffffffff811681146103ac57600080fd5b919050565b6000602082840312156103c357600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146103e757600080fd5b9392505050565b60006020828403121561040057600080fd5b5035919050565b60006020828403121561041957600080fd5b5051919050565b6000806040838503121561043357600080fd5b8235915060208084013567ffffffffffffffff8082111561045357600080fd5b818601915086601f83011261046757600080fd5b81358181111561047957610479610583565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156104bc576104bc610583565b604052828152858101935084860182860187018b10156104db57600080fd5b600095505b838610156104fe5780358552600195909501949386019386016104e0565b508096505050505050509250929050565b600080600080600060a0868803121561052757600080fd5b853567ffffffffffffffff8116811461053f57600080fd5b945061054d60208701610398565b9350604086013561ffff8116811461056457600080fd5b925061057260608701610398565b949793965091946080013592915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFExternalSubOwnerExampleABI = VRFExternalSubOwnerExampleMetaData.ABI

var VRFExternalSubOwnerExampleBin = VRFExternalSubOwnerExampleMetaData.Bin

func DeployVRFExternalSubOwnerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFExternalSubOwnerExample, error) {
	parsed, err := VRFExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFExternalSubOwnerExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFExternalSubOwnerExample{VRFExternalSubOwnerExampleCaller: VRFExternalSubOwnerExampleCaller{contract: contract}, VRFExternalSubOwnerExampleTransactor: VRFExternalSubOwnerExampleTransactor{contract: contract}, VRFExternalSubOwnerExampleFilterer: VRFExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

type VRFExternalSubOwnerExample struct {
	address common.Address
	abi     abi.ABI
	VRFExternalSubOwnerExampleCaller
	VRFExternalSubOwnerExampleTransactor
	VRFExternalSubOwnerExampleFilterer
}

type VRFExternalSubOwnerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFExternalSubOwnerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFExternalSubOwnerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFExternalSubOwnerExampleSession struct {
	Contract     *VRFExternalSubOwnerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFExternalSubOwnerExampleCallerSession struct {
	Contract *VRFExternalSubOwnerExampleCaller
	CallOpts bind.CallOpts
}

type VRFExternalSubOwnerExampleTransactorSession struct {
	Contract     *VRFExternalSubOwnerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFExternalSubOwnerExampleRaw struct {
	Contract *VRFExternalSubOwnerExample
}

type VRFExternalSubOwnerExampleCallerRaw struct {
	Contract *VRFExternalSubOwnerExampleCaller
}

type VRFExternalSubOwnerExampleTransactorRaw struct {
	Contract *VRFExternalSubOwnerExampleTransactor
}

func NewVRFExternalSubOwnerExample(address common.Address, backend bind.ContractBackend) (*VRFExternalSubOwnerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFExternalSubOwnerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExample{address: address, abi: abi, VRFExternalSubOwnerExampleCaller: VRFExternalSubOwnerExampleCaller{contract: contract}, VRFExternalSubOwnerExampleTransactor: VRFExternalSubOwnerExampleTransactor{contract: contract}, VRFExternalSubOwnerExampleFilterer: VRFExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

func NewVRFExternalSubOwnerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFExternalSubOwnerExampleCaller, error) {
	contract, err := bindVRFExternalSubOwnerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExampleCaller{contract: contract}, nil
}

func NewVRFExternalSubOwnerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFExternalSubOwnerExampleTransactor, error) {
	contract, err := bindVRFExternalSubOwnerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExampleTransactor{contract: contract}, nil
}

func NewVRFExternalSubOwnerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFExternalSubOwnerExampleFilterer, error) {
	contract, err := bindVRFExternalSubOwnerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExampleFilterer{contract: contract}, nil
}

func bindVRFExternalSubOwnerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFExternalSubOwnerExample.Contract.VRFExternalSubOwnerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.VRFExternalSubOwnerExampleTransactor.contract.Transfer(opts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.VRFExternalSubOwnerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFExternalSubOwnerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.contract.Transfer(opts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFExternalSubOwnerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRandomWords(&_VRFExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRandomWords(&_VRFExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFExternalSubOwnerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRequestId(&_VRFExternalSubOwnerExample.CallOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRequestId(&_VRFExternalSubOwnerExample.CallOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.contract.Transact(opts, "requestRandomWords", subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFExternalSubOwnerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFExternalSubOwnerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.TransferOwnership(&_VRFExternalSubOwnerExample.TransactOpts, newOwner)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.TransferOwnership(&_VRFExternalSubOwnerExample.TransactOpts, newOwner)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExample) Address() common.Address {
	return _VRFExternalSubOwnerExample.address
}

type VRFExternalSubOwnerExampleInterface interface {
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	Address() common.Address
}
