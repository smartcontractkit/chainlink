package vrf

import (
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/theodesp/go-heaps/pairing"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	v1 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v1"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Delegate struct {
	q            pg.Q
	pr           pipeline.Runner
	porm         pipeline.ORM
	ks           keystore.Master
	legacyChains evm.LegacyChainContainer
	lggr         logger.Logger
	mailMon      *utils.MailboxMonitor
}

func NewDelegate(
	db *sqlx.DB,
	ks keystore.Master,
	pr pipeline.Runner,
	porm pipeline.ORM,
	legacyChains evm.LegacyChainContainer,
	lggr logger.Logger,
	cfg pg.QConfig,
	mailMon *utils.MailboxMonitor) *Delegate {
	return &Delegate{
		q:            pg.NewQ(db, lggr, cfg),
		ks:           ks,
		pr:           pr,
		porm:         porm,
		legacyChains: legacyChains,
		lggr:         lggr.Named("VRF"),
		mailMon:      mailMon,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) BeforeJobCreated(spec job.Job)                {}
func (d *Delegate) AfterJobCreated(spec job.Job)                 {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)                {}
func (d *Delegate) OnDeleteJob(spec job.Job, q pg.Queryer) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	if jb.VRFSpec == nil || jb.PipelineSpec == nil {
		return nil, errors.Errorf("vrf.Delegate expects a VRFSpec and PipelineSpec to be present, got %+v", jb)
	}
	pl, err := jb.PipelineSpec.Pipeline()
	if err != nil {
		return nil, err
	}
	chain, err := d.legacyChains.Get(jb.VRFSpec.EVMChainID.String())
	if err != nil {
		return nil, err
	}
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}
	coordinatorV2Plus, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}

	// If the batch coordinator address is not provided, we will fall back to non-batched
	var batchCoordinatorV2 *batch_vrf_coordinator_v2.BatchVRFCoordinatorV2
	if jb.VRFSpec.BatchCoordinatorAddress != nil {
		batchCoordinatorV2, err = batch_vrf_coordinator_v2.NewBatchVRFCoordinatorV2(
			jb.VRFSpec.BatchCoordinatorAddress.Address(), chain.Client())
		if err != nil {
			return nil, errors.Wrap(err, "create batch coordinator wrapper")
		}
	}

	var vrfOwner *vrf_owner.VRFOwner
	if jb.VRFSpec.VRFOwnerAddress != nil {
		vrfOwner, err = vrf_owner.NewVRFOwner(
			jb.VRFSpec.VRFOwnerAddress.Address(), chain.Client(),
		)
		if err != nil {
			return nil, errors.Wrap(err, "create vrf owner wrapper")
		}
	}

	l := d.lggr.Named(jb.ExternalJobID.String()).With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	)
	lV1 := l.Named("VRFListener")
	lV2 := l.Named("VRFListenerV2")
	lV2Plus := l.Named("VRFListenerV2Plus")

	for _, task := range pl.Tasks {
		if _, ok := task.(*pipeline.VRFTaskV2Plus); ok {
			if err2 := CheckFromAddressesExist(jb, d.ks.Eth()); err != nil {
				return nil, err2
			}

			if !FromAddressMaxGasPricesAllEqual(jb, chain.Config().EVM().GasEstimator().PriceMaxKey) {
				return nil, errors.New("key-specific max gas prices of all fromAddresses are not equal, please set them to equal values")
			}

			if err2 := CheckFromAddressMaxGasPrices(jb, chain.Config().EVM().GasEstimator().PriceMaxKey); err != nil {
				return nil, err2
			}
			if vrfOwner != nil {
				return nil, errors.New("VRF Owner is not supported for VRF V2 Plus")
			}

			// Get the LINKNATIVEFEED address with retries
			// This is needed because the RPC endpoint may be down so we need to
			// switch over to another one.
			var linkNativeFeedAddress common.Address
			err = retry.Do(func() error {
				linkNativeFeedAddress, err = coordinatorV2Plus.LINKNATIVEFEED(nil)
				return err
			}, retry.Attempts(10), retry.Delay(500*time.Millisecond))
			if err != nil {
				return nil, errors.Wrap(err, "can't call LINKNATIVEFEED")
			}

			aggregator, err2 := aggregator_v3_interface.NewAggregatorV3Interface(linkNativeFeedAddress, chain.Client())
			if err2 != nil {
				return nil, errors.Wrap(err2, "NewAggregatorV3Interface")
			}

			return []job.ServiceCtx{v2.New(
				chain.Config().EVM(),
				chain.Config().EVM().GasEstimator(),
				lV2Plus,
				chain,
				chain.ID(),
				d.q,
				v2.NewCoordinatorV2_5(coordinatorV2Plus),
				batchCoordinatorV2,
				vrfOwner,
				aggregator,
				d.pr,
				d.ks.Eth(),
				jb,
				d.mailMon,
				utils.NewHighCapacityMailbox[log.Broadcast](),
				func() {},
				vrfcommon.NewLogDeduper(int(chain.Config().EVM().FinalityDepth())))}, nil
		}
		if _, ok := task.(*pipeline.VRFTaskV2); ok {
			if err2 := CheckFromAddressesExist(jb, d.ks.Eth()); err != nil {
				return nil, err2
			}

			if !FromAddressMaxGasPricesAllEqual(jb, chain.Config().EVM().GasEstimator().PriceMaxKey) {
				return nil, errors.New("key-specific max gas prices of all fromAddresses are not equal, please set them to equal values")
			}

			if err2 := CheckFromAddressMaxGasPrices(jb, chain.Config().EVM().GasEstimator().PriceMaxKey); err != nil {
				return nil, err2
			}

			// Get the LINKETHFEED address with retries
			// This is needed because the RPC endpoint may be down so we need to
			// switch over to another one.
			var linkEthFeedAddress common.Address
			err = retry.Do(func() error {
				linkEthFeedAddress, err = coordinatorV2.LINKETHFEED(nil)
				return err
			}, retry.Attempts(10), retry.Delay(500*time.Millisecond))
			if err != nil {
				return nil, errors.Wrap(err, "LINKETHFEED")
			}
			aggregator, err := aggregator_v3_interface.NewAggregatorV3Interface(linkEthFeedAddress, chain.Client())
			if err != nil {
				return nil, errors.Wrap(err, "NewAggregatorV3Interface")
			}
			if vrfOwner == nil {
				lV2.Infow("Running without VRFOwnerAddress set on the spec")
			}

			return []job.ServiceCtx{v2.New(
				chain.Config().EVM(),
				chain.Config().EVM().GasEstimator(),
				lV2,
				chain,
				chain.ID(),
				d.q,
				v2.NewCoordinatorV2(coordinatorV2),
				batchCoordinatorV2,
				vrfOwner,
				aggregator,
				d.pr,
				d.ks.Eth(),
				jb,
				d.mailMon,
				utils.NewHighCapacityMailbox[log.Broadcast](),
				func() {},
				vrfcommon.NewLogDeduper(int(chain.Config().EVM().FinalityDepth())))}, nil
		}
		if _, ok := task.(*pipeline.VRFTask); ok {
			return []job.ServiceCtx{&v1.Listener{
				Cfg:            chain.Config().EVM(),
				FeeCfg:         chain.Config().EVM().GasEstimator(),
				L:              logger.Sugared(lV1),
				Q:              d.q,
				Coordinator:    coordinator,
				PipelineRunner: d.pr,
				GethKs:         d.ks.Eth(),
				Job:            jb,
				MailMon:        d.mailMon,
				// Note the mailbox size effectively sets a limit on how many logs we can replay
				// in the event of a VRF outage.
				ReqLogs:            utils.NewHighCapacityMailbox[log.Broadcast](),
				ChStop:             make(chan struct{}),
				WaitOnStop:         make(chan struct{}),
				NewHead:            make(chan struct{}, 1),
				BlockNumberToReqID: pairing.New(),
				ReqAdded:           func() {},
				Deduper:            vrfcommon.NewLogDeduper(int(chain.Config().EVM().FinalityDepth())),
				Chain:              chain,
			}}, nil
		}
	}
	return nil, errors.New("invalid job spec expected a vrf task")
}

// CheckFromAddressesExist returns an error if and only if one of the addresses
// in the VRF spec's fromAddresses field does not exist in the keystore.
func CheckFromAddressesExist(jb job.Job, gethks keystore.Eth) (err error) {
	for _, a := range jb.VRFSpec.FromAddresses {
		_, err2 := gethks.Get(a.Hex())
		err = multierr.Append(err, err2)
	}
	return
}

// CheckFromAddressMaxGasPrices checks if the provided gas price in the job spec gas lane parameter
// matches what is set for the  provided from addresses.
// If they don't match, this is a configuration error. An error is returned with all the keys that do
// not match the provided gas lane price.
func CheckFromAddressMaxGasPrices(jb job.Job, keySpecificMaxGas keySpecificMaxGasFn) (err error) {
	if jb.VRFSpec.GasLanePrice != nil {
		for _, a := range jb.VRFSpec.FromAddresses {
			if keySpecific := keySpecificMaxGas(a.Address()); !keySpecific.Equal(jb.VRFSpec.GasLanePrice) {
				err = multierr.Append(err,
					fmt.Errorf(
						"key-specific max gas price of from address %s (%s) does not match gasLanePriceGWei (%s) specified in job spec",
						a.Hex(), keySpecific.String(), jb.VRFSpec.GasLanePrice.String()))
			}
		}
	}
	return
}

type keySpecificMaxGasFn func(common.Address) *assets.Wei

// FromAddressMaxGasPricesAllEqual returns true if and only if all the specified from
// addresses in the fromAddresses field of the VRF v2 job have the same key-specific max
// gas price.
func FromAddressMaxGasPricesAllEqual(jb job.Job, keySpecificMaxGasPriceWei keySpecificMaxGasFn) (allEqual bool) {
	allEqual = true
	for i := range jb.VRFSpec.FromAddresses {
		allEqual = allEqual && keySpecificMaxGasPriceWei(jb.VRFSpec.FromAddresses[i].Address()).Equal(
			keySpecificMaxGasPriceWei(jb.VRFSpec.FromAddresses[0].Address()),
		)
	}
	return
}
