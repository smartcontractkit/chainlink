package job_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func NewTestORM(t *testing.T, ds sqlutil.DataSource, pipelineORM pipeline.ORM, bridgeORM bridges.ORM, keyStore keystore.Master) job.ORM {
	o := job.NewORM(ds, pipelineORM, bridgeORM, keyStore, logger.TestLogger(t))
	t.Cleanup(func() { assert.NoError(t, o.Close()) })
	return o
}

func TestLoadConfigVarsLocalOCR(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	chainConfig := evmtest.NewChainScopedConfig(t, config)
	jobSpec := &job.OCROracleSpec{}

	jobSpec = job.LoadConfigVarsLocalOCR(chainConfig.EVM().OCR(), *jobSpec, config.OCR())

	require.Equal(t, models.Interval(config.OCR().ObservationTimeout()), jobSpec.ObservationTimeout)
	require.Equal(t, models.Interval(config.OCR().BlockchainTimeout()), jobSpec.BlockchainTimeout)
	require.Equal(t, models.Interval(config.OCR().ContractSubscribeInterval()), jobSpec.ContractConfigTrackerSubscribeInterval)
	require.Equal(t, models.Interval(config.OCR().ContractPollInterval()), jobSpec.ContractConfigTrackerPollInterval)
	require.Equal(t, config.OCR().CaptureEATelemetry(), jobSpec.CaptureEATelemetry)

	require.Equal(t, chainConfig.EVM().OCR().ContractConfirmations(), jobSpec.ContractConfigConfirmations)
	require.Equal(t, models.Interval(chainConfig.EVM().OCR().DatabaseTimeout()), *jobSpec.DatabaseTimeout)
	require.Equal(t, models.Interval(chainConfig.EVM().OCR().ObservationGracePeriod()), *jobSpec.ObservationGracePeriod)
	require.Equal(t, models.Interval(chainConfig.EVM().OCR().ContractTransmitterTransmitTimeout()), *jobSpec.ContractTransmitterTransmitTimeout)
}

func TestSetDRMinIncomingConfirmations(t *testing.T) {
	t.Parallel()

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		hundred := uint32(100)
		c.EVM[0].MinIncomingConfirmations = &hundred
	})
	chainConfig := evmtest.NewChainScopedConfig(t, config)

	jobSpec10 := job.DirectRequestSpec{
		MinIncomingConfirmations: clnull.Uint32From(10),
	}

	drs10 := job.SetDRMinIncomingConfirmations(chainConfig.EVM().MinIncomingConfirmations(), jobSpec10)
	assert.Equal(t, uint32(100), drs10.MinIncomingConfirmations.Uint32)

	jobSpec200 := job.DirectRequestSpec{
		MinIncomingConfirmations: clnull.Uint32From(200),
	}

	drs200 := job.SetDRMinIncomingConfirmations(chainConfig.EVM().MinIncomingConfirmations(), jobSpec200)
	assert.True(t, drs200.MinIncomingConfirmations.Valid)
	assert.Equal(t, uint32(200), drs200.MinIncomingConfirmations.Uint32)
}
