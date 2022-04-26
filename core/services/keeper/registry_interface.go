package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	registry1_1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

type RegistryVersion int32

const (
	RegistryVersion_1_0 RegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
)

// RegistryWrapper implements a layer on top of different versions of registry wrappers
// to provide a unified layer to rest of the codebase
type RegistryWrapper struct {
	Address     ethkey.EIP55Address
	Version     RegistryVersion
	contract1_1 *registry1_1.KeeperRegistry
	contract1_2 *registry1_2.KeeperRegistry
}

func NewRegistryWrapper(address ethkey.EIP55Address, backend bind.ContractBackend) (*RegistryWrapper, error) {
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

	version, err := getRegistryVersion(contract1_1)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine version of keeper registry contract")
	}

	return &RegistryWrapper{
		Address:     address,
		Version:     *version,
		contract1_1: contract1_1,
		contract1_2: contract1_2,
	}, nil
}

func getRegistryVersion(contract1_1 *registry1_1.KeeperRegistry) (*RegistryVersion, error) {
	// Use registry 1_1 wrapper to get version information
	typeAndVersion, err := contract1_1.TypeAndVersion(nil)
	if err != nil {
		jsonErr := evmclient.ExtractRPCError(err)
		if jsonErr != nil {
			// Version 1.0 does not support typeAndVersion interface, hence gives a json error on this call
			version := RegistryVersion_1_0
			return &version, nil
		}
		return nil, errors.Wrap(err, "unable to fetch version of registry")
	}
	switch typeAndVersion {
	case "KeeperRegistry 1.1.0":
		version := RegistryVersion_1_1
		return &version, nil
	case "KeeperRegistry 1.2.0":
		version := RegistryVersion_1_2
		return &version, nil
	default:
		return nil, errors.Errorf("Registry type and version %s not supported", typeAndVersion)
	}
}

/*
 * This returns the numebr of registered upkeeps on contract, which can be different
 * from number of active upkeeps as some might be cancelled
 */
func (rw *RegistryWrapper) GetRegisteredUpkeepCount(opts *bind.CallOpts) (*big.Int, error) {
	// only applicable on 1.0 and 1.1.
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.GetUpkeepCount(opts)
	default:
		return nil, errors.Errorf("Registry version %d does not support GetUpkeepCount", rw.Version)
	}
}

func (rw *RegistryWrapper) GetCanceledUpkeepList(opts *bind.CallOpts) ([]*big.Int, error) {
	// only applicable on 1.0 and 1.1
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.GetCanceledUpkeepList(opts)
	default:
		return *new([]*big.Int), errors.Errorf("Registry version %d does not support GetCanceledUpkeepList", rw.Version)
	}
}

func (rw *RegistryWrapper) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	// only applicable on 1.2
	switch rw.Version {
	case RegistryVersion_1_2:
		return rw.contract1_2.GetActiveUpkeepIDs(opts, startIndex, maxCount)
	default:
		return *new([]*big.Int), errors.Errorf("Registry version %d does not support GetActiveUpkeepIDs", rw.Version)
	}
}

type UpkeepConfig struct {
	ExecuteGas uint32
	CheckData  []byte
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
		}, nil
	case RegistryVersion_1_2:
		upkeep, err := rw.contract1_2.GetUpkeep(opts, id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get upkeep config")
		}
		return &UpkeepConfig{
			ExecuteGas: upkeep.ExecuteGas,
			CheckData:  upkeep.CheckData,
		}, nil
	default:
		return nil, errors.Errorf("Registry version %d does not support GetUpkeep", rw.Version)
	}
}

type RegistryConfig struct {
	BlockCountPerTurn int32
	CheckGas          int32
	KeeperAddresses   []common.Address
}

func (rw *RegistryWrapper) GetConfig(opts *bind.CallOpts) (*RegistryConfig, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		config, err := rw.contract1_1.GetConfig(opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get contract config")
		}
		keeperAddresses, err := rw.contract1_1.GetKeeperList(nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keeper list")
		}
		return &RegistryConfig{
			BlockCountPerTurn: int32(config.BlockCountPerTurn.Int64()),
			CheckGas:          int32(config.CheckGasLimit),
			KeeperAddresses:   keeperAddresses,
		}, nil
	case RegistryVersion_1_2:
		state, err := rw.contract1_2.GetState(nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get contract state")
		}

		return &RegistryConfig{
			BlockCountPerTurn: int32(state.Config.BlockCountPerTurn.Int64()),
			CheckGas:          int32(state.Config.CheckGasLimit),
			KeeperAddresses:   state.Keepers,
		}, nil
	default:
		return nil, errors.Errorf("Registry version %d does not support GetConfig", rw.Version)
	}
}

func (rw *RegistryWrapper) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.SetKeepers(opts, keepers, payees)
	case RegistryVersion_1_2:
		return rw.contract1_2.SetKeepers(opts, keepers, payees)
	default:
		return nil, errors.Errorf("Registry version %d does not support SetKeepers", rw.Version)
	}
}

func (rw *RegistryWrapper) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.RegisterUpkeep(opts, target, gasLimit, admin, checkData)
	case RegistryVersion_1_2:
		return rw.contract1_2.RegisterUpkeep(opts, target, gasLimit, admin, checkData)
	default:
		return nil, errors.Errorf("Registry version %d does not support SetKeepers", rw.Version)
	}
}

func (rw *RegistryWrapper) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.AddFunds(opts, id, amount)
	case RegistryVersion_1_2:
		return rw.contract1_2.AddFunds(opts, id, amount)
	default:
		return nil, errors.Errorf("Registry version %d does not support SetKeepers", rw.Version)
	}
}

func (rw *RegistryWrapper) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.PerformUpkeep(opts, id, performData)
	case RegistryVersion_1_2:
		return rw.contract1_2.PerformUpkeep(opts, id, performData)
	default:
		return nil, errors.Errorf("Registry version %d does not support SetKeepers", rw.Version)
	}
}

func (rw *RegistryWrapper) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	switch rw.Version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rw.contract1_1.CancelUpkeep(opts, id)
	case RegistryVersion_1_2:
		return rw.contract1_2.CancelUpkeep(opts, id)
	default:
		return nil, errors.Errorf("Registry version %d does not support SetKeepers", rw.Version)
	}
}
