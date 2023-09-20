package legacygasstation

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"

	forwarder "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/legacygasstation/generated/legacy_gas_station_forwarder"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	"github.com/smartcontractkit/capital-markets-projects/lib/web/jsonrpc"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/eal"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type (
	Delegate struct {
		lggr logger.Logger
		cc   evm.LegacyChainContainer
		ks   keystore.Eth
		q    pg.Q
		db   *sqlx.DB
		rr   *legacygasstation.RequestRouter
	}

	RequestRouter interface {
		SendTransaction(*gin.Context, types.SendTransactionRequest) (*types.SendTransactionResponse, *jsonrpc.Error)
	}
)

func NewServerDelegate(lggr logger.Logger, cc evm.LegacyChainContainer, ks keystore.Eth, db *sqlx.DB, cfg pg.QConfig) *Delegate {
	return &Delegate{
		lggr: lggr,
		cc:   cc,
		ks:   ks,
		q:    pg.NewQ(db, lggr, cfg),
		db:   db,
		rr:   legacygasstation.NewRequestRouter(lggr),
	}
}

func (d *Delegate) JobType() job.Type {
	return job.LegacyGasStationServer
}

func (d *Delegate) BeforeJobCreated(spec job.Job)                {}
func (d *Delegate) AfterJobCreated(spec job.Job)                 {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)                {}
func (d *Delegate) OnDeleteJob(spec job.Job, q pg.Queryer) error { return nil }

func (d *Delegate) RequestRouter() RequestRouter {
	return d.rr
}

func (d *Delegate) ServicesForSpec(jb job.Job, qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	if jb.LegacyGasStationServerSpec == nil {
		return nil, errors.Errorf("ServicesForSpec expects a LegacyGasStationServerSpec, got %+v", jb)
	}
	service := &gasStationService{
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

type gasStationService struct {
	spec job.Job
	rr   *legacygasstation.RequestRouter
	cc   evm.LegacyChainContainer
	ks   keystore.Eth
	q    pg.Q
	db   *sqlx.DB
	lggr logger.Logger
}

// Start starts gasStationService.
func (s *gasStationService) Start(context.Context) error {
	l := s.lggr.Named("Legacy Gas Station Server").With(
		"jobID", s.spec.ID,
		"externalJobID", s.spec.ExternalJobID,
		"chainID", s.spec.LegacyGasStationServerSpec.EVMChainID.ToInt().Uint64(),
		"ccipChainSelector", s.spec.LegacyGasStationServerSpec.CCIPChainSelector.ToInt().Uint64(),
	)
	chain, err := s.cc.Get(s.spec.LegacyGasStationServerSpec.EVMChainID.String())
	if err != nil {
		return err
	}
	forwarder, err := forwarder.NewLegacyGasStationForwarder(s.spec.LegacyGasStationServerSpec.ForwarderAddress.Address(), chain.Client())
	if err != nil {
		return errors.Wrap(err, "initializing forwarder")
	}
	if err = checkFromAddressesExist(s.spec, s.ks); err != nil {
		return err
	}

	orm := NewORM(s.db, l, chain.Config().Database())

	transactor, err := NewBlockchainTransactor(
		s.lggr,
		s.db,
		chain.TxManager(),
		s.ks,
		s.spec.LegacyGasStationServerSpec.FromAddresses,
		s.spec.LegacyGasStationServerSpec.EVMChainID.ToInt().Uint64(),
		orm,
	)
	if err != nil {
		return err
	}

	client, err := eal.NewBlockchainClient(
		s.lggr,
		chain.TxManager(),
		s.ks,
		s.spec.LegacyGasStationServerSpec.FromAddresses,
		s.spec.LegacyGasStationServerSpec.EVMChainID.ToInt().Uint64(),
		chain.Config().EVM(),
		chain.Client(),
	)
	if err != nil {
		return err
	}

	reqHandler, err := legacygasstation.NewRequestHandler(
		l,
		forwarder,
		chain.ID().Uint64(),
		s.spec.LegacyGasStationServerSpec.CCIPChainSelector.ToInt().Uint64(),
		client,
		transactor,
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

func (s *gasStationService) Close() error {
	s.rr.DeregisterHandler(s.spec.LegacyGasStationServerSpec.CCIPChainSelector.ToInt())
	return nil
}

// CheckFromAddressesExist returns an error if and only if one of the addresses
// in the LegacyGasStationServerSpec spec's fromAddresses field does not exist in the keystore.
func checkFromAddressesExist(jb job.Job, gethks keystore.Eth) (err error) {
	for _, a := range jb.LegacyGasStationServerSpec.FromAddresses {
		_, err2 := gethks.Get(a.Hex())
		err = multierr.Append(err, err2)
	}
	return
}
