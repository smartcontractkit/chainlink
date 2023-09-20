package legacygasstation

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/ccip/generated/evm_2_evm_off_ramp"
	forwarder "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/legacygasstation/generated/legacy_gas_station_forwarder"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type SidecarDelegate struct {
	logger logger.Logger
	chains evm.LegacyChainContainer
	ks     keystore.Eth
	db     *sqlx.DB
}

// JobType satisfies the job.Delegate interface.
func (d *SidecarDelegate) JobType() job.Type {
	return job.LegacyGasStationSidecar
}

// NewDelegate creates a new Delegate.
func NewSidecarDelegate(
	logger logger.Logger,
	chains evm.LegacyChainContainer,
	ks keystore.Eth,
	db *sqlx.DB,
) *SidecarDelegate {
	return &SidecarDelegate{
		logger: logger,
		chains: chains,
		ks:     ks,
		db:     db,
	}
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *SidecarDelegate) ServicesForSpec(jb job.Job, qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	if jb.LegacyGasStationSidecarSpec == nil {
		return nil, errors.Errorf(
			"legacygasstation.Delegate expects a LegacyGasStationSidecarSpec to be present, got %+v", jb)
	}

	chain, err := d.chains.Get(jb.LegacyGasStationSidecarSpec.EVMChainID.String())
	if err != nil {
		return nil, err
	}

	log := d.logger.Named("Legacy Gas Station Sidecar").With("jobID", jb.ID, "externalJobID", jb.ExternalJobID)

	forwarder, err := forwarder.NewLegacyGasStationForwarder(jb.LegacyGasStationSidecarSpec.ForwarderAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "initializing forwarder")
	}

	offramp, err := evm_2_evm_off_ramp.NewEVM2EVMOffRamp(jb.LegacyGasStationSidecarSpec.OffRampAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "initializing off ramp")
	}

	if jb.LegacyGasStationSidecarSpec.LookbackBlocks < int32(chain.Config().EVM().FinalityDepth()) {
		return nil, fmt.Errorf(
			"waitBlocks must be greater than or equal to chain's finality depth (%d), currently %d",
			chain.Config().EVM().FinalityDepth(), jb.LegacyGasStationSidecarSpec.LookbackBlocks)
	}

	orm := NewORM(d.db, d.logger, chain.Config().Database())

	var (
		mtlsCertificate string
		mtlsKey         string
	)

	if jb.LegacyGasStationSidecarSpec.ClientCertificate != nil && jb.LegacyGasStationSidecarSpec.ClientKey != nil {
		// temporary workaround until mtls auth config can be fetched inside nodes
		mtlsCertificate = *jb.LegacyGasStationSidecarSpec.ClientCertificate
		mtlsKey = *jb.LegacyGasStationSidecarSpec.ClientKey
	}

	su, err := legacygasstation.NewStatusUpdater(
		jb.LegacyGasStationSidecarSpec.StatusUpdateURL,
		mtlsCertificate,
		mtlsKey,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new status updater")
	}

	sidecar, err := legacygasstation.NewSidecar(
		log,
		NewLogPollerAdapter(chain.LogPoller()),
		forwarder,
		offramp,
		jb.LegacyGasStationSidecarSpec.CCIPChainSelector.ToInt().Uint64(),
		chain.Config().EVM().FinalityDepth(),
		uint32(jb.LegacyGasStationSidecarSpec.LookbackBlocks),
		orm,
		su,
	)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{&service{
		sidecar:    sidecar,
		pollPeriod: jb.LegacyGasStationSidecarSpec.PollPeriod,
		runTimeout: jb.LegacyGasStationSidecarSpec.RunTimeout,
		logger:     log,
		done:       make(chan struct{}),
	}}, nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *SidecarDelegate) AfterJobCreated(spec job.Job) {}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *SidecarDelegate) BeforeJobCreated(spec job.Job) {}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *SidecarDelegate) BeforeJobDeleted(spec job.Job) {}

// OnDeleteJob satisfies the job.Delegate interface.
func (d *SidecarDelegate) OnDeleteJob(spec job.Job, q pg.Queryer) error { return nil }

// service is a job.Service that runs the Gasless Transaction Sidecar every pollPeriod.
type service struct {
	utils.StartStopOnce
	sidecar    *legacygasstation.Sidecar
	done       chan struct{}
	pollPeriod time.Duration
	runTimeout time.Duration
	logger     logger.Logger
	parentCtx  context.Context
	cancel     context.CancelFunc
}

// Start the Gasless Transaction Sidecar, satisfying the job.Service interface.
func (s *service) Start(context.Context) error {
	return s.StartOnce("Gasless Transaction Sidecar", func() error {
		s.logger.Infow("Gasless Transaction Sidecar")
		ticker := time.NewTicker(utils.WithJitter(s.pollPeriod))
		s.parentCtx, s.cancel = context.WithCancel(context.Background())
		go func() {
			defer close(s.done)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.runSidecar()
				case <-s.parentCtx.Done():
					return
				}
			}
		}()
		return nil
	})
}

// Close the gasless transaction sidecar, satisfying the job.Service interface.
func (s *service) Close() error {
	return s.StopOnce("Gasless Transaction Sidecar", func() error {
		s.logger.Infow("Stopping Gasless Transaction Sidecar")
		s.cancel()
		<-s.done
		return nil
	})
}

func (s *service) runSidecar() {
	s.logger.Debugw("Running Gasless Transaction Sidecar")
	ctx, cancel := context.WithTimeout(s.parentCtx, s.runTimeout)
	defer cancel()
	err := s.sidecar.Run(ctx)
	if err == nil {
		s.logger.Debugw("Gasless Transaction Sidecar run completed successfully")
	} else {
		s.logger.Errorw("Gasless Transaction Sidecar run was at least partially unsuccessful",
			"error", err)
	}
}
