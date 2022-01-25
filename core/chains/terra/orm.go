//go:build terra
// +build terra

package terra

import (
	"database/sql"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

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
	return &orm{q: pg.NewQ(db, lggr.Named("ORM"), cfg)}
}

func (o *orm) Chain(id string, qopts ...pg.QOpt) (dbchain db.Chain, err error) {
	q := o.q.WithOpts(qopts...)
	chainSQL := `SELECT * FROM terra_chains WHERE id = $1;`
	if err = q.Get(&dbchain, chainSQL, id); err != nil {
		return
	}
	nodesSQL := `SELECT * FROM terra_nodes WHERE terra_chain_id = $1 ORDER BY created_at, id;`
	if err = q.Select(&dbchain.Nodes, nodesSQL, id); err != nil {
		return
	}
	return
}

func (o *orm) CreateChain(id string, config db.ChainCfg, qopts ...pg.QOpt) (chain db.Chain, err error) {
	q := o.q.WithOpts(qopts...)
	sql := `INSERT INTO terra_chains (id, cfg, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	err = q.Get(&chain, sql, id, config)
	return
}

func (o *orm) UpdateChain(id string, enabled bool, config db.ChainCfg, qopts ...pg.QOpt) (chain db.Chain, err error) {
	q := o.q.WithOpts(qopts...)
	sql := `UPDATE terra_chains SET enabled = $1, cfg = $2, updated_at = now() WHERE id = $3 RETURNING *`
	err = q.Get(&chain, sql, enabled, config, id)
	return
}

func (o *orm) DeleteChain(id string, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := `DELETE FROM terra_chains WHERE id = $1`
	result, err := q.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (o *orm) Chains(offset, limit int, qopts ...pg.QOpt) (chains []db.Chain, count int, err error) {
	q := o.q.WithOpts(qopts...)
	if err = q.Get(&count, "SELECT COUNT(*) FROM terra_chains"); err != nil {
		return
	}

	sql := `SELECT * FROM terra_chains ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = q.Select(&chains, sql, limit, offset); err != nil {
		return
	}

	return
}

func (o *orm) EnabledChainsWithNodes(qopts ...pg.QOpt) (chains []db.Chain, err error) {
	q := o.q.WithOpts(qopts...)
	chainsSQL := `SELECT * FROM terra_chains WHERE enabled ORDER BY created_at, id;`
	if err = q.Select(&chains, chainsSQL); err != nil {
		return
	}
	var nodes []db.Node
	nodesSQL := `SELECT * FROM terra_nodes ORDER BY created_at, id;`
	if err = q.Select(&nodes, nodesSQL); err != nil {
		return
	}
	nodemap := make(map[string][]db.Node)
	for _, n := range nodes {
		nodemap[n.TerraChainID] = append(nodemap[n.TerraChainID], n)
	}
	for i, c := range chains {
		chains[i].Nodes = nodemap[c.ID]
	}
	return
}

func (o *orm) CreateNode(data types.NewNode, qopts ...pg.QOpt) (node db.Node, err error) {
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
	query := `DELETE FROM terra_nodes WHERE id = $1`
	result, err := q.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (o *orm) Node(id int32, qopts ...pg.QOpt) (node db.Node, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, "SELECT * FROM terra_nodes WHERE id = $1;", id)

	return
}

func (o *orm) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error) {
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

func (o *orm) NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error) {
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
