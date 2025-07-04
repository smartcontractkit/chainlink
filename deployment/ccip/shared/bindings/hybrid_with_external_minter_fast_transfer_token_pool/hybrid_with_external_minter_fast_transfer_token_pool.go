// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hybrid_with_external_minter_fast_transfer_token_pool

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
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated"
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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type FastTransferTokenPoolAbstractDestChainConfig struct {
	MaxFillAmountPerRequest  *big.Int
	FillerAllowlistEnabled   bool
	FastTransferFillerFeeBps uint16
	FastTransferPoolFeeBps   uint16
	SettlementOverheadGas    uint32
	DestinationPool          []byte
	CustomExtraArgs          []byte
}

type FastTransferTokenPoolAbstractDestChainConfigUpdateArgs struct {
	FillerAllowlistEnabled   bool
	FastTransferFillerFeeBps uint16
	FastTransferPoolFeeBps   uint16
	SettlementOverheadGas    uint32
	RemoteChainSelector      uint64
	ChainFamilySelector      [4]byte
	MaxFillAmountPerRequest  *big.Int
	DestinationPool          []byte
	CustomExtraArgs          []byte
}

type FastTransferTokenPoolAbstractFillInfo struct {
	State  uint8
	Filler common.Address
}

type HybridWithExternalMinterFastTransferTokenPoolGroupUpdate struct {
	RemoteChainSelector uint64
	Group               uint8
	RemoteChainSupply   *big.Int
}

type IFastTransferPoolQuote struct {
	CcipSettlementFee *big.Int
	FastTransferFee   *big.Int
}

type PoolLockOrBurnInV1 struct {
	Receiver            []byte
	RemoteChainSelector uint64
	OriginalSender      common.Address
	Amount              *big.Int
	LocalToken          common.Address
}

type PoolLockOrBurnOutV1 struct {
	DestTokenAddress []byte
	DestPoolData     []byte
}

type PoolReleaseOrMintInV1 struct {
	OriginalSender          []byte
	RemoteChainSelector     uint64
	Receiver                common.Address
	SourceDenominatedAmount *big.Int
	LocalToken              common.Address
	SourcePoolAddress       []byte
	SourcePoolData          []byte
	OffchainTokenData       []byte
}

type PoolReleaseOrMintOutV1 struct {
	DestinationAmount *big.Int
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

type TokenPoolChainUpdate struct {
	RemoteChainSelector       uint64
	RemotePoolAddresses       [][]byte
	RemoteTokenAddress        []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var HybridWithExternalMinterFastTransferTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSendToken\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"computeFillId\",\"inputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"fastFill\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAccumulatedPoolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedFillers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCcipSendTokenFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIFastTransferPool.Quote\",\"components\":[{\"name\":\"ccipSettlementFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfig\",\"components\":[{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFillInfo\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.FillInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIFastTransferPool.FillState\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGroup\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumHybridWithExternalMinterFastTransferTokenPool.Group\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockedTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRebalancer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAllowedFiller\",\"inputs\":[{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"provideLiquidity\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeLiquidity\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRebalancer\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateDestChainConfig\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfigUpdateArgs[]\",\"components\":[{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFillerAllowList\",\"inputs\":[{\"name\":\"fillersToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"fillersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateGroups\",\"inputs\":[{\"name\":\"groupUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structHybridWithExternalMinterFastTransferTokenPool.GroupUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"internalType\":\"enumHybridWithExternalMinterFastTransferTokenPool.Group\"},{\"name\":\"remoteChainSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPoolFees\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestinationPoolUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destinationPool\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"filler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferRequested\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fillerReimbursementAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFeeAccumulated\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"prevState\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIFastTransferPool.FillState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FillerAllowListUpdated\",\"inputs\":[{\"name\":\"addFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"removeFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GroupUpdated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumHybridWithExternalMinterFastTransferTokenPool.Group\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityAdded\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityMigrated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumHybridWithExternalMinterFastTransferTokenPool.Group\"},{\"name\":\"remoteChainSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityRemoved\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolFeeWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RebalancerSet\",\"inputs\":[{\"name\":\"oldRebalancer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRebalancer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadySettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"FillerNotAllowlisted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientLiquidity\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientPoolFees\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFillId\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGroupUpdate\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"internalType\":\"enumHybridWithExternalMinterFastTransferTokenPool.Group\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LiquidityAmountCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"QuoteFeeExceedsUserMaxLimit\",\"inputs\":[{\"name\":\"quoteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferAmountExceedsMaxFillAmount\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461046a576166c2803803809161001d828561052f565b8339810160a08282031261046a5761003482610552565b9061004160208401610566565b60408401516001600160401b03811161046a5784019180601f8401121561046a578251926001600160401b038411610519578360051b90602082019461008a604051968761052f565b855260208086019282010192831161046a57602001905b828210610501575050506100c360806100bc60608701610552565b9501610552565b6040516321df0da760e01b81526001600160a01b03909416949093602081600481895afa9081156104f5576000916104bb575b506001600160a01b03169033156104aa57600180546001600160a01b0319163317905581158015610499575b8015610488575b610477578160209160049360805260c0526040519283809263313ce56760e01b82525afa60009181610436575b5061040b575b5060a052600480546001600160a01b0319166001600160a01b0384169081179091558151151560e08190529091906102e8575b50156102d2576101005261012052604051615fad908161071582396080518181816103d101528181610de1015281816116eb01528181611748015281816118a001528181611a4d015281816124e001528181612a12015281816131680152818161332a01528181613416015281816139210152818161396e01528181613d2a01528181614a8001528181615133015281816152f20152615edb015260a051818181611a91015281816136c3015281816138d701528181613bea01528181613e3901528181614d120152614d7c015260c051818181610f5801528181611908015281816127f1015281816131d1015281816135c20152613af4015260e051818181610f1301528181612f460152615c8a01526101005181613ef801526101205181818161027001528181610ced01528181610dba01528181614ea301526152ca0152f35b6335fdcccd60e21b600052600060045260246000fd5b602091604051916102f9848461052f565b60008352600036813760e051156103fa5760005b8351811015610374576001906001600160a01b0361032b8287610574565b511686610337826105b6565b610344575b50500161030d565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a1388661033c565b5091509260005b82518110156103ef576001906001600160a01b036103998286610574565b511680156103e957856103ab826106b4565b6103b9575b50505b0161037b565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a138856103b0565b506103b3565b50929150503861018f565b6335f4a7b360e01b60005260046000fd5b60ff1660ff821681810361041f575061015c565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d60201161046f575b816104526020938361052f565b8101031261046a5761046390610566565b9038610156565b600080fd5b3d9150610445565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03811615610129565b506001600160a01b03851615610122565b639b15e16f60e01b60005260046000fd5b90506020813d6020116104ed575b816104d66020938361052f565b8101031261046a576104e790610552565b386100f6565b3d91506104c9565b6040513d6000823e3d90fd5b6020809161050e84610552565b8152019101906100a1565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b0382119082101761051957604052565b51906001600160a01b038216820361046a57565b519060ff8216820361046a57565b80518210156105885760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156105885760005260206000200190600090565b60008181526003602052604090205480156106ad5760001981018181116106975760025460001981019190821161069757818103610646575b5050506002548015610630576000190161060a81600261059e565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b61067f61065761066893600261059e565b90549060031b1c928392600261059e565b819391549060031b91821b91600019901b19161790565b905560005260036020526040600020553880806105ef565b634e487b7160e01b600052601160045260246000fd5b5050600090565b8060005260036020526040600020541560001461070e5760025468010000000000000000811015610519576106f5610668826001859401600255600261059e565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a71461404c57508063055befd414613a18578063181f5a771461399257806321df0da71461394e578063240028e8146138fb57806324f65ee7146138bd5780632b2c0eb41461389f5780632e7aa8c8146134c8578063319ac101146134845780633317bbcc146133e257806339077537146130df578063432a6ba3146130b85780634c5ef0ed1461307457806354c8a4f314612f1457806362ddd3c414612eab5780636609f59914612e8f5780636cfd155314612dfa5780636d3d1a5814612dd35780636def4ce714612c9057806378b410f214612c5657806379ba509714612bb05780637d54534e14612b3b57806385572ffb146125c157806387f060d01461232b5780638926f54f146122fb5780638a18dcbd14611e895780638da5cb5b14611e62578063929ea5ba14611d5a578063962d402014611c215780639a4575b9146118655780639c8f9f231461171a5780639fe280f514611687578063a42a7b8b14611553578063a7cd63b7146114e5578063abe1c1e814611476578063acfecf911461136b578063af58d59f14611322578063b0f479a1146112fb578063b7946580146112c3578063c0d786551461122a578063c4bffe2b1461111a578063c75eea9c1461107b578063cf7401f314610f7c578063dc0bd97114610f38578063e0351e1314610efb578063e7e62f8514610b92578063e8a1da171461047b578063eb521a4c1461039d578063eeebc67414610346578063f2fde38b146102995763f36675171461025057600080fd5b346102945760003660031901126102945760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b34610294576020366003190112610294576001600160a01b036102ba6141d4565b6102c2614eed565b1633811461031c578073ffffffffffffffffffffffffffffffffffffffff1960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102945760803660031901126102945760443560ff81168103610294576064356001600160401b0381116102945760209161038961039592369060040161434d565b906024356004356149f1565b604051908152f35b34610294576020366003190112610294576004358015610451576001600160a01b03600e54163303610423576103f58130337f0000000000000000000000000000000000000000000000000000000000000000614acc565b6040519081527fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208860203392a2005b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b7fa90c0d190000000000000000000000000000000000000000000000000000000060005260046000fd5b34610294576104893661439b565b919092610494614eed565b6000905b828210610a075750505060009063ffffffff42165b8183106104b657005b6104c1838386614850565b926101208436031261029457604051936104da85614253565b6104e381614193565b855260208101356001600160401b0381116102945781019336601f860112156102945784356105118161449a565b9561051f60405197886142da565b81875260208088019260051b820101903682116102945760208101925b8284106109d957505050506020860194855260408201356001600160401b0381116102945761056e903690840161434d565b90604087019182526105986105863660608601614571565b936060890194855260c0369101614571565b94608088019586526105aa84516153fb565b6105b486516153fb565b825151156109af576105cf6001600160401b03895116615981565b15610977576001600160401b03885116600052600760205260406000206106bf85516001600160801b03604082015116906106926001600160801b036020830151169151151583608060405161062481614253565b858152602081018a905260408101849052606081018690520152855460ff60a01b91151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166001600160801b0384161763ffffffff60801b608089901b1617178555565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176001830155565b61079487516001600160801b03604082015116906107676001600160801b03602083015116915115158360806040516106f781614253565b858152602081018a9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166001600160801b0385161763ffffffff60801b60808a901b161791151560a01b60ff60a01b16919091179055565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176003830155565b600484519101908051906001600160401b038211610961576107c0826107ba8554614766565b856149ac565b602090601f83116001146108fa576107f19291600091836108ef575b50508160011b916000199060031b1c19161790565b90555b60005b8751805182101561082b579061082560019261081e836001600160401b038e511692614889565b5190614f2b565b016107f7565b505097967f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c29392919650946108e46001600160401b0360019751169251935191516108b961088d60405196879687526101006020880152610100870190614221565b9360408601906001600160801b0360408092805115158552826020820151166020860152015116910152565b60a08401906001600160801b0360408092805115158552826020820151166020860152015116910152565b0390a10191926104ad565b015190508d806107dc565b90601f1983169184600052816000209260005b8181106109495750908460019594939210610930575b505050811b0190556107f4565b015160001960f88460031b161c191690558c8080610923565b9293602060018192878601518155019501930161090d565b634e487b7160e01b600052604160045260246000fd5b6001600160401b038851167f1d5ad3c50000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b83356001600160401b038111610294576020916109fc839283369187010161434d565b81520193019261053c565b909291936001600160401b03610a26610a2186888661489d565b614716565b1692610a3184615d91565b15610b7d57836000526007602052610a4f6005604060002001615864565b9260005b8451811015610a8b57600190866000526007602052610a846005604060002001610a7d8389614889565b5190615e25565b5001610a53565b5093909491959250806000526007602052600560406000206000815560006001820155600060028201556000600382015560048101610aca8154614766565b9081610b3a575b5050018054906000815581610b19575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a1019091610498565b6000526020600020908101905b81811015610ae15760008155600101610b26565b81601f60009311600114610b525750555b8880610ad1565b81835260208320610b6d91601f01861c810190600101614982565b8082528160208120915555610b4b565b83631e670e4b60e01b60005260045260246000fd5b34610294576020366003190112610294576004356001600160401b03811161029457610bc2903690600401614520565b610bca614eed565b60005b818110610bd657005b610be18183856148ad565b906001600160401b03610bf383614716565b16600052601060205260ff60406000205416602083013590600282108015610294576000916002811015610ee75783148015610eb8575b610e575750610294576001600160401b03610c4484614716565b16600052601060205260406000209260009360ff1981541660ff8416179055604081013580610cc1575b50610c7890614716565b92610294577f1d1eeb97006356bf772500dc592e232d913119a3143e8452f60e5c98b6a29ca160206001600160401b03600195610cb86040518096614246565b1692a201610bcd565b6000945082610dad576040516340c10f1960e01b815230600482015260248101829052602081604481897f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165af18015610da257918491610c789493610d74575b505b610d3583614716565b96507fbbaa9aea43e3358cd56e894ad9620b8a065abcffab21357fb0702f222480fccc60206001600160401b036000996040519485521692a390610c6e565b610d949060203d8111610d9b575b610d8c81836142da565b81019061496a565b5089610d2a565b503d610d82565b6040513d88823e3d90fd5b8460206001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016610e0584827f0000000000000000000000000000000000000000000000000000000000000000614b36565b602460405180948193630852cd8d60e31b83528760048401525af18015610da257918491610c789493610e39575b50610d2c565b610e509060203d8111610d9b57610d8c81836142da565b5089610e33565b906001600160401b0390610e6a86614716565b90507fe2017d61000000000000000000000000000000000000000000000000000000006000521660045215610ea25760245260446000fd5b634e487b7160e01b600052602160045260246000fd5b50610ee16001600160401b03610ecd87614716565b166000526006602052604060002054151590565b15610c2a565b602483634e487b7160e01b81526021600452fd5b346102945760003660031901126102945760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102945760003660031901126102945760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102945760e036600319011261029457610f9561417d565b606036602319011261029457604051610fad816142bf565b60243580151581036102945781526044356001600160801b03811681036102945760208201526064356001600160801b038116810361029457604082015260603660831901126102945760405190611004826142bf565b608435801515810361029457825260a4356001600160801b038116810361029457602083015260c4356001600160801b03811681036102945760408301526001600160a01b036009541633141580611066575b61042357611064926151b5565b005b506001600160a01b0360015416331415611057565b34610294576020366003190112610294576001600160401b0361109c61417d565b6110a46148ca565b501660005260076020526111166110c66110c160406000206148f5565b61538e565b6040519182918291909160806001600160801b038160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b34610294576000366003190112610294576040516005548082528160208101600560005260206000209260005b81811061121157505061115c925003826142da565b80519061118161116b8361449a565b9261117960405194856142da565b80845261449a565b602083019190601f190136833760005b81518110156111c257806001600160401b036111af60019385614889565b51166111bb8287614889565b5201611191565b5050906040519182916020830190602084525180915260408301919060005b8181106111ef575050500390f35b82516001600160401b03168452859450602093840193909201916001016111e1565b8454835260019485019486945060209093019201611147565b34610294576020366003190112610294576112436141d4565b61124b614eed565b6001600160a01b0381169081156109af576004805473ffffffffffffffffffffffffffffffffffffffff1981169093179055604080516001600160a01b0393841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491819081015b0390a1005b34610294576020366003190112610294576111166112e76112e261417d565b614949565b604051918291602083526020830190614221565b346102945760003660031901126102945760206001600160a01b0360045416604051908152f35b34610294576020366003190112610294576001600160401b0361134361417d565b61134b6148ca565b501660005260076020526111166110c66110c160026040600020016148f5565b34610294576001600160401b03611381366143eb565b92909161138c614eed565b16906113a5826000526006602052604060002054151590565b15611461578160005260076020526113d660056040600020016113c9368685614316565b6020815191012090615e25565b1561141a577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d769192611415604051928392602084526020840191614696565b0390a2005b61145d906040519384937f74f23c7c0000000000000000000000000000000000000000000000000000000085526004850152604060248501526044840191614696565b0390fd5b50631e670e4b60e01b60005260045260246000fd5b346102945760203660031901126102945761148f6146b7565b50600435600052600d6020526040806000206001600160a01b038251916114b58361426e565b546114c360ff821684614844565b81602084019160081c1681526114dc8451809451614550565b51166020820152f35b34610294576000366003190112610294576040516002548082526020820190600260005260206000209060005b81811061153d5761111685611529818703826142da565b60405191829160208352602083019061442a565b8254845260209093019260019283019201611512565b34610294576020366003190112610294576001600160401b0361157461417d565b16600052600760205261158d6005604060002001615864565b8051906115998261449a565b916115a760405193846142da565b8083526115b6601f199161449a565b0160005b81811061167657505060005b815181101561160e57806115dc60019284614889565b5160005260086020526115f260406000206147a0565b6115fc8286614889565b526116078185614889565b50016115c6565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b82821061164757505050500390f35b919360019193955060206116668192603f198a82030186528851614221565b9601920192018594939192611638565b8060606020809387010152016115ba565b34610294576020366003190112610294576116a06141d4565b6116a8614eed565b600f5490816116b357005b60206001600160a01b037f738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599926000600f5561170f85827f000000000000000000000000000000000000000000000000000000000000000061533c565b6040519485521692a2005b34610294576020366003190112610294576004358015610451576001600160a01b03600e54163303610423577f00000000000000000000000000000000000000000000000000000000000000006040516370a0823160e01b81523060048201526020816024816001600160a01b0386165afa90811561185957600091611824575b50600f546117a981856148bd565b82106117e957836117bb81338661533c565b6040519081527fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171960203392a2005b6117f390846148bd565b907fa17e11d50000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b906020823d602011611851575b8161183e602093836142da565b8101031261184e5750518361179b565b80fd5b3d9150611831565b6040513d6000823e3d90fd5b346102945761187336614467565b606060206040516118838161426e565b82815201526080810161189581614702565b6001600160a01b03807f000000000000000000000000000000000000000000000000000000000000000016911603611be257506020810167ffffffffffffffff60801b6118e182614716565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561185957600091611bc3575b50611b995761195261194d60408401614702565b615c88565b6001600160401b0361196382614716565b1661197b816000526006602052604060002054151590565b15611b855760206001600160a01b0360045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa90811561185957600091611b38575b506001600160a01b03163303611b0a576112e281611af793611a0460606119fa611a8796614716565b9201358092614a38565b611a0d816152b9565b7ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae106001600160401b03611a3f84614716565b604080516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152336020820152908101949094521691606090a2614716565b61111660405160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260208152611ac56040826142da565b60405192611ad28461426e565b8352602083019081526040519384936020855251604060208601526060850190614221565b9051838203601f19016040850152614221565b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d602011611b7d575b81611b51602093836142da565b81010312611b795751906001600160a01b038216820361184e57506001600160a01b036119d1565b5080fd5b3d9150611b44565b6354c8163f60e11b60005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b611bdc915060203d602011610d9b57610d8c81836142da565b83611939565b611bf36001600160a01b0391614702565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b34610294576060366003190112610294576004356001600160401b03811161029457611c5190369060040161436b565b906024356001600160401b03811161029457611c71903690600401614520565b906044356001600160401b03811161029457611c91903690600401614520565b6001600160a01b036009541633141580611d45575b61042357838614801590611d3b575b611d115760005b868110611cc557005b80611d0b611cd9610a216001948b8b61489d565b611ce48389896148ad565b611d05611cfd611cf586898b6148ad565b923690614571565b913690614571565b916151b5565b01611cbc565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b5080861415611cb5565b506001600160a01b0360015416331415611ca6565b34610294576040366003190112610294576004356001600160401b03811161029457611d8a903690600401614505565b6024356001600160401b03811161029457611da9903690600401614505565b90611db2614eed565b60005b8151811015611de45780611ddd6001600160a01b03611dd660019486614889565b5116615948565b5001611db5565b5060005b8251811015611e175780611e106001600160a01b03611e0960019487614889565b5116615a37565b5001611de8565b7ffd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193611e54836112be8660405193849360408552604085019061442a565b90838203602085015261442a565b346102945760003660031901126102945760206001600160a01b0360015416604051908152f35b34610294576020366003190112610294576004356001600160401b03811161029457611eb990369060040161436b565b611ec1614eed565b60005b818110611ecd57005b611ed8818385614850565b60a081017f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000611f278361515b565b16146122d3575b60208201611f3b81615199565b90604084019161ffff80611f4e85615199565b1691160161ffff81116122bd5761ffff612710911610156122ac5760808401936001600160401b03611f7f86614716565b16600052600a60205260406000209060e0810194611f9d86836146d0565b60028501916001600160401b03821161096157611fbe826107ba8554614766565b600090601f831160011461224857611fee92916000918361223d5750508160011b916000199060031b1c19161790565b90555b611ffa84615199565b96600184019384549461200c88615199565b60181b64ffff0000001695612020866151a8565b151560c087013597888555606088019c6120398e615188565b60281b68ffffffff0000000000169368ffffffff0000000000199160081b62ffff00169064ffffffffff19161716179060ff1617179055610100840161207f90856146d0565b9091600301916001600160401b038211610961576120a1826107ba8554614766565b600090601f83116001146121c85791806120d7926120de9695946000926121bd5750508160011b916000199060031b1c19161790565b9055614716565b936120e890615199565b946120f290615199565b956120fd90836146d0565b90916121089061515b565b9761211290615188565b9261211c906151a8565b936040519761ffff899816885261ffff16602088015260408701526060860160e0905260e086019061214d92614696565b957fffffffff0000000000000000000000000000000000000000000000000000000016608085015263ffffffff1660a0840152151560c08301526001600160401b031692037f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd91a2600101611ec4565b0135905038806107dc565b838252602082209a9e9d9c9b9a91601f198416815b8181106122255750919e9f9b9c9d9e60019391856120de989796941061220b575b505050811b019055614716565b0135600019600384901b60f8161c191690558f80806121fe565b919360206001819287870135815501950192016121dd565b013590508e806107dc565b8382526020822091601f198416815b818110612294575090846001959493921061227a575b505050811b019055611ff1565b0135600019600384901b60f8161c191690558d808061226d565b83830135855560019094019360209283019201612257565b631c1604c160e11b60005260046000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff6122e460608401615188565b1615611f2e57631c1604c160e11b60005260046000fd5b346102945760203660031901126102945760206123216001600160401b03610ecd61417d565b6040519015158152f35b346102945760c0366003190112610294576004356024356044356001600160401b0381169182820361029457606435926084359160ff831683036102945760a435926001600160a01b038416928385036102945780600052600a60205260ff60016040600020015416612575575b506123bd604051846020820152602081526123b56040826142da565b8288856149f1565b87036125475786600052600d60205260406000206001600160a01b03604051916123e68361426e565b546123f460ff821684614844565b60081c16602082015251956003871015610ea25760009661251b576124239161241c91614d79565b80956150e8565b604051956124308761426e565b600187526020870196338852818752600d602052604087209051976003891015612507578798612504985060ff80198454169116178255517fffffffffffffffffffffff0000000000000000000000000000000000000000ff74ffffffffffffffffffffffffffffffffffffffff0083549260081b1691161790556040519285845260208401527fd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf160403394a4337f0000000000000000000000000000000000000000000000000000000000000000614acc565b80f35b602488634e487b7160e01b81526021600452fd5b602487897fcee81443000000000000000000000000000000000000000000000000000000008252600452fd5b867fcb537aa40000000000000000000000000000000000000000000000000000000060005260045260246000fd5b61258c33600052600c602052604060002054151590565b612399577f6c46a9b5000000000000000000000000000000000000000000000000000000006000526004523360245260446000fd5b34610294576125cf36614467565b6001600160a01b03600454163303612b0d5760a081360312610294576040516125f781614253565b8135815261260760208301614193565b906020810191825260408301356001600160401b0381116102945761262f903690850161434d565b916040820192835260608401356001600160401b03811161029457612657903690860161434d565b93606083019485526080810135906001600160401b038211610294570136601f8201121561029457803561268a8161449a565b9161269860405193846142da565b81835260208084019260061b8201019036821161029457602001915b818310612ad5575050506080830152516001600160401b0381169151925193519182518301946020860193602081880312610294576020810151906001600160401b03821161029457019560a09087900312610294576040519261271784614253565b6020870151845261272a604088016150d9565b916020850192835261273e606089016150d9565b916040860192835260808901519860ff8a168a036102945760608701998a5260a08101516001600160401b03811161029457602091010187601f8201121561029457805161278b816142fb565b986127996040519a8b6142da565b818a5260208284010111610294576127b7916020808b0191016141fe565b6080860196875267ffffffffffffffff60801b60405191632cbc26bb60e01b835260801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561185957600091612ab6575b50611b9957612831818561472a565b15612a9157509560ff61285861289393612887989961ffff8089519351169151169161550c565b61288261286d889a939a518587511690614d79565b9961287b8587511684614d79565b9851614689565b614689565b915116855191886149f1565b9384600052600d602052604060002091604051926128b08461426e565b54936128bf60ff861685614844565b6001600160a01b03602085019560081c16855260009484516003811015612a7d576129995750506000946128f383836150e8565b516020818051810103126129955760200151906001600160a01b038216809203612995579061293192918652601060205260ff604087205416615ec5565b83600052600d6020526040600020600260ff1982541617905551906003821015610ea2576129926060927f33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c9460405193845260208401526040830190614550565ba3005b8580fd5b9490955083929192516003811015612a6957600103612a3d57506129c5856001600160a01b0392614689565b9351169060005260106020526129ef60ff604060002054166129e786866148bd565b903090615ec5565b6129fb84600f546148bd565b600f558280612a0c575b5050612931565b612a36917f000000000000000000000000000000000000000000000000000000000000000061533c565b8582612a05565b80877fb196a44a0000000000000000000000000000000000000000000000000000000060249352600452fd5b602482634e487b7160e01b81526021600452fd5b602487634e487b7160e01b81526021600452fd5b61145d906040519182916324eb47e560e01b8352602060048401526024830190614221565b612acf915060203d602011610d9b57610d8c81836142da565b89612822565b6040833603126102945760206040918251612aef8161426e565b612af8866141ea565b815282860135838201528152019201916126b4565b7fd7f73334000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b34610294576020366003190112610294577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917460206001600160a01b03612b7f6141d4565b612b87614eed565b168073ffffffffffffffffffffffffffffffffffffffff196009541617600955604051908152a1005b34610294576000366003190112610294576000546001600160a01b0381163303612c2c5773ffffffffffffffffffffffffffffffffffffffff19600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346102945760203660031901126102945760206123216001600160a01b03612c7c6141d4565b16600052600c602052604060002054151590565b34610294576020366003190112610294576001600160401b03612cb161417d565b606060c0604051612cc1816142a4565b600081526000602082015260006040820152600083820152600060808201528260a0820152015216600052600a60205260606040600020611116612d03615819565b611e54604051612d12816142a4565b84548152612dbf60018601549563ffffffff602084019760ff81161515895261ffff60408601818360081c168152818c880191818560181c1683528560808a019560281c168552612d786003612d6a60028a016147a0565b9860a08c01998a52016147a0565b9860c08101998a526040519e8f9e8f9260408452516040840152511515910152511660808c0152511660a08a0152511660c08801525160e080880152610120870190614221565b9051858203603f1901610100870152614221565b346102945760003660031901126102945760206001600160a01b0360095416604051908152f35b34610294576020366003190112610294577f64187bd7b97e66658c91904f3021d7c28de967281d18b1a20742348afdd6a6b36001600160a01b03612e3c6141d4565b612e44614eed565b6112be600e549183811673ffffffffffffffffffffffffffffffffffffffff19841617600e5560405193849316839092916001600160a01b0360209181604085019616845216910152565b3461029457600036600319011261029457611116611529615819565b3461029457612eb9366143eb565b612ec4929192614eed565b6001600160401b038216612ee5816000526006602052604060002054151590565b15612f00575061106492612efa913691614316565b90614f2b565b631e670e4b60e01b60005260045260246000fd5b3461029457612f3c612f44612f283661439b565b9491612f35939193614eed565b36916144b1565b9236916144b1565b7f00000000000000000000000000000000000000000000000000000000000000001561304a5760005b8251811015612fd357806001600160a01b03612f8b60019386614889565b5116612f9681615cfd565b612fa2575b5001612f6d565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a184612f9b565b5060005b815181101561106457806001600160a01b03612ff560019385614889565b511680156130445761300681615909565b613013575b505b01612fd7565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a18361300b565b5061300d565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b346102945760403660031901126102945761308d61417d565b6024356001600160401b038111610294576020916130b261232192369060040161434d565b9061472a565b346102945760003660031901126102945760206001600160a01b03600e5416604051908152f35b34610294576020366003190112610294576004356001600160401b0381116102945780600401610100600319833603011261029457600060405161312281614289565b5261314f61314561314061313960c48601856146d0565b3691614316565b614c9e565b6064840135614d79565b906084830161315d81614702565b6001600160a01b03807f000000000000000000000000000000000000000000000000000000000000000016911603611be25750602483019067ffffffffffffffff60801b6131aa83614716565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611859576000916133c3575b50611b99576001600160401b0361321883614716565b16613230816000526006602052604060002054151590565b15611b855760206001600160a01b0360045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa908115611859576000916133a4575b5015611b0a5761329b82614716565b906132b160a48601926130b261313985856146d0565b156133765750507ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc060806001600160401b0361331a6133146044602098613300896132fb8a614716565b6150e8565b0195610a218861330f89614702565b614e68565b94614702565b936001600160a01b0360405195817f000000000000000000000000000000000000000000000000000000000000000016875233898801521660408601528560608601521692a28060405161336d81614289565b52604051908152f35b61338092506146d0565b61145d6040519283926324eb47e560e01b8452602060048501526024840191614696565b6133bd915060203d602011610d9b57610d8c81836142da565b8561328c565b6133dc915060203d602011610d9b57610d8c81836142da565b85613202565b34610294576000366003190112610294576040516370a0823160e01b81523060048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa801561185957600090613451575b602090604051908152f35b506020813d60201161347c575b8161346b602093836142da565b810103126102945760209051613446565b3d915061345e565b34610294576020366003190112610294576001600160401b036134a561417d565b166000526010602052602060ff604060002054166134c66040518092614246565bf35b346102945760a0366003190112610294576134e161417d565b602435906044356001600160401b038111610294576135049036906004016141a7565b9091606435926001600160a01b038416809403610294576084356001600160401b0381116102945761353a9036906004016141a7565b50506135446146b7565b50604051946135528661426e565b60008652600060208701526060608060405161356d81614253565b82815282602082015282604082015260008382015201526001600160401b03831693604051632cbc26bb60e01b815267ffffffffffffffff60801b8560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561185957600091613880575b50611b995761360133615c88565b613618856000526006602052604060002054151590565b1561386b5784600052600a6020526040600020948554831161385257509263ffffffff9261373c869361372e8a9760ff60016137a99c9b01549161ffff8360081c169983602061367b61367561ffff8f9860181c1680988b61550c565b906148bd565b9d019c8d5260281c1680613800575061ffff6136ec61369c60038c016147a0565b985b604051976136ab89614253565b8852602088019c8d52604088019586526060880193857f00000000000000000000000000000000000000000000000000000000000000001685523691614316565b9360808701948552816040519c8d986020808b01525160408a01525116606088015251166080860152511660a08401525160a060c084015260e0830190614221565b03601f1981018652856142da565b60209586946040519061374f87836142da565b6000825261376b60026040519761376589614253565b016147a0565b8652868601526040850152606084015260808301526001600160a01b0360045416906040518097819482936320487ded60e01b8452600484016145bd565b03915afa928315611859576000936137ce575b50826040945283519283525190820152f35b9392508184813d83116137f9575b6137e681836142da565b81010312610294576040935192936137bc565b503d6137dc565b6136ec61ffff91604051906138148261426e565b8152602081016001815260405191630181dcf160e41b602084015251602483015251151560448201526044815261384c6064826142da565b9861369e565b90506358dd87c560e01b60005260045260245260446000fd5b846354c8163f60e11b60005260045260246000fd5b613899915060203d602011610d9b57610d8c81836142da565b886135f3565b34610294576000366003190112610294576020600f54604051908152f35b3461029457600036600319011261029457602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102945760203660031901126102945760206139166141d4565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102945760003660031901126102945760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b34610294576000366003190112610294576111166040516139b46060826142da565b603381527f4879627269645769746845787465726e616c4d696e746572466173745472616e60208201527f73666572546f6b656e506f6f6c20312e362e30000000000000000000000000006040820152604051918291602083526020830190614221565b60c036600319011261029457613a2c61417d565b6064356001600160401b03811161029457613a4b9036906004016141a7565b9091608435916001600160a01b03831683036102945760a4356001600160401b03811161029457613a809036906004016141a7565b505060405190613a8f8261426e565b600082526000602083015260606080604051613aaa81614253565b8281528260208201528260408201526000838201520152604051632cbc26bb60e01b815267ffffffffffffffff60801b8460801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156118595760009161402d575b50611b9957613b3333615c88565b613b536001600160401b0384166000526006602052604060002054151590565b1561400f576001600160401b038316600052600a602052604060002092835460243511613feb5760018401549561ffff8760081c169463ffffffff61ffff8960181c1698613ba76136758b8a60243561550c565b602088015260281c1680613fa05750613bc2600382016147a0565b965b604051613bd081614253565b60243581526020810197885260408101998a5260608101997f000000000000000000000000000000000000000000000000000000000000000060ff169a8b815236613c1c908988614316565b91608084019283526040519a8b9460208601602090525160408601525161ffff1660608501525161ffff1660808401525160ff1660a08301525160c0820160a0905260e08201613c6b91614221565b03601f1981018852613c7d90886142da565b602097604051613c8d8a826142da565b600080825298613ca560026040519661376588614253565b85528a85015260408401526001600160a01b038216606084015260808301526001600160a01b03600454168860405180926320487ded60e01b82528180613cf0888b600484016145bd565b03915afa908115613f95578891613f64575b508652613d1160243585614a38565b60208601516044358111613f345750613d4e60243530337f0000000000000000000000000000000000000000000000000000000000000000614acc565b6001600160401b03841687526010885260ff6040882054166002811015612507579188916001613dcf9414613f24575b6001600160a01b038116613ed1575b506001600160a01b036004541660405180809581947f96f4e9f900000000000000000000000000000000000000000000000000000000835289600484016145bd565b039134905af1958615613ec5578096613e91575b5050957f240a1286fd41f1034c4032dcd6b93fc09e81be4a0b64c7ecee6260b605a8e01691613e868697986001600160401b03613e266020890151602435614689565b936020613e5f613e37368b87614316565b7f0000000000000000000000000000000000000000000000000000000000000000888e6149f1565b99015160405196879687528d87015260408601526080606086015216956080840191614696565b0390a4604051908152f35b909195508682813d8311613ebe575b613eaa81836142da565b8101031261184e5750519381613e86613de3565b503d613ea0565b604051903d90823e3d90fd5b613f1e90613eeb895130336001600160a01b038516614acc565b8851906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000009116614b36565b8a613d8d565b613f2f6024356152b9565b613d7e565b7f61acdb930000000000000000000000000000000000000000000000000000000088526004526044803560245287fd5b90508881813d8311613f8e575b613f7b81836142da565b81010312613f8a57518a613d02565b8780fd5b503d613f71565b6040513d8a823e3d90fd5b60405190613fad8261426e565b8152602081016001815260405191630181dcf160e41b6020840152516024830152511515604482015260448152613fe56064826142da565b96613bc4565b6001600160401b03906358dd87c560e01b6000521660045260243560245260446000fd5b6001600160401b03836354c8163f60e11b6000521660045260246000fd5b614046915060203d602011610d9b57610d8c81836142da565b86613b25565b3461029457602036600319011261029457600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361029457817ff6f46ff9000000000000000000000000000000000000000000000000000000006020931490811561410b575b81156140c7575b5015158152f35b7f85572ffb000000000000000000000000000000000000000000000000000000008114915081156140fa575b50836140c0565b6301ffc9a760e01b915014836140f3565b90507faff2afbf0000000000000000000000000000000000000000000000000000000081148015614154575b8015614144575b906140b9565b506301ffc9a760e01b811461413e565b507f0e64dd29000000000000000000000000000000000000000000000000000000008114614137565b600435906001600160401b038216820361029457565b35906001600160401b038216820361029457565b9181601f84011215610294578235916001600160401b038311610294576020838186019501011161029457565b600435906001600160a01b038216820361029457565b35906001600160a01b038216820361029457565b60005b8381106142115750506000910152565b8181015183820152602001614201565b9060209161423a815180928185528580860191016141fe565b601f01601f1916010190565b906002821015610ea25752565b60a081019081106001600160401b0382111761096157604052565b604081019081106001600160401b0382111761096157604052565b602081019081106001600160401b0382111761096157604052565b60e081019081106001600160401b0382111761096157604052565b606081019081106001600160401b0382111761096157604052565b90601f801991011681019081106001600160401b0382111761096157604052565b6001600160401b03811161096157601f01601f191660200190565b929192614322826142fb565b9161433060405193846142da565b829481845281830111610294578281602093846000960137010152565b9080601f830112156102945781602061436893359101614316565b90565b9181601f84011215610294578235916001600160401b038311610294576020808501948460051b01011161029457565b6040600319820112610294576004356001600160401b03811161029457816143c59160040161436b565b92909291602435906001600160401b038211610294576143e79160040161436b565b9091565b906040600319830112610294576004356001600160401b03811681036102945791602435906001600160401b038211610294576143e7916004016141a7565b906020808351928381520192019060005b8181106144485750505090565b82516001600160a01b031684526020938401939092019160010161443b565b602060031982011261029457600435906001600160401b0382116102945760a09082900360031901126102945760040190565b6001600160401b0381116109615760051b60200190565b9291906144bd8161449a565b936144cb60405195866142da565b602085838152019160051b810192831161029457905b8282106144ed57505050565b602080916144fa846141ea565b8152019101906144e1565b9080601f8301121561029457816020614368933591016144b1565b9181601f84011215610294578235916001600160401b038311610294576020808501946060850201011161029457565b906003821015610ea25752565b35906001600160801b038216820361029457565b919082606091031261029457604051614589816142bf565b809280359081151582036102945760406145b891819385526145ad6020820161455d565b60208601520161455d565b910152565b906001600160401b0390939293168152604060208201526146036145ed845160a0604085015260e0840190614221565b6020850151838203603f19016060850152614221565b90604084015191603f198282030160808301526020808451928381520193019060005b81811061465e575050506080846001600160a01b036060614368969701511660a084015201519060c0603f1982850301910152614221565b825180516001600160a01b031686526020908101518187015260409095019490920191600101614626565b919082039182116122bd57565b908060209392818452848401376000828201840152601f01601f1916010190565b604051906146c48261426e565b60006020838281520152565b903590601e198136030182121561029457018035906001600160401b0382116102945760200191813603831361029457565b356001600160a01b03811681036102945790565b356001600160401b03811681036102945790565b906001600160401b0361436892166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b90600182811c92168015614796575b602083101461478057565b634e487b7160e01b600052602260045260246000fd5b91607f1691614775565b90604051918260008254926147b484614766565b808452936001811690811561482257506001146147db575b506147d9925003836142da565b565b90506000929192526020600020906000915b8183106148065750509060206147d992820101386147cc565b60209193508060019154838589010152019101909184926147ed565b9050602092506147d994915060ff191682840152151560051b820101386147cc565b6003821015610ea25752565b91908110156148735760051b8101359061011e1981360301821215610294570190565b634e487b7160e01b600052603260045260246000fd5b80518210156148735760209160051b010190565b91908110156148735760051b0190565b9190811015614873576060020190565b919082018092116122bd57565b604051906148d782614253565b60006080838281528260208201528260408201528260608201520152565b9060405161490281614253565b60806001829460ff81546001600160801b038116865263ffffffff81861c16602087015260a01c161515604085015201546001600160801b0381166060840152811c910152565b6001600160401b0316600052600760205261436860046040600020016147a0565b90816020910312610294575180151581036102945790565b81811061498d575050565b60008155600101614982565b818102929181159184041417156122bd57565b9190601f81116149bb57505050565b6147d9926000526020600020906020601f840160051c830193106149e7575b601f0160051c0190614982565b90915081906149da565b9290614a24614a329260ff60405195869460208601988952604086015216606084015260808084015260a0830190614221565b03601f1981018352826142da565b51902090565b6001600160401b037fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da817894491169182600052600760205280614aa860406000206001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615533565b604080516001600160a01b039092168252602082019290925290819081015b0390a2565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000060208201526001600160a01b03928316602482015292909116604483015260648201929092526147d991614b3182608481015b03601f1981018452836142da565b615700565b91909181158015614c04575b15614b9a576040517f095ea7b30000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909316602484015260448301919091526147d99190614b318260648101614b23565b608460405162461bcd60e51b815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b0384166024820152602081806044810103816001600160a01b0386165afa90811561185957600091614c6c575b5015614b42565b90506020813d602011614c96575b81614c87602093836142da565b81010312610294575138614c65565b3d9150614c7a565b80518015614d0e57602003614cd057805160208281019183018390031261029457519060ff8211614cd0575060ff1690565b61145d906040519182917f953576f7000000000000000000000000000000000000000000000000000000008352602060048401526024830190614221565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff82116122bd57565b60ff16604d81116122bd57600a0a90565b8115614d63570490565b634e487b7160e01b600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff811692828414614e6157828411614e375790614dbe91614d34565b91604d60ff8416118015614e1c575b614de657505090614de061436892614d48565b90614999565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b50614e2683614d48565b8015614d6357600019048411614dcd565b614e4091614d34565b91604d60ff841611614de657505090614e5b61436892614d48565b90614d59565b5050505090565b6040516340c10f1960e01b81526001600160a01b03909116600482015260248101919091526020818060448101038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1801561185957614ed25750565b614eea9060203d602011610d9b57610d8c81836142da565b50565b6001600160a01b03600154163303614f0157565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b908051156109af576001600160401b0381516020830120921691826000526007602052614f5f8160056040600020016159ba565b15615095576000526008602052604060002081516001600160401b03811161096157614f9581614f8f8454614766565b846149ac565b6020601f821160011461500b5791614fea827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea9593614ac795600091615000575b508160011b916000199060031b1c19161790565b9055604051918291602083526020830190614221565b905084015138614fd6565b601f1982169083600052806000209160005b81811061507d575092614ac79492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea989610615064575b5050811b0190556112e7565b85015160001960f88460031b161c191690553880615058565b9192602060018192868a01518155019401920161501d565b509061145d6040519283927f393b8ad20000000000000000000000000000000000000000000000000000000084526004840152604060248401526044830190614221565b519061ffff8216820361029457565b6001600160401b037f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c91169182600052600760205280614aa860026040600020016001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615533565b357fffffffff00000000000000000000000000000000000000000000000000000000811681036102945790565b3563ffffffff811681036102945790565b3561ffff811681036102945790565b3580151581036102945790565b6001600160401b031660008181526006602052604090205490929190156152a457916152a160e0926152768561520b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b976153fb565b846000526007602052615222816040600020615acb565b61522b836153fb565b846000526007602052615245836002604060002001615acb565b60405194855260208501906001600160801b0360408092805115158552826020820151166020860152015116910152565b60808301906001600160801b0360408092805115158552826020820151166020860152015116910152565ba1565b82631e670e4b60e01b60005260045260246000fd5b602060009160246001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169161531681847f0000000000000000000000000000000000000000000000000000000000000000614b36565b6040519485938492630852cd8d60e31b845260048401525af1801561185957614ed25750565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909216602483015260448201929092526147d991614b318260648101614b23565b6153966148ca565b506001600160801b036060820151166001600160801b0380835116916153db60208501936136756153ce63ffffffff87511642614689565b8560808901511690614999565b808210156153f457505b16825263ffffffff4216905290565b90506153e5565b805115615480576001600160801b036040820151166001600160801b03602083015116106154265750565b60649061547e604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565bfd5b6001600160801b03604082015116158015906154f6575b61549e5750565b60649061547e604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565b506001600160801b036020820151161515615497565b61552f9061ffff6127106155268282969897981684614999565b04951690614999565b0490565b9182549060ff8260a01c161580156156f8575b6156f2576001600160801b038216916001850190815461557963ffffffff6001600160801b0383169360801c1642614689565b908161566c575b505084811061562d57508383106155c25750506155a66001600160801b03928392614689565b16166fffffffffffffffffffffffffffffffff19825416179055565b5460801c916155d18185614689565b926000198101908082116122bd576155f46155f9926001600160a01b03966148bd565b614d59565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b82856001600160a01b03927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b8286929396116156c857615687926136759160801c90614999565b808410156156c35750825b855473ffffffff0000000000000000000000000000000019164260801b63ffffffff60801b16178655923880615580565b615692565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b508215615546565b6001600160a01b0361578291169160409260008085519361572187866142da565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d15615811573d91615766836142fb565b92615773875194856142da565b83523d6000602085013e615f08565b8051908161578f57505050565b6020806157a093830101910161496a565b156157a85750565b6084905162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606091615f08565b60405190600b548083528260208101600b60005260206000209260005b81811061584b5750506147d9925003836142da565b8454835260019485019487945060209093019201615836565b906040519182815491828252602082019060005260206000209260005b8181106158965750506147d9925003836142da565b8454835260019485019487945060209093019201615881565b80548210156148735760005260206000200190600090565b8054906801000000000000000082101561096157816158ee916001615905940181556158af565b819391549060031b91821b91600019901b19161790565b9055565b806000526003602052604060002054156000146159425761592b8160026158c7565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146159425761596a81600b6158c7565b600b5490600052600c602052604060002055600190565b80600052600660205260406000205415600014615942576159a38160056158c7565b600554906000526006602052604060002055600190565b60008281526001820160205260409020546159f157806159dc836001936158c7565b80549260005201602052604060002055600190565b5050600090565b80548015615a21576000190190615a0f82826158af565b8154906000199060031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c602052604090205480156159f15760001981018181116122bd57600b546000198101919082116122bd57808203615a91575b505050615a7d600b6159f8565b600052600c60205260006040812055600190565b615ab3615aa26158ee93600b6158af565b90549060031b1c928392600b6158af565b9055600052600c602052604060002055388080615a70565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1991615bc36060928054615b0863ffffffff8260801c1642614689565b9081615bf9575b50506001600160801b036001816020860151169282815416808510600014615bf157508280855b16166fffffffffffffffffffffffffffffffff19825416178155615b8f8651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff60ff60a01b835492151560a01b169116179055565b60408601516fffffffffffffffffffffffffffffffff1960809190911b16939092166001600160801b031692909217910155565b6152a160405180926001600160801b0360408092805115158552826020820151166020860152015116910152565b838091615b36565b6001600160801b0391615c25839283615c1e6001880154948286169560801c90614999565b91166148bd565b80821015615c8157505b835473ffffffff000000000000000000000000000000001992909116929092161673ffffffffffffffffffffffffffffffffffffffff19909116174260801b63ffffffff60801b161781553880615b0f565b9050615c2f565b7f0000000000000000000000000000000000000000000000000000000000000000615cb05750565b6001600160a01b031680600052600360205260406000205415615cd05750565b7fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b60008181526003602052604090205480156159f15760001981018181116122bd576002546000198101919082116122bd57818103615d57575b505050615d4360026159f8565b600052600360205260006040812055600190565b615d79615d686158ee9360026158af565b90549060031b1c92839260026158af565b90556000526003602052604060002055388080615d36565b60008181526006602052604090205480156159f15760001981018181116122bd576005546000198101919082116122bd57818103615deb575b505050615dd760056159f8565b600052600660205260006040812055600190565b615e0d615dfc6158ee9360056158af565b90549060031b1c92839260056158af565b90556000526006602052604060002055388080615dca565b906001820191816000528260205260406000205490811515600014615ebc576000198201918083116122bd57815460001981019081116122bd578381615e739503615e85575b5050506159f8565b60005260205260006040812055600190565b615ea5615e956158ee93866158af565b90549060031b1c928392866158af565b905560005284602052604060002055388080615e6b565b50505050600090565b9190916002811015610ea257615eff576147d9917f000000000000000000000000000000000000000000000000000000000000000061533c565b6147d991614e68565b91929015615f695750815115615f1c575090565b3b15615f255790565b606460405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015615f7c5750805190602001fd5b61145d9060405191829162461bcd60e51b835260206004840152602483019061422156fea164736f6c634300081a000a",
}

var HybridWithExternalMinterFastTransferTokenPoolABI = HybridWithExternalMinterFastTransferTokenPoolMetaData.ABI

var HybridWithExternalMinterFastTransferTokenPoolBin = HybridWithExternalMinterFastTransferTokenPoolMetaData.Bin

func DeployHybridWithExternalMinterFastTransferTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, minter common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *HybridWithExternalMinterFastTransferTokenPool, error) {
	parsed, err := HybridWithExternalMinterFastTransferTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HybridWithExternalMinterFastTransferTokenPoolBin), backend, minter, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HybridWithExternalMinterFastTransferTokenPool{address: address, abi: *parsed, HybridWithExternalMinterFastTransferTokenPoolCaller: HybridWithExternalMinterFastTransferTokenPoolCaller{contract: contract}, HybridWithExternalMinterFastTransferTokenPoolTransactor: HybridWithExternalMinterFastTransferTokenPoolTransactor{contract: contract}, HybridWithExternalMinterFastTransferTokenPoolFilterer: HybridWithExternalMinterFastTransferTokenPoolFilterer{contract: contract}}, nil
}

type HybridWithExternalMinterFastTransferTokenPool struct {
	address common.Address
	abi     abi.ABI
	HybridWithExternalMinterFastTransferTokenPoolCaller
	HybridWithExternalMinterFastTransferTokenPoolTransactor
	HybridWithExternalMinterFastTransferTokenPoolFilterer
}

type HybridWithExternalMinterFastTransferTokenPoolCaller struct {
	contract *bind.BoundContract
}

type HybridWithExternalMinterFastTransferTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type HybridWithExternalMinterFastTransferTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type HybridWithExternalMinterFastTransferTokenPoolSession struct {
	Contract     *HybridWithExternalMinterFastTransferTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type HybridWithExternalMinterFastTransferTokenPoolCallerSession struct {
	Contract *HybridWithExternalMinterFastTransferTokenPoolCaller
	CallOpts bind.CallOpts
}

type HybridWithExternalMinterFastTransferTokenPoolTransactorSession struct {
	Contract     *HybridWithExternalMinterFastTransferTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type HybridWithExternalMinterFastTransferTokenPoolRaw struct {
	Contract *HybridWithExternalMinterFastTransferTokenPool
}

type HybridWithExternalMinterFastTransferTokenPoolCallerRaw struct {
	Contract *HybridWithExternalMinterFastTransferTokenPoolCaller
}

type HybridWithExternalMinterFastTransferTokenPoolTransactorRaw struct {
	Contract *HybridWithExternalMinterFastTransferTokenPoolTransactor
}

func NewHybridWithExternalMinterFastTransferTokenPool(address common.Address, backend bind.ContractBackend) (*HybridWithExternalMinterFastTransferTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(HybridWithExternalMinterFastTransferTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindHybridWithExternalMinterFastTransferTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPool{address: address, abi: abi, HybridWithExternalMinterFastTransferTokenPoolCaller: HybridWithExternalMinterFastTransferTokenPoolCaller{contract: contract}, HybridWithExternalMinterFastTransferTokenPoolTransactor: HybridWithExternalMinterFastTransferTokenPoolTransactor{contract: contract}, HybridWithExternalMinterFastTransferTokenPoolFilterer: HybridWithExternalMinterFastTransferTokenPoolFilterer{contract: contract}}, nil
}

func NewHybridWithExternalMinterFastTransferTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*HybridWithExternalMinterFastTransferTokenPoolCaller, error) {
	contract, err := bindHybridWithExternalMinterFastTransferTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolCaller{contract: contract}, nil
}

func NewHybridWithExternalMinterFastTransferTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*HybridWithExternalMinterFastTransferTokenPoolTransactor, error) {
	contract, err := bindHybridWithExternalMinterFastTransferTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolTransactor{contract: contract}, nil
}

func NewHybridWithExternalMinterFastTransferTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*HybridWithExternalMinterFastTransferTokenPoolFilterer, error) {
	contract, err := bindHybridWithExternalMinterFastTransferTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolFilterer{contract: contract}, nil
}

func bindHybridWithExternalMinterFastTransferTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HybridWithExternalMinterFastTransferTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.HybridWithExternalMinterFastTransferTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.HybridWithExternalMinterFastTransferTokenPoolTransactor.contract.Transfer(opts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.HybridWithExternalMinterFastTransferTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.contract.Transfer(opts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) ComputeFillId(opts *bind.CallOpts, settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "computeFillId", settlementId, sourceAmountNetFee, sourceDecimals, receiver)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) ComputeFillId(settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ComputeFillId(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, settlementId, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) ComputeFillId(settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ComputeFillId(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, settlementId, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAccumulatedPoolFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetAccumulatedPoolFees() (*big.Int, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAccumulatedPoolFees(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetAccumulatedPoolFees() (*big.Int, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAccumulatedPoolFees(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAllowList(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAllowList(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAllowListEnabled(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAllowListEnabled(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAllowedFillers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetAllowedFillers() ([]common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAllowedFillers(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetAllowedFillers() ([]common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetAllowedFillers(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getCcipSendTokenFee", destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)

	if err != nil {
		return *new(IFastTransferPoolQuote), err
	}

	out0 := *abi.ConvertType(out[0], new(IFastTransferPoolQuote)).(*IFastTransferPoolQuote)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetCcipSendTokenFee(destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetCcipSendTokenFee(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetCcipSendTokenFee(destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetCcipSendTokenFee(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetCurrentInboundRateLimiterState(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetCurrentInboundRateLimiterState(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getDestChainConfig", remoteChainSelector)

	if err != nil {
		return *new(FastTransferTokenPoolAbstractDestChainConfig), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(FastTransferTokenPoolAbstractDestChainConfig)).(*FastTransferTokenPoolAbstractDestChainConfig)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetDestChainConfig(remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetDestChainConfig(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetDestChainConfig(remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetDestChainConfig(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetFillInfo(opts *bind.CallOpts, fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getFillInfo", fillId)

	if err != nil {
		return *new(FastTransferTokenPoolAbstractFillInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FastTransferTokenPoolAbstractFillInfo)).(*FastTransferTokenPoolAbstractFillInfo)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetFillInfo(fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetFillInfo(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, fillId)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetFillInfo(fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetFillInfo(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, fillId)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetGroup(opts *bind.CallOpts, remoteChainSelector uint64) (uint8, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getGroup", remoteChainSelector)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetGroup(remoteChainSelector uint64) (uint8, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetGroup(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetGroup(remoteChainSelector uint64) (uint8, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetGroup(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetLockedTokens(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getLockedTokens")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetLockedTokens() (*big.Int, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetLockedTokens(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetLockedTokens() (*big.Int, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetLockedTokens(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetMinter() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetMinter(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetMinter() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetMinter(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRateLimitAdmin(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRateLimitAdmin(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetRebalancer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRebalancer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetRebalancer() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRebalancer(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetRebalancer() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRebalancer(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRemotePools(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRemotePools(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRemoteToken(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRemoteToken(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRmnProxy(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRmnProxy(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetRouter() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRouter(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetRouter(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetSupportedChains(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetSupportedChains(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetToken() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetToken(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetToken(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetTokenDecimals(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.GetTokenDecimals(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isAllowedFiller", filler)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) IsAllowedFiller(filler common.Address) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsAllowedFiller(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, filler)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) IsAllowedFiller(filler common.Address) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsAllowedFiller(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, filler)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsRemotePool(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsRemotePool(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsSupportedChain(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsSupportedChain(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsSupportedToken(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, token)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.IsSupportedToken(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, token)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) Owner() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.Owner(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) Owner() (common.Address, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.Owner(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SupportsInterface(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, interfaceId)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SupportsInterface(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts, interfaceId)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HybridWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) TypeAndVersion() (string, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.TypeAndVersion(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.TypeAndVersion(&_HybridWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.AcceptOwnership(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.AcceptOwnership(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.AddRemotePool(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.AddRemotePool(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ApplyAllowListUpdates(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, removes, adds)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ApplyAllowListUpdates(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, removes, adds)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ApplyChainUpdates(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ApplyChainUpdates(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipReceive", message)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.CcipReceive(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, message)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.CcipReceive(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, message)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipSendToken", destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) FastFill(opts *bind.TransactOpts, fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "fastFill", fillId, settlementId, sourceChainSelector, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) FastFill(fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.FastFill(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, fillId, settlementId, sourceChainSelector, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) FastFill(fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.FastFill(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, fillId, settlementId, sourceChainSelector, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.LockOrBurn(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, lockOrBurnIn)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.LockOrBurn(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, lockOrBurnIn)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "provideLiquidity", amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ProvideLiquidity(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ProvideLiquidity(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ReleaseOrMint(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, releaseOrMintIn)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.ReleaseOrMint(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, releaseOrMintIn)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) RemoveLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "removeLiquidity", amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) RemoveLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.RemoveLiquidity(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) RemoveLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.RemoveLiquidity(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.RemoveRemotePool(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.RemoveRemotePool(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfig(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfig(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setChainRateLimiterConfigs", remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfigs(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfigs(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetRateLimitAdmin(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, rateLimitAdmin)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetRateLimitAdmin(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, rateLimitAdmin)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setRebalancer", rebalancer)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetRebalancer(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, rebalancer)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetRebalancer(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, rebalancer)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetRouter(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, newRouter)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.SetRouter(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, newRouter)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.TransferOwnership(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, to)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.TransferOwnership(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, to)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) UpdateDestChainConfig(opts *bind.TransactOpts, destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "updateDestChainConfig", destChainConfigArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) UpdateDestChainConfig(destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateDestChainConfig(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, destChainConfigArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) UpdateDestChainConfig(destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateDestChainConfig(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, destChainConfigArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "updateFillerAllowList", fillersToAdd, fillersToRemove)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) UpdateFillerAllowList(fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateFillerAllowList(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, fillersToAdd, fillersToRemove)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) UpdateFillerAllowList(fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateFillerAllowList(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, fillersToAdd, fillersToRemove)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) UpdateGroups(opts *bind.TransactOpts, groupUpdates []HybridWithExternalMinterFastTransferTokenPoolGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "updateGroups", groupUpdates)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) UpdateGroups(groupUpdates []HybridWithExternalMinterFastTransferTokenPoolGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateGroups(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, groupUpdates)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) UpdateGroups(groupUpdates []HybridWithExternalMinterFastTransferTokenPoolGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateGroups(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, groupUpdates)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "withdrawPoolFees", recipient)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) WithdrawPoolFees(recipient common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.WithdrawPoolFees(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, recipient)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) WithdrawPoolFees(recipient common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.WithdrawPoolFees(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, recipient)
}

type HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolAllowListAdd)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolAllowListAdd)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolAllowListAdd)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolAllowListAdd, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolAllowListAdd)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolAllowListRemove)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolAllowListRemove)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolAllowListRemove)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolAllowListRemove, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolAllowListRemove)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolChainAdded)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolChainAdded)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolChainAdded)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseChainAdded(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolChainAdded, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolChainAdded)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolChainConfigured)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolChainConfigured)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolChainConfigured)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseChainConfigured(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolChainConfigured, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolChainConfigured)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolChainRemoved)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolChainRemoved)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolChainRemoved)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseChainRemoved(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolChainRemoved, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolChainRemoved)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolConfigChanged)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolConfigChanged)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolConfigChanged)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseConfigChanged(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolConfigChanged, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolConfigChanged)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated struct {
	DestinationChainSelector uint64
	FastTransferFillerFeeBps uint16
	FastTransferPoolFeeBps   uint16
	MaxFillAmountPerRequest  *big.Int
	DestinationPool          []byte
	ChainFamilySelector      [4]byte
	SettlementOverheadGas    *big.Int
	FillerAllowlistEnabled   bool
	Raw                      types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterDestChainConfigUpdated(opts *bind.FilterOpts, destinationChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator, error) {

	var destinationChainSelectorRule []interface{}
	for _, destinationChainSelectorItem := range destinationChainSelector {
		destinationChainSelectorRule = append(destinationChainSelectorRule, destinationChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "DestChainConfigUpdated", destinationChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "DestChainConfigUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, destinationChainSelector []uint64) (event.Subscription, error) {

	var destinationChainSelectorRule []interface{}
	for _, destinationChainSelectorItem := range destinationChainSelector {
		destinationChainSelectorRule = append(destinationChainSelectorRule, destinationChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "DestChainConfigUpdated", destinationChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseDestChainConfigUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated struct {
	DestChainSelector uint64
	DestinationPool   common.Address
	Raw               types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterDestinationPoolUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "DestinationPoolUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "DestinationPoolUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchDestinationPoolUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "DestinationPoolUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestinationPoolUpdated", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseDestinationPoolUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestinationPoolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled struct {
	FillId       [32]byte
	SettlementId [32]byte
	Filler       common.Address
	DestAmount   *big.Int
	Receiver     common.Address
	Raw          types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterFastTransferFilled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator, error) {

	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}
	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FastTransferFilled", fillIdRule, settlementIdRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "FastTransferFilled", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchFastTransferFilled(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (event.Subscription, error) {

	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}
	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FastTransferFilled", fillIdRule, settlementIdRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferFilled", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseFastTransferFilled(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferFilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested struct {
	DestinationChainSelector uint64
	FillId                   [32]byte
	SettlementId             [32]byte
	SourceAmountNetFee       *big.Int
	SourceDecimals           uint8
	FastTransferFee          *big.Int
	Receiver                 []byte
	Raw                      types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator, error) {

	var destinationChainSelectorRule []interface{}
	for _, destinationChainSelectorItem := range destinationChainSelector {
		destinationChainSelectorRule = append(destinationChainSelectorRule, destinationChainSelectorItem)
	}
	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FastTransferRequested", destinationChainSelectorRule, fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "FastTransferRequested", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchFastTransferRequested(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error) {

	var destinationChainSelectorRule []interface{}
	for _, destinationChainSelectorItem := range destinationChainSelector {
		destinationChainSelectorRule = append(destinationChainSelectorRule, destinationChainSelectorItem)
	}
	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FastTransferRequested", destinationChainSelectorRule, fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferRequested", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseFastTransferRequested(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled struct {
	FillId                    [32]byte
	SettlementId              [32]byte
	FillerReimbursementAmount *big.Int
	PoolFeeAccumulated        *big.Int
	PrevState                 uint8
	Raw                       types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterFastTransferSettled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator, error) {

	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FastTransferSettled", fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "FastTransferSettled", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchFastTransferSettled(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error) {

	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FastTransferSettled", fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferSettled", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseFastTransferSettled(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated struct {
	AddFillers    []common.Address
	RemoveFillers []common.Address
	Raw           types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterFillerAllowListUpdated(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FillerAllowListUpdated")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "FillerAllowListUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchFillerAllowListUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FillerAllowListUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FillerAllowListUpdated", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseFillerAllowListUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FillerAllowListUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolGroupUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolGroupUpdated)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolGroupUpdated)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolGroupUpdated struct {
	RemoteChainSelector uint64
	Group               uint8
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterGroupUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "GroupUpdated", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "GroupUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchGroupUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolGroupUpdated, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "GroupUpdated", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolGroupUpdated)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "GroupUpdated", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseGroupUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolGroupUpdated, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolGroupUpdated)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "GroupUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "InboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseInboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded struct {
	Rebalancer common.Address
	Amount     *big.Int
	Raw        types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterLiquidityAdded(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "LiquidityAdded", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "LiquidityAdded", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded, rebalancer []common.Address) (event.Subscription, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "LiquidityAdded", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseLiquidityAdded(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated struct {
	RemoteChainSelector uint64
	Group               uint8
	RemoteChainSupply   *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterLiquidityMigrated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "LiquidityMigrated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "LiquidityMigrated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchLiquidityMigrated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "LiquidityMigrated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LiquidityMigrated", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseLiquidityMigrated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LiquidityMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved struct {
	Rebalancer common.Address
	Amount     *big.Int
	Raw        types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterLiquidityRemoved(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "LiquidityRemoved", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "LiquidityRemoved", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved, rebalancer []common.Address) (event.Subscription, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "LiquidityRemoved", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseLiquidityRemoved(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "LockedOrBurned", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseLockedOrBurned(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "OutboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseOutboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterPoolFeeWithdrawn(opts *bind.FilterOpts, recipient []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "PoolFeeWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "PoolFeeWithdrawn", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchPoolFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "PoolFeeWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "PoolFeeWithdrawn", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParsePoolFeeWithdrawn(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "PoolFeeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolRebalancerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRebalancerSet)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRebalancerSet)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolRebalancerSet struct {
	OldRebalancer common.Address
	NewRebalancer common.Address
	Raw           types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterRebalancerSet(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RebalancerSet")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "RebalancerSet", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchRebalancerSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRebalancerSet) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RebalancerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolRebalancerSet)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RebalancerSet", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseRebalancerSet(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRebalancerSet, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolRebalancerSet)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RebalancerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Recipient           common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "ReleasedOrMinted", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseReleasedOrMinted(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator struct {
	Event *HybridWithExternalMinterFastTransferTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRouterUpdated)
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
		it.Event = new(HybridWithExternalMinterFastTransferTokenPoolRouterUpdated)
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

func (it *HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterFastTransferTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterFastTransferTokenPoolRouterUpdated)
				if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRouterUpdated, error) {
	event := new(HybridWithExternalMinterFastTransferTokenPoolRouterUpdated)
	if err := _HybridWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["AllowListAdd"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseAllowListAdd(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["AllowListRemove"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseAllowListRemove(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["ChainAdded"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseChainAdded(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["ChainConfigured"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseChainConfigured(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["ChainRemoved"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseChainRemoved(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["ConfigChanged"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseConfigChanged(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["DestChainConfigUpdated"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseDestChainConfigUpdated(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["DestinationPoolUpdated"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseDestinationPoolUpdated(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["FastTransferFilled"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseFastTransferFilled(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["FastTransferRequested"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseFastTransferRequested(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["FastTransferSettled"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseFastTransferSettled(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["FillerAllowListUpdated"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseFillerAllowListUpdated(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["GroupUpdated"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseGroupUpdated(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["InboundRateLimitConsumed"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseInboundRateLimitConsumed(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["LiquidityAdded"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseLiquidityAdded(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["LiquidityMigrated"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseLiquidityMigrated(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["LiquidityRemoved"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseLiquidityRemoved(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["LockedOrBurned"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseLockedOrBurned(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["OutboundRateLimitConsumed"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseOutboundRateLimitConsumed(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseOwnershipTransferRequested(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseOwnershipTransferred(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["PoolFeeWithdrawn"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParsePoolFeeWithdrawn(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseRateLimitAdminSet(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["RebalancerSet"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseRebalancerSet(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["ReleasedOrMinted"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseReleasedOrMinted(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseRemotePoolAdded(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseRemotePoolRemoved(log)
	case _HybridWithExternalMinterFastTransferTokenPool.abi.Events["RouterUpdated"].ID:
		return _HybridWithExternalMinterFastTransferTokenPool.ParseRouterUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (HybridWithExternalMinterFastTransferTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (HybridWithExternalMinterFastTransferTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (HybridWithExternalMinterFastTransferTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (HybridWithExternalMinterFastTransferTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (HybridWithExternalMinterFastTransferTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (HybridWithExternalMinterFastTransferTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd")
}

func (HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated) Topic() common.Hash {
	return common.HexToHash("0xb760e03fa04c0e86fcff6d0046cdcf22fb5d5b6a17d1e6f890b3456e81c40fd8")
}

func (HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled) Topic() common.Hash {
	return common.HexToHash("0xd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf1")
}

func (HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x240a1286fd41f1034c4032dcd6b93fc09e81be4a0b64c7ecee6260b605a8e016")
}

func (HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled) Topic() common.Hash {
	return common.HexToHash("0x33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c")
}

func (HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated) Topic() common.Hash {
	return common.HexToHash("0xfd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193")
}

func (HybridWithExternalMinterFastTransferTokenPoolGroupUpdated) Topic() common.Hash {
	return common.HexToHash("0x1d1eeb97006356bf772500dc592e232d913119a3143e8452f60e5c98b6a29ca1")
}

func (HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0x50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c")
}

func (HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded) Topic() common.Hash {
	return common.HexToHash("0xc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088")
}

func (HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated) Topic() common.Hash {
	return common.HexToHash("0xbbaa9aea43e3358cd56e894ad9620b8a065abcffab21357fb0702f222480fccc")
}

func (HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved) Topic() common.Hash {
	return common.HexToHash("0xc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719")
}

func (HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned) Topic() common.Hash {
	return common.HexToHash("0xf33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae10")
}

func (HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0xff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da8178944")
}

func (HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599")
}

func (HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (HybridWithExternalMinterFastTransferTokenPoolRebalancerSet) Topic() common.Hash {
	return common.HexToHash("0x64187bd7b97e66658c91904f3021d7c28de967281d18b1a20742348afdd6a6b3")
}

func (HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted) Topic() common.Hash {
	return common.HexToHash("0xfc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0")
}

func (HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (HybridWithExternalMinterFastTransferTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPool) Address() common.Address {
	return _HybridWithExternalMinterFastTransferTokenPool.address
}

type HybridWithExternalMinterFastTransferTokenPoolInterface interface {
	ComputeFillId(opts *bind.CallOpts, settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error)

	GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error)

	GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error)

	GetFillInfo(opts *bind.CallOpts, fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error)

	GetGroup(opts *bind.CallOpts, remoteChainSelector uint64) (uint8, error)

	GetLockedTokens(opts *bind.CallOpts) (*big.Int, error)

	GetMinter(opts *bind.CallOpts) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRebalancer(opts *bind.CallOpts) (common.Address, error)

	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	GetTokenDecimals(opts *bind.CallOpts) (uint8, error)

	IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error)

	IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)

	IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error)

	IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error)

	FastFill(opts *bind.TransactOpts, fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDestChainConfig(opts *bind.TransactOpts, destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error)

	UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error)

	UpdateGroups(opts *bind.TransactOpts, groupUpdates []HybridWithExternalMinterFastTransferTokenPoolGroupUpdate) (*types.Transaction, error)

	WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolAllowListRemove, error)

	FilterChainAdded(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolConfigChanged, error)

	FilterDestChainConfigUpdated(opts *bind.FilterOpts, destinationChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator, error)

	WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, destinationChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, error)

	FilterDestinationPoolUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator, error)

	WatchDestinationPoolUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParseDestinationPoolUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, error)

	FilterFastTransferFilled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator, error)

	WatchFastTransferFilled(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (event.Subscription, error)

	ParseFastTransferFilled(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferFilled, error)

	FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator, error)

	WatchFastTransferRequested(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error)

	ParseFastTransferRequested(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferRequested, error)

	FilterFastTransferSettled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator, error)

	WatchFastTransferSettled(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error)

	ParseFastTransferSettled(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFastTransferSettled, error)

	FilterFillerAllowListUpdated(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator, error)

	WatchFillerAllowListUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated) (event.Subscription, error)

	ParseFillerAllowListUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated, error)

	FilterGroupUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator, error)

	WatchGroupUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolGroupUpdated, remoteChainSelector []uint64) (event.Subscription, error)

	ParseGroupUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolGroupUpdated, error)

	FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator, error)

	WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseInboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, error)

	FilterLiquidityAdded(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityAddedIterator, error)

	WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded, rebalancer []common.Address) (event.Subscription, error)

	ParseLiquidityAdded(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityAdded, error)

	FilterLiquidityMigrated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityMigratedIterator, error)

	WatchLiquidityMigrated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error)

	ParseLiquidityMigrated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityMigrated, error)

	FilterLiquidityRemoved(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityRemovedIterator, error)

	WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved, rebalancer []common.Address) (event.Subscription, error)

	ParseLiquidityRemoved(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLiquidityRemoved, error)

	FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator, error)

	WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error)

	ParseLockedOrBurned(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolLockedOrBurned, error)

	FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator, error)

	WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseOutboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolOwnershipTransferred, error)

	FilterPoolFeeWithdrawn(opts *bind.FilterOpts, recipient []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator, error)

	WatchPoolFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, recipient []common.Address) (event.Subscription, error)

	ParsePoolFeeWithdrawn(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRateLimitAdminSet, error)

	FilterRebalancerSet(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator, error)

	WatchRebalancerSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRebalancerSet) (event.Subscription, error)

	ParseRebalancerSet(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRebalancerSet, error)

	FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator, error)

	WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error)

	ParseReleasedOrMinted(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolReleasedOrMinted, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*HybridWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*HybridWithExternalMinterFastTransferTokenPoolRouterUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
