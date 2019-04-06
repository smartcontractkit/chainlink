package old

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// IndexableBlockNumber captured in old state for tests.
type IndexableBlockNumber struct {
	Number models.Big  `json:"number" gorm:"index;type:varchar(255);not null"`
	Digits int         `json:"digits" gorm:"index"`
	Hash   common.Hash `json:"hash"`
}
