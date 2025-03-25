// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_fee_manager_v0_5_0

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

type IRewardManagerFeePayment struct {
	PoolId [32]byte
	Amount *big.Int
}

var MockFeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_linkAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_nativeAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_proxyAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rewardManagerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getFeeAndReward\",\"inputs\":[{\"name\":\"subscriber\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quoteAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCommon.Asset\",\"components\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCommon.Asset\",\"components\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"i_linkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"i_nativeAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"i_proxyAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"i_rewardManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRewardManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"linkAvailableForPayment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payLinkDeficit\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"processFee\",\"inputs\":[{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"parameterPayload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"subscriber\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"processFeeBulk\",\"inputs\":[{\"name\":\"payloads\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"parameterPayload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"subscriber\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"s_globalDiscounts\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_linkDeficit\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_nativeSurcharge\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_subscriberDiscounts\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setFeeRecipients\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\",\"internalType\":\"structCommon.AddressAndWeight[]\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"weight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNativeSurcharge\",\"inputs\":[{\"name\":\"surcharge\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateSubscriberDiscount\",\"inputs\":[{\"name\":\"subscriber\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"discount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSubscriberGlobalDiscount\",\"inputs\":[{\"name\":\"subscriber\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"discount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quantity\",\"type\":\"uint192\",\"internalType\":\"uint192\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DiscountApplied\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"subscriber\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCommon.Asset\",\"components\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"reward\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCommon.Asset\",\"components\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"appliedDiscount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientLink\",\"inputs\":[{\"name\":\"rewards\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIRewardManager.FeePayment[]\",\"components\":[{\"name\":\"poolId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint192\",\"internalType\":\"uint192\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LinkDeficitCleared\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"linkQuantity\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NativeSurchargeUpdated\",\"inputs\":[{\"name\":\"newSurcharge\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubscriberDiscountUpdated\",\"inputs\":[{\"name\":\"subscriber\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"discount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"adminAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"assetAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quantity\",\"type\":\"uint192\",\"indexed\":false,\"internalType\":\"uint192\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ExpiredReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDiscount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidQuote\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceivingAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSurcharge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroDeficit\",\"inputs\":[]}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200259638038062002596833981016040819052620000359162000288565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001c0565b5050506001600160a01b0384161580620000e057506001600160a01b038316155b80620000f357506001600160a01b038216155b806200010657506001600160a01b038116155b15620001255760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03848116608081905284821660a05283821660c05290821660e081905260405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156200018f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b59190620002e5565b505050505062000310565b336001600160a01b038216036200021a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200028357600080fd5b919050565b600080600080608085870312156200029f57600080fd5b620002aa856200026b565b9350620002ba602086016200026b565b9250620002ca604086016200026b565b9150620002da606086016200026b565b905092959194509250565b600060208284031215620002f857600080fd5b815180151581146200030957600080fd5b9392505050565b60805160a05160c05160e0516121f86200039e600039600081816102c301526111d30152600081816103830152818161079301528181610ee4015261111901526000818161033c01528181610a5001528181610d070152610faf0152600081816105280152818161090e015281816109f901528181610cb001528181610e47015261100601526121f86000f3fe6080604052600436106101805760003560e01c806376cf3187116100d6578063dba45fe01161007f578063ea4b861b11610059578063ea4b861b14610516578063f2fde38b1461054a578063f65df9621461056a57600080fd5b8063dba45fe014610478578063e03dab1a1461048b578063e389d9a4146104f657600080fd5b80638da5cb5b116100b05780638da5cb5b14610418578063ce7817d114610443578063d09dc3391461046357600080fd5b806376cf3187146103a557806379ba5097146103c557806387d6d843146103da57600080fd5b806332f5f746116101385780636387866811610112578063638786681461032a5780636c2f1a171461035e5780636d1342cb1461037157600080fd5b806332f5f7461461029b5780633aa5ac07146102b1578063505380941461030a57600080fd5b8063181f5a7711610169578063181f5a77146101f55780631cc7f2d8146102415780631d4d84a21461027957600080fd5b8063013f542b1461018557806301ffc9a7146101c5575b600080fd5b34801561019157600080fd5b506101b26101a0366004611905565b60046020526000908152604090205481565b6040519081526020015b60405180910390f35b3480156101d157600080fd5b506101e56101e036600461191e565b61058a565b60405190151581526020016101bc565b34801561020157600080fd5b50604080518082018252601081527f4665654d616e6167657220322e312e3000000000000000000000000000000000602082015290516101bc919061198b565b34801561024d57600080fd5b506101b261025c366004611a0e565b600360209081526000928352604080842090915290825290205481565b34801561028557600080fd5b50610299610294366004611a47565b610623565b005b3480156102a757600080fd5b506101b260055481565b3480156102bd57600080fd5b506102e57f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101bc565b34801561031657600080fd5b50610299610325366004611ac7565b6106e1565b34801561033657600080fd5b506102e57f000000000000000000000000000000000000000000000000000000000000000081565b61029961036c366004611b2b565b61077b565b34801561037d57600080fd5b506102e57f000000000000000000000000000000000000000000000000000000000000000081565b3480156103b157600080fd5b506102996103c0366004611bda565b6109a3565b3480156103d157600080fd5b50610299610b58565b3480156103e657600080fd5b506101b26103f5366004611c21565b600260209081526000938452604080852082529284528284209052825290205481565b34801561042457600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166102e5565b34801561044f57600080fd5b5061029961045e366004611c58565b610c5a565b34801561046f57600080fd5b506101b2610e16565b610299610486366004611ca9565b610ecc565b34801561049757600080fd5b506104ab6104a6366004611e07565b610f5a565b60408051845173ffffffffffffffffffffffffffffffffffffffff9081168252602095860151868301528451169181019190915292909101516060830152608082015260a0016101bc565b34801561050257600080fd5b50610299610511366004611905565b61109c565b34801561052257600080fd5b506102e57f000000000000000000000000000000000000000000000000000000000000000081565b34801561055657600080fd5b50610299610565366004611e60565b6110ed565b34801561057657600080fd5b50610299610585366004611e7d565b611101565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fdba45fe000000000000000000000000000000000000000000000000000000000148061061d57507fffffffff0000000000000000000000000000000000000000000000000000000082167f6c2f1a1700000000000000000000000000000000000000000000000000000000145b92915050565b61062b611243565b61066673ffffffffffffffffffffffffffffffffffffffff84168377ffffffffffffffffffffffffffffffffffffffffffffffff84166112c6565b6040805133815273ffffffffffffffffffffffffffffffffffffffff848116602083015285168183015277ffffffffffffffffffffffffffffffffffffffffffffffff8316606082015290517f7ff78a71698bdb18dcca96f52ab25e0a1b146fb6a49adf8e6845299e49021f299181900360800190a1505050565b6106e9611243565b670de0b6b3a764000067ffffffffffffffff82161115610735576040517f05e8ac2900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660058190556040519081527f08f7c0d17932ddb8523bc06754d42ff19ebc77d76a8b9bfde02c28ab1ed3d6399060200160405180910390a150565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146107ea576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008467ffffffffffffffff81111561080557610805611d2d565b60405190808252806020026020018201604052801561083e57816020015b61082b611878565b8152602001906001900390816108235790505b5090506000806000805b888110156109895760008060006108848d8d8681811061086a5761086a611efc565b905060200281019061087c9190611f2b565b8d8d8d611358565b92509250925082602001516000146109755760405180608001604052808e8e878181106108b3576108b3611efc565b90506020028101906108c59190611f2b565b6108ce91611f90565b8152602001848152602001838152602001828152508886806108ef90611fcc565b97508151811061090157610901611efc565b60200260200101819052507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16836000015173ffffffffffffffffffffffffffffffffffffffff160361096e57866001019650610975565b8560010195505b5050508061098290611fcc565b9050610848565b508215158061099757508115155b50505050505050505050565b6109ab611243565b670de0b6b3a764000067ffffffffffffffff821611156109f7576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614158015610a9f57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b15610ad6576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff838116600081815260036020908152604080832094871680845294825280832067ffffffffffffffff87169081905581519586529185019190915290927f5eba5a8afa39780f0f99b6cbeb95f3da6a7040ca00abd46bdc91a0a060134139910160405180910390a3505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610bde576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610c62611243565b670de0b6b3a764000067ffffffffffffffff82161115610cae576040517f997ea36000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614158015610d5657507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b15610d8d576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff848116600081815260026020908152604080832088845282528083209487168084529482529182902067ffffffffffffffff86169081905582519485529084015285927f5eba5a8afa39780f0f99b6cbeb95f3da6a7040ca00abd46bdc91a0a060134139910160405180910390a350505050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610ea3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ec7919061202b565b905090565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610f3b576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000610f4d8888888888611358565b5050505050505050505050565b60408051808201909152600080825260208201526040805180820190915260008082526020820152604080518082019091526000808252602082018190529060408051808201909152600080825260208201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161415801561105557507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614155b1561108c576040517ff861803000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9097909650600095509350505050565b6110a4611243565b60008181526004602090815260409182902054915182815283917f843f0b103e50b42b08f9d30f12f961845a6d02623730872e24644899c0dd9895910160405180910390a25050565b6110f5611243565b6110fe81611468565b50565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161480159061115f575060005473ffffffffffffffffffffffffffffffffffffffff163314155b15611196576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f14060f2300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906314060f239061120c90869086908690600401612044565b600060405180830381600087803b15801561122657600080fd5b505af115801561123a573d6000803e3d6000fd5b50505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146112c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bd5565b565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb0000000000000000000000000000000000000000000000000000000017905261135390849061155d565b505050565b6040805180820190915260008082526020820152604080518082019091526000808252602082015260003073ffffffffffffffffffffffffffffffffffffffff8516036113d1576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006113df888a018a6120c4565b9150506000816113ee90612168565b905060007e010000000000000000000000000000000000000000000000000000000000007fffff00000000000000000000000000000000000000000000000000000000000083161461144957611446888a018a611e60565b90505b611454878483610f5a565b955095509550505050955095509592505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036114e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bd5565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006115bf826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166116699092919063ffffffff16565b80519091501561135357808060200190518101906115dd91906121ad565b611353576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610bd5565b60606116788484600085611680565b949350505050565b606082471015611712576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610bd5565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161173b91906121cf565b60006040518083038185875af1925050503d8060008114611778576040519150601f19603f3d011682016040523d82523d6000602084013e61177d565b606091505b509150915061178e87838387611799565b979650505050505050565b6060831561182f5782516000036118285773ffffffffffffffffffffffffffffffffffffffff85163b611828576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610bd5565b5081611678565b61167883838151156118445781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bd5919061198b565b6040518060800160405280600080191681526020016118c06040518060400160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600081525090565b81526020016118f86040518060400160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600081525090565b8152602001600081525090565b60006020828403121561191757600080fd5b5035919050565b60006020828403121561193057600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461196057600080fd5b9392505050565b60005b8381101561198257818101518382015260200161196a565b50506000910152565b60208152600082518060208401526119aa816040850160208701611967565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b73ffffffffffffffffffffffffffffffffffffffff811681146110fe57600080fd5b8035611a09816119dc565b919050565b60008060408385031215611a2157600080fd5b8235611a2c816119dc565b91506020830135611a3c816119dc565b809150509250929050565b600080600060608486031215611a5c57600080fd5b8335611a67816119dc565b92506020840135611a77816119dc565b9150604084013577ffffffffffffffffffffffffffffffffffffffffffffffff81168114611aa457600080fd5b809150509250925092565b803567ffffffffffffffff81168114611a0957600080fd5b600060208284031215611ad957600080fd5b61196082611aaf565b60008083601f840112611af457600080fd5b50813567ffffffffffffffff811115611b0c57600080fd5b602083019150836020828501011115611b2457600080fd5b9250929050565b600080600080600060608688031215611b4357600080fd5b853567ffffffffffffffff80821115611b5b57600080fd5b818801915088601f830112611b6f57600080fd5b813581811115611b7e57600080fd5b8960208260051b8501011115611b9357600080fd5b602092830197509550908701359080821115611bae57600080fd5b50611bbb88828901611ae2565b9094509250611bce9050604087016119fe565b90509295509295909350565b600080600060608486031215611bef57600080fd5b8335611bfa816119dc565b92506020840135611c0a816119dc565b9150611c1860408501611aaf565b90509250925092565b600080600060608486031215611c3657600080fd5b8335611c41816119dc565b9250602084013591506040840135611aa4816119dc565b60008060008060808587031215611c6e57600080fd5b8435611c79816119dc565b9350602085013592506040850135611c90816119dc565b9150611c9e60608601611aaf565b905092959194509250565b600080600080600060608688031215611cc157600080fd5b853567ffffffffffffffff80821115611cd957600080fd5b611ce589838a01611ae2565b90975095506020880135915080821115611cfe57600080fd5b50611d0b88828901611ae2565b9094509250506040860135611d1f816119dc565b809150509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112611d6d57600080fd5b813567ffffffffffffffff80821115611d8857611d88611d2d565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715611dce57611dce611d2d565b81604052838152866020858801011115611de757600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600060608486031215611e1c57600080fd5b8335611e27816119dc565b9250602084013567ffffffffffffffff811115611e4357600080fd5b611e4f86828701611d5c565b9250506040840135611aa4816119dc565b600060208284031215611e7257600080fd5b8135611960816119dc565b600080600060408486031215611e9257600080fd5b83359250602084013567ffffffffffffffff80821115611eb157600080fd5b818601915086601f830112611ec557600080fd5b813581811115611ed457600080fd5b8760208260061b8501011115611ee957600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611f6057600080fd5b83018035915067ffffffffffffffff821115611f7b57600080fd5b602001915036819003821315611b2457600080fd5b8035602083101561061d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203612024577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b60006020828403121561203d57600080fd5b5051919050565b8381526040602080830182905282820184905260009190859060608501845b878110156120b7578335612076816119dc565b73ffffffffffffffffffffffffffffffffffffffff16825267ffffffffffffffff6120a2858501611aaf565b16828401529284019290840190600101612063565b5098975050505050505050565b600080608083850312156120d757600080fd5b83601f8401126120e657600080fd5b6040516060810167ffffffffffffffff828210818311171561210a5761210a611d2d565b81604052829150606086018781111561212257600080fd5b865b8181101561213c578035845260209384019301612124565b509294509135918083111561215057600080fd5b505061215e85828601611d5c565b9150509250929050565b805160208083015191908110156121a7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b6000602082840312156121bf57600080fd5b8151801515811461196057600080fd5b600082516121e1818460208701611967565b919091019291505056fea164736f6c6343000813000a",
}

var MockFeeManagerABI = MockFeeManagerMetaData.ABI

var MockFeeManagerBin = MockFeeManagerMetaData.Bin

func DeployMockFeeManager(auth *bind.TransactOpts, backend bind.ContractBackend, _linkAddress common.Address, _nativeAddress common.Address, _proxyAddress common.Address, _rewardManagerAddress common.Address) (common.Address, *types.Transaction, *MockFeeManager, error) {
	parsed, err := MockFeeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockFeeManagerBin), backend, _linkAddress, _nativeAddress, _proxyAddress, _rewardManagerAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockFeeManager{address: address, abi: *parsed, MockFeeManagerCaller: MockFeeManagerCaller{contract: contract}, MockFeeManagerTransactor: MockFeeManagerTransactor{contract: contract}, MockFeeManagerFilterer: MockFeeManagerFilterer{contract: contract}}, nil
}

type MockFeeManager struct {
	address common.Address
	abi     abi.ABI
	MockFeeManagerCaller
	MockFeeManagerTransactor
	MockFeeManagerFilterer
}

type MockFeeManagerCaller struct {
	contract *bind.BoundContract
}

type MockFeeManagerTransactor struct {
	contract *bind.BoundContract
}

type MockFeeManagerFilterer struct {
	contract *bind.BoundContract
}

type MockFeeManagerSession struct {
	Contract     *MockFeeManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockFeeManagerCallerSession struct {
	Contract *MockFeeManagerCaller
	CallOpts bind.CallOpts
}

type MockFeeManagerTransactorSession struct {
	Contract     *MockFeeManagerTransactor
	TransactOpts bind.TransactOpts
}

type MockFeeManagerRaw struct {
	Contract *MockFeeManager
}

type MockFeeManagerCallerRaw struct {
	Contract *MockFeeManagerCaller
}

type MockFeeManagerTransactorRaw struct {
	Contract *MockFeeManagerTransactor
}

func NewMockFeeManager(address common.Address, backend bind.ContractBackend) (*MockFeeManager, error) {
	abi, err := abi.JSON(strings.NewReader(MockFeeManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockFeeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockFeeManager{address: address, abi: abi, MockFeeManagerCaller: MockFeeManagerCaller{contract: contract}, MockFeeManagerTransactor: MockFeeManagerTransactor{contract: contract}, MockFeeManagerFilterer: MockFeeManagerFilterer{contract: contract}}, nil
}

func NewMockFeeManagerCaller(address common.Address, caller bind.ContractCaller) (*MockFeeManagerCaller, error) {
	contract, err := bindMockFeeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerCaller{contract: contract}, nil
}

func NewMockFeeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*MockFeeManagerTransactor, error) {
	contract, err := bindMockFeeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerTransactor{contract: contract}, nil
}

func NewMockFeeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*MockFeeManagerFilterer, error) {
	contract, err := bindMockFeeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerFilterer{contract: contract}, nil
}

func bindMockFeeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockFeeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockFeeManager *MockFeeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockFeeManager.Contract.MockFeeManagerCaller.contract.Call(opts, result, method, params...)
}

func (_MockFeeManager *MockFeeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockFeeManager.Contract.MockFeeManagerTransactor.contract.Transfer(opts)
}

func (_MockFeeManager *MockFeeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockFeeManager.Contract.MockFeeManagerTransactor.contract.Transact(opts, method, params...)
}

func (_MockFeeManager *MockFeeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockFeeManager.Contract.contract.Call(opts, result, method, params...)
}

func (_MockFeeManager *MockFeeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockFeeManager.Contract.contract.Transfer(opts)
}

func (_MockFeeManager *MockFeeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockFeeManager.Contract.contract.Transact(opts, method, params...)
}

func (_MockFeeManager *MockFeeManagerCaller) GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quoteAddress common.Address) (CommonAsset, CommonAsset, *big.Int, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "getFeeAndReward", subscriber, report, quoteAddress)

	if err != nil {
		return *new(CommonAsset), *new(CommonAsset), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(CommonAsset)).(*CommonAsset)
	out1 := *abi.ConvertType(out[1], new(CommonAsset)).(*CommonAsset)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

func (_MockFeeManager *MockFeeManagerSession) GetFeeAndReward(subscriber common.Address, report []byte, quoteAddress common.Address) (CommonAsset, CommonAsset, *big.Int, error) {
	return _MockFeeManager.Contract.GetFeeAndReward(&_MockFeeManager.CallOpts, subscriber, report, quoteAddress)
}

func (_MockFeeManager *MockFeeManagerCallerSession) GetFeeAndReward(subscriber common.Address, report []byte, quoteAddress common.Address) (CommonAsset, CommonAsset, *big.Int, error) {
	return _MockFeeManager.Contract.GetFeeAndReward(&_MockFeeManager.CallOpts, subscriber, report, quoteAddress)
}

func (_MockFeeManager *MockFeeManagerCaller) ILinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "i_linkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) ILinkAddress() (common.Address, error) {
	return _MockFeeManager.Contract.ILinkAddress(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) ILinkAddress() (common.Address, error) {
	return _MockFeeManager.Contract.ILinkAddress(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) INativeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "i_nativeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) INativeAddress() (common.Address, error) {
	return _MockFeeManager.Contract.INativeAddress(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) INativeAddress() (common.Address, error) {
	return _MockFeeManager.Contract.INativeAddress(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) IProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "i_proxyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) IProxyAddress() (common.Address, error) {
	return _MockFeeManager.Contract.IProxyAddress(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) IProxyAddress() (common.Address, error) {
	return _MockFeeManager.Contract.IProxyAddress(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) IRewardManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "i_rewardManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) IRewardManager() (common.Address, error) {
	return _MockFeeManager.Contract.IRewardManager(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) IRewardManager() (common.Address, error) {
	return _MockFeeManager.Contract.IRewardManager(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _MockFeeManager.Contract.LinkAvailableForPayment(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _MockFeeManager.Contract.LinkAvailableForPayment(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) Owner() (common.Address, error) {
	return _MockFeeManager.Contract.Owner(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) Owner() (common.Address, error) {
	return _MockFeeManager.Contract.Owner(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) SGlobalDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "s_globalDiscounts", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) SGlobalDiscounts(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _MockFeeManager.Contract.SGlobalDiscounts(&_MockFeeManager.CallOpts, arg0, arg1)
}

func (_MockFeeManager *MockFeeManagerCallerSession) SGlobalDiscounts(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _MockFeeManager.Contract.SGlobalDiscounts(&_MockFeeManager.CallOpts, arg0, arg1)
}

func (_MockFeeManager *MockFeeManagerCaller) SLinkDeficit(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "s_linkDeficit", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) SLinkDeficit(arg0 [32]byte) (*big.Int, error) {
	return _MockFeeManager.Contract.SLinkDeficit(&_MockFeeManager.CallOpts, arg0)
}

func (_MockFeeManager *MockFeeManagerCallerSession) SLinkDeficit(arg0 [32]byte) (*big.Int, error) {
	return _MockFeeManager.Contract.SLinkDeficit(&_MockFeeManager.CallOpts, arg0)
}

func (_MockFeeManager *MockFeeManagerCaller) SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "s_nativeSurcharge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) SNativeSurcharge() (*big.Int, error) {
	return _MockFeeManager.Contract.SNativeSurcharge(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) SNativeSurcharge() (*big.Int, error) {
	return _MockFeeManager.Contract.SNativeSurcharge(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCaller) SSubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "s_subscriberDiscounts", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) SSubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _MockFeeManager.Contract.SSubscriberDiscounts(&_MockFeeManager.CallOpts, arg0, arg1, arg2)
}

func (_MockFeeManager *MockFeeManagerCallerSession) SSubscriberDiscounts(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _MockFeeManager.Contract.SSubscriberDiscounts(&_MockFeeManager.CallOpts, arg0, arg1, arg2)
}

func (_MockFeeManager *MockFeeManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MockFeeManager.Contract.SupportsInterface(&_MockFeeManager.CallOpts, interfaceId)
}

func (_MockFeeManager *MockFeeManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MockFeeManager.Contract.SupportsInterface(&_MockFeeManager.CallOpts, interfaceId)
}

func (_MockFeeManager *MockFeeManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockFeeManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockFeeManager *MockFeeManagerSession) TypeAndVersion() (string, error) {
	return _MockFeeManager.Contract.TypeAndVersion(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerCallerSession) TypeAndVersion() (string, error) {
	return _MockFeeManager.Contract.TypeAndVersion(&_MockFeeManager.CallOpts)
}

func (_MockFeeManager *MockFeeManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "acceptOwnership")
}

func (_MockFeeManager *MockFeeManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockFeeManager.Contract.AcceptOwnership(&_MockFeeManager.TransactOpts)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockFeeManager.Contract.AcceptOwnership(&_MockFeeManager.TransactOpts)
}

func (_MockFeeManager *MockFeeManagerTransactor) PayLinkDeficit(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "payLinkDeficit", configDigest)
}

func (_MockFeeManager *MockFeeManagerSession) PayLinkDeficit(configDigest [32]byte) (*types.Transaction, error) {
	return _MockFeeManager.Contract.PayLinkDeficit(&_MockFeeManager.TransactOpts, configDigest)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) PayLinkDeficit(configDigest [32]byte) (*types.Transaction, error) {
	return _MockFeeManager.Contract.PayLinkDeficit(&_MockFeeManager.TransactOpts, configDigest)
}

func (_MockFeeManager *MockFeeManagerTransactor) ProcessFee(opts *bind.TransactOpts, payload []byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "processFee", payload, parameterPayload, subscriber)
}

func (_MockFeeManager *MockFeeManagerSession) ProcessFee(payload []byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MockFeeManager.Contract.ProcessFee(&_MockFeeManager.TransactOpts, payload, parameterPayload, subscriber)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) ProcessFee(payload []byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MockFeeManager.Contract.ProcessFee(&_MockFeeManager.TransactOpts, payload, parameterPayload, subscriber)
}

func (_MockFeeManager *MockFeeManagerTransactor) ProcessFeeBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "processFeeBulk", payloads, parameterPayload, subscriber)
}

func (_MockFeeManager *MockFeeManagerSession) ProcessFeeBulk(payloads [][]byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MockFeeManager.Contract.ProcessFeeBulk(&_MockFeeManager.TransactOpts, payloads, parameterPayload, subscriber)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) ProcessFeeBulk(payloads [][]byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error) {
	return _MockFeeManager.Contract.ProcessFeeBulk(&_MockFeeManager.TransactOpts, payloads, parameterPayload, subscriber)
}

func (_MockFeeManager *MockFeeManagerTransactor) SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "setFeeRecipients", configDigest, rewardRecipientAndWeights)
}

func (_MockFeeManager *MockFeeManagerSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MockFeeManager.Contract.SetFeeRecipients(&_MockFeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) SetFeeRecipients(configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MockFeeManager.Contract.SetFeeRecipients(&_MockFeeManager.TransactOpts, configDigest, rewardRecipientAndWeights)
}

func (_MockFeeManager *MockFeeManagerTransactor) SetNativeSurcharge(opts *bind.TransactOpts, surcharge uint64) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "setNativeSurcharge", surcharge)
}

func (_MockFeeManager *MockFeeManagerSession) SetNativeSurcharge(surcharge uint64) (*types.Transaction, error) {
	return _MockFeeManager.Contract.SetNativeSurcharge(&_MockFeeManager.TransactOpts, surcharge)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) SetNativeSurcharge(surcharge uint64) (*types.Transaction, error) {
	return _MockFeeManager.Contract.SetNativeSurcharge(&_MockFeeManager.TransactOpts, surcharge)
}

func (_MockFeeManager *MockFeeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "transferOwnership", to)
}

func (_MockFeeManager *MockFeeManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockFeeManager.Contract.TransferOwnership(&_MockFeeManager.TransactOpts, to)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockFeeManager.Contract.TransferOwnership(&_MockFeeManager.TransactOpts, to)
}

func (_MockFeeManager *MockFeeManagerTransactor) UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "updateSubscriberDiscount", subscriber, feedId, token, discount)
}

func (_MockFeeManager *MockFeeManagerSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	return _MockFeeManager.Contract.UpdateSubscriberDiscount(&_MockFeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	return _MockFeeManager.Contract.UpdateSubscriberDiscount(&_MockFeeManager.TransactOpts, subscriber, feedId, token, discount)
}

func (_MockFeeManager *MockFeeManagerTransactor) UpdateSubscriberGlobalDiscount(opts *bind.TransactOpts, subscriber common.Address, token common.Address, discount uint64) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "updateSubscriberGlobalDiscount", subscriber, token, discount)
}

func (_MockFeeManager *MockFeeManagerSession) UpdateSubscriberGlobalDiscount(subscriber common.Address, token common.Address, discount uint64) (*types.Transaction, error) {
	return _MockFeeManager.Contract.UpdateSubscriberGlobalDiscount(&_MockFeeManager.TransactOpts, subscriber, token, discount)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) UpdateSubscriberGlobalDiscount(subscriber common.Address, token common.Address, discount uint64) (*types.Transaction, error) {
	return _MockFeeManager.Contract.UpdateSubscriberGlobalDiscount(&_MockFeeManager.TransactOpts, subscriber, token, discount)
}

func (_MockFeeManager *MockFeeManagerTransactor) Withdraw(opts *bind.TransactOpts, assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _MockFeeManager.contract.Transact(opts, "withdraw", assetAddress, recipient, quantity)
}

func (_MockFeeManager *MockFeeManagerSession) Withdraw(assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _MockFeeManager.Contract.Withdraw(&_MockFeeManager.TransactOpts, assetAddress, recipient, quantity)
}

func (_MockFeeManager *MockFeeManagerTransactorSession) Withdraw(assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error) {
	return _MockFeeManager.Contract.Withdraw(&_MockFeeManager.TransactOpts, assetAddress, recipient, quantity)
}

type MockFeeManagerDiscountAppliedIterator struct {
	Event *MockFeeManagerDiscountApplied

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerDiscountAppliedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerDiscountApplied)
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
		it.Event = new(MockFeeManagerDiscountApplied)
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

func (it *MockFeeManagerDiscountAppliedIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerDiscountAppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerDiscountApplied struct {
	ConfigDigest    [32]byte
	Subscriber      common.Address
	Fee             CommonAsset
	Reward          CommonAsset
	AppliedDiscount *big.Int
	Raw             types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterDiscountApplied(opts *bind.FilterOpts, configDigest [][32]byte, subscriber []common.Address) (*MockFeeManagerDiscountAppliedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}
	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "DiscountApplied", configDigestRule, subscriberRule)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerDiscountAppliedIterator{contract: _MockFeeManager.contract, event: "DiscountApplied", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchDiscountApplied(opts *bind.WatchOpts, sink chan<- *MockFeeManagerDiscountApplied, configDigest [][32]byte, subscriber []common.Address) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}
	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "DiscountApplied", configDigestRule, subscriberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerDiscountApplied)
				if err := _MockFeeManager.contract.UnpackLog(event, "DiscountApplied", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseDiscountApplied(log types.Log) (*MockFeeManagerDiscountApplied, error) {
	event := new(MockFeeManagerDiscountApplied)
	if err := _MockFeeManager.contract.UnpackLog(event, "DiscountApplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerInsufficientLinkIterator struct {
	Event *MockFeeManagerInsufficientLink

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerInsufficientLinkIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerInsufficientLink)
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
		it.Event = new(MockFeeManagerInsufficientLink)
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

func (it *MockFeeManagerInsufficientLinkIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerInsufficientLinkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerInsufficientLink struct {
	Rewards []IRewardManagerFeePayment
	Raw     types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterInsufficientLink(opts *bind.FilterOpts) (*MockFeeManagerInsufficientLinkIterator, error) {

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "InsufficientLink")
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerInsufficientLinkIterator{contract: _MockFeeManager.contract, event: "InsufficientLink", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *MockFeeManagerInsufficientLink) (event.Subscription, error) {

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "InsufficientLink")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerInsufficientLink)
				if err := _MockFeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseInsufficientLink(log types.Log) (*MockFeeManagerInsufficientLink, error) {
	event := new(MockFeeManagerInsufficientLink)
	if err := _MockFeeManager.contract.UnpackLog(event, "InsufficientLink", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerLinkDeficitClearedIterator struct {
	Event *MockFeeManagerLinkDeficitCleared

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerLinkDeficitClearedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerLinkDeficitCleared)
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
		it.Event = new(MockFeeManagerLinkDeficitCleared)
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

func (it *MockFeeManagerLinkDeficitClearedIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerLinkDeficitClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerLinkDeficitCleared struct {
	ConfigDigest [32]byte
	LinkQuantity *big.Int
	Raw          types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterLinkDeficitCleared(opts *bind.FilterOpts, configDigest [][32]byte) (*MockFeeManagerLinkDeficitClearedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "LinkDeficitCleared", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerLinkDeficitClearedIterator{contract: _MockFeeManager.contract, event: "LinkDeficitCleared", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchLinkDeficitCleared(opts *bind.WatchOpts, sink chan<- *MockFeeManagerLinkDeficitCleared, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "LinkDeficitCleared", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerLinkDeficitCleared)
				if err := _MockFeeManager.contract.UnpackLog(event, "LinkDeficitCleared", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseLinkDeficitCleared(log types.Log) (*MockFeeManagerLinkDeficitCleared, error) {
	event := new(MockFeeManagerLinkDeficitCleared)
	if err := _MockFeeManager.contract.UnpackLog(event, "LinkDeficitCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerNativeSurchargeUpdatedIterator struct {
	Event *MockFeeManagerNativeSurchargeUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerNativeSurchargeUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerNativeSurchargeUpdated)
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
		it.Event = new(MockFeeManagerNativeSurchargeUpdated)
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

func (it *MockFeeManagerNativeSurchargeUpdatedIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerNativeSurchargeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerNativeSurchargeUpdated struct {
	NewSurcharge uint64
	Raw          types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterNativeSurchargeUpdated(opts *bind.FilterOpts) (*MockFeeManagerNativeSurchargeUpdatedIterator, error) {

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "NativeSurchargeUpdated")
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerNativeSurchargeUpdatedIterator{contract: _MockFeeManager.contract, event: "NativeSurchargeUpdated", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchNativeSurchargeUpdated(opts *bind.WatchOpts, sink chan<- *MockFeeManagerNativeSurchargeUpdated) (event.Subscription, error) {

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "NativeSurchargeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerNativeSurchargeUpdated)
				if err := _MockFeeManager.contract.UnpackLog(event, "NativeSurchargeUpdated", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseNativeSurchargeUpdated(log types.Log) (*MockFeeManagerNativeSurchargeUpdated, error) {
	event := new(MockFeeManagerNativeSurchargeUpdated)
	if err := _MockFeeManager.contract.UnpackLog(event, "NativeSurchargeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerOwnershipTransferRequestedIterator struct {
	Event *MockFeeManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerOwnershipTransferRequested)
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
		it.Event = new(MockFeeManagerOwnershipTransferRequested)
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

func (it *MockFeeManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockFeeManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerOwnershipTransferRequestedIterator{contract: _MockFeeManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockFeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerOwnershipTransferRequested)
				if err := _MockFeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*MockFeeManagerOwnershipTransferRequested, error) {
	event := new(MockFeeManagerOwnershipTransferRequested)
	if err := _MockFeeManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerOwnershipTransferredIterator struct {
	Event *MockFeeManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerOwnershipTransferred)
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
		it.Event = new(MockFeeManagerOwnershipTransferred)
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

func (it *MockFeeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockFeeManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerOwnershipTransferredIterator{contract: _MockFeeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockFeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerOwnershipTransferred)
				if err := _MockFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*MockFeeManagerOwnershipTransferred, error) {
	event := new(MockFeeManagerOwnershipTransferred)
	if err := _MockFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerSubscriberDiscountUpdatedIterator struct {
	Event *MockFeeManagerSubscriberDiscountUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerSubscriberDiscountUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerSubscriberDiscountUpdated)
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
		it.Event = new(MockFeeManagerSubscriberDiscountUpdated)
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

func (it *MockFeeManagerSubscriberDiscountUpdatedIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerSubscriberDiscountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerSubscriberDiscountUpdated struct {
	Subscriber common.Address
	FeedId     [32]byte
	Token      common.Address
	Discount   uint64
	Raw        types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*MockFeeManagerSubscriberDiscountUpdatedIterator, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerSubscriberDiscountUpdatedIterator{contract: _MockFeeManager.contract, event: "SubscriberDiscountUpdated", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *MockFeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error) {

	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "SubscriberDiscountUpdated", subscriberRule, feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerSubscriberDiscountUpdated)
				if err := _MockFeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseSubscriberDiscountUpdated(log types.Log) (*MockFeeManagerSubscriberDiscountUpdated, error) {
	event := new(MockFeeManagerSubscriberDiscountUpdated)
	if err := _MockFeeManager.contract.UnpackLog(event, "SubscriberDiscountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockFeeManagerWithdrawIterator struct {
	Event *MockFeeManagerWithdraw

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockFeeManagerWithdrawIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockFeeManagerWithdraw)
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
		it.Event = new(MockFeeManagerWithdraw)
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

func (it *MockFeeManagerWithdrawIterator) Error() error {
	return it.fail
}

func (it *MockFeeManagerWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockFeeManagerWithdraw struct {
	AdminAddress common.Address
	Recipient    common.Address
	AssetAddress common.Address
	Quantity     *big.Int
	Raw          types.Log
}

func (_MockFeeManager *MockFeeManagerFilterer) FilterWithdraw(opts *bind.FilterOpts) (*MockFeeManagerWithdrawIterator, error) {

	logs, sub, err := _MockFeeManager.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &MockFeeManagerWithdrawIterator{contract: _MockFeeManager.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

func (_MockFeeManager *MockFeeManagerFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *MockFeeManagerWithdraw) (event.Subscription, error) {

	logs, sub, err := _MockFeeManager.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockFeeManagerWithdraw)
				if err := _MockFeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

func (_MockFeeManager *MockFeeManagerFilterer) ParseWithdraw(log types.Log) (*MockFeeManagerWithdraw, error) {
	event := new(MockFeeManagerWithdraw)
	if err := _MockFeeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MockFeeManager *MockFeeManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockFeeManager.abi.Events["DiscountApplied"].ID:
		return _MockFeeManager.ParseDiscountApplied(log)
	case _MockFeeManager.abi.Events["InsufficientLink"].ID:
		return _MockFeeManager.ParseInsufficientLink(log)
	case _MockFeeManager.abi.Events["LinkDeficitCleared"].ID:
		return _MockFeeManager.ParseLinkDeficitCleared(log)
	case _MockFeeManager.abi.Events["NativeSurchargeUpdated"].ID:
		return _MockFeeManager.ParseNativeSurchargeUpdated(log)
	case _MockFeeManager.abi.Events["OwnershipTransferRequested"].ID:
		return _MockFeeManager.ParseOwnershipTransferRequested(log)
	case _MockFeeManager.abi.Events["OwnershipTransferred"].ID:
		return _MockFeeManager.ParseOwnershipTransferred(log)
	case _MockFeeManager.abi.Events["SubscriberDiscountUpdated"].ID:
		return _MockFeeManager.ParseSubscriberDiscountUpdated(log)
	case _MockFeeManager.abi.Events["Withdraw"].ID:
		return _MockFeeManager.ParseWithdraw(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockFeeManagerDiscountApplied) Topic() common.Hash {
	return common.HexToHash("0x88b15eb682210089cddf967648e2cb2a4535aeadc8f8f36050922e33c04e7125")
}

func (MockFeeManagerInsufficientLink) Topic() common.Hash {
	return common.HexToHash("0xf52e5907b69d97c33392936c12d78b494463b78c5b72df50b4c497eee5720b67")
}

func (MockFeeManagerLinkDeficitCleared) Topic() common.Hash {
	return common.HexToHash("0x843f0b103e50b42b08f9d30f12f961845a6d02623730872e24644899c0dd9895")
}

func (MockFeeManagerNativeSurchargeUpdated) Topic() common.Hash {
	return common.HexToHash("0x08f7c0d17932ddb8523bc06754d42ff19ebc77d76a8b9bfde02c28ab1ed3d639")
}

func (MockFeeManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MockFeeManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MockFeeManagerSubscriberDiscountUpdated) Topic() common.Hash {
	return common.HexToHash("0x5eba5a8afa39780f0f99b6cbeb95f3da6a7040ca00abd46bdc91a0a060134139")
}

func (MockFeeManagerWithdraw) Topic() common.Hash {
	return common.HexToHash("0x7ff78a71698bdb18dcca96f52ab25e0a1b146fb6a49adf8e6845299e49021f29")
}

func (_MockFeeManager *MockFeeManager) Address() common.Address {
	return _MockFeeManager.address
}

type MockFeeManagerInterface interface {
	GetFeeAndReward(opts *bind.CallOpts, subscriber common.Address, report []byte, quoteAddress common.Address) (CommonAsset, CommonAsset, *big.Int, error)

	ILinkAddress(opts *bind.CallOpts) (common.Address, error)

	INativeAddress(opts *bind.CallOpts) (common.Address, error)

	IProxyAddress(opts *bind.CallOpts) (common.Address, error)

	IRewardManager(opts *bind.CallOpts) (common.Address, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SGlobalDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error)

	SLinkDeficit(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)

	SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error)

	SSubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	PayLinkDeficit(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	ProcessFee(opts *bind.TransactOpts, payload []byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error)

	ProcessFeeBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte, subscriber common.Address) (*types.Transaction, error)

	SetFeeRecipients(opts *bind.TransactOpts, configDigest [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetNativeSurcharge(opts *bind.TransactOpts, surcharge uint64) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscriberDiscount(opts *bind.TransactOpts, subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error)

	UpdateSubscriberGlobalDiscount(opts *bind.TransactOpts, subscriber common.Address, token common.Address, discount uint64) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, assetAddress common.Address, recipient common.Address, quantity *big.Int) (*types.Transaction, error)

	FilterDiscountApplied(opts *bind.FilterOpts, configDigest [][32]byte, subscriber []common.Address) (*MockFeeManagerDiscountAppliedIterator, error)

	WatchDiscountApplied(opts *bind.WatchOpts, sink chan<- *MockFeeManagerDiscountApplied, configDigest [][32]byte, subscriber []common.Address) (event.Subscription, error)

	ParseDiscountApplied(log types.Log) (*MockFeeManagerDiscountApplied, error)

	FilterInsufficientLink(opts *bind.FilterOpts) (*MockFeeManagerInsufficientLinkIterator, error)

	WatchInsufficientLink(opts *bind.WatchOpts, sink chan<- *MockFeeManagerInsufficientLink) (event.Subscription, error)

	ParseInsufficientLink(log types.Log) (*MockFeeManagerInsufficientLink, error)

	FilterLinkDeficitCleared(opts *bind.FilterOpts, configDigest [][32]byte) (*MockFeeManagerLinkDeficitClearedIterator, error)

	WatchLinkDeficitCleared(opts *bind.WatchOpts, sink chan<- *MockFeeManagerLinkDeficitCleared, configDigest [][32]byte) (event.Subscription, error)

	ParseLinkDeficitCleared(log types.Log) (*MockFeeManagerLinkDeficitCleared, error)

	FilterNativeSurchargeUpdated(opts *bind.FilterOpts) (*MockFeeManagerNativeSurchargeUpdatedIterator, error)

	WatchNativeSurchargeUpdated(opts *bind.WatchOpts, sink chan<- *MockFeeManagerNativeSurchargeUpdated) (event.Subscription, error)

	ParseNativeSurchargeUpdated(log types.Log) (*MockFeeManagerNativeSurchargeUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockFeeManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockFeeManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MockFeeManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockFeeManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockFeeManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MockFeeManagerOwnershipTransferred, error)

	FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*MockFeeManagerSubscriberDiscountUpdatedIterator, error)

	WatchSubscriberDiscountUpdated(opts *bind.WatchOpts, sink chan<- *MockFeeManagerSubscriberDiscountUpdated, subscriber []common.Address, feedId [][32]byte) (event.Subscription, error)

	ParseSubscriberDiscountUpdated(log types.Log) (*MockFeeManagerSubscriberDiscountUpdated, error)

	FilterWithdraw(opts *bind.FilterOpts) (*MockFeeManagerWithdrawIterator, error)

	WatchWithdraw(opts *bind.WatchOpts, sink chan<- *MockFeeManagerWithdraw) (event.Subscription, error)

	ParseWithdraw(log types.Log) (*MockFeeManagerWithdraw, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
