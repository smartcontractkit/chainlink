package chains

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// Chain is a generic DB chain for any configuration type.
// CFG normally implements sql.Scanner and driver.Valuer, but that is not enforced here.
type Chain[CFG any] struct {
	ID        string
	Cfg       CFG
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

type ChainsORM[CFG any, CH Chain[CFG]] struct {
	q     pg.Q
	table string
}

// NewChainsORM returns an ChainsORM backed by db.
func NewChainsORM[CFG any, CH Chain[CFG]](q pg.Q, table string) *ChainsORM[CFG, CH] {
	return &ChainsORM[CFG, CH]{q: q, table: table}
}

func (o *ChainsORM[CFG, CH]) Chain(id string, qopts ...pg.QOpt) (dbchain CH, err error) {
	q := o.q.WithOpts(qopts...)
	chainSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1;`, o.table)
	err = q.Get(&dbchain, chainSQL, id)
	return
}

func (o *ChainsORM[CFG, CH]) CreateChain(id string, config CFG, qopts ...pg.QOpt) (chain CH, err error) {
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf(`INSERT INTO %s (id, cfg, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`, o.table)
	err = q.Get(&chain, sql, id, config)
	return
}

func (o *ChainsORM[CFG, CH]) UpdateChain(id string, enabled bool, config CFG, qopts ...pg.QOpt) (chain CH, err error) {
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf(`UPDATE %s SET enabled = $1, cfg = $2, updated_at = now() WHERE id = $3 RETURNING *`, o.table)
	err = q.Get(&chain, sql, enabled, config, id)
	return
}

func (o *ChainsORM[CFG, CH]) DeleteChain(id string, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, o.table)
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

func (o *ChainsORM[CFG, CH]) Chains(offset, limit int, qopts ...pg.QOpt) (chains []CH, count int, err error) {
	q := o.q.WithOpts(qopts...)
	if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", o.table)); err != nil {
		return
	}

	sql := fmt.Sprintf(`SELECT * FROM %s ORDER BY created_at, id LIMIT $1 OFFSET $2;`, o.table)
	if err = q.Select(&chains, sql, limit, offset); err != nil {
		return
	}

	return
}

func (o *ChainsORM[CFG, CH]) EnabledChains(qopts ...pg.QOpt) (chains []CH, err error) {
	q := o.q.WithOpts(qopts...)
	chainsSQL := fmt.Sprintf(`SELECT * FROM %s WHERE enabled ORDER BY created_at, id;`, o.table)
	if err = q.Select(&chains, chainsSQL); err != nil {
		return
	}
	return
}

type NodesORM[NEW, N any] struct {
	q             pg.Q
	table         string
	chainID       string
	createNodeSQL string
}

// NewNodesORM returns a NodesORM backed by db.
func NewNodesORM[NEW, N any](q pg.Q, table, chainID, createNodeSQL string) *NodesORM[NEW, N] {
	return &NodesORM[NEW, N]{q: q, table: table, chainID: chainID, createNodeSQL: createNodeSQL}
}

func (o *NodesORM[NEW, N]) CreateNode(data NEW, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	stmt, err := q.PrepareNamed(o.createNodeSQL)
	if err != nil {
		return node, err
	}
	err = stmt.Get(&node, data)
	return node, err
}

func (o *NodesORM[NEW, N]) DeleteNode(id int32, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, o.table)
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

func (o *NodesORM[NEW, N]) Node(id int32, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, fmt.Sprintf("SELECT * FROM %s WHERE id = $1;", o.table), id)

	return
}

func (o *NodesORM[NEW, N]) NodeNamed(name string, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, fmt.Sprintf("SELECT * FROM %s WHERE name = $1;", o.table), name)

	return
}

func (o *NodesORM[NEW, N]) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error) {
	q := o.q.WithOpts(qopts...)
	if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", o.table)); err != nil {
		return
	}

	sql := fmt.Sprintf(`SELECT * FROM %s ORDER BY created_at, id LIMIT $1 OFFSET $2;`, o.table)
	if err = q.Select(&nodes, sql, limit, offset); err != nil {
		return
	}

	return
}

func (o *NodesORM[NEW, N]) NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error) {
	q := o.q.WithOpts(qopts...)
	if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", o.table, o.chainID), chainID); err != nil {
		return
	}

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3;`, o.table, o.chainID)
	if err = q.Select(&nodes, sql, chainID, limit, offset); err != nil {
		return
	}

	return
}
