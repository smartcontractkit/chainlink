// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_fee_manager

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

var MercuryFeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nativeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proxyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardManagerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ExpiredReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDeposit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDiscount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidQuote\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSurcharge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nativeQuantity\",\"type\":\"uint256\"}],\"name\":\"InsufficientLink\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSurcharge\",\"type\":\"uint256\"}],\"name\":\"NativeSurchargeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"SubscriberDiscountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"quoteAddress\",\"type\":\"address\"}],\"internalType\":\"structIFeeManager.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"getFeeAndReward\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nativeSurcharge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"}],\"name\":\"processFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"surcharge\",\"type\":\"uint256\"}],\"name\":\"setNativeSurcharge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriberDiscounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"updateSubscriberDiscount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b50604051620024d4380380620024d4833981016040819052620000359162000288565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001c0565b5050506001600160a01b0384161580620000e057506001600160a01b038316155b80620000f357506001600160a01b038216155b806200010657506001600160a01b038116155b15620001255760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03848116608081905284821660a05283821660c05290821660e081905260405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156200018f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b59190620002e5565b505050505062000310565b336001600160a01b038216036200021a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200028357600080fd5b919050565b600080600080608085870312156200029f57600080fd5b620002aa856200026b565b9350620002ba602086016200026b565b9250620002ca604086016200026b565b9150620002da606086016200026b565b905092959194509250565b600060208284031215620002f857600080fd5b815180151581146200030957600080fd5b9392505050565b60805160a05160c05160e051612128620003ac6000396000818161044c01528181610e5401526110b20152600081816103820152610b1b0152600081816106af015281816107020152818161090f01528181610c6901528181610d3001526112170152600081816106d40152818161075d0152818161089a01528181610a5d01528181610def01528181610fa301526111c001526121286000f3fe6080604052600436106100dd5760003560e01c806393798be71161007f578063f1387e1611610059578063f1387e16146102d6578063f237f1a8146102e9578063f2fde38b14610309578063f3fef3a31461032957600080fd5b806393798be714610255578063c541cbde14610293578063d09dc339146102c157600080fd5b80636b54d8a6116100bb5780636b54d8a6146101c757806379ba5097146101e75780637d75cc49146101fc5780638da5cb5b1461022057600080fd5b806301ffc9a7146100e2578063181f5a771461015957806369fd2b34146101a5575b600080fd5b3480156100ee57600080fd5b506101446100fd36600461167b565b7fffffffff00000000000000000000000000000000000000000000000000000000167ff1387e16000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b34801561016557600080fd5b50604080518082018252601081527f4665654d616e6167657220302e302e31000000000000000000000000000000006020820152905161015091906116c4565b3480156101b157600080fd5b506101c56101c0366004611730565b610349565b005b3480156101d357600080fd5b506101c56101e23660046117af565b6104bc565b3480156101f357600080fd5b506101c5610541565b34801561020857600080fd5b5061021260035481565b604051908152602001610150565b34801561022c57600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610150565b34801561026157600080fd5b506102126102703660046117ea565b600260209081526000938452604080852082529284528284209052825290205481565b34801561029f57600080fd5b506102b36102ae366004611961565b61063e565b604051610150929190611a02565b3480156102cd57600080fd5b50610212610a2c565b6101c56102e4366004611a56565b610ae2565b3480156102f557600080fd5b506101c5610304366004611ace565b611174565b34801561031557600080fd5b506101c5610324366004611b16565b61131b565b34801561033557600080fd5b506101c5610344366004611b33565b61132f565b60005473ffffffffffffffffffffffffffffffffffffffff163314806103a457503373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016145b61040f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4f6e6c79206f776e6572206f722070726f78790000000000000000000000000060448201526064015b60405180910390fd5b6040517f633b5f6e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063633b5f6e9061048590869086908690600401611b5f565b600060405180830381600087803b15801561049f57600080fd5b505af11580156104b3573d6000803e3d6000fd5b50505050505050565b6104c46114c9565b670de0b6b3a7640000811115610506576040517f05e8ac2900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60038190556040518181527f4fbcc2b7f4cb5518be923bd7c7c887e29e07b001e2a4a0fdd47c494696d8a1479060200160405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610406565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60408051808201909152600080825260208201526040805180820190915260008082526020820152604080518082019091526000808252602082015260408051808201909152600080825260208201528551610120106107005773ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811683527f00000000000000000000000000000000000000000000000000000000000000001681529092509050610a24565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16856000015173ffffffffffffffffffffffffffffffffffffffff16141580156107b057507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16856000015173ffffffffffffffffffffffffffffffffffffffff1614155b156107e7576040517ff861803000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000888060200190518101906108009190611c39565b63ffffffff169b5077ffffffffffffffffffffffffffffffffffffffffffffffff169b5077ffffffffffffffffffffffffffffffffffffffffffffffff169b5050505050505050505042811015610883576040517fb6c405f500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116808652602086018590528951909116036108f857835173ffffffffffffffffffffffffffffffffffffffff16855260208085015190860152610968565b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001685526003546109629061094a90670de0b6b3a7640000611d36565b6109549084611d49565b670de0b6b3a764000061154c565b60208601525b60006109738a611d86565b73ffffffffffffffffffffffffffffffffffffffff808d16600090815260026020908152604080832085845282528083208e519094168352928152919020549088015191925090670de0b6b3a7640000906109cf908390611d49565b6109d99190611dcb565b87602001516109e89190611e06565b876020018181525050610a048187602001516109549190611d49565b8660200151610a139190611e06565b602087015250949650929450505050505b935093915050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610ab9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610add9190611e19565b905090565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610b3d57503373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016145b610ba3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4f6e6c79206f776e6572206f722070726f7879000000000000000000000000006044820152606401610406565b3073ffffffffffffffffffffffffffffffffffffffff821603610bf2576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610c0083850185611ea0565b604080516020810190915260008152909250905081516101201015610c4c576000610c2d85870187611f6f565b9550505050505080806020019051810190610c489190612037565b9150505b600080610c5a33858561063e565b909250905060003415610dbd577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16836000015173ffffffffffffffffffffffffffffffffffffffff1614610cf0576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3483602001511115610d2e576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db084602001516040518263ffffffff1660e01b81526004016000604051808303818588803b158015610d9a57600080fd5b505af1158015610dae573d6000803e3d6000fd5b50505050508260200151340390505b6000610dc9888a612065565b60208501519091501561111e57835173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911603610ec4576040517f84afb76e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906384afb76e90610e8d9084908b9088906004016120a1565b600060405180830381600087803b158015610ea757600080fd5b505af1158015610ebb573d6000803e3d6000fd5b5050505061111e565b34600003610f7557835160208501516040517f23b872dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a8116600483015230602483015260448201929092529116906323b872dd906064016020604051808303816000875af1158015610f4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f7391906120f9565b505b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610fff573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110239190611e19565b836020015111156110755760208084015185820151604080519283529282015282917feb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7910160405180910390a261111e565b6040517f84afb76e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906384afb76e906110eb908490309088906004016120a1565b600060405180830381600087803b15801561110557600080fd5b505af1158015611119573d6000803e3d6000fd5b505050505b81156111695760405173ffffffffffffffffffffffffffffffffffffffff88169083156108fc029084906000818181858888f19350505050158015611167573d6000803e3d6000fd5b505b505050505050505050565b61117c6114c9565b670de0b6b3a76400008111156111be576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415801561126657507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b1561129d576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84811660008181526002602090815260408083208884528252808320948716808452948252918290208590558151938452830184905285927f41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84910160405180910390a350505050565b6113236114c9565b61132c81611586565b50565b6113376114c9565b73ffffffffffffffffffffffffffffffffffffffff821661139d576000805460405173ffffffffffffffffffffffffffffffffffffffff9091169183156108fc02918491818181858888f19350505050158015611398573d6000803e3d6000fd5b505050565b8173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6113d860005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602481018490526044016020604051808303816000875af115801561144a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061146e91906120f9565b506040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201529081018290527f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9060600160405180910390a15050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461154a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610406565b565b6000821561157a5781611560600185611e06565b61156a9190611dcb565b611575906001611d36565b61157d565b60005b90505b92915050565b3373ffffffffffffffffffffffffffffffffffffffff821603611605576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610406565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561168d57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146116bd57600080fd5b9392505050565b600060208083528351808285015260005b818110156116f1578581018301518582016040015282016116d5565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60008060006040848603121561174557600080fd5b83359250602084013567ffffffffffffffff8082111561176457600080fd5b818601915086601f83011261177857600080fd5b81358181111561178757600080fd5b8760208260061b850101111561179c57600080fd5b6020830194508093505050509250925092565b6000602082840312156117c157600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461132c57600080fd5b6000806000606084860312156117ff57600080fd5b833561180a816117c8565b9250602084013591506040840135611821816117c8565b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810167ffffffffffffffff8111828210171561187e5761187e61182c565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156118cb576118cb61182c565b604052919050565b600082601f8301126118e457600080fd5b813567ffffffffffffffff8111156118fe576118fe61182c565b61192f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611884565b81815284602083860101111561194457600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000838503606081121561197757600080fd5b8435611982816117c8565b9350602085013567ffffffffffffffff81111561199e57600080fd5b6119aa878288016118d3565b93505060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0820112156119dd57600080fd5b506119e661185b565b60408501356119f4816117c8565b815292959194509192509050565b825173ffffffffffffffffffffffffffffffffffffffff1681526020808401519082015260808101825173ffffffffffffffffffffffffffffffffffffffff166040830152602083015160608301526116bd565b600080600060408486031215611a6b57600080fd5b833567ffffffffffffffff80821115611a8357600080fd5b818601915086601f830112611a9757600080fd5b813581811115611aa657600080fd5b876020828501011115611ab857600080fd5b60209283019550935050840135611821816117c8565b60008060008060808587031215611ae457600080fd5b8435611aef816117c8565b9350602085013592506040850135611b06816117c8565b9396929550929360600135925050565b600060208284031215611b2857600080fd5b81356116bd816117c8565b60008060408385031215611b4657600080fd5b8235611b51816117c8565b946020939093013593505050565b8381526040602080830182905282820184905260009190859060608501845b87811015611bc1578335611b91816117c8565b73ffffffffffffffffffffffffffffffffffffffff16825283830135838301529284019290840190600101611b7e565b5098975050505050505050565b805163ffffffff81168114611be257600080fd5b919050565b8051601781900b8114611be257600080fd5b805167ffffffffffffffff81168114611be257600080fd5b805177ffffffffffffffffffffffffffffffffffffffffffffffff81168114611be257600080fd5b6000806000806000806000806000806000806101808d8f031215611c5c57600080fd5b8c519b50611c6c60208e01611bce565b9a50611c7a60408e01611be7565b9950611c8860608e01611be7565b9850611c9660808e01611be7565b9750611ca460a08e01611bf9565b965060c08d01519550611cb960e08e01611bf9565b9450611cc86101008e01611bf9565b9350611cd76101208e01611c11565b9250611ce66101408e01611c11565b9150611cf56101608e01611bce565b90509295989b509295989b509295989b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561158057611580611d07565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611d8157611d81611d07565b500290565b80516020808301519190811015611dc5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b600082611e01577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b8181038181111561158057611580611d07565b600060208284031215611e2b57600080fd5b5051919050565b600082601f830112611e4357600080fd5b6040516060810181811067ffffffffffffffff82111715611e6657611e6661182c565b604052806060840185811115611e7b57600080fd5b845b81811015611e95578035835260209283019201611e7d565b509195945050505050565b60008060808385031215611eb357600080fd5b611ebd8484611e32565b9150606083013567ffffffffffffffff811115611ed957600080fd5b611ee5858286016118d3565b9150509250929050565b600082601f830112611f0057600080fd5b8135602067ffffffffffffffff821115611f1c57611f1c61182c565b8160051b611f2b828201611884565b9283528481018201928281019087851115611f4557600080fd5b83870192505b84831015611f6457823582529183019190830190611f4b565b979650505050505050565b6000806000806000806101008789031215611f8957600080fd5b611f938888611e32565b9550606087013567ffffffffffffffff80821115611fb057600080fd5b611fbc8a838b016118d3565b96506080890135915080821115611fd257600080fd5b611fde8a838b01611eef565b955060a0890135915080821115611ff457600080fd5b6120008a838b01611eef565b945060c0890135935060e089013591508082111561201d57600080fd5b5061202a89828a016118d3565b9150509295509295509295565b60006020828403121561204957600080fd5b61205161185b565b825161205c816117c8565b81529392505050565b80356020831015611580577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b83815273ffffffffffffffffffffffffffffffffffffffff83166020820152608081016120f16040830184805173ffffffffffffffffffffffffffffffffffffffff168252602090810151910152565b949350505050565b60006020828403121561210b57600080fd5b815180151581146116bd57600080fdfea164736f6c6343000810000a",
}

var MercuryFeeManagerABI = MercuryFeeManagerMetaData.ABI

var MercuryFeeManagerBin = MercuryFeeManagerMetaData.Bin

func DeployMercuryFeeManager(auth *bind.TransactOpts, backend bind.ContractBackend, _linkAddress common.Address, _nativeAddress common.Address, _proxyAddress common.Address, _rewardManagerAddress common.Address) (common.Address, *types.Transaction, *MercuryFeeManager, error) {
	parsed, err := MercuryFeeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryFeeManagerBin), backend, _linkAddress, _nativeAddress, _proxyAddress, _rewardManagerAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryFeeManager{MercuryFeeManagerCaller: MercuryFeeManagerCaller{contract: contract}, MercuryFeeManagerTransactor: MercuryFeeManagerTransactor{contract: contract}, MercuryFeeManagerFilterer: MercuryFeeManagerFilterer{contract: contract}}, nil
}

type MercuryFeeManager struct {
	address common.Address
	abi     abi.ABI
	MercuryFeeManagerCaller
	MercuryFeeManagerTransactor
	MercuryFeeManagerFilterer
}

type MercuryFeeManagerCaller struct {
	contract *bind.BoundContract
}

type MercuryFeeManagerTransactor struct {
	contract *bind.BoundContract
}

type MercuryFeeManagerFilterer struct {
	contract *bind.BoundContract
}

type MercuryFeeManagerSession struct {
	Contract     *MercuryFeeManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryFeeManagerCallerSession struct {
	Contract *MercuryFeeManagerCaller
	CallOpts bind.CallOpts
}

type MercuryFeeManagerTransactorSession struct {
	Contract     *MercuryFeeManagerTransactor
	TransactOpts bind.TransactOpts
}

type MercuryFeeManagerRaw struct {
	Contract *MercuryFeeManager
}

type MercuryFeeManagerCallerRaw struct {
	Contract *MercuryFeeManagerCaller
}

type MercuryFeeManagerTransactorRaw struct {
	Contract *MercuryFeeManagerTransactor
}

func NewMercuryFeeManager(address common.Address, backend bind.ContractBackend) (*MercuryFeeManager, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryFeeManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryFeeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManager{address: address, abi: abi, MercuryFeeManagerCaller: MercuryFeeManagerCaller{contract: contract}, MercuryFeeManagerTransactor: MercuryFeeManagerTransactor{contract: contract}, MercuryFeeManagerFilterer: MercuryFeeManagerFilterer{contract: contract}}, nil
}

func NewMercuryFeeManagerCaller(address common.Address, caller bind.ContractCaller) (*MercuryFeeManagerCaller, error) {
	contract, err := bindMercuryFeeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerCaller{contract: contract}, nil
}

func NewMercuryFeeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryFeeManagerTransactor, error) {
	contract, err := bindMercuryFeeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerTransactor{contract: contract}, nil
}

func NewMercuryFeeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryFeeManagerFilterer, error) {
	contract, err := bindMercuryFeeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerFilterer{contract: contract}, nil
}

func bindMercuryFeeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryFeeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryFeeManager *MercuryFeeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryFeeManager.Contract.MercuryFeeManagerCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryFeeManager *MercuryFeeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.MercuryFeeManagerTransactor.contract.Transfer(opts)
}

func (_MercuryFeeManager *MercuryFeeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.MercuryFeeManagerTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryFeeManager.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.contract.Transfer(opts)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "getFeeAndReward", subscriber, report, quote)

	if err != nil {
		return *new(CommonAsset), *new(CommonAsset), err
	}

	out0 := *abi.ConvertType(out[0], new(CommonAsset)).(*CommonAsset)
	out1 := *abi.ConvertType(out[1], new(CommonAsset)).(*CommonAsset)

	return out0, out1, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) GetFeeAndReward(subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	return _MercuryFeeManager.Contract.GetFeeAndReward(&_MercuryFeeManager.CallOpts, subscriber, report, quote)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) GetFeeAndReward(subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	return _MercuryFeeManager.Contract.GetFeeAndReward(&_MercuryFeeManager.CallOpts, subscriber, report, quote)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _MercuryFeeManager.Contract.LinkAvailableForPayment(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _MercuryFeeManager.Contract.LinkAvailableForPayment(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) NativeSurcharge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "nativeSurcharge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) NativeSurcharge() (*big.Int, error) {
	return _MercuryFeeManager.Contract.NativeSurcharge(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) NativeSurcharge() (*big.Int, error) {
	return _MercuryFeeManager.Contract.NativeSurcharge(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) Owner() (common.Address, error) {
	return _MercuryFeeManager.Contract.Owner(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) Owner() (common.Address, error) {
	return _MercuryFeeManager.Contract.Owner(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) SubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "subscriberDiscounts", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) SubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _MercuryFeeManager.Contract.SubscriberDiscounts(&_MercuryFeeManager.CallOpts, arg0, arg1, arg2)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) SubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _MercuryFeeManager.Contract.SubscriberDiscounts(&_MercuryFeeManager.CallOpts, arg0, arg1, arg2)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MercuryFeeManager.Contract.SupportsInterface(&_MercuryFeeManager.CallOpts, interfaceId)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MercuryFeeManager.Contract.SupportsInterface(&_MercuryFeeManager.CallOpts, interfaceId)
}

func (_MercuryFeeManager *MercuryFeeManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryFeeManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryFeeManager *MercuryFeeManagerSession) TypeAndVersion() (string, error) {
	return _MercuryFeeManager.Contract.TypeAndVersion(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerCallerSession) TypeAndVersion() (string, error) {
	return _MercuryFeeManager.Contract.TypeAndVersion(&_MercuryFeeManager.CallOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "acceptOwnership")
}

func (_MercuryFeeManager *MercuryFeeManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.AcceptOwnership(&_MercuryFeeManager.TransactOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.AcceptOwnership(&_MercuryFeeManager.TransactOpts)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) ProcessFee(opts *bind.TransactOpts, payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "processFee", payload, subscriber)
}

func (_MercuryFeeManager *MercuryFeeManagerSession) ProcessFee(payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.ProcessFee(&_MercuryFeeManager.TransactOpts, payload, subscriber)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) ProcessFee(payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.ProcessFee(&_MercuryFeeManager.TransactOpts, payload, subscriber)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "setFeeRecipients", configDigest, rewardRecipientAndWeights)
}

func (_MercuryFeeManager *MercuryFeeManagerSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.SetFeeRecipients(&_MercuryFeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.SetFeeRecipients(&_MercuryFeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) SetNativeSurcharge(opts *bind.TransactOpts, surcharge *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "setNativeSurcharge", surcharge)
}

func (_MercuryFeeManager *MercuryFeeManagerSession) SetNativeSurcharge(surcharge *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.SetNativeSurcharge(&_MercuryFeeManager.TransactOpts, surcharge)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) SetNativeSurcharge(surcharge *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.SetNativeSurcharge(&_MercuryFeeManager.TransactOpts, surcharge)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "transferOwnership", to)
}

func (_MercuryFeeManager *MercuryFeeManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.TransferOwnership(&_MercuryFeeManager.TransactOpts, to)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.TransferOwnership(&_MercuryFeeManager.TransactOpts, to)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "updateSubscriberDiscount", subscriber, feedId, token, discount)
}

func (_MercuryFeeManager *MercuryFeeManagerSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.UpdateSubscriberDiscount(&_MercuryFeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.UpdateSubscriberDiscount(&_MercuryFeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactor) Withdraw(opts *bind.TransactOpts, assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.contract.Transact(opts, "withdraw", assetAddress, quantity)
}

func (_MercuryFeeManager *MercuryFeeManagerSession) Withdraw(assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.Withdraw(&_MercuryFeeManager.TransactOpts, assetAddress, quantity)
}

func (_MercuryFeeManager *MercuryFeeManagerTransactorSession) Withdraw(assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _MercuryFeeManager.Contract.Withdraw(&_MercuryFeeManager.TransactOpts, assetAddress, quantity)
}

type MercuryFeeManagerInsufficientLinkIterator struct {
	Event *MercuryFeeManagerInsufficientLink

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryFeeManagerInsufficientLinkIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryFeeManagerInsufficientLink)
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
		it.Event = new(MercuryFeeManagerInsufficientLink)
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

func (it *MercuryFeeManagerInsufficientLinkIterator) Error() error {
	return it.fail
}

func (it *MercuryFeeManagerInsufficientLinkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryFeeManagerInsufficientLink struct {
	ConfigDigest   [32]byte
	LinkQuantity   *big.Int
	NativeQuantity *big.Int
	Raw            types.Log
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*MercuryFeeManagerInsufficientLinkIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.FilterLogs(opts, "InsufficientLink", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerInsufficientLinkIterator{contract: _MercuryFeeManager.contract, event: "InsufficientLink", logs: logs, sub: sub}, nil
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.WatchLogs(opts, "InsufficientLink", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryFeeManagerInsufficientLink)
				if err := _MercuryFeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
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

func (_MercuryFeeManager *MercuryFeeManagerFilterer) ParseInsufficientLink(log types.Log) (*MercuryFeeManagerInsufficientLink, error) {
	event := new(MercuryFeeManagerInsufficientLink)
	if err := _MercuryFeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryFeeManagerNativeSurchargeSetIterator struct {
	Event *MercuryFeeManagerNativeSurchargeSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryFeeManagerNativeSurchargeSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryFeeManagerNativeSurchargeSet)
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
		it.Event = new(MercuryFeeManagerNativeSurchargeSet)
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

func (it *MercuryFeeManagerNativeSurchargeSetIterator) Error() error {
	return it.fail
}

func (it *MercuryFeeManagerNativeSurchargeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryFeeManagerNativeSurchargeSet struct {
	NewSurcharge *big.Int
	Raw          types.Log
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) FilterNativeSurchargeSet(opts *bind.FilterOpts) (*MercuryFeeManagerNativeSurchargeSetIterator, error) {

	logs, sub, err := _MercuryFeeManager.contract.FilterLogs(opts, "NativeSurchargeSet")
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerNativeSurchargeSetIterator{contract: _MercuryFeeManager.contract, event: "NativeSurchargeSet", logs: logs, sub: sub}, nil
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) WatchNativeSurchargeSet(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerNativeSurchargeSet) (event.Subscription, error) {

	logs, sub, err := _MercuryFeeManager.contract.WatchLogs(opts, "NativeSurchargeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryFeeManagerNativeSurchargeSet)
				if err := _MercuryFeeManager.contract.UnpackLog(event, "NativeSurchargeSet", log); err != nil {
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

func (_MercuryFeeManager *MercuryFeeManagerFilterer) ParseNativeSurchargeSet(log types.Log) (*MercuryFeeManagerNativeSurchargeSet, error) {
	event := new(MercuryFeeManagerNativeSurchargeSet)
	if err := _MercuryFeeManager.contract.UnpackLog(event, "NativeSurchargeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryFeeManagerOwnershipTransferRequestedIterator struct {
	Event *MercuryFeeManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryFeeManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryFeeManagerOwnershipTransferRequested)
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
		it.Event = new(MercuryFeeManagerOwnershipTransferRequested)
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

func (it *MercuryFeeManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MercuryFeeManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryFeeManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryFeeManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerOwnershipTransferRequestedIterator{contract: _MercuryFeeManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryFeeManagerOwnershipTransferRequested)
				if err := _MercuryFeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MercuryFeeManager *MercuryFeeManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*MercuryFeeManagerOwnershipTransferRequested, error) {
	event := new(MercuryFeeManagerOwnershipTransferRequested)
	if err := _MercuryFeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryFeeManagerOwnershipTransferredIterator struct {
	Event *MercuryFeeManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryFeeManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryFeeManagerOwnershipTransferred)
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
		it.Event = new(MercuryFeeManagerOwnershipTransferred)
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

func (it *MercuryFeeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MercuryFeeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryFeeManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryFeeManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerOwnershipTransferredIterator{contract: _MercuryFeeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryFeeManagerOwnershipTransferred)
				if err := _MercuryFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MercuryFeeManager *MercuryFeeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*MercuryFeeManagerOwnershipTransferred, error) {
	event := new(MercuryFeeManagerOwnershipTransferred)
	if err := _MercuryFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryFeeManagerSubscriberDiscountUpdatedIterator struct {
	Event *MercuryFeeManagerSubscriberDiscountUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryFeeManagerSubscriberDiscountUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryFeeManagerSubscriberDiscountUpdated)
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
		it.Event = new(MercuryFeeManagerSubscriberDiscountUpdated)
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

func (it *MercuryFeeManagerSubscriberDiscountUpdatedIterator) Error() error {
	return it.fail
}

func (it *MercuryFeeManagerSubscriberDiscountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryFeeManagerSubscriberDiscountUpdated struct {
	Subscriber common.Address
	FeedId     [32]byte
	Token      common.Address
	Discount   *big.Int
	Raw        types.Log
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*MercuryFeeManagerSubscriberDiscountUpdatedIterator, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.FilterLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerSubscriberDiscountUpdatedIterator{contract: _MercuryFeeManager.contract, event: "SubscriberDiscountUpdated", logs: logs, sub: sub}, nil
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MercuryFeeManager.contract.WatchLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryFeeManagerSubscriberDiscountUpdated)
				if err := _MercuryFeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
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

func (_MercuryFeeManager *MercuryFeeManagerFilterer) ParseSubscriberDiscountUpdated(log types.Log) (*MercuryFeeManagerSubscriberDiscountUpdated, error) {
	event := new(MercuryFeeManagerSubscriberDiscountUpdated)
	if err := _MercuryFeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryFeeManagerWithdrawIterator struct {
	Event *MercuryFeeManagerWithdraw

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryFeeManagerWithdrawIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryFeeManagerWithdraw)
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
		it.Event = new(MercuryFeeManagerWithdraw)
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

func (it *MercuryFeeManagerWithdrawIterator) Error() error {
	return it.fail
}

func (it *MercuryFeeManagerWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryFeeManagerWithdraw struct {
	AdminAddress common.Address
	AssetAddress common.Address
	Quantity     *big.Int
	Raw          types.Log
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) FilterWithdraw(opts *bind.FilterOpts) (*MercuryFeeManagerWithdrawIterator, error) {

	logs, sub, err := _MercuryFeeManager.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &MercuryFeeManagerWithdrawIterator{contract: _MercuryFeeManager.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

func (_MercuryFeeManager *MercuryFeeManagerFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerWithdraw) (event.Subscription, error) {

	logs, sub, err := _MercuryFeeManager.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryFeeManagerWithdraw)
				if err := _MercuryFeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

func (_MercuryFeeManager *MercuryFeeManagerFilterer) ParseWithdraw(log types.Log) (*MercuryFeeManagerWithdraw, error) {
	event := new(MercuryFeeManagerWithdraw)
	if err := _MercuryFeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MercuryFeeManager *MercuryFeeManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MercuryFeeManager.abi.Events["InsufficientLink"].ID:
		return _MercuryFeeManager.ParseInsufficientLink(log)
	case _MercuryFeeManager.abi.Events["NativeSurchargeSet"].ID:
		return _MercuryFeeManager.ParseNativeSurchargeSet(log)
	case _MercuryFeeManager.abi.Events["OwnershipTransferRequested"].ID:
		return _MercuryFeeManager.ParseOwnershipTransferRequested(log)
	case _MercuryFeeManager.abi.Events["OwnershipTransferred"].ID:
		return _MercuryFeeManager.ParseOwnershipTransferred(log)
	case _MercuryFeeManager.abi.Events["SubscriberDiscountUpdated"].ID:
		return _MercuryFeeManager.ParseSubscriberDiscountUpdated(log)
	case _MercuryFeeManager.abi.Events["Withdraw"].ID:
		return _MercuryFeeManager.ParseWithdraw(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MercuryFeeManagerInsufficientLink) Topic() common.Hash {
	return common.HexToHash("0xeb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7")
}

func (MercuryFeeManagerNativeSurchargeSet) Topic() common.Hash {
	return common.HexToHash("0x4fbcc2b7f4cb5518be923bd7c7c887e29e07b001e2a4a0fdd47c494696d8a147")
}

func (MercuryFeeManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MercuryFeeManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MercuryFeeManagerSubscriberDiscountUpdated) Topic() common.Hash {
	return common.HexToHash("0x41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84")
}

func (MercuryFeeManagerWithdraw) Topic() common.Hash {
	return common.HexToHash("0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb")
}

func (_MercuryFeeManager *MercuryFeeManager) Address() common.Address {
	return _MercuryFeeManager.address
}

type MercuryFeeManagerInterface interface {
	GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	NativeSurcharge(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ProcessFee(opts *bind.TransactOpts, payload []byte, subscriber common.Address) (*types.Transaction, error)

	SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetNativeSurcharge(opts *bind.TransactOpts, surcharge *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, assetAddress common.Address, quantity *big.Int) (*types.Transaction, error)

	FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*MercuryFeeManagerInsufficientLinkIterator, error)

	WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error)

	ParseInsufficientLink(log types.Log) (*MercuryFeeManagerInsufficientLink, error)

	FilterNativeSurchargeSet(opts *bind.FilterOpts) (*MercuryFeeManagerNativeSurchargeSetIterator, error)

	WatchNativeSurchargeSet(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerNativeSurchargeSet) (event.Subscription, error)

	ParseNativeSurchargeSet(log types.Log) (*MercuryFeeManagerNativeSurchargeSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryFeeManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MercuryFeeManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryFeeManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MercuryFeeManagerOwnershipTransferred, error)

	FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*MercuryFeeManagerSubscriberDiscountUpdatedIterator, error)

	WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error)

	ParseSubscriberDiscountUpdated(log types.Log) (*MercuryFeeManagerSubscriberDiscountUpdated, error)

	FilterWithdraw(opts *bind.FilterOpts) (*MercuryFeeManagerWithdrawIterator, error)

	WatchWithdraw(opts *bind.WatchOpts, sink chan<- *MercuryFeeManagerWithdraw) (event.Subscription, error)

	ParseWithdraw(log types.Log) (*MercuryFeeManagerWithdraw, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
