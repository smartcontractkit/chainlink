package v2

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func counterValue(t *testing.T, c prometheus.Counter) float64 {
	t.Helper()
	var m dto.Metric
	require.NoError(t, c.Write(&m))
	return m.GetCounter().GetValue()
}

func gaugeValue(t *testing.T, g prometheus.Gauge) float64 {
	t.Helper()
	var m dto.Metric
	require.NoError(t, g.Write(&m))
	return m.GetGauge().GetValue()
}

func TestCacheMetricsPrometheusExport(t *testing.T) { //nolint:paralleltest // package-global promauto metrics
	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	reloadDiskBefore := counterValue(t, promModuleCacheReloadTotal.WithLabelValues("disk"))
	reloadWeakRefBefore := counterValue(t, promModuleCacheReloadTotal.WithLabelValues("weak_ref"))
	evictionBefore := counterValue(t, promModuleCacheEvictionTotal)
	versionMismatchBefore := counterValue(t, promModuleCacheVersionMismatchTotal)
	pinExhaustedBefore := counterValue(t, promModuleCachePinExhaustedTotal)
	tryAcquireExhaustedBefore := counterValue(t, promModuleCacheTryAcquireExhaustedTotal)

	ctx := t.Context()
	cm.recordReload(ctx, "disk")
	cm.recordReload(ctx, "weak_ref")
	cm.recordEviction(ctx, 2)
	cm.recordLoaded(ctx, 3)
	cm.recordMemorySaved(ctx, 4096)
	cm.recordVersionMismatch(ctx)
	cm.recordPinExhausted(ctx)
	cm.recordTryAcquireExhausted(ctx)

	require.InDelta(t, reloadDiskBefore+1, counterValue(t, promModuleCacheReloadTotal.WithLabelValues("disk")), 0)
	require.InDelta(t, reloadWeakRefBefore+1, counterValue(t, promModuleCacheReloadTotal.WithLabelValues("weak_ref")), 0)
	require.InDelta(t, evictionBefore+2, counterValue(t, promModuleCacheEvictionTotal), 0)
	require.InDelta(t, 3, gaugeValue(t, promModuleCacheLoaded), 0)
	require.InDelta(t, 4096, gaugeValue(t, promModuleCacheMemorySavedBytes), 0)
	require.InDelta(t, versionMismatchBefore+1, counterValue(t, promModuleCacheVersionMismatchTotal), 0)
	require.InDelta(t, pinExhaustedBefore+1, counterValue(t, promModuleCachePinExhaustedTotal), 0)
	require.InDelta(t, tryAcquireExhaustedBefore+1, counterValue(t, promModuleCacheTryAcquireExhaustedTotal), 0)
}
