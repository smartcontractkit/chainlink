package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	medianconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	ErrNoSuchKeyBundle      = errors.New("no such key bundle exists")
	ErrNoSuchTransmitterKey = errors.New("no such transmitter key exists")
	ErrNoSuchPublicKey      = errors.New("no such public key exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	InsertWebhookSpec(webhookSpec *WebhookSpec, qopts ...pg.QOpt) error
	InsertJob(job *Job, qopts ...pg.QOpt) error
	CreateJob(jb *Job, qopts ...pg.QOpt) error
	FindJobs(offset, limit int) ([]Job, int, error)
	FindJobTx(id int32) (Job, error)
	FindJob(ctx context.Context, id int32) (Job, error)
	FindJobByExternalJobID(uuid uuid.UUID, qopts ...pg.QOpt) (Job, error)
	FindJobIDByAddress(address ethkey.EIP55Address, qopts ...pg.QOpt) (int32, error)
	FindJobIDsWithBridge(name string) ([]int32, error)
	DeleteJob(id int32, qopts ...pg.QOpt) error
	RecordError(jobID int32, description string, qopts ...pg.QOpt) error
	// TryRecordError is a helper which calls RecordError and logs the returned error if present.
	TryRecordError(jobID int32, description string, qopts ...pg.QOpt)
	DismissError(ctx context.Context, errorID int64) error
	FindSpecError(id int64, qopts ...pg.QOpt) (SpecError, error)
	Close() error
	PipelineRuns(jobID *int32, offset, size int) ([]pipeline.Run, int, error)

	FindPipelineRunIDsByJobID(jobID int32, offset, limit int) (ids []int64, err error)
	FindPipelineRunsByIDs(ids []int64) (runs []pipeline.Run, err error)
	CountPipelineRunsByJobID(jobID int32) (count int32, err error)

	FindJobsByPipelineSpecIDs(ids []int32) ([]Job, error)
	FindPipelineRunByID(id int64) (pipeline.Run, error)
}

type orm struct {
	q           pg.Q
	chainSet    evm.ChainSet
	keyStore    keystore.Master
	pipelineORM pipeline.ORM
	lggr        logger.Logger
	bridgeORM   bridges.ORM
}

var _ ORM = (*orm)(nil)

func NewORM(
	db *sqlx.DB,
	chainSet evm.ChainSet,
	pipelineORM pipeline.ORM,
	keyStore keystore.Master, // needed to validation key properties on new job creation
	lggr logger.Logger,
	cfg pg.LogConfig,
) *orm {
	namedLogger := lggr.Named("JobORM")
	return &orm{
		q:           pg.NewQ(db, namedLogger, cfg),
		chainSet:    chainSet,
		keyStore:    keyStore,
		pipelineORM: pipelineORM,
		bridgeORM:   bridges.NewORM(db, lggr, cfg),
		lggr:        namedLogger,
	}
}
func (o *orm) Close() error {
	return nil
}

func (o *orm) assertBridgesExist(p pipeline.Pipeline) error {
	var bridgeNames = make(map[bridges.BridgeName]struct{})
	var uniqueBridges []bridges.BridgeName
	for _, task := range p.Tasks {
		if task.Type() == pipeline.TaskTypeBridge {
			// Bridge must exist
			name := task.(*pipeline.BridgeTask).Name
			bridge, err := bridges.ParseBridgeName(name)
			if err != nil {
				return err
			}
			if _, have := bridgeNames[bridge]; have {
				continue
			}
			bridgeNames[bridge] = struct{}{}
			uniqueBridges = append(uniqueBridges, bridge)
		}
	}
	if len(uniqueBridges) != 0 {
		_, err := o.bridgeORM.FindBridges(uniqueBridges)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateJob creates the job, and it's associated spec record.
// Expects an unmarshalled job spec as the jb argument i.e. output from ValidatedXX.
// Scans all persisted records back into jb
func (o *orm) CreateJob(jb *Job, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	p := jb.Pipeline
	if err := o.assertBridgesExist(p); err != nil {
		return err
	}

	var jobID int32
	err := q.Transaction(func(tx pg.Queryer) error {
		// Autogenerate a job ID if not specified
		if jb.ExternalJobID == (uuid.UUID{}) {
			jb.ExternalJobID = uuid.NewV4()
		}

		switch jb.Type {
		case DirectRequest:
			var specID int32
			sql := `INSERT INTO direct_request_specs (contract_address, min_incoming_confirmations, requesters, min_contract_payment, evm_chain_id, created_at, updated_at)
			VALUES (:contract_address, :min_incoming_confirmations, :requesters, :min_contract_payment, :evm_chain_id, now(), now())
			RETURNING id;`
			if err := pg.PrepareQueryRowx(tx, sql, &specID, jb.DirectRequestSpec); err != nil {
				return errors.Wrap(err, "failed to create DirectRequestSpec")
			}
			jb.DirectRequestSpecID = &specID
		case FluxMonitor:
			var specID int32
			sql := `INSERT INTO flux_monitor_specs (contract_address, threshold, absolute_threshold, poll_timer_period, poll_timer_disabled, idle_timer_period, idle_timer_disabled,
					drumbeat_schedule, drumbeat_random_delay, drumbeat_enabled, min_payment, evm_chain_id, created_at, updated_at)
			VALUES (:contract_address, :threshold, :absolute_threshold, :poll_timer_period, :poll_timer_disabled, :idle_timer_period, :idle_timer_disabled,
					:drumbeat_schedule, :drumbeat_random_delay, :drumbeat_enabled, :min_payment, :evm_chain_id, NOW(), NOW())
			RETURNING id;`
			if err := pg.PrepareQueryRowx(tx, sql, &specID, jb.FluxMonitorSpec); err != nil {
				return errors.Wrap(err, "failed to create FluxMonitorSpec")
			}
			jb.FluxMonitorSpecID = &specID
		case OffchainReporting:
			var specID int32
			if jb.OCROracleSpec.EncryptedOCRKeyBundleID != nil {
				_, err := o.keyStore.OCR().Get(jb.OCROracleSpec.EncryptedOCRKeyBundleID.String())
				if err != nil {
					return errors.Wrapf(ErrNoSuchKeyBundle, "%v", jb.OCROracleSpec.EncryptedOCRKeyBundleID)
				}
			}
			if jb.OCROracleSpec.TransmitterAddress != nil {
				_, err := o.keyStore.Eth().Get(jb.OCROracleSpec.TransmitterAddress.Hex())
				if err != nil {
					return errors.Wrapf(ErrNoSuchTransmitterKey, "%v", jb.OCROracleSpec.TransmitterAddress)
				}
			}

			existingSpec := new(OCROracleSpec)
			err := tx.Get(existingSpec, `SELECT * FROM ocr_oracle_specs WHERE contract_address = $1 and (evm_chain_id = $2 or evm_chain_id IS NULL) LIMIT 1;`,
				jb.OCROracleSpec.ContractAddress, jb.OCROracleSpec.EVMChainID,
			)
			if !errors.Is(err, sql.ErrNoRows) {
				if err != nil {
					return errors.Wrap(err, "failed to validate OffchainreportingOracleSpec on creation")
				}

				matchErr := errors.Errorf("a job with contract address %s already exists for chain ID %d", jb.OCROracleSpec.ContractAddress, jb.OCROracleSpec.EVMChainID.ToInt())
				if existingSpec.EVMChainID == nil {
					chain, err2 := o.chainSet.Default()
					if err2 != nil {
						return errors.Wrap(err2, "failed to validate OffchainreportingOracleSpec on creation")
					}
					if jb.OCROracleSpec.EVMChainID.Equal((*utils.Big)(chain.ID())) {
						return matchErr
					}
				} else {
					return matchErr
				}
			}

			sql := `INSERT INTO ocr_oracle_specs (contract_address, p2p_bootstrap_peers, is_bootstrap_peer, encrypted_ocr_key_bundle_id, transmitter_address,
					observation_timeout, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, evm_chain_id,
					created_at, updated_at, database_timeout, observation_grace_period, contract_transmitter_transmit_timeout)
			VALUES (:contract_address, :p2p_bootstrap_peers, :is_bootstrap_peer, :encrypted_ocr_key_bundle_id, :transmitter_address,
					:observation_timeout, :blockchain_timeout, :contract_config_tracker_subscribe_interval, :contract_config_tracker_poll_interval, :contract_config_confirmations, :evm_chain_id,
					NOW(), NOW(), :database_timeout, :observation_grace_period, :contract_transmitter_transmit_timeout)
			RETURNING id;`
			err = pg.PrepareQueryRowx(tx, sql, &specID, jb.OCROracleSpec)
			if err != nil {
				return errors.Wrap(err, "failed to create OffchainreportingOracleSpec")
			}
			jb.OCROracleSpecID = &specID
		case OffchainReporting2:
			var specID int32
			if jb.OCR2OracleSpec.OCRKeyBundleID.Valid {
				_, err := o.keyStore.OCR2().Get(jb.OCR2OracleSpec.OCRKeyBundleID.String)
				if err != nil {
					return errors.Wrapf(ErrNoSuchKeyBundle, "%v", jb.OCR2OracleSpec.OCRKeyBundleID)
				}
			}
			if jb.OCR2OracleSpec.TransmitterID.Valid {
				switch jb.OCR2OracleSpec.Relay {
				case relaytypes.EVM:
					_, err := o.keyStore.Eth().Get(jb.OCR2OracleSpec.TransmitterID.String)
					if err != nil {
						return errors.Wrapf(ErrNoSuchTransmitterKey, "%v", jb.OCR2OracleSpec.TransmitterID)
					}
				case relaytypes.Solana:
					_, err := o.keyStore.Solana().Get(jb.OCR2OracleSpec.TransmitterID.String)
					if err != nil {
						return errors.Wrapf(ErrNoSuchTransmitterKey, "%v", jb.OCR2OracleSpec.TransmitterID)
					}
				case relaytypes.Terra:
					_, err := o.keyStore.Terra().Get(jb.OCR2OracleSpec.TransmitterID.String)
					if err != nil {
						return errors.Wrapf(ErrNoSuchTransmitterKey, "%v", jb.OCR2OracleSpec.TransmitterID)
					}
				}
			}
			switch jb.OCR2OracleSpec.PluginType {
			case Median:
				var cfg medianconfig.PluginConfig
				err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &cfg)
				if err != nil {
					return errors.Wrap(err, "failed to parse plugin config")
				}
				feePipeline, err := pipeline.Parse(cfg.JuelsPerFeeCoinPipeline)
				if err != nil {
					return err
				}
				if err2 := o.assertBridgesExist(*feePipeline); err2 != nil {
					return err2
				}
			}

			sql := `INSERT INTO ocr2_oracle_specs (contract_id, relay, relay_config, plugin_type, plugin_config, p2p_bootstrap_peers, ocr_key_bundle_id, transmitter_id,
					blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations,
					created_at, updated_at)
			VALUES (:contract_id, :relay, :relay_config, :plugin_type, :plugin_config, :p2p_bootstrap_peers, :ocr_key_bundle_id, :transmitter_id,
					 :blockchain_timeout, :contract_config_tracker_poll_interval, :contract_config_confirmations,
					NOW(), NOW())
			RETURNING id;`
			err := pg.PrepareQueryRowx(tx, sql, &specID, jb.OCR2OracleSpec)
			if err != nil {
				return errors.Wrap(err, "failed to create Offchainreporting2OracleSpec")
			}
			jb.OCR2OracleSpecID = &specID
		case Keeper:
			var specID int32
			sql := `INSERT INTO keeper_specs (contract_address, from_address, evm_chain_id, created_at, updated_at)
			VALUES (:contract_address, :from_address, :evm_chain_id, NOW(), NOW())
			RETURNING id;`
			if err := pg.PrepareQueryRowx(tx, sql, &specID, jb.KeeperSpec); err != nil {
				return errors.Wrap(err, "failed to create KeeperSpec")
			}
			jb.KeeperSpecID = &specID
		case Cron:
			var specID int32
			sql := `INSERT INTO cron_specs (cron_schedule, created_at, updated_at)
			VALUES (:cron_schedule, NOW(), NOW())
			RETURNING id;`
			if err := pg.PrepareQueryRowx(tx, sql, &specID, jb.CronSpec); err != nil {
				return errors.Wrap(err, "failed to create CronSpec")
			}
			jb.CronSpecID = &specID
		case VRF:
			var specID int32
			sql := `INSERT INTO vrf_specs (
				coordinator_address, public_key, min_incoming_confirmations, 
				evm_chain_id, from_addresses, poll_period, requested_confs_delay, 
				request_timeout, chunk_size, batch_coordinator_address, batch_fulfillment_enabled, 
				batch_fulfillment_gas_multiplier, backoff_initial_delay, backoff_max_delay,
				created_at, updated_at)
			VALUES (
				:coordinator_address, :public_key, :min_incoming_confirmations, 
				:evm_chain_id, :from_addresses, :poll_period, :requested_confs_delay, 
				:request_timeout, :chunk_size, :batch_coordinator_address, :batch_fulfillment_enabled,
				:batch_fulfillment_gas_multiplier, :backoff_initial_delay, :backoff_max_delay,
				NOW(), NOW())
			RETURNING id;`

			err := pg.PrepareQueryRowx(tx, sql, &specID, toVRFSpecRow(jb.VRFSpec))
			var pqErr *pgconn.PgError
			ok := errors.As(err, &pqErr)
			if err != nil && ok && pqErr.Code == "23503" {
				if pqErr.ConstraintName == "vrf_specs_public_key_fkey" {
					return errors.Wrapf(ErrNoSuchPublicKey, "%s", jb.VRFSpec.PublicKey.String())
				}
			}
			if err != nil {
				return errors.Wrap(err, "failed to create VRFSpec")
			}
			jb.VRFSpecID = &specID
		case Webhook:
			err := o.InsertWebhookSpec(jb.WebhookSpec, pg.WithQueryer(tx))
			if err != nil {
				return errors.Wrap(err, "failed to create WebhookSpec")
			}
			jb.WebhookSpecID = &jb.WebhookSpec.ID

			if len(jb.WebhookSpec.ExternalInitiatorWebhookSpecs) > 0 {
				for i := range jb.WebhookSpec.ExternalInitiatorWebhookSpecs {
					jb.WebhookSpec.ExternalInitiatorWebhookSpecs[i].WebhookSpecID = jb.WebhookSpec.ID
				}
				sql := `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec)
			VALUES (:external_initiator_id, :webhook_spec_id, :spec);`
				query, args, err := tx.BindNamed(sql, jb.WebhookSpec.ExternalInitiatorWebhookSpecs)
				if err != nil {
					return errors.Wrap(err, "failed to bindquery for ExternalInitiatorWebhookSpecs")
				}
				if _, err = tx.Exec(query, args...); err != nil {
					return errors.Wrap(err, "failed to create ExternalInitiatorWebhookSpecs")
				}
			}
		case BlockhashStore:
			var specID int32
			sql := `INSERT INTO blockhash_store_specs (coordinator_v1_address, coordinator_v2_address, wait_blocks, lookback_blocks, blockhash_store_address, poll_period, run_timeout, evm_chain_id, from_address, created_at, updated_at)
			VALUES (:coordinator_v1_address, :coordinator_v2_address, :wait_blocks, :lookback_blocks, :blockhash_store_address, :poll_period, :run_timeout, :evm_chain_id, :from_address, NOW(), NOW())
			RETURNING id;`
			if err := pg.PrepareQueryRowx(tx, sql, &specID, jb.BlockhashStoreSpec); err != nil {
				return errors.Wrap(err, "failed to create BlockhashStore spec")
			}
			jb.BlockhashStoreSpecID = &specID
		case Bootstrap:
			var specID int32
			sql := `INSERT INTO bootstrap_specs (contract_id, relay, relay_config, monitoring_endpoint,
					blockchain_timeout, contract_config_tracker_poll_interval, 
					contract_config_confirmations, created_at, updated_at)
			VALUES (:contract_id, :relay, :relay_config, :monitoring_endpoint, 
					:blockchain_timeout, :contract_config_tracker_poll_interval, 
					:contract_config_confirmations, NOW(), NOW())
			RETURNING id;`
			if err := pg.PrepareQueryRowx(tx, sql, &specID, jb.BootstrapSpec); err != nil {
				return errors.Wrap(err, "failed to create BootstrapSpec for jobSpec")
			}
			jb.BootstrapSpecID = &specID
		default:
			o.lggr.Panicf("Unsupported jb.Type: %v", jb.Type)
		}

		pipelineSpecID, err := o.pipelineORM.CreateSpec(p, jb.MaxTaskDuration, pg.WithQueryer(tx))
		if err != nil {
			return errors.Wrap(err, "failed to create pipeline spec")
		}

		jb.PipelineSpecID = pipelineSpecID
		err = o.InsertJob(jb, pg.WithQueryer(tx))
		jobID = jb.ID
		return errors.Wrap(err, "failed to insert job")
	})
	if err != nil {
		return errors.Wrap(err, "CreateJobFailed")
	}

	return o.findJob(jb, "id", jobID, qopts...)
}

func (o *orm) InsertWebhookSpec(webhookSpec *WebhookSpec, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := `INSERT INTO webhook_specs (created_at, updated_at)
			VALUES (NOW(), NOW())
			RETURNING *;`
	return q.GetNamed(query, webhookSpec, webhookSpec)
}

func (o *orm) InsertJob(job *Job, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := `INSERT INTO jobs (pipeline_spec_id, name, schema_version, type, max_task_duration, ocr_oracle_spec_id, ocr2_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id,
				keeper_spec_id, cron_spec_id, vrf_spec_id, webhook_spec_id, blockhash_store_spec_id, bootstrap_spec_id, external_job_id, created_at)
		VALUES (:pipeline_spec_id, :name, :schema_version, :type, :max_task_duration, :ocr_oracle_spec_id, :ocr2_oracle_spec_id, :direct_request_spec_id, :flux_monitor_spec_id,
				:keeper_spec_id, :cron_spec_id, :vrf_spec_id, :webhook_spec_id, :blockhash_store_spec_id, :bootstrap_spec_id, :external_job_id, NOW())
		RETURNING *;`
	return q.GetNamed(query, job, job)
}

// DeleteJob removes a job
func (o *orm) DeleteJob(id int32, qopts ...pg.QOpt) error {
	o.lggr.Debugw("Deleting job", "jobID", id)
	// Added a 1 minute timeout to this query since this can take a long time as data increases.
	// This was added specifically due to an issue with a database that had a millions of pipeline_runs and pipeline_task_runs
	// and this query was taking ~40secs.
	qopts = append(qopts, pg.WithLongQueryTimeout())
	q := o.q.WithOpts(qopts...)
	query := `
		WITH deleted_jobs AS (
			DELETE FROM jobs WHERE id = $1 RETURNING
				pipeline_spec_id,
				ocr_oracle_spec_id,
				ocr2_oracle_spec_id,
				keeper_spec_id,
				cron_spec_id,
				flux_monitor_spec_id,
				vrf_spec_id,
				webhook_spec_id,
				direct_request_spec_id,
				blockhash_store_spec_id,
				bootstrap_spec_id
		),
		deleted_oracle_specs AS (
			DELETE FROM ocr_oracle_specs WHERE id IN (SELECT ocr_oracle_spec_id FROM deleted_jobs)
		),
		deleted_oracle2_specs AS (
			DELETE FROM ocr2_oracle_specs WHERE id IN (SELECT ocr2_oracle_spec_id FROM deleted_jobs)
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
		),
		deleted_blockhash_store_specs AS (
			DELETE FROM blockhash_store_specs WHERE id IN (SELECT blockhash_store_spec_id FROM deleted_jobs)
		),
		deleted_bootstrap_specs AS (
			DELETE FROM bootstrap_specs WHERE id IN (SELECT bootstrap_spec_id FROM deleted_jobs)
		)
		DELETE FROM pipeline_specs WHERE id IN (SELECT pipeline_spec_id FROM deleted_jobs)`
	res, cancel, err := q.ExecQIter(query, id)
	defer cancel()
	if err != nil {
		return errors.Wrap(err, "DeleteJob failed to delete job")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "DeleteJob failed getting RowsAffected")
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	o.lggr.Debugw("Deleted job", "jobID", id)
	return nil
}

func (o *orm) RecordError(jobID int32, description string, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	sql := `INSERT INTO job_spec_errors (job_id, description, occurrences, created_at, updated_at)
	VALUES ($1, $2, 1, $3, $3)
	ON CONFLICT (job_id, description) DO UPDATE SET
	occurrences = job_spec_errors.occurrences + 1,
	updated_at = excluded.updated_at`
	err := q.ExecQ(sql, jobID, description, time.Now())
	// Noop if the job has been deleted.
	var pqErr *pgconn.PgError
	ok := errors.As(err, &pqErr)
	if err != nil && ok && pqErr.Code == "23503" {
		if pqErr.ConstraintName == "job_spec_errors_v2_job_id_fkey" {
			return nil
		}
	}
	return err
}
func (o *orm) TryRecordError(jobID int32, description string, qopts ...pg.QOpt) {
	err := o.RecordError(jobID, description, qopts...)
	o.lggr.ErrorIf(err, fmt.Sprintf("Error creating SpecError %v", description))
}

func (o *orm) DismissError(ctx context.Context, ID int64) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	res, cancel, err := q.ExecQIter("DELETE FROM job_spec_errors WHERE id = $1", ID)
	defer cancel()
	if err != nil {
		return errors.Wrap(err, "failed to dismiss error")
	}
	n, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to dismiss error")
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (o *orm) FindSpecError(id int64, qopts ...pg.QOpt) (SpecError, error) {
	stmt := `SELECT * FROM job_spec_errors WHERE id = $1;`

	specErr := new(SpecError)
	err := o.q.WithOpts(qopts...).Get(specErr, stmt, id)

	return *specErr, errors.Wrap(err, "FindSpecError failed")
}

func (o *orm) FindJobs(offset, limit int) (jobs []Job, count int, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		sql := `SELECT count(*) FROM jobs;`
		err = tx.QueryRowx(sql).Scan(&count)
		if err != nil {
			return err
		}

		sql = `SELECT * FROM jobs ORDER BY created_at DESC, id DESC OFFSET $1 LIMIT $2;`
		err = tx.Select(&jobs, sql, offset, limit)
		if err != nil {
			return err
		}

		err = LoadAllJobsTypes(tx, jobs)
		if err != nil {
			return err
		}
		for i := range jobs {
			err = o.LoadEnvConfigVars(&jobs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	return jobs, int(count), err
}

func (o *orm) LoadEnvConfigVars(jb *Job) error {
	if jb.OCROracleSpec != nil {
		ch, err := o.chainSet.Get(jb.OCROracleSpec.EVMChainID.ToInt())
		if err != nil {
			return err
		}
		newSpec, err := LoadEnvConfigVarsOCR(ch.Config(), o.keyStore.P2P(), *jb.OCROracleSpec)
		if err != nil {
			return err
		}
		jb.OCROracleSpec = newSpec
	} else if jb.VRFSpec != nil {
		ch, err := o.chainSet.Get(jb.VRFSpec.EVMChainID.ToInt())
		if err != nil {
			return err
		}
		jb.VRFSpec = LoadEnvConfigVarsVRF(ch.Config(), *jb.VRFSpec)
	} else if jb.DirectRequestSpec != nil {
		ch, err := o.chainSet.Get(jb.DirectRequestSpec.EVMChainID.ToInt())
		if err != nil {
			return err
		}
		jb.DirectRequestSpec = LoadEnvConfigVarsDR(ch.Config(), *jb.DirectRequestSpec)
	}
	return nil
}

type DRSpecConfig interface {
	MinIncomingConfirmations() uint32
}

func LoadEnvConfigVarsVRF(cfg DRSpecConfig, vrfs VRFSpec) *VRFSpec {
	// Take the larger of the global vs specific.
	// Note that the v2 vrf requests specify their own confirmation requirements.
	// We wait for max(minIncomingConfirmations, request required confs) to be safe.
	minIncomingConfirmations := cfg.MinIncomingConfirmations()
	if vrfs.MinIncomingConfirmations <= minIncomingConfirmations {
		vrfs.ConfirmationsEnv = true
		vrfs.MinIncomingConfirmations = minIncomingConfirmations
	}

	if vrfs.PollPeriod == 0 {
		vrfs.PollPeriodEnv = true
		vrfs.PollPeriod = 5 * time.Second
	}

	return &vrfs
}

func LoadEnvConfigVarsDR(cfg DRSpecConfig, drs DirectRequestSpec) *DirectRequestSpec {
	minIncomingConfirmations := cfg.MinIncomingConfirmations()
	if drs.MinIncomingConfirmations.Uint32 > minIncomingConfirmations {
		drs.MinIncomingConfirmationsEnv = true
		drs.MinIncomingConfirmations = null.Uint32From(minIncomingConfirmations)
	}

	return &drs
}

type OCRSpecConfig interface {
	P2PPeerID() p2pkey.PeerID
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRObservationTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRObservationGracePeriod() time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRTransmitterAddress() (ethkey.EIP55Address, error)
	OCRKeyBundleID() (string, error)
}

// LoadEnvConfigVarsLocalOCR loads local OCR env vars into the OCROracleSpec.
func LoadEnvConfigVarsLocalOCR(cfg OCRSpecConfig, os OCROracleSpec) *OCROracleSpec {
	if os.ObservationTimeout == 0 {
		os.ObservationTimeoutEnv = true
		os.ObservationTimeout = models.Interval(cfg.OCRObservationTimeout())
	}
	if os.BlockchainTimeout == 0 {
		os.BlockchainTimeoutEnv = true
		os.BlockchainTimeout = models.Interval(cfg.OCRBlockchainTimeout())
	}
	if os.ContractConfigTrackerSubscribeInterval == 0 {
		os.ContractConfigTrackerSubscribeIntervalEnv = true
		os.ContractConfigTrackerSubscribeInterval = models.Interval(cfg.OCRContractSubscribeInterval())
	}
	if os.ContractConfigTrackerPollInterval == 0 {
		os.ContractConfigTrackerPollIntervalEnv = true
		os.ContractConfigTrackerPollInterval = models.Interval(cfg.OCRContractPollInterval())
	}
	if os.ContractConfigConfirmations == 0 {
		os.ContractConfigConfirmationsEnv = true
		os.ContractConfigConfirmations = cfg.OCRContractConfirmations()
	}
	if os.DatabaseTimeout == nil {
		os.DatabaseTimeoutEnv = true
		os.DatabaseTimeout = models.NewInterval(cfg.OCRDatabaseTimeout())
	}
	if os.ObservationGracePeriod == nil {
		os.ObservationGracePeriodEnv = true
		os.ObservationGracePeriod = models.NewInterval(cfg.OCRObservationGracePeriod())
	}
	if os.ContractTransmitterTransmitTimeout == nil {
		os.ContractTransmitterTransmitTimeoutEnv = true
		os.ContractTransmitterTransmitTimeout = models.NewInterval(cfg.OCRContractTransmitterTransmitTimeout())
	}
	return &os
}

// LoadEnvConfigVarsOCR loads OCR env vars into the OCROracleSpec.
func LoadEnvConfigVarsOCR(cfg OCRSpecConfig, p2pStore keystore.P2P, os OCROracleSpec) (*OCROracleSpec, error) {
	if os.TransmitterAddress == nil {
		ta, err := cfg.OCRTransmitterAddress()
		if !errors.Is(errors.Cause(err), config.ErrUnset) {
			if err != nil {
				return nil, err
			}
			os.TransmitterAddressEnv = true
			os.TransmitterAddress = &ta
		}
	}

	if os.EncryptedOCRKeyBundleID == nil {
		kb, err := cfg.OCRKeyBundleID()
		if err != nil {
			return nil, err
		}
		encryptedOCRKeyBundleID, err := models.Sha256HashFromHex(kb)
		if err != nil {
			return nil, err
		}
		os.EncryptedOCRKeyBundleIDEnv = true
		os.EncryptedOCRKeyBundleID = &encryptedOCRKeyBundleID
	}

	return LoadEnvConfigVarsLocalOCR(cfg, os), nil
}

func (o *orm) FindJobTx(id int32) (Job, error) {
	ctx, cancel := pg.DefaultQueryCtx()
	defer cancel()
	return o.FindJob(ctx, id)
}

// FindJob returns job by ID, with all relations preloaded
func (o *orm) FindJob(ctx context.Context, id int32) (jb Job, err error) {
	err = o.findJob(&jb, "id", id, pg.WithParentCtx(ctx))
	return
}

func (o *orm) FindJobByExternalJobID(externalJobID uuid.UUID, qopts ...pg.QOpt) (jb Job, err error) {
	err = o.findJob(&jb, "external_job_id", externalJobID, qopts...)
	return
}

// FindJobIDByAddress - finds a job id by contract address. Currently only OCR and FM jobs are supported
func (o *orm) FindJobIDByAddress(address ethkey.EIP55Address, qopts ...pg.QOpt) (jobID int32, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT jobs.id
FROM jobs
LEFT JOIN ocr_oracle_specs ocrspec on ocrspec.contract_address = $1 AND ocrspec.id = jobs.ocr_oracle_spec_id
LEFT JOIN flux_monitor_specs fmspec on fmspec.contract_address = $1 AND fmspec.id = jobs.flux_monitor_spec_id
WHERE ocrspec.id IS NOT NULL OR fmspec.id IS NOT NULL
`
		err = tx.Get(&jobID, stmt, address)

		if !errors.Is(err, sql.ErrNoRows) {
			if err != nil {
				return errors.Wrap(err, "error searching for job by contract address")
			}
			return nil
		}

		return err
	})

	return jobID, errors.Wrap(err, "FindJobIDByAddress failed")
}

func (o *orm) findJob(jb *Job, col string, arg interface{}, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.Transaction(func(tx pg.Queryer) error {
		sql := fmt.Sprintf(`SELECT * FROM jobs WHERE %s = $1 LIMIT 1`, col)
		err := tx.Get(jb, sql, arg)
		if err != nil {
			return errors.Wrap(err, "failed to load job")
		}

		if err = LoadAllJobTypes(tx, jb); err != nil {
			return err
		}

		return loadJobSpecErrors(tx, jb)
	})
	if err != nil {
		return errors.Wrap(err, "findJob failed")
	}
	return o.LoadEnvConfigVars(jb)
}

func (o *orm) FindJobIDsWithBridge(name string) (jids []int32, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		query := `SELECT jobs.id, dot_dag_source FROM jobs JOIN pipeline_specs ON pipeline_specs.id = jobs.pipeline_spec_id WHERE dot_dag_source ILIKE '%' || $1 || '%' ORDER BY id`
		var rows *sqlx.Rows
		rows, err = tx.Queryx(query, name)
		if err != nil {
			return err
		}
		defer rows.Close()
		var ids []int32
		var sources []string
		for rows.Next() {
			var id int32
			var source string
			if err = rows.Scan(&id, &source); err != nil {
				return err
			}
			ids = append(jids, id)
			sources = append(sources, source)
		}

		for i, id := range ids {
			var p *pipeline.Pipeline
			p, err = pipeline.Parse(sources[i])
			if err != nil {
				return errors.Wrapf(err, "could not parse dag for job %d", id)
			}
			for _, task := range p.Tasks {
				if task.Type() == pipeline.TaskTypeBridge {
					if task.(*pipeline.BridgeTask).Name == name {
						jids = append(jids, id)
					}
				}
			}
		}
		return nil
	})
	return jids, errors.Wrap(err, "FindJobIDsWithBridge failed")
}

// PipelineRunsByJobsIDs returns pipeline runs for multiple jobs, not preloading data
func (o *orm) PipelineRunsByJobsIDs(ids []int32) (runs []pipeline.Run, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		stmt := `SELECT pipeline_runs.* FROM pipeline_runs INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id WHERE jobs.id = ANY($1)
		ORDER BY pipeline_runs.created_at DESC, pipeline_runs.id DESC;`
		if err = tx.Select(&runs, stmt, ids); err != nil {
			return errors.Wrap(err, "error loading runs")
		}

		runs, err = o.loadPipelineRunsRelations(runs, tx)

		return err
	})

	return runs, errors.Wrap(err, "PipelineRunsByJobsIDs failed")
}

// FindPipelineRunIDsByJobID fetches the ids of pipeline runs for a job.
func (o *orm) FindPipelineRunIDsByJobID(jobID int32, offset, limit int) (ids []int64, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT pipeline_runs.id
FROM pipeline_runs
WHERE pipeline_runs.pipeline_spec_id = (SELECT jobs.pipeline_spec_id FROM JOBS WHERE jobs.id = $1)
ORDER BY pipeline_runs.created_at DESC, pipeline_runs.id DESC
OFFSET $2
LIMIT $3
`
		if err = tx.Select(&ids, stmt, jobID, offset, limit); err != nil {
			return errors.Wrap(err, "error loading runs")
		}

		return err
	})

	return ids, errors.Wrap(err, "PipelineRunsByJobIDs failed")
}

// FindPipelineRunsByIDs returns pipeline runs with the ids.
func (o *orm) FindPipelineRunsByIDs(ids []int64) (runs []pipeline.Run, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT pipeline_runs.*
FROM pipeline_runs
WHERE id = ANY($1)
`

		if err = tx.Select(&runs, stmt, ids); err != nil {
			return errors.Wrap(err, "error loading runs")
		}

		runs, err = o.loadPipelineRunsRelations(runs, tx)

		return err
	})

	return runs, errors.Wrap(err, "FindPipelineRunsByIDs failed")
}

// FindPipelineRunByID returns pipeline run with the id.
func (o *orm) FindPipelineRunByID(id int64) (pipeline.Run, error) {
	var run pipeline.Run

	err := o.q.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT pipeline_runs.*
FROM pipeline_runs
WHERE id = $1
`

		if err := tx.Get(&run, stmt, id); err != nil {
			return errors.Wrap(err, "error loading run")
		}

		runs, err := o.loadPipelineRunsRelations([]pipeline.Run{run}, tx)

		run = runs[0]

		return err
	})

	return run, errors.Wrap(err, "FindPipelineRunByID failed")
}

// CountPipelineRunsByJobID returns the total number of pipeline runs for a job.
func (o *orm) CountPipelineRunsByJobID(jobID int32) (count int32, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT COUNT(*)
FROM pipeline_runs
WHERE pipeline_runs.pipeline_spec_id = (SELECT jobs.pipeline_spec_id FROM JOBS WHERE jobs.id = $1)
`
		if err = tx.Get(&count, stmt, jobID); err != nil {
			return errors.Wrap(err, "error counting runs")
		}

		return err
	})

	return count, errors.Wrap(err, "PipelineRunsByJobsIDs failed")
}

func (o *orm) FindJobsByPipelineSpecIDs(ids []int32) ([]Job, error) {
	var jbs []Job

	err := o.q.Transaction(func(tx pg.Queryer) error {
		stmt := `SELECT * FROM jobs WHERE jobs.pipeline_spec_id = ANY($1) ORDER BY id ASC
`
		if err := tx.Select(&jbs, stmt, ids); err != nil {
			return errors.Wrap(err, "error fetching jobs by pipeline spec IDs")
		}

		err := LoadAllJobsTypes(tx, jbs)
		if err != nil {
			return err
		}
		for i := range jbs {
			err = o.LoadEnvConfigVars(&jbs[i])
			if err != nil {
				return err
			}
		}

		return nil
	})

	return jbs, errors.Wrap(err, "FindJobsByPipelineSpecIDs failed")
}

// PipelineRuns returns pipeline runs for a job, with spec and taskruns loaded, latest first
// If jobID is nil, returns all pipeline runs
func (o *orm) PipelineRuns(jobID *int32, offset, size int) (runs []pipeline.Run, count int, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		var args []interface{}
		var filter string
		if jobID != nil {
			filter = "JOIN jobs USING(pipeline_spec_id) WHERE jobs.id = $1" // TODO:  add support for more than 1 jobID?
			args = append(args, *jobID)
		}
		sql := fmt.Sprintf(`SELECT count(*) FROM pipeline_runs %s`, filter)
		if err = tx.QueryRowx(sql, args...).Scan(&count); err != nil {
			return errors.Wrap(err, "error counting runs")
		}

		sql = fmt.Sprintf(`SELECT pipeline_runs.* FROM pipeline_runs %s
		ORDER BY pipeline_runs.created_at DESC, pipeline_runs.id DESC
		OFFSET $%d LIMIT $%d;`, filter, len(args)+1, len(args)+2)

		if err = tx.Select(&runs, sql, append(args, offset, size)...); err != nil {
			return errors.Wrap(err, "error loading runs")
		}

		runs, err = o.loadPipelineRunsRelations(runs, tx)

		return err
	})

	return runs, count, errors.Wrap(err, "PipelineRuns failed")
}

func (o *orm) loadPipelineRunsRelations(runs []pipeline.Run, tx pg.Queryer) ([]pipeline.Run, error) {
	// Postload PipelineSpecs
	// TODO: We should pull this out into a generic preload function once go has generics
	specM := make(map[int32]pipeline.Spec)
	for _, run := range runs {
		if _, exists := specM[run.PipelineSpecID]; !exists {
			specM[run.PipelineSpecID] = pipeline.Spec{}
		}
	}
	specIDs := make([]int32, len(specM))
	for specID := range specM {
		specIDs = append(specIDs, specID)
	}
	stmt := `SELECT pipeline_specs.*, jobs.id AS job_id FROM pipeline_specs JOIN jobs ON pipeline_specs.id = jobs.pipeline_spec_id WHERE pipeline_specs.id = ANY($1);`
	var specs []pipeline.Spec
	if err := o.q.Select(&specs, stmt, specIDs); err != nil {
		return nil, errors.Wrap(err, "error loading specs")
	}
	for _, spec := range specs {
		specM[spec.ID] = spec
	}
	runM := make(map[int64]*pipeline.Run, len(runs))
	for i, run := range runs {
		runs[i].PipelineSpec = specM[run.PipelineSpecID]
		runM[run.ID] = &runs[i]
	}

	// Postload PipelineTaskRuns
	runIDs := make([]int64, len(runs))
	for i, run := range runs {
		runIDs[i] = run.ID
	}
	var taskRuns []pipeline.TaskRun
	stmt = `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = ANY($1) ORDER BY pipeline_run_id, created_at, id;`
	if err := tx.Select(&taskRuns, stmt, runIDs); err != nil {
		return nil, errors.Wrap(err, "error loading pipeline_task_runs")
	}
	for _, taskRun := range taskRuns {
		run := runM[taskRun.PipelineRunID]
		run.PipelineTaskRuns = append(run.PipelineTaskRuns, taskRun)
	}

	return runs, nil
}

// NOTE: N+1 query, be careful of performance
// This is not easily fixable without complicating the logic a lot, since we
// only use it in the GUI it's probably acceptable
func LoadAllJobsTypes(tx pg.Queryer, jobs []Job) error {
	for i := range jobs {
		err := LoadAllJobTypes(tx, &jobs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadAllJobTypes(tx pg.Queryer, job *Job) error {
	return multierr.Combine(
		loadJobType(tx, job, "PipelineSpec", "pipeline_specs", &job.PipelineSpecID),
		loadJobType(tx, job, "FluxMonitorSpec", "flux_monitor_specs", job.FluxMonitorSpecID),
		loadJobType(tx, job, "DirectRequestSpec", "direct_request_specs", job.DirectRequestSpecID),
		loadJobType(tx, job, "OCROracleSpec", "ocr_oracle_specs", job.OCROracleSpecID),
		loadJobType(tx, job, "OCR2OracleSpec", "ocr2_oracle_specs", job.OCR2OracleSpecID),
		loadJobType(tx, job, "KeeperSpec", "keeper_specs", job.KeeperSpecID),
		loadJobType(tx, job, "CronSpec", "cron_specs", job.CronSpecID),
		loadJobType(tx, job, "WebhookSpec", "webhook_specs", job.WebhookSpecID),
		loadVRFJob(tx, job, job.VRFSpecID),
		loadJobType(tx, job, "BlockhashStoreSpec", "blockhash_store_specs", job.BlockhashStoreSpecID),
		loadJobType(tx, job, "BootstrapSpec", "bootstrap_specs", job.BootstrapSpecID),
	)
}

func loadJobType(tx pg.Queryer, job *Job, field, table string, id *int32) error {
	if id == nil {
		return nil
	}

	// The abomination below allows us to initialise and then scan into the
	// type of the field without hardcoding for each individual field
	// My LIFE for generics...
	r := reflect.ValueOf(job)
	t := reflect.Indirect(r).FieldByName(field).Type().Elem()
	destVal := reflect.New(t)
	dest := destVal.Interface()

	err := tx.Get(dest, fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, table), *id)

	if err != nil {
		return errors.Wrapf(err, "failed to load job type %s with id %d", table, *id)
	}
	reflect.ValueOf(job).Elem().FieldByName(field).Set(destVal)
	return nil
}

func loadVRFJob(tx pg.Queryer, job *Job, id *int32) error {
	if id == nil {
		return nil
	}

	var row vrfSpecRow
	err := tx.Get(&row, `SELECT * FROM vrf_specs WHERE id = $1`, *id)
	if err != nil {
		return errors.Wrapf(err, `failed to load job type VRFSpec with id %d`, *id)
	}

	job.VRFSpec = row.toVRFSpec()
	return nil
}

// vrfSpecRow is a helper type for reading and writing VRF specs to the database. This is necessary
// because the bytea[] in the DB is not automatically convertible to or from the spec's
// FromAddresses field. pq.ByteaArray must be used instead.
type vrfSpecRow struct {
	*VRFSpec
	FromAddresses pq.ByteaArray
}

func toVRFSpecRow(spec *VRFSpec) vrfSpecRow {
	addresses := make(pq.ByteaArray, len(spec.FromAddresses))
	for i, a := range spec.FromAddresses {
		addresses[i] = a.Bytes()
	}
	return vrfSpecRow{VRFSpec: spec, FromAddresses: addresses}
}

func (r vrfSpecRow) toVRFSpec() *VRFSpec {
	for _, a := range r.FromAddresses {
		r.VRFSpec.FromAddresses = append(r.VRFSpec.FromAddresses,
			ethkey.EIP55AddressFromAddress(common.BytesToAddress(a)))
	}
	return r.VRFSpec
}

func loadJobSpecErrors(tx pg.Queryer, jb *Job) error {
	return errors.Wrapf(tx.Select(&jb.JobSpecErrors, `SELECT * FROM job_spec_errors WHERE job_id = $1`, jb.ID), "failed to load job spec errors for job %d", jb.ID)
}
