package models_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
)

func TestCreateLogConsumption_Happy(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job1 := cltest.NewJob()
	err := store.ORM.CreateJob(&job1)
	require.NoError(t, err)
	job2 := cltest.NewJob()
	err = store.ORM.CreateJob(&job2)
	require.NoError(t, err)

	logConsumption1 := models.LogConsumption{
		BlockHash: cltest.NewHash(),
		LogIndex:  0,
		JobID:     job1.ID,
	}

	err = store.ORM.CreateLogConsumption(&logConsumption1)
	require.NoError(t, err)

	tests := []struct {
		description string
		BlockHash   common.Hash
		LogIndex    uint
		JobID       *models.ID
	}{
		{"different blockhash", cltest.NewHash(), logConsumption1.LogIndex, logConsumption1.JobID},
		{"different log", logConsumption1.BlockHash, 1, logConsumption1.JobID},
		{"different consumer", logConsumption1.BlockHash, logConsumption1.LogIndex, job2.ID},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			logConsumption2 := models.LogConsumption{
				BlockHash: test.BlockHash,
				LogIndex:  test.LogIndex,
				JobID:     test.JobID,
			}

			err = store.ORM.CreateLogConsumption(&logConsumption2)
			require.NoError(t, err)
		})
	}

}

func TestCreateLogConsumption_Errors(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job1 := cltest.NewJob()
	err := store.ORM.CreateJob(&job1)
	require.NoError(t, err)
	job2 := cltest.NewJob()
	err = store.ORM.CreateJob(&job2)
	require.NoError(t, err)

	logConsumption1 := models.LogConsumption{
		BlockHash: cltest.NewHash(),
		LogIndex:  0,
		JobID:     job1.ID,
	}

	err = store.ORM.CreateLogConsumption(&logConsumption1)
	require.NoError(t, err)

	tests := []struct {
		description string
		BlockHash   common.Hash
		LogIndex    uint
		JobID       *models.ID
	}{
		{"non existent job", cltest.NewHash(), 0, models.NewID()},
		{"duplicate record", logConsumption1.BlockHash, logConsumption1.LogIndex, logConsumption1.JobID},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			logConsumption2 := models.LogConsumption{
				BlockHash: test.BlockHash,
				LogIndex:  test.LogIndex,
				JobID:     test.JobID,
			}
			err = store.ORM.CreateLogConsumption(&logConsumption2)
			require.Error(t, err)
		})
	}
}
