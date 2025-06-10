// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_with_external_minter_fast_transfer_token_pool

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
	OriginalSender      []byte
	RemoteChainSelector uint64
	Receiver            common.Address
	Amount              *big.Int
	LocalToken          common.Address
	SourcePoolAddress   []byte
	SourcePoolData      []byte
	OffchainTokenData   []byte
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

var BurnMintWithExternalMinterFastTransferTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSendToken\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"computeFillId\",\"inputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"fastFill\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAccumulatedPoolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedFillers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCcipSendTokenFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIFastTransferPool.Quote\",\"components\":[{\"name\":\"ccipSettlementFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfig\",\"components\":[{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFillInfo\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.FillInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIFastTransferPool.FillState\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAllowedFiller\",\"inputs\":[{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateDestChainConfig\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfigUpdateArgs[]\",\"components\":[{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFillerAllowList\",\"inputs\":[{\"name\":\"fillersToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"fillersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPoolFees\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestinationPoolUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destinationPool\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"filler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferRequested\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fillerReimbursementAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFeeAccumulated\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"prevState\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIFastTransferPool.FillState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FillerAllowListUpdated\",\"inputs\":[{\"name\":\"addFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"removeFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolFeeWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadySettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"FillerNotAllowlisted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientPoolFees\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFillId\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferAmountExceedsMaxFillAmount\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x61014080604052346104395761645e803803809161001d82856104fe565b8339810160a0828203126104395761003482610521565b9061004160208401610535565b60408401516001600160401b0381116104395784019180601f84011215610439578251926001600160401b0384116104e8578360051b90602082019461008a60405196876104fe565b855260208086019282010192831161043957602001905b8282106104d0575050506100c360806100bc60608701610521565b9501610521565b6040516321df0da760e01b81526001600160a01b03909416949093602081600481895afa9081156104c45760009161048a575b506001600160a01b031690331561047957600180546001600160a01b0319163317905581158015610468575b8015610457575b610446578160209160049360805260c0526040519283809263313ce56760e01b82525afa60009181610405575b506103da575b5060a052600480546001600160a01b0319166001600160a01b0384169081179091558151151560e08190529091906102b7575b50156102a1576101005261012052604051615d7a90816106e48239608051818181611370015281816113da015281816115af01528181611d01015281816126b101528181612c29015281816133280152818161353a01528181613a9501528181613ae2015281816141ec01528181614bde01528181614ccc0152615037015260a0518181816115f301528181611b85015281816137d301528181613a4b015281816147bd0152614807015260c051818181610b4c0152818161146801528181611ab301528181612a0d015281816133b601526136cc015260e051818181610b07015281816131570152615a6601526101005181611e7d01526101205181818161022301528181614947015261500f0152f35b6335fdcccd60e21b600052600060045260246000fd5b602091604051916102c884846104fe565b60008352600036813760e051156103c95760005b8351811015610343576001906001600160a01b036102fa8287610543565b51168661030682610585565b610313575b5050016102dc565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a1388661030b565b5091509260005b82518110156103be576001906001600160a01b036103688286610543565b511680156103b8578561037a82610683565b610388575b50505b0161034a565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a1388561037f565b50610382565b50929150503861018f565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036103ee575061015c565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d60201161043e575b81610421602093836104fe565b810103126104395761043290610535565b9038610156565b600080fd5b3d9150610414565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03811615610129565b506001600160a01b03851615610122565b639b15e16f60e01b60005260046000fd5b90506020813d6020116104bc575b816104a5602093836104fe565b81010312610439576104b690610521565b386100f6565b3d9150610498565b6040513d6000823e3d90fd5b602080916104dd84610521565b8152019101906100a1565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b038211908210176104e857604052565b51906001600160a01b038216820361043957565b519060ff8216820361043957565b80518210156105575760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156105575760005260206000200190600090565b600081815260036020526040902054801561067c5760001981018181116106665760025460001981019190821161066657818103610615575b50505060025480156105ff57600019016105d981600261056d565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b61064e61062661063793600261056d565b90549060031b1c928392600261056d565b819391549060031b91821b91600019901b19161790565b905560005260036020526040600020553880806105be565b634e487b7160e01b600052601160045260246000fd5b5050600090565b806000526003602052604060002054156000146106dd57600254680100000000000000008110156104e8576106c4610637826001859401600255600261056d565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a714613b8c57508063181f5a7714613b0657806321df0da714613ac2578063240028e814613a6f57806324f65ee714613a315780632b2c0eb414613a165780632e7aa8c81461360b57806339077537146132ca5780634c5ef0ed1461328557806354c8a4f31461312557806362ddd3c4146130a25780636609f599146130865780636d3d1a581461305f5780636def4ce714612f1b57806378b410f214612ee157806379ba509714612e305780637d54534e14612db057806385572ffb146127a857806387f060d0146124fa5780638926f54f146124b55780638a18dcbd146120005780638da5cb5b14611fd957806392575c3b14611a01578063929ea5ba146118f7578063962d4020146117bb5780639a4575b91461139f5780639fe280f51461130c578063a42a7b8b146111da578063a7cd63b71461116c578063abe1c1e8146110fd578063acfecf9114610fd8578063af58d59f14610f8e578063b0f479a114610f67578063b794658014610f2f578063c0d7865514610e8b578063c4bffe2b14610d79578063c75eea9c14610cd9578063cf7401f314610b70578063dc0bd97114610b2c578063e0351e1314610aef578063e8a1da171461035c578063eeebc67414610304578063f2fde38b1461024c5763f36675171461020357600080fd5b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b34610247576020600319360112610247576001600160a01b0361026d613d37565b610275614991565b163381146102da57807fffffffffffffffffffffffff000000000000000000000000000000000000000060005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760806003193601126102475760443560ff811681036102475760643567ffffffffffffffff811161024757602091610348610354923690600401613f3d565b906024356004356146b4565b604051908152f35b346102475761036a36613f8c565b919092610375614991565b6000905b82821061094a5750505060009063ffffffff42165b81831061039757005b6103a2838386614409565b926101208436031261024757604051936103bb85613e3b565b6103c481613d78565b8552602081013567ffffffffffffffff81116102475781019336601f860112156102475784356103f38161408b565b956104016040519788613ec7565b81875260208088019260051b820101903682116102475760208101925b82841061091b575050505060208601948552604082013567ffffffffffffffff8111610247576104519036908401613f3d565b906040870191825261047b6104693660608601614164565b936060890194855260c0369101614164565b946080880195865261048d845161515f565b610497865161515f565b825151156108f1576104b367ffffffffffffffff89511661570c565b156108b85767ffffffffffffffff885116600052600760205260406000206105d185516001600160801b03604082015116906105956001600160801b036020830151169151151583608060405161050981613e3b565b858152602081018a905260408101849052606081018690520152855474ff000000000000000000000000000000000000000091151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166001600160801b0384161773ffffffff00000000000000000000000000000000608089901b1617178555565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166001600160801b0391909116176001830155565b6106d387516001600160801b03604082015116906106976001600160801b036020830151169151151583608060405161060981613e3b565b858152602081018a9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166001600160801b0385161773ffffffff0000000000000000000000000000000060808a901b161791151560a01b74ff000000000000000000000000000000000000000016919091179055565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166001600160801b0391909116176003830155565b6004845191019080519067ffffffffffffffff82116108a257610700826106fa855461431f565b8561466f565b602090601f831160011461083b57610730929160009183610830575b50506000198260011b9260031b1c19161790565b90555b60005b8751805182101561076b579061076560019261075e8367ffffffffffffffff8e511692614570565b51906149cf565b01610736565b505097967f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c293929196509461082567ffffffffffffffff60019751169251935191516107fa6107ce60405196879687526101006020880152610100870190613d12565b9360408601906001600160801b0360408092805115158552826020820151166020860152015116910152565b60a08401906001600160801b0360408092805115158552826020820151166020860152015116910152565b0390a101919261038e565b015190508d8061071c565b90601f1983169184600052816000209260005b81811061088a5750908460019594939210610871575b505050811b019055610733565b015160001960f88460031b161c191690558c8080610864565b9293602060018192878601518155019501930161084e565b634e487b7160e01b600052604160045260246000fd5b67ffffffffffffffff8851167f1d5ad3c50000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b833567ffffffffffffffff81116102475760209161093f8392833691870101613f3d565b81520193019261041e565b9092919367ffffffffffffffff61096a610965868886614584565b6142cd565b169261097584615b6d565b15610ac15783600052600760205261099360056040600020016152bb565b9260005b84518110156109cf576001908660005260076020526109c860056040600020016109c18389614570565b5190615c01565b5001610997565b5093909491959250806000526007602052600560406000206000815560006001820155600060028201556000600382015560048101610a0e815461431f565b9081610a7e575b5050018054906000815581610a5d575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a1019091610379565b6000526020600020908101905b81811015610a255760008155600101610a6a565b81601f60009311600114610a965750555b8880610a15565b81835260208320610ab191601f01861c810190600101614645565b8082528160208120915555610a8f565b837f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102475760006003193601126102475760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760e060031936011261024757610b89613d61565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc36011261024757604051610bbf81613eab565b60243580151581036102475781526044356001600160801b03811681036102475760208201526064356001600160801b038116810361024757604082015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7c3601126102475760405190610c3482613eab565b608435801515810361024757825260a4356001600160801b038116810361024757602083015260c4356001600160801b03811681036102475760408301526001600160a01b036009541633141580610cc4575b610c9657610c9492614ee0565b005b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b506001600160a01b0360015416331415610c87565b346102475760206003193601126102475767ffffffffffffffff610cfb613d61565b610d036145a4565b50166000526007602052610d75610d25610d2060406000206145cf565b6150ec565b6040519182918291909160806001600160801b038160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b34610247576000600319360112610247576040516005548082528160208101600560005260206000209260005b818110610e72575050610dbb92500382613ec7565b805190610de0610dca8361408b565b92610dd86040519485613ec7565b80845261408b565b90601f1960208401920136833760005b8151811015610e22578067ffffffffffffffff610e0f60019385614570565b5116610e1b8287614570565b5201610df0565b5050906040519182916020830190602084525180915260408301919060005b818110610e4f575050500390f35b825167ffffffffffffffff16845285945060209384019390920191600101610e41565b8454835260019485019486945060209093019201610da6565b3461024757602060031936011261024757610ea4613d37565b610eac614991565b6001600160a01b0381169081156108f157600480547fffffffffffffffffffffffff000000000000000000000000000000000000000081169093179055604080516001600160a01b0393841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491819081015b0390a1005b3461024757602060031936011261024757610d75610f53610f4e613d61565b614623565b604051918291602083526020830190613d12565b346102475760006003193601126102475760206001600160a01b0360045416604051908152f35b346102475760206003193601126102475767ffffffffffffffff610fb0613d61565b610fb86145a4565b50166000526007602052610d75610d25610d2060026040600020016145cf565b346102475767ffffffffffffffff610fef36613fda565b929091610ffa614991565b1690611013826000526006602052604060002054151590565b156110cf578160005260076020526110446005604060002001611037368685613f06565b6020815191012090615c01565b15611088577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76919261108360405192839260208452602084019161454f565b0390a2005b6110cb906040519384937f74f23c7c000000000000000000000000000000000000000000000000000000008552600485015260406024850152604484019161454f565b0390fd5b507f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102475760206003193601126102475761111661424f565b50600435600052600d6020526040806000206001600160a01b0382519161113c83613e57565b5461114a60ff8216846143fd565b81602084019160081c1681526111638451809451614143565b51166020820152f35b34610247576000600319360112610247576040516002548082526020820190600260005260206000209060005b8181106111c457610d75856111b081870382613ec7565b60405191829160208352602083019061401b565b8254845260209093019260019283019201611199565b346102475760206003193601126102475767ffffffffffffffff6111fc613d61565b16600052600760205261121560056040600020016152bb565b805190601f1961123d6112278461408b565b936112356040519586613ec7565b80855261408b565b0160005b8181106112fb57505060005b8151811015611295578061126360019284614570565b5160005260086020526112796040600020614359565b6112838286614570565b5261128e8185614570565b500161124d565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b8282106112ce57505050500390f35b919360206112eb82603f1960019597998495030186528851613d12565b96019201920185949391926112bf565b806060602080938701015201611241565b3461024757602060031936011261024757611325613d37565b61132d614991565b6113356141b0565b908161133d57005b60206001600160a01b0382611394857f738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599957f000000000000000000000000000000000000000000000000000000000000000061509a565b6040519485521692a2005b34610247576113ad36614058565b606060206040516113bd81613e57565b8281520152608081016113cf816142b9565b6001600160a01b03807f00000000000000000000000000000000000000000000000000000000000000001691160361177c57506020810177ffffffffffffffff00000000000000000000000000000000611428826142cd565b60801b16604051907f2cbc26bb00000000000000000000000000000000000000000000000000000000825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116ea5760009161174d575b50611723576114b26114ad604084016142b9565b615a64565b67ffffffffffffffff6114c4826142cd565b166114dc816000526006602052604060002054151590565b156116f65760206001600160a01b0360045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa9081156116ea5760009161169a575b506001600160a01b0316330361166c57610f4e8161165993611565606061155b6115e9966142cd565b9201358092614c83565b61156e81614ffe565b7ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae1067ffffffffffffffff6115a1846142cd565b604080516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152336020820152908101949094521691606090a26142cd565b610d7560405160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260208152611627604082613ec7565b6040519261163484613e57565b8352602083019081526040519384936020855251604060208601526060850190613d12565b9051601f19848303016040850152613d12565b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d6020116116e2575b816116b360209383613ec7565b810103126116de5751906001600160a01b03821682036116db57506001600160a01b03611532565b80fd5b5080fd5b3d91506116a6565b6040513d6000823e3d90fd5b7fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b61176f915060203d602011611775575b6117678183613ec7565b810190614731565b83611499565b503d61175d565b61178d6001600160a01b03916142b9565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b346102475760606003193601126102475760043567ffffffffffffffff8111610247576117ec903690600401613f5b565b9060243567ffffffffffffffff81116102475761180d903690600401614112565b9060443567ffffffffffffffff81116102475761182e903690600401614112565b6001600160a01b0360095416331415806118e2575b610c96578386148015906118d8575b6118ae5760005b86811061186257005b806118a86118766109656001948b8b614584565b611881838989614594565b6118a261189a61189286898b614594565b923690614164565b913690614164565b91614ee0565b01611859565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b5080861415611852565b506001600160a01b0360015416331415611843565b346102475760406003193601126102475760043567ffffffffffffffff8111610247576119289036906004016140f7565b60243567ffffffffffffffff8111610247576119489036906004016140f7565b90611951614991565b60005b8151811015611983578061197c6001600160a01b0361197560019486614570565b51166156d3565b5001611954565b5060005b82518110156119b657806119af6001600160a01b036119a860019487614570565b51166157c1565b5001611987565b7ffd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e241936119f383610f2a8660405193849360408552604085019061401b565b90838203602085015261401b565b611a0a36613dbb565b505092909391604051611a1c81613e57565b600081526000602082015260606080604051611a3781613e3b565b828152826020820152826040820152600083820152015267ffffffffffffffff8416916040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008660801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116ea57600091611fba575b5061172357611af233615a64565b611b09836000526006602052604060002054151590565b15611f8c5782600052600a60205260406000209687548511611f59576001880154948560081c61ffff16968660181c61ffff1699611b478b8a6146fb565b61ffff16611b55908461465c565b6127109004976020870198895260281c63ffffffff168015600014611ef55750611b8160038201614359565b915b7f000000000000000000000000000000000000000000000000000000000000000096604051611bb181613e3b565b858152602081019b8c52604081019d8e52606081019d60ff8a169e8f815236611bdb908a8c613f06565b91608084019283526040519e8f9460208601602090525160408601525161ffff1660608501525161ffff1660808401525160ff1660a08301525160c0820160a0905260e08201611c2a91613d12565b03601f1981018c52611c3c908c613ec7565b6020809c604051611c4d8382613ec7565b60008082529d6001600160a01b0393611c74600260405199611c6e8b613e3b565b01614359565b88528701526040860152169384606085015260808401526001600160a01b0360045416938c60405180967f20487ded0000000000000000000000000000000000000000000000000000000082528180611cd189896004840161445f565b03915afa8015611eea57908d949392918d90611ea9575b611d7896508252611cf98784614c83565b611d258730337f0000000000000000000000000000000000000000000000000000000000000000614cf4565b611d2e87614ffe565b80611e68575b50506001600160a01b0360045416906040518095819482937f96f4e9f90000000000000000000000000000000000000000000000000000000084526004840161445f565b039134905af1978815611e5c578098611e0a575b5050869798927f240a1286fd41f1034c4032dcd6b93fc09e81be4a0b64c7ecee6260b605a8e0169492611ddb611dc7611dff948a519061452c565b94611dd3368486613f06565b90868c6146b4565b975160405195869586528c860152604085015260806060850152608084019161454f565b0390a4604051908152f35b909197508882813d8311611e55575b611e238183613ec7565b810103126116db57505195827f240a1286fd41f1034c4032dcd6b93fc09e81be4a0b64c7ecee6260b605a8e016611d8c565b503d611e19565b604051903d90823e3d90fd5b81611e79611ea29351303385614cf4565b51907f000000000000000000000000000000000000000000000000000000000000000090614d5e565b8d80611d34565b509193858d8297939597903d8411611ee2575b611ec68284613ec7565b5081010312611ede57918c9391611d78959351611ce8565b8b80fd5b3d9150611ebc565b6040513d8e823e3d90fd5b60405190611f0282613e57565b81526020810160018152604051917f181dcf10000000000000000000000000000000000000000000000000000000006020840152516024830152511515604482015260448152611f53606482613ec7565b91611b83565b5050507f58dd87c50000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b827fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b611fd3915060203d602011611775576117678183613ec7565b88611ae4565b346102475760006003193601126102475760206001600160a01b0360015416604051908152f35b346102475760206003193601126102475760043567ffffffffffffffff811161024757612031903690600401613f5b565b612039614991565b60005b81811061204557005b612050818385614409565b60a081017f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000061209f83614c29565b1614612474575b602082016120b381614c67565b9061271061ffff6120d160408701946120cb86614c67565b906146fb565b161161244a576080840167ffffffffffffffff6120ed826142cd565b16600052600a60205260406000209460e081019461210b8683614268565b600289019167ffffffffffffffff82116108a25761212d826106fa855461431f565b600090601f83116001146123e65761215c9291600091836123db5750506000198260011b9260031b1c19161790565b90555b61216884614c67565b92600188019788549861217a88614c67565b60181b64ffff000000169561218e86614c76565b151560c087013597888555606088019c6121a78e614c56565b60281b68ffffffff0000000000169360081b62ffff0016907fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016177fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff16179060ff1617179055610100840161221c9085614268565b90916003019167ffffffffffffffff82116108a25761223f826106fa855461431f565b600090601f83116001146123665791806122749261227b96959460009261235b5750506000198260011b9260031b1c19161790565b90556142cd565b9361228590614c67565b9461228f90614c67565b9561229a9083614268565b90916122a590614c29565b976122af90614c56565b926122b990614c76565b936040519761ffff899816885261ffff16602088015260408701526060860160e0905260e08601906122ea9261454f565b957fffffffff0000000000000000000000000000000000000000000000000000000016608085015263ffffffff1660a0840152151560c083015267ffffffffffffffff1692037f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd91a260010161203c565b01359050388061071c565b838252602082209a9e9d9c9b9a91601f198416815b8181106123c35750919e9f9b9c9d9e600193918561227b98979694106123a9575b505050811b0190556142cd565b60001960f88560031b161c199101351690558f808061239c565b9193602060018192878701358155019501920161237b565b013590508e8061071c565b8382526020822091601f198416815b8181106124325750908460019594939210612418575b505050811b01905561215f565b60001960f88560031b161c199101351690558d808061240b565b838301358555600190940193602092830192016123f5565b7f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b63ffffffff61248560608401614c56565b16156120a6577f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760206003193601126102475760206124f067ffffffffffffffff6124dc613d61565b166000526006602052604060002054151590565b6040519015158152f35b346102475760c06003193601126102475760043560243560443567ffffffffffffffff81169182820361024757606435926084359160ff831683036102475760a435926001600160a01b038416928385036102475780600052600a60205260ff6001604060002001541661275c575b5061258d60405184602082015260208152612585604082613ec7565b8288856146b4565b870361272e5786600052600d60205260406000206001600160a01b03604051916125b683613e57565b546125c460ff8216846143fd565b60081c16602082015251956003871015612718576000966126ec576125f3916125ec91614804565b8095614b92565b6040519561260087613e57565b600187526020870196338852818752600d6020526040872090519760038910156126d85787986126d5985060ff60ff198454169116178255517fffffffffffffffffffffff0000000000000000000000000000000000000000ff74ffffffffffffffffffffffffffffffffffffffff0083549260081b1691161790556040519285845260208401527fd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf160403394a4337f0000000000000000000000000000000000000000000000000000000000000000614cf4565b80f35b602488634e487b7160e01b81526021600452fd5b602487897fcee81443000000000000000000000000000000000000000000000000000000008252600452fd5b634e487b7160e01b600052602160045260246000fd5b867fcb537aa40000000000000000000000000000000000000000000000000000000060005260045260246000fd5b61277333600052600c602052604060002054151590565b612569577f6c46a9b5000000000000000000000000000000000000000000000000000000006000526004523360245260446000fd5b34610247576127b636614058565b6001600160a01b03600454163303612d825760a08136031261024757604051906127df82613e3b565b803582526127ef60208201613d78565b9060208301918252604081013567ffffffffffffffff8111610247576128189036908301613f3d565b9060408401918252606081013567ffffffffffffffff8111610247576128419036908301613f3d565b906060850191825260808101359067ffffffffffffffff8211610247570136601f820112156102475780356128758161408b565b916128836040519384613ec7565b81835260208084019260061b8201019036821161024757602001915b818310612d4a575050506080850152600092519367ffffffffffffffff85169051925191519182518301926020840190602081860312612d465760208101519067ffffffffffffffff8211612d4257019360a09085900312612d3e576040519061290882613e3b565b6020850151825261291b60408601614b83565b926020830193845261292f60608701614b83565b9860408401998a5260808701519660ff88168803612d3a576060850197885260a08101519067ffffffffffffffff8211612d365790602091010183601f82011215612d3a5780519061298082613eea565b9461298e6040519687613ec7565b82865260208383010111612d3657906129ad9160208087019101613cef565b6080840192835277ffffffffffffffff00000000000000000000000000000000604051917f2cbc26bb00000000000000000000000000000000000000000000000000000000835260801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115612d2b578991612d0c575b50612ce457612a4c81866142e2565b15612ca657509060ff612ad8612acd612ac18a9b612710612ab181612aa1612abc9e9f6001600160a01b0390612afa9c519050602081519101519060208110612c94575b50169b61ffff8b519151169061465c565b049261ffff89519151169061465c565b049a8b91875161452c565b61452c565b93518389511690614804565b978288511690614804565b9551166040519184602084015260208352612af4604084613ec7565b886146b4565b93848752600d602052604087209160405192612b1584613e57565b54612b2360ff8216856143fd565b6001600160a01b03602085019160081c168152889584516003811015612c8057612bd657505090612b5e91612b59828a96614b92565b6148f3565b838652600d60205260408620600260ff1982541617905551906003821015612bc25791612bbe6060927f33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c9460405193845260208401526040830190614143565ba380f35b602486634e487b7160e01b81526021600452fd5b9450945050815160038110156126d857600103612c5457612bff836001600160a01b039261452c565b935116612c15612c0f8486615306565b306148f3565b8380612c23575b5050612b5e565b612c4d917f000000000000000000000000000000000000000000000000000000000000000061509a565b8683612c1c565b602487867fb196a44a000000000000000000000000000000000000000000000000000000008252600452fd5b60248b634e487b7160e01b81526021600452fd5b6000199060200360031b1b1638612a90565b6110cb906040519182917f24eb47e5000000000000000000000000000000000000000000000000000000008352602060048401526024830190613d12565b6004887f53ad11d8000000000000000000000000000000000000000000000000000000008152fd5b612d25915060203d602011611775576117678183613ec7565b8a612a3d565b6040513d8b823e3d90fd5b8a80fd5b8980fd5b8580fd5b8780fd5b8680fd5b6040833603126102475760206040918251612d6481613e57565b612d6d86613d4d565b8152828601358382015281520192019161289f565b7fd7f73334000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b34610247576020600319360112610247577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917460206001600160a01b03612df4613d37565b612dfc614991565b16807fffffffffffffffffffffffff00000000000000000000000000000000000000006009541617600955604051908152a1005b34610247576000600319360112610247576000546001600160a01b0381163303612eb7577fffffffffffffffffffffffff0000000000000000000000000000000000000000600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760206003193601126102475760206124f06001600160a01b03612f07613d37565b16600052600c602052604060002054151590565b346102475760206003193601126102475767ffffffffffffffff612f3d613d61565b606060c0604051612f4d81613e8f565b600081526000602082015260006040820152600083820152600060808201528260a0820152015216600052600a60205260606040600020610d75612f8f615270565b6119f3604051612f9e81613e8f565b8454815261304b60018601549563ffffffff602084019760ff81161515895261ffff60408601818360081c168152818c880191818560181c1683528560808a019560281c1685526130046003612ff660028a01614359565b9860a08c01998a5201614359565b9860c08101998a526040519e8f9e8f9260408452516040840152511515910152511660808c0152511660a08a0152511660c08801525160e080880152610120870190613d12565b9051603f1986830301610100870152613d12565b346102475760006003193601126102475760206001600160a01b0360095416604051908152f35b3461024757600060031936011261024757610d756111b0615270565b34610247576130b036613fda565b6130bb929192614991565b67ffffffffffffffff82166130dd816000526006602052604060002054151590565b156130f85750610c94926130f2913691613f06565b906149cf565b7f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102475761314d61315561313936613f8c565b9491613146939193614991565b36916140a3565b9236916140a3565b7f00000000000000000000000000000000000000000000000000000000000000001561325b5760005b82518110156131e457806001600160a01b0361319c60019386614570565b51166131a781615ad9565b6131b3575b500161317e565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a1846131ac565b5060005b8151811015610c9457806001600160a01b0361320660019385614570565b511680156132555761321781615694565b613224575b505b016131e8565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a18361321c565b5061321e565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760406003193601126102475761329e613d61565b60243567ffffffffffffffff8111610247576020916132c46124f0923690600401613f3d565b906142e2565b346102475760206003193601126102475760043567ffffffffffffffff8111610247578060040190610100600319823603011261024757600060405161330f81613e73565b526084810161331d816142b9565b6001600160a01b03807f00000000000000000000000000000000000000000000000000000000000000001691160361177c57506024810177ffffffffffffffff00000000000000000000000000000000613376826142cd565b60801b16604051907f2cbc26bb00000000000000000000000000000000000000000000000000000000825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116ea576000916135ec575b506117235767ffffffffffffffff6133fe826142cd565b16613416816000526006602052604060002054151590565b156116f65760206001600160a01b0360045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa9081156116ea576000916135cd575b501561166c57613481816142cd565b61349d60a48401916132c46134968488614268565b3691613f06565b1561358657507ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0608067ffffffffffffffff61352a61352460446135148861350e61350961349660209d60c46134f28e6142cd565b9561350260648201358098614b92565b0190614268565b614749565b90614804565b97019561096588612b59896142b9565b946142b9565b936001600160a01b0360405195817f000000000000000000000000000000000000000000000000000000000000000016875233898801521660408601528560608601521692a28060405161357d81613e73565b52604051908152f35b6135909084614268565b6110cb6040519283927f24eb47e500000000000000000000000000000000000000000000000000000000845260206004850152602484019161454f565b6135e6915060203d602011611775576117678183613ec7565b84613472565b613605915060203d602011611775576117678183613ec7565b846133e7565b346102475761361936613dbb565b50509293909161362761424f565b506040519461363586613e57565b60008652600060208701526060608060405161365081613e3b565b828152826020820152826040820152600083820152015267ffffffffffffffff8316936040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116ea576000916139f7575b506117235761370b33615a64565b613722856000526006602052604060002054151590565b156139c95784600052600a6020526040600020948554831161399757509263ffffffff9261384c869361383e8a9760ff60016138d59c9b01549161ffff8360081c169983602061271061378a61ffff8f98816137839160181c16809a6146fb565b168a61465c565b049d019c8d5260281c168061392c575061ffff6137fc6137ac60038c01614359565b985b604051976137bb89613e3b565b8852602088019c8d52604088019586526060880193857f00000000000000000000000000000000000000000000000000000000000000001685523691613f06565b9360808701948552816040519c8d986020808b01525160408a01525116606088015251166080860152511660a08401525160a060c084015260e0830190613d12565b03601f198101865285613ec7565b6001600160a01b036020968795604051906138678883613ec7565b6000825261387d600260405198611c6e8a613e3b565b875287870152604086015216606084015260808301526001600160a01b0360045416906040518097819482937f20487ded0000000000000000000000000000000000000000000000000000000084526004840161445f565b03915afa9283156116ea576000936138fa575b50826040945283519283525190820152f35b9392508184813d8311613925575b6139128183613ec7565b81010312610247576040935192936138e8565b503d613908565b6137fc61ffff916040519061394082613e57565b81526020810160018152604051917f181dcf10000000000000000000000000000000000000000000000000000000006020840152516024830152511515604482015260448152613991606482613ec7565b986137ae565b90507f58dd87c50000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b847fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b613a10915060203d602011611775576117678183613ec7565b886136fd565b346102475760006003193601126102475760206103546141b0565b3461024757600060031936011261024757602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b34610247576020600319360112610247576020613a8a613d37565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b3461024757600060031936011261024757610d75604051613b28606082613ec7565b603581527f4275726e4d696e745769746845787465726e616c4d696e74657246617374547260208201527f616e73666572546f6b656e506f6f6c20312e362e3000000000000000000000006040820152604051918291602083526020830190613d12565b3461024757602060031936011261024757600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361024757817f61f8dc160000000000000000000000000000000000000000000000000000000060209314908115613c64575b8115613c07575b5015158152f35b7f85572ffb00000000000000000000000000000000000000000000000000000000811491508115613c3a575b5083613c00565b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483613c33565b90507faff2afbf0000000000000000000000000000000000000000000000000000000081148015613cc6575b8015613c9d575b90613bf9565b507f01ffc9a7000000000000000000000000000000000000000000000000000000008114613c97565b507f0e64dd29000000000000000000000000000000000000000000000000000000008114613c90565b60005b838110613d025750506000910152565b8181015183820152602001613cf2565b90601f19601f602093613d3081518092818752878088019101613cef565b0116010190565b600435906001600160a01b038216820361024757565b35906001600160a01b038216820361024757565b6004359067ffffffffffffffff8216820361024757565b359067ffffffffffffffff8216820361024757565b9181601f840112156102475782359167ffffffffffffffff8311610247576020838186019501011161024757565b9060a06003198301126102475760043567ffffffffffffffff8116810361024757916024359160443567ffffffffffffffff81116102475782613e0091600401613d8d565b929092916064356001600160a01b038116810361024757916084359067ffffffffffffffff821161024757613e3791600401613d8d565b9091565b60a0810190811067ffffffffffffffff8211176108a257604052565b6040810190811067ffffffffffffffff8211176108a257604052565b6020810190811067ffffffffffffffff8211176108a257604052565b60e0810190811067ffffffffffffffff8211176108a257604052565b6060810190811067ffffffffffffffff8211176108a257604052565b90601f601f19910116810190811067ffffffffffffffff8211176108a257604052565b67ffffffffffffffff81116108a257601f01601f191660200190565b929192613f1282613eea565b91613f206040519384613ec7565b829481845281830111610247578281602093846000960137010152565b9080601f8301121561024757816020613f5893359101613f06565b90565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501948460051b01011161024757565b60406003198201126102475760043567ffffffffffffffff81116102475781613fb791600401613f5b565b929092916024359067ffffffffffffffff821161024757613e3791600401613f5b565b9060406003198301126102475760043567ffffffffffffffff8116810361024757916024359067ffffffffffffffff821161024757613e3791600401613d8d565b906020808351928381520192019060005b8181106140395750505090565b82516001600160a01b031684526020938401939092019160010161402c565b6020600319820112610247576004359067ffffffffffffffff8211610247576003198260a0920301126102475760040190565b67ffffffffffffffff81116108a25760051b60200190565b9291906140af8161408b565b936140bd6040519586613ec7565b602085838152019160051b810192831161024757905b8282106140df57505050565b602080916140ec84613d4d565b8152019101906140d3565b9080601f8301121561024757816020613f58933591016140a3565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501946060850201011161024757565b9060038210156127185752565b35906001600160801b038216820361024757565b91908260609103126102475760405161417c81613eab565b809280359081151582036102475760406141ab91819385526141a060208201614150565b602086015201614150565b910152565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116ea57600091614220575090565b90506020813d602011614247575b8161423b60209383613ec7565b81010312610247575190565b3d915061422e565b6040519061425c82613e57565b60006020838281520152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610247570180359067ffffffffffffffff82116102475760200191813603831361024757565b356001600160a01b03811681036102475790565b3567ffffffffffffffff811681036102475790565b9067ffffffffffffffff613f5892166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b90600182811c9216801561434f575b602083101461433957565b634e487b7160e01b600052602260045260246000fd5b91607f169161432e565b906040519182600082549261436d8461431f565b80845293600181169081156143db5750600114614394575b5061439292500383613ec7565b565b90506000929192526020600020906000915b8183106143bf5750509060206143929282010138614385565b60209193508060019154838589010152019101909184926143a6565b6020935061439295925060ff1991501682840152151560051b82010138614385565b60038210156127185752565b91908110156144495760051b810135907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee181360301821215610247570190565b634e487b7160e01b600052603260045260246000fd5b9067ffffffffffffffff90939293168152604060208201526144a6614490845160a0604085015260e0840190613d12565b6020850151603f19848303016060850152613d12565b90604084015191603f198282030160808301526020808451928381520193019060005b818110614501575050506080846001600160a01b036060613f58969701511660a084015201519060c0603f1982850301910152613d12565b825180516001600160a01b0316865260209081015181870152604090950194909201916001016144c9565b9190820391821161453957565b634e487b7160e01b600052601160045260246000fd5b601f8260209493601f19938186528686013760008582860101520116010190565b80518210156144495760209160051b010190565b91908110156144495760051b0190565b9190811015614449576060020190565b604051906145b182613e3b565b60006080838281528260208201528260408201528260608201520152565b906040516145dc81613e3b565b60806001829460ff81546001600160801b038116865263ffffffff81861c16602087015260a01c161515604085015201546001600160801b0381166060840152811c910152565b67ffffffffffffffff166000526007602052613f586004604060002001614359565b818110614650575050565b60008155600101614645565b8181029291811591840414171561453957565b9190601f811161467e57505050565b614392926000526020600020906020601f840160051c830193106146aa575b601f0160051c0190614645565b909150819061469d565b92906146e76146f59260ff60405195869460208601988952604086015216606084015260808084015260a0830190613d12565b03601f198101835282613ec7565b51902090565b9061ffff8091169116019061ffff821161453957565b811561471b570490565b634e487b7160e01b600052601260045260246000fd5b90816020910312610247575180151581036102475790565b805180156147b95760200361477b57805160208281019183018390031261024757519060ff821161477b575060ff1690565b6110cb906040519182917f953576f7000000000000000000000000000000000000000000000000000000008352602060048401526024830190613d12565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff821161453957565b60ff16604d811161453957600a0a90565b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff8116928284146148ec578284116148c25790614849916147df565b91604d60ff84161180156148a7575b6148715750509061486b613f58926147f3565b9061465c565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b506148b1836147f3565b801561471b57600019048411614858565b6148cb916147df565b91604d60ff841611614871575050906148e6613f58926147f3565b90614711565b5050505090565b6040517f40c10f190000000000000000000000000000000000000000000000000000000081526001600160a01b03909116600482015260248101919091526020818060448101038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af180156116ea576149765750565b61498e9060203d602011611775576117678183613ec7565b50565b6001600160a01b036001541633036149a557565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b908051156108f15767ffffffffffffffff81516020830120921691826000526007602052614a04816005604060002001615745565b15614b3f5760005260086020526040600020815167ffffffffffffffff81116108a257614a3b81614a35845461431f565b8461466f565b6020601f8211600114614ab55791614a8f827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea9593614aa595600091614aaa575b506000198260011b9260031b1c19161790565b9055604051918291602083526020830190613d12565b0390a2565b905084015138614a7c565b601f1982169083600052806000209160005b818110614b27575092614aa59492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea989610614b0e575b5050811b019055610f53565b85015160001960f88460031b161c191690553880614b02565b9192602060018192868a015181550194019201614ac7565b50906110cb6040519283927f393b8ad20000000000000000000000000000000000000000000000000000000084526004840152604060248401526044830190613d12565b519061ffff8216820361024757565b67ffffffffffffffff7f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c91169182600052600760205280614c0660026040600020016001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615313565b604080516001600160a01b03909216825260208201929092529081908101614aa5565b357fffffffff00000000000000000000000000000000000000000000000000000000811681036102475790565b3563ffffffff811681036102475790565b3561ffff811681036102475790565b3580151581036102475790565b67ffffffffffffffff7fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da817894491169182600052600760205280614c0660406000206001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615313565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000060208201526001600160a01b039283166024820152929091166044830152606482019290925261439291614d5982608481015b03601f198101845283613ec7565b615507565b91909181158015614e46575b15614dc2576040517f095ea7b30000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909316602484015260448301919091526143929190614d598260648101614d4b565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b0384166024820152602081806044810103816001600160a01b0386165afa9081156116ea57600091614eae575b5015614d6a565b90506020813d602011614ed8575b81614ec960209383613ec7565b81010312610247575138614ea7565b3d9150614ebc565b67ffffffffffffffff166000818152600660205260409020549092919015614fd05791614fcd60e092614fa285614f377f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b9761515f565b846000526007602052614f4e816040600020615855565b614f578361515f565b846000526007602052614f71836002604060002001615855565b60405194855260208501906001600160801b0360408092805115158552826020820151166020860152015116910152565b60808301906001600160801b0360408092805115158552826020820151166020860152015116910152565ba1565b827f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b602060009160246001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169161505b81847f0000000000000000000000000000000000000000000000000000000000000000614d5e565b60405194859384927f42966c6800000000000000000000000000000000000000000000000000000000845260048401525af180156116ea576149765750565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000060208201526001600160a01b039092166024830152604482019290925261439291614d598260648101614d4b565b6150f46145a4565b506001600160801b036060820151166001600160801b03808351169161513f602085019361513961512c63ffffffff8751164261452c565b856080890151169061465c565b90615306565b8082101561515857505b16825263ffffffff4216905290565b9050615149565b8051156151e4576001600160801b036040820151166001600160801b036020830151161061518a5750565b6064906151e2604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565bfd5b6001600160801b036040820151161580159061525a575b6152025750565b6064906151e2604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565b506001600160801b0360208201511615156151fb565b60405190600b548083528260208101600b60005260206000209260005b8181106152a257505061439292500383613ec7565b845483526001948501948794506020909301920161528d565b906040519182815491828252602082019060005260206000209260005b8181106152ed57505061439292500383613ec7565b84548352600194850194879450602090930192016152d8565b9190820180921161453957565b9182549060ff8260a01c161580156154ff575b6154f9576001600160801b038216916001850190815461535963ffffffff6001600160801b0383169360801c164261452c565b908161545b575b505084811061541c57508383106153b15750506153866001600160801b0392839261452c565b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416179055565b5460801c916153c0818561452c565b92600019810190808211614539576153e36153e8926001600160a01b0396615306565b614711565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b82856001600160a01b03927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b8286929396116154cf57615476926151399160801c9061465c565b808410156154ca5750825b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff0000000000000000000000000000000016178655923880615360565b615481565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b508215615326565b6001600160a01b036155899116916040926000808551936155288786613ec7565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d15615632573d9161556d83613eea565b9261557a87519485613ec7565b83523d6000602085013e615ca1565b8051908161559657505050565b6020806155a7938301019101614731565b156155af5750565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606091615ca1565b80548210156144495760005260206000200190600090565b805490680100000000000000008210156108a257816156799160016156909401815561563a565b81939154906000199060031b92831b921b19161790565b9055565b806000526003602052604060002054156000146156cd576156b6816002615652565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146156cd576156f581600b615652565b600b5490600052600c602052604060002055600190565b806000526006602052604060002054156000146156cd5761572e816005615652565b600554906000526006602052604060002055600190565b600082815260018201602052604090205461577c578061576783600193615652565b80549260005201602052604060002055600190565b5050600090565b805480156157ab57600019019061579a828261563a565b60001982549160031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c6020526040902054801561577c57600019810181811161453957600b549060001982019182116145395780820361581b575b505050615807600b615783565b600052600c60205260006040812055600190565b61583d61582c61567993600b61563a565b90549060031b1c928392600b61563a565b9055600052600c6020526040600020553880806157fa565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199161597c606092805461589263ffffffff8260801c164261452c565b90816159b2575b50506001600160801b0360018160208601511692828154168085106000146159aa57508280855b16167fffffffffffffffffffffffffffffffff000000000000000000000000000000008254161781556159398651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601517fffffffffffffffffffffffffffffffff0000000000000000000000000000000060809190911b16939092166001600160801b031692909217910155565b614fcd60405180926001600160801b0360408092805115158552826020820151166020860152015116910152565b8380916158c0565b6001600160801b03916159de8392836159d76001880154948286169560801c9061465c565b9116615306565b80821015615a5d57505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9290911692909216167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116174260801b73ffffffff00000000000000000000000000000000161781553880615899565b90506159e8565b7f0000000000000000000000000000000000000000000000000000000000000000615a8c5750565b6001600160a01b031680600052600360205260406000205415615aac5750565b7fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b600081815260036020526040902054801561577c5760001981018181116145395760025490600019820191821161453957818103615b33575b505050615b1f6002615783565b600052600360205260006040812055600190565b615b55615b4461567993600261563a565b90549060031b1c928392600261563a565b90556000526003602052604060002055388080615b12565b600081815260066020526040902054801561577c5760001981018181116145395760055490600019820191821161453957818103615bc7575b505050615bb36005615783565b600052600660205260006040812055600190565b615be9615bd861567993600561563a565b90549060031b1c928392600561563a565b90556000526006602052604060002055388080615ba6565b906001820191816000528260205260406000205490811515600014615c98576000198201918083116145395781546000198101908111614539578381615c4f9503615c61575b505050615783565b60005260205260006040812055600190565b615c81615c71615679938661563a565b90549060031b1c9283928661563a565b905560005284602052604060002055388080615c47565b50505050600090565b91929015615d1c5750815115615cb5575090565b3b15615cbe5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015615d2f5750805190602001fd5b6110cb906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190613d1256fea164736f6c634300081a000a",
}

var BurnMintWithExternalMinterFastTransferTokenPoolABI = BurnMintWithExternalMinterFastTransferTokenPoolMetaData.ABI

var BurnMintWithExternalMinterFastTransferTokenPoolBin = BurnMintWithExternalMinterFastTransferTokenPoolMetaData.Bin

func DeployBurnMintWithExternalMinterFastTransferTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, minter common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnMintWithExternalMinterFastTransferTokenPool, error) {
	parsed, err := BurnMintWithExternalMinterFastTransferTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintWithExternalMinterFastTransferTokenPoolBin), backend, minter, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintWithExternalMinterFastTransferTokenPool{address: address, abi: *parsed, BurnMintWithExternalMinterFastTransferTokenPoolCaller: BurnMintWithExternalMinterFastTransferTokenPoolCaller{contract: contract}, BurnMintWithExternalMinterFastTransferTokenPoolTransactor: BurnMintWithExternalMinterFastTransferTokenPoolTransactor{contract: contract}, BurnMintWithExternalMinterFastTransferTokenPoolFilterer: BurnMintWithExternalMinterFastTransferTokenPoolFilterer{contract: contract}}, nil
}

type BurnMintWithExternalMinterFastTransferTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnMintWithExternalMinterFastTransferTokenPoolCaller
	BurnMintWithExternalMinterFastTransferTokenPoolTransactor
	BurnMintWithExternalMinterFastTransferTokenPoolFilterer
}

type BurnMintWithExternalMinterFastTransferTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnMintWithExternalMinterFastTransferTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnMintWithExternalMinterFastTransferTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnMintWithExternalMinterFastTransferTokenPoolSession struct {
	Contract     *BurnMintWithExternalMinterFastTransferTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintWithExternalMinterFastTransferTokenPoolCallerSession struct {
	Contract *BurnMintWithExternalMinterFastTransferTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession struct {
	Contract     *BurnMintWithExternalMinterFastTransferTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnMintWithExternalMinterFastTransferTokenPoolRaw struct {
	Contract *BurnMintWithExternalMinterFastTransferTokenPool
}

type BurnMintWithExternalMinterFastTransferTokenPoolCallerRaw struct {
	Contract *BurnMintWithExternalMinterFastTransferTokenPoolCaller
}

type BurnMintWithExternalMinterFastTransferTokenPoolTransactorRaw struct {
	Contract *BurnMintWithExternalMinterFastTransferTokenPoolTransactor
}

func NewBurnMintWithExternalMinterFastTransferTokenPool(address common.Address, backend bind.ContractBackend) (*BurnMintWithExternalMinterFastTransferTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintWithExternalMinterFastTransferTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintWithExternalMinterFastTransferTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPool{address: address, abi: abi, BurnMintWithExternalMinterFastTransferTokenPoolCaller: BurnMintWithExternalMinterFastTransferTokenPoolCaller{contract: contract}, BurnMintWithExternalMinterFastTransferTokenPoolTransactor: BurnMintWithExternalMinterFastTransferTokenPoolTransactor{contract: contract}, BurnMintWithExternalMinterFastTransferTokenPoolFilterer: BurnMintWithExternalMinterFastTransferTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnMintWithExternalMinterFastTransferTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnMintWithExternalMinterFastTransferTokenPoolCaller, error) {
	contract, err := bindBurnMintWithExternalMinterFastTransferTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolCaller{contract: contract}, nil
}

func NewBurnMintWithExternalMinterFastTransferTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintWithExternalMinterFastTransferTokenPoolTransactor, error) {
	contract, err := bindBurnMintWithExternalMinterFastTransferTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolTransactor{contract: contract}, nil
}

func NewBurnMintWithExternalMinterFastTransferTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintWithExternalMinterFastTransferTokenPoolFilterer, error) {
	contract, err := bindBurnMintWithExternalMinterFastTransferTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolFilterer{contract: contract}, nil
}

func bindBurnMintWithExternalMinterFastTransferTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintWithExternalMinterFastTransferTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.BurnMintWithExternalMinterFastTransferTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.BurnMintWithExternalMinterFastTransferTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.BurnMintWithExternalMinterFastTransferTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) ComputeFillId(opts *bind.CallOpts, settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "computeFillId", settlementId, sourceAmountNetFee, sourceDecimals, receiver)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) ComputeFillId(settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ComputeFillId(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, settlementId, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) ComputeFillId(settlementId [32]byte, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver []byte) ([32]byte, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ComputeFillId(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, settlementId, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAccumulatedPoolFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetAccumulatedPoolFees() (*big.Int, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAccumulatedPoolFees(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetAccumulatedPoolFees() (*big.Int, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAccumulatedPoolFees(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAllowList(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAllowList(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAllowListEnabled(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAllowListEnabled(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getAllowedFillers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetAllowedFillers() ([]common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAllowedFillers(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetAllowedFillers() ([]common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetAllowedFillers(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getCcipSendTokenFee", destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)

	if err != nil {
		return *new(IFastTransferPoolQuote), err
	}

	out0 := *abi.ConvertType(out[0], new(IFastTransferPoolQuote)).(*IFastTransferPoolQuote)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetCcipSendTokenFee(destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetCcipSendTokenFee(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetCcipSendTokenFee(destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (IFastTransferPoolQuote, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetCcipSendTokenFee(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getDestChainConfig", remoteChainSelector)

	if err != nil {
		return *new(FastTransferTokenPoolAbstractDestChainConfig), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(FastTransferTokenPoolAbstractDestChainConfig)).(*FastTransferTokenPoolAbstractDestChainConfig)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetDestChainConfig(remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetDestChainConfig(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetDestChainConfig(remoteChainSelector uint64) (FastTransferTokenPoolAbstractDestChainConfig, []common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetDestChainConfig(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetFillInfo(opts *bind.CallOpts, fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getFillInfo", fillId)

	if err != nil {
		return *new(FastTransferTokenPoolAbstractFillInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FastTransferTokenPoolAbstractFillInfo)).(*FastTransferTokenPoolAbstractFillInfo)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetFillInfo(fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetFillInfo(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, fillId)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetFillInfo(fillId [32]byte) (FastTransferTokenPoolAbstractFillInfo, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetFillInfo(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, fillId)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetMinter() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetMinter(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetMinter() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetMinter(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRateLimitAdmin(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRateLimitAdmin(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRemotePools(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRemotePools(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRemoteToken(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRemoteToken(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRmnProxy(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRmnProxy(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRouter(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetRouter(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetSupportedChains(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetSupportedChains(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetToken(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetToken(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetTokenDecimals(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.GetTokenDecimals(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isAllowedFiller", filler)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) IsAllowedFiller(filler common.Address) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsAllowedFiller(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, filler)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) IsAllowedFiller(filler common.Address) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsAllowedFiller(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, filler)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsRemotePool(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsRemotePool(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsSupportedChain(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsSupportedChain(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsSupportedToken(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, token)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.IsSupportedToken(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, token)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) Owner() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.Owner(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.Owner(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SupportsInterface(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, interfaceId)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SupportsInterface(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts, interfaceId)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.TypeAndVersion(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.TypeAndVersion(&_BurnMintWithExternalMinterFastTransferTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.AcceptOwnership(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.AcceptOwnership(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.AddRemotePool(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.AddRemotePool(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ApplyChainUpdates(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ApplyChainUpdates(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipReceive", message)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipReceive(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, message)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipReceive(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, message)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipSendToken", destinationChainSelector, amount, receiver, feeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, receiver, feeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, receiver, feeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) FastFill(opts *bind.TransactOpts, fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "fastFill", fillId, settlementId, sourceChainSelector, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) FastFill(fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.FastFill(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, fillId, settlementId, sourceChainSelector, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) FastFill(fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.FastFill(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, fillId, settlementId, sourceChainSelector, sourceAmountNetFee, sourceDecimals, receiver)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.LockOrBurn(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.LockOrBurn(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ReleaseOrMint(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.ReleaseOrMint(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.RemoveRemotePool(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.RemoveRemotePool(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setChainRateLimiterConfigs", remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfigs(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetChainRateLimiterConfigs(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetRateLimitAdmin(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetRateLimitAdmin(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetRouter(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, newRouter)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.SetRouter(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, newRouter)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.TransferOwnership(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, to)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.TransferOwnership(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, to)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) UpdateDestChainConfig(opts *bind.TransactOpts, destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "updateDestChainConfig", destChainConfigArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) UpdateDestChainConfig(destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.UpdateDestChainConfig(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destChainConfigArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) UpdateDestChainConfig(destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.UpdateDestChainConfig(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destChainConfigArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "updateFillerAllowList", fillersToAdd, fillersToRemove)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) UpdateFillerAllowList(fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.UpdateFillerAllowList(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, fillersToAdd, fillersToRemove)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) UpdateFillerAllowList(fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.UpdateFillerAllowList(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, fillersToAdd, fillersToRemove)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "withdrawPoolFees", recipient)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) WithdrawPoolFees(recipient common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.WithdrawPoolFees(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, recipient)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) WithdrawPoolFees(recipient common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.WithdrawPoolFees(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, recipient)
}

type BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolChainAdded)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolChainAdded)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolChainAdded)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolChainAdded, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolChainAdded)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated struct {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterDestChainConfigUpdated(opts *bind.FilterOpts, destinationChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator, error) {

	var destinationChainSelectorRule []interface{}
	for _, destinationChainSelectorItem := range destinationChainSelector {
		destinationChainSelectorRule = append(destinationChainSelectorRule, destinationChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "DestChainConfigUpdated", destinationChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "DestChainConfigUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, destinationChainSelector []uint64) (event.Subscription, error) {

	var destinationChainSelectorRule []interface{}
	for _, destinationChainSelectorItem := range destinationChainSelector {
		destinationChainSelectorRule = append(destinationChainSelectorRule, destinationChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "DestChainConfigUpdated", destinationChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseDestChainConfigUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated struct {
	DestChainSelector uint64
	DestinationPool   common.Address
	Raw               types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterDestinationPoolUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "DestinationPoolUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "DestinationPoolUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchDestinationPoolUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "DestinationPoolUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestinationPoolUpdated", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseDestinationPoolUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "DestinationPoolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled struct {
	FillId       [32]byte
	SettlementId [32]byte
	Filler       common.Address
	DestAmount   *big.Int
	Receiver     common.Address
	Raw          types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterFastTransferFilled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator, error) {

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

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FastTransferFilled", fillIdRule, settlementIdRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "FastTransferFilled", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchFastTransferFilled(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FastTransferFilled", fillIdRule, settlementIdRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferFilled", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseFastTransferFilled(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferFilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested struct {
	DestinationChainSelector uint64
	FillId                   [32]byte
	SettlementId             [32]byte
	SourceAmountNetFee       *big.Int
	SourceDecimals           uint8
	FastTransferFee          *big.Int
	Receiver                 []byte
	Raw                      types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator, error) {

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

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FastTransferRequested", destinationChainSelectorRule, fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "FastTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchFastTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FastTransferRequested", destinationChainSelectorRule, fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferRequested", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseFastTransferRequested(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled struct {
	FillId                    [32]byte
	SettlementId              [32]byte
	FillerReimbursementAmount *big.Int
	PoolFeeAccumulated        *big.Int
	PrevState                 uint8
	Raw                       types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterFastTransferSettled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator, error) {

	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FastTransferSettled", fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "FastTransferSettled", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchFastTransferSettled(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error) {

	var fillIdRule []interface{}
	for _, fillIdItem := range fillId {
		fillIdRule = append(fillIdRule, fillIdItem)
	}
	var settlementIdRule []interface{}
	for _, settlementIdItem := range settlementId {
		settlementIdRule = append(settlementIdRule, settlementIdItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FastTransferSettled", fillIdRule, settlementIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferSettled", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseFastTransferSettled(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FastTransferSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated struct {
	AddFillers    []common.Address
	RemoveFillers []common.Address
	Raw           types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterFillerAllowListUpdated(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "FillerAllowListUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "FillerAllowListUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchFillerAllowListUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "FillerAllowListUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FillerAllowListUpdated", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseFillerAllowListUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "FillerAllowListUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "InboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseInboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "LockedOrBurned", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseLockedOrBurned(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "OutboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseOutboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterPoolFeeWithdrawn(opts *bind.FilterOpts, recipient []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "PoolFeeWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "PoolFeeWithdrawn", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchPoolFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "PoolFeeWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "PoolFeeWithdrawn", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParsePoolFeeWithdrawn(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "PoolFeeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Recipient           common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "ReleasedOrMinted", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseReleasedOrMinted(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator struct {
	Event *BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated)
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
		it.Event = new(BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated)
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

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator{contract: _BurnMintWithExternalMinterFastTransferTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated)
				if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated, error) {
	event := new(BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated)
	if err := _BurnMintWithExternalMinterFastTransferTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseAllowListAdd(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseAllowListRemove(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseChainAdded(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseChainConfigured(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseChainRemoved(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseConfigChanged(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["DestChainConfigUpdated"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseDestChainConfigUpdated(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["DestinationPoolUpdated"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseDestinationPoolUpdated(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["FastTransferFilled"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseFastTransferFilled(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["FastTransferRequested"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseFastTransferRequested(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["FastTransferSettled"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseFastTransferSettled(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["FillerAllowListUpdated"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseFillerAllowListUpdated(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["InboundRateLimitConsumed"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseInboundRateLimitConsumed(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["LockedOrBurned"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseLockedOrBurned(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["OutboundRateLimitConsumed"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseOutboundRateLimitConsumed(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseOwnershipTransferred(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["PoolFeeWithdrawn"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParsePoolFeeWithdrawn(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseRateLimitAdminSet(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["ReleasedOrMinted"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseReleasedOrMinted(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseRemotePoolAdded(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseRemotePoolRemoved(log)
	case _BurnMintWithExternalMinterFastTransferTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnMintWithExternalMinterFastTransferTokenPool.ParseRouterUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated) Topic() common.Hash {
	return common.HexToHash("0xb760e03fa04c0e86fcff6d0046cdcf22fb5d5b6a17d1e6f890b3456e81c40fd8")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled) Topic() common.Hash {
	return common.HexToHash("0xd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf1")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x240a1286fd41f1034c4032dcd6b93fc09e81be4a0b64c7ecee6260b605a8e016")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled) Topic() common.Hash {
	return common.HexToHash("0x33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated) Topic() common.Hash {
	return common.HexToHash("0xfd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0x50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned) Topic() common.Hash {
	return common.HexToHash("0xf33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae10")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0xff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da8178944")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted) Topic() common.Hash {
	return common.HexToHash("0xfc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPool) Address() common.Address {
	return _BurnMintWithExternalMinterFastTransferTokenPool.address
}

type BurnMintWithExternalMinterFastTransferTokenPoolInterface interface {
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

	GetMinter(opts *bind.CallOpts) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

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

	CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error)

	FastFill(opts *bind.TransactOpts, fillId [32]byte, settlementId [32]byte, sourceChainSelector uint64, sourceAmountNetFee *big.Int, sourceDecimals uint8, receiver common.Address) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDestChainConfig(opts *bind.TransactOpts, destChainConfigArgs []FastTransferTokenPoolAbstractDestChainConfigUpdateArgs) (*types.Transaction, error)

	UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error)

	WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolAllowListRemove, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolConfigChanged, error)

	FilterDestChainConfigUpdated(opts *bind.FilterOpts, destinationChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdatedIterator, error)

	WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, destinationChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolDestChainConfigUpdated, error)

	FilterDestinationPoolUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdatedIterator, error)

	WatchDestinationPoolUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParseDestinationPoolUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolDestinationPoolUpdated, error)

	FilterFastTransferFilled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilledIterator, error)

	WatchFastTransferFilled(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled, fillId [][32]byte, settlementId [][32]byte, filler []common.Address) (event.Subscription, error)

	ParseFastTransferFilled(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferFilled, error)

	FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator, error)

	WatchFastTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested, destinationChainSelector []uint64, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error)

	ParseFastTransferRequested(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequested, error)

	FilterFastTransferSettled(opts *bind.FilterOpts, fillId [][32]byte, settlementId [][32]byte) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettledIterator, error)

	WatchFastTransferSettled(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled, fillId [][32]byte, settlementId [][32]byte) (event.Subscription, error)

	ParseFastTransferSettled(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFastTransferSettled, error)

	FilterFillerAllowListUpdated(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdatedIterator, error)

	WatchFillerAllowListUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated) (event.Subscription, error)

	ParseFillerAllowListUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolFillerAllowListUpdated, error)

	FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumedIterator, error)

	WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseInboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolInboundRateLimitConsumed, error)

	FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurnedIterator, error)

	WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error)

	ParseLockedOrBurned(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolLockedOrBurned, error)

	FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumedIterator, error)

	WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseOutboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolOutboundRateLimitConsumed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolOwnershipTransferred, error)

	FilterPoolFeeWithdrawn(opts *bind.FilterOpts, recipient []common.Address) (*BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawnIterator, error)

	WatchPoolFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, recipient []common.Address) (event.Subscription, error)

	ParsePoolFeeWithdrawn(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolPoolFeeWithdrawn, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRateLimitAdminSet, error)

	FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMintedIterator, error)

	WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error)

	ParseReleasedOrMinted(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolReleasedOrMinted, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnMintWithExternalMinterFastTransferTokenPoolRouterUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
