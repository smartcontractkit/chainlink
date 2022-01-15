package terra

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	terraconfig "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type orm struct {
	q pg.Q
}

var _ types.ORM = (*orm)(nil)

// NewORM returns an ORM backed by db.
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) types.ORM {
	return &orm{q: pg.NewQ(db, lggr.Named("TerraORM"), cfg)}
}

// ErrNoRowsAffected is returned when rows should have been affected but were not.
var ErrNoRowsAffected = errors.New("no rows affected")

var defaultCfg = terraconfig.ChainCfg{
	FallbackGasPriceULuna: "0.01",
	GasLimitMultiplier:    1.5,
}

func (o *orm) EnabledChainsWithNodes(qopts ...pg.QOpt) (chains []types.Chain, err error) {
	q := o.q.WithOpts(qopts...)
	var nodes []types.Node
	nodesSQL := `SELECT * FROM terra_nodes ORDER BY created_at, id;`
	if err = q.Select(&nodes, nodesSQL); err != nil {
		return
	}
	nodemap := make(map[string][]types.Node)
	for _, n := range nodes {
		nodemap[n.TerraChainID] = append(nodemap[n.TerraChainID], n)
	}
	for id, ns := range nodemap {
		chains = append(chains, types.Chain{
			ID:    id,
			Nodes: ns,
			Cfg:   defaultCfg,
		})
	}
	return chains, nil
}

func (o *orm) Chain(id string, qopts ...pg.QOpt) (types.Chain, error) {
	q := o.q.WithOpts(qopts...)
	var nodes []types.Node
	nodesSQL := `SELECT * FROM terra_nodes WHERE terra_chain_id = $1 ORDER BY created_at, id;`
	if err := q.Select(&nodes, nodesSQL, id); err != nil {
		return types.Chain{}, err
	}
	return types.Chain{
		ID:    id,
		Nodes: nodes,
		Cfg:   defaultCfg,
	}, nil
}

func (o *orm) CreateNode(data types.NewNode, qopts ...pg.QOpt) (node types.Node, err error) {
	q := o.q.WithOpts(qopts...)
	sql := `INSERT INTO terra_nodes (name, terra_chain_id, tendermint_url, fcd_url, created_at, updated_at)
	VALUES (:name, :terra_chain_id, :tendermint_url, :fcd_url, now(), now())
	RETURNING *;`
	stmt, err := q.PrepareNamed(sql)
	if err != nil {
		return node, err
	}
	err = stmt.Get(&node, data)
	return node, err
}

func (o *orm) DeleteNode(id int32, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	sql := `DELETE FROM terra_nodes WHERE id = $1`
	result, err := q.Exec(sql, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (o *orm) Node(id int32, qopts ...pg.QOpt) (node types.Node, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, "SELECT * FROM terra_nodes WHERE id = $1;", id)

	return
}

func (o *orm) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []types.Node, count int, err error) {
	q := o.q.WithOpts(qopts...)
	if err = q.Get(&count, "SELECT COUNT(*) FROM terra_nodes"); err != nil {
		return
	}

	sql := `SELECT * FROM terra_nodes ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = q.Select(&nodes, sql, limit, offset); err != nil {
		return
	}

	return
}

func (o *orm) NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []types.Node, count int, err error) {
	q := o.q.WithOpts(qopts...)
	if err = q.Get(&count, "SELECT COUNT(*) FROM terra_nodes WHERE terra_chain_id = $1", chainID); err != nil {
		return
	}

	sql := `SELECT * FROM terra_nodes WHERE terra_chain_id = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3;`
	if err = q.Select(&nodes, sql, chainID, limit, offset); err != nil {
		return
	}

	return
}
