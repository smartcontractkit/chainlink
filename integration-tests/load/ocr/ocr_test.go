package ocr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/k8s"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	CommonTestLabels = map[string]string{
		"branch": "ocr_healthcheck_local",
		"commit": "ocr_healthcheck_local",
	}
)

func TestOCRLoad(t *testing.T) {
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Load", tc.OCR)
	require.NoError(t, err)

	evmNetwork, msClient, bootstrapNode, workerNodes, err := k8s.ConnectRemote()
	require.NoError(t, err)

	seth, err := actions_seth.GetChainClient(config, *evmNetwork)
	require.NoError(t, err, "Error creating seth client")

	lta, err := SetupCluster(l, seth, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := SetupFeed(l, seth, lta, msClient, bootstrapNode, workerNodes)
	require.NoError(t, err)

	cfg := config.OCR
	cfgl := config.Logging.Loki
	SimulateEAActivity(l, cfg.Load.EAChangeInterval.Duration, ocrInstances, workerNodes, msClient)

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:                     t,
		GenName:               "ocr",
		LoadType:              wasp.RPS,
		CallTimeout:           cfg.Load.VerificationTimeout.Duration,
		RateLimitUnitDuration: cfg.Load.RateLimitUnitDuration.Duration,
		Schedule:              wasp.Plain(*cfg.Load.Rate, cfg.Load.TestDuration.Duration),
		Gun:                   NewGun(l, seth, ocrInstances),
		Labels:                CommonTestLabels,
		LokiConfig:            wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
	}))
	_, err = p.Run(true)
	require.NoError(t, err)
}

func TestOCRVolume(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Volume", tc.OCR)
	require.NoError(t, err)

	evmNetwork, msClient, bootstrapNode, workerNodes, err := k8s.ConnectRemote()
	require.NoError(t, err)

	seth, err := actions_seth.GetChainClient(config, *evmNetwork)
	require.NoError(t, err, "Error creating seth client")

	lta, err := SetupCluster(l, seth, workerNodes)
	require.NoError(t, err)

	cfg := config.OCR
	cfgl := config.Logging.Loki

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "ocr",
		LoadType:    wasp.VU,
		CallTimeout: cfg.Volume.VerificationTimeout.Duration,
		Schedule:    wasp.Plain(*cfg.Volume.Rate, cfg.Volume.TestDuration.Duration),
		VU:          NewVU(l, seth, *cfg.Volume.VURequestsPerUnit, cfg.Volume.RateLimitUnitDuration.Duration, lta, bootstrapNode, workerNodes, msClient),
		Labels:      CommonTestLabels,
		LokiConfig:  wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
	}))
	_, err = p.Run(true)
	require.NoError(t, err)
}
