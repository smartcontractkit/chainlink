package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	evmkeystore "github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	medianconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	ErrNoSuchKeyBundle      = errors.New("no such key bundle exists")
	ErrNoSuchTransmitterKey = errors.New("no such transmitter key exists")
	ErrNoSuchSendingKey     = errors.New("no such sending key exists")
	ErrNoSuchPublicKey      = errors.New("no such public key exists")
)

type ORM interface {
	InsertWebhookSpec(ctx context.Context, webhookSpec *WebhookSpec) error
	InsertJob(ctx context.Context, job *Job) error
	CreateJob(ctx context.Context, jb *Job) error
	FindJobs(ctx context.Context, offset, limit int) ([]Job, int, error)
	FindJob(ctx context.Context, id int32) (Job, error)
	FindJobByExternalJobID(ctx context.Context, uuid uuid.UUID) (Job, error)
	FindJobIDByAddress(ctx context.Context, address evmtypes.EIP55Address, evmChainID *big.Big) (int32, error)
	FindOCR2JobIDByAddress(ctx context.Context, contractID string, feedID *common.Hash) (int32, error)
	FindJobIDsWithBridge(ctx context.Context, name string) ([]int32, error)
	DeleteJob(ctx context.Context, id int32, jobType Type) error
	RecordError(ctx context.Context, jobID int32, description string) error
	// TryRecordError is a helper which calls RecordError and logs the returned error if present.
	TryRecordError(ctx context.Context, jobID int32, description string)
	DismissError(ctx context.Context, errorID int64) error
	FindSpecError(ctx context.Context, id int64) (SpecError, error)
	Close() error
	PipelineRuns(ctx context.Context, jobID *int32, offset, size int) ([]pipeline.Run, int, error)

	FindPipelineRunIDsByJobID(ctx context.Context, jobID int32, offset, limit int) (ids []int64, err error)
	FindPipelineRunsByIDs(ctx context.Context, ids []int64) (runs []pipeline.Run, err error)
	CountPipelineRunsByJobID(ctx context.Context, jobID int32) (count int32, err error)

	FindJobsByPipelineSpecIDs(ctx context.Context, ids []int32) ([]Job, error)
	FindPipelineRunByID(ctx context.Context, id int64) (pipeline.Run, error)

	FindSpecErrorsByJobIDs(ctx context.Context, ids []int32) ([]SpecError, error)
	FindJobWithoutSpecErrors(ctx context.Context, id int32) (jb Job, err error)

	FindTaskResultByRunIDAndTaskName(ctx context.Context, runID int64, taskName string) ([]byte, error)
	AssertBridgesExist(ctx context.Context, p pipeline.Pipeline) error

	DataSource() sqlutil.DataSource
	WithDataSource(source sqlutil.DataSource) ORM

	FindJobIDByWorkflow(ctx context.Context, spec WorkflowSpec) (int32, error)
	// TODO rename function to indicate it is CCIP-specific, not generic?
	FindJobIDByCapabilityNameAndVersion(ctx context.Context, spec CCIPSpec) (int32, error)
	FindStandardCapabilityJobID(ctx context.Context, spec StandardCapabilitiesSpec) (int32, error)
	FindGatewayJobID(ctx context.Context, spec GatewaySpec) (int32, error)

	FindJobIDByStreamID(ctx context.Context, streamID uint32) (int32, error)
}

type ORMConfig interface {
	DatabaseDefaultQueryTimeout() time.Duration
}

type orm struct {
	ds          sqlutil.DataSource
	keyStore    keystore.Master
	pipelineORM pipeline.ORM
	lggr        logger.SugaredLogger
	bridgeORM   bridges.ORM
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, pipelineORM pipeline.ORM, bridgeORM bridges.ORM, keyStore keystore.Master, lggr logger.Logger) *orm {
	namedLogger := logger.Sugared(lggr.Named("JobORM"))
	return &orm{
		ds:          ds,
		keyStore:    keyStore,
		pipelineORM: pipelineORM,
		bridgeORM:   bridgeORM,
		lggr:        namedLogger,
	}
}

func (o *orm) Close() error {
	return nil
}

func (o *orm) DataSource() sqlutil.DataSource {
	return o.ds
}

func (o *orm) WithDataSource(ds sqlutil.DataSource) ORM { return o.withDataSource(ds) }

func (o *orm) withDataSource(ds sqlutil.DataSource) *orm {
	n := &orm{
		ds:       ds,
		lggr:     o.lggr,
		keyStore: o.keyStore,
	}
	if o.bridgeORM != nil {
		n.bridgeORM = o.bridgeORM.WithDataSource(ds)
	}
	if o.pipelineORM != nil {
		n.pipelineORM = o.pipelineORM.WithDataSource(ds)
	}
	return n
}

func (o *orm) transact(ctx context.Context, readOnly bool, fn func(*orm) error) error {
	opts := &sqlutil.TxOptions{TxOptions: sql.TxOptions{ReadOnly: readOnly}}
	return sqlutil.Transact(ctx, o.withDataSource, o.ds, opts, fn)
}

func (o *orm) AssertBridgesExist(ctx context.Context, p pipeline.Pipeline) error {
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
		_, err := o.bridgeORM.FindBridges(ctx, uniqueBridges)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateJob creates the job, and it's associated spec record.
// Expects an unmarshalled job spec as the jb argument i.e. output from ValidatedXX.
// Scans all persisted records back into jb
func (o *orm) CreateJob(ctx context.Context, jb *Job) error {
	p := jb.Pipeline
	if err := o.AssertBridgesExist(ctx, p); err != nil {
		return err
	}

	var jobID int32
	err := o.transact(ctx, false, func(tx *orm) error {
		// Autogenerate a job ID if not specified
		if jb.ExternalJobID == (uuid.UUID{}) {
			jb.ExternalJobID = uuid.New()
		}

		switch jb.Type {
		case DirectRequest:
			if jb.DirectRequestSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertDirectRequestSpec(ctx, jb.DirectRequestSpec)
			if err != nil {
				return fmt.Errorf("failed to create DirectRequestSpec for jobSpec: %w", err)
			}
			jb.DirectRequestSpecID = &specID
		case FluxMonitor:
			if jb.FluxMonitorSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertFluxMonitorSpec(ctx, jb.FluxMonitorSpec)
			if err != nil {
				return fmt.Errorf("failed to create FluxMonitorSpec for jobSpec: %w", err)
			}
			jb.FluxMonitorSpecID = &specID
		case OffchainReporting:
			if jb.OCROracleSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}

			if jb.OCROracleSpec.EncryptedOCRKeyBundleID != nil {
				_, err := tx.keyStore.OCR().Get(jb.OCROracleSpec.EncryptedOCRKeyBundleID.String())
				if err != nil {
					return errors.Wrapf(ErrNoSuchKeyBundle, "no key bundle with id: %x", jb.OCROracleSpec.EncryptedOCRKeyBundleID)
				}
			}
			if jb.OCROracleSpec.TransmitterAddress != nil {
				_, err := tx.keyStore.Eth().Get(ctx, jb.OCROracleSpec.TransmitterAddress.Hex())
				if err != nil {
					return errors.Wrapf(ErrNoSuchTransmitterKey, "no key matching transmitter address: %s", jb.OCROracleSpec.TransmitterAddress.Hex())
				}
			}

			newChainID := jb.OCROracleSpec.EVMChainID
			existingSpec := new(OCROracleSpec)
			err := tx.ds.GetContext(ctx, existingSpec, `SELECT * FROM ocr_oracle_specs WHERE contract_address = $1 and (evm_chain_id = $2 or evm_chain_id IS NULL) LIMIT 1;`,
				jb.OCROracleSpec.ContractAddress, newChainID,
			)

			if !errors.Is(err, sql.ErrNoRows) {
				if err != nil {
					return errors.Wrap(err, "failed to validate OffchainreportingOracleSpec on creation")
				}

				return errors.Errorf("a job with contract address %s already exists for chain ID %s", jb.OCROracleSpec.ContractAddress, newChainID)
			}

			specID, err := tx.insertOCROracleSpec(ctx, jb.OCROracleSpec)
			if err != nil {
				return fmt.Errorf("failed to create OCROracleSpec for jobSpec: %w", err)
			}
			jb.OCROracleSpecID = &specID
		case OffchainReporting2:
			if jb.OCR2OracleSpec.OCRKeyBundleID.Valid {
				_, err := tx.keyStore.OCR2().Get(jb.OCR2OracleSpec.OCRKeyBundleID.String)
				if err != nil {
					return errors.Wrapf(ErrNoSuchKeyBundle, "no key bundle with id: %q", jb.OCR2OracleSpec.OCRKeyBundleID.ValueOrZero())
				}
			}

			if jb.OCR2OracleSpec.RelayConfig["sendingKeys"] != nil && jb.OCR2OracleSpec.TransmitterID.Valid {
				return errors.New("sending keys and transmitter ID can't both be defined")
			}

			// checks if they are present and if they are valid
			sendingKeysDefined, err := areSendingKeysDefined(ctx, jb, tx.keyStore)
			if err != nil {
				return err
			}

			if !sendingKeysDefined && !jb.OCR2OracleSpec.TransmitterID.Valid {
				return errors.New("neither sending keys nor transmitter ID is defined")
			}

			if !sendingKeysDefined {
				if err = ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, tx.keyStore, jb.OCR2OracleSpec.TransmitterID.String); err != nil {
					return errors.Wrap(ErrNoSuchTransmitterKey, err.Error())
				}
			}

			if jb.ForwardingAllowed && !slices.Contains(ForwardersSupportedPlugins, jb.OCR2OracleSpec.PluginType) {
				return errors.Errorf("forwarding is not currently supported for %s jobs", jb.OCR2OracleSpec.PluginType)
			}

			if jb.OCR2OracleSpec.PluginType == types.Mercury {
				if jb.OCR2OracleSpec.FeedID == nil {
					return errors.New("feed ID is required for mercury plugin type")
				}
			} else {
				if jb.OCR2OracleSpec.FeedID != nil {
					return errors.New("feed ID is not currently supported for non-mercury jobs")
				}
			}

			if jb.OCR2OracleSpec.PluginType == types.Median {
				var cfg medianconfig.PluginConfig

				validatePipeline := func(p string) error {
					pipeline, pipelineErr := pipeline.Parse(p)
					if pipelineErr != nil {
						return pipelineErr
					}
					return tx.AssertBridgesExist(ctx, *pipeline)
				}

				errUnmarshal := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &cfg)
				if errUnmarshal != nil {
					return errors.Wrap(errUnmarshal, "failed to parse plugin config")
				}

				if errFeePipeline := validatePipeline(cfg.JuelsPerFeeCoinPipeline); errFeePipeline != nil {
					return errFeePipeline
				}

				if cfg.HasGasPriceSubunitsPipeline() {
					if errGasPipeline := validatePipeline(cfg.GasPriceSubunitsPipeline); errGasPipeline != nil {
						return errGasPipeline
					}
				}
			}

			if enableDualTransmission, ok := jb.OCR2OracleSpec.RelayConfig["enableDualTransmission"]; ok && enableDualTransmission != nil {
				if jb.OCR2OracleSpec.Relay != relay.NetworkEVM {
					return errors.New("dual transmission is enabled only for EVM")
				}

				rawDualTransmissionConfig, ok := jb.OCR2OracleSpec.RelayConfig["dualTransmission"]
				if !ok {
					return errors.New("dual transmission is enabled but no dual transmission config present")
				}

				dualTransmissionConfig, ok := rawDualTransmissionConfig.(map[string]interface{})
				if !ok {
					return errors.New("invalid dual transmission config")
				}

				dtContractAddress, ok := dualTransmissionConfig["contractAddress"].(string)
				if !ok || !common.IsHexAddress(dtContractAddress) {
					return errors.New("invalid contract address in dual transmission config")
				}

				dtTransmitterAddress, ok := dualTransmissionConfig["transmitterAddress"].(string)
				if !ok || !common.IsHexAddress(dtTransmitterAddress) {
					return errors.New("invalid transmitter address in dual transmission config")
				}

				if _, ok := dualTransmissionConfig["meta"].(map[string]interface{}); !ok {
					return errors.New("invalid dual transmission meta")
				}

				if err = validateKeyStoreMatchForRelay(ctx, jb.OCR2OracleSpec.Relay, tx.keyStore, dtTransmitterAddress); err != nil {
					return errors.Wrap(err, "unknown dual transmission transmitterAddress")
				}

				// Check if secondary transmitter address is used as primary somewhere else
				if checkIfKeyHasLock(ctx, tx.keyStore.Eth(), common.HexToAddress(dtTransmitterAddress), evmkeystore.TXMv1) {
					return errors.Errorf("key %s cannot be a secondary transmitter address because it's used a primary transmitter in another job", dtTransmitterAddress)
				}
			}

			// Check if primary transmitter address is used as secondary somewhere else, don't check for mercury as it uses CSA keys for transmitters
			if jb.OCR2OracleSpec.PluginType != types.Mercury {
				if checkIfKeyHasLock(ctx, tx.keyStore.Eth(), common.HexToAddress(jb.OCR2OracleSpec.TransmitterID.String), evmkeystore.TXMv2) {
					return errors.Errorf("key %s cannot be a (primary) transmitter address because it's used a secondary transmitter address in another job", jb.OCR2OracleSpec.TransmitterID.String)
				}
			}

			specID, err := tx.insertOCR2OracleSpec(ctx, jb.OCR2OracleSpec)
			if err != nil {
				return fmt.Errorf("failed to create OCR2OracleSpec for jobSpec: %w", err)
			}
			jb.OCR2OracleSpecID = &specID
		case Keeper:
			if jb.KeeperSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertKeeperSpec(ctx, jb.KeeperSpec)
			if err != nil {
				return fmt.Errorf("failed to create KeeperSpec for jobSpec: %w", err)
			}
			jb.KeeperSpecID = &specID
		case Cron:
			specID, err := tx.insertCronSpec(ctx, jb.CronSpec)
			if err != nil {
				return fmt.Errorf("failed to create CronSpec for jobSpec: %w", err)
			}
			jb.CronSpecID = &specID
		case VRF:
			if jb.VRFSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertVRFSpec(ctx, jb.VRFSpec)
			var pqErr *pgconn.PgError
			ok := errors.As(err, &pqErr)
			if err != nil && ok && pqErr.Code == "23503" {
				if pqErr.ConstraintName == "vrf_specs_public_key_fkey" {
					return errors.Wrapf(ErrNoSuchPublicKey, "%s", jb.VRFSpec.PublicKey.String())
				}
			}
			if err != nil {
				return fmt.Errorf("failed to create VRFSpec for jobSpec: %w", err)
			}
			jb.VRFSpecID = &specID
		case Webhook:
			err := tx.InsertWebhookSpec(ctx, jb.WebhookSpec)
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
				if _, err := tx.ds.NamedExecContext(ctx, sql, jb.WebhookSpec.ExternalInitiatorWebhookSpecs); err != nil {
					return errors.Wrap(err, "failed to create ExternalInitiatorWebhookSpecs")
				}
			}
		case BlockhashStore:
			if jb.BlockhashStoreSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertBlockhashStoreSpec(ctx, jb.BlockhashStoreSpec)
			if err != nil {
				return fmt.Errorf("failed to create BlockhashStoreSpec for jobSpec: %w", err)
			}
			jb.BlockhashStoreSpecID = &specID
		case BlockHeaderFeeder:
			if jb.BlockHeaderFeederSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertBlockHeaderFeederSpec(ctx, jb.BlockHeaderFeederSpec)
			if err != nil {
				return fmt.Errorf("failed to create BlockHeaderFeederSpec for jobSpec: %w", err)
			}
			jb.BlockHeaderFeederSpecID = &specID
		case LegacyGasStationServer:
			if jb.LegacyGasStationServerSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertLegacyGasStationServerSpec(ctx, jb.LegacyGasStationServerSpec)
			if err != nil {
				return fmt.Errorf("failed to create LegacyGasStationServerSpec for jobSpec: %w", err)
			}
			jb.LegacyGasStationServerSpecID = &specID
		case LegacyGasStationSidecar:
			if jb.LegacyGasStationSidecarSpec.EVMChainID == nil {
				return errors.New("evm chain id must be defined")
			}
			specID, err := tx.insertLegacyGasStationSidecarSpec(ctx, jb.LegacyGasStationSidecarSpec)
			if err != nil {
				return fmt.Errorf("failed to create LegacyGasStationSidecarSpec for jobSpec: %w", err)
			}
			jb.LegacyGasStationSidecarSpecID = &specID
		case Bootstrap:
			specID, err := tx.insertBootstrapSpec(ctx, jb.BootstrapSpec)
			if err != nil {
				return fmt.Errorf("failed to create BootstrapSpec for jobSpec: %w", err)
			}
			jb.BootstrapSpecID = &specID
		case Gateway:
			specID, err := tx.insertGatewaySpec(ctx, jb.GatewaySpec)
			if err != nil {
				return fmt.Errorf("failed to create GatewaySpec for jobSpec: %w", err)
			}
			jb.GatewaySpecID = &specID
		case Stream:
			// 'stream' type has no associated spec, nothing to do here
		case Workflow:
			sql := `INSERT INTO workflow_specs (workflow, workflow_id, workflow_owner, workflow_name, binary_url, config_url, secrets_id, created_at, updated_at, spec_type, config)
			VALUES (:workflow, :workflow_id, :workflow_owner, :workflow_name, :binary_url, :config_url, :secrets_id, NOW(), NOW(), :spec_type, :config)
			RETURNING id;`
			specID, err := tx.prepareQuerySpecID(ctx, sql, jb.WorkflowSpec)
			if err != nil {
				return fmt.Errorf("failed to create WorkflowSpec for jobSpec given %v: %w", *jb.WorkflowSpec, err)
			}
			jb.WorkflowSpecID = &specID
		case StandardCapabilities:
			sql := `INSERT INTO standardcapabilities_specs (command, config, oracle_factory, created_at, updated_at)
			VALUES (:command, :config, :oracle_factory, NOW(), NOW())
			RETURNING id;`
			specID, err := tx.prepareQuerySpecID(ctx, sql, jb.StandardCapabilitiesSpec)
			if err != nil {
				return errors.Wrap(err, "failed to create StandardCapabilities for jobSpec")
			}
			jb.StandardCapabilitiesSpecID = &specID
		case CCIP:
			sql := `INSERT INTO ccip_specs (
				capability_version,
				capability_labelled_name,
				ocr_key_bundle_ids,
				p2p_key_id,
				p2pv2_bootstrappers,
				relay_configs,
				plugin_config,
				created_at,
				updated_at
			) VALUES (
				:capability_version,
				:capability_labelled_name,
				:ocr_key_bundle_ids,
				:p2p_key_id,
				:p2pv2_bootstrappers,
				:relay_configs,
				:plugin_config,
				NOW(),
				NOW()
			)
			RETURNING id;`
			specID, err := tx.prepareQuerySpecID(ctx, sql, jb.CCIPSpec)
			if err != nil {
				return errors.Wrap(err, "failed to create CCIPSpec for jobSpec")
			}
			jb.CCIPSpecID = &specID
		default:
			o.lggr.Panicf("Unsupported jb.Type: %v", jb.Type)
		}

		pipelineSpecID, err := tx.pipelineORM.CreateSpec(ctx, p, jb.MaxTaskDuration)
		if err != nil {
			return errors.Wrap(err, "failed to create pipeline spec")
		}

		jb.PipelineSpecID = pipelineSpecID

		err = tx.InsertJob(ctx, jb)
		jobID = jb.ID
		return errors.Wrap(err, "failed to insert job")
	})
	if err != nil {
		return errors.Wrap(err, "CreateJobFailed")
	}

	return o.findJob(ctx, jb, "id", jobID)
}

func (o *orm) prepareQuerySpecID(ctx context.Context, sql string, arg any) (specID int32, err error) {
	var stmt *sqlx.NamedStmt
	stmt, err = o.ds.PrepareNamedContext(ctx, sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowxContext(ctx, arg).Scan(&specID)
	return
}

func (o *orm) insertDirectRequestSpec(ctx context.Context, spec *DirectRequestSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO direct_request_specs (contract_address, min_incoming_confirmations, requesters, min_contract_payment, evm_chain_id, created_at, updated_at)
			VALUES (:contract_address, :min_incoming_confirmations, :requesters, :min_contract_payment, :evm_chain_id, now(), now())
			RETURNING id;`, spec)
}

func (o *orm) insertFluxMonitorSpec(ctx context.Context, spec *FluxMonitorSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO flux_monitor_specs (contract_address, threshold, absolute_threshold, poll_timer_period, poll_timer_disabled, idle_timer_period, idle_timer_disabled,
					drumbeat_schedule, drumbeat_random_delay, drumbeat_enabled, min_payment, evm_chain_id, created_at, updated_at)
			VALUES (:contract_address, :threshold, :absolute_threshold, :poll_timer_period, :poll_timer_disabled, :idle_timer_period, :idle_timer_disabled,
					:drumbeat_schedule, :drumbeat_random_delay, :drumbeat_enabled, :min_payment, :evm_chain_id, NOW(), NOW())
			RETURNING id;`, spec)
}

func (o *orm) insertOCROracleSpec(ctx context.Context, spec *OCROracleSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO ocr_oracle_specs (contract_address, p2pv2_bootstrappers, is_bootstrap_peer, encrypted_ocr_key_bundle_id, transmitter_address,
					observation_timeout, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, evm_chain_id,
					created_at, updated_at, database_timeout, observation_grace_period, contract_transmitter_transmit_timeout)
			VALUES (:contract_address, :p2pv2_bootstrappers, :is_bootstrap_peer, :encrypted_ocr_key_bundle_id, :transmitter_address,
					:observation_timeout, :blockchain_timeout, :contract_config_tracker_subscribe_interval, :contract_config_tracker_poll_interval, :contract_config_confirmations, :evm_chain_id,
					NOW(), NOW(), :database_timeout, :observation_grace_period, :contract_transmitter_transmit_timeout)
			RETURNING id;`, spec)
}

func (o *orm) insertOCR2OracleSpec(ctx context.Context, spec *OCR2OracleSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO ocr2_oracle_specs (contract_id, feed_id, relay, relay_config, plugin_type, plugin_config, onchain_signing_strategy, p2pv2_bootstrappers, ocr_key_bundle_id, transmitter_id,
					blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations, allow_no_bootstrappers,
					created_at, updated_at)
			VALUES (:contract_id, :feed_id, :relay, :relay_config, :plugin_type, :plugin_config, :onchain_signing_strategy, :p2pv2_bootstrappers, :ocr_key_bundle_id, :transmitter_id,
					 :blockchain_timeout, :contract_config_tracker_poll_interval, :contract_config_confirmations, :allow_no_bootstrappers,
					NOW(), NOW())
			RETURNING id;`, spec)
}

func (o *orm) insertKeeperSpec(ctx context.Context, spec *KeeperSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO keeper_specs (contract_address, from_address, evm_chain_id, created_at, updated_at)
			VALUES (:contract_address, :from_address, :evm_chain_id, NOW(), NOW())
			RETURNING id;`, spec)
}

func (o *orm) insertCronSpec(ctx context.Context, spec *CronSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO cron_specs (cron_schedule, evm_chain_id, created_at, updated_at)
			VALUES (:cron_schedule, :evm_chain_id, NOW(), NOW())
			RETURNING id;`, spec)
}

func (o *orm) insertVRFSpec(ctx context.Context, spec *VRFSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO vrf_specs (
				coordinator_address, public_key, min_incoming_confirmations,
				evm_chain_id, from_addresses, poll_period, requested_confs_delay,
				request_timeout, chunk_size, batch_coordinator_address, batch_fulfillment_enabled,
				batch_fulfillment_gas_multiplier, backoff_initial_delay, backoff_max_delay, gas_lane_price,
                vrf_owner_address, custom_reverts_pipeline_enabled,
				created_at, updated_at)
			VALUES (
				:coordinator_address, :public_key, :min_incoming_confirmations,
				:evm_chain_id, :from_addresses, :poll_period, :requested_confs_delay,
				:request_timeout, :chunk_size, :batch_coordinator_address, :batch_fulfillment_enabled,
				:batch_fulfillment_gas_multiplier, :backoff_initial_delay, :backoff_max_delay, :gas_lane_price,
			    :vrf_owner_address, :custom_reverts_pipeline_enabled,
				NOW(), NOW())
			RETURNING id;`, toVRFSpecRow(spec))
}

func (o *orm) insertBlockhashStoreSpec(ctx context.Context, spec *BlockhashStoreSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO blockhash_store_specs (coordinator_v1_address, coordinator_v2_address, coordinator_v2_plus_address, trusted_blockhash_store_address, trusted_blockhash_store_batch_size, wait_blocks, lookback_blocks, heartbeat_period, blockhash_store_address, poll_period, run_timeout, evm_chain_id, from_addresses, created_at, updated_at)
			VALUES (:coordinator_v1_address, :coordinator_v2_address, :coordinator_v2_plus_address, :trusted_blockhash_store_address, :trusted_blockhash_store_batch_size, :wait_blocks, :lookback_blocks, :heartbeat_period, :blockhash_store_address, :poll_period, :run_timeout, :evm_chain_id, :from_addresses, NOW(), NOW())
			RETURNING id;`, toBlockhashStoreSpecRow(spec))
}

func (o *orm) insertBlockHeaderFeederSpec(ctx context.Context, spec *BlockHeaderFeederSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO block_header_feeder_specs (coordinator_v1_address, coordinator_v2_address, coordinator_v2_plus_address, wait_blocks, lookback_blocks, blockhash_store_address, batch_blockhash_store_address, poll_period, run_timeout, evm_chain_id, from_addresses, get_blockhashes_batch_size, store_blockhashes_batch_size, created_at, updated_at)
			VALUES (:coordinator_v1_address, :coordinator_v2_address, :coordinator_v2_plus_address, :wait_blocks, :lookback_blocks, :blockhash_store_address, :batch_blockhash_store_address, :poll_period, :run_timeout, :evm_chain_id, :from_addresses,  :get_blockhashes_batch_size, :store_blockhashes_batch_size, NOW(), NOW())
			RETURNING id;`, toBlockHeaderFeederSpecRow(spec))
}

func (o *orm) insertLegacyGasStationServerSpec(ctx context.Context, spec *LegacyGasStationServerSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO legacy_gas_station_server_specs (forwarder_address, evm_chain_id, ccip_chain_selector, from_addresses, created_at, updated_at)
			VALUES (:forwarder_address, :evm_chain_id, :ccip_chain_selector, :from_addresses, NOW(), NOW())
			RETURNING id;`, toLegacyGasStationServerSpecRow(spec))
}

func (o *orm) insertLegacyGasStationSidecarSpec(ctx context.Context, spec *LegacyGasStationSidecarSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO legacy_gas_station_sidecar_specs (forwarder_address, off_ramp_address, lookback_blocks, poll_period, run_timeout, evm_chain_id, ccip_chain_selector, created_at, updated_at)
			VALUES (:forwarder_address, :off_ramp_address, :lookback_blocks, :poll_period, :run_timeout, :evm_chain_id, :ccip_chain_selector, NOW(), NOW())
			RETURNING id;`, spec)
}

func (o *orm) insertBootstrapSpec(ctx context.Context, spec *BootstrapSpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO bootstrap_specs (contract_id, feed_id, relay, relay_config, monitoring_endpoint,
					blockchain_timeout, contract_config_tracker_poll_interval,
					contract_config_confirmations, created_at, updated_at)
			VALUES (:contract_id, :feed_id, :relay, :relay_config, :monitoring_endpoint,
					:blockchain_timeout, :contract_config_tracker_poll_interval,
					:contract_config_confirmations, NOW(), NOW())
			RETURNING id;`, spec)
}

func (o *orm) insertGatewaySpec(ctx context.Context, spec *GatewaySpec) (specID int32, err error) {
	return o.prepareQuerySpecID(ctx, `INSERT INTO gateway_specs (gateway_config, created_at, updated_at)
			VALUES (:gateway_config, NOW(), NOW())
			RETURNING id;`, spec)
}

// ValidateKeyStoreMatch confirms that the key has a valid match in the keystore
func ValidateKeyStoreMatch(ctx context.Context, spec *OCR2OracleSpec, keyStore keystore.Master, key string) (err error) {
	switch spec.PluginType {
	case types.Mercury, types.LLO:
		_, err = keyStore.CSA().Get(key)
		if err != nil {
			err = errors.Errorf("no CSA key matching: %q", key)
		}
	default:
		err = validateKeyStoreMatchForRelay(ctx, spec.Relay, keyStore, key)
	}
	return
}

func validateKeyStoreMatchForRelay(ctx context.Context, network string, keyStore keystore.Master, key string) error {
	switch network {
	case relay.NetworkEVM:
		_, err := keyStore.Eth().Get(ctx, key)
		if err != nil {
			return errors.Errorf("no EVM key matching: %q", key)
		}
	case relay.NetworkCosmos:
		_, err := keyStore.Cosmos().Get(key)
		if err != nil {
			return errors.Errorf("no Cosmos key matching: %q", key)
		}
	case relay.NetworkSolana:
		_, err := keyStore.Solana().Get(key)
		if err != nil {
			return errors.Errorf("no Solana key matching: %q", key)
		}
	case relay.NetworkStarkNet:
		_, err := keyStore.StarkNet().Get(key)
		if err != nil {
			return errors.Errorf("no Starknet key matching: %q", key)
		}
	case relay.NetworkAptos:
		_, err := keyStore.Aptos().Get(key)
		if err != nil {
			return errors.Errorf("no Aptos key matching: %q", key)
		}
	case relay.NetworkTron:
		_, err := keyStore.Tron().Get(key)
		if err != nil {
			return errors.Errorf("no Tron key matching: %q", key)
		}
	}
	return nil
}

func areSendingKeysDefined(ctx context.Context, jb *Job, keystore keystore.Master) (bool, error) {
	if jb.OCR2OracleSpec.RelayConfig["sendingKeys"] != nil {
		sendingKeys, err := SendingKeysForJob(jb)
		if err != nil {
			return false, err
		}

		for _, sendingKey := range sendingKeys {
			if err = ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keystore, sendingKey); err != nil {
				return false, errors.Wrap(ErrNoSuchSendingKey, err.Error())
			}
		}

		return true, nil
	}
	return false, nil
}

func (o *orm) InsertWebhookSpec(ctx context.Context, webhookSpec *WebhookSpec) error {
	query, args, err := o.ds.BindNamed(`INSERT INTO webhook_specs (created_at, updated_at)
			VALUES (NOW(), NOW())
			RETURNING *;`, webhookSpec)
	if err != nil {
		return fmt.Errorf("error binding arg: %w", err)
	}
	return o.ds.GetContext(ctx, webhookSpec, query, args...)
}

func (o *orm) InsertJob(ctx context.Context, job *Job) error {
	return o.transact(ctx, false, func(tx *orm) error {
		var query string

		// if job has id, emplace otherwise insert with a new id.
		if job.ID == 0 {
			query = `INSERT INTO jobs (name, stream_id, schema_version, type, max_task_duration, ocr_oracle_spec_id, ocr2_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id,
				keeper_spec_id, cron_spec_id, vrf_spec_id, webhook_spec_id, blockhash_store_spec_id, bootstrap_spec_id, block_header_feeder_spec_id, gateway_spec_id,
                legacy_gas_station_server_spec_id, legacy_gas_station_sidecar_spec_id, workflow_spec_id, standard_capabilities_spec_id, ccip_spec_id, external_job_id, gas_limit, forwarding_allowed, created_at)
		VALUES (:name, :stream_id, :schema_version, :type, :max_task_duration, :ocr_oracle_spec_id, :ocr2_oracle_spec_id, :direct_request_spec_id, :flux_monitor_spec_id,
				:keeper_spec_id, :cron_spec_id, :vrf_spec_id, :webhook_spec_id, :blockhash_store_spec_id, :bootstrap_spec_id, :block_header_feeder_spec_id, :gateway_spec_id,
				:legacy_gas_station_server_spec_id, :legacy_gas_station_sidecar_spec_id, :workflow_spec_id, :standard_capabilities_spec_id, :ccip_spec_id, :external_job_id, :gas_limit, :forwarding_allowed, NOW())
		RETURNING *;`
		} else {
			query = `INSERT INTO jobs (id, name, stream_id, schema_version, type, max_task_duration, ocr_oracle_spec_id, ocr2_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id,
			keeper_spec_id, cron_spec_id, vrf_spec_id, webhook_spec_id, blockhash_store_spec_id, bootstrap_spec_id, block_header_feeder_spec_id, gateway_spec_id,
                  legacy_gas_station_server_spec_id, legacy_gas_station_sidecar_spec_id, workflow_spec_id, standard_capabilities_spec_id, ccip_spec_id, external_job_id, gas_limit, forwarding_allowed, created_at)
		VALUES (:id, :name, :stream_id, :schema_version, :type, :max_task_duration, :ocr_oracle_spec_id, :ocr2_oracle_spec_id, :direct_request_spec_id, :flux_monitor_spec_id,
				:keeper_spec_id, :cron_spec_id, :vrf_spec_id, :webhook_spec_id, :blockhash_store_spec_id, :bootstrap_spec_id, :block_header_feeder_spec_id, :gateway_spec_id,
				:legacy_gas_station_server_spec_id, :legacy_gas_station_sidecar_spec_id, :workflow_spec_id, :standard_capabilities_spec_id, :ccip_spec_id, :external_job_id, :gas_limit, :forwarding_allowed, NOW())
		RETURNING *;`
		}
		query, args, err := tx.ds.BindNamed(query, job)
		if err != nil {
			return fmt.Errorf("error binding arg: %w", err)
		}
		err = tx.ds.GetContext(ctx, job, query, args...)
		if err != nil {
			return err
		}

		// Always inserts the `job_pipeline_specs` record as primary, since this is the first one for the job.
		sqlStmt := `INSERT INTO job_pipeline_specs (job_id, pipeline_spec_id, is_primary) VALUES ($1, $2, true)`
		_, err = tx.ds.ExecContext(ctx, sqlStmt, job.ID, job.PipelineSpecID)
		return errors.Wrap(err, "failed to insert job_pipeline_specs relationship")
	})
}

// DeleteJob removes a job
func (o *orm) DeleteJob(ctx context.Context, id int32, jobType Type) error {
	o.lggr.Debugw("Deleting job", "jobID", id)
	queries := map[Type]string{
		DirectRequest:        `DELETE FROM direct_request_specs WHERE id IN (SELECT direct_request_spec_id FROM deleted_jobs)`,
		FluxMonitor:          `DELETE FROM flux_monitor_specs WHERE id IN (SELECT flux_monitor_spec_id FROM deleted_jobs)`,
		OffchainReporting:    `DELETE FROM ocr_oracle_specs WHERE id IN (SELECT ocr_oracle_spec_id FROM deleted_jobs)`,
		OffchainReporting2:   `DELETE FROM ocr2_oracle_specs WHERE id IN (SELECT ocr2_oracle_spec_id FROM deleted_jobs)`,
		Keeper:               `DELETE FROM keeper_specs WHERE id IN (SELECT keeper_spec_id FROM deleted_jobs)`,
		Cron:                 `DELETE FROM cron_specs WHERE id IN (SELECT cron_spec_id FROM deleted_jobs)`,
		VRF:                  `DELETE FROM vrf_specs WHERE id IN (SELECT vrf_spec_id FROM deleted_jobs)`,
		Webhook:              `DELETE FROM webhook_specs WHERE id IN (SELECT webhook_spec_id FROM deleted_jobs)`,
		BlockhashStore:       `DELETE FROM blockhash_store_specs WHERE id IN (SELECT blockhash_store_spec_id FROM deleted_jobs)`,
		Bootstrap:            `DELETE FROM bootstrap_specs WHERE id IN (SELECT bootstrap_spec_id FROM deleted_jobs)`,
		BlockHeaderFeeder:    `DELETE FROM block_header_feeder_specs WHERE id IN (SELECT block_header_feeder_spec_id FROM deleted_jobs)`,
		Gateway:              `DELETE FROM gateway_specs WHERE id IN (SELECT gateway_spec_id FROM deleted_jobs)`,
		Workflow:             `DELETE FROM workflow_specs WHERE id in (SELECT workflow_spec_id FROM deleted_jobs)`,
		StandardCapabilities: `DELETE FROM standardcapabilities_specs WHERE id in (SELECT standard_capabilities_spec_id FROM deleted_jobs)`,
		CCIP:                 `DELETE FROM ccip_specs WHERE id in (SELECT ccip_spec_id FROM deleted_jobs)`,
		Stream:               ``,
	}
	q, ok := queries[jobType]
	if !ok {
		return errors.Errorf("job type %s not supported", jobType)
	}
	// Added a 1-minute timeout to this query since this can take a long time as data increases.
	// This was added specifically due to an issue with a database that had a million of pipeline_runs and pipeline_task_runs
	// and this query was taking ~40secs.
	ctx, cancel := context.WithTimeout(sqlutil.WithoutDefaultTimeout(ctx), time.Minute)
	defer cancel()
	query := `
		WITH deleted_jobs AS (
			DELETE FROM jobs WHERE id = $1 RETURNING
				id,
				ocr_oracle_spec_id,
				ocr2_oracle_spec_id,
				keeper_spec_id,
				cron_spec_id,
				flux_monitor_spec_id,
				vrf_spec_id,
				webhook_spec_id,
				direct_request_spec_id,
				blockhash_store_spec_id,
				bootstrap_spec_id,
				block_header_feeder_spec_id,
				gateway_spec_id,
				workflow_spec_id,
				standard_capabilities_spec_id,
				ccip_spec_id,
				stream_id
		),`
	if len(q) > 0 {
		query += fmt.Sprintf(`deleted_specific_specs AS (
								%s
							),`, q)
	}
	query += `
		deleted_job_pipeline_specs AS (
			DELETE FROM job_pipeline_specs WHERE job_id IN (SELECT id FROM deleted_jobs) RETURNING pipeline_spec_id
		)
		DELETE FROM pipeline_specs WHERE id IN (SELECT pipeline_spec_id FROM deleted_job_pipeline_specs)`
	res, err := o.ds.ExecContext(ctx, query, id)
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

func (o *orm) RecordError(ctx context.Context, jobID int32, description string) error {
	sql := `INSERT INTO job_spec_errors (job_id, description, occurrences, created_at, updated_at)
	VALUES ($1, $2, 1, $3, $3)
	ON CONFLICT (job_id, description) DO UPDATE SET
	occurrences = job_spec_errors.occurrences + 1,
	updated_at = excluded.updated_at`
	_, err := o.ds.ExecContext(ctx, sql, jobID, description, time.Now())
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
func (o *orm) TryRecordError(ctx context.Context, jobID int32, description string) {
	err := o.RecordError(ctx, jobID, description)
	o.lggr.ErrorIf(err, fmt.Sprintf("Error creating SpecError %v", description))
}

func (o *orm) DismissError(ctx context.Context, ID int64) error {
	res, err := o.ds.ExecContext(ctx, "DELETE FROM job_spec_errors WHERE id = $1", ID)
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

func (o *orm) FindSpecError(ctx context.Context, id int64) (SpecError, error) {
	stmt := `SELECT * FROM job_spec_errors WHERE id = $1;`

	specErr := new(SpecError)
	err := o.ds.GetContext(ctx, specErr, stmt, id)

	return *specErr, errors.Wrap(err, "FindSpecError failed")
}

func (o *orm) FindJobs(ctx context.Context, offset, limit int) (jobs []Job, count int, err error) {
	err = o.transact(ctx, false, func(tx *orm) error {
		sql := `SELECT count(*) FROM jobs;`
		err = tx.ds.QueryRowxContext(ctx, sql).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to query jobs count: %w", err)
		}

		sql = `SELECT jobs.*, job_pipeline_specs.pipeline_spec_id as pipeline_spec_id
			FROM jobs
			    JOIN job_pipeline_specs ON (jobs.id = job_pipeline_specs.job_id)
			ORDER BY jobs.created_at DESC, jobs.id DESC OFFSET $1 LIMIT $2;`
		err = tx.ds.SelectContext(ctx, &jobs, sql, offset, limit)
		if err != nil {
			return fmt.Errorf("failed to select jobs: %w", err)
		}

		err = tx.loadAllJobsTypes(ctx, jobs)
		if err != nil {
			return fmt.Errorf("failed to load job types: %w", err)
		}

		return nil
	})
	return jobs, count, err
}

func LoadDefaultVRFPollPeriod(vrfs VRFSpec) *VRFSpec {
	if vrfs.PollPeriod == 0 {
		vrfs.PollPeriod = 5 * time.Second
	}

	return &vrfs
}

// SetDRMinIncomingConfirmations takes the largest of the global vs specific.
func SetDRMinIncomingConfirmations(defaultMinIncomingConfirmations uint32, drs DirectRequestSpec) *DirectRequestSpec {
	if !drs.MinIncomingConfirmations.Valid || drs.MinIncomingConfirmations.Uint32 < defaultMinIncomingConfirmations {
		drs.MinIncomingConfirmations = null.Uint32From(defaultMinIncomingConfirmations)
	}
	return &drs
}

type OCRConfig interface {
	BlockchainTimeout() time.Duration
	CaptureEATelemetry() bool
	ContractPollInterval() time.Duration
	ContractSubscribeInterval() time.Duration
	KeyBundleID() (string, error)
	ObservationTimeout() time.Duration
	TransmitterAddress() (evmtypes.EIP55Address, error)
}

// LoadConfigVarsLocalOCR loads local OCR vars into the OCROracleSpec.
func LoadConfigVarsLocalOCR(evmOcrCfg evmconfig.OCR, os OCROracleSpec, ocrCfg OCRConfig) *OCROracleSpec {
	if os.ObservationTimeout == 0 {
		os.ObservationTimeout = models.Interval(ocrCfg.ObservationTimeout())
	}
	if os.BlockchainTimeout == 0 {
		os.BlockchainTimeout = models.Interval(ocrCfg.BlockchainTimeout())
	}
	if os.ContractConfigTrackerSubscribeInterval == 0 {
		os.ContractConfigTrackerSubscribeInterval = models.Interval(ocrCfg.ContractSubscribeInterval())
	}
	if os.ContractConfigTrackerPollInterval == 0 {
		os.ContractConfigTrackerPollInterval = models.Interval(ocrCfg.ContractPollInterval())
	}
	if os.ContractConfigConfirmations == 0 {
		os.ContractConfigConfirmations = evmOcrCfg.ContractConfirmations()
	}
	if os.DatabaseTimeout == nil {
		os.DatabaseTimeout = models.NewInterval(evmOcrCfg.DatabaseTimeout())
	}
	if os.ObservationGracePeriod == nil {
		os.ObservationGracePeriod = models.NewInterval(evmOcrCfg.ObservationGracePeriod())
	}
	if os.ContractTransmitterTransmitTimeout == nil {
		os.ContractTransmitterTransmitTimeout = models.NewInterval(evmOcrCfg.ContractTransmitterTransmitTimeout())
	}
	os.CaptureEATelemetry = ocrCfg.CaptureEATelemetry()

	return &os
}

// LoadConfigVarsOCR loads OCR config vars into the OCROracleSpec.
func LoadConfigVarsOCR(evmOcrCfg evmconfig.OCR, ocrCfg OCRConfig, os OCROracleSpec) (*OCROracleSpec, error) {
	if os.TransmitterAddress == nil {
		ta, err := ocrCfg.TransmitterAddress()
		if !errors.Is(errors.Cause(err), config.ErrEnvUnset) {
			if err != nil {
				return nil, err
			}
			os.TransmitterAddress = &ta
		}
	}

	if os.EncryptedOCRKeyBundleID == nil {
		kb, err := ocrCfg.KeyBundleID()
		if err != nil {
			return nil, err
		}
		encryptedOCRKeyBundleID, err := models.Sha256HashFromHex(kb)
		if err != nil {
			return nil, err
		}
		os.EncryptedOCRKeyBundleID = &encryptedOCRKeyBundleID
	}

	return LoadConfigVarsLocalOCR(evmOcrCfg, os, ocrCfg), nil
}

// FindJob returns job by ID, with all relations preloaded
func (o *orm) FindJob(ctx context.Context, id int32) (jb Job, err error) {
	err = o.findJob(ctx, &jb, "id", id)
	return
}

// FindJobWithoutSpecErrors returns a job by ID, without loading SpecVal Errors preloaded
func (o *orm) FindJobWithoutSpecErrors(ctx context.Context, id int32) (jb Job, err error) {
	err = o.transact(ctx, true, func(tx *orm) error {
		stmt := "SELECT jobs.*, job_pipeline_specs.pipeline_spec_id as pipeline_spec_id FROM jobs JOIN job_pipeline_specs ON (jobs.id = job_pipeline_specs.job_id) WHERE jobs.id = $1 LIMIT 1"
		err = tx.ds.GetContext(ctx, &jb, stmt, id)
		if err != nil {
			return errors.Wrap(err, "failed to load job")
		}

		if err = tx.loadAllJobTypes(ctx, &jb); err != nil {
			return errors.Wrap(err, "failed to load job types")
		}

		return nil
	})
	if err != nil {
		return jb, errors.Wrap(err, "FindJobWithoutSpecErrors failed")
	}

	return jb, nil
}

// FindSpecErrorsByJobIDs returns all jobs spec errors by jobs IDs
func (o *orm) FindSpecErrorsByJobIDs(ctx context.Context, ids []int32) ([]SpecError, error) {
	stmt := `SELECT * FROM job_spec_errors WHERE job_id = ANY($1);`

	var specErrs []SpecError
	err := o.ds.SelectContext(ctx, &specErrs, stmt, ids)

	return specErrs, errors.Wrap(err, "FindSpecErrorsByJobIDs failed")
}

func (o *orm) FindJobByExternalJobID(ctx context.Context, externalJobID uuid.UUID) (jb Job, err error) {
	err = o.findJob(ctx, &jb, "external_job_id", externalJobID)
	return
}

// FindJobIDByAddress - finds a job id by contract address. Currently only OCR and FM jobs are supported
func (o *orm) FindJobIDByAddress(ctx context.Context, address evmtypes.EIP55Address, evmChainID *big.Big) (jobID int32, err error) {
	stmt := `
SELECT jobs.id
FROM jobs
LEFT JOIN ocr_oracle_specs ocrspec on ocrspec.contract_address = $1 AND (ocrspec.evm_chain_id = $2 OR ocrspec.evm_chain_id IS NULL) AND ocrspec.id = jobs.ocr_oracle_spec_id
LEFT JOIN flux_monitor_specs fmspec on fmspec.contract_address = $1 AND (fmspec.evm_chain_id = $2 OR fmspec.evm_chain_id IS NULL) AND fmspec.id = jobs.flux_monitor_spec_id
WHERE ocrspec.id IS NOT NULL OR fmspec.id IS NOT NULL
`
	err = o.ds.GetContext(ctx, &jobID, stmt, address, evmChainID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = errors.Wrap(err, "error searching for job by contract address")
		}
		err = errors.Wrap(err, "FindJobIDByAddress failed")
		return
	}

	return
}

func (o *orm) FindOCR2JobIDByAddress(ctx context.Context, contractID string, feedID *common.Hash) (jobID int32, err error) {
	// NOTE: We want to explicitly match on NULL feed_id hence usage of `IS
	// NOT DISTINCT FROM` instead of `=`
	stmt := `
SELECT jobs.id
FROM jobs
LEFT JOIN ocr2_oracle_specs ocr2spec on ocr2spec.contract_id = $1 AND ocr2spec.feed_id IS NOT DISTINCT FROM $2 AND ocr2spec.id = jobs.ocr2_oracle_spec_id
LEFT JOIN bootstrap_specs bs on bs.contract_id = $1 AND bs.feed_id IS NOT DISTINCT FROM $2 AND bs.id = jobs.bootstrap_spec_id
WHERE ocr2spec.id IS NOT NULL OR bs.id IS NOT NULL
`
	err = o.ds.GetContext(ctx, &jobID, stmt, contractID, feedID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = errors.Wrapf(err, "error searching for job by contract id=%s and feed id=%s", contractID, feedID)
		}
		err = errors.Wrap(err, "FindOCR2JobIDByAddress failed")
		return
	}

	return
}

func (o *orm) findJob(ctx context.Context, jb *Job, col string, arg interface{}) error {
	err := o.transact(ctx, false, func(tx *orm) error {
		sql := fmt.Sprintf(`SELECT jobs.*, job_pipeline_specs.pipeline_spec_id FROM jobs JOIN job_pipeline_specs ON (jobs.id = job_pipeline_specs.job_id) WHERE jobs.%s = $1 AND job_pipeline_specs.is_primary = true LIMIT 1`, col)
		err := tx.ds.GetContext(ctx, jb, sql, arg)
		if err != nil {
			return errors.Wrap(err, "failed to load job")
		}

		if err = tx.loadAllJobTypes(ctx, jb); err != nil {
			return err
		}

		return tx.loadJobSpecErrors(ctx, jb)
	})
	if err != nil {
		return errors.Wrap(err, "findJob failed")
	}
	return nil
}

func (o *orm) FindJobIDsWithBridge(ctx context.Context, name string) (jids []int32, err error) {
	query := `SELECT
			jobs.id, pipeline_specs.dot_dag_source
		FROM jobs
		    JOIN job_pipeline_specs ON job_pipeline_specs.job_id = jobs.id
		    JOIN pipeline_specs ON pipeline_specs.id = job_pipeline_specs.pipeline_spec_id
		WHERE pipeline_specs.dot_dag_source ILIKE '%' || $1 || '%' ORDER BY id`
	var rows *sqlx.Rows
	rows, err = o.ds.QueryxContext(ctx, query, name)
	if err != nil {
		return
	}
	defer rows.Close()
	var ids []int32
	var sources []string
	for rows.Next() {
		var id int32
		var source string
		if err = rows.Scan(&id, &source); err != nil {
			return
		}
		ids = append(jids, id)
		sources = append(sources, source)
	}
	if err = rows.Err(); err != nil {
		return
	}

	for i, id := range ids {
		var p *pipeline.Pipeline
		p, err = pipeline.Parse(sources[i])
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse dag for job %d", id)
		}
		for _, task := range p.Tasks {
			if task.Type() == pipeline.TaskTypeBridge {
				if task.(*pipeline.BridgeTask).Name == name {
					jids = append(jids, id)
				}
			}
		}
	}

	return
}

func (o *orm) FindJobIDByWorkflow(ctx context.Context, spec WorkflowSpec) (jobID int32, err error) {
	stmt := `
SELECT jobs.id FROM jobs
INNER JOIN workflow_specs ws on jobs.workflow_spec_id = ws.id AND ws.workflow_owner = $1 AND ws.workflow_name = $2
`
	err = o.ds.GetContext(ctx, &jobID, stmt, spec.WorkflowOwner, spec.WorkflowName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("error searching for job by workflow (owner,name) ('%s','%s'): %w", spec.WorkflowOwner, spec.WorkflowName, err)
		}
		err = fmt.Errorf("FindJobIDByWorkflow failed: %w", err)
		return
	}

	return
}

func (o *orm) FindJobIDByCapabilityNameAndVersion(ctx context.Context, spec CCIPSpec) (jobID int32, err error) {
	stmt := `
SELECT jobs.id FROM jobs
INNER JOIN ccip_specs ccip on jobs.ccip_spec_id = ccip.id AND ccip.capability_labelled_name = $1 AND ccip.capability_version = $2
`
	err = o.ds.GetContext(ctx, &jobID, stmt, spec.CapabilityLabelledName, spec.CapabilityVersion)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("error searching for job for CCIP (capabilityName,capabilityVersion) ('%s','%s'): %w", spec.CapabilityLabelledName, spec.CapabilityVersion, err)
	}
	return
}

func (o *orm) FindStandardCapabilityJobID(ctx context.Context, spec StandardCapabilitiesSpec) (jobID int32, err error) {
	stmt := `
SELECT jobs.id FROM jobs
INNER JOIN standardcapabilities_specs sc on jobs.standard_capabilities_spec_id = sc.id AND sc.command = $1
`
	err = o.ds.GetContext(ctx, &jobID, stmt, spec.Command)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("error searching for job for standardcapabilities (command) ('%s'): %w", spec.Command, err)
	}
	return
}

func (o *orm) FindGatewayJobID(ctx context.Context, spec GatewaySpec) (jobID int32, err error) {
	stmt := `
SELECT jobs.id FROM jobs
INNER JOIN gateway_specs gs on jobs.gateway_spec_id = gs.id
WHERE gs.gateway_config @> jsonb_build_object('ConnectionManagerConfig', jsonb_build_object('AuthGatewayId', $1::text));`
	err = o.ds.GetContext(ctx, &jobID, stmt, spec.AuthGatewayID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("error searching for job for gateway (ConnectionManagerConfig.AuthGatewayId) ('%s'): %w", spec.AuthGatewayID(), err)
	}
	return
}

// PipelineRunsByJobsIDs returns pipeline runs for multiple jobs, not preloading data
func (o *orm) PipelineRunsByJobsIDs(ctx context.Context, ids []int32) (runs []pipeline.Run, err error) {
	err = o.transact(ctx, false, func(tx *orm) error {
		stmt := `SELECT pipeline_runs.* FROM pipeline_runs INNER JOIN job_pipeline_specs ON pipeline_runs.pipeline_spec_id = job_pipeline_specs.pipeline_spec_id WHERE jobs.id = ANY($1)
		ORDER BY pipeline_runs.created_at DESC, pipeline_runs.id DESC;`
		if err = tx.ds.SelectContext(ctx, &runs, stmt, ids); err != nil {
			return errors.Wrap(err, "error loading runs")
		}

		runs, err = tx.loadPipelineRunsRelations(ctx, runs)

		return err
	})

	return runs, errors.Wrap(err, "PipelineRunsByJobsIDs failed")
}

func (o *orm) loadPipelineRunIDs(ctx context.Context, jobID *int32, offset, limit int) (ids []int64, err error) {
	lggr := logger.Sugared(o.lggr)

	var res sql.NullInt64
	if err = o.ds.GetContext(ctx, &res, "SELECT MAX(id) FROM pipeline_runs"); err != nil {
		err = errors.Wrap(err, "error while loading runs")
		return
	} else if !res.Valid {
		// MAX() will return NULL if there are no rows in table.  This is not an error
		return
	}
	maxID := res.Int64

	var filter string
	if jobID != nil {
		filter = fmt.Sprintf("JOIN job_pipeline_specs USING(pipeline_spec_id) WHERE job_pipeline_specs.job_id = %d AND ", *jobID)
	} else {
		filter = "WHERE "
	}

	stmt := fmt.Sprintf(`SELECT p.id FROM pipeline_runs AS p %s p.id >= $3 AND p.id <= $4
			ORDER BY p.id DESC OFFSET $1 LIMIT $2`, filter)

	// Only search the most recent n pipeline runs (whether deleted or not), starting with n = 1000 and
	//  doubling only if we still need more.  Without this, large tables can result in the UI
	//  becoming unusably slow, continuously flashing, or timing out.  The ORDER BY in
	//  this query requires a sort of all runs matching jobID, so we restrict it to the
	//  range minID <-> maxID.

	for n := int64(1000); maxID > 0 && len(ids) < limit; n *= 2 {
		var batch []int64
		minID := maxID - n
		if err = o.ds.SelectContext(ctx, &batch, stmt, offset, limit-len(ids), minID, maxID); err != nil {
			err = errors.Wrap(err, "error loading runs")
			return
		}
		ids = append(ids, batch...)
		if offset > 0 {
			if len(ids) > 0 {
				// If we're already receiving rows back, then we no longer need an offset
				offset = 0
			} else {
				var skipped int
				// If no rows were returned, we need to know whether there were any ids skipped
				//  in this batch due to the offset, and reduce it for the next batch
				err = o.ds.GetContext(ctx, &skipped,
					fmt.Sprintf(
						`SELECT COUNT(p.id) FROM pipeline_runs AS p %s p.id >= $1 AND p.id <= $2`, filter,
					), minID, maxID,
				)
				if err != nil {
					err = errors.Wrap(err, "error loading from pipeline_runs")
					return
				}
				offset -= skipped
				if offset < 0 { // sanity assertion, if this ever happened it would probably mean db corruption or pg bug
					lggr.AssumptionViolationw("offset < 0 while reading pipeline_runs")
					err = errors.Wrap(err, "internal db error while reading pipeline_runs")
					return
				}
				lggr.Debugw("loadPipelineRunIDs empty batch", "minId", minID, "maxID", maxID, "n", n, "len(ids)", len(ids), "limit", limit, "offset", offset, "skipped", skipped)
			}
		}
		maxID = minID - 1
	}
	return
}

func (o *orm) FindTaskResultByRunIDAndTaskName(ctx context.Context, runID int64, taskName string) (result []byte, err error) {
	stmt := fmt.Sprintf("SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1 AND dot_id = '%s';", taskName)

	var taskRuns []pipeline.TaskRun
	if errB := o.ds.SelectContext(ctx, &taskRuns, stmt, runID); errB != nil {
		return nil, errB
	}
	if len(taskRuns) == 0 {
		return nil, fmt.Errorf("can't find task run with id: %v, taskName: %v", runID, taskName)
	}
	if len(taskRuns) > 1 {
		o.lggr.Errorf("found multiple task runs with id: %v, taskName: %v. Using the first one.", runID, taskName)
	}
	taskRun := taskRuns[0]
	if !taskRun.Error.IsZero() {
		return nil, errors.New(taskRun.Error.ValueOrZero())
	}
	resBytes, errB := taskRun.Output.MarshalJSON()
	if errB != nil {
		return
	}
	result = resBytes

	return
}

// FindPipelineRunIDsByJobID fetches the ids of pipeline runs for a job.
func (o *orm) FindPipelineRunIDsByJobID(ctx context.Context, jobID int32, offset, limit int) (ids []int64, err error) {
	err = o.transact(ctx, false, func(tx *orm) error {
		ids, err = tx.loadPipelineRunIDs(ctx, &jobID, offset, limit)
		return err
	})
	return ids, errors.Wrap(err, "FindPipelineRunIDsByJobID failed")
}

func (o *orm) loadPipelineRunsByID(ctx context.Context, ids []int64) (runs []pipeline.Run, err error) {
	stmt := `
		SELECT pipeline_runs.*
		FROM pipeline_runs
		WHERE id = ANY($1)
		ORDER BY created_at DESC, id DESC
	`
	if err = o.ds.SelectContext(ctx, &runs, stmt, ids); err != nil {
		err = errors.Wrap(err, "error loading runs")
		return
	}

	return o.loadPipelineRunsRelations(ctx, runs)
}

// FindPipelineRunsByIDs returns pipeline runs with the ids.
func (o *orm) FindPipelineRunsByIDs(ctx context.Context, ids []int64) (runs []pipeline.Run, err error) {
	err = o.transact(ctx, false, func(tx *orm) error {
		runs, err = tx.loadPipelineRunsByID(ctx, ids)
		return err
	})

	return runs, errors.Wrap(err, "FindPipelineRunsByIDs failed")
}

// FindPipelineRunByID returns pipeline run with the id.
func (o *orm) FindPipelineRunByID(ctx context.Context, id int64) (pipeline.Run, error) {
	var run pipeline.Run

	err := o.transact(ctx, false, func(tx *orm) error {
		stmt := `
SELECT pipeline_runs.*
FROM pipeline_runs
WHERE id = $1
`

		if err := tx.ds.GetContext(ctx, &run, stmt, id); err != nil {
			return errors.Wrap(err, "error loading run")
		}

		runs, err := tx.loadPipelineRunsRelations(ctx, []pipeline.Run{run})

		run = runs[0]

		return err
	})

	return run, errors.Wrap(err, "FindPipelineRunByID failed")
}

// CountPipelineRunsByJobID returns the total number of pipeline runs for a job.
func (o *orm) CountPipelineRunsByJobID(ctx context.Context, jobID int32) (count int32, err error) {
	stmt := "SELECT COUNT(*) FROM pipeline_runs JOIN job_pipeline_specs USING (pipeline_spec_id) WHERE job_pipeline_specs.job_id = $1"
	err = o.ds.GetContext(ctx, &count, stmt, jobID)

	return count, errors.Wrap(err, "CountPipelineRunsByJobID failed")
}

func (o *orm) FindJobsByPipelineSpecIDs(ctx context.Context, ids []int32) ([]Job, error) {
	var jbs []Job

	err := o.transact(ctx, false, func(tx *orm) error {
		stmt := `SELECT jobs.*, job_pipeline_specs.pipeline_spec_id FROM jobs JOIN job_pipeline_specs ON (jobs.id = job_pipeline_specs.job_id) WHERE job_pipeline_specs.pipeline_spec_id = ANY($1) ORDER BY jobs.id ASC
`
		if err := tx.ds.SelectContext(ctx, &jbs, stmt, ids); err != nil {
			return errors.Wrap(err, "error fetching jobs by pipeline spec IDs")
		}

		err := tx.loadAllJobsTypes(ctx, jbs)
		if err != nil {
			return err
		}

		return nil
	})

	return jbs, errors.Wrap(err, "FindJobsByPipelineSpecIDs failed")
}

func (o *orm) FindJobIDByStreamID(ctx context.Context, streamID uint32) (jobID int32, err error) {
	stmt := `SELECT id FROM jobs WHERE type = 'stream' AND stream_id = $1`
	err = o.ds.GetContext(ctx, &jobID, stmt, streamID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = errors.Wrap(err, "error searching for job by stream id")
		}
		err = errors.Wrap(err, "FindJobIDByStreamID failed")
		return
	}

	return
}

// PipelineRuns returns pipeline runs for a job, with spec and taskruns loaded, latest first
// If jobID is nil, returns all pipeline runs
func (o *orm) PipelineRuns(ctx context.Context, jobID *int32, offset, size int) (runs []pipeline.Run, count int, err error) {
	var filter string
	if jobID != nil {
		filter = fmt.Sprintf("JOIN job_pipeline_specs USING(pipeline_spec_id) WHERE job_pipeline_specs.job_id = %d", *jobID)
	}
	err = o.transact(ctx, false, func(tx *orm) error {
		sql := "SELECT count(*) FROM pipeline_runs " + filter
		if err = tx.ds.QueryRowxContext(ctx, sql).Scan(&count); err != nil {
			return errors.Wrap(err, "error counting runs")
		}

		var ids []int64
		ids, err = tx.loadPipelineRunIDs(ctx, jobID, offset, size)
		runs, err = tx.loadPipelineRunsByID(ctx, ids)

		return err
	})

	return runs, count, errors.Wrap(err, "PipelineRuns failed")
}

func (o *orm) loadPipelineRunsRelations(ctx context.Context, runs []pipeline.Run) ([]pipeline.Run, error) {
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
	stmt := `SELECT pipeline_specs.*, job_pipeline_specs.job_id AS job_id FROM pipeline_specs JOIN job_pipeline_specs ON pipeline_specs.id = job_pipeline_specs.pipeline_spec_id WHERE pipeline_specs.id = ANY($1);`
	var specs []pipeline.Spec
	if err := o.ds.SelectContext(ctx, &specs, stmt, specIDs); err != nil {
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
	if err := o.ds.SelectContext(ctx, &taskRuns, stmt, runIDs); err != nil {
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
func (o *orm) loadAllJobsTypes(ctx context.Context, jobs []Job) error {
	for i := range jobs {
		err := o.loadAllJobTypes(ctx, &jobs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *orm) loadAllJobTypes(ctx context.Context, job *Job) error {
	return multierr.Combine(
		o.loadJobPipelineSpec(ctx, job, &job.PipelineSpecID),
		o.loadJobType(ctx, job, "FluxMonitorSpec", "flux_monitor_specs", job.FluxMonitorSpecID),
		o.loadJobType(ctx, job, "DirectRequestSpec", "direct_request_specs", job.DirectRequestSpecID),
		o.loadJobType(ctx, job, "OCROracleSpec", "ocr_oracle_specs", job.OCROracleSpecID),
		o.loadJobType(ctx, job, "OCR2OracleSpec", "ocr2_oracle_specs", job.OCR2OracleSpecID),
		o.loadJobType(ctx, job, "KeeperSpec", "keeper_specs", job.KeeperSpecID),
		o.loadJobType(ctx, job, "CronSpec", "cron_specs", job.CronSpecID),
		o.loadJobType(ctx, job, "WebhookSpec", "webhook_specs", job.WebhookSpecID),
		o.loadVRFJob(ctx, job, job.VRFSpecID),
		o.loadBlockhashStoreJob(ctx, job, job.BlockhashStoreSpecID),
		o.loadBlockHeaderFeederJob(ctx, job, job.BlockHeaderFeederSpecID),
		o.loadLegacyGasStationServerJob(ctx, job, job.LegacyGasStationServerSpecID),
		o.loadJobType(ctx, job, "LegacyGasStationSidecarSpec", "legacy_gas_station_sidecar_specs", job.LegacyGasStationSidecarSpecID),
		o.loadJobType(ctx, job, "BootstrapSpec", "bootstrap_specs", job.BootstrapSpecID),
		o.loadJobType(ctx, job, "GatewaySpec", "gateway_specs", job.GatewaySpecID),
		o.loadJobType(ctx, job, "WorkflowSpec", "workflow_specs", job.WorkflowSpecID),
		o.loadJobType(ctx, job, "StandardCapabilitiesSpec", "standardcapabilities_specs", job.StandardCapabilitiesSpecID),
		o.loadJobType(ctx, job, "CCIPSpec", "ccip_specs", job.CCIPSpecID),
	)
}

func (o *orm) loadJobType(ctx context.Context, job *Job, field, table string, id *int32) error {
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

	err := o.ds.GetContext(ctx, dest, fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, table), *id)

	if err != nil {
		return errors.Wrapf(err, "failed to load job type %s with id %d", table, *id)
	}
	reflect.ValueOf(job).Elem().FieldByName(field).Set(destVal)
	return nil
}

func (o *orm) loadJobPipelineSpec(ctx context.Context, job *Job, id *int32) error {
	if id == nil {
		return nil
	}
	pipelineSpecRow := new(pipeline.Spec)
	if job.PipelineSpec != nil {
		pipelineSpecRow = job.PipelineSpec
	}
	err := o.ds.GetContext(
		ctx,
		pipelineSpecRow,
		`SELECT pipeline_specs.*, job_pipeline_specs.job_id as job_id
			FROM pipeline_specs
    		JOIN job_pipeline_specs ON(pipeline_specs.id = job_pipeline_specs.pipeline_spec_id)
        	WHERE job_pipeline_specs.job_id = $1 AND job_pipeline_specs.pipeline_spec_id = $2`,
		job.ID, *id,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to load job type PipelineSpec with id %d", *id)
	}
	job.PipelineSpec = pipelineSpecRow
	return nil
}

func (o *orm) loadVRFJob(ctx context.Context, job *Job, id *int32) error {
	if id == nil {
		return nil
	}

	var row vrfSpecRow
	err := o.ds.GetContext(ctx, &row, `SELECT * FROM vrf_specs WHERE id = $1`, *id)
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
			evmtypes.EIP55AddressFromAddress(common.BytesToAddress(a)))
	}
	return r.VRFSpec
}

func (o *orm) loadBlockhashStoreJob(ctx context.Context, job *Job, id *int32) error {
	if id == nil {
		return nil
	}

	var row blockhashStoreSpecRow
	err := o.ds.GetContext(ctx, &row, `SELECT * FROM blockhash_store_specs WHERE id = $1`, *id)
	if err != nil {
		return errors.Wrapf(err, `failed to load job type BlockhashStoreSpec with id %d`, *id)
	}

	job.BlockhashStoreSpec = row.toBlockhashStoreSpec()
	return nil
}

// blockhashStoreSpecRow is a helper type for reading and writing blockhashStore specs to the database. This is necessary
// because the bytea[] in the DB is not automatically convertible to or from the spec's
// FromAddresses field. pq.ByteaArray must be used instead.
type blockhashStoreSpecRow struct {
	*BlockhashStoreSpec
	FromAddresses pq.ByteaArray
}

func toBlockhashStoreSpecRow(spec *BlockhashStoreSpec) blockhashStoreSpecRow {
	addresses := make(pq.ByteaArray, len(spec.FromAddresses))
	for i, a := range spec.FromAddresses {
		addresses[i] = a.Bytes()
	}
	return blockhashStoreSpecRow{BlockhashStoreSpec: spec, FromAddresses: addresses}
}

func (r blockhashStoreSpecRow) toBlockhashStoreSpec() *BlockhashStoreSpec {
	for _, a := range r.FromAddresses {
		r.BlockhashStoreSpec.FromAddresses = append(r.BlockhashStoreSpec.FromAddresses,
			evmtypes.EIP55AddressFromAddress(common.BytesToAddress(a)))
	}
	return r.BlockhashStoreSpec
}

func (o *orm) loadBlockHeaderFeederJob(ctx context.Context, job *Job, id *int32) error {
	if id == nil {
		return nil
	}

	var row blockHeaderFeederSpecRow
	err := o.ds.GetContext(ctx, &row, `SELECT * FROM block_header_feeder_specs WHERE id = $1`, *id)
	if err != nil {
		return errors.Wrapf(err, `failed to load job type BlockHeaderFeederSpec with id %d`, *id)
	}

	job.BlockHeaderFeederSpec = row.toBlockHeaderFeederSpec()
	return nil
}

// blockHeaderFeederSpecRow is a helper type for reading and writing blockHeaderFeederSpec specs to the database. This is necessary
// because the bytea[] in the DB is not automatically convertible to or from the spec's
// FromAddresses field. pq.ByteaArray must be used instead.
type blockHeaderFeederSpecRow struct {
	*BlockHeaderFeederSpec
	FromAddresses pq.ByteaArray
}

func toBlockHeaderFeederSpecRow(spec *BlockHeaderFeederSpec) blockHeaderFeederSpecRow {
	addresses := make(pq.ByteaArray, len(spec.FromAddresses))
	for i, a := range spec.FromAddresses {
		addresses[i] = a.Bytes()
	}
	return blockHeaderFeederSpecRow{BlockHeaderFeederSpec: spec, FromAddresses: addresses}
}

func (r blockHeaderFeederSpecRow) toBlockHeaderFeederSpec() *BlockHeaderFeederSpec {
	for _, a := range r.FromAddresses {
		r.BlockHeaderFeederSpec.FromAddresses = append(r.BlockHeaderFeederSpec.FromAddresses,
			evmtypes.EIP55AddressFromAddress(common.BytesToAddress(a)))
	}
	return r.BlockHeaderFeederSpec
}

func (o *orm) loadLegacyGasStationServerJob(ctx context.Context, job *Job, id *int32) error {
	if id == nil {
		return nil
	}

	var row legacyGasStationServerSpecRow
	err := o.ds.GetContext(ctx, &row, `SELECT * FROM legacy_gas_station_server_specs WHERE id = $1`, *id)
	if err != nil {
		return errors.Wrapf(err, `failed to load job type LegacyGasStationServerSpec with id %d`, *id)
	}

	job.LegacyGasStationServerSpec = row.toLegacyGasStationServerSpec()
	return nil
}

// legacyGasStationServerSpecRow is a helper type for reading and writing legacyGasStationServerSpec specs to the database. This is necessary
// because the bytea[] in the DB is not automatically convertible to or from the spec's
// FromAddresses field. pq.ByteaArray must be used instead.
type legacyGasStationServerSpecRow struct {
	*LegacyGasStationServerSpec
	FromAddresses pq.ByteaArray
}

func toLegacyGasStationServerSpecRow(spec *LegacyGasStationServerSpec) legacyGasStationServerSpecRow {
	addresses := make(pq.ByteaArray, len(spec.FromAddresses))
	for i, a := range spec.FromAddresses {
		addresses[i] = a.Bytes()
	}
	return legacyGasStationServerSpecRow{LegacyGasStationServerSpec: spec, FromAddresses: addresses}
}

func (r legacyGasStationServerSpecRow) toLegacyGasStationServerSpec() *LegacyGasStationServerSpec {
	for _, a := range r.FromAddresses {
		r.LegacyGasStationServerSpec.FromAddresses = append(r.LegacyGasStationServerSpec.FromAddresses,
			evmtypes.EIP55AddressFromAddress(common.BytesToAddress(a)))
	}
	return r.LegacyGasStationServerSpec
}

func (o *orm) loadJobSpecErrors(ctx context.Context, jb *Job) error {
	return errors.Wrapf(o.ds.SelectContext(ctx, &jb.JobSpecErrors, `SELECT * FROM job_spec_errors WHERE job_id = $1`, jb.ID), "failed to load job spec errors for job %d", jb.ID)
}

func checkIfKeyHasLock(ctx context.Context, ks keystore.Eth, address common.Address, usage evmkeystore.ServiceType) bool {
	rm := ks.GetResourceMutex(ctx, address)

	return rm.IsLocked(usage)
}
