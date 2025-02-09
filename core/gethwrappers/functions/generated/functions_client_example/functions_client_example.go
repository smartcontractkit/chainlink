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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405260058054600160401b600160e01b0319167b234a5fb5bd614a7aa2ffab244d603abfa0ac5c5c000000000000000017905534801561004157600080fd5b506005546801000000000000000090046001600160a01b031660805233806000816100b35760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100e3576100e3816100eb565b505050610194565b336001600160a01b038216036101435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016100aa565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6080516118c16101b6600039600081816102140152610a4101526118c16000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806379ba509711610081578063f2fde38b1161005b578063f2fde38b146101ad578063f7b4c06f146101c0578063f887ea40146101d057600080fd5b806379ba50971461015d5780638da5cb5b14610165578063b1e21749146101a457600080fd5b80634b0795a8116100b25780634b0795a8146101375780635fa353e7146101405780636d9809a01461015357600080fd5b80630ca76175146100d95780633944ea3a146100ee57806342748b2a1461010a575b600080fd5b6100ec6100e736600461130d565b6101fc565b005b6100f760035481565b6040519081526020015b60405180910390f35b60055461012290640100000000900463ffffffff1681565b60405163ffffffff9091168152602001610101565b6100f760045481565b6100ec61014e3660046113e0565b6102a6565b6101226201117081565b6100ec6103b8565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610101565b6100f760025481565b6100ec6101bb3660046114c4565b6104ba565b6005546101229063ffffffff1681565b60055461017f9068010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461026b576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102768383836104ce565b60405183907f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e690600090a2505050565b6102ae61059c565b6102ef6040805160e0810190915280600081526020016000815260200160008152602001606081526020016060815260200160608152602001606081525090565b61033189898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250859392505061061f9050565b85156103795761037987878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506106309050565b83156103935761039361038c85876114fa565b829061067a565b6103aa61039f826106bd565b846201117085610a3c565b600255505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461043e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104c261059c565b6104cb81610b1b565b50565b826002541461050c576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610435565b61051582610c10565b6003558151600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90921691909117905561055781610c10565b600455516005805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461061d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610435565b565b61062c8260008084610c92565b5050565b805160000361066b576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b80516000036106b5576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b606060006106cc610100610d29565b90506107166040518060400160405280600c81526020017f636f64654c6f636174696f6e000000000000000000000000000000000000000081525082610d4a90919063ffffffff16565b825161073490600281111561072d5761072d611592565b8290610d68565b60408051808201909152600881527f6c616e67756167650000000000000000000000000000000000000000000000006020820152610773908290610d4a565b604083015161078a90801561072d5761072d611592565b60408051808201909152600681527f736f75726365000000000000000000000000000000000000000000000000000060208201526107c9908290610d4a565b60608301516107d9908290610d4a565b60a083015151156108865760408051808201909152600481527f61726773000000000000000000000000000000000000000000000000000000006020820152610823908290610d4a565b61082c81610da1565b60005b8360a001515181101561087c5761086c8460a001518281518110610855576108556115c1565b602002602001015183610d4a90919063ffffffff16565b6108758161161f565b905061082f565b5061088681610dc5565b60808301515115610987576000836020015160028111156108a9576108a9611592565b036108e0576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e0000000000000000000000000000000000602082015261091f908290610d4a565b6109388360200151600281111561072d5761072d611592565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610977908290610d4a565b6080830151610987908290610de3565b60c08301515115610a345760408051808201909152600981527f627974657341726773000000000000000000000000000000000000000000000060208201526109d1908290610d4a565b6109da81610da1565b60005b8360c0015151811015610a2a57610a1a8460c001518281518110610a0357610a036115c1565b602002602001015183610de390919063ffffffff16565b610a238161161f565b90506109dd565b50610a3481610dc5565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401610aa1959493929190611657565b6020604051808303816000875af1158015610ac0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ae491906116f7565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610b9a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610435565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060209050602083511015610c25575081515b60005b81811015610c8b57610c3b816008611710565b848281518110610c4d57610c4d6115c1565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9290921791610c848161161f565b9050610c28565b5050919050565b8051600003610ccd576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610ce057610ce0611592565b90816002811115610cf357610cf3611592565b90525060408401828015610d0957610d09611592565b90818015610d1957610d19611592565b9052506060909301929092525050565b610d316111c4565b8051610d3d9083610df0565b5060006020820152919050565b610d578260038351610e6a565b8151610d639082610f91565b505050565b8151610d759060c2610fb9565b5061062c8282604051602001610d8d91815260200190565b604051602081830303815290604052610de3565b610dac816004611022565b600181602001818151610dbf9190611727565b90525050565b610dd0816007611022565b600181602001818151610dbf919061173a565b610d578260028351610e6a565b604080518082019091526060815260006020820152610e1060208361174d565b15610e3857610e2060208361174d565b610e2b90602061173a565b610e359083611727565b91505b602080840183905260405180855260008152908184010181811015610e5c57600080fd5b604052508290505b92915050565b60178167ffffffffffffffff1611610e97578251610e919060e0600585901b168317610fb9565b50505050565b60ff8167ffffffffffffffff1611610ed9578251610ec0906018611fe0600586901b1617610fb9565b508251610e919067ffffffffffffffff83166001611039565b61ffff8167ffffffffffffffff1611610f1c578251610f03906019611fe0600586901b1617610fb9565b508251610e919067ffffffffffffffff83166002611039565b63ffffffff8167ffffffffffffffff1611610f61578251610f4890601a611fe0600586901b1617610fb9565b508251610e919067ffffffffffffffff83166004611039565b8251610f7890601b611fe0600586901b1617610fb9565b508251610e919067ffffffffffffffff83166008611039565b604080518082019091526060815260006020820152610fb2838384516110be565b9392505050565b6040805180820190915260608152600060208201528251516000610fde826001611727565b905084602001518210610fff57610fff85610ffa836002611710565b6111ad565b8451602083820101858153508051821115611018578181525b5093949350505050565b8151610d6390601f611fe0600585901b1617610fb9565b604080518082019091526060815260006020820152835151600061105d8285611727565b9050856020015181111561107a5761107a86610ffa836002611710565b6000600161108a866101006118a8565b611094919061173a565b905086518281018783198251161781525080518311156110b2578281525b50959695505050505050565b60408051808201909152606081526000602082015282518211156110e157600080fd5b83515160006110f08483611727565b9050856020015181111561110d5761110d86610ffa836002611710565b855180518382016020019160009180851115611127578482525b505050602086015b602086106111675780518252611146602083611727565b9150611153602082611727565b905061116060208761173a565b955061112f565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516111b98383610df0565b50610e918382610f91565b60405180604001604052806111ec604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561126f5761126f6111f9565b604052919050565b600067ffffffffffffffff831115611291576112916111f9565b6112c260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601611228565b90508281528383830111156112d657600080fd5b828260208301376000602084830101529392505050565b600082601f8301126112fe57600080fd5b610fb283833560208501611277565b60008060006060848603121561132257600080fd5b83359250602084013567ffffffffffffffff8082111561134157600080fd5b61134d878388016112ed565b9350604086013591508082111561136357600080fd5b50611370868287016112ed565b9150509250925092565b60008083601f84011261138c57600080fd5b50813567ffffffffffffffff8111156113a457600080fd5b6020830191508360208285010111156113bc57600080fd5b9250929050565b803567ffffffffffffffff811681146113db57600080fd5b919050565b60008060008060008060008060a0898b0312156113fc57600080fd5b883567ffffffffffffffff8082111561141457600080fd5b6114208c838d0161137a565b909a50985060208b013591508082111561143957600080fd5b6114458c838d0161137a565b909850965060408b013591508082111561145e57600080fd5b818b0191508b601f83011261147257600080fd5b81358181111561148157600080fd5b8c60208260051b850101111561149657600080fd5b6020830196508095505050506114ae60608a016113c3565b9150608089013590509295985092959890939650565b6000602082840312156114d657600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610fb257600080fd5b600067ffffffffffffffff80841115611515576115156111f9565b8360051b6020611526818301611228565b86815291850191818101903684111561153e57600080fd5b865b84811015611586578035868111156115585760008081fd5b880136601f82011261156a5760008081fd5b611578368235878401611277565b845250918301918301611540565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611650576116506115f0565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156116955788810183015185820160c001528201611679565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050506116de604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b60006020828403121561170957600080fd5b5051919050565b8082028115828204841417610e6457610e646115f0565b80820180821115610e6457610e646115f0565b81810381811115610e6457610e646115f0565b600082611783577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b600181815b808511156117e157817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156117c7576117c76115f0565b808516156117d457918102915b93841c939080029061178d565b509250929050565b6000826117f857506001610e64565b8161180557506000610e64565b816001811461181b576002811461182557611841565b6001915050610e64565b60ff841115611836576118366115f0565b50506001821b610e64565b5060208310610133831016604e8410600b8410161715611864575081810a610e64565b61186e8383611788565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156118a0576118a06115f0565b029392505050565b6000610fb283836117e956fea164736f6c6343000813000a",
}

var FunctionsClientExampleABI = FunctionsClientExampleMetaData.ABI

var FunctionsClientExampleBin = FunctionsClientExampleMetaData.Bin

func DeployFunctionsClientExample(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FunctionsClientExample, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsClientExampleBin), backend)
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

func (_FunctionsClientExample *FunctionsClientExampleCaller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) Router() (common.Address, error) {
	return _FunctionsClientExample.Contract.Router(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) Router() (common.Address, error) {
	return _FunctionsClientExample.Contract.Router(&_FunctionsClientExample.CallOpts)
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

	Router(opts *bind.CallOpts) (common.Address, error)

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
