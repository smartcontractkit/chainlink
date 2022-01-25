//go:build !terra
// +build !terra

package terra

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

type ChainSet struct {
}

type ChainSetOpts struct {
	Config           coreconfig.GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Terra
	EventBroadcaster pg.EventBroadcaster
	ORM              ORM
}

type ORM struct {
}

func (O ORM) Chain(s string, opt ...pg.QOpt) (db.Chain, error) {
	panic("implement me")
}

func (O ORM) Chains(offset, limit int, qopts ...pg.QOpt) ([]db.Chain, int, error) {
	panic("implement me")
}

func (O ORM) CreateChain(id string, config db.ChainCfg, qopts ...pg.QOpt) (db.Chain, error) {
	panic("implement me")
}

func (O ORM) UpdateChain(id string, enabled bool, config db.ChainCfg, qopts ...pg.QOpt) (db.Chain, error) {
	panic("implement me")
}

func (O ORM) DeleteChain(id string, qopts ...pg.QOpt) error {
	panic("implement me")
}

func (O ORM) EnabledChainsWithNodes(opt ...pg.QOpt) ([]db.Chain, error) {
	panic("implement me")
}

func (O ORM) CreateNode(node types.NewNode, opt ...pg.QOpt) (db.Node, error) {
	panic("implement me")
}

func (O ORM) DeleteNode(i int32, opt ...pg.QOpt) error {
	panic("implement me")
}

func (O ORM) Node(i int32, opt ...pg.QOpt) (db.Node, error) {
	panic("implement me")
}

func (O ORM) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error) {
	panic("implement me")
}

func (O ORM) NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error) {
	panic("implement me")
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) ORM {
	return ORM{}
}

func NewChainSet(opts ChainSetOpts) (ChainSet, error) {
	return ChainSet{}, nil
}
func (ChainSet) ORM() ORM {
	return ORM{}
}

func (ChainSet) Add(id string, config db.ChainCfg) (db.Chain, error) {
	return db.Chain{}, nil
}

func (ChainSet) Remove(id string) error {
	return nil
}
func (ChainSet) Configure(id string, enabled bool, config db.ChainCfg) (db.Chain, error) {
	return db.Chain{}, nil
}
