// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package channel_verifier

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

var ChannelVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadVerification\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InactiveFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"FeedDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"activateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"deactivateFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sourceAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newConfigCount\",\"type\":\"uint32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setConfigFromSource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"verifierResponse\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620020ed380380620020ed8339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611ef9620001f4600039600061051a0152611ef96000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063564a0a7a1161008c5780638da5cb5b116100665780638da5cb5b14610247578063afcb95d71461026f578063eb1dc803146102a4578063f2fde38b146102b757600080fd5b8063564a0a7a146101fc57806379ba50971461020f57806381ff70481461021757600080fd5b8063181f5a77116100c8578063181f5a77146101815780633d3ac1b5146101c35780633dd86430146101d657806354e68a81146101e957600080fd5b806301ffc9a7146100ef5780630d1d79af146101595780630f672ef41461016e575b600080fd5b6101446100fd3660046113d0565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61016c610167366004611419565b6102ca565b005b61016c61017c366004611419565b6103d0565b60408051808201909152601581527f4368616e6e656c566572696669657220302e302e30000000000000000000000060208201525b6040516101509190611496565b6101b66101d13660046114d2565b610500565b61016c6101e4366004611419565b610680565b61016c6101f7366004611866565b6106fd565b61016c61020a366004611419565b61080d565b61016c61088d565b6002546003546040805163ffffffff80851682526401000000009094049093166020840152820152606001610150565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610150565b600354600254604080516000815260208101939093526801000000000000000090910463ffffffff1690820152606001610150565b61016c6102b2366004611986565b61098a565b61016c6102c5366004611a79565b610a52565b6102d2610a63565b80610309576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526005602052604081205460ff16900361035b576040517f74eb4b93000000000000000000000000000000000000000000000000000000008152600481018290526024015b60405180910390fd5b6000818152600560205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055517fa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea906103c59083815260200190565b60405180910390a150565b6103d8610a63565b8061040f576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526005602052604081205460ff16900361045c576040517f74eb4b9300000000000000000000000000000000000000000000000000000000815260048101829052602401610352565b600354810361049a576040517f67863f4400000000000000000000000000000000000000000000000000000000815260048101829052602401610352565b6000818152600560205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055517f5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579906103c59083815260200190565b60603373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610571576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080610583888a018a611a94565b9450945094509450945060008461059990611b6f565b60008181526004602052604090205490915060ff16156105e8576040517f36dbe74800000000000000000000000000000000000000000000000000000000815260048101829052602401610352565b8551600081815260056020526040902061060482878784610ae6565b61060f886002610bda565b86516020880120610624818a89898987610c42565b60405173ffffffffffffffffffffffffffffffffffffffff8c16815284907f58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc59060200160405180910390a250959b9a5050505050505050505050565b610688610a63565b60008181526004602052604090205460ff166106fa5760008181526004602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690555182917ff438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b491a25b50565b86518560ff168060000361073d576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610782576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f6024820152604401610352565b61078d816003611c12565b82116107e5578161079f826003611c12565b6107aa906001611c2f565b6040517f9dd9e6d800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610352565b6107ed610a63565b6107ff8c8c8c8c8c8c8c8c8c8c610ebe565b505050505050505050505050565b610815610a63565b60008181526004602052604090205460ff16156106fa5760008181526004602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555182917ffc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd06191a250565b60015473ffffffffffffffffffffffffffffffffffffffff16331461090e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610352565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b86518560ff16806000036109ca576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610a0f576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f6024820152604401610352565b610a1a816003611c12565b8211610a2c578161079f826003611c12565b610a34610a63565b610a47463060008c8c8c8c8c8c8c610ebe565b505050505050505050565b610a5a610a63565b6106fa81611230565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ae4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610352565b565b8054600090610af99060ff166001611c42565b8254909150610100900460ff16610b3f576040517fd990d62100000000000000000000000000000000000000000000000000000000815260048101869052602401610352565b8060ff16845114610b8b5783516040517f5348a282000000000000000000000000000000000000000000000000000000008152600481019190915260ff82166024820152604401610352565b8251845114610bd357835183516040517ff0d3140800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610352565b5050505050565b6020820151815463ffffffff600883901c81169168010000000000000000900416811115610c3c5782547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff166801000000000000000063ffffffff8316021783555b50505050565b60008686604051602001610c57929190611c5b565b6040516020818303038152906040528051906020012090506000610c8b604080518082019091526000808252602082015290565b8651600090815b81811015610e5657600186898360208110610caf57610caf611bb4565b610cbc91901a601b611c42565b8c8481518110610cce57610cce611bb4565b60200260200101518c8581518110610ce857610ce8611bb4565b602002602001015160405160008152602001604052604051610d26949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610d48573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526001808d01602090815291859020848601909552845460ff808216865293995093955090850192610100900490911690811115610dcd57610dcd611c97565b6001811115610dde57610dde611c97565b9052509350600184602001516001811115610dfb57610dfb611c97565b14610e32576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836000015160080260ff166001901b8501945080610e4f90611cc6565b9050610c92565b50837e01010101010101010101010101010101010101010101010101010101010101851614610eb1576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b63ffffffff881615610eff57600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8a16179055610f35565b6002805463ffffffff16906000610f1583611cfe565b91906101000a81548163ffffffff021916908363ffffffff160217905550505b600254600090610f54908c908c9063ffffffff168b8b8b8b8b8b611325565b600081815260056020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660ff8a16176101001790559091505b88518160ff16101561118f576000898260ff1681518110610fb857610fb8611bb4565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611028576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600085815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452600190810190925290912054610100900460ff169081111561107957611079611c97565b14801591506110b4576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff841681526020810160019052600085815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684526001908101835292208351815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00821681178355928501519193919284927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216179061010090849081111561117457611174611c97565b021790555090505050508061118890611d21565b9050610f95565b506002546040517f1074b4b9a073f79bd1f7f5c808348125ce0f25c27188df7efcaa7a08276051b3916111e29163ffffffff640100000000830481169286929116908d908d908d908d908d908d90611dc1565b60405180910390a1600280547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1664010000000063ffffffff43160217905560035550505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036112af576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610352565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808a8a8a8a8a8a8a8a8a60405160200161134999989796959493929190611e57565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b6000602082840312156113e257600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461141257600080fd5b9392505050565b60006020828403121561142b57600080fd5b5035919050565b6000815180845260005b818110156114585760208185018101518683018201520161143c565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006114126020830184611432565b803573ffffffffffffffffffffffffffffffffffffffff811681146114cd57600080fd5b919050565b6000806000604084860312156114e757600080fd5b833567ffffffffffffffff808211156114ff57600080fd5b818601915086601f83011261151357600080fd5b81358181111561152257600080fd5b87602082850101111561153457600080fd5b60209283019550935061154a91860190506114a9565b90509250925092565b803563ffffffff811681146114cd57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156115b9576115b9611567565b60405290565b6040516060810167ffffffffffffffff811182821017156115b9576115b9611567565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561162957611629611567565b604052919050565b600067ffffffffffffffff82111561164b5761164b611567565b5060051b60200190565b600082601f83011261166657600080fd5b8135602061167b61167683611631565b6115e2565b82815260059290921b8401810191818101908684111561169a57600080fd5b8286015b848110156116bc576116af816114a9565b835291830191830161169e565b509695505050505050565b600082601f8301126116d857600080fd5b813560206116e861167683611631565b82815260059290921b8401810191818101908684111561170757600080fd5b8286015b848110156116bc578035835291830191830161170b565b803560ff811681146114cd57600080fd5b600082601f83011261174457600080fd5b813567ffffffffffffffff81111561175e5761175e611567565b61178f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016115e2565b8181528460208386010111156117a457600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff811681146114cd57600080fd5b600082601f8301126117ea57600080fd5b813560206117fa61167683611631565b82815260069290921b8401810191818101908684111561181957600080fd5b8286015b848110156116bc57604081890312156118365760008081fd5b61183e611596565b611847826114a9565b81526118548583016117c1565b8186015283529183019160400161181d565b6000806000806000806000806000806101408b8d03121561188657600080fd5b8a35995061189660208c016114a9565b98506118a460408c01611553565b975060608b013567ffffffffffffffff808211156118c157600080fd5b6118cd8e838f01611655565b985060808d01359150808211156118e357600080fd5b6118ef8e838f016116c7565b97506118fd60a08e01611722565b965060c08d013591508082111561191357600080fd5b61191f8e838f01611733565b955061192d60e08e016117c1565b94506101008d013591508082111561194457600080fd5b6119508e838f01611733565b93506101208d013591508082111561196757600080fd5b506119748d828e016117d9565b9150509295989b9194979a5092959850565b600080600080600080600060e0888a0312156119a157600080fd5b873567ffffffffffffffff808211156119b957600080fd5b6119c58b838c01611655565b985060208a01359150808211156119db57600080fd5b6119e78b838c016116c7565b97506119f560408b01611722565b965060608a0135915080821115611a0b57600080fd5b611a178b838c01611733565b9550611a2560808b016117c1565b945060a08a0135915080821115611a3b57600080fd5b611a478b838c01611733565b935060c08a0135915080821115611a5d57600080fd5b50611a6a8a828b016117d9565b91505092959891949750929550565b600060208284031215611a8b57600080fd5b611412826114a9565b600080600080600060e08688031215611aac57600080fd5b86601f870112611abb57600080fd5b611ac36115bf565b806060880189811115611ad557600080fd5b885b81811015611aef578035845260209384019301611ad7565b5090965035905067ffffffffffffffff80821115611b0c57600080fd5b611b1889838a01611733565b95506080880135915080821115611b2e57600080fd5b611b3a89838a016116c7565b945060a0880135915080821115611b5057600080fd5b50611b5d888289016116c7565b9598949750929560c001359392505050565b80516020808301519190811015611bae577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417611c2957611c29611be3565b92915050565b80820180821115611c2957611c29611be3565b60ff8181168382160190811115611c2957611c29611be3565b828152600060208083018460005b6003811015611c8657815183529183019190830190600101611c69565b505050506080820190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611cf757611cf7611be3565b5060010190565b600063ffffffff808316818103611d1757611d17611be3565b6001019392505050565b600060ff821660ff8103611d3757611d37611be3565b60010192915050565b600081518084526020808501945080840160005b83811015611d8657815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611d54565b509495945050505050565b600081518084526020808501945080840160005b83811015611d8657815187529582019590820190600101611da5565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611df18184018a611d40565b90508281036080840152611e058189611d91565b905060ff871660a084015282810360c0840152611e228187611432565b905067ffffffffffffffff851660e0840152828103610100840152611e478185611432565b9c9b505050505050505050505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152611e9e8285018b611d40565b91508382036080850152611eb2828a611d91565b915060ff881660a085015283820360c0850152611ecf8288611432565b90861660e08501528381036101008501529050611e47818561143256fea164736f6c6343000813000a",
}

var ChannelVerifierABI = ChannelVerifierMetaData.ABI

var ChannelVerifierBin = ChannelVerifierMetaData.Bin

func DeployChannelVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, verifierProxyAddr common.Address) (common.Address, *types.Transaction, *ChannelVerifier, error) {
	parsed, err := ChannelVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChannelVerifierBin), backend, verifierProxyAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChannelVerifier{address: address, abi: *parsed, ChannelVerifierCaller: ChannelVerifierCaller{contract: contract}, ChannelVerifierTransactor: ChannelVerifierTransactor{contract: contract}, ChannelVerifierFilterer: ChannelVerifierFilterer{contract: contract}}, nil
}

type ChannelVerifier struct {
	address common.Address
	abi     abi.ABI
	ChannelVerifierCaller
	ChannelVerifierTransactor
	ChannelVerifierFilterer
}

type ChannelVerifierCaller struct {
	contract *bind.BoundContract
}

type ChannelVerifierTransactor struct {
	contract *bind.BoundContract
}

type ChannelVerifierFilterer struct {
	contract *bind.BoundContract
}

type ChannelVerifierSession struct {
	Contract     *ChannelVerifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChannelVerifierCallerSession struct {
	Contract *ChannelVerifierCaller
	CallOpts bind.CallOpts
}

type ChannelVerifierTransactorSession struct {
	Contract     *ChannelVerifierTransactor
	TransactOpts bind.TransactOpts
}

type ChannelVerifierRaw struct {
	Contract *ChannelVerifier
}

type ChannelVerifierCallerRaw struct {
	Contract *ChannelVerifierCaller
}

type ChannelVerifierTransactorRaw struct {
	Contract *ChannelVerifierTransactor
}

func NewChannelVerifier(address common.Address, backend bind.ContractBackend) (*ChannelVerifier, error) {
	abi, err := abi.JSON(strings.NewReader(ChannelVerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChannelVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifier{address: address, abi: abi, ChannelVerifierCaller: ChannelVerifierCaller{contract: contract}, ChannelVerifierTransactor: ChannelVerifierTransactor{contract: contract}, ChannelVerifierFilterer: ChannelVerifierFilterer{contract: contract}}, nil
}

func NewChannelVerifierCaller(address common.Address, caller bind.ContractCaller) (*ChannelVerifierCaller, error) {
	contract, err := bindChannelVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierCaller{contract: contract}, nil
}

func NewChannelVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelVerifierTransactor, error) {
	contract, err := bindChannelVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierTransactor{contract: contract}, nil
}

func NewChannelVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelVerifierFilterer, error) {
	contract, err := bindChannelVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierFilterer{contract: contract}, nil
}

func bindChannelVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChannelVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChannelVerifier *ChannelVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelVerifier.Contract.ChannelVerifierCaller.contract.Call(opts, result, method, params...)
}

func (_ChannelVerifier *ChannelVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.ChannelVerifierTransactor.contract.Transfer(opts)
}

func (_ChannelVerifier *ChannelVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.ChannelVerifierTransactor.contract.Transact(opts, method, params...)
}

func (_ChannelVerifier *ChannelVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelVerifier.Contract.contract.Call(opts, result, method, params...)
}

func (_ChannelVerifier *ChannelVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.contract.Transfer(opts)
}

func (_ChannelVerifier *ChannelVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.contract.Transact(opts, method, params...)
}

func (_ChannelVerifier *ChannelVerifierCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _ChannelVerifier.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_ChannelVerifier *ChannelVerifierSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _ChannelVerifier.Contract.LatestConfigDetails(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _ChannelVerifier.Contract.LatestConfigDetails(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _ChannelVerifier.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_ChannelVerifier *ChannelVerifierSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _ChannelVerifier.Contract.LatestConfigDigestAndEpoch(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _ChannelVerifier.Contract.LatestConfigDigestAndEpoch(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChannelVerifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ChannelVerifier *ChannelVerifierSession) Owner() (common.Address, error) {
	return _ChannelVerifier.Contract.Owner(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCallerSession) Owner() (common.Address, error) {
	return _ChannelVerifier.Contract.Owner(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ChannelVerifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ChannelVerifier *ChannelVerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelVerifier.Contract.SupportsInterface(&_ChannelVerifier.CallOpts, interfaceId)
}

func (_ChannelVerifier *ChannelVerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelVerifier.Contract.SupportsInterface(&_ChannelVerifier.CallOpts, interfaceId)
}

func (_ChannelVerifier *ChannelVerifierCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ChannelVerifier.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ChannelVerifier *ChannelVerifierSession) TypeAndVersion() (string, error) {
	return _ChannelVerifier.Contract.TypeAndVersion(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierCallerSession) TypeAndVersion() (string, error) {
	return _ChannelVerifier.Contract.TypeAndVersion(&_ChannelVerifier.CallOpts)
}

func (_ChannelVerifier *ChannelVerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "acceptOwnership")
}

func (_ChannelVerifier *ChannelVerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelVerifier.Contract.AcceptOwnership(&_ChannelVerifier.TransactOpts)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelVerifier.Contract.AcceptOwnership(&_ChannelVerifier.TransactOpts)
}

func (_ChannelVerifier *ChannelVerifierTransactor) ActivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "activateConfig", configDigest)
}

func (_ChannelVerifier *ChannelVerifierSession) ActivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.ActivateConfig(&_ChannelVerifier.TransactOpts, configDigest)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) ActivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.ActivateConfig(&_ChannelVerifier.TransactOpts, configDigest)
}

func (_ChannelVerifier *ChannelVerifierTransactor) ActivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "activateFeed", feedId)
}

func (_ChannelVerifier *ChannelVerifierSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.ActivateFeed(&_ChannelVerifier.TransactOpts, feedId)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) ActivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.ActivateFeed(&_ChannelVerifier.TransactOpts, feedId)
}

func (_ChannelVerifier *ChannelVerifierTransactor) DeactivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "deactivateConfig", configDigest)
}

func (_ChannelVerifier *ChannelVerifierSession) DeactivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.DeactivateConfig(&_ChannelVerifier.TransactOpts, configDigest)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) DeactivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.DeactivateConfig(&_ChannelVerifier.TransactOpts, configDigest)
}

func (_ChannelVerifier *ChannelVerifierTransactor) DeactivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "deactivateFeed", feedId)
}

func (_ChannelVerifier *ChannelVerifierSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.DeactivateFeed(&_ChannelVerifier.TransactOpts, feedId)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) DeactivateFeed(feedId [32]byte) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.DeactivateFeed(&_ChannelVerifier.TransactOpts, feedId)
}

func (_ChannelVerifier *ChannelVerifierTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "setConfig", signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_ChannelVerifier *ChannelVerifierSession) SetConfig(signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.SetConfig(&_ChannelVerifier.TransactOpts, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) SetConfig(signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.SetConfig(&_ChannelVerifier.TransactOpts, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_ChannelVerifier *ChannelVerifierTransactor) SetConfigFromSource(opts *bind.TransactOpts, sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "setConfigFromSource", sourceChainId, sourceAddress, newConfigCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_ChannelVerifier *ChannelVerifierSession) SetConfigFromSource(sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.SetConfigFromSource(&_ChannelVerifier.TransactOpts, sourceChainId, sourceAddress, newConfigCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) SetConfigFromSource(sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.SetConfigFromSource(&_ChannelVerifier.TransactOpts, sourceChainId, sourceAddress, newConfigCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
}

func (_ChannelVerifier *ChannelVerifierTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "transferOwnership", to)
}

func (_ChannelVerifier *ChannelVerifierSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.TransferOwnership(&_ChannelVerifier.TransactOpts, to)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.TransferOwnership(&_ChannelVerifier.TransactOpts, to)
}

func (_ChannelVerifier *ChannelVerifierTransactor) Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _ChannelVerifier.contract.Transact(opts, "verify", signedReport, sender)
}

func (_ChannelVerifier *ChannelVerifierSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.Verify(&_ChannelVerifier.TransactOpts, signedReport, sender)
}

func (_ChannelVerifier *ChannelVerifierTransactorSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _ChannelVerifier.Contract.Verify(&_ChannelVerifier.TransactOpts, signedReport, sender)
}

type ChannelVerifierConfigActivatedIterator struct {
	Event *ChannelVerifierConfigActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierConfigActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierConfigActivated)
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
		it.Event = new(ChannelVerifierConfigActivated)
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

func (it *ChannelVerifierConfigActivatedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierConfigActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierConfigActivated struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts) (*ChannelVerifierConfigActivatedIterator, error) {

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierConfigActivatedIterator{contract: _ChannelVerifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierConfigActivated) (event.Subscription, error) {

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierConfigActivated)
				if err := _ChannelVerifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseConfigActivated(log types.Log) (*ChannelVerifierConfigActivated, error) {
	event := new(ChannelVerifierConfigActivated)
	if err := _ChannelVerifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierConfigDeactivatedIterator struct {
	Event *ChannelVerifierConfigDeactivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierConfigDeactivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierConfigDeactivated)
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
		it.Event = new(ChannelVerifierConfigDeactivated)
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

func (it *ChannelVerifierConfigDeactivatedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierConfigDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierConfigDeactivated struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterConfigDeactivated(opts *bind.FilterOpts) (*ChannelVerifierConfigDeactivatedIterator, error) {

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "ConfigDeactivated")
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierConfigDeactivatedIterator{contract: _ChannelVerifier.contract, event: "ConfigDeactivated", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierConfigDeactivated) (event.Subscription, error) {

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "ConfigDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierConfigDeactivated)
				if err := _ChannelVerifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseConfigDeactivated(log types.Log) (*ChannelVerifierConfigDeactivated, error) {
	event := new(ChannelVerifierConfigDeactivated)
	if err := _ChannelVerifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierConfigSetIterator struct {
	Event *ChannelVerifierConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierConfigSet)
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
		it.Event = new(ChannelVerifierConfigSet)
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

func (it *ChannelVerifierConfigSetIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierConfigSet struct {
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

func (_ChannelVerifier *ChannelVerifierFilterer) FilterConfigSet(opts *bind.FilterOpts) (*ChannelVerifierConfigSetIterator, error) {

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierConfigSetIterator{contract: _ChannelVerifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *ChannelVerifierConfigSet) (event.Subscription, error) {

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierConfigSet)
				if err := _ChannelVerifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseConfigSet(log types.Log) (*ChannelVerifierConfigSet, error) {
	event := new(ChannelVerifierConfigSet)
	if err := _ChannelVerifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierFeedActivatedIterator struct {
	Event *ChannelVerifierFeedActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierFeedActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierFeedActivated)
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
		it.Event = new(ChannelVerifierFeedActivated)
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

func (it *ChannelVerifierFeedActivatedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierFeedActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierFeedActivated struct {
	FeedId [32]byte
	Raw    types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterFeedActivated(opts *bind.FilterOpts, feedId [][32]byte) (*ChannelVerifierFeedActivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "FeedActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierFeedActivatedIterator{contract: _ChannelVerifier.contract, event: "FeedActivated", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchFeedActivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierFeedActivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "FeedActivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierFeedActivated)
				if err := _ChannelVerifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseFeedActivated(log types.Log) (*ChannelVerifierFeedActivated, error) {
	event := new(ChannelVerifierFeedActivated)
	if err := _ChannelVerifier.contract.UnpackLog(event, "FeedActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierFeedDeactivatedIterator struct {
	Event *ChannelVerifierFeedDeactivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierFeedDeactivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierFeedDeactivated)
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
		it.Event = new(ChannelVerifierFeedDeactivated)
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

func (it *ChannelVerifierFeedDeactivatedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierFeedDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierFeedDeactivated struct {
	FeedId [32]byte
	Raw    types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterFeedDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*ChannelVerifierFeedDeactivatedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "FeedDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierFeedDeactivatedIterator{contract: _ChannelVerifier.contract, event: "FeedDeactivated", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchFeedDeactivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierFeedDeactivated, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "FeedDeactivated", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierFeedDeactivated)
				if err := _ChannelVerifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseFeedDeactivated(log types.Log) (*ChannelVerifierFeedDeactivated, error) {
	event := new(ChannelVerifierFeedDeactivated)
	if err := _ChannelVerifier.contract.UnpackLog(event, "FeedDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierOwnershipTransferRequestedIterator struct {
	Event *ChannelVerifierOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierOwnershipTransferRequested)
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
		it.Event = new(ChannelVerifierOwnershipTransferRequested)
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

func (it *ChannelVerifierOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierOwnershipTransferRequestedIterator{contract: _ChannelVerifier.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierOwnershipTransferRequested)
				if err := _ChannelVerifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseOwnershipTransferRequested(log types.Log) (*ChannelVerifierOwnershipTransferRequested, error) {
	event := new(ChannelVerifierOwnershipTransferRequested)
	if err := _ChannelVerifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierOwnershipTransferredIterator struct {
	Event *ChannelVerifierOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierOwnershipTransferred)
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
		it.Event = new(ChannelVerifierOwnershipTransferred)
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

func (it *ChannelVerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierOwnershipTransferredIterator{contract: _ChannelVerifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierOwnershipTransferred)
				if err := _ChannelVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseOwnershipTransferred(log types.Log) (*ChannelVerifierOwnershipTransferred, error) {
	event := new(ChannelVerifierOwnershipTransferred)
	if err := _ChannelVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierReportVerifiedIterator struct {
	Event *ChannelVerifierReportVerified

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierReportVerifiedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierReportVerified)
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
		it.Event = new(ChannelVerifierReportVerified)
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

func (it *ChannelVerifierReportVerifiedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierReportVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierReportVerified struct {
	FeedId    [32]byte
	Requester common.Address
	Raw       types.Log
}

func (_ChannelVerifier *ChannelVerifierFilterer) FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*ChannelVerifierReportVerifiedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _ChannelVerifier.contract.FilterLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierReportVerifiedIterator{contract: _ChannelVerifier.contract, event: "ReportVerified", logs: logs, sub: sub}, nil
}

func (_ChannelVerifier *ChannelVerifierFilterer) WatchReportVerified(opts *bind.WatchOpts, sink chan<- *ChannelVerifierReportVerified, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _ChannelVerifier.contract.WatchLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierReportVerified)
				if err := _ChannelVerifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifierFilterer) ParseReportVerified(log types.Log) (*ChannelVerifierReportVerified, error) {
	event := new(ChannelVerifierReportVerified)
	if err := _ChannelVerifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

func (_ChannelVerifier *ChannelVerifier) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ChannelVerifier.abi.Events["ConfigActivated"].ID:
		return _ChannelVerifier.ParseConfigActivated(log)
	case _ChannelVerifier.abi.Events["ConfigDeactivated"].ID:
		return _ChannelVerifier.ParseConfigDeactivated(log)
	case _ChannelVerifier.abi.Events["ConfigSet"].ID:
		return _ChannelVerifier.ParseConfigSet(log)
	case _ChannelVerifier.abi.Events["FeedActivated"].ID:
		return _ChannelVerifier.ParseFeedActivated(log)
	case _ChannelVerifier.abi.Events["FeedDeactivated"].ID:
		return _ChannelVerifier.ParseFeedDeactivated(log)
	case _ChannelVerifier.abi.Events["OwnershipTransferRequested"].ID:
		return _ChannelVerifier.ParseOwnershipTransferRequested(log)
	case _ChannelVerifier.abi.Events["OwnershipTransferred"].ID:
		return _ChannelVerifier.ParseOwnershipTransferred(log)
	case _ChannelVerifier.abi.Events["ReportVerified"].ID:
		return _ChannelVerifier.ParseReportVerified(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ChannelVerifierConfigActivated) Topic() common.Hash {
	return common.HexToHash("0xa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea")
}

func (ChannelVerifierConfigDeactivated) Topic() common.Hash {
	return common.HexToHash("0x5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579")
}

func (ChannelVerifierConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1074b4b9a073f79bd1f7f5c808348125ce0f25c27188df7efcaa7a08276051b3")
}

func (ChannelVerifierFeedActivated) Topic() common.Hash {
	return common.HexToHash("0xf438564f793525caa89c6e3a26d41e16aa39d1e589747595751e3f3df75cb2b4")
}

func (ChannelVerifierFeedDeactivated) Topic() common.Hash {
	return common.HexToHash("0xfc4f79b8c65b6be1773063461984c0974400d1e99654c79477a092ace83fd061")
}

func (ChannelVerifierOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ChannelVerifierOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (ChannelVerifierReportVerified) Topic() common.Hash {
	return common.HexToHash("0x58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc5")
}

func (_ChannelVerifier *ChannelVerifier) Address() common.Address {
	return _ChannelVerifier.address
}

type ChannelVerifierInterface interface {
	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ActivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	ActivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error)

	DeactivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	DeactivateFeed(opts *bind.TransactOpts, feedId [32]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetConfigFromSource(opts *bind.TransactOpts, sourceChainId *big.Int, sourceAddress common.Address, newConfigCount uint32, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error)

	FilterConfigActivated(opts *bind.FilterOpts) (*ChannelVerifierConfigActivatedIterator, error)

	WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierConfigActivated) (event.Subscription, error)

	ParseConfigActivated(log types.Log) (*ChannelVerifierConfigActivated, error)

	FilterConfigDeactivated(opts *bind.FilterOpts) (*ChannelVerifierConfigDeactivatedIterator, error)

	WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierConfigDeactivated) (event.Subscription, error)

	ParseConfigDeactivated(log types.Log) (*ChannelVerifierConfigDeactivated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*ChannelVerifierConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *ChannelVerifierConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*ChannelVerifierConfigSet, error)

	FilterFeedActivated(opts *bind.FilterOpts, feedId [][32]byte) (*ChannelVerifierFeedActivatedIterator, error)

	WatchFeedActivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierFeedActivated, feedId [][32]byte) (event.Subscription, error)

	ParseFeedActivated(log types.Log) (*ChannelVerifierFeedActivated, error)

	FilterFeedDeactivated(opts *bind.FilterOpts, feedId [][32]byte) (*ChannelVerifierFeedDeactivatedIterator, error)

	WatchFeedDeactivated(opts *bind.WatchOpts, sink chan<- *ChannelVerifierFeedDeactivated, feedId [][32]byte) (event.Subscription, error)

	ParseFeedDeactivated(log types.Log) (*ChannelVerifierFeedDeactivated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ChannelVerifierOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ChannelVerifierOwnershipTransferred, error)

	FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*ChannelVerifierReportVerifiedIterator, error)

	WatchReportVerified(opts *bind.WatchOpts, sink chan<- *ChannelVerifierReportVerified, feedId [][32]byte) (event.Subscription, error)

	ParseReportVerified(log types.Log) (*ChannelVerifierReportVerified, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
