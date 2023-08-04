// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_client_example

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

var FunctionsClientExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001a8d38038062001a8d833981016040819052620000349162000180565b6001600160a01b0381166080523380600081620000985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cb57620000cb81620000d5565b50505050620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b6080516118b8620001d5600039600081816101c60152610c4f01526118b86000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80636d9809a011610081578063b1e217491161005b578063b1e2174914610182578063f2fde38b1461018b578063f7b4c06f1461019e57600080fd5b80636d9809a01461014857806379ba5097146101525780638da5cb5b1461015a57600080fd5b806342748b2a116100b257806342748b2a146100ff5780634b0795a81461012c5780635fa353e71461013557600080fd5b80630ca76175146100ce5780633944ea3a146100e3575b600080fd5b6100e16100dc366004611304565b6101ae565b005b6100ec60035481565b6040519081526020015b60405180910390f35b60055461011790640100000000900463ffffffff1681565b60405163ffffffff90911681526020016100f6565b6100ec60045481565b6100e16101433660046113d7565b61022d565b6101176201117081565b6100e161033f565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b6100ec60025481565b6100e16101993660046114bb565b610441565b6005546101179063ffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461021d576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610228838383610455565b505050565b610235610523565b61027e6040805161010081019091528060008152602001600081526020016000815260200160608152602001606081526020016060815260200160608152602001606081525090565b6102c089898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105a69050565b85156103085761030887878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105b79050565b83156103225761032261031b85876114f1565b8290610601565b61033181846201117085610644565b600255505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610449610523565b61045281610663565b50565b8260025414610493576040517fd068bf5b000000000000000000000000000000000000000000000000000000008152600481018490526024016103bc565b61049c82610758565b6003558151600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556104de81610758565b600455516005805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103bc565b565b6105b382600080846107da565b5050565b80516000036105f2576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b805160000361063c576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c090910152565b600061065a61065286610871565b858585610c4a565b95945050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036106e2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103bc565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000806020905060208351101561076d575081515b60005b818110156107d3576107838160086115e7565b848281518110610795576107956115fe565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c92909217916107cc8161162d565b9050610770565b5050919050565b8051600003610815576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8383600281111561082857610828611589565b9081600281111561083b5761083b611589565b9052506040840182801561085157610851611589565b9081801561086157610861611589565b9052506060909301929092525050565b60606000610880610100610d29565b90506108ca6040518060400160405280600c81526020017f636f64654c6f636174696f6e000000000000000000000000000000000000000081525082610d4a90919063ffffffff16565b82516108e89060028111156108e1576108e1611589565b8290610d63565b60408051808201909152600881527f6c616e67756167650000000000000000000000000000000000000000000000006020820152610927908290610d4a565b604083015161093e9080156108e1576108e1611589565b60408051808201909152600681527f736f757263650000000000000000000000000000000000000000000000000000602082015261097d908290610d4a565b606083015161098d908290610d4a565b60a083015151156109e75760408051808201909152601081527f726571756573745369676e61747572650000000000000000000000000000000060208201526109d7908290610d4a565b60a08301516109e7908290610d98565b60c08301515115610a945760408051808201909152600481527f61726773000000000000000000000000000000000000000000000000000000006020820152610a31908290610d4a565b610a3a81610da5565b60005b8360c0015151811015610a8a57610a7a8460c001518281518110610a6357610a636115fe565b602002602001015183610d4a90919063ffffffff16565b610a838161162d565b9050610a3d565b50610a9481610dc9565b60808301515115610b9557600083602001516002811115610ab757610ab7611589565b03610aee576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e00000000000000000000000000000000006020820152610b2d908290610d4a565b610b46836020015160028111156108e1576108e1611589565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610b85908290610d4a565b6080830151610b95908290610d98565b60e08301515115610c425760408051808201909152600981527f62797465734172677300000000000000000000000000000000000000000000006020820152610bdf908290610d4a565b610be881610da5565b60005b8360e0015151811015610c3857610c288460e001518281518110610c1157610c116115fe565b602002602001015183610d9890919063ffffffff16565b610c318161162d565b9050610beb565b50610c4281610dc9565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401610caf959493929190611665565b6020604051808303816000875af1158015610cce573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cf29190611705565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b610d316111bb565b8051610d3d9083610de7565b5060006020820152919050565b610d578260038351610e61565b81516102289082610f88565b8151610d709060c2610fb0565b506105b38282604051602001610d8891815260200190565b6040516020818303038152906040525b610d578260028351610e61565b610db0816004611019565b600181602001818151610dc3919061171e565b90525050565b610dd4816007611019565b600181602001818151610dc39190611731565b604080518082019091526060815260006020820152610e07602083611744565b15610e2f57610e17602083611744565b610e22906020611731565b610e2c908361171e565b91505b602080840183905260405180855260008152908184010181811015610e5357600080fd5b604052508290505b92915050565b60178167ffffffffffffffff1611610e8e578251610e889060e0600585901b168317610fb0565b50505050565b60ff8167ffffffffffffffff1611610ed0578251610eb7906018611fe0600586901b1617610fb0565b508251610e889067ffffffffffffffff83166001611030565b61ffff8167ffffffffffffffff1611610f13578251610efa906019611fe0600586901b1617610fb0565b508251610e889067ffffffffffffffff83166002611030565b63ffffffff8167ffffffffffffffff1611610f58578251610f3f90601a611fe0600586901b1617610fb0565b508251610e889067ffffffffffffffff83166004611030565b8251610f6f90601b611fe0600586901b1617610fb0565b508251610e889067ffffffffffffffff83166008611030565b604080518082019091526060815260006020820152610fa9838384516110b5565b9392505050565b6040805180820190915260608152600060208201528251516000610fd582600161171e565b905084602001518210610ff657610ff685610ff18360026115e7565b6111a4565b845160208382010185815350805182111561100f578181525b5093949350505050565b815161022890601f611fe0600585901b1617610fb0565b6040805180820190915260608152600060208201528351516000611054828561171e565b905085602001518111156110715761107186610ff18360026115e7565b600060016110818661010061189f565b61108b9190611731565b905086518281018783198251161781525080518311156110a9578281525b50959695505050505050565b60408051808201909152606081526000602082015282518211156110d857600080fd5b83515160006110e7848361171e565b905085602001518111156111045761110486610ff18360026115e7565b85518051838201602001916000918085111561111e578482525b505050602086015b6020861061115e578051825261113d60208361171e565b915061114a60208261171e565b9050611157602087611731565b9550611126565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516111b08383610de7565b50610e888382610f88565b60405180604001604052806111e3604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611266576112666111f0565b604052919050565b600067ffffffffffffffff831115611288576112886111f0565b6112b960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8601160161121f565b90508281528383830111156112cd57600080fd5b828260208301376000602084830101529392505050565b600082601f8301126112f557600080fd5b610fa98383356020850161126e565b60008060006060848603121561131957600080fd5b83359250602084013567ffffffffffffffff8082111561133857600080fd5b611344878388016112e4565b9350604086013591508082111561135a57600080fd5b50611367868287016112e4565b9150509250925092565b60008083601f84011261138357600080fd5b50813567ffffffffffffffff81111561139b57600080fd5b6020830191508360208285010111156113b357600080fd5b9250929050565b803567ffffffffffffffff811681146113d257600080fd5b919050565b60008060008060008060008060a0898b0312156113f357600080fd5b883567ffffffffffffffff8082111561140b57600080fd5b6114178c838d01611371565b909a50985060208b013591508082111561143057600080fd5b61143c8c838d01611371565b909850965060408b013591508082111561145557600080fd5b818b0191508b601f83011261146957600080fd5b81358181111561147857600080fd5b8c60208260051b850101111561148d57600080fd5b6020830196508095505050506114a560608a016113ba565b9150608089013590509295985092959890939650565b6000602082840312156114cd57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610fa957600080fd5b600067ffffffffffffffff8084111561150c5761150c6111f0565b8360051b602061151d81830161121f565b86815291850191818101903684111561153557600080fd5b865b8481101561157d5780358681111561154f5760008081fd5b880136601f8201126115615760008081fd5b61156f36823587840161126e565b845250918301918301611537565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610e5b57610e5b6115b8565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361165e5761165e6115b8565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156116a35788810183015185820160c001528201611687565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050506116ec604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b60006020828403121561171757600080fd5b5051919050565b80820180821115610e5b57610e5b6115b8565b81810381811115610e5b57610e5b6115b8565b60008261177a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b600181815b808511156117d857817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156117be576117be6115b8565b808516156117cb57918102915b93841c9390800290611784565b509250929050565b6000826117ef57506001610e5b565b816117fc57506000610e5b565b8160018114611812576002811461181c57611838565b6001915050610e5b565b60ff84111561182d5761182d6115b8565b50506001821b610e5b565b5060208310610133831016604e8410600b841016171561185b575081810a610e5b565b611865838361177f565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611897576118976115b8565b029392505050565b6000610fa983836117e056fea164736f6c6343000813000a",
}

var FunctionsClientExampleABI = FunctionsClientExampleMetaData.ABI

var FunctionsClientExampleBin = FunctionsClientExampleMetaData.Bin

func DeployFunctionsClientExample(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *FunctionsClientExample, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsClientExampleBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsClientExample{FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
}

type FunctionsClientExample struct {
	address common.Address
	abi     abi.ABI
	FunctionsClientExampleCaller
	FunctionsClientExampleTransactor
	FunctionsClientExampleFilterer
}

type FunctionsClientExampleCaller struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleTransactor struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleFilterer struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleSession struct {
	Contract     *FunctionsClientExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsClientExampleCallerSession struct {
	Contract *FunctionsClientExampleCaller
	CallOpts bind.CallOpts
}

type FunctionsClientExampleTransactorSession struct {
	Contract     *FunctionsClientExampleTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsClientExampleRaw struct {
	Contract *FunctionsClientExample
}

type FunctionsClientExampleCallerRaw struct {
	Contract *FunctionsClientExampleCaller
}

type FunctionsClientExampleTransactorRaw struct {
	Contract *FunctionsClientExampleTransactor
}

func NewFunctionsClientExample(address common.Address, backend bind.ContractBackend) (*FunctionsClientExample, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsClientExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsClientExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExample{address: address, abi: abi, FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
}

func NewFunctionsClientExampleCaller(address common.Address, caller bind.ContractCaller) (*FunctionsClientExampleCaller, error) {
	contract, err := bindFunctionsClientExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleCaller{contract: contract}, nil
}

func NewFunctionsClientExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsClientExampleTransactor, error) {
	contract, err := bindFunctionsClientExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleTransactor{contract: contract}, nil
}

func NewFunctionsClientExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsClientExampleFilterer, error) {
	contract, err := bindFunctionsClientExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleFilterer{contract: contract}, nil
}

func bindFunctionsClientExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClientExample.Contract.FunctionsClientExampleCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.FunctionsClientExampleTransactor.contract.Transfer(opts)
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.FunctionsClientExampleTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClientExample.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.contract.Transfer(opts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "MAX_CALLBACK_GAS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsClientExample.Contract.MAXCALLBACKGAS(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsClientExample.Contract.MAXCALLBACKGAS(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) Owner() (common.Address, error) {
	return _FunctionsClientExample.Contract.Owner(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) Owner() (common.Address, error) {
	return _FunctionsClientExample.Contract.Owner(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastError(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastError")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastError() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastError(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastError() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastError(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastErrorLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastErrorLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastErrorLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastErrorLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastErrorLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastErrorLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastRequestId() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastRequestId(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastRequestId() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastRequestId(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastResponse(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastResponse")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastResponse() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastResponse(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastResponse() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastResponse(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastResponseLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastResponseLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastResponseLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastResponseLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastResponseLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastResponseLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsClientExample *FunctionsClientExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.AcceptOwnership(&_FunctionsClientExample.TransactOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.AcceptOwnership(&_FunctionsClientExample.TransactOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.HandleOracleFulfillment(&_FunctionsClientExample.TransactOpts, requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.HandleOracleFulfillment(&_FunctionsClientExample.TransactOpts, requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "sendRequest", source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.SendRequest(&_FunctionsClientExample.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.SendRequest(&_FunctionsClientExample.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.TransferOwnership(&_FunctionsClientExample.TransactOpts, to)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.TransferOwnership(&_FunctionsClientExample.TransactOpts, to)
}

type FunctionsClientExampleOwnershipTransferRequestedIterator struct {
	Event *FunctionsClientExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleOwnershipTransferRequested)
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
		it.Event = new(FunctionsClientExampleOwnershipTransferRequested)
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

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleOwnershipTransferRequestedIterator{contract: _FunctionsClientExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleOwnershipTransferRequested)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsClientExampleOwnershipTransferRequested, error) {
	event := new(FunctionsClientExampleOwnershipTransferRequested)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleOwnershipTransferredIterator struct {
	Event *FunctionsClientExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleOwnershipTransferred)
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
		it.Event = new(FunctionsClientExampleOwnershipTransferred)
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

func (it *FunctionsClientExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleOwnershipTransferredIterator{contract: _FunctionsClientExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleOwnershipTransferred)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsClientExampleOwnershipTransferred, error) {
	event := new(FunctionsClientExampleOwnershipTransferred)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleRequestFulfilledIterator struct {
	Event *FunctionsClientExampleRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleRequestFulfilled)
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
		it.Event = new(FunctionsClientExampleRequestFulfilled)
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

func (it *FunctionsClientExampleRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleRequestFulfilledIterator{contract: _FunctionsClientExample.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleRequestFulfilled)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseRequestFulfilled(log types.Log) (*FunctionsClientExampleRequestFulfilled, error) {
	event := new(FunctionsClientExampleRequestFulfilled)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleRequestSentIterator struct {
	Event *FunctionsClientExampleRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleRequestSent)
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
		it.Event = new(FunctionsClientExampleRequestSent)
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

func (it *FunctionsClientExampleRequestSentIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleRequestSentIterator{contract: _FunctionsClientExample.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleRequestSent)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseRequestSent(log types.Log) (*FunctionsClientExampleRequestSent, error) {
	event := new(FunctionsClientExampleRequestSent)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsClientExample *FunctionsClientExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsClientExample.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsClientExample.ParseOwnershipTransferRequested(log)
	case _FunctionsClientExample.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsClientExample.ParseOwnershipTransferred(log)
	case _FunctionsClientExample.abi.Events["RequestFulfilled"].ID:
		return _FunctionsClientExample.ParseRequestFulfilled(log)
	case _FunctionsClientExample.abi.Events["RequestSent"].ID:
		return _FunctionsClientExample.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsClientExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsClientExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsClientExampleRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (FunctionsClientExampleRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_FunctionsClientExample *FunctionsClientExample) Address() common.Address {
	return _FunctionsClientExample.address
}

type FunctionsClientExampleInterface interface {
	MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SLastError(opts *bind.CallOpts) ([32]byte, error)

	SLastErrorLength(opts *bind.CallOpts) (uint32, error)

	SLastRequestId(opts *bind.CallOpts) ([32]byte, error)

	SLastResponse(opts *bind.CallOpts) ([32]byte, error)

	SLastResponseLength(opts *bind.CallOpts) (uint32, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsClientExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsClientExampleOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*FunctionsClientExampleRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*FunctionsClientExampleRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
