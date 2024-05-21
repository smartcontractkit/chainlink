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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCreatedFundedAndConsumerAdded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINKTOKEN\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"_subTopUpAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"name\":\"requestRandomWordsWithForceFulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c0604052600060035560006004556103e760055534801561002057600080fd5b5060405161133038038061133083398101604081905261003f91610059565b60601b6001600160601b031916608081905260a052610089565b60006020828403121561006b57600080fd5b81516001600160a01b038116811461008257600080fd5b9392505050565b60805160601c60a05160601c61124b6100e560003960008181610156015281816103e00152818161048c0152818161056b0152818161067d0152818161072a0152610a250152600081816102c4015261032c015261124b6000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806362cce24411610097578063b1e2174911610066578063b1e2174914610256578063d826f88f1461025f578063d8a4676f1461027e578063dc1670db146102a357600080fd5b806362cce244146101c6578063737144bc146101d957806374dba124146101e2578063a168fa89146101eb57600080fd5b80632be555da116100d35780632be555da1461013e5780633b2bcbf11461015157806355380dfb1461019d578063557d2e92146101bd57600080fd5b80631757f11c146100fa5780631fe543e314610116578063271095ef1461012b575b600080fd5b61010360045481565b6040519081526020015b60405180910390f35b610129610124366004610e08565b6102ac565b005b610129610139366004610f14565b61036b565b61012961014c366004610f83565b610381565b6101787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010d565b6000546101789073ffffffffffffffffffffffffffffffffffffffff1681565b61010360025481565b6101296101d4366004610d5e565b610488565b61010360035481565b61010360055481565b61022c6101f9366004610dd6565b6008602052600090815260409020805460028201546003830154600484015460059094015460ff90931693919290919085565b6040805195151586526020860194909452928401919091526060830152608082015260a00161010d565b61010360065481565b6101296000600381905560048190556103e76005556002819055600155565b61029161028c366004610dd6565b6107a7565b60405161010d96959493929190611059565b61010360015481565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461035d576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b610367828261088c565b5050565b6103798686868686866109b2565b505050505050565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040805167ffffffffffffffff86166020820152634000aea0917f0000000000000000000000000000000000000000000000000000000000000000918691016040516020818303038152906040526040518463ffffffff1660e01b815260040161043093929190610fc1565b602060405180830381600087803b15801561044a57600080fd5b505af115801561045e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104829190610d35565b50505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156104f257600080fd5b505af1158015610506573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061052a9190610ef7565b6040517f7341c10c00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201523060248201529091507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b1580156105c457600080fd5b505af11580156105d8573d6000803e3d6000fd5b505050506105e7818484610381565b6040805167ffffffffffffffff831681523060208201529081018490527f56c142509574e8340ca0190b029c74464b84037d2876278ea0ade3ffb1f0042c9060600160405180910390a161063f8189898989896109b2565b6040517f9f87fad700000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201523060248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639f87fad790604401600060405180830381600087803b1580156106d657600080fd5b505af11580156106ea573d6000803e3d6000fd5b50506040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201523360248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16925063d7ae1d309150604401600060405180830381600087803b15801561078557600080fd5b505af1158015610799573d6000803e3d6000fd5b505050505050505050505050565b6000818152600860209081526040808320815160c081018352815460ff161515815260018201805484518187028101870190955280855260609587958695869586958695919492938584019390929083018282801561082557602002820191906000526020600020905b815481526020019060010190808311610811575b505050505081526020016002820154815260200160038201548152602001600482015481526020016005820154815250509050806000015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b6000610896610bc2565b600084815260076020526040812054919250906108b39083611155565b905060006108c482620f4240611118565b90506004548211156108d65760048290555b60055482106108e7576005546108e9565b815b6005556001546108f9578061092b565b60018054610906916110c5565b816001546003546109179190611118565b61092191906110c5565b61092b91906110dd565b600355600085815260086020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660019081178255865161097d939290910191870190610c86565b50600085815260086020526040812042600382015560050184905560018054916109a68361118e565b91905055505050505050565b60005b8161ffff168161ffff161015610bb9576040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810186905267ffffffffffffffff8816602482015261ffff8716604482015263ffffffff8086166064830152841660848201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b158015610a7e57600080fd5b505af1158015610a92573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ab69190610def565b600681905590506000610ac7610bc2565b6040805160c08101825260008082528251818152602080820185528084019182524284860152606084018390526080840186905260a084018390528783526008815293909120825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169015151781559051805194955091939092610b55926001850192910190610c86565b506040820151600280830191909155606083015160038301556080830151600483015560a0909201516005909101558054906000610b928361118e565b90915550506000918252600760205260409091205580610bb18161116c565b9150506109b5565b50505050505050565b600046610bce81610c5f565b15610c5857606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610c1a57600080fd5b505afa158015610c2e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c529190610def565b91505090565b4391505090565b600061a4b1821480610c73575062066eed82145b80610c80575062066eee82145b92915050565b828054828255906000526020600020908101928215610cc1579160200282015b82811115610cc1578251825591602001919060010190610ca6565b50610ccd929150610cd1565b5090565b5b80821115610ccd5760008155600101610cd2565b803573ffffffffffffffffffffffffffffffffffffffff81168114610d0a57600080fd5b919050565b803561ffff81168114610d0a57600080fd5b803563ffffffff81168114610d0a57600080fd5b600060208284031215610d4757600080fd5b81518015158114610d5757600080fd5b9392505050565b600080600080600080600060e0888a031215610d7957600080fd5b610d8288610d0f565b965060208801359550610d9760408901610d21565b9450610da560608901610d21565b9350610db360808901610d0f565b925060a08801359150610dc860c08901610ce6565b905092959891949750929550565b600060208284031215610de857600080fd5b5035919050565b600060208284031215610e0157600080fd5b5051919050565b60008060408385031215610e1b57600080fd5b8235915060208084013567ffffffffffffffff80821115610e3b57600080fd5b818601915086601f830112610e4f57600080fd5b813581811115610e6157610e616111f6565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610ea457610ea46111f6565b604052828152858101935084860182860187018b1015610ec357600080fd5b600095505b83861015610ee6578035855260019590950194938601938601610ec8565b508096505050505050509250929050565b600060208284031215610f0957600080fd5b8151610d5781611225565b60008060008060008060c08789031215610f2d57600080fd5b8635610f3881611225565b9550610f4660208801610d0f565b945060408701359350610f5b60608801610d21565b9250610f6960808801610d21565b9150610f7760a08801610d0f565b90509295509295509295565b600080600060608486031215610f9857600080fd5b8335610fa381611225565b925060208401359150610fb860408501610ce6565b90509250925092565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b8181101561101157858101830151858201608001528201610ff5565b81811115611023576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b600060c082018815158352602060c08185015281895180845260e086019150828b01935060005b8181101561109c57845183529383019391830191600101611080565b505060408501989098525050506060810193909352608083019190915260a09091015292915050565b600082198211156110d8576110d86111c7565b500190565b600082611113577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611150576111506111c7565b500290565b600082821015611167576111676111c7565b500390565b600061ffff80831681811415611184576111846111c7565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156111c0576111c06111c7565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff8116811461123b57600080fd5b5056fea164736f6c6343000806000a",
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
	return address, tx, &VRFV2LoadTestWithMetrics{address: address, abi: *parsed, VRFV2LoadTestWithMetricsCaller: VRFV2LoadTestWithMetricsCaller{contract: contract}, VRFV2LoadTestWithMetricsTransactor: VRFV2LoadTestWithMetricsTransactor{contract: contract}, VRFV2LoadTestWithMetricsFilterer: VRFV2LoadTestWithMetricsFilterer{contract: contract}}, nil
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) LINKTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "LINKTOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) LINKTOKEN() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.LINKTOKEN(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) LINKTOKEN() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.LINKTOKEN(&_VRFV2LoadTestWithMetrics.CallOpts)
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) RequestRandomWordsWithForceFulfill(opts *bind.TransactOpts, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int, _link common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "requestRandomWordsWithForceFulfill", _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount, _subTopUpAmount, _link)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) RequestRandomWordsWithForceFulfill(_requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int, _link common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RequestRandomWordsWithForceFulfill(&_VRFV2LoadTestWithMetrics.TransactOpts, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount, _subTopUpAmount, _link)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) RequestRandomWordsWithForceFulfill(_requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int, _link common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RequestRandomWordsWithForceFulfill(&_VRFV2LoadTestWithMetrics.TransactOpts, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount, _subTopUpAmount, _link)
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) TopUpSubscription(opts *bind.TransactOpts, _subId uint64, _amount *big.Int, _link common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "topUpSubscription", _subId, _amount, _link)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) TopUpSubscription(_subId uint64, _amount *big.Int, _link common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.TopUpSubscription(&_VRFV2LoadTestWithMetrics.TransactOpts, _subId, _amount, _link)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) TopUpSubscription(_subId uint64, _amount *big.Int, _link common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.TopUpSubscription(&_VRFV2LoadTestWithMetrics.TransactOpts, _subId, _amount, _link)
}

type VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator struct {
	Event *VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded)
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
		it.Event = new(VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded)
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

func (it *VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) FilterSubscriptionCreatedFundedAndConsumerAdded(opts *bind.FilterOpts) (*VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator, error) {

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.FilterLogs(opts, "SubscriptionCreatedFundedAndConsumerAdded")
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator{contract: _VRFV2LoadTestWithMetrics.contract, event: "SubscriptionCreatedFundedAndConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) WatchSubscriptionCreatedFundedAndConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded) (event.Subscription, error) {

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.WatchLogs(opts, "SubscriptionCreatedFundedAndConsumerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded)
				if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "SubscriptionCreatedFundedAndConsumerAdded", log); err != nil {
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) ParseSubscriptionCreatedFundedAndConsumerAdded(log types.Log) (*VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded, error) {
	event := new(VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded)
	if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "SubscriptionCreatedFundedAndConsumerAdded", log); err != nil {
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
	case _VRFV2LoadTestWithMetrics.abi.Events["SubscriptionCreatedFundedAndConsumerAdded"].ID:
		return _VRFV2LoadTestWithMetrics.ParseSubscriptionCreatedFundedAndConsumerAdded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x56c142509574e8340ca0190b029c74464b84037d2876278ea0ade3ffb1f0042c")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetrics) Address() common.Address {
	return _VRFV2LoadTestWithMetrics.address
}

type VRFV2LoadTestWithMetricsInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINKTOKEN(opts *bind.CallOpts) (common.Address, error)

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SLastRequestId(opts *bind.CallOpts) (*big.Int, error)

	SRequestCount(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error)

	RequestRandomWordsWithForceFulfill(opts *bind.TransactOpts, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int, _link common.Address) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, _subId uint64, _amount *big.Int, _link common.Address) (*types.Transaction, error)

	FilterSubscriptionCreatedFundedAndConsumerAdded(opts *bind.FilterOpts) (*VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAddedIterator, error)

	WatchSubscriptionCreatedFundedAndConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded) (event.Subscription, error)

	ParseSubscriptionCreatedFundedAndConsumerAdded(log types.Log) (*VRFV2LoadTestWithMetricsSubscriptionCreatedFundedAndConsumerAdded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
