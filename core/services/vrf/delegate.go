package vrf

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"

	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type Delegate struct {
	cfg  Config
	db   *gorm.DB
	txm  bulletprooftxmanager.TxManager
	pr   pipeline.Runner
	porm pipeline.ORM
	ks   keystore.Master
	ec   eth.Client
	hb   httypes.HeadBroadcasterRegistry
	lb   log.Broadcaster
}

//go:generate mockery --name GethKeyStore --output mocks/ --case=underscore

type GethKeyStore interface {
	GetRoundRobinAddress(addresses ...common.Address) (common.Address, error)
}

type Config interface {
	MinIncomingConfirmations() uint32
	EvmGasLimitDefault() uint64
}

func NewDelegate(
	db *gorm.DB,
	txm bulletprooftxmanager.TxManager,
	ks keystore.Master,
	pr pipeline.Runner,
	porm pipeline.ORM,
	lb log.Broadcaster,
	hb httypes.HeadBroadcasterRegistry,
	ec eth.Client,
	cfg Config) *Delegate {
	return &Delegate{
		cfg:  cfg,
		db:   db,
		txm:  txm,
		ks:   ks,
		pr:   pr,
		porm: porm,
		hb:   hb,
		lb:   lb,
		ec:   ec,
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
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), d.ec)
	if err != nil {
		return nil, err
	}
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(jb.VRFSpec.CoordinatorAddress.Address(), d.ec)
	if err != nil {
		return nil, err
	}
	abi := eth.MustGetABI(solidity_vrf_coordinator_interface.VRFCoordinatorABI)
	abiV2 := eth.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2ABI)
	l := logger.CreateLogger(logger.Default.SugaredLogger.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	))

	vorm := keystore.NewVRFORM(d.db)
	for _, task := range pl.Tasks {
		if _, ok := task.(*pipeline.VRFTaskV2); ok {
			return []job.Service{&listenerV2{
				cfg:                d.cfg,
				l:                  *l,
				ethClient:          d.ec,
				logBroadcaster:     d.lb,
				headBroadcaster:    d.hb,
				db:                 d.db,
				abi:                abiV2,
				coordinator:        coordinatorV2,
				txm:                d.txm,
				pipelineRunner:     d.pr,
				vorm:               vorm,
				vrfks:              d.ks.VRF(),
				gethks:             d.ks.Eth(),
				pipelineORM:        d.porm,
				job:                jb,
				reqLogs:            utils.NewMailbox(100000),
				chStop:             make(chan struct{}),
				waitOnStop:         make(chan struct{}),
				newHead:            make(chan struct{}, 1),
				respCount:          GetStartingResponseCountsV2(d.db, l),
				blockNumberToReqID: pairing.New(),
				reqAdded:           func() {},
			}}, nil
		}
		if _, ok := task.(*pipeline.VRFTask); ok {
			return []job.Service{&listenerV1{
				cfg:             d.cfg,
				l:               *l,
				headBroadcaster: d.hb,
				logBroadcaster:  d.lb,
				db:              d.db,
				txm:             d.txm,
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
				reqLogs:            utils.NewMailbox(100000),
				chStop:             make(chan struct{}),
				waitOnStop:         make(chan struct{}),
				newHead:            make(chan struct{}, 1),
				respCount:          getStartingResponseCounts(d.db, l),
				blockNumberToReqID: pairing.New(),
				reqAdded:           func() {},
			}}, nil
		}
	}
	return nil, errors.New("invalid job spec expected a vrf task")
}

func getStartingResponseCounts(db *gorm.DB, l *logger.Logger) map[[32]byte]uint64 {
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
		l.Errorw("VRFListener: unable to read previous fulfillments", "err", err)
		return respCounts
	}
	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			l.Errorw("VRFListener: unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		var reqID [32]byte
		copy(reqID[:], b)
		respCounts[reqID] = uint64(c.Count)
	}
	return respCounts
}

func GetStartingResponseCountsV2(db *gorm.DB, l *logger.Logger) map[string]uint64 {
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
		l.Errorw("VRFListenerV2: unable to read previous fulfillments", "err", err)
		return respCounts
	}
	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			l.Errorw("VRFListenerV2: unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		var reqID [32]byte
		copy(reqID[:], b)
		bi := new(big.Int).SetBytes(b)
		respCounts[bi.String()] = uint64(c.Count)
	}
	return respCounts
}
