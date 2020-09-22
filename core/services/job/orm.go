package job

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	UnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error)
	CreateJob(jobSpec *models.JobSpecV2, pipelineSpec *pipeline.Spec) error
	DeleteJob(ctx context.Context, id int32) error
	Close() error
}

type orm struct {
	db           *gorm.DB
	advisoryLock *utils.PostgresAdvisoryLock
}

var _ ORM = (*orm)(nil)

func NewORM(o *gorm.DB, uri string) *orm {
	return &orm{o, &utils.PostgresAdvisoryLock{URI: uri}}
}

func (o *orm) Close() error {
	return o.advisoryLock.Close()
}

func (o *orm) UnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error) {
	var unclaimedJobs []models.JobSpecV2
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		var maybeJobs []models.JobSpecV2
		err := o.db.
			Preload("OffchainreportingOracleSpec").
			Preload("OffchainreportingOracleSpec.OffchainreportingKeyBundle").
			Find(&maybeJobs).Error
		if err != nil {
			return err
		}

		for _, job := range maybeJobs {
			err = o.advisoryLock.TryLock(ctx, utils.AdvisoryLockClassID_JobSpawner, job.ID)
			if err == nil {
				unclaimedJobs = append(unclaimedJobs, job)
			}
		}
		return nil
	})
	return unclaimedJobs, err
}

func (o *orm) CreateJob(jobSpec *models.JobSpecV2, pipelineSpec *pipeline.Spec) error {
	return utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		err := tx.Create(pipelineSpec).Error
		if err != nil {
			return err
		}
		jobSpec.PipelineSpecID = pipelineSpec.ID
		err = tx.Create(jobSpec.OffchainreportingOracleSpec.OffchainreportingKeyBundle).Error
		if err != nil {
			return err
		}
		jobSpec.OffchainreportingOracleSpec.OffchainreportingKeyBundleID = jobSpec.OffchainreportingOracleSpec.OffchainreportingKeyBundle.ID
		return tx.Create(jobSpec).Error
	})
}

func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	return utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		// If we can take the advisory lock, that means either we own this job or
		// nobody does.  That gives us permission to delete the job.  Note that we
		// have to unlock twice at the end (as we already have it).
		err := o.advisoryLock.TryLock(ctx, utils.AdvisoryLockClassID_JobSpawner, id)
		if err != nil {
			return err
		}
		defer o.advisoryLock.Unlock(ctx, utils.AdvisoryLockClassID_JobSpawner, id)
		defer o.advisoryLock.Unlock(ctx, utils.AdvisoryLockClassID_JobSpawner, id)
		return o.db.Where("id = ?", id).Delete(models.JobSpecV2{}).Error
	})
}
