package ocr

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/havoc"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/k8s"
)

var (
	CommonTestLabels = map[string]string{
		"branch": "ocr_healthcheck_local",
		"commit": "ocr_healthcheck_local",
	}
)

func runHavocMonkey(t *testing.T, l zerolog.Logger, ctx context.Context) {
	havoc.SetGlobalLogger(l)
	cfg, err := havoc.ReadConfig("config.toml")
	require.NoError(t, err)
	err = havoc.GenerateSpecs(
		"skudasov-crib",
		cfg.Havoc.Monkey.Dir,
		cfg,
	)
	require.NoError(t, err)
	m, err := havoc.NewMonkey(cfg)
	require.NoError(t, err)
	err = m.Run(ctx)
	require.NoError(t, err)
}

func TestOCRLoad(t *testing.T) {
	testCtx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})
	l := logging.GetTestLogger(t)
	cc, msClient, cd, bootstrapNode, workerNodes, err := k8s.ConnectRemote(l)
	require.NoError(t, err)
	lt, err := SetupCluster(cc, cd, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := SetupFeed(cc, msClient, cd, bootstrapNode, workerNodes, lt)
	require.NoError(t, err)
	cfg, err := ReadConfig()
	require.NoError(t, err)
	SimulateEAActivity(l, cfg.Load.EAChangeInterval.Duration(), ocrInstances, workerNodes, msClient)

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:                     t,
		GenName:               "ocr",
		LoadType:              wasp.RPS,
		CallTimeout:           cfg.Load.VerificationTimeout.Duration(),
		RateLimitUnitDuration: cfg.Load.RateLimitUnitDuration.Duration(),
		Schedule:              wasp.Plain(cfg.Load.Rate, cfg.Load.TestDuration.Duration()),
		Gun:                   NewGun(l, cc, ocrInstances),
		Labels:                CommonTestLabels,
		LokiConfig:            wasp.NewEnvLokiConfig(),
	}))
	go runHavocMonkey(t, l, testCtx)
	_, err = p.Run(true)
	require.NoError(t, err)
}

func TestOCRVolume(t *testing.T) {
	testCtx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})
	l := logging.GetTestLogger(t)
	cc, msClient, cd, bootstrapNode, workerNodes, err := k8s.ConnectRemote(l)
	require.NoError(t, err)
	lt, err := SetupCluster(cc, cd, workerNodes)
	require.NoError(t, err)
	cfg, err := ReadConfig()
	require.NoError(t, err)

	p := wasp.NewProfile()
	p.Add(wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "ocr",
		LoadType:    wasp.VU,
		CallTimeout: cfg.Volume.VerificationTimeout.Duration(),
		Schedule:    wasp.Plain(cfg.Volume.Rate, cfg.Volume.TestDuration.Duration()),
		VU:          NewVU(l, cfg.Volume.VURequestsPerUnit, cfg.Volume.RateLimitUnitDuration.Duration(), cc, lt, cd, bootstrapNode, workerNodes, msClient),
		Labels:      CommonTestLabels,
		LokiConfig:  wasp.NewEnvLokiConfig(),
	}))
	go runHavocMonkey(t, l, testCtx)
	_, err = p.Run(true)
	require.NoError(t, err)
}
