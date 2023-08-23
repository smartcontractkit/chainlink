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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ethUSDFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"ethUSD\",\"type\":\"int256\"}],\"name\":\"InvalidEthUSDPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkUSD\",\"type\":\"int256\"}],\"name\":\"InvalidLinkUSDPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feed\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"}],\"name\":\"InvalidUSDPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitEth\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeed\",\"type\":\"address\"}],\"name\":\"EthUSDFeedSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeed\",\"type\":\"address\"}],\"name\":\"LinkEthFeedSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeed\",\"type\":\"address\"}],\"name\":\"LinkUSDFeedSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"USD_FEE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitGas\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"calculatePaymentAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_ethUSDFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackUSDPerUnitEth\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackUSDPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkEthFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkUSDFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitEth\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethUsdFeed\",\"type\":\"address\"}],\"name\":\"setEthUsdFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"name\":\"setLINKETHFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkUsdFeed\",\"type\":\"address\"}],\"name\":\"setLINKUSDFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200160d3803806200160d8339810160408190526200003491620001cd565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000104565b5050600280546001600160a01b039586166001600160a01b0319918216179091556003805494861694821694909417909355506004805491909316911617905562000217565b6001600160a01b0381163314156200015f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c857600080fd5b919050565b600080600060608486031215620001e357600080fd5b620001ee84620001b0565b9250620001fe60208501620001b0565b91506200020e60408501620001b0565b90509250925092565b6113e680620002276000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806379ba509711610097578063de6a924811610066578063de6a92481461029b578063e6152d81146102ae578063e7ddbb8d146102de578063f2fde38b146102f157600080fd5b806379ba509714610252578063835c0dfc1461025a5780638da5cb5b14610263578063bb0697a51461028157600080fd5b806355a75475116100d357806355a75475146101f657806359392b6d1461020957806361d38666146102125780636af6890d1461023257600080fd5b8063043bd6ae146101055780630784e5d014610121578063088070f514610136578063180a4909146101b1575b600080fd5b61010e60065481565b6040519081526020015b60405180910390f35b61013461012f366004611045565b610304565b005b60055461017c9063ffffffff8082169164010000000081049091169064ffffffffff6801000000000000000082048116916d010000000000000000000000000090041684565b6040805163ffffffff958616815294909316602085015264ffffffffff91821692840192909252166060820152608001610118565b6004546101d19073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610118565b610134610204366004611045565b610393565b61010e60085481565b6002546101d19073ffffffffffffffffffffffffffffffffffffffff1681565b6003546101d19073ffffffffffffffffffffffffffffffffffffffff1681565b61013461041a565b61010e60075481565b60005473ffffffffffffffffffffffffffffffffffffffff166101d1565b610289600881565b60405160ff9091168152602001610118565b6101346102a93660046110d2565b61051c565b6102c16102bc366004611094565b610719565b6040516bffffffffffffffffffffffff9091168152602001610118565b6101346102ec366004611045565b610799565b6101346102ff366004611045565b610820565b61030c610834565b6002805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f15f61b91e528d42be960613d5606dbf13df3ef988e6a097b8543c9a58b2b7fd891015b60405180910390a15050565b61039b610834565b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527fe6a1e056cb2ec82c5f49294ff925bd5a0ab6a8ccbe8fdfdf7d9a333d9c12c5079101610387565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610524610834565b60008513610561576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101869052602401610497565b6000841361059e576040517f599d67e300000000000000000000000000000000000000000000000000000000815260048101859052602401610497565b600083136105db576040517f25b2499f00000000000000000000000000000000000000000000000000000000815260048101849052602401610497565b600685905560078490556008839055604080516080808201835263ffffffff8a8116808452908a16602080850182905264ffffffffff8881168688018190529088166060968701819052600580547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000168617640100000000909502949094177fffffffffffffffffffffffffffff00000000000000000000ffffffffffffffff166801000000000000000083027fffffffffffffffffffffffffffff0000000000ffffffffffffffffffffffffff16176d010000000000000000000000000082021790935586519384529083018b90529482018990529281018790529081019290925260a08201527fe5c285d336cb17bb08823b603864963ca7aedc5a4d3fea30d299112cb47ddd5a9060c00160405180910390a150505050505050565b6000811561075e57600554610757908590640100000000810463ffffffff16906d0100000000000000000000000000900464ffffffffff16866108b7565b9050610792565b60055461078f908590640100000000810463ffffffff169068010000000000000000900464ffffffffff1686610932565b90505b9392505050565b6107a1610834565b6003805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f23b99d3a969380aa9df8e7afd6d3dbff42d352acaae63d51ad0466d62a1a917d9101610387565b610828610834565b61083181610a63565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108b5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610497565b565b6000806108c2610b59565b905060005a6108d188886111b6565b6108db9190611370565b6108e59085611333565b60045490915060009061090f90879073ffffffffffffffffffffffffffffffffffffffff16610c00565b90508261091c82846111b6565b61092691906111b6565b98975050505050505050565b60008061093d610d55565b90506000811361097c576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610497565b6000610986610b59565b9050600082825a6109978b8b6111b6565b6109a19190611370565b6109ab9088611333565b6109b591906111b6565b6109c790670de0b6b3a7640000611333565b6109d191906111ce565b6002549091506000906109fb90889073ffffffffffffffffffffffffffffffffffffffff16610c00565b9050610a13816b033b2e3c9fd0803ce8000000611370565b821115610a4c576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5681836111b6565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610ae3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610497565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b1811480610b6e575062066eed81145b15610bf857606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b158015610bba57600080fd5b505afa158015610bce573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bf2919061107b565b91505090565b600091505090565b600064ffffffffff8316610c1657506000610d4f565b600080610c2284610e40565b909250905060008213610c80576040517fc3388fe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260248101839052604401610497565b600860ff82161015610cdb576000610c99826008611387565b9050610ca681600a61126a565b610cb09084611333565b610cc964ffffffffff8816670de0b6b3a7640000611333565b610cd391906111ce565b935050610d4c565b600860ff82161115610d25576000610cf4600883611387565b905082610d0282600a61126a565b610d1b64ffffffffff8916670de0b6b3a7640000611333565b610cc99190611333565b81610d3f64ffffffffff8716670de0b6b3a7640000611333565b610d4991906111ce565b92505b50505b92915050565b600554600254604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009363ffffffff1692831515928592839273ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b158015610dd257600080fd5b505afa158015610de6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e0a9190611143565b509450909250849150508015610e2e5750610e258242611370565b8463ffffffff16105b15610e3857506006545b949350505050565b600554604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600092839263ffffffff9091169182151591849173ffffffffffffffffffffffffffffffffffffffff88169163feaf968c9160048083019260a0929190829003018186803b158015610ebd57600080fd5b505afa158015610ed1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ef59190611143565b50919750909250839150508015610f1a5750610f118142611370565b8363ffffffff16105b15610f755760045473ffffffffffffffffffffffffffffffffffffffff87811691161415610f4c576007549450610f75565b60025473ffffffffffffffffffffffffffffffffffffffff878116911614156101005760065494505b8573ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b158015610fbb57600080fd5b505afa158015610fcf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ff39190611193565b9350505050915091565b803563ffffffff8116811461101157600080fd5b919050565b803564ffffffffff8116811461101157600080fd5b805169ffffffffffffffffffff8116811461101157600080fd5b60006020828403121561105757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461079257600080fd5b60006020828403121561108d57600080fd5b5051919050565b6000806000606084860312156110a957600080fd5b8335925060208401359150604084013580151581146110c757600080fd5b809150509250925092565b600080600080600080600060e0888a0312156110ed57600080fd5b6110f688610ffd565b965061110460208901610ffd565b955060408801359450606088013593506080880135925061112760a08901611016565b915061113560c08901611016565b905092959891949750929550565b600080600080600060a0868803121561115b57600080fd5b6111648661102b565b94506020860151935060408601519250606086015191506111876080870161102b565b90509295509295909350565b6000602082840312156111a557600080fd5b815160ff8116811461079257600080fd5b600082198211156111c9576111c96113aa565b500190565b600082611204577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600181815b8085111561126257817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611248576112486113aa565b8085161561125557918102915b93841c939080029061120e565b509250929050565b600061079260ff84168360008261128357506001610d4f565b8161129057506000610d4f565b81600181146112a657600281146112b0576112cc565b6001915050610d4f565b60ff8411156112c1576112c16113aa565b50506001821b610d4f565b5060208310610133831016604e8410600b84101617156112ef575081810a610d4f565b6112f98383611209565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561132b5761132b6113aa565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561136b5761136b6113aa565b500290565b600082821015611382576113826113aa565b500390565b600060ff821660ff8416808210156113a1576113a16113aa565b90039392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c6343000806000a",
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCaller) SLinkEthFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusPriceRegistry.contract.Call(opts, &out, "s_linkEthFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SLinkEthFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SLinkEthFeed(&_VRFV2PlusPriceRegistry.CallOpts)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryCallerSession) SLinkEthFeed() (common.Address, error) {
	return _VRFV2PlusPriceRegistry.Contract.SLinkEthFeed(&_VRFV2PlusPriceRegistry.CallOpts)
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactor) SetEthUsdFeed(opts *bind.TransactOpts, ethUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.contract.Transact(opts, "setEthUsdFeed", ethUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistrySession) SetEthUsdFeed(ethUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetEthUsdFeed(&_VRFV2PlusPriceRegistry.TransactOpts, ethUsdFeed)
}

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistryTransactorSession) SetEthUsdFeed(ethUsdFeed common.Address) (*types.Transaction, error) {
	return _VRFV2PlusPriceRegistry.Contract.SetEthUsdFeed(&_VRFV2PlusPriceRegistry.TransactOpts, ethUsdFeed)
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

type SConfig struct {
	StalenessSeconds           uint32
	GasAfterPaymentCalculation uint32
	FulfillmentFlatFeeLinkUSD  *big.Int
	FulfillmentFlatFeeEthUSD   *big.Int
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

func (_VRFV2PlusPriceRegistry *VRFV2PlusPriceRegistry) Address() common.Address {
	return _VRFV2PlusPriceRegistry.address
}

type VRFV2PlusPriceRegistryInterface interface {
	USDFEEDECIMALS(opts *bind.CallOpts) (uint8, error)

	CalculatePaymentAmount(opts *bind.CallOpts, startGas *big.Int, weiPerUnitGas *big.Int, nativePayment bool) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SEthUSDFeed(opts *bind.CallOpts) (common.Address, error)

	SFallbackUSDPerUnitEth(opts *bind.CallOpts) (*big.Int, error)

	SFallbackUSDPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SLinkEthFeed(opts *bind.CallOpts) (common.Address, error)

	SLinkUSDFeed(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, fulfillmentFlatFeeLinkUSD *big.Int, fulfillmentFlatFeeEthUSD *big.Int) (*types.Transaction, error)

	SetEthUsdFeed(opts *bind.TransactOpts, ethUsdFeed common.Address) (*types.Transaction, error)

	SetLINKETHFeed(opts *bind.TransactOpts, linkEthFeed common.Address) (*types.Transaction, error)

	SetLINKUSDFeed(opts *bind.TransactOpts, linkUsdFeed common.Address) (*types.Transaction, error)

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

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
