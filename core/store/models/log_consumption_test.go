package models_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
)

func TestMarkLogConsumed_Happy(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job1 := cltest.NewJob()
	err := store.ORM.CreateJob(&job1)
	require.NoError(t, err)
	job2 := cltest.NewJob()
	err = store.ORM.CreateJob(&job2)
	require.NoError(t, err)

	blockHash1 := cltest.NewHash()
	logIndex1 := uint(0)
	jobID1 := job1.ID

	err = store.ORM.MarkLogConsumed(blockHash1, logIndex1, jobID1)
	require.NoError(t, err)

	tests := []struct {
		description string
		BlockHash   common.Hash
		LogIndex    uint
		JobID       *models.ID
	}{
		{"different blockhash", cltest.NewHash(), logIndex1, jobID1},
		{"different log", blockHash1, 1, jobID1},
		{"different consumer", blockHash1, logIndex1, job2.ID},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err = store.ORM.MarkLogConsumed(test.BlockHash, test.LogIndex, test.JobID)
			require.NoError(t, err)
		})
	}

}

func TestMarkLogConsumed_Errors(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job1 := cltest.NewJob()
	err := store.ORM.CreateJob(&job1)
	require.NoError(t, err)
	job2 := cltest.NewJob()
	err = store.ORM.CreateJob(&job2)
	require.NoError(t, err)

	blockHash1 := cltest.NewHash()
	logIndex1 := uint(0)
	jobID1 := job1.ID

	err = store.ORM.MarkLogConsumed(blockHash1, logIndex1, jobID1)
	require.NoError(t, err)

	tests := []struct {
		description string
		BlockHash   common.Hash
		LogIndex    uint
		JobID       *models.ID
	}{
		{"non existent job", cltest.NewHash(), 0, models.NewID()},
		{"duplicate record", blockHash1, logIndex1, jobID1},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err = store.ORM.MarkLogConsumed(test.BlockHash, test.LogIndex, test.JobID)
			require.Error(t, err)
		})
	}
}
