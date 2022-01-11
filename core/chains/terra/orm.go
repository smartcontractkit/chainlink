package terra

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
)

type orm struct {
	db *sqlx.DB
}

var _ types.ORM = (*orm)(nil)

// NewORM returns an ORM backed by db.
func NewORM(db *sqlx.DB) types.ORM {
	return &orm{db}
}

// ErrNoRowsAffected is returned when rows should have been affected but were not.
var ErrNoRowsAffected = errors.New("no rows affected")

func (o orm) CreateNode(data types.NewNode) (node types.Node, err error) {
	sql := `INSERT INTO terra_nodes (name, terra_chain_id, tendermint_url, fcd_url, created_at, updated_at)
	VALUES (:name, :terra_chain_id, :tendermint_url, :fcd_url, now(), now())
	RETURNING *;`
	stmt, err := o.db.PrepareNamed(sql)
	if err != nil {
		return node, err
	}
	err = stmt.Get(&node, data)
	return node, err
}

func (o orm) DeleteNode(id int32) error {
	sql := `DELETE FROM terra_nodes WHERE id = $1`
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

func (o orm) Node(id int32) (node types.Node, err error) {
	err = o.db.Get(&node, "SELECT * FROM terra_nodes WHERE id = $1;", id)

	return
}

func (o *orm) Nodes(offset, limit int) (nodes []types.Node, count int, err error) {
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM terra_nodes"); err != nil {
		return
	}

	sql := `SELECT * FROM terra_nodes ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = o.db.Select(&nodes, sql, limit, offset); err != nil {
		return
	}

	return
}

func (o *orm) NodesForChain(chainID string, offset, limit int) (nodes []types.Node, count int, err error) {
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM terra_nodes WHERE terra_chain_id = $1", chainID); err != nil {
		return
	}

	sql := `SELECT * FROM terra_nodes WHERE terra_chain_id = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3;`
	if err = o.db.Select(&nodes, sql, chainID, limit, offset); err != nil {
		return
	}

	return
}
