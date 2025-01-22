// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier_v0_5_0

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadVerification\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"prevSigners\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"newSigners\",\"type\":\"address[]\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"prevSigners\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"newSigners\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"verifierResponse\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001e7038038062001e708339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611c75620001fb600039600081816107c20152610e110152611c756000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80633d3ac1b5116100815780638da5cb5b1161005b5780638da5cb5b146101e3578063e84f128e1461020b578063f2fde38b1461024657600080fd5b80633d3ac1b5146101b557806341e3df58146101c857806379ba5097146101db57600080fd5b80630e112e54116100b25780630e112e541461014d5780630f672ef414610160578063181f5a771461017357600080fd5b806301ffc9a7146100ce5780630d1d79af14610138575b600080fd5b6101236100dc366004611333565b7fffffffff00000000000000000000000000000000000000000000000000000000167f3d3ac1b5000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61014b61014636600461137c565b610259565b005b61014b61015b3660046113f7565b610351565b61014b61016e36600461137c565b6106bc565b60408051808201909152600e81527f566572696669657220322e302e3000000000000000000000000000000000000060208201525b60405161012f91906114e4565b6101a86101c336600461151b565b6107a8565b61014b6101d636600461168a565b6108cc565b61014b61098e565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012f565b61023161021936600461137c565b60009081526002602052604090205463ffffffff1690565b60405163ffffffff909116815260200161012f565b61014b6102543660046117a3565b610a8b565b610261610a9f565b6000818152600260205260409020816102a6576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805465010000000000900460ff166000036102f5576040517f74eb4b93000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff1664010000000017815560405182907fa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea90600090a25050565b8160ff82166000819003610391576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f8211156103d6576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044016102ec565b6103e18160036117ed565b821161043957816103f38260036117ed565b6103fe90600161180a565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260248201526044016102ec565b610441610a9f565b6000888152600260205260408120805490916501000000000090910460ff16900361049b576040517f74eb4b93000000000000000000000000000000000000000000000000000000008152600481018a90526024016102ec565b80546601000000000000900460ff1687146104e2576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8781101561061c5760008260010160008b8b858181106105075761050761184c565b905060200201602081019061051c91906117a3565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002054610100900460ff16600181111561055c5761055c61181d565b03610593576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8160010160008a8a848181106105ab576105ab61184c565b90506020020160208101906105c091906117a3565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690556106158161187b565b90506104e5565b5061067389878787600060405190808252806020026020018201604052801561066b57816020015b60408051808201909152600080825260208201528152602001906001900390816106445790505b506001610b22565b887fb0b75a854fab801413da6202fc07e875c54eaf371a1e3909fb2645364ba58616898989896040516106a99493929190611907565b60405180910390a2505050505050505050565b6106c4610a9f565b600081815260026020526040902081610709576040517fe332262700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805465010000000000900460ff16600003610753576040517f74eb4b93000000000000000000000000000000000000000000000000000000008152600481018390526024016102ec565b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff16815560405182907f5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d607657990600090a25050565b60603373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610819576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008080808061082b888a018a611a2d565b84516000818152600260205260409020959a50939850919650945092509061085582868684610ec2565b8551602087012061086a818988888887610fc2565b61087387611b08565b60405173ffffffffffffffffffffffffffffffffffffffff8c1681527f58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc59060200160405180910390a250949a9950505050505050505050565b8260ff8316600081900361090c576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610951576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044016102ec565b61095c8160036117ed565b821161096e57816103f38260036117ed565b610976610a9f565b61098587878787876000610b22565b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a0f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102ec565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a93610a9f565b610a9c8161123e565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102ec565b565b6000868152600260205260409020805465010000000000900460ff1615801590610b4a575081155b15610b81576040517f961dba8800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805460ff8681166601000000000000027fffffffffffffffffffffffffffffffffffffffffffffffffff00ff00ffffffff91871665010000000000027fffffffffffffffffffffffffffffffffffffffffffffffffffff00ff0000000090931663ffffffff43161792909217161764010000000017815560005b60ff8116861115610dce57600087878360ff16818110610c1d57610c1d61184c565b9050602002016020810190610c3291906117a3565b905073ffffffffffffffffffffffffffffffffffffffff8116610c81576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008073ffffffffffffffffffffffffffffffffffffffff831660009081526001868101602052604090912054610100900460ff1690811115610cc657610cc661181d565b1480159150610d01576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201825260ff85811682526001602080840182815273ffffffffffffffffffffffffffffffffffffffff881660009081528a84019092529490208351815493167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008416811782559451939490939284927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090911690911790610100908490811115610db357610db361181d565b0217905550905050505080610dc790611b4d565b9050610bfb565b5081610985576040517fb011b24700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063b011b24790610e4b906000908b908890600401611b6c565b600060405180830381600087803b158015610e6557600080fd5b505af1158015610e79573d6000803e3d6000fd5b50505050867f5b1f376eb2bda670fa39339616d0a73f45b61bec8faeba8ca834f2ebb49676e0878787604051610eb193929190611bec565b60405180910390a250505050505050565b8054600090610ede9065010000000000900460ff166001611c13565b8254909150640100000000900460ff16610f27576040517fd990d621000000000000000000000000000000000000000000000000000000008152600481018690526024016102ec565b8060ff16845114610f735783516040517f5348a282000000000000000000000000000000000000000000000000000000008152600481019190915260ff821660248201526044016102ec565b8251845114610fbb57835183516040517ff0d31408000000000000000000000000000000000000000000000000000000008152600481019290925260248201526044016102ec565b5050505050565b60008686604051602001610fd7929190611c2c565b604051602081830303815290604052805190602001209050600061100b604080518082019091526000808252602082015290565b8651600090815b818110156111d65760018689836020811061102f5761102f61184c565b61103c91901a601b611c13565b8c848151811061104e5761104e61184c565b60200260200101518c85815181106110685761106861184c565b6020026020010151604051600081526020016040526040516110a6949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156110c8573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526001808d01602090815291859020848601909552845460ff80821686529399509395509085019261010090049091169081111561114d5761114d61181d565b600181111561115e5761115e61181d565b905250935060018460200151600181111561117b5761117b61181d565b146111b2576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836000015160080260ff166001901b85019450806111cf9061187b565b9050611012565b50837e01010101010101010101010101010101010101010101010101010101010101851614611231576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036112bd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102ec565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561134557600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461137557600080fd5b9392505050565b60006020828403121561138e57600080fd5b5035919050565b60008083601f8401126113a757600080fd5b50813567ffffffffffffffff8111156113bf57600080fd5b6020830191508360208260051b85010111156113da57600080fd5b9250929050565b803560ff811681146113f257600080fd5b919050565b6000806000806000806080878903121561141057600080fd5b86359550602087013567ffffffffffffffff8082111561142f57600080fd5b61143b8a838b01611395565b9097509550604089013591508082111561145457600080fd5b5061146189828a01611395565b90945092506114749050606088016113e1565b90509295509295509295565b6000815180845260005b818110156114a65760208185018101518683018201520161148a565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006113756020830184611480565b803573ffffffffffffffffffffffffffffffffffffffff811681146113f257600080fd5b60008060006040848603121561153057600080fd5b833567ffffffffffffffff8082111561154857600080fd5b818601915086601f83011261155c57600080fd5b81358181111561156b57600080fd5b87602082850101111561157d57600080fd5b60209283019550935061159391860190506114f7565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156115ee576115ee61159c565b60405290565b6040516060810167ffffffffffffffff811182821017156115ee576115ee61159c565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561165e5761165e61159c565b604052919050565b600067ffffffffffffffff8211156116805761168061159c565b5060051b60200190565b6000806000806000608086880312156116a257600080fd5b8535945060208087013567ffffffffffffffff808211156116c257600080fd5b6116ce8a838b01611395565b9097509550604091506116e28983016113e1565b94506060890135818111156116f657600080fd5b8901601f81018b1361170757600080fd5b803561171a61171582611666565b611617565b81815260069190911b8201850190858101908d83111561173957600080fd5b928601925b8284101561178f5785848f0312156117565760008081fd5b61175e6115cb565b611767856114f7565b815287850135868116811461177c5760008081fd5b818901528252928501929086019061173e565b809750505050505050509295509295909350565b6000602082840312156117b557600080fd5b611375826114f7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417611804576118046117be565b92915050565b80820180821115611804576118046117be565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036118ac576118ac6117be565b5060010190565b8183526000602080850194508260005b858110156118fc5773ffffffffffffffffffffffffffffffffffffffff6118e9836114f7565b16875295820195908201906001016118c3565b509495945050505050565b60408152600061191b6040830186886118b3565b828103602084015261192e8185876118b3565b979650505050505050565b600082601f83011261194a57600080fd5b813567ffffffffffffffff8111156119645761196461159c565b61199560207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611617565b8181528460208386010111156119aa57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126119d857600080fd5b813560206119e861171583611666565b82815260059290921b84018101918181019086841115611a0757600080fd5b8286015b84811015611a225780358352918301918301611a0b565b509695505050505050565b600080600080600060e08688031215611a4557600080fd5b86601f870112611a5457600080fd5b611a5c6115f4565b806060880189811115611a6e57600080fd5b885b81811015611a88578035845260209384019301611a70565b5090965035905067ffffffffffffffff80821115611aa557600080fd5b611ab189838a01611939565b95506080880135915080821115611ac757600080fd5b611ad389838a016119c7565b945060a0880135915080821115611ae957600080fd5b50611af6888289016119c7565b9598949750929560c001359392505050565b80516020808301519190811015611b47577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b600060ff821660ff8103611b6357611b636117be565b60010192915050565b600060608201858352602085818501526040606081860152828651808552608087019150838801945060005b81811015611bdd578551805173ffffffffffffffffffffffffffffffffffffffff16845285015167ffffffffffffffff16858401529484019491830191600101611b98565b50909998505050505050505050565b604081526000611c006040830185876118b3565b905060ff83166020830152949350505050565b60ff8181168382160190811115611804576118046117be565b828152600060208083018460005b6003811015611c5757815183529183019190830190600101611c3a565b50505050608082019050939250505056fea164736f6c6343000813000a",
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

func (_Verifier *VerifierCaller) LatestConfigDetails(opts *bind.CallOpts, configDigest [32]byte) (uint32, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDetails", configDigest)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_Verifier *VerifierSession) LatestConfigDetails(configDigest [32]byte) (uint32, error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts, configDigest)
}

func (_Verifier *VerifierCallerSession) LatestConfigDetails(configDigest [32]byte) (uint32, error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts, configDigest)
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

func (_Verifier *VerifierTransactor) ActivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "activateConfig", configDigest)
}

func (_Verifier *VerifierSession) ActivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, configDigest)
}

func (_Verifier *VerifierTransactorSession) ActivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, configDigest)
}

func (_Verifier *VerifierTransactor) DeactivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "deactivateConfig", configDigest)
}

func (_Verifier *VerifierSession) DeactivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, configDigest)
}

func (_Verifier *VerifierTransactorSession) DeactivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, configDigest)
}

func (_Verifier *VerifierTransactor) SetConfig(opts *bind.TransactOpts, configDigest [32]byte, signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "setConfig", configDigest, signers, f, recipientAddressesAndWeights)
}

func (_Verifier *VerifierSession) SetConfig(configDigest [32]byte, signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, configDigest, signers, f, recipientAddressesAndWeights)
}

func (_Verifier *VerifierTransactorSession) SetConfig(configDigest [32]byte, signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, configDigest, signers, f, recipientAddressesAndWeights)
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

func (_Verifier *VerifierTransactor) UpdateConfig(opts *bind.TransactOpts, configDigest [32]byte, prevSigners []common.Address, newSigners []common.Address, f uint8) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "updateConfig", configDigest, prevSigners, newSigners, f)
}

func (_Verifier *VerifierSession) UpdateConfig(configDigest [32]byte, prevSigners []common.Address, newSigners []common.Address, f uint8) (*types.Transaction, error) {
	return _Verifier.Contract.UpdateConfig(&_Verifier.TransactOpts, configDigest, prevSigners, newSigners, f)
}

func (_Verifier *VerifierTransactorSession) UpdateConfig(configDigest [32]byte, prevSigners []common.Address, newSigners []common.Address, f uint8) (*types.Transaction, error) {
	return _Verifier.Contract.UpdateConfig(&_Verifier.TransactOpts, configDigest, prevSigners, newSigners, f)
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
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_Verifier *VerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigActivatedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigActivated", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigActivatedIterator{contract: _Verifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

func (_Verifier *VerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigActivated, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigActivated", configDigestRule)
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
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_Verifier *VerifierFilterer) FilterConfigDeactivated(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigDeactivatedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigDeactivated", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigDeactivatedIterator{contract: _Verifier.contract, event: "ConfigDeactivated", logs: logs, sub: sub}, nil
}

func (_Verifier *VerifierFilterer) WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigDeactivated, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigDeactivated", configDigestRule)
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
	ConfigDigest [32]byte
	Signers      []common.Address
	F            uint8
	Raw          types.Log
}

func (_Verifier *VerifierFilterer) FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigSetIterator{contract: _Verifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_Verifier *VerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VerifierConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigSet", configDigestRule)
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

type VerifierConfigUpdatedIterator struct {
	Event *VerifierConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigUpdated)
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
		it.Event = new(VerifierConfigUpdated)
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

func (it *VerifierConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *VerifierConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierConfigUpdated struct {
	ConfigDigest [32]byte
	PrevSigners  []common.Address
	NewSigners   []common.Address
	Raw          types.Log
}

func (_Verifier *VerifierFilterer) FilterConfigUpdated(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigUpdatedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigUpdated", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &VerifierConfigUpdatedIterator{contract: _Verifier.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

func (_Verifier *VerifierFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *VerifierConfigUpdated, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigUpdated", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierConfigUpdated)
				if err := _Verifier.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_Verifier *VerifierFilterer) ParseConfigUpdated(log types.Log) (*VerifierConfigUpdated, error) {
	event := new(VerifierConfigUpdated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_Verifier *Verifier) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Verifier.abi.Events["ConfigActivated"].ID:
		return _Verifier.ParseConfigActivated(log)
	case _Verifier.abi.Events["ConfigDeactivated"].ID:
		return _Verifier.ParseConfigDeactivated(log)
	case _Verifier.abi.Events["ConfigSet"].ID:
		return _Verifier.ParseConfigSet(log)
	case _Verifier.abi.Events["ConfigUpdated"].ID:
		return _Verifier.ParseConfigUpdated(log)
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
	return common.HexToHash("0xa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea")
}

func (VerifierConfigDeactivated) Topic() common.Hash {
	return common.HexToHash("0x5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579")
}

func (VerifierConfigSet) Topic() common.Hash {
	return common.HexToHash("0x5b1f376eb2bda670fa39339616d0a73f45b61bec8faeba8ca834f2ebb49676e0")
}

func (VerifierConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0xb0b75a854fab801413da6202fc07e875c54eaf371a1e3909fb2645364ba58616")
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
	LatestConfigDetails(opts *bind.CallOpts, configDigest [32]byte) (uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ActivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	DeactivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, configDigest [32]byte, signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, configDigest [32]byte, prevSigners []common.Address, newSigners []common.Address, f uint8) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error)

	FilterConfigActivated(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigActivatedIterator, error)

	WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigActivated, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigActivated(log types.Log) (*VerifierConfigActivated, error)

	FilterConfigDeactivated(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigDeactivatedIterator, error)

	WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigDeactivated, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigDeactivated(log types.Log) (*VerifierConfigDeactivated, error)

	FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VerifierConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VerifierConfigSet, error)

	FilterConfigUpdated(opts *bind.FilterOpts, configDigest [][32]byte) (*VerifierConfigUpdatedIterator, error)

	WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *VerifierConfigUpdated, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigUpdated(log types.Log) (*VerifierConfigUpdated, error)

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
