package ocrbootstrap

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/sqlx"

	relaylogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type RelayGetter interface {
	Get(relay.ID) (loop.Relayer, error)
}

// Delegate creates Bootstrap jobs
type Delegate struct {
	db          *sqlx.DB
	jobORM      job.ORM
	peerWrapper *ocrcommon.SingletonPeerWrapper
	ocr2Cfg     validate.OCR2Config
	insecureCfg validate.InsecureConfig
	lggr        logger.SugaredLogger
	RelayGetter
	isNewlyCreatedJob bool
}

// Extra fields to enable router proxy contract support. Must match field names of functions' PluginConfig.
type relayConfigRouterContractFields struct {
	DONID                           string `json:"donID"`
	ContractVersion                 uint32 `json:"contractVersion"`
	ContractUpdateCheckFrequencySec uint32 `json:"contractUpdateCheckFrequencySec"`
}

// NewDelegateBootstrap creates a new Delegate
func NewDelegateBootstrap(
	db *sqlx.DB,
	jobORM job.ORM,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	lggr logger.Logger,
	ocr2Cfg validate.OCR2Config,
	insecureCfg validate.InsecureConfig,
	relayers RelayGetter,
) *Delegate {
	return &Delegate{
		db:          db,
		jobORM:      jobORM,
		peerWrapper: peerWrapper,
		lggr:        logger.Sugared(lggr),
		ocr2Cfg:     ocr2Cfg,
		insecureCfg: insecureCfg,
		RelayGetter: relayers,
	}
}

// JobType satisfies the job.Delegate interface.
func (d *Delegate) JobType() job.Type {
	return job.Bootstrap
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	d.isNewlyCreatedJob = true
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(jb job.Job) (services []job.ServiceCtx, err error) {
	spec := jb.BootstrapSpec
	if spec == nil {
		return nil, errors.Errorf("Bootstrap.Delegate expects an *job.BootstrapSpec to be present, got %v", jb)
	}
	if d.peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !d.peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running p2p v2 peer")
	}
	s := spec.AsOCR2Spec()
	rid, err := s.RelayID()
	if err != nil {
		return nil, fmt.Errorf("ServicesForSpec: could not get relayer: %w", err)
	}

	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, fmt.Errorf("ServiceForSpec: failed to get relay %s is it enabled?: %w", rid.Name(), err)
	}
	if spec.FeedID != nil {
		spec.RelayConfig["feedID"] = *spec.FeedID
	}

	ctxVals := loop.ContextValues{
		JobID:      jb.ID,
		JobName:    jb.Name.ValueOrZero(),
		ContractID: spec.ContractID,
		FeedID:     spec.FeedID,
	}
	ctx := ctxVals.ContextWithValues(context.Background())

	var routerFields relayConfigRouterContractFields
	if err = json.Unmarshal(spec.RelayConfig.Bytes(), &routerFields); err != nil {
		return nil, err
	}

	var configProvider types.ConfigProvider
	if routerFields.DONID != "" {
		if routerFields.ContractVersion != 1 || routerFields.ContractUpdateCheckFrequencySec == 0 {
			return nil, errors.New("invalid router contract config")
		}
		configProvider, err = relayer.NewPluginProvider(
			ctx,
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         jb.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
				New:           d.isNewlyCreatedJob,
				ProviderType:  string(types.Functions),
			},
			types.PluginArgs{
				PluginConfig: spec.RelayConfig.Bytes(), // contains all necessary fields for config provider
			},
		)
	} else {
		configProvider, err = relayer.NewConfigProvider(ctx, types.RelayArgs{
			ExternalJobID: jb.ExternalJobID,
			JobID:         jb.ID,
			ContractID:    spec.ContractID,
			New:           d.isNewlyCreatedJob,
			RelayConfig:   spec.RelayConfig.Bytes(),
		})
	}

	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relayer.NewConfigWatcher'")
	}
	lc, err := validate.ToLocalConfig(d.ocr2Cfg, d.insecureCfg, spec.AsOCR2Spec())
	if err != nil {
		return nil, err
	}
	if err = ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	lggr := d.lggr.With(ctxVals.Args()...)
	lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)
	bootstrapNodeArgs := ocr.BootstrapperArgs{
		BootstrapperFactory:   d.peerWrapper.Peer2,
		ContractConfigTracker: configProvider.ContractConfigTracker(),
		Database:              NewDB(d.db.DB, spec.ID, lggr),
		LocalConfig:           lc,
		Logger: relaylogger.NewOCRWrapper(lggr.Named("OCRBootstrap"), d.ocr2Cfg.TraceLogging(), func(msg string) {
			logger.Sugared(lggr).ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
		}),
		OffchainConfigDigester: configProvider.OffchainConfigDigester(),
	}
	lggr.Debugw("Launching new bootstrap node", "args", bootstrapNodeArgs)
	bootstrapper, err := ocr.NewBootstrapper(bootstrapNodeArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewBootstrapNode")
	}
	return []job.ServiceCtx{configProvider, job.NewServiceAdapter(bootstrapper)}, nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) AfterJobCreated(spec job.Job) {
}

// BeforeJobDeleted satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// OnDeleteJob satisfies the job.Delegate interface.
func (d *Delegate) OnDeleteJob(spec job.Job, q pg.Queryer) error {
	return nil
}
