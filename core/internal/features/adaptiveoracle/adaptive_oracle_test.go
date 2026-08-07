package adaptiveoracle

import (
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/features/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const (
	numOracles    = 4
	referenceRate = 700
	initialPrice  = 500 // avg(700, 500) = 600
	changedPrice  = 550 // avg(700, 550) = 625, used to force a distinct round post-reset
)

// TestIntegration_AdaptiveOracle is a PoC smoke test answering: "does the v0.3-adaptive-oracle
// on-chain design (DualAggregator + AdaptiveOracle + AdaptiveRateLogic, from
// https://github.com/smartcontractkit/svr-auction-don) play nice with the existing off-chain OCR2
// median plugin?" A real 4-oracle DON transmits real signed reports into the DualAggregator; we
// verify the AdaptiveOracle::transformAnswer hook fires on each report, that the original
// (untransformed) market answer is preserved for off-chain node logic, and that resetAnchors takes
// effect immediately but does not freeze subsequent adaptation.
func TestIntegration_AdaptiveOracle(t *testing.T) {
	t.Parallel()

	owner, b, nodeConfig, contracts := SetupAdaptiveOracleContracts(t, referenceRate)

	lggr := logger.TestLogger(t)
	bootstrapNodePort := freeport.GetOne(t)
	bootstrapNode := ocr2.SetupNodeOCR2(t, owner, bootstrapNodePort, false /* useForwarders */, b, nil, nodeConfig)

	oracles := make([]confighelper2.OracleIdentityExtra, 0, numOracles)
	transmitters := make([]common.Address, 0, numOracles)
	kbs := make([]ocr2key.KeyBundle, 0, numOracles)
	apps := make([]*cltest.TestApplication, 0, numOracles)
	ports := freeport.GetN(t, numOracles)
	for i := range numOracles {
		node := ocr2.SetupNodeOCR2(t, owner, ports[i], false /* useForwarders */, b, []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.PeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, nodeConfig)

		kbs = append(kbs, node.KeyBundle)
		apps = append(apps, node.App)
		transmitters = append(transmitters, node.Transmitter)

		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  node.KeyBundle.PublicKey(),
				TransmitAccount:   ocrtypes2.Account(node.Transmitter.String()),
				OffchainPublicKey: node.KeyBundle.OffchainPublicKey(),
				PeerID:            node.PeerID,
			},
			ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
		})
	}

	blockBeforeConfig := InitAdaptiveOracle(t, lggr, b, contracts.DualAggregator, owner, bootstrapNode, oracles, transmitters, transmitters, func(blockNum int64) string {
		return fmt.Sprintf(`
type				= "bootstrap"
name				= "bootstrap"
relay				= "evm"
schemaVersion		= 1
contractID			= "%s"
[relayConfig]
chainID 			= 1337
fromBlock = %d
`, contracts.DualAggregatorAddress, blockNum)
	})

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

	var priceLock sync.Mutex
	currentPrice := int64(initialPrice)
	servers := make([]*httptest.Server, numOracles)
	jids := make([]int32, 0, numOracles)

	for i := range numOracles {
		s := i
		require.NoError(t, apps[i].Start(t.Context()))

		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// Drain the request body before responding -- the bridge task sends metadata in the
			// POST body, and not reading it before writing the response can reset the connection
			// under keep-alive while the client is still uploading.
			_, _ = io.ReadAll(req.Body)

			priceLock.Lock()
			price := currentPrice
			priceLock.Unlock()
			res.WriteHeader(http.StatusOK)
			_, err := res.Write(fmt.Appendf(nil, `{"data":%d}`, price))
			require.NoError(t, err)
		}))
		t.Cleanup(func() { servers[s].Close() })

		u, err := url.Parse(servers[i].URL)
		require.NoError(t, err)
		bridgeName := fmt.Sprintf("bridge%d", i)
		require.NoError(t, apps[i].BridgeORM().CreateBridgeType(t.Context(), &bridges.BridgeType{
			Name: bridges.BridgeName(bridgeName),
			URL:  models.WebURL(*u),
		}))

		ocrJob, err := validate.ValidatedOracleSpecToml(
			t.Context(), apps[i].Config.OCR2(), apps[i].Config.Insecure(),
			medianJobToml(contracts.DualAggregatorAddress, kbs[i].ID(), transmitters[i].String(), bridgeName, blockBeforeConfig.Number().Int64()),
			nil,
		)
		require.NoError(t, err)
		require.NoError(t, apps[i].AddJobV2(t.Context(), &ocrJob))
		jids = append(jids, ocrJob.ID)
	}

	// Wait for all 4 oracles to have completed at least one pipeline run against the initial price.
	// Checked concurrently (not sequentially) so no node's wait is delayed by another's.
	//
	// expectedTaskRuns is 0, not the task count (3: ds1, ds1_parse, answer1): OCR-triggered runs are
	// persisted via ocrcommon.RunResultSaver, which hardcodes saveSuccessfulTaskRuns=false for every
	// run ("OCR runs very frequently so a lot of records are produced and the successful TaskRuns do
	// not provide value") -- so pipeline_task_runs rows are never written for these runs, and
	// pr.PipelineTaskRuns always loads back empty regardless of how many tasks the spec has.
	var wg sync.WaitGroup
	for i := range numOracles {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cltest.WaitForPipelineComplete(t, i, jids[i], 1, 0, apps[i].JobORM(), 2*time.Minute, 100*time.Millisecond)
		}(i)
	}
	wg.Wait()

	// The DualAggregator's own AggregatorInterface view exposes the adaptive rate. lastAdaptiveRate
	// starts at 0 (no reset yet), so avg(0, initialPrice) = avg(0, 500) = 250.
	require.Eventually(t, func() bool {
		answer, err := contracts.DualAggregator.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.Cmp(big.NewInt(initialPrice/2)) == 0
	}, 1*time.Minute, 200*time.Millisecond, "expected DualAggregator to report the adaptive rate")

	// AdaptiveOracle clamps to min(adaptive, reference); since adaptive (250) < reference (700),
	// consumers should also see 250.
	adaptiveOracleAnswer, err := contracts.AdaptiveOracle.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(initialPrice/2), adaptiveOracleAnswer)

	// latestTransmissionDetails() returns the ADAPTIVE answer, not the original market answer --
	// this is what the off-chain median plugin's deviation check reads (see
	// DualAggregator::latestTransmissionDetails for why).
	transmissionDetails, err := contracts.DualAggregator.LatestTransmissionDetails(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(initialPrice/2), transmissionDetails.LatestAnswer)

	// Reset the adaptive rate back to the reference rate -- must take effect immediately.
	_, err = contracts.AdaptiveOracle.ResetAnchors(owner)
	require.NoError(t, err)
	b.Commit()

	dualAnswerAfterReset, err := contracts.DualAggregator.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(referenceRate), dualAnswerAfterReset, "reset must take effect immediately")

	adaptiveOracleAnswerAfterReset, err := contracts.AdaptiveOracle.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(referenceRate), adaptiveOracleAnswerAfterReset)

	// Change the observed market price so the median plugin's deviation threshold triggers a new,
	// distinct round -- proving the adaptive rate keeps moving per AdaptiveRateLogic after a reset,
	// rather than staying frozen at the reference rate. lastAdaptiveRate is now 700 (from the
	// reset), so avg(700, changedPrice) = avg(700, 550) = 625.
	priceLock.Lock()
	currentPrice = changedPrice
	priceLock.Unlock()

	require.Eventually(t, func() bool {
		answer, err := contracts.DualAggregator.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.Cmp(big.NewInt((referenceRate+changedPrice)/2)) == 0
	}, 2*time.Minute, 200*time.Millisecond, "adaptive rate must continue moving after a reset, not stay frozen")

	finalTransmissionDetails, err := contracts.DualAggregator.LatestTransmissionDetails(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt((referenceRate+changedPrice)/2), finalTransmissionDetails.LatestAnswer, "latestTransmissionDetails must reflect the new adaptive rate")
}
