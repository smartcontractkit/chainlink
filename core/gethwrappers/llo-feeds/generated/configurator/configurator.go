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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ProductionConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"retiredConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"StagingConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"promoteStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setProductionConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611795806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80638da5cb5b1161005b5780638da5cb5b14610153578063dfb533d01461017b578063e6e7c5a41461018e578063f2fde38b146101a157600080fd5b806301ffc9a71461008d578063181f5a77146100f7578063790464e01461013657806379ba50971461014b575b600080fd5b6100e261009b3660046110aa565b7fffffffff00000000000000000000000000000000000000000000000000000000167f40569294000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601281527f436f6e666967757261746f7220302e352e300000000000000000000000000000602082015290516100ee9190611157565b6101496101443660046113b3565b6101b4565b005b61014961045c565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ee565b6101496101893660046113b3565b610559565b61014961019c36600461148b565b61086a565b6101496101af3660046114c0565b610c40565b85518460ff16806000036101f4576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f82111561023e576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b610249816003611525565b82116102a1578161025b826003611525565b610266906001611542565b6040517f9dd9e6d800000000000000000000000000000000000000000000000000000000815260048101929092526024820152604401610235565b6102a9610c54565b8451604014610314576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f496e76616c6964206f6e636861696e436f6e666967206c656e677468000000006044820152606401610235565b60208501516040860151600182146103ae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f556e737570706f72746564206f6e636861696e436f6e6669672076657273696f60448201527f6e000000000000000000000000000000000000000000000000000000000000006064820152608401610235565b801561043c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603b60248201527f7072656465636573736f72436f6e666967446967657374206d7573742062652060448201527f756e73657420666f722070726f64756374696f6e20636f6e66696700000000006064820152608401610235565b61044f8b46308d8d8d8d8d8d6001610cd7565b5050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104dd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610235565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b85518460ff1680600003610599576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f8211156105de576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f6024820152604401610235565b6105e9816003611525565b82116105fb578161025b826003611525565b610603610c54565b845160401461066e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f496e76616c6964206f6e636861696e436f6e666967206c656e677468000000006044820152606401610235565b6020850151604086015160018214610708576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f556e737570706f72746564206f6e636861696e436f6e6669672076657273696f60448201527f6e000000000000000000000000000000000000000000000000000000000000006064820152608401610235565b60008b81526002602081815260408084208151608081018352815467ffffffffffffffff8116825268010000000000000000810463ffffffff16948201949094526c0100000000000000000000000090930460ff161515838301528151808301928390529293909260608501929091600185019182845b81548152602001906001019080831161077f575050505050815250509050600081604001516107af5760006107b2565b60015b60ff169050600260008e815260200190815260200160002060010181600281106107de576107de611555565b01548314610848576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f496e76616c6964207072656465636573736f72436f6e666967446967657374006044820152606401610235565b61085b8d46308f8f8f8f8f8f6000610cd7565b50505050505050505050505050565b610872610c54565b600082815260026020526040902080546c01000000000000000000000000900460ff1615158215151461094d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604160248201527f50726f6d6f746553746167696e67436f6e6669673a206973477265656e50726f60448201527f64756374696f6e206d757374206d6174636820636f6e7472616374207374617460648201527f6500000000000000000000000000000000000000000000000000000000000000608482015260a401610235565b805467ffffffffffffffff16610a0b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604260248201527f50726f6d6f746553746167696e67436f6e6669673a20436f6e6669672068617360448201527f206e65766572206265656e2073657420666f72207468697320636f6e6669672060648201527f4944000000000000000000000000000000000000000000000000000000000000608482015260a401610235565b60006001820183610a1d576001610a20565b60005b60ff1660028110610a3357610a33611555565b015403610ae8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604660248201527f50726f6d6f746553746167696e67436f6e6669673a20436f6e6669672064696760448201527f657374206d7573742062652073657420666f72207468652073746167696e672060648201527f636f6e6669670000000000000000000000000000000000000000000000000000608482015260a401610235565b60008160010183610afa576000610afd565b60015b60ff1660028110610b1057610b10611555565b0154905080610bc7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152605260248201527f50726f6d6f746553746167696e67436f6e6669673a20436f6e6669672064696760448201527f657374206d7573742062652073657420666f7220746865207265746972696e6760648201527f2070726f64756374696f6e20636f6e6669670000000000000000000000000000608482015260a401610235565b81547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff1683156c010000000000000000000000008102919091178355604051908152819085907f1062aa08ac6046a0e69e3eafdf12d1eba63a67b71a874623e86eb06348a1d84f9060200160405180910390a350505050565b610c48610c54565b610c5181610f07565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610cd5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610235565b565b60008a8152600260205260408120805490919082908290610d019067ffffffffffffffff16611584565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905590506000610d3c8d8d8d858e8e8e8e8e8e610ffc565b90508315610e04578c7f261b20c2ecd99d86d6e936279e4f78db34603a3de3a4a84d6f3d4e0dd55e24788460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff16604051610dab9a9998979695949392919061163b565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16610de7576000610dea565b60015b60ff1660028110610dfd57610dfd611555565b0155610ec0565b8c7fef1b5f9d1b927b0fe871b12c7e7846457602d67b2bc36b0bc95feaf480e890568460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff16604051610e6b9a9998979695949392919061163b565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16610ea7576001610eaa565b60005b60ff1660028110610ebd57610ebd611555565b01555b505080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000004363ffffffff160217905550505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610f86576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610235565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808b8b8b8b8b8b8b8b8b8b6040516020016110229a999897969594939291906116db565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b6000602082840312156110bc57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146110ec57600080fd5b9392505050565b6000815180845260005b81811015611119576020818501810151868301820152016110fd565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006110ec60208301846110f3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111e0576111e061116a565b604052919050565b600067ffffffffffffffff8211156112025761120261116a565b5060051b60200190565b600082601f83011261121d57600080fd5b813567ffffffffffffffff8111156112375761123761116a565b61126860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611199565b81815284602083860101111561127d57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126112ab57600080fd5b813560206112c06112bb836111e8565b611199565b82815260059290921b840181019181810190868411156112df57600080fd5b8286015b8481101561131f57803567ffffffffffffffff8111156113035760008081fd5b6113118986838b010161120c565b8452509183019183016112e3565b509695505050505050565b600082601f83011261133b57600080fd5b8135602061134b6112bb836111e8565b82815260059290921b8401810191818101908684111561136a57600080fd5b8286015b8481101561131f578035835291830191830161136e565b803560ff8116811461139657600080fd5b919050565b803567ffffffffffffffff8116811461139657600080fd5b600080600080600080600060e0888a0312156113ce57600080fd5b87359650602088013567ffffffffffffffff808211156113ed57600080fd5b6113f98b838c0161129a565b975060408a013591508082111561140f57600080fd5b61141b8b838c0161132a565b965061142960608b01611385565b955060808a013591508082111561143f57600080fd5b61144b8b838c0161120c565b945061145960a08b0161139b565b935060c08a013591508082111561146f57600080fd5b5061147c8a828b0161120c565b91505092959891949750929550565b6000806040838503121561149e57600080fd5b82359150602083013580151581146114b557600080fd5b809150509250929050565b6000602082840312156114d257600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146110ec57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761153c5761153c6114f6565b92915050565b8082018082111561153c5761153c6114f6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600067ffffffffffffffff8083168181036115a1576115a16114f6565b6001019392505050565b600081518084526020808501808196508360051b8101915082860160005b858110156115f35782840389526115e18483516110f3565b988501989350908401906001016115c9565b5091979650505050505050565b600081518084526020808501945080840160005b8381101561163057815187529582019590820190600101611614565b509495945050505050565b600061014063ffffffff8d1683528b602084015267ffffffffffffffff808c1660408501528160608501526116728285018c6115ab565b91508382036080850152611686828b611600565b915060ff891660a085015283820360c08501526116a382896110f3565b90871660e085015283810361010085015290506116c081866110f3565b9150508215156101208301529b9a5050505050505050505050565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b1660608501528160808501526117288285018b6115ab565b915083820360a085015261173c828a611600565b915060ff881660c085015283820360e085015261175982886110f3565b908616610100850152838103610120850152905061177781856110f3565b9d9c5050505050505050505050505056fea164736f6c6343000813000a",
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
