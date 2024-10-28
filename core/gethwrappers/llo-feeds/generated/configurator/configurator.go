// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package configurator

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

var ConfiguratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"}],\"name\":\"ConfigUnset\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ConfigUnsetProduction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ConfigUnsetStaging\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"onchainConfigLength\",\"type\":\"uint256\"}],\"name\":\"InvalidOnchainLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"predecessorConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"InvalidPredecessorConfigDigest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProductionContractState\",\"type\":\"bool\"}],\"name\":\"IsGreenProductionMustMatchContractState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"predecessorConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"NonZeroPredecessorConfigDigest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"UnsupportedOnchainConfigVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ProductionConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"retiredConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"StagingConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"promoteStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setProductionConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61144a806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80638da5cb5b1161005b5780638da5cb5b14610153578063dfb533d01461017b578063e6e7c5a41461018e578063f2fde38b146101a157600080fd5b806301ffc9a71461008d578063181f5a77146100f7578063790464e01461013657806379ba50971461014b575b600080fd5b6100e261009b366004610d62565b7fffffffff00000000000000000000000000000000000000000000000000000000167f40569294000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601281527f436f6e666967757261746f7220302e352e300000000000000000000000000000602082015290516100ee9190610e0f565b61014961014436600461106b565b6101b4565b005b61014961038d565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ee565b61014961018936600461106b565b61048a565b61014961019c366004611143565b6106ec565b6101496101af366004611178565b6108f8565b85518460ff16806000036101f4576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f82111561023e576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b6102498160036111dd565b82116102a1578161025b8260036111dd565b6102669060016111fa565b6040517f9dd9e6d800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610235565b6102a961090c565b6040855110156102ea5784516040517f3e936ca800000000000000000000000000000000000000000000000000000000815260040161023591815260200190565b602085015160408601516001821015610332576040517f8f01e0d700000000000000000000000000000000000000000000000000000000815260048101839052602401610235565b801561036d576040517fb96bb76000000000000000000000000000000000000000000000000000000000815260048101829052602401610235565b6103808b46308d8d8d8d8d8d600161098f565b5050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461040e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610235565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b85518460ff16806000036104ca576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f82111561050f576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f6024820152604401610235565b61051a8160036111dd565b821161052c578161025b8260036111dd565b61053461090c565b6040855110156105755784516040517f3e936ca800000000000000000000000000000000000000000000000000000000815260040161023591815260200190565b6020850151604086015160018210156105bd576040517f8f01e0d700000000000000000000000000000000000000000000000000000000815260048101839052602401610235565b60008b81526002602081815260408084208151608081018352815467ffffffffffffffff8116825268010000000000000000810463ffffffff16948201949094526c0100000000000000000000000090930460ff161515838301528151808301928390529293909260608501929091600185019182845b815481526020019060010190808311610634575050505050815250509050600260008d8152602001908152602001600020600101816040015161067857600061067b565b60015b60ff166002811061068e5761068e61120d565b015482146106cb576040517f7d78c2a100000000000000000000000000000000000000000000000000000000815260048101839052602401610235565b6106de8c46308e8e8e8e8e8e600061098f565b505050505050505050505050565b6106f461090c565b600082815260026020526040902080546c01000000000000000000000000900460ff1615158215151461075d576040517f85fa3a370000000000000000000000000000000000000000000000000000000081526004810184905282156024820152604401610235565b805467ffffffffffffffff166000036107a5576040517f90e6f6dc00000000000000000000000000000000000000000000000000000000815260048101849052602401610235565b600060018201836107b75760016107ba565b60005b60ff16600281106107cd576107cd61120d565b015403610811576040517f5b7f6357000000000000000000000000000000000000000000000000000000008152600481018490528215156024820152604401610235565b60008160010183610823576000610826565b60015b60ff16600281106108395761083961120d565b015490508061087f576040517fcaf1e773000000000000000000000000000000000000000000000000000000008152600481018590528315156024820152604401610235565b81547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff1683156c010000000000000000000000008102919091178355604051908152819085907f1062aa08ac6046a0e69e3eafdf12d1eba63a67b71a874623e86eb06348a1d84f9060200160405180910390a350505050565b61090061090c565b61090981610bbf565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461098d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610235565b565b60008a81526002602052604081208054909190829082906109b99067ffffffffffffffff1661123c565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055905060006109f48d8d8d858e8e8e8e8e8e610cb4565b90508315610abc578c7f261b20c2ecd99d86d6e936279e4f78db34603a3de3a4a84d6f3d4e0dd55e24788460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff16604051610a639a999897969594939291906112f0565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16610a9f576000610aa2565b60015b60ff1660028110610ab557610ab561120d565b0155610b78565b8c7fef1b5f9d1b927b0fe871b12c7e7846457602d67b2bc36b0bc95feaf480e890568460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff16604051610b239a999897969594939291906112f0565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16610b5f576001610b62565b60005b60ff1660028110610b7557610b7561120d565b01555b505080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000004363ffffffff160217905550505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610c3e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610235565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808b8b8b8b8b8b8b8b8b8b604051602001610cda9a99989796959493929190611390565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b600060208284031215610d7457600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610da457600080fd5b9392505050565b6000815180845260005b81811015610dd157602081850181015186830182015201610db5565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610da46020830184610dab565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610e9857610e98610e22565b604052919050565b600067ffffffffffffffff821115610eba57610eba610e22565b5060051b60200190565b600082601f830112610ed557600080fd5b813567ffffffffffffffff811115610eef57610eef610e22565b610f2060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610e51565b818152846020838601011115610f3557600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f830112610f6357600080fd5b81356020610f78610f7383610ea0565b610e51565b82815260059290921b84018101918181019086841115610f9757600080fd5b8286015b84811015610fd757803567ffffffffffffffff811115610fbb5760008081fd5b610fc98986838b0101610ec4565b845250918301918301610f9b565b509695505050505050565b600082601f830112610ff357600080fd5b81356020611003610f7383610ea0565b82815260059290921b8401810191818101908684111561102257600080fd5b8286015b84811015610fd75780358352918301918301611026565b803560ff8116811461104e57600080fd5b919050565b803567ffffffffffffffff8116811461104e57600080fd5b600080600080600080600060e0888a03121561108657600080fd5b87359650602088013567ffffffffffffffff808211156110a557600080fd5b6110b18b838c01610f52565b975060408a01359150808211156110c757600080fd5b6110d38b838c01610fe2565b96506110e160608b0161103d565b955060808a01359150808211156110f757600080fd5b6111038b838c01610ec4565b945061111160a08b01611053565b935060c08a013591508082111561112757600080fd5b506111348a828b01610ec4565b91505092959891949750929550565b6000806040838503121561115657600080fd5b823591506020830135801515811461116d57600080fd5b809150509250929050565b60006020828403121561118a57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610da457600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176111f4576111f46111ae565b92915050565b808201808211156111f4576111f46111ae565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600067ffffffffffffffff808316818103611259576112596111ae565b6001019392505050565b6000815180845260208085019450848260051b860182860160005b858110156112a8578383038952611296838351610dab565b9885019892509084019060010161127e565b5090979650505050505050565b600081518084526020808501945080840160005b838110156112e5578151875295820195908201906001016112c9565b509495945050505050565b600061014063ffffffff8d1683528b602084015267ffffffffffffffff808c1660408501528160608501526113278285018c611263565b9150838203608085015261133b828b6112b5565b915060ff891660a085015283820360c08501526113588289610dab565b90871660e085015283810361010085015290506113758186610dab565b9150508215156101208301529b9a5050505050505050505050565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b1660608501528160808501526113dd8285018b611263565b915083820360a08501526113f1828a6112b5565b915060ff881660c085015283820360e085015261140e8288610dab565b908616610100850152838103610120850152905061142c8185610dab565b9d9c5050505050505050505050505056fea164736f6c6343000813000a",
}

var ConfiguratorABI = ConfiguratorMetaData.ABI

var ConfiguratorBin = ConfiguratorMetaData.Bin

func DeployConfigurator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Configurator, error) {
	parsed, err := ConfiguratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfiguratorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Configurator{address: address, abi: *parsed, ConfiguratorCaller: ConfiguratorCaller{contract: contract}, ConfiguratorTransactor: ConfiguratorTransactor{contract: contract}, ConfiguratorFilterer: ConfiguratorFilterer{contract: contract}}, nil
}

type Configurator struct {
	address common.Address
	abi     abi.ABI
	ConfiguratorCaller
	ConfiguratorTransactor
	ConfiguratorFilterer
}

type ConfiguratorCaller struct {
	contract *bind.BoundContract
}

type ConfiguratorTransactor struct {
	contract *bind.BoundContract
}

type ConfiguratorFilterer struct {
	contract *bind.BoundContract
}

type ConfiguratorSession struct {
	Contract     *Configurator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ConfiguratorCallerSession struct {
	Contract *ConfiguratorCaller
	CallOpts bind.CallOpts
}

type ConfiguratorTransactorSession struct {
	Contract     *ConfiguratorTransactor
	TransactOpts bind.TransactOpts
}

type ConfiguratorRaw struct {
	Contract *Configurator
}

type ConfiguratorCallerRaw struct {
	Contract *ConfiguratorCaller
}

type ConfiguratorTransactorRaw struct {
	Contract *ConfiguratorTransactor
}

func NewConfigurator(address common.Address, backend bind.ContractBackend) (*Configurator, error) {
	abi, err := abi.JSON(strings.NewReader(ConfiguratorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindConfigurator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Configurator{address: address, abi: abi, ConfiguratorCaller: ConfiguratorCaller{contract: contract}, ConfiguratorTransactor: ConfiguratorTransactor{contract: contract}, ConfiguratorFilterer: ConfiguratorFilterer{contract: contract}}, nil
}

func NewConfiguratorCaller(address common.Address, caller bind.ContractCaller) (*ConfiguratorCaller, error) {
	contract, err := bindConfigurator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorCaller{contract: contract}, nil
}

func NewConfiguratorTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfiguratorTransactor, error) {
	contract, err := bindConfigurator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorTransactor{contract: contract}, nil
}

func NewConfiguratorFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfiguratorFilterer, error) {
	contract, err := bindConfigurator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorFilterer{contract: contract}, nil
}

func bindConfigurator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfiguratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Configurator *ConfiguratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Configurator.Contract.ConfiguratorCaller.contract.Call(opts, result, method, params...)
}

func (_Configurator *ConfiguratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Configurator.Contract.ConfiguratorTransactor.contract.Transfer(opts)
}

func (_Configurator *ConfiguratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Configurator.Contract.ConfiguratorTransactor.contract.Transact(opts, method, params...)
}

func (_Configurator *ConfiguratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Configurator.Contract.contract.Call(opts, result, method, params...)
}

func (_Configurator *ConfiguratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Configurator.Contract.contract.Transfer(opts)
}

func (_Configurator *ConfiguratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Configurator.Contract.contract.Transact(opts, method, params...)
}

func (_Configurator *ConfiguratorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Configurator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Configurator *ConfiguratorSession) Owner() (common.Address, error) {
	return _Configurator.Contract.Owner(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorCallerSession) Owner() (common.Address, error) {
	return _Configurator.Contract.Owner(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Configurator.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Configurator *ConfiguratorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Configurator.Contract.SupportsInterface(&_Configurator.CallOpts, interfaceId)
}

func (_Configurator *ConfiguratorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Configurator.Contract.SupportsInterface(&_Configurator.CallOpts, interfaceId)
}

func (_Configurator *ConfiguratorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Configurator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Configurator *ConfiguratorSession) TypeAndVersion() (string, error) {
	return _Configurator.Contract.TypeAndVersion(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorCallerSession) TypeAndVersion() (string, error) {
	return _Configurator.Contract.TypeAndVersion(&_Configurator.CallOpts)
}

func (_Configurator *ConfiguratorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "acceptOwnership")
}

func (_Configurator *ConfiguratorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Configurator.Contract.AcceptOwnership(&_Configurator.TransactOpts)
}

func (_Configurator *ConfiguratorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Configurator.Contract.AcceptOwnership(&_Configurator.TransactOpts)
}

func (_Configurator *ConfiguratorTransactor) PromoteStagingConfig(opts *bind.TransactOpts, configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "promoteStagingConfig", configId, isGreenProduction)
}

func (_Configurator *ConfiguratorSession) PromoteStagingConfig(configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _Configurator.Contract.PromoteStagingConfig(&_Configurator.TransactOpts, configId, isGreenProduction)
}

func (_Configurator *ConfiguratorTransactorSession) PromoteStagingConfig(configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _Configurator.Contract.PromoteStagingConfig(&_Configurator.TransactOpts, configId, isGreenProduction)
}

func (_Configurator *ConfiguratorTransactor) SetProductionConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "setProductionConfig", configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorSession) SetProductionConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.Contract.SetProductionConfig(&_Configurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorTransactorSession) SetProductionConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.Contract.SetProductionConfig(&_Configurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorTransactor) SetStagingConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "setStagingConfig", configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorSession) SetStagingConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.Contract.SetStagingConfig(&_Configurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorTransactorSession) SetStagingConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Configurator.Contract.SetStagingConfig(&_Configurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_Configurator *ConfiguratorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Configurator.contract.Transact(opts, "transferOwnership", to)
}

func (_Configurator *ConfiguratorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Configurator.Contract.TransferOwnership(&_Configurator.TransactOpts, to)
}

func (_Configurator *ConfiguratorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Configurator.Contract.TransferOwnership(&_Configurator.TransactOpts, to)
}

type ConfiguratorOwnershipTransferRequestedIterator struct {
	Event *ConfiguratorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorOwnershipTransferRequested)
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
		it.Event = new(ConfiguratorOwnershipTransferRequested)
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

func (it *ConfiguratorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorOwnershipTransferRequestedIterator{contract: _Configurator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorOwnershipTransferRequested)
				if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfiguratorOwnershipTransferRequested, error) {
	event := new(ConfiguratorOwnershipTransferRequested)
	if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConfiguratorOwnershipTransferredIterator struct {
	Event *ConfiguratorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorOwnershipTransferred)
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
		it.Event = new(ConfiguratorOwnershipTransferred)
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

func (it *ConfiguratorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorOwnershipTransferredIterator{contract: _Configurator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorOwnershipTransferred)
				if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseOwnershipTransferred(log types.Log) (*ConfiguratorOwnershipTransferred, error) {
	event := new(ConfiguratorOwnershipTransferred)
	if err := _Configurator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConfiguratorProductionConfigSetIterator struct {
	Event *ConfiguratorProductionConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorProductionConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorProductionConfigSet)
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
		it.Event = new(ConfiguratorProductionConfigSet)
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

func (it *ConfiguratorProductionConfigSetIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorProductionConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorProductionConfigSet struct {
	ConfigId                  [32]byte
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   [][]byte
	OffchainTransmitters      [][32]byte
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	IsGreenProduction         bool
	Raw                       types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterProductionConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ConfiguratorProductionConfigSetIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "ProductionConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorProductionConfigSetIterator{contract: _Configurator.contract, event: "ProductionConfigSet", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchProductionConfigSet(opts *bind.WatchOpts, sink chan<- *ConfiguratorProductionConfigSet, configId [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "ProductionConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorProductionConfigSet)
				if err := _Configurator.contract.UnpackLog(event, "ProductionConfigSet", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseProductionConfigSet(log types.Log) (*ConfiguratorProductionConfigSet, error) {
	event := new(ConfiguratorProductionConfigSet)
	if err := _Configurator.contract.UnpackLog(event, "ProductionConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConfiguratorPromoteStagingConfigIterator struct {
	Event *ConfiguratorPromoteStagingConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorPromoteStagingConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorPromoteStagingConfig)
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
		it.Event = new(ConfiguratorPromoteStagingConfig)
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

func (it *ConfiguratorPromoteStagingConfigIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorPromoteStagingConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorPromoteStagingConfig struct {
	ConfigId            [32]byte
	RetiredConfigDigest [32]byte
	IsGreenProduction   bool
	Raw                 types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterPromoteStagingConfig(opts *bind.FilterOpts, configId [][32]byte, retiredConfigDigest [][32]byte) (*ConfiguratorPromoteStagingConfigIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}
	var retiredConfigDigestRule []interface{}
	for _, retiredConfigDigestItem := range retiredConfigDigest {
		retiredConfigDigestRule = append(retiredConfigDigestRule, retiredConfigDigestItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "PromoteStagingConfig", configIdRule, retiredConfigDigestRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorPromoteStagingConfigIterator{contract: _Configurator.contract, event: "PromoteStagingConfig", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *ConfiguratorPromoteStagingConfig, configId [][32]byte, retiredConfigDigest [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}
	var retiredConfigDigestRule []interface{}
	for _, retiredConfigDigestItem := range retiredConfigDigest {
		retiredConfigDigestRule = append(retiredConfigDigestRule, retiredConfigDigestItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "PromoteStagingConfig", configIdRule, retiredConfigDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorPromoteStagingConfig)
				if err := _Configurator.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParsePromoteStagingConfig(log types.Log) (*ConfiguratorPromoteStagingConfig, error) {
	event := new(ConfiguratorPromoteStagingConfig)
	if err := _Configurator.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConfiguratorStagingConfigSetIterator struct {
	Event *ConfiguratorStagingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConfiguratorStagingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorStagingConfigSet)
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
		it.Event = new(ConfiguratorStagingConfigSet)
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

func (it *ConfiguratorStagingConfigSetIterator) Error() error {
	return it.fail
}

func (it *ConfiguratorStagingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConfiguratorStagingConfigSet struct {
	ConfigId                  [32]byte
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   [][]byte
	OffchainTransmitters      [][32]byte
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	IsGreenProduction         bool
	Raw                       types.Log
}

func (_Configurator *ConfiguratorFilterer) FilterStagingConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ConfiguratorStagingConfigSetIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _Configurator.contract.FilterLogs(opts, "StagingConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorStagingConfigSetIterator{contract: _Configurator.contract, event: "StagingConfigSet", logs: logs, sub: sub}, nil
}

func (_Configurator *ConfiguratorFilterer) WatchStagingConfigSet(opts *bind.WatchOpts, sink chan<- *ConfiguratorStagingConfigSet, configId [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _Configurator.contract.WatchLogs(opts, "StagingConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConfiguratorStagingConfigSet)
				if err := _Configurator.contract.UnpackLog(event, "StagingConfigSet", log); err != nil {
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

func (_Configurator *ConfiguratorFilterer) ParseStagingConfigSet(log types.Log) (*ConfiguratorStagingConfigSet, error) {
	event := new(ConfiguratorStagingConfigSet)
	if err := _Configurator.contract.UnpackLog(event, "StagingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Configurator *Configurator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Configurator.abi.Events["OwnershipTransferRequested"].ID:
		return _Configurator.ParseOwnershipTransferRequested(log)
	case _Configurator.abi.Events["OwnershipTransferred"].ID:
		return _Configurator.ParseOwnershipTransferred(log)
	case _Configurator.abi.Events["ProductionConfigSet"].ID:
		return _Configurator.ParseProductionConfigSet(log)
	case _Configurator.abi.Events["PromoteStagingConfig"].ID:
		return _Configurator.ParsePromoteStagingConfig(log)
	case _Configurator.abi.Events["StagingConfigSet"].ID:
		return _Configurator.ParseStagingConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ConfiguratorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ConfiguratorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (ConfiguratorProductionConfigSet) Topic() common.Hash {
	return common.HexToHash("0x261b20c2ecd99d86d6e936279e4f78db34603a3de3a4a84d6f3d4e0dd55e2478")
}

func (ConfiguratorPromoteStagingConfig) Topic() common.Hash {
	return common.HexToHash("0x1062aa08ac6046a0e69e3eafdf12d1eba63a67b71a874623e86eb06348a1d84f")
}

func (ConfiguratorStagingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xef1b5f9d1b927b0fe871b12c7e7846457602d67b2bc36b0bc95feaf480e89056")
}

func (_Configurator *Configurator) Address() common.Address {
	return _Configurator.address
}

type ConfiguratorInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	PromoteStagingConfig(opts *bind.TransactOpts, configId [32]byte, isGreenProduction bool) (*types.Transaction, error)

	SetProductionConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetStagingConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ConfiguratorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfiguratorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfiguratorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ConfiguratorOwnershipTransferred, error)

	FilterProductionConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ConfiguratorProductionConfigSetIterator, error)

	WatchProductionConfigSet(opts *bind.WatchOpts, sink chan<- *ConfiguratorProductionConfigSet, configId [][32]byte) (event.Subscription, error)

	ParseProductionConfigSet(log types.Log) (*ConfiguratorProductionConfigSet, error)

	FilterPromoteStagingConfig(opts *bind.FilterOpts, configId [][32]byte, retiredConfigDigest [][32]byte) (*ConfiguratorPromoteStagingConfigIterator, error)

	WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *ConfiguratorPromoteStagingConfig, configId [][32]byte, retiredConfigDigest [][32]byte) (event.Subscription, error)

	ParsePromoteStagingConfig(log types.Log) (*ConfiguratorPromoteStagingConfig, error)

	FilterStagingConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ConfiguratorStagingConfigSetIterator, error)

	WatchStagingConfigSet(opts *bind.WatchOpts, sink chan<- *ConfiguratorStagingConfigSet, configId [][32]byte) (event.Subscription, error)

	ParseStagingConfigSet(log types.Log) (*ConfiguratorStagingConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
