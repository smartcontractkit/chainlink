// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_registration_requests_wrapper

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var UpkeepRegistrationRequestsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"RegistrationRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"windowStart\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"approvedInCurrentWindow\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200220638038062002206833981810160405260408110156200003757600080fd5b508051602090910151338060008162000097576040805162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f0000000000000000604482015290519081900360640190fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000ca57620000ca81620000e9565b50505060609190911b6001600160601b03191660805260025562000199565b6001600160a01b03811633141562000148576040805162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60805160601c61203b620001cb600039806108ad5280610c83528061101852806116e35280611a0e525061203b6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806388b12d5511610081578063c4110e5c1161005b578063c4110e5c14610473578063c4d252f514610605578063f2fde38b14610622576100d4565b806388b12d551461037f5780638da5cb5b146103d9578063a4c0ed36146103e1576100d4565b80635772ac92116100b25780635772ac92146102b357806379ba50971461030a578063850af0cb14610312576100d4565b8063181f5a77146100d9578063183310b3146101565780631b6b6d2314610282575b600080fd5b6100e1610655565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561011b578181015183820152602001610103565b50505050905090810190601f1680156101485780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610280600480360360c081101561016c57600080fd5b81019060208101813564010000000081111561018757600080fd5b82018360208201111561019957600080fd5b803590602001918460018302840111640100000000831117156101bb57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929573ffffffffffffffffffffffffffffffffffffffff853581169663ffffffff602088013516966040810135909216955091935090915060808101906060013564010000000081111561024157600080fd5b82018360208201111561025357600080fd5b8035906020019184600183028401116401000000008311171561027557600080fd5b91935091503561068e565b005b61028a6108ab565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610280600480360360a08110156102c957600080fd5b50803515159063ffffffff6020820135169061ffff6040820135169073ffffffffffffffffffffffffffffffffffffffff60608201351690608001356108cf565b610280610a4c565b61031a610b4e565b60408051971515885263ffffffff909616602088015261ffff9485168787015273ffffffffffffffffffffffffffffffffffffffff9093166060870152608086019190915267ffffffffffffffff1660a08501521660c0830152519081900360e00190f35b61039c6004803603602081101561039557600080fd5b5035610be9565b6040805173ffffffffffffffffffffffffffffffffffffffff90931683526bffffffffffffffffffffffff90911660208301528051918290030190f35b61028a610c4f565b610280600480360360608110156103f757600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561043457600080fd5b82018360208201111561044657600080fd5b8035906020019184600183028401116401000000008311171561046857600080fd5b509092509050610c6b565b610280600480360361010081101561048a57600080fd5b8101906020810181356401000000008111156104a557600080fd5b8201836020820111156104b757600080fd5b803590602001918460018302840111640100000000831117156104d957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929594936020810193503591505064010000000081111561052c57600080fd5b82018360208201111561053e57600080fd5b8035906020019184600183028401116401000000008311171561056057600080fd5b9193909273ffffffffffffffffffffffffffffffffffffffff833581169363ffffffff60208201351693604082013590921692906080810190606001356401000000008111156105af57600080fd5b8201836020820111156105c157600080fd5b803590602001918460018302840111640100000000831117156105e357600080fd5b919350915080356bffffffffffffffffffffffff16906020013560ff16611000565b6102806004803603602081101561061b57600080fd5b50356114e2565b6102806004803603602081101561063857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166117f0565b6040518060400160405280602081526020017f55706b656570526567697374726174696f6e526571756573747320312e302e3081525081565b610696611804565b60008181526003602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169183019190915261076457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015290519081900360640190fd5b60008787878787604051602001808673ffffffffffffffffffffffffffffffffffffffff1681526020018563ffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001806020018281038252848482818152602001925080828437600081840152601f19601f820116905080830192505050965050505050505060405160208183030381529060405280519060200120905080831461087357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6861736820616e64207061796c6f616420646f206e6f74206d61746368000000604482015290519081900360640190fd5b60008381526003602090815260408220919091558201516108a0908a908a908a908a908a908a908a61188c565b505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6108d7611804565b6040805160a0808201835287151580835261ffff8716602080850182905263ffffffff8a1685870181905260006060808801829052608097880191909152600480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff166101008602177fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff1663010000008402177fffffffffffffffffffffffffffffff00000000000000000000ffffffffffffff1690556002899055600580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8c16908117909155885195865292850191909152838701929092529082015291820184905291517f421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53929181900390910190a15050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610ad257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a08101825260045460ff8116151580835261ffff610100830481166020850181905263ffffffff630100000085041695850186905267ffffffffffffffff670100000000000000850416606086018190526f0100000000000000000000000000000090940490911660809094018490526005546002549296919473ffffffffffffffffffffffffffffffffffffffff9091169391565b60009081526003602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169290910182905291565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610d0f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101517fffffffff0000000000000000000000000000000000000000000000000000000081167fc4110e5c0000000000000000000000000000000000000000000000000000000014610dfa57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b8484848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060e4810151828114610ea457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f416d6f756e74206d69736d617463680000000000000000000000000000000000604482015290519081900360640190fd5b600254881015610f1557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e73756666696369656e74207061796d656e74000000000000000000000000604482015290519081900360640190fd5b60003073ffffffffffffffffffffffffffffffffffffffff1688886040518083838082843760405192019450600093509091505080830381855af49150503d8060008114610f7f576040519150601f19603f3d011682016040523d82523d6000602084013e610f84565b606091505b5050905080610ff457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f556e61626c6520746f2063726561746520726571756573740000000000000000604482015290519081900360640190fd5b50505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110a457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff851661112657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f696e76616c69642061646d696e20616464726573730000000000000000000000604482015290519081900360640190fd5b60008787878787604051602001808673ffffffffffffffffffffffffffffffffffffffff1681526020018563ffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505096505050505050506040516020818303038152906040528051906020012090508160ff168873ffffffffffffffffffffffffffffffffffffffff16827fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788e8e8e8d8d8d8d8d6040518080602001806020018863ffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff16815260200180602001856bffffffffffffffffffffffff16815260200184810384528c818151815260200191508051906020019080838360005b8381101561128d578181015183820152602001611275565b50505050905090810190601f1680156112ba5780820380516001836020036101000a031916815260200191505b5084810383528a81526020018b8b80828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528781526020019050878780828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039d50909b505050505050505050505050a46040805160a08101825260045460ff811615801580845261ffff61010084048116602086015263ffffffff63010000008504169585019590955267ffffffffffffffff67010000000000000084041660608501526f01000000000000000000000000000000909204909316608083015290916113d457506113d481611bf6565b156113f7576113e281611c2a565b6113f28c8a8a8a8a8a8a8961188c565b6114d4565b600082815260036020526040812054611436907401000000000000000000000000000000000000000090046bffffffffffffffffffffffff1686611d5e565b60408051808201825273ffffffffffffffffffffffffffffffffffffffff8b811682526bffffffffffffffffffffffff938416602080840191825260008981526003909152939093209151825493517fffffffffffffffffffffffff00000000000000000000000000000000000000009094169082161716740100000000000000000000000000000000000000009290931691909102919091179055505b505050505050505050505050565b60008181526003602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff16918301919091523314806115845750611555610c4f565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b6115ef57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792061646d696e202f206f776e65722063616e2063616e63656c000000604482015290519081900360640190fd5b805173ffffffffffffffffffffffffffffffffffffffff1661167257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015290519081900360640190fd5b60008281526003602090815260408083208390558382015181517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff9091166024820152905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb93604480850194919392918390030190829087803b15801561172a57600080fd5b505af115801561173e573d6000803e3d6000fd5b505050506040513d602081101561175457600080fd5b50516117c157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4c494e4b20746f6b656e207472616e73666572206661696c6564000000000000604482015290519081900360640190fd5b60405182907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a25050565b6117f8611804565b61180181611dea565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461188a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b565b6005546040517fda5c674100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8981166004830190815263ffffffff8a1660248401528882166044840152608060648401908152608484018890529190931692600092849263da5c6741928d928d928d928d928d929060a401848480828437600081840152601f19601f8201169050808301925050509650505050505050602060405180830381600087803b15801561195557600080fd5b505af1158015611969573d6000803e3d6000fd5b505050506040513d602081101561197f57600080fd5b5051604080516020808201849052825180830382018152828401938490527f4000aea00000000000000000000000000000000000000000000000000000000090935273ffffffffffffffffffffffffffffffffffffffff868116604484019081526bffffffffffffffffffffffff8a166064850152606060848501908152855160a486015285519697506000967f000000000000000000000000000000000000000000000000000000000000000090931695634000aea0958a958d959294939260c490920191908501908083838d5b83811015611a66578181015183820152602001611a4e565b50505050905090810190601f168015611a935780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015611ab457600080fd5b505af1158015611ac8573d6000803e3d6000fd5b505050506040513d6020811015611ade57600080fd5b5051905080611b4e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f6661696c656420746f2066756e642075706b6565700000000000000000000000604482015290519081900360640190fd5b81847fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b8d6040518080602001828103825283818151815260200191508051906020019080838360005b83811015611baf578181015183820152602001611b97565b50505050905090810190601f168015611bdc5780820380516001836020036101000a031916815260200191505b509250505060405180910390a35050505050505050505050565b6000611c0182611ee5565b816020015161ffff16826080015161ffff161015611c2157506001611c25565b5060005b919050565b60808101805161ffff6001909101811691829052825160048054602086015160408701516060909701516f010000000000000000000000000000009096027fffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffff67ffffffffffffffff909716670100000000000000027fffffffffffffffffffffffffffffffffff0000000000000000ffffffffffffff63ffffffff9099166301000000027fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff93909716610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff9615157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009095169490941795909516929092171693909317949094161791909116179055565b60008282016bffffffffffffffffffffffff8085169082161015611de357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415611e6f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000816060015167ffffffffffffffff1643039050816040015163ffffffff168167ffffffffffffffff161061202a574367ffffffffffffffff166060830181905260006080840152825160048054602086015160408701517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00909216931515939093177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff1661010061ffff90941693909302929092177fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff16630100000063ffffffff90931692909202919091177fffffffffffffffffffffffffffffffffff0000000000000000ffffffffffffff16670100000000000000909202919091177fffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffff1690555b505056fea164736f6c6343000706000a",
}

var UpkeepRegistrationRequestsABI = UpkeepRegistrationRequestsMetaData.ABI

var UpkeepRegistrationRequestsBin = UpkeepRegistrationRequestsMetaData.Bin

func DeployUpkeepRegistrationRequests(auth *bind.TransactOpts, backend bind.ContractBackend, LINKAddress common.Address, minimumLINKJuels *big.Int) (common.Address, *types.Transaction, *UpkeepRegistrationRequests, error) {
	parsed, err := UpkeepRegistrationRequestsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepRegistrationRequestsBin), backend, LINKAddress, minimumLINKJuels)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepRegistrationRequests{UpkeepRegistrationRequestsCaller: UpkeepRegistrationRequestsCaller{contract: contract}, UpkeepRegistrationRequestsTransactor: UpkeepRegistrationRequestsTransactor{contract: contract}, UpkeepRegistrationRequestsFilterer: UpkeepRegistrationRequestsFilterer{contract: contract}}, nil
}

type UpkeepRegistrationRequests struct {
	address common.Address
	abi     abi.ABI
	UpkeepRegistrationRequestsCaller
	UpkeepRegistrationRequestsTransactor
	UpkeepRegistrationRequestsFilterer
}

type UpkeepRegistrationRequestsCaller struct {
	contract *bind.BoundContract
}

type UpkeepRegistrationRequestsTransactor struct {
	contract *bind.BoundContract
}

type UpkeepRegistrationRequestsFilterer struct {
	contract *bind.BoundContract
}

type UpkeepRegistrationRequestsSession struct {
	Contract     *UpkeepRegistrationRequests
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepRegistrationRequestsCallerSession struct {
	Contract *UpkeepRegistrationRequestsCaller
	CallOpts bind.CallOpts
}

type UpkeepRegistrationRequestsTransactorSession struct {
	Contract     *UpkeepRegistrationRequestsTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepRegistrationRequestsRaw struct {
	Contract *UpkeepRegistrationRequests
}

type UpkeepRegistrationRequestsCallerRaw struct {
	Contract *UpkeepRegistrationRequestsCaller
}

type UpkeepRegistrationRequestsTransactorRaw struct {
	Contract *UpkeepRegistrationRequestsTransactor
}

func NewUpkeepRegistrationRequests(address common.Address, backend bind.ContractBackend) (*UpkeepRegistrationRequests, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepRegistrationRequestsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepRegistrationRequests(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequests{address: address, abi: abi, UpkeepRegistrationRequestsCaller: UpkeepRegistrationRequestsCaller{contract: contract}, UpkeepRegistrationRequestsTransactor: UpkeepRegistrationRequestsTransactor{contract: contract}, UpkeepRegistrationRequestsFilterer: UpkeepRegistrationRequestsFilterer{contract: contract}}, nil
}

func NewUpkeepRegistrationRequestsCaller(address common.Address, caller bind.ContractCaller) (*UpkeepRegistrationRequestsCaller, error) {
	contract, err := bindUpkeepRegistrationRequests(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsCaller{contract: contract}, nil
}

func NewUpkeepRegistrationRequestsTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepRegistrationRequestsTransactor, error) {
	contract, err := bindUpkeepRegistrationRequests(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsTransactor{contract: contract}, nil
}

func NewUpkeepRegistrationRequestsFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepRegistrationRequestsFilterer, error) {
	contract, err := bindUpkeepRegistrationRequests(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsFilterer{contract: contract}, nil
}

func bindUpkeepRegistrationRequests(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepRegistrationRequestsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepRegistrationRequests.Contract.UpkeepRegistrationRequestsCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.UpkeepRegistrationRequestsTransactor.contract.Transfer(opts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.UpkeepRegistrationRequestsTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepRegistrationRequests.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.contract.Transfer(opts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) LINK() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.LINK(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) LINK() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.LINK(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) GetPendingRequest(opts *bind.CallOpts, hash [32]byte) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "getPendingRequest", hash)

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _UpkeepRegistrationRequests.Contract.GetPendingRequest(&_UpkeepRegistrationRequests.CallOpts, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _UpkeepRegistrationRequests.Contract.GetPendingRequest(&_UpkeepRegistrationRequests.CallOpts, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

	error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "getRegistrationConfig")

	outstruct := new(GetRegistrationConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Enabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.WindowSizeInBlocks = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.AllowedPerWindow = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.KeeperRegistry = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.MinLINKJuels = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.WindowStart = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.ApprovedInCurrentWindow = *abi.ConvertType(out[6], new(uint16)).(*uint16)

	return *outstruct, err

}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _UpkeepRegistrationRequests.Contract.GetRegistrationConfig(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _UpkeepRegistrationRequests.Contract.GetRegistrationConfig(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Owner() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.Owner(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) Owner() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.Owner(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) TypeAndVersion() (string, error) {
	return _UpkeepRegistrationRequests.Contract.TypeAndVersion(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) TypeAndVersion() (string, error) {
	return _UpkeepRegistrationRequests.Contract.TypeAndVersion(&_UpkeepRegistrationRequests.CallOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "acceptOwnership")
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) AcceptOwnership() (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.AcceptOwnership(&_UpkeepRegistrationRequests.TransactOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.AcceptOwnership(&_UpkeepRegistrationRequests.TransactOpts)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "approve", name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Approve(&_UpkeepRegistrationRequests.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Approve(&_UpkeepRegistrationRequests.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "cancel", hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Cancel(&_UpkeepRegistrationRequests.TransactOpts, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Cancel(&_UpkeepRegistrationRequests.TransactOpts, hash)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.OnTokenTransfer(&_UpkeepRegistrationRequests.TransactOpts, arg0, amount, data)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.OnTokenTransfer(&_UpkeepRegistrationRequests.TransactOpts, arg0, amount, data)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "register", name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Register(&_UpkeepRegistrationRequests.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Register(&_UpkeepRegistrationRequests.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) SetRegistrationConfig(opts *bind.TransactOpts, enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "setRegistrationConfig", enabled, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) SetRegistrationConfig(enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.SetRegistrationConfig(&_UpkeepRegistrationRequests.TransactOpts, enabled, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) SetRegistrationConfig(enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.SetRegistrationConfig(&_UpkeepRegistrationRequests.TransactOpts, enabled, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "transferOwnership", to)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.TransferOwnership(&_UpkeepRegistrationRequests.TransactOpts, to)
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.TransferOwnership(&_UpkeepRegistrationRequests.TransactOpts, to)
}

type UpkeepRegistrationRequestsConfigChangedIterator struct {
	Event *UpkeepRegistrationRequestsConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepRegistrationRequestsConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsConfigChanged)
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
		it.Event = new(UpkeepRegistrationRequestsConfigChanged)
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

func (it *UpkeepRegistrationRequestsConfigChangedIterator) Error() error {
	return it.fail
}

func (it *UpkeepRegistrationRequestsConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepRegistrationRequestsConfigChanged struct {
	Enabled            bool
	WindowSizeInBlocks uint32
	AllowedPerWindow   uint16
	KeeperRegistry     common.Address
	MinLINKJuels       *big.Int
	Raw                types.Log
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*UpkeepRegistrationRequestsConfigChangedIterator, error) {

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsConfigChangedIterator{contract: _UpkeepRegistrationRequests.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsConfigChanged) (event.Subscription, error) {

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepRegistrationRequestsConfigChanged)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseConfigChanged(log types.Log) (*UpkeepRegistrationRequestsConfigChanged, error) {
	event := new(UpkeepRegistrationRequestsConfigChanged)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type UpkeepRegistrationRequestsOwnershipTransferRequestedIterator struct {
	Event *UpkeepRegistrationRequestsOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepRegistrationRequestsOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsOwnershipTransferRequested)
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
		it.Event = new(UpkeepRegistrationRequestsOwnershipTransferRequested)
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

func (it *UpkeepRegistrationRequestsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *UpkeepRegistrationRequestsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepRegistrationRequestsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpkeepRegistrationRequestsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsOwnershipTransferRequestedIterator{contract: _UpkeepRegistrationRequests.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepRegistrationRequestsOwnershipTransferRequested)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseOwnershipTransferRequested(log types.Log) (*UpkeepRegistrationRequestsOwnershipTransferRequested, error) {
	event := new(UpkeepRegistrationRequestsOwnershipTransferRequested)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type UpkeepRegistrationRequestsOwnershipTransferredIterator struct {
	Event *UpkeepRegistrationRequestsOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepRegistrationRequestsOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsOwnershipTransferred)
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
		it.Event = new(UpkeepRegistrationRequestsOwnershipTransferred)
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

func (it *UpkeepRegistrationRequestsOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *UpkeepRegistrationRequestsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepRegistrationRequestsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpkeepRegistrationRequestsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsOwnershipTransferredIterator{contract: _UpkeepRegistrationRequests.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepRegistrationRequestsOwnershipTransferred)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseOwnershipTransferred(log types.Log) (*UpkeepRegistrationRequestsOwnershipTransferred, error) {
	event := new(UpkeepRegistrationRequestsOwnershipTransferred)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type UpkeepRegistrationRequestsRegistrationApprovedIterator struct {
	Event *UpkeepRegistrationRequestsRegistrationApproved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepRegistrationRequestsRegistrationApprovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsRegistrationApproved)
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
		it.Event = new(UpkeepRegistrationRequestsRegistrationApproved)
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

func (it *UpkeepRegistrationRequestsRegistrationApprovedIterator) Error() error {
	return it.fail
}

func (it *UpkeepRegistrationRequestsRegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepRegistrationRequestsRegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*UpkeepRegistrationRequestsRegistrationApprovedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsRegistrationApprovedIterator{contract: _UpkeepRegistrationRequests.contract, event: "RegistrationApproved", logs: logs, sub: sub}, nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepRegistrationRequestsRegistrationApproved)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
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

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseRegistrationApproved(log types.Log) (*UpkeepRegistrationRequestsRegistrationApproved, error) {
	event := new(UpkeepRegistrationRequestsRegistrationApproved)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type UpkeepRegistrationRequestsRegistrationRejectedIterator struct {
	Event *UpkeepRegistrationRequestsRegistrationRejected

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepRegistrationRequestsRegistrationRejectedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsRegistrationRejected)
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
		it.Event = new(UpkeepRegistrationRequestsRegistrationRejected)
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

func (it *UpkeepRegistrationRequestsRegistrationRejectedIterator) Error() error {
	return it.fail
}

func (it *UpkeepRegistrationRequestsRegistrationRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepRegistrationRequestsRegistrationRejected struct {
	Hash [32]byte
	Raw  types.Log
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*UpkeepRegistrationRequestsRegistrationRejectedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsRegistrationRejectedIterator{contract: _UpkeepRegistrationRequests.contract, event: "RegistrationRejected", logs: logs, sub: sub}, nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationRejected, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepRegistrationRequestsRegistrationRejected)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
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

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseRegistrationRejected(log types.Log) (*UpkeepRegistrationRequestsRegistrationRejected, error) {
	event := new(UpkeepRegistrationRequestsRegistrationRejected)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type UpkeepRegistrationRequestsRegistrationRequestedIterator struct {
	Event *UpkeepRegistrationRequestsRegistrationRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepRegistrationRequestsRegistrationRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsRegistrationRequested)
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
		it.Event = new(UpkeepRegistrationRequestsRegistrationRequested)
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

func (it *UpkeepRegistrationRequestsRegistrationRequestedIterator) Error() error {
	return it.fail
}

func (it *UpkeepRegistrationRequestsRegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepRegistrationRequestsRegistrationRequested struct {
	Hash           [32]byte
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	Amount         *big.Int
	Source         uint8
	Raw            types.Log
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*UpkeepRegistrationRequestsRegistrationRequestedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsRegistrationRequestedIterator{contract: _UpkeepRegistrationRequests.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepRegistrationRequestsRegistrationRequested)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseRegistrationRequested(log types.Log) (*UpkeepRegistrationRequestsRegistrationRequested, error) {
	event := new(UpkeepRegistrationRequestsRegistrationRequested)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRegistrationConfig struct {
	Enabled                 bool
	WindowSizeInBlocks      uint32
	AllowedPerWindow        uint16
	KeeperRegistry          common.Address
	MinLINKJuels            *big.Int
	WindowStart             uint64
	ApprovedInCurrentWindow uint16
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequests) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepRegistrationRequests.abi.Events["ConfigChanged"].ID:
		return _UpkeepRegistrationRequests.ParseConfigChanged(log)
	case _UpkeepRegistrationRequests.abi.Events["OwnershipTransferRequested"].ID:
		return _UpkeepRegistrationRequests.ParseOwnershipTransferRequested(log)
	case _UpkeepRegistrationRequests.abi.Events["OwnershipTransferred"].ID:
		return _UpkeepRegistrationRequests.ParseOwnershipTransferred(log)
	case _UpkeepRegistrationRequests.abi.Events["RegistrationApproved"].ID:
		return _UpkeepRegistrationRequests.ParseRegistrationApproved(log)
	case _UpkeepRegistrationRequests.abi.Events["RegistrationRejected"].ID:
		return _UpkeepRegistrationRequests.ParseRegistrationRejected(log)
	case _UpkeepRegistrationRequests.abi.Events["RegistrationRequested"].ID:
		return _UpkeepRegistrationRequests.ParseRegistrationRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepRegistrationRequestsConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53")
}

func (UpkeepRegistrationRequestsOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (UpkeepRegistrationRequestsOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (UpkeepRegistrationRequestsRegistrationApproved) Topic() common.Hash {
	return common.HexToHash("0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b")
}

func (UpkeepRegistrationRequestsRegistrationRejected) Topic() common.Hash {
	return common.HexToHash("0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22")
}

func (UpkeepRegistrationRequestsRegistrationRequested) Topic() common.Hash {
	return common.HexToHash("0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78")
}

func (_UpkeepRegistrationRequests *UpkeepRegistrationRequests) Address() common.Address {
	return _UpkeepRegistrationRequests.address
}

type UpkeepRegistrationRequestsInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	GetPendingRequest(opts *bind.CallOpts, hash [32]byte) (common.Address, *big.Int, error)

	GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error)

	Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error)

	SetRegistrationConfig(opts *bind.TransactOpts, enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*UpkeepRegistrationRequestsConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*UpkeepRegistrationRequestsConfigChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpkeepRegistrationRequestsOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*UpkeepRegistrationRequestsOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpkeepRegistrationRequestsOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*UpkeepRegistrationRequestsOwnershipTransferred, error)

	FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*UpkeepRegistrationRequestsRegistrationApprovedIterator, error)

	WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error)

	ParseRegistrationApproved(log types.Log) (*UpkeepRegistrationRequestsRegistrationApproved, error)

	FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*UpkeepRegistrationRequestsRegistrationRejectedIterator, error)

	WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationRejected, hash [][32]byte) (event.Subscription, error)

	ParseRegistrationRejected(log types.Log) (*UpkeepRegistrationRequestsRegistrationRejected, error)

	FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*UpkeepRegistrationRequestsRegistrationRequestedIterator, error)

	WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error)

	ParseRegistrationRequested(log types.Log) (*UpkeepRegistrationRequestsRegistrationRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
