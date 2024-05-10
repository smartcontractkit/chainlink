// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package price_registry

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

type IPriceRegistryTokenPriceFeedConfig struct {
	DataFeedAddress common.Address
	TokenDecimals   uint8
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalTimestampedPackedUint224 struct {
	Value     *big.Int
	Timestamp uint32
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type PriceRegistryTokenPriceFeedUpdate struct {
	SourceToken common.Address
	FeedConfig  IPriceRegistryTokenPriceFeedConfig
}

var PriceRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"priceUpdaters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokens\",\"type\":\"address[]\"},{\"internalType\":\"uint32\",\"name\":\"stalenessThreshold\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeeds\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chain\",\"type\":\"uint64\"}],\"name\":\"ChainNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DataFeedValueOutOfUint224Range\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStalenessThreshold\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpdaterOrOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timePassed\",\"type\":\"uint256\"}],\"name\":\"StaleGasPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timePassed\",\"type\":\"uint256\"}],\"name\":\"StaleTokenPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenNotSupported\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"indexed\":false,\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"priceFeedConfig\",\"type\":\"tuple\"}],\"name\":\"PriceFeedPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"priceUpdater\",\"type\":\"address\"}],\"name\":\"PriceUpdaterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"priceUpdater\",\"type\":\"address\"}],\"name\":\"PriceUpdaterSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChain\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerUnitGasUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"feeTokensToAdd\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokensToRemove\",\"type\":\"address[]\"}],\"name\":\"applyFeeTokensUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"priceUpdatersToAdd\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"priceUpdatersToRemove\",\"type\":\"address[]\"}],\"name\":\"applyPriceUpdatersUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fromTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"toToken\",\"type\":\"address\"}],\"name\":\"convertTokenAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestinationChainGasPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPriceUpdaters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStalenessThreshold\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getTokenAndGasPrices\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"tokenPrice\",\"type\":\"uint224\"},{\"internalType\":\"uint224\",\"name\":\"gasPriceValue\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPriceFeedConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getTokenPrices\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getValidatedTokenPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"updatePrices\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeedUpdates\",\"type\":\"tuple[]\"}],\"name\":\"updateTokenPriceFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002b7038038062002b7083398101604081905262000034916200080a565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be816200013f565b5050604080516000815260208101909152620000dd91508590620001ea565b604080516000815260208101909152620000f990849062000332565b620001048162000475565b8163ffffffff166000036200012c57604051631151410960e11b815260040160405180910390fd5b5063ffffffff1660805250620009d99050565b336001600160a01b03821603620001995760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b82518110156200028b57620002298382815181106200021057620002106200098b565b602002602001015160056200054660201b90919060201c565b1562000282578281815181106200024457620002446200098b565b60200260200101516001600160a01b03167f34a02290b7920078c19f58e94b78c77eb9cc10195b20676e19bd3b82085893b860405160405180910390a25b600101620001ed565b5060005b81518110156200032d57620002cb828281518110620002b257620002b26200098b565b602002602001015160056200056660201b90919060201c565b156200032457818181518110620002e657620002e66200098b565b60200260200101516001600160a01b03167fff7dbb85c77ca68ca1f894d6498570e3d5095cd19466f07ee8d222b337e4068c60405160405180910390a25b6001016200028f565b505050565b60005b8251811015620003d357620003718382815181106200035857620003586200098b565b602002602001015160076200054660201b90919060201c565b15620003ca578281815181106200038c576200038c6200098b565b60200260200101516001600160a01b03167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba2360405160405180910390a25b60010162000335565b5060005b81518110156200032d5762000413828281518110620003fa57620003fa6200098b565b602002602001015160076200056660201b90919060201c565b156200046c578181815181106200042e576200042e6200098b565b60200260200101516001600160a01b03167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9160405160405180910390a25b600101620003d7565b60005b8151811015620005425760008282815181106200049957620004996200098b565b6020908102919091018101518051818301516001600160a01b0380831660008181526004875260409081902084518154868a018051929096166001600160a81b03199091168117600160a01b60ff9384160217909255825191825293519093169683019690965293955091939092917f08a5f7f5bb38a81d8e43aca13ecd76431dbf8816ae4699affff7b00b2fc1c464910160405180910390a250505080600101905062000478565b5050565b60006200055d836001600160a01b0384166200057d565b90505b92915050565b60006200055d836001600160a01b038416620005cf565b6000818152600183016020526040812054620005c65750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000560565b50600062000560565b60008181526001830160205260408120548015620006c8576000620005f6600183620009a1565b85549091506000906200060c90600190620009a1565b9050818114620006785760008660000182815481106200063057620006306200098b565b90600052602060002001549050808760000184815481106200065657620006566200098b565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200068c576200068c620009c3565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000560565b600091505062000560565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156200070e576200070e620006d3565b60405290565b604051601f8201601f191681016001600160401b03811182821017156200073f576200073f620006d3565b604052919050565b60006001600160401b03821115620007635762000763620006d3565b5060051b60200190565b80516001600160a01b03811681146200078557600080fd5b919050565b600082601f8301126200079c57600080fd5b81516020620007b5620007af8362000747565b62000714565b8083825260208201915060208460051b870101935086841115620007d857600080fd5b602086015b84811015620007ff57620007f1816200076d565b8352918301918301620007dd565b509695505050505050565b600080600080608085870312156200082157600080fd5b84516001600160401b03808211156200083957600080fd5b62000847888389016200078a565b95506020915081870151818111156200085f57600080fd5b6200086d89828a016200078a565b95505060408088015163ffffffff811681146200088957600080fd5b8095505060608089015183811115620008a157600080fd5b89019250601f83018a13620008b557600080fd5b8251620008c6620007af8262000747565b81815260609091028401850190858101908c831115620008e557600080fd5b948601945b828610156200097a57858d0384811215620009055760008081fd5b6200090f620006e9565b6200091a886200076d565b815286601f19830112156200092f5760008081fd5b62000939620006e9565b9150620009488989016200076d565b82528688015160ff811681146200095f5760008081fd5b828a01528089019190915282529483019490860190620008ea565b999c989b5096995050505050505050565b634e487b7160e01b600052603260045260246000fd5b818103818111156200056057634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805161215f62000a116000396000818161033401528181610c6b01528181610cd401528181610dd10152610e46015261215f6000f3fe608060405234801561001057600080fd5b50600436106101355760003560e01c80637afac322116100b2578063bfcd456611610081578063d02641a011610066578063d02641a014610428578063f2fde38b1461043b578063ffdb4b371461044e57600080fd5b8063bfcd45661461040b578063cdc73d511461042057600080fd5b80637afac322146102ef5780638da5cb5b14610302578063a6c94a731461032a578063bf78e03f1461035e57600080fd5b806345ac924d11610109578063514e8cff116100ee578063514e8cff1461023157806352877af0146102d457806379ba5097146102e757600080fd5b806345ac924d146101d15780634ab35b0b146101f157600080fd5b806241e5be1461013a578063181f5a77146101605780633937306f146101a9578063407e1086146101be575b600080fd5b61014d6101483660046117d2565b610496565b6040519081526020015b60405180910390f35b61019c6040518060400160405280601781526020017f5072696365526567697374727920312e362e302d64657600000000000000000081525081565b604051610157919061180e565b6101bc6101b736600461187b565b610502565b005b6101bc6101cc366004611990565b610816565b6101e46101df366004611aa9565b61082a565b6040516101579190611b1e565b6102046101ff366004611b99565b6108f5565b6040517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9091168152602001610157565b6102c761023f366004611bcc565b60408051808201909152600080825260208201525067ffffffffffffffff166000908152600260209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811683527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169082015290565b6040516101579190611be7565b6101bc6102e2366004611c93565b610900565b6101bc610916565b6101bc6102fd366004611c93565b610a18565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610157565b60405163ffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610157565b6103d761036c366004611b99565b6040805180820182526000808252602091820181905273ffffffffffffffffffffffffffffffffffffffff93841681526004825282902082518084019093525492831682527401000000000000000000000000000000000000000090920460ff169181019190915290565b60408051825173ffffffffffffffffffffffffffffffffffffffff16815260209283015160ff169281019290925201610157565b610413610a2a565b6040516101579190611cf7565b610413610a3b565b6102c7610436366004611b99565b610a47565b6101bc610449366004611b99565b610b86565b61046161045c366004611d51565b610b97565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff938416815292909116602083015201610157565b60006104a182610d22565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166104c885610d22565b6104f0907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1685611db3565b6104fa9190611dca565b949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906105325750610530600533610e82565b155b15610569576040517f46f0815400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006105758280611e05565b9050905060005b818110156106bf5760006105908480611e05565b838181106105a0576105a0611e6d565b9050604002018036038101906105b69190611ec8565b604080518082018252602080840180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff908116845263ffffffff42818116858701908152885173ffffffffffffffffffffffffffffffffffffffff9081166000908152600390975295889020965190519092167c010000000000000000000000000000000000000000000000000000000002919092161790935584519051935194955016927f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a926106ae9290917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b60405180910390a25060010161057c565b5060006106cf6020840184611e05565b9050905060005b818110156108105760006106ed6020860186611e05565b838181106106fd576106fd611e6d565b9050604002018036038101906107139190611f05565b604080518082018252602080840180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff908116845263ffffffff42818116858701908152885167ffffffffffffffff9081166000908152600290975295889020965190519092167c010000000000000000000000000000000000000000000000000000000002919092161790935584519051935194955016927fdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e926107ff9290917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b60405180910390a2506001016106d6565b50505050565b61081e610eb4565b61082781610f37565b50565b60608160008167ffffffffffffffff811115610848576108486118b6565b60405190808252806020026020018201604052801561088d57816020015b60408051808201909152600080825260208201528152602001906001900390816108665790505b50905060005b828110156108ea576108c58686838181106108b0576108b0611e6d565b90506020020160208101906104369190611b99565b8282815181106108d7576108d7611e6d565b6020908102919091010152600101610893565b509150505b92915050565b60006108ef82610d22565b610908610eb4565b6109128282611035565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461099c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a20610eb4565b6109128282611181565b6060610a3660056112c8565b905090565b6060610a3660076112c8565b604080518082019091526000808252602082015273ffffffffffffffffffffffffffffffffffffffff8281166000908152600460209081526040918290208251808401909352549283168083527401000000000000000000000000000000000000000090930460ff169082015290610b3a57505073ffffffffffffffffffffffffffffffffffffffff166000908152600360209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811683527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169082015290565b600080610b46836112d5565b604080518082019091527bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909216825263ffffffff16602082015295945050505050565b610b8e610eb4565b610827816114d2565b67ffffffffffffffff811660009081526002602090815260408083208151808301909252547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811682527c0100000000000000000000000000000000000000000000000000000000900463ffffffff1691810182905282918203610c4f576040517f2e59db3a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610993565b6000816020015163ffffffff1642610c679190611f28565b90507f000000000000000000000000000000000000000000000000000000000000000063ffffffff16811115610d08576040517ff08bcb3e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015263ffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260448101829052606401610993565b610d1186610d22565b9151919350909150505b9250929050565b600080610d2e83610a47565b9050806020015163ffffffff1660001480610d66575080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16155b15610db5576040517f06439c6b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610993565b6000816020015163ffffffff1642610dcd9190611f28565b90507f000000000000000000000000000000000000000000000000000000000000000063ffffffff16811115610e7a576040517fc65fdfca00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015263ffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260448101829052606401610993565b505192915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610f35576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610993565b565b60005b8151811015610912576000828281518110610f5757610f57611e6d565b60209081029190910181015180518183015173ffffffffffffffffffffffffffffffffffffffff80831660008181526004875260409081902084518154868a018051929096167fffffffffffffffffffffff00000000000000000000000000000000000000000090911681177401000000000000000000000000000000000000000060ff9384160217909255825191825293519093169683019690965293955091939092917f08a5f7f5bb38a81d8e43aca13ecd76431dbf8816ae4699affff7b00b2fc1c464910160405180910390a2505050806001019050610f3a565b60005b82518110156110d85761106e83828151811061105657611056611e6d565b602002602001015160056115c790919063ffffffff16565b156110d05782818151811061108557611085611e6d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167f34a02290b7920078c19f58e94b78c77eb9cc10195b20676e19bd3b82085893b860405160405180910390a25b600101611038565b5060005b815181101561117c576111128282815181106110fa576110fa611e6d565b602002602001015160056115e990919063ffffffff16565b156111745781818151811061112957611129611e6d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167fff7dbb85c77ca68ca1f894d6498570e3d5095cd19466f07ee8d222b337e4068c60405160405180910390a25b6001016110dc565b505050565b60005b8251811015611224576111ba8382815181106111a2576111a2611e6d565b602002602001015160076115c790919063ffffffff16565b1561121c578281815181106111d1576111d1611e6d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba2360405160405180910390a25b600101611184565b5060005b815181101561117c5761125e82828151811061124657611246611e6d565b602002602001015160076115e990919063ffffffff16565b156112c05781818151811061127557611275611e6d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9160405160405180910390a25b600101611228565b60606000610ead8361160b565b60008060008360000151905060008173ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561132e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113529190611f55565b5050509150506000811215611393576040517f10cb51d100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000819050600086602001518473ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156113ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061140e9190611fa5565b6114189190611fc2565b905060248160ff16111561144d57611431602482611fdb565b61143c90600a612114565b6114469083611dca565b9150611470565b611458816024611fdb565b61146390600a612114565b61146d9083611db3565b91505b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8211156114c6576040517f10cb51d100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50954295509350505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611551576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610993565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000610ead8373ffffffffffffffffffffffffffffffffffffffff8416611667565b6000610ead8373ffffffffffffffffffffffffffffffffffffffff84166116b6565b60608160000180548060200260200160405190810160405280929190818152602001828054801561165b57602002820191906000526020600020905b815481526020019060010190808311611647575b50505050509050919050565b60008181526001830160205260408120546116ae575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556108ef565b5060006108ef565b6000818152600183016020526040812054801561179f5760006116da600183611f28565b85549091506000906116ee90600190611f28565b905081811461175357600086600001828154811061170e5761170e611e6d565b906000526020600020015490508087600001848154811061173157611731611e6d565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061176457611764612123565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506108ef565b60009150506108ef565b803573ffffffffffffffffffffffffffffffffffffffff811681146117cd57600080fd5b919050565b6000806000606084860312156117e757600080fd5b6117f0846117a9565b925060208401359150611805604085016117a9565b90509250925092565b60006020808352835180602085015260005b8181101561183c57858101830151858201604001528201611820565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561188d57600080fd5b813567ffffffffffffffff8111156118a457600080fd5b820160408185031215610ead57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611908576119086118b6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611955576119556118b6565b604052919050565b600067ffffffffffffffff821115611977576119776118b6565b5060051b60200190565b60ff8116811461082757600080fd5b600060208083850312156119a357600080fd5b823567ffffffffffffffff8111156119ba57600080fd5b8301601f810185136119cb57600080fd5b80356119de6119d98261195d565b61190e565b818152606091820283018401918482019190888411156119fd57600080fd5b938501935b83851015611a9d5784890381811215611a1b5760008081fd5b611a236118e5565b611a2c876117a9565b81526040807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084011215611a605760008081fd5b611a686118e5565b9250611a758989016117a9565b8352870135611a8381611981565b828901528088019190915283529384019391850191611a02565b50979650505050505050565b60008060208385031215611abc57600080fd5b823567ffffffffffffffff80821115611ad457600080fd5b818501915085601f830112611ae857600080fd5b813581811115611af757600080fd5b8660208260051b8501011115611b0c57600080fd5b60209290920196919550909350505050565b602080825282518282018190526000919060409081850190868401855b82811015611b8c57611b7c84835180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16825260209081015163ffffffff16910152565b9284019290850190600101611b3b565b5091979650505050505050565b600060208284031215611bab57600080fd5b610ead826117a9565b803567ffffffffffffffff811681146117cd57600080fd5b600060208284031215611bde57600080fd5b610ead82611bb4565b81517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260208083015163ffffffff1690820152604081016108ef565b600082601f830112611c3357600080fd5b81356020611c436119d98361195d565b8083825260208201915060208460051b870101935086841115611c6557600080fd5b602086015b84811015611c8857611c7b816117a9565b8352918301918301611c6a565b509695505050505050565b60008060408385031215611ca657600080fd5b823567ffffffffffffffff80821115611cbe57600080fd5b611cca86838701611c22565b93506020850135915080821115611ce057600080fd5b50611ced85828601611c22565b9150509250929050565b6020808252825182820181905260009190848201906040850190845b81811015611d4557835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611d13565b50909695505050505050565b60008060408385031215611d6457600080fd5b611d6d836117a9565b9150611d7b60208401611bb4565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176108ef576108ef611d84565b600082611e00577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611e3a57600080fd5b83018035915067ffffffffffffffff821115611e5557600080fd5b6020019150600681901b3603821315610d1b57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146117cd57600080fd5b600060408284031215611eda57600080fd5b611ee26118e5565b611eeb836117a9565b8152611ef960208401611e9c565b60208201529392505050565b600060408284031215611f1757600080fd5b611f1f6118e5565b611eeb83611bb4565b818103818111156108ef576108ef611d84565b805169ffffffffffffffffffff811681146117cd57600080fd5b600080600080600060a08688031215611f6d57600080fd5b611f7686611f3b565b9450602086015193506040860151925060608601519150611f9960808701611f3b565b90509295509295909350565b600060208284031215611fb757600080fd5b8151610ead81611981565b60ff81811683821601908111156108ef576108ef611d84565b60ff82811682821603908111156108ef576108ef611d84565b600181815b8085111561204d57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561203357612033611d84565b8085161561204057918102915b93841c9390800290611ff9565b509250929050565b600082612064575060016108ef565b81612071575060006108ef565b81600181146120875760028114612091576120ad565b60019150506108ef565b60ff8411156120a2576120a2611d84565b50506001821b6108ef565b5060208310610133831016604e8410600b84101617156120d0575081810a6108ef565b6120da8383611ff4565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561210c5761210c611d84565b029392505050565b6000610ead60ff841683612055565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var PriceRegistryABI = PriceRegistryMetaData.ABI

var PriceRegistryBin = PriceRegistryMetaData.Bin

func DeployPriceRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, priceUpdaters []common.Address, feeTokens []common.Address, stalenessThreshold uint32, tokenPriceFeeds []PriceRegistryTokenPriceFeedUpdate) (common.Address, *types.Transaction, *PriceRegistry, error) {
	parsed, err := PriceRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PriceRegistryBin), backend, priceUpdaters, feeTokens, stalenessThreshold, tokenPriceFeeds)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PriceRegistry{address: address, abi: *parsed, PriceRegistryCaller: PriceRegistryCaller{contract: contract}, PriceRegistryTransactor: PriceRegistryTransactor{contract: contract}, PriceRegistryFilterer: PriceRegistryFilterer{contract: contract}}, nil
}

type PriceRegistry struct {
	address common.Address
	abi     abi.ABI
	PriceRegistryCaller
	PriceRegistryTransactor
	PriceRegistryFilterer
}

type PriceRegistryCaller struct {
	contract *bind.BoundContract
}

type PriceRegistryTransactor struct {
	contract *bind.BoundContract
}

type PriceRegistryFilterer struct {
	contract *bind.BoundContract
}

type PriceRegistrySession struct {
	Contract     *PriceRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PriceRegistryCallerSession struct {
	Contract *PriceRegistryCaller
	CallOpts bind.CallOpts
}

type PriceRegistryTransactorSession struct {
	Contract     *PriceRegistryTransactor
	TransactOpts bind.TransactOpts
}

type PriceRegistryRaw struct {
	Contract *PriceRegistry
}

type PriceRegistryCallerRaw struct {
	Contract *PriceRegistryCaller
}

type PriceRegistryTransactorRaw struct {
	Contract *PriceRegistryTransactor
}

func NewPriceRegistry(address common.Address, backend bind.ContractBackend) (*PriceRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(PriceRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPriceRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PriceRegistry{address: address, abi: abi, PriceRegistryCaller: PriceRegistryCaller{contract: contract}, PriceRegistryTransactor: PriceRegistryTransactor{contract: contract}, PriceRegistryFilterer: PriceRegistryFilterer{contract: contract}}, nil
}

func NewPriceRegistryCaller(address common.Address, caller bind.ContractCaller) (*PriceRegistryCaller, error) {
	contract, err := bindPriceRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryCaller{contract: contract}, nil
}

func NewPriceRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*PriceRegistryTransactor, error) {
	contract, err := bindPriceRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryTransactor{contract: contract}, nil
}

func NewPriceRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*PriceRegistryFilterer, error) {
	contract, err := bindPriceRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryFilterer{contract: contract}, nil
}

func bindPriceRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PriceRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_PriceRegistry *PriceRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceRegistry.Contract.PriceRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_PriceRegistry *PriceRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceRegistry.Contract.PriceRegistryTransactor.contract.Transfer(opts)
}

func (_PriceRegistry *PriceRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceRegistry.Contract.PriceRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_PriceRegistry *PriceRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_PriceRegistry *PriceRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceRegistry.Contract.contract.Transfer(opts)
}

func (_PriceRegistry *PriceRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_PriceRegistry *PriceRegistryCaller) ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "convertTokenAmount", fromToken, fromTokenAmount, toToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.ConvertTokenAmount(&_PriceRegistry.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_PriceRegistry *PriceRegistryCallerSession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.ConvertTokenAmount(&_PriceRegistry.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_PriceRegistry *PriceRegistryCaller) GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getDestinationChainGasPrice", destChainSelector)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetDestinationChainGasPrice(&_PriceRegistry.CallOpts, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetDestinationChainGasPrice(&_PriceRegistry.CallOpts, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCaller) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getFeeTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetFeeTokens() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetFeeTokens(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetFeeTokens() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetFeeTokens(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) GetPriceUpdaters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getPriceUpdaters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetPriceUpdaters() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetPriceUpdaters(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetPriceUpdaters() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetPriceUpdaters(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) GetStalenessThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getStalenessThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetStalenessThreshold() (*big.Int, error) {
	return _PriceRegistry.Contract.GetStalenessThreshold(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetStalenessThreshold() (*big.Int, error) {
	return _PriceRegistry.Contract.GetStalenessThreshold(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenAndGasPrices", token, destChainSelector)

	outstruct := new(GetTokenAndGasPrices)
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenPrice = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.GasPriceValue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _PriceRegistry.Contract.GetTokenAndGasPrices(&_PriceRegistry.CallOpts, token, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _PriceRegistry.Contract.GetTokenAndGasPrices(&_PriceRegistry.CallOpts, token, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenPrice", token)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (IPriceRegistryTokenPriceFeedConfig, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenPriceFeedConfig", token)

	if err != nil {
		return *new(IPriceRegistryTokenPriceFeedConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IPriceRegistryTokenPriceFeedConfig)).(*IPriceRegistryTokenPriceFeedConfig)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenPriceFeedConfig(token common.Address) (IPriceRegistryTokenPriceFeedConfig, error) {
	return _PriceRegistry.Contract.GetTokenPriceFeedConfig(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenPriceFeedConfig(token common.Address) (IPriceRegistryTokenPriceFeedConfig, error) {
	return _PriceRegistry.Contract.GetTokenPriceFeedConfig(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenPrices", tokens)

	if err != nil {
		return *new([]InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalTimestampedPackedUint224)).(*[]InternalTimestampedPackedUint224)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrices(&_PriceRegistry.CallOpts, tokens)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrices(&_PriceRegistry.CallOpts, tokens)
}

func (_PriceRegistry *PriceRegistryCaller) GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getValidatedTokenPrice", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.GetValidatedTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.GetValidatedTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) Owner() (common.Address, error) {
	return _PriceRegistry.Contract.Owner(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) Owner() (common.Address, error) {
	return _PriceRegistry.Contract.Owner(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) TypeAndVersion() (string, error) {
	return _PriceRegistry.Contract.TypeAndVersion(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) TypeAndVersion() (string, error) {
	return _PriceRegistry.Contract.TypeAndVersion(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_PriceRegistry *PriceRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _PriceRegistry.Contract.AcceptOwnership(&_PriceRegistry.TransactOpts)
}

func (_PriceRegistry *PriceRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PriceRegistry.Contract.AcceptOwnership(&_PriceRegistry.TransactOpts)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyFeeTokensUpdates", feeTokensToAdd, feeTokensToRemove)
}

func (_PriceRegistry *PriceRegistrySession) ApplyFeeTokensUpdates(feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyFeeTokensUpdates(&_PriceRegistry.TransactOpts, feeTokensToAdd, feeTokensToRemove)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyFeeTokensUpdates(feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyFeeTokensUpdates(&_PriceRegistry.TransactOpts, feeTokensToAdd, feeTokensToRemove)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyPriceUpdatersUpdates(opts *bind.TransactOpts, priceUpdatersToAdd []common.Address, priceUpdatersToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyPriceUpdatersUpdates", priceUpdatersToAdd, priceUpdatersToRemove)
}

func (_PriceRegistry *PriceRegistrySession) ApplyPriceUpdatersUpdates(priceUpdatersToAdd []common.Address, priceUpdatersToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyPriceUpdatersUpdates(&_PriceRegistry.TransactOpts, priceUpdatersToAdd, priceUpdatersToRemove)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyPriceUpdatersUpdates(priceUpdatersToAdd []common.Address, priceUpdatersToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyPriceUpdatersUpdates(&_PriceRegistry.TransactOpts, priceUpdatersToAdd, priceUpdatersToRemove)
}

func (_PriceRegistry *PriceRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_PriceRegistry *PriceRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.TransferOwnership(&_PriceRegistry.TransactOpts, to)
}

func (_PriceRegistry *PriceRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.TransferOwnership(&_PriceRegistry.TransactOpts, to)
}

func (_PriceRegistry *PriceRegistryTransactor) UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "updatePrices", priceUpdates)
}

func (_PriceRegistry *PriceRegistrySession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdatePrices(&_PriceRegistry.TransactOpts, priceUpdates)
}

func (_PriceRegistry *PriceRegistryTransactorSession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdatePrices(&_PriceRegistry.TransactOpts, priceUpdates)
}

func (_PriceRegistry *PriceRegistryTransactor) UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "updateTokenPriceFeeds", tokenPriceFeedUpdates)
}

func (_PriceRegistry *PriceRegistrySession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdateTokenPriceFeeds(&_PriceRegistry.TransactOpts, tokenPriceFeedUpdates)
}

func (_PriceRegistry *PriceRegistryTransactorSession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdateTokenPriceFeeds(&_PriceRegistry.TransactOpts, tokenPriceFeedUpdates)
}

type PriceRegistryFeeTokenAddedIterator struct {
	Event *PriceRegistryFeeTokenAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryFeeTokenAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryFeeTokenAdded)
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
		it.Event = new(PriceRegistryFeeTokenAdded)
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

func (it *PriceRegistryFeeTokenAddedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryFeeTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryFeeTokenAdded struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenAddedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryFeeTokenAddedIterator{contract: _PriceRegistry.contract, event: "FeeTokenAdded", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenAdded, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryFeeTokenAdded)
				if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseFeeTokenAdded(log types.Log) (*PriceRegistryFeeTokenAdded, error) {
	event := new(PriceRegistryFeeTokenAdded)
	if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryFeeTokenRemovedIterator struct {
	Event *PriceRegistryFeeTokenRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryFeeTokenRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryFeeTokenRemoved)
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
		it.Event = new(PriceRegistryFeeTokenRemoved)
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

func (it *PriceRegistryFeeTokenRemovedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryFeeTokenRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryFeeTokenRemoved struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenRemovedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryFeeTokenRemovedIterator{contract: _PriceRegistry.contract, event: "FeeTokenRemoved", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryFeeTokenRemoved)
				if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseFeeTokenRemoved(log types.Log) (*PriceRegistryFeeTokenRemoved, error) {
	event := new(PriceRegistryFeeTokenRemoved)
	if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryOwnershipTransferRequestedIterator struct {
	Event *PriceRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryOwnershipTransferRequested)
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
		it.Event = new(PriceRegistryOwnershipTransferRequested)
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

func (it *PriceRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryOwnershipTransferRequestedIterator{contract: _PriceRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryOwnershipTransferRequested)
				if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*PriceRegistryOwnershipTransferRequested, error) {
	event := new(PriceRegistryOwnershipTransferRequested)
	if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryOwnershipTransferredIterator struct {
	Event *PriceRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryOwnershipTransferred)
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
		it.Event = new(PriceRegistryOwnershipTransferred)
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

func (it *PriceRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryOwnershipTransferredIterator{contract: _PriceRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryOwnershipTransferred)
				if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*PriceRegistryOwnershipTransferred, error) {
	event := new(PriceRegistryOwnershipTransferred)
	if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryPriceFeedPerTokenUpdatedIterator struct {
	Event *PriceRegistryPriceFeedPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryPriceFeedPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryPriceFeedPerTokenUpdated)
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
		it.Event = new(PriceRegistryPriceFeedPerTokenUpdated)
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

func (it *PriceRegistryPriceFeedPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryPriceFeedPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryPriceFeedPerTokenUpdated struct {
	Token           common.Address
	PriceFeedConfig IPriceRegistryTokenPriceFeedConfig
	Raw             types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryPriceFeedPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryPriceFeedPerTokenUpdatedIterator{contract: _PriceRegistry.contract, event: "PriceFeedPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryPriceFeedPerTokenUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParsePriceFeedPerTokenUpdated(log types.Log) (*PriceRegistryPriceFeedPerTokenUpdated, error) {
	event := new(PriceRegistryPriceFeedPerTokenUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryPriceUpdaterRemovedIterator struct {
	Event *PriceRegistryPriceUpdaterRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryPriceUpdaterRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryPriceUpdaterRemoved)
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
		it.Event = new(PriceRegistryPriceUpdaterRemoved)
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

func (it *PriceRegistryPriceUpdaterRemovedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryPriceUpdaterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryPriceUpdaterRemoved struct {
	PriceUpdater common.Address
	Raw          types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterPriceUpdaterRemoved(opts *bind.FilterOpts, priceUpdater []common.Address) (*PriceRegistryPriceUpdaterRemovedIterator, error) {

	var priceUpdaterRule []interface{}
	for _, priceUpdaterItem := range priceUpdater {
		priceUpdaterRule = append(priceUpdaterRule, priceUpdaterItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "PriceUpdaterRemoved", priceUpdaterRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryPriceUpdaterRemovedIterator{contract: _PriceRegistry.contract, event: "PriceUpdaterRemoved", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchPriceUpdaterRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceUpdaterRemoved, priceUpdater []common.Address) (event.Subscription, error) {

	var priceUpdaterRule []interface{}
	for _, priceUpdaterItem := range priceUpdater {
		priceUpdaterRule = append(priceUpdaterRule, priceUpdaterItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "PriceUpdaterRemoved", priceUpdaterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryPriceUpdaterRemoved)
				if err := _PriceRegistry.contract.UnpackLog(event, "PriceUpdaterRemoved", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParsePriceUpdaterRemoved(log types.Log) (*PriceRegistryPriceUpdaterRemoved, error) {
	event := new(PriceRegistryPriceUpdaterRemoved)
	if err := _PriceRegistry.contract.UnpackLog(event, "PriceUpdaterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryPriceUpdaterSetIterator struct {
	Event *PriceRegistryPriceUpdaterSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryPriceUpdaterSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryPriceUpdaterSet)
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
		it.Event = new(PriceRegistryPriceUpdaterSet)
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

func (it *PriceRegistryPriceUpdaterSetIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryPriceUpdaterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryPriceUpdaterSet struct {
	PriceUpdater common.Address
	Raw          types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterPriceUpdaterSet(opts *bind.FilterOpts, priceUpdater []common.Address) (*PriceRegistryPriceUpdaterSetIterator, error) {

	var priceUpdaterRule []interface{}
	for _, priceUpdaterItem := range priceUpdater {
		priceUpdaterRule = append(priceUpdaterRule, priceUpdaterItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "PriceUpdaterSet", priceUpdaterRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryPriceUpdaterSetIterator{contract: _PriceRegistry.contract, event: "PriceUpdaterSet", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchPriceUpdaterSet(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceUpdaterSet, priceUpdater []common.Address) (event.Subscription, error) {

	var priceUpdaterRule []interface{}
	for _, priceUpdaterItem := range priceUpdater {
		priceUpdaterRule = append(priceUpdaterRule, priceUpdaterItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "PriceUpdaterSet", priceUpdaterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryPriceUpdaterSet)
				if err := _PriceRegistry.contract.UnpackLog(event, "PriceUpdaterSet", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParsePriceUpdaterSet(log types.Log) (*PriceRegistryPriceUpdaterSet, error) {
	event := new(PriceRegistryPriceUpdaterSet)
	if err := _PriceRegistry.contract.UnpackLog(event, "PriceUpdaterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryUsdPerTokenUpdatedIterator struct {
	Event *PriceRegistryUsdPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryUsdPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryUsdPerTokenUpdated)
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
		it.Event = new(PriceRegistryUsdPerTokenUpdated)
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

func (it *PriceRegistryUsdPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryUsdPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryUsdPerTokenUpdated struct {
	Token     common.Address
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryUsdPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryUsdPerTokenUpdatedIterator{contract: _PriceRegistry.contract, event: "UsdPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryUsdPerTokenUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseUsdPerTokenUpdated(log types.Log) (*PriceRegistryUsdPerTokenUpdated, error) {
	event := new(PriceRegistryUsdPerTokenUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryUsdPerUnitGasUpdatedIterator struct {
	Event *PriceRegistryUsdPerUnitGasUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryUsdPerUnitGasUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryUsdPerUnitGasUpdated)
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
		it.Event = new(PriceRegistryUsdPerUnitGasUpdated)
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

func (it *PriceRegistryUsdPerUnitGasUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryUsdPerUnitGasUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryUsdPerUnitGasUpdated struct {
	DestChain uint64
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*PriceRegistryUsdPerUnitGasUpdatedIterator, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryUsdPerUnitGasUpdatedIterator{contract: _PriceRegistry.contract, event: "UsdPerUnitGasUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryUsdPerUnitGasUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseUsdPerUnitGasUpdated(log types.Log) (*PriceRegistryUsdPerUnitGasUpdated, error) {
	event := new(PriceRegistryUsdPerUnitGasUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetTokenAndGasPrices struct {
	TokenPrice    *big.Int
	GasPriceValue *big.Int
}

func (_PriceRegistry *PriceRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _PriceRegistry.abi.Events["FeeTokenAdded"].ID:
		return _PriceRegistry.ParseFeeTokenAdded(log)
	case _PriceRegistry.abi.Events["FeeTokenRemoved"].ID:
		return _PriceRegistry.ParseFeeTokenRemoved(log)
	case _PriceRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _PriceRegistry.ParseOwnershipTransferRequested(log)
	case _PriceRegistry.abi.Events["OwnershipTransferred"].ID:
		return _PriceRegistry.ParseOwnershipTransferred(log)
	case _PriceRegistry.abi.Events["PriceFeedPerTokenUpdated"].ID:
		return _PriceRegistry.ParsePriceFeedPerTokenUpdated(log)
	case _PriceRegistry.abi.Events["PriceUpdaterRemoved"].ID:
		return _PriceRegistry.ParsePriceUpdaterRemoved(log)
	case _PriceRegistry.abi.Events["PriceUpdaterSet"].ID:
		return _PriceRegistry.ParsePriceUpdaterSet(log)
	case _PriceRegistry.abi.Events["UsdPerTokenUpdated"].ID:
		return _PriceRegistry.ParseUsdPerTokenUpdated(log)
	case _PriceRegistry.abi.Events["UsdPerUnitGasUpdated"].ID:
		return _PriceRegistry.ParseUsdPerUnitGasUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PriceRegistryFeeTokenAdded) Topic() common.Hash {
	return common.HexToHash("0xdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23")
}

func (PriceRegistryFeeTokenRemoved) Topic() common.Hash {
	return common.HexToHash("0x1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91")
}

func (PriceRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PriceRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (PriceRegistryPriceFeedPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x08a5f7f5bb38a81d8e43aca13ecd76431dbf8816ae4699affff7b00b2fc1c464")
}

func (PriceRegistryPriceUpdaterRemoved) Topic() common.Hash {
	return common.HexToHash("0xff7dbb85c77ca68ca1f894d6498570e3d5095cd19466f07ee8d222b337e4068c")
}

func (PriceRegistryPriceUpdaterSet) Topic() common.Hash {
	return common.HexToHash("0x34a02290b7920078c19f58e94b78c77eb9cc10195b20676e19bd3b82085893b8")
}

func (PriceRegistryUsdPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a")
}

func (PriceRegistryUsdPerUnitGasUpdated) Topic() common.Hash {
	return common.HexToHash("0xdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e")
}

func (_PriceRegistry *PriceRegistry) Address() common.Address {
	return _PriceRegistry.address
}

type PriceRegistryInterface interface {
	ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error)

	GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error)

	GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetPriceUpdaters(opts *bind.CallOpts) ([]common.Address, error)

	GetStalenessThreshold(opts *bind.CallOpts) (*big.Int, error)

	GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

		error)

	GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error)

	GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (IPriceRegistryTokenPriceFeedConfig, error)

	GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error)

	GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error)

	ApplyPriceUpdatersUpdates(opts *bind.TransactOpts, priceUpdatersToAdd []common.Address, priceUpdatersToRemove []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error)

	UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error)

	FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenAddedIterator, error)

	WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenAdded, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenAdded(log types.Log) (*PriceRegistryFeeTokenAdded, error)

	FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenRemovedIterator, error)

	WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenRemoved(log types.Log) (*PriceRegistryFeeTokenRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PriceRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PriceRegistryOwnershipTransferred, error)

	FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryPriceFeedPerTokenUpdatedIterator, error)

	WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParsePriceFeedPerTokenUpdated(log types.Log) (*PriceRegistryPriceFeedPerTokenUpdated, error)

	FilterPriceUpdaterRemoved(opts *bind.FilterOpts, priceUpdater []common.Address) (*PriceRegistryPriceUpdaterRemovedIterator, error)

	WatchPriceUpdaterRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceUpdaterRemoved, priceUpdater []common.Address) (event.Subscription, error)

	ParsePriceUpdaterRemoved(log types.Log) (*PriceRegistryPriceUpdaterRemoved, error)

	FilterPriceUpdaterSet(opts *bind.FilterOpts, priceUpdater []common.Address) (*PriceRegistryPriceUpdaterSetIterator, error)

	WatchPriceUpdaterSet(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceUpdaterSet, priceUpdater []common.Address) (event.Subscription, error)

	ParsePriceUpdaterSet(log types.Log) (*PriceRegistryPriceUpdaterSet, error)

	FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryUsdPerTokenUpdatedIterator, error)

	WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParseUsdPerTokenUpdated(log types.Log) (*PriceRegistryUsdPerTokenUpdated, error)

	FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*PriceRegistryUsdPerUnitGasUpdatedIterator, error)

	WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error)

	ParseUsdPerUnitGasUpdated(log types.Log) (*PriceRegistryUsdPerUnitGasUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
