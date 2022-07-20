package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (rs *RegistrySynchronizer) ExportedFullSync() {
	rs.fullSync()
}

func (rs *RegistrySynchronizer) ExportedProcessLogs() {
	rs.processLogs()
}

func (rw *RegistryWrapper) GetUpkeepIdFromRawRegistrationLog(rawLog types.Log) (*big.Int, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		parsedLog, err := rw.contract1_1.ParseUpkeepRegistered(rawLog)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get parse UpkeepRegistered log")
		}
		return parsedLog.Id, nil
	case RegistryVersion_1_2:
		parsedLog, err := rw.contract1_2.ParseUpkeepRegistered(rawLog)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get parse UpkeepRegistered log")
		}
		return parsedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromRawRegistrationLog", rw.Version)
	}
}
