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

type FunctionsRequest struct {
	CodeLocation    uint8
	SecretsLocation uint8
	Language        uint8
	Source          string
	Secrets         []byte
	Args            []string
}

var OCR2DRClientExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsNotPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderIsNotRegistry\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"SendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumFunctions.Location\",\"name\":\"codeLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.Location\",\"name\":\"secretsLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.CodeLanguage\",\"name\":\"language\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"}],\"internalType\":\"structFunctions.Request\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001f0338038062001f03833981016040819052620000349162000199565b600080546001600160a01b0319166001600160a01b038316178155339081906001600160a01b038216620000af5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600280546001600160a01b0319166001600160a01b0384811691909117909155811615620000e257620000e281620000ec565b50505050620001cb565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a6565b600380546001600160a01b0319166001600160a01b03838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600060208284031215620001ac57600080fd5b81516001600160a01b0381168114620001c457600080fd5b9392505050565b611d2880620001db6000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80638da5cb5b1161008c578063d4b3917511610066578063d4b39175146101aa578063d769717e146101da578063f2fde38b146101ed578063fc2a88c31461020057600080fd5b80638da5cb5b1461015d578063b48cffea14610185578063d328a91e1461019557600080fd5b806362747e42116100bd57806362747e42146101425780636d9809a01461014b57806379ba50971461015557600080fd5b80630ca76175146100e457806329f0de3f146100f95780632c29166b14610115575b600080fd5b6100f76100f2366004611577565b610209565b005b61010260065481565b6040519081526020015b60405180910390f35b60075461012d90640100000000900463ffffffff1681565b60405163ffffffff909116815260200161010c565b61010260055481565b61012d6201117081565b6100f76102d4565b60025460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b60075461012d9063ffffffff1681565b61019d6103da565b60405161010c91906118b7565b6101bd6101b836600461172d565b6104a3565b6040516bffffffffffffffffffffffff909116815260200161010c565b6100f76101e8366004611652565b610546565b6100f76101fb366004611524565b610647565b61010260045481565b600083815260016020526040902054839073ffffffffffffffffffffffffffffffffffffffff163314610268576040517fa0c5ec6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e691a26102ce84848461065b565b50505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461035a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560038054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b60008054604080517fd328a91e000000000000000000000000000000000000000000000000000000008152905160609373ffffffffffffffffffffffffffffffffffffffff9093169263d328a91e9260048082019391829003018186803b15801561044457600080fd5b505afa158015610458573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261049e91908101906115e4565b905090565b6000805473ffffffffffffffffffffffffffffffffffffffff1663d227d245856104cc88610729565b86866040518563ffffffff1660e01b81526004016104ed9493929190611903565b60206040518083038186803b15801561050557600080fd5b505afa158015610519573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053d919061183f565b95945050505050565b61054e6109f5565b6105886040805160c08101909152806000815260200160008152602001600081526020016060815260200160608152602001606081525090565b6105ca88888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a789050565b84156106125761061286868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a899050565b821561062c5761062c6106258486611b90565b8290610ad0565b61063a818362011170610b10565b6004555050505050505050565b61064f6109f5565b61065881610cce565b50565b8260045414610699576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610351565b6106a282610dc5565b6005558151600780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556106e481610dc5565b600655516007805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b6060610733611355565b805161074190610100610e4d565b50610781816040518060400160405280600c81526020017f636f64654c6f636174696f6e0000000000000000000000000000000000000000815250610ec7565b6107a0818460000151600281111561079b5761079b611c6c565b610ee5565b6107df816040518060400160405280600881526020017f6c616e6775616765000000000000000000000000000000000000000000000000815250610ec7565b6107f9818460400151600081111561079b5761079b611c6c565b610838816040518060400160405280600681526020017f736f757263650000000000000000000000000000000000000000000000000000815250610ec7565b610846818460600151610ec7565b60a083015151156108ec57610890816040518060400160405280600481526020017f6172677300000000000000000000000000000000000000000000000000000000815250610ec7565b61089981610f1e565b60005b8360a00151518110156108e2576108d0828560a0015183815181106108c3576108c3611c9b565b6020026020010151610ec7565b806108da81611bc9565b91505061089c565b506108ec81610f42565b608083015151156109ed5760008360200151600281111561090f5761090f611c6c565b1415610947576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610986816040518060400160405280600f81526020017f736563726574734c6f636174696f6e0000000000000000000000000000000000815250610ec7565b6109a0818460200151600281111561079b5761079b611c6c565b6109df816040518060400160405280600781526020017f7365637265747300000000000000000000000000000000000000000000000000815250610ec7565b6109ed818460800151610f60565b515192915050565b60025473ffffffffffffffffffffffffffffffffffffffff163314610a76576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610351565b565b610a858260008084610f6d565b5050565b8051610ac1576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b8051610b08576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60008054819073ffffffffffffffffffffffffffffffffffffffff166328242b0485610b3b88610729565b866040518463ffffffff1660e01b8152600401610b5a939291906118ca565b602060405180830381600087803b158015610b7457600080fd5b505af1158015610b88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bac919061155e565b905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635ab1bd536040518163ffffffff1660e01b815260040160206040518083038186803b158015610c1457600080fd5b505afa158015610c28573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c4c9190611541565b60008281526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9490941693909317909255905182917f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db891a2949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610d4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610351565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600080600060209050602084511015610ddc575082515b60005b81811015610e4457610df2816008611b3c565b858281518110610e0457610e04611c9b565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c929092179180610e3c81611bc9565b915050610ddf565b50909392505050565b604080518082019091526060815260006020820152610e6d602083611c02565b15610e9557610e7d602083611c02565b610e88906020611b79565b610e9290836119fd565b91505b602080840183905260405180855260008152908184010181811015610eb957600080fd5b604052508290505b92915050565b610ed48260038351611001565b8151610ee09082611122565b505050565b8151610ef29060c261114a565b50610a858282604051602001610f0a91815260200190565b604051602081830303815290604052610f60565b610f298160046111b3565b600181602001818151610f3c91906119fd565b90525050565b610f4d8160076111b3565b600181602001818151610f3c9190611b79565b610ed48260028351611001565b8051610fa5576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610fb857610fb8611c6c565b90816002811115610fcb57610fcb611c6c565b90525060408401828015610fe157610fe1611c6c565b90818015610ff157610ff1611c6c565b9052506060909301929092525050565b60178167ffffffffffffffff16116110285782516102ce9060e0600585901b16831761114a565b60ff8167ffffffffffffffff161161106a578251611051906018611fe0600586901b161761114a565b5082516102ce9067ffffffffffffffff831660016111ca565b61ffff8167ffffffffffffffff16116110ad578251611094906019611fe0600586901b161761114a565b5082516102ce9067ffffffffffffffff831660026111ca565b63ffffffff8167ffffffffffffffff16116110f25782516110d990601a611fe0600586901b161761114a565b5082516102ce9067ffffffffffffffff831660046111ca565b825161110990601b611fe0600586901b161761114a565b5082516102ce9067ffffffffffffffff831660086111ca565b6040805180820190915260608152600060208201526111438383845161124f565b9392505050565b604080518082019091526060815260006020820152825151600061116f8260016119fd565b905084602001518210611190576111908561118b836002611b3c565b61133e565b84516020838201018581535080518211156111a9578181525b5093949350505050565b8151610ee090601f611fe0600585901b161761114a565b60408051808201909152606081526000602082015283515160006111ee82856119fd565b9050856020015181111561120b5761120b8661118b836002611b3c565b6000600161121b86610100611a76565b6112259190611b79565b90508651828101878319825116178152508051831115611243578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561127257600080fd5b835151600061128184836119fd565b9050856020015181111561129e5761129e8661118b836002611b3c565b8551805183820160200191600091808511156112b8578482525b505050602086015b602086106112f857805182526112d76020836119fd565b91506112e46020826119fd565b90506112f1602087611b79565b95506112c0565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b815161134a8383610e4d565b506102ce8382611122565b604051806040016040528061137d604051806040016040528060608152602001600081525090565b8152602001600081525090565b600067ffffffffffffffff808411156113a5576113a5611cca565b8360051b60206113b6818301611968565b8681529350808401858381018910156113ce57600080fd5b60009350835b88811015611409578135868111156113ea578586fd5b6113f68b828b0161147f565b84525091830191908301906001016113d4565b5050505050509392505050565b600082601f83011261142757600080fd5b6111438383356020850161138a565b60008083601f84011261144857600080fd5b50813567ffffffffffffffff81111561146057600080fd5b60208301915083602082850101111561147857600080fd5b9250929050565b600082601f83011261149057600080fd5b81356114a361149e826119b7565b611968565b8181528460208386010111156114b857600080fd5b816020850160208301376000918101602001919091529392505050565b8035600181106114e457600080fd5b919050565b8035600381106114e457600080fd5b803563ffffffff811681146114e457600080fd5b803567ffffffffffffffff811681146114e457600080fd5b60006020828403121561153657600080fd5b813561114381611cf9565b60006020828403121561155357600080fd5b815161114381611cf9565b60006020828403121561157057600080fd5b5051919050565b60008060006060848603121561158c57600080fd5b83359250602084013567ffffffffffffffff808211156115ab57600080fd5b6115b78783880161147f565b935060408601359150808211156115cd57600080fd5b506115da8682870161147f565b9150509250925092565b6000602082840312156115f657600080fd5b815167ffffffffffffffff81111561160d57600080fd5b8201601f8101841361161e57600080fd5b805161162c61149e826119b7565b81815285602083850101111561164157600080fd5b61053d826020830160208601611b9d565b60008060008060008060006080888a03121561166d57600080fd5b873567ffffffffffffffff8082111561168557600080fd5b6116918b838c01611436565b909950975060208a01359150808211156116aa57600080fd5b6116b68b838c01611436565b909750955060408a01359150808211156116cf57600080fd5b818a0191508a601f8301126116e357600080fd5b8135818111156116f257600080fd5b8b60208260051b850101111561170757600080fd5b60208301955080945050505061171f6060890161150c565b905092959891949750929550565b6000806000806080858703121561174357600080fd5b843567ffffffffffffffff8082111561175b57600080fd5b9086019060c0828903121561176f57600080fd5b61177761193f565b611780836114e9565b815261178e602084016114e9565b602082015261179f604084016114d5565b60408201526060830135828111156117b657600080fd5b6117c28a82860161147f565b6060830152506080830135828111156117da57600080fd5b6117e68a82860161147f565b60808301525060a0830135828111156117fe57600080fd5b61180a8a828601611416565b60a08301525095506118219150506020860161150c565b925061182f604086016114f8565b9396929550929360600135925050565b60006020828403121561185157600080fd5b81516bffffffffffffffffffffffff8116811461114357600080fd5b60008151808452611885816020860160208601611b9d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611143602083018461186d565b67ffffffffffffffff841681526060602082015260006118ed606083018561186d565b905063ffffffff83166040830152949350505050565b67ffffffffffffffff85168152608060208201526000611926608083018661186d565b63ffffffff949094166040830152506060015292915050565b60405160c0810167ffffffffffffffff8111828210171561196257611962611cca565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156119af576119af611cca565b604052919050565b600067ffffffffffffffff8211156119d1576119d1611cca565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60008219821115611a1057611a10611c3d565b500190565b600181815b80851115611a6e57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611a5457611a54611c3d565b80851615611a6157918102915b93841c9390800290611a1a565b509250929050565b60006111438383600082611a8c57506001610ec1565b81611a9957506000610ec1565b8160018114611aaf5760028114611ab957611ad5565b6001915050610ec1565b60ff841115611aca57611aca611c3d565b50506001821b610ec1565b5060208310610133831016604e8410600b8410161715611af8575081810a610ec1565b611b028383611a15565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611b3457611b34611c3d565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611b7457611b74611c3d565b500290565b600082821015611b8b57611b8b611c3d565b500390565b600061114336848461138a565b60005b83811015611bb8578181015183820152602001611ba0565b838111156102ce5750506000910152565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611bfb57611bfb611c3d565b5060010190565b600082611c38577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461065857600080fdfea164736f6c6343000806000a",
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
	parsed, err := OCR2DRClientExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
