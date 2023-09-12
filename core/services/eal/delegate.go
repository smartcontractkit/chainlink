package eal

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/eal/generated/forwarder"
	eallib "github.com/smartcontractkit/capital-markets-projects/lib/services/eal"
	ealtypes "github.com/smartcontractkit/capital-markets-projects/lib/services/eal/types"
	"github.com/smartcontractkit/capital-markets-projects/lib/web/jsonrpc"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"go.uber.org/multierr"
)

type (
	Delegate struct {
		lggr logger.Logger
		cc   evm.LegacyChainContainer
		ks   keystore.Eth
		q    pg.Q
		db   *sqlx.DB
		rr   *eallib.RequestRouter
	}

	RequestRouter interface {
		SendTransaction(*gin.Context, ealtypes.SendTransactionRequest) (*ealtypes.SendTransactionResponse, *jsonrpc.Error)
	}
)

func NewDelegate(lggr logger.Logger, cc evm.LegacyChainContainer, ks keystore.Eth) *Delegate {
	return &Delegate{
		lggr: lggr,
		cc:   cc,
		ks:   ks,
		rr:   eallib.NewRequestRouter(lggr),
	}
}

func (d *Delegate) JobType() job.Type {
	return job.EAL
}

func (d *Delegate) BeforeJobCreated(spec job.Job)                {}
func (d *Delegate) AfterJobCreated(spec job.Job)                 {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)                {}
func (d *Delegate) OnDeleteJob(spec job.Job, q pg.Queryer) error { return nil }

func (d *Delegate) RequestRouter() RequestRouter {
	return d.rr
}

func (d *Delegate) ServicesForSpec(jb job.Job, qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	if jb.EALSpec == nil {
		return nil, errors.Errorf("ServicesForSpec expects a EALSpec, got %+v", jb)
	}
	service := &ealService{
		spec: jb,
		rr:   d.rr,
		cc:   d.cc,
		ks:   d.ks,
		q:    d.q,
		db:   d.db,
		lggr: d.lggr,
	}
	return []job.ServiceCtx{service}, nil
}

type ealService struct {
	spec job.Job
	rr   *eallib.RequestRouter
	cc   evm.LegacyChainContainer
	ks   keystore.Eth
	q    pg.Q
	db   *sqlx.DB
	lggr logger.Logger
}

// Start starts ealService.
func (s *ealService) Start(context.Context) error {
	l := s.lggr.Named("EAL").With(
		"jobID", s.spec.ID,
		"externalJobID", s.spec.ExternalJobID,
		"chainID", s.spec.EALSpec.EVMChainID.ToInt().Uint64(),
		"ccipChainSelector", s.spec.EALSpec.CCIPChainSelector.ToInt().Uint64(),
	)
	chain, err := s.cc.Get(s.spec.EALSpec.EVMChainID.String())
	if err != nil {
		return err
	}
	forwarder, err := forwarder.NewForwarder(s.spec.EALSpec.ForwarderAddress.Address(), chain.Client())
	if err != nil {
		return errors.Wrap(err, "initializing forwarder")
	}
	if err = checkFromAddressesExist(s.spec, s.ks); err != nil {
		return err
	}

	client, err := NewBlockchainClient(
		l,
		chain.TxManager(),
		s.ks,
		*s.spec.EALSpec,
		chain.ID().Uint64(),
	)

	if err != nil {
		return err
	}

	reqHandler, err := eallib.NewRequestHandler(
		l,
		forwarder,
		chain.ID().Uint64(),
		s.spec.EALSpec.CCIPChainSelector.ToInt().Uint64(),
		client,
	)
	if err != nil {
		return err
	}
	err = s.rr.RegisterHandler(reqHandler)
	if err != nil {
		return err
	}
	return err
}

func (s *ealService) Close() error {
	s.rr.DeregisterHandler(s.spec.EALSpec.CCIPChainSelector.ToInt())
	return nil
}

// CheckFromAddressesExist returns an error if and only if one of the addresses
// in the EALSpec spec's fromAddresses field does not exist in the keystore.
func checkFromAddressesExist(jb job.Job, gethks keystore.Eth) (err error) {
	for _, a := range jb.EALSpec.FromAddresses {
		_, err2 := gethks.Get(a.Hex())
		err = multierr.Append(err, err2)
	}
	return
}
