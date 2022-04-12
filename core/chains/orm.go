package chains

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// Chain is a generic DB chain for any configuration type.
// CFG normally implements sql.Scanner and driver.Valuer, but that is not enforced here.
type Chain[ID, CFG any] struct {
	ID        ID
	Cfg       CFG
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

//TODO doc
type ChainsORM[ID, CFG any, CH Chain[ID, CFG]] struct {
	q     pg.Q
	table string
}

// NewChainsORM returns an ChainsORM backed by db.
func NewChainsORM[ID, CFG any, CH Chain[ID, CFG]](q pg.Q, table string) *ChainsORM[ID, CFG, CH] {
	return &ChainsORM[ID, CFG, CH]{q: q, table: table}
}

func (o *ChainsORM[ID, CFG, CH]) Chain(id ID, qopts ...pg.QOpt) (dbchain CH, err error) {
	q := o.q.WithOpts(qopts...)
	chainSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1;`, o.table)
	err = q.Get(&dbchain, chainSQL, id)
	return
}

func (o *ChainsORM[ID, CFG, CH]) GetChainsByIDs(ids []ID) (chains []CH, err error) {
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE id = ANY($1) ORDER BY created_at, id;`, o.table)

	chainIDs := pq.Array(ids)
	if err = o.q.Select(&chains, sql, chainIDs); err != nil {
		return nil, err
	}

	return chains, nil
}

func (o *ChainsORM[ID, CFG, CH]) CreateChain(id ID, config CFG, qopts ...pg.QOpt) (chain CH, err error) {
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf(`INSERT INTO %s (id, cfg, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`, o.table)
	err = q.Get(&chain, sql, id, config)
	return
}

func (o *ChainsORM[ID, CFG, CH]) UpdateChain(id ID, enabled bool, config CFG, qopts ...pg.QOpt) (chain CH, err error) {
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf(`UPDATE %s SET enabled = $1, cfg = $2, updated_at = now() WHERE id = $3 RETURNING *`, o.table)
	err = q.Get(&chain, sql, enabled, config, id)
	return
}

// StoreString saves a string value into the config for the given chain and key
func (o *ChainsORM[ID, CFG, CH]) StoreString(chainID ID, name, val string) error {
	s := fmt.Sprintf(`UPDATE %s SET cfg = cfg || jsonb_build_object($1::text, $2::text) WHERE id = $3`, o.table)
	res, err := o.q.Exec(s, name, val, chainID)
	if err != nil {
		return errors.Wrapf(err, "failed to store chain config for chain ID %v", chainID)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.Wrapf(sql.ErrNoRows, "no chain found with ID %v", chainID)
	}
	return nil
}

// Clear deletes a config value for the given chain and key
func (o *ChainsORM[ID, CFG, CH]) Clear(chainID ID, name string) error {
	s := fmt.Sprintf(`UPDATE %s SET cfg = cfg - $1 WHERE id = $2`, o.table)
	res, err := o.q.Exec(s, name, chainID)
	if err != nil {
		return errors.Wrapf(err, "failed to clear chain config for chain ID %v", chainID)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.Wrapf(sql.ErrNoRows, "no chain found with ID %v", chainID)
	}
	return nil
}

func (o *ChainsORM[ID, CFG, CH]) DeleteChain(id ID, qopts ...pg.QOpt) error {
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

func (o *ChainsORM[ID, CFG, CH]) Chains(offset, limit int, qopts ...pg.QOpt) (chains []CH, count int, err error) {
	err = o.q.WithOpts(qopts...).Transaction(func(q pg.Queryer) error {
		if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", o.table)); err != nil {
			return errors.Wrap(err, "failed to fetch chains count")
		}

		sql := fmt.Sprintf(`SELECT * FROM %s ORDER BY created_at, id LIMIT $1 OFFSET $2;`, o.table)
		err = q.Select(&chains, sql, pg.Limit(limit), offset)
		return errors.Wrap(err, "failed to fetch chains")
	}, pg.OptReadOnlyTx())

	return
}

func (o *ChainsORM[ID, CFG, CH]) EnabledChains(qopts ...pg.QOpt) (chains []CH, err error) {
	q := o.q.WithOpts(qopts...)
	chainsSQL := fmt.Sprintf(`SELECT * FROM %s WHERE enabled ORDER BY created_at, id;`, o.table)
	if err = q.Select(&chains, chainsSQL); err != nil {
		return
	}
	return
}

//TODO doc
type NodesORM[ID, NEW, N any] struct {
	q             pg.Q
	table         string
	chainID       string
	createNodeSQL string
}

// NewNodesORM returns a NodesORM backed by db.
func NewNodesORM[ID, NEW, N any](q pg.Q, table, chainID, createNodeSQL string) *NodesORM[ID, NEW, N] {
	return &NodesORM[ID, NEW, N]{q: q, table: table, chainID: chainID, createNodeSQL: createNodeSQL}
}

func (o *NodesORM[ID, NEW, N]) CreateNode(data NEW, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	stmt, err := q.PrepareNamed(o.createNodeSQL)
	if err != nil {
		return node, err
	}
	err = stmt.Get(&node, data)
	return node, err
}

func (o *NodesORM[ID, NEW, N]) DeleteNode(id int32, qopts ...pg.QOpt) error {
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

func (o *NodesORM[ID, NEW, N]) Node(id int32, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, fmt.Sprintf("SELECT * FROM %s WHERE id = $1;", o.table), id)

	return
}

func (o *NodesORM[ID, NEW, N]) NodeNamed(name string, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, fmt.Sprintf("SELECT * FROM %s WHERE name = $1;", o.table), name)

	return
}

func (o *NodesORM[ID, NEW, N]) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error) {
	err = o.q.WithOpts(qopts...).Transaction(func(q pg.Queryer) error {
		if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", o.table)); err != nil {
			return errors.Wrap(err, "failed to fetch nodes count")
		}

		sql := fmt.Sprintf(`SELECT * FROM %s ORDER BY created_at, id LIMIT $1 OFFSET $2;`, o.table)
		err = q.Select(&nodes, sql, pg.Limit(limit), offset)
		return errors.Wrap(err, "failed to fetch nodes")
	}, pg.OptReadOnlyTx())

	return
}

func (o *NodesORM[ID, NEW, N]) NodesForChain(chainID ID, offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error) {
	err = o.q.WithOpts(qopts...).Transaction(func(q pg.Queryer) error {
		if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", o.table, o.chainID), chainID); err != nil {
			return errors.Wrap(err, "failed to fetch nodes count")
		}

		sql := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3;`, o.table, o.chainID)
		err = q.Select(&nodes, sql, chainID, pg.Limit(limit), offset)
		return errors.Wrap(err, "failed to fetch nodes")
	}, pg.OptReadOnlyTx())

	return
}

func (o *NodesORM[ID, NEW, N]) GetNodesByChainIDs(chainIDs []ID, qopts ...pg.QOpt) (nodes []N, err error) {
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE %s = ANY($1) ORDER BY created_at, id;`, o.table, o.chainID)

	cids := pq.Array(chainIDs)
	if err = o.q.WithOpts(qopts...).Select(&nodes, sql, cids); err != nil {
		return nil, err
	}

	return nodes, nil
}
