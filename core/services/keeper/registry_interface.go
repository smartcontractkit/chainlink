package keeper

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	registry1_1 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry2_0 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	type_and_version "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/type_and_version_interface_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type RegistryVersion int32

const (
	RegistryVersion_1_0 RegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
	RegistryVersion_2_0
)

func (rv RegistryVersion) String() string {
	switch rv {
	case RegistryVersion_1_0, RegistryVersion_1_1, RegistryVersion_1_2:
		return fmt.Sprintf("v1.%d", rv)
	case RegistryVersion_2_0:
		return "v2.0"
	default:
		return "unknown registry version"
	}
}

const ActiveUpkeepIDBatchSize int64 = 10000

// upkeepGetter is declared as a private interface as it is only needed
// internally to the keeper package for now
type upkeepGetter interface {
	GetUpkeep(*bind.CallOpts, *big.Int) (*UpkeepConfig, error)
}

// RegistryWrapper implements a layer on top of different versions of registry wrappers
// to provide a unified layer to rest of the codebase
type RegistryWrapper struct {
	Address     ethkey.EIP55Address
	Version     RegistryVersion
	contract1_1 *registry1_1.KeeperRegistry
	contract1_2 *registry1_2.KeeperRegistry
	contract2_0 *registry2_0.KeeperRegistry
	evmClient   bind.ContractBackend
}

func NewRegistryWrapper(address ethkey.EIP55Address, backend bind.ContractBackend) (*RegistryWrapper, error) {
	interface_wrapper, err := type_and_version.NewTypeAndVersionInterface(
		address.Address(),
		backend,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create type and interface wrapper")
	}

	version, err := getRegistryVersion(interface_wrapper)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine version of keeper registry contract")
	}

	contract1_1, err := registry1_1.NewKeeperRegistry(
		address.Address(),
		backend,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry 1_1 contract wrapper")
	}

	contract1_2, err := registry1_2.NewKeeperRegistry(
		address.Address(),
		backend,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry 1_2 contract wrapper")
	}

	contract2_0, err := registry2_0.NewKeeperRegistry(
		address.Address(),
		backend,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry 2_0 contract wrapper")
	}

	return &RegistryWrapper{
		Address:     address,
		Version:     *version,
		contract1_1: contract1_1,
		contract1_2: contract1_2,
		contract2_0: contract2_0,
		evmClient:   backend,
	}, nil
}

func getRegistryVersion(contract *type_and_version.TypeAndVersionInterface) (*RegistryVersion, error) {
	typeAndVersion, err := contract.TypeAndVersion(nil)
	if err != nil {
		jsonErr := evmclient.ExtractRPCError(err)
		if jsonErr != nil {
			// Version 1.0 does not support typeAndVersion interface, hence gives a json error on this call
			version := RegistryVersion_1_0
			return &version, nil
		}
		return nil, errors.Wrap(err, "unable to fetch version of registry")
	}

	switch {
	case strings.HasPrefix(typeAndVersion, "KeeperRegistry 1.1"):
		version := RegistryVersion_1_1
		return &version, nil
	case strings.HasPrefix(typeAndVersion, "KeeperRegistry 1.2"):
		version := RegistryVersion_1_2
		return &version, nil
	case strings.HasPrefix(typeAndVersion, "KeeperRegistry 2.0"):
		version := RegistryVersion_2_0
		return &version, nil
	default:
		return nil, errors.Errorf("Registry type and version %s not supported", typeAndVersion)
	}
}

func newUnsupportedVersionError(functionName string, version RegistryVersion) error {
	return errors.Errorf("Registry version %d does not support %s", version, functionName)
}

func (rw *RegistryWrapper) GetActiveUpkeepIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		upkeepCount, err := rw.contract1_1.GetUpkeepCount(opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get upkeep count")
		}
		cancelledUpkeeps, err := rw.contract1_1.GetCanceledUpkeepList(opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get cancelled upkeeps")
		}
		cancelledSet := make(map[int64]bool)
		for _, upkeepID := range cancelledUpkeeps {
			cancelledSet[upkeepID.Int64()] = true
		}
		// Active upkeep IDs are 0,1 ... upkeepCount-1, removing the cancelled ones
		activeUpkeeps := make([]*big.Int, 0)
		for i := int64(0); i < upkeepCount.Int64(); i++ {
			if _, found := cancelledSet[i]; !found {
				activeUpkeeps = append(activeUpkeeps, big.NewInt(i))
			}
		}
		return activeUpkeeps, nil
	case RegistryVersion_1_2:
		if opts == nil || opts.BlockNumber.Int64() == 0 {
			// fetch the current block number so batched GetActiveUpkeepIDs calls can be performed on the same block
			header, err := rw.evmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch EVM block header")
			}
			if opts != nil {
				opts.BlockNumber = header.Number
			} else {
				opts = &bind.CallOpts{
					BlockNumber: header.Number,
				}
			}
		}

		state, err := rw.contract1_2.GetState(opts)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get contract state at block number %d", opts.BlockNumber.Int64())
		}

		activeUpkeepIDs := make([]*big.Int, 0)
		for int64(len(activeUpkeepIDs)) < state.State.NumUpkeeps.Int64() {
			startIndex := int64(len(activeUpkeepIDs))
			maxCount := state.State.NumUpkeeps.Int64() - int64(len(activeUpkeepIDs))
			if maxCount > ActiveUpkeepIDBatchSize {
				maxCount = ActiveUpkeepIDBatchSize
			}
			activeUpkeepIDBatch, err := rw.contract1_2.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get active upkeep IDs from index %d to %d (both inclusive)", startIndex, startIndex+maxCount-1)
			}
			activeUpkeepIDs = append(activeUpkeepIDs, activeUpkeepIDBatch...)
		}

		return activeUpkeepIDs, nil
	case RegistryVersion_2_0:
		if opts == nil || opts.BlockNumber.Int64() == 0 {
			// fetch the current block number so batched GetActiveUpkeepIDs calls can be performed on the same block
			header, err := rw.evmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch EVM block header")
			}
			if opts != nil {
				opts.BlockNumber = header.Number
			} else {
				opts = &bind.CallOpts{
					BlockNumber: header.Number,
				}
			}
		}

		state, err := rw.contract2_0.GetState(opts)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get contract state at block number %d", opts.BlockNumber.Int64())
		}

		activeUpkeepIDs := make([]*big.Int, 0)
		for int64(len(activeUpkeepIDs)) < state.State.NumUpkeeps.Int64() {
			startIndex := int64(len(activeUpkeepIDs))
			maxCount := state.State.NumUpkeeps.Int64() - int64(len(activeUpkeepIDs))
			if maxCount > ActiveUpkeepIDBatchSize {
				maxCount = ActiveUpkeepIDBatchSize
			}
			activeUpkeepIDBatch, err := rw.contract1_2.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get active upkeep IDs from index %d to %d (both inclusive)", startIndex, startIndex+maxCount-1)
			}
			activeUpkeepIDs = append(activeUpkeepIDs, activeUpkeepIDBatch...)
		}

		return activeUpkeepIDs, nil
	default:
		return nil, newUnsupportedVersionError("GetActiveUpkeepIDs", rw.Version)
	}
}

type UpkeepConfig struct {
	ExecuteGas uint32
	CheckData  []byte
	LastKeeper common.Address
}

func (rw *RegistryWrapper) GetUpkeep(opts *bind.CallOpts, id *big.Int) (*UpkeepConfig, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		upkeep, err := rw.contract1_1.GetUpkeep(opts, id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get upkeep config")
		}
		return &UpkeepConfig{
			ExecuteGas: upkeep.ExecuteGas,
			CheckData:  upkeep.CheckData,
			LastKeeper: upkeep.LastKeeper,
		}, nil
	case RegistryVersion_1_2:
		upkeep, err := rw.contract1_2.GetUpkeep(opts, id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get upkeep config")
		}
		return &UpkeepConfig{
			ExecuteGas: upkeep.ExecuteGas,
			CheckData:  upkeep.CheckData,
			LastKeeper: upkeep.LastKeeper,
		}, nil
	case RegistryVersion_2_0:
		upkeep, err := rw.contract2_0.GetUpkeep(opts, id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get upkeep config")
		}
		return &UpkeepConfig{
			ExecuteGas: upkeep.ExecuteGas,
			CheckData:  upkeep.CheckData,
			LastKeeper: utils.ZeroAddress, // TODO: ...
		}, nil
	default:
		return nil, newUnsupportedVersionError("GetUpkeep", rw.Version)
	}
}

type RegistryConfig struct {
	BlockCountPerTurn int32
	CheckGas          uint32
	KeeperAddresses   []common.Address
}

func (rw *RegistryWrapper) GetConfig(opts *bind.CallOpts) (*RegistryConfig, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		config, err := rw.contract1_1.GetConfig(opts)
		if err != nil {
			// TODO: error wrapping with %w should be done here to preserve the error type as it bubbles up
			// pkg/errors doesn't support the native errors.Is/As capabilities
			// using pkg/errors produces a stack trace in the logs and this behavior is too valuable to let go
			return nil, errors.Errorf("%s [%s]: getConfig %s", ErrContractCallFailure, err, rw.Version)
		}

		keeperAddresses, err := rw.contract1_1.GetKeeperList(opts)
		if err != nil {
			return nil, errors.Errorf("%s [%s]: getKeeperList %s", ErrContractCallFailure, err, rw.Version)
		}

		return &RegistryConfig{
			BlockCountPerTurn: int32(config.BlockCountPerTurn.Int64()),
			CheckGas:          config.CheckGasLimit,
			KeeperAddresses:   keeperAddresses,
		}, nil
	case RegistryVersion_1_2:
		state, err := rw.contract1_2.GetState(opts)
		if err != nil {
			return nil, errors.Errorf("%s [%s]: getState %s", ErrContractCallFailure, err, rw.Version)
		}

		return &RegistryConfig{
			BlockCountPerTurn: int32(state.Config.BlockCountPerTurn.Int64()),
			CheckGas:          state.Config.CheckGasLimit,
			KeeperAddresses:   state.Keepers,
		}, nil
	case RegistryVersion_2_0:
		state, err := rw.contract2_0.GetState(opts)
		if err != nil {
			return nil, errors.Errorf("%s [%s]: getState %s", ErrContractCallFailure, err, rw.Version)
		}

		return &RegistryConfig{
			BlockCountPerTurn: 1, // TODO: ...
			CheckGas:          state.Config.CheckGasLimit,
			KeeperAddresses:   state.Transmitters,
		}, nil
	default:
		return nil, newUnsupportedVersionError("GetConfig", rw.Version)
	}
}

func (rw *RegistryWrapper) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.SetKeepers(opts, keepers, payees)
	case RegistryVersion_1_2:
		return rw.contract1_2.SetKeepers(opts, keepers, payees)
	case RegistryVersion_2_0:
		panic("not implemented")
	default:
		return nil, newUnsupportedVersionError("SetKeepers", rw.Version)
	}
}

func (rw *RegistryWrapper) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.RegisterUpkeep(opts, target, gasLimit, admin, checkData)
	case RegistryVersion_1_2:
		return rw.contract1_2.RegisterUpkeep(opts, target, gasLimit, admin, checkData)
	case RegistryVersion_2_0:
		// TODO: Don't skip sig verification
		return rw.contract2_0.RegisterUpkeep(opts, target, gasLimit, admin, true, checkData)
	default:
		return nil, newUnsupportedVersionError("RegisterUpkeep", rw.Version)
	}
}

func (rw *RegistryWrapper) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.AddFunds(opts, id, amount)
	case RegistryVersion_1_2:
		return rw.contract1_2.AddFunds(opts, id, amount)
	case RegistryVersion_2_0:
		return rw.contract2_0.AddFunds(opts, id, amount)
	default:
		return nil, newUnsupportedVersionError("AddFunds", rw.Version)
	}
}

func (rw *RegistryWrapper) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.PerformUpkeep(opts, id, performData)
	case RegistryVersion_1_2:
		return rw.contract1_2.PerformUpkeep(opts, id, performData)
	case RegistryVersion_2_0:
		panic("does not support") // TODO: ...
	default:
		return nil, newUnsupportedVersionError("PerformUpkeep", rw.Version)
	}
}

func (rw *RegistryWrapper) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.CancelUpkeep(opts, id)
	case RegistryVersion_1_2:
		return rw.contract1_2.CancelUpkeep(opts, id)
	case RegistryVersion_2_0:
		return rw.contract2_0.CancelUpkeep(opts, id)
	default:
		return nil, newUnsupportedVersionError("CancelUpkeep", rw.Version)
	}
}
