package ccip

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	sourceChainId = 1337
	destChainId   = 2337
	srcChainName  = "sourceChain"
	destChainName = "destChain"
)

func Test_SequenceNumbers(t *testing.T) {
	t.Parallel()
	collector, _ := NewPluginMetricsCollector("test", sourceChainId, destChainId, srcChainName, destChainName)

	collector.SequenceNumber(Report, 10, "0xabc")
	assert.Equal(t, float64(10), testutil.ToFloat64(maxSequenceNumber.WithLabelValues("test", "1337", "2337", "report", "0xabc", "sourceChain", "destChain")))

	collector.SequenceNumber(Report, 0, "0xabc")
	assert.Equal(t, float64(10), testutil.ToFloat64(maxSequenceNumber.WithLabelValues("test", "1337", "2337", "report", "0xabc", "sourceChain", "destChain")))
}

func Test_NumberOfMessages(t *testing.T) {
	t.Parallel()
	collector, _ := NewPluginMetricsCollector("test", sourceChainId, destChainId, srcChainName, destChainName)
	collector2, _ := NewPluginMetricsCollector("test2", destChainId, sourceChainId, destChainName, srcChainName)

	collector.NumberOfMessagesBasedOnInterval(Observation, 1, 10)
	assert.Equal(t, float64(10), testutil.ToFloat64(messagesProcessed.WithLabelValues("test", "1337", "2337", "observation", "sourceChain", "destChain")))

	collector.NumberOfMessagesBasedOnInterval(Report, 5, 30)
	assert.Equal(t, float64(26), testutil.ToFloat64(messagesProcessed.WithLabelValues("test", "1337", "2337", "report", "sourceChain", "destChain")))

	collector2.NumberOfMessagesProcessed(Report, 15)
	assert.Equal(t, float64(15), testutil.ToFloat64(messagesProcessed.WithLabelValues("test2", "2337", "1337", "report", "destChain", "sourceChain")))
}

func Test_UnexpiredCommitRoots(t *testing.T) {
	t.Parallel()
	collector, _ := NewPluginMetricsCollector("test", sourceChainId, destChainId, srcChainName, destChainName)

	collector.UnexpiredCommitRoots(10)
	assert.Equal(t, float64(10), testutil.ToFloat64(unexpiredCommitRoots.WithLabelValues("test", "1337", "2337", "sourceChain", "destChain")))

	collector.UnexpiredCommitRoots(5)
	assert.Equal(t, float64(5), testutil.ToFloat64(unexpiredCommitRoots.WithLabelValues("test", "1337", "2337", "sourceChain", "destChain")))
}
