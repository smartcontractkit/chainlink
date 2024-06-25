// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package liquiditymanager

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

type ILiquidityManagerCrossChainRebalancerArgs struct {
	RemoteRebalancer    common.Address
	LocalBridge         common.Address
	RemoteToken         common.Address
	RemoteChainSelector uint64
	Enabled             bool
}

type LiquidityManagerCrossChainRebalancer struct {
	RemoteRebalancer common.Address
	LocalBridge      common.Address
	RemoteToken      common.Address
	Enabled          bool
}

var LiquidityManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"finance\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidRemoteChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"CrossChainRebalancerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"FinalizationFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"}],\"name\":\"FinalizationStepCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"financeRole\",\"type\":\"address\"}],\"name\":\"FinanceRoleSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAddedToContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newLiquidityContainer\",\"type\":\"address\"}],\"name\":\"LiquidityContainerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemovedFromContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"fromChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"toChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeReturnData\",\"type\":\"bytes\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"MinimumLiquiditySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"NativeDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"NativeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainRebalancers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getCrossChainRebalancer\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityManager.CrossChainRebalancer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFinanceRole\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalLiquidityContainer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedDestChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_localToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nativeBridgeFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"rebalanceLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"shouldWrapNative\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"receiveLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs\",\"name\":\"crossChainLiqManager\",\"type\":\"tuple\"}],\"name\":\"setCrossChainRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"crossChainRebalancers\",\"type\":\"tuple[]\"}],\"name\":\"setCrossChainRebalancers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"finance\",\"type\":\"address\"}],\"name\":\"setFinanceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"}],\"name\":\"setLocalLiquidityContainer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"}],\"name\":\"setMinimumLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004aa638038062004aa6833981016040819052620000349162000239565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000175565b505046608052506001600160401b038416600003620000f05760405163f89d762960e01b815260040160405180910390fd5b6001600160a01b03851615806200010e57506001600160a01b038316155b156200012d5760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0394851660a0526001600160401b0390931660c052600b80549285166001600160a01b0319938416179055600855600c8054929093169116179055620002b8565b336001600160a01b03821603620001cf5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200023657600080fd5b50565b600080600080600060a086880312156200025257600080fd5b85516200025f8162000220565b60208701519095506001600160401b03811681146200027d57600080fd5b6040870151909450620002908162000220565b606087015160808801519194509250620002aa8162000220565b809150509295509295909350565b60805160a05160c0516147786200032e6000396000818161317d015261335201526000818161045a01528181610757015281816109f401528181610a3a015281816116c00152818161305b015281816130dc0152818161320501526132a301526000818161185d01526118a901526147786000f3fe6080604052600436106101b05760003560e01c8063791781f5116100ec578063b7e7fa051161008a578063f1c0461611610064578063f1c046161461066b578063f2fde38b146106aa578063f8c2d8fa146106ca578063fe65d5af146106ea57600080fd5b8063b7e7fa0514610600578063b8ca8dd814610620578063da9c0f961461064057600080fd5b806383d34afe116100c657806383d34afe146105805780638da5cb5b146105955780639c8f9f23146105c0578063b1dc65a4146105e057600080fd5b8063791781f51461050357806379ba50971461052e57806381ff70481461054357600080fd5b806350a197d7116101595780636511d919116101335780636511d91914610448578063666cab8d146104a15780636a11ee90146104c3578063706bf645146104e357600080fd5b806350a197d7146102f657806351c6590a14610316578063568446e71461033657600080fd5b80633275636e1161018a5780633275636e14610294578063348759c1146102b45780634f814d04146102d657600080fd5b80630910a510146101f4578063181f5a771461021c578063282567b41461027257600080fd5b366101ef57604080513481523360208201527f3c597f6ac9fe7f0ed6da50b07618f5850a642e459ad587f7fab491a71f8b0ab8910160405180910390a1005b600080fd5b34801561020057600080fd5b5061020961070c565b6040519081526020015b60405180910390f35b34801561022857600080fd5b506102656040518060400160405280601a81526020017f4c69717569646974794d616e6167657220312e302e302d64657600000000000081525081565b6040516102139190613929565b34801561027e57600080fd5b5061029261028d366004613943565b6107c7565b005b3480156102a057600080fd5b506102926102af36600461395c565b610814565b3480156102c057600080fd5b506102c9610828565b6040516102139190613974565b3480156102e257600080fd5b506102926102f13660046139e4565b6108b4565b34801561030257600080fd5b50610292610311366004613a75565b610935565b34801561032257600080fd5b50610292610331366004613943565b6109da565b34801561034257600080fd5b506103f9610351366004613ae6565b6040805160808101825260008082526020820181905291810182905260608101919091525067ffffffffffffffff166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff161515606082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff9081168252602080850151821690830152838301511691810191909152606091820151151591810191909152608001610213565b34801561045457600080fd5b5061047c7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610213565b3480156104ad57600080fd5b506104b6610b17565b6040516102139190613b53565b3480156104cf57600080fd5b506102926104de366004613d74565b610b85565b3480156104ef57600080fd5b506102926104fe3660046139e4565b6113b9565b34801561050f57600080fd5b50600b5473ffffffffffffffffffffffffffffffffffffffff1661047c565b34801561053a57600080fd5b5061029261147d565b34801561054f57600080fd5b506004546002546040805163ffffffff80851682526401000000009094049093166020840152820152606001610213565b34801561058c57600080fd5b50600854610209565b3480156105a157600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661047c565b3480156105cc57600080fd5b506102926105db366004613943565b61157a565b3480156105ec57600080fd5b506102926105fb366004613e86565b61171a565b34801561060c57600080fd5b5061029261061b366004613f3d565b611d8b565b34801561062c57600080fd5b5061029261063b366004613fb2565b611dcb565b34801561064c57600080fd5b50600c5473ffffffffffffffffffffffffffffffffffffffff1661047c565b34801561067757600080fd5b5060045468010000000000000000900467ffffffffffffffff1660405167ffffffffffffffff9091168152602001610213565b3480156106b657600080fd5b506102926106c53660046139e4565b611f09565b3480156106d657600080fd5b506102926106e5366004613fe2565b611f1a565b3480156106f657600080fd5b506106ff611fb6565b604051610213919061402d565b600b546040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526000917f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa15801561079e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107c291906140c2565b905090565b6107cf612177565b600880549082905560408051828152602081018490527ff97e758c8b3d81df7b0e1b7327a6a7fcf09a41536b2d274b9103015d715f11eb910160405180910390a15050565b61081c612177565b610825816121fa565b50565b6060600a8054806020026020016040519081016040528092919081815260200182805480156108aa57602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116108655790505b5050505050905090565b6108bc612177565b600c80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f58024d20c07d3ebb87b192861d337d3a60995665acc5b8ce29596458b1f251709060200160405180910390a150565b600c5473ffffffffffffffffffffffffffffffffffffffff163314610986576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109d3858584848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525089925067ffffffffffffffff91506125f69050565b5050505050565b610a1c73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633308461287b565b600b54610a639073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691168361295d565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b158015610acf57600080fd5b505af1158015610ae3573d6000803e3d6000fd5b50506040518392503391507f5414b81d05ac3542606f164e16a9a107d05d21e906539cc5ceb61d7b6b707eb590600090a350565b606060078054806020026020016040519081016040528092919081815260200182805480156108aa57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610b51575050505050905090565b855185518560ff16601f831115610bfd576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b80600003610c67576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610bf4565b818314610cf5576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610bf4565b610d0081600361410a565b8311610d68576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610bf4565b610d70612177565b60065460005b81811015610e64576005600060068381548110610d9557610d95614127565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560078054600592919084908110610e0557610e05614127565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600101610d76565b50895160005b818110156112375760008c8281518110610e8657610e86614127565b6020026020010151905060006002811115610ea357610ea3614156565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff166002811115610ee257610ee2614156565b14610f49576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610bf4565b73ffffffffffffffffffffffffffffffffffffffff8116610f96576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561104657611046614156565b021790555090505060008c838151811061106257611062614127565b602002602001015190506000600281111561107f5761107f614156565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff1660028111156110be576110be614156565b14611125576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610bf4565b73ffffffffffffffffffffffffffffffffffffffff8116611172576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561122257611222614156565b02179055509050505050806001019050610e6a565b508a5161124b9060069060208e01906137fd565b50895161125f9060079060208d01906137fd565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c1617179055600480546112e59146913091906000906112b79063ffffffff16614185565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e612adf565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055506000600460086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f6040516113a3999897969594939291906141a8565b60405180910390a1505050505050505050505050565b6113c1612177565b73ffffffffffffffffffffffffffffffffffffffff811661140e576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517f07dc474694ac40123aadcd2445f1b38d2eb353edd9319dcea043548ab34990ec90600090a250565b60015473ffffffffffffffffffffffffffffffffffffffff1633146114fe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610bf4565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600c5473ffffffffffffffffffffffffffffffffffffffff1633146115cb576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006115d561070c565b905081811015611622576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018390526024810182905260006044820152606401610bf4565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b15801561168e57600080fd5b505af11580156116a2573d6000803e3d6000fd5b506116e992505073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690503384612b8a565b604051829033907f2bda316674f8d73d289689d7a3acdf8e353b7a142fb5a68ac2aa475104039c1890600090a35050565b60045460208901359067ffffffffffffffff6801000000000000000090910481169082161161179d57600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff80851693820193909352680100000000000000009091049091166024820152604401610bf4565b6117a8888883612be0565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a3591821461185a5780516040517f93df584c000000000000000000000000000000000000000000000000000000008152600481019190915260248101839052604401610bf4565b467f0000000000000000000000000000000000000000000000000000000000000000146118db576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610bf4565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a1602081015161192d90600161423e565b60ff168714611968576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8685146119a1576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156119e4576119e4614156565b60028111156119f5576119f5614156565b9052509050600281602001516002811115611a1257611a12614156565b148015611a5957506007816000015160ff1681548110611a3457611a34614127565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611a8f576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000611a9d86602061410a565b611aa889602061410a565b611ab48c610144614257565b611abe9190614257565b611ac89190614257565b9050368114611b0c576040517f8e1192e100000000000000000000000000000000000000000000000000000000815260048101829052366024820152604401610bf4565b5060008a8a604051611b1f92919061426a565b604051908190038120611b36918e9060200161427a565b604051602081830303815290604052805190602001209050611b56613887565b8860005b81811015611d7a5760006001858a8460208110611b7957611b79614127565b611b8691901a601b61423e565b8f8f86818110611b9857611b98614127565b905060200201358e8e87818110611bb157611bb1614127565b9050602002013560405160008152602001604052604051611bee949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611c10573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff8082168652939750919550929392840191610100909104166002811115611c9357611c93614156565b6002811115611ca457611ca4614156565b9052509050600181602001516002811115611cc157611cc1614156565b14611cf8576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f8110611d0f57611d0f614127565b602002015115611d4b576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f8110611d6657611d66614127565b911515602090920201525050600101611b5a565b505050505050505050505050505050565b611d93612177565b60005b81811015611dc657611dbe838383818110611db357611db3614127565b905060a002016121fa565b600101611d96565b505050565b600c5473ffffffffffffffffffffffffffffffffffffffff163314611e1c576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008173ffffffffffffffffffffffffffffffffffffffff168360405160006040518083038185875af1925050503d8060008114611e76576040519150601f19603f3d011682016040523d82523d6000602084013e611e7b565b606091505b5050905080611eb6576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805184815273ffffffffffffffffffffffffffffffffffffffff841660208201527f6b84d241b711af111ecfa0e518239e6ca212da442a76548fe8a1f4e77518256a910160405180910390a1505050565b611f11612177565b61082581612d8f565b600c5473ffffffffffffffffffffffffffffffffffffffff163314611f6b576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109d385858567ffffffffffffffff86868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612e8492505050565b600a5460609060008167ffffffffffffffff811115611fd757611fd7613b66565b60405190808252806020026020018201604052801561204e57816020015b6040805160a0810182526000808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181611ff55790505b50905060005b82811015612170576000600a828154811061207157612071614127565b6000918252602080832060048304015460039092166008026101000a90910467ffffffffffffffff1680835260098252604092839020835160808082018652825473ffffffffffffffffffffffffffffffffffffffff9081168352600184015481168387019081526002909401548082168489019081527401000000000000000000000000000000000000000090910460ff1615156060808601918252895160a081018b5286518516815296518416988701989098529051909116968401969096529382018390529351151592810192909252855190935085908590811061215b5761215b614127565b60209081029190910101525050600101612054565b5092915050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146121f8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bf4565b565b61220a6080820160608301613ae6565b67ffffffffffffffff1660000361224d576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061225c60208301836139e4565b73ffffffffffffffffffffffffffffffffffffffff1614806122a35750600061228b60408301602084016139e4565b73ffffffffffffffffffffffffffffffffffffffff16145b806122d3575060006122bb60608301604084016139e4565b73ffffffffffffffffffffffffffffffffffffffff16145b1561230a576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060098161231f6080850160608601613ae6565b67ffffffffffffffff16815260208101919091526040016000206002015473ffffffffffffffffffffffffffffffffffffffff16036123ab57600a61236a6080830160608401613ae6565b8154600181018355600092835260209092206004830401805460039093166008026101000a67ffffffffffffffff8181021990941692909316929092021790555b6040805160808101909152806123c460208401846139e4565b73ffffffffffffffffffffffffffffffffffffffff1681526020018260200160208101906123f291906139e4565b73ffffffffffffffffffffffffffffffffffffffff16815260200161241d60608401604085016139e4565b73ffffffffffffffffffffffffffffffffffffffff16815260200161244860a084016080850161428e565b15159052600960006124606080850160608601613ae6565b67ffffffffffffffff16815260208082019290925260409081016000208351815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161783559385015160018301805491831691909516179093559083015160029091018054606094850151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090911692909316919091179190911790556125419060808301908301613ae6565b67ffffffffffffffff167fab9bd0e4888101232b8f09dae2952ff59a6eea4a19fbddf2a8ca7b23f0e4cb4061257c60408401602085016139e4565b61258c60608501604086016139e4565b61259960208601866139e4565b6125a960a087016080880161428e565b6040516125eb949392919073ffffffffffffffffffffffffffffffffffffffff9485168152928416602084015292166040820152901515606082015260800190565b60405180910390a250565b67ffffffffffffffff85166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff161515606082018190526126bb576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152602401610bf4565b602081015181516040517f38314bb200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909216916338314bb29161271791309089906004016142ab565b6020604051808303816000875af1925050508015612770575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261276d918101906142ed565b60015b6127f9573d80801561279e576040519150601f19603f3d011682016040523d82523d6000602084013e6127a3565b606091505b508667ffffffffffffffff168367ffffffffffffffff167fa481d91c3f9574c23ee84fef85246354b760a0527a535d6382354e4684703ce387846040516127eb92919061430a565b60405180910390a350612866565b80156128115761280c86848988886131fd565b61285f565b8667ffffffffffffffff168367ffffffffffffffff167f8d3121fe961b40270f336aa75feb1213f1c979a33993311c60da4dd0f24526cf876040516128569190613929565b60405180910390a35b50506109d3565b61287385838887876131fd565b505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526129579085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526133e4565b50505050565b8015806129fd57506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa1580156129d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129fb91906140c2565b155b612a89576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608401610bf4565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052611dc69084907f095ea7b300000000000000000000000000000000000000000000000000000000906064016128d5565b6000808a8a8a8a8a8a8a8a8a604051602001612b039998979695949392919061432f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052611dc69084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064016128d5565b6000612bee838501856144cc565b8051516020820151519192509081158015612c07575080155b15612c3d576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82811015612ce157612cd984600001518281518110612c6157612c61614127565b60200260200101516040015185600001518381518110612c8357612c83614127565b60200260200101516000015186600001518481518110612ca557612ca5614127565b6020026020010151602001518888600001518681518110612cc857612cc8614127565b602002602001015160600151612e84565b600101612c40565b5060005b81811015612d8657612d7e84602001518281518110612d0657612d06614127565b60200260200101516020015185602001518381518110612d2857612d28614127565b60200260200101516000015186602001518481518110612d4a57612d4a614127565b60200260200101516060015187602001518581518110612d6c57612d6c614127565b602002602001015160400151896125f6565b600101612ce5565b50505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612e0e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bf4565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612e8e61070c565b60085490915080821080612eaa575085612ea88284614640565b105b15612ef2576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018790526024810183905260448101829052606401610bf4565b67ffffffffffffffff87166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052612fb7576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff89166004820152602401610bf4565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810189905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b15801561302357600080fd5b505af1158015613037573d6000803e3d6000fd5b5050506020820151613082915073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016908961295d565b6020810151604080830151835191517fa71d98b700000000000000000000000000000000000000000000000000000000815260009373ffffffffffffffffffffffffffffffffffffffff169263a71d98b7928b92613109927f000000000000000000000000000000000000000000000000000000000000000092918f908d90600401614653565b60006040518083038185885af1158015613127573d6000803e3d6000fd5b50505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261316e919081019061469a565b90508867ffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168767ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a885600001518c8a876040516131ea9493929190614708565b60405180910390a4505050505050505050565b8015613285577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0866040518263ffffffff1660e01b81526004016000604051808303818588803b15801561326b57600080fd5b505af115801561327f573d6000803e3d6000fd5b50505050505b600b546132cc9073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169116876134f0565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810187905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b15801561333857600080fd5b505af115801561334c573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168367ffffffffffffffff168567ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a8308987604051806020016040528060008152506040516133d59493929190614708565b60405180910390a45050505050565b6000613446826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166135ee9092919063ffffffff16565b805190915015611dc6578080602001905181019061346491906142ed565b611dc6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610bf4565b6040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff8381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa158015613567573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061358b91906140c2565b6135959190614257565b60405173ffffffffffffffffffffffffffffffffffffffff85166024820152604481018290529091506129579085907f095ea7b300000000000000000000000000000000000000000000000000000000906064016128d5565b60606135fd8484600085613605565b949350505050565b606082471015613697576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610bf4565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516136c0919061474f565b60006040518083038185875af1925050503d80600081146136fd576040519150601f19603f3d011682016040523d82523d6000602084013e613702565b606091505b50915091506137138783838761371e565b979650505050505050565b606083156137b45782516000036137ad5773ffffffffffffffffffffffffffffffffffffffff85163b6137ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610bf4565b50816135fd565b6135fd83838151156137c95781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bf49190613929565b828054828255906000526020600020908101928215613877579160200282015b8281111561387757825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90911617825560209092019160019091019061381d565b506138839291506138a6565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b8082111561388357600081556001016138a7565b60005b838110156138d65781810151838201526020016138be565b50506000910152565b600081518084526138f78160208601602086016138bb565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061393c60208301846138df565b9392505050565b60006020828403121561395557600080fd5b5035919050565b600060a0828403121561396e57600080fd5b50919050565b6020808252825182820181905260009190848201906040850190845b818110156139b657835167ffffffffffffffff1683529284019291840191600101613990565b50909695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461082557600080fd5b6000602082840312156139f657600080fd5b813561393c816139c2565b803567ffffffffffffffff81168114613a1957600080fd5b919050565b801515811461082557600080fd5b60008083601f840112613a3e57600080fd5b50813567ffffffffffffffff811115613a5657600080fd5b602083019150836020828501011115613a6e57600080fd5b9250929050565b600080600080600060808688031215613a8d57600080fd5b613a9686613a01565b9450602086013593506040860135613aad81613a1e565b9250606086013567ffffffffffffffff811115613ac957600080fd5b613ad588828901613a2c565b969995985093965092949392505050565b600060208284031215613af857600080fd5b61393c82613a01565b60008151808452602080850194506020840160005b83811015613b4857815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613b16565b509495945050505050565b60208152600061393c6020830184613b01565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613bb857613bb8613b66565b60405290565b6040805190810167ffffffffffffffff81118282101715613bb857613bb8613b66565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613c2857613c28613b66565b604052919050565b600067ffffffffffffffff821115613c4a57613c4a613b66565b5060051b60200190565b600082601f830112613c6557600080fd5b81356020613c7a613c7583613c30565b613be1565b8083825260208201915060208460051b870101935086841115613c9c57600080fd5b602086015b84811015613cc1578035613cb4816139c2565b8352918301918301613ca1565b509695505050505050565b803560ff81168114613a1957600080fd5b600067ffffffffffffffff821115613cf757613cf7613b66565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613d3457600080fd5b8135613d42613c7582613cdd565b818152846020838601011115613d5757600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215613d8d57600080fd5b863567ffffffffffffffff80821115613da557600080fd5b613db18a838b01613c54565b97506020890135915080821115613dc757600080fd5b613dd38a838b01613c54565b9650613de160408a01613ccc565b95506060890135915080821115613df757600080fd5b613e038a838b01613d23565b9450613e1160808a01613a01565b935060a0890135915080821115613e2757600080fd5b50613e3489828a01613d23565b9150509295509295509295565b60008083601f840112613e5357600080fd5b50813567ffffffffffffffff811115613e6b57600080fd5b6020830191508360208260051b8501011115613a6e57600080fd5b60008060008060008060008060e0898b031215613ea257600080fd5b606089018a811115613eb357600080fd5b8998503567ffffffffffffffff80821115613ecd57600080fd5b613ed98c838d01613a2c565b909950975060808b0135915080821115613ef257600080fd5b613efe8c838d01613e41565b909750955060a08b0135915080821115613f1757600080fd5b50613f248b828c01613e41565b999c989b50969995989497949560c00135949350505050565b60008060208385031215613f5057600080fd5b823567ffffffffffffffff80821115613f6857600080fd5b818501915085601f830112613f7c57600080fd5b813581811115613f8b57600080fd5b86602060a083028501011115613fa057600080fd5b60209290920196919550909350505050565b60008060408385031215613fc557600080fd5b823591506020830135613fd7816139c2565b809150509250929050565b600080600080600060808688031215613ffa57600080fd5b61400386613a01565b94506020860135935060408601359250606086013567ffffffffffffffff811115613ac957600080fd5b602080825282518282018190526000919060409081850190868401855b828110156140b5578151805173ffffffffffffffffffffffffffffffffffffffff90811686528782015181168887015286820151168686015260608082015167ffffffffffffffff169086015260809081015115159085015260a0909301929085019060010161404a565b5091979650505050505050565b6000602082840312156140d457600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417614121576141216140db565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff80831681810361419e5761419e6140db565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526141d88184018a613b01565b905082810360808401526141ec8189613b01565b905060ff871660a084015282810360c084015261420981876138df565b905067ffffffffffffffff851660e084015282810361010084015261422e81856138df565b9c9b505050505050505050505050565b60ff8181168382160190811115614121576141216140db565b80820180821115614121576141216140db565b8183823760009101908152919050565b828152606082602083013760800192915050565b6000602082840312156142a057600080fd5b813561393c81613a1e565b600073ffffffffffffffffffffffffffffffffffffffff8086168352808516602084015250606060408301526142e460608301846138df565b95945050505050565b6000602082840312156142ff57600080fd5b815161393c81613a1e565b60408152600061431d60408301856138df565b82810360208401526142e481856138df565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526143768285018b613b01565b9150838203608085015261438a828a613b01565b915060ff881660a085015283820360c08501526143a782886138df565b90861660e0850152838103610100850152905061422e81856138df565b600082601f8301126143d557600080fd5b813560206143e5613c7583613c30565b82815260059290921b8401810191818101908684111561440457600080fd5b8286015b84811015613cc157803567ffffffffffffffff808211156144295760008081fd5b81890191506080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156144625760008081fd5b61446a613b95565b878401358152604061447d818601613a01565b8983015260608086013561449081613a1e565b838301529285013592848411156144a957600091508182fd5b6144b78e8b86890101613d23565b90830152508652505050918301918301614408565b600060208083850312156144df57600080fd5b823567ffffffffffffffff808211156144f757600080fd5b908401906040828703121561450b57600080fd5b614513613bbe565b82358281111561452257600080fd5b8301601f8101881361453357600080fd5b8035614541613c7582613c30565b81815260059190911b8201860190868101908a83111561456057600080fd5b8784015b8381101561460c5780358781111561457b57600080fd5b85016080818e037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00112156145af57600080fd5b6145b7613b95565b8a820135815260408201358b8201526145d260608301613a01565b60408201526080820135898111156145ea5760008081fd5b6145f88f8d83860101613d23565b606083015250845250918801918801614564565b508452505050828401358281111561462357600080fd5b61462f888286016143c4565b948201949094529695505050505050565b81810381811115614121576141216140db565b600073ffffffffffffffffffffffffffffffffffffffff8088168352808716602084015280861660408401525083606083015260a0608083015261371360a08301846138df565b6000602082840312156146ac57600080fd5b815167ffffffffffffffff8111156146c357600080fd5b8201601f810184136146d457600080fd5b80516146e2613c7582613cdd565b8181528560208385010111156146f757600080fd5b6142e48260208301602086016138bb565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260806040820152600061473d60808301856138df565b828103606084015261371381856138df565b600082516147618184602087016138bb565b919091019291505056fea164736f6c6343000818000a",
}

var LiquidityManagerABI = LiquidityManagerMetaData.ABI

var LiquidityManagerBin = LiquidityManagerMetaData.Bin

func DeployLiquidityManager(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localChainSelector uint64, localLiquidityContainer common.Address, minimumLiquidity *big.Int, finance common.Address) (common.Address, *types.Transaction, *LiquidityManager, error) {
	parsed, err := LiquidityManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidityManagerBin), backend, token, localChainSelector, localLiquidityContainer, minimumLiquidity, finance)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LiquidityManager{address: address, abi: *parsed, LiquidityManagerCaller: LiquidityManagerCaller{contract: contract}, LiquidityManagerTransactor: LiquidityManagerTransactor{contract: contract}, LiquidityManagerFilterer: LiquidityManagerFilterer{contract: contract}}, nil
}

type LiquidityManager struct {
	address common.Address
	abi     abi.ABI
	LiquidityManagerCaller
	LiquidityManagerTransactor
	LiquidityManagerFilterer
}

type LiquidityManagerCaller struct {
	contract *bind.BoundContract
}

type LiquidityManagerTransactor struct {
	contract *bind.BoundContract
}

type LiquidityManagerFilterer struct {
	contract *bind.BoundContract
}

type LiquidityManagerSession struct {
	Contract     *LiquidityManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LiquidityManagerCallerSession struct {
	Contract *LiquidityManagerCaller
	CallOpts bind.CallOpts
}

type LiquidityManagerTransactorSession struct {
	Contract     *LiquidityManagerTransactor
	TransactOpts bind.TransactOpts
}

type LiquidityManagerRaw struct {
	Contract *LiquidityManager
}

type LiquidityManagerCallerRaw struct {
	Contract *LiquidityManagerCaller
}

type LiquidityManagerTransactorRaw struct {
	Contract *LiquidityManagerTransactor
}

func NewLiquidityManager(address common.Address, backend bind.ContractBackend) (*LiquidityManager, error) {
	abi, err := abi.JSON(strings.NewReader(LiquidityManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLiquidityManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LiquidityManager{address: address, abi: abi, LiquidityManagerCaller: LiquidityManagerCaller{contract: contract}, LiquidityManagerTransactor: LiquidityManagerTransactor{contract: contract}, LiquidityManagerFilterer: LiquidityManagerFilterer{contract: contract}}, nil
}

func NewLiquidityManagerCaller(address common.Address, caller bind.ContractCaller) (*LiquidityManagerCaller, error) {
	contract, err := bindLiquidityManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerCaller{contract: contract}, nil
}

func NewLiquidityManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidityManagerTransactor, error) {
	contract, err := bindLiquidityManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerTransactor{contract: contract}, nil
}

func NewLiquidityManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidityManagerFilterer, error) {
	contract, err := bindLiquidityManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerFilterer{contract: contract}, nil
}

func bindLiquidityManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LiquidityManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LiquidityManager *LiquidityManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidityManager.Contract.LiquidityManagerCaller.contract.Call(opts, result, method, params...)
}

func (_LiquidityManager *LiquidityManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityManager.Contract.LiquidityManagerTransactor.contract.Transfer(opts)
}

func (_LiquidityManager *LiquidityManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidityManager.Contract.LiquidityManagerTransactor.contract.Transact(opts, method, params...)
}

func (_LiquidityManager *LiquidityManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidityManager.Contract.contract.Call(opts, result, method, params...)
}

func (_LiquidityManager *LiquidityManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityManager.Contract.contract.Transfer(opts)
}

func (_LiquidityManager *LiquidityManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidityManager.Contract.contract.Transact(opts, method, params...)
}

func (_LiquidityManager *LiquidityManagerCaller) GetAllCrossChainRebalancers(opts *bind.CallOpts) ([]ILiquidityManagerCrossChainRebalancerArgs, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getAllCrossChainRebalancers")

	if err != nil {
		return *new([]ILiquidityManagerCrossChainRebalancerArgs), err
	}

	out0 := *abi.ConvertType(out[0], new([]ILiquidityManagerCrossChainRebalancerArgs)).(*[]ILiquidityManagerCrossChainRebalancerArgs)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetAllCrossChainRebalancers() ([]ILiquidityManagerCrossChainRebalancerArgs, error) {
	return _LiquidityManager.Contract.GetAllCrossChainRebalancers(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetAllCrossChainRebalancers() ([]ILiquidityManagerCrossChainRebalancerArgs, error) {
	return _LiquidityManager.Contract.GetAllCrossChainRebalancers(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) GetCrossChainRebalancer(opts *bind.CallOpts, chainSelector uint64) (LiquidityManagerCrossChainRebalancer, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getCrossChainRebalancer", chainSelector)

	if err != nil {
		return *new(LiquidityManagerCrossChainRebalancer), err
	}

	out0 := *abi.ConvertType(out[0], new(LiquidityManagerCrossChainRebalancer)).(*LiquidityManagerCrossChainRebalancer)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetCrossChainRebalancer(chainSelector uint64) (LiquidityManagerCrossChainRebalancer, error) {
	return _LiquidityManager.Contract.GetCrossChainRebalancer(&_LiquidityManager.CallOpts, chainSelector)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetCrossChainRebalancer(chainSelector uint64) (LiquidityManagerCrossChainRebalancer, error) {
	return _LiquidityManager.Contract.GetCrossChainRebalancer(&_LiquidityManager.CallOpts, chainSelector)
}

func (_LiquidityManager *LiquidityManagerCaller) GetFinanceRole(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getFinanceRole")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetFinanceRole() (common.Address, error) {
	return _LiquidityManager.Contract.GetFinanceRole(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetFinanceRole() (common.Address, error) {
	return _LiquidityManager.Contract.GetFinanceRole(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) GetLiquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getLiquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetLiquidity() (*big.Int, error) {
	return _LiquidityManager.Contract.GetLiquidity(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetLiquidity() (*big.Int, error) {
	return _LiquidityManager.Contract.GetLiquidity(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) GetLocalLiquidityContainer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getLocalLiquidityContainer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetLocalLiquidityContainer() (common.Address, error) {
	return _LiquidityManager.Contract.GetLocalLiquidityContainer(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetLocalLiquidityContainer() (common.Address, error) {
	return _LiquidityManager.Contract.GetLocalLiquidityContainer(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) GetMinimumLiquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getMinimumLiquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetMinimumLiquidity() (*big.Int, error) {
	return _LiquidityManager.Contract.GetMinimumLiquidity(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetMinimumLiquidity() (*big.Int, error) {
	return _LiquidityManager.Contract.GetMinimumLiquidity(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) GetSupportedDestChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getSupportedDestChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetSupportedDestChains() ([]uint64, error) {
	return _LiquidityManager.Contract.GetSupportedDestChains(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetSupportedDestChains() ([]uint64, error) {
	return _LiquidityManager.Contract.GetSupportedDestChains(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) GetTransmitters() ([]common.Address, error) {
	return _LiquidityManager.Contract.GetTransmitters(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) GetTransmitters() ([]common.Address, error) {
	return _LiquidityManager.Contract.GetTransmitters(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) ILocalToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "i_localToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) ILocalToken() (common.Address, error) {
	return _LiquidityManager.Contract.ILocalToken(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) ILocalToken() (common.Address, error) {
	return _LiquidityManager.Contract.ILocalToken(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_LiquidityManager *LiquidityManagerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _LiquidityManager.Contract.LatestConfigDetails(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _LiquidityManager.Contract.LatestConfigDetails(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) LatestSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "latestSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) LatestSequenceNumber() (uint64, error) {
	return _LiquidityManager.Contract.LatestSequenceNumber(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) LatestSequenceNumber() (uint64, error) {
	return _LiquidityManager.Contract.LatestSequenceNumber(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) Owner() (common.Address, error) {
	return _LiquidityManager.Contract.Owner(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) Owner() (common.Address, error) {
	return _LiquidityManager.Contract.Owner(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LiquidityManager *LiquidityManagerSession) TypeAndVersion() (string, error) {
	return _LiquidityManager.Contract.TypeAndVersion(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) TypeAndVersion() (string, error) {
	return _LiquidityManager.Contract.TypeAndVersion(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "acceptOwnership")
}

func (_LiquidityManager *LiquidityManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _LiquidityManager.Contract.AcceptOwnership(&_LiquidityManager.TransactOpts)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LiquidityManager.Contract.AcceptOwnership(&_LiquidityManager.TransactOpts)
}

func (_LiquidityManager *LiquidityManagerTransactor) AddLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "addLiquidity", amount)
}

func (_LiquidityManager *LiquidityManagerSession) AddLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.Contract.AddLiquidity(&_LiquidityManager.TransactOpts, amount)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) AddLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.Contract.AddLiquidity(&_LiquidityManager.TransactOpts, amount)
}

func (_LiquidityManager *LiquidityManagerTransactor) RebalanceLiquidity(opts *bind.TransactOpts, chainSelector uint64, amount *big.Int, nativeBridgeFee *big.Int, bridgeSpecificPayload []byte) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "rebalanceLiquidity", chainSelector, amount, nativeBridgeFee, bridgeSpecificPayload)
}

func (_LiquidityManager *LiquidityManagerSession) RebalanceLiquidity(chainSelector uint64, amount *big.Int, nativeBridgeFee *big.Int, bridgeSpecificPayload []byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.RebalanceLiquidity(&_LiquidityManager.TransactOpts, chainSelector, amount, nativeBridgeFee, bridgeSpecificPayload)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) RebalanceLiquidity(chainSelector uint64, amount *big.Int, nativeBridgeFee *big.Int, bridgeSpecificPayload []byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.RebalanceLiquidity(&_LiquidityManager.TransactOpts, chainSelector, amount, nativeBridgeFee, bridgeSpecificPayload)
}

func (_LiquidityManager *LiquidityManagerTransactor) ReceiveLiquidity(opts *bind.TransactOpts, remoteChainSelector uint64, amount *big.Int, shouldWrapNative bool, bridgeSpecificPayload []byte) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "receiveLiquidity", remoteChainSelector, amount, shouldWrapNative, bridgeSpecificPayload)
}

func (_LiquidityManager *LiquidityManagerSession) ReceiveLiquidity(remoteChainSelector uint64, amount *big.Int, shouldWrapNative bool, bridgeSpecificPayload []byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.ReceiveLiquidity(&_LiquidityManager.TransactOpts, remoteChainSelector, amount, shouldWrapNative, bridgeSpecificPayload)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) ReceiveLiquidity(remoteChainSelector uint64, amount *big.Int, shouldWrapNative bool, bridgeSpecificPayload []byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.ReceiveLiquidity(&_LiquidityManager.TransactOpts, remoteChainSelector, amount, shouldWrapNative, bridgeSpecificPayload)
}

func (_LiquidityManager *LiquidityManagerTransactor) RemoveLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "removeLiquidity", amount)
}

func (_LiquidityManager *LiquidityManagerSession) RemoveLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.Contract.RemoveLiquidity(&_LiquidityManager.TransactOpts, amount)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) RemoveLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.Contract.RemoveLiquidity(&_LiquidityManager.TransactOpts, amount)
}

func (_LiquidityManager *LiquidityManagerTransactor) SetCrossChainRebalancer(opts *bind.TransactOpts, crossChainLiqManager ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "setCrossChainRebalancer", crossChainLiqManager)
}

func (_LiquidityManager *LiquidityManagerSession) SetCrossChainRebalancer(crossChainLiqManager ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetCrossChainRebalancer(&_LiquidityManager.TransactOpts, crossChainLiqManager)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) SetCrossChainRebalancer(crossChainLiqManager ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetCrossChainRebalancer(&_LiquidityManager.TransactOpts, crossChainLiqManager)
}

func (_LiquidityManager *LiquidityManagerTransactor) SetCrossChainRebalancers(opts *bind.TransactOpts, crossChainRebalancers []ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "setCrossChainRebalancers", crossChainRebalancers)
}

func (_LiquidityManager *LiquidityManagerSession) SetCrossChainRebalancers(crossChainRebalancers []ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetCrossChainRebalancers(&_LiquidityManager.TransactOpts, crossChainRebalancers)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) SetCrossChainRebalancers(crossChainRebalancers []ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetCrossChainRebalancers(&_LiquidityManager.TransactOpts, crossChainRebalancers)
}

func (_LiquidityManager *LiquidityManagerTransactor) SetFinanceRole(opts *bind.TransactOpts, finance common.Address) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "setFinanceRole", finance)
}

func (_LiquidityManager *LiquidityManagerSession) SetFinanceRole(finance common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetFinanceRole(&_LiquidityManager.TransactOpts, finance)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) SetFinanceRole(finance common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetFinanceRole(&_LiquidityManager.TransactOpts, finance)
}

func (_LiquidityManager *LiquidityManagerTransactor) SetLocalLiquidityContainer(opts *bind.TransactOpts, localLiquidityContainer common.Address) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "setLocalLiquidityContainer", localLiquidityContainer)
}

func (_LiquidityManager *LiquidityManagerSession) SetLocalLiquidityContainer(localLiquidityContainer common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetLocalLiquidityContainer(&_LiquidityManager.TransactOpts, localLiquidityContainer)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) SetLocalLiquidityContainer(localLiquidityContainer common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetLocalLiquidityContainer(&_LiquidityManager.TransactOpts, localLiquidityContainer)
}

func (_LiquidityManager *LiquidityManagerTransactor) SetMinimumLiquidity(opts *bind.TransactOpts, minimumLiquidity *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "setMinimumLiquidity", minimumLiquidity)
}

func (_LiquidityManager *LiquidityManagerSession) SetMinimumLiquidity(minimumLiquidity *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetMinimumLiquidity(&_LiquidityManager.TransactOpts, minimumLiquidity)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) SetMinimumLiquidity(minimumLiquidity *big.Int) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetMinimumLiquidity(&_LiquidityManager.TransactOpts, minimumLiquidity)
}

func (_LiquidityManager *LiquidityManagerTransactor) SetOCR3Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "setOCR3Config", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_LiquidityManager *LiquidityManagerSession) SetOCR3Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetOCR3Config(&_LiquidityManager.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) SetOCR3Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.SetOCR3Config(&_LiquidityManager.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_LiquidityManager *LiquidityManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "transferOwnership", to)
}

func (_LiquidityManager *LiquidityManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.TransferOwnership(&_LiquidityManager.TransactOpts, to)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.TransferOwnership(&_LiquidityManager.TransactOpts, to)
}

func (_LiquidityManager *LiquidityManagerTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_LiquidityManager *LiquidityManagerSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.Transmit(&_LiquidityManager.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _LiquidityManager.Contract.Transmit(&_LiquidityManager.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_LiquidityManager *LiquidityManagerTransactor) WithdrawNative(opts *bind.TransactOpts, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "withdrawNative", amount, destination)
}

func (_LiquidityManager *LiquidityManagerSession) WithdrawNative(amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.WithdrawNative(&_LiquidityManager.TransactOpts, amount, destination)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) WithdrawNative(amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.WithdrawNative(&_LiquidityManager.TransactOpts, amount, destination)
}

func (_LiquidityManager *LiquidityManagerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityManager.contract.RawTransact(opts, nil)
}

func (_LiquidityManager *LiquidityManagerSession) Receive() (*types.Transaction, error) {
	return _LiquidityManager.Contract.Receive(&_LiquidityManager.TransactOpts)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) Receive() (*types.Transaction, error) {
	return _LiquidityManager.Contract.Receive(&_LiquidityManager.TransactOpts)
}

type LiquidityManagerConfigSetIterator struct {
	Event *LiquidityManagerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerConfigSet)
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
		it.Event = new(LiquidityManagerConfigSet)
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

func (it *LiquidityManagerConfigSetIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*LiquidityManagerConfigSetIterator, error) {

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerConfigSetIterator{contract: _LiquidityManager.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerConfigSet) (event.Subscription, error) {

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerConfigSet)
				if err := _LiquidityManager.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseConfigSet(log types.Log) (*LiquidityManagerConfigSet, error) {
	event := new(LiquidityManagerConfigSet)
	if err := _LiquidityManager.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerCrossChainRebalancerSetIterator struct {
	Event *LiquidityManagerCrossChainRebalancerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerCrossChainRebalancerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerCrossChainRebalancerSet)
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
		it.Event = new(LiquidityManagerCrossChainRebalancerSet)
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

func (it *LiquidityManagerCrossChainRebalancerSetIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerCrossChainRebalancerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerCrossChainRebalancerSet struct {
	RemoteChainSelector uint64
	LocalBridge         common.Address
	RemoteToken         common.Address
	RemoteRebalancer    common.Address
	Enabled             bool
	Raw                 types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterCrossChainRebalancerSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LiquidityManagerCrossChainRebalancerSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "CrossChainRebalancerSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerCrossChainRebalancerSetIterator{contract: _LiquidityManager.contract, event: "CrossChainRebalancerSet", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchCrossChainRebalancerSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerCrossChainRebalancerSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "CrossChainRebalancerSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerCrossChainRebalancerSet)
				if err := _LiquidityManager.contract.UnpackLog(event, "CrossChainRebalancerSet", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseCrossChainRebalancerSet(log types.Log) (*LiquidityManagerCrossChainRebalancerSet, error) {
	event := new(LiquidityManagerCrossChainRebalancerSet)
	if err := _LiquidityManager.contract.UnpackLog(event, "CrossChainRebalancerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerFinalizationFailedIterator struct {
	Event *LiquidityManagerFinalizationFailed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerFinalizationFailedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerFinalizationFailed)
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
		it.Event = new(LiquidityManagerFinalizationFailed)
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

func (it *LiquidityManagerFinalizationFailedIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerFinalizationFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerFinalizationFailed struct {
	OcrSeqNum           uint64
	RemoteChainSelector uint64
	BridgeSpecificData  []byte
	Reason              []byte
	Raw                 types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterFinalizationFailed(opts *bind.FilterOpts, ocrSeqNum []uint64, remoteChainSelector []uint64) (*LiquidityManagerFinalizationFailedIterator, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "FinalizationFailed", ocrSeqNumRule, remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerFinalizationFailedIterator{contract: _LiquidityManager.contract, event: "FinalizationFailed", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchFinalizationFailed(opts *bind.WatchOpts, sink chan<- *LiquidityManagerFinalizationFailed, ocrSeqNum []uint64, remoteChainSelector []uint64) (event.Subscription, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "FinalizationFailed", ocrSeqNumRule, remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerFinalizationFailed)
				if err := _LiquidityManager.contract.UnpackLog(event, "FinalizationFailed", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseFinalizationFailed(log types.Log) (*LiquidityManagerFinalizationFailed, error) {
	event := new(LiquidityManagerFinalizationFailed)
	if err := _LiquidityManager.contract.UnpackLog(event, "FinalizationFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerFinalizationStepCompletedIterator struct {
	Event *LiquidityManagerFinalizationStepCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerFinalizationStepCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerFinalizationStepCompleted)
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
		it.Event = new(LiquidityManagerFinalizationStepCompleted)
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

func (it *LiquidityManagerFinalizationStepCompletedIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerFinalizationStepCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerFinalizationStepCompleted struct {
	OcrSeqNum           uint64
	RemoteChainSelector uint64
	BridgeSpecificData  []byte
	Raw                 types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterFinalizationStepCompleted(opts *bind.FilterOpts, ocrSeqNum []uint64, remoteChainSelector []uint64) (*LiquidityManagerFinalizationStepCompletedIterator, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "FinalizationStepCompleted", ocrSeqNumRule, remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerFinalizationStepCompletedIterator{contract: _LiquidityManager.contract, event: "FinalizationStepCompleted", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchFinalizationStepCompleted(opts *bind.WatchOpts, sink chan<- *LiquidityManagerFinalizationStepCompleted, ocrSeqNum []uint64, remoteChainSelector []uint64) (event.Subscription, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "FinalizationStepCompleted", ocrSeqNumRule, remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerFinalizationStepCompleted)
				if err := _LiquidityManager.contract.UnpackLog(event, "FinalizationStepCompleted", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseFinalizationStepCompleted(log types.Log) (*LiquidityManagerFinalizationStepCompleted, error) {
	event := new(LiquidityManagerFinalizationStepCompleted)
	if err := _LiquidityManager.contract.UnpackLog(event, "FinalizationStepCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerFinanceRoleSetIterator struct {
	Event *LiquidityManagerFinanceRoleSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerFinanceRoleSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerFinanceRoleSet)
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
		it.Event = new(LiquidityManagerFinanceRoleSet)
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

func (it *LiquidityManagerFinanceRoleSetIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerFinanceRoleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerFinanceRoleSet struct {
	FinanceRole common.Address
	Raw         types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterFinanceRoleSet(opts *bind.FilterOpts) (*LiquidityManagerFinanceRoleSetIterator, error) {

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "FinanceRoleSet")
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerFinanceRoleSetIterator{contract: _LiquidityManager.contract, event: "FinanceRoleSet", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchFinanceRoleSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerFinanceRoleSet) (event.Subscription, error) {

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "FinanceRoleSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerFinanceRoleSet)
				if err := _LiquidityManager.contract.UnpackLog(event, "FinanceRoleSet", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseFinanceRoleSet(log types.Log) (*LiquidityManagerFinanceRoleSet, error) {
	event := new(LiquidityManagerFinanceRoleSet)
	if err := _LiquidityManager.contract.UnpackLog(event, "FinanceRoleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerLiquidityAddedToContainerIterator struct {
	Event *LiquidityManagerLiquidityAddedToContainer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerLiquidityAddedToContainerIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerLiquidityAddedToContainer)
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
		it.Event = new(LiquidityManagerLiquidityAddedToContainer)
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

func (it *LiquidityManagerLiquidityAddedToContainerIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerLiquidityAddedToContainerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerLiquidityAddedToContainer struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterLiquidityAddedToContainer(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LiquidityManagerLiquidityAddedToContainerIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "LiquidityAddedToContainer", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerLiquidityAddedToContainerIterator{contract: _LiquidityManager.contract, event: "LiquidityAddedToContainer", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchLiquidityAddedToContainer(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityAddedToContainer, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "LiquidityAddedToContainer", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerLiquidityAddedToContainer)
				if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityAddedToContainer", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseLiquidityAddedToContainer(log types.Log) (*LiquidityManagerLiquidityAddedToContainer, error) {
	event := new(LiquidityManagerLiquidityAddedToContainer)
	if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityAddedToContainer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerLiquidityContainerSetIterator struct {
	Event *LiquidityManagerLiquidityContainerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerLiquidityContainerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerLiquidityContainerSet)
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
		it.Event = new(LiquidityManagerLiquidityContainerSet)
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

func (it *LiquidityManagerLiquidityContainerSetIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerLiquidityContainerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerLiquidityContainerSet struct {
	NewLiquidityContainer common.Address
	Raw                   types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterLiquidityContainerSet(opts *bind.FilterOpts, newLiquidityContainer []common.Address) (*LiquidityManagerLiquidityContainerSetIterator, error) {

	var newLiquidityContainerRule []interface{}
	for _, newLiquidityContainerItem := range newLiquidityContainer {
		newLiquidityContainerRule = append(newLiquidityContainerRule, newLiquidityContainerItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "LiquidityContainerSet", newLiquidityContainerRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerLiquidityContainerSetIterator{contract: _LiquidityManager.contract, event: "LiquidityContainerSet", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchLiquidityContainerSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityContainerSet, newLiquidityContainer []common.Address) (event.Subscription, error) {

	var newLiquidityContainerRule []interface{}
	for _, newLiquidityContainerItem := range newLiquidityContainer {
		newLiquidityContainerRule = append(newLiquidityContainerRule, newLiquidityContainerItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "LiquidityContainerSet", newLiquidityContainerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerLiquidityContainerSet)
				if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityContainerSet", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseLiquidityContainerSet(log types.Log) (*LiquidityManagerLiquidityContainerSet, error) {
	event := new(LiquidityManagerLiquidityContainerSet)
	if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityContainerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerLiquidityRemovedFromContainerIterator struct {
	Event *LiquidityManagerLiquidityRemovedFromContainer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerLiquidityRemovedFromContainerIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerLiquidityRemovedFromContainer)
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
		it.Event = new(LiquidityManagerLiquidityRemovedFromContainer)
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

func (it *LiquidityManagerLiquidityRemovedFromContainerIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerLiquidityRemovedFromContainerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerLiquidityRemovedFromContainer struct {
	Remover common.Address
	Amount  *big.Int
	Raw     types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterLiquidityRemovedFromContainer(opts *bind.FilterOpts, remover []common.Address, amount []*big.Int) (*LiquidityManagerLiquidityRemovedFromContainerIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "LiquidityRemovedFromContainer", removerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerLiquidityRemovedFromContainerIterator{contract: _LiquidityManager.contract, event: "LiquidityRemovedFromContainer", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchLiquidityRemovedFromContainer(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityRemovedFromContainer, remover []common.Address, amount []*big.Int) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "LiquidityRemovedFromContainer", removerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerLiquidityRemovedFromContainer)
				if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityRemovedFromContainer", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseLiquidityRemovedFromContainer(log types.Log) (*LiquidityManagerLiquidityRemovedFromContainer, error) {
	event := new(LiquidityManagerLiquidityRemovedFromContainer)
	if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityRemovedFromContainer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerLiquidityTransferredIterator struct {
	Event *LiquidityManagerLiquidityTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerLiquidityTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerLiquidityTransferred)
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
		it.Event = new(LiquidityManagerLiquidityTransferred)
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

func (it *LiquidityManagerLiquidityTransferredIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerLiquidityTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerLiquidityTransferred struct {
	OcrSeqNum          uint64
	FromChainSelector  uint64
	ToChainSelector    uint64
	To                 common.Address
	Amount             *big.Int
	BridgeSpecificData []byte
	BridgeReturnData   []byte
	Raw                types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterLiquidityTransferred(opts *bind.FilterOpts, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (*LiquidityManagerLiquidityTransferredIterator, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var fromChainSelectorRule []interface{}
	for _, fromChainSelectorItem := range fromChainSelector {
		fromChainSelectorRule = append(fromChainSelectorRule, fromChainSelectorItem)
	}
	var toChainSelectorRule []interface{}
	for _, toChainSelectorItem := range toChainSelector {
		toChainSelectorRule = append(toChainSelectorRule, toChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "LiquidityTransferred", ocrSeqNumRule, fromChainSelectorRule, toChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerLiquidityTransferredIterator{contract: _LiquidityManager.contract, event: "LiquidityTransferred", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityTransferred, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (event.Subscription, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var fromChainSelectorRule []interface{}
	for _, fromChainSelectorItem := range fromChainSelector {
		fromChainSelectorRule = append(fromChainSelectorRule, fromChainSelectorItem)
	}
	var toChainSelectorRule []interface{}
	for _, toChainSelectorItem := range toChainSelector {
		toChainSelectorRule = append(toChainSelectorRule, toChainSelectorItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "LiquidityTransferred", ocrSeqNumRule, fromChainSelectorRule, toChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerLiquidityTransferred)
				if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseLiquidityTransferred(log types.Log) (*LiquidityManagerLiquidityTransferred, error) {
	event := new(LiquidityManagerLiquidityTransferred)
	if err := _LiquidityManager.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerMinimumLiquiditySetIterator struct {
	Event *LiquidityManagerMinimumLiquiditySet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerMinimumLiquiditySetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerMinimumLiquiditySet)
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
		it.Event = new(LiquidityManagerMinimumLiquiditySet)
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

func (it *LiquidityManagerMinimumLiquiditySetIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerMinimumLiquiditySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerMinimumLiquiditySet struct {
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterMinimumLiquiditySet(opts *bind.FilterOpts) (*LiquidityManagerMinimumLiquiditySetIterator, error) {

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "MinimumLiquiditySet")
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerMinimumLiquiditySetIterator{contract: _LiquidityManager.contract, event: "MinimumLiquiditySet", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchMinimumLiquiditySet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerMinimumLiquiditySet) (event.Subscription, error) {

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "MinimumLiquiditySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerMinimumLiquiditySet)
				if err := _LiquidityManager.contract.UnpackLog(event, "MinimumLiquiditySet", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseMinimumLiquiditySet(log types.Log) (*LiquidityManagerMinimumLiquiditySet, error) {
	event := new(LiquidityManagerMinimumLiquiditySet)
	if err := _LiquidityManager.contract.UnpackLog(event, "MinimumLiquiditySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerNativeDepositedIterator struct {
	Event *LiquidityManagerNativeDeposited

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerNativeDepositedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerNativeDeposited)
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
		it.Event = new(LiquidityManagerNativeDeposited)
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

func (it *LiquidityManagerNativeDepositedIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerNativeDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerNativeDeposited struct {
	Amount    *big.Int
	Depositor common.Address
	Raw       types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterNativeDeposited(opts *bind.FilterOpts) (*LiquidityManagerNativeDepositedIterator, error) {

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "NativeDeposited")
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerNativeDepositedIterator{contract: _LiquidityManager.contract, event: "NativeDeposited", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchNativeDeposited(opts *bind.WatchOpts, sink chan<- *LiquidityManagerNativeDeposited) (event.Subscription, error) {

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "NativeDeposited")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerNativeDeposited)
				if err := _LiquidityManager.contract.UnpackLog(event, "NativeDeposited", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseNativeDeposited(log types.Log) (*LiquidityManagerNativeDeposited, error) {
	event := new(LiquidityManagerNativeDeposited)
	if err := _LiquidityManager.contract.UnpackLog(event, "NativeDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerNativeWithdrawnIterator struct {
	Event *LiquidityManagerNativeWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerNativeWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerNativeWithdrawn)
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
		it.Event = new(LiquidityManagerNativeWithdrawn)
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

func (it *LiquidityManagerNativeWithdrawnIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerNativeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerNativeWithdrawn struct {
	Amount      *big.Int
	Destination common.Address
	Raw         types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterNativeWithdrawn(opts *bind.FilterOpts) (*LiquidityManagerNativeWithdrawnIterator, error) {

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "NativeWithdrawn")
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerNativeWithdrawnIterator{contract: _LiquidityManager.contract, event: "NativeWithdrawn", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchNativeWithdrawn(opts *bind.WatchOpts, sink chan<- *LiquidityManagerNativeWithdrawn) (event.Subscription, error) {

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "NativeWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerNativeWithdrawn)
				if err := _LiquidityManager.contract.UnpackLog(event, "NativeWithdrawn", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseNativeWithdrawn(log types.Log) (*LiquidityManagerNativeWithdrawn, error) {
	event := new(LiquidityManagerNativeWithdrawn)
	if err := _LiquidityManager.contract.UnpackLog(event, "NativeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerOwnershipTransferRequestedIterator struct {
	Event *LiquidityManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerOwnershipTransferRequested)
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
		it.Event = new(LiquidityManagerOwnershipTransferRequested)
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

func (it *LiquidityManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidityManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerOwnershipTransferRequestedIterator{contract: _LiquidityManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LiquidityManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerOwnershipTransferRequested)
				if err := _LiquidityManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*LiquidityManagerOwnershipTransferRequested, error) {
	event := new(LiquidityManagerOwnershipTransferRequested)
	if err := _LiquidityManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerOwnershipTransferredIterator struct {
	Event *LiquidityManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerOwnershipTransferred)
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
		it.Event = new(LiquidityManagerOwnershipTransferred)
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

func (it *LiquidityManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidityManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerOwnershipTransferredIterator{contract: _LiquidityManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LiquidityManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerOwnershipTransferred)
				if err := _LiquidityManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseOwnershipTransferred(log types.Log) (*LiquidityManagerOwnershipTransferred, error) {
	event := new(LiquidityManagerOwnershipTransferred)
	if err := _LiquidityManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LiquidityManagerTransmittedIterator struct {
	Event *LiquidityManagerTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LiquidityManagerTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityManagerTransmitted)
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
		it.Event = new(LiquidityManagerTransmitted)
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

func (it *LiquidityManagerTransmittedIterator) Error() error {
	return it.fail
}

func (it *LiquidityManagerTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LiquidityManagerTransmitted struct {
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_LiquidityManager *LiquidityManagerFilterer) FilterTransmitted(opts *bind.FilterOpts) (*LiquidityManagerTransmittedIterator, error) {

	logs, sub, err := _LiquidityManager.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &LiquidityManagerTransmittedIterator{contract: _LiquidityManager.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_LiquidityManager *LiquidityManagerFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *LiquidityManagerTransmitted) (event.Subscription, error) {

	logs, sub, err := _LiquidityManager.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LiquidityManagerTransmitted)
				if err := _LiquidityManager.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_LiquidityManager *LiquidityManagerFilterer) ParseTransmitted(log types.Log) (*LiquidityManagerTransmitted, error) {
	event := new(LiquidityManagerTransmitted)
	if err := _LiquidityManager.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}

func (_LiquidityManager *LiquidityManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LiquidityManager.abi.Events["ConfigSet"].ID:
		return _LiquidityManager.ParseConfigSet(log)
	case _LiquidityManager.abi.Events["CrossChainRebalancerSet"].ID:
		return _LiquidityManager.ParseCrossChainRebalancerSet(log)
	case _LiquidityManager.abi.Events["FinalizationFailed"].ID:
		return _LiquidityManager.ParseFinalizationFailed(log)
	case _LiquidityManager.abi.Events["FinalizationStepCompleted"].ID:
		return _LiquidityManager.ParseFinalizationStepCompleted(log)
	case _LiquidityManager.abi.Events["FinanceRoleSet"].ID:
		return _LiquidityManager.ParseFinanceRoleSet(log)
	case _LiquidityManager.abi.Events["LiquidityAddedToContainer"].ID:
		return _LiquidityManager.ParseLiquidityAddedToContainer(log)
	case _LiquidityManager.abi.Events["LiquidityContainerSet"].ID:
		return _LiquidityManager.ParseLiquidityContainerSet(log)
	case _LiquidityManager.abi.Events["LiquidityRemovedFromContainer"].ID:
		return _LiquidityManager.ParseLiquidityRemovedFromContainer(log)
	case _LiquidityManager.abi.Events["LiquidityTransferred"].ID:
		return _LiquidityManager.ParseLiquidityTransferred(log)
	case _LiquidityManager.abi.Events["MinimumLiquiditySet"].ID:
		return _LiquidityManager.ParseMinimumLiquiditySet(log)
	case _LiquidityManager.abi.Events["NativeDeposited"].ID:
		return _LiquidityManager.ParseNativeDeposited(log)
	case _LiquidityManager.abi.Events["NativeWithdrawn"].ID:
		return _LiquidityManager.ParseNativeWithdrawn(log)
	case _LiquidityManager.abi.Events["OwnershipTransferRequested"].ID:
		return _LiquidityManager.ParseOwnershipTransferRequested(log)
	case _LiquidityManager.abi.Events["OwnershipTransferred"].ID:
		return _LiquidityManager.ParseOwnershipTransferred(log)
	case _LiquidityManager.abi.Events["Transmitted"].ID:
		return _LiquidityManager.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LiquidityManagerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (LiquidityManagerCrossChainRebalancerSet) Topic() common.Hash {
	return common.HexToHash("0xab9bd0e4888101232b8f09dae2952ff59a6eea4a19fbddf2a8ca7b23f0e4cb40")
}

func (LiquidityManagerFinalizationFailed) Topic() common.Hash {
	return common.HexToHash("0xa481d91c3f9574c23ee84fef85246354b760a0527a535d6382354e4684703ce3")
}

func (LiquidityManagerFinalizationStepCompleted) Topic() common.Hash {
	return common.HexToHash("0x8d3121fe961b40270f336aa75feb1213f1c979a33993311c60da4dd0f24526cf")
}

func (LiquidityManagerFinanceRoleSet) Topic() common.Hash {
	return common.HexToHash("0x58024d20c07d3ebb87b192861d337d3a60995665acc5b8ce29596458b1f25170")
}

func (LiquidityManagerLiquidityAddedToContainer) Topic() common.Hash {
	return common.HexToHash("0x5414b81d05ac3542606f164e16a9a107d05d21e906539cc5ceb61d7b6b707eb5")
}

func (LiquidityManagerLiquidityContainerSet) Topic() common.Hash {
	return common.HexToHash("0x07dc474694ac40123aadcd2445f1b38d2eb353edd9319dcea043548ab34990ec")
}

func (LiquidityManagerLiquidityRemovedFromContainer) Topic() common.Hash {
	return common.HexToHash("0x2bda316674f8d73d289689d7a3acdf8e353b7a142fb5a68ac2aa475104039c18")
}

func (LiquidityManagerLiquidityTransferred) Topic() common.Hash {
	return common.HexToHash("0x2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a8")
}

func (LiquidityManagerMinimumLiquiditySet) Topic() common.Hash {
	return common.HexToHash("0xf97e758c8b3d81df7b0e1b7327a6a7fcf09a41536b2d274b9103015d715f11eb")
}

func (LiquidityManagerNativeDeposited) Topic() common.Hash {
	return common.HexToHash("0x3c597f6ac9fe7f0ed6da50b07618f5850a642e459ad587f7fab491a71f8b0ab8")
}

func (LiquidityManagerNativeWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x6b84d241b711af111ecfa0e518239e6ca212da442a76548fe8a1f4e77518256a")
}

func (LiquidityManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LiquidityManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LiquidityManagerTransmitted) Topic() common.Hash {
	return common.HexToHash("0xe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2")
}

func (_LiquidityManager *LiquidityManager) Address() common.Address {
	return _LiquidityManager.address
}

type LiquidityManagerInterface interface {
	GetAllCrossChainRebalancers(opts *bind.CallOpts) ([]ILiquidityManagerCrossChainRebalancerArgs, error)

	GetCrossChainRebalancer(opts *bind.CallOpts, chainSelector uint64) (LiquidityManagerCrossChainRebalancer, error)

	GetFinanceRole(opts *bind.CallOpts) (common.Address, error)

	GetLiquidity(opts *bind.CallOpts) (*big.Int, error)

	GetLocalLiquidityContainer(opts *bind.CallOpts) (common.Address, error)

	GetMinimumLiquidity(opts *bind.CallOpts) (*big.Int, error)

	GetSupportedDestChains(opts *bind.CallOpts) ([]uint64, error)

	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	ILocalToken(opts *bind.CallOpts) (common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestSequenceNumber(opts *bind.CallOpts) (uint64, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RebalanceLiquidity(opts *bind.TransactOpts, chainSelector uint64, amount *big.Int, nativeBridgeFee *big.Int, bridgeSpecificPayload []byte) (*types.Transaction, error)

	ReceiveLiquidity(opts *bind.TransactOpts, remoteChainSelector uint64, amount *big.Int, shouldWrapNative bool, bridgeSpecificPayload []byte) (*types.Transaction, error)

	RemoveLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	SetCrossChainRebalancer(opts *bind.TransactOpts, crossChainLiqManager ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error)

	SetCrossChainRebalancers(opts *bind.TransactOpts, crossChainRebalancers []ILiquidityManagerCrossChainRebalancerArgs) (*types.Transaction, error)

	SetFinanceRole(opts *bind.TransactOpts, finance common.Address) (*types.Transaction, error)

	SetLocalLiquidityContainer(opts *bind.TransactOpts, localLiquidityContainer common.Address) (*types.Transaction, error)

	SetMinimumLiquidity(opts *bind.TransactOpts, minimumLiquidity *big.Int) (*types.Transaction, error)

	SetOCR3Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, amount *big.Int, destination common.Address) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*LiquidityManagerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*LiquidityManagerConfigSet, error)

	FilterCrossChainRebalancerSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LiquidityManagerCrossChainRebalancerSetIterator, error)

	WatchCrossChainRebalancerSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerCrossChainRebalancerSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseCrossChainRebalancerSet(log types.Log) (*LiquidityManagerCrossChainRebalancerSet, error)

	FilterFinalizationFailed(opts *bind.FilterOpts, ocrSeqNum []uint64, remoteChainSelector []uint64) (*LiquidityManagerFinalizationFailedIterator, error)

	WatchFinalizationFailed(opts *bind.WatchOpts, sink chan<- *LiquidityManagerFinalizationFailed, ocrSeqNum []uint64, remoteChainSelector []uint64) (event.Subscription, error)

	ParseFinalizationFailed(log types.Log) (*LiquidityManagerFinalizationFailed, error)

	FilterFinalizationStepCompleted(opts *bind.FilterOpts, ocrSeqNum []uint64, remoteChainSelector []uint64) (*LiquidityManagerFinalizationStepCompletedIterator, error)

	WatchFinalizationStepCompleted(opts *bind.WatchOpts, sink chan<- *LiquidityManagerFinalizationStepCompleted, ocrSeqNum []uint64, remoteChainSelector []uint64) (event.Subscription, error)

	ParseFinalizationStepCompleted(log types.Log) (*LiquidityManagerFinalizationStepCompleted, error)

	FilterFinanceRoleSet(opts *bind.FilterOpts) (*LiquidityManagerFinanceRoleSetIterator, error)

	WatchFinanceRoleSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerFinanceRoleSet) (event.Subscription, error)

	ParseFinanceRoleSet(log types.Log) (*LiquidityManagerFinanceRoleSet, error)

	FilterLiquidityAddedToContainer(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LiquidityManagerLiquidityAddedToContainerIterator, error)

	WatchLiquidityAddedToContainer(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityAddedToContainer, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityAddedToContainer(log types.Log) (*LiquidityManagerLiquidityAddedToContainer, error)

	FilterLiquidityContainerSet(opts *bind.FilterOpts, newLiquidityContainer []common.Address) (*LiquidityManagerLiquidityContainerSetIterator, error)

	WatchLiquidityContainerSet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityContainerSet, newLiquidityContainer []common.Address) (event.Subscription, error)

	ParseLiquidityContainerSet(log types.Log) (*LiquidityManagerLiquidityContainerSet, error)

	FilterLiquidityRemovedFromContainer(opts *bind.FilterOpts, remover []common.Address, amount []*big.Int) (*LiquidityManagerLiquidityRemovedFromContainerIterator, error)

	WatchLiquidityRemovedFromContainer(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityRemovedFromContainer, remover []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityRemovedFromContainer(log types.Log) (*LiquidityManagerLiquidityRemovedFromContainer, error)

	FilterLiquidityTransferred(opts *bind.FilterOpts, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (*LiquidityManagerLiquidityTransferredIterator, error)

	WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *LiquidityManagerLiquidityTransferred, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (event.Subscription, error)

	ParseLiquidityTransferred(log types.Log) (*LiquidityManagerLiquidityTransferred, error)

	FilterMinimumLiquiditySet(opts *bind.FilterOpts) (*LiquidityManagerMinimumLiquiditySetIterator, error)

	WatchMinimumLiquiditySet(opts *bind.WatchOpts, sink chan<- *LiquidityManagerMinimumLiquiditySet) (event.Subscription, error)

	ParseMinimumLiquiditySet(log types.Log) (*LiquidityManagerMinimumLiquiditySet, error)

	FilterNativeDeposited(opts *bind.FilterOpts) (*LiquidityManagerNativeDepositedIterator, error)

	WatchNativeDeposited(opts *bind.WatchOpts, sink chan<- *LiquidityManagerNativeDeposited) (event.Subscription, error)

	ParseNativeDeposited(log types.Log) (*LiquidityManagerNativeDeposited, error)

	FilterNativeWithdrawn(opts *bind.FilterOpts) (*LiquidityManagerNativeWithdrawnIterator, error)

	WatchNativeWithdrawn(opts *bind.WatchOpts, sink chan<- *LiquidityManagerNativeWithdrawn) (event.Subscription, error)

	ParseNativeWithdrawn(log types.Log) (*LiquidityManagerNativeWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidityManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LiquidityManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LiquidityManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidityManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LiquidityManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LiquidityManagerOwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts) (*LiquidityManagerTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *LiquidityManagerTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*LiquidityManagerTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
