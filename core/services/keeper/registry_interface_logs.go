package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	registry1_1 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry1_3 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
)

func (rw *RegistryWrapper) GetLogListenerOpts(minIncomingConfirmations uint32, upkeepPerformedFilter [][]log.Topic) (*log.ListenerOpts, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return &log.ListenerOpts{
			Contract: rw.contract1_1.Address(),
			ParseLog: rw.contract1_1.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				registry1_1.KeeperRegistryKeepersUpdated{}.Topic():   nil,
				registry1_1.KeeperRegistryConfigSet{}.Topic():        nil,
				registry1_1.KeeperRegistryUpkeepCanceled{}.Topic():   nil,
				registry1_1.KeeperRegistryUpkeepRegistered{}.Topic(): nil,
				registry1_1.KeeperRegistryUpkeepPerformed{}.Topic():  upkeepPerformedFilter,
			},
			MinIncomingConfirmations: minIncomingConfirmations,
		}, nil
	case RegistryVersion_1_2:
		return &log.ListenerOpts{
			Contract: rw.contract1_2.Address(),
			ParseLog: rw.contract1_2.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				registry1_2.KeeperRegistryKeepersUpdated{}.Topic():    nil,
				registry1_2.KeeperRegistryConfigSet{}.Topic():         nil,
				registry1_2.KeeperRegistryUpkeepCanceled{}.Topic():    nil,
				registry1_2.KeeperRegistryUpkeepRegistered{}.Topic():  nil,
				registry1_2.KeeperRegistryUpkeepPerformed{}.Topic():   upkeepPerformedFilter,
				registry1_2.KeeperRegistryUpkeepGasLimitSet{}.Topic(): nil,
				registry1_2.KeeperRegistryUpkeepMigrated{}.Topic():    nil,
				registry1_2.KeeperRegistryUpkeepReceived{}.Topic():    nil,
			},
			MinIncomingConfirmations: minIncomingConfirmations,
		}, nil
	case RegistryVersion_1_3:
		return &log.ListenerOpts{
			Contract: rw.contract1_3.Address(),
			ParseLog: rw.contract1_3.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				registry1_3.KeeperRegistryKeepersUpdated{}.Topic():         nil,
				registry1_3.KeeperRegistryConfigSet{}.Topic():              nil,
				registry1_3.KeeperRegistryUpkeepCanceled{}.Topic():         nil,
				registry1_3.KeeperRegistryUpkeepRegistered{}.Topic():       nil,
				registry1_3.KeeperRegistryUpkeepPerformed{}.Topic():        upkeepPerformedFilter,
				registry1_3.KeeperRegistryUpkeepGasLimitSet{}.Topic():      nil,
				registry1_3.KeeperRegistryUpkeepMigrated{}.Topic():         nil,
				registry1_3.KeeperRegistryUpkeepReceived{}.Topic():         nil,
				registry1_3.KeeperRegistryUpkeepPaused{}.Topic():           nil,
				registry1_3.KeeperRegistryUpkeepUnpaused{}.Topic():         nil,
				registry1_3.KeeperRegistryUpkeepCheckDataUpdated{}.Topic(): nil,
			},
			MinIncomingConfirmations: minIncomingConfirmations,
		}, nil
	default:
		return nil, newUnsupportedVersionError("GetLogListenerOpts", rw.Version)
	}
}

func (rw *RegistryWrapper) GetCancelledUpkeepIDFromLog(broadcast log.Broadcast) (*big.Int, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_1.KeeperRegistryUpkeepCanceled)
		if !ok {
			return nil, errors.Errorf("expected UpkeepCanceled log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepCanceled)
		if !ok {
			return nil, errors.Errorf("expected UpkeepCanceled log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepCanceled)
		if !ok {
			return nil, errors.Errorf("expected UpkeepCanceled log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetCancelledUpkeepIDFromLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetUpkeepIdFromRegistrationLog(broadcast log.Broadcast) (*big.Int, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_1.KeeperRegistryUpkeepRegistered)
		if !ok {
			return nil, errors.Errorf("expected UpkeepRegistered log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepRegistered)
		if !ok {
			return nil, errors.Errorf("expected UpkeepRegistered log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepRegistered)
		if !ok {
			return nil, errors.Errorf("expected UpkeepRegistered log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromRegistrationLog", rw.Version)
	}
}

type UpkeepPerformedLog struct {
	UpkeepID   *big.Int
	FromKeeper common.Address
}

func (rw *RegistryWrapper) ParseUpkeepPerformedLog(broadcast log.Broadcast) (*UpkeepPerformedLog, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_1.KeeperRegistryUpkeepPerformed)
		if !ok {
			return nil, errors.Errorf("expected UpkeepPerformed log but got %T", broadcastedLog)
		}
		return &UpkeepPerformedLog{
			UpkeepID:   broadcastedLog.Id,
			FromKeeper: broadcastedLog.From,
		}, nil
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepPerformed)
		if !ok {
			return nil, errors.Errorf("expected UpkeepPerformed log but got %T", broadcastedLog)
		}
		return &UpkeepPerformedLog{
			UpkeepID:   broadcastedLog.Id,
			FromKeeper: broadcastedLog.From,
		}, nil
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepPerformed)
		if !ok {
			return nil, errors.Errorf("expected UpkeepPerformed log but got %T", broadcastedLog)
		}
		return &UpkeepPerformedLog{
			UpkeepID:   broadcastedLog.Id,
			FromKeeper: broadcastedLog.From,
		}, nil
	default:
		return nil, newUnsupportedVersionError("ParseUpkeepPerformedLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetIDFromGasLimitSetLog(broadcast log.Broadcast) (*big.Int, error) {
	// Only supported on 1.2 and 1.3
	switch rw.Version {
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepGasLimitSet)
		if !ok {
			return nil, errors.Errorf("expected UpkeepGasLimitSetlog but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepGasLimitSet)
		if !ok {
			return nil, errors.Errorf("expected UpkeepGasLimitSetlog but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetIDFromGasLimitSetLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetUpkeepIdFromReceivedLog(broadcast log.Broadcast) (*big.Int, error) {
	// Only supported on 1.2 and 1.3
	switch rw.Version {
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepReceived)
		if !ok {
			return nil, errors.Errorf("expected UpkeepReceived log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepReceived)
		if !ok {
			return nil, errors.Errorf("expected UpkeepReceived log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromReceivedLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetUpkeepIdFromMigratedLog(broadcast log.Broadcast) (*big.Int, error) {
	// Only supported on 1.2 and 1.3
	switch rw.Version {
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepMigrated)
		if !ok {
			return nil, errors.Errorf("expected UpkeepMigrated log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepMigrated)
		if !ok {
			return nil, errors.Errorf("expected UpkeepMigrated log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromMigratedLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetUpkeepIdFromUpkeepPausedLog(broadcast log.Broadcast) (*big.Int, error) {
	// Only supported on 1.3
	switch rw.Version {
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepPaused)
		if !ok {
			return nil, errors.Errorf("expected UpkeepPaused log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromUpkeepPausedLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetUpkeepIdFromUpkeepUnpausedLog(broadcast log.Broadcast) (*big.Int, error) {
	// Only supported on 1.3
	switch rw.Version {
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepUnpaused)
		if !ok {
			return nil, errors.Errorf("expected UpkeepUnpaused log but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromUpkeepUnpausedLog", rw.Version)
	}
}

type UpkeepCheckDataUpdatedLog struct {
	UpkeepID     *big.Int
	NewCheckData []byte
}

func (rw *RegistryWrapper) ParseUpkeepCheckDataUpdatedLog(broadcast log.Broadcast) (*UpkeepCheckDataUpdatedLog, error) {
	// Only supported on 1.3
	switch rw.Version {
	case RegistryVersion_1_3:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_3.KeeperRegistryUpkeepCheckDataUpdated)
		if !ok {
			return nil, errors.Errorf("expected UpkeepCheckDataUpdated log but got %T", broadcastedLog)
		}
		return &UpkeepCheckDataUpdatedLog{
			UpkeepID:     broadcastedLog.Id,
			NewCheckData: broadcastedLog.NewCheckData,
		}, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeepIdFromUpkeepPausedLog", rw.Version)
	}
}
