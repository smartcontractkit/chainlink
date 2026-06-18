package v2

import (
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

func benchLoadedModule(wfID string) *EvictableModule {
	em := NewEvictableModule(&fakeModule{}, &host.ModuleConfig{}, nil, wfID, "", nil, nil, 1024)
	em.started.Store(true)
	return em
}

func BenchmarkModuleLRU_Register(b *testing.B) {
	clock := clockwork.NewFakeClock()
	lru := NewModuleLRU(clock, WithIdleTimeout(time.Hour))
	em := benchLoadedModule("wf-bench")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wfID := fmt.Sprintf("wf-%d", i)
		lru.Register(wfID, em)
	}
}

func BenchmarkModuleLRU_Contains(b *testing.B) {
	clock := clockwork.NewFakeClock()
	lru := NewModuleLRU(clock, WithIdleTimeout(time.Hour))
	const n = 256
	for i := 0; i < n; i++ {
		wfID := fmt.Sprintf("wf-%d", i)
		lru.Register(wfID, benchLoadedModule(wfID))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Contains(fmt.Sprintf("wf-%d", i%n))
	}
}

// BenchmarkModuleLRU_reap_cap measures cap enforcement cost (sort.Slice over loaded entries).
// Complexity is O(n log n) which is acceptable for our current scale (hundreds to low thousands of modules) and 30s reap interval.
func BenchmarkModuleLRU_reap_cap(b *testing.B) {
	sizes := []int{8, 64, 256}
	for _, n := range sizes {
		b.Run(fmt.Sprintf("loaded=%d", n), func(b *testing.B) {
			clock := clockwork.NewFakeClock()
			reap := make(chan time.Time, 1)
			done := make(chan struct{}, 1)
			capLimit := n / 2
			if capLimit < 1 {
				capLimit = 1
			}
			lru := NewModuleLRU(clock,
				WithMaxLoadedModules(capLimit),
				WithIdleTimeout(time.Hour),
				WithReapTicker(reap),
				WithOnReaped(done),
			)
			lru.Start()
			defer lru.Close()

			for j := 0; j < n; j++ {
				wfID := fmt.Sprintf("wf-%d", j)
				em := benchLoadedModule(wfID)
				em.lastUsed.Store(clock.Now().Add(-time.Duration(j) * time.Second).UnixNano())
				lru.Register(wfID, em)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reap <- clock.Now()
				<-done
			}
		})
	}
}
