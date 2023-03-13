package ocrcommon

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultMetricVec(t *testing.T) {

	mv, err := NewDefaultMetricVec("test_metric", "help_msg", "label1")
	require.NoError(t, err)

	// check that the underlying prom impl was registered correctly
	// the only way to do that at this layer to try unregistering it
	dmv := (mv).(*DefaultMetricVec)
	assert.True(t, prometheus.DefaultRegisterer.Unregister(dmv))

	// ensure GetMetricsWith works as expected
	labels := map[string]string{"label1": "x"}
	metric, err := mv.GetMetricWith(labels)
	require.NoError(t, err)
	metric.Add(1)
	// in order to use the prom test utils, but cast to underlying prom type
	g := (metric).(prometheus.Gauge)
	assert.Equal(t, float64(1), testutil.ToFloat64(g))

}
