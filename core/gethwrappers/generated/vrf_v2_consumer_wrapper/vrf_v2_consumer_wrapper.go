// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2_consumer_wrapper

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

var VRFv2ConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requestIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610e31380380610e3183398101604081905261002f9161019a565b6001600160601b0319606082901b1660805233806000816100975760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100c7576100c7816100f0565b5050600380546001600160a01b0319166001600160a01b039390931692909217909155506101ca565b6001600160a01b0381163314156101495760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161008e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156101ac57600080fd5b81516001600160a01b03811681146101c357600080fd5b9392505050565b60805160601c610c426101ef600039600081816101be01526102260152610c426000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80639561f02311610076578063d8a4676f1161005b578063d8a4676f14610169578063f2fde38b1461018a578063fc2a88c31461019d57600080fd5b80639561f02314610113578063a168fa891461012657600080fd5b80631fe543e3146100a857806379ba5097146100bd5780638796ba8c146100c55780638da5cb5b146100eb575b600080fd5b6100bb6100b6366004610a2c565b6101a6565b005b6100bb610266565b6100d86100d33660046109fa565b610363565b6040519081526020015b60405180910390f35b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e2565b6100d8610121366004610b1b565b610384565b6101526101343660046109fa565b60026020526000908152604090205460ff8082169161010090041682565b6040805192151583529015156020830152016100e2565b61017c6101773660046109fa565b610593565b6040516100e2929190610bca565b6100bb6101983660046109bd565b6106ad565b6100d860055481565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610258576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61026282826106c1565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161024f565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6004818154811061037357600080fd5b600091825260209091200154905081565b600061038e6107cb565b6003546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810184905267ffffffffffffffff8816602482015261ffff8616604482015263ffffffff80881660648301528516608482015273ffffffffffffffffffffffffffffffffffffffff90911690635d3b1d309060a401602060405180830381600087803b15801561042857600080fd5b505af115801561043c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104609190610a13565b604080516060810182526000808252600160208084018281528551848152808301875285870190815287855260028352959093208451815494517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161761010094151594909402939093178355935180519596509294919361050f9391850192910190610944565b5050600480546001810182556000919091527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b018290555060058190556040805182815263ffffffff851660208201527fcc58b13ad3eab50626c6a6300b1d139cd6ebb1688a7cced9461c2f7e762665ee910160405180910390a195945050505050565b600081815260026020526040812054606090610100900460ff16610613576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015260640161024f565b60008381526002602090815260408083208151606081018352815460ff808216151583526101009091041615158185015260018201805484518187028101870186528181529295939486019383018282801561068e57602002820191906000526020600020905b81548152602001906001019080831161067a575b5050505050815250509050806000015181604001519250925050915091565b6106b56107cb565b6106be8161084e565b50565b600082815260026020526040902054610100900460ff1661073e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015260640161024f565b600082815260026020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660019081178255835161078d939290910191840190610944565b507ffe2e2d779dba245964d4e3ef9b994be63856fd568bf7d3ca9e224755cb1bd54d82826040516107bf929190610bed565b60405180910390a15050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461084c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161024f565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156108ce576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161024f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821561097f579160200282015b8281111561097f578251825591602001919060010190610964565b5061098b92915061098f565b5090565b5b8082111561098b5760008155600101610990565b803563ffffffff811681146109b857600080fd5b919050565b6000602082840312156109cf57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146109f357600080fd5b9392505050565b600060208284031215610a0c57600080fd5b5035919050565b600060208284031215610a2557600080fd5b5051919050565b60008060408385031215610a3f57600080fd5b8235915060208084013567ffffffffffffffff80821115610a5f57600080fd5b818601915086601f830112610a7357600080fd5b813581811115610a8557610a85610c06565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610ac857610ac8610c06565b604052828152858101935084860182860187018b1015610ae757600080fd5b600095505b83861015610b0a578035855260019590950194938601938601610aec565b508096505050505050509250929050565b600080600080600060a08688031215610b3357600080fd5b853567ffffffffffffffff81168114610b4b57600080fd5b9450610b59602087016109a4565b9350604086013561ffff81168114610b7057600080fd5b9250610b7e606087016109a4565b949793965091946080013592915050565b600081518084526020808501945080840160005b83811015610bbf57815187529582019590820190600101610ba3565b509495945050505050565b8215158152604060208201526000610be56040830184610b8f565b949350505050565b828152604060208201526000610be56040830184610b8f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFv2ConsumerABI = VRFv2ConsumerMetaData.ABI

var VRFv2ConsumerBin = VRFv2ConsumerMetaData.Bin

func DeployVRFv2Consumer(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFv2Consumer, error) {
	parsed, err := VRFv2ConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFv2ConsumerBin), backend, vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFv2Consumer{address: address, abi: *parsed, VRFv2ConsumerCaller: VRFv2ConsumerCaller{contract: contract}, VRFv2ConsumerTransactor: VRFv2ConsumerTransactor{contract: contract}, VRFv2ConsumerFilterer: VRFv2ConsumerFilterer{contract: contract}}, nil
}

type VRFv2Consumer struct {
	address common.Address
	abi     abi.ABI
	VRFv2ConsumerCaller
	VRFv2ConsumerTransactor
	VRFv2ConsumerFilterer
}

type VRFv2ConsumerCaller struct {
	contract *bind.BoundContract
}

type VRFv2ConsumerTransactor struct {
	contract *bind.BoundContract
}

type VRFv2ConsumerFilterer struct {
	contract *bind.BoundContract
}

type VRFv2ConsumerSession struct {
	Contract     *VRFv2Consumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFv2ConsumerCallerSession struct {
	Contract *VRFv2ConsumerCaller
	CallOpts bind.CallOpts
}

type VRFv2ConsumerTransactorSession struct {
	Contract     *VRFv2ConsumerTransactor
	TransactOpts bind.TransactOpts
}

type VRFv2ConsumerRaw struct {
	Contract *VRFv2Consumer
}

type VRFv2ConsumerCallerRaw struct {
	Contract *VRFv2ConsumerCaller
}

type VRFv2ConsumerTransactorRaw struct {
	Contract *VRFv2ConsumerTransactor
}

func NewVRFv2Consumer(address common.Address, backend bind.ContractBackend) (*VRFv2Consumer, error) {
	abi, err := abi.JSON(strings.NewReader(VRFv2ConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFv2Consumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFv2Consumer{address: address, abi: abi, VRFv2ConsumerCaller: VRFv2ConsumerCaller{contract: contract}, VRFv2ConsumerTransactor: VRFv2ConsumerTransactor{contract: contract}, VRFv2ConsumerFilterer: VRFv2ConsumerFilterer{contract: contract}}, nil
}

func NewVRFv2ConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFv2ConsumerCaller, error) {
	contract, err := bindVRFv2Consumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerCaller{contract: contract}, nil
}

func NewVRFv2ConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFv2ConsumerTransactor, error) {
	contract, err := bindVRFv2Consumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerTransactor{contract: contract}, nil
}

func NewVRFv2ConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFv2ConsumerFilterer, error) {
	contract, err := bindVRFv2Consumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerFilterer{contract: contract}, nil
}

func bindVRFv2Consumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFv2ConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFv2Consumer *VRFv2ConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFv2Consumer.Contract.VRFv2ConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFv2Consumer *VRFv2ConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.VRFv2ConsumerTransactor.contract.Transfer(opts)
}

func (_VRFv2Consumer *VRFv2ConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.VRFv2ConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFv2Consumer *VRFv2ConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFv2Consumer.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.contract.Transfer(opts)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.contract.Transact(opts, method, params...)
}

func (_VRFv2Consumer *VRFv2ConsumerCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

func (_VRFv2Consumer *VRFv2ConsumerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFv2Consumer.Contract.GetRequestStatus(&_VRFv2Consumer.CallOpts, _requestId)
}

func (_VRFv2Consumer *VRFv2ConsumerCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFv2Consumer.Contract.GetRequestStatus(&_VRFv2Consumer.CallOpts, _requestId)
}

func (_VRFv2Consumer *VRFv2ConsumerCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFv2Consumer *VRFv2ConsumerSession) LastRequestId() (*big.Int, error) {
	return _VRFv2Consumer.Contract.LastRequestId(&_VRFv2Consumer.CallOpts)
}

func (_VRFv2Consumer *VRFv2ConsumerCallerSession) LastRequestId() (*big.Int, error) {
	return _VRFv2Consumer.Contract.LastRequestId(&_VRFv2Consumer.CallOpts)
}

func (_VRFv2Consumer *VRFv2ConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFv2Consumer *VRFv2ConsumerSession) Owner() (common.Address, error) {
	return _VRFv2Consumer.Contract.Owner(&_VRFv2Consumer.CallOpts)
}

func (_VRFv2Consumer *VRFv2ConsumerCallerSession) Owner() (common.Address, error) {
	return _VRFv2Consumer.Contract.Owner(&_VRFv2Consumer.CallOpts)
}

func (_VRFv2Consumer *VRFv2ConsumerCaller) RequestIds(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "requestIds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFv2Consumer *VRFv2ConsumerSession) RequestIds(arg0 *big.Int) (*big.Int, error) {
	return _VRFv2Consumer.Contract.RequestIds(&_VRFv2Consumer.CallOpts, arg0)
}

func (_VRFv2Consumer *VRFv2ConsumerCallerSession) RequestIds(arg0 *big.Int) (*big.Int, error) {
	return _VRFv2Consumer.Contract.RequestIds(&_VRFv2Consumer.CallOpts, arg0)
}

func (_VRFv2Consumer *VRFv2ConsumerCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Exists = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFv2Consumer *VRFv2ConsumerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFv2Consumer.Contract.SRequests(&_VRFv2Consumer.CallOpts, arg0)
}

func (_VRFv2Consumer *VRFv2ConsumerCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFv2Consumer.Contract.SRequests(&_VRFv2Consumer.CallOpts, arg0)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "acceptOwnership")
}

func (_VRFv2Consumer *VRFv2ConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.AcceptOwnership(&_VRFv2Consumer.TransactOpts)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.AcceptOwnership(&_VRFv2Consumer.TransactOpts)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFv2Consumer *VRFv2ConsumerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RawFulfillRandomWords(&_VRFv2Consumer.TransactOpts, requestId, randomWords)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RawFulfillRandomWords(&_VRFv2Consumer.TransactOpts, requestId, randomWords)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactor) RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "requestRandomWords", subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

func (_VRFv2Consumer *VRFv2ConsumerSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RequestRandomWords(&_VRFv2Consumer.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RequestRandomWords(&_VRFv2Consumer.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFv2Consumer *VRFv2ConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.TransferOwnership(&_VRFv2Consumer.TransactOpts, to)
}

func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.TransferOwnership(&_VRFv2Consumer.TransactOpts, to)
}

type VRFv2ConsumerOwnershipTransferRequestedIterator struct {
	Event *VRFv2ConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFv2ConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerOwnershipTransferRequested)
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
		it.Event = new(VRFv2ConsumerOwnershipTransferRequested)
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

func (it *VRFv2ConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFv2ConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFv2ConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFv2ConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerOwnershipTransferRequestedIterator{contract: _VRFv2Consumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFv2ConsumerOwnershipTransferRequested)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFv2ConsumerOwnershipTransferRequested, error) {
	event := new(VRFv2ConsumerOwnershipTransferRequested)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFv2ConsumerOwnershipTransferredIterator struct {
	Event *VRFv2ConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFv2ConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerOwnershipTransferred)
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
		it.Event = new(VRFv2ConsumerOwnershipTransferred)
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

func (it *VRFv2ConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFv2ConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFv2ConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFv2ConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerOwnershipTransferredIterator{contract: _VRFv2Consumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFv2ConsumerOwnershipTransferred)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFv2ConsumerOwnershipTransferred, error) {
	event := new(VRFv2ConsumerOwnershipTransferred)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFv2ConsumerRequestFulfilledIterator struct {
	Event *VRFv2ConsumerRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFv2ConsumerRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerRequestFulfilled)
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
		it.Event = new(VRFv2ConsumerRequestFulfilled)
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

func (it *VRFv2ConsumerRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFv2ConsumerRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFv2ConsumerRequestFulfilled struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Raw         types.Log
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterRequestFulfilled(opts *bind.FilterOpts) (*VRFv2ConsumerRequestFulfilledIterator, error) {

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "RequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerRequestFulfilledIterator{contract: _VRFv2Consumer.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "RequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFv2ConsumerRequestFulfilled)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseRequestFulfilled(log types.Log) (*VRFv2ConsumerRequestFulfilled, error) {
	event := new(VRFv2ConsumerRequestFulfilled)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFv2ConsumerRequestSentIterator struct {
	Event *VRFv2ConsumerRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFv2ConsumerRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerRequestSent)
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
		it.Event = new(VRFv2ConsumerRequestSent)
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

func (it *VRFv2ConsumerRequestSentIterator) Error() error {
	return it.fail
}

func (it *VRFv2ConsumerRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFv2ConsumerRequestSent struct {
	RequestId *big.Int
	NumWords  uint32
	Raw       types.Log
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterRequestSent(opts *bind.FilterOpts) (*VRFv2ConsumerRequestSentIterator, error) {

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "RequestSent")
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerRequestSentIterator{contract: _VRFv2Consumer.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerRequestSent) (event.Subscription, error) {

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "RequestSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFv2ConsumerRequestSent)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseRequestSent(log types.Log) (*VRFv2ConsumerRequestSent, error) {
	event := new(VRFv2ConsumerRequestSent)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Fulfilled   bool
	RandomWords []*big.Int
}
type SRequests struct {
	Fulfilled bool
	Exists    bool
}

func (_VRFv2Consumer *VRFv2Consumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFv2Consumer.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFv2Consumer.ParseOwnershipTransferRequested(log)
	case _VRFv2Consumer.abi.Events["OwnershipTransferred"].ID:
		return _VRFv2Consumer.ParseOwnershipTransferred(log)
	case _VRFv2Consumer.abi.Events["RequestFulfilled"].ID:
		return _VRFv2Consumer.ParseRequestFulfilled(log)
	case _VRFv2Consumer.abi.Events["RequestSent"].ID:
		return _VRFv2Consumer.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFv2ConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFv2ConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFv2ConsumerRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0xfe2e2d779dba245964d4e3ef9b994be63856fd568bf7d3ca9e224755cb1bd54d")
}

func (VRFv2ConsumerRequestSent) Topic() common.Hash {
	return common.HexToHash("0xcc58b13ad3eab50626c6a6300b1d139cd6ebb1688a7cced9461c2f7e762665ee")
}

func (_VRFv2Consumer *VRFv2Consumer) Address() common.Address {
	return _VRFv2Consumer.address
}

type VRFv2ConsumerInterface interface {
	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	LastRequestId(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	RequestIds(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFv2ConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFv2ConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFv2ConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFv2ConsumerOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts) (*VRFv2ConsumerRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerRequestFulfilled) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*VRFv2ConsumerRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts) (*VRFv2ConsumerRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerRequestSent) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*VRFv2ConsumerRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
