package evm_test

import (
	"encoding/json"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	capabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	pollermocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	relayevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewTestJobORM(t *testing.T, ds sqlutil.DataSource, pipelineORM pipeline.ORM, bridgeORM bridges.ORM, keyStore keystore.Master) job.ORM {
	o := job.NewORM(ds, pipelineORM, bridgeORM, keyStore, logger.TestLogger(t))
	t.Cleanup(func() { assert.NoError(t, o.Close()) })
	return o
}

type requestRoundTracker interface {
	RequestRoundTracker() *evm.RequestRoundTracker
}

// This is a regression test verifying that when we instantiate
// a relayer with a jobID and oracleSpecID, we pass the correct
// ID into medianContract, and in particular the request_round_tracker.
// Previously we erroneously passed in jobID to the request_round_db
// rather than the oracleSpecID, causing foreign key violations as the
// an object with that id did not exist.
func TestMedian_RequestRoundTracker(t *testing.T) {
	lggr := logger.TestLogger(t)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	_, addr := cltest.MustInsertRandomKey(t, keyStore.Eth())

	chain := mocks.NewChain(t)
	chainID := testutils.NewRandomEVMChainID()
	chain.On("ID").Return(chainID)

	evmClient := evmclimocks.NewClient(t)
	poller := pollermocks.NewLogPoller(t)
	txManager := txmmocks.NewMockEvmTxManager(t)
	logBroadcaster := logmocks.NewBroadcaster(t)

	cfg := configtest.NewTestGeneralConfig(t)
	evmCfg := evmtest.NewChainScopedConfig(t, cfg)

	chain.On("Config").Return(evmCfg)
	chain.On("Client").Return(evmClient)
	chain.On("LogPoller").Return(poller)
	chain.On("TxManager").Return(txManager)
	chain.On("LogBroadcaster").Return(logBroadcaster)

	poller.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)

	contractID := gethCommon.HexToAddress("0x03bd0d5d39629423979f8a0e53dbce78c1791ebf")
	relayer, err := relayevm.NewRelayer(testutils.Context(t), lggr, chain, relayevm.RelayerOpts{
		DS:                   db,
		CSAETHKeystore:       keyStore,
		CapabilitiesRegistry: capabilities.NewRegistry(lggr),
	})
	require.NoError(t, err)

	pargs := commontypes.PluginArgs{}

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	borm := bridges.NewORM(db)
	orm := NewTestJobORM(t, db, pipelineORM, borm, keyStore)

	job := &job.Job{
		Type:          job.OffchainReporting2,
		SchemaVersion: 1,
		OCR2OracleSpec: &job.OCR2OracleSpec{
			RelayConfig: map[string]any{
				"sendingKeys": []any{addr.String()},
			},
			P2PV2Bootstrappers: []string{"aBootstrapper"},
		},
	}
	require.NoError(t, orm.CreateJob(tests.Context(t), job))

	relayConfig := evmtypes.RelayConfig{
		ChainID:                big.New(chainID),
		EffectiveTransmitterID: null.StringFrom("something"),
		SendingKeys:            []string{addr.String()},
	}
	rc, err := json.Marshal(&relayConfig)
	rargs := commontypes.RelayArgs{
		ContractID:   contractID.String(),
		RelayConfig:  rc,
		JobID:        job.ID,
		OracleSpecID: job.OCR2OracleSpec.ID,
	}

	require.NoError(t, err)
	md, err := relayer.NewMedianProvider(testutils.Context(t), rargs, pargs)
	require.NoError(t, err)

	// Cast medianContract to requestRoundTracker so we can access the relevant round tracker.
	// and from there its methods to emit a log.
	rrt := md.MedianContract().(requestRoundTracker).RequestRoundTracker()

	logBroadcast := logmocks.NewBroadcast(t)

	rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
	logBroadcast.On("RawLog").Return(rawLog).Maybe()
	logBroadcast.On("String").Return("").Maybe()

	logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	rrt.HandleLog(tests.Context(t), logBroadcast)

	configDigest, epoch, round, err := rrt.LatestRoundRequested(testutils.Context(t), 0)
	require.NoError(t, err)
	assert.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", configDigest.Hex())
	assert.Equal(t, 1, int(epoch))
	assert.Equal(t, 1, int(round))
	require.NoError(t, err)
}
