package ccip

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	sourceChainId = 1337
	destChainId   = 2337
	srcChainName  = "sourceChain"
	destChainName = "destChain"
)

func Test_SequenceNumbers(t *testing.T) {
	// setup
	t.Parallel()
	var b strings.Builder

	bhClient, err := beholder.NewWriterClient(&b)
	require.NoError(t, err)
	collector, _ := NewPluginMetricsCollector("test", *bhClient, sourceChainId, destChainId, srcChainName, destChainName)

	collector.SequenceNumber(Report, 10, "0xabc")
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(10), testutil.ToFloat64(maxSequenceNumber.WithLabelValues("test", "1337", "2337", "report", "0xabc", "sourceChain", "destChain")))

	collector.SequenceNumber(Report, 0, "0xabc")
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(10), testutil.ToFloat64(maxSequenceNumber.WithLabelValues("test", "1337", "2337", "report", "0xabc", "sourceChain", "destChain")))

	bhClient.Close()
	assert.Contains(t, b.String(), "ccip_max_sequence_number")
}

func Test_NumberOfMessages(t *testing.T) {
	// setup
	t.Parallel()
	var b strings.Builder
	bhClient, err := beholder.NewWriterClient(&b)
	require.NoError(t, err)

	collector, _ := NewPluginMetricsCollector("test", *bhClient, sourceChainId, destChainId, srcChainName, destChainName)
	collector2, _ := NewPluginMetricsCollector("test2", *bhClient, destChainId, sourceChainId, destChainName, srcChainName)

	collector.NumberOfMessagesBasedOnInterval(Observation, 1, 10)
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(10), testutil.ToFloat64(messagesProcessed.WithLabelValues("test", "1337", "2337", "observation", "sourceChain", "destChain")))

	collector.NumberOfMessagesBasedOnInterval(Report, 5, 30)
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(26), testutil.ToFloat64(messagesProcessed.WithLabelValues("test", "1337", "2337", "report", "sourceChain", "destChain")))

	collector2.NumberOfMessagesProcessed(Report, 15)
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(15), testutil.ToFloat64(messagesProcessed.WithLabelValues("test2", "2337", "1337", "report", "destChain", "sourceChain")))

	bhClient.Close()
	assert.Contains(t, b.String(), "ccip_number_of_messages_processed")
}

func Test_UnexpiredCommitRoots(t *testing.T) {
	// setup
	t.Parallel()
	var b strings.Builder
	bhClient, err := beholder.NewWriterClient(&b)
	require.NoError(t, err)

	collector, _ := NewPluginMetricsCollector("test", *bhClient, sourceChainId, destChainId, srcChainName, destChainName)

	collector.UnexpiredCommitRoots(10)
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(10), testutil.ToFloat64(unexpiredCommitRoots.WithLabelValues("test", "1337", "2337", "sourceChain", "destChain")))

	collector.UnexpiredCommitRoots(5)
	//nolint:testifylint // no need for indelta
	assert.Equal(t, float64(5), testutil.ToFloat64(unexpiredCommitRoots.WithLabelValues("test", "1337", "2337", "sourceChain", "destChain")))

	bhClient.Close()
	assert.Contains(t, b.String(), "ccip_unexpired_commit_roots")
}
