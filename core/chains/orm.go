package chains

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type ChainConfigs[I ID, CFG Config, C ChainConfig[I, CFG]] interface {
	Chain(I, ...pg.QOpt) (C, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]C, int, error)
	GetChainsByIDs(ids []I) (chains []C, err error)
}

type NodeConfigs[I ID, N Node] interface {
	GetNodesByChainIDs(chainIDs []I, qopts ...pg.QOpt) (nodes []N, err error)
	NodeNamed(string, ...pg.QOpt) (N, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error)
	NodesForChain(chainID I, offset, limit int, qopts ...pg.QOpt) (nodes []N, count int, err error)
}

// ORM manages chains and nodes.
type ORM[I ID, C Config, N Node] interface {
	ChainConfigs[I, C, ChainConfig[I, C]]

	NodeConfigs[I, N]

	// EnsureChains creates an entry for any chain IDs which don't already exist.
	// TODO remove with https://smartcontract-it.atlassian.net/browse/BCF-1474
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

// ChainConfig is a generic DB chain for an ID and Config.
//
// A ChainConfig type alias can be used for convenience:
//
//	type ChainConfig = chains.ChainConfig[string, pkg.ChainCfg]
type ChainConfig[I ID, C Config] struct {
	ID      I
	Cfg     C
	Enabled bool
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

func (o *chainsORM[I, C]) Chain(id I, qopts ...pg.QOpt) (cc ChainConfig[I, C], err error) {
	q := o.q.WithOpts(qopts...)
	chainSQL := fmt.Sprintf(`SELECT * FROM %s_chains WHERE id = $1;`, o.prefix)
	err = q.Get(&cc, chainSQL, id)
	return
}

func (o *chainsORM[I, C]) GetChainsByIDs(ids []I) (chains []ChainConfig[I, C], err error) {
	sql := fmt.Sprintf(`SELECT * FROM %s_chains WHERE id = ANY($1) ORDER BY created_at, id;`, o.prefix)

	chainIDs := pq.Array(ids)
	if err = o.q.Select(&chains, sql, chainIDs); err != nil {
		return nil, err
	}

	return chains, nil
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

func (o *chainsORM[I, C]) Chains(offset, limit int, qopts ...pg.QOpt) (chains []ChainConfig[I, C], count int, err error) {
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
