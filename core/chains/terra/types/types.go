package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// ORM manages terra chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (Chain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]Chain, int, error)
	CreateChain(id string, config ChainCfg, qopts ...pg.QOpt) (Chain, error)
	UpdateChain(id string, enabled bool, config ChainCfg, qopts ...pg.QOpt) (Chain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error

	// EnabledChainsWithNodes returns enabled chains with nodes (if any) included.
	EnabledChainsWithNodes(...pg.QOpt) ([]Chain, error)

	CreateNode(NewNode, ...pg.QOpt) (Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	Node(int32, ...pg.QOpt) (Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []Node, count int, err error)
}

// ChainCfg is configuration parameters for a terra chain.
type ChainCfg struct {
	ConfirmMaxPolls       null.Int
	ConfirmPollPeriod     *models.Duration
	FallbackGasPriceULuna null.String
	GasLimitMultiplier    null.Float
}

// Scan deserializes JSON from the database.
func (c *ChainCfg) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, c)
}

// Value serializes JSON for the database.
func (c ChainCfg) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	TerraChainID  string `json:"terraChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
	FCDURL        string `json:"fcdURL" db:"fcd_url"`
}

// Node is an existing node.
type Node struct {
	ID            int32
	Name          string
	TerraChainID  string
	TendermintURL string `db:"tendermint_url"`
	FCDURL        string `db:"fcd_url"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Chain is a an existing chain.
type Chain struct {
	ID        string
	Nodes     []Node
	Cfg       ChainCfg
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

func (Chain) TableName() string {
	return "terra_chains"
}
