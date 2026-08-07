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
	cappedReferenceRate  = 700
	cappedMarketBelowRef = 500        // < reference: adaptive rate should track it exactly, in one report
	cappedMarketAboveRef = 900        // > reference: adaptive rate should be capped at the reference rate
	cappedAlphaPPB       = 10_000_000 // 1% deviation threshold
	cappedHeartbeat      = 2 * time.Minute
)

// TestIntegration_AdaptiveOracle_CappedLogic answers the same on-chain/off-chain integration
// question as TestIntegration_AdaptiveOracle, but with AdaptiveOracle wired to
// CappedAdaptiveRateLogic instead of AdaptiveRateLogic. Because CappedAdaptiveRateLogic has no
// memory of its own -- it always serves min(marketRate, referenceRate) directly -- a single
// landed report should bring DualAggregator::latestTransmissionDetails() (which returns the
// adaptive answer, see that function's doc comment) fully in line with the market rate. That
// means, unlike the geometric-convergence AdaptiveRateLogic, the DON should NOT keep
// resubmitting on its own afterward: the very first report already satisfies the off-chain
// deviation check.
func TestIntegration_AdaptiveOracle_CappedLogic(t *testing.T) {
	t.Parallel()

	owner, b, nodeConfig, contracts := SetupAdaptiveOracleContractsWithCappedLogic(t, cappedReferenceRate)

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

	// A long DeltaC heartbeat, same rationale as the deviation-convergence test: any report that
	// lands after the first must be deviation-triggered, not heartbeat-triggered, so "it didn't
	// report again" is actually meaningful.
	blockBeforeConfig := InitAdaptiveOracleWithDeltaC(
		t, lggr, b, contracts.DualAggregator, owner, bootstrapNode, oracles, transmitters, transmitters,
		func(blockNum int64) string {
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
		},
		cappedAlphaPPB,
		cappedHeartbeat,
	)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

	var marketRateLock sync.Mutex
	marketRate := int64(cappedMarketBelowRef)
	servers := make([]*httptest.Server, numOracles)
	jids := make([]int32, 0, numOracles)
	for i := range numOracles {
		s := i
		require.NoError(t, apps[i].Start(t.Context()))

		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			_, _ = io.ReadAll(req.Body)
			marketRateLock.Lock()
			price := marketRate
			marketRateLock.Unlock()
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

	// First report: market rate (500) is below the reference rate (700), so the adaptive rate
	// should land at exactly 500 -- no multi-round convergence needed.
	require.Eventually(t, func() bool {
		answer, err := contracts.DualAggregator.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.Cmp(big.NewInt(cappedMarketBelowRef)) == 0
	}, 2*time.Minute, 200*time.Millisecond, "expected DualAggregator to report the market rate directly, uncapped")

	settledRound, err := contracts.DualAggregator.LatestRound(nil)
	require.NoError(t, err)

	// With the adaptive rate already exactly equal to the market rate, the off-chain deviation
	// check has nothing left to trigger on. Confirm no further report lands even though the DON
	// keeps observing the same (already-matching) market rate.
	time.Sleep(15 * time.Second)

	finalAnswer, err := contracts.DualAggregator.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(cappedMarketBelowRef), finalAnswer, "adaptive rate must not move once it already equals the market rate")

	finalRound, err := contracts.DualAggregator.LatestRound(nil)
	require.NoError(t, err)
	require.Equal(t, settledRound, finalRound, "capped logic converges in a single report; no further round should land")

	// Now push the market rate above the reference rate. AdaptiveRateLogic (via
	// CappedAdaptiveRateLogic) should clamp the adaptive rate at the reference rate, and
	// AdaptiveOracle's own min(adaptive, reference) clamp -- which is redundant here but exercised
	// independently -- should agree.
	marketRateLock.Lock()
	marketRate = cappedMarketAboveRef
	marketRateLock.Unlock()

	require.Eventually(t, func() bool {
		answer, err := contracts.DualAggregator.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.Cmp(big.NewInt(cappedReferenceRate)) == 0
	}, 2*time.Minute, 200*time.Millisecond, "adaptive rate must be capped at the reference rate once the market rate exceeds it")

	adaptiveOracleAnswer, err := contracts.AdaptiveOracle.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(cappedReferenceRate), adaptiveOracleAnswer)

	transmissionDetails, err := contracts.DualAggregator.LatestTransmissionDetails(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(cappedReferenceRate), transmissionDetails.LatestAnswer,
		"latestTransmissionDetails must reflect the capped adaptive rate, not the uncapped market rate")
}
