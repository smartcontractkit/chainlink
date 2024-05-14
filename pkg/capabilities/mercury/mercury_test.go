package mercury_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
)

const (
	testFeedID               = mercury.FeedID("0x1111111111111111111100000000000000000000000000000000000000000000")
	testFullReport           = "0x1234"
	testObservationTimestamp = int64(3)
)

func TestFeedID_Validate(t *testing.T) {
	_, err := mercury.NewFeedID("012345678901234567890123456789012345678901234567890123456789000000")
	require.Error(t, err)

	_, err = mercury.NewFeedID("0x1234")
	require.Error(t, err)

	_, err = mercury.NewFeedID("0x123zzz")
	require.Error(t, err)

	_, err = mercury.NewFeedID("0x0001013ebd4ed3f5889FB5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	require.Error(t, err)

	_, err = mercury.NewFeedID("0x0001013ebd4ed3f5889fb5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	require.NoError(t, err)
}
