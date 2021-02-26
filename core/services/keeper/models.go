package keeper

import "github.com/smartcontractkit/chainlink/core/store/models"

type Registry struct {
	ID                int32 `gorm:"primary_key"`
	ContractAddress   models.EIP55Address
	BlockCountPerTurn int32
	CheckGas          int32
	FromAddress       models.EIP55Address
	JobID             int32 `gorm:"default:null"`
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

// func (reg Registry) SyncFromContract(contract *keeper_registry_contract.KeeperRegistryContract) (Registry, error) {
// 	config, err := contract.GetConfig(nil)
// 	if err != nil {
// 		return Registry{}, err
// 	}
// 	reg.CheckGas = config.CheckGasLimit
// 	reg.BlockCountPerTurn = uint32(config.BlockCountPerTurn.Uint64())
// 	keeperAddresses, err := contract.GetKeeperList(nil)
// 	if err != nil {
// 		return Registry{}, err
// 	}
// 	found := false
// 	for idx, address := range keeperAddresses {
// 		if address == reg.From {
// 			reg.KeeperIndex = uint32(idx)
// 			found = true
// 		}
// 	}
// 	if !found {
// 		return Registry{}, fmt.Errorf("unable to find %s in keeper list on registry %s", reg.From.Hex(), reg.Address.Hex())
// 	}

// 	reg.NumKeepers = uint32(len(keeperAddresses))

// 	return reg, nil
// }

type UpkeepRegistration struct {
	ID                  int32 `gorm:"primary_key"`
	CheckData           []byte
	ExecuteGas          int32
	RegistryID          int32
	Registry            Registry `gorm:"association_autoupdate:false"`
	UpkeepID            int64
	PositioningConstant int32
}
