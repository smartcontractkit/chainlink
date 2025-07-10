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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSendToken\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"computeFillId\",\"inputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"fastFill\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAccumulatedPoolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedFillers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCcipSendTokenFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structIFastTransferPool.Quote\",\"components\":[{\"name\":\"ccipSettlementFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfig\",\"components\":[{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFillInfo\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.FillInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIFastTransferPool.FillState\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAllowedFiller\",\"inputs\":[{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateDestChainConfig\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfigUpdateArgs[]\",\"components\":[{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFillerAllowList\",\"inputs\":[{\"name\":\"fillersToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"fillersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPoolFees\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestinationPoolUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destinationPool\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"filler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferRequested\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"fillerFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fillerReimbursementAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFeeAccumulated\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"prevState\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIFastTransferPool.FillState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FillerAllowListUpdated\",\"inputs\":[{\"name\":\"addFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"removeFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolFeeWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyFilledOrSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadySettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"FillerNotAllowlisted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientPoolFees\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFillId\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"QuoteFeeExceedsUserMaxLimit\",\"inputs\":[{\"name\":\"quoteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferAmountExceedsMaxFillAmount\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461032857616698803803809161001d8285610559565b8339810160c082820312610328576100348261057c565b60208301516001600160a01b03811693919290918483036103285761005b60408201610590565b60608201519091906001600160401b0381116103285781019280601f85011215610328578351936001600160401b038511610543578460051b9060208201956100a76040519788610559565b865260208087019282010192831161032857602001905b82821061052b575050506100e060a06100d96080840161057c565b920161057c565b93331561051a57600180546001600160a01b0319163317905586158015610509575b80156104f8575b6104e75760805260c05260405163313ce56760e01b8152602081600481895afa600091816104ab575b50610480575b5060a052600480546001600160a01b0319166001600160a01b0384169081179091558151151560e0819052909190610357575b501561034157610120526001600160a01b03166101008190526040516321df0da760e01b815290602090829060049082905afa908115610335576000916102f6575b506001600160a01b0316908181036102df57604051615f59908161073f8239608051818181611352015281816113bc015281816115280152818161215a0152818161279701528181612e9d015281816130f8015281816136d20152818161371f01528181613b220152818161463601528181614a8a0152615779015260a05181818161163e015281816133ff01528181613688015281816139cb01528181613c7e01528181613cd601528181614d360152614da0015260c051818181610b2e015281816114310152818161249401528181612f130152818161330201526138c0015260e051818181610ae901528181612c9f0152615c45015261010051818181610223015281816115540152818161262401528181612755015281816130780152613b5c01526101205181613d980152f35b63f902523f60e01b60005260045260245260446000fd5b90506020813d60201161032d575b8161031160209383610559565b81010312610328576103229061057c565b386101ad565b600080fd5b3d9150610304565b6040513d6000823e3d90fd5b6335fdcccd60e21b600052600060045260246000fd5b9192906020926040519261036b8585610559565b60008452600036813760e0511561046f5760005b84518110156103e6576001906001600160a01b0361039d828861059e565b5116876103a9826105e0565b6103b6575b50500161037f565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a138876103ae565b50919490925060005b8351811015610463576001906001600160a01b0361040d828761059e565b5116801561045d578661041f826106de565b61042d575b50505b016103ef565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a13886610424565b50610427565b5092509290503861016b565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036104945750610138565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d6020116104df575b816104c760209383610559565b81010312610328576104d890610590565b9038610132565b3d91506104ba565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03821615610109565b506001600160a01b03851615610102565b639b15e16f60e01b60005260046000fd5b602080916105388461057c565b8152019101906100be565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b0382119082101761054357604052565b51906001600160a01b038216820361032857565b519060ff8216820361032857565b80518210156105b25760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156105b25760005260206000200190600090565b60008181526003602052604090205480156106d75760001981018181116106c1576002546000198101919082116106c157818103610670575b505050600254801561065a57600019016106348160026105c8565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b6106a96106816106929360026105c8565b90549060031b1c92839260026105c8565b819391549060031b91821b91600019901b19161790565b90556000526003602052604060002055388080610619565b634e487b7160e01b600052601160045260246000fd5b5050600090565b8060005260036020526040600020541560001461073857600254680100000000000000008110156105435761071f61069282600185940160025560026105c8565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a714613f3b57508063055befd4146137c9578063181f5a771461374357806321df0da7146136ff578063240028e8146136ac57806324f65ee71461366e5780632b2c0eb4146136535780632e7aa8c8146131e65780633907753714612e125780634c5ef0ed14612dcd57806354c8a4f314612c6d57806362ddd3c414612bea5780636609f59914612bce5780636d3d1a5814612ba75780636def4ce714612a6357806378b410f214612a2957806379ba5097146129785780637d54534e146128f857806385572ffb1461225157806387f060d014611fac5780638926f54f14611f675780638a18dcbd14611a925780638da5cb5b14611a6b578063929ea5ba14611961578063962d4020146118255780639a4575b9146113815780639fe280f5146112ee578063a42a7b8b146111bc578063a7cd63b71461114e578063abe1c1e8146110df578063acfecf9114610fba578063af58d59f14610f70578063b0f479a114610f49578063b794658014610f11578063c0d7865514610e6d578063c4bffe2b14610d5b578063c75eea9c14610cbb578063cf7401f314610b52578063dc0bd97114610b0e578063e0351e1314610ad1578063e8a1da171461035c578063eeebc67414610304578063f2fde38b1461024c5763f36675171461020357600080fd5b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b34610247576020600319360112610247576001600160a01b0361026d6140f8565b610275614ea4565b163381146102da57807fffffffffffffffffffffffff000000000000000000000000000000000000000060005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760806003193601126102475760443560ff811681036102475760643567ffffffffffffffff811161024757602091610348610354923690600401614288565b906024356004356149ed565b604051908152f35b346102475761036a366142d7565b919092610375614ea4565b6000905b82821061092c5750505060009063ffffffff42165b81831061039757005b6103a2838386614853565b926101208436031261024757604051936103bb85614186565b6103c4816140b5565b8552602081013567ffffffffffffffff81116102475781019336601f860112156102475784356103f3816143da565b956104016040519788614212565b81875260208088019260051b820101903682116102475760208101925b8284106108fd575050505060208601948552604082013567ffffffffffffffff8111610247576104519036908401614288565b906040870191825261047b61046936606086016144b3565b936060890194855260c03691016144b3565b946080880195865261048d84516152dd565b61049786516152dd565b825151156108d3576104b367ffffffffffffffff895116615909565b1561089a5767ffffffffffffffff885116600052600760205260406000206105c285516001600160801b03604082015116906105956001600160801b036020830151169151151583608060405161050981614186565b858152602081018a905260408101849052606081018690520152855474ff000000000000000000000000000000000000000091151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166001600160801b0384161773ffffffff00000000000000000000000000000000608089901b1617178555565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176001830155565b6106b587516001600160801b03604082015116906106886001600160801b03602083015116915115158360806040516105fa81614186565b858152602081018a9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166001600160801b0385161773ffffffff0000000000000000000000000000000060808a901b161791151560a01b74ff000000000000000000000000000000000000000016919091179055565b60809190911b6fffffffffffffffffffffffffffffffff19166001600160801b0391909116176003830155565b6004845191019080519067ffffffffffffffff8211610884576106e2826106dc8554614769565b856149a8565b602090601f831160011461081d57610712929160009183610812575b50506000198260011b9260031b1c19161790565b90555b60005b8751805182101561074d57906107476001926107408367ffffffffffffffff8e5116926148a9565b5190614ee2565b01610718565b505097967f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c293929196509461080767ffffffffffffffff60019751169251935191516107dc6107b060405196879687526101006020880152610100870190614145565b9360408601906001600160801b0360408092805115158552826020820151166020860152015116910152565b60a08401906001600160801b0360408092805115158552826020820151166020860152015116910152565b0390a101919261038e565b015190508d806106fe565b90601f1983169184600052816000209260005b81811061086c5750908460019594939210610853575b505050811b019055610715565b015160001960f88460031b161c191690558c8080610846565b92936020600181928786015181550195019301610830565b634e487b7160e01b600052604160045260246000fd5b67ffffffffffffffff8851167f1d5ad3c50000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b833567ffffffffffffffff8111610247576020916109218392833691870101614288565b81520193019261041e565b9092919367ffffffffffffffff61094c6109478688866148bd565b614717565b169261095784615d4c565b15610aa35783600052600760205261097560056040600020016157ec565b9260005b84518110156109b1576001908660005260076020526109aa60056040600020016109a383896148a9565b5190615de0565b5001610979565b50939094919592508060005260076020526005604060002060008155600060018201556000600282015560006003820155600481016109f08154614769565b9081610a60575b5050018054906000815581610a3f575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a1019091610379565b6000526020600020908101905b81811015610a075760008155600101610a4c565b81601f60009311600114610a785750555b88806109f7565b81835260208320610a9391601f01861c81019060010161497e565b8082528160208120915555610a71565b837f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102475760006003193601126102475760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760e060031936011261024757610b6b61409e565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc36011261024757604051610ba1816141f6565b60243580151581036102475781526044356001600160801b03811681036102475760208201526064356001600160801b038116810361024757604082015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7c3601126102475760405190610c16826141f6565b608435801515810361024757825260a4356001600160801b038116810361024757602083015260c4356001600160801b03811681036102475760408301526001600160a01b036009541633141580610ca6575b610c7857610c76926150fa565b005b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b506001600160a01b0360015416331415610c69565b346102475760206003193601126102475767ffffffffffffffff610cdd61409e565b610ce56148dd565b50166000526007602052610d57610d07610d026040600020614908565b61526a565b6040519182918291909160806001600160801b038160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b34610247576000600319360112610247576040516005548082528160208101600560005260206000209260005b818110610e54575050610d9d92500382614212565b805190610dc2610dac836143da565b92610dba6040519485614212565b8084526143da565b90601f1960208401920136833760005b8151811015610e04578067ffffffffffffffff610df1600193856148a9565b5116610dfd82876148a9565b5201610dd2565b5050906040519182916020830190602084525180915260408301919060005b818110610e31575050500390f35b825167ffffffffffffffff16845285945060209384019390920191600101610e23565b8454835260019485019486945060209093019201610d88565b3461024757602060031936011261024757610e866140f8565b610e8e614ea4565b6001600160a01b0381169081156108d357600480547fffffffffffffffffffffffff000000000000000000000000000000000000000081169093179055604080516001600160a01b0393841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491819081015b0390a1005b3461024757602060031936011261024757610d57610f35610f3061409e565b61495c565b604051918291602083526020830190614145565b346102475760006003193601126102475760206001600160a01b0360045416604051908152f35b346102475760206003193601126102475767ffffffffffffffff610f9261409e565b610f9a6148dd565b50166000526007602052610d57610d07610d026002604060002001614908565b346102475767ffffffffffffffff610fd136614329565b929091610fdc614ea4565b1690610ff5826000526006602052604060002054151590565b156110b1578160005260076020526110266005604060002001611019368685614251565b6020815191012090615de0565b1561106a577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7691926110656040519283926020845260208401916145d9565b0390a2005b6110ad906040519384937f74f23c7c00000000000000000000000000000000000000000000000000000000855260048501526040602485015260448401916145d9565b0390fd5b507f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b34610247576020600319360112610247576110f8614699565b50600435600052600d6020526040806000206001600160a01b0382519161111e836141a2565b5461112c60ff821684614847565b81602084019160081c1681526111458451809451614492565b51166020820152f35b34610247576000600319360112610247576040516002548082526020820190600260005260206000209060005b8181106111a657610d578561119281870382614212565b60405191829160208352602083019061436a565b825484526020909301926001928301920161117b565b346102475760206003193601126102475767ffffffffffffffff6111de61409e565b1660005260076020526111f760056040600020016157ec565b805190601f1961121f611209846143da565b936112176040519586614212565b8085526143da565b0160005b8181106112dd57505060005b81518110156112775780611245600192846148a9565b51600052600860205261125b60406000206147a3565b61126582866148a9565b5261127081856148a9565b500161122f565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b8282106112b057505050500390f35b919360206112cd82603f1960019597998495030186528851614145565b96019201920185949391926112a1565b806060602080938701015201611223565b34610247576020600319360112610247576113076140f8565b61130f614ea4565b6113176145fa565b908161131f57005b60206001600160a01b0382611376857f738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599957f0000000000000000000000000000000000000000000000000000000000000000615218565b6040519485521692a2005b346102475761138f366143a7565b6060602060405161139f816141a2565b8281520152608081016113b181614703565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116036117e657506020810177ffffffffffffffff0000000000000000000000000000000061140a82614717565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116e6576000916117c7575b5061179d5761147b61147660408401614703565b615c43565b67ffffffffffffffff61148d82614717565b166114a5816000526006602052604060002054151590565b156117705760206001600160a01b0360045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa9081156116e657600091611720575b506001600160a01b031633036116f257600091611526606061151c84614717565b9201358092614a41565b7f00000000000000000000000000000000000000000000000000000000000000009160206001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001661157f848287614b40565b6024604051809781937f42966c680000000000000000000000000000000000000000000000000000000083528760048401525af19283156116e6577ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae1067ffffffffffffffff610f30946116a497611634976116b7575b5061162c61160286614717565b604080516001600160a01b0390971687523360208801528601929092529116929081906060820190565b0390a2614717565b610d5760405160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260208152611672604082614212565b6040519261167f846141a2565b8352602083019081526040519384936020855251604060208601526060850190614145565b9051601f19848303016040850152614145565b6116d89060203d6020116116df575b6116d08183614212565b810190614e8c565b50886115f5565b503d6116c6565b6040513d6000823e3d90fd5b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d602011611768575b8161173960209383614212565b810103126117645751906001600160a01b038216820361176157506001600160a01b036114fb565b80fd5b5080fd5b3d915061172c565b7fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b6117e0915060203d6020116116df576116d08183614212565b83611462565b6117f76001600160a01b0391614703565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b346102475760606003193601126102475760043567ffffffffffffffff8111610247576118569036906004016142a6565b9060243567ffffffffffffffff811161024757611877903690600401614461565b9060443567ffffffffffffffff811161024757611898903690600401614461565b6001600160a01b03600954163314158061194c575b610c7857838614801590611942575b6119185760005b8681106118cc57005b806119126118e06109476001948b8b6148bd565b6118eb8389896148cd565b61190c6119046118fc86898b6148cd565b9236906144b3565b9136906144b3565b916150fa565b016118c3565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b50808614156118bc565b506001600160a01b03600154163314156118ad565b346102475760406003193601126102475760043567ffffffffffffffff811161024757611992903690600401614446565b60243567ffffffffffffffff8111610247576119b2903690600401614446565b906119bb614ea4565b60005b81518110156119ed57806119e66001600160a01b036119df600194866148a9565b51166158d0565b50016119be565b5060005b8251811015611a205780611a196001600160a01b03611a12600194876148a9565b51166159be565b50016119f1565b7ffd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193611a5d83610f0c8660405193849360408552604085019061436a565b90838203602085015261436a565b346102475760006003193601126102475760206001600160a01b0360015416604051908152f35b346102475760206003193601126102475760043567ffffffffffffffff811161024757611ac39036906004016142a6565b611acb614ea4565b60005b818110611ad757005b611ae2818385614853565b60a081017f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000611b31836150a0565b1614611f26575b60208201611b45816150de565b90604084019161ffff80611b58856150de565b1691160161ffff8111611f105761ffff61271091161015611ee6576080840167ffffffffffffffff611b8982614717565b16600052600a60205260406000209460e0810194611ba786836146b2565b600289019167ffffffffffffffff821161088457611bc9826106dc8554614769565b600090601f8311600114611e8257611bf8929160009183611e775750506000198260011b9260031b1c19161790565b90555b611c04846150de565b926001880197885498611c16886150de565b60181b64ffff0000001695611c2a866150ed565b151560c087013597888555606088019c611c438e6150cd565b60281b68ffffffff0000000000169360081b62ffff0016907fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016177fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff16179060ff16171790556101008401611cb890856146b2565b90916003019167ffffffffffffffff821161088457611cdb826106dc8554614769565b600090601f8311600114611e02579180611d1092611d17969594600092611df75750506000198260011b9260031b1c19161790565b9055614717565b93611d21906150de565b94611d2b906150de565b95611d3690836146b2565b9091611d41906150a0565b97611d4b906150cd565b92611d55906150ed565b936040519761ffff899816885261ffff16602088015260408701526060860160e0905260e0860190611d86926145d9565b957fffffffff0000000000000000000000000000000000000000000000000000000016608085015263ffffffff1660a0840152151560c083015267ffffffffffffffff1692037f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd91a2600101611ace565b0135905038806106fe565b838252602082209a9e9d9c9b9a91601f198416815b818110611e5f5750919e9f9b9c9d9e6001939185611d179897969410611e45575b505050811b019055614717565b60001960f88560031b161c199101351690558f8080611e38565b91936020600181928787013581550195019201611e17565b013590508e806106fe565b8382526020822091601f198416815b818110611ece5750908460019594939210611eb4575b505050811b019055611bfb565b60001960f88560031b161c199101351690558d8080611ea7565b83830135855560019094019360209283019201611e91565b7f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff611f37606084016150cd565b1615611b38577f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b34610247576020600319360112610247576020611fa267ffffffffffffffff611f8e61409e565b166000526006602052604060002054151590565b6040519015158152f35b346102475760c0600319360112610247576004356024356044359067ffffffffffffffff821680920361024757606435916084359060ff821682036102475760a435916001600160a01b038316918284036102475780600052600a60205260ff60016040600020015416612205575b5061203f60405183602082015260208152612037604082614212565b8287876149ed565b86036121d75785600052600d60205260406000206001600160a01b0360405191612068836141a2565b5461207660ff821684614847565b60081c166020820152519460038610156121c157600095612195579061209b91614d9d565b92604051956120a9876141a2565b600187526020870196338852818752600d60205260408720905197600389101561218157879861217e985060ff60ff198454169116178255517fffffffffffffffffffffff0000000000000000000000000000000000000000ff74ffffffffffffffffffffffffffffffffffffffff0083549260081b1691161790556040519285845260208401527fd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf160403394a4337f0000000000000000000000000000000000000000000000000000000000000000614ad6565b80f35b602488634e487b7160e01b81526021600452fd5b602486887f9b91b78c000000000000000000000000000000000000000000000000000000008252600452fd5b634e487b7160e01b600052602160045260246000fd5b857fcb537aa40000000000000000000000000000000000000000000000000000000060005260045260246000fd5b61221c33600052600c602052604060002054151590565b61201b577f6c46a9b5000000000000000000000000000000000000000000000000000000006000526004523360245260446000fd5b346102475761225f366143a7565b6001600160a01b036004541633036128ca5760a0813603126102475760405161228781614186565b81358152612297602083016140b5565b9060208101918252604083013567ffffffffffffffff8111610247576122c09036908501614288565b9160408201928352606084013567ffffffffffffffff8111610247576122e99036908601614288565b936060830194855260808101359067ffffffffffffffff8211610247570136601f8201121561024757803561231d816143da565b9161232b6040519384614212565b81835260208084019260061b8201019036821161024757602001915b818310612892575050506080830152519067ffffffffffffffff8216905192519351918251830194602086019360208188031261024757602081015167ffffffffffffffff811161024757019560a0908790031261024757604051916123ac83614186565b602087015183526123bf60408801615091565b91602084019283526123d360608901615091565b916040850192835260808901519860ff8a168a036102475760608601998a5260a081015167ffffffffffffffff811161024757602091010187601f8201121561024757805161242181614235565b9861242f6040519a8b614212565b818a52602082840101116102475761244d916020808b019101614122565b6080850196875277ffffffffffffffff0000000000000000000000000000000060405191632cbc26bb60e01b835260801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116e657600091612873575b5061179d576124d4818661472c565b1561283557509560ff6124fb6125369361252a989961ffff808851935116915116916153ee565b612525612510879a939a518587511690614d9d565b9961251e8587511684614d9d565b97516145cc565b6145cc565b915116855191886149ed565b9384600052600d6020526040600020916125978260405194612557866141a2565b549561256660ff881687614847565b6001600160a01b03602087019760081c16875288600052600d6020526040600020600260ff1982541617905561572d565b60009383516003811015612821576126cc575050600093516020818051810103126126c857602001516001600160a01b0381168091036126c8576040517f40c10f190000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602481019190915260208180604481010381876001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af180156126bd5761269e575b505b519060038210156121c15761269b6060927f33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c9460405193845260208401526040830190614492565ba3005b6126b69060203d6020116116df576116d08183614212565b5085612651565b6040513d86823e3d90fd5b8480fd5b939094508251600381101561280d576001036127e157506126f5846001600160a01b03926145cc565b92511661274760206127078686614a34565b6040517f40c10f19000000000000000000000000000000000000000000000000000000008152306004820152602481019190915291829081906044820190565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af180156116e6576127c2575b508280612791575b5050612653565b6127bb917f0000000000000000000000000000000000000000000000000000000000000000615218565b858261278a565b6127da9060203d6020116116df576116d08183614212565b5086612782565b80867fb196a44a0000000000000000000000000000000000000000000000000000000060249352600452fd5b602482634e487b7160e01b81526021600452fd5b602486634e487b7160e01b81526021600452fd5b6110ad906040519182917f24eb47e5000000000000000000000000000000000000000000000000000000008352602060048401526024830190614145565b61288c915060203d6020116116df576116d08183614212565b896124c5565b60408336031261024757602060409182516128ac816141a2565b6128b58661410e565b81528286013583820152815201920191612347565b7fd7f73334000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b34610247576020600319360112610247577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917460206001600160a01b0361293c6140f8565b612944614ea4565b16807fffffffffffffffffffffffff00000000000000000000000000000000000000006009541617600955604051908152a1005b34610247576000600319360112610247576000546001600160a01b03811633036129ff577fffffffffffffffffffffffff0000000000000000000000000000000000000000600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b34610247576020600319360112610247576020611fa26001600160a01b03612a4f6140f8565b16600052600c602052604060002054151590565b346102475760206003193601126102475767ffffffffffffffff612a8561409e565b606060c0604051612a95816141da565b600081526000602082015260006040820152600083820152600060808201528260a0820152015216600052600a60205260606040600020610d57612ad76157a1565b611a5d604051612ae6816141da565b84548152612b9360018601549563ffffffff602084019760ff81161515895261ffff60408601818360081c168152818c880191818560181c1683528560808a019560281c168552612b4c6003612b3e60028a016147a3565b9860a08c01998a52016147a3565b9860c08101998a526040519e8f9e8f9260408452516040840152511515910152511660808c0152511660a08a0152511660c08801525160e080880152610120870190614145565b9051603f1986830301610100870152614145565b346102475760006003193601126102475760206001600160a01b0360095416604051908152f35b3461024757600060031936011261024757610d576111926157a1565b3461024757612bf836614329565b612c03929192614ea4565b67ffffffffffffffff8216612c25816000526006602052604060002054151590565b15612c405750610c7692612c3a913691614251565b90614ee2565b7f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b3461024757612c95612c9d612c81366142d7565b9491612c8e939193614ea4565b36916143f2565b9236916143f2565b7f000000000000000000000000000000000000000000000000000000000000000015612da35760005b8251811015612d2c57806001600160a01b03612ce4600193866148a9565b5116612cef81615cb8565b612cfb575b5001612cc6565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a184612cf4565b5060005b8151811015610c7657806001600160a01b03612d4e600193856148a9565b51168015612d9d57612d5f81615891565b612d6c575b505b01612d30565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a183612d64565b50612d66565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b3461024757604060031936011261024757612de661409e565b60243567ffffffffffffffff811161024757602091612e0c611fa2923690600401614288565b9061472c565b346102475760206003193601126102475760043567ffffffffffffffff81116102475780600401906101006003198236030112610247576000604051612e57816141be565b52612e84612e7a612e75612e6e60c48501866146b2565b3691614251565b614cc2565b6064830135614d9d565b9060848101612e9281614703565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116036117e65750602481019277ffffffffffffffff00000000000000000000000000000000612eec85614717565b60801b1660405190632cbc26bb60e01b825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116e6576000916131c7575b5061179d5767ffffffffffffffff612f5b85614717565b16612f73816000526006602052604060002054151590565b156117705760206001600160a01b0360045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa9081156116e6576000916131a8575b50156116f257612fde84614717565b90612ff460a4840192612e0c612e6e85856146b2565b156131615750506044906130108361300b86614717565b61572d565b0161306a60208361302084614703565b60405193849283927f40c10f1900000000000000000000000000000000000000000000000000000000845260048401602090939291936001600160a01b0360408201951681520152565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af180156116e65767ffffffffffffffff6130e86130e26020977ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc095608095613144575b50614717565b94614703565b936001600160a01b0360405195817f000000000000000000000000000000000000000000000000000000000000000016875233898801521660408601528560608601521692a28060405161313b816141be565b52604051908152f35b61315a908a3d8c116116df576116d08183614212565b50896130dc565b61316b92506146b2565b6110ad6040519283927f24eb47e50000000000000000000000000000000000000000000000000000000084526020600485015260248401916145d9565b6131c1915060203d6020116116df576116d08183614212565b85612fcf565b6131e0915060203d6020116116df576116d08183614212565b85612f44565b346102475760a0600319360112610247576131ff61409e565b6024359060443567ffffffffffffffff8111610247576132239036906004016140ca565b91606435916001600160a01b0383168093036102475760843567ffffffffffffffff8111610247576132599036906004016140ca565b5050613263614699565b50604051936132718561416a565b6000855260208501926000845260408601926000845260608701600081526060608060405161329f81614186565b828152826020820152826040820152600083820152015267ffffffffffffffff831693604051632cbc26bb60e01b815277ffffffffffffffff000000000000000000000000000000008560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116e657600091613634575b5061179d5761334133615c43565b613358856000526006602052604060002054151590565b156136065784600052600a60205260406000209485548b116135d457509861346a6134789260ff6134fe9b9c63ffffffff60018a01546133ba61ffff8260081c169c8d966133af61ffff8560181c1680998c6153ee565b928382935252614a34565b8d5260281c1680613569575061ffff6134286133d860038c016147a3565b985b604051976133e789614186565b8852602088019c8d52604088019586526060880193857f00000000000000000000000000000000000000000000000000000000000000001685523691614251565b9360808701948552816040519c8d986020808b01525160408a01525116606088015251166080860152511660a08401525160a060c084015260e0830190614145565b03601f198101865285614212565b60209586946040519061348b8783614212565b600082526134a76002604051976134a189614186565b016147a3565b8652868601526040850152606084015260808301526001600160a01b0360045416906040518097819482937f20487ded000000000000000000000000000000000000000000000000000000008452600484016144ff565b03915afa9283156116e657600093613537575b50826040945251818451613524816141a2565b8481520190815283519283525190820152f35b9392508184813d8311613562575b61354f8183614212565b8101031261024757604093519293613511565b503d613545565b61342861ffff916040519061357d826141a2565b81526020810160018152604051917f181dcf100000000000000000000000000000000000000000000000000000000060208401525160248301525115156044820152604481526135ce606482614212565b986133da565b8a907f58dd87c50000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b847fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b61364d915060203d6020116116df576116d08183614212565b8b613333565b346102475760006003193601126102475760206103546145fa565b3461024757600060031936011261024757602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760206003193601126102475760206136c76140f8565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b3461024757600060031936011261024757610d57604051613765606082614212565b603581527f4275726e4d696e745769746845787465726e616c4d696e74657246617374547260208201527f616e73666572546f6b656e506f6f6c20312e362e3000000000000000000000006040820152604051918291602083526020830190614145565b60c0600319360112610247576137dd61409e565b60643567ffffffffffffffff8111610247576137fd9036906004016140ca565b608435916001600160a01b03831683036102475760a43567ffffffffffffffff8111610247576138319036906004016140ca565b5050604051916138408361416a565b600083526000602084015260006040840152600060608401526060608060405161386981614186565b8281528260208201528260408201526000838201520152604051632cbc26bb60e01b815277ffffffffffffffff000000000000000000000000000000008660801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116e657600091613f1c575b5061179d576138ff33615c43565b61392067ffffffffffffffff86166000526006602052604060002054151590565b15613ee45767ffffffffffffffff8516600052600a602052604060002093845460243511613ea65794849560016080960154613a5e61ffff8260081c1663ffffffff61ffff8460181c169361398e60408b61397e88876024356153ee565b9282846060849501520152614a34565b60208b015260281c1680613e4257506139a960038b016147a3565b925b604051916139b883614186565b60243583526020830152604082015260ff7f00000000000000000000000000000000000000000000000000000000000000001660608201526139fb368789614251565b89820152613a50604051998a926020808501528051604085015261ffff602082015116606085015261ffff6040820151168285015260ff60608201511660a0850152015160a060c084015260e0830190614145565b03601f198101895288614212565b604051602097613a6e8983614212565b60008252613a8460026040519b6134a18d614186565b8a52888a015260408901526001600160a01b038216606089015260808801526001600160a01b03600454168660405180927f20487ded0000000000000000000000000000000000000000000000000000000082528180613ae88d89600484016144ff565b03915afa9081156116e657600091613e15575b508552613b0a60243583614a41565b6020850151966044358811613de157866000969798507f0000000000000000000000000000000000000000000000000000000000000000613b4f602435303384614ad6565b613b876001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016918260243591614b40565b6024604051809981937f42966c68000000000000000000000000000000000000000000000000000000008352833560048401525af19182156116e657613c20968993613dc4575b506001600160a01b038116613d71575b506001600160a01b036004541660405180809881947f96f4e9f900000000000000000000000000000000000000000000000000000000835287600484016144ff565b039134905af19384156116e657600094613d22575b50847f662a290835d430973477690029020949d2975f8badb76f0593a6d573b08ef8f391613d17613ca4613c706020899a01516024356145cc565b613c7b36888a614251565b907f0000000000000000000000000000000000000000000000000000000000000000908a6149ed565b9567ffffffffffffffff613cbe60208601516024356145cc565b936060604087015196015160405196879687528d60ff7f000000000000000000000000000000000000000000000000000000000000000016908801526040870152606086015260a06080860152169560a08401916145d9565b0390a4604051908152f35b9093508581813d8311613d6a575b613d3a8183614212565b810103126102475751927f662a290835d430973477690029020949d2975f8badb76f0593a6d573b08ef8f3613c35565b503d613d30565b613dbe90613d8b895130336001600160a01b038516614ad6565b8851906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000009116614b40565b88613bde565b613dda90843d86116116df576116d08183614212565b5089613bce565b877f61acdb930000000000000000000000000000000000000000000000000000000060005260045260443560245260446000fd5b90508681813d8311613e3b575b613e2c8183614212565b81010312610247575188613afb565b503d613e22565b60405190613e4f826141a2565b81526020810160018152604051917f181dcf10000000000000000000000000000000000000000000000000000000006020840152516024830152511515604482015260448152613ea0606482614212565b926139ab565b67ffffffffffffffff867f58dd87c5000000000000000000000000000000000000000000000000000000006000521660045260243560245260446000fd5b67ffffffffffffffff857fa9902c7e000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b613f35915060203d6020116116df576116d08183614212565b866138f1565b3461024757602060031936011261024757600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361024757817ff6f46ff90000000000000000000000000000000000000000000000000000000060209314908115614013575b8115613fb6575b5015158152f35b7f85572ffb00000000000000000000000000000000000000000000000000000000811491508115613fe9575b5083613faf565b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483613fe2565b90507faff2afbf0000000000000000000000000000000000000000000000000000000081148015614075575b801561404c575b90613fa8565b507f01ffc9a7000000000000000000000000000000000000000000000000000000008114614046565b507f0e64dd2900000000000000000000000000000000000000000000000000000000811461403f565b6004359067ffffffffffffffff8216820361024757565b359067ffffffffffffffff8216820361024757565b9181601f840112156102475782359167ffffffffffffffff8311610247576020838186019501011161024757565b600435906001600160a01b038216820361024757565b35906001600160a01b038216820361024757565b60005b8381106141355750506000910152565b8181015183820152602001614125565b90601f19601f60209361416381518092818752878088019101614122565b0116010190565b6080810190811067ffffffffffffffff82111761088457604052565b60a0810190811067ffffffffffffffff82111761088457604052565b6040810190811067ffffffffffffffff82111761088457604052565b6020810190811067ffffffffffffffff82111761088457604052565b60e0810190811067ffffffffffffffff82111761088457604052565b6060810190811067ffffffffffffffff82111761088457604052565b90601f601f19910116810190811067ffffffffffffffff82111761088457604052565b67ffffffffffffffff811161088457601f01601f191660200190565b92919261425d82614235565b9161426b6040519384614212565b829481845281830111610247578281602093846000960137010152565b9080601f83011215610247578160206142a393359101614251565b90565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501948460051b01011161024757565b60406003198201126102475760043567ffffffffffffffff81116102475781614302916004016142a6565b929092916024359067ffffffffffffffff821161024757614325916004016142a6565b9091565b9060406003198301126102475760043567ffffffffffffffff8116810361024757916024359067ffffffffffffffff821161024757614325916004016140ca565b906020808351928381520192019060005b8181106143885750505090565b82516001600160a01b031684526020938401939092019160010161437b565b6020600319820112610247576004359067ffffffffffffffff8211610247576003198260a0920301126102475760040190565b67ffffffffffffffff81116108845760051b60200190565b9291906143fe816143da565b9361440c6040519586614212565b602085838152019160051b810192831161024757905b82821061442e57505050565b6020809161443b8461410e565b815201910190614422565b9080601f83011215610247578160206142a3933591016143f2565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501946060850201011161024757565b9060038210156121c15752565b35906001600160801b038216820361024757565b9190826060910312610247576040516144cb816141f6565b809280359081151582036102475760406144fa91819385526144ef6020820161449f565b60208601520161449f565b910152565b9067ffffffffffffffff9093929316815260406020820152614546614530845160a0604085015260e0840190614145565b6020850151603f19848303016060850152614145565b90604084015191603f198282030160808301526020808451928381520193019060005b8181106145a1575050506080846001600160a01b0360606142a3969701511660a084015201519060c0603f1982850301910152614145565b825180516001600160a01b031686526020908101518187015260409095019490920191600101614569565b91908203918211611f1057565b601f8260209493601f19938186528686013760008582860101520116010190565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156116e65760009161466a575090565b90506020813d602011614691575b8161468560209383614212565b81010312610247575190565b3d9150614678565b604051906146a6826141a2565b60006020838281520152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610247570180359067ffffffffffffffff82116102475760200191813603831361024757565b356001600160a01b03811681036102475790565b3567ffffffffffffffff811681036102475790565b9067ffffffffffffffff6142a392166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b90600182811c92168015614799575b602083101461478357565b634e487b7160e01b600052602260045260246000fd5b91607f1691614778565b90604051918260008254926147b784614769565b808452936001811690811561482557506001146147de575b506147dc92500383614212565b565b90506000929192526020600020906000915b8183106148095750509060206147dc92820101386147cf565b60209193508060019154838589010152019101909184926147f0565b602093506147dc95925060ff1991501682840152151560051b820101386147cf565b60038210156121c15752565b91908110156148935760051b810135907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee181360301821215610247570190565b634e487b7160e01b600052603260045260246000fd5b80518210156148935760209160051b010190565b91908110156148935760051b0190565b9190811015614893576060020190565b604051906148ea82614186565b60006080838281528260208201528260408201528260608201520152565b9060405161491581614186565b60806001829460ff81546001600160801b038116865263ffffffff81861c16602087015260a01c161515604085015201546001600160801b0381166060840152811c910152565b67ffffffffffffffff1660005260076020526142a360046040600020016147a3565b818110614989575050565b6000815560010161497e565b81810292918115918404141715611f1057565b9190601f81116149b757505050565b6147dc926000526020600020906020601f840160051c830193106149e3575b601f0160051c019061497e565b90915081906149d6565b9290614a20614a2e9260ff60405195869460208601988952604086015216606084015260808084015260a0830190614145565b03601f198101835282614212565b51902090565b91908201809211611f1057565b67ffffffffffffffff7fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da817894491169182600052600760205280614ab260406000206001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615415565b604080516001600160a01b039092168252602082019290925290819081015b0390a2565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000060208201526001600160a01b03928316602482015292909116604483015260648201929092526147dc91614b3b82608481015b03601f198101845283614212565b6155fa565b91909181158015614c28575b15614ba4576040517f095ea7b30000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909316602484015260448301919091526147dc9190614b3b8260648101614b2d565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b0384166024820152602081806044810103816001600160a01b0386165afa9081156116e657600091614c90575b5015614b4c565b90506020813d602011614cba575b81614cab60209383614212565b81010312610247575138614c89565b3d9150614c9e565b80518015614d3257602003614cf457805160208281019183018390031261024757519060ff8211614cf4575060ff1690565b6110ad906040519182917f953576f7000000000000000000000000000000000000000000000000000000008352602060048401526024830190614145565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff8211611f1057565b60ff16604d8111611f1057600a0a90565b8115614d87570490565b634e487b7160e01b600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff811692828414614e8557828411614e5b5790614de291614d58565b91604d60ff8416118015614e40575b614e0a57505090614e046142a392614d6c565b90614995565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b50614e4a83614d6c565b8015614d8757600019048411614df1565b614e6491614d58565b91604d60ff841611614e0a57505090614e7f6142a392614d6c565b90614d7d565b5050505090565b90816020910312610247575180151581036102475790565b6001600160a01b03600154163303614eb857565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b908051156108d35767ffffffffffffffff81516020830120921691826000526007602052614f17816005604060002001615942565b1561504d5760005260086020526040600020815167ffffffffffffffff811161088457614f4e81614f488454614769565b846149a8565b6020601f8211600114614fc35791614fa2827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea9593614ad195600091614fb8575b506000198260011b9260031b1c19161790565b9055604051918291602083526020830190614145565b905084015138614f8f565b601f1982169083600052806000209160005b818110615035575092614ad19492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea98961061501c575b5050811b019055610f35565b85015160001960f88460031b161c191690553880615010565b9192602060018192868a015181550194019201614fd5565b50906110ad6040519283927f393b8ad20000000000000000000000000000000000000000000000000000000084526004840152604060248401526044830190614145565b519061ffff8216820361024757565b357fffffffff00000000000000000000000000000000000000000000000000000000811681036102475790565b3563ffffffff811681036102475790565b3561ffff811681036102475790565b3580151581036102475790565b67ffffffffffffffff1660008181526006602052604090205490929190156151ea57916151e760e0926151bc856151517f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b976152dd565b846000526007602052615168816040600020615a52565b615171836152dd565b84600052600760205261518b836002604060002001615a52565b60405194855260208501906001600160801b0360408092805115158552826020820151166020860152015116910152565b60808301906001600160801b0360408092805115158552826020820151166020860152015116910152565ba1565b827f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909216602483015260448201929092526147dc91614b3b8260648101614b2d565b6152726148dd565b506001600160801b036060820151166001600160801b0380835116916152bd60208501936152b76152aa63ffffffff875116426145cc565b8560808901511690614995565b90614a34565b808210156152d657505b16825263ffffffff4216905290565b90506152c7565b805115615362576001600160801b036040820151166001600160801b03602083015116106153085750565b606490615360604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565bfd5b6001600160801b03604082015116158015906153d8575b6153805750565b606490615360604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906001600160801b0360408092805115158552826020820151166020860152015116910152565b506001600160801b036020820151161515615379565b6154119061ffff6127106154088282969897981684614995565b04951690614995565b0490565b9182549060ff8260a01c161580156155f2575b6155ec576001600160801b038216916001850190815461545b63ffffffff6001600160801b0383169360801c16426145cc565b908161554e575b505084811061550f57508383106154a45750506154886001600160801b039283926145cc565b16166fffffffffffffffffffffffffffffffff19825416179055565b5460801c916154b381856145cc565b92600019810190808211611f10576154d66154db926001600160a01b0396614a34565b614d7d565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b82856001600160a01b03927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b8286929396116155c257615569926152b79160801c90614995565b808410156155bd5750825b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff0000000000000000000000000000000016178655923880615462565b615574565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b508215615428565b6001600160a01b0361567c91169160409260008085519361561b8786614212565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d15615725573d9161566083614235565b9261566d87519485614212565b83523d6000602085013e615e80565b8051908161568957505050565b60208061569a938301019101614e8c565b156156a25750565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606091615e80565b67ffffffffffffffff7f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c91169182600052600760205280614ab260026040600020016001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615415565b60405190600b548083528260208101600b60005260206000209260005b8181106157d35750506147dc92500383614212565b84548352600194850194879450602090930192016157be565b906040519182815491828252602082019060005260206000209260005b81811061581e5750506147dc92500383614212565b8454835260019485019487945060209093019201615809565b80548210156148935760005260206000200190600090565b80549068010000000000000000821015610884578161587691600161588d94018155615837565b81939154906000199060031b92831b921b19161790565b9055565b806000526003602052604060002054156000146158ca576158b381600261584f565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146158ca576158f281600b61584f565b600b5490600052600c602052604060002055600190565b806000526006602052604060002054156000146158ca5761592b81600561584f565b600554906000526006602052604060002055600190565b600082815260018201602052604090205461597957806159648360019361584f565b80549260005201602052604060002055600190565b5050600090565b805480156159a85760001901906159978282615837565b60001982549160031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c60205260409020548015615979576000198101818111611f1057600b54906000198201918211611f1057808203615a18575b505050615a04600b615980565b600052600c60205260006040812055600190565b615a3a615a2961587693600b615837565b90549060031b1c928392600b615837565b9055600052600c6020526040600020553880806159f7565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1991615b5b6060928054615a8f63ffffffff8260801c16426145cc565b9081615b91575b50506001600160801b036001816020860151169282815416808510600014615b8957508280855b16166fffffffffffffffffffffffffffffffff19825416178155615b278651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601516fffffffffffffffffffffffffffffffff1960809190911b16939092166001600160801b031692909217910155565b6151e760405180926001600160801b0360408092805115158552826020820151166020860152015116910152565b838091615abd565b6001600160801b0391615bbd839283615bb66001880154948286169560801c90614995565b9116614a34565b80821015615c3c57505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9290911692909216167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116174260801b73ffffffff00000000000000000000000000000000161781553880615a96565b9050615bc7565b7f0000000000000000000000000000000000000000000000000000000000000000615c6b5750565b6001600160a01b031680600052600360205260406000205415615c8b5750565b7fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6000818152600360205260409020548015615979576000198101818111611f1057600254906000198201918211611f1057818103615d12575b505050615cfe6002615980565b600052600360205260006040812055600190565b615d34615d23615876936002615837565b90549060031b1c9283926002615837565b90556000526003602052604060002055388080615cf1565b6000818152600660205260409020548015615979576000198101818111611f1057600554906000198201918211611f1057818103615da6575b505050615d926005615980565b600052600660205260006040812055600190565b615dc8615db7615876936005615837565b90549060031b1c9283926005615837565b90556000526006602052604060002055388080615d85565b906001820191816000528260205260406000205490811515600014615e7757600019820191808311611f105781546000198101908111611f10578381615e2e9503615e40575b505050615980565b60005260205260006040812055600190565b615e60615e506158769386615837565b90549060031b1c92839286615837565b905560005284602052604060002055388080615e26565b50505050600090565b91929015615efb5750815115615e94575090565b3b15615e9d5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015615f0e5750805190602001fd5b6110ad906040519182917f08c379a000000000000000000000000000000000000000000000000000000000835260206004840152602483019061414556fea164736f6c634300081a000a",
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
	return common.HexToHash("0x662a290835d430973477690029020949d2975f8badb76f0593a6d573b08ef8f3")
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
