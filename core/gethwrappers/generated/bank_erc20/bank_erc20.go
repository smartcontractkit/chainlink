// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bank_erc20

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

var BankERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ccipRouter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ccipFeeProvider\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ccipChainId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"MustBeCCIPFeeProvider\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"MustBeTrustedForwarder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawFailure\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPChainId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPFeeProvider\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPRouter\",\"outputs\":[{\"internalType\":\"contractIRouterClient\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTrustedForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"destinationChainId\",\"type\":\"uint64\"}],\"name\":\"metaTransfer\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162001f1638038062001f16833981016040819052620000359162000404565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c38162000185565b5050506001600160a01b0384161580620000e457506001600160a01b038316155b80620000f757506001600160a01b038216155b15620001165760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0393841660805291831660a05290911660c0526001600160401b031660e05260056200014a888262000559565b50600662000159878262000559565b5062000178620001716000546001600160a01b031690565b8662000230565b505050505050506200064c565b336001600160a01b03821603620001df5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b038216620002885760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640162000087565b80600460008282546200029c919062000625565b90915550506001600160a01b03821660009081526002602052604081208054839290620002cb90849062000625565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b505050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200034257600080fd5b81516001600160401b03808211156200035f576200035f6200031a565b604051601f8301601f19908116603f011681019082821181831017156200038a576200038a6200031a565b81604052838152602092508683858801011115620003a757600080fd5b600091505b83821015620003cb5785820183015181830184015290820190620003ac565b83821115620003dd5760008385830101525b9695505050505050565b80516001600160a01b0381168114620003ff57600080fd5b919050565b600080600080600080600060e0888a0312156200042057600080fd5b87516001600160401b03808211156200043857600080fd5b620004468b838c0162000330565b985060208a01519150808211156200045d57600080fd5b6200046b8b838c0162000330565b975060408a015196506200048260608b01620003e7565b95506200049260808b01620003e7565b9450620004a260a08b01620003e7565b935060c08a015191508082168214620004ba57600080fd5b508091505092959891949750929550565b600181811c90821680620004e057607f821691505b6020821081036200050157634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200031557600081815260208120601f850160051c81016020861015620005305750805b601f850160051c820191505b8181101562000551578281556001016200053c565b505050505050565b81516001600160401b038111156200057557620005756200031a565b6200058d81620005868454620004cb565b8462000507565b602080601f831160018114620005c55760008415620005ac5750858301515b600019600386901b1c1916600185901b17855562000551565b600085815260208120601f198616915b82811015620005f657888601518255948401946001909101908401620005d5565b5085821015620006155787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b600082198211156200064757634e487b7160e01b600052601160045260246000fd5b500190565b60805160a05160c05160e051611851620006c56000396000818161035f01526106670152600081816104f601528181610b130152610b8b0152600081816102ff015281816107e80152818161084d01526108c90152600081816102b301528181610430015281816105ed0152610e6401526118516000f3fe6080604052600436106101795760003560e01c806370a08231116100cb578063a457c2d71161007f578063dd62ed3e11610059578063dd62ed3e14610494578063ef76fb3e146104e7578063f2fde38b1461051a57600080fd5b8063a457c2d714610454578063a9059cbb14610474578063ce1b815f1461042157600080fd5b80638da5cb5b116100b05780638da5cb5b146103e157806395d89b411461040c578063a00425261461042157600080fd5b806370a082311461038957806379ba5097146103cc57600080fd5b8063313ce5671161012d578063572b6c0511610107578063572b6c0514610296578063588cbd0e146102f0578063593a61b01461034457600080fd5b8063313ce56714610243578063395093511461025f57806350431ce41461027f57600080fd5b8063178293441161015e57806317829344146101e057806318160ddd1461020e57806323b872dd1461022357600080fd5b806306fdde0314610185578063095ea7b3146101b057600080fd5b3661018057005b600080fd5b34801561019157600080fd5b5061019a61053a565b6040516101a7919061151b565b60405180910390f35b3480156101bc57600080fd5b506101d06101cb366004611557565b6105cc565b60405190151581526020016101a7565b3480156101ec57600080fd5b506102006101fb366004611581565b6105e9565b6040519081526020016101a7565b34801561021a57600080fd5b50600454610200565b34801561022f57600080fd5b506101d061023e3660046115ce565b610973565b34801561024f57600080fd5b50604051601281526020016101a7565b34801561026b57600080fd5b506101d061027a366004611557565b610a9a565b34801561028b57600080fd5b50610294610afb565b005b3480156102a257600080fd5b506101d06102b136600461160a565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b3480156102fc57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101a7565b34801561035057600080fd5b5060405167ffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101a7565b34801561039557600080fd5b506102006103a436600461160a565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205490565b3480156103d857600080fd5b50610294610c2a565b3480156103ed57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661031f565b34801561041857600080fd5b5061019a610d27565b34801561042d57600080fd5b507f000000000000000000000000000000000000000000000000000000000000000061031f565b34801561046057600080fd5b506101d061046f366004611557565b610d36565b34801561048057600080fd5b506101d061048f366004611557565b610e2c565b3480156104a057600080fd5b506102006104af366004611625565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260036020908152604080832093909416825291909152205490565b3480156104f357600080fd5b507f000000000000000000000000000000000000000000000000000000000000000061031f565b34801561052657600080fd5b5061029461053536600461160a565b610e40565b60606005805461054990611658565b80601f016020809104026020016040519081016040528092919081815260200182805461057590611658565b80156105c25780601f10610597576101008083540402835291602001916105c2565b820191906000526020600020905b8154815290600101906020018083116105a557829003601f168201915b5050505050905090565b60006105e06105d9610e54565b8484610ed1565b50600192915050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163314610661576040517fa2f64cc50000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b610699827f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff9081169116141590565b6106b7576106af6106a8610e54565b8585611084565b50600061096c565b604080516001808252818301909252600091816020015b60408051808201909152600080825260208201528152602001906001900390816106ce57905050905060405180604001604052803073ffffffffffffffffffffffffffffffffffffffff1681526020018581525081600081518110610735576107356116ab565b60209081029190910101526040805160a0810190915273ffffffffffffffffffffffffffffffffffffffff861660c08201526000908060e081016040516020818303038152906040528152602001604051806020016040528060008152508152602001838152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020016040518060200160405280600081525081525090506107e26107db610e54565b3087611084565b61080d307f000000000000000000000000000000000000000000000000000000000000000087610ed1565b6040517f20487ded00000000000000000000000000000000000000000000000000000000815260009073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906320487ded9061088490889086906004016116da565b602060405180830381865afa1580156108a1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108c591906117ec565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166396f4e9f98287856040518463ffffffff1660e01b81526004016109239291906116da565b60206040518083038185885af1158015610941573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061096691906117ec565b93505050505b9392505050565b6000610980848484611084565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600360205260408120816109ae610e54565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905082811015610a7b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206160448201527f6c6c6f77616e63650000000000000000000000000000000000000000000000006064820152608401610658565b610a8f85610a87610e54565b858403610ed1565b506001949350505050565b60006105e0610aa7610e54565b848460036000610ab5610e54565b73ffffffffffffffffffffffffffffffffffffffff908116825260208083019390935260409182016000908120918b1681529252902054610af69190611805565b610ed1565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610b6c576040517fb0c9d2e7000000000000000000000000000000000000000000000000000000008152336004820152602401610658565b604051479060009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169083908381818185875af1925050503d8060008114610be6576040519150601f19603f3d011682016040523d82523d6000602084013e610beb565b606091505b5050905080610c26576040517f1a0263ed00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610cab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610658565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60606006805461054990611658565b60008060036000610d45610e54565b73ffffffffffffffffffffffffffffffffffffffff90811682526020808301939093526040918201600090812091881681529252902054905082811015610e0e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f0000000000000000000000000000000000000000000000000000006064820152608401610658565b610e22610e19610e54565b85858403610ed1565b5060019392505050565b60006105e0610e39610e54565b8484611084565b610e48611338565b610e51816113bb565b50565b600060143610801590610e9c57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1633145b15610ecc57507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec36013560601c90565b503390565b73ffffffffffffffffffffffffffffffffffffffff8316610f73576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f72657373000000000000000000000000000000000000000000000000000000006064820152608401610658565b73ffffffffffffffffffffffffffffffffffffffff8216611016576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152608401610658565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316611127576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152608401610658565b73ffffffffffffffffffffffffffffffffffffffff82166111ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f65737300000000000000000000000000000000000000000000000000000000006064820152608401610658565b73ffffffffffffffffffffffffffffffffffffffff831660009081526002602052604090205481811015611280576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e636500000000000000000000000000000000000000000000000000006064820152608401610658565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600260205260408082208585039055918516815290812080548492906112c4908490611805565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161132a91815260200190565b60405180910390a350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146113b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610658565b565b3373ffffffffffffffffffffffffffffffffffffffff82160361143a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610658565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000815180845260005b818110156114d6576020818501810151868301820152016114ba565b818111156114e8576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061096c60208301846114b0565b803573ffffffffffffffffffffffffffffffffffffffff8116811461155257600080fd5b919050565b6000806040838503121561156a57600080fd5b6115738361152e565b946020939093013593505050565b60008060006060848603121561159657600080fd5b61159f8461152e565b925060208401359150604084013567ffffffffffffffff811681146115c357600080fd5b809150509250925092565b6000806000606084860312156115e357600080fd5b6115ec8461152e565b92506115fa6020850161152e565b9150604084013590509250925092565b60006020828403121561161c57600080fd5b61096c8261152e565b6000806040838503121561163857600080fd5b6116418361152e565b915061164f6020840161152e565b90509250929050565b600181811c9082168061166c57607f821691505b6020821081036116a5577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000604067ffffffffffffffff8516835260208181850152845160a08386015261170760e08601826114b0565b9050818601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08087840301606088015261174283836114b0565b88860151888203830160808a01528051808352908601945060009350908501905b808410156117a2578451805173ffffffffffffffffffffffffffffffffffffffff16835286015186830152938501936001939093019290860190611763565b50606089015173ffffffffffffffffffffffffffffffffffffffff1660a08901526080890151888203830160c08a015295506117de81876114b0565b9a9950505050505050505050565b6000602082840312156117fe57600080fd5b5051919050565b6000821982111561183f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50019056fea164736f6c634300080f000a",
}

var BankERC20ABI = BankERC20MetaData.ABI

var BankERC20Bin = BankERC20MetaData.Bin

func DeployBankERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string, totalSupply_ *big.Int, forwarder common.Address, ccipRouter common.Address, ccipFeeProvider common.Address, ccipChainId uint64) (common.Address, *types.Transaction, *BankERC20, error) {
	parsed, err := BankERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BankERC20Bin), backend, name_, symbol_, totalSupply_, forwarder, ccipRouter, ccipFeeProvider, ccipChainId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BankERC20{BankERC20Caller: BankERC20Caller{contract: contract}, BankERC20Transactor: BankERC20Transactor{contract: contract}, BankERC20Filterer: BankERC20Filterer{contract: contract}}, nil
}

type BankERC20 struct {
	address common.Address
	abi     abi.ABI
	BankERC20Caller
	BankERC20Transactor
	BankERC20Filterer
}

type BankERC20Caller struct {
	contract *bind.BoundContract
}

type BankERC20Transactor struct {
	contract *bind.BoundContract
}

type BankERC20Filterer struct {
	contract *bind.BoundContract
}

type BankERC20Session struct {
	Contract     *BankERC20
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BankERC20CallerSession struct {
	Contract *BankERC20Caller
	CallOpts bind.CallOpts
}

type BankERC20TransactorSession struct {
	Contract     *BankERC20Transactor
	TransactOpts bind.TransactOpts
}

type BankERC20Raw struct {
	Contract *BankERC20
}

type BankERC20CallerRaw struct {
	Contract *BankERC20Caller
}

type BankERC20TransactorRaw struct {
	Contract *BankERC20Transactor
}

func NewBankERC20(address common.Address, backend bind.ContractBackend) (*BankERC20, error) {
	abi, err := abi.JSON(strings.NewReader(BankERC20ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBankERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BankERC20{address: address, abi: abi, BankERC20Caller: BankERC20Caller{contract: contract}, BankERC20Transactor: BankERC20Transactor{contract: contract}, BankERC20Filterer: BankERC20Filterer{contract: contract}}, nil
}

func NewBankERC20Caller(address common.Address, caller bind.ContractCaller) (*BankERC20Caller, error) {
	contract, err := bindBankERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BankERC20Caller{contract: contract}, nil
}

func NewBankERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*BankERC20Transactor, error) {
	contract, err := bindBankERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BankERC20Transactor{contract: contract}, nil
}

func NewBankERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*BankERC20Filterer, error) {
	contract, err := bindBankERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BankERC20Filterer{contract: contract}, nil
}

func bindBankERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BankERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BankERC20 *BankERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BankERC20.Contract.BankERC20Caller.contract.Call(opts, result, method, params...)
}

func (_BankERC20 *BankERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankERC20.Contract.BankERC20Transactor.contract.Transfer(opts)
}

func (_BankERC20 *BankERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BankERC20.Contract.BankERC20Transactor.contract.Transact(opts, method, params...)
}

func (_BankERC20 *BankERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BankERC20.Contract.contract.Call(opts, result, method, params...)
}

func (_BankERC20 *BankERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankERC20.Contract.contract.Transfer(opts)
}

func (_BankERC20 *BankERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BankERC20.Contract.contract.Transact(opts, method, params...)
}

func (_BankERC20 *BankERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BankERC20 *BankERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BankERC20.Contract.Allowance(&_BankERC20.CallOpts, owner, spender)
}

func (_BankERC20 *BankERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BankERC20.Contract.Allowance(&_BankERC20.CallOpts, owner, spender)
}

func (_BankERC20 *BankERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BankERC20 *BankERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _BankERC20.Contract.BalanceOf(&_BankERC20.CallOpts, account)
}

func (_BankERC20 *BankERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _BankERC20.Contract.BalanceOf(&_BankERC20.CallOpts, account)
}

func (_BankERC20 *BankERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BankERC20 *BankERC20Session) Decimals() (uint8, error) {
	return _BankERC20.Contract.Decimals(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) Decimals() (uint8, error) {
	return _BankERC20.Contract.Decimals(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) GetCCIPChainId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "getCCIPChainId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_BankERC20 *BankERC20Session) GetCCIPChainId() (uint64, error) {
	return _BankERC20.Contract.GetCCIPChainId(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) GetCCIPChainId() (uint64, error) {
	return _BankERC20.Contract.GetCCIPChainId(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) GetCCIPFeeProvider(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "getCCIPFeeProvider")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BankERC20 *BankERC20Session) GetCCIPFeeProvider() (common.Address, error) {
	return _BankERC20.Contract.GetCCIPFeeProvider(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) GetCCIPFeeProvider() (common.Address, error) {
	return _BankERC20.Contract.GetCCIPFeeProvider(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) GetCCIPRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "getCCIPRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BankERC20 *BankERC20Session) GetCCIPRouter() (common.Address, error) {
	return _BankERC20.Contract.GetCCIPRouter(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) GetCCIPRouter() (common.Address, error) {
	return _BankERC20.Contract.GetCCIPRouter(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) GetForwarder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "getForwarder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BankERC20 *BankERC20Session) GetForwarder() (common.Address, error) {
	return _BankERC20.Contract.GetForwarder(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) GetForwarder() (common.Address, error) {
	return _BankERC20.Contract.GetForwarder(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) GetTrustedForwarder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "getTrustedForwarder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BankERC20 *BankERC20Session) GetTrustedForwarder() (common.Address, error) {
	return _BankERC20.Contract.GetTrustedForwarder(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) GetTrustedForwarder() (common.Address, error) {
	return _BankERC20.Contract.GetTrustedForwarder(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BankERC20 *BankERC20Session) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _BankERC20.Contract.IsTrustedForwarder(&_BankERC20.CallOpts, forwarder)
}

func (_BankERC20 *BankERC20CallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _BankERC20.Contract.IsTrustedForwarder(&_BankERC20.CallOpts, forwarder)
}

func (_BankERC20 *BankERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BankERC20 *BankERC20Session) Name() (string, error) {
	return _BankERC20.Contract.Name(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) Name() (string, error) {
	return _BankERC20.Contract.Name(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BankERC20 *BankERC20Session) Owner() (common.Address, error) {
	return _BankERC20.Contract.Owner(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) Owner() (common.Address, error) {
	return _BankERC20.Contract.Owner(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BankERC20 *BankERC20Session) Symbol() (string, error) {
	return _BankERC20.Contract.Symbol(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) Symbol() (string, error) {
	return _BankERC20.Contract.Symbol(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BankERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BankERC20 *BankERC20Session) TotalSupply() (*big.Int, error) {
	return _BankERC20.Contract.TotalSupply(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _BankERC20.Contract.TotalSupply(&_BankERC20.CallOpts)
}

func (_BankERC20 *BankERC20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "acceptOwnership")
}

func (_BankERC20 *BankERC20Session) AcceptOwnership() (*types.Transaction, error) {
	return _BankERC20.Contract.AcceptOwnership(&_BankERC20.TransactOpts)
}

func (_BankERC20 *BankERC20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BankERC20.Contract.AcceptOwnership(&_BankERC20.TransactOpts)
}

func (_BankERC20 *BankERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "approve", spender, amount)
}

func (_BankERC20 *BankERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.Approve(&_BankERC20.TransactOpts, spender, amount)
}

func (_BankERC20 *BankERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.Approve(&_BankERC20.TransactOpts, spender, amount)
}

func (_BankERC20 *BankERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_BankERC20 *BankERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.DecreaseAllowance(&_BankERC20.TransactOpts, spender, subtractedValue)
}

func (_BankERC20 *BankERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.DecreaseAllowance(&_BankERC20.TransactOpts, spender, subtractedValue)
}

func (_BankERC20 *BankERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_BankERC20 *BankERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.IncreaseAllowance(&_BankERC20.TransactOpts, spender, addedValue)
}

func (_BankERC20 *BankERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.IncreaseAllowance(&_BankERC20.TransactOpts, spender, addedValue)
}

func (_BankERC20 *BankERC20Transactor) MetaTransfer(opts *bind.TransactOpts, receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "metaTransfer", receiver, amount, destinationChainId)
}

func (_BankERC20 *BankERC20Session) MetaTransfer(receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error) {
	return _BankERC20.Contract.MetaTransfer(&_BankERC20.TransactOpts, receiver, amount, destinationChainId)
}

func (_BankERC20 *BankERC20TransactorSession) MetaTransfer(receiver common.Address, amount *big.Int, destinationChainId uint64) (*types.Transaction, error) {
	return _BankERC20.Contract.MetaTransfer(&_BankERC20.TransactOpts, receiver, amount, destinationChainId)
}

func (_BankERC20 *BankERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "transfer", recipient, amount)
}

func (_BankERC20 *BankERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.Transfer(&_BankERC20.TransactOpts, recipient, amount)
}

func (_BankERC20 *BankERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.Transfer(&_BankERC20.TransactOpts, recipient, amount)
}

func (_BankERC20 *BankERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

func (_BankERC20 *BankERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.TransferFrom(&_BankERC20.TransactOpts, sender, recipient, amount)
}

func (_BankERC20 *BankERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BankERC20.Contract.TransferFrom(&_BankERC20.TransactOpts, sender, recipient, amount)
}

func (_BankERC20 *BankERC20Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "transferOwnership", to)
}

func (_BankERC20 *BankERC20Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BankERC20.Contract.TransferOwnership(&_BankERC20.TransactOpts, to)
}

func (_BankERC20 *BankERC20TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BankERC20.Contract.TransferOwnership(&_BankERC20.TransactOpts, to)
}

func (_BankERC20 *BankERC20Transactor) WithdrawNative(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankERC20.contract.Transact(opts, "withdrawNative")
}

func (_BankERC20 *BankERC20Session) WithdrawNative() (*types.Transaction, error) {
	return _BankERC20.Contract.WithdrawNative(&_BankERC20.TransactOpts)
}

func (_BankERC20 *BankERC20TransactorSession) WithdrawNative() (*types.Transaction, error) {
	return _BankERC20.Contract.WithdrawNative(&_BankERC20.TransactOpts)
}

func (_BankERC20 *BankERC20Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankERC20.contract.RawTransact(opts, nil)
}

func (_BankERC20 *BankERC20Session) Receive() (*types.Transaction, error) {
	return _BankERC20.Contract.Receive(&_BankERC20.TransactOpts)
}

func (_BankERC20 *BankERC20TransactorSession) Receive() (*types.Transaction, error) {
	return _BankERC20.Contract.Receive(&_BankERC20.TransactOpts)
}

type BankERC20ApprovalIterator struct {
	Event *BankERC20Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BankERC20ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BankERC20Approval)
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
		it.Event = new(BankERC20Approval)
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

func (it *BankERC20ApprovalIterator) Error() error {
	return it.fail
}

func (it *BankERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BankERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_BankERC20 *BankERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BankERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BankERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &BankERC20ApprovalIterator{contract: _BankERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_BankERC20 *BankERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BankERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BankERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BankERC20Approval)
				if err := _BankERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_BankERC20 *BankERC20Filterer) ParseApproval(log types.Log) (*BankERC20Approval, error) {
	event := new(BankERC20Approval)
	if err := _BankERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BankERC20OwnershipTransferRequestedIterator struct {
	Event *BankERC20OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BankERC20OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BankERC20OwnershipTransferRequested)
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
		it.Event = new(BankERC20OwnershipTransferRequested)
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

func (it *BankERC20OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BankERC20OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BankERC20OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BankERC20 *BankERC20Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BankERC20OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BankERC20.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BankERC20OwnershipTransferRequestedIterator{contract: _BankERC20.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BankERC20 *BankERC20Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BankERC20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BankERC20.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BankERC20OwnershipTransferRequested)
				if err := _BankERC20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BankERC20 *BankERC20Filterer) ParseOwnershipTransferRequested(log types.Log) (*BankERC20OwnershipTransferRequested, error) {
	event := new(BankERC20OwnershipTransferRequested)
	if err := _BankERC20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BankERC20OwnershipTransferredIterator struct {
	Event *BankERC20OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BankERC20OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BankERC20OwnershipTransferred)
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
		it.Event = new(BankERC20OwnershipTransferred)
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

func (it *BankERC20OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BankERC20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BankERC20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BankERC20 *BankERC20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BankERC20OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BankERC20.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BankERC20OwnershipTransferredIterator{contract: _BankERC20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BankERC20 *BankERC20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BankERC20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BankERC20.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BankERC20OwnershipTransferred)
				if err := _BankERC20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BankERC20 *BankERC20Filterer) ParseOwnershipTransferred(log types.Log) (*BankERC20OwnershipTransferred, error) {
	event := new(BankERC20OwnershipTransferred)
	if err := _BankERC20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BankERC20TransferIterator struct {
	Event *BankERC20Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BankERC20TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BankERC20Transfer)
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
		it.Event = new(BankERC20Transfer)
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

func (it *BankERC20TransferIterator) Error() error {
	return it.fail
}

func (it *BankERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BankERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_BankERC20 *BankERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BankERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BankERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BankERC20TransferIterator{contract: _BankERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_BankERC20 *BankERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BankERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BankERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BankERC20Transfer)
				if err := _BankERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_BankERC20 *BankERC20Filterer) ParseTransfer(log types.Log) (*BankERC20Transfer, error) {
	event := new(BankERC20Transfer)
	if err := _BankERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BankERC20 *BankERC20) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BankERC20.abi.Events["Approval"].ID:
		return _BankERC20.ParseApproval(log)
	case _BankERC20.abi.Events["OwnershipTransferRequested"].ID:
		return _BankERC20.ParseOwnershipTransferRequested(log)
	case _BankERC20.abi.Events["OwnershipTransferred"].ID:
		return _BankERC20.ParseOwnershipTransferred(log)
	case _BankERC20.abi.Events["Transfer"].ID:
		return _BankERC20.ParseTransfer(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BankERC20Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (BankERC20OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BankERC20OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BankERC20Transfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (_BankERC20 *BankERC20) Address() common.Address {
	return _BankERC20.address
}

type BankERC20Interface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	GetCCIPChainId(opts *bind.CallOpts) (uint64, error)

	GetCCIPFeeProvider(opts *bind.CallOpts) (common.Address, error)

	GetCCIPRouter(opts *bind.CallOpts) (common.Address, error)

	GetForwarder(opts *bind.CallOpts) (common.Address, error)

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

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BankERC20ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *BankERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*BankERC20Approval, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BankERC20OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BankERC20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BankERC20OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BankERC20OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BankERC20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BankERC20OwnershipTransferred, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BankERC20TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *BankERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*BankERC20Transfer, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
