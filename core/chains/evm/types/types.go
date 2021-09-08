package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	null "gopkg.in/guregu/null.v4"
)

type ChainCfg struct {
	BlockHistoryEstimatorBlockDelay       null.Int
	BlockHistoryEstimatorBlockHistorySize null.Int
	EthTxResendAfterThreshold             *models.Duration
	EvmFinalityDepth                      null.Int
	EvmGasBumpPercent                     null.Int
	EvmGasBumpTxDepth                     null.Int
	EvmGasBumpWei                         *utils.Big
	EvmGasLimitDefault                    null.Int
	EvmGasLimitMultiplier                 null.Float
	EvmGasPriceDefault                    *utils.Big
	EvmHeadTrackerHistoryDepth            null.Int
	EvmHeadTrackerMaxBufferSize           null.Int
	EvmHeadTrackerSamplingInterval        *models.Duration
	EvmLogBackfillBatchSize               null.Int
	EvmMaxGasPriceWei                     *utils.Big
	EvmNonceAutoSync                      null.Bool
	EvmRPCDefaultBatchSize                null.Int
	FlagsContractAddress                  null.String
	GasEstimatorMode                      null.String
	MinRequiredOutgoingConfirmations      null.Int
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
	ID        utils.Big `gorm:"primary_key"`
	Nodes     []Node    `gorm:"->;foreignKey:EVMChainID;references:ID"`
	Cfg       ChainCfg
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Chain) TableName() string {
	return "evm_chains"
}
func (c Chain) IsL2() bool {
	return IsL2(c.ID.ToInt())
}
func (c Chain) IsArbitrum() bool {
	return IsArbitrum(c.ID.ToInt())
}
func (c Chain) IsOptimism() bool {
	return IsOptimism(c.ID.ToInt())
}

// IsArbitrum returns true if the chain is arbitrum mainnet or testnet
func IsArbitrum(id *big.Int) bool {
	return id.Cmp(big.NewInt(42161)) == 0 || id.Cmp(big.NewInt(421611)) == 0
}

// IsOptimism returns true if the chain is optimism mainnet or testnet
func IsOptimism(id *big.Int) bool {
	return id.Cmp(big.NewInt(10)) == 0 || id.Cmp(big.NewInt(69)) == 0
}

// IsL2 returns true if this chain is an L2 chain, notably that the block
// numbers used for log searching are different from calling block.number
func IsL2(id *big.Int) bool {
	return IsOptimism(id) || IsArbitrum(id)
}

func IsExChain(id *big.Int) bool {
	return id.Cmp(big.NewInt(65)) == 0 || id.Cmp(big.NewInt(66)) == 0
}

type Node struct {
	ID         int64 `gorm:"primary_key"`
	Name       string
	EVMChain   Chain
	EVMChainID utils.Big   `gorm:"column:evm_chain_id"`
	WSURL      null.String `gorm:"column:ws_url" db:"ws_url"`
	HTTPURL    string      `gorm:"column:http_url" db:"http_url"`
	SendOnly   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
