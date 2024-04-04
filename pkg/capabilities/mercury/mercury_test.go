package mercury_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
)

func TestFeedID_Validate(t *testing.T) {
	noPrefix := mercury.FeedID("012345678901234567890123456789012345678901234567890123456789000000")
	require.Error(t, noPrefix.Validate())

	badLength := mercury.FeedID("0x1234")
	require.Error(t, badLength.Validate())

	badChars := mercury.FeedID("0x123zzz")
	require.Error(t, badChars.Validate())

	notLowercase := mercury.FeedID("0x0001013ebd4ed3f5889FB5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	require.Error(t, notLowercase.Validate())

	correct := mercury.FeedID("0x0001013ebd4ed3f5889fb5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	require.NoError(t, correct.Validate())
}

func TestCodec(t *testing.T) {
	// Test WrapMercuryTriggerEvent
	const testID = "test-id-1"
	const testTimestamp = "2021-01-01T00:00:00Z"
	var testFeedID = mercury.FeedID("0x1111111111111111111100000000000000000000000000000000000000000000")
	const testFullReport = "0x1234"
	const testBenchmarkPrice = int64(2)
	const testObservationTimestamp = int64(3)
	te := capabilities.TriggerEvent{
		TriggerType: "mercury",
		ID:          testID,
		Timestamp:   testTimestamp,
		BatchedPayload: map[string]any{
			testFeedID.String(): mercury.FeedReport{
				FeedID:               string(testFeedID),
				FullReport:           []byte(testFullReport),
				BenchmarkPrice:       testBenchmarkPrice,
				ObservationTimestamp: testObservationTimestamp,
			},
		},
	}
	wrappedTE, err := mercury.Codec{}.WrapMercuryTriggerEvent(te)
	require.NoError(t, err)

	// Test UnwrapMercuryTriggerEvent
	unwrappedTE, err := mercury.Codec{}.UnwrapMercuryTriggerEvent(wrappedTE)
	require.NoError(t, err)
	require.Equal(t, te, unwrappedTE)
}
