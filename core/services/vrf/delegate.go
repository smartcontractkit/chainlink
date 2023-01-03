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

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
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

//go:generate mockery --quiet --name GethKeyStore --output ./mocks/ --case=underscore
type GethKeyStore interface {
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (common.Address, error)
}

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore
type Config interface {
	EvmFinalityDepth() uint32
	EvmGasLimitDefault() uint32
	EvmGasLimitVRFJobType() *uint32
	KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei
	MinIncomingConfirmations() uint32
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

func (d *Delegate) BeforeJobCreated(spec job.Job) {}
func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
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
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
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

	l := d.lggr.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	)
	lV1 := l.Named("VRFListener")
	lV2 := l.Named("VRFListenerV2")

	for _, task := range pl.Tasks {
		if _, ok := task.(*pipeline.VRFTaskV2); ok {
			if err := CheckFromAddressMaxGasPrices(jb, chain.Config()); err != nil {
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

			return []job.ServiceCtx{newListenerV2(
				chain.Config(),
				lV2,
				chain.Client(),
				chain.ID(),
				chain.LogBroadcaster(),
				d.q,
				coordinatorV2,
				batchCoordinatorV2,
				aggregator,
				chain.TxManager(),
				d.pr,
				d.ks.Eth(),
				jb,
				d.mailMon,
				utils.NewHighCapacityMailbox[log.Broadcast](),
				func() {},
				GetStartingResponseCountsV2(d.q, lV2, chain.Client().ChainID().Uint64(), chain.Config().EvmFinalityDepth()),
				chain.HeadBroadcaster(),
				newLogDeduper(int(chain.Config().EvmFinalityDepth())))}, nil
		}
		if _, ok := task.(*pipeline.VRFTask); ok {
			return []job.ServiceCtx{&listenerV1{
				cfg:             chain.Config(),
				l:               logger.Sugared(lV1),
				headBroadcaster: chain.HeadBroadcaster(),
				logBroadcaster:  chain.LogBroadcaster(),
				q:               d.q,
				txm:             chain.TxManager(),
				coordinator:     coordinator,
				pipelineRunner:  d.pr,
				gethks:          d.ks.Eth(),
				job:             jb,
				mailMon:         d.mailMon,
				// Note the mailbox size effectively sets a limit on how many logs we can replay
				// in the event of a VRF outage.
				reqLogs:            utils.NewHighCapacityMailbox[log.Broadcast](),
				chStop:             make(chan struct{}),
				waitOnStop:         make(chan struct{}),
				newHead:            make(chan struct{}, 1),
				respCount:          GetStartingResponseCountsV1(d.q, lV1, chain.Client().ChainID().Uint64(), chain.Config().EvmFinalityDepth()),
				blockNumberToReqID: pairing.New(),
				reqAdded:           func() {},
				deduper:            newLogDeduper(int(chain.Config().EvmFinalityDepth())),
			}}, nil
		}
	}
	return nil, errors.New("invalid job spec expected a vrf task")
}

// CheckFromAddressMaxGasPrices checks if the provided gas price in the job spec gas lane parameter
// matches what is set for the  provided from addresses.
// If they don't match, this is a configuration error. An error is returned with all the keys that do
// not match the provided gas lane price.
func CheckFromAddressMaxGasPrices(jb job.Job, cfg Config) (err error) {
	if jb.VRFSpec.GasLanePrice != nil {
		for _, a := range jb.VRFSpec.FromAddresses {
			if keySpecific := cfg.KeySpecificMaxGasPriceWei(a.Address()); !keySpecific.Equal(jb.VRFSpec.GasLanePrice) {
				err = multierr.Append(err,
					fmt.Errorf(
						"key-specific max gas price of from address %s (%s) does not match gasLanePriceGWei (%s) specified in job spec",
						a.Hex(), keySpecific.String(), jb.VRFSpec.GasLanePrice.String()))
			}
		}
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
