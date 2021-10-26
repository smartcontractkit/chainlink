package vrf

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"

	"github.com/theodesp/go-heaps/pairing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type Delegate struct {
	db   *gorm.DB
	pr   pipeline.Runner
	porm pipeline.ORM
	ks   keystore.Master
	cc   evm.ChainSet
	lggr logger.Logger
}

//go:generate mockery --name GethKeyStore --output mocks/ --case=underscore

type GethKeyStore interface {
	GetRoundRobinAddress(addresses ...common.Address) (common.Address, error)
}

type Config interface {
	MinIncomingConfirmations() uint32
	EvmGasLimitDefault() uint64
	EvmMaxGasPriceWei() *big.Int
	MinRequiredOutgoingConfirmations() uint64
}

func NewDelegate(
	db *gorm.DB,
	ks keystore.Master,
	pr pipeline.Runner,
	porm pipeline.ORM,
	chainSet evm.ChainSet,
	lggr logger.Logger) *Delegate {
	return &Delegate{
		db:   db,
		ks:   ks,
		pr:   pr,
		porm: porm,
		cc:   chainSet,
		lggr: lggr,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.Service, error) {
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
	abi := eth.MustGetABI(solidity_vrf_coordinator_interface.VRFCoordinatorABI)
	abiV2 := eth.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2ABI)
	l := d.lggr.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	)
	lV1 := l.Named("VRFListener")
	lV2 := l.Named("VRFListenerV2")

	vorm := keystore.NewVRFORM(d.db)
	for _, task := range pl.Tasks {
		if _, ok := task.(*pipeline.VRFTaskV2); ok {
			return []job.Service{&listenerV2{
				cfg:                chain.Config(),
				l:                  lV2,
				ethClient:          chain.Client(),
				logBroadcaster:     chain.LogBroadcaster(),
				db:                 d.db,
				abi:                abiV2,
				coordinator:        coordinatorV2,
				txm:                chain.TxManager(),
				pipelineRunner:     d.pr,
				vorm:               vorm,
				vrfks:              d.ks.VRF(),
				gethks:             d.ks.Eth(),
				pipelineORM:        d.porm,
				job:                jb,
				reqLogs:            utils.NewHighCapacityMailbox(),
				chStop:             make(chan struct{}),
				waitOnStop:         make(chan struct{}),
				respCount:          GetStartingResponseCountsV2(d.db, lV2),
				blockNumberToReqID: pairing.New(),
				reqAdded:           func() {},
			}}, nil
		}
		if _, ok := task.(*pipeline.VRFTask); ok {
			return []job.Service{&listenerV1{
				cfg:             chain.Config(),
				l:               lV1,
				headBroadcaster: chain.HeadBroadcaster(),
				logBroadcaster:  chain.LogBroadcaster(),
				db:              d.db,
				txm:             chain.TxManager(),
				abi:             abi,
				coordinator:     coordinator,
				pipelineRunner:  d.pr,
				vorm:            vorm,
				vrfks:           d.ks.VRF(),
				gethks:          d.ks.Eth(),
				pipelineORM:     d.porm,
				job:             jb,
				// Note the mailbox size effectively sets a limit on how many logs we can replay
				// in the event of a VRF outage.
				reqLogs:            utils.NewHighCapacityMailbox(),
				chStop:             make(chan struct{}),
				waitOnStop:         make(chan struct{}),
				newHead:            make(chan struct{}, 1),
				respCount:          getStartingResponseCounts(d.db, lV1),
				blockNumberToReqID: pairing.New(),
				reqAdded:           func() {},
			}}, nil
		}
	}
	return nil, errors.New("invalid job spec expected a vrf task")
}

func getStartingResponseCounts(db *gorm.DB, l logger.Logger) map[[32]byte]uint64 {
	respCounts := make(map[[32]byte]uint64)
	var counts []struct {
		RequestID string
		Count     int
	}
	// Allow any state, not just confirmed, on purpose.
	// We assume once a ethtx is queued it will go through.
	err := db.Raw(`SELECT meta->'RequestID' AS request_id, count(meta->'RequestID') as count
			FROM eth_txes
			WHERE meta->'RequestID' IS NOT NULL
		    GROUP BY meta->'RequestID'`).Scan(&counts).Error
	if err != nil {
		// Continue with an empty map, do not block job on this.
		l.Errorw("Unable to read previous fulfillments", "err", err)
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

func GetStartingResponseCountsV2(db *gorm.DB, l logger.Logger) map[string]uint64 {
	respCounts := make(map[string]uint64)
	var counts []struct {
		RequestID string
		Count     int
	}
	// Allow any state, not just confirmed, on purpose.
	// We assume once a ethtx is queued it will go through.
	err := db.Raw(`SELECT meta->'RequestID' AS request_id, count(meta->'RequestID') as count
			FROM eth_txes
			WHERE meta->'RequestID' IS NOT NULL
		    GROUP BY meta->'RequestID'`).Scan(&counts).Error
	if err != nil {
		// Continue with an empty map, do not block job on this.
		l.Errorw("Unable to read previous fulfillments", "err", err)
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
		bi := new(big.Int).SetBytes(b)
		respCounts[bi.String()] = uint64(c.Count)
	}
	return respCounts
}
