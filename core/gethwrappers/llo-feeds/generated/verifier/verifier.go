// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier

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

type CommonAddressAndWeight struct {
	Addr   common.Address
	Weight uint64
}

var VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadVerification\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InactiveFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"activateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"deactivateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sourceAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newConfigCount\",\"type\":\"uint32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setConfigFromSource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"verifierResponse\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620023d5380380620023d58339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b6080516121da620001fb6000396000818161033b015261131501526121da6000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806394d959801161008c578063e7db9c2a11610066578063e7db9c2a1461028b578063e84f128e1461029e578063f0107221146102fb578063f2fde38b1461030e57600080fd5b806394d9598014610206578063b70d929d14610219578063ded6307c1461027857600080fd5b80633dd86430116100c85780633dd86430146101ae578063564a0a7a146101c357806379ba5097146101d65780638da5cb5b146101de57600080fd5b806301ffc9a7146100ef578063181f5a77146101595780633d3ac1b51461019b575b600080fd5b6101446100fd3660046115da565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b60408051808201909152600e81527f566572696669657220312e322e3000000000000000000000000000000000000060208201525b6040516101509190611687565b61018e6101a93660046116c3565b610321565b6101c16101bc366004611744565b6104bb565b005b6101c16101d1366004611744565b61056d565b6101c161062e565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610150565b6101c161021436600461175d565b61072b565b610255610227366004611744565b6000908152600260205260408120600181015490549192909168010000000000000000900463ffffffff1690565b604080519315158452602084019290925263ffffffff1690820152606001610150565b6101c161028636600461175d565b61088c565b6101c1610299366004611a92565b61099d565b6102d86102ac366004611744565b6000908152600260205260409020805460019091015463ffffffff808316936401000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610150565b6101c1610309366004611bc5565b610aaf565b6101c161031c366004611cc2565b610b79565b60603373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610392576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080806103a4888a018a611cdd565b945094509450945094506000846103ba90611db8565b60008181526002602052604090208054919250906c01000000000000000000000000900460ff1615610420576040517f36dbe748000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b86516000818152600283016020526040902061043f8483898985610b8d565b6104498984610c89565b8751602089012061045e818b8a8a8a87610cf1565b60405173ffffffffffffffffffffffffffffffffffffffff8d16815285907f58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc59060200160405180910390a250969c9b505050505050505050505050565b6104c3610f6d565b60008181526002602052604081208054909163ffffffff9091169003610518576040517fa25b0b9600000000000000000000000000000000000000000000000000000000815260048101839052602401610417565b80547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff16815560405182907ff438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b490600090a25050565b610575610f6d565b60008181526002602052604081208054909163ffffffff90911690036105ca576040517fa25b0b9600000000000000000000000000000000000000000000000000000000815260048101839052602401610417565b80547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff166c0100000000000000000000000017815560405182907ffc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd06190600090a25050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146106af576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610417565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610733610f6d565b600082815260026020526040902081610778576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604081205460ff1690036107ce576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610417565b80600101548203610815576040517fa403c0160000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610417565b60008281526002820160205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555183907f0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c9061087f9085815260200190565b60405180910390a2505050565b610894610f6d565b6000828152600260205260409020816108d9576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604081205460ff16900361092f576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610417565b60008281526002820160205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555183907f54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe77469061087f9085815260200190565b86518560ff16806000036109dd576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610a22576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f6024820152604401610417565b610a2d816003611e5b565b8211610a855781610a3f826003611e5b565b610a4a906001611e78565b6040517f9dd9e6d800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610417565b610a8d610f6d565b610aa08d8d8d8d8d8d8d8d8d8d8d610ff0565b50505050505050505050505050565b86518560ff1680600003610aef576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610b34576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f6024820152604401610417565b610b3f816003611e5b565b8211610b515781610a3f826003611e5b565b610b59610f6d565b610b6d8a463060008d8d8d8d8d8d8d610ff0565b50505050505050505050565b610b81610f6d565b610b8a81611437565b50565b8054600090610ba09060ff166001611e8b565b8254909150610100900460ff16610bed576040517ffc10a2830000000000000000000000000000000000000000000000000000000081526004810187905260248101869052604401610417565b8060ff16845114610c395783516040517f5348a282000000000000000000000000000000000000000000000000000000008152600481019190915260ff82166024820152604401610417565b8251845114610c8157835183516040517ff0d3140800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610417565b505050505050565b6020820151815463ffffffff600883901c81169168010000000000000000900416811115610ceb5782547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff166801000000000000000063ffffffff8316021783555b50505050565b60008686604051602001610d06929190611ea4565b6040516020818303038152906040528051906020012090506000610d3a604080518082019091526000808252602082015290565b8651600090815b81811015610f0557600186898360208110610d5e57610d5e611dfd565b610d6b91901a601b611e8b565b8c8481518110610d7d57610d7d611dfd565b60200260200101518c8581518110610d9757610d97611dfd565b602002602001015160405160008152602001604052604051610dd5949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610df7573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526001808d01602090815291859020848601909552845460ff808216865293995093955090850192610100900490911690811115610e7c57610e7c611ee0565b6001811115610e8d57610e8d611ee0565b9052509350600184602001516001811115610eaa57610eaa611ee0565b14610ee1576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836000015160080260ff166001901b8501945080610efe90611f0f565b9050610d41565b50837e01010101010101010101010101010101010101010101010101010101010101851614610f60576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610fee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610417565b565b60008b815260026020526040902063ffffffff89161561103d5780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8a16178155611071565b805463ffffffff1681600061105183611f47565b91906101000a81548163ffffffff021916908363ffffffff160217905550505b8054600090611091908e908e908e9063ffffffff168d8d8d8d8d8d61152c565b6000818152600284016020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660ff8b16176101001790559091505b89518160ff1610156112d25760008a8260ff16815181106110f7576110f7611dfd565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611167576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452600190810190925290912054610100900460ff16908111156111ba576111ba611ee0565b14801591506111f5576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff8416815260208101600190526000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684526001908101835292208351815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00821681178355928501519193919284927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090921617906101009084908111156112b7576112b7611ee0565b02179055509050505050806112cb90611f6a565b90506110d4565b5060018201546040517fb011b24700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163b011b2479161134d919085908890600401611f89565b600060405180830381600087803b15801561136757600080fd5b505af115801561137b573d6000803e3d6000fd5b505050508c7fa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da8360000160049054906101000a900463ffffffff16838560000160009054906101000a900463ffffffff168d8d8d8d8d8d6040516113e79998979695949392919061208a565b60405180910390a281547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1664010000000063ffffffff4316021782556001909101555050505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036114b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610417565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808b8b8b8b8b8b8b8b8b8b6040516020016115529a99989796959493929190612120565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e06000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b6000602082840312156115ec57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461161c57600080fd5b9392505050565b6000815180845260005b818110156116495760208185018101518683018201520161162d565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061161c6020830184611623565b803573ffffffffffffffffffffffffffffffffffffffff811681146116be57600080fd5b919050565b6000806000604084860312156116d857600080fd5b833567ffffffffffffffff808211156116f057600080fd5b818601915086601f83011261170457600080fd5b81358181111561171357600080fd5b87602082850101111561172557600080fd5b60209283019550935061173b918601905061169a565b90509250925092565b60006020828403121561175657600080fd5b5035919050565b6000806040838503121561177057600080fd5b50508035926020909101359150565b803563ffffffff811681146116be57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156117e5576117e5611793565b60405290565b6040516060810167ffffffffffffffff811182821017156117e5576117e5611793565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561185557611855611793565b604052919050565b600067ffffffffffffffff82111561187757611877611793565b5060051b60200190565b600082601f83011261189257600080fd5b813560206118a76118a28361185d565b61180e565b82815260059290921b840181019181810190868411156118c657600080fd5b8286015b848110156118e8576118db8161169a565b83529183019183016118ca565b509695505050505050565b600082601f83011261190457600080fd5b813560206119146118a28361185d565b82815260059290921b8401810191818101908684111561193357600080fd5b8286015b848110156118e85780358352918301918301611937565b803560ff811681146116be57600080fd5b600082601f83011261197057600080fd5b813567ffffffffffffffff81111561198a5761198a611793565b6119bb60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161180e565b8181528460208386010111156119d057600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff811681146116be57600080fd5b600082601f830112611a1657600080fd5b81356020611a266118a28361185d565b82815260069290921b84018101918181019086841115611a4557600080fd5b8286015b848110156118e85760408189031215611a625760008081fd5b611a6a6117c2565b611a738261169a565b8152611a808583016119ed565b81860152835291830191604001611a49565b60008060008060008060008060008060006101608c8e031215611ab457600080fd5b8b359a5060208c01359950611acb60408d0161169a565b9850611ad960608d0161177f565b975067ffffffffffffffff8060808e01351115611af557600080fd5b611b058e60808f01358f01611881565b97508060a08e01351115611b1857600080fd5b611b288e60a08f01358f016118f3565b9650611b3660c08e0161194e565b95508060e08e01351115611b4957600080fd5b611b598e60e08f01358f0161195f565b9450611b686101008e016119ed565b9350806101208e01351115611b7c57600080fd5b611b8d8e6101208f01358f0161195f565b9250806101408e01351115611ba157600080fd5b50611bb38d6101408e01358e01611a05565b90509295989b509295989b9093969950565b600080600080600080600080610100898b031215611be257600080fd5b88359750602089013567ffffffffffffffff80821115611c0157600080fd5b611c0d8c838d01611881565b985060408b0135915080821115611c2357600080fd5b611c2f8c838d016118f3565b9750611c3d60608c0161194e565b965060808b0135915080821115611c5357600080fd5b611c5f8c838d0161195f565b9550611c6d60a08c016119ed565b945060c08b0135915080821115611c8357600080fd5b611c8f8c838d0161195f565b935060e08b0135915080821115611ca557600080fd5b50611cb28b828c01611a05565b9150509295985092959890939650565b600060208284031215611cd457600080fd5b61161c8261169a565b600080600080600060e08688031215611cf557600080fd5b86601f870112611d0457600080fd5b611d0c6117eb565b806060880189811115611d1e57600080fd5b885b81811015611d38578035845260209384019301611d20565b5090965035905067ffffffffffffffff80821115611d5557600080fd5b611d6189838a0161195f565b95506080880135915080821115611d7757600080fd5b611d8389838a016118f3565b945060a0880135915080821115611d9957600080fd5b50611da6888289016118f3565b9598949750929560c001359392505050565b80516020808301519190811015611df7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417611e7257611e72611e2c565b92915050565b80820180821115611e7257611e72611e2c565b60ff8181168382160190811115611e7257611e72611e2c565b828152600060208083018460005b6003811015611ecf57815183529183019190830190600101611eb2565b505050506080820190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611f4057611f40611e2c565b5060010190565b600063ffffffff808316818103611f6057611f60611e2c565b6001019392505050565b600060ff821660ff8103611f8057611f80611e2c565b60010192915050565b600060608201858352602085818501526040606081860152828651808552608087019150838801945060005b81811015611ffa578551805173ffffffffffffffffffffffffffffffffffffffff16845285015167ffffffffffffffff16858401529484019491830191600101611fb5565b50909998505050505050505050565b600081518084526020808501945080840160005b8381101561204f57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161201d565b509495945050505050565b600081518084526020808501945080840160005b8381101561204f5781518752958201959082019060010161206e565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526120ba8184018a612009565b905082810360808401526120ce818961205a565b905060ff871660a084015282810360c08401526120eb8187611623565b905067ffffffffffffffff851660e08401528281036101008401526121108185611623565b9c9b505050505050505050505050565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b16606085015281608085015261216d8285018b612009565b915083820360a0850152612181828a61205a565b915060ff881660c085015283820360e085015261219e8288611623565b90861661010085015283810361012085015290506121bc8185611623565b9d9c5050505050505050505050505056fea164736f6c6343000813000a",
}

var VerifierABI = VerifierMetaData.ABI

var VerifierBin = VerifierMetaData.Bin

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
	return address, tx, &Verifier{address: address, abi: *parsed, VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

type Verifier struct {
	address common.Address
	abi     abi.ABI
	VerifierCaller
	VerifierTransactor
	VerifierFilterer
}

type VerifierCaller struct {
	contract *bind.BoundContract
}

type VerifierTransactor struct {
	contract *bind.BoundContract
}

type VerifierFilterer struct {
	contract *bind.BoundContract
}

type VerifierSession struct {
	Contract     *Verifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifierCallerSession struct {
	Contract *VerifierCaller
	CallOpts bind.CallOpts
}

type VerifierTransactorSession struct {
	Contract     *VerifierTransactor
	TransactOpts bind.TransactOpts
}

type VerifierRaw struct {
	Contract *Verifier
}

type VerifierCallerRaw struct {
	Contract *VerifierCaller
}

type VerifierTransactorRaw struct {
	Contract *VerifierTransactor
}

func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	abi, err := abi.JSON(strings.NewReader(VerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{address: address, abi: abi, VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

func (_Verifier *VerifierCaller) LatestConfigDetails(opts *bind.CallOpts, feedId [32]byte) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDetails", feedId)

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_Verifier *VerifierSession) LatestConfigDetails(feedId [32]byte) (LatestConfigDetails,

	error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts, feedId)
}

func (_Verifier *VerifierCallerSession) LatestConfigDetails(feedId [32]byte) (LatestConfigDetails,

	error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts, feedId)
}

func (_Verifier *VerifierCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts, feedId [32]byte) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDigestAndEpoch", feedId)

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_Verifier *VerifierSession) LatestConfigDigestAndEpoch(feedId [32]byte) (LatestConfigDigestAndEpoch,

	error) {
	return _Verifier.Contract.LatestConfigDigestAndEpoch(&_Verifier.CallOpts, feedId)
}

func (_Verifier *VerifierCallerSession) LatestConfigDigestAndEpoch(feedId [32]byte) (LatestConfigDigestAndEpoch,

	error) {
	return _Verifier.Contract.LatestConfigDigestAndEpoch(&_Verifier.CallOpts, feedId)
}

func (_Verifier *VerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Verifier *VerifierSession) Owner() (common.Address, error) {
	return _Verifier.Contract.Owner(&_Verifier.CallOpts)
}

func (_Verifier *VerifierCallerSession) Owner() (common.Address, error) {
	return _Verifier.Contract.Owner(&_Verifier.CallOpts)
}

func (_Verifier *VerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Verifier *VerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Verifier.Contract.SupportsInterface(&_Verifier.CallOpts, interfaceId)
}

func (_Verifier *VerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Verifier.Contract.SupportsInterface(&_Verifier.CallOpts, interfaceId)
}

func (_Verifier *VerifierCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Verifier *VerifierSession) TypeAndVersion() (string, error) {
	return _Verifier.Contract.TypeAndVersion(&_Verifier.CallOpts)
}

func (_Verifier *VerifierCallerSession) TypeAndVersion() (string, error) {
	return _Verifier.Contract.TypeAndVersion(&_Verifier.CallOpts)
}

func (_Verifier *VerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "acceptOwnership")
}

func (_Verifier *VerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _Verifier.Contract.AcceptOwnership(&_Verifier.TransactOpts)
}

func (_Verifier *VerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Verifier.Contract.AcceptOwnership(&_Verifier.TransactOpts)
}

func (_Verifier *VerifierTransactor) ActivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "activateConfig", feedId, configDigest)
}

func (_Verifier *VerifierSession) ActivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

func (_Verifier *VerifierTransactorSession) ActivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

func (_Verifier *VerifierTransactor) ActivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "activateFeed", feedId)
}

func (_Verifier *VerifierSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateFeed(&_Verifier.TransactOpts, feedId)
}

func (_Verifier *VerifierTransactorSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateFeed(&_Verifier.TransactOpts, feedId)
}

func (_Verifier *VerifierTransactor) DeactivateConfig(opts *bind.TransactOpts, feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "deactivateConfig", feedId, configDigest)
}

func (_Verifier *VerifierSession) DeactivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

func (_Verifier *VerifierTransactorSession) DeactivateConfig(feedId [32]byte, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, feedId, configDigest)
}

func (_Verifier *VerifierTransactor) DeactivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "deactivateFeed", feedId)
}

func (_Verifier *VerifierSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateFeed(&_Verifier.TransactOpts, feedId)
}

func (_Verifier *VerifierTransactorSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateFeed(&_Verifier.TransactOpts, feedId)
}

func (_Verifier *VerifierTransactor) SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "setConfig", feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_Verifier *VerifierSession) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_Verifier *VerifierTransactorSession) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_Verifier *VerifierTransactor) SetConfigFromSource(opts *bind.TransactOpts, feedId [32]byte, sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "setConfigFromSource", feedId, sourceChainId, sourceAddress, newConfigCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_Verifier *VerifierSession) SetConfigFromSource(feedId [32]byte, sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfigFromSource(&_Verifier.TransactOpts, feedId, sourceChainId, sourceAddress, newConfigCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_Verifier *VerifierTransactorSession) SetConfigFromSource(feedId [32]byte, sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfigFromSource(&_Verifier.TransactOpts, feedId, sourceChainId, sourceAddress, newConfigCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_Verifier *VerifierTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "transferOwnership", to)
}

func (_Verifier *VerifierSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.TransferOwnership(&_Verifier.TransactOpts, to)
}

func (_Verifier *VerifierTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.TransferOwnership(&_Verifier.TransactOpts, to)
}

func (_Verifier *VerifierTransactor) Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "verify", signedReport, sender)
}

func (_Verifier *VerifierSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.Verify(&_Verifier.TransactOpts, signedReport, sender)
}

func (_Verifier *VerifierTransactorSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.Verify(&_Verifier.TransactOpts, signedReport, sender)
}

type VerifierConfigActivatedIterator struct {
	Event *VerifierConfigActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierConfigActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierConfigActivatedIterator) Error() error {
	return it.fail
}

func (it *VerifierConfigActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierConfigActivated struct {
	FeedId       [32]byte
	ConfigDigest [32]byte
	Raw          types.Log
}

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

func (_Verifier *VerifierFilterer) ParseConfigActivated(log types.Log) (*VerifierConfigActivated, error) {
	event := new(VerifierConfigActivated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierConfigDeactivatedIterator struct {
	Event *VerifierConfigDeactivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierConfigDeactivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierConfigDeactivatedIterator) Error() error {
	return it.fail
}

func (it *VerifierConfigDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierConfigDeactivated struct {
	FeedId       [32]byte
	ConfigDigest [32]byte
	Raw          types.Log
}

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

func (_Verifier *VerifierFilterer) ParseConfigDeactivated(log types.Log) (*VerifierConfigDeactivated, error) {
	event := new(VerifierConfigDeactivated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierConfigSetIterator struct {
	Event *VerifierConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierConfigSetIterator) Error() error {
	return it.fail
}

func (it *VerifierConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

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
	Raw                       types.Log
}

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

func (_Verifier *VerifierFilterer) ParseConfigSet(log types.Log) (*VerifierConfigSet, error) {
	event := new(VerifierConfigSet)
	if err := _Verifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierFeedActivatedIterator struct {
	Event *VerifierFeedActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierFeedActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierFeedActivatedIterator) Error() error {
	return it.fail
}

func (it *VerifierFeedActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierFeedActivated struct {
	FeedId [32]byte
	Raw    types.Log
}

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

func (_Verifier *VerifierFilterer) ParseFeedActivated(log types.Log) (*VerifierFeedActivated, error) {
	event := new(VerifierFeedActivated)
	if err := _Verifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierFeedDeactivatedIterator struct {
	Event *VerifierFeedDeactivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierFeedDeactivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierFeedDeactivatedIterator) Error() error {
	return it.fail
}

func (it *VerifierFeedDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierFeedDeactivated struct {
	FeedId [32]byte
	Raw    types.Log
}

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

func (_Verifier *VerifierFilterer) ParseFeedDeactivated(log types.Log) (*VerifierFeedDeactivated, error) {
	event := new(VerifierFeedDeactivated)
	if err := _Verifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierOwnershipTransferRequestedIterator struct {
	Event *VerifierOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifierOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

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

func (_Verifier *VerifierFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifierOwnershipTransferRequested, error) {
	event := new(VerifierOwnershipTransferRequested)
	if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierOwnershipTransferredIterator struct {
	Event *VerifierOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

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

func (_Verifier *VerifierFilterer) ParseOwnershipTransferred(log types.Log) (*VerifierOwnershipTransferred, error) {
	event := new(VerifierOwnershipTransferred)
	if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierReportVerifiedIterator struct {
	Event *VerifierReportVerified

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierReportVerifiedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VerifierReportVerifiedIterator) Error() error {
	return it.fail
}

func (it *VerifierReportVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierReportVerified struct {
	FeedId    [32]byte
	Requester common.Address
	Raw       types.Log
}

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

func (_Verifier *VerifierFilterer) ParseReportVerified(log types.Log) (*VerifierReportVerified, error) {
	event := new(VerifierReportVerified)
	if err := _Verifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

func (_Verifier *Verifier) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Verifier.abi.Events["ConfigActivated"].ID:
		return _Verifier.ParseConfigActivated(log)
	case _Verifier.abi.Events["ConfigDeactivated"].ID:
		return _Verifier.ParseConfigDeactivated(log)
	case _Verifier.abi.Events["ConfigSet"].ID:
		return _Verifier.ParseConfigSet(log)
	case _Verifier.abi.Events["FeedActivated"].ID:
		return _Verifier.ParseFeedActivated(log)
	case _Verifier.abi.Events["FeedDeactivated"].ID:
		return _Verifier.ParseFeedDeactivated(log)
	case _Verifier.abi.Events["OwnershipTransferRequested"].ID:
		return _Verifier.ParseOwnershipTransferRequested(log)
	case _Verifier.abi.Events["OwnershipTransferred"].ID:
		return _Verifier.ParseOwnershipTransferred(log)
	case _Verifier.abi.Events["ReportVerified"].ID:
		return _Verifier.ParseReportVerified(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifierConfigActivated) Topic() common.Hash {
	return common.HexToHash("0x54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe7746")
}

func (VerifierConfigDeactivated) Topic() common.Hash {
	return common.HexToHash("0x0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c")
}

func (VerifierConfigSet) Topic() common.Hash {
	return common.HexToHash("0xa23a88453230b183877098801ff5a8f771a120e2573eea559ce6c4c2e305a4da")
}

func (VerifierFeedActivated) Topic() common.Hash {
	return common.HexToHash("0xf438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b4")
}

func (VerifierFeedDeactivated) Topic() common.Hash {
	return common.HexToHash("0xfc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd061")
}

func (VerifierOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifierOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifierReportVerified) Topic() common.Hash {
	return common.HexToHash("0x58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc5")
}

func (_Verifier *Verifier) Address() common.Address {
	return _Verifier.address
}

type VerifierInterface interface {
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

	SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetConfigFromSource(opts *bind.TransactOpts, feedId [32]byte, sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error)

	FilterConfigActivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierConfigActivatedIterator, error)

	WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigActivated, feedId [][32]byte) (event.Subscription, error)

	ParseConfigActivated(log types.Log) (*VerifierConfigActivated, error)

	FilterConfigDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierConfigDeactivatedIterator, error)

	WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigDeactivated, feedId [][32]byte) (event.Subscription, error)

	ParseConfigDeactivated(log types.Log) (*VerifierConfigDeactivated, error)

	FilterConfigSet(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VerifierConfigSet, feedId [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VerifierConfigSet, error)

	FilterFeedActivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierFeedActivatedIterator, error)

	WatchFeedActivated(opts *bind.WatchOpts, sink chan<- *VerifierFeedActivated, feedId [][32]byte) (event.Subscription, error)

	ParseFeedActivated(log types.Log) (*VerifierFeedActivated, error)

	FilterFeedDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierFeedDeactivatedIterator, error)

	WatchFeedDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierFeedDeactivated, feedId [][32]byte) (event.Subscription, error)

	ParseFeedDeactivated(log types.Log) (*VerifierFeedDeactivated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifierOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifierOwnershipTransferred, error)

	FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*VerifierReportVerifiedIterator, error)

	WatchReportVerified(opts *bind.WatchOpts, sink chan<- *VerifierReportVerified, feedId [][32]byte) (event.Subscription, error)

	ParseReportVerified(log types.Log) (*VerifierReportVerified, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
