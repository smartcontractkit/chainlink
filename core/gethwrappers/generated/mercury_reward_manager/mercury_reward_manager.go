// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_reward_manager

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

var MercuryRewardManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWeights\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"FeeManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"RewardRecipientsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"RewardsClaimed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"poolIds\",\"type\":\"bytes32[]\"}],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getAvailableRewardPoolIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.Asset\",\"name\":\"fee\",\"type\":\"tuple\"}],\"name\":\"onFeePaid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"}],\"name\":\"payRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"registeredPoolIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"rewardRecipientWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"totalRewardRecipientFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"updateRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001944380380620019448339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b6080516117426200020260003960008181610877015281816109600152610daa01526117426000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806379ba509711610097578063a241424b11610066578063a241424b14610275578063e9cd6ff3146102a0578063f2fde38b146102c0578063f34517aa146102d357600080fd5b806379ba50971461021157806384afb76e1461021957806389762bd21461022c5780638da5cb5b1461024d57600080fd5b8063472d35b9116100d3578063472d35b9146101b85780634d322084146101cb578063633b5f6e146101de5780636992922f146101f157600080fd5b806301ffc9a7146100fa5780630f3c34d114610164578063181f5a7714610179575b600080fd5b61014f6101083660046111b8565b7fffffffff00000000000000000000000000000000000000000000000000000000167f84afb76e000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b610177610172366004611278565b6102e6565b005b604080518082018252601381527f5265776172644d616e6167657220302e302e31000000000000000000000000006020820152905161015b919061131e565b6101776101c63660046113b3565b6102f4565b6101776101d93660046113d5565b610375565b6101776101ec366004611454565b6104ac565b6102046101ff3660046113b3565b610662565b60405161015b91906114c0565b61017761075e565b610177610227366004611504565b610860565b61023f61023a36600461156a565b6109d5565b60405190815260200161015b565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161015b565b61023f610283366004611583565b600460209081526000928352604080842090915290825290205481565b61023f6102ae36600461156a565b60026020526000908152604090205481565b6101776102ce3660046113b3565b6109f6565b6101776102e1366004611454565b610a0a565b6102f03382610c16565b5050565b6102fc610e22565b600780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b49060200160405180910390a150565b8261039560005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156103e757506000818152600460209081526040808320338452909152902054155b1561041e576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516001808252818301909252600091602080830190803683370190505090508481600081518110610454576104546115af565b60200260200101818152505060005b838110156104a45761049b858583818110610480576104806115af565b905060200201602081019061049591906113b3565b83610c16565b50600101610463565b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906104ec575060075473ffffffffffffffffffffffffffffffffffffffff163314155b15610523576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081900361055e576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526005602052604090205460ff16156105a7576040517f0afa7ee800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6006805460018181019092557ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01849055600084815260056020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055610623838383670de0b6b3a7640000610ea5565b827fe5ca4131deaeb848b5d6a9c0f85795efc54e8a0252eb38d3e77a1efc2b38419683836040516106559291906115de565b60405180910390a2505050565b60065460609060008167ffffffffffffffff811115610683576106836111fa565b6040519080825280602002602001820160405280156106ac578160200160208202803683370190505b5090506000805b83811015610754576000600682815481106106d0576106d06115af565b600091825260208083209091015480835260048252604080842073ffffffffffffffffffffffffffffffffffffffff8c1685529092529120549091501561074b576000818152600260205260409020541561074b5780848481518110610738576107386115af565b6020026020010181815250508260010192505b506001016106b3565b5090949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146107e4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166108a460208301836113b3565b73ffffffffffffffffffffffffffffffffffffffff16146108f1576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600083815260026020908152604091829020805491840135918201905590517f23b872dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015230602483015260448201929092527f0000000000000000000000000000000000000000000000000000000000000000909116906323b872dd906064016020604051808303816000875af11580156109ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109cf9190611642565b50505050565b600681815481106109e557600080fd5b600091825260209091200154905081565b6109fe610e22565b610a078161100c565b50565b610a12610e22565b610a6d8282808060200260200160405190810160405280939291908181526020016000905b82821015610a6357610a5460408302860136819003810190611664565b81526020019060010190610a37565b5050505050611101565b15610aa4576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516001808252818301909252600091602080830190803683370190505090508381600081518110610ada57610ada6115af565b6020026020010181815250506000805b83811015610bc8576000858583818110610b0657610b066115af565b610b1c92602060409092020190810191506113b3565b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff85168452909152812054919250819003610b88576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610bb9878785818110610b9d57610b9d6115af565b610bb392602060409092020190810191506113b3565b86610c16565b50929092019150600101610aea565b50610bd585858584610ea5565b847fe5ca4131deaeb848b5d6a9c0f85795efc54e8a0252eb38d3e77a1efc2b3841968585604051610c079291906115de565b60405180910390a25050505050565b60008060005b8351811015610d57576000848281518110610c3957610c396115af565b6020908102919091018101516000818152600283526040808220546003855281832073ffffffffffffffffffffffffffffffffffffffff8c16845290945281205491935090820390819003610c9057505050610d47565b600083815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8c16808552908352818420548785526003845282852082865290935292208490558851670de0b6b3a764000091840291909104968701969190899087908110610d0157610d016115af565b60200260200101517ffec539b3c42d74cd492c3ecf7966bf014b327beef98e935cdc3ec5e54a6901c788604051610d3a91815260200190565b60405180910390a3505050505b610d50816116ea565b9050610c1c565b508015610e19576040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152602482018390527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015610df3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e179190611642565b505b90505b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ea3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016107db565b565b6000805b83811015610fcb576000858583818110610ec557610ec56115af565b9050604002016020013590506000868684818110610ee557610ee56115af565b610efb92602060409092020190810191506113b3565b905073ffffffffffffffffffffffffffffffffffffffff8116610f4a576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600003610f84576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff9094168352929052208190559190910190610fc4816116ea565b9050610ea9565b50818114611005576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361108b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016107db565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000805b82518110156111af57600061111b826001611722565b90505b83518110156111a657838181518110611139576111396115af565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1684838151811061116d5761116d6115af565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff160361119e575060019392505050565b60010161111e565b50600101611105565b50600092915050565b6000602082840312156111ca57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610e1957600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611270576112706111fa565b604052919050565b6000602080838503121561128b57600080fd5b823567ffffffffffffffff808211156112a357600080fd5b818501915085601f8301126112b757600080fd5b8135818111156112c9576112c96111fa565b8060051b91506112da848301611229565b81815291830184019184810190888411156112f457600080fd5b938501935b83851015611312578435825293850193908501906112f9565b98975050505050505050565b600060208083528351808285015260005b8181101561134b5785810183015185820160400152820161132f565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146113ae57600080fd5b919050565b6000602082840312156113c557600080fd5b6113ce8261138a565b9392505050565b6000806000604084860312156113ea57600080fd5b83359250602084013567ffffffffffffffff8082111561140957600080fd5b818601915086601f83011261141d57600080fd5b81358181111561142c57600080fd5b8760208260051b850101111561144157600080fd5b6020830194508093505050509250925092565b60008060006040848603121561146957600080fd5b83359250602084013567ffffffffffffffff8082111561148857600080fd5b818601915086601f83011261149c57600080fd5b8135818111156114ab57600080fd5b8760208260061b850101111561144157600080fd5b6020808252825182820181905260009190848201906040850190845b818110156114f8578351835292840192918401916001016114dc565b50909695505050505050565b6000806000838503608081121561151a57600080fd5b8435935061152a6020860161138a565b925060407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08201121561155c57600080fd5b506040840190509250925092565b60006020828403121561157c57600080fd5b5035919050565b6000806040838503121561159657600080fd5b823591506115a66020840161138a565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020808252818101839052600090604080840186845b878110156116355773ffffffffffffffffffffffffffffffffffffffff61161a8361138a565b168352818501358584015291830191908301906001016115f4565b5090979650505050505050565b60006020828403121561165457600080fd5b81518015158114610e1957600080fd5b60006040828403121561167657600080fd5b6040516040810181811067ffffffffffffffff82111715611699576116996111fa565b6040526116a58361138a565b8152602083013560208201528091505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361171b5761171b6116bb565b5060010190565b80820180821115610e1c57610e1c6116bb56fea164736f6c6343000810000a",
}

var MercuryRewardManagerABI = MercuryRewardManagerMetaData.ABI

var MercuryRewardManagerBin = MercuryRewardManagerMetaData.Bin

func DeployMercuryRewardManager(auth *bind.TransactOpts, backend bind.ContractBackend, linkAddress common.Address) (common.Address, *types.Transaction, *MercuryRewardManager, error) {
	parsed, err := MercuryRewardManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryRewardManagerBin), backend, linkAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryRewardManager{MercuryRewardManagerCaller: MercuryRewardManagerCaller{contract: contract}, MercuryRewardManagerTransactor: MercuryRewardManagerTransactor{contract: contract}, MercuryRewardManagerFilterer: MercuryRewardManagerFilterer{contract: contract}}, nil
}

type MercuryRewardManager struct {
	address common.Address
	abi     abi.ABI
	MercuryRewardManagerCaller
	MercuryRewardManagerTransactor
	MercuryRewardManagerFilterer
}

type MercuryRewardManagerCaller struct {
	contract *bind.BoundContract
}

type MercuryRewardManagerTransactor struct {
	contract *bind.BoundContract
}

type MercuryRewardManagerFilterer struct {
	contract *bind.BoundContract
}

type MercuryRewardManagerSession struct {
	Contract     *MercuryRewardManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryRewardManagerCallerSession struct {
	Contract *MercuryRewardManagerCaller
	CallOpts bind.CallOpts
}

type MercuryRewardManagerTransactorSession struct {
	Contract     *MercuryRewardManagerTransactor
	TransactOpts bind.TransactOpts
}

type MercuryRewardManagerRaw struct {
	Contract *MercuryRewardManager
}

type MercuryRewardManagerCallerRaw struct {
	Contract *MercuryRewardManagerCaller
}

type MercuryRewardManagerTransactorRaw struct {
	Contract *MercuryRewardManagerTransactor
}

func NewMercuryRewardManager(address common.Address, backend bind.ContractBackend) (*MercuryRewardManager, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryRewardManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryRewardManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManager{address: address, abi: abi, MercuryRewardManagerCaller: MercuryRewardManagerCaller{contract: contract}, MercuryRewardManagerTransactor: MercuryRewardManagerTransactor{contract: contract}, MercuryRewardManagerFilterer: MercuryRewardManagerFilterer{contract: contract}}, nil
}

func NewMercuryRewardManagerCaller(address common.Address, caller bind.ContractCaller) (*MercuryRewardManagerCaller, error) {
	contract, err := bindMercuryRewardManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerCaller{contract: contract}, nil
}

func NewMercuryRewardManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryRewardManagerTransactor, error) {
	contract, err := bindMercuryRewardManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerTransactor{contract: contract}, nil
}

func NewMercuryRewardManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryRewardManagerFilterer, error) {
	contract, err := bindMercuryRewardManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerFilterer{contract: contract}, nil
}

func bindMercuryRewardManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryRewardManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryRewardManager *MercuryRewardManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryRewardManager.Contract.MercuryRewardManagerCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryRewardManager *MercuryRewardManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.MercuryRewardManagerTransactor.contract.Transfer(opts)
}

func (_MercuryRewardManager *MercuryRewardManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.MercuryRewardManagerTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryRewardManager.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.contract.Transfer(opts)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "getAvailableRewardPoolIds", recipient)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) GetAvailableRewardPoolIds(recipient common.Address) ([][32]byte, error) {
	return _MercuryRewardManager.Contract.GetAvailableRewardPoolIds(&_MercuryRewardManager.CallOpts, recipient)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) GetAvailableRewardPoolIds(recipient common.Address) ([][32]byte, error) {
	return _MercuryRewardManager.Contract.GetAvailableRewardPoolIds(&_MercuryRewardManager.CallOpts, recipient)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) Owner() (common.Address, error) {
	return _MercuryRewardManager.Contract.Owner(&_MercuryRewardManager.CallOpts)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) Owner() (common.Address, error) {
	return _MercuryRewardManager.Contract.Owner(&_MercuryRewardManager.CallOpts)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) RegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "registeredPoolIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) RegisteredPoolIds(arg0 *big.Int) ([32]byte, error) {
	return _MercuryRewardManager.Contract.RegisteredPoolIds(&_MercuryRewardManager.CallOpts, arg0)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) RegisteredPoolIds(arg0 *big.Int) ([32]byte, error) {
	return _MercuryRewardManager.Contract.RegisteredPoolIds(&_MercuryRewardManager.CallOpts, arg0)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) RewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "rewardRecipientWeights", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) RewardRecipientWeights(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _MercuryRewardManager.Contract.RewardRecipientWeights(&_MercuryRewardManager.CallOpts, arg0, arg1)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) RewardRecipientWeights(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _MercuryRewardManager.Contract.RewardRecipientWeights(&_MercuryRewardManager.CallOpts, arg0, arg1)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MercuryRewardManager.Contract.SupportsInterface(&_MercuryRewardManager.CallOpts, interfaceId)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MercuryRewardManager.Contract.SupportsInterface(&_MercuryRewardManager.CallOpts, interfaceId)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) TotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "totalRewardRecipientFees", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) TotalRewardRecipientFees(arg0 [32]byte) (*big.Int, error) {
	return _MercuryRewardManager.Contract.TotalRewardRecipientFees(&_MercuryRewardManager.CallOpts, arg0)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) TotalRewardRecipientFees(arg0 [32]byte) (*big.Int, error) {
	return _MercuryRewardManager.Contract.TotalRewardRecipientFees(&_MercuryRewardManager.CallOpts, arg0)
}

func (_MercuryRewardManager *MercuryRewardManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryRewardManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryRewardManager *MercuryRewardManagerSession) TypeAndVersion() (string, error) {
	return _MercuryRewardManager.Contract.TypeAndVersion(&_MercuryRewardManager.CallOpts)
}

func (_MercuryRewardManager *MercuryRewardManagerCallerSession) TypeAndVersion() (string, error) {
	return _MercuryRewardManager.Contract.TypeAndVersion(&_MercuryRewardManager.CallOpts)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "acceptOwnership")
}

func (_MercuryRewardManager *MercuryRewardManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.AcceptOwnership(&_MercuryRewardManager.TransactOpts)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.AcceptOwnership(&_MercuryRewardManager.TransactOpts)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) ClaimRewards(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "claimRewards", poolIds)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) ClaimRewards(poolIds [][32]byte) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.ClaimRewards(&_MercuryRewardManager.TransactOpts, poolIds)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) ClaimRewards(poolIds [][32]byte) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.ClaimRewards(&_MercuryRewardManager.TransactOpts, poolIds)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) OnFeePaid(opts *bind.TransactOpts, poolId [32]byte, payee common.Address, fee CommonAsset) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "onFeePaid", poolId, payee, fee)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) OnFeePaid(poolId [32]byte, payee common.Address, fee CommonAsset) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.OnFeePaid(&_MercuryRewardManager.TransactOpts, poolId, payee, fee)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) OnFeePaid(poolId [32]byte, payee common.Address, fee CommonAsset) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.OnFeePaid(&_MercuryRewardManager.TransactOpts, poolId, payee, fee)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) PayRecipients(opts *bind.TransactOpts, poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "payRecipients", poolId, recipients)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) PayRecipients(poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.PayRecipients(&_MercuryRewardManager.TransactOpts, poolId, recipients)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) PayRecipients(poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.PayRecipients(&_MercuryRewardManager.TransactOpts, poolId, recipients)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) SetFeeManager(opts *bind.TransactOpts, newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "setFeeManager", newFeeManagerAddress)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) SetFeeManager(newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.SetFeeManager(&_MercuryRewardManager.TransactOpts, newFeeManagerAddress)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) SetFeeManager(newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.SetFeeManager(&_MercuryRewardManager.TransactOpts, newFeeManagerAddress)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) SetRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "setRewardRecipients", poolId, rewardRecipientAndWeights)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) SetRewardRecipients(poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.SetRewardRecipients(&_MercuryRewardManager.TransactOpts, poolId, rewardRecipientAndWeights)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) SetRewardRecipients(poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.SetRewardRecipients(&_MercuryRewardManager.TransactOpts, poolId, rewardRecipientAndWeights)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "transferOwnership", to)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.TransferOwnership(&_MercuryRewardManager.TransactOpts, to)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.TransferOwnership(&_MercuryRewardManager.TransactOpts, to)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactor) UpdateRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryRewardManager.contract.Transact(opts, "updateRewardRecipients", poolId, newRewardRecipients)
}

func (_MercuryRewardManager *MercuryRewardManagerSession) UpdateRewardRecipients(poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.UpdateRewardRecipients(&_MercuryRewardManager.TransactOpts, poolId, newRewardRecipients)
}

func (_MercuryRewardManager *MercuryRewardManagerTransactorSession) UpdateRewardRecipients(poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryRewardManager.Contract.UpdateRewardRecipients(&_MercuryRewardManager.TransactOpts, poolId, newRewardRecipients)
}

type MercuryRewardManagerFeeManagerUpdatedIterator struct {
	Event *MercuryRewardManagerFeeManagerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryRewardManagerFeeManagerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryRewardManagerFeeManagerUpdated)
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
		it.Event = new(MercuryRewardManagerFeeManagerUpdated)
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

func (it *MercuryRewardManagerFeeManagerUpdatedIterator) Error() error {
	return it.fail
}

func (it *MercuryRewardManagerFeeManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryRewardManagerFeeManagerUpdated struct {
	NewFeeManagerAddress common.Address
	Raw                  types.Log
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) FilterFeeManagerUpdated(opts *bind.FilterOpts) (*MercuryRewardManagerFeeManagerUpdatedIterator, error) {

	logs, sub, err := _MercuryRewardManager.contract.FilterLogs(opts, "FeeManagerUpdated")
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerFeeManagerUpdatedIterator{contract: _MercuryRewardManager.contract, event: "FeeManagerUpdated", logs: logs, sub: sub}, nil
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) WatchFeeManagerUpdated(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerFeeManagerUpdated) (event.Subscription, error) {

	logs, sub, err := _MercuryRewardManager.contract.WatchLogs(opts, "FeeManagerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryRewardManagerFeeManagerUpdated)
				if err := _MercuryRewardManager.contract.UnpackLog(event, "FeeManagerUpdated", log); err != nil {
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

func (_MercuryRewardManager *MercuryRewardManagerFilterer) ParseFeeManagerUpdated(log types.Log) (*MercuryRewardManagerFeeManagerUpdated, error) {
	event := new(MercuryRewardManagerFeeManagerUpdated)
	if err := _MercuryRewardManager.contract.UnpackLog(event, "FeeManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryRewardManagerOwnershipTransferRequestedIterator struct {
	Event *MercuryRewardManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryRewardManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryRewardManagerOwnershipTransferRequested)
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
		it.Event = new(MercuryRewardManagerOwnershipTransferRequested)
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

func (it *MercuryRewardManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MercuryRewardManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryRewardManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryRewardManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerOwnershipTransferRequestedIterator{contract: _MercuryRewardManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryRewardManagerOwnershipTransferRequested)
				if err := _MercuryRewardManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MercuryRewardManager *MercuryRewardManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*MercuryRewardManagerOwnershipTransferRequested, error) {
	event := new(MercuryRewardManagerOwnershipTransferRequested)
	if err := _MercuryRewardManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryRewardManagerOwnershipTransferredIterator struct {
	Event *MercuryRewardManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryRewardManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryRewardManagerOwnershipTransferred)
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
		it.Event = new(MercuryRewardManagerOwnershipTransferred)
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

func (it *MercuryRewardManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MercuryRewardManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryRewardManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryRewardManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerOwnershipTransferredIterator{contract: _MercuryRewardManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryRewardManagerOwnershipTransferred)
				if err := _MercuryRewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MercuryRewardManager *MercuryRewardManagerFilterer) ParseOwnershipTransferred(log types.Log) (*MercuryRewardManagerOwnershipTransferred, error) {
	event := new(MercuryRewardManagerOwnershipTransferred)
	if err := _MercuryRewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryRewardManagerRewardRecipientsUpdatedIterator struct {
	Event *MercuryRewardManagerRewardRecipientsUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryRewardManagerRewardRecipientsUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryRewardManagerRewardRecipientsUpdated)
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
		it.Event = new(MercuryRewardManagerRewardRecipientsUpdated)
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

func (it *MercuryRewardManagerRewardRecipientsUpdatedIterator) Error() error {
	return it.fail
}

func (it *MercuryRewardManagerRewardRecipientsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryRewardManagerRewardRecipientsUpdated struct {
	PoolId              [32]byte
	NewRewardRecipients []CommonAddressAndWeight
	Raw                 types.Log
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*MercuryRewardManagerRewardRecipientsUpdatedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.FilterLogs(opts, "RewardRecipientsUpdated", poolIdRule)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerRewardRecipientsUpdatedIterator{contract: _MercuryRewardManager.contract, event: "RewardRecipientsUpdated", logs: logs, sub: sub}, nil
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) WatchRewardRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerRewardRecipientsUpdated, poolId [][32]byte) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.WatchLogs(opts, "RewardRecipientsUpdated", poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryRewardManagerRewardRecipientsUpdated)
				if err := _MercuryRewardManager.contract.UnpackLog(event, "RewardRecipientsUpdated", log); err != nil {
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

func (_MercuryRewardManager *MercuryRewardManagerFilterer) ParseRewardRecipientsUpdated(log types.Log) (*MercuryRewardManagerRewardRecipientsUpdated, error) {
	event := new(MercuryRewardManagerRewardRecipientsUpdated)
	if err := _MercuryRewardManager.contract.UnpackLog(event, "RewardRecipientsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryRewardManagerRewardsClaimedIterator struct {
	Event *MercuryRewardManagerRewardsClaimed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryRewardManagerRewardsClaimedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryRewardManagerRewardsClaimed)
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
		it.Event = new(MercuryRewardManagerRewardsClaimed)
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

func (it *MercuryRewardManagerRewardsClaimedIterator) Error() error {
	return it.fail
}

func (it *MercuryRewardManagerRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryRewardManagerRewardsClaimed struct {
	PoolId    [32]byte
	Recipient common.Address
	Quantity  *big.Int
	Raw       types.Log
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*MercuryRewardManagerRewardsClaimedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.FilterLogs(opts, "RewardsClaimed", poolIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &MercuryRewardManagerRewardsClaimedIterator{contract: _MercuryRewardManager.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

func (_MercuryRewardManager *MercuryRewardManagerFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerRewardsClaimed, poolId [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MercuryRewardManager.contract.WatchLogs(opts, "RewardsClaimed", poolIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryRewardManagerRewardsClaimed)
				if err := _MercuryRewardManager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

func (_MercuryRewardManager *MercuryRewardManagerFilterer) ParseRewardsClaimed(log types.Log) (*MercuryRewardManagerRewardsClaimed, error) {
	event := new(MercuryRewardManagerRewardsClaimed)
	if err := _MercuryRewardManager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MercuryRewardManager *MercuryRewardManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MercuryRewardManager.abi.Events["FeeManagerUpdated"].ID:
		return _MercuryRewardManager.ParseFeeManagerUpdated(log)
	case _MercuryRewardManager.abi.Events["OwnershipTransferRequested"].ID:
		return _MercuryRewardManager.ParseOwnershipTransferRequested(log)
	case _MercuryRewardManager.abi.Events["OwnershipTransferred"].ID:
		return _MercuryRewardManager.ParseOwnershipTransferred(log)
	case _MercuryRewardManager.abi.Events["RewardRecipientsUpdated"].ID:
		return _MercuryRewardManager.ParseRewardRecipientsUpdated(log)
	case _MercuryRewardManager.abi.Events["RewardsClaimed"].ID:
		return _MercuryRewardManager.ParseRewardsClaimed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MercuryRewardManagerFeeManagerUpdated) Topic() common.Hash {
	return common.HexToHash("0xe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b4")
}

func (MercuryRewardManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MercuryRewardManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MercuryRewardManagerRewardRecipientsUpdated) Topic() common.Hash {
	return common.HexToHash("0xe5ca4131deaeb848b5d6a9c0f85795efc54e8a0252eb38d3e77a1efc2b384196")
}

func (MercuryRewardManagerRewardsClaimed) Topic() common.Hash {
	return common.HexToHash("0xfec539b3c42d74cd492c3ecf7966bf014b327beef98e935cdc3ec5e54a6901c7")
}

func (_MercuryRewardManager *MercuryRewardManager) Address() common.Address {
	return _MercuryRewardManager.address
}

type MercuryRewardManagerInterface interface {
	GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address) ([][32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	RegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	RewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ClaimRewards(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error)

	OnFeePaid(opts *bind.TransactOpts, poolId [32]byte, payee common.Address, fee CommonAsset) (*types.Transaction, error)

	PayRecipients(opts *bind.TransactOpts, poolId [32]byte, recipients []common.Address) (*types.Transaction, error)

	SetFeeManager(opts *bind.TransactOpts, newFeeManagerAddress common.Address) (*types.Transaction, error)

	SetRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error)

	FilterFeeManagerUpdated(opts *bind.FilterOpts) (*MercuryRewardManagerFeeManagerUpdatedIterator, error)

	WatchFeeManagerUpdated(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerFeeManagerUpdated) (event.Subscription, error)

	ParseFeeManagerUpdated(log types.Log) (*MercuryRewardManagerFeeManagerUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryRewardManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MercuryRewardManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryRewardManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MercuryRewardManagerOwnershipTransferred, error)

	FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*MercuryRewardManagerRewardRecipientsUpdatedIterator, error)

	WatchRewardRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerRewardRecipientsUpdated, poolId [][32]byte) (event.Subscription, error)

	ParseRewardRecipientsUpdated(log types.Log) (*MercuryRewardManagerRewardRecipientsUpdated, error)

	FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*MercuryRewardManagerRewardsClaimedIterator, error)

	WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *MercuryRewardManagerRewardsClaimed, poolId [][32]byte, recipient []common.Address) (event.Subscription, error)

	ParseRewardsClaimed(log types.Log) (*MercuryRewardManagerRewardsClaimed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
