package ocr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink/integration-tests/k8s"
)

var (
	CommonTestLabels = map[string]string{
		"branch": "ocr_healthcheck_local",
		"commit": "ocr_healthcheck_local",
	}
)

func TestOCRLoad(t *testing.T) {
	l := logging.GetTestLogger(t)
	cc, msClient, cd, bootstrapNode, workerNodes, err := k8s.ConnectRemote(l)
	require.NoError(t, err)
	lt, err := SetupCluster(cc, cd, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := SetupFeed(cc, msClient, cd, bootstrapNode, workerNodes, lt)
	require.NoError(t, err)

	config, err := tc.GetConfig("Load", tc.OCR)
	require.NoError(t, err)

	cfg := config.OCR
	SimulateEAActivity(l, cfg.Volume.EAChangeInterval.Duration, ocrInstances, workerNodes, msClient)

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:                     t,
		GenName:               "ocr",
		LoadType:              wasp.RPS,
		CallTimeout:           cfg.Volume.VerificationTimeout.Duration,
		RateLimitUnitDuration: cfg.Volume.RateLimitUnitDuration.Duration,
		Schedule:              wasp.Plain(*cfg.Volume.Rate, cfg.Common.TestDuration.Duration),
		Gun:                   NewGun(l, cc, ocrInstances),
		Labels:                CommonTestLabels,
		LokiConfig:            tc.LokiConfigFromToml(&config),
	}))
	_, err = p.Run(true)
	require.NoError(t, err)
}

func TestOCRVolume(t *testing.T) {
	l := logging.GetTestLogger(t)
	cc, msClient, cd, bootstrapNode, workerNodes, err := k8s.ConnectRemote(l)
	require.NoError(t, err)
	lt, err := SetupCluster(cc, cd, workerNodes)
	require.NoError(t, err)
	config, err := tc.GetConfig("Volume", tc.OCR)
	require.NoError(t, err)

	cfg := config.OCR

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "ocr",
		LoadType:    wasp.VU,
		CallTimeout: cfg.Volume.VerificationTimeout.Duration,
		Schedule:    wasp.Plain(*cfg.Volume.Rate, cfg.Common.TestDuration.Duration),
		VU:          NewVU(l, *cfg.Volume.VURequestsPerUnit, cfg.Volume.RateLimitUnitDuration.Duration, cc, lt, cd, bootstrapNode, workerNodes, msClient),
		Labels:      CommonTestLabels,
		LokiConfig:  tc.LokiConfigFromToml(&config),
	}))
	_, err = p.Run(true)
	require.NoError(t, err)
}
