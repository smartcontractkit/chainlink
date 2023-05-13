// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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

// VerifierMetaData contains all meta data concerning the Verifier contract.
var VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InactiveFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"activateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"deactivateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001a3038038062001a308339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611835620001fb600039600081816102ae015261072f01526118356000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806379ba50971161008c578063b70d929d11610066578063b70d929d146101c4578063ded6307c1461021e578063e84f128e14610231578063f2fde38b1461028e57600080fd5b806379ba50971461018e5780638da5cb5b1461019657806394d95980146101b157600080fd5b806301ffc9a7146100d4578063181f5a771461010d5780633d3ac1b5146101405780633dd864301461015357806344a0b2ad14610168578063564a0a7a1461017b575b600080fd5b6100f86100e2366004610fd7565b6001600160e01b031916633d3ac1b560e01b1490565b60405190151581526020015b60405180910390f35b60408051808201909152600e81526d2b32b934b334b2b910181718171960911b60208201525b604051610104919061104e565b61013361014e36600461107d565b6102a1565b6101666101613660046110fe565b6103e9565b005b610166610176366004611311565b610468565b6101666101893660046110fe565b610837565b6101666108b9565b6000546040516001600160a01b039091168152602001610104565b6101666101bf3660046113e9565b610963565b6101fb6101d23660046110fe565b60009081526002602052604081206001810154905491929091600160401b900463ffffffff1690565b604080519315158452602084019290925263ffffffff1690820152606001610104565b61016661022c3660046113e9565b610a5e565b61026b61023f3660046110fe565b6000908152600260205260409020805460019091015463ffffffff808316936401000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610104565b61016661029c36600461140b565b610b21565b6060336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146102ec57604051631decfebb60e31b815260040160405180910390fd5b6000808080806102fe888a018a611426565b9450945094509450945060008461031490611501565b6000818152600260208190526040909120908101549192509060ff1615610356576040516306db7ce960e31b8152600481018390526024015b60405180910390fd5b8651600081815260038301602052604090206103758483898985610b35565b61037f8984610c13565b87516020890120610394818b8a8a8a87610c5e565b604080518281526001600160a01b038e16602082015286917f5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269910160405180910390a250969c9b505050505050505050505050565b6103f1610e7d565b60008181526002602052604081208054909163ffffffff909116900361042d5760405163512d85cb60e11b81526004810183905260240161034d565b60028101805460ff1916905560405182907ff438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b490600090a25050565b85518460ff168060000361048f576040516303a1dd7360e11b815260040160405180910390fd5b601f8211156104bb57604051630185d43d60e61b815260048101839052601f602482015260440161034d565b6104c6816003611554565b821161050557816104d8826003611554565b6104e3906001611573565b6040516313bb3cdb60e31b81526004810192909252602482015260440161034d565b61050d610e7d565b60008981526002602052604081208054909163ffffffff9091169082906105338361158c565b82546101009290920a63ffffffff818102199093169183160217909155825460009250610568918d91168c8c8c8c8c8c610ed2565b60008181526003840160205260408120805460ff8c1660ff199182161782556002909101805490911660011790559091505b8a518160ff1610156106f95760008b8260ff16815181106105bd576105bd611528565b6020026020010151905060006001600160a01b0316816001600160a01b0316036105fa5760405163d92e233d60e01b815260040160405180910390fd5b600080600085815260038701602090815260408083206001600160a01b0387168452600190810190925290912054610100900460ff1690811115610640576106406115af565b148015915061066257604051633d9ef1f160e21b815260040160405180910390fd5b6040805180820190915260ff841681526020810160019052600085815260038701602090815260408083206001600160a01b03871684526001908101835292208351815460ff90911660ff198216811783559285015191939192849261ffff1990921617906101009084908111156106dc576106dc6115af565b0217905550905050505080806106f1906115c5565b91505061059a565b508154600163ffffffff9091161115610794576001820154604051632cc9947760e01b81526004810191909152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cc9947790604401600060405180830381600087803b15801561077b57600080fd5b505af115801561078f573d6000803e3d6000fd5b505050505b8a7fa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da8360000160049054906101000a900463ffffffff16838560000160009054906101000a900463ffffffff168e8e8e8e8e8e6040516107fc99989796959493929190611658565b60405180910390a281546bffffffffffffffff00000000191664010000000063ffffffff431602178255600190910155505050505050505050565b61083f610e7d565b60008181526002602052604081208054909163ffffffff909116900361087b5760405163512d85cb60e11b81526004810183905260240161034d565b60028101805460ff1916600117905560405182907ffc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd06190600090a25050565b6001546001600160a01b0316331461090c5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015260640161034d565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61096b610e7d565b6000828152600260205260409020816109975760405163e332262760e01b815260040160405180910390fd5b600082815260038201602052604081205460ff1690036109d457604051638bca631160e01b8152600481018490526024810183905260440161034d565b80600101548203610a0257604051635201e00b60e11b8152600481018490526024810183905260440161034d565b600082815260038201602052604090819020600201805460ff191690555183907f0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c90610a519085815260200190565b60405180910390a2505050565b610a66610e7d565b600082815260026020526040902081610a925760405163e332262760e01b815260040160405180910390fd5b600082815260038201602052604081205460ff169003610acf57604051638bca631160e01b8152600481018490526024810183905260440161034d565b600082815260038201602052604090819020600201805460ff191660011790555183907f54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe774690610a519085815260200190565b610b29610e7d565b610b3281610f2e565b50565b8054600090610b489060ff1660016116ee565b825490915060ff16600003610b7a57604051638bca631160e01b8152600481018790526024810186905260440161034d565b600282015460ff16610ba95760405163fc10a28360e01b8152600481018790526024810186905260440161034d565b8060ff16845114610bdc5783516040516329a4514160e11b8152600481019190915260ff8216602482015260440161034d565b8251845114610c0b5783518351604051631e1a628160e31b81526004810192909252602482015260440161034d565b505050505050565b6020820151815463ffffffff600883901c811691600160401b900416811115610c585782546bffffffff00000000000000001916600160401b63ffffffff8316021783555b50505050565b60008686604051602001610c73929190611707565b6040516020818303038152906040528051906020012090506000610ca7604080518082019091526000808252602082015290565b8651600090815b81811015610e2e57600186898360208110610ccb57610ccb611528565b610cd891901a601b6116ee565b8c8481518110610cea57610cea611528565b60200260200101518c8581518110610d0457610d04611528565b602002602001015160405160008152602001604052604051610d42949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610d64573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526001808d01602090815291859020848601909552845460ff808216865293995093955090850192610100900490911690811115610dbe57610dbe6115af565b6001811115610dcf57610dcf6115af565b9052509350600184602001516001811115610dec57610dec6115af565b14610e0a57604051631decfebb60e31b815260040160405180910390fd5b836000015160080260ff166001901b8501945080610e2790611743565b9050610cae565b50837e01010101010101010101010101010101010101010101010101010101010101851614610e7057604051633d9ef1f160e21b815260040160405180910390fd5b5050505050505050505050565b6000546001600160a01b03163314610ed05760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015260640161034d565b565b6000808946308b8b8b8b8b8b8b604051602001610ef89a9998979695949392919061175c565b60408051601f1981840301815291905280516020909101206001600160f01b0316600160f01b1791505098975050505050505050565b336001600160a01b03821603610f865760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161034d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215610fe957600080fd5b81356001600160e01b03198116811461100157600080fd5b9392505050565b6000815180845260005b8181101561102e57602081850181015186830182015201611012565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006110016020830184611008565b80356001600160a01b038116811461107857600080fd5b919050565b60008060006040848603121561109257600080fd5b833567ffffffffffffffff808211156110aa57600080fd5b818601915086601f8301126110be57600080fd5b8135818111156110cd57600080fd5b8760208285010111156110df57600080fd5b6020928301955093506110f59186019050611061565b90509250925092565b60006020828403121561111057600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff8111828210171561115057611150611117565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561117f5761117f611117565b604052919050565b600067ffffffffffffffff8211156111a1576111a1611117565b5060051b60200190565b600082601f8301126111bc57600080fd5b813560206111d16111cc83611187565b611156565b82815260059290921b840181019181810190868411156111f057600080fd5b8286015b848110156112125761120581611061565b83529183019183016111f4565b509695505050505050565b600082601f83011261122e57600080fd5b8135602061123e6111cc83611187565b82815260059290921b8401810191818101908684111561125d57600080fd5b8286015b848110156112125780358352918301918301611261565b803560ff8116811461107857600080fd5b600082601f83011261129a57600080fd5b813567ffffffffffffffff8111156112b4576112b4611117565b6112c7601f8201601f1916602001611156565b8181528460208386010111156112dc57600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461107857600080fd5b600080600080600080600060e0888a03121561132c57600080fd5b87359650602088013567ffffffffffffffff8082111561134b57600080fd5b6113578b838c016111ab565b975060408a013591508082111561136d57600080fd5b6113798b838c0161121d565b965061138760608b01611278565b955060808a013591508082111561139d57600080fd5b6113a98b838c01611289565b94506113b760a08b016112f9565b935060c08a01359150808211156113cd57600080fd5b506113da8a828b01611289565b91505092959891949750929550565b600080604083850312156113fc57600080fd5b50508035926020909101359150565b60006020828403121561141d57600080fd5b61100182611061565b600080600080600060e0868803121561143e57600080fd5b86601f87011261144d57600080fd5b61145561112d565b80606088018981111561146757600080fd5b885b81811015611481578035845260209384019301611469565b5090965035905067ffffffffffffffff8082111561149e57600080fd5b6114aa89838a01611289565b955060808801359150808211156114c057600080fd5b6114cc89838a0161121d565b945060a08801359150808211156114e257600080fd5b506114ef8882890161121d565b9598949750929560c001359392505050565b80516020808301519190811015611522576000198160200360031b1b821691505b50919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600081600019048311821515161561156e5761156e61153e565b500290565b808201808211156115865761158661153e565b92915050565b600063ffffffff8083168181036115a5576115a561153e565b6001019392505050565b634e487b7160e01b600052602160045260246000fd5b600060ff821660ff81036115db576115db61153e565b60010192915050565b600081518084526020808501945080840160005b8381101561161d5781516001600160a01b0316875295820195908201906001016115f8565b509495945050505050565b600081518084526020808501945080840160005b8381101561161d5781518752958201959082019060010161163c565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526116888184018a6115e4565b9050828103608084015261169c8189611628565b905060ff871660a084015282810360c08401526116b98187611008565b905067ffffffffffffffff851660e08401528281036101008401526116de8185611008565b9c9b505050505050505050505050565b60ff81811683821601908111156115865761158661153e565b828152600060208083018460005b600381101561173257815183529183019190830190600101611715565b505050506080820190509392505050565b6000600182016117555761175561153e565b5060010190565b8a8152602081018a90526001600160a01b038916604082015267ffffffffffffffff88811660608301526101406080830181905260009161179f8483018b6115e4565b915083820360a08501526117b3828a611628565b915060ff881660c085015283820360e08501526117d08288611008565b90861661010085015283810361012085015290506117ee8185611008565b9d9c5050505050505050505050505056fea2646970667358221220e10b508a084e3a946a51a45c6669b19f9170e45255769502b9ef76451c7d92f364736f6c63430008100033",
}

// VerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierMetaData.ABI instead.
var VerifierABI = VerifierMetaData.ABI

// VerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierMetaData.Bin instead.
var VerifierBin = VerifierMetaData.Bin

// DeployVerifier deploys a new Ethereum contract, binding an instance of Verifier to it.
func DeployVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, verifierProxyAddr common.Address) (common.Address, *types.Transaction, *Verifier, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierBin), backend, verifierProxyAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// Verifier is an auto generated Go binding around an Ethereum contract.
type Verifier struct {
	VerifierCaller     // Read-only binding to the contract
	VerifierTransactor // Write-only binding to the contract
	VerifierFilterer   // Log filterer for contract events
}

// VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierSession struct {
	Contract     *Verifier         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierCallerSession struct {
	Contract *VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierTransactorSession struct {
	Contract     *VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierRaw struct {
	Contract *Verifier // Generic contract binding to access the raw methods on
}

// VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierCallerRaw struct {
	Contract *VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierTransactorRaw struct {
	Contract *VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifier creates a new instance of Verifier, bound to a specific deployed contract.
func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// NewVerifierCaller creates a new read-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

// NewVerifierTransactor creates a new write-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

// NewVerifierFilterer creates a new log filterer instance of Verifier, bound to a specific deployed contract.
func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

// bindVerifier binds a generic wrapper to an already deployed contract.
func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0xe84f128e.
//
// Solidity: function latestConfigDetails(bytes32 feedId) view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_Verifier *VerifierCaller) LatestConfigDetails(opts *bind.CallOpts, feedId [32]byte) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDetails", feedId)

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0xe84f128e.
//
// Solidity: function latestConfigDetails(bytes32 feedId) view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_Verifier *VerifierSession) LatestConfigDetails(feedId [32]byte) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts, feedId)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0xe84f128e.
//
// Solidity: function latestConfigDetails(bytes32 feedId) view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_Verifier *VerifierCallerSession) LatestConfigDetails(feedId [32]byte) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts, feedId)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xb70d929d.
//
// Solidity: function latestConfigDigestAndEpoch(bytes32 feedId) view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts, feedId [32]byte) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDigestAndEpoch", feedId)

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xb70d929d.
//
// Solidity: function latestConfigDigestAndEpoch(bytes32 feedId) view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierSession) LatestConfigDigestAndEpoch(feedId [32]byte) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _Verifier.Contract.LatestConfigDigestAndEpoch(&_Verifier.CallOpts, feedId)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xb70d929d.
//
// Solidity: function latestConfigDigestAndEpoch(bytes32 feedId) view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierCallerSession) LatestConfigDigestAndEpoch(feedId [32]byte) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _Verifier.Contract.LatestConfigDigestAndEpoch(&_Verifier.CallOpts, feedId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Verifier *VerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Verifier *VerifierSession) Owner() (common.Address, error) {
	return _Verifier.Contract.Owner(&_Verifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Verifier *VerifierCallerSession) Owner() (common.Address, error) {
	return _Verifier.Contract.Owner(&_Verifier.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool isVerifier)
func (_Verifier *VerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool isVerifier)
func (_Verifier *VerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Verifier.Contract.SupportsInterface(&_Verifier.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool isVerifier)
func (_Verifier *VerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Verifier.Contract.SupportsInterface(&_Verifier.CallOpts, interfaceId)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Verifier *VerifierCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Verifier *VerifierSession) TypeAndVersion() (string, error) {
	return _Verifier.Contract.TypeAndVersion(&_Verifier.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Verifier *VerifierCallerSession) TypeAndVersion() (string, error) {
	return _Verifier.Contract.TypeAndVersion(&_Verifier.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Verifier *VerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Verifier *VerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _Verifier.Contract.AcceptOwnership(&_Verifier.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Verifier *VerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Verifier.Contract.AcceptOwnership(&_Verifier.TransactOpts)
}

// ActivateConfig is a paid mutator transaction binding the contract method 0xded6307c.
//
// Solidity: function activateConfig(bytes32 feedId, bytes32 configDigest) returns()
func (_Verifier *VerifierTransactor) ActivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "activateConfig", feedId, configDigest)
}

// ActivateConfig is a paid mutator transaction binding the contract method 0xded6307c.
//
// Solidity: function activateConfig(bytes32 feedId, bytes32 configDigest) returns()
func (_Verifier *VerifierSession) ActivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

// ActivateConfig is a paid mutator transaction binding the contract method 0xded6307c.
//
// Solidity: function activateConfig(bytes32 feedId, bytes32 configDigest) returns()
func (_Verifier *VerifierTransactorSession) ActivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

// ActivateFeed is a paid mutator transaction binding the contract method 0x3dd86430.
//
// Solidity: function activateFeed(bytes32 feedId) returns()
func (_Verifier *VerifierTransactor) ActivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "activateFeed", feedId)
}

// ActivateFeed is a paid mutator transaction binding the contract method 0x3dd86430.
//
// Solidity: function activateFeed(bytes32 feedId) returns()
func (_Verifier *VerifierSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateFeed(&_Verifier.TransactOpts, feedId)
}

// ActivateFeed is a paid mutator transaction binding the contract method 0x3dd86430.
//
// Solidity: function activateFeed(bytes32 feedId) returns()
func (_Verifier *VerifierTransactorSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateFeed(&_Verifier.TransactOpts, feedId)
}

// DeactivateConfig is a paid mutator transaction binding the contract method 0x94d95980.
//
// Solidity: function deactivateConfig(bytes32 feedId, bytes32 configDigest) returns()
func (_Verifier *VerifierTransactor) DeactivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "deactivateConfig", feedId, configDigest)
}

// DeactivateConfig is a paid mutator transaction binding the contract method 0x94d95980.
//
// Solidity: function deactivateConfig(bytes32 feedId, bytes32 configDigest) returns()
func (_Verifier *VerifierSession) DeactivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

// DeactivateConfig is a paid mutator transaction binding the contract method 0x94d95980.
//
// Solidity: function deactivateConfig(bytes32 feedId, bytes32 configDigest) returns()
func (_Verifier *VerifierTransactorSession) DeactivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

// DeactivateFeed is a paid mutator transaction binding the contract method 0x564a0a7a.
//
// Solidity: function deactivateFeed(bytes32 feedId) returns()
func (_Verifier *VerifierTransactor) DeactivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "deactivateFeed", feedId)
}

// DeactivateFeed is a paid mutator transaction binding the contract method 0x564a0a7a.
//
// Solidity: function deactivateFeed(bytes32 feedId) returns()
func (_Verifier *VerifierSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateFeed(&_Verifier.TransactOpts, feedId)
}

// DeactivateFeed is a paid mutator transaction binding the contract method 0x564a0a7a.
//
// Solidity: function deactivateFeed(bytes32 feedId) returns()
func (_Verifier *VerifierTransactorSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateFeed(&_Verifier.TransactOpts, feedId)
}

// SetConfig is a paid mutator transaction binding the contract method 0x44a0b2ad.
//
// Solidity: function setConfig(bytes32 feedId, address[] signers, bytes32[] offchainTransmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_Verifier *VerifierTransactor) SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "setConfig", feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0x44a0b2ad.
//
// Solidity: function setConfig(bytes32 feedId, address[] signers, bytes32[] offchainTransmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_Verifier *VerifierSession) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0x44a0b2ad.
//
// Solidity: function setConfig(bytes32 feedId, address[] signers, bytes32[] offchainTransmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_Verifier *VerifierTransactorSession) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Verifier *VerifierTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Verifier *VerifierSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.TransferOwnership(&_Verifier.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Verifier *VerifierTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.TransferOwnership(&_Verifier.TransactOpts, to)
}

// Verify is a paid mutator transaction binding the contract method 0x3d3ac1b5.
//
// Solidity: function verify(bytes signedReport, address sender) returns(bytes response)
func (_Verifier *VerifierTransactor) Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "verify", signedReport, sender)
}

// Verify is a paid mutator transaction binding the contract method 0x3d3ac1b5.
//
// Solidity: function verify(bytes signedReport, address sender) returns(bytes response)
func (_Verifier *VerifierSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.Verify(&_Verifier.TransactOpts, signedReport, sender)
}

// Verify is a paid mutator transaction binding the contract method 0x3d3ac1b5.
//
// Solidity: function verify(bytes signedReport, address sender) returns(bytes response)
func (_Verifier *VerifierTransactorSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.Verify(&_Verifier.TransactOpts, signedReport, sender)
}

// VerifierConfigActivatedIterator is returned from FilterConfigActivated and is used to iterate over the raw logs and unpacked data for ConfigActivated events raised by the Verifier contract.
type VerifierConfigActivatedIterator struct {
	Event *VerifierConfigActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierConfigActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigActivated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierConfigActivated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierConfigActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierConfigActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierConfigActivated represents a ConfigActivated event raised by the Verifier contract.
type VerifierConfigActivated struct {
	FeedId       [32]byte
	ConfigDigest [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConfigActivated is a free log retrieval operation binding the contract event 0x54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe7746.
//
// Solidity: event ConfigActivated(bytes32 indexed feedId, bytes32 configDigest)
func (_Verifier *VerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierConfigActivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigActivatedIterator{contract: _Verifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

// WatchConfigActivated is a free log subscription operation binding the contract event 0x54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe7746.
//
// Solidity: event ConfigActivated(bytes32 indexed feedId, bytes32 configDigest)
func (_Verifier *VerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigActivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierConfigActivated)
				if err := _Verifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
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

// ParseConfigActivated is a log parse operation binding the contract event 0x54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe7746.
//
// Solidity: event ConfigActivated(bytes32 indexed feedId, bytes32 configDigest)
func (_Verifier *VerifierFilterer) ParseConfigActivated(log types.Log) (*VerifierConfigActivated, error) {
	event := new(VerifierConfigActivated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierConfigDeactivatedIterator is returned from FilterConfigDeactivated and is used to iterate over the raw logs and unpacked data for ConfigDeactivated events raised by the Verifier contract.
type VerifierConfigDeactivatedIterator struct {
	Event *VerifierConfigDeactivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierConfigDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigDeactivated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierConfigDeactivated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierConfigDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierConfigDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierConfigDeactivated represents a ConfigDeactivated event raised by the Verifier contract.
type VerifierConfigDeactivated struct {
	FeedId       [32]byte
	ConfigDigest [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConfigDeactivated is a free log retrieval operation binding the contract event 0x0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c.
//
// Solidity: event ConfigDeactivated(bytes32 indexed feedId, bytes32 configDigest)
func (_Verifier *VerifierFilterer) FilterConfigDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierConfigDeactivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigDeactivatedIterator{contract: _Verifier.contract, event: "ConfigDeactivated", logs: logs, sub: sub}, nil
}

// WatchConfigDeactivated is a free log subscription operation binding the contract event 0x0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c.
//
// Solidity: event ConfigDeactivated(bytes32 indexed feedId, bytes32 configDigest)
func (_Verifier *VerifierFilterer) WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigDeactivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierConfigDeactivated)
				if err := _Verifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
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

// ParseConfigDeactivated is a log parse operation binding the contract event 0x0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c.
//
// Solidity: event ConfigDeactivated(bytes32 indexed feedId, bytes32 configDigest)
func (_Verifier *VerifierFilterer) ParseConfigDeactivated(log types.Log) (*VerifierConfigDeactivated, error) {
	event := new(VerifierConfigDeactivated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the Verifier contract.
type VerifierConfigSetIterator struct {
	Event *VerifierConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierConfigSet represents a ConfigSet event raised by the Verifier contract.
type VerifierConfigSet struct {
	FeedId                    [32]byte
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	OffchainTransmitters      [][32]byte
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da.
//
// Solidity: event ConfigSet(bytes32 indexed feedId, uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, bytes32[] offchainTransmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_Verifier *VerifierFilterer) FilterConfigSet(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierConfigSetIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigSet", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigSetIterator{contract: _Verifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da.
//
// Solidity: event ConfigSet(bytes32 indexed feedId, uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, bytes32[] offchainTransmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_Verifier *VerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VerifierConfigSet, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigSet", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierConfigSet)
				if err := _Verifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0xa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da.
//
// Solidity: event ConfigSet(bytes32 indexed feedId, uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, bytes32[] offchainTransmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_Verifier *VerifierFilterer) ParseConfigSet(log types.Log) (*VerifierConfigSet, error) {
	event := new(VerifierConfigSet)
	if err := _Verifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierFeedActivatedIterator is returned from FilterFeedActivated and is used to iterate over the raw logs and unpacked data for FeedActivated events raised by the Verifier contract.
type VerifierFeedActivatedIterator struct {
	Event *VerifierFeedActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierFeedActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierFeedActivated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierFeedActivated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierFeedActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierFeedActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierFeedActivated represents a FeedActivated event raised by the Verifier contract.
type VerifierFeedActivated struct {
	FeedId [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeedActivated is a free log retrieval operation binding the contract event 0xf438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b4.
//
// Solidity: event FeedActivated(bytes32 indexed feedId)
func (_Verifier *VerifierFilterer) FilterFeedActivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierFeedActivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "FeedActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifierFeedActivatedIterator{contract: _Verifier.contract, event: "FeedActivated", logs: logs, sub: sub}, nil
}

// WatchFeedActivated is a free log subscription operation binding the contract event 0xf438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b4.
//
// Solidity: event FeedActivated(bytes32 indexed feedId)
func (_Verifier *VerifierFilterer) WatchFeedActivated(opts *bind.WatchOpts, sink chan<- *VerifierFeedActivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "FeedActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierFeedActivated)
				if err := _Verifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
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

// ParseFeedActivated is a log parse operation binding the contract event 0xf438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b4.
//
// Solidity: event FeedActivated(bytes32 indexed feedId)
func (_Verifier *VerifierFilterer) ParseFeedActivated(log types.Log) (*VerifierFeedActivated, error) {
	event := new(VerifierFeedActivated)
	if err := _Verifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierFeedDeactivatedIterator is returned from FilterFeedDeactivated and is used to iterate over the raw logs and unpacked data for FeedDeactivated events raised by the Verifier contract.
type VerifierFeedDeactivatedIterator struct {
	Event *VerifierFeedDeactivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierFeedDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierFeedDeactivated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierFeedDeactivated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierFeedDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierFeedDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierFeedDeactivated represents a FeedDeactivated event raised by the Verifier contract.
type VerifierFeedDeactivated struct {
	FeedId [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeedDeactivated is a free log retrieval operation binding the contract event 0xfc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd061.
//
// Solidity: event FeedDeactivated(bytes32 indexed feedId)
func (_Verifier *VerifierFilterer) FilterFeedDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierFeedDeactivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "FeedDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifierFeedDeactivatedIterator{contract: _Verifier.contract, event: "FeedDeactivated", logs: logs, sub: sub}, nil
}

// WatchFeedDeactivated is a free log subscription operation binding the contract event 0xfc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd061.
//
// Solidity: event FeedDeactivated(bytes32 indexed feedId)
func (_Verifier *VerifierFilterer) WatchFeedDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierFeedDeactivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "FeedDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierFeedDeactivated)
				if err := _Verifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
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

// ParseFeedDeactivated is a log parse operation binding the contract event 0xfc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd061.
//
// Solidity: event FeedDeactivated(bytes32 indexed feedId)
func (_Verifier *VerifierFilterer) ParseFeedDeactivated(log types.Log) (*VerifierFeedDeactivated, error) {
	event := new(VerifierFeedDeactivated)
	if err := _Verifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Verifier contract.
type VerifierOwnershipTransferRequestedIterator struct {
	Event *VerifierOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Verifier contract.
type VerifierOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierOwnershipTransferRequestedIterator{contract: _Verifier.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierOwnershipTransferRequested)
				if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifierOwnershipTransferRequested, error) {
	event := new(VerifierOwnershipTransferRequested)
	if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Verifier contract.
type VerifierOwnershipTransferredIterator struct {
	Event *VerifierOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierOwnershipTransferred represents a OwnershipTransferred event raised by the Verifier contract.
type VerifierOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierOwnershipTransferredIterator{contract: _Verifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierOwnershipTransferred)
				if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) ParseOwnershipTransferred(log types.Log) (*VerifierOwnershipTransferred, error) {
	event := new(VerifierOwnershipTransferred)
	if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierReportVerifiedIterator is returned from FilterReportVerified and is used to iterate over the raw logs and unpacked data for ReportVerified events raised by the Verifier contract.
type VerifierReportVerifiedIterator struct {
	Event *VerifierReportVerified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierReportVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierReportVerified)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierReportVerified)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierReportVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierReportVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierReportVerified represents a ReportVerified event raised by the Verifier contract.
type VerifierReportVerified struct {
	FeedId     [32]byte
	ReportHash [32]byte
	Requester  common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReportVerified is a free log retrieval operation binding the contract event 0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269.
//
// Solidity: event ReportVerified(bytes32 indexed feedId, bytes32 reportHash, address requester)
func (_Verifier *VerifierFilterer) FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierReportVerifiedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifierReportVerifiedIterator{contract: _Verifier.contract, event: "ReportVerified", logs: logs, sub: sub}, nil
}

// WatchReportVerified is a free log subscription operation binding the contract event 0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269.
//
// Solidity: event ReportVerified(bytes32 indexed feedId, bytes32 reportHash, address requester)
func (_Verifier *VerifierFilterer) WatchReportVerified(opts *bind.WatchOpts, sink chan<- *VerifierReportVerified, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierReportVerified)
				if err := _Verifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

// ParseReportVerified is a log parse operation binding the contract event 0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269.
//
// Solidity: event ReportVerified(bytes32 indexed feedId, bytes32 reportHash, address requester)
func (_Verifier *VerifierFilterer) ParseReportVerified(log types.Log) (*VerifierReportVerified, error) {
	event := new(VerifierReportVerified)
	if err := _Verifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
