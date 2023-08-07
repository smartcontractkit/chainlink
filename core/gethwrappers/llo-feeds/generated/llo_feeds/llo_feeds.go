// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package llo_feeds

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

var LLOFeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nativeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proxyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardManagerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ExpiredReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDeposit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDiscount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidQuote\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSurcharge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"linkQuantity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nativeQuantity\",\"type\":\"uint256\"}],\"name\":\"InsufficientLink\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSurcharge\",\"type\":\"uint256\"}],\"name\":\"NativeSurchargeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"SubscriberDiscountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"quoteAddress\",\"type\":\"address\"}],\"internalType\":\"structIFeeManager.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"getFeeAndReward\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nativeSurcharge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"}],\"name\":\"processFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"surcharge\",\"type\":\"uint256\"}],\"name\":\"setNativeSurcharge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriberDiscounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subscriber\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"discount\",\"type\":\"uint256\"}],\"name\":\"updateSubscriberDiscount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b50604051620024d4380380620024d4833981016040819052620000359162000288565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001c0565b5050506001600160a01b0384161580620000e057506001600160a01b038316155b80620000f357506001600160a01b038216155b806200010657506001600160a01b038116155b15620001255760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03848116608081905284821660a05283821660c05290821660e081905260405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156200018f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b59190620002e5565b505050505062000310565b336001600160a01b038216036200021a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200028357600080fd5b919050565b600080600080608085870312156200029f57600080fd5b620002aa856200026b565b9350620002ba602086016200026b565b9250620002ca604086016200026b565b9150620002da606086016200026b565b905092959194509250565b600060208284031215620002f857600080fd5b815180151581146200030957600080fd5b9392505050565b60805160a05160c05160e051612128620003ac6000396000818161044c01528181610e5401526110b20152600081816103820152610b1b0152600081816106af015281816107020152818161090f01528181610c6901528181610d3001526112170152600081816106d40152818161075d0152818161089a01528181610a5d01528181610def01528181610fa301526111c001526121286000f3fe6080604052600436106100dd5760003560e01c806393798be71161007f578063f1387e1611610059578063f1387e16146102d6578063f237f1a8146102e9578063f2fde38b14610309578063f3fef3a31461032957600080fd5b806393798be714610255578063c541cbde14610293578063d09dc339146102c157600080fd5b80636b54d8a6116100bb5780636b54d8a6146101c757806379ba5097146101e75780637d75cc49146101fc5780638da5cb5b1461022057600080fd5b806301ffc9a7146100e2578063181f5a771461015957806369fd2b34146101a5575b600080fd5b3480156100ee57600080fd5b506101446100fd36600461167b565b7fffffffff00000000000000000000000000000000000000000000000000000000167ff1387e16000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b34801561016557600080fd5b50604080518082018252601081527f4665654d616e6167657220302e302e31000000000000000000000000000000006020820152905161015091906116c4565b3480156101b157600080fd5b506101c56101c0366004611730565b610349565b005b3480156101d357600080fd5b506101c56101e23660046117af565b6104bc565b3480156101f357600080fd5b506101c5610541565b34801561020857600080fd5b5061021260035481565b604051908152602001610150565b34801561022c57600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610150565b34801561026157600080fd5b506102126102703660046117ea565b600260209081526000938452604080852082529284528284209052825290205481565b34801561029f57600080fd5b506102b36102ae366004611961565b61063e565b604051610150929190611a02565b3480156102cd57600080fd5b50610212610a2c565b6101c56102e4366004611a56565b610ae2565b3480156102f557600080fd5b506101c5610304366004611ace565b611174565b34801561031557600080fd5b506101c5610324366004611b16565b61131b565b34801561033557600080fd5b506101c5610344366004611b33565b61132f565b60005473ffffffffffffffffffffffffffffffffffffffff163314806103a457503373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016145b61040f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4f6e6c79206f776e6572206f722070726f78790000000000000000000000000060448201526064015b60405180910390fd5b6040517f633b5f6e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063633b5f6e9061048590869086908690600401611b5f565b600060405180830381600087803b15801561049f57600080fd5b505af11580156104b3573d6000803e3d6000fd5b50505050505050565b6104c46114c9565b670de0b6b3a7640000811115610506576040517f05e8ac2900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60038190556040518181527f4fbcc2b7f4cb5518be923bd7c7c887e29e07b001e2a4a0fdd47c494696d8a1479060200160405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610406565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60408051808201909152600080825260208201526040805180820190915260008082526020820152604080518082019091526000808252602082015260408051808201909152600080825260208201528551610120106107005773ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811683527f00000000000000000000000000000000000000000000000000000000000000001681529092509050610a24565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16856000015173ffffffffffffffffffffffffffffffffffffffff16141580156107b057507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16856000015173ffffffffffffffffffffffffffffffffffffffff1614155b156107e7576040517ff861803000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000888060200190518101906108009190611c39565b63ffffffff169b5077ffffffffffffffffffffffffffffffffffffffffffffffff169b5077ffffffffffffffffffffffffffffffffffffffffffffffff169b5050505050505050505042811015610883576040517fb6c405f500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116808652602086018590528951909116036108f857835173ffffffffffffffffffffffffffffffffffffffff16855260208085015190860152610968565b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001685526003546109629061094a90670de0b6b3a7640000611d36565b6109549084611d49565b670de0b6b3a764000061154c565b60208601525b60006109738a611d86565b73ffffffffffffffffffffffffffffffffffffffff808d16600090815260026020908152604080832085845282528083208e519094168352928152919020549088015191925090670de0b6b3a7640000906109cf908390611d49565b6109d99190611dcb565b87602001516109e89190611e06565b876020018181525050610a048187602001516109549190611d49565b8660200151610a139190611e06565b602087015250949650929450505050505b935093915050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610ab9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610add9190611e19565b905090565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610b3d57503373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016145b610ba3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4f6e6c79206f776e6572206f722070726f7879000000000000000000000000006044820152606401610406565b3073ffffffffffffffffffffffffffffffffffffffff821603610bf2576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610c0083850185611ea0565b604080516020810190915260008152909250905081516101201015610c4c576000610c2d85870187611f6f565b9550505050505080806020019051810190610c489190612037565b9150505b600080610c5a33858561063e565b909250905060003415610dbd577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16836000015173ffffffffffffffffffffffffffffffffffffffff1614610cf0576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3483602001511115610d2e576040517fb2e532de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db084602001516040518263ffffffff1660e01b81526004016000604051808303818588803b158015610d9a57600080fd5b505af1158015610dae573d6000803e3d6000fd5b50505050508260200151340390505b6000610dc9888a612065565b60208501519091501561111e57835173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911603610ec4576040517f84afb76e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906384afb76e90610e8d9084908b9088906004016120a1565b600060405180830381600087803b158015610ea757600080fd5b505af1158015610ebb573d6000803e3d6000fd5b5050505061111e565b34600003610f7557835160208501516040517f23b872dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a8116600483015230602483015260448201929092529116906323b872dd906064016020604051808303816000875af1158015610f4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f7391906120f9565b505b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610fff573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110239190611e19565b836020015111156110755760208084015185820151604080519283529282015282917feb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7910160405180910390a261111e565b6040517f84afb76e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906384afb76e906110eb908490309088906004016120a1565b600060405180830381600087803b15801561110557600080fd5b505af1158015611119573d6000803e3d6000fd5b505050505b81156111695760405173ffffffffffffffffffffffffffffffffffffffff88169083156108fc029084906000818181858888f19350505050158015611167573d6000803e3d6000fd5b505b505050505050505050565b61117c6114c9565b670de0b6b3a76400008111156111be576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415801561126657507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b1561129d576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84811660008181526002602090815260408083208884528252808320948716808452948252918290208590558151938452830184905285927f41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84910160405180910390a350505050565b6113236114c9565b61132c81611586565b50565b6113376114c9565b73ffffffffffffffffffffffffffffffffffffffff821661139d576000805460405173ffffffffffffffffffffffffffffffffffffffff9091169183156108fc02918491818181858888f19350505050158015611398573d6000803e3d6000fd5b505050565b8173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6113d860005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602481018490526044016020604051808303816000875af115801561144a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061146e91906120f9565b506040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201529081018290527f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9060600160405180910390a15050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461154a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610406565b565b6000821561157a5781611560600185611e06565b61156a9190611dcb565b611575906001611d36565b61157d565b60005b90505b92915050565b3373ffffffffffffffffffffffffffffffffffffffff821603611605576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610406565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561168d57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146116bd57600080fd5b9392505050565b600060208083528351808285015260005b818110156116f1578581018301518582016040015282016116d5565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60008060006040848603121561174557600080fd5b83359250602084013567ffffffffffffffff8082111561176457600080fd5b818601915086601f83011261177857600080fd5b81358181111561178757600080fd5b8760208260061b850101111561179c57600080fd5b6020830194508093505050509250925092565b6000602082840312156117c157600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461132c57600080fd5b6000806000606084860312156117ff57600080fd5b833561180a816117c8565b9250602084013591506040840135611821816117c8565b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810167ffffffffffffffff8111828210171561187e5761187e61182c565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156118cb576118cb61182c565b604052919050565b600082601f8301126118e457600080fd5b813567ffffffffffffffff8111156118fe576118fe61182c565b61192f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611884565b81815284602083860101111561194457600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000838503606081121561197757600080fd5b8435611982816117c8565b9350602085013567ffffffffffffffff81111561199e57600080fd5b6119aa878288016118d3565b93505060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0820112156119dd57600080fd5b506119e661185b565b60408501356119f4816117c8565b815292959194509192509050565b825173ffffffffffffffffffffffffffffffffffffffff1681526020808401519082015260808101825173ffffffffffffffffffffffffffffffffffffffff166040830152602083015160608301526116bd565b600080600060408486031215611a6b57600080fd5b833567ffffffffffffffff80821115611a8357600080fd5b818601915086601f830112611a9757600080fd5b813581811115611aa657600080fd5b876020828501011115611ab857600080fd5b60209283019550935050840135611821816117c8565b60008060008060808587031215611ae457600080fd5b8435611aef816117c8565b9350602085013592506040850135611b06816117c8565b9396929550929360600135925050565b600060208284031215611b2857600080fd5b81356116bd816117c8565b60008060408385031215611b4657600080fd5b8235611b51816117c8565b946020939093013593505050565b8381526040602080830182905282820184905260009190859060608501845b87811015611bc1578335611b91816117c8565b73ffffffffffffffffffffffffffffffffffffffff16825283830135838301529284019290840190600101611b7e565b5098975050505050505050565b805163ffffffff81168114611be257600080fd5b919050565b8051601781900b8114611be257600080fd5b805167ffffffffffffffff81168114611be257600080fd5b805177ffffffffffffffffffffffffffffffffffffffffffffffff81168114611be257600080fd5b6000806000806000806000806000806000806101808d8f031215611c5c57600080fd5b8c519b50611c6c60208e01611bce565b9a50611c7a60408e01611be7565b9950611c8860608e01611be7565b9850611c9660808e01611be7565b9750611ca460a08e01611bf9565b965060c08d01519550611cb960e08e01611bf9565b9450611cc86101008e01611bf9565b9350611cd76101208e01611c11565b9250611ce66101408e01611c11565b9150611cf56101608e01611bce565b90509295989b509295989b509295989b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561158057611580611d07565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611d8157611d81611d07565b500290565b80516020808301519190811015611dc5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b600082611e01577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b8181038181111561158057611580611d07565b600060208284031215611e2b57600080fd5b5051919050565b600082601f830112611e4357600080fd5b6040516060810181811067ffffffffffffffff82111715611e6657611e6661182c565b604052806060840185811115611e7b57600080fd5b845b81811015611e95578035835260209283019201611e7d565b509195945050505050565b60008060808385031215611eb357600080fd5b611ebd8484611e32565b9150606083013567ffffffffffffffff811115611ed957600080fd5b611ee5858286016118d3565b9150509250929050565b600082601f830112611f0057600080fd5b8135602067ffffffffffffffff821115611f1c57611f1c61182c565b8160051b611f2b828201611884565b9283528481018201928281019087851115611f4557600080fd5b83870192505b84831015611f6457823582529183019190830190611f4b565b979650505050505050565b6000806000806000806101008789031215611f8957600080fd5b611f938888611e32565b9550606087013567ffffffffffffffff80821115611fb057600080fd5b611fbc8a838b016118d3565b96506080890135915080821115611fd257600080fd5b611fde8a838b01611eef565b955060a0890135915080821115611ff457600080fd5b6120008a838b01611eef565b945060c0890135935060e089013591508082111561201d57600080fd5b5061202a89828a016118d3565b9150509295509295509295565b60006020828403121561204957600080fd5b61205161185b565b825161205c816117c8565b81529392505050565b80356020831015611580577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b83815273ffffffffffffffffffffffffffffffffffffffff83166020820152608081016120f16040830184805173ffffffffffffffffffffffffffffffffffffffff168252602090810151910152565b949350505050565b60006020828403121561210b57600080fd5b815180151581146116bd57600080fdfea164736f6c6343000810000a",
}

var LLOFeeManagerABI = LLOFeeManagerMetaData.ABI

var LLOFeeManagerBin = LLOFeeManagerMetaData.Bin

func DeployLLOFeeManager(auth *bind.TransactOpts, backend bind.ContractBackend, _linkAddress common.Address, _nativeAddress common.Address, _proxyAddress common.Address, _rewardManagerAddress common.Address) (common.Address, *types.Transaction, *LLOFeeManager, error) {
	parsed, err := LLOFeeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LLOFeeManagerBin), backend, _linkAddress, _nativeAddress, _proxyAddress, _rewardManagerAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LLOFeeManager{LLOFeeManagerCaller: LLOFeeManagerCaller{contract: contract}, LLOFeeManagerTransactor: LLOFeeManagerTransactor{contract: contract}, LLOFeeManagerFilterer: LLOFeeManagerFilterer{contract: contract}}, nil
}

type LLOFeeManager struct {
	address common.Address
	abi     abi.ABI
	LLOFeeManagerCaller
	LLOFeeManagerTransactor
	LLOFeeManagerFilterer
}

type LLOFeeManagerCaller struct {
	contract *bind.BoundContract
}

type LLOFeeManagerTransactor struct {
	contract *bind.BoundContract
}

type LLOFeeManagerFilterer struct {
	contract *bind.BoundContract
}

type LLOFeeManagerSession struct {
	Contract     *LLOFeeManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LLOFeeManagerCallerSession struct {
	Contract *LLOFeeManagerCaller
	CallOpts bind.CallOpts
}

type LLOFeeManagerTransactorSession struct {
	Contract     *LLOFeeManagerTransactor
	TransactOpts bind.TransactOpts
}

type LLOFeeManagerRaw struct {
	Contract *LLOFeeManager
}

type LLOFeeManagerCallerRaw struct {
	Contract *LLOFeeManagerCaller
}

type LLOFeeManagerTransactorRaw struct {
	Contract *LLOFeeManagerTransactor
}

func NewLLOFeeManager(address common.Address, backend bind.ContractBackend) (*LLOFeeManager, error) {
	abi, err := abi.JSON(strings.NewReader(LLOFeeManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLLOFeeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManager{address: address, abi: abi, LLOFeeManagerCaller: LLOFeeManagerCaller{contract: contract}, LLOFeeManagerTransactor: LLOFeeManagerTransactor{contract: contract}, LLOFeeManagerFilterer: LLOFeeManagerFilterer{contract: contract}}, nil
}

func NewLLOFeeManagerCaller(address common.Address, caller bind.ContractCaller) (*LLOFeeManagerCaller, error) {
	contract, err := bindLLOFeeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerCaller{contract: contract}, nil
}

func NewLLOFeeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*LLOFeeManagerTransactor, error) {
	contract, err := bindLLOFeeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerTransactor{contract: contract}, nil
}

func NewLLOFeeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*LLOFeeManagerFilterer, error) {
	contract, err := bindLLOFeeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerFilterer{contract: contract}, nil
}

func bindLLOFeeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LLOFeeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LLOFeeManager *LLOFeeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LLOFeeManager.Contract.LLOFeeManagerCaller.contract.Call(opts, result, method, params...)
}

func (_LLOFeeManager *LLOFeeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.LLOFeeManagerTransactor.contract.Transfer(opts)
}

func (_LLOFeeManager *LLOFeeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.LLOFeeManagerTransactor.contract.Transact(opts, method, params...)
}

func (_LLOFeeManager *LLOFeeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LLOFeeManager.Contract.contract.Call(opts, result, method, params...)
}

func (_LLOFeeManager *LLOFeeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.contract.Transfer(opts)
}

func (_LLOFeeManager *LLOFeeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.contract.Transact(opts, method, params...)
}

func (_LLOFeeManager *LLOFeeManagerCaller) GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "getFeeAndReward", subscriber, report, quote)

	if err != nil {
		return *new(CommonAsset), *new(CommonAsset), err
	}

	out0 := *abi.ConvertType(out[0], new(CommonAsset)).(*CommonAsset)
	out1 := *abi.ConvertType(out[1], new(CommonAsset)).(*CommonAsset)

	return out0, out1, err

}

func (_LLOFeeManager *LLOFeeManagerSession) GetFeeAndReward(subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	return _LLOFeeManager.Contract.GetFeeAndReward(&_LLOFeeManager.CallOpts, subscriber, report, quote)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) GetFeeAndReward(subscriber common.Address, report []byte, quote IFeeManagerQuote) (CommonAsset, CommonAsset, error) {
	return _LLOFeeManager.Contract.GetFeeAndReward(&_LLOFeeManager.CallOpts, subscriber, report, quote)
}

func (_LLOFeeManager *LLOFeeManagerCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LLOFeeManager *LLOFeeManagerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _LLOFeeManager.Contract.LinkAvailableForPayment(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _LLOFeeManager.Contract.LinkAvailableForPayment(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCaller) NativeSurcharge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "nativeSurcharge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LLOFeeManager *LLOFeeManagerSession) NativeSurcharge() (*big.Int, error) {
	return _LLOFeeManager.Contract.NativeSurcharge(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) NativeSurcharge() (*big.Int, error) {
	return _LLOFeeManager.Contract.NativeSurcharge(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LLOFeeManager *LLOFeeManagerSession) Owner() (common.Address, error) {
	return _LLOFeeManager.Contract.Owner(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) Owner() (common.Address, error) {
	return _LLOFeeManager.Contract.Owner(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCaller) SubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "subscriberDiscounts", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LLOFeeManager *LLOFeeManagerSession) SubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _LLOFeeManager.Contract.SubscriberDiscounts(&_LLOFeeManager.CallOpts, arg0, arg1, arg2)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) SubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _LLOFeeManager.Contract.SubscriberDiscounts(&_LLOFeeManager.CallOpts, arg0, arg1, arg2)
}

func (_LLOFeeManager *LLOFeeManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LLOFeeManager *LLOFeeManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LLOFeeManager.Contract.SupportsInterface(&_LLOFeeManager.CallOpts, interfaceId)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LLOFeeManager.Contract.SupportsInterface(&_LLOFeeManager.CallOpts, interfaceId)
}

func (_LLOFeeManager *LLOFeeManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LLOFeeManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LLOFeeManager *LLOFeeManagerSession) TypeAndVersion() (string, error) {
	return _LLOFeeManager.Contract.TypeAndVersion(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerCallerSession) TypeAndVersion() (string, error) {
	return _LLOFeeManager.Contract.TypeAndVersion(&_LLOFeeManager.CallOpts)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "acceptOwnership")
}

func (_LLOFeeManager *LLOFeeManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _LLOFeeManager.Contract.AcceptOwnership(&_LLOFeeManager.TransactOpts)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LLOFeeManager.Contract.AcceptOwnership(&_LLOFeeManager.TransactOpts)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) ProcessFee(opts *bind.TransactOpts, payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "processFee", payload, subscriber)
}

func (_LLOFeeManager *LLOFeeManagerSession) ProcessFee(payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.ProcessFee(&_LLOFeeManager.TransactOpts, payload, subscriber)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) ProcessFee(payload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.ProcessFee(&_LLOFeeManager.TransactOpts, payload, subscriber)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "setFeeRecipients", configDigest, rewardRecipientAndWeights)
}

func (_LLOFeeManager *LLOFeeManagerSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.SetFeeRecipients(&_LLOFeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.SetFeeRecipients(&_LLOFeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) SetNativeSurcharge(opts *bind.TransactOpts, surcharge *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "setNativeSurcharge", surcharge)
}

func (_LLOFeeManager *LLOFeeManagerSession) SetNativeSurcharge(surcharge *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.SetNativeSurcharge(&_LLOFeeManager.TransactOpts, surcharge)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) SetNativeSurcharge(surcharge *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.SetNativeSurcharge(&_LLOFeeManager.TransactOpts, surcharge)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "transferOwnership", to)
}

func (_LLOFeeManager *LLOFeeManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.TransferOwnership(&_LLOFeeManager.TransactOpts, to)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.TransferOwnership(&_LLOFeeManager.TransactOpts, to)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "updateSubscriberDiscount", subscriber, feedId, token, discount)
}

func (_LLOFeeManager *LLOFeeManagerSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.UpdateSubscriberDiscount(&_LLOFeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.UpdateSubscriberDiscount(&_LLOFeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_LLOFeeManager *LLOFeeManagerTransactor) Withdraw(opts *bind.TransactOpts, assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.contract.Transact(opts, "withdraw", assetAddress, quantity)
}

func (_LLOFeeManager *LLOFeeManagerSession) Withdraw(assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.Withdraw(&_LLOFeeManager.TransactOpts, assetAddress, quantity)
}

func (_LLOFeeManager *LLOFeeManagerTransactorSession) Withdraw(assetAddress common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _LLOFeeManager.Contract.Withdraw(&_LLOFeeManager.TransactOpts, assetAddress, quantity)
}

type LLOFeeManagerInsufficientLinkIterator struct {
	Event *LLOFeeManagerInsufficientLink

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOFeeManagerInsufficientLinkIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOFeeManagerInsufficientLink)
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
		it.Event = new(LLOFeeManagerInsufficientLink)
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

func (it *LLOFeeManagerInsufficientLinkIterator) Error() error {
	return it.fail
}

func (it *LLOFeeManagerInsufficientLinkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOFeeManagerInsufficientLink struct {
	ConfigDigest   [32]byte
	LinkQuantity   *big.Int
	NativeQuantity *big.Int
	Raw            types.Log
}

func (_LLOFeeManager *LLOFeeManagerFilterer) FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*LLOFeeManagerInsufficientLinkIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _LLOFeeManager.contract.FilterLogs(opts, "InsufficientLink", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerInsufficientLinkIterator{contract: _LLOFeeManager.contract, event: "InsufficientLink", logs: logs, sub: sub}, nil
}

func (_LLOFeeManager *LLOFeeManagerFilterer) WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _LLOFeeManager.contract.WatchLogs(opts, "InsufficientLink", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOFeeManagerInsufficientLink)
				if err := _LLOFeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
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

func (_LLOFeeManager *LLOFeeManagerFilterer) ParseInsufficientLink(log types.Log) (*LLOFeeManagerInsufficientLink, error) {
	event := new(LLOFeeManagerInsufficientLink)
	if err := _LLOFeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOFeeManagerNativeSurchargeSetIterator struct {
	Event *LLOFeeManagerNativeSurchargeSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOFeeManagerNativeSurchargeSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOFeeManagerNativeSurchargeSet)
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
		it.Event = new(LLOFeeManagerNativeSurchargeSet)
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

func (it *LLOFeeManagerNativeSurchargeSetIterator) Error() error {
	return it.fail
}

func (it *LLOFeeManagerNativeSurchargeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOFeeManagerNativeSurchargeSet struct {
	NewSurcharge *big.Int
	Raw          types.Log
}

func (_LLOFeeManager *LLOFeeManagerFilterer) FilterNativeSurchargeSet(opts *bind.FilterOpts) (*LLOFeeManagerNativeSurchargeSetIterator, error) {

	logs, sub, err := _LLOFeeManager.contract.FilterLogs(opts, "NativeSurchargeSet")
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerNativeSurchargeSetIterator{contract: _LLOFeeManager.contract, event: "NativeSurchargeSet", logs: logs, sub: sub}, nil
}

func (_LLOFeeManager *LLOFeeManagerFilterer) WatchNativeSurchargeSet(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerNativeSurchargeSet) (event.Subscription, error) {

	logs, sub, err := _LLOFeeManager.contract.WatchLogs(opts, "NativeSurchargeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOFeeManagerNativeSurchargeSet)
				if err := _LLOFeeManager.contract.UnpackLog(event, "NativeSurchargeSet", log); err != nil {
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

func (_LLOFeeManager *LLOFeeManagerFilterer) ParseNativeSurchargeSet(log types.Log) (*LLOFeeManagerNativeSurchargeSet, error) {
	event := new(LLOFeeManagerNativeSurchargeSet)
	if err := _LLOFeeManager.contract.UnpackLog(event, "NativeSurchargeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOFeeManagerOwnershipTransferRequestedIterator struct {
	Event *LLOFeeManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOFeeManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOFeeManagerOwnershipTransferRequested)
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
		it.Event = new(LLOFeeManagerOwnershipTransferRequested)
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

func (it *LLOFeeManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LLOFeeManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOFeeManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LLOFeeManager *LLOFeeManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOFeeManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOFeeManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerOwnershipTransferRequestedIterator{contract: _LLOFeeManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LLOFeeManager *LLOFeeManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOFeeManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOFeeManagerOwnershipTransferRequested)
				if err := _LLOFeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LLOFeeManager *LLOFeeManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*LLOFeeManagerOwnershipTransferRequested, error) {
	event := new(LLOFeeManagerOwnershipTransferRequested)
	if err := _LLOFeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOFeeManagerOwnershipTransferredIterator struct {
	Event *LLOFeeManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOFeeManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOFeeManagerOwnershipTransferred)
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
		it.Event = new(LLOFeeManagerOwnershipTransferred)
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

func (it *LLOFeeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LLOFeeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOFeeManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LLOFeeManager *LLOFeeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOFeeManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOFeeManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerOwnershipTransferredIterator{contract: _LLOFeeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LLOFeeManager *LLOFeeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOFeeManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOFeeManagerOwnershipTransferred)
				if err := _LLOFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LLOFeeManager *LLOFeeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*LLOFeeManagerOwnershipTransferred, error) {
	event := new(LLOFeeManagerOwnershipTransferred)
	if err := _LLOFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOFeeManagerSubscriberDiscountUpdatedIterator struct {
	Event *LLOFeeManagerSubscriberDiscountUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOFeeManagerSubscriberDiscountUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOFeeManagerSubscriberDiscountUpdated)
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
		it.Event = new(LLOFeeManagerSubscriberDiscountUpdated)
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

func (it *LLOFeeManagerSubscriberDiscountUpdatedIterator) Error() error {
	return it.fail
}

func (it *LLOFeeManagerSubscriberDiscountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOFeeManagerSubscriberDiscountUpdated struct {
	Subscriber common.Address
	FeedId     [32]byte
	Token      common.Address
	Discount   *big.Int
	Raw        types.Log
}

func (_LLOFeeManager *LLOFeeManagerFilterer) FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*LLOFeeManagerSubscriberDiscountUpdatedIterator, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _LLOFeeManager.contract.FilterLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerSubscriberDiscountUpdatedIterator{contract: _LLOFeeManager.contract, event: "SubscriberDiscountUpdated", logs: logs, sub: sub}, nil
}

func (_LLOFeeManager *LLOFeeManagerFilterer) WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _LLOFeeManager.contract.WatchLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOFeeManagerSubscriberDiscountUpdated)
				if err := _LLOFeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
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

func (_LLOFeeManager *LLOFeeManagerFilterer) ParseSubscriberDiscountUpdated(log types.Log) (*LLOFeeManagerSubscriberDiscountUpdated, error) {
	event := new(LLOFeeManagerSubscriberDiscountUpdated)
	if err := _LLOFeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOFeeManagerWithdrawIterator struct {
	Event *LLOFeeManagerWithdraw

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOFeeManagerWithdrawIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOFeeManagerWithdraw)
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
		it.Event = new(LLOFeeManagerWithdraw)
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

func (it *LLOFeeManagerWithdrawIterator) Error() error {
	return it.fail
}

func (it *LLOFeeManagerWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOFeeManagerWithdraw struct {
	AdminAddress common.Address
	AssetAddress common.Address
	Quantity     *big.Int
	Raw          types.Log
}

func (_LLOFeeManager *LLOFeeManagerFilterer) FilterWithdraw(opts *bind.FilterOpts) (*LLOFeeManagerWithdrawIterator, error) {

	logs, sub, err := _LLOFeeManager.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &LLOFeeManagerWithdrawIterator{contract: _LLOFeeManager.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

func (_LLOFeeManager *LLOFeeManagerFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerWithdraw) (event.Subscription, error) {

	logs, sub, err := _LLOFeeManager.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOFeeManagerWithdraw)
				if err := _LLOFeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

func (_LLOFeeManager *LLOFeeManagerFilterer) ParseWithdraw(log types.Log) (*LLOFeeManagerWithdraw, error) {
	event := new(LLOFeeManagerWithdraw)
	if err := _LLOFeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LLOFeeManager *LLOFeeManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LLOFeeManager.abi.Events["InsufficientLink"].ID:
		return _LLOFeeManager.ParseInsufficientLink(log)
	case _LLOFeeManager.abi.Events["NativeSurchargeSet"].ID:
		return _LLOFeeManager.ParseNativeSurchargeSet(log)
	case _LLOFeeManager.abi.Events["OwnershipTransferRequested"].ID:
		return _LLOFeeManager.ParseOwnershipTransferRequested(log)
	case _LLOFeeManager.abi.Events["OwnershipTransferred"].ID:
		return _LLOFeeManager.ParseOwnershipTransferred(log)
	case _LLOFeeManager.abi.Events["SubscriberDiscountUpdated"].ID:
		return _LLOFeeManager.ParseSubscriberDiscountUpdated(log)
	case _LLOFeeManager.abi.Events["Withdraw"].ID:
		return _LLOFeeManager.ParseWithdraw(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LLOFeeManagerInsufficientLink) Topic() common.Hash {
	return common.HexToHash("0xeb6f22018570d97db6df12dc94f202b4e2b2888a6a5d4bd179422c91b29dcdf7")
}

func (LLOFeeManagerNativeSurchargeSet) Topic() common.Hash {
	return common.HexToHash("0x4fbcc2b7f4cb5518be923bd7c7c887e29e07b001e2a4a0fdd47c494696d8a147")
}

func (LLOFeeManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LLOFeeManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LLOFeeManagerSubscriberDiscountUpdated) Topic() common.Hash {
	return common.HexToHash("0x41eb9ccd292d5906dc1f0ec108bed3e2b966e3071e033df938f7215f6d30ca84")
}

func (LLOFeeManagerWithdraw) Topic() common.Hash {
	return common.HexToHash("0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb")
}

func (_LLOFeeManager *LLOFeeManager) Address() common.Address {
	return _LLOFeeManager.address
}

type LLOFeeManagerInterface interface {
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

	FilterInsufficientLink(opts *bind.FilterOpts, configDigest [][32]byte) (*LLOFeeManagerInsufficientLinkIterator, error)

	WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerInsufficientLink, configDigest [][32]byte) (event.Subscription, error)

	ParseInsufficientLink(log types.Log) (*LLOFeeManagerInsufficientLink, error)

	FilterNativeSurchargeSet(opts *bind.FilterOpts) (*LLOFeeManagerNativeSurchargeSetIterator, error)

	WatchNativeSurchargeSet(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerNativeSurchargeSet) (event.Subscription, error)

	ParseNativeSurchargeSet(log types.Log) (*LLOFeeManagerNativeSurchargeSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOFeeManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LLOFeeManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOFeeManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LLOFeeManagerOwnershipTransferred, error)

	FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*LLOFeeManagerSubscriberDiscountUpdatedIterator, error)

	WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error)

	ParseSubscriberDiscountUpdated(log types.Log) (*LLOFeeManagerSubscriberDiscountUpdated, error)

	FilterWithdraw(opts *bind.FilterOpts) (*LLOFeeManagerWithdrawIterator, error)

	WatchWithdraw(opts *bind.WatchOpts, sink chan<- *LLOFeeManagerWithdraw) (event.Subscription, error)

	ParseWithdraw(log types.Log) (*LLOFeeManagerWithdraw, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
