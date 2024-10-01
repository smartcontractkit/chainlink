// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_erc677

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

var BurnMintERC677MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"maxSupply_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supplyAfterMint\",\"type\":\"uint256\"}],\"name\":\"MaxSupplyExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotBurner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotMinter\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"CCIPAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"grantBurnRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burnAndMinter\",\"type\":\"address\"}],\"name\":\"grantMintAndBurnRoles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"grantMintRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"isBurner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"revokeBurnRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"revokeMintRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620023fa380380620023fa833981016040819052620000349162000279565b8383838333806000868660036200004c838262000393565b5060046200005b828262000393565b5050506001600160a01b038216620000ba5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600580546001600160a01b0319166001600160a01b0384811691909117909155811615620000ed57620000ed8162000108565b50505060ff90911660805260a052506200045f945050505050565b336001600160a01b03821603620001625760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000b1565b600680546001600160a01b0319166001600160a01b03838116918217909255600554604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620001dc57600080fd5b81516001600160401b0380821115620001f957620001f9620001b4565b604051601f8301601f19908116603f01168101908282118183101715620002245762000224620001b4565b816040528381526020925086838588010111156200024157600080fd5b600091505b8382101562000265578582018301518183018401529082019062000246565b600093810190920192909252949350505050565b600080600080608085870312156200029057600080fd5b84516001600160401b0380821115620002a857600080fd5b620002b688838901620001ca565b95506020870151915080821115620002cd57600080fd5b50620002dc87828801620001ca565b935050604085015160ff81168114620002f457600080fd5b6060959095015193969295505050565b600181811c908216806200031957607f821691505b6020821081036200033a57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200038e57600081815260208120601f850160051c81016020861015620003695750805b601f850160051c820191505b818110156200038a5782815560010162000375565b5050505b505050565b81516001600160401b03811115620003af57620003af620001b4565b620003c781620003c0845462000304565b8462000340565b602080601f831160018114620003ff5760008415620003e65750858301515b600019600386901b1c1916600185901b1785556200038a565b600085815260208120601f198616915b8281101562000430578886015182559484019460019091019084016200040f565b50858210156200044f5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a051611f6762000493600039600081816104c50152818161086c0152610896015260006102a70152611f676000f3fe608060405234801561001057600080fd5b50600436106102265760003560e01c806386fe8b431161012a578063aa271e1a116100bd578063d5abeb011161008c578063dd62ed3e11610071578063dd62ed3e146104fc578063f2fde38b14610542578063f81094f31461055557600080fd5b8063d5abeb01146104c3578063d73dd623146104e957600080fd5b8063aa271e1a14610477578063c2e3273d1461048a578063c630948d1461049d578063c64d0ebc146104b057600080fd5b80639dc29fac116100f95780639dc29fac1461042b578063a457c2d71461043e578063a8fa343c14610451578063a9059cbb1461046457600080fd5b806386fe8b43146103be5780638da5cb5b146103c65780638fd6a6ac1461040557806395d89b411461042357600080fd5b806340c10f19116101bd578063661884631161018c57806370a082311161017157806370a082311461036d57806379ba5097146103a357806379cc6790146103ab57600080fd5b806366188463146103455780636b32810b1461035857600080fd5b806340c10f19146102f757806342966c681461030c5780634334614a1461031f5780634f5632f81461033257600080fd5b806323b872dd116101f957806323b872dd1461028d578063313ce567146102a057806339509351146102d15780634000aea0146102e457600080fd5b806301ffc9a71461022b57806306fdde0314610253578063095ea7b31461026857806318160ddd1461027b575b600080fd5b61023e610239366004611ad4565b610568565b60405190151581526020015b60405180910390f35b61025b6105c4565b60405161024a9190611b7a565b61023e610276366004611bb6565b610656565b6002545b60405190815260200161024a565b61023e61029b366004611be0565b61066e565b60405160ff7f000000000000000000000000000000000000000000000000000000000000000016815260200161024a565b61023e6102df366004611bb6565b610692565b61023e6102f2366004611c4b565b6106de565b61030a610305366004611bb6565b610801565b005b61030a61031a366004611d34565b610928565b61023e61032d366004611d4d565b610975565b61030a610340366004611d4d565b610982565b61023e610353366004611bb6565b6109de565b6103606109f1565b60405161024a9190611d68565b61027f61037b366004611d4d565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b61030a610a02565b61030a6103b9366004611bb6565b610b03565b610360610b52565b60055473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161024a565b600b5473ffffffffffffffffffffffffffffffffffffffff166103e0565b61025b610b5e565b61030a610439366004611bb6565b610b6d565b61023e61044c366004611bb6565b610b77565b61030a61045f366004611d4d565b610c48565b61023e610472366004611bb6565b610cd6565b61023e610485366004611d4d565b610ce4565b61030a610498366004611d4d565b610cf1565b61030a6104ab366004611d4d565b610d4d565b61030a6104be366004611d4d565b610d5b565b7f000000000000000000000000000000000000000000000000000000000000000061027f565b61030a6104f7366004611bb6565b610db7565b61027f61050a366004611dc2565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b61030a610550366004611d4d565b610dc1565b61030a610563366004611d4d565b610dd2565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f4000aea00000000000000000000000000000000000000000000000000000000014806105be57506105be82610e2e565b92915050565b6060600380546105d390611df5565b80601f01602080910402602001604051908101604052809291908181526020018280546105ff90611df5565b801561064c5780601f106106215761010080835404028352916020019161064c565b820191906000526020600020905b81548152906001019060200180831161062f57829003601f168201915b5050505050905090565b600033610664818585610f12565b5060019392505050565b60003361067c858285610f46565b610687858585611017565b506001949350505050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919061066490829086906106d9908790611e77565b610f12565b60006106ea8484610cd6565b508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fe19260aff97b920c7df27010903aeb9c8d2be5d310a2c67824cf3f15396e4c16858560405161074a929190611e8a565b60405180910390a373ffffffffffffffffffffffffffffffffffffffff84163b15610664576040517fa4c0ed3600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85169063a4c0ed36906107c590339087908790600401611eab565b600060405180830381600087803b1580156107df57600080fd5b505af11580156107f3573d6000803e3d6000fd5b505050505060019392505050565b61080a33610ce4565b610847576040517fe2c8c9d50000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b813073ffffffffffffffffffffffffffffffffffffffff82160361086a57600080fd5b7f0000000000000000000000000000000000000000000000000000000000000000158015906108cb57507f0000000000000000000000000000000000000000000000000000000000000000826108bf60025490565b6108c99190611e77565b115b1561091957816108da60025490565b6108e49190611e77565b6040517fcbbf111300000000000000000000000000000000000000000000000000000000815260040161083e91815260200190565b6109238383611045565b505050565b61093133610975565b610969576040517fc820b10b00000000000000000000000000000000000000000000000000000000815233600482015260240161083e565b61097281611138565b50565b60006105be600983611142565b61098a611171565b6109956009826111f4565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907f0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c90600090a250565b60006109ea8383610b77565b9392505050565b60606109fd6007611216565b905090565b60065473ffffffffffffffffffffffffffffffffffffffff163314610a83576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161083e565b600580547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560068054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b610b0c33610975565b610b44576040517fc820b10b00000000000000000000000000000000000000000000000000000000815233600482015260240161083e565b610b4e8282611223565b5050565b60606109fd6009611216565b6060600480546105d390611df5565b610b4e8282610b03565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610c3b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f000000000000000000000000000000000000000000000000000000606482015260840161083e565b6106878286868403610f12565b610c50611171565b600b805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527fe6aed30f8fcb01c95cfd4d0512b436882bb2d9005792b8d1b791bc7505614d21910160405180910390a15050565b600033610664818585611017565b60006105be600783611142565b610cf9611171565b610d04600782611238565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907fe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea90600090a250565b610d5681610cf1565b610972815b610d63611171565b610d6e600982611238565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907f92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad90600090a250565b6109238282610692565b610dc9611171565b6109728161125a565b610dda611171565b610de56007826111f4565b156109725760405173ffffffffffffffffffffffffffffffffffffffff8216907fed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e990600090a250565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f36372b07000000000000000000000000000000000000000000000000000000001480610ec157507fffffffff0000000000000000000000000000000000000000000000000000000082167fe6599b4d00000000000000000000000000000000000000000000000000000000145b806105be57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b813073ffffffffffffffffffffffffffffffffffffffff821603610f3557600080fd5b610f40848484611350565b50505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610f40578181101561100a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161083e565b610f408484848403610f12565b813073ffffffffffffffffffffffffffffffffffffffff82160361103a57600080fd5b610f40848484611503565b73ffffffffffffffffffffffffffffffffffffffff82166110c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640161083e565b80600260008282546110d49190611e77565b909155505073ffffffffffffffffffffffffffffffffffffffff8216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b6109723382611772565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156109ea565b60055473ffffffffffffffffffffffffffffffffffffffff1633146111f2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161083e565b565b60006109ea8373ffffffffffffffffffffffffffffffffffffffff8416611936565b606060006109ea83611a29565b61122e823383610f46565b610b4e8282611772565b60006109ea8373ffffffffffffffffffffffffffffffffffffffff8416611a85565b3373ffffffffffffffffffffffffffffffffffffffff8216036112d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161083e565b600680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600554604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b73ffffffffffffffffffffffffffffffffffffffff83166113f2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f7265737300000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8216611495576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f7373000000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff83166115a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f6472657373000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8216611649576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260208190526040902054818110156116ff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e63650000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610f40565b73ffffffffffffffffffffffffffffffffffffffff8216611815576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f7300000000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260208190526040902054818110156118cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f6365000000000000000000000000000000000000000000000000000000000000606482015260840161083e565b73ffffffffffffffffffffffffffffffffffffffff83166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3505050565b60008181526001830160205260408120548015611a1f57600061195a600183611ee9565b855490915060009061196e90600190611ee9565b90508181146119d357600086600001828154811061198e5761198e611efc565b90600052602060002001549050808760000184815481106119b1576119b1611efc565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806119e4576119e4611f2b565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105be565b60009150506105be565b606081600001805480602002602001604051908101604052809291908181526020018280548015611a7957602002820191906000526020600020905b815481526020019060010190808311611a65575b50505050509050919050565b6000818152600183016020526040812054611acc575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105be565b5060006105be565b600060208284031215611ae657600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146109ea57600080fd5b6000815180845260005b81811015611b3c57602081850181015186830182015201611b20565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006109ea6020830184611b16565b803573ffffffffffffffffffffffffffffffffffffffff81168114611bb157600080fd5b919050565b60008060408385031215611bc957600080fd5b611bd283611b8d565b946020939093013593505050565b600080600060608486031215611bf557600080fd5b611bfe84611b8d565b9250611c0c60208501611b8d565b9150604084013590509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600060608486031215611c6057600080fd5b611c6984611b8d565b925060208401359150604084013567ffffffffffffffff80821115611c8d57600080fd5b818601915086601f830112611ca157600080fd5b813581811115611cb357611cb3611c1c565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715611cf957611cf9611c1c565b81604052828152896020848701011115611d1257600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b600060208284031215611d4657600080fd5b5035919050565b600060208284031215611d5f57600080fd5b6109ea82611b8d565b6020808252825182820181905260009190848201906040850190845b81811015611db657835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611d84565b50909695505050505050565b60008060408385031215611dd557600080fd5b611dde83611b8d565b9150611dec60208401611b8d565b90509250929050565b600181811c90821680611e0957607f821691505b602082108103611e42577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156105be576105be611e48565b828152604060208201526000611ea36040830184611b16565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff84168152826020820152606060408201526000611ee06060830184611b16565b95945050505050565b818103818111156105be576105be611e48565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var BurnMintERC677ABI = BurnMintERC677MetaData.ABI

var BurnMintERC677Bin = BurnMintERC677MetaData.Bin

func DeployBurnMintERC677(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals_ uint8, maxSupply_ *big.Int) (common.Address, *types.Transaction, *BurnMintERC677, error) {
	parsed, err := BurnMintERC677MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintERC677Bin), backend, name, symbol, decimals_, maxSupply_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintERC677{address: address, abi: *parsed, BurnMintERC677Caller: BurnMintERC677Caller{contract: contract}, BurnMintERC677Transactor: BurnMintERC677Transactor{contract: contract}, BurnMintERC677Filterer: BurnMintERC677Filterer{contract: contract}}, nil
}

type BurnMintERC677 struct {
	address common.Address
	abi     abi.ABI
	BurnMintERC677Caller
	BurnMintERC677Transactor
	BurnMintERC677Filterer
}

type BurnMintERC677Caller struct {
	contract *bind.BoundContract
}

type BurnMintERC677Transactor struct {
	contract *bind.BoundContract
}

type BurnMintERC677Filterer struct {
	contract *bind.BoundContract
}

type BurnMintERC677Session struct {
	Contract     *BurnMintERC677
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintERC677CallerSession struct {
	Contract *BurnMintERC677Caller
	CallOpts bind.CallOpts
}

type BurnMintERC677TransactorSession struct {
	Contract     *BurnMintERC677Transactor
	TransactOpts bind.TransactOpts
}

type BurnMintERC677Raw struct {
	Contract *BurnMintERC677
}

type BurnMintERC677CallerRaw struct {
	Contract *BurnMintERC677Caller
}

type BurnMintERC677TransactorRaw struct {
	Contract *BurnMintERC677Transactor
}

func NewBurnMintERC677(address common.Address, backend bind.ContractBackend) (*BurnMintERC677, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintERC677ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintERC677(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677{address: address, abi: abi, BurnMintERC677Caller: BurnMintERC677Caller{contract: contract}, BurnMintERC677Transactor: BurnMintERC677Transactor{contract: contract}, BurnMintERC677Filterer: BurnMintERC677Filterer{contract: contract}}, nil
}

func NewBurnMintERC677Caller(address common.Address, caller bind.ContractCaller) (*BurnMintERC677Caller, error) {
	contract, err := bindBurnMintERC677(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677Caller{contract: contract}, nil
}

func NewBurnMintERC677Transactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintERC677Transactor, error) {
	contract, err := bindBurnMintERC677(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677Transactor{contract: contract}, nil
}

func NewBurnMintERC677Filterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintERC677Filterer, error) {
	contract, err := bindBurnMintERC677(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677Filterer{contract: contract}, nil
}

func bindBurnMintERC677(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintERC677MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintERC677 *BurnMintERC677Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintERC677.Contract.BurnMintERC677Caller.contract.Call(opts, result, method, params...)
}

func (_BurnMintERC677 *BurnMintERC677Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.BurnMintERC677Transactor.contract.Transfer(opts)
}

func (_BurnMintERC677 *BurnMintERC677Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.BurnMintERC677Transactor.contract.Transact(opts, method, params...)
}

func (_BurnMintERC677 *BurnMintERC677CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintERC677.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintERC677 *BurnMintERC677TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.contract.Transfer(opts)
}

func (_BurnMintERC677 *BurnMintERC677TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintERC677 *BurnMintERC677Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BurnMintERC677.Contract.Allowance(&_BurnMintERC677.CallOpts, owner, spender)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BurnMintERC677.Contract.Allowance(&_BurnMintERC677.CallOpts, owner, spender)
}

func (_BurnMintERC677 *BurnMintERC677Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _BurnMintERC677.Contract.BalanceOf(&_BurnMintERC677.CallOpts, account)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _BurnMintERC677.Contract.BalanceOf(&_BurnMintERC677.CallOpts, account)
}

func (_BurnMintERC677 *BurnMintERC677Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) Decimals() (uint8, error) {
	return _BurnMintERC677.Contract.Decimals(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) Decimals() (uint8, error) {
	return _BurnMintERC677.Contract.Decimals(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) GetBurners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "getBurners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) GetBurners() ([]common.Address, error) {
	return _BurnMintERC677.Contract.GetBurners(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) GetBurners() ([]common.Address, error) {
	return _BurnMintERC677.Contract.GetBurners(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "getCCIPAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) GetCCIPAdmin() (common.Address, error) {
	return _BurnMintERC677.Contract.GetCCIPAdmin(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) GetCCIPAdmin() (common.Address, error) {
	return _BurnMintERC677.Contract.GetCCIPAdmin(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) GetMinters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "getMinters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) GetMinters() ([]common.Address, error) {
	return _BurnMintERC677.Contract.GetMinters(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) GetMinters() ([]common.Address, error) {
	return _BurnMintERC677.Contract.GetMinters(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) IsBurner(opts *bind.CallOpts, burner common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "isBurner", burner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) IsBurner(burner common.Address) (bool, error) {
	return _BurnMintERC677.Contract.IsBurner(&_BurnMintERC677.CallOpts, burner)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) IsBurner(burner common.Address) (bool, error) {
	return _BurnMintERC677.Contract.IsBurner(&_BurnMintERC677.CallOpts, burner)
}

func (_BurnMintERC677 *BurnMintERC677Caller) IsMinter(opts *bind.CallOpts, minter common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "isMinter", minter)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) IsMinter(minter common.Address) (bool, error) {
	return _BurnMintERC677.Contract.IsMinter(&_BurnMintERC677.CallOpts, minter)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) IsMinter(minter common.Address) (bool, error) {
	return _BurnMintERC677.Contract.IsMinter(&_BurnMintERC677.CallOpts, minter)
}

func (_BurnMintERC677 *BurnMintERC677Caller) MaxSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "maxSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) MaxSupply() (*big.Int, error) {
	return _BurnMintERC677.Contract.MaxSupply(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) MaxSupply() (*big.Int, error) {
	return _BurnMintERC677.Contract.MaxSupply(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) Name() (string, error) {
	return _BurnMintERC677.Contract.Name(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) Name() (string, error) {
	return _BurnMintERC677.Contract.Name(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) Owner() (common.Address, error) {
	return _BurnMintERC677.Contract.Owner(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) Owner() (common.Address, error) {
	return _BurnMintERC677.Contract.Owner(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintERC677.Contract.SupportsInterface(&_BurnMintERC677.CallOpts, interfaceId)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintERC677.Contract.SupportsInterface(&_BurnMintERC677.CallOpts, interfaceId)
}

func (_BurnMintERC677 *BurnMintERC677Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) Symbol() (string, error) {
	return _BurnMintERC677.Contract.Symbol(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) Symbol() (string, error) {
	return _BurnMintERC677.Contract.Symbol(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC677.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC677 *BurnMintERC677Session) TotalSupply() (*big.Int, error) {
	return _BurnMintERC677.Contract.TotalSupply(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677CallerSession) TotalSupply() (*big.Int, error) {
	return _BurnMintERC677.Contract.TotalSupply(&_BurnMintERC677.CallOpts)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintERC677 *BurnMintERC677Session) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintERC677.Contract.AcceptOwnership(&_BurnMintERC677.TransactOpts)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintERC677.Contract.AcceptOwnership(&_BurnMintERC677.TransactOpts)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "approve", spender, amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Approve(&_BurnMintERC677.TransactOpts, spender, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Approve(&_BurnMintERC677.TransactOpts, spender, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "burn", amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) Burn(amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Burn(&_BurnMintERC677.TransactOpts, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Burn(&_BurnMintERC677.TransactOpts, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "burn0", account, amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Burn0(&_BurnMintERC677.TransactOpts, account, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Burn0(&_BurnMintERC677.TransactOpts, account, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "burnFrom", account, amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.BurnFrom(&_BurnMintERC677.TransactOpts, account, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.BurnFrom(&_BurnMintERC677.TransactOpts, account, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_BurnMintERC677 *BurnMintERC677Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.DecreaseAllowance(&_BurnMintERC677.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.DecreaseAllowance(&_BurnMintERC677.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "decreaseApproval", spender, subtractedValue)
}

func (_BurnMintERC677 *BurnMintERC677Session) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.DecreaseApproval(&_BurnMintERC677.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.DecreaseApproval(&_BurnMintERC677.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) GrantBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "grantBurnRole", burner)
}

func (_BurnMintERC677 *BurnMintERC677Session) GrantBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.GrantBurnRole(&_BurnMintERC677.TransactOpts, burner)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) GrantBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.GrantBurnRole(&_BurnMintERC677.TransactOpts, burner)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "grantMintAndBurnRoles", burnAndMinter)
}

func (_BurnMintERC677 *BurnMintERC677Session) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.GrantMintAndBurnRoles(&_BurnMintERC677.TransactOpts, burnAndMinter)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.GrantMintAndBurnRoles(&_BurnMintERC677.TransactOpts, burnAndMinter)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) GrantMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "grantMintRole", minter)
}

func (_BurnMintERC677 *BurnMintERC677Session) GrantMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.GrantMintRole(&_BurnMintERC677.TransactOpts, minter)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) GrantMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.GrantMintRole(&_BurnMintERC677.TransactOpts, minter)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_BurnMintERC677 *BurnMintERC677Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.IncreaseAllowance(&_BurnMintERC677.TransactOpts, spender, addedValue)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.IncreaseAllowance(&_BurnMintERC677.TransactOpts, spender, addedValue)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "increaseApproval", spender, addedValue)
}

func (_BurnMintERC677 *BurnMintERC677Session) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.IncreaseApproval(&_BurnMintERC677.TransactOpts, spender, addedValue)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.IncreaseApproval(&_BurnMintERC677.TransactOpts, spender, addedValue)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "mint", account, amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Mint(&_BurnMintERC677.TransactOpts, account, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Mint(&_BurnMintERC677.TransactOpts, account, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) RevokeBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "revokeBurnRole", burner)
}

func (_BurnMintERC677 *BurnMintERC677Session) RevokeBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.RevokeBurnRole(&_BurnMintERC677.TransactOpts, burner)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) RevokeBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.RevokeBurnRole(&_BurnMintERC677.TransactOpts, burner)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) RevokeMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "revokeMintRole", minter)
}

func (_BurnMintERC677 *BurnMintERC677Session) RevokeMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.RevokeMintRole(&_BurnMintERC677.TransactOpts, minter)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) RevokeMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.RevokeMintRole(&_BurnMintERC677.TransactOpts, minter)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "setCCIPAdmin", newAdmin)
}

func (_BurnMintERC677 *BurnMintERC677Session) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.SetCCIPAdmin(&_BurnMintERC677.TransactOpts, newAdmin)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.SetCCIPAdmin(&_BurnMintERC677.TransactOpts, newAdmin)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "transfer", to, amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Transfer(&_BurnMintERC677.TransactOpts, to, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.Transfer(&_BurnMintERC677.TransactOpts, to, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "transferAndCall", to, amount, data)
}

func (_BurnMintERC677 *BurnMintERC677Session) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.TransferAndCall(&_BurnMintERC677.TransactOpts, to, amount, data)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.TransferAndCall(&_BurnMintERC677.TransactOpts, to, amount, data)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "transferFrom", from, to, amount)
}

func (_BurnMintERC677 *BurnMintERC677Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.TransferFrom(&_BurnMintERC677.TransactOpts, from, to, amount)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.TransferFrom(&_BurnMintERC677.TransactOpts, from, to, amount)
}

func (_BurnMintERC677 *BurnMintERC677Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintERC677 *BurnMintERC677Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.TransferOwnership(&_BurnMintERC677.TransactOpts, to)
}

func (_BurnMintERC677 *BurnMintERC677TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintERC677.Contract.TransferOwnership(&_BurnMintERC677.TransactOpts, to)
}

type BurnMintERC677ApprovalIterator struct {
	Event *BurnMintERC677Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677Approval)
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
		it.Event = new(BurnMintERC677Approval)
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

func (it *BurnMintERC677ApprovalIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC677ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677ApprovalIterator{contract: _BurnMintERC677.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC677Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677Approval)
				if err := _BurnMintERC677.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseApproval(log types.Log) (*BurnMintERC677Approval, error) {
	event := new(BurnMintERC677Approval)
	if err := _BurnMintERC677.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677BurnAccessGrantedIterator struct {
	Event *BurnMintERC677BurnAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677BurnAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677BurnAccessGranted)
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
		it.Event = new(BurnMintERC677BurnAccessGranted)
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

func (it *BurnMintERC677BurnAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677BurnAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677BurnAccessGranted struct {
	Burner common.Address
	Raw    types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterBurnAccessGranted(opts *bind.FilterOpts, burner []common.Address) (*BurnMintERC677BurnAccessGrantedIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "BurnAccessGranted", burnerRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677BurnAccessGrantedIterator{contract: _BurnMintERC677.contract, event: "BurnAccessGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC677BurnAccessGranted, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "BurnAccessGranted", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677BurnAccessGranted)
				if err := _BurnMintERC677.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseBurnAccessGranted(log types.Log) (*BurnMintERC677BurnAccessGranted, error) {
	event := new(BurnMintERC677BurnAccessGranted)
	if err := _BurnMintERC677.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677BurnAccessRevokedIterator struct {
	Event *BurnMintERC677BurnAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677BurnAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677BurnAccessRevoked)
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
		it.Event = new(BurnMintERC677BurnAccessRevoked)
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

func (it *BurnMintERC677BurnAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677BurnAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677BurnAccessRevoked struct {
	Burner common.Address
	Raw    types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterBurnAccessRevoked(opts *bind.FilterOpts, burner []common.Address) (*BurnMintERC677BurnAccessRevokedIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "BurnAccessRevoked", burnerRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677BurnAccessRevokedIterator{contract: _BurnMintERC677.contract, event: "BurnAccessRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC677BurnAccessRevoked, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "BurnAccessRevoked", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677BurnAccessRevoked)
				if err := _BurnMintERC677.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseBurnAccessRevoked(log types.Log) (*BurnMintERC677BurnAccessRevoked, error) {
	event := new(BurnMintERC677BurnAccessRevoked)
	if err := _BurnMintERC677.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677CCIPAdminUpdatedIterator struct {
	Event *BurnMintERC677CCIPAdminUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677CCIPAdminUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677CCIPAdminUpdated)
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
		it.Event = new(BurnMintERC677CCIPAdminUpdated)
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

func (it *BurnMintERC677CCIPAdminUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677CCIPAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677CCIPAdminUpdated struct {
	OldAdmin common.Address
	NewAdmin common.Address
	Raw      types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterCCIPAdminUpdated(opts *bind.FilterOpts) (*BurnMintERC677CCIPAdminUpdatedIterator, error) {

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "CCIPAdminUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677CCIPAdminUpdatedIterator{contract: _BurnMintERC677.contract, event: "CCIPAdminUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchCCIPAdminUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintERC677CCIPAdminUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "CCIPAdminUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677CCIPAdminUpdated)
				if err := _BurnMintERC677.contract.UnpackLog(event, "CCIPAdminUpdated", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseCCIPAdminUpdated(log types.Log) (*BurnMintERC677CCIPAdminUpdated, error) {
	event := new(BurnMintERC677CCIPAdminUpdated)
	if err := _BurnMintERC677.contract.UnpackLog(event, "CCIPAdminUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677MintAccessGrantedIterator struct {
	Event *BurnMintERC677MintAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677MintAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677MintAccessGranted)
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
		it.Event = new(BurnMintERC677MintAccessGranted)
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

func (it *BurnMintERC677MintAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677MintAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677MintAccessGranted struct {
	Minter common.Address
	Raw    types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterMintAccessGranted(opts *bind.FilterOpts, minter []common.Address) (*BurnMintERC677MintAccessGrantedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "MintAccessGranted", minterRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677MintAccessGrantedIterator{contract: _BurnMintERC677.contract, event: "MintAccessGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC677MintAccessGranted, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "MintAccessGranted", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677MintAccessGranted)
				if err := _BurnMintERC677.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseMintAccessGranted(log types.Log) (*BurnMintERC677MintAccessGranted, error) {
	event := new(BurnMintERC677MintAccessGranted)
	if err := _BurnMintERC677.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677MintAccessRevokedIterator struct {
	Event *BurnMintERC677MintAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677MintAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677MintAccessRevoked)
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
		it.Event = new(BurnMintERC677MintAccessRevoked)
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

func (it *BurnMintERC677MintAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677MintAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677MintAccessRevoked struct {
	Minter common.Address
	Raw    types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterMintAccessRevoked(opts *bind.FilterOpts, minter []common.Address) (*BurnMintERC677MintAccessRevokedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "MintAccessRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677MintAccessRevokedIterator{contract: _BurnMintERC677.contract, event: "MintAccessRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC677MintAccessRevoked, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "MintAccessRevoked", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677MintAccessRevoked)
				if err := _BurnMintERC677.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseMintAccessRevoked(log types.Log) (*BurnMintERC677MintAccessRevoked, error) {
	event := new(BurnMintERC677MintAccessRevoked)
	if err := _BurnMintERC677.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677OwnershipTransferRequestedIterator struct {
	Event *BurnMintERC677OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677OwnershipTransferRequested)
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
		it.Event = new(BurnMintERC677OwnershipTransferRequested)
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

func (it *BurnMintERC677OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677OwnershipTransferRequestedIterator{contract: _BurnMintERC677.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintERC677OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677OwnershipTransferRequested)
				if err := _BurnMintERC677.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintERC677OwnershipTransferRequested, error) {
	event := new(BurnMintERC677OwnershipTransferRequested)
	if err := _BurnMintERC677.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677OwnershipTransferredIterator struct {
	Event *BurnMintERC677OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677OwnershipTransferred)
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
		it.Event = new(BurnMintERC677OwnershipTransferred)
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

func (it *BurnMintERC677OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677OwnershipTransferredIterator{contract: _BurnMintERC677.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC677OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677OwnershipTransferred)
				if err := _BurnMintERC677.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseOwnershipTransferred(log types.Log) (*BurnMintERC677OwnershipTransferred, error) {
	event := new(BurnMintERC677OwnershipTransferred)
	if err := _BurnMintERC677.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677TransferIterator struct {
	Event *BurnMintERC677Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677Transfer)
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
		it.Event = new(BurnMintERC677Transfer)
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

func (it *BurnMintERC677TransferIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Data  []byte
	Raw   types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677TransferIterator{contract: _BurnMintERC677.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnMintERC677Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677Transfer)
				if err := _BurnMintERC677.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseTransfer(log types.Log) (*BurnMintERC677Transfer, error) {
	event := new(BurnMintERC677Transfer)
	if err := _BurnMintERC677.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC677Transfer0Iterator struct {
	Event *BurnMintERC677Transfer0

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC677Transfer0Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC677Transfer0)
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
		it.Event = new(BurnMintERC677Transfer0)
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

func (it *BurnMintERC677Transfer0Iterator) Error() error {
	return it.fail
}

func (it *BurnMintERC677Transfer0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC677Transfer0 struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_BurnMintERC677 *BurnMintERC677Filterer) FilterTransfer0(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677Transfer0Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.FilterLogs(opts, "Transfer0", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC677Transfer0Iterator{contract: _BurnMintERC677.contract, event: "Transfer0", logs: logs, sub: sub}, nil
}

func (_BurnMintERC677 *BurnMintERC677Filterer) WatchTransfer0(opts *bind.WatchOpts, sink chan<- *BurnMintERC677Transfer0, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC677.contract.WatchLogs(opts, "Transfer0", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC677Transfer0)
				if err := _BurnMintERC677.contract.UnpackLog(event, "Transfer0", log); err != nil {
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

func (_BurnMintERC677 *BurnMintERC677Filterer) ParseTransfer0(log types.Log) (*BurnMintERC677Transfer0, error) {
	event := new(BurnMintERC677Transfer0)
	if err := _BurnMintERC677.contract.UnpackLog(event, "Transfer0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintERC677 *BurnMintERC677) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintERC677.abi.Events["Approval"].ID:
		return _BurnMintERC677.ParseApproval(log)
	case _BurnMintERC677.abi.Events["BurnAccessGranted"].ID:
		return _BurnMintERC677.ParseBurnAccessGranted(log)
	case _BurnMintERC677.abi.Events["BurnAccessRevoked"].ID:
		return _BurnMintERC677.ParseBurnAccessRevoked(log)
	case _BurnMintERC677.abi.Events["CCIPAdminUpdated"].ID:
		return _BurnMintERC677.ParseCCIPAdminUpdated(log)
	case _BurnMintERC677.abi.Events["MintAccessGranted"].ID:
		return _BurnMintERC677.ParseMintAccessGranted(log)
	case _BurnMintERC677.abi.Events["MintAccessRevoked"].ID:
		return _BurnMintERC677.ParseMintAccessRevoked(log)
	case _BurnMintERC677.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintERC677.ParseOwnershipTransferRequested(log)
	case _BurnMintERC677.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintERC677.ParseOwnershipTransferred(log)
	case _BurnMintERC677.abi.Events["Transfer"].ID:
		return _BurnMintERC677.ParseTransfer(log)
	case _BurnMintERC677.abi.Events["Transfer0"].ID:
		return _BurnMintERC677.ParseTransfer0(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintERC677Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (BurnMintERC677BurnAccessGranted) Topic() common.Hash {
	return common.HexToHash("0x92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad")
}

func (BurnMintERC677BurnAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0x0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c")
}

func (BurnMintERC677CCIPAdminUpdated) Topic() common.Hash {
	return common.HexToHash("0xe6aed30f8fcb01c95cfd4d0512b436882bb2d9005792b8d1b791bc7505614d21")
}

func (BurnMintERC677MintAccessGranted) Topic() common.Hash {
	return common.HexToHash("0xe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea")
}

func (BurnMintERC677MintAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0xed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e9")
}

func (BurnMintERC677OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintERC677OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintERC677Transfer) Topic() common.Hash {
	return common.HexToHash("0xe19260aff97b920c7df27010903aeb9c8d2be5d310a2c67824cf3f15396e4c16")
}

func (BurnMintERC677Transfer0) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (_BurnMintERC677 *BurnMintERC677) Address() common.Address {
	return _BurnMintERC677.address
}

type BurnMintERC677Interface interface {
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

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC677ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC677Approval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*BurnMintERC677Approval, error)

	FilterBurnAccessGranted(opts *bind.FilterOpts, burner []common.Address) (*BurnMintERC677BurnAccessGrantedIterator, error)

	WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC677BurnAccessGranted, burner []common.Address) (event.Subscription, error)

	ParseBurnAccessGranted(log types.Log) (*BurnMintERC677BurnAccessGranted, error)

	FilterBurnAccessRevoked(opts *bind.FilterOpts, burner []common.Address) (*BurnMintERC677BurnAccessRevokedIterator, error)

	WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC677BurnAccessRevoked, burner []common.Address) (event.Subscription, error)

	ParseBurnAccessRevoked(log types.Log) (*BurnMintERC677BurnAccessRevoked, error)

	FilterCCIPAdminUpdated(opts *bind.FilterOpts) (*BurnMintERC677CCIPAdminUpdatedIterator, error)

	WatchCCIPAdminUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintERC677CCIPAdminUpdated) (event.Subscription, error)

	ParseCCIPAdminUpdated(log types.Log) (*BurnMintERC677CCIPAdminUpdated, error)

	FilterMintAccessGranted(opts *bind.FilterOpts, minter []common.Address) (*BurnMintERC677MintAccessGrantedIterator, error)

	WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC677MintAccessGranted, minter []common.Address) (event.Subscription, error)

	ParseMintAccessGranted(log types.Log) (*BurnMintERC677MintAccessGranted, error)

	FilterMintAccessRevoked(opts *bind.FilterOpts, minter []common.Address) (*BurnMintERC677MintAccessRevokedIterator, error)

	WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC677MintAccessRevoked, minter []common.Address) (event.Subscription, error)

	ParseMintAccessRevoked(log types.Log) (*BurnMintERC677MintAccessRevoked, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintERC677OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintERC677OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC677OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintERC677OwnershipTransferred, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnMintERC677Transfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*BurnMintERC677Transfer, error)

	FilterTransfer0(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC677Transfer0Iterator, error)

	WatchTransfer0(opts *bind.WatchOpts, sink chan<- *BurnMintERC677Transfer0, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer0(log types.Log) (*BurnMintERC677Transfer0, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
