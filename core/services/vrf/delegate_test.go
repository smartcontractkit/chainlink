package vrf_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/vrfkey"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"

	log_mocks "github.com/smartcontractkit/chainlink/v2/common/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	vrf_mocks "github.com/smartcontractkit/chainlink/v2/core/services/vrf/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

type vrfUniverse struct {
	jrm          job.ORM
	pr           pipeline.Runner
	prm          pipeline.ORM
	lb           *log_mocks.Broadcaster
	ec           *clienttest.Client
	ks           keystore.Master
	vrfkey       vrfkey.KeyV2
	submitter    common.Address
	txm          *txmgr.TxManager
	hb           heads.Broadcaster
	legacyChains legacyevm.LegacyChainContainer
	cid          big.Int
}

func buildVrfUni(t *testing.T, db *sqlx.DB, cfg chainlink.GeneralConfig) vrfUniverse {
	ctx := testutils.Context(t)
	lb := log_mocks.NewBroadcaster(t)
	servicetest.SetupNoOpMock(lb)
	lb.On("AddDependents", 1).Maybe()
	lb.On("Register", mock.Anything, mock.Anything).Return(func() {}).Maybe()
	ec := cltest.NewEthMocksWithStartupAssertions(t)
	ec.On("ConfiguredChainID").Return(testutils.FixtureChainID)
	ec.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(51), nil).Maybe()
	lggr := logger.TestLogger(t)
	hb := heads.NewBroadcaster(lggr)

	prm := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	btORM := bridges.NewORM(db)
	ks := keystore.NewInMemory(db, commonkeystore.FastScryptParams, lggr.Infof)
	require.NoError(t, ks.Unlock(ctx, testutils.Password))
	_, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	evmKs := keys.NewChainStore(keystore.NewEthSigner(ks.Eth(), ec.ConfiguredChainID()), ec.ConfiguredChainID())
	estimator, err := gas.NewEstimator(logger.TestLogger(t), ec, evmConfig.ChainType(), ec.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	ht := headstest.NewSimulatedHeadTracker(ec, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), ec, lggr, ht, lpOpts)
	txm, err := txmgr.NewTxm(db, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), nil, dbConfig, dbConfig.Listener(), ec, logger.TestLogger(t), lp, evmKs, estimator, nil, nil, false)
	require.NoError(t, err)
	orm := heads.NewORM(*testutils.FixtureChainID, db, 0)
	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(51)))
	jrm := job.NewORM(db, prm, btORM, ks, lggr)
	t.Cleanup(func() { assert.NoError(t, jrm.Close()) })
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		LogBroadcaster: lb,
		Client:         ec,
		DB:             db,
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		TxManager:      txm,
		KeyStore:       ks.Eth(),
	})
	pr := pipeline.NewRunner(prm, btORM, cfg.JobPipeline(), cfg.WebServer(), legacyChains, ks.Eth(), ks.VRF(), lggr, nil, nil)
	require.NoError(t, ks.Unlock(ctx, testutils.Password))
	k, err2 := ks.Eth().Create(testutils.Context(t), testutils.FixtureChainID)
	require.NoError(t, err2)
	submitter := k.Address
	vrfkey, err3 := ks.VRF().Create(ctx)
	require.NoError(t, err3)

	return vrfUniverse{
		jrm:          jrm,
		pr:           pr,
		prm:          prm,
		lb:           lb,
		ec:           ec,
		ks:           ks,
		vrfkey:       vrfkey,
		submitter:    submitter,
		txm:          &txm,
		hb:           hb,
		legacyChains: legacyChains,
		cid:          *ec.ConfiguredChainID(),
	}
}

func Test_CheckFromAddressMaxGasPrices(t *testing.T) {
	t.Run("returns nil error if gasLanePrice not set in job spec", func(tt *testing.T) {
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				PublicKey:        "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800",
				FromAddresses:    []string{"0x1111111111111111111111111111111111111111"},
				OmitGasLanePrice: true,
			}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewFeeConfig(t)
		require.NoError(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg.PriceMaxKey))
	})

	t.Run("returns nil error on valid gas lane <=> key specific gas price setting", func(tt *testing.T) {
		var fromAddresses []string
		for range 3 {
			fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		}

		cfg := vrf_mocks.NewFeeConfig(t)
		for _, a := range fromAddresses {
			cfg.On("PriceMaxKey", common.HexToAddress(a)).Return(assets.GWei(100)).Once()
		}
		defer cfg.AssertExpectations(tt)

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		require.NoError(t, err)

		require.NoError(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg.PriceMaxKey))
	})

	t.Run("returns error on invalid setting", func(tt *testing.T) {
		var fromAddresses []string
		for range 3 {
			fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		}

		cfg := vrf_mocks.NewFeeConfig(t)
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[0])).Return(assets.GWei(100)).Once()
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[1])).Return(assets.GWei(100)).Once()
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[2])).Return(assets.GWei(50)).Once()
		defer cfg.AssertExpectations(tt)

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		require.NoError(t, err)

		require.Error(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg.PriceMaxKey))
	})
}

func Test_CheckFromAddressesExist(t *testing.T) {
	t.Run("from addresses exist", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		lggr := logger.TestLogger(t)
		ks := keystore.NewInMemory(db, commonkeystore.FastScryptParams, lggr.Infof)
		require.NoError(t, ks.Unlock(ctx, testutils.Password))

		var fromAddresses []string
		for range 3 {
			k, err := ks.Eth().Create(testutils.Context(t), big.NewInt(1337))
			assert.NoError(t, err)
			fromAddresses = append(fromAddresses, k.Address.Hex())
		}
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		assert.NoError(t, err)

		assert.NoError(t, vrf.CheckFromAddressesExist(testutils.Context(t), jb, ks.Eth()))
	})

	t.Run("one of from addresses doesn't exist", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		lggr := logger.TestLogger(t)
		ks := keystore.NewInMemory(db, commonkeystore.FastScryptParams, lggr.Infof)
		require.NoError(t, ks.Unlock(ctx, testutils.Password))

		var fromAddresses []string
		for range 3 {
			k, err := ks.Eth().Create(testutils.Context(t), big.NewInt(1337))
			assert.NoError(t, err)
			fromAddresses = append(fromAddresses, k.Address.Hex())
		}
		fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		assert.NoError(t, err)

		assert.Error(t, vrf.CheckFromAddressesExist(testutils.Context(t), jb, ks.Eth()))
	})
}

func Test_FromAddressMaxGasPricesAllEqual(t *testing.T) {
	t.Run("all max gas prices equal", func(tt *testing.T) {
		fromAddresses := []string{
			"0x498C2Dce1d3aEDE31A8c808c511C38a809e67684",
			"0x253b01b9CaAfbB9dC138d7D8c3ACBCDd47144b4B",
			"0xD94E6AD557277c6E3e163cefF90F52AB51A95143",
		}

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewFeeConfig(t)
		for _, a := range fromAddresses {
			cfg.On("PriceMaxKey", common.HexToAddress(a)).Return(assets.GWei(100))
		}
		defer cfg.AssertExpectations(tt)

		assert.True(tt, vrf.FromAddressMaxGasPricesAllEqual(jb, cfg.PriceMaxKey))
	})

	t.Run("one max gas price not equal to others", func(tt *testing.T) {
		fromAddresses := []string{
			"0x498C2Dce1d3aEDE31A8c808c511C38a809e67684",
			"0x253b01b9CaAfbB9dC138d7D8c3ACBCDd47144b4B",
			"0xD94E6AD557277c6E3e163cefF90F52AB51A95143",
			"0x86E7c45Bf013Bf1Df3C22c14d5fd6fc3051AC569",
		}

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewFeeConfig(t)
		for _, a := range fromAddresses[:3] {
			cfg.On("PriceMaxKey", common.HexToAddress(a)).Return(assets.GWei(100))
		}
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[len(fromAddresses)-1])).
			Return(assets.GWei(200))
		defer cfg.AssertExpectations(tt)

		assert.False(tt, vrf.FromAddressMaxGasPricesAllEqual(jb, cfg.PriceMaxKey))
	})
}

func Test_VRFV2PlusServiceFailsWhenVRFOwnerProvided(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	vuni := buildVrfUni(t, db, cfg)

	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

	vd := vrf.NewDelegate(
		db,
		vuni.ks,
		vuni.pr,
		vuni.prm,
		vuni.legacyChains,
		logger.TestLogger(t),
		mailMon)
	chainService, err := vuni.legacyChains.Get(testutils.FixtureChainID.String())
	require.NoError(t, err)
	chain, ok := chainService.(legacyevm.Chain)
	require.True(t, ok)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		VRFVersion:    vrfcommon.V2Plus,
		PublicKey:     vuni.vrfkey.PublicKey.String(),
		FromAddresses: []string{vuni.submitter.Hex()},
		GasLanePrice:  chain.Config().EVM().GasEstimator().PriceMax(),
		EVMChainID:    testutils.FixtureChainID.String(),
	})
	toml := "vrfOwnerAddress=\"0xF62fEFb54a0af9D32CDF0Db21C52710844c7eddb\"\n" + vs.Toml()
	jb, err := vrfcommon.ValidatedVRFSpec(toml)
	require.NoError(t, err)
	ctx := testutils.Context(t)
	err = vuni.jrm.CreateJob(ctx, &jb)
	require.NoError(t, err)
	_, err = vd.ServicesForSpec(testutils.Context(t), jb)
	require.Error(t, err)
	require.Equal(t, "VRF Owner is not supported for VRF V2 Plus", err.Error())
}
