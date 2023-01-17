package web2

import (
	"sync"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/lottery_consumer"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

type Delegate struct {
	ks   keystore.Master
	cc   evm.ChainSet
	lggr logger.Logger
	db   *sqlx.DB
	cfg  pg.QConfig
}

func NewDelegate(
	db *sqlx.DB,
	ks keystore.Master,
	chainSet evm.ChainSet,
	lggr logger.Logger,
	cfg pg.QConfig) *Delegate {
	return &Delegate{
		ks:   ks,
		cc:   chainSet,
		lggr: lggr,
		db:   db,
		cfg:  cfg,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRFWeb2
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}
func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	chain, err := d.cc.Get(jb.VRFWeb2Spec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	lotteryAbi := evmtypes.MustGetABI(lottery_consumer.LotteryConsumerMetaData.ABI)

	lotteryConsumer, err := lottery_consumer.NewLotteryConsumer(
		jb.VRFWeb2Spec.LotteryConsumerAddress.Address(),
		chain.Client(),
	)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{
		&vrfServer{
			txManager:              chain.TxManager(),
			logBroadcaster:         chain.LogBroadcaster(),
			lotteryConsumerABI:     lotteryAbi,
			lotteryConsumerAddress: jb.VRFWeb2Spec.LotteryConsumerAddress.Address(),
			orm:                    newORM(chain.ID(), d.db, d.lggr, d.cfg),
			lggr:                   d.lggr.Named("VRFWeb2Server"),
			wg:                     &sync.WaitGroup{},
			j:                      jb,
			gethks:                 d.ks.Eth(),
			chainID:                chain.ID(),
			lotteryConsumer:        lotteryConsumer,
			q:                      pg.NewQ(d.db, d.lggr.Named("VRFWeb2Q"), d.cfg),
		},
	}, nil
}
