package v1

import (
	"testing"
	"time"
)

type ListenerV1 = Listener

func (l *Listener) SetReqAdded(fn func()) {
	l.ReqAdded = fn
}

func (l *Listener) RunLogListener(unsubs []func(), minConfs uint32) {
	l.runLogListener(unsubs, minConfs)
}

func (l *Listener) RunHeadListener(fn func()) {
	l.runHeadListener(fn)
}

func (l *Listener) Stop(t *testing.T) {
	l.ChStop <- struct{}{}
	select {
	case <-l.WaitOnStop:
	case <-time.After(time.Second):
		t.Error("did not clean up properly")
	}
}

func (l *Listener) ReqsConfirmedAt() (us []uint64) {
	for i := range l.Reqs {
		us = append(us, l.Reqs[i].confirmedAtBlock)
	}
	return us
}

func (l *Listener) RespCount(reqIDBytes [32]byte) uint64 {
	return l.RespCount[reqIDBytes]
}

func (l *Listener) SetRespCount(reqIDBytes [32]byte, c uint64) {
	l.RespCount[reqIDBytes] = c
}
