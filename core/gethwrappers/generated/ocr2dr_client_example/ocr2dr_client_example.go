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
	Bin: "0x60806040523480156200001157600080fd5b5060405162001edf38038062001edf833981016040819052620000349162000199565b600080546001600160a01b0319166001600160a01b038316178155339081906001600160a01b038216620000af5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600280546001600160a01b0319166001600160a01b0384811691909117909155811615620000e257620000e281620000ec565b50505050620001cb565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a6565b600380546001600160a01b0319166001600160a01b03838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600060208284031215620001ac57600080fd5b81516001600160a01b0381168114620001c457600080fd5b9392505050565b611d0480620001db6000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80638da5cb5b1161008c578063d4b3917511610066578063d4b39175146101aa578063d769717e146101da578063f2fde38b146101ed578063fc2a88c31461020057600080fd5b80638da5cb5b1461015d578063b48cffea14610185578063d328a91e1461019557600080fd5b806362747e42116100bd57806362747e42146101425780636d9809a01461014b57806379ba50971461015557600080fd5b80630ca76175146100e457806329f0de3f146100f95780632c29166b14610115575b600080fd5b6100f76100f2366004611553565b610209565b005b61010260065481565b6040519081526020015b60405180910390f35b60075461012d90640100000000900463ffffffff1681565b60405163ffffffff909116815260200161010c565b61010260055481565b61012d6201117081565b6100f76102d4565b60025460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b60075461012d9063ffffffff1681565b61019d6103da565b60405161010c9190611893565b6101bd6101b8366004611709565b6104a3565b6040516bffffffffffffffffffffffff909116815260200161010c565b6100f76101e836600461162e565b610546565b6100f76101fb366004611500565b610647565b61010260045481565b600083815260016020526040902054839073ffffffffffffffffffffffffffffffffffffffff163314610268576040517fa0c5ec6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e691a26102ce84848461065b565b50505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461035a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560038054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b60008054604080517fd328a91e000000000000000000000000000000000000000000000000000000008152905160609373ffffffffffffffffffffffffffffffffffffffff9093169263d328a91e9260048082019391829003018186803b15801561044457600080fd5b505afa158015610458573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261049e91908101906115c0565b905090565b6000805473ffffffffffffffffffffffffffffffffffffffff1663d227d245856104cc88610729565b86866040518563ffffffff1660e01b81526004016104ed94939291906118df565b60206040518083038186803b15801561050557600080fd5b505afa158015610519573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053d919061181b565b95945050505050565b61054e6109bf565b6105886040805160c08101909152806000815260200160008152602001600081526020016060815260200160608152602001606081525090565b6105ca88888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a429050565b84156106125761061286868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a539050565b821561062c5761062c6106258486611b6c565b8290610a9a565b61063a818362011170610ada565b6004555050505050505050565b61064f6109bf565b61065881610c98565b50565b8260045414610699576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610351565b6106a282610d8f565b6005558151600780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556106e481610d8f565b600655516007805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b6060610748604051806040016040528060608152602001600081525090565b61075481610100610e17565b5060408051808201909152600c81527f636f64654c6f636174696f6e00000000000000000000000000000000000000006020820152610794908290610e82565b82516107b29060018111156107ab576107ab611c48565b8290610e9e565b60408051808201909152600881527f6c616e677561676500000000000000000000000000000000000000000000000060208201526107f1908290610e82565b60408301516108089080156107ab576107ab611c48565b60408051808201909152600681527f736f7572636500000000000000000000000000000000000000000000000000006020820152610847908290610e82565b6060830151610857908290610e82565b60a083015151156109065760408051808201909152600481527f617267730000000000000000000000000000000000000000000000000000000060208201526108a1908290610e82565b6108aa81610ec4565b60005b8360a00151518110156108fc576108ea8460a0015182815181106108d3576108d3611c77565b602002602001015183610e8290919063ffffffff16565b806108f481611ba5565b9150506108ad565b5061090681610ecf565b608083015151156109b85760408051808201909152600f81527f736563726574734c6f636174696f6e00000000000000000000000000000000006020820152610950908290610e82565b610969836020015160018111156107ab576107ab611c48565b60408051808201909152600781527f736563726574730000000000000000000000000000000000000000000000000060208201526109a8908290610e82565b60808301516109b8908290610eda565b5192915050565b60025473ffffffffffffffffffffffffffffffffffffffff163314610a40576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610351565b565b610a4f8260008084610ee7565b5050565b8051610a8b576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006020830152608090910152565b8051610ad2576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60008054819073ffffffffffffffffffffffffffffffffffffffff166328242b0485610b0588610729565b866040518463ffffffff1660e01b8152600401610b24939291906118a6565b602060405180830381600087803b158015610b3e57600080fd5b505af1158015610b52573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b76919061153a565b905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635ab1bd536040518163ffffffff1660e01b815260040160206040518083038186803b158015610bde57600080fd5b505afa158015610bf2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c16919061151d565b60008281526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9490941693909317909255905182917f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db891a2949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610d18576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610351565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600080600060209050602084511015610da6575082515b60005b81811015610e0e57610dbc816008611b18565b858281518110610dce57610dce611c77565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c929092179180610e0681611ba5565b915050610da9565b50909392505050565b604080518082019091526060815260006020820152610e37602083611bde565b15610e5f57610e47602083611bde565b610e52906020611b55565b610e5c90836119d9565b91505b506020808301829052604080518085526000815283019091019052815b92915050565b610e8f8260038351610f7b565b610e99828261108a565b505050565b67ffffffffffffffff811115610eb857610a4f82826110b8565b610a4f82600083610f7b565b6106588160046110ef565b6106588160076110ef565b610e8f8260028351610f7b565b8051610f1f576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836001811115610f3257610f32611c48565b90816001811115610f4557610f45611c48565b90525060408401828015610f5b57610f5b611c48565b90818015610f6b57610f6b611c48565b9052506060909301929092525050565b60178167ffffffffffffffff1611610fa0576102ce8360e0600585901b168317611100565b60ff8167ffffffffffffffff1611610fde57610fc7836018611fe0600586901b1617611100565b506102ce8367ffffffffffffffff83166001611125565b61ffff8167ffffffffffffffff161161101d57611006836019611fe0600586901b1617611100565b506102ce8367ffffffffffffffff83166002611125565b63ffffffff8167ffffffffffffffff161161105e5761104783601a611fe0600586901b1617611100565b506102ce8367ffffffffffffffff83166004611125565b61107383601b611fe0600586901b1617611100565b506102ce8367ffffffffffffffff83166008611125565b6040805180820190915260608152600060208201526110b183846000015151848551611153565b9392505050565b6110c38260c2611100565b50610a4f82826040516020016110db91815260200190565b604051602081830303815290604052610eda565b610e9982601f611fe0600585901b16175b6040805180820190915260608152600060208201526110b1838460000151518461125b565b60408051808201909152606081526000602082015261114b8485600001515185856112b7565b949350505050565b604080518082019091526060815260006020820152825182111561117657600080fd5b602085015161118583866119d9565b11156111b8576111b8856111a8876020015187866111a391906119d9565b611338565b6111b3906002611b18565b61134f565b6000808651805187602083010193508088870111156111d75787860182525b505050602084015b6020841061121757805182526111f66020836119d9565b91506112036020826119d9565b9050611210602085611b55565b93506111df565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b604080518082019091526060815260006020820152836020015183106112905761129084856020015160026111b39190611b18565b8351805160208583010184815350808514156112ad576001810182525b5093949350505050565b60408051808201909152606081526000602082015260208501516112db85846119d9565b11156112ef576112ef856111a886856119d9565b600060016112ff84610100611a52565b6113099190611b55565b905085518386820101858319825116178152508051848701111561132d5783860181525b509495945050505050565b600081831115611349575081610e7c565b50919050565b815161135b8383610e17565b506102ce838261108a565b600067ffffffffffffffff8084111561138157611381611ca6565b8360051b6020611392818301611944565b8681529350808401858381018910156113aa57600080fd5b60009350835b888110156113e5578135868111156113c6578586fd5b6113d28b828b0161145b565b84525091830191908301906001016113b0565b5050505050509392505050565b600082601f83011261140357600080fd5b6110b183833560208501611366565b60008083601f84011261142457600080fd5b50813567ffffffffffffffff81111561143c57600080fd5b60208301915083602082850101111561145457600080fd5b9250929050565b600082601f83011261146c57600080fd5b813561147f61147a82611993565b611944565b81815284602083860101111561149457600080fd5b816020850160208301376000918101602001919091529392505050565b8035600181106114c057600080fd5b919050565b8035600281106114c057600080fd5b803563ffffffff811681146114c057600080fd5b803567ffffffffffffffff811681146114c057600080fd5b60006020828403121561151257600080fd5b81356110b181611cd5565b60006020828403121561152f57600080fd5b81516110b181611cd5565b60006020828403121561154c57600080fd5b5051919050565b60008060006060848603121561156857600080fd5b83359250602084013567ffffffffffffffff8082111561158757600080fd5b6115938783880161145b565b935060408601359150808211156115a957600080fd5b506115b68682870161145b565b9150509250925092565b6000602082840312156115d257600080fd5b815167ffffffffffffffff8111156115e957600080fd5b8201601f810184136115fa57600080fd5b805161160861147a82611993565b81815285602083850101111561161d57600080fd5b61053d826020830160208601611b79565b60008060008060008060006080888a03121561164957600080fd5b873567ffffffffffffffff8082111561166157600080fd5b61166d8b838c01611412565b909950975060208a013591508082111561168657600080fd5b6116928b838c01611412565b909750955060408a01359150808211156116ab57600080fd5b818a0191508a601f8301126116bf57600080fd5b8135818111156116ce57600080fd5b8b60208260051b85010111156116e357600080fd5b6020830195508094505050506116fb606089016114e8565b905092959891949750929550565b6000806000806080858703121561171f57600080fd5b843567ffffffffffffffff8082111561173757600080fd5b9086019060c0828903121561174b57600080fd5b61175361191b565b61175c836114c5565b815261176a602084016114c5565b602082015261177b604084016114b1565b604082015260608301358281111561179257600080fd5b61179e8a82860161145b565b6060830152506080830135828111156117b657600080fd5b6117c28a82860161145b565b60808301525060a0830135828111156117da57600080fd5b6117e68a8286016113f2565b60a08301525095506117fd915050602086016114e8565b925061180b604086016114d4565b9396929550929360600135925050565b60006020828403121561182d57600080fd5b81516bffffffffffffffffffffffff811681146110b157600080fd5b60008151808452611861816020860160208601611b79565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006110b16020830184611849565b67ffffffffffffffff841681526060602082015260006118c96060830185611849565b905063ffffffff83166040830152949350505050565b67ffffffffffffffff851681526080602082015260006119026080830186611849565b63ffffffff949094166040830152506060015292915050565b60405160c0810167ffffffffffffffff8111828210171561193e5761193e611ca6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561198b5761198b611ca6565b604052919050565b600067ffffffffffffffff8211156119ad576119ad611ca6565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082198211156119ec576119ec611c19565b500190565b600181815b80851115611a4a57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611a3057611a30611c19565b80851615611a3d57918102915b93841c93908002906119f6565b509250929050565b60006110b18383600082611a6857506001610e7c565b81611a7557506000610e7c565b8160018114611a8b5760028114611a9557611ab1565b6001915050610e7c565b60ff841115611aa657611aa6611c19565b50506001821b610e7c565b5060208310610133831016604e8410600b8410161715611ad4575081810a610e7c565b611ade83836119f1565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611b1057611b10611c19565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611b5057611b50611c19565b500290565b600082821015611b6757611b67611c19565b500390565b60006110b1368484611366565b60005b83811015611b94578181015183820152602001611b7c565b838111156102ce5750506000910152565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611bd757611bd7611c19565b5060010190565b600082611c14577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461065857600080fdfea164736f6c6343000806000a",
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
