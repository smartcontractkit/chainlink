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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nativeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proxyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardManagerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ExpiredReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDeposit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDiscount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidQuote\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSurcharge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nativeQuantity\",\"type\":\"uint256\"}],\"name\":\"InsufficientLink\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSurcharge\",\"type\":\"uint256\"}],\"name\":\"NativeSurchargeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"SubscriberDiscountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"quoteAddress\",\"type\":\"address\"}],\"internalType\":\"structIFeeManager.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"getFeeAndReward\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"}],\"name\":\"processFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nativeSurcharge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_subscriberDiscounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"surcharge\",\"type\":\"uint256\"}],\"name\":\"setNativeSurcharge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"updateSubscriberDiscount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200242338038062002423833981016040819052620000359162000288565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001c0565b5050506001600160a01b0384161580620000e057506001600160a01b038316155b80620000f357506001600160a01b038216155b806200010657506001600160a01b038116155b15620001255760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03848116608081905284821660a05283821660c05290821660e081905260405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156200018f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b59190620002e5565b505050505062000310565b336001600160a01b038216036200021a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200028357600080fd5b919050565b600080600080608085870312156200029f57600080fd5b620002aa856200026b565b9350620002ba602086016200026b565b9250620002ca604086016200026b565b9150620002da606086016200026b565b905092959194509250565b600060208284031215620002f857600080fd5b815180151581146200030957600080fd5b9392505050565b60805160a05160c05160e051612077620003ac6000396000818161041b01528181610e6d01526110c30152600081816103840152610afa0152600081816106b301528181610706015281816108fa01528181610c6501528181610d2c015261122f0152600081816106d8015281816107610152818161088501528181610a3a01528181610dec01528181610faf01526111d801526120776000f3fe6080604052600436106100dd5760003560e01c80638da5cb5b1161007f578063f1387e1611610059578063f1387e16146102d6578063f237f1a8146102e9578063f2fde38b14610309578063f3fef3a31461032957600080fd5b80638da5cb5b1461025e578063c541cbde14610293578063d09dc339146102c157600080fd5b806369fd2b34116100bb57806369fd2b34146101c95780636b54d8a6146101eb57806379ba50971461020b57806387d6d8431461022057600080fd5b806301ffc9a7146100e2578063181f5a771461015957806332f5f746146101a5575b600080fd5b3480156100ee57600080fd5b506101446100fd366004611693565b7fffffffff00000000000000000000000000000000000000000000000000000000167ff1387e16000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b34801561016557600080fd5b50604080518082018252601081527f4665654d616e6167657220302e302e31000000000000000000000000000000006020820152905161015091906116dc565b3480156101b157600080fd5b506101bb60035481565b604051908152602001610150565b3480156101d557600080fd5b506101e96101e4366004611748565b610349565b005b3480156101f757600080fd5b506101e96102063660046117c7565b61048b565b34801561021757600080fd5b506101e9610510565b34801561022c57600080fd5b506101bb61023b366004611802565b600260209081526000938452604080852082529284528284209052825290205481565b34801561026a57600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610150565b34801561029f57600080fd5b506102b36102ae366004611979565b610612565b604051610150929190611a1a565b3480156102cd57600080fd5b506101bb610a09565b6101e96102e4366004611a6e565b610abf565b3480156102f557600080fd5b506101e9610304366004611ae6565b61118c565b34801561031557600080fd5b506101e9610324366004611b2e565b611333565b34801561033557600080fd5b506101e9610344366004611b4b565b611347565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906103a757503373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614155b156103de576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f633b5f6e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063633b5f6e9061045490869086908690600401611b77565b600060405180830381600087803b15801561046e57600080fd5b505af1158015610482573d6000803e3d6000fd5b50505050505050565b6104936114e1565b670de0b6b3a76400008111156104d5576040517f05e8ac2900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60038190556040518181527fa320ddc8288125050af9b9b75fd872443f8c722e70b14fc8475dbd630a4916e19060200160405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff163314610596576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805180820190915260008082526020820152604080518082019091526000808252602082015260408051808201909152600080825260208201526040805180820190915260008082526020820152600061066d87611be6565b90507fffff00000000000000000000000000000000000000000000000000000000000080821690810161070457505073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811683527f00000000000000000000000000000000000000000000000000000000000000001681529092509050610a01565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16876000015173ffffffffffffffffffffffffffffffffffffffff16141580156107b457507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16876000015173ffffffffffffffffffffffffffffffffffffffff1614155b156107eb576040517ff861803000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060008a8060200190518101906108049190611c6c565b77ffffffffffffffffffffffffffffffffffffffffffffffff92831699509116965063ffffffff169450505042831015915061086e9050576040517fb6c405f500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116808852602088018590528b51909116036108e357855173ffffffffffffffffffffffffffffffffffffffff16875260208087015190880152610953565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016875260035461094d9061093590670de0b6b3a7640000611d22565b61093f9084611d35565b670de0b6b3a7640000611564565b60208801525b73ffffffffffffffffffffffffffffffffffffffff808d16600090815260026020908152604080832089845282528083208e5190941683529281529190205490880151670de0b6b3a7640000906109ab908390611d35565b6109b59190611d72565b88602001516109c49190611dad565b8860200181815250506109e081886020015161093f9190611d35565b87602001516109ef9190611dad565b60208801525095975093955050505050505b935093915050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610a96573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aba9190611dc0565b905090565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610b1d57503373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614155b15610b54576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff821603610ba3576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610bb183850185611e47565b915050600081610bc090611be6565b6040805160208101909152600081529091507e010000000000000000000000000000000000000000000000000000000000007fffff000000000000000000000000000000000000000000000000000000000000831614610c47576000610c2886880188611f16565b9550505050505080806020019051810190610c439190611fde565b9150505b600080610c55338685610612565b91509150600034600014610db9577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16836000015173ffffffffffffffffffffffffffffffffffffffff1614610cec576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3483602001511115610d2a576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db084602001516040518263ffffffff1660e01b81526004016000604051808303818588803b158015610d9657600080fd5b505af1158015610daa573d6000803e3d6000fd5b50505050508260200151340390505b6000610dc5898b61200c565b9050836020015160001461113557835173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911603610ed05760208301516040517fefb03fe90000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff8a8116602483015260448201929092527f00000000000000000000000000000000000000000000000000000000000000009091169063efb03fe990606401600060405180830381600087803b158015610eb357600080fd5b505af1158015610ec7573d6000803e3d6000fd5b50505050611135565b34600003610f8157835160208501516040517f23b872dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8b8116600483015230602483015260448201929092529116906323b872dd906064016020604051808303816000875af1158015610f5b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f7f9190612048565b505b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561100b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061102f9190611dc0565b836020015111156110815760208084015185820151604080519283529282015282917feb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7910160405180910390a2611135565b60208301516040517fefb03fe90000000000000000000000000000000000000000000000000000000081526004810183905230602482015260448101919091527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063efb03fe990606401600060405180830381600087803b15801561111c57600080fd5b505af1158015611130573d6000803e3d6000fd5b505050505b81156111805760405173ffffffffffffffffffffffffffffffffffffffff89169083156108fc029084906000818181858888f1935050505015801561117e573d6000803e3d6000fd5b505b50505050505050505050565b6111946114e1565b670de0b6b3a76400008111156111d6576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415801561127e57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b156112b5576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84811660008181526002602090815260408083208884528252808320948716808452948252918290208590558151938452830184905285927f41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84910160405180910390a350505050565b61133b6114e1565b6113448161159e565b50565b61134f6114e1565b73ffffffffffffffffffffffffffffffffffffffff82166113b5576000805460405173ffffffffffffffffffffffffffffffffffffffff9091169183156108fc02918491818181858888f193505050501580156113b0573d6000803e3d6000fd5b505050565b8173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6113f060005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602481018490526044016020604051808303816000875af1158015611462573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114869190612048565b506040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201529081018290527f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9060600160405180910390a15050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611562576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161058d565b565b600082156115925781611578600185611dad565b6115829190611d72565b61158d906001611d22565b611595565b60005b90505b92915050565b3373ffffffffffffffffffffffffffffffffffffffff82160361161d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161058d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156116a557600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146116d557600080fd5b9392505050565b600060208083528351808285015260005b81811015611709578581018301518582016040015282016116ed565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60008060006040848603121561175d57600080fd5b83359250602084013567ffffffffffffffff8082111561177c57600080fd5b818601915086601f83011261179057600080fd5b81358181111561179f57600080fd5b8760208260061b85010111156117b457600080fd5b6020830194508093505050509250925092565b6000602082840312156117d957600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461134457600080fd5b60008060006060848603121561181757600080fd5b8335611822816117e0565b9250602084013591506040840135611839816117e0565b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810167ffffffffffffffff8111828210171561189657611896611844565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156118e3576118e3611844565b604052919050565b600082601f8301126118fc57600080fd5b813567ffffffffffffffff81111561191657611916611844565b61194760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161189c565b81815284602083860101111561195c57600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000838503606081121561198f57600080fd5b843561199a816117e0565b9350602085013567ffffffffffffffff8111156119b657600080fd5b6119c2878288016118eb565b93505060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0820112156119f557600080fd5b506119fe611873565b6040850135611a0c816117e0565b815292959194509192509050565b825173ffffffffffffffffffffffffffffffffffffffff1681526020808401519082015260808101825173ffffffffffffffffffffffffffffffffffffffff166040830152602083015160608301526116d5565b600080600060408486031215611a8357600080fd5b833567ffffffffffffffff80821115611a9b57600080fd5b818601915086601f830112611aaf57600080fd5b813581811115611abe57600080fd5b876020828501011115611ad057600080fd5b60209283019550935050840135611839816117e0565b60008060008060808587031215611afc57600080fd5b8435611b07816117e0565b9350602085013592506040850135611b1e816117e0565b9396929550929360600135925050565b600060208284031215611b4057600080fd5b81356116d5816117e0565b60008060408385031215611b5e57600080fd5b8235611b69816117e0565b946020939093013593505050565b8381526040602080830182905282820184905260009190859060608501845b87811015611bd9578335611ba9816117e0565b73ffffffffffffffffffffffffffffffffffffffff16825283830135838301529284019290840190600101611b96565b5098975050505050505050565b80516020808301519190811015611c25577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b805163ffffffff81168114611c3f57600080fd5b919050565b805177ffffffffffffffffffffffffffffffffffffffffffffffff81168114611c3f57600080fd5b600080600080600080600060e0888a031215611c8757600080fd5b87519650611c9760208901611c2b565b9550611ca560408901611c2b565b945060608801518060170b8114611cbb57600080fd5b9350611cc960808901611c44565b9250611cd760a08901611c44565b9150611ce560c08901611c2b565b905092959891949750929550565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561159857611598611cf3565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611d6d57611d6d611cf3565b500290565b600082611da8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b8181038181111561159857611598611cf3565b600060208284031215611dd257600080fd5b5051919050565b600082601f830112611dea57600080fd5b6040516060810181811067ffffffffffffffff82111715611e0d57611e0d611844565b604052806060840185811115611e2257600080fd5b845b81811015611e3c578035835260209283019201611e24565b509195945050505050565b60008060808385031215611e5a57600080fd5b611e648484611dd9565b9150606083013567ffffffffffffffff811115611e8057600080fd5b611e8c858286016118eb565b9150509250929050565b600082601f830112611ea757600080fd5b8135602067ffffffffffffffff821115611ec357611ec3611844565b8160051b611ed282820161189c565b9283528481018201928281019087851115611eec57600080fd5b83870192505b84831015611f0b57823582529183019190830190611ef2565b979650505050505050565b6000806000806000806101008789031215611f3057600080fd5b611f3a8888611dd9565b9550606087013567ffffffffffffffff80821115611f5757600080fd5b611f638a838b016118eb565b96506080890135915080821115611f7957600080fd5b611f858a838b01611e96565b955060a0890135915080821115611f9b57600080fd5b611fa78a838b01611e96565b945060c0890135935060e0890135915080821115611fc457600080fd5b50611fd189828a016118eb565b9150509295509295509295565b600060208284031215611ff057600080fd5b611ff8611873565b8251612003816117e0565b81529392505050565b80356020831015611598577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b60006020828403121561205a57600080fd5b815180151581146116d557600080fdfea164736f6c6343000810000a",
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
	return common.HexToHash("0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb")
}

func (_FeeManager *FeeManager) Address() common.Address {
	return _FeeManager.address
}

type FeeManagerInterface interface {
	GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error)

	SSubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ProcessFee(opts *bind.TransactOpts, payload []byte, subscriber common.Address) (*types.Transaction, error)

	SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetNativeSurcharge(opts *bind.TransactOpts, surcharge *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, assetAddress common.Address, quantity *big.Int) (*types.Transaction, error)

	FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*FeeManagerInsufficientLinkIterator, error)

	WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *FeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error)

	ParseInsufficientLink(log types.Log) (*FeeManagerInsufficientLink, error)

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
