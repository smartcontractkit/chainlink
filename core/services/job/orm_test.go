package job

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/sqlx"
)

func NewTestORM(t *testing.T, db *sqlx.DB, chainSet evm.ChainSet, pipelineORM pipeline.ORM, keyStore keystore.Master, cfg pg.LogConfig) ORM {
	o := NewORM(db, chainSet, pipelineORM, keyStore, logger.TestLogger(t), cfg)
	t.Cleanup(func() { o.Close() })
	return o
}

func TestLoadEnvConfigVarsLocalOCR(t *testing.T) {
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
