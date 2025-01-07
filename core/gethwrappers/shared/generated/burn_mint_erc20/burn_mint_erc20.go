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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"maxSupply_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"preMint\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BURNER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MINTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCCIPAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantMintAndBurnRoles\",\"inputs\":[{\"name\":\"burnAndMinter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCCIPAdmin\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CCIPAdminTransferred\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MaxSupplyExceeded\",\"inputs\":[{\"name\":\"supplyAfterMint\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162002260380380620022608339810160408190526200003491620002e7565b848460036200004483826200040b565b5060046200005382826200040b565b50505060ff831660805260a0829052600680546001600160a01b03191633179055801562000087576200008733826200009f565b6200009460003362000166565b5050505050620004f9565b6001600160a01b038216620000fa5760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b80600260008282546200010e9190620004d7565b90915550506001600160a01b038216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35b5050565b620001728282620001f5565b620001625760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055620001ac3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b505050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff165b92915050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200024a57600080fd5b81516001600160401b038082111562000267576200026762000222565b604051601f8301601f19908116603f0116810190828211818310171562000292576200029262000222565b81604052838152602092508683858801011115620002af57600080fd5b600091505b83821015620002d35785820183015181830184015290820190620002b4565b600093810190920192909252949350505050565b600080600080600060a086880312156200030057600080fd5b85516001600160401b03808211156200031857600080fd5b6200032689838a0162000238565b965060208801519150808211156200033d57600080fd5b506200034c8882890162000238565b945050604086015160ff811681146200036457600080fd5b6060870151608090970151959894975095949392505050565b600181811c908216806200039257607f821691505b602082108103620003b357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620001f057600081815260208120601f850160051c81016020861015620003e25750805b601f850160051c820191505b818110156200040357828155600101620003ee565b505050505050565b81516001600160401b0381111562000427576200042762000222565b6200043f816200043884546200037d565b84620003b9565b602080601f8311600181146200047757600084156200045e5750858301515b600019600386901b1c1916600185901b17855562000403565b600085815260208120601f198616915b82811015620004a85788860151825594840194600190910190840162000487565b5085821015620004c75787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b808201808211156200021c57634e487b7160e01b600052601160045260246000fd5b60805160a051611d336200052d6000396000818161047c015281816108f2015261091c015260006102a40152611d336000f3fe608060405234801561001057600080fd5b50600436106101c45760003560e01c806379cc6790116100f9578063a8fa343c11610097578063d539139311610071578063d539139314610440578063d547741f14610467578063d5abeb011461047a578063dd62ed3e146104a057600080fd5b8063a8fa343c14610407578063a9059cbb1461041a578063c630948d1461042d57600080fd5b806395d89b41116100d357806395d89b41146103d15780639dc29fac146103d9578063a217fddf146103ec578063a457c2d7146103f457600080fd5b806379cc6790146103505780638fd6a6ac1461036357806391d148541461038b57600080fd5b80632f2ff15d11610166578063395093511161014057806339509351146102e157806340c10f19146102f457806342966c681461030757806370a082311461031a57600080fd5b80632f2ff15d14610288578063313ce5671461029d57806336568abe146102ce57600080fd5b806318160ddd116101a257806318160ddd1461021957806323b872dd1461022b578063248a9ca31461023e578063282c51f31461026157600080fd5b806301ffc9a7146101c957806306fdde03146101f1578063095ea7b314610206575b600080fd5b6101dc6101d7366004611996565b6104e6565b60405190151581526020015b60405180910390f35b6101f9610663565b6040516101e891906119fc565b6101dc610214366004611a76565b6106f5565b6002545b6040519081526020016101e8565b6101dc610239366004611aa0565b61070d565b61021d61024c366004611adc565b60009081526005602052604090206001015490565b61021d7f3c11d16cbaffd01df69ce1c404f6340ee057498f5f00246190ea54220576a84881565b61029b610296366004611af5565b610731565b005b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101e8565b61029b6102dc366004611af5565b61075b565b6101dc6102ef366004611a76565b610813565b61029b610302366004611a76565b61085f565b61029b610315366004611adc565b6109a9565b61021d610328366004611b21565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b61029b61035e366004611a76565b6109dc565b60065460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101e8565b6101dc610399366004611af5565b600091825260056020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b6101f9610a10565b61029b6103e7366004611a76565b610a1f565b61021d600081565b6101dc610402366004611a76565b610a29565b61029b610415366004611b21565b610afa565b6101dc610428366004611a76565b610b7d565b61029b61043b366004611b21565b610b8b565b61021d7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a681565b61029b610475366004611af5565b610be2565b7f000000000000000000000000000000000000000000000000000000000000000061021d565b61021d6104ae366004611b3c565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f36372b0700000000000000000000000000000000000000000000000000000000148061057957507fffffffff0000000000000000000000000000000000000000000000000000000082167fe6599b4d00000000000000000000000000000000000000000000000000000000145b806105c557507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b8061061157507fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000145b8061065d57507fffffffff0000000000000000000000000000000000000000000000000000000082167f8fd6a6ac00000000000000000000000000000000000000000000000000000000145b92915050565b60606003805461067290611b66565b80601f016020809104026020016040519081016040528092919081815260200182805461069e90611b66565b80156106eb5780601f106106c0576101008083540402835291602001916106eb565b820191906000526020600020905b8154815290600101906020018083116106ce57829003601f168201915b5050505050905090565b600033610703818585610c07565b5060019392505050565b60003361071b858285610c79565b610726858585610d50565b506001949350505050565b60008281526005602052604090206001015461074c81610dc2565b6107568383610dcc565b505050565b73ffffffffffffffffffffffffffffffffffffffff81163314610805576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c66000000000000000000000000000000000060648201526084015b60405180910390fd5b61080f8282610ec0565b5050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152812054909190610703908290869061085a908790611be8565b610c07565b7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a661088981610dc2565b3073ffffffffffffffffffffffffffffffffffffffff8416036108f0576040517f17858bbe00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841660048201526024016107fc565b7f00000000000000000000000000000000000000000000000000000000000000001580159061095157507f00000000000000000000000000000000000000000000000000000000000000008261094560025490565b61094f9190611be8565b115b1561099f578161096060025490565b61096a9190611be8565b6040517fcbbf11130000000000000000000000000000000000000000000000000000000081526004016107fc91815260200190565b6107568383610f7b565b7f3c11d16cbaffd01df69ce1c404f6340ee057498f5f00246190ea54220576a8486109d381610dc2565b61080f8261106e565b7f3c11d16cbaffd01df69ce1c404f6340ee057498f5f00246190ea54220576a848610a0681610dc2565b6107568383611078565b60606004805461067290611b66565b61080f82826109dc565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610aed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f00000000000000000000000000000000000000000000000000000060648201526084016107fc565b6107268286868403610c07565b6000610b0581610dc2565b6006805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f9524c9e4b0b61eb018dd58a1cd856e3e74009528328ab4a613b434fa631d724290600090a3505050565b600033610703818585610d50565b610bb57f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a682610731565b610bdf7f3c11d16cbaffd01df69ce1c404f6340ee057498f5f00246190ea54220576a84882610731565b50565b600082815260056020526040902060010154610bfd81610dc2565b6107568383610ec0565b3073ffffffffffffffffffffffffffffffffffffffff831603610c6e576040517f17858bbe00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016107fc565b61075683838361108d565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610d4a5781811015610d3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e636500000060448201526064016107fc565b610d4a8484848403610c07565b50505050565b3073ffffffffffffffffffffffffffffffffffffffff831603610db7576040517f17858bbe00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016107fc565b610756838383611240565b610bdf81336114af565b600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1661080f57600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff85168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610e623390565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff161561080f57600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b73ffffffffffffffffffffffffffffffffffffffff8216610ff8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016107fc565b806002600082825461100a9190611be8565b909155505073ffffffffffffffffffffffffffffffffffffffff8216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b610bdf3382611569565b611083823383610c79565b61080f8282611569565b73ffffffffffffffffffffffffffffffffffffffff831661112f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff82166111d2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff83166112e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f647265737300000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff8216611386576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f657373000000000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602081905260409020548181101561143c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e6365000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610d4a565b600082815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1661080f576114ef8161172d565b6114fa83602061174c565b60405160200161150b929190611bfb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f08c379a00000000000000000000000000000000000000000000000000000000082526107fc916004016119fc565b73ffffffffffffffffffffffffffffffffffffffff821661160c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f730000000000000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260208190526040902054818110156116c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f636500000000000000000000000000000000000000000000000000000000000060648201526084016107fc565b73ffffffffffffffffffffffffffffffffffffffff83166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3505050565b606061065d73ffffffffffffffffffffffffffffffffffffffff831660145b6060600061175b836002611c7c565b611766906002611be8565b67ffffffffffffffff81111561177e5761177e611c93565b6040519080825280601f01601f1916602001820160405280156117a8576020820181803683370190505b5090507f3000000000000000000000000000000000000000000000000000000000000000816000815181106117df576117df611cc2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061184257611842611cc2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600061187e846002611c7c565b611889906001611be8565b90505b6001811115611926577f303132333435363738396162636465660000000000000000000000000000000085600f16601081106118ca576118ca611cc2565b1a60f81b8282815181106118e0576118e0611cc2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060049490941c9361191f81611cf1565b905061188c565b50831561198f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016107fc565b9392505050565b6000602082840312156119a857600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461198f57600080fd5b60005b838110156119f35781810151838201526020016119db565b50506000910152565b6020815260008251806020840152611a1b8160408501602087016119d8565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611a7157600080fd5b919050565b60008060408385031215611a8957600080fd5b611a9283611a4d565b946020939093013593505050565b600080600060608486031215611ab557600080fd5b611abe84611a4d565b9250611acc60208501611a4d565b9150604084013590509250925092565b600060208284031215611aee57600080fd5b5035919050565b60008060408385031215611b0857600080fd5b82359150611b1860208401611a4d565b90509250929050565b600060208284031215611b3357600080fd5b61198f82611a4d565b60008060408385031215611b4f57600080fd5b611b5883611a4d565b9150611b1860208401611a4d565b600181811c90821680611b7a57607f821691505b602082108103611bb3577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561065d5761065d611bb9565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000815260008351611c338160178501602088016119d8565b7f206973206d697373696e6720726f6c65200000000000000000000000000000006017918401918201528351611c708160288401602088016119d8565b01602801949350505050565b808202811582820484141761065d5761065d611bb9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081611d0057611d00611bb9565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019056fea164736f6c6343000813000a",
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
	case _BurnMintERC20.abi.Events["CCIPAdminTransferred"].ID:
		return _BurnMintERC20.ParseCCIPAdminTransferred(log)
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

func (BurnMintERC20CCIPAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x9524c9e4b0b61eb018dd58a1cd856e3e74009528328ab4a613b434fa631d7242")
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

	GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error)

	GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC20ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*BurnMintERC20Approval, error)

	FilterCCIPAdminTransferred(opts *bind.FilterOpts, previousAdmin []common.Address, newAdmin []common.Address) (*BurnMintERC20CCIPAdminTransferredIterator, error)

	WatchCCIPAdminTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20CCIPAdminTransferred, previousAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseCCIPAdminTransferred(log types.Log) (*BurnMintERC20CCIPAdminTransferred, error)

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
