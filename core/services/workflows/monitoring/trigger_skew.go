package monitoring

import (
	"context"
	"sync"
	"time"
)

type triggerSkewContextKey struct{}

// Module load sources for platform_engine_trigger_queue_to_execution_start_seconds.
// Aligns with platform_workflow_module_cache_reload_total where applicable.
const (
	ModuleLoadSourceLoaded  = "loaded"   // module already pinned in L1 (no reload counter)
	ModuleLoadSourceWeakRef = "weak_ref" // resurrected from in-memory weak ref
	ModuleLoadSourceDisk    = "disk"     // read binary from disk and re-instantiated
	ModuleLoadSourceDirect  = "direct"   // non-evictable module wrapper (no cache path)
)

// TriggerSkewRecorder records enqueue-to-execution-ready skew once per execution.
type TriggerSkewRecorder struct {
	EnqueueTime time.Time
	Record      func(ctx context.Context, seconds float64, source string)
	once        sync.Once
}

func ContextWithTriggerSkewRecorder(ctx context.Context, rec *TriggerSkewRecorder) context.Context {
	if rec == nil {
		return ctx
	}
	return context.WithValue(ctx, triggerSkewContextKey{}, rec)
}

func TriggerSkewRecorderFromContext(ctx context.Context) *TriggerSkewRecorder {
	rec, _ := ctx.Value(triggerSkewContextKey{}).(*TriggerSkewRecorder)
	return rec
}

func (r *TriggerSkewRecorder) RecordReady(ctx context.Context, source string) {
	if r == nil {
		return
	}
	r.once.Do(func() {
		if r.Record == nil {
			return
		}
		r.Record(ctx, time.Since(r.EnqueueTime).Seconds(), source)
	})
}
