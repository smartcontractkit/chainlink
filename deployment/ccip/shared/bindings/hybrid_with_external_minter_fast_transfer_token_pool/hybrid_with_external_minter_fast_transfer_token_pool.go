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

type HybridTokenPoolAbstractGroupUpdate struct {
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSendToken\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"computeFillId\",\"inputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"fastFill\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAccumulatedPoolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedFillers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCcipSendTokenFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structIFastTransferPool.Quote\",\"components\":[{\"name\":\"ccipSettlementFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfig\",\"components\":[{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFillInfo\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.FillInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIFastTransferPool.FillState\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGroup\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumHybridTokenPoolAbstract.Group\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockedTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRebalancer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAllowedFiller\",\"inputs\":[{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"provideLiquidity\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRebalancer\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateDestChainConfig\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfigUpdateArgs[]\",\"components\":[{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFillerAllowList\",\"inputs\":[{\"name\":\"fillersToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"fillersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateGroups\",\"inputs\":[{\"name\":\"groupUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structHybridTokenPoolAbstract.GroupUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"internalType\":\"enumHybridTokenPoolAbstract.Group\"},{\"name\":\"remoteChainSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawLiquidity\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPoolFees\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestinationPoolUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destinationPool\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"filler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferRequested\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"fillerFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fillerReimbursementAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFeeAccumulated\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"prevState\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIFastTransferPool.FillState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FillerAllowListUpdated\",\"inputs\":[{\"name\":\"addFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"removeFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GroupUpdated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumHybridTokenPoolAbstract.Group\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityAdded\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityMigrated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumHybridTokenPoolAbstract.Group\"},{\"name\":\"remoteChainSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityRemoved\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolFeeWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RebalancerSet\",\"inputs\":[{\"name\":\"oldRebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyFilledOrSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadySettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"FillerNotAllowlisted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientLiquidity\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientLiquidityForGroupUpdate\",\"inputs\":[{\"name\":\"balanceBeforeMigration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balanceAfterMigration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"accumulatedPoolFees\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientPoolFees\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFillId\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGroupUpdate\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"internalType\":\"enumHybridTokenPoolAbstract.Group\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LiquidityAmountCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"QuoteFeeExceedsUserMaxLimit\",\"inputs\":[{\"name\":\"quoteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferAmountExceedsMaxFillAmount\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x6101408060405234610344576164ac803803809161001d8285610575565b8339810160c0828203126103445761003482610598565b60208301516001600160a01b03811693919290918483036103445761005b604082016105ac565b60608201519091906001600160401b0381116103445781019280601f85011215610344578351936001600160401b03851161055f578460051b9060208201956100a76040519788610575565b865260208087019282010192831161034457602001905b828210610547575050506100e060a06100d960808401610598565b9201610598565b93331561053657600180546001600160a01b0319163317905586158015610525575b8015610514575b6105035760805260c05260405163313ce56760e01b8152602081600481895afa600091816104c7575b5061049c575b5060a052600480546001600160a01b0319166001600160a01b0384169081179091558151151560e0819052909190610373575b501561035d57610100526101208190526040516321df0da760e01b815290602090829060049082906001600160a01b03165afa90811561035157600091610312575b506001600160a01b0316908181036102fb57604051615d51908161075b82396080518181816103a101528181610b0d0152818161170d015281816117740152818161192b015281816123500152818161283801528181612ef5015281816130bb0152818161319f0152818161374a015281816137970152818161387801528181613cad0152818161498101528181615350015281816155510152615c83015260a05181818161196f0152818161348e0152818161370001528181613b7001528181613db101528181613de301528181614df90152615519015260c051818181610f7e015281816117dc0152818161260f01528181612f5e015281816133910152613a61015260e051818181610f3901528181612cdc0152615bde01526101005181613e9901526101205181818161027001528181610d1601528181610de50152818161532901526155ad0152f35b63f902523f60e01b60005260045260245260446000fd5b90506020813d602011610349575b8161032d60209383610575565b810103126103445761033e90610598565b386101ad565b600080fd5b3d9150610320565b6040513d6000823e3d90fd5b6335fdcccd60e21b600052600060045260246000fd5b919290602092604051926103878585610575565b60008452600036813760e0511561048b5760005b8451811015610402576001906001600160a01b036103b982886105ba565b5116876103c5826105fc565b6103d2575b50500161039b565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a138876103ca565b50919490925060005b835181101561047f576001906001600160a01b0361042982876105ba565b51168015610479578661043b826106fa565b610449575b50505b0161040b565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a13886610440565b50610443565b5092509290503861016b565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036104b05750610138565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d6020116104fb575b816104e360209383610575565b81010312610344576104f4906105ac565b9038610132565b3d91506104d6565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03821615610109565b506001600160a01b03851615610102565b639b15e16f60e01b60005260046000fd5b6020809161055484610598565b8152019101906100be565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b0382119082101761055f57604052565b51906001600160a01b038216820361034457565b519060ff8216820361034457565b80518210156105ce5760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156105ce5760005260206000200190600090565b60008181526003602052604090205480156106f35760001981018181116106dd576002546000198101919082116106dd5781810361068c575b505050600254801561067657600019016106508160026105e4565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b6106c561069d6106ae9360026105e4565b90549060031b1c92839260026105e4565b819391549060031b91821b91600019901b19161790565b90556000526003602052604060002055388080610635565b634e487b7160e01b600052601160045260246000fd5b5050600090565b80600052600360205260406000205415600014610754576002546801000000000000000081101561055f5761073b6106ae82600185940160025560026105e4565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a71461400157508063055befd41461395b5780630a861f2a14613841578063181f5a77146137bb57806321df0da714613777578063240028e81461372457806324f65ee7146136e65780632b2c0eb4146136c85780632e7aa8c814613261578063319ac1011461321b5780633317bbcc146131735780633907753714612e5c578063432a6ba314612e355780634c5ef0ed14612df157806354c8a4f314612caa57806362ddd3c414612c415780636609f59914612c255780636cfd155314612bc45780636d3d1a5814612b9d5780636def4ce714612a5a57806378b410f214612a2057806379ba5097146129985780637d54534e1461293057806385572ffb146123e757806387f060d0146121ae5780638926f54f1461217e5780638a18dcbd14611d035780638da5cb5b14611cdc578063929ea5ba14611bd4578063962d402014611ab45780639a4575b91461173c5780639fe280f5146116a9578063a42a7b8b14611575578063a7cd63b714611507578063abe1c1e814611498578063acfecf91146113a6578063af58d59f1461135d578063b0f479a114611336578063b7946580146112fe578063c0d7865514611265578063c4bffe2b14611155578063c75eea9c146110b6578063cf7401f314610fa2578063dc0bd97114610f5e578063e0351e1314610f21578063e7e62f8514610ada578063e8a1da1714610404578063eb521a4c14610377578063eeebc67414610320578063f2fde38b146102995763f36675171461025057600080fd5b346102945760003660031901126102945760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b34610294576020366003190112610294576001600160a01b036102ba61410c565b6102c2614be8565b1633811461030f57806001600160a01b031960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b636d6c4ee560e11b60005260046000fd5b346102945760803660031901126102945760443560ff81168103610294576064356001600160401b0381116102945760209161036361036f923690600401614293565b906024356004356148f2565b604051908152f35b346102945760203660031901126102945760043580156103f357610399614b55565b6103c58130337f00000000000000000000000000000000000000000000000000000000000000006149cd565b6040519081527fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208860203392a2005b63a90c0d1960e01b60005260046000fd5b3461029457610412366142e1565b91909261041d614be8565b6000905b82821061094f5750505060009063ffffffff42165b81831061043f57005b61044a838386614728565b9261012084360312610294576040519361046385614199565b61046c816140cb565b855260208101356001600160401b0381116102945781019336601f8601121561029457843561049a816143e0565b956104a86040519788614220565b81875260208088019260051b820101903682116102945760208101925b82841061092157505050506020860194855260408201356001600160401b038111610294576104f79036908401614293565b906040870191825261052161050f36606086016144b7565b936060890194855260c03691016144b7565b9460808801958652610533845161509d565b61053d865161509d565b82515115610910576105586001600160401b03895116615744565b156108f1576001600160401b038851166000526007602052604060002061063e85516001600160801b03604082015116906106116001600160801b03602083015116915115158360806040516105ad81614199565b858152602081018a905260408101849052606081018690520152855460ff60a01b91151560a01b9190911674ffffffffffffffffffffffffffffffffffffffffff199091166001600160801b0384161763ffffffff60801b608089901b1617178555565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176001830155565b61070e87516001600160801b03604082015116906106e16001600160801b036020830151169151151583608060405161067681614199565b858152602081018a90526040810184905260608101869052015260028601805460ff60a01b92151560a01b9290921674ffffffffffffffffffffffffffffffffffffffffff199092166001600160801b0385161763ffffffff60801b60808a901b1617919091179055565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176003830155565b600484519101908051906001600160401b0382116108db5761073a82610734855461463e565b856148ad565b602090601f83116001146108745761076b929160009183610869575b50508160011b916000199060031b1c19161790565b90555b60005b875180518210156107a5579061079f600192610798836001600160401b038e511692614761565b5190614c0d565b01610771565b505097967f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c293929196509461085e6001600160401b03600197511692519351915161083361080760405196879687526101006020880152610100870190614159565b9360408601906001600160801b0360408092805115158552826020820151166020860152015116910152565b60a08401906001600160801b0360408092805115158552826020820151166020860152015116910152565b0390a1019192610436565b015190508d80610756565b90601f1983169184600052816000209260005b8181106108c357509084600195949392106108aa575b505050811b01905561076e565b015160001960f88460031b161c191690558c808061089d565b92936020600181928786015181550195019301610887565b634e487b7160e01b600052604160045260246000fd5b6001600160401b03885116631d5ad3c560e01b60005260045260246000fd5b6342bcdf7f60e11b60005260046000fd5b83356001600160401b038111610294576020916109448392833691870101614293565b8152019301926104c5565b909291936001600160401b0361096e610969868886614775565b614785565b169261097984615aa8565b15610ac5578360005260076020526109976005604060002001615627565b9260005b84518110156109d3576001908660005260076020526109cc60056040600020016109c58389614761565b5190615b3c565b500161099b565b5093909491959250806000526007602052600560406000206000815560006001820155600060028201556000600382015560048101610a12815461463e565b9081610a82575b5050018054906000815581610a61575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a1019091610421565b6000526020600020908101905b81811015610a295760008155600101610a6e565b81601f60009311600114610a9a5750555b8880610a19565b81835260208320610ab591601f01861c810190600101614883565b8082528160208120915555610a93565b83631e670e4b60e01b60005260045260246000fd5b34610294576020366003190112610294576004356001600160401b03811161029457610b0a903690600401614466565b907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03811691604051936370a0823160e01b8552306004860152602085602481875afa948515610c0057600095610eed575b50610b6c614be8565b60005b818110610c0c576040516370a0823160e01b815230600482015286906020816024818a5afa908115610c0057600091610bcb575b5060105491828210610bb157005b630e70aaa360e41b60005260045260245260445260646000fd5b906020823d602011610bf8575b81610be560209383614220565b81010312610bf557505182610ba3565b80fd5b3d9150610bd8565b6040513d6000823e3d90fd5b610c17818385614799565b906001600160401b03610c2983614785565b16600052600f60205260ff60406000205416602083013590600282108015610294576000916002811015610ed95783148015610eaa575b610e625750610294576001600160401b03610c7a84614785565b16600052600f60205260406000209260009360ff1981541660ff8416179055604081013580610cea575b50610cae90614785565b92610294576001926001600160401b03167f1d1eeb97006356bf772500dc592e232d913119a3143e8452f60e5c98b6a29ca1600080a301610b6f565b6000945082610dd6576040516340c10f1960e01b815230600482015260248101829052602081604481897f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165af18015610dcb57918491610cae9493610d9d575b505b610d5e83614785565b96507fbbaa9aea43e3358cd56e894ad9620b8a065abcffab21357fb0702f222480fccc60206001600160401b036000996040519485521692a390610ca4565b610dbd9060203d8111610dc4575b610db58183614220565b810190615085565b508c610d53565b503d610dab565b6040513d88823e3d90fd5b8460206001600160a01b0360247f0000000000000000000000000000000000000000000000000000000000000000610e0f86828f614a1e565b604051630852cd8d60e31b8152600481018790529485938492165af18015610dcb57918491610cae9493610e44575b50610d55565b610e5b9060203d8111610dc457610db58183614220565b508c610e3e565b906001600160401b0390610e7586614785565b905063e2017d6160e01b6000521660045215610e945760245260446000fd5b634e487b7160e01b600052602160045260246000fd5b50610ed36001600160401b03610ebf87614785565b166000526006602052604060002054151590565b15610c60565b634e487b7160e01b83526021600452602483fd5b9094506020813d602011610f19575b81610f0960209383614220565b8101031261029457519385610b63565b3d9150610efc565b346102945760003660031901126102945760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102945760003660031901126102945760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102945760e036600319011261029457610fbb6140b5565b606036602319011261029457604051610fd381614205565b60243580151581036102945781526044356001600160801b03811681036102945760208201526064356001600160801b03811681036102945760408201526060366083190112610294576040519061102a82614205565b608435801515810361029457825260a4356001600160801b038116810361029457602083015260c4356001600160801b03811681036102945760408301526001600160a01b0360095416331415806110a1575b61108c5761108a92614f0e565b005b63472511eb60e11b6000523360045260246000fd5b506001600160a01b036001541633141561107d565b34610294576020366003190112610294576001600160401b036110d76140b5565b6110df6147e3565b501660005260076020526111516111016110fc604060002061480e565b615012565b6040519182918291909160806001600160801b038160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b34610294576000366003190112610294576040516005548082528160208101600560005260206000209260005b81811061124c57505061119792500382614220565b8051906111bc6111a6836143e0565b926111b46040519485614220565b8084526143e0565b602083019190601f190136833760005b81518110156111fd57806001600160401b036111ea60019385614761565b51166111f68287614761565b52016111cc565b5050906040519182916020830190602084525180915260408301919060005b81811061122a575050500390f35b82516001600160401b031684528594506020938401939092019160010161121c565b8454835260019485019486945060209093019201611182565b346102945760203660031901126102945761127e61410c565b611286614be8565b6001600160a01b038116908115610910576004805473ffffffffffffffffffffffffffffffffffffffff1981169093179055604080516001600160a01b0393841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491819081015b0390a1005b346102945760203660031901126102945761115161132261131d6140b5565b614862565b604051918291602083526020830190614159565b346102945760003660031901126102945760206001600160a01b0360045416604051908152f35b34610294576020366003190112610294576001600160401b0361137e6140b5565b6113866147e3565b501660005260076020526111516111016110fc600260406000200161480e565b34610294576001600160401b036113bc36614331565b9290916113c7614be8565b16906113e0826000526006602052604060002054151590565b1561148357816000526007602052611411600560406000200161140436868561425c565b6020815191012090615b3c565b15611455577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7691926114506040519283926020845260208401916147c2565b0390a2005b61147f90604051938493631d3c8f1f60e21b855260048501526040602485015260448401916147c2565b0390fd5b50631e670e4b60e01b60005260045260246000fd5b34610294576020366003190112610294576114b16145e9565b50600435600052600d6020526040806000206001600160a01b038251916114d7836141b4565b546114e560ff82168461471c565b81602084019160081c1681526114fe8451809451614496565b51166020820152f35b34610294576000366003190112610294576040516002548082526020820190600260005260206000209060005b81811061155f576111518561154b81870382614220565b604051918291602083526020830190614370565b8254845260209093019260019283019201611534565b34610294576020366003190112610294576001600160401b036115966140b5565b1660005260076020526115af6005604060002001615627565b8051906115bb826143e0565b916115c96040519384614220565b8083526115d8601f19916143e0565b0160005b81811061169857505060005b815181101561163057806115fe60019284614761565b5160005260086020526116146040600020614678565b61161e8286614761565b526116298185614761565b50016115e8565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b82821061166957505050500390f35b919360019193955060206116888192603f198a82030186528851614159565b960192019201859493919261165a565b8060606020809387010152016115dc565b34610294576020366003190112610294576116c261410c565b6116ca614be8565b60105490816116d557005b60206001600160a01b037f738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b85959992600060105561173185827f0000000000000000000000000000000000000000000000000000000000000000614b69565b6040519485521692a2005b346102945761174a366143ad565b6117526147a9565b5061175b6147a9565b506080810161176981614bd4565b6001600160a01b03807f000000000000000000000000000000000000000000000000000000000000000016911603611a8e57506020810167ffffffffffffffff60801b6117b582614785565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610c0057600091611a6f575b50611a5e5761182661182160408401614bd4565b615bdc565b6001600160401b0361183782614785565b1661184f816000526006602052604060002054151590565b15611a4a5760206001600160a01b03600454169160246040518094819363a8d87a3b60e01b835260048301525afa908115610c00576000916119fd575b506001600160a01b031633036119e85761131d816119d5936118bf60606118b561196596614785565b9201358092614939565b6001600160401b036118d083614785565b16600052600f6020526118eb8160ff60406000205416615307565b7ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae106001600160401b0361191d84614785565b604080516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152336020820152908101949094521691606090a2614785565b61115160405160ff7f0000000000000000000000000000000000000000000000000000000000000000166020820152602081526119a3604082614220565b604051926119b0846141b4565b8352602083019081526040519384936020855251604060208601526060850190614159565b9051838203601f19016040850152614159565b63728fe07b60e01b6000523360045260246000fd5b6020813d602011611a42575b81611a1660209383614220565b81010312611a3e5751906001600160a01b0382168203610bf557506001600160a01b0361188c565b5080fd5b3d9150611a09565b6354c8163f60e11b60005260045260246000fd5b630a75a23b60e31b60005260046000fd5b611a88915060203d602011610dc457610db58183614220565b8361180d565b611a9f6001600160a01b0391614bd4565b63961c9a4f60e01b6000521660045260246000fd5b34610294576060366003190112610294576004356001600160401b03811161029457611ae49036906004016142b1565b906024356001600160401b03811161029457611b04903690600401614466565b906044356001600160401b03811161029457611b24903690600401614466565b6001600160a01b036009541633141580611bbf575b61108c57838614801590611bb5575b611ba45760005b868110611b5857005b80611b9e611b6c6109696001948b8b614775565b611b77838989614799565b611b98611b90611b8886898b614799565b9236906144b7565b9136906144b7565b91614f0e565b01611b4f565b632b477e7160e11b60005260046000fd5b5080861415611b48565b506001600160a01b0360015416331415611b39565b34610294576040366003190112610294576004356001600160401b03811161029457611c0490369060040161444b565b6024356001600160401b03811161029457611c2390369060040161444b565b90611c2c614be8565b60005b8151811015611c5e5780611c576001600160a01b03611c5060019486614761565b511661570b565b5001611c2f565b5060005b8251811015611c915780611c8a6001600160a01b03611c8360019487614761565b51166157fa565b5001611c62565b7ffd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193611cce836112f986604051938493604085526040850190614370565b908382036020850152614370565b346102945760003660031901126102945760206001600160a01b0360015416604051908152f35b34610294576020366003190112610294576004356001600160401b03811161029457611d339036906004016142b1565b611d3b614be8565b6000915b818310611d4857005b611d53838383614728565b916307842f7160e21b6001600160e01b0319611d7160a08601614ecc565b1614612156575b611d8460208401614ef2565b61ffff80611d9460408701614ef2565b1691160161ffff81116121405761ffff6127109116101561212f576001600160401b03611dc360808501614785565b16600052600a602052604060002093611ddf60e0850185614ba2565b6001600160401b0381969296116108db57611e0a81611e0160028a015461463e565b60028a016148ad565b600095601f82116001146120c157611e3e929394959682916000926120b65750508160011b916000199060031b1c19161790565b60028601555b611e5060208201614ef2565b600186019081549060ff64ffff000000611e6c60408701614ef2565b60181b16611e7986614f01565b15159260c08701358b5568ffffffff0000000000611e9960608901614ee1565b60281b169462ffff0068ffffffff0000000000199260081b169064ffffffffff1916171617911617179055611ed2610100820182614ba2565b6001600160401b0381979297116108db57611efd81611ef4600385015461463e565b600385016148ad565b6000601f821160011461201f57927f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd926003611f5c8461ffff95611ff09860019b9c9d6000926120145750508160011b916000199060031b1c19161790565b9101555b6001600160401b03611f7460808301614785565b611f8060208401614ef2565b9263ffffffff611f9260408301614ef2565b91611fa060e0820182614ba2565b9990611fae60a08401614ecc565b9a60c0611fbd60608601614ee1565b948b611fc882614f01565b986040519d8e9d168d521660208c0152013560408a015260e060608a015260e08901916147c2565b988260e01b1660808701521660a0850152151560c084015216930390a20191611d3f565b013590508d80610756565b6003830181526020812090805b601f198416811061209e5750926003600184611ff09794829a9b9c7f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd9861ffff98601f19811610612086575b505050811b01910155611f60565b013560001983861b60f8161c191690558c8080612078565b9091602060018192858d01358155019301910161202c565b013590508880610756565b6002880187526020872090875b601f198416891061211757600194959697985083601f198116106120fd575b505050811b016002860155611e44565b0135600019600384901b60f8161c191690558780806120ed565b916020600181928585013581550193019801976120ce565b631c1604c160e11b60005260046000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff61216760608501614ee1565b1615611d7857631c1604c160e11b60005260046000fd5b346102945760203660031901126102945760206121a46001600160401b03610ebf6140b5565b6040519015158152f35b346102945760c036600319011261029457600435602435604435906001600160401b03821680920361029457606435916084359060ff821682036102945760a435916001600160a01b038316918284036102945780600052600a60205260ff600160406000200154166123b4575b5061224060405183602082015260208152612238604082614220565b8287876148f2565b860361239f5785600052600d60205260406000206001600160a01b0360405191612269836141b4565b5461227760ff82168461471c565b60081c16602082015251946003861015610e945760009561238b579061229c91614df6565b92604051956122aa876141b4565b600187526020870196338852818752600d602052604087209051976003891015612377578798612374985060ff801984541691161782555174ffffffffffffffffffffffffffffffffffffffff0082549160081b169074ffffffffffffffffffffffffffffffffffffffff0019161790556040519285845260208401527fd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf160403394a4337f00000000000000000000000000000000000000000000000000000000000000006149cd565b80f35b634e487b7160e01b88526021600452602488fd5b6326e46de360e21b86526004879052602486fd5b856332d4dea960e21b60005260045260246000fd5b6123cb33600052600c602052604060002054151590565b61221c57636c46a9b560e01b6000526004523360245260446000fd5b34610294576123f5366143ad565b6001600160a01b0360045416330361291b5760a0813603126102945760405161241d81614199565b8135815261242d602083016140cb565b906020810191825260408301356001600160401b038111610294576124559036908501614293565b916040820192835260608401356001600160401b0381116102945761247d9036908601614293565b93606083019485526080810135906001600160401b038211610294570136601f820112156102945780356124b0816143e0565b916124be6040519384614220565b81835260208084019260061b8201019036821161029457602001915b8183106128e3575050506080830152516001600160401b0381169151925193519182518301946020860193602081880312610294576020810151906001600160401b03821161029457019560a09087900312610294576040519261253d84614199565b6020870151845261255060408801614da2565b916020850192835261256460608901614da2565b916040860192835260808901519860ff8a168a036102945760608701998a5260a08101516001600160401b03811161029457602091010187601f820112156102945780516125b181614241565b986125bf6040519a8b614220565b818a5260208284010111610294576125dd916020808b019101614136565b6080868101978852604051632cbc26bb60e01b815291901b67ffffffffffffffff60801b1660048201526020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa908115610c00576000916128c4575b50611a5e576126578185614602565b1561289f57509560ff61267e6126b9936126ad989961ffff80895193511691511691615174565b6126a8612693889a939a518587511690614df6565b996126a18587511684614df6565b98516145cf565b6145cf565b915116855191886148f2565b9384600052600d602052604060002091604051926126d6846141b4565b54936126e560ff86168561471c565b6001600160a01b03602085019560081c16855286600052600d6020526040600020600260ff1982541617905561271b8383615c38565b6000948451600381101561288b576127bf575050600094516020818051810103126127bb5760200151906001600160a01b0382168092036127bb579061277092918652600f60205260ff60408720541661553b565b51906003821015610e94576127b86060927f33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c9460405193845260208401526040830190614496565ba3005b8580fd5b94909550839291925160038110156128775760010361286357506127eb856001600160a01b03926145cf565b93511690600052600f60205261281560ff6040600020541661280d86866145dc565b90309061553b565b612821846010546145dc565b6010558280612832575b5050612770565b61285c917f0000000000000000000000000000000000000000000000000000000000000000614b69565b858261282b565b6358cb522560e11b81526004879052602490fd5b634e487b7160e01b82526021600452602482fd5b634e487b7160e01b87526021600452602487fd5b6040516324eb47e560e01b81526020600482015290819061147f906024830190614159565b6128dd915060203d602011610dc457610db58183614220565b89612648565b60408336031261029457602060409182516128fd816141b4565b61290686614122565b815282860135838201528152019201916124da565b6335fdcccd60e21b6000523360045260246000fd5b34610294576020366003190112610294577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917460206001600160a01b0361297461410c565b61297c614be8565b16806001600160a01b03196009541617600955604051908152a1005b34610294576000366003190112610294576000546001600160a01b0381163303612a0f5760015490336001600160a01b03198316176001556001600160a01b0319166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b63015aa1e360e11b60005260046000fd5b346102945760203660031901126102945760206121a46001600160a01b03612a4661410c565b16600052600c602052604060002054151590565b34610294576020366003190112610294576001600160401b03612a7b6140b5565b606060c0604051612a8b816141ea565b600081526000602082015260006040820152600083820152600060808201528260a0820152015216600052600a60205260606040600020611151612acd6155dc565b611cce604051612adc816141ea565b84548152612b8960018601549563ffffffff602084019760ff81161515895261ffff60408601818360081c168152818c880191818560181c1683528560808a019560281c168552612b426003612b3460028a01614678565b9860a08c01998a5201614678565b9860c08101998a526040519e8f9e8f9260408452516040840152511515910152511660808c0152511660a08a0152511660c08801525160e080880152610120870190614159565b9051858203603f1901610100870152614159565b346102945760003660031901126102945760206001600160a01b0360095416604051908152f35b3461029457602036600319011261029457612bdd61410c565b612be5614be8565b6001600160a01b0380600e54921691828219821617600e55167f64187bd7b97e66658c91904f3021d7c28de967281d18b1a20742348afdd6a6b3600080a3005b346102945760003660031901126102945761115161154b6155dc565b3461029457612c4f36614331565b612c5a929192614be8565b6001600160401b038216612c7b816000526006602052604060002054151590565b15612c96575061108a92612c9091369161425c565b90614c0d565b631e670e4b60e01b60005260045260246000fd5b3461029457612cd2612cda612cbe366142e1565b9491612ccb939193614be8565b36916143f7565b9236916143f7565b7f000000000000000000000000000000000000000000000000000000000000000015612de05760005b8251811015612d6957806001600160a01b03612d2160019386614761565b5116612d2c81615a14565b612d38575b5001612d03565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a184612d31565b5060005b815181101561108a57806001600160a01b03612d8b60019385614761565b51168015612dda57612d9c816156cc565b612da9575b505b01612d6d565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a183612da1565b50612da3565b6335f4a7b360e01b60005260046000fd5b3461029457604036600319011261029457612e0a6140b5565b6024356001600160401b03811161029457602091612e2f6121a4923690600401614293565b90614602565b346102945760003660031901126102945760206001600160a01b03600e5416604051908152f35b34610294576020366003190112610294576004356001600160401b0381116102945780600401906101006003198236030112610294576000604051612ea0816141cf565b526000604051612eaf816141cf565b52612edc612ed2612ecd612ec660c4850186614ba2565b369161425c565b6154be565b6064830135614df6565b9060848101612eea81614bd4565b6001600160a01b03807f000000000000000000000000000000000000000000000000000000000000000016911603611a8e5750602481019267ffffffffffffffff60801b612f3785614785565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610c0057600091613154575b50611a5e576001600160401b03612fa585614785565b16612fbd816000526006602052604060002054151590565b15611a4a5760206001600160a01b0360045416916044604051809481936383826b2b60e01b835260048301523360248301525afa908115610c0057600091613135575b50156119e85761300f84614785565b9061302560a4840192612e2f612ec68585614ba2565b156131075750507ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc060806001600160401b036130ab6130a5876130728861306d60209b614785565b615c38565b8361307c82614785565b16600052600f895261096988604460ff6040600020541699019861309f8a614bd4565b9061553b565b94614bd4565b936001600160a01b0360405195817f000000000000000000000000000000000000000000000000000000000000000016875233898801521660408601528560608601521692a2806040516130fe816141cf565b52604051908152f35b6131119250614ba2565b61147f6040519283926324eb47e560e01b84526020600485015260248401916147c2565b61314e915060203d602011610dc457610db58183614220565b85613000565b61316d915060203d602011610dc457610db58183614220565b85612f8f565b34610294576000366003190112610294576040516370a0823160e01b81523060048201526020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa908115610c00576000916131e7575b602061036f83601054906145cf565b90506020813d602011613213575b8161320260209383614220565b81010312610294575161036f6131d8565b3d91506131f5565b34610294576020366003190112610294576001600160401b0361323c6140b5565b16600052600f60205260ff604060002054166040516002821015610e94576020918152f35b346102945760a03660031901126102945761327a6140b5565b602435906044356001600160401b0381116102945761329d9036906004016140df565b91606435916001600160a01b038316809303610294576084356001600160401b038111610294576132d29036906004016140df565b50506132dc6145e9565b50604051936132ea8561417e565b6000855260208501926000845260408601926000845260608701600081526060608060405161331881614199565b8281528260208201528260408201526000838201520152831580156136be575b61367e57600082356136b6575b602085116136a4575b1561367e576001600160401b03831693604051632cbc26bb60e01b815267ffffffffffffffff60801b8560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610c005760009161365f575b50611a5e576133d033615bdc565b6133e7856000526006602052604060002054151590565b1561364a5784600052600a60205260406000209485548b116136315750986134f96135079260ff6135749b9c63ffffffff60018a015461344961ffff8260081c169c8d9661343e61ffff8560181c1680998c615174565b9283829352526145dc565b8d5260281c16806135df575061ffff6134b761346760038c01614678565b985b6040519761347689614199565b8852602088019c8d52604088019586526060880193857f0000000000000000000000000000000000000000000000000000000000000000168552369161425c565b9360808701948552816040519c8d986020808b01525160408a01525116606088015251166080860152511660a08401525160a060c084015260e0830190614159565b03601f198101865285614220565b60209586946040519061351a8783614220565b6000825261353660026040519761353089614199565b01614678565b8652868601526040850152606084015260808301526001600160a01b0360045416906040518097819482936320487ded60e01b845260048401614503565b03915afa928315610c00576000936135ad575b5082604094525181845161359a816141b4565b8481520190815283519283525190820152f35b9392508184813d83116135d8575b6135c58183614220565b8101031261029457604093519293613587565b503d6135bb565b6134b761ffff91604051906135f3826141b4565b8152602081016001815260405191630181dcf160e41b602084015251602483015251151560448201526044815261362b606482614220565b98613469565b8a906358dd87c560e01b60005260045260245260446000fd5b846354c8163f60e11b60005260045260246000fd5b613678915060203d602011610dc457610db58183614220565b8b6133c2565b508261147f60405192839263a3c8cf0960e01b84526020600485015260248401916147c2565b60208301351561334e5750600161334e565b506001613345565b5060408411613338565b34610294576000366003190112610294576020601054604051908152f35b3461029457600036600319011261029457602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b3461029457602036600319011261029457602061373f61410c565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102945760003660031901126102945760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b34610294576000366003190112610294576111516040516137dd606082614220565b603381527f4879627269645769746845787465726e616c4d696e746572466173745472616e60208201527f73666572546f6b656e506f6f6c20312e362e30000000000000000000000000006040820152604051918291602083526020830190614159565b346102945760203660031901126102945760043580156103f357613863614b55565b6040516370a0823160e01b81523060048201527f0000000000000000000000000000000000000000000000000000000000000000906020816024816001600160a01b0386165afa908115610c0057600091613929575b506010546138c781856145dc565b821061390757836138d9813386614b69565b6040519081527fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171960203392a2005b61391190846145dc565b9063a17e11d560e01b60005260045260245260446000fd5b906020823d602011613953575b8161394360209383614220565b81010312610bf5575051836138b9565b3d9150613936565b60c03660031901126102945761396f6140b5565b6064356001600160401b0381116102945761398e9036906004016140df565b60843592916001600160a01b03841684036102945760a4356001600160401b038111610294576139c29036906004016140df565b5050604051906139d18261417e565b60008252600060208301526000604083015260006060830152606060806040516139fa81614199565b828152826020820152826040820152600083820152015282158015613ff7575b613fb85760008135613fef575b60208411613fdd575b15613fb857604051632cbc26bb60e01b8152608085901b67ffffffffffffffff60801b1660048201526020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa908115610c0057600091613f99575b50611a5e57613aa833615bdc565b613ac86001600160401b0385166000526006602052604060002054151590565b15613f7b576001600160401b038416600052600a6020526040600020805460243511613f575760018101549561ffff8760081c169663ffffffff61ffff8260181c1691613b2d613b1b848c602435615174565b908160608b01528060408b01526145dc565b602088015260281c1680613f0c5750613b4860038401614678565b905b60405190613b5782614199565b602435825260208201998a5260408201908152606082017f000000000000000000000000000000000000000000000000000000000000000060ff16815236613ba0908a8961425c565b91608084019283526040519b8c9460208601602090525160408601525161ffff1660608501525161ffff1660808401525160ff1660a08301525160c0820160a0905260e08201613bef91614159565b03601f1981018952613c019089614220565b604051602098613c118a83614220565b60008252613c2760026040519661353088614199565b85528985015260408401526001600160a01b038216606084015260808301526001600160a01b03600454168760405180926320487ded60e01b82528180613c72888d60048401614503565b03915afa908115610c0057600091613edf575b508452613c9460243587614939565b60208401516044358111613ec55750613cd160243530337f00000000000000000000000000000000000000000000000000000000000000006149cd565b6001600160401b038616600052600f8752613cf660243560ff60406000205416615307565b6001600160a01b038116613e72575b506001600160a01b0360045416948660405180976396f4e9f960e01b82528180613d33878760048401614503565b039134905af1948515610c00578796600096613e3d575b5091613e328697927ffa7d3740fa7611df3f0d8d8c3aa1ed57c4fffaf2dcd0c47535f18a4774b44acd946001600160401b03613e22613d8f60208b01516024356145cf565b95606060408c01519b0151905190613dd68d613dad8d88369161425c565b908a7f0000000000000000000000000000000000000000000000000000000000000000916148f2565b9b604051998a998a5260ff7f000000000000000000000000000000000000000000000000000000000000000016908a01526040890152606088015260c0608088015260c0870190614159565b9285840360a087015216966147c2565b0390a4604051908152f35b878193989297503d8311613e6b575b613e568183614220565b81010312610294575186959094613e32613d4a565b503d613e4c565b613ebf90613e8c855130336001600160a01b0385166149cd565b8451906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000009116614a1e565b86613d05565b6361acdb9360e01b60005260045260443560245260446000fd5b90508781813d8311613f05575b613ef68183614220565b81010312610294575188613c85565b503d613eec565b60405190613f19826141b4565b8152602081016001815260405191630181dcf160e41b6020840152516024830152511515604482015260448152613f51606482614220565b90613b4a565b6001600160401b03856358dd87c560e01b6000521660045260243560245260446000fd5b6001600160401b03846354c8163f60e11b6000521660045260246000fd5b613fb2915060203d602011610dc457610db58183614220565b86613a9a565b8261147f60405192839263a3c8cf0960e01b84526020600485015260248401916147c2565b602082013515613a3057506001613a30565b506001613a27565b5060408311613a1a565b34610294576020366003190112610294576004359063ffffffff60e01b82168092036102945760209163f6f46ff960e01b8114908115614075575b811561404a575b5015158152f35b6385572ffb60e01b811491508115614064575b5083614043565b6301ffc9a760e01b1490508361405d565b905063aff2afbf60e01b811480156140a5575b8015614095575b9061403c565b506301ffc9a760e01b811461408f565b50630e64dd2960e01b8114614088565b600435906001600160401b038216820361029457565b35906001600160401b038216820361029457565b9181601f84011215610294578235916001600160401b038311610294576020838186019501011161029457565b600435906001600160a01b038216820361029457565b35906001600160a01b038216820361029457565b60005b8381106141495750506000910152565b8181015183820152602001614139565b9060209161417281518092818552858086019101614136565b601f01601f1916010190565b608081019081106001600160401b038211176108db57604052565b60a081019081106001600160401b038211176108db57604052565b604081019081106001600160401b038211176108db57604052565b602081019081106001600160401b038211176108db57604052565b60e081019081106001600160401b038211176108db57604052565b606081019081106001600160401b038211176108db57604052565b90601f801991011681019081106001600160401b038211176108db57604052565b6001600160401b0381116108db57601f01601f191660200190565b92919261426882614241565b916142766040519384614220565b829481845281830111610294578281602093846000960137010152565b9080601f83011215610294578160206142ae9335910161425c565b90565b9181601f84011215610294578235916001600160401b038311610294576020808501948460051b01011161029457565b6040600319820112610294576004356001600160401b038111610294578161430b916004016142b1565b92909291602435906001600160401b0382116102945761432d916004016142b1565b9091565b906040600319830112610294576004356001600160401b03811681036102945791602435906001600160401b0382116102945761432d916004016140df565b906020808351928381520192019060005b81811061438e5750505090565b82516001600160a01b0316845260209384019390920191600101614381565b602060031982011261029457600435906001600160401b0382116102945760a09082900360031901126102945760040190565b6001600160401b0381116108db5760051b60200190565b929190614403816143e0565b936144116040519586614220565b602085838152019160051b810192831161029457905b82821061443357505050565b6020809161444084614122565b815201910190614427565b9080601f83011215610294578160206142ae933591016143f7565b9181601f84011215610294578235916001600160401b038311610294576020808501946060850201011161029457565b906003821015610e945752565b35906001600160801b038216820361029457565b9190826060910312610294576040516144cf81614205565b809280359081151582036102945760406144fe91819385526144f3602082016144a3565b6020860152016144a3565b910152565b906001600160401b039093929316815260406020820152614549614533845160a0604085015260e0840190614159565b6020850151838203603f19016060850152614159565b90604084015191603f198282030160808301526020808451928381520193019060005b8181106145a4575050506080846001600160a01b0360606142ae969701511660a084015201519060c0603f1982850301910152614159565b825180516001600160a01b03168652602090810151818701526040909501949092019160010161456c565b9190820391821161214057565b9190820180921161214057565b604051906145f6826141b4565b60006020838281520152565b906001600160401b036142ae92166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b90600182811c9216801561466e575b602083101461465857565b634e487b7160e01b600052602260045260246000fd5b91607f169161464d565b906040519182600082549261468c8461463e565b80845293600181169081156146fa57506001146146b3575b506146b192500383614220565b565b90506000929192526020600020906000915b8183106146de5750509060206146b192820101386146a4565b60209193508060019154838589010152019101909184926146c5565b9050602092506146b194915060ff191682840152151560051b820101386146a4565b6003821015610e945752565b919081101561474b5760051b8101359061011e1981360301821215610294570190565b634e487b7160e01b600052603260045260246000fd5b805182101561474b5760209160051b010190565b919081101561474b5760051b0190565b356001600160401b03811681036102945790565b919081101561474b576060020190565b604051906147b6826141b4565b60606020838281520152565b908060209392818452848401376000828201840152601f01601f1916010190565b604051906147f082614199565b60006080838281528260208201528260408201528260608201520152565b9060405161481b81614199565b60806001829460ff81546001600160801b038116865263ffffffff81861c16602087015260a01c161515604085015201546001600160801b0381166060840152811c910152565b6001600160401b031660005260076020526142ae6004604060002001614678565b81811061488e575050565b60008155600101614883565b8181029291811591840414171561214057565b9190601f81116148bc57505050565b6146b1926000526020600020906020601f840160051c830193106148e8575b601f0160051c0190614883565b90915081906148db565b92906149256149339260ff60405195869460208601988952604086015216606084015260808084015260a0830190614159565b03601f198101835282614220565b51902090565b6001600160401b037fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da8178944911691826000526007602052806149a960406000206001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001692839161519b565b604080516001600160a01b039092168252602082019290925290819081015b0390a2565b6040516323b872dd60e01b60208201526001600160a01b03928316602482015292909116604483015260648201929092526146b191614a1982608481015b03601f198101845283614220565b6153b8565b91909181158015614ad4575b15614a695760405163095ea7b360e01b60208201526001600160a01b03909316602484015260448301919091526146b19190614a198260648101614a0b565b60405162461bcd60e51b815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608490fd5b50604051636eb1769f60e11b81523060048201526001600160a01b0384166024820152602081806044810103816001600160a01b0386165afa908115610c0057600091614b23575b5015614a2a565b90506020813d602011614b4d575b81614b3e60209383614220565b81010312610294575138614b1c565b3d9150614b31565b6001600160a01b03600e5416330361108c57565b60405163a9059cbb60e01b60208201526001600160a01b03909216602483015260448201929092526146b191614a198260648101614a0b565b903590601e198136030182121561029457018035906001600160401b0382116102945760200191813603831361029457565b356001600160a01b03811681036102945790565b6001600160a01b03600154163303614bfc57565b6315ae3a6f60e11b60005260046000fd5b90805115610910576001600160401b0381516020830120921691826000526007602052614c4181600560406000200161577d565b15614d77576000526008602052604060002081516001600160401b0381116108db57614c7781614c71845461463e565b846148ad565b6020601f8211600114614ced5791614ccc827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea95936149c895600091614ce2575b508160011b916000199060031b1c19161790565b9055604051918291602083526020830190614159565b905084015138614cb8565b601f1982169083600052806000209160005b818110614d5f5750926149c89492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea989610614d46575b5050811b019055611322565b85015160001960f88460031b161c191690553880614d3a565b9192602060018192868a015181550194019201614cff565b509061147f604051928392631c9dc56960e11b84526004840152604060248401526044830190614159565b519061ffff8216820361029457565b9060ff8091169116039060ff821161214057565b60ff16604d811161214057600a0a90565b8115614de0570490565b634e487b7160e01b600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff811692828414614ec557828411614e9b5790614e3b91614db1565b91604d60ff8416118015614e80575b614e6357505090614e5d6142ae92614dc5565b9061489a565b90915063a9cb113d60e01b60005260045260245260445260646000fd5b50614e8a83614dc5565b8015614de057600019048411614e4a565b614ea491614db1565b91604d60ff841611614e6357505090614ebf6142ae92614dc5565b90614dd6565b5050505090565b356001600160e01b0319811681036102945790565b3563ffffffff811681036102945790565b3561ffff811681036102945790565b3580151581036102945790565b6001600160401b03166000818152600660205260409020549092919015614ffd5791614ffa60e092614fcf85614f647f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b9761509d565b846000526007602052614f7b81604060002061588e565b614f848361509d565b846000526007602052614f9e83600260406000200161588e565b60405194855260208501906001600160801b0360408092805115158552826020820151166020860152015116910152565b60808301906001600160801b0360408092805115158552826020820151166020860152015116910152565ba1565b82631e670e4b60e01b60005260045260246000fd5b61501a6147e3565b506001600160801b036060820151166001600160801b038083511691615065602085019361505f61505263ffffffff875116426145cf565b856080890151169061489a565b906145dc565b8082101561507e57505b16825263ffffffff4216905290565b905061506f565b90816020910312610294575180151581036102945790565b805115615104576001600160801b036040820151166001600160801b03602083015116106150c85750565b60408051632008344960e21b815282511515600482015260208301516001600160801b0390811660248301529190920151166044820152606490fd5b6001600160801b036040820151161580159061515e575b6151225750565b604080516335a2be7360e21b815282511515600482015260208301516001600160801b0390811660248301529190920151166044820152606490fd5b506001600160801b03602082015116151561511b565b6151979061ffff61271061518e828296989798168461489a565b0495169061489a565b0490565b9182549060ff8260a01c161580156152ff575b6152f9576001600160801b03821691600185019081546151e163ffffffff6001600160801b0383169360801c16426145cf565b9081615299575b5050848110615273575083831061522157505061520e6001600160801b039283926145cf565b16166001600160801b0319825416179055565b5460801c9161523081856145cf565b9260001981019080821161214057615253615258926001600160a01b03966145dc565b614dd6565b636864691d60e11b6000526004526024521660445260646000fd5b82856001600160a01b0392630d3b2b9560e11b6000526004526024521660445260646000fd5b8286929396116152e8576152b49261505f9160801c9061489a565b808410156152e35750825b855463ffffffff60801b19164260801b63ffffffff60801b161786559238806151e8565b6152bf565b634b92ca1560e11b60005260046000fd5b50505050565b5082156151ae565b6002811015610e945760011461531a5750565b60206001600160a01b039160247f00000000000000000000000000000000000000000000000000000000000000009161537481847f0000000000000000000000000000000000000000000000000000000000000000614a1e565b60006040519586948593630852cd8d60e31b85526004850152165af18015610c005761539d5750565b6153b59060203d602011610dc457610db58183614220565b50565b6001600160a01b0361543a9116916040926000808551936153d98786614220565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d156154b6573d9161541e83614241565b9261542b87519485614220565b83523d6000602085013e615cab565b8051908161544757505050565b602080615458938301019101615085565b156154605750565b5162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b6064820152608490fd5b606091615cab565b80518015615515576020036154f057805160208281019183018390031261029457519060ff82116154f0575060ff1690565b60405163953576f760e01b81526020600482015290819061147f906024830190614159565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9190916002811015610e9457615575576146b1917f0000000000000000000000000000000000000000000000000000000000000000614b69565b6040516340c10f1960e01b81526001600160a01b03909216600483015260248201526020818060448101038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015610c005761539d5750565b60405190600b548083528260208101600b60005260206000209260005b81811061560e5750506146b192500383614220565b84548352600194850194879450602090930192016155f9565b906040519182815491828252602082019060005260206000209260005b8181106156595750506146b192500383614220565b8454835260019485019487945060209093019201615644565b805482101561474b5760005260206000200190600090565b805490680100000000000000008210156108db57816156b19160016156c894018155615672565b819391549060031b91821b91600019901b19161790565b9055565b80600052600360205260406000205415600014615705576156ee81600261568a565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146157055761572d81600b61568a565b600b5490600052600c602052604060002055600190565b806000526006602052604060002054156000146157055761576681600561568a565b600554906000526006602052604060002055600190565b60008281526001820160205260409020546157b4578061579f8360019361568a565b80549260005201602052604060002055600190565b5050600090565b805480156157e45760001901906157d28282615672565b8154906000199060031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c602052604090205480156157b457600019810181811161214057600b5460001981019190821161214057808203615854575b505050615840600b6157bb565b600052600c60205260006040812055600190565b6158766158656156b193600b615672565b90549060031b1c928392600b615672565b9055600052600c602052604060002055388080615833565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199161595c60609280546158cb63ffffffff8260801c16426145cf565b9081615992575b50506001600160801b03600181602086015116928281541680851060001461598a57508280855b1616831982541617815561592886511515829081549060ff60a01b90151560a01b169060ff60a01b1916179055565b60408601516fffffffffffffffffffffffffffffffff1960809190911b16939092166001600160801b031692909217910155565b614ffa60405180926001600160801b0360408092805115158552826020820151166020860152015116910152565b8380916158f9565b6001600160801b03916159be8392836159b76001880154948286169560801c9061489a565b91166145dc565b80821015615a0d57505b835463ffffffff60801b1992909116929092161673ffffffffffffffffffffffffffffffffffffffff19909116174260801b63ffffffff60801b1617815538806158d2565b90506159c8565b60008181526003602052604090205480156157b45760001981018181116121405760025460001981019190821161214057818103615a6e575b505050615a5a60026157bb565b600052600360205260006040812055600190565b615a90615a7f6156b1936002615672565b90549060031b1c9283926002615672565b90556000526003602052604060002055388080615a4d565b60008181526006602052604090205480156157b45760001981018181116121405760055460001981019190821161214057818103615b02575b505050615aee60056157bb565b600052600660205260006040812055600190565b615b24615b136156b1936005615672565b90549060031b1c9283926005615672565b90556000526006602052604060002055388080615ae1565b906001820191816000528260205260406000205490811515600014615bd3576000198201918083116121405781546000198101908111612140578381615b8a9503615b9c575b5050506157bb565b60005260205260006040812055600190565b615bbc615bac6156b19386615672565b90549060031b1c92839286615672565b905560005284602052604060002055388080615b82565b50505050600090565b7f0000000000000000000000000000000000000000000000000000000000000000615c045750565b6001600160a01b031680600052600360205260406000205415615c245750565b6368692cbb60e11b60005260045260246000fd5b6001600160401b037f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c911691826000526007602052806149a960026040600020016001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001692839161519b565b91929015615d0d5750815115615cbf575090565b3b15615cc85790565b60405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606490fd5b825190915015615d205750805190602001fd5b60405162461bcd60e51b81526020600482015290819061147f90602483019061415956fea164736f6c634300081a000a",
}

var HybridWithExternalMinterFastTransferTokenPoolABI = HybridWithExternalMinterFastTransferTokenPoolMetaData.ABI

var HybridWithExternalMinterFastTransferTokenPoolBin = HybridWithExternalMinterFastTransferTokenPoolMetaData.Bin

func DeployHybridWithExternalMinterFastTransferTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, minter common.Address, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *HybridWithExternalMinterFastTransferTokenPool, error) {
	parsed, err := HybridWithExternalMinterFastTransferTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HybridWithExternalMinterFastTransferTokenPoolBin), backend, minter, token, localTokenDecimals, allowlist, rmnProxy, router)
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipSendToken", destinationChainSelector, amount, maxFastTransferFee, receiver, settlementFeeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, settlementFeeToken, extraArgs)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, settlementFeeToken, extraArgs)
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) UpdateGroups(opts *bind.TransactOpts, groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "updateGroups", groupUpdates)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) UpdateGroups(groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateGroups(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, groupUpdates)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) UpdateGroups(groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.UpdateGroups(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, groupUpdates)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactor) WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "withdrawLiquidity", amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.WithdrawLiquidity(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolTransactorSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterFastTransferTokenPool.Contract.WithdrawLiquidity(&_HybridWithExternalMinterFastTransferTokenPool.TransactOpts, amount)
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
	FillerFee                *big.Int
	PoolFee                  *big.Int
	DestinationPool          []byte
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterGroupUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "GroupUpdated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "GroupUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchGroupUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolGroupUpdated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "GroupUpdated", remoteChainSelectorRule, groupRule)
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

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) FilterRebalancerSet(opts *bind.FilterOpts, oldRebalancer []common.Address, newRebalancer []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator, error) {

	var oldRebalancerRule []interface{}
	for _, oldRebalancerItem := range oldRebalancer {
		oldRebalancerRule = append(oldRebalancerRule, oldRebalancerItem)
	}
	var newRebalancerRule []interface{}
	for _, newRebalancerItem := range newRebalancer {
		newRebalancerRule = append(newRebalancerRule, newRebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RebalancerSet", oldRebalancerRule, newRebalancerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator{contract: _HybridWithExternalMinterFastTransferTokenPool.contract, event: "RebalancerSet", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterFastTransferTokenPool *HybridWithExternalMinterFastTransferTokenPoolFilterer) WatchRebalancerSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRebalancerSet, oldRebalancer []common.Address, newRebalancer []common.Address) (event.Subscription, error) {

	var oldRebalancerRule []interface{}
	for _, oldRebalancerItem := range oldRebalancer {
		oldRebalancerRule = append(oldRebalancerRule, oldRebalancerItem)
	}
	var newRebalancerRule []interface{}
	for _, newRebalancerItem := range newRebalancer {
		newRebalancerRule = append(newRebalancerRule, newRebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RebalancerSet", oldRebalancerRule, newRebalancerRule)
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
	return common.HexToHash("0xfa7d3740fa7611df3f0d8d8c3aa1ed57c4fffaf2dcd0c47535f18a4774b44acd")
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

	CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error)

	FastFill(opts *bind.TransactOpts, fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDestChainConfig(opts *bind.TransactOpts, destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error)

	UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error)

	UpdateGroups(opts *bind.TransactOpts, groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error)

	WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

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

	FilterGroupUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterFastTransferTokenPoolGroupUpdatedIterator, error)

	WatchGroupUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolGroupUpdated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error)

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

	FilterRebalancerSet(opts *bind.FilterOpts, oldRebalancer []common.Address, newRebalancer []common.Address) (*HybridWithExternalMinterFastTransferTokenPoolRebalancerSetIterator, error)

	WatchRebalancerSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterFastTransferTokenPoolRebalancerSet, oldRebalancer []common.Address, newRebalancer []common.Address) (event.Subscription, error)

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
