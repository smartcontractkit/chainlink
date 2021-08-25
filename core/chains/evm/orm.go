package evm

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"
	null "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ORM interface {
	CreateChain(id utils.Big, config types.ChainCfg) (types.Chain, error)
	DeleteChain(id utils.Big) error
	Chains(offset, limit int) ([]types.Chain, int, error)
	CreateNode(data NewNode) (types.Node, error)
	DeleteNode(id int64) error
	Nodes(offset, limit int) ([]types.Node, int, error)
	NodesForChain(chainID utils.Big, offset, limit int) ([]types.Node, int, error)
}

type orm struct {
	db *sqlx.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB) ORM {
	return &orm{db}
}

var ErrNoRowsAffected = errors.New("no rows affected")

func (o *orm) CreateChain(id utils.Big, config types.ChainCfg) (chain types.Chain, err error) {
	sql := `INSERT INTO evm_chains (id, cfg, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	err = o.db.Get(&chain, sql, id, config)
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

type NewNode struct {
	Name       string      `json:"name"`
	EVMChainID utils.Big   `json:"evmChainId"`
	WSURL      null.String `json:"wsURL" db:"ws_url"`
	HTTPURL    string      `json:"httpURL" db:"http_url"`
	SendOnly   bool        `json:"sendOnly"`
}

func (o *orm) CreateNode(data NewNode) (node types.Node, err error) {
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
