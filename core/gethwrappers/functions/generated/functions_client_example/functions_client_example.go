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
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001a4838038062001a48833981016040819052620000349162000180565b6001600160a01b0381166080523380600081620000985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cb57620000cb81620000d5565b50505050620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b608051611873620001d5600039600081816101c601526109f301526118736000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80636d9809a011610081578063b1e217491161005b578063b1e2174914610182578063f2fde38b1461018b578063f7b4c06f1461019e57600080fd5b80636d9809a01461014857806379ba5097146101525780638da5cb5b1461015a57600080fd5b806342748b2a116100b257806342748b2a146100ff5780634b0795a81461012c5780635fa353e71461013557600080fd5b80630ca76175146100ce5780633944ea3a146100e3575b600080fd5b6100e16100dc3660046112bf565b6101ae565b005b6100ec60035481565b6040519081526020015b60405180910390f35b60055461011790640100000000900463ffffffff1681565b60405163ffffffff90911681526020016100f6565b6100ec60045481565b6100e1610143366004611392565b610258565b6101176201117081565b6100e161036a565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b6100ec60025481565b6100e1610199366004611476565b61046c565b6005546101179063ffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461021d576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610228838383610480565b60405183907f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e690600090a2505050565b61026061054e565b6102a16040805160e0810190915280600081526020016000815260200160008152602001606081526020016060815260200160608152602001606081525090565b6102e389898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105d19050565b851561032b5761032b87878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105e29050565b83156103455761034561033e85876114ac565b829061062c565b61035c6103518261066f565b8462011170856109ee565b600255505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103f0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61047461054e565b61047d81610acd565b50565b82600254146104be576040517fd068bf5b000000000000000000000000000000000000000000000000000000008152600481018490526024016103e7565b6104c782610bc2565b6003558151600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90921691909117905561050981610bc2565b600455516005805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105cf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103e7565b565b6105de8260008084610c44565b5050565b805160000361061d576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b8051600003610667576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b6060600061067e610100610cdb565b90506106c86040518060400160405280600c81526020017f636f64654c6f636174696f6e000000000000000000000000000000000000000081525082610cfc90919063ffffffff16565b82516106e69060028111156106df576106df611544565b8290610d1a565b60408051808201909152600881527f6c616e67756167650000000000000000000000000000000000000000000000006020820152610725908290610cfc565b604083015161073c9080156106df576106df611544565b60408051808201909152600681527f736f757263650000000000000000000000000000000000000000000000000000602082015261077b908290610cfc565b606083015161078b908290610cfc565b60a083015151156108385760408051808201909152600481527f617267730000000000000000000000000000000000000000000000000000000060208201526107d5908290610cfc565b6107de81610d53565b60005b8360a001515181101561082e5761081e8460a00151828151811061080757610807611573565b602002602001015183610cfc90919063ffffffff16565b610827816115d1565b90506107e1565b5061083881610d77565b608083015151156109395760008360200151600281111561085b5761085b611544565b03610892576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e000000000000000000000000000000000060208201526108d1908290610cfc565b6108ea836020015160028111156106df576106df611544565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610929908290610cfc565b6080830151610939908290610d95565b60c083015151156109e65760408051808201909152600981527f62797465734172677300000000000000000000000000000000000000000000006020820152610983908290610cfc565b61098c81610d53565b60005b8360c00151518110156109dc576109cc8460c0015182815181106109b5576109b5611573565b602002602001015183610d9590919063ffffffff16565b6109d5816115d1565b905061098f565b506109e681610d77565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401610a53959493929190611609565b6020604051808303816000875af1158015610a72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a9691906116a9565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610b4c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103e7565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060209050602083511015610bd7575081515b60005b81811015610c3d57610bed8160086116c2565b848281518110610bff57610bff611573565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9290921791610c36816115d1565b9050610bda565b5050919050565b8051600003610c7f576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610c9257610c92611544565b90816002811115610ca557610ca5611544565b90525060408401828015610cbb57610cbb611544565b90818015610ccb57610ccb611544565b9052506060909301929092525050565b610ce3611176565b8051610cef9083610da2565b5060006020820152919050565b610d098260038351610e1c565b8151610d159082610f43565b505050565b8151610d279060c2610f6b565b506105de8282604051602001610d3f91815260200190565b604051602081830303815290604052610d95565b610d5e816004610fd4565b600181602001818151610d7191906116d9565b90525050565b610d82816007610fd4565b600181602001818151610d7191906116ec565b610d098260028351610e1c565b604080518082019091526060815260006020820152610dc26020836116ff565b15610dea57610dd26020836116ff565b610ddd9060206116ec565b610de790836116d9565b91505b602080840183905260405180855260008152908184010181811015610e0e57600080fd5b604052508290505b92915050565b60178167ffffffffffffffff1611610e49578251610e439060e0600585901b168317610f6b565b50505050565b60ff8167ffffffffffffffff1611610e8b578251610e72906018611fe0600586901b1617610f6b565b508251610e439067ffffffffffffffff83166001610feb565b61ffff8167ffffffffffffffff1611610ece578251610eb5906019611fe0600586901b1617610f6b565b508251610e439067ffffffffffffffff83166002610feb565b63ffffffff8167ffffffffffffffff1611610f13578251610efa90601a611fe0600586901b1617610f6b565b508251610e439067ffffffffffffffff83166004610feb565b8251610f2a90601b611fe0600586901b1617610f6b565b508251610e439067ffffffffffffffff83166008610feb565b604080518082019091526060815260006020820152610f6483838451611070565b9392505050565b6040805180820190915260608152600060208201528251516000610f908260016116d9565b905084602001518210610fb157610fb185610fac8360026116c2565b61115f565b8451602083820101858153508051821115610fca578181525b5093949350505050565b8151610d1590601f611fe0600585901b1617610f6b565b604080518082019091526060815260006020820152835151600061100f82856116d9565b9050856020015181111561102c5761102c86610fac8360026116c2565b6000600161103c8661010061185a565b61104691906116ec565b90508651828101878319825116178152508051831115611064578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561109357600080fd5b83515160006110a284836116d9565b905085602001518111156110bf576110bf86610fac8360026116c2565b8551805183820160200191600091808511156110d9578482525b505050602086015b6020861061111957805182526110f86020836116d9565b91506111056020826116d9565b90506111126020876116ec565b95506110e1565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b815161116b8383610da2565b50610e438382610f43565b604051806040016040528061119e604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611221576112216111ab565b604052919050565b600067ffffffffffffffff831115611243576112436111ab565b61127460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116016111da565b905082815283838301111561128857600080fd5b828260208301376000602084830101529392505050565b600082601f8301126112b057600080fd5b610f6483833560208501611229565b6000806000606084860312156112d457600080fd5b83359250602084013567ffffffffffffffff808211156112f357600080fd5b6112ff8783880161129f565b9350604086013591508082111561131557600080fd5b506113228682870161129f565b9150509250925092565b60008083601f84011261133e57600080fd5b50813567ffffffffffffffff81111561135657600080fd5b60208301915083602082850101111561136e57600080fd5b9250929050565b803567ffffffffffffffff8116811461138d57600080fd5b919050565b60008060008060008060008060a0898b0312156113ae57600080fd5b883567ffffffffffffffff808211156113c657600080fd5b6113d28c838d0161132c565b909a50985060208b01359150808211156113eb57600080fd5b6113f78c838d0161132c565b909850965060408b013591508082111561141057600080fd5b818b0191508b601f83011261142457600080fd5b81358181111561143357600080fd5b8c60208260051b850101111561144857600080fd5b60208301965080955050505061146060608a01611375565b9150608089013590509295985092959890939650565b60006020828403121561148857600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610f6457600080fd5b600067ffffffffffffffff808411156114c7576114c76111ab565b8360051b60206114d88183016111da565b8681529185019181810190368411156114f057600080fd5b865b848110156115385780358681111561150a5760008081fd5b880136601f82011261151c5760008081fd5b61152a368235878401611229565b8452509183019183016114f2565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611602576116026115a2565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156116475788810183015185820160c00152820161162b565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010192505050611690604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b6000602082840312156116bb57600080fd5b5051919050565b8082028115828204841417610e1657610e166115a2565b80820180821115610e1657610e166115a2565b81810381811115610e1657610e166115a2565b600082611735577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b600181815b8085111561179357817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611779576117796115a2565b8085161561178657918102915b93841c939080029061173f565b509250929050565b6000826117aa57506001610e16565b816117b757506000610e16565b81600181146117cd57600281146117d7576117f3565b6001915050610e16565b60ff8411156117e8576117e86115a2565b50506001821b610e16565b5060208310610133831016604e8410600b8410161715611816575081810a610e16565b611820838361173a565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611852576118526115a2565b029392505050565b6000610f64838361179b56fea164736f6c6343000813000a",
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
	return address, tx, &FunctionsClientExample{address: address, abi: *parsed, FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
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
