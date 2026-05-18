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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"finance\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidRemoteChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"CrossChainRebalancerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"FinalizationFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"}],\"name\":\"FinalizationStepCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"financeRole\",\"type\":\"address\"}],\"name\":\"FinanceRoleSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAddedToContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newLiquidityContainer\",\"type\":\"address\"}],\"name\":\"LiquidityContainerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemovedFromContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"fromChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"toChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeReturnData\",\"type\":\"bytes\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"MinimumLiquiditySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"NativeDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"NativeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainRebalancers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getCrossChainRebalancer\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityManager.CrossChainRebalancer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFinanceRole\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalLiquidityContainer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedDestChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_localToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nativeBridgeFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"rebalanceLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"shouldWrapNative\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"receiveLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs\",\"name\":\"crossChainLiqManager\",\"type\":\"tuple\"}],\"name\":\"setCrossChainRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"crossChainRebalancers\",\"type\":\"tuple[]\"}],\"name\":\"setCrossChainRebalancers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"finance\",\"type\":\"address\"}],\"name\":\"setFinanceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"}],\"name\":\"setLocalLiquidityContainer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"}],\"name\":\"setMinimumLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004b8538038062004b85833981016040819052620000349162000239565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000175565b505046608052506001600160401b038416600003620000f05760405163f89d762960e01b815260040160405180910390fd5b6001600160a01b03851615806200010e57506001600160a01b038316155b156200012d5760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0394851660a0526001600160401b0390931660c052600b80549285166001600160a01b0319938416179055600855600c8054929093169116179055620002b8565b336001600160a01b03821603620001cf5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200023657600080fd5b50565b600080600080600060a086880312156200025257600080fd5b85516200025f8162000220565b60208701519095506001600160401b03811681146200027d57600080fd5b6040870151909450620002908162000220565b606087015160808801519194509250620002aa8162000220565b809150509295509295909350565b60805160a05160c0516148576200032e6000396000818161321a01526133ef0152600081816104850152818161078201528181610a1f01528181610a6501528181611762015281816130f801528181613179015281816132a201526133400152600081816118ff015261194b01526148576000f3fe6080604052600436106101bb5760003560e01c8063791781f5116100ec578063b7e7fa051161008a578063f1c0461611610064578063f1c0461614610696578063f2fde38b146106d5578063f8c2d8fa146106f5578063fe65d5af1461071557600080fd5b8063b7e7fa051461062b578063b8ca8dd81461064b578063da9c0f961461066b57600080fd5b806383d34afe116100c657806383d34afe146105ab5780638da5cb5b146105c05780639c8f9f23146105eb578063b1dc65a41461060b57600080fd5b8063791781f51461052e57806379ba50971461055957806381ff70481461056e57600080fd5b806351c6590a116101595780636511d919116101335780636511d91914610473578063666cab8d146104cc5780636a11ee90146104ee578063706bf6451461050e57600080fd5b806351c6590a14610321578063568446e7146103415780635fc3ea0b1461045357600080fd5b80633275636e116101955780633275636e1461029f578063348759c1146102bf5780634f814d04146102e157806350a197d71461030157600080fd5b80630910a510146101ff578063181f5a7714610227578063282567b41461027d57600080fd5b366101fa57604080513481523360208201527f3c597f6ac9fe7f0ed6da50b07618f5850a642e459ad587f7fab491a71f8b0ab8910160405180910390a1005b600080fd5b34801561020b57600080fd5b50610214610737565b6040519081526020015b60405180910390f35b34801561023357600080fd5b506102706040518060400160405280601a81526020017f4c69717569646974794d616e6167657220312e302e302d64657600000000000081525081565b60405161021e91906139c6565b34801561028957600080fd5b5061029d6102983660046139e0565b6107f2565b005b3480156102ab57600080fd5b5061029d6102ba3660046139f9565b61083f565b3480156102cb57600080fd5b506102d4610853565b60405161021e9190613a11565b3480156102ed57600080fd5b5061029d6102fc366004613a81565b6108df565b34801561030d57600080fd5b5061029d61031c366004613b12565b610960565b34801561032d57600080fd5b5061029d61033c3660046139e0565b610a05565b34801561034d57600080fd5b5061040461035c366004613b83565b6040805160808101825260008082526020820181905291810182905260608101919091525067ffffffffffffffff166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff161515606082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff908116825260208085015182169083015283830151169181019190915260609182015115159181019190915260800161021e565b34801561045f57600080fd5b5061029d61046e366004613b9e565b610b42565b34801561047f57600080fd5b506104a77f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161021e565b3480156104d857600080fd5b506104e1610bb9565b60405161021e9190613c32565b3480156104fa57600080fd5b5061029d610509366004613e53565b610c27565b34801561051a57600080fd5b5061029d610529366004613a81565b61145b565b34801561053a57600080fd5b50600b5473ffffffffffffffffffffffffffffffffffffffff166104a7565b34801561056557600080fd5b5061029d61151f565b34801561057a57600080fd5b506004546002546040805163ffffffff8085168252640100000000909404909316602084015282015260600161021e565b3480156105b757600080fd5b50600854610214565b3480156105cc57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166104a7565b3480156105f757600080fd5b5061029d6106063660046139e0565b61161c565b34801561061757600080fd5b5061029d610626366004613f65565b6117bc565b34801561063757600080fd5b5061029d61064636600461401c565b611e2d565b34801561065757600080fd5b5061029d610666366004614091565b611e68565b34801561067757600080fd5b50600c5473ffffffffffffffffffffffffffffffffffffffff166104a7565b3480156106a257600080fd5b5060045468010000000000000000900467ffffffffffffffff1660405167ffffffffffffffff909116815260200161021e565b3480156106e157600080fd5b5061029d6106f0366004613a81565b611fa6565b34801561070157600080fd5b5061029d6107103660046140c1565b611fb7565b34801561072157600080fd5b5061072a612053565b60405161021e919061410c565b600b546040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526000917f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa1580156107c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107ed91906141a1565b905090565b6107fa612214565b600880549082905560408051828152602081018490527ff97e758c8b3d81df7b0e1b7327a6a7fcf09a41536b2d274b9103015d715f11eb910160405180910390a15050565b610847612214565b61085081612297565b50565b6060600a8054806020026020016040519081016040528092919081815260200182805480156108d557602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116108905790505b5050505050905090565b6108e7612214565b600c80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f58024d20c07d3ebb87b192861d337d3a60995665acc5b8ce29596458b1f251709060200160405180910390a150565b600c5473ffffffffffffffffffffffffffffffffffffffff1633146109b1576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109fe858584848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525089925067ffffffffffffffff91506126939050565b5050505050565b610a4773ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333084612918565b600b54610a8e9073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169116836129fa565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b158015610afa57600080fd5b505af1158015610b0e573d6000803e3d6000fd5b50506040518392503391507f5414b81d05ac3542606f164e16a9a107d05d21e906539cc5ceb61d7b6b707eb590600090a350565b600c5473ffffffffffffffffffffffffffffffffffffffff163314610b93576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610bb473ffffffffffffffffffffffffffffffffffffffff84168284612b7c565b505050565b606060078054806020026020016040519081016040528092919081815260200182805480156108d557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610bf3575050505050905090565b855185518560ff16601f831115610c9f576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b80600003610d09576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610c96565b818314610d97576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610c96565b610da28160036141e9565b8311610e0a576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610c96565b610e12612214565b60065460005b81811015610f06576005600060068381548110610e3757610e37614206565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560078054600592919084908110610ea757610ea7614206565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600101610e18565b50895160005b818110156112d95760008c8281518110610f2857610f28614206565b6020026020010151905060006002811115610f4557610f45614235565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff166002811115610f8457610f84614235565b14610feb576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610c96565b73ffffffffffffffffffffffffffffffffffffffff8116611038576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156110e8576110e8614235565b021790555090505060008c838151811061110457611104614206565b602002602001015190506000600281111561112157611121614235565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff16600281111561116057611160614235565b146111c7576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610c96565b73ffffffffffffffffffffffffffffffffffffffff8116611214576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156112c4576112c4614235565b02179055509050505050806001019050610f0c565b508a516112ed9060069060208e019061389a565b5089516113019060079060208d019061389a565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c1617179055600480546113879146913091906000906113599063ffffffff16614264565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e612bd2565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055506000600460086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f60405161144599989796959493929190614287565b60405180910390a1505050505050505050505050565b611463612214565b73ffffffffffffffffffffffffffffffffffffffff81166114b0576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517f07dc474694ac40123aadcd2445f1b38d2eb353edd9319dcea043548ab34990ec90600090a250565b60015473ffffffffffffffffffffffffffffffffffffffff1633146115a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610c96565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600c5473ffffffffffffffffffffffffffffffffffffffff16331461166d576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611677610737565b9050818110156116c4576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018390526024810182905260006044820152606401610c96565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b15801561173057600080fd5b505af1158015611744573d6000803e3d6000fd5b5061178b92505073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690503384612b7c565b604051829033907f2bda316674f8d73d289689d7a3acdf8e353b7a142fb5a68ac2aa475104039c1890600090a35050565b60045460208901359067ffffffffffffffff6801000000000000000090910481169082161161183f57600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff80851693820193909352680100000000000000009091049091166024820152604401610c96565b61184a888883612c7d565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a359182146118fc5780516040517f93df584c000000000000000000000000000000000000000000000000000000008152600481019190915260248101839052604401610c96565b467f00000000000000000000000000000000000000000000000000000000000000001461197d576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610c96565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a160208101516119cf90600161431d565b60ff168714611a0a576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b868514611a43576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115611a8657611a86614235565b6002811115611a9757611a97614235565b9052509050600281602001516002811115611ab457611ab4614235565b148015611afb57506007816000015160ff1681548110611ad657611ad6614206565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611b31576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000611b3f8660206141e9565b611b4a8960206141e9565b611b568c610144614336565b611b609190614336565b611b6a9190614336565b9050368114611bae576040517f8e1192e100000000000000000000000000000000000000000000000000000000815260048101829052366024820152604401610c96565b5060008a8a604051611bc1929190614349565b604051908190038120611bd8918e90602001614359565b604051602081830303815290604052805190602001209050611bf8613924565b8860005b81811015611e1c5760006001858a8460208110611c1b57611c1b614206565b611c2891901a601b61431d565b8f8f86818110611c3a57611c3a614206565b905060200201358e8e87818110611c5357611c53614206565b9050602002013560405160008152602001604052604051611c90949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611cb2573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff8082168652939750919550929392840191610100909104166002811115611d3557611d35614235565b6002811115611d4657611d46614235565b9052509050600181602001516002811115611d6357611d63614235565b14611d9a576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f8110611db157611db1614206565b602002015115611ded576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f8110611e0857611e08614206565b911515602090920201525050600101611bfc565b505050505050505050505050505050565b611e35612214565b60005b81811015610bb457611e60838383818110611e5557611e55614206565b905060a00201612297565b600101611e38565b600c5473ffffffffffffffffffffffffffffffffffffffff163314611eb9576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008173ffffffffffffffffffffffffffffffffffffffff168360405160006040518083038185875af1925050503d8060008114611f13576040519150601f19603f3d011682016040523d82523d6000602084013e611f18565b606091505b5050905080611f53576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805184815273ffffffffffffffffffffffffffffffffffffffff841660208201527f6b84d241b711af111ecfa0e518239e6ca212da442a76548fe8a1f4e77518256a910160405180910390a1505050565b611fae612214565b61085081612e2c565b600c5473ffffffffffffffffffffffffffffffffffffffff163314612008576040517fb2a59b2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109fe85858567ffffffffffffffff86868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612f2192505050565b600a5460609060008167ffffffffffffffff81111561207457612074613c45565b6040519080825280602002602001820160405280156120eb57816020015b6040805160a0810182526000808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816120925790505b50905060005b8281101561220d576000600a828154811061210e5761210e614206565b6000918252602080832060048304015460039092166008026101000a90910467ffffffffffffffff1680835260098252604092839020835160808082018652825473ffffffffffffffffffffffffffffffffffffffff9081168352600184015481168387019081526002909401548082168489019081527401000000000000000000000000000000000000000090910460ff1615156060808601918252895160a081018b528651851681529651841698870198909852905190911696840196909652938201839052935115159281019290925285519093508590859081106121f8576121f8614206565b602090810291909101015250506001016120f1565b5092915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612295576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610c96565b565b6122a76080820160608301613b83565b67ffffffffffffffff166000036122ea576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006122f96020830183613a81565b73ffffffffffffffffffffffffffffffffffffffff161480612340575060006123286040830160208401613a81565b73ffffffffffffffffffffffffffffffffffffffff16145b80612370575060006123586060830160408401613a81565b73ffffffffffffffffffffffffffffffffffffffff16145b156123a7576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006009816123bc6080850160608601613b83565b67ffffffffffffffff16815260208101919091526040016000206002015473ffffffffffffffffffffffffffffffffffffffff160361244857600a6124076080830160608401613b83565b8154600181018355600092835260209092206004830401805460039093166008026101000a67ffffffffffffffff8181021990941692909316929092021790555b6040805160808101909152806124616020840184613a81565b73ffffffffffffffffffffffffffffffffffffffff16815260200182602001602081019061248f9190613a81565b73ffffffffffffffffffffffffffffffffffffffff1681526020016124ba6060840160408501613a81565b73ffffffffffffffffffffffffffffffffffffffff1681526020016124e560a084016080850161436d565b15159052600960006124fd6080850160608601613b83565b67ffffffffffffffff16815260208082019290925260409081016000208351815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161783559385015160018301805491831691909516179093559083015160029091018054606094850151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090911692909316919091179190911790556125de9060808301908301613b83565b67ffffffffffffffff167fab9bd0e4888101232b8f09dae2952ff59a6eea4a19fbddf2a8ca7b23f0e4cb406126196040840160208501613a81565b6126296060850160408601613a81565b6126366020860186613a81565b61264660a087016080880161436d565b604051612688949392919073ffffffffffffffffffffffffffffffffffffffff9485168152928416602084015292166040820152901515606082015260800190565b60405180910390a250565b67ffffffffffffffff85166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052612758576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152602401610c96565b602081015181516040517f38314bb200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909216916338314bb2916127b4913090899060040161438a565b6020604051808303816000875af192505050801561280d575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261280a918101906143cc565b60015b612896573d80801561283b576040519150601f19603f3d011682016040523d82523d6000602084013e612840565b606091505b508667ffffffffffffffff168367ffffffffffffffff167fa481d91c3f9574c23ee84fef85246354b760a0527a535d6382354e4684703ce387846040516128889291906143e9565b60405180910390a350612903565b80156128ae576128a9868489888861329a565b6128fc565b8667ffffffffffffffff168367ffffffffffffffff167f8d3121fe961b40270f336aa75feb1213f1c979a33993311c60da4dd0f24526cf876040516128f391906139c6565b60405180910390a35b50506109fe565b612910858388878761329a565b505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526129f49085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152613481565b50505050565b801580612a9a57506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa158015612a74573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a9891906141a1565b155b612b26576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608401610c96565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052610bb49084907f095ea7b30000000000000000000000000000000000000000000000000000000090606401612972565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052610bb49084907fa9059cbb0000000000000000000000000000000000000000000000000000000090606401612972565b6000808a8a8a8a8a8a8a8a8a604051602001612bf69998979695949392919061440e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b6000612c8b838501856145ab565b8051516020820151519192509081158015612ca4575080155b15612cda576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82811015612d7e57612d7684600001518281518110612cfe57612cfe614206565b60200260200101516040015185600001518381518110612d2057612d20614206565b60200260200101516000015186600001518481518110612d4257612d42614206565b6020026020010151602001518888600001518681518110612d6557612d65614206565b602002602001015160600151612f21565b600101612cdd565b5060005b81811015612e2357612e1b84602001518281518110612da357612da3614206565b60200260200101516020015185602001518381518110612dc557612dc5614206565b60200260200101516000015186602001518481518110612de757612de7614206565b60200260200101516060015187602001518581518110612e0957612e09614206565b60200260200101516040015189612693565b600101612d82565b50505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612eab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610c96565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612f2b610737565b60085490915080821080612f47575085612f45828461471f565b105b15612f8f576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018790526024810183905260448101829052606401610c96565b67ffffffffffffffff87166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052613054576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff89166004820152602401610c96565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810189905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b1580156130c057600080fd5b505af11580156130d4573d6000803e3d6000fd5b505050602082015161311f915073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690896129fa565b6020810151604080830151835191517fa71d98b700000000000000000000000000000000000000000000000000000000815260009373ffffffffffffffffffffffffffffffffffffffff169263a71d98b7928b926131a6927f000000000000000000000000000000000000000000000000000000000000000092918f908d90600401614732565b60006040518083038185885af11580156131c4573d6000803e3d6000fd5b50505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261320b9190810190614779565b90508867ffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168767ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a885600001518c8a8760405161328794939291906147e7565b60405180910390a4505050505050505050565b8015613322577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0866040518263ffffffff1660e01b81526004016000604051808303818588803b15801561330857600080fd5b505af115801561331c573d6000803e3d6000fd5b50505050505b600b546133699073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691168761358d565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810187905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b1580156133d557600080fd5b505af11580156133e9573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168367ffffffffffffffff168567ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a83089876040518060200160405280600081525060405161347294939291906147e7565b60405180910390a45050505050565b60006134e3826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661368b9092919063ffffffff16565b805190915015610bb4578080602001905181019061350191906143cc565b610bb4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610c96565b6040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff8381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa158015613604573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061362891906141a1565b6136329190614336565b60405173ffffffffffffffffffffffffffffffffffffffff85166024820152604481018290529091506129f49085907f095ea7b30000000000000000000000000000000000000000000000000000000090606401612972565b606061369a84846000856136a2565b949350505050565b606082471015613734576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610c96565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161375d919061482e565b60006040518083038185875af1925050503d806000811461379a576040519150601f19603f3d011682016040523d82523d6000602084013e61379f565b606091505b50915091506137b0878383876137bb565b979650505050505050565b6060831561385157825160000361384a5773ffffffffffffffffffffffffffffffffffffffff85163b61384a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610c96565b508161369a565b61369a83838151156138665781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c9691906139c6565b828054828255906000526020600020908101928215613914579160200282015b8281111561391457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906138ba565b50613920929150613943565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b808211156139205760008155600101613944565b60005b8381101561397357818101518382015260200161395b565b50506000910152565b60008151808452613994816020860160208601613958565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006139d9602083018461397c565b9392505050565b6000602082840312156139f257600080fd5b5035919050565b600060a08284031215613a0b57600080fd5b50919050565b6020808252825182820181905260009190848201906040850190845b81811015613a5357835167ffffffffffffffff1683529284019291840191600101613a2d565b50909695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461085057600080fd5b600060208284031215613a9357600080fd5b81356139d981613a5f565b803567ffffffffffffffff81168114613ab657600080fd5b919050565b801515811461085057600080fd5b60008083601f840112613adb57600080fd5b50813567ffffffffffffffff811115613af357600080fd5b602083019150836020828501011115613b0b57600080fd5b9250929050565b600080600080600060808688031215613b2a57600080fd5b613b3386613a9e565b9450602086013593506040860135613b4a81613abb565b9250606086013567ffffffffffffffff811115613b6657600080fd5b613b7288828901613ac9565b969995985093965092949392505050565b600060208284031215613b9557600080fd5b6139d982613a9e565b600080600060608486031215613bb357600080fd5b8335613bbe81613a5f565b9250602084013591506040840135613bd581613a5f565b809150509250925092565b60008151808452602080850194506020840160005b83811015613c2757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613bf5565b509495945050505050565b6020815260006139d96020830184613be0565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613c9757613c97613c45565b60405290565b6040805190810167ffffffffffffffff81118282101715613c9757613c97613c45565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613d0757613d07613c45565b604052919050565b600067ffffffffffffffff821115613d2957613d29613c45565b5060051b60200190565b600082601f830112613d4457600080fd5b81356020613d59613d5483613d0f565b613cc0565b8083825260208201915060208460051b870101935086841115613d7b57600080fd5b602086015b84811015613da0578035613d9381613a5f565b8352918301918301613d80565b509695505050505050565b803560ff81168114613ab657600080fd5b600067ffffffffffffffff821115613dd657613dd6613c45565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613e1357600080fd5b8135613e21613d5482613dbc565b818152846020838601011115613e3657600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215613e6c57600080fd5b863567ffffffffffffffff80821115613e8457600080fd5b613e908a838b01613d33565b97506020890135915080821115613ea657600080fd5b613eb28a838b01613d33565b9650613ec060408a01613dab565b95506060890135915080821115613ed657600080fd5b613ee28a838b01613e02565b9450613ef060808a01613a9e565b935060a0890135915080821115613f0657600080fd5b50613f1389828a01613e02565b9150509295509295509295565b60008083601f840112613f3257600080fd5b50813567ffffffffffffffff811115613f4a57600080fd5b6020830191508360208260051b8501011115613b0b57600080fd5b60008060008060008060008060e0898b031215613f8157600080fd5b606089018a811115613f9257600080fd5b8998503567ffffffffffffffff80821115613fac57600080fd5b613fb88c838d01613ac9565b909950975060808b0135915080821115613fd157600080fd5b613fdd8c838d01613f20565b909750955060a08b0135915080821115613ff657600080fd5b506140038b828c01613f20565b999c989b50969995989497949560c00135949350505050565b6000806020838503121561402f57600080fd5b823567ffffffffffffffff8082111561404757600080fd5b818501915085601f83011261405b57600080fd5b81358181111561406a57600080fd5b86602060a08302850101111561407f57600080fd5b60209290920196919550909350505050565b600080604083850312156140a457600080fd5b8235915060208301356140b681613a5f565b809150509250929050565b6000806000806000608086880312156140d957600080fd5b6140e286613a9e565b94506020860135935060408601359250606086013567ffffffffffffffff811115613b6657600080fd5b602080825282518282018190526000919060409081850190868401855b82811015614194578151805173ffffffffffffffffffffffffffffffffffffffff90811686528782015181168887015286820151168686015260608082015167ffffffffffffffff169086015260809081015115159085015260a09093019290850190600101614129565b5091979650505050505050565b6000602082840312156141b357600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417614200576142006141ba565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff80831681810361427d5761427d6141ba565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526142b78184018a613be0565b905082810360808401526142cb8189613be0565b905060ff871660a084015282810360c08401526142e8818761397c565b905067ffffffffffffffff851660e084015282810361010084015261430d818561397c565b9c9b505050505050505050505050565b60ff8181168382160190811115614200576142006141ba565b80820180821115614200576142006141ba565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006020828403121561437f57600080fd5b81356139d981613abb565b600073ffffffffffffffffffffffffffffffffffffffff8086168352808516602084015250606060408301526143c3606083018461397c565b95945050505050565b6000602082840312156143de57600080fd5b81516139d981613abb565b6040815260006143fc604083018561397c565b82810360208401526143c3818561397c565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526144558285018b613be0565b91508382036080850152614469828a613be0565b915060ff881660a085015283820360c0850152614486828861397c565b90861660e0850152838103610100850152905061430d818561397c565b600082601f8301126144b457600080fd5b813560206144c4613d5483613d0f565b82815260059290921b840181019181810190868411156144e357600080fd5b8286015b84811015613da057803567ffffffffffffffff808211156145085760008081fd5b81890191506080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156145415760008081fd5b614549613c74565b878401358152604061455c818601613a9e565b8983015260608086013561456f81613abb565b8383015292850135928484111561458857600091508182fd5b6145968e8b86890101613e02565b908301525086525050509183019183016144e7565b600060208083850312156145be57600080fd5b823567ffffffffffffffff808211156145d657600080fd5b90840190604082870312156145ea57600080fd5b6145f2613c9d565b82358281111561460157600080fd5b8301601f8101881361461257600080fd5b8035614620613d5482613d0f565b81815260059190911b8201860190868101908a83111561463f57600080fd5b8784015b838110156146eb5780358781111561465a57600080fd5b85016080818e037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001121561468e57600080fd5b614696613c74565b8a820135815260408201358b8201526146b160608301613a9e565b60408201526080820135898111156146c95760008081fd5b6146d78f8d83860101613e02565b606083015250845250918801918801614643565b508452505050828401358281111561470257600080fd5b61470e888286016144a3565b948201949094529695505050505050565b81810381811115614200576142006141ba565b600073ffffffffffffffffffffffffffffffffffffffff8088168352808716602084015280861660408401525083606083015260a060808301526137b060a083018461397c565b60006020828403121561478b57600080fd5b815167ffffffffffffffff8111156147a257600080fd5b8201601f810184136147b357600080fd5b80516147c1613d5482613dbc565b8181528560208385010111156147d657600080fd5b6143c3826020830160208601613958565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260806040820152600061481c608083018561397c565b82810360608401526137b0818561397c565b60008251614840818460208701613958565b919091019291505056fea164736f6c6343000818000a",
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

func (_LiquidityManager *LiquidityManagerTransactor) WithdrawERC20(opts *bind.TransactOpts, token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _LiquidityManager.contract.Transact(opts, "withdrawERC20", token, amount, destination)
}

func (_LiquidityManager *LiquidityManagerSession) WithdrawERC20(token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.WithdrawERC20(&_LiquidityManager.TransactOpts, token, amount, destination)
}

func (_LiquidityManager *LiquidityManagerTransactorSession) WithdrawERC20(token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _LiquidityManager.Contract.WithdrawERC20(&_LiquidityManager.TransactOpts, token, amount, destination)
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

	WithdrawERC20(opts *bind.TransactOpts, token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error)

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
