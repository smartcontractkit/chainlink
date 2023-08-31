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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"_nativePayment\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052600060055560006006556103e760075534801561002057600080fd5b5060405161139238038061139283398101604081905261003f916101a9565b8033806000816100965760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100c6576100c6816100ff565b5050600280546001600160a01b0319166001600160a01b0393909316929092179091555060601b6001600160601b0319166080526101d9565b6001600160a01b0381163314156101585760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161008d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156101bb57600080fd5b81516001600160a01b03811681146101d257600080fd5b9392505050565b60805160601c6111946101fe6000396000818161013b01526103f701526111946000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80638da5cb5b11610097578063d826f88f11610066578063d826f88f1461025d578063d8a4676f1461027c578063dc1670db146102a1578063f2fde38b146102aa57600080fd5b80638da5cb5b146101b85780638ea98117146101d6578063a168fa89146101e9578063b1e217491461025457600080fd5b80636846de20116100d35780636846de201461018b578063737144bc1461019e57806374dba124146101a757806379ba5097146101b057600080fd5b80631757f11c146101055780631fe543e3146101215780633b2bcbf114610136578063557d2e9214610182575b600080fd5b61010e60065481565b6040519081526020015b60405180910390f35b61013461012f366004610d8d565b6102bd565b005b61015d7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610118565b61010e60045481565b610134610199366004610e7c565b610343565b61010e60055481565b61010e60075481565b61013461058c565b60005473ffffffffffffffffffffffffffffffffffffffff1661015d565b6101346101e4366004610d1e565b610689565b61022a6101f7366004610d5b565b600a602052600090815260409020805460028201546003830154600484015460059094015460ff90931693919290919085565b6040805195151586526020860194909452928401919091526060830152608082015260a001610118565b61010e60085481565b6101346000600581905560068190556103e76007556004819055600355565b61028f61028a366004610d5b565b610794565b60405161011896959493929190610efb565b61010e60035481565b6101346102b8366004610d1e565b610879565b60025473ffffffffffffffffffffffffffffffffffffffff163314610335576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61033f828261088d565b5050565b61034b6109b8565b60005b8161ffff168161ffff1610156105825760006040518060c001604052808881526020018a81526020018961ffff1681526020018763ffffffff1681526020018563ffffffff1681526020016103b26040518060200160405280891515815250610a3b565b90526040517f9b1c385e00000000000000000000000000000000000000000000000000000000815290915060009073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690639b1c385e9061042c908590600401610f67565b602060405180830381600087803b15801561044657600080fd5b505af115801561045a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061047e9190610d74565b60088190559050600061048f610af7565b6040805160c08101825260008082528251818152602080820185528084019182524284860152606084018390526080840186905260a08401839052878352600a815293909120825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016901515178155905180519495509193909261051d926001850192910190610c93565b506040820151600282015560608201516003820155608082015160048083019190915560a090920151600590910155805490600061055a836110f0565b909155505060009182526009602052604090912055508061057a816110ce565b91505061034e565b5050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461060d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161032c565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906106c9575060025473ffffffffffffffffffffffffffffffffffffffff163314155b1561074d57336106ee60005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9384166004820152918316602483015291909116604482015260640161032c565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000818152600a60209081526040808320815160c081018352815460ff161515815260018201805484518187028101870190955280855260609587958695869586958695919492938584019390929083018282801561081257602002820191906000526020600020905b8154815260200190600101908083116107fe575b505050505081526020016002820154815260200160038201548152602001600482015481526020016005820154815250509050806000015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b6108816109b8565b61088a81610b9d565b50565b6000610897610af7565b600084815260096020526040812054919250906108b490836110b7565b905060006108c582620f424061107a565b90506006548211156108d75760068290555b60075482106108e8576007546108ea565b815b6007556003546108fa578061092d565b600354610908906001611027565b81600354600554610919919061107a565b6109239190611027565b61092d919061103f565b6005556000858152600a6020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660019081178255865161097f939290910191870190610c93565b506000858152600a602052604081204260038083019190915560059091018590558054916109ac836110f0565b91905055505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a39576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161032c565b565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610a7491511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b60004661a4b1811480610b0c575062066eed81145b15610b9657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610b5857600080fd5b505afa158015610b6c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b909190610d74565b91505090565b4391505090565b73ffffffffffffffffffffffffffffffffffffffff8116331415610c1d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161032c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610cce579160200282015b82811115610cce578251825591602001919060010190610cb3565b50610cda929150610cde565b5090565b5b80821115610cda5760008155600101610cdf565b803561ffff81168114610d0557600080fd5b919050565b803563ffffffff81168114610d0557600080fd5b600060208284031215610d3057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610d5457600080fd5b9392505050565b600060208284031215610d6d57600080fd5b5035919050565b600060208284031215610d8657600080fd5b5051919050565b60008060408385031215610da057600080fd5b8235915060208084013567ffffffffffffffff80821115610dc057600080fd5b818601915086601f830112610dd457600080fd5b813581811115610de657610de6611158565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610e2957610e29611158565b604052828152858101935084860182860187018b1015610e4857600080fd5b600095505b83861015610e6b578035855260019590950194938601938601610e4d565b508096505050505050509250929050565b600080600080600080600060e0888a031215610e9757600080fd5b87359650610ea760208901610cf3565b955060408801359450610ebc60608901610d0a565b935060808801358015158114610ed157600080fd5b9250610edf60a08901610d0a565b9150610eed60c08901610cf3565b905092959891949750929550565b600060c082018815158352602060c08185015281895180845260e086019150828b01935060005b81811015610f3e57845183529383019391830191600101610f22565b505060408501989098525050506060810193909352608083019190915260a09091015292915050565b6000602080835283518184015280840151604084015261ffff6040850151166060840152606084015163ffffffff80821660808601528060808701511660a0860152505060a084015160c08085015280518060e086015260005b81811015610fde5782810184015186820161010001528301610fc1565b81811115610ff157600061010083880101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169390930161010001949350505050565b6000821982111561103a5761103a611129565b500190565b600082611075577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156110b2576110b2611129565b500290565b6000828210156110c9576110c9611129565b500390565b600061ffff808316818114156110e6576110e6611129565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561112257611122611129565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	return address, tx, &VRFV2PlusLoadTestWithMetrics{VRFV2PlusLoadTestWithMetricsCaller: VRFV2PlusLoadTestWithMetricsCaller{contract: contract}, VRFV2PlusLoadTestWithMetricsTransactor: VRFV2PlusLoadTestWithMetricsTransactor{contract: contract}, VRFV2PlusLoadTestWithMetricsFilterer: VRFV2PlusLoadTestWithMetricsFilterer{contract: contract}}, nil
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

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusLoadTestWithMetrics.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsSession) COORDINATOR() (common.Address, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.COORDINATOR(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
}

func (_VRFV2PlusLoadTestWithMetrics *VRFV2PlusLoadTestWithMetricsCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2PlusLoadTestWithMetrics.Contract.COORDINATOR(&_VRFV2PlusLoadTestWithMetrics.CallOpts)
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
