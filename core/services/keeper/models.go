package keeper

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type KeeperIndexMap map[ethkey.EIP55Address]int32

type Registry struct {
	ID                int64
	BlockCountPerTurn int32
	CheckGas          int32
	ContractAddress   ethkey.EIP55Address
	FromAddress       ethkey.EIP55Address
	JobID             int32
	KeeperIndex       int32
	NumKeepers        int32
	KeeperIndexMap    KeeperIndexMap
}

type UpkeepRegistration struct {
	ID                  int32
	CheckData           []byte
	ExecuteGas          uint64
	LastRunBlockHeight  int64
	RegistryID          int64
	Registry            Registry
	UpkeepID            *utils.Big
	LastKeeperIndex     null.Int64
	PositioningConstant int32
}

func (k *KeeperIndexMap) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &k)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &k)
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (k *KeeperIndexMap) Value() (driver.Value, error) {
	return json.Marshal(&k)
}

func (upkeep UpkeepRegistration) PrettyID() string {
	result, err := utils.Uint256ToBytes(upkeep.UpkeepID.ToInt())
	if err != nil {
		panic(errors.Wrap(err, "invariant, invalid upkeepID"))
	}
	return "UPx" + hex.EncodeToString(result)
}
