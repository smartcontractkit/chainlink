// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exposed_configurator

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

type ConfiguratorConfigurationState struct {
	ConfigCount             uint64
	LatestConfigBlockNumber uint32
	IsGreenProduction       bool
	ConfigDigest            [2][32]byte
}

var ExposedConfiguratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"}],\"name\":\"ConfigUnset\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ConfigUnsetProduction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ConfigUnsetStaging\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"onchainConfigLength\",\"type\":\"uint256\"}],\"name\":\"InvalidOnchainLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"predecessorConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"InvalidPredecessorConfigDigest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProductionContractState\",\"type\":\"bool\"}],\"name\":\"IsGreenProductionMustMatchContractState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"predecessorConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"NonZeroPredecessorConfigDigest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"UnsupportedOnchainConfigVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ProductionConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"retiredConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"StagingConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_configId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"_signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"}],\"name\":\"exposedReadConfigurationStates\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"},{\"internalType\":\"bytes32[2]\",\"name\":\"configDigest\",\"type\":\"bytes32[2]\"}],\"internalType\":\"structConfigurator.ConfigurationState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"},{\"internalType\":\"bytes32[2]\",\"name\":\"configDigest\",\"type\":\"bytes32[2]\"}],\"internalType\":\"structConfigurator.ConfigurationState\",\"name\":\"state\",\"type\":\"tuple\"}],\"name\":\"exposedSetConfigurationState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"exposedSetIsGreenProduction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"promoteStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setProductionConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611ac080620001586000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806379ba509711610081578063dfb533d01161005b578063dfb533d014610278578063e6e7c5a41461028b578063f2fde38b1461029e57600080fd5b806379ba5097146102285780638da5cb5b1461023057806399a073401461025857600080fd5b8063639fec28116100b2578063639fec28146101a357806369a120eb146101b8578063790464e01461021557600080fd5b806301ffc9a7146100d9578063181f5a771461014357806360e72ec914610182575b600080fd5b61012e6100e73660046110bc565b7fffffffff00000000000000000000000000000000000000000000000000000000167f40569294000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601281527f436f6e666967757261746f7220302e352e3000000000000000000000000000006020820152905161013a9190611169565b61019561019036600461147e565b6102b1565b60405190815260200161013a565b6101b66101b136600461159f565b61030d565b005b6101b66101c6366004611684565b60009182526002602052604090912080549115156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b6101b66102233660046116b0565b6103cc565b6101b66105a5565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161013a565b61026b610266366004611788565b6106a2565b60405161013a91906117a1565b6101b66102863660046116b0565b610745565b6101b6610299366004611684565b6109a7565b6101b66102ac366004611806565b610bb3565b60006102fd8c8c8c8c8c8c8c8c8c8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508e92508d9150610bc79050565b9c9b505050505050505050505050565b60008281526002602081815260409283902084518154928601519486015115156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff63ffffffff90961668010000000000000000027fffffffffffffffffffffffffffffffffffffffff00000000000000000000000090941667ffffffffffffffff90921691909117929092179390931617825560608301518392916103c591600184019161101d565b5050505050565b85518460ff168060000361040c576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610456576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b610461816003611850565b82116104b95781610473826003611850565b61047e90600161186d565b6040517f9dd9e6d80000000000000000000000000000000000000000000000000000000081526004810192909252602482015260440161044d565b6104c1610c75565b6040855110156105025784516040517f3e936ca800000000000000000000000000000000000000000000000000000000815260040161044d91815260200190565b60208501516040860151600182101561054a576040517f8f01e0d70000000000000000000000000000000000000000000000000000000081526004810183905260240161044d565b8015610585576040517fb96bb7600000000000000000000000000000000000000000000000000000000081526004810182905260240161044d565b6105988b46308d8d8d8d8d8d6001610cf8565b5050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610626576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161044d565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6106aa61105b565b6000828152600260208181526040928390208351608081018552815467ffffffffffffffff8116825268010000000000000000810463ffffffff16938201939093526c0100000000000000000000000090920460ff161515828501528351808501948590529193909260608501929160018501919082845b815481526020019060010190808311610722575050505050815250509050919050565b85518460ff1680600003610785576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f8211156107ca576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f602482015260440161044d565b6107d5816003611850565b82116107e75781610473826003611850565b6107ef610c75565b6040855110156108305784516040517f3e936ca800000000000000000000000000000000000000000000000000000000815260040161044d91815260200190565b602085015160408601516001821015610878576040517f8f01e0d70000000000000000000000000000000000000000000000000000000081526004810183905260240161044d565b60008b81526002602081815260408084208151608081018352815467ffffffffffffffff8116825268010000000000000000810463ffffffff16948201949094526c0100000000000000000000000090930460ff161515838301528151808301928390529293909260608501929091600185019182845b8154815260200190600101908083116108ef575050505050815250509050600260008d81526020019081526020016000206001018160400151610933576000610936565b60015b60ff166002811061094957610949611880565b01548214610986576040517f7d78c2a10000000000000000000000000000000000000000000000000000000081526004810183905260240161044d565b6109998c46308e8e8e8e8e8e6000610cf8565b505050505050505050505050565b6109af610c75565b600082815260026020526040902080546c01000000000000000000000000900460ff16151582151514610a18576040517f85fa3a37000000000000000000000000000000000000000000000000000000008152600481018490528215602482015260440161044d565b805467ffffffffffffffff16600003610a60576040517f90e6f6dc0000000000000000000000000000000000000000000000000000000081526004810184905260240161044d565b60006001820183610a72576001610a75565b60005b60ff1660028110610a8857610a88611880565b015403610acc576040517f5b7f635700000000000000000000000000000000000000000000000000000000815260048101849052821515602482015260440161044d565b60008160010183610ade576000610ae1565b60015b60ff1660028110610af457610af4611880565b0154905080610b3a576040517fcaf1e77300000000000000000000000000000000000000000000000000000000815260048101859052831515602482015260440161044d565b81547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff1683156c010000000000000000000000008102919091178355604051908152819085907f1062aa08ac6046a0e69e3eafdf12d1eba63a67b71a874623e86eb06348a1d84f9060200160405180910390a350505050565b610bbb610c75565b610bc481610f28565b50565b6000808b8b8b8b8b8b8b8b8b8b604051602001610bed9a9998979695949392919061193f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610cf6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161044d565b565b60008a8152600260205260408120805490919082908290610d229067ffffffffffffffff166119ec565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905590506000610d5d8d8d8d858e8e8e8e8e8e610bc7565b90508315610e25578c7f261b20c2ecd99d86d6e936279e4f78db34603a3de3a4a84d6f3d4e0dd55e24788460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff16604051610dcc9a99989796959493929190611a13565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16610e08576000610e0b565b60015b60ff1660028110610e1e57610e1e611880565b0155610ee1565b8c7fef1b5f9d1b927b0fe871b12c7e7846457602d67b2bc36b0bc95feaf480e890568460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff16604051610e8c9a99989796959493929190611a13565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16610ec8576001610ecb565b60005b60ff1660028110610ede57610ede611880565b01555b505080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000004363ffffffff160217905550505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610fa7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161044d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b826002810192821561104b579160200282015b8281111561104b578251825591602001919060010190611030565b50611057929150611089565b5090565b60408051608081018252600080825260208201819052918101919091526060810161108461109e565b905290565b5b80821115611057576000815560010161108a565b60405180604001604052806002906020820280368337509192915050565b6000602082840312156110ce57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146110fe57600080fd5b9392505050565b6000815180845260005b8181101561112b5760208185018101518683018201520161110f565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006110fe6020830184611105565b803573ffffffffffffffffffffffffffffffffffffffff811681146111a057600080fd5b919050565b803567ffffffffffffffff811681146111a057600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff8111828210171561120f5761120f6111bd565b60405290565b6040805190810167ffffffffffffffff8111828210171561120f5761120f6111bd565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561127f5761127f6111bd565b604052919050565b600067ffffffffffffffff8211156112a1576112a16111bd565b5060051b60200190565b600082601f8301126112bc57600080fd5b813567ffffffffffffffff8111156112d6576112d66111bd565b61130760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611238565b81815284602083860101111561131c57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261134a57600080fd5b8135602061135f61135a83611287565b611238565b82815260059290921b8401810191818101908684111561137e57600080fd5b8286015b848110156113be57803567ffffffffffffffff8111156113a25760008081fd5b6113b08986838b01016112ab565b845250918301918301611382565b509695505050505050565b600082601f8301126113da57600080fd5b813560206113ea61135a83611287565b82815260059290921b8401810191818101908684111561140957600080fd5b8286015b848110156113be578035835291830191830161140d565b803560ff811681146111a057600080fd5b60008083601f84011261144757600080fd5b50813567ffffffffffffffff81111561145f57600080fd5b60208301915083602082850101111561147757600080fd5b9250929050565b60008060008060008060008060008060006101408c8e0312156114a057600080fd5b8b359a5060208c013599506114b760408d0161117c565b98506114c560608d016111a5565b975067ffffffffffffffff8060808e013511156114e157600080fd5b6114f18e60808f01358f01611339565b97508060a08e0135111561150457600080fd5b6115148e60a08f01358f016113c9565b965061152260c08e01611424565b95508060e08e0135111561153557600080fd5b6115458e60e08f01358f01611435565b90955093506115576101008e016111a5565b9250806101208e0135111561156b57600080fd5b5061157d8d6101208e01358e016112ab565b90509295989b509295989b9093969950565b803580151581146111a057600080fd5b60008082840360c08112156115b357600080fd5b83359250602060a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0830112156115e957600080fd5b6115f16111ec565b91506115fe8186016111a5565b8252604085013563ffffffff8116811461161757600080fd5b828201526116276060860161158f565b604083015285609f86011261163b57600080fd5b611643611215565b8060c087018881111561165557600080fd5b608088015b81811015611671578035845292840192840161165a565b5050606084015250929590945092505050565b6000806040838503121561169757600080fd5b823591506116a76020840161158f565b90509250929050565b600080600080600080600060e0888a0312156116cb57600080fd5b87359650602088013567ffffffffffffffff808211156116ea57600080fd5b6116f68b838c01611339565b975060408a013591508082111561170c57600080fd5b6117188b838c016113c9565b965061172660608b01611424565b955060808a013591508082111561173c57600080fd5b6117488b838c016112ab565b945061175660a08b016111a5565b935060c08a013591508082111561176c57600080fd5b506117798a828b016112ab565b91505092959891949750929550565b60006020828403121561179a57600080fd5b5035919050565b600060a08201905067ffffffffffffffff8351168252602063ffffffff81850151168184015260408401511515604084015260608401516060840160005b60028110156117fc578251825291830191908301906001016117df565b5050505092915050565b60006020828403121561181857600080fd5b6110fe8261117c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761186757611867611821565b92915050565b8082018082111561186757611867611821565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081518084526020808501808196508360051b8101915082860160005b858110156118f75782840389526118e5848351611105565b988501989350908401906001016118cd565b5091979650505050505050565b600081518084526020808501945080840160005b8381101561193457815187529582019590820190600101611918565b509495945050505050565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b16606085015281608085015261198c8285018b6118af565b915083820360a08501526119a0828a611904565b915060ff881660c085015283820360e08501526119bd8288611105565b90861661010085015283810361012085015290506119db8185611105565b9d9c50505050505050505050505050565b600067ffffffffffffffff808316818103611a0957611a09611821565b6001019392505050565b600061014063ffffffff8d1683528b602084015267ffffffffffffffff808c166040850152816060850152611a4a8285018c6118af565b91508382036080850152611a5e828b611904565b915060ff891660a085015283820360c0850152611a7b8289611105565b90871660e08501528381036101008501529050611a988186611105565b9150508215156101208301529b9a505050505050505050505056fea164736f6c6343000813000a",
}

var ExposedConfiguratorABI = ExposedConfiguratorMetaData.ABI

var ExposedConfiguratorBin = ExposedConfiguratorMetaData.Bin

func DeployExposedConfigurator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExposedConfigurator, error) {
	parsed, err := ExposedConfiguratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExposedConfiguratorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExposedConfigurator{address: address, abi: *parsed, ExposedConfiguratorCaller: ExposedConfiguratorCaller{contract: contract}, ExposedConfiguratorTransactor: ExposedConfiguratorTransactor{contract: contract}, ExposedConfiguratorFilterer: ExposedConfiguratorFilterer{contract: contract}}, nil
}

type ExposedConfigurator struct {
	address common.Address
	abi     abi.ABI
	ExposedConfiguratorCaller
	ExposedConfiguratorTransactor
	ExposedConfiguratorFilterer
}

type ExposedConfiguratorCaller struct {
	contract *bind.BoundContract
}

type ExposedConfiguratorTransactor struct {
	contract *bind.BoundContract
}

type ExposedConfiguratorFilterer struct {
	contract *bind.BoundContract
}

type ExposedConfiguratorSession struct {
	Contract     *ExposedConfigurator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ExposedConfiguratorCallerSession struct {
	Contract *ExposedConfiguratorCaller
	CallOpts bind.CallOpts
}

type ExposedConfiguratorTransactorSession struct {
	Contract     *ExposedConfiguratorTransactor
	TransactOpts bind.TransactOpts
}

type ExposedConfiguratorRaw struct {
	Contract *ExposedConfigurator
}

type ExposedConfiguratorCallerRaw struct {
	Contract *ExposedConfiguratorCaller
}

type ExposedConfiguratorTransactorRaw struct {
	Contract *ExposedConfiguratorTransactor
}

func NewExposedConfigurator(address common.Address, backend bind.ContractBackend) (*ExposedConfigurator, error) {
	abi, err := abi.JSON(strings.NewReader(ExposedConfiguratorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindExposedConfigurator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExposedConfigurator{address: address, abi: abi, ExposedConfiguratorCaller: ExposedConfiguratorCaller{contract: contract}, ExposedConfiguratorTransactor: ExposedConfiguratorTransactor{contract: contract}, ExposedConfiguratorFilterer: ExposedConfiguratorFilterer{contract: contract}}, nil
}

func NewExposedConfiguratorCaller(address common.Address, caller bind.ContractCaller) (*ExposedConfiguratorCaller, error) {
	contract, err := bindExposedConfigurator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorCaller{contract: contract}, nil
}

func NewExposedConfiguratorTransactor(address common.Address, transactor bind.ContractTransactor) (*ExposedConfiguratorTransactor, error) {
	contract, err := bindExposedConfigurator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorTransactor{contract: contract}, nil
}

func NewExposedConfiguratorFilterer(address common.Address, filterer bind.ContractFilterer) (*ExposedConfiguratorFilterer, error) {
	contract, err := bindExposedConfigurator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorFilterer{contract: contract}, nil
}

func bindExposedConfigurator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExposedConfiguratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ExposedConfigurator *ExposedConfiguratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExposedConfigurator.Contract.ExposedConfiguratorCaller.contract.Call(opts, result, method, params...)
}

func (_ExposedConfigurator *ExposedConfiguratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.ExposedConfiguratorTransactor.contract.Transfer(opts)
}

func (_ExposedConfigurator *ExposedConfiguratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.ExposedConfiguratorTransactor.contract.Transact(opts, method, params...)
}

func (_ExposedConfigurator *ExposedConfiguratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExposedConfigurator.Contract.contract.Call(opts, result, method, params...)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.contract.Transfer(opts)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.contract.Transact(opts, method, params...)
}

func (_ExposedConfigurator *ExposedConfiguratorCaller) ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _configId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers [][]byte, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	var out []interface{}
	err := _ExposedConfigurator.contract.Call(opts, &out, "exposedConfigDigestFromConfigData", _configId, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ExposedConfigurator *ExposedConfiguratorSession) ExposedConfigDigestFromConfigData(_configId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers [][]byte, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _ExposedConfigurator.Contract.ExposedConfigDigestFromConfigData(&_ExposedConfigurator.CallOpts, _configId, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorCallerSession) ExposedConfigDigestFromConfigData(_configId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers [][]byte, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _ExposedConfigurator.Contract.ExposedConfigDigestFromConfigData(&_ExposedConfigurator.CallOpts, _configId, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorCaller) ExposedReadConfigurationStates(opts *bind.CallOpts, configId [32]byte) (ConfiguratorConfigurationState, error) {
	var out []interface{}
	err := _ExposedConfigurator.contract.Call(opts, &out, "exposedReadConfigurationStates", configId)

	if err != nil {
		return *new(ConfiguratorConfigurationState), err
	}

	out0 := *abi.ConvertType(out[0], new(ConfiguratorConfigurationState)).(*ConfiguratorConfigurationState)

	return out0, err

}

func (_ExposedConfigurator *ExposedConfiguratorSession) ExposedReadConfigurationStates(configId [32]byte) (ConfiguratorConfigurationState, error) {
	return _ExposedConfigurator.Contract.ExposedReadConfigurationStates(&_ExposedConfigurator.CallOpts, configId)
}

func (_ExposedConfigurator *ExposedConfiguratorCallerSession) ExposedReadConfigurationStates(configId [32]byte) (ConfiguratorConfigurationState, error) {
	return _ExposedConfigurator.Contract.ExposedReadConfigurationStates(&_ExposedConfigurator.CallOpts, configId)
}

func (_ExposedConfigurator *ExposedConfiguratorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExposedConfigurator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ExposedConfigurator *ExposedConfiguratorSession) Owner() (common.Address, error) {
	return _ExposedConfigurator.Contract.Owner(&_ExposedConfigurator.CallOpts)
}

func (_ExposedConfigurator *ExposedConfiguratorCallerSession) Owner() (common.Address, error) {
	return _ExposedConfigurator.Contract.Owner(&_ExposedConfigurator.CallOpts)
}

func (_ExposedConfigurator *ExposedConfiguratorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ExposedConfigurator.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ExposedConfigurator *ExposedConfiguratorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ExposedConfigurator.Contract.SupportsInterface(&_ExposedConfigurator.CallOpts, interfaceId)
}

func (_ExposedConfigurator *ExposedConfiguratorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ExposedConfigurator.Contract.SupportsInterface(&_ExposedConfigurator.CallOpts, interfaceId)
}

func (_ExposedConfigurator *ExposedConfiguratorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ExposedConfigurator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ExposedConfigurator *ExposedConfiguratorSession) TypeAndVersion() (string, error) {
	return _ExposedConfigurator.Contract.TypeAndVersion(&_ExposedConfigurator.CallOpts)
}

func (_ExposedConfigurator *ExposedConfiguratorCallerSession) TypeAndVersion() (string, error) {
	return _ExposedConfigurator.Contract.TypeAndVersion(&_ExposedConfigurator.CallOpts)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "acceptOwnership")
}

func (_ExposedConfigurator *ExposedConfiguratorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.AcceptOwnership(&_ExposedConfigurator.TransactOpts)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.AcceptOwnership(&_ExposedConfigurator.TransactOpts)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) ExposedSetConfigurationState(opts *bind.TransactOpts, configId [32]byte, state ConfiguratorConfigurationState) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "exposedSetConfigurationState", configId, state)
}

func (_ExposedConfigurator *ExposedConfiguratorSession) ExposedSetConfigurationState(configId [32]byte, state ConfiguratorConfigurationState) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.ExposedSetConfigurationState(&_ExposedConfigurator.TransactOpts, configId, state)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) ExposedSetConfigurationState(configId [32]byte, state ConfiguratorConfigurationState) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.ExposedSetConfigurationState(&_ExposedConfigurator.TransactOpts, configId, state)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) ExposedSetIsGreenProduction(opts *bind.TransactOpts, configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "exposedSetIsGreenProduction", configId, isGreenProduction)
}

func (_ExposedConfigurator *ExposedConfiguratorSession) ExposedSetIsGreenProduction(configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.ExposedSetIsGreenProduction(&_ExposedConfigurator.TransactOpts, configId, isGreenProduction)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) ExposedSetIsGreenProduction(configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.ExposedSetIsGreenProduction(&_ExposedConfigurator.TransactOpts, configId, isGreenProduction)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) PromoteStagingConfig(opts *bind.TransactOpts, configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "promoteStagingConfig", configId, isGreenProduction)
}

func (_ExposedConfigurator *ExposedConfiguratorSession) PromoteStagingConfig(configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.PromoteStagingConfig(&_ExposedConfigurator.TransactOpts, configId, isGreenProduction)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) PromoteStagingConfig(configId [32]byte, isGreenProduction bool) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.PromoteStagingConfig(&_ExposedConfigurator.TransactOpts, configId, isGreenProduction)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) SetProductionConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "setProductionConfig", configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorSession) SetProductionConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.SetProductionConfig(&_ExposedConfigurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) SetProductionConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.SetProductionConfig(&_ExposedConfigurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) SetStagingConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "setStagingConfig", configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorSession) SetStagingConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.SetStagingConfig(&_ExposedConfigurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) SetStagingConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.SetStagingConfig(&_ExposedConfigurator.TransactOpts, configId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ExposedConfigurator.contract.Transact(opts, "transferOwnership", to)
}

func (_ExposedConfigurator *ExposedConfiguratorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.TransferOwnership(&_ExposedConfigurator.TransactOpts, to)
}

func (_ExposedConfigurator *ExposedConfiguratorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ExposedConfigurator.Contract.TransferOwnership(&_ExposedConfigurator.TransactOpts, to)
}

type ExposedConfiguratorOwnershipTransferRequestedIterator struct {
	Event *ExposedConfiguratorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ExposedConfiguratorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedConfiguratorOwnershipTransferRequested)
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
		it.Event = new(ExposedConfiguratorOwnershipTransferRequested)
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

func (it *ExposedConfiguratorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ExposedConfiguratorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ExposedConfiguratorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedConfiguratorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorOwnershipTransferRequestedIterator{contract: _ExposedConfigurator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ExposedConfiguratorOwnershipTransferRequested)
				if err := _ExposedConfigurator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) ParseOwnershipTransferRequested(log types.Log) (*ExposedConfiguratorOwnershipTransferRequested, error) {
	event := new(ExposedConfiguratorOwnershipTransferRequested)
	if err := _ExposedConfigurator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ExposedConfiguratorOwnershipTransferredIterator struct {
	Event *ExposedConfiguratorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ExposedConfiguratorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedConfiguratorOwnershipTransferred)
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
		it.Event = new(ExposedConfiguratorOwnershipTransferred)
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

func (it *ExposedConfiguratorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ExposedConfiguratorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ExposedConfiguratorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedConfiguratorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorOwnershipTransferredIterator{contract: _ExposedConfigurator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ExposedConfiguratorOwnershipTransferred)
				if err := _ExposedConfigurator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) ParseOwnershipTransferred(log types.Log) (*ExposedConfiguratorOwnershipTransferred, error) {
	event := new(ExposedConfiguratorOwnershipTransferred)
	if err := _ExposedConfigurator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ExposedConfiguratorProductionConfigSetIterator struct {
	Event *ExposedConfiguratorProductionConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ExposedConfiguratorProductionConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedConfiguratorProductionConfigSet)
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
		it.Event = new(ExposedConfiguratorProductionConfigSet)
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

func (it *ExposedConfiguratorProductionConfigSetIterator) Error() error {
	return it.fail
}

func (it *ExposedConfiguratorProductionConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ExposedConfiguratorProductionConfigSet struct {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) FilterProductionConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ExposedConfiguratorProductionConfigSetIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.FilterLogs(opts, "ProductionConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorProductionConfigSetIterator{contract: _ExposedConfigurator.contract, event: "ProductionConfigSet", logs: logs, sub: sub}, nil
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) WatchProductionConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorProductionConfigSet, configId [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.WatchLogs(opts, "ProductionConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ExposedConfiguratorProductionConfigSet)
				if err := _ExposedConfigurator.contract.UnpackLog(event, "ProductionConfigSet", log); err != nil {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) ParseProductionConfigSet(log types.Log) (*ExposedConfiguratorProductionConfigSet, error) {
	event := new(ExposedConfiguratorProductionConfigSet)
	if err := _ExposedConfigurator.contract.UnpackLog(event, "ProductionConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ExposedConfiguratorPromoteStagingConfigIterator struct {
	Event *ExposedConfiguratorPromoteStagingConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ExposedConfiguratorPromoteStagingConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedConfiguratorPromoteStagingConfig)
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
		it.Event = new(ExposedConfiguratorPromoteStagingConfig)
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

func (it *ExposedConfiguratorPromoteStagingConfigIterator) Error() error {
	return it.fail
}

func (it *ExposedConfiguratorPromoteStagingConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ExposedConfiguratorPromoteStagingConfig struct {
	ConfigId            [32]byte
	RetiredConfigDigest [32]byte
	IsGreenProduction   bool
	Raw                 types.Log
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) FilterPromoteStagingConfig(opts *bind.FilterOpts, configId [][32]byte, retiredConfigDigest [][32]byte) (*ExposedConfiguratorPromoteStagingConfigIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}
	var retiredConfigDigestRule []interface{}
	for _, retiredConfigDigestItem := range retiredConfigDigest {
		retiredConfigDigestRule = append(retiredConfigDigestRule, retiredConfigDigestItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.FilterLogs(opts, "PromoteStagingConfig", configIdRule, retiredConfigDigestRule)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorPromoteStagingConfigIterator{contract: _ExposedConfigurator.contract, event: "PromoteStagingConfig", logs: logs, sub: sub}, nil
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorPromoteStagingConfig, configId [][32]byte, retiredConfigDigest [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}
	var retiredConfigDigestRule []interface{}
	for _, retiredConfigDigestItem := range retiredConfigDigest {
		retiredConfigDigestRule = append(retiredConfigDigestRule, retiredConfigDigestItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.WatchLogs(opts, "PromoteStagingConfig", configIdRule, retiredConfigDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ExposedConfiguratorPromoteStagingConfig)
				if err := _ExposedConfigurator.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) ParsePromoteStagingConfig(log types.Log) (*ExposedConfiguratorPromoteStagingConfig, error) {
	event := new(ExposedConfiguratorPromoteStagingConfig)
	if err := _ExposedConfigurator.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ExposedConfiguratorStagingConfigSetIterator struct {
	Event *ExposedConfiguratorStagingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ExposedConfiguratorStagingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedConfiguratorStagingConfigSet)
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
		it.Event = new(ExposedConfiguratorStagingConfigSet)
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

func (it *ExposedConfiguratorStagingConfigSetIterator) Error() error {
	return it.fail
}

func (it *ExposedConfiguratorStagingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ExposedConfiguratorStagingConfigSet struct {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) FilterStagingConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ExposedConfiguratorStagingConfigSetIterator, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.FilterLogs(opts, "StagingConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return &ExposedConfiguratorStagingConfigSetIterator{contract: _ExposedConfigurator.contract, event: "StagingConfigSet", logs: logs, sub: sub}, nil
}

func (_ExposedConfigurator *ExposedConfiguratorFilterer) WatchStagingConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorStagingConfigSet, configId [][32]byte) (event.Subscription, error) {

	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _ExposedConfigurator.contract.WatchLogs(opts, "StagingConfigSet", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ExposedConfiguratorStagingConfigSet)
				if err := _ExposedConfigurator.contract.UnpackLog(event, "StagingConfigSet", log); err != nil {
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

func (_ExposedConfigurator *ExposedConfiguratorFilterer) ParseStagingConfigSet(log types.Log) (*ExposedConfiguratorStagingConfigSet, error) {
	event := new(ExposedConfiguratorStagingConfigSet)
	if err := _ExposedConfigurator.contract.UnpackLog(event, "StagingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ExposedConfigurator *ExposedConfigurator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ExposedConfigurator.abi.Events["OwnershipTransferRequested"].ID:
		return _ExposedConfigurator.ParseOwnershipTransferRequested(log)
	case _ExposedConfigurator.abi.Events["OwnershipTransferred"].ID:
		return _ExposedConfigurator.ParseOwnershipTransferred(log)
	case _ExposedConfigurator.abi.Events["ProductionConfigSet"].ID:
		return _ExposedConfigurator.ParseProductionConfigSet(log)
	case _ExposedConfigurator.abi.Events["PromoteStagingConfig"].ID:
		return _ExposedConfigurator.ParsePromoteStagingConfig(log)
	case _ExposedConfigurator.abi.Events["StagingConfigSet"].ID:
		return _ExposedConfigurator.ParseStagingConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ExposedConfiguratorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ExposedConfiguratorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (ExposedConfiguratorProductionConfigSet) Topic() common.Hash {
	return common.HexToHash("0x261b20c2ecd99d86d6e936279e4f78db34603a3de3a4a84d6f3d4e0dd55e2478")
}

func (ExposedConfiguratorPromoteStagingConfig) Topic() common.Hash {
	return common.HexToHash("0x1062aa08ac6046a0e69e3eafdf12d1eba63a67b71a874623e86eb06348a1d84f")
}

func (ExposedConfiguratorStagingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xef1b5f9d1b927b0fe871b12c7e7846457602d67b2bc36b0bc95feaf480e89056")
}

func (_ExposedConfigurator *ExposedConfigurator) Address() common.Address {
	return _ExposedConfigurator.address
}

type ExposedConfiguratorInterface interface {
	ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _configId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers [][]byte, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error)

	ExposedReadConfigurationStates(opts *bind.CallOpts, configId [32]byte) (ConfiguratorConfigurationState, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ExposedSetConfigurationState(opts *bind.TransactOpts, configId [32]byte, state ConfiguratorConfigurationState) (*types.Transaction, error)

	ExposedSetIsGreenProduction(opts *bind.TransactOpts, configId [32]byte, isGreenProduction bool) (*types.Transaction, error)

	PromoteStagingConfig(opts *bind.TransactOpts, configId [32]byte, isGreenProduction bool) (*types.Transaction, error)

	SetProductionConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetStagingConfig(opts *bind.TransactOpts, configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedConfiguratorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ExposedConfiguratorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedConfiguratorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ExposedConfiguratorOwnershipTransferred, error)

	FilterProductionConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ExposedConfiguratorProductionConfigSetIterator, error)

	WatchProductionConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorProductionConfigSet, configId [][32]byte) (event.Subscription, error)

	ParseProductionConfigSet(log types.Log) (*ExposedConfiguratorProductionConfigSet, error)

	FilterPromoteStagingConfig(opts *bind.FilterOpts, configId [][32]byte, retiredConfigDigest [][32]byte) (*ExposedConfiguratorPromoteStagingConfigIterator, error)

	WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorPromoteStagingConfig, configId [][32]byte, retiredConfigDigest [][32]byte) (event.Subscription, error)

	ParsePromoteStagingConfig(log types.Log) (*ExposedConfiguratorPromoteStagingConfig, error)

	FilterStagingConfigSet(opts *bind.FilterOpts, configId [][32]byte) (*ExposedConfiguratorStagingConfigSetIterator, error)

	WatchStagingConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedConfiguratorStagingConfigSet, configId [][32]byte) (event.Subscription, error)

	ParseStagingConfigSet(log types.Log) (*ExposedConfiguratorStagingConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
