package blockheaderfeeder

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	v1 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_coordinator_interface"
	v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	v2plus "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

var _ job.ServiceCtx = &service{}

type Config interface {
	Feature() config.Feature
	Database() config.Database
}

type Delegate struct {
	cfg          Config
	logger       logger.Logger
	legacyChains legacyevm.LegacyChainContainer
	ks           keystore.Eth
}

func NewDelegate(
	cfg Config,
	logger logger.Logger,
	legacyChains legacyevm.LegacyChainContainer,
	ks keystore.Eth,
) *Delegate {
	return &Delegate{
		cfg:          cfg,
		logger:       logger,
		legacyChains: legacyChains,
		ks:           ks,
	}
}

// JobType satisfies the job.Delegate interface.
func (d *Delegate) JobType() job.Type {
	return job.BlockHeaderFeeder
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, jb job.Job) ([]job.ServiceCtx, error) {
	if jb.BlockHeaderFeederSpec == nil {
		return nil, errors.Errorf("Delegate expects a BlockHeaderFeederSpec to be present, got %+v", jb)
	}
	marshalledJob, err := json.MarshalIndent(jb.BlockHeaderFeederSpec, "", " ")
	if err != nil {
		return nil, err
	}
	d.logger.Debugw("Creating services for job spec", "job", string(marshalledJob))

	cid := jb.BlockHeaderFeederSpec.EVMChainID.ToInt()
	chain, err := d.legacyChains.Get(cid.String())
	if err != nil {
		return nil, fmt.Errorf(
			"getting chain ID %s: %w", cid, err)
	}

	if !d.cfg.Feature().LogPoller() {
		return nil, errors.New("log poller must be enabled to run blockheaderfeeder")
	}

	if jb.BlockHeaderFeederSpec.LookbackBlocks < int32(chain.Config().EVM().FinalityDepth()) {
		return nil, fmt.Errorf(
			"lookbackBlocks must be greater than or equal to chain's finality depth (%d), currently %d",
			chain.Config().EVM().FinalityDepth(), jb.BlockHeaderFeederSpec.LookbackBlocks)
	}

	ks := keys.NewChainStore(keystore.NewEthSigner(d.ks, cid), cid)
	enabled, err := ks.EnabledAddresses(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting sending keys")
	}
	if len(enabled) == 0 {
		return nil, fmt.Errorf("missing sending keys for chain ID: %v", chain.ID())
	}
	if err = CheckFromAddressesExist(jb, enabled); err != nil {
		return nil, err
	}
	fromAddresses := jb.BlockHeaderFeederSpec.FromAddresses

	bhs, err := blockhash_store.NewBlockhashStore(
		jb.BlockHeaderFeederSpec.BlockhashStoreAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "building BHS")
	}

	batchBlockhashStore, err := batch_blockhash_store.NewBatchBlockhashStore(
		jb.BlockHeaderFeederSpec.BatchBlockhashStoreAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "building batch BHS")
	}

	lp := chain.LogPoller()
	var coordinators []blockhashstore.Coordinator
	if jb.BlockHeaderFeederSpec.CoordinatorV1Address != nil {
		var c *v1.VRFCoordinator
		if c, err = v1.NewVRFCoordinator(
			jb.BlockHeaderFeederSpec.CoordinatorV1Address.Address(), chain.Client()); err != nil {
			return nil, errors.Wrap(err, "building V1 coordinator")
		}
		var coord *blockhashstore.V1Coordinator
		coord, err = blockhashstore.NewV1Coordinator(ctx, c, lp)
		if err != nil {
			return nil, errors.Wrap(err, "building V1 coordinator")
		}
		coordinators = append(coordinators, coord)
	}
	if jb.BlockHeaderFeederSpec.CoordinatorV2Address != nil {
		var c *v2.VRFCoordinatorV2
		if c, err = v2.NewVRFCoordinatorV2(
			jb.BlockHeaderFeederSpec.CoordinatorV2Address.Address(), chain.Client()); err != nil {
			return nil, errors.Wrap(err, "building V2 coordinator")
		}
		var coord *blockhashstore.V2Coordinator
		coord, err = blockhashstore.NewV2Coordinator(ctx, c, lp)
		if err != nil {
			return nil, errors.Wrap(err, "building V2 coordinator")
		}
		coordinators = append(coordinators, coord)
	}
	if jb.BlockHeaderFeederSpec.CoordinatorV2PlusAddress != nil {
		var c v2plus.IVRFCoordinatorV2PlusInternalInterface
		if c, err = v2plus.NewIVRFCoordinatorV2PlusInternal(
			jb.BlockHeaderFeederSpec.CoordinatorV2PlusAddress.Address(), chain.Client()); err != nil {
			return nil, errors.Wrap(err, "building V2 plus coordinator")
		}
		var coord *blockhashstore.V2PlusCoordinator
		coord, err = blockhashstore.NewV2PlusCoordinator(ctx, c, lp)
		if err != nil {
			return nil, errors.Wrap(err, "building V2 plus coordinator")
		}
		coordinators = append(coordinators, coord)
	}

	bpBHS, err := blockhashstore.NewBulletproofBHS(
		chain.Config().EVM().GasEstimator(),
		d.cfg.Database(),
		fromAddresses,
		chain.TxManager(),
		bhs,
		nil,
		ks,
	)
	if err != nil {
		return nil, errors.Wrap(err, "building bulletproof bhs")
	}

	batchBHS, err := blockhashstore.NewBatchBHS(
		chain.Config().EVM().GasEstimator(),
		chain.TxManager(),
		batchBlockhashStore,
	)
	if err != nil {
		return nil, errors.Wrap(err, "building batchBHS")
	}

	log := d.logger.Named("BlockHeaderFeeder").With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"bhsAddress", bhs.Address(),
		"batchBHSAddress", batchBlockhashStore.Address(),
	)

	blockHeaderProvider := NewGethBlockHeaderProvider(chain.Client())

	feeder := NewBlockHeaderFeeder(
		log,
		blockhashstore.NewMultiCoordinator(coordinators...),
		bpBHS,
		batchBHS,
		blockHeaderProvider,
		int(jb.BlockHeaderFeederSpec.WaitBlocks),
		int(jb.BlockHeaderFeederSpec.LookbackBlocks),
		func(ctx context.Context) (uint64, error) {
			head, err := chain.Client().HeadByNumber(ctx, nil)
			if err != nil {
				return 0, errors.Wrap(err, "getting chain head")
			}
			return uint64(head.Number), nil
		},
		ks,
		jb.BlockHeaderFeederSpec.GetBlockhashesBatchSize,
		jb.BlockHeaderFeederSpec.StoreBlockhashesBatchSize,
		fromAddresses,
	)

	services := []job.ServiceCtx{&service{
		feeder:     feeder,
		pollPeriod: jb.BlockHeaderFeederSpec.PollPeriod,
		runTimeout: jb.BlockHeaderFeederSpec.RunTimeout,
		logger:     log,
		done:       make(chan struct{}),
	}}

	return services, nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}

// BeforeJobDeleted satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// OnDeleteJob satisfies the job.Delegate interface.
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// service is a job.Service that runs the BHS feeder every pollPeriod.
type service struct {
	services.StateMachine
	feeder     *BlockHeaderFeeder
	done       chan struct{}
	pollPeriod time.Duration
	runTimeout time.Duration
	logger     logger.Logger
	stopCh     services.StopChan
}

// Start the BHS feeder service, satisfying the job.Service interface.
func (s *service) Start(context.Context) error {
	return s.StartOnce("Block Header Feeder Service", func() error {
		s.logger.Infow("Starting BlockHeaderFeeder")
		s.stopCh = make(chan struct{})
		go func() {
			defer close(s.done)
			ctx, cancel := s.stopCh.NewCtx()
			defer cancel()
			ticker := services.NewTicker(s.pollPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.runFeeder(ctx)
				case <-ctx.Done():
					return
				}
			}
		}()
		return nil
	})
}

// Close the BHS feeder service, satisfying the job.Service interface.
func (s *service) Close() error {
	return s.StopOnce("Block Header Feeder Service", func() error {
		s.logger.Infow("Stopping BlockHeaderFeeder")
		close(s.stopCh)
		<-s.done
		return nil
	})
}

func (s *service) runFeeder(ctx context.Context) {
	s.logger.Debugw("Running BlockHeaderFeeder")
	ctx, cancel := context.WithTimeout(ctx, s.runTimeout)
	defer cancel()
	err := s.feeder.Run(ctx)
	if err == nil {
		s.logger.Debugw("BlockHeaderFeeder run completed successfully")
	} else {
		s.logger.Errorw("BlockHeaderFeeder run was at least partially unsuccessful",
			"err", err)
	}
}

// CheckFromAddressesExist returns an error if and only if one of the addresses
// in the BlockHeaderFeeder spec's fromAddresses field does not exist in the keystore.
func CheckFromAddressesExist(jb job.Job, enabled []common.Address) (err error) {
	for _, a := range jb.BlockHeaderFeederSpec.FromAddresses {
		if !slices.Contains(enabled, a.Address()) {
			err = multierr.Append(err, fmt.Errorf("address not enabled: %s", a.Hex()))
		}
	}
	return
}
