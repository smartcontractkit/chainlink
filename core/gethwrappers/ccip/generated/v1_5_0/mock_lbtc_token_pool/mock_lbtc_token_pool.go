// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_lbtc_token_pool

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

var MockLBTCTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expected\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"actual\",\"type\":\"uint8\"}],\"name\":\"InvalidDecimalArgs\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"remoteDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"localDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"remoteAmount\",\"type\":\"uint256\"}],\"name\":\"OverflowDetected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_destPoolData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162003a5138038062003a51833981016040819052620000359162000663565b846008858585336000816200005d57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038481169190911790915581161562000090576200009081620001fe565b50506001600160a01b0385161580620000b057506001600160a01b038116155b80620000c357506001600160a01b038216155b15620000e2576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808616608081905290831660c0526040805163313ce56760e01b8152905163313ce567916004808201926020929091908290030181865afa92505050801562000152575060408051601f3d908101601f191682019092526200014f9181019062000785565b60015b1562000192578060ff168560ff161462000190576040516332ad3e0760e11b815260ff80871660048301528216602482015260440160405180910390fd5b505b60ff841660a052600480546001600160a01b0319166001600160a01b038316179055825115801560e052620001dc57604080516000815260208101909152620001dc908462000278565b505050505080600a9081620001f2919062000841565b5050505050506200095b565b336001600160a01b038216036200022857604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e05162000299576040516335f4a7b360e01b815260040160405180910390fd5b60005b825181101562000324576000838281518110620002bd57620002bd6200090d565b60209081029190910101519050620002d7600282620003d5565b156200031a576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016200029c565b5060005b8151811015620003d05760008282815181106200034957620003496200090d565b6020026020010151905060006001600160a01b0316816001600160a01b031603620003755750620003c7565b62000382600282620003f5565b15620003c5576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b60010162000328565b505050565b6000620003ec836001600160a01b0384166200040c565b90505b92915050565b6000620003ec836001600160a01b03841662000510565b60008181526001830160205260408120548015620005055760006200043360018362000923565b8554909150600090620004499060019062000923565b9050808214620004b55760008660000182815481106200046d576200046d6200090d565b90600052602060002001549050808760000184815481106200049357620004936200090d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620004c957620004c962000945565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003ef565b6000915050620003ef565b60008181526001830160205260408120546200055957508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003ef565b506000620003ef565b6001600160a01b03811681146200057857600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620005bc57620005bc6200057b565b604052919050565b8051620005d18162000562565b919050565b600082601f830112620005e857600080fd5b81516001600160401b038111156200060457620006046200057b565b60206200061a601f8301601f1916820162000591565b82815285828487010111156200062f57600080fd5b60005b838110156200064f57858101830151828201840152820162000632565b506000928101909101919091529392505050565b600080600080600060a086880312156200067c57600080fd5b8551620006898162000562565b602087810151919650906001600160401b0380821115620006a957600080fd5b818901915089601f830112620006be57600080fd5b815181811115620006d357620006d36200057b565b8060051b620006e485820162000591565b918252838101850191858101908d841115620006ff57600080fd5b948601945b838610156200072d57855192506200071c8362000562565b828252948601949086019062000704565b9950620007419250505060408a01620005c4565b95506200075160608a01620005c4565b945060808901519250808311156200076857600080fd5b50506200077888828901620005d6565b9150509295509295909350565b6000602082840312156200079857600080fd5b815160ff81168114620007aa57600080fd5b9392505050565b600181811c90821680620007c657607f821691505b602082108103620007e757634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620003d0576000816000526020600020601f850160051c81016020861015620008185750805b601f850160051c820191505b81811015620008395782815560010162000824565b505050505050565b81516001600160401b038111156200085d576200085d6200057b565b62000875816200086e8454620007b1565b84620007ed565b602080601f831160018114620008ad5760008415620008945750858301515b600019600386901b1c1916600185901b17855562000839565b600085815260208120601f198616915b82811015620008de57888601518255948401946001909101908401620008bd565b5085821015620008fd5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003ef57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05160e0516130a0620009b160003960008181610562015261199f0152600061053c015260006102eb015260008181610252015281816102a7015281816107450152610ba201526130a06000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80639a4575b911610104578063c0d78655116100a2578063dc0bd97111610071578063dc0bd9711461053a578063e0351e1314610560578063e8a1da1714610586578063f2fde38b1461059957600080fd5b8063c0d78655146104ec578063c4bffe2b146104ff578063c75eea9c14610514578063cf7401f31461052757600080fd5b8063acfecf91116100de578063acfecf9114610439578063af58d59f1461044c578063b0f479a1146104bb578063b7946580146104d957600080fd5b80639a4575b9146103e4578063a42a7b8b14610404578063a7cd63b71461042457600080fd5b80634c5ef0ed1161017c57806379ba50971161014b57806379ba5097146103985780637d54534e146103a05780638926f54f146103b35780638da5cb5b146103c657600080fd5b80634c5ef0ed1461033f57806354c8a4f31461035257806362ddd3c4146103675780636d3d1a581461037a57600080fd5b8063240028e8116101b8578063240028e81461029757806324f65ee7146102e457806332a7a82214610315578063390775371461031d57600080fd5b806301ffc9a7146101df578063181f5a771461020757806321df0da714610250575b600080fd5b6101f26101ed36600461243d565b6105ac565b60405190151581526020015b60405180910390f35b6102436040518060400160405280601781526020017f4d6f636b4c425443546f6b656e506f6f6c20312e352e3100000000000000000081525081565b6040516101fe91906124e3565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fe565b6101f26102a53660046124f6565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101fe565b610243610691565b61033061032b36600461252c565b61071f565b604051905181526020016101fe565b6101f261034d366004612585565b610894565b610365610360366004612654565b6108de565b005b610365610375366004612585565b610959565b60095473ffffffffffffffffffffffffffffffffffffffff16610272565b6103656109f6565b6103656103ae3660046124f6565b610ac4565b6101f26103c13660046126c0565b610b45565b60015473ffffffffffffffffffffffffffffffffffffffff16610272565b6103f76103f23660046126db565b610b5c565b6040516101fe9190612716565b6104176104123660046126c0565b610d07565b6040516101fe919061276d565b61042c610e72565b6040516101fe91906127ef565b610365610447366004612585565b610e83565b61045f61045a3660046126c0565b610f9b565b6040516101fe919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610272565b6102436104e73660046126c0565b611070565b6103656104fa3660046124f6565b611120565b6105076111fb565b6040516101fe9190612849565b61045f6105223660046126c0565b6112b3565b6103656105353660046129c8565b611385565b7f0000000000000000000000000000000000000000000000000000000000000000610272565b7f00000000000000000000000000000000000000000000000000000000000000006101f2565b610365610594366004612654565b611409565b6103656105a73660046124f6565b61191b565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061063f57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061068b57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b600a805461069e90612a0d565b80601f01602080910402602001604051908101604052809291908181526020018280546106ca90612a0d565b80156107175780601f106106ec57610100808354040283529160200191610717565b820191906000526020600020905b8154815290600101906020018083116106fa57829003601f168201915b505050505081565b60408051602081019091526000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f1961077a60608501604086016124f6565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260608501356024820152604401600060405180830381600087803b1580156107ea57600080fd5b505af11580156107fe573d6000803e3d6000fd5b506108139250505060608301604084016124f6565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0846060013560405161087591815260200190565b60405180910390a3506040805160208101909152606090910135815290565b60006108d683836040516108a9929190612a60565b604080519182900390912067ffffffffffffffff871660009081526007602052919091206005019061192f565b949350505050565b6108e661194a565b6109538484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061199d92505050565b50505050565b61096161194a565b61096a83610b45565b6109b1576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b6109f18383838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b5392505050565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a47576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610acc61194a565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b600061068b600567ffffffffffffffff841661192f565b60408051808201909152606080825260208201526040517f42966c68000000000000000000000000000000000000000000000000000000008152606083013560048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906342966c6890602401600060405180830381600087803b158015610bfb57600080fd5b505af1158015610c0f573d6000803e3d6000fd5b5050604051606085013581523392507f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7915060200160405180910390a26040518060400160405280610c6d8460200160208101906104e791906126c0565b8152602001600a8054610c7f90612a0d565b80601f0160208091040260200160405190810160405280929190818152602001828054610cab90612a0d565b8015610cf85780601f10610ccd57610100808354040283529160200191610cf8565b820191906000526020600020905b815481529060010190602001808311610cdb57829003601f168201915b50505050508152509050919050565b67ffffffffffffffff8116600090815260076020526040812060609190610d3090600501611c4d565b90506000815167ffffffffffffffff811115610d4e57610d4e61288b565b604051908082528060200260200182016040528015610d8157816020015b6060815260200190600190039081610d6c5790505b50905060005b8251811015610e6a5760086000848381518110610da657610da6612a70565b602002602001015181526020019081526020016000208054610dc790612a0d565b80601f0160208091040260200160405190810160405280929190818152602001828054610df390612a0d565b8015610e405780601f10610e1557610100808354040283529160200191610e40565b820191906000526020600020905b815481529060010190602001808311610e2357829003601f168201915b5050505050828281518110610e5757610e57612a70565b6020908102919091010152600101610d87565b509392505050565b6060610e7e6002611c4d565b905090565b610e8b61194a565b610e9483610b45565b610ed6576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024016109a8565b610f168282604051610ee9929190612a60565b604080519182900390912067ffffffffffffffff8616600090815260076020529190912060050190611c5a565b610f52578282826040517f74f23c7c0000000000000000000000000000000000000000000000000000000081526004016109a893929190612ae8565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051610f8e929190612b0c565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261068b90611c66565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061109b90612a0d565b80601f01602080910402602001604051908101604052809291908181526020018280546110c790612a0d565b80156111145780601f106110e957610100808354040283529160200191611114565b820191906000526020600020905b8154815290600101906020018083116110f757829003601f168201915b50505050509050919050565b61112861194a565b73ffffffffffffffffffffffffffffffffffffffff8116611175576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006112096005611c4d565b90506000815167ffffffffffffffff8111156112275761122761288b565b604051908082528060200260200182016040528015611250578160200160208202803683370190505b50905060005b82518110156112ac5782818151811061127157611271612a70565b602002602001015182828151811061128b5761128b612a70565b67ffffffffffffffff90921660209283029190910190910152600101611256565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261068b90611c66565b60095473ffffffffffffffffffffffffffffffffffffffff1633148015906113c5575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156113fe576040517f8e4a23d60000000000000000000000000000000000000000000000000000000081523360048201526024016109a8565b6109f1838383611d18565b61141161194a565b60005b838110156115fe57600085858381811061143057611430612a70565b905060200201602081019061144591906126c0565b905061145c600567ffffffffffffffff8316611c5a565b61149e576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201526024016109a8565b67ffffffffffffffff811660009081526007602052604081206114c390600501611c4d565b905060005b815181101561152f576115268282815181106114e6576114e6612a70565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600501611c5a90919063ffffffff16565b506001016114c8565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff0000000000000000000000000000000000000000009081168255600182018390556002820180549091169055600381018290559061159860048301826123d0565b60058201600081816115aa828261240a565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916945060200192506115ec915050565b60405180910390a15050600101611414565b5060005b8181101561191457600083838381811061161e5761161e612a70565b90506020028101906116309190612b20565b61163990612bec565b905061164a81606001516000611e02565b61165981608001516000611e02565b806040015151600003611698576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516116b09060059067ffffffffffffffff16611f3f565b6116f55780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024016109a8565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a9091029990931617179094169590951790925590920290911760038201559082015160048201906118789082612d63565b5060005b8260200151518110156118bc576118b48360000151846020015183815181106118a7576118a7612a70565b6020026020010151611b53565b60010161187c565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c282600001518360400151846060015185608001516040516119029493929190612e7d565b60405180910390a15050600101611602565b5050505050565b61192361194a565b61192c81611f4b565b50565b600081815260018301602052604081205415155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461199b576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f00000000000000000000000000000000000000000000000000000000000000006119f4576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611a8a576000838281518110611a1457611a14612a70565b60200260200101519050611a3281600261200f90919063ffffffff16565b15611a815760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016119f7565b5060005b81518110156109f1576000828281518110611aab57611aab612a70565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611aef5750611b4b565b611afa600282612031565b15611b495760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611a8e565b8051600003611b8e576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff8416600090815260079092526040909120611bc09060050182611f3f565b611bfa5782826040517f393b8ad20000000000000000000000000000000000000000000000000000000081526004016109a8929190612f16565b6000818152600860205260409020611c128382612d63565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea83604051610f8e91906124e3565b6060600061194383612053565b600061194383836120ae565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611cf482606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611cd89190612f68565b85608001516fffffffffffffffffffffffffffffffff166121a1565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b611d2183610b45565b611d63576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024016109a8565b611d6e826000611e02565b67ffffffffffffffff83166000908152600760205260409020611d9190836121c9565b611d9c816000611e02565b67ffffffffffffffff83166000908152600760205260409020611dc290600201826121c9565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b838383604051611df593929190612f7b565b60405180910390a1505050565b815115611ecd5781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580611e58575060408201516fffffffffffffffffffffffffffffffff16155b15611e9157816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016109a89190612ffe565b8015611ec9576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580611f06575060208201516fffffffffffffffffffffffffffffffff1615155b15611ec957816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016109a89190612ffe565b6000611943838361236b565b3373ffffffffffffffffffffffffffffffffffffffff821603611f9a576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006119438373ffffffffffffffffffffffffffffffffffffffff84166120ae565b60006119438373ffffffffffffffffffffffffffffffffffffffff841661236b565b60608160000180548060200260200160405190810160405280929190818152602001828054801561111457602002820191906000526020600020905b81548152602001906001019080831161208f5750505050509050919050565b600081815260018301602052604081205480156121975760006120d2600183612f68565b85549091506000906120e690600190612f68565b905080821461214b57600086600001828154811061210657612106612a70565b906000526020600020015490508087600001848154811061212957612129612a70565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061215c5761215c61303a565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061068b565b600091505061068b565b60006121c0856121b18486613069565b6121bb9087613080565b6123ba565b95945050505050565b81546000906121f290700100000000000000000000000000000000900463ffffffff1642612f68565b90508015612294576001830154835461223a916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166121a1565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546122ba916fffffffffffffffffffffffffffffffff90811691166123ba565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990611df5908490612ffe565b60008181526001830160205260408120546123b25750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561068b565b50600061068b565b60008183106123c95781611943565b5090919050565b5080546123dc90612a0d565b6000825580601f106123ec575050565b601f01602090049060005260206000209081019061192c9190612424565b508054600082559060005260206000209081019061192c91905b5b808211156124395760008155600101612425565b5090565b60006020828403121561244f57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461194357600080fd5b6000815180845260005b818110156124a557602081850181015186830182015201612489565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611943602083018461247f565b60006020828403121561250857600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461194357600080fd5b60006020828403121561253e57600080fd5b813567ffffffffffffffff81111561255557600080fd5b8201610100818503121561194357600080fd5b803567ffffffffffffffff8116811461258057600080fd5b919050565b60008060006040848603121561259a57600080fd5b6125a384612568565b9250602084013567ffffffffffffffff808211156125c057600080fd5b818601915086601f8301126125d457600080fd5b8135818111156125e357600080fd5b8760208285010111156125f557600080fd5b6020830194508093505050509250925092565b60008083601f84011261261a57600080fd5b50813567ffffffffffffffff81111561263257600080fd5b6020830191508360208260051b850101111561264d57600080fd5b9250929050565b6000806000806040858703121561266a57600080fd5b843567ffffffffffffffff8082111561268257600080fd5b61268e88838901612608565b909650945060208701359150808211156126a757600080fd5b506126b487828801612608565b95989497509550505050565b6000602082840312156126d257600080fd5b61194382612568565b6000602082840312156126ed57600080fd5b813567ffffffffffffffff81111561270457600080fd5b820160a0818503121561194357600080fd5b602081526000825160406020840152612732606084018261247f565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160408501526121c0828261247f565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b828110156127e2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526127d085835161247f565b94509285019290850190600101612796565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561283d57835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161280b565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561283d57835167ffffffffffffffff1683529284019291840191600101612865565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff811182821017156128dd576128dd61288b565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561292a5761292a61288b565b604052919050565b80356fffffffffffffffffffffffffffffffff8116811461258057600080fd5b60006060828403121561296457600080fd5b6040516060810181811067ffffffffffffffff821117156129875761298761288b565b6040529050808235801515811461299d57600080fd5b81526129ab60208401612932565b60208201526129bc60408401612932565b60408201525092915050565b600080600060e084860312156129dd57600080fd5b6129e684612568565b92506129f58560208601612952565b9150612a048560808601612952565b90509250925092565b600181811c90821680612a2157607f821691505b602082108103612a5a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff841681526040602082015260006121c0604083018486612a9f565b6020815260006108d6602083018486612a9f565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee1833603018112612b5457600080fd5b9190910192915050565b600082601f830112612b6f57600080fd5b813567ffffffffffffffff811115612b8957612b8961288b565b612bba60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016128e3565b818152846020838601011115612bcf57600080fd5b816020850160208301376000918101602001919091529392505050565b60006101208236031215612bff57600080fd5b612c076128ba565b612c1083612568565b815260208084013567ffffffffffffffff80821115612c2e57600080fd5b9085019036601f830112612c4157600080fd5b813581811115612c5357612c5361288b565b8060051b612c628582016128e3565b9182528381018501918581019036841115612c7c57600080fd5b86860192505b83831015612cb857823585811115612c9a5760008081fd5b612ca83689838a0101612b5e565b8352509186019190860190612c82565b8087890152505050506040860135925080831115612cd557600080fd5b5050612ce336828601612b5e565b604083015250612cf63660608501612952565b6060820152612d083660c08501612952565b608082015292915050565b601f8211156109f1576000816000526020600020601f850160051c81016020861015612d3c5750805b601f850160051c820191505b81811015612d5b57828155600101612d48565b505050505050565b815167ffffffffffffffff811115612d7d57612d7d61288b565b612d9181612d8b8454612a0d565b84612d13565b602080601f831160018114612de45760008415612dae5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555612d5b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015612e3157888601518255948401946001909101908401612e12565b5085821015612e6d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152612ea18184018761247f565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150612edf9050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e08301526121c0565b67ffffffffffffffff831681526040602082015260006108d6604083018461247f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561068b5761068b612f39565b67ffffffffffffffff8416815260e08101612fc760208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c08301526108d6565b6060810161068b82848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b808202811582820484141761068b5761068b612f39565b8082018082111561068b5761068b612f3956fea164736f6c6343000818000a",
}

var MockLBTCTokenPoolABI = MockLBTCTokenPoolMetaData.ABI

var MockLBTCTokenPoolBin = MockLBTCTokenPoolMetaData.Bin

func DeployMockLBTCTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address, destPoolData []byte) (common.Address, *types.Transaction, *MockLBTCTokenPool, error) {
	parsed, err := MockLBTCTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockLBTCTokenPoolBin), backend, token, allowlist, rmnProxy, router, destPoolData)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockLBTCTokenPool{address: address, abi: *parsed, MockLBTCTokenPoolCaller: MockLBTCTokenPoolCaller{contract: contract}, MockLBTCTokenPoolTransactor: MockLBTCTokenPoolTransactor{contract: contract}, MockLBTCTokenPoolFilterer: MockLBTCTokenPoolFilterer{contract: contract}}, nil
}

type MockLBTCTokenPool struct {
	address common.Address
	abi     abi.ABI
	MockLBTCTokenPoolCaller
	MockLBTCTokenPoolTransactor
	MockLBTCTokenPoolFilterer
}

type MockLBTCTokenPoolCaller struct {
	contract *bind.BoundContract
}

type MockLBTCTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type MockLBTCTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type MockLBTCTokenPoolSession struct {
	Contract     *MockLBTCTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockLBTCTokenPoolCallerSession struct {
	Contract *MockLBTCTokenPoolCaller
	CallOpts bind.CallOpts
}

type MockLBTCTokenPoolTransactorSession struct {
	Contract     *MockLBTCTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type MockLBTCTokenPoolRaw struct {
	Contract *MockLBTCTokenPool
}

type MockLBTCTokenPoolCallerRaw struct {
	Contract *MockLBTCTokenPoolCaller
}

type MockLBTCTokenPoolTransactorRaw struct {
	Contract *MockLBTCTokenPoolTransactor
}

func NewMockLBTCTokenPool(address common.Address, backend bind.ContractBackend) (*MockLBTCTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(MockLBTCTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockLBTCTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPool{address: address, abi: abi, MockLBTCTokenPoolCaller: MockLBTCTokenPoolCaller{contract: contract}, MockLBTCTokenPoolTransactor: MockLBTCTokenPoolTransactor{contract: contract}, MockLBTCTokenPoolFilterer: MockLBTCTokenPoolFilterer{contract: contract}}, nil
}

func NewMockLBTCTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*MockLBTCTokenPoolCaller, error) {
	contract, err := bindMockLBTCTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolCaller{contract: contract}, nil
}

func NewMockLBTCTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*MockLBTCTokenPoolTransactor, error) {
	contract, err := bindMockLBTCTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolTransactor{contract: contract}, nil
}

func NewMockLBTCTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*MockLBTCTokenPoolFilterer, error) {
	contract, err := bindMockLBTCTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolFilterer{contract: contract}, nil
}

func bindMockLBTCTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockLBTCTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockLBTCTokenPool.Contract.MockLBTCTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.MockLBTCTokenPoolTransactor.contract.Transfer(opts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.MockLBTCTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockLBTCTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.contract.Transfer(opts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetAllowList(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetAllowList(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _MockLBTCTokenPool.Contract.GetAllowListEnabled(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _MockLBTCTokenPool.Contract.GetAllowListEnabled(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _MockLBTCTokenPool.Contract.GetCurrentInboundRateLimiterState(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _MockLBTCTokenPool.Contract.GetCurrentInboundRateLimiterState(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _MockLBTCTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _MockLBTCTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetRateLimitAdmin(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetRateLimitAdmin(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _MockLBTCTokenPool.Contract.GetRemotePools(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _MockLBTCTokenPool.Contract.GetRemotePools(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _MockLBTCTokenPool.Contract.GetRemoteToken(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _MockLBTCTokenPool.Contract.GetRemoteToken(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetRmnProxy(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetRmnProxy(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetRouter() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetRouter(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetRouter(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _MockLBTCTokenPool.Contract.GetSupportedChains(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _MockLBTCTokenPool.Contract.GetSupportedChains(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetToken() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetToken(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.GetToken(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _MockLBTCTokenPool.Contract.GetTokenDecimals(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _MockLBTCTokenPool.Contract.GetTokenDecimals(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) IDestPoolData(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "i_destPoolData")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) IDestPoolData() ([]byte, error) {
	return _MockLBTCTokenPool.Contract.IDestPoolData(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) IDestPoolData() ([]byte, error) {
	return _MockLBTCTokenPool.Contract.IDestPoolData(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _MockLBTCTokenPool.Contract.IsRemotePool(&_MockLBTCTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _MockLBTCTokenPool.Contract.IsRemotePool(&_MockLBTCTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _MockLBTCTokenPool.Contract.IsSupportedChain(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _MockLBTCTokenPool.Contract.IsSupportedChain(&_MockLBTCTokenPool.CallOpts, remoteChainSelector)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _MockLBTCTokenPool.Contract.IsSupportedToken(&_MockLBTCTokenPool.CallOpts, token)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _MockLBTCTokenPool.Contract.IsSupportedToken(&_MockLBTCTokenPool.CallOpts, token)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) Owner() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.Owner(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) Owner() (common.Address, error) {
	return _MockLBTCTokenPool.Contract.Owner(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MockLBTCTokenPool.Contract.SupportsInterface(&_MockLBTCTokenPool.CallOpts, interfaceId)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MockLBTCTokenPool.Contract.SupportsInterface(&_MockLBTCTokenPool.CallOpts, interfaceId)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockLBTCTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) TypeAndVersion() (string, error) {
	return _MockLBTCTokenPool.Contract.TypeAndVersion(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _MockLBTCTokenPool.Contract.TypeAndVersion(&_MockLBTCTokenPool.CallOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.AcceptOwnership(&_MockLBTCTokenPool.TransactOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.AcceptOwnership(&_MockLBTCTokenPool.TransactOpts)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.AddRemotePool(&_MockLBTCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.AddRemotePool(&_MockLBTCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.ApplyAllowListUpdates(&_MockLBTCTokenPool.TransactOpts, removes, adds)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.ApplyAllowListUpdates(&_MockLBTCTokenPool.TransactOpts, removes, adds)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.ApplyChainUpdates(&_MockLBTCTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.ApplyChainUpdates(&_MockLBTCTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.LockOrBurn(&_MockLBTCTokenPool.TransactOpts, lockOrBurnIn)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.LockOrBurn(&_MockLBTCTokenPool.TransactOpts, lockOrBurnIn)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.ReleaseOrMint(&_MockLBTCTokenPool.TransactOpts, releaseOrMintIn)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.ReleaseOrMint(&_MockLBTCTokenPool.TransactOpts, releaseOrMintIn)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.RemoveRemotePool(&_MockLBTCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.RemoveRemotePool(&_MockLBTCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.SetChainRateLimiterConfig(&_MockLBTCTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.SetChainRateLimiterConfig(&_MockLBTCTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.SetRateLimitAdmin(&_MockLBTCTokenPool.TransactOpts, rateLimitAdmin)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.SetRateLimitAdmin(&_MockLBTCTokenPool.TransactOpts, rateLimitAdmin)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.SetRouter(&_MockLBTCTokenPool.TransactOpts, newRouter)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.SetRouter(&_MockLBTCTokenPool.TransactOpts, newRouter)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.TransferOwnership(&_MockLBTCTokenPool.TransactOpts, to)
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MockLBTCTokenPool.Contract.TransferOwnership(&_MockLBTCTokenPool.TransactOpts, to)
}

type MockLBTCTokenPoolAllowListAddIterator struct {
	Event *MockLBTCTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolAllowListAdd)
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
		it.Event = new(MockLBTCTokenPoolAllowListAdd)
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

func (it *MockLBTCTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*MockLBTCTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolAllowListAddIterator{contract: _MockLBTCTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolAllowListAdd)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*MockLBTCTokenPoolAllowListAdd, error) {
	event := new(MockLBTCTokenPoolAllowListAdd)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolAllowListRemoveIterator struct {
	Event *MockLBTCTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolAllowListRemove)
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
		it.Event = new(MockLBTCTokenPoolAllowListRemove)
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

func (it *MockLBTCTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*MockLBTCTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolAllowListRemoveIterator{contract: _MockLBTCTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolAllowListRemove)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*MockLBTCTokenPoolAllowListRemove, error) {
	event := new(MockLBTCTokenPoolAllowListRemove)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolBurnedIterator struct {
	Event *MockLBTCTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolBurned)
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
		it.Event = new(MockLBTCTokenPoolBurned)
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

func (it *MockLBTCTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*MockLBTCTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolBurnedIterator{contract: _MockLBTCTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolBurned)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseBurned(log types.Log) (*MockLBTCTokenPoolBurned, error) {
	event := new(MockLBTCTokenPoolBurned)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolChainAddedIterator struct {
	Event *MockLBTCTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolChainAdded)
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
		it.Event = new(MockLBTCTokenPoolChainAdded)
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

func (it *MockLBTCTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*MockLBTCTokenPoolChainAddedIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolChainAddedIterator{contract: _MockLBTCTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolChainAdded)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseChainAdded(log types.Log) (*MockLBTCTokenPoolChainAdded, error) {
	event := new(MockLBTCTokenPoolChainAdded)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolChainConfiguredIterator struct {
	Event *MockLBTCTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolChainConfigured)
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
		it.Event = new(MockLBTCTokenPoolChainConfigured)
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

func (it *MockLBTCTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*MockLBTCTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolChainConfiguredIterator{contract: _MockLBTCTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolChainConfigured)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseChainConfigured(log types.Log) (*MockLBTCTokenPoolChainConfigured, error) {
	event := new(MockLBTCTokenPoolChainConfigured)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolChainRemovedIterator struct {
	Event *MockLBTCTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolChainRemoved)
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
		it.Event = new(MockLBTCTokenPoolChainRemoved)
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

func (it *MockLBTCTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*MockLBTCTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolChainRemovedIterator{contract: _MockLBTCTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolChainRemoved)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseChainRemoved(log types.Log) (*MockLBTCTokenPoolChainRemoved, error) {
	event := new(MockLBTCTokenPoolChainRemoved)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolConfigChangedIterator struct {
	Event *MockLBTCTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolConfigChanged)
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
		it.Event = new(MockLBTCTokenPoolConfigChanged)
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

func (it *MockLBTCTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*MockLBTCTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolConfigChangedIterator{contract: _MockLBTCTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolConfigChanged)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseConfigChanged(log types.Log) (*MockLBTCTokenPoolConfigChanged, error) {
	event := new(MockLBTCTokenPoolConfigChanged)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolLockedIterator struct {
	Event *MockLBTCTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolLocked)
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
		it.Event = new(MockLBTCTokenPoolLocked)
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

func (it *MockLBTCTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*MockLBTCTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolLockedIterator{contract: _MockLBTCTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolLocked)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseLocked(log types.Log) (*MockLBTCTokenPoolLocked, error) {
	event := new(MockLBTCTokenPoolLocked)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolMintedIterator struct {
	Event *MockLBTCTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolMinted)
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
		it.Event = new(MockLBTCTokenPoolMinted)
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

func (it *MockLBTCTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*MockLBTCTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolMintedIterator{contract: _MockLBTCTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolMinted)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseMinted(log types.Log) (*MockLBTCTokenPoolMinted, error) {
	event := new(MockLBTCTokenPoolMinted)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolOwnershipTransferRequestedIterator struct {
	Event *MockLBTCTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolOwnershipTransferRequested)
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
		it.Event = new(MockLBTCTokenPoolOwnershipTransferRequested)
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

func (it *MockLBTCTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockLBTCTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolOwnershipTransferRequestedIterator{contract: _MockLBTCTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolOwnershipTransferRequested)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*MockLBTCTokenPoolOwnershipTransferRequested, error) {
	event := new(MockLBTCTokenPoolOwnershipTransferRequested)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolOwnershipTransferredIterator struct {
	Event *MockLBTCTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolOwnershipTransferred)
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
		it.Event = new(MockLBTCTokenPoolOwnershipTransferred)
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

func (it *MockLBTCTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockLBTCTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolOwnershipTransferredIterator{contract: _MockLBTCTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolOwnershipTransferred)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*MockLBTCTokenPoolOwnershipTransferred, error) {
	event := new(MockLBTCTokenPoolOwnershipTransferred)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolRateLimitAdminSetIterator struct {
	Event *MockLBTCTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolRateLimitAdminSet)
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
		it.Event = new(MockLBTCTokenPoolRateLimitAdminSet)
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

func (it *MockLBTCTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*MockLBTCTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolRateLimitAdminSetIterator{contract: _MockLBTCTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolRateLimitAdminSet)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*MockLBTCTokenPoolRateLimitAdminSet, error) {
	event := new(MockLBTCTokenPoolRateLimitAdminSet)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolReleasedIterator struct {
	Event *MockLBTCTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolReleased)
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
		it.Event = new(MockLBTCTokenPoolReleased)
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

func (it *MockLBTCTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*MockLBTCTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolReleasedIterator{contract: _MockLBTCTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolReleased)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseReleased(log types.Log) (*MockLBTCTokenPoolReleased, error) {
	event := new(MockLBTCTokenPoolReleased)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolRemotePoolAddedIterator struct {
	Event *MockLBTCTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolRemotePoolAdded)
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
		it.Event = new(MockLBTCTokenPoolRemotePoolAdded)
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

func (it *MockLBTCTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*MockLBTCTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolRemotePoolAddedIterator{contract: _MockLBTCTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolRemotePoolAdded)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*MockLBTCTokenPoolRemotePoolAdded, error) {
	event := new(MockLBTCTokenPoolRemotePoolAdded)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolRemotePoolRemovedIterator struct {
	Event *MockLBTCTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolRemotePoolRemoved)
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
		it.Event = new(MockLBTCTokenPoolRemotePoolRemoved)
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

func (it *MockLBTCTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*MockLBTCTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolRemotePoolRemovedIterator{contract: _MockLBTCTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolRemotePoolRemoved)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*MockLBTCTokenPoolRemotePoolRemoved, error) {
	event := new(MockLBTCTokenPoolRemotePoolRemoved)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockLBTCTokenPoolRouterUpdatedIterator struct {
	Event *MockLBTCTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockLBTCTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockLBTCTokenPoolRouterUpdated)
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
		it.Event = new(MockLBTCTokenPoolRouterUpdated)
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

func (it *MockLBTCTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *MockLBTCTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockLBTCTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*MockLBTCTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &MockLBTCTokenPoolRouterUpdatedIterator{contract: _MockLBTCTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _MockLBTCTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockLBTCTokenPoolRouterUpdated)
				if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_MockLBTCTokenPool *MockLBTCTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*MockLBTCTokenPoolRouterUpdated, error) {
	event := new(MockLBTCTokenPoolRouterUpdated)
	if err := _MockLBTCTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MockLBTCTokenPool *MockLBTCTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockLBTCTokenPool.abi.Events["AllowListAdd"].ID:
		return _MockLBTCTokenPool.ParseAllowListAdd(log)
	case _MockLBTCTokenPool.abi.Events["AllowListRemove"].ID:
		return _MockLBTCTokenPool.ParseAllowListRemove(log)
	case _MockLBTCTokenPool.abi.Events["Burned"].ID:
		return _MockLBTCTokenPool.ParseBurned(log)
	case _MockLBTCTokenPool.abi.Events["ChainAdded"].ID:
		return _MockLBTCTokenPool.ParseChainAdded(log)
	case _MockLBTCTokenPool.abi.Events["ChainConfigured"].ID:
		return _MockLBTCTokenPool.ParseChainConfigured(log)
	case _MockLBTCTokenPool.abi.Events["ChainRemoved"].ID:
		return _MockLBTCTokenPool.ParseChainRemoved(log)
	case _MockLBTCTokenPool.abi.Events["ConfigChanged"].ID:
		return _MockLBTCTokenPool.ParseConfigChanged(log)
	case _MockLBTCTokenPool.abi.Events["Locked"].ID:
		return _MockLBTCTokenPool.ParseLocked(log)
	case _MockLBTCTokenPool.abi.Events["Minted"].ID:
		return _MockLBTCTokenPool.ParseMinted(log)
	case _MockLBTCTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _MockLBTCTokenPool.ParseOwnershipTransferRequested(log)
	case _MockLBTCTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _MockLBTCTokenPool.ParseOwnershipTransferred(log)
	case _MockLBTCTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _MockLBTCTokenPool.ParseRateLimitAdminSet(log)
	case _MockLBTCTokenPool.abi.Events["Released"].ID:
		return _MockLBTCTokenPool.ParseReleased(log)
	case _MockLBTCTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _MockLBTCTokenPool.ParseRemotePoolAdded(log)
	case _MockLBTCTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _MockLBTCTokenPool.ParseRemotePoolRemoved(log)
	case _MockLBTCTokenPool.abi.Events["RouterUpdated"].ID:
		return _MockLBTCTokenPool.ParseRouterUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockLBTCTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (MockLBTCTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (MockLBTCTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (MockLBTCTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (MockLBTCTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (MockLBTCTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (MockLBTCTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (MockLBTCTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (MockLBTCTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (MockLBTCTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MockLBTCTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MockLBTCTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (MockLBTCTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (MockLBTCTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (MockLBTCTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (MockLBTCTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (_MockLBTCTokenPool *MockLBTCTokenPool) Address() common.Address {
	return _MockLBTCTokenPool.address
}

type MockLBTCTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	GetTokenDecimals(opts *bind.CallOpts) (uint8, error)

	IDestPoolData(opts *bind.CallOpts) ([]byte, error)

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

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*MockLBTCTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*MockLBTCTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*MockLBTCTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*MockLBTCTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*MockLBTCTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*MockLBTCTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*MockLBTCTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*MockLBTCTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*MockLBTCTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*MockLBTCTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*MockLBTCTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*MockLBTCTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*MockLBTCTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*MockLBTCTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*MockLBTCTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*MockLBTCTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*MockLBTCTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*MockLBTCTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockLBTCTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MockLBTCTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockLBTCTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MockLBTCTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*MockLBTCTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*MockLBTCTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*MockLBTCTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*MockLBTCTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*MockLBTCTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*MockLBTCTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*MockLBTCTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*MockLBTCTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*MockLBTCTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *MockLBTCTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*MockLBTCTokenPoolRouterUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
