// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr_client_example

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

type FunctionsRequest struct {
	CodeLocation    uint8
	SecretsLocation uint8
	Language        uint8
	Source          string
	Secrets         []byte
	Args            []string
}

var OCR2DRClientExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsNotPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderIsNotRegistry\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"SendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumFunctions.Location\",\"name\":\"codeLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.Location\",\"name\":\"secretsLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.CodeLanguage\",\"name\":\"language\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"}],\"internalType\":\"structFunctions.Request\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001eb338038062001eb3833981016040819052620000349162000199565b600080546001600160a01b0319166001600160a01b038316178155339081906001600160a01b038216620000af5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600280546001600160a01b0319166001600160a01b0384811691909117909155811615620000e257620000e281620000ec565b50505050620001cb565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a6565b600380546001600160a01b0319166001600160a01b03838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600060208284031215620001ac57600080fd5b81516001600160a01b0381168114620001c457600080fd5b9392505050565b611cd880620001db6000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80638da5cb5b1161008c578063d4b3917511610066578063d4b39175146101aa578063d769717e146101da578063f2fde38b146101ed578063fc2a88c31461020057600080fd5b80638da5cb5b1461015d578063b48cffea14610185578063d328a91e1461019557600080fd5b806362747e42116100bd57806362747e42146101425780636d9809a01461014b57806379ba50971461015557600080fd5b80630ca76175146100e457806329f0de3f146100f95780632c29166b14610115575b600080fd5b6100f76100f2366004611527565b610209565b005b61010260065481565b6040519081526020015b60405180910390f35b60075461012d90640100000000900463ffffffff1681565b60405163ffffffff909116815260200161010c565b61010260055481565b61012d6201117081565b6100f76102d4565b60025460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b60075461012d9063ffffffff1681565b61019d6103da565b60405161010c9190611867565b6101bd6101b83660046116dd565b6104a3565b6040516bffffffffffffffffffffffff909116815260200161010c565b6100f76101e8366004611602565b610546565b6100f76101fb3660046114d4565b610647565b61010260045481565b600083815260016020526040902054839073ffffffffffffffffffffffffffffffffffffffff163314610268576040517fa0c5ec6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e691a26102ce84848461065b565b50505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461035a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560038054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b60008054604080517fd328a91e000000000000000000000000000000000000000000000000000000008152905160609373ffffffffffffffffffffffffffffffffffffffff9093169263d328a91e9260048082019391829003018186803b15801561044457600080fd5b505afa158015610458573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261049e9190810190611594565b905090565b6000805473ffffffffffffffffffffffffffffffffffffffff1663d227d245856104cc88610729565b86866040518563ffffffff1660e01b81526004016104ed94939291906118b3565b60206040518083038186803b15801561050557600080fd5b505afa158015610519573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053d91906117ef565b95945050505050565b61054e6109a5565b6105886040805160c08101909152806000815260200160008152602001600081526020016060815260200160608152602001606081525090565b6105ca88888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a289050565b84156106125761061286868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a399050565b821561062c5761062c6106258486611b40565b8290610a80565b61063a818362011170610ac0565b6004555050505050505050565b61064f6109a5565b61065881610c7e565b50565b8260045414610699576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610351565b6106a282610d75565b6005558151600780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556106e481610d75565b600655516007805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b6060610733611305565b805161074190610100610dfd565b50610781816040518060400160405280600c81526020017f636f64654c6f636174696f6e0000000000000000000000000000000000000000815250610e77565b6107a0818460000151600181111561079b5761079b611c1c565b610e95565b6107df816040518060400160405280600881526020017f6c616e6775616765000000000000000000000000000000000000000000000000815250610e77565b6107f9818460400151600081111561079b5761079b611c1c565b610838816040518060400160405280600681526020017f736f757263650000000000000000000000000000000000000000000000000000815250610e77565b610846818460600151610e77565b60a083015151156108ec57610890816040518060400160405280600481526020017f6172677300000000000000000000000000000000000000000000000000000000815250610e77565b61089981610ece565b60005b8360a00151518110156108e2576108d0828560a0015183815181106108c3576108c3611c4b565b6020026020010151610e77565b806108da81611b79565b91505061089c565b506108ec81610ef2565b6080830151511561099d57610936816040518060400160405280600f81526020017f736563726574734c6f636174696f6e0000000000000000000000000000000000815250610e77565b610950818460200151600181111561079b5761079b611c1c565b61098f816040518060400160405280600781526020017f7365637265747300000000000000000000000000000000000000000000000000815250610e77565b61099d818460800151610f10565b515192915050565b60025473ffffffffffffffffffffffffffffffffffffffff163314610a26576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610351565b565b610a358260008084610f1d565b5050565b8051610a71576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006020830152608090910152565b8051610ab8576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60008054819073ffffffffffffffffffffffffffffffffffffffff166328242b0485610aeb88610729565b866040518463ffffffff1660e01b8152600401610b0a9392919061187a565b602060405180830381600087803b158015610b2457600080fd5b505af1158015610b38573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b5c919061150e565b905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635ab1bd536040518163ffffffff1660e01b815260040160206040518083038186803b158015610bc457600080fd5b505afa158015610bd8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bfc91906114f1565b60008281526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9490941693909317909255905182917f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db891a2949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610cfe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610351565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600080600060209050602084511015610d8c575082515b60005b81811015610df457610da2816008611aec565b858281518110610db457610db4611c4b565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c929092179180610dec81611b79565b915050610d8f565b50909392505050565b604080518082019091526060815260006020820152610e1d602083611bb2565b15610e4557610e2d602083611bb2565b610e38906020611b29565b610e4290836119ad565b91505b602080840183905260405180855260008152908184010181811015610e6957600080fd5b604052508290505b92915050565b610e848260038351610fb1565b8151610e9090826110d2565b505050565b8151610ea29060c26110fa565b50610a358282604051602001610eba91815260200190565b604051602081830303815290604052610f10565b610ed9816004611163565b600181602001818151610eec91906119ad565b90525050565b610efd816007611163565b600181602001818151610eec9190611b29565b610e848260028351610fb1565b8051610f55576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836001811115610f6857610f68611c1c565b90816001811115610f7b57610f7b611c1c565b90525060408401828015610f9157610f91611c1c565b90818015610fa157610fa1611c1c565b9052506060909301929092525050565b60178167ffffffffffffffff1611610fd85782516102ce9060e0600585901b1683176110fa565b60ff8167ffffffffffffffff161161101a578251611001906018611fe0600586901b16176110fa565b5082516102ce9067ffffffffffffffff8316600161117a565b61ffff8167ffffffffffffffff161161105d578251611044906019611fe0600586901b16176110fa565b5082516102ce9067ffffffffffffffff8316600261117a565b63ffffffff8167ffffffffffffffff16116110a257825161108990601a611fe0600586901b16176110fa565b5082516102ce9067ffffffffffffffff8316600461117a565b82516110b990601b611fe0600586901b16176110fa565b5082516102ce9067ffffffffffffffff8316600861117a565b6040805180820190915260608152600060208201526110f3838384516111ff565b9392505050565b604080518082019091526060815260006020820152825151600061111f8260016119ad565b905084602001518210611140576111408561113b836002611aec565b6112ee565b8451602083820101858153508051821115611159578181525b5093949350505050565b8151610e9090601f611fe0600585901b16176110fa565b604080518082019091526060815260006020820152835151600061119e82856119ad565b905085602001518111156111bb576111bb8661113b836002611aec565b600060016111cb86610100611a26565b6111d59190611b29565b905086518281018783198251161781525080518311156111f3578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561122257600080fd5b835151600061123184836119ad565b9050856020015181111561124e5761124e8661113b836002611aec565b855180518382016020019160009180851115611268578482525b505050602086015b602086106112a857805182526112876020836119ad565b91506112946020826119ad565b90506112a1602087611b29565b9550611270565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516112fa8383610dfd565b506102ce83826110d2565b604051806040016040528061132d604051806040016040528060608152602001600081525090565b8152602001600081525090565b600067ffffffffffffffff8084111561135557611355611c7a565b8360051b6020611366818301611918565b86815293508084018583810189101561137e57600080fd5b60009350835b888110156113b95781358681111561139a578586fd5b6113a68b828b0161142f565b8452509183019190830190600101611384565b5050505050509392505050565b600082601f8301126113d757600080fd5b6110f38383356020850161133a565b60008083601f8401126113f857600080fd5b50813567ffffffffffffffff81111561141057600080fd5b60208301915083602082850101111561142857600080fd5b9250929050565b600082601f83011261144057600080fd5b813561145361144e82611967565b611918565b81815284602083860101111561146857600080fd5b816020850160208301376000918101602001919091529392505050565b80356001811061149457600080fd5b919050565b80356002811061149457600080fd5b803563ffffffff8116811461149457600080fd5b803567ffffffffffffffff8116811461149457600080fd5b6000602082840312156114e657600080fd5b81356110f381611ca9565b60006020828403121561150357600080fd5b81516110f381611ca9565b60006020828403121561152057600080fd5b5051919050565b60008060006060848603121561153c57600080fd5b83359250602084013567ffffffffffffffff8082111561155b57600080fd5b6115678783880161142f565b9350604086013591508082111561157d57600080fd5b5061158a8682870161142f565b9150509250925092565b6000602082840312156115a657600080fd5b815167ffffffffffffffff8111156115bd57600080fd5b8201601f810184136115ce57600080fd5b80516115dc61144e82611967565b8181528560208385010111156115f157600080fd5b61053d826020830160208601611b4d565b60008060008060008060006080888a03121561161d57600080fd5b873567ffffffffffffffff8082111561163557600080fd5b6116418b838c016113e6565b909950975060208a013591508082111561165a57600080fd5b6116668b838c016113e6565b909750955060408a013591508082111561167f57600080fd5b818a0191508a601f83011261169357600080fd5b8135818111156116a257600080fd5b8b60208260051b85010111156116b757600080fd5b6020830195508094505050506116cf606089016114bc565b905092959891949750929550565b600080600080608085870312156116f357600080fd5b843567ffffffffffffffff8082111561170b57600080fd5b9086019060c0828903121561171f57600080fd5b6117276118ef565b61173083611499565b815261173e60208401611499565b602082015261174f60408401611485565b604082015260608301358281111561176657600080fd5b6117728a82860161142f565b60608301525060808301358281111561178a57600080fd5b6117968a82860161142f565b60808301525060a0830135828111156117ae57600080fd5b6117ba8a8286016113c6565b60a08301525095506117d1915050602086016114bc565b92506117df604086016114a8565b9396929550929360600135925050565b60006020828403121561180157600080fd5b81516bffffffffffffffffffffffff811681146110f357600080fd5b60008151808452611835816020860160208601611b4d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006110f3602083018461181d565b67ffffffffffffffff8416815260606020820152600061189d606083018561181d565b905063ffffffff83166040830152949350505050565b67ffffffffffffffff851681526080602082015260006118d6608083018661181d565b63ffffffff949094166040830152506060015292915050565b60405160c0810167ffffffffffffffff8111828210171561191257611912611c7a565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561195f5761195f611c7a565b604052919050565b600067ffffffffffffffff82111561198157611981611c7a565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082198211156119c0576119c0611bed565b500190565b600181815b80851115611a1e57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611a0457611a04611bed565b80851615611a1157918102915b93841c93908002906119ca565b509250929050565b60006110f38383600082611a3c57506001610e71565b81611a4957506000610e71565b8160018114611a5f5760028114611a6957611a85565b6001915050610e71565b60ff841115611a7a57611a7a611bed565b50506001821b610e71565b5060208310610133831016604e8410600b8410161715611aa8575081810a610e71565b611ab283836119c5565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611ae457611ae4611bed565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611b2457611b24611bed565b500290565b600082821015611b3b57611b3b611bed565b500390565b60006110f336848461133a565b60005b83811015611b68578181015183820152602001611b50565b838111156102ce5750506000910152565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611bab57611bab611bed565b5060010190565b600082611be8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461065857600080fdfea164736f6c6343000806000a",
}

var OCR2DRClientExampleABI = OCR2DRClientExampleMetaData.ABI

var OCR2DRClientExampleBin = OCR2DRClientExampleMetaData.Bin

func DeployOCR2DRClientExample(auth *bind.TransactOpts, backend bind.ContractBackend, oracle common.Address) (common.Address, *types.Transaction, *OCR2DRClientExample, error) {
	parsed, err := OCR2DRClientExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DRClientExampleBin), backend, oracle)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2DRClientExample{OCR2DRClientExampleCaller: OCR2DRClientExampleCaller{contract: contract}, OCR2DRClientExampleTransactor: OCR2DRClientExampleTransactor{contract: contract}, OCR2DRClientExampleFilterer: OCR2DRClientExampleFilterer{contract: contract}}, nil
}

type OCR2DRClientExample struct {
	address common.Address
	abi     abi.ABI
	OCR2DRClientExampleCaller
	OCR2DRClientExampleTransactor
	OCR2DRClientExampleFilterer
}

type OCR2DRClientExampleCaller struct {
	contract *bind.BoundContract
}

type OCR2DRClientExampleTransactor struct {
	contract *bind.BoundContract
}

type OCR2DRClientExampleFilterer struct {
	contract *bind.BoundContract
}

type OCR2DRClientExampleSession struct {
	Contract     *OCR2DRClientExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DRClientExampleCallerSession struct {
	Contract *OCR2DRClientExampleCaller
	CallOpts bind.CallOpts
}

type OCR2DRClientExampleTransactorSession struct {
	Contract     *OCR2DRClientExampleTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DRClientExampleRaw struct {
	Contract *OCR2DRClientExample
}

type OCR2DRClientExampleCallerRaw struct {
	Contract *OCR2DRClientExampleCaller
}

type OCR2DRClientExampleTransactorRaw struct {
	Contract *OCR2DRClientExampleTransactor
}

func NewOCR2DRClientExample(address common.Address, backend bind.ContractBackend) (*OCR2DRClientExample, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DRClientExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DRClientExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExample{address: address, abi: abi, OCR2DRClientExampleCaller: OCR2DRClientExampleCaller{contract: contract}, OCR2DRClientExampleTransactor: OCR2DRClientExampleTransactor{contract: contract}, OCR2DRClientExampleFilterer: OCR2DRClientExampleFilterer{contract: contract}}, nil
}

func NewOCR2DRClientExampleCaller(address common.Address, caller bind.ContractCaller) (*OCR2DRClientExampleCaller, error) {
	contract, err := bindOCR2DRClientExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleCaller{contract: contract}, nil
}

func NewOCR2DRClientExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DRClientExampleTransactor, error) {
	contract, err := bindOCR2DRClientExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleTransactor{contract: contract}, nil
}

func NewOCR2DRClientExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DRClientExampleFilterer, error) {
	contract, err := bindOCR2DRClientExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleFilterer{contract: contract}, nil
}

func bindOCR2DRClientExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2DRClientExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRClientExample.Contract.OCR2DRClientExampleCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.OCR2DRClientExampleTransactor.contract.Transfer(opts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.OCR2DRClientExampleTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRClientExample.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.contract.Transfer(opts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "MAX_CALLBACK_GAS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) MAXCALLBACKGAS() (uint32, error) {
	return _OCR2DRClientExample.Contract.MAXCALLBACKGAS(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) MAXCALLBACKGAS() (uint32, error) {
	return _OCR2DRClientExample.Contract.MAXCALLBACKGAS(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) EstimateCost(opts *bind.CallOpts, req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "estimateCost", req, subscriptionId, gasLimit, gasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) EstimateCost(req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _OCR2DRClientExample.Contract.EstimateCost(&_OCR2DRClientExample.CallOpts, req, subscriptionId, gasLimit, gasPrice)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) EstimateCost(req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _OCR2DRClientExample.Contract.EstimateCost(&_OCR2DRClientExample.CallOpts, req, subscriptionId, gasLimit, gasPrice)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) GetDONPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "getDONPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DRClientExample.Contract.GetDONPublicKey(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DRClientExample.Contract.GetDONPublicKey(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastError(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastError")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastError() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastError(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastError() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastError(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastErrorLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastErrorLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastErrorLength() (uint32, error) {
	return _OCR2DRClientExample.Contract.LastErrorLength(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastErrorLength() (uint32, error) {
	return _OCR2DRClientExample.Contract.LastErrorLength(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastRequestId() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastRequestId(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastRequestId() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastRequestId(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastResponse(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastResponse")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastResponse() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastResponse(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastResponse() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastResponse(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastResponseLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastResponseLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastResponseLength() (uint32, error) {
	return _OCR2DRClientExample.Contract.LastResponseLength(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastResponseLength() (uint32, error) {
	return _OCR2DRClientExample.Contract.LastResponseLength(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) Owner() (common.Address, error) {
	return _OCR2DRClientExample.Contract.Owner(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) Owner() (common.Address, error) {
	return _OCR2DRClientExample.Contract.Owner(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) SendRequest(opts *bind.TransactOpts, source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "SendRequest", source, secrets, args, subscriptionId)
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) SendRequest(source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.SendRequest(&_OCR2DRClientExample.TransactOpts, source, secrets, args, subscriptionId)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) SendRequest(source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.SendRequest(&_OCR2DRClientExample.TransactOpts, source, secrets, args, subscriptionId)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "acceptOwnership")
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.AcceptOwnership(&_OCR2DRClientExample.TransactOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.AcceptOwnership(&_OCR2DRClientExample.TransactOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.HandleOracleFulfillment(&_OCR2DRClientExample.TransactOpts, requestId, response, err)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.HandleOracleFulfillment(&_OCR2DRClientExample.TransactOpts, requestId, response, err)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.TransferOwnership(&_OCR2DRClientExample.TransactOpts, to)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.TransferOwnership(&_OCR2DRClientExample.TransactOpts, to)
}

type OCR2DRClientExampleOwnershipTransferRequestedIterator struct {
	Event *OCR2DRClientExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleOwnershipTransferRequested)
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
		it.Event = new(OCR2DRClientExampleOwnershipTransferRequested)
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

func (it *OCR2DRClientExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleOwnershipTransferRequestedIterator{contract: _OCR2DRClientExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleOwnershipTransferRequested)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2DRClientExampleOwnershipTransferRequested, error) {
	event := new(OCR2DRClientExampleOwnershipTransferRequested)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientExampleOwnershipTransferredIterator struct {
	Event *OCR2DRClientExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleOwnershipTransferred)
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
		it.Event = new(OCR2DRClientExampleOwnershipTransferred)
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

func (it *OCR2DRClientExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleOwnershipTransferredIterator{contract: _OCR2DRClientExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleOwnershipTransferred)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2DRClientExampleOwnershipTransferred, error) {
	event := new(OCR2DRClientExampleOwnershipTransferred)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientExampleRequestFulfilledIterator struct {
	Event *OCR2DRClientExampleRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleRequestFulfilled)
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
		it.Event = new(OCR2DRClientExampleRequestFulfilled)
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

func (it *OCR2DRClientExampleRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleRequestFulfilledIterator{contract: _OCR2DRClientExample.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleRequestFulfilled)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseRequestFulfilled(log types.Log) (*OCR2DRClientExampleRequestFulfilled, error) {
	event := new(OCR2DRClientExampleRequestFulfilled)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientExampleRequestSentIterator struct {
	Event *OCR2DRClientExampleRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleRequestSent)
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
		it.Event = new(OCR2DRClientExampleRequestSent)
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

func (it *OCR2DRClientExampleRequestSentIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleRequestSentIterator{contract: _OCR2DRClientExample.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleRequestSent)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseRequestSent(log types.Log) (*OCR2DRClientExampleRequestSent, error) {
	event := new(OCR2DRClientExampleRequestSent)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OCR2DRClientExample *OCR2DRClientExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2DRClientExample.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR2DRClientExample.ParseOwnershipTransferRequested(log)
	case _OCR2DRClientExample.abi.Events["OwnershipTransferred"].ID:
		return _OCR2DRClientExample.ParseOwnershipTransferred(log)
	case _OCR2DRClientExample.abi.Events["RequestFulfilled"].ID:
		return _OCR2DRClientExample.ParseRequestFulfilled(log)
	case _OCR2DRClientExample.abi.Events["RequestSent"].ID:
		return _OCR2DRClientExample.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2DRClientExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR2DRClientExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCR2DRClientExampleRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (OCR2DRClientExampleRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_OCR2DRClientExample *OCR2DRClientExample) Address() common.Address {
	return _OCR2DRClientExample.address
}

type OCR2DRClientExampleInterface interface {
	MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error)

	EstimateCost(opts *bind.CallOpts, req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error)

	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	LastError(opts *bind.CallOpts) ([32]byte, error)

	LastErrorLength(opts *bind.CallOpts) (uint32, error)

	LastRequestId(opts *bind.CallOpts) ([32]byte, error)

	LastResponse(opts *bind.CallOpts) ([32]byte, error)

	LastResponseLength(opts *bind.CallOpts) (uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SendRequest(opts *bind.TransactOpts, source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR2DRClientExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR2DRClientExampleOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*OCR2DRClientExampleRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*OCR2DRClientExampleRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
