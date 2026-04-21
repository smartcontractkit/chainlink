package cre

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	smokecre "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

/*
Test_CRE_PoR_MemoryLeakSoak is a long-running soak test that:
 1. Registers 20 V2 PoR workflows, each with a unique feed ID and staggered cron schedule.
 2. Runs a time-bounded loop that detects on-chain price updates for every workflow.
 3. After the soak, asserts Prometheus CPU / memory metrics haven't exceeded configured thresholds.

Prerequisites (handled by the CI workflow):
  - Local CRE environment is running (go run . env start).
  - Observability stack is running (./bin/ctf obs up -f) — required by the leak package.

Environment variables:
  - CRE_SOAK_DURATION: duration string (e.g. "2h"). Defaults to 2h.
  - CTF_CONFIGS: path to the environment config TOML file.
*/
func Test_CRE_PoR_MemoryLeakSoak(t *testing.T) {
	const numWorkflows = 20

	start := time.Now()

	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetDefaultTestConfig(t),
	)

	t.Cleanup(func() {
		if t.Failed() {
			_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
			require.NoError(t, cErr)
		}
	})

	// ── Phase 1: Single shared price provider ──────────────────────────────────
	// IMPORTANT: must be called once for ALL feed IDs before any workflow is
	// registered.  Subsequent calls to setupFakeDataProvider would overwrite the
	// HTTP handler and make previously-registered feeds return 400.
	allFeedIDs := smokecre.GenerateSoakFeedIDs(numWorkflows)

	logFileName := fmt.Sprintf("./%s-%s/soak-fake-price-provider.log", framework.DefaultCTFLogsDir, t.Name())
	if _, err := os.Stat(filepath.Dir(logFileName)); os.IsNotExist(err) {
		require.NoError(t, os.MkdirAll(filepath.Dir(logFileName), 0755), "failed to create directory %s", filepath.Dir(logFileName))
	}
	ppLogFile, openErr := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	require.NoError(t, openErr, "failed to open %d", logFileName)
	oldGinDefaultWriter := gin.DefaultWriter
	oldGinDefaultErrorWriter := gin.DefaultErrorWriter
	// fake.NewFakeDataProvider uses gin.Default(), so redirect Gin's global request/error
	// output to the soak provider log file to avoid flooding the main test logs.
	gin.DefaultWriter = ppLogFile
	gin.DefaultErrorWriter = ppLogFile
	t.Cleanup(func() {
		gin.DefaultWriter = oldGinDefaultWriter
		gin.DefaultErrorWriter = oldGinDefaultErrorWriter
	})

	soakPPLogger := zerolog.New(ppLogFile).
		Level(framework.L.GetLevel()).
		With().
		Timestamp().
		Str("component", "soak-fake-price-provider").
		Logger()
	framework.L.Info().Str("path", logFileName).Msg("redirecting soak fake price provider logs to file")

	priceProvider, err := smokecre.NewFakePriceProviderForSoak(
		soakPPLogger, testEnv.Config.Fake, "", allFeedIDs,
	)
	require.NoError(t, err, "failed to create soak price provider")

	// ── Phase 2: Sequential workflow registration ──────────────────────────────
	type soakWorkflow struct {
		feedID       string
		cacheAddress common.Address
		lastPrice    *big.Int
		updateCount  int
	}

	soakWorkflows := make([]soakWorkflow, numWorkflows)
	for i := range numWorkflows {
		// Stagger: first 10 use the default offset, second 10 are shifted by 15s
		// so that not all workflows fire at the same instant.
		cronSchedule := "*/30 * * * * *"
		if i >= 10 {
			cronSchedule = "15/30 * * * * *"
		}

		wfConfig := smokecre.WorkflowTestConfig{
			WorkflowName:         fmt.Sprintf("por-soak-%02d", i),
			WorkflowFileLocation: smokecre.PoRWFV2Location,
			FeedIDs:              []string{allFeedIDs[i]},
			CronSchedule:         cronSchedule,
		}

		cacheAddr, setupErr := smokecre.SetupPoRWorkflowForSoak(t, testEnv, priceProvider, wfConfig)
		require.NoError(t, setupErr, "failed to set up soak workflow %d", i)

		soakWorkflows[i] = soakWorkflow{feedID: allFeedIDs[i], cacheAddress: cacheAddr}
		framework.L.Info().Msgf("Workflow %d/%d registered (%s)", i+1, numWorkflows, wfConfig.WorkflowName)
	}

	// ── Phase 3: Time-bounded verification + soak loop ─────────────────────────
	soakDuration := parseDuration(os.Getenv("CRE_SOAK_DURATION"), 2*time.Hour)
	framework.L.Info().Msgf("Soak duration: %s", soakDuration)

	ctx, cancel := context.WithTimeout(t.Context(), soakDuration)
	defer cancel()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Obtain the EVM client from the first blockchain in the environment
	require.NotEmpty(t, testEnv.CreEnvironment.Blockchains, "no blockchains in environment")
	evmBC, ok := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)
	require.True(t, ok, "first blockchain is not an EVM blockchain")

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			framework.L.Info().Msgf("Soak progress: elapsed=%s / %s", time.Since(start).Round(time.Second), soakDuration)
			for i := range soakWorkflows {
				wf := &soakWorkflows[i]
				latest := pollOnChainPrice(evmBC, wf.cacheAddress, wf.feedID)
				if latest == nil || latest.Sign() == 0 {
					continue
				}
				if wf.lastPrice != nil && latest.Cmp(wf.lastPrice) == 0 {
					continue
				}
				// New price landed on-chain — advance the provider to serve the next one
				priceProvider.NextPrice(wf.feedID, latest, time.Since(start))
				wf.updateCount++
				wf.lastPrice = new(big.Int).Set(latest)
				framework.L.Info().Msgf("  wf[%02d] confirmed update #%d, price=%s", i, wf.updateCount, latest)
			}
		}
	}

	// ── Phase 4: Verify each workflow ran ──────────────────────────────────────
	for i, wf := range soakWorkflows {
		if wf.updateCount == 0 {
			framework.L.Error().Msgf("workflow[%02d] (feedID=%s) had zero confirmed on-chain updates — never ran", i, wf.feedID)
			t.FailNow()
		} else {
			framework.L.Info().Msgf("workflow[%02d]: %d confirmed on-chain updates", i, wf.updateCount)
		}
	}

	// ── Phase 5: Memory / CPU leak check ──────────────────────────────────────
	leakDetector, ldErr := leak.NewCLNodesLeakDetector(leak.NewResourceLeakChecker(), leak.WithNodesetName("workflow"))
	require.NoError(t, ldErr, "failed to create CL nodes leak detector")

	checkErr := leakDetector.Check(&leak.CLNodesCheck{
		ComparisonMode:  leak.ComparisonModeAbsolute,
		NumNodes:        len(testEnv.Dons.MustWorkflowDON().Nodes),
		Start:           start,
		End:             time.Now(),
		WarmUpDuration:  30 * time.Minute,
		CPUThreshold:    25.0,
		MemoryThreshold: 3350.0, // CRE uses quite a bit of memory
	})
	require.NoError(t, checkErr, "resource leak check failed")
}

// pollOnChainPrice reads the latest answer from a DataFeedsCache contract for the given
// feed ID. Returns nil on any error or when no price is set yet.
func pollOnChainPrice(bc *evm.Blockchain, cacheAddress common.Address, feedID string) *big.Int {
	instance, err := data_feeds_cache.NewDataFeedsCache(cacheAddress, bc.SethClient.Client)
	if err != nil {
		return nil
	}
	// feed ID strings are 32 hex chars (16 bytes); use first 16 bytes via slice conversion
	feedIDBytes := [16]byte(common.Hex2Bytes(feedID))
	price, err := instance.GetLatestAnswer(bc.SethClient.NewCallOpts(), feedIDBytes)
	if err != nil {
		return nil
	}
	return price
}

// parseDuration parses s as a time.Duration. Returns defaultVal when s is empty or invalid.
func parseDuration(s string, defaultVal time.Duration) time.Duration {
	if s == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultVal
	}
	return d
}
