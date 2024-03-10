package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (lsn *Listener) SetReqAdded(fn func()) {
	lsn.ReqAdded = fn
}

func (lsn *Listener) Stop(t *testing.T) {
	assert.NoError(t, lsn.Close())
	select {
	case <-lsn.WaitOnStop:
	case <-time.After(time.Second):
		t.Error("did not clean up properly")
	}
}

func (lsn *Listener) ReqsConfirmedAt() (us []uint64) {
	for i := range lsn.Reqs {
		us = append(us, lsn.Reqs[i].confirmedAtBlock)
	}
	return us
}

func (lsn *Listener) RespCount(reqIDBytes [32]byte) uint64 {
	return lsn.ResponseCount[reqIDBytes]
}

func (lsn *Listener) SetRespCount(reqIDBytes [32]byte, c uint64) {
	lsn.ResponseCount[reqIDBytes] = c
}
