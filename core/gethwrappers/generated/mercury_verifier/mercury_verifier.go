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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001e2538038062001e258339810160408190526200003491620001ab565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000ff565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b60601b6001600160601b031916608052620001dd565b6001600160a01b0381163314156200015a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001be57600080fd5b81516001600160a01b0381168114620001d657600080fd5b9392505050565b60805160601c611c2262000203600039600081816106af01526107db0152611c226000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638da5cb5b11610081578063ded6307c1161005b578063ded6307c14610244578063e84f128e14610257578063f2fde38b146102b457600080fd5b80638da5cb5b146101aa57806394d95980146101d2578063b70d929d146101e557600080fd5b80632374291a116100b25780632374291a1461017a5780633d3ac1b51461018f57806379ba5097146101a257600080fd5b806301ffc9a7146100ce578063181f5a7714610138575b600080fd5b6101236100dc366004611696565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b60408051808201909152600e81527f566572696669657220302e302e3200000000000000000000000000000000000060208201525b60405161012f91906118d9565b61018d61018836600461155c565b6102c7565b005b61016d61019d3660046116d8565b6107c1565b61018d610911565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012f565b61018d6101e0366004611674565b610a0e565b6102216101f3366004611543565b6000908152600260205260408120600181015490549192909168010000000000000000900463ffffffff1690565b604080519315158452602084019290925263ffffffff169082015260600161012f565b61018d610252366004611674565b610b77565b610291610265366004611543565b6000908152600260205260409020805460019091015463ffffffff808316936401000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600161012f565b61018d6102c236600461143c565b610c8e565b845160ff851680610304576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f82111561034e576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b610359816003611a5a565b82116103b1578161036b826003611a5a565b610376906001611a1d565b6040517f9dd9e6d800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610345565b6103b9610ca2565b60008881526002602052604081208054909163ffffffff9091169082906103df83611b15565b82546101009290920a63ffffffff81810219909316918316021790915582546000925061041191168a8a8a8a8a610d25565b6000818152600280850160205260408220805460ff8d167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009182161782559101805490911660011790559091505b89518160ff1610156106605760008a8260ff168151811061048257610482611bb7565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156104f3576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452600190810190925290912054610100900460ff169081111561054657610546611b88565b1480159150610581576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff8416815260208101600190526000858152600287016020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684526001908101835292208351815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00821681178355928501519193919284927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216179061010090849081111561064357610643611b88565b02179055509050505050808061065890611b39565b91505061045f565b508154600163ffffffff90911611156107215760018201546040517f2cc994770000000000000000000000000000000000000000000000000000000081526004810191909152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cc9947790604401600060405180830381600087803b15801561070857600080fd5b505af115801561071c573d6000803e3d6000fd5b505050505b81546040517f6358f01514013bac7818694a3c4cb5f20eae72edee18877977c884280e9370b591610774918d9163ffffffff6401000000008204811692879291909116908f908f908f908f908f90611851565b60405180910390a181547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1664010000000063ffffffff4316021782556001909101555050505050505050565b60603373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610832576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080610844888a018a61145e565b9450945094509450945060008461085a90611a97565b60008181526002602081815260408084208b51808652938101909252909220929350909161088b8483898985610dd1565b6108958984610f10565b875160208901206108aa818b8a8a8a87610f78565b604080518681526020810183905273ffffffffffffffffffffffffffffffffffffffff8e168183015290517f5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a2699181900360600190a150969c9b505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610992576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610345565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a16610ca2565b600082815260026020526040902081610a5b576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604090205460ff16610aaf576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610345565b8060010154821415610af7576040517fa403c0160000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610345565b600082815260028083016020526040918290200180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517f0e173bea63a8c59ec70bf87043f2a729693790183f16a1a54b705de9e989cc4c90610b6a9085908590918252602082015260400190565b60405180910390a1505050565b610b7f610ca2565b600082815260026020526040902081610bc4576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260028201602052604090205460ff16610c18576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810184905260248101839052604401610345565b600082815260028083016020526040918290200180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f54f8872b9b94ebea6577f33576d55847bd8ea22641ccc886b965f6e50bfe774690610b6a9085908590918252602082015260400190565b610c96610ca2565b610c9f816111f4565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d23576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610345565b565b60006060600046308a8a858b8b8b8b604051602001610d4c999897969594939291906118ec565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e0100000000000000000000000000000000000000000000000000000000000017925050509695505050505050565b8054600090610de49060ff166001611a35565b825490915060ff16610e2c576040517f8bca63110000000000000000000000000000000000000000000000000000000081526004810187905260248101869052604401610345565b600282015460ff16610e74576040517ffc10a2830000000000000000000000000000000000000000000000000000000081526004810187905260248101869052604401610345565b8060ff16845114610ec05783516040517f5348a282000000000000000000000000000000000000000000000000000000008152600481019190915260ff82166024820152604401610345565b8251845114610f0857835183516040517ff0d3140800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610345565b505050505050565b6020820151815463ffffffff600883901c81169168010000000000000000900416811115610f725782547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff166801000000000000000063ffffffff8316021783555b50505050565b60008686604051602001610f8d929190611815565b6040516020818303038152906040528051906020012090506000610fc1604080518082019091526000808252602082015290565b8651600090815b8181101561118c57600186898360208110610fe557610fe5611bb7565b610ff291901a601b611a35565b8c848151811061100457611004611bb7565b60200260200101518c858151811061101e5761101e611bb7565b60200260200101516040516000815260200160405260405161105c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561107e573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526001808d01602090815291859020848601909552845460ff80821686529399509395509085019261010090049091169081111561110357611103611b88565b600181111561111457611114611b88565b905250935060018460200151600181111561113157611131611b88565b14611168576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836000015160080260ff166001901b850194508061118590611adc565b9050610fc8565b50837e010101010101010101010101010101010101010101010101010101010101018516146111e7576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415611274576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610345565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803573ffffffffffffffffffffffffffffffffffffffff8116811461130e57600080fd5b919050565b600082601f83011261132457600080fd5b81356020611339611334836119f9565b6119aa565b80838252828201915082860187848660051b890101111561135957600080fd5b60005b858110156113785781358452928401929084019060010161135c565b5090979650505050505050565b600082601f83011261139657600080fd5b813567ffffffffffffffff8111156113b0576113b0611be6565b6113e160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016119aa565b8181528460208386010111156113f657600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461130e57600080fd5b803560ff8116811461130e57600080fd5b60006020828403121561144e57600080fd5b611457826112ea565b9392505050565b600080600080600060e0868803121561147657600080fd5b86601f87011261148557600080fd5b61148d611981565b8087606089018a8111156114a057600080fd5b60005b60038110156114c25782358552602094850194909201916001016114a3565b509197505035905067ffffffffffffffff808211156114e057600080fd5b6114ec89838a01611385565b9550608088013591508082111561150257600080fd5b61150e89838a01611313565b945060a088013591508082111561152457600080fd5b5061153188828901611313565b9598949750929560c001359392505050565b60006020828403121561155557600080fd5b5035919050565b60008060008060008060c0878903121561157557600080fd5b8635955060208088013567ffffffffffffffff8082111561159557600080fd5b818a0191508a601f8301126115a957600080fd5b81356115b7611334826119f9565b8082825285820191508585018e878560051b88010111156115d757600080fd5b600095505b83861015611601576115ed816112ea565b8352600195909501949186019186016115dc565b50995061161391505060408b0161142b565b965060608a013592508083111561162957600080fd5b6116358b848c01611385565b955061164360808b01611413565b945060a08a013592508083111561165957600080fd5b505061166789828a01611385565b9150509295509295509295565b6000806040838503121561168757600080fd5b50508035926020909101359150565b6000602082840312156116a857600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461145757600080fd5b6000806000604084860312156116ed57600080fd5b833567ffffffffffffffff8082111561170557600080fd5b818601915086601f83011261171957600080fd5b81358181111561172857600080fd5b87602082850101111561173a57600080fd5b60209283019550935061175091860190506112ea565b90509250925092565b600081518084526020808501945080840160005b8381101561179f57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161176d565b509495945050505050565b6000815180845260005b818110156117d0576020818501810151868301820152016117b4565b818111156117e2576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b828152600060208083018460005b600381101561184057815183529183019190830190600101611823565b505050506080820190509392505050565b60006101208b835263ffffffff808c1660208501528a6040850152808a1660608501525080608084015261188781840189611759565b905060ff871660a084015282810360c08401526118a481876117aa565b905067ffffffffffffffff851660e08401528281036101008401526118c981856117aa565b9c9b505050505050505050505050565b60208152600061145760208301846117aa565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526119338285018b611759565b91508382036080850152611947828a611759565b915060ff881660a085015283820360c085015261196482886117aa565b90861660e085015283810361010085015290506118c981856117aa565b6040516060810167ffffffffffffffff811182821017156119a4576119a4611be6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156119f1576119f1611be6565b604052919050565b600067ffffffffffffffff821115611a1357611a13611be6565b5060051b60200190565b60008219821115611a3057611a30611b59565b500190565b600060ff821660ff84168060ff03821115611a5257611a52611b59565b019392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611a9257611a92611b59565b500290565b80516020808301519190811015611ad6577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611b0e57611b0e611b59565b5060010190565b600063ffffffff80831681811415611b2f57611b2f611b59565b6001019392505050565b600060ff821660ff811415611b5057611b50611b59565b60010192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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

func (_MercuryVerifier *MercuryVerifierTransactor) SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _MercuryVerifier.contract.Transact(opts, "setConfig", feedId, signers, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_MercuryVerifier *MercuryVerifierSession) SetConfig(feedId [32]byte, signers []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.SetConfig(&_MercuryVerifier.TransactOpts, feedId, signers, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_MercuryVerifier *MercuryVerifierTransactorSession) SetConfig(feedId [32]byte, signers []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _MercuryVerifier.Contract.SetConfig(&_MercuryVerifier.TransactOpts, feedId, signers, f, onchainConfig, offchainConfigVersion, offchainConfig)
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

func (_MercuryVerifier *MercuryVerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts) (*MercuryVerifierConfigActivatedIterator, error) {

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierConfigActivatedIterator{contract: _MercuryVerifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigActivated) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ConfigActivated")
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

func (_MercuryVerifier *MercuryVerifierFilterer) FilterConfigDeactivated(opts *bind.FilterOpts) (*MercuryVerifierConfigDeactivatedIterator, error) {

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ConfigDeactivated")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierConfigDeactivatedIterator{contract: _MercuryVerifier.contract, event: "ConfigDeactivated", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigDeactivated) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ConfigDeactivated")
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
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_MercuryVerifier *MercuryVerifierFilterer) FilterConfigSet(opts *bind.FilterOpts) (*MercuryVerifierConfigSetIterator, error) {

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierConfigSetIterator{contract: _MercuryVerifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigSet) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ConfigSet")
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

func (_MercuryVerifier *MercuryVerifierFilterer) FilterReportVerified(opts *bind.FilterOpts) (*MercuryVerifierReportVerifiedIterator, error) {

	logs, sub, err := _MercuryVerifier.contract.FilterLogs(opts, "ReportVerified")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierReportVerifiedIterator{contract: _MercuryVerifier.contract, event: "ReportVerified", logs: logs, sub: sub}, nil
}

func (_MercuryVerifier *MercuryVerifierFilterer) WatchReportVerified(opts *bind.WatchOpts, sink chan<- *MercuryVerifierReportVerified) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifier.contract.WatchLogs(opts, "ReportVerified")
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
	return common.HexToHash("0x6358f01514013bac7818694a3c4cb5f20eae72edee18877977c884280e9370b5")
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

	SetConfig(opts *bind.TransactOpts, feedId [32]byte, signers []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error)

	FilterConfigActivated(opts *bind.FilterOpts) (*MercuryVerifierConfigActivatedIterator, error)

	WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigActivated) (event.Subscription, error)

	ParseConfigActivated(log types.Log) (*MercuryVerifierConfigActivated, error)

	FilterConfigDeactivated(opts *bind.FilterOpts) (*MercuryVerifierConfigDeactivatedIterator, error)

	WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigDeactivated) (event.Subscription, error)

	ParseConfigDeactivated(log types.Log) (*MercuryVerifierConfigDeactivated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*MercuryVerifierConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*MercuryVerifierConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MercuryVerifierOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MercuryVerifierOwnershipTransferred, error)

	FilterReportVerified(opts *bind.FilterOpts) (*MercuryVerifierReportVerifiedIterator, error)

	WatchReportVerified(opts *bind.WatchOpts, sink chan<- *MercuryVerifierReportVerified) (event.Subscription, error)

	ParseReportVerified(log types.Log) (*MercuryVerifierReportVerified, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
