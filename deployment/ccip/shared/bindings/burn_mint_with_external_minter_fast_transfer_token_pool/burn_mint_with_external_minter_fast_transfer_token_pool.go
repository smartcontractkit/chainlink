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

var BurnMintWithExternalMinterFastTransferTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSendToken\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"computeFillId\",\"inputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"fastFill\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAccumulatedPoolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedFillers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCcipSendTokenFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structIFastTransferPool.Quote\",\"components\":[{\"name\":\"ccipSettlementFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfig\",\"components\":[{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFillInfo\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.FillInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIFastTransferPool.FillState\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAllowedFiller\",\"inputs\":[{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateDestChainConfig\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfigUpdateArgs[]\",\"components\":[{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFillerAllowList\",\"inputs\":[{\"name\":\"fillersToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"fillersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPoolFees\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestinationPoolUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destinationPool\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"filler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferRequested\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"fillerFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fillerReimbursementAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFeeAccumulated\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"prevState\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIFastTransferPool.FillState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FillerAllowListUpdated\",\"inputs\":[{\"name\":\"addFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"removeFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolFeeWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyFilledOrSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadySettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"FillerNotAllowlisted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientPoolFees\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFillId\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"QuoteFeeExceedsUserMaxLimit\",\"inputs\":[{\"name\":\"quoteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferAmountExceedsMaxFillAmount\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461032857616684803803809161001d8285610559565b8339810160c082820312610328576100348261057c565b60208301516001600160a01b03811693919290918483036103285761005b60408201610590565b60608201519091906001600160401b0381116103285781019280601f85011215610328578351936001600160401b038511610543578460051b9060208201956100a76040519788610559565b865260208087019282010192831161032857602001905b82821061052b575050506100e060a06100d96080840161057c565b920161057c565b93331561051a57600180546001600160a01b0319163317905586158015610509575b80156104f8575b6104e75760805260c05260405163313ce56760e01b8152602081600481895afa600091816104ab575b50610480575b5060a052600480546001600160a01b0319166001600160a01b0384169081179091558151151560e0819052909190610357575b501561034157610120526001600160a01b03166101008190526040516321df0da760e01b815290602090829060049082905afa908115610335576000916102f6575b506001600160a01b0316908181036102df57604051615f45908161073f823960805181818161130a01528181611374015281816114e0015281816121120152818161274f01528181612e2601528181613081015281816136e20152818161372f01528181613b4c0152818161468d01528181614b0201526157a4015260a0518181816115f6015281816133ac01528181613698015281816139fb01528181613cc701528181613cf901528181614d940152614dfe015260c051818181610b0a015281816113e90152818161244c01528181612e9c015281816132af01526138f5015260e051818181610ac501528181612c280152615c650152610100518181816102230152818161150c015281816125dc0152818161270d015281816130010152613b8601526101205181613daf0152f35b63f902523f60e01b60005260045260245260446000fd5b90506020813d60201161032d575b8161031160209383610559565b81010312610328576103229061057c565b386101ad565b600080fd5b3d9150610304565b6040513d6000823e3d90fd5b6335fdcccd60e21b600052600060045260246000fd5b9192906020926040519261036b8585610559565b60008452600036813760e0511561046f5760005b84518110156103e6576001906001600160a01b0361039d828861059e565b5116876103a9826105e0565b6103b6575b50500161037f565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a138876103ae565b50919490925060005b8351811015610463576001906001600160a01b0361040d828761059e565b5116801561045d578661041f826106de565b61042d575b50505b016103ef565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a13886610424565b50610427565b5092509290503861016b565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036104945750610138565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d6020116104df575b816104c760209383610559565b81010312610328576104d890610590565b9038610132565b3d91506104ba565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03821615610109565b506001600160a01b03851615610102565b639b15e16f60e01b60005260046000fd5b602080916105388461057c565b8152019101906100be565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b0382119082101761054357604052565b51906001600160a01b038216820361032857565b519060ff8216820361032857565b80518210156105b25760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156105b25760005260206000200190600090565b60008181526003602052604090205480156106d75760001981018181116106c1576002546000198101919082116106c157818103610670575b505050600254801561065a57600019016106348160026105c8565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b6106a96106816106929360026105c8565b90549060031b1c92839260026105c8565b819391549060031b91821b91600019901b19161790565b90556000526003602052604060002055388080610619565b634e487b7160e01b600052601160045260246000fd5b5050600090565b8060005260036020526040600020541560001461073857600254680100000000000000008110156105435761071f61069282600185940160025560026105c8565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a714613fb357508063055befd4146137d9578063181f5a771461375357806321df0da71461370f578063240028e8146136bc57806324f65ee71461367e5780632b2c0eb4146136635780632e7aa8c81461316f5780633907753714612d9b5780634c5ef0ed14612d5657806354c8a4f314612bf657806362ddd3c414612b8c5780636609f59914612b705780636d3d1a5814612b495780636def4ce714612a0557806378b410f2146129cb57806379ba5097146129255780637d54534e146128b057806385572ffb1461220957806387f060d014611f645780638926f54f14611f1f5780638a18dcbd14611a4a5780638da5cb5b14611a23578063929ea5ba14611919578063962d4020146117dd5780639a4575b9146113395780639fe280f5146112a6578063a42a7b8b14611174578063a7cd63b714611106578063abe1c1e814611097578063acfecf9114610f8b578063af58d59f14610f41578063b0f479a114610f1a578063b794658014610ee2578063c0d7865514610e49578063c4bffe2b14610d37578063c75eea9c14610c97578063cf7401f314610b2e578063dc0bd97114610aea578063e0351e1314610aad578063e8a1da1714610351578063eeebc674146102f9578063f2fde38b1461024c5763f36675171461020357600080fd5b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b34610247576020600319360112610247576001600160a01b0361026d614170565b610275614f02565b163381146102cf578073ffffffffffffffffffffffffffffffffffffffff1960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760806003193601126102475760443560ff811681036102475760643567ffffffffffffffff81116102475760209161033d610349923690600401614300565b90602435600435614a65565b604051908152f35b346102475761035f3661434f565b91909261036a614f02565b6000905b8282106109215750505060009063ffffffff42165b81831061038c57005b6103978383866148aa565b926101208436031261024757604051936103b0856141fe565b6103b98161412d565b8552602081013567ffffffffffffffff81116102475781019336601f860112156102475784356103e881614452565b956103f6604051978861428a565b81875260208088019260051b820101903682116102475760208101925b8284106108f2575050505060208601948552604082013567ffffffffffffffff8111610247576104469036908401614300565b906040870191825261047061045e366060860161452b565b936060890194855260c036910161452b565b94608088019586526104828451615322565b61048c8651615322565b825151156108c8576104a867ffffffffffffffff895116615934565b1561088f5767ffffffffffffffff885116600052600760205260406000206105b785516001600160801b036040820151169061058a6001600160801b03602083015116915115158360806040516104fe816141fe565b858152602081018a905260408101849052606081018690520152855474ff000000000000000000000000000000000000000091151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166001600160801b0384161773ffffffff00000000000000000000000000000000608089901b1617178555565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176001830155565b6106aa87516001600160801b036040820151169061067d6001600160801b03602083015116915115158360806040516105ef816141fe565b858152602081018a9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166001600160801b0385161773ffffffff0000000000000000000000000000000060808a901b161791151560a01b74ff000000000000000000000000000000000000000016919091179055565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176003830155565b6004845191019080519067ffffffffffffffff8211610879576106d7826106d185546147c0565b85614a20565b602090601f831160011461081257610707929160009183610807575b50506000198260011b9260031b1c19161790565b90555b60005b87518051821015610742579061073c6001926107358367ffffffffffffffff8e511692614900565b5190614f40565b0161070d565b505097967f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c29392919650946107fc67ffffffffffffffff60019751169251935191516107d16107a5604051968796875261010060208801526101008701906141bd565b9360408601906001600160801b0360408092805115158552826020820151166020860152015116910152565b60a08401906001600160801b0360408092805115158552826020820151166020860152015116910152565b0390a1019192610383565b015190508d806106f3565b90601f1983169184600052816000209260005b8181106108615750908460019594939210610848575b505050811b01905561070a565b015160001960f88460031b161c191690558c808061083b565b92936020600181928786015181550195019301610825565b634e487b7160e01b600052604160045260246000fd5b67ffffffffffffffff8851167f1d5ad3c50000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b833567ffffffffffffffff8111610247576020916109168392833691870101614300565b815201930192610413565b9092919367ffffffffffffffff61094161093c868886614914565b61476e565b169261094c84615d6c565b15610a985783600052600760205261096a6005604060002001615817565b9260005b84518110156109a65760019086600052600760205261099f60056040600020016109988389614900565b5190615e00565b500161096e565b50939094919592508060005260076020526005604060002060008155600060018201556000600282015560006003820155600481016109e581546147c0565b9081610a55575b5050018054906000815581610a34575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a101909161036e565b6000526020600020908101905b818110156109fc5760008155600101610a41565b81601f60009311600114610a6d5750555b88806109ec565b81835260208320610a8891601f01861c8101906001016149f6565b8082528160208120915555610a66565b83631e670e4b60e01b60005260045260246000fd5b346102475760006003193601126102475760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760e060031936011261024757610b47614116565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc36011261024757604051610b7d8161426e565b60243580151581036102475781526044356001600160801b03811681036102475760208201526064356001600160801b038116810361024757604082015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7c3601126102475760405190610bf28261426e565b608435801515810361024757825260a4356001600160801b038116810361024757602083015260c4356001600160801b03811681036102475760408301526001600160a01b036009541633141580610c82575b610c5457610c5292615158565b005b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b506001600160a01b0360015416331415610c45565b346102475760206003193601126102475767ffffffffffffffff610cb9614116565b610cc1614955565b50166000526007602052610d33610ce3610cde6040600020614980565b6152af565b6040519182918291909160806001600160801b038160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b34610247576000600319360112610247576040516005548082528160208101600560005260206000209260005b818110610e30575050610d799250038261428a565b805190610d9e610d8883614452565b92610d96604051948561428a565b808452614452565b90601f1960208401920136833760005b8151811015610de0578067ffffffffffffffff610dcd60019385614900565b5116610dd98287614900565b5201610dae565b5050906040519182916020830190602084525180915260408301919060005b818110610e0d575050500390f35b825167ffffffffffffffff16845285945060209384019390920191600101610dff565b8454835260019485019486945060209093019201610d64565b3461024757602060031936011261024757610e62614170565b610e6a614f02565b6001600160a01b0381169081156108c8576004805473ffffffffffffffffffffffffffffffffffffffff1981169093179055604080516001600160a01b0393841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491819081015b0390a1005b3461024757602060031936011261024757610d33610f06610f01614116565b6149d4565b6040519182916020835260208301906141bd565b346102475760006003193601126102475760206001600160a01b0360045416604051908152f35b346102475760206003193601126102475767ffffffffffffffff610f63614116565b610f6b614955565b50166000526007602052610d33610ce3610cde6002604060002001614980565b346102475767ffffffffffffffff610fa2366143a1565b929091610fad614f02565b1690610fc6826000526006602052604060002054151590565b1561108257816000526007602052610ff76005604060002001610fea3686856142c9565b6020815191012090615e00565b1561103b577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d769192611036604051928392602084526020840191614934565b0390a2005b61107e906040519384937f74f23c7c0000000000000000000000000000000000000000000000000000000085526004850152604060248501526044840191614934565b0390fd5b50631e670e4b60e01b60005260045260246000fd5b34610247576020600319360112610247576110b06146f0565b50600435600052600d6020526040806000206001600160a01b038251916110d68361421a565b546110e460ff82168461489e565b81602084019160081c1681526110fd845180945161450a565b51166020820152f35b34610247576000600319360112610247576040516002548082526020820190600260005260206000209060005b81811061115e57610d338561114a8187038261428a565b6040519182916020835260208301906143e2565b8254845260209093019260019283019201611133565b346102475760206003193601126102475767ffffffffffffffff611196614116565b1660005260076020526111af6005604060002001615817565b805190601f196111d76111c184614452565b936111cf604051958661428a565b808552614452565b0160005b81811061129557505060005b815181101561122f57806111fd60019284614900565b51600052600860205261121360406000206147fa565b61121d8286614900565b526112288185614900565b50016111e7565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b82821061126857505050500390f35b9193602061128582603f19600195979984950301865288516141bd565b9601920192018594939192611259565b8060606020809387010152016111db565b34610247576020600319360112610247576112bf614170565b6112c7614f02565b6112cf614651565b90816112d757005b60206001600160a01b038261132e857f738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599957f000000000000000000000000000000000000000000000000000000000000000061525d565b6040519485521692a2005b34610247576113473661441f565b606060206040516113578161421a565b8281520152608081016113698161475a565b6001600160a01b03807f00000000000000000000000000000000000000000000000000000000000000001691160361179e57506020810177ffffffffffffffff000000000000000000000000000000006113c28261476e565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561169e5760009161177f575b506117555761143361142e6040840161475a565b615c63565b67ffffffffffffffff6114458261476e565b1661145d816000526006602052604060002054151590565b156117285760206001600160a01b0360045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa90811561169e576000916116d8575b506001600160a01b031633036116aa576000916114de60606114d48461476e565b9201358092614ab9565b7f00000000000000000000000000000000000000000000000000000000000000009160206001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016611537848287614bb8565b6024604051809781937f42966c680000000000000000000000000000000000000000000000000000000083528760048401525af192831561169e577ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae1067ffffffffffffffff610f019461165c976115ec9761166f575b506115e46115ba8661476e565b604080516001600160a01b0390971687523360208801528601929092529116929081906060820190565b0390a261476e565b610d3360405160ff7f00000000000000000000000000000000000000000000000000000000000000001660208201526020815261162a60408261428a565b604051926116378461421a565b83526020830190815260405193849360208552516040602086015260608501906141bd565b9051601f198483030160408501526141bd565b6116909060203d602011611697575b611688818361428a565b810190614eea565b50886115ad565b503d61167e565b6040513d6000823e3d90fd5b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d602011611720575b816116f16020938361428a565b8101031261171c5751906001600160a01b038216820361171957506001600160a01b036114b3565b80fd5b5080fd5b3d91506116e4565b7fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b611798915060203d60201161169757611688818361428a565b8361141a565b6117af6001600160a01b039161475a565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b346102475760606003193601126102475760043567ffffffffffffffff81116102475761180e90369060040161431e565b9060243567ffffffffffffffff81116102475761182f9036906004016144d9565b9060443567ffffffffffffffff8111610247576118509036906004016144d9565b6001600160a01b036009541633141580611904575b610c54578386148015906118fa575b6118d05760005b86811061188457005b806118ca61189861093c6001948b8b614914565b6118a3838989614924565b6118c46118bc6118b486898b614924565b92369061452b565b91369061452b565b91615158565b0161187b565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b5080861415611874565b506001600160a01b0360015416331415611865565b346102475760406003193601126102475760043567ffffffffffffffff81116102475761194a9036906004016144be565b60243567ffffffffffffffff81116102475761196a9036906004016144be565b90611973614f02565b60005b81518110156119a5578061199e6001600160a01b0361199760019486614900565b51166158fb565b5001611976565b5060005b82518110156119d857806119d16001600160a01b036119ca60019487614900565b51166159e9565b50016119a9565b7ffd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193611a1583610edd866040519384936040855260408501906143e2565b9083820360208501526143e2565b346102475760006003193601126102475760206001600160a01b0360015416604051908152f35b346102475760206003193601126102475760043567ffffffffffffffff811161024757611a7b90369060040161431e565b611a83614f02565b60005b818110611a8f57005b611a9a8183856148aa565b60a081017f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000611ae9836150fe565b1614611ede575b60208201611afd8161513c565b90604084019161ffff80611b108561513c565b1691160161ffff8111611ec85761ffff61271091161015611e9e576080840167ffffffffffffffff611b418261476e565b16600052600a60205260406000209460e0810194611b5f8683614709565b600289019167ffffffffffffffff821161087957611b81826106d185546147c0565b600090601f8311600114611e3a57611bb0929160009183611e2f5750506000198260011b9260031b1c19161790565b90555b611bbc8461513c565b926001880197885498611bce8861513c565b60181b64ffff0000001695611be28661514b565b151560c087013597888555606088019c611bfb8e61512b565b60281b68ffffffff0000000000169360081b62ffff0016907fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016177fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff16179060ff16171790556101008401611c709085614709565b90916003019167ffffffffffffffff821161087957611c93826106d185546147c0565b600090601f8311600114611dba579180611cc892611ccf969594600092611daf5750506000198260011b9260031b1c19161790565b905561476e565b93611cd99061513c565b94611ce39061513c565b95611cee9083614709565b9091611cf9906150fe565b97611d039061512b565b92611d0d9061514b565b936040519761ffff899816885261ffff16602088015260408701526060860160e0905260e0860190611d3e92614934565b957fffffffff0000000000000000000000000000000000000000000000000000000016608085015263ffffffff1660a0840152151560c083015267ffffffffffffffff1692037f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd91a2600101611a86565b0135905038806106f3565b838252602082209a9e9d9c9b9a91601f198416815b818110611e175750919e9f9b9c9d9e6001939185611ccf9897969410611dfd575b505050811b01905561476e565b60001960f88560031b161c199101351690558f8080611df0565b91936020600181928787013581550195019201611dcf565b013590508e806106f3565b8382526020822091601f198416815b818110611e865750908460019594939210611e6c575b505050811b019055611bb3565b60001960f88560031b161c199101351690558d8080611e5f565b83830135855560019094019360209283019201611e49565b7f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff611eef6060840161512b565b1615611af0577f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b34610247576020600319360112610247576020611f5a67ffffffffffffffff611f46614116565b166000526006602052604060002054151590565b6040519015158152f35b346102475760c0600319360112610247576004356024356044359067ffffffffffffffff821680920361024757606435916084359060ff821682036102475760a435916001600160a01b038316918284036102475780600052600a60205260ff600160406000200154166121bd575b50611ff760405183602082015260208152611fef60408261428a565b828787614a65565b860361218f5785600052600d60205260406000206001600160a01b03604051916120208361421a565b5461202e60ff82168461489e565b60081c166020820152519460038610156121795760009561214d579061205391614dfb565b92604051956120618761421a565b600187526020870196338852818752600d602052604087209051976003891015612139578798612136985060ff60ff198454169116178255517fffffffffffffffffffffff0000000000000000000000000000000000000000ff74ffffffffffffffffffffffffffffffffffffffff0083549260081b1691161790556040519285845260208401527fd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf160403394a4337f0000000000000000000000000000000000000000000000000000000000000000614b4e565b80f35b602488634e487b7160e01b81526021600452fd5b602486887f9b91b78c000000000000000000000000000000000000000000000000000000008252600452fd5b634e487b7160e01b600052602160045260246000fd5b857fcb537aa40000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6121d433600052600c602052604060002054151590565b611fd3577f6c46a9b5000000000000000000000000000000000000000000000000000000006000526004523360245260446000fd5b34610247576122173661441f565b6001600160a01b036004541633036128825760a0813603126102475760405161223f816141fe565b8135815261224f6020830161412d565b9060208101918252604083013567ffffffffffffffff8111610247576122789036908501614300565b9160408201928352606084013567ffffffffffffffff8111610247576122a19036908601614300565b936060830194855260808101359067ffffffffffffffff8211610247570136601f820112156102475780356122d581614452565b916122e3604051938461428a565b81835260208084019260061b8201019036821161024757602001915b81831061284a575050506080830152519067ffffffffffffffff8216905192519351918251830194602086019360208188031261024757602081015167ffffffffffffffff811161024757019560a090879003126102475760405191612364836141fe565b60208701518352612377604088016150ef565b916020840192835261238b606089016150ef565b916040850192835260808901519860ff8a168a036102475760608601998a5260a081015167ffffffffffffffff811161024757602091010187601f820112156102475780516123d9816142ad565b986123e76040519a8b61428a565b818a526020828401011161024757612405916020808b01910161419a565b6080850196875277ffffffffffffffff0000000000000000000000000000000060405191632cbc26bb60e01b835260801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561169e5760009161282b575b506117555761248c8186614783565b156127ed57509560ff6124b36124ee936124e2989961ffff80885193511691511691615433565b6124dd6124c8879a939a518587511690614dfb565b996124d68587511684614dfb565b9751614644565b614644565b91511685519188614a65565b9384600052600d60205260406000209161254f826040519461250f8661421a565b549561251e60ff88168761489e565b6001600160a01b03602087019760081c16875288600052600d6020526040600020600260ff19825416179055615758565b600093835160038110156127d9576126845750506000935160208180518101031261268057602001516001600160a01b038116809103612680576040517f40c10f190000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602481019190915260208180604481010381876001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1801561267557612656575b505b51906003821015612179576126536060927f33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c946040519384526020840152604083019061450a565ba3005b61266e9060203d60201161169757611688818361428a565b5085612609565b6040513d86823e3d90fd5b8480fd5b93909450825160038110156127c55760010361279957506126ad846001600160a01b0392614644565b9251166126ff60206126bf8686614aac565b6040517f40c10f19000000000000000000000000000000000000000000000000000000008152306004820152602481019190915291829081906044820190565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1801561169e5761277a575b508280612749575b505061260b565b612773917f000000000000000000000000000000000000000000000000000000000000000061525d565b8582612742565b6127929060203d60201161169757611688818361428a565b508661273a565b80867fb196a44a0000000000000000000000000000000000000000000000000000000060249352600452fd5b602482634e487b7160e01b81526021600452fd5b602486634e487b7160e01b81526021600452fd5b61107e906040519182917f24eb47e50000000000000000000000000000000000000000000000000000000083526020600484015260248301906141bd565b612844915060203d60201161169757611688818361428a565b8961247d565b60408336031261024757602060409182516128648161421a565b61286d86614186565b815282860135838201528152019201916122ff565b7fd7f73334000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b34610247576020600319360112610247577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917460206001600160a01b036128f4614170565b6128fc614f02565b168073ffffffffffffffffffffffffffffffffffffffff196009541617600955604051908152a1005b34610247576000600319360112610247576000546001600160a01b03811633036129a15773ffffffffffffffffffffffffffffffffffffffff19600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b34610247576020600319360112610247576020611f5a6001600160a01b036129f1614170565b16600052600c602052604060002054151590565b346102475760206003193601126102475767ffffffffffffffff612a27614116565b606060c0604051612a3781614252565b600081526000602082015260006040820152600083820152600060808201528260a0820152015216600052600a60205260606040600020610d33612a796157cc565b611a15604051612a8881614252565b84548152612b3560018601549563ffffffff602084019760ff81161515895261ffff60408601818360081c168152818c880191818560181c1683528560808a019560281c168552612aee6003612ae060028a016147fa565b9860a08c01998a52016147fa565b9860c08101998a526040519e8f9e8f9260408452516040840152511515910152511660808c0152511660a08a0152511660c08801525160e0808801526101208701906141bd565b9051603f19868303016101008701526141bd565b346102475760006003193601126102475760206001600160a01b0360095416604051908152f35b3461024757600060031936011261024757610d3361114a6157cc565b3461024757612b9a366143a1565b612ba5929192614f02565b67ffffffffffffffff8216612bc7816000526006602052604060002054151590565b15612be25750610c5292612bdc9136916142c9565b90614f40565b631e670e4b60e01b60005260045260246000fd5b3461024757612c1e612c26612c0a3661434f565b9491612c17939193614f02565b369161446a565b92369161446a565b7f000000000000000000000000000000000000000000000000000000000000000015612d2c5760005b8251811015612cb557806001600160a01b03612c6d60019386614900565b5116612c7881615cd8565b612c84575b5001612c4f565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a184612c7d565b5060005b8151811015610c5257806001600160a01b03612cd760019385614900565b51168015612d2657612ce8816158bc565b612cf5575b505b01612cb9565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a183612ced565b50612cef565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b3461024757604060031936011261024757612d6f614116565b60243567ffffffffffffffff811161024757602091612d95611f5a923690600401614300565b90614783565b346102475760206003193601126102475760043567ffffffffffffffff81116102475780600401906101006003198236030112610247576000604051612de081614236565b52612e0d612e03612dfe612df760c4850186614709565b36916142c9565b614d20565b6064830135614dfb565b9060848101612e1b8161475a565b6001600160a01b03807f00000000000000000000000000000000000000000000000000000000000000001691160361179e5750602481019277ffffffffffffffff00000000000000000000000000000000612e758561476e565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561169e57600091613150575b506117555767ffffffffffffffff612ee48561476e565b16612efc816000526006602052604060002054151590565b156117285760206001600160a01b0360045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa90811561169e57600091613131575b50156116aa57612f678461476e565b90612f7d60a4840192612d95612df78585614709565b156130ea575050604490612f9983612f948661476e565b615758565b01612ff3602083612fa98461475a565b60405193849283927f40c10f1900000000000000000000000000000000000000000000000000000000845260048401602090939291936001600160a01b0360408201951681520152565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1801561169e5767ffffffffffffffff61307161306b6020977ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0956080956130cd575b5061476e565b9461475a565b936001600160a01b0360405195817f000000000000000000000000000000000000000000000000000000000000000016875233898801521660408601528560608601521692a2806040516130c481614236565b52604051908152f35b6130e3908a3d8c1161169757611688818361428a565b5089613065565b6130f49250614709565b61107e6040519283927f24eb47e5000000000000000000000000000000000000000000000000000000008452602060048501526024840191614934565b61314a915060203d60201161169757611688818361428a565b85612f58565b613169915060203d60201161169757611688818361428a565b85612ecd565b346102475760a060031936011261024757613188614116565b6024359060443567ffffffffffffffff8111610247576131ac903690600401614142565b91606435916001600160a01b0383168093036102475760843567ffffffffffffffff8111610247576131e2903690600401614142565b50506131ec6146f0565b50604051936131fa856141e2565b60008552602085019260008452604086019260008452606087016000815260606080604051613228816141fe565b828152826020820152826040820152600083820152015283158015613659575b6136005760008235613651575b6020851161363f575b156136005767ffffffffffffffff831693604051632cbc26bb60e01b815277ffffffffffffffff000000000000000000000000000000008560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561169e576000916135e1575b50611755576132ee33615c63565b613305856000526006602052604060002054151590565b156135b35784600052600a60205260406000209485548b116135815750986134176134259260ff6134ab9b9c63ffffffff60018a015461336761ffff8260081c169c8d9661335c61ffff8560181c1680998c615433565b928382935252614aac565b8d5260281c1680613516575061ffff6133d561338560038c016147fa565b985b60405197613394896141fe565b8852602088019c8d52604088019586526060880193857f000000000000000000000000000000000000000000000000000000000000000016855236916142c9565b9360808701948552816040519c8d986020808b01525160408a01525116606088015251166080860152511660a08401525160a060c084015260e08301906141bd565b03601f19810186528561428a565b602095869460405190613438878361428a565b6000825261345460026040519761344e896141fe565b016147fa565b8652868601526040850152606084015260808301526001600160a01b0360045416906040518097819482937f20487ded00000000000000000000000000000000000000000000000000000000845260048401614577565b03915afa92831561169e576000936134e4575b508260409452518184516134d18161421a565b8481520190815283519283525190820152f35b9392508184813d831161350f575b6134fc818361428a565b81010312610247576040935192936134be565b503d6134f2565b6133d561ffff916040519061352a8261421a565b81526020810160018152604051917f181dcf1000000000000000000000000000000000000000000000000000000000602084015251602483015251151560448201526044815261357b60648261428a565b98613387565b8a907f58dd87c50000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b847fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6135fa915060203d60201161169757611688818361428a565b8b6132e0565b508261107e6040519283927fa3c8cf09000000000000000000000000000000000000000000000000000000008452602060048501526024840191614934565b60208301351561325e5750600161325e565b506001613255565b5060408411613248565b34610247576000600319360112610247576020610349614651565b3461024757600060031936011261024757602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760206003193601126102475760206136d7614170565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b3461024757600060031936011261024757610d3360405161377560608261428a565b603581527f4275726e4d696e745769746845787465726e616c4d696e74657246617374547260208201527f616e73666572546f6b656e506f6f6c20312e362e30000000000000000000000060408201526040519182916020835260208301906141bd565b60c0600319360112610247576137ed614116565b60643567ffffffffffffffff81116102475761380d903690600401614142565b60843592916001600160a01b03841684036102475760a43567ffffffffffffffff811161024757613842903690600401614142565b505060405190613851826141e2565b600082526000602083015260006040830152600060608301526060608060405161387a816141fe565b828152826020820152826040820152600083820152015282158015613fa9575b613f515760008135613fa1575b60208411613f8f575b15613f5157604051632cbc26bb60e01b815277ffffffffffffffff000000000000000000000000000000008560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561169e57600091613f32575b506117555761393433615c63565b61395567ffffffffffffffff85166000526006602052604060002054151590565b15613efa5767ffffffffffffffff8416600052600a6020526040600020805460243511613ebc576060956001820154613a8c61ffff8260081c1663ffffffff61ffff8460181c16936139be6139ad8685602435615433565b90818f8d01528060408d0152614aac565b60208a015260281c1680613e5857506139d9600386016147fa565b925b604051916139e8836141fe565b60243583526020830152604082015260ff7f0000000000000000000000000000000000000000000000000000000000000000168a820152613a7e613a2d368a896142c9565b6080830190815260ff6040519c8d946020808701528051604087015261ffff6020820151168287015261ffff604082015116608087015201511660a08401525160a060c084015260e08301906141bd565b03601f1981018a528961428a565b604051602098613a9c8a8361428a565b60008252613ab260026040519661344e886141fe565b85528985015260408401526001600160a01b038216606084015260808301526001600160a01b03600454168760405180927f20487ded0000000000000000000000000000000000000000000000000000000082528180613b16888d60048401614577565b03915afa90811561169e57600091613e2b575b508452613b3860243587614ab9565b60208401516044358111613df857506000877f0000000000000000000000000000000000000000000000000000000000000000613b79602435303384614b4e565b613bb16001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016918260243591614bb8565b6024604051809481937f42966c68000000000000000000000000000000000000000000000000000000008352833560048401525af1801561169e57613ddb575b506001600160a01b038116613d88575b506001600160a01b0360045416948660405180977f96f4e9f90000000000000000000000000000000000000000000000000000000082528180613c48878760048401614577565b039134905af194851561169e578796600096613d53575b5091613d488697927ffa7d3740fa7611df3f0d8d8c3aa1ed57c4fffaf2dcd0c47535f18a4774b44acd9467ffffffffffffffff613d38613ca560208b0151602435614644565b95606060408c01519b0151905190613cec8d613cc38d8836916142c9565b908a7f000000000000000000000000000000000000000000000000000000000000000091614a65565b9b604051998a998a5260ff7f000000000000000000000000000000000000000000000000000000000000000016908a01526040890152606088015260c0608088015260c08701906141bd565b9285840360a08701521696614934565b0390a4604051908152f35b878193989297503d8311613d81575b613d6c818361428a565b81010312610247575186959094613d48613c5f565b503d613d62565b613dd590613da2855130336001600160a01b038516614b4e565b8451906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000009116614bb8565b86613c01565b613df190883d8a1161169757611688818361428a565b5087613bf1565b7f61acdb930000000000000000000000000000000000000000000000000000000060005260045260443560245260446000fd5b90508781813d8311613e51575b613e42818361428a565b81010312610247575188613b29565b503d613e38565b60405190613e658261421a565b81526020810160018152604051917f181dcf10000000000000000000000000000000000000000000000000000000006020840152516024830152511515604482015260448152613eb660648261428a565b926139db565b67ffffffffffffffff857f58dd87c5000000000000000000000000000000000000000000000000000000006000521660045260243560245260446000fd5b67ffffffffffffffff847fa9902c7e000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b613f4b915060203d60201161169757611688818361428a565b86613926565b8261107e6040519283927fa3c8cf09000000000000000000000000000000000000000000000000000000008452602060048501526024840191614934565b6020820135156138b0575060016138b0565b5060016138a7565b506040831161389a565b3461024757602060031936011261024757600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361024757817ff6f46ff9000000000000000000000000000000000000000000000000000000006020931490811561408b575b811561402e575b5015158152f35b7f85572ffb00000000000000000000000000000000000000000000000000000000811491508115614061575b5083614027565b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150148361405a565b90507faff2afbf00000000000000000000000000000000000000000000000000000000811480156140ed575b80156140c4575b90614020565b507f01ffc9a70000000000000000000000000000000000000000000000000000000081146140be565b507f0e64dd290000000000000000000000000000000000000000000000000000000081146140b7565b6004359067ffffffffffffffff8216820361024757565b359067ffffffffffffffff8216820361024757565b9181601f840112156102475782359167ffffffffffffffff8311610247576020838186019501011161024757565b600435906001600160a01b038216820361024757565b35906001600160a01b038216820361024757565b60005b8381106141ad5750506000910152565b818101518382015260200161419d565b90601f19601f6020936141db8151809281875287808801910161419a565b0116010190565b6080810190811067ffffffffffffffff82111761087957604052565b60a0810190811067ffffffffffffffff82111761087957604052565b6040810190811067ffffffffffffffff82111761087957604052565b6020810190811067ffffffffffffffff82111761087957604052565b60e0810190811067ffffffffffffffff82111761087957604052565b6060810190811067ffffffffffffffff82111761087957604052565b90601f601f19910116810190811067ffffffffffffffff82111761087957604052565b67ffffffffffffffff811161087957601f01601f191660200190565b9291926142d5826142ad565b916142e3604051938461428a565b829481845281830111610247578281602093846000960137010152565b9080601f830112156102475781602061431b933591016142c9565b90565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501948460051b01011161024757565b60406003198201126102475760043567ffffffffffffffff8111610247578161437a9160040161431e565b929092916024359067ffffffffffffffff82116102475761439d9160040161431e565b9091565b9060406003198301126102475760043567ffffffffffffffff8116810361024757916024359067ffffffffffffffff82116102475761439d91600401614142565b906020808351928381520192019060005b8181106144005750505090565b82516001600160a01b03168452602093840193909201916001016143f3565b6020600319820112610247576004359067ffffffffffffffff8211610247576003198260a0920301126102475760040190565b67ffffffffffffffff81116108795760051b60200190565b92919061447681614452565b93614484604051958661428a565b602085838152019160051b810192831161024757905b8282106144a657505050565b602080916144b384614186565b81520191019061449a565b9080601f830112156102475781602061431b9335910161446a565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501946060850201011161024757565b9060038210156121795752565b35906001600160801b038216820361024757565b9190826060910312610247576040516145438161426e565b80928035908115158203610247576040614572918193855261456760208201614517565b602086015201614517565b910152565b9067ffffffffffffffff90939293168152604060208201526145be6145a8845160a0604085015260e08401906141bd565b6020850151603f198483030160608501526141bd565b90604084015191603f198282030160808301526020808451928381520193019060005b818110614619575050506080846001600160a01b03606061431b969701511660a084015201519060c0603f19828503019101526141bd565b825180516001600160a01b0316865260209081015181870152604090950194909201916001016145e1565b91908203918211611ec857565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561169e576000916146c1575090565b90506020813d6020116146e8575b816146dc6020938361428a565b81010312610247575190565b3d91506146cf565b604051906146fd8261421a565b60006020838281520152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610247570180359067ffffffffffffffff82116102475760200191813603831361024757565b356001600160a01b03811681036102475790565b3567ffffffffffffffff811681036102475790565b9067ffffffffffffffff61431b92166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b90600182811c921680156147f0575b60208310146147da57565b634e487b7160e01b600052602260045260246000fd5b91607f16916147cf565b906040519182600082549261480e846147c0565b808452936001811690811561487c5750600114614835575b506148339250038361428a565b565b90506000929192526020600020906000915b8183106148605750509060206148339282010138614826565b6020919350806001915483858901015201910190918492614847565b6020935061483395925060ff1991501682840152151560051b82010138614826565b60038210156121795752565b91908110156148ea5760051b810135907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee181360301821215610247570190565b634e487b7160e01b600052603260045260246000fd5b80518210156148ea5760209160051b010190565b91908110156148ea5760051b0190565b91908110156148ea576060020190565b601f8260209493601f19938186528686013760008582860101520116010190565b60405190614962826141fe565b60006080838281528260208201528260408201528260608201520152565b9060405161498d816141fe565b60806001829460ff81546001600160801b038116865263ffffffff81861c16602087015260a01c161515604085015201546001600160801b0381166060840152811c910152565b67ffffffffffffffff16600052600760205261431b60046040600020016147fa565b818110614a01575050565b600081556001016149f6565b81810292918115918404141715611ec857565b9190601f8111614a2f57505050565b614833926000526020600020906020601f840160051c83019310614a5b575b601f0160051c01906149f6565b9091508190614a4e565b9290614a98614aa69260ff60405195869460208601988952604086015216606084015260808084015260a08301906141bd565b03601f19810183528261428a565b51902090565b91908201809211611ec857565b67ffffffffffffffff7fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da817894491169182600052600760205280614b2a60406000206001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001692839161545a565b604080516001600160a01b039092168252602082019290925290819081015b0390a2565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000060208201526001600160a01b039283166024820152929091166044830152606482019290925261483391614bb382608481015b03601f19810184528361428a565b61563f565b91909181158015614c86575b15614c1c576040517f095ea7b30000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909316602484015260448301919091526148339190614bb38260648101614ba5565b608460405162461bcd60e51b815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b0384166024820152602081806044810103816001600160a01b0386165afa90811561169e57600091614cee575b5015614bc4565b90506020813d602011614d18575b81614d096020938361428a565b81010312610247575138614ce7565b3d9150614cfc565b80518015614d9057602003614d5257805160208281019183018390031261024757519060ff8211614d52575060ff1690565b61107e906040519182917f953576f70000000000000000000000000000000000000000000000000000000083526020600484015260248301906141bd565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff8211611ec857565b60ff16604d8111611ec857600a0a90565b8115614de5570490565b634e487b7160e01b600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff811692828414614ee357828411614eb95790614e4091614db6565b91604d60ff8416118015614e9e575b614e6857505090614e6261431b92614dca565b90614a0d565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b50614ea883614dca565b8015614de557600019048411614e4f565b614ec291614db6565b91604d60ff841611614e6857505090614edd61431b92614dca565b90614ddb565b5050505090565b90816020910312610247575180151581036102475790565b6001600160a01b03600154163303614f1657565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b908051156108c85767ffffffffffffffff81516020830120921691826000526007602052614f7581600560406000200161596d565b156150ab5760005260086020526040600020815167ffffffffffffffff811161087957614fac81614fa684546147c0565b84614a20565b6020601f82116001146150215791615000827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea9593614b4995600091615016575b506000198260011b9260031b1c19161790565b90556040519182916020835260208301906141bd565b905084015138614fed565b601f1982169083600052806000209160005b818110615093575092614b499492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea98961061507a575b5050811b019055610f06565b85015160001960f88460031b161c19169055388061506e565b9192602060018192868a015181550194019201615033565b509061107e6040519283927f393b8ad200000000000000000000000000000000000000000000000000000000845260048401526040602484015260448301906141bd565b519061ffff8216820361024757565b357fffffffff00000000000000000000000000000000000000000000000000000000811681036102475790565b3563ffffffff811681036102475790565b3561ffff811681036102475790565b3580151581036102475790565b67ffffffffffffffff166000818152600660205260409020549092919015615248579161524560e09261521a856151af7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b97615322565b8460005260076020526151c6816040600020615a7d565b6151cf83615322565b8460005260076020526151e9836002604060002001615a7d565b60405194855260208501906001600160801b0360408092805115158552826020820151166020860152015116910152565b60808301906001600160801b0360408092805115158552826020820151166020860152015116910152565ba1565b82631e670e4b60e01b60005260045260246000fd5b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000060208201526001600160a01b039092166024830152604482019290925261483391614bb38260648101614ba5565b6152b7614955565b506001600160801b036060820151166001600160801b03808351169161530260208501936152fc6152ef63ffffffff87511642614644565b8560808901511690614a0d565b90614aac565b8082101561531b57505b16825263ffffffff4216905290565b905061530c565b8051156153a7576001600160801b036040820151166001600160801b036020830151161061534d5750565b6064906153a5604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565bfd5b6001600160801b036040820151161580159061541d575b6153c55750565b6064906153a5604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565b506001600160801b0360208201511615156153be565b6154569061ffff61271061544d8282969897981684614a0d565b04951690614a0d565b0490565b9182549060ff8260a01c16158015615637575b615631576001600160801b03821691600185019081546154a063ffffffff6001600160801b0383169360801c1642614644565b9081615593575b505084811061555457508383106154e95750506154cd6001600160801b03928392614644565b16166fffffffffffffffffffffffffffffffff19825416179055565b5460801c916154f88185614644565b92600019810190808211611ec85761551b615520926001600160a01b0396614aac565b614ddb565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b82856001600160a01b03927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b828692939611615607576155ae926152fc9160801c90614a0d565b808410156156025750825b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff00000000000000000000000000000000161786559238806154a7565b6155b9565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b50821561546d565b6001600160a01b036156c1911691604092600080855193615660878661428a565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d15615750573d916156a5836142ad565b926156b28751948561428a565b83523d6000602085013e615ea0565b805190816156ce57505050565b6020806156df938301019101614eea565b156156e75750565b6084905162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606091615ea0565b67ffffffffffffffff7f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c91169182600052600760205280614b2a60026040600020016001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001692839161545a565b60405190600b548083528260208101600b60005260206000209260005b8181106157fe5750506148339250038361428a565b84548352600194850194879450602090930192016157e9565b906040519182815491828252602082019060005260206000209260005b8181106158495750506148339250038361428a565b8454835260019485019487945060209093019201615834565b80548210156148ea5760005260206000200190600090565b8054906801000000000000000082101561087957816158a19160016158b894018155615862565b81939154906000199060031b92831b921b19161790565b9055565b806000526003602052604060002054156000146158f5576158de81600261587a565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146158f55761591d81600b61587a565b600b5490600052600c602052604060002055600190565b806000526006602052604060002054156000146158f55761595681600561587a565b600554906000526006602052604060002055600190565b60008281526001820160205260409020546159a4578061598f8360019361587a565b80549260005201602052604060002055600190565b5050600090565b805480156159d35760001901906159c28282615862565b60001982549160031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c602052604090205480156159a4576000198101818111611ec857600b54906000198201918211611ec857808203615a43575b505050615a2f600b6159ab565b600052600c60205260006040812055600190565b615a65615a546158a193600b615862565b90549060031b1c928392600b615862565b9055600052600c602052604060002055388080615a22565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1991615b866060928054615aba63ffffffff8260801c1642614644565b9081615bbc575b50506001600160801b036001816020860151169282815416808510600014615bb457508280855b16166fffffffffffffffffffffffffffffffff19825416178155615b528651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601516fffffffffffffffffffffffffffffffff1960809190911b16939092166001600160801b031692909217910155565b61524560405180926001600160801b0360408092805115158552826020820151166020860152015116910152565b838091615ae8565b6001600160801b0391615be8839283615be16001880154948286169560801c90614a0d565b9116614aac565b80821015615c5c57505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff92909116929092161673ffffffffffffffffffffffffffffffffffffffff19909116174260801b73ffffffff00000000000000000000000000000000161781553880615ac1565b9050615bf2565b7f0000000000000000000000000000000000000000000000000000000000000000615c8b5750565b6001600160a01b031680600052600360205260406000205415615cab5750565b7fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b60008181526003602052604090205480156159a4576000198101818111611ec857600254906000198201918211611ec857818103615d32575b505050615d1e60026159ab565b600052600360205260006040812055600190565b615d54615d436158a1936002615862565b90549060031b1c9283926002615862565b90556000526003602052604060002055388080615d11565b60008181526006602052604090205480156159a4576000198101818111611ec857600554906000198201918211611ec857818103615dc6575b505050615db260056159ab565b600052600660205260006040812055600190565b615de8615dd76158a1936005615862565b90549060031b1c9283926005615862565b90556000526006602052604060002055388080615da5565b906001820191816000528260205260406000205490811515600014615e9757600019820191808311611ec85781546000198101908111611ec8578381615e4e9503615e60575b5050506159ab565b60005260205260006040812055600190565b615e80615e706158a19386615862565b90549060031b1c92839286615862565b905560005284602052604060002055388080615e46565b50505050600090565b91929015615f015750815115615eb4575090565b3b15615ebd5790565b606460405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015615f145750805190602001fd5b61107e9060405191829162461bcd60e51b83526020600484015260248301906141bd56fea164736f6c634300081a000a",
}

var BurnMintWithExternalMinterFastTransferTokenPoolABI = BurnMintWithExternalMinterFastTransferTokenPoolMetaData.ABI

var BurnMintWithExternalMinterFastTransferTokenPoolBin = BurnMintWithExternalMinterFastTransferTokenPoolMetaData.Bin

func DeployBurnMintWithExternalMinterFastTransferTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, minter common.Address, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnMintWithExternalMinterFastTransferTokenPool, error) {
	parsed, err := BurnMintWithExternalMinterFastTransferTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintWithExternalMinterFastTransferTokenPoolBin), backend, minter, token, localTokenDecimals, allowlist, rmnProxy, router)
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipSendToken", destinationChainSelector, amount, maxFastTransferFee, receiver, settlementFeeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, settlementFeeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, settlementFeeToken, extraArgs)
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
	FillerFee                *big.Int
	PoolFee                  *big.Int
	DestinationPool          []byte
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
	return common.HexToHash("0xfa7d3740fa7611df3f0d8d8c3aa1ed57c4fffaf2dcd0c47535f18a4774b44acd")
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

	CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (*types.Transaction, error)

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
