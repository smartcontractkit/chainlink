package pipeline

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	// Spawner
	UnclaimedJobs() ([]JobSpec, error)
	CreatePipelineSpec(spec PipelineSpec) error
	LoadPipelineSpec(id int64) (PipelineSpec, error)

	// Runner
	CreatePipelineRun(prun PipelineRun) (int64, error)
	WithNextUnclaimedTaskRun(f func(tx *gorm.DB, ptRun PipelineTaskRun, predecessors []PipelineTaskRun) error) error
	MarkTaskRunCompleted(tx *gorm.DB, taskRunID int64, output sql.Scanner, err error) error
	NotifyCompletion(pipelineRunID int64) error
}

type orm struct {
	db database
}

type database interface {
	Exec(sql string, values ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Transaction(fn func(db *gorm.DB) error) error
}

func NewORM(o database) *orm {
	return &orm{o}
}

func (o *orm) CreatePipelineSpec(spec *PipelineSpec) error {
	return o.db.Create(spec)
}

func (o *orm) LoadPipelineSpec(id int64) (spec PipelineSpec, err error) {
	err = o.db.Preload("TaskSpecs").First(&spec, id).Error
	return j, err
}

// CreatePipelineRun adds a PipelineRun record to the DB, and one PipelineTaskRun
// per PipelineTaskSpec associated with the given PipelineSpec.  Processing of the
// TaskRuns is maximally parallelized across all of the Chainlink nodes in the
// cluster.
func (o *orm) CreatePipelineRun(specID int64) (int64, error) {
	return o.db.Transaction(func(tx *gorm.DB) error {
		prun := PipelineRun{PipelineSpecID: specID}

		err := tx.Create(&prun).Error
		if err != nil {
			return errors.Wrap(err, "could not create pipeline run")
		}

		err = tx.Exec(`
            INSERT INTO pipeline_task_runs (
                pipeline_run_id, pipeline_task_spec_id, created_at
            )
            SELECT ? AS pipeline_run_id, id AS pipeline_task_spec_id, NOW() AS created_at
            FROM pipeline_task_specs
            WHERE pipeline_spec_id = ?
        `, prun.ID, prun.PipelineSpecID).Error
		return errors.Wrap(err, "could not create pipeline task runs")
	})
}

// WithNextUnclaimedTaskRun chooses any arbitrary incomplete TaskRun from the DB
// whose parent TaskRuns have already been processed.
func (o *orm) WithNextUnclaimedTaskRun(fn func(ptRun PipelineTaskRun, predecessors []PipelineTaskRun) Result) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		var ptRun PipelineTaskRun
		var predecessors []PipelineTaskRun

		// NOTE: This could conceivably be done in pure SQL and made marginally more efficient by preloading with joins

		// Find the next unlocked, unfinished pipeline_task_run with no uncompleted predecessors
		err := tx.Table("pipeline_task_runs AS successor").
			Set("gorm:query_option", "FOR UPDATE OF successor SKIP LOCKED").
			Joins("LEFT JOIN pipeline_task_runs AS predecessors ON successor.id = predecessor.successor_id AND predecessors.finished_at IS NULL").
			Where("predecessors.id IS NULL").
			Where("successor.finished_at IS NULL").
			Preload("PipelineTaskSpec").
			Order("id ASC").
			First(&ptRun).
			Error
		if err != nil {
			return errors.Wrap(err, "error finding next task run")
		}

		// Find all the predecessors
		err = tx.Where("successor_id = ?", ptRun.ID).Find(&predecessors).Error
		if err != nil {
			return errors.Wrap(err, "error finding task run predecessors")
		}

		// Call the callback and convert its output to a format appropriate for the DB
		result := fn(ptRun, predecessors)
		var out *JSONSerializable
		var errString null.String
		if result.Value != nil {
			out = &JSONSerializable{Value: result.Value}
		} else if result.Error != nil {
			errString = null.StringFrom(err.Error())
		}

		// Update the task run record with the output and error
		err = tx.Exec(`
            UPDATE pipeline_task_runs
            SET output = ?, error = ?, finished_at = ?
            WHERE id = ?
        `, out, errString, time.Now(), ptRun.ID).Error

		return errors.Wrap(err, "could not mark pipeline_task_run as finished")
	})
}

func (o *orm) NotifyCompletion(pipelineRunID int64) error {
	return o.db.Exec(`
    $$
    BEGIN
        IF (SELECT bool_and(pipeline_task_runs.error IS NOT NULL OR pipeline_task_runs.output IS NOT NULL) FROM pipeline_job_runs JOIN pipeline_task_runs ON pipeline_task_runs.pipeline_job_run_id = pipeline_job_runs.id WHERE pipeline_job_runs.id = $1)
            PERFORM pg_notify('pipeline_job_run_completed', $1::text);
        END IF;
    END;
    $$ LANGUAGE plpgsql;
    )`, pipelineRunID).Error
}
