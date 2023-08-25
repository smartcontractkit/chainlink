// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fee_manager

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

type CommonAddressAndWeight struct {
	Addr   common.Address
	Weight uint64
}

type CommonAsset struct {
	AssetAddress common.Address
	Amount       *big.Int
}

type IFeeManagerQuote struct {
	QuoteAddress common.Address
}

type IRewardManagerFeePayment struct {
	PoolId [32]byte
	Amount *big.Int
}

var FeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nativeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proxyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardManagerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ExpiredReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDeposit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDiscount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidQuote\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSurcharge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroDeficit\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nativeQuantity\",\"type\":\"uint256\"}],\"name\":\"InsufficientLink\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"}],\"name\":\"LinkDeficitCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSurcharge\",\"type\":\"uint256\"}],\"name\":\"NativeSurchargeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"SubscriberDiscountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"quoteAddress\",\"type\":\"address\"}],\"internalType\":\"structIFeeManager.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"getFeeAndReward\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"payLinkDeficit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"}],\"name\":\"processFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_linkDeficit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nativeSurcharge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_subscriberDiscounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"surcharge\",\"type\":\"uint256\"}],\"name\":\"setNativeSurcharge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"updateSubscriberDiscount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200293f3803806200293f833981016040819052620000359162000288565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001c0565b5050506001600160a01b0384161580620000e057506001600160a01b038316155b80620000f357506001600160a01b038216155b806200010657506001600160a01b038116155b15620001255760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03848116608081905284821660a05283821660c05290821660e081905260405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156200018f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b59190620002e5565b505050505062000310565b336001600160a01b038216036200021a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200028357600080fd5b919050565b600080600080608085870312156200029f57600080fd5b620002aa856200026b565b9350620002ba602086016200026b565b9250620002ca604086016200026b565b9150620002da606086016200026b565b905092959194509250565b600060208284031215620002f857600080fd5b815180151581146200030957600080fd5b9392505050565b60805160a05160c05160e051612585620003ba6000396000818161047e01528181610bd30152818161103601526112070152600081816103e70152610cbd01526000818161071601528181610769015281816109fc01528181610e4201528181610eb0015261137401526000818161073b015281816107c40152818161090e0152818161095901528181610ab001528181610f88015281816110c7015261131d01526125856000f3fe6080604052600436106100f35760003560e01c80638da5cb5b1161008a578063f1387e1611610059578063f1387e1614610339578063f237f1a81461034c578063f2fde38b1461036c578063f3fef3a31461038c57600080fd5b80638da5cb5b146102a1578063c541cbde146102d6578063d09dc33914610304578063e389d9a41461031957600080fd5b806369fd2b34116100c657806369fd2b341461020c5780636b54d8a61461022e57806379ba50971461024e57806387d6d8431461026357600080fd5b8063013f542b146100f857806301ffc9a714610138578063181f5a77146101aa57806332f5f746146101f6575b600080fd5b34801561010457600080fd5b50610125610113366004611b9d565b60036020526000908152604090205481565b6040519081526020015b60405180910390f35b34801561014457600080fd5b5061019a610153366004611bb6565b7fffffffff00000000000000000000000000000000000000000000000000000000167ff1387e16000000000000000000000000000000000000000000000000000000001490565b604051901515815260200161012f565b3480156101b657600080fd5b50604080518082018252601081527f4665654d616e6167657220302e302e31000000000000000000000000000000006020820152905161012f9190611c1c565b34801561020257600080fd5b5061012560045481565b34801561021857600080fd5b5061022c610227366004611c6d565b6103ac565b005b34801561023a57600080fd5b5061022c610249366004611b9d565b6104ee565b34801561025a57600080fd5b5061022c610573565b34801561026f57600080fd5b5061012561027e366004611d0e565b600260209081526000938452604080852082529284528284209052825290205481565b3480156102ad57600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012f565b3480156102e257600080fd5b506102f66102f1366004611e85565b610675565b60405161012f929190611f26565b34801561031057600080fd5b50610125610a7f565b34801561032557600080fd5b5061022c610334366004611b9d565b610b35565b61022c610347366004611f7a565b610c82565b34801561035857600080fd5b5061022c610367366004611ff2565b6112d1565b34801561037857600080fd5b5061022c61038736600461203a565b611478565b34801561039857600080fd5b5061022c6103a7366004612057565b61148c565b60005473ffffffffffffffffffffffffffffffffffffffff16331480159061040a57503373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614155b15610441576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f633b5f6e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063633b5f6e906104b790869086908690600401612083565b600060405180830381600087803b1580156104d157600080fd5b505af11580156104e5573d6000803e3d6000fd5b50505050505050565b6104f6611594565b670de0b6b3a7640000811115610538576040517f05e8ac2900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60048190556040518181527fa320ddc8288125050af9b9b75fd872443f8c722e70b14fc8475dbd630a4916e19060200160405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b604080518082019091526000808252602082015260408051808201909152600080825260208201526040805180820190915260008082526020820152604080518082019091526000808252602082015260006106d0876120f2565b90507fffff00000000000000000000000000000000000000000000000000000000000080821690810161076757505073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811683527f00000000000000000000000000000000000000000000000000000000000000001681529092509050610a77565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16876000015173ffffffffffffffffffffffffffffffffffffffff161415801561081757507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16876000015173ffffffffffffffffffffffffffffffffffffffff1614155b1561084e576040517ff861803000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060008a8060200190518101906108679190612178565b77ffffffffffffffffffffffffffffffffffffffffffffffff92831698509116955063ffffffff16935050504282101590506108cf576040517fb6c405f500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff808d16600090815260026020908152604080832089845282528083208e51851684528252909120547f00000000000000000000000000000000000000000000000000000000000000008316808a529189018690528c51909216036109b25773ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168852670de0b6b3a764000061098e8286612214565b6109989190612251565b6109a2908561228c565b6020808a01829052880152610a6a565b60006109e5600454670de0b6b3a76400006109cd919061229f565b6109d79086612214565b670de0b6b3a7640000611617565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168a529050670de0b6b3a7640000610a338383612214565b610a3d9190612251565b610a47908261228c565b60208a0152610a596109d78387612214565b610a63908661228c565b6020890152505b5095975093955050505050505b935093915050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610b0c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b3091906122b2565b905090565b610b3d611594565b60008181526003602052604081205490819003610b86576040517f03aad31200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526003602052604080822091909155517fefb03fe900000000000000000000000000000000000000000000000000000000815260048101839052306024820152604481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063efb03fe990606401600060405180830381600087803b158015610c2c57600080fd5b505af1158015610c40573d6000803e3d6000fd5b50505050817f843f0b103e50b42b08f9d30f12f961845a6d02623730872e24644899c0dd989582604051610c7691815260200190565b60405180910390a25050565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610ce057503373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614155b15610d17576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff821603610d66576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610d7483850185612339565b915050600081610d83906120f2565b6040805160208101909152600081529091507e010000000000000000000000000000000000000000000000000000000000007fffff000000000000000000000000000000000000000000000000000000000000831614610e0a576000610deb86880188612408565b9550505050505080806020019051810190610e0691906124d0565b9150505b600080610e18868685610675565b6020820151919350915034901561127b57825173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911603610f6e573415610f41573483602001511115610eae576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db084602001516040518263ffffffff1660e01b81526004016000604051808303818588803b158015610f1a57600080fd5b505af1158015610f2e573d6000803e3d6000fd5b5050505050826020015181039050610f6e565b60208301518351610f6e9173ffffffffffffffffffffffffffffffffffffffff9091169089903090611651565b60208201511561127b576000610f84898b6124fe565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16846000015173ffffffffffffffffffffffffffffffffffffffff16036110995760208301516040517fefb03fe90000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff8a8116602483015260448201929092527f00000000000000000000000000000000000000000000000000000000000000009091169063efb03fe990606401600060405180830381600087803b15801561107c57600080fd5b505af1158015611090573d6000803e3d6000fd5b50505050611279565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015611123573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061114791906122b2565b836020015111156111c5578260200151600360008381526020019081526020016000206000828254611179919061229f565b909155505060208381015185820151604080519283529282015282917feb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7910160405180910390a2611279565b60208301516040517fefb03fe90000000000000000000000000000000000000000000000000000000081526004810183905230602482015260448101919091527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063efb03fe990606401600060405180830381600087803b15801561126057600080fd5b505af1158015611274573d6000803e3d6000fd5b505050505b505b80156112c65760405173ffffffffffffffffffffffffffffffffffffffff88169082156108fc029083906000818181858888f193505050501580156112c4573d6000803e3d6000fd5b505b505050505050505050565b6112d9611594565b670de0b6b3a764000081111561131b576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141580156113c357507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b156113fa576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84811660008181526002602090815260408083208884528252808320948716808452948252918290208590558151938452830184905285927f41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84910160405180910390a350505050565b611480611594565b61148981611733565b50565b611494611594565b73ffffffffffffffffffffffffffffffffffffffff82166114fa576000805460405173ffffffffffffffffffffffffffffffffffffffff9091169183156108fc02918491818181858888f193505050501580156114f5573d6000803e3d6000fd5b505050565b61153a61151c60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff84169083611828565b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201529081018290527f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9060600160405180910390a15050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611615576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105f0565b565b60008215611645578161162b60018561228c565b6116359190612251565b61164090600161229f565b611648565b60005b90505b92915050565b60405173ffffffffffffffffffffffffffffffffffffffff8085166024830152831660448201526064810182905261172d9085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915261187e565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036117b2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105f0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526114f59084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064016116ab565b60006118e0826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661198a9092919063ffffffff16565b8051909150156114f557808060200190518101906118fe919061253a565b6114f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016105f0565b606061199984846000856119a3565b90505b9392505050565b606082471015611a35576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016105f0565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051611a5e919061255c565b60006040518083038185875af1925050503d8060008114611a9b576040519150601f19603f3d011682016040523d82523d6000602084013e611aa0565b606091505b5091509150611ab187838387611abe565b925050505b949350505050565b60608315611b54578251600003611b4d5773ffffffffffffffffffffffffffffffffffffffff85163b611b4d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016105f0565b5081611ab6565b611ab68383815115611b695781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105f09190611c1c565b600060208284031215611baf57600080fd5b5035919050565b600060208284031215611bc857600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461199c57600080fd5b60005b83811015611c13578181015183820152602001611bfb565b50506000910152565b6020815260008251806020840152611c3b816040850160208701611bf8565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b600080600060408486031215611c8257600080fd5b83359250602084013567ffffffffffffffff80821115611ca157600080fd5b818601915086601f830112611cb557600080fd5b813581811115611cc457600080fd5b8760208260061b8501011115611cd957600080fd5b6020830194508093505050509250925092565b73ffffffffffffffffffffffffffffffffffffffff8116811461148957600080fd5b600080600060608486031215611d2357600080fd5b8335611d2e81611cec565b9250602084013591506040840135611d4581611cec565b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810167ffffffffffffffff81118282101715611da257611da2611d50565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611def57611def611d50565b604052919050565b600082601f830112611e0857600080fd5b813567ffffffffffffffff811115611e2257611e22611d50565b611e5360207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611da8565b818152846020838601011115611e6857600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008385036060811215611e9b57600080fd5b8435611ea681611cec565b9350602085013567ffffffffffffffff811115611ec257600080fd5b611ece87828801611df7565b93505060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082011215611f0157600080fd5b50611f0a611d7f565b6040850135611f1881611cec565b815292959194509192509050565b825173ffffffffffffffffffffffffffffffffffffffff1681526020808401519082015260808101825173ffffffffffffffffffffffffffffffffffffffff1660408301526020830151606083015261199c565b600080600060408486031215611f8f57600080fd5b833567ffffffffffffffff80821115611fa757600080fd5b818601915086601f830112611fbb57600080fd5b813581811115611fca57600080fd5b876020828501011115611fdc57600080fd5b60209283019550935050840135611d4581611cec565b6000806000806080858703121561200857600080fd5b843561201381611cec565b935060208501359250604085013561202a81611cec565b9396929550929360600135925050565b60006020828403121561204c57600080fd5b813561199c81611cec565b6000806040838503121561206a57600080fd5b823561207581611cec565b946020939093013593505050565b8381526040602080830182905282820184905260009190859060608501845b878110156120e55783356120b581611cec565b73ffffffffffffffffffffffffffffffffffffffff168252838301358383015292840192908401906001016120a2565b5098975050505050505050565b80516020808301519190811015612131577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b805163ffffffff8116811461214b57600080fd5b919050565b805177ffffffffffffffffffffffffffffffffffffffffffffffff8116811461214b57600080fd5b60008060008060008060c0878903121561219157600080fd5b865195506121a160208801612137565b94506121af60408801612137565b93506121bd60608801612150565b92506121cb60808801612150565b91506121d960a08801612137565b90509295509295509295565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561224c5761224c6121e5565b500290565b600082612287577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b8181038181111561164b5761164b6121e5565b8082018082111561164b5761164b6121e5565b6000602082840312156122c457600080fd5b5051919050565b600082601f8301126122dc57600080fd5b6040516060810181811067ffffffffffffffff821117156122ff576122ff611d50565b60405280606084018581111561231457600080fd5b845b8181101561232e578035835260209283019201612316565b509195945050505050565b6000806080838503121561234c57600080fd5b61235684846122cb565b9150606083013567ffffffffffffffff81111561237257600080fd5b61237e85828601611df7565b9150509250929050565b600082601f83011261239957600080fd5b8135602067ffffffffffffffff8211156123b5576123b5611d50565b8160051b6123c4828201611da8565b92835284810182019282810190878511156123de57600080fd5b83870192505b848310156123fd578235825291830191908301906123e4565b979650505050505050565b600080600080600080610100878903121561242257600080fd5b61242c88886122cb565b9550606087013567ffffffffffffffff8082111561244957600080fd5b6124558a838b01611df7565b9650608089013591508082111561246b57600080fd5b6124778a838b01612388565b955060a089013591508082111561248d57600080fd5b6124998a838b01612388565b945060c0890135935060e08901359150808211156124b657600080fd5b506124c389828a01611df7565b9150509295509295509295565b6000602082840312156124e257600080fd5b6124ea611d7f565b82516124f581611cec565b81529392505050565b8035602083101561164b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b60006020828403121561254c57600080fd5b8151801515811461199c57600080fd5b6000825161256e818460208701611bf8565b919091019291505056fea164736f6c6343000810000a",
}

var FeeManagerABI = FeeManagerMetaData.ABI

var FeeManagerBin = FeeManagerMetaData.Bin

func DeployFeeManager(auth *bind.TransactOpts, backend bind.ContractBackend, _linkAddress common.Address, _nativeAddress common.Address, _proxyAddress common.Address, _rewardManagerAddress common.Address) (common.Address, *types.Transaction, *FeeManager, error) {
	parsed, err := FeeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeeManagerBin), backend, _linkAddress, _nativeAddress, _proxyAddress, _rewardManagerAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FeeManager{FeeManagerCaller: FeeManagerCaller{contract: contract}, FeeManagerTransactor: FeeManagerTransactor{contract: contract}, FeeManagerFilterer: FeeManagerFilterer{contract: contract}}, nil
}

type FeeManager struct {
	address common.Address
	abi     abi.ABI
	FeeManagerCaller
	FeeManagerTransactor
	FeeManagerFilterer
}

type FeeManagerCaller struct {
	contract *bind.BoundContract
}

type FeeManagerTransactor struct {
	contract *bind.BoundContract
}

type FeeManagerFilterer struct {
	contract *bind.BoundContract
}

type FeeManagerSession struct {
	Contract     *FeeManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FeeManagerCallerSession struct {
	Contract *FeeManagerCaller
	CallOpts bind.CallOpts
}

type FeeManagerTransactorSession struct {
	Contract     *FeeManagerTransactor
	TransactOpts bind.TransactOpts
}

type FeeManagerRaw struct {
	Contract *FeeManager
}

type FeeManagerCallerRaw struct {
	Contract *FeeManagerCaller
}

type FeeManagerTransactorRaw struct {
	Contract *FeeManagerTransactor
}

func NewFeeManager(address common.Address, backend bind.ContractBackend) (*FeeManager, error) {
	abi, err := abi.JSON(strings.NewReader(FeeManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFeeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeManager{address: address, abi: abi, FeeManagerCaller: FeeManagerCaller{contract: contract}, FeeManagerTransactor: FeeManagerTransactor{contract: contract}, FeeManagerFilterer: FeeManagerFilterer{contract: contract}}, nil
}

func NewFeeManagerCaller(address common.Address, caller bind.ContractCaller) (*FeeManagerCaller, error) {
	contract, err := bindFeeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeManagerCaller{contract: contract}, nil
}

func NewFeeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeManagerTransactor, error) {
	contract, err := bindFeeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeManagerTransactor{contract: contract}, nil
}

func NewFeeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeManagerFilterer, error) {
	contract, err := bindFeeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeManagerFilterer{contract: contract}, nil
}

func bindFeeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FeeManager *FeeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeManager.Contract.FeeManagerCaller.contract.Call(opts, result, method, params...)
}

func (_FeeManager *FeeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeManager.Contract.FeeManagerTransactor.contract.Transfer(opts)
}

func (_FeeManager *FeeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeManager.Contract.FeeManagerTransactor.contract.Transact(opts, method, params...)
}

func (_FeeManager *FeeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeManager.Contract.contract.Call(opts, result, method, params...)
}

func (_FeeManager *FeeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeManager.Contract.contract.Transfer(opts)
}

func (_FeeManager *FeeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeManager.Contract.contract.Transact(opts, method, params...)
}

func (_FeeManager *FeeManagerCaller) GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "getFeeAndReward", subscriber, report, quote)

	if err != nil {
		return *new(CommonAsset), *new(CommonAsset), err
	}

	out0 := *abi.ConvertType(out[0], new(CommonAsset)).(*CommonAsset)
	out1 := *abi.ConvertType(out[1], new(CommonAsset)).(*CommonAsset)

	return out0, out1, err

}

func (_FeeManager *FeeManagerSession) GetFeeAndReward(subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	return _FeeManager.Contract.GetFeeAndReward(&_FeeManager.CallOpts, subscriber, report, quote)
}

func (_FeeManager *FeeManagerCallerSession) GetFeeAndReward(subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	return _FeeManager.Contract.GetFeeAndReward(&_FeeManager.CallOpts, subscriber, report, quote)
}

func (_FeeManager *FeeManagerCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeManager *FeeManagerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _FeeManager.Contract.LinkAvailableForPayment(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _FeeManager.Contract.LinkAvailableForPayment(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FeeManager *FeeManagerSession) Owner() (common.Address, error) {
	return _FeeManager.Contract.Owner(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCallerSession) Owner() (common.Address, error) {
	return _FeeManager.Contract.Owner(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCaller) SLinkDeficit(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "s_linkDeficit", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeManager *FeeManagerSession) SLinkDeficit(arg0 [32]byte) (*big.Int, error) {
	return _FeeManager.Contract.SLinkDeficit(&_FeeManager.CallOpts, arg0)
}

func (_FeeManager *FeeManagerCallerSession) SLinkDeficit(arg0 [32]byte) (*big.Int, error) {
	return _FeeManager.Contract.SLinkDeficit(&_FeeManager.CallOpts, arg0)
}

func (_FeeManager *FeeManagerCaller) SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "s_nativeSurcharge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeManager *FeeManagerSession) SNativeSurcharge() (*big.Int, error) {
	return _FeeManager.Contract.SNativeSurcharge(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCallerSession) SNativeSurcharge() (*big.Int, error) {
	return _FeeManager.Contract.SNativeSurcharge(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCaller) SSubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "s_subscriberDiscounts", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeManager *FeeManagerSession) SSubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _FeeManager.Contract.SSubscriberDiscounts(&_FeeManager.CallOpts, arg0, arg1, arg2)
}

func (_FeeManager *FeeManagerCallerSession) SSubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _FeeManager.Contract.SSubscriberDiscounts(&_FeeManager.CallOpts, arg0, arg1, arg2)
}

func (_FeeManager *FeeManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FeeManager *FeeManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeManager.Contract.SupportsInterface(&_FeeManager.CallOpts, interfaceId)
}

func (_FeeManager *FeeManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeManager.Contract.SupportsInterface(&_FeeManager.CallOpts, interfaceId)
}

func (_FeeManager *FeeManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_FeeManager *FeeManagerSession) TypeAndVersion() (string, error) {
	return _FeeManager.Contract.TypeAndVersion(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerCallerSession) TypeAndVersion() (string, error) {
	return _FeeManager.Contract.TypeAndVersion(&_FeeManager.CallOpts)
}

func (_FeeManager *FeeManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "acceptOwnership")
}

func (_FeeManager *FeeManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeeManager.Contract.AcceptOwnership(&_FeeManager.TransactOpts)
}

func (_FeeManager *FeeManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeeManager.Contract.AcceptOwnership(&_FeeManager.TransactOpts)
}

func (_FeeManager *FeeManagerTransactor) PayLinkDeficit(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "payLinkDeficit", configDigest)
}

func (_FeeManager *FeeManagerSession) PayLinkDeficit(configDigest [32]byte) (*types.Transaction, error) {
	return _FeeManager.Contract.PayLinkDeficit(&_FeeManager.TransactOpts, configDigest)
}

func (_FeeManager *FeeManagerTransactorSession) PayLinkDeficit(configDigest [32]byte) (*types.Transaction, error) {
	return _FeeManager.Contract.PayLinkDeficit(&_FeeManager.TransactOpts, configDigest)
}

func (_FeeManager *FeeManagerTransactor) ProcessFee(opts *bind.TransactOpts, payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "processFee", payload, subscriber)
}

func (_FeeManager *FeeManagerSession) ProcessFee(payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.ProcessFee(&_FeeManager.TransactOpts, payload, subscriber)
}

func (_FeeManager *FeeManagerTransactorSession) ProcessFee(payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.ProcessFee(&_FeeManager.TransactOpts, payload, subscriber)
}

func (_FeeManager *FeeManagerTransactor) ProcessFeeBulk(opts *bind.TransactOpts, payloads [][]byte, subscriber common.Address) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "processFeeBulk", payloads, subscriber)
}

func (_FeeManager *FeeManagerSession) ProcessFeeBulk(payloads [][]byte, subscriber common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.ProcessFeeBulk(&_FeeManager.TransactOpts, payloads, subscriber)
}

func (_FeeManager *FeeManagerTransactorSession) ProcessFeeBulk(payloads [][]byte, subscriber common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.ProcessFeeBulk(&_FeeManager.TransactOpts, payloads, subscriber)
}

func (_FeeManager *FeeManagerTransactor) SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "setFeeRecipients", configDigest, rewardRecipientAndWeights)
}

func (_FeeManager *FeeManagerSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _FeeManager.Contract.SetFeeRecipients(&_FeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_FeeManager *FeeManagerTransactorSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _FeeManager.Contract.SetFeeRecipients(&_FeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_FeeManager *FeeManagerTransactor) SetNativeSurcharge(opts *bind.TransactOpts, surcharge uint64) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "setNativeSurcharge", surcharge)
}

func (_FeeManager *FeeManagerSession) SetNativeSurcharge(surcharge uint64) (*types.Transaction, error) {
	return _FeeManager.Contract.SetNativeSurcharge(&_FeeManager.TransactOpts, surcharge)
}

func (_FeeManager *FeeManagerTransactorSession) SetNativeSurcharge(surcharge uint64) (*types.Transaction, error) {
	return _FeeManager.Contract.SetNativeSurcharge(&_FeeManager.TransactOpts, surcharge)
}

func (_FeeManager *FeeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "transferOwnership", to)
}

func (_FeeManager *FeeManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.TransferOwnership(&_FeeManager.TransactOpts, to)
}

func (_FeeManager *FeeManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.TransferOwnership(&_FeeManager.TransactOpts, to)
}

func (_FeeManager *FeeManagerTransactor) UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "updateSubscriberDiscount", subscriber, feedId, token, discount)
}

func (_FeeManager *FeeManagerSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	return _FeeManager.Contract.UpdateSubscriberDiscount(&_FeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_FeeManager *FeeManagerTransactorSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	return _FeeManager.Contract.UpdateSubscriberDiscount(&_FeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_FeeManager *FeeManagerTransactor) Withdraw(opts *bind.TransactOpts, assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "withdraw", assetAddress, quantity)
}

func (_FeeManager *FeeManagerSession) Withdraw(assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.Withdraw(&_FeeManager.TransactOpts, assetAddress, quantity)
}

func (_FeeManager *FeeManagerTransactorSession) Withdraw(assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.Withdraw(&_FeeManager.TransactOpts, assetAddress, quantity)
}

type FeeManagerInsufficientLinkIterator struct {
	Event *FeeManagerInsufficientLink

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerInsufficientLinkIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerInsufficientLink)
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
		it.Event = new(FeeManagerInsufficientLink)
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

func (it *FeeManagerInsufficientLinkIterator) Error() error {
	return it.fail
}

func (it *FeeManagerInsufficientLinkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerInsufficientLink struct {
	Rewards []IRewardManagerFeePayment
	Raw     types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterInsufficientLink(opts *bind.FilterOpts) (*FeeManagerInsufficientLinkIterator, error) {

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "InsufficientLink")
	if err != nil {
		return nil, err
	}
	return &FeeManagerInsufficientLinkIterator{contract: _FeeManager.contract, event: "InsufficientLink", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *FeeManagerInsufficientLink) (event.Subscription, error) {

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "InsufficientLink")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerInsufficientLink)
				if err := _FeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseInsufficientLink(log types.Log) (*FeeManagerInsufficientLink, error) {
	event := new(FeeManagerInsufficientLink)
	if err := _FeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeManagerLinkDeficitClearedIterator struct {
	Event *FeeManagerLinkDeficitCleared

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerLinkDeficitClearedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerLinkDeficitCleared)
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
		it.Event = new(FeeManagerLinkDeficitCleared)
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

func (it *FeeManagerLinkDeficitClearedIterator) Error() error {
	return it.fail
}

func (it *FeeManagerLinkDeficitClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerLinkDeficitCleared struct {
	ConfigDigest [32]byte
	LinkQuantity *big.Int
	Raw          types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterLinkDeficitCleared(opts *bind.FilterOpts, configDigest [][32]byte) (*FeeManagerLinkDeficitClearedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "LinkDeficitCleared", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerLinkDeficitClearedIterator{contract: _FeeManager.contract, event: "LinkDeficitCleared", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchLinkDeficitCleared(opts *bind.WatchOpts, sink chan<- *FeeManagerLinkDeficitCleared, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "LinkDeficitCleared", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerLinkDeficitCleared)
				if err := _FeeManager.contract.UnpackLog(event, "LinkDeficitCleared", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseLinkDeficitCleared(log types.Log) (*FeeManagerLinkDeficitCleared, error) {
	event := new(FeeManagerLinkDeficitCleared)
	if err := _FeeManager.contract.UnpackLog(event, "LinkDeficitCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeManagerNativeSurchargeUpdatedIterator struct {
	Event *FeeManagerNativeSurchargeUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerNativeSurchargeUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerNativeSurchargeUpdated)
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
		it.Event = new(FeeManagerNativeSurchargeUpdated)
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

func (it *FeeManagerNativeSurchargeUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeManagerNativeSurchargeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerNativeSurchargeUpdated struct {
	NewSurcharge uint64
	Raw          types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterNativeSurchargeUpdated(opts *bind.FilterOpts) (*FeeManagerNativeSurchargeUpdatedIterator, error) {

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "NativeSurchargeUpdated")
	if err != nil {
		return nil, err
	}
	return &FeeManagerNativeSurchargeUpdatedIterator{contract: _FeeManager.contract, event: "NativeSurchargeUpdated", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchNativeSurchargeUpdated(opts *bind.WatchOpts, sink chan<- *FeeManagerNativeSurchargeUpdated) (event.Subscription, error) {

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "NativeSurchargeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerNativeSurchargeUpdated)
				if err := _FeeManager.contract.UnpackLog(event, "NativeSurchargeUpdated", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseNativeSurchargeUpdated(log types.Log) (*FeeManagerNativeSurchargeUpdated, error) {
	event := new(FeeManagerNativeSurchargeUpdated)
	if err := _FeeManager.contract.UnpackLog(event, "NativeSurchargeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeManagerOwnershipTransferRequestedIterator struct {
	Event *FeeManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerOwnershipTransferRequested)
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
		it.Event = new(FeeManagerOwnershipTransferRequested)
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

func (it *FeeManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FeeManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerOwnershipTransferRequestedIterator{contract: _FeeManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerOwnershipTransferRequested)
				if err := _FeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*FeeManagerOwnershipTransferRequested, error) {
	event := new(FeeManagerOwnershipTransferRequested)
	if err := _FeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeManagerOwnershipTransferredIterator struct {
	Event *FeeManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerOwnershipTransferred)
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
		it.Event = new(FeeManagerOwnershipTransferred)
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

func (it *FeeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FeeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerOwnershipTransferredIterator{contract: _FeeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerOwnershipTransferred)
				if err := _FeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*FeeManagerOwnershipTransferred, error) {
	event := new(FeeManagerOwnershipTransferred)
	if err := _FeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeManagerSubscriberDiscountUpdatedIterator struct {
	Event *FeeManagerSubscriberDiscountUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerSubscriberDiscountUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerSubscriberDiscountUpdated)
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
		it.Event = new(FeeManagerSubscriberDiscountUpdated)
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

func (it *FeeManagerSubscriberDiscountUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeManagerSubscriberDiscountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerSubscriberDiscountUpdated struct {
	Subscriber common.Address
	FeedId     [32]byte
	Token      common.Address
	Discount   uint64
	Raw        types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*FeeManagerSubscriberDiscountUpdatedIterator, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerSubscriberDiscountUpdatedIterator{contract: _FeeManager.contract, event: "SubscriberDiscountUpdated", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *FeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerSubscriberDiscountUpdated)
				if err := _FeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseSubscriberDiscountUpdated(log types.Log) (*FeeManagerSubscriberDiscountUpdated, error) {
	event := new(FeeManagerSubscriberDiscountUpdated)
	if err := _FeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeManagerWithdrawIterator struct {
	Event *FeeManagerWithdraw

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeManagerWithdrawIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerWithdraw)
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
		it.Event = new(FeeManagerWithdraw)
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

func (it *FeeManagerWithdrawIterator) Error() error {
	return it.fail
}

func (it *FeeManagerWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeManagerWithdraw struct {
	AdminAddress common.Address
	AssetAddress common.Address
	Quantity     *big.Int
	Raw          types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterWithdraw(opts *bind.FilterOpts) (*FeeManagerWithdrawIterator, error) {

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &FeeManagerWithdrawIterator{contract: _FeeManager.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *FeeManagerWithdraw) (event.Subscription, error) {

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeManagerWithdraw)
				if err := _FeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

func (_FeeManager *FeeManagerFilterer) ParseWithdraw(log types.Log) (*FeeManagerWithdraw, error) {
	event := new(FeeManagerWithdraw)
	if err := _FeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FeeManager *FeeManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FeeManager.abi.Events["InsufficientLink"].ID:
		return _FeeManager.ParseInsufficientLink(log)
	case _FeeManager.abi.Events["LinkDeficitCleared"].ID:
		return _FeeManager.ParseLinkDeficitCleared(log)
	case _FeeManager.abi.Events["NativeSurchargeUpdated"].ID:
		return _FeeManager.ParseNativeSurchargeUpdated(log)
	case _FeeManager.abi.Events["OwnershipTransferRequested"].ID:
		return _FeeManager.ParseOwnershipTransferRequested(log)
	case _FeeManager.abi.Events["OwnershipTransferred"].ID:
		return _FeeManager.ParseOwnershipTransferred(log)
	case _FeeManager.abi.Events["SubscriberDiscountUpdated"].ID:
		return _FeeManager.ParseSubscriberDiscountUpdated(log)
	case _FeeManager.abi.Events["Withdraw"].ID:
		return _FeeManager.ParseWithdraw(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FeeManagerInsufficientLink) Topic() common.Hash {
	return common.HexToHash("0xf52e5907b69d97c33392936c12d78b494463b78c5b72df50b4c497eee5720b67")
}

func (FeeManagerLinkDeficitCleared) Topic() common.Hash {
	return common.HexToHash("0x843f0b103e50b42b08f9d30f12f961845a6d02623730872e24644899c0dd9895")
}

func (FeeManagerNativeSurchargeUpdated) Topic() common.Hash {
	return common.HexToHash("0x08f7c0d17932ddb8523bc06754d42ff19ebc77d76a8b9bfde02c28ab1ed3d639")
}

func (FeeManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FeeManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FeeManagerSubscriberDiscountUpdated) Topic() common.Hash {
	return common.HexToHash("0x5eba5a8afa39780f0f99b6cbeb95f3da6a7040ca00abd46bdc91a0a060134139")
}

func (FeeManagerWithdraw) Topic() common.Hash {
	return common.HexToHash("0x72608e45b52a95a12c2ac7f15ff53f92fc9572c9d84b6e6b5d7f0f7826cf3271")
}

func (_FeeManager *FeeManager) Address() common.Address {
	return _FeeManager.address
}

type FeeManagerInterface interface {
	GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SLinkDeficit(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)

	SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error)

	SSubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	PayLinkDeficit(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	ProcessFee(opts *bind.TransactOpts, payload []byte, subscriber common.Address) (*types.Transaction, error)

	ProcessFeeBulk(opts *bind.TransactOpts, payloads [][]byte, subscriber common.Address) (*types.Transaction, error)

	SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetNativeSurcharge(opts *bind.TransactOpts, surcharge uint64) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, assetAddress common.Address, quantity *big.Int) (*types.Transaction, error)

	FilterInsufficientLink(opts *bind.FilterOpts) (*FeeManagerInsufficientLinkIterator, error)

	WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *FeeManagerInsufficientLink) (event.Subscription, error)

	ParseInsufficientLink(log types.Log) (*FeeManagerInsufficientLink, error)

	FilterLinkDeficitCleared(opts *bind.FilterOpts, configDigest [][32]byte) (*FeeManagerLinkDeficitClearedIterator, error)

	WatchLinkDeficitCleared(opts *bind.WatchOpts, sink chan<- *FeeManagerLinkDeficitCleared, configDigest [][32]byte) (event.Subscription, error)

	ParseLinkDeficitCleared(log types.Log) (*FeeManagerLinkDeficitCleared, error)

	FilterNativeSurchargeUpdated(opts *bind.FilterOpts) (*FeeManagerNativeSurchargeUpdatedIterator, error)

	WatchNativeSurchargeUpdated(opts *bind.WatchOpts, sink chan<- *FeeManagerNativeSurchargeUpdated) (event.Subscription, error)

	ParseNativeSurchargeUpdated(log types.Log) (*FeeManagerNativeSurchargeUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FeeManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FeeManagerOwnershipTransferred, error)

	FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*FeeManagerSubscriberDiscountUpdatedIterator, error)

	WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *FeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error)

	ParseSubscriberDiscountUpdated(log types.Log) (*FeeManagerSubscriberDiscountUpdated, error)

	FilterWithdraw(opts *bind.FilterOpts) (*FeeManagerWithdrawIterator, error)

	WatchWithdraw(opts *bind.WatchOpts, sink chan<- *FeeManagerWithdraw) (event.Subscription, error)

	ParseWithdraw(log types.Log) (*FeeManagerWithdraw, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
