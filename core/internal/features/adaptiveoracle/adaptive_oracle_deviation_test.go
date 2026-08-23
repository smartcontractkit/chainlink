package adaptiveoracle

import (
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	deviationInitialAdaptiveRate = 1_000_000_000 // "1.0" at 9 decimals
	deviationMarketRate          = 900_000_000   // "0.9" -- fixed for the entire test, never changed
	deviationAlphaPPB            = 10_000_000    // 1% deviation threshold
	deviationHeartbeat           = 2 * time.Minute
)

// TestIntegration_AdaptiveOracle_DeviationConvergence answers a specific design question: since
// DualAggregator::latestTransmissionDetails() now returns the ADAPTIVE answer (not the original
// market answer -- see that function's doc comment for why), the off-chain median plugin's
// deviation check compares each node's fresh market observation against the adaptive rate, not the
// last-written market price. That means: does the DON keep resubmitting on its own, purely because
// the adaptive rate hasn't caught up to the real market rate yet, and does it correctly STOP once
// it has, without us ever touching the reported market price?
//
// The scenario mirrors a concrete walk-through: current (reset) answer is 1.0, the DON observes a
// real market rate of 0.9, so it writes an update; AdaptiveRateLogic computes a new adaptive rate
// of 0.95 (avg(1.0, 0.9)); 0.95 is still > 1% away from 0.9, so the DON immediately tries again;
// this repeats, each step closing the gap, until the adaptive rate is within the configured 1%
// deviation threshold of 0.9 -- at which point the DON must stop updating on its own.
func TestIntegration_AdaptiveOracle_DeviationConvergence(t *testing.T) {
	t.Parallel()

	owner, b, nodeConfig, contracts := SetupAdaptiveOracleContracts(t, deviationInitialAdaptiveRate)

	// Establish a known starting adaptive rate ("current answer is 1.0").
	_, err := contracts.AdaptiveOracle.ResetAnchors(owner)
	require.NoError(t, err)
	b.Commit()
	startingAnswer, err := contracts.DualAggregator.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(deviationInitialAdaptiveRate), startingAnswer, "resetAnchors must establish the starting adaptive rate")

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

	// alphaPPB=1% and a long DeltaC heartbeat: any report that lands must be deviation-triggered,
	// not heartbeat-triggered, so "it stopped" is actually meaningful.
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
		deviationAlphaPPB,
		deviationHeartbeat,
	)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

	servers := make([]*httptest.Server, numOracles)
	jids := make([]int32, 0, numOracles)
	for i := range numOracles {
		s := i
		require.NoError(t, apps[i].Start(t.Context()))

		// The market rate is fixed for the entire test -- any further on-chain movement must come
		// purely from the DON's own deviation-vs-adaptive-rate retriggering, not from us.
		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			_, _ = io.ReadAll(req.Body)
			res.WriteHeader(http.StatusOK)
			_, err := res.Write(fmt.Appendf(nil, `{"data":%d}`, deviationMarketRate))
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

	// Poll DualAggregator's reported adaptive rate, recording it every time it changes, until it
	// hasn't changed for `stableFor` -- i.e. the DON has stopped updating it.
	history := []*big.Int{startingAnswer}
	stableSince := time.Now()
	require.Eventually(t, func() bool {
		current, err := contracts.DualAggregator.LatestAnswer(nil)
		require.NoError(t, err)

		if current.Cmp(history[len(history)-1]) != 0 {
			history = append(history, new(big.Int).Set(current))
			stableSince = time.Now()
			t.Logf("adaptive rate moved to %s (deviation-triggered report #%d)", current.String(), len(history)-1)
			return false
		}
		// Require at least one real movement before we're willing to call it "stable" -- otherwise
		// we'd trivially "converge" during the warmup period before the first report lands.
		return len(history) > 1 && time.Since(stableSince) >= 5*time.Second
	}, 3*time.Minute, 200*time.Millisecond, "DON must converge the adaptive rate toward the market rate and then stop")

	require.Greater(t, len(history), 2, "expected multiple distinct deviation-triggered reports, not just one jump")
	require.Equal(t, big.NewInt(950_000_000), history[1], "first report: avg(1.0, 0.9) = 0.95, exactly as expected")

	// Each successive report must move strictly closer to the market rate than the last -- proving
	// genuine convergence, not oscillation or divergence.
	marketRate := big.NewInt(deviationMarketRate)
	prevDistance := new(big.Int).Abs(new(big.Int).Sub(history[0], marketRate))
	for i := 1; i < len(history); i++ {
		distance := new(big.Int).Abs(new(big.Int).Sub(history[i], marketRate))
		require.Less(t, distance.Cmp(prevDistance), 0, "report #%d must be strictly closer to the market rate than the previous one", i)
		prevDistance = distance
	}

	settledAnswer := history[len(history)-1]
	require.False(t, exceedsDeviationThreshold(settledAnswer, marketRate, deviationAlphaPPB),
		"the DON must only stop once the adaptive rate is within the configured deviation threshold of the market rate")

	settledRound, err := contracts.DualAggregator.LatestRound(nil)
	require.NoError(t, err)

	// It must actually STAY stopped: wait well past the point where a wrongly-computed deviation
	// would have triggered another report, and confirm nothing further landed. DeltaC (the
	// heartbeat) is configured to 2 minutes, comfortably longer than this wait, so any change here
	// could only be deviation-triggered.
	time.Sleep(15 * time.Second)

	finalAnswer, err := contracts.DualAggregator.LatestAnswer(nil)
	require.NoError(t, err)
	require.Equal(t, settledAnswer, finalAnswer, "adaptive rate must not move once within the deviation threshold")

	finalRound, err := contracts.DualAggregator.LatestRound(nil)
	require.NoError(t, err)
	require.Equal(t, settledRound, finalRound, "no further round should land once the DON has stopped updating")
}

// exceedsDeviationThreshold mirrors the off-chain median plugin's default deviation check:
// abs((offchainMedian - contractMedian) / contractMedian) >= alphaPPB / 1e9.
func exceedsDeviationThreshold(contractMedian, offchainMedian *big.Int, alphaPPB uint64) bool {
	diff := new(big.Int).Abs(new(big.Int).Sub(offchainMedian, contractMedian))
	// diff/contractMedian >= alphaPPB/1e9  <=>  diff*1e9 >= alphaPPB*contractMedian
	lhs := new(big.Int).Mul(diff, big.NewInt(1_000_000_000))
	rhs := new(big.Int).Mul(big.NewInt(0).SetUint64(alphaPPB), contractMedian)
	return lhs.Cmp(rhs) >= 0
}
