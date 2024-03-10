// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_owner

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

type FeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}

type VRFTypesProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

type VRFTypesRequestCommitment struct {
	BlockNum         uint64
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
}

var VRFOwnerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptySendersList\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAllowedToSetSenders\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsForced\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptVRFOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRFTypes.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structVRFTypes.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVRFCoordinator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structFeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferVRFOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162002007380380620020078339810160408190526200003491620001fc565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000150565b5050506001600160a01b0381166200012a5760405162461bcd60e51b815260206004820152602860248201527f76726620636f6f7264696e61746f722061646472657373206d757374206265206044820152676e6f6e2d7a65726f60c01b606482015260840162000082565b600580546001600160a01b0319166001600160a01b03929092169190911790556200022e565b6001600160a01b038116331415620001ab5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200020f57600080fd5b81516001600160a01b03811681146200022757600080fd5b9392505050565b611dc9806200023e6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063984e804711610097578063e72f6e3011610066578063e72f6e30146101f9578063ee56997b1461020c578063f2fde38b1461021f578063fa00763a1461023257600080fd5b8063984e8047146101ad578063a378f371146101b5578063af198b97146101d3578063c2df03e4146101e657600080fd5b80634cb48a54116100d35780634cb48a54146101405780636f64f03f1461015357806379ba5097146101665780638da5cb5b1461016e57600080fd5b806302bcc5b6146100fa57806308821d581461010f5780632408afaa14610122575b600080fd5b61010d6101083660046118f6565b610255565b005b61010d61011d366004611591565b6102ee565b61012a61034c565b6040516101379190611a0a565b60405180910390f35b61010d61014e3660046116ed565b6103bb565b61010d6101613660046114e8565b61045d565b61010d6104f3565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610137565b61010d6105f5565b60055473ffffffffffffffffffffffffffffffffffffffff16610188565b61010d6101e13660046115c6565b610681565b61010d6101f43660046114cd565b610859565b61010d6102073660046114cd565b6108b9565b61010d61021a36600461151c565b610919565b61010d61022d3660046114cd565b610a8c565b6102456102403660046114cd565b610aa0565b6040519015158152602001610137565b61025d610ab3565b6005546040517f02bcc5b600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8316600482015273ffffffffffffffffffffffffffffffffffffffff909116906302bcc5b6906024015b600060405180830381600087803b1580156102d357600080fd5b505af11580156102e7573d6000803e3d6000fd5b5050505050565b6102f6610ab3565b6005546040517f08821d5800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116906308821d58906102b9908490600401611a64565b606060048054806020026020016040519081016040528092919081815260200182805480156103b157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610386575b5050505050905090565b6103c3610ab3565b6005546040517f4cb48a5400000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634cb48a549061042390899089908990899089908990600401611b82565b600060405180830381600087803b15801561043d57600080fd5b505af1158015610451573d6000803e3d6000fd5b50505050505050505050565b610465610ab3565b6005546040517f6f64f03f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636f64f03f906104bd9085908590600401611962565b600060405180830381600087803b1580156104d757600080fd5b505af11580156104eb573d6000803e3d6000fd5b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610579576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6105fd610ab3565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166379ba50976040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561066757600080fd5b505af115801561067b573d6000803e3d6000fd5b50505050565b610689610b36565b600061069d83600001518460800151610b75565b905060006106a9610c7f565b805160208083015160408051610120810182526000808252938101849052908101839052606081018390526080810183905260a0810183905260c0810183905260e08101839052610100810183905293945061072b9390916001917f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff906103c3565b6005546040517faf198b9700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063af198b97906107839087908790600401611a85565b602060405180830381600087803b15801561079d57600080fd5b505af11580156107b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107d59190611911565b506107fc816000015182602001518360400151846060015185608001518660a001516103c3565b826080015173ffffffffffffffffffffffffffffffffffffffff16836020015167ffffffffffffffff16837fabbcd646b939d78de3053d035798eb5c9818ea1836a2fbdbad335331df51e01d60405160405180910390a450505050565b610861610ab3565b6005546040517ff2fde38b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301529091169063f2fde38b906024016102b9565b6108c1610ab3565b6005546040517fe72f6e3000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301529091169063e72f6e30906024016102b9565b610921610fcb565b610957576040517fad77f06100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8061098e576040517f75158c3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b6004548110156109ee576109db600482815481106109b1576109b1611d2b565b60009182526020909120015460029073ffffffffffffffffffffffffffffffffffffffff16611009565b50806109e681611c94565b915050610991565b5060005b81811015610a3f57610a2c838383818110610a0f57610a0f611d2b565b9050602002016020810190610a2491906114cd565b600290611032565b5080610a3781611c94565b9150506109f2565b50610a4c600483836112bb565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0828233604051610a8093929190611992565b60405180910390a15050565b610a94610ab3565b610a9d81611054565b50565b6000610aad60028361114a565b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b34576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610570565b565b610b3f33610aa0565b610b34576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005546040517fcaf70c4a000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff9091169063caf70c4a90610bd0908790600401611a77565b60206040518083038186803b158015610be857600080fd5b505afa158015610bfc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c2091906115ad565b905060008184604051602001610c40929190918252602082015260400190565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152919052805160209091012095945050505050565b6040805160c080820183526000808352602080840182905283850182905260608085018390526080808601849052865161012081018852848152928301849052958201839052810182905293840181905260a080850182905291840181905260e08401819052610100840152810191909152600080600080600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c3f909d46040518163ffffffff1660e01b815260040160806040518083038186803b158015610d5f57600080fd5b505afa158015610d73573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d97919061168e565b93509350935093506000806000806000806000806000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635fbbc0d26040518163ffffffff1660e01b81526004016101206040518083038186803b158015610e1657600080fd5b505afa158015610e2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e4e9190611832565b9850985098509850985098509850985098506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663356dac716040518163ffffffff1660e01b815260040160206040518083038186803b158015610eca57600080fd5b505afa158015610ede573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f0291906115ad565b90506040518060c001604052808f61ffff1681526020018e63ffffffff1681526020018d63ffffffff1681526020018c63ffffffff1681526020018281526020016040518061012001604052808d63ffffffff1681526020018c63ffffffff1681526020018b63ffffffff1681526020018a63ffffffff1681526020018963ffffffff1681526020018862ffffff1681526020018762ffffff1681526020018662ffffff1681526020018562ffffff168152508152509e50505050505050505050505050505090565b600033610fed60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b600061102b8373ffffffffffffffffffffffffffffffffffffffff8416611179565b9392505050565b600061102b8373ffffffffffffffffffffffffffffffffffffffff841661126c565b73ffffffffffffffffffffffffffffffffffffffff81163314156110d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610570565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600183016020526040812054151561102b565b6000818152600183016020526040812054801561126257600061119d600183611c7d565b85549091506000906111b190600190611c7d565b90508181146112165760008660000182815481106111d1576111d1611d2b565b90600052602060002001549050808760000184815481106111f4576111f4611d2b565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061122757611227611cfc565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610aad565b6000915050610aad565b60008181526001830160205260408120546112b357508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610aad565b506000610aad565b828054828255906000526020600020908101928215611333579160200282015b828111156113335781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8435161782556020909201916001909101906112db565b5061133f929150611343565b5090565b5b8082111561133f5760008155600101611344565b803573ffffffffffffffffffffffffffffffffffffffff8116811461137c57600080fd5b919050565b8060408101831015610aad57600080fd5b600082601f8301126113a357600080fd5b6040516040810181811067ffffffffffffffff821117156113c6576113c6611d5a565b80604052508083856040860111156113dd57600080fd5b60005b60028110156113ff5781358352602092830192909101906001016113e0565b509195945050505050565b600060a0828403121561141c57600080fd5b60405160a0810181811067ffffffffffffffff8211171561143f5761143f611d5a565b60405290508061144e836114b5565b815261145c602084016114b5565b6020820152604083013561146f81611daa565b6040820152606083013561148281611daa565b606082015261149360808401611358565b60808201525092915050565b803561137c81611d99565b803561137c81611daa565b803567ffffffffffffffff8116811461137c57600080fd5b6000602082840312156114df57600080fd5b61102b82611358565b600080606083850312156114fb57600080fd5b61150483611358565b91506115138460208501611381565b90509250929050565b6000806020838503121561152f57600080fd5b823567ffffffffffffffff8082111561154757600080fd5b818501915085601f83011261155b57600080fd5b81358181111561156a57600080fd5b8660208260051b850101111561157f57600080fd5b60209290920196919550909350505050565b6000604082840312156115a357600080fd5b61102b8383611381565b6000602082840312156115bf57600080fd5b5051919050565b6000808284036102408112156115db57600080fd5b6101a0808212156115eb57600080fd5b6115f3611c53565b91506115ff8686611392565b825261160e8660408701611392565b60208301526080850135604083015260a0850135606083015260c0850135608083015261163d60e08601611358565b60a083015261010061165187828801611392565b60c0840152611664876101408801611392565b60e084015261018086013581840152508193506116838682870161140a565b925050509250929050565b600080600080608085870312156116a457600080fd5b84516116af81611d89565b60208601519094506116c081611daa565b60408601519093506116d181611daa565b60608601519092506116e281611daa565b939692955090935050565b6000806000806000808688036101c081121561170857600080fd5b873561171381611d89565b9650602088013561172381611daa565b9550604088013561173381611daa565b9450606088013561174381611daa565b9350608088013592506101207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60820181131561177e57600080fd5b611786611c53565b915061179460a08a016114aa565b82526117a260c08a016114aa565b60208301526117b360e08a016114aa565b60408301526101006117c6818b016114aa565b60608401526117d6828b016114aa565b60808401526117e86101408b0161149f565b60a08401526117fa6101608b0161149f565b60c084015261180c6101808b0161149f565b60e084015261181e6101a08b0161149f565b818401525050809150509295509295509295565b60008060008060008060008060006101208a8c03121561185157600080fd5b895161185c81611daa565b60208b015190995061186d81611daa565b60408b015190985061187e81611daa565b60608b015190975061188f81611daa565b60808b01519096506118a081611daa565b60a08b01519095506118b181611d99565b60c08b01519094506118c281611d99565b60e08b01519093506118d381611d99565b6101008b01519092506118e581611d99565b809150509295985092959850929598565b60006020828403121561190857600080fd5b61102b826114b5565b60006020828403121561192357600080fd5b81516bffffffffffffffffffffffff8116811461102b57600080fd5b8060005b600281101561067b578151845260209384019390910190600101611943565b73ffffffffffffffffffffffffffffffffffffffff83168152606081016040836020840137600081529392505050565b6040808252810183905260008460608301825b868110156119e05773ffffffffffffffffffffffffffffffffffffffff6119cb84611358565b168252602092830192909101906001016119a5565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b6020808252825182820181905260009190848201906040850190845b81811015611a5857835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611a26565b50909695505050505050565b6040818101908383376000815292915050565b60408101610aad828461193f565b600061024082019050611a9982855161193f565b6020840151611aab604084018261193f565b5060408401516080830152606084015160a0830152608084015160c083015273ffffffffffffffffffffffffffffffffffffffff60a08501511660e083015260c0840151610100611afe8185018361193f565b60e08601519150611b1361014085018361193f565b85015161018084015250825167ffffffffffffffff9081166101a08401526020840151166101c0830152604083015163ffffffff9081166101e0840152606084015116610200830152608083015173ffffffffffffffffffffffffffffffffffffffff1661022083015261102b565b60006101c08201905061ffff8816825263ffffffff8088166020840152808716604084015280861660608401528460808401528084511660a08401528060208501511660c0840152506040830151611be260e084018263ffffffff169052565b506060830151610100611bfc8185018363ffffffff169052565b608085015163ffffffff1661012085015260a085015162ffffff90811661014086015260c0860151811661016086015260e086015181166101808601529401519093166101a0909201919091529695505050505050565b604051610120810167ffffffffffffffff81118282101715611c7757611c77611d5a565b60405290565b600082821015611c8f57611c8f611ccd565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611cc657611cc6611ccd565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61ffff81168114610a9d57600080fd5b62ffffff81168114610a9d57600080fd5b63ffffffff81168114610a9d57600080fdfea164736f6c6343000806000a",
}

var VRFOwnerABI = VRFOwnerMetaData.ABI

var VRFOwnerBin = VRFOwnerMetaData.Bin

func DeployVRFOwner(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFOwner, error) {
	parsed, err := VRFOwnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFOwnerBin), backend, _vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFOwner{address: address, abi: *parsed, VRFOwnerCaller: VRFOwnerCaller{contract: contract}, VRFOwnerTransactor: VRFOwnerTransactor{contract: contract}, VRFOwnerFilterer: VRFOwnerFilterer{contract: contract}}, nil
}

type VRFOwner struct {
	address common.Address
	abi     abi.ABI
	VRFOwnerCaller
	VRFOwnerTransactor
	VRFOwnerFilterer
}

type VRFOwnerCaller struct {
	contract *bind.BoundContract
}

type VRFOwnerTransactor struct {
	contract *bind.BoundContract
}

type VRFOwnerFilterer struct {
	contract *bind.BoundContract
}

type VRFOwnerSession struct {
	Contract     *VRFOwner
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFOwnerCallerSession struct {
	Contract *VRFOwnerCaller
	CallOpts bind.CallOpts
}

type VRFOwnerTransactorSession struct {
	Contract     *VRFOwnerTransactor
	TransactOpts bind.TransactOpts
}

type VRFOwnerRaw struct {
	Contract *VRFOwner
}

type VRFOwnerCallerRaw struct {
	Contract *VRFOwnerCaller
}

type VRFOwnerTransactorRaw struct {
	Contract *VRFOwnerTransactor
}

func NewVRFOwner(address common.Address, backend bind.ContractBackend) (*VRFOwner, error) {
	abi, err := abi.JSON(strings.NewReader(VRFOwnerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFOwner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFOwner{address: address, abi: abi, VRFOwnerCaller: VRFOwnerCaller{contract: contract}, VRFOwnerTransactor: VRFOwnerTransactor{contract: contract}, VRFOwnerFilterer: VRFOwnerFilterer{contract: contract}}, nil
}

func NewVRFOwnerCaller(address common.Address, caller bind.ContractCaller) (*VRFOwnerCaller, error) {
	contract, err := bindVRFOwner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerCaller{contract: contract}, nil
}

func NewVRFOwnerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFOwnerTransactor, error) {
	contract, err := bindVRFOwner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerTransactor{contract: contract}, nil
}

func NewVRFOwnerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFOwnerFilterer, error) {
	contract, err := bindVRFOwner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerFilterer{contract: contract}, nil
}

func bindVRFOwner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFOwnerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFOwner *VRFOwnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFOwner.Contract.VRFOwnerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFOwner *VRFOwnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFOwner.Contract.VRFOwnerTransactor.contract.Transfer(opts)
}

func (_VRFOwner *VRFOwnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFOwner.Contract.VRFOwnerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFOwner *VRFOwnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFOwner.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFOwner *VRFOwnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFOwner.Contract.contract.Transfer(opts)
}

func (_VRFOwner *VRFOwnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFOwner.Contract.contract.Transact(opts, method, params...)
}

func (_VRFOwner *VRFOwnerCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _VRFOwner.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_VRFOwner *VRFOwnerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _VRFOwner.Contract.GetAuthorizedSenders(&_VRFOwner.CallOpts)
}

func (_VRFOwner *VRFOwnerCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _VRFOwner.Contract.GetAuthorizedSenders(&_VRFOwner.CallOpts)
}

func (_VRFOwner *VRFOwnerCaller) GetVRFCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFOwner.contract.Call(opts, &out, "getVRFCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFOwner *VRFOwnerSession) GetVRFCoordinator() (common.Address, error) {
	return _VRFOwner.Contract.GetVRFCoordinator(&_VRFOwner.CallOpts)
}

func (_VRFOwner *VRFOwnerCallerSession) GetVRFCoordinator() (common.Address, error) {
	return _VRFOwner.Contract.GetVRFCoordinator(&_VRFOwner.CallOpts)
}

func (_VRFOwner *VRFOwnerCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _VRFOwner.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFOwner *VRFOwnerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _VRFOwner.Contract.IsAuthorizedSender(&_VRFOwner.CallOpts, sender)
}

func (_VRFOwner *VRFOwnerCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _VRFOwner.Contract.IsAuthorizedSender(&_VRFOwner.CallOpts, sender)
}

func (_VRFOwner *VRFOwnerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFOwner.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFOwner *VRFOwnerSession) Owner() (common.Address, error) {
	return _VRFOwner.Contract.Owner(&_VRFOwner.CallOpts)
}

func (_VRFOwner *VRFOwnerCallerSession) Owner() (common.Address, error) {
	return _VRFOwner.Contract.Owner(&_VRFOwner.CallOpts)
}

func (_VRFOwner *VRFOwnerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "acceptOwnership")
}

func (_VRFOwner *VRFOwnerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFOwner.Contract.AcceptOwnership(&_VRFOwner.TransactOpts)
}

func (_VRFOwner *VRFOwnerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFOwner.Contract.AcceptOwnership(&_VRFOwner.TransactOpts)
}

func (_VRFOwner *VRFOwnerTransactor) AcceptVRFOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "acceptVRFOwnership")
}

func (_VRFOwner *VRFOwnerSession) AcceptVRFOwnership() (*types.Transaction, error) {
	return _VRFOwner.Contract.AcceptVRFOwnership(&_VRFOwner.TransactOpts)
}

func (_VRFOwner *VRFOwnerTransactorSession) AcceptVRFOwnership() (*types.Transaction, error) {
	return _VRFOwner.Contract.AcceptVRFOwnership(&_VRFOwner.TransactOpts)
}

func (_VRFOwner *VRFOwnerTransactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFOwner *VRFOwnerSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFOwner.Contract.DeregisterProvingKey(&_VRFOwner.TransactOpts, publicProvingKey)
}

func (_VRFOwner *VRFOwnerTransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFOwner.Contract.DeregisterProvingKey(&_VRFOwner.TransactOpts, publicProvingKey)
}

func (_VRFOwner *VRFOwnerTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFTypesProof, rc VRFTypesRequestCommitment) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_VRFOwner *VRFOwnerSession) FulfillRandomWords(proof VRFTypesProof, rc VRFTypesRequestCommitment) (*types.Transaction, error) {
	return _VRFOwner.Contract.FulfillRandomWords(&_VRFOwner.TransactOpts, proof, rc)
}

func (_VRFOwner *VRFOwnerTransactorSession) FulfillRandomWords(proof VRFTypesProof, rc VRFTypesRequestCommitment) (*types.Transaction, error) {
	return _VRFOwner.Contract.FulfillRandomWords(&_VRFOwner.TransactOpts, proof, rc)
}

func (_VRFOwner *VRFOwnerTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFOwner *VRFOwnerSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFOwner.Contract.OwnerCancelSubscription(&_VRFOwner.TransactOpts, subId)
}

func (_VRFOwner *VRFOwnerTransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFOwner.Contract.OwnerCancelSubscription(&_VRFOwner.TransactOpts, subId)
}

func (_VRFOwner *VRFOwnerTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFOwner *VRFOwnerSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.RecoverFunds(&_VRFOwner.TransactOpts, to)
}

func (_VRFOwner *VRFOwnerTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.RecoverFunds(&_VRFOwner.TransactOpts, to)
}

func (_VRFOwner *VRFOwnerTransactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_VRFOwner *VRFOwnerSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFOwner.Contract.RegisterProvingKey(&_VRFOwner.TransactOpts, oracle, publicProvingKey)
}

func (_VRFOwner *VRFOwnerTransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFOwner.Contract.RegisterProvingKey(&_VRFOwner.TransactOpts, oracle, publicProvingKey)
}

func (_VRFOwner *VRFOwnerTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_VRFOwner *VRFOwnerSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.SetAuthorizedSenders(&_VRFOwner.TransactOpts, senders)
}

func (_VRFOwner *VRFOwnerTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.SetAuthorizedSenders(&_VRFOwner.TransactOpts, senders)
}

func (_VRFOwner *VRFOwnerTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig FeeConfig) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFOwner *VRFOwnerSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig FeeConfig) (*types.Transaction, error) {
	return _VRFOwner.Contract.SetConfig(&_VRFOwner.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFOwner *VRFOwnerTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig FeeConfig) (*types.Transaction, error) {
	return _VRFOwner.Contract.SetConfig(&_VRFOwner.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFOwner *VRFOwnerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFOwner *VRFOwnerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.TransferOwnership(&_VRFOwner.TransactOpts, to)
}

func (_VRFOwner *VRFOwnerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.TransferOwnership(&_VRFOwner.TransactOpts, to)
}

func (_VRFOwner *VRFOwnerTransactor) TransferVRFOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFOwner.contract.Transact(opts, "transferVRFOwnership", to)
}

func (_VRFOwner *VRFOwnerSession) TransferVRFOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.TransferVRFOwnership(&_VRFOwner.TransactOpts, to)
}

func (_VRFOwner *VRFOwnerTransactorSession) TransferVRFOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFOwner.Contract.TransferVRFOwnership(&_VRFOwner.TransactOpts, to)
}

type VRFOwnerAuthorizedSendersChangedIterator struct {
	Event *VRFOwnerAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFOwnerAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFOwnerAuthorizedSendersChanged)
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
		it.Event = new(VRFOwnerAuthorizedSendersChanged)
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

func (it *VRFOwnerAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *VRFOwnerAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFOwnerAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_VRFOwner *VRFOwnerFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*VRFOwnerAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _VRFOwner.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &VRFOwnerAuthorizedSendersChangedIterator{contract: _VRFOwner.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_VRFOwner *VRFOwnerFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *VRFOwnerAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _VRFOwner.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFOwnerAuthorizedSendersChanged)
				if err := _VRFOwner.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_VRFOwner *VRFOwnerFilterer) ParseAuthorizedSendersChanged(log types.Log) (*VRFOwnerAuthorizedSendersChanged, error) {
	event := new(VRFOwnerAuthorizedSendersChanged)
	if err := _VRFOwner.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFOwnerOwnershipTransferRequestedIterator struct {
	Event *VRFOwnerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFOwnerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFOwnerOwnershipTransferRequested)
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
		it.Event = new(VRFOwnerOwnershipTransferRequested)
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

func (it *VRFOwnerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFOwnerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFOwnerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFOwner *VRFOwnerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFOwnerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFOwner.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerOwnershipTransferRequestedIterator{contract: _VRFOwner.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFOwner *VRFOwnerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFOwner.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFOwnerOwnershipTransferRequested)
				if err := _VRFOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFOwner *VRFOwnerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFOwnerOwnershipTransferRequested, error) {
	event := new(VRFOwnerOwnershipTransferRequested)
	if err := _VRFOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFOwnerOwnershipTransferredIterator struct {
	Event *VRFOwnerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFOwnerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFOwnerOwnershipTransferred)
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
		it.Event = new(VRFOwnerOwnershipTransferred)
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

func (it *VRFOwnerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFOwnerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFOwnerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFOwner *VRFOwnerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFOwnerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFOwner.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerOwnershipTransferredIterator{contract: _VRFOwner.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFOwner *VRFOwnerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFOwner.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFOwnerOwnershipTransferred)
				if err := _VRFOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFOwner *VRFOwnerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFOwnerOwnershipTransferred, error) {
	event := new(VRFOwnerOwnershipTransferred)
	if err := _VRFOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFOwnerRandomWordsForcedIterator struct {
	Event *VRFOwnerRandomWordsForced

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFOwnerRandomWordsForcedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFOwnerRandomWordsForced)
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
		it.Event = new(VRFOwnerRandomWordsForced)
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

func (it *VRFOwnerRandomWordsForcedIterator) Error() error {
	return it.fail
}

func (it *VRFOwnerRandomWordsForcedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFOwnerRandomWordsForced struct {
	RequestId *big.Int
	SubId     uint64
	Sender    common.Address
	Raw       types.Log
}

func (_VRFOwner *VRFOwnerFilterer) FilterRandomWordsForced(opts *bind.FilterOpts, requestId []*big.Int, subId []uint64, sender []common.Address) (*VRFOwnerRandomWordsForcedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFOwner.contract.FilterLogs(opts, "RandomWordsForced", requestIdRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerRandomWordsForcedIterator{contract: _VRFOwner.contract, event: "RandomWordsForced", logs: logs, sub: sub}, nil
}

func (_VRFOwner *VRFOwnerFilterer) WatchRandomWordsForced(opts *bind.WatchOpts, sink chan<- *VRFOwnerRandomWordsForced, requestId []*big.Int, subId []uint64, sender []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFOwner.contract.WatchLogs(opts, "RandomWordsForced", requestIdRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFOwnerRandomWordsForced)
				if err := _VRFOwner.contract.UnpackLog(event, "RandomWordsForced", log); err != nil {
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

func (_VRFOwner *VRFOwnerFilterer) ParseRandomWordsForced(log types.Log) (*VRFOwnerRandomWordsForced, error) {
	event := new(VRFOwnerRandomWordsForced)
	if err := _VRFOwner.contract.UnpackLog(event, "RandomWordsForced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFOwner *VRFOwner) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFOwner.abi.Events["AuthorizedSendersChanged"].ID:
		return _VRFOwner.ParseAuthorizedSendersChanged(log)
	case _VRFOwner.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFOwner.ParseOwnershipTransferRequested(log)
	case _VRFOwner.abi.Events["OwnershipTransferred"].ID:
		return _VRFOwner.ParseOwnershipTransferred(log)
	case _VRFOwner.abi.Events["RandomWordsForced"].ID:
		return _VRFOwner.ParseRandomWordsForced(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFOwnerAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (VRFOwnerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFOwnerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFOwnerRandomWordsForced) Topic() common.Hash {
	return common.HexToHash("0xabbcd646b939d78de3053d035798eb5c9818ea1836a2fbdbad335331df51e01d")
}

func (_VRFOwner *VRFOwner) Address() common.Address {
	return _VRFOwner.address
}

type VRFOwnerInterface interface {
	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetVRFCoordinator(opts *bind.CallOpts) (common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptVRFOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFTypesProof, rc VRFTypesRequestCommitment) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig FeeConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferVRFOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*VRFOwnerAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *VRFOwnerAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*VRFOwnerAuthorizedSendersChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFOwnerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFOwnerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFOwnerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFOwnerOwnershipTransferred, error)

	FilterRandomWordsForced(opts *bind.FilterOpts, requestId []*big.Int, subId []uint64, sender []common.Address) (*VRFOwnerRandomWordsForcedIterator, error)

	WatchRandomWordsForced(opts *bind.WatchOpts, sink chan<- *VRFOwnerRandomWordsForced, requestId []*big.Int, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsForced(log types.Log) (*VRFOwnerRandomWordsForced, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
