// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package link_token

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

var LinkTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supplyAfterMint\",\"type\":\"uint256\"}],\"name\":\"MaxSupplyExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotBurner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotMinter\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"CCIPAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"grantBurnRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burnAndMinter\",\"type\":\"address\"}],\"name\":\"grantMintAndBurnRoles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"grantMintRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"isBurner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"revokeBurnRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"revokeMintRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040518060400160405280600f81526020016e21b430b4b72634b735902a37b5b2b760891b815250604051806040016040528060048152602001634c494e4b60e01b81525060126b033b2e3c9fd0803ce80000008383838333806000868681600390816200008191906200028e565b5060046200009082826200028e565b5050506001600160a01b038216620000ef5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600580546001600160a01b0319166001600160a01b0384811691909117909155811615620001225762000122816200013d565b50505060ff90911660805260a052506200035a945050505050565b336001600160a01b03821603620001975760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000e6565b600680546001600160a01b0319166001600160a01b03838116918217909255600554604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200021457607f821691505b6020821081036200023557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200028957600081815260208120601f850160051c81016020861015620002645750805b601f850160051c820191505b81811015620002855782815560010162000270565b5050505b505050565b81516001600160401b03811115620002aa57620002aa620001e9565b620002c281620002bb8454620001ff565b846200023b565b602080601f831160018114620002fa5760008415620002e15750858301515b600019600386901b1c1916600185901b17855562000285565b600085815260208120601f198616915b828110156200032b578886015182559484019460019091019084016200030a565b50858210156200034a5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a051611f676200038e600039600081816104c50152818161086c0152610896015260006102a70152611f676000f3fe608060405234801561001057600080fd5b50600436106102265760003560e01c806386fe8b431161012a578063aa271e1a116100bd578063d5abeb011161008c578063dd62ed3e11610071578063dd62ed3e146104fc578063f2fde38b14610542578063f81094f31461055557600080fd5b8063d5abeb01146104c3578063d73dd623146104e957600080fd5b8063aa271e1a14610477578063c2e3273d1461048a578063c630948d1461049d578063c64d0ebc146104b057600080fd5b80639dc29fac116100f95780639dc29fac1461042b578063a457c2d71461043e578063a8fa343c14610451578063a9059cbb1461046457600080fd5b806386fe8b43146103be5780638da5cb5b146103c65780638fd6a6ac1461040557806395d89b411461042357600080fd5b806340c10f19116101bd578063661884631161018c57806370a082311161017157806370a082311461036d57806379ba5097146103a357806379cc6790146103ab57600080fd5b806366188463146103455780636b32810b1461035857600080fd5b806340c10f19146102f757806342966c681461030c5780634334614a1461031f5780634f5632f81461033257600080fd5b806323b872dd116101f957806323b872dd1461028d578063313ce567146102a057806339509351146102d15780634000aea0146102e457600080fd5b806301ffc9a71461022b57806306fdde0314610253578063095ea7b31461026857806318160ddd1461027b575b600080fd5b61023e610239366004611ad4565b610568565b60405190151581526020015b60405180910390f35b61025b6105c4565b60405161024a9190611b7a565b61023e610276366004611bb6565b610656565b6002545b60405190815260200161024a565b61023e61029b366004611be0565b61066e565b60405160ff7f000000000000000000000000000000000000000000000000000000000000000016815260200161024a565b61023e6102df366004611bb6565b610692565b61023e6102f2366004611c4b565b6106de565b61030a610305366004611bb6565b610801565b005b61030a61031a366004611d34565b610928565b61023e61032d366004611d4d565b610975565b61030a610340366004611d4d565b610982565b61023e610353366004611bb6565b6109de565b6103606109f1565b60405161024a9190611d68565b61027f61037b366004611d4d565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b61030a610a02565b61030a6103b9366004611bb6565b610b03565b610360610b52565b60055473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161024a565b600b5473ffffffffffffffffffffffffffffffffffffffff166103e0565b61025b610b5e565b61030a610439366004611bb6565b610b6d565b61023e61044c366004611bb6565b610b77565b61030a61045f366004611d4d565b610c48565b61023e610472366004611bb6565b610cd6565b61023e610485366004611d4d565b610ce4565b61030a610498366004611d4d565b610cf1565b61030a6104ab366004611d4d565b610d4d565b61030a6104be366004611d4d565b610d5b565b7f000000000000000000000000000000000000000000000000000000000000000061027f565b61030a6104f7366004611bb6565b610db7565b61027f61050a366004611dc2565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b61030a610550366004611d4d565b610dc1565b61030a610563366004611d4d565b610dd2565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f4000aea00000000000000000000000000000000000000000000000000000000014806105be57506105be82610e2e565b92915050565b6060600380546105d390611df5565b80601f01602080910402602001604051908101604052809291908181526020018280546105ff90611df5565b801561064c5780601f106106215761010080835404028352916020019161064c565b820191906000526020600020905b81548152906001019060200180831161062f57829003601f168201915b5050505050905090565b600033610664818585610f12565b5060019392505050565b60003361067c858285610f46565b610687858585611017565b506001949350505050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919061066490829086906106d9908790611e77565b610f12565b60006106ea8484610cd6565b508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fe19260aff97b920c7df27010903aeb9c8d2be5d310a2c67824cf3f15396e4c16858560405161074a929190611e8a565b60405180910390a373ffffffffffffffffffffffffffffffffffffffff84163b15610664576040517fa4c0ed3600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85169063a4c0ed36906107c590339087908790600401611eab565b600060405180830381600087803b1580156107df57600080fd5b505af11580156107f3573d6000803e3d6000fd5b505050505060019392505050565b61080a33610ce4565b610847576040517fe2c8c9d50000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b813073ffffffffffffffffffffffffffffffffffffffff82160361086a57600080fd5b7f0000000000000000000000000000000000000000000000000000000000000000158015906108cb57507f0000000000000000000000000000000000000000000000000000000000000000826108bf60025490565b6108c99190611e77565b115b1561091957816108da60025490565b6108e49190611e77565b6040517fcbbf111300000000000000000000000000000000000000000000000000000000815260040161083e91815260200190565b6109238383611045565b505050565b61093133610975565b610969576040517fc820b10b00000000000000000000000000000000000000000000000000000000815233600482015260240161083e565b61097281611138565b50565b60006105be600983611142565b61098a611171565b6109956009826111f4565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907f0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c90600090a250565b60006109ea8383610b77565b9392505050565b60606109fd6007611216565b905090565b60065473ffffffffffffffffffffffffffffffffffffffff163314610a83576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161083e565b600580547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560068054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b610b0c33610975565b610b44576040517fc820b10b00000000000000000000000000000000000000000000000000000000815233600482015260240161083e565b610b4e8282611223565b5050565b60606109fd6009611216565b6060600480546105d390611df5565b610b4e8282610b03565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610c3b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f000000000000000000000000000000000000000000000000000000606482015260840161083e565b6106878286868403610f12565b610c50611171565b600b805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527fe6aed30f8fcb01c95cfd4d0512b436882bb2d9005792b8d1b791bc7505614d21910160405180910390a15050565b600033610664818585611017565b60006105be600783611142565b610cf9611171565b610d04600782611238565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907fe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea90600090a250565b610d5681610cf1565b610972815b610d63611171565b610d6e600982611238565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907f92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad90600090a250565b6109238282610692565b610dc9611171565b6109728161125a565b610dda611171565b610de56007826111f4565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907fed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e990600090a250565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f36372b07000000000000000000000000000000000000000000000000000000001480610ec157507fffffffff0000000000000000000000000000000000000000000000000000000082167fe6599b4d00000000000000000000000000000000000000000000000000000000145b806105be57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b813073ffffffffffffffffffffffffffffffffffffffff821603610f3557600080fd5b610f40848484611350565b50505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610f40578181101561100a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161083e565b610f408484848403610f12565b813073ffffffffffffffffffffffffffffffffffffffff82160361103a57600080fd5b610f40848484611503565b73ffffffffffffffffffffffffffffffffffffffff82166110c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640161083e565b80600260008282546110d49190611e77565b909155505073ffffffffffffffffffffffffffffffffffffffff8216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b6109723382611772565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156109ea565b60055473ffffffffffffffffffffffffffffffffffffffff1633146111f2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161083e565b565b60006109ea8373ffffffffffffffffffffffffffffffffffffffff8416611936565b606060006109ea83611a29565b61122e823383610f46565b610b4e8282611772565b60006109ea8373ffffffffffffffffffffffffffffffffffffffff8416611a85565b3373ffffffffffffffffffffffffffffffffffffffff8216036112d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161083e565b600680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600554604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b73ffffffffffffffffffffffffffffffffffffffff83166113f2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f7265737300000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8216611495576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f7373000000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff83166115a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f6472657373000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8216611649576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260208190526040902054818110156116ff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e63650000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610f40565b73ffffffffffffffffffffffffffffffffffffffff8216611815576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f7300000000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260208190526040902054818110156118cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f6365000000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff83166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3505050565b60008181526001830160205260408120548015611a1f57600061195a600183611ee9565b855490915060009061196e90600190611ee9565b90508181146119d357600086600001828154811061198e5761198e611efc565b90600052602060002001549050808760000184815481106119b1576119b1611efc565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806119e4576119e4611f2b565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105be565b60009150506105be565b606081600001805480602002602001604051908101604052809291908181526020018280548015611a7957602002820191906000526020600020905b815481526020019060010190808311611a65575b50505050509050919050565b6000818152600183016020526040812054611acc575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105be565b5060006105be565b600060208284031215611ae657600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146109ea57600080fd5b6000815180845260005b81811015611b3c57602081850181015186830182015201611b20565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006109ea6020830184611b16565b803573ffffffffffffffffffffffffffffffffffffffff81168114611bb157600080fd5b919050565b60008060408385031215611bc957600080fd5b611bd283611b8d565b946020939093013593505050565b600080600060608486031215611bf557600080fd5b611bfe84611b8d565b9250611c0c60208501611b8d565b9150604084013590509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600060608486031215611c6057600080fd5b611c6984611b8d565b925060208401359150604084013567ffffffffffffffff80821115611c8d57600080fd5b818601915086601f830112611ca157600080fd5b813581811115611cb357611cb3611c1c565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715611cf957611cf9611c1c565b81604052828152896020848701011115611d1257600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b600060208284031215611d4657600080fd5b5035919050565b600060208284031215611d5f57600080fd5b6109ea82611b8d565b6020808252825182820181905260009190848201906040850190845b81811015611db657835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611d84565b50909695505050505050565b60008060408385031215611dd557600080fd5b611dde83611b8d565b9150611dec60208401611b8d565b90509250929050565b600181811c90821680611e0957607f821691505b602082108103611e42577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156105be576105be611e48565b828152604060208201526000611ea36040830184611b16565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff84168152826020820152606060408201526000611ee06060830184611b16565b95945050505050565b818103818111156105be576105be611e48565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var LinkTokenABI = LinkTokenMetaData.ABI

var LinkTokenBin = LinkTokenMetaData.Bin

func DeployLinkToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LinkToken, error) {
	parsed, err := LinkTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LinkTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LinkToken{address: address, abi: *parsed, LinkTokenCaller: LinkTokenCaller{contract: contract}, LinkTokenTransactor: LinkTokenTransactor{contract: contract}, LinkTokenFilterer: LinkTokenFilterer{contract: contract}}, nil
}

type LinkToken struct {
	address common.Address
	abi     abi.ABI
	LinkTokenCaller
	LinkTokenTransactor
	LinkTokenFilterer
}

type LinkTokenCaller struct {
	contract *bind.BoundContract
}

type LinkTokenTransactor struct {
	contract *bind.BoundContract
}

type LinkTokenFilterer struct {
	contract *bind.BoundContract
}

type LinkTokenSession struct {
	Contract     *LinkToken
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LinkTokenCallerSession struct {
	Contract *LinkTokenCaller
	CallOpts bind.CallOpts
}

type LinkTokenTransactorSession struct {
	Contract     *LinkTokenTransactor
	TransactOpts bind.TransactOpts
}

type LinkTokenRaw struct {
	Contract *LinkToken
}

type LinkTokenCallerRaw struct {
	Contract *LinkTokenCaller
}

type LinkTokenTransactorRaw struct {
	Contract *LinkTokenTransactor
}

func NewLinkToken(address common.Address, backend bind.ContractBackend) (*LinkToken, error) {
	abi, err := abi.JSON(strings.NewReader(LinkTokenABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLinkToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkToken{address: address, abi: abi, LinkTokenCaller: LinkTokenCaller{contract: contract}, LinkTokenTransactor: LinkTokenTransactor{contract: contract}, LinkTokenFilterer: LinkTokenFilterer{contract: contract}}, nil
}

func NewLinkTokenCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenCaller, error) {
	contract, err := bindLinkToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenCaller{contract: contract}, nil
}

func NewLinkTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenTransactor, error) {
	contract, err := bindLinkToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenTransactor{contract: contract}, nil
}

func NewLinkTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenFilterer, error) {
	contract, err := bindLinkToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenFilterer{contract: contract}, nil
}

func bindLinkToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LinkTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LinkToken *LinkTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkToken.Contract.LinkTokenCaller.contract.Call(opts, result, method, params...)
}

func (_LinkToken *LinkTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkToken.Contract.LinkTokenTransactor.contract.Transfer(opts)
}

func (_LinkToken *LinkTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkToken.Contract.LinkTokenTransactor.contract.Transact(opts, method, params...)
}

func (_LinkToken *LinkTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkToken.Contract.contract.Call(opts, result, method, params...)
}

func (_LinkToken *LinkTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkToken.Contract.contract.Transfer(opts)
}

func (_LinkToken *LinkTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkToken.Contract.contract.Transact(opts, method, params...)
}

func (_LinkToken *LinkTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LinkToken *LinkTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkToken.Contract.Allowance(&_LinkToken.CallOpts, owner, spender)
}

func (_LinkToken *LinkTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkToken.Contract.Allowance(&_LinkToken.CallOpts, owner, spender)
}

func (_LinkToken *LinkTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LinkToken *LinkTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LinkToken.Contract.BalanceOf(&_LinkToken.CallOpts, account)
}

func (_LinkToken *LinkTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LinkToken.Contract.BalanceOf(&_LinkToken.CallOpts, account)
}

func (_LinkToken *LinkTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_LinkToken *LinkTokenSession) Decimals() (uint8, error) {
	return _LinkToken.Contract.Decimals(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) Decimals() (uint8, error) {
	return _LinkToken.Contract.Decimals(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) GetBurners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "getBurners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_LinkToken *LinkTokenSession) GetBurners() ([]common.Address, error) {
	return _LinkToken.Contract.GetBurners(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) GetBurners() ([]common.Address, error) {
	return _LinkToken.Contract.GetBurners(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "getCCIPAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LinkToken *LinkTokenSession) GetCCIPAdmin() (common.Address, error) {
	return _LinkToken.Contract.GetCCIPAdmin(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) GetCCIPAdmin() (common.Address, error) {
	return _LinkToken.Contract.GetCCIPAdmin(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) GetMinters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "getMinters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_LinkToken *LinkTokenSession) GetMinters() ([]common.Address, error) {
	return _LinkToken.Contract.GetMinters(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) GetMinters() ([]common.Address, error) {
	return _LinkToken.Contract.GetMinters(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) IsBurner(opts *bind.CallOpts, burner common.Address) (bool, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "isBurner", burner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LinkToken *LinkTokenSession) IsBurner(burner common.Address) (bool, error) {
	return _LinkToken.Contract.IsBurner(&_LinkToken.CallOpts, burner)
}

func (_LinkToken *LinkTokenCallerSession) IsBurner(burner common.Address) (bool, error) {
	return _LinkToken.Contract.IsBurner(&_LinkToken.CallOpts, burner)
}

func (_LinkToken *LinkTokenCaller) IsMinter(opts *bind.CallOpts, minter common.Address) (bool, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "isMinter", minter)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LinkToken *LinkTokenSession) IsMinter(minter common.Address) (bool, error) {
	return _LinkToken.Contract.IsMinter(&_LinkToken.CallOpts, minter)
}

func (_LinkToken *LinkTokenCallerSession) IsMinter(minter common.Address) (bool, error) {
	return _LinkToken.Contract.IsMinter(&_LinkToken.CallOpts, minter)
}

func (_LinkToken *LinkTokenCaller) MaxSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "maxSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LinkToken *LinkTokenSession) MaxSupply() (*big.Int, error) {
	return _LinkToken.Contract.MaxSupply(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) MaxSupply() (*big.Int, error) {
	return _LinkToken.Contract.MaxSupply(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LinkToken *LinkTokenSession) Name() (string, error) {
	return _LinkToken.Contract.Name(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) Name() (string, error) {
	return _LinkToken.Contract.Name(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LinkToken *LinkTokenSession) Owner() (common.Address, error) {
	return _LinkToken.Contract.Owner(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) Owner() (common.Address, error) {
	return _LinkToken.Contract.Owner(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LinkToken *LinkTokenSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LinkToken.Contract.SupportsInterface(&_LinkToken.CallOpts, interfaceId)
}

func (_LinkToken *LinkTokenCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LinkToken.Contract.SupportsInterface(&_LinkToken.CallOpts, interfaceId)
}

func (_LinkToken *LinkTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LinkToken *LinkTokenSession) Symbol() (string, error) {
	return _LinkToken.Contract.Symbol(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) Symbol() (string, error) {
	return _LinkToken.Contract.Symbol(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LinkToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LinkToken *LinkTokenSession) TotalSupply() (*big.Int, error) {
	return _LinkToken.Contract.TotalSupply(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _LinkToken.Contract.TotalSupply(&_LinkToken.CallOpts)
}

func (_LinkToken *LinkTokenTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "acceptOwnership")
}

func (_LinkToken *LinkTokenSession) AcceptOwnership() (*types.Transaction, error) {
	return _LinkToken.Contract.AcceptOwnership(&_LinkToken.TransactOpts)
}

func (_LinkToken *LinkTokenTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LinkToken.Contract.AcceptOwnership(&_LinkToken.TransactOpts)
}

func (_LinkToken *LinkTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "approve", spender, amount)
}

func (_LinkToken *LinkTokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Approve(&_LinkToken.TransactOpts, spender, amount)
}

func (_LinkToken *LinkTokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Approve(&_LinkToken.TransactOpts, spender, amount)
}

func (_LinkToken *LinkTokenTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "burn", amount)
}

func (_LinkToken *LinkTokenSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Burn(&_LinkToken.TransactOpts, amount)
}

func (_LinkToken *LinkTokenTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Burn(&_LinkToken.TransactOpts, amount)
}

func (_LinkToken *LinkTokenTransactor) Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "burn0", account, amount)
}

func (_LinkToken *LinkTokenSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Burn0(&_LinkToken.TransactOpts, account, amount)
}

func (_LinkToken *LinkTokenTransactorSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Burn0(&_LinkToken.TransactOpts, account, amount)
}

func (_LinkToken *LinkTokenTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "burnFrom", account, amount)
}

func (_LinkToken *LinkTokenSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.BurnFrom(&_LinkToken.TransactOpts, account, amount)
}

func (_LinkToken *LinkTokenTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.BurnFrom(&_LinkToken.TransactOpts, account, amount)
}

func (_LinkToken *LinkTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_LinkToken *LinkTokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.DecreaseAllowance(&_LinkToken.TransactOpts, spender, subtractedValue)
}

func (_LinkToken *LinkTokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.DecreaseAllowance(&_LinkToken.TransactOpts, spender, subtractedValue)
}

func (_LinkToken *LinkTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "decreaseApproval", spender, subtractedValue)
}

func (_LinkToken *LinkTokenSession) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.DecreaseApproval(&_LinkToken.TransactOpts, spender, subtractedValue)
}

func (_LinkToken *LinkTokenTransactorSession) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.DecreaseApproval(&_LinkToken.TransactOpts, spender, subtractedValue)
}

func (_LinkToken *LinkTokenTransactor) GrantBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "grantBurnRole", burner)
}

func (_LinkToken *LinkTokenSession) GrantBurnRole(burner common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.GrantBurnRole(&_LinkToken.TransactOpts, burner)
}

func (_LinkToken *LinkTokenTransactorSession) GrantBurnRole(burner common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.GrantBurnRole(&_LinkToken.TransactOpts, burner)
}

func (_LinkToken *LinkTokenTransactor) GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "grantMintAndBurnRoles", burnAndMinter)
}

func (_LinkToken *LinkTokenSession) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.GrantMintAndBurnRoles(&_LinkToken.TransactOpts, burnAndMinter)
}

func (_LinkToken *LinkTokenTransactorSession) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.GrantMintAndBurnRoles(&_LinkToken.TransactOpts, burnAndMinter)
}

func (_LinkToken *LinkTokenTransactor) GrantMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "grantMintRole", minter)
}

func (_LinkToken *LinkTokenSession) GrantMintRole(minter common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.GrantMintRole(&_LinkToken.TransactOpts, minter)
}

func (_LinkToken *LinkTokenTransactorSession) GrantMintRole(minter common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.GrantMintRole(&_LinkToken.TransactOpts, minter)
}

func (_LinkToken *LinkTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_LinkToken *LinkTokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.IncreaseAllowance(&_LinkToken.TransactOpts, spender, addedValue)
}

func (_LinkToken *LinkTokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.IncreaseAllowance(&_LinkToken.TransactOpts, spender, addedValue)
}

func (_LinkToken *LinkTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "increaseApproval", spender, addedValue)
}

func (_LinkToken *LinkTokenSession) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.IncreaseApproval(&_LinkToken.TransactOpts, spender, addedValue)
}

func (_LinkToken *LinkTokenTransactorSession) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.IncreaseApproval(&_LinkToken.TransactOpts, spender, addedValue)
}

func (_LinkToken *LinkTokenTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "mint", account, amount)
}

func (_LinkToken *LinkTokenSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Mint(&_LinkToken.TransactOpts, account, amount)
}

func (_LinkToken *LinkTokenTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Mint(&_LinkToken.TransactOpts, account, amount)
}

func (_LinkToken *LinkTokenTransactor) RevokeBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "revokeBurnRole", burner)
}

func (_LinkToken *LinkTokenSession) RevokeBurnRole(burner common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.RevokeBurnRole(&_LinkToken.TransactOpts, burner)
}

func (_LinkToken *LinkTokenTransactorSession) RevokeBurnRole(burner common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.RevokeBurnRole(&_LinkToken.TransactOpts, burner)
}

func (_LinkToken *LinkTokenTransactor) RevokeMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "revokeMintRole", minter)
}

func (_LinkToken *LinkTokenSession) RevokeMintRole(minter common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.RevokeMintRole(&_LinkToken.TransactOpts, minter)
}

func (_LinkToken *LinkTokenTransactorSession) RevokeMintRole(minter common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.RevokeMintRole(&_LinkToken.TransactOpts, minter)
}

func (_LinkToken *LinkTokenTransactor) SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "setCCIPAdmin", newAdmin)
}

func (_LinkToken *LinkTokenSession) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.SetCCIPAdmin(&_LinkToken.TransactOpts, newAdmin)
}

func (_LinkToken *LinkTokenTransactorSession) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.SetCCIPAdmin(&_LinkToken.TransactOpts, newAdmin)
}

func (_LinkToken *LinkTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "transfer", to, amount)
}

func (_LinkToken *LinkTokenSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Transfer(&_LinkToken.TransactOpts, to, amount)
}

func (_LinkToken *LinkTokenTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.Transfer(&_LinkToken.TransactOpts, to, amount)
}

func (_LinkToken *LinkTokenTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "transferAndCall", to, amount, data)
}

func (_LinkToken *LinkTokenSession) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkToken.Contract.TransferAndCall(&_LinkToken.TransactOpts, to, amount, data)
}

func (_LinkToken *LinkTokenTransactorSession) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkToken.Contract.TransferAndCall(&_LinkToken.TransactOpts, to, amount, data)
}

func (_LinkToken *LinkTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "transferFrom", from, to, amount)
}

func (_LinkToken *LinkTokenSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.TransferFrom(&_LinkToken.TransactOpts, from, to, amount)
}

func (_LinkToken *LinkTokenTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LinkToken.Contract.TransferFrom(&_LinkToken.TransactOpts, from, to, amount)
}

func (_LinkToken *LinkTokenTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LinkToken.contract.Transact(opts, "transferOwnership", to)
}

func (_LinkToken *LinkTokenSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.TransferOwnership(&_LinkToken.TransactOpts, to)
}

func (_LinkToken *LinkTokenTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LinkToken.Contract.TransferOwnership(&_LinkToken.TransactOpts, to)
}

type LinkTokenApprovalIterator struct {
	Event *LinkTokenApproval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenApproval)
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
		it.Event = new(LinkTokenApproval)
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

func (it *LinkTokenApprovalIterator) Error() error {
	return it.fail
}

func (it *LinkTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LinkTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenApprovalIterator{contract: _LinkToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LinkTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenApproval)
				if err := _LinkToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseApproval(log types.Log) (*LinkTokenApproval, error) {
	event := new(LinkTokenApproval)
	if err := _LinkToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenBurnAccessGrantedIterator struct {
	Event *LinkTokenBurnAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenBurnAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenBurnAccessGranted)
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
		it.Event = new(LinkTokenBurnAccessGranted)
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

func (it *LinkTokenBurnAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *LinkTokenBurnAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenBurnAccessGranted struct {
	Burner common.Address
	Raw    types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterBurnAccessGranted(opts *bind.FilterOpts, burner []common.Address) (*LinkTokenBurnAccessGrantedIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "BurnAccessGranted", burnerRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenBurnAccessGrantedIterator{contract: _LinkToken.contract, event: "BurnAccessGranted", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *LinkTokenBurnAccessGranted, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "BurnAccessGranted", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenBurnAccessGranted)
				if err := _LinkToken.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseBurnAccessGranted(log types.Log) (*LinkTokenBurnAccessGranted, error) {
	event := new(LinkTokenBurnAccessGranted)
	if err := _LinkToken.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenBurnAccessRevokedIterator struct {
	Event *LinkTokenBurnAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenBurnAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenBurnAccessRevoked)
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
		it.Event = new(LinkTokenBurnAccessRevoked)
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

func (it *LinkTokenBurnAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *LinkTokenBurnAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenBurnAccessRevoked struct {
	Burner common.Address
	Raw    types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterBurnAccessRevoked(opts *bind.FilterOpts, burner []common.Address) (*LinkTokenBurnAccessRevokedIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "BurnAccessRevoked", burnerRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenBurnAccessRevokedIterator{contract: _LinkToken.contract, event: "BurnAccessRevoked", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *LinkTokenBurnAccessRevoked, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "BurnAccessRevoked", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenBurnAccessRevoked)
				if err := _LinkToken.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseBurnAccessRevoked(log types.Log) (*LinkTokenBurnAccessRevoked, error) {
	event := new(LinkTokenBurnAccessRevoked)
	if err := _LinkToken.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenCCIPAdminUpdatedIterator struct {
	Event *LinkTokenCCIPAdminUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenCCIPAdminUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenCCIPAdminUpdated)
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
		it.Event = new(LinkTokenCCIPAdminUpdated)
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

func (it *LinkTokenCCIPAdminUpdatedIterator) Error() error {
	return it.fail
}

func (it *LinkTokenCCIPAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenCCIPAdminUpdated struct {
	OldAdmin common.Address
	NewAdmin common.Address
	Raw      types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterCCIPAdminUpdated(opts *bind.FilterOpts) (*LinkTokenCCIPAdminUpdatedIterator, error) {

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "CCIPAdminUpdated")
	if err != nil {
		return nil, err
	}
	return &LinkTokenCCIPAdminUpdatedIterator{contract: _LinkToken.contract, event: "CCIPAdminUpdated", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchCCIPAdminUpdated(opts *bind.WatchOpts, sink chan<- *LinkTokenCCIPAdminUpdated) (event.Subscription, error) {

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "CCIPAdminUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenCCIPAdminUpdated)
				if err := _LinkToken.contract.UnpackLog(event, "CCIPAdminUpdated", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseCCIPAdminUpdated(log types.Log) (*LinkTokenCCIPAdminUpdated, error) {
	event := new(LinkTokenCCIPAdminUpdated)
	if err := _LinkToken.contract.UnpackLog(event, "CCIPAdminUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenMintAccessGrantedIterator struct {
	Event *LinkTokenMintAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenMintAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenMintAccessGranted)
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
		it.Event = new(LinkTokenMintAccessGranted)
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

func (it *LinkTokenMintAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *LinkTokenMintAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenMintAccessGranted struct {
	Minter common.Address
	Raw    types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterMintAccessGranted(opts *bind.FilterOpts, minter []common.Address) (*LinkTokenMintAccessGrantedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "MintAccessGranted", minterRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenMintAccessGrantedIterator{contract: _LinkToken.contract, event: "MintAccessGranted", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *LinkTokenMintAccessGranted, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "MintAccessGranted", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenMintAccessGranted)
				if err := _LinkToken.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseMintAccessGranted(log types.Log) (*LinkTokenMintAccessGranted, error) {
	event := new(LinkTokenMintAccessGranted)
	if err := _LinkToken.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenMintAccessRevokedIterator struct {
	Event *LinkTokenMintAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenMintAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenMintAccessRevoked)
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
		it.Event = new(LinkTokenMintAccessRevoked)
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

func (it *LinkTokenMintAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *LinkTokenMintAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenMintAccessRevoked struct {
	Minter common.Address
	Raw    types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterMintAccessRevoked(opts *bind.FilterOpts, minter []common.Address) (*LinkTokenMintAccessRevokedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "MintAccessRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenMintAccessRevokedIterator{contract: _LinkToken.contract, event: "MintAccessRevoked", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *LinkTokenMintAccessRevoked, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "MintAccessRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenMintAccessRevoked)
				if err := _LinkToken.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseMintAccessRevoked(log types.Log) (*LinkTokenMintAccessRevoked, error) {
	event := new(LinkTokenMintAccessRevoked)
	if err := _LinkToken.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenOwnershipTransferRequestedIterator struct {
	Event *LinkTokenOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenOwnershipTransferRequested)
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
		it.Event = new(LinkTokenOwnershipTransferRequested)
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

func (it *LinkTokenOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LinkTokenOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenOwnershipTransferRequestedIterator{contract: _LinkToken.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LinkTokenOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenOwnershipTransferRequested)
				if err := _LinkToken.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseOwnershipTransferRequested(log types.Log) (*LinkTokenOwnershipTransferRequested, error) {
	event := new(LinkTokenOwnershipTransferRequested)
	if err := _LinkToken.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenOwnershipTransferredIterator struct {
	Event *LinkTokenOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenOwnershipTransferred)
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
		it.Event = new(LinkTokenOwnershipTransferred)
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

func (it *LinkTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LinkTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenOwnershipTransferredIterator{contract: _LinkToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LinkTokenOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenOwnershipTransferred)
				if err := _LinkToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseOwnershipTransferred(log types.Log) (*LinkTokenOwnershipTransferred, error) {
	event := new(LinkTokenOwnershipTransferred)
	if err := _LinkToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenTransferIterator struct {
	Event *LinkTokenTransfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenTransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenTransfer)
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
		it.Event = new(LinkTokenTransfer)
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

func (it *LinkTokenTransferIterator) Error() error {
	return it.fail
}

func (it *LinkTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Data  []byte
	Raw   types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenTransferIterator{contract: _LinkToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LinkTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenTransfer)
				if err := _LinkToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseTransfer(log types.Log) (*LinkTokenTransfer, error) {
	event := new(LinkTokenTransfer)
	if err := _LinkToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LinkTokenTransfer0Iterator struct {
	Event *LinkTokenTransfer0

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LinkTokenTransfer0Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LinkTokenTransfer0)
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
		it.Event = new(LinkTokenTransfer0)
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

func (it *LinkTokenTransfer0Iterator) Error() error {
	return it.fail
}

func (it *LinkTokenTransfer0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LinkTokenTransfer0 struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_LinkToken *LinkTokenFilterer) FilterTransfer0(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenTransfer0Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.FilterLogs(opts, "Transfer0", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LinkTokenTransfer0Iterator{contract: _LinkToken.contract, event: "Transfer0", logs: logs, sub: sub}, nil
}

func (_LinkToken *LinkTokenFilterer) WatchTransfer0(opts *bind.WatchOpts, sink chan<- *LinkTokenTransfer0, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LinkToken.contract.WatchLogs(opts, "Transfer0", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LinkTokenTransfer0)
				if err := _LinkToken.contract.UnpackLog(event, "Transfer0", log); err != nil {
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

func (_LinkToken *LinkTokenFilterer) ParseTransfer0(log types.Log) (*LinkTokenTransfer0, error) {
	event := new(LinkTokenTransfer0)
	if err := _LinkToken.contract.UnpackLog(event, "Transfer0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LinkToken *LinkToken) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LinkToken.abi.Events["Approval"].ID:
		return _LinkToken.ParseApproval(log)
	case _LinkToken.abi.Events["BurnAccessGranted"].ID:
		return _LinkToken.ParseBurnAccessGranted(log)
	case _LinkToken.abi.Events["BurnAccessRevoked"].ID:
		return _LinkToken.ParseBurnAccessRevoked(log)
	case _LinkToken.abi.Events["CCIPAdminUpdated"].ID:
		return _LinkToken.ParseCCIPAdminUpdated(log)
	case _LinkToken.abi.Events["MintAccessGranted"].ID:
		return _LinkToken.ParseMintAccessGranted(log)
	case _LinkToken.abi.Events["MintAccessRevoked"].ID:
		return _LinkToken.ParseMintAccessRevoked(log)
	case _LinkToken.abi.Events["OwnershipTransferRequested"].ID:
		return _LinkToken.ParseOwnershipTransferRequested(log)
	case _LinkToken.abi.Events["OwnershipTransferred"].ID:
		return _LinkToken.ParseOwnershipTransferred(log)
	case _LinkToken.abi.Events["Transfer"].ID:
		return _LinkToken.ParseTransfer(log)
	case _LinkToken.abi.Events["Transfer0"].ID:
		return _LinkToken.ParseTransfer0(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LinkTokenApproval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (LinkTokenBurnAccessGranted) Topic() common.Hash {
	return common.HexToHash("0x92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad")
}

func (LinkTokenBurnAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0x0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c")
}

func (LinkTokenCCIPAdminUpdated) Topic() common.Hash {
	return common.HexToHash("0xe6aed30f8fcb01c95cfd4d0512b436882bb2d9005792b8d1b791bc7505614d21")
}

func (LinkTokenMintAccessGranted) Topic() common.Hash {
	return common.HexToHash("0xe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea")
}

func (LinkTokenMintAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0xed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e9")
}

func (LinkTokenOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LinkTokenOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LinkTokenTransfer) Topic() common.Hash {
	return common.HexToHash("0xe19260aff97b920c7df27010903aeb9c8d2be5d310a2c67824cf3f15396e4c16")
}

func (LinkTokenTransfer0) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (_LinkToken *LinkToken) Address() common.Address {
	return _LinkToken.address
}

type LinkTokenInterface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	GetBurners(opts *bind.CallOpts) ([]common.Address, error)

	GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error)

	GetMinters(opts *bind.CallOpts) ([]common.Address, error)

	IsBurner(opts *bind.CallOpts, burner common.Address) (bool, error)

	IsMinter(opts *bind.CallOpts, minter common.Address) (bool, error)

	MaxSupply(opts *bind.CallOpts) (*big.Int, error)

	Name(opts *bind.CallOpts) (string, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	GrantBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error)

	GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error)

	GrantMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	RevokeBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error)

	RevokeMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error)

	SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferAndCall(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LinkTokenApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *LinkTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*LinkTokenApproval, error)

	FilterBurnAccessGranted(opts *bind.FilterOpts, burner []common.Address) (*LinkTokenBurnAccessGrantedIterator, error)

	WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *LinkTokenBurnAccessGranted, burner []common.Address) (event.Subscription, error)

	ParseBurnAccessGranted(log types.Log) (*LinkTokenBurnAccessGranted, error)

	FilterBurnAccessRevoked(opts *bind.FilterOpts, burner []common.Address) (*LinkTokenBurnAccessRevokedIterator, error)

	WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *LinkTokenBurnAccessRevoked, burner []common.Address) (event.Subscription, error)

	ParseBurnAccessRevoked(log types.Log) (*LinkTokenBurnAccessRevoked, error)

	FilterCCIPAdminUpdated(opts *bind.FilterOpts) (*LinkTokenCCIPAdminUpdatedIterator, error)

	WatchCCIPAdminUpdated(opts *bind.WatchOpts, sink chan<- *LinkTokenCCIPAdminUpdated) (event.Subscription, error)

	ParseCCIPAdminUpdated(log types.Log) (*LinkTokenCCIPAdminUpdated, error)

	FilterMintAccessGranted(opts *bind.FilterOpts, minter []common.Address) (*LinkTokenMintAccessGrantedIterator, error)

	WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *LinkTokenMintAccessGranted, minter []common.Address) (event.Subscription, error)

	ParseMintAccessGranted(log types.Log) (*LinkTokenMintAccessGranted, error)

	FilterMintAccessRevoked(opts *bind.FilterOpts, minter []common.Address) (*LinkTokenMintAccessRevokedIterator, error)

	WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *LinkTokenMintAccessRevoked, minter []common.Address) (event.Subscription, error)

	ParseMintAccessRevoked(log types.Log) (*LinkTokenMintAccessRevoked, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LinkTokenOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LinkTokenOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LinkTokenOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LinkTokenOwnershipTransferred, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenTransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *LinkTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*LinkTokenTransfer, error)

	FilterTransfer0(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LinkTokenTransfer0Iterator, error)

	WatchTransfer0(opts *bind.WatchOpts, sink chan<- *LinkTokenTransfer0, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer0(log types.Log) (*LinkTokenTransfer0, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
