package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	cltypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	registrylogicc23 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_c_wrapper_2_3"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arbitrum_module"
	acutils "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_consumer_benchmark"
	automationForwarderLogic "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_forwarder_logic"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	registrar23 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_3"
	registrylogica22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_a_wrapper_2_2"
	registrylogicb22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_b_wrapper_2_2"
	registry22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_wrapper_2_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_module_base"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	iregistry22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"

	registrylogica23 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_a_wrapper_2_3"
	registrylogicb23 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_b_wrapper_2_3"
	registry23 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_wrapper_2_3"
	iregistry23 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_chain_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_consumer_performance_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic2_0"
	registrylogica21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_a_wrapper_2_1"
	registrylogicb21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	registry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_triggered_streams_lookup_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/optimism_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/perform_data_checker_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/scroll_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/streams_lookup_upkeep_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_transcoder"
)

// EthereumUpkeepTranscoder represents the transcoder which is used to perform migrations
// of upkeeps from one registry to another.
type EthereumUpkeepTranscoder struct {
	client     *seth.Client
	transcoder *upkeep_transcoder.UpkeepTranscoder
	address    *common.Address
}

func (v *EthereumUpkeepTranscoder) Address() string {
	return v.address.Hex()
}

func DeployUpkeepTranscoder(client *seth.Client) (*EthereumUpkeepTranscoder, error) {
	abi, err := upkeep_transcoder.UpkeepTranscoderMetaData.GetAbi()
	if err != nil {
		return &EthereumUpkeepTranscoder{}, fmt.Errorf("failed to get UpkeepTranscoder ABI: %w", err)
	}
	transcoderDeploymentData, err := client.DeployContract(client.NewTXOpts(), "UpkeepTranscoder", *abi, common.FromHex(upkeep_transcoder.UpkeepTranscoderMetaData.Bin))
	if err != nil {
		return &EthereumUpkeepTranscoder{}, fmt.Errorf("UpkeepTranscoder instance deployment have failed: %w", err)
	}

	transcoder, err := upkeep_transcoder.NewUpkeepTranscoder(transcoderDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumUpkeepTranscoder{}, fmt.Errorf("failed to instantiate UpkeepTranscoder instance: %w", err)
	}

	return &EthereumUpkeepTranscoder{
		client:     client,
		transcoder: transcoder,
		address:    &transcoderDeploymentData.Address,
	}, nil
}

func LoadUpkeepTranscoder(client *seth.Client, address common.Address) (*EthereumUpkeepTranscoder, error) {
	abi, err := upkeep_transcoder.UpkeepTranscoderMetaData.GetAbi()
	if err != nil {
		return &EthereumUpkeepTranscoder{}, fmt.Errorf("failed to get UpkeepTranscoder ABI: %w", err)
	}

	client.ContractStore.AddABI("UpkeepTranscoder", *abi)
	client.ContractStore.AddBIN("UpkeepTranscoder", common.FromHex(upkeep_transcoder.UpkeepTranscoderMetaData.Bin))

	transcoder, err := upkeep_transcoder.NewUpkeepTranscoder(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumUpkeepTranscoder{}, fmt.Errorf("failed to instantiate UpkeepTranscoder instance: %w", err)
	}

	return &EthereumUpkeepTranscoder{
		client:     client,
		transcoder: transcoder,
		address:    &address,
	}, nil
}

// EthereumKeeperRegistry represents keeper registry contract
type EthereumKeeperRegistry struct {
	client      *seth.Client
	version     ethereum.KeeperRegistryVersion
	registry1_1 *keeper_registry_wrapper1_1.KeeperRegistry
	registry1_2 *keeper_registry_wrapper1_2.KeeperRegistry
	registry1_3 *keeper_registry_wrapper1_3.KeeperRegistry
	registry2_0 *keeper_registry_wrapper2_0.KeeperRegistry
	registry2_1 *i_keeper_registry_master_wrapper_2_1.IKeeperRegistryMaster
	registry2_2 *i_automation_registry_master_wrapper_2_2.IAutomationRegistryMaster
	registry2_3 *i_automation_registry_master_wrapper_2_3.IAutomationRegistryMaster23
	chainModule *i_chain_module.IChainModule
	address     *common.Address
	l           zerolog.Logger
}

func (v *EthereumKeeperRegistry) ReorgProtectionEnabled() bool {
	chainId := v.client.ChainID
	// reorg protection is disabled in polygon zkEVM and Scroll bc currently there is no way to get the block hash onchain
	return v.version < ethereum.RegistryVersion_2_2 || (chainId != 1101 && chainId != 1442 && chainId != 2442 && chainId != 534352 && chainId != 534351)
}

func (v *EthereumKeeperRegistry) ChainModuleAddress() common.Address {
	if v.version >= ethereum.RegistryVersion_2_2 {
		return v.chainModule.Address()
	}
	return common.Address{}
}

func (v *EthereumKeeperRegistry) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperRegistry) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}

func (v *EthereumKeeperRegistry) RegistryOwnerAddress() common.Address {
	callOpts := &bind.CallOpts{
		Pending: false,
	}

	switch v.version {
	case ethereum.RegistryVersion_2_3:
		ownerAddress, _ := v.registry2_3.Owner(callOpts)
		return ownerAddress
	case ethereum.RegistryVersion_2_2:
		ownerAddress, _ := v.registry2_2.Owner(callOpts)
		return ownerAddress
	case ethereum.RegistryVersion_2_1:
		ownerAddress, _ := v.registry2_1.Owner(callOpts)
		return ownerAddress
	case ethereum.RegistryVersion_2_0:
		ownerAddress, _ := v.registry2_0.Owner(callOpts)
		return ownerAddress
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1, ethereum.RegistryVersion_1_2, ethereum.RegistryVersion_1_3:
		return v.client.MustGetRootKeyAddress()
	default:
		return v.client.MustGetRootKeyAddress()
	}
}

func (v *EthereumKeeperRegistry) SetConfigTypeSafe(ocrConfig OCRv2Config) error {
	txOpts := v.client.NewTXOpts()
	var err error
	var decodedTx *seth.DecodedTransaction

	switch v.version {
	case ethereum.RegistryVersion_2_1:
		decodedTx, err = v.client.Decode(v.registry2_1.SetConfigTypeSafe(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.TypedOnchainConfig21,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		))
	case ethereum.RegistryVersion_2_2:
		decodedTx, err = v.client.Decode(v.registry2_2.SetConfigTypeSafe(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.TypedOnchainConfig22,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		))
	case ethereum.RegistryVersion_2_3:
		decodedTx, err = v.client.Decode(v.registry2_3.SetConfigTypeSafe(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.TypedOnchainConfig23,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
			ocrConfig.BillingTokens,
			ocrConfig.BillingConfigs,
		))
	default:
		return fmt.Errorf("SetConfigTypeSafe is not supported in keeper registry version %d", v.version)
	}
	v.l.Debug().Interface("decodedTx", decodedTx).Msg("SetConfigTypeSafe")
	return err
}

func (v *EthereumKeeperRegistry) SetConfig(config KeeperRegistrySettings, ocrConfig OCRv2Config) error {
	txOpts := v.client.NewTXOpts()
	callOpts := bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: nil,
	}

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		_, err := v.client.Decode(v.registry1_1.SetConfig(
			txOpts,
			config.PaymentPremiumPPB,
			config.FlatFeeMicroLINK,
			config.BlockCountPerTurn,
			config.CheckGasLimit,
			config.StalenessSeconds,
			config.GasCeilingMultiplier,
			config.FallbackGasPrice,
			config.FallbackLinkPrice,
		))
		return err
	case ethereum.RegistryVersion_1_2:
		state, err := v.registry1_2.GetState(&callOpts)
		if err != nil {
			return err
		}

		_, err = v.client.Decode(v.registry1_2.SetConfig(txOpts, keeper_registry_wrapper1_2.Config{
			PaymentPremiumPPB:    config.PaymentPremiumPPB,
			FlatFeeMicroLink:     config.FlatFeeMicroLINK,
			BlockCountPerTurn:    config.BlockCountPerTurn,
			CheckGasLimit:        config.CheckGasLimit,
			StalenessSeconds:     config.StalenessSeconds,
			GasCeilingMultiplier: config.GasCeilingMultiplier,
			MinUpkeepSpend:       config.MinUpkeepSpend,
			MaxPerformGas:        config.MaxPerformGas,
			FallbackGasPrice:     config.FallbackGasPrice,
			FallbackLinkPrice:    config.FallbackLinkPrice,
			// Keep the transcoder and registrar same. They have separate setters
			Transcoder: state.Config.Transcoder,
			Registrar:  state.Config.Registrar,
		}))
		return err
	case ethereum.RegistryVersion_1_3:
		state, err := v.registry1_3.GetState(&callOpts)
		if err != nil {
			return err
		}

		_, err = v.client.Decode(v.registry1_3.SetConfig(txOpts, keeper_registry_wrapper1_3.Config{
			PaymentPremiumPPB:    config.PaymentPremiumPPB,
			FlatFeeMicroLink:     config.FlatFeeMicroLINK,
			BlockCountPerTurn:    config.BlockCountPerTurn,
			CheckGasLimit:        config.CheckGasLimit,
			StalenessSeconds:     config.StalenessSeconds,
			GasCeilingMultiplier: config.GasCeilingMultiplier,
			MinUpkeepSpend:       config.MinUpkeepSpend,
			MaxPerformGas:        config.MaxPerformGas,
			FallbackGasPrice:     config.FallbackGasPrice,
			FallbackLinkPrice:    config.FallbackLinkPrice,
			// Keep the transcoder and registrar same. They have separate setters
			Transcoder: state.Config.Transcoder,
			Registrar:  state.Config.Registrar,
		}))
		return err
	case ethereum.RegistryVersion_2_0:
		_, err := v.client.Decode(v.registry2_0.SetConfig(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.OnchainConfig,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		))
		return err
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		return fmt.Errorf("registry version 2.1 2.2 and 2.3 must use setConfigTypeSafe function")
	default:
		return fmt.Errorf("keeper registry version %d is not supported", v.version)
	}
}

func (v *EthereumKeeperRegistry) SetUpkeepOffchainConfig(id *big.Int, offchainConfig []byte) error {
	switch v.version {
	case ethereum.RegistryVersion_2_0:
		_, err := v.client.Decode(v.registry2_0.SetUpkeepOffchainConfig(v.client.NewTXOpts(), id, offchainConfig))
		return err
	case ethereum.RegistryVersion_2_1:
		_, err := v.client.Decode(v.registry2_1.SetUpkeepOffchainConfig(v.client.NewTXOpts(), id, offchainConfig))
		return err
	case ethereum.RegistryVersion_2_2:
		_, err := v.client.Decode(v.registry2_2.SetUpkeepOffchainConfig(v.client.NewTXOpts(), id, offchainConfig))
		return err
	case ethereum.RegistryVersion_2_3:
		_, err := v.client.Decode(v.registry2_3.SetUpkeepOffchainConfig(v.client.NewTXOpts(), id, offchainConfig))
		return err
	default:
		return fmt.Errorf("SetUpkeepOffchainConfig is not supported by keeper registry version %d", v.version)
	}
}

// Pause pauses the registry.
func (v *EthereumKeeperRegistry) Pause() error {
	txOpts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		_, err = v.client.Decode(v.registry1_1.Pause(txOpts))
	case ethereum.RegistryVersion_1_2:
		_, err = v.client.Decode(v.registry1_2.Pause(txOpts))
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.Pause(txOpts))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.Pause(txOpts))
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.Pause(txOpts))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.Pause(txOpts))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.Pause(txOpts))
	default:
		return fmt.Errorf("keeper registry version %d is not supported", v.version)
	}

	return err
}

// Migrate performs a migration of the given upkeep ids to the specific destination passed as parameter.
func (v *EthereumKeeperRegistry) Migrate(upkeepIDs []*big.Int, destinationAddress common.Address) error {
	if v.version != ethereum.RegistryVersion_1_2 {
		return fmt.Errorf("migration of upkeeps is only available for version 1.2 of the registries")
	}

	_, err := v.client.Decode(v.registry1_2.MigrateUpkeeps(v.client.NewTXOpts(), upkeepIDs, destinationAddress))
	return err
}

// SetMigrationPermissions sets the permissions of another registry to allow migrations between the two.
func (v *EthereumKeeperRegistry) SetMigrationPermissions(peerAddress common.Address, permission uint8) error {
	if v.version != ethereum.RegistryVersion_1_2 {
		return fmt.Errorf("migration of upkeeps is only available for version 1.2 of the registries")
	}

	_, err := v.client.Decode(v.registry1_2.SetPeerRegistryMigrationPermission(v.client.NewTXOpts(), peerAddress, permission))
	return err
}

func (v *EthereumKeeperRegistry) SetRegistrar(registrarAddr string) error {
	if v.version == ethereum.RegistryVersion_2_0 {
		// we short circuit and exit, so we don't create a new txs messing up the nonce before exiting
		return fmt.Errorf("please use set config")
	}

	txOpts := v.client.NewTXOpts()
	callOpts := bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: nil,
	}

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		_, err := v.client.Decode(v.registry1_1.SetRegistrar(txOpts, common.HexToAddress(registrarAddr)))
		return err
	case ethereum.RegistryVersion_1_2:
		state, err := v.registry1_2.GetState(&callOpts)
		if err != nil {
			return err
		}
		newConfig := state.Config
		newConfig.Registrar = common.HexToAddress(registrarAddr)
		_, err = v.client.Decode(v.registry1_2.SetConfig(txOpts, newConfig))
		return err
	case ethereum.RegistryVersion_1_3:
		state, err := v.registry1_3.GetState(&callOpts)
		if err != nil {
			return err
		}
		newConfig := state.Config
		newConfig.Registrar = common.HexToAddress(registrarAddr)
		_, err = v.client.Decode(v.registry1_3.SetConfig(txOpts, newConfig))
		return err
	default:
		return fmt.Errorf("keeper registry version %d is not supported", v.version)
	}
}

// AddUpkeepFunds adds link for particular upkeep id
func (v *EthereumKeeperRegistry) AddUpkeepFundsFromKey(id *big.Int, amount *big.Int, keyNum int) error {
	opts := v.client.NewTXKeyOpts(keyNum)
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		_, err = v.client.Decode(v.registry1_1.AddFunds(opts, id, amount))
	case ethereum.RegistryVersion_1_2:
		_, err = v.client.Decode(v.registry1_2.AddFunds(opts, id, amount))
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.AddFunds(opts, id, amount))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.AddFunds(opts, id, amount))
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.AddFunds(opts, id, amount))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.AddFunds(opts, id, amount))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.AddFunds(opts, id, amount))
	}

	return err
}

// AddUpkeepFunds adds link for particular upkeep id
func (v *EthereumKeeperRegistry) AddUpkeepFunds(id *big.Int, amount *big.Int) error {
	return v.AddUpkeepFundsFromKey(id, amount, 0)
}

// GetUpkeepInfo gets upkeep info
func (v *EthereumKeeperRegistry) GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		uk, err := v.registry1_1.GetUpkeep(opts, id)
		if err != nil {
			return nil, err
		}
		return &UpkeepInfo{
			Target:              uk.Target.Hex(),
			ExecuteGas:          uk.ExecuteGas,
			CheckData:           uk.CheckData,
			Balance:             uk.Balance,
			LastKeeper:          uk.LastKeeper.Hex(),
			Admin:               uk.Admin.Hex(),
			MaxValidBlocknumber: uk.MaxValidBlocknumber,
		}, nil
	case ethereum.RegistryVersion_1_2:
		uk, err := v.registry1_2.GetUpkeep(opts, id)
		if err != nil {
			return nil, err
		}
		return &UpkeepInfo{
			Target:              uk.Target.Hex(),
			ExecuteGas:          uk.ExecuteGas,
			CheckData:           uk.CheckData,
			Balance:             uk.Balance,
			LastKeeper:          uk.LastKeeper.Hex(),
			Admin:               uk.Admin.Hex(),
			MaxValidBlocknumber: uk.MaxValidBlocknumber,
		}, nil
	case ethereum.RegistryVersion_1_3:
		uk, err := v.registry1_3.GetUpkeep(opts, id)
		if err != nil {
			return nil, err
		}
		return &UpkeepInfo{
			Target:              uk.Target.Hex(),
			ExecuteGas:          uk.ExecuteGas,
			CheckData:           uk.CheckData,
			Balance:             uk.Balance,
			LastKeeper:          uk.LastKeeper.Hex(),
			Admin:               uk.Admin.Hex(),
			MaxValidBlocknumber: uk.MaxValidBlocknumber,
		}, nil
	case ethereum.RegistryVersion_2_0:
		uk, err := v.registry2_0.GetUpkeep(opts, id)
		if err != nil {
			return nil, err
		}
		return &UpkeepInfo{
			Target:                 uk.Target.Hex(),
			ExecuteGas:             uk.ExecuteGas,
			CheckData:              uk.CheckData,
			Balance:                uk.Balance,
			Admin:                  uk.Admin.Hex(),
			MaxValidBlocknumber:    uk.MaxValidBlocknumber,
			LastPerformBlockNumber: uk.LastPerformBlockNumber,
			AmountSpent:            uk.AmountSpent,
			Paused:                 uk.Paused,
			OffchainConfig:         uk.OffchainConfig,
		}, nil
	case ethereum.RegistryVersion_2_1:
		uk, err := v.registry2_1.GetUpkeep(opts, id)
		if err != nil {
			return nil, err
		}
		return &UpkeepInfo{
			Target:                 uk.Target.Hex(),
			ExecuteGas:             uk.PerformGas,
			CheckData:              uk.CheckData,
			Balance:                uk.Balance,
			Admin:                  uk.Admin.Hex(),
			MaxValidBlocknumber:    uk.MaxValidBlocknumber,
			LastPerformBlockNumber: uk.LastPerformedBlockNumber,
			AmountSpent:            uk.AmountSpent,
			Paused:                 uk.Paused,
			OffchainConfig:         uk.OffchainConfig,
		}, nil
	case ethereum.RegistryVersion_2_2:
		return v.getUpkeepInfo22(opts, id)
	case ethereum.RegistryVersion_2_3:
		return v.getUpkeepInfo23(opts, id)
	}

	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

func (v *EthereumKeeperRegistry) getUpkeepInfo22(opts *bind.CallOpts, id *big.Int) (*UpkeepInfo, error) {
	uk, err := v.registry2_2.GetUpkeep(opts, id)
	if err != nil {
		return nil, err
	}
	return &UpkeepInfo{
		Target:                 uk.Target.Hex(),
		ExecuteGas:             uk.PerformGas,
		CheckData:              uk.CheckData,
		Balance:                uk.Balance,
		Admin:                  uk.Admin.Hex(),
		MaxValidBlocknumber:    uk.MaxValidBlocknumber,
		LastPerformBlockNumber: uk.LastPerformedBlockNumber,
		AmountSpent:            uk.AmountSpent,
		Paused:                 uk.Paused,
		OffchainConfig:         uk.OffchainConfig,
	}, nil
}

func (v *EthereumKeeperRegistry) getUpkeepInfo23(opts *bind.CallOpts, id *big.Int) (*UpkeepInfo, error) {
	uk, err := v.registry2_3.GetUpkeep(opts, id)
	if err != nil {
		return nil, err
	}
	return &UpkeepInfo{
		Target:                 uk.Target.Hex(),
		ExecuteGas:             uk.PerformGas,
		CheckData:              uk.CheckData,
		Balance:                uk.Balance,
		Admin:                  uk.Admin.Hex(),
		MaxValidBlocknumber:    uk.MaxValidBlocknumber,
		LastPerformBlockNumber: uk.LastPerformedBlockNumber,
		AmountSpent:            uk.AmountSpent,
		Paused:                 uk.Paused,
		OffchainConfig:         uk.OffchainConfig,
	}, nil
}

func (v *EthereumKeeperRegistry) GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	var info struct {
		Payee   common.Address
		Active  bool
		Balance *big.Int
	}
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		info, err = v.registry1_1.GetKeeperInfo(opts, common.HexToAddress(keeperAddr))
	case ethereum.RegistryVersion_1_2:
		info, err = v.registry1_2.GetKeeperInfo(opts, common.HexToAddress(keeperAddr))
	case ethereum.RegistryVersion_1_3:
		info, err = v.registry1_3.GetKeeperInfo(opts, common.HexToAddress(keeperAddr))
	case ethereum.RegistryVersion_2_0, ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		// this is not used anywhere
		return nil, fmt.Errorf("not supported")
	}

	if err != nil {
		return nil, err
	}
	return &KeeperInfo{
		Payee:   info.Payee.Hex(),
		Active:  info.Active,
		Balance: info.Balance,
	}, nil
}

func (v *EthereumKeeperRegistry) SetKeepers(keepers []string, payees []string, ocrConfig OCRv2Config) error {
	opts := v.client.NewTXOpts()
	var err error

	keepersAddresses := make([]common.Address, 0)
	for _, k := range keepers {
		keepersAddresses = append(keepersAddresses, common.HexToAddress(k))
	}
	payeesAddresses := make([]common.Address, 0)
	for _, p := range payees {
		payeesAddresses = append(payeesAddresses, common.HexToAddress(p))
	}

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		_, err = v.client.Decode(v.registry1_1.SetKeepers(opts, keepersAddresses, payeesAddresses))
	case ethereum.RegistryVersion_1_2:
		_, err = v.client.Decode(v.registry1_2.SetKeepers(opts, keepersAddresses, payeesAddresses))
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.SetKeepers(opts, keepersAddresses, payeesAddresses))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.SetConfig(opts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.OnchainConfig,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		))
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		return fmt.Errorf("not supported")
	}

	return err
}

// RegisterUpkeep registers contract to perform upkeep
func (v *EthereumKeeperRegistry) RegisterUpkeep(target string, gasLimit uint32, admin string, checkData []byte) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		_, err = v.client.Decode(v.registry1_1.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
		))
	case ethereum.RegistryVersion_1_2:
		_, err = v.client.Decode(v.registry1_2.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
		))
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
		))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
			nil, //offchain config
		))
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		return fmt.Errorf("not supported")
	}

	return err
}

// CancelUpkeep cancels the given upkeep ID
func (v *EthereumKeeperRegistry) CancelUpkeep(id *big.Int) error {
	opts := v.client.NewTXOpts()
	var err error
	var tx *seth.DecodedTransaction

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err = v.client.Decode(v.registry1_1.CancelUpkeep(opts, id))
	case ethereum.RegistryVersion_1_2:
		tx, err = v.client.Decode(v.registry1_2.CancelUpkeep(opts, id))
	case ethereum.RegistryVersion_1_3:
		tx, err = v.client.Decode(v.registry1_3.CancelUpkeep(opts, id))
	case ethereum.RegistryVersion_2_0:
		tx, err = v.client.Decode(v.registry2_0.CancelUpkeep(opts, id))
	case ethereum.RegistryVersion_2_1:
		tx, err = v.client.Decode(v.registry2_1.CancelUpkeep(opts, id))
	case ethereum.RegistryVersion_2_2:
		tx, err = v.client.Decode(v.registry2_2.CancelUpkeep(opts, id))
	case ethereum.RegistryVersion_2_3:
		tx, err = v.client.Decode(v.registry2_3.CancelUpkeep(opts, id))
	}

	txHash := "none"
	if err == nil && tx != nil {
		txHash = tx.Hash
	}

	v.l.Info().
		Str("Upkeep ID", strconv.FormatInt(id.Int64(), 10)).
		Str("From", v.client.MustGetRootKeyAddress().Hex()).
		Str("TX Hash", txHash).
		Msg("Cancel Upkeep tx")

	return err
}

// SetUpkeepGasLimit sets the perform gas limit for a given upkeep ID
func (v *EthereumKeeperRegistry) SetUpkeepGasLimit(id *big.Int, gas uint32) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_2:
		_, err = v.client.Decode(v.registry1_2.SetUpkeepGasLimit(opts, id, gas))
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.SetUpkeepGasLimit(opts, id, gas))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.SetUpkeepGasLimit(opts, id, gas))
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.SetUpkeepGasLimit(opts, id, gas))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.SetUpkeepGasLimit(opts, id, gas))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.SetUpkeepGasLimit(opts, id, gas))
	default:
		return fmt.Errorf("keeper registry version %d is not supported for SetUpkeepGasLimit", v.version)
	}

	return err
}

// GetKeeperList get list of all registered keeper addresses
func (v *EthereumKeeperRegistry) GetKeeperList(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	var list []common.Address
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		list, err = v.registry1_1.GetKeeperList(opts)
	case ethereum.RegistryVersion_1_2:
		state, err := v.registry1_2.GetState(opts)
		if err != nil {
			return []string{}, err
		}
		list = state.Keepers
	case ethereum.RegistryVersion_1_3:
		state, err := v.registry1_3.GetState(opts)
		if err != nil {
			return []string{}, err
		}
		list = state.Keepers
	case ethereum.RegistryVersion_2_0:
		state, err := v.registry2_0.GetState(opts)
		if err != nil {
			return []string{}, err
		}
		list = state.Transmitters
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		return nil, fmt.Errorf("not supported")
	}

	if err != nil {
		return []string{}, err
	}
	addrs := make([]string, 0)
	for _, ca := range list {
		addrs = append(addrs, ca.Hex())
	}
	return addrs, nil
}

// UpdateCheckData updates the check data of an upkeep
func (v *EthereumKeeperRegistry) UpdateCheckData(id *big.Int, newCheckData []byte) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.UpdateCheckData(opts, id, newCheckData))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.UpdateCheckData(opts, id, newCheckData))
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.SetUpkeepCheckData(opts, id, newCheckData))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.SetUpkeepCheckData(opts, id, newCheckData))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.SetUpkeepCheckData(opts, id, newCheckData))
	default:
		return fmt.Errorf("UpdateCheckData is not supported by keeper registry version %d", v.version)
	}

	return err
}

// SetUpkeepTriggerConfig updates the trigger config of an upkeep (only for version 2.1)
func (v *EthereumKeeperRegistry) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.SetUpkeepTriggerConfig(opts, id, triggerConfig))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.SetUpkeepTriggerConfig(opts, id, triggerConfig))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.SetUpkeepTriggerConfig(opts, id, triggerConfig))
	default:
		return fmt.Errorf("SetUpkeepTriggerConfig is not supported by keeper registry version %d", v.version)
	}

	return err
}

// SetUpkeepPrivilegeConfig sets the privilege config of an upkeep (only for version 2.1)
func (v *EthereumKeeperRegistry) SetUpkeepPrivilegeConfig(id *big.Int, privilegeConfig []byte) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.SetUpkeepPrivilegeConfig(opts, id, privilegeConfig))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.SetUpkeepPrivilegeConfig(opts, id, privilegeConfig))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.SetUpkeepPrivilegeConfig(opts, id, privilegeConfig))
	default:
		return fmt.Errorf("SetUpkeepPrivilegeConfig is not supported by keeper registry version %d", v.version)
	}

	return err
}

// PauseUpkeep stops an upkeep from an upkeep
func (v *EthereumKeeperRegistry) PauseUpkeep(id *big.Int) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.PauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.PauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.PauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.PauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.PauseUpkeep(opts, id))
	default:
		return fmt.Errorf("PauseUpkeep is not supported by keeper registry version %d", v.version)
	}

	return err
}

// UnpauseUpkeep get list of all registered keeper addresses
func (v *EthereumKeeperRegistry) UnpauseUpkeep(id *big.Int) error {
	opts := v.client.NewTXOpts()
	var err error

	switch v.version {
	case ethereum.RegistryVersion_1_3:
		_, err = v.client.Decode(v.registry1_3.UnpauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_0:
		_, err = v.client.Decode(v.registry2_0.UnpauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_1:
		_, err = v.client.Decode(v.registry2_1.UnpauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_2:
		_, err = v.client.Decode(v.registry2_2.UnpauseUpkeep(opts, id))
	case ethereum.RegistryVersion_2_3:
		_, err = v.client.Decode(v.registry2_3.UnpauseUpkeep(opts, id))
	default:
		return fmt.Errorf("UnpauseUpkeep is not supported by keeper registry version %d", v.version)
	}

	return err
}

// Parses upkeep performed log
func (v *EthereumKeeperRegistry) ParseUpkeepPerformedLog(log *types.Log) (*UpkeepPerformedLog, error) {
	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		parsedLog, err := v.registry1_1.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    parsedLog.From,
		}, nil
	case ethereum.RegistryVersion_1_2:
		parsedLog, err := v.registry1_2.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    parsedLog.From,
		}, nil
	case ethereum.RegistryVersion_1_3:
		parsedLog, err := v.registry1_3.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    parsedLog.From,
		}, nil
	case ethereum.RegistryVersion_2_0:
		parsedLog, err := v.registry2_0.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    utils.ZeroAddress,
		}, nil
	case ethereum.RegistryVersion_2_1:
		parsedLog, err := v.registry2_1.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    utils.ZeroAddress,
		}, nil
	case ethereum.RegistryVersion_2_2:
		parsedLog, err := v.registry2_2.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    utils.ZeroAddress,
		}, nil
	case ethereum.RegistryVersion_2_3:
		parsedLog, err := v.registry2_3.ParseUpkeepPerformed(*log)
		if err != nil {
			return nil, err
		}
		return &UpkeepPerformedLog{
			Id:      parsedLog.Id,
			Success: parsedLog.Success,
			From:    utils.ZeroAddress,
		}, nil
	}
	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

// ParseStaleUpkeepReportLog Parses Stale upkeep report log
func (v *EthereumKeeperRegistry) ParseStaleUpkeepReportLog(log *types.Log) (*StaleUpkeepReportLog, error) {
	//nolint:exhaustive
	switch v.version {
	case ethereum.RegistryVersion_2_0:
		parsedLog, err := v.registry2_0.ParseStaleUpkeepReport(*log)
		if err != nil {
			return nil, err
		}
		return &StaleUpkeepReportLog{
			Id: parsedLog.Id,
		}, nil
	case ethereum.RegistryVersion_2_1:
		parsedLog, err := v.registry2_1.ParseStaleUpkeepReport(*log)
		if err != nil {
			return nil, err
		}
		return &StaleUpkeepReportLog{
			Id: parsedLog.Id,
		}, nil
	case ethereum.RegistryVersion_2_2:
		parsedLog, err := v.registry2_2.ParseStaleUpkeepReport(*log)
		if err != nil {
			return nil, err
		}
		return &StaleUpkeepReportLog{
			Id: parsedLog.Id,
		}, nil
	case ethereum.RegistryVersion_2_3:
		parsedLog, err := v.registry2_3.ParseStaleUpkeepReport(*log)
		if err != nil {
			return nil, err
		}
		return &StaleUpkeepReportLog{
			Id: parsedLog.Id,
		}, nil
	}
	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

// Parses the upkeep ID from an 'UpkeepRegistered' log, returns error on any other log
func (v *EthereumKeeperRegistry) ParseUpkeepIdFromRegisteredLog(log *types.Log) (*big.Int, error) {
	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		parsedLog, err := v.registry1_1.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	case ethereum.RegistryVersion_1_2:
		parsedLog, err := v.registry1_2.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	case ethereum.RegistryVersion_1_3:
		parsedLog, err := v.registry1_3.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	case ethereum.RegistryVersion_2_0:
		parsedLog, err := v.registry2_0.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	case ethereum.RegistryVersion_2_1:
		parsedLog, err := v.registry2_1.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	case ethereum.RegistryVersion_2_2:
		parsedLog, err := v.registry2_2.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	case ethereum.RegistryVersion_2_3:
		parsedLog, err := v.registry2_3.ParseUpkeepRegistered(*log)
		if err != nil {
			return nil, err
		}
		return parsedLog.Id, nil
	}

	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

func DeployKeeperRegistry(
	client *seth.Client,
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	var mode uint8
	switch client.ChainID {
	//Arbitrum payment model
	case networks.ArbitrumMainnet.ChainID, networks.ArbitrumSepolia.ChainID:
		mode = uint8(1)
	//Optimism payment model
	case networks.OptimismMainnet.ChainID, networks.OptimismSepolia.ChainID:
		mode = uint8(2)
	//Base
	case networks.BaseMainnet.ChainID, networks.BaseSepolia.ChainID:
		mode = uint8(2)
	default:
		mode = uint8(0)
	}
	registryGasOverhead := big.NewInt(80000)
	switch opts.RegistryVersion {
	case eth_contracts.RegistryVersion_1_0, eth_contracts.RegistryVersion_1_1:
		return deployRegistry10_11(client, opts)
	case eth_contracts.RegistryVersion_1_2:
		return deployRegistry12(client, opts)
	case eth_contracts.RegistryVersion_1_3:
		return deployRegistry13(client, opts, mode, registryGasOverhead)
	case eth_contracts.RegistryVersion_2_0:
		return deployRegistry20(client, opts, mode)
	case eth_contracts.RegistryVersion_2_1:
		return deployRegistry21(client, opts, mode)
	case eth_contracts.RegistryVersion_2_2:
		return deployRegistry22(client, opts)
	case eth_contracts.RegistryVersion_2_3:
		return deployRegistry23(client, opts)
	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", opts.RegistryVersion)
	}
}

func deployRegistry10_11(client *seth.Client, opts *KeeperRegistryOpts) (KeeperRegistry, error) {
	abi, err := keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_1 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistry1_1", *abi, common.FromHex(keeper_registry_wrapper1_1.KeeperRegistryMetaData.Bin),
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.ETHFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
		opts.Settings.PaymentPremiumPPB,
		opts.Settings.FlatFeeMicroLINK,
		opts.Settings.BlockCountPerTurn,
		opts.Settings.CheckGasLimit,
		opts.Settings.StalenessSeconds,
		opts.Settings.GasCeilingMultiplier,
		opts.Settings.FallbackGasPrice,
		opts.Settings.FallbackLinkPrice,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistry1_1 instance deployment have failed: %w", err)
	}

	instance, err := keeper_registry_wrapper1_1.NewKeeperRegistry(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry1_1 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_1_1,
		registry1_1: instance,
		registry1_2: nil,
		registry1_3: nil,
		address:     &data.Address,
	}, err
}

func deployRegistry12(client *seth.Client, opts *KeeperRegistryOpts) (KeeperRegistry, error) {
	abi, err := keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_2 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistry1_2", *abi, common.FromHex(keeper_registry_wrapper1_2.KeeperRegistryMetaData.Bin),
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.ETHFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
		keeper_registry_wrapper1_2.Config{
			PaymentPremiumPPB:    opts.Settings.PaymentPremiumPPB,
			FlatFeeMicroLink:     opts.Settings.FlatFeeMicroLINK,
			BlockCountPerTurn:    opts.Settings.BlockCountPerTurn,
			CheckGasLimit:        opts.Settings.CheckGasLimit,
			StalenessSeconds:     opts.Settings.StalenessSeconds,
			GasCeilingMultiplier: opts.Settings.GasCeilingMultiplier,
			MinUpkeepSpend:       opts.Settings.MinUpkeepSpend,
			MaxPerformGas:        opts.Settings.MaxPerformGas,
			FallbackGasPrice:     opts.Settings.FallbackGasPrice,
			FallbackLinkPrice:    opts.Settings.FallbackLinkPrice,
			Transcoder:           common.HexToAddress(opts.TranscoderAddr),
			Registrar:            common.HexToAddress(opts.RegistrarAddr),
		},
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistry1_2 instance deployment have failed: %w", err)
	}

	instance, err := keeper_registry_wrapper1_2.NewKeeperRegistry(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry1_2 instance: %w", err)
	}
	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_1_2,
		registry1_1: nil,
		registry1_2: instance,
		registry1_3: nil,
		address:     &data.Address,
	}, err
}

func deployRegistry13(client *seth.Client, opts *KeeperRegistryOpts, mode uint8, registryGasOverhead *big.Int) (KeeperRegistry, error) {
	logicAbi, err := keeper_registry_logic1_3.KeeperRegistryLogicMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistryLogic1_3 ABI: %w", err)
	}
	logicData, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistryLogic1_3", *logicAbi, common.FromHex(keeper_registry_logic1_3.KeeperRegistryLogicMetaData.Bin),
		mode,                // Default payment model
		registryGasOverhead, // Registry gas overhead
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.ETHFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistryLogic1_3 instance deployment have failed: %w", err)
	}

	abi, err := keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_3 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistry1_3", *abi, common.FromHex(keeper_registry_wrapper1_3.KeeperRegistryMetaData.Bin),
		logicData.Address,
		keeper_registry_wrapper1_3.Config{
			PaymentPremiumPPB:    opts.Settings.PaymentPremiumPPB,
			FlatFeeMicroLink:     opts.Settings.FlatFeeMicroLINK,
			BlockCountPerTurn:    opts.Settings.BlockCountPerTurn,
			CheckGasLimit:        opts.Settings.CheckGasLimit,
			StalenessSeconds:     opts.Settings.StalenessSeconds,
			GasCeilingMultiplier: opts.Settings.GasCeilingMultiplier,
			MinUpkeepSpend:       opts.Settings.MinUpkeepSpend,
			MaxPerformGas:        opts.Settings.MaxPerformGas,
			FallbackGasPrice:     opts.Settings.FallbackGasPrice,
			FallbackLinkPrice:    opts.Settings.FallbackLinkPrice,
			Transcoder:           common.HexToAddress(opts.TranscoderAddr),
			Registrar:            common.HexToAddress(opts.RegistrarAddr),
		},
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistry1_3 instance deployment have failed: %w", err)
	}

	instance, err := keeper_registry_wrapper1_3.NewKeeperRegistry(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry1_3 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_1_3,
		registry1_1: nil,
		registry1_2: nil,
		registry1_3: instance,
		address:     &data.Address,
	}, err
}

func deployRegistry20(client *seth.Client, opts *KeeperRegistryOpts, mode uint8) (KeeperRegistry, error) {
	logicAbi, err := keeper_registry_logic2_0.KeeperRegistryLogicMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistryLogic2_0 ABI: %w", err)
	}
	logicData, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistryLogic2_0", *logicAbi, common.FromHex(keeper_registry_logic2_0.KeeperRegistryLogicMetaData.Bin),
		mode, // Default payment model
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.ETHFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistryLogic2_0 instance deployment have failed: %w", err)
	}

	abi, err := keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_3 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistry2_0", *abi, common.FromHex(keeper_registry_wrapper2_0.KeeperRegistryMetaData.Bin),
		logicData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistry2_0 instance deployment have failed: %w", err)
	}

	instance, err := keeper_registry_wrapper2_0.NewKeeperRegistry(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry2_0 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_2_0,
		registry2_0: instance,
		address:     &data.Address,
	}, err
}

func deployRegistry21(client *seth.Client, opts *KeeperRegistryOpts, mode uint8) (KeeperRegistry, error) {
	automationForwarderLogicAddr, err := deployAutomationForwarderLogicSeth(client)
	if err != nil {
		return nil, err
	}

	logicBAbi, err := registrylogicb21.KeeperRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistryLogicB2_1 ABI: %w", err)
	}
	logicBData, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistryLogicB2_1", *logicBAbi, common.FromHex(registrylogicb21.KeeperRegistryLogicBMetaData.Bin),
		mode,
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.ETHFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
		automationForwarderLogicAddr,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistryLogicB2_1 instance deployment have failed: %w", err)
	}

	logicAAbi, err := registrylogica21.KeeperRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistryLogicA2_1 ABI: %w", err)
	}
	logicAData, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistryLogicA2_1", *logicAAbi, common.FromHex(registrylogica21.KeeperRegistryLogicAMetaData.Bin),
		logicBData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistryLogicA2_1 instance deployment have failed: %w", err)
	}

	abi, err := registry21.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry2_1 ABI: %w", err)
	}

	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistry2_1", *abi, common.FromHex(registry21.KeeperRegistryMetaData.Bin),
		logicAData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("KeeperRegistry2_1 instance deployment have failed: %w", err)
	}

	instance, err := iregistry21.NewIKeeperRegistryMaster(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry2_1 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_2_1,
		registry2_1: instance,
		address:     &data.Address,
	}, err
}

func deployRegistry22(client *seth.Client, opts *KeeperRegistryOpts) (KeeperRegistry, error) {
	var chainModuleAddr common.Address
	var err error
	chainId := client.ChainID

	if chainId == networks.ScrollMainnet.ChainID || chainId == networks.ScrollSepolia.ChainID {
		chainModuleAddr, err = deployScrollModule(client)
	} else if chainId == networks.ArbitrumMainnet.ChainID || chainId == networks.ArbitrumSepolia.ChainID {
		chainModuleAddr, err = deployArbitrumModule(client)
	} else if chainId == networks.OptimismMainnet.ChainID || chainId == networks.OptimismSepolia.ChainID {
		chainModuleAddr, err = deployOptimismModule(client)
	} else {
		chainModuleAddr, err = deployBaseModule(client)
	}
	if err != nil {
		return nil, err
	}

	automationForwarderLogicAddr, err := deployAutomationForwarderLogicSeth(client)
	if err != nil {
		return nil, err
	}

	allowedReadOnlyAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	logicBAbi, err := registrylogicb22.AutomationRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistryLogicB2_2 ABI: %w", err)
	}

	logicBData, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistryLogicB2_2", *logicBAbi, common.FromHex(registrylogicb22.AutomationRegistryLogicBMetaData.Bin),
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.ETHFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
		automationForwarderLogicAddr,
		allowedReadOnlyAddress,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistryLogicB2_2 instance deployment have failed: %w", err)
	}

	logicAAbi, err := registrylogica22.AutomationRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistryLogicA2_2 ABI: %w", err)
	}
	logicAData, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistryLogicA2_2", *logicAAbi, common.FromHex(registrylogica22.AutomationRegistryLogicAMetaData.Bin),
		logicBData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistryLogicA2_2 instance deployment have failed: %w", err)
	}

	abi, err := registry22.AutomationRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistry2_2 ABI: %w", err)
	}

	data, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistry2_2", *abi, common.FromHex(registry22.AutomationRegistryMetaData.Bin),
		logicAData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistry2_2 instance deployment have failed: %w", err)
	}

	instance, err := iregistry22.NewIAutomationRegistryMaster(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate AutomationRegistry2_2 instance: %w", err)
	}

	chainModule, err := i_chain_module.NewIChainModule(
		chainModuleAddr,
		wrappers.MustNewWrappedContractBackend(nil, client),
	)

	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_2_2,
		registry2_2: instance,
		chainModule: chainModule,
		address:     &data.Address,
	}, err
}

func deployRegistry23(client *seth.Client, opts *KeeperRegistryOpts) (KeeperRegistry, error) {
	var chainModuleAddr common.Address
	var err error
	chainId := client.ChainID

	if chainId == networks.ScrollMainnet.ChainID || chainId == networks.ScrollSepolia.ChainID {
		chainModuleAddr, err = deployScrollModule(client)
	} else if chainId == networks.ArbitrumMainnet.ChainID || chainId == networks.ArbitrumSepolia.ChainID {
		chainModuleAddr, err = deployArbitrumModule(client)
	} else if chainId == networks.OptimismMainnet.ChainID || chainId == networks.OptimismSepolia.ChainID {
		chainModuleAddr, err = deployOptimismModule(client)
	} else {
		chainModuleAddr, err = deployBaseModule(client)
	}
	if err != nil {
		return nil, err
	}

	automationForwarderLogicAddr, err := deployAutomationForwarderLogicSeth(client)
	if err != nil {
		return nil, err
	}

	allowedReadOnlyAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	logicCAbi, err := registrylogicc23.AutomationRegistryLogicCMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistryLogicC2_3 ABI: %w", err)
	}

	logicCData, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistryLogicC2_3", *logicCAbi, common.FromHex(registrylogicc23.AutomationRegistryLogicCMetaData.Bin),
		common.HexToAddress(opts.LinkAddr),
		common.HexToAddress(opts.LinkUSDFeedAddr),
		common.HexToAddress(opts.NativeUSDFeedAddr),
		common.HexToAddress(opts.GasFeedAddr),
		automationForwarderLogicAddr,
		allowedReadOnlyAddress,
		uint8(0), // onchain payout mode
		common.HexToAddress(opts.WrappedNativeAddr),
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistryLogicC2_3 instance deployment have failed: %w", err)
	}

	logicBAbi, err := registrylogicb23.AutomationRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistryLogicB2_3 ABI: %w", err)
	}

	logicBData, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistryLogicB2_3", *logicBAbi, common.FromHex(registrylogicb23.AutomationRegistryLogicBMetaData.Bin),
		logicCData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistryLogicB2_3 instance deployment have failed: %w", err)
	}

	logicAAbi, err := registrylogica23.AutomationRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistryLogicA2_3 ABI: %w", err)
	}
	logicAData, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistryLogicA2_3", *logicAAbi, common.FromHex(registrylogica23.AutomationRegistryLogicAMetaData.Bin),
		logicBData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistryLogicA2_3 instance deployment have failed: %w", err)
	}

	abi, err := registry23.AutomationRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistry2_3 ABI: %w", err)
	}

	data, err := client.DeployContract(client.NewTXOpts(), "AutomationRegistry2_3", *abi, common.FromHex(registry23.AutomationRegistryMetaData.Bin),
		logicAData.Address,
	)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("AutomationRegistry2_3 instance deployment have failed: %w", err)
	}

	instance, err := iregistry23.NewIAutomationRegistryMaster23(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate AutomationRegistry2_3 instance: %w", err)
	}

	chainModule, err := i_chain_module.NewIChainModule(
		chainModuleAddr,
		wrappers.MustNewWrappedContractBackend(nil, client),
	)

	return &EthereumKeeperRegistry{
		client:      client,
		version:     eth_contracts.RegistryVersion_2_3,
		registry2_3: instance,
		chainModule: chainModule,
		address:     &data.Address,
	}, err
}

// LoadKeeperRegistry returns deployed on given address EthereumKeeperRegistry
func LoadKeeperRegistry(l zerolog.Logger, client *seth.Client, address common.Address, registryVersion eth_contracts.KeeperRegistryVersion, chainModuleAddress common.Address) (KeeperRegistry, error) {
	var keeper *EthereumKeeperRegistry
	var err error
	switch registryVersion {
	case eth_contracts.RegistryVersion_1_1:
		keeper, err = loadRegistry1_1(client, address)
	case eth_contracts.RegistryVersion_1_2:
		keeper, err = loadRegistry1_2(client, address)
	case eth_contracts.RegistryVersion_1_3:
		keeper, err = loadRegistry1_3(client, address)
	case eth_contracts.RegistryVersion_2_0:
		keeper, err = loadRegistry2_0(client, address)
	case eth_contracts.RegistryVersion_2_1:
		keeper, err = loadRegistry2_1(client, address)
	case eth_contracts.RegistryVersion_2_2: // why the contract name is not the same as the actual contract name?
		keeper, err = loadRegistry2_2(client, address)
	case eth_contracts.RegistryVersion_2_3:
		keeper, err = loadRegistry2_3(client, address, chainModuleAddress)
	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", registryVersion)
	}

	if keeper != nil {
		keeper.version = registryVersion
		keeper.l = l
	}
	return keeper, err
}

func loadRegistry1_1(client *seth.Client, address common.Address) (*EthereumKeeperRegistry, error) {
	abi, err := keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_1 ABI: %w", err)
	}

	client.ContractStore.AddABI("KeeperRegistry1_1", *abi)
	client.ContractStore.AddBIN("KeeperRegistry1_1", common.FromHex(keeper_registry_wrapper1_1.KeeperRegistryMetaData.Bin))

	instance, err := keeper_registry_wrapper1_1.NewKeeperRegistry(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry1_1 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry1_1: instance,
	}, nil
}

func loadRegistry1_2(client *seth.Client, address common.Address) (*EthereumKeeperRegistry, error) {
	abi, err := keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_2 ABI: %w", err)
	}

	client.ContractStore.AddABI("KeeperRegistry1_2", *abi)
	client.ContractStore.AddBIN("KeeperRegistry1_2", common.FromHex(keeper_registry_wrapper1_2.KeeperRegistryMetaData.Bin))

	instance, err := keeper_registry_wrapper1_2.NewKeeperRegistry(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry1_2 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry1_2: instance,
	}, nil
}

func loadRegistry1_3(client *seth.Client, address common.Address) (*EthereumKeeperRegistry, error) {
	abi, err := keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry1_3 ABI: %w", err)
	}

	client.ContractStore.AddABI("KeeperRegistry1_3", *abi)
	client.ContractStore.AddBIN("KeeperRegistry1_3", common.FromHex(keeper_registry_wrapper1_3.KeeperRegistryMetaData.Bin))

	instance, err := keeper_registry_wrapper1_3.NewKeeperRegistry(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry1_3 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry1_3: instance,
	}, nil
}

func loadRegistry2_0(client *seth.Client, address common.Address) (*EthereumKeeperRegistry, error) {
	abi, err := keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry2_0 ABI: %w", err)
	}

	client.ContractStore.AddABI("KeeperRegistry2_0", *abi)
	client.ContractStore.AddBIN("KeeperRegistry2_0", common.FromHex(keeper_registry_wrapper2_0.KeeperRegistryMetaData.Bin))

	instance, err := keeper_registry_wrapper2_0.NewKeeperRegistry(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry2_0 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry2_0: instance,
	}, nil
}

func loadRegistry2_1(client *seth.Client, address common.Address) (*EthereumKeeperRegistry, error) {
	abi, err := ac.IAutomationV21PlusCommonMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get KeeperRegistry2_1 ABI: %w", err)
	}

	client.ContractStore.AddABI("KeeperRegistry2_1", *abi)
	client.ContractStore.AddBIN("KeeperRegistry2_1", common.FromHex(ac.IAutomationV21PlusCommonMetaData.Bin))

	var instance interface{}

	instance, err = ac.NewIAutomationV21PlusCommon(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate KeeperRegistry2_1 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry2_1: instance.(*iregistry21.IKeeperRegistryMaster),
	}, nil
}

func loadRegistry2_2(client *seth.Client, address common.Address) (*EthereumKeeperRegistry, error) {
	abi, err := iregistry22.IAutomationRegistryMasterMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to get AutomationRegistry2_2 ABI: %w", err)
	}

	client.ContractStore.AddABI("AutomationRegistry2_2", *abi)
	client.ContractStore.AddBIN("AutomationRegistry2_2", common.FromHex(iregistry22.IAutomationRegistryMasterMetaData.Bin))

	instance, err := iregistry22.NewIAutomationRegistryMaster(address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to instantiate AutomationRegistry2_2 instance: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry2_2: instance,
	}, nil
}

func loadRegistry2_3(client *seth.Client, address, chainModuleAddress common.Address) (*EthereumKeeperRegistry, error) {

	loader := seth.NewContractLoader[iregistry23.IAutomationRegistryMaster23](client)
	instance, err := loader.LoadContract("AutomationRegistry2_3", address, iregistry23.IAutomationRegistryMaster23MetaData.GetAbi, iregistry23.NewIAutomationRegistryMaster23)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to load AutomationRegistry2_3 instance: %w", err)
	}

	chainModule, err := loadChainModule(client, chainModuleAddress)
	if err != nil {
		return &EthereumKeeperRegistry{}, fmt.Errorf("failed to load chain module: %w", err)
	}

	return &EthereumKeeperRegistry{
		address:     &address,
		client:      client,
		registry2_3: instance,
		chainModule: chainModule,
	}, nil
}

func deployAutomationForwarderLogicSeth(client *seth.Client) (common.Address, error) {
	abi, err := automationForwarderLogic.AutomationForwarderLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get AutomationForwarderLogic ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "AutomationForwarderLogic", *abi, common.FromHex(automationForwarderLogic.AutomationForwarderLogicMetaData.Bin))
	if err != nil {
		return common.Address{}, fmt.Errorf("AutomationForwarderLogic instance deployment have failed: %w", err)
	}

	return data.Address, nil
}

func deployScrollModule(client *seth.Client) (common.Address, error) {
	abi, err := scroll_module.ScrollModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get ScrollModule ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "ScrollModule", *abi, common.FromHex(scroll_module.ScrollModuleMetaData.Bin))
	if err != nil {
		return common.Address{}, fmt.Errorf("ScrollModule instance deployment have failed: %w", err)
	}

	return data.Address, nil
}

func deployArbitrumModule(client *seth.Client) (common.Address, error) {
	abi, err := arbitrum_module.ArbitrumModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get ArbitrumModule ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "ArbitrumModule", *abi, common.FromHex(arbitrum_module.ArbitrumModuleMetaData.Bin))
	if err != nil {
		return common.Address{}, fmt.Errorf("ArbitrumModule instance deployment have failed: %w", err)
	}

	return data.Address, nil
}

func deployOptimismModule(client *seth.Client) (common.Address, error) {
	abi, err := optimism_module.OptimismModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get OptimismModule ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "OptimismModule", *abi, common.FromHex(optimism_module.OptimismModuleMetaData.Bin))
	if err != nil {
		return common.Address{}, fmt.Errorf("OptimismModule instance deployment have failed: %w", err)
	}

	return data.Address, nil
}

func loadChainModule(client *seth.Client, address common.Address) (*i_chain_module.IChainModule, error) {
	abi, err := i_chain_module.IChainModuleMetaData.GetAbi()
	if err != nil {
		return &i_chain_module.IChainModule{}, fmt.Errorf("failed to get IChainModule ABI: %w", err)
	}

	client.ContractStore.AddABI("IChainModule", *abi)
	client.ContractStore.AddBIN("IChainModule", common.FromHex(i_chain_module.IChainModuleMetaData.Bin))

	chainModule, err := i_chain_module.NewIChainModule(
		address,
		wrappers.MustNewWrappedContractBackend(nil, client),
	)
	if err != nil {
		return &i_chain_module.IChainModule{}, fmt.Errorf("failed to instantiate IChainModule instance: %w", err)
	}

	return chainModule, nil
}

func deployBaseModule(client *seth.Client) (common.Address, error) {
	abi, err := chain_module_base.ChainModuleBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get BaseModule ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "BaseModule", *abi, common.FromHex(chain_module_base.ChainModuleBaseMetaData.Bin))
	if err != nil {
		return common.Address{}, fmt.Errorf("BaseModule instance deployment have failed: %w", err)
	}

	return data.Address, nil
}

// EthereumKeeperRegistrar corresponds to the registrar which is used to send requests to the registry when
// registering new upkeeps.
type EthereumKeeperRegistrar struct {
	client      *seth.Client
	registrar   *keeper_registrar_wrapper1_2.KeeperRegistrar
	registrar20 *keeper_registrar_wrapper2_0.KeeperRegistrar
	registrar21 *registrar21.AutomationRegistrar
	registrar23 *registrar23.AutomationRegistrar
	address     *common.Address
}

func (v *EthereumKeeperRegistrar) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperRegistrar) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}

// register Upkeep with native token, only available from v2.3
func (v *EthereumKeeperRegistrar) RegisterUpkeepFromKey(keyNum int, name string, email []byte, upkeepAddr string, gasLimit uint32, adminAddr string, checkData []byte, amount *big.Int, wethTokenAddr string, isLogTrigger bool, isMercury bool) (*types.Transaction, error) {
	if v.registrar23 == nil {
		return nil, fmt.Errorf("RegisterUpkeepFromKey with native token is only supported in registrar version v2.3")
	}

	registrarABI = cltypes.MustGetABI(registrar23.AutomationRegistrarABI)
	txOpts := v.client.NewTXKeyOpts(keyNum, seth.WithValue(amount))

	if isLogTrigger {
		var topic0InBytes [32]byte
		// bytes representation of 0x0000000000000000000000000000000000000000000000000000000000000000
		bytes0 := [32]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
		if isMercury {
			// bytes representation of 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd
			topic0InBytes = [32]byte{209, 255, 233, 228, 85, 129, 193, 29, 125, 159, 46, 213, 247, 82, 23, 205, 75, 233, 248, 183, 238, 230, 175, 15, 109, 3, 244, 109, 229, 57, 86, 205}
		} else {
			// bytes representation of 0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d
			topic0InBytes = [32]byte{
				61, 83, 163, 149, 80, 224, 70, 136,
				6, 88, 39, 243, 187, 134, 88, 76,
				176, 7, 171, 158, 188, 167, 235,
				213, 40, 231, 48, 28, 156, 49, 235, 93,
			}
		}

		logTriggerConfigStruct := acutils.IAutomationV21PlusCommonLogTriggerConfig{
			ContractAddress: common.HexToAddress(upkeepAddr),
			FilterSelector:  0,
			Topic0:          topic0InBytes,
			Topic1:          bytes0,
			Topic2:          bytes0,
			Topic3:          bytes0,
		}
		encodedLogTriggerConfig, err := compatibleUtils.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
		if err != nil {
			return nil, err
		}

		params := registrar23.AutomationRegistrar23RegistrationParams{
			UpkeepContract: common.HexToAddress(upkeepAddr),
			Amount:         amount,
			AdminAddress:   common.HexToAddress(adminAddr),
			GasLimit:       gasLimit,
			TriggerType:    uint8(1),                           // trigger type
			BillingToken:   common.HexToAddress(wethTokenAddr), // native
			Name:           name,
			EncryptedEmail: email,
			CheckData:      checkData,
			TriggerConfig:  encodedLogTriggerConfig, // log trigger upkeep
			OffchainConfig: []byte{},
		}

		decodedTx, err := v.client.Decode(v.registrar23.RegisterUpkeep(txOpts,
			params,
		))
		return decodedTx.Transaction, err
	}

	params := registrar23.AutomationRegistrar23RegistrationParams{
		UpkeepContract: common.HexToAddress(upkeepAddr),
		Amount:         amount,
		AdminAddress:   common.HexToAddress(adminAddr),
		GasLimit:       gasLimit,
		TriggerType:    uint8(0),                           // trigger type
		BillingToken:   common.HexToAddress(wethTokenAddr), // native
		Name:           name,
		EncryptedEmail: email,
		CheckData:      checkData,
		TriggerConfig:  []byte{}, // conditional upkeep
		OffchainConfig: []byte{},
	}

	decodedTx, err := v.client.Decode(v.registrar23.RegisterUpkeep(txOpts,
		params,
	))
	return decodedTx.Transaction, err
}

// EncodeRegisterRequest encodes register request to call it through link token TransferAndCall
func (v *EthereumKeeperRegistrar) EncodeRegisterRequest(name string, email []byte, upkeepAddr string, gasLimit uint32, adminAddr string, checkData []byte, amount *big.Int, source uint8, senderAddr string, isLogTrigger bool, isMercury bool, linkTokenAddr string) ([]byte, error) {
	if v.registrar20 != nil {
		registryABI, err := abi.JSON(strings.NewReader(keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.ABI))
		if err != nil {
			return nil, err
		}
		req, err := registryABI.Pack(
			"register",
			name,
			email,
			common.HexToAddress(upkeepAddr),
			gasLimit,
			common.HexToAddress(adminAddr),
			checkData,
			[]byte{}, //offchainConfig
			amount,
			common.HexToAddress(senderAddr),
		)

		if err != nil {
			return nil, err
		}
		return req, nil
	} else if v.registrar21 != nil {
		if isLogTrigger {
			var topic0InBytes [32]byte
			// bytes representation of 0x0000000000000000000000000000000000000000000000000000000000000000
			bytes0 := [32]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			}
			if isMercury {
				// bytes representation of 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd
				topic0InBytes = [32]byte{209, 255, 233, 228, 85, 129, 193, 29, 125, 159, 46, 213, 247, 82, 23, 205, 75, 233, 248, 183, 238, 230, 175, 15, 109, 3, 244, 109, 229, 57, 86, 205}
			} else {
				// bytes representation of 0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d
				topic0InBytes = [32]byte{
					61, 83, 163, 149, 80, 224, 70, 136,
					6, 88, 39, 243, 187, 134, 88, 76,
					176, 7, 171, 158, 188, 167, 235,
					213, 40, 231, 48, 28, 156, 49, 235, 93,
				}
			}

			logTriggerConfigStruct := acutils.IAutomationV21PlusCommonLogTriggerConfig{
				ContractAddress: common.HexToAddress(upkeepAddr),
				FilterSelector:  0,
				Topic0:          topic0InBytes,
				Topic1:          bytes0,
				Topic2:          bytes0,
				Topic3:          bytes0,
			}
			encodedLogTriggerConfig, err := compatibleUtils.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
			if err != nil {
				return nil, err
			}

			req, err := registrarABI.Pack(
				"register",
				name,
				email,
				common.HexToAddress(upkeepAddr),
				gasLimit,
				common.HexToAddress(adminAddr),
				uint8(1), // trigger type
				checkData,
				encodedLogTriggerConfig, // triggerConfig
				[]byte{},                // offchainConfig
				amount,
				common.HexToAddress(senderAddr),
			)

			return req, err
		}
		req, err := registrarABI.Pack(
			"register",
			name,
			email,
			common.HexToAddress(upkeepAddr),
			gasLimit,
			common.HexToAddress(adminAddr),
			uint8(0), // trigger type
			checkData,
			[]byte{}, // triggerConfig
			[]byte{}, // offchainConfig
			amount,
			common.HexToAddress(senderAddr),
		)
		return req, err
	} else if v.registrar23 != nil {
		registrarABI = cltypes.MustGetABI(registrar23.AutomationRegistrarABI)

		if isLogTrigger {
			var topic0InBytes [32]byte
			// bytes representation of 0x0000000000000000000000000000000000000000000000000000000000000000
			bytes0 := [32]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			}
			if isMercury {
				// bytes representation of 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd
				topic0InBytes = [32]byte{209, 255, 233, 228, 85, 129, 193, 29, 125, 159, 46, 213, 247, 82, 23, 205, 75, 233, 248, 183, 238, 230, 175, 15, 109, 3, 244, 109, 229, 57, 86, 205}
			} else {
				// bytes representation of 0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d
				topic0InBytes = [32]byte{
					61, 83, 163, 149, 80, 224, 70, 136,
					6, 88, 39, 243, 187, 134, 88, 76,
					176, 7, 171, 158, 188, 167, 235,
					213, 40, 231, 48, 28, 156, 49, 235, 93,
				}
			}

			logTriggerConfigStruct := acutils.IAutomationV21PlusCommonLogTriggerConfig{
				ContractAddress: common.HexToAddress(upkeepAddr),
				FilterSelector:  0,
				Topic0:          topic0InBytes,
				Topic1:          bytes0,
				Topic2:          bytes0,
				Topic3:          bytes0,
			}
			encodedLogTriggerConfig, err := compatibleUtils.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
			if err != nil {
				return nil, err
			}

			params := registrar23.AutomationRegistrar23RegistrationParams{
				UpkeepContract: common.HexToAddress(upkeepAddr),
				Amount:         amount,
				AdminAddress:   common.HexToAddress(adminAddr),
				GasLimit:       gasLimit,
				TriggerType:    uint8(1), // trigger type
				BillingToken:   common.HexToAddress(linkTokenAddr),
				Name:           name,
				EncryptedEmail: email,
				CheckData:      checkData,
				TriggerConfig:  encodedLogTriggerConfig,
				OffchainConfig: []byte{},
			}

			req, err := registrarABI.Methods["registerUpkeep"].Inputs.Pack(&params)
			return req, err
		}

		params := registrar23.AutomationRegistrar23RegistrationParams{
			UpkeepContract: common.HexToAddress(upkeepAddr),
			Amount:         amount,
			AdminAddress:   common.HexToAddress(adminAddr),
			GasLimit:       gasLimit,
			TriggerType:    uint8(0), // trigger type
			BillingToken:   common.HexToAddress(linkTokenAddr),
			Name:           name,
			EncryptedEmail: email,
			CheckData:      checkData,
			TriggerConfig:  []byte{},
			OffchainConfig: []byte{},
		}

		encodedRegistrationParamsStruct, err := registrarABI.Methods["registerUpkeep"].Inputs.Pack(&params)

		return encodedRegistrationParamsStruct, err
	}
	registryABI, err := abi.JSON(strings.NewReader(keeper_registrar_wrapper1_2.KeeperRegistrarMetaData.ABI))
	if err != nil {
		return nil, err
	}
	req, err := registryABI.Pack(
		"register",
		name,
		email,
		common.HexToAddress(upkeepAddr),
		gasLimit,
		common.HexToAddress(adminAddr),
		checkData,
		amount,
		source,
		common.HexToAddress(senderAddr),
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DeployKeeperRegistrar(client *seth.Client, registryVersion eth_contracts.KeeperRegistryVersion, linkAddr string, registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error) {
	if registryVersion == eth_contracts.RegistryVersion_2_0 {
		abi, err := keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.GetAbi()
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to get KeeperRegistrar2_0 ABI: %w", err)
		}
		data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistrar2_0", *abi, common.FromHex(keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.Bin),
			common.HexToAddress(linkAddr),
			registrarSettings.AutoApproveConfigType,
			registrarSettings.AutoApproveMaxAllowed,
			common.HexToAddress(registrarSettings.RegistryAddr),
			registrarSettings.MinLinkJuels,
		)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("KeeperRegistrar2_0 instance deployment have failed: %w", err)
		}

		instance, err := keeper_registrar_wrapper2_0.NewKeeperRegistrar(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to instantiate KeeperRegistrar2_0 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			client:      client,
			registrar20: instance,
			address:     &data.Address,
		}, nil
	} else if registryVersion == eth_contracts.RegistryVersion_2_1 || registryVersion == eth_contracts.RegistryVersion_2_2 { // both 2.1 and 2.2 registry use registrar 2.1
		abi, err := registrar21.AutomationRegistrarMetaData.GetAbi()
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to get KeeperRegistrar2_1 ABI: %w", err)
		}
		// set default TriggerType to 0(conditional), AutoApproveConfigType to 2(auto approve enabled), AutoApproveMaxAllowed to 1000
		triggerConfigs := []registrar21.AutomationRegistrar21InitialTriggerConfig{
			{TriggerType: 0, AutoApproveType: registrarSettings.AutoApproveConfigType,
				AutoApproveMaxAllowed: uint32(registrarSettings.AutoApproveMaxAllowed)},
			{TriggerType: 1, AutoApproveType: registrarSettings.AutoApproveConfigType,
				AutoApproveMaxAllowed: uint32(registrarSettings.AutoApproveMaxAllowed)},
		}

		data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistrar2_1", *abi, common.FromHex(registrar21.AutomationRegistrarMetaData.Bin),
			common.HexToAddress(linkAddr),
			common.HexToAddress(registrarSettings.RegistryAddr),
			registrarSettings.MinLinkJuels,
			triggerConfigs,
		)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("KeeperRegistrar2_1 instance deployment have failed: %w", err)
		}

		instance, err := registrar21.NewAutomationRegistrar(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to instantiate KeeperRegistrar2_1 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			client:      client,
			registrar21: instance,
			address:     &data.Address,
		}, nil
	} else if registryVersion == eth_contracts.RegistryVersion_2_3 {
		abi, err := registrar23.AutomationRegistrarMetaData.GetAbi()
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to get KeeperRegistrar2_3 ABI: %w", err)
		}
		// set default TriggerType to 0(conditional), AutoApproveConfigType to 2(auto approve enabled), AutoApproveMaxAllowed to 1000
		triggerConfigs := []registrar23.AutomationRegistrar23InitialTriggerConfig{
			{TriggerType: 0, AutoApproveType: registrarSettings.AutoApproveConfigType,
				AutoApproveMaxAllowed: uint32(registrarSettings.AutoApproveMaxAllowed)},
			{TriggerType: 1, AutoApproveType: registrarSettings.AutoApproveConfigType,
				AutoApproveMaxAllowed: uint32(registrarSettings.AutoApproveMaxAllowed)},
		}

		billingTokens := []common.Address{
			common.HexToAddress(linkAddr),
			common.HexToAddress(registrarSettings.WETHTokenAddr),
		}
		minRegistrationFees := []*big.Int{
			big.NewInt(10),
			big.NewInt(10),
		}

		data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistrar2_3", *abi, common.FromHex(registrar23.AutomationRegistrarMetaData.Bin),
			common.HexToAddress(linkAddr),
			common.HexToAddress(registrarSettings.RegistryAddr),
			triggerConfigs,
			billingTokens,
			minRegistrationFees,
			common.HexToAddress(registrarSettings.WETHTokenAddr),
		)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("KeeperRegistrar2_3 instance deployment have failed: %w", err)
		}

		instance, err := registrar23.NewAutomationRegistrar(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to instantiate KeeperRegistrar2_3 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			client:      client,
			registrar23: instance,
			address:     &data.Address,
		}, nil
	}

	// non OCR registrar
	abi, err := keeper_registrar_wrapper1_2.KeeperRegistrarMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to get KeeperRegistrar1_2 ABI: %w", err)
	}

	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistrar1_2", *abi, common.FromHex(keeper_registrar_wrapper1_2.KeeperRegistrarMetaData.Bin),
		common.HexToAddress(linkAddr),
		registrarSettings.AutoApproveConfigType,
		registrarSettings.AutoApproveMaxAllowed,
		common.HexToAddress(registrarSettings.RegistryAddr),
		registrarSettings.MinLinkJuels,
	)
	if err != nil {
		return &EthereumKeeperRegistrar{}, fmt.Errorf("KeeperRegistrar1_2 instance deployment have failed: %w", err)
	}

	instance, err := keeper_registrar_wrapper1_2.NewKeeperRegistrar(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to instantiate KeeperRegistrar1_2 instance: %w", err)
	}

	return &EthereumKeeperRegistrar{
		client:    client,
		registrar: instance,
		address:   &data.Address,
	}, nil
}

// LoadKeeperRegistrar returns deployed on given address EthereumKeeperRegistrar
func LoadKeeperRegistrar(client *seth.Client, address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistrar, error) {
	if registryVersion == eth_contracts.RegistryVersion_1_1 || registryVersion == eth_contracts.RegistryVersion_1_2 ||
		registryVersion == eth_contracts.RegistryVersion_1_3 {

		loader := seth.NewContractLoader[keeper_registrar_wrapper1_2.KeeperRegistrar](client)
		instance, err := loader.LoadContract("KeeperRegistrar1_2", address, keeper_registrar_wrapper1_2.KeeperRegistrarMetaData.GetAbi, keeper_registrar_wrapper1_2.NewKeeperRegistrar)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to load KeeperRegistrar1_2 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			address:   &address,
			client:    client,
			registrar: instance,
		}, err
	} else if registryVersion == eth_contracts.RegistryVersion_2_0 {
		loader := seth.NewContractLoader[keeper_registrar_wrapper2_0.KeeperRegistrar](client)
		instance, err := loader.LoadContract("KeeperRegistrar2_0", address, keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.GetAbi, keeper_registrar_wrapper2_0.NewKeeperRegistrar)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to load KeeperRegistrar2_0 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			address:     &address,
			client:      client,
			registrar20: instance,
		}, nil
	} else if registryVersion == eth_contracts.RegistryVersion_2_1 || registryVersion == eth_contracts.RegistryVersion_2_2 {
		loader := seth.NewContractLoader[registrar21.AutomationRegistrar](client)
		instance, err := loader.LoadContract("KeeperRegistrar2_1", address, registrar21.AutomationRegistrarMetaData.GetAbi, registrar21.NewAutomationRegistrar)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to load KeeperRegistrar2_1 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			address:     &address,
			client:      client,
			registrar21: instance,
		}, nil
	} else if registryVersion == eth_contracts.RegistryVersion_2_3 {
		loader := seth.NewContractLoader[registrar23.AutomationRegistrar](client)
		instance, err := loader.LoadContract("KeeperRegistrar2_3", address, registrar23.AutomationRegistrarMetaData.GetAbi, registrar23.NewAutomationRegistrar)
		if err != nil {
			return &EthereumKeeperRegistrar{}, fmt.Errorf("failed to load KeeperRegistrar2_3 instance: %w", err)
		}

		return &EthereumKeeperRegistrar{
			address:     &address,
			client:      client,
			registrar23: instance,
		}, nil
	}
	return &EthereumKeeperRegistrar{}, fmt.Errorf("unsupported registry version: %v", registryVersion)
}

type EthereumAutomationKeeperConsumer struct {
	client   *seth.Client
	consumer *log_upkeep_counter_wrapper.LogUpkeepCounter
	address  *common.Address
}

func (e EthereumAutomationKeeperConsumer) Address() string {
	return e.address.Hex()
}

func (e EthereumAutomationKeeperConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return e.consumer.Counter(&bind.CallOpts{
		From:    e.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (e EthereumAutomationKeeperConsumer) Start() error {
	_, err := e.client.Decode(e.consumer.Start(e.client.NewTXOpts()))
	return err
}

func LoadKeeperConsumer(client *seth.Client, address common.Address) (*EthereumAutomationKeeperConsumer, error) {
	loader := seth.NewContractLoader[log_upkeep_counter_wrapper.LogUpkeepCounter](client)
	instance, err := loader.LoadContract("KeeperConsumer", address, log_upkeep_counter_wrapper.LogUpkeepCounterMetaData.GetAbi, log_upkeep_counter_wrapper.NewLogUpkeepCounter)
	if err != nil {
		return &EthereumAutomationKeeperConsumer{}, fmt.Errorf("failed to load KeeperConsumerMetaData instance: %w", err)
	}

	return &EthereumAutomationKeeperConsumer{
		client:   client,
		consumer: instance,
		address:  &address,
	}, nil
}

type EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer struct {
	client   *seth.Client
	consumer *log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookup
	address  *common.Address
}

func (v *EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer) Address() string {
	return v.address.Hex()
}

// Kick off the log trigger event. The contract uses Mercury v0.2 so no need to set ParamKeys
func (v *EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer) Start() error {
	_, err := v.client.Decode(v.consumer.Start(v.client.NewTXOpts()))
	return err
}

func (v *EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func DeployAutomationLogTriggeredStreamsLookupUpkeepConsumerFromKey(client *seth.Client, keyNum int) (KeeperConsumer, error) {
	abi, err := log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookupMetaData.GetAbi()
	if err != nil {
		return &EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer{}, fmt.Errorf("failed to get LogTriggeredStreamsLookupUpkeep ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "LogTriggeredStreamsLookupUpkeep", *abi, common.FromHex(log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookupMetaData.Bin), false, false, false)
	if err != nil {
		return &EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer{}, fmt.Errorf("LogTriggeredStreamsLookupUpkeep instance deployment have failed: %w", err)
	}

	instance, err := log_triggered_streams_lookup_wrapper.NewLogTriggeredStreamsLookup(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer{}, fmt.Errorf("failed to instantiate LogTriggeredStreamsLookupUpkeep instance: %w", err)
	}

	return &EthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

func DeployAutomationLogTriggeredStreamsLookupUpkeepConsumer(client *seth.Client) (KeeperConsumer, error) {
	return DeployAutomationLogTriggeredStreamsLookupUpkeepConsumerFromKey(client, 0)
}

type EthereumAutomationStreamsLookupUpkeepConsumer struct {
	client   *seth.Client
	consumer *streams_lookup_upkeep_wrapper.StreamsLookupUpkeep
	address  *common.Address
}

func (v *EthereumAutomationStreamsLookupUpkeepConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumAutomationStreamsLookupUpkeepConsumer) Start() error {
	_, err := v.client.Decode(v.consumer.SetParamKeys(v.client.NewTXOpts(), "feedIdHex", "blockNumber"))
	return err
}

func (v *EthereumAutomationStreamsLookupUpkeepConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func DeployAutomationStreamsLookupUpkeepConsumerFromKey(client *seth.Client, keyNum int, testRange *big.Int, interval *big.Int, useArbBlock bool, staging bool, verify bool) (KeeperConsumer, error) {
	abi, err := streams_lookup_upkeep_wrapper.StreamsLookupUpkeepMetaData.GetAbi()
	if err != nil {
		return &EthereumAutomationStreamsLookupUpkeepConsumer{}, fmt.Errorf("failed to get StreamsLookupUpkeep ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "StreamsLookupUpkeep", *abi, common.FromHex(streams_lookup_upkeep_wrapper.StreamsLookupUpkeepMetaData.Bin),
		testRange,
		interval,
		useArbBlock,
		staging,
		verify,
	)
	if err != nil {
		return &EthereumAutomationStreamsLookupUpkeepConsumer{}, fmt.Errorf("StreamsLookupUpkeep instance deployment have failed: %w", err)
	}

	instance, err := streams_lookup_upkeep_wrapper.NewStreamsLookupUpkeep(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumAutomationStreamsLookupUpkeepConsumer{}, fmt.Errorf("failed to instantiate StreamsLookupUpkeep instance: %w", err)
	}

	return &EthereumAutomationStreamsLookupUpkeepConsumer{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

func DeployAutomationStreamsLookupUpkeepConsumer(client *seth.Client, testRange *big.Int, interval *big.Int, useArbBlock bool, staging bool, verify bool) (KeeperConsumer, error) {
	return DeployAutomationStreamsLookupUpkeepConsumerFromKey(client, 0, testRange, interval, useArbBlock, staging, verify)
}

type EthereumAutomationLogCounterConsumer struct {
	client   *seth.Client
	consumer *log_upkeep_counter_wrapper.LogUpkeepCounter
	address  *common.Address
}

func (v *EthereumAutomationLogCounterConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumAutomationLogCounterConsumer) Start() error {
	_, err := v.client.Decode(v.consumer.Start(v.client.NewTXOpts()))
	return err
}

func (v *EthereumAutomationLogCounterConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func DeployAutomationLogTriggerConsumerFromKey(client *seth.Client, keyNum int, testInterval *big.Int) (KeeperConsumer, error) {
	abi, err := log_upkeep_counter_wrapper.LogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return &EthereumAutomationLogCounterConsumer{}, fmt.Errorf("failed to get LogUpkeepCounter ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "LogUpkeepCounter", *abi, common.FromHex(log_upkeep_counter_wrapper.LogUpkeepCounterMetaData.Bin), testInterval)
	if err != nil {
		return &EthereumAutomationLogCounterConsumer{}, fmt.Errorf("LogUpkeepCounter instance deployment have failed: %w", err)
	}

	instance, err := log_upkeep_counter_wrapper.NewLogUpkeepCounter(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumAutomationLogCounterConsumer{}, fmt.Errorf("failed to instantiate LogUpkeepCounter instance: %w", err)
	}

	return &EthereumAutomationLogCounterConsumer{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

func DeployAutomationLogTriggerConsumer(client *seth.Client, testInterval *big.Int) (KeeperConsumer, error) {
	return DeployAutomationLogTriggerConsumerFromKey(client, 0, testInterval)
}

// EthereumUpkeepCounter represents keeper consumer (upkeep) counter contract
type EthereumUpkeepCounter struct {
	client   *seth.Client
	consumer *upkeep_counter_wrapper.UpkeepCounter
	address  *common.Address
}

func (v *EthereumUpkeepCounter) Address() string {
	return v.address.Hex()
}

func (v *EthereumUpkeepCounter) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}
func (v *EthereumUpkeepCounter) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumUpkeepCounter) SetSpread(testRange *big.Int, interval *big.Int) error {
	_, err := v.client.Decode(v.consumer.SetSpread(v.client.NewTXOpts(), testRange, interval))
	return err
}

// Just pass for non-logtrigger
func (v *EthereumUpkeepCounter) Start() error {
	return nil
}

func DeployUpkeepCounterFromKey(client *seth.Client, keyNum int, testRange *big.Int, interval *big.Int) (UpkeepCounter, error) {
	abi, err := upkeep_counter_wrapper.UpkeepCounterMetaData.GetAbi()
	if err != nil {
		return &EthereumUpkeepCounter{}, fmt.Errorf("failed to get UpkeepCounter ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "UpkeepCounter", *abi, common.FromHex(upkeep_counter_wrapper.UpkeepCounterMetaData.Bin), testRange, interval)
	if err != nil {
		return &EthereumUpkeepCounter{}, fmt.Errorf("UpkeepCounter instance deployment have failed: %w", err)
	}

	instance, err := upkeep_counter_wrapper.NewUpkeepCounter(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumUpkeepCounter{}, fmt.Errorf("failed to instantiate UpkeepCounter instance: %w", err)
	}

	return &EthereumUpkeepCounter{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

func DeployUpkeepCounter(client *seth.Client, testRange *big.Int, interval *big.Int) (UpkeepCounter, error) {
	return DeployUpkeepCounterFromKey(client, 0, testRange, interval)
}

// EthereumUpkeepPerformCounterRestrictive represents keeper consumer (upkeep) counter contract
type EthereumUpkeepPerformCounterRestrictive struct {
	client   *seth.Client
	consumer *upkeep_perform_counter_restrictive_wrapper.UpkeepPerformCounterRestrictive
	address  *common.Address
}

func (v *EthereumUpkeepPerformCounterRestrictive) Address() string {
	return v.address.Hex()
}

func (v *EthereumUpkeepPerformCounterRestrictive) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}
func (v *EthereumUpkeepPerformCounterRestrictive) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.GetCountPerforms(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumUpkeepPerformCounterRestrictive) SetSpread(testRange *big.Int, interval *big.Int) error {
	_, err := v.client.Decode(v.consumer.SetSpread(v.client.NewTXOpts(), testRange, interval))
	return err
}

func DeployUpkeepPerformCounterRestrictive(client *seth.Client, testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error) {
	abi, err := upkeep_perform_counter_restrictive_wrapper.UpkeepPerformCounterRestrictiveMetaData.GetAbi()
	if err != nil {
		return &EthereumUpkeepCounter{}, fmt.Errorf("failed to get UpkeepPerformCounterRestrictive ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "UpkeepPerformCounterRestrictive", *abi, common.FromHex(upkeep_perform_counter_restrictive_wrapper.UpkeepPerformCounterRestrictiveMetaData.Bin), testRange, averageEligibilityCadence)
	if err != nil {
		return &EthereumUpkeepCounter{}, fmt.Errorf("UpkeepPerformCounterRestrictive instance deployment have failed: %w", err)
	}

	instance, err := upkeep_perform_counter_restrictive_wrapper.NewUpkeepPerformCounterRestrictive(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumUpkeepCounter{}, fmt.Errorf("failed to instantiate UpkeepPerformCounterRestrictive instance: %w", err)
	}

	return &EthereumUpkeepPerformCounterRestrictive{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

// EthereumKeeperPerformDataCheckerConsumer represents keeper perform data checker contract
type EthereumKeeperPerformDataCheckerConsumer struct {
	client             *seth.Client
	performDataChecker *perform_data_checker_wrapper.PerformDataChecker
	address            *common.Address
}

func (v *EthereumKeeperPerformDataCheckerConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperPerformDataCheckerConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return v.performDataChecker.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumKeeperPerformDataCheckerConsumer) SetExpectedData(_ context.Context, expectedData []byte) error {
	_, err := v.client.Decode(v.performDataChecker.SetExpectedData(v.client.NewTXOpts(), expectedData))
	return err
}

func DeployKeeperPerformDataChecker(client *seth.Client, expectedData []byte) (KeeperPerformDataChecker, error) {
	abi, err := perform_data_checker_wrapper.PerformDataCheckerMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperPerformDataCheckerConsumer{}, fmt.Errorf("failed to get PerformDataChecker ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "PerformDataChecker", *abi, common.FromHex(perform_data_checker_wrapper.PerformDataCheckerMetaData.Bin), expectedData)
	if err != nil {
		return &EthereumKeeperPerformDataCheckerConsumer{}, fmt.Errorf("PerformDataChecker instance deployment have failed: %w", err)
	}

	instance, err := perform_data_checker_wrapper.NewPerformDataChecker(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperPerformDataCheckerConsumer{}, fmt.Errorf("failed to instantiate PerformDataChecker instance: %w", err)
	}

	return &EthereumKeeperPerformDataCheckerConsumer{
		client:             client,
		performDataChecker: instance,
		address:            &data.Address,
	}, nil
}

// EthereumKeeperConsumerPerformance represents a more complicated keeper consumer contract, one intended only for
// performance tests.
type EthereumKeeperConsumerPerformance struct {
	client   *seth.Client
	consumer *keeper_consumer_performance_wrapper.KeeperConsumerPerformance
	address  *common.Address
}

func (v *EthereumKeeperConsumerPerformance) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperConsumerPerformance) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}

func (v *EthereumKeeperConsumerPerformance) CheckEligible(ctx context.Context) (bool, error) {
	return v.consumer.CheckEligible(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumKeeperConsumerPerformance) GetUpkeepCount(ctx context.Context) (*big.Int, error) {
	return v.consumer.GetCountPerforms(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumKeeperConsumerPerformance) SetCheckGasToBurn(_ context.Context, gas *big.Int) error {
	_, err := v.client.Decode(v.consumer.SetCheckGasToBurn(v.client.NewTXOpts(), gas))
	return err
}

func (v *EthereumKeeperConsumerPerformance) SetPerformGasToBurn(_ context.Context, gas *big.Int) error {
	_, err := v.client.Decode(v.consumer.SetPerformGasToBurn(v.client.NewTXOpts(), gas))
	return err
}

func DeployKeeperConsumerPerformance(
	client *seth.Client,
	testBlockRange,
	averageCadence,
	checkGasToBurn,
	performGasToBurn *big.Int,
) (KeeperConsumerPerformance, error) {
	abi, err := keeper_consumer_performance_wrapper.KeeperConsumerPerformanceMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperConsumerPerformance{}, fmt.Errorf("failed to get KeeperConsumerPerformance ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "KeeperConsumerPerformance", *abi, common.FromHex(keeper_consumer_performance_wrapper.KeeperConsumerPerformanceMetaData.Bin),
		testBlockRange,
		averageCadence,
		checkGasToBurn,
		performGasToBurn)
	if err != nil {
		return &EthereumKeeperConsumerPerformance{}, fmt.Errorf("KeeperConsumerPerformance instance deployment have failed: %w", err)
	}

	instance, err := keeper_consumer_performance_wrapper.NewKeeperConsumerPerformance(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperConsumerPerformance{}, fmt.Errorf("failed to instantiate KeeperConsumerPerformance instance: %w", err)
	}

	return &EthereumKeeperConsumerPerformance{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

type EthereumAutomationSimpleLogCounterConsumer struct {
	client   *seth.Client
	consumer *simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounter
	address  *common.Address
}

func (v *EthereumAutomationSimpleLogCounterConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumAutomationSimpleLogCounterConsumer) Start() error {
	return nil
}

func (v *EthereumAutomationSimpleLogCounterConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func DeployAutomationSimpleLogTriggerConsumer(client *seth.Client, isStreamsLookup bool) (KeeperConsumer, error) {
	return DeployAutomationSimpleLogTriggerConsumerFromKey(client, isStreamsLookup, 0)
}

func DeployAutomationSimpleLogTriggerConsumerFromKey(client *seth.Client, isStreamsLookup bool, keyNum int) (KeeperConsumer, error) {
	abi, err := simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return &EthereumAutomationSimpleLogCounterConsumer{}, fmt.Errorf("failed to get SimpleLogUpkeepCounter ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "SimpleLogUpkeepCounter", *abi, common.FromHex(simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.Bin), isStreamsLookup)
	if err != nil {
		return &EthereumAutomationSimpleLogCounterConsumer{}, fmt.Errorf("SimpleLogUpkeepCounter instance deployment have failed: %w", err)
	}

	instance, err := simple_log_upkeep_counter_wrapper.NewSimpleLogUpkeepCounter(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumAutomationSimpleLogCounterConsumer{}, fmt.Errorf("failed to instantiate SimpleLogUpkeepCounter instance: %w", err)
	}

	return &EthereumAutomationSimpleLogCounterConsumer{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

// EthereumAutomationConsumerBenchmark represents a more complicated keeper consumer contract, one intended only for
// Benchmark tests.
type EthereumAutomationConsumerBenchmark struct {
	client   *seth.Client
	consumer *automation_consumer_benchmark.AutomationConsumerBenchmark
	address  *common.Address
}

func (v *EthereumAutomationConsumerBenchmark) Address() string {
	return v.address.Hex()
}

func (v *EthereumAutomationConsumerBenchmark) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}

func (v *EthereumAutomationConsumerBenchmark) CheckEligible(ctx context.Context, id *big.Int, _range *big.Int, firstEligibleBuffer *big.Int) (bool, error) {
	return v.consumer.CheckEligible(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, id, _range, firstEligibleBuffer)
}

func (v *EthereumAutomationConsumerBenchmark) GetUpkeepCount(ctx context.Context, id *big.Int) (*big.Int, error) {
	return v.consumer.GetCountPerforms(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, id)
}

// DeployAutomationConsumerBenchmark deploys a keeper consumer benchmark contract with a standard contract backend
func DeployAutomationConsumerBenchmark(client *seth.Client) (AutomationConsumerBenchmark, error) {
	return deployAutomationConsumerBenchmarkWithWrapperFn(client, func(client *seth.Client) *wrappers.WrappedContractBackend {
		return wrappers.MustNewWrappedContractBackend(nil, client)
	})
}

func LoadAutomationConsumerBenchmark(client *seth.Client, address common.Address) (*EthereumAutomationConsumerBenchmark, error) {
	loader := seth.NewContractLoader[automation_consumer_benchmark.AutomationConsumerBenchmark](client)
	instance, err := loader.LoadContract("AutomationConsumerBenchmark", address, automation_consumer_benchmark.AutomationConsumerBenchmarkMetaData.GetAbi, automation_consumer_benchmark.NewAutomationConsumerBenchmark)
	if err != nil {
		return &EthereumAutomationConsumerBenchmark{}, fmt.Errorf("failed to load AutomationConsumerBenchmark instance: %w", err)
	}

	return &EthereumAutomationConsumerBenchmark{
		client:   client,
		consumer: instance,
		address:  &address,
	}, nil
}

// DeployAutomationConsumerBenchmarkWithRetry deploys a keeper consumer benchmark contract with a read-only operations retrying contract backend
func DeployAutomationConsumerBenchmarkWithRetry(client *seth.Client, logger zerolog.Logger, maxAttempts uint, retryDelay time.Duration) (AutomationConsumerBenchmark, error) {
	return deployAutomationConsumerBenchmarkWithWrapperFn(client, func(client *seth.Client) *wrappers.WrappedContractBackend {
		return wrappers.MustNewRetryingWrappedContractBackend(client, logger, maxAttempts, retryDelay)
	})
}

func deployAutomationConsumerBenchmarkWithWrapperFn(client *seth.Client, wrapperConstrFn func(client *seth.Client) *wrappers.WrappedContractBackend) (AutomationConsumerBenchmark, error) {
	abi, err := automation_consumer_benchmark.AutomationConsumerBenchmarkMetaData.GetAbi()
	if err != nil {
		return &EthereumAutomationConsumerBenchmark{}, fmt.Errorf("failed to get AutomationConsumerBenchmark ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "AutomationConsumerBenchmark", *abi, common.FromHex(automation_consumer_benchmark.AutomationConsumerBenchmarkMetaData.Bin))
	if err != nil {
		return &EthereumAutomationConsumerBenchmark{}, fmt.Errorf("AutomationConsumerBenchmark instance deployment have failed: %w", err)
	}

	instance, err := automation_consumer_benchmark.NewAutomationConsumerBenchmark(data.Address, wrapperConstrFn(client))
	if err != nil {
		return &EthereumAutomationConsumerBenchmark{}, fmt.Errorf("failed to instantiate AutomationConsumerBenchmark instance: %w", err)
	}

	return &EthereumAutomationConsumerBenchmark{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

// AutomationConsumerBenchmarkUpkeepObserver is a header subscription that awaits for a round of upkeeps
type AutomationConsumerBenchmarkUpkeepObserver struct {
	instance AutomationConsumerBenchmark
	registry KeeperRegistry
	upkeepID *big.Int

	firstBlockNum       uint64                                     // Records the number of the first block that came in
	lastBlockNum        uint64                                     // Records the number of the last block that came in
	blockRange          int64                                      // How many blocks to watch upkeeps for
	upkeepSLA           int64                                      // SLA after which an upkeep is counted as 'missed'
	metricsReporter     *testreporters.KeeperBenchmarkTestReporter // Testreporter to track results
	upkeepIndex         int64
	firstEligibleBuffer int64

	// State variables, changes as we get blocks
	blocksSinceSubscription int64   // How many blocks have passed since subscribing
	blocksSinceEligible     int64   // How many blocks have come in since upkeep has been eligible for check
	countEligible           int64   // Number of times the upkeep became eligible
	countMissed             int64   // Number of times we missed SLA for performing upkeep
	upkeepCount             int64   // The count of upkeeps done so far
	allCheckDelays          []int64 // Tracks the amount of blocks missed before an upkeep since it became eligible
	complete                bool
	l                       zerolog.Logger
}

// NewAutomationConsumerBenchmarkUpkeepObserver provides a new instance of a NewAutomationConsumerBenchmarkUpkeepObserver
// Used to track and log benchmark test results for keepers
func NewAutomationConsumerBenchmarkUpkeepObserver(
	contract AutomationConsumerBenchmark,
	registry KeeperRegistry,
	upkeepID *big.Int,
	blockRange int64,
	upkeepSLA int64,
	metricsReporter *testreporters.KeeperBenchmarkTestReporter,
	upkeepIndex int64,
	firstEligibleBuffer int64,
	logger zerolog.Logger,
) *AutomationConsumerBenchmarkUpkeepObserver {
	return &AutomationConsumerBenchmarkUpkeepObserver{
		instance:                contract,
		registry:                registry,
		upkeepID:                upkeepID,
		blockRange:              blockRange,
		upkeepSLA:               upkeepSLA,
		blocksSinceSubscription: 0,
		blocksSinceEligible:     0,
		upkeepCount:             0,
		allCheckDelays:          []int64{},
		metricsReporter:         metricsReporter,
		complete:                false,
		lastBlockNum:            0,
		upkeepIndex:             upkeepIndex,
		firstBlockNum:           0,
		firstEligibleBuffer:     firstEligibleBuffer,
		l:                       logger,
	}
}

// ReceiveHeader will query the latest Keeper round and check to see whether upkeep was performed, it returns
// true when observation has finished.
func (o *AutomationConsumerBenchmarkUpkeepObserver) ReceiveHeader(receivedHeader *blockchain.SafeEVMHeader) (bool, error) {
	if receivedHeader.Number.Uint64() <= o.lastBlockNum { // Uncle / reorg we won't count
		return false, nil
	}
	if o.firstBlockNum == 0 {
		o.firstBlockNum = receivedHeader.Number.Uint64()
	}
	o.lastBlockNum = receivedHeader.Number.Uint64()
	// Increment block counters
	o.blocksSinceSubscription++

	upkeepCount, err := o.instance.GetUpkeepCount(context.Background(), big.NewInt(o.upkeepIndex))
	if err != nil {
		return false, err
	}

	if upkeepCount.Int64() > o.upkeepCount { // A new upkeep was done
		if upkeepCount.Int64() != o.upkeepCount+1 {
			return false, errors.New("upkeep count increased by more than 1 in a single block")
		}
		o.l.Info().
			Uint64("Block_Number", receivedHeader.Number.Uint64()).
			Str("Upkeep_ID", o.upkeepID.String()).
			Str("Contract_Address", o.instance.Address()).
			Int64("Upkeep_Count", upkeepCount.Int64()).
			Int64("Blocks_since_eligible", o.blocksSinceEligible).
			Str("Registry_Address", o.registry.Address()).
			Msg("Upkeep Performed")

		if o.blocksSinceEligible > o.upkeepSLA {
			o.l.Warn().
				Uint64("Block_Number", receivedHeader.Number.Uint64()).
				Str("Upkeep_ID", o.upkeepID.String()).
				Str("Contract_Address", o.instance.Address()).
				Int64("Blocks_since_eligible", o.blocksSinceEligible).
				Str("Registry_Address", o.registry.Address()).
				Msg("Upkeep Missed SLA")
			o.countMissed++
		}

		o.allCheckDelays = append(o.allCheckDelays, o.blocksSinceEligible)
		o.upkeepCount++
		o.blocksSinceEligible = 0
	}

	isEligible, err := o.instance.CheckEligible(context.Background(), big.NewInt(o.upkeepIndex), big.NewInt(o.blockRange), big.NewInt(o.firstEligibleBuffer))
	if err != nil {
		return false, err
	}
	if isEligible {
		if o.blocksSinceEligible == 0 {
			// First time this upkeep became eligible
			o.countEligible++
			o.l.Info().
				Uint64("Block_Number", receivedHeader.Number.Uint64()).
				Str("Upkeep_ID", o.upkeepID.String()).
				Str("Contract_Address", o.instance.Address()).
				Str("Registry_Address", o.registry.Address()).
				Msg("Upkeep Now Eligible")
		}
		o.blocksSinceEligible++
	}

	if o.blocksSinceSubscription >= o.blockRange || int64(o.lastBlockNum-o.firstBlockNum) >= o.blockRange {
		if o.blocksSinceEligible > 0 {
			if o.blocksSinceEligible > o.upkeepSLA {
				o.l.Warn().
					Uint64("Block_Number", receivedHeader.Number.Uint64()).
					Str("Upkeep_ID", o.upkeepID.String()).
					Str("Contract_Address", o.instance.Address()).
					Int64("Blocks_since_eligible", o.blocksSinceEligible).
					Str("Registry_Address", o.registry.Address()).
					Msg("Upkeep remained eligible at end of test and missed SLA")
				o.countMissed++
			} else {
				o.l.Info().
					Uint64("Block_Number", receivedHeader.Number.Uint64()).
					Str("Upkeep_ID", o.upkeepID.String()).
					Str("Contract_Address", o.instance.Address()).
					Int64("Upkeep_Count", upkeepCount.Int64()).
					Int64("Blocks_since_eligible", o.blocksSinceEligible).
					Str("Registry_Address", o.registry.Address()).
					Msg("Upkeep remained eligible at end of test and was within SLA")
			}
			o.allCheckDelays = append(o.allCheckDelays, o.blocksSinceEligible)
		}

		o.l.Info().
			Uint64("Block_Number", receivedHeader.Number.Uint64()).
			Str("Upkeep_ID", o.upkeepID.String()).
			Str("Contract_Address", o.instance.Address()).
			Int64("Upkeeps_Performed", upkeepCount.Int64()).
			Int64("Total_Blocks_Watched", o.blocksSinceSubscription).
			Str("Registry_Address", o.registry.Address()).
			Msg("Finished Watching for Upkeeps")

		o.complete = true
		return true, nil
	}
	return false, nil
}

// Complete returns whether watching for upkeeps has completed
func (o *AutomationConsumerBenchmarkUpkeepObserver) Complete() bool {
	return o.complete
}

// LogDetails logs the results of the benchmark test to testreporter
func (o *AutomationConsumerBenchmarkUpkeepObserver) LogDetails() {
	report := testreporters.KeeperBenchmarkTestReport{
		ContractAddress:       o.instance.Address(),
		TotalEligibleCount:    o.countEligible,
		TotalSLAMissedUpkeeps: o.countMissed,
		TotalPerformedUpkeeps: o.upkeepCount,
		AllCheckDelays:        o.allCheckDelays,
		RegistryAddress:       o.registry.Address(),
	}
	o.metricsReporter.ReportMutex.Lock()
	o.metricsReporter.Reports = append(o.metricsReporter.Reports, report)
	defer o.metricsReporter.ReportMutex.Unlock()
}
