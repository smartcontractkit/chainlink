package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func NewTestORM(t *testing.T, db *sqlx.DB, chainSet evm.ChainSet, pipelineORM pipeline.ORM, keyStore keystore.Master, cfg pg.LogConfig) ORM {
	o := NewORM(db, chainSet, pipelineORM, keyStore, logger.TestLogger(t), cfg)
	t.Cleanup(func() { o.Close() })
	return o
}

func TestLoadEnvConfigVarsLocalOCR(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	chainConfig := evmtest.NewChainScopedConfig(t, config)
	jobSpec := &OCROracleSpec{}

	jobSpec = LoadEnvConfigVarsLocalOCR(chainConfig, *jobSpec)

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

	config := configtest.NewTestGeneralConfig(t)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(100)
	chainConfig := evmtest.NewChainScopedConfig(t, config)

	jobSpec10 := DirectRequestSpec{
		MinIncomingConfirmations: clnull.Uint32From(10),
	}

	drs10 := LoadEnvConfigVarsDR(chainConfig, jobSpec10)
	assert.True(t, drs10.MinIncomingConfirmationsEnv)

	jobSpec200 := DirectRequestSpec{
		MinIncomingConfirmations: clnull.Uint32From(200),
	}

	drs200 := LoadEnvConfigVarsDR(chainConfig, jobSpec200)
	assert.False(t, drs200.MinIncomingConfirmationsEnv)
	assert.True(t, drs200.MinIncomingConfirmations.Valid)
	assert.Equal(t, uint32(200), drs200.MinIncomingConfirmations.Uint32)
}
