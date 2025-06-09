package v1_5

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	config2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"

	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	v1_5changeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	plugintesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

func AddLanes(t *testing.T, e cldf.Environment, state stateview.CCIPOnChainState, pairs []testhelpers.SourceDestPair) cldf.Environment {
	addLanesCfg, commitOCR2Configs, execOCR2Configs, jobspecs := LaneConfigsForChains(t, e, state, pairs)
	var err error
	e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5changeset.DeployLanesChangeset),
		v1_5changeset.DeployLanesConfig{Configs: addLanesCfg},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5changeset.SetOCR2ConfigForTestChangeset),
		v1_5changeset.OCR2Config{CommitConfigs: commitOCR2Configs, ExecConfigs: execOCR2Configs},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5changeset.JobSpecsForLanesChangeset),
		v1_5changeset.JobSpecsForLanesConfig{Configs: jobspecs},
	))
	require.NoError(t, err)
	return e
}

func LaneConfigsForChains(t *testing.T, env cldf.Environment, state stateview.CCIPOnChainState, pairs []testhelpers.SourceDestPair) (
	[]v1_5changeset.DeployLaneConfig,
	[]v1_5changeset.CommitOCR2ConfigParams,
	[]v1_5changeset.ExecuteOCR2ConfigParams,
	[]v1_5changeset.JobSpecInput,
) {
	addLanesCfg := make([]v1_5changeset.DeployLaneConfig, 0)
	commitOCR2Configs := make([]v1_5changeset.CommitOCR2ConfigParams, 0)
	execOCR2Configs := make([]v1_5changeset.ExecuteOCR2ConfigParams, 0)
	jobSpecs := make([]v1_5changeset.JobSpecInput, 0)
	for _, pair := range pairs {
		dest := pair.DestChainSelector
		src := pair.SourceChainSelector
		sourceChainState := state.MustGetEVMChainState(src)
		destChainState := state.MustGetEVMChainState(dest)
		_, err := sourceChainState.LinkTokenAddress()
		require.NoError(t, err)
		require.NotNil(t, sourceChainState.RMNProxy)
		require.NotNil(t, sourceChainState.TokenAdminRegistry)
		require.NotNil(t, sourceChainState.Router)
		require.NotNil(t, sourceChainState.PriceRegistry)
		require.NotNil(t, sourceChainState.Weth9)
		_, err = destChainState.LinkTokenAddress()
		require.NoError(t, err)
		require.NotNil(t, destChainState.RMNProxy)
		require.NotNil(t, destChainState.TokenAdminRegistry)
		priceGetterConfig := CreatePriceGetterConfig(t, state, src, dest)
		block, err := env.BlockChains.EVMChains()[dest].Client.HeaderByNumber(context.Background(), nil)
		require.NoError(t, err)
		destEVMChainIdStr, err := chain_selectors.GetChainIDFromSelector(dest)
		require.NoError(t, err)
		destEVMChainId, err := strconv.ParseUint(destEVMChainIdStr, 10, 64)
		require.NoError(t, err)
		jobSpecs = append(jobSpecs, v1_5changeset.JobSpecInput{
			SourceChainSelector:      src,
			DestinationChainSelector: dest,
			DestEVMChainID:           destEVMChainId,
			PriceGetterConfigJson:    priceGetterConfig,
			DestinationStartBlock:    block.Number.Uint64(),
		})
		srcLinkTokenAddr, err := sourceChainState.LinkTokenAddress()
		require.NoError(t, err)
		addLanesCfg = append(addLanesCfg, v1_5changeset.DeployLaneConfig{
			SourceChainSelector:      src,
			DestinationChainSelector: dest,
			OnRampStaticCfg: evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
				LinkToken:          srcLinkTokenAddr,
				ChainSelector:      src,
				DestChainSelector:  dest,
				DefaultTxGasLimit:  200_000,
				MaxNopFeesJuels:    big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
				PrevOnRamp:         common.Address{},
				RmnProxy:           sourceChainState.RMNProxy.Address(),
				TokenAdminRegistry: sourceChainState.TokenAdminRegistry.Address(),
			},
			OnRampDynamicCfg: evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
				Router:                            sourceChainState.Router.Address(),
				MaxNumberOfTokensPerMsg:           5,
				DestGasOverhead:                   350_000,
				DestGasPerPayloadByte:             16,
				DestDataAvailabilityOverheadGas:   33_596,
				DestGasPerDataAvailabilityByte:    16,
				DestDataAvailabilityMultiplierBps: 6840,
				PriceRegistry:                     sourceChainState.PriceRegistry.Address(),
				MaxDataBytes:                      1e5,
				MaxPerMsgGasLimit:                 4_000_000,
				DefaultTokenFeeUSDCents:           50,
				DefaultTokenDestGasOverhead:       125_000,
			},
			OnRampFeeTokenArgs: []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
				{
					Token:                      srcLinkTokenAddr,
					NetworkFeeUSDCents:         1_00,
					GasMultiplierWeiPerEth:     1e18,
					PremiumMultiplierWeiPerEth: 9e17,
					Enabled:                    true,
				},
				{
					Token:                      sourceChainState.Weth9.Address(),
					NetworkFeeUSDCents:         1_00,
					GasMultiplierWeiPerEth:     1e18,
					PremiumMultiplierWeiPerEth: 1e18,
					Enabled:                    true,
				},
			},
			OnRampTransferTokenCfgs: []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				{
					Token:                     srcLinkTokenAddr,
					MinFeeUSDCents:            50,           // $0.5
					MaxFeeUSDCents:            1_000_000_00, // $ 1 million
					DeciBps:                   5_0,          // 5 bps
					DestGasOverhead:           350_000,
					DestBytesOverhead:         32,
					AggregateRateLimitEnabled: true,
				},
			},
			OnRampNopsAndWeight: []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
			OnRampRateLimiterCfg: evm_2_evm_onramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  plugintesthelpers.LinkUSDValue(100),
				Rate:      plugintesthelpers.LinkUSDValue(1),
			},
			OffRampRateLimiterCfg: evm_2_evm_offramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  plugintesthelpers.LinkUSDValue(100),
				Rate:      plugintesthelpers.LinkUSDValue(1),
			},
			InitialTokenPrices: []price_registry_1_2_0.InternalTokenPriceUpdate{
				{
					SourceToken: srcLinkTokenAddr,
					UsdPerToken: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(20)),
				},
				{
					SourceToken: sourceChainState.Weth9.Address(),
					UsdPerToken: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2000)),
				},
			},
			GasPriceUpdates: []price_registry_1_2_0.InternalGasPriceUpdate{
				{
					DestChainSelector: dest,
					UsdPerUnitGas:     big.NewInt(20000e9),
				},
			},
		})
		commitOCR2Configs = append(commitOCR2Configs, v1_5changeset.CommitOCR2ConfigParams{
			SourceChainSelector:      src,
			DestinationChainSelector: dest,
			OCR2ConfigParams:         DefaultOCRParams(),
			GasPriceHeartBeat:        *config.MustNewDuration(10 * time.Second),
			DAGasPriceDeviationPPB:   1,
			ExecGasPriceDeviationPPB: 1,
			TokenPriceHeartBeat:      *config.MustNewDuration(10 * time.Second),
			TokenPriceDeviationPPB:   1,
			InflightCacheExpiry:      *config.MustNewDuration(5 * time.Second),
			PriceReportingDisabled:   false,
		})
		execOCR2Configs = append(execOCR2Configs, v1_5changeset.ExecuteOCR2ConfigParams{
			DestinationChainSelector:    dest,
			SourceChainSelector:         src,
			DestOptimisticConfirmations: 1,
			BatchGasLimit:               5_000_000,
			RelativeBoostPerWaitHour:    0.07,
			InflightCacheExpiry:         *config.MustNewDuration(1 * time.Minute),
			RootSnoozeTime:              *config.MustNewDuration(1 * time.Minute),
			BatchingStrategyID:          0,
			MessageVisibilityInterval:   config.Duration{},
			ExecOnchainConfig: evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig{
				PermissionLessExecutionThresholdSeconds: uint32(24 * time.Hour.Seconds()),
				MaxDataBytes:                            1e5,
				MaxNumberOfTokensPerMsg:                 5,
			},
			OCR2ConfigParams: DefaultOCRParams(),
		})
	}
	return addLanesCfg, commitOCR2Configs, execOCR2Configs, jobSpecs
}

// CreatePriceGetterConfig returns price getter config as json string.
func CreatePriceGetterConfig(t *testing.T, state stateview.CCIPOnChainState, source, dest uint64) string {
	sourceRouter := state.MustGetEVMChainState(source).Router
	destRouter := state.MustGetEVMChainState(dest).Router
	destLinkAddr, err := state.MustGetEVMChainState(dest).LinkTokenAddress()
	require.NoError(t, err)

	linkPriceDest, ok := big.NewInt(0).SetString("8000000000000000000", 10)
	require.True(t, ok)

	ethPrice, ok := big.NewInt(0).SetString("1700000000000000000000", 10)
	require.True(t, ok)

	sourceWrappedNative, err := sourceRouter.GetWrappedNative(nil)
	require.NoError(t, err)

	destWrappedNative, err := destRouter.GetWrappedNative(nil)
	require.NoError(t, err)

	cfg := config2.DynamicPriceGetterConfig{
		TokenPrices: []config2.TokenPriceConfig{
			{
				TokenAddress:  destLinkAddr,
				ChainSelector: dest,
				StaticConfig:  &config2.StaticPriceConfig{ChainID: dest, Price: linkPriceDest},
			},
			{
				TokenAddress:  sourceWrappedNative, // eth
				ChainSelector: source,
				StaticConfig:  &config2.StaticPriceConfig{ChainID: source, Price: ethPrice},
			},
			{
				TokenAddress:  destWrappedNative, // eth
				ChainSelector: dest,
				StaticConfig:  &config2.StaticPriceConfig{ChainID: dest, Price: ethPrice},
			},
		},
	}

	cfgJSON, err := json.Marshal(cfg)
	require.NoError(t, err)

	return string(cfgJSON)
}

func DefaultOCRParams() confighelper.PublicConfig {
	return confighelper.PublicConfig{
		DeltaProgress:                           2 * time.Second,
		DeltaResend:                             1 * time.Second,
		DeltaRound:                              1 * time.Second,
		DeltaGrace:                              500 * time.Millisecond,
		DeltaStage:                              2 * time.Second,
		RMax:                                    3,
		MaxDurationInitialization:               nil,
		MaxDurationQuery:                        50 * time.Millisecond,
		MaxDurationObservation:                  1 * time.Second,
		MaxDurationReport:                       100 * time.Millisecond,
		MaxDurationShouldAcceptFinalizedReport:  100 * time.Millisecond,
		MaxDurationShouldTransmitAcceptedReport: 100 * time.Millisecond,
	}
}

func SendRequest(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	opts ...testhelpers.SendReqOpts,
) (*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested, error) {
	cfg := &testhelpers.CCIPSendReqConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// Set default sender if not provided
	if cfg.Sender == nil {
		cfg.Sender = e.BlockChains.EVMChains()[cfg.SourceChain].DeployerKey
	}
	t.Logf("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, cfg.Sender.From.String())
	tx, blockNum, err := testhelpers.CCIPSendRequest(e, state, cfg)
	if err != nil {
		return nil, err
	}

	onRamp := state.MustGetEVMChainState(cfg.SourceChain).EVM2EVMOnRamp[cfg.DestChain]

	it, err := onRamp.FilterCCIPSendRequested(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}

	require.True(t, it.Next())
	t.Logf("CCIP message (id %x) sent from chain selector %d to chain selector %d tx %s seqNum %d sender %s",
		it.Event.Message.MessageId[:],
		cfg.SourceChain,
		cfg.DestChain,
		tx.Hash().String(),
		it.Event.Message.SequenceNumber,
		it.Event.Message.Sender.String(),
	)
	return it.Event, nil
}

func WaitForCommit(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	commitStore *commit_store.CommitStore,
	seqNr uint64,
) {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			minSeqNr, err := commitStore.GetExpectedNextSequenceNumber(nil)
			require.NoError(t, err)
			t.Logf("Waiting for commit for sequence number %d, current min sequence number %d", seqNr, minSeqNr)
			if minSeqNr > seqNr {
				t.Logf("Commit for sequence number %d found", seqNr)
				return
			}
		case <-timer.C:
			t.Fatalf("timed out waiting for commit for sequence number %d for commit store %s ", seqNr, commitStore.Address().String())
			return
		}
	}
}

func WaitForNoCommit(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	commitStore *commit_store.CommitStore,
	seqNr uint64,
) {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			minSeqNr, err := commitStore.GetExpectedNextSequenceNumber(nil)
			require.NoError(t, err)
			t.Logf("Waiting for commit for sequence number %d, current min sequence number %d", seqNr, minSeqNr)
			if minSeqNr > seqNr {
				t.Fatalf("Commit for sequence number %d found while it was not expected", seqNr)
				return
			}
		case <-timer.C:
			t.Logf("Successfully observed no commit for sequence number %d for commit store %s during 30s period", seqNr, commitStore.Address().String())
			return
		}
	}
}

func WaitForExecute(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	offRamp *evm_2_evm_offramp.EVM2EVMOffRamp,
	seqNrs []uint64,
	blockNum uint64,
) {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t.Logf("Waiting for execute for sequence numbers %v", seqNrs)
			it, err := offRamp.FilterExecutionStateChanged(
				&bind.FilterOpts{
					Start: blockNum,
				}, seqNrs, [][32]byte{})
			require.NoError(t, err)
			for it.Next() {
				t.Logf("Execution state changed for sequence number=%d current state=%d", it.Event.SequenceNumber, it.Event.State)
				if cciptypes.MessageExecutionState(it.Event.State) == cciptypes.ExecutionStateSuccess {
					t.Logf("Execution for sequence number %d found", it.Event.SequenceNumber)
					return
				}
				t.Logf("Execution for sequence number %d resulted in status %d", it.Event.SequenceNumber, it.Event.State)
				t.Fail()
			}
		case <-timer.C:
			t.Fatalf("timed out waiting for execute for sequence numbers %v for offramp %s ", seqNrs, offRamp.Address().String())
			return
		}
	}
}
