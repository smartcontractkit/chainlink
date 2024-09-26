package testhelpers_1_4_0

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	burn_mint_token_pool "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_mint_token_pool_1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_2_0"
	evm_2_evm_offramp "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_2_0"
	evm_2_evm_onramp "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool_1_0_0"
	lock_release_token_pool "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool_1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

var (
	// Source
	SourcePool          = "source Link pool"
	SourcePriceRegistry = "source PriceRegistry"
	OnRamp              = "onramp"
	OnRampNative        = "onramp-native"
	SourceRouter        = "source router"

	// Dest
	OffRamp  = "offramp"
	DestPool = "dest Link pool"

	Receiver            = "receiver"
	Sender              = "sender"
	Link                = func(amount int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(amount)) }
	HundredLink         = Link(100)
	LinkUSDValue        = func(amount int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(amount)) }
	SourceChainID       = uint64(1000)
	SourceChainSelector = uint64(11787463284727550157)
	DestChainID         = uint64(1337)
	DestChainSelector   = uint64(3379446385462418246)
)

// Backwards compat, in principle these statuses are version dependent
// TODO: Adjust integration tests to be version agnostic using readers
var (
	ExecutionStateSuccess = MessageExecutionState(cciptypes.ExecutionStateSuccess)
	ExecutionStateFailure = MessageExecutionState(cciptypes.ExecutionStateFailure)
)

type MessageExecutionState cciptypes.MessageExecutionState
type CommitOffchainConfig struct {
	v1_2_0.JSONCommitOffchainConfig
}

func (c CommitOffchainConfig) Encode() ([]byte, error) {
	return ccipconfig.EncodeOffchainConfig(c.JSONCommitOffchainConfig)
}

func NewCommitOffchainConfig(
	GasPriceHeartBeat config.Duration,
	DAGasPriceDeviationPPB uint32,
	ExecGasPriceDeviationPPB uint32,
	TokenPriceHeartBeat config.Duration,
	TokenPriceDeviationPPB uint32,
	InflightCacheExpiry config.Duration,
	priceReportingDisabled bool) CommitOffchainConfig {
	return CommitOffchainConfig{v1_2_0.JSONCommitOffchainConfig{
		GasPriceHeartBeat:        GasPriceHeartBeat,
		DAGasPriceDeviationPPB:   DAGasPriceDeviationPPB,
		ExecGasPriceDeviationPPB: ExecGasPriceDeviationPPB,
		TokenPriceHeartBeat:      TokenPriceHeartBeat,
		TokenPriceDeviationPPB:   TokenPriceDeviationPPB,
		InflightCacheExpiry:      InflightCacheExpiry,
		PriceReportingDisabled:   priceReportingDisabled,
	}}
}

type CommitOnchainConfig struct {
	ccipdata.CommitOnchainConfig
}

func NewCommitOnchainConfig(
	PriceRegistry common.Address,
) CommitOnchainConfig {
	return CommitOnchainConfig{ccipdata.CommitOnchainConfig{
		PriceRegistry: PriceRegistry,
	}}
}

type ExecOnchainConfig struct {
	v1_2_0.ExecOnchainConfig
}

func NewExecOnchainConfig(
	PermissionLessExecutionThresholdSeconds uint32,
	Router common.Address,
	PriceRegistry common.Address,
	MaxNumberOfTokensPerMsg uint16,
	MaxDataBytes uint32,
	MaxPoolReleaseOrMintGas uint32,
) ExecOnchainConfig {
	return ExecOnchainConfig{v1_2_0.ExecOnchainConfig{
		PermissionLessExecutionThresholdSeconds: PermissionLessExecutionThresholdSeconds,
		Router:                                  Router,
		PriceRegistry:                           PriceRegistry,
		MaxNumberOfTokensPerMsg:                 MaxNumberOfTokensPerMsg,
		MaxDataBytes:                            MaxDataBytes,
		MaxPoolReleaseOrMintGas:                 MaxPoolReleaseOrMintGas,
	}}
}

type ExecOffchainConfig struct {
	v1_2_0.JSONExecOffchainConfig
}

func (c ExecOffchainConfig) Encode() ([]byte, error) {
	return ccipconfig.EncodeOffchainConfig(c.JSONExecOffchainConfig)
}

func NewExecOffchainConfig(
	DestOptimisticConfirmations uint32,
	BatchGasLimit uint32,
	RelativeBoostPerWaitHour float64,
	InflightCacheExpiry config.Duration,
	RootSnoozeTime config.Duration,
	BatchingStrategyID uint32,
) ExecOffchainConfig {
	return ExecOffchainConfig{v1_2_0.JSONExecOffchainConfig{
		DestOptimisticConfirmations: DestOptimisticConfirmations,
		BatchGasLimit:               BatchGasLimit,
		RelativeBoostPerWaitHour:    RelativeBoostPerWaitHour,
		InflightCacheExpiry:         InflightCacheExpiry,
		RootSnoozeTime:              RootSnoozeTime,
		BatchingStrategyID:          BatchingStrategyID,
	}}
}

type MaybeRevertReceiver struct {
	Receiver *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	Strict   bool
}

type Common struct {
	ChainID           uint64
	ChainSelector     uint64
	User              *bind.TransactOpts
	Chain             *backends.SimulatedBackend
	LinkToken         *link_token_interface.LinkToken
	LinkTokenPool     *lock_release_token_pool.LockReleaseTokenPool
	CustomToken       *link_token_interface.LinkToken
	WrappedNative     *weth9.WETH9
	WrappedNativePool *lock_release_token_pool_1_0_0.LockReleaseTokenPool
	ARM               *mock_rmn_contract.MockRMNContract
	ARMProxy          *rmn_proxy_contract.RMNProxyContract
	PriceRegistry     *price_registry_1_2_0.PriceRegistry
}

type SourceChain struct {
	Common
	Router *router.Router
	OnRamp *evm_2_evm_onramp.EVM2EVMOnRamp
}

type DestinationChain struct {
	Common

	CommitStore *commit_store_1_2_0.CommitStore
	Router      *router.Router
	OffRamp     *evm_2_evm_offramp.EVM2EVMOffRamp
	Receivers   []MaybeRevertReceiver
}

type OCR2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type BalanceAssertion struct {
	Name     string
	Address  common.Address
	Expected string
	Getter   func(t *testing.T, addr common.Address) *big.Int
	Within   string
}

type BalanceReq struct {
	Name   string
	Addr   common.Address
	Getter func(t *testing.T, addr common.Address) *big.Int
}

type CCIPContracts struct {
	Source  SourceChain
	Dest    DestinationChain
	Oracles []confighelper.OracleIdentityExtra

	commitOCRConfig, execOCRConfig *OCR2Config
}

func (c *CCIPContracts) DeployNewOffRamp(t *testing.T) {
	prevOffRamp := common.HexToAddress("")
	if c.Dest.OffRamp != nil {
		prevOffRamp = c.Dest.OffRamp.Address()
	}
	offRampAddress, _, _, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
		c.Dest.User,
		c.Dest.Chain,
		evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore:         c.Dest.CommitStore.Address(),
			ChainSelector:       c.Dest.ChainSelector,
			SourceChainSelector: c.Source.ChainSelector,
			OnRamp:              c.Source.OnRamp.Address(),
			PrevOffRamp:         prevOffRamp,
			ArmProxy:            c.Dest.ARMProxy.Address(),
		},
		[]common.Address{c.Source.LinkToken.Address()},   // source tokens
		[]common.Address{c.Dest.LinkTokenPool.Address()}, // pools
		evm_2_evm_offramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  LinkUSDValue(100),
			Rate:      LinkUSDValue(1),
		},
	)
	require.NoError(t, err)
	c.Dest.Chain.Commit()

	c.Dest.OffRamp, err = evm_2_evm_offramp.NewEVM2EVMOffRamp(offRampAddress, c.Dest.Chain)
	require.NoError(t, err)

	c.Dest.Chain.Commit()
	c.Source.Chain.Commit()
}

func (c *CCIPContracts) EnableOffRamp(t *testing.T) {
	_, err := c.Dest.Router.ApplyRampUpdates(c.Dest.User, nil, nil, []router.RouterOffRamp{{SourceChainSelector: SourceChainSelector, OffRamp: c.Dest.OffRamp.Address()}})
	require.NoError(t, err)
	c.Dest.Chain.Commit()

	onChainConfig := c.CreateDefaultExecOnchainConfig(t)
	offChainConfig := c.CreateDefaultExecOffchainConfig(t)

	c.SetupExecOCR2Config(t, onChainConfig, offChainConfig)
}

func (c *CCIPContracts) EnableCommitStore(t *testing.T) {
	onChainConfig := c.CreateDefaultCommitOnchainConfig(t)
	offChainConfig := c.CreateDefaultCommitOffchainConfig(t)

	c.SetupCommitOCR2Config(t, onChainConfig, offChainConfig)

	_, err := c.Dest.PriceRegistry.ApplyPriceUpdatersUpdates(c.Dest.User, []common.Address{c.Dest.CommitStore.Address()}, []common.Address{})
	require.NoError(t, err)
	c.Dest.Chain.Commit()
}

func (c *CCIPContracts) DeployNewOnRamp(t *testing.T) {
	t.Log("Deploying new onRamp")
	// find the last onRamp
	prevOnRamp := common.HexToAddress("")
	if c.Source.OnRamp != nil {
		prevOnRamp = c.Source.OnRamp.Address()
	}
	onRampAddress, _, _, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		c.Source.User,  // user
		c.Source.Chain, // client
		evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			LinkToken:         c.Source.LinkToken.Address(),
			ChainSelector:     c.Source.ChainSelector,
			DestChainSelector: c.Dest.ChainSelector,
			DefaultTxGasLimit: 200_000,
			MaxNopFeesJuels:   big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
			PrevOnRamp:        prevOnRamp,
			ArmProxy:          c.Source.ARM.Address(), // ARM
		},
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                            c.Source.Router.Address(),
			MaxNumberOfTokensPerMsg:           5,
			DestGasOverhead:                   350_000,
			DestGasPerPayloadByte:             16,
			DestDataAvailabilityOverheadGas:   33_596,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 6840, // 0.684
			PriceRegistry:                     c.Source.PriceRegistry.Address(),
			MaxDataBytes:                      1e5,
			MaxPerMsgGasLimit:                 4_000_000,
		},
		[]evm_2_evm_onramp.InternalPoolUpdate{
			{
				Token: c.Source.LinkToken.Address(),
				Pool:  c.Source.LinkTokenPool.Address(),
			},
		},
		evm_2_evm_onramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  LinkUSDValue(100),
			Rate:      LinkUSDValue(1),
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
			{
				Token:                      c.Source.LinkToken.Address(),
				NetworkFeeUSDCents:         1_00,
				GasMultiplierWeiPerEth:     1e18,
				PremiumMultiplierWeiPerEth: 9e17,
				Enabled:                    true,
			},
			{
				Token:                      c.Source.WrappedNative.Address(),
				NetworkFeeUSDCents:         1_00,
				GasMultiplierWeiPerEth:     1e18,
				PremiumMultiplierWeiPerEth: 1e18,
				Enabled:                    true,
			},
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
			{
				Token:             c.Source.LinkToken.Address(),
				MinFeeUSDCents:    50,           // $0.5
				MaxFeeUSDCents:    1_000_000_00, // $ 1 million
				DeciBps:           5_0,          // 5 bps
				DestGasOverhead:   34_000,
				DestBytesOverhead: 32,
			},
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
	)

	require.NoError(t, err)
	c.Source.Chain.Commit()
	c.Dest.Chain.Commit()
	c.Source.OnRamp, err = evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampAddress, c.Source.Chain)
	require.NoError(t, err)
	c.Source.Chain.Commit()
	c.Dest.Chain.Commit()
}

func (c *CCIPContracts) EnableOnRamp(t *testing.T) {
	t.Log("Setting onRamp on source router")
	_, err := c.Source.Router.ApplyRampUpdates(c.Source.User, []router.RouterOnRamp{{DestChainSelector: c.Dest.ChainSelector, OnRamp: c.Source.OnRamp.Address()}}, nil, nil)
	require.NoError(t, err)
	c.Source.Chain.Commit()
	c.Dest.Chain.Commit()
}

func (c *CCIPContracts) DeployNewCommitStore(t *testing.T) {
	commitStoreAddress, _, _, err := commit_store_1_2_0.DeployCommitStore(
		c.Dest.User,  // user
		c.Dest.Chain, // client
		commit_store_1_2_0.CommitStoreStaticConfig{
			ChainSelector:       c.Dest.ChainSelector,
			SourceChainSelector: c.Source.ChainSelector,
			OnRamp:              c.Source.OnRamp.Address(),
			ArmProxy:            c.Dest.ARMProxy.Address(),
		},
	)
	require.NoError(t, err)
	c.Dest.Chain.Commit()
	// since CommitStoreHelper derives from CommitStore, it's safe to instantiate both on same address
	c.Dest.CommitStore, err = commit_store_1_2_0.NewCommitStore(commitStoreAddress, c.Dest.Chain)
	require.NoError(t, err)
}

func (c *CCIPContracts) DeployNewPriceRegistry(t *testing.T) {
	t.Log("Deploying new Price Registry")
	destPricesAddress, _, _, err := price_registry_1_2_0.DeployPriceRegistry(
		c.Dest.User,
		c.Dest.Chain,
		[]common.Address{c.Dest.CommitStore.Address()},
		[]common.Address{c.Dest.LinkToken.Address()},
		60*60*24*14, // two weeks
	)
	require.NoError(t, err)
	c.Source.Chain.Commit()
	c.Dest.Chain.Commit()
	c.Dest.PriceRegistry, err = price_registry_1_2_0.NewPriceRegistry(destPricesAddress, c.Dest.Chain)
	require.NoError(t, err)

	priceUpdates := price_registry_1_2_0.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry_1_2_0.InternalTokenPriceUpdate{
			{
				SourceToken: c.Dest.LinkToken.Address(),
				UsdPerToken: big.NewInt(8e18), // 8usd
			},
			{
				SourceToken: c.Dest.WrappedNative.Address(),
				UsdPerToken: big.NewInt(1e18), // 1usd
			},
		},
		GasPriceUpdates: []price_registry_1_2_0.InternalGasPriceUpdate{
			{
				DestChainSelector: c.Source.ChainSelector,
				UsdPerUnitGas:     big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
			},
		},
	}
	_, err = c.Dest.PriceRegistry.UpdatePrices(c.Dest.User, priceUpdates)
	require.NoError(t, err)

	c.Source.Chain.Commit()
	c.Dest.Chain.Commit()

	t.Logf("New Price Registry deployed at %s", destPricesAddress.String())
}

func (c *CCIPContracts) SetNopsOnRamp(t *testing.T, nopsAndWeights []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight) {
	tx, err := c.Source.OnRamp.SetNops(c.Source.User, nopsAndWeights)
	require.NoError(t, err)
	c.Source.Chain.Commit()
	_, err = bind.WaitMined(context.Background(), c.Source.Chain, tx)
	require.NoError(t, err)
}

func (c *CCIPContracts) GetSourceLinkBalance(t *testing.T, addr common.Address) *big.Int {
	return GetBalance(t, c.Source.Chain, c.Source.LinkToken.Address(), addr)
}

func (c *CCIPContracts) GetDestLinkBalance(t *testing.T, addr common.Address) *big.Int {
	return GetBalance(t, c.Dest.Chain, c.Dest.LinkToken.Address(), addr)
}

func (c *CCIPContracts) GetSourceWrappedTokenBalance(t *testing.T, addr common.Address) *big.Int {
	return GetBalance(t, c.Source.Chain, c.Source.WrappedNative.Address(), addr)
}

func (c *CCIPContracts) GetDestWrappedTokenBalance(t *testing.T, addr common.Address) *big.Int {
	return GetBalance(t, c.Dest.Chain, c.Dest.WrappedNative.Address(), addr)
}

func (c *CCIPContracts) AssertBalances(t *testing.T, bas []BalanceAssertion) {
	for _, b := range bas {
		actual := b.Getter(t, b.Address)
		t.Log("Checking balance for", b.Name, "at", b.Address.Hex(), "got", actual)
		require.NotNil(t, actual, "%v getter return nil", b.Name)
		if b.Within == "" {
			require.Equal(t, b.Expected, actual.String(), "wrong balance for %s got %s want %s", b.Name, actual, b.Expected)
		} else {
			bi, _ := big.NewInt(0).SetString(b.Expected, 10)
			withinI, _ := big.NewInt(0).SetString(b.Within, 10)
			high := big.NewInt(0).Add(bi, withinI)
			low := big.NewInt(0).Sub(bi, withinI)
			require.Equal(t, -1, actual.Cmp(high), "wrong balance for %s got %s outside expected range [%s, %s]", b.Name, actual, low, high)
			require.Equal(t, 1, actual.Cmp(low), "wrong balance for %s got %s outside expected range [%s, %s]", b.Name, actual, low, high)
		}
	}
}

func AccountToAddress(accounts []ocr2types.Account) (addresses []common.Address, err error) {
	for _, signer := range accounts {
		bytes, err := hexutil.Decode(string(signer))
		if err != nil {
			return []common.Address{}, errors.Wrap(err, fmt.Sprintf("given address is not valid %s", signer))
		}
		if len(bytes) != 20 {
			return []common.Address{}, errors.Errorf("address is not 20 bytes %s", signer)
		}
		addresses = append(addresses, common.BytesToAddress(bytes))
	}
	return addresses, nil
}

func OnchainPublicKeyToAddress(publicKeys []ocrtypes.OnchainPublicKey) (addresses []common.Address, err error) {
	for _, signer := range publicKeys {
		if len(signer) != 20 {
			return []common.Address{}, errors.Errorf("address is not 20 bytes %s", signer)
		}
		addresses = append(addresses, common.BytesToAddress(signer))
	}
	return addresses, nil
}

func (c *CCIPContracts) DeriveOCR2Config(t *testing.T, oracles []confighelper.OracleIdentityExtra, rawOnchainConfig []byte, rawOffchainConfig []byte) *OCR2Config {
	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress
		1*time.Second,        // deltaResend
		1*time.Second,        // deltaRound
		500*time.Millisecond, // deltaGrace
		2*time.Second,        // deltaStage
		3,
		[]int{1, 1, 1, 1},
		oracles,
		rawOffchainConfig,
		50*time.Millisecond, // Max duration query
		1*time.Second,       // Max duration observation
		100*time.Millisecond,
		100*time.Millisecond,
		100*time.Millisecond,
		1, // faults
		rawOnchainConfig,
	)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	lggr.Infow("Setting Config on Oracle Contract",
		"signers", signers,
		"transmitters", transmitters,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", offchainConfigVersion,
	)
	signerAddresses, err := OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := AccountToAddress(transmitters)
	require.NoError(t, err)

	return &OCR2Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     threshold,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func (c *CCIPContracts) SetupCommitOCR2Config(t *testing.T, commitOnchainConfig, commitOffchainConfig []byte) {
	c.commitOCRConfig = c.DeriveOCR2Config(t, c.Oracles, commitOnchainConfig, commitOffchainConfig)
	// Set the DON on the commit store
	_, err := c.Dest.CommitStore.SetOCR2Config(
		c.Dest.User,
		c.commitOCRConfig.Signers,
		c.commitOCRConfig.Transmitters,
		c.commitOCRConfig.F,
		c.commitOCRConfig.OnchainConfig,
		c.commitOCRConfig.OffchainConfigVersion,
		c.commitOCRConfig.OffchainConfig,
	)
	require.NoError(t, err)
	c.Dest.Chain.Commit()
}

func (c *CCIPContracts) SetupExecOCR2Config(t *testing.T, execOnchainConfig, execOffchainConfig []byte) {
	c.execOCRConfig = c.DeriveOCR2Config(t, c.Oracles, execOnchainConfig, execOffchainConfig)
	// Same DON on the offramp
	_, err := c.Dest.OffRamp.SetOCR2Config(
		c.Dest.User,
		c.execOCRConfig.Signers,
		c.execOCRConfig.Transmitters,
		c.execOCRConfig.F,
		c.execOCRConfig.OnchainConfig,
		c.execOCRConfig.OffchainConfigVersion,
		c.execOCRConfig.OffchainConfig,
	)
	require.NoError(t, err)
	c.Dest.Chain.Commit()
}

func (c *CCIPContracts) SetupOnchainConfig(t *testing.T, commitOnchainConfig, commitOffchainConfig, execOnchainConfig, execOffchainConfig []byte) int64 {
	// Note We do NOT set the payees, payment is done in the OCR2Base implementation
	blockBeforeConfig, err := c.Dest.Chain.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)

	c.SetupCommitOCR2Config(t, commitOnchainConfig, commitOffchainConfig)
	c.SetupExecOCR2Config(t, execOnchainConfig, execOffchainConfig)

	return blockBeforeConfig.Number().Int64()
}

func (c *CCIPContracts) SetupLockAndMintTokenPool(
	sourceTokenAddress common.Address,
	wrappedTokenName,
	wrappedTokenSymbol string) (common.Address, *burn_mint_erc677.BurnMintERC677, error) {
	// Deploy dest token & pool
	destTokenAddress, _, _, err := burn_mint_erc677.DeployBurnMintERC677(c.Dest.User, c.Dest.Chain, wrappedTokenName, wrappedTokenSymbol, 18, big.NewInt(0))
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Dest.Chain.Commit()

	destToken, err := burn_mint_erc677.NewBurnMintERC677(destTokenAddress, c.Dest.Chain)
	if err != nil {
		return [20]byte{}, nil, err
	}

	destPoolAddress, _, destPool, err := burn_mint_token_pool.DeployBurnMintTokenPool(
		c.Dest.User,
		c.Dest.Chain,
		destTokenAddress,
		[]common.Address{}, // pool originalSender allowList
		c.Dest.ARMProxy.Address(),
		c.Dest.Router.Address(),
	)
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Dest.Chain.Commit()

	_, err = destToken.GrantMintAndBurnRoles(c.Dest.User, destPoolAddress)
	if err != nil {
		return [20]byte{}, nil, err
	}

	_, err = destPool.ApplyChainUpdates(c.Dest.User,
		[]burn_mint_token_pool.TokenPoolChainUpdate{
			{
				RemoteChainSelector: c.Source.ChainSelector,
				Allowed:             true,
				OutboundRateLimiterConfig: burn_mint_token_pool.RateLimiterConfig{
					IsEnabled: true,
					Capacity:  HundredLink,
					Rate:      big.NewInt(1e18),
				},
				InboundRateLimiterConfig: burn_mint_token_pool.RateLimiterConfig{
					IsEnabled: true,
					Capacity:  HundredLink,
					Rate:      big.NewInt(1e18),
				},
			},
		})
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Dest.Chain.Commit()

	sourcePoolAddress, _, sourcePool, err := lock_release_token_pool.DeployLockReleaseTokenPool(
		c.Source.User,
		c.Source.Chain,
		sourceTokenAddress,
		[]common.Address{}, // empty allowList at deploy time indicates pool has no original sender restrictions
		c.Source.ARMProxy.Address(),
		true,
		c.Source.Router.Address(),
	)
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Source.Chain.Commit()

	// set onRamp as valid caller for source pool
	_, err = sourcePool.ApplyChainUpdates(c.Source.User, []lock_release_token_pool.TokenPoolChainUpdate{
		{
			RemoteChainSelector: c.Dest.ChainSelector,
			Allowed:             true,
			OutboundRateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
			InboundRateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
		},
	})
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Source.Chain.Commit()

	wrappedNativeAddress, err := c.Source.Router.GetWrappedNative(nil)
	if err != nil {
		return [20]byte{}, nil, err
	}

	//native token is used as fee token
	_, err = c.Source.PriceRegistry.UpdatePrices(c.Source.User, price_registry_1_2_0.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry_1_2_0.InternalTokenPriceUpdate{
			{
				SourceToken: sourceTokenAddress,
				UsdPerToken: big.NewInt(5),
			},
		},
		GasPriceUpdates: []price_registry_1_2_0.InternalGasPriceUpdate{},
	})
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Source.Chain.Commit()

	_, err = c.Source.PriceRegistry.ApplyFeeTokensUpdates(c.Source.User, []common.Address{wrappedNativeAddress}, nil)
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Source.Chain.Commit()

	// add new token pool created above
	_, err = c.Source.OnRamp.ApplyPoolUpdates(c.Source.User, nil, []evm_2_evm_onramp.InternalPoolUpdate{
		{
			Token: sourceTokenAddress,
			Pool:  sourcePoolAddress,
		},
	})
	if err != nil {
		return [20]byte{}, nil, err
	}

	_, err = c.Dest.OffRamp.ApplyPoolUpdates(c.Dest.User, nil, []evm_2_evm_offramp.InternalPoolUpdate{
		{
			Token: sourceTokenAddress,
			Pool:  destPoolAddress,
		},
	})
	if err != nil {
		return [20]byte{}, nil, err
	}
	c.Dest.Chain.Commit()

	return sourcePoolAddress, destToken, err
}

func (c *CCIPContracts) SendMessage(t *testing.T, gasLimit, tokenAmount *big.Int, receiverAddr common.Address) {
	extraArgs, err := GetEVMExtraArgsV1(gasLimit, false)
	require.NoError(t, err)
	msg := router.ClientEVM2AnyMessage{
		Receiver: MustEncodeAddress(t, receiverAddr),
		Data:     []byte("hello"),
		TokenAmounts: []router.ClientEVMTokenAmount{
			{
				Token:  c.Source.LinkToken.Address(),
				Amount: tokenAmount,
			},
		},
		FeeToken:  c.Source.LinkToken.Address(),
		ExtraArgs: extraArgs,
	}
	fee, err := c.Source.Router.GetFee(nil, c.Dest.ChainSelector, msg)
	require.NoError(t, err)
	// Currently no overhead and 1gwei dest gas price. So fee is simply gasLimit * gasPrice.
	// require.Equal(t, new(big.Int).Mul(gasLimit, gasPrice).String(), fee.String())
	// Approve the fee amount + the token amount
	_, err = c.Source.LinkToken.Approve(c.Source.User, c.Source.Router.Address(), new(big.Int).Add(fee, tokenAmount))
	require.NoError(t, err)
	c.Source.Chain.Commit()
	c.SendRequest(t, msg)
}

func GetBalances(t *testing.T, brs []BalanceReq) (map[string]*big.Int, error) {
	m := make(map[string]*big.Int)
	for _, br := range brs {
		m[br.Name] = br.Getter(t, br.Addr)
		if m[br.Name] == nil {
			return nil, fmt.Errorf("%v getter return nil", br.Name)
		}
	}
	return m, nil
}

func MustAddBigInt(a *big.Int, b string) *big.Int {
	bi, _ := big.NewInt(0).SetString(b, 10)
	return big.NewInt(0).Add(a, bi)
}

func MustSubBigInt(a *big.Int, b string) *big.Int {
	bi, _ := big.NewInt(0).SetString(b, 10)
	return big.NewInt(0).Sub(a, bi)
}

func MustEncodeAddress(t *testing.T, address common.Address) []byte {
	bts, err := utils.ABIEncode(`[{"type":"address"}]`, address)
	require.NoError(t, err)
	return bts
}

func SetupCCIPContracts(t *testing.T, sourceChainID, sourceChainSelector, destChainID, destChainSelector uint64) CCIPContracts {
	sourceChain, sourceUser := testhelpers.SetupChain(t)
	destChain, destUser := testhelpers.SetupChain(t)

	armSourceAddress, _, _, err := mock_rmn_contract.DeployMockRMNContract(
		sourceUser,
		sourceChain,
	)
	require.NoError(t, err)
	sourceARM, err := mock_rmn_contract.NewMockRMNContract(armSourceAddress, sourceChain)
	require.NoError(t, err)
	armProxySourceAddress, _, _, err := rmn_proxy_contract.DeployRMNProxyContract(
		sourceUser,
		sourceChain,
		armSourceAddress,
	)
	require.NoError(t, err)
	sourceARMProxy, err := rmn_proxy_contract.NewRMNProxyContract(armProxySourceAddress, sourceChain)
	require.NoError(t, err)
	sourceChain.Commit()

	armDestAddress, _, _, err := mock_rmn_contract.DeployMockRMNContract(
		destUser,
		destChain,
	)
	require.NoError(t, err)
	armProxyDestAddress, _, _, err := rmn_proxy_contract.DeployRMNProxyContract(
		destUser,
		destChain,
		armDestAddress,
	)
	require.NoError(t, err)
	destChain.Commit()
	destARM, err := mock_rmn_contract.NewMockRMNContract(armDestAddress, destChain)
	require.NoError(t, err)
	destARMProxy, err := rmn_proxy_contract.NewRMNProxyContract(armProxyDestAddress, destChain)
	require.NoError(t, err)

	// Deploy link token and pool on source chain
	sourceLinkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(sourceUser, sourceChain)
	require.NoError(t, err)
	sourceChain.Commit()
	sourceLinkToken, err := link_token_interface.NewLinkToken(sourceLinkTokenAddress, sourceChain)
	require.NoError(t, err)

	// Create router
	sourceWeth9addr, _, _, err := weth9.DeployWETH9(sourceUser, sourceChain)
	require.NoError(t, err)
	sourceWrapped, err := weth9.NewWETH9(sourceWeth9addr, sourceChain)
	require.NoError(t, err)

	sourceRouterAddress, _, _, err := router.DeployRouter(sourceUser, sourceChain, sourceWeth9addr, armProxySourceAddress)
	require.NoError(t, err)
	sourceRouter, err := router.NewRouter(sourceRouterAddress, sourceChain)
	require.NoError(t, err)
	sourceChain.Commit()

	sourceWeth9PoolAddress, _, _, err := lock_release_token_pool_1_0_0.DeployLockReleaseTokenPool(
		sourceUser,
		sourceChain,
		sourceWeth9addr,
		[]common.Address{},
		armProxySourceAddress,
	)
	require.NoError(t, err)
	sourceChain.Commit()

	sourceWeth9Pool, err := lock_release_token_pool_1_0_0.NewLockReleaseTokenPool(sourceWeth9PoolAddress, sourceChain)
	require.NoError(t, err)

	sourcePoolAddress, _, _, err := lock_release_token_pool.DeployLockReleaseTokenPool(
		sourceUser,
		sourceChain,
		sourceLinkTokenAddress,
		[]common.Address{},
		armProxySourceAddress,
		true,
		sourceRouterAddress,
	)
	require.NoError(t, err)
	sourceChain.Commit()
	sourcePool, err := lock_release_token_pool.NewLockReleaseTokenPool(sourcePoolAddress, sourceChain)
	require.NoError(t, err)

	// Deploy custom token pool source
	sourceCustomTokenAddress, _, _, err := link_token_interface.DeployLinkToken(sourceUser, sourceChain) // Just re-use this, it's an ERC20.
	require.NoError(t, err)
	sourceCustomToken, err := link_token_interface.NewLinkToken(sourceCustomTokenAddress, sourceChain)
	require.NoError(t, err)
	destChain.Commit()

	// Deploy custom token pool dest
	destCustomTokenAddress, _, _, err := link_token_interface.DeployLinkToken(destUser, destChain) // Just re-use this, it's an ERC20.
	require.NoError(t, err)
	destCustomToken, err := link_token_interface.NewLinkToken(destCustomTokenAddress, destChain)
	require.NoError(t, err)
	destChain.Commit()

	// Deploy and configure onramp
	sourcePricesAddress, _, _, err := price_registry_1_2_0.DeployPriceRegistry(
		sourceUser,
		sourceChain,
		nil,
		[]common.Address{sourceLinkTokenAddress, sourceWeth9addr},
		60*60*24*14, // two weeks
	)
	require.NoError(t, err)

	srcPriceRegistry, err := price_registry_1_2_0.NewPriceRegistry(sourcePricesAddress, sourceChain)
	require.NoError(t, err)

	_, err = srcPriceRegistry.UpdatePrices(sourceUser, price_registry_1_2_0.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry_1_2_0.InternalTokenPriceUpdate{
			{
				SourceToken: sourceLinkTokenAddress,
				UsdPerToken: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(20)),
			},
			{
				SourceToken: sourceWeth9addr,
				UsdPerToken: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2000)),
			},
		},
		GasPriceUpdates: []price_registry_1_2_0.InternalGasPriceUpdate{
			{
				DestChainSelector: destChainSelector,
				UsdPerUnitGas:     big.NewInt(20000e9),
			},
		},
	})
	require.NoError(t, err)

	onRampAddress, _, _, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		sourceUser,  // user
		sourceChain, // client
		evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			LinkToken:         sourceLinkTokenAddress,
			ChainSelector:     sourceChainSelector,
			DestChainSelector: destChainSelector,
			DefaultTxGasLimit: 200_000,
			MaxNopFeesJuels:   big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
			PrevOnRamp:        common.HexToAddress(""),
			ArmProxy:          armProxySourceAddress, // ARM
		},
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                            sourceRouterAddress,
			MaxNumberOfTokensPerMsg:           5,
			DestGasOverhead:                   350_000,
			DestGasPerPayloadByte:             16,
			DestDataAvailabilityOverheadGas:   33_596,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 6840, // 0.684
			PriceRegistry:                     sourcePricesAddress,
			MaxDataBytes:                      1e5,
			MaxPerMsgGasLimit:                 4_000_000,
		},
		[]evm_2_evm_onramp.InternalPoolUpdate{
			{
				Token: sourceLinkTokenAddress,
				Pool:  sourcePoolAddress,
			},
			{
				Token: sourceWeth9addr,
				Pool:  sourceWeth9PoolAddress,
			},
		},
		evm_2_evm_onramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  LinkUSDValue(100),
			Rate:      LinkUSDValue(1),
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
			{
				Token:                      sourceLinkTokenAddress,
				NetworkFeeUSDCents:         1_00,
				GasMultiplierWeiPerEth:     1e18,
				PremiumMultiplierWeiPerEth: 9e17,
				Enabled:                    true,
			},
			{
				Token:                      sourceWeth9addr,
				NetworkFeeUSDCents:         1_00,
				GasMultiplierWeiPerEth:     1e18,
				PremiumMultiplierWeiPerEth: 1e18,
				Enabled:                    true,
			},
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
			{
				Token:             sourceLinkTokenAddress,
				MinFeeUSDCents:    50,           // $0.5
				MaxFeeUSDCents:    1_000_000_00, // $ 1 million
				DeciBps:           5_0,          // 5 bps
				DestGasOverhead:   34_000,
				DestBytesOverhead: 32,
			},
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
	)
	require.NoError(t, err)
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampAddress, sourceChain)
	require.NoError(t, err)
	_, err = sourcePool.ApplyChainUpdates(
		sourceUser,
		[]lock_release_token_pool.TokenPoolChainUpdate{{
			RemoteChainSelector: DestChainSelector,
			Allowed:             true,
			OutboundRateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
			InboundRateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
		}},
	)
	require.NoError(t, err)
	_, err = sourceWeth9Pool.ApplyRampUpdates(sourceUser,
		[]lock_release_token_pool_1_0_0.TokenPoolRampUpdate{{Ramp: onRampAddress, Allowed: true,
			RateLimiterConfig: lock_release_token_pool_1_0_0.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
		}},
		[]lock_release_token_pool_1_0_0.TokenPoolRampUpdate{},
	)
	require.NoError(t, err)
	sourceChain.Commit()
	_, err = sourceRouter.ApplyRampUpdates(sourceUser, []router.RouterOnRamp{{DestChainSelector: destChainSelector, OnRamp: onRampAddress}}, nil, nil)
	require.NoError(t, err)
	sourceChain.Commit()

	destWethaddr, _, _, err := weth9.DeployWETH9(destUser, destChain)
	require.NoError(t, err)
	destWrapped, err := weth9.NewWETH9(destWethaddr, destChain)
	require.NoError(t, err)

	// Create dest router
	destRouterAddress, _, _, err := router.DeployRouter(destUser, destChain, destWethaddr, armProxyDestAddress)
	require.NoError(t, err)
	destChain.Commit()
	destRouter, err := router.NewRouter(destRouterAddress, destChain)
	require.NoError(t, err)

	// Deploy link token and pool on destination chain
	destLinkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(destUser, destChain)
	require.NoError(t, err)
	destChain.Commit()
	destLinkToken, err := link_token_interface.NewLinkToken(destLinkTokenAddress, destChain)
	require.NoError(t, err)
	destPoolAddress, _, _, err := lock_release_token_pool.DeployLockReleaseTokenPool(
		destUser,
		destChain,
		destLinkTokenAddress,
		[]common.Address{},
		armProxyDestAddress,
		true,
		destRouterAddress,
	)
	require.NoError(t, err)
	destChain.Commit()
	destPool, err := lock_release_token_pool.NewLockReleaseTokenPool(destPoolAddress, destChain)
	require.NoError(t, err)
	destChain.Commit()

	// Float the offramp pool
	o, err := destPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, destUser.From.String(), o.String())
	_, err = destPool.SetRebalancer(destUser, destUser.From)
	require.NoError(t, err)
	_, err = destLinkToken.Approve(destUser, destPoolAddress, Link(200))
	require.NoError(t, err)
	_, err = destPool.ProvideLiquidity(destUser, Link(200))
	require.NoError(t, err)
	destChain.Commit()

	destWrappedPoolAddress, _, _, err := lock_release_token_pool_1_0_0.DeployLockReleaseTokenPool(
		destUser,
		destChain,
		destWethaddr,
		[]common.Address{},
		armProxyDestAddress,
	)
	require.NoError(t, err)
	destWrappedPool, err := lock_release_token_pool_1_0_0.NewLockReleaseTokenPool(destWrappedPoolAddress, destChain)
	require.NoError(t, err)

	poolFloatValue := big.NewInt(1e18)

	destUser.Value = poolFloatValue
	_, err = destWrapped.Deposit(destUser)
	require.NoError(t, err)
	destChain.Commit()
	destUser.Value = nil

	_, err = destWrapped.Transfer(destUser, destWrappedPool.Address(), poolFloatValue)
	require.NoError(t, err)
	destChain.Commit()

	// Deploy and configure ge offramp.
	destPricesAddress, _, _, err := price_registry_1_2_0.DeployPriceRegistry(
		destUser,
		destChain,
		nil,
		[]common.Address{destLinkTokenAddress},
		60*60*24*14, // two weeks
	)
	require.NoError(t, err)
	destPriceRegistry, err := price_registry_1_2_0.NewPriceRegistry(destPricesAddress, destChain)
	require.NoError(t, err)

	// Deploy commit store.
	commitStoreAddress, _, _, err := commit_store_1_2_0.DeployCommitStore(
		destUser,  // user
		destChain, // client
		commit_store_1_2_0.CommitStoreStaticConfig{
			ChainSelector:       destChainSelector,
			SourceChainSelector: sourceChainSelector,
			OnRamp:              onRamp.Address(),
			ArmProxy:            destARMProxy.Address(),
		},
	)
	require.NoError(t, err)
	destChain.Commit()
	commitStore, err := commit_store_1_2_0.NewCommitStore(commitStoreAddress, destChain)
	require.NoError(t, err)

	offRampAddress, _, _, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
		destUser,
		destChain,
		evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore:         commitStore.Address(),
			ChainSelector:       destChainSelector,
			SourceChainSelector: sourceChainSelector,
			OnRamp:              onRampAddress,
			PrevOffRamp:         common.HexToAddress(""),
			ArmProxy:            armProxyDestAddress,
		},
		[]common.Address{sourceLinkTokenAddress, sourceWeth9addr},
		[]common.Address{destPoolAddress, destWrappedPool.Address()},
		evm_2_evm_offramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  LinkUSDValue(100),
			Rate:      LinkUSDValue(1),
		},
	)
	require.NoError(t, err)
	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(offRampAddress, destChain)
	require.NoError(t, err)
	_, err = destPool.ApplyChainUpdates(destUser,
		[]lock_release_token_pool.TokenPoolChainUpdate{{
			RemoteChainSelector: sourceChainSelector,
			Allowed:             true,
			OutboundRateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
			InboundRateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
		}},
	)
	require.NoError(t, err)

	_, err = destWrappedPool.ApplyRampUpdates(destUser,
		[]lock_release_token_pool_1_0_0.TokenPoolRampUpdate{},
		[]lock_release_token_pool_1_0_0.TokenPoolRampUpdate{{
			Ramp:    offRampAddress,
			Allowed: true,
			RateLimiterConfig: lock_release_token_pool_1_0_0.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  HundredLink,
				Rate:      big.NewInt(1e18),
			},
		}},
	)
	require.NoError(t, err)

	destChain.Commit()
	_, err = destPriceRegistry.ApplyPriceUpdatersUpdates(destUser, []common.Address{commitStoreAddress}, []common.Address{})
	require.NoError(t, err)
	_, err = destRouter.ApplyRampUpdates(destUser, nil,
		nil, []router.RouterOffRamp{{SourceChainSelector: sourceChainSelector, OffRamp: offRampAddress}})
	require.NoError(t, err)

	// Deploy 2 revertable (one SS one non-SS)
	revertingMessageReceiver1Address, _, _, err := maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(destUser, destChain, false)
	require.NoError(t, err)
	revertingMessageReceiver1, _ := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(revertingMessageReceiver1Address, destChain)
	revertingMessageReceiver2Address, _, _, err := maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(destUser, destChain, false)
	require.NoError(t, err)
	revertingMessageReceiver2, _ := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(revertingMessageReceiver2Address, destChain)
	// Need to commit here, or we will hit the block gas limit when deploying the executor
	sourceChain.Commit()
	destChain.Commit()

	// Ensure we have at least finality blocks.
	for i := 0; i < 50; i++ {
		sourceChain.Commit()
		destChain.Commit()
	}

	source := SourceChain{
		Common: Common{
			ChainID:           sourceChainID,
			ChainSelector:     sourceChainSelector,
			User:              sourceUser,
			Chain:             sourceChain,
			LinkToken:         sourceLinkToken,
			LinkTokenPool:     sourcePool,
			CustomToken:       sourceCustomToken,
			ARM:               sourceARM,
			ARMProxy:          sourceARMProxy,
			PriceRegistry:     srcPriceRegistry,
			WrappedNative:     sourceWrapped,
			WrappedNativePool: sourceWeth9Pool,
		},
		Router: sourceRouter,
		OnRamp: onRamp,
	}
	dest := DestinationChain{
		Common: Common{
			ChainID:           destChainID,
			ChainSelector:     destChainSelector,
			User:              destUser,
			Chain:             destChain,
			LinkToken:         destLinkToken,
			LinkTokenPool:     destPool,
			CustomToken:       destCustomToken,
			ARM:               destARM,
			ARMProxy:          destARMProxy,
			PriceRegistry:     destPriceRegistry,
			WrappedNative:     destWrapped,
			WrappedNativePool: destWrappedPool,
		},
		CommitStore: commitStore,
		Router:      destRouter,
		OffRamp:     offRamp,
		Receivers:   []MaybeRevertReceiver{{Receiver: revertingMessageReceiver1, Strict: false}, {Receiver: revertingMessageReceiver2, Strict: true}},
	}

	return CCIPContracts{
		Source: source,
		Dest:   dest,
	}
}

func (c *CCIPContracts) SendRequest(t *testing.T, msg router.ClientEVM2AnyMessage) *types.Transaction {
	tx, err := c.Source.Router.CcipSend(c.Source.User, c.Dest.ChainSelector, msg)
	require.NoError(t, err)
	testhelpers.ConfirmTxs(t, []*types.Transaction{tx}, c.Source.Chain)
	return tx
}

func (c *CCIPContracts) AssertExecState(t *testing.T, log logpoller.Log, state MessageExecutionState, offRampOpts ...common.Address) {
	var offRamp *evm_2_evm_offramp.EVM2EVMOffRamp
	var err error
	if len(offRampOpts) > 0 {
		offRamp, err = evm_2_evm_offramp.NewEVM2EVMOffRamp(offRampOpts[0], c.Dest.Chain)
		require.NoError(t, err)
	} else {
		require.NotNil(t, c.Dest.OffRamp, "no offRamp configured")
		offRamp = c.Dest.OffRamp
	}
	executionStateChanged, err := offRamp.ParseExecutionStateChanged(log.ToGethLog())
	require.NoError(t, err)
	if MessageExecutionState(executionStateChanged.State) != state {
		t.Log("Execution failed", hexutil.Encode(executionStateChanged.ReturnData))
		t.Fail()
	}
}

func GetEVMExtraArgsV1(gasLimit *big.Int, strict bool) ([]byte, error) {
	EVMV1Tag := []byte{0x97, 0xa6, 0x57, 0xc9}

	encodedArgs, err := utils.ABIEncode(`[{"type":"uint256"},{"type":"bool"}]`, gasLimit, strict)
	if err != nil {
		return nil, err
	}

	return append(EVMV1Tag, encodedArgs...), nil
}

type ManualExecArgs struct {
	SourceChainID, DestChainID uint64
	DestUser                   *bind.TransactOpts
	SourceChain, DestChain     bind.ContractBackend
	SourceStartBlock           *big.Int // the block in/after which failed ccip-send transaction was triggered
	DestStartBlock             uint64   // the start block for filtering ReportAccepted event (including the failed seq num)
	// in destination chain. if not provided to be derived by ApproxDestStartBlock method
	DestLatestBlockNum uint64 // current block number in destination
	DestDeployedAt     uint64 // destination block number for the initial destination contract deployment.
	// Can be any number before the tx was reverted in destination chain. Preferably this needs to be set up with
	// a value greater than zero to avoid performance issue in locating approximate destination block
	SendReqLogIndex uint   // log index of the CCIPSendRequested log in source chain
	SendReqTxHash   string // tx hash of the ccip-send transaction for which execution was reverted
	CommitStore     string
	OnRamp          string
	OffRamp         string
	SeqNr           uint64
	GasLimit        *big.Int
}

// ApproxDestStartBlock attempts to locate a block in destination chain with timestamp closest to the timestamp of the block
// in source chain in which ccip-send transaction was included
// it uses binary search to locate the block with the closest timestamp
// if the block located has a timestamp greater than the timestamp of mentioned source block
// it just returns the first block found with lesser timestamp of the source block
// providing a value of args.DestDeployedAt ensures better performance by reducing the range of block numbers to be traversed
func (args *ManualExecArgs) ApproxDestStartBlock() error {
	sourceBlockHdr, err := args.SourceChain.HeaderByNumber(context.Background(), args.SourceStartBlock)
	if err != nil {
		return err
	}
	sendTxTime := sourceBlockHdr.Time
	maxBlockNum := args.DestLatestBlockNum
	// setting this to an approx value of 1000 considering destination chain would have at least 1000 blocks before the transaction started
	minBlockNum := args.DestDeployedAt
	closestBlockNum := uint64(math.Floor((float64(maxBlockNum) + float64(minBlockNum)) / 2))
	var closestBlockHdr *types.Header
	closestBlockHdr, err = args.DestChain.HeaderByNumber(context.Background(), big.NewInt(int64(closestBlockNum)))
	if err != nil {
		return err
	}
	// to reduce the number of RPC calls increase the value of blockOffset
	blockOffset := uint64(10)
	for {
		blockNum := closestBlockHdr.Number.Uint64()
		if minBlockNum > maxBlockNum {
			break
		}
		timeDiff := math.Abs(float64(closestBlockHdr.Time - sendTxTime))
		// break if the difference in timestamp is lesser than 1 minute
		if timeDiff < 60 {
			break
		} else if closestBlockHdr.Time > sendTxTime {
			maxBlockNum = blockNum - 1
		} else {
			minBlockNum = blockNum + 1
		}
		closestBlockNum = uint64(math.Floor((float64(maxBlockNum) + float64(minBlockNum)) / 2))
		closestBlockHdr, err = args.DestChain.HeaderByNumber(context.Background(), big.NewInt(int64(closestBlockNum)))
		if err != nil {
			return err
		}
	}

	for closestBlockHdr.Time > sendTxTime {
		closestBlockNum = closestBlockNum - blockOffset
		if closestBlockNum <= 0 {
			return fmt.Errorf("approx destination blocknumber not found")
		}
		closestBlockHdr, err = args.DestChain.HeaderByNumber(context.Background(), big.NewInt(int64(closestBlockNum)))
		if err != nil {
			return err
		}
	}
	args.DestStartBlock = closestBlockHdr.Number.Uint64()
	fmt.Println("using approx destination start block number", args.DestStartBlock)
	return nil
}

func (args *ManualExecArgs) FindSeqNrFromCCIPSendRequested() (uint64, error) {
	var seqNr uint64
	onRampContract, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress(args.OnRamp), args.SourceChain)
	if err != nil {
		return seqNr, err
	}
	iterator, err := onRampContract.FilterCCIPSendRequested(&bind.FilterOpts{
		Start: args.SourceStartBlock.Uint64(),
	})
	if err != nil {
		return seqNr, err
	}
	for iterator.Next() {
		if iterator.Event.Raw.Index == args.SendReqLogIndex &&
			iterator.Event.Raw.TxHash.Hex() == args.SendReqTxHash {
			seqNr = iterator.Event.Message.SequenceNumber
			break
		}
	}
	if seqNr == 0 {
		return seqNr,
			fmt.Errorf("no CCIPSendRequested logs found for logIndex %d starting from block number %d", args.SendReqLogIndex, args.SourceStartBlock)
	}
	return seqNr, nil
}

func (args *ManualExecArgs) ExecuteManually() (*types.Transaction, error) {
	if args.SourceChainID == 0 ||
		args.DestChainID == 0 ||
		args.DestUser == nil {
		return nil, fmt.Errorf("chain ids and owners are mandatory for source and dest chain")
	}
	if !common.IsHexAddress(args.CommitStore) ||
		!common.IsHexAddress(args.OffRamp) ||
		!common.IsHexAddress(args.OnRamp) {
		return nil, fmt.Errorf("contract addresses must be valid hex address")
	}
	if args.SendReqTxHash == "" {
		return nil, fmt.Errorf("tx hash of ccip-send request are required")
	}
	if args.SourceStartBlock == nil {
		return nil, fmt.Errorf("must provide the value of source block in/after which ccip-send tx was included")
	}
	if args.SeqNr == 0 {
		if args.SendReqLogIndex == 0 {
			return nil, fmt.Errorf("must provide the value of log index of ccip-send request")
		}
		// locate seq nr from CCIPSendRequested log
		seqNr, err := args.FindSeqNrFromCCIPSendRequested()
		if err != nil {
			return nil, err
		}
		args.SeqNr = seqNr
	}
	commitStore, err := commit_store_1_2_0.NewCommitStore(common.HexToAddress(args.CommitStore), args.DestChain)
	if err != nil {
		return nil, err
	}
	if args.DestStartBlock < 1 {
		err = args.ApproxDestStartBlock()
		if err != nil {
			return nil, err
		}
	}
	iterator, err := commitStore.FilterReportAccepted(&bind.FilterOpts{Start: args.DestStartBlock})
	if err != nil {
		return nil, err
	}

	var commitReport *commit_store_1_2_0.CommitStoreCommitReport
	for iterator.Next() {
		if iterator.Event.Report.Interval.Min <= args.SeqNr && iterator.Event.Report.Interval.Max >= args.SeqNr {
			commitReport = &iterator.Event.Report
			fmt.Println("Found root")
			break
		}
	}
	if commitReport == nil {
		return nil, fmt.Errorf("unable to find seq num %d in commit report", args.SeqNr)
	}

	return args.execute(commitReport)
}

func (args *ManualExecArgs) execute(report *commit_store_1_2_0.CommitStoreCommitReport) (*types.Transaction, error) {
	log.Info().Msg("Executing request manually")
	seqNr := args.SeqNr
	// Build a merkle tree for the report
	mctx := hashutil.NewKeccak()
	onRampContract, err := evm_2_evm_onramp_1_2_0.NewEVM2EVMOnRamp(common.HexToAddress(args.OnRamp), args.SourceChain)
	if err != nil {
		return nil, err
	}
	leafHasher := v1_2_0.NewLeafHasher(args.SourceChainID, args.DestChainID, common.HexToAddress(args.OnRamp), mctx, onRampContract)
	if leafHasher == nil {
		return nil, fmt.Errorf("unable to create leaf hasher")
	}

	var leaves [][32]byte
	var curr, prove int
	var msgs []evm_2_evm_offramp.InternalEVM2EVMMessage
	var manualExecGasLimits []*big.Int
	var tokenData [][][]byte
	sendRequestedIterator, err := onRampContract.FilterCCIPSendRequested(&bind.FilterOpts{
		Start: args.SourceStartBlock.Uint64(),
	})
	if err != nil {
		return nil, err
	}
	for sendRequestedIterator.Next() {
		if sendRequestedIterator.Event.Message.SequenceNumber <= report.Interval.Max &&
			sendRequestedIterator.Event.Message.SequenceNumber >= report.Interval.Min {
			fmt.Println("Found seq num", sendRequestedIterator.Event.Message.SequenceNumber, report.Interval)
			hash, err2 := leafHasher.HashLeaf(sendRequestedIterator.Event.Raw)
			if err2 != nil {
				return nil, err2
			}
			leaves = append(leaves, hash)
			if sendRequestedIterator.Event.Message.SequenceNumber == seqNr {
				fmt.Printf("Found proving %d %+v\n", curr, sendRequestedIterator.Event.Message)
				var tokensAndAmounts []evm_2_evm_offramp.ClientEVMTokenAmount
				for _, tokenAndAmount := range sendRequestedIterator.Event.Message.TokenAmounts {
					tokensAndAmounts = append(tokensAndAmounts, evm_2_evm_offramp.ClientEVMTokenAmount{
						Token:  tokenAndAmount.Token,
						Amount: tokenAndAmount.Amount,
					})
				}
				msg := evm_2_evm_offramp.InternalEVM2EVMMessage{
					SourceChainSelector: sendRequestedIterator.Event.Message.SourceChainSelector,
					Sender:              sendRequestedIterator.Event.Message.Sender,
					Receiver:            sendRequestedIterator.Event.Message.Receiver,
					SequenceNumber:      sendRequestedIterator.Event.Message.SequenceNumber,
					GasLimit:            sendRequestedIterator.Event.Message.GasLimit,
					Strict:              sendRequestedIterator.Event.Message.Strict,
					Nonce:               sendRequestedIterator.Event.Message.Nonce,
					FeeToken:            sendRequestedIterator.Event.Message.FeeToken,
					FeeTokenAmount:      sendRequestedIterator.Event.Message.FeeTokenAmount,
					Data:                sendRequestedIterator.Event.Message.Data,
					TokenAmounts:        tokensAndAmounts,
					SourceTokenData:     sendRequestedIterator.Event.Message.SourceTokenData,
					MessageId:           sendRequestedIterator.Event.Message.MessageId,
				}
				msgs = append(msgs, msg)
				if args.GasLimit != nil {
					msg.GasLimit = args.GasLimit
				}
				manualExecGasLimits = append(manualExecGasLimits, msg.GasLimit)
				var msgTokenData [][]byte
				for range sendRequestedIterator.Event.Message.TokenAmounts {
					msgTokenData = append(msgTokenData, []byte{})
				}

				tokenData = append(tokenData, msgTokenData)
				prove = curr
			}
			curr++
		}
	}
	sendRequestedIterator.Close()
	if msgs == nil {
		return nil, fmt.Errorf("unable to find msg with seqNr %d", seqNr)
	}
	tree, err := merklemulti.NewTree(mctx, leaves)
	if err != nil {
		return nil, err
	}
	if tree.Root() != report.MerkleRoot {
		return nil, fmt.Errorf("root doesn't match")
	}

	proof, err := tree.Prove([]int{prove})
	if err != nil {
		return nil, err
	}

	offRampProof := evm_2_evm_offramp.InternalExecutionReport{
		Messages:          msgs,
		OffchainTokenData: tokenData,
		Proofs:            proof.Hashes,
		ProofFlagBits:     abihelpers.ProofFlagsToBits(proof.SourceFlags),
	}
	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(common.HexToAddress(args.OffRamp), args.DestChain)
	if err != nil {
		return nil, err
	}
	// Execute.
	return offRamp.ManuallyExecute(args.DestUser, offRampProof, manualExecGasLimits)
}

func (c *CCIPContracts) ExecuteMessage(
	t *testing.T,
	req logpoller.Log,
	txHash common.Hash,
	destStartBlock uint64,
) uint64 {
	t.Log("Executing request manually")
	sendReqReceipt, err := c.Source.Chain.TransactionReceipt(context.Background(), txHash)
	require.NoError(t, err)
	args := ManualExecArgs{
		SourceChainID:      c.Source.ChainID,
		DestChainID:        c.Dest.ChainID,
		DestUser:           c.Dest.User,
		SourceChain:        c.Source.Chain,
		DestChain:          c.Dest.Chain,
		SourceStartBlock:   sendReqReceipt.BlockNumber,
		DestStartBlock:     destStartBlock,
		DestLatestBlockNum: c.Dest.Chain.Blockchain().CurrentBlock().Number.Uint64(),
		SendReqLogIndex:    uint(req.LogIndex),
		SendReqTxHash:      txHash.String(),
		CommitStore:        c.Dest.CommitStore.Address().String(),
		OnRamp:             c.Source.OnRamp.Address().String(),
		OffRamp:            c.Dest.OffRamp.Address().String(),
	}
	tx, err := args.ExecuteManually()
	require.NoError(t, err)
	c.Dest.Chain.Commit()
	c.Source.Chain.Commit()
	rec, err := c.Dest.Chain.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), rec.Status, "manual execution failed")
	t.Logf("Manual Execution completed for seqNum %d", args.SeqNr)
	return args.SeqNr
}

func GetBalance(t *testing.T, chain bind.ContractBackend, tokenAddr common.Address, addr common.Address) *big.Int {
	token, err := link_token_interface.NewLinkToken(tokenAddr, chain)
	require.NoError(t, err)
	bal, err := token.BalanceOf(nil, addr)
	require.NoError(t, err)
	return bal
}
