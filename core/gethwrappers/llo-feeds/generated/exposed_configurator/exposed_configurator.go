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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"ProductionConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"retiredConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"StagingConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_configId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"_signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"}],\"name\":\"exposedReadConfigurationStates\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"},{\"internalType\":\"bytes32[2]\",\"name\":\"configDigest\",\"type\":\"bytes32[2]\"}],\"internalType\":\"structConfigurator.ConfigurationState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"},{\"internalType\":\"bytes32[2]\",\"name\":\"configDigest\",\"type\":\"bytes32[2]\"}],\"internalType\":\"structConfigurator.ConfigurationState\",\"name\":\"state\",\"type\":\"tuple\"}],\"name\":\"exposedSetConfigurationState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"exposedSetIsGreenProduction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isGreenProduction\",\"type\":\"bool\"}],\"name\":\"promoteStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setProductionConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000a5565b50505062000150565b336001600160a01b03821603620000ff5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611e0880620001606000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806379ba509711610081578063dfb533d01161005b578063dfb533d014610278578063e6e7c5a41461028b578063f2fde38b1461029e57600080fd5b806379ba5097146102285780638da5cb5b1461023057806399a073401461025857600080fd5b8063639fec28116100b2578063639fec28146101a357806369a120eb146101b8578063790464e01461021557600080fd5b806301ffc9a7146100d9578063181f5a771461014357806360e72ec914610182575b600080fd5b61012e6100e7366004611404565b7fffffffff00000000000000000000000000000000000000000000000000000000167f40569294000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601281527f436f6e666967757261746f7220302e352e3000000000000000000000000000006020820152905161013a91906114b1565b6101956101903660046117c6565b6102b1565b60405190815260200161013a565b6101b66101b13660046118e7565b61030d565b005b6101b66101c63660046119cc565b60009182526002602052604090912080549115156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b6101b66102233660046119f8565b6103cc565b6101b6610674565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161013a565b61026b610266366004611ad0565b610771565b60405161013a9190611ae9565b6101b66102863660046119f8565b610814565b6101b66102993660046119cc565b610b25565b6101b66102ac366004611b4e565b610efb565b60006102fd8c8c8c8c8c8c8c8c8c8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508e92508d9150610f0f9050565b9c9b505050505050505050505050565b60008281526002602081815260409283902084518154928601519486015115156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff63ffffffff90961668010000000000000000027fffffffffffffffffffffffffffffffffffffffff00000000000000000000000090941667ffffffffffffffff90921691909117929092179390931617825560608301518392916103c5916001840191611365565b5050505050565b85518460ff168060000361040c576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610456576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b610461816003611b98565b82116104b95781610473826003611b98565b61047e906001611bb5565b6040517f9dd9e6d80000000000000000000000000000000000000000000000000000000081526004810192909252602482015260440161044d565b6104c1610fbd565b845160401461052c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f496e76616c6964206f6e636861696e436f6e666967206c656e67746800000000604482015260640161044d565b60208501516040860151600182146105c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f556e737570706f72746564206f6e636861696e436f6e6669672076657273696f60448201527f6e00000000000000000000000000000000000000000000000000000000000000606482015260840161044d565b8015610654576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603b60248201527f7072656465636573736f72436f6e666967446967657374206d7573742062652060448201527f756e73657420666f722070726f64756374696f6e20636f6e6669670000000000606482015260840161044d565b6106678b46308d8d8d8d8d8d6001611040565b5050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146106f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161044d565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6107796113a3565b6000828152600260208181526040928390208351608081018552815467ffffffffffffffff8116825268010000000000000000810463ffffffff16938201939093526c0100000000000000000000000090920460ff161515828501528351808501948590529193909260608501929160018501919082845b8154815260200190600101908083116107f1575050505050815250509050919050565b85518460ff1680600003610854576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610899576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f602482015260440161044d565b6108a4816003611b98565b82116108b65781610473826003611b98565b6108be610fbd565b8451604014610929576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f496e76616c6964206f6e636861696e436f6e666967206c656e67746800000000604482015260640161044d565b60208501516040860151600182146109c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f556e737570706f72746564206f6e636861696e436f6e6669672076657273696f60448201527f6e00000000000000000000000000000000000000000000000000000000000000606482015260840161044d565b60008b81526002602081815260408084208151608081018352815467ffffffffffffffff8116825268010000000000000000810463ffffffff16948201949094526c0100000000000000000000000090930460ff161515838301528151808301928390529293909260608501929091600185019182845b815481526020019060010190808311610a3a57505050505081525050905060008160400151610a6a576000610a6d565b60015b60ff169050600260008e81526020019081526020016000206001018160028110610a9957610a99611bc8565b01548314610b03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f496e76616c6964207072656465636573736f72436f6e66696744696765737400604482015260640161044d565b610b168d46308f8f8f8f8f8f6000611040565b50505050505050505050505050565b610b2d610fbd565b600082815260026020526040902080546c01000000000000000000000000900460ff16151582151514610c08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604160248201527f50726f6d6f746553746167696e67436f6e6669673a206973477265656e50726f60448201527f64756374696f6e206d757374206d6174636820636f6e7472616374207374617460648201527f6500000000000000000000000000000000000000000000000000000000000000608482015260a40161044d565b805467ffffffffffffffff16610cc6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604260248201527f50726f6d6f746553746167696e67436f6e6669673a20436f6e6669672068617360448201527f206e65766572206265656e2073657420666f72207468697320636f6e6669672060648201527f4944000000000000000000000000000000000000000000000000000000000000608482015260a40161044d565b60006001820183610cd8576001610cdb565b60005b60ff1660028110610cee57610cee611bc8565b015403610da3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604660248201527f50726f6d6f746553746167696e67436f6e6669673a20436f6e6669672064696760448201527f657374206d7573742062652073657420666f72207468652073746167696e672060648201527f636f6e6669670000000000000000000000000000000000000000000000000000608482015260a40161044d565b60008160010183610db5576000610db8565b60015b60ff1660028110610dcb57610dcb611bc8565b0154905080610e82576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152605260248201527f50726f6d6f746553746167696e67436f6e6669673a20436f6e6669672064696760448201527f657374206d7573742062652073657420666f7220746865207265746972696e6760648201527f2070726f64756374696f6e20636f6e6669670000000000000000000000000000608482015260a40161044d565b81547fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff1683156c010000000000000000000000008102919091178355604051908152819085907f1062aa08ac6046a0e69e3eafdf12d1eba63a67b71a874623e86eb06348a1d84f9060200160405180910390a350505050565b610f03610fbd565b610f0c81611270565b50565b6000808b8b8b8b8b8b8b8b8b8b604051602001610f359a99989796959493929190611c87565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461103e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161044d565b565b60008a815260026020526040812080549091908290829061106a9067ffffffffffffffff16611d34565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055905060006110a58d8d8d858e8e8e8e8e8e610f0f565b9050831561116d578c7f261b20c2ecd99d86d6e936279e4f78db34603a3de3a4a84d6f3d4e0dd55e24788460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff166040516111149a99989796959493929190611d5b565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16611150576000611153565b60015b60ff166002811061116657611166611bc8565b0155611229565b8c7fef1b5f9d1b927b0fe871b12c7e7846457602d67b2bc36b0bc95feaf480e890568460000160089054906101000a900463ffffffff1683858e8e8e8e8e8e8d600001600c9054906101000a900460ff166040516111d49a99989796959493929190611d5b565b60405180910390a260008d815260026020526040902083548291600101906c01000000000000000000000000900460ff16611210576001611213565b60005b60ff166002811061122657611226611bc8565b01555b505080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000004363ffffffff160217905550505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036112ef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161044d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8260028101928215611393579160200282015b82811115611393578251825591602001919060010190611378565b5061139f9291506113d1565b5090565b6040805160808101825260008082526020820181905291810191909152606081016113cc6113e6565b905290565b5b8082111561139f57600081556001016113d2565b60405180604001604052806002906020820280368337509192915050565b60006020828403121561141657600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461144657600080fd5b9392505050565b6000815180845260005b8181101561147357602081850181015186830182015201611457565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611446602083018461144d565b803573ffffffffffffffffffffffffffffffffffffffff811681146114e857600080fd5b919050565b803567ffffffffffffffff811681146114e857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff8111828210171561155757611557611505565b60405290565b6040805190810167ffffffffffffffff8111828210171561155757611557611505565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156115c7576115c7611505565b604052919050565b600067ffffffffffffffff8211156115e9576115e9611505565b5060051b60200190565b600082601f83011261160457600080fd5b813567ffffffffffffffff81111561161e5761161e611505565b61164f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611580565b81815284602083860101111561166457600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261169257600080fd5b813560206116a76116a2836115cf565b611580565b82815260059290921b840181019181810190868411156116c657600080fd5b8286015b8481101561170657803567ffffffffffffffff8111156116ea5760008081fd5b6116f88986838b01016115f3565b8452509183019183016116ca565b509695505050505050565b600082601f83011261172257600080fd5b813560206117326116a2836115cf565b82815260059290921b8401810191818101908684111561175157600080fd5b8286015b848110156117065780358352918301918301611755565b803560ff811681146114e857600080fd5b60008083601f84011261178f57600080fd5b50813567ffffffffffffffff8111156117a757600080fd5b6020830191508360208285010111156117bf57600080fd5b9250929050565b60008060008060008060008060008060006101408c8e0312156117e857600080fd5b8b359a5060208c013599506117ff60408d016114c4565b985061180d60608d016114ed565b975067ffffffffffffffff8060808e0135111561182957600080fd5b6118398e60808f01358f01611681565b97508060a08e0135111561184c57600080fd5b61185c8e60a08f01358f01611711565b965061186a60c08e0161176c565b95508060e08e0135111561187d57600080fd5b61188d8e60e08f01358f0161177d565b909550935061189f6101008e016114ed565b9250806101208e013511156118b357600080fd5b506118c58d6101208e01358e016115f3565b90509295989b509295989b9093969950565b803580151581146114e857600080fd5b60008082840360c08112156118fb57600080fd5b83359250602060a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08301121561193157600080fd5b611939611534565b91506119468186016114ed565b8252604085013563ffffffff8116811461195f57600080fd5b8282015261196f606086016118d7565b604083015285609f86011261198357600080fd5b61198b61155d565b8060c087018881111561199d57600080fd5b608088015b818110156119b957803584529284019284016119a2565b5050606084015250929590945092505050565b600080604083850312156119df57600080fd5b823591506119ef602084016118d7565b90509250929050565b600080600080600080600060e0888a031215611a1357600080fd5b87359650602088013567ffffffffffffffff80821115611a3257600080fd5b611a3e8b838c01611681565b975060408a0135915080821115611a5457600080fd5b611a608b838c01611711565b9650611a6e60608b0161176c565b955060808a0135915080821115611a8457600080fd5b611a908b838c016115f3565b9450611a9e60a08b016114ed565b935060c08a0135915080821115611ab457600080fd5b50611ac18a828b016115f3565b91505092959891949750929550565b600060208284031215611ae257600080fd5b5035919050565b600060a08201905067ffffffffffffffff8351168252602063ffffffff81850151168184015260408401511515604084015260608401516060840160005b6002811015611b4457825182529183019190830190600101611b27565b5050505092915050565b600060208284031215611b6057600080fd5b611446826114c4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417611baf57611baf611b69565b92915050565b80820180821115611baf57611baf611b69565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081518084526020808501808196508360051b8101915082860160005b85811015611c3f578284038952611c2d84835161144d565b98850198935090840190600101611c15565b5091979650505050505050565b600081518084526020808501945080840160005b83811015611c7c57815187529582019590820190600101611c60565b509495945050505050565b60006101408c83528b602084015273ffffffffffffffffffffffffffffffffffffffff8b16604084015267ffffffffffffffff808b166060850152816080850152611cd48285018b611bf7565b915083820360a0850152611ce8828a611c4c565b915060ff881660c085015283820360e0850152611d05828861144d565b9086166101008501528381036101208501529050611d23818561144d565b9d9c50505050505050505050505050565b600067ffffffffffffffff808316818103611d5157611d51611b69565b6001019392505050565b600061014063ffffffff8d1683528b602084015267ffffffffffffffff808c166040850152816060850152611d928285018c611bf7565b91508382036080850152611da6828b611c4c565b915060ff891660a085015283820360c0850152611dc3828961144d565b90871660e08501528381036101008501529050611de0818661144d565b9150508215156101208301529b9a505050505050505050505056fea164736f6c6343000813000a",
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
