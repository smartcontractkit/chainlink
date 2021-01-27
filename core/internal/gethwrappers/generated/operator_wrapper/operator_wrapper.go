// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package operator_wrapper

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

const OperatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CancelOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cancelExpiration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXPIRY_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunc\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable[]\",\"name\":\"receivers\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"distributeFunds\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"fulfillOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"fulfillOracleRequest2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"operatorTransferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"oracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAuthorizedSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

var OperatorBin = "0x60a0604052600160045534801561001557600080fd5b506040516129433803806129438339818101604052604081101561003857600080fd5b508051602090910151600080546001600160a01b0319166001600160a01b03928316178155606083901b6001600160601b03191660805291169061289c906100a7903980610842528061090852806114d252806115b65280611c6352806120a2528061252b525061289c6000f3fe60806040526004361061010e5760003560e01c80636fadcf72116100a5578063e93f306611610074578063f3dfc2a911610059578063f3dfc2a914610772578063f3fef3a3146107ba578063fa00763a146108005761010e565b8063e93f306614610693578063f2fde38b146107325761010e565b80636fadcf72146104fa57806379ba5097146105945780638da5cb5b146105a9578063a4c0ed36146105be5761010e565b806350188301116100e157806350188301146102f15780636ae0bc76146103065780636bd59ec0146103da5780636ee4d5531461049c5761010e565b8063165d35e11461011357806340429946146101515780634ab0d190146102365780634b602282146102ca575b600080fd5b34801561011f57600080fd5b50610128610840565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b34801561015d57600080fd5b50610234600480360361010081101561017557600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235811692602081013592604082013592606083013516917fffffffff000000000000000000000000000000000000000000000000000000006080820135169160a08201359160c081013591810190610100810160e08201356401000000008111156101f557600080fd5b82018360208201111561020757600080fd5b8035906020019184600183028401116401000000008311171561022957600080fd5b509092509050610864565b005b34801561024257600080fd5b506102b6600480360360c081101561025957600080fd5b5080359060208101359073ffffffffffffffffffffffffffffffffffffffff604082013516907fffffffff000000000000000000000000000000000000000000000000000000006060820135169060808101359060a00135610ae6565b604080519115158252519081900360200190f35b3480156102d657600080fd5b506102df610de5565b60408051918252519081900360200190f35b3480156102fd57600080fd5b506102df610deb565b34801561031257600080fd5b506102b6600480360360c081101561032957600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516917fffffffff00000000000000000000000000000000000000000000000000000000606082013516916080820135919081019060c0810160a082013564010000000081111561039b57600080fd5b8201836020820111156103ad57600080fd5b803590602001918460018302840111640100000000831117156103cf57600080fd5b509092509050610dfa565b610234600480360360408110156103f057600080fd5b81019060208101813564010000000081111561040b57600080fd5b82018360208201111561041d57600080fd5b8035906020019184602083028401116401000000008311171561043f57600080fd5b91939092909160208101903564010000000081111561045d57600080fd5b82018360208201111561046f57600080fd5b8035906020019184602083028401116401000000008311171561049157600080fd5b50909250905061118a565b3480156104a857600080fd5b50610234600480360360808110156104bf57600080fd5b508035906020810135907fffffffff000000000000000000000000000000000000000000000000000000006040820135169060600135611328565b34801561050657600080fd5b506102346004803603604081101561051d57600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561055557600080fd5b82018360208201111561056757600080fd5b8035906020019184600183028401116401000000008311171561058957600080fd5b50909250905061154c565b3480156105a057600080fd5b5061023461173f565b3480156105b557600080fd5b50610128611841565b3480156105ca57600080fd5b50610234600480360360608110156105e157600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561061e57600080fd5b82018360208201111561063057600080fd5b8035906020019184600183028401116401000000008311171561065257600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061185d945050505050565b34801561069f57600080fd5b506102b6600480360360608110156106b657600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516916020810135918101906060810160408201356401000000008111156106f357600080fd5b82018360208201111561070557600080fd5b8035906020019184600183028401116401000000008311171561072757600080fd5b509092509050611b79565b34801561073e57600080fd5b506102346004803603602081101561075557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611d5c565b34801561077e57600080fd5b506102346004803603604081101561079557600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001351515611edd565b3480156107c657600080fd5b50610234600480360360408110156107dd57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135611fb9565b34801561080c57600080fd5b506102b66004803603602081101561082357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16612168565b7f000000000000000000000000000000000000000000000000000000000000000090565b61086c610840565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461090557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b857f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156109c157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f742063616c6c6261636b20746f204c494e4b000000000000000000604482015290519081900360640190fd5b6000806109d28c8c8b8b8b8b612193565b91509150897fd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c658d848e8d8d878d8d8d604051808a73ffffffffffffffffffffffffffffffffffffffff1681526020018981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff168152602001867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001858152602001848152602001806020018281038252848482818152602001925080828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039c50909a5050505050505050505050a2505050505050505050505050565b3360009081526003602052604081205460ff16610b4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a81526020018061283d602a913960400191505060405180910390fd5b600087815260026020526040902054879060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016610bef57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b610bfe8888888888600161237c565b60405188907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a262061a805a1015610c9b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60408051602481018a9052604480820186905282518083039091018152606490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000089161781529151815160009373ffffffffffffffffffffffffffffffffffffffff8b169392918291908083835b60208310610d6e57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610d31565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610dd0576040519150601f19603f3d011682016040523d82523d6000602084013e610dd5565b606091505b50909a9950505050505050505050565b61012c81565b6000610df561250a565b905090565b3360009081526003602052604081205460ff16610e62576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a81526020018061283d602a913960400191505060405180910390fd5b600088815260026020526040902054889060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016610f0357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b8884848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050506020810151828114610fad57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f466972737420776f7264206d7573742062652072657175657374496400000000604482015290519081900360640190fd5b610fbc8c8c8c8c8c600261237c565b6040518c907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a262061a805a101561105957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60008a73ffffffffffffffffffffffffffffffffffffffff168a898960405160200180847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526004018383808284378083019250505093505050506040516020818303038152906040526040518082805190602001908083835b6020831061110f57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016110d2565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611171576040519150601f19603f3d011682016040523d82523d6000602084013e611176565b606091505b50909e9d5050505050505050505050505050565b821580159061119857508281145b61120357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f496e76616c6964206172726179206c656e677468287329000000000000000000604482015290519081900360640190fd5b3460005b848110156112b357600084848381811061121d57fe5b90506020020135905061123981846125e890919063ffffffff16565b925086868381811061124757fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f193505050501580156112a9573d6000803e3d6000fd5b5050600101611207565b50801561132157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f546f6f206d756368204554482073656e74000000000000000000000000000000604482015290519081900360640190fd5b5050505050565b60006113368433858561265f565b60008681526002602052604090205490915060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00908116908216146113de57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b4282111561144d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f52657175657374206973206e6f74206578706972656400000000000000000000604482015290519081900360640190fd5b6000858152600260205260408082208290555186917fa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e9391a2604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815233600482015260248101869052905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163a9059cbb9160448083019260209291908290030181600087803b15801561151a57600080fd5b505af115801561152e573d6000803e3d6000fd5b505050506040513d602081101561154457600080fd5b505161132157fe5b3360009081526003602052604090205460ff166115b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a81526020018061283d602a913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611659576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260328152602001806127d66032913960400191505060405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff168383604051808383808284376040519201945060009350909150508083038183865af19150503d80600081146116c4576040519150601f19603f3d011682016040523d82523d6000602084013e6116c9565b606091505b505090508061173957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f466f727761726465642063616c6c206661696c65642e00000000000000000000604482015290519081900360640190fd5b50505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146117c557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b611865610840565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146118fe57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b805181906044111561197157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c69642072657175657374206c656e67746800000000000000000000604482015290519081900360640190fd5b602082015182907fffffffff0000000000000000000000000000000000000000000000000000000081167f404299460000000000000000000000000000000000000000000000000000000014611a2857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b85602485015284604485015260003073ffffffffffffffffffffffffffffffffffffffff16856040518082805190602001908083835b60208310611a9b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101611a5e565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d8060008114611afb576040519150601f19603f3d011682016040523d82523d6000602084013e611b00565b606091505b5050905080611b7057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f556e61626c6520746f2063726561746520726571756573740000000000000000604482015290519081900360640190fd5b50505050505050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314611c0057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b8380611c0a61250a565b1015611c61576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806128086035913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea0878787876040518563ffffffff1660e01b8152600401808573ffffffffffffffffffffffffffffffffffffffff168152602001848152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505095505050505050602060405180830381600087803b158015611d2657600080fd5b505af1158015611d3a573d6000803e3d6000fd5b505050506040513d6020811015611d5057600080fd5b50519695505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611de257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8116331415611e6757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314611f6357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff91909116600090815260036020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff16331461203f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b808061204961250a565b10156120a0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806128086035913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561213157600080fd5b505af1158015612145573d6000803e3d6000fd5b505050506040513d602081101561215b57600080fd5b505161216357fe5b505050565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205460ff1690565b60408051606088901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660208083019190915260348083018690528351808403909101815260549092018352815191810191909120600081815260029092529181205460081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00161561228857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4d75737420757365206120756e69717565204944000000000000000000000000604482015290519081900360640190fd5b6122944261012c6126e4565b905060006122a48888888561265f565b905060405180604001604052808260ff191681526020016122c48661275f565b60ff90811690915260008581526002602090815260409091208351815494909201519092167f01000000000000000000000000000000000000000000000000000000000000000260089190911c7fff00000000000000000000000000000000000000000000000000000000000000909316929092177effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1691909117905560045461236d90896126e4565b60045550965096945050505050565b600061238a8686868661265f565b60008881526002602052604090205490915060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009081169082161461243257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b61243b8261275f565b60008881526002602052604090205460ff9182167f010000000000000000000000000000000000000000000000000000000000000090910490911611156124e357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f446174612076657273696f6e73206d757374206d617463680000000000000000604482015290519081900360640190fd5b6004546124f090876125e8565b600455505050600093845250506002602052506040812055565b60008061252360016004546125e890919063ffffffff16565b90506125e2817f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156125b057600080fd5b505afa1580156125c4573d6000803e3d6000fd5b505050506040513d60208110156125da57600080fd5b5051906125e8565b91505090565b60008282111561265957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6040805160208082019690965260609490941b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016848201527fffffffff000000000000000000000000000000000000000000000000000000009290921660548401526058808401919091528151808403909101815260789092019052805191012090565b60008282018381101561275857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b600061010082106127d157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f6e756d62657220746f6f2062696720746f206361737400000000000000000000604482015290519081900360640190fd5b509056fe43616e6e6f74207573652023666f727761726420746f2073656e64206d6573736167657320746f204c696e6b20746f6b656e416d6f756e74207265717565737465642069732067726561746572207468616e20776974686472617761626c652062616c616e63654e6f7420616e20617574686f72697a6564206e6f646520746f2066756c66696c6c207265717565737473a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

func DeployOperator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, owner common.Address) (common.Address, *types.Transaction, *Operator, error) {
	parsed, err := abi.JSON(strings.NewReader(OperatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OperatorBin), backend, link, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Operator{OperatorCaller: OperatorCaller{contract: contract}, OperatorTransactor: OperatorTransactor{contract: contract}, OperatorFilterer: OperatorFilterer{contract: contract}}, nil
}

type Operator struct {
	address common.Address
	OperatorCaller
	OperatorTransactor
	OperatorFilterer
}

type OperatorCaller struct {
	contract *bind.BoundContract
}

type OperatorTransactor struct {
	contract *bind.BoundContract
}

type OperatorFilterer struct {
	contract *bind.BoundContract
}

type OperatorSession struct {
	Contract     *Operator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OperatorCallerSession struct {
	Contract *OperatorCaller
	CallOpts bind.CallOpts
}

type OperatorTransactorSession struct {
	Contract     *OperatorTransactor
	TransactOpts bind.TransactOpts
}

type OperatorRaw struct {
	Contract *Operator
}

type OperatorCallerRaw struct {
	Contract *OperatorCaller
}

type OperatorTransactorRaw struct {
	Contract *OperatorTransactor
}

func NewOperator(address common.Address, backend bind.ContractBackend) (*Operator, error) {
	contract, err := bindOperator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Operator{address: address, OperatorCaller: OperatorCaller{contract: contract}, OperatorTransactor: OperatorTransactor{contract: contract}, OperatorFilterer: OperatorFilterer{contract: contract}}, nil
}

func NewOperatorCaller(address common.Address, caller bind.ContractCaller) (*OperatorCaller, error) {
	contract, err := bindOperator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OperatorCaller{contract: contract}, nil
}

func NewOperatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OperatorTransactor, error) {
	contract, err := bindOperator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OperatorTransactor{contract: contract}, nil
}

func NewOperatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OperatorFilterer, error) {
	contract, err := bindOperator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OperatorFilterer{contract: contract}, nil
}

func bindOperator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OperatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Operator *OperatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Operator.Contract.OperatorCaller.contract.Call(opts, result, method, params...)
}

func (_Operator *OperatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operator.Contract.OperatorTransactor.contract.Transfer(opts)
}

func (_Operator *OperatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Operator.Contract.OperatorTransactor.contract.Transact(opts, method, params...)
}

func (_Operator *OperatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Operator.Contract.contract.Call(opts, result, method, params...)
}

func (_Operator *OperatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operator.Contract.contract.Transfer(opts)
}

func (_Operator *OperatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Operator.Contract.contract.Transact(opts, method, params...)
}

func (_Operator *OperatorCaller) EXPIRYTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Operator.contract.Call(opts, &out, "EXPIRY_TIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Operator *OperatorSession) EXPIRYTIME() (*big.Int, error) {
	return _Operator.Contract.EXPIRYTIME(&_Operator.CallOpts)
}

func (_Operator *OperatorCallerSession) EXPIRYTIME() (*big.Int, error) {
	return _Operator.Contract.EXPIRYTIME(&_Operator.CallOpts)
}

func (_Operator *OperatorCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Operator.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Operator *OperatorSession) GetChainlinkToken() (common.Address, error) {
	return _Operator.Contract.GetChainlinkToken(&_Operator.CallOpts)
}

func (_Operator *OperatorCallerSession) GetChainlinkToken() (common.Address, error) {
	return _Operator.Contract.GetChainlinkToken(&_Operator.CallOpts)
}

func (_Operator *OperatorCaller) IsAuthorizedSender(opts *bind.CallOpts, node common.Address) (bool, error) {
	var out []interface{}
	err := _Operator.contract.Call(opts, &out, "isAuthorizedSender", node)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Operator *OperatorSession) IsAuthorizedSender(node common.Address) (bool, error) {
	return _Operator.Contract.IsAuthorizedSender(&_Operator.CallOpts, node)
}

func (_Operator *OperatorCallerSession) IsAuthorizedSender(node common.Address) (bool, error) {
	return _Operator.Contract.IsAuthorizedSender(&_Operator.CallOpts, node)
}

func (_Operator *OperatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Operator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Operator *OperatorSession) Owner() (common.Address, error) {
	return _Operator.Contract.Owner(&_Operator.CallOpts)
}

func (_Operator *OperatorCallerSession) Owner() (common.Address, error) {
	return _Operator.Contract.Owner(&_Operator.CallOpts)
}

func (_Operator *OperatorCaller) Withdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Operator.contract.Call(opts, &out, "withdrawable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Operator *OperatorSession) Withdrawable() (*big.Int, error) {
	return _Operator.Contract.Withdrawable(&_Operator.CallOpts)
}

func (_Operator *OperatorCallerSession) Withdrawable() (*big.Int, error) {
	return _Operator.Contract.Withdrawable(&_Operator.CallOpts)
}

func (_Operator *OperatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "acceptOwnership")
}

func (_Operator *OperatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Operator.Contract.AcceptOwnership(&_Operator.TransactOpts)
}

func (_Operator *OperatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Operator.Contract.AcceptOwnership(&_Operator.TransactOpts)
}

func (_Operator *OperatorTransactor) CancelOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "cancelOracleRequest", requestId, payment, callbackFunc, expiration)
}

func (_Operator *OperatorSession) CancelOracleRequest(requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.CancelOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackFunc, expiration)
}

func (_Operator *OperatorTransactorSession) CancelOracleRequest(requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.CancelOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackFunc, expiration)
}

func (_Operator *OperatorTransactor) DistributeFunds(opts *bind.TransactOpts, receivers []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "distributeFunds", receivers, amounts)
}

func (_Operator *OperatorSession) DistributeFunds(receivers []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Operator.Contract.DistributeFunds(&_Operator.TransactOpts, receivers, amounts)
}

func (_Operator *OperatorTransactorSession) DistributeFunds(receivers []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Operator.Contract.DistributeFunds(&_Operator.TransactOpts, receivers, amounts)
}

func (_Operator *OperatorTransactor) Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "forward", to, data)
}

func (_Operator *OperatorSession) Forward(to common.Address, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.Forward(&_Operator.TransactOpts, to, data)
}

func (_Operator *OperatorTransactorSession) Forward(to common.Address, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.Forward(&_Operator.TransactOpts, to, data)
}

func (_Operator *OperatorTransactor) FulfillOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "fulfillOracleRequest", requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

func (_Operator *OperatorSession) FulfillOracleRequest(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

func (_Operator *OperatorTransactorSession) FulfillOracleRequest(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

func (_Operator *OperatorTransactor) FulfillOracleRequest2(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "fulfillOracleRequest2", requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

func (_Operator *OperatorSession) FulfillOracleRequest2(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest2(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

func (_Operator *OperatorTransactorSession) FulfillOracleRequest2(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest2(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

func (_Operator *OperatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_Operator *OperatorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OnTokenTransfer(&_Operator.TransactOpts, sender, amount, data)
}

func (_Operator *OperatorTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OnTokenTransfer(&_Operator.TransactOpts, sender, amount, data)
}

func (_Operator *OperatorTransactor) OperatorTransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "operatorTransferAndCall", to, value, data)
}

func (_Operator *OperatorSession) OperatorTransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OperatorTransferAndCall(&_Operator.TransactOpts, to, value, data)
}

func (_Operator *OperatorTransactorSession) OperatorTransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OperatorTransferAndCall(&_Operator.TransactOpts, to, value, data)
}

func (_Operator *OperatorTransactor) OracleRequest(opts *bind.TransactOpts, sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "oracleRequest", sender, payment, specId, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

func (_Operator *OperatorSession) OracleRequest(sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OracleRequest(&_Operator.TransactOpts, sender, payment, specId, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

func (_Operator *OperatorTransactorSession) OracleRequest(sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OracleRequest(&_Operator.TransactOpts, sender, payment, specId, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

func (_Operator *OperatorTransactor) SetAuthorizedSender(opts *bind.TransactOpts, node common.Address, allowed bool) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "setAuthorizedSender", node, allowed)
}

func (_Operator *OperatorSession) SetAuthorizedSender(node common.Address, allowed bool) (*types.Transaction, error) {
	return _Operator.Contract.SetAuthorizedSender(&_Operator.TransactOpts, node, allowed)
}

func (_Operator *OperatorTransactorSession) SetAuthorizedSender(node common.Address, allowed bool) (*types.Transaction, error) {
	return _Operator.Contract.SetAuthorizedSender(&_Operator.TransactOpts, node, allowed)
}

func (_Operator *OperatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "transferOwnership", to)
}

func (_Operator *OperatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Operator.Contract.TransferOwnership(&_Operator.TransactOpts, to)
}

func (_Operator *OperatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Operator.Contract.TransferOwnership(&_Operator.TransactOpts, to)
}

func (_Operator *OperatorTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "withdraw", recipient, amount)
}

func (_Operator *OperatorSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.Withdraw(&_Operator.TransactOpts, recipient, amount)
}

func (_Operator *OperatorTransactorSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.Withdraw(&_Operator.TransactOpts, recipient, amount)
}

type OperatorCancelOracleRequestIterator struct {
	Event *OperatorCancelOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OperatorCancelOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OperatorCancelOracleRequest)
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
		it.Event = new(OperatorCancelOracleRequest)
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

func (it *OperatorCancelOracleRequestIterator) Error() error {
	return it.fail
}

func (it *OperatorCancelOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OperatorCancelOracleRequest struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_Operator *OperatorFilterer) FilterCancelOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*OperatorCancelOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Operator.contract.FilterLogs(opts, "CancelOracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OperatorCancelOracleRequestIterator{contract: _Operator.contract, event: "CancelOracleRequest", logs: logs, sub: sub}, nil
}

func (_Operator *OperatorFilterer) WatchCancelOracleRequest(opts *bind.WatchOpts, sink chan<- *OperatorCancelOracleRequest, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Operator.contract.WatchLogs(opts, "CancelOracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OperatorCancelOracleRequest)
				if err := _Operator.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
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

func (_Operator *OperatorFilterer) ParseCancelOracleRequest(log types.Log) (*OperatorCancelOracleRequest, error) {
	event := new(OperatorCancelOracleRequest)
	if err := _Operator.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OperatorOracleRequestIterator struct {
	Event *OperatorOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OperatorOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OperatorOracleRequest)
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
		it.Event = new(OperatorOracleRequest)
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

func (it *OperatorOracleRequestIterator) Error() error {
	return it.fail
}

func (it *OperatorOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OperatorOracleRequest struct {
	SpecId             [32]byte
	Requester          common.Address
	RequestId          [32]byte
	Payment            *big.Int
	CallbackAddr       common.Address
	CallbackFunctionId [4]byte
	CancelExpiration   *big.Int
	DataVersion        *big.Int
	Data               []byte
	Raw                types.Log
}

func (_Operator *OperatorFilterer) FilterOracleRequest(opts *bind.FilterOpts, specId [][32]byte) (*OperatorOracleRequestIterator, error) {

	var specIdRule []interface{}
	for _, specIdItem := range specId {
		specIdRule = append(specIdRule, specIdItem)
	}

	logs, sub, err := _Operator.contract.FilterLogs(opts, "OracleRequest", specIdRule)
	if err != nil {
		return nil, err
	}
	return &OperatorOracleRequestIterator{contract: _Operator.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_Operator *OperatorFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OperatorOracleRequest, specId [][32]byte) (event.Subscription, error) {

	var specIdRule []interface{}
	for _, specIdItem := range specId {
		specIdRule = append(specIdRule, specIdItem)
	}

	logs, sub, err := _Operator.contract.WatchLogs(opts, "OracleRequest", specIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OperatorOracleRequest)
				if err := _Operator.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_Operator *OperatorFilterer) ParseOracleRequest(log types.Log) (*OperatorOracleRequest, error) {
	event := new(OperatorOracleRequest)
	if err := _Operator.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OperatorOracleResponseIterator struct {
	Event *OperatorOracleResponse

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OperatorOracleResponseIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OperatorOracleResponse)
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
		it.Event = new(OperatorOracleResponse)
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

func (it *OperatorOracleResponseIterator) Error() error {
	return it.fail
}

func (it *OperatorOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OperatorOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_Operator *OperatorFilterer) FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*OperatorOracleResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Operator.contract.FilterLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OperatorOracleResponseIterator{contract: _Operator.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

func (_Operator *OperatorFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *OperatorOracleResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Operator.contract.WatchLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OperatorOracleResponse)
				if err := _Operator.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

func (_Operator *OperatorFilterer) ParseOracleResponse(log types.Log) (*OperatorOracleResponse, error) {
	event := new(OperatorOracleResponse)
	if err := _Operator.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OperatorOwnershipTransferRequestedIterator struct {
	Event *OperatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OperatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OperatorOwnershipTransferRequested)
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
		it.Event = new(OperatorOwnershipTransferRequested)
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

func (it *OperatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OperatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OperatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Operator *OperatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OperatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Operator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OperatorOwnershipTransferRequestedIterator{contract: _Operator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_Operator *OperatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OperatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Operator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OperatorOwnershipTransferRequested)
				if err := _Operator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_Operator *OperatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OperatorOwnershipTransferRequested, error) {
	event := new(OperatorOwnershipTransferRequested)
	if err := _Operator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OperatorOwnershipTransferredIterator struct {
	Event *OperatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OperatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OperatorOwnershipTransferred)
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
		it.Event = new(OperatorOwnershipTransferred)
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

func (it *OperatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OperatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OperatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Operator *OperatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OperatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Operator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OperatorOwnershipTransferredIterator{contract: _Operator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Operator *OperatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OperatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Operator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OperatorOwnershipTransferred)
				if err := _Operator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Operator *OperatorFilterer) ParseOwnershipTransferred(log types.Log) (*OperatorOwnershipTransferred, error) {
	event := new(OperatorOwnershipTransferred)
	if err := _Operator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Operator *Operator) UnpackLog(out interface{}, event string, log types.Log) error {
	return _Operator.OperatorFilterer.contract.UnpackLog(out, event, log)
}

func (_Operator *Operator) Address() common.Address {
	return _Operator.address
}

type OperatorInterface interface {
	EXPIRYTIME(opts *bind.CallOpts) (*big.Int, error)

	GetChainlinkToken(opts *bind.CallOpts) (common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, node common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Withdrawable(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CancelOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error)

	DistributeFunds(opts *bind.TransactOpts, receivers []common.Address, amounts []*big.Int) (*types.Transaction, error)

	Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error)

	FulfillOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error)

	FulfillOracleRequest2(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OperatorTransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error)

	OracleRequest(opts *bind.TransactOpts, sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error)

	SetAuthorizedSender(opts *bind.TransactOpts, node common.Address, allowed bool) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	FilterCancelOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*OperatorCancelOracleRequestIterator, error)

	WatchCancelOracleRequest(opts *bind.WatchOpts, sink chan<- *OperatorCancelOracleRequest, requestId [][32]byte) (event.Subscription, error)

	ParseCancelOracleRequest(log types.Log) (*OperatorCancelOracleRequest, error)

	FilterOracleRequest(opts *bind.FilterOpts, specId [][32]byte) (*OperatorOracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OperatorOracleRequest, specId [][32]byte) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*OperatorOracleRequest, error)

	FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*OperatorOracleResponseIterator, error)

	WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *OperatorOracleResponse, requestId [][32]byte) (event.Subscription, error)

	ParseOracleResponse(log types.Log) (*OperatorOracleResponse, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OperatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OperatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OperatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OperatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OperatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OperatorOwnershipTransferred, error)

	UnpackLog(out interface{}, event string, log types.Log) error

	Address() common.Address
}
