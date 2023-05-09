// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cross_chain_erc20_extension

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

var CrossChainERC20ExtensionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ccipRouter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_test_only_force_cross_chain_transfer\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"MustBeTrustedForwarder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawFailure\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTrustedForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"destinationChainId\",\"type\":\"uint64\"}],\"name\":\"metaTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001d1d38038062001d1d8339810160408190526200003491620003e2565b338060008888600362000048838262000522565b50600462000057828262000522565b5050506001600160a01b038216620000b65760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600580546001600160a01b0319166001600160a01b0384811691909117909155811615620000e957620000e98162000153565b5050506200010762000100620001ff60201b60201c565b856200020e565b600780546001600160a01b039485166001600160a01b031990911617905560088054911515600160a01b026001600160a81b031990921692909316919091171790555062000615915050565b336001600160a01b03821603620001ad5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000ad565b600680546001600160a01b0319166001600160a01b03838116918217909255600554604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b6005546001600160a01b031690565b6001600160a01b038216620002665760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f2061646472657373006044820152606401620000ad565b80600260008282546200027a9190620005ee565b90915550506001600160a01b03821660009081526020819052604081208054839290620002a9908490620005ee565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b505050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200032057600080fd5b81516001600160401b03808211156200033d576200033d620002f8565b604051601f8301601f19908116603f01168101908282118183101715620003685762000368620002f8565b816040528381526020925086838588010111156200038557600080fd5b600091505b83821015620003a957858201830151818301840152908201906200038a565b83821115620003bb5760008385830101525b9695505050505050565b80516001600160a01b0381168114620003dd57600080fd5b919050565b60008060008060008060c08789031215620003fc57600080fd5b86516001600160401b03808211156200041457600080fd5b620004228a838b016200030e565b975060208901519150808211156200043957600080fd5b506200044889828a016200030e565b955050604087015193506200046060608801620003c5565b92506200047060808801620003c5565b915060a087015180151581146200048657600080fd5b809150509295509295509295565b600181811c90821680620004a957607f821691505b602082108103620004ca57634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f357600081815260208120601f850160051c81016020861015620004f95750805b601f850160051c820191505b818110156200051a5782815560010162000505565b505050505050565b81516001600160401b038111156200053e576200053e620002f8565b62000556816200054f845462000494565b84620004d0565b602080601f8311600181146200058e5760008415620005755750858301515b600019600386901b1c1916600185901b1785556200051a565b600085815260208120601f198616915b82811015620005bf578886015182559484019460019091019084016200059e565b5085821015620005de5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b600082198211156200061057634e487b7160e01b600052601160045260246000fd5b500190565b6116f880620006256000396000f3fe60806040526004361061012d5760003560e01c806370a08231116100a5578063a457c2d711610074578063ce1b815f11610059578063ce1b815f1461037b578063dd62ed3e146103a6578063f2fde38b146103f957600080fd5b8063a457c2d71461033b578063a9059cbb1461035b57600080fd5b806370a082311461028257806379ba5097146102c55780638da5cb5b146102da57806395d89b411461032657600080fd5b806323b872dd116100fc57806339509351116100e1578063395093511461021157806350431ce414610231578063572b6c051461024657600080fd5b806323b872dd146101d5578063313ce567146101f557600080fd5b806306fdde0314610139578063095ea7b314610164578063178293441461019457806318160ddd146101b657600080fd5b3661013457005b600080fd5b34801561014557600080fd5b5061014e610419565b60405161015b91906113bb565b60405180910390f35b34801561017057600080fd5b5061018461017f3660046113fe565b6104ab565b604051901515815260200161015b565b3480156101a057600080fd5b506101b46101af366004611428565b6104c8565b005b3480156101c257600080fd5b506002545b60405190815260200161015b565b3480156101e157600080fd5b506101846101f0366004611475565b61085e565b34801561020157600080fd5b506040516012815260200161015b565b34801561021d57600080fd5b5061018461022c3660046113fe565b610985565b34801561023d57600080fd5b506101b46109e6565b34801561025257600080fd5b506101846102613660046114b1565b60075473ffffffffffffffffffffffffffffffffffffffff91821691161490565b34801561028e57600080fd5b506101c761029d3660046114b1565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b3480156102d157600080fd5b506101b4610aab565b3480156102e657600080fd5b5060055473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161015b565b34801561033257600080fd5b5061014e610bac565b34801561034757600080fd5b506101846103563660046113fe565b610bbb565b34801561036757600080fd5b506101846103763660046113fe565b610cb1565b34801561038757600080fd5b5060075473ffffffffffffffffffffffffffffffffffffffff16610301565b3480156103b257600080fd5b506101c76103c13660046114cc565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b34801561040557600080fd5b506101b46104143660046114b1565b610cbe565b606060038054610428906114ff565b80601f0160208091040260200160405190810160405280929190818152602001828054610454906114ff565b80156104a15780601f10610476576101008083540402835291602001916104a1565b820191906000526020600020905b81548152906001019060200180831161048457829003601f168201915b5050505050905090565b60006104bf6104b8610cd2565b8484610d31565b50600192915050565b60075473ffffffffffffffffffffffffffffffffffffffff163314610520576040517fa2f64cc50000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b61052981610ee4565b6105445761053f610538610cd2565b8484610f23565b505050565b604080516001808252818301909252600091816020015b604080518082019091526000808252602082015281526020019060019003908161055b57905050905060405180604001604052803073ffffffffffffffffffffffffffffffffffffffff16815260200184815250816000815181106105c2576105c2611552565b60209081029190910101526040805160a0810190915273ffffffffffffffffffffffffffffffffffffffff851660c08201526000908060e081016040516020818303038152906040528152602001604051806020016040528060008152508152602001838152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020016106de604051806040016040528062030d408152602001600015158152506040805182516024820152602092830151151560448083019190915282518083039091018152606490910190915290810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f97a657c90000000000000000000000000000000000000000000000000000000017905290565b90526008546040517f20487ded00000000000000000000000000000000000000000000000000000000815291925060009173ffffffffffffffffffffffffffffffffffffffff909116906320487ded9061073e9087908690600401611581565b602060405180830381865afa15801561075b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061077f9190611693565b905061079361078c610cd2565b3087610f23565b6008546107b890309073ffffffffffffffffffffffffffffffffffffffff1687610d31565b6008546040517f96f4e9f900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116906396f4e9f99083906108129088908790600401611581565b60206040518083038185885af1158015610830573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906108559190611693565b50505050505050565b600061086b848484610f23565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260016020526040812081610899610cd2565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905082811015610966576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206160448201527f6c6c6f77616e63650000000000000000000000000000000000000000000000006064820152608401610517565b61097a85610972610cd2565b858403610d31565b506001949350505050565b60006104bf610992610cd2565b8484600160006109a0610cd2565b73ffffffffffffffffffffffffffffffffffffffff908116825260208083019390935260409182016000908120918b16815292529020546109e191906116ac565b610d31565b6109ee6111d7565b476000610a1060055473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114610a67576040519150601f19603f3d011682016040523d82523d6000602084013e610a6c565b606091505b5050905080610aa7576040517f1a0263ed00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60065473ffffffffffffffffffffffffffffffffffffffff163314610b2c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610517565b600580547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560068054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b606060048054610428906114ff565b60008060016000610bca610cd2565b73ffffffffffffffffffffffffffffffffffffffff90811682526020808301939093526040918201600090812091881681529252902054905082811015610c93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f0000000000000000000000000000000000000000000000000000006064820152608401610517565b610ca7610c9e610cd2565b85858403610d31565b5060019392505050565b60006104bf610538610cd2565b610cc66111d7565b610ccf8161125a565b50565b600060143610801590610cfc575060075473ffffffffffffffffffffffffffffffffffffffff1633145b15610d2c57507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec36013560601c90565b503390565b73ffffffffffffffffffffffffffffffffffffffff8316610dd3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f72657373000000000000000000000000000000000000000000000000000000006064820152608401610517565b73ffffffffffffffffffffffffffffffffffffffff8216610e76576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152608401610517565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b60085460009074010000000000000000000000000000000000000000900460ff1615610f1257506001919050565b5067ffffffffffffffff1646141590565b73ffffffffffffffffffffffffffffffffffffffff8316610fc6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152608401610517565b73ffffffffffffffffffffffffffffffffffffffff8216611069576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f65737300000000000000000000000000000000000000000000000000000000006064820152608401610517565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602081905260409020548181101561111f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e636500000000000000000000000000000000000000000000000000006064820152608401610517565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152602081905260408082208585039055918516815290812080548492906111639084906116ac565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516111c991815260200190565b60405180910390a350505050565b60055473ffffffffffffffffffffffffffffffffffffffff163314611258576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610517565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036112d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610517565b600680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600554604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b6000815180845260005b818110156113765760208185018101518683018201520161135a565b81811115611388576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006113ce6020830184611350565b9392505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146113f957600080fd5b919050565b6000806040838503121561141157600080fd5b61141a836113d5565b946020939093013593505050565b60008060006060848603121561143d57600080fd5b611446846113d5565b925060208401359150604084013567ffffffffffffffff8116811461146a57600080fd5b809150509250925092565b60008060006060848603121561148a57600080fd5b611493846113d5565b92506114a1602085016113d5565b9150604084013590509250925092565b6000602082840312156114c357600080fd5b6113ce826113d5565b600080604083850312156114df57600080fd5b6114e8836113d5565b91506114f6602084016113d5565b90509250929050565b600181811c9082168061151357607f821691505b60208210810361154c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000604067ffffffffffffffff8516835260208181850152845160a0838601526115ae60e0860182611350565b9050818601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0808784030160608801526115e98383611350565b88860151888203830160808a01528051808352908601945060009350908501905b80841015611649578451805173ffffffffffffffffffffffffffffffffffffffff1683528601518683015293850193600193909301929086019061160a565b50606089015173ffffffffffffffffffffffffffffffffffffffff1660a08901526080890151888203830160c08a015295506116858187611350565b9a9950505050505050505050565b6000602082840312156116a557600080fd5b5051919050565b600082198211156116e6577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50019056fea164736f6c634300080f000a",
}

var CrossChainERC20ExtensionABI = CrossChainERC20ExtensionMetaData.ABI

var CrossChainERC20ExtensionBin = CrossChainERC20ExtensionMetaData.Bin

func DeployCrossChainERC20Extension(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, totalSupply *big.Int, forwarder common.Address, ccipRouter common.Address, _test_only_force_cross_chain_transfer bool) (common.Address, *types.Transaction, *CrossChainERC20Extension, error) {
	parsed, err := CrossChainERC20ExtensionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainERC20ExtensionBin), backend, name, symbol, totalSupply, forwarder, ccipRouter, _test_only_force_cross_chain_transfer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainERC20Extension{CrossChainERC20ExtensionCaller: CrossChainERC20ExtensionCaller{contract: contract}, CrossChainERC20ExtensionTransactor: CrossChainERC20ExtensionTransactor{contract: contract}, CrossChainERC20ExtensionFilterer: CrossChainERC20ExtensionFilterer{contract: contract}}, nil
}

type CrossChainERC20Extension struct {
	address common.Address
	abi     abi.ABI
	CrossChainERC20ExtensionCaller
	CrossChainERC20ExtensionTransactor
	CrossChainERC20ExtensionFilterer
}

type CrossChainERC20ExtensionCaller struct {
	contract *bind.BoundContract
}

type CrossChainERC20ExtensionTransactor struct {
	contract *bind.BoundContract
}

type CrossChainERC20ExtensionFilterer struct {
	contract *bind.BoundContract
}

type CrossChainERC20ExtensionSession struct {
	Contract     *CrossChainERC20Extension
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CrossChainERC20ExtensionCallerSession struct {
	Contract *CrossChainERC20ExtensionCaller
	CallOpts bind.CallOpts
}

type CrossChainERC20ExtensionTransactorSession struct {
	Contract     *CrossChainERC20ExtensionTransactor
	TransactOpts bind.TransactOpts
}

type CrossChainERC20ExtensionRaw struct {
	Contract *CrossChainERC20Extension
}

type CrossChainERC20ExtensionCallerRaw struct {
	Contract *CrossChainERC20ExtensionCaller
}

type CrossChainERC20ExtensionTransactorRaw struct {
	Contract *CrossChainERC20ExtensionTransactor
}

func NewCrossChainERC20Extension(address common.Address, backend bind.ContractBackend) (*CrossChainERC20Extension, error) {
	abi, err := abi.JSON(strings.NewReader(CrossChainERC20ExtensionABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCrossChainERC20Extension(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20Extension{address: address, abi: abi, CrossChainERC20ExtensionCaller: CrossChainERC20ExtensionCaller{contract: contract}, CrossChainERC20ExtensionTransactor: CrossChainERC20ExtensionTransactor{contract: contract}, CrossChainERC20ExtensionFilterer: CrossChainERC20ExtensionFilterer{contract: contract}}, nil
}

func NewCrossChainERC20ExtensionCaller(address common.Address, caller bind.ContractCaller) (*CrossChainERC20ExtensionCaller, error) {
	contract, err := bindCrossChainERC20Extension(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionCaller{contract: contract}, nil
}

func NewCrossChainERC20ExtensionTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainERC20ExtensionTransactor, error) {
	contract, err := bindCrossChainERC20Extension(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionTransactor{contract: contract}, nil
}

func NewCrossChainERC20ExtensionFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainERC20ExtensionFilterer, error) {
	contract, err := bindCrossChainERC20Extension(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionFilterer{contract: contract}, nil
}

func bindCrossChainERC20Extension(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrossChainERC20ExtensionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainERC20Extension.Contract.CrossChainERC20ExtensionCaller.contract.Call(opts, result, method, params...)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.CrossChainERC20ExtensionTransactor.contract.Transfer(opts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.CrossChainERC20ExtensionTransactor.contract.Transact(opts, method, params...)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainERC20Extension.Contract.contract.Call(opts, result, method, params...)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.contract.Transfer(opts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.contract.Transact(opts, method, params...)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CrossChainERC20Extension.Contract.Allowance(&_CrossChainERC20Extension.CallOpts, owner, spender)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CrossChainERC20Extension.Contract.Allowance(&_CrossChainERC20Extension.CallOpts, owner, spender)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _CrossChainERC20Extension.Contract.BalanceOf(&_CrossChainERC20Extension.CallOpts, account)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _CrossChainERC20Extension.Contract.BalanceOf(&_CrossChainERC20Extension.CallOpts, account)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Decimals() (uint8, error) {
	return _CrossChainERC20Extension.Contract.Decimals(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) Decimals() (uint8, error) {
	return _CrossChainERC20Extension.Contract.Decimals(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) GetTrustedForwarder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "getTrustedForwarder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) GetTrustedForwarder() (common.Address, error) {
	return _CrossChainERC20Extension.Contract.GetTrustedForwarder(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) GetTrustedForwarder() (common.Address, error) {
	return _CrossChainERC20Extension.Contract.GetTrustedForwarder(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _CrossChainERC20Extension.Contract.IsTrustedForwarder(&_CrossChainERC20Extension.CallOpts, forwarder)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _CrossChainERC20Extension.Contract.IsTrustedForwarder(&_CrossChainERC20Extension.CallOpts, forwarder)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Name() (string, error) {
	return _CrossChainERC20Extension.Contract.Name(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) Name() (string, error) {
	return _CrossChainERC20Extension.Contract.Name(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Owner() (common.Address, error) {
	return _CrossChainERC20Extension.Contract.Owner(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) Owner() (common.Address, error) {
	return _CrossChainERC20Extension.Contract.Owner(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Symbol() (string, error) {
	return _CrossChainERC20Extension.Contract.Symbol(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) Symbol() (string, error) {
	return _CrossChainERC20Extension.Contract.Symbol(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrossChainERC20Extension.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) TotalSupply() (*big.Int, error) {
	return _CrossChainERC20Extension.Contract.TotalSupply(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionCallerSession) TotalSupply() (*big.Int, error) {
	return _CrossChainERC20Extension.Contract.TotalSupply(&_CrossChainERC20Extension.CallOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "acceptOwnership")
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) AcceptOwnership() (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.AcceptOwnership(&_CrossChainERC20Extension.TransactOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.AcceptOwnership(&_CrossChainERC20Extension.TransactOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "approve", spender, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.Approve(&_CrossChainERC20Extension.TransactOpts, spender, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.Approve(&_CrossChainERC20Extension.TransactOpts, spender, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.DecreaseAllowance(&_CrossChainERC20Extension.TransactOpts, spender, subtractedValue)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.DecreaseAllowance(&_CrossChainERC20Extension.TransactOpts, spender, subtractedValue)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.IncreaseAllowance(&_CrossChainERC20Extension.TransactOpts, spender, addedValue)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.IncreaseAllowance(&_CrossChainERC20Extension.TransactOpts, spender, addedValue)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) MetaTransfer(opts *bind.TransactOpts, receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "metaTransfer", receiver, amount, destinationChainId)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) MetaTransfer(receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.MetaTransfer(&_CrossChainERC20Extension.TransactOpts, receiver, amount, destinationChainId)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) MetaTransfer(receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.MetaTransfer(&_CrossChainERC20Extension.TransactOpts, receiver, amount, destinationChainId)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "transfer", recipient, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.Transfer(&_CrossChainERC20Extension.TransactOpts, recipient, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.Transfer(&_CrossChainERC20Extension.TransactOpts, recipient, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.TransferFrom(&_CrossChainERC20Extension.TransactOpts, sender, recipient, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.TransferFrom(&_CrossChainERC20Extension.TransactOpts, sender, recipient, amount)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "transferOwnership", to)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.TransferOwnership(&_CrossChainERC20Extension.TransactOpts, to)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.TransferOwnership(&_CrossChainERC20Extension.TransactOpts, to)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) WithdrawNative(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.Transact(opts, "withdrawNative")
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) WithdrawNative() (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.WithdrawNative(&_CrossChainERC20Extension.TransactOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) WithdrawNative() (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.WithdrawNative(&_CrossChainERC20Extension.TransactOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainERC20Extension.contract.RawTransact(opts, nil)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionSession) Receive() (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.Receive(&_CrossChainERC20Extension.TransactOpts)
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionTransactorSession) Receive() (*types.Transaction, error) {
	return _CrossChainERC20Extension.Contract.Receive(&_CrossChainERC20Extension.TransactOpts)
}

type CrossChainERC20ExtensionApprovalIterator struct {
	Event *CrossChainERC20ExtensionApproval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CrossChainERC20ExtensionApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainERC20ExtensionApproval)
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
		it.Event = new(CrossChainERC20ExtensionApproval)
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

func (it *CrossChainERC20ExtensionApprovalIterator) Error() error {
	return it.fail
}

func (it *CrossChainERC20ExtensionApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CrossChainERC20ExtensionApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CrossChainERC20ExtensionApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionApprovalIterator{contract: _CrossChainERC20Extension.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CrossChainERC20ExtensionApproval)
				if err := _CrossChainERC20Extension.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) ParseApproval(log types.Log) (*CrossChainERC20ExtensionApproval, error) {
	event := new(CrossChainERC20ExtensionApproval)
	if err := _CrossChainERC20Extension.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CrossChainERC20ExtensionOwnershipTransferRequestedIterator struct {
	Event *CrossChainERC20ExtensionOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CrossChainERC20ExtensionOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainERC20ExtensionOwnershipTransferRequested)
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
		it.Event = new(CrossChainERC20ExtensionOwnershipTransferRequested)
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

func (it *CrossChainERC20ExtensionOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CrossChainERC20ExtensionOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CrossChainERC20ExtensionOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrossChainERC20ExtensionOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionOwnershipTransferRequestedIterator{contract: _CrossChainERC20Extension.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CrossChainERC20ExtensionOwnershipTransferRequested)
				if err := _CrossChainERC20Extension.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) ParseOwnershipTransferRequested(log types.Log) (*CrossChainERC20ExtensionOwnershipTransferRequested, error) {
	event := new(CrossChainERC20ExtensionOwnershipTransferRequested)
	if err := _CrossChainERC20Extension.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CrossChainERC20ExtensionOwnershipTransferredIterator struct {
	Event *CrossChainERC20ExtensionOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CrossChainERC20ExtensionOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainERC20ExtensionOwnershipTransferred)
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
		it.Event = new(CrossChainERC20ExtensionOwnershipTransferred)
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

func (it *CrossChainERC20ExtensionOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CrossChainERC20ExtensionOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CrossChainERC20ExtensionOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrossChainERC20ExtensionOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionOwnershipTransferredIterator{contract: _CrossChainERC20Extension.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CrossChainERC20ExtensionOwnershipTransferred)
				if err := _CrossChainERC20Extension.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) ParseOwnershipTransferred(log types.Log) (*CrossChainERC20ExtensionOwnershipTransferred, error) {
	event := new(CrossChainERC20ExtensionOwnershipTransferred)
	if err := _CrossChainERC20Extension.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CrossChainERC20ExtensionTransferIterator struct {
	Event *CrossChainERC20ExtensionTransfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CrossChainERC20ExtensionTransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainERC20ExtensionTransfer)
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
		it.Event = new(CrossChainERC20ExtensionTransfer)
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

func (it *CrossChainERC20ExtensionTransferIterator) Error() error {
	return it.fail
}

func (it *CrossChainERC20ExtensionTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CrossChainERC20ExtensionTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrossChainERC20ExtensionTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainERC20ExtensionTransferIterator{contract: _CrossChainERC20Extension.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrossChainERC20Extension.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CrossChainERC20ExtensionTransfer)
				if err := _CrossChainERC20Extension.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_CrossChainERC20Extension *CrossChainERC20ExtensionFilterer) ParseTransfer(log types.Log) (*CrossChainERC20ExtensionTransfer, error) {
	event := new(CrossChainERC20ExtensionTransfer)
	if err := _CrossChainERC20Extension.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CrossChainERC20Extension *CrossChainERC20Extension) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CrossChainERC20Extension.abi.Events["Approval"].ID:
		return _CrossChainERC20Extension.ParseApproval(log)
	case _CrossChainERC20Extension.abi.Events["OwnershipTransferRequested"].ID:
		return _CrossChainERC20Extension.ParseOwnershipTransferRequested(log)
	case _CrossChainERC20Extension.abi.Events["OwnershipTransferred"].ID:
		return _CrossChainERC20Extension.ParseOwnershipTransferred(log)
	case _CrossChainERC20Extension.abi.Events["Transfer"].ID:
		return _CrossChainERC20Extension.ParseTransfer(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CrossChainERC20ExtensionApproval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (CrossChainERC20ExtensionOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CrossChainERC20ExtensionOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (CrossChainERC20ExtensionTransfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (_CrossChainERC20Extension *CrossChainERC20Extension) Address() common.Address {
	return _CrossChainERC20Extension.address
}

type CrossChainERC20ExtensionInterface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	GetTrustedForwarder(opts *bind.CallOpts) (common.Address, error)

	IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error)

	Name(opts *bind.CallOpts) (string, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	MetaTransfer(opts *bind.TransactOpts, receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CrossChainERC20ExtensionApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionApproval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*CrossChainERC20ExtensionApproval, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrossChainERC20ExtensionOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CrossChainERC20ExtensionOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrossChainERC20ExtensionOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CrossChainERC20ExtensionOwnershipTransferred, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrossChainERC20ExtensionTransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *CrossChainERC20ExtensionTransfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*CrossChainERC20ExtensionTransfer, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
