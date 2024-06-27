// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_load_test_with_metrics

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

var VRFV2LoadTestWithMetricsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c0604052600060045560006005556103e760065534801561002057600080fd5b5060405161111138038061111183398101604081905261003f91610199565b6001600160601b0319606082901b1660805233806000816100a75760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100d7576100d7816100ef565b50505060601b6001600160601b03191660a0526101c9565b6001600160a01b0381163314156101485760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161009e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156101ab57600080fd5b81516001600160a01b03811681146101c257600080fd5b9392505050565b60805160601c60a05160601c610f0f6102026000396000818161014301526103da0152600081816102b7015261031f0152610f0f6000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806379ba509711610097578063d826f88f11610066578063d826f88f1461023f578063d8a4676f1461025e578063dc1670db14610283578063f2fde38b1461028c57600080fd5b806379ba5097146101a55780638da5cb5b146101ad578063a168fa89146101cb578063b1e217491461023657600080fd5b80633b2bcbf1116100d35780633b2bcbf11461013e578063557d2e921461018a578063737144bc1461019357806374dba1241461019c57600080fd5b80631757f11c146100fa5780631fe543e314610116578063271095ef1461012b575b600080fd5b61010360055481565b6040519081526020015b60405180910390f35b610129610124366004610bcb565b61029f565b005b610129610139366004610cba565b61035f565b6101657f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010d565b61010360035481565b61010360045481565b61010360065481565b610129610577565b60005473ffffffffffffffffffffffffffffffffffffffff16610165565b61020c6101d9366004610b99565b6009602052600090815260409020805460028201546003830154600484015460059094015460ff90931693919290919085565b6040805195151586526020860194909452928401919091526060830152608082015260a00161010d565b61010360075481565b6101296000600481905560058190556103e76006556003819055600255565b61027161026c366004610b99565b610674565b60405161010d96959493929190610d36565b61010360025481565b61012961029a366004610b5c565b610759565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610351576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61035b828261076d565b5050565b610367610894565b60005b8161ffff168161ffff16101561056e576040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810186905267ffffffffffffffff8816602482015261ffff8716604482015263ffffffff8086166064830152841660848201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561043357600080fd5b505af1158015610447573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061046b9190610bb2565b60078190559050600061047c610917565b6040805160c08101825260008082528251818152602080820185528084019182524284860152606084018390526080840186905260a084018390528783526009815293909120825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016901515178155905180519495509193909261050a926001850192910190610ad1565b506040820151600282015560608201516003808301919091556080830151600483015560a090920151600590910155805490600061054783610e6b565b9091555050600091825260086020526040909120558061056681610e49565b91505061036a565b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105f8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610348565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000818152600960209081526040808320815160c081018352815460ff16151581526001820180548451818702810187019095528085526060958795869586958695869591949293858401939092908301828280156106f257602002820191906000526020600020905b8154815260200190600101908083116106de575b505050505081526020016002820154815260200160038201548152602001600482015481526020016005820154815250509050806000015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b610761610894565b61076a816109b4565b50565b6000610777610917565b600084815260086020526040812054919250906107949083610e32565b905060006107a582620f4240610df5565b90506005548211156107b75760058290555b60065482106107c8576006546107ca565b815b6006556002546107da578061080d565b6002546107e8906001610da2565b816002546004546107f99190610df5565b6108039190610da2565b61080d9190610dba565b600455600085815260096020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660019081178255865161085f939290910191870190610ad1565b506000858152600960205260408120426003820155600501849055600280549161088883610e6b565b91905055505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610915576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610348565b565b60004661092381610aaa565b156109ad57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561096f57600080fd5b505afa158015610983573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109a79190610bb2565b91505090565b4391505090565b73ffffffffffffffffffffffffffffffffffffffff8116331415610a34576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610348565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061a4b1821480610abe575062066eed82145b80610acb575062066eee82145b92915050565b828054828255906000526020600020908101928215610b0c579160200282015b82811115610b0c578251825591602001919060010190610af1565b50610b18929150610b1c565b5090565b5b80821115610b185760008155600101610b1d565b803561ffff81168114610b4357600080fd5b919050565b803563ffffffff81168114610b4357600080fd5b600060208284031215610b6e57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610b9257600080fd5b9392505050565b600060208284031215610bab57600080fd5b5035919050565b600060208284031215610bc457600080fd5b5051919050565b60008060408385031215610bde57600080fd5b8235915060208084013567ffffffffffffffff80821115610bfe57600080fd5b818601915086601f830112610c1257600080fd5b813581811115610c2457610c24610ed3565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610c6757610c67610ed3565b604052828152858101935084860182860187018b1015610c8657600080fd5b600095505b83861015610ca9578035855260019590950194938601938601610c8b565b508096505050505050509250929050565b60008060008060008060c08789031215610cd357600080fd5b863567ffffffffffffffff81168114610ceb57600080fd5b9550610cf960208801610b31565b945060408701359350610d0e60608801610b48565b9250610d1c60808801610b48565b9150610d2a60a08801610b31565b90509295509295509295565b600060c082018815158352602060c08185015281895180845260e086019150828b01935060005b81811015610d7957845183529383019391830191600101610d5d565b505060408501989098525050506060810193909352608083019190915260a09091015292915050565b60008219821115610db557610db5610ea4565b500190565b600082610df0577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610e2d57610e2d610ea4565b500290565b600082821015610e4457610e44610ea4565b500390565b600061ffff80831681811415610e6157610e61610ea4565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610e9d57610e9d610ea4565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2LoadTestWithMetricsABI = VRFV2LoadTestWithMetricsMetaData.ABI

var VRFV2LoadTestWithMetricsBin = VRFV2LoadTestWithMetricsMetaData.Bin

func DeployVRFV2LoadTestWithMetrics(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFV2LoadTestWithMetrics, error) {
	parsed, err := VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2LoadTestWithMetricsBin), backend, _vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2LoadTestWithMetrics{VRFV2LoadTestWithMetricsCaller: VRFV2LoadTestWithMetricsCaller{contract: contract}, VRFV2LoadTestWithMetricsTransactor: VRFV2LoadTestWithMetricsTransactor{contract: contract}, VRFV2LoadTestWithMetricsFilterer: VRFV2LoadTestWithMetricsFilterer{contract: contract}}, nil
}

type VRFV2LoadTestWithMetrics struct {
	address common.Address
	abi     abi.ABI
	VRFV2LoadTestWithMetricsCaller
	VRFV2LoadTestWithMetricsTransactor
	VRFV2LoadTestWithMetricsFilterer
}

type VRFV2LoadTestWithMetricsCaller struct {
	contract *bind.BoundContract
}

type VRFV2LoadTestWithMetricsTransactor struct {
	contract *bind.BoundContract
}

type VRFV2LoadTestWithMetricsFilterer struct {
	contract *bind.BoundContract
}

type VRFV2LoadTestWithMetricsSession struct {
	Contract     *VRFV2LoadTestWithMetrics
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2LoadTestWithMetricsCallerSession struct {
	Contract *VRFV2LoadTestWithMetricsCaller
	CallOpts bind.CallOpts
}

type VRFV2LoadTestWithMetricsTransactorSession struct {
	Contract     *VRFV2LoadTestWithMetricsTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2LoadTestWithMetricsRaw struct {
	Contract *VRFV2LoadTestWithMetrics
}

type VRFV2LoadTestWithMetricsCallerRaw struct {
	Contract *VRFV2LoadTestWithMetricsCaller
}

type VRFV2LoadTestWithMetricsTransactorRaw struct {
	Contract *VRFV2LoadTestWithMetricsTransactor
}

func NewVRFV2LoadTestWithMetrics(address common.Address, backend bind.ContractBackend) (*VRFV2LoadTestWithMetrics, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2LoadTestWithMetricsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2LoadTestWithMetrics(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetrics{address: address, abi: abi, VRFV2LoadTestWithMetricsCaller: VRFV2LoadTestWithMetricsCaller{contract: contract}, VRFV2LoadTestWithMetricsTransactor: VRFV2LoadTestWithMetricsTransactor{contract: contract}, VRFV2LoadTestWithMetricsFilterer: VRFV2LoadTestWithMetricsFilterer{contract: contract}}, nil
}

func NewVRFV2LoadTestWithMetricsCaller(address common.Address, caller bind.ContractCaller) (*VRFV2LoadTestWithMetricsCaller, error) {
	contract, err := bindVRFV2LoadTestWithMetrics(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsCaller{contract: contract}, nil
}

func NewVRFV2LoadTestWithMetricsTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2LoadTestWithMetricsTransactor, error) {
	contract, err := bindVRFV2LoadTestWithMetrics(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsTransactor{contract: contract}, nil
}

func NewVRFV2LoadTestWithMetricsFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2LoadTestWithMetricsFilterer, error) {
	contract, err := bindVRFV2LoadTestWithMetrics(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsFilterer{contract: contract}, nil
}

func bindVRFV2LoadTestWithMetrics(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2LoadTestWithMetrics.Contract.VRFV2LoadTestWithMetricsCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.VRFV2LoadTestWithMetricsTransactor.contract.Transfer(opts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.VRFV2LoadTestWithMetricsTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2LoadTestWithMetrics.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.contract.Transfer(opts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) COORDINATOR() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.COORDINATOR(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.COORDINATOR(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.RequestTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2LoadTestWithMetrics.Contract.GetRequestStatus(&_VRFV2LoadTestWithMetrics.CallOpts, _requestId)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2LoadTestWithMetrics.Contract.GetRequestStatus(&_VRFV2LoadTestWithMetrics.CallOpts, _requestId)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) Owner() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Owner(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) Owner() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Owner(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SAverageFulfillmentInMillions(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SAverageFulfillmentInMillions(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SFastestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SFastestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SLastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SLastRequestId(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SLastRequestId(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SRequestCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SRequestCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SRequestCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RequestTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2LoadTestWithMetrics.Contract.SRequests(&_VRFV2LoadTestWithMetrics.CallOpts, arg0)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2LoadTestWithMetrics.Contract.SRequests(&_VRFV2LoadTestWithMetrics.CallOpts, arg0)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SResponseCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SResponseCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SResponseCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SSlowestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SSlowestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.AcceptOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.AcceptOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RawFulfillRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, requestId, randomWords)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RawFulfillRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, requestId, randomWords)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "requestRandomWords", _subId, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) RequestRandomWords(_subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RequestRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, _subId, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) RequestRandomWords(_subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RequestRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, _subId, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "reset")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) Reset() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Reset(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) Reset() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Reset(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.TransferOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts, to)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.TransferOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts, to)
}

type VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator struct {
	Event *VRFV2LoadTestWithMetricsOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
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
		it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
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

func (it *VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2LoadTestWithMetricsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator{contract: _VRFV2LoadTestWithMetrics.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
				if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferRequested, error) {
	event := new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
	if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2LoadTestWithMetricsOwnershipTransferredIterator struct {
	Event *VRFV2LoadTestWithMetricsOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferred)
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
		it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferred)
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

func (it *VRFV2LoadTestWithMetricsOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2LoadTestWithMetricsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsOwnershipTransferredIterator{contract: _VRFV2LoadTestWithMetrics.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2LoadTestWithMetricsOwnershipTransferred)
				if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferred, error) {
	event := new(VRFV2LoadTestWithMetricsOwnershipTransferred)
	if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Fulfilled             bool
	RandomWords           []*big.Int
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}
type SRequests struct {
	Fulfilled             bool
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetrics) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2LoadTestWithMetrics.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2LoadTestWithMetrics.ParseOwnershipTransferRequested(log)
	case _VRFV2LoadTestWithMetrics.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2LoadTestWithMetrics.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2LoadTestWithMetricsOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2LoadTestWithMetricsOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetrics) Address() common.Address {
	return _VRFV2LoadTestWithMetrics.address
}

type VRFV2LoadTestWithMetricsInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SLastRequestId(opts *bind.CallOpts) (*big.Int, error)

	SRequestCount(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
