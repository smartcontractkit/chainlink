package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/sqlx"

	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNoSuchPeerID             = errors.New("no such peer id exists")
	ErrNoSuchKeyBundle          = errors.New("no such key bundle exists")
	ErrNoSuchTransmitterAddress = errors.New("no such transmitter address exists")
	ErrNoSuchPublicKey          = errors.New("no such public key exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

var (
	ErrViolatesForeignKeyConstraint = errors.New("violates foreign key constraint")
)

type ORM interface {
	CreateJob(ctx context.Context, jobSpec *Job, pipeline pipeline.Pipeline) (Job, error)
	JobsV2(offset, limit int) ([]Job, int, error)
	FindJobTx(id int32) (Job, error)
	FindJob(ctx context.Context, id int32) (Job, error)
	FindJobByExternalJobID(ctx context.Context, uuid uuid.UUID) (Job, error)
	FindJobIDsWithBridge(name string) ([]int32, error)
	DeleteJob(ctx context.Context, id int32) error
	RecordError(ctx context.Context, jobID int32, description string)
	DismissError(ctx context.Context, errorID int32) error
	Close() error
	PipelineRuns(offset, size int) ([]pipeline.Run, int, error)
	PipelineRunsByJobID(jobID int32, offset, size int) ([]pipeline.Run, int, error)
}

type orm struct {
	db          *gorm.DB
	chainSet    evm.ChainSet
	keyStore    keystore.Master
	pipelineORM pipeline.ORM
}

var _ ORM = (*orm)(nil)

func NewORM(
	db *gorm.DB,
	chainSet evm.ChainSet,
	pipelineORM pipeline.ORM,
	keyStore keystore.Master, // needed to validation key properties on new job creation
) *orm {
	return &orm{
		db:          db,
		chainSet:    chainSet,
		keyStore:    keyStore,
		pipelineORM: pipelineORM,
	}
}

func PreloadAllJobTypes(db *gorm.DB) *gorm.DB {
	return db.
		Preload("PipelineSpec").
		Preload("FluxMonitorSpec").
		Preload("DirectRequestSpec").
		Preload("OffchainreportingOracleSpec").
		Preload("KeeperSpec").
		Preload("PipelineSpec").
		Preload("CronSpec").
		Preload("WebhookSpec").
		Preload("VRFSpec")
}

func (o *orm) Close() error {
	return nil
}

// CreateJob creates the job and it's associated spec record.
//
// NOTE: This is not wrapped in a db transaction so if you call this, you should
// use postgres.TransactionManager to create the transaction in the context.
// Expects an unmarshaled job spec as the jobSpec argument i.e. output from ValidatedXX.
// Returns a fully populated Job.
func (o *orm) CreateJob(ctx context.Context, jobSpec *Job, p pipeline.Pipeline) (Job, error) {
	var jb Job
	for _, task := range p.Tasks {
		if task.Type() == pipeline.TaskTypeBridge {
			// Bridge must exist
			name := task.(*pipeline.BridgeTask).Name
			bt := bridges.BridgeType{}
			if err := o.db.First(&bt, "name = ?", name).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return jb, errors.Wrap(pipeline.ErrNoSuchBridge, name)
				}
				return jb, err
			}
		}
	}

	tx := postgres.TxFromContext(ctx, o.db)

	// Autogenerate a job ID if not specified
	if jobSpec.ExternalJobID == (uuid.UUID{}) {
		jobSpec.ExternalJobID = uuid.NewV4()
	}

	switch jobSpec.Type {
	case DirectRequest:
		err := tx.Create(&jobSpec.DirectRequestSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create DirectRequestSpec for jobSpec")
		}
		jobSpec.DirectRequestSpecID = &jobSpec.DirectRequestSpec.ID
	case FluxMonitor:
		err := tx.Create(&jobSpec.FluxMonitorSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create FluxMonitorSpec for jobSpec")
		}
		jobSpec.FluxMonitorSpecID = &jobSpec.FluxMonitorSpec.ID
	case OffchainReporting:
		if jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID.Valid {
			_, err := o.keyStore.OCR().Get(jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID.String)
			if err != nil {
				return jb, errors.Wrapf(ErrNoSuchKeyBundle, "%v", jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
			}
		}
		if jobSpec.OffchainreportingOracleSpec.P2PPeerID != nil {
			_, err := o.keyStore.P2P().Get(jobSpec.OffchainreportingOracleSpec.P2PPeerID.Raw())
			if err != nil {
				return jb, errors.Wrapf(ErrNoSuchPeerID, "%v", jobSpec.OffchainreportingOracleSpec.P2PPeerID)
			}
		}
		if jobSpec.OffchainreportingOracleSpec.TransmitterAddress != nil {
			_, err := o.keyStore.Eth().Get(jobSpec.OffchainreportingOracleSpec.TransmitterAddress.Hex())
			if err != nil {
				return jb, errors.Wrapf(ErrNoSuchTransmitterAddress, "%v", jobSpec.OffchainreportingOracleSpec.TransmitterAddress)
			}
		}

		err := tx.Create(&jobSpec.OffchainreportingOracleSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create OffchainreportingOracleSpec for jobSpec")
		}
		jobSpec.OffchainreportingOracleSpecID = &jobSpec.OffchainreportingOracleSpec.ID
	case Keeper:
		err := tx.Create(&jobSpec.KeeperSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create KeeperSpec for jobSpec")
		}
		jobSpec.KeeperSpecID = &jobSpec.KeeperSpec.ID
	case Cron:
		err := tx.Create(&jobSpec.CronSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create CronSpec for jobSpec")
		}
		jobSpec.CronSpecID = &jobSpec.CronSpec.ID
	case VRF:
		err := tx.Create(&jobSpec.VRFSpec).Error
		pqErr, ok := err.(*pgconn.PgError)
		if err != nil && ok && pqErr.Code == "23503" {
			if pqErr.ConstraintName == "vrf_specs_public_key_fkey" {
				return jb, errors.Wrapf(ErrNoSuchPublicKey, "%s", jobSpec.VRFSpec.PublicKey.String())
			}
		}
		if err != nil {
			return jb, errors.Wrap(err, "failed to create VRFSpec for jobSpec")
		}
		jobSpec.VRFSpecID = &jobSpec.VRFSpec.ID
	case Webhook:
		err := tx.Create(&jobSpec.WebhookSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create WebhookSpec for jobSpec")
		}
		jobSpec.WebhookSpecID = &jobSpec.WebhookSpec.ID
		for i, eiWS := range jobSpec.WebhookSpec.ExternalInitiatorWebhookSpecs {
			jobSpec.WebhookSpec.ExternalInitiatorWebhookSpecs[i].WebhookSpecID = jobSpec.WebhookSpec.ID
			err := tx.Create(&jobSpec.WebhookSpec.ExternalInitiatorWebhookSpecs[i]).Error
			if err != nil {
				return jb, errors.Wrapf(err, "failed to create ExternalInitiatorWebhookSpec for WebhookSpec: %#v", eiWS)
			}
		}
	default:
		logger.Fatalf("Unsupported jobSpec.Type: %v", jobSpec.Type)
	}

	pipelineSpecID, err := o.pipelineORM.CreateSpec(ctx, tx, p, jobSpec.MaxTaskDuration)
	if err != nil {
		return jb, errors.Wrap(err, "failed to create pipeline spec")
	}
	jobSpec.PipelineSpecID = pipelineSpecID
	err = tx.Create(jobSpec).Error
	if err != nil {
		return jb, errors.Wrap(err, "failed to create job")
	}

	return o.FindJob(ctx, jobSpec.ID)
}

// DeleteJob removes a job
func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	tx := postgres.TxFromContext(ctx, o.db)

	err := tx.Exec(`
		WITH deleted_jobs AS (
			DELETE FROM jobs WHERE id = ? RETURNING
				pipeline_spec_id,
				offchainreporting_oracle_spec_id,
				keeper_spec_id,
				cron_spec_id,
				flux_monitor_spec_id,
				vrf_spec_id,
				webhook_spec_id,
				direct_request_spec_id
		),
		deleted_oracle_specs AS (
			DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT offchainreporting_oracle_spec_id FROM deleted_jobs)
		),
		deleted_keeper_specs AS (
			DELETE FROM keeper_specs WHERE id IN (SELECT keeper_spec_id FROM deleted_jobs)
		),
		deleted_cron_specs AS (
			DELETE FROM cron_specs WHERE id IN (SELECT cron_spec_id FROM deleted_jobs)
		),
		deleted_fm_specs AS (
			DELETE FROM flux_monitor_specs WHERE id IN (SELECT flux_monitor_spec_id FROM deleted_jobs)
		),
		deleted_vrf_specs AS (
			DELETE FROM vrf_specs WHERE id IN (SELECT vrf_spec_id FROM deleted_jobs)
		),
		deleted_webhook_specs AS (
			DELETE FROM webhook_specs WHERE id IN (SELECT webhook_spec_id FROM deleted_jobs)
		),
		deleted_dr_specs AS (
			DELETE FROM direct_request_specs WHERE id IN (SELECT direct_request_spec_id FROM deleted_jobs)
		)
		DELETE FROM pipeline_specs WHERE id IN (SELECT pipeline_spec_id FROM deleted_jobs)
	`, id).Error
	if err != nil {
		return errors.Wrap(err, "DeleteJob failed to delete job")
	}

	return nil
}

func (o *orm) RecordError(ctx context.Context, jobID int32, description string) {
	pse := SpecError{JobID: jobID, Description: description, Occurrences: 1}
	err := o.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "job_id"}, {Name: "description"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"occurrences": gorm.Expr("job_spec_errors.occurrences + 1"),
				"updated_at":  gorm.Expr("excluded.updated_at"),
			}),
		}).
		Create(&pse).
		Error
	// Noop if the job has been deleted.
	if err != nil && strings.Contains(err.Error(), ErrViolatesForeignKeyConstraint.Error()) {
		return
	}
	logger.ErrorIf(err, fmt.Sprintf("error creating SpecError %v", description))
}

func (o *orm) DismissError(ctx context.Context, ID int32) error {
	result := o.db.Exec("DELETE FROM job_spec_errors WHERE id = ?", ID)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (o *orm) JobsV2(offset, limit int) ([]Job, int, error) {
	var count int64
	var jobs []Job
	err := postgres.GormTransactionWithDefaultContext(o.db, func(tx *gorm.DB) error {
		err := tx.
			Model(Job{}).
			Count(&count).
			Error

		if err != nil {
			return err
		}

		err = PreloadAllJobTypes(tx).
			Preload("JobSpecErrors").
			Limit(limit).
			Offset(offset).
			Order("id ASC").
			Find(&jobs).
			Error
		if err != nil {
			return err
		}
		for i := range jobs {
			if jobs[i].OffchainreportingOracleSpec != nil {
				ch, err := o.chainSet.Get(jobs[i].OffchainreportingOracleSpec.EVMChainID.ToInt())
				if err != nil {
					return err
				}
				jobs[i].OffchainreportingOracleSpec = LoadDynamicConfigVars(ch.Config(), *jobs[i].OffchainreportingOracleSpec)
			}
		}
		return nil
	})
	return jobs, int(count), err
}

type OCRSpecConfig interface {
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRObservationTimeout() time.Duration
}

func LoadDynamicConfigVars(cfg OCRSpecConfig, os OffchainReportingOracleSpec) *OffchainReportingOracleSpec {
	if os.ObservationTimeout == 0 {
		os.ObservationTimeout = models.Interval(cfg.OCRObservationTimeout())
	}
	if os.BlockchainTimeout == 0 {
		os.BlockchainTimeout = models.Interval(cfg.OCRBlockchainTimeout())
	}
	if os.ContractConfigTrackerSubscribeInterval == 0 {
		os.ContractConfigTrackerSubscribeInterval = models.Interval(cfg.OCRContractSubscribeInterval())
	}
	if os.ContractConfigTrackerPollInterval == 0 {
		os.ContractConfigTrackerPollInterval = models.Interval(cfg.OCRContractPollInterval())
	}
	if os.ContractConfigConfirmations == 0 {
		os.ContractConfigConfirmations = cfg.OCRContractConfirmations()
	}
	return &os
}

func (o *orm) FindJobTx(id int32) (Job, error) {
	var jb Job
	var err error
	txm := postgres.NewGormTransactionManager(o.db)
	err = txm.Transact(func(ctx context.Context) error {
		jb, err = o.FindJob(ctx, id)
		return err
	})
	return jb, err
}

// FindJob returns job by ID
func (o *orm) FindJob(ctx context.Context, id int32) (jb Job, err error) {
	tx := postgres.TxFromContext(ctx, o.db)
	err = PreloadAllJobTypes(tx).
		Preload("JobSpecErrors").
		First(&jb, "jobs.id = ?", id).
		Error
	if err != nil {
		return jb, err
	}

	if jb.OffchainreportingOracleSpec != nil {
		var ch evm.Chain
		ch, err = o.chainSet.Get(jb.OffchainreportingOracleSpec.EVMChainID.ToInt())
		if err != nil {
			return jb, err
		}
		jb.OffchainreportingOracleSpec = LoadDynamicConfigVars(ch.Config(), *jb.OffchainreportingOracleSpec)
	}
	return jb, err
}

func (o *orm) FindJobByExternalJobID(ctx context.Context, externalJobID uuid.UUID) (jb Job, err error) {
	tx := postgres.TxFromContext(ctx, o.db)
	err = PreloadAllJobTypes(tx).
		Preload("JobSpecErrors").
		First(&jb, "jobs.external_job_id = ?", externalJobID).
		Error
	if err != nil {
		return jb, err
	}

	if jb.OffchainreportingOracleSpec != nil {
		var ch evm.Chain
		ch, err = o.chainSet.Get(jb.OffchainreportingOracleSpec.EVMChainID.ToInt())
		if err != nil {
			return jb, err
		}
		jb.OffchainreportingOracleSpec = LoadDynamicConfigVars(ch.Config(), *jb.OffchainreportingOracleSpec)
	}
	return jb, err
}

func (o *orm) FindJobIDsWithBridge(name string) ([]int32, error) {
	var jobs []Job
	err := o.db.Preload("PipelineSpec").Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	var jids []int32
	for _, job := range jobs {
		p, err := pipeline.Parse(job.PipelineSpec.DotDagSource)
		if err != nil {
			return nil, err
		}
		for _, task := range p.Tasks {
			if task.Type() == pipeline.TaskTypeBridge {
				if task.(*pipeline.BridgeTask).Name == name {
					jids = append(jids, job.ID)
				}
			}
		}
	}
	return jids, nil
}

// Preload PipelineSpec.JobID for each Run
func (o *orm) preloadJobIDs(runs []pipeline.Run) error {
	db := postgres.UnwrapGormDB(o.db)

	ids := make([]int32, 0, len(runs))
	for _, run := range runs {
		ids = append(ids, run.PipelineSpecID)
	}

	// construct a WHERE IN query
	sql := `SELECT id, pipeline_spec_id FROM jobs WHERE pipeline_spec_id IN (?);`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	query = db.Rebind(query)
	var results []struct {
		ID             int32
		PipelineSpecID int32
	}
	if err := db.Select(&results, query, args...); err != nil {
		return err
	}

	// fill in fields
	for i := range runs {
		for _, result := range results {
			if result.PipelineSpecID == runs[i].PipelineSpecID {
				runs[i].PipelineSpec.JobID = result.ID
			}
		}
	}

	return nil
}

// PipelineRunsByJobID returns all pipeline runs
func (o *orm) PipelineRuns(offset, size int) ([]pipeline.Run, int, error) {
	var pipelineRuns []pipeline.Run
	var count int64
	err := o.db.
		Model(pipeline.Run{}).
		Count(&count).
		Error

	if err != nil {
		return pipelineRuns, 0, err
	}

	err = o.db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		}).
		Limit(size).
		Offset(offset).
		Order("created_at DESC, id DESC").
		Find(&pipelineRuns).
		Error

	if err != nil {
		return pipelineRuns, int(count), err
	}

	err = o.preloadJobIDs(pipelineRuns)

	return pipelineRuns, int(count), err
}

// PipelineRunsByJobID returns pipeline runs for a job
func (o *orm) PipelineRunsByJobID(jobID int32, offset, size int) ([]pipeline.Run, int, error) {
	var pipelineRuns []pipeline.Run
	var count int64
	err := o.db.
		Model(pipeline.Run{}).
		Joins("INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id").
		Where("jobs.id = ?", jobID).
		Count(&count).
		Error

	if err != nil {
		return pipelineRuns, 0, err
	}

	err = o.db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		}).
		Joins("INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id").
		Where("jobs.id = ?", jobID).
		Limit(size).
		Offset(offset).
		Order("created_at DESC, id DESC").
		Find(&pipelineRuns).
		Error

	// can skip preloadJobIDs since we already know the jobID
	for i := range pipelineRuns {
		pipelineRuns[i].PipelineSpec.JobID = jobID
	}

	return pipelineRuns, int(count), err
}
