package keeper

import "github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"

type Registry struct {
	ID                int64
	BlockCountPerTurn int32
	CheckGas          int32
	ContractAddress   ethkey.EIP55Address
	FromAddress       ethkey.EIP55Address
	JobID             int32
	KeeperIndex       int32
	NumKeepers        int32
}

func (Registry) TableName() string {
	return "keeper_registries"
}

type UpkeepRegistration struct {
	ID                  int32
	CheckData           []byte
	ExecuteGas          uint64
	LastRunBlockHeight  int64
	RegistryID          int64
	Registry            Registry
	UpkeepID            int64
	PositioningConstant int32
}
