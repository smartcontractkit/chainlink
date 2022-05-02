package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	registry1_1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
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
	default:
		return nil, newUnsupportedVersionError("ParseUpkeepPerformedLog", rw.Version)
	}
}

func (rw *RegistryWrapper) GetIDFromGasLimitSetLog(broadcast log.Broadcast) (*big.Int, error) {
	// Only supported on 1.2
	switch rw.Version {
	case RegistryVersion_1_2:
		broadcastedLog, ok := broadcast.DecodedLog().(*registry1_2.KeeperRegistryUpkeepGasLimitSet)
		if !ok {
			return nil, errors.Errorf("expected UpkeepGasLimitSetlog but got %T", broadcastedLog)
		}
		return broadcastedLog.Id, nil
	default:
		return nil, newUnsupportedVersionError("GetIDFromGasLimitSetLog", rw.Version)
	}
}
