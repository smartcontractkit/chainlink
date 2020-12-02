// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package operator

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OperatorABI is the input ABI used to generate the binding from.
const OperatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CancelOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cancelExpiration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXPIRY_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunc\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"fulfillOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"fulfillOracleRequest2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"getAuthorizationStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"oracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setFulfillmentPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OperatorBin is the compiled bytecode used for deploying new contracts.
var OperatorBin = "0x60a0604052600160055534801561001557600080fd5b506040516121ae3803806121ae8339818101604052602081101561003857600080fd5b5051600080546001600160a01b031916331790556001600160601b031960609190911b1660805260805160601c61211e6100906000398061063b5280610701528061113552806112205280611b6c525061211e6000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80636fadcf7211610097578063a4c0ed3611610066578063a4c0ed36146104d2578063d3e9c3141461059a578063f2fde38b146105cd578063f3fef3a314610600576100f5565b80636fadcf72146103fa57806379ba5097146104875780637fcd56db1461048f5780638da5cb5b146104ca576100f5565b80634b602282116100d35780634b6022821461028a57806350188301146102a45780636ae0bc76146102ac5780636ee4d553146103a9576100f5565b8063165d35e1146100fa578063404299461461012b5780634ab0d19014610203575b600080fd5b610102610639565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610201600480360361010081101561014257600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235811692602081013592604082013592606083013516917fffffffff000000000000000000000000000000000000000000000000000000006080820135169160a08201359160c081013591810190610100810160e08201356401000000008111156101c257600080fd5b8201836020820111156101d457600080fd5b803590602001918460018302840111640100000000831117156101f657600080fd5b50909250905061065d565b005b610276600480360360c081101561021957600080fd5b5080359060208101359073ffffffffffffffffffffffffffffffffffffffff604082013516907fffffffff000000000000000000000000000000000000000000000000000000006060820135169060808101359060a001356108df565b604080519115158252519081900360200190f35b610292610bb9565b60408051918252519081900360200190f35b610292610bbf565b610276600480360360c08110156102c257600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516917fffffffff00000000000000000000000000000000000000000000000000000000606082013516916080820135919081019060c0810160a082013564010000000081111561033457600080fd5b82018360208201111561034657600080fd5b8035906020019184600183028401116401000000008311171561036857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610bd5945050505050565b610201600480360360808110156103bf57600080fd5b508035906020810135907fffffffff000000000000000000000000000000000000000000000000000000006040820135169060600135610f68565b6102016004803603604081101561041057600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561044857600080fd5b82018360208201111561045a57600080fd5b8035906020019184600183028401116401000000008311171561047c57600080fd5b5090925090506111b6565b6102016113a9565b610201600480360360408110156104a557600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020013515156114ab565b610102611587565b610201600480360360608110156104e857600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561052557600080fd5b82018360208201111561053757600080fd5b8035906020019184600183028401116401000000008311171561055957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506115a3945050505050565b610276600480360360208110156105b057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118bf565b610201600480360360208110156105e357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118ea565b6102016004803603604081101561061657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135611a6b565b7f000000000000000000000000000000000000000000000000000000000000000090565b610665610639565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106fe57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b857f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156107ba57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f742063616c6c6261636b20746f204c494e4b000000000000000000604482015290519081900360640190fd5b6000806107cb8c8c8b8b8b8b611c32565b91509150897fd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c658d848e8d8d878d8d8d604051808a73ffffffffffffffffffffffffffffffffffffffff1681526020018981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff168152602001867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001858152602001848152602001806020018281038252848482818152602001925080828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039c50909a5050505050505050505050a2505050505050505050505050565b3360009081526003602052604081205460ff16610947576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001806120bf602a913960400191505060405180910390fd5b60008781526002602052604090205487906109c357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b6109d288888888886001611db5565b60405188907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a262061a805a1015610a6f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60408051602481018a9052604480820186905282518083039091018152606490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000089161781529151815160009373ffffffffffffffffffffffffffffffffffffffff8b169392918291908083835b60208310610b4257805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610b05565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610ba4576040519150601f19603f3d011682016040523d82523d6000602084013e610ba9565b606091505b50909a9950505050505050505050565b61012c81565b600554600090610bd0906001611f65565b905090565b3360009081526003602052604081205460ff16610c3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001806120bf602a913960400191505060405180910390fd5b6000878152600260205260409020548790610cb957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b602083015188908490808314610d3057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f466972737420776f7264206d7573742062652072657175657374496400000000604482015290519081900360640190fd5b610d3f8b8b8b8b8b6002611db5565b6040518b907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a262061a805a1015610ddc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60008973ffffffffffffffffffffffffffffffffffffffff16898860405160200180837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260040182805190602001908083835b60208310610e6c57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610e2f565b6001836020036101000a038019825116818451168082178552505050505050905001925050506040516020818303038152906040526040518082805190602001908083835b60208310610eee57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610eb1565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610f50576040519150601f19603f3d011682016040523d82523d6000602084013e610f55565b606091505b50909d9c50505050505050505050505050565b6040805160208082018690523360601b828401527fffffffff0000000000000000000000000000000000000000000000000000000085166054830152605880830185905283518084039091018152607890920183528151918101919091206000878152600290925291902054811461104157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b428211156110b057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f52657175657374206973206e6f74206578706972656400000000000000000000604482015290519081900360640190fd5b6000858152600260205260408082208290555186917fa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e9391a2604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815233600482015260248101869052905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163a9059cbb9160448083019260209291908290030181600087803b15801561117d57600080fd5b505af1158015611191573d6000803e3d6000fd5b505050506040513d60208110156111a757600080fd5b50516111af57fe5b5050505050565b3360009081526003602052604090205460ff1661121e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001806120bf602a913960400191505060405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156112c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260328152602001806120586032913960400191505060405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff168383604051808383808284376040519201945060009350909150508083038183865af19150503d806000811461132e576040519150601f19603f3d011682016040523d82523d6000602084013e611333565b606091505b50509050806113a357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f466f727761726465642063616c6c206661696c65642e00000000000000000000604482015290519081900360640190fd5b50505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461142f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461153157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff91909116600090815260036020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b6115ab610639565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461164457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b80518190604411156116b757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c69642072657175657374206c656e67746800000000000000000000604482015290519081900360640190fd5b602082015182907fffffffff0000000000000000000000000000000000000000000000000000000081167f40429946000000000000000000000000000000000000000000000000000000001461176e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b85602485015284604485015260003073ffffffffffffffffffffffffffffffffffffffff16856040518082805190602001908083835b602083106117e157805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016117a4565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d8060008114611841576040519150601f19603f3d011682016040523d82523d6000602084013e611846565b606091505b50509050806118b657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f556e61626c6520746f2063726561746520726571756573740000000000000000604482015290519081900360640190fd5b50505050505050565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205460ff1690565b60005473ffffffffffffffffffffffffffffffffffffffff16331461197057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81163314156119f557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314611af157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b80611afd816001611fdc565b6005541015611b57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603581526020018061208a6035913960400191505060405180910390fd5b600554611b649083611f65565b6005819055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015611bfb57600080fd5b505af1158015611c0f573d6000803e3d6000fd5b505050506040513d6020811015611c2557600080fd5b5051611c2d57fe5b505050565b604080517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606089901b1660208083019190915260348083018690528351808403909101815260549092018352815191810191909120600081815260029092529181205415611d0257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4d75737420757365206120756e69717565204944000000000000000000000000604482015290519081900360640190fd5b611d0e4261012c611fdc565b6040805160208082018b90527fffffffffffffffffffffffffffffffffffffffff00000000000000000000000060608b901b16828401527fffffffff000000000000000000000000000000000000000000000000000000008916605483015260588083018590528351808403909101815260789092018352815191810191909120600086815260028352838120919091556004909152208490559050965096945050505050565b6040805160208082018890527fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606088901b16828401527fffffffff00000000000000000000000000000000000000000000000000000000861660548301526058808301869052835180840390910181526078909201835281519181019190912060008981526002909252919020548114611eb157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b600087815260046020526040902054821015611f2e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e636f727265637420646174612076657273696f6e00000000000000000000604482015290519081900360640190fd5b600554611f3b9087611fdc565b60055550505060009384525050600260209081526040808420849055600490915282209190915550565b600082821115611fd657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b60008282018381101561205057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b939250505056fe43616e6e6f74207573652023666f727761726420746f2073656e64206d6573736167657320746f204c696e6b20746f6b656e416d6f756e74207265717565737465642069732067726561746572207468616e20776974686472617761626c652062616c616e63654e6f7420616e20617574686f72697a6564206e6f646520746f2066756c66696c6c207265717565737473a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

// DeployOperator deploys a new Ethereum contract, binding an instance of Operator to it.
func DeployOperator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address) (common.Address, *types.Transaction, *Operator, error) {
	parsed, err := abi.JSON(strings.NewReader(OperatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OperatorBin), backend, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Operator{OperatorCaller: OperatorCaller{contract: contract}, OperatorTransactor: OperatorTransactor{contract: contract}, OperatorFilterer: OperatorFilterer{contract: contract}}, nil
}

// Operator is an auto generated Go binding around an Ethereum contract.
type Operator struct {
	OperatorCaller     // Read-only binding to the contract
	OperatorTransactor // Write-only binding to the contract
	OperatorFilterer   // Log filterer for contract events
}

// OperatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OperatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OperatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OperatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OperatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OperatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OperatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OperatorSession struct {
	Contract     *Operator         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OperatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OperatorCallerSession struct {
	Contract *OperatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// OperatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OperatorTransactorSession struct {
	Contract     *OperatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OperatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OperatorRaw struct {
	Contract *Operator // Generic contract binding to access the raw methods on
}

// OperatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OperatorCallerRaw struct {
	Contract *OperatorCaller // Generic read-only contract binding to access the raw methods on
}

// OperatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OperatorTransactorRaw struct {
	Contract *OperatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOperator creates a new instance of Operator, bound to a specific deployed contract.
func NewOperator(address common.Address, backend bind.ContractBackend) (*Operator, error) {
	contract, err := bindOperator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Operator{OperatorCaller: OperatorCaller{contract: contract}, OperatorTransactor: OperatorTransactor{contract: contract}, OperatorFilterer: OperatorFilterer{contract: contract}}, nil
}

// NewOperatorCaller creates a new read-only instance of Operator, bound to a specific deployed contract.
func NewOperatorCaller(address common.Address, caller bind.ContractCaller) (*OperatorCaller, error) {
	contract, err := bindOperator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OperatorCaller{contract: contract}, nil
}

// NewOperatorTransactor creates a new write-only instance of Operator, bound to a specific deployed contract.
func NewOperatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OperatorTransactor, error) {
	contract, err := bindOperator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OperatorTransactor{contract: contract}, nil
}

// NewOperatorFilterer creates a new log filterer instance of Operator, bound to a specific deployed contract.
func NewOperatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OperatorFilterer, error) {
	contract, err := bindOperator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OperatorFilterer{contract: contract}, nil
}

// bindOperator binds a generic wrapper to an already deployed contract.
func bindOperator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OperatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Operator *OperatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Operator.Contract.OperatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Operator *OperatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operator.Contract.OperatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Operator *OperatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Operator.Contract.OperatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Operator *OperatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Operator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Operator *OperatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Operator *OperatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Operator.Contract.contract.Transact(opts, method, params...)
}

// EXPIRYTIME is a free data retrieval call binding the contract method 0x4b602282.
//
// Solidity: function EXPIRY_TIME() view returns(uint256)
func (_Operator *OperatorCaller) EXPIRYTIME(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Operator.contract.Call(opts, out, "EXPIRY_TIME")
	return *ret0, err
}

// EXPIRYTIME is a free data retrieval call binding the contract method 0x4b602282.
//
// Solidity: function EXPIRY_TIME() view returns(uint256)
func (_Operator *OperatorSession) EXPIRYTIME() (*big.Int, error) {
	return _Operator.Contract.EXPIRYTIME(&_Operator.CallOpts)
}

// EXPIRYTIME is a free data retrieval call binding the contract method 0x4b602282.
//
// Solidity: function EXPIRY_TIME() view returns(uint256)
func (_Operator *OperatorCallerSession) EXPIRYTIME() (*big.Int, error) {
	return _Operator.Contract.EXPIRYTIME(&_Operator.CallOpts)
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address node) view returns(bool)
func (_Operator *OperatorCaller) GetAuthorizationStatus(opts *bind.CallOpts, node common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Operator.contract.Call(opts, out, "getAuthorizationStatus", node)
	return *ret0, err
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address node) view returns(bool)
func (_Operator *OperatorSession) GetAuthorizationStatus(node common.Address) (bool, error) {
	return _Operator.Contract.GetAuthorizationStatus(&_Operator.CallOpts, node)
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address node) view returns(bool)
func (_Operator *OperatorCallerSession) GetAuthorizationStatus(node common.Address) (bool, error) {
	return _Operator.Contract.GetAuthorizationStatus(&_Operator.CallOpts, node)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Operator *OperatorCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Operator.contract.Call(opts, out, "getChainlinkToken")
	return *ret0, err
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Operator *OperatorSession) GetChainlinkToken() (common.Address, error) {
	return _Operator.Contract.GetChainlinkToken(&_Operator.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Operator *OperatorCallerSession) GetChainlinkToken() (common.Address, error) {
	return _Operator.Contract.GetChainlinkToken(&_Operator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Operator *OperatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Operator.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Operator *OperatorSession) Owner() (common.Address, error) {
	return _Operator.Contract.Owner(&_Operator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Operator *OperatorCallerSession) Owner() (common.Address, error) {
	return _Operator.Contract.Owner(&_Operator.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_Operator *OperatorCaller) Withdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Operator.contract.Call(opts, out, "withdrawable")
	return *ret0, err
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_Operator *OperatorSession) Withdrawable() (*big.Int, error) {
	return _Operator.Contract.Withdrawable(&_Operator.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_Operator *OperatorCallerSession) Withdrawable() (*big.Int, error) {
	return _Operator.Contract.Withdrawable(&_Operator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Operator *OperatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Operator *OperatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Operator.Contract.AcceptOwnership(&_Operator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Operator *OperatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Operator.Contract.AcceptOwnership(&_Operator.TransactOpts)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 requestId, uint256 payment, bytes4 callbackFunc, uint256 expiration) returns()
func (_Operator *OperatorTransactor) CancelOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "cancelOracleRequest", requestId, payment, callbackFunc, expiration)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 requestId, uint256 payment, bytes4 callbackFunc, uint256 expiration) returns()
func (_Operator *OperatorSession) CancelOracleRequest(requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.CancelOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackFunc, expiration)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 requestId, uint256 payment, bytes4 callbackFunc, uint256 expiration) returns()
func (_Operator *OperatorTransactorSession) CancelOracleRequest(requestId [32]byte, payment *big.Int, callbackFunc [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.CancelOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackFunc, expiration)
}

// Forward is a paid mutator transaction binding the contract method 0x6fadcf72.
//
// Solidity: function forward(address _to, bytes _data) returns()
func (_Operator *OperatorTransactor) Forward(opts *bind.TransactOpts, _to common.Address, _data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "forward", _to, _data)
}

// Forward is a paid mutator transaction binding the contract method 0x6fadcf72.
//
// Solidity: function forward(address _to, bytes _data) returns()
func (_Operator *OperatorSession) Forward(_to common.Address, _data []byte) (*types.Transaction, error) {
	return _Operator.Contract.Forward(&_Operator.TransactOpts, _to, _data)
}

// Forward is a paid mutator transaction binding the contract method 0x6fadcf72.
//
// Solidity: function forward(address _to, bytes _data) returns()
func (_Operator *OperatorTransactorSession) Forward(_to common.Address, _data []byte) (*types.Transaction, error) {
	return _Operator.Contract.Forward(&_Operator.TransactOpts, _to, _data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data) returns(bool)
func (_Operator *OperatorTransactor) FulfillOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "fulfillOracleRequest", requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data) returns(bool)
func (_Operator *OperatorSession) FulfillOracleRequest(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data) returns(bool)
func (_Operator *OperatorTransactorSession) FulfillOracleRequest(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest2 is a paid mutator transaction binding the contract method 0x6ae0bc76.
//
// Solidity: function fulfillOracleRequest2(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes data) returns(bool)
func (_Operator *OperatorTransactor) FulfillOracleRequest2(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "fulfillOracleRequest2", requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest2 is a paid mutator transaction binding the contract method 0x6ae0bc76.
//
// Solidity: function fulfillOracleRequest2(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes data) returns(bool)
func (_Operator *OperatorSession) FulfillOracleRequest2(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest2(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest2 is a paid mutator transaction binding the contract method 0x6ae0bc76.
//
// Solidity: function fulfillOracleRequest2(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes data) returns(bool)
func (_Operator *OperatorTransactorSession) FulfillOracleRequest2(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.FulfillOracleRequest2(&_Operator.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_Operator *OperatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "onTokenTransfer", _sender, _amount, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_Operator *OperatorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OnTokenTransfer(&_Operator.TransactOpts, _sender, _amount, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_Operator *OperatorTransactorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OnTokenTransfer(&_Operator.TransactOpts, _sender, _amount, _data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address sender, uint256 payment, bytes32 specId, address callbackAddress, bytes4 callbackFunctionId, uint256 nonce, uint256 dataVersion, bytes data) returns()
func (_Operator *OperatorTransactor) OracleRequest(opts *bind.TransactOpts, sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "oracleRequest", sender, payment, specId, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address sender, uint256 payment, bytes32 specId, address callbackAddress, bytes4 callbackFunctionId, uint256 nonce, uint256 dataVersion, bytes data) returns()
func (_Operator *OperatorSession) OracleRequest(sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OracleRequest(&_Operator.TransactOpts, sender, payment, specId, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address sender, uint256 payment, bytes32 specId, address callbackAddress, bytes4 callbackFunctionId, uint256 nonce, uint256 dataVersion, bytes data) returns()
func (_Operator *OperatorTransactorSession) OracleRequest(sender common.Address, payment *big.Int, specId [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _Operator.Contract.OracleRequest(&_Operator.TransactOpts, sender, payment, specId, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address node, bool allowed) returns()
func (_Operator *OperatorTransactor) SetFulfillmentPermission(opts *bind.TransactOpts, node common.Address, allowed bool) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "setFulfillmentPermission", node, allowed)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address node, bool allowed) returns()
func (_Operator *OperatorSession) SetFulfillmentPermission(node common.Address, allowed bool) (*types.Transaction, error) {
	return _Operator.Contract.SetFulfillmentPermission(&_Operator.TransactOpts, node, allowed)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address node, bool allowed) returns()
func (_Operator *OperatorTransactorSession) SetFulfillmentPermission(node common.Address, allowed bool) (*types.Transaction, error) {
	return _Operator.Contract.SetFulfillmentPermission(&_Operator.TransactOpts, node, allowed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Operator *OperatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Operator *OperatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Operator.Contract.TransferOwnership(&_Operator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Operator *OperatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Operator.Contract.TransferOwnership(&_Operator.TransactOpts, _to)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_Operator *OperatorTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Operator.contract.Transact(opts, "withdraw", recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_Operator *OperatorSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.Withdraw(&_Operator.TransactOpts, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_Operator *OperatorTransactorSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Operator.Contract.Withdraw(&_Operator.TransactOpts, recipient, amount)
}

// OperatorCancelOracleRequestIterator is returned from FilterCancelOracleRequest and is used to iterate over the raw logs and unpacked data for CancelOracleRequest events raised by the Operator contract.
type OperatorCancelOracleRequestIterator struct {
	Event *OperatorCancelOracleRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OperatorCancelOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OperatorCancelOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OperatorCancelOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OperatorCancelOracleRequest represents a CancelOracleRequest event raised by the Operator contract.
type OperatorCancelOracleRequest struct {
	RequestId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCancelOracleRequest is a free log retrieval operation binding the contract event 0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93.
//
// Solidity: event CancelOracleRequest(bytes32 indexed requestId)
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

// WatchCancelOracleRequest is a free log subscription operation binding the contract event 0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93.
//
// Solidity: event CancelOracleRequest(bytes32 indexed requestId)
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
				// New log arrived, parse the event and forward to the user
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

// ParseCancelOracleRequest is a log parse operation binding the contract event 0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93.
//
// Solidity: event CancelOracleRequest(bytes32 indexed requestId)
func (_Operator *OperatorFilterer) ParseCancelOracleRequest(log types.Log) (*OperatorCancelOracleRequest, error) {
	event := new(OperatorCancelOracleRequest)
	if err := _Operator.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OperatorOracleRequestIterator is returned from FilterOracleRequest and is used to iterate over the raw logs and unpacked data for OracleRequest events raised by the Operator contract.
type OperatorOracleRequestIterator struct {
	Event *OperatorOracleRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OperatorOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OperatorOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OperatorOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OperatorOracleRequest represents a OracleRequest event raised by the Operator contract.
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
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOracleRequest is a free log retrieval operation binding the contract event 0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65.
//
// Solidity: event OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)
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

// WatchOracleRequest is a free log subscription operation binding the contract event 0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65.
//
// Solidity: event OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOracleRequest is a log parse operation binding the contract event 0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65.
//
// Solidity: event OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)
func (_Operator *OperatorFilterer) ParseOracleRequest(log types.Log) (*OperatorOracleRequest, error) {
	event := new(OperatorOracleRequest)
	if err := _Operator.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OperatorOracleResponseIterator is returned from FilterOracleResponse and is used to iterate over the raw logs and unpacked data for OracleResponse events raised by the Operator contract.
type OperatorOracleResponseIterator struct {
	Event *OperatorOracleResponse // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OperatorOracleResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OperatorOracleResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OperatorOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OperatorOracleResponse represents a OracleResponse event raised by the Operator contract.
type OperatorOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleResponse is a free log retrieval operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
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

// WatchOracleResponse is a free log subscription operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOracleResponse is a log parse operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
func (_Operator *OperatorFilterer) ParseOracleResponse(log types.Log) (*OperatorOracleResponse, error) {
	event := new(OperatorOracleResponse)
	if err := _Operator.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OperatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Operator contract.
type OperatorOwnershipTransferRequestedIterator struct {
	Event *OperatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OperatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OperatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OperatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OperatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Operator contract.
type OperatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Operator *OperatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OperatorOwnershipTransferRequested, error) {
	event := new(OperatorOwnershipTransferRequested)
	if err := _Operator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OperatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Operator contract.
type OperatorOwnershipTransferredIterator struct {
	Event *OperatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OperatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OperatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OperatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OperatorOwnershipTransferred represents a OwnershipTransferred event raised by the Operator contract.
type OperatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Operator *OperatorFilterer) ParseOwnershipTransferred(log types.Log) (*OperatorOwnershipTransferred, error) {
	event := new(OperatorOwnershipTransferred)
	if err := _Operator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
