package keeper

import "github.com/smartcontractkit/chainlink/core/store/models"

type Registry struct {
	ID                int32 `gorm:"primary_key"`
	BlockCountPerTurn int32
	CheckGas          int32
	ContractAddress   models.EIP55Address
	FromAddress       models.EIP55Address
	JobID             int32
	KeeperIndex       int32
	NumKeepers        int32
}

func NewRegistry(address models.EIP55Address, from models.EIP55Address, jobID int32) Registry {
	return Registry{
		ContractAddress: address,
		FromAddress:     from,
		JobID:           jobID,
	}
}

func (Registry) TableName() string {
	return "keeper_registries"
}

// todo - upkeep
type UpkeepRegistration struct {
	ID                  int32 `gorm:"primary_key"`
	CheckData           []byte
	ExecuteGas          int32
	RegistryID          int32
	Registry            Registry
	UpkeepID            int64
	PositioningConstant int32
}
