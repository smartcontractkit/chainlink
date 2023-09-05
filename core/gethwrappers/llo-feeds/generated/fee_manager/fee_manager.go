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
	Weight *big.Int
}

type CommonAsset struct {
	AssetAddress common.Address
	Amount       *big.Int
}

type IFeeManagerQuote struct {
	QuoteAddress common.Address
}

var FeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nativeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proxyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardManagerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ExpiredReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDeposit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDiscount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidQuote\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReceivingAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSurcharge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroDeficit\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nativeQuantity\",\"type\":\"uint256\"}],\"name\":\"InsufficientLink\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"}],\"name\":\"LinkDeficitCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSurcharge\",\"type\":\"uint256\"}],\"name\":\"NativeSurchargeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"SubscriberDiscountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"quoteAddress\",\"type\":\"address\"}],\"internalType\":\"structIFeeManager.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"getFeeAndReward\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"payLinkDeficit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"}],\"name\":\"processFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_linkDeficit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nativeSurcharge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_subscriberDiscounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"surcharge\",\"type\":\"uint256\"}],\"name\":\"setNativeSurcharge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"updateSubscriberDiscount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200296738038062002967833981016040819052620000359162000288565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001c0565b5050506001600160a01b0384161580620000e057506001600160a01b038316155b80620000f357506001600160a01b038216155b806200010657506001600160a01b038116155b15620001255760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03848116608081905284821660a05283821660c05290821660e081905260405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156200018f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b59190620002e5565b505050505062000310565b336001600160a01b038216036200021a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200028357600080fd5b919050565b600080600080608085870312156200029f57600080fd5b620002aa856200026b565b9350620002ba602086016200026b565b9250620002ca604086016200026b565b9150620002da606086016200026b565b905092959194509250565b600060208284031215620002f857600080fd5b815180151581146200030957600080fd5b9392505050565b60805160a05160c05160e0516125ad620003ba6000396000818161047e01528181610d1a0152818161115701526113280152600081816103e70152610de101526000818161071601528181610769015281816109fc01528181610f6301528181610fd1015261149501526000818161073b015281816107c40152818161090e0152818161095901528181610ab0015281816110a9015281816111e8015261143e01526125ad6000f3fe6080604052600436106100f35760003560e01c80638da5cb5b1161008a578063e389d9a411610059578063e389d9a414610339578063f1387e1614610359578063f237f1a81461036c578063f2fde38b1461038c57600080fd5b80638da5cb5b146102a1578063c541cbde146102d6578063d09dc33914610304578063d9caed121461031957600080fd5b806369fd2b34116100c657806369fd2b341461020c5780636b54d8a61461022e57806379ba50971461024e57806387d6d8431461026357600080fd5b8063013f542b146100f857806301ffc9a714610138578063181f5a77146101aa57806332f5f746146101f6575b600080fd5b34801561010457600080fd5b50610125610113366004611bb0565b60036020526000908152604090205481565b6040519081526020015b60405180910390f35b34801561014457600080fd5b5061019a610153366004611bc9565b7fffffffff00000000000000000000000000000000000000000000000000000000167ff1387e16000000000000000000000000000000000000000000000000000000001490565b604051901515815260200161012f565b3480156101b657600080fd5b50604080518082018252601081527f4665654d616e6167657220302e302e31000000000000000000000000000000006020820152905161012f9190611c2f565b34801561020257600080fd5b5061012560045481565b34801561021857600080fd5b5061022c610227366004611c80565b6103ac565b005b34801561023a57600080fd5b5061022c610249366004611bb0565b6104ee565b34801561025a57600080fd5b5061022c610573565b34801561026f57600080fd5b5061012561027e366004611d21565b600260209081526000938452604080852082529284528284209052825290205481565b3480156102ad57600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012f565b3480156102e257600080fd5b506102f66102f1366004611e98565b610675565b60405161012f929190611f39565b34801561031057600080fd5b50610125610a7f565b34801561032557600080fd5b5061022c610334366004611f8d565b610b35565b34801561034557600080fd5b5061022c610354366004611bb0565b610c7c565b61022c610367366004611fce565b610dc9565b34801561037857600080fd5b5061022c610387366004612046565b6113f2565b34801561039857600080fd5b5061022c6103a736600461208e565b611599565b60005473ffffffffffffffffffffffffffffffffffffffff16331480159061040a57503373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614155b15610441576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f633b5f6e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063633b5f6e906104b7908690869086906004016120ab565b600060405180830381600087803b1580156104d157600080fd5b505af11580156104e5573d6000803e3d6000fd5b50505050505050565b6104f66115ad565b670de0b6b3a7640000811115610538576040517f05e8ac2900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60048190556040518181527fa320ddc8288125050af9b9b75fd872443f8c722e70b14fc8475dbd630a4916e19060200160405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b604080518082019091526000808252602082015260408051808201909152600080825260208201526040805180820190915260008082526020820152604080518082019091526000808252602082015260006106d08761211a565b90507fffff00000000000000000000000000000000000000000000000000000000000080821690810161076757505073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811683527f00000000000000000000000000000000000000000000000000000000000000001681529092509050610a77565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16876000015173ffffffffffffffffffffffffffffffffffffffff161415801561081757507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16876000015173ffffffffffffffffffffffffffffffffffffffff1614155b1561084e576040517ff861803000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060008a80602001905181019061086791906121a0565b77ffffffffffffffffffffffffffffffffffffffffffffffff92831698509116955063ffffffff16935050504282101590506108cf576040517fb6c405f500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff808d16600090815260026020908152604080832089845282528083208e51851684528252909120547f00000000000000000000000000000000000000000000000000000000000000008316808a529189018690528c51909216036109b25773ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168852670de0b6b3a764000061098e828661223c565b6109989190612279565b6109a290856122b4565b6020808a01829052880152610a6a565b60006109e5600454670de0b6b3a76400006109cd91906122c7565b6109d7908661223c565b670de0b6b3a7640000611630565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168a529050670de0b6b3a7640000610a33838361223c565b610a3d9190612279565b610a4790826122b4565b60208a0152610a596109d7838761223c565b610a6390866122b4565b6020890152505b5095975093955050505050505b935093915050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610b0c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b3091906122da565b905090565b610b3d6115ad565b73ffffffffffffffffffffffffffffffffffffffff8316610bf85760008273ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114610bb2576040519150601f19603f3d011682016040523d82523d6000602084013e610bb7565b606091505b5050905080610bf2576040517fef2af20100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050565b610c1973ffffffffffffffffffffffffffffffffffffffff8416838361166a565b6040805133815273ffffffffffffffffffffffffffffffffffffffff84811660208301528516818301526060810183905290517f3115d1449a7b732c986cba18244e897a450f61e1bb8d589cd2e69e6c8924f9f79181900360800190a15b505050565b610c846115ad565b60008181526003602052604081205490819003610ccd576040517f03aad31200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526003602052604080822091909155517fefb03fe900000000000000000000000000000000000000000000000000000000815260048101839052306024820152604481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063efb03fe990606401600060405180830381600087803b158015610d7357600080fd5b505af1158015610d87573d6000803e3d6000fd5b50505050817f843f0b103e50b42b08f9d30f12f961845a6d02623730872e24644899c0dd989582604051610dbd91815260200190565b60405180910390a25050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610e38576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff821603610e87576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610e9583850185612361565b915050600081610ea49061211a565b6040805160208101909152600081529091507e010000000000000000000000000000000000000000000000000000000000007fffff000000000000000000000000000000000000000000000000000000000000831614610f2b576000610f0c86880188612430565b9550505050505080806020019051810190610f2791906124f8565b9150505b600080610f39868685610675565b6020820151919350915034901561139c57825173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691160361108f573415611062573483602001511115610fcf576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db084602001516040518263ffffffff1660e01b81526004016000604051808303818588803b15801561103b57600080fd5b505af115801561104f573d6000803e3d6000fd5b505050505082602001518103905061108f565b6020830151835161108f9173ffffffffffffffffffffffffffffffffffffffff909116908990309061173e565b60208201511561139c5760006110a5898b612526565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16846000015173ffffffffffffffffffffffffffffffffffffffff16036111ba5760208301516040517fefb03fe90000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff8a8116602483015260448201929092527f00000000000000000000000000000000000000000000000000000000000000009091169063efb03fe990606401600060405180830381600087803b15801561119d57600080fd5b505af11580156111b1573d6000803e3d6000fd5b5050505061139a565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015611244573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061126891906122da565b836020015111156112e657826020015160036000838152602001908152602001600020600082825461129a91906122c7565b909155505060208381015185820151604080519283529282015282917feb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7910160405180910390a261139a565b60208301516040517fefb03fe90000000000000000000000000000000000000000000000000000000081526004810183905230602482015260448101919091527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063efb03fe990606401600060405180830381600087803b15801561138157600080fd5b505af1158015611395573d6000803e3d6000fd5b505050505b505b80156113e75760405173ffffffffffffffffffffffffffffffffffffffff88169082156108fc029083906000818181858888f193505050501580156113e5573d6000803e3d6000fd5b505b505050505050505050565b6113fa6115ad565b670de0b6b3a764000081111561143c576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141580156114e457507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b1561151b576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84811660008181526002602090815260408083208884528252808320948716808452948252918290208590558151938452830184905285927f41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84910160405180910390a350505050565b6115a16115ad565b6115aa8161179c565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461162e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105f0565b565b6000821561165e57816116446001856122b4565b61164e9190612279565b6116599060016122c7565b611661565b60005b90505b92915050565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052610c779084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152611891565b60405173ffffffffffffffffffffffffffffffffffffffff80851660248301528316604482015260648101829052610bf29085907f23b872dd00000000000000000000000000000000000000000000000000000000906084016116bc565b3373ffffffffffffffffffffffffffffffffffffffff82160361181b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105f0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006118f3826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661199d9092919063ffffffff16565b805190915015610c7757808060200190518101906119119190612562565b610c77576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016105f0565b60606119ac84846000856119b6565b90505b9392505050565b606082471015611a48576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016105f0565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051611a719190612584565b60006040518083038185875af1925050503d8060008114611aae576040519150601f19603f3d011682016040523d82523d6000602084013e611ab3565b606091505b5091509150611ac487838387611ad1565b925050505b949350505050565b60608315611b67578251600003611b605773ffffffffffffffffffffffffffffffffffffffff85163b611b60576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016105f0565b5081611ac9565b611ac98383815115611b7c5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105f09190611c2f565b600060208284031215611bc257600080fd5b5035919050565b600060208284031215611bdb57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146119af57600080fd5b60005b83811015611c26578181015183820152602001611c0e565b50506000910152565b6020815260008251806020840152611c4e816040850160208701611c0b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b600080600060408486031215611c9557600080fd5b83359250602084013567ffffffffffffffff80821115611cb457600080fd5b818601915086601f830112611cc857600080fd5b813581811115611cd757600080fd5b8760208260061b8501011115611cec57600080fd5b6020830194508093505050509250925092565b73ffffffffffffffffffffffffffffffffffffffff811681146115aa57600080fd5b600080600060608486031215611d3657600080fd5b8335611d4181611cff565b9250602084013591506040840135611d5881611cff565b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810167ffffffffffffffff81118282101715611db557611db5611d63565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611e0257611e02611d63565b604052919050565b600082601f830112611e1b57600080fd5b813567ffffffffffffffff811115611e3557611e35611d63565b611e6660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611dbb565b818152846020838601011115611e7b57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008385036060811215611eae57600080fd5b8435611eb981611cff565b9350602085013567ffffffffffffffff811115611ed557600080fd5b611ee187828801611e0a565b93505060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082011215611f1457600080fd5b50611f1d611d92565b6040850135611f2b81611cff565b815292959194509192509050565b825173ffffffffffffffffffffffffffffffffffffffff1681526020808401519082015260808101825173ffffffffffffffffffffffffffffffffffffffff166040830152602083015160608301526119af565b600080600060608486031215611fa257600080fd5b8335611fad81611cff565b92506020840135611fbd81611cff565b929592945050506040919091013590565b600080600060408486031215611fe357600080fd5b833567ffffffffffffffff80821115611ffb57600080fd5b818601915086601f83011261200f57600080fd5b81358181111561201e57600080fd5b87602082850101111561203057600080fd5b60209283019550935050840135611d5881611cff565b6000806000806080858703121561205c57600080fd5b843561206781611cff565b935060208501359250604085013561207e81611cff565b9396929550929360600135925050565b6000602082840312156120a057600080fd5b81356119af81611cff565b8381526040602080830182905282820184905260009190859060608501845b8781101561210d5783356120dd81611cff565b73ffffffffffffffffffffffffffffffffffffffff168252838301358383015292840192908401906001016120ca565b5098975050505050505050565b80516020808301519190811015612159577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b805163ffffffff8116811461217357600080fd5b919050565b805177ffffffffffffffffffffffffffffffffffffffffffffffff8116811461217357600080fd5b60008060008060008060c087890312156121b957600080fd5b865195506121c96020880161215f565b94506121d76040880161215f565b93506121e560608801612178565b92506121f360808801612178565b915061220160a0880161215f565b90509295509295509295565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156122745761227461220d565b500290565b6000826122af577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b818103818111156116645761166461220d565b808201808211156116645761166461220d565b6000602082840312156122ec57600080fd5b5051919050565b600082601f83011261230457600080fd5b6040516060810181811067ffffffffffffffff8211171561232757612327611d63565b60405280606084018581111561233c57600080fd5b845b8181101561235657803583526020928301920161233e565b509195945050505050565b6000806080838503121561237457600080fd5b61237e84846122f3565b9150606083013567ffffffffffffffff81111561239a57600080fd5b6123a685828601611e0a565b9150509250929050565b600082601f8301126123c157600080fd5b8135602067ffffffffffffffff8211156123dd576123dd611d63565b8160051b6123ec828201611dbb565b928352848101820192828101908785111561240657600080fd5b83870192505b848310156124255782358252918301919083019061240c565b979650505050505050565b600080600080600080610100878903121561244a57600080fd5b61245488886122f3565b9550606087013567ffffffffffffffff8082111561247157600080fd5b61247d8a838b01611e0a565b9650608089013591508082111561249357600080fd5b61249f8a838b016123b0565b955060a08901359150808211156124b557600080fd5b6124c18a838b016123b0565b945060c0890135935060e08901359150808211156124de57600080fd5b506124eb89828a01611e0a565b9150509295509295509295565b60006020828403121561250a57600080fd5b612512611d92565b825161251d81611cff565b81529392505050565b80356020831015611664577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b60006020828403121561257457600080fd5b815180151581146119af57600080fd5b60008251612596818460208701611c0b565b919091019291505056fea164736f6c6343000810000a",
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

func (_FeeManager *FeeManagerTransactor) SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "setFeeRecipients", configDigest, rewardRecipientAndWeights)
}

func (_FeeManager *FeeManagerSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _FeeManager.Contract.SetFeeRecipients(&_FeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_FeeManager *FeeManagerTransactorSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _FeeManager.Contract.SetFeeRecipients(&_FeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_FeeManager *FeeManagerTransactor) SetNativeSurcharge(opts *bind.TransactOpts, surcharge *big.Int) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "setNativeSurcharge", surcharge)
}

func (_FeeManager *FeeManagerSession) SetNativeSurcharge(surcharge *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.SetNativeSurcharge(&_FeeManager.TransactOpts, surcharge)
}

func (_FeeManager *FeeManagerTransactorSession) SetNativeSurcharge(surcharge *big.Int) (*types.Transaction, error) {
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

func (_FeeManager *FeeManagerTransactor) UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "updateSubscriberDiscount", subscriber, feedId, token, discount)
}

func (_FeeManager *FeeManagerSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.UpdateSubscriberDiscount(&_FeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_FeeManager *FeeManagerTransactorSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.UpdateSubscriberDiscount(&_FeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_FeeManager *FeeManagerTransactor) Withdraw(opts *bind.TransactOpts, assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "withdraw", assetAddress, recipient, quantity)
}

func (_FeeManager *FeeManagerSession) Withdraw(assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.Withdraw(&_FeeManager.TransactOpts, assetAddress, recipient, quantity)
}

func (_FeeManager *FeeManagerTransactorSession) Withdraw(assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.Withdraw(&_FeeManager.TransactOpts, assetAddress, recipient, quantity)
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
	ConfigDigest   [32]byte
	LinkQuantity   *big.Int
	NativeQuantity *big.Int
	Raw            types.Log
}

func (_FeeManager *FeeManagerFilterer) FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*FeeManagerInsufficientLinkIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "InsufficientLink", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerInsufficientLinkIterator{contract: _FeeManager.contract, event: "InsufficientLink", logs: logs, sub: sub}, nil
}

func (_FeeManager *FeeManagerFilterer) WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *FeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "InsufficientLink", configDigestRule)
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
	NewSurcharge *big.Int
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
	Discount   *big.Int
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
	Recipient    common.Address
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
	return common.HexToHash("0xeb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7")
}

func (FeeManagerLinkDeficitCleared) Topic() common.Hash {
	return common.HexToHash("0x843f0b103e50b42b08f9d30f12f961845a6d02623730872e24644899c0dd9895")
}

func (FeeManagerNativeSurchargeUpdated) Topic() common.Hash {
	return common.HexToHash("0xa320ddc8288125050af9b9b75fd872443f8c722e70b14fc8475dbd630a4916e1")
}

func (FeeManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FeeManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FeeManagerSubscriberDiscountUpdated) Topic() common.Hash {
	return common.HexToHash("0x41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84")
}

func (FeeManagerWithdraw) Topic() common.Hash {
	return common.HexToHash("0x3115d1449a7b732c986cba18244e897a450f61e1bb8d589cd2e69e6c8924f9f7")
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

	SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetNativeSurcharge(opts *bind.TransactOpts, surcharge *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error)

	FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*FeeManagerInsufficientLinkIterator, error)

	WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *FeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error)

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
