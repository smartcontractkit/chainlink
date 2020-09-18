package job_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestORM(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := job.NewORM(store.ORM.DB)

	jobTypeRegistrations := map[job.Type]job.Registration{
		models.OffchainReportingOracleSpec{}.JobType(): {Spec: &models.OffchainReportingOracleSpec{}},
	}
	services := map[int32][]job.Service{}
	specs := []*models.OffchainReportingOracleSpec{
		&models.OffchainReportingOracleSpec{JID: int32(1)},
		&models.OffchainReportingOracleSpec{JID: int32(2)},
		&models.OffchainReportingOracleSpec{JID: int32(3)},
	}

	t.Run("it creates job specs", func(t *testing.T) {
		for _, spec := range specs {
			err = orm.CreateJob(spec)
			require.NoError(t, err)
		}
	})

	t.Run("it correctly returns the unclaimed jobs in the DB", func(t *testing.T) {
		unclaimed, err := orm.UnclaimedJobs(jobTypeRegistrations, services)
		require.NoError(t, err)

		require.Len(t, unclaimed, len(specs))
		for _, spec := range specs {
			require.Contains(t, unclaimed, spec)
		}

		// Now simulate that a job has been claimed
		services[*specs[0].UUID] = []job.Service{}

		unclaimed, err = orm.UnclaimedJobs(jobTypeRegistrations, services)
		require.NoError(t, err)
		require.Len(t, unclaimed, len(specs)-1)
		for _, spec := range specs[1:] {
			require.Contains(t, unclaimed, spec)
		}
	})
}
