// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keepers_vrf_consumer

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
	_ = abi.ConvertType
)

var KeepersVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"upkeepInterval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"KEY_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUBSCRIPTION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPKEEP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastTimeStamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfRequestCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfResponseCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x61014060405234801561001157600080fd5b50604051610bab380380610bab83398101604081905261003091610089565b60609490941b6001600160601b031916608081905260c09081529290921b6001600160c01b03191660e05260f09190911b6001600160f01b031916610100526101205260a0524260009081556001819055600255610109565b600080600080600060a086880312156100a157600080fd5b85516001600160a01b03811681146100b857600080fd5b60208701519095506001600160401b03811681146100d557600080fd5b60408701516060880151919550935061ffff811681146100f457600080fd5b80925050608086015190509295509295909350565b60805160601c60a05160c05160601c60e05160c01c6101005160f01c61012051610a17610194600039600081816101d501526105370152600081816101fc015261059001526000818160de015261056601526000818161017601526105ca0152600081816101230152818161039201526103d70152600081816102e801526103500152610a176000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806351dc86a5116100815780639579291e1161005b5780639579291e14610252578063a168fa891461025b578063f7818645146102c757600080fd5b806351dc86a5146101d057806367f082b0146101f75780636e04ff0d1461023157600080fd5b806334854043116100b257806334854043146101685780633b2bcbf1146101715780634585e33b146101bd57600080fd5b8063030932bb146100d9578063035262101461011e5780631fe543e314610153575b600080fd5b6101007f000000000000000000000000000000000000000000000000000000000000000081565b60405167ffffffffffffffff90911681526020015b60405180910390f35b6101457f000000000000000000000000000000000000000000000000000000000000000081565b604051908152602001610115565b6101666101613660046107c1565b6102d0565b005b61014560005481565b6101987f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610115565b6101666101cb36600461071d565b610390565b6101457f000000000000000000000000000000000000000000000000000000000000000081565b61021e7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610115565b61024461023f36600461071d565b6103d1565b6040516101159291906108b0565b61014560025481565b61029b61026936600461078f565b600360205260009081526040902080546001820154600290920154909160ff81169161010090910463ffffffff169084565b6040516101159493929190938452911515602084015263ffffffff166040830152606082015260800190565b61014560015481565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610382576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61038c828261040e565b5050565b7f0000000000000000000000000000000000000000000000000000000000000000600054426103bf919061092d565b111561038c574260005561038c61050f565b600060607f000000000000000000000000000000000000000000000000000000000000000060005442610404919061092d565b1191509250929050565b60008281526003602090815260409182902082516080810184528154808252600183015460ff811615159483019490945261010090930463ffffffff169381019390935260020154606083015283146104c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f72657175657374204944206e6f7420666f756e6420696e206d617000000000006044820152606401610379565b816000815181106104d6576104d66109ac565b602090810291909101810151600085815260039092526040822060029081019190915580549161050583610944565b9190505550505050565b6040517f5d3b1d300000000000000000000000000000000000000000000000000000000081527f0000000000000000000000000000000000000000000000000000000000000000600482015267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015261ffff7f0000000000000000000000000000000000000000000000000000000000000000166044820152620249f06064820152600160848201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561062357600080fd5b505af1158015610637573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061065b91906107a8565b6040805160808101825282815260006020808301828152620249f0848601908152606085018481528785526003909352948320935184555160018481018054965163ffffffff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff931515939093167fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000909716969096179190911790945551600290920191909155815492935061071583610944565b919050555050565b6000806020838503121561073057600080fd5b823567ffffffffffffffff8082111561074857600080fd5b818501915085601f83011261075c57600080fd5b81358181111561076b57600080fd5b86602082850101111561077d57600080fd5b60209290920196919550909350505050565b6000602082840312156107a157600080fd5b5035919050565b6000602082840312156107ba57600080fd5b5051919050565b600080604083850312156107d457600080fd5b8235915060208084013567ffffffffffffffff808211156107f457600080fd5b818601915086601f83011261080857600080fd5b81358181111561081a5761081a6109db565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561085d5761085d6109db565b604052828152858101935084860182860187018b101561087c57600080fd5b600095505b8386101561089f578035855260019590950194938601938601610881565b508096505050505050509250929050565b821515815260006020604081840152835180604085015260005b818110156108e6578581018301518582016060015282016108ca565b818111156108f8576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b60008282101561093f5761093f61097d565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156109765761097661097d565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeepersVRFConsumerABI = KeepersVRFConsumerMetaData.ABI

var KeepersVRFConsumerBin = KeepersVRFConsumerMetaData.Bin

func DeployKeepersVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, subscriptionId uint64, keyHash [32]byte, requestConfirmations uint16, upkeepInterval *big.Int) (common.Address, *types.Transaction, *KeepersVRFConsumer, error) {
	parsed, err := KeepersVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeepersVRFConsumerBin), backend, vrfCoordinator, subscriptionId, keyHash, requestConfirmations, upkeepInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeepersVRFConsumer{address: address, abi: *parsed, KeepersVRFConsumerCaller: KeepersVRFConsumerCaller{contract: contract}, KeepersVRFConsumerTransactor: KeepersVRFConsumerTransactor{contract: contract}, KeepersVRFConsumerFilterer: KeepersVRFConsumerFilterer{contract: contract}}, nil
}

type KeepersVRFConsumer struct {
	address common.Address
	abi     abi.ABI
	KeepersVRFConsumerCaller
	KeepersVRFConsumerTransactor
	KeepersVRFConsumerFilterer
}

type KeepersVRFConsumerCaller struct {
	contract *bind.BoundContract
}

type KeepersVRFConsumerTransactor struct {
	contract *bind.BoundContract
}

type KeepersVRFConsumerFilterer struct {
	contract *bind.BoundContract
}

type KeepersVRFConsumerSession struct {
	Contract     *KeepersVRFConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeepersVRFConsumerCallerSession struct {
	Contract *KeepersVRFConsumerCaller
	CallOpts bind.CallOpts
}

type KeepersVRFConsumerTransactorSession struct {
	Contract     *KeepersVRFConsumerTransactor
	TransactOpts bind.TransactOpts
}

type KeepersVRFConsumerRaw struct {
	Contract *KeepersVRFConsumer
}

type KeepersVRFConsumerCallerRaw struct {
	Contract *KeepersVRFConsumerCaller
}

type KeepersVRFConsumerTransactorRaw struct {
	Contract *KeepersVRFConsumerTransactor
}

func NewKeepersVRFConsumer(address common.Address, backend bind.ContractBackend) (*KeepersVRFConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(KeepersVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeepersVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeepersVRFConsumer{address: address, abi: abi, KeepersVRFConsumerCaller: KeepersVRFConsumerCaller{contract: contract}, KeepersVRFConsumerTransactor: KeepersVRFConsumerTransactor{contract: contract}, KeepersVRFConsumerFilterer: KeepersVRFConsumerFilterer{contract: contract}}, nil
}

func NewKeepersVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*KeepersVRFConsumerCaller, error) {
	contract, err := bindKeepersVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeepersVRFConsumerCaller{contract: contract}, nil
}

func NewKeepersVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeepersVRFConsumerTransactor, error) {
	contract, err := bindKeepersVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeepersVRFConsumerTransactor{contract: contract}, nil
}

func NewKeepersVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeepersVRFConsumerFilterer, error) {
	contract, err := bindKeepersVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeepersVRFConsumerFilterer{contract: contract}, nil
}

func bindKeepersVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeepersVRFConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeepersVRFConsumer *KeepersVRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeepersVRFConsumer.Contract.KeepersVRFConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.KeepersVRFConsumerTransactor.contract.Transfer(opts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.KeepersVRFConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeepersVRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.contract.Transfer(opts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) COORDINATOR() (common.Address, error) {
	return _KeepersVRFConsumer.Contract.COORDINATOR(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) COORDINATOR() (common.Address, error) {
	return _KeepersVRFConsumer.Contract.COORDINATOR(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) KEYHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "KEY_HASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) KEYHASH() ([32]byte, error) {
	return _KeepersVRFConsumer.Contract.KEYHASH(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) KEYHASH() ([32]byte, error) {
	return _KeepersVRFConsumer.Contract.KEYHASH(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) REQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) REQUESTCONFIRMATIONS() (uint16, error) {
	return _KeepersVRFConsumer.Contract.REQUESTCONFIRMATIONS(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) REQUESTCONFIRMATIONS() (uint16, error) {
	return _KeepersVRFConsumer.Contract.REQUESTCONFIRMATIONS(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) SUBSCRIPTIONID(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "SUBSCRIPTION_ID")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) SUBSCRIPTIONID() (uint64, error) {
	return _KeepersVRFConsumer.Contract.SUBSCRIPTIONID(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) SUBSCRIPTIONID() (uint64, error) {
	return _KeepersVRFConsumer.Contract.SUBSCRIPTIONID(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) UPKEEPINTERVAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "UPKEEP_INTERVAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) UPKEEPINTERVAL() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.UPKEEPINTERVAL(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) UPKEEPINTERVAL() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.UPKEEPINTERVAL(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) CheckUpkeep(opts *bind.CallOpts, arg0 []byte) (bool, []byte, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "checkUpkeep", arg0)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) CheckUpkeep(arg0 []byte) (bool, []byte, error) {
	return _KeepersVRFConsumer.Contract.CheckUpkeep(&_KeepersVRFConsumer.CallOpts, arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) CheckUpkeep(arg0 []byte) (bool, []byte, error) {
	return _KeepersVRFConsumer.Contract.CheckUpkeep(&_KeepersVRFConsumer.CallOpts, arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) SLastTimeStamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "s_lastTimeStamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) SLastTimeStamp() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.SLastTimeStamp(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) SLastTimeStamp() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.SLastTimeStamp(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.Randomness = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _KeepersVRFConsumer.Contract.SRequests(&_KeepersVRFConsumer.CallOpts, arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _KeepersVRFConsumer.Contract.SRequests(&_KeepersVRFConsumer.CallOpts, arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) SVrfRequestCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "s_vrfRequestCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) SVrfRequestCounter() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.SVrfRequestCounter(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) SVrfRequestCounter() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.SVrfRequestCounter(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCaller) SVrfResponseCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeepersVRFConsumer.contract.Call(opts, &out, "s_vrfResponseCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) SVrfResponseCounter() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.SVrfResponseCounter(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerCallerSession) SVrfResponseCounter() (*big.Int, error) {
	return _KeepersVRFConsumer.Contract.SVrfResponseCounter(&_KeepersVRFConsumer.CallOpts)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerTransactor) PerformUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error) {
	return _KeepersVRFConsumer.contract.Transact(opts, "performUpkeep", arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.PerformUpkeep(&_KeepersVRFConsumer.TransactOpts, arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerTransactorSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.PerformUpkeep(&_KeepersVRFConsumer.TransactOpts, arg0)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _KeepersVRFConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.RawFulfillRandomWords(&_KeepersVRFConsumer.TransactOpts, requestId, randomWords)
}

func (_KeepersVRFConsumer *KeepersVRFConsumerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _KeepersVRFConsumer.Contract.RawFulfillRandomWords(&_KeepersVRFConsumer.TransactOpts, requestId, randomWords)
}

type SRequests struct {
	RequestId        *big.Int
	Fulfilled        bool
	CallbackGasLimit uint32
	Randomness       *big.Int
}

func (_KeepersVRFConsumer *KeepersVRFConsumer) Address() common.Address {
	return _KeepersVRFConsumer.address
}

type KeepersVRFConsumerInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	KEYHASH(opts *bind.CallOpts) ([32]byte, error)

	REQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	SUBSCRIPTIONID(opts *bind.CallOpts) (uint64, error)

	UPKEEPINTERVAL(opts *bind.CallOpts) (*big.Int, error)

	CheckUpkeep(opts *bind.CallOpts, arg0 []byte) (bool, []byte, error)

	SLastTimeStamp(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SVrfRequestCounter(opts *bind.CallOpts) (*big.Int, error)

	SVrfResponseCounter(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	Address() common.Address
}
