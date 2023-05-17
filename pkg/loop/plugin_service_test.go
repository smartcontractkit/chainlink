package loop

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

const KeepAliveTickDuration = keepAliveTickDuration

// TestHook returns a TestPluginService.
// It must only be called once, and before Start.
func (r *pluginService[P, S]) TestHook() TestPluginService[P, S] {
	r.testInterrupt = make(chan func(*pluginService[P, S]))
	return r.testInterrupt
}

// TestPluginService supports Killing & Resetting a running *pluginService.
type TestPluginService[P grpcPlugin, S types.Service] chan<- func(*pluginService[P, S])

func (ch TestPluginService[P, S]) Kill() {
	done := make(chan struct{})
	ch <- func(s *pluginService[P, S]) {
		defer close(done)
		if s.client != nil {
			s.client.Kill()
		}
	}
	<-done
}

func (ch TestPluginService[P, S]) Reset() {
	done := make(chan struct{})
	ch <- func(r *pluginService[P, S]) {
		defer close(done)
		if r.client != nil {
			r.client.Kill()
		}
		r.client = nil
		r.clientProtocol = nil
	}
	<-done
}
