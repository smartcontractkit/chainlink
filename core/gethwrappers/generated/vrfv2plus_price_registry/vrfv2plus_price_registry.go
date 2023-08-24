// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_price_registry

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

var VRFV2PlusPriceRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ethUSDFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"ethUSD\",\"type\":\"int256\"}],\"name\":\"InvalidEthUSDPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"got\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"expected1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"expected2\",\"type\":\"address\"}],\"name\":\"InvalidInput\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkUSD\",\"type\":\"int256\"}],\"name\":\"InvalidLinkUSDPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feed\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"}],\"name\":\"InvalidUSDPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitEth\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeed\",\"type\":\"address\"}],\"name\":\"EthUSDFeedSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeed\",\"type\":\"address\"}],\"name\":\"LinkEthFeedSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeed\",\"type\":\"address\"}],\"name\":\"LinkUSDFeedSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"wrapperGasOverhead\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"coordinatorGasOverhead\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentTxSizeBytes\",\"type\":\"uint32\"}],\"name\":\"WrapperConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"USD_FEE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitGas\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"calculatePaymentAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"calculateRequestPriceNativeWrapper\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"calculateRequestPriceWrapper\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_requestGasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateRequestPriceNativeWrapper\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_requestGasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateRequestPriceWrapper\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_ethUSDFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackUSDPerUnitEth\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackUSDPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkETHFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkUSDFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_wrapperConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"coordinatorGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentTxSizeBytes\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitEth\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethUsdFeed\",\"type\":\"address\"}],\"name\":\"setETHUSDFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"name\":\"setLINKETHFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkUsdFeed\",\"type\":\"address\"}],\"name\":\"setLINKUSDFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"coordinatorGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentTxSizeBytes\",\"type\":\"uint32\"}],\"name\":\"setWrapperConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001e0f38038062001e0f8339810160408190526200003491620001cd565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000104565b5050600280546001600160a01b039586166001600160a01b0319918216179091556003805494861694821694909417909355506004805491909316911617905562000217565b6001600160a01b0381163314156200015f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c857600080fd5b919050565b600080600060608486031215620001e357600080fd5b620001ee84620001b0565b9250620001fe60208501620001b0565b91506200020e60408501620001b0565b90509250925092565b611be880620002276000396000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c80638da5cb5b116100d8578063de6a92481161008c578063e7ddbb8d11610066578063e7ddbb8d146103bf578063e993b4aa146103d2578063f2fde38b1461044157600080fd5b8063de6a924814610369578063e16ad7cf1461037c578063e6152d811461038f57600080fd5b8063bb0697a5116100bd578063bb0697a514610329578063cfbb957a14610343578063d1e28cec1461035657600080fd5b80638da5cb5b146102f8578063907d64591461031657600080fd5b806359392b6d1161013a578063723276d611610114578063723276d6146102c757806379ba5097146102e7578063835c0dfc146102ef57600080fd5b806359392b6d1461028b57806367c77a54146102945780636af6890d146102a757600080fd5b80630784e5d01161016b5780630784e5d0146101b8578063088070f5146101cb578063180a49091461024657600080fd5b80630396525514610187578063043bd6ae1461019c575b600080fd5b61019a6101953660046118c1565b610454565b005b6101a560075481565b6040519081526020015b60405180910390f35b61019a6101c6366004611734565b61056f565b6005546102119063ffffffff8082169164010000000081049091169064ffffffffff6801000000000000000082048116916d010000000000000000000000000090041684565b6040805163ffffffff958616815294909316602085015264ffffffffff918216928401929092521660608201526080016101af565b6004546102669073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101af565b6101a560095481565b61019a6102a2366004611734565b6105fe565b6003546102669073ffffffffffffffffffffffffffffffffffffffff1681565b6002546102669073ffffffffffffffffffffffffffffffffffffffff1681565b61019a610685565b6101a560085481565b60005473ffffffffffffffffffffffffffffffffffffffff16610266565b6101a561032436600461180b565b610787565b610331600881565b60405160ff90911681526020016101af565b6101a5610351366004611826565b6107d4565b6101a5610364366004611826565b610822565b61019a610377366004611850565b610871565b6101a561038a36600461180b565b610a6e565b6103a261039d366004611783565b610a83565b6040516bffffffffffffffffffffffff90911681526020016101af565b61019a6103cd366004611734565b610b01565b60065461040e9063ffffffff80821691640100000000810482169160ff6801000000000000000083041691690100000000000000000090041684565b6040805163ffffffff9586168152938516602085015260ff909216918301919091529190911660608201526080016101af565b61019a61044f366004611734565b610b88565b61045c610b9c565b604080516080808201835263ffffffff878116808452878216602080860182905260ff89168688018190529388166060968701819052600680547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001685176401000000008502177fffffffffffffffffffffffffffffffffffffff0000000000ffffffffffffffff166801000000000000000087027fffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffff16176901000000000000000000830217905587519384529083019190915294810191909152918201929092527fb76e5cdd0dc2bc17df6a1911db92060a11b138d8002f77267013944123a4475d910160405180910390a150505050565b610577610b9c565b6002805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f15f61b91e528d42be960613d5606dbf13df3ef988e6a097b8543c9a58b2b7fd891015b60405180910390a15050565b610606610b9c565b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527fe6a1e056cb2ec82c5f49294ff925bd5a0ab6a8ccbe8fdfdf7d9a333d9c12c50791016105f2565b60015473ffffffffffffffffffffffffffffffffffffffff16331461070b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6006546000906107ce9063ffffffff848116913a9180821691640100000000810482169169010000000000000000008204169068010000000000000000900460ff16610c1f565b92915050565b60065460009061081b9063ffffffff85811691859180821691640100000000810482169169010000000000000000008204169068010000000000000000900460ff16610c1f565b9392505050565b600061081b8363ffffffff1683610837610ce8565b60065463ffffffff80821691640100000000810482169169010000000000000000008204169068010000000000000000900460ff16610dd2565b610879610b9c565b600085136108b6576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101869052602401610702565b600084136108f3576040517f599d67e300000000000000000000000000000000000000000000000000000000815260048101859052602401610702565b60008313610930576040517f25b2499f00000000000000000000000000000000000000000000000000000000815260048101849052602401610702565b600785905560088490556009839055604080516080808201835263ffffffff8a8116808452908a16602080850182905264ffffffffff8881168688018190529088166060968701819052600580547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000168617640100000000909502949094177fffffffffffffffffffffffffffff00000000000000000000ffffffffffffffff166801000000000000000083027fffffffffffffffffffffffffffff0000000000ffffffffffffffffffffffffff16176d010000000000000000000000000082021790935586519384529083018b90529482018990529281018790529081019290925260a08201527fe5c285d336cb17bb08823b603864963ca7aedc5a4d3fea30d299112cb47ddd5a9060c00160405180910390a150505050505050565b60006107ce8263ffffffff163a610837610ce8565b60008115610ac857600554610ac1908590640100000000810463ffffffff16906d0100000000000000000000000000900464ffffffffff1686610ebc565b905061081b565b600554610af9908590640100000000810463ffffffff169068010000000000000000900464ffffffffff1686610f37565b949350505050565b610b09610b9c565b6003805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f23b99d3a969380aa9df8e7afd6d3dbff42d352acaae63d51ad0466d62a1a917d91016105f2565b610b90610b9c565b610b9981611068565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610c1d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610702565b565b600080610c318463ffffffff1661115e565b8563ffffffff168763ffffffff168a610c4a9190611984565b610c549190611984565b610c5e9089611b26565b610c689190611984565b90508060006064610c79868261199c565b610c869060ff1684611b26565b610c9091906119c1565b600554600454919250600091610ccf9168010000000000000000900464ffffffffff169073ffffffffffffffffffffffffffffffffffffffff16611220565b610cd99083611984565b9b9a5050505050505050505050565b600554600254604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009363ffffffff1692831515928592839273ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b158015610d6557600080fd5b505afa158015610d79573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d9d9190611917565b509450909250849150508015610dc15750610db88242611b63565b8463ffffffff16105b15610af95750600754949350505050565b600080610de48463ffffffff1661115e565b8563ffffffff168763ffffffff168b610dfd9190611984565b610e079190611984565b610e11908a611b26565b610e1b9190611984565b9050600087610e3283670de0b6b3a7640000611b26565b610e3c91906119c1565b905060006064610e4c868261199c565b610e599060ff1684611b26565b610e6391906119c1565b600554600354919250600091610ea29168010000000000000000900464ffffffffff169073ffffffffffffffffffffffffffffffffffffffff16611220565b610eac9083611984565b9c9b505050505050505050505050565b600080610ec7611374565b905060005a610ed68888611984565b610ee09190611b63565b610eea9085611b26565b600454909150600090610f1490879073ffffffffffffffffffffffffffffffffffffffff16611220565b905082610f218284611984565b610f2b9190611984565b98975050505050505050565b600080610f42610ce8565b905060008113610f81576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610702565b6000610f8b611374565b9050600082825a610f9c8b8b611984565b610fa69190611b63565b610fb09088611b26565b610fba9190611984565b610fcc90670de0b6b3a7640000611b26565b610fd691906119c1565b60025490915060009061100090889073ffffffffffffffffffffffffffffffffffffffff16611220565b9050611018816b033b2e3c9fd0803ce8000000611b63565b821115611051576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61105b8183611984565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156110e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610702565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b1811480611173575062066eed81145b15611217576000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c06040518083038186803b1580156111c157600080fd5b505afa1580156111d5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111f991906117c1565b5050505091505083608c61120d9190611984565b610af99082611b26565b50600092915050565b600064ffffffffff8316611236575060006107ce565b6000806112428461141b565b9092509050600082136112a0576040517fc3388fe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260248101839052604401610702565b600860ff821610156112fb5760006112b9826008611b7a565b90506112c681600a611a5d565b6112d09084611b26565b6112e964ffffffffff8816670de0b6b3a7640000611b26565b6112f391906119c1565b93505061136c565b600860ff82161115611345576000611314600883611b7a565b90508261132282600a611a5d565b61133b64ffffffffff8916670de0b6b3a7640000611b26565b6112e99190611b26565b8161135f64ffffffffff8716670de0b6b3a7640000611b26565b61136991906119c1565b92505b505092915050565b60004661a4b1811480611389575062066eed81145b1561141357606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156113d557600080fd5b505afa1580156113e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061140d919061176a565b91505090565b600091505090565b600354600090819073ffffffffffffffffffffffffffffffffffffffff848116911614801590611466575060045473ffffffffffffffffffffffffffffffffffffffff848116911614155b156114cd57600354600480546040517f76266ef600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff808816938201939093529282166024840152166044820152606401610702565b600554604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905163ffffffff909216918215159160009173ffffffffffffffffffffffffffffffffffffffff88169163feaf968c9160048083019260a0929190829003018186803b15801561154657600080fd5b505afa15801561155a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061157e9190611917565b509197509092508391505080156115a3575061159a8142611b63565b8363ffffffff16105b156116645760045473ffffffffffffffffffffffffffffffffffffffff878116911614156115d5576008549450611664565b60035473ffffffffffffffffffffffffffffffffffffffff87811691161415611602576009549450611664565b600354600480546040517f76266ef600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff808b16938201939093529282166024840152166044820152606401610702565b8573ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b1580156116aa57600080fd5b505afa1580156116be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116e29190611967565b9350505050915091565b803563ffffffff8116811461170057600080fd5b919050565b803564ffffffffff8116811461170057600080fd5b805169ffffffffffffffffffff8116811461170057600080fd5b60006020828403121561174657600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461081b57600080fd5b60006020828403121561177c57600080fd5b5051919050565b60008060006060848603121561179857600080fd5b8335925060208401359150604084013580151581146117b657600080fd5b809150509250925092565b60008060008060008060c087890312156117da57600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b60006020828403121561181d57600080fd5b61081b826116ec565b6000806040838503121561183957600080fd5b611842836116ec565b946020939093013593505050565b600080600080600080600060e0888a03121561186b57600080fd5b611874886116ec565b9650611882602089016116ec565b95506040880135945060608801359350608088013592506118a560a08901611705565b91506118b360c08901611705565b905092959891949750929550565b600080600080608085870312156118d757600080fd5b6118e0856116ec565b93506118ee602086016116ec565b925060408501356118fe81611bcc565b915061190c606086016116ec565b905092959194509250565b600080600080600060a0868803121561192f57600080fd5b6119388661171a565b945060208601519350604086015192506060860151915061195b6080870161171a565b90509295509295909350565b60006020828403121561197957600080fd5b815161081b81611bcc565b6000821982111561199757611997611b9d565b500190565b600060ff821660ff84168060ff038211156119b9576119b9611b9d565b019392505050565b6000826119f7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600181815b80851115611a5557817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611a3b57611a3b611b9d565b80851615611a4857918102915b93841c9390800290611a01565b509250929050565b600061081b60ff841683600082611a76575060016107ce565b81611a83575060006107ce565b8160018114611a995760028114611aa357611abf565b60019150506107ce565b60ff841115611ab457611ab4611b9d565b50506001821b6107ce565b5060208310610133831016604e8410600b8410161715611ae2575081810a6107ce565b611aec83836119fc565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611b1e57611b1e611b9d565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611b5e57611b5e611b9d565b500290565b600082821015611b7557611b75611b9d565b500390565b600060ff821660ff841680821015611b9457611b94611b9d565b90039392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff81168114610b9957600080fdfea164736f6c6343000806000a",
}

var VRFV2PlusPriceRegistryABI = VRFV2PlusPriceRegistryMetaData.ABI

var VRFV2PlusPriceRegistryBin = VRFV2PlusPriceRegistryMetaData.Bin

func DeployVRFV2PlusPriceRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, linkEthFeed common.Address, linkUSDFeed common.Address, ethUSDFeed common.Address) (common.Address, *types.Transaction, *VRFV2PlusPriceRegistry, error) {
	parsed, err := VRFV2PlusPriceRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusPriceRegistryBin), backend, linkEthFeed, linkUSDFeed, ethUSDFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusPriceRegistry{VRFV2PlusPriceRegistryCaller: VRFV2PlusPriceRegistryCaller{contract: contract}, VRFV2PlusPriceRegistryTransactor: VRFV2PlusPriceRegistryTransactor{contract: contract}, VRFV2PlusPriceRegistryFilterer: VRFV2PlusPriceRegistryFilterer{contract: contract}}, nil
}

type VRFV2PlusPriceRegistry struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusPriceRegistryCaller
	VRFV2PlusPriceRegistryTransactor
	VRFV2PlusPriceRegistryFilterer
}

type VRFV2PlusPriceRegistryCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusPriceRegistryTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusPriceRegistryFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusPriceRegistrySession struct {
	Contract     *VRFV2PlusPriceRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusPriceRegistryCallerSession struct {
	Contract *VRFV2PlusPriceRegistryCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusPriceRegistryTransactorSession struct {
	Contract     *VRFV2PlusPriceRegistryTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusPriceRegistryRaw struct {
	Contract *VRFV2PlusPriceRegistry
}

type VRFV2PlusPriceRegistryCallerRaw struct {
	Contract *VRFV2PlusPriceRegistryCaller
}

type VRFV2PlusPriceRegistryTransactorRaw struct {
	Contract *VRFV2PlusPriceRegistryTransactor
}

func NewVRFV2PlusPriceRegistry(address common.Address, backend bind.ContractBackend) (*VRFV2PlusPriceRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusPriceRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusPriceRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistry{address: address, abi: abi, VRFV2PlusPriceRegistryCaller: VRFV2PlusPriceRegistryCaller{contract: contract}, VRFV2PlusPriceRegistryTransactor: VRFV2PlusPriceRegistryTransactor{contract: contract}, VRFV2PlusPriceRegistryFilterer: VRFV2PlusPriceRegistryFilterer{contract: contract}}, nil
}

func NewVRFV2PlusPriceRegistryCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusPriceRegistryCaller, error) {
	contract, err := bindVRFV2PlusPriceRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryCaller{contract: contract}, nil
}

func NewVRFV2PlusPriceRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusPriceRegistryTransactor, error) {
	contract, err := bindVRFV2PlusPriceRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryTransactor{contract: contract}, nil
}

func NewVRFV2PlusPriceRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusPriceRegistryFilterer, error) {
	contract, err := bindVRFV2PlusPriceRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryFilterer{contract: contract}, nil
}

func bindVRFV2PlusPriceRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusPriceRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusPriceRegistry.Contract.VRFV2PlusPriceRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.VRFV2PlusPriceRegistryTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.VRFV2PlusPriceRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusPriceRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) USDFEEDECIMALS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "USD_FEE_DECIMALS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) USDFEEDECIMALS() (uint8, error) {
	return _VRFV2PlusPriceRegistry.Contract.USDFEEDECIMALS(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) USDFEEDECIMALS() (uint8, error) {
	return _VRFV2PlusPriceRegistry.Contract.USDFEEDECIMALS(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) CalculatePaymentAmount(opts *bind.CallOpts, startGas *big.Int, weiPerUnitGas *big.Int, nativePayment bool) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "calculatePaymentAmount", startGas, weiPerUnitGas, nativePayment)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) CalculatePaymentAmount(startGas *big.Int, weiPerUnitGas *big.Int, nativePayment bool) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.CalculatePaymentAmount(&_VRFV2PlusPriceRegistry.CallOpts, startGas, weiPerUnitGas, nativePayment)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) CalculatePaymentAmount(startGas *big.Int, weiPerUnitGas *big.Int, nativePayment bool) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.CalculatePaymentAmount(&_VRFV2PlusPriceRegistry.CallOpts, startGas, weiPerUnitGas, nativePayment)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) CalculateRequestPriceNativeWrapper(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "calculateRequestPriceNativeWrapper", _callbackGasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) CalculateRequestPriceNativeWrapper(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.CalculateRequestPriceNativeWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) CalculateRequestPriceNativeWrapper(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.CalculateRequestPriceNativeWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) CalculateRequestPriceWrapper(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "calculateRequestPriceWrapper", _callbackGasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) CalculateRequestPriceWrapper(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.CalculateRequestPriceWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) CalculateRequestPriceWrapper(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.CalculateRequestPriceWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) EstimateRequestPriceNativeWrapper(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "estimateRequestPriceNativeWrapper", _callbackGasLimit, _requestGasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) EstimateRequestPriceNativeWrapper(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.EstimateRequestPriceNativeWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) EstimateRequestPriceNativeWrapper(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.EstimateRequestPriceNativeWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) EstimateRequestPriceWrapper(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "estimateRequestPriceWrapper", _callbackGasLimit, _requestGasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) EstimateRequestPriceWrapper(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.EstimateRequestPriceWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) EstimateRequestPriceWrapper(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.EstimateRequestPriceWrapper(&_VRFV2PlusPriceRegistry.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) Owner() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.Owner(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.Owner(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.StalenessSeconds = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkUSD = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FulfillmentFlatFeeEthUSD = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SConfig() (SConfig,

	error) {
	return _VRFV2PlusPriceRegistry.Contract.SConfig(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SConfig() (SConfig,

	error) {
	return _VRFV2PlusPriceRegistry.Contract.SConfig(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SEthUSDFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_ethUSDFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SEthUSDFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SEthUSDFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SEthUSDFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SEthUSDFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SFallbackUSDPerUnitEth(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_fallbackUSDPerUnitEth")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SFallbackUSDPerUnitEth() (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.SFallbackUSDPerUnitEth(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SFallbackUSDPerUnitEth() (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.SFallbackUSDPerUnitEth(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SFallbackUSDPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_fallbackUSDPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SFallbackUSDPerUnitLink() (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.SFallbackUSDPerUnitLink(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SFallbackUSDPerUnitLink() (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.SFallbackUSDPerUnitLink(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.SFallbackWeiPerUnitLink(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFV2PlusPriceRegistry.Contract.SFallbackWeiPerUnitLink(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SLinkETHFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_linkETHFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SLinkETHFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SLinkETHFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SLinkETHFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SLinkETHFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SLinkUSDFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_linkUSDFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SLinkUSDFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SLinkUSDFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SLinkUSDFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SLinkUSDFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SWrapperConfig(opts *bind.CallOpts) (SWrapperConfig,

	error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_wrapperConfig")

	outstruct := new(SWrapperConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.WrapperGasOverhead = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.CoordinatorGasOverhead = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.WrapperPremiumPercentage = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.FulfillmentTxSizeBytes = *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SWrapperConfig() (SWrapperConfig,

	error) {
	return _VRFV2PlusPriceRegistry.Contract.SWrapperConfig(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SWrapperConfig() (SWrapperConfig,

	error) {
	return _VRFV2PlusPriceRegistry.Contract.SWrapperConfig(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.AcceptOwnership(&_VRFV2PlusPriceRegistry.TransactOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.AcceptOwnership(&_VRFV2PlusPriceRegistry.TransactOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, fulfillmentFlatFeeLinkUSD *big.Int, fulfillmentFlatFeeEthUSD *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "setConfig", stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fallbackUSDPerUnitEth, fallbackUSDPerUnitLink, fulfillmentFlatFeeLinkUSD, fulfillmentFlatFeeEthUSD)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SetConfig(stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, fulfillmentFlatFeeLinkUSD *big.Int, fulfillmentFlatFeeEthUSD *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetConfig(&_VRFV2PlusPriceRegistry.TransactOpts, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fallbackUSDPerUnitEth, fallbackUSDPerUnitLink, fulfillmentFlatFeeLinkUSD, fulfillmentFlatFeeEthUSD)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) SetConfig(stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, fulfillmentFlatFeeLinkUSD *big.Int, fulfillmentFlatFeeEthUSD *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetConfig(&_VRFV2PlusPriceRegistry.TransactOpts, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fallbackUSDPerUnitEth, fallbackUSDPerUnitLink, fulfillmentFlatFeeLinkUSD, fulfillmentFlatFeeEthUSD)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) SetETHUSDFeed(opts *bind.TransactOpts, ethUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "setETHUSDFeed", ethUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SetETHUSDFeed(ethUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetETHUSDFeed(&_VRFV2PlusPriceRegistry.TransactOpts, ethUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) SetETHUSDFeed(ethUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetETHUSDFeed(&_VRFV2PlusPriceRegistry.TransactOpts, ethUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) SetLINKETHFeed(opts *bind.TransactOpts, linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "setLINKETHFeed", linkEthFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SetLINKETHFeed(linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetLINKETHFeed(&_VRFV2PlusPriceRegistry.TransactOpts, linkEthFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) SetLINKETHFeed(linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetLINKETHFeed(&_VRFV2PlusPriceRegistry.TransactOpts, linkEthFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) SetLINKUSDFeed(opts *bind.TransactOpts, linkUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "setLINKUSDFeed", linkUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SetLINKUSDFeed(linkUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetLINKUSDFeed(&_VRFV2PlusPriceRegistry.TransactOpts, linkUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) SetLINKUSDFeed(linkUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetLINKUSDFeed(&_VRFV2PlusPriceRegistry.TransactOpts, linkUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) SetWrapperConfig(opts *bind.TransactOpts, wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, fulfillmentTxSizeBytes uint32) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "setWrapperConfig", wrapperGasOverhead, coordinatorGasOverhead, wrapperPremiumPercentage, fulfillmentTxSizeBytes)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SetWrapperConfig(wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, fulfillmentTxSizeBytes uint32) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetWrapperConfig(&_VRFV2PlusPriceRegistry.TransactOpts, wrapperGasOverhead, coordinatorGasOverhead, wrapperPremiumPercentage, fulfillmentTxSizeBytes)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) SetWrapperConfig(wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, fulfillmentTxSizeBytes uint32) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetWrapperConfig(&_VRFV2PlusPriceRegistry.TransactOpts, wrapperGasOverhead, coordinatorGasOverhead, wrapperPremiumPercentage, fulfillmentTxSizeBytes)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.TransferOwnership(&_VRFV2PlusPriceRegistry.TransactOpts, to)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.TransferOwnership(&_VRFV2PlusPriceRegistry.TransactOpts, to)
}

type VRFV2PlusPriceRegistryConfigSetIterator struct {
	Event *VRFV2PlusPriceRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryConfigSet)
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
		it.Event = new(VRFV2PlusPriceRegistryConfigSet)
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

func (it *VRFV2PlusPriceRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryConfigSet struct {
	StalenessSeconds          uint32
	FallbackWeiPerUnitLink    *big.Int
	FallbackUSDPerUnitEth     *big.Int
	FallbackUSDPerUnitLink    *big.Int
	FulfillmentFlatFeeLinkUSD *big.Int
	FulfillmentFlatFeeEthUSD  *big.Int
	Raw                       types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryConfigSetIterator, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryConfigSetIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryConfigSet)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseConfigSet(log types.Log) (*VRFV2PlusPriceRegistryConfigSet, error) {
	event := new(VRFV2PlusPriceRegistryConfigSet)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusPriceRegistryEthUSDFeedSetIterator struct {
	Event *VRFV2PlusPriceRegistryEthUSDFeedSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryEthUSDFeedSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryEthUSDFeedSet)
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
		it.Event = new(VRFV2PlusPriceRegistryEthUSDFeedSet)
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

func (it *VRFV2PlusPriceRegistryEthUSDFeedSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryEthUSDFeedSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryEthUSDFeedSet struct {
	OldFeed common.Address
	NewFeed common.Address
	Raw     types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterEthUSDFeedSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryEthUSDFeedSetIterator, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "EthUSDFeedSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryEthUSDFeedSetIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "EthUSDFeedSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchEthUSDFeedSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryEthUSDFeedSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "EthUSDFeedSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryEthUSDFeedSet)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "EthUSDFeedSet", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseEthUSDFeedSet(log types.Log) (*VRFV2PlusPriceRegistryEthUSDFeedSet, error) {
	event := new(VRFV2PlusPriceRegistryEthUSDFeedSet)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "EthUSDFeedSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusPriceRegistryLinkEthFeedSetIterator struct {
	Event *VRFV2PlusPriceRegistryLinkEthFeedSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryLinkEthFeedSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryLinkEthFeedSet)
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
		it.Event = new(VRFV2PlusPriceRegistryLinkEthFeedSet)
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

func (it *VRFV2PlusPriceRegistryLinkEthFeedSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryLinkEthFeedSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryLinkEthFeedSet struct {
	OldFeed common.Address
	NewFeed common.Address
	Raw     types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterLinkEthFeedSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryLinkEthFeedSetIterator, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "LinkEthFeedSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryLinkEthFeedSetIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "LinkEthFeedSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchLinkEthFeedSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryLinkEthFeedSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "LinkEthFeedSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryLinkEthFeedSet)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "LinkEthFeedSet", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseLinkEthFeedSet(log types.Log) (*VRFV2PlusPriceRegistryLinkEthFeedSet, error) {
	event := new(VRFV2PlusPriceRegistryLinkEthFeedSet)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "LinkEthFeedSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusPriceRegistryLinkUSDFeedSetIterator struct {
	Event *VRFV2PlusPriceRegistryLinkUSDFeedSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryLinkUSDFeedSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryLinkUSDFeedSet)
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
		it.Event = new(VRFV2PlusPriceRegistryLinkUSDFeedSet)
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

func (it *VRFV2PlusPriceRegistryLinkUSDFeedSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryLinkUSDFeedSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryLinkUSDFeedSet struct {
	OldFeed common.Address
	NewFeed common.Address
	Raw     types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterLinkUSDFeedSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryLinkUSDFeedSetIterator, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "LinkUSDFeedSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryLinkUSDFeedSetIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "LinkUSDFeedSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchLinkUSDFeedSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryLinkUSDFeedSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "LinkUSDFeedSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryLinkUSDFeedSet)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "LinkUSDFeedSet", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseLinkUSDFeedSet(log types.Log) (*VRFV2PlusPriceRegistryLinkUSDFeedSet, error) {
	event := new(VRFV2PlusPriceRegistryLinkUSDFeedSet)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "LinkUSDFeedSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusPriceRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusPriceRegistryOwnershipTransferRequested)
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

func (it *VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryOwnershipTransferRequested)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusPriceRegistryOwnershipTransferRequested, error) {
	event := new(VRFV2PlusPriceRegistryOwnershipTransferRequested)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusPriceRegistryOwnershipTransferredIterator struct {
	Event *VRFV2PlusPriceRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryOwnershipTransferred)
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
		it.Event = new(VRFV2PlusPriceRegistryOwnershipTransferred)
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

func (it *VRFV2PlusPriceRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusPriceRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryOwnershipTransferredIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryOwnershipTransferred)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusPriceRegistryOwnershipTransferred, error) {
	event := new(VRFV2PlusPriceRegistryOwnershipTransferred)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusPriceRegistryWrapperConfigSetIterator struct {
	Event *VRFV2PlusPriceRegistryWrapperConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusPriceRegistryWrapperConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusPriceRegistryWrapperConfigSet)
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
		it.Event = new(VRFV2PlusPriceRegistryWrapperConfigSet)
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

func (it *VRFV2PlusPriceRegistryWrapperConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusPriceRegistryWrapperConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusPriceRegistryWrapperConfigSet struct {
	WrapperGasOverhead       uint32
	CoordinatorGasOverhead   uint32
	WrapperPremiumPercentage uint8
	FulfillmentTxSizeBytes   uint32
	Raw                      types.Log
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) FilterWrapperConfigSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryWrapperConfigSetIterator, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.FilterLogs(opts, "WrapperConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusPriceRegistryWrapperConfigSetIterator{contract: _VRFV2PlusPriceRegistry.contract, event: "WrapperConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) WatchWrapperConfigSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryWrapperConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusPriceRegistry.contract.WatchLogs(opts, "WrapperConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusPriceRegistryWrapperConfigSet)
				if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "WrapperConfigSet", log); err != nil {
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryFilterer) ParseWrapperConfigSet(log types.Log) (*VRFV2PlusPriceRegistryWrapperConfigSet, error) {
	event := new(VRFV2PlusPriceRegistryWrapperConfigSet)
	if err := _VRFV2PlusPriceRegistry.contract.UnpackLog(event, "WrapperConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SConfig struct {
	StalenessSeconds           uint32
	GasAfterPaymentCalculation uint32
	FulfillmentFlatFeeLinkUSD  *big.Int
	FulfillmentFlatFeeEthUSD   *big.Int
}
type SWrapperConfig struct {
	WrapperGasOverhead       uint32
	CoordinatorGasOverhead   uint32
	WrapperPremiumPercentage uint8
	FulfillmentTxSizeBytes   uint32
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusPriceRegistry.abi.Events["ConfigSet"].ID:
		return _VRFV2PlusPriceRegistry.ParseConfigSet(log)
	case _VRFV2PlusPriceRegistry.abi.Events["EthUSDFeedSet"].ID:
		return _VRFV2PlusPriceRegistry.ParseEthUSDFeedSet(log)
	case _VRFV2PlusPriceRegistry.abi.Events["LinkEthFeedSet"].ID:
		return _VRFV2PlusPriceRegistry.ParseLinkEthFeedSet(log)
	case _VRFV2PlusPriceRegistry.abi.Events["LinkUSDFeedSet"].ID:
		return _VRFV2PlusPriceRegistry.ParseLinkUSDFeedSet(log)
	case _VRFV2PlusPriceRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusPriceRegistry.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusPriceRegistry.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusPriceRegistry.ParseOwnershipTransferred(log)
	case _VRFV2PlusPriceRegistry.abi.Events["WrapperConfigSet"].ID:
		return _VRFV2PlusPriceRegistry.ParseWrapperConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusPriceRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0xe5c285d336cb17bb08823b603864963ca7aedc5a4d3fea30d299112cb47ddd5a")
}

func (VRFV2PlusPriceRegistryEthUSDFeedSet) Topic() common.Hash {
	return common.HexToHash("0xe6a1e056cb2ec82c5f49294ff925bd5a0ab6a8ccbe8fdfdf7d9a333d9c12c507")
}

func (VRFV2PlusPriceRegistryLinkEthFeedSet) Topic() common.Hash {
	return common.HexToHash("0x15f61b91e528d42be960613d5606dbf13df3ef988e6a097b8543c9a58b2b7fd8")
}

func (VRFV2PlusPriceRegistryLinkUSDFeedSet) Topic() common.Hash {
	return common.HexToHash("0x23b99d3a969380aa9df8e7afd6d3dbff42d352acaae63d51ad0466d62a1a917d")
}

func (VRFV2PlusPriceRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusPriceRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2PlusPriceRegistryWrapperConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb76e5cdd0dc2bc17df6a1911db92060a11b138d8002f77267013944123a4475d")
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistry) Address() common.Address {
	return _VRFV2PlusPriceRegistry.address
}

type VRFV2PlusPriceRegistryInterface interface {
	USDFEEDECIMALS(opts *bind.CallOpts) (uint8, error)

	CalculatePaymentAmount(opts *bind.CallOpts, startGas *big.Int, weiPerUnitGas *big.Int, nativePayment bool) (*big.Int, error)

	CalculateRequestPriceNativeWrapper(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error)

	CalculateRequestPriceWrapper(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error)

	EstimateRequestPriceNativeWrapper(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error)

	EstimateRequestPriceWrapper(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SEthUSDFeed(opts *bind.CallOpts) (common.Address, error)

	SFallbackUSDPerUnitEth(opts *bind.CallOpts) (*big.Int, error)

	SFallbackUSDPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SLinkETHFeed(opts *bind.CallOpts) (common.Address, error)

	SLinkUSDFeed(opts *bind.CallOpts) (common.Address, error)

	SWrapperConfig(opts *bind.CallOpts) (SWrapperConfig,

		error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, fulfillmentFlatFeeLinkUSD *big.Int, fulfillmentFlatFeeEthUSD *big.Int) (*types.Transaction, error)

	SetETHUSDFeed(opts *bind.TransactOpts, ethUsdFeed common.Address) (*types.Transaction, error)

	SetLINKETHFeed(opts *bind.TransactOpts, linkEthFeed common.Address) (*types.Transaction, error)

	SetLINKUSDFeed(opts *bind.TransactOpts, linkUsdFeed common.Address) (*types.Transaction, error)

	SetWrapperConfig(opts *bind.TransactOpts, wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, fulfillmentTxSizeBytes uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFV2PlusPriceRegistryConfigSet, error)

	FilterEthUSDFeedSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryEthUSDFeedSetIterator, error)

	WatchEthUSDFeedSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryEthUSDFeedSet) (event.Subscription, error)

	ParseEthUSDFeedSet(log types.Log) (*VRFV2PlusPriceRegistryEthUSDFeedSet, error)

	FilterLinkEthFeedSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryLinkEthFeedSetIterator, error)

	WatchLinkEthFeedSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryLinkEthFeedSet) (event.Subscription, error)

	ParseLinkEthFeedSet(log types.Log) (*VRFV2PlusPriceRegistryLinkEthFeedSet, error)

	FilterLinkUSDFeedSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryLinkUSDFeedSetIterator, error)

	WatchLinkUSDFeedSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryLinkUSDFeedSet) (event.Subscription, error)

	ParseLinkUSDFeedSet(log types.Log) (*VRFV2PlusPriceRegistryLinkUSDFeedSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusPriceRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusPriceRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusPriceRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusPriceRegistryOwnershipTransferred, error)

	FilterWrapperConfigSet(opts *bind.FilterOpts) (*VRFV2PlusPriceRegistryWrapperConfigSetIterator, error)

	WatchWrapperConfigSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusPriceRegistryWrapperConfigSet) (event.Subscription, error)

	ParseWrapperConfigSet(log types.Log) (*VRFV2PlusPriceRegistryWrapperConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
