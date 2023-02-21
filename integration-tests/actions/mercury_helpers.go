package actions

//revive:disable:dot-imports
import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/stretchr/testify/require"
)

func BuildMercuryOCR2Config(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
) contracts.OCRConfig {
	onchainConfig, err := (median.StandardOnchainConfigCodec{}).Encode(median.OnchainConfig{median.MinValue(), median.MaxValue()})
	require.NoError(t, err, "Shouldn't fail encoding config")

	alphaPPB := uint64(1000)

	return BuildGeneralOCR2Config(
		t,
		chainlinkNodes,
		2*time.Second,        // deltaProgress time.Duration,
		20*time.Second,       // deltaResend time.Duration,
		100*time.Millisecond, // deltaRound time.Duration,
		0,                    // deltaGrace time.Duration,
		1*time.Minute,        // deltaStage time.Duration,
		100,                  // rMax uint8,
		[]int{len(chainlinkNodes)},
		median.OffchainConfig{
			false,
			alphaPPB,
			false,
			alphaPPB,
			0,
		}.Encode(),
		0*time.Millisecond,   // maxDurationQuery time.Duration,
		250*time.Millisecond, // maxDurationObservation time.Duration,
		250*time.Millisecond, // maxDurationReport time.Duration,
		250*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
		250*time.Millisecond, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                    // f int,
		onchainConfig,
	)
}
