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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSendToken\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"computeFillId\",\"inputs\":[{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"fastFill\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAccumulatedPoolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedFillers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCcipSendTokenFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"settlementFeeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIFastTransferPool.Quote\",\"components\":[{\"name\":\"ccipSettlementFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfig\",\"components\":[{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFillInfo\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFastTransferTokenPoolAbstract.FillInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIFastTransferPool.FillState\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAllowedFiller\",\"inputs\":[{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateDestChainConfig\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFastTransferTokenPoolAbstract.DestChainConfigUpdateArgs[]\",\"components\":[{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"customExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFillerAllowList\",\"inputs\":[{\"name\":\"fillersToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"fillersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPoolFees\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fastTransferFillerFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"fastTransferPoolFeeBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"maxFillAmountPerRequest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationPool\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"settlementOverheadGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fillerAllowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestinationPoolUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destinationPool\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"filler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferRequested\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceAmountNetFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sourceDecimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"fastTransferFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FastTransferSettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"settlementId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fillerReimbursementAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"poolFeeAccumulated\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"prevState\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIFastTransferPool.FillState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FillerAllowListUpdated\",\"inputs\":[{\"name\":\"addFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"removeFillers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolFeeWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyFilled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadySettled\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"FillerNotAllowlisted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"filler\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientPoolFees\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFillId\",\"inputs\":[{\"name\":\"fillId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"QuoteFeeExceedsUserMaxLimit\",\"inputs\":[{\"name\":\"quoteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFastTransferFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferAmountExceedsMaxFillAmount\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x61014080604052346104405761666a803803809161001d8285610505565b8339810160a0828203126104405761003482610528565b906100416020840161053c565b60408401516001600160401b0381116104405784019180601f84011215610440578251926001600160401b0384116104ef578360051b90602082019461008a6040519687610505565b855260208086019282010192831161044057602001905b8282106104d7575050506100c360806100bc60608701610528565b9501610528565b6040516321df0da760e01b81526001600160a01b03909416949093602081600481895afa9081156104cb57600091610491575b506001600160a01b031690331561048057600180546001600160a01b031916331790558115801561046f575b801561045e575b61044d578160209160049360805260c0526040519283809263313ce56760e01b82525afa6000918161040c575b506103e1575b5060a052600480546001600160a01b0319166001600160a01b0384169081179091558151151560e08190529091906102be575b50156102a8576101005261012052604051615f7f90816106eb82396080518181816113f701528181611461015281816116360152818161219e0152818161270e01528181612e230152818161300d015281816135c10152818161360e01528181613a1301528181614432015281816148b601528181614fb601526151a1015260a05181818161167a0152818161330201528181613577015281816138b401528181613afe01528181614b620152614bcc015260c051818181610ba6015281816114ef015281816124f801528181612eb20152818161320101526137bc015260e051818181610b6101528181612c260152615c6b01526101005181613bbd01526101205181818161022301528181614d2401526151790152f35b6335fdcccd60e21b600052600060045260246000fd5b602091604051916102cf8484610505565b60008352600036813760e051156103d05760005b835181101561034a576001906001600160a01b03610301828761054a565b51168661030d8261058c565b61031a575b5050016102e3565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a13886610312565b5091509260005b82518110156103c5576001906001600160a01b0361036f828661054a565b511680156103bf57856103818261068a565b61038f575b50505b01610351565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a13885610386565b50610389565b50929150503861018f565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036103f5575061015c565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d602011610445575b8161042860209383610505565b81010312610440576104399061053c565b9038610156565b600080fd5b3d915061041b565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03811615610129565b506001600160a01b03851615610122565b639b15e16f60e01b60005260046000fd5b90506020813d6020116104c3575b816104ac60209383610505565b81010312610440576104bd90610528565b386100f6565b3d915061049f565b6040513d6000823e3d90fd5b602080916104e484610528565b8152019101906100a1565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b038211908210176104ef57604052565b51906001600160a01b038216820361044057565b519060ff8216820361044057565b805182101561055e5760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b805482101561055e5760005260206000200190600090565b600081815260036020526040902054801561068357600019810181811161066d5760025460001981019190821161066d5781810361061c575b505050600254801561060657600019016105e0816002610574565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b61065561062d61063e936002610574565b90549060031b1c9283926002610574565b819391549060031b91821b91600019901b19161790565b905560005260036020526040600020553880806105c5565b634e487b7160e01b600052601160045260246000fd5b5050600090565b806000526003602052604060002054156000146106e457600254680100000000000000008110156104ef576106cb61063e8260018594016002556002610574565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a714613d4a57508063055befd4146136b8578063181f5a771461363257806321df0da7146135ee578063240028e81461359b57806324f65ee71461355d5780632b2c0eb4146135425780632e7aa8c8146130de5780633907753714612d995780634c5ef0ed14612d5457806354c8a4f314612bf457806362ddd3c414612b715780636609f59914612b555780636d3d1a5814612b2e5780636def4ce7146129ea57806378b410f2146129b057806379ba5097146128ff5780637d54534e1461287f57806385572ffb1461229557806387f060d014611fc95780638926f54f14611f845780638a18dcbd14611aaf5780638da5cb5b14611a88578063929ea5ba1461197e578063962d4020146118425780639a4575b9146114265780639fe280f514611393578063a42a7b8b14611261578063a7cd63b7146111f3578063abe1c1e814611184578063acfecf911461105f578063af58d59f14611015578063b0f479a114610fee578063b794658014610fb6578063c0d7865514610f12578063c4bffe2b14610e00578063c75eea9c14610d57578063cf7401f314610bca578063dc0bd97114610b86578063e0351e1314610b49578063e8a1da171461035c578063eeebc67414610304578063f2fde38b1461024c5763f36675171461020357600080fd5b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b34610247576020600319360112610247576001600160a01b0361026d613f07565b610275614d6e565b163381146102da57807fffffffffffffffffffffffff000000000000000000000000000000000000000060005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102475760806003193601126102475760443560ff811681036102475760643567ffffffffffffffff81116102475760209161034861035492369060040161407b565b90602435600435614819565b604051908152f35b346102475761036a366140ca565b919092610375614d6e565b6000905b8282106109a45750505060009063ffffffff42165b81831061039757005b6103a283838661466d565b926101208436031261024757604051936103bb85613f31565b6103c481613ec4565b8552602081013567ffffffffffffffff81116102475781019336601f860112156102475784356103f3816141cd565b956104016040519788613fbd565b81875260208088019260051b820101903682116102475760208101925b828410610975575050505060208601948552604082013567ffffffffffffffff811161024757610451903690840161407b565b906040870191825261047b61046936606086016142af565b936060890194855260c03691016142af565b946080880195865261048d84516152d5565b61049786516152d5565b8251511561094b576104b367ffffffffffffffff8951166158ed565b156109125767ffffffffffffffff885116600052600760205260406000206105f585516fffffffffffffffffffffffffffffffff604082015116906105b06fffffffffffffffffffffffffffffffff6020830151169151151583608060405161051b81613f31565b858152602081018a905260408101849052606081018690520152855474ff000000000000000000000000000000000000000091151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff84161773ffffffff00000000000000000000000000000000608089901b1617178555565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff91909116176001830155565b61071b87516fffffffffffffffffffffffffffffffff604082015116906106d66fffffffffffffffffffffffffffffffff6020830151169151151583608060405161063f81613f31565b858152602081018a9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166fffffffffffffffffffffffffffffffff85161773ffffffff0000000000000000000000000000000060808a901b161791151560a01b74ff000000000000000000000000000000000000000016919091179055565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff91909116176003830155565b6004845191019080519067ffffffffffffffff82116108fc57610748826107428554614565565b856147d4565b602090601f83116001146108955761077892916000918361088a575b50506000198260011b9260031b1c19161790565b90555b60005b875180518210156107b357906107ad6001926107a68367ffffffffffffffff8e5116926146c3565b5190614dac565b0161077e565b505097967f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c293929196509461087f67ffffffffffffffff600197511692519351915161084b6108166040519687968752610100602088015261010087019061401f565b9360408601906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b60a08401906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b0390a101919261038e565b015190508d80610764565b90601f1983169184600052816000209260005b8181106108e457509084600195949392106108cb575b505050811b01905561077b565b015160001960f88460031b161c191690558c80806108be565b929360206001819287860151815501950193016108a8565b634e487b7160e01b600052604160045260246000fd5b67ffffffffffffffff8851167f1d5ad3c50000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b833567ffffffffffffffff811161024757602091610999839283369187010161407b565b81520193019261041e565b9092919367ffffffffffffffff6109c46109bf8688866146d7565b614513565b16926109cf84615d72565b15610b1b578360005260076020526109ed60056040600020016157d0565b9260005b8451811015610a2957600190866000526007602052610a226005604060002001610a1b83896146c3565b5190615e06565b50016109f1565b5093909491959250806000526007602052600560406000206000815560006001820155600060028201556000600382015560048101610a688154614565565b9081610ad8575b5050018054906000815581610ab7575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a1019091610379565b6000526020600020908101905b81811015610a7f5760008155600101610ac4565b81601f60009311600114610af05750555b8880610a6f565b81835260208320610b0b91601f01861c8101906001016147aa565b8082528160208120915555610ae9565b837f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102475760006003193601126102475760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760e060031936011261024757610be3613ead565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc36011261024757604051610c1981613fa1565b60243580151581036102475781526044356fffffffffffffffffffffffffffffffff811681036102475760208201526064356fffffffffffffffffffffffffffffffff8116810361024757604082015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7c3601126102475760405190610ca082613fa1565b608435801515810361024757825260a4356fffffffffffffffffffffffffffffffff8116810361024757602083015260c4356fffffffffffffffffffffffffffffffff811681036102475760408301526001600160a01b036009541633141580610d42575b610d1457610d1292615038565b005b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b506001600160a01b0360015416331415610d05565b346102475760206003193601126102475767ffffffffffffffff610d79613ead565b610d816146f7565b50166000526007602052610dfc610da3610d9e6040600020614722565b615256565b6040519182918291909160806fffffffffffffffffffffffffffffffff8160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b34610247576000600319360112610247576040516005548082528160208101600560005260206000209260005b818110610ef9575050610e4292500382613fbd565b805190610e67610e51836141cd565b92610e5f6040519485613fbd565b8084526141cd565b90601f1960208401920136833760005b8151811015610ea9578067ffffffffffffffff610e96600193856146c3565b5116610ea282876146c3565b5201610e77565b5050906040519182916020830190602084525180915260408301919060005b818110610ed6575050500390f35b825167ffffffffffffffff16845285945060209384019390920191600101610ec8565b8454835260019485019486945060209093019201610e2d565b3461024757602060031936011261024757610f2b613f07565b610f33614d6e565b6001600160a01b03811690811561094b57600480547fffffffffffffffffffffffff000000000000000000000000000000000000000081169093179055604080516001600160a01b0393841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491819081015b0390a1005b3461024757602060031936011261024757610dfc610fda610fd5613ead565b614788565b60405191829160208352602083019061401f565b346102475760006003193601126102475760206001600160a01b0360045416604051908152f35b346102475760206003193601126102475767ffffffffffffffff611037613ead565b61103f6146f7565b50166000526007602052610dfc610da3610d9e6002604060002001614722565b346102475767ffffffffffffffff6110763661411c565b929091611081614d6e565b169061109a826000526006602052604060002054151590565b15611156578160005260076020526110cb60056040600020016110be368685614044565b6020815191012090615e06565b1561110f577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76919261110a6040519283926020845260208401916143d5565b0390a2005b611152906040519384937f74f23c7c00000000000000000000000000000000000000000000000000000000855260048501526040602485015260448401916143d5565b0390fd5b507f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102475760206003193601126102475761119d614495565b50600435600052600d6020526040806000206001600160a01b038251916111c383613f4d565b546111d160ff821684614661565b81602084019160081c1681526111ea8451809451614285565b51166020820152f35b34610247576000600319360112610247576040516002548082526020820190600260005260206000209060005b81811061124b57610dfc8561123781870382613fbd565b60405191829160208352602083019061415d565b8254845260209093019260019283019201611220565b346102475760206003193601126102475767ffffffffffffffff611283613ead565b16600052600760205261129c60056040600020016157d0565b805190601f196112c46112ae846141cd565b936112bc6040519586613fbd565b8085526141cd565b0160005b81811061138257505060005b815181101561131c57806112ea600192846146c3565b516000526008602052611300604060002061459f565b61130a82866146c3565b5261131581856146c3565b50016112d4565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b82821061135557505050500390f35b9193602061137282603f196001959799849503018652885161401f565b9601920192018594939192611346565b8060606020809387010152016112c8565b34610247576020600319360112610247576113ac613f07565b6113b4614d6e565b6113bc6143f6565b90816113c457005b60206001600160a01b038261141b857f738b39462909f2593b7546a62adee9bc4e5cadde8e0e0f80686198081b859599957f0000000000000000000000000000000000000000000000000000000000000000615204565b6040519485521692a2005b34610247576114343661419a565b6060602060405161144481613f4d565b828152015260808101611456816144ff565b6001600160a01b03807f00000000000000000000000000000000000000000000000000000000000000001691160361180357506020810177ffffffffffffffff000000000000000000000000000000006114af82614513565b60801b16604051907f2cbc26bb00000000000000000000000000000000000000000000000000000000825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611771576000916117d4575b506117aa57611539611534604084016144ff565b615c69565b67ffffffffffffffff61154b82614513565b16611563816000526006602052604060002054151590565b1561177d5760206001600160a01b0360045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa90811561177157600091611721575b506001600160a01b031633036116f357610fd5816116e0936115ec60606115e261167096614513565b920135809261486d565b6115f581615168565b7ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae1067ffffffffffffffff61162884614513565b604080516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152336020820152908101949094521691606090a2614513565b610dfc60405160ff7f0000000000000000000000000000000000000000000000000000000000000000166020820152602081526116ae604082613fbd565b604051926116bb84613f4d565b835260208301908152604051938493602085525160406020860152606085019061401f565b9051601f1984830301604085015261401f565b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d602011611769575b8161173a60209383613fbd565b810103126117655751906001600160a01b038216820361176257506001600160a01b036115b9565b80fd5b5080fd5b3d915061172d565b6040513d6000823e3d90fd5b7fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b6117f6915060203d6020116117fc575b6117ee8183613fbd565b810190614cb8565b83611520565b503d6117e4565b6118146001600160a01b03916144ff565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b346102475760606003193601126102475760043567ffffffffffffffff811161024757611873903690600401614099565b9060243567ffffffffffffffff811161024757611894903690600401614254565b9060443567ffffffffffffffff8111610247576118b5903690600401614254565b6001600160a01b036009541633141580611969575b610d145783861480159061195f575b6119355760005b8681106118e957005b8061192f6118fd6109bf6001948b8b6146d7565b6119088389896146e7565b61192961192161191986898b6146e7565b9236906142af565b9136906142af565b91615038565b016118e0565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b50808614156118d9565b506001600160a01b03600154163314156118ca565b346102475760406003193601126102475760043567ffffffffffffffff8111610247576119af903690600401614239565b60243567ffffffffffffffff8111610247576119cf903690600401614239565b906119d8614d6e565b60005b8151811015611a0a5780611a036001600160a01b036119fc600194866146c3565b51166158b4565b50016119db565b5060005b8251811015611a3d5780611a366001600160a01b03611a2f600194876146c3565b51166159a2565b5001611a0e565b7ffd35c599d42a981cbb1bbf7d3e6d9855a59f5c994ec6b427118ee0c260e24193611a7a83610fb18660405193849360408552604085019061415d565b90838203602085015261415d565b346102475760006003193601126102475760206001600160a01b0360015416604051908152f35b346102475760206003193601126102475760043567ffffffffffffffff811161024757611ae0903690600401614099565b611ae8614d6e565b60005b818110611af457005b611aff81838561466d565b60a081017f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000611b4e83614fde565b1614611f43575b60208201611b628161501c565b90604084019161ffff80611b758561501c565b1691160161ffff8111611f2d5761ffff61271091161015611f03576080840167ffffffffffffffff611ba682614513565b16600052600a60205260406000209460e0810194611bc486836144ae565b600289019167ffffffffffffffff82116108fc57611be6826107428554614565565b600090601f8311600114611e9f57611c15929160009183611e945750506000198260011b9260031b1c19161790565b90555b611c218461501c565b926001880197885498611c338861501c565b60181b64ffff0000001695611c478661502b565b151560c087013597888555606088019c611c608e61500b565b60281b68ffffffff0000000000169360081b62ffff0016907fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016177fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff16179060ff16171790556101008401611cd590856144ae565b90916003019167ffffffffffffffff82116108fc57611cf8826107428554614565565b600090601f8311600114611e1f579180611d2d92611d34969594600092611e145750506000198260011b9260031b1c19161790565b9055614513565b93611d3e9061501c565b94611d489061501c565b95611d5390836144ae565b9091611d5e90614fde565b97611d689061500b565b92611d729061502b565b936040519761ffff899816885261ffff16602088015260408701526060860160e0905260e0860190611da3926143d5565b957fffffffff0000000000000000000000000000000000000000000000000000000016608085015263ffffffff1660a0840152151560c083015267ffffffffffffffff1692037f6cfec31453105612e33aed8011f0e249b68d55e4efa65374322eb7ceeee76fbd91a2600101611aeb565b013590503880610764565b838252602082209a9e9d9c9b9a91601f198416815b818110611e7c5750919e9f9b9c9d9e6001939185611d349897969410611e62575b505050811b019055614513565b60001960f88560031b161c199101351690558f8080611e55565b91936020600181928787013581550195019201611e34565b013590508e80610764565b8382526020822091601f198416815b818110611eeb5750908460019594939210611ed1575b505050811b019055611c18565b60001960f88560031b161c199101351690558d8080611ec4565b83830135855560019094019360209283019201611eae565b7f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff611f546060840161500b565b1615611b55577f382c09820000000000000000000000000000000000000000000000000000000060005260046000fd5b34610247576020600319360112610247576020611fbf67ffffffffffffffff611fab613ead565b166000526006602052604060002054151590565b6040519015158152f35b346102475760c06003193601126102475760043560243560443567ffffffffffffffff81169182820361024757606435926084359160ff831683036102475760a435926001600160a01b038416928385036102475780600052600a60205260ff60016040600020015416612249575b5061205c60405184602082015260208152612054604082613fbd565b828885614819565b870361221b5786600052600d60205260406000206001600160a01b036040519161208583613f4d565b5461209360ff821684614661565b60081c16602082015251956003871015612205576000966121d9576120c2916120bb91614bc9565b8095614f6a565b604051956120cf87613f4d565b600187526020870196338852818752600d6020526040872090519760038910156121c55787986121c2985060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008454169116178255517fffffffffffffffffffffff0000000000000000000000000000000000000000ff74ffffffffffffffffffffffffffffffffffffffff0083549260081b1691161790556040519285845260208401527fd6f70fb263bfe7d01ec6802b3c07b6bd32579760fe9fcb4e248a036debb8cdf160403394a4337f0000000000000000000000000000000000000000000000000000000000000000614902565b80f35b602488634e487b7160e01b81526021600452fd5b602487897fcee81443000000000000000000000000000000000000000000000000000000008252600452fd5b634e487b7160e01b600052602160045260246000fd5b867fcb537aa40000000000000000000000000000000000000000000000000000000060005260045260246000fd5b61226033600052600c602052604060002054151590565b612038577f6c46a9b5000000000000000000000000000000000000000000000000000000006000526004523360245260446000fd5b34610247576122a33661419a565b6001600160a01b036004541633036128515760a081360312610247576040516122cb81613f31565b813581526122db60208301613ec4565b9060208101918252604083013567ffffffffffffffff811161024757612304903690850161407b565b9060408101918252606084013567ffffffffffffffff81116102475761232d903690860161407b565b936060820194855260808101359067ffffffffffffffff8211610247570136601f82011215610247578035612361816141cd565b9161236f6040519384613fbd565b81835260208084019260061b8201019036821161024757602001915b818310612819575050506080820152600092519067ffffffffffffffff8216905192519451805181019560208701916020818903126128155760208101519067ffffffffffffffff82116126b757019660a0908890031261281157604051936123f385613f31565b6020880151855261240660408901614f5b565b916020860192835261241a60608a01614f5b565b916040870192835260808a01519960ff8b168b0361280d57606088019a8b5260a08101519067ffffffffffffffff82116128095790602091010185601f8201121561280d5780519061246b82613fe0565b966124796040519889613fbd565b8288526020838301011161280957906124989160208089019101613ffc565b6080870194855277ffffffffffffffff00000000000000000000000000000000604051917f2cbc26bb00000000000000000000000000000000000000000000000000000000835260801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156127fe5789916127df575b506127b7576125378186614528565b1561277957509660ff61255f61259a9361258e98999a61ffff808a519351169151169161541c565b612589612574899a939a518587511690614bc9565b996125828587511684614bc9565b99516143c8565b6143c8565b91511683519188614819565b93848752600d6020526040872091604051926125b584613f4d565b546125c360ff821685614661565b6001600160a01b03602085019160081c168152889584516003811015612765576126bb5750506125f4828995614f6a565b516020818051810103126126b75760200151906001600160a01b0382168092036126b75761262191614cd0565b838652600d6020526040862060027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055519060038210156126a3579161269f6060927f33e17439bb4d31426d9168fc32af3a69cfce0467ba0d532fa804c27b5ff2189c9460405193845260208401526040830190614285565ba380f35b602486634e487b7160e01b81526021600452fd5b8780fd5b9450945050815160038110156121c557600103612739576126e4836001600160a01b03926143c8565b9351166126fa6126f48486614860565b30614cd0565b8380612708575b5050612621565b612732917f0000000000000000000000000000000000000000000000000000000000000000615204565b8683612701565b602487867fb196a44a000000000000000000000000000000000000000000000000000000008252600452fd5b60248b634e487b7160e01b81526021600452fd5b611152906040519182917f24eb47e500000000000000000000000000000000000000000000000000000000835260206004840152602483019061401f565b6004887f53ad11d8000000000000000000000000000000000000000000000000000000008152fd5b6127f8915060203d6020116117fc576117ee8183613fbd565b8a612528565b6040513d8b823e3d90fd5b8a80fd5b8980fd5b8580fd5b8680fd5b604083360312610247576020604091825161283381613f4d565b61283c86613f1d565b8152828601358382015281520192019161238b565b7fd7f73334000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b34610247576020600319360112610247577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917460206001600160a01b036128c3613f07565b6128cb614d6e565b16807fffffffffffffffffffffffff00000000000000000000000000000000000000006009541617600955604051908152a1005b34610247576000600319360112610247576000546001600160a01b0381163303612986577fffffffffffffffffffffffff0000000000000000000000000000000000000000600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b34610247576020600319360112610247576020611fbf6001600160a01b036129d6613f07565b16600052600c602052604060002054151590565b346102475760206003193601126102475767ffffffffffffffff612a0c613ead565b606060c0604051612a1c81613f85565b600081526000602082015260006040820152600083820152600060808201528260a0820152015216600052600a60205260606040600020610dfc612a5e615785565b611a7a604051612a6d81613f85565b84548152612b1a60018601549563ffffffff602084019760ff81161515895261ffff60408601818360081c168152818c880191818560181c1683528560808a019560281c168552612ad36003612ac560028a0161459f565b9860a08c01998a520161459f565b9860c08101998a526040519e8f9e8f9260408452516040840152511515910152511660808c0152511660a08a0152511660c08801525160e08088015261012087019061401f565b9051603f198683030161010087015261401f565b346102475760006003193601126102475760206001600160a01b0360095416604051908152f35b3461024757600060031936011261024757610dfc611237615785565b3461024757612b7f3661411c565b612b8a929192614d6e565b67ffffffffffffffff8216612bac816000526006602052604060002054151590565b15612bc75750610d1292612bc1913691614044565b90614dac565b7f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b3461024757612c1c612c24612c08366140ca565b9491612c15939193614d6e565b36916141e5565b9236916141e5565b7f000000000000000000000000000000000000000000000000000000000000000015612d2a5760005b8251811015612cb357806001600160a01b03612c6b600193866146c3565b5116612c7681615cde565b612c82575b5001612c4d565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a184612c7b565b5060005b8151811015610d1257806001600160a01b03612cd5600193856146c3565b51168015612d2457612ce681615875565b612cf3575b505b01612cb7565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a183612ceb565b50612ced565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b3461024757604060031936011261024757612d6d613ead565b60243567ffffffffffffffff811161024757602091612d93611fbf92369060040161407b565b90614528565b346102475760206003193601126102475760043567ffffffffffffffff811161024757806004016101006003198336030112610247576000604051612ddd81613f69565b52612e0a612e00612dfb612df460c48601856144ae565b3691614044565b614aee565b6064840135614bc9565b9060848301612e18816144ff565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116036118035750602483019077ffffffffffffffff00000000000000000000000000000000612e7283614513565b60801b16604051907f2cbc26bb00000000000000000000000000000000000000000000000000000000825260048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611771576000916130bf575b506117aa5767ffffffffffffffff612efa83614513565b16612f12816000526006602052604060002054151590565b1561177d5760206001600160a01b0360045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa908115611771576000916130a0575b50156116f357612f7d82614513565b90612f9360a4860192612d93612df485856144ae565b156130595750507ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0608067ffffffffffffffff612ffd612ff76044602098612fe389612fde8a614513565b614f6a565b01956109bf88612ff2896144ff565b614cd0565b946144ff565b936001600160a01b0360405195817f000000000000000000000000000000000000000000000000000000000000000016875233898801521660408601528560608601521692a28060405161305081613f69565b52604051908152f35b61306392506144ae565b6111526040519283927f24eb47e50000000000000000000000000000000000000000000000000000000084526020600485015260248401916143d5565b6130b9915060203d6020116117fc576117ee8183613fbd565b85612f6e565b6130d8915060203d6020116117fc576117ee8183613fbd565b85612ee3565b346102475760a0600319360112610247576130f7613ead565b6024359060443567ffffffffffffffff81116102475761311b903690600401613ed9565b9091606435926001600160a01b0384168094036102475760843567ffffffffffffffff811161024757613152903690600401613ed9565b505061315c614495565b506040519461316a86613f4d565b60008652600060208701526060608060405161318581613f31565b828152826020820152826040820152600083820152015267ffffffffffffffff8316936040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008560801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561177157600091613523575b506117aa5761324033615c69565b613257856000526006602052604060002054151590565b156134f55784600052600a602052604060002094855483116134c357509263ffffffff9261337b869361336d8a9760ff60016134019c9b01549161ffff8360081c16998360206132ba6132b461ffff8f9860181c1680988b61541c565b90614860565b9d019c8d5260281c1680613458575061ffff61332b6132db60038c0161459f565b985b604051976132ea89613f31565b8852602088019c8d52604088019586526060880193857f00000000000000000000000000000000000000000000000000000000000000001685523691614044565b9360808701948552816040519c8d986020808b01525160408a01525116606088015251166080860152511660a08401525160a060c084015260e083019061401f565b03601f198101865285613fbd565b60209586946040519061338e8783613fbd565b600082526133aa6002604051976133a489613f31565b0161459f565b8652868601526040850152606084015260808301526001600160a01b0360045416906040518097819482937f20487ded000000000000000000000000000000000000000000000000000000008452600484016142fb565b03915afa92831561177157600093613426575b50826040945283519283525190820152f35b9392508184813d8311613451575b61343e8183613fbd565b8101031261024757604093519293613414565b503d613434565b61332b61ffff916040519061346c82613f4d565b81526020810160018152604051917f181dcf100000000000000000000000000000000000000000000000000000000060208401525160248301525115156044820152604481526134bd606482613fbd565b986132dd565b90507f58dd87c50000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b847fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b61353c915060203d6020116117fc576117ee8183613fbd565b88613232565b346102475760006003193601126102475760206103546143f6565b3461024757600060031936011261024757602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102475760206003193601126102475760206135b6613f07565b6001600160a01b03807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102475760006003193601126102475760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b3461024757600060031936011261024757610dfc604051613654606082613fbd565b603581527f4275726e4d696e745769746845787465726e616c4d696e74657246617374547260208201527f616e73666572546f6b656e506f6f6c20312e362e300000000000000000000000604082015260405191829160208352602083019061401f565b60c0600319360112610247576136cc613ead565b60643567ffffffffffffffff8111610247576136ec903690600401613ed9565b9091608435916001600160a01b03831683036102475760a43567ffffffffffffffff811161024757613722903690600401613ed9565b50506040519061373182613f4d565b60008252600060208301526060608060405161374c81613f31565b82815282602082015282604082015260008382015201526040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008460801b1660048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561177157600091613d2b575b506117aa576137fb33615c69565b61381c67ffffffffffffffff84166000526006602052604060002054151590565b15613cf35767ffffffffffffffff8316600052600a602052604060002092835460243511613cb55760018401549561ffff8760081c169463ffffffff61ffff8960181c16986138716132b48b8a60243561541c565b602088015260281c1680613c51575061388c6003820161459f565b965b60405161389a81613f31565b60243581526020810197885260408101998a5260608101997f000000000000000000000000000000000000000000000000000000000000000060ff169a8b8152366138e6908988614044565b91608084019283526040519a8b9460208601602090525160408601525161ffff1660608501525161ffff1660808401525160ff1660a08301525160c0820160a0905260e082016139359161401f565b03601f19810188526139479088613fbd565b6020976040516139578a82613fbd565b60008082529861396f6002604051966133a488613f31565b85528a85015260408401526001600160a01b038216606084015260808301526001600160a01b03600454168860405180927f20487ded00000000000000000000000000000000000000000000000000000000825281806139d3888b600484016142fb565b03915afa908115613c46578891613c19575b5086526139f46024358561486d565b60208601516044358111613be95750613a93918891613a3760243530337f0000000000000000000000000000000000000000000000000000000000000000614902565b613a42602435615168565b6001600160a01b038116613b96575b506001600160a01b036004541660405180809581947f96f4e9f900000000000000000000000000000000000000000000000000000000835289600484016142fb565b039134905af1958615613b8a578096613b56575b5050957f240a1286fd41f1034c4032dcd6b93fc09e81be4a0b64c7ecee6260b605a8e01691613b4b86979867ffffffffffffffff613aeb60208901516024356143c8565b936020613b24613afc368b87614044565b7f0000000000000000000000000000000000000000000000000000000000000000888e614819565b99015160405196879687528d870152604086015260806060860152169560808401916143d5565b0390a4604051908152f35b909195508682813d8311613b83575b613b6f8183613fbd565b810103126117625750519381613b4b613aa7565b503d613b65565b604051903d90823e3d90fd5b613be390613bb0895130336001600160a01b038516614902565b8851906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000911661496c565b8a613a51565b7f61acdb930000000000000000000000000000000000000000000000000000000088526004526044803560245287fd5b90508881813d8311613c3f575b613c308183613fbd565b810103126126b757518a6139e5565b503d613c26565b6040513d8a823e3d90fd5b60405190613c5e82613f4d565b81526020810160018152604051917f181dcf10000000000000000000000000000000000000000000000000000000006020840152516024830152511515604482015260448152613caf606482613fbd565b9661388e565b67ffffffffffffffff907f58dd87c5000000000000000000000000000000000000000000000000000000006000521660045260243560245260446000fd5b67ffffffffffffffff837fa9902c7e000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b613d44915060203d6020116117fc576117ee8183613fbd565b866137ed565b3461024757602060031936011261024757600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361024757817ff6f46ff90000000000000000000000000000000000000000000000000000000060209314908115613e22575b8115613dc5575b5015158152f35b7f85572ffb00000000000000000000000000000000000000000000000000000000811491508115613df8575b5083613dbe565b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483613df1565b90507faff2afbf0000000000000000000000000000000000000000000000000000000081148015613e84575b8015613e5b575b90613db7565b507f01ffc9a7000000000000000000000000000000000000000000000000000000008114613e55565b507f0e64dd29000000000000000000000000000000000000000000000000000000008114613e4e565b6004359067ffffffffffffffff8216820361024757565b359067ffffffffffffffff8216820361024757565b9181601f840112156102475782359167ffffffffffffffff8311610247576020838186019501011161024757565b600435906001600160a01b038216820361024757565b35906001600160a01b038216820361024757565b60a0810190811067ffffffffffffffff8211176108fc57604052565b6040810190811067ffffffffffffffff8211176108fc57604052565b6020810190811067ffffffffffffffff8211176108fc57604052565b60e0810190811067ffffffffffffffff8211176108fc57604052565b6060810190811067ffffffffffffffff8211176108fc57604052565b90601f601f19910116810190811067ffffffffffffffff8211176108fc57604052565b67ffffffffffffffff81116108fc57601f01601f191660200190565b60005b83811061400f5750506000910152565b8181015183820152602001613fff565b90601f19601f60209361403d81518092818752878088019101613ffc565b0116010190565b92919261405082613fe0565b9161405e6040519384613fbd565b829481845281830111610247578281602093846000960137010152565b9080601f830112156102475781602061409693359101614044565b90565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501948460051b01011161024757565b60406003198201126102475760043567ffffffffffffffff811161024757816140f591600401614099565b929092916024359067ffffffffffffffff82116102475761411891600401614099565b9091565b9060406003198301126102475760043567ffffffffffffffff8116810361024757916024359067ffffffffffffffff82116102475761411891600401613ed9565b906020808351928381520192019060005b81811061417b5750505090565b82516001600160a01b031684526020938401939092019160010161416e565b6020600319820112610247576004359067ffffffffffffffff8211610247576003198260a0920301126102475760040190565b67ffffffffffffffff81116108fc5760051b60200190565b9291906141f1816141cd565b936141ff6040519586613fbd565b602085838152019160051b810192831161024757905b82821061422157505050565b6020809161422e84613f1d565b815201910190614215565b9080601f8301121561024757816020614096933591016141e5565b9181601f840112156102475782359167ffffffffffffffff8311610247576020808501946060850201011161024757565b9060038210156122055752565b35906fffffffffffffffffffffffffffffffff8216820361024757565b9190826060910312610247576040516142c781613fa1565b809280359081151582036102475760406142f691819385526142eb60208201614292565b602086015201614292565b910152565b9067ffffffffffffffff909392931681526040602082015261434261432c845160a0604085015260e084019061401f565b6020850151603f1984830301606085015261401f565b90604084015191603f198282030160808301526020808451928381520193019060005b81811061439d575050506080846001600160a01b036060614096969701511660a084015201519060c0603f198285030191015261401f565b825180516001600160a01b031686526020908101518187015260409095019490920191600101614365565b91908203918211611f2d57565b601f8260209493601f19938186528686013760008582860101520116010190565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526020816024816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa90811561177157600091614466575090565b90506020813d60201161448d575b8161448160209383613fbd565b81010312610247575190565b3d9150614474565b604051906144a282613f4d565b60006020838281520152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610247570180359067ffffffffffffffff82116102475760200191813603831361024757565b356001600160a01b03811681036102475790565b3567ffffffffffffffff811681036102475790565b9067ffffffffffffffff61409692166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b90600182811c92168015614595575b602083101461457f57565b634e487b7160e01b600052602260045260246000fd5b91607f1691614574565b90604051918260008254926145b384614565565b808452936001811690811561462157506001146145da575b506145d892500383613fbd565b565b90506000929192526020600020906000915b8183106146055750509060206145d892820101386145cb565b60209193508060019154838589010152019101909184926145ec565b602093506145d89592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b820101386145cb565b60038210156122055752565b91908110156146ad5760051b810135907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee181360301821215610247570190565b634e487b7160e01b600052603260045260246000fd5b80518210156146ad5760209160051b010190565b91908110156146ad5760051b0190565b91908110156146ad576060020190565b6040519061470482613f31565b60006080838281528260208201528260408201528260608201520152565b9060405161472f81613f31565b60806001829460ff81546fffffffffffffffffffffffffffffffff8116865263ffffffff81861c16602087015260a01c161515604085015201546fffffffffffffffffffffffffffffffff81166060840152811c910152565b67ffffffffffffffff166000526007602052614096600460406000200161459f565b8181106147b5575050565b600081556001016147aa565b81810292918115918404141715611f2d57565b9190601f81116147e357505050565b6145d8926000526020600020906020601f840160051c8301931061480f575b601f0160051c01906147aa565b9091508190614802565b929061484c61485a9260ff60405195869460208601988952604086015216606084015260808084015260a083019061401f565b03601f198101835282613fbd565b51902090565b91908201809211611f2d57565b67ffffffffffffffff7fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da8178944911691826000526007602052806148de60406000206001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615443565b604080516001600160a01b039092168252602082019290925290819081015b0390a2565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000060208201526001600160a01b03928316602482015292909116604483015260648201929092526145d89161496782608481015b03601f198101845283613fbd565b615652565b91909181158015614a54575b156149d0576040517f095ea7b30000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909316602484015260448301919091526145d891906149678260648101614959565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b0384166024820152602081806044810103816001600160a01b0386165afa90811561177157600091614abc575b5015614978565b90506020813d602011614ae6575b81614ad760209383613fbd565b81010312610247575138614ab5565b3d9150614aca565b80518015614b5e57602003614b2057805160208281019183018390031261024757519060ff8211614b20575060ff1690565b611152906040519182917f953576f700000000000000000000000000000000000000000000000000000000835260206004840152602483019061401f565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff8211611f2d57565b60ff16604d8111611f2d57600a0a90565b8115614bb3570490565b634e487b7160e01b600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff811692828414614cb157828411614c875790614c0e91614b84565b91604d60ff8416118015614c6c575b614c3657505090614c3061409692614b98565b906147c1565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b50614c7683614b98565b8015614bb357600019048411614c1d565b614c9091614b84565b91604d60ff841611614c3657505090614cab61409692614b98565b90614ba9565b5050505090565b90816020910312610247575180151581036102475790565b6040517f40c10f190000000000000000000000000000000000000000000000000000000081526001600160a01b03909116600482015260248101919091526020818060448101038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1801561177157614d535750565b614d6b9060203d6020116117fc576117ee8183613fbd565b50565b6001600160a01b03600154163303614d8257565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b9080511561094b5767ffffffffffffffff81516020830120921691826000526007602052614de1816005604060002001615926565b15614f175760005260086020526040600020815167ffffffffffffffff81116108fc57614e1881614e128454614565565b846147d4565b6020601f8211600114614e8d5791614e6c827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea95936148fd95600091614e82575b506000198260011b9260031b1c19161790565b905560405191829160208352602083019061401f565b905084015138614e59565b601f1982169083600052806000209160005b818110614eff5750926148fd9492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea989610614ee6575b5050811b019055610fda565b85015160001960f88460031b161c191690553880614eda565b9192602060018192868a015181550194019201614e9f565b50906111526040519283927f393b8ad2000000000000000000000000000000000000000000000000000000008452600484015260406024840152604483019061401f565b519061ffff8216820361024757565b67ffffffffffffffff7f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c911691826000526007602052806148de60026040600020016001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016928391615443565b357fffffffff00000000000000000000000000000000000000000000000000000000811681036102475790565b3563ffffffff811681036102475790565b3561ffff811681036102475790565b3580151581036102475790565b67ffffffffffffffff16600081815260066020526040902054909291901561513a579161513760e0926151038561508f7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b976152d5565b8460005260076020526150a6816040600020615a36565b6150af836152d5565b8460005260076020526150c9836002604060002001615a36565b60405194855260208501906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b60808301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565ba1565b827f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b602060009160246001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016916151c581847f000000000000000000000000000000000000000000000000000000000000000061496c565b60405194859384927f42966c6800000000000000000000000000000000000000000000000000000000845260048401525af1801561177157614d535750565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000060208201526001600160a01b03909216602483015260448201929092526145d8916149678260648101614959565b61525e6146f7565b506fffffffffffffffffffffffffffffffff6060820151166fffffffffffffffffffffffffffffffff80835116916152b560208501936132b46152a863ffffffff875116426143c8565b85608089015116906147c1565b808210156152ce57505b16825263ffffffff4216905290565b90506152bf565b805115615375576fffffffffffffffffffffffffffffffff6040820151166fffffffffffffffffffffffffffffffff602083015116106153125750565b606490615373604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565bfd5b6fffffffffffffffffffffffffffffffff604082015116158015906153fd575b61539c5750565b606490615373604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b506fffffffffffffffffffffffffffffffff6020820151161515615395565b61543f9061ffff61271061543682829698979816846147c1565b049516906147c1565b0490565b9182549060ff8260a01c1615801561564a575b615644576fffffffffffffffffffffffffffffffff8216916001850190815461549b63ffffffff6fffffffffffffffffffffffffffffffff83169360801c16426143c8565b90816155a6575b505084811061556757508383106154fc5750506154d16fffffffffffffffffffffffffffffffff9283926143c8565b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416179055565b5460801c9161550b81856143c8565b92600019810190808211611f2d5761552e615533926001600160a01b0396614860565b614ba9565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b82856001600160a01b03927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b82869293961161561a576155c1926132b49160801c906147c1565b808410156156155750825b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff00000000000000000000000000000000161786559238806154a2565b6155cc565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b508215615456565b6001600160a01b036156d49116916040926000808551936156738786613fbd565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d1561577d573d916156b883613fe0565b926156c587519485613fbd565b83523d6000602085013e615ea6565b805190816156e157505050565b6020806156f2938301019101614cb8565b156156fa5750565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606091615ea6565b60405190600b548083528260208101600b60005260206000209260005b8181106157b75750506145d892500383613fbd565b84548352600194850194879450602090930192016157a2565b906040519182815491828252602082019060005260206000209260005b8181106158025750506145d892500383613fbd565b84548352600194850194879450602090930192016157ed565b80548210156146ad5760005260206000200190600090565b805490680100000000000000008210156108fc578161585a9160016158719401815561581b565b81939154906000199060031b92831b921b19161790565b9055565b806000526003602052604060002054156000146158ae57615897816002615833565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146158ae576158d681600b615833565b600b5490600052600c602052604060002055600190565b806000526006602052604060002054156000146158ae5761590f816005615833565b600554906000526006602052604060002055600190565b600082815260018201602052604090205461595d578061594883600193615833565b80549260005201602052604060002055600190565b5050600090565b8054801561598c57600019019061597b828261581b565b60001982549160031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c6020526040902054801561595d576000198101818111611f2d57600b54906000198201918211611f2d578082036159fc575b5050506159e8600b615964565b600052600c60205260006040812055600190565b615a1e615a0d61585a93600b61581b565b90549060031b1c928392600b61581b565b9055600052600c6020526040600020553880806159db565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1991615b6f6060928054615a7363ffffffff8260801c16426143c8565b9081615bae575b50506fffffffffffffffffffffffffffffffff6001816020860151169282815416808510600014615ba657508280855b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416178155615b238651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601517fffffffffffffffffffffffffffffffff0000000000000000000000000000000060809190911b16939092166fffffffffffffffffffffffffffffffff1692909217910155565b61513760405180926fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b838091615aaa565b6fffffffffffffffffffffffffffffffff91615be3839283615bdc6001880154948286169560801c906147c1565b9116614860565b80821015615c6257505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9290911692909216167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116174260801b73ffffffff00000000000000000000000000000000161781553880615a7a565b9050615bed565b7f0000000000000000000000000000000000000000000000000000000000000000615c915750565b6001600160a01b031680600052600360205260406000205415615cb15750565b7fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b600081815260036020526040902054801561595d576000198101818111611f2d57600254906000198201918211611f2d57818103615d38575b505050615d246002615964565b600052600360205260006040812055600190565b615d5a615d4961585a93600261581b565b90549060031b1c928392600261581b565b90556000526003602052604060002055388080615d17565b600081815260066020526040902054801561595d576000198101818111611f2d57600554906000198201918211611f2d57818103615dcc575b505050615db86005615964565b600052600660205260006040812055600190565b615dee615ddd61585a93600561581b565b90549060031b1c928392600561581b565b90556000526006602052604060002055388080615dab565b906001820191816000528260205260406000205490811515600014615e9d57600019820191808311611f2d5781546000198101908111611f2d578381615e549503615e66575b505050615964565b60005260205260006040812055600190565b615e86615e7661585a938661581b565b90549060031b1c9283928661581b565b905560005284602052604060002055388080615e4c565b50505050600090565b91929015615f215750815115615eba575090565b3b15615ec35790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015615f345750805190602001fd5b611152906040519182917f08c379a000000000000000000000000000000000000000000000000000000000835260206004840152602483019061401f56fea164736f6c634300081a000a",
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

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactor) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.contract.Transact(opts, "ccipSendToken", destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}

func (_BurnMintWithExternalMinterFastTransferTokenPool *BurnMintWithExternalMinterFastTransferTokenPoolTransactorSession) CcipSendToken(destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterFastTransferTokenPool.Contract.CcipSendToken(&_BurnMintWithExternalMinterFastTransferTokenPool.TransactOpts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
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

	CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error)

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
