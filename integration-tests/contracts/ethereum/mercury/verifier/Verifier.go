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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620018b4380380620018b48339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b6080516116b9620001fb60003960008181610262015261063701526116b96000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80638da5cb5b116100715780638da5cb5b1461014a57806394d9598014610165578063b70d929d14610178578063ded6307c146101d2578063e84f128e146101e5578063f2fde38b1461024257600080fd5b806301ffc9a7146100ae578063181f5a77146100e75780633d3ac1b51461011a57806344a0b2ad1461012d57806379ba509714610142575b600080fd5b6100d26100bc366004610e5b565b6001600160e01b031916633d3ac1b560e01b1490565b60405190151581526020015b60405180910390f35b60408051808201909152600e81526d2b32b934b334b2b910181718171960911b60208201525b6040516100de9190610ed2565b61010d610128366004610f01565b610255565b61014061013b36600461117c565b61036d565b005b61014061073f565b6000546040516001600160a01b0390911681526020016100de565b610140610173366004611254565b6107e9565b6101af610186366004611276565b60009081526002602052604081206001810154905491929091600160401b900463ffffffff1690565b604080519315158452602084019290925263ffffffff16908201526060016100de565b6101406101e0366004611254565b6108e3565b61021f6101f3366004611276565b6000908152600260205260409020805460019091015463ffffffff808316936401000000009093041691565b6040805163ffffffff9485168152939092166020840152908201526060016100de565b61014061025036600461128f565b6109a5565b6060336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146102a057604051631decfebb60e31b815260040160405180910390fd5b6000808080806102b2888a018a6112aa565b945094509450945094506000846102c890611385565b60008181526002602081815260408084208b5180865293810190925290922092935090916102f984838989856109b9565b6103038984610a97565b87516020890120610318818b8a8a8a87610ae2565b604080518281526001600160a01b038e16602082015286917f5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269910160405180910390a250969c9b505050505050505050505050565b85518460ff1680600003610394576040516303a1dd7360e11b815260040160405180910390fd5b601f8211156103c557604051630185d43d60e61b815260048101839052601f60248201526044015b60405180910390fd5b6103d08160036113d8565b821161040f57816103e28260036113d8565b6103ed9060016113f7565b6040516313bb3cdb60e31b8152600481019290925260248201526044016103bc565b610417610d01565b60008981526002602052604081208054909163ffffffff90911690829061043d83611410565b82546101009290920a63ffffffff818102199093169183160217909155825460009250610472918d91168c8c8c8c8c8c610d56565b6000818152600280850160205260408220805460ff8d1660ff199182161782559101805490911660011790559091505b8a518160ff1610156106015760008b8260ff16815181106104c5576104c56113ac565b6020026020010151905060006001600160a01b0316816001600160a01b0316036105025760405163d92e233d60e01b815260040160405180910390fd5b600080600085815260028701602090815260408083206001600160a01b0387168452600190810190925290912054610100900460ff169081111561054857610548611433565b148015915061056a57604051633d9ef1f160e21b815260040160405180910390fd5b6040805180820190915260ff841681526020810160019052600085815260028701602090815260408083206001600160a01b03871684526001908101835292208351815460ff90911660ff198216811783559285015191939192849261ffff1990921617906101009084908111156105e4576105e4611433565b0217905550905050505080806105f990611449565b9150506104a2565b508154600163ffffffff909116111561069c576001820154604051632cc9947760e01b81526004810191909152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cc9947790604401600060405180830381600087803b15801561068357600080fd5b505af1158015610697573d6000803e3d6000fd5b505050505b8a7fa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da8360000160049054906101000a900463ffffffff16838560000160009054906101000a900463ffffffff168e8e8e8e8e8e604051610704999897969594939291906114dc565b60405180910390a281546bffffffffffffffff00000000191664010000000063ffffffff431602178255600190910155505050505050505050565b6001546001600160a01b031633146107925760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016103bc565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6107f1610d01565b60008281526002602052604090208161081d5760405163e332262760e01b815260040160405180910390fd5b600082815260028201602052604081205460ff16900361085a57604051638bca631160e01b815260048101849052602481018390526044016103bc565b8060010154820361088857604051635201e00b60e11b815260048101849052602481018390526044016103bc565b6000828152600280830160205260409182902001805460ff191690555183907f0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c906108d69085815260200190565b60405180910390a2505050565b6108eb610d01565b6000828152600260205260409020816109175760405163e332262760e01b815260040160405180910390fd5b600082815260028201602052604081205460ff16900361095457604051638bca631160e01b815260048101849052602481018390526044016103bc565b6000828152600280830160205260409182902001805460ff191660011790555183907f54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe7746906108d69085815260200190565b6109ad610d01565b6109b681610db2565b50565b80546000906109cc9060ff166001611572565b825490915060ff166000036109fe57604051638bca631160e01b815260048101879052602481018690526044016103bc565b600282015460ff16610a2d5760405163fc10a28360e01b815260048101879052602481018690526044016103bc565b8060ff16845114610a605783516040516329a4514160e11b8152600481019190915260ff821660248201526044016103bc565b8251845114610a8f5783518351604051631e1a628160e31b8152600481019290925260248201526044016103bc565b505050505050565b6020820151815463ffffffff600883901c811691600160401b900416811115610adc5782546bffffffff00000000000000001916600160401b63ffffffff8316021783555b50505050565b60008686604051602001610af792919061158b565b6040516020818303038152906040528051906020012090506000610b2b604080518082019091526000808252602082015290565b8651600090815b81811015610cb257600186898360208110610b4f57610b4f6113ac565b610b5c91901a601b611572565b8c8481518110610b6e57610b6e6113ac565b60200260200101518c8581518110610b8857610b886113ac565b602002602001015160405160008152602001604052604051610bc6949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610be8573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526001808d01602090815291859020848601909552845460ff808216865293995093955090850192610100900490911690811115610c4257610c42611433565b6001811115610c5357610c53611433565b9052509350600184602001516001811115610c7057610c70611433565b14610c8e57604051631decfebb60e31b815260040160405180910390fd5b836000015160080260ff166001901b8501945080610cab906115c7565b9050610b32565b50837e01010101010101010101010101010101010101010101010101010101010101851614610cf457604051633d9ef1f160e21b815260040160405180910390fd5b5050505050505050505050565b6000546001600160a01b03163314610d545760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016103bc565b565b6000808946308b8b8b8b8b8b8b604051602001610d7c9a999897969594939291906115e0565b60408051601f1981840301815291905280516020909101206001600160f01b0316600160f01b1791505098975050505050505050565b336001600160a01b03821603610e0a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103bc565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215610e6d57600080fd5b81356001600160e01b031981168114610e8557600080fd5b9392505050565b6000815180845260005b81811015610eb257602081850181015186830182015201610e96565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000610e856020830184610e8c565b80356001600160a01b0381168114610efc57600080fd5b919050565b600080600060408486031215610f1657600080fd5b833567ffffffffffffffff80821115610f2e57600080fd5b818601915086601f830112610f4257600080fd5b813581811115610f5157600080fd5b876020828501011115610f6357600080fd5b602092830195509350610f799186019050610ee5565b90509250925092565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610fbb57610fbb610f82565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610fea57610fea610f82565b604052919050565b600067ffffffffffffffff82111561100c5761100c610f82565b5060051b60200190565b600082601f83011261102757600080fd5b8135602061103c61103783610ff2565b610fc1565b82815260059290921b8401810191818101908684111561105b57600080fd5b8286015b8481101561107d5761107081610ee5565b835291830191830161105f565b509695505050505050565b600082601f83011261109957600080fd5b813560206110a961103783610ff2565b82815260059290921b840181019181810190868411156110c857600080fd5b8286015b8481101561107d57803583529183019183016110cc565b803560ff81168114610efc57600080fd5b600082601f83011261110557600080fd5b813567ffffffffffffffff81111561111f5761111f610f82565b611132601f8201601f1916602001610fc1565b81815284602083860101111561114757600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114610efc57600080fd5b600080600080600080600060e0888a03121561119757600080fd5b87359650602088013567ffffffffffffffff808211156111b657600080fd5b6111c28b838c01611016565b975060408a01359150808211156111d857600080fd5b6111e48b838c01611088565b96506111f260608b016110e3565b955060808a013591508082111561120857600080fd5b6112148b838c016110f4565b945061122260a08b01611164565b935060c08a013591508082111561123857600080fd5b506112458a828b016110f4565b91505092959891949750929550565b6000806040838503121561126757600080fd5b50508035926020909101359150565b60006020828403121561128857600080fd5b5035919050565b6000602082840312156112a157600080fd5b610e8582610ee5565b600080600080600060e086880312156112c257600080fd5b86601f8701126112d157600080fd5b6112d9610f98565b8060608801898111156112eb57600080fd5b885b818110156113055780358452602093840193016112ed565b5090965035905067ffffffffffffffff8082111561132257600080fd5b61132e89838a016110f4565b9550608088013591508082111561134457600080fd5b61135089838a01611088565b945060a088013591508082111561136657600080fd5b5061137388828901611088565b9598949750929560c001359392505050565b805160208083015191908110156113a6576000198160200360031b1b821691505b50919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60008160001904831182151516156113f2576113f26113c2565b500290565b8082018082111561140a5761140a6113c2565b92915050565b600063ffffffff808316818103611429576114296113c2565b6001019392505050565b634e487b7160e01b600052602160045260246000fd5b600060ff821660ff810361145f5761145f6113c2565b60010192915050565b600081518084526020808501945080840160005b838110156114a15781516001600160a01b03168752958201959082019060010161147c565b509495945050505050565b600081518084526020808501945080840160005b838110156114a1578151875295820195908201906001016114c0565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261150c8184018a611468565b9050828103608084015261152081896114ac565b905060ff871660a084015282810360c084015261153d8187610e8c565b905067ffffffffffffffff851660e08401528281036101008401526115628185610e8c565b9c9b505050505050505050505050565b60ff818116838216019081111561140a5761140a6113c2565b828152600060208083018460005b60038110156115b657815183529183019190830190600101611599565b505050506080820190509392505050565b6000600182016115d9576115d96113c2565b5060010190565b8a8152602081018a90526001600160a01b038916604082015267ffffffffffffffff8881166060830152610140608083018190526000916116238483018b611468565b915083820360a0850152611637828a6114ac565b915060ff881660c085015283820360e08501526116548288610e8c565b90861661010085015283810361012085015290506116728185610e8c565b9d9c5050505050505050505050505056fea26469706673582212209d5004dcecd6031e5ce7e3afbfb3b18bc5546731db67ae0a31cdc5750f919bea64736f6c63430008100033",
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
