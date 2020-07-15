package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type OracleRequest struct {
	RequestID          common.Hash `gorm:"primary_key"`
	SpecID             uuid.UUID
	Requester          common.Address
	Payment            assets.Link
	CallbackAddr       common.Address
	CallbackFunctionID FunctionSelector
	CancelExpiration   time.Time
	DataVersion        utils.Big
	Data               []byte
	CreatedAt          time.Time
}
