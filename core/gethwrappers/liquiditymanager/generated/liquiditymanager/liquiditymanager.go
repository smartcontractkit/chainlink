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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"finance\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidRemoteChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"CrossChainRebalancerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"FinalizationFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"}],\"name\":\"FinalizationStepCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"financeRole\",\"type\":\"address\"}],\"name\":\"FinanceRoleSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAddedToContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newLiquidityContainer\",\"type\":\"address\"}],\"name\":\"LiquidityContainerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemovedFromContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"fromChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"toChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeReturnData\",\"type\":\"bytes\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"MinimumLiquiditySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"NativeDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"NativeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainRebalancers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getCrossChainRebalancer\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityManager.CrossChainRebalancer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFinanceRole\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalLiquidityContainer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedDestChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_localToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nativeBridgeFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"rebalanceLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"shouldWrapNative\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"receiveLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs\",\"name\":\"crossChainLiqManager\",\"type\":\"tuple\"}],\"name\":\"setCrossChainRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"crossChainRebalancers\",\"type\":\"tuple[]\"}],\"name\":\"setCrossChainRebalancers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"finance\",\"type\":\"address\"}],\"name\":\"setFinanceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"}],\"name\":\"setLocalLiquidityContainer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"}],\"name\":\"setMinimumLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004a8338038062004a83833981016040819052620000349162000239565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000175565b505046608052506001600160401b038416600003620000f05760405163f89d762960e01b815260040160405180910390fd5b6001600160a01b03851615806200010e57506001600160a01b038316155b156200012d5760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0394851660a0526001600160401b0390931660c052600b80549285166001600160a01b0319938416179055600855600c8054929093169116179055620002b8565b336001600160a01b03821603620001cf5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200023657600080fd5b50565b600080600080600060a086880312156200025257600080fd5b85516200025f8162000220565b60208701519095506001600160401b03811681146200027d57600080fd5b6040870151909450620002908162000220565b606087015160808801519194509250620002aa8162000220565b809150509295509295909350565b60805160a05160c0516147556200032e6000396000818161315a015261332f01526000818161045a0152818161075e015281816109fb01528181610a410152818161169d01528181613038015281816130b9015281816131e2015261328001526000818161183a015261188601526147556000f3fe6080604052600436106101b05760003560e01c8063791781f5116100ec578063b1dc65a41161008a578063da9c0f9611610064578063da9c0f9614610686578063f2fde38b146106b1578063f8c2d8fa146106d1578063fe65d5af146106f157600080fd5b8063b1dc65a414610626578063b7e7fa0514610646578063b8ca8dd81461066657600080fd5b806383d34afe116100c657806383d34afe146105805780638da5cb5b146105955780639c8f9f23146105c0578063afcb95d7146105e057600080fd5b8063791781f51461050357806379ba50971461052e57806381ff70481461054357600080fd5b806350a197d7116101595780636511d919116101335780636511d91914610448578063666cab8d146104a15780636a11ee90146104c3578063706bf645146104e357600080fd5b806350a197d7146102f657806351c6590a14610316578063568446e71461033657600080fd5b80633275636e1161018a5780633275636e14610294578063348759c1146102b45780634f814d04146102d657600080fd5b80630910a510146101f4578063181f5a771461021c578063282567b41461027257600080fd5b366101ef57604080513481523360208201527f3c597f6ac9fe7f0ed6da50b07618f5850a642e459ad587f7fab491a71f8b0ab8910160405180910390a1005b600080fd5b34801561020057600080fd5b50610209610713565b6040519081526020015b60405180910390f35b34801561022857600080fd5b506102656040518060400160405280601a81526020017f4c69717569646974794d616e6167657220312e302e302d64657600000000000081525081565b6040516102139190613906565b34801561027e57600080fd5b5061029261028d366004613920565b6107ce565b005b3480156102a057600080fd5b506102926102af366004613939565b61081b565b3480156102c057600080fd5b506102c961082f565b6040516102139190613951565b3480156102e257600080fd5b506102926102f13660046139c1565b6108bb565b34801561030257600080fd5b50610292610311366004613a52565b61093c565b34801561032257600080fd5b50610292610331366004613920565b6109e1565b34801561034257600080fd5b506103f9610351366004613ac3565b6040805160808101825260008082526020820181905291810182905260608101919091525067ffffffffffffffff166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff161515606082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff9081168252602080850151821690830152838301511691810191909152606091820151151591810191909152608001610213565b34801561045457600080fd5b5061047c7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610213565b3480156104ad57600080fd5b506104b6610b1e565b6040516102139190613b30565b3480156104cf57600080fd5b506102926104de366004613d51565b610b8c565b3480156104ef57600080fd5b506102926104fe3660046139c1565b611396565b34801561050f57600080fd5b50600b5473ffffffffffffffffffffffffffffffffffffffff1661047c565b34801561053a57600080fd5b5061029261145a565b34801561054f57600080fd5b506004546002546040805163ffffffff80851682526401000000009094049093166020840152820152606001610213565b34801561058c57600080fd5b50600854610209565b3480156105a157600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661047c565b3480156105cc57600080fd5b506102926105db366004613920565b611557565b3480156105ec57600080fd5b50600254600454604080516001815260208101939093526801000000000000000090910467ffffffffffffffff1690820152606001610213565b34801561063257600080fd5b50610292610641366004613e63565b6116f7565b34801561065257600080fd5b50610292610661366004613f1a565b611d68565b34801561067257600080fd5b50610292610681366004613f8f565b611da8565b34801561069257600080fd5b50600c5473ffffffffffffffffffffffffffffffffffffffff1661047c565b3480156106bd57600080fd5b506102926106cc3660046139c1565b611ee6565b3480156106dd57600080fd5b506102926106ec366004613fbf565b611ef7565b3480156106fd57600080fd5b50610706611f93565b604051610213919061400a565b600b546040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526000917f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa1580156107a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107c9919061409f565b905090565b6107d6612154565b600880549082905560408051828152602081018490527ff97e758c8b3d81df7b0e1b7327a6a7fcf09a41536b2d274b9103015d715f11eb910160405180910390a15050565b610823612154565b61082c816121d7565b50565b6060600a8054806020026020016040519081016040528092919081815260200182805480156108b157602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff168152602001906008019060208260070104928301926001038202915080841161086c5790505b5050505050905090565b6108c3612154565b600c80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f58024d20c07d3ebb87b192861d337d3a60995665acc5b8ce29596458b1f251709060200160405180910390a150565b600c5473ffffffffffffffffffffffffffffffffffffffff16331461098d576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109da858584848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525089925067ffffffffffffffff91506125d39050565b5050505050565b610a2373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333084612858565b600b54610a6a9073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691168361293a565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b158015610ad657600080fd5b505af1158015610aea573d6000803e3d6000fd5b50506040518392503391507f5414b81d05ac3542606f164e16a9a107d05d21e906539cc5ceb61d7b6b707eb590600090a350565b606060078054806020026020016040519081016040528092919081815260200182805480156108b157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610b58575050505050905090565b855185518560ff16601f831115610c04576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b80600003610c6e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610bfb565b818314610cfc576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610bfb565b610d078160036140e7565b8311610d6f576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610bfb565b610d77612154565b60065460005b81811015610e6b576005600060068381548110610d9c57610d9c614104565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560078054600592919084908110610e0c57610e0c614104565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600101610d7d565b50895160005b8181101561123e5760008c8281518110610e8d57610e8d614104565b6020026020010151905060006002811115610eaa57610eaa614133565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff166002811115610ee957610ee9614133565b14610f50576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610bfb565b73ffffffffffffffffffffffffffffffffffffffff8116610f9d576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561104d5761104d614133565b021790555090505060008c838151811061106957611069614104565b602002602001015190506000600281111561108657611086614133565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff1660028111156110c5576110c5614133565b1461112c576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610bfb565b73ffffffffffffffffffffffffffffffffffffffff8116611179576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561122957611229614133565b02179055509050505050806001019050610e71565b508a516112529060069060208e01906137da565b5089516112669060079060208d01906137da565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c1617179055600480546112ec9146913091906000906112be9063ffffffff16614162565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e612abc565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f60405161138099989796959493929190614185565b60405180910390a1505050505050505050505050565b61139e612154565b73ffffffffffffffffffffffffffffffffffffffff81166113eb576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517f07dc474694ac40123aadcd2445f1b38d2eb353edd9319dcea043548ab34990ec90600090a250565b60015473ffffffffffffffffffffffffffffffffffffffff1633146114db576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610bfb565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600c5473ffffffffffffffffffffffffffffffffffffffff1633146115a8576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006115b2610713565b9050818110156115ff576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018390526024810182905260006044820152606401610bfb565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b15801561166b57600080fd5b505af115801561167f573d6000803e3d6000fd5b506116c692505073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690503384612b67565b604051829033907f2bda316674f8d73d289689d7a3acdf8e353b7a142fb5a68ac2aa475104039c1890600090a35050565b60045460208901359067ffffffffffffffff6801000000000000000090910481169082161161177a57600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff80851693820193909352680100000000000000009091049091166024820152604401610bfb565b611785888883612bbd565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a359182146118375780516040517f93df584c000000000000000000000000000000000000000000000000000000008152600481019190915260248101839052604401610bfb565b467f0000000000000000000000000000000000000000000000000000000000000000146118b8576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610bfb565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a1602081015161190a90600161421b565b60ff168714611945576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b86851461197e576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156119c1576119c1614133565b60028111156119d2576119d2614133565b90525090506002816020015160028111156119ef576119ef614133565b148015611a3657506007816000015160ff1681548110611a1157611a11614104565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611a6c576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000611a7a8660206140e7565b611a858960206140e7565b611a918c610144614234565b611a9b9190614234565b611aa59190614234565b9050368114611ae9576040517f8e1192e100000000000000000000000000000000000000000000000000000000815260048101829052366024820152604401610bfb565b5060008a8a604051611afc929190614247565b604051908190038120611b13918e90602001614257565b604051602081830303815290604052805190602001209050611b33613864565b8860005b81811015611d575760006001858a8460208110611b5657611b56614104565b611b6391901a601b61421b565b8f8f86818110611b7557611b75614104565b905060200201358e8e87818110611b8e57611b8e614104565b9050602002013560405160008152602001604052604051611bcb949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611bed573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff8082168652939750919550929392840191610100909104166002811115611c7057611c70614133565b6002811115611c8157611c81614133565b9052509050600181602001516002811115611c9e57611c9e614133565b14611cd5576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f8110611cec57611cec614104565b602002015115611d28576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f8110611d4357611d43614104565b911515602090920201525050600101611b37565b505050505050505050505050505050565b611d70612154565b60005b81811015611da357611d9b838383818110611d9057611d90614104565b905060a002016121d7565b600101611d73565b505050565b600c5473ffffffffffffffffffffffffffffffffffffffff163314611df9576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008173ffffffffffffffffffffffffffffffffffffffff168360405160006040518083038185875af1925050503d8060008114611e53576040519150601f19603f3d011682016040523d82523d6000602084013e611e58565b606091505b5050905080611e93576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805184815273ffffffffffffffffffffffffffffffffffffffff841660208201527f6b84d241b711af111ecfa0e518239e6ca212da442a76548fe8a1f4e77518256a910160405180910390a1505050565b611eee612154565b61082c81612d6c565b600c5473ffffffffffffffffffffffffffffffffffffffff163314611f48576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109da85858567ffffffffffffffff86868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612e6192505050565b600a5460609060008167ffffffffffffffff811115611fb457611fb4613b43565b60405190808252806020026020018201604052801561202b57816020015b6040805160a0810182526000808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181611fd25790505b50905060005b8281101561214d576000600a828154811061204e5761204e614104565b6000918252602080832060048304015460039092166008026101000a90910467ffffffffffffffff1680835260098252604092839020835160808082018652825473ffffffffffffffffffffffffffffffffffffffff9081168352600184015481168387019081526002909401548082168489019081527401000000000000000000000000000000000000000090910460ff1615156060808601918252895160a081018b5286518516815296518416988701989098529051909116968401969096529382018390529351151592810192909252855190935085908590811061213857612138614104565b60209081029190910101525050600101612031565b5092915050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146121d5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bfb565b565b6121e76080820160608301613ac3565b67ffffffffffffffff1660000361222a576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061223960208301836139c1565b73ffffffffffffffffffffffffffffffffffffffff1614806122805750600061226860408301602084016139c1565b73ffffffffffffffffffffffffffffffffffffffff16145b806122b05750600061229860608301604084016139c1565b73ffffffffffffffffffffffffffffffffffffffff16145b156122e7576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006009816122fc6080850160608601613ac3565b67ffffffffffffffff16815260208101919091526040016000206002015473ffffffffffffffffffffffffffffffffffffffff160361238857600a6123476080830160608401613ac3565b8154600181018355600092835260209092206004830401805460039093166008026101000a67ffffffffffffffff8181021990941692909316929092021790555b6040805160808101909152806123a160208401846139c1565b73ffffffffffffffffffffffffffffffffffffffff1681526020018260200160208101906123cf91906139c1565b73ffffffffffffffffffffffffffffffffffffffff1681526020016123fa60608401604085016139c1565b73ffffffffffffffffffffffffffffffffffffffff16815260200161242560a084016080850161426b565b151590526009600061243d6080850160608601613ac3565b67ffffffffffffffff16815260208082019290925260409081016000208351815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161783559385015160018301805491831691909516179093559083015160029091018054606094850151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff000000000000000000000000000000000000000000909116929093169190911791909117905561251e9060808301908301613ac3565b67ffffffffffffffff167fab9bd0e4888101232b8f09dae2952ff59a6eea4a19fbddf2a8ca7b23f0e4cb4061255960408401602085016139c1565b61256960608501604086016139c1565b61257660208601866139c1565b61258660a087016080880161426b565b6040516125c8949392919073ffffffffffffffffffffffffffffffffffffffff9485168152928416602084015292166040820152901515606082015260800190565b60405180910390a250565b67ffffffffffffffff85166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052612698576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152602401610bfb565b602081015181516040517f38314bb200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909216916338314bb2916126f49130908990600401614288565b6020604051808303816000875af192505050801561274d575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261274a918101906142ca565b60015b6127d6573d80801561277b576040519150601f19603f3d011682016040523d82523d6000602084013e612780565b606091505b508667ffffffffffffffff168367ffffffffffffffff167fa481d91c3f9574c23ee84fef85246354b760a0527a535d6382354e4684703ce387846040516127c89291906142e7565b60405180910390a350612843565b80156127ee576127e986848988886131da565b61283c565b8667ffffffffffffffff168367ffffffffffffffff167f8d3121fe961b40270f336aa75feb1213f1c979a33993311c60da4dd0f24526cf876040516128339190613906565b60405180910390a35b50506109da565b61285085838887876131da565b505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526129349085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526133c1565b50505050565b8015806129da57506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa1580156129b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129d8919061409f565b155b612a66576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608401610bfb565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052611da39084907f095ea7b300000000000000000000000000000000000000000000000000000000906064016128b2565b6000808a8a8a8a8a8a8a8a8a604051602001612ae09998979695949392919061430c565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052611da39084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064016128b2565b6000612bcb838501856144a9565b8051516020820151519192509081158015612be4575080155b15612c1a576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82811015612cbe57612cb684600001518281518110612c3e57612c3e614104565b60200260200101516040015185600001518381518110612c6057612c60614104565b60200260200101516000015186600001518481518110612c8257612c82614104565b6020026020010151602001518888600001518681518110612ca557612ca5614104565b602002602001015160600151612e61565b600101612c1d565b5060005b81811015612d6357612d5b84602001518281518110612ce357612ce3614104565b60200260200101516020015185602001518381518110612d0557612d05614104565b60200260200101516000015186602001518481518110612d2757612d27614104565b60200260200101516060015187602001518581518110612d4957612d49614104565b602002602001015160400151896125d3565b600101612cc2565b50505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612deb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bfb565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612e6b610713565b60085490915080821080612e87575085612e85828461461d565b105b15612ecf576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018790526024810183905260448101829052606401610bfb565b67ffffffffffffffff87166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052612f94576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff89166004820152602401610bfb565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810189905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b15801561300057600080fd5b505af1158015613014573d6000803e3d6000fd5b505050602082015161305f915073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016908961293a565b6020810151604080830151835191517fa71d98b700000000000000000000000000000000000000000000000000000000815260009373ffffffffffffffffffffffffffffffffffffffff169263a71d98b7928b926130e6927f000000000000000000000000000000000000000000000000000000000000000092918f908d90600401614630565b60006040518083038185885af1158015613104573d6000803e3d6000fd5b50505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261314b9190810190614677565b90508867ffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168767ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a885600001518c8a876040516131c794939291906146e5565b60405180910390a4505050505050505050565b8015613262577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0866040518263ffffffff1660e01b81526004016000604051808303818588803b15801561324857600080fd5b505af115801561325c573d6000803e3d6000fd5b50505050505b600b546132a99073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169116876134cd565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810187905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b15801561331557600080fd5b505af1158015613329573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168367ffffffffffffffff168567ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a8308987604051806020016040528060008152506040516133b294939291906146e5565b60405180910390a45050505050565b6000613423826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166135cb9092919063ffffffff16565b805190915015611da3578080602001905181019061344191906142ca565b611da3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610bfb565b6040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff8381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa158015613544573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613568919061409f565b6135729190614234565b60405173ffffffffffffffffffffffffffffffffffffffff85166024820152604481018290529091506129349085907f095ea7b300000000000000000000000000000000000000000000000000000000906064016128b2565b60606135da84846000856135e2565b949350505050565b606082471015613674576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610bfb565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161369d919061472c565b60006040518083038185875af1925050503d80600081146136da576040519150601f19603f3d011682016040523d82523d6000602084013e6136df565b606091505b50915091506136f0878383876136fb565b979650505050505050565b6060831561379157825160000361378a5773ffffffffffffffffffffffffffffffffffffffff85163b61378a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610bfb565b50816135da565b6135da83838151156137a65781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bfb9190613906565b828054828255906000526020600020908101928215613854579160200282015b8281111561385457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906137fa565b50613860929150613883565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b808211156138605760008155600101613884565b60005b838110156138b357818101518382015260200161389b565b50506000910152565b600081518084526138d4816020860160208601613898565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061391960208301846138bc565b9392505050565b60006020828403121561393257600080fd5b5035919050565b600060a0828403121561394b57600080fd5b50919050565b6020808252825182820181905260009190848201906040850190845b8181101561399357835167ffffffffffffffff168352928401929184019160010161396d565b50909695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461082c57600080fd5b6000602082840312156139d357600080fd5b81356139198161399f565b803567ffffffffffffffff811681146139f657600080fd5b919050565b801515811461082c57600080fd5b60008083601f840112613a1b57600080fd5b50813567ffffffffffffffff811115613a3357600080fd5b602083019150836020828501011115613a4b57600080fd5b9250929050565b600080600080600060808688031215613a6a57600080fd5b613a73866139de565b9450602086013593506040860135613a8a816139fb565b9250606086013567ffffffffffffffff811115613aa657600080fd5b613ab288828901613a09565b969995985093965092949392505050565b600060208284031215613ad557600080fd5b613919826139de565b60008151808452602080850194506020840160005b83811015613b2557815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613af3565b509495945050505050565b6020815260006139196020830184613ade565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613b9557613b95613b43565b60405290565b6040805190810167ffffffffffffffff81118282101715613b9557613b95613b43565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613c0557613c05613b43565b604052919050565b600067ffffffffffffffff821115613c2757613c27613b43565b5060051b60200190565b600082601f830112613c4257600080fd5b81356020613c57613c5283613c0d565b613bbe565b8083825260208201915060208460051b870101935086841115613c7957600080fd5b602086015b84811015613c9e578035613c918161399f565b8352918301918301613c7e565b509695505050505050565b803560ff811681146139f657600080fd5b600067ffffffffffffffff821115613cd457613cd4613b43565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613d1157600080fd5b8135613d1f613c5282613cba565b818152846020838601011115613d3457600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215613d6a57600080fd5b863567ffffffffffffffff80821115613d8257600080fd5b613d8e8a838b01613c31565b97506020890135915080821115613da457600080fd5b613db08a838b01613c31565b9650613dbe60408a01613ca9565b95506060890135915080821115613dd457600080fd5b613de08a838b01613d00565b9450613dee60808a016139de565b935060a0890135915080821115613e0457600080fd5b50613e1189828a01613d00565b9150509295509295509295565b60008083601f840112613e3057600080fd5b50813567ffffffffffffffff811115613e4857600080fd5b6020830191508360208260051b8501011115613a4b57600080fd5b60008060008060008060008060e0898b031215613e7f57600080fd5b606089018a811115613e9057600080fd5b8998503567ffffffffffffffff80821115613eaa57600080fd5b613eb68c838d01613a09565b909950975060808b0135915080821115613ecf57600080fd5b613edb8c838d01613e1e565b909750955060a08b0135915080821115613ef457600080fd5b50613f018b828c01613e1e565b999c989b50969995989497949560c00135949350505050565b60008060208385031215613f2d57600080fd5b823567ffffffffffffffff80821115613f4557600080fd5b818501915085601f830112613f5957600080fd5b813581811115613f6857600080fd5b86602060a083028501011115613f7d57600080fd5b60209290920196919550909350505050565b60008060408385031215613fa257600080fd5b823591506020830135613fb48161399f565b809150509250929050565b600080600080600060808688031215613fd757600080fd5b613fe0866139de565b94506020860135935060408601359250606086013567ffffffffffffffff811115613aa657600080fd5b602080825282518282018190526000919060409081850190868401855b82811015614092578151805173ffffffffffffffffffffffffffffffffffffffff90811686528782015181168887015286820151168686015260608082015167ffffffffffffffff169086015260809081015115159085015260a09093019290850190600101614027565b5091979650505050505050565b6000602082840312156140b157600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176140fe576140fe6140b8565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff80831681810361417b5761417b6140b8565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526141b58184018a613ade565b905082810360808401526141c98189613ade565b905060ff871660a084015282810360c08401526141e681876138bc565b905067ffffffffffffffff851660e084015282810361010084015261420b81856138bc565b9c9b505050505050505050505050565b60ff81811683821601908111156140fe576140fe6140b8565b808201808211156140fe576140fe6140b8565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006020828403121561427d57600080fd5b8135613919816139fb565b600073ffffffffffffffffffffffffffffffffffffffff8086168352808516602084015250606060408301526142c160608301846138bc565b95945050505050565b6000602082840312156142dc57600080fd5b8151613919816139fb565b6040815260006142fa60408301856138bc565b82810360208401526142c181856138bc565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526143538285018b613ade565b91508382036080850152614367828a613ade565b915060ff881660a085015283820360c085015261438482886138bc565b90861660e0850152838103610100850152905061420b81856138bc565b600082601f8301126143b257600080fd5b813560206143c2613c5283613c0d565b82815260059290921b840181019181810190868411156143e157600080fd5b8286015b84811015613c9e57803567ffffffffffffffff808211156144065760008081fd5b81890191506080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d0301121561443f5760008081fd5b614447613b72565b878401358152604061445a8186016139de565b8983015260608086013561446d816139fb565b8383015292850135928484111561448657600091508182fd5b6144948e8b86890101613d00565b908301525086525050509183019183016143e5565b600060208083850312156144bc57600080fd5b823567ffffffffffffffff808211156144d457600080fd5b90840190604082870312156144e857600080fd5b6144f0613b9b565b8235828111156144ff57600080fd5b8301601f8101881361451057600080fd5b803561451e613c5282613c0d565b81815260059190911b8201860190868101908a83111561453d57600080fd5b8784015b838110156145e95780358781111561455857600080fd5b85016080818e037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001121561458c57600080fd5b614594613b72565b8a820135815260408201358b8201526145af606083016139de565b60408201526080820135898111156145c75760008081fd5b6145d58f8d83860101613d00565b606083015250845250918801918801614541565b508452505050828401358281111561460057600080fd5b61460c888286016143a1565b948201949094529695505050505050565b818103818111156140fe576140fe6140b8565b600073ffffffffffffffffffffffffffffffffffffffff8088168352808716602084015280861660408401525083606083015260a060808301526136f060a08301846138bc565b60006020828403121561468957600080fd5b815167ffffffffffffffff8111156146a057600080fd5b8201601f810184136146b157600080fd5b80516146bf613c5282613cba565b8181528560208385010111156146d457600080fd5b6142c1826020830160208601613898565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260806040820152600061471a60808301856138bc565b82810360608401526136f081856138bc565b6000825161473e818460208701613898565b919091019291505056fea164736f6c6343000818000a",
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

func (_LiquidityManager *LiquidityManagerCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _LiquidityManager.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.SequenceNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_LiquidityManager *LiquidityManagerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _LiquidityManager.Contract.LatestConfigDigestAndEpoch(&_LiquidityManager.CallOpts)
}

func (_LiquidityManager *LiquidityManagerCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _LiquidityManager.Contract.LatestConfigDigestAndEpoch(&_LiquidityManager.CallOpts)
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
type LatestConfigDigestAndEpoch struct {
	ScanLogs       bool
	ConfigDigest   [32]byte
	SequenceNumber uint64
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

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

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
