package vrf

import (
	"testing"
	"time"
)

type ListenerV1 = listenerV1

func (l *listenerV1) SetReqAdded(fn func()) {
	l.reqAdded = fn
}

func (l *listenerV1) RunLogListener(unsubs []func(), minConfs uint32) {
	l.runLogListener(unsubs, minConfs)
}

func (l *listenerV1) RunHeadListener(fn func()) {
	l.runHeadListener(fn)
}

func (l *listenerV1) Stop(t *testing.T) {
	l.chStop <- struct{}{}
	select {
	case <-l.waitOnStop:
	case <-time.After(time.Second):
		t.Error("did not clean up properly")
	}
}

func (l *listenerV1) ReqsConfirmedAt() (us []uint64) {
	for i := range l.reqs {
		us = append(us, l.reqs[i].confirmedAtBlock)
	}
	return us
}

func (l *listenerV1) RespCount(reqIDBytes [32]byte) uint64 {
	return l.respCount[reqIDBytes]
}

func (l *listenerV1) SetRespCount(reqIDBytes [32]byte, c uint64) {
	l.respCount[reqIDBytes] = c
}
