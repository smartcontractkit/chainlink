package log_test

import (
	"math/big"
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
}

func (listener simpleLogListener) HandleLog(lb log.Broadcast, err error) {
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

func createJob(t *testing.T, store *store.Store) models.JobSpec {
	t.Helper()

	job := cltest.NewJob()
	err := store.ORM.CreateJob(&job)
	require.NoError(t, err)
	return job
}
