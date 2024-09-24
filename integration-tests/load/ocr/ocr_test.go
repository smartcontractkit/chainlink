package ocr

import (
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/crib"
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

	config, err := tc.GetConfig([]string{"Load"}, tc.OCR)
	require.NoError(t, err)

	sethClient, msClient, bootstrapNode, workerNodes, _, err := crib.ConnectRemote()
	require.NoError(t, err)

	lta, err := actions.SetupOCRv1Cluster(l, sethClient, config.OCR, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := actions.SetupOCRv1Feed(l, sethClient, lta, config.OCR, msClient, bootstrapNode, workerNodes)
	require.NoError(t, err)

	cfg := config.OCR
	cfgl := config.Logging.Loki
	actions.SimulateOCRv1EAActivity(l, cfg.Load.EAChangeInterval.Duration, ocrInstances, workerNodes, msClient)

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:                     t,
		GenName:               "ocr",
		LoadType:              wasp.RPS,
		CallTimeout:           cfg.Load.VerificationTimeout.Duration,
		RateLimitUnitDuration: cfg.Load.RateLimitUnitDuration.Duration,
		Schedule:              wasp.Plain(*cfg.Load.Rate, cfg.Load.TestDuration.Duration),
		Gun:                   NewGun(l, sethClient, ocrInstances),
		Labels:                CommonTestLabels,
		LokiConfig:            wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
	}))
	_, err = p.Run(true)
	require.NoError(t, err)
}

func TestOCRVolume(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Volume"}, tc.OCR)
	require.NoError(t, err)

	sethClient, msClient, bootstrapNode, workerNodes, _, err := crib.ConnectRemote()
	require.NoError(t, err)

	lta, err := actions.SetupOCRv1Cluster(l, sethClient, config.OCR, workerNodes)
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
		VU:          NewVU(l, sethClient, cfg, *cfg.Volume.VURequestsPerUnit, cfg.Volume.RateLimitUnitDuration.Duration, lta, bootstrapNode, workerNodes, msClient),
		Labels:      CommonTestLabels,
		LokiConfig:  wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
	}))
	_, err = p.Run(true)
	require.NoError(t, err)
}
