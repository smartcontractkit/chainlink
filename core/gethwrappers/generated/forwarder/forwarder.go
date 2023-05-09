// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package forwarder

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

type IForwarderForwardRequest struct {
	From           common.Address
	Target         common.Address
	Nonce          *big.Int
	Data           []byte
	ValidUntilTime *big.Int
}

var ForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"ForwardFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"domainSeparator\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"domainValue\",\"type\":\"bytes\"}],\"name\":\"DomainRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"domainSeparator\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnValue\",\"type\":\"bytes\"}],\"name\":\"ForwardSucceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"typeHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"typeStr\",\"type\":\"string\"}],\"name\":\"RequestTypeRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EIP712_DOMAIN_TYPE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GENERIC_PARAMS\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validUntilTime\",\"type\":\"uint256\"}],\"internalType\":\"structIForwarder.ForwardRequest\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"requestTypeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"suffixData\",\"type\":\"bytes\"}],\"name\":\"_getEncoded\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"domains\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validUntilTime\",\"type\":\"uint256\"}],\"internalType\":\"structIForwarder.ForwardRequest\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"domainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"requestTypeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"suffixData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"ret\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getDomainSeparator\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"registerDomainSeparator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"typeName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"typeSuffix\",\"type\":\"string\"}],\"name\":\"registerRequestType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"typeHashes\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validUntilTime\",\"type\":\"uint256\"}],\"internalType\":\"structIForwarder.ForwardRequest\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"domainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"requestTypeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"suffixData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000f1565b50505060006040518060800160405280604b815260200162001eb0604b9139604051602001620000cd919062000234565b60408051601f198184030181529190529050620000ea816200019c565b50620002ad565b336001600160a01b038216036200014b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b805160208083019190912060008181526002909252604091829020805460ff19166001179055905181907f64d6bce64323458c44643c51fe45113efc882082f7b7fd5f09f0d69d2eedb20290620001f590859062000278565b60405180910390a25050565b60005b838110156200021e57818101518382015260200162000204565b838111156200022e576000848401525b50505050565b6e08cdee4eec2e4c8a4cae2eacae6e85608b1b8152600082516200026081600f85016020870162000201565b602960f81b600f939091019283015250601001919050565b60208152600082518060208401526200029981604085016020870162000201565b601f01601f19169190910160400192915050565b611bf380620002bd6000396000f3fe6080604052600436106100ec5760003560e01c80638da5cb5b1161008a578063c3f28abd11610059578063c3f28abd146102bd578063c722f177146102d2578063d9210be514610302578063f2fde38b1461032257600080fd5b80638da5cb5b146102285780639a0e89291461025d5780639c7b45921461027d578063b5aa37671461029d57600080fd5b80632d0335ab116100c65780632d0335ab1461017f57806371714e17146101d057806379ba5097146101f05780637c9547901461020757600080fd5b806301ffc9a7146100f8578063066a310c1461012d57806321fe98df1461014f57600080fd5b366100f357005b600080fd5b34801561010457600080fd5b50610118610113366004611585565b610342565b60405190151581526020015b60405180910390f35b34801561013957600080fd5b506101426103db565b6040516101249190611648565b34801561015b57600080fd5b5061011861016a36600461165b565b60026020526000908152604090205460ff1681565b34801561018b57600080fd5b506101c261019a366004611674565b73ffffffffffffffffffffffffffffffffffffffff1660009081526004602052604090205490565b604051908152602001610124565b3480156101dc57600080fd5b506101426101eb366004611704565b6103f7565b3480156101fc57600080fd5b5061020561049f565b005b61021a610215366004611774565b6105a1565b60405161012492919061181c565b34801561023457600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610124565b34801561026957600080fd5b50610205610278366004611774565b610855565b34801561028957600080fd5b5061020561029836600461183f565b610876565b3480156102a957600080fd5b506101426102b836600461183f565b610915565b3480156102c957600080fd5b506101426109ae565b3480156102de57600080fd5b506101186102ed36600461165b565b60036020526000908152604090205460ff1681565b34801561030e57600080fd5b5061020561031d36600461183f565b6109ca565b34801561032e57600080fd5b5061020561033d366004611674565b610b52565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f8ec2b5650000000000000000000000000000000000000000000000000000000014806103d557507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6040518060800160405280604b8152602001611b4a604b913981565b6060836104076020870187611674565b73ffffffffffffffffffffffffffffffffffffffff1661042d6040880160208901611674565b73ffffffffffffffffffffffffffffffffffffffff16604088013561045560608a018a611892565b6040516104639291906118f7565b604051908190038120610486959493929160808c0135908a908a90602001611907565b6040516020818303038152906040529050949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610525576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600060606105b489898989898989610b66565b6105bd89610dd9565b608089013515806105d15750428960800135115b610637576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4657443a20726571756573742065787069726564000000000000000000000000604482015260640161051c565b600061064660608b018b611892565b61065360208d018d611674565b60405160200161066593929190611944565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181529091506106a5908b0160208c01611674565b73ffffffffffffffffffffffffffffffffffffffff16816040516106c9919061197d565b6000604051808303816000865af19150503d8060008114610706576040519150601f19603f3d011682016040523d82523d6000602084013e61070b565b606091505b509093509150826107af5781516000036107a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f466f727761726465642063616c6c20726576657274656420776974686f75742060448201527f726561736f6e0000000000000000000000000000000000000000000000000000606482015260840161051c565b815182602001fd5b886107c060408c0160208d01611674565b73ffffffffffffffffffffffffffffffffffffffff166107e360208d018d611674565b73ffffffffffffffffffffffffffffffffffffffff167f94e72b9bdade69580d865c1293804239e91513c10e675ca3345fc848f9daba2e8d604001358e806060019061082f9190611892565b886040516108409493929190611999565b60405180910390a45097509795505050505050565b61085e87610e8e565b61086d87878787878787610b66565b50505050505050565b61087e610f46565b600061088c85858585610915565b80516020808301919091206000818152600390925260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590519192509081907f4bc68689cbe89a4a6333a3ab0a70093874da3e5bfb71e93102027f3f073687d890610905908590611648565b60405180910390a2505050505050565b60606000469050604051806080016040528060528152602001611b956052913980519060200120868660405161094c9291906118f7565b604051809103902085856040516109649291906118f7565b6040805191829003822060208301949094528101919091526060810191909152608081018290523060a082015260c001604051602081830303815290604052915050949350505050565b604051806080016040528060528152602001611b956052913981565b60005b83811015610afd5760008585838181106109e9576109e9611a06565b909101357fff00000000000000000000000000000000000000000000000000000000000000169150507f28000000000000000000000000000000000000000000000000000000000000008114801590610a8457507f29000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821614155b610aea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f4657443a20696e76616c696420747970656e616d650000000000000000000000604482015260640161051c565b5080610af581611a35565b9150506109cd565b50600084846040518060800160405280604b8152602001611b4a604b91398585604051602001610b31959493929190611a94565b6040516020818303038152906040529050610b4b81610fc9565b5050505050565b610b5a610f46565b610b638161104a565b50565b60008681526003602052604090205460ff16610bde576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f4657443a20756e7265676973746572656420646f6d61696e207365702e000000604482015260640161051c565b60008581526002602052604090205460ff16610c56576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4657443a20756e72656769737465726564207479706568617368000000000000604482015260640161051c565b600086610c65898888886103f7565b8051602091820120604051610cac9392017f190100000000000000000000000000000000000000000000000000000000000081526002810192909252602282015260420190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209050321580610d695750610cf96020890189611674565b73ffffffffffffffffffffffffffffffffffffffff16610d5184848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869392505061113f9050565b73ffffffffffffffffffffffffffffffffffffffff16145b610dcf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4657443a207369676e6174757265206d69736d61746368000000000000000000604482015260640161051c565b5050505050505050565b604081013560046000610def6020850185611674565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000908120805491610e2383611a35565b9190505514610b63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4657443a206e6f6e6365206d69736d6174636800000000000000000000000000604482015260640161051c565b604081013560046000610ea46020850185611674565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414610b63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4657443a206e6f6e6365206d69736d6174636800000000000000000000000000604482015260640161051c565b60005473ffffffffffffffffffffffffffffffffffffffff163314610fc7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161051c565b565b80516020808301919091206000818152600290925260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055905181907f64d6bce64323458c44643c51fe45113efc882082f7b7fd5f09f0d69d2eedb2029061103e908590611648565b60405180910390a25050565b3373ffffffffffffffffffffffffffffffffffffffff8216036110c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161051c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080600061114e8585611163565b9150915061115b816111d1565b509392505050565b60008082516041036111995760208301516040840151606085015160001a61118d87828585611425565b945094505050506111ca565b82516040036111c257602083015160408401516111b786838361153d565b9350935050506111ca565b506000905060025b9250929050565b60008160048111156111e5576111e5611b1a565b036111ed5750565b600181600481111561120157611201611b1a565b03611268576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604482015260640161051c565b600281600481111561127c5761127c611b1a565b036112e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604482015260640161051c565b60038160048111156112f7576112f7611b1a565b03611384576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c60448201527f7565000000000000000000000000000000000000000000000000000000000000606482015260840161051c565b600481600481111561139857611398611b1a565b03610b63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c60448201527f7565000000000000000000000000000000000000000000000000000000000000606482015260840161051c565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a083111561145c5750600090506003611534565b8460ff16601b1415801561147457508460ff16601c14155b156114855750600090506004611534565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa1580156114d9573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811661152d57600060019250925050611534565b9150600090505b94509492505050565b6000807f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831660ff84901c601b0161157787828885611425565b935093505050935093915050565b60006020828403121561159757600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146115c757600080fd5b9392505050565b60005b838110156115e95781810151838201526020016115d1565b838111156115f8576000848401525b50505050565b600081518084526116168160208601602086016115ce565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006115c760208301846115fe565b60006020828403121561166d57600080fd5b5035919050565b60006020828403121561168657600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146115c757600080fd5b600060a082840312156116bc57600080fd5b50919050565b60008083601f8401126116d457600080fd5b50813567ffffffffffffffff8111156116ec57600080fd5b6020830191508360208285010111156111ca57600080fd5b6000806000806060858703121561171a57600080fd5b843567ffffffffffffffff8082111561173257600080fd5b61173e888389016116aa565b955060208701359450604087013591508082111561175b57600080fd5b50611768878288016116c2565b95989497509550505050565b600080600080600080600060a0888a03121561178f57600080fd5b873567ffffffffffffffff808211156117a757600080fd5b6117b38b838c016116aa565b985060208a0135975060408a0135965060608a01359150808211156117d757600080fd5b6117e38b838c016116c2565b909650945060808a01359150808211156117fc57600080fd5b506118098a828b016116c2565b989b979a50959850939692959293505050565b821515815260406020820152600061183760408301846115fe565b949350505050565b6000806000806040858703121561185557600080fd5b843567ffffffffffffffff8082111561186d57600080fd5b611879888389016116c2565b9096509450602087013591508082111561175b57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126118c757600080fd5b83018035915067ffffffffffffffff8211156118e257600080fd5b6020019150368190038213156111ca57600080fd5b8183823760009101908152919050565b8881528760208201528660408201528560608201528460808201528360a0820152818360c08301376000910160c001908152979650505050505050565b8284823760609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000169101908152601401919050565b6000825161198f8184602087016115ce565b9190910192915050565b848152606060208201528260608201528284608083013760006080848301015260007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f850116820160808382030160408401526119fb60808201856115fe565b979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611a8d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b8486823760008582017f280000000000000000000000000000000000000000000000000000000000000081528551611ad3816001840160208a016115ce565b7f2c00000000000000000000000000000000000000000000000000000000000000600192909101918201528385600283013760009301600201928352509095945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fdfe616464726573732066726f6d2c61646472657373207461726765742c75696e74323536206e6f6e63652c627974657320646174612c75696e743235362076616c6964556e74696c54696d65454950373132446f6d61696e28737472696e67206e616d652c737472696e672076657273696f6e2c75696e7432353620636861696e49642c6164647265737320766572696679696e67436f6e747261637429a164736f6c634300080f000a616464726573732066726f6d2c61646472657373207461726765742c75696e74323536206e6f6e63652c627974657320646174612c75696e743235362076616c6964556e74696c54696d65",
}

var ForwarderABI = ForwarderMetaData.ABI

var ForwarderBin = ForwarderMetaData.Bin

func DeployForwarder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Forwarder, error) {
	parsed, err := ForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ForwarderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Forwarder{ForwarderCaller: ForwarderCaller{contract: contract}, ForwarderTransactor: ForwarderTransactor{contract: contract}, ForwarderFilterer: ForwarderFilterer{contract: contract}}, nil
}

type Forwarder struct {
	address common.Address
	abi     abi.ABI
	ForwarderCaller
	ForwarderTransactor
	ForwarderFilterer
}

type ForwarderCaller struct {
	contract *bind.BoundContract
}

type ForwarderTransactor struct {
	contract *bind.BoundContract
}

type ForwarderFilterer struct {
	contract *bind.BoundContract
}

type ForwarderSession struct {
	Contract     *Forwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ForwarderCallerSession struct {
	Contract *ForwarderCaller
	CallOpts bind.CallOpts
}

type ForwarderTransactorSession struct {
	Contract     *ForwarderTransactor
	TransactOpts bind.TransactOpts
}

type ForwarderRaw struct {
	Contract *Forwarder
}

type ForwarderCallerRaw struct {
	Contract *ForwarderCaller
}

type ForwarderTransactorRaw struct {
	Contract *ForwarderTransactor
}

func NewForwarder(address common.Address, backend bind.ContractBackend) (*Forwarder, error) {
	abi, err := abi.JSON(strings.NewReader(ForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Forwarder{address: address, abi: abi, ForwarderCaller: ForwarderCaller{contract: contract}, ForwarderTransactor: ForwarderTransactor{contract: contract}, ForwarderFilterer: ForwarderFilterer{contract: contract}}, nil
}

func NewForwarderCaller(address common.Address, caller bind.ContractCaller) (*ForwarderCaller, error) {
	contract, err := bindForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ForwarderCaller{contract: contract}, nil
}

func NewForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*ForwarderTransactor, error) {
	contract, err := bindForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ForwarderTransactor{contract: contract}, nil
}

func NewForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*ForwarderFilterer, error) {
	contract, err := bindForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ForwarderFilterer{contract: contract}, nil
}

func bindForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Forwarder *ForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Forwarder.Contract.ForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_Forwarder *ForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.Contract.ForwarderTransactor.contract.Transfer(opts)
}

func (_Forwarder *ForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Forwarder.Contract.ForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_Forwarder *ForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Forwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_Forwarder *ForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.Contract.contract.Transfer(opts)
}

func (_Forwarder *ForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Forwarder.Contract.contract.Transact(opts, method, params...)
}

func (_Forwarder *ForwarderCaller) EIP712DOMAINTYPE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "EIP712_DOMAIN_TYPE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Forwarder *ForwarderSession) EIP712DOMAINTYPE() (string, error) {
	return _Forwarder.Contract.EIP712DOMAINTYPE(&_Forwarder.CallOpts)
}

func (_Forwarder *ForwarderCallerSession) EIP712DOMAINTYPE() (string, error) {
	return _Forwarder.Contract.EIP712DOMAINTYPE(&_Forwarder.CallOpts)
}

func (_Forwarder *ForwarderCaller) GENERICPARAMS(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "GENERIC_PARAMS")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Forwarder *ForwarderSession) GENERICPARAMS() (string, error) {
	return _Forwarder.Contract.GENERICPARAMS(&_Forwarder.CallOpts)
}

func (_Forwarder *ForwarderCallerSession) GENERICPARAMS() (string, error) {
	return _Forwarder.Contract.GENERICPARAMS(&_Forwarder.CallOpts)
}

func (_Forwarder *ForwarderCaller) GetEncoded(opts *bind.CallOpts, req IForwarderForwardRequest, requestTypeHash [32]byte, suffixData []byte) ([]byte, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "_getEncoded", req, requestTypeHash, suffixData)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_Forwarder *ForwarderSession) GetEncoded(req IForwarderForwardRequest, requestTypeHash [32]byte, suffixData []byte) ([]byte, error) {
	return _Forwarder.Contract.GetEncoded(&_Forwarder.CallOpts, req, requestTypeHash, suffixData)
}

func (_Forwarder *ForwarderCallerSession) GetEncoded(req IForwarderForwardRequest, requestTypeHash [32]byte, suffixData []byte) ([]byte, error) {
	return _Forwarder.Contract.GetEncoded(&_Forwarder.CallOpts, req, requestTypeHash, suffixData)
}

func (_Forwarder *ForwarderCaller) Domains(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "domains", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Forwarder *ForwarderSession) Domains(arg0 [32]byte) (bool, error) {
	return _Forwarder.Contract.Domains(&_Forwarder.CallOpts, arg0)
}

func (_Forwarder *ForwarderCallerSession) Domains(arg0 [32]byte) (bool, error) {
	return _Forwarder.Contract.Domains(&_Forwarder.CallOpts, arg0)
}

func (_Forwarder *ForwarderCaller) GetDomainSeparator(opts *bind.CallOpts, name string, version string) ([]byte, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "getDomainSeparator", name, version)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_Forwarder *ForwarderSession) GetDomainSeparator(name string, version string) ([]byte, error) {
	return _Forwarder.Contract.GetDomainSeparator(&_Forwarder.CallOpts, name, version)
}

func (_Forwarder *ForwarderCallerSession) GetDomainSeparator(name string, version string) ([]byte, error) {
	return _Forwarder.Contract.GetDomainSeparator(&_Forwarder.CallOpts, name, version)
}

func (_Forwarder *ForwarderCaller) GetNonce(opts *bind.CallOpts, from common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "getNonce", from)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Forwarder *ForwarderSession) GetNonce(from common.Address) (*big.Int, error) {
	return _Forwarder.Contract.GetNonce(&_Forwarder.CallOpts, from)
}

func (_Forwarder *ForwarderCallerSession) GetNonce(from common.Address) (*big.Int, error) {
	return _Forwarder.Contract.GetNonce(&_Forwarder.CallOpts, from)
}

func (_Forwarder *ForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Forwarder *ForwarderSession) Owner() (common.Address, error) {
	return _Forwarder.Contract.Owner(&_Forwarder.CallOpts)
}

func (_Forwarder *ForwarderCallerSession) Owner() (common.Address, error) {
	return _Forwarder.Contract.Owner(&_Forwarder.CallOpts)
}

func (_Forwarder *ForwarderCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Forwarder *ForwarderSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Forwarder.Contract.SupportsInterface(&_Forwarder.CallOpts, interfaceId)
}

func (_Forwarder *ForwarderCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Forwarder.Contract.SupportsInterface(&_Forwarder.CallOpts, interfaceId)
}

func (_Forwarder *ForwarderCaller) TypeHashes(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "typeHashes", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Forwarder *ForwarderSession) TypeHashes(arg0 [32]byte) (bool, error) {
	return _Forwarder.Contract.TypeHashes(&_Forwarder.CallOpts, arg0)
}

func (_Forwarder *ForwarderCallerSession) TypeHashes(arg0 [32]byte) (bool, error) {
	return _Forwarder.Contract.TypeHashes(&_Forwarder.CallOpts, arg0)
}

func (_Forwarder *ForwarderCaller) Verify(opts *bind.CallOpts, req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) error {
	var out []interface{}
	err := _Forwarder.contract.Call(opts, &out, "verify", req, domainSeparator, requestTypeHash, suffixData, sig)

	if err != nil {
		return err
	}

	return err

}

func (_Forwarder *ForwarderSession) Verify(req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) error {
	return _Forwarder.Contract.Verify(&_Forwarder.CallOpts, req, domainSeparator, requestTypeHash, suffixData, sig)
}

func (_Forwarder *ForwarderCallerSession) Verify(req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) error {
	return _Forwarder.Contract.Verify(&_Forwarder.CallOpts, req, domainSeparator, requestTypeHash, suffixData, sig)
}

func (_Forwarder *ForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.contract.Transact(opts, "acceptOwnership")
}

func (_Forwarder *ForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _Forwarder.Contract.AcceptOwnership(&_Forwarder.TransactOpts)
}

func (_Forwarder *ForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Forwarder.Contract.AcceptOwnership(&_Forwarder.TransactOpts)
}

func (_Forwarder *ForwarderTransactor) Execute(opts *bind.TransactOpts, req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) (*types.Transaction, error) {
	return _Forwarder.contract.Transact(opts, "execute", req, domainSeparator, requestTypeHash, suffixData, sig)
}

func (_Forwarder *ForwarderSession) Execute(req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) (*types.Transaction, error) {
	return _Forwarder.Contract.Execute(&_Forwarder.TransactOpts, req, domainSeparator, requestTypeHash, suffixData, sig)
}

func (_Forwarder *ForwarderTransactorSession) Execute(req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) (*types.Transaction, error) {
	return _Forwarder.Contract.Execute(&_Forwarder.TransactOpts, req, domainSeparator, requestTypeHash, suffixData, sig)
}

func (_Forwarder *ForwarderTransactor) RegisterDomainSeparator(opts *bind.TransactOpts, name string, version string) (*types.Transaction, error) {
	return _Forwarder.contract.Transact(opts, "registerDomainSeparator", name, version)
}

func (_Forwarder *ForwarderSession) RegisterDomainSeparator(name string, version string) (*types.Transaction, error) {
	return _Forwarder.Contract.RegisterDomainSeparator(&_Forwarder.TransactOpts, name, version)
}

func (_Forwarder *ForwarderTransactorSession) RegisterDomainSeparator(name string, version string) (*types.Transaction, error) {
	return _Forwarder.Contract.RegisterDomainSeparator(&_Forwarder.TransactOpts, name, version)
}

func (_Forwarder *ForwarderTransactor) RegisterRequestType(opts *bind.TransactOpts, typeName string, typeSuffix string) (*types.Transaction, error) {
	return _Forwarder.contract.Transact(opts, "registerRequestType", typeName, typeSuffix)
}

func (_Forwarder *ForwarderSession) RegisterRequestType(typeName string, typeSuffix string) (*types.Transaction, error) {
	return _Forwarder.Contract.RegisterRequestType(&_Forwarder.TransactOpts, typeName, typeSuffix)
}

func (_Forwarder *ForwarderTransactorSession) RegisterRequestType(typeName string, typeSuffix string) (*types.Transaction, error) {
	return _Forwarder.Contract.RegisterRequestType(&_Forwarder.TransactOpts, typeName, typeSuffix)
}

func (_Forwarder *ForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Forwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_Forwarder *ForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Forwarder.Contract.TransferOwnership(&_Forwarder.TransactOpts, to)
}

func (_Forwarder *ForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Forwarder.Contract.TransferOwnership(&_Forwarder.TransactOpts, to)
}

func (_Forwarder *ForwarderTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.contract.RawTransact(opts, nil)
}

func (_Forwarder *ForwarderSession) Receive() (*types.Transaction, error) {
	return _Forwarder.Contract.Receive(&_Forwarder.TransactOpts)
}

func (_Forwarder *ForwarderTransactorSession) Receive() (*types.Transaction, error) {
	return _Forwarder.Contract.Receive(&_Forwarder.TransactOpts)
}

type ForwarderDomainRegisteredIterator struct {
	Event *ForwarderDomainRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ForwarderDomainRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ForwarderDomainRegistered)
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
		it.Event = new(ForwarderDomainRegistered)
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

func (it *ForwarderDomainRegisteredIterator) Error() error {
	return it.fail
}

func (it *ForwarderDomainRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ForwarderDomainRegistered struct {
	DomainSeparator [32]byte
	DomainValue     []byte
	Raw             types.Log
}

func (_Forwarder *ForwarderFilterer) FilterDomainRegistered(opts *bind.FilterOpts, domainSeparator [][32]byte) (*ForwarderDomainRegisteredIterator, error) {

	var domainSeparatorRule []interface{}
	for _, domainSeparatorItem := range domainSeparator {
		domainSeparatorRule = append(domainSeparatorRule, domainSeparatorItem)
	}

	logs, sub, err := _Forwarder.contract.FilterLogs(opts, "DomainRegistered", domainSeparatorRule)
	if err != nil {
		return nil, err
	}
	return &ForwarderDomainRegisteredIterator{contract: _Forwarder.contract, event: "DomainRegistered", logs: logs, sub: sub}, nil
}

func (_Forwarder *ForwarderFilterer) WatchDomainRegistered(opts *bind.WatchOpts, sink chan<- *ForwarderDomainRegistered, domainSeparator [][32]byte) (event.Subscription, error) {

	var domainSeparatorRule []interface{}
	for _, domainSeparatorItem := range domainSeparator {
		domainSeparatorRule = append(domainSeparatorRule, domainSeparatorItem)
	}

	logs, sub, err := _Forwarder.contract.WatchLogs(opts, "DomainRegistered", domainSeparatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ForwarderDomainRegistered)
				if err := _Forwarder.contract.UnpackLog(event, "DomainRegistered", log); err != nil {
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

func (_Forwarder *ForwarderFilterer) ParseDomainRegistered(log types.Log) (*ForwarderDomainRegistered, error) {
	event := new(ForwarderDomainRegistered)
	if err := _Forwarder.contract.UnpackLog(event, "DomainRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ForwarderForwardSucceededIterator struct {
	Event *ForwarderForwardSucceeded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ForwarderForwardSucceededIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ForwarderForwardSucceeded)
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
		it.Event = new(ForwarderForwardSucceeded)
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

func (it *ForwarderForwardSucceededIterator) Error() error {
	return it.fail
}

func (it *ForwarderForwardSucceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ForwarderForwardSucceeded struct {
	From            common.Address
	Target          common.Address
	DomainSeparator [32]byte
	Nonce           *big.Int
	Data            []byte
	ReturnValue     []byte
	Raw             types.Log
}

func (_Forwarder *ForwarderFilterer) FilterForwardSucceeded(opts *bind.FilterOpts, from []common.Address, target []common.Address, domainSeparator [][32]byte) (*ForwarderForwardSucceededIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}
	var domainSeparatorRule []interface{}
	for _, domainSeparatorItem := range domainSeparator {
		domainSeparatorRule = append(domainSeparatorRule, domainSeparatorItem)
	}

	logs, sub, err := _Forwarder.contract.FilterLogs(opts, "ForwardSucceeded", fromRule, targetRule, domainSeparatorRule)
	if err != nil {
		return nil, err
	}
	return &ForwarderForwardSucceededIterator{contract: _Forwarder.contract, event: "ForwardSucceeded", logs: logs, sub: sub}, nil
}

func (_Forwarder *ForwarderFilterer) WatchForwardSucceeded(opts *bind.WatchOpts, sink chan<- *ForwarderForwardSucceeded, from []common.Address, target []common.Address, domainSeparator [][32]byte) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}
	var domainSeparatorRule []interface{}
	for _, domainSeparatorItem := range domainSeparator {
		domainSeparatorRule = append(domainSeparatorRule, domainSeparatorItem)
	}

	logs, sub, err := _Forwarder.contract.WatchLogs(opts, "ForwardSucceeded", fromRule, targetRule, domainSeparatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ForwarderForwardSucceeded)
				if err := _Forwarder.contract.UnpackLog(event, "ForwardSucceeded", log); err != nil {
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

func (_Forwarder *ForwarderFilterer) ParseForwardSucceeded(log types.Log) (*ForwarderForwardSucceeded, error) {
	event := new(ForwarderForwardSucceeded)
	if err := _Forwarder.contract.UnpackLog(event, "ForwardSucceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ForwarderOwnershipTransferRequestedIterator struct {
	Event *ForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ForwarderOwnershipTransferRequested)
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
		it.Event = new(ForwarderOwnershipTransferRequested)
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

func (it *ForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Forwarder *ForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Forwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ForwarderOwnershipTransferRequestedIterator{contract: _Forwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_Forwarder *ForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Forwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ForwarderOwnershipTransferRequested)
				if err := _Forwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_Forwarder *ForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*ForwarderOwnershipTransferRequested, error) {
	event := new(ForwarderOwnershipTransferRequested)
	if err := _Forwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ForwarderOwnershipTransferredIterator struct {
	Event *ForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ForwarderOwnershipTransferred)
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
		it.Event = new(ForwarderOwnershipTransferred)
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

func (it *ForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Forwarder *ForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Forwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ForwarderOwnershipTransferredIterator{contract: _Forwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Forwarder *ForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Forwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ForwarderOwnershipTransferred)
				if err := _Forwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Forwarder *ForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*ForwarderOwnershipTransferred, error) {
	event := new(ForwarderOwnershipTransferred)
	if err := _Forwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ForwarderRequestTypeRegisteredIterator struct {
	Event *ForwarderRequestTypeRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ForwarderRequestTypeRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ForwarderRequestTypeRegistered)
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
		it.Event = new(ForwarderRequestTypeRegistered)
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

func (it *ForwarderRequestTypeRegisteredIterator) Error() error {
	return it.fail
}

func (it *ForwarderRequestTypeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ForwarderRequestTypeRegistered struct {
	TypeHash [32]byte
	TypeStr  string
	Raw      types.Log
}

func (_Forwarder *ForwarderFilterer) FilterRequestTypeRegistered(opts *bind.FilterOpts, typeHash [][32]byte) (*ForwarderRequestTypeRegisteredIterator, error) {

	var typeHashRule []interface{}
	for _, typeHashItem := range typeHash {
		typeHashRule = append(typeHashRule, typeHashItem)
	}

	logs, sub, err := _Forwarder.contract.FilterLogs(opts, "RequestTypeRegistered", typeHashRule)
	if err != nil {
		return nil, err
	}
	return &ForwarderRequestTypeRegisteredIterator{contract: _Forwarder.contract, event: "RequestTypeRegistered", logs: logs, sub: sub}, nil
}

func (_Forwarder *ForwarderFilterer) WatchRequestTypeRegistered(opts *bind.WatchOpts, sink chan<- *ForwarderRequestTypeRegistered, typeHash [][32]byte) (event.Subscription, error) {

	var typeHashRule []interface{}
	for _, typeHashItem := range typeHash {
		typeHashRule = append(typeHashRule, typeHashItem)
	}

	logs, sub, err := _Forwarder.contract.WatchLogs(opts, "RequestTypeRegistered", typeHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ForwarderRequestTypeRegistered)
				if err := _Forwarder.contract.UnpackLog(event, "RequestTypeRegistered", log); err != nil {
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

func (_Forwarder *ForwarderFilterer) ParseRequestTypeRegistered(log types.Log) (*ForwarderRequestTypeRegistered, error) {
	event := new(ForwarderRequestTypeRegistered)
	if err := _Forwarder.contract.UnpackLog(event, "RequestTypeRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Forwarder *Forwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Forwarder.abi.Events["DomainRegistered"].ID:
		return _Forwarder.ParseDomainRegistered(log)
	case _Forwarder.abi.Events["ForwardSucceeded"].ID:
		return _Forwarder.ParseForwardSucceeded(log)
	case _Forwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _Forwarder.ParseOwnershipTransferRequested(log)
	case _Forwarder.abi.Events["OwnershipTransferred"].ID:
		return _Forwarder.ParseOwnershipTransferred(log)
	case _Forwarder.abi.Events["RequestTypeRegistered"].ID:
		return _Forwarder.ParseRequestTypeRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ForwarderDomainRegistered) Topic() common.Hash {
	return common.HexToHash("0x4bc68689cbe89a4a6333a3ab0a70093874da3e5bfb71e93102027f3f073687d8")
}

func (ForwarderForwardSucceeded) Topic() common.Hash {
	return common.HexToHash("0x94e72b9bdade69580d865c1293804239e91513c10e675ca3345fc848f9daba2e")
}

func (ForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (ForwarderRequestTypeRegistered) Topic() common.Hash {
	return common.HexToHash("0x64d6bce64323458c44643c51fe45113efc882082f7b7fd5f09f0d69d2eedb202")
}

func (_Forwarder *Forwarder) Address() common.Address {
	return _Forwarder.address
}

type ForwarderInterface interface {
	EIP712DOMAINTYPE(opts *bind.CallOpts) (string, error)

	GENERICPARAMS(opts *bind.CallOpts) (string, error)

	GetEncoded(opts *bind.CallOpts, req IForwarderForwardRequest, requestTypeHash [32]byte, suffixData []byte) ([]byte, error)

	Domains(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	GetDomainSeparator(opts *bind.CallOpts, name string, version string) ([]byte, error)

	GetNonce(opts *bind.CallOpts, from common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeHashes(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	Verify(opts *bind.CallOpts, req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) error

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, req IForwarderForwardRequest, domainSeparator [32]byte, requestTypeHash [32]byte, suffixData []byte, sig []byte) (*types.Transaction, error)

	RegisterDomainSeparator(opts *bind.TransactOpts, name string, version string) (*types.Transaction, error)

	RegisterRequestType(opts *bind.TransactOpts, typeName string, typeSuffix string) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterDomainRegistered(opts *bind.FilterOpts, domainSeparator [][32]byte) (*ForwarderDomainRegisteredIterator, error)

	WatchDomainRegistered(opts *bind.WatchOpts, sink chan<- *ForwarderDomainRegistered, domainSeparator [][32]byte) (event.Subscription, error)

	ParseDomainRegistered(log types.Log) (*ForwarderDomainRegistered, error)

	FilterForwardSucceeded(opts *bind.FilterOpts, from []common.Address, target []common.Address, domainSeparator [][32]byte) (*ForwarderForwardSucceededIterator, error)

	WatchForwardSucceeded(opts *bind.WatchOpts, sink chan<- *ForwarderForwardSucceeded, from []common.Address, target []common.Address, domainSeparator [][32]byte) (event.Subscription, error)

	ParseForwardSucceeded(log types.Log) (*ForwarderForwardSucceeded, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ForwarderOwnershipTransferred, error)

	FilterRequestTypeRegistered(opts *bind.FilterOpts, typeHash [][32]byte) (*ForwarderRequestTypeRegisteredIterator, error)

	WatchRequestTypeRegistered(opts *bind.WatchOpts, sink chan<- *ForwarderRequestTypeRegistered, typeHash [][32]byte) (event.Subscription, error)

	ParseRequestTypeRegistered(log types.Log) (*ForwarderRequestTypeRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
