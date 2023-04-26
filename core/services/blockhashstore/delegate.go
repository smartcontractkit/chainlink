package blockhashstore

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	v1 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	v2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ job.ServiceCtx = &service{}

// Delegate creates BlockhashStore feeder jobs.
type Delegate struct {
	logger logger.Logger
	chains evm.ChainSet
	ks     keystore.Eth
}

// NewDelegate creates a new Delegate.
func NewDelegate(
	logger logger.Logger,
	chains evm.ChainSet,
	ks keystore.Eth,
) *Delegate {
	return &Delegate{
		logger: logger,
		chains: chains,
		ks:     ks,
	}
}

// JobType satisfies the job.Delegate interface.
func (d *Delegate) JobType() job.Type {
	return job.BlockhashStore
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	if jb.BlockhashStoreSpec == nil {
		return nil, errors.Errorf(
			"blockhashstore.Delegate expects a BlockhashStoreSpec to be present, got %+v", jb)
	}

	chain, err := d.chains.Get(jb.BlockhashStoreSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, fmt.Errorf(
			"getting chain ID %d: %w", jb.BlockhashStoreSpec.EVMChainID.ToInt(), err)
	}

	if !chain.Config().FeatureLogPoller() {
		return nil, errors.New("log poller must be enabled to run blockhashstore")
	}

	if jb.BlockhashStoreSpec.WaitBlocks < int32(chain.Config().EvmFinalityDepth()) {
		return nil, fmt.Errorf(
			"waitBlocks must be greater than or equal to chain's finality depth (%d), currently %d",
			chain.Config().EvmFinalityDepth(), jb.BlockhashStoreSpec.WaitBlocks)
	}

	keys, err := d.ks.EnabledKeysForChain(chain.ID())
	if err != nil {
		return nil, errors.Wrap(err, "getting sending keys")
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("missing sending keys for chain ID: %v", chain.ID())
	}
	fromAddresses := []ethkey.EIP55Address{keys[0].EIP55Address}
	if jb.BlockhashStoreSpec.FromAddresses != nil {
		fromAddresses = jb.BlockhashStoreSpec.FromAddresses
	}

	bhs, err := blockhash_store.NewBlockhashStore(
		jb.BlockhashStoreSpec.BlockhashStoreAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "building BHS")
	}

	lp := chain.LogPoller()
	var coordinators []Coordinator
	if jb.BlockhashStoreSpec.CoordinatorV1Address != nil {
		var c *v1.VRFCoordinator
		if c, err = v1.NewVRFCoordinator(
			jb.BlockhashStoreSpec.CoordinatorV1Address.Address(), chain.Client()); err != nil {

			return nil, errors.Wrap(err, "building V1 coordinator")
		}

		var coord *V1Coordinator
		coord, err = NewV1Coordinator(c, lp)
		if err != nil {
			return nil, errors.Wrap(err, "building V1 coordinator")
		}
		coordinators = append(coordinators, coord)
	}
	if jb.BlockhashStoreSpec.CoordinatorV2Address != nil {
		var c *v2.VRFCoordinatorV2
		if c, err = v2.NewVRFCoordinatorV2(
			jb.BlockhashStoreSpec.CoordinatorV2Address.Address(), chain.Client()); err != nil {

			return nil, errors.Wrap(err, "building V2 coordinator")
		}

		var coord *V2Coordinator
		coord, err = NewV2Coordinator(c, lp)
		if err != nil {
			return nil, errors.Wrap(err, "building V2 coordinator")
		}
		coordinators = append(coordinators, coord)
	}

	bpBHS, err := NewBulletproofBHS(chain.Config(), fromAddresses, chain.TxManager(), bhs, chain.ID(), d.ks)
	if err != nil {
		return nil, errors.Wrap(err, "building bulletproof bhs")
	}

	log := d.logger.Named("BHS Feeder").With("jobID", jb.ID, "externalJobID", jb.ExternalJobID)
	feeder := NewFeeder(
		log,
		NewMultiCoordinator(coordinators...),
		bpBHS,
		int(jb.BlockhashStoreSpec.WaitBlocks),
		int(jb.BlockhashStoreSpec.LookbackBlocks),
		func(ctx context.Context) (uint64, error) {
			head, err := chain.Client().HeadByNumber(ctx, nil)
			if err != nil {
				return 0, errors.Wrap(err, "getting chain head")
			}
			return uint64(head.Number), nil
		})

	return []job.ServiceCtx{&service{
		feeder:     feeder,
		pollPeriod: jb.BlockhashStoreSpec.PollPeriod,
		runTimeout: jb.BlockhashStoreSpec.RunTimeout,
		logger:     log,
		done:       make(chan struct{}),
	}}, nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) AfterJobCreated(spec job.Job) {}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobCreated(spec job.Job) {}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// OnDeleteJob satisfies the job.Delegate interface.
func (d *Delegate) OnDeleteJob(spec job.Job, q pg.Queryer) error { return nil }

// service is a job.Service that runs the BHS feeder every pollPeriod.
type service struct {
	utils.StartStopOnce
	feeder     *Feeder
	done       chan struct{}
	pollPeriod time.Duration
	runTimeout time.Duration
	logger     logger.Logger
	parentCtx  context.Context
	cancel     context.CancelFunc
}

// Start the BHS feeder service, satisfying the job.Service interface.
func (s *service) Start(context.Context) error {
	return s.StartOnce("BHS Feeder Service", func() error {
		s.logger.Infow("Starting BHS feeder")
		ticker := time.NewTicker(utils.WithJitter(s.pollPeriod))
		s.parentCtx, s.cancel = context.WithCancel(context.Background())
		go func() {
			defer close(s.done)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.runFeeder()
				case <-s.parentCtx.Done():
					return
				}
			}
		}()
		return nil
	})
}

// Close the BHS feeder service, satisfying the job.Service interface.
func (s *service) Close() error {
	return s.StopOnce("BHS Feeder Service", func() error {
		s.logger.Infow("Stopping BHS feeder")
		s.cancel()
		<-s.done
		return nil
	})
}

func (s *service) runFeeder() {
	s.logger.Debugw("Running BHS feeder")
	ctx, cancel := context.WithTimeout(s.parentCtx, s.runTimeout)
	defer cancel()
	err := s.feeder.Run(ctx)
	if err == nil {
		s.logger.Debugw("BHS feeder run completed successfully")
	} else {
		s.logger.Errorw("BHS feeder run was at least partially unsuccessful",
			"error", err)
	}
}
