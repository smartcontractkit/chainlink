package chains

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type ChainsORM[I ID, CFG Config, C DBChain[I, CFG]] interface {
	Chain(I, ...pg.QOpt) (C, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]C, int, error)
	CreateChain(id I, config CFG, qopts ...pg.QOpt) (C, error)
	UpdateChain(id I, enabled bool, config CFG, qopts ...pg.QOpt) (C, error)
	DeleteChain(id I, qopts ...pg.QOpt) error
	GetChainsByIDs(ids []I) (chains []C, err error)
	EnabledChains(...pg.QOpt) ([]C, error)
}

type NodesORM[I ID, N Node] interface {
	CreateNode(N, ...pg.QOpt) (N, error)
	DeleteNode(int32, ...pg.QOpt) error
	GetNodesByChainIDs(chainIDs []I, qopts ...pg.QOpt) (nodes []N, err error)
	NodeNamed(string, ...pg.QOpt) (N, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error)
	NodesForChain(chainID I, offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error)
}

// ORM manages chains and nodes.
type ORM[I ID, C Config, N Node] interface {
	ChainsORM[I, C, DBChain[I, C]]

	NodesORM[I, N]

	StoreString(chainID I, key, val string) error
	Clear(chainID I, key string) error

	// SetupNodes is a shim to help with configuring multiple nodes via ENV.
	// All existing nodes are dropped, and any missing chains are automatically created.
	// Then all nodes are inserted, and conflicts are ignored.
	SetupNodes(nodes []N, chainIDs []I) error
	EnsureChains([]I, ...pg.QOpt) error
}

type orm[I ID, C Config, N Node] struct {
	*chainsORM[I, C]
	*nodesORM[I, N]
}

// NewORM returns an ORM backed by q, for the tables <prefix>_chains and <prefix>_nodes with column <prefix>_chain_id.
// Additional Node fields should be included in nodeCols.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewORM[I ID, C Config, N Node](q pg.Q, prefix string, nodeCols ...string) ORM[I, C, N] {
	return orm[I, C, N]{
		newChainsORM[I, C](q, prefix),
		newNodesORM[I, N](q, prefix, nodeCols...),
	}
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func (o orm[I, C, N]) SetupNodes(nodes []N, ids []I) error {
	return o.chainsORM.q.Transaction(func(q pg.Queryer) error {
		tx := pg.WithQueryer(q)
		if err := o.truncateNodes(tx); err != nil {
			return err
		}

		if err := o.EnsureChains(ids, tx); err != nil {
			return err
		}

		return o.ensureNodes(nodes, tx)
	})
}

// DBChain is a generic DB chain for an ID and Config.
//
// A DBChain type alias can be used for convenience:
//
//	type DBChain = chains.DBChain[string, pkg.ChainCfg]
type DBChain[I ID, C Config] struct {
	ID        I
	Cfg       C
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

// chainsORM is a generic ORM for chains.
type chainsORM[I ID, C Config] struct {
	q      pg.Q
	prefix string
}

// newChainsORM returns an chainsORM backed by q, for the table <prefix>_chains.
func newChainsORM[I ID, C Config](q pg.Q, prefix string) *chainsORM[I, C] {
	return &chainsORM[I, C]{q: q, prefix: prefix}
}

func (o *chainsORM[I, C]) Chain(id I, qopts ...pg.QOpt) (dbchain DBChain[I, C], err error) {
	q := o.q.WithOpts(qopts...)
	chainSQL := fmt.Sprintf(`SELECT * FROM %s_chains WHERE id = $1;`, o.prefix)
	err = q.Get(&dbchain, chainSQL, id)
	return
}

func (o *chainsORM[I, C]) GetChainsByIDs(ids []I) (chains []DBChain[I, C], err error) {
	sql := fmt.Sprintf(`SELECT * FROM %s_chains WHERE id = ANY($1) ORDER BY created_at, id;`, o.prefix)

	chainIDs := pq.Array(ids)
	if err = o.q.Select(&chains, sql, chainIDs); err != nil {
		return nil, err
	}

	return chains, nil
}

func (o *chainsORM[I, C]) CreateChain(id I, config C, qopts ...pg.QOpt) (chain DBChain[I, C], err error) {
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf(`INSERT INTO %s_chains (id, cfg, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`, o.prefix)
	err = q.Get(&chain, sql, id, config)
	return
}

func (o *chainsORM[I, C]) EnsureChains(ids []I, qopts ...pg.QOpt) (err error) {
	named := make([]struct{ ID I }, len(ids))
	for i, id := range ids {
		named[i].ID = id
	}
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf("INSERT INTO %s_chains (id, created_at, updated_at) VALUES (:id, NOW(), NOW()) ON CONFLICT DO NOTHING;", o.prefix)

	if _, err := q.NamedExec(sql, named); err != nil {
		return errors.Wrapf(err, "failed to insert chains %v", ids)
	}

	return nil
}

func (o *chainsORM[I, C]) UpdateChain(id I, enabled bool, config C, qopts ...pg.QOpt) (chain DBChain[I, C], err error) {
	q := o.q.WithOpts(qopts...)
	sql := fmt.Sprintf(`UPDATE %s_chains SET enabled = $1, cfg = $2, updated_at = now() WHERE id = $3 RETURNING *`, o.prefix)
	err = q.Get(&chain, sql, enabled, config, id)
	return
}

// StoreString saves a string value into the config for the given chain and key
func (o *chainsORM[I, C]) StoreString(chainID I, name, val string) error {
	s := fmt.Sprintf(`UPDATE %s_chains SET cfg = cfg || jsonb_build_object($1::text, $2::text) WHERE id = $3`, o.prefix)
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
func (o *chainsORM[I, C]) Clear(chainID I, name string) error {
	s := fmt.Sprintf(`UPDATE %s_chains SET cfg = cfg - $1 WHERE id = $2`, o.prefix)
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

func (o *chainsORM[I, C]) DeleteChain(id I, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := fmt.Sprintf(`DELETE FROM %s_chains WHERE id = $1`, o.prefix)
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

func (o *chainsORM[I, C]) Chains(offset, limit int, qopts ...pg.QOpt) (chains []DBChain[I, C], count int, err error) {
	err = o.q.WithOpts(qopts...).Transaction(func(q pg.Queryer) error {
		if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s_chains", o.prefix)); err != nil {
			return errors.Wrap(err, "failed to fetch chains count")
		}

		sql := fmt.Sprintf(`SELECT * FROM %s_chains ORDER BY created_at, id LIMIT $1 OFFSET $2;`, o.prefix)
		err = q.Select(&chains, sql, pg.Limit(limit), offset)
		return errors.Wrap(err, "failed to fetch chains")
	}, pg.OptReadOnlyTx())

	return
}

func (o *chainsORM[I, C]) EnabledChains(qopts ...pg.QOpt) (chains []DBChain[I, C], err error) {
	q := o.q.WithOpts(qopts...)
	chainsSQL := fmt.Sprintf(`SELECT * FROM %s_chains WHERE enabled ORDER BY created_at, id;`, o.prefix)
	if err = q.Select(&chains, chainsSQL); err != nil {
		return
	}
	return
}

// nodesORM is a generic ORM for nodes.
type nodesORM[I ID, N Node] struct {
	q           pg.Q
	prefix      string
	createNodeQ string
	ensureNodeQ string
}

func newNodesORM[I ID, N Node](q pg.Q, prefix string, nodeCols ...string) *nodesORM[I, N] {
	// pre-compute query for CreateNode
	var withColon []string
	for _, c := range nodeCols {
		withColon = append(withColon, ":"+c)
	}
	query := fmt.Sprintf(`INSERT INTO %s_nodes (name, %s_chain_id, %s, created_at, updated_at)
		VALUES (:name, :%s_chain_id, %s, now(), now())`,
		prefix, prefix, strings.Join(nodeCols, ", "), prefix, strings.Join(withColon, ", "))

	return &nodesORM[I, N]{q: q, prefix: prefix,
		createNodeQ: query + ` RETURNING *;`,
		ensureNodeQ: query + ` ON CONFLICT DO NOTHING;`,
	}
}

func (o *nodesORM[I, N]) ensureNodes(nodes []N, qopts ...pg.QOpt) (err error) {
	q := o.q.WithOpts(qopts...)
	_, err = q.NamedExec(o.ensureNodeQ, nodes)
	err = errors.Wrap(err, "failed to insert nodes")
	return
}

func (o *nodesORM[I, N]) CreateNode(data N, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	stmt, err := q.PrepareNamed(o.createNodeQ)
	if err != nil {
		return node, err
	}
	err = stmt.Get(&node, data)
	return node, err
}

func (o *nodesORM[I, N]) DeleteNode(id int32, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	query := fmt.Sprintf(`DELETE FROM %s_nodes WHERE id = $1`, o.prefix)
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

func (o *nodesORM[I, N]) truncateNodes(qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(fmt.Sprintf(`TRUNCATE %s_nodes;`, o.prefix))
	if err != nil {
		return errors.Wrapf(err, "failed to truncate %s_nodes table", o.prefix)
	}
	return nil
}

func (o *nodesORM[I, N]) NodeNamed(name string, qopts ...pg.QOpt) (node N, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Get(&node, fmt.Sprintf("SELECT * FROM %s_nodes WHERE name = $1;", o.prefix), name)

	return
}

func (o *nodesORM[I, N]) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error) {
	err = o.q.WithOpts(qopts...).Transaction(func(q pg.Queryer) error {
		if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s_nodes", o.prefix)); err != nil {
			return errors.Wrap(err, "failed to fetch nodes count")
		}

		sql := fmt.Sprintf(`SELECT * FROM %s_nodes ORDER BY created_at, id LIMIT $1 OFFSET $2;`, o.prefix)
		err = q.Select(&nodes, sql, pg.Limit(limit), offset)
		return errors.Wrap(err, "failed to fetch nodes")
	}, pg.OptReadOnlyTx())

	return
}

func (o *nodesORM[I, N]) NodesForChain(chainID I, offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error) {
	err = o.q.WithOpts(qopts...).Transaction(func(q pg.Queryer) error {
		if err = q.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s_nodes WHERE %s_chain_id = $1", o.prefix, o.prefix), chainID); err != nil {
			return errors.Wrap(err, "failed to fetch nodes count")
		}

		sql := fmt.Sprintf(`SELECT * FROM %s_nodes WHERE %s_chain_id = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3;`, o.prefix, o.prefix)
		err = q.Select(&nodes, sql, chainID, pg.Limit(limit), offset)
		return errors.Wrap(err, "failed to fetch nodes")
	}, pg.OptReadOnlyTx())

	return
}

func (o *nodesORM[I, N]) GetNodesByChainIDs(chainIDs []I, qopts ...pg.QOpt) (nodes []N, err error) {
	sql := fmt.Sprintf(`SELECT * FROM %s_nodes WHERE %s_chain_id = ANY($1) ORDER BY created_at, id;`, o.prefix, o.prefix)

	cids := pq.Array(chainIDs)
	if err = o.q.WithOpts(qopts...).Select(&nodes, sql, cids); err != nil {
		return nil, err
	}

	return nodes, nil
}
