// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_load_test_with_metrics

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

var VRFV2PlusLoadTestWithMetricsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"_nativePayment\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060055560006006556103e760075534801561002057600080fd5b5060405161136538038061136583398101604081905261003f9161019b565b8033806000816100965760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100c6576100c6816100f1565b5050600280546001600160a01b0319166001600160a01b039390931692909217909155506101cb9050565b6001600160a01b03811633141561014a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161008d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156101ad57600080fd5b81516001600160a01b03811681146101c457600080fd5b9392505050565b61118b806101da6000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80638ea9811711610097578063d826f88f11610066578063d826f88f14610252578063d8a4676f14610271578063dc1670db14610296578063f2fde38b1461029f57600080fd5b80638ea98117146101ab5780639eccacf6146101be578063a168fa89146101de578063b1e217491461024957600080fd5b8063737144bc116100d3578063737144bc1461015257806374dba1241461015b57806379ba5097146101645780638da5cb5b1461016c57600080fd5b80631757f11c146101055780631fe543e314610121578063557d2e92146101365780636846de201461013f575b600080fd5b61010e60065481565b6040519081526020015b60405180910390f35b61013461012f366004610d84565b6102b2565b005b61010e60045481565b61013461014d366004610e73565b610338565b61010e60055481565b61010e60075481565b610134610565565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610118565b6101346101b9366004610d15565b610662565b6002546101869073ffffffffffffffffffffffffffffffffffffffff1681565b61021f6101ec366004610d52565b600a602052600090815260409020805460028201546003830154600484015460059094015460ff90931693919290919085565b6040805195151586526020860194909452928401919091526060830152608082015260a001610118565b61010e60085481565b6101346000600581905560068190556103e76007556004819055600355565b61028461027f366004610d52565b61076d565b60405161011896959493929190610ef2565b61010e60035481565b6101346102ad366004610d15565b610852565b60025473ffffffffffffffffffffffffffffffffffffffff16331461032a576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6103348282610866565b5050565b610340610991565b60005b8161ffff168161ffff16101561055b5760006040518060c001604052808881526020018a81526020018961ffff1681526020018763ffffffff1681526020018563ffffffff1681526020016103a76040518060200160405280891515815250610a14565b90526002546040517f9b1c385e00000000000000000000000000000000000000000000000000000000815291925060009173ffffffffffffffffffffffffffffffffffffffff90911690639b1c385e90610405908590600401610f5e565b602060405180830381600087803b15801561041f57600080fd5b505af1158015610433573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104579190610d6b565b600881905590506000610468610ad0565b6040805160c08101825260008082528251818152602080820185528084019182524284860152606084018390526080840186905260a08401839052878352600a815293909120825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690151517815590518051949550919390926104f6926001850192910190610c8a565b506040820151600282015560608201516003820155608082015160048083019190915560a0909201516005909101558054906000610533836110e7565b9091555050600091825260096020526040909120555080610553816110c5565b915050610343565b5050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610321565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906106a2575060025473ffffffffffffffffffffffffffffffffffffffff163314155b1561072657336106c760005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff93841660048201529183166024830152919091166044820152606401610321565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000818152600a60209081526040808320815160c081018352815460ff16151581526001820180548451818702810187019095528085526060958795869586958695869591949293858401939092908301828280156107eb57602002820191906000526020600020905b8154815260200190600101908083116107d7575b505050505081526020016002820154815260200160038201548152602001600482015481526020016005820154815250509050806000015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b61085a610991565b61086381610b6d565b50565b6000610870610ad0565b6000848152600960205260408120549192509061088d90836110ae565b9050600061089e82620f4240611071565b90506006548211156108b05760068290555b60075482106108c1576007546108c3565b815b6007556003546108d35780610906565b6003546108e190600161101e565b816003546005546108f29190611071565b6108fc919061101e565b6109069190611036565b6005556000858152600a6020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600190811782558651610958939290910191870190610c8a565b506000858152600a60205260408120426003808301919091556005909101859055805491610985836110e7565b91905055505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a12576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610321565b565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610a4d91511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b600046610adc81610c63565b15610b6657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610b2857600080fd5b505afa158015610b3c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b609190610d6b565b91505090565b4391505090565b73ffffffffffffffffffffffffffffffffffffffff8116331415610bed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610321565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061a4b1821480610c77575062066eed82145b80610c84575062066eee82145b92915050565b828054828255906000526020600020908101928215610cc5579160200282015b82811115610cc5578251825591602001919060010190610caa565b50610cd1929150610cd5565b5090565b5b80821115610cd15760008155600101610cd6565b803561ffff81168114610cfc57600080fd5b919050565b803563ffffffff81168114610cfc57600080fd5b600060208284031215610d2757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610d4b57600080fd5b9392505050565b600060208284031215610d6457600080fd5b5035919050565b600060208284031215610d7d57600080fd5b5051919050565b60008060408385031215610d9757600080fd5b8235915060208084013567ffffffffffffffff80821115610db757600080fd5b818601915086601f830112610dcb57600080fd5b813581811115610ddd57610ddd61114f565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610e2057610e2061114f565b604052828152858101935084860182860187018b1015610e3f57600080fd5b600095505b83861015610e62578035855260019590950194938601938601610e44565b508096505050505050509250929050565b600080600080600080600060e0888a031215610e8e57600080fd5b87359650610e9e60208901610cea565b955060408801359450610eb360608901610d01565b935060808801358015158114610ec857600080fd5b9250610ed660a08901610d01565b9150610ee460c08901610cea565b905092959891949750929550565b600060c082018815158352602060c08185015281895180845260e086019150828b01935060005b81811015610f3557845183529383019391830191600101610f19565b505060408501989098525050506060810193909352608083019190915260a09091015292915050565b6000602080835283518184015280840151604084015261ffff6040850151166060840152606084015163ffffffff80821660808601528060808701511660a0860152505060a084015160c08085015280518060e086015260005b81811015610fd55782810184015186820161010001528301610fb8565b81811115610fe857600061010083880101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169390930161010001949350505050565b6000821982111561103157611031611120565b500190565b60008261106c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156110a9576110a9611120565b500290565b6000828210156110c0576110c0611120565b500390565b600061ffff808316818114156110dd576110dd611120565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561111957611119611120565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusLoadTestWithMetricsABI = VRFV2PlusLoadTestWithMetricsMetaData.ABI

var VRFV2PlusLoadTestWithMetricsBin = VRFV2PlusLoadTestWithMetricsMetaData.Bin

func DeployVRFV2PlusLoadTestWithMetrics(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFV2PlusLoadTestWithMetrics, error) {
	parsed, err := VRFV2PlusLoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusLoadTestWithMetricsBin), backend, _vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusLoadTestWithMetrics{address: address, abi: *parsed, VRFV2PlusLoadTestWithMetricsCaller: VRFV2PlusLoadTestWithMetricsCaller{contract: contract}, VRFV2PlusLoadTestWithMetricsTransactor: VRFV2PlusLoadTestWithMetricsTransactor{contract: contract}, VRFV2PlusLoadTestWithMetricsFilterer: VRFV2PlusLoadTestWithMetricsFilterer{contract: contract}}, nil
}

type VRFV2PlusLoadTestWithMetrics struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusLoadTestWithMetricsCaller
	VRFV2PlusLoadTestWithMetricsTransactor
	VRFV2PlusLoadTestWithMetricsFilterer
}

type VRFV2PlusLoadTestWithMetricsCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusLoadTestWithMetricsTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusLoadTestWithMetricsFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusLoadTestWithMetricsSession struct {
	Contract     *VRFV2PlusLoadTestWithMetrics
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusLoadTestWithMetricsCallerSession struct {
	Contract *VRFV2PlusLoadTestWithMetricsCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusLoadTestWithMetricsTransactorSession struct {
	Contract     *VRFV2PlusLoadTestWithMetricsTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusLoadTestWithMetricsRaw struct {
	Contract *VRFV2PlusLoadTestWithMetrics
}

type VRFV2PlusLoadTestWithMetricsCallerRaw struct {
	Contract *VRFV2PlusLoadTestWithMetricsCaller
}

type VRFV2PlusLoadTestWithMetricsTransactorRaw struct {
	Contract *VRFV2PlusLoadTestWithMetricsTransactor
}

func NewVRFV2PlusLoadTestWithMetrics(address common.Address, backend bind.ContractBackend) (*VRFV2PlusLoadTestWithMetrics, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusLoadTestWithMetricsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusLoadTestWithMetrics(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusLoadTestWithMetrics{address: address, abi: abi, VRFV2PlusLoadTestWithMetricsCaller: VRFV2PlusLoadTestWithMetricsCaller{contract: contract}, VRFV2PlusLoadTestWithMetricsTransactor: VRFV2PlusLoadTestWithMetricsTransactor{contract: contract}, VRFV2PlusLoadTestWithMetricsFilterer: VRFV2PlusLoadTestWithMetricsFilterer{contract: contract}}, nil
}

func NewVRFV2PlusLoadTestWithMetricsCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusLoadTestWithMetricsCaller, error) {
	contract, err := bindVRFV2PlusLoadTestWithMetrics(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusLoadTestWithMetricsCaller{contract: contract}, nil
}

func NewVRFV2PlusLoadTestWithMetricsTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusLoadTestWithMetricsTransactor, error) {
	contract, err := bindVRFV2PlusLoadTestWithMetrics(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusLoadTestWithMetricsTransactor{contract: contract}, nil
}

func NewVRFV2PlusLoadTestWithMetricsFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusLoadTestWithMetricsFilterer, error) {
	contract, err := bindVRFV2PlusLoadTestWithMetrics(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusLoadTestWithMetricsFilterer{contract: contract}, nil
}

func bindVRFV2PlusLoadTestWithMetrics(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusLoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusLoadTestWithMetrics.Contract.VRFV2PlusLoadTestWithMetricsCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.VRFV2PlusLoadTestWithMetricsTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.VRFV2PlusLoadTestWithMetricsTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusLoadTestWithMetrics.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "getRequestStatus", _requestId)

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

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.GetRequestStatus(&_VRFV2PlusLoadTestWithMetrics.CallOpts, _requestId)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.GetRequestStatus(&_VRFV2PlusLoadTestWithMetrics.CallOpts, _requestId)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) Owner() (common.Address, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.Owner(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.Owner(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SAverageFulfillmentInMillions(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SAverageFulfillmentInMillions(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SFastestFulfillment(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SFastestFulfillment(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SLastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SLastRequestId(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SLastRequestId(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SRequestCount() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SRequestCount(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SRequestCount(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_requests", arg0)

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

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SRequests(&_VRFV2PlusLoadTestWithMetrics.CallOpts, arg0)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SRequests(&_VRFV2PlusLoadTestWithMetrics.CallOpts, arg0)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SResponseCount() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SResponseCount(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SResponseCount(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SSlowestFulfillment(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SSlowestFulfillment(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SVrfCoordinator(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SVrfCoordinator(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.AcceptOwnership(&_VRFV2PlusLoadTestWithMetrics.TransactOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.AcceptOwnership(&_VRFV2PlusLoadTestWithMetrics.TransactOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.RawFulfillRandomWords(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.RawFulfillRandomWords(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactor) RequestRandomWords(opts *bind.TransactOpts, _subId *big.Int, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _nativePayment bool, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.contract.Transact(opts, "requestRandomWords", _subId, _requestConfirmations, _keyHash, _callbackGasLimit, _nativePayment, _numWords, _requestCount)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) RequestRandomWords(_subId *big.Int, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _nativePayment bool, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.RequestRandomWords(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, _subId, _requestConfirmations, _keyHash, _callbackGasLimit, _nativePayment, _numWords, _requestCount)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorSession) RequestRandomWords(_subId *big.Int, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _nativePayment bool, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.RequestRandomWords(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, _subId, _requestConfirmations, _keyHash, _callbackGasLimit, _nativePayment, _numWords, _requestCount)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.contract.Transact(opts, "reset")
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) Reset() (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.Reset(&_VRFV2PlusLoadTestWithMetrics.TransactOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorSession) Reset() (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.Reset(&_VRFV2PlusLoadTestWithMetrics.TransactOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactor) SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.contract.Transact(opts, "setCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SetCoordinator(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.SetCoordinator(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.TransferOwnership(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, to)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.TransferOwnership(&_VRFV2PlusLoadTestWithMetrics.TransactOpts, to)
}

type VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested)
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

func (it *VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusLoadTestWithMetrics.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator{contract: _VRFV2PlusLoadTestWithMetrics.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusLoadTestWithMetrics.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested)
				if err := _VRFV2PlusLoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested, error) {
	event := new(VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested)
	if err := _VRFV2PlusLoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator struct {
	Event *VRFV2PlusLoadTestWithMetricsOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusLoadTestWithMetricsOwnershipTransferred)
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
		it.Event = new(VRFV2PlusLoadTestWithMetricsOwnershipTransferred)
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

func (it *VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusLoadTestWithMetricsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusLoadTestWithMetrics.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator{contract: _VRFV2PlusLoadTestWithMetrics.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusLoadTestWithMetricsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusLoadTestWithMetrics.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusLoadTestWithMetricsOwnershipTransferred)
				if err := _VRFV2PlusLoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferred, error) {
	event := new(VRFV2PlusLoadTestWithMetricsOwnershipTransferred)
	if err := _VRFV2PlusLoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetrics) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusLoadTestWithMetrics.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusLoadTestWithMetrics.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusLoadTestWithMetrics.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusLoadTestWithMetrics.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusLoadTestWithMetricsOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetrics) Address() common.Address {
	return _VRFV2PlusLoadTestWithMetrics.address
}

type VRFV2PlusLoadTestWithMetricsInterface interface {
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

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _subId *big.Int, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _nativePayment bool, _numWords uint32, _requestCount uint16) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusLoadTestWithMetricsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusLoadTestWithMetricsOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
