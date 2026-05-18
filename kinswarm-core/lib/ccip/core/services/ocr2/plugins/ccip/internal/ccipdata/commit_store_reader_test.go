package ccipdata_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_helper_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_helper_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func TestCommitOffchainConfig_Encoding(t *testing.T) {
	tests := map[string]struct {
		want      v1_2_0.JSONCommitOffchainConfig
		expectErr bool
	}{
		"encodes and decodes config with all fields set": {
			want: v1_2_0.JSONCommitOffchainConfig{
				SourceFinalityDepth:      3,
				DestFinalityDepth:        3,
				GasPriceHeartBeat:        *config.MustNewDuration(1 * time.Hour),
				DAGasPriceDeviationPPB:   5e7,
				ExecGasPriceDeviationPPB: 5e7,
				TokenPriceHeartBeat:      *config.MustNewDuration(1 * time.Hour),
				TokenPriceDeviationPPB:   5e7,
				InflightCacheExpiry:      *config.MustNewDuration(23456 * time.Second),
			},
		},
		"fails decoding when all fields present but with 0 values": {
			want: v1_2_0.JSONCommitOffchainConfig{
				SourceFinalityDepth:      0,
				DestFinalityDepth:        0,
				GasPriceHeartBeat:        *config.MustNewDuration(0),
				DAGasPriceDeviationPPB:   0,
				ExecGasPriceDeviationPPB: 0,
				TokenPriceHeartBeat:      *config.MustNewDuration(0),
				TokenPriceDeviationPPB:   0,
				InflightCacheExpiry:      *config.MustNewDuration(0),
			},
			expectErr: true,
		},
		"fails decoding when all fields are missing": {
			want:      v1_2_0.JSONCommitOffchainConfig{},
			expectErr: true,
		},
		"fails decoding when some fields are missing": {
			want: v1_2_0.JSONCommitOffchainConfig{
				SourceFinalityDepth:      3,
				GasPriceHeartBeat:        *config.MustNewDuration(1 * time.Hour),
				DAGasPriceDeviationPPB:   5e7,
				ExecGasPriceDeviationPPB: 5e7,
				TokenPriceHeartBeat:      *config.MustNewDuration(1 * time.Hour),
				TokenPriceDeviationPPB:   5e7,
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encode, err := ccipconfig.EncodeOffchainConfig(tc.want)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[v1_2_0.JSONCommitOffchainConfig](encode)

			if tc.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestCommitOnchainConfig(t *testing.T) {
	tests := []struct {
		name      string
		want      ccipdata.CommitOnchainConfig
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ccipdata.CommitOnchainConfig{
				PriceRegistry: utils.RandomAddress(),
			},
			expectErr: false,
		},
		{
			name:      "encodes and fails decoding config with missing fields",
			want:      ccipdata.CommitOnchainConfig{},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := abihelpers.EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := abihelpers.DecodeAbiStruct[ccipdata.CommitOnchainConfig](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}

func TestCommitStoreReaders(t *testing.T) {
	user, ec := newSim(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	headTracker := headtracker.NewSimulatedHeadTracker(ec, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	if lpOpts.PollPeriod == 0 {
		lpOpts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, pgtest.NewSqlxDB(t), lggr), ec, lggr, headTracker, lpOpts)

	// Deploy 2 commit store versions
	onramp1 := utils.RandomAddress()
	onramp2 := utils.RandomAddress()
	// Report
	rep := cciptypes.CommitStoreReport{
		TokenPrices: []cciptypes.TokenPrice{{Token: ccipcalc.EvmAddrToGeneric(utils.RandomAddress()), Value: big.NewInt(1)}},
		GasPrices:   []cciptypes.GasPrice{{DestChainSelector: 1, Value: big.NewInt(1)}},
		Interval:    cciptypes.CommitStoreInterval{Min: 1, Max: 10},
		MerkleRoot:  common.HexToHash("0x1"),
	}
	er := big.NewInt(1)
	armAddr, _, arm, err := mock_rmn_contract.DeployMockRMNContract(user, ec)
	require.NoError(t, err)
	addr, _, ch, err := commit_store_helper_1_0_0.DeployCommitStoreHelper(user, ec, commit_store_helper_1_0_0.CommitStoreStaticConfig{
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onramp1,
		ArmProxy:            armAddr,
	})
	require.NoError(t, err)
	addr2, _, ch2, err := commit_store_helper_1_2_0.DeployCommitStoreHelper(user, ec, commit_store_helper_1_2_0.CommitStoreStaticConfig{
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onramp2,
		ArmProxy:            armAddr,
	})
	require.NoError(t, err)
	commitAndGetBlockTs(ec) // Deploy these
	pr, _, _, err := price_registry_1_0_0.DeployPriceRegistry(user, ec, []common.Address{addr}, nil, 1e6)
	require.NoError(t, err)
	pr2, _, _, err := price_registry_1_2_0.DeployPriceRegistry(user, ec, []common.Address{addr2}, nil, 1e6)
	require.NoError(t, err)
	commitAndGetBlockTs(ec) // Deploy these
	ge := new(gasmocks.EvmFeeEstimator)
	lm := new(rollupMocks.L1Oracle)
	ge.On("L1Oracle").Return(lm)

	feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)
	feeEstimatorConfig.On(
		"ModifyGasPriceComponents",
		mock.AnythingOfType("context.backgroundCtx"),
		mock.AnythingOfType("*big.Int"),
		mock.AnythingOfType("*big.Int"),
	).Return(func(ctx context.Context, x, y *big.Int) (*big.Int, *big.Int, error) {
		return x, y, nil
	})
	maxGasPrice := big.NewInt(1e8)
	c10r, err := factory.NewCommitStoreReader(lggr, factory.NewEvmVersionFinder(), ccipcalc.EvmAddrToGeneric(addr), ec, lp, feeEstimatorConfig) // ge, maxGasPrice
	require.NoError(t, err)
	err = c10r.SetGasEstimator(ctx, ge)
	require.NoError(t, err)
	err = c10r.SetSourceMaxGasPrice(ctx, maxGasPrice)
	require.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(c10r).String(), reflect.TypeOf(&v1_0_0.CommitStore{}).String())
	c12r, err := factory.NewCommitStoreReader(lggr, factory.NewEvmVersionFinder(), ccipcalc.EvmAddrToGeneric(addr2), ec, lp, feeEstimatorConfig)
	require.NoError(t, err)
	err = c12r.SetGasEstimator(ctx, ge)
	require.NoError(t, err)
	err = c12r.SetSourceMaxGasPrice(ctx, maxGasPrice)
	require.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(c12r).String(), reflect.TypeOf(&v1_2_0.CommitStore{}).String())

	// Apply config
	signers := []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()}
	transmitters := []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()}
	onchainConfig, err := abihelpers.EncodeAbiStruct[ccipdata.CommitOnchainConfig](ccipdata.CommitOnchainConfig{
		PriceRegistry: pr,
	})
	require.NoError(t, err)

	sourceFinalityDepth := uint32(1)
	destFinalityDepth := uint32(2)
	commonOffchain := cciptypes.CommitOffchainConfig{
		GasPriceDeviationPPB:   1e6,
		GasPriceHeartBeat:      1 * time.Hour,
		TokenPriceDeviationPPB: 1e6,
		TokenPriceHeartBeat:    1 * time.Hour,
		InflightCacheExpiry:    3 * time.Hour,
		PriceReportingDisabled: false,
	}
	offchainConfig, err := ccipconfig.EncodeOffchainConfig[v1_0_0.CommitOffchainConfig](v1_0_0.CommitOffchainConfig{
		SourceFinalityDepth:   sourceFinalityDepth,
		DestFinalityDepth:     destFinalityDepth,
		FeeUpdateHeartBeat:    *config.MustNewDuration(commonOffchain.GasPriceHeartBeat),
		FeeUpdateDeviationPPB: commonOffchain.GasPriceDeviationPPB,
		InflightCacheExpiry:   *config.MustNewDuration(commonOffchain.InflightCacheExpiry),
	})
	require.NoError(t, err)
	_, err = ch.SetOCR2Config(user, signers, transmitters, 1, onchainConfig, 1, []byte{})
	require.NoError(t, err)
	onchainConfig2, err := abihelpers.EncodeAbiStruct[ccipdata.CommitOnchainConfig](ccipdata.CommitOnchainConfig{
		PriceRegistry: pr2,
	})
	require.NoError(t, err)
	offchainConfig2, err := ccipconfig.EncodeOffchainConfig[v1_2_0.JSONCommitOffchainConfig](v1_2_0.JSONCommitOffchainConfig{
		SourceFinalityDepth:      sourceFinalityDepth,
		DestFinalityDepth:        destFinalityDepth,
		GasPriceHeartBeat:        *config.MustNewDuration(commonOffchain.GasPriceHeartBeat),
		DAGasPriceDeviationPPB:   1e7,
		ExecGasPriceDeviationPPB: commonOffchain.GasPriceDeviationPPB,
		TokenPriceDeviationPPB:   commonOffchain.TokenPriceDeviationPPB,
		TokenPriceHeartBeat:      *config.MustNewDuration(commonOffchain.TokenPriceHeartBeat),
		InflightCacheExpiry:      *config.MustNewDuration(commonOffchain.InflightCacheExpiry),
	})
	require.NoError(t, err)
	_, err = ch2.SetOCR2Config(user, signers, transmitters, 1, onchainConfig2, 1, []byte{})
	require.NoError(t, err)
	commitAndGetBlockTs(ec)

	// Apply report
	b, err := c10r.EncodeCommitReport(ctx, rep)
	require.NoError(t, err)
	_, err = ch.Report(user, b, er)
	require.NoError(t, err)
	b, err = c12r.EncodeCommitReport(ctx, rep)
	require.NoError(t, err)
	_, err = ch2.Report(user, b, er)
	require.NoError(t, err)
	commitAndGetBlockTs(ec)

	// Capture all logs.
	lp.PollAndSaveLogs(context.Background(), 1)

	configs := map[string][][]byte{
		ccipdata.V1_0_0: {onchainConfig, offchainConfig},
		ccipdata.V1_2_0: {onchainConfig2, offchainConfig2},
	}
	crs := map[string]ccipdata.CommitStoreReader{
		ccipdata.V1_0_0: c10r,
		ccipdata.V1_2_0: c12r,
	}
	prs := map[string]common.Address{
		ccipdata.V1_0_0: pr,
		ccipdata.V1_2_0: pr2,
	}
	gasPrice := big.NewInt(10)
	daPrice := big.NewInt(20)
	ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, assets.NewWei(maxGasPrice), (*common.Address)(nil), (*common.Address)(nil)).Return(gas.EvmFee{Legacy: assets.NewWei(gasPrice)}, uint64(0), nil)
	lm.On("GasPrice", mock.Anything).Return(assets.NewWei(daPrice), nil)

	for v, cr := range crs {
		cr := cr
		t.Run("CommitStoreReader "+v, func(t *testing.T) {
			// Static config.
			cfg, err := cr.GetCommitStoreStaticConfig(context.Background())
			require.NoError(t, err)
			require.NotNil(t, cfg)

			// Assert encoding
			b, err := cr.EncodeCommitReport(ctx, rep)
			require.NoError(t, err)
			d, err := cr.DecodeCommitReport(ctx, b)
			require.NoError(t, err)
			assert.Equal(t, d, rep)

			// Assert reading
			latest, err := cr.GetLatestPriceEpochAndRound(context.Background())
			require.NoError(t, err)
			assert.Equal(t, er.Uint64(), latest)

			// Assert cursing
			down, err := cr.IsDown(context.Background())
			require.NoError(t, err)
			assert.False(t, down)
			_, err = arm.VoteToCurse(user, [32]byte{})
			require.NoError(t, err)
			ec.Commit()
			down, err = cr.IsDown(context.Background())
			require.NoError(t, err)
			assert.True(t, down)
			_, err = arm.OwnerUnvoteToCurse0(user, nil)
			require.NoError(t, err)
			ec.Commit()

			seqNr, err := cr.GetExpectedNextSequenceNumber(context.Background())
			require.NoError(t, err)
			assert.Equal(t, rep.Interval.Max+1, seqNr)

			reps, err := cr.GetCommitReportMatchingSeqNum(context.Background(), rep.Interval.Max+1, 0)
			require.NoError(t, err)
			assert.Len(t, reps, 0)

			reps, err = cr.GetCommitReportMatchingSeqNum(context.Background(), rep.Interval.Max, 0)
			require.NoError(t, err)
			require.Len(t, reps, 1)
			assert.Equal(t, reps[0].Interval, rep.Interval)
			assert.Equal(t, reps[0].MerkleRoot, rep.MerkleRoot)
			assert.Equal(t, reps[0].GasPrices, rep.GasPrices)
			assert.Equal(t, reps[0].TokenPrices, rep.TokenPrices)

			reps, err = cr.GetCommitReportMatchingSeqNum(context.Background(), rep.Interval.Min, 0)
			require.NoError(t, err)
			require.Len(t, reps, 1)
			assert.Equal(t, reps[0].Interval, rep.Interval)
			assert.Equal(t, reps[0].MerkleRoot, rep.MerkleRoot)
			assert.Equal(t, reps[0].GasPrices, rep.GasPrices)
			assert.Equal(t, reps[0].TokenPrices, rep.TokenPrices)

			reps, err = cr.GetCommitReportMatchingSeqNum(context.Background(), rep.Interval.Min-1, 0)
			require.NoError(t, err)
			require.Len(t, reps, 0)

			// Sanity
			reps, err = cr.GetAcceptedCommitReportsGteTimestamp(context.Background(), time.Unix(0, 0), 0)
			require.NoError(t, err)
			require.Len(t, reps, 1)
			assert.Equal(t, reps[0].Interval, rep.Interval)
			assert.Equal(t, reps[0].MerkleRoot, rep.MerkleRoot)
			assert.Equal(t, reps[0].GasPrices, rep.GasPrices)
			assert.Equal(t, reps[0].TokenPrices, rep.TokenPrices)

			// Until we detect the config, we'll have empty offchain config
			c1, err := cr.OffchainConfig(ctx)
			require.NoError(t, err)
			assert.Equal(t, c1, cciptypes.CommitOffchainConfig{})
			newPr, err := cr.ChangeConfig(ctx, configs[v][0], configs[v][1])
			require.NoError(t, err)
			assert.Equal(t, ccipcalc.EvmAddrToGeneric(prs[v]), newPr)

			c2, err := cr.OffchainConfig(ctx)
			require.NoError(t, err)
			assert.Equal(t, commonOffchain, c2)
			// We should be able to query for gas prices now.

			gpe, err := cr.GasPriceEstimator(ctx)
			require.NoError(t, err)

			gp, err := gpe.GetGasPrice(context.Background())
			require.NoError(t, err)
			assert.True(t, gp.Cmp(big.NewInt(0)) > 0)
		})
	}
}

func TestNewCommitStoreReader(t *testing.T) {
	var tt = []struct {
		typeAndVersion string
		expectedErr    string
	}{
		{
			typeAndVersion: "blah",
			expectedErr:    "unable to read type and version: invalid type and version blah",
		},
		{
			typeAndVersion: "EVM2EVMOffRamp 1.0.0",
			expectedErr:    "expected CommitStore got EVM2EVMOffRamp",
		},
		{
			typeAndVersion: "CommitStore 1.2.0",
			expectedErr:    "",
		},
		{
			typeAndVersion: "CommitStore 2.0.0",
			expectedErr:    "unsupported commit store version 2.0.0",
		},
	}
	for _, tc := range tt {
		t.Run(tc.typeAndVersion, func(t *testing.T) {
			b, err := utils.ABIEncode(`[{"type":"string"}]`, tc.typeAndVersion)
			require.NoError(t, err)
			c := evmclientmocks.NewClient(t)
			c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(b, nil)
			addr := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
			lp := lpmocks.NewLogPoller(t)
			if tc.expectedErr == "" {
				lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
			}

			feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

			_, err = factory.NewCommitStoreReader(logger.TestLogger(t), factory.NewEvmVersionFinder(), addr, c, lp, feeEstimatorConfig)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
