package ccip

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	sourceChainId = 1337
	destChainId   = 2337
)

func Test_SequenceNumbers(t *testing.T) {
	t.Parallel()
	collector := NewPluginMetricsCollector("test", sourceChainId, destChainId)

	collector.SequenceNumber(Report, 10)
	assert.Equal(t, float64(10), testutil.ToFloat64(sequenceNumberCounter.WithLabelValues("test", "1337", "2337", "report")))

	collector.SequenceNumber(Report, 0)
	assert.Equal(t, float64(10), testutil.ToFloat64(sequenceNumberCounter.WithLabelValues("test", "1337", "2337", "report")))
}

func Test_NumberOfMessages(t *testing.T) {
	t.Parallel()
	collector := NewPluginMetricsCollector("test", sourceChainId, destChainId)
	collector2 := NewPluginMetricsCollector("test2", destChainId, sourceChainId)

	collector.NumberOfMessagesBasedOnInterval(Observation, 1, 10)
	assert.Equal(t, float64(10), testutil.ToFloat64(messagesProcessed.WithLabelValues("test", "1337", "2337", "observation")))

	collector.NumberOfMessagesBasedOnInterval(Report, 5, 30)
	assert.Equal(t, float64(26), testutil.ToFloat64(messagesProcessed.WithLabelValues("test", "1337", "2337", "report")))

	collector2.NumberOfMessagesProcessed(Report, 15)
	assert.Equal(t, float64(15), testutil.ToFloat64(messagesProcessed.WithLabelValues("test2", "2337", "1337", "report")))
}

func Test_UnexpiredCommitRoots(t *testing.T) {
	t.Parallel()
	collector := NewPluginMetricsCollector("test", sourceChainId, destChainId)

	collector.UnexpiredCommitRoots(10)
	assert.Equal(t, float64(10), testutil.ToFloat64(unexpiredCommitRoots.WithLabelValues("test", "1337", "2337")))

	collector.UnexpiredCommitRoots(5)
	assert.Equal(t, float64(5), testutil.ToFloat64(unexpiredCommitRoots.WithLabelValues("test", "1337", "2337")))
}
