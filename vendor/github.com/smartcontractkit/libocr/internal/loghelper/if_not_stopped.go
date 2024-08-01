package loghelper

import (
	"time"

	"github.com/smartcontractkit/libocr/subprocesses"
)

type IfNotStopped struct {
	chStop chan struct{}
	subs   subprocesses.Subprocesses
}

// If Stop is called prior to expiry of d, f won't be executed. Otherwise, f
// will be executed and Stop will block until f returns. That makes it different
// from the standard library's time.AfterFunc() whose Stop() function will
// return while f is still running.
func NewIfNotStopped(d time.Duration, f func()) *IfNotStopped {
	ins := IfNotStopped{
		make(chan struct{}, 1),
		subprocesses.Subprocesses{},
	}
	ins.subs.Go(func() {
		t := time.NewTimer(d)
		defer t.Stop()
		select {
		case <-t.C:
			f()
		case <-ins.chStop:
		}
	})
	return &ins
}

func (ins *IfNotStopped) Stop() {
	select {
	case <-ins.chStop:
		// chStop has been closed, don't close again
	default:
		close(ins.chStop)
	}

	ins.subs.Wait()
}
