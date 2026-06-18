package pipeline_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const (
	linkTransferRevertData = "0xa9059cbb000000000000000000000000526485b5abdd8ae9c6a63548e0215a83e7135e6100000000000000000000000000000000000000000000000db069932ea4fe1400"
	noopCallData           = "0xdeadbeef"
)

type asyncEthTxEnv struct {
	t                *testing.T
	backend          evmtypes.Backend
	linkToken        common.Address
	sender           common.Address
	runner           pipeline.Runner
	orm              pipeline.ORM
	jobORM           job.ORM
	minConfirmations uint32
}

func newAsyncEthTxEnv(t *testing.T) *asyncEthTxEnv {
	t.Helper()

	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)

	owner := evmtestutils.MustNewSimTransactor(t)
	genesisData := gethtypes.GenesisAlloc{
		owner.From: {Balance: assets.Ether(1000).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, 2*ethconfig.Defaults.Miner.GasCeil)

	linkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(owner, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	cfg := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, _ *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})

	keyStore := cltest.NewKeyStore(t, db)
	_, senderAddr := cltest.MustInsertRandomKey(t, keyStore.Eth(), *sqlutil.New(testutils.SimulatedChainID))

	n, err := backend.Client().NonceAt(ctx, owner.From, nil)
	require.NoError(t, err)
	tx := evmtestutils.NewLegacyTransaction(n, senderAddr, assets.Ether(100).ToInt(), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	require.NoError(t, backend.Client().SendTransaction(ctx, signedTx))
	backend.Commit()

	ethClient := evmclient.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		Client:         ethClient,
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	require.NoError(t, pipelineORM.Start(ctx))
	t.Cleanup(func() { require.NoError(t, pipelineORM.Close()) })

	btORM := bridges.NewORM(db)
	jobORM := job.NewORM(db, pipelineORM, btORM, keyStore, lggr)
	httpClient := clhttptest.NewTestLocalOnlyHTTPClient()
	runner := pipeline.NewRunner(
		pipelineORM,
		btORM,
		cfg.JobPipeline(),
		cfg.WebServer(),
		legacyChains,
		keyStore.Eth(),
		keyStore.VRF(),
		lggr,
		httpClient,
		httpClient,
	)
	servicetest.Run(t, runner)

	chain, err := legacyChains.Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	legacyChain, ok := chain.(legacyevm.Chain)
	require.True(t, ok)
	legacyChain.TxManager().RegisterResumeCallback(runner.ResumeRun)

	return &asyncEthTxEnv{
		t:                t,
		backend:          backend,
		linkToken:        linkTokenAddress,
		sender:           senderAddr,
		runner:           runner,
		orm:              pipelineORM,
		jobORM:           jobORM,
		minConfirmations: 2,
	}
}

func cronSpecWithObservationSource(observationSource string) string {
	return fmt.Sprintf(`
type                = "cron"
schemaVersion       = 1
schedule            = "CRON_TZ=UTC * 0 0 1 1 *"
externalJobID       = "%s"
observationSource   = '''
%s
'''
`, uuid.New(), observationSource)
}

func (e *asyncEthTxEnv) runUntilFinished(t *testing.T, data string, failOnRevert bool) pipeline.Run {
	t.Helper()
	ctx := testutils.Context(t)

	dot := fmt.Sprintf(`
submit_tx [type=ethtx to="%s"
	data="%s"
	minConfirmations="%d"
	failOnRevert=%t
	evmChainID="%s"
	from="[\"%s\"]"
]
`, e.linkToken, data, e.minConfirmations, failOnRevert, testutils.SimulatedChainID, e.sender)

	jb, err := cron.ValidatedCronSpec(cronSpecWithObservationSource(dot))
	require.NoError(t, err)
	require.NoError(t, e.jobORM.CreateJob(ctx, &jb))

	vars := pipeline.NewVarsFrom(map[string]any{
		"jobSpec": map[string]any{
			"databaseID":    jb.ID,
			"externalJobID": jb.ExternalJobID,
			"name":          jb.Name.ValueOrZero(),
		},
		"jobRun": map[string]any{"meta": map[string]any{}},
	})

	run := pipeline.NewRun(*jb.PipelineSpec, vars)
	_, err = e.runner.Run(ctx, run, false, nil)
	require.NoError(t, err)
	require.NotZero(t, run.ID)

	require.Eventually(t, func() bool {
		e.backend.Commit()
		stored, findErr := e.orm.FindRun(ctx, run.ID)
		if findErr != nil {
			return false
		}
		return stored.State.Finished()
	}, testutils.WaitTimeout(t), 200*time.Millisecond)

	stored, err := e.orm.FindRun(ctx, run.ID)
	require.NoError(t, err)
	return stored
}

func receiptFromRun(t *testing.T, run pipeline.Run) map[string]any {
	t.Helper()
	require.Len(t, run.Outputs.Val, 1)
	outputs, ok := run.Outputs.Val.([]any)
	require.True(t, ok)
	receipt, ok := outputs[0].(map[string]any)
	require.True(t, ok, "expected receipt map output, got %T", outputs[0])
	return receipt
}

func TestETHTxTask_AsyncIntegration(t *testing.T) {
	t.Run("with failOnRevert false, run completes for successful call", func(t *testing.T) {
		t.Parallel()
		env := newAsyncEthTxEnv(t)
		run := env.runUntilFinished(t, noopCallData, false)
		require.Equal(t, pipeline.RunStatusCompleted, run.State)
		require.Len(t, run.PipelineTaskRuns, 1)
		assert.True(t, run.PipelineTaskRuns[0].Error.IsZero())

		receipt := receiptFromRun(t, run)
		assert.NotEmpty(t, receipt["blockNumber"])
		assert.NotEmpty(t, receipt["gasUsed"])
	})

	t.Run("with failOnRevert true, run errors when transaction reverts", func(t *testing.T) {
		t.Parallel()
		env := newAsyncEthTxEnv(t)
		run := env.runUntilFinished(t, linkTransferRevertData, true)
		require.Equal(t, pipeline.RunStatusErrored, run.State)
		require.Len(t, run.PipelineTaskRuns, 1)
		assert.False(t, run.PipelineTaskRuns[0].Error.IsZero())
	})

	t.Run("with failOnRevert false, run completes with reverted receipt", func(t *testing.T) {
		t.Parallel()
		env := newAsyncEthTxEnv(t)
		run := env.runUntilFinished(t, linkTransferRevertData, false)
		require.Equal(t, pipeline.RunStatusCompleted, run.State)
		require.Len(t, run.PipelineTaskRuns, 1)
		assert.True(t, run.PipelineTaskRuns[0].Error.IsZero())

		receipt := receiptFromRun(t, run)
		assert.Equal(t, "0x0", receipt["status"])
	})
}
