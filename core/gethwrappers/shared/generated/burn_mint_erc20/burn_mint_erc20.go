// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_erc20

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

var BurnMintERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"maxSupply_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preMint\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supplyAfterMint\",\"type\":\"uint256\"}],\"name\":\"MaxSupplyExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotBurner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotMinter\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"CCIPAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BURNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burnAndMinter\",\"type\":\"address\"}],\"name\":\"grantMintAndBurnRoles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620022a6380380620022a68339810160408190526200003491620002e7565b848460036200004483826200040b565b5060046200005382826200040b565b50505060ff831660805260a0829052600680546001600160a01b03191633179055801562000087576200008733826200009f565b6200009460003362000166565b5050505050620004f9565b6001600160a01b038216620000fa5760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b80600260008282546200010e9190620004d7565b90915550506001600160a01b038216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35b5050565b620001728282620001f5565b620001625760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055620001ac3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b505050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff165b92915050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200024a57600080fd5b81516001600160401b038082111562000267576200026762000222565b604051601f8301601f19908116603f0116810190828211818310171562000292576200029262000222565b81604052838152602092508683858801011115620002af57600080fd5b600091505b83821015620002d35785820183015181830184015290820190620002b4565b600093810190920192909252949350505050565b600080600080600060a086880312156200030057600080fd5b85516001600160401b03808211156200031857600080fd5b6200032689838a0162000238565b965060208801519150808211156200033d57600080fd5b506200034c8882890162000238565b945050604086015160ff811681146200036457600080fd5b6060870151608090970151959894975095949392505050565b600181811c908216806200039257607f821691505b602082108103620003b357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620001f057600081815260208120601f850160051c81016020861015620003e25750805b601f850160051c820191505b818110156200040357828155600101620003ee565b505050505050565b81516001600160401b0381111562000427576200042762000222565b6200043f816200043884546200037d565b84620003b9565b602080601f8311600181146200047757600084156200045e5750858301515b600019600386901b1c1916600185901b17855562000403565b600085815260208120601f198616915b82811015620004a85788860151825594840194600190910190840162000487565b5085821015620004c75787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b808201808211156200021c57634e487b7160e01b600052601160045260246000fd5b60805160a051611d796200052d600039600081816104a50152818161092a0152610954015260006102ba0152611d796000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c806379cc679011610104578063a8fa343c116100a2578063d547741f11610071578063d547741f14610490578063d5abeb01146104a3578063d73dd623146104c9578063dd62ed3e146104dc57600080fd5b8063a8fa343c14610430578063a9059cbb14610443578063c630948d14610456578063d53913931461046957600080fd5b806395d89b41116100de57806395d89b41146103fa5780639dc29fac14610402578063a217fddf14610415578063a457c2d71461041d57600080fd5b806379cc6790146103795780638fd6a6ac1461038c57806391d14854146103b457600080fd5b80632f2ff15d1161017c57806340c10f191161014b57806340c10f191461030a57806342966c681461031d578063661884631461033057806370a082311461034357600080fd5b80632f2ff15d1461029e578063313ce567146102b357806336568abe146102e457806339509351146102f757600080fd5b806318160ddd116101b857806318160ddd1461022f57806323b872dd14610241578063248a9ca314610254578063282c51f31461027757600080fd5b806301ffc9a7146101df57806306fdde0314610207578063095ea7b31461021c575b600080fd5b6101f26101ed3660046119dc565b610522565b60405190151581526020015b60405180910390f35b61020f61069f565b6040516101fe9190611a42565b6101f261022a366004611abc565b610731565b6002545b6040519081526020016101fe565b6101f261024f366004611ae6565b610749565b610233610262366004611b22565b60009081526005602052604090206001015490565b6102337f3c11d16cbaffd01df69ce1c404f6340ee057498f5f00246190ea54220576a84881565b6102b16102ac366004611b3b565b61076d565b005b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101fe565b6102b16102f2366004611b3b565b610797565b6101f2610305366004611abc565b61084f565b6102b1610318366004611abc565b61089b565b6102b161032b366004611b22565b6109e1565b6101f261033e366004611abc565b610a57565b610233610351366004611b67565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b6102b1610387366004611abc565b610a6a565b60065460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fe565b6101f26103c2366004611b3b565b600091825260056020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b61020f610ade565b6102b1610410366004611abc565b610aed565b610233600081565b6101f261042b366004611abc565b610af7565b6102b161043e366004611b67565b610bc8565b6101f2610451366004611abc565b610c4b565b6102b1610464366004611b67565b610c59565b6102337f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a681565b6102b161049e366004611b3b565b610cad565b7f0000000000000000000000000000000000000000000000000000000000000000610233565b6102b16104d7366004611abc565b610cd2565b6102336104ea366004611b82565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f36372b070000000000000000000000000000000000000000000000000000000014806105b557507fffffffff0000000000000000000000000000000000000000000000000000000082167fe6599b4d00000000000000000000000000000000000000000000000000000000145b8061060157507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b8061064d57507fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000145b8061069957507fffffffff0000000000000000000000000000000000000000000000000000000082167f8fd6a6ac00000000000000000000000000000000000000000000000000000000145b92915050565b6060600380546106ae90611bac565b80601f01602080910402602001604051908101604052809291908181526020018280546106da90611bac565b80156107275780601f106106fc57610100808354040283529160200191610727565b820191906000526020600020905b81548152906001019060200180831161070a57829003601f168201915b5050505050905090565b60003361073f818585610cdc565b5060019392505050565b600033610757858285610d10565b610762858585610de1565b506001949350505050565b60008281526005602052604090206001015461078881610e0f565b6107928383610e19565b505050565b73ffffffffffffffffffffffffffffffffffffffff81163314610841576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c66000000000000000000000000000000000060648201526084015b60405180910390fd5b61084b8282610f0d565b5050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919061073f9082908690610896908790611c2e565b610cdc565b3360009081527f15a28d26fa1bf736cf7edc9922607171ccb09c3c73b808e7772a3013e068a522602052604090205460ff16610905576040517fe2c8c9d5000000000000000000000000000000000000000000000000000000008152336004820152602401610838565b813073ffffffffffffffffffffffffffffffffffffffff82160361092857600080fd5b7f00000000000000000000000000000000000000000000000000000000000000001580159061098957507f00000000000000000000000000000000000000000000000000000000000000008261097d60025490565b6109879190611c2e565b115b156109d7578161099860025490565b6109a29190611c2e565b6040517fcbbf111300000000000000000000000000000000000000000000000000000000815260040161083891815260200190565b6107928383610fc8565b3360009081527f847f481f687befb06ed3511f1a8dcef57e83007c0147ae5047583d7056170937602052604090205460ff16610a4b576040517fc820b10b000000000000000000000000000000000000000000000000000000008152336004820152602401610838565b610a54816110bb565b50565b6000610a638383610af7565b9392505050565b3360009081527f847f481f687befb06ed3511f1a8dcef57e83007c0147ae5047583d7056170937602052604090205460ff16610ad4576040517fc820b10b000000000000000000000000000000000000000000000000000000008152336004820152602401610838565b61084b82826110c5565b6060600480546106ae90611bac565b61084b8282610a6a565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610bbb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f0000000000000000000000000000000000000000000000000000006064820152608401610838565b6107628286868403610cdc565b6000610bd381610e0f565b6006805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f9524c9e4b0b61eb018dd58a1cd856e3e74009528328ab4a613b434fa631d724290600090a3505050565b60003361073f818585610de1565b610c837f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a68261076d565b610a547f3c11d16cbaffd01df69ce1c404f6340ee057498f5f00246190ea54220576a8488261076d565b600082815260056020526040902060010154610cc881610e0f565b6107928383610f0d565b610792828261084f565b813073ffffffffffffffffffffffffffffffffffffffff821603610cff57600080fd5b610d0a8484846110da565b50505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610d0a5781811015610dd4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e63650000006044820152606401610838565b610d0a8484848403610cdc565b813073ffffffffffffffffffffffffffffffffffffffff821603610e0457600080fd5b610d0a84848461128d565b610a5481336114fc565b600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1661084b57600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff85168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610eaf3390565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff161561084b57600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b73ffffffffffffffffffffffffffffffffffffffff8216611045576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f2061646472657373006044820152606401610838565b80600260008282546110579190611c2e565b909155505073ffffffffffffffffffffffffffffffffffffffff8216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b610a5433826115b6565b6110d0823383610d10565b61084b82826115b6565b73ffffffffffffffffffffffffffffffffffffffff831661117c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f72657373000000000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff821661121f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316611330576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff82166113d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f65737300000000000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff831660009081526020819052604090205481811015611489576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e636500000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610d0a565b600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1661084b5761153c8161177a565b611547836020611799565b604051602001611558929190611c41565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f08c379a000000000000000000000000000000000000000000000000000000000825261083891600401611a42565b73ffffffffffffffffffffffffffffffffffffffff8216611659576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f73000000000000000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260409020548181101561170f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f63650000000000000000000000000000000000000000000000000000000000006064820152608401610838565b73ffffffffffffffffffffffffffffffffffffffff83166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3505050565b606061069973ffffffffffffffffffffffffffffffffffffffff831660145b606060006117a8836002611cc2565b6117b3906002611c2e565b67ffffffffffffffff8111156117cb576117cb611cd9565b6040519080825280601f01601f1916602001820160405280156117f5576020820181803683370190505b5090507f30000000000000000000000000000000000000000000000000000000000000008160008151811061182c5761182c611d08565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061188f5761188f611d08565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060006118cb846002611cc2565b6118d6906001611c2e565b90505b6001811115611973577f303132333435363738396162636465660000000000000000000000000000000085600f166010811061191757611917611d08565b1a60f81b82828151811061192d5761192d611d08565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060049490941c9361196c81611d37565b90506118d9565b508315610a63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152606401610838565b6000602082840312156119ee57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610a6357600080fd5b60005b83811015611a39578181015183820152602001611a21565b50506000910152565b6020815260008251806020840152611a61816040850160208701611a1e565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611ab757600080fd5b919050565b60008060408385031215611acf57600080fd5b611ad883611a93565b946020939093013593505050565b600080600060608486031215611afb57600080fd5b611b0484611a93565b9250611b1260208501611a93565b9150604084013590509250925092565b600060208284031215611b3457600080fd5b5035919050565b60008060408385031215611b4e57600080fd5b82359150611b5e60208401611a93565b90509250929050565b600060208284031215611b7957600080fd5b610a6382611a93565b60008060408385031215611b9557600080fd5b611b9e83611a93565b9150611b5e60208401611a93565b600181811c90821680611bc057607f821691505b602082108103611bf9577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561069957610699611bff565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000815260008351611c79816017850160208801611a1e565b7f206973206d697373696e6720726f6c65200000000000000000000000000000006017918401918201528351611cb6816028840160208801611a1e565b01602801949350505050565b808202811582820484141761069957610699611bff565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081611d4657611d46611bff565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019056fea164736f6c6343000813000a",
}

var BurnMintERC20ABI = BurnMintERC20MetaData.ABI

var BurnMintERC20Bin = BurnMintERC20MetaData.Bin

func DeployBurnMintERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals_ uint8, maxSupply_ *big.Int, preMint *big.Int) (common.Address, *types.Transaction, *BurnMintERC20, error) {
	parsed, err := BurnMintERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintERC20Bin), backend, name, symbol, decimals_, maxSupply_, preMint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintERC20{address: address, abi: *parsed, BurnMintERC20Caller: BurnMintERC20Caller{contract: contract}, BurnMintERC20Transactor: BurnMintERC20Transactor{contract: contract}, BurnMintERC20Filterer: BurnMintERC20Filterer{contract: contract}}, nil
}

type BurnMintERC20 struct {
	address common.Address
	abi     abi.ABI
	BurnMintERC20Caller
	BurnMintERC20Transactor
	BurnMintERC20Filterer
}

type BurnMintERC20Caller struct {
	contract *bind.BoundContract
}

type BurnMintERC20Transactor struct {
	contract *bind.BoundContract
}

type BurnMintERC20Filterer struct {
	contract *bind.BoundContract
}

type BurnMintERC20Session struct {
	Contract     *BurnMintERC20
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintERC20CallerSession struct {
	Contract *BurnMintERC20Caller
	CallOpts bind.CallOpts
}

type BurnMintERC20TransactorSession struct {
	Contract     *BurnMintERC20Transactor
	TransactOpts bind.TransactOpts
}

type BurnMintERC20Raw struct {
	Contract *BurnMintERC20
}

type BurnMintERC20CallerRaw struct {
	Contract *BurnMintERC20Caller
}

type BurnMintERC20TransactorRaw struct {
	Contract *BurnMintERC20Transactor
}

func NewBurnMintERC20(address common.Address, backend bind.ContractBackend) (*BurnMintERC20, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintERC20ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20{address: address, abi: abi, BurnMintERC20Caller: BurnMintERC20Caller{contract: contract}, BurnMintERC20Transactor: BurnMintERC20Transactor{contract: contract}, BurnMintERC20Filterer: BurnMintERC20Filterer{contract: contract}}, nil
}

func NewBurnMintERC20Caller(address common.Address, caller bind.ContractCaller) (*BurnMintERC20Caller, error) {
	contract, err := bindBurnMintERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20Caller{contract: contract}, nil
}

func NewBurnMintERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintERC20Transactor, error) {
	contract, err := bindBurnMintERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20Transactor{contract: contract}, nil
}

func NewBurnMintERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintERC20Filterer, error) {
	contract, err := bindBurnMintERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20Filterer{contract: contract}, nil
}

func bindBurnMintERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintERC20 *BurnMintERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintERC20.Contract.BurnMintERC20Caller.contract.Call(opts, result, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnMintERC20Transactor.contract.Transfer(opts)
}

func (_BurnMintERC20 *BurnMintERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnMintERC20Transactor.contract.Transact(opts, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintERC20.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.contract.Transfer(opts)
}

func (_BurnMintERC20 *BurnMintERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20Caller) BURNERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "BURNER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) BURNERROLE() ([32]byte, error) {
	return _BurnMintERC20.Contract.BURNERROLE(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) BURNERROLE() ([32]byte, error) {
	return _BurnMintERC20.Contract.BURNERROLE(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _BurnMintERC20.Contract.DEFAULTADMINROLE(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BurnMintERC20.Contract.DEFAULTADMINROLE(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) MINTERROLE() ([32]byte, error) {
	return _BurnMintERC20.Contract.MINTERROLE(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) MINTERROLE() ([32]byte, error) {
	return _BurnMintERC20.Contract.MINTERROLE(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.Allowance(&_BurnMintERC20.CallOpts, owner, spender)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.Allowance(&_BurnMintERC20.CallOpts, owner, spender)
}

func (_BurnMintERC20 *BurnMintERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.BalanceOf(&_BurnMintERC20.CallOpts, account)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.BalanceOf(&_BurnMintERC20.CallOpts, account)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Decimals() (uint8, error) {
	return _BurnMintERC20.Contract.Decimals(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Decimals() (uint8, error) {
	return _BurnMintERC20.Contract.Decimals(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "getCCIPAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) GetCCIPAdmin() (common.Address, error) {
	return _BurnMintERC20.Contract.GetCCIPAdmin(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) GetCCIPAdmin() (common.Address, error) {
	return _BurnMintERC20.Contract.GetCCIPAdmin(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BurnMintERC20.Contract.GetRoleAdmin(&_BurnMintERC20.CallOpts, role)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BurnMintERC20.Contract.GetRoleAdmin(&_BurnMintERC20.CallOpts, role)
}

func (_BurnMintERC20 *BurnMintERC20Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BurnMintERC20.Contract.HasRole(&_BurnMintERC20.CallOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BurnMintERC20.Contract.HasRole(&_BurnMintERC20.CallOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20Caller) MaxSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "maxSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) MaxSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.MaxSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) MaxSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.MaxSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Name() (string, error) {
	return _BurnMintERC20.Contract.Name(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Name() (string, error) {
	return _BurnMintERC20.Contract.Name(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintERC20.Contract.SupportsInterface(&_BurnMintERC20.CallOpts, interfaceId)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintERC20.Contract.SupportsInterface(&_BurnMintERC20.CallOpts, interfaceId)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Symbol() (string, error) {
	return _BurnMintERC20.Contract.Symbol(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Symbol() (string, error) {
	return _BurnMintERC20.Contract.Symbol(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) TotalSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.TotalSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.TotalSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "approve", spender, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Approve(&_BurnMintERC20.TransactOpts, spender, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Approve(&_BurnMintERC20.TransactOpts, spender, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "burn", amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Burn(amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn(&_BurnMintERC20.TransactOpts, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn(&_BurnMintERC20.TransactOpts, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "burn0", account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn0(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn0(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "burnFrom", account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnFrom(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnFrom(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseAllowance(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseAllowance(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "decreaseApproval", spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseApproval(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseApproval(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "grantMintAndBurnRoles", burnAndMinter)
}

func (_BurnMintERC20 *BurnMintERC20Session) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantMintAndBurnRoles(&_BurnMintERC20.TransactOpts, burnAndMinter)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantMintAndBurnRoles(&_BurnMintERC20.TransactOpts, burnAndMinter)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "grantRole", role, account)
}

func (_BurnMintERC20 *BurnMintERC20Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantRole(&_BurnMintERC20.TransactOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantRole(&_BurnMintERC20.TransactOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseAllowance(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseAllowance(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "increaseApproval", spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseApproval(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseApproval(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "mint", account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Mint(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Mint(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "renounceRole", role, account)
}

func (_BurnMintERC20 *BurnMintERC20Session) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RenounceRole(&_BurnMintERC20.TransactOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RenounceRole(&_BurnMintERC20.TransactOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "revokeRole", role, account)
}

func (_BurnMintERC20 *BurnMintERC20Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RevokeRole(&_BurnMintERC20.TransactOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RevokeRole(&_BurnMintERC20.TransactOpts, role, account)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "setCCIPAdmin", newAdmin)
}

func (_BurnMintERC20 *BurnMintERC20Session) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.SetCCIPAdmin(&_BurnMintERC20.TransactOpts, newAdmin)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.SetCCIPAdmin(&_BurnMintERC20.TransactOpts, newAdmin)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "transfer", to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Transfer(&_BurnMintERC20.TransactOpts, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Transfer(&_BurnMintERC20.TransactOpts, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "transferFrom", from, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.TransferFrom(&_BurnMintERC20.TransactOpts, from, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.TransferFrom(&_BurnMintERC20.TransactOpts, from, to, amount)
}

type BurnMintERC20ApprovalIterator struct {
	Event *BurnMintERC20Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20Approval)
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
		it.Event = new(BurnMintERC20Approval)
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

func (it *BurnMintERC20ApprovalIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20ApprovalIterator{contract: _BurnMintERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20Approval)
				if err := _BurnMintERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseApproval(log types.Log) (*BurnMintERC20Approval, error) {
	event := new(BurnMintERC20Approval)
	if err := _BurnMintERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20BurnAccessGrantedIterator struct {
	Event *BurnMintERC20BurnAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20BurnAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20BurnAccessGranted)
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
		it.Event = new(BurnMintERC20BurnAccessGranted)
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

func (it *BurnMintERC20BurnAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20BurnAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20BurnAccessGranted struct {
	Burner common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterBurnAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessGrantedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "BurnAccessGranted")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20BurnAccessGrantedIterator{contract: _BurnMintERC20.contract, event: "BurnAccessGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessGranted) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "BurnAccessGranted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20BurnAccessGranted)
				if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseBurnAccessGranted(log types.Log) (*BurnMintERC20BurnAccessGranted, error) {
	event := new(BurnMintERC20BurnAccessGranted)
	if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20BurnAccessRevokedIterator struct {
	Event *BurnMintERC20BurnAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20BurnAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20BurnAccessRevoked)
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
		it.Event = new(BurnMintERC20BurnAccessRevoked)
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

func (it *BurnMintERC20BurnAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20BurnAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20BurnAccessRevoked struct {
	Burner common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterBurnAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessRevokedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "BurnAccessRevoked")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20BurnAccessRevokedIterator{contract: _BurnMintERC20.contract, event: "BurnAccessRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessRevoked) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "BurnAccessRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20BurnAccessRevoked)
				if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseBurnAccessRevoked(log types.Log) (*BurnMintERC20BurnAccessRevoked, error) {
	event := new(BurnMintERC20BurnAccessRevoked)
	if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20CCIPAdminTransferredIterator struct {
	Event *BurnMintERC20CCIPAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20CCIPAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20CCIPAdminTransferred)
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
		it.Event = new(BurnMintERC20CCIPAdminTransferred)
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

func (it *BurnMintERC20CCIPAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20CCIPAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20CCIPAdminTransferred struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterCCIPAdminTransferred(opts *bind.FilterOpts, previousAdmin []common.Address, newAdmin []common.Address) (*BurnMintERC20CCIPAdminTransferredIterator, error) {

	var previousAdminRule []interface{}
	for _, previousAdminItem := range previousAdmin {
		previousAdminRule = append(previousAdminRule, previousAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "CCIPAdminTransferred", previousAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20CCIPAdminTransferredIterator{contract: _BurnMintERC20.contract, event: "CCIPAdminTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchCCIPAdminTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20CCIPAdminTransferred, previousAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var previousAdminRule []interface{}
	for _, previousAdminItem := range previousAdmin {
		previousAdminRule = append(previousAdminRule, previousAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "CCIPAdminTransferred", previousAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20CCIPAdminTransferred)
				if err := _BurnMintERC20.contract.UnpackLog(event, "CCIPAdminTransferred", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseCCIPAdminTransferred(log types.Log) (*BurnMintERC20CCIPAdminTransferred, error) {
	event := new(BurnMintERC20CCIPAdminTransferred)
	if err := _BurnMintERC20.contract.UnpackLog(event, "CCIPAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20MintAccessGrantedIterator struct {
	Event *BurnMintERC20MintAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20MintAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20MintAccessGranted)
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
		it.Event = new(BurnMintERC20MintAccessGranted)
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

func (it *BurnMintERC20MintAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20MintAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20MintAccessGranted struct {
	Minter common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterMintAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20MintAccessGrantedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "MintAccessGranted")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20MintAccessGrantedIterator{contract: _BurnMintERC20.contract, event: "MintAccessGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessGranted) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "MintAccessGranted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20MintAccessGranted)
				if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseMintAccessGranted(log types.Log) (*BurnMintERC20MintAccessGranted, error) {
	event := new(BurnMintERC20MintAccessGranted)
	if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20MintAccessRevokedIterator struct {
	Event *BurnMintERC20MintAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20MintAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20MintAccessRevoked)
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
		it.Event = new(BurnMintERC20MintAccessRevoked)
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

func (it *BurnMintERC20MintAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20MintAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20MintAccessRevoked struct {
	Minter common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterMintAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20MintAccessRevokedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "MintAccessRevoked")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20MintAccessRevokedIterator{contract: _BurnMintERC20.contract, event: "MintAccessRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessRevoked) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "MintAccessRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20MintAccessRevoked)
				if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseMintAccessRevoked(log types.Log) (*BurnMintERC20MintAccessRevoked, error) {
	event := new(BurnMintERC20MintAccessRevoked)
	if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20RoleAdminChangedIterator struct {
	Event *BurnMintERC20RoleAdminChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20RoleAdminChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20RoleAdminChanged)
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
		it.Event = new(BurnMintERC20RoleAdminChanged)
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

func (it *BurnMintERC20RoleAdminChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BurnMintERC20RoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20RoleAdminChangedIterator{contract: _BurnMintERC20.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BurnMintERC20RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20RoleAdminChanged)
				if err := _BurnMintERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseRoleAdminChanged(log types.Log) (*BurnMintERC20RoleAdminChanged, error) {
	event := new(BurnMintERC20RoleAdminChanged)
	if err := _BurnMintERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20RoleGrantedIterator struct {
	Event *BurnMintERC20RoleGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20RoleGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20RoleGranted)
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
		it.Event = new(BurnMintERC20RoleGranted)
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

func (it *BurnMintERC20RoleGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BurnMintERC20RoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20RoleGrantedIterator{contract: _BurnMintERC20.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20RoleGranted)
				if err := _BurnMintERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseRoleGranted(log types.Log) (*BurnMintERC20RoleGranted, error) {
	event := new(BurnMintERC20RoleGranted)
	if err := _BurnMintERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20RoleRevokedIterator struct {
	Event *BurnMintERC20RoleRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20RoleRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20RoleRevoked)
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
		it.Event = new(BurnMintERC20RoleRevoked)
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

func (it *BurnMintERC20RoleRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BurnMintERC20RoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20RoleRevokedIterator{contract: _BurnMintERC20.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20RoleRevoked)
				if err := _BurnMintERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseRoleRevoked(log types.Log) (*BurnMintERC20RoleRevoked, error) {
	event := new(BurnMintERC20RoleRevoked)
	if err := _BurnMintERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20TransferIterator struct {
	Event *BurnMintERC20Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20Transfer)
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
		it.Event = new(BurnMintERC20Transfer)
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

func (it *BurnMintERC20TransferIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20TransferIterator{contract: _BurnMintERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20Transfer)
				if err := _BurnMintERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseTransfer(log types.Log) (*BurnMintERC20Transfer, error) {
	event := new(BurnMintERC20Transfer)
	if err := _BurnMintERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintERC20 *BurnMintERC20) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintERC20.abi.Events["Approval"].ID:
		return _BurnMintERC20.ParseApproval(log)
	case _BurnMintERC20.abi.Events["BurnAccessGranted"].ID:
		return _BurnMintERC20.ParseBurnAccessGranted(log)
	case _BurnMintERC20.abi.Events["BurnAccessRevoked"].ID:
		return _BurnMintERC20.ParseBurnAccessRevoked(log)
	case _BurnMintERC20.abi.Events["CCIPAdminTransferred"].ID:
		return _BurnMintERC20.ParseCCIPAdminTransferred(log)
	case _BurnMintERC20.abi.Events["MintAccessGranted"].ID:
		return _BurnMintERC20.ParseMintAccessGranted(log)
	case _BurnMintERC20.abi.Events["MintAccessRevoked"].ID:
		return _BurnMintERC20.ParseMintAccessRevoked(log)
	case _BurnMintERC20.abi.Events["RoleAdminChanged"].ID:
		return _BurnMintERC20.ParseRoleAdminChanged(log)
	case _BurnMintERC20.abi.Events["RoleGranted"].ID:
		return _BurnMintERC20.ParseRoleGranted(log)
	case _BurnMintERC20.abi.Events["RoleRevoked"].ID:
		return _BurnMintERC20.ParseRoleRevoked(log)
	case _BurnMintERC20.abi.Events["Transfer"].ID:
		return _BurnMintERC20.ParseTransfer(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintERC20Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (BurnMintERC20BurnAccessGranted) Topic() common.Hash {
	return common.HexToHash("0x92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad")
}

func (BurnMintERC20BurnAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0x0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c")
}

func (BurnMintERC20CCIPAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x9524c9e4b0b61eb018dd58a1cd856e3e74009528328ab4a613b434fa631d7242")
}

func (BurnMintERC20MintAccessGranted) Topic() common.Hash {
	return common.HexToHash("0xe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea")
}

func (BurnMintERC20MintAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0xed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e9")
}

func (BurnMintERC20RoleAdminChanged) Topic() common.Hash {
	return common.HexToHash("0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff")
}

func (BurnMintERC20RoleGranted) Topic() common.Hash {
	return common.HexToHash("0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d")
}

func (BurnMintERC20RoleRevoked) Topic() common.Hash {
	return common.HexToHash("0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b")
}

func (BurnMintERC20Transfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (_BurnMintERC20 *BurnMintERC20) Address() common.Address {
	return _BurnMintERC20.address
}

type BurnMintERC20Interface interface {
	BURNERROLE(opts *bind.CallOpts) ([32]byte, error)

	DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error)

	MINTERROLE(opts *bind.CallOpts) ([32]byte, error)

	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error)

	HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error)

	MaxSupply(opts *bind.CallOpts) (*big.Int, error)

	Name(opts *bind.CallOpts) (string, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error)

	GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC20ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*BurnMintERC20Approval, error)

	FilterBurnAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessGrantedIterator, error)

	WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessGranted) (event.Subscription, error)

	ParseBurnAccessGranted(log types.Log) (*BurnMintERC20BurnAccessGranted, error)

	FilterBurnAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessRevokedIterator, error)

	WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessRevoked) (event.Subscription, error)

	ParseBurnAccessRevoked(log types.Log) (*BurnMintERC20BurnAccessRevoked, error)

	FilterCCIPAdminTransferred(opts *bind.FilterOpts, previousAdmin []common.Address, newAdmin []common.Address) (*BurnMintERC20CCIPAdminTransferredIterator, error)

	WatchCCIPAdminTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20CCIPAdminTransferred, previousAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseCCIPAdminTransferred(log types.Log) (*BurnMintERC20CCIPAdminTransferred, error)

	FilterMintAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20MintAccessGrantedIterator, error)

	WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessGranted) (event.Subscription, error)

	ParseMintAccessGranted(log types.Log) (*BurnMintERC20MintAccessGranted, error)

	FilterMintAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20MintAccessRevokedIterator, error)

	WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessRevoked) (event.Subscription, error)

	ParseMintAccessRevoked(log types.Log) (*BurnMintERC20MintAccessRevoked, error)

	FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BurnMintERC20RoleAdminChangedIterator, error)

	WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BurnMintERC20RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error)

	ParseRoleAdminChanged(log types.Log) (*BurnMintERC20RoleAdminChanged, error)

	FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BurnMintERC20RoleGrantedIterator, error)

	WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error)

	ParseRoleGranted(log types.Log) (*BurnMintERC20RoleGranted, error)

	FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BurnMintERC20RoleRevokedIterator, error)

	WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error)

	ParseRoleRevoked(log types.Log) (*BurnMintERC20RoleRevoked, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*BurnMintERC20Transfer, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
