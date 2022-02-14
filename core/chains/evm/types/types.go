package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type NewNode struct {
	Name       string      `json:"name"`
	EVMChainID utils.Big   `json:"evmChainId"`
	WSURL      null.String `json:"wsURL" db:"ws_url"`
	HTTPURL    null.String `json:"httpURL" db:"http_url"`
	SendOnly   bool        `json:"sendOnly"`
}

type ChainConfigORM interface {
	StoreString(chainID *big.Int, key, val string) error
	Clear(chainID *big.Int, key string) error
}

//go:generate mockery --name ORM --output ./../mocks/ --case=underscore
type ORM interface {
	EnabledChainsWithNodes() ([]Chain, error)
	Chain(id utils.Big) (chain Chain, err error)
	CreateChain(id utils.Big, config ChainCfg) (Chain, error)
	UpdateChain(id utils.Big, enabled bool, config ChainCfg) (Chain, error)
	DeleteChain(id utils.Big) error
	Chains(offset, limit int) ([]Chain, int, error)
	CreateNode(data NewNode) (Node, error)
	DeleteNode(id int64) error
	GetChainsByIDs(ids []utils.Big) (chains []Chain, err error)
	GetNodesByChainIDs(chainIDs []utils.Big) (nodes []Node, err error)
	Node(id int32) (Node, error)
	Nodes(offset, limit int) ([]Node, int, error)
	NodesForChain(chainID utils.Big, offset, limit int) ([]Node, int, error)
	ChainConfigORM
}

type ChainCfg struct {
	BlockHistoryEstimatorBlockDelay                null.Int
	BlockHistoryEstimatorBlockHistorySize          null.Int
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks null.Int
	ChainType                                      null.String
	EthTxReaperThreshold                           *models.Duration
	EthTxResendAfterThreshold                      *models.Duration
	EvmEIP1559DynamicFees                          null.Bool
	EvmFinalityDepth                               null.Int
	EvmGasBumpPercent                              null.Int
	EvmGasBumpTxDepth                              null.Int
	EvmGasBumpWei                                  *utils.Big
	EvmGasFeeCapDefault                            *utils.Big
	EvmGasLimitDefault                             null.Int
	EvmGasLimitMultiplier                          null.Float
	EvmGasPriceDefault                             *utils.Big
	EvmGasTipCapDefault                            *utils.Big
	EvmGasTipCapMinimum                            *utils.Big
	EvmHeadTrackerHistoryDepth                     null.Int
	EvmHeadTrackerMaxBufferSize                    null.Int
	EvmHeadTrackerSamplingInterval                 *models.Duration
	EvmLogBackfillBatchSize                        null.Int
	EvmMaxGasPriceWei                              *utils.Big
	EvmNonceAutoSync                               null.Bool
	EvmRPCDefaultBatchSize                         null.Int
	FlagsContractAddress                           null.String
	GasEstimatorMode                               null.String
	KeySpecific                                    map[string]ChainCfg
	LinkContractAddress                            null.String
	MinIncomingConfirmations                       null.Int
	MinRequiredOutgoingConfirmations               null.Int
	MinimumContractPayment                         *assets.Link
	OCRObservationTimeout                          *models.Duration
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

type Chain struct {
	ID        utils.Big
	Nodes     []Node
	Cfg       ChainCfg
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}
type Node struct {
	ID         int32
	Name       string
	EVMChainID utils.Big
	WSURL      null.String `db:"ws_url"`
	HTTPURL    null.String `db:"http_url"`
	SendOnly   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
