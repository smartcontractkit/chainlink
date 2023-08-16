package vrf

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
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
	q       pg.Q
	pr      pipeline.Runner
	porm    pipeline.ORM
	ks      keystore.Master
	cc      evm.ChainSet
	lggr    logger.Logger
	mailMon *utils.MailboxMonitor
}

func NewDelegate(
	db *sqlx.DB,
	ks keystore.Master,
	pr pipeline.Runner,
	porm pipeline.ORM,
	chainSet evm.ChainSet,
	lggr logger.Logger,
	cfg pg.QConfig,
	mailMon *utils.MailboxMonitor) *Delegate {
	return &Delegate{
		q:       pg.NewQ(db, lggr, cfg),
		ks:      ks,
		pr:      pr,
		porm:    porm,
		cc:      chainSet,
		lggr:    lggr,
		mailMon: mailMon,
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
func (d *Delegate) ServicesForSpec(jb job.Job, qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	if jb.VRFSpec == nil || jb.PipelineSpec == nil {
		return nil, errors.Errorf("vrf.Delegate expects a VRFSpec and PipelineSpec to be present, got %+v", jb)
	}
	pl, err := jb.PipelineSpec.Pipeline()
	if err != nil {
		return nil, err
	}
	chain, err := d.cc.Get(jb.VRFSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}
	chainId := chain.Client().ConfiguredChainID()
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}
	coordinatorV2Plus, err := vrf_coordinator_v2plus.NewVRFCoordinatorV2Plus(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
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

	l := d.lggr.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	)
	lV1 := l.Named("VRFListener")
	lV2 := l.Named("VRFListenerV2")
	lV2Plus := l.Named("VRFListenerV2Plus")

	for _, task := range pl.Tasks {
		if _, ok := task.(*pipeline.VRFTaskV2Plus); ok {
			if err := CheckFromAddressesExist(jb, d.ks.Eth()); err != nil {
				return nil, err
			}

			if !FromAddressMaxGasPricesAllEqual(jb, chain.Config().EVM().GasEstimator().PriceMaxKey) {
				return nil, errors.New("key-specific max gas prices of all fromAddresses are not equal, please set them to equal values")
			}

			if err := CheckFromAddressMaxGasPrices(jb, chain.Config().EVM().GasEstimator().PriceMaxKey); err != nil {
				return nil, err
			}
			if vrfOwner != nil {
				return nil, errors.New("VRF Owner is not supported for VRF V2 Plus")
			}
			linkEthFeedAddress, err := coordinatorV2Plus.LINKETHFEED(nil)
			if err != nil {
				return nil, errors.Wrap(err, "LINKETHFEED")
			}
			aggregator, err := aggregator_v3_interface.NewAggregatorV3Interface(linkEthFeedAddress, chain.Client())
			if err != nil {
				return nil, errors.Wrap(err, "NewAggregatorV3Interface")
			}

			return []job.ServiceCtx{v2.New(
				chain.Config().EVM(),
				chain.Config().EVM().GasEstimator(),
				lV2Plus,
				chain.Client(),
				chain.ID(),
				chain.LogBroadcaster(),
				d.q,
				v2.NewCoordinatorV2Plus(coordinatorV2Plus),
				batchCoordinatorV2,
				vrfOwner,
				aggregator,
				chain.TxManager(),
				d.pr,
				d.ks.Eth(),
				jb,
				d.mailMon,
				utils.NewHighCapacityMailbox[log.Broadcast](),
				func() {},
				GetStartingResponseCountsV2(d.q, lV2Plus, chainId.Uint64(), chain.Config().EVM().FinalityDepth()),
				chain.HeadBroadcaster(),
				vrfcommon.NewLogDeduper(int(chain.Config().EVM().FinalityDepth())))}, nil
		}
		if _, ok := task.(*pipeline.VRFTaskV2); ok {
			if err := CheckFromAddressesExist(jb, d.ks.Eth()); err != nil {
				return nil, err
			}

			if !FromAddressMaxGasPricesAllEqual(jb, chain.Config().EVM().GasEstimator().PriceMaxKey) {
				return nil, errors.New("key-specific max gas prices of all fromAddresses are not equal, please set them to equal values")
			}

			if err := CheckFromAddressMaxGasPrices(jb, chain.Config().EVM().GasEstimator().PriceMaxKey); err != nil {
				return nil, err
			}

			linkEthFeedAddress, err := coordinatorV2.LINKETHFEED(nil)
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
				chain.Client(),
				chain.ID(),
				chain.LogBroadcaster(),
				d.q,
				v2.NewCoordinatorV2(coordinatorV2),
				batchCoordinatorV2,
				vrfOwner,
				aggregator,
				chain.TxManager(),
				d.pr,
				d.ks.Eth(),
				jb,
				d.mailMon,
				utils.NewHighCapacityMailbox[log.Broadcast](),
				func() {},
				GetStartingResponseCountsV2(d.q, lV2, chainId.Uint64(), chain.Config().EVM().FinalityDepth()),
				chain.HeadBroadcaster(),
				vrfcommon.NewLogDeduper(int(chain.Config().EVM().FinalityDepth())))}, nil
		}
		if _, ok := task.(*pipeline.VRFTask); ok {
			return []job.ServiceCtx{&v1.Listener{
				Cfg:             chain.Config().EVM(),
				FeeCfg:          chain.Config().EVM().GasEstimator(),
				L:               logger.Sugared(lV1),
				HeadBroadcaster: chain.HeadBroadcaster(),
				LogBroadcaster:  chain.LogBroadcaster(),
				Q:               d.q,
				Txm:             chain.TxManager(),
				Coordinator:     coordinator,
				PipelineRunner:  d.pr,
				GethKs:          d.ks.Eth(),
				Job:             jb,
				MailMon:         d.mailMon,
				// Note the mailbox size effectively sets a limit on how many logs we can replay
				// in the event of a VRF outage.
				ReqLogs:            utils.NewHighCapacityMailbox[log.Broadcast](),
				ChStop:             make(chan struct{}),
				WaitOnStop:         make(chan struct{}),
				NewHead:            make(chan struct{}, 1),
				ResponseCount:      GetStartingResponseCountsV1(d.q, lV1, chainId.Uint64(), chain.Config().EVM().FinalityDepth()),
				BlockNumberToReqID: pairing.New(),
				ReqAdded:           func() {},
				Deduper:            vrfcommon.NewLogDeduper(int(chain.Config().EVM().FinalityDepth())),
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

func GetStartingResponseCountsV1(q pg.Q, l logger.Logger, chainID uint64, evmFinalityDepth uint32) map[[32]byte]uint64 {
	respCounts := map[[32]byte]uint64{}

	// Only check as far back as the evm finality depth for completed transactions.
	counts, err := getRespCounts(q, chainID, evmFinalityDepth)
	if err != nil {
		// Continue with an empty map, do not block job on this.
		l.Errorw("Unable to read previous confirmed fulfillments", "err", err)
		return respCounts
	}

	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			l.Errorw("Unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		var reqID [32]byte
		copy(reqID[:], b)
		respCounts[reqID] = uint64(c.Count)
	}

	return respCounts
}

func GetStartingResponseCountsV2(
	q pg.Q,
	l logger.Logger,
	chainID uint64,
	evmFinalityDepth uint32,
) map[string]uint64 {
	respCounts := map[string]uint64{}

	// Only check as far back as the evm finality depth for completed transactions.
	counts, err := getRespCounts(q, chainID, evmFinalityDepth)
	if err != nil {
		// Continue with an empty map, do not block job on this.
		l.Errorw("Unable to read previous confirmed fulfillments", "err", err)
		return respCounts
	}

	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			l.Errorw("Unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		bi := new(big.Int).SetBytes(b)
		respCounts[bi.String()] = uint64(c.Count)
	}
	return respCounts
}

func getRespCounts(q pg.Q, chainID uint64, evmFinalityDepth uint32) (
	[]struct {
		RequestID string
		Count     int
	},
	error,
) {
	counts := []struct {
		RequestID string
		Count     int
	}{}
	// This query should use the idx_eth_txes_state_from_address_evm_chain_id
	// index, since the quantity of unconfirmed/unstarted/in_progress transactions _should_ be small
	// relative to the rest of the data.
	unconfirmedQuery := `
SELECT meta->'RequestID' AS request_id, count(meta->'RequestID') AS count
FROM eth_txes et
WHERE et.meta->'RequestID' IS NOT NULL
AND et.state IN ('unconfirmed', 'unstarted', 'in_progress')
GROUP BY meta->'RequestID'
	`
	// Fetch completed transactions only as far back as the given cutoffBlockNumber. This avoids
	// a table scan of the eth_txes table, which could be large if it is unpruned.
	confirmedQuery := `
SELECT meta->'RequestID' AS request_id, count(meta->'RequestID') AS count
FROM eth_txes et JOIN eth_tx_attempts eta on et.id = eta.eth_tx_id
	join eth_receipts er on eta.hash = er.tx_hash
WHERE et.meta->'RequestID' is not null
AND er.block_number >= (SELECT number FROM evm_heads WHERE evm_chain_id = $1 ORDER BY number DESC LIMIT 1) - $2
GROUP BY meta->'RequestID'
	`
	query := unconfirmedQuery + "\nUNION ALL\n" + confirmedQuery
	err := q.Select(&counts, query, chainID, evmFinalityDepth)
	if err != nil {
		return nil, err
	}
	return counts, nil
}
