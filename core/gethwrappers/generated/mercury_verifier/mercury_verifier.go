// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_verifier

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

var MercuryVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001e7c38038062001e7c8339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611c81620001fb600039600081816102e101526107ff0152611c816000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638da5cb5b11610081578063ded6307c1161005b578063ded6307c14610244578063e84f128e14610257578063f2fde38b146102b457600080fd5b80638da5cb5b146101aa57806394d95980146101d2578063b70d929d146101e557600080fd5b80633d3ac1b5116100b25780633d3ac1b51461017a57806344a0b2ad1461018d57806379ba5097146101a257600080fd5b806301ffc9a7146100ce578063181f5a7714610138575b600080fd5b6101236100dc3660046112f7565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b60408051808201909152600e81527f566572696669657220302e302e3200000000000000000000000000000000000060208201525b60405161012f91906113a4565b61016d6101883660046113e0565b6102c7565b6101a061019b3660046116b0565b610412565b005b6101a0610927565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012f565b6101a06101e0366004611788565b610a24565b6102216101f33660046117aa565b6000908152600260205260408120600181015490549192909168010000000000000000900463ffffffff1690565b604080519315158452602084019290925263ffffffff169082015260600161012f565b6101a0610252366004611788565b610b87565b6102916102653660046117aa565b6000908152600260205260409020805460019091015463ffffffff808316936401000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600161012f565b6101a06102c23660046117c3565b610c99565b60603373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610338576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008080808061034a888a018a6117de565b94509450945094509450600084610360906118b9565b60008181526002602081815260408084208b5180865293810190925290922092935090916103918483898985610cad565b61039b8984610def565b875160208901206103b0818b8a8a8a87610e57565b6040805182815273ffffffffffffffffffffffffffffffffffffffff8e16602082015286917f5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269910160405180910390a250969c9b505050505050505050505050565b85518460ff1680600003610452576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f82111561049c576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b6104a781600361195c565b82116104ff57816104b982600361195c565b6104c4906001611999565b6040517f9dd9e6d800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610493565b6105076110d3565b60008981526002602052604081208054909163ffffffff90911690829061052d836119b2565b82546101009290920a63ffffffff818102199093169183160217909155825460009250610562918d91168c8c8c8c8c8c611156565b6000818152600280850160205260408220805460ff8d167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009182161782559101805490911660011790559091505b8a518160ff1610156107b05760008b8260ff16815181106105d3576105d36118fe565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610643576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452600190810190925290912054610100900460ff1690811115610696576106966119d5565b14801591506106d1576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff8416815260208101600190526000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684526001908101835292208351815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00821681178355928501519193919284927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009092161790610100908490811115610793576107936119d5565b0217905550905050505080806107a890611a04565b9150506105b0565b508154600163ffffffff90911611156108715760018201546040517f2cc994770000000000000000000000000000000000000000000000000000000081526004810191909152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cc9947790604401600060405180830381600087803b15801561085857600080fd5b505af115801561086c573d6000803e3d6000fd5b505050505b8a7fa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da8360000160049054906101000a900463ffffffff16838560000160009054906101000a900463ffffffff168e8e8e8e8e8e6040516108d999989796959493929190611aa4565b60405180910390a281547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1664010000000063ffffffff431602178255600190910155505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610493565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a2c6110d3565b600082815260026020526040902081610a71576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604081205460ff169003610ac7576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610493565b80600101548203610b0e576040517fa403c0160000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610493565b600082815260028083016020526040918290200180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690555183907f0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c90610b7a9085815260200190565b60405180910390a2505050565b610b8f6110d3565b600082815260026020526040902081610bd4576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604081205460ff169003610c2a576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610493565b600082815260028083016020526040918290200180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555183907f54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe774690610b7a9085815260200190565b610ca16110d3565b610caa81611202565b50565b8054600090610cc09060ff166001611b3a565b825490915060ff16600003610d0b576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810187905260248101869052604401610493565b600282015460ff16610d53576040517ffc10a2830000000000000000000000000000000000000000000000000000000081526004810187905260248101869052604401610493565b8060ff16845114610d9f5783516040517f5348a282000000000000000000000000000000000000000000000000000000008152600481019190915260ff82166024820152604401610493565b8251845114610de757835183516040517ff0d3140800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610493565b505050505050565b6020820151815463ffffffff600883901c81169168010000000000000000900416811115610e515782547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff166801000000000000000063ffffffff8316021783555b50505050565b60008686604051602001610e6c929190611b53565b6040516020818303038152906040528051906020012090506000610ea0604080518082019091526000808252602082015290565b8651600090815b8181101561106b57600186898360208110610ec457610ec46118fe565b610ed191901a601b611b3a565b8c8481518110610ee357610ee36118fe565b60200260200101518c8581518110610efd57610efd6118fe565b602002602001015160405160008152602001604052604051610f3b949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610f5d573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526001808d01602090815291859020848601909552845460ff808216865293995093955090850192610100900490911690811115610fe257610fe26119d5565b6001811115610ff357610ff36119d5565b9052509350600184602001516001811115611010576110106119d5565b14611047576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836000015160080260ff166001901b850194508061106490611b8f565b9050610ea7565b50837e010101010101010101010101010101010101010101010101010101010101018516146110c6576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611154576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610493565b565b6000808946308b8b8b8b8b8b8b60405160200161117c9a99989796959493929190611bc7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e010000000000000000000000000000000000000000000000000000000000001791505098975050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611281576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610493565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561130957600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461133957600080fd5b9392505050565b6000815180845260005b818110156113665760208185018101518683018201520161134a565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006113396020830184611340565b803573ffffffffffffffffffffffffffffffffffffffff811681146113db57600080fd5b919050565b6000806000604084860312156113f557600080fd5b833567ffffffffffffffff8082111561140d57600080fd5b818601915086601f83011261142157600080fd5b81358181111561143057600080fd5b87602082850101111561144257600080fd5b60209283019550935061145891860190506113b7565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156114b3576114b3611461565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561150057611500611461565b604052919050565b600067ffffffffffffffff82111561152257611522611461565b5060051b60200190565b600082601f83011261153d57600080fd5b8135602061155261154d83611508565b6114b9565b82815260059290921b8401810191818101908684111561157157600080fd5b8286015b8481101561159357611586816113b7565b8352918301918301611575565b509695505050505050565b600082601f8301126115af57600080fd5b813560206115bf61154d83611508565b82815260059290921b840181019181810190868411156115de57600080fd5b8286015b8481101561159357803583529183019183016115e2565b803560ff811681146113db57600080fd5b600082601f83011261161b57600080fd5b813567ffffffffffffffff81111561163557611635611461565b61166660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016114b9565b81815284602083860101111561167b57600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff811681146113db57600080fd5b600080600080600080600060e0888a0312156116cb57600080fd5b87359650602088013567ffffffffffffffff808211156116ea57600080fd5b6116f68b838c0161152c565b975060408a013591508082111561170c57600080fd5b6117188b838c0161159e565b965061172660608b016115f9565b955060808a013591508082111561173c57600080fd5b6117488b838c0161160a565b945061175660a08b01611698565b935060c08a013591508082111561176c57600080fd5b506117798a828b0161160a565b91505092959891949750929550565b6000806040838503121561179b57600080fd5b50508035926020909101359150565b6000602082840312156117bc57600080fd5b5035919050565b6000602082840312156117d557600080fd5b611339826113b7565b600080600080600060e086880312156117f657600080fd5b86601f87011261180557600080fd5b61180d611490565b80606088018981111561181f57600080fd5b885b81811015611839578035845260209384019301611821565b5090965035905067ffffffffffffffff8082111561185657600080fd5b61186289838a0161160a565b9550608088013591508082111561187857600080fd5b61188489838a0161159e565b945060a088013591508082111561189a57600080fd5b506118a78882890161159e565b9598949750929560c001359392505050565b805160208083015191908110156118f8577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156119945761199461192d565b500290565b808201808211156119ac576119ac61192d565b92915050565b600063ffffffff8083168181036119cb576119cb61192d565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600060ff821660ff8103611a1a57611a1a61192d565b60010192915050565b600081518084526020808501945080840160005b83811015611a6957815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611a37565b509495945050505050565b600081518084526020808501945080840160005b83811015611a6957815187529582019590820190600101611a88565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611ad48184018a611a23565b90508281036080840152611ae88189611a74565b905060ff871660a084015282810360c0840152611b058187611340565b905067ffffffffffffffff851660e0840152828103610100840152611b2a8185611340565b9c9b505050505050505050505050565b60ff81811683821601908111156119ac576119ac61192d565b828152600060208083018460005b6003811015611b7e57815183529183019190830190600101611b61565b505050506080820190509392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611bc057611bc061192d565b5060010190565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b166060850152816080850152611c148285018b611a23565b915083820360a0850152611c28828a611a74565b915060ff881660c085015283820360e0850152611c458288611340565b9086166101008501528381036101208501529050611c638185611340565b9d9c5050505050505050505050505056fea164736f6c6343000810000a",
}

var MercuryVerifierABI = MercuryVerifierMetaData.ABI

var MercuryVerifierBin = MercuryVerifierMetaData.Bin

func DeployMercuryVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, verifierProxyAddr common.Address) (common.Address, *types.Transaction, *MercuryVerifier, error) {
	parsed, err := MercuryVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryVerifierBin), backend, verifierProxyAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryVerifier{MercuryVerifierCaller: MercuryVerifierCaller{contract: contract}, MercuryVerifierTransactor: MercuryVerifierTransactor{contract: contract}, MercuryVerifierFilterer: MercuryVerifierFilterer{contract: contract}}, nil
}

type MercuryVerifier struct {
	address common.Address
	abi     abi.ABI
	MercuryVerifierCaller
	MercuryVerifierTransactor
	MercuryVerifierFilterer
}

type MercuryVerifierCaller struct {
	contract *bind.BoundContract
}

type MercuryVerifierTransactor struct {
	contract *bind.BoundContract
}

type MercuryVerifierFilterer struct {
	contract *bind.BoundContract
}

type MercuryVerifierSession struct {
	Contract     *MercuryVerifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryVerifierCallerSession struct {
	Contract *MercuryVerifierCaller
	CallOpts bind.CallOpts
}

type MercuryVerifierTransactorSession struct {
	Contract     *MercuryVerifierTransactor
	TransactOpts bind.TransactOpts
}

type MercuryVerifierRaw struct {
	Contract *MercuryVerifier
}

type MercuryVerifierCallerRaw struct {
	Contract *MercuryVerifierCaller
}

type MercuryVerifierTransactorRaw struct {
	Contract *MercuryVerifierTransactor
}

func NewMercuryVerifier(address common.Address, backend bind.ContractBackend) (*MercuryVerifier, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryVerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifier{address: address, abi: abi, MercuryVerifierCaller: MercuryVerifierCaller{contract: contract}, MercuryVerifierTransactor: MercuryVerifierTransactor{contract: contract}, MercuryVerifierFilterer: MercuryVerifierFilterer{contract: contract}}, nil
}

func NewMercuryVerifierCaller(address common.Address, caller bind.ContractCaller) (*MercuryVerifierCaller, error) {
	contract, err := bindMercuryVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierCaller{contract: contract}, nil
}

func NewMercuryVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryVerifierTransactor, error) {
	contract, err := bindMercuryVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierTransactor{contract: contract}, nil
}

func NewMercuryVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryVerifierFilterer, error) {
	contract, err := bindMercuryVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierFilterer{contract: contract}, nil
}

func bindMercuryVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MercuryVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_MercuryVerifier *MercuryVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryVerifier.Contract.MercuryVerifierCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryVerifier *MercuryVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.MercuryVerifierTransactor.contract.Transfer(opts)
}

func (_MercuryVerifier *MercuryVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.MercuryVerifierTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryVerifier *MercuryVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryVerifier.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryVerifier *MercuryVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.contract.Transfer(opts)
}

func (_MercuryVerifier *MercuryVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryVerifier *MercuryVerifierCaller) LatestConfigDetails(opts *bind.CallOpts, feedId [32]byte) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _MercuryVerifier.contract.Call(opts, &out, "latestConfigDetails", feedId)

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_MercuryVerifier *MercuryVerifierSession) LatestConfigDetails(feedId [32]byte) (LatestConfigDetails,

	error) {
	return _MercuryVerifier.Contract.LatestConfigDetails(&_MercuryVerifier.CallOpts, feedId)
}

func (_MercuryVerifier *MercuryVerifierCallerSession) LatestConfigDetails(feedId [32]byte) (LatestConfigDetails,

	error) {
	return _MercuryVerifier.Contract.LatestConfigDetails(&_MercuryVerifier.CallOpts, feedId)
}

func (_MercuryVerifier *MercuryVerifierCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts, feedId [32]byte) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _MercuryVerifier.contract.Call(opts, &out, "latestConfigDigestAndEpoch", feedId)

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_MercuryVerifier *MercuryVerifierSession) LatestConfigDigestAndEpoch(feedId [32]byte) (LatestConfigDigestAndEpoch,

	error) {
	return _MercuryVerifier.Contract.LatestConfigDigestAndEpoch(&_MercuryVerifier.CallOpts, feedId)
}

func (_MercuryVerifier *MercuryVerifierCallerSession) LatestConfigDigestAndEpoch(feedId [32]byte) (LatestConfigDigestAndEpoch,

	error) {
	return _MercuryVerifier.Contract.LatestConfigDigestAndEpoch(&_MercuryVerifier.CallOpts, feedId)
}

func (_MercuryVerifier *MercuryVerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MercuryVerifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryVerifier *MercuryVerifierSession) Owner() (common.Address, error) {
	return _MercuryVerifier.Contract.Owner(&_MercuryVerifier.CallOpts)
}

func (_MercuryVerifier *MercuryVerifierCallerSession) Owner() (common.Address, error) {
	return _MercuryVerifier.Contract.Owner(&_MercuryVerifier.CallOpts)
}

func (_MercuryVerifier *MercuryVerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MercuryVerifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryVerifier *MercuryVerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MercuryVerifier.Contract.SupportsInterface(&_MercuryVerifier.CallOpts, interfaceId)
}

func (_MercuryVerifier *MercuryVerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MercuryVerifier.Contract.SupportsInterface(&_MercuryVerifier.CallOpts, interfaceId)
}

func (_MercuryVerifier *MercuryVerifierCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryVerifier.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryVerifier *MercuryVerifierSession) TypeAndVersion() (string, error) {
	return _MercuryVerifier.Contract.TypeAndVersion(&_MercuryVerifier.CallOpts)
}

func (_MercuryVerifier *MercuryVerifierCallerSession) TypeAndVersion() (string, error) {
	return _MercuryVerifier.Contract.TypeAndVersion(&_MercuryVerifier.CallOpts)
}

func (_MercuryVerifier *MercuryVerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "acceptOwnership")
}

func (_MercuryVerifier *MercuryVerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryVerifier.Contract.AcceptOwnership(&_MercuryVerifier.TransactOpts)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryVerifier.Contract.AcceptOwnership(&_MercuryVerifier.TransactOpts)
}

func (_MercuryVerifier *MercuryVerifierTransactor) ActivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "activateConfig", feedId, configDigest)
}

func (_MercuryVerifier *MercuryVerifierSession) ActivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.ActivateConfig(&_MercuryVerifier.TransactOpts, feedId, configDigest)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) ActivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.ActivateConfig(&_MercuryVerifier.TransactOpts, feedId, configDigest)
}

func (_MercuryVerifier *MercuryVerifierTransactor) DeactivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "deactivateConfig", feedId, configDigest)
}

func (_MercuryVerifier *MercuryVerifierSession) DeactivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.DeactivateConfig(&_MercuryVerifier.TransactOpts, feedId, configDigest)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) DeactivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.DeactivateConfig(&_MercuryVerifier.TransactOpts, feedId, configDigest)
}

func (_MercuryVerifier *MercuryVerifierTransactor) SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "setConfig", feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_MercuryVerifier *MercuryVerifierSession) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.SetConfig(&_MercuryVerifier.TransactOpts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.SetConfig(&_MercuryVerifier.TransactOpts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_MercuryVerifier *MercuryVerifierTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "transferOwnership", to)
}

func (_MercuryVerifier *MercuryVerifierSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.TransferOwnership(&_MercuryVerifier.TransactOpts, to)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.TransferOwnership(&_MercuryVerifier.TransactOpts, to)
}

func (_MercuryVerifier *MercuryVerifierTransactor) Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "verify", signedReport, sender)
}

func (_MercuryVerifier *MercuryVerifierSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.Verify(&_MercuryVerifier.TransactOpts, signedReport, sender)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.Verify(&_MercuryVerifier.TransactOpts, signedReport, sender)
}

type MercuryVerifierConfigActivatedIterator struct {
	Event *MercuryVerifierConfigActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierConfigActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierConfigActivated)
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
		it.Event = new(MercuryVerifierConfigActivated)
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

func (it *MercuryVerifierConfigActivatedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierConfigActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierConfigActivated struct {
	FeedId       [32]byte
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierConfigActivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ConfigActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierConfigActivatedIterator{contract: _MercuryVerifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigActivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ConfigActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierConfigActivated)
				if err := _MercuryVerifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseConfigActivated(log types.Log) (*MercuryVerifierConfigActivated, error) {
	event := new(MercuryVerifierConfigActivated)
	if err := _MercuryVerifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierConfigDeactivatedIterator struct {
	Event *MercuryVerifierConfigDeactivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierConfigDeactivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierConfigDeactivated)
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
		it.Event = new(MercuryVerifierConfigDeactivated)
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

func (it *MercuryVerifierConfigDeactivatedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierConfigDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierConfigDeactivated struct {
	FeedId       [32]byte
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterConfigDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierConfigDeactivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ConfigDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierConfigDeactivatedIterator{contract: _MercuryVerifier.contract, event: "ConfigDeactivated", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigDeactivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ConfigDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierConfigDeactivated)
				if err := _MercuryVerifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseConfigDeactivated(log types.Log) (*MercuryVerifierConfigDeactivated, error) {
	event := new(MercuryVerifierConfigDeactivated)
	if err := _MercuryVerifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierConfigSetIterator struct {
	Event *MercuryVerifierConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierConfigSet)
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
		it.Event = new(MercuryVerifierConfigSet)
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

func (it *MercuryVerifierConfigSetIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierConfigSet struct {
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
	Raw                       types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterConfigSet(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierConfigSetIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ConfigSet", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierConfigSetIterator{contract: _MercuryVerifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigSet, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ConfigSet", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierConfigSet)
				if err := _MercuryVerifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseConfigSet(log types.Log) (*MercuryVerifierConfigSet, error) {
	event := new(MercuryVerifierConfigSet)
	if err := _MercuryVerifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierOwnershipTransferRequestedIterator struct {
	Event *MercuryVerifierOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierOwnershipTransferRequested)
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
		it.Event = new(MercuryVerifierOwnershipTransferRequested)
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

func (it *MercuryVerifierOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierOwnershipTransferRequestedIterator{contract: _MercuryVerifier.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierOwnershipTransferRequested)
				if err := _MercuryVerifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseOwnershipTransferRequested(log types.Log) (*MercuryVerifierOwnershipTransferRequested, error) {
	event := new(MercuryVerifierOwnershipTransferRequested)
	if err := _MercuryVerifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierOwnershipTransferredIterator struct {
	Event *MercuryVerifierOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierOwnershipTransferred)
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
		it.Event = new(MercuryVerifierOwnershipTransferred)
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

func (it *MercuryVerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierOwnershipTransferredIterator{contract: _MercuryVerifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierOwnershipTransferred)
				if err := _MercuryVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseOwnershipTransferred(log types.Log) (*MercuryVerifierOwnershipTransferred, error) {
	event := new(MercuryVerifierOwnershipTransferred)
	if err := _MercuryVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierReportVerifiedIterator struct {
	Event *MercuryVerifierReportVerified

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierReportVerifiedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierReportVerified)
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
		it.Event = new(MercuryVerifierReportVerified)
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

func (it *MercuryVerifierReportVerifiedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierReportVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierReportVerified struct {
	FeedId     [32]byte
	ReportHash [32]byte
	Requester  common.Address
	Raw        types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierReportVerifiedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierReportVerifiedIterator{contract: _MercuryVerifier.contract, event: "ReportVerified", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchReportVerified(opts *bind.WatchOpts, sink chan<- *MercuryVerifierReportVerified, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierReportVerified)
				if err := _MercuryVerifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseReportVerified(log types.Log) (*MercuryVerifierReportVerified, error) {
	event := new(MercuryVerifierReportVerified)
	if err := _MercuryVerifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}

func (_MercuryVerifier *MercuryVerifier) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MercuryVerifier.abi.Events["ConfigActivated"].ID:
		return _MercuryVerifier.ParseConfigActivated(log)
	case _MercuryVerifier.abi.Events["ConfigDeactivated"].ID:
		return _MercuryVerifier.ParseConfigDeactivated(log)
	case _MercuryVerifier.abi.Events["ConfigSet"].ID:
		return _MercuryVerifier.ParseConfigSet(log)
	case _MercuryVerifier.abi.Events["OwnershipTransferRequested"].ID:
		return _MercuryVerifier.ParseOwnershipTransferRequested(log)
	case _MercuryVerifier.abi.Events["OwnershipTransferred"].ID:
		return _MercuryVerifier.ParseOwnershipTransferred(log)
	case _MercuryVerifier.abi.Events["ReportVerified"].ID:
		return _MercuryVerifier.ParseReportVerified(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MercuryVerifierConfigActivated) Topic() common.Hash {
	return common.HexToHash("0x54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe7746")
}

func (MercuryVerifierConfigDeactivated) Topic() common.Hash {
	return common.HexToHash("0x0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c")
}

func (MercuryVerifierConfigSet) Topic() common.Hash {
	return common.HexToHash("0xa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da")
}

func (MercuryVerifierOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MercuryVerifierOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MercuryVerifierReportVerified) Topic() common.Hash {
	return common.HexToHash("0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269")
}

func (_MercuryVerifier *MercuryVerifier) Address() common.Address {
	return _MercuryVerifier.address
}

type MercuryVerifierInterface interface {
	LatestConfigDetails(opts *bind.CallOpts, feedId [32]byte) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts, feedId [32]byte) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ActivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error)

	DeactivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error)

	FilterConfigActivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierConfigActivatedIterator, error)

	WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigActivated, feedId [][32]byte) (event.Subscription, error)

	ParseConfigActivated(log types.Log) (*MercuryVerifierConfigActivated, error)

	FilterConfigDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierConfigDeactivatedIterator, error)

	WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigDeactivated, feedId [][32]byte) (event.Subscription, error)

	ParseConfigDeactivated(log types.Log) (*MercuryVerifierConfigDeactivated, error)

	FilterConfigSet(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigSet, feedId [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*MercuryVerifierConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MercuryVerifierOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MercuryVerifierOwnershipTransferred, error)

	FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierReportVerifiedIterator, error)

	WatchReportVerified(opts *bind.WatchOpts, sink chan<- *MercuryVerifierReportVerified, feedId [][32]byte) (event.Subscription, error)

	ParseReportVerified(log types.Log) (*MercuryVerifierReportVerified, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
