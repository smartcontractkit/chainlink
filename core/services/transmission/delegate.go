package transmission

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/transmission/handler"
)

type Delegate struct {
	lggr     logger.Logger
	chainSet evm.ChainSet
	q        pg.Q
	txmORM   txmgr.ORM
}

func NewDelegate(lggr logger.Logger, chainSet evm.ChainSet, q pg.Q, txmORM txmgr.ORM) *Delegate {
	return &Delegate{
		lggr:     lggr,
		chainSet: chainSet,
		q:        q,
		txmORM:   txmORM,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}
func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	if jb.TransmissionSpec == nil {
		return nil, errors.Errorf("transmission.Delegate expects a Transmission Spec, got %+v", jb)
	}

	chain, err := d.chainSet.Get(jb.TransmissionSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	handler := handler.NewHandler(d.lggr, chain, jb.TransmissionSpec.FromAddresses, d.q, d.txmORM)
	server := NewServer(*handler, jb.TransmissionSpec.RPCPort, d.lggr)
	return []job.ServiceCtx{
		server,
	}, nil
}
