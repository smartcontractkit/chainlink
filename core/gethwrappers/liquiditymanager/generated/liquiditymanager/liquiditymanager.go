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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidRemoteChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"CrossChainRebalancerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"FinalizationFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"}],\"name\":\"FinalizationStepCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAddedToContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newLiquidityContainer\",\"type\":\"address\"}],\"name\":\"LiquidityContainerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemovedFromContainer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"fromChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"toChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeSpecificData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bridgeReturnData\",\"type\":\"bytes\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"MinimumLiquiditySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"NativeDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"NativeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainRebalancers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getCrossChainRebalancer\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityManager.CrossChainRebalancer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLiquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalLiquidityContainer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedDestChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_localToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nativeBridgeFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"rebalanceLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"shouldWrapNative\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"bridgeSpecificPayload\",\"type\":\"bytes\"}],\"name\":\"receiveLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs\",\"name\":\"crossChainLiqManager\",\"type\":\"tuple\"}],\"name\":\"setCrossChainRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteRebalancer\",\"type\":\"address\"},{\"internalType\":\"contractIBridgeAdapter\",\"name\":\"localBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainRebalancerArgs[]\",\"name\":\"crossChainRebalancers\",\"type\":\"tuple[]\"}],\"name\":\"setCrossChainRebalancers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractILiquidityContainer\",\"name\":\"localLiquidityContainer\",\"type\":\"address\"}],\"name\":\"setLocalLiquidityContainer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minimumLiquidity\",\"type\":\"uint256\"}],\"name\":\"setMinimumLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620048573803806200485783398101604081905262000034916200022d565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000169565b505046608052506001600160401b038316600003620000f05760405163f89d762960e01b815260040160405180910390fd5b6001600160a01b03841615806200010e57506001600160a01b038216155b156200012d5760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0393841660a0526001600160401b039290921660c052600b80546001600160a01b031916919093161790915560085562000292565b336001600160a01b03821603620001c35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200022a57600080fd5b50565b600080600080608085870312156200024457600080fd5b8451620002518162000214565b60208601519094506001600160401b03811681146200026f57600080fd5b6040860151909350620002828162000214565b6060959095015193969295505050565b60805160a05160c05161454f6200030860003960008181612f540152613129015260008181610424015281816106fd015281816108d0015281816109160152818161152901528181612e3201528181612eb301528181612fdc015261307a0152600081816116c60152611712015261454f6000f3fe60806040526004361061019a5760003560e01c8063791781f5116100e1578063afcb95d71161008a578063b8ca8dd811610064578063b8ca8dd814610630578063f2fde38b14610650578063f8c2d8fa14610670578063fe65d5af1461069057600080fd5b8063afcb95d7146105aa578063b1dc65a4146105f0578063b7e7fa051461061057600080fd5b806383d34afe116100bb57806383d34afe1461054a5780638da5cb5b1461055f5780639c8f9f231461058a57600080fd5b8063791781f5146104cd57806379ba5097146104f857806381ff70481461050d57600080fd5b806351c6590a11610143578063666cab8d1161011d578063666cab8d1461046b5780636a11ee901461048d578063706bf645146104ad57600080fd5b806351c6590a146102e0578063568446e7146103005780636511d9191461041257600080fd5b80633275636e116101745780633275636e1461027e578063348759c11461029e57806350a197d7146102c057600080fd5b80630910a510146101de578063181f5a7714610206578063282567b41461025c57600080fd5b366101d957604080513481523360208201527f3c597f6ac9fe7f0ed6da50b07618f5850a642e459ad587f7fab491a71f8b0ab8910160405180910390a1005b600080fd5b3480156101ea57600080fd5b506101f36106b2565b6040519081526020015b60405180910390f35b34801561021257600080fd5b5061024f6040518060400160405280601a81526020017f4c69717569646974794d616e6167657220312e302e302d64657600000000000081525081565b6040516101fd9190613700565b34801561026857600080fd5b5061027c61027736600461371a565b61076d565b005b34801561028a57600080fd5b5061027c610299366004613733565b6107ba565b3480156102aa57600080fd5b506102b36107ce565b6040516101fd919061374b565b3480156102cc57600080fd5b5061027c6102db36600461380d565b61085a565b3480156102ec57600080fd5b5061027c6102fb36600461371a565b6108b6565b34801561030c57600080fd5b506103c361031b36600461387e565b6040805160808101825260008082526020820181905291810182905260608101919091525067ffffffffffffffff166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff161515606082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff90811682526020808501518216908301528383015116918101919091526060918201511515918101919091526080016101fd565b34801561041e57600080fd5b506104467f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fd565b34801561047757600080fd5b506104806109f3565b6040516101fd91906138eb565b34801561049957600080fd5b5061027c6104a8366004613b2e565b610a61565b3480156104b957600080fd5b5061027c6104c8366004613bfb565b61126b565b3480156104d957600080fd5b50600b5473ffffffffffffffffffffffffffffffffffffffff16610446565b34801561050457600080fd5b5061027c61132f565b34801561051957600080fd5b506004546002546040805163ffffffff808516825264010000000090940490931660208401528201526060016101fd565b34801561055657600080fd5b506008546101f3565b34801561056b57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610446565b34801561059657600080fd5b5061027c6105a536600461371a565b61142c565b3480156105b657600080fd5b50600254600454604080516001815260208101939093526801000000000000000090910467ffffffffffffffff16908201526060016101fd565b3480156105fc57600080fd5b5061027c61060b366004613c5d565b611583565b34801561061c57600080fd5b5061027c61062b366004613d14565b611bf4565b34801561063c57600080fd5b5061027c61064b366004613d89565b611c34565b34801561065c57600080fd5b5061027c61066b366004613bfb565b611d29565b34801561067c57600080fd5b5061027c61068b366004613db9565b611d3a565b34801561069c57600080fd5b506106a5611d8d565b6040516101fd9190613e04565b600b546040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526000917f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa158015610744573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107689190613e99565b905090565b610775611f4e565b600880549082905560408051828152602081018490527ff97e758c8b3d81df7b0e1b7327a6a7fcf09a41536b2d274b9103015d715f11eb910160405180910390a15050565b6107c2611f4e565b6107cb81611fd1565b50565b6060600a80548060200260200160405190810160405280929190818152602001828054801561085057602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff168152602001906008019060208260070104928301926001038202915080841161080b5790505b5050505050905090565b610862611f4e565b6108af858584848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525089925067ffffffffffffffff91506123cd9050565b5050505050565b6108f873ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333084612652565b600b5461093f9073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911683612734565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b1580156109ab57600080fd5b505af11580156109bf573d6000803e3d6000fd5b50506040518392503391507f5414b81d05ac3542606f164e16a9a107d05d21e906539cc5ceb61d7b6b707eb590600090a350565b6060600780548060200260200160405190810160405280929190818152602001828054801561085057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610a2d575050505050905090565b855185518560ff16601f831115610ad9576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b80600003610b43576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610ad0565b818314610bd1576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610ad0565b610bdc816003613ee1565b8311610c44576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610ad0565b610c4c611f4e565b60065460005b81811015610d40576005600060068381548110610c7157610c71613efe565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560078054600592919084908110610ce157610ce1613efe565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600101610c52565b50895160005b818110156111135760008c8281518110610d6257610d62613efe565b6020026020010151905060006002811115610d7f57610d7f613f2d565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff166002811115610dbe57610dbe613f2d565b14610e25576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610ad0565b73ffffffffffffffffffffffffffffffffffffffff8116610e72576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115610f2257610f22613f2d565b021790555090505060008c8381518110610f3e57610f3e613efe565b6020026020010151905060006002811115610f5b57610f5b613f2d565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff166002811115610f9a57610f9a613f2d565b14611001576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610ad0565b73ffffffffffffffffffffffffffffffffffffffff811661104e576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156110fe576110fe613f2d565b02179055509050505050806001019050610d46565b508a516111279060069060208e01906135d4565b50895161113b9060079060208d01906135d4565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c1617179055600480546111c19146913091906000906111939063ffffffff16613f5c565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e6128b6565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f60405161125599989796959493929190613f7f565b60405180910390a1505050505050505050505050565b611273611f4e565b73ffffffffffffffffffffffffffffffffffffffff81166112c0576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517f07dc474694ac40123aadcd2445f1b38d2eb353edd9319dcea043548ab34990ec90600090a250565b60015473ffffffffffffffffffffffffffffffffffffffff1633146113b0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610ad0565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611434611f4e565b600061143e6106b2565b90508181101561148b576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018390526024810182905260006044820152606401610ad0565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b1580156114f757600080fd5b505af115801561150b573d6000803e3d6000fd5b5061155292505073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690503384612961565b604051829033907f2bda316674f8d73d289689d7a3acdf8e353b7a142fb5a68ac2aa475104039c1890600090a35050565b60045460208901359067ffffffffffffffff6801000000000000000090910481169082161161160657600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff80851693820193909352680100000000000000009091049091166024820152604401610ad0565b6116118888836129b7565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a359182146116c35780516040517f93df584c000000000000000000000000000000000000000000000000000000008152600481019190915260248101839052604401610ad0565b467f000000000000000000000000000000000000000000000000000000000000000014611744576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610ad0565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a16020810151611796906001614015565b60ff1687146117d1576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b86851461180a576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561184d5761184d613f2d565b600281111561185e5761185e613f2d565b905250905060028160200151600281111561187b5761187b613f2d565b1480156118c257506007816000015160ff168154811061189d5761189d613efe565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b6118f8576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000611906866020613ee1565b611911896020613ee1565b61191d8c61014461402e565b611927919061402e565b611931919061402e565b9050368114611975576040517f8e1192e100000000000000000000000000000000000000000000000000000000815260048101829052366024820152604401610ad0565b5060008a8a604051611988929190614041565b60405190819003812061199f918e90602001614051565b6040516020818303038152906040528051906020012090506119bf61365e565b8860005b81811015611be35760006001858a84602081106119e2576119e2613efe565b6119ef91901a601b614015565b8f8f86818110611a0157611a01613efe565b905060200201358e8e87818110611a1a57611a1a613efe565b9050602002013560405160008152602001604052604051611a57949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611a79573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff8082168652939750919550929392840191610100909104166002811115611afc57611afc613f2d565b6002811115611b0d57611b0d613f2d565b9052509050600181602001516002811115611b2a57611b2a613f2d565b14611b61576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f8110611b7857611b78613efe565b602002015115611bb4576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f8110611bcf57611bcf613efe565b9115156020909202015250506001016119c3565b505050505050505050505050505050565b611bfc611f4e565b60005b81811015611c2f57611c27838383818110611c1c57611c1c613efe565b905060a00201611fd1565b600101611bff565b505050565b611c3c611f4e565b60008173ffffffffffffffffffffffffffffffffffffffff168360405160006040518083038185875af1925050503d8060008114611c96576040519150601f19603f3d011682016040523d82523d6000602084013e611c9b565b606091505b5050905080611cd6576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805184815273ffffffffffffffffffffffffffffffffffffffff841660208201527f6b84d241b711af111ecfa0e518239e6ca212da442a76548fe8a1f4e77518256a910160405180910390a1505050565b611d31611f4e565b6107cb81612b66565b611d42611f4e565b6108af85858567ffffffffffffffff86868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612c5b92505050565b600a5460609060008167ffffffffffffffff811115611dae57611dae6138fe565b604051908082528060200260200182016040528015611e2557816020015b6040805160a0810182526000808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181611dcc5790505b50905060005b82811015611f47576000600a8281548110611e4857611e48613efe565b6000918252602080832060048304015460039092166008026101000a90910467ffffffffffffffff1680835260098252604092839020835160808082018652825473ffffffffffffffffffffffffffffffffffffffff9081168352600184015481168387019081526002909401548082168489019081527401000000000000000000000000000000000000000090910460ff1615156060808601918252895160a081018b52865185168152965184169887019890985290519091169684019690965293820183905293511515928101929092528551909350859085908110611f3257611f32613efe565b60209081029190910101525050600101611e2b565b5092915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611fcf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610ad0565b565b611fe1608082016060830161387e565b67ffffffffffffffff16600003612024576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006120336020830183613bfb565b73ffffffffffffffffffffffffffffffffffffffff16148061207a575060006120626040830160208401613bfb565b73ffffffffffffffffffffffffffffffffffffffff16145b806120aa575060006120926060830160408401613bfb565b73ffffffffffffffffffffffffffffffffffffffff16145b156120e1576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006009816120f6608085016060860161387e565b67ffffffffffffffff16815260208101919091526040016000206002015473ffffffffffffffffffffffffffffffffffffffff160361218257600a612141608083016060840161387e565b8154600181018355600092835260209092206004830401805460039093166008026101000a67ffffffffffffffff8181021990941692909316929092021790555b60408051608081019091528061219b6020840184613bfb565b73ffffffffffffffffffffffffffffffffffffffff1681526020018260200160208101906121c99190613bfb565b73ffffffffffffffffffffffffffffffffffffffff1681526020016121f46060840160408501613bfb565b73ffffffffffffffffffffffffffffffffffffffff16815260200161221f60a0840160808501614065565b1515905260096000612237608085016060860161387e565b67ffffffffffffffff16815260208082019290925260409081016000208351815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161783559385015160018301805491831691909516179093559083015160029091018054606094850151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff0000000000000000000000000000000000000000009091169290931691909117919091179055612318906080830190830161387e565b67ffffffffffffffff167fab9bd0e4888101232b8f09dae2952ff59a6eea4a19fbddf2a8ca7b23f0e4cb406123536040840160208501613bfb565b6123636060850160408601613bfb565b6123706020860186613bfb565b61238060a0870160808801614065565b6040516123c2949392919073ffffffffffffffffffffffffffffffffffffffff9485168152928416602084015292166040820152901515606082015260800190565b60405180910390a250565b67ffffffffffffffff85166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052612492576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152602401610ad0565b602081015181516040517f38314bb200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909216916338314bb2916124ee9130908990600401614082565b6020604051808303816000875af1925050508015612547575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201909252612544918101906140c4565b60015b6125d0573d808015612575576040519150601f19603f3d011682016040523d82523d6000602084013e61257a565b606091505b508667ffffffffffffffff168367ffffffffffffffff167fa481d91c3f9574c23ee84fef85246354b760a0527a535d6382354e4684703ce387846040516125c29291906140e1565b60405180910390a35061263d565b80156125e8576125e38684898888612fd4565b612636565b8667ffffffffffffffff168367ffffffffffffffff167f8d3121fe961b40270f336aa75feb1213f1c979a33993311c60da4dd0f24526cf8760405161262d9190613700565b60405180910390a35b50506108af565b61264a8583888787612fd4565b505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff8085166024830152831660448201526064810182905261272e9085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526131bb565b50505050565b8015806127d457506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa1580156127ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127d29190613e99565b155b612860576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608401610ad0565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052611c2f9084907f095ea7b300000000000000000000000000000000000000000000000000000000906064016126ac565b6000808a8a8a8a8a8a8a8a8a6040516020016128da99989796959493929190614106565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052611c2f9084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064016126ac565b60006129c5838501856142a3565b80515160208201515191925090811580156129de575080155b15612a14576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82811015612ab857612ab084600001518281518110612a3857612a38613efe565b60200260200101516040015185600001518381518110612a5a57612a5a613efe565b60200260200101516000015186600001518481518110612a7c57612a7c613efe565b6020026020010151602001518888600001518681518110612a9f57612a9f613efe565b602002602001015160600151612c5b565b600101612a17565b5060005b81811015612b5d57612b5584602001518281518110612add57612add613efe565b60200260200101516020015185602001518381518110612aff57612aff613efe565b60200260200101516000015186602001518481518110612b2157612b21613efe565b60200260200101516060015187602001518581518110612b4357612b43613efe565b602002602001015160400151896123cd565b600101612abc565b50505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612be5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ad0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612c656106b2565b60085490915080821080612c81575085612c7f8284614417565b105b15612cc9576040517fd54d0fc4000000000000000000000000000000000000000000000000000000008152600481018790526024810183905260448101829052606401610ad0565b67ffffffffffffffff87166000908152600960209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301548116938201939093526002909101549182169281019290925274010000000000000000000000000000000000000000900460ff16151560608201819052612d8e576040517fc9ff038f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff89166004820152602401610ad0565b600b546040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810189905273ffffffffffffffffffffffffffffffffffffffff90911690630a861f2a90602401600060405180830381600087803b158015612dfa57600080fd5b505af1158015612e0e573d6000803e3d6000fd5b5050506020820151612e59915073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169089612734565b6020810151604080830151835191517fa71d98b700000000000000000000000000000000000000000000000000000000815260009373ffffffffffffffffffffffffffffffffffffffff169263a71d98b7928b92612ee0927f000000000000000000000000000000000000000000000000000000000000000092918f908d9060040161442a565b60006040518083038185885af1158015612efe573d6000803e3d6000fd5b50505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612f459190810190614471565b90508867ffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168767ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a885600001518c8a87604051612fc194939291906144df565b60405180910390a4505050505050505050565b801561305c577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0866040518263ffffffff1660e01b81526004016000604051808303818588803b15801561304257600080fd5b505af1158015613056573d6000803e3d6000fd5b50505050505b600b546130a39073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169116876132c7565b600b546040517feb521a4c0000000000000000000000000000000000000000000000000000000081526004810187905273ffffffffffffffffffffffffffffffffffffffff9091169063eb521a4c90602401600060405180830381600087803b15801561310f57600080fd5b505af1158015613123573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168367ffffffffffffffff168567ffffffffffffffff167f2a0b69eaf1b415ca57005b4f87582ddefc6d960325ff30dc62a9b3e1e1e5b8a8308987604051806020016040528060008152506040516131ac94939291906144df565b60405180910390a45050505050565b600061321d826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166133c59092919063ffffffff16565b805190915015611c2f578080602001905181019061323b91906140c4565b611c2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610ad0565b6040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff8381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa15801561333e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133629190613e99565b61336c919061402e565b60405173ffffffffffffffffffffffffffffffffffffffff851660248201526044810182905290915061272e9085907f095ea7b300000000000000000000000000000000000000000000000000000000906064016126ac565b60606133d484846000856133dc565b949350505050565b60608247101561346e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610ad0565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516134979190614526565b60006040518083038185875af1925050503d80600081146134d4576040519150601f19603f3d011682016040523d82523d6000602084013e6134d9565b606091505b50915091506134ea878383876134f5565b979650505050505050565b6060831561358b5782516000036135845773ffffffffffffffffffffffffffffffffffffffff85163b613584576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610ad0565b50816133d4565b6133d483838151156135a05781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ad09190613700565b82805482825590600052602060002090810192821561364e579160200282015b8281111561364e57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906135f4565b5061365a92915061367d565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b8082111561365a576000815560010161367e565b60005b838110156136ad578181015183820152602001613695565b50506000910152565b600081518084526136ce816020860160208601613692565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061371360208301846136b6565b9392505050565b60006020828403121561372c57600080fd5b5035919050565b600060a0828403121561374557600080fd5b50919050565b6020808252825182820181905260009190848201906040850190845b8181101561378d57835167ffffffffffffffff1683529284019291840191600101613767565b50909695505050505050565b803567ffffffffffffffff811681146137b157600080fd5b919050565b80151581146107cb57600080fd5b60008083601f8401126137d657600080fd5b50813567ffffffffffffffff8111156137ee57600080fd5b60208301915083602082850101111561380657600080fd5b9250929050565b60008060008060006080868803121561382557600080fd5b61382e86613799565b9450602086013593506040860135613845816137b6565b9250606086013567ffffffffffffffff81111561386157600080fd5b61386d888289016137c4565b969995985093965092949392505050565b60006020828403121561389057600080fd5b61371382613799565b60008151808452602080850194506020840160005b838110156138e057815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016138ae565b509495945050505050565b6020815260006137136020830184613899565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613950576139506138fe565b60405290565b6040805190810167ffffffffffffffff81118282101715613950576139506138fe565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156139c0576139c06138fe565b604052919050565b600067ffffffffffffffff8211156139e2576139e26138fe565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff811681146107cb57600080fd5b600082601f830112613a1f57600080fd5b81356020613a34613a2f836139c8565b613979565b8083825260208201915060208460051b870101935086841115613a5657600080fd5b602086015b84811015613a7b578035613a6e816139ec565b8352918301918301613a5b565b509695505050505050565b803560ff811681146137b157600080fd5b600067ffffffffffffffff821115613ab157613ab16138fe565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613aee57600080fd5b8135613afc613a2f82613a97565b818152846020838601011115613b1157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215613b4757600080fd5b863567ffffffffffffffff80821115613b5f57600080fd5b613b6b8a838b01613a0e565b97506020890135915080821115613b8157600080fd5b613b8d8a838b01613a0e565b9650613b9b60408a01613a86565b95506060890135915080821115613bb157600080fd5b613bbd8a838b01613add565b9450613bcb60808a01613799565b935060a0890135915080821115613be157600080fd5b50613bee89828a01613add565b9150509295509295509295565b600060208284031215613c0d57600080fd5b8135613713816139ec565b60008083601f840112613c2a57600080fd5b50813567ffffffffffffffff811115613c4257600080fd5b6020830191508360208260051b850101111561380657600080fd5b60008060008060008060008060e0898b031215613c7957600080fd5b606089018a811115613c8a57600080fd5b8998503567ffffffffffffffff80821115613ca457600080fd5b613cb08c838d016137c4565b909950975060808b0135915080821115613cc957600080fd5b613cd58c838d01613c18565b909750955060a08b0135915080821115613cee57600080fd5b50613cfb8b828c01613c18565b999c989b50969995989497949560c00135949350505050565b60008060208385031215613d2757600080fd5b823567ffffffffffffffff80821115613d3f57600080fd5b818501915085601f830112613d5357600080fd5b813581811115613d6257600080fd5b86602060a083028501011115613d7757600080fd5b60209290920196919550909350505050565b60008060408385031215613d9c57600080fd5b823591506020830135613dae816139ec565b809150509250929050565b600080600080600060808688031215613dd157600080fd5b613dda86613799565b94506020860135935060408601359250606086013567ffffffffffffffff81111561386157600080fd5b602080825282518282018190526000919060409081850190868401855b82811015613e8c578151805173ffffffffffffffffffffffffffffffffffffffff90811686528782015181168887015286820151168686015260608082015167ffffffffffffffff169086015260809081015115159085015260a09093019290850190600101613e21565b5091979650505050505050565b600060208284031215613eab57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417613ef857613ef8613eb2565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff808316818103613f7557613f75613eb2565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152613faf8184018a613899565b90508281036080840152613fc38189613899565b905060ff871660a084015282810360c0840152613fe081876136b6565b905067ffffffffffffffff851660e084015282810361010084015261400581856136b6565b9c9b505050505050505050505050565b60ff8181168382160190811115613ef857613ef8613eb2565b80820180821115613ef857613ef8613eb2565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006020828403121561407757600080fd5b8135613713816137b6565b600073ffffffffffffffffffffffffffffffffffffffff8086168352808516602084015250606060408301526140bb60608301846136b6565b95945050505050565b6000602082840312156140d657600080fd5b8151613713816137b6565b6040815260006140f460408301856136b6565b82810360208401526140bb81856136b6565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b16604085015281606085015261414d8285018b613899565b91508382036080850152614161828a613899565b915060ff881660a085015283820360c085015261417e82886136b6565b90861660e0850152838103610100850152905061400581856136b6565b600082601f8301126141ac57600080fd5b813560206141bc613a2f836139c8565b82815260059290921b840181019181810190868411156141db57600080fd5b8286015b84811015613a7b57803567ffffffffffffffff808211156142005760008081fd5b81890191506080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156142395760008081fd5b61424161392d565b8784013581526040614254818601613799565b89830152606080860135614267816137b6565b8383015292850135928484111561428057600091508182fd5b61428e8e8b86890101613add565b908301525086525050509183019183016141df565b600060208083850312156142b657600080fd5b823567ffffffffffffffff808211156142ce57600080fd5b90840190604082870312156142e257600080fd5b6142ea613956565b8235828111156142f957600080fd5b8301601f8101881361430a57600080fd5b8035614318613a2f826139c8565b81815260059190911b8201860190868101908a83111561433757600080fd5b8784015b838110156143e35780358781111561435257600080fd5b85016080818e037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001121561438657600080fd5b61438e61392d565b8a820135815260408201358b8201526143a960608301613799565b60408201526080820135898111156143c15760008081fd5b6143cf8f8d83860101613add565b60608301525084525091880191880161433b565b50845250505082840135828111156143fa57600080fd5b6144068882860161419b565b948201949094529695505050505050565b81810381811115613ef857613ef8613eb2565b600073ffffffffffffffffffffffffffffffffffffffff8088168352808716602084015280861660408401525083606083015260a060808301526134ea60a08301846136b6565b60006020828403121561448357600080fd5b815167ffffffffffffffff81111561449a57600080fd5b8201601f810184136144ab57600080fd5b80516144b9613a2f82613a97565b8181528560208385010111156144ce57600080fd5b6140bb826020830160208601613692565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260806040820152600061451460808301856136b6565b82810360608401526134ea81856136b6565b60008251614538818460208701613692565b919091019291505056fea164736f6c6343000818000a",
}

var LiquidityManagerABI = LiquidityManagerMetaData.ABI

var LiquidityManagerBin = LiquidityManagerMetaData.Bin

func DeployLiquidityManager(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localChainSelector uint64, localLiquidityContainer common.Address, minimumLiquidity *big.Int) (common.Address, *types.Transaction, *LiquidityManager, error) {
	parsed, err := LiquidityManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidityManagerBin), backend, token, localChainSelector, localLiquidityContainer, minimumLiquidity)
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
