package globalconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashGaugeValue(t *testing.T) {
	t.Parallel()
	assert.Equal(t, int64(0xabcdef012345), HashGaugeValue("abcdef0123456789"))
	assert.Equal(t, HashGaugeValue("abcdef012345aaaa"), HashGaugeValue("abcdef012345bbbb"), "only the first 48 bits count")
	assert.Equal(t, int64(0), HashGaugeValue("short"))
	assert.Equal(t, int64(0), HashGaugeValue("zzzzzzzzzzzzzz"))
	assert.Equal(t, int64(0xffffffffffff), HashGaugeValue("ffffffffffffffff"), "fits in an int64 and a float64 exactly")
}

func TestMetrics_NilIsNoop(t *testing.T) {
	t.Parallel()
	var m *Metrics
	ctx := t.Context()
	m.RecordApplyError(ctx, "d", "e")
	m.RecordApplied(ctx, "d", "e", "abcdef0123456789", time.Unix(1, 0))
}
