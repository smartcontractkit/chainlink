package evm

import (
	"math/big"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type orm struct {
	db *sqlx.DB
}

var _ types.ORM = (*orm)(nil)

func NewORM(db *sqlx.DB) types.ORM {
	return &orm{db}
}

var ErrNoRowsAffected = errors.New("no rows affected")

func (o *orm) Chain(id utils.Big) (chain types.Chain, err error) {
	sql := `SELECT * FROM evm_chains WHERE id = $1`
	err = o.db.Get(&chain, sql, id)
	return chain, err
}

func (o *orm) CreateChain(id utils.Big, config types.ChainCfg) (chain types.Chain, err error) {
	sql := `INSERT INTO evm_chains (id, cfg, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	err = o.db.Get(&chain, sql, id, config)
	return chain, err
}

func (o *orm) UpdateChain(id utils.Big, enabled bool, config types.ChainCfg) (chain types.Chain, err error) {
	sql := `UPDATE evm_chains SET enabled = $1, cfg = $2, updated_at = now() WHERE id = $3 RETURNING *`
	err = o.db.Get(&chain, sql, enabled, config, id)
	return chain, err
}

func (o *orm) DeleteChain(id utils.Big) error {
	sql := `DELETE FROM evm_chains WHERE id = $1`
	result, err := o.db.Exec(sql, id)
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

func (o *orm) Chains(offset, limit int) (chains []types.Chain, count int, err error) {
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM evm_chains"); err != nil {
		return
	}

	sql := `SELECT * FROM evm_chains ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = o.db.Select(&chains, sql, limit, offset); err != nil {
		return
	}

	return
}

// GetNodesByChainID fetches allow nodes for the given chain ids.
func (o *orm) GetChainsByIDs(ids []utils.Big) (chains []types.Chain, err error) {
	sql := `SELECT * FROM evm_chains WHERE id = ANY($1) ORDER BY created_at, id;`

	chainIDs := pq.Array(ids)
	if err = o.db.Select(&chains, sql, chainIDs); err != nil {
		return nil, err
	}

	return chains, nil
}

func (o *orm) CreateNode(data types.NewNode) (node types.Node, err error) {
	sql := `INSERT INTO nodes (name, evm_chain_id, ws_url, http_url, send_only, created_at, updated_at)
	VALUES (:name, :evm_chain_id, :ws_url, :http_url, :send_only, now(), now())
	RETURNING *;`
	stmt, err := o.db.PrepareNamed(sql)
	if err != nil {
		return node, err
	}
	err = stmt.Get(&node, data)
	return node, err
}

func (o *orm) DeleteNode(id int64) error {
	sql := `DELETE FROM nodes WHERE id = $1`
	result, err := o.db.Exec(sql, id)
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

func (o *orm) EnabledChainsWithNodes() (chains []types.Chain, err error) {
	var nodes []types.Node
	chainsSQL := `SELECT * FROM evm_chains WHERE enabled ORDER BY created_at, id;`
	if err = o.db.Select(&chains, chainsSQL); err != nil {
		return
	}
	nodesSQL := `SELECT * FROM nodes ORDER BY created_at, id;`
	if err = o.db.Select(&nodes, nodesSQL); err != nil {
		return
	}
	nodemap := make(map[string][]types.Node)
	for _, n := range nodes {
		nodemap[n.EVMChainID.String()] = append(nodemap[n.EVMChainID.String()], n)
	}
	for i, c := range chains {
		chains[i].Nodes = nodemap[c.ID.String()]
	}
	return chains, nil
}

func (o *orm) Nodes(offset, limit int) (nodes []types.Node, count int, err error) {
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM nodes"); err != nil {
		return
	}

	sql := `SELECT * FROM nodes ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = o.db.Select(&nodes, sql, limit, offset); err != nil {
		return
	}

	return
}

// GetNodesByChainID fetches allow nodes for the given chain ids.
func (o *orm) GetNodesByChainIDs(chainIDs []utils.Big) (nodes []types.Node, err error) {
	sql := `SELECT * FROM nodes WHERE evm_chain_id = ANY($1) ORDER BY created_at, id;`

	cids := pq.Array(chainIDs)
	if err = o.db.Select(&nodes, sql, cids); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (o *orm) NodesForChain(chainID utils.Big, offset, limit int) (nodes []types.Node, count int, err error) {
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM nodes WHERE evm_chain_id = $1", chainID); err != nil {
		return
	}

	sql := `SELECT * FROM nodes WHERE evm_chain_id = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3;`
	if err = o.db.Select(&nodes, sql, chainID, limit, offset); err != nil {
		return
	}

	return
}

// StoreString saves a string value into the config for the given chain and key
func (o *orm) StoreString(chainID *big.Int, name, val string) error {
	res, err := o.db.Exec(`UPDATE evm_chains SET cfg = cfg || jsonb_build_object($1::text, $2::text) WHERE id = $3`, name, val, utils.NewBig(chainID))
	if err != nil {
		return errors.Wrapf(err, "failed to store chain config for chain ID %s", chainID.String())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.Wrapf(ErrNoRowsAffected, "no chain found with ID %s", chainID.String())
	}
	return nil
}

// Clear deletes a config value for the given chain and key
func (o *orm) Clear(chainID *big.Int, name string) error {
	res, err := o.db.Exec(`UPDATE evm_chains SET cfg = cfg - $1 WHERE id = $2`, name, utils.NewBig(chainID))
	if err != nil {
		return errors.Wrapf(err, "failed to clear chain config for chain ID %s", chainID.String())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.Wrapf(ErrNoRowsAffected, "no chain found with ID %s", chainID.String())
	}
	return nil
}
