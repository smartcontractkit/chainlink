package queuelimits

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlobQuota(t *testing.T) {
	t.Parallel()

	const (
		timeout  = 10 * time.Minute
		grace    = time.Minute
		perRound = 1 << 20
	)

	t.Run("covers the event lifetime plus reaping slack", func(t *testing.T) {
		t.Parallel()
		count, payload := BlobQuota(time.Second, timeout, grace, perRound)
		// 690 rounds of event lifetime (10m + 1m + 30s minimum) plus 60 rounds of reaping slack.
		assert.Equal(t, uint32(7_500), count)
		assert.Equal(t, uint32(786_432_000), payload)
	})

	t.Run("faster rounds keep proportionally more blobs", func(t *testing.T) {
		t.Parallel()
		count, payload := BlobQuota(200*time.Millisecond, timeout, grace, perRound)
		// 3450 rounds of event lifetime plus 300 rounds of reaping slack.
		assert.Equal(t, uint32(37_500), count)
		assert.Equal(t, uint32(3_932_160_000), payload)
	})

	t.Run("keeps blobs for at least three rounds", func(t *testing.T) {
		t.Parallel()
		count, _ := BlobQuota(time.Hour, 0, 0, 0)
		// The lifetime is MinBlobLifetime, which rounds to zero at a one hour round interval, so the minimum of
		// three applies, plus one reaping round.
		assert.Equal(t, uint32(4*MaxBlobHandlesPerObservation), count)
	})

	t.Run("non-positive round interval means one second", func(t *testing.T) {
		t.Parallel()
		count, payload := BlobQuota(0, timeout, grace, perRound)
		wantCount, wantPayload := BlobQuota(time.Second, timeout, grace, perRound)
		assert.Equal(t, wantCount, count)
		assert.Equal(t, wantPayload, payload)
	})

	t.Run("saturates at the config's uint32 maximum", func(t *testing.T) {
		t.Parallel()
		count, payload := BlobQuota(time.Nanosecond, timeout, grace, math.MaxInt32)
		assert.Equal(t, uint32(math.MaxUint32), count)
		assert.Equal(t, uint32(math.MaxUint32), payload)
	})

	t.Run("negative lifetimes are clamped", func(t *testing.T) {
		t.Parallel()
		count, _ := BlobQuota(time.Second, -time.Hour, -time.Hour, 0)
		// Only MinBlobLifetime remains: 30 rounds at a one second interval, plus 60 rounds of reaping slack.
		assert.Equal(t, uint32(10*(30+60)), count)
	})
}
