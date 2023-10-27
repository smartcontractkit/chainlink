package job_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func NewTestORM(t *testing.T, db *sqlx.DB, legacyChains evm.LegacyChainContainer, pipelineORM pipeline.ORM, bridgeORM bridges.ORM, keyStore keystore.Master, cfg pg.QConfig) job.ORM {
	o := job.NewORM(db, legacyChains, pipelineORM, bridgeORM, keyStore, logger.TestLogger(t), cfg)
	t.Cleanup(func() { o.Close() })
	return o
}

func TestLoadEnvConfigVarsLocalOCR(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	chainConfig := evmtest.NewChainScopedConfig(t, config)
	jobSpec := &job.OCROracleSpec{}

	jobSpec = job.LoadEnvConfigVarsLocalOCR(chainConfig.EVM().OCR(), *jobSpec, chainConfig.OCR())

	require.True(t, jobSpec.ObservationTimeoutEnv)
	require.True(t, jobSpec.BlockchainTimeoutEnv)
	require.True(t, jobSpec.ContractConfigTrackerSubscribeIntervalEnv)
	require.True(t, jobSpec.ContractConfigTrackerPollIntervalEnv)
	require.True(t, jobSpec.ContractConfigConfirmationsEnv)
	require.True(t, jobSpec.DatabaseTimeoutEnv)
	require.True(t, jobSpec.ObservationGracePeriodEnv)
	require.True(t, jobSpec.ContractTransmitterTransmitTimeoutEnv)
}

func TestLoadEnvConfigVarsDR(t *testing.T) {
	t.Parallel()

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		hundred := uint32(100)
		c.EVM[0].MinIncomingConfirmations = &hundred
	})
	chainConfig := evmtest.NewChainScopedConfig(t, config)

	jobSpec10 := job.DirectRequestSpec{
		MinIncomingConfirmations: clnull.Uint32From(10),
	}

	drs10 := job.LoadEnvConfigVarsDR(chainConfig.EVM().MinIncomingConfirmations(), jobSpec10)
	assert.Equal(t, uint32(100), drs10.MinIncomingConfirmations.Uint32)

	jobSpec200 := job.DirectRequestSpec{
		MinIncomingConfirmations: clnull.Uint32From(200),
	}

	drs200 := job.LoadEnvConfigVarsDR(chainConfig.EVM().MinIncomingConfirmations(), jobSpec200)
	assert.True(t, drs200.MinIncomingConfirmations.Valid)
	assert.Equal(t, uint32(200), drs200.MinIncomingConfirmations.Uint32)
}
