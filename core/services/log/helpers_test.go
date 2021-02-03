package log_test

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type LogNewRound struct {
	types.Log
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
}

type simpleLogListener struct {
	handler    func(lb log.Broadcast, err error)
	consumerID *models.ID
	recvd      []log.Broadcast
	recvdMu    sync.Mutex
}

func (listener simpleLogListener) receivedBroadcasts() []log.Broadcast {
	listener.recvdMu.Lock()
	defer listener.recvdMu.Unlock()
	cp := make([]log.Broadcast, len(listener.recvd))
	copy(cp, listener.recvd)
	return cp
}

func (listener simpleLogListener) receivedLogs() []types.Log {
	broadcasts := listener.receivedBroadcasts()
	var logs []types.Log
	for _, b := range broadcasts {
		logs = append(logs, b.RawLog())
	}
	return logs
}

func (listener simpleLogListener) HandleLog(lb log.Broadcast, err error) {
	listener.recvdMu.Lock()
	defer listener.recvdMu.Unlock()
	listener.recvd = append(listener.recvd, lb)
	listener.handler(lb, err)
}
func (listener simpleLogListener) OnConnect()    {}
func (listener simpleLogListener) OnDisconnect() {}
func (listener simpleLogListener) JobID() *models.ID {
	return listener.consumerID
}
func (listener simpleLogListener) IsV2Job() bool {
	return false
}
func (listener simpleLogListener) JobIDV2() int32 {
	return 0
}

type logBroadcastRow struct {
	BlockHash   common.Hash
	BlockNumber uint64
	LogIndex    uint
	JobID       *models.ID
	JobIDV2     int32
	Consumed    bool
}

type mockListener struct {
	jobID   *models.ID
	jobIDV2 int32
}

func (l *mockListener) JobID() *models.ID              { return l.jobID }
func (l *mockListener) JobIDV2() int32                 { return l.jobIDV2 }
func (l *mockListener) IsV2Job() bool                  { return l.jobID == nil }
func (l *mockListener) OnConnect()                     {}
func (l *mockListener) OnDisconnect()                  {}
func (l *mockListener) HandleLog(log.Broadcast, error) {}

func createJob(t *testing.T, store *store.Store) models.JobSpec {
	t.Helper()

	job := cltest.NewJob()
	err := store.ORM.CreateJob(&job)
	require.NoError(t, err)
	return job
}
