package terra

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	TerraChainID  string `json:"terraChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
}

type Chain struct {
	ID        string
	Cfg       ChainCfg
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

type Node struct {
	ID            int32
	Name          string
	TerraChainID  string
	TendermintURL string `db:"tendermint_url"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ChainCfg struct {
	BlockRate             *models.Duration
	BlocksUntilTxTimeout  null.Int
	ConfirmPollPeriod     *models.Duration
	FallbackGasPriceULuna null.String
	FCDURL                null.String `db:"fcd_url"`
	GasLimitMultiplier    null.Float
	MaxMsgsPerBatch       null.Int
	OCR2CachePollPeriod   *models.Duration
	OCR2CacheTTL          *models.Duration
}

func (c *ChainCfg) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, c)
}

func (c ChainCfg) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// State represents the state of a given terra msg
// Happy path: Unstarted->Broadcasted->Confirmed
type State string

var (
	// Unstarted means queued but not processed.
	// Valid next states: Broadcasted, Errored (sim fails)
	Unstarted State = "unstarted"
	// Broadcasted means included in the mempool of a node.
	// Valid next states: Confirmed (found onchain), Errored (tx expired waiting for confirmation)
	Broadcasted State = "broadcasted"
	// Confirmed means we're able to retrieve the txhash of the tx which broadcasted the msg.
	// Valid next states: none, terminal state
	Confirmed State = "confirmed"
	// Errored means the msg reverted in simulation OR the tx containing the message timed out waiting to be confirmed
	// TODO: when we add gas bumping, we'll address that timeout case
	// Valid next states, none, terminal state
	Errored State = "errored"
)

type Msg struct {
	ID         int64
	ChainID    string `db:"terra_chain_id"`
	ContractID string
	State      State
	Raw        []byte // serialized msg
	TxHash     *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
