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

const OperatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CancelOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cancelExpiration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXPIRY_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunc\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"fulfillOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"fulfillOracleRequest2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"operatorTransferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"oracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAuthorizedSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

var OperatorBin = "0x60a0604052600160045534801561001557600080fd5b5060405161260c38038061260c8339818101604052604081101561003857600080fd5b508051602090910151600080546001600160a01b0319166001600160a01b03928316178155606083901b6001600160601b031916608052911690612565906100a79039806106a252806107685280611194528061127f528061192c5280611d6b52806121f452506125656000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806379ba509711610097578063f2fde38b11610066578063f2fde38b146105c6578063f3dfc2a9146105f9578063f3fef3a314610634578063fa00763a1461066d57610100565b806379ba50971461045c5780638da5cb5b14610464578063a4c0ed361461046c578063e93f30661461053457610100565b806350188301116100d357806350188301146102af5780636ae0bc76146102b75780636ee4d5531461037e5780636fadcf72146103cf57610100565b8063165d35e11461010557806340429946146101365780634ab0d1901461020e5780634b60228214610295575b600080fd5b61010d6106a0565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61020c600480360361010081101561014d57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235811692602081013592604082013592606083013516917fffffffff000000000000000000000000000000000000000000000000000000006080820135169160a08201359160c081013591810190610100810160e08201356401000000008111156101cd57600080fd5b8201836020820111156101df57600080fd5b8035906020019184600183028401116401000000008311171561020157600080fd5b5090925090506106c4565b005b610281600480360360c081101561022457600080fd5b5080359060208101359073ffffffffffffffffffffffffffffffffffffffff604082013516907fffffffff000000000000000000000000000000000000000000000000000000006060820135169060808101359060a00135610946565b604080519115158252519081900360200190f35b61029d610c45565b60408051918252519081900360200190f35b61029d610c4b565b610281600480360360c08110156102cd57600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516917fffffffff00000000000000000000000000000000000000000000000000000000606082013516916080820135919081019060c0810160a082013564010000000081111561033f57600080fd5b82018360208201111561035157600080fd5b8035906020019184600183028401116401000000008311171561037357600080fd5b509092509050610c5a565b61020c6004803603608081101561039457600080fd5b508035906020810135907fffffffff000000000000000000000000000000000000000000000000000000006040820135169060600135610fea565b61020c600480360360408110156103e557600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561041d57600080fd5b82018360208201111561042f57600080fd5b8035906020019184600183028401116401000000008311171561045157600080fd5b509092509050611215565b61020c611408565b61010d61150a565b61020c6004803603606081101561048257600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516916020810135918101906060810160408201356401000000008111156104bf57600080fd5b8201836020820111156104d157600080fd5b803590602001918460018302840111640100000000831117156104f357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611526945050505050565b6102816004803603606081101561054a57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561058757600080fd5b82018360208201111561059957600080fd5b803590602001918460018302840111640100000000831117156105bb57600080fd5b509092509050611842565b61020c600480360360208110156105dc57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611a25565b61020c6004803603604081101561060f57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001351515611ba6565b61020c6004803603604081101561064a57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135611c82565b6102816004803603602081101561068357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611e31565b7f000000000000000000000000000000000000000000000000000000000000000090565b6106cc6106a0565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461076557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b857f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561082157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f742063616c6c6261636b20746f204c494e4b000000000000000000604482015290519081900360640190fd5b6000806108328c8c8b8b8b8b611e5c565b91509150897fd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c658d848e8d8d878d8d8d604051808a73ffffffffffffffffffffffffffffffffffffffff1681526020018981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff168152602001867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001858152602001848152602001806020018281038252848482818152602001925080828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039c50909a5050505050505050505050a2505050505050505050505050565b3360009081526003602052604081205460ff166109ae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180612506602a913960400191505060405180910390fd5b600087815260026020526040902054879060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016610a4f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b610a5e88888888886001612045565b60405188907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a262061a805a1015610afb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60408051602481018a9052604480820186905282518083039091018152606490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000089161781529151815160009373ffffffffffffffffffffffffffffffffffffffff8b169392918291908083835b60208310610bce57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610b91565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610c30576040519150601f19603f3d011682016040523d82523d6000602084013e610c35565b606091505b50909a9950505050505050505050565b61012c81565b6000610c556121d3565b905090565b3360009081526003602052604081205460ff16610cc2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180612506602a913960400191505060405180910390fd5b600088815260026020526040902054889060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016610d6357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b8884848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050506020810151828114610e0d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f466972737420776f7264206d7573742062652072657175657374496400000000604482015290519081900360640190fd5b610e1c8c8c8c8c8c6002612045565b6040518c907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a262061a805a1015610eb957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60008a73ffffffffffffffffffffffffffffffffffffffff168a898960405160200180847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526004018383808284378083019250505093505050506040516020818303038152906040526040518082805190602001908083835b60208310610f6f57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610f32565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610fd1576040519150601f19603f3d011682016040523d82523d6000602084013e610fd6565b606091505b50909e9d5050505050505050505050505050565b6000610ff8843385856122b1565b60008681526002602052604090205490915060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00908116908216146110a057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b4282111561110f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f52657175657374206973206e6f74206578706972656400000000000000000000604482015290519081900360640190fd5b6000858152600260205260408082208290555186917fa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e9391a2604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815233600482015260248101869052905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163a9059cbb9160448083019260209291908290030181600087803b1580156111dc57600080fd5b505af11580156111f0573d6000803e3d6000fd5b505050506040513d602081101561120657600080fd5b505161120e57fe5b5050505050565b3360009081526003602052604090205460ff1661127d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180612506602a913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611322576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603281526020018061249f6032913960400191505060405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff168383604051808383808284376040519201945060009350909150508083038183865af19150503d806000811461138d576040519150601f19603f3d011682016040523d82523d6000602084013e611392565b606091505b505090508061140257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f466f727761726465642063616c6c206661696c65642e00000000000000000000604482015290519081900360640190fd5b50505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461148e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b61152e6106a0565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146115c757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b805181906044111561163a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c69642072657175657374206c656e67746800000000000000000000604482015290519081900360640190fd5b602082015182907fffffffff0000000000000000000000000000000000000000000000000000000081167f4042994600000000000000000000000000000000000000000000000000000000146116f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b85602485015284604485015260003073ffffffffffffffffffffffffffffffffffffffff16856040518082805190602001908083835b6020831061176457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101611727565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146117c4576040519150601f19603f3d011682016040523d82523d6000602084013e6117c9565b606091505b505090508061183957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f556e61626c6520746f2063726561746520726571756573740000000000000000604482015290519081900360640190fd5b50505050505050565b6000805473ffffffffffffffffffffffffffffffffffffffff1633146118c957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b83806118d36121d3565b101561192a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806124d16035913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea0878787876040518563ffffffff1660e01b8152600401808573ffffffffffffffffffffffffffffffffffffffff168152602001848152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505095505050505050602060405180830381600087803b1580156119ef57600080fd5b505af1158015611a03573d6000803e3d6000fd5b505050506040513d6020811015611a1957600080fd5b50519695505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611aab57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8116331415611b3057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314611c2c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff91909116600090815260036020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff163314611d0857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b8080611d126121d3565b1015611d69576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806124d16035913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015611dfa57600080fd5b505af1158015611e0e573d6000803e3d6000fd5b505050506040513d6020811015611e2457600080fd5b5051611e2c57fe5b505050565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205460ff1690565b60408051606088901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660208083019190915260348083018690528351808403909101815260549092018352815191810191909120600081815260029092529181205460081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001615611f5157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4d75737420757365206120756e69717565204944000000000000000000000000604482015290519081900360640190fd5b611f5d4261012c612336565b90506000611f6d888888856122b1565b905060405180604001604052808260ff19168152602001611f8d866123b1565b60ff90811690915260008581526002602090815260409091208351815494909201519092167f01000000000000000000000000000000000000000000000000000000000000000260089190911c7fff00000000000000000000000000000000000000000000000000000000000000909316929092177effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff169190911790556004546120369089612336565b60045550965096945050505050565b6000612053868686866122b1565b60008881526002602052604090205490915060081b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00908116908216146120fb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b612104826123b1565b60008881526002602052604090205460ff9182167f010000000000000000000000000000000000000000000000000000000000000090910490911611156121ac57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f446174612076657273696f6e73206d757374206d617463680000000000000000604482015290519081900360640190fd5b6004546121b99087612427565b600455505050600093845250506002602052506040812055565b6000806121ec600160045461242790919063ffffffff16565b90506122ab817f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561227957600080fd5b505afa15801561228d573d6000803e3d6000fd5b505050506040513d60208110156122a357600080fd5b505190612427565b91505090565b6040805160208082019690965260609490941b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016848201527fffffffff000000000000000000000000000000000000000000000000000000009290921660548401526058808401919091528151808403909101815260789092019052805191012090565b6000828201838110156123aa57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6000610100821061242357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f6e756d62657220746f6f2062696720746f206361737400000000000000000000604482015290519081900360640190fd5b5090565b60008282111561249857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b5090039056fe43616e6e6f74207573652023666f727761726420746f2073656e64206d6573736167657320746f204c696e6b20746f6b656e416d6f756e74207265717565737465642069732067726561746572207468616e20776974686472617761626c652062616c616e63654e6f7420616e20617574686f72697a6564206e6f646520746f2066756c66696c6c207265717565737473a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

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
