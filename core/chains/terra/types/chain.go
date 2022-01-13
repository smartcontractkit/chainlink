package types

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	terraconfig "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"

	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type ChainSetOpts struct {
	Config           config.GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Terra
	EventBroadcaster pg.EventBroadcaster
	ORM              ORM
}

var _ terra.ChainSet = (*chainSet)(nil)

type chainSet struct {
	chains map[string]terra.Chain
}

func NewChainSet(opts ChainSetOpts) (terra.ChainSet, error) {
	dbchains, err := opts.ORM.EnabledChainsWithNodes()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	cs := &chainSet{
		chains: make(map[string]terra.Chain),
	}
	for _, c := range dbchains {
		n := c.Nodes[0] //TODO client pool
		var err2 error
		cs.chains[c.ID], err2 = NewChain(opts.DB, opts.KeyStore, n, opts.Config, opts.EventBroadcaster, c.Cfg, opts.Logger)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}
	return cs, err
}

func (c *chainSet) Get(id string) (terra.Chain, error) {
	return c.chains[id], nil
}

var _ terra.Chain = (*chain)(nil)

type chain struct {
	id     string
	cfg    terraconfig.ChainCfg
	client *terraclient.Client
	txm    *terratxm.Txm
}

func NewChain(db *sqlx.DB, ks keystore.Terra, node Node, logCfg pg.LogConfig, eb pg.EventBroadcaster, cfg terraconfig.ChainCfg, lggr logger.Logger) (terra.Chain, error) {
	id := node.TerraChainID
	client, err := terraclient.NewClient(id,
		node.TendermintURL, node.FCDURL, 10, lggr)
	if err != nil {
		return nil, err
	}
	txm, err := terratxm.NewTxm(db, client, cfg.FallbackGasPriceULuna, cfg.GasLimitMultiplier, ks, lggr.Named(id), logCfg, eb, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &chain{
		id:     id,
		cfg:    cfg,
		client: client,
		txm:    txm,
	}, nil
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() terraconfig.ChainCfg {
	return c.cfg
}

func (c *chain) MsgEnqueuer() terra.MsgEnqueuer {
	return c.txm
}

func (c *chain) Reader() terraclient.Reader {
	return c.client
}

func (c *chain) Start() error {
	//TODO implement me
	panic("implement me")
}

func (c *chain) Close() error {
	//TODO implement me
	panic("implement me")
}

func (c *chain) Ready() error {
	//TODO implement me
	panic("implement me")
}

func (c *chain) Healthy() error {
	//TODO implement me
	panic("implement me")
}
