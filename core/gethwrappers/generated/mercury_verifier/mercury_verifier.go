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

var MercuryVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InactiveFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"activateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"deactivateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002059380380620020598339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611e5e620001fb6000396000818161031d015261091d0152611e5e6000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c806379ba50971161008c578063b70d929d11610066578063b70d929d14610221578063ded6307c14610280578063e84f128e14610293578063f2fde38b146102f057600080fd5b806379ba5097146101de5780638da5cb5b146101e657806394d959801461020e57600080fd5b80633dd86430116100bd5780633dd86430146101a357806344a0b2ad146101b8578063564a0a7a146101cb57600080fd5b806301ffc9a7146100e4578063181f5a771461014e5780633d3ac1b514610190575b600080fd5b6101396100f23660046114d4565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b60408051808201909152600e81527f566572696669657220302e302e3200000000000000000000000000000000000060208201525b6040516101459190611581565b61018361019e3660046115bd565b610303565b6101b66101b136600461163e565b6104a2565b005b6101b66101c63660046118a6565b610554565b6101b66101d936600461163e565b610a44565b6101b6610b05565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610145565b6101b661021c36600461197e565b610c02565b61025d61022f36600461163e565b6000908152600260205260408120600181015490549192909168010000000000000000900463ffffffff1690565b604080519315158452602084019290925263ffffffff1690820152606001610145565b6101b661028e36600461197e565b610d63565b6102cd6102a136600461163e565b6000908152600260205260409020805460019091015463ffffffff808316936401000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610145565b6101b66102fe3660046119a0565b610e74565b60603373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610374576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080610386888a018a6119bb565b9450945094509450945060008461039c90611a96565b60008181526002602052604090208054919250906c01000000000000000000000000900460ff1615610402576040517f36dbe748000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b8651600081815260028301602052604090206104218483898985610e88565b61042b8984610fcc565b87516020890120610440818b8a8a8a87611034565b6040805182815273ffffffffffffffffffffffffffffffffffffffff8e16602082015286917f5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269910160405180910390a250969c9b505050505050505050505050565b6104aa6112b0565b60008181526002602052604081208054909163ffffffff90911690036104ff576040517fa25b0b96000000000000000000000000000000000000000000000000000000008152600481018390526024016103f9565b80547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff16815560405182907ff438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b490600090a25050565b85518460ff1680600003610594576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f8211156105d9576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044016103f9565b6105e4816003611b39565b821161063c57816105f6826003611b39565b610601906001611b76565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260248201526044016103f9565b6106446112b0565b60008981526002602052604081208054909163ffffffff90911690829061066a83611b8f565b82546101009290920a63ffffffff81810219909316918316021790915582546000925061069f918d91168c8c8c8c8c8c611333565b6000818152600284016020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660ff8c16176101001790559091505b8a518160ff1610156108e05760008b8260ff168151811061070557610705611adb565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610775576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452600190810190925290912054610100900460ff16908111156107c8576107c8611bb2565b1480159150610803576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff8416815260208101600190526000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684526001908101835292208351815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00821681178355928501519193919284927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090921617906101009084908111156108c5576108c5611bb2565b02179055509050505050806108d990611be1565b90506106e2565b5060018201546040517f2cc994770000000000000000000000000000000000000000000000000000000081526004810191909152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cc9947790604401600060405180830381600087803b15801561097657600080fd5b505af115801561098a573d6000803e3d6000fd5b505050508a7fa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da8360000160049054906101000a900463ffffffff16838560000160009054906101000a900463ffffffff168e8e8e8e8e8e6040516109f699989796959493929190611c81565b60405180910390a281547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1664010000000063ffffffff431602178255600190910155505050505050505050565b610a4c6112b0565b60008181526002602052604081208054909163ffffffff9091169003610aa1576040517fa25b0b96000000000000000000000000000000000000000000000000000000008152600481018390526024016103f9565b80547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff166c0100000000000000000000000017815560405182907ffc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd06190600090a25050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b86576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103f9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610c0a6112b0565b600082815260026020526040902081610c4f576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604081205460ff169003610ca5576040517f8bca631100000000000000000000000000000000000000000000000000000000815260048101849052602481018390526044016103f9565b80600101548203610cec576040517fa403c01600000000000000000000000000000000000000000000000000000000815260048101849052602481018390526044016103f9565b60008281526002820160205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555183907f0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c90610d569085815260200190565b60405180910390a2505050565b610d6b6112b0565b600082815260026020526040902081610db0576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604081205460ff169003610e06576040517f8bca631100000000000000000000000000000000000000000000000000000000815260048101849052602481018390526044016103f9565b60008281526002820160205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555183907f54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe774690610d569085815260200190565b610e7c6112b0565b610e85816113df565b50565b8054600090610e9b9060ff166001611d17565b825490915060ff16600003610ee6576040517f8bca631100000000000000000000000000000000000000000000000000000000815260048101879052602481018690526044016103f9565b8154610100900460ff16610f30576040517ffc10a28300000000000000000000000000000000000000000000000000000000815260048101879052602481018690526044016103f9565b8060ff16845114610f7c5783516040517f5348a282000000000000000000000000000000000000000000000000000000008152600481019190915260ff821660248201526044016103f9565b8251845114610fc457835183516040517ff0d31408000000000000000000000000000000000000000000000000000000008152600481019290925260248201526044016103f9565b505050505050565b6020820151815463ffffffff600883901c8116916801000000000000000090041681111561102e5782547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff166801000000000000000063ffffffff8316021783555b50505050565b60008686604051602001611049929190611d30565b604051602081830303815290604052805190602001209050600061107d604080518082019091526000808252602082015290565b8651600090815b81811015611248576001868983602081106110a1576110a1611adb565b6110ae91901a601b611d17565b8c84815181106110c0576110c0611adb565b60200260200101518c85815181106110da576110da611adb565b602002602001015160405160008152602001604052604051611118949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561113a573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526001808d01602090815291859020848601909552845460ff8082168652939950939550908501926101009004909116908111156111bf576111bf611bb2565b60018111156111d0576111d0611bb2565b90525093506001846020015160018111156111ed576111ed611bb2565b14611224576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836000015160080260ff166001901b850194508061124190611d6c565b9050611084565b50837e010101010101010101010101010101010101010101010101010101010101018516146112a3576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611331576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103f9565b565b6000808946308b8b8b8b8b8b8b6040516020016113599a99989796959493929190611da4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e010000000000000000000000000000000000000000000000000000000000001791505098975050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361145e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103f9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156114e657600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461151657600080fd5b9392505050565b6000815180845260005b8181101561154357602081850181015186830182015201611527565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611516602083018461151d565b803573ffffffffffffffffffffffffffffffffffffffff811681146115b857600080fd5b919050565b6000806000604084860312156115d257600080fd5b833567ffffffffffffffff808211156115ea57600080fd5b818601915086601f8301126115fe57600080fd5b81358181111561160d57600080fd5b87602082850101111561161f57600080fd5b6020928301955093506116359186019050611594565b90509250925092565b60006020828403121561165057600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156116a9576116a9611657565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156116f6576116f6611657565b604052919050565b600067ffffffffffffffff82111561171857611718611657565b5060051b60200190565b600082601f83011261173357600080fd5b81356020611748611743836116fe565b6116af565b82815260059290921b8401810191818101908684111561176757600080fd5b8286015b848110156117895761177c81611594565b835291830191830161176b565b509695505050505050565b600082601f8301126117a557600080fd5b813560206117b5611743836116fe565b82815260059290921b840181019181810190868411156117d457600080fd5b8286015b8481101561178957803583529183019183016117d8565b803560ff811681146115b857600080fd5b600082601f83011261181157600080fd5b813567ffffffffffffffff81111561182b5761182b611657565b61185c60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016116af565b81815284602083860101111561187157600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff811681146115b857600080fd5b600080600080600080600060e0888a0312156118c157600080fd5b87359650602088013567ffffffffffffffff808211156118e057600080fd5b6118ec8b838c01611722565b975060408a013591508082111561190257600080fd5b61190e8b838c01611794565b965061191c60608b016117ef565b955060808a013591508082111561193257600080fd5b61193e8b838c01611800565b945061194c60a08b0161188e565b935060c08a013591508082111561196257600080fd5b5061196f8a828b01611800565b91505092959891949750929550565b6000806040838503121561199157600080fd5b50508035926020909101359150565b6000602082840312156119b257600080fd5b61151682611594565b600080600080600060e086880312156119d357600080fd5b86601f8701126119e257600080fd5b6119ea611686565b8060608801898111156119fc57600080fd5b885b81811015611a165780358452602093840193016119fe565b5090965035905067ffffffffffffffff80821115611a3357600080fd5b611a3f89838a01611800565b95506080880135915080821115611a5557600080fd5b611a6189838a01611794565b945060a0880135915080821115611a7757600080fd5b50611a8488828901611794565b9598949750929560c001359392505050565b80516020808301519190811015611ad5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611b7157611b71611b0a565b500290565b80820180821115611b8957611b89611b0a565b92915050565b600063ffffffff808316818103611ba857611ba8611b0a565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600060ff821660ff8103611bf757611bf7611b0a565b60010192915050565b600081518084526020808501945080840160005b83811015611c4657815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611c14565b509495945050505050565b600081518084526020808501945080840160005b83811015611c4657815187529582019590820190600101611c65565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611cb18184018a611c00565b90508281036080840152611cc58189611c51565b905060ff871660a084015282810360c0840152611ce2818761151d565b905067ffffffffffffffff851660e0840152828103610100840152611d07818561151d565b9c9b505050505050505050505050565b60ff8181168382160190811115611b8957611b89611b0a565b828152600060208083018460005b6003811015611d5b57815183529183019190830190600101611d3e565b505050506080820190509392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611d9d57611d9d611b0a565b5060010190565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b166060850152816080850152611df18285018b611c00565b915083820360a0850152611e05828a611c51565b915060ff881660c085015283820360e0850152611e22828861151d565b9086166101008501528381036101208501529050611e40818561151d565b9d9c5050505050505050505050505056fea164736f6c6343000810000a",
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
	parsed, err := MercuryVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

func (_MercuryVerifier *MercuryVerifierTransactor) ActivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "activateFeed", feedId)
}

func (_MercuryVerifier *MercuryVerifierSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.ActivateFeed(&_MercuryVerifier.TransactOpts, feedId)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.ActivateFeed(&_MercuryVerifier.TransactOpts, feedId)
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

func (_MercuryVerifier *MercuryVerifierTransactor) DeactivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "deactivateFeed", feedId)
}

func (_MercuryVerifier *MercuryVerifierSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.DeactivateFeed(&_MercuryVerifier.TransactOpts, feedId)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.DeactivateFeed(&_MercuryVerifier.TransactOpts, feedId)
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

type MercuryVerifierFeedActivatedIterator struct {
	Event *MercuryVerifierFeedActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierFeedActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierFeedActivated)
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
		it.Event = new(MercuryVerifierFeedActivated)
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

func (it *MercuryVerifierFeedActivatedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierFeedActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierFeedActivated struct {
	FeedId [32]byte
	Raw    types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterFeedActivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierFeedActivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "FeedActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierFeedActivatedIterator{contract: _MercuryVerifier.contract, event: "FeedActivated", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchFeedActivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierFeedActivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "FeedActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierFeedActivated)
				if err := _MercuryVerifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseFeedActivated(log types.Log) (*MercuryVerifierFeedActivated, error) {
	event := new(MercuryVerifierFeedActivated)
	if err := _MercuryVerifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierFeedDeactivatedIterator struct {
	Event *MercuryVerifierFeedDeactivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierFeedDeactivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierFeedDeactivated)
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
		it.Event = new(MercuryVerifierFeedDeactivated)
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

func (it *MercuryVerifierFeedDeactivatedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierFeedDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierFeedDeactivated struct {
	FeedId [32]byte
	Raw    types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterFeedDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierFeedDeactivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "FeedDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierFeedDeactivatedIterator{contract: _MercuryVerifier.contract, event: "FeedDeactivated", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchFeedDeactivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierFeedDeactivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "FeedDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierFeedDeactivated)
				if err := _MercuryVerifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
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

func (_MercuryVerifier *MercuryVerifierFilterer) ParseFeedDeactivated(log types.Log) (*MercuryVerifierFeedDeactivated, error) {
	event := new(MercuryVerifierFeedDeactivated)
	if err := _MercuryVerifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
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
	case _MercuryVerifier.abi.Events["FeedActivated"].ID:
		return _MercuryVerifier.ParseFeedActivated(log)
	case _MercuryVerifier.abi.Events["FeedDeactivated"].ID:
		return _MercuryVerifier.ParseFeedDeactivated(log)
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

func (MercuryVerifierFeedActivated) Topic() common.Hash {
	return common.HexToHash("0xf438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b4")
}

func (MercuryVerifierFeedDeactivated) Topic() common.Hash {
	return common.HexToHash("0xfc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd061")
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

	ActivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error)

	DeactivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error)

	DeactivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error)

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

	FilterFeedActivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierFeedActivatedIterator, error)

	WatchFeedActivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierFeedActivated, feedId [][32]byte) (event.Subscription, error)

	ParseFeedActivated(log types.Log) (*MercuryVerifierFeedActivated, error)

	FilterFeedDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*MercuryVerifierFeedDeactivatedIterator, error)

	WatchFeedDeactivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierFeedDeactivated, feedId [][32]byte) (event.Subscription, error)

	ParseFeedDeactivated(log types.Log) (*MercuryVerifierFeedDeactivated, error)

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
