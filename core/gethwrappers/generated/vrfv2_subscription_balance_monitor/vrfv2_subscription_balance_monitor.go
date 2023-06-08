// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2_subscription_balance_monitor

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

var VRFSubscriptionBalanceMonitorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"keeperRegistryAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minWaitPeriodSeconds\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"duplicate\",\"type\":\"uint64\"}],\"name\":\"DuplicateSubcriptionId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWatchList\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyKeeperRegistry\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountAdded\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountWithdrawn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"KeeperRegistryAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"LinkTokenAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMinWaitPeriod\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMinWaitPeriod\",\"type\":\"uint256\"}],\"name\":\"MinWaitPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"OutOfGas\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"TopUpFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"TopUpSucceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"VRFCoordinatorV2AddressUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINKTOKEN\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscriptionInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"},{\"internalType\":\"uint96\",\"name\":\"minBalanceJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"topUpAmountJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint56\",\"name\":\"lastTopUpTimestamp\",\"type\":\"uint56\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnderfundedSubscriptions\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWatchList\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keeperRegistryAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_minWaitPeriodSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_watchList\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeperRegistryAddress\",\"type\":\"address\"}],\"name\":\"setKeeperRegistryAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkTokenAddress\",\"type\":\"address\"}],\"name\":\"setLinkTokenAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setMinWaitPeriodSeconds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"setVRFCoordinatorV2Address\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"subscriptionIds\",\"type\":\"uint64[]\"},{\"internalType\":\"uint96[]\",\"name\":\"minBalancesJuels\",\"type\":\"uint96[]\"},{\"internalType\":\"uint96[]\",\"name\":\"topUpAmountsJuels\",\"type\":\"uint96[]\"}],\"name\":\"setWatchList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"needsFunding\",\"type\":\"uint64[]\"}],\"name\":\"topUp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162002a8f38038062002a8f83398101604081905262000034916200040b565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000104565b50506001805460ff60a01b1916905550620000d984620001b0565b620000e48362000237565b620000ef82620002be565b620000fa8162000345565b505050506200045d565b6001600160a01b0381163314156200015f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001ba62000390565b6001600160a01b038116620001ce57600080fd5b600354604080516001600160a01b03928316815291831660208301527fee7e95e098f422f231397e2532a2752a013a51f4122cee6c30b18f930cda91dc910160405180910390a1600380546001600160a01b0319166001600160a01b0392909216919091179055565b6200024162000390565b6001600160a01b0381166200025557600080fd5b600254604080516001600160a01b03928316815291831660208301527f4490e5bf8542a3633d8f268c3733706aa29acea979152e5fe5befd9cc7ad43a9910160405180910390a1600280546001600160a01b0319166001600160a01b0392909216919091179055565b620002c862000390565b6001600160a01b038116620002dc57600080fd5b600454604080516001600160a01b03928316815291831660208301527fb732223055abcde751d7a24272ffc8a3aa571cb72b443969a4199b7ecd59f8b9910160405180910390a1600480546001600160a01b0319166001600160a01b0392909216919091179055565b6200034f62000390565b60055460408051918252602082018390527f04330086c73b1fe1e13cd47a61c692e7c4399b5de08ed94b7ab824684af09323910160405180910390a1600555565b6000546001600160a01b03163314620003ec5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200040657600080fd5b919050565b600080600080608085870312156200042257600080fd5b6200042d85620003ee565b93506200043d60208601620003ee565b92506200044d60408601620003ee565b6060959095015193969295505050565b612622806200046d6000396000f3fe608060405234801561001057600080fd5b506004361061018c5760003560e01c8063728584b7116100e35780638da5cb5b1161008c578063c36805b411610066578063c36805b414610459578063d23d815b1461046c578063f2fde38b1461047f57600080fd5b80638da5cb5b1461041557806394555114146104335780639c4376c51461044657600080fd5b80637ac24d5a116100bd5780637ac24d5a146103ee5780638456cb59146103f657806385879755146103fe57600080fd5b8063728584b7146103b1578063771b081e146103c657806379ba5097146103e657600080fd5b80634585e33b116101455780635c975abb1161011f5780635c975abb1461027d5780636dd2702f146102ab5780636e04ff0d1461039057600080fd5b80634585e33b1461021e578063482081e51461023157806355380dfb1461025d57600080fd5b80633b2bcbf1116101765780633b2bcbf1146101b95780633f4ba83a146102035780633f85861f1461020b57600080fd5b8062f714ce146101915780630a61af93146101a6575b600080fd5b6101a461019f3660046121ec565b610492565b005b6101a46101b4366004611fe6565b6105b9565b6002546101d99073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6101a4610acc565b6101a46102193660046121ba565b610ade565b6101a461022c366004612148565b610b27565b61024461023f3660046121ba565b610c16565b60405167ffffffffffffffff90911681526020016101fa565b6003546101d99073ffffffffffffffffffffffffffffffffffffffff1681565b60015474010000000000000000000000000000000000000000900460ff1660405190151581526020016101fa565b6103506102b936600461221c565b67ffffffffffffffff166000908152600760209081526040918290208251608081018452905460ff8116151580835261010082046bffffffffffffffffffffffff9081169484018590526d010000000000000000000000000083041694830185905279010000000000000000000000000000000000000000000000000090910466ffffffffffffff16606090920182905293919291565b6040516101fa949392919093151584526bffffffffffffffffffffffff92831660208501529116604083015266ffffffffffffff16606082015260800190565b6103a361039e366004612148565b610c54565b6040516101fa929190612438565b6103b9610d22565b6040516101fa91906123ea565b6004546101d99073ffffffffffffffffffffffffffffffffffffffff1681565b6101a4610dae565b6103b9610eab565b6101a46112f5565b61040760055481565b6040519081526020016101fa565b60005473ffffffffffffffffffffffffffffffffffffffff166101d9565b6101a4610441366004611fc2565b611305565b6101a4610454366004611fc2565b6113c8565b6101a4610467366004611fc2565b61148b565b6101a461047a366004612080565b61154e565b6101a461048d366004611fc2565b611b36565b61049a611b47565b73ffffffffffffffffffffffffffffffffffffffff81166104ba57600080fd5b6040805183815273ffffffffffffffffffffffffffffffffffffffff831660208201527f6141b54b56b8a52a8c6f5cd2a857f6117b18ffbf4d46bd3106f300a839cbf5ea910160405180910390a16003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018590529091169063a9059cbb90604401602060405180830381600087803b15801561057c57600080fd5b505af1158015610590573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105b49190612126565b505050565b6105c1611b47565b84831415806105d05750848114155b15610607576040517f3869bbe600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600680548060200260200160405190810160405280929190818152602001828054801561068957602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116106445790505b5050505050905060005b8151811015610717576000600760008484815181106106b4576106b4612565565b60209081029190910181015167ffffffffffffffff16825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061070f816124fd565b915050610693565b5060005b86811015610ab5576007600089898481811061073957610739612565565b905060200201602081019061074e919061221c565b67ffffffffffffffff16815260208101919091526040016000205460ff16156107dc5787878281811061078357610783612565565b9050602002016020810190610798919061221c565b6040517f7ceeb8ad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024015b60405180910390fd5b8787828181106107ee576107ee612565565b9050602002016020810190610803919061221c565b67ffffffffffffffff16610843576040517f3869bbe600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85858281811061085557610855612565565b905060200201602081019061086a9190612239565b6bffffffffffffffffffffffff1684848381811061088a5761088a612565565b905060200201602081019061089f9190612239565b6bffffffffffffffffffffffff16116108e4576040517f3869bbe600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806080016040528060011515815260200187878481811061090a5761090a612565565b905060200201602081019061091f9190612239565b6bffffffffffffffffffffffff16815260200185858481811061094457610944612565565b90506020020160208101906109599190612239565b6bffffffffffffffffffffffff168152602001600066ffffffffffffff16815250600760008a8a8581811061099057610990612565565b90506020020160208101906109a5919061221c565b67ffffffffffffffff168152602080820192909252604090810160002083518154938501519285015160609095015166ffffffffffffff167901000000000000000000000000000000000000000000000000000278ffffffffffffffffffffffffffffffffffffffffffffffffff6bffffffffffffffffffffffff9687166d010000000000000000000000000002166cffffffffffffffffffffffffff96909416610100027fffffffffffffffffffffffffffffffffffffff000000000000000000000000ff921515929092167fffffffffffffffffffffffffffffffffffffff000000000000000000000000009095169490941717939093161717905580610aad816124fd565b91505061071b565b50610ac260068888611ea3565b5050505050505050565b610ad4611b47565b610adc611bc8565b565b610ae6611b47565b60055460408051918252602082018390527f04330086c73b1fe1e13cd47a61c692e7c4399b5de08ed94b7ab824684af09323910160405180910390a1600555565b60045473ffffffffffffffffffffffffffffffffffffffff163314610b78576040517fd3a6803400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015474010000000000000000000000000000000000000000900460ff1615610bfd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a207061757365640000000000000000000000000000000060448201526064016107d3565b6000610c0b82840184612080565b90506105b48161154e565b60068181548110610c2657600080fd5b9060005260206000209060049182820401919006600802915054906101000a900467ffffffffffffffff1681565b60006060610c7d60015460ff740100000000000000000000000000000000000000009091041690565b15610ce4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a207061757365640000000000000000000000000000000060448201526064016107d3565b6000610cee610eab565b90506000815111925080604051602001610d0891906123ea565b6040516020818303038152906040529150505b9250929050565b60606006805480602002602001604051908101604052809291908181526020018280548015610da457602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020019060080190602082600701049283019260010382029150808411610d5f5790505b5050505050905090565b60015473ffffffffffffffffffffffffffffffffffffffff163314610e2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016107d3565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b606060006006805480602002602001604051908101604052809291908181526020018280548015610f2f57602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020019060080190602082600701049283019260010382029150808411610eea5790505b505050505090506000815167ffffffffffffffff811115610f5257610f52612594565b604051908082528060200260200182016040528015610f7b578160200160208202803683370190505b506005546003546040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152929350600092839173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b158015610fee57600080fd5b505afa158015611002573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061102691906121d3565b60408051608081018252600080825260208201819052918101829052606081018290529192505b86518110156112dc576007600088838151811061106c5761106c612565565b60209081029190910181015167ffffffffffffffff16825281810192909252604090810160009081208251608081018452905460ff81161515825261010081046bffffffffffffffffffffffff908116958301959095526d010000000000000000000000000081049094169281019290925279010000000000000000000000000000000000000000000000000090920466ffffffffffffff166060820152600254895191945073ffffffffffffffffffffffffffffffffffffffff169063a47c7696908a908590811061114157611141612565565b60200260200101516040518263ffffffff1660e01b8152600401611175919067ffffffffffffffff91909116815260200190565b60006040518083038186803b15801561118d57600080fd5b505afa1580156111a1573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526111e79190810190612256565b50505090504285846060015166ffffffffffffff1661120691906124ce565b11158015611226575082604001516bffffffffffffffffffffffff168410155b8015611251575082602001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16105b156112c95787828151811061126857611268612565565b602002602001015187878151811061128257611282612565565b67ffffffffffffffff90921660209283029190910190910152856112a5816124fd565b96505082604001516bffffffffffffffffffffffff16846112c691906124e6565b93505b50806112d4816124fd565b91505061104d565b5085518410156112ea578385525b509295945050505050565b6112fd611b47565b610adc611cc1565b61130d611b47565b73ffffffffffffffffffffffffffffffffffffffff811661132d57600080fd5b6004546040805173ffffffffffffffffffffffffffffffffffffffff928316815291831660208301527fb732223055abcde751d7a24272ffc8a3aa571cb72b443969a4199b7ecd59f8b9910160405180910390a1600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6113d0611b47565b73ffffffffffffffffffffffffffffffffffffffff81166113f057600080fd5b6002546040805173ffffffffffffffffffffffffffffffffffffffff928316815291831660208301527f4490e5bf8542a3633d8f268c3733706aa29acea979152e5fe5befd9cc7ad43a9910160405180910390a1600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b611493611b47565b73ffffffffffffffffffffffffffffffffffffffff81166114b357600080fd5b6003546040805173ffffffffffffffffffffffffffffffffffffffff928316815291831660208301527fee7e95e098f422f231397e2532a2752a013a51f4122cee6c30b18f930cda91dc910160405180910390a1600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015474010000000000000000000000000000000000000000900460ff16156115d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a207061757365640000000000000000000000000000000060448201526064016107d3565b6005546003546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561164057600080fd5b505afa158015611654573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061167891906121d3565b60408051608081018252600080825260208201819052918101829052606081018290529192505b8451811015611b2e57600760008683815181106116be576116be612565565b60209081029190910181015167ffffffffffffffff16825281810192909252604090810160009081208251608081018452905460ff81161515825261010081046bffffffffffffffffffffffff908116958301959095526d010000000000000000000000000081049094169281019290925279010000000000000000000000000000000000000000000000000090920466ffffffffffffff166060820152600254875191945073ffffffffffffffffffffffffffffffffffffffff169063a47c76969088908590811061179357611793612565565b60200260200101516040518263ffffffff1660e01b81526004016117c7919067ffffffffffffffff91909116815260200190565b60006040518083038186803b1580156117df57600080fd5b505afa1580156117f3573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526118399190810190612256565b5050845191925050801561186557504285846060015166ffffffffffffff1661186291906124ce565b11155b8015611890575082602001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16105b80156118ae575082604001516bffffffffffffffffffffffff168410155b15611ad6576003546002546040850151885160009373ffffffffffffffffffffffffffffffffffffffff90811693634000aea0939116918b90889081106118f7576118f7612565565b602002602001015160405160200161191f919067ffffffffffffffff91909116815260200190565b6040516020818303038152906040526040518463ffffffff1660e01b815260040161194c9392919061239e565b602060405180830381600087803b15801561196657600080fd5b505af115801561197a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061199e9190612126565b90508015611a835742600760008986815181106119bd576119bd612565565b602002602001015167ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060000160196101000a81548166ffffffffffffff021916908366ffffffffffffff16021790555083604001516bffffffffffffffffffffffff1685611a2c91906124e6565b9450868381518110611a4057611a40612565565b602002602001015167ffffffffffffffff167fef9c49dfa5fd8a638d79bc4a4c1edfce3d6c0a30a86e1273de3bafa32a5029cc60405160405180910390a2611ad4565b868381518110611a9557611a95612565565b602002602001015167ffffffffffffffff167fd6fe53f9994bdd53d2797bd6980c9c0004d7124a8e334de87dbe36d32cd0180160405160405180910390a25b505b61d6d85a1015611b1b576040518281527f8cc56b4ad3a81fec179b269a0784cb483821e9a835ec2b23594495f305bf77e19060200160405180910390a1505050505050565b5080611b26816124fd565b91505061169f565b505050505b50565b611b3e611b47565b611b3381611dad565b60005473ffffffffffffffffffffffffffffffffffffffff163314610adc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016107d3565b60015474010000000000000000000000000000000000000000900460ff16611c4c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f742070617573656400000000000000000000000060448201526064016107d3565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b60015474010000000000000000000000000000000000000000900460ff1615611d46576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a207061757365640000000000000000000000000000000060448201526064016107d3565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258611c973390565b73ffffffffffffffffffffffffffffffffffffffff8116331415611e2d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016107d3565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090600301600490048101928215611f585791602002820160005b83821115611f2257833567ffffffffffffffff1683826101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055509260200192600801602081600701049283019260010302611ecc565b8015611f565782816101000a81549067ffffffffffffffff0219169055600801602081600701049283019260010302611f22565b505b50611f64929150611f68565b5090565b5b80821115611f645760008155600101611f69565b60008083601f840112611f8f57600080fd5b50813567ffffffffffffffff811115611fa757600080fd5b6020830191508360208260051b8501011115610d1b57600080fd5b600060208284031215611fd457600080fd5b8135611fdf816125c3565b9392505050565b60008060008060008060608789031215611fff57600080fd5b863567ffffffffffffffff8082111561201757600080fd5b6120238a838b01611f7d565b9098509650602089013591508082111561203c57600080fd5b6120488a838b01611f7d565b9096509450604089013591508082111561206157600080fd5b5061206e89828a01611f7d565b979a9699509497509295939492505050565b6000602080838503121561209357600080fd5b823567ffffffffffffffff8111156120aa57600080fd5b8301601f810185136120bb57600080fd5b80356120ce6120c9826124aa565b61245b565b80828252848201915084840188868560051b87010111156120ee57600080fd5b600094505b8385101561211a578035612106816125e5565b8352600194909401939185019185016120f3565b50979650505050505050565b60006020828403121561213857600080fd5b81518015158114611fdf57600080fd5b6000806020838503121561215b57600080fd5b823567ffffffffffffffff8082111561217357600080fd5b818501915085601f83011261218757600080fd5b81358181111561219657600080fd5b8660208285010111156121a857600080fd5b60209290920196919550909350505050565b6000602082840312156121cc57600080fd5b5035919050565b6000602082840312156121e557600080fd5b5051919050565b600080604083850312156121ff57600080fd5b823591506020830135612211816125c3565b809150509250929050565b60006020828403121561222e57600080fd5b8135611fdf816125e5565b60006020828403121561224b57600080fd5b8135611fdf816125fb565b6000806000806080858703121561226c57600080fd5b8451612277816125fb565b8094505060208086015161228a816125e5565b604087015190945061229b816125c3565b606087015190935067ffffffffffffffff8111156122b857600080fd5b8601601f810188136122c957600080fd5b80516122d76120c9826124aa565b8082825284820191508484018b868560051b87010111156122f757600080fd5b600094505b8385101561232357805161230f816125c3565b8352600194909401939185019185016122fc565b50979a9699509497505050505050565b6000815180845260005b818110156123595760208185018101518683018201520161233d565b8181111561236b576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff831660208201526060604082015260006123e16060830184612333565b95945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561242c57835167ffffffffffffffff1683529284019291840191600101612406565b50909695505050505050565b82151581526040602082015260006124536040830184612333565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156124a2576124a2612594565b604052919050565b600067ffffffffffffffff8211156124c4576124c4612594565b5060051b60200190565b600082198211156124e1576124e1612536565b500190565b6000828210156124f8576124f8612536565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561252f5761252f612536565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff81168114611b3357600080fd5b67ffffffffffffffff81168114611b3357600080fd5b6bffffffffffffffffffffffff81168114611b3357600080fdfea164736f6c6343000806000a",
}

var VRFSubscriptionBalanceMonitorABI = VRFSubscriptionBalanceMonitorMetaData.ABI

var VRFSubscriptionBalanceMonitorBin = VRFSubscriptionBalanceMonitorMetaData.Bin

func DeployVRFSubscriptionBalanceMonitor(auth *bind.TransactOpts, backend bind.ContractBackend, linkTokenAddress common.Address, coordinatorAddress common.Address, keeperRegistryAddress common.Address, minWaitPeriodSeconds *big.Int) (common.Address, *types.Transaction, *VRFSubscriptionBalanceMonitor, error) {
	parsed, err := VRFSubscriptionBalanceMonitorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFSubscriptionBalanceMonitorBin), backend, linkTokenAddress, coordinatorAddress, keeperRegistryAddress, minWaitPeriodSeconds)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFSubscriptionBalanceMonitor{VRFSubscriptionBalanceMonitorCaller: VRFSubscriptionBalanceMonitorCaller{contract: contract}, VRFSubscriptionBalanceMonitorTransactor: VRFSubscriptionBalanceMonitorTransactor{contract: contract}, VRFSubscriptionBalanceMonitorFilterer: VRFSubscriptionBalanceMonitorFilterer{contract: contract}}, nil
}

type VRFSubscriptionBalanceMonitor struct {
	address common.Address
	abi     abi.ABI
	VRFSubscriptionBalanceMonitorCaller
	VRFSubscriptionBalanceMonitorTransactor
	VRFSubscriptionBalanceMonitorFilterer
}

type VRFSubscriptionBalanceMonitorCaller struct {
	contract *bind.BoundContract
}

type VRFSubscriptionBalanceMonitorTransactor struct {
	contract *bind.BoundContract
}

type VRFSubscriptionBalanceMonitorFilterer struct {
	contract *bind.BoundContract
}

type VRFSubscriptionBalanceMonitorSession struct {
	Contract     *VRFSubscriptionBalanceMonitor
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFSubscriptionBalanceMonitorCallerSession struct {
	Contract *VRFSubscriptionBalanceMonitorCaller
	CallOpts bind.CallOpts
}

type VRFSubscriptionBalanceMonitorTransactorSession struct {
	Contract     *VRFSubscriptionBalanceMonitorTransactor
	TransactOpts bind.TransactOpts
}

type VRFSubscriptionBalanceMonitorRaw struct {
	Contract *VRFSubscriptionBalanceMonitor
}

type VRFSubscriptionBalanceMonitorCallerRaw struct {
	Contract *VRFSubscriptionBalanceMonitorCaller
}

type VRFSubscriptionBalanceMonitorTransactorRaw struct {
	Contract *VRFSubscriptionBalanceMonitorTransactor
}

func NewVRFSubscriptionBalanceMonitor(address common.Address, backend bind.ContractBackend) (*VRFSubscriptionBalanceMonitor, error) {
	abi, err := abi.JSON(strings.NewReader(VRFSubscriptionBalanceMonitorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFSubscriptionBalanceMonitor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitor{address: address, abi: abi, VRFSubscriptionBalanceMonitorCaller: VRFSubscriptionBalanceMonitorCaller{contract: contract}, VRFSubscriptionBalanceMonitorTransactor: VRFSubscriptionBalanceMonitorTransactor{contract: contract}, VRFSubscriptionBalanceMonitorFilterer: VRFSubscriptionBalanceMonitorFilterer{contract: contract}}, nil
}

func NewVRFSubscriptionBalanceMonitorCaller(address common.Address, caller bind.ContractCaller) (*VRFSubscriptionBalanceMonitorCaller, error) {
	contract, err := bindVRFSubscriptionBalanceMonitor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorCaller{contract: contract}, nil
}

func NewVRFSubscriptionBalanceMonitorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFSubscriptionBalanceMonitorTransactor, error) {
	contract, err := bindVRFSubscriptionBalanceMonitor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorTransactor{contract: contract}, nil
}

func NewVRFSubscriptionBalanceMonitorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFSubscriptionBalanceMonitorFilterer, error) {
	contract, err := bindVRFSubscriptionBalanceMonitor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorFilterer{contract: contract}, nil
}

func bindVRFSubscriptionBalanceMonitor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFSubscriptionBalanceMonitorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFSubscriptionBalanceMonitor.Contract.VRFSubscriptionBalanceMonitorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.VRFSubscriptionBalanceMonitorTransactor.contract.Transfer(opts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.VRFSubscriptionBalanceMonitorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFSubscriptionBalanceMonitor.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.contract.Transfer(opts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.contract.Transact(opts, method, params...)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) COORDINATOR() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.COORDINATOR(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.COORDINATOR(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) LINKTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "LINKTOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) LINKTOKEN() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.LINKTOKEN(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) LINKTOKEN() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.LINKTOKEN(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) CheckUpkeep(opts *bind.CallOpts, arg0 []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "checkUpkeep", arg0)

	outstruct := new(CheckUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) CheckUpkeep(arg0 []byte) (CheckUpkeep,

	error) {
	return _VRFSubscriptionBalanceMonitor.Contract.CheckUpkeep(&_VRFSubscriptionBalanceMonitor.CallOpts, arg0)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) CheckUpkeep(arg0 []byte) (CheckUpkeep,

	error) {
	return _VRFSubscriptionBalanceMonitor.Contract.CheckUpkeep(&_VRFSubscriptionBalanceMonitor.CallOpts, arg0)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) GetSubscriptionInfo(opts *bind.CallOpts, subscriptionId uint64) (GetSubscriptionInfo,

	error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "getSubscriptionInfo", subscriptionId)

	outstruct := new(GetSubscriptionInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsActive = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.MinBalanceJuels = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TopUpAmountJuels = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastTopUpTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) GetSubscriptionInfo(subscriptionId uint64) (GetSubscriptionInfo,

	error) {
	return _VRFSubscriptionBalanceMonitor.Contract.GetSubscriptionInfo(&_VRFSubscriptionBalanceMonitor.CallOpts, subscriptionId)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) GetSubscriptionInfo(subscriptionId uint64) (GetSubscriptionInfo,

	error) {
	return _VRFSubscriptionBalanceMonitor.Contract.GetSubscriptionInfo(&_VRFSubscriptionBalanceMonitor.CallOpts, subscriptionId)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) GetUnderfundedSubscriptions(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "getUnderfundedSubscriptions")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) GetUnderfundedSubscriptions() ([]uint64, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.GetUnderfundedSubscriptions(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) GetUnderfundedSubscriptions() ([]uint64, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.GetUnderfundedSubscriptions(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) GetWatchList(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "getWatchList")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) GetWatchList() ([]uint64, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.GetWatchList(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) GetWatchList() ([]uint64, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.GetWatchList(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) Owner() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Owner(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) Owner() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Owner(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) Paused() (bool, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Paused(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) Paused() (bool, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Paused(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) SKeeperRegistryAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "s_keeperRegistryAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SKeeperRegistryAddress() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SKeeperRegistryAddress(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) SKeeperRegistryAddress() (common.Address, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SKeeperRegistryAddress(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) SMinWaitPeriodSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "s_minWaitPeriodSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SMinWaitPeriodSeconds() (*big.Int, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SMinWaitPeriodSeconds(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) SMinWaitPeriodSeconds() (*big.Int, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SMinWaitPeriodSeconds(&_VRFSubscriptionBalanceMonitor.CallOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCaller) SWatchList(opts *bind.CallOpts, arg0 *big.Int) (uint64, error) {
	var out []interface{}
	err := _VRFSubscriptionBalanceMonitor.contract.Call(opts, &out, "s_watchList", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SWatchList(arg0 *big.Int) (uint64, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SWatchList(&_VRFSubscriptionBalanceMonitor.CallOpts, arg0)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorCallerSession) SWatchList(arg0 *big.Int) (uint64, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SWatchList(&_VRFSubscriptionBalanceMonitor.CallOpts, arg0)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "acceptOwnership")
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.AcceptOwnership(&_VRFSubscriptionBalanceMonitor.TransactOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.AcceptOwnership(&_VRFSubscriptionBalanceMonitor.TransactOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "pause")
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) Pause() (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Pause(&_VRFSubscriptionBalanceMonitor.TransactOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) Pause() (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Pause(&_VRFSubscriptionBalanceMonitor.TransactOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "performUpkeep", performData)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.PerformUpkeep(&_VRFSubscriptionBalanceMonitor.TransactOpts, performData)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.PerformUpkeep(&_VRFSubscriptionBalanceMonitor.TransactOpts, performData)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) SetKeeperRegistryAddress(opts *bind.TransactOpts, keeperRegistryAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "setKeeperRegistryAddress", keeperRegistryAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SetKeeperRegistryAddress(keeperRegistryAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetKeeperRegistryAddress(&_VRFSubscriptionBalanceMonitor.TransactOpts, keeperRegistryAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) SetKeeperRegistryAddress(keeperRegistryAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetKeeperRegistryAddress(&_VRFSubscriptionBalanceMonitor.TransactOpts, keeperRegistryAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) SetLinkTokenAddress(opts *bind.TransactOpts, linkTokenAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "setLinkTokenAddress", linkTokenAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SetLinkTokenAddress(linkTokenAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetLinkTokenAddress(&_VRFSubscriptionBalanceMonitor.TransactOpts, linkTokenAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) SetLinkTokenAddress(linkTokenAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetLinkTokenAddress(&_VRFSubscriptionBalanceMonitor.TransactOpts, linkTokenAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) SetMinWaitPeriodSeconds(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "setMinWaitPeriodSeconds", period)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SetMinWaitPeriodSeconds(period *big.Int) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetMinWaitPeriodSeconds(&_VRFSubscriptionBalanceMonitor.TransactOpts, period)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) SetMinWaitPeriodSeconds(period *big.Int) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetMinWaitPeriodSeconds(&_VRFSubscriptionBalanceMonitor.TransactOpts, period)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) SetVRFCoordinatorV2Address(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "setVRFCoordinatorV2Address", coordinatorAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SetVRFCoordinatorV2Address(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetVRFCoordinatorV2Address(&_VRFSubscriptionBalanceMonitor.TransactOpts, coordinatorAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) SetVRFCoordinatorV2Address(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetVRFCoordinatorV2Address(&_VRFSubscriptionBalanceMonitor.TransactOpts, coordinatorAddress)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) SetWatchList(opts *bind.TransactOpts, subscriptionIds []uint64, minBalancesJuels []*big.Int, topUpAmountsJuels []*big.Int) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "setWatchList", subscriptionIds, minBalancesJuels, topUpAmountsJuels)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) SetWatchList(subscriptionIds []uint64, minBalancesJuels []*big.Int, topUpAmountsJuels []*big.Int) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetWatchList(&_VRFSubscriptionBalanceMonitor.TransactOpts, subscriptionIds, minBalancesJuels, topUpAmountsJuels)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) SetWatchList(subscriptionIds []uint64, minBalancesJuels []*big.Int, topUpAmountsJuels []*big.Int) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.SetWatchList(&_VRFSubscriptionBalanceMonitor.TransactOpts, subscriptionIds, minBalancesJuels, topUpAmountsJuels)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) TopUp(opts *bind.TransactOpts, needsFunding []uint64) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "topUp", needsFunding)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) TopUp(needsFunding []uint64) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.TopUp(&_VRFSubscriptionBalanceMonitor.TransactOpts, needsFunding)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) TopUp(needsFunding []uint64) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.TopUp(&_VRFSubscriptionBalanceMonitor.TransactOpts, needsFunding)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.TransferOwnership(&_VRFSubscriptionBalanceMonitor.TransactOpts, to)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.TransferOwnership(&_VRFSubscriptionBalanceMonitor.TransactOpts, to)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "unpause")
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) Unpause() (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Unpause(&_VRFSubscriptionBalanceMonitor.TransactOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) Unpause() (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Unpause(&_VRFSubscriptionBalanceMonitor.TransactOpts)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, payee common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.contract.Transact(opts, "withdraw", amount, payee)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorSession) Withdraw(amount *big.Int, payee common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Withdraw(&_VRFSubscriptionBalanceMonitor.TransactOpts, amount, payee)
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorTransactorSession) Withdraw(amount *big.Int, payee common.Address) (*types.Transaction, error) {
	return _VRFSubscriptionBalanceMonitor.Contract.Withdraw(&_VRFSubscriptionBalanceMonitor.TransactOpts, amount, payee)
}

type VRFSubscriptionBalanceMonitorFundsAddedIterator struct {
	Event *VRFSubscriptionBalanceMonitorFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorFundsAdded)
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
		it.Event = new(VRFSubscriptionBalanceMonitorFundsAdded)
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

func (it *VRFSubscriptionBalanceMonitorFundsAddedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorFundsAdded struct {
	AmountAdded *big.Int
	NewBalance  *big.Int
	Sender      common.Address
	Raw         types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterFundsAdded(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorFundsAddedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorFundsAddedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorFundsAdded) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorFundsAdded)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseFundsAdded(log types.Log) (*VRFSubscriptionBalanceMonitorFundsAdded, error) {
	event := new(VRFSubscriptionBalanceMonitorFundsAdded)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorFundsWithdrawnIterator struct {
	Event *VRFSubscriptionBalanceMonitorFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorFundsWithdrawn)
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
		it.Event = new(VRFSubscriptionBalanceMonitorFundsWithdrawn)
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

func (it *VRFSubscriptionBalanceMonitorFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorFundsWithdrawn struct {
	AmountWithdrawn *big.Int
	Payee           common.Address
	Raw             types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorFundsWithdrawnIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "FundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorFundsWithdrawnIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "FundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorFundsWithdrawn)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseFundsWithdrawn(log types.Log) (*VRFSubscriptionBalanceMonitorFundsWithdrawn, error) {
	event := new(VRFSubscriptionBalanceMonitorFundsWithdrawn)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator struct {
	Event *VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated)
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
		it.Event = new(VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated)
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

func (it *VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterKeeperRegistryAddressUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "KeeperRegistryAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "KeeperRegistryAddressUpdated", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchKeeperRegistryAddressUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "KeeperRegistryAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "KeeperRegistryAddressUpdated", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseKeeperRegistryAddressUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated, error) {
	event := new(VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "KeeperRegistryAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator struct {
	Event *VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated)
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
		it.Event = new(VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated)
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

func (it *VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterLinkTokenAddressUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "LinkTokenAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "LinkTokenAddressUpdated", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchLinkTokenAddressUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "LinkTokenAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "LinkTokenAddressUpdated", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseLinkTokenAddressUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated, error) {
	event := new(VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "LinkTokenAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator struct {
	Event *VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated)
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
		it.Event = new(VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated)
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

func (it *VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated struct {
	OldMinWaitPeriod *big.Int
	NewMinWaitPeriod *big.Int
	Raw              types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterMinWaitPeriodUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "MinWaitPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "MinWaitPeriodUpdated", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchMinWaitPeriodUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "MinWaitPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "MinWaitPeriodUpdated", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseMinWaitPeriodUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated, error) {
	event := new(VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "MinWaitPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorOutOfGasIterator struct {
	Event *VRFSubscriptionBalanceMonitorOutOfGas

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorOutOfGasIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorOutOfGas)
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
		it.Event = new(VRFSubscriptionBalanceMonitorOutOfGas)
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

func (it *VRFSubscriptionBalanceMonitorOutOfGasIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorOutOfGasIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorOutOfGas struct {
	LastId *big.Int
	Raw    types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterOutOfGas(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorOutOfGasIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "OutOfGas")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorOutOfGasIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "OutOfGas", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchOutOfGas(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorOutOfGas) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "OutOfGas")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorOutOfGas)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "OutOfGas", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseOutOfGas(log types.Log) (*VRFSubscriptionBalanceMonitorOutOfGas, error) {
	event := new(VRFSubscriptionBalanceMonitorOutOfGas)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "OutOfGas", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator struct {
	Event *VRFSubscriptionBalanceMonitorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorOwnershipTransferRequested)
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
		it.Event = new(VRFSubscriptionBalanceMonitorOwnershipTransferRequested)
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

func (it *VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorOwnershipTransferRequested)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFSubscriptionBalanceMonitorOwnershipTransferRequested, error) {
	event := new(VRFSubscriptionBalanceMonitorOwnershipTransferRequested)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorOwnershipTransferredIterator struct {
	Event *VRFSubscriptionBalanceMonitorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorOwnershipTransferred)
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
		it.Event = new(VRFSubscriptionBalanceMonitorOwnershipTransferred)
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

func (it *VRFSubscriptionBalanceMonitorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFSubscriptionBalanceMonitorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorOwnershipTransferredIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorOwnershipTransferred)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseOwnershipTransferred(log types.Log) (*VRFSubscriptionBalanceMonitorOwnershipTransferred, error) {
	event := new(VRFSubscriptionBalanceMonitorOwnershipTransferred)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorPausedIterator struct {
	Event *VRFSubscriptionBalanceMonitorPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorPaused)
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
		it.Event = new(VRFSubscriptionBalanceMonitorPaused)
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

func (it *VRFSubscriptionBalanceMonitorPausedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterPaused(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorPausedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorPausedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorPaused) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorPaused)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParsePaused(log types.Log) (*VRFSubscriptionBalanceMonitorPaused, error) {
	event := new(VRFSubscriptionBalanceMonitorPaused)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorTopUpFailedIterator struct {
	Event *VRFSubscriptionBalanceMonitorTopUpFailed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorTopUpFailedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorTopUpFailed)
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
		it.Event = new(VRFSubscriptionBalanceMonitorTopUpFailed)
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

func (it *VRFSubscriptionBalanceMonitorTopUpFailedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorTopUpFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorTopUpFailed struct {
	SubscriptionId uint64
	Raw            types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterTopUpFailed(opts *bind.FilterOpts, subscriptionId []uint64) (*VRFSubscriptionBalanceMonitorTopUpFailedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "TopUpFailed", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorTopUpFailedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "TopUpFailed", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchTopUpFailed(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorTopUpFailed, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "TopUpFailed", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorTopUpFailed)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "TopUpFailed", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseTopUpFailed(log types.Log) (*VRFSubscriptionBalanceMonitorTopUpFailed, error) {
	event := new(VRFSubscriptionBalanceMonitorTopUpFailed)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "TopUpFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorTopUpSucceededIterator struct {
	Event *VRFSubscriptionBalanceMonitorTopUpSucceeded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorTopUpSucceededIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorTopUpSucceeded)
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
		it.Event = new(VRFSubscriptionBalanceMonitorTopUpSucceeded)
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

func (it *VRFSubscriptionBalanceMonitorTopUpSucceededIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorTopUpSucceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorTopUpSucceeded struct {
	SubscriptionId uint64
	Raw            types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterTopUpSucceeded(opts *bind.FilterOpts, subscriptionId []uint64) (*VRFSubscriptionBalanceMonitorTopUpSucceededIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "TopUpSucceeded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorTopUpSucceededIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "TopUpSucceeded", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchTopUpSucceeded(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorTopUpSucceeded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "TopUpSucceeded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorTopUpSucceeded)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "TopUpSucceeded", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseTopUpSucceeded(log types.Log) (*VRFSubscriptionBalanceMonitorTopUpSucceeded, error) {
	event := new(VRFSubscriptionBalanceMonitorTopUpSucceeded)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "TopUpSucceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorUnpausedIterator struct {
	Event *VRFSubscriptionBalanceMonitorUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorUnpaused)
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
		it.Event = new(VRFSubscriptionBalanceMonitorUnpaused)
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

func (it *VRFSubscriptionBalanceMonitorUnpausedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorUnpausedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorUnpausedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorUnpaused) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorUnpaused)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseUnpaused(log types.Log) (*VRFSubscriptionBalanceMonitorUnpaused, error) {
	event := new(VRFSubscriptionBalanceMonitorUnpaused)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator struct {
	Event *VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated)
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
		it.Event = new(VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated)
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

func (it *VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator) Error() error {
	return it.fail
}

func (it *VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) FilterVRFCoordinatorV2AddressUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.FilterLogs(opts, "VRFCoordinatorV2AddressUpdated")
	if err != nil {
		return nil, err
	}
	return &VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator{contract: _VRFSubscriptionBalanceMonitor.contract, event: "VRFCoordinatorV2AddressUpdated", logs: logs, sub: sub}, nil
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) WatchVRFCoordinatorV2AddressUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated) (event.Subscription, error) {

	logs, sub, err := _VRFSubscriptionBalanceMonitor.contract.WatchLogs(opts, "VRFCoordinatorV2AddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated)
				if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "VRFCoordinatorV2AddressUpdated", log); err != nil {
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

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitorFilterer) ParseVRFCoordinatorV2AddressUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated, error) {
	event := new(VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated)
	if err := _VRFSubscriptionBalanceMonitor.contract.UnpackLog(event, "VRFCoordinatorV2AddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CheckUpkeep struct {
	UpkeepNeeded bool
	PerformData  []byte
}
type GetSubscriptionInfo struct {
	IsActive           bool
	MinBalanceJuels    *big.Int
	TopUpAmountJuels   *big.Int
	LastTopUpTimestamp *big.Int
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitor) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFSubscriptionBalanceMonitor.abi.Events["FundsAdded"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseFundsAdded(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["FundsWithdrawn"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseFundsWithdrawn(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["KeeperRegistryAddressUpdated"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseKeeperRegistryAddressUpdated(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["LinkTokenAddressUpdated"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseLinkTokenAddressUpdated(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["MinWaitPeriodUpdated"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseMinWaitPeriodUpdated(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["OutOfGas"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseOutOfGas(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseOwnershipTransferRequested(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["OwnershipTransferred"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseOwnershipTransferred(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["Paused"].ID:
		return _VRFSubscriptionBalanceMonitor.ParsePaused(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["TopUpFailed"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseTopUpFailed(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["TopUpSucceeded"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseTopUpSucceeded(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["Unpaused"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseUnpaused(log)
	case _VRFSubscriptionBalanceMonitor.abi.Events["VRFCoordinatorV2AddressUpdated"].ID:
		return _VRFSubscriptionBalanceMonitor.ParseVRFCoordinatorV2AddressUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFSubscriptionBalanceMonitorFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xc6f3fb0fec49e4877342d4625d77a632541f55b7aae0f9d0b34c69b3478706dc")
}

func (VRFSubscriptionBalanceMonitorFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x6141b54b56b8a52a8c6f5cd2a857f6117b18ffbf4d46bd3106f300a839cbf5ea")
}

func (VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated) Topic() common.Hash {
	return common.HexToHash("0xb732223055abcde751d7a24272ffc8a3aa571cb72b443969a4199b7ecd59f8b9")
}

func (VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated) Topic() common.Hash {
	return common.HexToHash("0xee7e95e098f422f231397e2532a2752a013a51f4122cee6c30b18f930cda91dc")
}

func (VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated) Topic() common.Hash {
	return common.HexToHash("0x04330086c73b1fe1e13cd47a61c692e7c4399b5de08ed94b7ab824684af09323")
}

func (VRFSubscriptionBalanceMonitorOutOfGas) Topic() common.Hash {
	return common.HexToHash("0x8cc56b4ad3a81fec179b269a0784cb483821e9a835ec2b23594495f305bf77e1")
}

func (VRFSubscriptionBalanceMonitorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFSubscriptionBalanceMonitorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFSubscriptionBalanceMonitorPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (VRFSubscriptionBalanceMonitorTopUpFailed) Topic() common.Hash {
	return common.HexToHash("0xd6fe53f9994bdd53d2797bd6980c9c0004d7124a8e334de87dbe36d32cd01801")
}

func (VRFSubscriptionBalanceMonitorTopUpSucceeded) Topic() common.Hash {
	return common.HexToHash("0xef9c49dfa5fd8a638d79bc4a4c1edfce3d6c0a30a86e1273de3bafa32a5029cc")
}

func (VRFSubscriptionBalanceMonitorUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated) Topic() common.Hash {
	return common.HexToHash("0x4490e5bf8542a3633d8f268c3733706aa29acea979152e5fe5befd9cc7ad43a9")
}

func (_VRFSubscriptionBalanceMonitor *VRFSubscriptionBalanceMonitor) Address() common.Address {
	return _VRFSubscriptionBalanceMonitor.address
}

type VRFSubscriptionBalanceMonitorInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINKTOKEN(opts *bind.CallOpts) (common.Address, error)

	CheckUpkeep(opts *bind.CallOpts, arg0 []byte) (CheckUpkeep,

		error)

	GetSubscriptionInfo(opts *bind.CallOpts, subscriptionId uint64) (GetSubscriptionInfo,

		error)

	GetUnderfundedSubscriptions(opts *bind.CallOpts) ([]uint64, error)

	GetWatchList(opts *bind.CallOpts) ([]uint64, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	SKeeperRegistryAddress(opts *bind.CallOpts) (common.Address, error)

	SMinWaitPeriodSeconds(opts *bind.CallOpts) (*big.Int, error)

	SWatchList(opts *bind.CallOpts, arg0 *big.Int) (uint64, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetKeeperRegistryAddress(opts *bind.TransactOpts, keeperRegistryAddress common.Address) (*types.Transaction, error)

	SetLinkTokenAddress(opts *bind.TransactOpts, linkTokenAddress common.Address) (*types.Transaction, error)

	SetMinWaitPeriodSeconds(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error)

	SetVRFCoordinatorV2Address(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error)

	SetWatchList(opts *bind.TransactOpts, subscriptionIds []uint64, minBalancesJuels []*big.Int, topUpAmountsJuels []*big.Int) (*types.Transaction, error)

	TopUp(opts *bind.TransactOpts, needsFunding []uint64) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, amount *big.Int, payee common.Address) (*types.Transaction, error)

	FilterFundsAdded(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorFundsAdded) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*VRFSubscriptionBalanceMonitorFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorFundsWithdrawn) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*VRFSubscriptionBalanceMonitorFundsWithdrawn, error)

	FilterKeeperRegistryAddressUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdatedIterator, error)

	WatchKeeperRegistryAddressUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated) (event.Subscription, error)

	ParseKeeperRegistryAddressUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorKeeperRegistryAddressUpdated, error)

	FilterLinkTokenAddressUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorLinkTokenAddressUpdatedIterator, error)

	WatchLinkTokenAddressUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated) (event.Subscription, error)

	ParseLinkTokenAddressUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorLinkTokenAddressUpdated, error)

	FilterMinWaitPeriodUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorMinWaitPeriodUpdatedIterator, error)

	WatchMinWaitPeriodUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated) (event.Subscription, error)

	ParseMinWaitPeriodUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorMinWaitPeriodUpdated, error)

	FilterOutOfGas(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorOutOfGasIterator, error)

	WatchOutOfGas(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorOutOfGas) (event.Subscription, error)

	ParseOutOfGas(log types.Log) (*VRFSubscriptionBalanceMonitorOutOfGas, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFSubscriptionBalanceMonitorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFSubscriptionBalanceMonitorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFSubscriptionBalanceMonitorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFSubscriptionBalanceMonitorOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*VRFSubscriptionBalanceMonitorPaused, error)

	FilterTopUpFailed(opts *bind.FilterOpts, subscriptionId []uint64) (*VRFSubscriptionBalanceMonitorTopUpFailedIterator, error)

	WatchTopUpFailed(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorTopUpFailed, subscriptionId []uint64) (event.Subscription, error)

	ParseTopUpFailed(log types.Log) (*VRFSubscriptionBalanceMonitorTopUpFailed, error)

	FilterTopUpSucceeded(opts *bind.FilterOpts, subscriptionId []uint64) (*VRFSubscriptionBalanceMonitorTopUpSucceededIterator, error)

	WatchTopUpSucceeded(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorTopUpSucceeded, subscriptionId []uint64) (event.Subscription, error)

	ParseTopUpSucceeded(log types.Log) (*VRFSubscriptionBalanceMonitorTopUpSucceeded, error)

	FilterUnpaused(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*VRFSubscriptionBalanceMonitorUnpaused, error)

	FilterVRFCoordinatorV2AddressUpdated(opts *bind.FilterOpts) (*VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdatedIterator, error)

	WatchVRFCoordinatorV2AddressUpdated(opts *bind.WatchOpts, sink chan<- *VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated) (event.Subscription, error)

	ParseVRFCoordinatorV2AddressUpdated(log types.Log) (*VRFSubscriptionBalanceMonitorVRFCoordinatorV2AddressUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
