package ccipdata_test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_helper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_helper_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func assertFilterRegistration(t *testing.T, lp *lpmocks.LogPoller, buildCloser func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer, numFilter int) {
	// Expected filter properties for a closer:
	// - Should be the same filter set registered that is unregistered
	// - Should be registered to the address specified
	// - Number of events specific to this component should be registered
	addr := common.HexToAddress("0x1234")
	var filters []logpoller.Filter

	lp.On("RegisterFilter", mock.Anything).Run(func(args mock.Arguments) {
		f := args.Get(0).(logpoller.Filter)
		require.Equal(t, len(f.Addresses), 1)
		require.Equal(t, f.Addresses[0], addr)
		filters = append(filters, f)
	}).Return(nil).Times(numFilter)

	c := buildCloser(lp, addr)
	for _, filter := range filters {
		lp.On("UnregisterFilter", filter.Name).Return(nil)
	}

	require.NoError(t, c.Close())
	lp.AssertExpectations(t)
}

func TestCommitFilters(t *testing.T) {
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewCommitStoreV1_0_0(logger.TestLogger(t), addr, new(mocks.Client), lp, nil)
		require.NoError(t, err)
		return c
	}, 1)
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewCommitStoreV1_2_0(logger.TestLogger(t), addr, new(mocks.Client), lp, nil)
		require.NoError(t, err)
		return c
	}, 1)
}

func TestCommitOffchainConfig_Encoding(t *testing.T) {
	tests := map[string]struct {
		want      ccipdata.CommitOffchainConfigV1_2_0
		expectErr bool
	}{
		"encodes and decodes config with all fields set": {
			want: ccipdata.CommitOffchainConfigV1_2_0{
				SourceFinalityDepth:      3,
				DestFinalityDepth:        3,
				GasPriceHeartBeat:        models.MustMakeDuration(1 * time.Hour),
				DAGasPriceDeviationPPB:   5e7,
				ExecGasPriceDeviationPPB: 5e7,
				TokenPriceHeartBeat:      models.MustMakeDuration(1 * time.Hour),
				TokenPriceDeviationPPB:   5e7,
				MaxGasPrice:              200e9,
				InflightCacheExpiry:      models.MustMakeDuration(23456 * time.Second),
			},
		},
		"fails decoding when all fields present but with 0 values": {
			want: ccipdata.CommitOffchainConfigV1_2_0{
				SourceFinalityDepth:      0,
				DestFinalityDepth:        0,
				GasPriceHeartBeat:        models.MustMakeDuration(0),
				DAGasPriceDeviationPPB:   0,
				ExecGasPriceDeviationPPB: 0,
				TokenPriceHeartBeat:      models.MustMakeDuration(0),
				TokenPriceDeviationPPB:   0,
				MaxGasPrice:              0,
				InflightCacheExpiry:      models.MustMakeDuration(0),
			},
			expectErr: true,
		},
		"fails decoding when all fields are missing": {
			want:      ccipdata.CommitOffchainConfigV1_2_0{},
			expectErr: true,
		},
		"fails decoding when some fields are missing": {
			want: ccipdata.CommitOffchainConfigV1_2_0{
				SourceFinalityDepth:      3,
				GasPriceHeartBeat:        models.MustMakeDuration(1 * time.Hour),
				DAGasPriceDeviationPPB:   5e7,
				ExecGasPriceDeviationPPB: 5e7,
				TokenPriceHeartBeat:      models.MustMakeDuration(1 * time.Hour),
				TokenPriceDeviationPPB:   5e7,
				MaxGasPrice:              200e9,
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encode, err := ccipconfig.EncodeOffchainConfig(tc.want)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[ccipdata.CommitOffchainConfigV1_2_0](encode)

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
	lggr := logger.TestLogger(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, pgtest.NewSqlxDB(t), lggr, pgtest.NewQConfig(true)), ec, lggr, 100*time.Millisecond, false, 2, 3, 2, 1000)

	// Deploy 2 commit store versions
	onramp1 := utils.RandomAddress()
	onramp2 := utils.RandomAddress()
	// Report
	rep := ccipdata.CommitStoreReport{
		TokenPrices: []ccipdata.TokenPrice{{Token: utils.RandomAddress(), Value: big.NewInt(1)}},
		GasPrices:   []ccipdata.GasPrice{{DestChainSelector: 1, Value: big.NewInt(1)}},
		Interval:    ccipdata.CommitStoreInterval{Min: 1, Max: 10},
		MerkleRoot:  common.HexToHash("0x1"),
	}
	er := big.NewInt(1)
	armAddr, _, arm, err := mock_arm_contract.DeployMockARMContract(user, ec)
	require.NoError(t, err)
	addr, _, ch, err := commit_store_helper_1_0_0.DeployCommitStoreHelper(user, ec, commit_store_helper_1_0_0.CommitStoreStaticConfig{
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onramp1,
		ArmProxy:            armAddr,
	})
	require.NoError(t, err)
	addr2, _, ch2, err := commit_store_helper.DeployCommitStoreHelper(user, ec, commit_store_helper.CommitStoreStaticConfig{
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onramp2,
		ArmProxy:            armAddr,
	})
	require.NoError(t, err)
	commitAndGetBlockTs(ec) // Deploy these
	pr, _, _, err := price_registry_1_0_0.DeployPriceRegistry(user, ec, []common.Address{addr}, nil, 1e6)
	require.NoError(t, err)
	pr2, _, _, err := price_registry.DeployPriceRegistry(user, ec, []common.Address{addr2}, nil, 1e6)
	require.NoError(t, err)
	commitAndGetBlockTs(ec) // Deploy these
	ge := new(gasmocks.EvmFeeEstimator)
	c10r, err := ccipdata.NewCommitStoreReader(lggr, addr, ec, lp, ge)
	require.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(c10r).String(), reflect.TypeOf(&ccipdata.CommitStoreV1_0_0{}).String())
	c12r, err := ccipdata.NewCommitStoreReader(lggr, addr2, ec, lp, ge)
	require.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(c12r).String(), reflect.TypeOf(&ccipdata.CommitStoreV1_2_0{}).String())

	// Apply config
	signers := []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()}
	transmitters := []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()}
	onchainConfig, err := abihelpers.EncodeAbiStruct[ccipdata.CommitOnchainConfig](ccipdata.CommitOnchainConfig{
		PriceRegistry: pr,
	})
	require.NoError(t, err)

	commonOffchain := ccipdata.CommitOffchainConfig{
		SourceFinalityDepth:    1,
		GasPriceDeviationPPB:   1e6,
		GasPriceHeartBeat:      1 * time.Hour,
		TokenPriceDeviationPPB: 1e6,
		TokenPriceHeartBeat:    1 * time.Hour,
		InflightCacheExpiry:    3 * time.Hour,
		DestFinalityDepth:      2,
	}
	maxGas := uint64(1e9)
	offchainConfig, err := ccipconfig.EncodeOffchainConfig[ccipdata.CommitOffchainConfigV1_0_0](ccipdata.CommitOffchainConfigV1_0_0{
		SourceFinalityDepth:   commonOffchain.SourceFinalityDepth,
		DestFinalityDepth:     commonOffchain.DestFinalityDepth,
		FeeUpdateHeartBeat:    models.MustMakeDuration(commonOffchain.GasPriceHeartBeat),
		FeeUpdateDeviationPPB: commonOffchain.GasPriceDeviationPPB,
		MaxGasPrice:           maxGas,
		InflightCacheExpiry:   models.MustMakeDuration(commonOffchain.InflightCacheExpiry),
	})
	require.NoError(t, err)
	_, err = ch.SetOCR2Config(user, signers, transmitters, 1, onchainConfig, 1, []byte{})
	require.NoError(t, err)
	onchainConfig2, err := abihelpers.EncodeAbiStruct[ccipdata.CommitOnchainConfig](ccipdata.CommitOnchainConfig{
		PriceRegistry: pr2,
	})
	require.NoError(t, err)
	offchainConfig2, err := ccipconfig.EncodeOffchainConfig[ccipdata.CommitOffchainConfigV1_2_0](ccipdata.CommitOffchainConfigV1_2_0{
		SourceFinalityDepth:      commonOffchain.SourceFinalityDepth,
		DestFinalityDepth:        commonOffchain.DestFinalityDepth,
		MaxGasPrice:              maxGas,
		GasPriceHeartBeat:        models.MustMakeDuration(commonOffchain.GasPriceHeartBeat),
		DAGasPriceDeviationPPB:   1e7,
		ExecGasPriceDeviationPPB: commonOffchain.GasPriceDeviationPPB,
		TokenPriceDeviationPPB:   commonOffchain.TokenPriceDeviationPPB,
		TokenPriceHeartBeat:      models.MustMakeDuration(commonOffchain.TokenPriceHeartBeat),
		InflightCacheExpiry:      models.MustMakeDuration(commonOffchain.InflightCacheExpiry),
	})
	require.NoError(t, err)
	_, err = ch2.SetOCR2Config(user, signers, transmitters, 1, onchainConfig2, 1, []byte{})
	require.NoError(t, err)
	commitAndGetBlockTs(ec)

	// Apply report
	b, err := c10r.EncodeCommitReport(rep)
	require.NoError(t, err)
	_, err = ch.Report(user, b, er)
	require.NoError(t, err)
	b, err = c12r.EncodeCommitReport(rep)
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
	expectedGas := map[string]string{
		ccipdata.V1_0_0: gasPrice.String(),
		ccipdata.V1_2_0: fmt.Sprintf("DA Price: %s, Exec Price: %s", daPrice, gasPrice),
	}
	ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, assets.NewWei(big.NewInt(int64(maxGas)))).Return(gas.EvmFee{Legacy: assets.NewWei(gasPrice)}, uint32(0), nil)
	lm := new(rollupMocks.L1Oracle)
	lm.On("GasPrice", mock.Anything).Return(assets.NewWei(daPrice), nil)
	ge.On("L1Oracle").Return(lm)

	for v, cr := range crs {
		cr := cr
		t.Run("CommitStoreReader "+v, func(t *testing.T) {

			// Static config.
			cfg, err := cr.GetCommitStoreStaticConfig(context.Background())
			require.NoError(t, err)
			require.NotNil(t, cfg)

			// Assert encoding
			b, err := cr.EncodeCommitReport(rep)
			require.NoError(t, err)
			d, err := cr.DecodeCommitReport(b)
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
			_, err = arm.OwnerUnvoteToCurse(user, nil)
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
			assert.Equal(t, reps[0].Data, rep)

			reps, err = cr.GetCommitReportMatchingSeqNum(context.Background(), rep.Interval.Min, 0)
			require.NoError(t, err)
			require.Len(t, reps, 1)
			assert.Equal(t, reps[0].Data, rep)

			reps, err = cr.GetCommitReportMatchingSeqNum(context.Background(), rep.Interval.Min-1, 0)
			require.NoError(t, err)
			require.Len(t, reps, 0)

			// Sanity
			reps, err = cr.GetAcceptedCommitReportsGteTimestamp(context.Background(), time.Unix(0, 0), 0)
			require.NoError(t, err)
			require.Len(t, reps, 1)
			assert.Equal(t, reps[0].Data, rep)

			// Until we detect the config, we'll have empty offchain config
			assert.Equal(t, cr.OffchainConfig(), ccipdata.CommitOffchainConfig{})
			newPr, err := cr.ChangeConfig(configs[v][0], configs[v][1])
			require.NoError(t, err)
			assert.Equal(t, newPr, prs[v])
			assert.Equal(t, commonOffchain, cr.OffchainConfig())
			// We should be able to query for gas prices now.
			gp, err := cr.GasPriceEstimator().GetGasPrice(context.Background())
			require.NoError(t, err)
			assert.Equal(t, expectedGas[v], cr.GasPriceEstimator().String(gp))

		})
	}
}
