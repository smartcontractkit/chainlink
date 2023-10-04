package loop

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/services"
)

const KeepAliveTickDuration = keepAliveTickDuration

// TestHook returns a TestPluginService.
// It must only be called once, and before Start.
func (s *pluginService[P, S]) TestHook() TestPluginService[P, S] {
	s.testInterrupt = make(chan func(*pluginService[P, S]))
	return s.testInterrupt
}

// TestPluginService supports Killing & Resetting a running *pluginService.
type TestPluginService[P grpcPlugin, S services.Service] chan<- func(*pluginService[P, S])

func (ch TestPluginService[P, S]) Kill() {
	done := make(chan struct{})
	ch <- func(s *pluginService[P, S]) {
		defer close(done)
		_ = s.closeClient()
	}
	<-done
}

func (ch TestPluginService[P, S]) Reset() {
	done := make(chan struct{})
	ch <- func(r *pluginService[P, S]) {
		defer close(done)
		_ = r.closeClient()
		r.client = nil
		r.clientProtocol = nil
	}
	<-done
}
