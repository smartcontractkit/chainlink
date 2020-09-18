package job

import (
	"github.com/jinzhu/gorm"

	"github.com/smartcontractkit/chainlink/core/store/models"
	ormpkg "github.com/smartcontractkit/chainlink/core/store/orm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	UnclaimedJobs(jobTypeRegistrations map[Type]Registration, claimedJobs map[int32][]Service) ([]Spec, error)
	CreateJob(spec Spec) error
	DeleteJob(spec Spec) error
}

type orm struct {
	db *gorm.DB
}

var _ ORM = (*orm)(nil)

func NewORM(o *gorm.DB) *orm {
	return &orm{o}
}

func (o *orm) UnclaimedJobs(jobTypeRegistrations map[Type]Registration, claimedJobs map[int32][]Service) ([]Spec, error) {
	var specs []Spec
	ormpkg.GormTransaction(o.db, func(tx *gorm.DB) error {
		// Loop over each job type and fetch any unclaimed specs
		for _, reg := range jobTypeRegistrations {
			var concreteJobs []models.JobSpecV2

			// Make a slice of this job type's specs.  Gorm will take care of the rest.
			// specType := reflect.TypeOf(reg.Spec)
			// specSliceType := reflect.SliceOf(specType)
			// specSlice := reflect.Zero(specSliceType)

			err := o.db.
				Preload("OffchainreportingOracleSpec").
				Scan(&concreteJobs).Error
			if err != nil {
				return err
			}

			for _, j := range concreteJobs {
				// Skip claimed jobs
				if _, exists := claimedJobs[j.ID]; exists {
					continue
				}

				// Resolve DB polymorphism
				var spec Spec
				switch {
				case j.OffchainreportingOracleSpec != nil:
					spec = j.OffchainreportingOracleSpec
				default:
					continue
				}

				specs = append(specs, spec)
			}
		}
		return nil
	})
	return specs, nil
}

func (o *orm) CreateJob(spec Spec) error {
	return ormpkg.GormTransaction(o.db, func(tx *gorm.DB) error {
		err := tx.Create(spec).Error
		if err != nil {
			return err
		}

		pipelineSpec, err := spec.TaskDAG().ToPipelineSpec()
		if err != nil {
			return err
		}
		pipelineSpec.JobSpecID = spec.JobID()

		return tx.Create(pipelineSpec).Error
	})
}

func (o *orm) DeleteJob(spec Spec) error {
	return o.db.Delete(spec).Error
}
